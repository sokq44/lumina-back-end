<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/stmt.go - Go Documentation Server</title>

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
<a href="stmt.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">stmt.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/types">go/types</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2012 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements typechecking of statements.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package types
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/constant&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;internal/buildcfg&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>func (check *Checker) funcBody(decl *declInfo, name string, sig *Signature, body *ast.BlockStmt, iota constant.Value) {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	if check.conf.IgnoreFuncBodies {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>		panic(&#34;function body not ignored&#34;)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	}
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	if check.conf._Trace {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>		check.trace(body.Pos(), &#34;-- %s: %s&#34;, name, sig)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	}
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// save/restore current environment and set up function environment</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// (and use 0 indentation at function start)</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	defer func(env environment, indent int) {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		check.environment = env
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		check.indent = indent
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	}(check.environment, check.indent)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	check.environment = environment{
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		decl:  decl,
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		scope: sig.scope,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		iota:  iota,
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		sig:   sig,
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	check.indent = 0
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	check.stmtList(0, body.List)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	if check.hasLabel {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		check.labels(body)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	if sig.results.Len() &gt; 0 &amp;&amp; !check.isTerminating(body, &#34;&#34;) {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		check.error(atPos(body.Rbrace), MissingReturn, &#34;missing return&#34;)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;Implementation restriction: A compiler may make it illegal to</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// declare a variable inside a function body if the variable is never used.&#34;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	check.usage(sig.scope)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>func (check *Checker) usage(scope *Scope) {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	var unused []*Var
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	for name, elem := range scope.elems {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		elem = resolve(name, elem)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		if v, _ := elem.(*Var); v != nil &amp;&amp; !v.used {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			unused = append(unused, v)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	sort.Slice(unused, func(i, j int) bool {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		return cmpPos(unused[i].pos, unused[j].pos) &lt; 0
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	})
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	for _, v := range unused {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		check.softErrorf(v, UnusedVar, &#34;%s declared and not used&#34;, v.name)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	for _, scope := range scope.children {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		<span class="comment">// Don&#39;t go inside function literal scopes a second time;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		<span class="comment">// they are handled explicitly by funcBody.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		if !scope.isFunc {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			check.usage(scope)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// stmtContext is a bitset describing which</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// control-flow statements are permissible,</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// and provides additional context information</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// for better error messages.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>type stmtContext uint
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>const (
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// permissible control-flow statements</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	breakOk stmtContext = 1 &lt;&lt; iota
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	continueOk
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	fallthroughOk
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// additional context information</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	finalSwitchCase
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	inTypeSwitch
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func (check *Checker) simpleStmt(s ast.Stmt) {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	if s != nil {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		check.stmt(0, s)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>func trimTrailingEmptyStmts(list []ast.Stmt) []ast.Stmt {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	for i := len(list); i &gt; 0; i-- {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		if _, ok := list[i-1].(*ast.EmptyStmt); !ok {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			return list[:i]
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	return nil
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>func (check *Checker) stmtList(ctxt stmtContext, list []ast.Stmt) {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	ok := ctxt&amp;fallthroughOk != 0
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	inner := ctxt &amp;^ fallthroughOk
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	list = trimTrailingEmptyStmts(list) <span class="comment">// trailing empty statements are &#34;invisible&#34; to fallthrough analysis</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	for i, s := range list {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		inner := inner
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		if ok &amp;&amp; i+1 == len(list) {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			inner |= fallthroughOk
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		check.stmt(inner, s)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>func (check *Checker) multipleDefaults(list []ast.Stmt) {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	var first ast.Stmt
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	for _, s := range list {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		var d ast.Stmt
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		switch c := s.(type) {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		case *ast.CaseClause:
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			if len(c.List) == 0 {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>				d = s
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		case *ast.CommClause:
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			if c.Comm == nil {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>				d = s
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		default:
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			check.error(s, InvalidSyntaxTree, &#34;case/communication clause expected&#34;)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		if d != nil {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			if first != nil {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>				check.errorf(d, DuplicateDefault, &#34;multiple defaults (first at %s)&#34;, check.fset.Position(first.Pos()))
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			} else {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>				first = d
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>func (check *Checker) openScope(node ast.Node, comment string) {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	scope := NewScope(check.scope, node.Pos(), node.End(), comment)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	check.recordScope(node, scope)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	check.scope = scope
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>func (check *Checker) closeScope() {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	check.scope = check.scope.Parent()
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>func assignOp(op token.Token) token.Token {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// token_test.go verifies the token ordering this function relies on</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if token.ADD_ASSIGN &lt;= op &amp;&amp; op &lt;= token.AND_NOT_ASSIGN {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		return op + (token.ADD - token.ADD_ASSIGN)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	return token.ILLEGAL
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>func (check *Checker) suspendedCall(keyword string, call *ast.CallExpr) {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	var x operand
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	var msg string
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	var code Code
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	switch check.rawExpr(nil, &amp;x, call, nil, false) {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	case conversion:
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		msg = &#34;requires function call, not conversion&#34;
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		code = InvalidDefer
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		if keyword == &#34;go&#34; {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			code = InvalidGo
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	case expression:
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		msg = &#34;discards result of&#34;
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		code = UnusedResults
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	case statement:
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		return
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	default:
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		unreachable()
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	check.errorf(&amp;x, code, &#34;%s %s %s&#34;, keyword, msg, &amp;x)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// goVal returns the Go value for val, or nil.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>func goVal(val constant.Value) any {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// val should exist, but be conservative and check</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	if val == nil {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		return nil
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// Match implementation restriction of other compilers.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// gc only checks duplicates for integer, floating-point</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// and string values, so only create Go values for these</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// types.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	switch val.Kind() {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	case constant.Int:
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		if x, ok := constant.Int64Val(val); ok {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			return x
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		if x, ok := constant.Uint64Val(val); ok {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			return x
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	case constant.Float:
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		if x, ok := constant.Float64Val(val); ok {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			return x
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	case constant.String:
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		return constant.StringVal(val)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	return nil
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">// A valueMap maps a case value (of a basic Go type) to a list of positions</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">// where the same case value appeared, together with the corresponding case</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// types.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// Since two case values may have the same &#34;underlying&#34; value but different</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// types we need to also check the value&#39;s types (e.g., byte(1) vs myByte(1))</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// when the switch expression is of interface type.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>type (
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	valueMap  map[any][]valueType <span class="comment">// underlying Go value -&gt; valueType</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	valueType struct {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		pos token.Pos
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		typ Type
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>func (check *Checker) caseValues(x *operand, values []ast.Expr, seen valueMap) {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>L:
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	for _, e := range values {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		var v operand
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		check.expr(nil, &amp;v, e)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		if x.mode == invalid || v.mode == invalid {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			continue L
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		check.convertUntyped(&amp;v, x.typ)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		if v.mode == invalid {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			continue L
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		<span class="comment">// Order matters: By comparing v against x, error positions are at the case values.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		res := v <span class="comment">// keep original v unchanged</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		check.comparison(&amp;res, x, token.EQL, true)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		if res.mode == invalid {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			continue L
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		if v.mode != constant_ {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			continue L <span class="comment">// we&#39;re done</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		<span class="comment">// look for duplicate values</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		if val := goVal(v.val); val != nil {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			<span class="comment">// look for duplicate types for a given value</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			<span class="comment">// (quadratic algorithm, but these lists tend to be very short)</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			for _, vt := range seen[val] {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>				if Identical(v.typ, vt.typ) {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>					check.errorf(&amp;v, DuplicateCase, &#34;duplicate case %s in expression switch&#34;, &amp;v)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>					check.error(atPos(vt.pos), DuplicateCase, &#34;\tprevious case&#34;) <span class="comment">// secondary error, \t indented</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>					continue L
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>				}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			seen[val] = append(seen[val], valueType{v.Pos(), v.typ})
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>}
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span><span class="comment">// isNil reports whether the expression e denotes the predeclared value nil.</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>func (check *Checker) isNil(e ast.Expr) bool {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	<span class="comment">// The only way to express the nil value is by literally writing nil (possibly in parentheses).</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	if name, _ := unparen(e).(*ast.Ident); name != nil {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		_, ok := check.lookup(name.Name).(*Nil)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		return ok
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	return false
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">// If the type switch expression is invalid, x is nil.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>func (check *Checker) caseTypes(x *operand, types []ast.Expr, seen map[Type]ast.Expr) (T Type) {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	var dummy operand
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>L:
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	for _, e := range types {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		<span class="comment">// The spec allows the value nil instead of a type.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		if check.isNil(e) {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			T = nil
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			check.expr(nil, &amp;dummy, e) <span class="comment">// run e through expr so we get the usual Info recordings</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		} else {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			T = check.varType(e)
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			if !isValid(T) {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>				continue L
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		<span class="comment">// look for duplicate types</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		<span class="comment">// (quadratic algorithm, but type switches tend to be reasonably small)</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		for t, other := range seen {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			if T == nil &amp;&amp; t == nil || T != nil &amp;&amp; t != nil &amp;&amp; Identical(T, t) {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				<span class="comment">// talk about &#34;case&#34; rather than &#34;type&#34; because of nil case</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>				Ts := &#34;nil&#34;
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>				if T != nil {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>					Ts = TypeString(T, check.qualifier)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>				}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>				check.errorf(e, DuplicateCase, &#34;duplicate case %s in type switch&#34;, Ts)
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>				check.error(other, DuplicateCase, &#34;\tprevious case&#34;) <span class="comment">// secondary error, \t indented</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>				continue L
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		seen[T] = e
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		if x != nil &amp;&amp; T != nil {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>			check.typeAssertion(e, x, T, true)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	return
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">// TODO(gri) Once we are certain that typeHash is correct in all situations, use this version of caseTypes instead.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">// (Currently it may be possible that different types have identical names and import paths due to ImporterFrom.)</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span><span class="comment">// func (check *Checker) caseTypes(x *operand, xtyp *Interface, types []ast.Expr, seen map[string]ast.Expr) (T Type) {</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span><span class="comment">// 	var dummy operand</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span><span class="comment">// L:</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span><span class="comment">// 	for _, e := range types {</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span><span class="comment">// 		// The spec allows the value nil instead of a type.</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span><span class="comment">// 		var hash string</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">// 		if check.isNil(e) {</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span><span class="comment">// 			check.expr(nil, &amp;dummy, e) // run e through expr so we get the usual Info recordings</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// 			T = nil</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// 			hash = &#34;&lt;nil&gt;&#34; // avoid collision with a type named nil</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">// 		} else {</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span><span class="comment">// 			T = check.varType(e)</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// 			if !isValid(T) {</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">// 				continue L</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">// 			}</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// 			hash = typeHash(T, nil)</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">// 		}</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// 		// look for duplicate types</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">// 		if other := seen[hash]; other != nil {</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// 			// talk about &#34;case&#34; rather than &#34;type&#34; because of nil case</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// 			Ts := &#34;nil&#34;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">// 			if T != nil {</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span><span class="comment">// 				Ts = TypeString(T, check.qualifier)</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">// 			}</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">// 			var err error_</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">//			err.code = DuplicateCase</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// 			err.errorf(e, &#34;duplicate case %s in type switch&#34;, Ts)</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// 			err.errorf(other, &#34;previous case&#34;)</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// 			check.report(&amp;err)</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// 			continue L</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// 		}</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// 		seen[hash] = e</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// 		if T != nil {</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// 			check.typeAssertion(e.Pos(), x, xtyp, T)</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// 		}</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// 	}</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// 	return</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// }</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// stmt typechecks statement s.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>func (check *Checker) stmt(ctxt stmtContext, s ast.Stmt) {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	<span class="comment">// statements must end with the same top scope as they started with</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	if debug {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		defer func(scope *Scope) {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			<span class="comment">// don&#39;t check if code is panicking</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			if p := recover(); p != nil {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>				panic(p)
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			assert(scope == check.scope)
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		}(check.scope)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	<span class="comment">// process collected function literals before scope changes</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	defer check.processDelayed(len(check.delayed))
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	<span class="comment">// reset context for statements of inner blocks</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	inner := ctxt &amp;^ (fallthroughOk | finalSwitchCase | inTypeSwitch)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	switch s := s.(type) {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	case *ast.BadStmt, *ast.EmptyStmt:
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		<span class="comment">// ignore</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	case *ast.DeclStmt:
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		check.declStmt(s.Decl)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	case *ast.LabeledStmt:
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		check.hasLabel = true
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		check.stmt(ctxt, s.Stmt)
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	case *ast.ExprStmt:
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;With the exception of specific built-in functions,</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		<span class="comment">// function and method calls and receive operations can appear</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		<span class="comment">// in statement context. Such statements may be parenthesized.&#34;</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		var x operand
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		kind := check.rawExpr(nil, &amp;x, s.X, nil, false)
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		var msg string
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		var code Code
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		switch x.mode {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		default:
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			if kind == statement {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>				return
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			msg = &#34;is not used&#34;
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			code = UnusedExpr
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		case builtin:
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			msg = &#34;must be called&#34;
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			code = UncalledBuiltin
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		case typexpr:
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			msg = &#34;is not an expression&#34;
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			code = NotAnExpr
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		check.errorf(&amp;x, code, &#34;%s %s&#34;, &amp;x, msg)
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	case *ast.SendStmt:
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		var ch, val operand
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		check.expr(nil, &amp;ch, s.Chan)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		check.expr(nil, &amp;val, s.Value)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		if ch.mode == invalid || val.mode == invalid {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			return
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		}
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		u := coreType(ch.typ)
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		if u == nil {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			check.errorf(inNode(s, s.Arrow), InvalidSend, invalidOp+&#34;cannot send to %s: no core type&#34;, &amp;ch)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			return
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		uch, _ := u.(*Chan)
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		if uch == nil {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>			check.errorf(inNode(s, s.Arrow), InvalidSend, invalidOp+&#34;cannot send to non-channel %s&#34;, &amp;ch)
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			return
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		if uch.dir == RecvOnly {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>			check.errorf(inNode(s, s.Arrow), InvalidSend, invalidOp+&#34;cannot send to receive-only channel %s&#34;, &amp;ch)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			return
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		check.assignment(&amp;val, uch.elem, &#34;send&#34;)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	case *ast.IncDecStmt:
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		var op token.Token
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		switch s.Tok {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		case token.INC:
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			op = token.ADD
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		case token.DEC:
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			op = token.SUB
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		default:
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			check.errorf(inNode(s, s.TokPos), InvalidSyntaxTree, &#34;unknown inc/dec operation %s&#34;, s.Tok)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			return
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		var x operand
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		check.expr(nil, &amp;x, s.X)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			return
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		if !allNumeric(x.typ) {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>			check.errorf(s.X, NonNumericIncDec, invalidOp+&#34;%s%s (non-numeric type %s)&#34;, s.X, s.Tok, x.typ)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			return
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		Y := &amp;ast.BasicLit{ValuePos: s.X.Pos(), Kind: token.INT, Value: &#34;1&#34;} <span class="comment">// use x&#39;s position</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		check.binary(&amp;x, nil, s.X, Y, op, s.TokPos)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			return
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		check.assignVar(s.X, nil, &amp;x, &#34;assignment&#34;)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	case *ast.AssignStmt:
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		switch s.Tok {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		case token.ASSIGN, token.DEFINE:
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>			if len(s.Lhs) == 0 {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>				check.error(s, InvalidSyntaxTree, &#34;missing lhs in assignment&#34;)
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>				return
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			if s.Tok == token.DEFINE {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>				check.shortVarDecl(inNode(s, s.TokPos), s.Lhs, s.Rhs)
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			} else {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>				<span class="comment">// regular assignment</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>				check.assignVars(s.Lhs, s.Rhs)
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		default:
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			<span class="comment">// assignment operations</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>				check.errorf(inNode(s, s.TokPos), MultiValAssignOp, &#34;assignment operation %s requires single-valued expressions&#34;, s.Tok)
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>				return
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			op := assignOp(s.Tok)
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			if op == token.ILLEGAL {
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>				check.errorf(atPos(s.TokPos), InvalidSyntaxTree, &#34;unknown assignment operation %s&#34;, s.Tok)
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>				return
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>			}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>			var x operand
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>			check.binary(&amp;x, nil, s.Lhs[0], s.Rhs[0], op, s.TokPos)
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>			if x.mode == invalid {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>				return
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			}
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			check.assignVar(s.Lhs[0], nil, &amp;x, &#34;assignment&#34;)
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	case *ast.GoStmt:
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		check.suspendedCall(&#34;go&#34;, s.Call)
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	case *ast.DeferStmt:
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		check.suspendedCall(&#34;defer&#34;, s.Call)
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	case *ast.ReturnStmt:
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		res := check.sig.results
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		<span class="comment">// Return with implicit results allowed for function with named results.</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		<span class="comment">// (If one is named, all are named.)</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		if len(s.Results) == 0 &amp;&amp; res.Len() &gt; 0 &amp;&amp; res.vars[0].name != &#34;&#34; {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			<span class="comment">// spec: &#34;Implementation restriction: A compiler may disallow an empty expression</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			<span class="comment">// list in a &#34;return&#34; statement if a different entity (constant, type, or variable)</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			<span class="comment">// with the same name as a result parameter is in scope at the place of the return.&#34;</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			for _, obj := range res.vars {
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>				if alt := check.lookup(obj.name); alt != nil &amp;&amp; alt != obj {
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>					check.errorf(s, OutOfScopeResult, &#34;result parameter %s not in scope at return&#34;, obj.name)
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>					check.errorf(alt, OutOfScopeResult, &#34;\tinner declaration of %s&#34;, obj)
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>					<span class="comment">// ok to continue</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>				}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			}
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		} else {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			var lhs []*Var
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			if res.Len() &gt; 0 {
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>				lhs = res.vars
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>			}
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>			check.initVars(lhs, s.Results, s)
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		}
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	case *ast.BranchStmt:
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		if s.Label != nil {
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			check.hasLabel = true
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			return <span class="comment">// checked in 2nd pass (check.labels)</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		}
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		switch s.Tok {
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		case token.BREAK:
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>			if ctxt&amp;breakOk == 0 {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>				check.error(s, MisplacedBreak, &#34;break not in for, switch, or select statement&#34;)
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		case token.CONTINUE:
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>			if ctxt&amp;continueOk == 0 {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>				check.error(s, MisplacedContinue, &#34;continue not in for statement&#34;)
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			}
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		case token.FALLTHROUGH:
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			if ctxt&amp;fallthroughOk == 0 {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>				var msg string
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>				switch {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>				case ctxt&amp;finalSwitchCase != 0:
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>					msg = &#34;cannot fallthrough final case in switch&#34;
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>				case ctxt&amp;inTypeSwitch != 0:
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>					msg = &#34;cannot fallthrough in type switch&#34;
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>				default:
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>					msg = &#34;fallthrough statement out of place&#34;
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>				}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>				check.error(s, MisplacedFallthrough, msg)
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>			}
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		default:
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			check.errorf(s, InvalidSyntaxTree, &#34;branch statement: %s&#34;, s.Tok)
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	case *ast.BlockStmt:
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		check.openScope(s, &#34;block&#34;)
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		defer check.closeScope()
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		check.stmtList(inner, s.List)
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	case *ast.IfStmt:
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		check.openScope(s, &#34;if&#34;)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		defer check.closeScope()
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		check.simpleStmt(s.Init)
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		var x operand
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		check.expr(nil, &amp;x, s.Cond)
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		if x.mode != invalid &amp;&amp; !allBoolean(x.typ) {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>			check.error(s.Cond, InvalidCond, &#34;non-boolean condition in if statement&#34;)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		check.stmt(inner, s.Body)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		<span class="comment">// The parser produces a correct AST but if it was modified</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		<span class="comment">// elsewhere the else branch may be invalid. Check again.</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		switch s.Else.(type) {
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		case nil, *ast.BadStmt:
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>			<span class="comment">// valid or error already reported</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		case *ast.IfStmt, *ast.BlockStmt:
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>			check.stmt(inner, s.Else)
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		default:
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>			check.error(s.Else, InvalidSyntaxTree, &#34;invalid else branch in if statement&#34;)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	case *ast.SwitchStmt:
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		inner |= breakOk
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		check.openScope(s, &#34;switch&#34;)
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		defer check.closeScope()
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		check.simpleStmt(s.Init)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		var x operand
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		if s.Tag != nil {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			check.expr(nil, &amp;x, s.Tag)
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>			<span class="comment">// By checking assignment of x to an invisible temporary</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>			<span class="comment">// (as a compiler would), we get all the relevant checks.</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>			check.assignment(&amp;x, nil, &#34;switch expression&#34;)
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>			if x.mode != invalid &amp;&amp; !Comparable(x.typ) &amp;&amp; !hasNil(x.typ) {
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>				check.errorf(&amp;x, InvalidExprSwitch, &#34;cannot switch on %s (%s is not comparable)&#34;, &amp;x, x.typ)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>				x.mode = invalid
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>			}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		} else {
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>			<span class="comment">// spec: &#34;A missing switch expression is</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>			<span class="comment">// equivalent to the boolean value true.&#34;</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			x.mode = constant_
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>			x.typ = Typ[Bool]
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>			x.val = constant.MakeBool(true)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>			x.expr = &amp;ast.Ident{NamePos: s.Body.Lbrace, Name: &#34;true&#34;}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		}
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		check.multipleDefaults(s.Body.List)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		seen := make(valueMap) <span class="comment">// map of seen case values to positions and types</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		for i, c := range s.Body.List {
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>			clause, _ := c.(*ast.CaseClause)
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>			if clause == nil {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>				check.error(c, InvalidSyntaxTree, &#34;incorrect expression switch case&#34;)
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>				continue
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>			}
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			check.caseValues(&amp;x, clause.List, seen)
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			check.openScope(clause, &#34;case&#34;)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>			inner := inner
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>			if i+1 &lt; len(s.Body.List) {
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>				inner |= fallthroughOk
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>			} else {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>				inner |= finalSwitchCase
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			check.stmtList(inner, clause.Body)
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>			check.closeScope()
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	case *ast.TypeSwitchStmt:
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>		inner |= breakOk | inTypeSwitch
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		check.openScope(s, &#34;type switch&#34;)
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		defer check.closeScope()
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		check.simpleStmt(s.Init)
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		<span class="comment">// A type switch guard must be of the form:</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		<span class="comment">//     TypeSwitchGuard = [ identifier &#34;:=&#34; ] PrimaryExpr &#34;.&#34; &#34;(&#34; &#34;type&#34; &#34;)&#34; .</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		<span class="comment">// The parser is checking syntactic correctness;</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		<span class="comment">// remaining syntactic errors are considered AST errors here.</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) better factoring of error handling (invalid ASTs)</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		var lhs *ast.Ident <span class="comment">// lhs identifier or nil</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		var rhs ast.Expr
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		switch guard := s.Assign.(type) {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		case *ast.ExprStmt:
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>			rhs = guard.X
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		case *ast.AssignStmt:
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>			if len(guard.Lhs) != 1 || guard.Tok != token.DEFINE || len(guard.Rhs) != 1 {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>				check.error(s, InvalidSyntaxTree, &#34;incorrect form of type switch guard&#34;)
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>				return
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>			}
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>			lhs, _ = guard.Lhs[0].(*ast.Ident)
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>			if lhs == nil {
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>				check.error(s, InvalidSyntaxTree, &#34;incorrect form of type switch guard&#34;)
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>				return
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>			}
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>			if lhs.Name == &#34;_&#34; {
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>				<span class="comment">// _ := x.(type) is an invalid short variable declaration</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>				check.softErrorf(lhs, NoNewVar, &#34;no new variable on left side of :=&#34;)
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>				lhs = nil <span class="comment">// avoid declared and not used error below</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>			} else {
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>				check.recordDef(lhs, nil) <span class="comment">// lhs variable is implicitly declared in each cause clause</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>			}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>			rhs = guard.Rhs[0]
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		default:
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>			check.error(s, InvalidSyntaxTree, &#34;incorrect form of type switch guard&#34;)
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>			return
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		}
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		<span class="comment">// rhs must be of the form: expr.(type) and expr must be an ordinary interface</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		expr, _ := rhs.(*ast.TypeAssertExpr)
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		if expr == nil || expr.Type != nil {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			check.error(s, InvalidSyntaxTree, &#34;incorrect form of type switch guard&#34;)
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>			return
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		var x operand
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		check.expr(nil, &amp;x, expr.X)
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>			return
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		}
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) we may want to permit type switches on type parameter values at some point</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		var sx *operand <span class="comment">// switch expression against which cases are compared against; nil if invalid</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		if isTypeParam(x.typ) {
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>			check.errorf(&amp;x, InvalidTypeSwitch, &#34;cannot use type switch on type parameter value %s&#34;, &amp;x)
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		} else {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>			if _, ok := under(x.typ).(*Interface); ok {
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>				sx = &amp;x
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>			} else {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>				check.errorf(&amp;x, InvalidTypeSwitch, &#34;%s is not an interface&#34;, &amp;x)
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			}
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>		}
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		check.multipleDefaults(s.Body.List)
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		var lhsVars []*Var              <span class="comment">// list of implicitly declared lhs variables</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		seen := make(map[Type]ast.Expr) <span class="comment">// map of seen types to positions</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		for _, s := range s.Body.List {
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>			clause, _ := s.(*ast.CaseClause)
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			if clause == nil {
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>				check.error(s, InvalidSyntaxTree, &#34;incorrect type switch case&#34;)
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>				continue
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>			<span class="comment">// Check each type in this type switch case.</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>			T := check.caseTypes(sx, clause.List, seen)
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>			check.openScope(clause, &#34;case&#34;)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>			<span class="comment">// If lhs exists, declare a corresponding variable in the case-local scope.</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>			if lhs != nil {
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>				<span class="comment">// spec: &#34;The TypeSwitchGuard may include a short variable declaration.</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>				<span class="comment">// When that form is used, the variable is declared at the beginning of</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>				<span class="comment">// the implicit block in each clause. In clauses with a case listing</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>				<span class="comment">// exactly one type, the variable has that type; otherwise, the variable</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>				<span class="comment">// has the type of the expression in the TypeSwitchGuard.&#34;</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>				if len(clause.List) != 1 || T == nil {
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>					T = x.typ
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>				}
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>				obj := NewVar(lhs.Pos(), check.pkg, lhs.Name, T)
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>				scopePos := clause.Pos() + token.Pos(len(&#34;default&#34;)) <span class="comment">// for default clause (len(List) == 0)</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>				if n := len(clause.List); n &gt; 0 {
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>					scopePos = clause.List[n-1].End()
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>				}
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>				check.declare(check.scope, nil, obj, scopePos)
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>				check.recordImplicit(clause, obj)
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>				<span class="comment">// For the &#34;declared and not used&#34; error, all lhs variables act as</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>				<span class="comment">// one; i.e., if any one of them is &#39;used&#39;, all of them are &#39;used&#39;.</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>				<span class="comment">// Collect them for later analysis.</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>				lhsVars = append(lhsVars, obj)
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>			}
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			check.stmtList(inner, clause.Body)
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>			check.closeScope()
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		<span class="comment">// If lhs exists, we must have at least one lhs variable that was used.</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		if lhs != nil {
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>			var used bool
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>			for _, v := range lhsVars {
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>				if v.used {
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>					used = true
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>				}
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>				v.used = true <span class="comment">// avoid usage error when checking entire function</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>			}
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>			if !used {
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>				check.softErrorf(lhs, UnusedVar, &#34;%s declared and not used&#34;, lhs.Name)
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>			}
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>		}
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	case *ast.SelectStmt:
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		inner |= breakOk
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		check.multipleDefaults(s.Body.List)
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		for _, s := range s.Body.List {
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>			clause, _ := s.(*ast.CommClause)
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>			if clause == nil {
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>				continue <span class="comment">// error reported before</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>			<span class="comment">// clause.Comm must be a SendStmt, RecvStmt, or default case</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>			valid := false
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>			var rhs ast.Expr <span class="comment">// rhs of RecvStmt, or nil</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>			switch s := clause.Comm.(type) {
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>			case nil, *ast.SendStmt:
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>				valid = true
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>			case *ast.AssignStmt:
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>				if len(s.Rhs) == 1 {
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>					rhs = s.Rhs[0]
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>				}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>			case *ast.ExprStmt:
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>				rhs = s.X
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>			}
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>			<span class="comment">// if present, rhs must be a receive operation</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>			if rhs != nil {
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>				if x, _ := unparen(rhs).(*ast.UnaryExpr); x != nil &amp;&amp; x.Op == token.ARROW {
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>					valid = true
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>				}
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>			}
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>			if !valid {
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>				check.error(clause.Comm, InvalidSelectCase, &#34;select case must be send or receive (possibly with assignment)&#34;)
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>				continue
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>			}
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>			check.openScope(s, &#34;case&#34;)
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>			if clause.Comm != nil {
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>				check.stmt(inner, clause.Comm)
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>			}
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>			check.stmtList(inner, clause.Body)
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>			check.closeScope()
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		}
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	case *ast.ForStmt:
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		inner |= breakOk | continueOk
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		check.openScope(s, &#34;for&#34;)
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		defer check.closeScope()
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>		check.simpleStmt(s.Init)
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		if s.Cond != nil {
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>			var x operand
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>			check.expr(nil, &amp;x, s.Cond)
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>			if x.mode != invalid &amp;&amp; !allBoolean(x.typ) {
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>				check.error(s.Cond, InvalidCond, &#34;non-boolean condition in for statement&#34;)
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>			}
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>		}
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		check.simpleStmt(s.Post)
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;The init statement may be a short variable</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>		<span class="comment">// declaration, but the post statement must not.&#34;</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		if s, _ := s.Post.(*ast.AssignStmt); s != nil &amp;&amp; s.Tok == token.DEFINE {
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>			check.softErrorf(s, InvalidPostDecl, &#34;cannot declare in post statement&#34;)
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>			<span class="comment">// Don&#39;t call useLHS here because we want to use the lhs in</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>			<span class="comment">// this erroneous statement so that we don&#39;t get errors about</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>			<span class="comment">// these lhs variables being declared and not used.</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>			check.use(s.Lhs...) <span class="comment">// avoid follow-up errors</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>		}
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		check.stmt(inner, s.Body)
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	case *ast.RangeStmt:
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		inner |= breakOk | continueOk
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		check.rangeStmt(inner, s)
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	default:
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		check.error(s, InvalidSyntaxTree, &#34;invalid statement&#34;)
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	}
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>}
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>func (check *Checker) rangeStmt(inner stmtContext, s *ast.RangeStmt) {
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	<span class="comment">// Convert go/ast form to local variables.</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>	type Expr = ast.Expr
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	type identType = ast.Ident
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	identName := func(n *identType) string { return n.Name }
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	sKey, sValue := s.Key, s.Value
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	var sExtra ast.Expr = nil
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	isDef := s.Tok == token.DEFINE
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	rangeVar := s.X
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	noNewVarPos := inNode(s, s.TokPos)
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	<span class="comment">// Everything from here on is shared between cmd/compile/internal/types2 and go/types.</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	<span class="comment">// check expression to iterate over</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	var x operand
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	check.expr(nil, &amp;x, rangeVar)
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	<span class="comment">// determine key/value types</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	var key, val Type
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	if x.mode != invalid {
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		<span class="comment">// Ranging over a type parameter is permitted if it has a core type.</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>		k, v, cause, isFunc, ok := rangeKeyVal(x.typ, func(v goVersion) bool {
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>			return check.allowVersion(check.pkg, x.expr, v)
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>		})
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		switch {
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		case !ok &amp;&amp; cause != &#34;&#34;:
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>			check.softErrorf(&amp;x, InvalidRangeExpr, &#34;cannot range over %s: %s&#34;, &amp;x, cause)
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		case !ok:
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>			check.softErrorf(&amp;x, InvalidRangeExpr, &#34;cannot range over %s&#34;, &amp;x)
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		case k == nil &amp;&amp; sKey != nil:
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>			check.softErrorf(sKey, InvalidIterVar, &#34;range over %s permits no iteration variables&#34;, &amp;x)
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		case v == nil &amp;&amp; sValue != nil:
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>			check.softErrorf(sValue, InvalidIterVar, &#34;range over %s permits only one iteration variable&#34;, &amp;x)
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		case sExtra != nil:
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>			check.softErrorf(sExtra, InvalidIterVar, &#34;range clause permits at most two iteration variables&#34;)
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>		case isFunc &amp;&amp; ((k == nil) != (sKey == nil) || (v == nil) != (sValue == nil)):
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>			var count string
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>			switch {
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>			case k == nil:
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>				count = &#34;no iteration variables&#34;
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>			case v == nil:
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>				count = &#34;one iteration variable&#34;
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>			default:
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>				count = &#34;two iteration variables&#34;
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>			}
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>			check.softErrorf(&amp;x, InvalidIterVar, &#34;range over %s must have %s&#34;, &amp;x, count)
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>		}
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>		key, val = k, v
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	}
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	<span class="comment">// Open the for-statement block scope now, after the range clause.</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	<span class="comment">// Iteration variables declared with := need to go in this scope (was go.dev/issue/51437).</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	check.openScope(s, &#34;range&#34;)
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	defer check.closeScope()
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>	<span class="comment">// check assignment to/declaration of iteration variables</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>	<span class="comment">// (irregular assignment, cannot easily map to existing assignment checks)</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	<span class="comment">// lhs expressions and initialization value (rhs) types</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	lhs := [2]Expr{sKey, sValue} <span class="comment">// sKey, sValue may be nil</span>
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	rhs := [2]Type{key, val}     <span class="comment">// key, val may be nil</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	constIntRange := x.mode == constant_ &amp;&amp; isInteger(x.typ)
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	if isDef {
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>		<span class="comment">// short variable declaration</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		var vars []*Var
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		for i, lhs := range lhs {
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>			if lhs == nil {
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>				continue
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>			}
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>			<span class="comment">// determine lhs variable</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>			var obj *Var
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>			if ident, _ := lhs.(*identType); ident != nil {
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>				<span class="comment">// declare new variable</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>				name := identName(ident)
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>				obj = NewVar(ident.Pos(), check.pkg, name, nil)
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>				check.recordDef(ident, obj)
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>				<span class="comment">// _ variables don&#39;t count as new variables</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>				if name != &#34;_&#34; {
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>					vars = append(vars, obj)
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>				}
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>			} else {
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>				check.errorf(lhs, InvalidSyntaxTree, &#34;cannot declare %s&#34;, lhs)
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>				obj = NewVar(lhs.Pos(), check.pkg, &#34;_&#34;, nil) <span class="comment">// dummy variable</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>			}
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>			<span class="comment">// initialize lhs variable</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>			if constIntRange {
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>				check.initVar(obj, &amp;x, &#34;range clause&#34;)
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>			} else if typ := rhs[i]; typ != nil {
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>				x.mode = value
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>				x.expr = lhs <span class="comment">// we don&#39;t have a better rhs expression to use here</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>				x.typ = typ
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>				check.initVar(obj, &amp;x, &#34;assignment&#34;) <span class="comment">// error is on variable, use &#34;assignment&#34; not &#34;range clause&#34;</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>			} else {
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>				obj.typ = Typ[Invalid]
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>				obj.used = true <span class="comment">// don&#39;t complain about unused variable</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>			}
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>		}
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>		<span class="comment">// declare variables</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>		if len(vars) &gt; 0 {
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>			scopePos := s.Body.Pos()
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>			for _, obj := range vars {
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>				check.declare(check.scope, nil <span class="comment">/* recordDef already called */</span>, obj, scopePos)
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>			}
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		} else {
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>			check.error(noNewVarPos, NoNewVar, &#34;no new variables on left side of :=&#34;)
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>		}
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>	} else if sKey != nil <span class="comment">/* lhs[0] != nil */</span> {
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>		<span class="comment">// ordinary assignment</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>		for i, lhs := range lhs {
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>			if lhs == nil {
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>				continue
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>			}
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			if constIntRange {
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>				check.assignVar(lhs, nil, &amp;x, &#34;range clause&#34;)
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>			} else if typ := rhs[i]; typ != nil {
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>				x.mode = value
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>				x.expr = lhs <span class="comment">// we don&#39;t have a better rhs expression to use here</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>				x.typ = typ
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>				check.assignVar(lhs, nil, &amp;x, &#34;assignment&#34;) <span class="comment">// error is on variable, use &#34;assignment&#34; not &#34;range clause&#34;</span>
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>			}
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>		}
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	} else if constIntRange {
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>		<span class="comment">// If we don&#39;t have any iteration variables, we still need to</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>		<span class="comment">// check that a (possibly untyped) integer range expression x</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>		<span class="comment">// is valid.</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>		<span class="comment">// We do this by checking the assignment _ = x. This ensures</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		<span class="comment">// that an untyped x can be converted to a value of type int.</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>		check.assignment(&amp;x, nil, &#34;range clause&#34;)
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>	}
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	check.stmt(inner, s.Body)
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>}
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span><span class="comment">// rangeKeyVal returns the key and value type produced by a range clause</span>
<span id="L974" class="ln">   974&nbsp;&nbsp;</span><span class="comment">// over an expression of type typ.</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span><span class="comment">// If allowVersion != nil, it is used to check the required language version.</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span><span class="comment">// If the range clause is not permitted, rangeKeyVal returns ok = false.</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span><span class="comment">// When ok = false, rangeKeyVal may also return a reason in cause.</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>func rangeKeyVal(typ Type, allowVersion func(goVersion) bool) (key, val Type, cause string, isFunc, ok bool) {
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>	bad := func(cause string) (Type, Type, string, bool, bool) {
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>		return Typ[Invalid], Typ[Invalid], cause, false, false
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	}
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	toSig := func(t Type) *Signature {
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		sig, _ := coreType(t).(*Signature)
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>		return sig
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	}
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	orig := typ
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>	switch typ := arrayPtrDeref(coreType(typ)).(type) {
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	case nil:
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>		return bad(&#34;no core type&#34;)
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>	case *Basic:
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>		if isString(typ) {
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>			return Typ[Int], universeRune, &#34;&#34;, false, true <span class="comment">// use &#39;rune&#39; name</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>		}
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		if isInteger(typ) {
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>			if allowVersion != nil &amp;&amp; !allowVersion(go1_22) {
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>				return bad(&#34;requires go1.22 or later&#34;)
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>			}
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>			return orig, nil, &#34;&#34;, false, true
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>		}
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>	case *Array:
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>		return Typ[Int], typ.elem, &#34;&#34;, false, true
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>	case *Slice:
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>		return Typ[Int], typ.elem, &#34;&#34;, false, true
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	case *Map:
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>		return typ.key, typ.elem, &#34;&#34;, false, true
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	case *Chan:
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>		if typ.dir == SendOnly {
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>			return bad(&#34;receive from send-only channel&#34;)
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>		}
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>		return typ.elem, nil, &#34;&#34;, false, true
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	case *Signature:
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) when this becomes enabled permanently, add version check</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>		if !buildcfg.Experiment.RangeFunc {
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>			break
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>		}
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>		assert(typ.Recv() == nil)
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>		switch {
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>		case typ.Params().Len() != 1:
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>			return bad(&#34;func must be func(yield func(...) bool): wrong argument count&#34;)
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		case toSig(typ.Params().At(0).Type()) == nil:
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>			return bad(&#34;func must be func(yield func(...) bool): argument is not func&#34;)
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		case typ.Results().Len() != 0:
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>			return bad(&#34;func must be func(yield func(...) bool): unexpected results&#34;)
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		}
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>		cb := toSig(typ.Params().At(0).Type())
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>		assert(cb.Recv() == nil)
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>		switch {
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>		case cb.Params().Len() &gt; 2:
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>			return bad(&#34;func must be func(yield func(...) bool): yield func has too many parameters&#34;)
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>		case cb.Results().Len() != 1 || !isBoolean(cb.Results().At(0).Type()):
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>			return bad(&#34;func must be func(yield func(...) bool): yield func does not return bool&#34;)
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		}
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>		if cb.Params().Len() &gt;= 1 {
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>			key = cb.Params().At(0).Type()
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>		}
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>		if cb.Params().Len() &gt;= 2 {
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>			val = cb.Params().At(1).Type()
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>		}
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>		return key, val, &#34;&#34;, true, true
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>	}
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>	return
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>}
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>
</pre><p><a href="stmt.go?m=text">View as plain text</a></p>

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
