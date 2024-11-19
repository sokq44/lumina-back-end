<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/rand/normal.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../index.html">GoDoc</a></div>
<a href="normal.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<a href="http://localhost:8080/src/math/rand">rand</a>/<span class="text-muted">normal.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math/rand">math/rand</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package rand
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>)
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">/*
<span id="L12" class="ln">    12&nbsp;&nbsp;</span> * Normal distribution
<span id="L13" class="ln">    13&nbsp;&nbsp;</span> *
<span id="L14" class="ln">    14&nbsp;&nbsp;</span> * See &#34;The Ziggurat Method for Generating Random Variables&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span> * (Marsaglia &amp; Tsang, 2000)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span> * http://www.jstatsoft.org/v05/i08/paper [pdf]
<span id="L17" class="ln">    17&nbsp;&nbsp;</span> */</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>const (
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	rn = 3.442619855899
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>func absInt32(i int32) uint32 {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		return uint32(-i)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	}
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	return uint32(i)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// NormFloat64 returns a normally distributed float64 in</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// the range -[math.MaxFloat64] through +[math.MaxFloat64] inclusive,</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// with standard normal distribution (mean = 0, stddev = 1).</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// To produce a different normal distribution, callers can</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// adjust the output using:</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//	sample = NormFloat64() * desiredStdDev + desiredMean</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>func (r *Rand) NormFloat64() float64 {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	for {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		j := int32(r.Uint32()) <span class="comment">// Possibly negative</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		i := j &amp; 0x7F
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		x := float64(j) * float64(wn[i])
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		if absInt32(j) &lt; kn[i] {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>			<span class="comment">// This case should be hit better than 99% of the time.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>			return x
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		if i == 0 {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			<span class="comment">// This extra work is only required for the base strip.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			for {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>				x = -math.Log(r.Float64()) * (1.0 / rn)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>				y := -math.Log(r.Float64())
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>				if y+y &gt;= x*x {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>					break
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>				}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>			}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			if j &gt; 0 {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>				return rn + x
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			return -rn - x
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		if fn[i]+float32(r.Float64())*(fn[i-1]-fn[i]) &lt; float32(math.Exp(-.5*x*x)) {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			return x
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>var kn = [128]uint32{
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	0x76ad2212, 0x0, 0x600f1b53, 0x6ce447a6, 0x725b46a2,
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	0x7560051d, 0x774921eb, 0x789a25bd, 0x799045c3, 0x7a4bce5d,
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	0x7adf629f, 0x7b5682a6, 0x7bb8a8c6, 0x7c0ae722, 0x7c50cce7,
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	0x7c8cec5b, 0x7cc12cd6, 0x7ceefed2, 0x7d177e0b, 0x7d3b8883,
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	0x7d5bce6c, 0x7d78dd64, 0x7d932886, 0x7dab0e57, 0x7dc0dd30,
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	0x7dd4d688, 0x7de73185, 0x7df81cea, 0x7e07c0a3, 0x7e163efa,
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	0x7e23b587, 0x7e303dfd, 0x7e3beec2, 0x7e46db77, 0x7e51155d,
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	0x7e5aabb3, 0x7e63abf7, 0x7e6c222c, 0x7e741906, 0x7e7b9a18,
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	0x7e82adfa, 0x7e895c63, 0x7e8fac4b, 0x7e95a3fb, 0x7e9b4924,
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	0x7ea0a0ef, 0x7ea5b00d, 0x7eaa7ac3, 0x7eaf04f3, 0x7eb3522a,
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	0x7eb765a5, 0x7ebb4259, 0x7ebeeafd, 0x7ec2620a, 0x7ec5a9c4,
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	0x7ec8c441, 0x7ecbb365, 0x7ece78ed, 0x7ed11671, 0x7ed38d62,
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	0x7ed5df12, 0x7ed80cb4, 0x7eda175c, 0x7edc0005, 0x7eddc78e,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	0x7edf6ebf, 0x7ee0f647, 0x7ee25ebe, 0x7ee3a8a9, 0x7ee4d473,
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	0x7ee5e276, 0x7ee6d2f5, 0x7ee7a620, 0x7ee85c10, 0x7ee8f4cd,
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	0x7ee97047, 0x7ee9ce59, 0x7eea0eca, 0x7eea3147, 0x7eea3568,
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	0x7eea1aab, 0x7ee9e071, 0x7ee98602, 0x7ee90a88, 0x7ee86d08,
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	0x7ee7ac6a, 0x7ee6c769, 0x7ee5bc9c, 0x7ee48a67, 0x7ee32efc,
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	0x7ee1a857, 0x7edff42f, 0x7ede0ffa, 0x7edbf8d9, 0x7ed9ab94,
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	0x7ed7248d, 0x7ed45fae, 0x7ed1585c, 0x7ece095f, 0x7eca6ccb,
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	0x7ec67be2, 0x7ec22eee, 0x7ebd7d1a, 0x7eb85c35, 0x7eb2c075,
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	0x7eac9c20, 0x7ea5df27, 0x7e9e769f, 0x7e964c16, 0x7e8d44ba,
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	0x7e834033, 0x7e781728, 0x7e6b9933, 0x7e5d8a1a, 0x7e4d9ded,
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	0x7e3b737a, 0x7e268c2f, 0x7e0e3ff5, 0x7df1aa5d, 0x7dcf8c72,
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	0x7da61a1e, 0x7d72a0fb, 0x7d30e097, 0x7cd9b4ab, 0x7c600f1a,
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	0x7ba90bdc, 0x7a722176, 0x77d664e5,
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>var wn = [128]float32{
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	1.7290405e-09, 1.2680929e-10, 1.6897518e-10, 1.9862688e-10,
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	2.2232431e-10, 2.4244937e-10, 2.601613e-10, 2.7611988e-10,
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	2.9073963e-10, 3.042997e-10, 3.1699796e-10, 3.289802e-10,
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	3.4035738e-10, 3.5121603e-10, 3.616251e-10, 3.7164058e-10,
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	3.8130857e-10, 3.9066758e-10, 3.9975012e-10, 4.08584e-10,
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	4.1719309e-10, 4.2559822e-10, 4.338176e-10, 4.418672e-10,
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	4.497613e-10, 4.5751258e-10, 4.651324e-10, 4.7263105e-10,
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	4.8001775e-10, 4.87301e-10, 4.944885e-10, 5.015873e-10,
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	5.0860405e-10, 5.155446e-10, 5.2241467e-10, 5.2921934e-10,
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	5.359635e-10, 5.426517e-10, 5.4928817e-10, 5.5587696e-10,
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	5.624219e-10, 5.6892646e-10, 5.753941e-10, 5.818282e-10,
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	5.882317e-10, 5.946077e-10, 6.00959e-10, 6.072884e-10,
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	6.135985e-10, 6.19892e-10, 6.2617134e-10, 6.3243905e-10,
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	6.386974e-10, 6.449488e-10, 6.511956e-10, 6.5744005e-10,
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	6.6368433e-10, 6.699307e-10, 6.7618144e-10, 6.824387e-10,
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	6.8870465e-10, 6.949815e-10, 7.012715e-10, 7.075768e-10,
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	7.1389966e-10, 7.202424e-10, 7.266073e-10, 7.329966e-10,
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	7.394128e-10, 7.4585826e-10, 7.5233547e-10, 7.58847e-10,
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	7.653954e-10, 7.719835e-10, 7.7861395e-10, 7.852897e-10,
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	7.920138e-10, 7.987892e-10, 8.0561924e-10, 8.125073e-10,
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	8.194569e-10, 8.2647167e-10, 8.3355556e-10, 8.407127e-10,
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	8.479473e-10, 8.55264e-10, 8.6266755e-10, 8.7016316e-10,
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	8.777562e-10, 8.8545243e-10, 8.932582e-10, 9.0117996e-10,
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	9.09225e-10, 9.174008e-10, 9.2571584e-10, 9.341788e-10,
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	9.427997e-10, 9.515889e-10, 9.605579e-10, 9.697193e-10,
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	9.790869e-10, 9.88676e-10, 9.985036e-10, 1.0085882e-09,
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	1.0189509e-09, 1.0296151e-09, 1.0406069e-09, 1.0519566e-09,
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	1.063698e-09, 1.0758702e-09, 1.0885183e-09, 1.1016947e-09,
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	1.1154611e-09, 1.1298902e-09, 1.1450696e-09, 1.1611052e-09,
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	1.1781276e-09, 1.1962995e-09, 1.2158287e-09, 1.2369856e-09,
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	1.2601323e-09, 1.2857697e-09, 1.3146202e-09, 1.347784e-09,
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	1.3870636e-09, 1.4357403e-09, 1.5008659e-09, 1.6030948e-09,
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>var fn = [128]float32{
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	1, 0.9635997, 0.9362827, 0.9130436, 0.89228165, 0.87324303,
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	0.8555006, 0.8387836, 0.8229072, 0.8077383, 0.793177,
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	0.7791461, 0.7655842, 0.7524416, 0.73967725, 0.7272569,
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	0.7151515, 0.7033361, 0.69178915, 0.68049186, 0.6694277,
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	0.658582, 0.6479418, 0.63749546, 0.6272325, 0.6171434,
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	0.6072195, 0.5974532, 0.58783704, 0.5783647, 0.56903,
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	0.5598274, 0.5507518, 0.54179835, 0.5329627, 0.52424055,
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	0.5156282, 0.50712204, 0.49871865, 0.49041483, 0.48220766,
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	0.4740943, 0.46607214, 0.4581387, 0.45029163, 0.44252872,
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	0.43484783, 0.427247, 0.41972435, 0.41227803, 0.40490642,
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	0.39760786, 0.3903808, 0.3832238, 0.37613547, 0.36911446,
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	0.3621595, 0.35526937, 0.34844297, 0.34167916, 0.33497685,
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	0.3283351, 0.3217529, 0.3152294, 0.30876362, 0.30235484,
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	0.29600215, 0.28970486, 0.2834622, 0.2772735, 0.27113807,
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	0.2650553, 0.25902456, 0.2530453, 0.24711695, 0.241239,
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	0.23541094, 0.22963232, 0.2239027, 0.21822165, 0.21258877,
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	0.20700371, 0.20146611, 0.19597565, 0.19053204, 0.18513499,
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	0.17978427, 0.17447963, 0.1692209, 0.16400786, 0.15884037,
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	0.15371831, 0.14864157, 0.14361008, 0.13862377, 0.13368265,
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	0.12878671, 0.12393598, 0.119130544, 0.11437051, 0.10965602,
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	0.104987256, 0.10036444, 0.095787846, 0.0912578, 0.08677467,
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	0.0823389, 0.077950984, 0.073611505, 0.06932112, 0.06508058,
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	0.06089077, 0.056752663, 0.0526674, 0.048636295, 0.044660863,
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	0.040742867, 0.03688439, 0.033087887, 0.029356318,
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	0.025693292, 0.022103304, 0.018592102, 0.015167298,
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	0.011839478, 0.008624485, 0.005548995, 0.0026696292,
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
</pre><p><a href="normal.go?m=text">View as plain text</a></p>

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
