<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/strconv/decimal.go - Go Documentation Server</title>

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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/strconv">strconv</a>/<span class="text-muted">decimal.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/strconv">strconv</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Multiprecision decimal numbers.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// For floating-point formatting only; not general purpose.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Only operations are assign and (binary) left/right shift.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// Can do binary floating point in multiprecision decimal precisely</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// because 2 divides 10; cannot do decimal floating point</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// in multiprecision binary precisely.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>package strconv
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>type decimal struct {
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	d     [800]byte <span class="comment">// digits, big-endian representation</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	nd    int       <span class="comment">// number of digits used</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	dp    int       <span class="comment">// decimal point</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	neg   bool      <span class="comment">// negative flag</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	trunc bool      <span class="comment">// discarded nonzero digits beyond d[:nd]</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>}
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>func (a *decimal) String() string {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	n := 10 + a.nd
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	if a.dp &gt; 0 {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		n += a.dp
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	}
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	if a.dp &lt; 0 {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		n += -a.dp
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	buf := make([]byte, n)
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	w := 0
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	switch {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	case a.nd == 0:
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		return &#34;0&#34;
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	case a.dp &lt;= 0:
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		<span class="comment">// zeros fill space between decimal point and digits</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		buf[w] = &#39;0&#39;
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		w++
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		buf[w] = &#39;.&#39;
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		w++
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		w += digitZero(buf[w : w+-a.dp])
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		w += copy(buf[w:], a.d[0:a.nd])
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	case a.dp &lt; a.nd:
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		<span class="comment">// decimal point in middle of digits</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		w += copy(buf[w:], a.d[0:a.dp])
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		buf[w] = &#39;.&#39;
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		w++
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		w += copy(buf[w:], a.d[a.dp:a.nd])
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	default:
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		<span class="comment">// zeros fill space between digits and decimal point</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		w += copy(buf[w:], a.d[0:a.nd])
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		w += digitZero(buf[w : w+a.dp-a.nd])
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	return string(buf[0:w])
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>func digitZero(dst []byte) int {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	for i := range dst {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		dst[i] = &#39;0&#39;
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	return len(dst)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// trim trailing zeros from number.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// (They are meaningless; the decimal point is tracked</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// independent of the number of digits.)</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>func trim(a *decimal) {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	for a.nd &gt; 0 &amp;&amp; a.d[a.nd-1] == &#39;0&#39; {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		a.nd--
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	if a.nd == 0 {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		a.dp = 0
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// Assign v to a.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>func (a *decimal) Assign(v uint64) {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	var buf [24]byte
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// Write reversed decimal in buf.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	n := 0
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	for v &gt; 0 {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		v1 := v / 10
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		v -= 10 * v1
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		buf[n] = byte(v + &#39;0&#39;)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		n++
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		v = v1
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// Reverse again to produce forward decimal in a.d.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	a.nd = 0
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	for n--; n &gt;= 0; n-- {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		a.d[a.nd] = buf[n]
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		a.nd++
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	a.dp = a.nd
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	trim(a)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// Maximum shift that we can do in one pass without overflow.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// A uint has 32 or 64 bits, and we have to be able to accommodate 9&lt;&lt;k.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>const uintSize = 32 &lt;&lt; (^uint(0) &gt;&gt; 63)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>const maxShift = uintSize - 4
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// Binary shift right (/ 2) by k bits.  k &lt;= maxShift to avoid overflow.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>func rightShift(a *decimal, k uint) {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	r := 0 <span class="comment">// read pointer</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	w := 0 <span class="comment">// write pointer</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// Pick up enough leading digits to cover first shift.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	var n uint
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	for ; n&gt;&gt;k == 0; r++ {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		if r &gt;= a.nd {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			if n == 0 {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>				<span class="comment">// a == 0; shouldn&#39;t get here, but handle anyway.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>				a.nd = 0
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>				return
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			for n&gt;&gt;k == 0 {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>				n = n * 10
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>				r++
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			break
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		c := uint(a.d[r])
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		n = n*10 + c - &#39;0&#39;
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	a.dp -= r - 1
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	var mask uint = (1 &lt;&lt; k) - 1
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// Pick up a digit, put down a digit.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	for ; r &lt; a.nd; r++ {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		c := uint(a.d[r])
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		dig := n &gt;&gt; k
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		n &amp;= mask
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		a.d[w] = byte(dig + &#39;0&#39;)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		w++
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		n = n*10 + c - &#39;0&#39;
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// Put down extra digits.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	for n &gt; 0 {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		dig := n &gt;&gt; k
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		n &amp;= mask
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		if w &lt; len(a.d) {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			a.d[w] = byte(dig + &#39;0&#39;)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			w++
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		} else if dig &gt; 0 {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			a.trunc = true
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		n = n * 10
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	a.nd = w
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	trim(a)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// Cheat sheet for left shift: table indexed by shift count giving</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// number of new digits that will be introduced by that shift.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// For example, leftcheats[4] = {2, &#34;625&#34;}.  That means that</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// if we are shifting by 4 (multiplying by 16), it will add 2 digits</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span><span class="comment">// when the string prefix is &#34;625&#34; through &#34;999&#34;, and one fewer digit</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// if the string prefix is &#34;000&#34; through &#34;624&#34;.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// Credit for this trick goes to Ken.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>type leftCheat struct {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	delta  int    <span class="comment">// number of new digits</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	cutoff string <span class="comment">// minus one digit if original &lt; a.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>var leftcheats = []leftCheat{
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// Leading digits of 1/2^i = 5^i.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// 5^23 is not an exact 64-bit floating point number,</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// so have to use bc for the math.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// Go up to 60 to be large enough for 32bit and 64bit platforms.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">/*
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		seq 60 | sed &#39;s/^/5^/&#39; | bc |
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		awk &#39;BEGIN{ print &#34;\t{ 0, \&#34;\&#34; },&#34; }
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		{
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			log2 = log(2)/log(10)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			printf(&#34;\t{ %d, \&#34;%s\&#34; },\t// * %d\n&#34;,
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				int(log2*NR+1), $0, 2**NR)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		}&#39;
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	*/</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	{0, &#34;&#34;},
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	{1, &#34;5&#34;},                                           <span class="comment">// * 2</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	{1, &#34;25&#34;},                                          <span class="comment">// * 4</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	{1, &#34;125&#34;},                                         <span class="comment">// * 8</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	{2, &#34;625&#34;},                                         <span class="comment">// * 16</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	{2, &#34;3125&#34;},                                        <span class="comment">// * 32</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	{2, &#34;15625&#34;},                                       <span class="comment">// * 64</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	{3, &#34;78125&#34;},                                       <span class="comment">// * 128</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	{3, &#34;390625&#34;},                                      <span class="comment">// * 256</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	{3, &#34;1953125&#34;},                                     <span class="comment">// * 512</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	{4, &#34;9765625&#34;},                                     <span class="comment">// * 1024</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	{4, &#34;48828125&#34;},                                    <span class="comment">// * 2048</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	{4, &#34;244140625&#34;},                                   <span class="comment">// * 4096</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	{4, &#34;1220703125&#34;},                                  <span class="comment">// * 8192</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	{5, &#34;6103515625&#34;},                                  <span class="comment">// * 16384</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	{5, &#34;30517578125&#34;},                                 <span class="comment">// * 32768</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	{5, &#34;152587890625&#34;},                                <span class="comment">// * 65536</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	{6, &#34;762939453125&#34;},                                <span class="comment">// * 131072</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	{6, &#34;3814697265625&#34;},                               <span class="comment">// * 262144</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	{6, &#34;19073486328125&#34;},                              <span class="comment">// * 524288</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	{7, &#34;95367431640625&#34;},                              <span class="comment">// * 1048576</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	{7, &#34;476837158203125&#34;},                             <span class="comment">// * 2097152</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	{7, &#34;2384185791015625&#34;},                            <span class="comment">// * 4194304</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	{7, &#34;11920928955078125&#34;},                           <span class="comment">// * 8388608</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	{8, &#34;59604644775390625&#34;},                           <span class="comment">// * 16777216</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	{8, &#34;298023223876953125&#34;},                          <span class="comment">// * 33554432</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	{8, &#34;1490116119384765625&#34;},                         <span class="comment">// * 67108864</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	{9, &#34;7450580596923828125&#34;},                         <span class="comment">// * 134217728</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	{9, &#34;37252902984619140625&#34;},                        <span class="comment">// * 268435456</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	{9, &#34;186264514923095703125&#34;},                       <span class="comment">// * 536870912</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	{10, &#34;931322574615478515625&#34;},                      <span class="comment">// * 1073741824</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	{10, &#34;4656612873077392578125&#34;},                     <span class="comment">// * 2147483648</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	{10, &#34;23283064365386962890625&#34;},                    <span class="comment">// * 4294967296</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	{10, &#34;116415321826934814453125&#34;},                   <span class="comment">// * 8589934592</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	{11, &#34;582076609134674072265625&#34;},                   <span class="comment">// * 17179869184</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	{11, &#34;2910383045673370361328125&#34;},                  <span class="comment">// * 34359738368</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	{11, &#34;14551915228366851806640625&#34;},                 <span class="comment">// * 68719476736</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	{12, &#34;72759576141834259033203125&#34;},                 <span class="comment">// * 137438953472</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	{12, &#34;363797880709171295166015625&#34;},                <span class="comment">// * 274877906944</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	{12, &#34;1818989403545856475830078125&#34;},               <span class="comment">// * 549755813888</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	{13, &#34;9094947017729282379150390625&#34;},               <span class="comment">// * 1099511627776</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	{13, &#34;45474735088646411895751953125&#34;},              <span class="comment">// * 2199023255552</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	{13, &#34;227373675443232059478759765625&#34;},             <span class="comment">// * 4398046511104</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	{13, &#34;1136868377216160297393798828125&#34;},            <span class="comment">// * 8796093022208</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	{14, &#34;5684341886080801486968994140625&#34;},            <span class="comment">// * 17592186044416</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	{14, &#34;28421709430404007434844970703125&#34;},           <span class="comment">// * 35184372088832</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	{14, &#34;142108547152020037174224853515625&#34;},          <span class="comment">// * 70368744177664</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	{15, &#34;710542735760100185871124267578125&#34;},          <span class="comment">// * 140737488355328</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	{15, &#34;3552713678800500929355621337890625&#34;},         <span class="comment">// * 281474976710656</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	{15, &#34;17763568394002504646778106689453125&#34;},        <span class="comment">// * 562949953421312</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	{16, &#34;88817841970012523233890533447265625&#34;},        <span class="comment">// * 1125899906842624</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	{16, &#34;444089209850062616169452667236328125&#34;},       <span class="comment">// * 2251799813685248</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	{16, &#34;2220446049250313080847263336181640625&#34;},      <span class="comment">// * 4503599627370496</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	{16, &#34;11102230246251565404236316680908203125&#34;},     <span class="comment">// * 9007199254740992</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	{17, &#34;55511151231257827021181583404541015625&#34;},     <span class="comment">// * 18014398509481984</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	{17, &#34;277555756156289135105907917022705078125&#34;},    <span class="comment">// * 36028797018963968</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	{17, &#34;1387778780781445675529539585113525390625&#34;},   <span class="comment">// * 72057594037927936</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	{18, &#34;6938893903907228377647697925567626953125&#34;},   <span class="comment">// * 144115188075855872</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	{18, &#34;34694469519536141888238489627838134765625&#34;},  <span class="comment">// * 288230376151711744</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	{18, &#34;173472347597680709441192448139190673828125&#34;}, <span class="comment">// * 576460752303423488</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	{19, &#34;867361737988403547205962240695953369140625&#34;}, <span class="comment">// * 1152921504606846976</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// Is the leading prefix of b lexicographically less than s?</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>func prefixIsLessThan(b []byte, s string) bool {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i++ {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		if i &gt;= len(b) {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			return true
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		if b[i] != s[i] {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			return b[i] &lt; s[i]
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	return false
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">// Binary shift left (* 2) by k bits.  k &lt;= maxShift to avoid overflow.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>func leftShift(a *decimal, k uint) {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	delta := leftcheats[k].delta
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	if prefixIsLessThan(a.d[0:a.nd], leftcheats[k].cutoff) {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		delta--
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	r := a.nd         <span class="comment">// read index</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	w := a.nd + delta <span class="comment">// write index</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	<span class="comment">// Pick up a digit, put down a digit.</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	var n uint
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	for r--; r &gt;= 0; r-- {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		n += (uint(a.d[r]) - &#39;0&#39;) &lt;&lt; k
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		quo := n / 10
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		rem := n - 10*quo
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		w--
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		if w &lt; len(a.d) {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			a.d[w] = byte(rem + &#39;0&#39;)
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		} else if rem != 0 {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			a.trunc = true
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		n = quo
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// Put down extra digits.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	for n &gt; 0 {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		quo := n / 10
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		rem := n - 10*quo
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		w--
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		if w &lt; len(a.d) {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			a.d[w] = byte(rem + &#39;0&#39;)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		} else if rem != 0 {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			a.trunc = true
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		n = quo
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	a.nd += delta
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if a.nd &gt;= len(a.d) {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		a.nd = len(a.d)
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	a.dp += delta
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	trim(a)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">// Binary shift left (k &gt; 0) or right (k &lt; 0).</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>func (a *decimal) Shift(k int) {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	switch {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	case a.nd == 0:
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		<span class="comment">// nothing to do: a == 0</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	case k &gt; 0:
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		for k &gt; maxShift {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			leftShift(a, maxShift)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>			k -= maxShift
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		leftShift(a, uint(k))
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	case k &lt; 0:
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		for k &lt; -maxShift {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			rightShift(a, maxShift)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			k += maxShift
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		rightShift(a, uint(-k))
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">// If we chop a at nd digits, should we round up?</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>func shouldRoundUp(a *decimal, nd int) bool {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	if nd &lt; 0 || nd &gt;= a.nd {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		return false
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	if a.d[nd] == &#39;5&#39; &amp;&amp; nd+1 == a.nd { <span class="comment">// exactly halfway - round to even</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		<span class="comment">// if we truncated, a little higher than what&#39;s recorded - always round up</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		if a.trunc {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			return true
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		return nd &gt; 0 &amp;&amp; (a.d[nd-1]-&#39;0&#39;)%2 != 0
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	<span class="comment">// not halfway - digit tells all</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	return a.d[nd] &gt;= &#39;5&#39;
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// Round a to nd digits (or fewer).</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// If nd is zero, it means we&#39;re rounding</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// just to the left of the digits, as in</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// 0.09 -&gt; 0.1.</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>func (a *decimal) Round(nd int) {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	if nd &lt; 0 || nd &gt;= a.nd {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		return
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	if shouldRoundUp(a, nd) {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		a.RoundUp(nd)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	} else {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		a.RoundDown(nd)
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// Round a down to nd digits (or fewer).</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>func (a *decimal) RoundDown(nd int) {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	if nd &lt; 0 || nd &gt;= a.nd {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		return
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	}
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	a.nd = nd
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	trim(a)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span><span class="comment">// Round a up to nd digits (or fewer).</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>func (a *decimal) RoundUp(nd int) {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	if nd &lt; 0 || nd &gt;= a.nd {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		return
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	<span class="comment">// round up</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	for i := nd - 1; i &gt;= 0; i-- {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		c := a.d[i]
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		if c &lt; &#39;9&#39; { <span class="comment">// can stop after this digit</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			a.d[i]++
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>			a.nd = i + 1
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			return
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">// Number is all 9s.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	<span class="comment">// Change to single 1 with adjusted decimal point.</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	a.d[0] = &#39;1&#39;
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	a.nd = 1
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	a.dp++
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>}
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span><span class="comment">// Extract integer part, rounded appropriately.</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span><span class="comment">// No guarantees about overflow.</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>func (a *decimal) RoundedInteger() uint64 {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	if a.dp &gt; 20 {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		return 0xFFFFFFFFFFFFFFFF
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	var i int
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	n := uint64(0)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	for i = 0; i &lt; a.dp &amp;&amp; i &lt; a.nd; i++ {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		n = n*10 + uint64(a.d[i]-&#39;0&#39;)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	for ; i &lt; a.dp; i++ {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		n *= 10
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	if shouldRoundUp(a, a.dp) {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		n++
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	return n
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>
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
