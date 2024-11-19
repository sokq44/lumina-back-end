<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/aes/block.go - Go Documentation Server</title>

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
<a href="block.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/aes">aes</a>/<span class="text-muted">block.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/aes">crypto/aes</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This Go implementation is derived in part from the reference</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// ANSI C implementation, which carries the following notice:</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//	rijndael-alg-fst.c</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//	@version 3.0 (December 2000)</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//	Optimised ANSI C code for the Rijndael cipher (now AES)</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//	@author Vincent Rijmen &lt;vincent.rijmen@esat.kuleuven.ac.be&gt;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//	@author Antoon Bosselaers &lt;antoon.bosselaers@esat.kuleuven.ac.be&gt;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//	@author Paulo Barreto &lt;paulo.barreto@terra.com.br&gt;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//	This code is hereby placed in the public domain.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//	THIS SOFTWARE IS PROVIDED BY THE AUTHORS &#39;&#39;AS IS&#39;&#39; AND ANY EXPRESS</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//	OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//	WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//	ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHORS OR CONTRIBUTORS BE</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//	LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//	CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//	SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//	BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//	WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//	OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//	EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// See FIPS 197 for specification, and see Daemen and Rijmen&#39;s Rijndael submission</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// for implementation details.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//	https://csrc.nist.gov/csrc/media/publications/fips/197/final/documents/fips-197.pdf</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//	https://csrc.nist.gov/archive/aes/rijndael/Rijndael-ammended.pdf</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>package aes
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>import (
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	&#34;encoding/binary&#34;
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// Encrypt one block from src into dst, using the expanded key xk.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>func encryptBlockGo(xk []uint32, dst, src []byte) {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	_ = src[15] <span class="comment">// early bounds check</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	s0 := binary.BigEndian.Uint32(src[0:4])
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	s1 := binary.BigEndian.Uint32(src[4:8])
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	s2 := binary.BigEndian.Uint32(src[8:12])
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	s3 := binary.BigEndian.Uint32(src[12:16])
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// First round just XORs input with key.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	s0 ^= xk[0]
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	s1 ^= xk[1]
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	s2 ^= xk[2]
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	s3 ^= xk[3]
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// Middle rounds shuffle using tables.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// Number of rounds is set by length of expanded key.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	nr := len(xk)/4 - 2 <span class="comment">// - 2: one above, one more below</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	k := 4
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	var t0, t1, t2, t3 uint32
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	for r := 0; r &lt; nr; r++ {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		t0 = xk[k+0] ^ te0[uint8(s0&gt;&gt;24)] ^ te1[uint8(s1&gt;&gt;16)] ^ te2[uint8(s2&gt;&gt;8)] ^ te3[uint8(s3)]
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		t1 = xk[k+1] ^ te0[uint8(s1&gt;&gt;24)] ^ te1[uint8(s2&gt;&gt;16)] ^ te2[uint8(s3&gt;&gt;8)] ^ te3[uint8(s0)]
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		t2 = xk[k+2] ^ te0[uint8(s2&gt;&gt;24)] ^ te1[uint8(s3&gt;&gt;16)] ^ te2[uint8(s0&gt;&gt;8)] ^ te3[uint8(s1)]
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		t3 = xk[k+3] ^ te0[uint8(s3&gt;&gt;24)] ^ te1[uint8(s0&gt;&gt;16)] ^ te2[uint8(s1&gt;&gt;8)] ^ te3[uint8(s2)]
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		k += 4
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		s0, s1, s2, s3 = t0, t1, t2, t3
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// Last round uses s-box directly and XORs to produce output.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	s0 = uint32(sbox0[t0&gt;&gt;24])&lt;&lt;24 | uint32(sbox0[t1&gt;&gt;16&amp;0xff])&lt;&lt;16 | uint32(sbox0[t2&gt;&gt;8&amp;0xff])&lt;&lt;8 | uint32(sbox0[t3&amp;0xff])
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	s1 = uint32(sbox0[t1&gt;&gt;24])&lt;&lt;24 | uint32(sbox0[t2&gt;&gt;16&amp;0xff])&lt;&lt;16 | uint32(sbox0[t3&gt;&gt;8&amp;0xff])&lt;&lt;8 | uint32(sbox0[t0&amp;0xff])
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	s2 = uint32(sbox0[t2&gt;&gt;24])&lt;&lt;24 | uint32(sbox0[t3&gt;&gt;16&amp;0xff])&lt;&lt;16 | uint32(sbox0[t0&gt;&gt;8&amp;0xff])&lt;&lt;8 | uint32(sbox0[t1&amp;0xff])
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	s3 = uint32(sbox0[t3&gt;&gt;24])&lt;&lt;24 | uint32(sbox0[t0&gt;&gt;16&amp;0xff])&lt;&lt;16 | uint32(sbox0[t1&gt;&gt;8&amp;0xff])&lt;&lt;8 | uint32(sbox0[t2&amp;0xff])
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	s0 ^= xk[k+0]
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	s1 ^= xk[k+1]
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	s2 ^= xk[k+2]
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	s3 ^= xk[k+3]
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	_ = dst[15] <span class="comment">// early bounds check</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	binary.BigEndian.PutUint32(dst[0:4], s0)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	binary.BigEndian.PutUint32(dst[4:8], s1)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	binary.BigEndian.PutUint32(dst[8:12], s2)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	binary.BigEndian.PutUint32(dst[12:16], s3)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// Decrypt one block from src into dst, using the expanded key xk.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>func decryptBlockGo(xk []uint32, dst, src []byte) {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	_ = src[15] <span class="comment">// early bounds check</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	s0 := binary.BigEndian.Uint32(src[0:4])
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	s1 := binary.BigEndian.Uint32(src[4:8])
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	s2 := binary.BigEndian.Uint32(src[8:12])
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	s3 := binary.BigEndian.Uint32(src[12:16])
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// First round just XORs input with key.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	s0 ^= xk[0]
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	s1 ^= xk[1]
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	s2 ^= xk[2]
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	s3 ^= xk[3]
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// Middle rounds shuffle using tables.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// Number of rounds is set by length of expanded key.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	nr := len(xk)/4 - 2 <span class="comment">// - 2: one above, one more below</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	k := 4
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	var t0, t1, t2, t3 uint32
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	for r := 0; r &lt; nr; r++ {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		t0 = xk[k+0] ^ td0[uint8(s0&gt;&gt;24)] ^ td1[uint8(s3&gt;&gt;16)] ^ td2[uint8(s2&gt;&gt;8)] ^ td3[uint8(s1)]
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		t1 = xk[k+1] ^ td0[uint8(s1&gt;&gt;24)] ^ td1[uint8(s0&gt;&gt;16)] ^ td2[uint8(s3&gt;&gt;8)] ^ td3[uint8(s2)]
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		t2 = xk[k+2] ^ td0[uint8(s2&gt;&gt;24)] ^ td1[uint8(s1&gt;&gt;16)] ^ td2[uint8(s0&gt;&gt;8)] ^ td3[uint8(s3)]
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		t3 = xk[k+3] ^ td0[uint8(s3&gt;&gt;24)] ^ td1[uint8(s2&gt;&gt;16)] ^ td2[uint8(s1&gt;&gt;8)] ^ td3[uint8(s0)]
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		k += 4
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		s0, s1, s2, s3 = t0, t1, t2, t3
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// Last round uses s-box directly and XORs to produce output.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	s0 = uint32(sbox1[t0&gt;&gt;24])&lt;&lt;24 | uint32(sbox1[t3&gt;&gt;16&amp;0xff])&lt;&lt;16 | uint32(sbox1[t2&gt;&gt;8&amp;0xff])&lt;&lt;8 | uint32(sbox1[t1&amp;0xff])
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	s1 = uint32(sbox1[t1&gt;&gt;24])&lt;&lt;24 | uint32(sbox1[t0&gt;&gt;16&amp;0xff])&lt;&lt;16 | uint32(sbox1[t3&gt;&gt;8&amp;0xff])&lt;&lt;8 | uint32(sbox1[t2&amp;0xff])
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	s2 = uint32(sbox1[t2&gt;&gt;24])&lt;&lt;24 | uint32(sbox1[t1&gt;&gt;16&amp;0xff])&lt;&lt;16 | uint32(sbox1[t0&gt;&gt;8&amp;0xff])&lt;&lt;8 | uint32(sbox1[t3&amp;0xff])
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	s3 = uint32(sbox1[t3&gt;&gt;24])&lt;&lt;24 | uint32(sbox1[t2&gt;&gt;16&amp;0xff])&lt;&lt;16 | uint32(sbox1[t1&gt;&gt;8&amp;0xff])&lt;&lt;8 | uint32(sbox1[t0&amp;0xff])
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	s0 ^= xk[k+0]
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	s1 ^= xk[k+1]
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	s2 ^= xk[k+2]
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	s3 ^= xk[k+3]
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	_ = dst[15] <span class="comment">// early bounds check</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	binary.BigEndian.PutUint32(dst[0:4], s0)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	binary.BigEndian.PutUint32(dst[4:8], s1)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	binary.BigEndian.PutUint32(dst[8:12], s2)
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	binary.BigEndian.PutUint32(dst[12:16], s3)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">// Apply sbox0 to each byte in w.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>func subw(w uint32) uint32 {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	return uint32(sbox0[w&gt;&gt;24])&lt;&lt;24 |
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		uint32(sbox0[w&gt;&gt;16&amp;0xff])&lt;&lt;16 |
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		uint32(sbox0[w&gt;&gt;8&amp;0xff])&lt;&lt;8 |
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		uint32(sbox0[w&amp;0xff])
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// Rotate</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>func rotw(w uint32) uint32 { return w&lt;&lt;8 | w&gt;&gt;24 }
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// Key expansion algorithm. See FIPS-197, Figure 11.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// Their rcon[i] is our powx[i-1] &lt;&lt; 24.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>func expandKeyGo(key []byte, enc, dec []uint32) {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// Encryption key setup.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	var i int
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	nk := len(key) / 4
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	for i = 0; i &lt; nk; i++ {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		enc[i] = binary.BigEndian.Uint32(key[4*i:])
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	for ; i &lt; len(enc); i++ {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		t := enc[i-1]
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		if i%nk == 0 {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			t = subw(rotw(t)) ^ (uint32(powx[i/nk-1]) &lt;&lt; 24)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		} else if nk &gt; 6 &amp;&amp; i%nk == 4 {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			t = subw(t)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		enc[i] = enc[i-nk] ^ t
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// Derive decryption key from encryption key.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// Reverse the 4-word round key sets from enc to produce dec.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// All sets but the first and last get the MixColumn transform applied.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	if dec == nil {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		return
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	n := len(enc)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i += 4 {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		ei := n - i - 4
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		for j := 0; j &lt; 4; j++ {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			x := enc[ei+j]
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			if i &gt; 0 &amp;&amp; i+4 &lt; n {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>				x = td0[sbox0[x&gt;&gt;24]] ^ td1[sbox0[x&gt;&gt;16&amp;0xff]] ^ td2[sbox0[x&gt;&gt;8&amp;0xff]] ^ td3[sbox0[x&amp;0xff]]
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			dec[i+j] = x
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
</pre><p><a href="block.go?m=text">View as plain text</a></p>

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
