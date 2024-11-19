<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/trig_reduce.go - Go Documentation Server</title>

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
<a href="trig_reduce.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<span class="text-muted">trig_reduce.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math">math</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2018 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package math
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;math/bits&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>)
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// reduceThreshold is the maximum value of x where the reduction using Pi/4</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// in 3 float64 parts still gives accurate results. This threshold</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// is set by y*C being representable as a float64 without error</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// where y is given by y = floor(x * (4 / Pi)) and C is the leading partial</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// terms of 4/Pi. Since the leading terms (PI4A and PI4B in sin.go) have 30</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// and 32 trailing zero bits, y should have less than 30 significant bits.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//	y &lt; 1&lt;&lt;30  -&gt; floor(x*4/Pi) &lt; 1&lt;&lt;30 -&gt; x &lt; (1&lt;&lt;30 - 1) * Pi/4</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// So, conservatively we can take x &lt; 1&lt;&lt;29.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// Above this threshold Payne-Hanek range reduction must be used.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>const reduceThreshold = 1 &lt;&lt; 29
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// trigReduce implements Payne-Hanek range reduction by Pi/4</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// for x &gt; 0. It returns the integer part mod 8 (j) and</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// the fractional part (z) of x / (Pi/4).</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// The implementation is based on:</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// &#34;ARGUMENT REDUCTION FOR HUGE ARGUMENTS: Good to the Last Bit&#34;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// K. C. Ng et al, March 24, 1992</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// The simulated multi-precision calculation of x*B uses 64-bit integer arithmetic.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>func trigReduce(x float64) (j uint64, z float64) {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	const PI4 = Pi / 4
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	if x &lt; PI4 {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		return 0, x
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// Extract out the integer and exponent such that,</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// x = ix * 2 ** exp.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	ix := Float64bits(x)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	exp := int(ix&gt;&gt;shift&amp;mask) - bias - shift
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	ix &amp;^= mask &lt;&lt; shift
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	ix |= 1 &lt;&lt; shift
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// Use the exponent to extract the 3 appropriate uint64 digits from mPi4,</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// B ~ (z0, z1, z2), such that the product leading digit has the exponent -61.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// Note, exp &gt;= -53 since x &gt;= PI4 and exp &lt; 971 for maximum float64.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	digit, bitshift := uint(exp+61)/64, uint(exp+61)%64
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	z0 := (mPi4[digit] &lt;&lt; bitshift) | (mPi4[digit+1] &gt;&gt; (64 - bitshift))
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	z1 := (mPi4[digit+1] &lt;&lt; bitshift) | (mPi4[digit+2] &gt;&gt; (64 - bitshift))
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	z2 := (mPi4[digit+2] &lt;&lt; bitshift) | (mPi4[digit+3] &gt;&gt; (64 - bitshift))
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// Multiply mantissa by the digits and extract the upper two digits (hi, lo).</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	z2hi, _ := bits.Mul64(z2, ix)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	z1hi, z1lo := bits.Mul64(z1, ix)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	z0lo := z0 * ix
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	lo, c := bits.Add64(z1lo, z2hi, 0)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	hi, _ := bits.Add64(z0lo, z1hi, c)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// The top 3 bits are j.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	j = hi &gt;&gt; 61
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// Extract the fraction and find its magnitude.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	hi = hi&lt;&lt;3 | lo&gt;&gt;61
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	lz := uint(bits.LeadingZeros64(hi))
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	e := uint64(bias - (lz + 1))
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// Clear implicit mantissa bit and shift into place.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	hi = (hi &lt;&lt; (lz + 1)) | (lo &gt;&gt; (64 - (lz + 1)))
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	hi &gt;&gt;= 64 - shift
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// Include the exponent and convert to a float.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	hi |= e &lt;&lt; shift
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	z = Float64frombits(hi)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// Map zeros to origin.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	if j&amp;1 == 1 {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		j++
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		j &amp;= 7
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		z--
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// Multiply the fractional part by pi/4.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return j, z * PI4
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// mPi4 is the binary digits of 4/pi as a uint64 array,</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// that is, 4/pi = Sum mPi4[i]*2^(-64*i)</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// 19 64-bit digits and the leading one bit give 1217 bits</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// of precision to handle the largest possible float64 exponent.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>var mPi4 = [...]uint64{
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	0x0000000000000001,
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	0x45f306dc9c882a53,
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	0xf84eafa3ea69bb81,
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	0xb6c52b3278872083,
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	0xfca2c757bd778ac3,
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	0x6e48dc74849ba5c0,
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	0x0c925dd413a32439,
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	0xfc3bd63962534e7d,
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	0xd1046bea5d768909,
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	0xd338e04d68befc82,
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	0x7323ac7306a673e9,
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	0x3908bf177bf25076,
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	0x3ff12fffbc0b301f,
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	0xde5e2316b414da3e,
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	0xda6cfd9e4f96136e,
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	0x9e8c7ecd3cbfd45a,
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	0xea4f758fd7cbe2f6,
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	0x7a0e73ef14a525d4,
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	0xd7f6bf623f1aba10,
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	0xac06608df8f6d757,
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
</pre><p><a href="trig_reduce.go?m=text">View as plain text</a></p>

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
