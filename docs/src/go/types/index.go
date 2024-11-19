<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/index.go - Go Documentation Server</title>

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
<a href="index.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">index.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/types">go/types</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2021 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements typechecking of index/slice expressions.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package types
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/constant&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;go/internal/typeparams&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// If e is a valid function instantiation, indexExpr returns true.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// In that case x represents the uninstantiated function value and</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// it is the caller&#39;s responsibility to instantiate the function.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>func (check *Checker) indexExpr(x *operand, e *typeparams.IndexExpr) (isFuncInst bool) {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	check.exprOrType(x, e.X, true)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// x may be generic</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	switch x.mode {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	case invalid:
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		check.use(e.Indices...)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		return false
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	case typexpr:
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		<span class="comment">// type instantiation</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) here we re-evaluate e.X - try to avoid this</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		x.typ = check.varType(e.Orig)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		if isValid(x.typ) {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>			x.mode = typexpr
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		return false
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	case value:
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		if sig, _ := under(x.typ).(*Signature); sig != nil &amp;&amp; sig.TypeParams().Len() &gt; 0 {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>			<span class="comment">// function instantiation</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>			return true
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		}
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// x should not be generic at this point, but be safe and check</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	check.nonGeneric(nil, x)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	if x.mode == invalid {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		return false
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// ordinary index expression</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	valid := false
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	length := int64(-1) <span class="comment">// valid if &gt;= 0</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	switch typ := under(x.typ).(type) {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	case *Basic:
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		if isString(typ) {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			valid = true
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			if x.mode == constant_ {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>				length = int64(len(constant.StringVal(x.val)))
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			<span class="comment">// an indexed string always yields a byte value</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			<span class="comment">// (not a constant) even if the string and the</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			<span class="comment">// index are constant</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			x.mode = value
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			x.typ = universeByte <span class="comment">// use &#39;byte&#39; name</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	case *Array:
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		valid = true
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		length = typ.len
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		if x.mode != variable {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			x.mode = value
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		x.typ = typ.elem
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	case *Pointer:
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		if typ, _ := under(typ.base).(*Array); typ != nil {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			valid = true
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			length = typ.len
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			x.mode = variable
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			x.typ = typ.elem
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	case *Slice:
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		valid = true
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		x.mode = variable
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		x.typ = typ.elem
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	case *Map:
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		index := check.singleIndex(e)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		if index == nil {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			return false
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		var key operand
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		check.expr(nil, &amp;key, index)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		check.assignment(&amp;key, typ.key, &#34;map index&#34;)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		<span class="comment">// ok to continue even if indexing failed - map element type is known</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		x.mode = mapindex
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		x.typ = typ.elem
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		x.expr = e.Orig
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		return false
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	case *Interface:
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		if !isTypeParam(x.typ) {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			break
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) report detailed failure cause for better error messages</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		var key, elem Type <span class="comment">// key != nil: we must have all maps</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		mode := variable   <span class="comment">// non-maps result mode</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) factor out closure and use it for non-typeparam cases as well</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		if typ.typeSet().underIs(func(u Type) bool {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			l := int64(-1) <span class="comment">// valid if &gt;= 0</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			var k, e Type  <span class="comment">// k is only set for maps</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			switch t := u.(type) {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			case *Basic:
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>				if isString(t) {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>					e = universeByte
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>					mode = value
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>				}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			case *Array:
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>				l = t.len
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>				e = t.elem
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>				if x.mode != variable {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>					mode = value
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>				}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			case *Pointer:
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>				if t, _ := under(t.base).(*Array); t != nil {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>					l = t.len
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>					e = t.elem
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>				}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			case *Slice:
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>				e = t.elem
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			case *Map:
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>				k = t.key
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>				e = t.elem
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			if e == nil {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>				return false
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			if elem == nil {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>				<span class="comment">// first type</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>				length = l
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>				key, elem = k, e
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>				return true
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			<span class="comment">// all map keys must be identical (incl. all nil)</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			<span class="comment">// (that is, we cannot mix maps with other types)</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			if !Identical(key, k) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>				return false
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			<span class="comment">// all element types must be identical</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			if !Identical(elem, e) {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>				return false
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			<span class="comment">// track the minimal length for arrays, if any</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			if l &gt;= 0 &amp;&amp; l &lt; length {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>				length = l
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			return true
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		}) {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			<span class="comment">// For maps, the index expression must be assignable to the map key type.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			if key != nil {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>				index := check.singleIndex(e)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>				if index == nil {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>					x.mode = invalid
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>					return false
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>				}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>				var k operand
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>				check.expr(nil, &amp;k, index)
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>				check.assignment(&amp;k, key, &#34;map index&#34;)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>				<span class="comment">// ok to continue even if indexing failed - map element type is known</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>				x.mode = mapindex
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>				x.typ = elem
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>				x.expr = e
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>				return false
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			<span class="comment">// no maps</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			valid = true
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			x.mode = mode
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			x.typ = elem
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	if !valid {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		<span class="comment">// types2 uses the position of &#39;[&#39; for the error</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		check.errorf(x, NonIndexableOperand, invalidOp+&#34;cannot index %s&#34;, x)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		check.use(e.Indices...)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		return false
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	index := check.singleIndex(e)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	if index == nil {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		return false
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// In pathological (invalid) cases (e.g.: type T1 [][[]T1{}[0][0]]T0)</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// the element type may be accessed before it&#39;s set. Make sure we have</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// a valid type.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	if x.typ == nil {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		x.typ = Typ[Invalid]
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	check.index(index, length)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	return false
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>func (check *Checker) sliceExpr(x *operand, e *ast.SliceExpr) {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	check.expr(nil, x, e.X)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	if x.mode == invalid {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		check.use(e.Low, e.High, e.Max)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		return
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	valid := false
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	length := int64(-1) <span class="comment">// valid if &gt;= 0</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	switch u := coreString(x.typ).(type) {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	case nil:
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		check.errorf(x, NonSliceableOperand, invalidOp+&#34;cannot slice %s: %s has no core type&#34;, x, x.typ)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		return
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	case *Basic:
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		if isString(u) {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			if e.Slice3 {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>				at := e.Max
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>				if at == nil {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>					at = e <span class="comment">// e.Index[2] should be present but be careful</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>				}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>				check.error(at, InvalidSliceExpr, invalidOp+&#34;3-index slice of string&#34;)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>				x.mode = invalid
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>				return
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			valid = true
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			if x.mode == constant_ {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>				length = int64(len(constant.StringVal(x.val)))
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			<span class="comment">// spec: &#34;For untyped string operands the result</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			<span class="comment">// is a non-constant value of type string.&#34;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			if isUntyped(x.typ) {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>				x.typ = Typ[String]
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	case *Array:
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		valid = true
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		length = u.len
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		if x.mode != variable {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			check.errorf(x, NonSliceableOperand, invalidOp+&#34;cannot slice %s (value not addressable)&#34;, x)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			return
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		x.typ = &amp;Slice{elem: u.elem}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	case *Pointer:
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		if u, _ := under(u.base).(*Array); u != nil {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			valid = true
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			length = u.len
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			x.typ = &amp;Slice{elem: u.elem}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	case *Slice:
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		valid = true
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		<span class="comment">// x.typ doesn&#39;t change</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	}
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	if !valid {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		check.errorf(x, NonSliceableOperand, invalidOp+&#34;cannot slice %s&#34;, x)
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		return
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	x.mode = value
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;Only the first index may be omitted; it defaults to 0.&#34;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	if e.Slice3 &amp;&amp; (e.High == nil || e.Max == nil) {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		check.error(inNode(e, e.Rbrack), InvalidSyntaxTree, &#34;2nd and 3rd index required in 3-index slice&#34;)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		return
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">// check indices</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	var ind [3]int64
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	for i, expr := range []ast.Expr{e.Low, e.High, e.Max} {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		x := int64(-1)
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		switch {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		case expr != nil:
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			<span class="comment">// The &#34;capacity&#34; is only known statically for strings, arrays,</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			<span class="comment">// and pointers to arrays, and it is the same as the length for</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			<span class="comment">// those types.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			max := int64(-1)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			if length &gt;= 0 {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>				max = length + 1
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			if _, v := check.index(expr, max); v &gt;= 0 {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				x = v
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		case i == 0:
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			<span class="comment">// default is 0 for the first index</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			x = 0
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		case length &gt;= 0:
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			<span class="comment">// default is length (== capacity) otherwise</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			x = length
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		ind[i] = x
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	<span class="comment">// constant indices must be in range</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	<span class="comment">// (check.index already checks that existing indices &gt;= 0)</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>L:
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	for i, x := range ind[:len(ind)-1] {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		if x &gt; 0 {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>			for j, y := range ind[i+1:] {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>				if y &gt;= 0 &amp;&amp; y &lt; x {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>					<span class="comment">// The value y corresponds to the expression e.Index[i+1+j].</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>					<span class="comment">// Because y &gt;= 0, it must have been set from the expression</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>					<span class="comment">// when checking indices and thus e.Index[i+1+j] is not nil.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>					at := []ast.Expr{e.Low, e.High, e.Max}[i+1+j]
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>					check.errorf(at, SwappedSliceIndices, &#34;invalid slice indices: %d &lt; %d&#34;, y, x)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>					break L <span class="comment">// only report one error, ok to continue</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>				}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">// singleIndex returns the (single) index from the index expression e.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span><span class="comment">// If the index is missing, or if there are multiple indices, an error</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// is reported and the result is nil.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>func (check *Checker) singleIndex(expr *typeparams.IndexExpr) ast.Expr {
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	if len(expr.Indices) == 0 {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		check.errorf(expr.Orig, InvalidSyntaxTree, &#34;index expression %v with 0 indices&#34;, expr)
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		return nil
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	if len(expr.Indices) &gt; 1 {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		<span class="comment">// TODO(rFindley) should this get a distinct error code?</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		check.error(expr.Indices[1], InvalidIndex, invalidOp+&#34;more than one index&#34;)
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	return expr.Indices[0]
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// index checks an index expression for validity.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// If max &gt;= 0, it is the upper bound for index.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// If the result typ is != Typ[Invalid], index is valid and typ is its (possibly named) integer type.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// If the result val &gt;= 0, index is valid and val is its constant int value.</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>func (check *Checker) index(index ast.Expr, max int64) (typ Type, val int64) {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	typ = Typ[Invalid]
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	val = -1
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	var x operand
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	check.expr(nil, &amp;x, index)
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	if !check.isValidIndex(&amp;x, InvalidIndex, &#34;index&#34;, false) {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		return
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	if x.mode != constant_ {
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		return x.typ, -1
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	if x.val.Kind() == constant.Unknown {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		return
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	v, ok := constant.Int64Val(x.val)
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	assert(ok)
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	if max &gt;= 0 &amp;&amp; v &gt;= max {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		check.errorf(&amp;x, InvalidIndex, invalidArg+&#34;index %s out of bounds [0:%d]&#34;, x.val.String(), max)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		return
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	<span class="comment">// 0 &lt;= v [ &amp;&amp; v &lt; max ]</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	return x.typ, v
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>func (check *Checker) isValidIndex(x *operand, code Code, what string, allowNegative bool) bool {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	if x.mode == invalid {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		return false
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;a constant index that is untyped is given type int&#34;</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	check.convertUntyped(x, Typ[Int])
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	if x.mode == invalid {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		return false
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;the index x must be of integer type or an untyped constant&#34;</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	if !allInteger(x.typ) {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		check.errorf(x, code, invalidArg+&#34;%s %s must be integer&#34;, what, x)
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		return false
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	if x.mode == constant_ {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;a constant index must be non-negative ...&#34;</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		if !allowNegative &amp;&amp; constant.Sign(x.val) &lt; 0 {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			check.errorf(x, code, invalidArg+&#34;%s %s must not be negative&#34;, what, x)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			return false
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;... and representable by a value of type int&#34;</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		if !representableConst(x.val, check, Typ[Int], &amp;x.val) {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			check.errorf(x, code, invalidArg+&#34;%s %s overflows int&#34;, what, x)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			return false
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	return true
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span><span class="comment">// indexedElts checks the elements (elts) of an array or slice composite literal</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span><span class="comment">// against the literal&#39;s element type (typ), and the element indices against</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span><span class="comment">// the literal length if known (length &gt;= 0). It returns the length of the</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span><span class="comment">// literal (maximum index value + 1).</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>func (check *Checker) indexedElts(elts []ast.Expr, typ Type, length int64) int64 {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	visited := make(map[int64]bool, len(elts))
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	var index, max int64
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	for _, e := range elts {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		<span class="comment">// determine and check index</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		validIndex := false
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		eval := e
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		if kv, _ := e.(*ast.KeyValueExpr); kv != nil {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			if typ, i := check.index(kv.Key, length); isValid(typ) {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>				if i &gt;= 0 {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>					index = i
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>					validIndex = true
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>				} else {
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>					check.errorf(e, InvalidLitIndex, &#34;index %s must be integer constant&#34;, kv.Key)
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>				}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>			}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			eval = kv.Value
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		} else if length &gt;= 0 &amp;&amp; index &gt;= length {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			check.errorf(e, OversizeArrayLit, &#34;index %d is out of bounds (&gt;= %d)&#34;, index, length)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		} else {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>			validIndex = true
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		<span class="comment">// if we have a valid index, check for duplicate entries</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		if validIndex {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>			if visited[index] {
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>				check.errorf(e, DuplicateLitKey, &#34;duplicate index %d in array or slice literal&#34;, index)
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>			}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			visited[index] = true
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		index++
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		if index &gt; max {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>			max = index
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		<span class="comment">// check element against composite literal element type</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		var x operand
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		check.exprWithHint(&amp;x, eval, typ)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		check.assignment(&amp;x, typ, &#34;array or slice literal&#34;)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	return max
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>
</pre><p><a href="index.go?m=text">View as plain text</a></p>

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
