<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/assignments.go - Go Documentation Server</title>

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
<a href="assignments.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">assignments.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements initialization and assignment checks.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package types
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// assignment reports whether x can be assigned to a variable of type T,</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// if necessary by attempting to convert untyped values to the appropriate</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// type. context describes the context in which the assignment takes place.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// Use T == nil to indicate assignment to an untyped blank identifier.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// If the assignment check fails, x.mode is set to invalid.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>func (check *Checker) assignment(x *operand, T Type, context string) {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	check.singleValue(x)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	switch x.mode {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	case invalid:
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		return <span class="comment">// error reported before</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	case constant_, variable, mapindex, value, commaok, commaerr:
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		<span class="comment">// ok</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	default:
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		<span class="comment">// we may get here because of other problems (go.dev/issue/39634, crash 12)</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) do we need a new &#34;generic&#34; error code here?</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		check.errorf(x, IncompatibleAssign, &#34;cannot assign %s to %s in %s&#34;, x, T, context)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		return
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	if isUntyped(x.typ) {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		target := T
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;If an untyped constant is assigned to a variable of interface</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		<span class="comment">// type or the blank identifier, the constant is first converted to type</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		<span class="comment">// bool, rune, int, float64, complex128 or string respectively, depending</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		<span class="comment">// on whether the value is a boolean, rune, integer, floating-point,</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		<span class="comment">// complex, or string constant.&#34;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		if T == nil || isNonTypeParamInterface(T) {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>			if T == nil &amp;&amp; x.typ == Typ[UntypedNil] {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>				check.errorf(x, UntypedNilUse, &#34;use of untyped nil in %s&#34;, context)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>				x.mode = invalid
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>				return
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			target = Default(x.typ)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		newType, val, code := check.implicitTypeAndValue(x, target)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		if code != 0 {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			msg := check.sprintf(&#34;cannot use %s as %s value in %s&#34;, x, target, context)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>			switch code {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			case TruncatedFloat:
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>				msg += &#34; (truncated)&#34;
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			case NumericOverflow:
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>				msg += &#34; (overflows)&#34;
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			default:
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>				code = IncompatibleAssign
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			check.error(x, code, msg)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			return
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		if val != nil {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			x.val = val
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			check.updateExprVal(x.expr, val)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		if newType != x.typ {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			x.typ = newType
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			check.updateExprType(x.expr, newType, false)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// x.typ is typed</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// A generic (non-instantiated) function value cannot be assigned to a variable.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	if sig, _ := under(x.typ).(*Signature); sig != nil &amp;&amp; sig.TypeParams().Len() &gt; 0 {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		check.errorf(x, WrongTypeArgCount, &#34;cannot use generic function %s without instantiation in %s&#34;, x, context)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		return
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;If a left-hand side is the blank identifier, any typed or</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// non-constant value except for the predeclared identifier nil may</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// be assigned to it.&#34;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	if T == nil {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		return
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	cause := &#34;&#34;
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	if ok, code := x.assignableTo(check, T, &amp;cause); !ok {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		if cause != &#34;&#34; {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			check.errorf(x, code, &#34;cannot use %s as %s value in %s: %s&#34;, x, T, context, cause)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		} else {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			check.errorf(x, code, &#34;cannot use %s as %s value in %s&#34;, x, T, context)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>func (check *Checker) initConst(lhs *Const, x *operand) {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	if x.mode == invalid || !isValid(x.typ) || !isValid(lhs.typ) {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		if lhs.typ == nil {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			lhs.typ = Typ[Invalid]
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		return
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// rhs must be a constant</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	if x.mode != constant_ {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		check.errorf(x, InvalidConstInit, &#34;%s is not constant&#34;, x)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		if lhs.typ == nil {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			lhs.typ = Typ[Invalid]
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		return
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	assert(isConstType(x.typ))
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// If the lhs doesn&#39;t have a type yet, use the type of x.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	if lhs.typ == nil {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		lhs.typ = x.typ
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	check.assignment(x, lhs.typ, &#34;constant declaration&#34;)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if x.mode == invalid {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		return
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	lhs.val = x.val
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">// initVar checks the initialization lhs = x in a variable declaration.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">// If lhs doesn&#39;t have a type yet, it is given the type of x,</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// or Typ[Invalid] in case of an error.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// If the initialization check fails, x.mode is set to invalid.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>func (check *Checker) initVar(lhs *Var, x *operand, context string) {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	if x.mode == invalid || !isValid(x.typ) || !isValid(lhs.typ) {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		if lhs.typ == nil {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			lhs.typ = Typ[Invalid]
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		return
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">// If lhs doesn&#39;t have a type yet, use the type of x.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	if lhs.typ == nil {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		typ := x.typ
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		if isUntyped(typ) {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			<span class="comment">// convert untyped types to default types</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			if typ == Typ[UntypedNil] {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>				check.errorf(x, UntypedNilUse, &#34;use of untyped nil in %s&#34;, context)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>				lhs.typ = Typ[Invalid]
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>				x.mode = invalid
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>				return
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			typ = Default(typ)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		lhs.typ = typ
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	check.assignment(x, lhs.typ, context)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// lhsVar checks a lhs variable in an assignment and returns its type.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// lhsVar takes care of not counting a lhs identifier as a &#34;use&#34; of</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span><span class="comment">// that identifier. The result is nil if it is the blank identifier,</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// and Typ[Invalid] if it is an invalid lhs expression.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>func (check *Checker) lhsVar(lhs ast.Expr) Type {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// Determine if the lhs is a (possibly parenthesized) identifier.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	ident, _ := unparen(lhs).(*ast.Ident)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t evaluate lhs if it is the blank identifier.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	if ident != nil &amp;&amp; ident.Name == &#34;_&#34; {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		check.recordDef(ident, nil)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		return nil
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// If the lhs is an identifier denoting a variable v, this reference</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// is not a &#39;use&#39; of v. Remember current value of v.used and restore</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// after evaluating the lhs via check.expr.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	var v *Var
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	var v_used bool
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if ident != nil {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		if obj := check.lookup(ident.Name); obj != nil {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			<span class="comment">// It&#39;s ok to mark non-local variables, but ignore variables</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			<span class="comment">// from other packages to avoid potential race conditions with</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			<span class="comment">// dot-imported variables.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			if w, _ := obj.(*Var); w != nil &amp;&amp; w.pkg == check.pkg {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>				v = w
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>				v_used = v.used
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	var x operand
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	check.expr(nil, &amp;x, lhs)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	if v != nil {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		v.used = v_used <span class="comment">// restore v.used</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	if x.mode == invalid || !isValid(x.typ) {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		return Typ[Invalid]
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;Each left-hand side operand must be addressable, a map index</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">// expression, or the blank identifier. Operands may be parenthesized.&#34;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	switch x.mode {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	case invalid:
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		return Typ[Invalid]
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	case variable, mapindex:
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		<span class="comment">// ok</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	default:
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		if sel, ok := x.expr.(*ast.SelectorExpr); ok {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			var op operand
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			check.expr(nil, &amp;op, sel.X)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			if op.mode == mapindex {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>				check.errorf(&amp;x, UnaddressableFieldAssign, &#34;cannot assign to struct field %s in map&#34;, ExprString(x.expr))
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>				return Typ[Invalid]
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		check.errorf(&amp;x, UnassignableOperand, &#34;cannot assign to %s (neither addressable nor a map index expression)&#34;, x.expr)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		return Typ[Invalid]
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	return x.typ
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// assignVar checks the assignment lhs = rhs (if x == nil), or lhs = x (if x != nil).</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// If x != nil, it must be the evaluation of rhs (and rhs will be ignored).</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">// If the assignment check fails and x != nil, x.mode is set to invalid.</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>func (check *Checker) assignVar(lhs, rhs ast.Expr, x *operand, context string) {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	T := check.lhsVar(lhs) <span class="comment">// nil if lhs is _</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	if !isValid(T) {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		if x != nil {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		} else {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			check.use(rhs)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		return
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	if x == nil {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		var target *target
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		<span class="comment">// avoid calling ExprString if not needed</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		if T != nil {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			if _, ok := under(T).(*Signature); ok {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>				target = newTarget(T, ExprString(lhs))
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		x = new(operand)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		check.expr(target, x, rhs)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	if T == nil &amp;&amp; context == &#34;assignment&#34; {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		context = &#34;assignment to _ identifier&#34;
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	check.assignment(x, T, context)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">// operandTypes returns the list of types for the given operands.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>func operandTypes(list []*operand) (res []Type) {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	for _, x := range list {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		res = append(res, x.typ)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	return res
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span><span class="comment">// varTypes returns the list of types for the given variables.</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>func varTypes(list []*Var) (res []Type) {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	for _, x := range list {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		res = append(res, x.typ)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	return res
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">// typesSummary returns a string of the form &#34;(t1, t2, ...)&#34; where the</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">// ti&#39;s are user-friendly string representations for the given types.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// If variadic is set and the last type is a slice, its string is of</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">// the form &#34;...E&#34; where E is the slice&#39;s element type.</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>func (check *Checker) typesSummary(list []Type, variadic bool) string {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	var res []string
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	for i, t := range list {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		var s string
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		switch {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		case t == nil:
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			fallthrough <span class="comment">// should not happen but be cautious</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		case !isValid(t):
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			s = &#34;unknown type&#34;
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		case isUntyped(t):
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			if isNumeric(t) {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>				<span class="comment">// Do not imply a specific type requirement:</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>				<span class="comment">// &#34;have number, want float64&#34; is better than</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>				<span class="comment">// &#34;have untyped int, want float64&#34; or</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>				<span class="comment">// &#34;have int, want float64&#34;.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>				s = &#34;number&#34;
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			} else {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>				<span class="comment">// If we don&#39;t have a number, omit the &#34;untyped&#34; qualifier</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>				<span class="comment">// for compactness.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>				s = strings.Replace(t.(*Basic).name, &#34;untyped &#34;, &#34;&#34;, -1)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		case variadic &amp;&amp; i == len(list)-1:
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			s = check.sprintf(&#34;...%s&#34;, t.(*Slice).elem)
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		if s == &#34;&#34; {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>			s = check.sprintf(&#34;%s&#34;, t)
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		res = append(res, s)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	}
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	return &#34;(&#34; + strings.Join(res, &#34;, &#34;) + &#34;)&#34;
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>func measure(x int, unit string) string {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	if x != 1 {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		unit += &#34;s&#34;
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;%d %s&#34;, x, unit)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>func (check *Checker) assignError(rhs []ast.Expr, l, r int) {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	vars := measure(l, &#34;variable&#34;)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	vals := measure(r, &#34;value&#34;)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	rhs0 := rhs[0]
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	if len(rhs) == 1 {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		if call, _ := unparen(rhs0).(*ast.CallExpr); call != nil {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>			check.errorf(rhs0, WrongAssignCount, &#34;assignment mismatch: %s but %s returns %s&#34;, vars, call.Fun, vals)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>			return
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	check.errorf(rhs0, WrongAssignCount, &#34;assignment mismatch: %s but %s&#34;, vars, vals)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func (check *Checker) returnError(at positioner, lhs []*Var, rhs []*operand) {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	l, r := len(lhs), len(rhs)
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	qualifier := &#34;not enough&#34;
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	if r &gt; l {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		at = rhs[l] <span class="comment">// report at first extra value</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		qualifier = &#34;too many&#34;
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	} else if r &gt; 0 {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		at = rhs[r-1] <span class="comment">// report at last value</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	var err error_
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	err.code = WrongResultCount
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	err.errorf(at.Pos(), &#34;%s return values&#34;, qualifier)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	err.errorf(nopos, &#34;have %s&#34;, check.typesSummary(operandTypes(rhs), false))
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	err.errorf(nopos, &#34;want %s&#34;, check.typesSummary(varTypes(lhs), false))
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	check.report(&amp;err)
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// initVars type-checks assignments of initialization expressions orig_rhs</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// to variables lhs.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// If returnStmt is non-nil, initVars type-checks the implicit assignment</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// of result expressions orig_rhs to function result parameters lhs.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>func (check *Checker) initVars(lhs []*Var, orig_rhs []ast.Expr, returnStmt ast.Stmt) {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	context := &#34;assignment&#34;
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	if returnStmt != nil {
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		context = &#34;return statement&#34;
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	l, r := len(lhs), len(orig_rhs)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	<span class="comment">// If l == 1 and the rhs is a single call, for a better</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	<span class="comment">// error message don&#39;t handle it as n:n mapping below.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	isCall := false
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	if r == 1 {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		_, isCall = unparen(orig_rhs[0]).(*ast.CallExpr)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	<span class="comment">// If we have a n:n mapping from lhs variable to rhs expression,</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	<span class="comment">// each value can be assigned to its corresponding variable.</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	if l == r &amp;&amp; !isCall {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		var x operand
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		for i, lhs := range lhs {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			desc := lhs.name
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			if returnStmt != nil &amp;&amp; desc == &#34;&#34; {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>				desc = &#34;result variable&#34;
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			check.expr(newTarget(lhs.typ, desc), &amp;x, orig_rhs[i])
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			check.initVar(lhs, &amp;x, context)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		}
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		return
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	<span class="comment">// If we don&#39;t have an n:n mapping, the rhs must be a single expression</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	<span class="comment">// resulting in 2 or more values; otherwise we have an assignment mismatch.</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	if r != 1 {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		<span class="comment">// Only report a mismatch error if there are no other errors on the rhs.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		if check.use(orig_rhs...) {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			if returnStmt != nil {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>				rhs := check.exprList(orig_rhs)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>				check.returnError(returnStmt, lhs, rhs)
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			} else {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>				check.assignError(orig_rhs, l, r)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			}
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		<span class="comment">// ensure that LHS variables have a type</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		for _, v := range lhs {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			if v.typ == nil {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>				v.typ = Typ[Invalid]
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			}
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		return
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	rhs, commaOk := check.multiExpr(orig_rhs[0], l == 2 &amp;&amp; returnStmt == nil)
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	r = len(rhs)
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	if l == r {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		for i, lhs := range lhs {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			check.initVar(lhs, rhs[i], context)
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		<span class="comment">// Only record comma-ok expression if both initializations succeeded</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		<span class="comment">// (go.dev/issue/59371).</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		if commaOk &amp;&amp; rhs[0].mode != invalid &amp;&amp; rhs[1].mode != invalid {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			check.recordCommaOkTypes(orig_rhs[0], rhs)
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		return
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	<span class="comment">// In all other cases we have an assignment mismatch.</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	<span class="comment">// Only report a mismatch error if there are no other errors on the rhs.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	if rhs[0].mode != invalid {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		if returnStmt != nil {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			check.returnError(returnStmt, lhs, rhs)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		} else {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			check.assignError(orig_rhs, l, r)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	<span class="comment">// ensure that LHS variables have a type</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	for _, v := range lhs {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		if v.typ == nil {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			v.typ = Typ[Invalid]
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	}
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	<span class="comment">// orig_rhs[0] was already evaluated</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span><span class="comment">// assignVars type-checks assignments of expressions orig_rhs to variables lhs.</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>func (check *Checker) assignVars(lhs, orig_rhs []ast.Expr) {
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	l, r := len(lhs), len(orig_rhs)
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	<span class="comment">// If l == 1 and the rhs is a single call, for a better</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	<span class="comment">// error message don&#39;t handle it as n:n mapping below.</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	isCall := false
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	if r == 1 {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		_, isCall = unparen(orig_rhs[0]).(*ast.CallExpr)
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	<span class="comment">// If we have a n:n mapping from lhs variable to rhs expression,</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	<span class="comment">// each value can be assigned to its corresponding variable.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	if l == r &amp;&amp; !isCall {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		for i, lhs := range lhs {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			check.assignVar(lhs, orig_rhs[i], nil, &#34;assignment&#34;)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		return
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	<span class="comment">// If we don&#39;t have an n:n mapping, the rhs must be a single expression</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	<span class="comment">// resulting in 2 or more values; otherwise we have an assignment mismatch.</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	if r != 1 {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		<span class="comment">// Only report a mismatch error if there are no other errors on the lhs or rhs.</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		okLHS := check.useLHS(lhs...)
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		okRHS := check.use(orig_rhs...)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		if okLHS &amp;&amp; okRHS {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>			check.assignError(orig_rhs, l, r)
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		return
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	rhs, commaOk := check.multiExpr(orig_rhs[0], l == 2)
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	r = len(rhs)
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	if l == r {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		for i, lhs := range lhs {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			check.assignVar(lhs, nil, rhs[i], &#34;assignment&#34;)
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		<span class="comment">// Only record comma-ok expression if both assignments succeeded</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		<span class="comment">// (go.dev/issue/59371).</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		if commaOk &amp;&amp; rhs[0].mode != invalid &amp;&amp; rhs[1].mode != invalid {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			check.recordCommaOkTypes(orig_rhs[0], rhs)
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		return
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	<span class="comment">// In all other cases we have an assignment mismatch.</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	<span class="comment">// Only report a mismatch error if there are no other errors on the rhs.</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	if rhs[0].mode != invalid {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		check.assignError(orig_rhs, l, r)
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	check.useLHS(lhs...)
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	<span class="comment">// orig_rhs[0] was already evaluated</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>func (check *Checker) shortVarDecl(pos positioner, lhs, rhs []ast.Expr) {
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	top := len(check.delayed)
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	scope := check.scope
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	<span class="comment">// collect lhs variables</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	seen := make(map[string]bool, len(lhs))
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	lhsVars := make([]*Var, len(lhs))
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	newVars := make([]*Var, 0, len(lhs))
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	hasErr := false
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	for i, lhs := range lhs {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		ident, _ := lhs.(*ast.Ident)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		if ident == nil {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>			check.useLHS(lhs)
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			<span class="comment">// TODO(rFindley) this is redundant with a parser error. Consider omitting?</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			check.errorf(lhs, BadDecl, &#34;non-name %s on left side of :=&#34;, lhs)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			hasErr = true
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			continue
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		name := ident.Name
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		if name != &#34;_&#34; {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>			if seen[name] {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>				check.errorf(lhs, RepeatedDecl, &#34;%s repeated on left side of :=&#34;, lhs)
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>				hasErr = true
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>				continue
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			seen[name] = true
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		<span class="comment">// Use the correct obj if the ident is redeclared. The</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		<span class="comment">// variable&#39;s scope starts after the declaration; so we</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		<span class="comment">// must use Scope.Lookup here and call Scope.Insert</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		<span class="comment">// (via check.declare) later.</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		if alt := scope.Lookup(name); alt != nil {
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			check.recordUse(ident, alt)
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			<span class="comment">// redeclared object must be a variable</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>			if obj, _ := alt.(*Var); obj != nil {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>				lhsVars[i] = obj
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>			} else {
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>				check.errorf(lhs, UnassignableOperand, &#34;cannot assign to %s&#34;, lhs)
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>				hasErr = true
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>			continue
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		}
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		<span class="comment">// declare new variable</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		obj := NewVar(ident.Pos(), check.pkg, name, nil)
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		lhsVars[i] = obj
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		if name != &#34;_&#34; {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>			newVars = append(newVars, obj)
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		check.recordDef(ident, obj)
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	<span class="comment">// create dummy variables where the lhs is invalid</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	for i, obj := range lhsVars {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		if obj == nil {
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>			lhsVars[i] = NewVar(lhs[i].Pos(), check.pkg, &#34;_&#34;, nil)
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		}
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	}
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	check.initVars(lhsVars, rhs, nil)
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	<span class="comment">// process function literals in rhs expressions before scope changes</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	check.processDelayed(top)
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	if len(newVars) == 0 &amp;&amp; !hasErr {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		check.softErrorf(pos, NoNewVar, &#34;no new variables on left side of :=&#34;)
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		return
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	}
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	<span class="comment">// declare new variables</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;The scope of a constant or variable identifier declared inside</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	<span class="comment">// a function begins at the end of the ConstSpec or VarSpec (ShortVarDecl</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	<span class="comment">// for short variable declarations) and ends at the end of the innermost</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	<span class="comment">// containing block.&#34;</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	scopePos := rhs[len(rhs)-1].End()
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	for _, obj := range newVars {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		check.declare(scope, nil, obj, scopePos) <span class="comment">// id = nil: recordDef already called</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>}
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>
</pre><p><a href="assignments.go?m=text">View as plain text</a></p>

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
