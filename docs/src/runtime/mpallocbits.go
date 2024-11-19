<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mpallocbits.go - Go Documentation Server</title>

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
<a href="mpallocbits.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mpallocbits.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2019 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>)
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// pageBits is a bitmap representing one bit per page in a palloc chunk.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>type pageBits [pallocChunkPages / 64]uint64
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// get returns the value of the i&#39;th bit in the bitmap.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>func (b *pageBits) get(i uint) uint {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	return uint((b[i/64] &gt;&gt; (i % 64)) &amp; 1)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>}
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// block64 returns the 64-bit aligned block of bits containing the i&#39;th bit.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>func (b *pageBits) block64(i uint) uint64 {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	return b[i/64]
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// set sets bit i of pageBits.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>func (b *pageBits) set(i uint) {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	b[i/64] |= 1 &lt;&lt; (i % 64)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// setRange sets bits in the range [i, i+n).</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>func (b *pageBits) setRange(i, n uint) {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	_ = b[i/64]
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	if n == 1 {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		<span class="comment">// Fast path for the n == 1 case.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		b.set(i)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		return
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// Set bits [i, j].</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	j := i + n - 1
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	if i/64 == j/64 {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		b[i/64] |= ((uint64(1) &lt;&lt; n) - 1) &lt;&lt; (i % 64)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		return
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	}
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	_ = b[j/64]
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// Set leading bits.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	b[i/64] |= ^uint64(0) &lt;&lt; (i % 64)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	for k := i/64 + 1; k &lt; j/64; k++ {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		b[k] = ^uint64(0)
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// Set trailing bits.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	b[j/64] |= (uint64(1) &lt;&lt; (j%64 + 1)) - 1
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// setAll sets all the bits of b.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func (b *pageBits) setAll() {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	for i := range b {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		b[i] = ^uint64(0)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// setBlock64 sets the 64-bit aligned block of bits containing the i&#39;th bit that</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// are set in v.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>func (b *pageBits) setBlock64(i uint, v uint64) {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	b[i/64] |= v
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// clear clears bit i of pageBits.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>func (b *pageBits) clear(i uint) {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	b[i/64] &amp;^= 1 &lt;&lt; (i % 64)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// clearRange clears bits in the range [i, i+n).</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func (b *pageBits) clearRange(i, n uint) {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	_ = b[i/64]
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	if n == 1 {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		<span class="comment">// Fast path for the n == 1 case.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		b.clear(i)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		return
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// Clear bits [i, j].</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	j := i + n - 1
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	if i/64 == j/64 {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		b[i/64] &amp;^= ((uint64(1) &lt;&lt; n) - 1) &lt;&lt; (i % 64)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		return
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	_ = b[j/64]
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// Clear leading bits.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	b[i/64] &amp;^= ^uint64(0) &lt;&lt; (i % 64)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	for k := i/64 + 1; k &lt; j/64; k++ {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		b[k] = 0
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// Clear trailing bits.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	b[j/64] &amp;^= (uint64(1) &lt;&lt; (j%64 + 1)) - 1
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// clearAll frees all the bits of b.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func (b *pageBits) clearAll() {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	for i := range b {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		b[i] = 0
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// clearBlock64 clears the 64-bit aligned block of bits containing the i&#39;th bit that</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// are set in v.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>func (b *pageBits) clearBlock64(i uint, v uint64) {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	b[i/64] &amp;^= v
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// popcntRange counts the number of set bits in the</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// range [i, i+n).</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>func (b *pageBits) popcntRange(i, n uint) (s uint) {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	if n == 1 {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		return uint((b[i/64] &gt;&gt; (i % 64)) &amp; 1)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	_ = b[i/64]
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	j := i + n - 1
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	if i/64 == j/64 {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		return uint(sys.OnesCount64((b[i/64] &gt;&gt; (i % 64)) &amp; ((1 &lt;&lt; n) - 1)))
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	_ = b[j/64]
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	s += uint(sys.OnesCount64(b[i/64] &gt;&gt; (i % 64)))
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	for k := i/64 + 1; k &lt; j/64; k++ {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		s += uint(sys.OnesCount64(b[k]))
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	s += uint(sys.OnesCount64(b[j/64] &amp; ((1 &lt;&lt; (j%64 + 1)) - 1)))
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	return
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// pallocBits is a bitmap that tracks page allocations for at most one</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// palloc chunk.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">// The precise representation is an implementation detail, but for the</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// sake of documentation, 0s are free pages and 1s are allocated pages.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>type pallocBits pageBits
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">// summarize returns a packed summary of the bitmap in pallocBits.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>func (b *pallocBits) summarize() pallocSum {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	var start, most, cur uint
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	const notSetYet = ^uint(0) <span class="comment">// sentinel for start value</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	start = notSetYet
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	for i := 0; i &lt; len(b); i++ {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		x := b[i]
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		if x == 0 {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			cur += 64
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			continue
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		t := uint(sys.TrailingZeros64(x))
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		l := uint(sys.LeadingZeros64(x))
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		<span class="comment">// Finish any region spanning the uint64s</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		cur += t
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		if start == notSetYet {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			start = cur
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		most = max(most, cur)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		<span class="comment">// Final region that might span to next uint64</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		cur = l
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	if start == notSetYet {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		<span class="comment">// Made it all the way through without finding a single 1 bit.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		const n = uint(64 * len(b))
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		return packPallocSum(n, n, n)
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	most = max(most, cur)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	if most &gt;= 64-2 {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		<span class="comment">// There is no way an internal run of zeros could beat max.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		return packPallocSum(start, most, cur)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// Now look inside each uint64 for runs of zeros.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// All uint64s must be nonzero, or we would have aborted above.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>outer:
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	for i := 0; i &lt; len(b); i++ {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		x := b[i]
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		<span class="comment">// Look inside this uint64. We have a pattern like</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		<span class="comment">// 000000 1xxxxx1 000000</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		<span class="comment">// We need to look inside the 1xxxxx1 for any contiguous</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		<span class="comment">// region of zeros.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		<span class="comment">// We already know the trailing zeros are no larger than max. Remove them.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		x &gt;&gt;= sys.TrailingZeros64(x) &amp; 63
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		if x&amp;(x+1) == 0 { <span class="comment">// no more zeros (except at the top).</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			continue
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		<span class="comment">// Strategy: shrink all runs of zeros by max. If any runs of zero</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		<span class="comment">// remain, then we&#39;ve identified a larger maximum zero run.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		p := most    <span class="comment">// number of zeros we still need to shrink by.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		k := uint(1) <span class="comment">// current minimum length of runs of ones in x.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		for {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			<span class="comment">// Shrink all runs of zeros by p places (except the top zeros).</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			for p &gt; 0 {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>				if p &lt;= k {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>					<span class="comment">// Shift p ones down into the top of each run of zeros.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>					x |= x &gt;&gt; (p &amp; 63)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>					if x&amp;(x+1) == 0 { <span class="comment">// no more zeros (except at the top).</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>						continue outer
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>					}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>					break
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>				}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>				<span class="comment">// Shift k ones down into the top of each run of zeros.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>				x |= x &gt;&gt; (k &amp; 63)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>				if x&amp;(x+1) == 0 { <span class="comment">// no more zeros (except at the top).</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>					continue outer
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>				}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>				p -= k
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>				<span class="comment">// We&#39;ve just doubled the minimum length of 1-runs.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>				<span class="comment">// This allows us to shift farther in the next iteration.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>				k *= 2
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			<span class="comment">// The length of the lowest-order zero run is an increment to our maximum.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			j := uint(sys.TrailingZeros64(^x)) <span class="comment">// count contiguous trailing ones</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			x &gt;&gt;= j &amp; 63                       <span class="comment">// remove trailing ones</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			j = uint(sys.TrailingZeros64(x))   <span class="comment">// count contiguous trailing zeros</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			x &gt;&gt;= j &amp; 63                       <span class="comment">// remove zeros</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			most += j                          <span class="comment">// we have a new maximum!</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			if x&amp;(x+1) == 0 {                  <span class="comment">// no more zeros (except at the top).</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>				continue outer
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			p = j <span class="comment">// remove j more zeros from each zero run.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	return packPallocSum(start, most, cur)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// find searches for npages contiguous free pages in pallocBits and returns</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// the index where that run starts, as well as the index of the first free page</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// it found in the search. searchIdx represents the first known free page and</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// where to begin the next search from.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// If find fails to find any free space, it returns an index of ^uint(0) and</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">// the new searchIdx should be ignored.</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">// Note that if npages == 1, the two returned values will always be identical.</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>func (b *pallocBits) find(npages uintptr, searchIdx uint) (uint, uint) {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	if npages == 1 {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		addr := b.find1(searchIdx)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		return addr, addr
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	} else if npages &lt;= 64 {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		return b.findSmallN(npages, searchIdx)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	return b.findLargeN(npages, searchIdx)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">// find1 is a helper for find which searches for a single free page</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// in the pallocBits and returns the index.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">// See find for an explanation of the searchIdx parameter.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>func (b *pallocBits) find1(searchIdx uint) uint {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	_ = b[0] <span class="comment">// lift nil check out of loop</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	for i := searchIdx / 64; i &lt; uint(len(b)); i++ {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		x := b[i]
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		if ^x == 0 {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			continue
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		return i*64 + uint(sys.TrailingZeros64(^x))
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	return ^uint(0)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span><span class="comment">// findSmallN is a helper for find which searches for npages contiguous free pages</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">// in this pallocBits and returns the index where that run of contiguous pages</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span><span class="comment">// starts as well as the index of the first free page it finds in its search.</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">// See find for an explanation of the searchIdx parameter.</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">// Returns a ^uint(0) index on failure and the new searchIdx should be ignored.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span><span class="comment">// findSmallN assumes npages &lt;= 64, where any such contiguous run of pages</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span><span class="comment">// crosses at most one aligned 64-bit boundary in the bits.</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>func (b *pallocBits) findSmallN(npages uintptr, searchIdx uint) (uint, uint) {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	end, newSearchIdx := uint(0), ^uint(0)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	for i := searchIdx / 64; i &lt; uint(len(b)); i++ {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		bi := b[i]
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		if ^bi == 0 {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			end = 0
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			continue
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		<span class="comment">// First see if we can pack our allocation in the trailing</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		<span class="comment">// zeros plus the end of the last 64 bits.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		if newSearchIdx == ^uint(0) {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			<span class="comment">// The new searchIdx is going to be at these 64 bits after any</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			<span class="comment">// 1s we file, so count trailing 1s.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			newSearchIdx = i*64 + uint(sys.TrailingZeros64(^bi))
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		start := uint(sys.TrailingZeros64(bi))
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		if end+start &gt;= uint(npages) {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			return i*64 - end, newSearchIdx
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		<span class="comment">// Next, check the interior of the 64-bit chunk.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		j := findBitRange64(^bi, uint(npages))
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		if j &lt; 64 {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			return i*64 + j, newSearchIdx
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		end = uint(sys.LeadingZeros64(bi))
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	return ^uint(0), newSearchIdx
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">// findLargeN is a helper for find which searches for npages contiguous free pages</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span><span class="comment">// in this pallocBits and returns the index where that run starts, as well as the</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span><span class="comment">// index of the first free page it found it its search.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">// See alloc for an explanation of the searchIdx parameter.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">// Returns a ^uint(0) index on failure and the new searchIdx should be ignored.</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// findLargeN assumes npages &gt; 64, where any such run of free pages</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// crosses at least one aligned 64-bit boundary in the bits.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>func (b *pallocBits) findLargeN(npages uintptr, searchIdx uint) (uint, uint) {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	start, size, newSearchIdx := ^uint(0), uint(0), ^uint(0)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	for i := searchIdx / 64; i &lt; uint(len(b)); i++ {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		x := b[i]
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		if x == ^uint64(0) {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>			size = 0
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			continue
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		if newSearchIdx == ^uint(0) {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			<span class="comment">// The new searchIdx is going to be at these 64 bits after any</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			<span class="comment">// 1s we file, so count trailing 1s.</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>			newSearchIdx = i*64 + uint(sys.TrailingZeros64(^x))
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		if size == 0 {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			size = uint(sys.LeadingZeros64(x))
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			start = i*64 + 64 - size
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			continue
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		s := uint(sys.TrailingZeros64(x))
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		if s+size &gt;= uint(npages) {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			size += s
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			return start, newSearchIdx
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		if s &lt; 64 {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>			size = uint(sys.LeadingZeros64(x))
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			start = i*64 + 64 - size
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			continue
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		size += 64
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	if size &lt; uint(npages) {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		return ^uint(0), newSearchIdx
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	return start, newSearchIdx
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// allocRange allocates the range [i, i+n).</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>func (b *pallocBits) allocRange(i, n uint) {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	(*pageBits)(b).setRange(i, n)
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// allocAll allocates all the bits of b.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>func (b *pallocBits) allocAll() {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	(*pageBits)(b).setAll()
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>}
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// free1 frees a single page in the pallocBits at i.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>func (b *pallocBits) free1(i uint) {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	(*pageBits)(b).clear(i)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// free frees the range [i, i+n) of pages in the pallocBits.</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>func (b *pallocBits) free(i, n uint) {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	(*pageBits)(b).clearRange(i, n)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">// freeAll frees all the bits of b.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>func (b *pallocBits) freeAll() {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	(*pageBits)(b).clearAll()
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span><span class="comment">// pages64 returns a 64-bit bitmap representing a block of 64 pages aligned</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span><span class="comment">// to 64 pages. The returned block of pages is the one containing the i&#39;th</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span><span class="comment">// page in this pallocBits. Each bit represents whether the page is in-use.</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>func (b *pallocBits) pages64(i uint) uint64 {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	return (*pageBits)(b).block64(i)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span><span class="comment">// allocPages64 allocates a 64-bit block of 64 pages aligned to 64 pages according</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span><span class="comment">// to the bits set in alloc. The block set is the one containing the i&#39;th page.</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>func (b *pallocBits) allocPages64(i uint, alloc uint64) {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	(*pageBits)(b).setBlock64(i, alloc)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>}
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span><span class="comment">// findBitRange64 returns the bit index of the first set of</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span><span class="comment">// n consecutive 1 bits. If no consecutive set of 1 bits of</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span><span class="comment">// size n may be found in c, then it returns an integer &gt;= 64.</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span><span class="comment">// n must be &gt; 0.</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>func findBitRange64(c uint64, n uint) uint {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">// This implementation is based on shrinking the length of</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	<span class="comment">// runs of contiguous 1 bits. We remove the top n-1 1 bits</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	<span class="comment">// from each run of 1s, then look for the first remaining 1 bit.</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	p := n - 1   <span class="comment">// number of 1s we want to remove.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	k := uint(1) <span class="comment">// current minimum width of runs of 0 in c.</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	for p &gt; 0 {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		if p &lt;= k {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			<span class="comment">// Shift p 0s down into the top of each run of 1s.</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			c &amp;= c &gt;&gt; (p &amp; 63)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			break
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		<span class="comment">// Shift k 0s down into the top of each run of 1s.</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		c &amp;= c &gt;&gt; (k &amp; 63)
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		if c == 0 {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			return 64
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		p -= k
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		<span class="comment">// We&#39;ve just doubled the minimum length of 0-runs.</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		<span class="comment">// This allows us to shift farther in the next iteration.</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		k *= 2
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	<span class="comment">// Find first remaining 1.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	<span class="comment">// Since we shrunk from the top down, the first 1 is in</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	<span class="comment">// its correct original position.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	return uint(sys.TrailingZeros64(c))
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span><span class="comment">// pallocData encapsulates pallocBits and a bitmap for</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span><span class="comment">// whether or not a given page is scavenged in a single</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span><span class="comment">// structure. It&#39;s effectively a pallocBits with</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span><span class="comment">// additional functionality.</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span><span class="comment">// Update the comment on (*pageAlloc).chunks should this</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span><span class="comment">// structure change.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>type pallocData struct {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	pallocBits
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	scavenged pageBits
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span><span class="comment">// allocRange sets bits [i, i+n) in the bitmap to 1 and</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span><span class="comment">// updates the scavenged bits appropriately.</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>func (m *pallocData) allocRange(i, n uint) {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	<span class="comment">// Clear the scavenged bits when we alloc the range.</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	m.pallocBits.allocRange(i, n)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	m.scavenged.clearRange(i, n)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span><span class="comment">// allocAll sets every bit in the bitmap to 1 and updates</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span><span class="comment">// the scavenged bits appropriately.</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>func (m *pallocData) allocAll() {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	<span class="comment">// Clear the scavenged bits when we alloc the range.</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	m.pallocBits.allocAll()
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	m.scavenged.clearAll()
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>
</pre><p><a href="mpallocbits.go?m=text">View as plain text</a></p>

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
