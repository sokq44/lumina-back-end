<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/const.go - Go Documentation Server</title>

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
<a href="const.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<span class="text-muted">const.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math">math</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package math provides basic constants and mathematical functions.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This package does not guarantee bit-identical results across architectures.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>package math
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// Mathematical constants.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>const (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	E   = 2.71828182845904523536028747135266249775724709369995957496696763 <span class="comment">// https://oeis.org/A001113</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	Pi  = 3.14159265358979323846264338327950288419716939937510582097494459 <span class="comment">// https://oeis.org/A000796</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	Phi = 1.61803398874989484820458683436563811772030917980576286213544862 <span class="comment">// https://oeis.org/A001622</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	Sqrt2   = 1.41421356237309504880168872420969807856967187537694807317667974 <span class="comment">// https://oeis.org/A002193</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	SqrtE   = 1.64872127070012814684865078781416357165377610071014801157507931 <span class="comment">// https://oeis.org/A019774</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	SqrtPi  = 1.77245385090551602729816748334114518279754945612238712821380779 <span class="comment">// https://oeis.org/A002161</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	SqrtPhi = 1.27201964951406896425242246173749149171560804184009624861664038 <span class="comment">// https://oeis.org/A139339</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	Ln2    = 0.693147180559945309417232121458176568075500134360255254120680009 <span class="comment">// https://oeis.org/A002162</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	Log2E  = 1 / Ln2
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	Ln10   = 2.30258509299404568401799145468436420760110148862877297603332790 <span class="comment">// https://oeis.org/A002392</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	Log10E = 1 / Ln10
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// Floating-point limit values.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// Max is the largest finite value representable by the type.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// SmallestNonzero is the smallest positive, non-zero value representable by the type.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>const (
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	MaxFloat32             = 0x1p127 * (1 + (1 - 0x1p-23)) <span class="comment">// 3.40282346638528859811704183484516925440e+38</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	SmallestNonzeroFloat32 = 0x1p-126 * 0x1p-23            <span class="comment">// 1.401298464324817070923729583289916131280e-45</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	MaxFloat64             = 0x1p1023 * (1 + (1 - 0x1p-52)) <span class="comment">// 1.79769313486231570814527423731704356798070e+308</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	SmallestNonzeroFloat64 = 0x1p-1022 * 0x1p-52            <span class="comment">// 4.9406564584124654417656879286822137236505980e-324</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// Integer limit values.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>const (
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	intSize = 32 &lt;&lt; (^uint(0) &gt;&gt; 63) <span class="comment">// 32 or 64</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	MaxInt    = 1&lt;&lt;(intSize-1) - 1  <span class="comment">// MaxInt32 or MaxInt64 depending on intSize.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	MinInt    = -1 &lt;&lt; (intSize - 1) <span class="comment">// MinInt32 or MinInt64 depending on intSize.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	MaxInt8   = 1&lt;&lt;7 - 1            <span class="comment">// 127</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	MinInt8   = -1 &lt;&lt; 7             <span class="comment">// -128</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	MaxInt16  = 1&lt;&lt;15 - 1           <span class="comment">// 32767</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	MinInt16  = -1 &lt;&lt; 15            <span class="comment">// -32768</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	MaxInt32  = 1&lt;&lt;31 - 1           <span class="comment">// 2147483647</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	MinInt32  = -1 &lt;&lt; 31            <span class="comment">// -2147483648</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	MaxInt64  = 1&lt;&lt;63 - 1           <span class="comment">// 9223372036854775807</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	MinInt64  = -1 &lt;&lt; 63            <span class="comment">// -9223372036854775808</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	MaxUint   = 1&lt;&lt;intSize - 1      <span class="comment">// MaxUint32 or MaxUint64 depending on intSize.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	MaxUint8  = 1&lt;&lt;8 - 1            <span class="comment">// 255</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	MaxUint16 = 1&lt;&lt;16 - 1           <span class="comment">// 65535</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	MaxUint32 = 1&lt;&lt;32 - 1           <span class="comment">// 4294967295</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	MaxUint64 = 1&lt;&lt;64 - 1           <span class="comment">// 18446744073709551615</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
</pre><p><a href="const.go?m=text">View as plain text</a></p>

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
