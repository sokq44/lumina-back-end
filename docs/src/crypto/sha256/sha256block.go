<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/sha256/sha256block.go - Go Documentation Server</title>

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
<a href="sha256block.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/sha256">sha256</a>/<span class="text-muted">sha256block.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/sha256">crypto/sha256</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// SHA256 block step.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// In its own file so that a faster assembly or C version</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// can be substituted easily.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package sha256
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import &#34;math/bits&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>var _K = []uint32{
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	0x428a2f98,
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	0x71374491,
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	0xb5c0fbcf,
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	0xe9b5dba5,
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	0x3956c25b,
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	0x59f111f1,
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	0x923f82a4,
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	0xab1c5ed5,
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	0xd807aa98,
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	0x12835b01,
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	0x243185be,
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	0x550c7dc3,
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	0x72be5d74,
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	0x80deb1fe,
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	0x9bdc06a7,
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	0xc19bf174,
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	0xe49b69c1,
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	0xefbe4786,
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	0x0fc19dc6,
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	0x240ca1cc,
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	0x2de92c6f,
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	0x4a7484aa,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	0x5cb0a9dc,
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	0x76f988da,
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	0x983e5152,
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	0xa831c66d,
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	0xb00327c8,
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	0xbf597fc7,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	0xc6e00bf3,
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	0xd5a79147,
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	0x06ca6351,
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	0x14292967,
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	0x27b70a85,
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	0x2e1b2138,
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	0x4d2c6dfc,
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	0x53380d13,
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	0x650a7354,
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	0x766a0abb,
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	0x81c2c92e,
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	0x92722c85,
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	0xa2bfe8a1,
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	0xa81a664b,
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	0xc24b8b70,
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	0xc76c51a3,
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	0xd192e819,
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	0xd6990624,
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	0xf40e3585,
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	0x106aa070,
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	0x19a4c116,
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	0x1e376c08,
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	0x2748774c,
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	0x34b0bcb5,
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	0x391c0cb3,
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	0x4ed8aa4a,
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	0x5b9cca4f,
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	0x682e6ff3,
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	0x748f82ee,
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	0x78a5636f,
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	0x84c87814,
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	0x8cc70208,
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	0x90befffa,
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	0xa4506ceb,
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	0xbef9a3f7,
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	0xc67178f2,
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>func blockGeneric(dig *digest, p []byte) {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	var w [64]uint32
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	h0, h1, h2, h3, h4, h5, h6, h7 := dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4], dig.h[5], dig.h[6], dig.h[7]
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	for len(p) &gt;= chunk {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		<span class="comment">// Can interlace the computation of w with the</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		<span class="comment">// rounds below if needed for speed.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		for i := 0; i &lt; 16; i++ {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			j := i * 4
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			w[i] = uint32(p[j])&lt;&lt;24 | uint32(p[j+1])&lt;&lt;16 | uint32(p[j+2])&lt;&lt;8 | uint32(p[j+3])
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		for i := 16; i &lt; 64; i++ {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			v1 := w[i-2]
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			t1 := (bits.RotateLeft32(v1, -17)) ^ (bits.RotateLeft32(v1, -19)) ^ (v1 &gt;&gt; 10)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			v2 := w[i-15]
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			t2 := (bits.RotateLeft32(v2, -7)) ^ (bits.RotateLeft32(v2, -18)) ^ (v2 &gt;&gt; 3)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			w[i] = t1 + w[i-7] + t2 + w[i-16]
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		a, b, c, d, e, f, g, h := h0, h1, h2, h3, h4, h5, h6, h7
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		for i := 0; i &lt; 64; i++ {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			t1 := h + ((bits.RotateLeft32(e, -6)) ^ (bits.RotateLeft32(e, -11)) ^ (bits.RotateLeft32(e, -25))) + ((e &amp; f) ^ (^e &amp; g)) + _K[i] + w[i]
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			t2 := ((bits.RotateLeft32(a, -2)) ^ (bits.RotateLeft32(a, -13)) ^ (bits.RotateLeft32(a, -22))) + ((a &amp; b) ^ (a &amp; c) ^ (b &amp; c))
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			h = g
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			g = f
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			f = e
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			e = d + t1
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			d = c
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			c = b
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			b = a
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			a = t1 + t2
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		h0 += a
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		h1 += b
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		h2 += c
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		h3 += d
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		h4 += e
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		h5 += f
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		h6 += g
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		h7 += h
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		p = p[chunk:]
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4], dig.h[5], dig.h[6], dig.h[7] = h0, h1, h2, h3, h4, h5, h6, h7
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
</pre><p><a href="sha256block.go?m=text">View as plain text</a></p>

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
