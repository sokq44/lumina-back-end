<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/des/block.go - Go Documentation Server</title>

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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/des">des</a>/<span class="text-muted">block.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/des">crypto/des</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package des
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;encoding/binary&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>func cryptBlock(subkeys []uint64, dst, src []byte, decrypt bool) {
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	b := binary.BigEndian.Uint64(src)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	b = permuteInitialBlock(b)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	left, right := uint32(b&gt;&gt;32), uint32(b)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	left = (left &lt;&lt; 1) | (left &gt;&gt; 31)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	right = (right &lt;&lt; 1) | (right &gt;&gt; 31)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	if decrypt {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>		for i := 0; i &lt; 8; i++ {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>			left, right = feistel(left, right, subkeys[15-2*i], subkeys[15-(2*i+1)])
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>		}
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	} else {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		for i := 0; i &lt; 8; i++ {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>			left, right = feistel(left, right, subkeys[2*i], subkeys[2*i+1])
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	left = (left &lt;&lt; 31) | (left &gt;&gt; 1)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	right = (right &lt;&lt; 31) | (right &gt;&gt; 1)
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// switch left &amp; right and perform final permutation</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	preOutput := (uint64(right) &lt;&lt; 32) | uint64(left)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	binary.BigEndian.PutUint64(dst, permuteFinalBlock(preOutput))
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// DES Feistel function. feistelBox must be initialized via</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// feistelBoxOnce.Do(initFeistelBox) first.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>func feistel(l, r uint32, k0, k1 uint64) (lout, rout uint32) {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	var t uint32
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	t = r ^ uint32(k0&gt;&gt;32)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	l ^= feistelBox[7][t&amp;0x3f] ^
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		feistelBox[5][(t&gt;&gt;8)&amp;0x3f] ^
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		feistelBox[3][(t&gt;&gt;16)&amp;0x3f] ^
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		feistelBox[1][(t&gt;&gt;24)&amp;0x3f]
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	t = ((r &lt;&lt; 28) | (r &gt;&gt; 4)) ^ uint32(k0)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	l ^= feistelBox[6][(t)&amp;0x3f] ^
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		feistelBox[4][(t&gt;&gt;8)&amp;0x3f] ^
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		feistelBox[2][(t&gt;&gt;16)&amp;0x3f] ^
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		feistelBox[0][(t&gt;&gt;24)&amp;0x3f]
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	t = l ^ uint32(k1&gt;&gt;32)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	r ^= feistelBox[7][t&amp;0x3f] ^
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		feistelBox[5][(t&gt;&gt;8)&amp;0x3f] ^
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		feistelBox[3][(t&gt;&gt;16)&amp;0x3f] ^
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		feistelBox[1][(t&gt;&gt;24)&amp;0x3f]
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	t = ((l &lt;&lt; 28) | (l &gt;&gt; 4)) ^ uint32(k1)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	r ^= feistelBox[6][(t)&amp;0x3f] ^
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		feistelBox[4][(t&gt;&gt;8)&amp;0x3f] ^
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		feistelBox[2][(t&gt;&gt;16)&amp;0x3f] ^
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		feistelBox[0][(t&gt;&gt;24)&amp;0x3f]
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	return l, r
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// feistelBox[s][16*i+j] contains the output of permutationFunction</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// for sBoxes[s][i][j] &lt;&lt; 4*(7-s)</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>var feistelBox [8][64]uint32
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>var feistelBoxOnce sync.Once
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// general purpose function to perform DES block permutations.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>func permuteBlock(src uint64, permutation []uint8) (block uint64) {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	for position, n := range permutation {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		bit := (src &gt;&gt; n) &amp; 1
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		block |= bit &lt;&lt; uint((len(permutation)-1)-position)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	return
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>func initFeistelBox() {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	for s := range sBoxes {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		for i := 0; i &lt; 4; i++ {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			for j := 0; j &lt; 16; j++ {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>				f := uint64(sBoxes[s][i][j]) &lt;&lt; (4 * (7 - uint(s)))
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>				f = permuteBlock(f, permutationFunction[:])
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>				<span class="comment">// Row is determined by the 1st and 6th bit.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>				<span class="comment">// Column is the middle four bits.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>				row := uint8(((i &amp; 2) &lt;&lt; 4) | i&amp;1)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>				col := uint8(j &lt;&lt; 1)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>				t := row | col
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>				<span class="comment">// The rotation was performed in the feistel rounds, being factored out and now mixed into the feistelBox.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>				f = (f &lt;&lt; 1) | (f &gt;&gt; 31)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>				feistelBox[s][t] = uint32(f)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// permuteInitialBlock is equivalent to the permutation defined</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// by initialPermutation.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>func permuteInitialBlock(block uint64) uint64 {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// block = b7 b6 b5 b4 b3 b2 b1 b0 (8 bytes)</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	b1 := block &gt;&gt; 48
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	b2 := block &lt;&lt; 48
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	block ^= b1 ^ b2 ^ b1&lt;&lt;48 ^ b2&gt;&gt;48
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// block = b1 b0 b5 b4 b3 b2 b7 b6</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	b1 = block &gt;&gt; 32 &amp; 0xff00ff
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	b2 = (block &amp; 0xff00ff00)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	block ^= b1&lt;&lt;32 ^ b2 ^ b1&lt;&lt;8 ^ b2&lt;&lt;24 <span class="comment">// exchange b0 b4 with b3 b7</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// block is now b1 b3 b5 b7 b0 b2 b4 b6, the permutation:</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">//                  ...  8</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">//                  ... 24</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">//                  ... 40</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">//                  ... 56</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">//  7  6  5  4  3  2  1  0</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// 23 22 21 20 19 18 17 16</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">//                  ... 32</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">//                  ... 48</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// exchange 4,5,6,7 with 32,33,34,35 etc.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	b1 = block &amp; 0x0f0f00000f0f0000
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	b2 = block &amp; 0x0000f0f00000f0f0
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	block ^= b1 ^ b2 ^ b1&gt;&gt;12 ^ b2&lt;&lt;12
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// block is the permutation:</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">//   [+8]         [+40]</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">//  7  6  5  4</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// 23 22 21 20</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">//  3  2  1  0</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// 19 18 17 16    [+32]</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// exchange 0,1,4,5 with 18,19,22,23</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	b1 = block &amp; 0x3300330033003300
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	b2 = block &amp; 0x00cc00cc00cc00cc
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	block ^= b1 ^ b2 ^ b1&gt;&gt;6 ^ b2&lt;&lt;6
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// block is the permutation:</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// 15 14</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// 13 12</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// 11 10</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">//  9  8</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">//  7  6</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">//  5  4</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">//  3  2</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">//  1  0 [+16] [+32] [+64]</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// exchange 0,2,4,6 with 9,11,13,15:</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	b1 = block &amp; 0xaaaaaaaa55555555
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	block ^= b1 ^ b1&gt;&gt;33 ^ b1&lt;&lt;33
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// block is the permutation:</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// 6 14 22 30 38 46 54 62</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// 4 12 20 28 36 44 52 60</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// 2 10 18 26 34 42 50 58</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// 0  8 16 24 32 40 48 56</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// 7 15 23 31 39 47 55 63</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// 5 13 21 29 37 45 53 61</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// 3 11 19 27 35 43 51 59</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// 1  9 17 25 33 41 49 57</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	return block
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// permuteFinalBlock is equivalent to the permutation defined</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// by finalPermutation.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>func permuteFinalBlock(block uint64) uint64 {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// Perform the same bit exchanges as permuteInitialBlock</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// but in reverse order.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	b1 := block &amp; 0xaaaaaaaa55555555
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	block ^= b1 ^ b1&gt;&gt;33 ^ b1&lt;&lt;33
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	b1 = block &amp; 0x3300330033003300
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	b2 := block &amp; 0x00cc00cc00cc00cc
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	block ^= b1 ^ b2 ^ b1&gt;&gt;6 ^ b2&lt;&lt;6
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	b1 = block &amp; 0x0f0f00000f0f0000
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	b2 = block &amp; 0x0000f0f00000f0f0
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	block ^= b1 ^ b2 ^ b1&gt;&gt;12 ^ b2&lt;&lt;12
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	b1 = block &gt;&gt; 32 &amp; 0xff00ff
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	b2 = (block &amp; 0xff00ff00)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	block ^= b1&lt;&lt;32 ^ b2 ^ b1&lt;&lt;8 ^ b2&lt;&lt;24
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	b1 = block &gt;&gt; 48
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	b2 = block &lt;&lt; 48
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	block ^= b1 ^ b2 ^ b1&lt;&lt;48 ^ b2&gt;&gt;48
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	return block
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// creates 16 28-bit blocks rotated according</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">// to the rotation schedule.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>func ksRotate(in uint32) (out []uint32) {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	out = make([]uint32, 16)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	last := in
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	for i := 0; i &lt; 16; i++ {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		<span class="comment">// 28-bit circular left shift</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		left := (last &lt;&lt; (4 + ksRotations[i])) &gt;&gt; 4
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		right := (last &lt;&lt; 4) &gt;&gt; (32 - ksRotations[i])
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		out[i] = left | right
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		last = out[i]
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	return
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// creates 16 56-bit subkeys from the original key.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>func (c *desCipher) generateSubkeys(keyBytes []byte) {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	feistelBoxOnce.Do(initFeistelBox)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// apply PC1 permutation to key</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	key := binary.BigEndian.Uint64(keyBytes)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	permutedKey := permuteBlock(key, permutedChoice1[:])
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">// rotate halves of permuted key according to the rotation schedule</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	leftRotations := ksRotate(uint32(permutedKey &gt;&gt; 28))
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	rightRotations := ksRotate(uint32(permutedKey&lt;&lt;4) &gt;&gt; 4)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// generate subkeys</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	for i := 0; i &lt; 16; i++ {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		<span class="comment">// combine halves to form 56-bit input to PC2</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		pc2Input := uint64(leftRotations[i])&lt;&lt;28 | uint64(rightRotations[i])
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		<span class="comment">// apply PC2 permutation to 7 byte input</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		c.subkeys[i] = unpack(permuteBlock(pc2Input, permutedChoice2[:]))
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">// Expand 48-bit input to 64-bit, with each 6-bit block padded by extra two bits at the top.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">// By doing so, we can have the input blocks (four bits each), and the key blocks (six bits each) well-aligned without</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">// extra shifts/rotations for alignments.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>func unpack(x uint64) uint64 {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	return ((x&gt;&gt;(6*1))&amp;0xff)&lt;&lt;(8*0) |
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		((x&gt;&gt;(6*3))&amp;0xff)&lt;&lt;(8*1) |
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		((x&gt;&gt;(6*5))&amp;0xff)&lt;&lt;(8*2) |
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		((x&gt;&gt;(6*7))&amp;0xff)&lt;&lt;(8*3) |
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		((x&gt;&gt;(6*0))&amp;0xff)&lt;&lt;(8*4) |
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		((x&gt;&gt;(6*2))&amp;0xff)&lt;&lt;(8*5) |
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		((x&gt;&gt;(6*4))&amp;0xff)&lt;&lt;(8*6) |
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		((x&gt;&gt;(6*6))&amp;0xff)&lt;&lt;(8*7)
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
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
