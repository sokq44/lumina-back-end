<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/url/url.go - Go Documentation Server</title>

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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<a href="http://localhost:8080/src/net/url">url</a>/<span class="text-muted">url.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net/url">net/url</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package url parses URLs and implements query escaping.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>package url
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// See RFC 3986. This package generally follows RFC 3986, except where</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// it deviates for compatibility reasons. When sending changes, first</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// search old issues for history on decisions. Unit tests should also</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// contain references to issue numbers with details.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>import (
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;path&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// Error reports an error and the operation and URL that caused it.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>type Error struct {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	Op  string
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	URL string
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	Err error
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>func (e *Error) Unwrap() error { return e.Err }
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>func (e *Error) Error() string { return fmt.Sprintf(&#34;%s %q: %s&#34;, e.Op, e.URL, e.Err) }
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>func (e *Error) Timeout() bool {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	t, ok := e.Err.(interface {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		Timeout() bool
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	})
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	return ok &amp;&amp; t.Timeout()
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>func (e *Error) Temporary() bool {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	t, ok := e.Err.(interface {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		Temporary() bool
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	})
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	return ok &amp;&amp; t.Temporary()
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>const upperhex = &#34;0123456789ABCDEF&#34;
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>func ishex(c byte) bool {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	switch {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	case &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39;:
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		return true
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	case &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;f&#39;:
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		return true
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	case &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;F&#39;:
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		return true
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	return false
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func unhex(c byte) byte {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	switch {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	case &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39;:
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		return c - &#39;0&#39;
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	case &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;f&#39;:
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		return c - &#39;a&#39; + 10
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	case &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;F&#39;:
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		return c - &#39;A&#39; + 10
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	return 0
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>type encoding int
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>const (
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	encodePath encoding = 1 + iota
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	encodePathSegment
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	encodeHost
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	encodeZone
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	encodeUserPassword
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	encodeQueryComponent
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	encodeFragment
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>type EscapeError string
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>func (e EscapeError) Error() string {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	return &#34;invalid URL escape &#34; + strconv.Quote(string(e))
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>type InvalidHostError string
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>func (e InvalidHostError) Error() string {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	return &#34;invalid character &#34; + strconv.Quote(string(e)) + &#34; in host name&#34;
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// Return true if the specified character should be escaped when</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// appearing in a URL string, according to RFC 3986.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// Please be informed that for now shouldEscape does not check all</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// reserved characters correctly. See golang.org/issue/5684.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>func shouldEscape(c byte, mode encoding) bool {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// §2.3 Unreserved characters (alphanum)</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	if &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;z&#39; || &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;Z&#39; || &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39; {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		return false
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	if mode == encodeHost || mode == encodeZone {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		<span class="comment">// §3.2.2 Host allows</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		<span class="comment">//	sub-delims = &#34;!&#34; / &#34;$&#34; / &#34;&amp;&#34; / &#34;&#39;&#34; / &#34;(&#34; / &#34;)&#34; / &#34;*&#34; / &#34;+&#34; / &#34;,&#34; / &#34;;&#34; / &#34;=&#34;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		<span class="comment">// as part of reg-name.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		<span class="comment">// We add : because we include :port as part of host.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		<span class="comment">// We add [ ] because we include [ipv6]:port as part of host.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		<span class="comment">// We add &lt; &gt; because they&#39;re the only characters left that</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		<span class="comment">// we could possibly allow, and Parse will reject them if we</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		<span class="comment">// escape them (because hosts can&#39;t use %-encoding for</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		<span class="comment">// ASCII bytes).</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		switch c {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		case &#39;!&#39;, &#39;$&#39;, &#39;&amp;&#39;, &#39;\&#39;&#39;, &#39;(&#39;, &#39;)&#39;, &#39;*&#39;, &#39;+&#39;, &#39;,&#39;, &#39;;&#39;, &#39;=&#39;, &#39;:&#39;, &#39;[&#39;, &#39;]&#39;, &#39;&lt;&#39;, &#39;&gt;&#39;, &#39;&#34;&#39;:
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			return false
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	switch c {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	case &#39;-&#39;, &#39;_&#39;, &#39;.&#39;, &#39;~&#39;: <span class="comment">// §2.3 Unreserved characters (mark)</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		return false
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	case &#39;$&#39;, &#39;&amp;&#39;, &#39;+&#39;, &#39;,&#39;, &#39;/&#39;, &#39;:&#39;, &#39;;&#39;, &#39;=&#39;, &#39;?&#39;, &#39;@&#39;: <span class="comment">// §2.2 Reserved characters (reserved)</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		<span class="comment">// Different sections of the URL allow a few of</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		<span class="comment">// the reserved characters to appear unescaped.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		switch mode {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		case encodePath: <span class="comment">// §3.3</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			<span class="comment">// The RFC allows : @ &amp; = + $ but saves / ; , for assigning</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			<span class="comment">// meaning to individual path segments. This package</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			<span class="comment">// only manipulates the path as a whole, so we allow those</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			<span class="comment">// last three as well. That leaves only ? to escape.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			return c == &#39;?&#39;
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		case encodePathSegment: <span class="comment">// §3.3</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			<span class="comment">// The RFC allows : @ &amp; = + $ but saves / ; , for assigning</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			<span class="comment">// meaning to individual path segments.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			return c == &#39;/&#39; || c == &#39;;&#39; || c == &#39;,&#39; || c == &#39;?&#39;
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		case encodeUserPassword: <span class="comment">// §3.2.1</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			<span class="comment">// The RFC allows &#39;;&#39;, &#39;:&#39;, &#39;&amp;&#39;, &#39;=&#39;, &#39;+&#39;, &#39;$&#39;, and &#39;,&#39; in</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			<span class="comment">// userinfo, so we must escape only &#39;@&#39;, &#39;/&#39;, and &#39;?&#39;.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			<span class="comment">// The parsing of userinfo treats &#39;:&#39; as special so we must escape</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			<span class="comment">// that too.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			return c == &#39;@&#39; || c == &#39;/&#39; || c == &#39;?&#39; || c == &#39;:&#39;
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		case encodeQueryComponent: <span class="comment">// §3.4</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			<span class="comment">// The RFC reserves (so we must escape) everything.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			return true
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		case encodeFragment: <span class="comment">// §4.1</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			<span class="comment">// The RFC text is silent but the grammar allows</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			<span class="comment">// everything, so escape nothing.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			return false
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if mode == encodeFragment {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		<span class="comment">// RFC 3986 §2.2 allows not escaping sub-delims. A subset of sub-delims are</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		<span class="comment">// included in reserved from RFC 2396 §2.2. The remaining sub-delims do not</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		<span class="comment">// need to be escaped. To minimize potential breakage, we apply two restrictions:</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		<span class="comment">// (1) we always escape sub-delims outside of the fragment, and (2) we always</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		<span class="comment">// escape single quote to avoid breaking callers that had previously assumed that</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		<span class="comment">// single quotes would be escaped. See issue #19917.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		switch c {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		case &#39;!&#39;, &#39;(&#39;, &#39;)&#39;, &#39;*&#39;:
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			return false
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// Everything else must be escaped.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	return true
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">// QueryUnescape does the inverse transformation of [QueryEscape],</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">// converting each 3-byte encoded substring of the form &#34;%AB&#34; into the</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span><span class="comment">// hex-decoded byte 0xAB.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">// It returns an error if any % is not followed by two hexadecimal</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">// digits.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>func QueryUnescape(s string) (string, error) {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	return unescape(s, encodeQueryComponent)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// PathUnescape does the inverse transformation of [PathEscape],</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// converting each 3-byte encoded substring of the form &#34;%AB&#34; into the</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// hex-decoded byte 0xAB. It returns an error if any % is not followed</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// by two hexadecimal digits.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// PathUnescape is identical to [QueryUnescape] except that it does not</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">// unescape &#39;+&#39; to &#39; &#39; (space).</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>func PathUnescape(s string) (string, error) {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	return unescape(s, encodePathSegment)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// unescape unescapes a string; the mode specifies</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// which section of the URL string is being unescaped.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>func unescape(s string, mode encoding) (string, error) {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// Count %, check that they&#39;re well-formed.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	n := 0
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	hasPlus := false
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		switch s[i] {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		case &#39;%&#39;:
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			n++
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			if i+2 &gt;= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>				s = s[i:]
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>				if len(s) &gt; 3 {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>					s = s[:3]
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>				}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>				return &#34;&#34;, EscapeError(s)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			<span class="comment">// Per https://tools.ietf.org/html/rfc3986#page-21</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			<span class="comment">// in the host component %-encoding can only be used</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			<span class="comment">// for non-ASCII bytes.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			<span class="comment">// But https://tools.ietf.org/html/rfc6874#section-2</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			<span class="comment">// introduces %25 being allowed to escape a percent sign</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			<span class="comment">// in IPv6 scoped-address literals. Yay.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			if mode == encodeHost &amp;&amp; unhex(s[i+1]) &lt; 8 &amp;&amp; s[i:i+3] != &#34;%25&#34; {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>				return &#34;&#34;, EscapeError(s[i : i+3])
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			if mode == encodeZone {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>				<span class="comment">// RFC 6874 says basically &#34;anything goes&#34; for zone identifiers</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>				<span class="comment">// and that even non-ASCII can be redundantly escaped,</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>				<span class="comment">// but it seems prudent to restrict %-escaped bytes here to those</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>				<span class="comment">// that are valid host name bytes in their unescaped form.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>				<span class="comment">// That is, you can use escaping in the zone identifier but not</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>				<span class="comment">// to introduce bytes you couldn&#39;t just write directly.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>				<span class="comment">// But Windows puts spaces here! Yay.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>				v := unhex(s[i+1])&lt;&lt;4 | unhex(s[i+2])
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>				if s[i:i+3] != &#34;%25&#34; &amp;&amp; v != &#39; &#39; &amp;&amp; shouldEscape(v, encodeHost) {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>					return &#34;&#34;, EscapeError(s[i : i+3])
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>				}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			i += 3
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		case &#39;+&#39;:
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			hasPlus = mode == encodeQueryComponent
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			i++
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		default:
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			if (mode == encodeHost || mode == encodeZone) &amp;&amp; s[i] &lt; 0x80 &amp;&amp; shouldEscape(s[i], mode) {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>				return &#34;&#34;, InvalidHostError(s[i : i+1])
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			i++
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	if n == 0 &amp;&amp; !hasPlus {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		return s, nil
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	var t strings.Builder
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	t.Grow(len(s) - 2*n)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i++ {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		switch s[i] {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		case &#39;%&#39;:
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			t.WriteByte(unhex(s[i+1])&lt;&lt;4 | unhex(s[i+2]))
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			i += 2
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		case &#39;+&#39;:
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			if mode == encodeQueryComponent {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				t.WriteByte(&#39; &#39;)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			} else {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>				t.WriteByte(&#39;+&#39;)
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		default:
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			t.WriteByte(s[i])
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		}
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	return t.String(), nil
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span><span class="comment">// QueryEscape escapes the string so it can be safely placed</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span><span class="comment">// inside a [URL] query.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>func QueryEscape(s string) string {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	return escape(s, encodeQueryComponent)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">// PathEscape escapes the string so it can be safely placed inside a [URL] path segment,</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">// replacing special characters (including /) with %XX sequences as needed.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>func PathEscape(s string) string {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	return escape(s, encodePathSegment)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>func escape(s string, mode encoding) string {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	spaceCount, hexCount := 0, 0
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i++ {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		c := s[i]
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		if shouldEscape(c, mode) {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			if c == &#39; &#39; &amp;&amp; mode == encodeQueryComponent {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>				spaceCount++
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			} else {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>				hexCount++
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	if spaceCount == 0 &amp;&amp; hexCount == 0 {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		return s
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	var buf [64]byte
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	var t []byte
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	required := len(s) + 2*hexCount
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	if required &lt;= len(buf) {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		t = buf[:required]
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	} else {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		t = make([]byte, required)
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	}
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	if hexCount == 0 {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		copy(t, s)
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		for i := 0; i &lt; len(s); i++ {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>			if s[i] == &#39; &#39; {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>				t[i] = &#39;+&#39;
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		return string(t)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	j := 0
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i++ {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		switch c := s[i]; {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		case c == &#39; &#39; &amp;&amp; mode == encodeQueryComponent:
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			t[j] = &#39;+&#39;
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			j++
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		case shouldEscape(c, mode):
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>			t[j] = &#39;%&#39;
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>			t[j+1] = upperhex[c&gt;&gt;4]
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			t[j+2] = upperhex[c&amp;15]
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			j += 3
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		default:
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>			t[j] = s[i]
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>			j++
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	return string(t)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">// A URL represents a parsed URL (technically, a URI reference).</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">// The general form represented is:</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">//	[scheme:][//[userinfo@]host][/]path[?query][#fragment]</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// URLs that do not start with a slash after the scheme are interpreted as:</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">//	scheme:opaque[?query][#fragment]</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// The Host field contains the host and port subcomponents of the URL.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// When the port is present, it is separated from the host with a colon.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// When the host is an IPv6 address, it must be enclosed in square brackets:</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// &#34;[fe80::1]:80&#34;. The [net.JoinHostPort] function combines a host and port</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// into a string suitable for the Host field, adding square brackets to</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// the host when necessary.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">// Note that the Path field is stored in decoded form: /%47%6f%2f becomes /Go/.</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// A consequence is that it is impossible to tell which slashes in the Path were</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span><span class="comment">// slashes in the raw URL and which were %2f. This distinction is rarely important,</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// but when it is, the code should use the [URL.EscapedPath] method, which preserves</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// the original encoding of Path.</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">// The RawPath field is an optional field which is only set when the default</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// encoding of Path is different from the escaped path. See the EscapedPath method</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// for more details.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">// URL&#39;s String method uses the EscapedPath method to obtain the path.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>type URL struct {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	Scheme      string
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	Opaque      string    <span class="comment">// encoded opaque data</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	User        *Userinfo <span class="comment">// username and password information</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	Host        string    <span class="comment">// host or host:port (see Hostname and Port methods)</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	Path        string    <span class="comment">// path (relative paths may omit leading slash)</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	RawPath     string    <span class="comment">// encoded path hint (see EscapedPath method)</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	OmitHost    bool      <span class="comment">// do not emit empty host (authority)</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	ForceQuery  bool      <span class="comment">// append a query (&#39;?&#39;) even if RawQuery is empty</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	RawQuery    string    <span class="comment">// encoded query values, without &#39;?&#39;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	Fragment    string    <span class="comment">// fragment for references, without &#39;#&#39;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	RawFragment string    <span class="comment">// encoded fragment hint (see EscapedFragment method)</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span><span class="comment">// User returns a [Userinfo] containing the provided username</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span><span class="comment">// and no password set.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>func User(username string) *Userinfo {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	return &amp;Userinfo{username, &#34;&#34;, false}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span><span class="comment">// UserPassword returns a [Userinfo] containing the provided username</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span><span class="comment">// and password.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span><span class="comment">// This functionality should only be used with legacy web sites.</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">// RFC 2396 warns that interpreting Userinfo this way</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">// “is NOT RECOMMENDED, because the passing of authentication</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// information in clear text (such as URI) has proven to be a</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">// security risk in almost every case where it has been used.”</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>func UserPassword(username, password string) *Userinfo {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	return &amp;Userinfo{username, password, true}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span><span class="comment">// The Userinfo type is an immutable encapsulation of username and</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span><span class="comment">// password details for a [URL]. An existing Userinfo value is guaranteed</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span><span class="comment">// to have a username set (potentially empty, as allowed by RFC 2396),</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span><span class="comment">// and optionally a password.</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>type Userinfo struct {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	username    string
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	password    string
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	passwordSet bool
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span><span class="comment">// Username returns the username.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>func (u *Userinfo) Username() string {
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	if u == nil {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	return u.username
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span><span class="comment">// Password returns the password in case it is set, and whether it is set.</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>func (u *Userinfo) Password() (string, bool) {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	if u == nil {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		return &#34;&#34;, false
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	return u.password, u.passwordSet
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span><span class="comment">// String returns the encoded userinfo information in the standard form</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span><span class="comment">// of &#34;username[:password]&#34;.</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>func (u *Userinfo) String() string {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	if u == nil {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	s := escape(u.username, encodeUserPassword)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	if u.passwordSet {
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		s += &#34;:&#34; + escape(u.password, encodeUserPassword)
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	}
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	return s
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span><span class="comment">// Maybe rawURL is of the form scheme:path.</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span><span class="comment">// (Scheme must be [a-zA-Z][a-zA-Z0-9+.-]*)</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span><span class="comment">// If so, return scheme, path; else return &#34;&#34;, rawURL.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>func getScheme(rawURL string) (scheme, path string, err error) {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	for i := 0; i &lt; len(rawURL); i++ {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		c := rawURL[i]
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		switch {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		case &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;z&#39; || &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;Z&#39;:
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		<span class="comment">// do nothing</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		case &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39; || c == &#39;+&#39; || c == &#39;-&#39; || c == &#39;.&#39;:
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>			if i == 0 {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>				return &#34;&#34;, rawURL, nil
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		case c == &#39;:&#39;:
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>			if i == 0 {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>				return &#34;&#34;, &#34;&#34;, errors.New(&#34;missing protocol scheme&#34;)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>			}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			return rawURL[:i], rawURL[i+1:], nil
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		default:
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>			<span class="comment">// we have encountered an invalid character,</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			<span class="comment">// so there is no valid scheme</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			return &#34;&#34;, rawURL, nil
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	return &#34;&#34;, rawURL, nil
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span><span class="comment">// Parse parses a raw url into a [URL] structure.</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span><span class="comment">// The url may be relative (a path, without a host) or absolute</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span><span class="comment">// (starting with a scheme). Trying to parse a hostname and path</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span><span class="comment">// without a scheme is invalid but may not necessarily return an</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span><span class="comment">// error, due to parsing ambiguities.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>func Parse(rawURL string) (*URL, error) {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	<span class="comment">// Cut off #frag</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	u, frag, _ := strings.Cut(rawURL, &#34;#&#34;)
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	url, err := parse(u, false)
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	if err != nil {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		return nil, &amp;Error{&#34;parse&#34;, u, err}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	}
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	if frag == &#34;&#34; {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		return url, nil
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	if err = url.setFragment(frag); err != nil {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		return nil, &amp;Error{&#34;parse&#34;, rawURL, err}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	return url, nil
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span><span class="comment">// ParseRequestURI parses a raw url into a [URL] structure. It assumes that</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span><span class="comment">// url was received in an HTTP request, so the url is interpreted</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span><span class="comment">// only as an absolute URI or an absolute path.</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span><span class="comment">// The string url is assumed not to have a #fragment suffix.</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span><span class="comment">// (Web browsers strip #fragment before sending the URL to a web server.)</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>func ParseRequestURI(rawURL string) (*URL, error) {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	url, err := parse(rawURL, true)
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	if err != nil {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		return nil, &amp;Error{&#34;parse&#34;, rawURL, err}
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	return url, nil
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>}
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span><span class="comment">// parse parses a URL from a string in one of two contexts. If</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span><span class="comment">// viaRequest is true, the URL is assumed to have arrived via an HTTP request,</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span><span class="comment">// in which case only absolute URLs or path-absolute relative URLs are allowed.</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span><span class="comment">// If viaRequest is false, all forms of relative URLs are allowed.</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>func parse(rawURL string, viaRequest bool) (*URL, error) {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	var rest string
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	var err error
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	if stringContainsCTLByte(rawURL) {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		return nil, errors.New(&#34;net/url: invalid control character in URL&#34;)
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	if rawURL == &#34;&#34; &amp;&amp; viaRequest {
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		return nil, errors.New(&#34;empty url&#34;)
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	url := new(URL)
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	if rawURL == &#34;*&#34; {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		url.Path = &#34;*&#34;
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		return url, nil
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	<span class="comment">// Split off possible leading &#34;http:&#34;, &#34;mailto:&#34;, etc.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	<span class="comment">// Cannot contain escaped characters.</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	if url.Scheme, rest, err = getScheme(rawURL); err != nil {
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		return nil, err
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	}
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	url.Scheme = strings.ToLower(url.Scheme)
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	if strings.HasSuffix(rest, &#34;?&#34;) &amp;&amp; strings.Count(rest, &#34;?&#34;) == 1 {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		url.ForceQuery = true
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		rest = rest[:len(rest)-1]
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	} else {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		rest, url.RawQuery, _ = strings.Cut(rest, &#34;?&#34;)
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	if !strings.HasPrefix(rest, &#34;/&#34;) {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		if url.Scheme != &#34;&#34; {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			<span class="comment">// We consider rootless paths per RFC 3986 as opaque.</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>			url.Opaque = rest
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			return url, nil
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		if viaRequest {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>			return nil, errors.New(&#34;invalid URI for request&#34;)
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		}
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		<span class="comment">// Avoid confusion with malformed schemes, like cache_object:foo/bar.</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		<span class="comment">// See golang.org/issue/16822.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		<span class="comment">// RFC 3986, §3.3:</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		<span class="comment">// In addition, a URI reference (Section 4.1) may be a relative-path reference,</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		<span class="comment">// in which case the first path segment cannot contain a colon (&#34;:&#34;) character.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		if segment, _, _ := strings.Cut(rest, &#34;/&#34;); strings.Contains(segment, &#34;:&#34;) {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			<span class="comment">// First path segment has colon. Not allowed in relative URL.</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>			return nil, errors.New(&#34;first path segment in URL cannot contain colon&#34;)
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	if (url.Scheme != &#34;&#34; || !viaRequest &amp;&amp; !strings.HasPrefix(rest, &#34;///&#34;)) &amp;&amp; strings.HasPrefix(rest, &#34;//&#34;) {
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		var authority string
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		authority, rest = rest[2:], &#34;&#34;
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		if i := strings.Index(authority, &#34;/&#34;); i &gt;= 0 {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>			authority, rest = authority[:i], authority[i:]
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		}
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		url.User, url.Host, err = parseAuthority(authority)
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		if err != nil {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>			return nil, err
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	} else if url.Scheme != &#34;&#34; &amp;&amp; strings.HasPrefix(rest, &#34;/&#34;) {
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		<span class="comment">// OmitHost is set to true when rawURL has an empty host (authority).</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		<span class="comment">// See golang.org/issue/46059.</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		url.OmitHost = true
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	<span class="comment">// Set Path and, optionally, RawPath.</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	<span class="comment">// RawPath is a hint of the encoding of Path. We don&#39;t want to set it if</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	<span class="comment">// the default escaping of Path is equivalent, to help make sure that people</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	<span class="comment">// don&#39;t rely on it in general.</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	if err := url.setPath(rest); err != nil {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		return nil, err
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	}
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	return url, nil
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>func parseAuthority(authority string) (user *Userinfo, host string, err error) {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	i := strings.LastIndex(authority, &#34;@&#34;)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		host, err = parseHost(authority)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	} else {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		host, err = parseHost(authority[i+1:])
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	}
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	if err != nil {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		return nil, &#34;&#34;, err
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	}
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		return nil, host, nil
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	userinfo := authority[:i]
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	if !validUserinfo(userinfo) {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		return nil, &#34;&#34;, errors.New(&#34;net/url: invalid userinfo&#34;)
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	if !strings.Contains(userinfo, &#34;:&#34;) {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		if userinfo, err = unescape(userinfo, encodeUserPassword); err != nil {
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			return nil, &#34;&#34;, err
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		user = User(userinfo)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	} else {
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		username, password, _ := strings.Cut(userinfo, &#34;:&#34;)
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		if username, err = unescape(username, encodeUserPassword); err != nil {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>			return nil, &#34;&#34;, err
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		if password, err = unescape(password, encodeUserPassword); err != nil {
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>			return nil, &#34;&#34;, err
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		}
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		user = UserPassword(username, password)
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	}
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	return user, host, nil
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>}
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span><span class="comment">// parseHost parses host as an authority without user</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span><span class="comment">// information. That is, as host[:port].</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>func parseHost(host string) (string, error) {
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	if strings.HasPrefix(host, &#34;[&#34;) {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		<span class="comment">// Parse an IP-Literal in RFC 3986 and RFC 6874.</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		<span class="comment">// E.g., &#34;[fe80::1]&#34;, &#34;[fe80::1%25en0]&#34;, &#34;[fe80::1]:80&#34;.</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		i := strings.LastIndex(host, &#34;]&#34;)
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		if i &lt; 0 {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>			return &#34;&#34;, errors.New(&#34;missing &#39;]&#39; in host&#34;)
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		colonPort := host[i+1:]
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>		if !validOptionalPort(colonPort) {
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>			return &#34;&#34;, fmt.Errorf(&#34;invalid port %q after host&#34;, colonPort)
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		}
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		<span class="comment">// RFC 6874 defines that %25 (%-encoded percent) introduces</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		<span class="comment">// the zone identifier, and the zone identifier can use basically</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		<span class="comment">// any %-encoding it likes. That&#39;s different from the host, which</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		<span class="comment">// can only %-encode non-ASCII bytes.</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		<span class="comment">// We do impose some restrictions on the zone, to avoid stupidity</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		<span class="comment">// like newlines.</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>		zone := strings.Index(host[:i], &#34;%25&#34;)
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		if zone &gt;= 0 {
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>			host1, err := unescape(host[:zone], encodeHost)
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>			if err != nil {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>				return &#34;&#34;, err
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>			}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>			host2, err := unescape(host[zone:i], encodeZone)
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>			if err != nil {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>				return &#34;&#34;, err
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>			}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>			host3, err := unescape(host[i:], encodeHost)
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>			if err != nil {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>				return &#34;&#34;, err
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>			}
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>			return host1 + host2 + host3, nil
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		}
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	} else if i := strings.LastIndex(host, &#34;:&#34;); i != -1 {
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		colonPort := host[i:]
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		if !validOptionalPort(colonPort) {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>			return &#34;&#34;, fmt.Errorf(&#34;invalid port %q after host&#34;, colonPort)
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		}
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	}
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	var err error
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	if host, err = unescape(host, encodeHost); err != nil {
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		return &#34;&#34;, err
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	}
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	return host, nil
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span><span class="comment">// setPath sets the Path and RawPath fields of the URL based on the provided</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span><span class="comment">// escaped path p. It maintains the invariant that RawPath is only specified</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span><span class="comment">// when it differs from the default encoding of the path.</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span><span class="comment">// For example:</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span><span class="comment">// - setPath(&#34;/foo/bar&#34;)   will set Path=&#34;/foo/bar&#34; and RawPath=&#34;&#34;</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span><span class="comment">// - setPath(&#34;/foo%2fbar&#34;) will set Path=&#34;/foo/bar&#34; and RawPath=&#34;/foo%2fbar&#34;</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span><span class="comment">// setPath will return an error only if the provided path contains an invalid</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span><span class="comment">// escaping.</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>func (u *URL) setPath(p string) error {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	path, err := unescape(p, encodePath)
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	if err != nil {
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		return err
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	u.Path = path
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	if escp := escape(path, encodePath); p == escp {
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		<span class="comment">// Default encoding is fine.</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		u.RawPath = &#34;&#34;
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	} else {
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		u.RawPath = p
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	}
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	return nil
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>}
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span><span class="comment">// EscapedPath returns the escaped form of u.Path.</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span><span class="comment">// In general there are multiple possible escaped forms of any path.</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span><span class="comment">// EscapedPath returns u.RawPath when it is a valid escaping of u.Path.</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span><span class="comment">// Otherwise EscapedPath ignores u.RawPath and computes an escaped</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span><span class="comment">// form on its own.</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span><span class="comment">// The [URL.String] and [URL.RequestURI] methods use EscapedPath to construct</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span><span class="comment">// their results.</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span><span class="comment">// In general, code should call EscapedPath instead of</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span><span class="comment">// reading u.RawPath directly.</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>func (u *URL) EscapedPath() string {
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	if u.RawPath != &#34;&#34; &amp;&amp; validEncoded(u.RawPath, encodePath) {
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		p, err := unescape(u.RawPath, encodePath)
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		if err == nil &amp;&amp; p == u.Path {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			return u.RawPath
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		}
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	}
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	if u.Path == &#34;*&#34; {
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>		return &#34;*&#34; <span class="comment">// don&#39;t escape (Issue 11202)</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	}
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	return escape(u.Path, encodePath)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span><span class="comment">// validEncoded reports whether s is a valid encoded path or fragment,</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span><span class="comment">// according to mode.</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span><span class="comment">// It must not contain any bytes that require escaping during encoding.</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>func validEncoded(s string, mode encoding) bool {
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i++ {
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		<span class="comment">// RFC 3986, Appendix A.</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		<span class="comment">// pchar = unreserved / pct-encoded / sub-delims / &#34;:&#34; / &#34;@&#34;.</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		<span class="comment">// shouldEscape is not quite compliant with the RFC,</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		<span class="comment">// so we check the sub-delims ourselves and let</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		<span class="comment">// shouldEscape handle the others.</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		switch s[i] {
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>		case &#39;!&#39;, &#39;$&#39;, &#39;&amp;&#39;, &#39;\&#39;&#39;, &#39;(&#39;, &#39;)&#39;, &#39;*&#39;, &#39;+&#39;, &#39;,&#39;, &#39;;&#39;, &#39;=&#39;, &#39;:&#39;, &#39;@&#39;:
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>			<span class="comment">// ok</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		case &#39;[&#39;, &#39;]&#39;:
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>			<span class="comment">// ok - not specified in RFC 3986 but left alone by modern browsers</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		case &#39;%&#39;:
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>			<span class="comment">// ok - percent encoded, will decode</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		default:
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>			if shouldEscape(s[i], mode) {
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>				return false
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			}
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		}
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	return true
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>}
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span><span class="comment">// setFragment is like setPath but for Fragment/RawFragment.</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>func (u *URL) setFragment(f string) error {
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	frag, err := unescape(f, encodeFragment)
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	if err != nil {
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		return err
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	}
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	u.Fragment = frag
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	if escf := escape(frag, encodeFragment); f == escf {
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		<span class="comment">// Default encoding is fine.</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>		u.RawFragment = &#34;&#34;
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	} else {
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		u.RawFragment = f
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	}
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	return nil
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>}
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span><span class="comment">// EscapedFragment returns the escaped form of u.Fragment.</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span><span class="comment">// In general there are multiple possible escaped forms of any fragment.</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span><span class="comment">// EscapedFragment returns u.RawFragment when it is a valid escaping of u.Fragment.</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span><span class="comment">// Otherwise EscapedFragment ignores u.RawFragment and computes an escaped</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span><span class="comment">// form on its own.</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span><span class="comment">// The [URL.String] method uses EscapedFragment to construct its result.</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span><span class="comment">// In general, code should call EscapedFragment instead of</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span><span class="comment">// reading u.RawFragment directly.</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>func (u *URL) EscapedFragment() string {
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	if u.RawFragment != &#34;&#34; &amp;&amp; validEncoded(u.RawFragment, encodeFragment) {
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		f, err := unescape(u.RawFragment, encodeFragment)
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		if err == nil &amp;&amp; f == u.Fragment {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			return u.RawFragment
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>		}
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	}
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	return escape(u.Fragment, encodeFragment)
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span><span class="comment">// validOptionalPort reports whether port is either an empty string</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span><span class="comment">// or matches /^:\d*$/</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>func validOptionalPort(port string) bool {
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	if port == &#34;&#34; {
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>		return true
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	if port[0] != &#39;:&#39; {
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>		return false
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	}
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	for _, b := range port[1:] {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		if b &lt; &#39;0&#39; || b &gt; &#39;9&#39; {
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>			return false
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>		}
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	}
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	return true
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>}
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span><span class="comment">// String reassembles the [URL] into a valid URL string.</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span><span class="comment">// The general form of the result is one of:</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span><span class="comment">//	scheme:opaque?query#fragment</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span><span class="comment">//	scheme://userinfo@host/path?query#fragment</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span><span class="comment">// If u.Opaque is non-empty, String uses the first form;</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span><span class="comment">// otherwise it uses the second form.</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span><span class="comment">// Any non-ASCII characters in host are escaped.</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span><span class="comment">// To obtain the path, String uses u.EscapedPath().</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span><span class="comment">// In the second form, the following rules apply:</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span><span class="comment">//   - if u.Scheme is empty, scheme: is omitted.</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span><span class="comment">//   - if u.User is nil, userinfo@ is omitted.</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span><span class="comment">//   - if u.Host is empty, host/ is omitted.</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span><span class="comment">//   - if u.Scheme and u.Host are empty and u.User is nil,</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span><span class="comment">//     the entire scheme://userinfo@host/ is omitted.</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span><span class="comment">//   - if u.Host is non-empty and u.Path begins with a /,</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span><span class="comment">//     the form host/path does not add its own /.</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span><span class="comment">//   - if u.RawQuery is empty, ?query is omitted.</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span><span class="comment">//   - if u.Fragment is empty, #fragment is omitted.</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>func (u *URL) String() string {
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	var buf strings.Builder
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	if u.Scheme != &#34;&#34; {
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		buf.WriteString(u.Scheme)
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		buf.WriteByte(&#39;:&#39;)
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	}
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	if u.Opaque != &#34;&#34; {
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>		buf.WriteString(u.Opaque)
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	} else {
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		if u.Scheme != &#34;&#34; || u.Host != &#34;&#34; || u.User != nil {
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>			if u.OmitHost &amp;&amp; u.Host == &#34;&#34; &amp;&amp; u.User == nil {
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>				<span class="comment">// omit empty host</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>			} else {
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>				if u.Host != &#34;&#34; || u.Path != &#34;&#34; || u.User != nil {
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>					buf.WriteString(&#34;//&#34;)
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>				}
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>				if ui := u.User; ui != nil {
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>					buf.WriteString(ui.String())
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>					buf.WriteByte(&#39;@&#39;)
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>				}
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>				if h := u.Host; h != &#34;&#34; {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>					buf.WriteString(escape(h, encodeHost))
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>				}
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>			}
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		}
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		path := u.EscapedPath()
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		if path != &#34;&#34; &amp;&amp; path[0] != &#39;/&#39; &amp;&amp; u.Host != &#34;&#34; {
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>			buf.WriteByte(&#39;/&#39;)
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>		}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>		if buf.Len() == 0 {
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			<span class="comment">// RFC 3986 §4.2</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>			<span class="comment">// A path segment that contains a colon character (e.g., &#34;this:that&#34;)</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>			<span class="comment">// cannot be used as the first segment of a relative-path reference, as</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>			<span class="comment">// it would be mistaken for a scheme name. Such a segment must be</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>			<span class="comment">// preceded by a dot-segment (e.g., &#34;./this:that&#34;) to make a relative-</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>			<span class="comment">// path reference.</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>			if segment, _, _ := strings.Cut(path, &#34;/&#34;); strings.Contains(segment, &#34;:&#34;) {
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>				buf.WriteString(&#34;./&#34;)
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>			}
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		}
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>		buf.WriteString(path)
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	}
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	if u.ForceQuery || u.RawQuery != &#34;&#34; {
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		buf.WriteByte(&#39;?&#39;)
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		buf.WriteString(u.RawQuery)
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	}
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	if u.Fragment != &#34;&#34; {
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		buf.WriteByte(&#39;#&#39;)
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		buf.WriteString(u.EscapedFragment())
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	}
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>	return buf.String()
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>}
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span><span class="comment">// Redacted is like [URL.String] but replaces any password with &#34;xxxxx&#34;.</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span><span class="comment">// Only the password in u.User is redacted.</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>func (u *URL) Redacted() string {
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	if u == nil {
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	}
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	ru := *u
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	if _, has := ru.User.Password(); has {
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>		ru.User = UserPassword(ru.User.Username(), &#34;xxxxx&#34;)
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	}
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	return ru.String()
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>}
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span><span class="comment">// Values maps a string key to a list of values.</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span><span class="comment">// It is typically used for query parameters and form values.</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span><span class="comment">// Unlike in the http.Header map, the keys in a Values map</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span><span class="comment">// are case-sensitive.</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>type Values map[string][]string
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span><span class="comment">// Get gets the first value associated with the given key.</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span><span class="comment">// If there are no values associated with the key, Get returns</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span><span class="comment">// the empty string. To access multiple values, use the map</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span><span class="comment">// directly.</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>func (v Values) Get(key string) string {
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	vs := v[key]
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	if len(vs) == 0 {
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	}
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	return vs[0]
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>}
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span><span class="comment">// Set sets the key to value. It replaces any existing</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span><span class="comment">// values.</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>func (v Values) Set(key, value string) {
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>	v[key] = []string{value}
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>}
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span><span class="comment">// Add adds the value to key. It appends to any existing</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span><span class="comment">// values associated with key.</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>func (v Values) Add(key, value string) {
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	v[key] = append(v[key], value)
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>}
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span><span class="comment">// Del deletes the values associated with key.</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>func (v Values) Del(key string) {
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	delete(v, key)
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>}
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>
<span id="L917" class="ln">   917&nbsp;&nbsp;</span><span class="comment">// Has checks whether a given key is set.</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>func (v Values) Has(key string) bool {
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>	_, ok := v[key]
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>	return ok
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>}
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span><span class="comment">// ParseQuery parses the URL-encoded query string and returns</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span><span class="comment">// a map listing the values specified for each key.</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span><span class="comment">// ParseQuery always returns a non-nil map containing all the</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span><span class="comment">// valid query parameters found; err describes the first decoding error</span>
<span id="L927" class="ln">   927&nbsp;&nbsp;</span><span class="comment">// encountered, if any.</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span><span class="comment">// Query is expected to be a list of key=value settings separated by ampersands.</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span><span class="comment">// A setting without an equals sign is interpreted as a key set to an empty</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span><span class="comment">// value.</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span><span class="comment">// Settings containing a non-URL-encoded semicolon are considered invalid.</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>func ParseQuery(query string) (Values, error) {
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	m := make(Values)
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>	err := parseQuery(m, query)
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>	return m, err
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>}
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>func parseQuery(m Values, query string) (err error) {
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	for query != &#34;&#34; {
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		var key string
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		key, query, _ = strings.Cut(query, &#34;&amp;&#34;)
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>		if strings.Contains(key, &#34;;&#34;) {
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>			err = fmt.Errorf(&#34;invalid semicolon separator in query&#34;)
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>			continue
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>		}
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>		if key == &#34;&#34; {
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>			continue
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		}
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>		key, value, _ := strings.Cut(key, &#34;=&#34;)
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>		key, err1 := QueryUnescape(key)
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>		if err1 != nil {
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>			if err == nil {
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>				err = err1
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>			}
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>			continue
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>		}
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>		value, err1 = QueryUnescape(value)
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		if err1 != nil {
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>			if err == nil {
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>				err = err1
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>			}
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>			continue
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>		}
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>		m[key] = append(m[key], value)
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>	}
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>	return err
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>}
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span><span class="comment">// Encode encodes the values into “URL encoded” form</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span><span class="comment">// (&#34;bar=baz&amp;foo=quux&#34;) sorted by key.</span>
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>func (v Values) Encode() string {
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	if len(v) == 0 {
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>	}
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	var buf strings.Builder
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>	keys := make([]string, 0, len(v))
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	for k := range v {
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>		keys = append(keys, k)
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	}
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	sort.Strings(keys)
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	for _, k := range keys {
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		vs := v[k]
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>		keyEscaped := QueryEscape(k)
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>		for _, v := range vs {
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>			if buf.Len() &gt; 0 {
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>				buf.WriteByte(&#39;&amp;&#39;)
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>			}
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>			buf.WriteString(keyEscaped)
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>			buf.WriteByte(&#39;=&#39;)
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>			buf.WriteString(QueryEscape(v))
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>		}
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	}
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	return buf.String()
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>}
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span><span class="comment">// resolvePath applies special path segments from refs and applies</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span><span class="comment">// them to base, per RFC 3986.</span>
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>func resolvePath(base, ref string) string {
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	var full string
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>	if ref == &#34;&#34; {
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>		full = base
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>	} else if ref[0] != &#39;/&#39; {
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>		i := strings.LastIndex(base, &#34;/&#34;)
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>		full = base[:i+1] + ref
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	} else {
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>		full = ref
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	}
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>	if full == &#34;&#34; {
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>	}
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	var (
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>		elem string
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>		dst  strings.Builder
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	)
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	first := true
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	remaining := full
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>	<span class="comment">// We want to return a leading &#39;/&#39;, so write it now.</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>	dst.WriteByte(&#39;/&#39;)
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	found := true
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	for found {
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		elem, remaining, found = strings.Cut(remaining, &#34;/&#34;)
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>		if elem == &#34;.&#34; {
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>			first = false
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>			<span class="comment">// drop</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>			continue
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>		}
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>		if elem == &#34;..&#34; {
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>			<span class="comment">// Ignore the leading &#39;/&#39; we already wrote.</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>			str := dst.String()[1:]
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>			index := strings.LastIndexByte(str, &#39;/&#39;)
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>			dst.Reset()
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>			dst.WriteByte(&#39;/&#39;)
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>			if index == -1 {
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>				first = true
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>			} else {
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>				dst.WriteString(str[:index])
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>			}
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		} else {
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>			if !first {
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>				dst.WriteByte(&#39;/&#39;)
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>			}
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>			dst.WriteString(elem)
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>			first = false
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>		}
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>	}
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>	if elem == &#34;.&#34; || elem == &#34;..&#34; {
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>		dst.WriteByte(&#39;/&#39;)
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>	}
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>	<span class="comment">// We wrote an initial &#39;/&#39;, but we don&#39;t want two.</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>	r := dst.String()
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>	if len(r) &gt; 1 &amp;&amp; r[1] == &#39;/&#39; {
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>		r = r[1:]
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>	}
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>	return r
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>}
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span><span class="comment">// IsAbs reports whether the [URL] is absolute.</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span><span class="comment">// Absolute means that it has a non-empty scheme.</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>func (u *URL) IsAbs() bool {
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	return u.Scheme != &#34;&#34;
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>}
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span><span class="comment">// Parse parses a [URL] in the context of the receiver. The provided URL</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span><span class="comment">// may be relative or absolute. Parse returns nil, err on parse</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span><span class="comment">// failure, otherwise its return value is the same as [URL.ResolveReference].</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>func (u *URL) Parse(ref string) (*URL, error) {
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	refURL, err := Parse(ref)
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	if err != nil {
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>		return nil, err
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>	}
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	return u.ResolveReference(refURL), nil
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>}
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span><span class="comment">// ResolveReference resolves a URI reference to an absolute URI from</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span><span class="comment">// an absolute base URI u, per RFC 3986 Section 5.2. The URI reference</span>
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span><span class="comment">// may be relative or absolute. ResolveReference always returns a new</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span><span class="comment">// [URL] instance, even if the returned URL is identical to either the</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span><span class="comment">// base or reference. If ref is an absolute URL, then ResolveReference</span>
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span><span class="comment">// ignores base and returns a copy of ref.</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>func (u *URL) ResolveReference(ref *URL) *URL {
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	url := *ref
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>	if ref.Scheme == &#34;&#34; {
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>		url.Scheme = u.Scheme
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>	}
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>	if ref.Scheme != &#34;&#34; || ref.Host != &#34;&#34; || ref.User != nil {
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>		<span class="comment">// The &#34;absoluteURI&#34; or &#34;net_path&#34; cases.</span>
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>		<span class="comment">// We can ignore the error from setPath since we know we provided a</span>
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>		<span class="comment">// validly-escaped path.</span>
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>		url.setPath(resolvePath(ref.EscapedPath(), &#34;&#34;))
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>		return &amp;url
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>	}
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>	if ref.Opaque != &#34;&#34; {
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>		url.User = nil
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>		url.Host = &#34;&#34;
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>		url.Path = &#34;&#34;
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>		return &amp;url
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>	}
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>	if ref.Path == &#34;&#34; &amp;&amp; !ref.ForceQuery &amp;&amp; ref.RawQuery == &#34;&#34; {
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>		url.RawQuery = u.RawQuery
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>		if ref.Fragment == &#34;&#34; {
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>			url.Fragment = u.Fragment
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>			url.RawFragment = u.RawFragment
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>		}
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>	}
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>	<span class="comment">// The &#34;abs_path&#34; or &#34;rel_path&#34; cases.</span>
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>	url.Host = u.Host
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>	url.User = u.User
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>	url.setPath(resolvePath(u.EscapedPath(), ref.EscapedPath()))
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>	return &amp;url
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>}
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span><span class="comment">// Query parses RawQuery and returns the corresponding values.</span>
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span><span class="comment">// It silently discards malformed value pairs.</span>
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span><span class="comment">// To check errors use [ParseQuery].</span>
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>func (u *URL) Query() Values {
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	v, _ := ParseQuery(u.RawQuery)
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	return v
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>}
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span><span class="comment">// RequestURI returns the encoded path?query or opaque?query</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span><span class="comment">// string that would be used in an HTTP request for u.</span>
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>func (u *URL) RequestURI() string {
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>	result := u.Opaque
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>	if result == &#34;&#34; {
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>		result = u.EscapedPath()
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>		if result == &#34;&#34; {
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>			result = &#34;/&#34;
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>		}
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>	} else {
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>		if strings.HasPrefix(result, &#34;//&#34;) {
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>			result = u.Scheme + &#34;:&#34; + result
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>		}
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>	}
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>	if u.ForceQuery || u.RawQuery != &#34;&#34; {
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>		result += &#34;?&#34; + u.RawQuery
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	}
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	return result
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>}
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span><span class="comment">// Hostname returns u.Host, stripping any valid port number if present.</span>
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span><span class="comment">// If the result is enclosed in square brackets, as literal IPv6 addresses are,</span>
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span><span class="comment">// the square brackets are removed from the result.</span>
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>func (u *URL) Hostname() string {
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>	host, _ := splitHostPort(u.Host)
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>	return host
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>}
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span><span class="comment">// Port returns the port part of u.Host, without the leading colon.</span>
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span><span class="comment">// If u.Host doesn&#39;t contain a valid numeric port, Port returns an empty string.</span>
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>func (u *URL) Port() string {
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>	_, port := splitHostPort(u.Host)
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>	return port
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>}
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span><span class="comment">// splitHostPort separates host and port. If the port is not valid, it returns</span>
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span><span class="comment">// the entire input as host, and it doesn&#39;t check the validity of the host.</span>
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span><span class="comment">// Unlike net.SplitHostPort, but per RFC 3986, it requires ports to be numeric.</span>
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>func splitHostPort(hostPort string) (host, port string) {
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>	host = hostPort
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>	colon := strings.LastIndexByte(host, &#39;:&#39;)
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>	if colon != -1 &amp;&amp; validOptionalPort(host[colon:]) {
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>		host, port = host[:colon], host[colon+1:]
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>	}
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>	if strings.HasPrefix(host, &#34;[&#34;) &amp;&amp; strings.HasSuffix(host, &#34;]&#34;) {
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>		host = host[1 : len(host)-1]
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>	}
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>	return
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>}
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span><span class="comment">// Marshaling interface implementations.</span>
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span><span class="comment">// Would like to implement MarshalText/UnmarshalText but that will change the JSON representation of URLs.</span>
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>func (u *URL) MarshalBinary() (text []byte, err error) {
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>	return []byte(u.String()), nil
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>}
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>func (u *URL) UnmarshalBinary(text []byte) error {
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>	u1, err := Parse(string(text))
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>	if err != nil {
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>		return err
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>	}
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>	*u = *u1
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>	return nil
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>}
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span><span class="comment">// JoinPath returns a new [URL] with the provided path elements joined to</span>
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span><span class="comment">// any existing path and the resulting path cleaned of any ./ or ../ elements.</span>
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span><span class="comment">// Any sequences of multiple / characters will be reduced to a single /.</span>
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>func (u *URL) JoinPath(elem ...string) *URL {
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>	elem = append([]string{u.EscapedPath()}, elem...)
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>	var p string
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>	if !strings.HasPrefix(elem[0], &#34;/&#34;) {
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>		<span class="comment">// Return a relative path if u is relative,</span>
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>		<span class="comment">// but ensure that it contains no ../ elements.</span>
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>		elem[0] = &#34;/&#34; + elem[0]
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>		p = path.Join(elem...)[1:]
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>	} else {
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>		p = path.Join(elem...)
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>	}
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>	<span class="comment">// path.Join will remove any trailing slashes.</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>	<span class="comment">// Preserve at least one.</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>	if strings.HasSuffix(elem[len(elem)-1], &#34;/&#34;) &amp;&amp; !strings.HasSuffix(p, &#34;/&#34;) {
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>		p += &#34;/&#34;
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>	}
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>	url := *u
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>	url.setPath(p)
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>	return &amp;url
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>}
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span><span class="comment">// validUserinfo reports whether s is a valid userinfo string per RFC 3986</span>
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span><span class="comment">// Section 3.2.1:</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span><span class="comment">//	userinfo    = *( unreserved / pct-encoded / sub-delims / &#34;:&#34; )</span>
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span><span class="comment">//	unreserved  = ALPHA / DIGIT / &#34;-&#34; / &#34;.&#34; / &#34;_&#34; / &#34;~&#34;</span>
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span><span class="comment">//	sub-delims  = &#34;!&#34; / &#34;$&#34; / &#34;&amp;&#34; / &#34;&#39;&#34; / &#34;(&#34; / &#34;)&#34;</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span><span class="comment">//	              / &#34;*&#34; / &#34;+&#34; / &#34;,&#34; / &#34;;&#34; / &#34;=&#34;</span>
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span><span class="comment">// It doesn&#39;t validate pct-encoded. The caller does that via func unescape.</span>
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>func validUserinfo(s string) bool {
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	for _, r := range s {
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>		if &#39;A&#39; &lt;= r &amp;&amp; r &lt;= &#39;Z&#39; {
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>			continue
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>		}
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>		if &#39;a&#39; &lt;= r &amp;&amp; r &lt;= &#39;z&#39; {
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>			continue
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>		}
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>		if &#39;0&#39; &lt;= r &amp;&amp; r &lt;= &#39;9&#39; {
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>			continue
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>		}
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>		switch r {
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>		case &#39;-&#39;, &#39;.&#39;, &#39;_&#39;, &#39;:&#39;, &#39;~&#39;, &#39;!&#39;, &#39;$&#39;, &#39;&amp;&#39;, &#39;\&#39;&#39;,
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>			&#39;(&#39;, &#39;)&#39;, &#39;*&#39;, &#39;+&#39;, &#39;,&#39;, &#39;;&#39;, &#39;=&#39;, &#39;%&#39;, &#39;@&#39;:
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>			continue
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>		default:
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>			return false
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>		}
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>	}
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>	return true
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>}
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span><span class="comment">// stringContainsCTLByte reports whether s contains any ASCII control character.</span>
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>func stringContainsCTLByte(s string) bool {
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i++ {
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>		b := s[i]
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>		if b &lt; &#39; &#39; || b == 0x7f {
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>			return true
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>		}
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>	}
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>	return false
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>}
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span><span class="comment">// JoinPath returns a [URL] string with the provided path elements joined to</span>
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span><span class="comment">// the existing path of base and the resulting path cleaned of any ./ or ../ elements.</span>
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>func JoinPath(base string, elem ...string) (result string, err error) {
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>	url, err := Parse(base)
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>	if err != nil {
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>		return
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>	}
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	result = url.JoinPath(elem...).String()
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	return
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>}
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>
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
