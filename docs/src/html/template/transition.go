<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/html/template/transition.go - Go Documentation Server</title>

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
<a href="transition.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/html">html</a>/<a href="http://localhost:8080/src/html/template">template</a>/<span class="text-muted">transition.go</span>
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
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// transitionFunc is the array of context transition functions for text nodes.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// A transition function takes a context and template text input, and returns</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// the updated context and the number of bytes consumed from the front of the</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// input.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>var transitionFunc = [...]func(context, []byte) (context, int){
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	stateText:           tText,
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	stateTag:            tTag,
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	stateAttrName:       tAttrName,
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	stateAfterName:      tAfterName,
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	stateBeforeValue:    tBeforeValue,
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	stateHTMLCmt:        tHTMLCmt,
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	stateRCDATA:         tSpecialTagEnd,
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	stateAttr:           tAttr,
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	stateURL:            tURL,
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	stateSrcset:         tURL,
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	stateJS:             tJS,
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	stateJSDqStr:        tJSDelimited,
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	stateJSSqStr:        tJSDelimited,
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	stateJSRegexp:       tJSDelimited,
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	stateJSTmplLit:      tJSTmpl,
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	stateJSBlockCmt:     tBlockCmt,
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	stateJSLineCmt:      tLineCmt,
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	stateJSHTMLOpenCmt:  tLineCmt,
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	stateJSHTMLCloseCmt: tLineCmt,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	stateCSS:            tCSS,
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	stateCSSDqStr:       tCSSStr,
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	stateCSSSqStr:       tCSSStr,
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	stateCSSDqURL:       tCSSStr,
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	stateCSSSqURL:       tCSSStr,
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	stateCSSURL:         tCSSStr,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	stateCSSBlockCmt:    tBlockCmt,
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	stateCSSLineCmt:     tLineCmt,
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	stateError:          tError,
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>var commentStart = []byte(&#34;&lt;!--&#34;)
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>var commentEnd = []byte(&#34;--&gt;&#34;)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// tText is the context transition function for the text state.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>func tText(c context, s []byte) (context, int) {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	k := 0
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	for {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		i := k + bytes.IndexByte(s[k:], &#39;&lt;&#39;)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		if i &lt; k || i+1 == len(s) {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			return c, len(s)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		} else if i+4 &lt;= len(s) &amp;&amp; bytes.Equal(commentStart, s[i:i+4]) {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			return context{state: stateHTMLCmt}, i + 4
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		i++
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		end := false
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		if s[i] == &#39;/&#39; {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			if i+1 == len(s) {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>				return c, len(s)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			end, i = true, i+1
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		j, e := eatTagName(s, i)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		if j != i {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			if end {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>				e = elementNone
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			<span class="comment">// We&#39;ve found an HTML tag.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			return context{state: stateTag, element: e}, j
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		k = j
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>var elementContentType = [...]state{
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	elementNone:     stateText,
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	elementScript:   stateJS,
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	elementStyle:    stateCSS,
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	elementTextarea: stateRCDATA,
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	elementTitle:    stateRCDATA,
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// tTag is the context transition function for the tag state.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>func tTag(c context, s []byte) (context, int) {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// Find the attribute name.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	i := eatWhiteSpace(s, 0)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	if i == len(s) {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		return c, len(s)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	if s[i] == &#39;&gt;&#39; {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		return context{
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			state:   elementContentType[c.element],
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			element: c.element,
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		}, i + 1
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	j, err := eatAttrName(s, i)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	if err != nil {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		return context{state: stateError, err: err}, len(s)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	state, attr := stateTag, attrNone
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	if i == j {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		return context{
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			state: stateError,
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			err:   errorf(ErrBadHTML, nil, 0, &#34;expected space, attr name, or end of tag, but got %q&#34;, s[i:]),
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		}, len(s)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	attrName := strings.ToLower(string(s[i:j]))
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	if c.element == elementScript &amp;&amp; attrName == &#34;type&#34; {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		attr = attrScriptType
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	} else {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		switch attrType(attrName) {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		case contentTypeURL:
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			attr = attrURL
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		case contentTypeCSS:
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			attr = attrStyle
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		case contentTypeJS:
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			attr = attrScript
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		case contentTypeSrcset:
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			attr = attrSrcset
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if j == len(s) {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		state = stateAttrName
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	} else {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		state = stateAfterName
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	return context{state: state, element: c.element, attr: attr}, j
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// tAttrName is the context transition function for stateAttrName.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>func tAttrName(c context, s []byte) (context, int) {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	i, err := eatAttrName(s, 0)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	if err != nil {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		return context{state: stateError, err: err}, len(s)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	} else if i != len(s) {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		c.state = stateAfterName
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	return c, i
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// tAfterName is the context transition function for stateAfterName.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>func tAfterName(c context, s []byte) (context, int) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// Look for the start of the value.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	i := eatWhiteSpace(s, 0)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	if i == len(s) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		return c, len(s)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	} else if s[i] != &#39;=&#39; {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		<span class="comment">// Occurs due to tag ending &#39;&gt;&#39;, and valueless attribute.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		c.state = stateTag
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return c, i
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	c.state = stateBeforeValue
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// Consume the &#34;=&#34;.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	return c, i + 1
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>var attrStartStates = [...]state{
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	attrNone:       stateAttr,
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	attrScript:     stateJS,
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	attrScriptType: stateAttr,
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	attrStyle:      stateCSS,
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	attrURL:        stateURL,
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	attrSrcset:     stateSrcset,
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// tBeforeValue is the context transition function for stateBeforeValue.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>func tBeforeValue(c context, s []byte) (context, int) {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	i := eatWhiteSpace(s, 0)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	if i == len(s) {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		return c, len(s)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// Find the attribute delimiter.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	delim := delimSpaceOrTagEnd
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	switch s[i] {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	case &#39;\&#39;&#39;:
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		delim, i = delimSingleQuote, i+1
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	case &#39;&#34;&#39;:
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		delim, i = delimDoubleQuote, i+1
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	c.state, c.delim = attrStartStates[c.attr], delim
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	return c, i
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// tHTMLCmt is the context transition function for stateHTMLCmt.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>func tHTMLCmt(c context, s []byte) (context, int) {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	if i := bytes.Index(s, commentEnd); i != -1 {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		return context{}, i + 3
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	return c, len(s)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// specialTagEndMarkers maps element types to the character sequence that</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// case-insensitively signals the end of the special tag body.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>var specialTagEndMarkers = [...][]byte{
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	elementScript:   []byte(&#34;script&#34;),
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	elementStyle:    []byte(&#34;style&#34;),
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	elementTextarea: []byte(&#34;textarea&#34;),
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	elementTitle:    []byte(&#34;title&#34;),
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>var (
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	specialTagEndPrefix = []byte(&#34;&lt;/&#34;)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	tagEndSeparators    = []byte(&#34;&gt; \t\n\f/&#34;)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// tSpecialTagEnd is the context transition function for raw text and RCDATA</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// element states.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>func tSpecialTagEnd(c context, s []byte) (context, int) {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	if c.element != elementNone {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		<span class="comment">// script end tags (&#34;&lt;/script&#34;) within script literals are ignored, so that</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// we can properly escape them.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		if c.element == elementScript &amp;&amp; (isInScriptLiteral(c.state) || isComment(c.state)) {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			return c, len(s)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		if i := indexTagEnd(s, specialTagEndMarkers[c.element]); i != -1 {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			return context{}, i
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	return c, len(s)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// indexTagEnd finds the index of a special tag end in a case insensitive way, or returns -1</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>func indexTagEnd(s []byte, tag []byte) int {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	res := 0
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	plen := len(specialTagEndPrefix)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	for len(s) &gt; 0 {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		<span class="comment">// Try to find the tag end prefix first</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		i := bytes.Index(s, specialTagEndPrefix)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		if i == -1 {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			return i
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		s = s[i+plen:]
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		<span class="comment">// Try to match the actual tag if there is still space for it</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		if len(tag) &lt;= len(s) &amp;&amp; bytes.EqualFold(tag, s[:len(tag)]) {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			s = s[len(tag):]
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			<span class="comment">// Check the tag is followed by a proper separator</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			if len(s) &gt; 0 &amp;&amp; bytes.IndexByte(tagEndSeparators, s[0]) != -1 {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>				return res + i
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			res += len(tag)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		res += i + plen
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	return -1
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// tAttr is the context transition function for the attribute state.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>func tAttr(c context, s []byte) (context, int) {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	return c, len(s)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">// tURL is the context transition function for the URL state.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>func tURL(c context, s []byte) (context, int) {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	if bytes.ContainsAny(s, &#34;#?&#34;) {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		c.urlPart = urlPartQueryOrFrag
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	} else if len(s) != eatWhiteSpace(s, 0) &amp;&amp; c.urlPart == urlPartNone {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		<span class="comment">// HTML5 uses &#34;Valid URL potentially surrounded by spaces&#34; for</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		<span class="comment">// attrs: https://www.w3.org/TR/html5/index.html#attributes-1</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		c.urlPart = urlPartPreQuery
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	return c, len(s)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span><span class="comment">// tJS is the context transition function for the JS state.</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>func tJS(c context, s []byte) (context, int) {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	i := bytes.IndexAny(s, &#34;\&#34;`&#39;/{}&lt;-#&#34;)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	if i == -1 {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		<span class="comment">// Entire input is non string, comment, regexp tokens.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		c.jsCtx = nextJSCtx(s, c.jsCtx)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		return c, len(s)
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	c.jsCtx = nextJSCtx(s[:i], c.jsCtx)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	switch s[i] {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	case &#39;&#34;&#39;:
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		c.state, c.jsCtx = stateJSDqStr, jsCtxRegexp
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	case &#39;\&#39;&#39;:
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		c.state, c.jsCtx = stateJSSqStr, jsCtxRegexp
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	case &#39;`&#39;:
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		c.state, c.jsCtx = stateJSTmplLit, jsCtxRegexp
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	case &#39;/&#39;:
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		switch {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		case i+1 &lt; len(s) &amp;&amp; s[i+1] == &#39;/&#39;:
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			c.state, i = stateJSLineCmt, i+1
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		case i+1 &lt; len(s) &amp;&amp; s[i+1] == &#39;*&#39;:
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			c.state, i = stateJSBlockCmt, i+1
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		case c.jsCtx == jsCtxRegexp:
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			c.state = stateJSRegexp
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		case c.jsCtx == jsCtxDivOp:
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			c.jsCtx = jsCtxRegexp
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		default:
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			return context{
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				state: stateError,
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>				err:   errorf(ErrSlashAmbig, nil, 0, &#34;&#39;/&#39; could start a division or regexp: %.32q&#34;, s[i:]),
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			}, len(s)
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	<span class="comment">// ECMAScript supports HTML style comments for legacy reasons, see Appendix</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	<span class="comment">// B.1.1 &#34;HTML-like Comments&#34;. The handling of these comments is somewhat</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">// confusing. Multi-line comments are not supported, i.e. anything on lines</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	<span class="comment">// between the opening and closing tokens is not considered a comment, but</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	<span class="comment">// anything following the opening or closing token, on the same line, is</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	<span class="comment">// ignored. As such we simply treat any line prefixed with &#34;&lt;!--&#34; or &#34;--&gt;&#34;</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	<span class="comment">// as if it were actually prefixed with &#34;//&#34; and move on.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	case &#39;&lt;&#39;:
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		if i+3 &lt; len(s) &amp;&amp; bytes.Equal(commentStart, s[i:i+4]) {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>			c.state, i = stateJSHTMLOpenCmt, i+3
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	case &#39;-&#39;:
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		if i+2 &lt; len(s) &amp;&amp; bytes.Equal(commentEnd, s[i:i+3]) {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>			c.state, i = stateJSHTMLCloseCmt, i+2
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	<span class="comment">// ECMAScript also supports &#34;hashbang&#34; comment lines, see Section 12.5.</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	case &#39;#&#39;:
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		if i+1 &lt; len(s) &amp;&amp; s[i+1] == &#39;!&#39; {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			c.state, i = stateJSLineCmt, i+1
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	case &#39;{&#39;:
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		<span class="comment">// We only care about tracking brace depth if we are inside of a</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		<span class="comment">// template literal.</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		if len(c.jsBraceDepth) == 0 {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			return c, i + 1
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		c.jsBraceDepth[len(c.jsBraceDepth)-1]++
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	case &#39;}&#39;:
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		if len(c.jsBraceDepth) == 0 {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			return c, i + 1
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		<span class="comment">// There are no cases where a brace can be escaped in the JS context</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		<span class="comment">// that are not syntax errors, it seems. Because of this we can just</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		<span class="comment">// count &#34;\}&#34; as &#34;}&#34; and move on, the script is already broken as</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		<span class="comment">// fully fledged parsers will just fail anyway.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		c.jsBraceDepth[len(c.jsBraceDepth)-1]--
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		if c.jsBraceDepth[len(c.jsBraceDepth)-1] &gt;= 0 {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			return c, i + 1
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		c.jsBraceDepth = c.jsBraceDepth[:len(c.jsBraceDepth)-1]
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		c.state = stateJSTmplLit
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	default:
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		panic(&#34;unreachable&#34;)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	return c, i + 1
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>func tJSTmpl(c context, s []byte) (context, int) {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	var k int
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	for {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		i := k + bytes.IndexAny(s[k:], &#34;`\\$&#34;)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		if i &lt; k {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			break
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		switch s[i] {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		case &#39;\\&#39;:
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			i++
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			if i == len(s) {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>				return context{
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>					state: stateError,
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>					err:   errorf(ErrPartialEscape, nil, 0, &#34;unfinished escape sequence in JS string: %q&#34;, s),
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>				}, len(s)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		case &#39;$&#39;:
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			if len(s) &gt;= i+2 &amp;&amp; s[i+1] == &#39;{&#39; {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>				c.jsBraceDepth = append(c.jsBraceDepth, 0)
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>				c.state = stateJS
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>				return c, i + 2
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>			}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		case &#39;`&#39;:
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			<span class="comment">// end</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			c.state = stateJS
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			return c, i + 1
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		k = i + 1
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	return c, len(s)
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span><span class="comment">// tJSDelimited is the context transition function for the JS string and regexp</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span><span class="comment">// states.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>func tJSDelimited(c context, s []byte) (context, int) {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	specials := `\&#34;`
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	switch c.state {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	case stateJSSqStr:
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		specials = `\&#39;`
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	case stateJSRegexp:
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		specials = `\/[]`
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	k, inCharset := 0, false
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	for {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		i := k + bytes.IndexAny(s[k:], specials)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		if i &lt; k {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			break
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		switch s[i] {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		case &#39;\\&#39;:
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			i++
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			if i == len(s) {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>				return context{
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>					state: stateError,
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>					err:   errorf(ErrPartialEscape, nil, 0, &#34;unfinished escape sequence in JS string: %q&#34;, s),
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>				}, len(s)
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		case &#39;[&#39;:
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			inCharset = true
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		case &#39;]&#39;:
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			inCharset = false
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		case &#39;/&#39;:
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			<span class="comment">// If &#34;&lt;/script&#34; appears in a regex literal, the &#39;/&#39; should not</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			<span class="comment">// close the regex literal, and it will later be escaped to</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			<span class="comment">// &#34;\x3C/script&#34; in escapeText.</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			if i &gt; 0 &amp;&amp; i+7 &lt;= len(s) &amp;&amp; bytes.Compare(bytes.ToLower(s[i-1:i+7]), []byte(&#34;&lt;/script&#34;)) == 0 {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>				i++
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			} else if !inCharset {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>				c.state, c.jsCtx = stateJS, jsCtxDivOp
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>				return c, i + 1
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		default:
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			<span class="comment">// end delimiter</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>			if !inCharset {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>				c.state, c.jsCtx = stateJS, jsCtxDivOp
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>				return c, i + 1
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		k = i + 1
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	if inCharset {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		<span class="comment">// This can be fixed by making context richer if interpolation</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		<span class="comment">// into charsets is desired.</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		return context{
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			state: stateError,
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>			err:   errorf(ErrPartialCharset, nil, 0, &#34;unfinished JS regexp charset: %q&#34;, s),
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		}, len(s)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	return c, len(s)
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>var blockCommentEnd = []byte(&#34;*/&#34;)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span><span class="comment">// tBlockCmt is the context transition function for /*comment*/ states.</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>func tBlockCmt(c context, s []byte) (context, int) {
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	i := bytes.Index(s, blockCommentEnd)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	if i == -1 {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		return c, len(s)
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	switch c.state {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	case stateJSBlockCmt:
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		c.state = stateJS
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	case stateCSSBlockCmt:
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		c.state = stateCSS
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	default:
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		panic(c.state.String())
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	}
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	return c, i + 2
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">// tLineCmt is the context transition function for //comment states, and the JS HTML-like comment state.</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>func tLineCmt(c context, s []byte) (context, int) {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	var lineTerminators string
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	var endState state
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	switch c.state {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	case stateJSLineCmt, stateJSHTMLOpenCmt, stateJSHTMLCloseCmt:
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		lineTerminators, endState = &#34;\n\r\u2028\u2029&#34;, stateJS
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	case stateCSSLineCmt:
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		lineTerminators, endState = &#34;\n\f\r&#34;, stateCSS
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		<span class="comment">// Line comments are not part of any published CSS standard but</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		<span class="comment">// are supported by the 4 major browsers.</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		<span class="comment">// This defines line comments as</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		<span class="comment">//     LINECOMMENT ::= &#34;//&#34; [^\n\f\d]*</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		<span class="comment">// since https://www.w3.org/TR/css3-syntax/#SUBTOK-nl defines</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		<span class="comment">// newlines:</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		<span class="comment">//     nl ::= #xA | #xD #xA | #xD | #xC</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	default:
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		panic(c.state.String())
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	i := bytes.IndexAny(s, lineTerminators)
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	if i == -1 {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		return c, len(s)
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	c.state = endState
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	<span class="comment">// Per section 7.4 of EcmaScript 5 : https://es5.github.io/#x7.4</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	<span class="comment">// &#34;However, the LineTerminator at the end of the line is not</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	<span class="comment">// considered to be part of the single-line comment; it is</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	<span class="comment">// recognized separately by the lexical grammar and becomes part</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	<span class="comment">// of the stream of input elements for the syntactic grammar.&#34;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	return c, i
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span><span class="comment">// tCSS is the context transition function for the CSS state.</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>func tCSS(c context, s []byte) (context, int) {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	<span class="comment">// CSS quoted strings are almost never used except for:</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	<span class="comment">// (1) URLs as in background: &#34;/foo.png&#34;</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	<span class="comment">// (2) Multiword font-names as in font-family: &#34;Times New Roman&#34;</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	<span class="comment">// (3) List separators in content values as in inline-lists:</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	<span class="comment">//    &lt;style&gt;</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	<span class="comment">//    ul.inlineList { list-style: none; padding:0 }</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	<span class="comment">//    ul.inlineList &gt; li { display: inline }</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	<span class="comment">//    ul.inlineList &gt; li:before { content: &#34;, &#34; }</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	<span class="comment">//    ul.inlineList &gt; li:first-child:before { content: &#34;&#34; }</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	<span class="comment">//    &lt;/style&gt;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	<span class="comment">//    &lt;ul class=inlineList&gt;&lt;li&gt;One&lt;li&gt;Two&lt;li&gt;Three&lt;/ul&gt;</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	<span class="comment">// (4) Attribute value selectors as in a[href=&#34;http://example.com/&#34;]</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	<span class="comment">// We conservatively treat all strings as URLs, but make some</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	<span class="comment">// allowances to avoid confusion.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	<span class="comment">// In (1), our conservative assumption is justified.</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	<span class="comment">// In (2), valid font names do not contain &#39;:&#39;, &#39;?&#39;, or &#39;#&#39;, so our</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	<span class="comment">// conservative assumption is fine since we will never transition past</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	<span class="comment">// urlPartPreQuery.</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	<span class="comment">// In (3), our protocol heuristic should not be tripped, and there</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	<span class="comment">// should not be non-space content after a &#39;?&#39; or &#39;#&#39;, so as long as</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	<span class="comment">// we only %-encode RFC 3986 reserved characters we are ok.</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	<span class="comment">// In (4), we should URL escape for URL attributes, and for others we</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	<span class="comment">// have the attribute name available if our conservative assumption</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	<span class="comment">// proves problematic for real code.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	k := 0
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	for {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		i := k + bytes.IndexAny(s[k:], `(&#34;&#39;/`)
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		if i &lt; k {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			return c, len(s)
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		}
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		switch s[i] {
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		case &#39;(&#39;:
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>			<span class="comment">// Look for url to the left.</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			p := bytes.TrimRight(s[:i], &#34;\t\n\f\r &#34;)
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			if endsWithCSSKeyword(p, &#34;url&#34;) {
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>				j := len(s) - len(bytes.TrimLeft(s[i+1:], &#34;\t\n\f\r &#34;))
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>				switch {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>				case j != len(s) &amp;&amp; s[j] == &#39;&#34;&#39;:
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>					c.state, j = stateCSSDqURL, j+1
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>				case j != len(s) &amp;&amp; s[j] == &#39;\&#39;&#39;:
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>					c.state, j = stateCSSSqURL, j+1
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>				default:
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>					c.state = stateCSSURL
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>				}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>				return c, j
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>			}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		case &#39;/&#39;:
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>			if i+1 &lt; len(s) {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>				switch s[i+1] {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>				case &#39;/&#39;:
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>					c.state = stateCSSLineCmt
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>					return c, i + 2
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>				case &#39;*&#39;:
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>					c.state = stateCSSBlockCmt
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>					return c, i + 2
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>				}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		case &#39;&#34;&#39;:
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			c.state = stateCSSDqStr
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>			return c, i + 1
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		case &#39;\&#39;&#39;:
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>			c.state = stateCSSSqStr
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>			return c, i + 1
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		}
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		k = i + 1
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span><span class="comment">// tCSSStr is the context transition function for the CSS string and URL states.</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>func tCSSStr(c context, s []byte) (context, int) {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	var endAndEsc string
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	switch c.state {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	case stateCSSDqStr, stateCSSDqURL:
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		endAndEsc = `\&#34;`
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	case stateCSSSqStr, stateCSSSqURL:
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		endAndEsc = `\&#39;`
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	case stateCSSURL:
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		<span class="comment">// Unquoted URLs end with a newline or close parenthesis.</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		<span class="comment">// The below includes the wc (whitespace character) and nl.</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		endAndEsc = &#34;\\\t\n\f\r )&#34;
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	default:
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		panic(c.state.String())
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	k := 0
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	for {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		i := k + bytes.IndexAny(s[k:], endAndEsc)
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		if i &lt; k {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>			c, nread := tURL(c, decodeCSS(s[k:]))
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>			return c, k + nread
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		}
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		if s[i] == &#39;\\&#39; {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			i++
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>			if i == len(s) {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>				return context{
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>					state: stateError,
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>					err:   errorf(ErrPartialEscape, nil, 0, &#34;unfinished escape sequence in CSS string: %q&#34;, s),
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>				}, len(s)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>			}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		} else {
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>			c.state = stateCSS
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>			return c, i + 1
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		}
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		c, _ = tURL(c, decodeCSS(s[:i+1]))
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		k = i + 1
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	}
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span><span class="comment">// tError is the context transition function for the error state.</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>func tError(c context, s []byte) (context, int) {
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	return c, len(s)
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>}
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span><span class="comment">// eatAttrName returns the largest j such that s[i:j] is an attribute name.</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span><span class="comment">// It returns an error if s[i:] does not look like it begins with an</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span><span class="comment">// attribute name, such as encountering a quote mark without a preceding</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span><span class="comment">// equals sign.</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>func eatAttrName(s []byte, i int) (int, *Error) {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	for j := i; j &lt; len(s); j++ {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		switch s[j] {
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		case &#39; &#39;, &#39;\t&#39;, &#39;\n&#39;, &#39;\f&#39;, &#39;\r&#39;, &#39;=&#39;, &#39;&gt;&#39;:
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>			return j, nil
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		case &#39;\&#39;&#39;, &#39;&#34;&#39;, &#39;&lt;&#39;:
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>			<span class="comment">// These result in a parse warning in HTML5 and are</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>			<span class="comment">// indicative of serious problems if seen in an attr</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			<span class="comment">// name in a template.</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			return -1, errorf(ErrBadHTML, nil, 0, &#34;%q in attribute name: %.32q&#34;, s[j:j+1], s)
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		default:
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>			<span class="comment">// No-op.</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		}
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	return len(s), nil
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>}
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>var elementNameMap = map[string]element{
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	&#34;script&#34;:   elementScript,
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	&#34;style&#34;:    elementStyle,
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	&#34;textarea&#34;: elementTextarea,
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	&#34;title&#34;:    elementTitle,
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span><span class="comment">// asciiAlpha reports whether c is an ASCII letter.</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>func asciiAlpha(c byte) bool {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	return &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;Z&#39; || &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;z&#39;
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span><span class="comment">// asciiAlphaNum reports whether c is an ASCII letter or digit.</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>func asciiAlphaNum(c byte) bool {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	return asciiAlpha(c) || &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39;
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span><span class="comment">// eatTagName returns the largest j such that s[i:j] is a tag name and the tag type.</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>func eatTagName(s []byte, i int) (int, element) {
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	if i == len(s) || !asciiAlpha(s[i]) {
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		return i, elementNone
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	}
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	j := i + 1
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	for j &lt; len(s) {
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		x := s[j]
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		if asciiAlphaNum(x) {
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>			j++
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>			continue
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		}
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		<span class="comment">// Allow &#34;x-y&#34; or &#34;x:y&#34; but not &#34;x-&#34;, &#34;-y&#34;, or &#34;x--y&#34;.</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		if (x == &#39;:&#39; || x == &#39;-&#39;) &amp;&amp; j+1 &lt; len(s) &amp;&amp; asciiAlphaNum(s[j+1]) {
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>			j += 2
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>			continue
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		}
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		break
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	}
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	return j, elementNameMap[strings.ToLower(string(s[i:j]))]
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span><span class="comment">// eatWhiteSpace returns the largest j such that s[i:j] is white space.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>func eatWhiteSpace(s []byte, i int) int {
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	for j := i; j &lt; len(s); j++ {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		switch s[j] {
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		case &#39; &#39;, &#39;\t&#39;, &#39;\n&#39;, &#39;\f&#39;, &#39;\r&#39;:
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>			<span class="comment">// No-op.</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		default:
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			return j
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		}
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	return len(s)
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>
</pre><p><a href="transition.go?m=text">View as plain text</a></p>

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
