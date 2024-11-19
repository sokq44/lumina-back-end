<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/call.go - Go Documentation Server</title>

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
<a href="call.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">call.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/types">go/types</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2013 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements typechecking of call and selector expressions.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package types
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/internal/typeparams&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// funcInst type-checks a function instantiation.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// The incoming x must be a generic function.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// If ix != nil, it provides some or all of the type arguments (ix.Indices).</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// If target != nil, it may be used to infer missing type arguments of x, if any.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// At least one of T or ix must be provided.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// There are two modes of operation:</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//  1. If infer == true, funcInst infers missing type arguments as needed and</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//     instantiates the function x. The returned results are nil.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//  2. If infer == false and inst provides all type arguments, funcInst</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//     instantiates the function x. The returned results are nil.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//     If inst doesn&#39;t provide enough type arguments, funcInst returns the</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//     available arguments and the corresponding expression list; x remains</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//     unchanged.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// If an error (other than a version error) occurs in any case, it is reported</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// and x.mode is set to invalid.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>func (check *Checker) funcInst(T *target, pos token.Pos, x *operand, ix *typeparams.IndexExpr, infer bool) ([]Type, []ast.Expr) {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	assert(T != nil || ix != nil)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	var instErrPos positioner
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	if ix != nil {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		instErrPos = inNode(ix.Orig, ix.Lbrack)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	} else {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		instErrPos = atPos(pos)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	versionErr := !check.verifyVersionf(instErrPos, go1_18, &#34;function instantiation&#34;)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// targs and xlist are the type arguments and corresponding type expressions, or nil.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	var targs []Type
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	var xlist []ast.Expr
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	if ix != nil {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		xlist = ix.Indices
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		targs = check.typeList(xlist)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		if targs == nil {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			x.expr = ix
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			return nil, nil
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		assert(len(targs) == len(xlist))
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// Check the number of type arguments (got) vs number of type parameters (want).</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// Note that x is a function value, not a type expression, so we don&#39;t need to</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// call under below.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	sig := x.typ.(*Signature)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	got, want := len(targs), sig.TypeParams().Len()
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	if got &gt; want {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		<span class="comment">// Providing too many type arguments is always an error.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		check.errorf(ix.Indices[got-1], WrongTypeArgCount, &#34;got %d type arguments but want %d&#34;, got, want)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		x.expr = ix.Orig
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		return nil, nil
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	if got &lt; want {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		if !infer {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			return targs, xlist
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		<span class="comment">// If the uninstantiated or partially instantiated function x is used in</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		<span class="comment">// an assignment (tsig != nil), infer missing type arguments by treating</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		<span class="comment">// the assignment</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		<span class="comment">//    var tvar tsig = x</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		<span class="comment">// like a call g(tvar) of the synthetic generic function g</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		<span class="comment">//    func g[type_parameters_of_x](func_type_of_x)</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		var args []*operand
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		var params []*Var
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		var reverse bool
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		if T != nil &amp;&amp; sig.tparams != nil {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			if !versionErr &amp;&amp; !check.allowVersion(check.pkg, instErrPos, go1_21) {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>				if ix != nil {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>					check.versionErrorf(instErrPos, go1_21, &#34;partially instantiated function in assignment&#34;)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>				} else {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>					check.versionErrorf(instErrPos, go1_21, &#34;implicitly instantiated function in assignment&#34;)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>				}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			gsig := NewSignatureType(nil, nil, nil, sig.params, sig.results, sig.variadic)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			params = []*Var{NewVar(x.Pos(), check.pkg, &#34;&#34;, gsig)}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			<span class="comment">// The type of the argument operand is tsig, which is the type of the LHS in an assignment</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>			<span class="comment">// or the result type in a return statement. Create a pseudo-expression for that operand</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			<span class="comment">// that makes sense when reported in error messages from infer, below.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			expr := ast.NewIdent(T.desc)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			expr.NamePos = x.Pos() <span class="comment">// correct position</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			args = []*operand{{mode: value, expr: expr, typ: T.sig}}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			reverse = true
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		<span class="comment">// Rename type parameters to avoid problems with recursive instantiations.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		<span class="comment">// Note that NewTuple(params...) below is (*Tuple)(nil) if len(params) == 0, as desired.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		tparams, params2 := check.renameTParams(pos, sig.TypeParams().list(), NewTuple(params...))
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		targs = check.infer(atPos(pos), tparams, targs, params2.(*Tuple), args, reverse)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		if targs == nil {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			<span class="comment">// error was already reported</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			x.expr = ix <span class="comment">// TODO(gri) is this correct?</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			return nil, nil
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		got = len(targs)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	assert(got == want)
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// instantiate function signature</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	expr := x.expr <span class="comment">// if we don&#39;t have an index expression, keep the existing expression of x</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if ix != nil {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		expr = ix.Orig
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	sig = check.instantiateSignature(x.Pos(), expr, sig, targs, xlist)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	x.typ = sig
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	x.mode = value
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	x.expr = expr
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	return nil, nil
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>func (check *Checker) instantiateSignature(pos token.Pos, expr ast.Expr, typ *Signature, targs []Type, xlist []ast.Expr) (res *Signature) {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	assert(check != nil)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	assert(len(targs) == typ.TypeParams().Len())
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	if check.conf._Trace {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		check.trace(pos, &#34;-- instantiating signature %s with %s&#34;, typ, targs)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		check.indent++
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		defer func() {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			check.indent--
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			check.trace(pos, &#34;=&gt; %s (under = %s)&#34;, res, res.Underlying())
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		}()
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	inst := check.instance(pos, typ, targs, nil, check.context()).(*Signature)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	assert(inst.TypeParams().Len() == 0) <span class="comment">// signature is not generic anymore</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	check.recordInstance(expr, targs, inst)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	assert(len(xlist) &lt;= len(targs))
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// verify instantiation lazily (was go.dev/issue/50450)</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	check.later(func() {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		tparams := typ.TypeParams().list()
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		if i, err := check.verify(pos, tparams, targs, check.context()); err != nil {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			<span class="comment">// best position for error reporting</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			pos := pos
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			if i &lt; len(xlist) {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>				pos = xlist[i].Pos()
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>			check.softErrorf(atPos(pos), InvalidTypeArg, &#34;%s&#34;, err)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		} else {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			check.mono.recordInstance(check.pkg, pos, tparams, targs, xlist)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	}).describef(atPos(pos), &#34;verify instantiation&#34;)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	return inst
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>func (check *Checker) callExpr(x *operand, call *ast.CallExpr) exprKind {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	ix := typeparams.UnpackIndexExpr(call.Fun)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	if ix != nil {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		if check.indexExpr(x, ix) {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			<span class="comment">// Delay function instantiation to argument checking,</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			<span class="comment">// where we combine type and value arguments for type</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			<span class="comment">// inference.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			assert(x.mode == value)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		} else {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			ix = nil
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		x.expr = call.Fun
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		check.record(x)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	} else {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		check.exprOrType(x, call.Fun, true)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// x.typ may be generic</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	switch x.mode {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	case invalid:
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		check.use(call.Args...)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		x.expr = call
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		return statement
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	case typexpr:
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		<span class="comment">// conversion</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		check.nonGeneric(nil, x)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			return conversion
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		T := x.typ
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		switch n := len(call.Args); n {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		case 0:
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			check.errorf(inNode(call, call.Rparen), WrongArgCount, &#34;missing argument in conversion to %s&#34;, T)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		case 1:
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			check.expr(nil, x, call.Args[0])
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			if x.mode != invalid {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>				if call.Ellipsis.IsValid() {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>					check.errorf(call.Args[0], BadDotDotDotSyntax, &#34;invalid use of ... in conversion to %s&#34;, T)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>					break
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>				}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>				if t, _ := under(T).(*Interface); t != nil &amp;&amp; !isTypeParam(T) {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>					if !t.IsMethodSet() {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>						check.errorf(call, MisplacedConstraintIface, &#34;cannot use interface %s in conversion (contains specific type constraints or is comparable)&#34;, T)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>						break
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>					}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>				}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				check.conversion(x, T)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		default:
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			check.use(call.Args...)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			check.errorf(call.Args[n-1], WrongArgCount, &#34;too many arguments in conversion to %s&#34;, T)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		x.expr = call
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		return conversion
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	case builtin:
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		<span class="comment">// no need to check for non-genericity here</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		id := x.id
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		if !check.builtin(x, call, id) {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		x.expr = call
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		<span class="comment">// a non-constant result implies a function call</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		if x.mode != invalid &amp;&amp; x.mode != constant_ {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			check.hasCallOrRecv = true
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		return predeclaredFuncs[id].kind
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	<span class="comment">// ordinary function/method call</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// signature may be generic</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	cgocall := x.mode == cgofunc
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// a type parameter may be &#34;called&#34; if all types have the same signature</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	sig, _ := coreType(x.typ).(*Signature)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	if sig == nil {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		check.errorf(x, InvalidCall, invalidOp+&#34;cannot call non-function %s&#34;, x)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		x.expr = call
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		return statement
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// Capture wasGeneric before sig is potentially instantiated below.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	wasGeneric := sig.TypeParams().Len() &gt; 0
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// evaluate type arguments, if any</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	var xlist []ast.Expr
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	var targs []Type
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	if ix != nil {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		xlist = ix.Indices
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		targs = check.typeList(xlist)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		if targs == nil {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			check.use(call.Args...)
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			x.expr = call
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			return statement
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		assert(len(targs) == len(xlist))
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		<span class="comment">// check number of type arguments (got) vs number of type parameters (want)</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		got, want := len(targs), sig.TypeParams().Len()
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		if got &gt; want {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			check.errorf(xlist[want], WrongTypeArgCount, &#34;got %d type arguments but want %d&#34;, got, want)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			check.use(call.Args...)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			x.expr = call
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			return statement
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		<span class="comment">// If sig is generic and all type arguments are provided, preempt function</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		<span class="comment">// argument type inference by explicitly instantiating the signature. This</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		<span class="comment">// ensures that we record accurate type information for sig, even if there</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		<span class="comment">// is an error checking its arguments (for example, if an incorrect number</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		<span class="comment">// of arguments is supplied).</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		if got == want &amp;&amp; want &gt; 0 {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			check.verifyVersionf(atPos(ix.Lbrack), go1_18, &#34;function instantiation&#34;)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			sig = check.instantiateSignature(ix.Pos(), ix.Orig, sig, targs, xlist)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			<span class="comment">// targs have been consumed; proceed with checking arguments of the</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			<span class="comment">// non-generic signature.</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			targs = nil
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			xlist = nil
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	<span class="comment">// evaluate arguments</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	args, atargs, atxlist := check.genericExprList(call.Args)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	sig = check.arguments(call, sig, targs, xlist, args, atargs, atxlist)
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	if wasGeneric &amp;&amp; sig.TypeParams().Len() == 0 {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		<span class="comment">// Update the recorded type of call.Fun to its instantiated type.</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		check.recordTypeAndValue(call.Fun, value, sig, nil)
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	<span class="comment">// determine result</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	switch sig.results.Len() {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	case 0:
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		x.mode = novalue
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	case 1:
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		if cgocall {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			x.mode = commaerr
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		} else {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			x.mode = value
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		x.typ = sig.results.vars[0].typ <span class="comment">// unpack tuple</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	default:
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		x.mode = value
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		x.typ = sig.results
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	x.expr = call
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	check.hasCallOrRecv = true
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	<span class="comment">// if type inference failed, a parameterized result must be invalidated</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">// (operands cannot have a parameterized type)</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	if x.mode == value &amp;&amp; sig.TypeParams().Len() &gt; 0 &amp;&amp; isParameterized(sig.TypeParams().list(), x.typ) {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	return statement
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">// exprList evaluates a list of expressions and returns the corresponding operands.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// A single-element expression list may evaluate to multiple operands.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>func (check *Checker) exprList(elist []ast.Expr) (xlist []*operand) {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	if n := len(elist); n == 1 {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		xlist, _ = check.multiExpr(elist[0], false)
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	} else if n &gt; 1 {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		<span class="comment">// multiple (possibly invalid) values</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		xlist = make([]*operand, n)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		for i, e := range elist {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			var x operand
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>			check.expr(nil, &amp;x, e)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>			xlist[i] = &amp;x
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	return
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// genericExprList is like exprList but result operands may be uninstantiated or partially</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// instantiated generic functions (where constraint information is insufficient to infer</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// the missing type arguments) for Go 1.21 and later.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">// For each non-generic or uninstantiated generic operand, the corresponding targsList and</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// xlistList elements do not exist (targsList and xlistList are nil) or the elements are nil.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span><span class="comment">// For each partially instantiated generic function operand, the corresponding targsList and</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// xlistList elements are the operand&#39;s partial type arguments and type expression lists.</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>func (check *Checker) genericExprList(elist []ast.Expr) (resList []*operand, targsList [][]Type, xlistList [][]ast.Expr) {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	if debug {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		defer func() {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			<span class="comment">// targsList and xlistList must have matching lengths</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>			assert(len(targsList) == len(xlistList))
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			<span class="comment">// type arguments must only exist for partially instantiated functions</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			for i, x := range resList {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>				if i &lt; len(targsList) {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>					if n := len(targsList[i]); n &gt; 0 {
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>						<span class="comment">// x must be a partially instantiated function</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>						assert(n &lt; x.typ.(*Signature).TypeParams().Len())
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>					}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>				}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			}
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		}()
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	<span class="comment">// Before Go 1.21, uninstantiated or partially instantiated argument functions are</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	<span class="comment">// nor permitted. Checker.funcInst must infer missing type arguments in that case.</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	infer := true <span class="comment">// for -lang &lt; go1.21</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	n := len(elist)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	if n &gt; 0 &amp;&amp; check.allowVersion(check.pkg, elist[0], go1_21) {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		infer = false
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	if n == 1 {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		<span class="comment">// single value (possibly a partially instantiated function), or a multi-valued expression</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		e := elist[0]
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		var x operand
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		if ix := typeparams.UnpackIndexExpr(e); ix != nil &amp;&amp; check.indexExpr(&amp;x, ix) {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			<span class="comment">// x is a generic function.</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			targs, xlist := check.funcInst(nil, x.Pos(), &amp;x, ix, infer)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			if targs != nil {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>				<span class="comment">// x was not instantiated: collect the (partial) type arguments.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>				targsList = [][]Type{targs}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>				xlistList = [][]ast.Expr{xlist}
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>				<span class="comment">// Update x.expr so that we can record the partially instantiated function.</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>				x.expr = ix.Orig
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			} else {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>				<span class="comment">// x was instantiated: we must record it here because we didn&#39;t</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>				<span class="comment">// use the usual expression evaluators.</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>				check.record(&amp;x)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			resList = []*operand{&amp;x}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		} else {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			<span class="comment">// x is not a function instantiation (it may still be a generic function).</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			check.rawExpr(nil, &amp;x, e, nil, true)
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			check.exclude(&amp;x, 1&lt;&lt;novalue|1&lt;&lt;builtin|1&lt;&lt;typexpr)
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			if t, ok := x.typ.(*Tuple); ok &amp;&amp; x.mode != invalid {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>				<span class="comment">// x is a function call returning multiple values; it cannot be generic.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>				resList = make([]*operand, t.Len())
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>				for i, v := range t.vars {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>					resList[i] = &amp;operand{mode: value, expr: e, typ: v.typ}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>				}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			} else {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>				<span class="comment">// x is exactly one value (possibly invalid or uninstantiated generic function).</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>				resList = []*operand{&amp;x}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			}
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	} else if n &gt; 1 {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		<span class="comment">// multiple values</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		resList = make([]*operand, n)
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		targsList = make([][]Type, n)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		xlistList = make([][]ast.Expr, n)
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		for i, e := range elist {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>			var x operand
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			if ix := typeparams.UnpackIndexExpr(e); ix != nil &amp;&amp; check.indexExpr(&amp;x, ix) {
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>				<span class="comment">// x is a generic function.</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>				targs, xlist := check.funcInst(nil, x.Pos(), &amp;x, ix, infer)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>				if targs != nil {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>					<span class="comment">// x was not instantiated: collect the (partial) type arguments.</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>					targsList[i] = targs
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>					xlistList[i] = xlist
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>					<span class="comment">// Update x.expr so that we can record the partially instantiated function.</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>					x.expr = ix.Orig
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>				} else {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>					<span class="comment">// x was instantiated: we must record it here because we didn&#39;t</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>					<span class="comment">// use the usual expression evaluators.</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>					check.record(&amp;x)
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>				}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			} else {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>				<span class="comment">// x is exactly one value (possibly invalid or uninstantiated generic function).</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>				check.genericExpr(&amp;x, e)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>			resList[i] = &amp;x
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	return
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span><span class="comment">// arguments type-checks arguments passed to a function call with the given signature.</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span><span class="comment">// The function and its arguments may be generic, and possibly partially instantiated.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span><span class="comment">// targs and xlist are the function&#39;s type arguments (and corresponding expressions).</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span><span class="comment">// args are the function arguments. If an argument args[i] is a partially instantiated</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span><span class="comment">// generic function, atargs[i] and atxlist[i] are the corresponding type arguments</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span><span class="comment">// (and corresponding expressions).</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span><span class="comment">// If the callee is variadic, arguments adjusts its signature to match the provided</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span><span class="comment">// arguments. The type parameters and arguments of the callee and all its arguments</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span><span class="comment">// are used together to infer any missing type arguments, and the callee and argument</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span><span class="comment">// functions are instantiated as necessary.</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span><span class="comment">// The result signature is the (possibly adjusted and instantiated) function signature.</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">// If an error occurred, the result signature is the incoming sig.</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>func (check *Checker) arguments(call *ast.CallExpr, sig *Signature, targs []Type, xlist []ast.Expr, args []*operand, atargs [][]Type, atxlist [][]ast.Expr) (rsig *Signature) {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	rsig = sig
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	<span class="comment">// Function call argument/parameter count requirements</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	<span class="comment">//               | standard call    | dotdotdot call |</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	<span class="comment">// --------------+------------------+----------------+</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	<span class="comment">// standard func | nargs == npars   | invalid        |</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	<span class="comment">// --------------+------------------+----------------+</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	<span class="comment">// variadic func | nargs &gt;= npars-1 | nargs == npars |</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	<span class="comment">// --------------+------------------+----------------+</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	nargs := len(args)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	npars := sig.params.Len()
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	ddd := call.Ellipsis.IsValid()
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	<span class="comment">// set up parameters</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	sigParams := sig.params <span class="comment">// adjusted for variadic functions (may be nil for empty parameter lists!)</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	adjusted := false       <span class="comment">// indicates if sigParams is different from sig.params</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	if sig.variadic {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		if ddd {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			<span class="comment">// variadic_func(a, b, c...)</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>			if len(call.Args) == 1 &amp;&amp; nargs &gt; 1 {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>				<span class="comment">// f()... is not permitted if f() is multi-valued</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>				check.errorf(inNode(call, call.Ellipsis), InvalidDotDotDot, &#34;cannot use ... with %d-valued %s&#34;, nargs, call.Args[0])
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>				return
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>			}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		} else {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>			<span class="comment">// variadic_func(a, b, c)</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			if nargs &gt;= npars-1 {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>				<span class="comment">// Create custom parameters for arguments: keep</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>				<span class="comment">// the first npars-1 parameters and add one for</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>				<span class="comment">// each argument mapping to the ... parameter.</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>				vars := make([]*Var, npars-1) <span class="comment">// npars &gt; 0 for variadic functions</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>				copy(vars, sig.params.vars)
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>				last := sig.params.vars[npars-1]
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>				typ := last.typ.(*Slice).elem
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>				for len(vars) &lt; nargs {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>					vars = append(vars, NewParam(last.pos, last.pkg, last.name, typ))
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>				}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>				sigParams = NewTuple(vars...) <span class="comment">// possibly nil!</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>				adjusted = true
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>				npars = nargs
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			} else {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>				<span class="comment">// nargs &lt; npars-1</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>				npars-- <span class="comment">// for correct error message below</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			}
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	} else {
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		if ddd {
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>			<span class="comment">// standard_func(a, b, c...)</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>			check.errorf(inNode(call, call.Ellipsis), NonVariadicDotDotDot, &#34;cannot use ... in call to non-variadic %s&#34;, call.Fun)
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			return
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		}
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		<span class="comment">// standard_func(a, b, c)</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	<span class="comment">// check argument count</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	if nargs != npars {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		var at positioner = call
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		qualifier := &#34;not enough&#34;
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		if nargs &gt; npars {
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			at = args[npars].expr <span class="comment">// report at first extra argument</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>			qualifier = &#34;too many&#34;
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		} else {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			at = atPos(call.Rparen) <span class="comment">// report at closing )</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		}
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		<span class="comment">// take care of empty parameter lists represented by nil tuples</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		var params []*Var
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		if sig.params != nil {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			params = sig.params.vars
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		err := newErrorf(at, WrongArgCount, &#34;%s arguments in call to %s&#34;, qualifier, call.Fun)
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		err.errorf(nopos, &#34;have %s&#34;, check.typesSummary(operandTypes(args), false))
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		err.errorf(nopos, &#34;want %s&#34;, check.typesSummary(varTypes(params), sig.variadic))
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		check.report(err)
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		return
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	}
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	<span class="comment">// collect type parameters of callee and generic function arguments</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	var tparams []*TypeParam
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	<span class="comment">// collect type parameters of callee</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	n := sig.TypeParams().Len()
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	if n &gt; 0 {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		if !check.allowVersion(check.pkg, call, go1_18) {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>			switch call.Fun.(type) {
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>			case *ast.IndexExpr, *ast.IndexListExpr:
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>				ix := typeparams.UnpackIndexExpr(call.Fun)
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>				check.versionErrorf(inNode(call.Fun, ix.Lbrack), go1_18, &#34;function instantiation&#34;)
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			default:
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>				check.versionErrorf(inNode(call, call.Lparen), go1_18, &#34;implicit function instantiation&#34;)
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>			}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		<span class="comment">// rename type parameters to avoid problems with recursive calls</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		var tmp Type
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		tparams, tmp = check.renameTParams(call.Pos(), sig.TypeParams().list(), sigParams)
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		sigParams = tmp.(*Tuple)
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		<span class="comment">// make sure targs and tparams have the same length</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		for len(targs) &lt; len(tparams) {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>			targs = append(targs, nil)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	assert(len(tparams) == len(targs))
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	<span class="comment">// collect type parameters from generic function arguments</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	var genericArgs []int <span class="comment">// indices of generic function arguments</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	if enableReverseTypeInference {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		for i, arg := range args {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>			<span class="comment">// generic arguments cannot have a defined (*Named) type - no need for underlying type below</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>			if asig, _ := arg.typ.(*Signature); asig != nil &amp;&amp; asig.TypeParams().Len() &gt; 0 {
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>				<span class="comment">// The argument type is a generic function signature. This type is</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>				<span class="comment">// pointer-identical with (it&#39;s copied from) the type of the generic</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>				<span class="comment">// function argument and thus the function object.</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>				<span class="comment">// Before we change the type (type parameter renaming, below), make</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>				<span class="comment">// a clone of it as otherwise we implicitly modify the object&#39;s type</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>				<span class="comment">// (go.dev/issues/63260).</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>				asig = clone(asig)
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>				<span class="comment">// Rename type parameters for cases like f(g, g); this gives each</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>				<span class="comment">// generic function argument a unique type identity (go.dev/issues/59956).</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>				<span class="comment">// TODO(gri) Consider only doing this if a function argument appears</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>				<span class="comment">//           multiple times, which is rare (possible optimization).</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>				atparams, tmp := check.renameTParams(call.Pos(), asig.TypeParams().list(), asig)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>				asig = tmp.(*Signature)
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>				asig.tparams = &amp;TypeParamList{atparams} <span class="comment">// renameTParams doesn&#39;t touch associated type parameters</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>				arg.typ = asig                          <span class="comment">// new type identity for the function argument</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>				tparams = append(tparams, atparams...)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>				<span class="comment">// add partial list of type arguments, if any</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>				if i &lt; len(atargs) {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>					targs = append(targs, atargs[i]...)
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>				}
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>				<span class="comment">// make sure targs and tparams have the same length</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>				for len(targs) &lt; len(tparams) {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>					targs = append(targs, nil)
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>				}
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>				genericArgs = append(genericArgs, i)
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>			}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	}
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	assert(len(tparams) == len(targs))
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	<span class="comment">// at the moment we only support implicit instantiations of argument functions</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	_ = len(genericArgs) &gt; 0 &amp;&amp; check.verifyVersionf(args[genericArgs[0]], go1_21, &#34;implicitly instantiated function as argument&#34;)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	<span class="comment">// tparams holds the type parameters of the callee and generic function arguments, if any:</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	<span class="comment">// the first n type parameters belong to the callee, followed by mi type parameters for each</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	<span class="comment">// of the generic function arguments, where mi = args[i].typ.(*Signature).TypeParams().Len().</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	<span class="comment">// infer missing type arguments of callee and function arguments</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	if len(tparams) &gt; 0 {
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		targs = check.infer(call, tparams, targs, sigParams, args, false)
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		if targs == nil {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) If infer inferred the first targs[:n], consider instantiating</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>			<span class="comment">//           the call signature for better error messages/gopls behavior.</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>			<span class="comment">//           Perhaps instantiate as much as we can, also for arguments.</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			<span class="comment">//           This will require changes to how infer returns its results.</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			return <span class="comment">// error already reported</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		<span class="comment">// update result signature: instantiate if needed</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		if n &gt; 0 {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>			rsig = check.instantiateSignature(call.Pos(), call.Fun, sig, targs[:n], xlist)
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			<span class="comment">// If the callee&#39;s parameter list was adjusted we need to update (instantiate)</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			<span class="comment">// it separately. Otherwise we can simply use the result signature&#39;s parameter</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>			<span class="comment">// list.</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>			if adjusted {
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>				sigParams = check.subst(call.Pos(), sigParams, makeSubstMap(tparams[:n], targs[:n]), nil, check.context()).(*Tuple)
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			} else {
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>				sigParams = rsig.params
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>			}
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		}
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		<span class="comment">// compute argument signatures: instantiate if needed</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		j := n
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		for _, i := range genericArgs {
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			arg := args[i]
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>			asig := arg.typ.(*Signature)
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>			k := j + asig.TypeParams().Len()
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>			<span class="comment">// targs[j:k] are the inferred type arguments for asig</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>			arg.typ = check.instantiateSignature(call.Pos(), arg.expr, asig, targs[j:k], nil) <span class="comment">// TODO(gri) provide xlist if possible (partial instantiations)</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>			check.record(arg)                                                                 <span class="comment">// record here because we didn&#39;t use the usual expr evaluators</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>			j = k
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	}
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	<span class="comment">// check arguments</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	if len(args) &gt; 0 {
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		context := check.sprintf(&#34;argument to %s&#34;, call.Fun)
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>		for i, a := range args {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>			check.assignment(a, sigParams.vars[i].typ, context)
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		}
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	}
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	return
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>}
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>var cgoPrefixes = [...]string{
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	&#34;_Ciconst_&#34;,
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	&#34;_Cfconst_&#34;,
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	&#34;_Csconst_&#34;,
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	&#34;_Ctype_&#34;,
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	&#34;_Cvar_&#34;, <span class="comment">// actually a pointer to the var</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	&#34;_Cfpvar_fp_&#34;,
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	&#34;_Cfunc_&#34;,
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	&#34;_Cmacro_&#34;, <span class="comment">// function to evaluate the expanded expression</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>func (check *Checker) selector(x *operand, e *ast.SelectorExpr, def *TypeName, wantType bool) {
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	<span class="comment">// these must be declared before the &#34;goto Error&#34; statements</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	var (
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>		obj      Object
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		index    []int
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		indirect bool
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	)
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	sel := e.Sel.Name
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	<span class="comment">// If the identifier refers to a package, handle everything here</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	<span class="comment">// so we don&#39;t need a &#34;package&#34; mode for operands: package names</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	<span class="comment">// can only appear in qualified identifiers which are mapped to</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	<span class="comment">// selector expressions.</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	if ident, ok := e.X.(*ast.Ident); ok {
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		obj := check.lookup(ident.Name)
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		if pname, _ := obj.(*PkgName); pname != nil {
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>			assert(pname.pkg == check.pkg)
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>			check.recordUse(ident, pname)
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>			pname.used = true
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>			pkg := pname.imported
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>			var exp Object
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>			funcMode := value
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>			if pkg.cgo {
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>				<span class="comment">// cgo special cases C.malloc: it&#39;s</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>				<span class="comment">// rewritten to _CMalloc and does not</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>				<span class="comment">// support two-result calls.</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>				if sel == &#34;malloc&#34; {
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>					sel = &#34;_CMalloc&#34;
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>				} else {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>					funcMode = cgofunc
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>				}
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>				for _, prefix := range cgoPrefixes {
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>					<span class="comment">// cgo objects are part of the current package (in file</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>					<span class="comment">// _cgo_gotypes.go). Use regular lookup.</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>					_, exp = check.scope.LookupParent(prefix+sel, check.pos)
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>					if exp != nil {
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>						break
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>					}
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>				}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>				if exp == nil {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>					check.errorf(e.Sel, UndeclaredImportedName, &#34;undefined: %s&#34;, ast.Expr(e)) <span class="comment">// cast to ast.Expr to silence vet</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>					goto Error
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>				}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>				check.objDecl(exp, nil)
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>			} else {
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>				exp = pkg.scope.Lookup(sel)
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>				if exp == nil {
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>					if !pkg.fake {
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>						check.errorf(e.Sel, UndeclaredImportedName, &#34;undefined: %s&#34;, ast.Expr(e))
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>					}
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>					goto Error
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>				}
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>				if !exp.Exported() {
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>					check.errorf(e.Sel, UnexportedName, &#34;%s not exported by package %s&#34;, sel, pkg.name)
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>					<span class="comment">// ok to continue</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>				}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>			}
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>			check.recordUse(e.Sel, exp)
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>			<span class="comment">// Simplified version of the code for *ast.Idents:</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>			<span class="comment">// - imported objects are always fully initialized</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>			switch exp := exp.(type) {
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>			case *Const:
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>				assert(exp.Val() != nil)
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>				x.mode = constant_
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>				x.typ = exp.typ
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>				x.val = exp.val
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>			case *TypeName:
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>				x.mode = typexpr
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>				x.typ = exp.typ
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>			case *Var:
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>				x.mode = variable
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>				x.typ = exp.typ
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>				if pkg.cgo &amp;&amp; strings.HasPrefix(exp.name, &#34;_Cvar_&#34;) {
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>					x.typ = x.typ.(*Pointer).base
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>				}
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>			case *Func:
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>				x.mode = funcMode
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>				x.typ = exp.typ
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>				if pkg.cgo &amp;&amp; strings.HasPrefix(exp.name, &#34;_Cmacro_&#34;) {
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>					x.mode = value
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>					x.typ = x.typ.(*Signature).results.vars[0].typ
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>				}
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>			case *Builtin:
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>				x.mode = builtin
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>				x.typ = exp.typ
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>				x.id = exp.id
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>			default:
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>				check.dump(&#34;%v: unexpected object %v&#34;, e.Sel.Pos(), exp)
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>				unreachable()
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>			}
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			x.expr = e
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>			return
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		}
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	}
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	check.exprOrType(x, e.X, false)
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	switch x.mode {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	case typexpr:
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>		<span class="comment">// don&#39;t crash for &#34;type T T.x&#34; (was go.dev/issue/51509)</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		if def != nil &amp;&amp; def.typ == x.typ {
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>			check.cycleError([]Object{def})
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>			goto Error
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		}
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>	case builtin:
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		<span class="comment">// types2 uses the position of &#39;.&#39; for the error</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		check.errorf(e.Sel, UncalledBuiltin, &#34;cannot select on %s&#34;, x)
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		goto Error
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	case invalid:
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>		goto Error
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	}
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	<span class="comment">// Avoid crashing when checking an invalid selector in a method declaration</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	<span class="comment">// (i.e., where def is not set):</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>	<span class="comment">//   type S[T any] struct{}</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	<span class="comment">//   type V = S[any]</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	<span class="comment">//   func (fs *S[T]) M(x V.M) {}</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	<span class="comment">// All codepaths below return a non-type expression. If we get here while</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	<span class="comment">// expecting a type expression, it is an error.</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	<span class="comment">// See go.dev/issue/57522 for more details.</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	<span class="comment">// TODO(rfindley): We should do better by refusing to check selectors in all cases where</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	<span class="comment">// x.typ is incomplete.</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>	if wantType {
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>		check.errorf(e.Sel, NotAType, &#34;%s is not a type&#34;, ast.Expr(e))
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		goto Error
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	}
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	obj, index, indirect = LookupFieldOrMethod(x.typ, x.mode == variable, check.pkg, sel)
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	if obj == nil {
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		<span class="comment">// Don&#39;t report another error if the underlying type was invalid (go.dev/issue/49541).</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>		if !isValid(under(x.typ)) {
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>			goto Error
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>		}
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		if index != nil {
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) should provide actual type where the conflict happens</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>			check.errorf(e.Sel, AmbiguousSelector, &#34;ambiguous selector %s.%s&#34;, x.expr, sel)
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>			goto Error
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>		}
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>		if indirect {
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>			if x.mode == typexpr {
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>				check.errorf(e.Sel, InvalidMethodExpr, &#34;invalid method expression %s.%s (needs pointer receiver (*%s).%s)&#34;, x.typ, sel, x.typ, sel)
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>			} else {
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>				check.errorf(e.Sel, InvalidMethodExpr, &#34;cannot call pointer method %s on %s&#34;, sel, x.typ)
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>			}
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>			goto Error
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		}
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		var why string
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		if isInterfacePtr(x.typ) {
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>			why = check.interfacePtrError(x.typ)
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		} else {
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>			why = check.sprintf(&#34;type %s has no field or method %s&#34;, x.typ, sel)
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>			<span class="comment">// Check if capitalization of sel matters and provide better error message in that case.</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) This code only looks at the first character but LookupFieldOrMethod should</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>			<span class="comment">//           have an (internal) mechanism for case-insensitive lookup that we should use</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>			<span class="comment">//           instead (see types2).</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>			if len(sel) &gt; 0 {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>				var changeCase string
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>				if r := rune(sel[0]); unicode.IsUpper(r) {
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>					changeCase = string(unicode.ToLower(r)) + sel[1:]
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>				} else {
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>					changeCase = string(unicode.ToUpper(r)) + sel[1:]
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>				}
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>				if obj, _, _ = LookupFieldOrMethod(x.typ, x.mode == variable, check.pkg, changeCase); obj != nil {
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>					why += &#34;, but does have &#34; + changeCase
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>				}
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			}
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		}
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		check.errorf(e.Sel, MissingFieldOrMethod, &#34;%s.%s undefined (%s)&#34;, x.expr, sel, why)
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>		goto Error
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	}
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	<span class="comment">// methods may not have a fully set up signature yet</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	if m, _ := obj.(*Func); m != nil {
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>		check.objDecl(m, nil)
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	}
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	if x.mode == typexpr {
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>		<span class="comment">// method expression</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		m, _ := obj.(*Func)
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		if m == nil {
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) should check if capitalization of sel matters and provide better error message in that case</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>			check.errorf(e.Sel, MissingFieldOrMethod, &#34;%s.%s undefined (type %s has no method %s)&#34;, x.expr, sel, x.typ, sel)
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>			goto Error
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		}
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		check.recordSelection(e, MethodExpr, x.typ, m, index, indirect)
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		sig := m.typ.(*Signature)
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>		if sig.recv == nil {
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>			check.error(e, InvalidDeclCycle, &#34;illegal cycle in method declaration&#34;)
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>			goto Error
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>		}
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		<span class="comment">// the receiver type becomes the type of the first function</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		<span class="comment">// argument of the method expression&#39;s function type</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>		var params []*Var
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>		if sig.params != nil {
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>			params = sig.params.vars
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>		}
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>		<span class="comment">// Be consistent about named/unnamed parameters. This is not needed</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>		<span class="comment">// for type-checking, but the newly constructed signature may appear</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>		<span class="comment">// in an error message and then have mixed named/unnamed parameters.</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>		<span class="comment">// (An alternative would be to not print parameter names in errors,</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>		<span class="comment">// but it&#39;s useful to see them; this is cheap and method expressions</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>		<span class="comment">// are rare.)</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>		name := &#34;&#34;
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>		if len(params) &gt; 0 &amp;&amp; params[0].name != &#34;&#34; {
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>			<span class="comment">// name needed</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>			name = sig.recv.name
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			if name == &#34;&#34; {
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>				name = &#34;_&#34;
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			}
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>		}
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>		params = append([]*Var{NewVar(sig.recv.pos, sig.recv.pkg, name, x.typ)}, params...)
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>		x.mode = value
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		x.typ = &amp;Signature{
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>			tparams:  sig.tparams,
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>			params:   NewTuple(params...),
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>			results:  sig.results,
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>			variadic: sig.variadic,
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		}
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>		check.addDeclDep(m)
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	} else {
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>		<span class="comment">// regular selector</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>		switch obj := obj.(type) {
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>		case *Var:
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>			check.recordSelection(e, FieldVal, x.typ, obj, index, indirect)
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>			if x.mode == variable || indirect {
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>				x.mode = variable
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>			} else {
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>				x.mode = value
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>			}
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>			x.typ = obj.typ
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>		case *Func:
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) If we needed to take into account the receiver&#39;s</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>			<span class="comment">// addressability, should we report the type &amp;(x.typ) instead?</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>			check.recordSelection(e, MethodVal, x.typ, obj, index, indirect)
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) The verification pass below is disabled for now because</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>			<span class="comment">//           method sets don&#39;t match method lookup in some cases.</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>			<span class="comment">//           For instance, if we made a copy above when creating a</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>			<span class="comment">//           custom method for a parameterized received type, the</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>			<span class="comment">//           method set method doesn&#39;t match (no copy there). There</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>			<span class="comment">///          may be other situations.</span>
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>			disabled := true
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>			if !disabled &amp;&amp; debug {
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>				<span class="comment">// Verify that LookupFieldOrMethod and MethodSet.Lookup agree.</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>				<span class="comment">// TODO(gri) This only works because we call LookupFieldOrMethod</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>				<span class="comment">// _before_ calling NewMethodSet: LookupFieldOrMethod completes</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>				<span class="comment">// any incomplete interfaces so they are available to NewMethodSet</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>				<span class="comment">// (which assumes that interfaces have been completed already).</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>				typ := x.typ
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>				if x.mode == variable {
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>					<span class="comment">// If typ is not an (unnamed) pointer or an interface,</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>					<span class="comment">// use *typ instead, because the method set of *typ</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>					<span class="comment">// includes the methods of typ.</span>
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>					<span class="comment">// Variables are addressable, so we can always take their</span>
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>					<span class="comment">// address.</span>
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>					if _, ok := typ.(*Pointer); !ok &amp;&amp; !IsInterface(typ) {
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>						typ = &amp;Pointer{base: typ}
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>					}
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>				}
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>				<span class="comment">// If we created a synthetic pointer type above, we will throw</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>				<span class="comment">// away the method set computed here after use.</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>				<span class="comment">// TODO(gri) Method set computation should probably always compute</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>				<span class="comment">// both, the value and the pointer receiver method set and represent</span>
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>				<span class="comment">// them in a single structure.</span>
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>				<span class="comment">// TODO(gri) Consider also using a method set cache for the lifetime</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>				<span class="comment">// of checker once we rely on MethodSet lookup instead of individual</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>				<span class="comment">// lookup.</span>
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>				mset := NewMethodSet(typ)
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>				if m := mset.Lookup(check.pkg, sel); m == nil || m.obj != obj {
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>					check.dump(&#34;%v: (%s).%v -&gt; %s&#34;, e.Pos(), typ, obj.name, m)
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>					check.dump(&#34;%s\n&#34;, mset)
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>					<span class="comment">// Caution: MethodSets are supposed to be used externally</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>					<span class="comment">// only (after all interface types were completed). It&#39;s</span>
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>					<span class="comment">// now possible that we get here incorrectly. Not urgent</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>					<span class="comment">// to fix since we only run this code in debug mode.</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>					<span class="comment">// TODO(gri) fix this eventually.</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>					panic(&#34;method sets and lookup don&#39;t agree&#34;)
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>				}
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>			}
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>			x.mode = value
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>			<span class="comment">// remove receiver</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>			sig := *obj.typ.(*Signature)
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>			sig.recv = nil
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>			x.typ = &amp;sig
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>			check.addDeclDep(obj)
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>		default:
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>			unreachable()
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>		}
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	}
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	<span class="comment">// everything went well</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	x.expr = e
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	return
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>Error:
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	x.mode = invalid
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>	x.expr = e
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>}
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span><span class="comment">// use type-checks each argument.</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span><span class="comment">// Useful to make sure expressions are evaluated</span>
<span id="L991" class="ln">   991&nbsp;&nbsp;</span><span class="comment">// (and variables are &#34;used&#34;) in the presence of</span>
<span id="L992" class="ln">   992&nbsp;&nbsp;</span><span class="comment">// other errors. Arguments may be nil.</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span><span class="comment">// Reports if all arguments evaluated without error.</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>func (check *Checker) use(args ...ast.Expr) bool { return check.useN(args, false) }
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span><span class="comment">// useLHS is like use, but doesn&#39;t &#34;use&#34; top-level identifiers.</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span><span class="comment">// It should be called instead of use if the arguments are</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span><span class="comment">// expressions on the lhs of an assignment.</span>
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>func (check *Checker) useLHS(args ...ast.Expr) bool { return check.useN(args, true) }
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>func (check *Checker) useN(args []ast.Expr, lhs bool) bool {
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>	ok := true
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>	for _, e := range args {
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>		if !check.use1(e, lhs) {
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>			ok = false
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>		}
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	}
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	return ok
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>}
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>func (check *Checker) use1(e ast.Expr, lhs bool) bool {
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	var x operand
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	x.mode = value <span class="comment">// anything but invalid</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	switch n := unparen(e).(type) {
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	case nil:
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>		<span class="comment">// nothing to do</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	case *ast.Ident:
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>		<span class="comment">// don&#39;t report an error evaluating blank</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>		if n.Name == &#34;_&#34; {
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>			break
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		}
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>		<span class="comment">// If the lhs is an identifier denoting a variable v, this assignment</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		<span class="comment">// is not a &#39;use&#39; of v. Remember current value of v.used and restore</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>		<span class="comment">// after evaluating the lhs via check.rawExpr.</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		var v *Var
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>		var v_used bool
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>		if lhs {
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>			if _, obj := check.scope.LookupParent(n.Name, nopos); obj != nil {
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>				<span class="comment">// It&#39;s ok to mark non-local variables, but ignore variables</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>				<span class="comment">// from other packages to avoid potential race conditions with</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>				<span class="comment">// dot-imported variables.</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>				if w, _ := obj.(*Var); w != nil &amp;&amp; w.pkg == check.pkg {
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>					v = w
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>					v_used = v.used
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>				}
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>			}
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>		}
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>		check.exprOrType(&amp;x, n, true)
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>		if v != nil {
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>			v.used = v_used <span class="comment">// restore v.used</span>
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>		}
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>	default:
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>		check.rawExpr(nil, &amp;x, e, nil, true)
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>	}
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>	return x.mode != invalid
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>}
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>
</pre><p><a href="call.go?m=text">View as plain text</a></p>

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
