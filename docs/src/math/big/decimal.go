<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/big/decimal.go - Go Documentation Server</title>

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
<a href="decimal.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<a href="http://localhost:8080/src/math/big">big</a>/<span class="text-muted">decimal.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math/big">math/big</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2015 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements multi-precision decimal numbers.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// The implementation is for float to decimal conversion only;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// not general purpose use.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// The only operations are precise conversion from binary to</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// decimal and rounding.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// The key observation and some code (shr) is borrowed from</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// strconv/decimal.go: conversion of binary fractional values can be done</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// precisely in multi-precision decimal because 2 divides 10 (required for</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// &gt;&gt; of mantissa); but conversion of decimal floating-point values cannot</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// be done precisely in binary representation.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// In contrast to strconv/decimal.go, only right shift is implemented in</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// decimal format - left shift can be done precisely in binary format.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>package big
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// A decimal represents an unsigned floating-point number in decimal representation.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// The value of a non-zero decimal d is d.mant * 10**d.exp with 0.1 &lt;= d.mant &lt; 1,</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// with the most-significant mantissa digit at index 0. For the zero decimal, the</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// mantissa length and exponent are 0.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// The zero value for decimal represents a ready-to-use 0.0.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>type decimal struct {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	mant []byte <span class="comment">// mantissa ASCII digits, big-endian</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	exp  int    <span class="comment">// exponent</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// at returns the i&#39;th mantissa digit, starting with the most significant digit at 0.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>func (d *decimal) at(i int) byte {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	if 0 &lt;= i &amp;&amp; i &lt; len(d.mant) {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		return d.mant[i]
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	return &#39;0&#39;
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// Maximum shift amount that can be done in one pass without overflow.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// A Word has _W bits and (1&lt;&lt;maxShift - 1)*10 + 9 must fit into Word.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>const maxShift = _W - 4
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// TODO(gri) Since we know the desired decimal precision when converting</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// a floating-point number, we may be able to limit the number of decimal</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// digits that need to be computed by init by providing an additional</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// precision argument and keeping track of when a number was truncated early</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// (equivalent of &#34;sticky bit&#34; in binary rounding).</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// TODO(gri) Along the same lines, enforce some limit to shift magnitudes</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// to avoid &#34;infinitely&#34; long running conversions (until we run out of space).</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// Init initializes x to the decimal representation of m &lt;&lt; shift (for</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// shift &gt;= 0), or m &gt;&gt; -shift (for shift &lt; 0).</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>func (x *decimal) init(m nat, shift int) {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// special case 0</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	if len(m) == 0 {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		x.mant = x.mant[:0]
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		x.exp = 0
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		return
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// Optimization: If we need to shift right, first remove any trailing</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// zero bits from m to reduce shift amount that needs to be done in</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// decimal format (since that is likely slower).</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	if shift &lt; 0 {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		ntz := m.trailingZeroBits()
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		s := uint(-shift)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		if s &gt;= ntz {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			s = ntz <span class="comment">// shift at most ntz bits</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		m = nat(nil).shr(m, s)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		shift += int(s)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// Do any shift left in binary representation.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	if shift &gt; 0 {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		m = nat(nil).shl(m, uint(shift))
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		shift = 0
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// Convert mantissa into decimal representation.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	s := m.utoa(10)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	n := len(s)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	x.exp = n
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// Trim trailing zeros; instead the exponent is tracking</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// the decimal point independent of the number of digits.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	for n &gt; 0 &amp;&amp; s[n-1] == &#39;0&#39; {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		n--
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	x.mant = append(x.mant[:0], s[:n]...)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// Do any (remaining) shift right in decimal representation.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	if shift &lt; 0 {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		for shift &lt; -maxShift {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			shr(x, maxShift)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			shift += maxShift
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		shr(x, uint(-shift))
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// shr implements x &gt;&gt; s, for s &lt;= maxShift.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>func shr(x *decimal, s uint) {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// Division by 1&lt;&lt;s using shift-and-subtract algorithm.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// pick up enough leading digits to cover first shift</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	r := 0 <span class="comment">// read index</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	var n Word
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	for n&gt;&gt;s == 0 &amp;&amp; r &lt; len(x.mant) {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		ch := Word(x.mant[r])
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		r++
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		n = n*10 + ch - &#39;0&#39;
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if n == 0 {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		<span class="comment">// x == 0; shouldn&#39;t get here, but handle anyway</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		x.mant = x.mant[:0]
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		return
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	for n&gt;&gt;s == 0 {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		r++
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		n *= 10
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	x.exp += 1 - r
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// read a digit, write a digit</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	w := 0 <span class="comment">// write index</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	mask := Word(1)&lt;&lt;s - 1
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	for r &lt; len(x.mant) {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		ch := Word(x.mant[r])
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		r++
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		d := n &gt;&gt; s
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		n &amp;= mask <span class="comment">// n -= d &lt;&lt; s</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		x.mant[w] = byte(d + &#39;0&#39;)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		w++
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		n = n*10 + ch - &#39;0&#39;
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// write extra digits that still fit</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	for n &gt; 0 &amp;&amp; w &lt; len(x.mant) {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		d := n &gt;&gt; s
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		n &amp;= mask
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		x.mant[w] = byte(d + &#39;0&#39;)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		w++
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		n = n * 10
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	x.mant = x.mant[:w] <span class="comment">// the number may be shorter (e.g. 1024 &gt;&gt; 10)</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// append additional digits that didn&#39;t fit</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	for n &gt; 0 {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		d := n &gt;&gt; s
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		n &amp;= mask
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		x.mant = append(x.mant, byte(d+&#39;0&#39;))
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		n = n * 10
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	trim(x)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func (x *decimal) String() string {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if len(x.mant) == 0 {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		return &#34;0&#34;
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	var buf []byte
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	switch {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	case x.exp &lt;= 0:
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// 0.00ddd</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		buf = make([]byte, 0, 2+(-x.exp)+len(x.mant))
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		buf = append(buf, &#34;0.&#34;...)
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		buf = appendZeros(buf, -x.exp)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		buf = append(buf, x.mant...)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	case <span class="comment">/* 0 &lt; */</span> x.exp &lt; len(x.mant):
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		<span class="comment">// dd.ddd</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		buf = make([]byte, 0, 1+len(x.mant))
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		buf = append(buf, x.mant[:x.exp]...)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		buf = append(buf, &#39;.&#39;)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		buf = append(buf, x.mant[x.exp:]...)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	default: <span class="comment">// len(x.mant) &lt;= x.exp</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		<span class="comment">// ddd00</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		buf = make([]byte, 0, x.exp)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		buf = append(buf, x.mant...)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		buf = appendZeros(buf, x.exp-len(x.mant))
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	return string(buf)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// appendZeros appends n 0 digits to buf and returns buf.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>func appendZeros(buf []byte, n int) []byte {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	for ; n &gt; 0; n-- {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		buf = append(buf, &#39;0&#39;)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	return buf
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// shouldRoundUp reports if x should be rounded up</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// if shortened to n digits. n must be a valid index</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// for x.mant.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>func shouldRoundUp(x *decimal, n int) bool {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	if x.mant[n] == &#39;5&#39; &amp;&amp; n+1 == len(x.mant) {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		<span class="comment">// exactly halfway - round to even</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		return n &gt; 0 &amp;&amp; (x.mant[n-1]-&#39;0&#39;)&amp;1 != 0
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// not halfway - digit tells all (x.mant has no trailing zeros)</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	return x.mant[n] &gt;= &#39;5&#39;
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// round sets x to (at most) n mantissa digits by rounding it</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// to the nearest even value with n (or fever) mantissa digits.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// If n &lt; 0, x remains unchanged.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>func (x *decimal) round(n int) {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	if n &lt; 0 || n &gt;= len(x.mant) {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		return <span class="comment">// nothing to do</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	if shouldRoundUp(x, n) {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		x.roundUp(n)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	} else {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		x.roundDown(n)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>func (x *decimal) roundUp(n int) {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	if n &lt; 0 || n &gt;= len(x.mant) {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		return <span class="comment">// nothing to do</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// 0 &lt;= n &lt; len(x.mant)</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">// find first digit &lt; &#39;9&#39;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	for n &gt; 0 &amp;&amp; x.mant[n-1] &gt;= &#39;9&#39; {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		n--
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	if n == 0 {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		<span class="comment">// all digits are &#39;9&#39;s =&gt; round up to &#39;1&#39; and update exponent</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		x.mant[0] = &#39;1&#39; <span class="comment">// ok since len(x.mant) &gt; n</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		x.mant = x.mant[:1]
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		x.exp++
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		return
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// n &gt; 0 &amp;&amp; x.mant[n-1] &lt; &#39;9&#39;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	x.mant[n-1]++
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	x.mant = x.mant[:n]
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// x already trimmed</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func (x *decimal) roundDown(n int) {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	if n &lt; 0 || n &gt;= len(x.mant) {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		return <span class="comment">// nothing to do</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	x.mant = x.mant[:n]
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	trim(x)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">// trim cuts off any trailing zeros from x&#39;s mantissa;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// they are meaningless for the value of x.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>func trim(x *decimal) {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	i := len(x.mant)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	for i &gt; 0 &amp;&amp; x.mant[i-1] == &#39;0&#39; {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		i--
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	x.mant = x.mant[:i]
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	if i == 0 {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		x.exp = 0
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
</pre><p><a href="decimal.go?m=text">View as plain text</a></p>

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
