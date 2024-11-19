<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/strconv/ctoa.go - Go Documentation Server</title>

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
<a href="ctoa.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/strconv">strconv</a>/<span class="text-muted">ctoa.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// FormatComplex converts the complex number c to a string of the</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// form (a+bi) where a and b are the real and imaginary parts,</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// formatted according to the format fmt and precision prec.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// The format fmt and precision prec have the same meaning as in FormatFloat.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// It rounds the result assuming that the original was obtained from a complex</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// value of bitSize bits, which must be 64 for complex64 and 128 for complex128.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>func FormatComplex(c complex128, fmt byte, prec, bitSize int) string {
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	if bitSize != 64 &amp;&amp; bitSize != 128 {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>		panic(&#34;invalid bitSize&#34;)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	}
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	bitSize &gt;&gt;= 1 <span class="comment">// complex64 uses float32 internally</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// Check if imaginary part has a sign. If not, add one.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	im := FormatFloat(imag(c), fmt, prec, bitSize)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	if im[0] != &#39;+&#39; &amp;&amp; im[0] != &#39;-&#39; {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>		im = &#34;+&#34; + im
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	return &#34;(&#34; + FormatFloat(real(c), fmt, prec, bitSize) + im + &#34;i)&#34;
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
</pre><p><a href="ctoa.go?m=text">View as plain text</a></p>

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
