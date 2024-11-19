<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/des/const.go - Go Documentation Server</title>

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
<a href="const.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/des">des</a>/<span class="text-muted">const.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/des">crypto/des</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2010 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package des implements the Data Encryption Standard (DES) and the</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// Triple Data Encryption Algorithm (TDEA) as defined</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// in U.S. Federal Information Processing Standards Publication 46-3.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// DES is cryptographically broken and should not be used for secure</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// applications.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>package des
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// Used to perform an initial permutation of a 64-bit input block.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>var initialPermutation = [64]byte{
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	6, 14, 22, 30, 38, 46, 54, 62,
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	4, 12, 20, 28, 36, 44, 52, 60,
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	2, 10, 18, 26, 34, 42, 50, 58,
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	0, 8, 16, 24, 32, 40, 48, 56,
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	7, 15, 23, 31, 39, 47, 55, 63,
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	5, 13, 21, 29, 37, 45, 53, 61,
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	3, 11, 19, 27, 35, 43, 51, 59,
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	1, 9, 17, 25, 33, 41, 49, 57,
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>}
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// Used to perform a final permutation of a 4-bit preoutput block. This is the</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// inverse of initialPermutation</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>var finalPermutation = [64]byte{
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	24, 56, 16, 48, 8, 40, 0, 32,
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	25, 57, 17, 49, 9, 41, 1, 33,
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	26, 58, 18, 50, 10, 42, 2, 34,
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	27, 59, 19, 51, 11, 43, 3, 35,
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	28, 60, 20, 52, 12, 44, 4, 36,
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	29, 61, 21, 53, 13, 45, 5, 37,
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	30, 62, 22, 54, 14, 46, 6, 38,
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	31, 63, 23, 55, 15, 47, 7, 39,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// Used to expand an input block of 32 bits, producing an output block of 48</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// bits.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>var expansionFunction = [48]byte{
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	0, 31, 30, 29, 28, 27, 28, 27,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	26, 25, 24, 23, 24, 23, 22, 21,
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	20, 19, 20, 19, 18, 17, 16, 15,
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	16, 15, 14, 13, 12, 11, 12, 11,
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	10, 9, 8, 7, 8, 7, 6, 5,
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	4, 3, 4, 3, 2, 1, 0, 31,
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// Yields a 32-bit output from a 32-bit input</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>var permutationFunction = [32]byte{
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	16, 25, 12, 11, 3, 20, 4, 15,
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	31, 17, 9, 6, 27, 14, 1, 22,
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	30, 24, 8, 18, 0, 5, 29, 23,
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	13, 19, 2, 26, 10, 21, 28, 7,
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// Used in the key schedule to select 56 bits</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// from a 64-bit input.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>var permutedChoice1 = [56]byte{
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	7, 15, 23, 31, 39, 47, 55, 63,
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	6, 14, 22, 30, 38, 46, 54, 62,
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	5, 13, 21, 29, 37, 45, 53, 61,
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	4, 12, 20, 28, 1, 9, 17, 25,
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	33, 41, 49, 57, 2, 10, 18, 26,
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	34, 42, 50, 58, 3, 11, 19, 27,
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	35, 43, 51, 59, 36, 44, 52, 60,
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// Used in the key schedule to produce each subkey by selecting 48 bits from</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// the 56-bit input</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>var permutedChoice2 = [48]byte{
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	42, 39, 45, 32, 55, 51, 53, 28,
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	41, 50, 35, 46, 33, 37, 44, 52,
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	30, 48, 40, 49, 29, 36, 43, 54,
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	15, 4, 25, 19, 9, 1, 26, 16,
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	5, 11, 23, 8, 12, 7, 17, 0,
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	22, 3, 10, 14, 6, 20, 27, 24,
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// 8 S-boxes composed of 4 rows and 16 columns</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// Used in the DES cipher function</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>var sBoxes = [8][4][16]uint8{
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// S-box 1</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	{
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		{14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7},
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		{0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8},
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		{4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0},
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		{15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13},
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	},
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// S-box 2</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	{
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		{15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10},
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		{3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5},
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		{0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15},
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		{13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9},
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	},
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// S-box 3</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	{
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		{10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8},
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		{13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1},
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		{13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7},
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		{1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12},
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	},
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// S-box 4</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	{
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		{7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15},
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		{13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9},
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		{10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4},
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		{3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14},
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	},
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// S-box 5</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	{
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		{2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9},
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		{14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6},
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		{4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14},
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		{11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3},
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	},
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// S-box 6</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	{
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		{12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11},
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		{10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8},
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		{9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6},
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		{4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13},
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	},
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// S-box 7</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	{
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		{4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1},
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		{13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6},
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		{1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2},
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		{6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12},
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	},
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// S-box 8</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	{
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		{13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7},
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		{1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2},
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		{7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8},
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		{2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11},
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	},
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// Size of left rotation per round in each half of the key schedule</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>var ksRotations = [16]uint8{1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
</pre><p><a href="const.go?m=text">View as plain text</a></p>

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
