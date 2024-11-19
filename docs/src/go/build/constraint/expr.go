<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/build/constraint/expr.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../../index.html">GoDoc</a></div>
<a href="expr.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/build">build</a>/<a href="http://localhost:8080/src/go/build/constraint">constraint</a>/<span class="text-muted">expr.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/build/constraint">go/build/constraint</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2020 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package constraint implements parsing and evaluation of build constraint lines.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// See https://golang.org/cmd/go/#hdr-Build_constraints for documentation about build constraints themselves.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// This package parses both the original “// +build” syntax and the “//go:build” syntax that was added in Go 1.17.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// See https://golang.org/design/draft-gobuild for details about the “//go:build” syntax.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>package constraint
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>import (
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// maxSize is a limit used to control the complexity of expressions, in order</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// to prevent stack exhaustion issues due to recursion.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>const maxSize = 1000
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// An Expr is a build tag constraint expression.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// The underlying concrete type is *[AndExpr], *[OrExpr], *[NotExpr], or *[TagExpr].</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>type Expr interface {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// String returns the string form of the expression,</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// using the boolean syntax used in //go:build lines.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	String() string
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// Eval reports whether the expression evaluates to true.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// It calls ok(tag) as needed to find out whether a given build tag</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// is satisfied by the current build configuration.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	Eval(ok func(tag string) bool) bool
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// The presence of an isExpr method explicitly marks the type as an Expr.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// Only implementations in this package should be used as Exprs.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	isExpr()
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// A TagExpr is an [Expr] for the single tag Tag.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>type TagExpr struct {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	Tag string <span class="comment">// for example, “linux” or “cgo”</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>func (x *TagExpr) isExpr() {}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>func (x *TagExpr) Eval(ok func(tag string) bool) bool {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	return ok(x.Tag)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>func (x *TagExpr) String() string {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	return x.Tag
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>func tag(tag string) Expr { return &amp;TagExpr{tag} }
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// A NotExpr represents the expression !X (the negation of X).</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>type NotExpr struct {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	X Expr
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>func (x *NotExpr) isExpr() {}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>func (x *NotExpr) Eval(ok func(tag string) bool) bool {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	return !x.X.Eval(ok)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>func (x *NotExpr) String() string {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	s := x.X.String()
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	switch x.X.(type) {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	case *AndExpr, *OrExpr:
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		s = &#34;(&#34; + s + &#34;)&#34;
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return &#34;!&#34; + s
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>func not(x Expr) Expr { return &amp;NotExpr{x} }
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// An AndExpr represents the expression X &amp;&amp; Y.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>type AndExpr struct {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	X, Y Expr
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>func (x *AndExpr) isExpr() {}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>func (x *AndExpr) Eval(ok func(tag string) bool) bool {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// Note: Eval both, to make sure ok func observes all tags.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	xok := x.X.Eval(ok)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	yok := x.Y.Eval(ok)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	return xok &amp;&amp; yok
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>func (x *AndExpr) String() string {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	return andArg(x.X) + &#34; &amp;&amp; &#34; + andArg(x.Y)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func andArg(x Expr) string {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	s := x.String()
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	if _, ok := x.(*OrExpr); ok {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		s = &#34;(&#34; + s + &#34;)&#34;
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	return s
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func and(x, y Expr) Expr {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	return &amp;AndExpr{x, y}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// An OrExpr represents the expression X || Y.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>type OrExpr struct {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	X, Y Expr
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>func (x *OrExpr) isExpr() {}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>func (x *OrExpr) Eval(ok func(tag string) bool) bool {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// Note: Eval both, to make sure ok func observes all tags.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	xok := x.X.Eval(ok)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	yok := x.Y.Eval(ok)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	return xok || yok
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>func (x *OrExpr) String() string {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	return orArg(x.X) + &#34; || &#34; + orArg(x.Y)
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>func orArg(x Expr) string {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	s := x.String()
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if _, ok := x.(*AndExpr); ok {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		s = &#34;(&#34; + s + &#34;)&#34;
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	return s
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>func or(x, y Expr) Expr {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	return &amp;OrExpr{x, y}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// A SyntaxError reports a syntax error in a parsed build expression.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>type SyntaxError struct {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	Offset int    <span class="comment">// byte offset in input where error was detected</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	Err    string <span class="comment">// description of error</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>func (e *SyntaxError) Error() string {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	return e.Err
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>var errNotConstraint = errors.New(&#34;not a build constraint&#34;)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// Parse parses a single build constraint line of the form “//go:build ...” or “// +build ...”</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">// and returns the corresponding boolean expression.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>func Parse(line string) (Expr, error) {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	if text, ok := splitGoBuild(line); ok {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		return parseExpr(text)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	if text, ok := splitPlusBuild(line); ok {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		return parsePlusBuildExpr(text)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	return nil, errNotConstraint
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// IsGoBuild reports whether the line of text is a “//go:build” constraint.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// It only checks the prefix of the text, not that the expression itself parses.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>func IsGoBuild(line string) bool {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	_, ok := splitGoBuild(line)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	return ok
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">// splitGoBuild splits apart the leading //go:build prefix in line from the build expression itself.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// It returns &#34;&#34;, false if the input is not a //go:build line or if the input contains multiple lines.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>func splitGoBuild(line string) (expr string, ok bool) {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// A single trailing newline is OK; otherwise multiple lines are not.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	if len(line) &gt; 0 &amp;&amp; line[len(line)-1] == &#39;\n&#39; {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		line = line[:len(line)-1]
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	if strings.Contains(line, &#34;\n&#34;) {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		return &#34;&#34;, false
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	if !strings.HasPrefix(line, &#34;//go:build&#34;) {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		return &#34;&#34;, false
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	line = strings.TrimSpace(line)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	line = line[len(&#34;//go:build&#34;):]
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">// If strings.TrimSpace finds more to trim after removing the //go:build prefix,</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// it means that the prefix was followed by a space, making this a //go:build line</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// (as opposed to a //go:buildsomethingelse line).</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// If line is empty, we had &#34;//go:build&#34; by itself, which also counts.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	trim := strings.TrimSpace(line)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	if len(line) == len(trim) &amp;&amp; line != &#34;&#34; {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		return &#34;&#34;, false
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	return trim, true
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// An exprParser holds state for parsing a build expression.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>type exprParser struct {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	s string <span class="comment">// input string</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	i int    <span class="comment">// next read location in s</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	tok   string <span class="comment">// last token read</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	isTag bool
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	pos   int <span class="comment">// position (start) of last token</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	size int
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// parseExpr parses a boolean build tag expression.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>func parseExpr(text string) (x Expr, err error) {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	defer func() {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		if e := recover(); e != nil {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			if e, ok := e.(*SyntaxError); ok {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>				err = e
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>				return
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			panic(e) <span class="comment">// unreachable unless parser has a bug</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}()
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	p := &amp;exprParser{s: text}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	x = p.or()
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	if p.tok != &#34;&#34; {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		panic(&amp;SyntaxError{Offset: p.pos, Err: &#34;unexpected token &#34; + p.tok})
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	return x, nil
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// or parses a sequence of || expressions.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">// On entry, the next input token has not yet been lexed.</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">// On exit, the next input token has been lexed and is in p.tok.</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>func (p *exprParser) or() Expr {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	x := p.and()
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	for p.tok == &#34;||&#34; {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		x = or(x, p.and())
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	return x
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span><span class="comment">// and parses a sequence of &amp;&amp; expressions.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">// On entry, the next input token has not yet been lexed.</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// On exit, the next input token has been lexed and is in p.tok.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>func (p *exprParser) and() Expr {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	x := p.not()
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	for p.tok == &#34;&amp;&amp;&#34; {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		x = and(x, p.not())
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	return x
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// not parses a ! expression.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// On entry, the next input token has not yet been lexed.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// On exit, the next input token has been lexed and is in p.tok.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>func (p *exprParser) not() Expr {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	p.size++
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	if p.size &gt; maxSize {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		panic(&amp;SyntaxError{Offset: p.pos, Err: &#34;build expression too large&#34;})
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	p.lex()
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	if p.tok == &#34;!&#34; {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		p.lex()
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		if p.tok == &#34;!&#34; {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			panic(&amp;SyntaxError{Offset: p.pos, Err: &#34;double negation not allowed&#34;})
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		return not(p.atom())
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	return p.atom()
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span><span class="comment">// atom parses a tag or a parenthesized expression.</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span><span class="comment">// On entry, the next input token HAS been lexed.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span><span class="comment">// On exit, the next input token has been lexed and is in p.tok.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>func (p *exprParser) atom() Expr {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	<span class="comment">// first token already in p.tok</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	if p.tok == &#34;(&#34; {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		pos := p.pos
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		defer func() {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			if e := recover(); e != nil {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>				if e, ok := e.(*SyntaxError); ok &amp;&amp; e.Err == &#34;unexpected end of expression&#34; {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>					e.Err = &#34;missing close paren&#34;
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>				}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>				panic(e)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		}()
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		x := p.or()
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		if p.tok != &#34;)&#34; {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			panic(&amp;SyntaxError{Offset: pos, Err: &#34;missing close paren&#34;})
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		p.lex()
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		return x
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	if !p.isTag {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		if p.tok == &#34;&#34; {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			panic(&amp;SyntaxError{Offset: p.pos, Err: &#34;unexpected end of expression&#34;})
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		panic(&amp;SyntaxError{Offset: p.pos, Err: &#34;unexpected token &#34; + p.tok})
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	tok := p.tok
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	p.lex()
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	return tag(tok)
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">// lex finds and consumes the next token in the input stream.</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// On return, p.tok is set to the token text,</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// p.isTag reports whether the token was a tag,</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// and p.pos records the byte offset of the start of the token in the input stream.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// If lex reaches the end of the input, p.tok is set to the empty string.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">// For any other syntax error, lex panics with a SyntaxError.</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>func (p *exprParser) lex() {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	p.isTag = false
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	for p.i &lt; len(p.s) &amp;&amp; (p.s[p.i] == &#39; &#39; || p.s[p.i] == &#39;\t&#39;) {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		p.i++
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	if p.i &gt;= len(p.s) {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		p.tok = &#34;&#34;
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		p.pos = p.i
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		return
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	switch p.s[p.i] {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	case &#39;(&#39;, &#39;)&#39;, &#39;!&#39;:
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		p.pos = p.i
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		p.i++
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		p.tok = p.s[p.pos:p.i]
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		return
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	case &#39;&amp;&#39;, &#39;|&#39;:
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		if p.i+1 &gt;= len(p.s) || p.s[p.i+1] != p.s[p.i] {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			panic(&amp;SyntaxError{Offset: p.i, Err: &#34;invalid syntax at &#34; + string(rune(p.s[p.i]))})
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		p.pos = p.i
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		p.i += 2
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		p.tok = p.s[p.pos:p.i]
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		return
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	tag := p.s[p.i:]
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	for i, c := range tag {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		if !unicode.IsLetter(c) &amp;&amp; !unicode.IsDigit(c) &amp;&amp; c != &#39;_&#39; &amp;&amp; c != &#39;.&#39; {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			tag = tag[:i]
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			break
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	if tag == &#34;&#34; {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		c, _ := utf8.DecodeRuneInString(p.s[p.i:])
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		panic(&amp;SyntaxError{Offset: p.i, Err: &#34;invalid syntax at &#34; + string(c)})
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	p.pos = p.i
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	p.i += len(tag)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	p.tok = p.s[p.pos:p.i]
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	p.isTag = true
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">// IsPlusBuild reports whether the line of text is a “// +build” constraint.</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// It only checks the prefix of the text, not that the expression itself parses.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>func IsPlusBuild(line string) bool {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	_, ok := splitPlusBuild(line)
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	return ok
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// splitPlusBuild splits apart the leading // +build prefix in line from the build expression itself.</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// It returns &#34;&#34;, false if the input is not a // +build line or if the input contains multiple lines.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>func splitPlusBuild(line string) (expr string, ok bool) {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	<span class="comment">// A single trailing newline is OK; otherwise multiple lines are not.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	if len(line) &gt; 0 &amp;&amp; line[len(line)-1] == &#39;\n&#39; {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		line = line[:len(line)-1]
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	if strings.Contains(line, &#34;\n&#34;) {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		return &#34;&#34;, false
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	if !strings.HasPrefix(line, &#34;//&#34;) {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		return &#34;&#34;, false
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	line = line[len(&#34;//&#34;):]
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	<span class="comment">// Note the space is optional; &#34;//+build&#34; is recognized too.</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	line = strings.TrimSpace(line)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	if !strings.HasPrefix(line, &#34;+build&#34;) {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		return &#34;&#34;, false
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	line = line[len(&#34;+build&#34;):]
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	<span class="comment">// If strings.TrimSpace finds more to trim after removing the +build prefix,</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	<span class="comment">// it means that the prefix was followed by a space, making this a +build line</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">// (as opposed to a +buildsomethingelse line).</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	<span class="comment">// If line is empty, we had &#34;// +build&#34; by itself, which also counts.</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	trim := strings.TrimSpace(line)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	if len(line) == len(trim) &amp;&amp; line != &#34;&#34; {
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		return &#34;&#34;, false
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	}
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	return trim, true
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span><span class="comment">// parsePlusBuildExpr parses a legacy build tag expression (as used with “// +build”).</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>func parsePlusBuildExpr(text string) (Expr, error) {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	<span class="comment">// Only allow up to 100 AND/OR operators for &#34;old&#34; syntax.</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	<span class="comment">// This is much less than the limit for &#34;new&#34; syntax,</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	<span class="comment">// but uses of old syntax were always very simple.</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	const maxOldSize = 100
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	size := 0
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	var x Expr
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	for _, clause := range strings.Fields(text) {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		var y Expr
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		for _, lit := range strings.Split(clause, &#34;,&#34;) {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			var z Expr
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			var neg bool
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			if strings.HasPrefix(lit, &#34;!!&#34;) || lit == &#34;!&#34; {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>				z = tag(&#34;ignore&#34;)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			} else {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>				if strings.HasPrefix(lit, &#34;!&#34;) {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>					neg = true
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>					lit = lit[len(&#34;!&#34;):]
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>				}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>				if isValidTag(lit) {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>					z = tag(lit)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>				} else {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>					z = tag(&#34;ignore&#34;)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>				}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>				if neg {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>					z = not(z)
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>				}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>			}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			if y == nil {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>				y = z
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			} else {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>				if size++; size &gt; maxOldSize {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>					return nil, errComplex
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>				}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>				y = and(y, z)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		if x == nil {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			x = y
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		} else {
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			if size++; size &gt; maxOldSize {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>				return nil, errComplex
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			x = or(x, y)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	if x == nil {
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		x = tag(&#34;ignore&#34;)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	return x, nil
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span><span class="comment">// isValidTag reports whether the word is a valid build tag.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span><span class="comment">// Tags must be letters, digits, underscores or dots.</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span><span class="comment">// Unlike in Go identifiers, all digits are fine (e.g., &#34;386&#34;).</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>func isValidTag(word string) bool {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	if word == &#34;&#34; {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		return false
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	}
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	for _, c := range word {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		if !unicode.IsLetter(c) &amp;&amp; !unicode.IsDigit(c) &amp;&amp; c != &#39;_&#39; &amp;&amp; c != &#39;.&#39; {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			return false
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		}
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	return true
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>var errComplex = errors.New(&#34;expression too complex for // +build lines&#34;)
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span><span class="comment">// PlusBuildLines returns a sequence of “// +build” lines that evaluate to the build expression x.</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span><span class="comment">// If the expression is too complex to convert directly to “// +build” lines, PlusBuildLines returns an error.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>func PlusBuildLines(x Expr) ([]string, error) {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	<span class="comment">// Push all NOTs to the expression leaves, so that //go:build !(x &amp;&amp; y) can be treated as !x || !y.</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	<span class="comment">// This rewrite is both efficient and commonly needed, so it&#39;s worth doing.</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	<span class="comment">// Essentially all other possible rewrites are too expensive and too rarely needed.</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	x = pushNot(x, false)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	<span class="comment">// Split into AND of ORs of ANDs of literals (tag or NOT tag).</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	var split [][][]Expr
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	for _, or := range appendSplitAnd(nil, x) {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		var ands [][]Expr
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		for _, and := range appendSplitOr(nil, or) {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			var lits []Expr
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			for _, lit := range appendSplitAnd(nil, and) {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>				switch lit.(type) {
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>				case *TagExpr, *NotExpr:
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>					lits = append(lits, lit)
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>				default:
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>					return nil, errComplex
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>				}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>			}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>			ands = append(ands, lits)
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		}
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		split = append(split, ands)
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	<span class="comment">// If all the ORs have length 1 (no actual OR&#39;ing going on),</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	<span class="comment">// push the top-level ANDs to the bottom level, so that we get</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	<span class="comment">// one // +build line instead of many.</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	maxOr := 0
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	for _, or := range split {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		if maxOr &lt; len(or) {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>			maxOr = len(or)
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	if maxOr == 1 {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		var lits []Expr
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		for _, or := range split {
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			lits = append(lits, or[0]...)
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		}
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		split = [][][]Expr{{lits}}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	<span class="comment">// Prepare the +build lines.</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	var lines []string
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	for _, or := range split {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		line := &#34;// +build&#34;
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		for _, and := range or {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			clause := &#34;&#34;
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			for i, lit := range and {
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>				if i &gt; 0 {
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>					clause += &#34;,&#34;
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>				}
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>				clause += lit.String()
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>			}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			line += &#34; &#34; + clause
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		}
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		lines = append(lines, line)
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	}
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	return lines, nil
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span><span class="comment">// pushNot applies DeMorgan&#39;s law to push negations down the expression,</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span><span class="comment">// so that only tags are negated in the result.</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span><span class="comment">// (It applies the rewrites !(X &amp;&amp; Y) =&gt; (!X || !Y) and !(X || Y) =&gt; (!X &amp;&amp; !Y).)</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>func pushNot(x Expr, not bool) Expr {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	switch x := x.(type) {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	default:
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		<span class="comment">// unreachable</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		return x
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	case *NotExpr:
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		if _, ok := x.X.(*TagExpr); ok &amp;&amp; !not {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>			return x
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		}
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		return pushNot(x.X, !not)
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	case *TagExpr:
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		if not {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>			return &amp;NotExpr{X: x}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		return x
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	case *AndExpr:
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		x1 := pushNot(x.X, not)
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		y1 := pushNot(x.Y, not)
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		if not {
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>			return or(x1, y1)
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		if x1 == x.X &amp;&amp; y1 == x.Y {
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			return x
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		return and(x1, y1)
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	case *OrExpr:
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		x1 := pushNot(x.X, not)
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		y1 := pushNot(x.Y, not)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		if not {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>			return and(x1, y1)
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		if x1 == x.X &amp;&amp; y1 == x.Y {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>			return x
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		return or(x1, y1)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span><span class="comment">// appendSplitAnd appends x to list while splitting apart any top-level &amp;&amp; expressions.</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span><span class="comment">// For example, appendSplitAnd({W}, X &amp;&amp; Y &amp;&amp; Z) = {W, X, Y, Z}.</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>func appendSplitAnd(list []Expr, x Expr) []Expr {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	if x, ok := x.(*AndExpr); ok {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		list = appendSplitAnd(list, x.X)
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		list = appendSplitAnd(list, x.Y)
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		return list
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	}
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	return append(list, x)
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span><span class="comment">// appendSplitOr appends x to list while splitting apart any top-level || expressions.</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span><span class="comment">// For example, appendSplitOr({W}, X || Y || Z) = {W, X, Y, Z}.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>func appendSplitOr(list []Expr, x Expr) []Expr {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	if x, ok := x.(*OrExpr); ok {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		list = appendSplitOr(list, x.X)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		list = appendSplitOr(list, x.Y)
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		return list
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	return append(list, x)
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>}
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>
</pre><p><a href="expr.go?m=text">View as plain text</a></p>

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
