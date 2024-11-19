<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/operand.go - Go Documentation Server</title>

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
<a href="operand.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">operand.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file defines operands and associated operations.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package types
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;go/constant&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// An operandMode specifies the (addressing) mode of an operand.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type operandMode byte
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>const (
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	invalid   operandMode = iota <span class="comment">// operand is invalid</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	novalue                      <span class="comment">// operand represents no value (result of a function call w/o result)</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	builtin                      <span class="comment">// operand is a built-in function</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	typexpr                      <span class="comment">// operand is a type</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	constant_                    <span class="comment">// operand is a constant; the operand&#39;s typ is a Basic type</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	variable                     <span class="comment">// operand is an addressable variable</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	mapindex                     <span class="comment">// operand is a map index expression (acts like a variable on lhs, commaok on rhs of an assignment)</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	value                        <span class="comment">// operand is a computed value</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	commaok                      <span class="comment">// like value, but operand may be used in a comma,ok expression</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	commaerr                     <span class="comment">// like commaok, but second value is error, not boolean</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	cgofunc                      <span class="comment">// operand is a cgo function</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>var operandModeString = [...]string{
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	invalid:   &#34;invalid operand&#34;,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	novalue:   &#34;no value&#34;,
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	builtin:   &#34;built-in&#34;,
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	typexpr:   &#34;type&#34;,
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	constant_: &#34;constant&#34;,
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	variable:  &#34;variable&#34;,
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	mapindex:  &#34;map index expression&#34;,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	value:     &#34;value&#34;,
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	commaok:   &#34;comma, ok expression&#34;,
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	commaerr:  &#34;comma, error expression&#34;,
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	cgofunc:   &#34;cgo function&#34;,
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// An operand represents an intermediate value during type checking.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// Operands have an (addressing) mode, the expression evaluating to</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// the operand, the operand&#39;s type, a value for constants, and an id</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// for built-in functions.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// The zero value of operand is a ready to use invalid operand.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>type operand struct {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	mode operandMode
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	expr ast.Expr
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	typ  Type
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	val  constant.Value
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	id   builtinId
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// Pos returns the position of the expression corresponding to x.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// If x is invalid the position is nopos.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>func (x *operand) Pos() token.Pos {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// x.expr may not be set if x is invalid</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if x.expr == nil {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		return nopos
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	return x.expr.Pos()
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// Operand string formats</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// (not all &#34;untyped&#34; cases can appear due to the type system,</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// but they fall out naturally here)</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// mode       format</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// invalid    &lt;expr&gt; (               &lt;mode&gt;                    )</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// novalue    &lt;expr&gt; (               &lt;mode&gt;                    )</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// builtin    &lt;expr&gt; (               &lt;mode&gt;                    )</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// typexpr    &lt;expr&gt; (               &lt;mode&gt;                    )</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// constant   &lt;expr&gt; (&lt;untyped kind&gt; &lt;mode&gt;                    )</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// constant   &lt;expr&gt; (               &lt;mode&gt;       of type &lt;typ&gt;)</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// constant   &lt;expr&gt; (&lt;untyped kind&gt; &lt;mode&gt; &lt;val&gt;              )</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// constant   &lt;expr&gt; (               &lt;mode&gt; &lt;val&gt; of type &lt;typ&gt;)</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// variable   &lt;expr&gt; (&lt;untyped kind&gt; &lt;mode&gt;                    )</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// variable   &lt;expr&gt; (               &lt;mode&gt;       of type &lt;typ&gt;)</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// mapindex   &lt;expr&gt; (&lt;untyped kind&gt; &lt;mode&gt;                    )</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// mapindex   &lt;expr&gt; (               &lt;mode&gt;       of type &lt;typ&gt;)</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// value      &lt;expr&gt; (&lt;untyped kind&gt; &lt;mode&gt;                    )</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// value      &lt;expr&gt; (               &lt;mode&gt;       of type &lt;typ&gt;)</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// commaok    &lt;expr&gt; (&lt;untyped kind&gt; &lt;mode&gt;                    )</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// commaok    &lt;expr&gt; (               &lt;mode&gt;       of type &lt;typ&gt;)</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// commaerr   &lt;expr&gt; (&lt;untyped kind&gt; &lt;mode&gt;                    )</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// commaerr   &lt;expr&gt; (               &lt;mode&gt;       of type &lt;typ&gt;)</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// cgofunc    &lt;expr&gt; (&lt;untyped kind&gt; &lt;mode&gt;                    )</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// cgofunc    &lt;expr&gt; (               &lt;mode&gt;       of type &lt;typ&gt;)</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>func operandString(x *operand, qf Qualifier) string {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// special-case nil</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	if x.mode == value &amp;&amp; x.typ == Typ[UntypedNil] {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		return &#34;nil&#34;
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	var buf bytes.Buffer
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	var expr string
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	if x.expr != nil {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		expr = ExprString(x.expr)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	} else {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		switch x.mode {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		case builtin:
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			expr = predeclaredFuncs[x.id].name
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		case typexpr:
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			expr = TypeString(x.typ, qf)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		case constant_:
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			expr = x.val.String()
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// &lt;expr&gt; (</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if expr != &#34;&#34; {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		buf.WriteString(expr)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		buf.WriteString(&#34; (&#34;)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// &lt;untyped kind&gt;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	hasType := false
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	switch x.mode {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	case invalid, novalue, builtin, typexpr:
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// no type</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	default:
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		<span class="comment">// should have a type, but be cautious (don&#39;t crash during printing)</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		if x.typ != nil {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			if isUntyped(x.typ) {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>				buf.WriteString(x.typ.(*Basic).name)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>				buf.WriteByte(&#39; &#39;)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>				break
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			hasType = true
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// &lt;mode&gt;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	buf.WriteString(operandModeString[x.mode])
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// &lt;val&gt;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if x.mode == constant_ {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		if s := x.val.String(); s != expr {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			buf.WriteByte(&#39; &#39;)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			buf.WriteString(s)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// &lt;typ&gt;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if hasType {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		if isValid(x.typ) {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			var intro string
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			if isGeneric(x.typ) {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>				intro = &#34; of generic type &#34;
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			} else {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>				intro = &#34; of type &#34;
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			buf.WriteString(intro)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			WriteType(&amp;buf, x.typ, qf)
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			if tpar, _ := x.typ.(*TypeParam); tpar != nil {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>				buf.WriteString(&#34; constrained by &#34;)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>				WriteType(&amp;buf, tpar.bound, qf) <span class="comment">// do not compute interface type sets here</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>				<span class="comment">// If we have the type set and it&#39;s empty, say so for better error messages.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>				if hasEmptyTypeset(tpar) {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>					buf.WriteString(&#34; with empty type set&#34;)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>				}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		} else {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			buf.WriteString(&#34; with invalid type&#34;)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// )</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if expr != &#34;&#34; {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		buf.WriteByte(&#39;)&#39;)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	return buf.String()
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>func (x *operand) String() string {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	return operandString(x, nil)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// setConst sets x to the untyped constant for literal lit.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>func (x *operand) setConst(tok token.Token, lit string) {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	var kind BasicKind
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	switch tok {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	case token.INT:
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		kind = UntypedInt
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	case token.FLOAT:
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		kind = UntypedFloat
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	case token.IMAG:
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		kind = UntypedComplex
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	case token.CHAR:
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		kind = UntypedRune
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	case token.STRING:
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		kind = UntypedString
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	default:
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		unreachable()
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	val := constant.MakeFromLiteral(lit, tok, 0)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	if val.Kind() == constant.Unknown {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		x.mode = invalid
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		x.typ = Typ[Invalid]
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		return
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	x.mode = constant_
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	x.typ = Typ[kind]
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	x.val = val
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// isNil reports whether x is the (untyped) nil value.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>func (x *operand) isNil() bool { return x.mode == value &amp;&amp; x.typ == Typ[UntypedNil] }
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// assignableTo reports whether x is assignable to a variable of type T. If the</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// result is false and a non-nil cause is provided, it may be set to a more</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// detailed explanation of the failure (result != &#34;&#34;). The returned error code</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// is only valid if the (first) result is false. The check parameter may be nil</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// if assignableTo is invoked through an exported API call, i.e., when all</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">// methods have been type-checked.</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>func (x *operand) assignableTo(check *Checker, T Type, cause *string) (bool, Code) {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	if x.mode == invalid || !isValid(T) {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		return true, 0 <span class="comment">// avoid spurious errors</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	V := x.typ
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// x&#39;s type is identical to T</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if Identical(V, T) {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		return true, 0
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	Vu := under(V)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	Tu := under(T)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	Vp, _ := V.(*TypeParam)
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	Tp, _ := T.(*TypeParam)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// x is an untyped value representable by a value of type T.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	if isUntyped(Vu) {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		assert(Vp == nil)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		if Tp != nil {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			<span class="comment">// T is a type parameter: x is assignable to T if it is</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			<span class="comment">// representable by each specific type in the type set of T.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			return Tp.is(func(t *term) bool {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>				if t == nil {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>					return false
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>				}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>				<span class="comment">// A term may be a tilde term but the underlying</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				<span class="comment">// type of an untyped value doesn&#39;t change so we</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>				<span class="comment">// don&#39;t need to do anything special.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>				newType, _, _ := check.implicitTypeAndValue(x, t.typ)
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>				return newType != nil
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			}), IncompatibleAssign
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		newType, _, _ := check.implicitTypeAndValue(x, T)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		return newType != nil, IncompatibleAssign
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	<span class="comment">// Vu is typed</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	<span class="comment">// x&#39;s type V and T have identical underlying types</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	<span class="comment">// and at least one of V or T is not a named type</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	<span class="comment">// and neither V nor T is a type parameter.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	if Identical(Vu, Tu) &amp;&amp; (!hasName(V) || !hasName(T)) &amp;&amp; Vp == nil &amp;&amp; Tp == nil {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		return true, 0
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	<span class="comment">// T is an interface type, but not a type parameter, and V implements T.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	<span class="comment">// Also handle the case where T is a pointer to an interface so that we get</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	<span class="comment">// the Checker.implements error cause.</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	if _, ok := Tu.(*Interface); ok &amp;&amp; Tp == nil || isInterfacePtr(Tu) {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		if check.implements(x.Pos(), V, T, false, cause) {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			return true, 0
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		<span class="comment">// V doesn&#39;t implement T but V may still be assignable to T if V</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		<span class="comment">// is a type parameter; do not report an error in that case yet.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		if Vp == nil {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			return false, InvalidIfaceAssign
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		if cause != nil {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			*cause = &#34;&#34;
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// If V is an interface, check if a missing type assertion is the problem.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	if Vi, _ := Vu.(*Interface); Vi != nil &amp;&amp; Vp == nil {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		if check.implements(x.Pos(), T, V, false, nil) {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			<span class="comment">// T implements V, so give hint about type assertion.</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			if cause != nil {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>				*cause = &#34;need type assertion&#34;
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			return false, IncompatibleAssign
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	<span class="comment">// x is a bidirectional channel value, T is a channel</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	<span class="comment">// type, x&#39;s type V and T have identical element types,</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	<span class="comment">// and at least one of V or T is not a named type.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	if Vc, ok := Vu.(*Chan); ok &amp;&amp; Vc.dir == SendRecv {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		if Tc, ok := Tu.(*Chan); ok &amp;&amp; Identical(Vc.elem, Tc.elem) {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>			return !hasName(V) || !hasName(T), InvalidChanAssign
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	<span class="comment">// optimization: if we don&#39;t have type parameters, we&#39;re done</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	if Vp == nil &amp;&amp; Tp == nil {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		return false, IncompatibleAssign
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	errorf := func(format string, args ...any) {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		if check != nil &amp;&amp; cause != nil {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			msg := check.sprintf(format, args...)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			if *cause != &#34;&#34; {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>				msg += &#34;\n\t&#34; + *cause
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			*cause = msg
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">// x&#39;s type V is not a named type and T is a type parameter, and</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// x is assignable to each specific type in T&#39;s type set.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	if !hasName(V) &amp;&amp; Tp != nil {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		ok := false
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		code := IncompatibleAssign
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		Tp.is(func(T *term) bool {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>			if T == nil {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>				return false <span class="comment">// no specific types</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			ok, code = x.assignableTo(check, T.typ, cause)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			if !ok {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>				errorf(&#34;cannot assign %s to %s (in %s)&#34;, x.typ, T.typ, Tp)
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>				return false
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>			return true
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		})
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		return ok, code
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	<span class="comment">// x&#39;s type V is a type parameter and T is not a named type,</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	<span class="comment">// and values x&#39; of each specific type in V&#39;s type set are</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	<span class="comment">// assignable to T.</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	if Vp != nil &amp;&amp; !hasName(T) {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		x := *x <span class="comment">// don&#39;t clobber outer x</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		ok := false
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		code := IncompatibleAssign
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		Vp.is(func(V *term) bool {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			if V == nil {
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>				return false <span class="comment">// no specific types</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>			}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>			x.typ = V.typ
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			ok, code = x.assignableTo(check, T, cause)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			if !ok {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>				errorf(&#34;cannot assign %s (in %s) to %s&#34;, V.typ, Vp, T)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>				return false
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			return true
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		})
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		return ok, code
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	return false, IncompatibleAssign
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
</pre><p><a href="operand.go?m=text">View as plain text</a></p>

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
