<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/builtins.go - Go Documentation Server</title>

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
<a href="builtins.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">builtins.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements typechecking of builtin function calls.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package types
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/constant&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// builtin type-checks a call to the built-in specified by id and</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// reports whether the call is valid, with *x holding the result;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// but x.expr is not set. If the call is invalid, the result is</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// false, and *x is undefined.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>func (check *Checker) builtin(x *operand, call *ast.CallExpr, id builtinId) (_ bool) {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	argList := call.Args
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// append is the only built-in that permits the use of ... for the last argument</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	bin := predeclaredFuncs[id]
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	if call.Ellipsis.IsValid() &amp;&amp; id != _Append {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		check.errorf(atPos(call.Ellipsis),
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>			InvalidDotDotDot,
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>			invalidOp+&#34;invalid use of ... with built-in %s&#34;, bin.name)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		check.use(argList...)
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		return
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// For len(x) and cap(x) we need to know if x contains any function calls or</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// receive operations. Save/restore current setting and set hasCallOrRecv to</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// false for the evaluation of x so that we can check it afterwards.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// Note: We must do this _before_ calling exprList because exprList evaluates</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">//       all arguments.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	if id == _Len || id == _Cap {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		defer func(b bool) {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>			check.hasCallOrRecv = b
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		}(check.hasCallOrRecv)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		check.hasCallOrRecv = false
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// Evaluate arguments for built-ins that use ordinary (value) arguments.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// For built-ins with special argument handling (make, new, etc.),</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// evaluation is done by the respective built-in code.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	var args []*operand <span class="comment">// not valid for _Make, _New, _Offsetof, _Trace</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	var nargs int
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	switch id {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	default:
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		<span class="comment">// check all arguments</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		args = check.exprList(argList)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		nargs = len(args)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		for _, a := range args {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			if a.mode == invalid {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>				return
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		<span class="comment">// first argument is always in x</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		if nargs &gt; 0 {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			*x = *args[0]
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	case _Make, _New, _Offsetof, _Trace:
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		<span class="comment">// arguments require special handling</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		nargs = len(argList)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// check argument count</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	{
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		msg := &#34;&#34;
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		if nargs &lt; bin.nargs {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			msg = &#34;not enough&#34;
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		} else if !bin.variadic &amp;&amp; nargs &gt; bin.nargs {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			msg = &#34;too many&#34;
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		if msg != &#34;&#34; {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			check.errorf(inNode(call, call.Rparen), WrongArgCount, invalidOp+&#34;%s arguments for %s (expected %d, found %d)&#34;, msg, call, bin.nargs, nargs)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			return
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	switch id {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	case _Append:
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		<span class="comment">// append(s S, x ...T) S, where T is the element type of S</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;The variadic function append appends zero or more values x to s of type</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		<span class="comment">// S, which must be a slice type, and returns the resulting slice, also of type S.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		<span class="comment">// The values x are passed to a parameter of type ...T where T is the element type</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		<span class="comment">// of S and the respective parameter passing rules apply.&#34;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		S := x.typ
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		var T Type
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		if s, _ := coreType(S).(*Slice); s != nil {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			T = s.elem
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		} else {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			var cause string
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			switch {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			case x.isNil():
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>				cause = &#34;have untyped nil&#34;
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			case isTypeParam(S):
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>				if u := coreType(S); u != nil {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>					cause = check.sprintf(&#34;%s has core type %s&#34;, x, u)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>				} else {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>					cause = check.sprintf(&#34;%s has no core type&#34;, x)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>				}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			default:
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>				cause = check.sprintf(&#34;have %s&#34;, x)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			<span class="comment">// don&#39;t use invalidArg prefix here as it would repeat &#34;argument&#34; in the error message</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			check.errorf(x, InvalidAppend, &#34;first argument to append must be a slice; %s&#34;, cause)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			return
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;As a special case, append also accepts a first argument assignable</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		<span class="comment">// to type []byte with a second argument of string type followed by ... .</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		<span class="comment">// This form appends the bytes of the string.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		if nargs == 2 &amp;&amp; call.Ellipsis.IsValid() {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			if ok, _ := x.assignableTo(check, NewSlice(universeByte), nil); ok {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>				y := args[1]
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>				if t := coreString(y.typ); t != nil &amp;&amp; isString(t) {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>					if check.recordTypes() {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>						sig := makeSig(S, S, y.typ)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>						sig.variadic = true
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>						check.recordBuiltinType(call.Fun, sig)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>					}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>					x.mode = value
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>					x.typ = S
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>					break
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>				}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		<span class="comment">// check general case by creating custom signature</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		sig := makeSig(S, S, NewSlice(T)) <span class="comment">// []T required for variadic signature</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		sig.variadic = true
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		check.arguments(call, sig, nil, nil, args, nil, nil) <span class="comment">// discard result (we know the result type)</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// ok to continue even if check.arguments reported errors</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		x.mode = value
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		x.typ = S
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, sig)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	case _Cap, _Len:
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		<span class="comment">// cap(x)</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		<span class="comment">// len(x)</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		mode := invalid
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		var val constant.Value
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		switch t := arrayPtrDeref(under(x.typ)).(type) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		case *Basic:
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			if isString(t) &amp;&amp; id == _Len {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>				if x.mode == constant_ {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>					mode = constant_
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>					val = constant.MakeInt64(int64(len(constant.StringVal(x.val))))
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>				} else {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>					mode = value
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>				}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		case *Array:
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			mode = value
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			<span class="comment">// spec: &#34;The expressions len(s) and cap(s) are constants</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			<span class="comment">// if the type of s is an array or pointer to an array and</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			<span class="comment">// the expression s does not contain channel receives or</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			<span class="comment">// function calls; in this case s is not evaluated.&#34;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			if !check.hasCallOrRecv {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>				mode = constant_
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>				if t.len &gt;= 0 {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>					val = constant.MakeInt64(t.len)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>				} else {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>					val = constant.MakeUnknown()
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>				}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		case *Slice, *Chan:
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			mode = value
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		case *Map:
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			if id == _Len {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>				mode = value
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		case *Interface:
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			if !isTypeParam(x.typ) {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>				break
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			if t.typeSet().underIs(func(t Type) bool {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				switch t := arrayPtrDeref(t).(type) {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				case *Basic:
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>					if isString(t) &amp;&amp; id == _Len {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>						return true
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>					}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>				case *Array, *Slice, *Chan:
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>					return true
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>				case *Map:
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>					if id == _Len {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>						return true
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>					}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>				}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>				return false
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			}) {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>				mode = value
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		if mode == invalid {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			<span class="comment">// avoid error if underlying type is invalid</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			if isValid(under(x.typ)) {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>				code := InvalidCap
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>				if id == _Len {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>					code = InvalidLen
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>				}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>				check.errorf(x, code, invalidArg+&#34;%s for %s&#34;, x, bin.name)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			return
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// record the signature before changing x.typ</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		if check.recordTypes() &amp;&amp; mode != constant_ {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(Typ[Int], x.typ))
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		x.mode = mode
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		x.typ = Typ[Int]
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		x.val = val
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	case _Clear:
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		<span class="comment">// clear(m)</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		check.verifyVersionf(call.Fun, go1_21, &#34;clear&#34;)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		if !underIs(x.typ, func(u Type) bool {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			switch u.(type) {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			case *Map, *Slice:
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>				return true
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			check.errorf(x, InvalidClear, invalidArg+&#34;cannot clear %s: argument must be (or constrained by) map or slice&#34;, x)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			return false
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		}) {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			return
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		x.mode = novalue
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(nil, x.typ))
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	case _Close:
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		<span class="comment">// close(c)</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		if !underIs(x.typ, func(u Type) bool {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			uch, _ := u.(*Chan)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			if uch == nil {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>				check.errorf(x, InvalidClose, invalidOp+&#34;cannot close non-channel %s&#34;, x)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>				return false
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			if uch.dir == RecvOnly {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>				check.errorf(x, InvalidClose, invalidOp+&#34;cannot close receive-only channel %s&#34;, x)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>				return false
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			return true
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		}) {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			return
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		x.mode = novalue
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(nil, x.typ))
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	case _Complex:
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		<span class="comment">// complex(x, y floatT) complexT</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		y := args[1]
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		<span class="comment">// convert or check untyped arguments</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		d := 0
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		if isUntyped(x.typ) {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			d |= 1
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		if isUntyped(y.typ) {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			d |= 2
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		switch d {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		case 0:
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			<span class="comment">// x and y are typed =&gt; nothing to do</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		case 1:
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			<span class="comment">// only x is untyped =&gt; convert to type of y</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			check.convertUntyped(x, y.typ)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		case 2:
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			<span class="comment">// only y is untyped =&gt; convert to type of x</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			check.convertUntyped(y, x.typ)
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		case 3:
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			<span class="comment">// x and y are untyped =&gt;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			<span class="comment">// 1) if both are constants, convert them to untyped</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			<span class="comment">//    floating-point numbers if possible,</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			<span class="comment">// 2) if one of them is not constant (possible because</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			<span class="comment">//    it contains a shift that is yet untyped), convert</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			<span class="comment">//    both of them to float64 since they must have the</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			<span class="comment">//    same type to succeed (this will result in an error</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			<span class="comment">//    because shifts of floats are not permitted)</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			if x.mode == constant_ &amp;&amp; y.mode == constant_ {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				toFloat := func(x *operand) {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>					if isNumeric(x.typ) &amp;&amp; constant.Sign(constant.Imag(x.val)) == 0 {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>						x.typ = Typ[UntypedFloat]
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>					}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>				}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>				toFloat(x)
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>				toFloat(y)
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			} else {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>				check.convertUntyped(x, Typ[Float64])
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>				check.convertUntyped(y, Typ[Float64])
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>				<span class="comment">// x and y should be invalid now, but be conservative</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>				<span class="comment">// and check below</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>			}
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		if x.mode == invalid || y.mode == invalid {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>			return
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		<span class="comment">// both argument types must be identical</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		if !Identical(x.typ, y.typ) {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			check.errorf(x, InvalidComplex, invalidOp+&#34;%v (mismatched types %s and %s)&#34;, call, x.typ, y.typ)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			return
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		<span class="comment">// the argument types must be of floating-point type</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		<span class="comment">// (applyTypeFunc never calls f with a type parameter)</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		f := func(typ Type) Type {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			assert(!isTypeParam(typ))
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			if t, _ := under(typ).(*Basic); t != nil {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>				switch t.kind {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>				case Float32:
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>					return Typ[Complex64]
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>				case Float64:
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>					return Typ[Complex128]
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>				case UntypedFloat:
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>					return Typ[UntypedComplex]
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>				}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			return nil
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		resTyp := check.applyTypeFunc(f, x, id)
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		if resTyp == nil {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			check.errorf(x, InvalidComplex, invalidArg+&#34;arguments have type %s, expected floating-point&#34;, x.typ)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			return
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		<span class="comment">// if both arguments are constants, the result is a constant</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		if x.mode == constant_ &amp;&amp; y.mode == constant_ {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			x.val = constant.BinaryOp(constant.ToFloat(x.val), token.ADD, constant.MakeImag(constant.ToFloat(y.val)))
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		} else {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>			x.mode = value
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		if check.recordTypes() &amp;&amp; x.mode != constant_ {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(resTyp, x.typ, x.typ))
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		x.typ = resTyp
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	case _Copy:
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		<span class="comment">// copy(x, y []T) int</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		dst, _ := coreType(x.typ).(*Slice)
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		y := args[1]
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		src0 := coreString(y.typ)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		if src0 != nil &amp;&amp; isString(src0) {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			src0 = NewSlice(universeByte)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		src, _ := src0.(*Slice)
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		if dst == nil || src == nil {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			check.errorf(x, InvalidCopy, invalidArg+&#34;copy expects slice arguments; found %s and %s&#34;, x, y)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>			return
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		if !Identical(dst.elem, src.elem) {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			check.errorf(x, InvalidCopy, invalidArg+&#34;arguments to copy %s and %s have different element types %s and %s&#34;, x, y, dst.elem, src.elem)
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			return
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(Typ[Int], x.typ, y.typ))
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		x.mode = value
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		x.typ = Typ[Int]
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	case _Delete:
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		<span class="comment">// delete(map_, key)</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		<span class="comment">// map_ must be a map type or a type parameter describing map types.</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		<span class="comment">// The key cannot be a type parameter for now.</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		map_ := x.typ
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		var key Type
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		if !underIs(map_, func(u Type) bool {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			map_, _ := u.(*Map)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			if map_ == nil {
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>				check.errorf(x, InvalidDelete, invalidArg+&#34;%s is not a map&#34;, x)
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>				return false
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			if key != nil &amp;&amp; !Identical(map_.key, key) {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>				check.errorf(x, InvalidDelete, invalidArg+&#34;maps of %s must have identical key types&#34;, x)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>				return false
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			key = map_.key
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			return true
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		}) {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			return
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		*x = *args[1] <span class="comment">// key</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		check.assignment(x, key, &#34;argument to delete&#34;)
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			return
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		}
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		x.mode = novalue
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(nil, map_, key))
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	case _Imag, _Real:
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		<span class="comment">// imag(complexT) floatT</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		<span class="comment">// real(complexT) floatT</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		<span class="comment">// convert or check untyped argument</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		if isUntyped(x.typ) {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			if x.mode == constant_ {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>				<span class="comment">// an untyped constant number can always be considered</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>				<span class="comment">// as a complex constant</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>				if isNumeric(x.typ) {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>					x.typ = Typ[UntypedComplex]
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>				}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			} else {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>				<span class="comment">// an untyped non-constant argument may appear if</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>				<span class="comment">// it contains a (yet untyped non-constant) shift</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>				<span class="comment">// expression: convert it to complex128 which will</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>				<span class="comment">// result in an error (shift of complex value)</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>				check.convertUntyped(x, Typ[Complex128])
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>				<span class="comment">// x should be invalid now, but be conservative and check</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>				if x.mode == invalid {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>					return
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>				}
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		<span class="comment">// the argument must be of complex type</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		<span class="comment">// (applyTypeFunc never calls f with a type parameter)</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		f := func(typ Type) Type {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>			assert(!isTypeParam(typ))
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>			if t, _ := under(typ).(*Basic); t != nil {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>				switch t.kind {
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>				case Complex64:
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>					return Typ[Float32]
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>				case Complex128:
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>					return Typ[Float64]
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>				case UntypedComplex:
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>					return Typ[UntypedFloat]
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>				}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>			}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			return nil
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		resTyp := check.applyTypeFunc(f, x, id)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		if resTyp == nil {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			code := InvalidImag
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>			if id == _Real {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>				code = InvalidReal
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			}
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>			check.errorf(x, code, invalidArg+&#34;argument has type %s, expected complex type&#34;, x.typ)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			return
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		<span class="comment">// if the argument is a constant, the result is a constant</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		if x.mode == constant_ {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			if id == _Real {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>				x.val = constant.Real(x.val)
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>			} else {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>				x.val = constant.Imag(x.val)
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		} else {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			x.mode = value
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		if check.recordTypes() &amp;&amp; x.mode != constant_ {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(resTyp, x.typ))
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		x.typ = resTyp
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	case _Make:
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		<span class="comment">// make(T, n)</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		<span class="comment">// make(T, n, m)</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		<span class="comment">// (no argument evaluated yet)</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		arg0 := argList[0]
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		T := check.varType(arg0)
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		if !isValid(T) {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>			return
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		}
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		var min int <span class="comment">// minimum number of arguments</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		switch coreType(T).(type) {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		case *Slice:
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			min = 2
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		case *Map, *Chan:
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>			min = 1
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		case nil:
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			check.errorf(arg0, InvalidMake, invalidArg+&#34;cannot make %s: no core type&#34;, arg0)
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>			return
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		default:
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			check.errorf(arg0, InvalidMake, invalidArg+&#34;cannot make %s; type must be slice, map, or channel&#34;, arg0)
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>			return
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		if nargs &lt; min || min+1 &lt; nargs {
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			check.errorf(call, WrongArgCount, invalidOp+&#34;%v expects %d or %d arguments; found %d&#34;, call, min, min+1, nargs)
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			return
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		types := []Type{T}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		var sizes []int64 <span class="comment">// constant integer arguments, if any</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		for _, arg := range argList[1:] {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			typ, size := check.index(arg, -1) <span class="comment">// ok to continue with typ == Typ[Invalid]</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			types = append(types, typ)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>			if size &gt;= 0 {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>				sizes = append(sizes, size)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		if len(sizes) == 2 &amp;&amp; sizes[0] &gt; sizes[1] {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>			check.error(argList[1], SwappedMakeArgs, invalidArg+&#34;length and capacity swapped&#34;)
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>			<span class="comment">// safe to continue</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		x.mode = value
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		x.typ = T
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(x.typ, types...))
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		}
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	case _Max, _Min:
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		<span class="comment">// max(x, ...)</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		<span class="comment">// min(x, ...)</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		check.verifyVersionf(call.Fun, go1_21, bin.name)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		op := token.LSS
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		if id == _Max {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			op = token.GTR
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		}
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		for i, a := range args {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>			if a.mode == invalid {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>				return
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>			}
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>			if !allOrdered(a.typ) {
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>				check.errorf(a, InvalidMinMaxOperand, invalidArg+&#34;%s cannot be ordered&#34;, a)
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>				return
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>			}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>			<span class="comment">// The first argument is already in x and there&#39;s nothing left to do.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>			if i &gt; 0 {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>				check.matchTypes(x, a)
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>				if x.mode == invalid {
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>					return
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>				}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>				if !Identical(x.typ, a.typ) {
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>					check.errorf(a, MismatchedTypes, invalidArg+&#34;mismatched types %s (previous argument) and %s (type of %s)&#34;, x.typ, a.typ, a.expr)
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>					return
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>				}
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>				if x.mode == constant_ &amp;&amp; a.mode == constant_ {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>					if constant.Compare(a.val, op, x.val) {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>						*x = *a
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>					}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>				} else {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>					x.mode = value
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>				}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>			}
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		<span class="comment">// If nargs == 1, make sure x.mode is either a value or a constant.</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		if x.mode != constant_ {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>			x.mode = value
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>			<span class="comment">// A value must not be untyped.</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>			check.assignment(x, &amp;emptyInterface, &#34;argument to &#34;+bin.name)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>			if x.mode == invalid {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>				return
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>			}
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		}
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		<span class="comment">// Use the final type computed above for all arguments.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		for _, a := range args {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>			check.updateExprType(a.expr, x.typ, true)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		if check.recordTypes() &amp;&amp; x.mode != constant_ {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>			types := make([]Type, nargs)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>			for i := range types {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>				types[i] = x.typ
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(x.typ, types...))
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		}
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	case _New:
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		<span class="comment">// new(T)</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		<span class="comment">// (no argument evaluated yet)</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		T := check.varType(argList[0])
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		if !isValid(T) {
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>			return
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		}
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		x.mode = value
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		x.typ = &amp;Pointer{base: T}
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(x.typ, T))
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		}
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	case _Panic:
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		<span class="comment">// panic(x)</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		<span class="comment">// record panic call if inside a function with result parameters</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		<span class="comment">// (for use in Checker.isTerminating)</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		if check.sig != nil &amp;&amp; check.sig.results.Len() &gt; 0 {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>			<span class="comment">// function has result parameters</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>			p := check.isPanic
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>			if p == nil {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>				<span class="comment">// allocate lazily</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>				p = make(map[*ast.CallExpr]bool)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>				check.isPanic = p
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>			}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>			p[call] = true
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		}
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		check.assignment(x, &amp;emptyInterface, &#34;argument to panic&#34;)
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>			return
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>		x.mode = novalue
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(nil, &amp;emptyInterface))
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		}
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	case _Print, _Println:
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		<span class="comment">// print(x, y, ...)</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		<span class="comment">// println(x, y, ...)</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		var params []Type
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		if nargs &gt; 0 {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>			params = make([]Type, nargs)
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>			for i, a := range args {
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>				check.assignment(a, nil, &#34;argument to &#34;+predeclaredFuncs[id].name)
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>				if a.mode == invalid {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>					return
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>				}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>				params[i] = a.typ
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>			}
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		}
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		x.mode = novalue
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(nil, params...))
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		}
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	case _Recover:
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		<span class="comment">// recover() interface{}</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		x.mode = value
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		x.typ = &amp;emptyInterface
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(x.typ))
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		}
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	case _Add:
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		<span class="comment">// unsafe.Add(ptr unsafe.Pointer, len IntegerType) unsafe.Pointer</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		check.verifyVersionf(call.Fun, go1_17, &#34;unsafe.Add&#34;)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		check.assignment(x, Typ[UnsafePointer], &#34;argument to unsafe.Add&#34;)
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>			return
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		}
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		y := args[1]
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>		if !check.isValidIndex(y, InvalidUnsafeAdd, &#34;length&#34;, true) {
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>			return
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		}
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		x.mode = value
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		x.typ = Typ[UnsafePointer]
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(x.typ, x.typ, y.typ))
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		}
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	case _Alignof:
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		<span class="comment">// unsafe.Alignof(x T) uintptr</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		check.assignment(x, nil, &#34;argument to unsafe.Alignof&#34;)
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>			return
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		}
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		if hasVarSize(x.typ, nil) {
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>			x.mode = value
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>			if check.recordTypes() {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>				check.recordBuiltinType(call.Fun, makeSig(Typ[Uintptr], x.typ))
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>			}
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		} else {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			x.mode = constant_
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			x.val = constant.MakeInt64(check.conf.alignof(x.typ))
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>			<span class="comment">// result is constant - no need to record signature</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		}
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		x.typ = Typ[Uintptr]
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	case _Offsetof:
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		<span class="comment">// unsafe.Offsetof(x T) uintptr, where x must be a selector</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		<span class="comment">// (no argument evaluated yet)</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		arg0 := argList[0]
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>		selx, _ := unparen(arg0).(*ast.SelectorExpr)
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		if selx == nil {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>			check.errorf(arg0, BadOffsetofSyntax, invalidArg+&#34;%s is not a selector expression&#34;, arg0)
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			check.use(arg0)
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>			return
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		}
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		check.expr(nil, x, selx.X)
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>			return
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>		}
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		base := derefStructPtr(x.typ)
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		sel := selx.Sel.Name
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		obj, index, indirect := LookupFieldOrMethod(base, false, check.pkg, sel)
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		switch obj.(type) {
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		case nil:
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>			check.errorf(x, MissingFieldOrMethod, invalidArg+&#34;%s has no single field %s&#34;, base, sel)
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>			return
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		case *Func:
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) Using derefStructPtr may result in methods being found</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>			<span class="comment">// that don&#39;t actually exist. An error either way, but the error</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>			<span class="comment">// message is confusing. See: https://play.golang.org/p/al75v23kUy ,</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>			<span class="comment">// but go/types reports: &#34;invalid argument: x.m is a method value&#34;.</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>			check.errorf(arg0, InvalidOffsetof, invalidArg+&#34;%s is a method value&#34;, arg0)
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>			return
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		}
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>		if indirect {
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>			check.errorf(x, InvalidOffsetof, invalidArg+&#34;field %s is embedded via a pointer in %s&#34;, sel, base)
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			return
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		}
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) Should we pass x.typ instead of base (and have indirect report if derefStructPtr indirected)?</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		check.recordSelection(selx, FieldVal, base, obj, index, false)
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>		<span class="comment">// record the selector expression (was bug - go.dev/issue/47895)</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		{
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>			mode := value
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>			if x.mode == variable || indirect {
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>				mode = variable
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>			}
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>			check.record(&amp;operand{mode, selx, obj.Type(), nil, 0})
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		}
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>		<span class="comment">// The field offset is considered a variable even if the field is declared before</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>		<span class="comment">// the part of the struct which is variable-sized. This makes both the rules</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		<span class="comment">// simpler and also permits (or at least doesn&#39;t prevent) a compiler from re-</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>		<span class="comment">// arranging struct fields if it wanted to.</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		if hasVarSize(base, nil) {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>			x.mode = value
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>			if check.recordTypes() {
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>				check.recordBuiltinType(call.Fun, makeSig(Typ[Uintptr], obj.Type()))
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>			}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		} else {
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>			offs := check.conf.offsetof(base, index)
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>			if offs &lt; 0 {
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>				check.errorf(x, TypeTooLarge, &#34;%s is too large&#34;, x)
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>				return
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>			}
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>			x.mode = constant_
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>			x.val = constant.MakeInt64(offs)
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>			<span class="comment">// result is constant - no need to record signature</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		}
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>		x.typ = Typ[Uintptr]
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	case _Sizeof:
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>		<span class="comment">// unsafe.Sizeof(x T) uintptr</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		check.assignment(x, nil, &#34;argument to unsafe.Sizeof&#34;)
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>			return
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		}
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		if hasVarSize(x.typ, nil) {
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>			x.mode = value
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>			if check.recordTypes() {
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>				check.recordBuiltinType(call.Fun, makeSig(Typ[Uintptr], x.typ))
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>			}
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>		} else {
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>			size := check.conf.sizeof(x.typ)
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>			if size &lt; 0 {
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>				check.errorf(x, TypeTooLarge, &#34;%s is too large&#34;, x)
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>				return
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>			}
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>			x.mode = constant_
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>			x.val = constant.MakeInt64(size)
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>			<span class="comment">// result is constant - no need to record signature</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>		}
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>		x.typ = Typ[Uintptr]
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	case _Slice:
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		<span class="comment">// unsafe.Slice(ptr *T, len IntegerType) []T</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		check.verifyVersionf(call.Fun, go1_17, &#34;unsafe.Slice&#34;)
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		ptr, _ := coreType(x.typ).(*Pointer)
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		if ptr == nil {
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>			check.errorf(x, InvalidUnsafeSlice, invalidArg+&#34;%s is not a pointer&#34;, x)
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>			return
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>		}
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>		y := args[1]
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		if !check.isValidIndex(y, InvalidUnsafeSlice, &#34;length&#34;, false) {
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>			return
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>		}
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>		x.mode = value
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		x.typ = NewSlice(ptr.base)
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(x.typ, ptr, y.typ))
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		}
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>	case _SliceData:
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		<span class="comment">// unsafe.SliceData(slice []T) *T</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>		check.verifyVersionf(call.Fun, go1_20, &#34;unsafe.SliceData&#34;)
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>		slice, _ := coreType(x.typ).(*Slice)
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		if slice == nil {
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>			check.errorf(x, InvalidUnsafeSliceData, invalidArg+&#34;%s is not a slice&#34;, x)
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>			return
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		}
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>		x.mode = value
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		x.typ = NewPointer(slice.elem)
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(x.typ, slice))
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		}
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	case _String:
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		<span class="comment">// unsafe.String(ptr *byte, len IntegerType) string</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		check.verifyVersionf(call.Fun, go1_20, &#34;unsafe.String&#34;)
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		check.assignment(x, NewPointer(universeByte), &#34;argument to unsafe.String&#34;)
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>			return
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		}
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>		y := args[1]
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>		if !check.isValidIndex(y, InvalidUnsafeString, &#34;length&#34;, false) {
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			return
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		}
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>		x.mode = value
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		x.typ = Typ[String]
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(x.typ, NewPointer(universeByte), y.typ))
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		}
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	case _StringData:
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>		<span class="comment">// unsafe.StringData(str string) *byte</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		check.verifyVersionf(call.Fun, go1_20, &#34;unsafe.StringData&#34;)
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		check.assignment(x, Typ[String], &#34;argument to unsafe.StringData&#34;)
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>			return
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		}
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		x.mode = value
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		x.typ = NewPointer(universeByte)
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		if check.recordTypes() {
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>			check.recordBuiltinType(call.Fun, makeSig(x.typ, Typ[String]))
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		}
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	case _Assert:
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>		<span class="comment">// assert(pred) causes a typechecker error if pred is false.</span>
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>		<span class="comment">// The result of assert is the value of pred if there is no error.</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		<span class="comment">// Note: assert is only available in self-test mode.</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		if x.mode != constant_ || !isBoolean(x.typ) {
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>			check.errorf(x, Test, invalidArg+&#34;%s is not a boolean constant&#34;, x)
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>			return
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>		}
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>		if x.val.Kind() != constant.Bool {
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>			check.errorf(x, Test, &#34;internal error: value of %s should be a boolean constant&#34;, x)
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>			return
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>		}
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>		if !constant.BoolVal(x.val) {
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>			check.errorf(call, Test, &#34;%v failed&#34;, call)
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>			<span class="comment">// compile-time assertion failure - safe to continue</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>		}
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>		<span class="comment">// result is constant - no need to record signature</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	case _Trace:
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>		<span class="comment">// trace(x, y, z, ...) dumps the positions, expressions, and</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>		<span class="comment">// values of its arguments. The result of trace is the value</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>		<span class="comment">// of the first argument.</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>		<span class="comment">// Note: trace is only available in self-test mode.</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>		<span class="comment">// (no argument evaluated yet)</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>		if nargs == 0 {
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>			check.dump(&#34;%v: trace() without arguments&#34;, call.Pos())
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>			x.mode = novalue
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>			break
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>		}
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		var t operand
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>		x1 := x
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		for _, arg := range argList {
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>			check.rawExpr(nil, x1, arg, nil, false) <span class="comment">// permit trace for types, e.g.: new(trace(T))</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>			check.dump(&#34;%v: %s&#34;, x1.Pos(), x1)
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>			x1 = &amp;t <span class="comment">// use incoming x only for first argument</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>		}
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>		if x.mode == invalid {
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>			return
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>		}
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>		<span class="comment">// trace is only available in test mode - no need to record signature</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	default:
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>		unreachable()
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	}
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	assert(x.mode != invalid)
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	return true
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>}
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span><span class="comment">// hasVarSize reports if the size of type t is variable due to type parameters</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span><span class="comment">// or if the type is infinitely-sized due to a cycle for which the type has not</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span><span class="comment">// yet been checked.</span>
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>func hasVarSize(t Type, seen map[*Named]bool) (varSized bool) {
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	<span class="comment">// Cycles are only possible through *Named types.</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	<span class="comment">// The seen map is used to detect cycles and track</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	<span class="comment">// the results of previously seen types.</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	if named := asNamed(t); named != nil {
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>		if v, ok := seen[named]; ok {
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>			return v
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>		}
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		if seen == nil {
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>			seen = make(map[*Named]bool)
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>		}
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>		seen[named] = true <span class="comment">// possibly cyclic until proven otherwise</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>		defer func() {
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>			seen[named] = varSized <span class="comment">// record final determination for named</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>		}()
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>	}
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	switch u := under(t).(type) {
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	case *Array:
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>		return hasVarSize(u.elem, seen)
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>	case *Struct:
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		for _, f := range u.fields {
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>			if hasVarSize(f.typ, seen) {
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>				return true
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>			}
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>		}
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>	case *Interface:
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>		return isTypeParam(t)
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>	case *Named, *Union:
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>		unreachable()
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>	}
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>	return false
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>}
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span><span class="comment">// applyTypeFunc applies f to x. If x is a type parameter,</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span><span class="comment">// the result is a type parameter constrained by a new</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span><span class="comment">// interface bound. The type bounds for that interface</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span><span class="comment">// are computed by applying f to each of the type bounds</span>
<span id="L959" class="ln">   959&nbsp;&nbsp;</span><span class="comment">// of x. If any of these applications of f return nil,</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span><span class="comment">// applyTypeFunc returns nil.</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span><span class="comment">// If x is not a type parameter, the result is f(x).</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>func (check *Checker) applyTypeFunc(f func(Type) Type, x *operand, id builtinId) Type {
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>	if tp, _ := x.typ.(*TypeParam); tp != nil {
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>		<span class="comment">// Test if t satisfies the requirements for the argument</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>		<span class="comment">// type and collect possible result types at the same time.</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		var terms []*Term
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>		if !tp.is(func(t *term) bool {
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>			if t == nil {
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>				return false
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>			}
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>			if r := f(t.typ); r != nil {
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>				terms = append(terms, NewTerm(t.tilde, r))
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>				return true
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>			}
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>			return false
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>		}) {
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>			return nil
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>		}
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>		<span class="comment">// We can type-check this fine but we&#39;re introducing a synthetic</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>		<span class="comment">// type parameter for the result. It&#39;s not clear what the API</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		<span class="comment">// implications are here. Report an error for 1.18 (see go.dev/issue/50912),</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		<span class="comment">// but continue type-checking.</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>		var code Code
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>		switch id {
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>		case _Real:
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>			code = InvalidReal
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>		case _Imag:
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>			code = InvalidImag
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>		case _Complex:
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>			code = InvalidComplex
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>		default:
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>			unreachable()
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>		}
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		check.softErrorf(x, code, &#34;%s not supported as argument to %s for go1.18 (see go.dev/issue/50937)&#34;, x, predeclaredFuncs[id].name)
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>		<span class="comment">// Construct a suitable new type parameter for the result type.</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>		<span class="comment">// The type parameter is placed in the current package so export/import</span>
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>		<span class="comment">// works as expected.</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>		tpar := NewTypeName(nopos, check.pkg, tp.obj.name, nil)
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>		ptyp := check.newTypeParam(tpar, NewInterfaceType(nil, []Type{NewUnion(terms)})) <span class="comment">// assigns type to tpar as a side-effect</span>
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>		ptyp.index = tp.index
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>		return ptyp
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	}
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	return f(x.typ)
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>}
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span><span class="comment">// makeSig makes a signature for the given argument and result types.</span>
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span><span class="comment">// Default types are used for untyped arguments, and res may be nil.</span>
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>func makeSig(res Type, args ...Type) *Signature {
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	list := make([]*Var, len(args))
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	for i, param := range args {
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>		list[i] = NewVar(nopos, nil, &#34;&#34;, Default(param))
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	}
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	params := NewTuple(list...)
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	var result *Tuple
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>	if res != nil {
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>		assert(!isUntyped(res))
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		result = NewTuple(NewVar(nopos, nil, &#34;&#34;, res))
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	}
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	return &amp;Signature{params: params, results: result}
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>}
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span><span class="comment">// arrayPtrDeref returns A if typ is of the form *A and A is an array;</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span><span class="comment">// otherwise it returns typ.</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>func arrayPtrDeref(typ Type) Type {
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	if p, ok := typ.(*Pointer); ok {
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>		if a, _ := under(p.base).(*Array); a != nil {
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>			return a
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>		}
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>	}
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>	return typ
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>}
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>func unparen(e ast.Expr) ast.Expr { return ast.Unparen(e) }
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>
</pre><p><a href="builtins.go?m=text">View as plain text</a></p>

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
