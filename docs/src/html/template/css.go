<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/html/template/css.go - Go Documentation Server</title>

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
<a href="css.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/html">html</a>/<a href="http://localhost:8080/src/html/template">template</a>/<span class="text-muted">css.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/html/template">html/template</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package template
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// endsWithCSSKeyword reports whether b ends with an ident that</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// case-insensitively matches the lower-case kw.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>func endsWithCSSKeyword(b []byte, kw string) bool {
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	i := len(b) - len(kw)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>		<span class="comment">// Too short.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>		return false
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	if i != 0 {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>		r, _ := utf8.DecodeLastRune(b[:i])
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		if isCSSNmchar(r) {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>			<span class="comment">// Too long.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>			return false
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// Many CSS keywords, such as &#34;!important&#34; can have characters encoded,</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// but the URI production does not allow that according to</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// https://www.w3.org/TR/css3-syntax/#TOK-URI</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// This does not attempt to recognize encoded keywords. For example,</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// given &#34;\75\72\6c&#34; and &#34;url&#34; this return false.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	return string(bytes.ToLower(b[i:])) == kw
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// isCSSNmchar reports whether rune is allowed anywhere in a CSS identifier.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>func isCSSNmchar(r rune) bool {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// Based on the CSS3 nmchar production but ignores multi-rune escape</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// sequences.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// https://www.w3.org/TR/css3-syntax/#SUBTOK-nmchar</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	return &#39;a&#39; &lt;= r &amp;&amp; r &lt;= &#39;z&#39; ||
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		&#39;A&#39; &lt;= r &amp;&amp; r &lt;= &#39;Z&#39; ||
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		&#39;0&#39; &lt;= r &amp;&amp; r &lt;= &#39;9&#39; ||
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		r == &#39;-&#39; ||
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		r == &#39;_&#39; ||
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		<span class="comment">// Non-ASCII cases below.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		0x80 &lt;= r &amp;&amp; r &lt;= 0xd7ff ||
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		0xe000 &lt;= r &amp;&amp; r &lt;= 0xfffd ||
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		0x10000 &lt;= r &amp;&amp; r &lt;= 0x10ffff
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// decodeCSS decodes CSS3 escapes given a sequence of stringchars.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// If there is no change, it returns the input, otherwise it returns a slice</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// backed by a new array.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// https://www.w3.org/TR/css3-syntax/#SUBTOK-stringchar defines stringchar.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func decodeCSS(s []byte) []byte {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	i := bytes.IndexByte(s, &#39;\\&#39;)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if i == -1 {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return s
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// The UTF-8 sequence for a codepoint is never longer than 1 + the</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// number hex digits need to represent that codepoint, so len(s) is an</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// upper bound on the output length.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	b := make([]byte, 0, len(s))
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	for len(s) != 0 {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		i := bytes.IndexByte(s, &#39;\\&#39;)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		if i == -1 {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			i = len(s)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		b, s = append(b, s[:i]...), s[i:]
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		if len(s) &lt; 2 {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			break
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		<span class="comment">// https://www.w3.org/TR/css3-syntax/#SUBTOK-escape</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		<span class="comment">// escape ::= unicode | &#39;\&#39; [#x20-#x7E#x80-#xD7FF#xE000-#xFFFD#x10000-#x10FFFF]</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		if isHex(s[1]) {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			<span class="comment">// https://www.w3.org/TR/css3-syntax/#SUBTOK-unicode</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			<span class="comment">//   unicode ::= &#39;\&#39; [0-9a-fA-F]{1,6} wc?</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			j := 2
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			for j &lt; len(s) &amp;&amp; j &lt; 7 &amp;&amp; isHex(s[j]) {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>				j++
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			r := hexDecode(s[1:j])
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			if r &gt; unicode.MaxRune {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>				r, j = r/16, j-1
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			n := utf8.EncodeRune(b[len(b):cap(b)], r)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			<span class="comment">// The optional space at the end allows a hex</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			<span class="comment">// sequence to be followed by a literal hex.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			<span class="comment">// string(decodeCSS([]byte(`\A B`))) == &#34;\nB&#34;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			b, s = b[:len(b)+n], skipCSSSpace(s[j:])
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		} else {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			<span class="comment">// `\\` decodes to `\` and `\&#34;` to `&#34;`.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			_, n := utf8.DecodeRune(s[1:])
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			b, s = append(b, s[1:1+n]...), s[1+n:]
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	return b
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// isHex reports whether the given character is a hex digit.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>func isHex(c byte) bool {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	return &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39; || &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;f&#39; || &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;F&#39;
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// hexDecode decodes a short hex digit sequence: &#34;10&#34; -&gt; 16.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>func hexDecode(s []byte) rune {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	n := &#39;\x00&#39;
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	for _, c := range s {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		n &lt;&lt;= 4
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		switch {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		case &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39;:
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			n |= rune(c - &#39;0&#39;)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		case &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;f&#39;:
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			n |= rune(c-&#39;a&#39;) + 10
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		case &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;F&#39;:
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			n |= rune(c-&#39;A&#39;) + 10
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		default:
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			panic(fmt.Sprintf(&#34;Bad hex digit in %q&#34;, s))
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	return n
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// skipCSSSpace returns a suffix of c, skipping over a single space.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>func skipCSSSpace(c []byte) []byte {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if len(c) == 0 {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		return c
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// wc ::= #x9 | #xA | #xC | #xD | #x20</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	switch c[0] {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	case &#39;\t&#39;, &#39;\n&#39;, &#39;\f&#39;, &#39; &#39;:
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		return c[1:]
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	case &#39;\r&#39;:
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		<span class="comment">// This differs from CSS3&#39;s wc production because it contains a</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		<span class="comment">// probable spec error whereby wc contains all the single byte</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		<span class="comment">// sequences in nl (newline) but not CRLF.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		if len(c) &gt;= 2 &amp;&amp; c[1] == &#39;\n&#39; {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			return c[2:]
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		return c[1:]
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	return c
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// isCSSSpace reports whether b is a CSS space char as defined in wc.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>func isCSSSpace(b byte) bool {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	switch b {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	case &#39;\t&#39;, &#39;\n&#39;, &#39;\f&#39;, &#39;\r&#39;, &#39; &#39;:
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		return true
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	return false
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// cssEscaper escapes HTML and CSS special characters using \&lt;hex&gt;+ escapes.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>func cssEscaper(args ...any) string {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	s, _ := stringify(args...)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	var b strings.Builder
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	r, w, written := rune(0), 0, 0
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i += w {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		<span class="comment">// See comment in htmlEscaper.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		r, w = utf8.DecodeRuneInString(s[i:])
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		var repl string
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		switch {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		case int(r) &lt; len(cssReplacementTable) &amp;&amp; cssReplacementTable[r] != &#34;&#34;:
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			repl = cssReplacementTable[r]
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		default:
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			continue
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		if written == 0 {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			b.Grow(len(s))
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		b.WriteString(s[written:i])
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		b.WriteString(repl)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		written = i + w
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		if repl != `\\` &amp;&amp; (written == len(s) || isHex(s[written]) || isCSSSpace(s[written])) {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			b.WriteByte(&#39; &#39;)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	if written == 0 {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		return s
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	b.WriteString(s[written:])
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	return b.String()
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>var cssReplacementTable = []string{
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	0:    `\0`,
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	&#39;\t&#39;: `\9`,
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	&#39;\n&#39;: `\a`,
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	&#39;\f&#39;: `\c`,
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	&#39;\r&#39;: `\d`,
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// Encode HTML specials as hex so the output can be embedded</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// in HTML attributes without further encoding.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	&#39;&#34;&#39;:  `\22`,
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	&#39;&amp;&#39;:  `\26`,
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	&#39;\&#39;&#39;: `\27`,
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	&#39;(&#39;:  `\28`,
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	&#39;)&#39;:  `\29`,
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	&#39;+&#39;:  `\2b`,
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	&#39;/&#39;:  `\2f`,
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	&#39;:&#39;:  `\3a`,
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	&#39;;&#39;:  `\3b`,
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	&#39;&lt;&#39;:  `\3c`,
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	&#39;&gt;&#39;:  `\3e`,
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	&#39;\\&#39;: `\\`,
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	&#39;{&#39;:  `\7b`,
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	&#39;}&#39;:  `\7d`,
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>var expressionBytes = []byte(&#34;expression&#34;)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>var mozBindingBytes = []byte(&#34;mozbinding&#34;)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// cssValueFilter allows innocuous CSS values in the output including CSS</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// quantities (10px or 25%), ID or class literals (#foo, .bar), keyword values</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// (inherit, blue), and colors (#888).</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">// It filters out unsafe values, such as those that affect token boundaries,</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">// and anything that might execute scripts.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>func cssValueFilter(args ...any) string {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	s, t := stringify(args...)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	if t == contentTypeCSS {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		return s
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	b, id := decodeCSS([]byte(s)), make([]byte, 0, 64)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// CSS3 error handling is specified as honoring string boundaries per</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// https://www.w3.org/TR/css3-syntax/#error-handling :</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">//     Malformed declarations. User agents must handle unexpected</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">//     tokens encountered while parsing a declaration by reading until</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">//     the end of the declaration, while observing the rules for</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">//     matching pairs of (), [], {}, &#34;&#34;, and &#39;&#39;, and correctly handling</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	<span class="comment">//     escapes. For example, a malformed declaration may be missing a</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">//     property, colon (:) or value.</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// So we need to make sure that values do not have mismatched bracket</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// or quote characters to prevent the browser from restarting parsing</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// inside a string that might embed JavaScript source.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	for i, c := range b {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		switch c {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		case 0, &#39;&#34;&#39;, &#39;\&#39;&#39;, &#39;(&#39;, &#39;)&#39;, &#39;/&#39;, &#39;;&#39;, &#39;@&#39;, &#39;[&#39;, &#39;\\&#39;, &#39;]&#39;, &#39;`&#39;, &#39;{&#39;, &#39;}&#39;, &#39;&lt;&#39;, &#39;&gt;&#39;:
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			return filterFailsafe
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		case &#39;-&#39;:
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			<span class="comment">// Disallow &lt;!-- or --&gt;.</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			<span class="comment">// -- should not appear in valid identifiers.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			if i != 0 &amp;&amp; b[i-1] == &#39;-&#39; {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>				return filterFailsafe
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		default:
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			if c &lt; utf8.RuneSelf &amp;&amp; isCSSNmchar(rune(c)) {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>				id = append(id, c)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	id = bytes.ToLower(id)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	if bytes.Contains(id, expressionBytes) || bytes.Contains(id, mozBindingBytes) {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		return filterFailsafe
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	return string(b)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>
</pre><p><a href="css.go?m=text">View as plain text</a></p>

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
