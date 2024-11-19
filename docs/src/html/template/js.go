<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/html/template/js.go - Go Documentation Server</title>

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
<a href="js.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/html">html</a>/<a href="http://localhost:8080/src/html/template">template</a>/<span class="text-muted">js.go</span>
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
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;encoding/json&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// jsWhitespace contains all of the JS whitespace characters, as defined</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// by the \s character class.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// See https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_expressions/Character_classes.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>const jsWhitespace = &#34;\f\n\r\t\v\u0020\u00a0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u2028\u2029\u202f\u205f\u3000\ufeff&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// nextJSCtx returns the context that determines whether a slash after the</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// given run of tokens starts a regular expression instead of a division</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// operator: / or /=.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// This assumes that the token run does not include any string tokens, comment</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// tokens, regular expression literal tokens, or division operators.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// This fails on some valid but nonsensical JavaScript programs like</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// &#34;x = ++/foo/i&#34; which is quite different than &#34;x++/foo/i&#34;, but is not known to</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// fail on any known useful programs. It is based on the draft</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// JavaScript 2.0 lexical grammar and requires one token of lookbehind:</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// https://www.mozilla.org/js/language/js20-2000-07/rationale/syntax.html</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>func nextJSCtx(s []byte, preceding jsCtx) jsCtx {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// Trim all JS whitespace characters</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	s = bytes.TrimRight(s, jsWhitespace)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	if len(s) == 0 {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		return preceding
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// All cases below are in the single-byte UTF-8 group.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	switch c, n := s[len(s)-1], len(s); c {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	case &#39;+&#39;, &#39;-&#39;:
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		<span class="comment">// ++ and -- are not regexp preceders, but + and - are whether</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		<span class="comment">// they are used as infix or prefix operators.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		start := n - 1
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		<span class="comment">// Count the number of adjacent dashes or pluses.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		for start &gt; 0 &amp;&amp; s[start-1] == c {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			start--
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		if (n-start)&amp;1 == 1 {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			<span class="comment">// Reached for trailing minus signs since &#34;---&#34; is the</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			<span class="comment">// same as &#34;-- -&#34;.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			return jsCtxRegexp
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		return jsCtxDivOp
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	case &#39;.&#39;:
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		<span class="comment">// Handle &#34;42.&#34;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		if n != 1 &amp;&amp; &#39;0&#39; &lt;= s[n-2] &amp;&amp; s[n-2] &lt;= &#39;9&#39; {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			return jsCtxDivOp
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return jsCtxRegexp
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// Suffixes for all punctuators from section 7.7 of the language spec</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// that only end binary operators not handled above.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	case &#39;,&#39;, &#39;&lt;&#39;, &#39;&gt;&#39;, &#39;=&#39;, &#39;*&#39;, &#39;%&#39;, &#39;&amp;&#39;, &#39;|&#39;, &#39;^&#39;, &#39;?&#39;:
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		return jsCtxRegexp
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// Suffixes for all punctuators from section 7.7 of the language spec</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// that are prefix operators not handled above.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	case &#39;!&#39;, &#39;~&#39;:
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		return jsCtxRegexp
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// Matches all the punctuators from section 7.7 of the language spec</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// that are open brackets not handled above.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	case &#39;(&#39;, &#39;[&#39;:
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		return jsCtxRegexp
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// Matches all the punctuators from section 7.7 of the language spec</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// that precede expression starts.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	case &#39;:&#39;, &#39;;&#39;, &#39;{&#39;:
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		return jsCtxRegexp
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// CAVEAT: the close punctuators (&#39;}&#39;, &#39;]&#39;, &#39;)&#39;) precede div ops and</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// are handled in the default except for &#39;}&#39; which can precede a</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// division op as in</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">//    ({ valueOf: function () { return 42 } } / 2</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// which is valid, but, in practice, developers don&#39;t divide object</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// literals, so our heuristic works well for code like</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">//    function () { ... }  /foo/.test(x) &amp;&amp; sideEffect();</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// The &#39;)&#39; punctuator can precede a regular expression as in</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">//     if (b) /foo/.test(x) &amp;&amp; ...</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// but this is much less likely than</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">//     (a + b) / c</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	case &#39;}&#39;:
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		return jsCtxRegexp
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	default:
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		<span class="comment">// Look for an IdentifierName and see if it is a keyword that</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		<span class="comment">// can precede a regular expression.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		j := n
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		for j &gt; 0 &amp;&amp; isJSIdentPart(rune(s[j-1])) {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			j--
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		if regexpPrecederKeywords[string(s[j:])] {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			return jsCtxRegexp
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// Otherwise is a punctuator not listed above, or</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// a string which precedes a div op, or an identifier</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// which precedes a div op.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	return jsCtxDivOp
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// regexpPrecederKeywords is a set of reserved JS keywords that can precede a</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// regular expression in JS source.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>var regexpPrecederKeywords = map[string]bool{
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	&#34;break&#34;:      true,
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	&#34;case&#34;:       true,
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	&#34;continue&#34;:   true,
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	&#34;delete&#34;:     true,
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	&#34;do&#34;:         true,
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	&#34;else&#34;:       true,
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	&#34;finally&#34;:    true,
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	&#34;in&#34;:         true,
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	&#34;instanceof&#34;: true,
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	&#34;return&#34;:     true,
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	&#34;throw&#34;:      true,
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	&#34;try&#34;:        true,
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	&#34;typeof&#34;:     true,
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	&#34;void&#34;:       true,
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>var jsonMarshalType = reflect.TypeFor[json.Marshaler]()
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// indirectToJSONMarshaler returns the value, after dereferencing as many times</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// as necessary to reach the base type (or nil) or an implementation of json.Marshal.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func indirectToJSONMarshaler(a any) any {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// text/template now supports passing untyped nil as a func call</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// argument, so we must support it. Otherwise we&#39;d panic below, as one</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// cannot call the Type or Interface methods on an invalid</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// reflect.Value. See golang.org/issue/18716.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	if a == nil {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		return nil
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	v := reflect.ValueOf(a)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	for !v.Type().Implements(jsonMarshalType) &amp;&amp; v.Kind() == reflect.Pointer &amp;&amp; !v.IsNil() {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		v = v.Elem()
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	return v.Interface()
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// jsValEscaper escapes its inputs to a JS Expression (section 11.14) that has</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// neither side-effects nor free variables outside (NaN, Infinity).</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>func jsValEscaper(args ...any) string {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	var a any
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	if len(args) == 1 {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		a = indirectToJSONMarshaler(args[0])
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		switch t := a.(type) {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		case JS:
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			return string(t)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		case JSStr:
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			<span class="comment">// TODO: normalize quotes.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			return `&#34;` + string(t) + `&#34;`
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		case json.Marshaler:
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			<span class="comment">// Do not treat as a Stringer.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		case fmt.Stringer:
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			a = t.String()
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	} else {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		for i, arg := range args {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			args[i] = indirectToJSONMarshaler(arg)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		a = fmt.Sprint(args...)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// TODO: detect cycles before calling Marshal which loops infinitely on</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// cyclic data. This may be an unacceptable DoS risk.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	b, err := json.Marshal(a)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	if err != nil {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		<span class="comment">// While the standard JSON marshaller does not include user controlled</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		<span class="comment">// information in the error message, if a type has a MarshalJSON method,</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		<span class="comment">// the content of the error message is not guaranteed. Since we insert</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		<span class="comment">// the error into the template, as part of a comment, we attempt to</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		<span class="comment">// prevent the error from either terminating the comment, or the script</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		<span class="comment">// block itself.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		<span class="comment">// In particular we:</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		<span class="comment">//   * replace &#34;*/&#34; comment end tokens with &#34;* /&#34;, which does not</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		<span class="comment">//     terminate the comment</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		<span class="comment">//   * replace &#34;&lt;/script&#34; with &#34;\x3C/script&#34;, and &#34;&lt;!--&#34; with</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		<span class="comment">//     &#34;\x3C!--&#34;, which prevents confusing script block termination</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		<span class="comment">//     semantics</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		<span class="comment">// We also put a space before the comment so that if it is flush against</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		<span class="comment">// a division operator it is not turned into a line comment:</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		<span class="comment">//     x/{{y}}</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		<span class="comment">// turning into</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		<span class="comment">//     x//* error marshaling y:</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		<span class="comment">//          second line of error message */null</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		errStr := err.Error()
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		errStr = strings.ReplaceAll(errStr, &#34;*/&#34;, &#34;* /&#34;)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		errStr = strings.ReplaceAll(errStr, &#34;&lt;/script&#34;, `\x3C/script`)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		errStr = strings.ReplaceAll(errStr, &#34;&lt;!--&#34;, `\x3C!--`)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		return fmt.Sprintf(&#34; /* %s */null &#34;, errStr)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// TODO: maybe post-process output to prevent it from containing</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// &#34;&lt;!--&#34;, &#34;--&gt;&#34;, &#34;&lt;![CDATA[&#34;, &#34;]]&gt;&#34;, or &#34;&lt;/script&#34;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	<span class="comment">// in case custom marshalers produce output containing those.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// Note: Do not use \x escaping to save bytes because it is not JSON compatible and this escaper</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// supports ld+json content-type.</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	if len(b) == 0 {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		<span class="comment">// In, `x=y/{{.}}*z` a json.Marshaler that produces &#34;&#34; should</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		<span class="comment">// not cause the output `x=y/*z`.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		return &#34; null &#34;
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	first, _ := utf8.DecodeRune(b)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	last, _ := utf8.DecodeLastRune(b)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	var buf strings.Builder
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// Prevent IdentifierNames and NumericLiterals from running into</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// keywords: in, instanceof, typeof, void</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	pad := isJSIdentPart(first) || isJSIdentPart(last)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	if pad {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		buf.WriteByte(&#39; &#39;)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	written := 0
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// Make sure that json.Marshal escapes codepoints U+2028 &amp; U+2029</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">// so it falls within the subset of JSON which is valid JS.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	for i := 0; i &lt; len(b); {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		rune, n := utf8.DecodeRune(b[i:])
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		repl := &#34;&#34;
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		if rune == 0x2028 {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			repl = `\u2028`
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		} else if rune == 0x2029 {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			repl = `\u2029`
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		if repl != &#34;&#34; {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			buf.Write(b[written:i])
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			buf.WriteString(repl)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			written = i + n
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		i += n
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	if buf.Len() != 0 {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		buf.Write(b[written:])
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		if pad {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			buf.WriteByte(&#39; &#39;)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		return buf.String()
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	return string(b)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// jsStrEscaper produces a string that can be included between quotes in</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">// JavaScript source, in JavaScript embedded in an HTML5 &lt;script&gt; element,</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">// or in an HTML5 event handler attribute such as onclick.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func jsStrEscaper(args ...any) string {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	s, t := stringify(args...)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	if t == contentTypeJSStr {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		return replace(s, jsStrNormReplacementTable)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	return replace(s, jsStrReplacementTable)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>func jsTmplLitEscaper(args ...any) string {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	s, _ := stringify(args...)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	return replace(s, jsBqStrReplacementTable)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span><span class="comment">// jsRegexpEscaper behaves like jsStrEscaper but escapes regular expression</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">// specials so the result is treated literally when included in a regular</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">// expression literal. /foo{{.X}}bar/ matches the string &#34;foo&#34; followed by</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span><span class="comment">// the literal text of {{.X}} followed by the string &#34;bar&#34;.</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>func jsRegexpEscaper(args ...any) string {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	s, _ := stringify(args...)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	s = replace(s, jsRegexpReplacementTable)
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	if s == &#34;&#34; {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		<span class="comment">// /{{.X}}/ should not produce a line comment when .X == &#34;&#34;.</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		return &#34;(?:)&#34;
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	return s
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span><span class="comment">// replace replaces each rune r of s with replacementTable[r], provided that</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">// r &lt; len(replacementTable). If replacementTable[r] is the empty string then</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">// no replacement is made.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// It also replaces runes U+2028 and U+2029 with the raw strings `\u2028` and</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">// `\u2029`.</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>func replace(s string, replacementTable []string) string {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	var b strings.Builder
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	r, w, written := rune(0), 0, 0
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i += w {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		<span class="comment">// See comment in htmlEscaper.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		r, w = utf8.DecodeRuneInString(s[i:])
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		var repl string
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		switch {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		case int(r) &lt; len(lowUnicodeReplacementTable):
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			repl = lowUnicodeReplacementTable[r]
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		case int(r) &lt; len(replacementTable) &amp;&amp; replacementTable[r] != &#34;&#34;:
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			repl = replacementTable[r]
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		case r == &#39;\u2028&#39;:
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			repl = `\u2028`
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		case r == &#39;\u2029&#39;:
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			repl = `\u2029`
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		default:
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			continue
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		if written == 0 {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			b.Grow(len(s))
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		b.WriteString(s[written:i])
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		b.WriteString(repl)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		written = i + w
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	if written == 0 {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		return s
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	}
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	b.WriteString(s[written:])
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	return b.String()
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>var lowUnicodeReplacementTable = []string{
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	0: `\u0000`, 1: `\u0001`, 2: `\u0002`, 3: `\u0003`, 4: `\u0004`, 5: `\u0005`, 6: `\u0006`,
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	&#39;\a&#39;: `\u0007`,
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	&#39;\b&#39;: `\u0008`,
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	&#39;\t&#39;: `\t`,
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	&#39;\n&#39;: `\n`,
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	&#39;\v&#39;: `\u000b`, <span class="comment">// &#34;\v&#34; == &#34;v&#34; on IE 6.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	&#39;\f&#39;: `\f`,
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	&#39;\r&#39;: `\r`,
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	0xe:  `\u000e`, 0xf: `\u000f`, 0x10: `\u0010`, 0x11: `\u0011`, 0x12: `\u0012`, 0x13: `\u0013`,
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	0x14: `\u0014`, 0x15: `\u0015`, 0x16: `\u0016`, 0x17: `\u0017`, 0x18: `\u0018`, 0x19: `\u0019`,
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	0x1a: `\u001a`, 0x1b: `\u001b`, 0x1c: `\u001c`, 0x1d: `\u001d`, 0x1e: `\u001e`, 0x1f: `\u001f`,
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>var jsStrReplacementTable = []string{
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	0:    `\u0000`,
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	&#39;\t&#39;: `\t`,
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	&#39;\n&#39;: `\n`,
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	&#39;\v&#39;: `\u000b`, <span class="comment">// &#34;\v&#34; == &#34;v&#34; on IE 6.</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	&#39;\f&#39;: `\f`,
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	&#39;\r&#39;: `\r`,
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// Encode HTML specials as hex so the output can be embedded</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// in HTML attributes without further encoding.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	&#39;&#34;&#39;:  `\u0022`,
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	&#39;`&#39;:  `\u0060`,
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	&#39;&amp;&#39;:  `\u0026`,
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	&#39;\&#39;&#39;: `\u0027`,
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	&#39;+&#39;:  `\u002b`,
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	&#39;/&#39;:  `\/`,
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	&#39;&lt;&#39;:  `\u003c`,
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	&#39;&gt;&#39;:  `\u003e`,
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	&#39;\\&#39;: `\\`,
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// jsBqStrReplacementTable is like jsStrReplacementTable except it also contains</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// the special characters for JS template literals: $, {, and }.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>var jsBqStrReplacementTable = []string{
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	0:    `\u0000`,
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	&#39;\t&#39;: `\t`,
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	&#39;\n&#39;: `\n`,
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	&#39;\v&#39;: `\u000b`, <span class="comment">// &#34;\v&#34; == &#34;v&#34; on IE 6.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	&#39;\f&#39;: `\f`,
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	&#39;\r&#39;: `\r`,
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	<span class="comment">// Encode HTML specials as hex so the output can be embedded</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	<span class="comment">// in HTML attributes without further encoding.</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	&#39;&#34;&#39;:  `\u0022`,
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	&#39;`&#39;:  `\u0060`,
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	&#39;&amp;&#39;:  `\u0026`,
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	&#39;\&#39;&#39;: `\u0027`,
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	&#39;+&#39;:  `\u002b`,
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	&#39;/&#39;:  `\/`,
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	&#39;&lt;&#39;:  `\u003c`,
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	&#39;&gt;&#39;:  `\u003e`,
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	&#39;\\&#39;: `\\`,
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	&#39;$&#39;:  `\u0024`,
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	&#39;{&#39;:  `\u007b`,
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	&#39;}&#39;:  `\u007d`,
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span><span class="comment">// jsStrNormReplacementTable is like jsStrReplacementTable but does not</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span><span class="comment">// overencode existing escapes since this table has no entry for `\`.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>var jsStrNormReplacementTable = []string{
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	0:    `\u0000`,
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	&#39;\t&#39;: `\t`,
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	&#39;\n&#39;: `\n`,
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	&#39;\v&#39;: `\u000b`, <span class="comment">// &#34;\v&#34; == &#34;v&#34; on IE 6.</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	&#39;\f&#39;: `\f`,
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	&#39;\r&#39;: `\r`,
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// Encode HTML specials as hex so the output can be embedded</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	<span class="comment">// in HTML attributes without further encoding.</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	&#39;&#34;&#39;:  `\u0022`,
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	&#39;&amp;&#39;:  `\u0026`,
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	&#39;\&#39;&#39;: `\u0027`,
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	&#39;`&#39;:  `\u0060`,
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	&#39;+&#39;:  `\u002b`,
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	&#39;/&#39;:  `\/`,
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	&#39;&lt;&#39;:  `\u003c`,
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	&#39;&gt;&#39;:  `\u003e`,
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>var jsRegexpReplacementTable = []string{
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	0:    `\u0000`,
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	&#39;\t&#39;: `\t`,
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	&#39;\n&#39;: `\n`,
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	&#39;\v&#39;: `\u000b`, <span class="comment">// &#34;\v&#34; == &#34;v&#34; on IE 6.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	&#39;\f&#39;: `\f`,
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	&#39;\r&#39;: `\r`,
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	<span class="comment">// Encode HTML specials as hex so the output can be embedded</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	<span class="comment">// in HTML attributes without further encoding.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	&#39;&#34;&#39;:  `\u0022`,
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	&#39;$&#39;:  `\$`,
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	&#39;&amp;&#39;:  `\u0026`,
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	&#39;\&#39;&#39;: `\u0027`,
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	&#39;(&#39;:  `\(`,
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	&#39;)&#39;:  `\)`,
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	&#39;*&#39;:  `\*`,
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	&#39;+&#39;:  `\u002b`,
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	&#39;-&#39;:  `\-`,
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	&#39;.&#39;:  `\.`,
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	&#39;/&#39;:  `\/`,
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	&#39;&lt;&#39;:  `\u003c`,
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	&#39;&gt;&#39;:  `\u003e`,
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	&#39;?&#39;:  `\?`,
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	&#39;[&#39;:  `\[`,
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	&#39;\\&#39;: `\\`,
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	&#39;]&#39;:  `\]`,
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	&#39;^&#39;:  `\^`,
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	&#39;{&#39;:  `\{`,
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	&#39;|&#39;:  `\|`,
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	&#39;}&#39;:  `\}`,
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span><span class="comment">// isJSIdentPart reports whether the given rune is a JS identifier part.</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span><span class="comment">// It does not handle all the non-Latin letters, joiners, and combining marks,</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span><span class="comment">// but it does handle every codepoint that can occur in a numeric literal or</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span><span class="comment">// a keyword.</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>func isJSIdentPart(r rune) bool {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	switch {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	case r == &#39;$&#39;:
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		return true
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	case &#39;0&#39; &lt;= r &amp;&amp; r &lt;= &#39;9&#39;:
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		return true
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	case &#39;A&#39; &lt;= r &amp;&amp; r &lt;= &#39;Z&#39;:
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		return true
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	case r == &#39;_&#39;:
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		return true
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	case &#39;a&#39; &lt;= r &amp;&amp; r &lt;= &#39;z&#39;:
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		return true
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	return false
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span><span class="comment">// isJSType reports whether the given MIME type should be considered JavaScript.</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span><span class="comment">// It is used to determine whether a script tag with a type attribute is a javascript container.</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>func isJSType(mimeType string) bool {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	<span class="comment">// per</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	<span class="comment">//   https://www.w3.org/TR/html5/scripting-1.html#attr-script-type</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	<span class="comment">//   https://tools.ietf.org/html/rfc7231#section-3.1.1</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	<span class="comment">//   https://tools.ietf.org/html/rfc4329#section-3</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	<span class="comment">//   https://www.ietf.org/rfc/rfc4627.txt</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	<span class="comment">// discard parameters</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	mimeType, _, _ = strings.Cut(mimeType, &#34;;&#34;)
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	mimeType = strings.ToLower(mimeType)
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	mimeType = strings.TrimSpace(mimeType)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	switch mimeType {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	case
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		&#34;application/ecmascript&#34;,
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		&#34;application/javascript&#34;,
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		&#34;application/json&#34;,
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		&#34;application/ld+json&#34;,
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		&#34;application/x-ecmascript&#34;,
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		&#34;application/x-javascript&#34;,
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		&#34;module&#34;,
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		&#34;text/ecmascript&#34;,
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		&#34;text/javascript&#34;,
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		&#34;text/javascript1.0&#34;,
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		&#34;text/javascript1.1&#34;,
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		&#34;text/javascript1.2&#34;,
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		&#34;text/javascript1.3&#34;,
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		&#34;text/javascript1.4&#34;,
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		&#34;text/javascript1.5&#34;,
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		&#34;text/jscript&#34;,
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		&#34;text/livescript&#34;,
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		&#34;text/x-ecmascript&#34;,
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		&#34;text/x-javascript&#34;:
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		return true
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	default:
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		return false
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>
</pre><p><a href="js.go?m=text">View as plain text</a></p>

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
