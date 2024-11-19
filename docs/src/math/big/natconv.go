<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/big/natconv.go - Go Documentation Server</title>

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
<a href="natconv.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<a href="http://localhost:8080/src/math/big">big</a>/<span class="text-muted">natconv.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements nat-to-string conversion functions.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package big
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;math/bits&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>const digits = &#34;0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// Note: MaxBase = len(digits), but it must remain an untyped rune constant</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//       for API compatibility.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// MaxBase is the largest number base accepted for string conversions.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>const MaxBase = 10 + (&#39;z&#39; - &#39;a&#39; + 1) + (&#39;Z&#39; - &#39;A&#39; + 1)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>const maxBaseSmall = 10 + (&#39;z&#39; - &#39;a&#39; + 1)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// maxPow returns (b**n, n) such that b**n is the largest power b**n &lt;= _M.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// For instance maxPow(10) == (1e19, 19) for 19 decimal digits in a 64bit Word.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// In other words, at most n digits in base b fit into a Word.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// TODO(gri) replace this with a table, generated at build time.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>func maxPow(b Word) (p Word, n int) {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	p, n = b, 1 <span class="comment">// assuming b &lt;= _M</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	for max := _M / b; p &lt;= max; {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		<span class="comment">// p == b**n &amp;&amp; p &lt;= max</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		p *= b
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		n++
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// p == b**n &amp;&amp; p &lt;= _M</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	return
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// pow returns x**n for n &gt; 0, and 1 otherwise.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func pow(x Word, n int) (p Word) {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// n == sum of bi * 2**i, for 0 &lt;= i &lt; imax, and bi is 0 or 1</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// thus x**n == product of x**(2**i) for all i where bi == 1</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// (Russian Peasant Method for exponentiation)</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	p = 1
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	for n &gt; 0 {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		if n&amp;1 != 0 {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			p *= x
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		x *= x
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		n &gt;&gt;= 1
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	return
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// scan errors</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>var (
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	errNoDigits = errors.New(&#34;number has no digits&#34;)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	errInvalSep = errors.New(&#34;&#39;_&#39; must separate successive digits&#34;)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// scan scans the number corresponding to the longest possible prefix</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// from r representing an unsigned number in a given conversion base.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// scan returns the corresponding natural number res, the actual base b,</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// a digit count, and a read or syntax error err, if any.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// For base 0, an underscore character “_” may appear between a base</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// prefix and an adjacent digit, and between successive digits; such</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// underscores do not change the value of the number, or the returned</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// digit count. Incorrect placement of underscores is reported as an</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// error if there are no other errors. If base != 0, underscores are</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// not recognized and thus terminate scanning like any other character</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// that is not a valid radix point or digit.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//	number    = mantissa | prefix pmantissa .</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//	prefix    = &#34;0&#34; [ &#34;b&#34; | &#34;B&#34; | &#34;o&#34; | &#34;O&#34; | &#34;x&#34; | &#34;X&#34; ] .</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//	mantissa  = digits &#34;.&#34; [ digits ] | digits | &#34;.&#34; digits .</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//	pmantissa = [ &#34;_&#34; ] digits &#34;.&#34; [ digits ] | [ &#34;_&#34; ] digits | &#34;.&#34; digits .</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//	digits    = digit { [ &#34;_&#34; ] digit } .</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//	digit     = &#34;0&#34; ... &#34;9&#34; | &#34;a&#34; ... &#34;z&#34; | &#34;A&#34; ... &#34;Z&#34; .</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// Unless fracOk is set, the base argument must be 0 or a value between</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// 2 and MaxBase. If fracOk is set, the base argument must be one of</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// 0, 2, 8, 10, or 16. Providing an invalid base argument leads to a run-</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// time panic.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// For base 0, the number prefix determines the actual base: A prefix of</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// “0b” or “0B” selects base 2, “0o” or “0O” selects base 8, and</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// “0x” or “0X” selects base 16. If fracOk is false, a “0” prefix</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// (immediately followed by digits) selects base 8 as well. Otherwise,</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// the selected base is 10 and no prefix is accepted.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// If fracOk is set, a period followed by a fractional part is permitted.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// The result value is computed as if there were no period present; and</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// the count value is used to determine the fractional part.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// For bases &lt;= 36, lower and upper case letters are considered the same:</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// The letters &#39;a&#39; to &#39;z&#39; and &#39;A&#39; to &#39;Z&#39; represent digit values 10 to 35.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// For bases &gt; 36, the upper case letters &#39;A&#39; to &#39;Z&#39; represent the digit</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// values 36 to 61.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// A result digit count &gt; 0 corresponds to the number of (non-prefix) digits</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// parsed. A digit count &lt;= 0 indicates the presence of a period (if fracOk</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// is set, only), and -count is the number of fractional digits found.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// In this case, the actual value of the scanned number is res * b**count.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>func (z nat) scan(r io.ByteScanner, base int, fracOk bool) (res nat, b, count int, err error) {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// reject invalid bases</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	baseOk := base == 0 ||
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		!fracOk &amp;&amp; 2 &lt;= base &amp;&amp; base &lt;= MaxBase ||
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		fracOk &amp;&amp; (base == 2 || base == 8 || base == 10 || base == 16)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	if !baseOk {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		panic(fmt.Sprintf(&#34;invalid number base %d&#34;, base))
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// prev encodes the previously seen char: it is one</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// of &#39;_&#39;, &#39;0&#39; (a digit), or &#39;.&#39; (anything else). A</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// valid separator &#39;_&#39; may only occur after a digit</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// and if base == 0.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	prev := &#39;.&#39;
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	invalSep := false
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// one char look-ahead</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	ch, err := r.ReadByte()
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// determine actual base</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	b, prefix := base, 0
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if base == 0 {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		<span class="comment">// actual base is 10 unless there&#39;s a base prefix</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		b = 10
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		if err == nil &amp;&amp; ch == &#39;0&#39; {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			prev = &#39;0&#39;
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			count = 1
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			ch, err = r.ReadByte()
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			if err == nil {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>				<span class="comment">// possibly one of 0b, 0B, 0o, 0O, 0x, 0X</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>				switch ch {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>				case &#39;b&#39;, &#39;B&#39;:
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>					b, prefix = 2, &#39;b&#39;
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>				case &#39;o&#39;, &#39;O&#39;:
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>					b, prefix = 8, &#39;o&#39;
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>				case &#39;x&#39;, &#39;X&#39;:
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>					b, prefix = 16, &#39;x&#39;
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>				default:
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>					if !fracOk {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>						b, prefix = 8, &#39;0&#39;
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>					}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>				}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>				if prefix != 0 {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>					count = 0 <span class="comment">// prefix is not counted</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>					if prefix != &#39;0&#39; {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>						ch, err = r.ReadByte()
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>					}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>				}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// convert string</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// Algorithm: Collect digits in groups of at most n digits in di</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// and then use mulAddWW for every such group to add them to the</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// result.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	z = z[:0]
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	b1 := Word(b)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	bn, n := maxPow(b1) <span class="comment">// at most n digits in base b1 fit into Word</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	di := Word(0)       <span class="comment">// 0 &lt;= di &lt; b1**i &lt; bn</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	i := 0              <span class="comment">// 0 &lt;= i &lt; n</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	dp := -1            <span class="comment">// position of decimal point</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	for err == nil {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		if ch == &#39;.&#39; &amp;&amp; fracOk {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			fracOk = false
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			if prev == &#39;_&#39; {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>				invalSep = true
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			prev = &#39;.&#39;
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			dp = count
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		} else if ch == &#39;_&#39; &amp;&amp; base == 0 {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			if prev != &#39;0&#39; {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>				invalSep = true
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			prev = &#39;_&#39;
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		} else {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			<span class="comment">// convert rune into digit value d1</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			var d1 Word
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			switch {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			case &#39;0&#39; &lt;= ch &amp;&amp; ch &lt;= &#39;9&#39;:
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				d1 = Word(ch - &#39;0&#39;)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			case &#39;a&#39; &lt;= ch &amp;&amp; ch &lt;= &#39;z&#39;:
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>				d1 = Word(ch - &#39;a&#39; + 10)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			case &#39;A&#39; &lt;= ch &amp;&amp; ch &lt;= &#39;Z&#39;:
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>				if b &lt;= maxBaseSmall {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>					d1 = Word(ch - &#39;A&#39; + 10)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>				} else {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>					d1 = Word(ch - &#39;A&#39; + maxBaseSmall)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>				}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			default:
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>				d1 = MaxBase + 1
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			if d1 &gt;= b1 {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>				r.UnreadByte() <span class="comment">// ch does not belong to number anymore</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>				break
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			prev = &#39;0&#39;
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			count++
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			<span class="comment">// collect d1 in di</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			di = di*b1 + d1
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			i++
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			<span class="comment">// if di is &#34;full&#34;, add it to the result</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			if i == n {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>				z = z.mulAddWW(z, bn, di)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>				di = 0
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>				i = 0
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		ch, err = r.ReadByte()
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	if err == io.EOF {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		err = nil
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// other errors take precedence over invalid separators</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	if err == nil &amp;&amp; (invalSep || prev == &#39;_&#39;) {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		err = errInvalSep
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	if count == 0 {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		<span class="comment">// no digits found</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		if prefix == &#39;0&#39; {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			<span class="comment">// there was only the octal prefix 0 (possibly followed by separators and digits &gt; 7);</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			<span class="comment">// interpret as decimal 0</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			return z[:0], 10, 1, err
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		err = errNoDigits <span class="comment">// fall through; result will be 0</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// add remaining digits to result</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if i &gt; 0 {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		z = z.mulAddWW(z, pow(b1, i), di)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	res = z.norm()
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	<span class="comment">// adjust count for fraction, if any</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	if dp &gt;= 0 {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		<span class="comment">// 0 &lt;= dp &lt;= count</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		count = dp - count
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	return
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// utoa converts x to an ASCII representation in the given base;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span><span class="comment">// base must be between 2 and MaxBase, inclusive.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>func (x nat) utoa(base int) []byte {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	return x.itoa(false, base)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span><span class="comment">// itoa is like utoa but it prepends a &#39;-&#39; if neg &amp;&amp; x != 0.</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>func (x nat) itoa(neg bool, base int) []byte {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	if base &lt; 2 || base &gt; MaxBase {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		panic(&#34;invalid base&#34;)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// x == 0</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	if len(x) == 0 {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		return []byte(&#34;0&#34;)
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	<span class="comment">// len(x) &gt; 0</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	<span class="comment">// allocate buffer for conversion</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	i := int(float64(x.bitLen())/math.Log2(float64(base))) + 1 <span class="comment">// off by 1 at most</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	if neg {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		i++
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	s := make([]byte, i)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	<span class="comment">// convert power of two and non power of two bases separately</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	if b := Word(base); b == b&amp;-b {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		<span class="comment">// shift is base b digit size in bits</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		shift := uint(bits.TrailingZeros(uint(b))) <span class="comment">// shift &gt; 0 because b &gt;= 2</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		mask := Word(1&lt;&lt;shift - 1)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		w := x[0]         <span class="comment">// current word</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		nbits := uint(_W) <span class="comment">// number of unprocessed bits in w</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		<span class="comment">// convert less-significant words (include leading zeros)</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		for k := 1; k &lt; len(x); k++ {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			<span class="comment">// convert full digits</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			for nbits &gt;= shift {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>				i--
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>				s[i] = digits[w&amp;mask]
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>				w &gt;&gt;= shift
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>				nbits -= shift
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			<span class="comment">// convert any partial leading digit and advance to next word</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			if nbits == 0 {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>				<span class="comment">// no partial digit remaining, just advance</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>				w = x[k]
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>				nbits = _W
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			} else {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>				<span class="comment">// partial digit in current word w (== x[k-1]) and next word x[k]</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>				w |= x[k] &lt;&lt; nbits
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>				i--
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>				s[i] = digits[w&amp;mask]
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>				<span class="comment">// advance</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>				w = x[k] &gt;&gt; (shift - nbits)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>				nbits = _W - (shift - nbits)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>			}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		<span class="comment">// convert digits of most-significant word w (omit leading zeros)</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		for w != 0 {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>			i--
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			s[i] = digits[w&amp;mask]
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			w &gt;&gt;= shift
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	} else {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		bb, ndigits := maxPow(b)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		<span class="comment">// construct table of successive squares of bb*leafSize to use in subdivisions</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		<span class="comment">// result (table != nil) &lt;=&gt; (len(x) &gt; leafSize &gt; 0)</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		table := divisors(len(x), b, ndigits, bb)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		<span class="comment">// preserve x, create local copy for use by convertWords</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		q := nat(nil).set(x)
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		<span class="comment">// convert q to string s in base b</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		q.convertWords(s, b, ndigits, bb, table)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		<span class="comment">// strip leading zeros</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		<span class="comment">// (x != 0; thus s must contain at least one non-zero digit</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		<span class="comment">// and the loop will terminate)</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		i = 0
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		for s[i] == &#39;0&#39; {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			i++
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	if neg {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		i--
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		s[i] = &#39;-&#39;
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	return s[i:]
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// Convert words of q to base b digits in s. If q is large, it is recursively &#34;split in half&#34;</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// by nat/nat division using tabulated divisors. Otherwise, it is converted iteratively using</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// repeated nat/Word division.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// The iterative method processes n Words by n divW() calls, each of which visits every Word in the</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">// incrementally shortened q for a total of n + (n-1) + (n-2) ... + 2 + 1, or n(n+1)/2 divW()&#39;s.</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// Recursive conversion divides q by its approximate square root, yielding two parts, each half</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span><span class="comment">// the size of q. Using the iterative method on both halves means 2 * (n/2)(n/2 + 1)/2 divW()&#39;s</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// plus the expensive long div(). Asymptotically, the ratio is favorable at 1/2 the divW()&#39;s, and</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// is made better by splitting the subblocks recursively. Best is to split blocks until one more</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span><span class="comment">// split would take longer (because of the nat/nat div()) than the twice as many divW()&#39;s of the</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">// iterative approach. This threshold is represented by leafSize. Benchmarking of leafSize in the</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// range 2..64 shows that values of 8 and 16 work well, with a 4x speedup at medium lengths and</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// ~30x for 20000 digits. Use nat_test.go&#39;s BenchmarkLeafSize tests to optimize leafSize for</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">// specific hardware.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>func (q nat) convertWords(s []byte, b Word, ndigits int, bb Word, table []divisor) {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// split larger blocks recursively</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	if table != nil {
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		<span class="comment">// len(q) &gt; leafSize &gt; 0</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		var r nat
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		index := len(table) - 1
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		for len(q) &gt; leafSize {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			<span class="comment">// find divisor close to sqrt(q) if possible, but in any case &lt; q</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			maxLength := q.bitLen()     <span class="comment">// ~= log2 q, or at of least largest possible q of this bit length</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			minLength := maxLength &gt;&gt; 1 <span class="comment">// ~= log2 sqrt(q)</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			for index &gt; 0 &amp;&amp; table[index-1].nbits &gt; minLength {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>				index-- <span class="comment">// desired</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			if table[index].nbits &gt;= maxLength &amp;&amp; table[index].bbb.cmp(q) &gt;= 0 {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>				index--
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>				if index &lt; 0 {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>					panic(&#34;internal inconsistency&#34;)
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>				}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			<span class="comment">// split q into the two digit number (q&#39;*bbb + r) to form independent subblocks</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			q, r = q.div(r, q, table[index].bbb)
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			<span class="comment">// convert subblocks and collect results in s[:h] and s[h:]</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			h := len(s) - table[index].ndigits
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			r.convertWords(s[h:], b, ndigits, bb, table[0:index])
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			s = s[:h] <span class="comment">// == q.convertWords(s, b, ndigits, bb, table[0:index+1])</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		}
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	<span class="comment">// having split any large blocks now process the remaining (small) block iteratively</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	i := len(s)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	var r Word
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	if b == 10 {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		<span class="comment">// hard-coding for 10 here speeds this up by 1.25x (allows for / and % by constants)</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		for len(q) &gt; 0 {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			<span class="comment">// extract least significant, base bb &#34;digit&#34;</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			q, r = q.divW(q, bb)
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			for j := 0; j &lt; ndigits &amp;&amp; i &gt; 0; j++ {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>				i--
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>				<span class="comment">// avoid % computation since r%10 == r - int(r/10)*10;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>				<span class="comment">// this appears to be faster for BenchmarkString10000Base10</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>				<span class="comment">// and smaller strings (but a bit slower for larger ones)</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>				t := r / 10
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>				s[i] = &#39;0&#39; + byte(r-t*10)
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>				r = t
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	} else {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		for len(q) &gt; 0 {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			<span class="comment">// extract least significant, base bb &#34;digit&#34;</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			q, r = q.divW(q, bb)
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			for j := 0; j &lt; ndigits &amp;&amp; i &gt; 0; j++ {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>				i--
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>				s[i] = digits[r%b]
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>				r /= b
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	<span class="comment">// prepend high-order zeros</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	for i &gt; 0 { <span class="comment">// while need more leading zeros</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		i--
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		s[i] = &#39;0&#39;
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span><span class="comment">// Split blocks greater than leafSize Words (or set to 0 to disable recursive conversion)</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span><span class="comment">// Benchmark and configure leafSize using: go test -bench=&#34;Leaf&#34;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span><span class="comment">//	8 and 16 effective on 3.0 GHz Xeon &#34;Clovertown&#34; CPU (128 byte cache lines)</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span><span class="comment">//	8 and 16 effective on 2.66 GHz Core 2 Duo &#34;Penryn&#34; CPU</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>var leafSize int = 8 <span class="comment">// number of Word-size binary values treat as a monolithic block</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>type divisor struct {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	bbb     nat <span class="comment">// divisor</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	nbits   int <span class="comment">// bit length of divisor (discounting leading zeros) ~= log2(bbb)</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	ndigits int <span class="comment">// digit length of divisor in terms of output base digits</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>var cacheBase10 struct {
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	sync.Mutex
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	table [64]divisor <span class="comment">// cached divisors for base 10</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span><span class="comment">// expWW computes x**y</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>func (z nat) expWW(x, y Word) nat {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	return z.expNN(nat(nil).setWord(x), nat(nil).setWord(y), nil, false)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span><span class="comment">// construct table of powers of bb*leafSize to use in subdivisions.</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>func divisors(m int, b Word, ndigits int, bb Word) []divisor {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	<span class="comment">// only compute table when recursive conversion is enabled and x is large</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	if leafSize == 0 || m &lt;= leafSize {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		return nil
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	<span class="comment">// determine k where (bb**leafSize)**(2**k) &gt;= sqrt(x)</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	k := 1
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	for words := leafSize; words &lt; m&gt;&gt;1 &amp;&amp; k &lt; len(cacheBase10.table); words &lt;&lt;= 1 {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		k++
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	<span class="comment">// reuse and extend existing table of divisors or create new table as appropriate</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	var table []divisor <span class="comment">// for b == 10, table overlaps with cacheBase10.table</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	if b == 10 {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		cacheBase10.Lock()
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		table = cacheBase10.table[0:k] <span class="comment">// reuse old table for this conversion</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	} else {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		table = make([]divisor, k) <span class="comment">// create new table for this conversion</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	<span class="comment">// extend table</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	if table[k-1].ndigits == 0 {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		<span class="comment">// add new entries as needed</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		var larger nat
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		for i := 0; i &lt; k; i++ {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			if table[i].ndigits == 0 {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>				if i == 0 {
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>					table[0].bbb = nat(nil).expWW(bb, Word(leafSize))
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>					table[0].ndigits = ndigits * leafSize
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>				} else {
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>					table[i].bbb = nat(nil).sqr(table[i-1].bbb)
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>					table[i].ndigits = 2 * table[i-1].ndigits
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>				}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>				<span class="comment">// optimization: exploit aggregated extra bits in macro blocks</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>				larger = nat(nil).set(table[i].bbb)
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>				for mulAddVWW(larger, larger, b, 0) == 0 {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>					table[i].bbb = table[i].bbb.set(larger)
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>					table[i].ndigits++
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>				}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>				table[i].nbits = table[i].bbb.bitLen()
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	if b == 10 {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		cacheBase10.Unlock()
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	return table
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>}
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>
</pre><p><a href="natconv.go?m=text">View as plain text</a></p>

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
