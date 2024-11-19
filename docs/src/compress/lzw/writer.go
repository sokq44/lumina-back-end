<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/compress/lzw/writer.go - Go Documentation Server</title>

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
<a href="writer.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/compress">compress</a>/<a href="http://localhost:8080/src/compress/lzw">lzw</a>/<span class="text-muted">writer.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/compress/lzw">compress/lzw</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package lzw
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bufio&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// A writer is a buffered, flushable writer.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>type writer interface {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	io.ByteWriter
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	Flush() error
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>}
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>const (
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// A code is a 12 bit value, stored as a uint32 when encoding to avoid</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// type conversions when shifting bits.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	maxCode     = 1&lt;&lt;12 - 1
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	invalidCode = 1&lt;&lt;32 - 1
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// There are 1&lt;&lt;12 possible codes, which is an upper bound on the number of</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// valid hash table entries at any given point in time. tableSize is 4x that.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	tableSize = 4 * 1 &lt;&lt; 12
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	tableMask = tableSize - 1
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// A hash table entry is a uint32. Zero is an invalid entry since the</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// lower 12 bits of a valid entry must be a non-literal code.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	invalidEntry = 0
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// Writer is an LZW compressor. It writes the compressed form of the data</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// to an underlying writer (see [NewWriter]).</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>type Writer struct {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// w is the writer that compressed bytes are written to.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	w writer
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// order, write, bits, nBits and width are the state for</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// converting a code stream into a byte stream.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	order Order
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	write func(*Writer, uint32) error
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	bits  uint32
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	nBits uint
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	width uint
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// litWidth is the width in bits of literal codes.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	litWidth uint
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// hi is the code implied by the next code emission.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// overflow is the code at which hi overflows the code width.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	hi, overflow uint32
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// savedCode is the accumulated code at the end of the most recent Write</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// call. It is equal to invalidCode if there was no such call.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	savedCode uint32
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// err is the first error encountered during writing. Closing the writer</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// will make any future Write calls return errClosed</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	err error
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// table is the hash table from 20-bit keys to 12-bit values. Each table</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// entry contains key&lt;&lt;12|val and collisions resolve by linear probing.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// The keys consist of a 12-bit code prefix and an 8-bit byte suffix.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// The values are a 12-bit code.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	table [tableSize]uint32
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// writeLSB writes the code c for &#34;Least Significant Bits first&#34; data.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>func (w *Writer) writeLSB(c uint32) error {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	w.bits |= c &lt;&lt; w.nBits
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	w.nBits += w.width
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	for w.nBits &gt;= 8 {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		if err := w.w.WriteByte(uint8(w.bits)); err != nil {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			return err
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		w.bits &gt;&gt;= 8
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		w.nBits -= 8
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	return nil
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// writeMSB writes the code c for &#34;Most Significant Bits first&#34; data.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func (w *Writer) writeMSB(c uint32) error {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	w.bits |= c &lt;&lt; (32 - w.width - w.nBits)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	w.nBits += w.width
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	for w.nBits &gt;= 8 {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		if err := w.w.WriteByte(uint8(w.bits &gt;&gt; 24)); err != nil {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			return err
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		w.bits &lt;&lt;= 8
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		w.nBits -= 8
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	return nil
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// errOutOfCodes is an internal error that means that the writer has run out</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// of unused codes and a clear code needs to be sent next.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>var errOutOfCodes = errors.New(&#34;lzw: out of codes&#34;)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// incHi increments e.hi and checks for both overflow and running out of</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// unused codes. In the latter case, incHi sends a clear code, resets the</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// writer state and returns errOutOfCodes.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>func (w *Writer) incHi() error {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	w.hi++
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	if w.hi == w.overflow {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		w.width++
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		w.overflow &lt;&lt;= 1
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	if w.hi == maxCode {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		clear := uint32(1) &lt;&lt; w.litWidth
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		if err := w.write(w, clear); err != nil {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			return err
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		w.width = w.litWidth + 1
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		w.hi = clear + 1
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		w.overflow = clear &lt;&lt; 1
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		for i := range w.table {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			w.table[i] = invalidEntry
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		return errOutOfCodes
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	return nil
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// Write writes a compressed representation of p to w&#39;s underlying writer.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>func (w *Writer) Write(p []byte) (n int, err error) {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if w.err != nil {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		return 0, w.err
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	if len(p) == 0 {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		return 0, nil
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if maxLit := uint8(1&lt;&lt;w.litWidth - 1); maxLit != 0xff {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		for _, x := range p {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			if x &gt; maxLit {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>				w.err = errors.New(&#34;lzw: input byte too large for the litWidth&#34;)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>				return 0, w.err
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	n = len(p)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	code := w.savedCode
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	if code == invalidCode {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		<span class="comment">// This is the first write; send a clear code.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		<span class="comment">// https://www.w3.org/Graphics/GIF/spec-gif89a.txt Appendix F</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		<span class="comment">// &#34;Variable-Length-Code LZW Compression&#34; says that &#34;Encoders should</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		<span class="comment">// output a Clear code as the first code of each image data stream&#34;.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		<span class="comment">// LZW compression isn&#39;t only used by GIF, but it&#39;s cheap to follow</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		<span class="comment">// that directive unconditionally.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		clear := uint32(1) &lt;&lt; w.litWidth
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		if err := w.write(w, clear); err != nil {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			return 0, err
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		<span class="comment">// After the starting clear code, the next code sent (for non-empty</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		<span class="comment">// input) is always a literal code.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		code, p = uint32(p[0]), p[1:]
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>loop:
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	for _, x := range p {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		literal := uint32(x)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		key := code&lt;&lt;8 | literal
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		<span class="comment">// If there is a hash table hit for this key then we continue the loop</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		<span class="comment">// and do not emit a code yet.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		hash := (key&gt;&gt;12 ^ key) &amp; tableMask
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		for h, t := hash, w.table[hash]; t != invalidEntry; {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			if key == t&gt;&gt;12 {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>				code = t &amp; maxCode
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>				continue loop
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>			h = (h + 1) &amp; tableMask
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			t = w.table[h]
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		<span class="comment">// Otherwise, write the current code, and literal becomes the start of</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		<span class="comment">// the next emitted code.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		if w.err = w.write(w, code); w.err != nil {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			return 0, w.err
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		code = literal
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		<span class="comment">// Increment e.hi, the next implied code. If we run out of codes, reset</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		<span class="comment">// the writer state (including clearing the hash table) and continue.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		if err1 := w.incHi(); err1 != nil {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			if err1 == errOutOfCodes {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>				continue
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			w.err = err1
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			return 0, w.err
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		<span class="comment">// Otherwise, insert key -&gt; e.hi into the map that e.table represents.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		for {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			if w.table[hash] == invalidEntry {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				w.table[hash] = (key &lt;&lt; 12) | w.hi
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				break
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			hash = (hash + 1) &amp; tableMask
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	w.savedCode = code
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	return n, nil
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// Close closes the [Writer], flushing any pending output. It does not close</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// w&#39;s underlying writer.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>func (w *Writer) Close() error {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	if w.err != nil {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		if w.err == errClosed {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			return nil
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		return w.err
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// Make any future calls to Write return errClosed.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	w.err = errClosed
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">// Write the savedCode if valid.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	if w.savedCode != invalidCode {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		if err := w.write(w, w.savedCode); err != nil {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			return err
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		if err := w.incHi(); err != nil &amp;&amp; err != errOutOfCodes {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			return err
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	} else {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// Write the starting clear code, as w.Write did not.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		clear := uint32(1) &lt;&lt; w.litWidth
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		if err := w.write(w, clear); err != nil {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			return err
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">// Write the eof code.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	eof := uint32(1)&lt;&lt;w.litWidth + 1
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	if err := w.write(w, eof); err != nil {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		return err
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// Write the final bits.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	if w.nBits &gt; 0 {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		if w.order == MSB {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			w.bits &gt;&gt;= 24
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		if err := w.w.WriteByte(uint8(w.bits)); err != nil {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			return err
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	return w.w.Flush()
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span><span class="comment">// Reset clears the [Writer]&#39;s state and allows it to be reused again</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span><span class="comment">// as a new [Writer].</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>func (w *Writer) Reset(dst io.Writer, order Order, litWidth int) {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	*w = Writer{}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	w.init(dst, order, litWidth)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// NewWriter creates a new [io.WriteCloser].</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">// Writes to the returned [io.WriteCloser] are compressed and written to w.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">// It is the caller&#39;s responsibility to call Close on the WriteCloser when</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span><span class="comment">// finished writing.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span><span class="comment">// The number of bits to use for literal codes, litWidth, must be in the</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span><span class="comment">// range [2,8] and is typically 8. Input bytes must be less than 1&lt;&lt;litWidth.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// It is guaranteed that the underlying type of the returned [io.WriteCloser]</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// is a *[Writer].</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>func NewWriter(w io.Writer, order Order, litWidth int) io.WriteCloser {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	return newWriter(w, order, litWidth)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>func newWriter(dst io.Writer, order Order, litWidth int) *Writer {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	w := new(Writer)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	w.init(dst, order, litWidth)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	return w
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>func (w *Writer) init(dst io.Writer, order Order, litWidth int) {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	switch order {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	case LSB:
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		w.write = (*Writer).writeLSB
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	case MSB:
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		w.write = (*Writer).writeMSB
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	default:
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		w.err = errors.New(&#34;lzw: unknown order&#34;)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		return
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	if litWidth &lt; 2 || 8 &lt; litWidth {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		w.err = fmt.Errorf(&#34;lzw: litWidth %d out of range&#34;, litWidth)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		return
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	bw, ok := dst.(writer)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	if !ok &amp;&amp; dst != nil {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		bw = bufio.NewWriter(dst)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	w.w = bw
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	lw := uint(litWidth)
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	w.order = order
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	w.width = 1 + lw
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	w.litWidth = lw
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	w.hi = 1&lt;&lt;lw + 1
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	w.overflow = 1 &lt;&lt; (lw + 1)
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	w.savedCode = invalidCode
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
</pre><p><a href="writer.go?m=text">View as plain text</a></p>

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
