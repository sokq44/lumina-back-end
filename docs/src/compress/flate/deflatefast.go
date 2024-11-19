<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/compress/flate/deflatefast.go - Go Documentation Server</title>

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
<a href="deflatefast.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/compress">compress</a>/<a href="http://localhost:8080/src/compress/flate">flate</a>/<span class="text-muted">deflatefast.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;math&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// This encoding algorithm, which prioritizes speed over output size, is</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// based on Snappy&#39;s LZ77-style encoder: github.com/golang/snappy</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>const (
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	tableBits  = 14             <span class="comment">// Bits used in the table.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	tableSize  = 1 &lt;&lt; tableBits <span class="comment">// Size of the table.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	tableMask  = tableSize - 1  <span class="comment">// Mask for table indices. Redundant, but can eliminate bounds checks.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	tableShift = 32 - tableBits <span class="comment">// Right-shift to get the tableBits most significant bits of a uint32.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// Reset the buffer offset when reaching this.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// Offsets are stored between blocks as int32 values.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// Since the offset we are checking against is at the beginning</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// of the buffer, we need to subtract the current and input</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// buffer to not risk overflowing the int32.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	bufferReset = math.MaxInt32 - maxStoreBlockSize*2
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>func load32(b []byte, i int32) uint32 {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	b = b[i : i+4 : len(b)] <span class="comment">// Help the compiler eliminate bounds checks on the next line.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	return uint32(b[0]) | uint32(b[1])&lt;&lt;8 | uint32(b[2])&lt;&lt;16 | uint32(b[3])&lt;&lt;24
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>func load64(b []byte, i int32) uint64 {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	b = b[i : i+8 : len(b)] <span class="comment">// Help the compiler eliminate bounds checks on the next line.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	return uint64(b[0]) | uint64(b[1])&lt;&lt;8 | uint64(b[2])&lt;&lt;16 | uint64(b[3])&lt;&lt;24 |
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		uint64(b[4])&lt;&lt;32 | uint64(b[5])&lt;&lt;40 | uint64(b[6])&lt;&lt;48 | uint64(b[7])&lt;&lt;56
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>func hash(u uint32) uint32 {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	return (u * 0x1e35a7bd) &gt;&gt; tableShift
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// These constants are defined by the Snappy implementation so that its</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// assembly implementation can fast-path some 16-bytes-at-a-time copies. They</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// aren&#39;t necessary in the pure Go implementation, as we don&#39;t use those same</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// optimizations, but using the same thresholds doesn&#39;t really hurt.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>const (
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	inputMargin            = 16 - 1
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	minNonLiteralBlockSize = 1 + 1 + inputMargin
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>type tableEntry struct {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	val    uint32 <span class="comment">// Value at destination</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	offset int32
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// deflateFast maintains the table for matches,</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// and the previous byte block for cross block matching.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>type deflateFast struct {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	table [tableSize]tableEntry
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	prev  []byte <span class="comment">// Previous block, zero length if unknown.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	cur   int32  <span class="comment">// Current match offset.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>func newDeflateFast() *deflateFast {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	return &amp;deflateFast{cur: maxStoreBlockSize, prev: make([]byte, 0, maxStoreBlockSize)}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// encode encodes a block given in src and appends tokens</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// to dst and returns the result.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>func (e *deflateFast) encode(dst []token, src []byte) []token {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// Ensure that e.cur doesn&#39;t wrap.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if e.cur &gt;= bufferReset {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		e.shiftOffsets()
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// This check isn&#39;t in the Snappy implementation, but there, the caller</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// instead of the callee handles this case.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	if len(src) &lt; minNonLiteralBlockSize {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		e.cur += maxStoreBlockSize
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		e.prev = e.prev[:0]
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		return emitLiteral(dst, src)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// sLimit is when to stop looking for offset/length copies. The inputMargin</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// lets us use a fast path for emitLiteral in the main loop, while we are</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// looking for copies.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	sLimit := int32(len(src) - inputMargin)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// nextEmit is where in src the next emitLiteral should start from.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	nextEmit := int32(0)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	s := int32(0)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	cv := load32(src, s)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	nextHash := hash(cv)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	for {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		<span class="comment">// Copied from the C++ snappy implementation:</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		<span class="comment">// Heuristic match skipping: If 32 bytes are scanned with no matches</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		<span class="comment">// found, start looking only at every other byte. If 32 more bytes are</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		<span class="comment">// scanned (or skipped), look at every third byte, etc.. When a match</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		<span class="comment">// is found, immediately go back to looking at every byte. This is a</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		<span class="comment">// small loss (~5% performance, ~0.1% density) for compressible data</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		<span class="comment">// due to more bookkeeping, but for non-compressible data (such as</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		<span class="comment">// JPEG) it&#39;s a huge win since the compressor quickly &#34;realizes&#34; the</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		<span class="comment">// data is incompressible and doesn&#39;t bother looking for matches</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		<span class="comment">// everywhere.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		<span class="comment">// The &#34;skip&#34; variable keeps track of how many bytes there are since</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		<span class="comment">// the last match; dividing it by 32 (ie. right-shifting by five) gives</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		<span class="comment">// the number of bytes to move ahead for each iteration.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		skip := int32(32)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		nextS := s
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		var candidate tableEntry
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		for {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			s = nextS
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			bytesBetweenHashLookups := skip &gt;&gt; 5
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			nextS = s + bytesBetweenHashLookups
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			skip += bytesBetweenHashLookups
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			if nextS &gt; sLimit {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>				goto emitRemainder
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			candidate = e.table[nextHash&amp;tableMask]
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			now := load32(src, nextS)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			e.table[nextHash&amp;tableMask] = tableEntry{offset: s + e.cur, val: cv}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			nextHash = hash(now)
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			offset := s - (candidate.offset - e.cur)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			if offset &gt; maxMatchOffset || cv != candidate.val {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>				<span class="comment">// Out of range or not matched.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>				cv = now
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>				continue
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			break
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// A 4-byte match has been found. We&#39;ll later see if more than 4 bytes</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		<span class="comment">// match. But, prior to the match, src[nextEmit:s] are unmatched. Emit</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		<span class="comment">// them as literal bytes.</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		dst = emitLiteral(dst, src[nextEmit:s])
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		<span class="comment">// Call emitCopy, and then see if another emitCopy could be our next</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		<span class="comment">// move. Repeat until we find no match for the input immediately after</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		<span class="comment">// what was consumed by the last emitCopy call.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		<span class="comment">// If we exit this loop normally then we need to call emitLiteral next,</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		<span class="comment">// though we don&#39;t yet know how big the literal will be. We handle that</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		<span class="comment">// by proceeding to the next iteration of the main loop. We also can</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		<span class="comment">// exit this loop via goto if we get close to exhausting the input.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		for {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			<span class="comment">// Invariant: we have a 4-byte match at s, and no need to emit any</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			<span class="comment">// literal bytes prior to s.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			<span class="comment">// Extend the 4-byte match as long as possible.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			s += 4
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			t := candidate.offset - e.cur + 4
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			l := e.matchLen(s, t, src)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			<span class="comment">// matchToken is flate&#39;s equivalent of Snappy&#39;s emitCopy. (length,offset)</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			dst = append(dst, matchToken(uint32(l+4-baseMatchLength), uint32(s-t-baseMatchOffset)))
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			s += l
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			nextEmit = s
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			if s &gt;= sLimit {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>				goto emitRemainder
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>			<span class="comment">// We could immediately start working at s now, but to improve</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			<span class="comment">// compression we first update the hash table at s-1 and at s. If</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			<span class="comment">// another emitCopy is not our next move, also calculate nextHash</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			<span class="comment">// at s+1. At least on GOARCH=amd64, these three hash calculations</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			<span class="comment">// are faster as one load64 call (with some shifts) instead of</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			<span class="comment">// three load32 calls.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			x := load64(src, s-1)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			prevHash := hash(uint32(x))
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			e.table[prevHash&amp;tableMask] = tableEntry{offset: e.cur + s - 1, val: uint32(x)}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			x &gt;&gt;= 8
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			currHash := hash(uint32(x))
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			candidate = e.table[currHash&amp;tableMask]
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			e.table[currHash&amp;tableMask] = tableEntry{offset: e.cur + s, val: uint32(x)}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			offset := s - (candidate.offset - e.cur)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			if offset &gt; maxMatchOffset || uint32(x) != candidate.val {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>				cv = uint32(x &gt;&gt; 8)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>				nextHash = hash(cv)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>				s++
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>				break
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>emitRemainder:
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	if int(nextEmit) &lt; len(src) {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		dst = emitLiteral(dst, src[nextEmit:])
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	e.cur += int32(len(src))
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	e.prev = e.prev[:len(src)]
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	copy(e.prev, src)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	return dst
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>func emitLiteral(dst []token, lit []byte) []token {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	for _, v := range lit {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		dst = append(dst, literalToken(uint32(v)))
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	return dst
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// matchLen returns the match length between src[s:] and src[t:].</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// t can be negative to indicate the match is starting in e.prev.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// We assume that src[s-4:s] and src[t-4:t] already match.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>func (e *deflateFast) matchLen(s, t int32, src []byte) int32 {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	s1 := int(s) + maxMatchLength - 4
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	if s1 &gt; len(src) {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		s1 = len(src)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	<span class="comment">// If we are inside the current block</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	if t &gt;= 0 {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		b := src[t:]
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		a := src[s:s1]
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		b = b[:len(a)]
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		<span class="comment">// Extend the match to be as long as possible.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		for i := range a {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			if a[i] != b[i] {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>				return int32(i)
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		return int32(len(a))
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">// We found a match in the previous block.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	tp := int32(len(e.prev)) + t
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	if tp &lt; 0 {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		return 0
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// Extend the match to be as long as possible.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	a := src[s:s1]
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	b := e.prev[tp:]
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	if len(b) &gt; len(a) {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		b = b[:len(a)]
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	a = a[:len(b)]
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	for i := range b {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		if a[i] != b[i] {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			return int32(i)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// If we reached our limit, we matched everything we are</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// allowed to in the previous block and we return.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	n := int32(len(b))
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	if int(s+n) == s1 {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		return n
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">// Continue looking for more matches in the current block.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	a = src[s+n : s1]
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	b = src[:len(a)]
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	for i := range a {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		if a[i] != b[i] {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			return int32(i) + n
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	return int32(len(a)) + n
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">// Reset resets the encoding history.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">// This ensures that no matches are made to the previous block.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>func (e *deflateFast) reset() {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	e.prev = e.prev[:0]
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	<span class="comment">// Bump the offset, so all matches will fail distance check.</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	<span class="comment">// Nothing should be &gt;= e.cur in the table.</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	e.cur += maxMatchOffset
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">// Protect against e.cur wraparound.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	if e.cur &gt;= bufferReset {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		e.shiftOffsets()
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">// shiftOffsets will shift down all match offset.</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">// This is only called in rare situations to prevent integer overflow.</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">// See https://golang.org/issue/18636 and https://github.com/golang/go/issues/34121.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>func (e *deflateFast) shiftOffsets() {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	if len(e.prev) == 0 {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		<span class="comment">// We have no history; just clear the table.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		for i := range e.table[:] {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			e.table[i] = tableEntry{}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		e.cur = maxMatchOffset + 1
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		return
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">// Shift down everything in the table that isn&#39;t already too far away.</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	for i := range e.table[:] {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		v := e.table[i].offset - e.cur + maxMatchOffset + 1
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		if v &lt; 0 {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			<span class="comment">// We want to reset e.cur to maxMatchOffset + 1, so we need to shift</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			<span class="comment">// all table entries down by (e.cur - (maxMatchOffset + 1)).</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			<span class="comment">// Because we ignore matches &gt; maxMatchOffset, we can cap</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			<span class="comment">// any negative offsets at 0.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			v = 0
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		e.table[i].offset = v
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	e.cur = maxMatchOffset + 1
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
</pre><p><a href="deflatefast.go?m=text">View as plain text</a></p>

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
