<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/lgamma.go - Go Documentation Server</title>

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
<a href="lgamma.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<span class="text-muted">lgamma.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	Floating-point logarithm of the Gamma function.
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>*/</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// The original C code and the long comment below are</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// from FreeBSD&#39;s /usr/src/lib/msun/src/e_lgamma_r.c and</span>
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
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// __ieee754_lgamma_r(x, signgamp)</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// Reentrant version of the logarithm of the Gamma function</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// with user provided pointer for the sign of Gamma(x).</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// Method:</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//   1. Argument Reduction for 0 &lt; x &lt;= 8</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//      Since gamma(1+s)=s*gamma(s), for x in [0,8], we may</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//      reduce x to a number in [1.5,2.5] by</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//              lgamma(1+s) = log(s) + lgamma(s)</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//      for example,</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//              lgamma(7.3) = log(6.3) + lgamma(6.3)</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//                          = log(6.3*5.3) + lgamma(5.3)</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//                          = log(6.3*5.3*4.3*3.3*2.3) + lgamma(2.3)</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//   2. Polynomial approximation of lgamma around its</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//      minimum (ymin=1.461632144968362245) to maintain monotonicity.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//      On [ymin-0.23, ymin+0.27] (i.e., [1.23164,1.73163]), use</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//              Let z = x-ymin;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//              lgamma(x) = -1.214862905358496078218 + z**2*poly(z)</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//              poly(z) is a 14 degree polynomial.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//   2. Rational approximation in the primary interval [2,3]</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//      We use the following approximation:</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//              s = x-2.0;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//              lgamma(x) = 0.5*s + s*P(s)/Q(s)</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//      with accuracy</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//              |P/Q - (lgamma(x)-0.5s)| &lt; 2**-61.71</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//      Our algorithms are based on the following observation</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//                             zeta(2)-1    2    zeta(3)-1    3</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// lgamma(2+s) = s*(1-Euler) + --------- * s  -  --------- * s  + ...</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//                                 2                 3</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//      where Euler = 0.5772156649... is the Euler constant, which</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//      is very close to 0.5.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//   3. For x&gt;=8, we have</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//      lgamma(x)~(x-0.5)log(x)-x+0.5*log(2pi)+1/(12x)-1/(360x**3)+....</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//      (better formula:</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//         lgamma(x)~(x-0.5)*(log(x)-1)-.5*(log(2pi)-1) + ...)</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//      Let z = 1/x, then we approximation</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//              f(z) = lgamma(x) - (x-0.5)(log(x)-1)</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//      by</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">//                                  3       5             11</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">//              w = w0 + w1*z + w2*z  + w3*z  + ... + w6*z</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//      where</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//              |w - f(z)| &lt; 2**-58.74</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//   4. For negative x, since (G is gamma function)</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//              -x*G(-x)*G(x) = pi/sin(pi*x),</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//      we have</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//              G(x) = pi/(sin(pi*x)*(-x)*G(-x))</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//      since G(-x) is positive, sign(G(x)) = sign(sin(pi*x)) for x&lt;0</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//      Hence, for x&lt;0, signgam = sign(sin(pi*x)) and</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//              lgamma(x) = log(|Gamma(x)|)</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//                        = log(pi/(|x*sin(pi*x)|)) - lgamma(-x);</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//      Note: one should avoid computing pi*(-x) directly in the</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//            computation of sin(pi*(-x)).</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//   5. Special Cases</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//              lgamma(2+s) ~ s*(1-Euler) for tiny s</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">//              lgamma(1)=lgamma(2)=0</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//              lgamma(x) ~ -log(x) for tiny x</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//              lgamma(0) = lgamma(inf) = inf</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">//              lgamma(-integer) = +-inf</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>var _lgamA = [...]float64{
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	7.72156649015328655494e-02, <span class="comment">// 0x3FB3C467E37DB0C8</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	3.22467033424113591611e-01, <span class="comment">// 0x3FD4A34CC4A60FAD</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	6.73523010531292681824e-02, <span class="comment">// 0x3FB13E001A5562A7</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	2.05808084325167332806e-02, <span class="comment">// 0x3F951322AC92547B</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	7.38555086081402883957e-03, <span class="comment">// 0x3F7E404FB68FEFE8</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	2.89051383673415629091e-03, <span class="comment">// 0x3F67ADD8CCB7926B</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	1.19270763183362067845e-03, <span class="comment">// 0x3F538A94116F3F5D</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	5.10069792153511336608e-04, <span class="comment">// 0x3F40B6C689B99C00</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	2.20862790713908385557e-04, <span class="comment">// 0x3F2CF2ECED10E54D</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	1.08011567247583939954e-04, <span class="comment">// 0x3F1C5088987DFB07</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	2.52144565451257326939e-05, <span class="comment">// 0x3EFA7074428CFA52</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	4.48640949618915160150e-05, <span class="comment">// 0x3F07858E90A45837</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>var _lgamR = [...]float64{
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	1.0,                        <span class="comment">// placeholder</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	1.39200533467621045958e+00, <span class="comment">// 0x3FF645A762C4AB74</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	7.21935547567138069525e-01, <span class="comment">// 0x3FE71A1893D3DCDC</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	1.71933865632803078993e-01, <span class="comment">// 0x3FC601EDCCFBDF27</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	1.86459191715652901344e-02, <span class="comment">// 0x3F9317EA742ED475</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	7.77942496381893596434e-04, <span class="comment">// 0x3F497DDACA41A95B</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	7.32668430744625636189e-06, <span class="comment">// 0x3EDEBAF7A5B38140</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>var _lgamS = [...]float64{
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	-7.72156649015328655494e-02, <span class="comment">// 0xBFB3C467E37DB0C8</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	2.14982415960608852501e-01,  <span class="comment">// 0x3FCB848B36E20878</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	3.25778796408930981787e-01,  <span class="comment">// 0x3FD4D98F4F139F59</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	1.46350472652464452805e-01,  <span class="comment">// 0x3FC2BB9CBEE5F2F7</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	2.66422703033638609560e-02,  <span class="comment">// 0x3F9B481C7E939961</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	1.84028451407337715652e-03,  <span class="comment">// 0x3F5E26B67368F239</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	3.19475326584100867617e-05,  <span class="comment">// 0x3F00BFECDD17E945</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>var _lgamT = [...]float64{
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	4.83836122723810047042e-01,  <span class="comment">// 0x3FDEF72BC8EE38A2</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	-1.47587722994593911752e-01, <span class="comment">// 0xBFC2E4278DC6C509</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	6.46249402391333854778e-02,  <span class="comment">// 0x3FB08B4294D5419B</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	-3.27885410759859649565e-02, <span class="comment">// 0xBFA0C9A8DF35B713</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	1.79706750811820387126e-02,  <span class="comment">// 0x3F9266E7970AF9EC</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	-1.03142241298341437450e-02, <span class="comment">// 0xBF851F9FBA91EC6A</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	6.10053870246291332635e-03,  <span class="comment">// 0x3F78FCE0E370E344</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	-3.68452016781138256760e-03, <span class="comment">// 0xBF6E2EFFB3E914D7</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	2.25964780900612472250e-03,  <span class="comment">// 0x3F6282D32E15C915</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	-1.40346469989232843813e-03, <span class="comment">// 0xBF56FE8EBF2D1AF1</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	8.81081882437654011382e-04,  <span class="comment">// 0x3F4CDF0CEF61A8E9</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	-5.38595305356740546715e-04, <span class="comment">// 0xBF41A6109C73E0EC</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	3.15632070903625950361e-04,  <span class="comment">// 0x3F34AF6D6C0EBBF7</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	-3.12754168375120860518e-04, <span class="comment">// 0xBF347F24ECC38C38</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	3.35529192635519073543e-04,  <span class="comment">// 0x3F35FD3EE8C2D3F4</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>var _lgamU = [...]float64{
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	-7.72156649015328655494e-02, <span class="comment">// 0xBFB3C467E37DB0C8</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	6.32827064025093366517e-01,  <span class="comment">// 0x3FE4401E8B005DFF</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	1.45492250137234768737e+00,  <span class="comment">// 0x3FF7475CD119BD6F</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	9.77717527963372745603e-01,  <span class="comment">// 0x3FEF497644EA8450</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	2.28963728064692451092e-01,  <span class="comment">// 0x3FCD4EAEF6010924</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	1.33810918536787660377e-02,  <span class="comment">// 0x3F8B678BBF2BAB09</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>var _lgamV = [...]float64{
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	1.0,
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	2.45597793713041134822e+00, <span class="comment">// 0x4003A5D7C2BD619C</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	2.12848976379893395361e+00, <span class="comment">// 0x40010725A42B18F5</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	7.69285150456672783825e-01, <span class="comment">// 0x3FE89DFBE45050AF</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	1.04222645593369134254e-01, <span class="comment">// 0x3FBAAE55D6537C88</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	3.21709242282423911810e-03, <span class="comment">// 0x3F6A5ABB57D0CF61</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>var _lgamW = [...]float64{
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	4.18938533204672725052e-01,  <span class="comment">// 0x3FDACFE390C97D69</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	8.33333333333329678849e-02,  <span class="comment">// 0x3FB555555555553B</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	-2.77777777728775536470e-03, <span class="comment">// 0xBF66C16C16B02E5C</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	7.93650558643019558500e-04,  <span class="comment">// 0x3F4A019F98CF38B6</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	-5.95187557450339963135e-04, <span class="comment">// 0xBF4380CB8C0FE741</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	8.36339918996282139126e-04,  <span class="comment">// 0x3F4B67BA4CDAD5D1</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	-1.63092934096575273989e-03, <span class="comment">// 0xBF5AB89D0B9E43E4</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// Lgamma returns the natural logarithm and sign (-1 or +1) of [Gamma](x).</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span><span class="comment">// Special cases are:</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">//	Lgamma(+Inf) = +Inf</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">//	Lgamma(0) = +Inf</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">//	Lgamma(-integer) = +Inf</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">//	Lgamma(-Inf) = -Inf</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">//	Lgamma(NaN) = NaN</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>func Lgamma(x float64) (lgamma float64, sign int) {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	const (
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		Ymin  = 1.461632144968362245
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		Two52 = 1 &lt;&lt; 52                     <span class="comment">// 0x4330000000000000 ~4.5036e+15</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		Two53 = 1 &lt;&lt; 53                     <span class="comment">// 0x4340000000000000 ~9.0072e+15</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		Two58 = 1 &lt;&lt; 58                     <span class="comment">// 0x4390000000000000 ~2.8823e+17</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		Tiny  = 1.0 / (1 &lt;&lt; 70)             <span class="comment">// 0x3b90000000000000 ~8.47033e-22</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		Tc    = 1.46163214496836224576e+00  <span class="comment">// 0x3FF762D86356BE3F</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		Tf    = -1.21486290535849611461e-01 <span class="comment">// 0xBFBF19B9BCC38A42</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		<span class="comment">// Tt = -(tail of Tf)</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		Tt = -3.63867699703950536541e-18 <span class="comment">// 0xBC50C7CAA48A971F</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// special cases</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	sign = 1
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	switch {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	case IsNaN(x):
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		lgamma = x
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		return
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	case IsInf(x, 0):
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		lgamma = x
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		return
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	case x == 0:
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		lgamma = Inf(1)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		return
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	neg := false
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		x = -x
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		neg = true
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	if x &lt; Tiny { <span class="comment">// if |x| &lt; 2**-70, return -log(|x|)</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		if neg {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			sign = -1
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		lgamma = -Log(x)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		return
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	var nadj float64
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	if neg {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		if x &gt;= Two52 { <span class="comment">// |x| &gt;= 2**52, must be -integer</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			lgamma = Inf(1)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			return
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		t := sinPi(x)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		if t == 0 {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			lgamma = Inf(1) <span class="comment">// -integer</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			return
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		nadj = Log(Pi / Abs(t*x))
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		if t &lt; 0 {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			sign = -1
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	switch {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	case x == 1 || x == 2: <span class="comment">// purge off 1 and 2</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		lgamma = 0
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		return
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	case x &lt; 2: <span class="comment">// use lgamma(x) = lgamma(x+1) - log(x)</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		var y float64
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		var i int
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		if x &lt;= 0.9 {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			lgamma = -Log(x)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			switch {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			case x &gt;= (Ymin - 1 + 0.27): <span class="comment">// 0.7316 &lt;= x &lt;=  0.9</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>				y = 1 - x
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>				i = 0
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			case x &gt;= (Ymin - 1 - 0.27): <span class="comment">// 0.2316 &lt;= x &lt; 0.7316</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>				y = x - (Tc - 1)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>				i = 1
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			default: <span class="comment">// 0 &lt; x &lt; 0.2316</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>				y = x
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>				i = 2
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		} else {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			lgamma = 0
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			switch {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			case x &gt;= (Ymin + 0.27): <span class="comment">// 1.7316 &lt;= x &lt; 2</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				y = 2 - x
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>				i = 0
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			case x &gt;= (Ymin - 0.27): <span class="comment">// 1.2316 &lt;= x &lt; 1.7316</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>				y = x - Tc
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>				i = 1
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			default: <span class="comment">// 0.9 &lt; x &lt; 1.2316</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>				y = x - 1
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				i = 2
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		switch i {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		case 0:
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			z := y * y
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			p1 := _lgamA[0] + z*(_lgamA[2]+z*(_lgamA[4]+z*(_lgamA[6]+z*(_lgamA[8]+z*_lgamA[10]))))
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			p2 := z * (_lgamA[1] + z*(+_lgamA[3]+z*(_lgamA[5]+z*(_lgamA[7]+z*(_lgamA[9]+z*_lgamA[11])))))
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			p := y*p1 + p2
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			lgamma += (p - 0.5*y)
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		case 1:
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			z := y * y
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			w := z * y
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			p1 := _lgamT[0] + w*(_lgamT[3]+w*(_lgamT[6]+w*(_lgamT[9]+w*_lgamT[12]))) <span class="comment">// parallel comp</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			p2 := _lgamT[1] + w*(_lgamT[4]+w*(_lgamT[7]+w*(_lgamT[10]+w*_lgamT[13])))
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			p3 := _lgamT[2] + w*(_lgamT[5]+w*(_lgamT[8]+w*(_lgamT[11]+w*_lgamT[14])))
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			p := z*p1 - (Tt - w*(p2+y*p3))
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			lgamma += (Tf + p)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		case 2:
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			p1 := y * (_lgamU[0] + y*(_lgamU[1]+y*(_lgamU[2]+y*(_lgamU[3]+y*(_lgamU[4]+y*_lgamU[5])))))
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			p2 := 1 + y*(_lgamV[1]+y*(_lgamV[2]+y*(_lgamV[3]+y*(_lgamV[4]+y*_lgamV[5]))))
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			lgamma += (-0.5*y + p1/p2)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	case x &lt; 8: <span class="comment">// 2 &lt;= x &lt; 8</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		i := int(x)
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		y := x - float64(i)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		p := y * (_lgamS[0] + y*(_lgamS[1]+y*(_lgamS[2]+y*(_lgamS[3]+y*(_lgamS[4]+y*(_lgamS[5]+y*_lgamS[6]))))))
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		q := 1 + y*(_lgamR[1]+y*(_lgamR[2]+y*(_lgamR[3]+y*(_lgamR[4]+y*(_lgamR[5]+y*_lgamR[6])))))
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		lgamma = 0.5*y + p/q
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		z := 1.0 <span class="comment">// Lgamma(1+s) = Log(s) + Lgamma(s)</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		switch i {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		case 7:
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			z *= (y + 6)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			fallthrough
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		case 6:
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			z *= (y + 5)
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			fallthrough
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		case 5:
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			z *= (y + 4)
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			fallthrough
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		case 4:
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			z *= (y + 3)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			fallthrough
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		case 3:
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			z *= (y + 2)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			lgamma += Log(z)
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	case x &lt; Two58: <span class="comment">// 8 &lt;= x &lt; 2**58</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		t := Log(x)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		z := 1 / x
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		y := z * z
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		w := _lgamW[0] + z*(_lgamW[1]+y*(_lgamW[2]+y*(_lgamW[3]+y*(_lgamW[4]+y*(_lgamW[5]+y*_lgamW[6])))))
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		lgamma = (x-0.5)*(t-1) + w
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	default: <span class="comment">// 2**58 &lt;= x &lt;= Inf</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		lgamma = x * (Log(x) - 1)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	if neg {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		lgamma = nadj - lgamma
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	return
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span><span class="comment">// sinPi(x) is a helper function for negative x</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>func sinPi(x float64) float64 {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	const (
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		Two52 = 1 &lt;&lt; 52 <span class="comment">// 0x4330000000000000 ~4.5036e+15</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		Two53 = 1 &lt;&lt; 53 <span class="comment">// 0x4340000000000000 ~9.0072e+15</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	if x &lt; 0.25 {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		return -Sin(Pi * x)
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// argument reduction</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	z := Floor(x)
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	var n int
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	if z != x { <span class="comment">// inexact</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		x = Mod(x, 2)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		n = int(x * 4)
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	} else {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		if x &gt;= Two53 { <span class="comment">// x must be even</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			x = 0
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			n = 0
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		} else {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			if x &lt; Two52 {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>				z = x + Two52 <span class="comment">// exact</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>			n = int(1 &amp; Float64bits(z))
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>			x = float64(n)
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>			n &lt;&lt;= 2
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	switch n {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	case 0:
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		x = Sin(Pi * x)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	case 1, 2:
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		x = Cos(Pi * (0.5 - x))
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	case 3, 4:
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		x = Sin(Pi * (1 - x))
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	case 5, 6:
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		x = -Cos(Pi * (x - 1.5))
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	default:
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		x = Sin(Pi * (x - 2))
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	}
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	return -x
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
</pre><p><a href="lgamma.go?m=text">View as plain text</a></p>

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
