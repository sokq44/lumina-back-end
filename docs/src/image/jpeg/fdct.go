<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/image/jpeg/fdct.go - Go Documentation Server</title>

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
<a href="fdct.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/image">image</a>/<a href="http://localhost:8080/src/image/jpeg">jpeg</a>/<span class="text-muted">fdct.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/image/jpeg">image/jpeg</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package jpeg
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This file implements a Forward Discrete Cosine Transformation.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">/*
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>It is based on the code in jfdctint.c from the Independent JPEG Group,
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>found at http://www.ijg.org/files/jpegsrc.v8c.tar.gz.
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>The &#34;LEGAL ISSUES&#34; section of the README in that archive says:
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>In plain English:
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>1. We don&#39;t promise that this software works.  (But if you find any bugs,
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>   please let us know!)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>2. You can use this software for whatever you want.  You don&#39;t have to pay us.
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>3. You may not pretend that you wrote this software.  If you use it in a
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>   program, you must acknowledge somewhere in your documentation that
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>   you&#39;ve used the IJG code.
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>In legalese:
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>The authors make NO WARRANTY or representation, either express or implied,
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>with respect to this software, its quality, accuracy, merchantability, or
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>fitness for a particular purpose.  This software is provided &#34;AS IS&#34;, and you,
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>its user, assume the entire risk as to its quality and accuracy.
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>This software is copyright (C) 1991-2011, Thomas G. Lane, Guido Vollbeding.
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>All Rights Reserved except as specified below.
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>Permission is hereby granted to use, copy, modify, and distribute this
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>software (or portions thereof) for any purpose, without fee, subject to these
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>conditions:
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>(1) If any part of the source code for this software is distributed, then this
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>README file must be included, with this copyright and no-warranty notice
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>unaltered; and any additions, deletions, or changes to the original files
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>must be clearly indicated in accompanying documentation.
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>(2) If only executable code is distributed, then the accompanying
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>documentation must state that &#34;this software is based in part on the work of
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>the Independent JPEG Group&#34;.
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>(3) Permission for use of this software is granted only if the user accepts
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>full responsibility for any undesirable consequences; the authors accept
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>NO LIABILITY for damages of any kind.
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>These conditions apply to any software derived from or based on the IJG code,
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>not just to the unmodified library.  If you use our work, you ought to
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>acknowledge us.
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>Permission is NOT granted for the use of any IJG author&#39;s name or company name
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>in advertising or publicity relating to this software or products derived from
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>it.  This software may be referred to only as &#34;the Independent JPEG Group&#39;s
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>software&#34;.
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>We specifically permit and encourage the use of this software as the basis of
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>commercial products, provided that all warranty or liability claims are
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>assumed by the product vendor.
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>*/</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// Trigonometric constants in 13-bit fixed point format.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>const (
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	fix_0_298631336 = 2446
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	fix_0_390180644 = 3196
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	fix_0_541196100 = 4433
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	fix_0_765366865 = 6270
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	fix_0_899976223 = 7373
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	fix_1_175875602 = 9633
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	fix_1_501321110 = 12299
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	fix_1_847759065 = 15137
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	fix_1_961570560 = 16069
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	fix_2_053119869 = 16819
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	fix_2_562915447 = 20995
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	fix_3_072711026 = 25172
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>const (
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	constBits     = 13
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	pass1Bits     = 2
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	centerJSample = 128
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// fdct performs a forward DCT on an 8x8 block of coefficients, including a</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// level shift.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>func fdct(b *block) {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// Pass 1: process rows.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	for y := 0; y &lt; 8; y++ {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		y8 := y * 8
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		s := b[y8 : y8+8 : y8+8] <span class="comment">// Small cap improves performance, see https://golang.org/issue/27857</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		x0 := s[0]
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		x1 := s[1]
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		x2 := s[2]
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		x3 := s[3]
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		x4 := s[4]
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		x5 := s[5]
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		x6 := s[6]
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		x7 := s[7]
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		tmp0 := x0 + x7
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		tmp1 := x1 + x6
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		tmp2 := x2 + x5
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		tmp3 := x3 + x4
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		tmp10 := tmp0 + tmp3
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		tmp12 := tmp0 - tmp3
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		tmp11 := tmp1 + tmp2
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		tmp13 := tmp1 - tmp2
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		tmp0 = x0 - x7
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		tmp1 = x1 - x6
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		tmp2 = x2 - x5
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		tmp3 = x3 - x4
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		s[0] = (tmp10 + tmp11 - 8*centerJSample) &lt;&lt; pass1Bits
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		s[4] = (tmp10 - tmp11) &lt;&lt; pass1Bits
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		z1 := (tmp12 + tmp13) * fix_0_541196100
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		z1 += 1 &lt;&lt; (constBits - pass1Bits - 1)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		s[2] = (z1 + tmp12*fix_0_765366865) &gt;&gt; (constBits - pass1Bits)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		s[6] = (z1 - tmp13*fix_1_847759065) &gt;&gt; (constBits - pass1Bits)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		tmp10 = tmp0 + tmp3
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		tmp11 = tmp1 + tmp2
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		tmp12 = tmp0 + tmp2
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		tmp13 = tmp1 + tmp3
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		z1 = (tmp12 + tmp13) * fix_1_175875602
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		z1 += 1 &lt;&lt; (constBits - pass1Bits - 1)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		tmp0 *= fix_1_501321110
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		tmp1 *= fix_3_072711026
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		tmp2 *= fix_2_053119869
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		tmp3 *= fix_0_298631336
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		tmp10 *= -fix_0_899976223
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		tmp11 *= -fix_2_562915447
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		tmp12 *= -fix_0_390180644
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		tmp13 *= -fix_1_961570560
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		tmp12 += z1
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		tmp13 += z1
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		s[1] = (tmp0 + tmp10 + tmp12) &gt;&gt; (constBits - pass1Bits)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		s[3] = (tmp1 + tmp11 + tmp13) &gt;&gt; (constBits - pass1Bits)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		s[5] = (tmp2 + tmp11 + tmp12) &gt;&gt; (constBits - pass1Bits)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		s[7] = (tmp3 + tmp10 + tmp13) &gt;&gt; (constBits - pass1Bits)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// Pass 2: process columns.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// We remove pass1Bits scaling, but leave results scaled up by an overall factor of 8.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	for x := 0; x &lt; 8; x++ {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		tmp0 := b[0*8+x] + b[7*8+x]
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		tmp1 := b[1*8+x] + b[6*8+x]
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		tmp2 := b[2*8+x] + b[5*8+x]
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		tmp3 := b[3*8+x] + b[4*8+x]
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		tmp10 := tmp0 + tmp3 + 1&lt;&lt;(pass1Bits-1)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		tmp12 := tmp0 - tmp3
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		tmp11 := tmp1 + tmp2
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		tmp13 := tmp1 - tmp2
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		tmp0 = b[0*8+x] - b[7*8+x]
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		tmp1 = b[1*8+x] - b[6*8+x]
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		tmp2 = b[2*8+x] - b[5*8+x]
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		tmp3 = b[3*8+x] - b[4*8+x]
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		b[0*8+x] = (tmp10 + tmp11) &gt;&gt; pass1Bits
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		b[4*8+x] = (tmp10 - tmp11) &gt;&gt; pass1Bits
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		z1 := (tmp12 + tmp13) * fix_0_541196100
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		z1 += 1 &lt;&lt; (constBits + pass1Bits - 1)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		b[2*8+x] = (z1 + tmp12*fix_0_765366865) &gt;&gt; (constBits + pass1Bits)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		b[6*8+x] = (z1 - tmp13*fix_1_847759065) &gt;&gt; (constBits + pass1Bits)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		tmp10 = tmp0 + tmp3
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		tmp11 = tmp1 + tmp2
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		tmp12 = tmp0 + tmp2
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		tmp13 = tmp1 + tmp3
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		z1 = (tmp12 + tmp13) * fix_1_175875602
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		z1 += 1 &lt;&lt; (constBits + pass1Bits - 1)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		tmp0 *= fix_1_501321110
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		tmp1 *= fix_3_072711026
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		tmp2 *= fix_2_053119869
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		tmp3 *= fix_0_298631336
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		tmp10 *= -fix_0_899976223
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		tmp11 *= -fix_2_562915447
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		tmp12 *= -fix_0_390180644
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		tmp13 *= -fix_1_961570560
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		tmp12 += z1
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		tmp13 += z1
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		b[1*8+x] = (tmp0 + tmp10 + tmp12) &gt;&gt; (constBits + pass1Bits)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		b[3*8+x] = (tmp1 + tmp11 + tmp13) &gt;&gt; (constBits + pass1Bits)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		b[5*8+x] = (tmp2 + tmp11 + tmp12) &gt;&gt; (constBits + pass1Bits)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		b[7*8+x] = (tmp3 + tmp10 + tmp13) &gt;&gt; (constBits + pass1Bits)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
</pre><p><a href="fdct.go?m=text">View as plain text</a></p>

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
