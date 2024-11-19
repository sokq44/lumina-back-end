<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/pow.go - Go Documentation Server</title>

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
<a href="pow.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<span class="text-muted">pow.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>func isOddInt(x float64) bool {
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	if Abs(x) &gt;= (1 &lt;&lt; 53) {
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>		<span class="comment">// 1 &lt;&lt; 53 is the largest exact integer in the float64 format.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>		<span class="comment">// Any number outside this range will be truncated before the decimal point and therefore will always be</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>		<span class="comment">// an even integer.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>		<span class="comment">// Without this check and if x overflows int64 the int64(xi) conversion below may produce incorrect results</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>		<span class="comment">// on some architectures (and does so on arm64). See issue #57465.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>		return false
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	}
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	xi, xf := Modf(x)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	return xf == 0 &amp;&amp; int64(xi)&amp;1 == 1
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>}
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// Special cases taken from FreeBSD&#39;s /usr/src/lib/msun/src/e_pow.c</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// updated by IEEE Std. 754-2008 &#34;Section 9.2.1 Special values&#34;.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// Pow returns x**y, the base-x exponential of y.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// Special cases are (in order):</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//	Pow(x, ±0) = 1 for any x</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//	Pow(1, y) = 1 for any y</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//	Pow(x, 1) = x for any x</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//	Pow(NaN, y) = NaN</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//	Pow(x, NaN) = NaN</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//	Pow(±0, y) = ±Inf for y an odd integer &lt; 0</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//	Pow(±0, -Inf) = +Inf</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//	Pow(±0, +Inf) = +0</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//	Pow(±0, y) = +Inf for finite y &lt; 0 and not an odd integer</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//	Pow(±0, y) = ±0 for y an odd integer &gt; 0</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//	Pow(±0, y) = +0 for finite y &gt; 0 and not an odd integer</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//	Pow(-1, ±Inf) = 1</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//	Pow(x, +Inf) = +Inf for |x| &gt; 1</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//	Pow(x, -Inf) = +0 for |x| &gt; 1</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//	Pow(x, +Inf) = +0 for |x| &lt; 1</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//	Pow(x, -Inf) = +Inf for |x| &lt; 1</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//	Pow(+Inf, y) = +Inf for y &gt; 0</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//	Pow(+Inf, y) = +0 for y &lt; 0</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//	Pow(-Inf, y) = Pow(-0, -y)</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//	Pow(x, y) = NaN for finite x &lt; 0 and finite non-integer y</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>func Pow(x, y float64) float64 {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	if haveArchPow {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		return archPow(x, y)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	return pow(x, y)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>func pow(x, y float64) float64 {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	switch {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	case y == 0 || x == 1:
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		return 1
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	case y == 1:
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		return x
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	case IsNaN(x) || IsNaN(y):
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		return NaN()
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	case x == 0:
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		switch {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		case y &lt; 0:
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			if Signbit(x) &amp;&amp; isOddInt(y) {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>				return Inf(-1)
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			return Inf(1)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		case y &gt; 0:
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>			if Signbit(x) &amp;&amp; isOddInt(y) {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>				return x
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			return 0
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	case IsInf(y, 0):
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		switch {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		case x == -1:
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			return 1
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		case (Abs(x) &lt; 1) == IsInf(y, 1):
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			return 0
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		default:
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			return Inf(1)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	case IsInf(x, 0):
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		if IsInf(x, -1) {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			return Pow(1/x, -y) <span class="comment">// Pow(-0, -y)</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		switch {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		case y &lt; 0:
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			return 0
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		case y &gt; 0:
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			return Inf(1)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	case y == 0.5:
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		return Sqrt(x)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	case y == -0.5:
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		return 1 / Sqrt(x)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	yi, yf := Modf(Abs(y))
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	if yf != 0 &amp;&amp; x &lt; 0 {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		return NaN()
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	if yi &gt;= 1&lt;&lt;63 {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">// yi is a large even int that will lead to overflow (or underflow to 0)</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		<span class="comment">// for all x except -1 (x == 1 was handled earlier)</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		switch {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		case x == -1:
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			return 1
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		case (Abs(x) &lt; 1) == (y &gt; 0):
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			return 0
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		default:
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			return Inf(1)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// ans = a1 * 2**ae (= 1 for now).</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	a1 := 1.0
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	ae := 0
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// ans *= x**yf</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if yf != 0 {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		if yf &gt; 0.5 {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			yf--
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			yi++
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		a1 = Exp(yf * Log(x))
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// ans *= x**yi</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// by multiplying in successive squarings</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// of x according to bits of yi.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// accumulate powers of two into exp.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	x1, xe := Frexp(x)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	for i := int64(yi); i != 0; i &gt;&gt;= 1 {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		if xe &lt; -1&lt;&lt;12 || 1&lt;&lt;12 &lt; xe {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			<span class="comment">// catch xe before it overflows the left shift below</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			<span class="comment">// Since i !=0 it has at least one bit still set, so ae will accumulate xe</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			<span class="comment">// on at least one more iteration, ae += xe is a lower bound on ae</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			<span class="comment">// the lower bound on ae exceeds the size of a float64 exp</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			<span class="comment">// so the final call to Ldexp will produce under/overflow (0/Inf)</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			ae += xe
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			break
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		if i&amp;1 == 1 {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			a1 *= x1
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			ae += xe
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		x1 *= x1
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		xe &lt;&lt;= 1
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		if x1 &lt; .5 {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			x1 += x1
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			xe--
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// ans = a1*2**ae</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// if y &lt; 0 { ans = 1 / ans }</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// but in the opposite order</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if y &lt; 0 {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		a1 = 1 / a1
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		ae = -ae
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	return Ldexp(a1, ae)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
</pre><p><a href="pow.go?m=text">View as plain text</a></p>

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
