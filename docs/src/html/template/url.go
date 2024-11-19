<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/html/template/url.go - Go Documentation Server</title>

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
<a href="url.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/html">html</a>/<a href="http://localhost:8080/src/html/template">template</a>/<span class="text-muted">url.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// urlFilter returns its input unless it contains an unsafe scheme in which</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// case it defangs the entire URL.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// Schemes that cause unintended side effects that are irreversible without user</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// interaction are considered unsafe. For example, clicking on a &#34;javascript:&#34;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// link can immediately trigger JavaScript code execution.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// This filter conservatively assumes that all schemes other than the following</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// are unsafe:</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//   - http:   Navigates to a new website, and may open a new window or tab.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//     These side effects can be reversed by navigating back to the</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//     previous website, or closing the window or tab. No irreversible</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//     changes will take place without further user interaction with</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//     the new website.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//   - https:  Same as http.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//   - mailto: Opens an email program and starts a new draft. This side effect</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//     is not irreversible until the user explicitly clicks send; it</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//     can be undone by closing the email program.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// To allow URLs containing other schemes to bypass this filter, developers must</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// explicitly indicate that such a URL is expected and safe by encapsulating it</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// in a template.URL value.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>func urlFilter(args ...any) string {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	s, t := stringify(args...)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	if t == contentTypeURL {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		return s
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	if !isSafeURL(s) {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		return &#34;#&#34; + filterFailsafe
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	return s
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// isSafeURL is true if s is a relative URL or if URL has a protocol in</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// (http, https, mailto).</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>func isSafeURL(s string) bool {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	if protocol, _, ok := strings.Cut(s, &#34;:&#34;); ok &amp;&amp; !strings.Contains(protocol, &#34;/&#34;) {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		if !strings.EqualFold(protocol, &#34;http&#34;) &amp;&amp; !strings.EqualFold(protocol, &#34;https&#34;) &amp;&amp; !strings.EqualFold(protocol, &#34;mailto&#34;) {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			return false
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	return true
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// urlEscaper produces an output that can be embedded in a URL query.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// The output can be embedded in an HTML attribute without further escaping.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func urlEscaper(args ...any) string {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	return urlProcessor(false, args...)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// urlNormalizer normalizes URL content so it can be embedded in a quote-delimited</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// string or parenthesis delimited url(...).</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// The normalizer does not encode all HTML specials. Specifically, it does not</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// encode &#39;&amp;&#39; so correct embedding in an HTML attribute requires escaping of</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// &#39;&amp;&#39; to &#39;&amp;amp;&#39;.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>func urlNormalizer(args ...any) string {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	return urlProcessor(true, args...)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// urlProcessor normalizes (when norm is true) or escapes its input to produce</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// a valid hierarchical or opaque URL part.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>func urlProcessor(norm bool, args ...any) string {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	s, t := stringify(args...)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	if t == contentTypeURL {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		norm = true
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	var b strings.Builder
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	if processURLOnto(s, norm, &amp;b) {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		return b.String()
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	return s
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// processURLOnto appends a normalized URL corresponding to its input to b</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// and reports whether the appended content differs from s.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>func processURLOnto(s string, norm bool, b *strings.Builder) bool {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	b.Grow(len(s) + 16)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	written := 0
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// The byte loop below assumes that all URLs use UTF-8 as the</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// content-encoding. This is similar to the URI to IRI encoding scheme</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// defined in section 3.1 of  RFC 3987, and behaves the same as the</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// EcmaScript builtin encodeURIComponent.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// It should not cause any misencoding of URLs in pages with</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// Content-type: text/html;charset=UTF-8.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	for i, n := 0, len(s); i &lt; n; i++ {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		c := s[i]
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		switch c {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		<span class="comment">// Single quote and parens are sub-delims in RFC 3986, but we</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		<span class="comment">// escape them so the output can be embedded in single</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		<span class="comment">// quoted attributes and unquoted CSS url(...) constructs.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		<span class="comment">// Single quotes are reserved in URLs, but are only used in</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		<span class="comment">// the obsolete &#34;mark&#34; rule in an appendix in RFC 3986</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		<span class="comment">// so can be safely encoded.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		case &#39;!&#39;, &#39;#&#39;, &#39;$&#39;, &#39;&amp;&#39;, &#39;*&#39;, &#39;+&#39;, &#39;,&#39;, &#39;/&#39;, &#39;:&#39;, &#39;;&#39;, &#39;=&#39;, &#39;?&#39;, &#39;@&#39;, &#39;[&#39;, &#39;]&#39;:
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			if norm {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>				continue
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		<span class="comment">// Unreserved according to RFC 3986 sec 2.3</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		<span class="comment">// &#34;For consistency, percent-encoded octets in the ranges of</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		<span class="comment">// ALPHA (%41-%5A and %61-%7A), DIGIT (%30-%39), hyphen (%2D),</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		<span class="comment">// period (%2E), underscore (%5F), or tilde (%7E) should not be</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		<span class="comment">// created by URI producers</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		case &#39;-&#39;, &#39;.&#39;, &#39;_&#39;, &#39;~&#39;:
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			continue
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		case &#39;%&#39;:
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			<span class="comment">// When normalizing do not re-encode valid escapes.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			if norm &amp;&amp; i+2 &lt; len(s) &amp;&amp; isHex(s[i+1]) &amp;&amp; isHex(s[i+2]) {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>				continue
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		default:
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			<span class="comment">// Unreserved according to RFC 3986 sec 2.3</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			if &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;z&#39; {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>				continue
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			if &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;Z&#39; {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>				continue
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			if &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39; {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>				continue
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		b.WriteString(s[written:i])
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		fmt.Fprintf(b, &#34;%%%02x&#34;, c)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		written = i + 1
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	b.WriteString(s[written:])
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	return written != 0
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// Filters and normalizes srcset values which are comma separated</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// URLs followed by metadata.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>func srcsetFilterAndEscaper(args ...any) string {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	s, t := stringify(args...)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	switch t {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	case contentTypeSrcset:
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		return s
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	case contentTypeURL:
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		<span class="comment">// Normalizing gets rid of all HTML whitespace</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		<span class="comment">// which separate the image URL from its metadata.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		var b strings.Builder
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		if processURLOnto(s, true, &amp;b) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			s = b.String()
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		<span class="comment">// Additionally, commas separate one source from another.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		return strings.ReplaceAll(s, &#34;,&#34;, &#34;%2c&#34;)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	var b strings.Builder
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	written := 0
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i++ {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		if s[i] == &#39;,&#39; {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			filterSrcsetElement(s, written, i, &amp;b)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			b.WriteString(&#34;,&#34;)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			written = i + 1
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	filterSrcsetElement(s, written, len(s), &amp;b)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	return b.String()
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// Derived from https://play.golang.org/p/Dhmj7FORT5</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>const htmlSpaceAndASCIIAlnumBytes = &#34;\x00\x36\x00\x00\x01\x00\xff\x03\xfe\xff\xff\x07\xfe\xff\xff\x07&#34;
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// isHTMLSpace is true iff c is a whitespace character per</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// https://infra.spec.whatwg.org/#ascii-whitespace</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>func isHTMLSpace(c byte) bool {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	return (c &lt;= 0x20) &amp;&amp; 0 != (htmlSpaceAndASCIIAlnumBytes[c&gt;&gt;3]&amp;(1&lt;&lt;uint(c&amp;0x7)))
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>func isHTMLSpaceOrASCIIAlnum(c byte) bool {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	return (c &lt; 0x80) &amp;&amp; 0 != (htmlSpaceAndASCIIAlnumBytes[c&gt;&gt;3]&amp;(1&lt;&lt;uint(c&amp;0x7)))
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>func filterSrcsetElement(s string, left int, right int, b *strings.Builder) {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	start := left
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	for start &lt; right &amp;&amp; isHTMLSpace(s[start]) {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		start++
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	end := right
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	for i := start; i &lt; right; i++ {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		if isHTMLSpace(s[i]) {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			end = i
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			break
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	if url := s[start:end]; isSafeURL(url) {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		<span class="comment">// If image metadata is only spaces or alnums then</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		<span class="comment">// we don&#39;t need to URL normalize it.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		metadataOk := true
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		for i := end; i &lt; right; i++ {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>			if !isHTMLSpaceOrASCIIAlnum(s[i]) {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>				metadataOk = false
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>				break
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		if metadataOk {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			b.WriteString(s[left:start])
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			processURLOnto(url, true, b)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			b.WriteString(s[end:right])
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			return
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	b.WriteString(&#34;#&#34;)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	b.WriteString(filterFailsafe)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
</pre><p><a href="url.go?m=text">View as plain text</a></p>

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
