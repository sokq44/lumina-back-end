<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/big/nat.go - Go Documentation Server</title>

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
<a href="nat.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<a href="http://localhost:8080/src/math/big">big</a>/<span class="text-muted">nat.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math/big">math/big</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements unsigned multi-precision integers (natural</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// numbers). They are the building blocks for the implementation</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// of signed integers, rationals, and floating-point numbers.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// Caution: This implementation relies on the function &#34;alias&#34;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//          which assumes that (nat) slice capacities are never</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//          changed (no 3-operand slice expressions). If that</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//          changes, alias needs to be updated for correctness.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>package big
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>import (
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;encoding/binary&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;math/bits&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;math/rand&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// An unsigned integer x of the form</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//	x = x[n-1]*_B^(n-1) + x[n-2]*_B^(n-2) + ... + x[1]*_B + x[0]</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// with 0 &lt;= x[i] &lt; _B and 0 &lt;= i &lt; n is stored in a slice of length n,</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// with the digits x[i] as the slice elements.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// A number is normalized if the slice contains no leading 0 digits.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// During arithmetic operations, denormalized values may occur but are</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// always normalized before returning the final result. The normalized</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// representation of 0 is the empty or nil slice (length = 0).</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>type nat []Word
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>var (
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	natOne  = nat{1}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	natTwo  = nat{2}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	natFive = nat{5}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	natTen  = nat{10}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func (z nat) String() string {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	return &#34;0x&#34; + string(z.itoa(false, 16))
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>func (z nat) clear() {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	for i := range z {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		z[i] = 0
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>func (z nat) norm() nat {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	i := len(z)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	for i &gt; 0 &amp;&amp; z[i-1] == 0 {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		i--
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	return z[0:i]
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>func (z nat) make(n int) nat {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	if n &lt;= cap(z) {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		return z[:n] <span class="comment">// reuse z</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if n == 1 {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		<span class="comment">// Most nats start small and stay that way; don&#39;t over-allocate.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		return make(nat, 1)
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// Choosing a good value for e has significant performance impact</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// because it increases the chance that a value can be reused.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	const e = 4 <span class="comment">// extra capacity</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	return make(nat, n, n+e)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>func (z nat) setWord(x Word) nat {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	if x == 0 {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		return z[:0]
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	z = z.make(1)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	z[0] = x
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	return z
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>func (z nat) setUint64(x uint64) nat {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// single-word value</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	if w := Word(x); uint64(w) == x {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		return z.setWord(w)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// 2-word value</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	z = z.make(2)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	z[1] = Word(x &gt;&gt; 32)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	z[0] = Word(x)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	return z
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func (z nat) set(x nat) nat {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	z = z.make(len(x))
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	copy(z, x)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	return z
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>func (z nat) add(x, y nat) nat {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	m := len(x)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	n := len(y)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	switch {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	case m &lt; n:
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		return z.add(y, x)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	case m == 0:
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		<span class="comment">// n == 0 because m &gt;= n; result is 0</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		return z[:0]
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	case n == 0:
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		<span class="comment">// result is x</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		return z.set(x)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// m &gt; 0</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	z = z.make(m + 1)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	c := addVV(z[0:n], x, y)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	if m &gt; n {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		c = addVW(z[n:m], x[n:], c)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	z[m] = c
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	return z.norm()
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>func (z nat) sub(x, y nat) nat {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	m := len(x)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	n := len(y)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	switch {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	case m &lt; n:
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		panic(&#34;underflow&#34;)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	case m == 0:
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// n == 0 because m &gt;= n; result is 0</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		return z[:0]
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	case n == 0:
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		<span class="comment">// result is x</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		return z.set(x)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// m &gt; 0</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	z = z.make(m)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	c := subVV(z[0:n], x, y)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	if m &gt; n {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		c = subVW(z[n:], x[n:], c)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	if c != 0 {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		panic(&#34;underflow&#34;)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	return z.norm()
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>func (x nat) cmp(y nat) (r int) {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	m := len(x)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	n := len(y)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if m != n || m == 0 {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		switch {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		case m &lt; n:
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			r = -1
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		case m &gt; n:
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			r = 1
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		return
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	i := m - 1
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	for i &gt; 0 &amp;&amp; x[i] == y[i] {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		i--
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	switch {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	case x[i] &lt; y[i]:
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		r = -1
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	case x[i] &gt; y[i]:
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		r = 1
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	return
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>func (z nat) mulAddWW(x nat, y, r Word) nat {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	m := len(x)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if m == 0 || y == 0 {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		return z.setWord(r) <span class="comment">// result is r</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">// m &gt; 0</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	z = z.make(m + 1)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	z[m] = mulAddVWW(z[0:m], x, y, r)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	return z.norm()
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// basicMul multiplies x and y and leaves the result in z.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// The (non-normalized) result is placed in z[0 : len(x) + len(y)].</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>func basicMul(z, x, y nat) {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	z[0 : len(x)+len(y)].clear() <span class="comment">// initialize z</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	for i, d := range y {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		if d != 0 {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>			z[len(x)+i] = addMulVVW(z[i:i+len(x)], x, d)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// montgomery computes z mod m = x*y*2**(-n*_W) mod m,</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// assuming k = -1/m mod 2**_W.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// z is used for storing the result which is returned;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// z must not alias x, y or m.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// See Gueron, &#34;Efficient Software Implementations of Modular Exponentiation&#34;.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// https://eprint.iacr.org/2011/239.pdf</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// In the terminology of that paper, this is an &#34;Almost Montgomery Multiplication&#34;:</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// x and y are required to satisfy 0 &lt;= z &lt; 2**(n*_W) and then the result</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// z is guaranteed to satisfy 0 &lt;= z &lt; 2**(n*_W), but it may not be &lt; m.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>func (z nat) montgomery(x, y, m nat, k Word, n int) nat {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	<span class="comment">// This code assumes x, y, m are all the same length, n.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	<span class="comment">// (required by addMulVVW and the for loop).</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	<span class="comment">// It also assumes that x, y are already reduced mod m,</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// or else the result will not be properly reduced.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	if len(x) != n || len(y) != n || len(m) != n {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		panic(&#34;math/big: mismatched montgomery number lengths&#34;)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	z = z.make(n * 2)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	z.clear()
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	var c Word
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		d := y[i]
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		c2 := addMulVVW(z[i:n+i], x, d)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		t := z[i] * k
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		c3 := addMulVVW(z[i:n+i], m, t)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		cx := c + c2
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		cy := cx + c3
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		z[n+i] = cy
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		if cx &lt; c2 || cy &lt; c3 {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			c = 1
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		} else {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			c = 0
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	if c != 0 {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		subVV(z[:n], z[n:], m)
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	} else {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		copy(z[:n], z[n:])
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	return z[:n]
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">// Fast version of z[0:n+n&gt;&gt;1].add(z[0:n+n&gt;&gt;1], x[0:n]) w/o bounds checks.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">// Factored out for readability - do not use outside karatsuba.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func karatsubaAdd(z, x nat, n int) {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	if c := addVV(z[0:n], z, x); c != 0 {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		addVW(z[n:n+n&gt;&gt;1], z[n:], c)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span><span class="comment">// Like karatsubaAdd, but does subtract.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>func karatsubaSub(z, x nat, n int) {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	if c := subVV(z[0:n], z, x); c != 0 {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		subVW(z[n:n+n&gt;&gt;1], z[n:], c)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span><span class="comment">// Operands that are shorter than karatsubaThreshold are multiplied using</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">// &#34;grade school&#34; multiplication; for longer operands the Karatsuba algorithm</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">// is used.</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>var karatsubaThreshold = 40 <span class="comment">// computed by calibrate_test.go</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">// karatsuba multiplies x and y and leaves the result in z.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span><span class="comment">// Both x and y must have the same length n and n must be a</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span><span class="comment">// power of 2. The result vector z must have len(z) &gt;= 6*n.</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span><span class="comment">// The (non-normalized) result is placed in z[0 : 2*n].</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>func karatsuba(z, x, y nat) {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	n := len(y)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">// Switch to basic multiplication if numbers are odd or small.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	<span class="comment">// (n is always even if karatsubaThreshold is even, but be</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	<span class="comment">// conservative)</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	if n&amp;1 != 0 || n &lt; karatsubaThreshold || n &lt; 2 {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		basicMul(z, x, y)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		return
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	<span class="comment">// n&amp;1 == 0 &amp;&amp; n &gt;= karatsubaThreshold &amp;&amp; n &gt;= 2</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">// Karatsuba multiplication is based on the observation that</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	<span class="comment">// for two numbers x and y with:</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	<span class="comment">//   x = x1*b + x0</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	<span class="comment">//   y = y1*b + y0</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// the product x*y can be obtained with 3 products z2, z1, z0</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// instead of 4:</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">//   x*y = x1*y1*b*b + (x1*y0 + x0*y1)*b + x0*y0</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">//       =    z2*b*b +              z1*b +    z0</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// with:</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">//   xd = x1 - x0</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">//   yd = y0 - y1</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	<span class="comment">//   z1 =      xd*yd                    + z2 + z0</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	<span class="comment">//      = (x1-x0)*(y0 - y1)             + z2 + z0</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	<span class="comment">//      = x1*y0 - x1*y1 - x0*y0 + x0*y1 + z2 + z0</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">//      = x1*y0 -    z2 -    z0 + x0*y1 + z2 + z0</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	<span class="comment">//      = x1*y0                 + x0*y1</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	<span class="comment">// split x, y into &#34;digits&#34;</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	n2 := n &gt;&gt; 1              <span class="comment">// n2 &gt;= 1</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	x1, x0 := x[n2:], x[0:n2] <span class="comment">// x = x1*b + y0</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	y1, y0 := y[n2:], y[0:n2] <span class="comment">// y = y1*b + y0</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	<span class="comment">// z is used for the result and temporary storage:</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	<span class="comment">//   6*n     5*n     4*n     3*n     2*n     1*n     0*n</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// z = [z2 copy|z0 copy| xd*yd | yd:xd | x1*y1 | x0*y0 ]</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	<span class="comment">// For each recursive call of karatsuba, an unused slice of</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	<span class="comment">// z is passed in that has (at least) half the length of the</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	<span class="comment">// caller&#39;s z.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	<span class="comment">// compute z0 and z2 with the result &#34;in place&#34; in z</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	karatsuba(z, x0, y0)     <span class="comment">// z0 = x0*y0</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	karatsuba(z[n:], x1, y1) <span class="comment">// z2 = x1*y1</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">// compute xd (or the negative value if underflow occurs)</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	s := 1 <span class="comment">// sign of product xd*yd</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	xd := z[2*n : 2*n+n2]
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	if subVV(xd, x1, x0) != 0 { <span class="comment">// x1-x0</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		s = -s
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		subVV(xd, x0, x1) <span class="comment">// x0-x1</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// compute yd (or the negative value if underflow occurs)</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	yd := z[2*n+n2 : 3*n]
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	if subVV(yd, y0, y1) != 0 { <span class="comment">// y0-y1</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		s = -s
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		subVV(yd, y1, y0) <span class="comment">// y1-y0</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	<span class="comment">// p = (x1-x0)*(y0-y1) == x1*y0 - x1*y1 - x0*y0 + x0*y1 for s &gt; 0</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	<span class="comment">// p = (x0-x1)*(y0-y1) == x0*y0 - x0*y1 - x1*y0 + x1*y1 for s &lt; 0</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	p := z[n*3:]
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	karatsuba(p, xd, yd)
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	<span class="comment">// save original z2:z0</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	<span class="comment">// (ok to use upper half of z since we&#39;re done recurring)</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	r := z[n*4:]
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	copy(r, z[:n*2])
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	<span class="comment">// add up all partial products</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	<span class="comment">//   2*n     n     0</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	<span class="comment">// z = [ z2  | z0  ]</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	<span class="comment">//   +    [ z0  ]</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	<span class="comment">//   +    [ z2  ]</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	<span class="comment">//   +    [  p  ]</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	karatsubaAdd(z[n2:], r, n)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	karatsubaAdd(z[n2:], r[n:], n)
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	if s &gt; 0 {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		karatsubaAdd(z[n2:], p, n)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	} else {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		karatsubaSub(z[n2:], p, n)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">// alias reports whether x and y share the same base array.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">// Note: alias assumes that the capacity of underlying arrays</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span><span class="comment">// is never changed for nat values; i.e. that there are</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span><span class="comment">// no 3-operand slice expressions in this code (or worse,</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span><span class="comment">// reflect-based operations to the same effect).</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>func alias(x, y nat) bool {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	return cap(x) &gt; 0 &amp;&amp; cap(y) &gt; 0 &amp;&amp; &amp;x[0:cap(x)][cap(x)-1] == &amp;y[0:cap(y)][cap(y)-1]
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span><span class="comment">// addAt implements z += x&lt;&lt;(_W*i); z must be long enough.</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span><span class="comment">// (we don&#39;t use nat.add because we need z to stay the same</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span><span class="comment">// slice, and we don&#39;t need to normalize z after each addition)</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>func addAt(z, x nat, i int) {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	if n := len(x); n &gt; 0 {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		if c := addVV(z[i:i+n], z[i:], x); c != 0 {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			j := i + n
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>			if j &lt; len(z) {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>				addVW(z[j:], z[j:], c)
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>			}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>}
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span><span class="comment">// karatsubaLen computes an approximation to the maximum k &lt;= n such that</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">// k = p&lt;&lt;i for a number p &lt;= threshold and an i &gt;= 0. Thus, the</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">// result is the largest number that can be divided repeatedly by 2 before</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// becoming about the value of threshold.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>func karatsubaLen(n, threshold int) int {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	i := uint(0)
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	for n &gt; threshold {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		n &gt;&gt;= 1
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		i++
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	return n &lt;&lt; i
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>}
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>func (z nat) mul(x, y nat) nat {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	m := len(x)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	n := len(y)
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	switch {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	case m &lt; n:
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		return z.mul(y, x)
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	case m == 0 || n == 0:
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		return z[:0]
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	case n == 1:
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		return z.mulAddWW(x, y[0], 0)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	<span class="comment">// m &gt;= n &gt; 1</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	<span class="comment">// determine if z can be reused</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	if alias(z, x) || alias(z, y) {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		z = nil <span class="comment">// z is an alias for x or y - cannot reuse</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	<span class="comment">// use basic multiplication if the numbers are small</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	if n &lt; karatsubaThreshold {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		z = z.make(m + n)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		basicMul(z, x, y)
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		return z.norm()
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	<span class="comment">// m &gt;= n &amp;&amp; n &gt;= karatsubaThreshold &amp;&amp; n &gt;= 2</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	<span class="comment">// determine Karatsuba length k such that</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	<span class="comment">//   x = xh*b + x0  (0 &lt;= x0 &lt; b)</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	<span class="comment">//   y = yh*b + y0  (0 &lt;= y0 &lt; b)</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	<span class="comment">//   b = 1&lt;&lt;(_W*k)  (&#34;base&#34; of digits xi, yi)</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	k := karatsubaLen(n, karatsubaThreshold)
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	<span class="comment">// k &lt;= n</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// multiply x0 and y0 via Karatsuba</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	x0 := x[0:k]              <span class="comment">// x0 is not normalized</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	y0 := y[0:k]              <span class="comment">// y0 is not normalized</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	z = z.make(max(6*k, m+n)) <span class="comment">// enough space for karatsuba of x0*y0 and full result of x*y</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	karatsuba(z, x0, y0)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	z = z[0 : m+n]  <span class="comment">// z has final length but may be incomplete</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	z[2*k:].clear() <span class="comment">// upper portion of z is garbage (and 2*k &lt;= m+n since k &lt;= n &lt;= m)</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	<span class="comment">// If xh != 0 or yh != 0, add the missing terms to z. For</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	<span class="comment">//   xh = xi*b^i + ... + x2*b^2 + x1*b (0 &lt;= xi &lt; b)</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	<span class="comment">//   yh =                         y1*b (0 &lt;= y1 &lt; b)</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	<span class="comment">// the missing terms are</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	<span class="comment">//   x0*y1*b and xi*y0*b^i, xi*y1*b^(i+1) for i &gt; 0</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	<span class="comment">// since all the yi for i &gt; 1 are 0 by choice of k: If any of them</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	<span class="comment">// were &gt; 0, then yh &gt;= b^2 and thus y &gt;= b^2. Then k&#39; = k*2 would</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	<span class="comment">// be a larger valid threshold contradicting the assumption about k.</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	if k &lt; n || m != n {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		tp := getNat(3 * k)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		t := *tp
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		<span class="comment">// add x0*y1*b</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		x0 := x0.norm()
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		y1 := y[k:]       <span class="comment">// y1 is normalized because y is</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		t = t.mul(x0, y1) <span class="comment">// update t so we don&#39;t lose t&#39;s underlying array</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		addAt(z, t, k)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		<span class="comment">// add xi*y0&lt;&lt;i, xi*y1*b&lt;&lt;(i+k)</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		y0 := y0.norm()
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		for i := k; i &lt; len(x); i += k {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			xi := x[i:]
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			if len(xi) &gt; k {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>				xi = xi[:k]
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			xi = xi.norm()
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			t = t.mul(xi, y0)
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			addAt(z, t, i)
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			t = t.mul(xi, y1)
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			addAt(z, t, i+k)
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		putNat(tp)
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	return z.norm()
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span><span class="comment">// basicSqr sets z = x*x and is asymptotically faster than basicMul</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span><span class="comment">// by about a factor of 2, but slower for small arguments due to overhead.</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span><span class="comment">// Requirements: len(x) &gt; 0, len(z) == 2*len(x)</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span><span class="comment">// The (non-normalized) result is placed in z.</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>func basicSqr(z, x nat) {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	n := len(x)
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	tp := getNat(2 * n)
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	t := *tp <span class="comment">// temporary variable to hold the products</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	t.clear()
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	z[1], z[0] = mulWW(x[0], x[0]) <span class="comment">// the initial square</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	for i := 1; i &lt; n; i++ {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		d := x[i]
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		<span class="comment">// z collects the squares x[i] * x[i]</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		z[2*i+1], z[2*i] = mulWW(d, d)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		<span class="comment">// t collects the products x[i] * x[j] where j &lt; i</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		t[2*i] = addMulVVW(t[i:2*i], x[0:i], d)
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	t[2*n-1] = shlVU(t[1:2*n-1], t[1:2*n-1], 1) <span class="comment">// double the j &lt; i products</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	addVV(z, z, t)                              <span class="comment">// combine the result</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	putNat(tp)
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span><span class="comment">// karatsubaSqr squares x and leaves the result in z.</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span><span class="comment">// len(x) must be a power of 2 and len(z) &gt;= 6*len(x).</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span><span class="comment">// The (non-normalized) result is placed in z[0 : 2*len(x)].</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span><span class="comment">// The algorithm and the layout of z are the same as for karatsuba.</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>func karatsubaSqr(z, x nat) {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	n := len(x)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	if n&amp;1 != 0 || n &lt; karatsubaSqrThreshold || n &lt; 2 {
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		basicSqr(z[:2*n], x)
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		return
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	}
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	n2 := n &gt;&gt; 1
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	x1, x0 := x[n2:], x[0:n2]
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	karatsubaSqr(z, x0)
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	karatsubaSqr(z[n:], x1)
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	<span class="comment">// s = sign(xd*yd) == -1 for xd != 0; s == 1 for xd == 0</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	xd := z[2*n : 2*n+n2]
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	if subVV(xd, x1, x0) != 0 {
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		subVV(xd, x0, x1)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	p := z[n*3:]
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	karatsubaSqr(p, xd)
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	r := z[n*4:]
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	copy(r, z[:n*2])
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	karatsubaAdd(z[n2:], r, n)
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	karatsubaAdd(z[n2:], r[n:], n)
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	karatsubaSub(z[n2:], p, n) <span class="comment">// s == -1 for p != 0; s == 1 for p == 0</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span><span class="comment">// Operands that are shorter than basicSqrThreshold are squared using</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span><span class="comment">// &#34;grade school&#34; multiplication; for operands longer than karatsubaSqrThreshold</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// we use the Karatsuba algorithm optimized for x == y.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>var basicSqrThreshold = 20      <span class="comment">// computed by calibrate_test.go</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>var karatsubaSqrThreshold = 260 <span class="comment">// computed by calibrate_test.go</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span><span class="comment">// z = x*x</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>func (z nat) sqr(x nat) nat {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	n := len(x)
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	switch {
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	case n == 0:
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		return z[:0]
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	case n == 1:
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		d := x[0]
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		z = z.make(2)
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		z[1], z[0] = mulWW(d, d)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		return z.norm()
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	if alias(z, x) {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		z = nil <span class="comment">// z is an alias for x - cannot reuse</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	if n &lt; basicSqrThreshold {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		z = z.make(2 * n)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		basicMul(z, x, x)
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		return z.norm()
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	}
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	if n &lt; karatsubaSqrThreshold {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		z = z.make(2 * n)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		basicSqr(z, x)
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		return z.norm()
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	}
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	<span class="comment">// Use Karatsuba multiplication optimized for x == y.</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	<span class="comment">// The algorithm and layout of z are the same as for mul.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	<span class="comment">// z = (x1*b + x0)^2 = x1^2*b^2 + 2*x1*x0*b + x0^2</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	k := karatsubaLen(n, karatsubaSqrThreshold)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	x0 := x[0:k]
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	z = z.make(max(6*k, 2*n))
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	karatsubaSqr(z, x0) <span class="comment">// z = x0^2</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	z = z[0 : 2*n]
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	z[2*k:].clear()
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	if k &lt; n {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		tp := getNat(2 * k)
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		t := *tp
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		x0 := x0.norm()
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		x1 := x[k:]
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		t = t.mul(x0, x1)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		addAt(z, t, k)
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		addAt(z, t, k) <span class="comment">// z = 2*x1*x0*b + x0^2</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		t = t.sqr(x1)
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		addAt(z, t, 2*k) <span class="comment">// z = x1^2*b^2 + 2*x1*x0*b + x0^2</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		putNat(tp)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	return z.norm()
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>}
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span><span class="comment">// mulRange computes the product of all the unsigned integers in the</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span><span class="comment">// range [a, b] inclusively. If a &gt; b (empty range), the result is 1.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>func (z nat) mulRange(a, b uint64) nat {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	switch {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	case a == 0:
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		<span class="comment">// cut long ranges short (optimization)</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		return z.setUint64(0)
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	case a &gt; b:
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		return z.setUint64(1)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	case a == b:
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		return z.setUint64(a)
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	case a+1 == b:
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		return z.mul(nat(nil).setUint64(a), nat(nil).setUint64(b))
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	m := a + (b-a)/2 <span class="comment">// avoid overflow</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	return z.mul(nat(nil).mulRange(a, m), nat(nil).mulRange(m+1, b))
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>}
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span><span class="comment">// getNat returns a *nat of len n. The contents may not be zero.</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span><span class="comment">// The pool holds *nat to avoid allocation when converting to interface{}.</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>func getNat(n int) *nat {
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	var z *nat
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	if v := natPool.Get(); v != nil {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		z = v.(*nat)
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	}
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	if z == nil {
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		z = new(nat)
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	}
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	*z = z.make(n)
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	if n &gt; 0 {
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		(*z)[0] = 0xfedcb <span class="comment">// break code expecting zero</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	}
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	return z
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>func putNat(x *nat) {
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	natPool.Put(x)
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>}
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>var natPool sync.Pool
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span><span class="comment">// bitLen returns the length of x in bits.</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span><span class="comment">// Unlike most methods, it works even if x is not normalized.</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>func (x nat) bitLen() int {
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	<span class="comment">// This function is used in cryptographic operations. It must not leak</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	<span class="comment">// anything but the Int&#39;s sign and bit size through side-channels. Any</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	<span class="comment">// changes must be reviewed by a security expert.</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	if i := len(x) - 1; i &gt;= 0 {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		<span class="comment">// bits.Len uses a lookup table for the low-order bits on some</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		<span class="comment">// architectures. Neutralize any input-dependent behavior by setting all</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		<span class="comment">// bits after the first one bit.</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		top := uint(x[i])
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		top |= top &gt;&gt; 1
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		top |= top &gt;&gt; 2
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		top |= top &gt;&gt; 4
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		top |= top &gt;&gt; 8
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		top |= top &gt;&gt; 16
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		top |= top &gt;&gt; 16 &gt;&gt; 16 <span class="comment">// &#34;&gt;&gt; 32&#34; doesn&#39;t compile on 32-bit architectures</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>		return i*_W + bits.Len(top)
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	}
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	return 0
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>}
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span><span class="comment">// trailingZeroBits returns the number of consecutive least significant zero</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span><span class="comment">// bits of x.</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>func (x nat) trailingZeroBits() uint {
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	if len(x) == 0 {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		return 0
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	}
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	var i uint
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	for x[i] == 0 {
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		i++
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	}
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	<span class="comment">// x[i] != 0</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	return i*_W + uint(bits.TrailingZeros(uint(x[i])))
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span><span class="comment">// isPow2 returns i, true when x == 2**i and 0, false otherwise.</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>func (x nat) isPow2() (uint, bool) {
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	var i uint
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	for x[i] == 0 {
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		i++
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	}
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	if i == uint(len(x))-1 &amp;&amp; x[i]&amp;(x[i]-1) == 0 {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		return i*_W + uint(bits.TrailingZeros(uint(x[i]))), true
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	}
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	return 0, false
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>}
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>func same(x, y nat) bool {
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	return len(x) == len(y) &amp;&amp; len(x) &gt; 0 &amp;&amp; &amp;x[0] == &amp;y[0]
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>}
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span><span class="comment">// z = x &lt;&lt; s</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>func (z nat) shl(x nat, s uint) nat {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	if s == 0 {
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		if same(z, x) {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>			return z
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>		}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>		if !alias(z, x) {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>			return z.set(x)
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		}
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	m := len(x)
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	if m == 0 {
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>		return z[:0]
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	}
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	<span class="comment">// m &gt; 0</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	n := m + int(s/_W)
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	z = z.make(n + 1)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	z[n] = shlVU(z[n-m:n], x, s%_W)
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	z[0 : n-m].clear()
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	return z.norm()
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>}
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span><span class="comment">// z = x &gt;&gt; s</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>func (z nat) shr(x nat, s uint) nat {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	if s == 0 {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		if same(z, x) {
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>			return z
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>		}
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>		if !alias(z, x) {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>			return z.set(x)
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	}
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	m := len(x)
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	n := m - int(s/_W)
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	if n &lt;= 0 {
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		return z[:0]
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	}
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	<span class="comment">// n &gt; 0</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	z = z.make(n)
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	shrVU(z, x[m-n:], s%_W)
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	return z.norm()
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>}
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>func (z nat) setBit(x nat, i uint, b uint) nat {
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	j := int(i / _W)
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	m := Word(1) &lt;&lt; (i % _W)
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	n := len(x)
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	switch b {
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	case 0:
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		z = z.make(n)
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		copy(z, x)
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		if j &gt;= n {
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			<span class="comment">// no need to grow</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>			return z
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		}
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>		z[j] &amp;^= m
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>		return z.norm()
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	case 1:
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		if j &gt;= n {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			z = z.make(j + 1)
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>			z[n:].clear()
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		} else {
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>			z = z.make(n)
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		copy(z, x)
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>		z[j] |= m
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		<span class="comment">// no need to normalize</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		return z
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	}
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	panic(&#34;set bit is not 0 or 1&#34;)
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span><span class="comment">// bit returns the value of the i&#39;th bit, with lsb == bit 0.</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>func (x nat) bit(i uint) uint {
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	j := i / _W
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	if j &gt;= uint(len(x)) {
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>		return 0
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	}
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	<span class="comment">// 0 &lt;= j &lt; len(x)</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	return uint(x[j] &gt;&gt; (i % _W) &amp; 1)
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>}
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span><span class="comment">// sticky returns 1 if there&#39;s a 1 bit within the</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span><span class="comment">// i least significant bits, otherwise it returns 0.</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>func (x nat) sticky(i uint) uint {
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	j := i / _W
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	if j &gt;= uint(len(x)) {
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		if len(x) == 0 {
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			return 0
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		}
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		return 1
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	}
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	<span class="comment">// 0 &lt;= j &lt; len(x)</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	for _, x := range x[:j] {
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		if x != 0 {
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>			return 1
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		}
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	}
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	if x[j]&lt;&lt;(_W-i%_W) != 0 {
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		return 1
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	}
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	return 0
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>}
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>func (z nat) and(x, y nat) nat {
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	m := len(x)
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>	n := len(y)
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	if m &gt; n {
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>		m = n
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	}
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	<span class="comment">// m &lt;= n</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	z = z.make(m)
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	for i := 0; i &lt; m; i++ {
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		z[i] = x[i] &amp; y[i]
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	}
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	return z.norm()
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>}
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span><span class="comment">// trunc returns z = x mod 2ⁿ.</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>func (z nat) trunc(x nat, n uint) nat {
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	w := (n + _W - 1) / _W
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	if uint(len(x)) &lt; w {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		return z.set(x)
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	}
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	z = z.make(int(w))
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	copy(z, x)
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	if n%_W != 0 {
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		z[len(z)-1] &amp;= 1&lt;&lt;(n%_W) - 1
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	return z.norm()
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>}
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>func (z nat) andNot(x, y nat) nat {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	m := len(x)
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	n := len(y)
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	if n &gt; m {
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		n = m
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	}
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	<span class="comment">// m &gt;= n</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	z = z.make(m)
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		z[i] = x[i] &amp;^ y[i]
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	}
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	copy(z[n:m], x[n:m])
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	return z.norm()
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>}
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>func (z nat) or(x, y nat) nat {
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	m := len(x)
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>	n := len(y)
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	s := x
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	if m &lt; n {
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>		n, m = m, n
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>		s = y
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	}
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	<span class="comment">// m &gt;= n</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	z = z.make(m)
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>		z[i] = x[i] | y[i]
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	}
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	copy(z[n:m], s[n:m])
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	return z.norm()
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>}
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>func (z nat) xor(x, y nat) nat {
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	m := len(x)
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	n := len(y)
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	s := x
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	if m &lt; n {
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>		n, m = m, n
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>		s = y
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>	}
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>	<span class="comment">// m &gt;= n</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	z = z.make(m)
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>		z[i] = x[i] ^ y[i]
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	}
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	copy(z[n:m], s[n:m])
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	return z.norm()
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>}
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span><span class="comment">// random creates a random integer in [0..limit), using the space in z if</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span><span class="comment">// possible. n is the bit length of limit.</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>func (z nat) random(rand *rand.Rand, limit nat, n int) nat {
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	if alias(z, limit) {
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>		z = nil <span class="comment">// z is an alias for limit - cannot reuse</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	}
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	z = z.make(len(limit))
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	bitLengthOfMSW := uint(n % _W)
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	if bitLengthOfMSW == 0 {
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>		bitLengthOfMSW = _W
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	}
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	mask := Word((1 &lt;&lt; bitLengthOfMSW) - 1)
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	for {
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>		switch _W {
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		case 32:
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>			for i := range z {
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>				z[i] = Word(rand.Uint32())
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>			}
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		case 64:
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>			for i := range z {
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>				z[i] = Word(rand.Uint32()) | Word(rand.Uint32())&lt;&lt;32
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>			}
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>		default:
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>			panic(&#34;unknown word size&#34;)
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>		}
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>		z[len(limit)-1] &amp;= mask
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		if z.cmp(limit) &lt; 0 {
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>			break
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>		}
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>	}
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	return z.norm()
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>}
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span><span class="comment">// If m != 0 (i.e., len(m) != 0), expNN sets z to x**y mod m;</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span><span class="comment">// otherwise it sets z to x**y. The result is the value of z.</span>
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>func (z nat) expNN(x, y, m nat, slow bool) nat {
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	if alias(z, x) || alias(z, y) {
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		<span class="comment">// We cannot allow in-place modification of x or y.</span>
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		z = nil
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>	}
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>	<span class="comment">// x**y mod 1 == 0</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>	if len(m) == 1 &amp;&amp; m[0] == 1 {
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>		return z.setWord(0)
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>	}
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>	<span class="comment">// m == 0 || m &gt; 1</span>
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>	<span class="comment">// x**0 == 1</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>	if len(y) == 0 {
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>		return z.setWord(1)
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>	}
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>	<span class="comment">// y &gt; 0</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	<span class="comment">// 0**y = 0</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>	if len(x) == 0 {
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		return z.setWord(0)
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>	}
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	<span class="comment">// x &gt; 0</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>	<span class="comment">// 1**y = 1</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	if len(x) == 1 &amp;&amp; x[0] == 1 {
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>		return z.setWord(1)
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>	}
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>	<span class="comment">// x &gt; 1</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>	<span class="comment">// x**1 == x</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	if len(y) == 1 &amp;&amp; y[0] == 1 {
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>		if len(m) != 0 {
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>			return z.rem(x, m)
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>		}
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>		return z.set(x)
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>	}
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	<span class="comment">// y &gt; 1</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	if len(m) != 0 {
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>		<span class="comment">// We likely end up being as long as the modulus.</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>		z = z.make(len(m))
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		<span class="comment">// If the exponent is large, we use the Montgomery method for odd values,</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		<span class="comment">// and a 4-bit, windowed exponentiation for powers of two,</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>		<span class="comment">// and a CRT-decomposed Montgomery method for the remaining values</span>
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>		<span class="comment">// (even values times non-trivial odd values, which decompose into one</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>		<span class="comment">// instance of each of the first two cases).</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>		if len(y) &gt; 1 &amp;&amp; !slow {
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>			if m[0]&amp;1 == 1 {
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>				return z.expNNMontgomery(x, y, m)
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>			}
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>			if logM, ok := m.isPow2(); ok {
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>				return z.expNNWindowed(x, y, logM)
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>			}
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>			return z.expNNMontgomeryEven(x, y, m)
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		}
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	}
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	z = z.set(x)
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	v := y[len(y)-1] <span class="comment">// v &gt; 0 because y is normalized and y &gt; 0</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	shift := nlz(v) + 1
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>	v &lt;&lt;= shift
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>	var q nat
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>	const mask = 1 &lt;&lt; (_W - 1)
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	<span class="comment">// We walk through the bits of the exponent one by one. Each time we</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	<span class="comment">// see a bit, we square, thus doubling the power. If the bit is a one,</span>
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	<span class="comment">// we also multiply by x, thus adding one to the power.</span>
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>	w := _W - int(shift)
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>	<span class="comment">// zz and r are used to avoid allocating in mul and div as</span>
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	<span class="comment">// otherwise the arguments would alias.</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	var zz, r nat
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	for j := 0; j &lt; w; j++ {
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>		zz = zz.sqr(z)
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>		zz, z = z, zz
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>		if v&amp;mask != 0 {
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>			zz = zz.mul(z, x)
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>			zz, z = z, zz
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		}
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		if len(m) != 0 {
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>			zz, r = zz.div(r, z, m)
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>			zz, r, q, z = q, z, zz, r
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>		}
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>		v &lt;&lt;= 1
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	}
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	for i := len(y) - 2; i &gt;= 0; i-- {
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>		v = y[i]
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>		for j := 0; j &lt; _W; j++ {
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>			zz = zz.sqr(z)
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>			zz, z = z, zz
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>			if v&amp;mask != 0 {
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>				zz = zz.mul(z, x)
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>				zz, z = z, zz
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>			}
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>			if len(m) != 0 {
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>				zz, r = zz.div(r, z, m)
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>				zz, r, q, z = q, z, zz, r
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>			}
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>			v &lt;&lt;= 1
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>		}
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>	}
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>	return z.norm()
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>}
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span><span class="comment">// expNNMontgomeryEven calculates x**y mod m where m = m1 × m2 for m1 = 2ⁿ and m2 odd.</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span><span class="comment">// It uses two recursive calls to expNN for x**y mod m1 and x**y mod m2</span>
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span><span class="comment">// and then uses the Chinese Remainder Theorem to combine the results.</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span><span class="comment">// The recursive call using m1 will use expNNWindowed,</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span><span class="comment">// while the recursive call using m2 will use expNNMontgomery.</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span><span class="comment">// For more details, see Ç. K. Koç, “Montgomery Reduction with Even Modulus”,</span>
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span><span class="comment">// IEE Proceedings: Computers and Digital Techniques, 141(5) 314-316, September 1994.</span>
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span><span class="comment">// http://www.people.vcu.edu/~jwang3/CMSC691/j34monex.pdf</span>
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>func (z nat) expNNMontgomeryEven(x, y, m nat) nat {
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>	<span class="comment">// Split m = m₁ × m₂ where m₁ = 2ⁿ</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>	n := m.trailingZeroBits()
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	m1 := nat(nil).shl(natOne, n)
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>	m2 := nat(nil).shr(m, n)
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>	<span class="comment">// We want z = x**y mod m.</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	<span class="comment">// z₁ = x**y mod m1 = (x**y mod m) mod m1 = z mod m1</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	<span class="comment">// z₂ = x**y mod m2 = (x**y mod m) mod m2 = z mod m2</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>	<span class="comment">// (We are using the math/big convention for names here,</span>
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	<span class="comment">// where the computation is z = x**y mod m, so its parts are z1 and z2.</span>
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	<span class="comment">// The paper is computing x = a**e mod n; it refers to these as x2 and z1.)</span>
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	z1 := nat(nil).expNN(x, y, m1, false)
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>	z2 := nat(nil).expNN(x, y, m2, false)
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>	<span class="comment">// Reconstruct z from z₁, z₂ using CRT, using algorithm from paper,</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	<span class="comment">// which uses only a single modInverse (and an easy one at that).</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>	<span class="comment">//	p = (z₁ - z₂) × m₂⁻¹ (mod m₁)</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>	<span class="comment">//	z = z₂ + p × m₂</span>
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>	<span class="comment">// The final addition is in range because:</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>	<span class="comment">//	z = z₂ + p × m₂</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>	<span class="comment">//	  ≤ z₂ + (m₁-1) × m₂</span>
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>	<span class="comment">//	  &lt; m₂ + (m₁-1) × m₂</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>	<span class="comment">//	  = m₁ × m₂</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	<span class="comment">//	  = m.</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>	z = z.set(z2)
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>	<span class="comment">// Compute (z₁ - z₂) mod m1 [m1 == 2**n] into z1.</span>
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>	z1 = z1.subMod2N(z1, z2, n)
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>	<span class="comment">// Reuse z2 for p = (z₁ - z₂) [in z1] * m2⁻¹ (mod m₁ [= 2ⁿ]).</span>
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>	m2inv := nat(nil).modInverse(m2, m1)
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>	z2 = z2.mul(z1, m2inv)
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>	z2 = z2.trunc(z2, n)
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>	<span class="comment">// Reuse z1 for p * m2.</span>
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	z = z.add(z, z1.mul(z2, m2))
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	return z
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>}
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span><span class="comment">// expNNWindowed calculates x**y mod m using a fixed, 4-bit window,</span>
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span><span class="comment">// where m = 2**logM.</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>func (z nat) expNNWindowed(x, y nat, logM uint) nat {
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>	if len(y) &lt;= 1 {
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>		panic(&#34;big: misuse of expNNWindowed&#34;)
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>	}
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>	if x[0]&amp;1 == 0 {
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>		<span class="comment">// len(y) &gt; 1, so y  &gt; logM.</span>
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>		<span class="comment">// x is even, so x**y is a multiple of 2**y which is a multiple of 2**logM.</span>
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>		return z.setWord(0)
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>	}
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>	if logM == 1 {
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>		return z.setWord(1)
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>	}
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>	<span class="comment">// zz is used to avoid allocating in mul as otherwise</span>
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>	<span class="comment">// the arguments would alias.</span>
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>	w := int((logM + _W - 1) / _W)
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	zzp := getNat(w)
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	zz := *zzp
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>	const n = 4
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>	<span class="comment">// powers[i] contains x^i.</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>	var powers [1 &lt;&lt; n]*nat
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>	for i := range powers {
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>		powers[i] = getNat(w)
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>	}
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>	*powers[0] = powers[0].set(natOne)
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>	*powers[1] = powers[1].trunc(x, logM)
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>	for i := 2; i &lt; 1&lt;&lt;n; i += 2 {
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>		p2, p, p1 := powers[i/2], powers[i], powers[i+1]
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>		*p = p.sqr(*p2)
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>		*p = p.trunc(*p, logM)
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>		*p1 = p1.mul(*p, x)
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>		*p1 = p1.trunc(*p1, logM)
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>	}
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	<span class="comment">// Because phi(2**logM) = 2**(logM-1), x**(2**(logM-1)) = 1,</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	<span class="comment">// so we can compute x**(y mod 2**(logM-1)) instead of x**y.</span>
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	<span class="comment">// That is, we can throw away all but the bottom logM-1 bits of y.</span>
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>	<span class="comment">// Instead of allocating a new y, we start reading y at the right word</span>
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>	<span class="comment">// and truncate it appropriately at the start of the loop.</span>
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>	i := len(y) - 1
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>	mtop := int((logM - 2) / _W) <span class="comment">// -2 because the top word of N bits is the (N-1)/W&#39;th word.</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>	mmask := ^Word(0)
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>	if mbits := (logM - 1) &amp; (_W - 1); mbits != 0 {
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>		mmask = (1 &lt;&lt; mbits) - 1
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>	}
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>	if i &gt; mtop {
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>		i = mtop
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>	}
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>	advance := false
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>	z = z.setWord(1)
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>	for ; i &gt;= 0; i-- {
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>		yi := y[i]
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>		if i == mtop {
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>			yi &amp;= mmask
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>		}
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>		for j := 0; j &lt; _W; j += n {
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>			if advance {
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>				<span class="comment">// Account for use of 4 bits in previous iteration.</span>
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>				<span class="comment">// Unrolled loop for significant performance</span>
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>				<span class="comment">// gain. Use go test -bench=&#34;.*&#34; in crypto/rsa</span>
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>				<span class="comment">// to check performance before making changes.</span>
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>				zz = zz.sqr(z)
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>				zz, z = z, zz
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>				z = z.trunc(z, logM)
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>				zz = zz.sqr(z)
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>				zz, z = z, zz
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>				z = z.trunc(z, logM)
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>				zz = zz.sqr(z)
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>				zz, z = z, zz
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>				z = z.trunc(z, logM)
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>				zz = zz.sqr(z)
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>				zz, z = z, zz
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>				z = z.trunc(z, logM)
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>			}
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>			zz = zz.mul(z, *powers[yi&gt;&gt;(_W-n)])
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>			zz, z = z, zz
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>			z = z.trunc(z, logM)
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>			yi &lt;&lt;= n
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>			advance = true
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>		}
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>	}
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>	*zzp = zz
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>	putNat(zzp)
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>	for i := range powers {
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>		putNat(powers[i])
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>	}
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>	return z.norm()
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>}
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span><span class="comment">// expNNMontgomery calculates x**y mod m using a fixed, 4-bit window.</span>
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span><span class="comment">// Uses Montgomery representation.</span>
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>func (z nat) expNNMontgomery(x, y, m nat) nat {
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>	numWords := len(m)
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>	<span class="comment">// We want the lengths of x and m to be equal.</span>
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>	<span class="comment">// It is OK if x &gt;= m as long as len(x) == len(m).</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>	if len(x) &gt; numWords {
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>		_, x = nat(nil).div(nil, x, m)
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>		<span class="comment">// Note: now len(x) &lt;= numWords, not guaranteed ==.</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>	}
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>	if len(x) &lt; numWords {
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>		rr := make(nat, numWords)
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>		copy(rr, x)
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>		x = rr
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>	}
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>	<span class="comment">// Ideally the precomputations would be performed outside, and reused</span>
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>	<span class="comment">// k0 = -m**-1 mod 2**_W. Algorithm from: Dumas, J.G. &#34;On Newton–Raphson</span>
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>	<span class="comment">// Iteration for Multiplicative Inverses Modulo Prime Powers&#34;.</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>	k0 := 2 - m[0]
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>	t := m[0] - 1
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>	for i := 1; i &lt; _W; i &lt;&lt;= 1 {
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>		t *= t
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>		k0 *= (t + 1)
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>	}
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>	k0 = -k0
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	<span class="comment">// RR = 2**(2*_W*len(m)) mod m</span>
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>	RR := nat(nil).setWord(1)
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>	zz := nat(nil).shl(RR, uint(2*numWords*_W))
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>	_, RR = nat(nil).div(RR, zz, m)
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>	if len(RR) &lt; numWords {
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>		zz = zz.make(numWords)
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>		copy(zz, RR)
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>		RR = zz
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>	}
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>	<span class="comment">// one = 1, with equal length to that of m</span>
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>	one := make(nat, numWords)
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>	one[0] = 1
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>	const n = 4
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>	<span class="comment">// powers[i] contains x^i</span>
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>	var powers [1 &lt;&lt; n]nat
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>	powers[0] = powers[0].montgomery(one, RR, m, k0, numWords)
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>	powers[1] = powers[1].montgomery(x, RR, m, k0, numWords)
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>	for i := 2; i &lt; 1&lt;&lt;n; i++ {
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>		powers[i] = powers[i].montgomery(powers[i-1], powers[1], m, k0, numWords)
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>	}
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>	<span class="comment">// initialize z = 1 (Montgomery 1)</span>
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>	z = z.make(numWords)
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>	copy(z, powers[0])
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>	zz = zz.make(numWords)
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>	<span class="comment">// same windowed exponent, but with Montgomery multiplications</span>
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>	for i := len(y) - 1; i &gt;= 0; i-- {
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>		yi := y[i]
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>		for j := 0; j &lt; _W; j += n {
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>			if i != len(y)-1 || j != 0 {
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>				zz = zz.montgomery(z, z, m, k0, numWords)
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>				z = z.montgomery(zz, zz, m, k0, numWords)
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>				zz = zz.montgomery(z, z, m, k0, numWords)
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>				z = z.montgomery(zz, zz, m, k0, numWords)
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>			}
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>			zz = zz.montgomery(z, powers[yi&gt;&gt;(_W-n)], m, k0, numWords)
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>			z, zz = zz, z
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>			yi &lt;&lt;= n
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>		}
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	}
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>	<span class="comment">// convert to regular number</span>
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>	zz = zz.montgomery(z, one, m, k0, numWords)
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>	<span class="comment">// One last reduction, just in case.</span>
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>	<span class="comment">// See golang.org/issue/13907.</span>
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>	if zz.cmp(m) &gt;= 0 {
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>		<span class="comment">// Common case is m has high bit set; in that case,</span>
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>		<span class="comment">// since zz is the same length as m, there can be just</span>
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>		<span class="comment">// one multiple of m to remove. Just subtract.</span>
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>		<span class="comment">// We think that the subtract should be sufficient in general,</span>
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>		<span class="comment">// so do that unconditionally, but double-check,</span>
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>		<span class="comment">// in case our beliefs are wrong.</span>
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>		<span class="comment">// The div is not expected to be reached.</span>
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>		zz = zz.sub(zz, m)
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>		if zz.cmp(m) &gt;= 0 {
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>			_, zz = nat(nil).div(nil, zz, m)
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>		}
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>	}
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>	return zz.norm()
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>}
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span><span class="comment">// bytes writes the value of z into buf using big-endian encoding.</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span><span class="comment">// The value of z is encoded in the slice buf[i:]. If the value of z</span>
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span><span class="comment">// cannot be represented in buf, bytes panics. The number i of unused</span>
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span><span class="comment">// bytes at the beginning of buf is returned as result.</span>
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>func (z nat) bytes(buf []byte) (i int) {
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>	<span class="comment">// This function is used in cryptographic operations. It must not leak</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>	<span class="comment">// anything but the Int&#39;s sign and bit size through side-channels. Any</span>
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>	<span class="comment">// changes must be reviewed by a security expert.</span>
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>	i = len(buf)
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>	for _, d := range z {
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>		for j := 0; j &lt; _S; j++ {
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>			i--
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>			if i &gt;= 0 {
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>				buf[i] = byte(d)
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>			} else if byte(d) != 0 {
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>				panic(&#34;math/big: buffer too small to fit value&#34;)
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>			}
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>			d &gt;&gt;= 8
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>		}
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>	}
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>		i = 0
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>	}
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>	for i &lt; len(buf) &amp;&amp; buf[i] == 0 {
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>		i++
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>	}
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>	return
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>}
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span><span class="comment">// bigEndianWord returns the contents of buf interpreted as a big-endian encoded Word value.</span>
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>func bigEndianWord(buf []byte) Word {
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>	if _W == 64 {
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>		return Word(binary.BigEndian.Uint64(buf))
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>	}
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>	return Word(binary.BigEndian.Uint32(buf))
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>}
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span><span class="comment">// setBytes interprets buf as the bytes of a big-endian unsigned</span>
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span><span class="comment">// integer, sets z to that value, and returns z.</span>
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>func (z nat) setBytes(buf []byte) nat {
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>	z = z.make((len(buf) + _S - 1) / _S)
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>	i := len(buf)
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>	for k := 0; i &gt;= _S; k++ {
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>		z[k] = bigEndianWord(buf[i-_S : i])
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>		i -= _S
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>	}
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>	if i &gt; 0 {
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>		var d Word
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>		for s := uint(0); i &gt; 0; s += 8 {
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>			d |= Word(buf[i-1]) &lt;&lt; s
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>			i--
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>		}
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>		z[len(z)-1] = d
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>	}
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>	return z.norm()
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>}
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span><span class="comment">// sqrt sets z = ⌊√x⌋</span>
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>func (z nat) sqrt(x nat) nat {
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>	if x.cmp(natOne) &lt;= 0 {
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>		return z.set(x)
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>	}
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>	if alias(z, x) {
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>		z = nil
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>	}
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>	<span class="comment">// Start with value known to be too large and repeat &#34;z = ⌊(z + ⌊x/z⌋)/2⌋&#34; until it stops getting smaller.</span>
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>	<span class="comment">// See Brent and Zimmermann, Modern Computer Arithmetic, Algorithm 1.13 (SqrtInt).</span>
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>	<span class="comment">// https://members.loria.fr/PZimmermann/mca/pub226.html</span>
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>	<span class="comment">// If x is one less than a perfect square, the sequence oscillates between the correct z and z+1;</span>
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>	<span class="comment">// otherwise it converges to the correct z and stays there.</span>
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>	var z1, z2 nat
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>	z1 = z
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>	z1 = z1.setUint64(1)
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>	z1 = z1.shl(z1, uint(x.bitLen()+1)/2) <span class="comment">// must be ≥ √x</span>
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>	for n := 0; ; n++ {
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>		z2, _ = z2.div(nil, x, z1)
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>		z2 = z2.add(z2, z1)
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>		z2 = z2.shr(z2, 1)
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>		if z2.cmp(z1) &gt;= 0 {
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>			<span class="comment">// z1 is answer.</span>
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>			<span class="comment">// Figure out whether z1 or z2 is currently aliased to z by looking at loop count.</span>
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>			if n&amp;1 == 0 {
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>				return z1
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>			}
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>			return z.set(z1)
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>		}
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>		z1, z2 = z2, z1
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>	}
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>}
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span><span class="comment">// subMod2N returns z = (x - y) mod 2ⁿ.</span>
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>func (z nat) subMod2N(x, y nat, n uint) nat {
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>	if uint(x.bitLen()) &gt; n {
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>		if alias(z, x) {
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>			<span class="comment">// ok to overwrite x in place</span>
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>			x = x.trunc(x, n)
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>		} else {
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>			x = nat(nil).trunc(x, n)
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>		}
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>	}
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>	if uint(y.bitLen()) &gt; n {
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>		if alias(z, y) {
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>			<span class="comment">// ok to overwrite y in place</span>
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>			y = y.trunc(y, n)
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>		} else {
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>			y = nat(nil).trunc(y, n)
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>		}
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>	}
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>	if x.cmp(y) &gt;= 0 {
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>		return z.sub(x, y)
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>	}
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>	<span class="comment">// x - y &lt; 0; x - y mod 2ⁿ = x - y + 2ⁿ = 2ⁿ - (y - x) = 1 + 2ⁿ-1 - (y - x) = 1 + ^(y - x).</span>
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>	z = z.sub(y, x)
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>	for uint(len(z))*_W &lt; n {
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>		z = append(z, 0)
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>	}
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>	for i := range z {
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>		z[i] = ^z[i]
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>	}
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>	z = z.trunc(z, n)
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>	return z.add(z, natOne)
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>}
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>
</pre><p><a href="nat.go?m=text">View as plain text</a></p>

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
