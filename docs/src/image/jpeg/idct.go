<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/image/jpeg/idct.go - Go Documentation Server</title>

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
<a href="idct.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/image">image</a>/<a href="http://localhost:8080/src/image/jpeg">jpeg</a>/<span class="text-muted">idct.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/image/jpeg">image/jpeg</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package jpeg
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This is a Go translation of idct.c from</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// http://standards.iso.org/ittf/PubliclyAvailableStandards/ISO_IEC_13818-4_2004_Conformance_Testing/Video/verifier/mpeg2decode_960109.tar.gz</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// which carries the following notice:</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">/* Copyright (C) 1996, MPEG Software Simulation Group. All Rights Reserved. */</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">/*
<span id="L16" class="ln">    16&nbsp;&nbsp;</span> * Disclaimer of Warranty
<span id="L17" class="ln">    17&nbsp;&nbsp;</span> *
<span id="L18" class="ln">    18&nbsp;&nbsp;</span> * These software programs are available to the user without any license fee or
<span id="L19" class="ln">    19&nbsp;&nbsp;</span> * royalty on an &#34;as is&#34; basis.  The MPEG Software Simulation Group disclaims
<span id="L20" class="ln">    20&nbsp;&nbsp;</span> * any and all warranties, whether express, implied, or statuary, including any
<span id="L21" class="ln">    21&nbsp;&nbsp;</span> * implied warranties or merchantability or of fitness for a particular
<span id="L22" class="ln">    22&nbsp;&nbsp;</span> * purpose.  In no event shall the copyright-holder be liable for any
<span id="L23" class="ln">    23&nbsp;&nbsp;</span> * incidental, punitive, or consequential damages of any kind whatsoever
<span id="L24" class="ln">    24&nbsp;&nbsp;</span> * arising from the use of these programs.
<span id="L25" class="ln">    25&nbsp;&nbsp;</span> *
<span id="L26" class="ln">    26&nbsp;&nbsp;</span> * This disclaimer of warranty extends to the user of these programs and user&#39;s
<span id="L27" class="ln">    27&nbsp;&nbsp;</span> * customers, employees, agents, transferees, successors, and assigns.
<span id="L28" class="ln">    28&nbsp;&nbsp;</span> *
<span id="L29" class="ln">    29&nbsp;&nbsp;</span> * The MPEG Software Simulation Group does not represent or warrant that the
<span id="L30" class="ln">    30&nbsp;&nbsp;</span> * programs furnished hereunder are free of infringement of any third-party
<span id="L31" class="ln">    31&nbsp;&nbsp;</span> * patents.
<span id="L32" class="ln">    32&nbsp;&nbsp;</span> *
<span id="L33" class="ln">    33&nbsp;&nbsp;</span> * Commercial implementations of MPEG-1 and MPEG-2 video, including shareware,
<span id="L34" class="ln">    34&nbsp;&nbsp;</span> * are subject to royalty fees to patent holders.  Many of these patents are
<span id="L35" class="ln">    35&nbsp;&nbsp;</span> * general enough such that they are unavoidable regardless of implementation
<span id="L36" class="ln">    36&nbsp;&nbsp;</span> * design.
<span id="L37" class="ln">    37&nbsp;&nbsp;</span> *
<span id="L38" class="ln">    38&nbsp;&nbsp;</span> */</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>const blockSize = 64 <span class="comment">// A DCT block is 8x8.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>type block [blockSize]int32
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>const (
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	w1 = 2841 <span class="comment">// 2048*sqrt(2)*cos(1*pi/16)</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	w2 = 2676 <span class="comment">// 2048*sqrt(2)*cos(2*pi/16)</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	w3 = 2408 <span class="comment">// 2048*sqrt(2)*cos(3*pi/16)</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	w5 = 1609 <span class="comment">// 2048*sqrt(2)*cos(5*pi/16)</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	w6 = 1108 <span class="comment">// 2048*sqrt(2)*cos(6*pi/16)</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	w7 = 565  <span class="comment">// 2048*sqrt(2)*cos(7*pi/16)</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	w1pw7 = w1 + w7
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	w1mw7 = w1 - w7
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	w2pw6 = w2 + w6
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	w2mw6 = w2 - w6
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	w3pw5 = w3 + w5
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	w3mw5 = w3 - w5
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	r2 = 181 <span class="comment">// 256/sqrt(2)</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// idct performs a 2-D Inverse Discrete Cosine Transformation.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// The input coefficients should already have been multiplied by the</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// appropriate quantization table. We use fixed-point computation, with the</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// number of bits for the fractional component varying over the intermediate</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// stages.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// For more on the actual algorithm, see Z. Wang, &#34;Fast algorithms for the</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// discrete W transform and for the discrete Fourier transform&#34;, IEEE Trans. on</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// ASSP, Vol. ASSP- 32, pp. 803-816, Aug. 1984.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func idct(src *block) {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// Horizontal 1-D IDCT.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	for y := 0; y &lt; 8; y++ {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		y8 := y * 8
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		s := src[y8 : y8+8 : y8+8] <span class="comment">// Small cap improves performance, see https://golang.org/issue/27857</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		<span class="comment">// If all the AC components are zero, then the IDCT is trivial.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		if s[1] == 0 &amp;&amp; s[2] == 0 &amp;&amp; s[3] == 0 &amp;&amp;
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			s[4] == 0 &amp;&amp; s[5] == 0 &amp;&amp; s[6] == 0 &amp;&amp; s[7] == 0 {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			dc := s[0] &lt;&lt; 3
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			s[0] = dc
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			s[1] = dc
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			s[2] = dc
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			s[3] = dc
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			s[4] = dc
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			s[5] = dc
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			s[6] = dc
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			s[7] = dc
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			continue
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		<span class="comment">// Prescale.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		x0 := (s[0] &lt;&lt; 11) + 128
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		x1 := s[4] &lt;&lt; 11
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		x2 := s[6]
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		x3 := s[2]
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		x4 := s[1]
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		x5 := s[7]
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		x6 := s[5]
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		x7 := s[3]
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		<span class="comment">// Stage 1.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		x8 := w7 * (x4 + x5)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		x4 = x8 + w1mw7*x4
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		x5 = x8 - w1pw7*x5
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		x8 = w3 * (x6 + x7)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		x6 = x8 - w3mw5*x6
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		x7 = x8 - w3pw5*x7
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		<span class="comment">// Stage 2.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		x8 = x0 + x1
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		x0 -= x1
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		x1 = w6 * (x3 + x2)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		x2 = x1 - w2pw6*x2
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		x3 = x1 + w2mw6*x3
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		x1 = x4 + x6
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		x4 -= x6
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		x6 = x5 + x7
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		x5 -= x7
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		<span class="comment">// Stage 3.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		x7 = x8 + x3
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		x8 -= x3
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		x3 = x0 + x2
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		x0 -= x2
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		x2 = (r2*(x4+x5) + 128) &gt;&gt; 8
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		x4 = (r2*(x4-x5) + 128) &gt;&gt; 8
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		<span class="comment">// Stage 4.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		s[0] = (x7 + x1) &gt;&gt; 8
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		s[1] = (x3 + x2) &gt;&gt; 8
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		s[2] = (x0 + x4) &gt;&gt; 8
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		s[3] = (x8 + x6) &gt;&gt; 8
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		s[4] = (x8 - x6) &gt;&gt; 8
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		s[5] = (x0 - x4) &gt;&gt; 8
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		s[6] = (x3 - x2) &gt;&gt; 8
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		s[7] = (x7 - x1) &gt;&gt; 8
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// Vertical 1-D IDCT.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	for x := 0; x &lt; 8; x++ {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		<span class="comment">// Similar to the horizontal 1-D IDCT case, if all the AC components are zero, then the IDCT is trivial.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		<span class="comment">// However, after performing the horizontal 1-D IDCT, there are typically non-zero AC components, so</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">// we do not bother to check for the all-zero case.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		s := src[x : x+57 : x+57] <span class="comment">// Small cap improves performance, see https://golang.org/issue/27857</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		<span class="comment">// Prescale.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		y0 := (s[8*0] &lt;&lt; 8) + 8192
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		y1 := s[8*4] &lt;&lt; 8
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		y2 := s[8*6]
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		y3 := s[8*2]
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		y4 := s[8*1]
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		y5 := s[8*7]
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		y6 := s[8*5]
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		y7 := s[8*3]
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// Stage 1.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		y8 := w7*(y4+y5) + 4
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		y4 = (y8 + w1mw7*y4) &gt;&gt; 3
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		y5 = (y8 - w1pw7*y5) &gt;&gt; 3
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		y8 = w3*(y6+y7) + 4
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		y6 = (y8 - w3mw5*y6) &gt;&gt; 3
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		y7 = (y8 - w3pw5*y7) &gt;&gt; 3
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		<span class="comment">// Stage 2.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		y8 = y0 + y1
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		y0 -= y1
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		y1 = w6*(y3+y2) + 4
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		y2 = (y1 - w2pw6*y2) &gt;&gt; 3
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		y3 = (y1 + w2mw6*y3) &gt;&gt; 3
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		y1 = y4 + y6
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		y4 -= y6
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		y6 = y5 + y7
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		y5 -= y7
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		<span class="comment">// Stage 3.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		y7 = y8 + y3
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		y8 -= y3
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		y3 = y0 + y2
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		y0 -= y2
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		y2 = (r2*(y4+y5) + 128) &gt;&gt; 8
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		y4 = (r2*(y4-y5) + 128) &gt;&gt; 8
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		<span class="comment">// Stage 4.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		s[8*0] = (y7 + y1) &gt;&gt; 14
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		s[8*1] = (y3 + y2) &gt;&gt; 14
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		s[8*2] = (y0 + y4) &gt;&gt; 14
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		s[8*3] = (y8 + y6) &gt;&gt; 14
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		s[8*4] = (y8 - y6) &gt;&gt; 14
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		s[8*5] = (y0 - y4) &gt;&gt; 14
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		s[8*6] = (y3 - y2) &gt;&gt; 14
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		s[8*7] = (y7 - y1) &gt;&gt; 14
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
</pre><p><a href="idct.go?m=text">View as plain text</a></p>

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
