<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/sqrt.go - Go Documentation Server</title>

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
<a href="sqrt.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<span class="text-muted">sqrt.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// The original C code and the long comment below are</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// from FreeBSD&#39;s /usr/src/lib/msun/src/e_sqrt.c and</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// came with this notice. The go code is a simplified</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// version of the original C.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// ====================================================</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// Developed at SunPro, a Sun Microsystems, Inc. business.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// Permission to use, copy, modify, and distribute this</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// software is freely granted, provided that this notice</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// is preserved.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// ====================================================</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// __ieee754_sqrt(x)</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// Return correctly rounded sqrt.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//           -----------------------------------------</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//           | Use the hardware sqrt if you have one |</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//           -----------------------------------------</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// Method:</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//   Bit by bit method using integer arithmetic. (Slow, but portable)</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//   1. Normalization</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//      Scale x to y in [1,4) with even powers of 2:</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//      find an integer k such that  1 &lt;= (y=x*2**(2k)) &lt; 4, then</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//              sqrt(x) = 2**k * sqrt(y)</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//   2. Bit by bit computation</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//      Let q  = sqrt(y) truncated to i bit after binary point (q = 1),</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//           i                                                   0</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//                                     i+1         2</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//          s  = 2*q , and      y  =  2   * ( y - q  ).          (1)</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//           i      i            i                 i</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//      To compute q    from q , one checks whether</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//                  i+1       i</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//                            -(i+1) 2</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//                      (q + 2      )  &lt;= y.                     (2)</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//                        i</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//                                                            -(i+1)</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//      If (2) is false, then q   = q ; otherwise q   = q  + 2      .</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//                             i+1   i             i+1   i</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//      With some algebraic manipulation, it is not difficult to see</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//      that (2) is equivalent to</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//                             -(i+1)</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//                      s  +  2       &lt;= y                       (3)</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//                       i                i</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//      The advantage of (3) is that s  and y  can be computed by</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//                                    i      i</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//      the following recurrence formula:</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//          if (3) is false</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//          s     =  s  ,       y    = y   ;                     (4)</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//           i+1      i          i+1    i</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//      otherwise,</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//                         -i                      -(i+1)</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//          s     =  s  + 2  ,  y    = y  -  s  - 2              (5)</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">//           i+1      i          i+1    i     i</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//      One may easily use induction to prove (4) and (5).</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//      Note. Since the left hand side of (3) contain only i+2 bits,</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">//            it is not necessary to do a full (53-bit) comparison</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//            in (3).</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//   3. Final rounding</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//      After generating the 53 bits result, we compute one more bit.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//      Together with the remainder, we can decide whether the</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//      result is exact, bigger than 1/2ulp, or less than 1/2ulp</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//      (it will never equal to 1/2ulp).</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//      The rounding mode can be detected by checking whether</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//      huge + tiny is equal to huge, and whether huge - tiny is</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//      equal to huge for some floating point number &#34;huge&#34; and &#34;tiny&#34;.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// Notes:  Rounding mode detection omitted. The constants &#34;mask&#34;, &#34;shift&#34;,</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// and &#34;bias&#34; are found in src/math/bits.go</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// Sqrt returns the square root of x.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// Special cases are:</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//	Sqrt(+Inf) = +Inf</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">//	Sqrt(±0) = ±0</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//	Sqrt(x &lt; 0) = NaN</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">//	Sqrt(NaN) = NaN</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>func Sqrt(x float64) float64 {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	return sqrt(x)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// Note: On systems where Sqrt is a single instruction, the compiler</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// may turn a direct call into a direct use of that instruction instead.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>func sqrt(x float64) float64 {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// special cases</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	switch {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	case x == 0 || IsNaN(x) || IsInf(x, 1):
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		return x
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	case x &lt; 0:
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		return NaN()
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	ix := Float64bits(x)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// normalize x</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	exp := int((ix &gt;&gt; shift) &amp; mask)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	if exp == 0 { <span class="comment">// subnormal x</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		for ix&amp;(1&lt;&lt;shift) == 0 {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			ix &lt;&lt;= 1
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			exp--
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		exp++
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	exp -= bias <span class="comment">// unbias exponent</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	ix &amp;^= mask &lt;&lt; shift
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	ix |= 1 &lt;&lt; shift
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	if exp&amp;1 == 1 { <span class="comment">// odd exp, double x to make it even</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		ix &lt;&lt;= 1
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	exp &gt;&gt;= 1 <span class="comment">// exp = exp/2, exponent of square root</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// generate sqrt(x) bit by bit</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	ix &lt;&lt;= 1
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	var q, s uint64               <span class="comment">// q = sqrt(x)</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	r := uint64(1 &lt;&lt; (shift + 1)) <span class="comment">// r = moving bit from MSB to LSB</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	for r != 0 {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		t := s + r
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		if t &lt;= ix {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			s = t + r
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			ix -= t
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			q += r
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		ix &lt;&lt;= 1
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		r &gt;&gt;= 1
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// final rounding</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	if ix != 0 { <span class="comment">// remainder, result not exact</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		q += q &amp; 1 <span class="comment">// round according to extra bit</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	ix = q&gt;&gt;1 + uint64(exp-1+bias)&lt;&lt;shift <span class="comment">// significand + biased exponent</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	return Float64frombits(ix)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
</pre><p><a href="sqrt.go?m=text">View as plain text</a></p>

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
