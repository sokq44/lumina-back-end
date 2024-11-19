<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/nextafter.go - Go Documentation Server</title>

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
<a href="nextafter.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<span class="text-muted">nextafter.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math">math</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2010 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package math
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Nextafter32 returns the next representable float32 value after x towards y.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// Special cases are:</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//	Nextafter32(x, x)   = x</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//	Nextafter32(NaN, y) = NaN</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//	Nextafter32(x, NaN) = NaN</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>func Nextafter32(x, y float32) (r float32) {
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	switch {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	case IsNaN(float64(x)) || IsNaN(float64(y)): <span class="comment">// special case</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>		r = float32(NaN())
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	case x == y:
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>		r = x
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	case x == 0:
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>		r = float32(Copysign(float64(Float32frombits(1)), float64(y)))
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	case (y &gt; x) == (x &gt; 0):
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>		r = Float32frombits(Float32bits(x) + 1)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	default:
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		r = Float32frombits(Float32bits(x) - 1)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	}
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	return
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// Nextafter returns the next representable float64 value after x towards y.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// Special cases are:</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//	Nextafter(x, x)   = x</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//	Nextafter(NaN, y) = NaN</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//	Nextafter(x, NaN) = NaN</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>func Nextafter(x, y float64) (r float64) {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	switch {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	case IsNaN(x) || IsNaN(y): <span class="comment">// special case</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		r = NaN()
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	case x == y:
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		r = x
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	case x == 0:
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		r = Copysign(Float64frombits(1), y)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	case (y &gt; x) == (x &gt; 0):
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		r = Float64frombits(Float64bits(x) + 1)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	default:
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		r = Float64frombits(Float64bits(x) - 1)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	return
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
</pre><p><a href="nextafter.go?m=text">View as plain text</a></p>

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
