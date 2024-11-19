<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/image/color/ycbcr.go - Go Documentation Server</title>

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
<a href="ycbcr.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/image">image</a>/<a href="http://localhost:8080/src/image/color">color</a>/<span class="text-muted">ycbcr.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/image/color">image/color</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package color
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// RGBToYCbCr converts an RGB triple to a Y&#39;CbCr triple.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>func RGBToYCbCr(r, g, b uint8) (uint8, uint8, uint8) {
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	<span class="comment">// The JFIF specification says:</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	<span class="comment">//	Y&#39; =  0.2990*R + 0.5870*G + 0.1140*B</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	<span class="comment">//	Cb = -0.1687*R - 0.3313*G + 0.5000*B + 128</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	<span class="comment">//	Cr =  0.5000*R - 0.4187*G - 0.0813*B + 128</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	<span class="comment">// https://www.w3.org/Graphics/JPEG/jfif3.pdf says Y but means Y&#39;.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	r1 := int32(r)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	g1 := int32(g)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	b1 := int32(b)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// yy is in range [0,0xff].</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// Note that 19595 + 38470 + 7471 equals 65536.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	yy := (19595*r1 + 38470*g1 + 7471*b1 + 1&lt;&lt;15) &gt;&gt; 16
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// The bit twiddling below is equivalent to</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// cb := (-11056*r1 - 21712*g1 + 32768*b1 + 257&lt;&lt;15) &gt;&gt; 16</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// if cb &lt; 0 {</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">//     cb = 0</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// } else if cb &gt; 0xff {</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">//     cb = ^int32(0)</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// }</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// but uses fewer branches and is faster.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// Note that the uint8 type conversion in the return</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// statement will convert ^int32(0) to 0xff.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// The code below to compute cr uses a similar pattern.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// Note that -11056 - 21712 + 32768 equals 0.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	cb := -11056*r1 - 21712*g1 + 32768*b1 + 257&lt;&lt;15
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	if uint32(cb)&amp;0xff000000 == 0 {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		cb &gt;&gt;= 16
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	} else {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		cb = ^(cb &gt;&gt; 31)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// Note that 32768 - 27440 - 5328 equals 0.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	cr := 32768*r1 - 27440*g1 - 5328*b1 + 257&lt;&lt;15
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	if uint32(cr)&amp;0xff000000 == 0 {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		cr &gt;&gt;= 16
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	} else {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		cr = ^(cr &gt;&gt; 31)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	return uint8(yy), uint8(cb), uint8(cr)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// YCbCrToRGB converts a Y&#39;CbCr triple to an RGB triple.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func YCbCrToRGB(y, cb, cr uint8) (uint8, uint8, uint8) {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// The JFIF specification says:</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">//	R = Y&#39; + 1.40200*(Cr-128)</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">//	G = Y&#39; - 0.34414*(Cb-128) - 0.71414*(Cr-128)</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">//	B = Y&#39; + 1.77200*(Cb-128)</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// https://www.w3.org/Graphics/JPEG/jfif3.pdf says Y but means Y&#39;.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// Those formulae use non-integer multiplication factors. When computing,</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// integer math is generally faster than floating point math. We multiply</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// all of those factors by 1&lt;&lt;16 and round to the nearest integer:</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">//	 91881 = roundToNearestInteger(1.40200 * 65536).</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">//	 22554 = roundToNearestInteger(0.34414 * 65536).</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">//	 46802 = roundToNearestInteger(0.71414 * 65536).</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">//	116130 = roundToNearestInteger(1.77200 * 65536).</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// Adding a rounding adjustment in the range [0, 1&lt;&lt;16-1] and then shifting</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// right by 16 gives us an integer math version of the original formulae.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">//	R = (65536*Y&#39; +  91881 *(Cr-128)                  + adjustment) &gt;&gt; 16</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">//	G = (65536*Y&#39; -  22554 *(Cb-128) - 46802*(Cr-128) + adjustment) &gt;&gt; 16</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">//	B = (65536*Y&#39; + 116130 *(Cb-128)                  + adjustment) &gt;&gt; 16</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// A constant rounding adjustment of 1&lt;&lt;15, one half of 1&lt;&lt;16, would mean</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// round-to-nearest when dividing by 65536 (shifting right by 16).</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// Similarly, a constant rounding adjustment of 0 would mean round-down.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// Defining YY1 = 65536*Y&#39; + adjustment simplifies the formulae and</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// requires fewer CPU operations:</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">//	R = (YY1 +  91881 *(Cr-128)                 ) &gt;&gt; 16</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">//	G = (YY1 -  22554 *(Cb-128) - 46802*(Cr-128)) &gt;&gt; 16</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">//	B = (YY1 + 116130 *(Cb-128)                 ) &gt;&gt; 16</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// The inputs (y, cb, cr) are 8 bit color, ranging in [0x00, 0xff]. In this</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// function, the output is also 8 bit color, but in the related YCbCr.RGBA</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// method, below, the output is 16 bit color, ranging in [0x0000, 0xffff].</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// Outputting 16 bit color simply requires changing the 16 to 8 in the &#34;R =</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// etc &gt;&gt; 16&#34; equation, and likewise for G and B.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// As mentioned above, a constant rounding adjustment of 1&lt;&lt;15 is a natural</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// choice, but there is an additional constraint: if c0 := YCbCr{Y: y, Cb:</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// 0x80, Cr: 0x80} and c1 := Gray{Y: y} then c0.RGBA() should equal</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// c1.RGBA(). Specifically, if y == 0 then &#34;R = etc &gt;&gt; 8&#34; should yield</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// 0x0000 and if y == 0xff then &#34;R = etc &gt;&gt; 8&#34; should yield 0xffff. If we</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// used a constant rounding adjustment of 1&lt;&lt;15, then it would yield 0x0080</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// and 0xff80 respectively.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// Note that when cb == 0x80 and cr == 0x80 then the formulae collapse to:</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">//	R = YY1 &gt;&gt; n</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">//	G = YY1 &gt;&gt; n</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">//	B = YY1 &gt;&gt; n</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	<span class="comment">// where n is 16 for this function (8 bit color output) and 8 for the</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// YCbCr.RGBA method (16 bit color output).</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// The solution is to make the rounding adjustment non-constant, and equal</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// to 257*Y&#39;, which ranges over [0, 1&lt;&lt;16-1] as Y&#39; ranges over [0, 255].</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// YY1 is then defined as:</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">//	YY1 = 65536*Y&#39; + 257*Y&#39;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// or equivalently:</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">//	YY1 = Y&#39; * 0x10101</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	yy1 := int32(y) * 0x10101
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	cb1 := int32(cb) - 128
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	cr1 := int32(cr) - 128
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// The bit twiddling below is equivalent to</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// r := (yy1 + 91881*cr1) &gt;&gt; 16</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// if r &lt; 0 {</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">//     r = 0</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// } else if r &gt; 0xff {</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">//     r = ^int32(0)</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// }</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// but uses fewer branches and is faster.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// Note that the uint8 type conversion in the return</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// statement will convert ^int32(0) to 0xff.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// The code below to compute g and b uses a similar pattern.</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	r := yy1 + 91881*cr1
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if uint32(r)&amp;0xff000000 == 0 {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		r &gt;&gt;= 16
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	} else {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		r = ^(r &gt;&gt; 31)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	g := yy1 - 22554*cb1 - 46802*cr1
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	if uint32(g)&amp;0xff000000 == 0 {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		g &gt;&gt;= 16
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	} else {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		g = ^(g &gt;&gt; 31)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	b := yy1 + 116130*cb1
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	if uint32(b)&amp;0xff000000 == 0 {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		b &gt;&gt;= 16
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	} else {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		b = ^(b &gt;&gt; 31)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	return uint8(r), uint8(g), uint8(b)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// YCbCr represents a fully opaque 24-bit Y&#39;CbCr color, having 8 bits each for</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// one luma and two chroma components.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// JPEG, VP8, the MPEG family and other codecs use this color model. Such</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// codecs often use the terms YUV and Y&#39;CbCr interchangeably, but strictly</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">// speaking, the term YUV applies only to analog video signals, and Y&#39; (luma)</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// is Y (luminance) after applying gamma correction.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// Conversion between RGB and Y&#39;CbCr is lossy and there are multiple, slightly</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">// different formulae for converting between the two. This package follows</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// the JFIF specification at https://www.w3.org/Graphics/JPEG/jfif3.pdf.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>type YCbCr struct {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	Y, Cb, Cr uint8
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>func (c YCbCr) RGBA() (uint32, uint32, uint32, uint32) {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// This code is a copy of the YCbCrToRGB function above, except that it</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// returns values in the range [0, 0xffff] instead of [0, 0xff]. There is a</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// subtle difference between doing this and having YCbCr satisfy the Color</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// interface by first converting to an RGBA. The latter loses some</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// information by going to and from 8 bits per channel.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// For example, this code:</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">//	const y, cb, cr = 0x7f, 0x7f, 0x7f</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">//	r, g, b := color.YCbCrToRGB(y, cb, cr)</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">//	r0, g0, b0, _ := color.YCbCr{y, cb, cr}.RGBA()</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">//	r1, g1, b1, _ := color.RGBA{r, g, b, 0xff}.RGBA()</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">//	fmt.Printf(&#34;0x%04x 0x%04x 0x%04x\n&#34;, r0, g0, b0)</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">//	fmt.Printf(&#34;0x%04x 0x%04x 0x%04x\n&#34;, r1, g1, b1)</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// prints:</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">//	0x7e18 0x808d 0x7db9</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">//	0x7e7e 0x8080 0x7d7d</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	yy1 := int32(c.Y) * 0x10101
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	cb1 := int32(c.Cb) - 128
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	cr1 := int32(c.Cr) - 128
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// The bit twiddling below is equivalent to</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// r := (yy1 + 91881*cr1) &gt;&gt; 8</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// if r &lt; 0 {</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">//     r = 0</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// } else if r &gt; 0xff {</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">//     r = 0xffff</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// }</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// but uses fewer branches and is faster.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	<span class="comment">// The code below to compute g and b uses a similar pattern.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	r := yy1 + 91881*cr1
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	if uint32(r)&amp;0xff000000 == 0 {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		r &gt;&gt;= 8
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	} else {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		r = ^(r &gt;&gt; 31) &amp; 0xffff
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	g := yy1 - 22554*cb1 - 46802*cr1
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	if uint32(g)&amp;0xff000000 == 0 {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		g &gt;&gt;= 8
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	} else {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		g = ^(g &gt;&gt; 31) &amp; 0xffff
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	b := yy1 + 116130*cb1
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	if uint32(b)&amp;0xff000000 == 0 {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		b &gt;&gt;= 8
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	} else {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		b = ^(b &gt;&gt; 31) &amp; 0xffff
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	return uint32(r), uint32(g), uint32(b), 0xffff
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// YCbCrModel is the [Model] for Y&#39;CbCr colors.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>var YCbCrModel Model = ModelFunc(yCbCrModel)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>func yCbCrModel(c Color) Color {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	if _, ok := c.(YCbCr); ok {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		return c
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	r, g, b, _ := c.RGBA()
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	y, u, v := RGBToYCbCr(uint8(r&gt;&gt;8), uint8(g&gt;&gt;8), uint8(b&gt;&gt;8))
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	return YCbCr{y, u, v}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span><span class="comment">// NYCbCrA represents a non-alpha-premultiplied Y&#39;CbCr-with-alpha color, having</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span><span class="comment">// 8 bits each for one luma, two chroma and one alpha component.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>type NYCbCrA struct {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	YCbCr
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	A uint8
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>func (c NYCbCrA) RGBA() (uint32, uint32, uint32, uint32) {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// The first part of this method is the same as YCbCr.RGBA.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	yy1 := int32(c.Y) * 0x10101
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	cb1 := int32(c.Cb) - 128
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	cr1 := int32(c.Cr) - 128
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// The bit twiddling below is equivalent to</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	<span class="comment">// r := (yy1 + 91881*cr1) &gt;&gt; 8</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// if r &lt; 0 {</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">//     r = 0</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	<span class="comment">// } else if r &gt; 0xff {</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">//     r = 0xffff</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// }</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// but uses fewer branches and is faster.</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// The code below to compute g and b uses a similar pattern.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	r := yy1 + 91881*cr1
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	if uint32(r)&amp;0xff000000 == 0 {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		r &gt;&gt;= 8
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	} else {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		r = ^(r &gt;&gt; 31) &amp; 0xffff
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	g := yy1 - 22554*cb1 - 46802*cr1
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	if uint32(g)&amp;0xff000000 == 0 {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		g &gt;&gt;= 8
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	} else {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		g = ^(g &gt;&gt; 31) &amp; 0xffff
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	b := yy1 + 116130*cb1
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	if uint32(b)&amp;0xff000000 == 0 {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		b &gt;&gt;= 8
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	} else {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		b = ^(b &gt;&gt; 31) &amp; 0xffff
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">// The second part of this method applies the alpha.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	a := uint32(c.A) * 0x101
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	return uint32(r) * a / 0xffff, uint32(g) * a / 0xffff, uint32(b) * a / 0xffff, a
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">// NYCbCrAModel is the [Model] for non-alpha-premultiplied Y&#39;CbCr-with-alpha</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span><span class="comment">// colors.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>var NYCbCrAModel Model = ModelFunc(nYCbCrAModel)
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>func nYCbCrAModel(c Color) Color {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	switch c := c.(type) {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	case NYCbCrA:
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		return c
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	case YCbCr:
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		return NYCbCrA{c, 0xff}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	r, g, b, a := c.RGBA()
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	<span class="comment">// Convert from alpha-premultiplied to non-alpha-premultiplied.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	if a != 0 {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		r = (r * 0xffff) / a
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		g = (g * 0xffff) / a
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		b = (b * 0xffff) / a
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	y, u, v := RGBToYCbCr(uint8(r&gt;&gt;8), uint8(g&gt;&gt;8), uint8(b&gt;&gt;8))
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	return NYCbCrA{YCbCr{Y: y, Cb: u, Cr: v}, uint8(a &gt;&gt; 8)}
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">// RGBToCMYK converts an RGB triple to a CMYK quadruple.</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>func RGBToCMYK(r, g, b uint8) (uint8, uint8, uint8, uint8) {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	rr := uint32(r)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	gg := uint32(g)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	bb := uint32(b)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	w := rr
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	if w &lt; gg {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		w = gg
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	if w &lt; bb {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		w = bb
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	if w == 0 {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		return 0, 0, 0, 0xff
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	c := (w - rr) * 0xff / w
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	m := (w - gg) * 0xff / w
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	y := (w - bb) * 0xff / w
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	return uint8(c), uint8(m), uint8(y), uint8(0xff - w)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// CMYKToRGB converts a [CMYK] quadruple to an RGB triple.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func CMYKToRGB(c, m, y, k uint8) (uint8, uint8, uint8) {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	w := 0xffff - uint32(k)*0x101
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	r := (0xffff - uint32(c)*0x101) * w / 0xffff
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	g := (0xffff - uint32(m)*0x101) * w / 0xffff
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	b := (0xffff - uint32(y)*0x101) * w / 0xffff
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	return uint8(r &gt;&gt; 8), uint8(g &gt;&gt; 8), uint8(b &gt;&gt; 8)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">// CMYK represents a fully opaque CMYK color, having 8 bits for each of cyan,</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// magenta, yellow and black.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// It is not associated with any particular color profile.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>type CMYK struct {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	C, M, Y, K uint8
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>func (c CMYK) RGBA() (uint32, uint32, uint32, uint32) {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	<span class="comment">// This code is a copy of the CMYKToRGB function above, except that it</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	<span class="comment">// returns values in the range [0, 0xffff] instead of [0, 0xff].</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	w := 0xffff - uint32(c.K)*0x101
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	r := (0xffff - uint32(c.C)*0x101) * w / 0xffff
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	g := (0xffff - uint32(c.M)*0x101) * w / 0xffff
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	b := (0xffff - uint32(c.Y)*0x101) * w / 0xffff
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	return r, g, b, 0xffff
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span><span class="comment">// CMYKModel is the [Model] for CMYK colors.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>var CMYKModel Model = ModelFunc(cmykModel)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>func cmykModel(c Color) Color {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	if _, ok := c.(CMYK); ok {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		return c
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	}
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	r, g, b, _ := c.RGBA()
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	cc, mm, yy, kk := RGBToCMYK(uint8(r&gt;&gt;8), uint8(g&gt;&gt;8), uint8(b&gt;&gt;8))
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	return CMYK{cc, mm, yy, kk}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
</pre><p><a href="ycbcr.go?m=text">View as plain text</a></p>

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
