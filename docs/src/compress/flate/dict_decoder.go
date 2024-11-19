<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/compress/flate/dict_decoder.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../index.html">GoDoc</a></div>
<a href="dict_decoder.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
<form method="GET" action="http://localhost:8080/search">
<div id="menu">

<span class="search-box"><input type="search" id="search" name="q" placeholder="Search" aria-label="Search" required><button type="submit"><span><!-- magnifying glass: --><svg width="24" height="24" viewBox="0 0 24 24"><title>submit search</title><path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/><path d="M0 0h24v24H0z" fill="none"/></svg></span></button></span>
</div>
</form>

</div></div>



<div id="page" class="wide">
<div class="container">


  <h1>
    Source file
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/compress">compress</a>/<a href="http://localhost:8080/src/compress/flate">flate</a>/<span class="text-muted">dict_decoder.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/compress/flate">compress/flate</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2016 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package flate
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// dictDecoder implements the LZ77 sliding dictionary as used in decompression.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// LZ77 decompresses data through sequences of two forms of commands:</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//   - Literal insertions: Runs of one or more symbols are inserted into the data</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//     stream as is. This is accomplished through the writeByte method for a</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//     single symbol, or combinations of writeSlice/writeMark for multiple symbols.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//     Any valid stream must start with a literal insertion if no preset dictionary</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//     is used.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//   - Backward copies: Runs of one or more symbols are copied from previously</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//     emitted data. Backward copies come as the tuple (dist, length) where dist</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//     determines how far back in the stream to copy from and length determines how</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//     many bytes to copy. Note that it is valid for the length to be greater than</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//     the distance. Since LZ77 uses forward copies, that situation is used to</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//     perform a form of run-length encoding on repeated runs of symbols.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//     The writeCopy and tryWriteCopy are used to implement this command.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// For performance reasons, this implementation performs little to no sanity</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// checks about the arguments. As such, the invariants documented for each</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// method call must be respected.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>type dictDecoder struct {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	hist []byte <span class="comment">// Sliding window history</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// Invariant: 0 &lt;= rdPos &lt;= wrPos &lt;= len(hist)</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	wrPos int  <span class="comment">// Current output position in buffer</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	rdPos int  <span class="comment">// Have emitted hist[:rdPos] already</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	full  bool <span class="comment">// Has a full window length been written yet?</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// init initializes dictDecoder to have a sliding window dictionary of the given</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// size. If a preset dict is provided, it will initialize the dictionary with</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// the contents of dict.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>func (dd *dictDecoder) init(size int, dict []byte) {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	*dd = dictDecoder{hist: dd.hist}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	if cap(dd.hist) &lt; size {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		dd.hist = make([]byte, size)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	dd.hist = dd.hist[:size]
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	if len(dict) &gt; len(dd.hist) {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		dict = dict[len(dict)-len(dd.hist):]
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	dd.wrPos = copy(dd.hist, dict)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	if dd.wrPos == len(dd.hist) {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		dd.wrPos = 0
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		dd.full = true
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	dd.rdPos = dd.wrPos
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// histSize reports the total amount of historical data in the dictionary.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>func (dd *dictDecoder) histSize() int {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if dd.full {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return len(dd.hist)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	return dd.wrPos
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// availRead reports the number of bytes that can be flushed by readFlush.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>func (dd *dictDecoder) availRead() int {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	return dd.wrPos - dd.rdPos
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// availWrite reports the available amount of output buffer space.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func (dd *dictDecoder) availWrite() int {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	return len(dd.hist) - dd.wrPos
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// writeSlice returns a slice of the available buffer to write data to.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// This invariant will be kept: len(s) &lt;= availWrite()</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func (dd *dictDecoder) writeSlice() []byte {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	return dd.hist[dd.wrPos:]
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// writeMark advances the writer pointer by cnt.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// This invariant must be kept: 0 &lt;= cnt &lt;= availWrite()</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>func (dd *dictDecoder) writeMark(cnt int) {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	dd.wrPos += cnt
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// writeByte writes a single byte to the dictionary.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// This invariant must be kept: 0 &lt; availWrite()</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>func (dd *dictDecoder) writeByte(c byte) {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	dd.hist[dd.wrPos] = c
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	dd.wrPos++
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// writeCopy copies a string at a given (dist, length) to the output.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// This returns the number of bytes copied and may be less than the requested</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// length if the available space in the output buffer is too small.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// This invariant must be kept: 0 &lt; dist &lt;= histSize()</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>func (dd *dictDecoder) writeCopy(dist, length int) int {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	dstBase := dd.wrPos
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	dstPos := dstBase
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	srcPos := dstPos - dist
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	endPos := dstPos + length
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	if endPos &gt; len(dd.hist) {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		endPos = len(dd.hist)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// Copy non-overlapping section after destination position.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// This section is non-overlapping in that the copy length for this section</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// is always less than or equal to the backwards distance. This can occur</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// if a distance refers to data that wraps-around in the buffer.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// Thus, a backwards copy is performed here; that is, the exact bytes in</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// the source prior to the copy is placed in the destination.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	if srcPos &lt; 0 {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		srcPos += len(dd.hist)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		dstPos += copy(dd.hist[dstPos:endPos], dd.hist[srcPos:])
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		srcPos = 0
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// Copy possibly overlapping section before destination position.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// This section can overlap if the copy length for this section is larger</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// than the backwards distance. This is allowed by LZ77 so that repeated</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// strings can be succinctly represented using (dist, length) pairs.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// Thus, a forwards copy is performed here; that is, the bytes copied is</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// possibly dependent on the resulting bytes in the destination as the copy</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// progresses along. This is functionally equivalent to the following:</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">//	for i := 0; i &lt; endPos-dstPos; i++ {</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">//		dd.hist[dstPos+i] = dd.hist[srcPos+i]</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">//	}</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">//	dstPos = endPos</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	for dstPos &lt; endPos {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		dstPos += copy(dd.hist[dstPos:endPos], dd.hist[srcPos:dstPos])
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	dd.wrPos = dstPos
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	return dstPos - dstBase
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// tryWriteCopy tries to copy a string at a given (distance, length) to the</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// output. This specialized version is optimized for short distances.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// This method is designed to be inlined for performance reasons.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">// This invariant must be kept: 0 &lt; dist &lt;= histSize()</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>func (dd *dictDecoder) tryWriteCopy(dist, length int) int {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	dstPos := dd.wrPos
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	endPos := dstPos + length
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if dstPos &lt; dist || endPos &gt; len(dd.hist) {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return 0
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	dstBase := dstPos
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	srcPos := dstPos - dist
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// Copy possibly overlapping section before destination position.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	for dstPos &lt; endPos {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		dstPos += copy(dd.hist[dstPos:endPos], dd.hist[srcPos:dstPos])
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	dd.wrPos = dstPos
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	return dstPos - dstBase
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// readFlush returns a slice of the historical buffer that is ready to be</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// emitted to the user. The data returned by readFlush must be fully consumed</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// before calling any other dictDecoder methods.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>func (dd *dictDecoder) readFlush() []byte {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	toRead := dd.hist[dd.rdPos:dd.wrPos]
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	dd.rdPos = dd.wrPos
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	if dd.wrPos == len(dd.hist) {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		dd.wrPos, dd.rdPos = 0, 0
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		dd.full = true
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	return toRead
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
</pre><p><a href="dict_decoder.go?m=text">View as plain text</a></p>

<div id="footer">
Build version go1.22.2.<br>
Except as <a href="https://developers.google.com/site-policies#restrictions">noted</a>,
the content of this page is licensed under the
Creative Commons Attribution 3.0 License,
and code is licensed under a <a href="http://localhost:8080/LICENSE">BSD license</a>.<br>
<a href="https://golang.org/doc/tos.html">Terms of Service</a> |
<a href="https://www.google.com/intl/en/policies/privacy/">Privacy Policy</a>
</div>

</div><!-- .container -->
</div><!-- #page -->
</body>
</html>
