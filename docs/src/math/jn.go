<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/jn.go - Go Documentation Server</title>

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
<a href="jn.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<span class="text-muted">jn.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">/*
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	Bessel function of the first and second kinds of order n.
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>*/</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// The original C code and the long comment below are</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// from FreeBSD&#39;s /usr/src/lib/msun/src/e_jn.c and</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// came with this notice. The go code is a simplified</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// version of the original C.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// ====================================================</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// Developed at SunPro, a Sun Microsystems, Inc. business.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// Permission to use, copy, modify, and distribute this</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// software is freely granted, provided that this notice</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// is preserved.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// ====================================================</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// __ieee754_jn(n, x), __ieee754_yn(n, x)</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// floating point Bessel&#39;s function of the 1st and 2nd kind</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// of order n</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// Special cases:</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//      y0(0)=y1(0)=yn(n,0) = -inf with division by zero signal;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//      y0(-ve)=y1(-ve)=yn(n,-ve) are NaN with invalid signal.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// Note 2. About jn(n,x), yn(n,x)</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//      For n=0, j0(x) is called,</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//      for n=1, j1(x) is called,</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//      for n&lt;x, forward recursion is used starting</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//      from values of j0(x) and j1(x).</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//      for n&gt;x, a continued fraction approximation to</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//      j(n,x)/j(n-1,x) is evaluated and then backward</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//      recursion is used starting from a supposed value</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//      for j(n,x). The resulting value of j(0,x) is</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//      compared with the actual value to correct the</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//      supposed value of j(n,x).</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//      yn(n,x) is similar in all respects, except</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//      that forward recursion is used for all</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//      values of n&gt;1.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// Jn returns the order-n Bessel function of the first kind.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// Special cases are:</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//	Jn(n, ±Inf) = 0</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//	Jn(n, NaN) = NaN</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func Jn(n int, x float64) float64 {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	const (
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		TwoM29 = 1.0 / (1 &lt;&lt; 29) <span class="comment">// 2**-29 0x3e10000000000000</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		Two302 = 1 &lt;&lt; 302        <span class="comment">// 2**302 0x52D0000000000000</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// special cases</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	switch {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	case IsNaN(x):
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		return x
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	case IsInf(x, 0):
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		return 0
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// J(-n, x) = (-1)**n * J(n, x), J(n, -x) = (-1)**n * J(n, x)</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// Thus, J(-n, x) = J(n, -x)</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	if n == 0 {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		return J0(x)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	if x == 0 {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		return 0
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	if n &lt; 0 {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		n, x = -n, -x
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	if n == 1 {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		return J1(x)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	sign := false
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		x = -x
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		if n&amp;1 == 1 {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			sign = true <span class="comment">// odd n and negative x</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	var b float64
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	if float64(n) &lt;= x {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		<span class="comment">// Safe to use J(n+1,x)=2n/x *J(n,x)-J(n-1,x)</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		if x &gt;= Two302 { <span class="comment">// x &gt; 2**302</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			<span class="comment">// (x &gt;&gt; n**2)</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			<span class="comment">//          Jn(x) = cos(x-(2n+1)*pi/4)*sqrt(2/x*pi)</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			<span class="comment">//          Yn(x) = sin(x-(2n+1)*pi/4)*sqrt(2/x*pi)</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			<span class="comment">//          Let s=sin(x), c=cos(x),</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			<span class="comment">//              xn=x-(2n+1)*pi/4, sqt2 = sqrt(2),then</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			<span class="comment">//                 n    sin(xn)*sqt2    cos(xn)*sqt2</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			<span class="comment">//              ----------------------------------</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			<span class="comment">//                 0     s-c             c+s</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			<span class="comment">//                 1    -s-c            -c+s</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			<span class="comment">//                 2    -s+c            -c-s</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>			<span class="comment">//                 3     s+c             c-s</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			var temp float64
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			switch s, c := Sincos(x); n &amp; 3 {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			case 0:
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>				temp = c + s
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			case 1:
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>				temp = -c + s
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			case 2:
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>				temp = -c - s
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			case 3:
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>				temp = c - s
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			b = (1 / SqrtPi) * temp / Sqrt(x)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		} else {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			b = J1(x)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			for i, a := 1, J0(x); i &lt; n; i++ {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>				a, b = b, b*(float64(i+i)/x)-a <span class="comment">// avoid underflow</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	} else {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		if x &lt; TwoM29 { <span class="comment">// x &lt; 2**-29</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			<span class="comment">// x is tiny, return the first Taylor expansion of J(n,x)</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			<span class="comment">// J(n,x) = 1/n!*(x/2)**n  - ...</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			if n &gt; 33 { <span class="comment">// underflow</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>				b = 0
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			} else {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>				temp := x * 0.5
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>				b = temp
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>				a := 1.0
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>				for i := 2; i &lt;= n; i++ {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>					a *= float64(i) <span class="comment">// a = n!</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>					b *= temp       <span class="comment">// b = (x/2)**n</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>				}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>				b /= a
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		} else {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			<span class="comment">// use backward recurrence</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			<span class="comment">//                      x      x**2      x**2</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			<span class="comment">//  J(n,x)/J(n-1,x) =  ----   ------   ------   .....</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			<span class="comment">//                      2n  - 2(n+1) - 2(n+2)</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			<span class="comment">//                      1      1        1</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			<span class="comment">//  (for large x)   =  ----  ------   ------   .....</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			<span class="comment">//                      2n   2(n+1)   2(n+2)</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			<span class="comment">//                      -- - ------ - ------ -</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			<span class="comment">//                       x     x         x</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			<span class="comment">// Let w = 2n/x and h=2/x, then the above quotient</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			<span class="comment">// is equal to the continued fraction:</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			<span class="comment">//                  1</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			<span class="comment">//      = -----------------------</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			<span class="comment">//                     1</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			<span class="comment">//         w - -----------------</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			<span class="comment">//                        1</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			<span class="comment">//              w+h - ---------</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			<span class="comment">//                     w+2h - ...</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			<span class="comment">// To determine how many terms needed, let</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			<span class="comment">// Q(0) = w, Q(1) = w(w+h) - 1,</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			<span class="comment">// Q(k) = (w+k*h)*Q(k-1) - Q(k-2),</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			<span class="comment">// When Q(k) &gt; 1e4	good for single</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>			<span class="comment">// When Q(k) &gt; 1e9	good for double</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			<span class="comment">// When Q(k) &gt; 1e17	good for quadruple</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			<span class="comment">// determine k</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			w := float64(n+n) / x
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			h := 2 / x
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			q0 := w
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			z := w + h
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			q1 := w*z - 1
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			k := 1
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			for q1 &lt; 1e9 {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>				k++
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>				z += h
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>				q0, q1 = q1, z*q1-q0
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			m := n + n
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			t := 0.0
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			for i := 2 * (n + k); i &gt;= m; i -= 2 {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>				t = 1 / (float64(i)/x - t)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			a := t
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			b = 1
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			<span class="comment">//  estimate log((2/x)**n*n!) = n*log(2/x)+n*ln(n)</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			<span class="comment">//  Hence, if n*(log(2n/x)) &gt; ...</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			<span class="comment">//  single 8.8722839355e+01</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			<span class="comment">//  double 7.09782712893383973096e+02</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			<span class="comment">//  long double 1.1356523406294143949491931077970765006170e+04</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			<span class="comment">//  then recurrent value may overflow and the result is</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			<span class="comment">//  likely underflow to zero</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			tmp := float64(n)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			v := 2 / x
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			tmp = tmp * Log(Abs(v*tmp))
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			if tmp &lt; 7.09782712893383973096e+02 {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>				for i := n - 1; i &gt; 0; i-- {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>					di := float64(i + i)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>					a, b = b, b*di/x-a
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>				}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			} else {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>				for i := n - 1; i &gt; 0; i-- {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>					di := float64(i + i)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>					a, b = b, b*di/x-a
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>					<span class="comment">// scale b to avoid spurious overflow</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>					if b &gt; 1e100 {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>						a /= b
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>						t /= b
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>						b = 1
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>					}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>				}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			b = t * J0(x) / b
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	if sign {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		return -b
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	return b
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// Yn returns the order-n Bessel function of the second kind.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// Special cases are:</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">//	Yn(n, +Inf) = 0</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">//	Yn(n ≥ 0, 0) = -Inf</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">//	Yn(n &lt; 0, 0) = +Inf if n is odd, -Inf if n is even</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">//	Yn(n, x &lt; 0) = NaN</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">//	Yn(n, NaN) = NaN</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>func Yn(n int, x float64) float64 {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	const Two302 = 1 &lt;&lt; 302 <span class="comment">// 2**302 0x52D0000000000000</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// special cases</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	switch {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	case x &lt; 0 || IsNaN(x):
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		return NaN()
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	case IsInf(x, 1):
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		return 0
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	if n == 0 {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		return Y0(x)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	if x == 0 {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		if n &lt; 0 &amp;&amp; n&amp;1 == 1 {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			return Inf(1)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		return Inf(-1)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	sign := false
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	if n &lt; 0 {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		n = -n
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		if n&amp;1 == 1 {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			sign = true <span class="comment">// sign true if n &lt; 0 &amp;&amp; |n| odd</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	if n == 1 {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		if sign {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			return -Y1(x)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		return Y1(x)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	var b float64
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	if x &gt;= Two302 { <span class="comment">// x &gt; 2**302</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		<span class="comment">// (x &gt;&gt; n**2)</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		<span class="comment">//	    Jn(x) = cos(x-(2n+1)*pi/4)*sqrt(2/x*pi)</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		<span class="comment">//	    Yn(x) = sin(x-(2n+1)*pi/4)*sqrt(2/x*pi)</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		<span class="comment">//	    Let s=sin(x), c=cos(x),</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		<span class="comment">//		xn=x-(2n+1)*pi/4, sqt2 = sqrt(2),then</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		<span class="comment">//		   n	sin(xn)*sqt2	cos(xn)*sqt2</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		<span class="comment">//		----------------------------------</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		<span class="comment">//		   0	 s-c		 c+s</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		<span class="comment">//		   1	-s-c 		-c+s</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		<span class="comment">//		   2	-s+c		-c-s</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		<span class="comment">//		   3	 s+c		 c-s</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		var temp float64
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		switch s, c := Sincos(x); n &amp; 3 {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		case 0:
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			temp = s - c
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		case 1:
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			temp = -s - c
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		case 2:
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			temp = -s + c
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		case 3:
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			temp = s + c
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		b = (1 / SqrtPi) * temp / Sqrt(x)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	} else {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		a := Y0(x)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		b = Y1(x)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		<span class="comment">// quit if b is -inf</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		for i := 1; i &lt; n &amp;&amp; !IsInf(b, -1); i++ {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			a, b = b, (float64(i+i)/x)*b-a
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	if sign {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		return -b
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	return b
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
</pre><p><a href="jn.go?m=text">View as plain text</a></p>

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
