<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/typexpr.go - Go Documentation Server</title>

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
<a href="typexpr.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">typexpr.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements type-checking of identifiers and type expressions.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package types
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;go/constant&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;go/internal/typeparams&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// ident type-checks identifier e and initializes x with the value or type of e.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// If an error occurred, x.mode is set to invalid.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// For the meaning of def, see Checker.definedType, below.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// If wantType is set, the identifier e is expected to denote a type.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>func (check *Checker) ident(x *operand, e *ast.Ident, def *TypeName, wantType bool) {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	x.mode = invalid
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	x.expr = e
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// Note that we cannot use check.lookup here because the returned scope</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// may be different from obj.Parent(). See also Scope.LookupParent doc.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	scope, obj := check.scope.LookupParent(e.Name, check.pos)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	switch obj {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	case nil:
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		if e.Name == &#34;_&#34; {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>			<span class="comment">// Blank identifiers are never declared, but the current identifier may</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>			<span class="comment">// be a placeholder for a receiver type parameter. In this case we can</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>			<span class="comment">// resolve its type and object from Checker.recvTParamMap.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>			if tpar := check.recvTParamMap[e]; tpar != nil {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>				x.mode = typexpr
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>				x.typ = tpar
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>			} else {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>				check.error(e, InvalidBlank, &#34;cannot use _ as value or type&#34;)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>			}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		} else {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>			check.errorf(e, UndeclaredName, &#34;undefined: %s&#34;, e.Name)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		return
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	case universeAny, universeComparable:
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		if !check.verifyVersionf(e, go1_18, &#34;predeclared %s&#34;, e.Name) {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			return <span class="comment">// avoid follow-on errors</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	check.recordUse(e, obj)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// Type-check the object.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// Only call Checker.objDecl if the object doesn&#39;t have a type yet</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// (in which case we must actually determine it) or the object is a</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// TypeName and we also want a type (in which case we might detect</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// a cycle which needs to be reported). Otherwise we can skip the</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// call and avoid a possible cycle error in favor of the more</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// informative &#34;not a type/value&#34; error that this function&#39;s caller</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// will issue (see go.dev/issue/25790).</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	typ := obj.Type()
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	if _, gotType := obj.(*TypeName); typ == nil || gotType &amp;&amp; wantType {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		check.objDecl(obj, def)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		typ = obj.Type() <span class="comment">// type must have been assigned by Checker.objDecl</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	assert(typ != nil)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// The object may have been dot-imported.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// If so, mark the respective package as used.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// (This code is only needed for dot-imports. Without them,</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// we only have to mark variables, see *Var case below).</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if pkgName := check.dotImportMap[dotImportKey{scope, obj.Name()}]; pkgName != nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		pkgName.used = true
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	switch obj := obj.(type) {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	case *PkgName:
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		check.errorf(e, InvalidPkgUse, &#34;use of package %s not in selector&#34;, obj.name)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		return
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	case *Const:
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		check.addDeclDep(obj)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		if !isValid(typ) {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			return
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		if obj == universeIota {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			if check.iota == nil {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>				check.error(e, InvalidIota, &#34;cannot use iota outside constant declaration&#34;)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>				return
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			x.val = check.iota
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		} else {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			x.val = obj.val
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		assert(x.val != nil)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		x.mode = constant_
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	case *TypeName:
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		if !check.enableAlias &amp;&amp; check.isBrokenAlias(obj) {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			check.errorf(e, InvalidDeclCycle, &#34;invalid use of type alias %s in recursive type (see go.dev/issue/50729)&#34;, obj.name)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			return
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		x.mode = typexpr
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	case *Var:
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		<span class="comment">// It&#39;s ok to mark non-local variables, but ignore variables</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">// from other packages to avoid potential race conditions with</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		<span class="comment">// dot-imported variables.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		if obj.pkg == check.pkg {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			obj.used = true
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		check.addDeclDep(obj)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		if !isValid(typ) {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			return
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		x.mode = variable
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	case *Func:
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		check.addDeclDep(obj)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		x.mode = value
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	case *Builtin:
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		x.id = obj.id
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		x.mode = builtin
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	case *Nil:
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		x.mode = value
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	default:
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		unreachable()
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	x.typ = typ
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">// typ type-checks the type expression e and returns its type, or Typ[Invalid].</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// The type must not be an (uninstantiated) generic type.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>func (check *Checker) typ(e ast.Expr) Type {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	return check.definedType(e, nil)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// varType type-checks the type expression e and returns its type, or Typ[Invalid].</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// The type must not be an (uninstantiated) generic type and it must not be a</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// constraint interface.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>func (check *Checker) varType(e ast.Expr) Type {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	typ := check.definedType(e, nil)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	check.validVarType(e, typ)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	return typ
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// validVarType reports an error if typ is a constraint interface.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// The expression e is used for error reporting, if any.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func (check *Checker) validVarType(e ast.Expr, typ Type) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// If we have a type parameter there&#39;s nothing to do.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	if isTypeParam(typ) {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		return
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// We don&#39;t want to call under() or complete interfaces while we are in</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// the middle of type-checking parameter declarations that might belong</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// to interface methods. Delay this check to the end of type-checking.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	check.later(func() {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		if t, _ := under(typ).(*Interface); t != nil {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			tset := computeInterfaceTypeSet(check, e.Pos(), t) <span class="comment">// TODO(gri) is this the correct position?</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			if !tset.IsMethodSet() {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>				if tset.comparable {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>					check.softErrorf(e, MisplacedConstraintIface, &#34;cannot use type %s outside a type constraint: interface is (or embeds) comparable&#34;, typ)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>				} else {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>					check.softErrorf(e, MisplacedConstraintIface, &#34;cannot use type %s outside a type constraint: interface contains type constraints&#34;, typ)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>				}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}).describef(e, &#34;check var type %s&#34;, typ)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// definedType is like typ but also accepts a type name def.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// If def != nil, e is the type specification for the type named def, declared</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">// in a type declaration, and def.typ.underlying will be set to the type of e</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">// before any components of e are type-checked.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>func (check *Checker) definedType(e ast.Expr, def *TypeName) Type {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	typ := check.typInternal(e, def)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	assert(isTyped(typ))
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	if isGeneric(typ) {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		check.errorf(e, WrongTypeArgCount, &#34;cannot use generic type %s without instantiation&#34;, typ)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		typ = Typ[Invalid]
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	check.recordTypeAndValue(e, typexpr, typ, nil)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	return typ
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// genericType is like typ but the type must be an (uninstantiated) generic</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// type. If cause is non-nil and the type expression was a valid type but not</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// generic, cause will be populated with a message describing the error.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>func (check *Checker) genericType(e ast.Expr, cause *string) Type {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	typ := check.typInternal(e, nil)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	assert(isTyped(typ))
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	if isValid(typ) &amp;&amp; !isGeneric(typ) {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		if cause != nil {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			*cause = check.sprintf(&#34;%s is not a generic type&#34;, typ)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		typ = Typ[Invalid]
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// TODO(gri) what is the correct call below?</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	check.recordTypeAndValue(e, typexpr, typ, nil)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	return typ
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// goTypeName returns the Go type name for typ and</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// removes any occurrences of &#34;types.&#34; from that name.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>func goTypeName(typ Type) string {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	return strings.ReplaceAll(fmt.Sprintf(&#34;%T&#34;, typ), &#34;types.&#34;, &#34;&#34;)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// typInternal drives type checking of types.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// Must only be called by definedType or genericType.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>func (check *Checker) typInternal(e0 ast.Expr, def *TypeName) (T Type) {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	if check.conf._Trace {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		check.trace(e0.Pos(), &#34;-- type %s&#34;, e0)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		check.indent++
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		defer func() {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			check.indent--
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			var under Type
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			if T != nil {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>				<span class="comment">// Calling under() here may lead to endless instantiations.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				<span class="comment">// Test case: type T[P any] *T[P]</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>				under = safeUnderlying(T)
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			if T == under {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>				check.trace(e0.Pos(), &#34;=&gt; %s // %s&#34;, T, goTypeName(T))
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			} else {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>				check.trace(e0.Pos(), &#34;=&gt; %s (under = %s) // %s&#34;, T, under, goTypeName(T))
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		}()
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	switch e := e0.(type) {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	case *ast.BadExpr:
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		<span class="comment">// ignore - error reported before</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	case *ast.Ident:
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		var x operand
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		check.ident(&amp;x, e, def, true)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		switch x.mode {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		case typexpr:
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			typ := x.typ
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			setDefType(def, typ)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			return typ
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		case invalid:
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			<span class="comment">// ignore - error reported before</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		case novalue:
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			check.errorf(&amp;x, NotAType, &#34;%s used as type&#34;, &amp;x)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		default:
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			check.errorf(&amp;x, NotAType, &#34;%s is not a type&#34;, &amp;x)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	case *ast.SelectorExpr:
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		var x operand
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		check.selector(&amp;x, e, def, true)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		switch x.mode {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		case typexpr:
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			typ := x.typ
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			setDefType(def, typ)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			return typ
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		case invalid:
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			<span class="comment">// ignore - error reported before</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		case novalue:
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			check.errorf(&amp;x, NotAType, &#34;%s used as type&#34;, &amp;x)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		default:
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			check.errorf(&amp;x, NotAType, &#34;%s is not a type&#34;, &amp;x)
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	case *ast.IndexExpr, *ast.IndexListExpr:
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		ix := typeparams.UnpackIndexExpr(e)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		check.verifyVersionf(inNode(e, ix.Lbrack), go1_18, &#34;type instantiation&#34;)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		return check.instantiatedType(ix, def)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	case *ast.ParenExpr:
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		<span class="comment">// Generic types must be instantiated before they can be used in any form.</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		<span class="comment">// Consequently, generic types cannot be parenthesized.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		return check.definedType(e.X, def)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	case *ast.ArrayType:
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		if e.Len == nil {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			typ := new(Slice)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			setDefType(def, typ)
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			typ.elem = check.varType(e.Elt)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			return typ
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		typ := new(Array)
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		setDefType(def, typ)
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		<span class="comment">// Provide a more specific error when encountering a [...] array</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		<span class="comment">// rather than leaving it to the handling of the ... expression.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		if _, ok := e.Len.(*ast.Ellipsis); ok {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			check.error(e.Len, BadDotDotDotSyntax, &#34;invalid use of [...] array (outside a composite literal)&#34;)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			typ.len = -1
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		} else {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			typ.len = check.arrayLength(e.Len)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		typ.elem = check.varType(e.Elt)
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		if typ.len &gt;= 0 {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			return typ
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		<span class="comment">// report error if we encountered [...]</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	case *ast.Ellipsis:
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		<span class="comment">// dots are handled explicitly where they are legal</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		<span class="comment">// (array composite literals and parameter lists)</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		check.error(e, InvalidDotDotDot, &#34;invalid use of &#39;...&#39;&#34;)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		check.use(e.Elt)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	case *ast.StructType:
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		typ := new(Struct)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		setDefType(def, typ)
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		check.structType(typ, e)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		return typ
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	case *ast.StarExpr:
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		typ := new(Pointer)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		typ.base = Typ[Invalid] <span class="comment">// avoid nil base in invalid recursive type declaration</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		setDefType(def, typ)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		typ.base = check.varType(e.X)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		return typ
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	case *ast.FuncType:
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		typ := new(Signature)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		setDefType(def, typ)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		check.funcType(typ, nil, e)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		return typ
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	case *ast.InterfaceType:
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		typ := check.newInterface()
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		setDefType(def, typ)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		check.interfaceType(typ, e, def)
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		return typ
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	case *ast.MapType:
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		typ := new(Map)
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		setDefType(def, typ)
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		typ.key = check.varType(e.Key)
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		typ.elem = check.varType(e.Value)
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;The comparison operators == and != must be fully defined</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		<span class="comment">// for operands of the key type; thus the key type must not be a</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		<span class="comment">// function, map, or slice.&#34;</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		<span class="comment">// Delay this check because it requires fully setup types;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		<span class="comment">// it is safe to continue in any case (was go.dev/issue/6667).</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		check.later(func() {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			if !Comparable(typ.key) {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>				var why string
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>				if isTypeParam(typ.key) {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>					why = &#34; (missing comparable constraint)&#34;
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>				}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>				check.errorf(e.Key, IncomparableMapKey, &#34;invalid map key type %s%s&#34;, typ.key, why)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		}).describef(e.Key, &#34;check map key %s&#34;, typ.key)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		return typ
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	case *ast.ChanType:
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		typ := new(Chan)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		setDefType(def, typ)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		dir := SendRecv
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		switch e.Dir {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		case ast.SEND | ast.RECV:
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			<span class="comment">// nothing to do</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		case ast.SEND:
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>			dir = SendOnly
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		case ast.RECV:
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			dir = RecvOnly
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		default:
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			check.errorf(e, InvalidSyntaxTree, &#34;unknown channel direction %d&#34;, e.Dir)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			<span class="comment">// ok to continue</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		typ.dir = dir
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		typ.elem = check.varType(e.Value)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		return typ
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	default:
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		check.errorf(e0, NotAType, &#34;%s is not a type&#34;, e0)
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		check.use(e0)
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	typ := Typ[Invalid]
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	setDefType(def, typ)
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	return typ
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>func setDefType(def *TypeName, typ Type) {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	if def != nil {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		switch t := def.typ.(type) {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		case *Alias:
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			<span class="comment">// t.fromRHS should always be set, either to an invalid type</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			<span class="comment">// in the beginning, or to typ in certain cyclic declarations.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			if t.fromRHS != Typ[Invalid] &amp;&amp; t.fromRHS != typ {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>				panic(sprintf(nil, nil, true, &#34;t.fromRHS = %s, typ = %s\n&#34;, t.fromRHS, typ))
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			t.fromRHS = typ
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		case *Basic:
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			assert(t == Typ[Invalid])
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		case *Named:
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			t.underlying = typ
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		default:
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			panic(fmt.Sprintf(&#34;unexpected type %T&#34;, t))
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	}
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>func (check *Checker) instantiatedType(ix *typeparams.IndexExpr, def *TypeName) (res Type) {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	if check.conf._Trace {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		check.trace(ix.Pos(), &#34;-- instantiating type %s with %s&#34;, ix.X, ix.Indices)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		check.indent++
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		defer func() {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			check.indent--
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			<span class="comment">// Don&#39;t format the underlying here. It will always be nil.</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>			check.trace(ix.Pos(), &#34;=&gt; %s&#34;, res)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		}()
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	var cause string
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	gtyp := check.genericType(ix.X, &amp;cause)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	if cause != &#34;&#34; {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		check.errorf(ix.Orig, NotAGenericType, invalidOp+&#34;%s (%s)&#34;, ix.Orig, cause)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	if !isValid(gtyp) {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		return gtyp <span class="comment">// error already reported</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	orig := asNamed(gtyp)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	if orig == nil {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		panic(fmt.Sprintf(&#34;%v: cannot instantiate %v&#34;, ix.Pos(), gtyp))
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	<span class="comment">// evaluate arguments</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	targs := check.typeList(ix.Indices)
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	if targs == nil {
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		setDefType(def, Typ[Invalid]) <span class="comment">// avoid errors later due to lazy instantiation</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		return Typ[Invalid]
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	<span class="comment">// create the instance</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	inst := asNamed(check.instance(ix.Pos(), orig, targs, nil, check.context()))
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	setDefType(def, inst)
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	<span class="comment">// orig.tparams may not be set up, so we need to do expansion later.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	check.later(func() {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		<span class="comment">// This is an instance from the source, not from recursive substitution,</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		<span class="comment">// and so it must be resolved during type-checking so that we can report</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		<span class="comment">// errors.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		check.recordInstance(ix.Orig, inst.TypeArgs().list(), inst)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		if check.validateTArgLen(ix.Pos(), inst.obj.name, inst.TypeParams().Len(), inst.TypeArgs().Len()) {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			if i, err := check.verify(ix.Pos(), inst.TypeParams().list(), inst.TypeArgs().list(), check.context()); err != nil {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>				<span class="comment">// best position for error reporting</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>				pos := ix.Pos()
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>				if i &lt; len(ix.Indices) {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>					pos = ix.Indices[i].Pos()
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>				}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>				check.softErrorf(atPos(pos), InvalidTypeArg, err.Error())
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>			} else {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>				check.mono.recordInstance(check.pkg, ix.Pos(), inst.TypeParams().list(), inst.TypeArgs().list(), ix.Indices)
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>			}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		<span class="comment">// TODO(rfindley): remove this call: we don&#39;t need to call validType here,</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		<span class="comment">// as cycles can only occur for types used inside a Named type declaration,</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		<span class="comment">// and so it suffices to call validType from declared types.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		check.validType(inst)
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	}).describef(ix, &#34;resolve instance %s&#34;, inst)
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	return inst
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span><span class="comment">// arrayLength type-checks the array length expression e</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span><span class="comment">// and returns the constant length &gt;= 0, or a value &lt; 0</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span><span class="comment">// to indicate an error (and thus an unknown length).</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>func (check *Checker) arrayLength(e ast.Expr) int64 {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	<span class="comment">// If e is an identifier, the array declaration might be an</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	<span class="comment">// attempt at a parameterized type declaration with missing</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	<span class="comment">// constraint. Provide an error message that mentions array</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	<span class="comment">// length.</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	if name, _ := e.(*ast.Ident); name != nil {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		obj := check.lookup(name.Name)
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		if obj == nil {
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>			check.errorf(name, InvalidArrayLen, &#34;undefined array length %s or missing type constraint&#34;, name.Name)
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>			return -1
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		if _, ok := obj.(*Const); !ok {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			check.errorf(name, InvalidArrayLen, &#34;invalid array length %s&#34;, name.Name)
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			return -1
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	}
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	var x operand
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	check.expr(nil, &amp;x, e)
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	if x.mode != constant_ {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		if x.mode != invalid {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			check.errorf(&amp;x, InvalidArrayLen, &#34;array length %s must be constant&#34;, &amp;x)
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		return -1
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	if isUntyped(x.typ) || isInteger(x.typ) {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		if val := constant.ToInt(x.val); val.Kind() == constant.Int {
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			if representableConst(val, check, Typ[Int], nil) {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>				if n, ok := constant.Int64Val(val); ok &amp;&amp; n &gt;= 0 {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>					return n
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>				}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>			}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		}
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	var msg string
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	if isInteger(x.typ) {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		msg = &#34;invalid array length %s&#34;
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	} else {
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		msg = &#34;array length %s must be integer&#34;
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	}
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	check.errorf(&amp;x, InvalidArrayLen, msg, &amp;x)
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	return -1
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span><span class="comment">// typeList provides the list of types corresponding to the incoming expression list.</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span><span class="comment">// If an error occurred, the result is nil, but all list elements were type-checked.</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>func (check *Checker) typeList(list []ast.Expr) []Type {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	res := make([]Type, len(list)) <span class="comment">// res != nil even if len(list) == 0</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	for i, x := range list {
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		t := check.varType(x)
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		if !isValid(t) {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			res = nil
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		if res != nil {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>			res[i] = t
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	}
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	return res
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>}
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>
</pre><p><a href="typexpr.go?m=text">View as plain text</a></p>

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
