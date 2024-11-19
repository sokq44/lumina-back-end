<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/hash/crc32/crc32_amd64.go - Go Documentation Server</title>

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
<a href="crc32_amd64.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/hash">hash</a>/<a href="http://localhost:8080/src/hash/crc32">crc32</a>/<span class="text-muted">crc32_amd64.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/hash/crc32">hash/crc32</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// AMD64-specific hardware-assisted CRC32 algorithms. See crc32.go for a</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// description of the interface that each architecture-specific file</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// implements.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package crc32
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/cpu&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// This file contains the code to call the SSE 4.2 version of the Castagnoli</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// and IEEE CRC.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// castagnoliSSE42 is defined in crc32_amd64.s and uses the SSE 4.2 CRC32</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// instruction.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>func castagnoliSSE42(crc uint32, p []byte) uint32
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// castagnoliSSE42Triple is defined in crc32_amd64.s and uses the SSE 4.2 CRC32</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// instruction.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>func castagnoliSSE42Triple(
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	crcA, crcB, crcC uint32,
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	a, b, c []byte,
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	rounds uint32,
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>) (retA uint32, retB uint32, retC uint32)
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// ieeeCLMUL is defined in crc_amd64.s and uses the PCLMULQDQ</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// instruction as well as SSE 4.1.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>func ieeeCLMUL(crc uint32, p []byte) uint32
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>const castagnoliK1 = 168
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>const castagnoliK2 = 1344
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>type sse42Table [4]Table
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>var castagnoliSSE42TableK1 *sse42Table
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>var castagnoliSSE42TableK2 *sse42Table
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>func archAvailableCastagnoli() bool {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	return cpu.X86.HasSSE42
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>func archInitCastagnoli() {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	if !cpu.X86.HasSSE42 {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		panic(&#34;arch-specific Castagnoli not available&#34;)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	castagnoliSSE42TableK1 = new(sse42Table)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	castagnoliSSE42TableK2 = new(sse42Table)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// See description in updateCastagnoli.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">//    t[0][i] = CRC(i000, O)</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">//    t[1][i] = CRC(0i00, O)</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">//    t[2][i] = CRC(00i0, O)</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">//    t[3][i] = CRC(000i, O)</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// where O is a sequence of K zeros.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	var tmp [castagnoliK2]byte
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	for b := 0; b &lt; 4; b++ {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		for i := 0; i &lt; 256; i++ {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			val := uint32(i) &lt;&lt; uint32(b*8)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			castagnoliSSE42TableK1[b][i] = castagnoliSSE42(val, tmp[:castagnoliK1])
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			castagnoliSSE42TableK2[b][i] = castagnoliSSE42(val, tmp[:])
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// castagnoliShift computes the CRC32-C of K1 or K2 zeroes (depending on the</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// table given) with the given initial crc value. This corresponds to</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// CRC(crc, O) in the description in updateCastagnoli.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>func castagnoliShift(table *sse42Table, crc uint32) uint32 {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	return table[3][crc&gt;&gt;24] ^
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		table[2][(crc&gt;&gt;16)&amp;0xFF] ^
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		table[1][(crc&gt;&gt;8)&amp;0xFF] ^
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		table[0][crc&amp;0xFF]
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>func archUpdateCastagnoli(crc uint32, p []byte) uint32 {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	if !cpu.X86.HasSSE42 {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		panic(&#34;not available&#34;)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// This method is inspired from the algorithm in Intel&#39;s white paper:</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">//    &#34;Fast CRC Computation for iSCSI Polynomial Using CRC32 Instruction&#34;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// The same strategy of splitting the buffer in three is used but the</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// combining calculation is different; the complete derivation is explained</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// below.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// -- The basic idea --</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// The CRC32 instruction (available in SSE4.2) can process 8 bytes at a</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// time. In recent Intel architectures the instruction takes 3 cycles;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// however the processor can pipeline up to three instructions if they</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// don&#39;t depend on each other.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// Roughly this means that we can process three buffers in about the same</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// time we can process one buffer.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	<span class="comment">// The idea is then to split the buffer in three, CRC the three pieces</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// separately and then combine the results.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// Combining the results requires precomputed tables, so we must choose a</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// fixed buffer length to optimize. The longer the length, the faster; but</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// only buffers longer than this length will use the optimization. We choose</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// two cutoffs and compute tables for both:</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">//  - one around 512: 168*3=504</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">//  - one around 4KB: 1344*3=4032</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// -- The nitty gritty --</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// Let CRC(I, X) be the non-inverted CRC32-C of the sequence X (with</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// initial non-inverted CRC I). This function has the following properties:</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">//   (a) CRC(I, AB) = CRC(CRC(I, A), B)</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">//   (b) CRC(I, A xor B) = CRC(I, A) xor CRC(0, B)</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// Say we want to compute CRC(I, ABC) where A, B, C are three sequences of</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// K bytes each, where K is a fixed constant. Let O be the sequence of K zero</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// bytes.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// CRC(I, ABC) = CRC(I, ABO xor C)</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">//             = CRC(I, ABO) xor CRC(0, C)</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">//             = CRC(CRC(I, AB), O) xor CRC(0, C)</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">//             = CRC(CRC(I, AO xor B), O) xor CRC(0, C)</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">//             = CRC(CRC(I, AO) xor CRC(0, B), O) xor CRC(0, C)</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">//             = CRC(CRC(CRC(I, A), O) xor CRC(0, B), O) xor CRC(0, C)</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// The castagnoliSSE42Triple function can compute CRC(I, A), CRC(0, B),</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// and CRC(0, C) efficiently.  We just need to find a way to quickly compute</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// CRC(uvwx, O) given a 4-byte initial value uvwx. We can precompute these</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// values; since we can&#39;t have a 32-bit table, we break it up into four</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// 8-bit tables:</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">//    CRC(uvwx, O) = CRC(u000, O) xor</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">//                   CRC(0v00, O) xor</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">//                   CRC(00w0, O) xor</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">//                   CRC(000x, O)</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// We can compute tables corresponding to the four terms for all 8-bit</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// values.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	crc = ^crc
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// If a buffer is long enough to use the optimization, process the first few</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// bytes to align the buffer to an 8 byte boundary (if necessary).</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	if len(p) &gt;= castagnoliK1*3 {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		delta := int(uintptr(unsafe.Pointer(&amp;p[0])) &amp; 7)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		if delta != 0 {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			delta = 8 - delta
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			crc = castagnoliSSE42(crc, p[:delta])
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			p = p[delta:]
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// Process 3*K2 at a time.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	for len(p) &gt;= castagnoliK2*3 {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		<span class="comment">// Compute CRC(I, A), CRC(0, B), and CRC(0, C).</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		crcA, crcB, crcC := castagnoliSSE42Triple(
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			crc, 0, 0,
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			p, p[castagnoliK2:], p[castagnoliK2*2:],
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>			castagnoliK2/24)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		<span class="comment">// CRC(I, AB) = CRC(CRC(I, A), O) xor CRC(0, B)</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		crcAB := castagnoliShift(castagnoliSSE42TableK2, crcA) ^ crcB
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		<span class="comment">// CRC(I, ABC) = CRC(CRC(I, AB), O) xor CRC(0, C)</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		crc = castagnoliShift(castagnoliSSE42TableK2, crcAB) ^ crcC
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		p = p[castagnoliK2*3:]
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// Process 3*K1 at a time.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	for len(p) &gt;= castagnoliK1*3 {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		<span class="comment">// Compute CRC(I, A), CRC(0, B), and CRC(0, C).</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		crcA, crcB, crcC := castagnoliSSE42Triple(
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			crc, 0, 0,
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			p, p[castagnoliK1:], p[castagnoliK1*2:],
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			castagnoliK1/24)
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		<span class="comment">// CRC(I, AB) = CRC(CRC(I, A), O) xor CRC(0, B)</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		crcAB := castagnoliShift(castagnoliSSE42TableK1, crcA) ^ crcB
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		<span class="comment">// CRC(I, ABC) = CRC(CRC(I, AB), O) xor CRC(0, C)</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		crc = castagnoliShift(castagnoliSSE42TableK1, crcAB) ^ crcC
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		p = p[castagnoliK1*3:]
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// Use the simple implementation for what&#39;s left.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	crc = castagnoliSSE42(crc, p)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	return ^crc
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>func archAvailableIEEE() bool {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	return cpu.X86.HasPCLMULQDQ &amp;&amp; cpu.X86.HasSSE41
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>var archIeeeTable8 *slicing8Table
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>func archInitIEEE() {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	if !cpu.X86.HasPCLMULQDQ || !cpu.X86.HasSSE41 {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		panic(&#34;not available&#34;)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// We still use slicing-by-8 for small buffers.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	archIeeeTable8 = slicingMakeTable(IEEE)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>func archUpdateIEEE(crc uint32, p []byte) uint32 {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if !cpu.X86.HasPCLMULQDQ || !cpu.X86.HasSSE41 {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		panic(&#34;not available&#34;)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	if len(p) &gt;= 64 {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		left := len(p) &amp; 15
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		do := len(p) - left
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		crc = ^ieeeCLMUL(^crc, p[:do])
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		p = p[do:]
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	if len(p) == 0 {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		return crc
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	return slicingUpdate(crc, archIeeeTable8, p)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
</pre><p><a href="crc32_amd64.go?m=text">View as plain text</a></p>

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
