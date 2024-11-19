<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/filippo.io/edwards25519/field/fe_generic.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../../index.html">GoDoc</a></div>
<a href="fe_generic.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/filippo.io">filippo.io</a>/<a href="http://localhost:8080/src/filippo.io/edwards25519">edwards25519</a>/<a href="http://localhost:8080/src/filippo.io/edwards25519/field">field</a>/<span class="text-muted">fe_generic.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/filippo.io/edwards25519/field">filippo.io/edwards25519/field</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright (c) 2017 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package field
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;math/bits&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// uint128 holds a 128-bit number as two 64-bit limbs, for use with the</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// bits.Mul64 and bits.Add64 intrinsics.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>type uint128 struct {
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	lo, hi uint64
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>}
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// mul64 returns a * b.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>func mul64(a, b uint64) uint128 {
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	hi, lo := bits.Mul64(a, b)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	return uint128{lo, hi}
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>}
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// addMul64 returns v + a * b.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>func addMul64(v uint128, a, b uint64) uint128 {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	hi, lo := bits.Mul64(a, b)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	lo, c := bits.Add64(lo, v.lo, 0)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	hi, _ = bits.Add64(hi, v.hi, c)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	return uint128{lo, hi}
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// shiftRightBy51 returns a &gt;&gt; 51. a is assumed to be at most 115 bits.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>func shiftRightBy51(a uint128) uint64 {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	return (a.hi &lt;&lt; (64 - 51)) | (a.lo &gt;&gt; 51)
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>func feMulGeneric(v, a, b *Element) {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	a0 := a.l0
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	a1 := a.l1
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	a2 := a.l2
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	a3 := a.l3
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	a4 := a.l4
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	b0 := b.l0
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	b1 := b.l1
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	b2 := b.l2
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	b3 := b.l3
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	b4 := b.l4
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// Limb multiplication works like pen-and-paper columnar multiplication, but</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// with 51-bit limbs instead of digits.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">//                          a4   a3   a2   a1   a0  x</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">//                          b4   b3   b2   b1   b0  =</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">//                         ------------------------</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">//                        a4b0 a3b0 a2b0 a1b0 a0b0  +</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">//                   a4b1 a3b1 a2b1 a1b1 a0b1       +</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">//              a4b2 a3b2 a2b2 a1b2 a0b2            +</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">//         a4b3 a3b3 a2b3 a1b3 a0b3                 +</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">//    a4b4 a3b4 a2b4 a1b4 a0b4                      =</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">//   ----------------------------------------------</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">//      r8   r7   r6   r5   r4   r3   r2   r1   r0</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// We can then use the reduction identity (a * 2²⁵⁵ + b = a * 19 + b) to</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// reduce the limbs that would overflow 255 bits. r5 * 2²⁵⁵ becomes 19 * r5,</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// r6 * 2³⁰⁶ becomes 19 * r6 * 2⁵¹, etc.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// Reduction can be carried out simultaneously to multiplication. For</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// example, we do not compute r5: whenever the result of a multiplication</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// belongs to r5, like a1b4, we multiply it by 19 and add the result to r0.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">//            a4b0    a3b0    a2b0    a1b0    a0b0  +</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">//            a3b1    a2b1    a1b1    a0b1 19×a4b1  +</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">//            a2b2    a1b2    a0b2 19×a4b2 19×a3b2  +</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">//            a1b3    a0b3 19×a4b3 19×a3b3 19×a2b3  +</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">//            a0b4 19×a4b4 19×a3b4 19×a2b4 19×a1b4  =</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">//           --------------------------------------</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">//              r4      r3      r2      r1      r0</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// Finally we add up the columns into wide, overlapping limbs.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	a1_19 := a1 * 19
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	a2_19 := a2 * 19
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	a3_19 := a3 * 19
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	a4_19 := a4 * 19
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// r0 = a0×b0 + 19×(a1×b4 + a2×b3 + a3×b2 + a4×b1)</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	r0 := mul64(a0, b0)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	r0 = addMul64(r0, a1_19, b4)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	r0 = addMul64(r0, a2_19, b3)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	r0 = addMul64(r0, a3_19, b2)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	r0 = addMul64(r0, a4_19, b1)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// r1 = a0×b1 + a1×b0 + 19×(a2×b4 + a3×b3 + a4×b2)</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	r1 := mul64(a0, b1)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	r1 = addMul64(r1, a1, b0)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	r1 = addMul64(r1, a2_19, b4)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	r1 = addMul64(r1, a3_19, b3)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	r1 = addMul64(r1, a4_19, b2)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// r2 = a0×b2 + a1×b1 + a2×b0 + 19×(a3×b4 + a4×b3)</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	r2 := mul64(a0, b2)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	r2 = addMul64(r2, a1, b1)
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	r2 = addMul64(r2, a2, b0)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	r2 = addMul64(r2, a3_19, b4)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	r2 = addMul64(r2, a4_19, b3)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// r3 = a0×b3 + a1×b2 + a2×b1 + a3×b0 + 19×a4×b4</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	r3 := mul64(a0, b3)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	r3 = addMul64(r3, a1, b2)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	r3 = addMul64(r3, a2, b1)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	r3 = addMul64(r3, a3, b0)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	r3 = addMul64(r3, a4_19, b4)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// r4 = a0×b4 + a1×b3 + a2×b2 + a3×b1 + a4×b0</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	r4 := mul64(a0, b4)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	r4 = addMul64(r4, a1, b3)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	r4 = addMul64(r4, a2, b2)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	r4 = addMul64(r4, a3, b1)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	r4 = addMul64(r4, a4, b0)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// After the multiplication, we need to reduce (carry) the five coefficients</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// to obtain a result with limbs that are at most slightly larger than 2⁵¹,</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// to respect the Element invariant.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// Overall, the reduction works the same as carryPropagate, except with</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// wider inputs: we take the carry for each coefficient by shifting it right</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// by 51, and add it to the limb above it. The top carry is multiplied by 19</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// according to the reduction identity and added to the lowest limb.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// The largest coefficient (r0) will be at most 111 bits, which guarantees</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// that all carries are at most 111 - 51 = 60 bits, which fits in a uint64.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">//     r0 = a0×b0 + 19×(a1×b4 + a2×b3 + a3×b2 + a4×b1)</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">//     r0 &lt; 2⁵²×2⁵² + 19×(2⁵²×2⁵² + 2⁵²×2⁵² + 2⁵²×2⁵² + 2⁵²×2⁵²)</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">//     r0 &lt; (1 + 19 × 4) × 2⁵² × 2⁵²</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">//     r0 &lt; 2⁷ × 2⁵² × 2⁵²</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">//     r0 &lt; 2¹¹¹</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// Moreover, the top coefficient (r4) is at most 107 bits, so c4 is at most</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// 56 bits, and c4 * 19 is at most 61 bits, which again fits in a uint64 and</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// allows us to easily apply the reduction identity.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">//     r4 = a0×b4 + a1×b3 + a2×b2 + a3×b1 + a4×b0</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">//     r4 &lt; 5 × 2⁵² × 2⁵²</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">//     r4 &lt; 2¹⁰⁷</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	c0 := shiftRightBy51(r0)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	c1 := shiftRightBy51(r1)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	c2 := shiftRightBy51(r2)
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	c3 := shiftRightBy51(r3)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	c4 := shiftRightBy51(r4)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	rr0 := r0.lo&amp;maskLow51Bits + c4*19
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	rr1 := r1.lo&amp;maskLow51Bits + c0
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	rr2 := r2.lo&amp;maskLow51Bits + c1
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	rr3 := r3.lo&amp;maskLow51Bits + c2
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	rr4 := r4.lo&amp;maskLow51Bits + c3
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// Now all coefficients fit into 64-bit registers but are still too large to</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// be passed around as an Element. We therefore do one last carry chain,</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// where the carries will be small enough to fit in the wiggle room above 2⁵¹.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	*v = Element{rr0, rr1, rr2, rr3, rr4}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	v.carryPropagate()
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>func feSquareGeneric(v, a *Element) {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	l0 := a.l0
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	l1 := a.l1
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	l2 := a.l2
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	l3 := a.l3
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	l4 := a.l4
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// Squaring works precisely like multiplication above, but thanks to its</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// symmetry we get to group a few terms together.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">//                          l4   l3   l2   l1   l0  x</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">//                          l4   l3   l2   l1   l0  =</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">//                         ------------------------</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">//                        l4l0 l3l0 l2l0 l1l0 l0l0  +</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">//                   l4l1 l3l1 l2l1 l1l1 l0l1       +</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">//              l4l2 l3l2 l2l2 l1l2 l0l2            +</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">//         l4l3 l3l3 l2l3 l1l3 l0l3                 +</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">//    l4l4 l3l4 l2l4 l1l4 l0l4                      =</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">//   ----------------------------------------------</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">//      r8   r7   r6   r5   r4   r3   r2   r1   r0</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">//            l4l0    l3l0    l2l0    l1l0    l0l0  +</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">//            l3l1    l2l1    l1l1    l0l1 19×l4l1  +</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">//            l2l2    l1l2    l0l2 19×l4l2 19×l3l2  +</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">//            l1l3    l0l3 19×l4l3 19×l3l3 19×l2l3  +</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">//            l0l4 19×l4l4 19×l3l4 19×l2l4 19×l1l4  =</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">//           --------------------------------------</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">//              r4      r3      r2      r1      r0</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// With precomputed 2×, 19×, and 2×19× terms, we can compute each limb with</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// only three Mul64 and four Add64, instead of five and eight.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	l0_2 := l0 * 2
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	l1_2 := l1 * 2
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	l1_38 := l1 * 38
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	l2_38 := l2 * 38
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	l3_38 := l3 * 38
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	l3_19 := l3 * 19
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	l4_19 := l4 * 19
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// r0 = l0×l0 + 19×(l1×l4 + l2×l3 + l3×l2 + l4×l1) = l0×l0 + 19×2×(l1×l4 + l2×l3)</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	r0 := mul64(l0, l0)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	r0 = addMul64(r0, l1_38, l4)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	r0 = addMul64(r0, l2_38, l3)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// r1 = l0×l1 + l1×l0 + 19×(l2×l4 + l3×l3 + l4×l2) = 2×l0×l1 + 19×2×l2×l4 + 19×l3×l3</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	r1 := mul64(l0_2, l1)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	r1 = addMul64(r1, l2_38, l4)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	r1 = addMul64(r1, l3_19, l3)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	<span class="comment">// r2 = l0×l2 + l1×l1 + l2×l0 + 19×(l3×l4 + l4×l3) = 2×l0×l2 + l1×l1 + 19×2×l3×l4</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	r2 := mul64(l0_2, l2)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	r2 = addMul64(r2, l1, l1)
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	r2 = addMul64(r2, l3_38, l4)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">// r3 = l0×l3 + l1×l2 + l2×l1 + l3×l0 + 19×l4×l4 = 2×l0×l3 + 2×l1×l2 + 19×l4×l4</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	r3 := mul64(l0_2, l3)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	r3 = addMul64(r3, l1_2, l2)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	r3 = addMul64(r3, l4_19, l4)
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	<span class="comment">// r4 = l0×l4 + l1×l3 + l2×l2 + l3×l1 + l4×l0 = 2×l0×l4 + 2×l1×l3 + l2×l2</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	r4 := mul64(l0_2, l4)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	r4 = addMul64(r4, l1_2, l3)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	r4 = addMul64(r4, l2, l2)
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	c0 := shiftRightBy51(r0)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	c1 := shiftRightBy51(r1)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	c2 := shiftRightBy51(r2)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	c3 := shiftRightBy51(r3)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	c4 := shiftRightBy51(r4)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	rr0 := r0.lo&amp;maskLow51Bits + c4*19
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	rr1 := r1.lo&amp;maskLow51Bits + c0
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	rr2 := r2.lo&amp;maskLow51Bits + c1
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	rr3 := r3.lo&amp;maskLow51Bits + c2
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	rr4 := r4.lo&amp;maskLow51Bits + c3
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	*v = Element{rr0, rr1, rr2, rr3, rr4}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	v.carryPropagate()
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// carryPropagateGeneric brings the limbs below 52 bits by applying the reduction</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">// identity (a * 2²⁵⁵ + b = a * 19 + b) to the l4 carry.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>func (v *Element) carryPropagateGeneric() *Element {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	c0 := v.l0 &gt;&gt; 51
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	c1 := v.l1 &gt;&gt; 51
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	c2 := v.l2 &gt;&gt; 51
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	c3 := v.l3 &gt;&gt; 51
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	c4 := v.l4 &gt;&gt; 51
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">// c4 is at most 64 - 51 = 13 bits, so c4*19 is at most 18 bits, and</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	<span class="comment">// the final l0 will be at most 52 bits. Similarly for the rest.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	v.l0 = v.l0&amp;maskLow51Bits + c4*19
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	v.l1 = v.l1&amp;maskLow51Bits + c0
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	v.l2 = v.l2&amp;maskLow51Bits + c1
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	v.l3 = v.l3&amp;maskLow51Bits + c2
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	v.l4 = v.l4&amp;maskLow51Bits + c3
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	return v
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
</pre><p><a href="fe_generic.go?m=text">View as plain text</a></p>

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
