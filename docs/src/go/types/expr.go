<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/expr.go - Go Documentation Server</title>

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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">expr.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements typechecking of expressions.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package types
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;go/constant&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;go/internal/typeparams&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">/*
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>Basic algorithm:
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>Expressions are checked recursively, top down. Expression checker functions
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>are generally of the form:
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>  func f(x *operand, e *ast.Expr, ...)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>where e is the expression to be checked, and x is the result of the check.
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>The check performed by f may fail in which case x.mode == invalid, and
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>related error messages will have been issued by f.
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>If a hint argument is present, it is the composite literal element type
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>of an outer composite literal; it is used to type-check composite literal
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>elements that have no explicit type specification in the source
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>(e.g.: []T{{...}, {...}}, the hint is the type T in this case).
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>All expressions are checked via rawExpr, which dispatches according
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>to expression kind. Upon returning, rawExpr is recording the types and
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>constant values for all expressions that have an untyped type (those types
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>may change on the way up in the expression tree). Usually these are constants,
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>but the results of comparisons or non-constant shifts of untyped constants
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>may also be untyped, but not constant.
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>Untyped expressions may eventually become fully typed (i.e., not untyped),
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>typically when the value is assigned to a variable, or is used otherwise.
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>The updateExprType method is used to record this final type and update
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>the recorded types: the type-checked expression tree is again traversed down,
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>and the new type is propagated as needed. Untyped constant expression values
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>that become fully typed must now be representable by the full type (constant
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>sub-expression trees are left alone except for their roots). This mechanism
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>ensures that a client sees the actual (run-time) type an untyped value would
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>have. It also permits type-checking of lhs shift operands &#34;as if the shift
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>were not present&#34;: when updateExprType visits an untyped lhs shift operand
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>and assigns it it&#39;s final type, that type must be an integer type, and a
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>constant lhs must be representable as an integer.
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>When an expression gets its final type, either on the way out from rawExpr,
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>on the way down in updateExprType, or at the end of the type checker run,
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>the type (and constant value, if any) is recorded via Info.Types, if present.
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>*/</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>type opPredicates map[token.Token]func(Type) bool
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>var unaryOpPredicates opPredicates
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>func init() {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// Setting unaryOpPredicates in init avoids declaration cycles.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	unaryOpPredicates = opPredicates{
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		token.ADD: allNumeric,
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		token.SUB: allNumeric,
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		token.XOR: allInteger,
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		token.NOT: allBoolean,
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>func (check *Checker) op(m opPredicates, x *operand, op token.Token) bool {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	if pred := m[op]; pred != nil {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		if !pred(x.typ) {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			check.errorf(x, UndefinedOp, invalidOp+&#34;operator %s not defined on %s&#34;, op, x)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			return false
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	} else {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		check.errorf(x, InvalidSyntaxTree, &#34;unknown operator %s&#34;, op)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		return false
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	return true
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// opName returns the name of the operation if x is an operation</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// that might overflow; otherwise it returns the empty string.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>func opName(e ast.Expr) string {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	switch e := e.(type) {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		if int(e.Op) &lt; len(op2str2) {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			return op2str2[e.Op]
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	case *ast.UnaryExpr:
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		if int(e.Op) &lt; len(op2str1) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			return op2str1[e.Op]
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	return &#34;&#34;
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>var op2str1 = [...]string{
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	token.XOR: &#34;bitwise complement&#34;,
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// This is only used for operations that may cause overflow.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>var op2str2 = [...]string{
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	token.ADD: &#34;addition&#34;,
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	token.SUB: &#34;subtraction&#34;,
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	token.XOR: &#34;bitwise XOR&#34;,
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	token.MUL: &#34;multiplication&#34;,
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	token.SHL: &#34;shift&#34;,
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// If typ is a type parameter, underIs returns the result of typ.underIs(f).</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// Otherwise, underIs returns the result of f(under(typ)).</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>func underIs(typ Type, f func(Type) bool) bool {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	if tpar, _ := typ.(*TypeParam); tpar != nil {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		return tpar.underIs(f)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	return f(under(typ))
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// The unary expression e may be nil. It&#39;s passed in for better error messages only.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>func (check *Checker) unary(x *operand, e *ast.UnaryExpr) {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	check.expr(nil, x, e.X)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	if x.mode == invalid {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		return
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	op := e.Op
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	switch op {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	case token.AND:
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;As an exception to the addressability</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// requirement x may also be a composite literal.&#34;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		if _, ok := unparen(e.X).(*ast.CompositeLit); !ok &amp;&amp; x.mode != variable {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			check.errorf(x, UnaddressableOperand, invalidOp+&#34;cannot take address of %s&#34;, x)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			return
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		x.mode = value
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		x.typ = &amp;Pointer{base: x.typ}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		return
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	case token.ARROW:
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		u := coreType(x.typ)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		if u == nil {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			check.errorf(x, InvalidReceive, invalidOp+&#34;cannot receive from %s (no core type)&#34;, x)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			return
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		ch, _ := u.(*Chan)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		if ch == nil {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			check.errorf(x, InvalidReceive, invalidOp+&#34;cannot receive from non-channel %s&#34;, x)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			return
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		if ch.dir == SendOnly {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			check.errorf(x, InvalidReceive, invalidOp+&#34;cannot receive from send-only channel %s&#34;, x)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			return
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		x.mode = commaok
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		x.typ = ch.elem
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		check.hasCallOrRecv = true
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		return
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	case token.TILDE:
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		<span class="comment">// Provide a better error position and message than what check.op below would do.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		if !allInteger(x.typ) {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			check.error(e, UndefinedOp, &#34;cannot use ~ outside of interface or type constraint&#34;)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			return
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		check.error(e, UndefinedOp, &#34;cannot use ~ outside of interface or type constraint (use ^ for bitwise complement)&#34;)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		op = token.XOR
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	if !check.op(unaryOpPredicates, x, op) {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		return
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	if x.mode == constant_ {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		if x.val.Kind() == constant.Unknown {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			<span class="comment">// nothing to do (and don&#39;t cause an error below in the overflow check)</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			return
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		var prec uint
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		if isUnsigned(x.typ) {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			prec = uint(check.conf.sizeof(x.typ) * 8)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		x.val = constant.UnaryOp(op, x.val, prec)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		x.expr = e
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		check.overflow(x, x.Pos())
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		return
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	x.mode = value
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// x.typ remains unchanged</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>func isShift(op token.Token) bool {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	return op == token.SHL || op == token.SHR
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>func isComparison(op token.Token) bool {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">// Note: tokens are not ordered well to make this much easier</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	switch op {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		return true
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	return false
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// updateExprType updates the type of x to typ and invokes itself</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">// recursively for the operands of x, depending on expression kind.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">// If typ is still an untyped and not the final type, updateExprType</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// only updates the recorded untyped type for x and possibly its</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// operands. Otherwise (i.e., typ is not an untyped type anymore,</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// or it is the final type for x), the type and value are recorded.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// Also, if x is a constant, it must be representable as a value of typ,</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// and if x is the (formerly untyped) lhs operand of a non-constant</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// shift, it must be an integer value.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>func (check *Checker) updateExprType(x ast.Expr, typ Type, final bool) {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	check.updateExprType0(nil, x, typ, final)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>func (check *Checker) updateExprType0(parent, x ast.Expr, typ Type, final bool) {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	old, found := check.untyped[x]
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	if !found {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		return <span class="comment">// nothing to do</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// update operands of x if necessary</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	switch x := x.(type) {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	case *ast.BadExpr,
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		*ast.FuncLit,
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		*ast.CompositeLit,
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		*ast.IndexExpr,
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		*ast.SliceExpr,
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		*ast.TypeAssertExpr,
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		*ast.StarExpr,
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		*ast.KeyValueExpr,
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		*ast.ArrayType,
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		*ast.StructType,
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		*ast.FuncType,
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		*ast.InterfaceType,
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		*ast.MapType,
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		*ast.ChanType:
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		<span class="comment">// These expression are never untyped - nothing to do.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		<span class="comment">// The respective sub-expressions got their final types</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		<span class="comment">// upon assignment or use.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		if debug {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			check.dump(&#34;%v: found old type(%s): %s (new: %s)&#34;, x.Pos(), x, old.typ, typ)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			unreachable()
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		return
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	case *ast.CallExpr:
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		<span class="comment">// Resulting in an untyped constant (e.g., built-in complex).</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		<span class="comment">// The respective calls take care of calling updateExprType</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		<span class="comment">// for the arguments if necessary.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	case *ast.Ident, *ast.BasicLit, *ast.SelectorExpr:
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		<span class="comment">// An identifier denoting a constant, a constant literal,</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		<span class="comment">// or a qualified identifier (imported untyped constant).</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		<span class="comment">// No operands to take care of.</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	case *ast.ParenExpr:
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		check.updateExprType0(x, x.X, typ, final)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	case *ast.UnaryExpr:
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		<span class="comment">// If x is a constant, the operands were constants.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		<span class="comment">// The operands don&#39;t need to be updated since they</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		<span class="comment">// never get &#34;materialized&#34; into a typed value. If</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		<span class="comment">// left in the untyped map, they will be processed</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		<span class="comment">// at the end of the type check.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		if old.val != nil {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			break
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		check.updateExprType0(x, x.X, typ, final)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		if old.val != nil {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			break <span class="comment">// see comment for unary expressions</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		if isComparison(x.Op) {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			<span class="comment">// The result type is independent of operand types</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			<span class="comment">// and the operand types must have final types.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		} else if isShift(x.Op) {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			<span class="comment">// The result type depends only on lhs operand.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			<span class="comment">// The rhs type was updated when checking the shift.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			check.updateExprType0(x, x.X, typ, final)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		} else {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			<span class="comment">// The operand types match the result type.</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			check.updateExprType0(x, x.X, typ, final)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			check.updateExprType0(x, x.Y, typ, final)
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	default:
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		unreachable()
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	<span class="comment">// If the new type is not final and still untyped, just</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	<span class="comment">// update the recorded type.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	if !final &amp;&amp; isUntyped(typ) {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		old.typ = under(typ).(*Basic)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		check.untyped[x] = old
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		return
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	<span class="comment">// Otherwise we have the final (typed or untyped type).</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// Remove it from the map of yet untyped expressions.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	delete(check.untyped, x)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	if old.isLhs {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		<span class="comment">// If x is the lhs of a shift, its final type must be integer.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		<span class="comment">// We already know from the shift check that it is representable</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		<span class="comment">// as an integer if it is a constant.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		if !allInteger(typ) {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			check.errorf(x, InvalidShiftOperand, invalidOp+&#34;shifted operand %s (type %s) must be integer&#34;, x, typ)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			return
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		<span class="comment">// Even if we have an integer, if the value is a constant we</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		<span class="comment">// still must check that it is representable as the specific</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		<span class="comment">// int type requested (was go.dev/issue/22969). Fall through here.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	if old.val != nil {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		<span class="comment">// If x is a constant, it must be representable as a value of typ.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		c := operand{old.mode, x, old.typ, old.val, 0}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		check.convertUntyped(&amp;c, typ)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		if c.mode == invalid {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			return
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">// Everything&#39;s fine, record final type and value for x.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	check.recordTypeAndValue(x, old.mode, typ, old.val)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">// updateExprVal updates the value of x to val.</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>func (check *Checker) updateExprVal(x ast.Expr, val constant.Value) {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	if info, ok := check.untyped[x]; ok {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		info.val = val
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		check.untyped[x] = info
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// implicitTypeAndValue returns the implicit type of x when used in a context</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// where the target type is expected. If no such implicit conversion is</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// possible, it returns a nil Type and non-zero error code.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// If x is a constant operand, the returned constant.Value will be the</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// representation of x in this context.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>func (check *Checker) implicitTypeAndValue(x *operand, target Type) (Type, constant.Value, Code) {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	if x.mode == invalid || isTyped(x.typ) || !isValid(target) {
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		return x.typ, nil, 0
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	<span class="comment">// x is untyped</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	if isUntyped(target) {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		<span class="comment">// both x and target are untyped</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		if m := maxType(x.typ, target); m != nil {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			return m, nil, 0
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		return nil, nil, InvalidUntypedConversion
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	switch u := under(target).(type) {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	case *Basic:
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		if x.mode == constant_ {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			v, code := check.representation(x, u)
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			if code != 0 {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>				return nil, nil, code
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			return target, v, code
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		<span class="comment">// Non-constant untyped values may appear as the</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		<span class="comment">// result of comparisons (untyped bool), intermediate</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		<span class="comment">// (delayed-checked) rhs operands of shifts, and as</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		<span class="comment">// the value nil.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		switch x.typ.(*Basic).kind {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		case UntypedBool:
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>			if !isBoolean(target) {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>				return nil, nil, InvalidUntypedConversion
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		case UntypedInt, UntypedRune, UntypedFloat, UntypedComplex:
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			if !isNumeric(target) {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>				return nil, nil, InvalidUntypedConversion
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		case UntypedString:
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			<span class="comment">// Non-constant untyped string values are not permitted by the spec and</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			<span class="comment">// should not occur during normal typechecking passes, but this path is</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			<span class="comment">// reachable via the AssignableTo API.</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			if !isString(target) {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>				return nil, nil, InvalidUntypedConversion
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		case UntypedNil:
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			<span class="comment">// Unsafe.Pointer is a basic type that includes nil.</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			if !hasNil(target) {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>				return nil, nil, InvalidUntypedConversion
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			<span class="comment">// Preserve the type of nil as UntypedNil: see go.dev/issue/13061.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			return Typ[UntypedNil], nil, 0
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		default:
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			return nil, nil, InvalidUntypedConversion
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	case *Interface:
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		if isTypeParam(target) {
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			if !u.typeSet().underIs(func(u Type) bool {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>				if u == nil {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>					return false
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>				}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>				t, _, _ := check.implicitTypeAndValue(x, u)
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>				return t != nil
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			}) {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>				return nil, nil, InvalidUntypedConversion
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>			}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			<span class="comment">// keep nil untyped (was bug go.dev/issue/39755)</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			if x.isNil() {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>				return Typ[UntypedNil], nil, 0
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>			}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			break
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		<span class="comment">// Values must have concrete dynamic types. If the value is nil,</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		<span class="comment">// keep it untyped (this is important for tools such as go vet which</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		<span class="comment">// need the dynamic type for argument checking of say, print</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		<span class="comment">// functions)</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		if x.isNil() {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>			return Typ[UntypedNil], nil, 0
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		<span class="comment">// cannot assign untyped values to non-empty interfaces</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		if !u.Empty() {
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			return nil, nil, InvalidUntypedConversion
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		return Default(x.typ), nil, 0
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	case *Pointer, *Signature, *Slice, *Map, *Chan:
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		if !x.isNil() {
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			return nil, nil, InvalidUntypedConversion
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		<span class="comment">// Keep nil untyped - see comment for interfaces, above.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		return Typ[UntypedNil], nil, 0
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	default:
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		return nil, nil, InvalidUntypedConversion
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	return target, nil, 0
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span><span class="comment">// If switchCase is true, the operator op is ignored.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>func (check *Checker) comparison(x, y *operand, op token.Token, switchCase bool) {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	<span class="comment">// Avoid spurious errors if any of the operands has an invalid type (go.dev/issue/54405).</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	if !isValid(x.typ) || !isValid(y.typ) {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		return
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	if switchCase {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		op = token.EQL
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	errOp := x  <span class="comment">// operand for which error is reported, if any</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	cause := &#34;&#34; <span class="comment">// specific error cause, if any</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;In any comparison, the first operand must be assignable</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	<span class="comment">// to the type of the second operand, or vice versa.&#34;</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	code := MismatchedTypes
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	ok, _ := x.assignableTo(check, y.typ, nil)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	if !ok {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		ok, _ = y.assignableTo(check, x.typ, nil)
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	if !ok {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		<span class="comment">// Report the error on the 2nd operand since we only</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		<span class="comment">// know after seeing the 2nd operand whether we have</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		<span class="comment">// a type mismatch.</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		errOp = y
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		cause = check.sprintf(&#34;mismatched types %s and %s&#34;, x.typ, y.typ)
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		goto Error
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	<span class="comment">// check if comparison is defined for operands</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	code = UndefinedOp
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	switch op {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	case token.EQL, token.NEQ:
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;The equality operators == and != apply to operands that are comparable.&#34;</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		switch {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		case x.isNil() || y.isNil():
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>			<span class="comment">// Comparison against nil requires that the other operand type has nil.</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>			typ := x.typ
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>			if x.isNil() {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>				typ = y.typ
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			}
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			if !hasNil(typ) {
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>				<span class="comment">// This case should only be possible for &#34;nil == nil&#34;.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>				<span class="comment">// Report the error on the 2nd operand since we only</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>				<span class="comment">// know after seeing the 2nd operand whether we have</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>				<span class="comment">// an invalid comparison.</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>				errOp = y
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>				goto Error
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		case !Comparable(x.typ):
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>			errOp = x
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			cause = check.incomparableCause(x.typ)
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>			goto Error
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		case !Comparable(y.typ):
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			errOp = y
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			cause = check.incomparableCause(y.typ)
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			goto Error
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	case token.LSS, token.LEQ, token.GTR, token.GEQ:
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		<span class="comment">// spec: The ordering operators &lt;, &lt;=, &gt;, and &gt;= apply to operands that are ordered.&#34;</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		switch {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		case !allOrdered(x.typ):
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>			errOp = x
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			goto Error
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		case !allOrdered(y.typ):
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>			errOp = y
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>			goto Error
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		}
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	default:
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		unreachable()
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	}
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	<span class="comment">// comparison is ok</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	if x.mode == constant_ &amp;&amp; y.mode == constant_ {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		x.val = constant.MakeBool(constant.Compare(x.val, op, y.val))
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		<span class="comment">// The operands are never materialized; no need to update</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		<span class="comment">// their types.</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	} else {
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		x.mode = value
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		<span class="comment">// The operands have now their final types, which at run-</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		<span class="comment">// time will be materialized. Update the expression trees.</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		<span class="comment">// If the current types are untyped, the materialized type</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		<span class="comment">// is the respective default type.</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		check.updateExprType(x.expr, Default(x.typ), true)
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		check.updateExprType(y.expr, Default(y.typ), true)
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;Comparison operators compare two operands and yield</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	<span class="comment">//        an untyped boolean value.&#34;</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	x.typ = Typ[UntypedBool]
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	return
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>Error:
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	<span class="comment">// We have an offending operand errOp and possibly an error cause.</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	if cause == &#34;&#34; {
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		if isTypeParam(x.typ) || isTypeParam(y.typ) {
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) should report the specific type causing the problem, if any</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			if !isTypeParam(x.typ) {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>				errOp = y
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>			}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			cause = check.sprintf(&#34;type parameter %s is not comparable with %s&#34;, errOp.typ, op)
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		} else {
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			cause = check.sprintf(&#34;operator %s not defined on %s&#34;, op, check.kindString(errOp.typ)) <span class="comment">// catch-all</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	if switchCase {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		check.errorf(x, code, &#34;invalid case %s in switch on %s (%s)&#34;, x.expr, y.expr, cause) <span class="comment">// error position always at 1st operand</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	} else {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		check.errorf(errOp, code, invalidOp+&#34;%s %s %s (%s)&#34;, x.expr, op, y.expr, cause)
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	x.mode = invalid
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span><span class="comment">// incomparableCause returns a more specific cause why typ is not comparable.</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span><span class="comment">// If there is no more specific cause, the result is &#34;&#34;.</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>func (check *Checker) incomparableCause(typ Type) string {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	switch under(typ).(type) {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	case *Slice, *Signature, *Map:
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		return check.kindString(typ) + &#34; can only be compared to nil&#34;
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	}
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	<span class="comment">// see if we can extract a more specific error</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	var cause string
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	comparable(typ, true, nil, func(format string, args ...interface{}) {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		cause = check.sprintf(format, args...)
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	})
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	return cause
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span><span class="comment">// kindString returns the type kind as a string.</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>func (check *Checker) kindString(typ Type) string {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	switch under(typ).(type) {
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	case *Array:
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		return &#34;array&#34;
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	case *Slice:
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		return &#34;slice&#34;
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	case *Struct:
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		return &#34;struct&#34;
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	case *Pointer:
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		return &#34;pointer&#34;
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	case *Signature:
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		return &#34;func&#34;
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	case *Interface:
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		if isTypeParam(typ) {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>			return check.sprintf(&#34;type parameter %s&#34;, typ)
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		return &#34;interface&#34;
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	case *Map:
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		return &#34;map&#34;
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	case *Chan:
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		return &#34;chan&#34;
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	default:
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		return check.sprintf(&#34;%s&#34;, typ) <span class="comment">// catch-all</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	}
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>}
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span><span class="comment">// If e != nil, it must be the shift expression; it may be nil for non-constant shifts.</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>func (check *Checker) shift(x, y *operand, e ast.Expr, op token.Token) {
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	<span class="comment">// TODO(gri) This function seems overly complex. Revisit.</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	var xval constant.Value
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	if x.mode == constant_ {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		xval = constant.ToInt(x.val)
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	}
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	if allInteger(x.typ) || isUntyped(x.typ) &amp;&amp; xval != nil &amp;&amp; xval.Kind() == constant.Int {
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		<span class="comment">// The lhs is of integer type or an untyped constant representable</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		<span class="comment">// as an integer. Nothing to do.</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	} else {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		<span class="comment">// shift has no chance</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		check.errorf(x, InvalidShiftOperand, invalidOp+&#34;shifted operand %s must be integer&#34;, x)
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		return
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;The right operand in a shift expression must have integer type</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	<span class="comment">// or be an untyped constant representable by a value of type uint.&#34;</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	<span class="comment">// Check that constants are representable by uint, but do not convert them</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	<span class="comment">// (see also go.dev/issue/47243).</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	var yval constant.Value
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	if y.mode == constant_ {
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		<span class="comment">// Provide a good error message for negative shift counts.</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		yval = constant.ToInt(y.val) <span class="comment">// consider -1, 1.0, but not -1.1</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		if yval.Kind() == constant.Int &amp;&amp; constant.Sign(yval) &lt; 0 {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>			check.errorf(y, InvalidShiftCount, invalidOp+&#34;negative shift count %s&#34;, y)
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>			return
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		}
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		if isUntyped(y.typ) {
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>			<span class="comment">// Caution: Check for representability here, rather than in the switch</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>			<span class="comment">// below, because isInteger includes untyped integers (was bug go.dev/issue/43697).</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>			check.representable(y, Typ[Uint])
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>			if y.mode == invalid {
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>				x.mode = invalid
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>				return
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>			}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		}
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	} else {
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		<span class="comment">// Check that RHS is otherwise at least of integer type.</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		switch {
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		case allInteger(y.typ):
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>			if !allUnsigned(y.typ) &amp;&amp; !check.verifyVersionf(y, go1_13, invalidOp+&#34;signed shift count %s&#34;, y) {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>				x.mode = invalid
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>				return
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>			}
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		case isUntyped(y.typ):
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>			<span class="comment">// This is incorrect, but preserves pre-existing behavior.</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>			<span class="comment">// See also go.dev/issue/47410.</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>			check.convertUntyped(y, Typ[Uint])
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>			if y.mode == invalid {
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>				x.mode = invalid
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>				return
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>			}
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		default:
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>			check.errorf(y, InvalidShiftCount, invalidOp+&#34;shift count %s must be integer&#34;, y)
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>			return
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		}
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	}
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	if x.mode == constant_ {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		if y.mode == constant_ {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>			<span class="comment">// if either x or y has an unknown value, the result is unknown</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			if x.val.Kind() == constant.Unknown || y.val.Kind() == constant.Unknown {
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>				x.val = constant.MakeUnknown()
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>				<span class="comment">// ensure the correct type - see comment below</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>				if !isInteger(x.typ) {
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>					x.typ = Typ[UntypedInt]
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>				}
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>				return
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>			}
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>			<span class="comment">// rhs must be within reasonable bounds in constant shifts</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>			const shiftBound = 1023 - 1 + 52 <span class="comment">// so we can express smallestFloat64 (see go.dev/issue/44057)</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>			s, ok := constant.Uint64Val(yval)
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>			if !ok || s &gt; shiftBound {
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>				check.errorf(y, InvalidShiftCount, invalidOp+&#34;invalid shift count %s&#34;, y)
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>				x.mode = invalid
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>				return
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>			}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			<span class="comment">// The lhs is representable as an integer but may not be an integer</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			<span class="comment">// (e.g., 2.0, an untyped float) - this can only happen for untyped</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>			<span class="comment">// non-integer numeric constants. Correct the type so that the shift</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>			<span class="comment">// result is of integer type.</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>			if !isInteger(x.typ) {
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>				x.typ = Typ[UntypedInt]
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>			}
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			<span class="comment">// x is a constant so xval != nil and it must be of Int kind.</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>			x.val = constant.Shift(xval, op, uint(s))
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>			x.expr = e
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			opPos := x.Pos()
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>			if b, _ := e.(*ast.BinaryExpr); b != nil {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>				opPos = b.OpPos
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>			check.overflow(x, opPos)
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>			return
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		}
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		<span class="comment">// non-constant shift with constant lhs</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		if isUntyped(x.typ) {
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>			<span class="comment">// spec: &#34;If the left operand of a non-constant shift</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>			<span class="comment">// expression is an untyped constant, the type of the</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>			<span class="comment">// constant is what it would be if the shift expression</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>			<span class="comment">// were replaced by its left operand alone.&#34;.</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>			<span class="comment">// Delay operand checking until we know the final type</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>			<span class="comment">// by marking the lhs expression as lhs shift operand.</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>			<span class="comment">// Usually (in correct programs), the lhs expression</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>			<span class="comment">// is in the untyped map. However, it is possible to</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>			<span class="comment">// create incorrect programs where the same expression</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>			<span class="comment">// is evaluated twice (via a declaration cycle) such</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>			<span class="comment">// that the lhs expression type is determined in the</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>			<span class="comment">// first round and thus deleted from the map, and then</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>			<span class="comment">// not found in the second round (double insertion of</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>			<span class="comment">// the same expr node still just leads to one entry for</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>			<span class="comment">// that node, and it can only be deleted once).</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>			<span class="comment">// Be cautious and check for presence of entry.</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>			<span class="comment">// Example: var e, f = int(1&lt;&lt;&#34;&#34;[f]) // go.dev/issue/11347</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			if info, found := check.untyped[x.expr]; found {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>				info.isLhs = true
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>				check.untyped[x.expr] = info
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>			}
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>			<span class="comment">// keep x&#39;s type</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			x.mode = value
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>			return
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		}
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	}
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	<span class="comment">// non-constant shift - lhs must be an integer</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	if !allInteger(x.typ) {
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		check.errorf(x, InvalidShiftOperand, invalidOp+&#34;shifted operand %s must be integer&#34;, x)
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		return
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	}
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	x.mode = value
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>}
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>var binaryOpPredicates opPredicates
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>func init() {
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	<span class="comment">// Setting binaryOpPredicates in init avoids declaration cycles.</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	binaryOpPredicates = opPredicates{
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		token.ADD: allNumericOrString,
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		token.SUB: allNumeric,
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>		token.MUL: allNumeric,
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		token.QUO: allNumeric,
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		token.REM: allInteger,
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>		token.AND:     allInteger,
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		token.OR:      allInteger,
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		token.XOR:     allInteger,
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>		token.AND_NOT: allInteger,
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		token.LAND: allBoolean,
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>		token.LOR:  allBoolean,
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>	}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>}
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span><span class="comment">// If e != nil, it must be the binary expression; it may be nil for non-constant expressions</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span><span class="comment">// (when invoked for an assignment operation where the binary expression is implicit).</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>func (check *Checker) binary(x *operand, e ast.Expr, lhs, rhs ast.Expr, op token.Token, opPos token.Pos) {
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	var y operand
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	check.expr(nil, x, lhs)
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	check.expr(nil, &amp;y, rhs)
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	if x.mode == invalid {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		return
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>	}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	if y.mode == invalid {
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>		x.expr = y.expr
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		return
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	}
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	if isShift(op) {
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		check.shift(x, &amp;y, e, op)
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>		return
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	}
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	check.matchTypes(x, &amp;y)
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	if x.mode == invalid {
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		return
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	}
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	if isComparison(op) {
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		check.comparison(x, &amp;y, op, false)
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>		return
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	}
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	if !Identical(x.typ, y.typ) {
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		<span class="comment">// only report an error if we have valid types</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>		<span class="comment">// (otherwise we had an error reported elsewhere already)</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		if isValid(x.typ) &amp;&amp; isValid(y.typ) {
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>			var posn positioner = x
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>			if e != nil {
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>				posn = e
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>			}
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>			if e != nil {
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>				check.errorf(posn, MismatchedTypes, invalidOp+&#34;%s (mismatched types %s and %s)&#34;, e, x.typ, y.typ)
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>			} else {
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>				check.errorf(posn, MismatchedTypes, invalidOp+&#34;%s %s= %s (mismatched types %s and %s)&#34;, lhs, op, rhs, x.typ, y.typ)
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>			}
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		}
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>		return
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	}
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>	if !check.op(binaryOpPredicates, x, op) {
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		return
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	}
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	if op == token.QUO || op == token.REM {
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		<span class="comment">// check for zero divisor</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		if (x.mode == constant_ || allInteger(x.typ)) &amp;&amp; y.mode == constant_ &amp;&amp; constant.Sign(y.val) == 0 {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>			check.error(&amp;y, DivByZero, invalidOp+&#34;division by zero&#34;)
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>			return
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		}
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		<span class="comment">// check for divisor underflow in complex division (see go.dev/issue/20227)</span>
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		if x.mode == constant_ &amp;&amp; y.mode == constant_ &amp;&amp; isComplex(x.typ) {
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>			re, im := constant.Real(y.val), constant.Imag(y.val)
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>			re2, im2 := constant.BinaryOp(re, token.MUL, re), constant.BinaryOp(im, token.MUL, im)
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			if constant.Sign(re2) == 0 &amp;&amp; constant.Sign(im2) == 0 {
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>				check.error(&amp;y, DivByZero, invalidOp+&#34;division by zero&#34;)
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>				x.mode = invalid
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>				return
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>			}
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	}
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	if x.mode == constant_ &amp;&amp; y.mode == constant_ {
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		<span class="comment">// if either x or y has an unknown value, the result is unknown</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>		if x.val.Kind() == constant.Unknown || y.val.Kind() == constant.Unknown {
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>			x.val = constant.MakeUnknown()
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>			<span class="comment">// x.typ is unchanged</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>			return
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		}
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>		<span class="comment">// force integer division of integer operands</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		if op == token.QUO &amp;&amp; isInteger(x.typ) {
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>			op = token.QUO_ASSIGN
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		}
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		x.val = constant.BinaryOp(x.val, op, y.val)
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		x.expr = e
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>		check.overflow(x, opPos)
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		return
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	}
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	x.mode = value
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	<span class="comment">// x.typ is unchanged</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>}
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span><span class="comment">// matchTypes attempts to convert any untyped types x and y such that they match.</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span><span class="comment">// If an error occurs, x.mode is set to invalid.</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>func (check *Checker) matchTypes(x, y *operand) {
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	<span class="comment">// mayConvert reports whether the operands x and y may</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	<span class="comment">// possibly have matching types after converting one</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	<span class="comment">// untyped operand to the type of the other.</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	<span class="comment">// If mayConvert returns true, we try to convert the</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	<span class="comment">// operands to each other&#39;s types, and if that fails</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	<span class="comment">// we report a conversion failure.</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	<span class="comment">// If mayConvert returns false, we continue without an</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	<span class="comment">// attempt at conversion, and if the operand types are</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	<span class="comment">// not compatible, we report a type mismatch error.</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	mayConvert := func(x, y *operand) bool {
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>		<span class="comment">// If both operands are typed, there&#39;s no need for an implicit conversion.</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>		if isTyped(x.typ) &amp;&amp; isTyped(y.typ) {
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			return false
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>		}
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>		<span class="comment">// An untyped operand may convert to its default type when paired with an empty interface</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) This should only matter for comparisons (the only binary operation that is</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>		<span class="comment">//           valid with interfaces), but in that case the assignability check should take</span>
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>		<span class="comment">//           care of the conversion. Verify and possibly eliminate this extra test.</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		if isNonTypeParamInterface(x.typ) || isNonTypeParamInterface(y.typ) {
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>			return true
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>		}
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		<span class="comment">// A boolean type can only convert to another boolean type.</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>		if allBoolean(x.typ) != allBoolean(y.typ) {
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>			return false
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		}
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>		<span class="comment">// A string type can only convert to another string type.</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>		if allString(x.typ) != allString(y.typ) {
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>			return false
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>		}
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>		<span class="comment">// Untyped nil can only convert to a type that has a nil.</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>		if x.isNil() {
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>			return hasNil(y.typ)
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>		}
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>		if y.isNil() {
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>			return hasNil(x.typ)
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>		}
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>		<span class="comment">// An untyped operand cannot convert to a pointer.</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) generalize to type parameters</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		if isPointer(x.typ) || isPointer(y.typ) {
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>			return false
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		}
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		return true
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>	}
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>	if mayConvert(x, y) {
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>		check.convertUntyped(x, y.typ)
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>			return
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>		}
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>		check.convertUntyped(y, x.typ)
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>		if y.mode == invalid {
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>			return
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>		}
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>	}
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>}
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span><span class="comment">// exprKind describes the kind of an expression; the kind</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span><span class="comment">// determines if an expression is valid in &#39;statement context&#39;.</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>type exprKind int
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>const (
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	conversion exprKind = iota
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	expression
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>	statement
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>)
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>
<span id="L944" class="ln">   944&nbsp;&nbsp;</span><span class="comment">// target represent the (signature) type and description of the LHS</span>
<span id="L945" class="ln">   945&nbsp;&nbsp;</span><span class="comment">// variable of an assignment, or of a function result variable.</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>type target struct {
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>	sig  *Signature
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>	desc string
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>}
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span><span class="comment">// newTarget creates a new target for the given type and description.</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span><span class="comment">// The result is nil if typ is not a signature.</span>
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>func newTarget(typ Type, desc string) *target {
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>	if typ != nil {
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>		if sig, _ := under(typ).(*Signature); sig != nil {
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>			return &amp;target{sig, desc}
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>		}
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>	}
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>	return nil
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>}
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span><span class="comment">// rawExpr typechecks expression e and initializes x with the expression</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span><span class="comment">// value or type. If an error occurred, x.mode is set to invalid.</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span><span class="comment">// If a non-nil target T is given and e is a generic function,</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span><span class="comment">// T is used to infer the type arguments for e.</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span><span class="comment">// If hint != nil, it is the type of a composite literal element.</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span><span class="comment">// If allowGeneric is set, the operand type may be an uninstantiated</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span><span class="comment">// parameterized type or function value.</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>func (check *Checker) rawExpr(T *target, x *operand, e ast.Expr, hint Type, allowGeneric bool) exprKind {
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	if check.conf._Trace {
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>		check.trace(e.Pos(), &#34;-- expr %s&#34;, e)
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>		check.indent++
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>		defer func() {
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>			check.indent--
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>			check.trace(e.Pos(), &#34;=&gt; %s&#34;, x)
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>		}()
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>	}
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>	kind := check.exprInternal(T, x, e, hint)
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	if !allowGeneric {
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		check.nonGeneric(T, x)
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>	}
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	check.record(x)
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	return kind
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>}
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span><span class="comment">// If x is a generic type, or a generic function whose type arguments cannot be inferred</span>
<span id="L991" class="ln">   991&nbsp;&nbsp;</span><span class="comment">// from a non-nil target T, nonGeneric reports an error and invalidates x.mode and x.typ.</span>
<span id="L992" class="ln">   992&nbsp;&nbsp;</span><span class="comment">// Otherwise it leaves x alone.</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>func (check *Checker) nonGeneric(T *target, x *operand) {
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	if x.mode == invalid || x.mode == novalue {
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		return
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	}
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	var what string
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	switch t := x.typ.(type) {
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	case *Named:
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>		if isGeneric(t) {
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>			what = &#34;type&#34;
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>		}
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>	case *Signature:
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>		if t.tparams != nil {
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>			if enableReverseTypeInference &amp;&amp; T != nil {
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>				check.funcInst(T, x.Pos(), x, nil, true)
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>				return
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>			}
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>			what = &#34;function&#34;
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>		}
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>	}
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	if what != &#34;&#34; {
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>		check.errorf(x.expr, WrongTypeArgCount, &#34;cannot use generic %s %s without instantiation&#34;, what, x.expr)
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>		x.typ = Typ[Invalid]
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	}
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>}
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span><span class="comment">// exprInternal contains the core of type checking of expressions.</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span><span class="comment">// Must only be called by rawExpr.</span>
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span><span class="comment">// (See rawExpr for an explanation of the parameters.)</span>
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>func (check *Checker) exprInternal(T *target, x *operand, e ast.Expr, hint Type) exprKind {
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	<span class="comment">// make sure x has a valid state in case of bailout</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>	<span class="comment">// (was go.dev/issue/5770)</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>	x.mode = invalid
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	x.typ = Typ[Invalid]
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	switch e := e.(type) {
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	case *ast.BadExpr:
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>		goto Error <span class="comment">// error was reported before</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	case *ast.Ident:
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		check.ident(x, e, nil, false)
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	case *ast.Ellipsis:
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>		<span class="comment">// ellipses are handled explicitly where they are legal</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>		<span class="comment">// (array composite literals and parameter lists)</span>
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>		check.error(e, BadDotDotDotSyntax, &#34;invalid use of &#39;...&#39;&#34;)
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>		goto Error
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>	case *ast.BasicLit:
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		switch e.Kind {
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>		case token.INT, token.FLOAT, token.IMAG:
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>			check.langCompat(e)
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>			<span class="comment">// The max. mantissa precision for untyped numeric values</span>
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>			<span class="comment">// is 512 bits, or 4048 bits for each of the two integer</span>
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>			<span class="comment">// parts of a fraction for floating-point numbers that are</span>
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>			<span class="comment">// represented accurately in the go/constant package.</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>			<span class="comment">// Constant literals that are longer than this many bits</span>
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>			<span class="comment">// are not meaningful; and excessively long constants may</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>			<span class="comment">// consume a lot of space and time for a useless conversion.</span>
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>			<span class="comment">// Cap constant length with a generous upper limit that also</span>
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>			<span class="comment">// allows for separators between all digits.</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>			const limit = 10000
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>			if len(e.Value) &gt; limit {
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>				check.errorf(e, InvalidConstVal, &#34;excessively long constant: %s... (%d chars)&#34;, e.Value[:10], len(e.Value))
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>				goto Error
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>			}
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>		}
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>		x.setConst(e.Kind, e.Value)
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>			<span class="comment">// The parser already establishes syntactic correctness.</span>
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>			<span class="comment">// If we reach here it&#39;s because of number under-/overflow.</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) setConst (and in turn the go/constant package)</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>			<span class="comment">// should return an error describing the issue.</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>			check.errorf(e, InvalidConstVal, &#34;malformed constant: %s&#34;, e.Value)
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>			goto Error
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>		}
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>		<span class="comment">// Ensure that integer values don&#39;t overflow (go.dev/issue/54280).</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>		check.overflow(x, e.Pos())
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>	case *ast.FuncLit:
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>		if sig, ok := check.typ(e.Type).(*Signature); ok {
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>			<span class="comment">// Set the Scope&#39;s extent to the complete &#34;func (...) {...}&#34;</span>
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>			<span class="comment">// so that Scope.Innermost works correctly.</span>
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>			sig.scope.pos = e.Pos()
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>			sig.scope.end = e.End()
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>			if !check.conf.IgnoreFuncBodies &amp;&amp; e.Body != nil {
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>				<span class="comment">// Anonymous functions are considered part of the</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>				<span class="comment">// init expression/func declaration which contains</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>				<span class="comment">// them: use existing package-level declaration info.</span>
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>				decl := check.decl <span class="comment">// capture for use in closure below</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>				iota := check.iota <span class="comment">// capture for use in closure below (go.dev/issue/22345)</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>				<span class="comment">// Don&#39;t type-check right away because the function may</span>
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>				<span class="comment">// be part of a type definition to which the function</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>				<span class="comment">// body refers. Instead, type-check as soon as possible,</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>				<span class="comment">// but before the enclosing scope contents changes (go.dev/issue/22992).</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>				check.later(func() {
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>					check.funcBody(decl, &#34;&lt;function literal&gt;&#34;, sig, e.Body, iota)
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>				}).describef(e, &#34;func literal&#34;)
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>			}
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>			x.mode = value
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>			x.typ = sig
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>		} else {
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>			check.errorf(e, InvalidSyntaxTree, &#34;invalid function literal %s&#34;, e)
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>			goto Error
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>		}
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	case *ast.CompositeLit:
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>		var typ, base Type
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>		switch {
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>		case e.Type != nil:
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>			<span class="comment">// composite literal type present - use it</span>
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>			<span class="comment">// [...]T array types may only appear with composite literals.</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>			<span class="comment">// Check for them here so we don&#39;t have to handle ... in general.</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>			if atyp, _ := e.Type.(*ast.ArrayType); atyp != nil &amp;&amp; atyp.Len != nil {
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>				if ellip, _ := atyp.Len.(*ast.Ellipsis); ellip != nil &amp;&amp; ellip.Elt == nil {
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>					<span class="comment">// We have an &#34;open&#34; [...]T array type.</span>
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>					<span class="comment">// Create a new ArrayType with unknown length (-1)</span>
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>					<span class="comment">// and finish setting it up after analyzing the literal.</span>
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>					typ = &amp;Array{len: -1, elem: check.varType(atyp.Elt)}
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>					base = typ
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>					break
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>				}
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>			}
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>			typ = check.typ(e.Type)
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>			base = typ
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>		case hint != nil:
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>			<span class="comment">// no composite literal type present - use hint (element type of enclosing type)</span>
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>			typ = hint
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>			base, _ = deref(coreType(typ)) <span class="comment">// *T implies &amp;T{}</span>
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>			if base == nil {
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>				check.errorf(e, InvalidLit, &#34;invalid composite literal element type %s (no core type)&#34;, typ)
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>				goto Error
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>			}
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>		default:
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) provide better error messages depending on context</span>
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>			check.error(e, UntypedLit, &#34;missing type in composite literal&#34;)
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>			goto Error
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>		}
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>		switch utyp := coreType(base).(type) {
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>		case *Struct:
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>			<span class="comment">// Prevent crash if the struct referred to is not yet set up.</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>			<span class="comment">// See analogous comment for *Array.</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>			if utyp.fields == nil {
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>				check.error(e, InvalidTypeCycle, &#34;invalid recursive type&#34;)
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>				goto Error
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>			}
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>			if len(e.Elts) == 0 {
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>				break
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>			}
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>			<span class="comment">// Convention for error messages on invalid struct literals:</span>
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>			<span class="comment">// we mention the struct type only if it clarifies the error</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>			<span class="comment">// (e.g., a duplicate field error doesn&#39;t need the struct type).</span>
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>			fields := utyp.fields
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>			if _, ok := e.Elts[0].(*ast.KeyValueExpr); ok {
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>				<span class="comment">// all elements must have keys</span>
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>				visited := make([]bool, len(fields))
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>				for _, e := range e.Elts {
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>					kv, _ := e.(*ast.KeyValueExpr)
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>					if kv == nil {
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>						check.error(e, MixedStructLit, &#34;mixture of field:value and value elements in struct literal&#34;)
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>						continue
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>					}
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>					key, _ := kv.Key.(*ast.Ident)
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>					<span class="comment">// do all possible checks early (before exiting due to errors)</span>
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>					<span class="comment">// so we don&#39;t drop information on the floor</span>
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>					check.expr(nil, x, kv.Value)
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>					if key == nil {
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>						check.errorf(kv, InvalidLitField, &#34;invalid field name %s in struct literal&#34;, kv.Key)
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>						continue
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>					}
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>					i := fieldIndex(utyp.fields, check.pkg, key.Name)
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>					if i &lt; 0 {
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>						check.errorf(kv, MissingLitField, &#34;unknown field %s in struct literal of type %s&#34;, key.Name, base)
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>						continue
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>					}
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>					fld := fields[i]
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>					check.recordUse(key, fld)
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>					etyp := fld.typ
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>					check.assignment(x, etyp, &#34;struct literal&#34;)
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>					<span class="comment">// 0 &lt;= i &lt; len(fields)</span>
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>					if visited[i] {
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>						check.errorf(kv, DuplicateLitField, &#34;duplicate field name %s in struct literal&#34;, key.Name)
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>						continue
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>					}
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>					visited[i] = true
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>				}
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>			} else {
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>				<span class="comment">// no element must have a key</span>
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>				for i, e := range e.Elts {
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>					if kv, _ := e.(*ast.KeyValueExpr); kv != nil {
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>						check.error(kv, MixedStructLit, &#34;mixture of field:value and value elements in struct literal&#34;)
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>						continue
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>					}
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>					check.expr(nil, x, e)
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>					if i &gt;= len(fields) {
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>						check.errorf(x, InvalidStructLit, &#34;too many values in struct literal of type %s&#34;, base)
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>						break <span class="comment">// cannot continue</span>
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>					}
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>					<span class="comment">// i &lt; len(fields)</span>
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>					fld := fields[i]
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>					if !fld.Exported() &amp;&amp; fld.pkg != check.pkg {
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>						check.errorf(x,
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>							UnexportedLitField,
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>							&#34;implicit assignment to unexported field %s in struct literal of type %s&#34;, fld.name, base)
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>						continue
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>					}
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>					etyp := fld.typ
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>					check.assignment(x, etyp, &#34;struct literal&#34;)
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>				}
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>				if len(e.Elts) &lt; len(fields) {
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>					check.errorf(inNode(e, e.Rbrace), InvalidStructLit, &#34;too few values in struct literal of type %s&#34;, base)
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>					<span class="comment">// ok to continue</span>
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>				}
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>			}
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>		case *Array:
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>			<span class="comment">// Prevent crash if the array referred to is not yet set up. Was go.dev/issue/18643.</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>			<span class="comment">// This is a stop-gap solution. Should use Checker.objPath to report entire</span>
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>			<span class="comment">// path starting with earliest declaration in the source. TODO(gri) fix this.</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>			if utyp.elem == nil {
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>				check.error(e, InvalidTypeCycle, &#34;invalid recursive type&#34;)
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>				goto Error
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>			}
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>			n := check.indexedElts(e.Elts, utyp.elem, utyp.len)
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>			<span class="comment">// If we have an array of unknown length (usually [...]T arrays, but also</span>
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>			<span class="comment">// arrays [n]T where n is invalid) set the length now that we know it and</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>			<span class="comment">// record the type for the array (usually done by check.typ which is not</span>
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>			<span class="comment">// called for [...]T). We handle [...]T arrays and arrays with invalid</span>
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>			<span class="comment">// length the same here because it makes sense to &#34;guess&#34; the length for</span>
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>			<span class="comment">// the latter if we have a composite literal; e.g. for [n]int{1, 2, 3}</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>			<span class="comment">// where n is invalid for some reason, it seems fair to assume it should</span>
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>			<span class="comment">// be 3 (see also Checked.arrayLength and go.dev/issue/27346).</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>			if utyp.len &lt; 0 {
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>				utyp.len = n
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>				<span class="comment">// e.Type is missing if we have a composite literal element</span>
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>				<span class="comment">// that is itself a composite literal with omitted type. In</span>
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>				<span class="comment">// that case there is nothing to record (there is no type in</span>
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>				<span class="comment">// the source at that point).</span>
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>				if e.Type != nil {
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>					check.recordTypeAndValue(e.Type, typexpr, utyp, nil)
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>				}
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>			}
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>		case *Slice:
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>			<span class="comment">// Prevent crash if the slice referred to is not yet set up.</span>
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>			<span class="comment">// See analogous comment for *Array.</span>
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>			if utyp.elem == nil {
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>				check.error(e, InvalidTypeCycle, &#34;invalid recursive type&#34;)
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>				goto Error
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>			}
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>			check.indexedElts(e.Elts, utyp.elem, -1)
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>		case *Map:
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>			<span class="comment">// Prevent crash if the map referred to is not yet set up.</span>
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>			<span class="comment">// See analogous comment for *Array.</span>
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>			if utyp.key == nil || utyp.elem == nil {
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>				check.error(e, InvalidTypeCycle, &#34;invalid recursive type&#34;)
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>				goto Error
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>			}
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>			<span class="comment">// If the map key type is an interface (but not a type parameter),</span>
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>			<span class="comment">// the type of a constant key must be considered when checking for</span>
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>			<span class="comment">// duplicates.</span>
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>			keyIsInterface := isNonTypeParamInterface(utyp.key)
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>			visited := make(map[any][]Type, len(e.Elts))
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>			for _, e := range e.Elts {
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>				kv, _ := e.(*ast.KeyValueExpr)
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>				if kv == nil {
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>					check.error(e, MissingLitKey, &#34;missing key in map literal&#34;)
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>					continue
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>				}
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>				check.exprWithHint(x, kv.Key, utyp.key)
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>				check.assignment(x, utyp.key, &#34;map literal&#34;)
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>				if x.mode == invalid {
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>					continue
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>				}
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>				if x.mode == constant_ {
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>					duplicate := false
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>					xkey := keyVal(x.val)
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>					if keyIsInterface {
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>						for _, vtyp := range visited[xkey] {
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>							if Identical(vtyp, x.typ) {
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>								duplicate = true
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>								break
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>							}
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>						}
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>						visited[xkey] = append(visited[xkey], x.typ)
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>					} else {
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>						_, duplicate = visited[xkey]
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>						visited[xkey] = nil
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>					}
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>					if duplicate {
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>						check.errorf(x, DuplicateLitKey, &#34;duplicate key %s in map literal&#34;, x.val)
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>						continue
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>					}
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>				}
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>				check.exprWithHint(x, kv.Value, utyp.elem)
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>				check.assignment(x, utyp.elem, &#34;map literal&#34;)
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>			}
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>		default:
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>			<span class="comment">// when &#34;using&#34; all elements unpack KeyValueExpr</span>
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>			<span class="comment">// explicitly because check.use doesn&#39;t accept them</span>
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>			for _, e := range e.Elts {
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>				if kv, _ := e.(*ast.KeyValueExpr); kv != nil {
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>					<span class="comment">// Ideally, we should also &#34;use&#34; kv.Key but we can&#39;t know</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>					<span class="comment">// if it&#39;s an externally defined struct key or not. Going</span>
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>					<span class="comment">// forward anyway can lead to other errors. Give up instead.</span>
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>					e = kv.Value
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>				}
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>				check.use(e)
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>			}
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>			<span class="comment">// if utyp is invalid, an error was reported before</span>
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>			if isValid(utyp) {
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>				check.errorf(e, InvalidLit, &#34;invalid composite literal type %s&#34;, typ)
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>				goto Error
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>			}
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>		}
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>		x.mode = value
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>		x.typ = typ
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>	case *ast.ParenExpr:
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>		<span class="comment">// type inference doesn&#39;t go past parentheses (targe type T = nil)</span>
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>		kind := check.rawExpr(nil, x, e.X, nil, false)
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>		x.expr = e
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>		return kind
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>	case *ast.SelectorExpr:
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>		check.selector(x, e, nil, false)
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>	case *ast.IndexExpr, *ast.IndexListExpr:
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>		ix := typeparams.UnpackIndexExpr(e)
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>		if check.indexExpr(x, ix) {
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>			if !enableReverseTypeInference {
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>				T = nil
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>			}
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>			check.funcInst(T, e.Pos(), x, ix, true)
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>		}
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>			goto Error
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>		}
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>	case *ast.SliceExpr:
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>		check.sliceExpr(x, e)
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>			goto Error
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>		}
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>	case *ast.TypeAssertExpr:
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>		check.expr(nil, x, e.X)
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>			goto Error
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>		}
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>		<span class="comment">// x.(type) expressions are handled explicitly in type switches</span>
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>		if e.Type == nil {
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>			<span class="comment">// Don&#39;t use InvalidSyntaxTree because this can occur in the AST produced by</span>
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>			<span class="comment">// go/parser.</span>
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>			check.error(e, BadTypeKeyword, &#34;use of .(type) outside type switch&#34;)
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>			goto Error
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>		}
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>		if isTypeParam(x.typ) {
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>			check.errorf(x, InvalidAssert, invalidOp+&#34;cannot use type assertion on type parameter value %s&#34;, x)
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>			goto Error
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>		}
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>		if _, ok := under(x.typ).(*Interface); !ok {
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>			check.errorf(x, InvalidAssert, invalidOp+&#34;%s is not an interface&#34;, x)
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>			goto Error
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>		}
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>		T := check.varType(e.Type)
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>		if !isValid(T) {
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>			goto Error
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>		}
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>		check.typeAssertion(e, x, T, false)
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>		x.mode = commaok
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>		x.typ = T
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>	case *ast.CallExpr:
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>		return check.callExpr(x, e)
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>	case *ast.StarExpr:
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>		check.exprOrType(x, e.X, false)
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>		switch x.mode {
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>		case invalid:
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>			goto Error
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>		case typexpr:
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>			check.validVarType(e.X, x.typ)
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>			x.typ = &amp;Pointer{base: x.typ}
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>		default:
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>			var base Type
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>			if !underIs(x.typ, func(u Type) bool {
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>				p, _ := u.(*Pointer)
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>				if p == nil {
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>					check.errorf(x, InvalidIndirection, invalidOp+&#34;cannot indirect %s&#34;, x)
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>					return false
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>				}
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>				if base != nil &amp;&amp; !Identical(p.base, base) {
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>					check.errorf(x, InvalidIndirection, invalidOp+&#34;pointers of %s must have identical base types&#34;, x)
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>					return false
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>				}
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>				base = p.base
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>				return true
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>			}) {
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>				goto Error
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>			}
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>			x.mode = variable
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>			x.typ = base
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>		}
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>	case *ast.UnaryExpr:
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>		check.unary(x, e)
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>			goto Error
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>		}
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>		if e.Op == token.ARROW {
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>			x.expr = e
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>			return statement <span class="comment">// receive operations may appear in statement context</span>
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>		}
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>		check.binary(x, e, e.X, e.Y, e.Op, e.OpPos)
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>			goto Error
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>		}
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>	case *ast.KeyValueExpr:
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>		<span class="comment">// key:value expressions are handled in composite literals</span>
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>		check.error(e, InvalidSyntaxTree, &#34;no key:value expected&#34;)
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>		goto Error
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>	case *ast.ArrayType, *ast.StructType, *ast.FuncType,
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>		*ast.InterfaceType, *ast.MapType, *ast.ChanType:
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>		x.mode = typexpr
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>		x.typ = check.typ(e)
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span>		<span class="comment">// Note: rawExpr (caller of exprInternal) will call check.recordTypeAndValue</span>
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span>		<span class="comment">// even though check.typ has already called it. This is fine as both</span>
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>		<span class="comment">// times the same expression and type are recorded. It is also not a</span>
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>		<span class="comment">// performance issue because we only reach here for composite literal</span>
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>		<span class="comment">// types, which are comparatively rare.</span>
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span>
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span>	default:
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span>		panic(fmt.Sprintf(&#34;%s: unknown expression type %T&#34;, check.fset.Position(e.Pos()), e))
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span>	}
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span>
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span>	<span class="comment">// everything went well</span>
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>	x.expr = e
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>	return expression
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>Error:
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>	x.mode = invalid
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>	x.expr = e
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>	return statement <span class="comment">// avoid follow-up errors</span>
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>}
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span><span class="comment">// keyVal maps a complex, float, integer, string or boolean constant value</span>
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span><span class="comment">// to the corresponding complex128, float64, int64, uint64, string, or bool</span>
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span><span class="comment">// Go value if possible; otherwise it returns x.</span>
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span><span class="comment">// A complex constant that can be represented as a float (such as 1.2 + 0i)</span>
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span><span class="comment">// is returned as a floating point value; if a floating point value can be</span>
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span><span class="comment">// represented as an integer (such as 1.0) it is returned as an integer value.</span>
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span><span class="comment">// This ensures that constants of different kind but equal value (such as</span>
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span><span class="comment">// 1.0 + 0i, 1.0, 1) result in the same value.</span>
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span>func keyVal(x constant.Value) interface{} {
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span>	switch x.Kind() {
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span>	case constant.Complex:
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span>		f := constant.ToFloat(x)
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span>		if f.Kind() != constant.Float {
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span>			r, _ := constant.Float64Val(constant.Real(x))
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span>			i, _ := constant.Float64Val(constant.Imag(x))
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span>			return complex(r, i)
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span>		}
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span>		x = f
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span>		fallthrough
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span>	case constant.Float:
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span>		i := constant.ToInt(x)
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span>		if i.Kind() != constant.Int {
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span>			v, _ := constant.Float64Val(x)
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span>			return v
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span>		}
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span>		x = i
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span>		fallthrough
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span>	case constant.Int:
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span>		if v, ok := constant.Int64Val(x); ok {
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span>			return v
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span>		}
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span>		if v, ok := constant.Uint64Val(x); ok {
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span>			return v
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span>		}
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span>	case constant.String:
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span>		return constant.StringVal(x)
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span>	case constant.Bool:
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span>		return constant.BoolVal(x)
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span>	}
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span>	return x
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span>}
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span>
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span><span class="comment">// typeAssertion checks x.(T). The type of x must be an interface.</span>
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span>func (check *Checker) typeAssertion(e ast.Expr, x *operand, T Type, typeSwitch bool) {
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span>	var cause string
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span>	if check.assertableTo(x.typ, T, &amp;cause) {
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span>		return <span class="comment">// success</span>
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span>	}
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span>
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span>	if typeSwitch {
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span>		check.errorf(e, ImpossibleAssert, &#34;impossible type switch case: %s\n\t%s cannot have dynamic type %s %s&#34;, e, x, T, cause)
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span>		return
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span>	}
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span>
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span>	check.errorf(e, ImpossibleAssert, &#34;impossible type assertion: %s\n\t%s does not implement %s %s&#34;, e, T, x.typ, cause)
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span>}
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span>
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span><span class="comment">// expr typechecks expression e and initializes x with the expression value.</span>
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span><span class="comment">// If a non-nil target T is given and e is a generic function or</span>
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span><span class="comment">// a function call, T is used to infer the type arguments for e.</span>
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span><span class="comment">// The result must be a single value.</span>
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span><span class="comment">// If an error occurred, x.mode is set to invalid.</span>
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span>func (check *Checker) expr(T *target, x *operand, e ast.Expr) {
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span>	check.rawExpr(T, x, e, nil, false)
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span>	check.exclude(x, 1&lt;&lt;novalue|1&lt;&lt;builtin|1&lt;&lt;typexpr)
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>	check.singleValue(x)
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>}
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span>
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span><span class="comment">// genericExpr is like expr but the result may also be generic.</span>
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span>func (check *Checker) genericExpr(x *operand, e ast.Expr) {
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span>	check.rawExpr(nil, x, e, nil, true)
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span>	check.exclude(x, 1&lt;&lt;novalue|1&lt;&lt;builtin|1&lt;&lt;typexpr)
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span>	check.singleValue(x)
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span>}
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span>
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span><span class="comment">// multiExpr typechecks e and returns its value (or values) in list.</span>
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span><span class="comment">// If allowCommaOk is set and e is a map index, comma-ok, or comma-err</span>
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span><span class="comment">// expression, the result is a two-element list containing the value</span>
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span><span class="comment">// of e, and an untyped bool value or an error value, respectively.</span>
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span><span class="comment">// If an error occurred, list[0] is not valid.</span>
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>func (check *Checker) multiExpr(e ast.Expr, allowCommaOk bool) (list []*operand, commaOk bool) {
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>	var x operand
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>	check.rawExpr(nil, &amp;x, e, nil, false)
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>	check.exclude(&amp;x, 1&lt;&lt;novalue|1&lt;&lt;builtin|1&lt;&lt;typexpr)
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span>
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span>	if t, ok := x.typ.(*Tuple); ok &amp;&amp; x.mode != invalid {
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span>		<span class="comment">// multiple values</span>
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span>		list = make([]*operand, t.Len())
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>		for i, v := range t.vars {
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>			list[i] = &amp;operand{mode: value, expr: e, typ: v.typ}
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>		}
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span>		return
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span>	}
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span>
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span>	<span class="comment">// exactly one (possibly invalid or comma-ok) value</span>
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>	list = []*operand{&amp;x}
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>	if allowCommaOk &amp;&amp; (x.mode == mapindex || x.mode == commaok || x.mode == commaerr) {
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>		x2 := &amp;operand{mode: value, expr: e, typ: Typ[UntypedBool]}
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>		if x.mode == commaerr {
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span>			x2.typ = universeError
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span>		}
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span>		list = append(list, x2)
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span>		commaOk = true
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span>	}
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span>
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span>	return
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span>}
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span>
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span><span class="comment">// exprWithHint typechecks expression e and initializes x with the expression value;</span>
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span><span class="comment">// hint is the type of a composite literal element.</span>
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span><span class="comment">// If an error occurred, x.mode is set to invalid.</span>
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span>func (check *Checker) exprWithHint(x *operand, e ast.Expr, hint Type) {
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span>	assert(hint != nil)
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span>	check.rawExpr(nil, x, e, hint, false)
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span>	check.exclude(x, 1&lt;&lt;novalue|1&lt;&lt;builtin|1&lt;&lt;typexpr)
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span>	check.singleValue(x)
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span>}
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span>
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span><span class="comment">// exprOrType typechecks expression or type e and initializes x with the expression value or type.</span>
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span><span class="comment">// If allowGeneric is set, the operand type may be an uninstantiated parameterized type or function</span>
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span><span class="comment">// value.</span>
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span><span class="comment">// If an error occurred, x.mode is set to invalid.</span>
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span>func (check *Checker) exprOrType(x *operand, e ast.Expr, allowGeneric bool) {
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span>	check.rawExpr(nil, x, e, nil, allowGeneric)
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span>	check.exclude(x, 1&lt;&lt;novalue)
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span>	check.singleValue(x)
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span>}
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span>
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span><span class="comment">// exclude reports an error if x.mode is in modeset and sets x.mode to invalid.</span>
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span><span class="comment">// The modeset may contain any of 1&lt;&lt;novalue, 1&lt;&lt;builtin, 1&lt;&lt;typexpr.</span>
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span>func (check *Checker) exclude(x *operand, modeset uint) {
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span>	if modeset&amp;(1&lt;&lt;x.mode) != 0 {
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span>		var msg string
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>		var code Code
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span>		switch x.mode {
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span>		case novalue:
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span>			if modeset&amp;(1&lt;&lt;typexpr) != 0 {
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span>				msg = &#34;%s used as value&#34;
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span>			} else {
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>				msg = &#34;%s used as value or type&#34;
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>			}
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>			code = TooManyValues
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>		case builtin:
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>			msg = &#34;%s must be called&#34;
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>			code = UncalledBuiltin
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>		case typexpr:
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span>			msg = &#34;%s is not an expression&#34;
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span>			code = NotAnExpr
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span>		default:
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span>			unreachable()
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span>		}
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>		check.errorf(x, code, msg, x)
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span>	}
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span>}
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span>
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span><span class="comment">// singleValue reports an error if x describes a tuple and sets x.mode to invalid.</span>
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span>func (check *Checker) singleValue(x *operand) {
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span>	if x.mode == value {
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span>		<span class="comment">// tuple types are never named - no need for underlying type below</span>
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span>		if t, ok := x.typ.(*Tuple); ok {
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span>			assert(t.Len() != 1)
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span>			check.errorf(x, TooManyValues, &#34;multiple-value %s in single-value context&#34;, x)
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span>		}
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>	}
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>}
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>
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
