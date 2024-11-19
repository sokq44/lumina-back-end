<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/expm1.go - Go Documentation Server</title>

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
<a href="expm1.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<span class="text-muted">expm1.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// The original C code, the long comment, and the constants</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// below are from FreeBSD&#39;s /usr/src/lib/msun/src/s_expm1.c</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// and came with this notice. The go code is a simplified</span>
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
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// expm1(x)</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// Returns exp(x)-1, the exponential of x minus 1.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// Method</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//   1. Argument reduction:</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//      Given x, find r and integer k such that</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//               x = k*ln2 + r,  |r| &lt;= 0.5*ln2 ~ 0.34658</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//      Here a correction term c will be computed to compensate</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//      the error in r when rounded to a floating-point number.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//   2. Approximating expm1(r) by a special rational function on</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//      the interval [0,0.34658]:</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//      Since</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//          r*(exp(r)+1)/(exp(r)-1) = 2+ r**2/6 - r**4/360 + ...</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//      we define R1(r*r) by</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//          r*(exp(r)+1)/(exp(r)-1) = 2+ r**2/6 * R1(r*r)</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//      That is,</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//          R1(r**2) = 6/r *((exp(r)+1)/(exp(r)-1) - 2/r)</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//                   = 6/r * ( 1 + 2.0*(1/(exp(r)-1) - 1/r))</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//                   = 1 - r**2/60 + r**4/2520 - r**6/100800 + ...</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//      We use a special Reme algorithm on [0,0.347] to generate</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//      a polynomial of degree 5 in r*r to approximate R1. The</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//      maximum error of this polynomial approximation is bounded</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//      by 2**-61. In other words,</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//          R1(z) ~ 1.0 + Q1*z + Q2*z**2 + Q3*z**3 + Q4*z**4 + Q5*z**5</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//      where   Q1  =  -1.6666666666666567384E-2,</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//              Q2  =   3.9682539681370365873E-4,</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//              Q3  =  -9.9206344733435987357E-6,</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//              Q4  =   2.5051361420808517002E-7,</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//              Q5  =  -6.2843505682382617102E-9;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//      (where z=r*r, and the values of Q1 to Q5 are listed below)</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//      with error bounded by</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//          |                  5           |     -61</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//          | 1.0+Q1*z+...+Q5*z   -  R1(z) | &lt;= 2</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//          |                              |</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//      expm1(r) = exp(r)-1 is then computed by the following</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//      specific way which minimize the accumulation rounding error:</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//                             2     3</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//                            r     r    [ 3 - (R1 + R1*r/2)  ]</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//            expm1(r) = r + --- + --- * [--------------------]</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//                            2     2    [ 6 - r*(3 - R1*r/2) ]</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">//      To compensate the error in the argument reduction, we use</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">//              expm1(r+c) = expm1(r) + c + expm1(r)*c</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//                         ~ expm1(r) + c + r*c</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//      Thus c+r*c will be added in as the correction terms for</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">//      expm1(r+c). Now rearrange the term to avoid optimization</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//      screw up:</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//                      (      2                                    2 )</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//                      ({  ( r    [ R1 -  (3 - R1*r/2) ]  )  }    r  )</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//       expm1(r+c)~r - ({r*(--- * [--------------------]-c)-c} - --- )</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//                      ({  ( 2    [ 6 - r*(3 - R1*r/2) ]  )  }    2  )</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//                      (                                             )</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//                 = r - E</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//   3. Scale back to obtain expm1(x):</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//      From step 1, we have</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//         expm1(x) = either 2**k*[expm1(r)+1] - 1</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//                  = or     2**k*[expm1(r) + (1-2**-k)]</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//   4. Implementation notes:</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">//      (A). To save one multiplication, we scale the coefficient Qi</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//           to Qi*2**i, and replace z by (x**2)/2.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//      (B). To achieve maximum accuracy, we compute expm1(x) by</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">//        (i)   if x &lt; -56*ln2, return -1.0, (raise inexact if x!=inf)</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//        (ii)  if k=0, return r-E</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//        (iii) if k=-1, return 0.5*(r-E)-0.5</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">//        (iv)  if k=1 if r &lt; -0.25, return 2*((r+0.5)- E)</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//                     else          return  1.0+2.0*(r-E);</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">//        (v)   if (k&lt;-2||k&gt;56) return 2**k(1-(E-r)) - 1 (or exp(x)-1)</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">//        (vi)  if k &lt;= 20, return 2**k((1-2**-k)-(E-r)), else</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">//        (vii) return 2**k(1-((E+2**-k)-r))</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// Special cases:</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">//      expm1(INF) is INF, expm1(NaN) is NaN;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//      expm1(-INF) is -1, and</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">//      for finite argument, only expm1(0)=0 is exact.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// Accuracy:</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">//      according to an error analysis, the error is always less than</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">//      1 ulp (unit in the last place).</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// Misc. info.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//      For IEEE double</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">//          if x &gt;  7.09782712893383973096e+02 then expm1(x) overflow</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// Constants:</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// The hexadecimal values are the intended ones for the following</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// constants. The decimal values may be used, provided that the</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// compiler will convert from decimal to binary accurately enough</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// to produce the hexadecimal values shown.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// Expm1 returns e**x - 1, the base-e exponential of x minus 1.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// It is more accurate than [Exp](x) - 1 when x is near zero.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// Special cases are:</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">//	Expm1(+Inf) = +Inf</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">//	Expm1(-Inf) = -1</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">//	Expm1(NaN) = NaN</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// Very large values overflow to -1 or +Inf.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>func Expm1(x float64) float64 {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if haveArchExpm1 {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		return archExpm1(x)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	return expm1(x)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>func expm1(x float64) float64 {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	const (
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		Othreshold = 7.09782712893383973096e+02 <span class="comment">// 0x40862E42FEFA39EF</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		Ln2X56     = 3.88162421113569373274e+01 <span class="comment">// 0x4043687a9f1af2b1</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		Ln2HalfX3  = 1.03972077083991796413e+00 <span class="comment">// 0x3ff0a2b23f3bab73</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		Ln2Half    = 3.46573590279972654709e-01 <span class="comment">// 0x3fd62e42fefa39ef</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		Ln2Hi      = 6.93147180369123816490e-01 <span class="comment">// 0x3fe62e42fee00000</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		Ln2Lo      = 1.90821492927058770002e-10 <span class="comment">// 0x3dea39ef35793c76</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		InvLn2     = 1.44269504088896338700e+00 <span class="comment">// 0x3ff71547652b82fe</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		Tiny       = 1.0 / (1 &lt;&lt; 54)            <span class="comment">// 2**-54 = 0x3c90000000000000</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		<span class="comment">// scaled coefficients related to expm1</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		Q1 = -3.33333333333331316428e-02 <span class="comment">// 0xBFA11111111110F4</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		Q2 = 1.58730158725481460165e-03  <span class="comment">// 0x3F5A01A019FE5585</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		Q3 = -7.93650757867487942473e-05 <span class="comment">// 0xBF14CE199EAADBB7</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		Q4 = 4.00821782732936239552e-06  <span class="comment">// 0x3ED0CFCA86E65239</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		Q5 = -2.01099218183624371326e-07 <span class="comment">// 0xBE8AFDB76E09C32D</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// special cases</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	switch {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	case IsInf(x, 1) || IsNaN(x):
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return x
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	case IsInf(x, -1):
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		return -1
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	absx := x
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	sign := false
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		absx = -absx
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		sign = true
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// filter out huge argument</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	if absx &gt;= Ln2X56 { <span class="comment">// if |x| &gt;= 56 * ln2</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		if sign {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			return -1 <span class="comment">// x &lt; -56*ln2, return -1</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		if absx &gt;= Othreshold { <span class="comment">// if |x| &gt;= 709.78...</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			return Inf(1)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// argument reduction</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	var c float64
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	var k int
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	if absx &gt; Ln2Half { <span class="comment">// if  |x| &gt; 0.5 * ln2</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		var hi, lo float64
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		if absx &lt; Ln2HalfX3 { <span class="comment">// and |x| &lt; 1.5 * ln2</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			if !sign {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>				hi = x - Ln2Hi
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>				lo = Ln2Lo
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>				k = 1
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			} else {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>				hi = x + Ln2Hi
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				lo = -Ln2Lo
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				k = -1
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		} else {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			if !sign {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>				k = int(InvLn2*x + 0.5)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			} else {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>				k = int(InvLn2*x - 0.5)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			t := float64(k)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			hi = x - t*Ln2Hi <span class="comment">// t * Ln2Hi is exact here</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			lo = t * Ln2Lo
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		x = hi - lo
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		c = (hi - x) - lo
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	} else if absx &lt; Tiny { <span class="comment">// when |x| &lt; 2**-54, return x</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		return x
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	} else {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		k = 0
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">// x is now in primary range</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	hfx := 0.5 * x
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	hxs := x * hfx
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	r1 := 1 + hxs*(Q1+hxs*(Q2+hxs*(Q3+hxs*(Q4+hxs*Q5))))
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	t := 3 - r1*hfx
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	e := hxs * ((r1 - t) / (6.0 - x*t))
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	if k == 0 {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		return x - (x*e - hxs) <span class="comment">// c is 0</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	e = (x*(e-c) - c)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	e -= hxs
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	switch {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	case k == -1:
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		return 0.5*(x-e) - 0.5
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	case k == 1:
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		if x &lt; -0.25 {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			return -2 * (e - (x + 0.5))
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		return 1 + 2*(x-e)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	case k &lt;= -2 || k &gt; 56: <span class="comment">// suffice to return exp(x)-1</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		y := 1 - (e - x)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		y = Float64frombits(Float64bits(y) + uint64(k)&lt;&lt;52) <span class="comment">// add k to y&#39;s exponent</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		return y - 1
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	if k &lt; 20 {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		t := Float64frombits(0x3ff0000000000000 - (0x20000000000000 &gt;&gt; uint(k))) <span class="comment">// t=1-2**-k</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		y := t - (e - x)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		y = Float64frombits(Float64bits(y) + uint64(k)&lt;&lt;52) <span class="comment">// add k to y&#39;s exponent</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		return y
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	t = Float64frombits(uint64(0x3ff-k) &lt;&lt; 52) <span class="comment">// 2**-k</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	y := x - (e + t)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	y++
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	y = Float64frombits(Float64bits(y) + uint64(k)&lt;&lt;52) <span class="comment">// add k to y&#39;s exponent</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	return y
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
</pre><p><a href="expm1.go?m=text">View as plain text</a></p>

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
