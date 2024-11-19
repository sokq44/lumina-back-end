<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/strconv/eisel_lemire.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../index.html">GoDoc</a></div>
<a href="eisel_lemire.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/strconv">strconv</a>/<span class="text-muted">eisel_lemire.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/strconv">strconv</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2020 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package strconv
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This file implements the Eisel-Lemire ParseFloat algorithm, published in</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// 2020 and discussed extensively at</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// https://nigeltao.github.io/blog/2020/eisel-lemire.html</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// The original C++ implementation is at</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// https://github.com/lemire/fast_double_parser/blob/644bef4306059d3be01a04e77d3cc84b379c596f/include/fast_double_parser.h#L840</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// This Go re-implementation closely follows the C re-implementation at</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// https://github.com/google/wuffs/blob/ba3818cb6b473a2ed0b38ecfc07dbbd3a97e8ae7/internal/cgen/base/floatconv-submodule-code.c#L990</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Additional testing (on over several million test strings) is done by</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// https://github.com/nigeltao/parse-number-fxx-test-data/blob/5280dcfccf6d0b02a65ae282dad0b6d9de50e039/script/test-go-strconv.go</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>import (
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;math/bits&#34;
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>func eiselLemire64(man uint64, exp10 int, neg bool) (f float64, ok bool) {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// The terse comments in this function body refer to sections of the</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// https://nigeltao.github.io/blog/2020/eisel-lemire.html blog post.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// Exp10 Range.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	if man == 0 {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		if neg {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>			f = math.Float64frombits(0x8000000000000000) <span class="comment">// Negative zero.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		return f, true
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	if exp10 &lt; detailedPowersOfTenMinExp10 || detailedPowersOfTenMaxExp10 &lt; exp10 {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		return 0, false
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// Normalization.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	clz := bits.LeadingZeros64(man)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	man &lt;&lt;= uint(clz)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	const float64ExponentBias = 1023
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	retExp2 := uint64(217706*exp10&gt;&gt;16+64+float64ExponentBias) - uint64(clz)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// Multiplication.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	xHi, xLo := bits.Mul64(man, detailedPowersOfTen[exp10-detailedPowersOfTenMinExp10][1])
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// Wider Approximation.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	if xHi&amp;0x1FF == 0x1FF &amp;&amp; xLo+man &lt; man {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		yHi, yLo := bits.Mul64(man, detailedPowersOfTen[exp10-detailedPowersOfTenMinExp10][0])
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		mergedHi, mergedLo := xHi, xLo+yHi
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		if mergedLo &lt; xLo {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			mergedHi++
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		if mergedHi&amp;0x1FF == 0x1FF &amp;&amp; mergedLo+1 == 0 &amp;&amp; yLo+man &lt; man {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			return 0, false
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		xHi, xLo = mergedHi, mergedLo
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// Shifting to 54 Bits.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	msb := xHi &gt;&gt; 63
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	retMantissa := xHi &gt;&gt; (msb + 9)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	retExp2 -= 1 ^ msb
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// Half-way Ambiguity.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	if xLo == 0 &amp;&amp; xHi&amp;0x1FF == 0 &amp;&amp; retMantissa&amp;3 == 1 {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		return 0, false
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// From 54 to 53 Bits.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	retMantissa += retMantissa &amp; 1
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	retMantissa &gt;&gt;= 1
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	if retMantissa&gt;&gt;53 &gt; 0 {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		retMantissa &gt;&gt;= 1
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		retExp2 += 1
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// retExp2 is a uint64. Zero or underflow means that we&#39;re in subnormal</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// float64 space. 0x7FF or above means that we&#39;re in Inf/NaN float64 space.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// The if block is equivalent to (but has fewer branches than):</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">//   if retExp2 &lt;= 0 || retExp2 &gt;= 0x7FF { etc }</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	if retExp2-1 &gt;= 0x7FF-1 {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		return 0, false
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	retBits := retExp2&lt;&lt;52 | retMantissa&amp;0x000FFFFFFFFFFFFF
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	if neg {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		retBits |= 0x8000000000000000
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	return math.Float64frombits(retBits), true
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>func eiselLemire32(man uint64, exp10 int, neg bool) (f float32, ok bool) {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// The terse comments in this function body refer to sections of the</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// https://nigeltao.github.io/blog/2020/eisel-lemire.html blog post.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// That blog post discusses the float64 flavor (11 exponent bits with a</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// -1023 bias, 52 mantissa bits) of the algorithm, but the same approach</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// applies to the float32 flavor (8 exponent bits with a -127 bias, 23</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// mantissa bits). The computation here happens with 64-bit values (e.g.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// man, xHi, retMantissa) before finally converting to a 32-bit float.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// Exp10 Range.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	if man == 0 {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		if neg {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			f = math.Float32frombits(0x80000000) <span class="comment">// Negative zero.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		return f, true
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	if exp10 &lt; detailedPowersOfTenMinExp10 || detailedPowersOfTenMaxExp10 &lt; exp10 {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		return 0, false
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// Normalization.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	clz := bits.LeadingZeros64(man)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	man &lt;&lt;= uint(clz)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	const float32ExponentBias = 127
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	retExp2 := uint64(217706*exp10&gt;&gt;16+64+float32ExponentBias) - uint64(clz)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// Multiplication.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	xHi, xLo := bits.Mul64(man, detailedPowersOfTen[exp10-detailedPowersOfTenMinExp10][1])
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// Wider Approximation.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	if xHi&amp;0x3FFFFFFFFF == 0x3FFFFFFFFF &amp;&amp; xLo+man &lt; man {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		yHi, yLo := bits.Mul64(man, detailedPowersOfTen[exp10-detailedPowersOfTenMinExp10][0])
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		mergedHi, mergedLo := xHi, xLo+yHi
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		if mergedLo &lt; xLo {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			mergedHi++
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		if mergedHi&amp;0x3FFFFFFFFF == 0x3FFFFFFFFF &amp;&amp; mergedLo+1 == 0 &amp;&amp; yLo+man &lt; man {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			return 0, false
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		xHi, xLo = mergedHi, mergedLo
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// Shifting to 54 Bits (and for float32, it&#39;s shifting to 25 bits).</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	msb := xHi &gt;&gt; 63
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	retMantissa := xHi &gt;&gt; (msb + 38)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	retExp2 -= 1 ^ msb
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// Half-way Ambiguity.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	if xLo == 0 &amp;&amp; xHi&amp;0x3FFFFFFFFF == 0 &amp;&amp; retMantissa&amp;3 == 1 {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		return 0, false
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">// From 54 to 53 Bits (and for float32, it&#39;s from 25 to 24 bits).</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	retMantissa += retMantissa &amp; 1
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	retMantissa &gt;&gt;= 1
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	if retMantissa&gt;&gt;24 &gt; 0 {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		retMantissa &gt;&gt;= 1
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		retExp2 += 1
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// retExp2 is a uint64. Zero or underflow means that we&#39;re in subnormal</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// float32 space. 0xFF or above means that we&#39;re in Inf/NaN float32 space.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// The if block is equivalent to (but has fewer branches than):</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">//   if retExp2 &lt;= 0 || retExp2 &gt;= 0xFF { etc }</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if retExp2-1 &gt;= 0xFF-1 {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		return 0, false
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	retBits := retExp2&lt;&lt;23 | retMantissa&amp;0x007FFFFF
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if neg {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		retBits |= 0x80000000
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	return math.Float32frombits(uint32(retBits)), true
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// detailedPowersOfTen{Min,Max}Exp10 is the power of 10 represented by the</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">// first and last rows of detailedPowersOfTen. Both bounds are inclusive.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>const (
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	detailedPowersOfTenMinExp10 = -348
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	detailedPowersOfTenMaxExp10 = +347
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// detailedPowersOfTen contains 128-bit mantissa approximations (rounded down)</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">// to the powers of 10. For example:</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">//   - 1e43 ≈ (0xE596B7B0_C643C719                   * (2 ** 79))</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span><span class="comment">//   - 1e43 = (0xE596B7B0_C643C719_6D9CCD05_D0000000 * (2 ** 15))</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">// The mantissas are explicitly listed. The exponents are implied by a linear</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">// expression with slope 217706.0/65536.0 ≈ log(10)/log(2).</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span><span class="comment">// The table was generated by</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// https://github.com/google/wuffs/blob/ba3818cb6b473a2ed0b38ecfc07dbbd3a97e8ae7/script/print-mpb-powers-of-10.go</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>var detailedPowersOfTen = [...][2]uint64{
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	{0x1732C869CD60E453, 0xFA8FD5A0081C0288}, <span class="comment">// 1e-348</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	{0x0E7FBD42205C8EB4, 0x9C99E58405118195}, <span class="comment">// 1e-347</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	{0x521FAC92A873B261, 0xC3C05EE50655E1FA}, <span class="comment">// 1e-346</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	{0xE6A797B752909EF9, 0xF4B0769E47EB5A78}, <span class="comment">// 1e-345</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	{0x9028BED2939A635C, 0x98EE4A22ECF3188B}, <span class="comment">// 1e-344</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	{0x7432EE873880FC33, 0xBF29DCABA82FDEAE}, <span class="comment">// 1e-343</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	{0x113FAA2906A13B3F, 0xEEF453D6923BD65A}, <span class="comment">// 1e-342</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	{0x4AC7CA59A424C507, 0x9558B4661B6565F8}, <span class="comment">// 1e-341</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	{0x5D79BCF00D2DF649, 0xBAAEE17FA23EBF76}, <span class="comment">// 1e-340</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	{0xF4D82C2C107973DC, 0xE95A99DF8ACE6F53}, <span class="comment">// 1e-339</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	{0x79071B9B8A4BE869, 0x91D8A02BB6C10594}, <span class="comment">// 1e-338</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	{0x9748E2826CDEE284, 0xB64EC836A47146F9}, <span class="comment">// 1e-337</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	{0xFD1B1B2308169B25, 0xE3E27A444D8D98B7}, <span class="comment">// 1e-336</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	{0xFE30F0F5E50E20F7, 0x8E6D8C6AB0787F72}, <span class="comment">// 1e-335</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	{0xBDBD2D335E51A935, 0xB208EF855C969F4F}, <span class="comment">// 1e-334</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	{0xAD2C788035E61382, 0xDE8B2B66B3BC4723}, <span class="comment">// 1e-333</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	{0x4C3BCB5021AFCC31, 0x8B16FB203055AC76}, <span class="comment">// 1e-332</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	{0xDF4ABE242A1BBF3D, 0xADDCB9E83C6B1793}, <span class="comment">// 1e-331</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	{0xD71D6DAD34A2AF0D, 0xD953E8624B85DD78}, <span class="comment">// 1e-330</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	{0x8672648C40E5AD68, 0x87D4713D6F33AA6B}, <span class="comment">// 1e-329</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	{0x680EFDAF511F18C2, 0xA9C98D8CCB009506}, <span class="comment">// 1e-328</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	{0x0212BD1B2566DEF2, 0xD43BF0EFFDC0BA48}, <span class="comment">// 1e-327</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	{0x014BB630F7604B57, 0x84A57695FE98746D}, <span class="comment">// 1e-326</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	{0x419EA3BD35385E2D, 0xA5CED43B7E3E9188}, <span class="comment">// 1e-325</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	{0x52064CAC828675B9, 0xCF42894A5DCE35EA}, <span class="comment">// 1e-324</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	{0x7343EFEBD1940993, 0x818995CE7AA0E1B2}, <span class="comment">// 1e-323</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	{0x1014EBE6C5F90BF8, 0xA1EBFB4219491A1F}, <span class="comment">// 1e-322</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	{0xD41A26E077774EF6, 0xCA66FA129F9B60A6}, <span class="comment">// 1e-321</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	{0x8920B098955522B4, 0xFD00B897478238D0}, <span class="comment">// 1e-320</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	{0x55B46E5F5D5535B0, 0x9E20735E8CB16382}, <span class="comment">// 1e-319</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	{0xEB2189F734AA831D, 0xC5A890362FDDBC62}, <span class="comment">// 1e-318</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	{0xA5E9EC7501D523E4, 0xF712B443BBD52B7B}, <span class="comment">// 1e-317</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	{0x47B233C92125366E, 0x9A6BB0AA55653B2D}, <span class="comment">// 1e-316</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	{0x999EC0BB696E840A, 0xC1069CD4EABE89F8}, <span class="comment">// 1e-315</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	{0xC00670EA43CA250D, 0xF148440A256E2C76}, <span class="comment">// 1e-314</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	{0x380406926A5E5728, 0x96CD2A865764DBCA}, <span class="comment">// 1e-313</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	{0xC605083704F5ECF2, 0xBC807527ED3E12BC}, <span class="comment">// 1e-312</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	{0xF7864A44C633682E, 0xEBA09271E88D976B}, <span class="comment">// 1e-311</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	{0x7AB3EE6AFBE0211D, 0x93445B8731587EA3}, <span class="comment">// 1e-310</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	{0x5960EA05BAD82964, 0xB8157268FDAE9E4C}, <span class="comment">// 1e-309</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	{0x6FB92487298E33BD, 0xE61ACF033D1A45DF}, <span class="comment">// 1e-308</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	{0xA5D3B6D479F8E056, 0x8FD0C16206306BAB}, <span class="comment">// 1e-307</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	{0x8F48A4899877186C, 0xB3C4F1BA87BC8696}, <span class="comment">// 1e-306</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	{0x331ACDABFE94DE87, 0xE0B62E2929ABA83C}, <span class="comment">// 1e-305</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	{0x9FF0C08B7F1D0B14, 0x8C71DCD9BA0B4925}, <span class="comment">// 1e-304</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	{0x07ECF0AE5EE44DD9, 0xAF8E5410288E1B6F}, <span class="comment">// 1e-303</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	{0xC9E82CD9F69D6150, 0xDB71E91432B1A24A}, <span class="comment">// 1e-302</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	{0xBE311C083A225CD2, 0x892731AC9FAF056E}, <span class="comment">// 1e-301</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	{0x6DBD630A48AAF406, 0xAB70FE17C79AC6CA}, <span class="comment">// 1e-300</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	{0x092CBBCCDAD5B108, 0xD64D3D9DB981787D}, <span class="comment">// 1e-299</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	{0x25BBF56008C58EA5, 0x85F0468293F0EB4E}, <span class="comment">// 1e-298</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	{0xAF2AF2B80AF6F24E, 0xA76C582338ED2621}, <span class="comment">// 1e-297</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	{0x1AF5AF660DB4AEE1, 0xD1476E2C07286FAA}, <span class="comment">// 1e-296</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	{0x50D98D9FC890ED4D, 0x82CCA4DB847945CA}, <span class="comment">// 1e-295</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	{0xE50FF107BAB528A0, 0xA37FCE126597973C}, <span class="comment">// 1e-294</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	{0x1E53ED49A96272C8, 0xCC5FC196FEFD7D0C}, <span class="comment">// 1e-293</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	{0x25E8E89C13BB0F7A, 0xFF77B1FCBEBCDC4F}, <span class="comment">// 1e-292</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	{0x77B191618C54E9AC, 0x9FAACF3DF73609B1}, <span class="comment">// 1e-291</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	{0xD59DF5B9EF6A2417, 0xC795830D75038C1D}, <span class="comment">// 1e-290</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	{0x4B0573286B44AD1D, 0xF97AE3D0D2446F25}, <span class="comment">// 1e-289</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	{0x4EE367F9430AEC32, 0x9BECCE62836AC577}, <span class="comment">// 1e-288</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	{0x229C41F793CDA73F, 0xC2E801FB244576D5}, <span class="comment">// 1e-287</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	{0x6B43527578C1110F, 0xF3A20279ED56D48A}, <span class="comment">// 1e-286</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	{0x830A13896B78AAA9, 0x9845418C345644D6}, <span class="comment">// 1e-285</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	{0x23CC986BC656D553, 0xBE5691EF416BD60C}, <span class="comment">// 1e-284</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	{0x2CBFBE86B7EC8AA8, 0xEDEC366B11C6CB8F}, <span class="comment">// 1e-283</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	{0x7BF7D71432F3D6A9, 0x94B3A202EB1C3F39}, <span class="comment">// 1e-282</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	{0xDAF5CCD93FB0CC53, 0xB9E08A83A5E34F07}, <span class="comment">// 1e-281</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	{0xD1B3400F8F9CFF68, 0xE858AD248F5C22C9}, <span class="comment">// 1e-280</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	{0x23100809B9C21FA1, 0x91376C36D99995BE}, <span class="comment">// 1e-279</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	{0xABD40A0C2832A78A, 0xB58547448FFFFB2D}, <span class="comment">// 1e-278</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	{0x16C90C8F323F516C, 0xE2E69915B3FFF9F9}, <span class="comment">// 1e-277</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	{0xAE3DA7D97F6792E3, 0x8DD01FAD907FFC3B}, <span class="comment">// 1e-276</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	{0x99CD11CFDF41779C, 0xB1442798F49FFB4A}, <span class="comment">// 1e-275</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	{0x40405643D711D583, 0xDD95317F31C7FA1D}, <span class="comment">// 1e-274</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	{0x482835EA666B2572, 0x8A7D3EEF7F1CFC52}, <span class="comment">// 1e-273</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	{0xDA3243650005EECF, 0xAD1C8EAB5EE43B66}, <span class="comment">// 1e-272</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	{0x90BED43E40076A82, 0xD863B256369D4A40}, <span class="comment">// 1e-271</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	{0x5A7744A6E804A291, 0x873E4F75E2224E68}, <span class="comment">// 1e-270</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	{0x711515D0A205CB36, 0xA90DE3535AAAE202}, <span class="comment">// 1e-269</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	{0x0D5A5B44CA873E03, 0xD3515C2831559A83}, <span class="comment">// 1e-268</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	{0xE858790AFE9486C2, 0x8412D9991ED58091}, <span class="comment">// 1e-267</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	{0x626E974DBE39A872, 0xA5178FFF668AE0B6}, <span class="comment">// 1e-266</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	{0xFB0A3D212DC8128F, 0xCE5D73FF402D98E3}, <span class="comment">// 1e-265</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	{0x7CE66634BC9D0B99, 0x80FA687F881C7F8E}, <span class="comment">// 1e-264</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	{0x1C1FFFC1EBC44E80, 0xA139029F6A239F72}, <span class="comment">// 1e-263</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	{0xA327FFB266B56220, 0xC987434744AC874E}, <span class="comment">// 1e-262</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	{0x4BF1FF9F0062BAA8, 0xFBE9141915D7A922}, <span class="comment">// 1e-261</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	{0x6F773FC3603DB4A9, 0x9D71AC8FADA6C9B5}, <span class="comment">// 1e-260</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	{0xCB550FB4384D21D3, 0xC4CE17B399107C22}, <span class="comment">// 1e-259</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	{0x7E2A53A146606A48, 0xF6019DA07F549B2B}, <span class="comment">// 1e-258</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	{0x2EDA7444CBFC426D, 0x99C102844F94E0FB}, <span class="comment">// 1e-257</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	{0xFA911155FEFB5308, 0xC0314325637A1939}, <span class="comment">// 1e-256</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	{0x793555AB7EBA27CA, 0xF03D93EEBC589F88}, <span class="comment">// 1e-255</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	{0x4BC1558B2F3458DE, 0x96267C7535B763B5}, <span class="comment">// 1e-254</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	{0x9EB1AAEDFB016F16, 0xBBB01B9283253CA2}, <span class="comment">// 1e-253</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	{0x465E15A979C1CADC, 0xEA9C227723EE8BCB}, <span class="comment">// 1e-252</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	{0x0BFACD89EC191EC9, 0x92A1958A7675175F}, <span class="comment">// 1e-251</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	{0xCEF980EC671F667B, 0xB749FAED14125D36}, <span class="comment">// 1e-250</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	{0x82B7E12780E7401A, 0xE51C79A85916F484}, <span class="comment">// 1e-249</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	{0xD1B2ECB8B0908810, 0x8F31CC0937AE58D2}, <span class="comment">// 1e-248</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	{0x861FA7E6DCB4AA15, 0xB2FE3F0B8599EF07}, <span class="comment">// 1e-247</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	{0x67A791E093E1D49A, 0xDFBDCECE67006AC9}, <span class="comment">// 1e-246</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	{0xE0C8BB2C5C6D24E0, 0x8BD6A141006042BD}, <span class="comment">// 1e-245</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	{0x58FAE9F773886E18, 0xAECC49914078536D}, <span class="comment">// 1e-244</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	{0xAF39A475506A899E, 0xDA7F5BF590966848}, <span class="comment">// 1e-243</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	{0x6D8406C952429603, 0x888F99797A5E012D}, <span class="comment">// 1e-242</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	{0xC8E5087BA6D33B83, 0xAAB37FD7D8F58178}, <span class="comment">// 1e-241</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	{0xFB1E4A9A90880A64, 0xD5605FCDCF32E1D6}, <span class="comment">// 1e-240</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	{0x5CF2EEA09A55067F, 0x855C3BE0A17FCD26}, <span class="comment">// 1e-239</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	{0xF42FAA48C0EA481E, 0xA6B34AD8C9DFC06F}, <span class="comment">// 1e-238</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	{0xF13B94DAF124DA26, 0xD0601D8EFC57B08B}, <span class="comment">// 1e-237</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	{0x76C53D08D6B70858, 0x823C12795DB6CE57}, <span class="comment">// 1e-236</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	{0x54768C4B0C64CA6E, 0xA2CB1717B52481ED}, <span class="comment">// 1e-235</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	{0xA9942F5DCF7DFD09, 0xCB7DDCDDA26DA268}, <span class="comment">// 1e-234</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	{0xD3F93B35435D7C4C, 0xFE5D54150B090B02}, <span class="comment">// 1e-233</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	{0xC47BC5014A1A6DAF, 0x9EFA548D26E5A6E1}, <span class="comment">// 1e-232</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	{0x359AB6419CA1091B, 0xC6B8E9B0709F109A}, <span class="comment">// 1e-231</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	{0xC30163D203C94B62, 0xF867241C8CC6D4C0}, <span class="comment">// 1e-230</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	{0x79E0DE63425DCF1D, 0x9B407691D7FC44F8}, <span class="comment">// 1e-229</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	{0x985915FC12F542E4, 0xC21094364DFB5636}, <span class="comment">// 1e-228</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	{0x3E6F5B7B17B2939D, 0xF294B943E17A2BC4}, <span class="comment">// 1e-227</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	{0xA705992CEECF9C42, 0x979CF3CA6CEC5B5A}, <span class="comment">// 1e-226</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	{0x50C6FF782A838353, 0xBD8430BD08277231}, <span class="comment">// 1e-225</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	{0xA4F8BF5635246428, 0xECE53CEC4A314EBD}, <span class="comment">// 1e-224</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	{0x871B7795E136BE99, 0x940F4613AE5ED136}, <span class="comment">// 1e-223</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	{0x28E2557B59846E3F, 0xB913179899F68584}, <span class="comment">// 1e-222</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	{0x331AEADA2FE589CF, 0xE757DD7EC07426E5}, <span class="comment">// 1e-221</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	{0x3FF0D2C85DEF7621, 0x9096EA6F3848984F}, <span class="comment">// 1e-220</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	{0x0FED077A756B53A9, 0xB4BCA50B065ABE63}, <span class="comment">// 1e-219</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	{0xD3E8495912C62894, 0xE1EBCE4DC7F16DFB}, <span class="comment">// 1e-218</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	{0x64712DD7ABBBD95C, 0x8D3360F09CF6E4BD}, <span class="comment">// 1e-217</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	{0xBD8D794D96AACFB3, 0xB080392CC4349DEC}, <span class="comment">// 1e-216</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	{0xECF0D7A0FC5583A0, 0xDCA04777F541C567}, <span class="comment">// 1e-215</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	{0xF41686C49DB57244, 0x89E42CAAF9491B60}, <span class="comment">// 1e-214</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	{0x311C2875C522CED5, 0xAC5D37D5B79B6239}, <span class="comment">// 1e-213</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	{0x7D633293366B828B, 0xD77485CB25823AC7}, <span class="comment">// 1e-212</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	{0xAE5DFF9C02033197, 0x86A8D39EF77164BC}, <span class="comment">// 1e-211</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	{0xD9F57F830283FDFC, 0xA8530886B54DBDEB}, <span class="comment">// 1e-210</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	{0xD072DF63C324FD7B, 0xD267CAA862A12D66}, <span class="comment">// 1e-209</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	{0x4247CB9E59F71E6D, 0x8380DEA93DA4BC60}, <span class="comment">// 1e-208</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	{0x52D9BE85F074E608, 0xA46116538D0DEB78}, <span class="comment">// 1e-207</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	{0x67902E276C921F8B, 0xCD795BE870516656}, <span class="comment">// 1e-206</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	{0x00BA1CD8A3DB53B6, 0x806BD9714632DFF6}, <span class="comment">// 1e-205</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	{0x80E8A40ECCD228A4, 0xA086CFCD97BF97F3}, <span class="comment">// 1e-204</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	{0x6122CD128006B2CD, 0xC8A883C0FDAF7DF0}, <span class="comment">// 1e-203</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	{0x796B805720085F81, 0xFAD2A4B13D1B5D6C}, <span class="comment">// 1e-202</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	{0xCBE3303674053BB0, 0x9CC3A6EEC6311A63}, <span class="comment">// 1e-201</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	{0xBEDBFC4411068A9C, 0xC3F490AA77BD60FC}, <span class="comment">// 1e-200</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	{0xEE92FB5515482D44, 0xF4F1B4D515ACB93B}, <span class="comment">// 1e-199</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	{0x751BDD152D4D1C4A, 0x991711052D8BF3C5}, <span class="comment">// 1e-198</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	{0xD262D45A78A0635D, 0xBF5CD54678EEF0B6}, <span class="comment">// 1e-197</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	{0x86FB897116C87C34, 0xEF340A98172AACE4}, <span class="comment">// 1e-196</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	{0xD45D35E6AE3D4DA0, 0x9580869F0E7AAC0E}, <span class="comment">// 1e-195</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	{0x8974836059CCA109, 0xBAE0A846D2195712}, <span class="comment">// 1e-194</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	{0x2BD1A438703FC94B, 0xE998D258869FACD7}, <span class="comment">// 1e-193</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	{0x7B6306A34627DDCF, 0x91FF83775423CC06}, <span class="comment">// 1e-192</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	{0x1A3BC84C17B1D542, 0xB67F6455292CBF08}, <span class="comment">// 1e-191</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	{0x20CABA5F1D9E4A93, 0xE41F3D6A7377EECA}, <span class="comment">// 1e-190</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	{0x547EB47B7282EE9C, 0x8E938662882AF53E}, <span class="comment">// 1e-189</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	{0xE99E619A4F23AA43, 0xB23867FB2A35B28D}, <span class="comment">// 1e-188</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	{0x6405FA00E2EC94D4, 0xDEC681F9F4C31F31}, <span class="comment">// 1e-187</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	{0xDE83BC408DD3DD04, 0x8B3C113C38F9F37E}, <span class="comment">// 1e-186</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	{0x9624AB50B148D445, 0xAE0B158B4738705E}, <span class="comment">// 1e-185</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	{0x3BADD624DD9B0957, 0xD98DDAEE19068C76}, <span class="comment">// 1e-184</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	{0xE54CA5D70A80E5D6, 0x87F8A8D4CFA417C9}, <span class="comment">// 1e-183</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	{0x5E9FCF4CCD211F4C, 0xA9F6D30A038D1DBC}, <span class="comment">// 1e-182</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	{0x7647C3200069671F, 0xD47487CC8470652B}, <span class="comment">// 1e-181</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	{0x29ECD9F40041E073, 0x84C8D4DFD2C63F3B}, <span class="comment">// 1e-180</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	{0xF468107100525890, 0xA5FB0A17C777CF09}, <span class="comment">// 1e-179</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	{0x7182148D4066EEB4, 0xCF79CC9DB955C2CC}, <span class="comment">// 1e-178</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	{0xC6F14CD848405530, 0x81AC1FE293D599BF}, <span class="comment">// 1e-177</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	{0xB8ADA00E5A506A7C, 0xA21727DB38CB002F}, <span class="comment">// 1e-176</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	{0xA6D90811F0E4851C, 0xCA9CF1D206FDC03B}, <span class="comment">// 1e-175</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	{0x908F4A166D1DA663, 0xFD442E4688BD304A}, <span class="comment">// 1e-174</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	{0x9A598E4E043287FE, 0x9E4A9CEC15763E2E}, <span class="comment">// 1e-173</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	{0x40EFF1E1853F29FD, 0xC5DD44271AD3CDBA}, <span class="comment">// 1e-172</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	{0xD12BEE59E68EF47C, 0xF7549530E188C128}, <span class="comment">// 1e-171</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	{0x82BB74F8301958CE, 0x9A94DD3E8CF578B9}, <span class="comment">// 1e-170</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	{0xE36A52363C1FAF01, 0xC13A148E3032D6E7}, <span class="comment">// 1e-169</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	{0xDC44E6C3CB279AC1, 0xF18899B1BC3F8CA1}, <span class="comment">// 1e-168</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	{0x29AB103A5EF8C0B9, 0x96F5600F15A7B7E5}, <span class="comment">// 1e-167</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	{0x7415D448F6B6F0E7, 0xBCB2B812DB11A5DE}, <span class="comment">// 1e-166</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	{0x111B495B3464AD21, 0xEBDF661791D60F56}, <span class="comment">// 1e-165</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	{0xCAB10DD900BEEC34, 0x936B9FCEBB25C995}, <span class="comment">// 1e-164</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	{0x3D5D514F40EEA742, 0xB84687C269EF3BFB}, <span class="comment">// 1e-163</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	{0x0CB4A5A3112A5112, 0xE65829B3046B0AFA}, <span class="comment">// 1e-162</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	{0x47F0E785EABA72AB, 0x8FF71A0FE2C2E6DC}, <span class="comment">// 1e-161</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	{0x59ED216765690F56, 0xB3F4E093DB73A093}, <span class="comment">// 1e-160</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	{0x306869C13EC3532C, 0xE0F218B8D25088B8}, <span class="comment">// 1e-159</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	{0x1E414218C73A13FB, 0x8C974F7383725573}, <span class="comment">// 1e-158</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	{0xE5D1929EF90898FA, 0xAFBD2350644EEACF}, <span class="comment">// 1e-157</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	{0xDF45F746B74ABF39, 0xDBAC6C247D62A583}, <span class="comment">// 1e-156</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	{0x6B8BBA8C328EB783, 0x894BC396CE5DA772}, <span class="comment">// 1e-155</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	{0x066EA92F3F326564, 0xAB9EB47C81F5114F}, <span class="comment">// 1e-154</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	{0xC80A537B0EFEFEBD, 0xD686619BA27255A2}, <span class="comment">// 1e-153</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	{0xBD06742CE95F5F36, 0x8613FD0145877585}, <span class="comment">// 1e-152</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	{0x2C48113823B73704, 0xA798FC4196E952E7}, <span class="comment">// 1e-151</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	{0xF75A15862CA504C5, 0xD17F3B51FCA3A7A0}, <span class="comment">// 1e-150</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	{0x9A984D73DBE722FB, 0x82EF85133DE648C4}, <span class="comment">// 1e-149</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	{0xC13E60D0D2E0EBBA, 0xA3AB66580D5FDAF5}, <span class="comment">// 1e-148</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	{0x318DF905079926A8, 0xCC963FEE10B7D1B3}, <span class="comment">// 1e-147</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	{0xFDF17746497F7052, 0xFFBBCFE994E5C61F}, <span class="comment">// 1e-146</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	{0xFEB6EA8BEDEFA633, 0x9FD561F1FD0F9BD3}, <span class="comment">// 1e-145</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	{0xFE64A52EE96B8FC0, 0xC7CABA6E7C5382C8}, <span class="comment">// 1e-144</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	{0x3DFDCE7AA3C673B0, 0xF9BD690A1B68637B}, <span class="comment">// 1e-143</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	{0x06BEA10CA65C084E, 0x9C1661A651213E2D}, <span class="comment">// 1e-142</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	{0x486E494FCFF30A62, 0xC31BFA0FE5698DB8}, <span class="comment">// 1e-141</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	{0x5A89DBA3C3EFCCFA, 0xF3E2F893DEC3F126}, <span class="comment">// 1e-140</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	{0xF89629465A75E01C, 0x986DDB5C6B3A76B7}, <span class="comment">// 1e-139</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	{0xF6BBB397F1135823, 0xBE89523386091465}, <span class="comment">// 1e-138</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	{0x746AA07DED582E2C, 0xEE2BA6C0678B597F}, <span class="comment">// 1e-137</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	{0xA8C2A44EB4571CDC, 0x94DB483840B717EF}, <span class="comment">// 1e-136</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	{0x92F34D62616CE413, 0xBA121A4650E4DDEB}, <span class="comment">// 1e-135</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	{0x77B020BAF9C81D17, 0xE896A0D7E51E1566}, <span class="comment">// 1e-134</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	{0x0ACE1474DC1D122E, 0x915E2486EF32CD60}, <span class="comment">// 1e-133</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	{0x0D819992132456BA, 0xB5B5ADA8AAFF80B8}, <span class="comment">// 1e-132</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	{0x10E1FFF697ED6C69, 0xE3231912D5BF60E6}, <span class="comment">// 1e-131</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	{0xCA8D3FFA1EF463C1, 0x8DF5EFABC5979C8F}, <span class="comment">// 1e-130</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	{0xBD308FF8A6B17CB2, 0xB1736B96B6FD83B3}, <span class="comment">// 1e-129</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	{0xAC7CB3F6D05DDBDE, 0xDDD0467C64BCE4A0}, <span class="comment">// 1e-128</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	{0x6BCDF07A423AA96B, 0x8AA22C0DBEF60EE4}, <span class="comment">// 1e-127</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	{0x86C16C98D2C953C6, 0xAD4AB7112EB3929D}, <span class="comment">// 1e-126</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	{0xE871C7BF077BA8B7, 0xD89D64D57A607744}, <span class="comment">// 1e-125</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	{0x11471CD764AD4972, 0x87625F056C7C4A8B}, <span class="comment">// 1e-124</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	{0xD598E40D3DD89BCF, 0xA93AF6C6C79B5D2D}, <span class="comment">// 1e-123</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	{0x4AFF1D108D4EC2C3, 0xD389B47879823479}, <span class="comment">// 1e-122</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	{0xCEDF722A585139BA, 0x843610CB4BF160CB}, <span class="comment">// 1e-121</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	{0xC2974EB4EE658828, 0xA54394FE1EEDB8FE}, <span class="comment">// 1e-120</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	{0x733D226229FEEA32, 0xCE947A3DA6A9273E}, <span class="comment">// 1e-119</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	{0x0806357D5A3F525F, 0x811CCC668829B887}, <span class="comment">// 1e-118</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	{0xCA07C2DCB0CF26F7, 0xA163FF802A3426A8}, <span class="comment">// 1e-117</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	{0xFC89B393DD02F0B5, 0xC9BCFF6034C13052}, <span class="comment">// 1e-116</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	{0xBBAC2078D443ACE2, 0xFC2C3F3841F17C67}, <span class="comment">// 1e-115</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	{0xD54B944B84AA4C0D, 0x9D9BA7832936EDC0}, <span class="comment">// 1e-114</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	{0x0A9E795E65D4DF11, 0xC5029163F384A931}, <span class="comment">// 1e-113</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	{0x4D4617B5FF4A16D5, 0xF64335BCF065D37D}, <span class="comment">// 1e-112</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	{0x504BCED1BF8E4E45, 0x99EA0196163FA42E}, <span class="comment">// 1e-111</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	{0xE45EC2862F71E1D6, 0xC06481FB9BCF8D39}, <span class="comment">// 1e-110</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	{0x5D767327BB4E5A4C, 0xF07DA27A82C37088}, <span class="comment">// 1e-109</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	{0x3A6A07F8D510F86F, 0x964E858C91BA2655}, <span class="comment">// 1e-108</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	{0x890489F70A55368B, 0xBBE226EFB628AFEA}, <span class="comment">// 1e-107</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	{0x2B45AC74CCEA842E, 0xEADAB0ABA3B2DBE5}, <span class="comment">// 1e-106</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	{0x3B0B8BC90012929D, 0x92C8AE6B464FC96F}, <span class="comment">// 1e-105</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	{0x09CE6EBB40173744, 0xB77ADA0617E3BBCB}, <span class="comment">// 1e-104</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	{0xCC420A6A101D0515, 0xE55990879DDCAABD}, <span class="comment">// 1e-103</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	{0x9FA946824A12232D, 0x8F57FA54C2A9EAB6}, <span class="comment">// 1e-102</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	{0x47939822DC96ABF9, 0xB32DF8E9F3546564}, <span class="comment">// 1e-101</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	{0x59787E2B93BC56F7, 0xDFF9772470297EBD}, <span class="comment">// 1e-100</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	{0x57EB4EDB3C55B65A, 0x8BFBEA76C619EF36}, <span class="comment">// 1e-99</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	{0xEDE622920B6B23F1, 0xAEFAE51477A06B03}, <span class="comment">// 1e-98</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	{0xE95FAB368E45ECED, 0xDAB99E59958885C4}, <span class="comment">// 1e-97</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	{0x11DBCB0218EBB414, 0x88B402F7FD75539B}, <span class="comment">// 1e-96</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	{0xD652BDC29F26A119, 0xAAE103B5FCD2A881}, <span class="comment">// 1e-95</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	{0x4BE76D3346F0495F, 0xD59944A37C0752A2}, <span class="comment">// 1e-94</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	{0x6F70A4400C562DDB, 0x857FCAE62D8493A5}, <span class="comment">// 1e-93</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	{0xCB4CCD500F6BB952, 0xA6DFBD9FB8E5B88E}, <span class="comment">// 1e-92</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	{0x7E2000A41346A7A7, 0xD097AD07A71F26B2}, <span class="comment">// 1e-91</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	{0x8ED400668C0C28C8, 0x825ECC24C873782F}, <span class="comment">// 1e-90</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	{0x728900802F0F32FA, 0xA2F67F2DFA90563B}, <span class="comment">// 1e-89</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	{0x4F2B40A03AD2FFB9, 0xCBB41EF979346BCA}, <span class="comment">// 1e-88</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	{0xE2F610C84987BFA8, 0xFEA126B7D78186BC}, <span class="comment">// 1e-87</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	{0x0DD9CA7D2DF4D7C9, 0x9F24B832E6B0F436}, <span class="comment">// 1e-86</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	{0x91503D1C79720DBB, 0xC6EDE63FA05D3143}, <span class="comment">// 1e-85</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	{0x75A44C6397CE912A, 0xF8A95FCF88747D94}, <span class="comment">// 1e-84</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	{0xC986AFBE3EE11ABA, 0x9B69DBE1B548CE7C}, <span class="comment">// 1e-83</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	{0xFBE85BADCE996168, 0xC24452DA229B021B}, <span class="comment">// 1e-82</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	{0xFAE27299423FB9C3, 0xF2D56790AB41C2A2}, <span class="comment">// 1e-81</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	{0xDCCD879FC967D41A, 0x97C560BA6B0919A5}, <span class="comment">// 1e-80</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	{0x5400E987BBC1C920, 0xBDB6B8E905CB600F}, <span class="comment">// 1e-79</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	{0x290123E9AAB23B68, 0xED246723473E3813}, <span class="comment">// 1e-78</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	{0xF9A0B6720AAF6521, 0x9436C0760C86E30B}, <span class="comment">// 1e-77</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	{0xF808E40E8D5B3E69, 0xB94470938FA89BCE}, <span class="comment">// 1e-76</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	{0xB60B1D1230B20E04, 0xE7958CB87392C2C2}, <span class="comment">// 1e-75</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	{0xB1C6F22B5E6F48C2, 0x90BD77F3483BB9B9}, <span class="comment">// 1e-74</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	{0x1E38AEB6360B1AF3, 0xB4ECD5F01A4AA828}, <span class="comment">// 1e-73</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	{0x25C6DA63C38DE1B0, 0xE2280B6C20DD5232}, <span class="comment">// 1e-72</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	{0x579C487E5A38AD0E, 0x8D590723948A535F}, <span class="comment">// 1e-71</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	{0x2D835A9DF0C6D851, 0xB0AF48EC79ACE837}, <span class="comment">// 1e-70</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	{0xF8E431456CF88E65, 0xDCDB1B2798182244}, <span class="comment">// 1e-69</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	{0x1B8E9ECB641B58FF, 0x8A08F0F8BF0F156B}, <span class="comment">// 1e-68</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	{0xE272467E3D222F3F, 0xAC8B2D36EED2DAC5}, <span class="comment">// 1e-67</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	{0x5B0ED81DCC6ABB0F, 0xD7ADF884AA879177}, <span class="comment">// 1e-66</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	{0x98E947129FC2B4E9, 0x86CCBB52EA94BAEA}, <span class="comment">// 1e-65</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	{0x3F2398D747B36224, 0xA87FEA27A539E9A5}, <span class="comment">// 1e-64</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	{0x8EEC7F0D19A03AAD, 0xD29FE4B18E88640E}, <span class="comment">// 1e-63</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	{0x1953CF68300424AC, 0x83A3EEEEF9153E89}, <span class="comment">// 1e-62</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	{0x5FA8C3423C052DD7, 0xA48CEAAAB75A8E2B}, <span class="comment">// 1e-61</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	{0x3792F412CB06794D, 0xCDB02555653131B6}, <span class="comment">// 1e-60</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	{0xE2BBD88BBEE40BD0, 0x808E17555F3EBF11}, <span class="comment">// 1e-59</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	{0x5B6ACEAEAE9D0EC4, 0xA0B19D2AB70E6ED6}, <span class="comment">// 1e-58</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	{0xF245825A5A445275, 0xC8DE047564D20A8B}, <span class="comment">// 1e-57</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	{0xEED6E2F0F0D56712, 0xFB158592BE068D2E}, <span class="comment">// 1e-56</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	{0x55464DD69685606B, 0x9CED737BB6C4183D}, <span class="comment">// 1e-55</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	{0xAA97E14C3C26B886, 0xC428D05AA4751E4C}, <span class="comment">// 1e-54</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	{0xD53DD99F4B3066A8, 0xF53304714D9265DF}, <span class="comment">// 1e-53</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	{0xE546A8038EFE4029, 0x993FE2C6D07B7FAB}, <span class="comment">// 1e-52</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	{0xDE98520472BDD033, 0xBF8FDB78849A5F96}, <span class="comment">// 1e-51</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	{0x963E66858F6D4440, 0xEF73D256A5C0F77C}, <span class="comment">// 1e-50</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	{0xDDE7001379A44AA8, 0x95A8637627989AAD}, <span class="comment">// 1e-49</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	{0x5560C018580D5D52, 0xBB127C53B17EC159}, <span class="comment">// 1e-48</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	{0xAAB8F01E6E10B4A6, 0xE9D71B689DDE71AF}, <span class="comment">// 1e-47</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	{0xCAB3961304CA70E8, 0x9226712162AB070D}, <span class="comment">// 1e-46</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	{0x3D607B97C5FD0D22, 0xB6B00D69BB55C8D1}, <span class="comment">// 1e-45</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	{0x8CB89A7DB77C506A, 0xE45C10C42A2B3B05}, <span class="comment">// 1e-44</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	{0x77F3608E92ADB242, 0x8EB98A7A9A5B04E3}, <span class="comment">// 1e-43</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	{0x55F038B237591ED3, 0xB267ED1940F1C61C}, <span class="comment">// 1e-42</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	{0x6B6C46DEC52F6688, 0xDF01E85F912E37A3}, <span class="comment">// 1e-41</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	{0x2323AC4B3B3DA015, 0x8B61313BBABCE2C6}, <span class="comment">// 1e-40</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	{0xABEC975E0A0D081A, 0xAE397D8AA96C1B77}, <span class="comment">// 1e-39</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	{0x96E7BD358C904A21, 0xD9C7DCED53C72255}, <span class="comment">// 1e-38</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	{0x7E50D64177DA2E54, 0x881CEA14545C7575}, <span class="comment">// 1e-37</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	{0xDDE50BD1D5D0B9E9, 0xAA242499697392D2}, <span class="comment">// 1e-36</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	{0x955E4EC64B44E864, 0xD4AD2DBFC3D07787}, <span class="comment">// 1e-35</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	{0xBD5AF13BEF0B113E, 0x84EC3C97DA624AB4}, <span class="comment">// 1e-34</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	{0xECB1AD8AEACDD58E, 0xA6274BBDD0FADD61}, <span class="comment">// 1e-33</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	{0x67DE18EDA5814AF2, 0xCFB11EAD453994BA}, <span class="comment">// 1e-32</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	{0x80EACF948770CED7, 0x81CEB32C4B43FCF4}, <span class="comment">// 1e-31</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	{0xA1258379A94D028D, 0xA2425FF75E14FC31}, <span class="comment">// 1e-30</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	{0x096EE45813A04330, 0xCAD2F7F5359A3B3E}, <span class="comment">// 1e-29</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	{0x8BCA9D6E188853FC, 0xFD87B5F28300CA0D}, <span class="comment">// 1e-28</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	{0x775EA264CF55347D, 0x9E74D1B791E07E48}, <span class="comment">// 1e-27</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	{0x95364AFE032A819D, 0xC612062576589DDA}, <span class="comment">// 1e-26</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	{0x3A83DDBD83F52204, 0xF79687AED3EEC551}, <span class="comment">// 1e-25</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	{0xC4926A9672793542, 0x9ABE14CD44753B52}, <span class="comment">// 1e-24</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	{0x75B7053C0F178293, 0xC16D9A0095928A27}, <span class="comment">// 1e-23</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	{0x5324C68B12DD6338, 0xF1C90080BAF72CB1}, <span class="comment">// 1e-22</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	{0xD3F6FC16EBCA5E03, 0x971DA05074DA7BEE}, <span class="comment">// 1e-21</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	{0x88F4BB1CA6BCF584, 0xBCE5086492111AEA}, <span class="comment">// 1e-20</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	{0x2B31E9E3D06C32E5, 0xEC1E4A7DB69561A5}, <span class="comment">// 1e-19</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	{0x3AFF322E62439FCF, 0x9392EE8E921D5D07}, <span class="comment">// 1e-18</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	{0x09BEFEB9FAD487C2, 0xB877AA3236A4B449}, <span class="comment">// 1e-17</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	{0x4C2EBE687989A9B3, 0xE69594BEC44DE15B}, <span class="comment">// 1e-16</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	{0x0F9D37014BF60A10, 0x901D7CF73AB0ACD9}, <span class="comment">// 1e-15</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	{0x538484C19EF38C94, 0xB424DC35095CD80F}, <span class="comment">// 1e-14</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	{0x2865A5F206B06FB9, 0xE12E13424BB40E13}, <span class="comment">// 1e-13</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	{0xF93F87B7442E45D3, 0x8CBCCC096F5088CB}, <span class="comment">// 1e-12</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	{0xF78F69A51539D748, 0xAFEBFF0BCB24AAFE}, <span class="comment">// 1e-11</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	{0xB573440E5A884D1B, 0xDBE6FECEBDEDD5BE}, <span class="comment">// 1e-10</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	{0x31680A88F8953030, 0x89705F4136B4A597}, <span class="comment">// 1e-9</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	{0xFDC20D2B36BA7C3D, 0xABCC77118461CEFC}, <span class="comment">// 1e-8</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	{0x3D32907604691B4C, 0xD6BF94D5E57A42BC}, <span class="comment">// 1e-7</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	{0xA63F9A49C2C1B10F, 0x8637BD05AF6C69B5}, <span class="comment">// 1e-6</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	{0x0FCF80DC33721D53, 0xA7C5AC471B478423}, <span class="comment">// 1e-5</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	{0xD3C36113404EA4A8, 0xD1B71758E219652B}, <span class="comment">// 1e-4</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	{0x645A1CAC083126E9, 0x83126E978D4FDF3B}, <span class="comment">// 1e-3</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	{0x3D70A3D70A3D70A3, 0xA3D70A3D70A3D70A}, <span class="comment">// 1e-2</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	{0xCCCCCCCCCCCCCCCC, 0xCCCCCCCCCCCCCCCC}, <span class="comment">// 1e-1</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	{0x0000000000000000, 0x8000000000000000}, <span class="comment">// 1e0</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	{0x0000000000000000, 0xA000000000000000}, <span class="comment">// 1e1</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	{0x0000000000000000, 0xC800000000000000}, <span class="comment">// 1e2</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	{0x0000000000000000, 0xFA00000000000000}, <span class="comment">// 1e3</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	{0x0000000000000000, 0x9C40000000000000}, <span class="comment">// 1e4</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	{0x0000000000000000, 0xC350000000000000}, <span class="comment">// 1e5</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	{0x0000000000000000, 0xF424000000000000}, <span class="comment">// 1e6</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	{0x0000000000000000, 0x9896800000000000}, <span class="comment">// 1e7</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	{0x0000000000000000, 0xBEBC200000000000}, <span class="comment">// 1e8</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	{0x0000000000000000, 0xEE6B280000000000}, <span class="comment">// 1e9</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	{0x0000000000000000, 0x9502F90000000000}, <span class="comment">// 1e10</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	{0x0000000000000000, 0xBA43B74000000000}, <span class="comment">// 1e11</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	{0x0000000000000000, 0xE8D4A51000000000}, <span class="comment">// 1e12</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	{0x0000000000000000, 0x9184E72A00000000}, <span class="comment">// 1e13</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	{0x0000000000000000, 0xB5E620F480000000}, <span class="comment">// 1e14</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	{0x0000000000000000, 0xE35FA931A0000000}, <span class="comment">// 1e15</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	{0x0000000000000000, 0x8E1BC9BF04000000}, <span class="comment">// 1e16</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	{0x0000000000000000, 0xB1A2BC2EC5000000}, <span class="comment">// 1e17</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	{0x0000000000000000, 0xDE0B6B3A76400000}, <span class="comment">// 1e18</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	{0x0000000000000000, 0x8AC7230489E80000}, <span class="comment">// 1e19</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	{0x0000000000000000, 0xAD78EBC5AC620000}, <span class="comment">// 1e20</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	{0x0000000000000000, 0xD8D726B7177A8000}, <span class="comment">// 1e21</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	{0x0000000000000000, 0x878678326EAC9000}, <span class="comment">// 1e22</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	{0x0000000000000000, 0xA968163F0A57B400}, <span class="comment">// 1e23</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	{0x0000000000000000, 0xD3C21BCECCEDA100}, <span class="comment">// 1e24</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	{0x0000000000000000, 0x84595161401484A0}, <span class="comment">// 1e25</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	{0x0000000000000000, 0xA56FA5B99019A5C8}, <span class="comment">// 1e26</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	{0x0000000000000000, 0xCECB8F27F4200F3A}, <span class="comment">// 1e27</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	{0x4000000000000000, 0x813F3978F8940984}, <span class="comment">// 1e28</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	{0x5000000000000000, 0xA18F07D736B90BE5}, <span class="comment">// 1e29</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	{0xA400000000000000, 0xC9F2C9CD04674EDE}, <span class="comment">// 1e30</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	{0x4D00000000000000, 0xFC6F7C4045812296}, <span class="comment">// 1e31</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	{0xF020000000000000, 0x9DC5ADA82B70B59D}, <span class="comment">// 1e32</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	{0x6C28000000000000, 0xC5371912364CE305}, <span class="comment">// 1e33</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	{0xC732000000000000, 0xF684DF56C3E01BC6}, <span class="comment">// 1e34</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	{0x3C7F400000000000, 0x9A130B963A6C115C}, <span class="comment">// 1e35</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	{0x4B9F100000000000, 0xC097CE7BC90715B3}, <span class="comment">// 1e36</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	{0x1E86D40000000000, 0xF0BDC21ABB48DB20}, <span class="comment">// 1e37</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	{0x1314448000000000, 0x96769950B50D88F4}, <span class="comment">// 1e38</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	{0x17D955A000000000, 0xBC143FA4E250EB31}, <span class="comment">// 1e39</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	{0x5DCFAB0800000000, 0xEB194F8E1AE525FD}, <span class="comment">// 1e40</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	{0x5AA1CAE500000000, 0x92EFD1B8D0CF37BE}, <span class="comment">// 1e41</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	{0xF14A3D9E40000000, 0xB7ABC627050305AD}, <span class="comment">// 1e42</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	{0x6D9CCD05D0000000, 0xE596B7B0C643C719}, <span class="comment">// 1e43</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	{0xE4820023A2000000, 0x8F7E32CE7BEA5C6F}, <span class="comment">// 1e44</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	{0xDDA2802C8A800000, 0xB35DBF821AE4F38B}, <span class="comment">// 1e45</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	{0xD50B2037AD200000, 0xE0352F62A19E306E}, <span class="comment">// 1e46</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	{0x4526F422CC340000, 0x8C213D9DA502DE45}, <span class="comment">// 1e47</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	{0x9670B12B7F410000, 0xAF298D050E4395D6}, <span class="comment">// 1e48</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	{0x3C0CDD765F114000, 0xDAF3F04651D47B4C}, <span class="comment">// 1e49</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	{0xA5880A69FB6AC800, 0x88D8762BF324CD0F}, <span class="comment">// 1e50</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	{0x8EEA0D047A457A00, 0xAB0E93B6EFEE0053}, <span class="comment">// 1e51</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	{0x72A4904598D6D880, 0xD5D238A4ABE98068}, <span class="comment">// 1e52</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	{0x47A6DA2B7F864750, 0x85A36366EB71F041}, <span class="comment">// 1e53</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	{0x999090B65F67D924, 0xA70C3C40A64E6C51}, <span class="comment">// 1e54</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	{0xFFF4B4E3F741CF6D, 0xD0CF4B50CFE20765}, <span class="comment">// 1e55</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	{0xBFF8F10E7A8921A4, 0x82818F1281ED449F}, <span class="comment">// 1e56</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	{0xAFF72D52192B6A0D, 0xA321F2D7226895C7}, <span class="comment">// 1e57</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	{0x9BF4F8A69F764490, 0xCBEA6F8CEB02BB39}, <span class="comment">// 1e58</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	{0x02F236D04753D5B4, 0xFEE50B7025C36A08}, <span class="comment">// 1e59</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	{0x01D762422C946590, 0x9F4F2726179A2245}, <span class="comment">// 1e60</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	{0x424D3AD2B7B97EF5, 0xC722F0EF9D80AAD6}, <span class="comment">// 1e61</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	{0xD2E0898765A7DEB2, 0xF8EBAD2B84E0D58B}, <span class="comment">// 1e62</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	{0x63CC55F49F88EB2F, 0x9B934C3B330C8577}, <span class="comment">// 1e63</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	{0x3CBF6B71C76B25FB, 0xC2781F49FFCFA6D5}, <span class="comment">// 1e64</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	{0x8BEF464E3945EF7A, 0xF316271C7FC3908A}, <span class="comment">// 1e65</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	{0x97758BF0E3CBB5AC, 0x97EDD871CFDA3A56}, <span class="comment">// 1e66</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	{0x3D52EEED1CBEA317, 0xBDE94E8E43D0C8EC}, <span class="comment">// 1e67</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	{0x4CA7AAA863EE4BDD, 0xED63A231D4C4FB27}, <span class="comment">// 1e68</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	{0x8FE8CAA93E74EF6A, 0x945E455F24FB1CF8}, <span class="comment">// 1e69</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	{0xB3E2FD538E122B44, 0xB975D6B6EE39E436}, <span class="comment">// 1e70</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	{0x60DBBCA87196B616, 0xE7D34C64A9C85D44}, <span class="comment">// 1e71</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	{0xBC8955E946FE31CD, 0x90E40FBEEA1D3A4A}, <span class="comment">// 1e72</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	{0x6BABAB6398BDBE41, 0xB51D13AEA4A488DD}, <span class="comment">// 1e73</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	{0xC696963C7EED2DD1, 0xE264589A4DCDAB14}, <span class="comment">// 1e74</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	{0xFC1E1DE5CF543CA2, 0x8D7EB76070A08AEC}, <span class="comment">// 1e75</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	{0x3B25A55F43294BCB, 0xB0DE65388CC8ADA8}, <span class="comment">// 1e76</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	{0x49EF0EB713F39EBE, 0xDD15FE86AFFAD912}, <span class="comment">// 1e77</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	{0x6E3569326C784337, 0x8A2DBF142DFCC7AB}, <span class="comment">// 1e78</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	{0x49C2C37F07965404, 0xACB92ED9397BF996}, <span class="comment">// 1e79</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	{0xDC33745EC97BE906, 0xD7E77A8F87DAF7FB}, <span class="comment">// 1e80</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	{0x69A028BB3DED71A3, 0x86F0AC99B4E8DAFD}, <span class="comment">// 1e81</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	{0xC40832EA0D68CE0C, 0xA8ACD7C0222311BC}, <span class="comment">// 1e82</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	{0xF50A3FA490C30190, 0xD2D80DB02AABD62B}, <span class="comment">// 1e83</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	{0x792667C6DA79E0FA, 0x83C7088E1AAB65DB}, <span class="comment">// 1e84</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	{0x577001B891185938, 0xA4B8CAB1A1563F52}, <span class="comment">// 1e85</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	{0xED4C0226B55E6F86, 0xCDE6FD5E09ABCF26}, <span class="comment">// 1e86</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	{0x544F8158315B05B4, 0x80B05E5AC60B6178}, <span class="comment">// 1e87</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	{0x696361AE3DB1C721, 0xA0DC75F1778E39D6}, <span class="comment">// 1e88</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	{0x03BC3A19CD1E38E9, 0xC913936DD571C84C}, <span class="comment">// 1e89</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	{0x04AB48A04065C723, 0xFB5878494ACE3A5F}, <span class="comment">// 1e90</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	{0x62EB0D64283F9C76, 0x9D174B2DCEC0E47B}, <span class="comment">// 1e91</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	{0x3BA5D0BD324F8394, 0xC45D1DF942711D9A}, <span class="comment">// 1e92</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	{0xCA8F44EC7EE36479, 0xF5746577930D6500}, <span class="comment">// 1e93</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	{0x7E998B13CF4E1ECB, 0x9968BF6ABBE85F20}, <span class="comment">// 1e94</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	{0x9E3FEDD8C321A67E, 0xBFC2EF456AE276E8}, <span class="comment">// 1e95</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	{0xC5CFE94EF3EA101E, 0xEFB3AB16C59B14A2}, <span class="comment">// 1e96</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	{0xBBA1F1D158724A12, 0x95D04AEE3B80ECE5}, <span class="comment">// 1e97</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	{0x2A8A6E45AE8EDC97, 0xBB445DA9CA61281F}, <span class="comment">// 1e98</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	{0xF52D09D71A3293BD, 0xEA1575143CF97226}, <span class="comment">// 1e99</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	{0x593C2626705F9C56, 0x924D692CA61BE758}, <span class="comment">// 1e100</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	{0x6F8B2FB00C77836C, 0xB6E0C377CFA2E12E}, <span class="comment">// 1e101</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	{0x0B6DFB9C0F956447, 0xE498F455C38B997A}, <span class="comment">// 1e102</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	{0x4724BD4189BD5EAC, 0x8EDF98B59A373FEC}, <span class="comment">// 1e103</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	{0x58EDEC91EC2CB657, 0xB2977EE300C50FE7}, <span class="comment">// 1e104</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	{0x2F2967B66737E3ED, 0xDF3D5E9BC0F653E1}, <span class="comment">// 1e105</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	{0xBD79E0D20082EE74, 0x8B865B215899F46C}, <span class="comment">// 1e106</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	{0xECD8590680A3AA11, 0xAE67F1E9AEC07187}, <span class="comment">// 1e107</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	{0xE80E6F4820CC9495, 0xDA01EE641A708DE9}, <span class="comment">// 1e108</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	{0x3109058D147FDCDD, 0x884134FE908658B2}, <span class="comment">// 1e109</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	{0xBD4B46F0599FD415, 0xAA51823E34A7EEDE}, <span class="comment">// 1e110</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	{0x6C9E18AC7007C91A, 0xD4E5E2CDC1D1EA96}, <span class="comment">// 1e111</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	{0x03E2CF6BC604DDB0, 0x850FADC09923329E}, <span class="comment">// 1e112</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	{0x84DB8346B786151C, 0xA6539930BF6BFF45}, <span class="comment">// 1e113</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	{0xE612641865679A63, 0xCFE87F7CEF46FF16}, <span class="comment">// 1e114</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	{0x4FCB7E8F3F60C07E, 0x81F14FAE158C5F6E}, <span class="comment">// 1e115</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	{0xE3BE5E330F38F09D, 0xA26DA3999AEF7749}, <span class="comment">// 1e116</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	{0x5CADF5BFD3072CC5, 0xCB090C8001AB551C}, <span class="comment">// 1e117</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	{0x73D9732FC7C8F7F6, 0xFDCB4FA002162A63}, <span class="comment">// 1e118</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	{0x2867E7FDDCDD9AFA, 0x9E9F11C4014DDA7E}, <span class="comment">// 1e119</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	{0xB281E1FD541501B8, 0xC646D63501A1511D}, <span class="comment">// 1e120</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	{0x1F225A7CA91A4226, 0xF7D88BC24209A565}, <span class="comment">// 1e121</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	{0x3375788DE9B06958, 0x9AE757596946075F}, <span class="comment">// 1e122</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	{0x0052D6B1641C83AE, 0xC1A12D2FC3978937}, <span class="comment">// 1e123</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	{0xC0678C5DBD23A49A, 0xF209787BB47D6B84}, <span class="comment">// 1e124</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	{0xF840B7BA963646E0, 0x9745EB4D50CE6332}, <span class="comment">// 1e125</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	{0xB650E5A93BC3D898, 0xBD176620A501FBFF}, <span class="comment">// 1e126</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	{0xA3E51F138AB4CEBE, 0xEC5D3FA8CE427AFF}, <span class="comment">// 1e127</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	{0xC66F336C36B10137, 0x93BA47C980E98CDF}, <span class="comment">// 1e128</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	{0xB80B0047445D4184, 0xB8A8D9BBE123F017}, <span class="comment">// 1e129</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	{0xA60DC059157491E5, 0xE6D3102AD96CEC1D}, <span class="comment">// 1e130</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	{0x87C89837AD68DB2F, 0x9043EA1AC7E41392}, <span class="comment">// 1e131</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	{0x29BABE4598C311FB, 0xB454E4A179DD1877}, <span class="comment">// 1e132</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	{0xF4296DD6FEF3D67A, 0xE16A1DC9D8545E94}, <span class="comment">// 1e133</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	{0x1899E4A65F58660C, 0x8CE2529E2734BB1D}, <span class="comment">// 1e134</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	{0x5EC05DCFF72E7F8F, 0xB01AE745B101E9E4}, <span class="comment">// 1e135</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	{0x76707543F4FA1F73, 0xDC21A1171D42645D}, <span class="comment">// 1e136</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	{0x6A06494A791C53A8, 0x899504AE72497EBA}, <span class="comment">// 1e137</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	{0x0487DB9D17636892, 0xABFA45DA0EDBDE69}, <span class="comment">// 1e138</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	{0x45A9D2845D3C42B6, 0xD6F8D7509292D603}, <span class="comment">// 1e139</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	{0x0B8A2392BA45A9B2, 0x865B86925B9BC5C2}, <span class="comment">// 1e140</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	{0x8E6CAC7768D7141E, 0xA7F26836F282B732}, <span class="comment">// 1e141</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	{0x3207D795430CD926, 0xD1EF0244AF2364FF}, <span class="comment">// 1e142</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	{0x7F44E6BD49E807B8, 0x8335616AED761F1F}, <span class="comment">// 1e143</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	{0x5F16206C9C6209A6, 0xA402B9C5A8D3A6E7}, <span class="comment">// 1e144</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	{0x36DBA887C37A8C0F, 0xCD036837130890A1}, <span class="comment">// 1e145</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	{0xC2494954DA2C9789, 0x802221226BE55A64}, <span class="comment">// 1e146</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	{0xF2DB9BAA10B7BD6C, 0xA02AA96B06DEB0FD}, <span class="comment">// 1e147</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	{0x6F92829494E5ACC7, 0xC83553C5C8965D3D}, <span class="comment">// 1e148</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	{0xCB772339BA1F17F9, 0xFA42A8B73ABBF48C}, <span class="comment">// 1e149</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	{0xFF2A760414536EFB, 0x9C69A97284B578D7}, <span class="comment">// 1e150</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	{0xFEF5138519684ABA, 0xC38413CF25E2D70D}, <span class="comment">// 1e151</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	{0x7EB258665FC25D69, 0xF46518C2EF5B8CD1}, <span class="comment">// 1e152</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	{0xEF2F773FFBD97A61, 0x98BF2F79D5993802}, <span class="comment">// 1e153</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	{0xAAFB550FFACFD8FA, 0xBEEEFB584AFF8603}, <span class="comment">// 1e154</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	{0x95BA2A53F983CF38, 0xEEAABA2E5DBF6784}, <span class="comment">// 1e155</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	{0xDD945A747BF26183, 0x952AB45CFA97A0B2}, <span class="comment">// 1e156</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	{0x94F971119AEEF9E4, 0xBA756174393D88DF}, <span class="comment">// 1e157</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	{0x7A37CD5601AAB85D, 0xE912B9D1478CEB17}, <span class="comment">// 1e158</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	{0xAC62E055C10AB33A, 0x91ABB422CCB812EE}, <span class="comment">// 1e159</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	{0x577B986B314D6009, 0xB616A12B7FE617AA}, <span class="comment">// 1e160</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	{0xED5A7E85FDA0B80B, 0xE39C49765FDF9D94}, <span class="comment">// 1e161</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	{0x14588F13BE847307, 0x8E41ADE9FBEBC27D}, <span class="comment">// 1e162</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	{0x596EB2D8AE258FC8, 0xB1D219647AE6B31C}, <span class="comment">// 1e163</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	{0x6FCA5F8ED9AEF3BB, 0xDE469FBD99A05FE3}, <span class="comment">// 1e164</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	{0x25DE7BB9480D5854, 0x8AEC23D680043BEE}, <span class="comment">// 1e165</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	{0xAF561AA79A10AE6A, 0xADA72CCC20054AE9}, <span class="comment">// 1e166</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	{0x1B2BA1518094DA04, 0xD910F7FF28069DA4}, <span class="comment">// 1e167</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	{0x90FB44D2F05D0842, 0x87AA9AFF79042286}, <span class="comment">// 1e168</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	{0x353A1607AC744A53, 0xA99541BF57452B28}, <span class="comment">// 1e169</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	{0x42889B8997915CE8, 0xD3FA922F2D1675F2}, <span class="comment">// 1e170</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	{0x69956135FEBADA11, 0x847C9B5D7C2E09B7}, <span class="comment">// 1e171</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	{0x43FAB9837E699095, 0xA59BC234DB398C25}, <span class="comment">// 1e172</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	{0x94F967E45E03F4BB, 0xCF02B2C21207EF2E}, <span class="comment">// 1e173</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	{0x1D1BE0EEBAC278F5, 0x8161AFB94B44F57D}, <span class="comment">// 1e174</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	{0x6462D92A69731732, 0xA1BA1BA79E1632DC}, <span class="comment">// 1e175</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	{0x7D7B8F7503CFDCFE, 0xCA28A291859BBF93}, <span class="comment">// 1e176</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	{0x5CDA735244C3D43E, 0xFCB2CB35E702AF78}, <span class="comment">// 1e177</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	{0x3A0888136AFA64A7, 0x9DEFBF01B061ADAB}, <span class="comment">// 1e178</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	{0x088AAA1845B8FDD0, 0xC56BAEC21C7A1916}, <span class="comment">// 1e179</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	{0x8AAD549E57273D45, 0xF6C69A72A3989F5B}, <span class="comment">// 1e180</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	{0x36AC54E2F678864B, 0x9A3C2087A63F6399}, <span class="comment">// 1e181</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	{0x84576A1BB416A7DD, 0xC0CB28A98FCF3C7F}, <span class="comment">// 1e182</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	{0x656D44A2A11C51D5, 0xF0FDF2D3F3C30B9F}, <span class="comment">// 1e183</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	{0x9F644AE5A4B1B325, 0x969EB7C47859E743}, <span class="comment">// 1e184</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	{0x873D5D9F0DDE1FEE, 0xBC4665B596706114}, <span class="comment">// 1e185</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	{0xA90CB506D155A7EA, 0xEB57FF22FC0C7959}, <span class="comment">// 1e186</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	{0x09A7F12442D588F2, 0x9316FF75DD87CBD8}, <span class="comment">// 1e187</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	{0x0C11ED6D538AEB2F, 0xB7DCBF5354E9BECE}, <span class="comment">// 1e188</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	{0x8F1668C8A86DA5FA, 0xE5D3EF282A242E81}, <span class="comment">// 1e189</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	{0xF96E017D694487BC, 0x8FA475791A569D10}, <span class="comment">// 1e190</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>	{0x37C981DCC395A9AC, 0xB38D92D760EC4455}, <span class="comment">// 1e191</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	{0x85BBE253F47B1417, 0xE070F78D3927556A}, <span class="comment">// 1e192</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	{0x93956D7478CCEC8E, 0x8C469AB843B89562}, <span class="comment">// 1e193</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	{0x387AC8D1970027B2, 0xAF58416654A6BABB}, <span class="comment">// 1e194</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	{0x06997B05FCC0319E, 0xDB2E51BFE9D0696A}, <span class="comment">// 1e195</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	{0x441FECE3BDF81F03, 0x88FCF317F22241E2}, <span class="comment">// 1e196</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	{0xD527E81CAD7626C3, 0xAB3C2FDDEEAAD25A}, <span class="comment">// 1e197</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	{0x8A71E223D8D3B074, 0xD60B3BD56A5586F1}, <span class="comment">// 1e198</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	{0xF6872D5667844E49, 0x85C7056562757456}, <span class="comment">// 1e199</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	{0xB428F8AC016561DB, 0xA738C6BEBB12D16C}, <span class="comment">// 1e200</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	{0xE13336D701BEBA52, 0xD106F86E69D785C7}, <span class="comment">// 1e201</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	{0xECC0024661173473, 0x82A45B450226B39C}, <span class="comment">// 1e202</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	{0x27F002D7F95D0190, 0xA34D721642B06084}, <span class="comment">// 1e203</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	{0x31EC038DF7B441F4, 0xCC20CE9BD35C78A5}, <span class="comment">// 1e204</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	{0x7E67047175A15271, 0xFF290242C83396CE}, <span class="comment">// 1e205</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	{0x0F0062C6E984D386, 0x9F79A169BD203E41}, <span class="comment">// 1e206</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	{0x52C07B78A3E60868, 0xC75809C42C684DD1}, <span class="comment">// 1e207</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	{0xA7709A56CCDF8A82, 0xF92E0C3537826145}, <span class="comment">// 1e208</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	{0x88A66076400BB691, 0x9BBCC7A142B17CCB}, <span class="comment">// 1e209</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	{0x6ACFF893D00EA435, 0xC2ABF989935DDBFE}, <span class="comment">// 1e210</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	{0x0583F6B8C4124D43, 0xF356F7EBF83552FE}, <span class="comment">// 1e211</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	{0xC3727A337A8B704A, 0x98165AF37B2153DE}, <span class="comment">// 1e212</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	{0x744F18C0592E4C5C, 0xBE1BF1B059E9A8D6}, <span class="comment">// 1e213</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	{0x1162DEF06F79DF73, 0xEDA2EE1C7064130C}, <span class="comment">// 1e214</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	{0x8ADDCB5645AC2BA8, 0x9485D4D1C63E8BE7}, <span class="comment">// 1e215</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	{0x6D953E2BD7173692, 0xB9A74A0637CE2EE1}, <span class="comment">// 1e216</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	{0xC8FA8DB6CCDD0437, 0xE8111C87C5C1BA99}, <span class="comment">// 1e217</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	{0x1D9C9892400A22A2, 0x910AB1D4DB9914A0}, <span class="comment">// 1e218</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	{0x2503BEB6D00CAB4B, 0xB54D5E4A127F59C8}, <span class="comment">// 1e219</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	{0x2E44AE64840FD61D, 0xE2A0B5DC971F303A}, <span class="comment">// 1e220</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	{0x5CEAECFED289E5D2, 0x8DA471A9DE737E24}, <span class="comment">// 1e221</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	{0x7425A83E872C5F47, 0xB10D8E1456105DAD}, <span class="comment">// 1e222</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	{0xD12F124E28F77719, 0xDD50F1996B947518}, <span class="comment">// 1e223</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	{0x82BD6B70D99AAA6F, 0x8A5296FFE33CC92F}, <span class="comment">// 1e224</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	{0x636CC64D1001550B, 0xACE73CBFDC0BFB7B}, <span class="comment">// 1e225</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	{0x3C47F7E05401AA4E, 0xD8210BEFD30EFA5A}, <span class="comment">// 1e226</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	{0x65ACFAEC34810A71, 0x8714A775E3E95C78}, <span class="comment">// 1e227</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	{0x7F1839A741A14D0D, 0xA8D9D1535CE3B396}, <span class="comment">// 1e228</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	{0x1EDE48111209A050, 0xD31045A8341CA07C}, <span class="comment">// 1e229</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	{0x934AED0AAB460432, 0x83EA2B892091E44D}, <span class="comment">// 1e230</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	{0xF81DA84D5617853F, 0xA4E4B66B68B65D60}, <span class="comment">// 1e231</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	{0x36251260AB9D668E, 0xCE1DE40642E3F4B9}, <span class="comment">// 1e232</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	{0xC1D72B7C6B426019, 0x80D2AE83E9CE78F3}, <span class="comment">// 1e233</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	{0xB24CF65B8612F81F, 0xA1075A24E4421730}, <span class="comment">// 1e234</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	{0xDEE033F26797B627, 0xC94930AE1D529CFC}, <span class="comment">// 1e235</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	{0x169840EF017DA3B1, 0xFB9B7CD9A4A7443C}, <span class="comment">// 1e236</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	{0x8E1F289560EE864E, 0x9D412E0806E88AA5}, <span class="comment">// 1e237</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	{0xF1A6F2BAB92A27E2, 0xC491798A08A2AD4E}, <span class="comment">// 1e238</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>	{0xAE10AF696774B1DB, 0xF5B5D7EC8ACB58A2}, <span class="comment">// 1e239</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	{0xACCA6DA1E0A8EF29, 0x9991A6F3D6BF1765}, <span class="comment">// 1e240</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>	{0x17FD090A58D32AF3, 0xBFF610B0CC6EDD3F}, <span class="comment">// 1e241</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>	{0xDDFC4B4CEF07F5B0, 0xEFF394DCFF8A948E}, <span class="comment">// 1e242</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>	{0x4ABDAF101564F98E, 0x95F83D0A1FB69CD9}, <span class="comment">// 1e243</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	{0x9D6D1AD41ABE37F1, 0xBB764C4CA7A4440F}, <span class="comment">// 1e244</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	{0x84C86189216DC5ED, 0xEA53DF5FD18D5513}, <span class="comment">// 1e245</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	{0x32FD3CF5B4E49BB4, 0x92746B9BE2F8552C}, <span class="comment">// 1e246</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	{0x3FBC8C33221DC2A1, 0xB7118682DBB66A77}, <span class="comment">// 1e247</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	{0x0FABAF3FEAA5334A, 0xE4D5E82392A40515}, <span class="comment">// 1e248</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	{0x29CB4D87F2A7400E, 0x8F05B1163BA6832D}, <span class="comment">// 1e249</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	{0x743E20E9EF511012, 0xB2C71D5BCA9023F8}, <span class="comment">// 1e250</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	{0x914DA9246B255416, 0xDF78E4B2BD342CF6}, <span class="comment">// 1e251</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>	{0x1AD089B6C2F7548E, 0x8BAB8EEFB6409C1A}, <span class="comment">// 1e252</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	{0xA184AC2473B529B1, 0xAE9672ABA3D0C320}, <span class="comment">// 1e253</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	{0xC9E5D72D90A2741E, 0xDA3C0F568CC4F3E8}, <span class="comment">// 1e254</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	{0x7E2FA67C7A658892, 0x8865899617FB1871}, <span class="comment">// 1e255</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	{0xDDBB901B98FEEAB7, 0xAA7EEBFB9DF9DE8D}, <span class="comment">// 1e256</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	{0x552A74227F3EA565, 0xD51EA6FA85785631}, <span class="comment">// 1e257</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	{0xD53A88958F87275F, 0x8533285C936B35DE}, <span class="comment">// 1e258</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	{0x8A892ABAF368F137, 0xA67FF273B8460356}, <span class="comment">// 1e259</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	{0x2D2B7569B0432D85, 0xD01FEF10A657842C}, <span class="comment">// 1e260</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	{0x9C3B29620E29FC73, 0x8213F56A67F6B29B}, <span class="comment">// 1e261</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	{0x8349F3BA91B47B8F, 0xA298F2C501F45F42}, <span class="comment">// 1e262</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>	{0x241C70A936219A73, 0xCB3F2F7642717713}, <span class="comment">// 1e263</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	{0xED238CD383AA0110, 0xFE0EFB53D30DD4D7}, <span class="comment">// 1e264</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	{0xF4363804324A40AA, 0x9EC95D1463E8A506}, <span class="comment">// 1e265</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	{0xB143C6053EDCD0D5, 0xC67BB4597CE2CE48}, <span class="comment">// 1e266</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	{0xDD94B7868E94050A, 0xF81AA16FDC1B81DA}, <span class="comment">// 1e267</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	{0xCA7CF2B4191C8326, 0x9B10A4E5E9913128}, <span class="comment">// 1e268</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	{0xFD1C2F611F63A3F0, 0xC1D4CE1F63F57D72}, <span class="comment">// 1e269</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	{0xBC633B39673C8CEC, 0xF24A01A73CF2DCCF}, <span class="comment">// 1e270</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	{0xD5BE0503E085D813, 0x976E41088617CA01}, <span class="comment">// 1e271</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	{0x4B2D8644D8A74E18, 0xBD49D14AA79DBC82}, <span class="comment">// 1e272</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	{0xDDF8E7D60ED1219E, 0xEC9C459D51852BA2}, <span class="comment">// 1e273</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	{0xCABB90E5C942B503, 0x93E1AB8252F33B45}, <span class="comment">// 1e274</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	{0x3D6A751F3B936243, 0xB8DA1662E7B00A17}, <span class="comment">// 1e275</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	{0x0CC512670A783AD4, 0xE7109BFBA19C0C9D}, <span class="comment">// 1e276</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	{0x27FB2B80668B24C5, 0x906A617D450187E2}, <span class="comment">// 1e277</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	{0xB1F9F660802DEDF6, 0xB484F9DC9641E9DA}, <span class="comment">// 1e278</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	{0x5E7873F8A0396973, 0xE1A63853BBD26451}, <span class="comment">// 1e279</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	{0xDB0B487B6423E1E8, 0x8D07E33455637EB2}, <span class="comment">// 1e280</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	{0x91CE1A9A3D2CDA62, 0xB049DC016ABC5E5F}, <span class="comment">// 1e281</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>	{0x7641A140CC7810FB, 0xDC5C5301C56B75F7}, <span class="comment">// 1e282</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	{0xA9E904C87FCB0A9D, 0x89B9B3E11B6329BA}, <span class="comment">// 1e283</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	{0x546345FA9FBDCD44, 0xAC2820D9623BF429}, <span class="comment">// 1e284</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	{0xA97C177947AD4095, 0xD732290FBACAF133}, <span class="comment">// 1e285</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	{0x49ED8EABCCCC485D, 0x867F59A9D4BED6C0}, <span class="comment">// 1e286</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	{0x5C68F256BFFF5A74, 0xA81F301449EE8C70}, <span class="comment">// 1e287</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	{0x73832EEC6FFF3111, 0xD226FC195C6A2F8C}, <span class="comment">// 1e288</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	{0xC831FD53C5FF7EAB, 0x83585D8FD9C25DB7}, <span class="comment">// 1e289</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	{0xBA3E7CA8B77F5E55, 0xA42E74F3D032F525}, <span class="comment">// 1e290</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	{0x28CE1BD2E55F35EB, 0xCD3A1230C43FB26F}, <span class="comment">// 1e291</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>	{0x7980D163CF5B81B3, 0x80444B5E7AA7CF85}, <span class="comment">// 1e292</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	{0xD7E105BCC332621F, 0xA0555E361951C366}, <span class="comment">// 1e293</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	{0x8DD9472BF3FEFAA7, 0xC86AB5C39FA63440}, <span class="comment">// 1e294</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	{0xB14F98F6F0FEB951, 0xFA856334878FC150}, <span class="comment">// 1e295</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	{0x6ED1BF9A569F33D3, 0x9C935E00D4B9D8D2}, <span class="comment">// 1e296</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	{0x0A862F80EC4700C8, 0xC3B8358109E84F07}, <span class="comment">// 1e297</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	{0xCD27BB612758C0FA, 0xF4A642E14C6262C8}, <span class="comment">// 1e298</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	{0x8038D51CB897789C, 0x98E7E9CCCFBD7DBD}, <span class="comment">// 1e299</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>	{0xE0470A63E6BD56C3, 0xBF21E44003ACDD2C}, <span class="comment">// 1e300</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	{0x1858CCFCE06CAC74, 0xEEEA5D5004981478}, <span class="comment">// 1e301</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	{0x0F37801E0C43EBC8, 0x95527A5202DF0CCB}, <span class="comment">// 1e302</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	{0xD30560258F54E6BA, 0xBAA718E68396CFFD}, <span class="comment">// 1e303</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	{0x47C6B82EF32A2069, 0xE950DF20247C83FD}, <span class="comment">// 1e304</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	{0x4CDC331D57FA5441, 0x91D28B7416CDD27E}, <span class="comment">// 1e305</span>
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	{0xE0133FE4ADF8E952, 0xB6472E511C81471D}, <span class="comment">// 1e306</span>
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	{0x58180FDDD97723A6, 0xE3D8F9E563A198E5}, <span class="comment">// 1e307</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>	{0x570F09EAA7EA7648, 0x8E679C2F5E44FF8F}, <span class="comment">// 1e308</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	{0x2CD2CC6551E513DA, 0xB201833B35D63F73}, <span class="comment">// 1e309</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>	{0xF8077F7EA65E58D1, 0xDE81E40A034BCF4F}, <span class="comment">// 1e310</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	{0xFB04AFAF27FAF782, 0x8B112E86420F6191}, <span class="comment">// 1e311</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	{0x79C5DB9AF1F9B563, 0xADD57A27D29339F6}, <span class="comment">// 1e312</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	{0x18375281AE7822BC, 0xD94AD8B1C7380874}, <span class="comment">// 1e313</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	{0x8F2293910D0B15B5, 0x87CEC76F1C830548}, <span class="comment">// 1e314</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	{0xB2EB3875504DDB22, 0xA9C2794AE3A3C69A}, <span class="comment">// 1e315</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	{0x5FA60692A46151EB, 0xD433179D9C8CB841}, <span class="comment">// 1e316</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	{0xDBC7C41BA6BCD333, 0x849FEEC281D7F328}, <span class="comment">// 1e317</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	{0x12B9B522906C0800, 0xA5C7EA73224DEFF3}, <span class="comment">// 1e318</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	{0xD768226B34870A00, 0xCF39E50FEAE16BEF}, <span class="comment">// 1e319</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	{0xE6A1158300D46640, 0x81842F29F2CCE375}, <span class="comment">// 1e320</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	{0x60495AE3C1097FD0, 0xA1E53AF46F801C53}, <span class="comment">// 1e321</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	{0x385BB19CB14BDFC4, 0xCA5E89B18B602368}, <span class="comment">// 1e322</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	{0x46729E03DD9ED7B5, 0xFCF62C1DEE382C42}, <span class="comment">// 1e323</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	{0x6C07A2C26A8346D1, 0x9E19DB92B4E31BA9}, <span class="comment">// 1e324</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	{0xC7098B7305241885, 0xC5A05277621BE293}, <span class="comment">// 1e325</span>
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	{0xB8CBEE4FC66D1EA7, 0xF70867153AA2DB38}, <span class="comment">// 1e326</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	{0x737F74F1DC043328, 0x9A65406D44A5C903}, <span class="comment">// 1e327</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	{0x505F522E53053FF2, 0xC0FE908895CF3B44}, <span class="comment">// 1e328</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>	{0x647726B9E7C68FEF, 0xF13E34AABB430A15}, <span class="comment">// 1e329</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	{0x5ECA783430DC19F5, 0x96C6E0EAB509E64D}, <span class="comment">// 1e330</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	{0xB67D16413D132072, 0xBC789925624C5FE0}, <span class="comment">// 1e331</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	{0xE41C5BD18C57E88F, 0xEB96BF6EBADF77D8}, <span class="comment">// 1e332</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	{0x8E91B962F7B6F159, 0x933E37A534CBAAE7}, <span class="comment">// 1e333</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	{0x723627BBB5A4ADB0, 0xB80DC58E81FE95A1}, <span class="comment">// 1e334</span>
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	{0xCEC3B1AAA30DD91C, 0xE61136F2227E3B09}, <span class="comment">// 1e335</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	{0x213A4F0AA5E8A7B1, 0x8FCAC257558EE4E6}, <span class="comment">// 1e336</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	{0xA988E2CD4F62D19D, 0xB3BD72ED2AF29E1F}, <span class="comment">// 1e337</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	{0x93EB1B80A33B8605, 0xE0ACCFA875AF45A7}, <span class="comment">// 1e338</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	{0xBC72F130660533C3, 0x8C6C01C9498D8B88}, <span class="comment">// 1e339</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	{0xEB8FAD7C7F8680B4, 0xAF87023B9BF0EE6A}, <span class="comment">// 1e340</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	{0xA67398DB9F6820E1, 0xDB68C2CA82ED2A05}, <span class="comment">// 1e341</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	{0x88083F8943A1148C, 0x892179BE91D43A43}, <span class="comment">// 1e342</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	{0x6A0A4F6B948959B0, 0xAB69D82E364948D4}, <span class="comment">// 1e343</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	{0x848CE34679ABB01C, 0xD6444E39C3DB9B09}, <span class="comment">// 1e344</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	{0xF2D80E0C0C0B4E11, 0x85EAB0E41A6940E5}, <span class="comment">// 1e345</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	{0x6F8E118F0F0E2195, 0xA7655D1D2103911F}, <span class="comment">// 1e346</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	{0x4B7195F2D2D1A9FB, 0xD13EB46469447567}, <span class="comment">// 1e347</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>}
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>
</pre><p><a href="eisel_lemire.go?m=text">View as plain text</a></p>

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
