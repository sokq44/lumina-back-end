<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/strconv/itoa.go - Go Documentation Server</title>

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
<a href="itoa.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/strconv">strconv</a>/<span class="text-muted">itoa.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package strconv
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;math/bits&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>const fastSmalls = true <span class="comment">// enable fast path for small integers</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// FormatUint returns the string representation of i in the given base,</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// for 2 &lt;= base &lt;= 36. The result uses the lower-case letters &#39;a&#39; to &#39;z&#39;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// for digit values &gt;= 10.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>func FormatUint(i uint64, base int) string {
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	if fastSmalls &amp;&amp; i &lt; nSmalls &amp;&amp; base == 10 {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>		return small(int(i))
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	}
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	_, s := formatBits(nil, i, base, false, false)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	return s
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>}
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// FormatInt returns the string representation of i in the given base,</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// for 2 &lt;= base &lt;= 36. The result uses the lower-case letters &#39;a&#39; to &#39;z&#39;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// for digit values &gt;= 10.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>func FormatInt(i int64, base int) string {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	if fastSmalls &amp;&amp; 0 &lt;= i &amp;&amp; i &lt; nSmalls &amp;&amp; base == 10 {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		return small(int(i))
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	_, s := formatBits(nil, uint64(i), base, i &lt; 0, false)
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	return s
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// Itoa is equivalent to FormatInt(int64(i), 10).</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>func Itoa(i int) string {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	return FormatInt(int64(i), 10)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// AppendInt appends the string form of the integer i,</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// as generated by FormatInt, to dst and returns the extended buffer.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>func AppendInt(dst []byte, i int64, base int) []byte {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	if fastSmalls &amp;&amp; 0 &lt;= i &amp;&amp; i &lt; nSmalls &amp;&amp; base == 10 {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		return append(dst, small(int(i))...)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	dst, _ = formatBits(dst, uint64(i), base, i &lt; 0, true)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	return dst
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// AppendUint appends the string form of the unsigned integer i,</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// as generated by FormatUint, to dst and returns the extended buffer.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>func AppendUint(dst []byte, i uint64, base int) []byte {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	if fastSmalls &amp;&amp; i &lt; nSmalls &amp;&amp; base == 10 {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		return append(dst, small(int(i))...)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	dst, _ = formatBits(dst, i, base, false, true)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	return dst
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// small returns the string for an i with 0 &lt;= i &lt; nSmalls.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>func small(i int) string {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if i &lt; 10 {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return digits[i : i+1]
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	return smallsString[i*2 : i*2+2]
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>const nSmalls = 100
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>const smallsString = &#34;00010203040506070809&#34; +
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	&#34;10111213141516171819&#34; +
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	&#34;20212223242526272829&#34; +
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	&#34;30313233343536373839&#34; +
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	&#34;40414243444546474849&#34; +
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	&#34;50515253545556575859&#34; +
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	&#34;60616263646566676869&#34; +
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	&#34;70717273747576777879&#34; +
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	&#34;80818283848586878889&#34; +
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	&#34;90919293949596979899&#34;
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>const host32bit = ^uint(0)&gt;&gt;32 == 0
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>const digits = &#34;0123456789abcdefghijklmnopqrstuvwxyz&#34;
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// formatBits computes the string representation of u in the given base.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// If neg is set, u is treated as negative int64 value. If append_ is</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// set, the string is appended to dst and the resulting byte slice is</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// returned as the first result value; otherwise the string is returned</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// as the second result value.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>func formatBits(dst []byte, u uint64, base int, neg, append_ bool) (d []byte, s string) {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	if base &lt; 2 || base &gt; len(digits) {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		panic(&#34;strconv: illegal AppendInt/FormatInt base&#34;)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// 2 &lt;= base &amp;&amp; base &lt;= len(digits)</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	var a [64 + 1]byte <span class="comment">// +1 for sign of 64bit value in base 2</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	i := len(a)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	if neg {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		u = -u
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// convert bits</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// We use uint values where we can because those will</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// fit into a single register even on a 32bit machine.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	if base == 10 {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		<span class="comment">// common case: use constants for / because</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">// the compiler can optimize it into a multiply+shift</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		if host32bit {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			<span class="comment">// convert the lower digits using 32bit operations</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			for u &gt;= 1e9 {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>				<span class="comment">// Avoid using r = a%b in addition to q = a/b</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>				<span class="comment">// since 64bit division and modulo operations</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>				<span class="comment">// are calculated by runtime functions on 32bit machines.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>				q := u / 1e9
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>				us := uint(u - q*1e9) <span class="comment">// u % 1e9 fits into a uint</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>				for j := 4; j &gt; 0; j-- {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>					is := us % 100 * 2
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>					us /= 100
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>					i -= 2
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>					a[i+1] = smallsString[is+1]
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>					a[i+0] = smallsString[is+0]
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>				}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>				<span class="comment">// us &lt; 10, since it contains the last digit</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>				<span class="comment">// from the initial 9-digit us.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>				i--
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>				a[i] = smallsString[us*2+1]
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>				u = q
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			<span class="comment">// u &lt; 1e9</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// u guaranteed to fit into a uint</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		us := uint(u)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		for us &gt;= 100 {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			is := us % 100 * 2
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			us /= 100
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			i -= 2
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			a[i+1] = smallsString[is+1]
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			a[i+0] = smallsString[is+0]
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">// us &lt; 100</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		is := us * 2
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		i--
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		a[i] = smallsString[is+1]
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		if us &gt;= 10 {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			i--
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			a[i] = smallsString[is]
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	} else if isPowerOfTwo(base) {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		<span class="comment">// Use shifts and masks instead of / and %.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		<span class="comment">// Base is a power of 2 and 2 &lt;= base &lt;= len(digits) where len(digits) is 36.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		<span class="comment">// The largest power of 2 below or equal to 36 is 32, which is 1 &lt;&lt; 5;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// i.e., the largest possible shift count is 5. By &amp;-ind that value with</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// the constant 7 we tell the compiler that the shift count is always</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		<span class="comment">// less than 8 which is smaller than any register width. This allows</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		<span class="comment">// the compiler to generate better code for the shift operation.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		shift := uint(bits.TrailingZeros(uint(base))) &amp; 7
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		b := uint64(base)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		m := uint(base) - 1 <span class="comment">// == 1&lt;&lt;shift - 1</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		for u &gt;= b {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			i--
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			a[i] = digits[uint(u)&amp;m]
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>			u &gt;&gt;= shift
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		<span class="comment">// u &lt; base</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		i--
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		a[i] = digits[uint(u)]
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	} else {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		<span class="comment">// general case</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		b := uint64(base)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		for u &gt;= b {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			i--
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			<span class="comment">// Avoid using r = a%b in addition to q = a/b</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			<span class="comment">// since 64bit division and modulo operations</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			<span class="comment">// are calculated by runtime functions on 32bit machines.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			q := u / b
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			a[i] = digits[uint(u-q*b)]
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			u = q
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		<span class="comment">// u &lt; base</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		i--
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		a[i] = digits[uint(u)]
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// add sign, if any</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	if neg {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		i--
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		a[i] = &#39;-&#39;
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	if append_ {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		d = append(dst, a[i:]...)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		return
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	s = string(a[i:])
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	return
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>func isPowerOfTwo(x int) bool {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	return x&amp;(x-1) == 0
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
</pre><p><a href="itoa.go?m=text">View as plain text</a></p>

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
