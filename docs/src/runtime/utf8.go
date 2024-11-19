<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/utf8.go - Go Documentation Server</title>

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
<a href="utf8.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">utf8.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2016 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Numbers fundamental to the encoding.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>const (
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	runeError = &#39;\uFFFD&#39;     <span class="comment">// the &#34;error&#34; Rune or &#34;Unicode replacement character&#34;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	runeSelf  = 0x80         <span class="comment">// characters below runeSelf are represented as themselves in a single byte.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	maxRune   = &#39;\U0010FFFF&#39; <span class="comment">// Maximum valid Unicode code point.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// Code points in the surrogate range are not valid for UTF-8.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>const (
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	surrogateMin = 0xD800
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	surrogateMax = 0xDFFF
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>const (
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	t1 = 0x00 <span class="comment">// 0000 0000</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	tx = 0x80 <span class="comment">// 1000 0000</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	t2 = 0xC0 <span class="comment">// 1100 0000</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	t3 = 0xE0 <span class="comment">// 1110 0000</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	t4 = 0xF0 <span class="comment">// 1111 0000</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	t5 = 0xF8 <span class="comment">// 1111 1000</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	maskx = 0x3F <span class="comment">// 0011 1111</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	mask2 = 0x1F <span class="comment">// 0001 1111</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	mask3 = 0x0F <span class="comment">// 0000 1111</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	mask4 = 0x07 <span class="comment">// 0000 0111</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	rune1Max = 1&lt;&lt;7 - 1
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	rune2Max = 1&lt;&lt;11 - 1
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	rune3Max = 1&lt;&lt;16 - 1
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// The default lowest and highest continuation byte.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	locb = 0x80 <span class="comment">// 1000 0000</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	hicb = 0xBF <span class="comment">// 1011 1111</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// countrunes returns the number of runes in s.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func countrunes(s string) int {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	n := 0
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	for range s {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		n++
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	return n
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// decoderune returns the non-ASCII rune at the start of</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// s[k:] and the index after the rune in s.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// decoderune assumes that caller has checked that</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// the to be decoded rune is a non-ASCII rune.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// If the string appears to be incomplete or decoding problems</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// are encountered (runeerror, k + 1) is returned to ensure</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// progress when decoderune is used to iterate over a string.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func decoderune(s string, k int) (r rune, pos int) {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	pos = k
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if k &gt;= len(s) {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		return runeError, k + 1
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	s = s[k:]
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	switch {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	case t2 &lt;= s[0] &amp;&amp; s[0] &lt; t3:
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		<span class="comment">// 0080-07FF two byte sequence</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		if len(s) &gt; 1 &amp;&amp; (locb &lt;= s[1] &amp;&amp; s[1] &lt;= hicb) {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			r = rune(s[0]&amp;mask2)&lt;&lt;6 | rune(s[1]&amp;maskx)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			pos += 2
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			if rune1Max &lt; r {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>				return
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	case t3 &lt;= s[0] &amp;&amp; s[0] &lt; t4:
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		<span class="comment">// 0800-FFFF three byte sequence</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		if len(s) &gt; 2 &amp;&amp; (locb &lt;= s[1] &amp;&amp; s[1] &lt;= hicb) &amp;&amp; (locb &lt;= s[2] &amp;&amp; s[2] &lt;= hicb) {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			r = rune(s[0]&amp;mask3)&lt;&lt;12 | rune(s[1]&amp;maskx)&lt;&lt;6 | rune(s[2]&amp;maskx)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			pos += 3
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			if rune2Max &lt; r &amp;&amp; !(surrogateMin &lt;= r &amp;&amp; r &lt;= surrogateMax) {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>				return
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	case t4 &lt;= s[0] &amp;&amp; s[0] &lt; t5:
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		<span class="comment">// 10000-1FFFFF four byte sequence</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		if len(s) &gt; 3 &amp;&amp; (locb &lt;= s[1] &amp;&amp; s[1] &lt;= hicb) &amp;&amp; (locb &lt;= s[2] &amp;&amp; s[2] &lt;= hicb) &amp;&amp; (locb &lt;= s[3] &amp;&amp; s[3] &lt;= hicb) {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			r = rune(s[0]&amp;mask4)&lt;&lt;18 | rune(s[1]&amp;maskx)&lt;&lt;12 | rune(s[2]&amp;maskx)&lt;&lt;6 | rune(s[3]&amp;maskx)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			pos += 4
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			if rune3Max &lt; r &amp;&amp; r &lt;= maxRune {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>				return
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	return runeError, k + 1
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// encoderune writes into p (which must be large enough) the UTF-8 encoding of the rune.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// It returns the number of bytes written.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>func encoderune(p []byte, r rune) int {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// Negative values are erroneous. Making it unsigned addresses the problem.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	switch i := uint32(r); {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	case i &lt;= rune1Max:
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		p[0] = byte(r)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		return 1
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	case i &lt;= rune2Max:
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		_ = p[1] <span class="comment">// eliminate bounds checks</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		p[0] = t2 | byte(r&gt;&gt;6)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		p[1] = tx | byte(r)&amp;maskx
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		return 2
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	case i &gt; maxRune, surrogateMin &lt;= i &amp;&amp; i &lt;= surrogateMax:
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		r = runeError
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		fallthrough
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	case i &lt;= rune3Max:
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		_ = p[2] <span class="comment">// eliminate bounds checks</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		p[0] = t3 | byte(r&gt;&gt;12)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		p[1] = tx | byte(r&gt;&gt;6)&amp;maskx
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		p[2] = tx | byte(r)&amp;maskx
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		return 3
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	default:
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		_ = p[3] <span class="comment">// eliminate bounds checks</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		p[0] = t4 | byte(r&gt;&gt;18)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		p[1] = tx | byte(r&gt;&gt;12)&amp;maskx
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		p[2] = tx | byte(r&gt;&gt;6)&amp;maskx
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		p[3] = tx | byte(r)&amp;maskx
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		return 4
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
</pre><p><a href="utf8.go?m=text">View as plain text</a></p>

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
