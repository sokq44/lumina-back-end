<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/erf.go - Go Documentation Server</title>

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
<a href="erf.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<span class="text-muted">erf.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	Floating-point error function and complementary error function.
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>*/</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// The original C code and the long comment below are</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// from FreeBSD&#39;s /usr/src/lib/msun/src/s_erf.c and</span>
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
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// double erf(double x)</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// double erfc(double x)</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//                           x</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//                    2      |\</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//     erf(x)  =  ---------  | exp(-t*t)dt</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//                 sqrt(pi) \|</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//                           0</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//     erfc(x) =  1-erf(x)</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//  Note that</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//              erf(-x) = -erf(x)</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//              erfc(-x) = 2 - erfc(x)</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// Method:</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//      1. For |x| in [0, 0.84375]</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//          erf(x)  = x + x*R(x**2)</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//          erfc(x) = 1 - erf(x)           if x in [-.84375,0.25]</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//                  = 0.5 + ((0.5-x)-x*R)  if x in [0.25,0.84375]</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//         where R = P/Q where P is an odd poly of degree 8 and</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//         Q is an odd poly of degree 10.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//                                               -57.90</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//                      | R - (erf(x)-x)/x | &lt;= 2</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//         Remark. The formula is derived by noting</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//          erf(x) = (2/sqrt(pi))*(x - x**3/3 + x**5/10 - x**7/42 + ....)</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//         and that</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//          2/sqrt(pi) = 1.128379167095512573896158903121545171688</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//         is close to one. The interval is chosen because the fix</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//         point of erf(x) is near 0.6174 (i.e., erf(x)=x when x is</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//         near 0.6174), and by some experiment, 0.84375 is chosen to</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//         guarantee the error is less than one ulp for erf.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//      2. For |x| in [0.84375,1.25], let s = |x| - 1, and</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//         c = 0.84506291151 rounded to single (24 bits)</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//              erf(x)  = sign(x) * (c  + P1(s)/Q1(s))</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//              erfc(x) = (1-c)  - P1(s)/Q1(s) if x &gt; 0</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//                        1+(c+P1(s)/Q1(s))    if x &lt; 0</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//              |P1/Q1 - (erf(|x|)-c)| &lt;= 2**-59.06</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//         Remark: here we use the taylor series expansion at x=1.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">//              erf(1+s) = erf(1) + s*Poly(s)</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">//                       = 0.845.. + P1(s)/Q1(s)</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//         That is, we use rational approximation to approximate</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//                      erf(1+s) - (c = (single)0.84506291151)</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">//         Note that |P1/Q1|&lt; 0.078 for x in [0.84375,1.25]</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//         where</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//              P1(s) = degree 6 poly in s</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//              Q1(s) = degree 6 poly in s</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//      3. For x in [1.25,1/0.35(~2.857143)],</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//              erfc(x) = (1/x)*exp(-x*x-0.5625+R1/S1)</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//              erf(x)  = 1 - erfc(x)</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//         where</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//              R1(z) = degree 7 poly in z, (z=1/x**2)</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//              S1(z) = degree 8 poly in z</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//      4. For x in [1/0.35,28]</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//              erfc(x) = (1/x)*exp(-x*x-0.5625+R2/S2) if x &gt; 0</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">//                      = 2.0 - (1/x)*exp(-x*x-0.5625+R2/S2) if -6&lt;x&lt;0</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//                      = 2.0 - tiny            (if x &lt;= -6)</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//              erf(x)  = sign(x)*(1.0 - erfc(x)) if x &lt; 6, else</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">//              erf(x)  = sign(x)*(1.0 - tiny)</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//         where</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//              R2(z) = degree 6 poly in z, (z=1/x**2)</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">//              S2(z) = degree 7 poly in z</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">//      Note1:</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">//         To compute exp(-x*x-0.5625+R/S), let s be a single</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">//         precision number and s := x; then</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">//              -x*x = -s*s + (s-x)*(s+x)</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">//              exp(-x*x-0.5626+R/S) =</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">//                      exp(-s*s-0.5625)*exp((s-x)*(s+x)+R/S);</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//      Note2:</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">//         Here 4 and 5 make use of the asymptotic series</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">//                        exp(-x*x)</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">//              erfc(x) ~ ---------- * ( 1 + Poly(1/x**2) )</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">//                        x*sqrt(pi)</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">//         We use rational approximation to approximate</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">//              g(s)=f(1/x**2) = log(erfc(x)*x) - x*x + 0.5625</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">//         Here is the error bound for R1/S1 and R2/S2</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//              |R1/S1 - f(x)|  &lt; 2**(-62.57)</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">//              |R2/S2 - f(x)|  &lt; 2**(-61.52)</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">//      5. For inf &gt; x &gt;= 28</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">//              erf(x)  = sign(x) *(1 - tiny)  (raise inexact)</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//              erfc(x) = tiny*tiny (raise underflow) if x &gt; 0</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//                      = 2 - tiny if x&lt;0</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">//      7. Special case:</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">//              erf(0)  = 0, erf(inf)  = 1, erf(-inf) = -1,</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">//              erfc(0) = 1, erfc(inf) = 0, erfc(-inf) = 2,</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">//              erfc/erf(NaN) is NaN</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>const (
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	erx = 8.45062911510467529297e-01 <span class="comment">// 0x3FEB0AC160000000</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// Coefficients for approximation to  erf in [0, 0.84375]</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	efx  = 1.28379167095512586316e-01  <span class="comment">// 0x3FC06EBA8214DB69</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	efx8 = 1.02703333676410069053e+00  <span class="comment">// 0x3FF06EBA8214DB69</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	pp0  = 1.28379167095512558561e-01  <span class="comment">// 0x3FC06EBA8214DB68</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	pp1  = -3.25042107247001499370e-01 <span class="comment">// 0xBFD4CD7D691CB913</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	pp2  = -2.84817495755985104766e-02 <span class="comment">// 0xBF9D2A51DBD7194F</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	pp3  = -5.77027029648944159157e-03 <span class="comment">// 0xBF77A291236668E4</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	pp4  = -2.37630166566501626084e-05 <span class="comment">// 0xBEF8EAD6120016AC</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	qq1  = 3.97917223959155352819e-01  <span class="comment">// 0x3FD97779CDDADC09</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	qq2  = 6.50222499887672944485e-02  <span class="comment">// 0x3FB0A54C5536CEBA</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	qq3  = 5.08130628187576562776e-03  <span class="comment">// 0x3F74D022C4D36B0F</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	qq4  = 1.32494738004321644526e-04  <span class="comment">// 0x3F215DC9221C1A10</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	qq5  = -3.96022827877536812320e-06 <span class="comment">// 0xBED09C4342A26120</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// Coefficients for approximation to  erf  in [0.84375, 1.25]</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	pa0 = -2.36211856075265944077e-03 <span class="comment">// 0xBF6359B8BEF77538</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	pa1 = 4.14856118683748331666e-01  <span class="comment">// 0x3FDA8D00AD92B34D</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	pa2 = -3.72207876035701323847e-01 <span class="comment">// 0xBFD7D240FBB8C3F1</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	pa3 = 3.18346619901161753674e-01  <span class="comment">// 0x3FD45FCA805120E4</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	pa4 = -1.10894694282396677476e-01 <span class="comment">// 0xBFBC63983D3E28EC</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	pa5 = 3.54783043256182359371e-02  <span class="comment">// 0x3FA22A36599795EB</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	pa6 = -2.16637559486879084300e-03 <span class="comment">// 0xBF61BF380A96073F</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	qa1 = 1.06420880400844228286e-01  <span class="comment">// 0x3FBB3E6618EEE323</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	qa2 = 5.40397917702171048937e-01  <span class="comment">// 0x3FE14AF092EB6F33</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	qa3 = 7.18286544141962662868e-02  <span class="comment">// 0x3FB2635CD99FE9A7</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	qa4 = 1.26171219808761642112e-01  <span class="comment">// 0x3FC02660E763351F</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	qa5 = 1.36370839120290507362e-02  <span class="comment">// 0x3F8BEDC26B51DD1C</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	qa6 = 1.19844998467991074170e-02  <span class="comment">// 0x3F888B545735151D</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// Coefficients for approximation to  erfc in [1.25, 1/0.35]</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	ra0 = -9.86494403484714822705e-03 <span class="comment">// 0xBF843412600D6435</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	ra1 = -6.93858572707181764372e-01 <span class="comment">// 0xBFE63416E4BA7360</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	ra2 = -1.05586262253232909814e+01 <span class="comment">// 0xC0251E0441B0E726</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	ra3 = -6.23753324503260060396e+01 <span class="comment">// 0xC04F300AE4CBA38D</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	ra4 = -1.62396669462573470355e+02 <span class="comment">// 0xC0644CB184282266</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	ra5 = -1.84605092906711035994e+02 <span class="comment">// 0xC067135CEBCCABB2</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	ra6 = -8.12874355063065934246e+01 <span class="comment">// 0xC054526557E4D2F2</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	ra7 = -9.81432934416914548592e+00 <span class="comment">// 0xC023A0EFC69AC25C</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	sa1 = 1.96512716674392571292e+01  <span class="comment">// 0x4033A6B9BD707687</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	sa2 = 1.37657754143519042600e+02  <span class="comment">// 0x4061350C526AE721</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	sa3 = 4.34565877475229228821e+02  <span class="comment">// 0x407B290DD58A1A71</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	sa4 = 6.45387271733267880336e+02  <span class="comment">// 0x40842B1921EC2868</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	sa5 = 4.29008140027567833386e+02  <span class="comment">// 0x407AD02157700314</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	sa6 = 1.08635005541779435134e+02  <span class="comment">// 0x405B28A3EE48AE2C</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	sa7 = 6.57024977031928170135e+00  <span class="comment">// 0x401A47EF8E484A93</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	sa8 = -6.04244152148580987438e-02 <span class="comment">// 0xBFAEEFF2EE749A62</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// Coefficients for approximation to  erfc in [1/.35, 28]</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	rb0 = -9.86494292470009928597e-03 <span class="comment">// 0xBF84341239E86F4A</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	rb1 = -7.99283237680523006574e-01 <span class="comment">// 0xBFE993BA70C285DE</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	rb2 = -1.77579549177547519889e+01 <span class="comment">// 0xC031C209555F995A</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	rb3 = -1.60636384855821916062e+02 <span class="comment">// 0xC064145D43C5ED98</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	rb4 = -6.37566443368389627722e+02 <span class="comment">// 0xC083EC881375F228</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	rb5 = -1.02509513161107724954e+03 <span class="comment">// 0xC09004616A2E5992</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	rb6 = -4.83519191608651397019e+02 <span class="comment">// 0xC07E384E9BDC383F</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	sb1 = 3.03380607434824582924e+01  <span class="comment">// 0x403E568B261D5190</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	sb2 = 3.25792512996573918826e+02  <span class="comment">// 0x40745CAE221B9F0A</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	sb3 = 1.53672958608443695994e+03  <span class="comment">// 0x409802EB189D5118</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	sb4 = 3.19985821950859553908e+03  <span class="comment">// 0x40A8FFB7688C246A</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	sb5 = 2.55305040643316442583e+03  <span class="comment">// 0x40A3F219CEDF3BE6</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	sb6 = 4.74528541206955367215e+02  <span class="comment">// 0x407DA874E79FE763</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	sb7 = -2.24409524465858183362e+01 <span class="comment">// 0xC03670E242712D62</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">// Erf returns the error function of x.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span><span class="comment">// Special cases are:</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">//	Erf(+Inf) = 1</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">//	Erf(-Inf) = -1</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">//	Erf(NaN) = NaN</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>func Erf(x float64) float64 {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	if haveArchErf {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		return archErf(x)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	return erf(x)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>func erf(x float64) float64 {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	const (
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		VeryTiny = 2.848094538889218e-306 <span class="comment">// 0x0080000000000000</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		Small    = 1.0 / (1 &lt;&lt; 28)        <span class="comment">// 2**-28</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	)
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// special cases</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	switch {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	case IsNaN(x):
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		return NaN()
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	case IsInf(x, 1):
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		return 1
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	case IsInf(x, -1):
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		return -1
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	sign := false
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		x = -x
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		sign = true
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	if x &lt; 0.84375 { <span class="comment">// |x| &lt; 0.84375</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		var temp float64
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		if x &lt; Small { <span class="comment">// |x| &lt; 2**-28</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			if x &lt; VeryTiny {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>				temp = 0.125 * (8.0*x + efx8*x) <span class="comment">// avoid underflow</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			} else {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>				temp = x + efx*x
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		} else {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			z := x * x
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			r := pp0 + z*(pp1+z*(pp2+z*(pp3+z*pp4)))
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			s := 1 + z*(qq1+z*(qq2+z*(qq3+z*(qq4+z*qq5))))
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			y := r / s
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			temp = x + x*y
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		if sign {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			return -temp
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		return temp
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	if x &lt; 1.25 { <span class="comment">// 0.84375 &lt;= |x| &lt; 1.25</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		s := x - 1
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		P := pa0 + s*(pa1+s*(pa2+s*(pa3+s*(pa4+s*(pa5+s*pa6)))))
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		Q := 1 + s*(qa1+s*(qa2+s*(qa3+s*(qa4+s*(qa5+s*qa6)))))
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		if sign {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			return -erx - P/Q
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		return erx + P/Q
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	if x &gt;= 6 { <span class="comment">// inf &gt; |x| &gt;= 6</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		if sign {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			return -1
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		return 1
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	s := 1 / (x * x)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	var R, S float64
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	if x &lt; 1/0.35 { <span class="comment">// |x| &lt; 1 / 0.35  ~ 2.857143</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		R = ra0 + s*(ra1+s*(ra2+s*(ra3+s*(ra4+s*(ra5+s*(ra6+s*ra7))))))
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		S = 1 + s*(sa1+s*(sa2+s*(sa3+s*(sa4+s*(sa5+s*(sa6+s*(sa7+s*sa8)))))))
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	} else { <span class="comment">// |x| &gt;= 1 / 0.35  ~ 2.857143</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		R = rb0 + s*(rb1+s*(rb2+s*(rb3+s*(rb4+s*(rb5+s*rb6)))))
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		S = 1 + s*(sb1+s*(sb2+s*(sb3+s*(sb4+s*(sb5+s*(sb6+s*sb7))))))
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	z := Float64frombits(Float64bits(x) &amp; 0xffffffff00000000) <span class="comment">// pseudo-single (20-bit) precision x</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	r := Exp(-z*z-0.5625) * Exp((z-x)*(z+x)+R/S)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	if sign {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		return r/x - 1
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	return 1 - r/x
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span><span class="comment">// Erfc returns the complementary error function of x.</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">// Special cases are:</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span><span class="comment">//	Erfc(+Inf) = 0</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span><span class="comment">//	Erfc(-Inf) = 2</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span><span class="comment">//	Erfc(NaN) = NaN</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>func Erfc(x float64) float64 {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	if haveArchErfc {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		return archErfc(x)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	return erfc(x)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>func erfc(x float64) float64 {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	const Tiny = 1.0 / (1 &lt;&lt; 56) <span class="comment">// 2**-56</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	<span class="comment">// special cases</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	switch {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	case IsNaN(x):
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		return NaN()
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	case IsInf(x, 1):
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		return 0
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	case IsInf(x, -1):
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		return 2
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	sign := false
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		x = -x
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		sign = true
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	if x &lt; 0.84375 { <span class="comment">// |x| &lt; 0.84375</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		var temp float64
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		if x &lt; Tiny { <span class="comment">// |x| &lt; 2**-56</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			temp = x
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		} else {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			z := x * x
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			r := pp0 + z*(pp1+z*(pp2+z*(pp3+z*pp4)))
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			s := 1 + z*(qq1+z*(qq2+z*(qq3+z*(qq4+z*qq5))))
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			y := r / s
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			if x &lt; 0.25 { <span class="comment">// |x| &lt; 1/4</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>				temp = x + x*y
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>			} else {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>				temp = 0.5 + (x*y + (x - 0.5))
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>			}
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		}
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		if sign {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>			return 1 + temp
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		return 1 - temp
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	if x &lt; 1.25 { <span class="comment">// 0.84375 &lt;= |x| &lt; 1.25</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		s := x - 1
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		P := pa0 + s*(pa1+s*(pa2+s*(pa3+s*(pa4+s*(pa5+s*pa6)))))
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		Q := 1 + s*(qa1+s*(qa2+s*(qa3+s*(qa4+s*(qa5+s*qa6)))))
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		if sign {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>			return 1 + erx + P/Q
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		return 1 - erx - P/Q
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	if x &lt; 28 { <span class="comment">// |x| &lt; 28</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		s := 1 / (x * x)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		var R, S float64
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		if x &lt; 1/0.35 { <span class="comment">// |x| &lt; 1 / 0.35 ~ 2.857143</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			R = ra0 + s*(ra1+s*(ra2+s*(ra3+s*(ra4+s*(ra5+s*(ra6+s*ra7))))))
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			S = 1 + s*(sa1+s*(sa2+s*(sa3+s*(sa4+s*(sa5+s*(sa6+s*(sa7+s*sa8)))))))
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		} else { <span class="comment">// |x| &gt;= 1 / 0.35 ~ 2.857143</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>			if sign &amp;&amp; x &gt; 6 {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>				return 2 <span class="comment">// x &lt; -6</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			R = rb0 + s*(rb1+s*(rb2+s*(rb3+s*(rb4+s*(rb5+s*rb6)))))
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>			S = 1 + s*(sb1+s*(sb2+s*(sb3+s*(sb4+s*(sb5+s*(sb6+s*sb7))))))
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		z := Float64frombits(Float64bits(x) &amp; 0xffffffff00000000) <span class="comment">// pseudo-single (20-bit) precision x</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		r := Exp(-z*z-0.5625) * Exp((z-x)*(z+x)+R/S)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		if sign {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			return 2 - r/x
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		return r / x
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	if sign {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		return 2
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	return 0
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
</pre><p><a href="erf.go?m=text">View as plain text</a></p>

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
