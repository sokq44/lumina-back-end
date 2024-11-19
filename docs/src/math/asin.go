<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/asin.go - Go Documentation Server</title>

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
<a href="asin.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<span class="text-muted">asin.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package math
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">/*
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	Floating-point arcsine and arccosine.
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	They are implemented by computing the arctangent
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	after appropriate range reduction.
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>*/</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// Asin returns the arcsine, in radians, of x.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// Special cases are:</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//	Asin(±0) = ±0</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//	Asin(x) = NaN if x &lt; -1 or x &gt; 1</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>func Asin(x float64) float64 {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	if haveArchAsin {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>		return archAsin(x)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	}
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	return asin(x)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>}
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>func asin(x float64) float64 {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	if x == 0 {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		return x <span class="comment">// special case</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	sign := false
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		x = -x
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		sign = true
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	if x &gt; 1 {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		return NaN() <span class="comment">// special case</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	temp := Sqrt(1 - x*x)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	if x &gt; 0.7 {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		temp = Pi/2 - satan(temp/x)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	} else {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		temp = satan(x / temp)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	if sign {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		temp = -temp
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	return temp
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// Acos returns the arccosine, in radians, of x.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// Special case is:</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//	Acos(x) = NaN if x &lt; -1 or x &gt; 1</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func Acos(x float64) float64 {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if haveArchAcos {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		return archAcos(x)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	return acos(x)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>func acos(x float64) float64 {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	return Pi/2 - Asin(x)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
</pre><p><a href="asin.go?m=text">View as plain text</a></p>

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
