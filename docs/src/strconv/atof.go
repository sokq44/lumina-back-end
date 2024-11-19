<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/strconv/atof.go - Go Documentation Server</title>

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
<a href="atof.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/strconv">strconv</a>/<span class="text-muted">atof.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// decimal to binary floating point conversion.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// Algorithm:</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//   1) Store input in multiprecision decimal.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//   2) Multiply/divide decimal by powers of two until in range [0.5, 1)</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//   3) Multiply by 2^precision and round to get mantissa.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>import &#34;math&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>var optimize = true <span class="comment">// set to false to force slow-path conversions for testing</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// commonPrefixLenIgnoreCase returns the length of the common</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// prefix of s and prefix, with the character case of s ignored.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// The prefix argument must be all lower-case.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>func commonPrefixLenIgnoreCase(s, prefix string) int {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	n := len(prefix)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	if n &gt; len(s) {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>		n = len(s)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		c := s[i]
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		if &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;Z&#39; {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>			c += &#39;a&#39; - &#39;A&#39;
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		if c != prefix[i] {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>			return i
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	return n
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// special returns the floating-point value for the special,</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// possibly signed floating-point representations inf, infinity,</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// and NaN. The result is ok if a prefix of s contains one</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// of these representations and n is the length of that prefix.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// The character case is ignored.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func special(s string) (f float64, n int, ok bool) {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	if len(s) == 0 {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		return 0, 0, false
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	sign := 1
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	nsign := 0
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	switch s[0] {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	case &#39;+&#39;, &#39;-&#39;:
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		if s[0] == &#39;-&#39; {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			sign = -1
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		nsign = 1
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		s = s[1:]
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		fallthrough
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	case &#39;i&#39;, &#39;I&#39;:
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		n := commonPrefixLenIgnoreCase(s, &#34;infinity&#34;)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		<span class="comment">// Anything longer than &#34;inf&#34; is ok, but if we</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		<span class="comment">// don&#39;t have &#34;infinity&#34;, only consume &#34;inf&#34;.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		if 3 &lt; n &amp;&amp; n &lt; 8 {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			n = 3
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		if n == 3 || n == 8 {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			return math.Inf(sign), nsign + n, true
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	case &#39;n&#39;, &#39;N&#39;:
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		if commonPrefixLenIgnoreCase(s, &#34;nan&#34;) == 3 {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			return math.NaN(), 3, true
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	return 0, 0, false
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>func (b *decimal) set(s string) (ok bool) {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	i := 0
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	b.neg = false
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	b.trunc = false
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// optional sign</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if i &gt;= len(s) {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		return
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	switch {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	case s[i] == &#39;+&#39;:
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		i++
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	case s[i] == &#39;-&#39;:
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		b.neg = true
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		i++
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// digits</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	sawdot := false
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	sawdigits := false
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	for ; i &lt; len(s); i++ {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		switch {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		case s[i] == &#39;_&#39;:
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			<span class="comment">// readFloat already checked underscores</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			continue
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		case s[i] == &#39;.&#39;:
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			if sawdot {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>				return
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			sawdot = true
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>			b.dp = b.nd
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			continue
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		case &#39;0&#39; &lt;= s[i] &amp;&amp; s[i] &lt;= &#39;9&#39;:
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			sawdigits = true
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			if s[i] == &#39;0&#39; &amp;&amp; b.nd == 0 { <span class="comment">// ignore leading zeros</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>				b.dp--
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>				continue
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			if b.nd &lt; len(b.d) {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>				b.d[b.nd] = s[i]
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>				b.nd++
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			} else if s[i] != &#39;0&#39; {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>				b.trunc = true
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			continue
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		break
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if !sawdigits {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		return
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	if !sawdot {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		b.dp = b.nd
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// optional exponent moves decimal point.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// if we read a very large, very long number,</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// just be sure to move the decimal point by</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// a lot (say, 100000).  it doesn&#39;t matter if it&#39;s</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// not the exact number.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	if i &lt; len(s) &amp;&amp; lower(s[i]) == &#39;e&#39; {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		i++
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		if i &gt;= len(s) {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			return
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		esign := 1
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		if s[i] == &#39;+&#39; {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			i++
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		} else if s[i] == &#39;-&#39; {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			i++
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			esign = -1
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		if i &gt;= len(s) || s[i] &lt; &#39;0&#39; || s[i] &gt; &#39;9&#39; {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			return
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		e := 0
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		for ; i &lt; len(s) &amp;&amp; (&#39;0&#39; &lt;= s[i] &amp;&amp; s[i] &lt;= &#39;9&#39; || s[i] == &#39;_&#39;); i++ {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			if s[i] == &#39;_&#39; {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>				<span class="comment">// readFloat already checked underscores</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>				continue
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			if e &lt; 10000 {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>				e = e*10 + int(s[i]) - &#39;0&#39;
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		b.dp += e * esign
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if i != len(s) {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		return
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	ok = true
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	return
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// readFloat reads a decimal or hexadecimal mantissa and exponent from a float</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// string representation in s; the number may be followed by other characters.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// readFloat reports the number of bytes consumed (i), and whether the number</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">// is valid (ok).</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>func readFloat(s string) (mantissa uint64, exp int, neg, trunc, hex bool, i int, ok bool) {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	underscores := false
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// optional sign</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	if i &gt;= len(s) {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		return
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	switch {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	case s[i] == &#39;+&#39;:
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		i++
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	case s[i] == &#39;-&#39;:
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		neg = true
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		i++
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// digits</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	base := uint64(10)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	maxMantDigits := 19 <span class="comment">// 10^19 fits in uint64</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	expChar := byte(&#39;e&#39;)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	if i+2 &lt; len(s) &amp;&amp; s[i] == &#39;0&#39; &amp;&amp; lower(s[i+1]) == &#39;x&#39; {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		base = 16
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		maxMantDigits = 16 <span class="comment">// 16^16 fits in uint64</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		i += 2
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		expChar = &#39;p&#39;
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		hex = true
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	sawdot := false
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	sawdigits := false
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	nd := 0
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	ndMant := 0
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	dp := 0
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>loop:
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	for ; i &lt; len(s); i++ {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		switch c := s[i]; true {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		case c == &#39;_&#39;:
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			underscores = true
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			continue
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		case c == &#39;.&#39;:
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			if sawdot {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>				break loop
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			sawdot = true
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			dp = nd
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			continue
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		case &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39;:
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			sawdigits = true
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			if c == &#39;0&#39; &amp;&amp; nd == 0 { <span class="comment">// ignore leading zeros</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				dp--
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>				continue
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			nd++
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			if ndMant &lt; maxMantDigits {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>				mantissa *= base
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>				mantissa += uint64(c - &#39;0&#39;)
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>				ndMant++
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			} else if c != &#39;0&#39; {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>				trunc = true
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			continue
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		case base == 16 &amp;&amp; &#39;a&#39; &lt;= lower(c) &amp;&amp; lower(c) &lt;= &#39;f&#39;:
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			sawdigits = true
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			nd++
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			if ndMant &lt; maxMantDigits {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>				mantissa *= 16
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>				mantissa += uint64(lower(c) - &#39;a&#39; + 10)
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>				ndMant++
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			} else {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>				trunc = true
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			continue
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		break
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	if !sawdigits {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		return
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	if !sawdot {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		dp = nd
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	if base == 16 {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		dp *= 4
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		ndMant *= 4
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// optional exponent moves decimal point.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// if we read a very large, very long number,</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	<span class="comment">// just be sure to move the decimal point by</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	<span class="comment">// a lot (say, 100000).  it doesn&#39;t matter if it&#39;s</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	<span class="comment">// not the exact number.</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	if i &lt; len(s) &amp;&amp; lower(s[i]) == expChar {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		i++
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		if i &gt;= len(s) {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			return
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		esign := 1
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		if s[i] == &#39;+&#39; {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			i++
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		} else if s[i] == &#39;-&#39; {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			i++
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			esign = -1
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		if i &gt;= len(s) || s[i] &lt; &#39;0&#39; || s[i] &gt; &#39;9&#39; {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			return
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		e := 0
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		for ; i &lt; len(s) &amp;&amp; (&#39;0&#39; &lt;= s[i] &amp;&amp; s[i] &lt;= &#39;9&#39; || s[i] == &#39;_&#39;); i++ {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			if s[i] == &#39;_&#39; {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>				underscores = true
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>				continue
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			if e &lt; 10000 {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>				e = e*10 + int(s[i]) - &#39;0&#39;
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		dp += e * esign
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	} else if base == 16 {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		<span class="comment">// Must have exponent.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		return
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	if mantissa != 0 {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		exp = dp - ndMant
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	if underscores &amp;&amp; !underscoreOK(s[:i]) {
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		return
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	ok = true
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	return
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// decimal power of ten to binary power of two.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>var powtab = []int{1, 3, 6, 9, 13, 16, 19, 23, 26}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>func (d *decimal) floatBits(flt *floatInfo) (b uint64, overflow bool) {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	var exp int
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	var mant uint64
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	<span class="comment">// Zero is always a special case.</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	if d.nd == 0 {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		mant = 0
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		exp = flt.bias
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		goto out
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// Obvious overflow/underflow.</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">// These bounds are for 64-bit floats.</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// Will have to change if we want to support 80-bit floats in the future.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if d.dp &gt; 310 {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		goto overflow
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	if d.dp &lt; -330 {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		<span class="comment">// zero</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		mant = 0
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		exp = flt.bias
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		goto out
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// Scale by powers of two until in range [0.5, 1.0)</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	exp = 0
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	for d.dp &gt; 0 {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		var n int
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		if d.dp &gt;= len(powtab) {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			n = 27
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		} else {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			n = powtab[d.dp]
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		d.Shift(-n)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		exp += n
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	for d.dp &lt; 0 || d.dp == 0 &amp;&amp; d.d[0] &lt; &#39;5&#39; {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		var n int
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		if -d.dp &gt;= len(powtab) {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>			n = 27
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		} else {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			n = powtab[-d.dp]
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		d.Shift(n)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		exp -= n
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	<span class="comment">// Our range is [0.5,1) but floating point range is [1,2).</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	exp--
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	<span class="comment">// Minimum representable exponent is flt.bias+1.</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	<span class="comment">// If the exponent is smaller, move it up and</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	<span class="comment">// adjust d accordingly.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	if exp &lt; flt.bias+1 {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		n := flt.bias + 1 - exp
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		d.Shift(-n)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		exp += n
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	if exp-flt.bias &gt;= 1&lt;&lt;flt.expbits-1 {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		goto overflow
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	}
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	<span class="comment">// Extract 1+flt.mantbits bits.</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	d.Shift(int(1 + flt.mantbits))
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	mant = d.RoundedInteger()
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	<span class="comment">// Rounding might have added a bit; shift down.</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	if mant == 2&lt;&lt;flt.mantbits {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		mant &gt;&gt;= 1
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		exp++
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		if exp-flt.bias &gt;= 1&lt;&lt;flt.expbits-1 {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			goto overflow
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">// Denormalized?</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	if mant&amp;(1&lt;&lt;flt.mantbits) == 0 {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		exp = flt.bias
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	goto out
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>overflow:
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	<span class="comment">// ±Inf</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	mant = 0
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	exp = 1&lt;&lt;flt.expbits - 1 + flt.bias
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	overflow = true
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>out:
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	<span class="comment">// Assemble bits.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	bits := mant &amp; (uint64(1)&lt;&lt;flt.mantbits - 1)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	bits |= uint64((exp-flt.bias)&amp;(1&lt;&lt;flt.expbits-1)) &lt;&lt; flt.mantbits
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	if d.neg {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		bits |= 1 &lt;&lt; flt.mantbits &lt;&lt; flt.expbits
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	return bits, overflow
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span><span class="comment">// Exact powers of 10.</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>var float64pow10 = []float64{
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	1e20, 1e21, 1e22,
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>var float32pow10 = []float32{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span><span class="comment">// If possible to convert decimal representation to 64-bit float f exactly,</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span><span class="comment">// entirely in floating-point math, do so, avoiding the expense of decimalToFloatBits.</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span><span class="comment">// Three common cases:</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span><span class="comment">//	value is exact integer</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span><span class="comment">//	value is exact integer * exact power of ten</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span><span class="comment">//	value is exact integer / exact power of ten</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span><span class="comment">// These all produce potentially inexact but correctly rounded answers.</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>func atof64exact(mantissa uint64, exp int, neg bool) (f float64, ok bool) {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	if mantissa&gt;&gt;float64info.mantbits != 0 {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		return
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	f = float64(mantissa)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	if neg {
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		f = -f
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	}
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	switch {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	case exp == 0:
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		<span class="comment">// an integer.</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		return f, true
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// Exact integers are &lt;= 10^15.</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	<span class="comment">// Exact powers of ten are &lt;= 10^22.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	case exp &gt; 0 &amp;&amp; exp &lt;= 15+22: <span class="comment">// int * 10^k</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		<span class="comment">// If exponent is big but number of digits is not,</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		<span class="comment">// can move a few zeros into the integer part.</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		if exp &gt; 22 {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>			f *= float64pow10[exp-22]
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>			exp = 22
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		if f &gt; 1e15 || f &lt; -1e15 {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			<span class="comment">// the exponent was really too large.</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			return
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		return f * float64pow10[exp], true
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	case exp &lt; 0 &amp;&amp; exp &gt;= -22: <span class="comment">// int / 10^k</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		return f / float64pow10[-exp], true
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	return
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span><span class="comment">// If possible to compute mantissa*10^exp to 32-bit float f exactly,</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span><span class="comment">// entirely in floating-point math, do so, avoiding the machinery above.</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>func atof32exact(mantissa uint64, exp int, neg bool) (f float32, ok bool) {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	if mantissa&gt;&gt;float32info.mantbits != 0 {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		return
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	f = float32(mantissa)
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	if neg {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		f = -f
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	switch {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	case exp == 0:
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		return f, true
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	<span class="comment">// Exact integers are &lt;= 10^7.</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	<span class="comment">// Exact powers of ten are &lt;= 10^10.</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	case exp &gt; 0 &amp;&amp; exp &lt;= 7+10: <span class="comment">// int * 10^k</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		<span class="comment">// If exponent is big but number of digits is not,</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		<span class="comment">// can move a few zeros into the integer part.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		if exp &gt; 10 {
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			f *= float32pow10[exp-10]
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			exp = 10
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		if f &gt; 1e7 || f &lt; -1e7 {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			<span class="comment">// the exponent was really too large.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			return
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		}
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		return f * float32pow10[exp], true
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	case exp &lt; 0 &amp;&amp; exp &gt;= -10: <span class="comment">// int / 10^k</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		return f / float32pow10[-exp], true
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	return
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span><span class="comment">// atofHex converts the hex floating-point string s</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span><span class="comment">// to a rounded float32 or float64 value (depending on flt==&amp;float32info or flt==&amp;float64info)</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span><span class="comment">// and returns it as a float64.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span><span class="comment">// The string s has already been parsed into a mantissa, exponent, and sign (neg==true for negative).</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span><span class="comment">// If trunc is true, trailing non-zero bits have been omitted from the mantissa.</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>func atofHex(s string, flt *floatInfo, mantissa uint64, exp int, neg, trunc bool) (float64, error) {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	maxExp := 1&lt;&lt;flt.expbits + flt.bias - 2
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	minExp := flt.bias + 1
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	exp += int(flt.mantbits) <span class="comment">// mantissa now implicitly divided by 2^mantbits.</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	<span class="comment">// Shift mantissa and exponent to bring representation into float range.</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	<span class="comment">// Eventually we want a mantissa with a leading 1-bit followed by mantbits other bits.</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	<span class="comment">// For rounding, we need two more, where the bottom bit represents</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	<span class="comment">// whether that bit or any later bit was non-zero.</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	<span class="comment">// (If the mantissa has already lost non-zero bits, trunc is true,</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	<span class="comment">// and we OR in a 1 below after shifting left appropriately.)</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	for mantissa != 0 &amp;&amp; mantissa&gt;&gt;(flt.mantbits+2) == 0 {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		mantissa &lt;&lt;= 1
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		exp--
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	if trunc {
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		mantissa |= 1
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	for mantissa&gt;&gt;(1+flt.mantbits+2) != 0 {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		mantissa = mantissa&gt;&gt;1 | mantissa&amp;1
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		exp++
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	<span class="comment">// If exponent is too negative,</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	<span class="comment">// denormalize in hopes of making it representable.</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	<span class="comment">// (The -2 is for the rounding bits.)</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	for mantissa &gt; 1 &amp;&amp; exp &lt; minExp-2 {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		mantissa = mantissa&gt;&gt;1 | mantissa&amp;1
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		exp++
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	}
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	<span class="comment">// Round using two bottom bits.</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	round := mantissa &amp; 3
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	mantissa &gt;&gt;= 2
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	round |= mantissa &amp; 1 <span class="comment">// round to even (round up if mantissa is odd)</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	exp += 2
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	if round == 3 {
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		mantissa++
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		if mantissa == 1&lt;&lt;(1+flt.mantbits) {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>			mantissa &gt;&gt;= 1
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>			exp++
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		}
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	}
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	if mantissa&gt;&gt;flt.mantbits == 0 { <span class="comment">// Denormal or zero.</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		exp = flt.bias
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	var err error
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	if exp &gt; maxExp { <span class="comment">// infinity and range error</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		mantissa = 1 &lt;&lt; flt.mantbits
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		exp = maxExp + 1
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		err = rangeError(fnParseFloat, s)
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	bits := mantissa &amp; (1&lt;&lt;flt.mantbits - 1)
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	bits |= uint64((exp-flt.bias)&amp;(1&lt;&lt;flt.expbits-1)) &lt;&lt; flt.mantbits
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	if neg {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		bits |= 1 &lt;&lt; flt.mantbits &lt;&lt; flt.expbits
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	if flt == &amp;float32info {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		return float64(math.Float32frombits(uint32(bits))), err
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	return math.Float64frombits(bits), err
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>const fnParseFloat = &#34;ParseFloat&#34;
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>func atof32(s string) (f float32, n int, err error) {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	if val, n, ok := special(s); ok {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		return float32(val), n, nil
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	mantissa, exp, neg, trunc, hex, n, ok := readFloat(s)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	if !ok {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		return 0, n, syntaxError(fnParseFloat, s)
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	if hex {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		f, err := atofHex(s[:n], &amp;float32info, mantissa, exp, neg, trunc)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		return float32(f), n, err
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	if optimize {
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		<span class="comment">// Try pure floating-point arithmetic conversion, and if that fails,</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		<span class="comment">// the Eisel-Lemire algorithm.</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		if !trunc {
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>			if f, ok := atof32exact(mantissa, exp, neg); ok {
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>				return f, n, nil
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>			}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		f, ok := eiselLemire32(mantissa, exp, neg)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		if ok {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>			if !trunc {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>				return f, n, nil
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>			}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			<span class="comment">// Even if the mantissa was truncated, we may</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>			<span class="comment">// have found the correct result. Confirm by</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>			<span class="comment">// converting the upper mantissa bound.</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>			fUp, ok := eiselLemire32(mantissa+1, exp, neg)
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>			if ok &amp;&amp; f == fUp {
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>				return f, n, nil
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>			}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	<span class="comment">// Slow fallback.</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	var d decimal
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	if !d.set(s[:n]) {
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		return 0, n, syntaxError(fnParseFloat, s)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	b, ovf := d.floatBits(&amp;float32info)
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	f = math.Float32frombits(uint32(b))
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	if ovf {
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		err = rangeError(fnParseFloat, s)
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	}
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	return f, n, err
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>}
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>func atof64(s string) (f float64, n int, err error) {
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	if val, n, ok := special(s); ok {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		return val, n, nil
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	}
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	mantissa, exp, neg, trunc, hex, n, ok := readFloat(s)
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	if !ok {
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		return 0, n, syntaxError(fnParseFloat, s)
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	}
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	if hex {
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		f, err := atofHex(s[:n], &amp;float64info, mantissa, exp, neg, trunc)
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		return f, n, err
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	if optimize {
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>		<span class="comment">// Try pure floating-point arithmetic conversion, and if that fails,</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		<span class="comment">// the Eisel-Lemire algorithm.</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		if !trunc {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>			if f, ok := atof64exact(mantissa, exp, neg); ok {
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>				return f, n, nil
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>			}
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		}
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		f, ok := eiselLemire64(mantissa, exp, neg)
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		if ok {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>			if !trunc {
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>				return f, n, nil
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>			}
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>			<span class="comment">// Even if the mantissa was truncated, we may</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>			<span class="comment">// have found the correct result. Confirm by</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>			<span class="comment">// converting the upper mantissa bound.</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>			fUp, ok := eiselLemire64(mantissa+1, exp, neg)
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>			if ok &amp;&amp; f == fUp {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>				return f, n, nil
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>			}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		}
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	}
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	<span class="comment">// Slow fallback.</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	var d decimal
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	if !d.set(s[:n]) {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		return 0, n, syntaxError(fnParseFloat, s)
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	}
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	b, ovf := d.floatBits(&amp;float64info)
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	f = math.Float64frombits(b)
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	if ovf {
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		err = rangeError(fnParseFloat, s)
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	}
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	return f, n, err
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span><span class="comment">// ParseFloat converts the string s to a floating-point number</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span><span class="comment">// with the precision specified by bitSize: 32 for float32, or 64 for float64.</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span><span class="comment">// When bitSize=32, the result still has type float64, but it will be</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span><span class="comment">// convertible to float32 without changing its value.</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span><span class="comment">// ParseFloat accepts decimal and hexadecimal floating-point numbers</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span><span class="comment">// as defined by the Go syntax for [floating-point literals].</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span><span class="comment">// If s is well-formed and near a valid floating-point number,</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span><span class="comment">// ParseFloat returns the nearest floating-point number rounded</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span><span class="comment">// using IEEE754 unbiased rounding.</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span><span class="comment">// (Parsing a hexadecimal floating-point value only rounds when</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span><span class="comment">// there are more bits in the hexadecimal representation than</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span><span class="comment">// will fit in the mantissa.)</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span><span class="comment">// The errors that ParseFloat returns have concrete type *NumError</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span><span class="comment">// and include err.Num = s.</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span><span class="comment">// If s is not syntactically well-formed, ParseFloat returns err.Err = ErrSyntax.</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span><span class="comment">// If s is syntactically well-formed but is more than 1/2 ULP</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span><span class="comment">// away from the largest floating point number of the given size,</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span><span class="comment">// ParseFloat returns f = ±Inf, err.Err = ErrRange.</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span><span class="comment">// ParseFloat recognizes the string &#34;NaN&#34;, and the (possibly signed) strings &#34;Inf&#34; and &#34;Infinity&#34;</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span><span class="comment">// as their respective special floating point values. It ignores case when matching.</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span><span class="comment">// [floating-point literals]: https://go.dev/ref/spec#Floating-point_literals</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>func ParseFloat(s string, bitSize int) (float64, error) {
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	f, n, err := parseFloatPrefix(s, bitSize)
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	if n != len(s) &amp;&amp; (err == nil || err.(*NumError).Err != ErrSyntax) {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		return 0, syntaxError(fnParseFloat, s)
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	}
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	return f, err
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>}
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>func parseFloatPrefix(s string, bitSize int) (float64, int, error) {
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	if bitSize == 32 {
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		f, n, err := atof32(s)
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		return float64(f), n, err
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	}
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	return atof64(s)
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>}
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>
</pre><p><a href="atof.go?m=text">View as plain text</a></p>

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
