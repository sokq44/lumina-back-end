<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/conversions.go - Go Documentation Server</title>

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
<a href="conversions.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">conversions.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements typechecking of conversions.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package types
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;go/constant&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// conversion type-checks the conversion T(x).</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// The result is in x.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>func (check *Checker) conversion(x *operand, T Type) {
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	constArg := x.mode == constant_
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	constConvertibleTo := func(T Type, val *constant.Value) bool {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>		switch t, _ := under(T).(*Basic); {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>		case t == nil:
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>			<span class="comment">// nothing to do</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>		case representableConst(x.val, check, t, val):
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>			return true
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		case isInteger(x.typ) &amp;&amp; isString(t):
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>			codepoint := unicode.ReplacementChar
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>			if i, ok := constant.Uint64Val(x.val); ok &amp;&amp; i &lt;= unicode.MaxRune {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>				codepoint = rune(i)
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>			}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>			if val != nil {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>				*val = constant.MakeString(string(codepoint))
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>			}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>			return true
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		return false
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	var ok bool
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	var cause string
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	switch {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	case constArg &amp;&amp; isConstType(T):
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		<span class="comment">// constant conversion</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		ok = constConvertibleTo(T, &amp;x.val)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		<span class="comment">// A conversion from an integer constant to an integer type</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		<span class="comment">// can only fail if there&#39;s overflow. Give a concise error.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		<span class="comment">// (go.dev/issue/63563)</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		if !ok &amp;&amp; isInteger(x.typ) &amp;&amp; isInteger(T) {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			check.errorf(x, InvalidConversion, &#34;constant %s overflows %s&#34;, x.val, T)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			x.mode = invalid
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			return
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	case constArg &amp;&amp; isTypeParam(T):
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		<span class="comment">// x is convertible to T if it is convertible</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		<span class="comment">// to each specific type in the type set of T.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		<span class="comment">// If T&#39;s type set is empty, or if it doesn&#39;t</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		<span class="comment">// have specific types, constant x cannot be</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		<span class="comment">// converted.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		ok = T.(*TypeParam).underIs(func(u Type) bool {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			<span class="comment">// u is nil if there are no specific type terms</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			if u == nil {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>				cause = check.sprintf(&#34;%s does not contain specific types&#34;, T)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>				return false
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			if isString(x.typ) &amp;&amp; isBytesOrRunes(u) {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>				return true
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			if !constConvertibleTo(u, nil) {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>				if isInteger(x.typ) &amp;&amp; isInteger(u) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>					<span class="comment">// see comment above on constant conversion</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>					cause = check.sprintf(&#34;constant %s overflows %s (in %s)&#34;, x.val, u, T)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>				} else {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>					cause = check.sprintf(&#34;cannot convert %s to type %s (in %s)&#34;, x, u, T)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>				}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>				return false
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			return true
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		})
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		x.mode = value <span class="comment">// type parameters are not constants</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	case x.convertibleTo(check, T, &amp;cause):
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		<span class="comment">// non-constant conversion</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		ok = true
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		x.mode = value
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	if !ok {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		if cause != &#34;&#34; {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			check.errorf(x, InvalidConversion, &#34;cannot convert %s to type %s: %s&#34;, x, T, cause)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		} else {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			check.errorf(x, InvalidConversion, &#34;cannot convert %s to type %s&#34;, x, T)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		return
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// The conversion argument types are final. For untyped values the</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// conversion provides the type, per the spec: &#34;A constant may be</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// given a type explicitly by a constant declaration or conversion,...&#34;.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	if isUntyped(x.typ) {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		final := T
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		<span class="comment">// - For conversions to interfaces, use the argument&#39;s default type.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		<span class="comment">// - For conversions of untyped constants to non-constant types, also</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		<span class="comment">//   use the default type (e.g., []byte(&#34;foo&#34;) should report string</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		<span class="comment">//   not []byte as type for the constant &#34;foo&#34;).</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		<span class="comment">// - Keep untyped nil for untyped nil arguments.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">// - For constant integer to string conversions, keep the argument type.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		<span class="comment">//   (See also the TODO below.)</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		if isNonTypeParamInterface(T) || constArg &amp;&amp; !isConstType(T) || x.isNil() {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			final = Default(x.typ) <span class="comment">// default type of untyped nil is untyped nil</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		} else if x.mode == constant_ &amp;&amp; isInteger(x.typ) &amp;&amp; allString(T) {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			final = x.typ
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		check.updateExprType(x.expr, final, true)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	x.typ = T
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// TODO(gri) convertibleTo checks if T(x) is valid. It assumes that the type</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// of x is fully known, but that&#39;s not the case for say string(1&lt;&lt;s + 1.0):</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// Here, the type of 1&lt;&lt;s + 1.0 will be UntypedFloat which will lead to the</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// (correct!) refusal of the conversion. But the reported error is essentially</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// &#34;cannot convert untyped float value to string&#34;, yet the correct error (per</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// the spec) is that we cannot shift a floating-point value: 1 in 1&lt;&lt;s should</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// be converted to UntypedFloat because of the addition of 1.0. Fixing this</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// is tricky because we&#39;d have to run updateExprType on the argument first.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// (go.dev/issue/21982.)</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// convertibleTo reports whether T(x) is valid. In the failure case, *cause</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// may be set to the cause for the failure.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">// The check parameter may be nil if convertibleTo is invoked through an</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// exported API call, i.e., when all methods have been type-checked.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>func (x *operand) convertibleTo(check *Checker, T Type, cause *string) bool {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// &#34;x is assignable to T&#34;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	if ok, _ := x.assignableTo(check, T, cause); ok {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		return true
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// &#34;V and T have identical underlying types if tags are ignored</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// and V and T are not type parameters&#34;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	V := x.typ
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	Vu := under(V)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	Tu := under(T)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	Vp, _ := V.(*TypeParam)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	Tp, _ := T.(*TypeParam)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	if IdenticalIgnoreTags(Vu, Tu) &amp;&amp; Vp == nil &amp;&amp; Tp == nil {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		return true
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// &#34;V and T are unnamed pointer types and their pointer base types</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// have identical underlying types if tags are ignored</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// and their pointer base types are not type parameters&#34;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if V, ok := V.(*Pointer); ok {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		if T, ok := T.(*Pointer); ok {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			if IdenticalIgnoreTags(under(V.base), under(T.base)) &amp;&amp; !isTypeParam(V.base) &amp;&amp; !isTypeParam(T.base) {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>				return true
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// &#34;V and T are both integer or floating point types&#34;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	if isIntegerOrFloat(Vu) &amp;&amp; isIntegerOrFloat(Tu) {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		return true
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// &#34;V and T are both complex types&#34;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	if isComplex(Vu) &amp;&amp; isComplex(Tu) {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		return true
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// &#34;V is an integer or a slice of bytes or runes and T is a string type&#34;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	if (isInteger(Vu) || isBytesOrRunes(Vu)) &amp;&amp; isString(Tu) {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		return true
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// &#34;V is a string and T is a slice of bytes or runes&#34;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	if isString(Vu) &amp;&amp; isBytesOrRunes(Tu) {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		return true
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// package unsafe:</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// &#34;any pointer or value of underlying type uintptr can be converted into a unsafe.Pointer&#34;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	if (isPointer(Vu) || isUintptr(Vu)) &amp;&amp; isUnsafePointer(Tu) {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		return true
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// &#34;and vice versa&#34;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	if isUnsafePointer(Vu) &amp;&amp; (isPointer(Tu) || isUintptr(Tu)) {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		return true
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// &#34;V is a slice, T is an array or pointer-to-array type,</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// and the slice and array types have identical element types.&#34;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	if s, _ := Vu.(*Slice); s != nil {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		switch a := Tu.(type) {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		case *Array:
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			if Identical(s.Elem(), a.Elem()) {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>				if check == nil || check.allowVersion(check.pkg, x, go1_20) {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>					return true
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>				}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>				<span class="comment">// check != nil</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>				if cause != nil {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>					<span class="comment">// TODO(gri) consider restructuring versionErrorf so we can use it here and below</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>					*cause = &#34;conversion of slices to arrays requires go1.20 or later&#34;
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>				}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>				return false
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		case *Pointer:
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			if a, _ := under(a.Elem()).(*Array); a != nil {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>				if Identical(s.Elem(), a.Elem()) {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>					if check == nil || check.allowVersion(check.pkg, x, go1_17) {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>						return true
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>					}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>					<span class="comment">// check != nil</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>					if cause != nil {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>						*cause = &#34;conversion of slices to array pointers requires go1.17 or later&#34;
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>					}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>					return false
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>				}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	<span class="comment">// optimization: if we don&#39;t have type parameters, we&#39;re done</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	if Vp == nil &amp;&amp; Tp == nil {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		return false
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	errorf := func(format string, args ...any) {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		if check != nil &amp;&amp; cause != nil {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			msg := check.sprintf(format, args...)
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			if *cause != &#34;&#34; {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>				msg += &#34;\n\t&#34; + *cause
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			*cause = msg
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// generic cases with specific type terms</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	<span class="comment">// (generic operands cannot be constants, so we can ignore x.val)</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	switch {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	case Vp != nil &amp;&amp; Tp != nil:
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		x := *x <span class="comment">// don&#39;t clobber outer x</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		return Vp.is(func(V *term) bool {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			if V == nil {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>				return false <span class="comment">// no specific types</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			x.typ = V.typ
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			return Tp.is(func(T *term) bool {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>				if T == nil {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>					return false <span class="comment">// no specific types</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>				}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>				if !x.convertibleTo(check, T.typ, cause) {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>					errorf(&#34;cannot convert %s (in %s) to type %s (in %s)&#34;, V.typ, Vp, T.typ, Tp)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>					return false
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>				return true
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			})
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		})
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	case Vp != nil:
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		x := *x <span class="comment">// don&#39;t clobber outer x</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		return Vp.is(func(V *term) bool {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			if V == nil {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>				return false <span class="comment">// no specific types</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			x.typ = V.typ
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			if !x.convertibleTo(check, T, cause) {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>				errorf(&#34;cannot convert %s (in %s) to type %s&#34;, V.typ, Vp, T)
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				return false
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			return true
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		})
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	case Tp != nil:
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		return Tp.is(func(T *term) bool {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			if T == nil {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>				return false <span class="comment">// no specific types</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			if !x.convertibleTo(check, T.typ, cause) {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>				errorf(&#34;cannot convert %s to type %s (in %s)&#34;, x.typ, T.typ, Tp)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>				return false
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			return true
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		})
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	return false
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>func isUintptr(typ Type) bool {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	t, _ := under(typ).(*Basic)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	return t != nil &amp;&amp; t.kind == Uintptr
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>func isUnsafePointer(typ Type) bool {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	t, _ := under(typ).(*Basic)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	return t != nil &amp;&amp; t.kind == UnsafePointer
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>func isPointer(typ Type) bool {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	_, ok := under(typ).(*Pointer)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	return ok
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>func isBytesOrRunes(typ Type) bool {
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	if s, _ := under(typ).(*Slice); s != nil {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		t, _ := under(s.elem).(*Basic)
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		return t != nil &amp;&amp; (t.kind == Byte || t.kind == Rune)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	return false
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
</pre><p><a href="conversions.go?m=text">View as plain text</a></p>

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
