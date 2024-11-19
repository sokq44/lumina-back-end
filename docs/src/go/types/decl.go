<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/decl.go - Go Documentation Server</title>

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
<a href="decl.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">decl.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/types">go/types</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2014 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package types
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;go/constant&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>func (check *Checker) reportAltDecl(obj Object) {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	if pos := obj.Pos(); pos.IsValid() {
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>		<span class="comment">// We use &#34;other&#34; rather than &#34;previous&#34; here because</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>		<span class="comment">// the first declaration seen may not be textually</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>		<span class="comment">// earlier in the source.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>		check.errorf(obj, DuplicateDecl, &#34;\tother declaration of %s&#34;, obj.Name()) <span class="comment">// secondary error, \t indented</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	}
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>func (check *Checker) declare(scope *Scope, id *ast.Ident, obj Object, pos token.Pos) {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;The blank identifier, represented by the underscore</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// character _, may be used in a declaration like any other</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// identifier but the declaration does not introduce a new</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// binding.&#34;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	if obj.Name() != &#34;_&#34; {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		if alt := scope.Insert(obj); alt != nil {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>			check.errorf(obj, DuplicateDecl, &#34;%s redeclared in this block&#34;, obj.Name())
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>			check.reportAltDecl(alt)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>			return
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		obj.setScopePos(pos)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	if id != nil {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		check.recordDef(id, obj)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// pathString returns a string of the form a-&gt;b-&gt; ... -&gt;g for a path [a, b, ... g].</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func pathString(path []Object) string {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	var s string
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	for i, p := range path {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		if i &gt; 0 {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			s += &#34;-&gt;&#34;
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		s += p.Name()
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	return s
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// objDecl type-checks the declaration of obj in its respective (file) environment.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// For the meaning of def, see Checker.definedType, in typexpr.go.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>func (check *Checker) objDecl(obj Object, def *TypeName) {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	if check.conf._Trace &amp;&amp; obj.Type() == nil {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		if check.indent == 0 {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			fmt.Println() <span class="comment">// empty line between top-level objects for readability</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		check.trace(obj.Pos(), &#34;-- checking %s (%s, objPath = %s)&#34;, obj, obj.color(), pathString(check.objPath))
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		check.indent++
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		defer func() {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			check.indent--
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			check.trace(obj.Pos(), &#34;=&gt; %s (%s)&#34;, obj, obj.color())
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		}()
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// Checking the declaration of obj means inferring its type</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// (and possibly its value, for constants).</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// An object&#39;s type (and thus the object) may be in one of</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// three states which are expressed by colors:</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// - an object whose type is not yet known is painted white (initial color)</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// - an object whose type is in the process of being inferred is painted grey</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// - an object whose type is fully inferred is painted black</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// During type inference, an object&#39;s color changes from white to grey</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// to black (pre-declared objects are painted black from the start).</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// A black object (i.e., its type) can only depend on (refer to) other black</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// ones. White and grey objects may depend on white and black objects.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// A dependency on a grey object indicates a cycle which may or may not be</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// valid.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// When objects turn grey, they are pushed on the object path (a stack);</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// they are popped again when they turn black. Thus, if a grey object (a</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// cycle) is encountered, it is on the object path, and all the objects</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// it depends on are the remaining objects on that path. Color encoding</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// is such that the color value of a grey object indicates the index of</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// that object in the object path.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// During type-checking, white objects may be assigned a type without</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// traversing through objDecl; e.g., when initializing constants and</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// variables. Update the colors of those objects here (rather than</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// everywhere where we set the type) to satisfy the color invariants.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if obj.color() == white &amp;&amp; obj.Type() != nil {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		obj.setColor(black)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		return
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	switch obj.color() {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	case white:
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		assert(obj.Type() == nil)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		<span class="comment">// All color values other than white and black are considered grey.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		<span class="comment">// Because black and white are &lt; grey, all values &gt;= grey are grey.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">// Use those values to encode the object&#39;s index into the object path.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		obj.setColor(grey + color(check.push(obj)))
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		defer func() {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			check.pop().setColor(black)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		}()
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	case black:
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		assert(obj.Type() != nil)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		return
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	default:
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		<span class="comment">// Color values other than white or black are considered grey.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		fallthrough
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	case grey:
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		<span class="comment">// We have a (possibly invalid) cycle.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		<span class="comment">// In the existing code, this is marked by a non-nil type</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		<span class="comment">// for the object except for constants and variables whose</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		<span class="comment">// type may be non-nil (known), or nil if it depends on the</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		<span class="comment">// not-yet known initialization value.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		<span class="comment">// In the former case, set the type to Typ[Invalid] because</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		<span class="comment">// we have an initialization cycle. The cycle error will be</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		<span class="comment">// reported later, when determining initialization order.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) Report cycle here and simplify initialization</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		<span class="comment">// order code.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		switch obj := obj.(type) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		case *Const:
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			if !check.validCycle(obj) || obj.typ == nil {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>				obj.typ = Typ[Invalid]
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		case *Var:
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			if !check.validCycle(obj) || obj.typ == nil {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>				obj.typ = Typ[Invalid]
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		case *TypeName:
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			if !check.validCycle(obj) {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>				<span class="comment">// break cycle</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>				<span class="comment">// (without this, calling underlying()</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>				<span class="comment">// below may lead to an endless loop</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>				<span class="comment">// if we have a cycle for a defined</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				<span class="comment">// (*Named) type)</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>				obj.typ = Typ[Invalid]
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		case *Func:
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			if !check.validCycle(obj) {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>				<span class="comment">// Don&#39;t set obj.typ to Typ[Invalid] here</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>				<span class="comment">// because plenty of code type-asserts that</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>				<span class="comment">// functions have a *Signature type. Grey</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>				<span class="comment">// functions have their type set to an empty</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>				<span class="comment">// signature which makes it impossible to</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>				<span class="comment">// initialize a variable with the function.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		default:
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			unreachable()
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		assert(obj.Type() != nil)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		return
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	d := check.objMap[obj]
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if d == nil {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		check.dump(&#34;%v: %s should have been declared&#34;, obj.Pos(), obj)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		unreachable()
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// save/restore current environment and set up object environment</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	defer func(env environment) {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		check.environment = env
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}(check.environment)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	check.environment = environment{
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		scope: d.file,
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// Const and var declarations must not have initialization</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// cycles. We track them by remembering the current declaration</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// in check.decl. Initialization expressions depending on other</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// consts, vars, or functions, add dependencies to the current</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// check.decl.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	switch obj := obj.(type) {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	case *Const:
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		check.decl = d <span class="comment">// new package-level const decl</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		check.constDecl(obj, d.vtyp, d.init, d.inherited)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	case *Var:
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		check.decl = d <span class="comment">// new package-level var decl</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		check.varDecl(obj, d.lhs, d.vtyp, d.init)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	case *TypeName:
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		<span class="comment">// invalid recursive types are detected via path</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		check.typeDecl(obj, d.tdecl, def)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		check.collectMethods(obj) <span class="comment">// methods can only be added to top-level types</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	case *Func:
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		<span class="comment">// functions may be recursive - no need to track dependencies</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		check.funcDecl(obj, d)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	default:
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		unreachable()
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// validCycle checks if the cycle starting with obj is valid and</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// reports an error if it is not.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>func (check *Checker) validCycle(obj Object) (valid bool) {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">// The object map contains the package scope objects and the non-interface methods.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if debug {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		info := check.objMap[obj]
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		inObjMap := info != nil &amp;&amp; (info.fdecl == nil || info.fdecl.Recv == nil) <span class="comment">// exclude methods</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		isPkgObj := obj.Parent() == check.pkg.scope
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		if isPkgObj != inObjMap {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			check.dump(&#34;%v: inconsistent object map for %s (isPkgObj = %v, inObjMap = %v)&#34;, obj.Pos(), obj, isPkgObj, inObjMap)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			unreachable()
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// Count cycle objects.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	assert(obj.color() &gt;= grey)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	start := obj.color() - grey <span class="comment">// index of obj in objPath</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	cycle := check.objPath[start:]
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	tparCycle := false <span class="comment">// if set, the cycle is through a type parameter list</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	nval := 0          <span class="comment">// number of (constant or variable) values in the cycle; valid if !generic</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	ndef := 0          <span class="comment">// number of type definitions in the cycle; valid if !generic</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>loop:
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	for _, obj := range cycle {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		switch obj := obj.(type) {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		case *Const, *Var:
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			nval++
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		case *TypeName:
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			<span class="comment">// If we reach a generic type that is part of a cycle</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			<span class="comment">// and we are in a type parameter list, we have a cycle</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			<span class="comment">// through a type parameter list, which is invalid.</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			if check.inTParamList &amp;&amp; isGeneric(obj.typ) {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>				tparCycle = true
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>				break loop
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			<span class="comment">// Determine if the type name is an alias or not. For</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			<span class="comment">// package-level objects, use the object map which</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			<span class="comment">// provides syntactic information (which doesn&#39;t rely</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			<span class="comment">// on the order in which the objects are set up). For</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			<span class="comment">// local objects, we can rely on the order, so use</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			<span class="comment">// the object&#39;s predicate.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) It would be less fragile to always access</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			<span class="comment">// the syntactic information. We should consider storing</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			<span class="comment">// this information explicitly in the object.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			var alias bool
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			if check.enableAlias {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>				alias = obj.IsAlias()
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			} else {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				if d := check.objMap[obj]; d != nil {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>					alias = d.tdecl.Assign.IsValid() <span class="comment">// package-level object</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>				} else {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>					alias = obj.IsAlias() <span class="comment">// function local object</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>				}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			if !alias {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				ndef++
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		case *Func:
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			<span class="comment">// ignored for now</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		default:
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			unreachable()
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		}
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	if check.conf._Trace {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		check.trace(obj.Pos(), &#34;## cycle detected: objPath = %s-&gt;%s (len = %d)&#34;, pathString(cycle), obj.Name(), len(cycle))
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		if tparCycle {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			check.trace(obj.Pos(), &#34;## cycle contains: generic type in a type parameter list&#34;)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		} else {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			check.trace(obj.Pos(), &#34;## cycle contains: %d values, %d type definitions&#34;, nval, ndef)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		defer func() {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			if valid {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>				check.trace(obj.Pos(), &#34;=&gt; cycle is valid&#34;)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			} else {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>				check.trace(obj.Pos(), &#34;=&gt; error: cycle is invalid&#34;)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		}()
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	if !tparCycle {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		<span class="comment">// A cycle involving only constants and variables is invalid but we</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		<span class="comment">// ignore them here because they are reported via the initialization</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		<span class="comment">// cycle check.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		if nval == len(cycle) {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			return true
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		<span class="comment">// A cycle involving only types (and possibly functions) must have at least</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		<span class="comment">// one type definition to be permitted: If there is no type definition, we</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		<span class="comment">// have a sequence of alias type names which will expand ad infinitum.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		if nval == 0 &amp;&amp; ndef &gt; 0 {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			return true
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	check.cycleError(cycle)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	return false
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">// cycleError reports a declaration cycle starting with</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// the object in cycle that is &#34;first&#34; in the source.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>func (check *Checker) cycleError(cycle []Object) {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	<span class="comment">// name returns the (possibly qualified) object name.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	<span class="comment">// This is needed because with generic types, cycles</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	<span class="comment">// may refer to imported types. See go.dev/issue/50788.</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	<span class="comment">// TODO(gri) Thus functionality is used elsewhere. Factor it out.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	name := func(obj Object) string {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		return packagePrefix(obj.Pkg(), check.qualifier) + obj.Name()
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	<span class="comment">// TODO(gri) Should we start with the last (rather than the first) object in the cycle</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	<span class="comment">//           since that is the earliest point in the source where we start seeing the</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	<span class="comment">//           cycle? That would be more consistent with other error messages.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	i := firstInSrc(cycle)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	obj := cycle[i]
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	objName := name(obj)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	<span class="comment">// If obj is a type alias, mark it as valid (not broken) in order to avoid follow-on errors.</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	tname, _ := obj.(*TypeName)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	if tname != nil &amp;&amp; tname.IsAlias() {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		<span class="comment">// If we use Alias nodes, it is initialized with Typ[Invalid].</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) Adjust this code if we initialize with nil.</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		if !check.enableAlias {
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>			check.validAlias(tname, Typ[Invalid])
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// report a more concise error for self references</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	if len(cycle) == 1 {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		if tname != nil {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			check.errorf(obj, InvalidDeclCycle, &#34;invalid recursive type: %s refers to itself&#34;, objName)
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		} else {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>			check.errorf(obj, InvalidDeclCycle, &#34;invalid cycle in declaration: %s refers to itself&#34;, objName)
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		return
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	if tname != nil {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		check.errorf(obj, InvalidDeclCycle, &#34;invalid recursive type %s&#34;, objName)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	} else {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		check.errorf(obj, InvalidDeclCycle, &#34;invalid cycle in declaration of %s&#34;, objName)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	for range cycle {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		check.errorf(obj, InvalidDeclCycle, &#34;\t%s refers to&#34;, objName) <span class="comment">// secondary error, \t indented</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		i++
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		if i &gt;= len(cycle) {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>			i = 0
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		obj = cycle[i]
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		objName = name(obj)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	check.errorf(obj, InvalidDeclCycle, &#34;\t%s&#34;, objName)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// firstInSrc reports the index of the object with the &#34;smallest&#34;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// source position in path. path must not be empty.</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>func firstInSrc(path []Object) int {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	fst, pos := 0, path[0].Pos()
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	for i, t := range path[1:] {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		if cmpPos(t.Pos(), pos) &lt; 0 {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			fst, pos = i+1, t.Pos()
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	}
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	return fst
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>type (
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	decl interface {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		node() ast.Node
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	importDecl struct{ spec *ast.ImportSpec }
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	constDecl  struct {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		spec      *ast.ValueSpec
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		iota      int
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		typ       ast.Expr
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		init      []ast.Expr
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		inherited bool
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	varDecl  struct{ spec *ast.ValueSpec }
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	typeDecl struct{ spec *ast.TypeSpec }
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	funcDecl struct{ decl *ast.FuncDecl }
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>)
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>func (d importDecl) node() ast.Node { return d.spec }
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>func (d constDecl) node() ast.Node  { return d.spec }
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>func (d varDecl) node() ast.Node    { return d.spec }
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>func (d typeDecl) node() ast.Node   { return d.spec }
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>func (d funcDecl) node() ast.Node   { return d.decl }
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>func (check *Checker) walkDecls(decls []ast.Decl, f func(decl)) {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	for _, d := range decls {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		check.walkDecl(d, f)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>func (check *Checker) walkDecl(d ast.Decl, f func(decl)) {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	switch d := d.(type) {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	case *ast.BadDecl:
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		<span class="comment">// ignore</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	case *ast.GenDecl:
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		var last *ast.ValueSpec <span class="comment">// last ValueSpec with type or init exprs seen</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		for iota, s := range d.Specs {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			switch s := s.(type) {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>			case *ast.ImportSpec:
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>				f(importDecl{s})
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			case *ast.ValueSpec:
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>				switch d.Tok {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>				case token.CONST:
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>					<span class="comment">// determine which initialization expressions to use</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>					inherited := true
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>					switch {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>					case s.Type != nil || len(s.Values) &gt; 0:
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>						last = s
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>						inherited = false
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>					case last == nil:
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>						last = new(ast.ValueSpec) <span class="comment">// make sure last exists</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>						inherited = false
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>					}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>					check.arityMatch(s, last)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>					f(constDecl{spec: s, iota: iota, typ: last.Type, init: last.Values, inherited: inherited})
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>				case token.VAR:
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>					check.arityMatch(s, nil)
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>					f(varDecl{s})
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>				default:
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>					check.errorf(s, InvalidSyntaxTree, &#34;invalid token %s&#34;, d.Tok)
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>				}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			case *ast.TypeSpec:
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>				f(typeDecl{s})
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>			default:
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>				check.errorf(s, InvalidSyntaxTree, &#34;unknown ast.Spec node %T&#34;, s)
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>			}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		}
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	case *ast.FuncDecl:
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		f(funcDecl{d})
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	default:
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		check.errorf(d, InvalidSyntaxTree, &#34;unknown ast.Decl node %T&#34;, d)
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>func (check *Checker) constDecl(obj *Const, typ, init ast.Expr, inherited bool) {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	assert(obj.typ == nil)
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	<span class="comment">// use the correct value of iota</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	defer func(iota constant.Value, errpos positioner) {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		check.iota = iota
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		check.errpos = errpos
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	}(check.iota, check.errpos)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	check.iota = obj.val
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	check.errpos = nil
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	<span class="comment">// provide valid constant value under all circumstances</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	obj.val = constant.MakeUnknown()
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	<span class="comment">// determine type, if any</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	if typ != nil {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		t := check.typ(typ)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		if !isConstType(t) {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>			<span class="comment">// don&#39;t report an error if the type is an invalid C (defined) type</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			<span class="comment">// (go.dev/issue/22090)</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>			if isValid(under(t)) {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>				check.errorf(typ, InvalidConstType, &#34;invalid constant type %s&#34;, t)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>			}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>			obj.typ = Typ[Invalid]
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			return
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		obj.typ = t
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	<span class="comment">// check initialization</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	var x operand
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	if init != nil {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		if inherited {
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			<span class="comment">// The initialization expression is inherited from a previous</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			<span class="comment">// constant declaration, and (error) positions refer to that</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			<span class="comment">// expression and not the current constant declaration. Use</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			<span class="comment">// the constant identifier position for any errors during</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			<span class="comment">// init expression evaluation since that is all we have</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			<span class="comment">// (see issues go.dev/issue/42991, go.dev/issue/42992).</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			check.errpos = atPos(obj.pos)
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		check.expr(nil, &amp;x, init)
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	check.initConst(obj, &amp;x)
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>func (check *Checker) varDecl(obj *Var, lhs []*Var, typ, init ast.Expr) {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	assert(obj.typ == nil)
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	<span class="comment">// determine type, if any</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	if typ != nil {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		obj.typ = check.varType(typ)
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		<span class="comment">// We cannot spread the type to all lhs variables if there</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		<span class="comment">// are more than one since that would mark them as checked</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		<span class="comment">// (see Checker.objDecl) and the assignment of init exprs,</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		<span class="comment">// if any, would not be checked.</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) If we have no init expr, we should distribute</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		<span class="comment">// a given type otherwise we need to re-evalate the type</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		<span class="comment">// expr for each lhs variable, leading to duplicate work.</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	<span class="comment">// check initialization</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	if init == nil {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		if typ == nil {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			<span class="comment">// error reported before by arityMatch</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>			obj.typ = Typ[Invalid]
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		return
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	if lhs == nil || len(lhs) == 1 {
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		assert(lhs == nil || lhs[0] == obj)
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		var x operand
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		check.expr(newTarget(obj.typ, obj.name), &amp;x, init)
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		check.initVar(obj, &amp;x, &#34;variable declaration&#34;)
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		return
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	}
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	if debug {
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		<span class="comment">// obj must be one of lhs</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		found := false
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		for _, lhs := range lhs {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			if obj == lhs {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>				found = true
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>				break
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>			}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		if !found {
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			panic(&#34;inconsistent lhs&#34;)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	}
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	<span class="comment">// We have multiple variables on the lhs and one init expr.</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	<span class="comment">// Make sure all variables have been given the same type if</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	<span class="comment">// one was specified, otherwise they assume the type of the</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	<span class="comment">// init expression values (was go.dev/issue/15755).</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	if typ != nil {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		for _, lhs := range lhs {
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>			lhs.typ = obj.typ
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	check.initVars(lhs, []ast.Expr{init}, nil)
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">// isImportedConstraint reports whether typ is an imported type constraint.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>func (check *Checker) isImportedConstraint(typ Type) bool {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	named := asNamed(typ)
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	if named == nil || named.obj.pkg == check.pkg || named.obj.pkg == nil {
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		return false
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	u, _ := named.under().(*Interface)
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	return u != nil &amp;&amp; !u.IsMethodSet()
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>func (check *Checker) typeDecl(obj *TypeName, tdecl *ast.TypeSpec, def *TypeName) {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	assert(obj.typ == nil)
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	var rhs Type
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	check.later(func() {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		if t := asNamed(obj.typ); t != nil { <span class="comment">// type may be invalid</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>			check.validType(t)
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		<span class="comment">// If typ is local, an error was already reported where typ is specified/defined.</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		_ = check.isImportedConstraint(rhs) &amp;&amp; check.verifyVersionf(tdecl.Type, go1_18, &#34;using type constraint %s&#34;, rhs)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	}).describef(obj, &#34;validType(%s)&#34;, obj.Name())
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	aliasDecl := tdecl.Assign.IsValid()
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	if aliasDecl &amp;&amp; tdecl.TypeParams.NumFields() != 0 {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		<span class="comment">// The parser will ensure this but we may still get an invalid AST.</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		<span class="comment">// Complain and continue as regular type definition.</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		check.error(atPos(tdecl.Assign), BadDecl, &#34;generic type cannot be alias&#34;)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		aliasDecl = false
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	}
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	<span class="comment">// alias declaration</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	if aliasDecl {
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		check.verifyVersionf(atPos(tdecl.Assign), go1_9, &#34;type aliases&#34;)
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		if check.enableAlias {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) Should be able to use nil instead of Typ[Invalid] to mark</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>			<span class="comment">//           the alias as incomplete. Currently this causes problems</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>			<span class="comment">//           with certain cycles. Investigate.</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>			alias := check.newAlias(obj, Typ[Invalid])
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>			setDefType(def, alias)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>			rhs = check.definedType(tdecl.Type, obj)
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>			assert(rhs != nil)
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			alias.fromRHS = rhs
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>			Unalias(alias) <span class="comment">// resolve alias.actual</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		} else {
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>			check.brokenAlias(obj)
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>			rhs = check.typ(tdecl.Type)
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>			check.validAlias(obj, rhs)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		return
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	<span class="comment">// type definition or generic type declaration</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	named := check.newNamed(obj, nil, nil)
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	setDefType(def, named)
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	if tdecl.TypeParams != nil {
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		check.openScope(tdecl, &#34;type parameters&#34;)
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		defer check.closeScope()
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		check.collectTypeParams(&amp;named.tparams, tdecl.TypeParams)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	<span class="comment">// determine underlying type of named</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	rhs = check.definedType(tdecl.Type, obj)
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	assert(rhs != nil)
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	named.fromRHS = rhs
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	<span class="comment">// If the underlying type was not set while type-checking the right-hand</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	<span class="comment">// side, it is invalid and an error should have been reported elsewhere.</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	if named.underlying == nil {
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		named.underlying = Typ[Invalid]
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	<span class="comment">// Disallow a lone type parameter as the RHS of a type declaration (go.dev/issue/45639).</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	<span class="comment">// We don&#39;t need this restriction anymore if we make the underlying type of a type</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	<span class="comment">// parameter its constraint interface: if the RHS is a lone type parameter, we will</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	<span class="comment">// use its underlying type (like we do for any RHS in a type declaration), and its</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	<span class="comment">// underlying type is an interface and the type declaration is well defined.</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	if isTypeParam(rhs) {
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		check.error(tdecl.Type, MisplacedTypeParam, &#34;cannot use a type parameter as RHS in type declaration&#34;)
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>		named.underlying = Typ[Invalid]
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	}
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>}
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>func (check *Checker) collectTypeParams(dst **TypeParamList, list *ast.FieldList) {
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	var tparams []*TypeParam
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	<span class="comment">// Declare type parameters up-front, with empty interface as type bound.</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	<span class="comment">// The scope of type parameters starts at the beginning of the type parameter</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	<span class="comment">// list (so we can have mutually recursive parameterized interfaces).</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	scopePos := list.Pos()
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	for _, f := range list.List {
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		tparams = check.declareTypeParams(tparams, f.Names, scopePos)
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	}
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	<span class="comment">// Set the type parameters before collecting the type constraints because</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	<span class="comment">// the parameterized type may be used by the constraints (go.dev/issue/47887).</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	<span class="comment">// Example: type T[P T[P]] interface{}</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	*dst = bindTParams(tparams)
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	<span class="comment">// Signal to cycle detection that we are in a type parameter list.</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	<span class="comment">// We can only be inside one type parameter list at any given time:</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	<span class="comment">// function closures may appear inside a type parameter list but they</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	<span class="comment">// cannot be generic, and their bodies are processed in delayed and</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	<span class="comment">// sequential fashion. Note that with each new declaration, we save</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	<span class="comment">// the existing environment and restore it when done; thus inTPList is</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	<span class="comment">// true exactly only when we are in a specific type parameter list.</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	assert(!check.inTParamList)
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	check.inTParamList = true
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	defer func() {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		check.inTParamList = false
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	}()
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	index := 0
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	for _, f := range list.List {
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		var bound Type
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		<span class="comment">// NOTE: we may be able to assert that f.Type != nil here, but this is not</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		<span class="comment">// an invariant of the AST, so we are cautious.</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		if f.Type != nil {
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>			bound = check.bound(f.Type)
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>			if isTypeParam(bound) {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>				<span class="comment">// We may be able to allow this since it is now well-defined what</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>				<span class="comment">// the underlying type and thus type set of a type parameter is.</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>				<span class="comment">// But we may need some additional form of cycle detection within</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>				<span class="comment">// type parameter lists.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>				check.error(f.Type, MisplacedTypeParam, &#34;cannot use a type parameter as constraint&#34;)
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>				bound = Typ[Invalid]
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>			}
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		} else {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>			bound = Typ[Invalid]
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		}
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		for i := range f.Names {
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>			tparams[index+i].bound = bound
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		index += len(f.Names)
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>}
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>func (check *Checker) bound(x ast.Expr) Type {
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	<span class="comment">// A type set literal of the form ~T and A|B may only appear as constraint;</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	<span class="comment">// embed it in an implicit interface so that only interface type-checking</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	<span class="comment">// needs to take care of such type expressions.</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	wrap := false
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	switch op := x.(type) {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	case *ast.UnaryExpr:
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		wrap = op.Op == token.TILDE
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		wrap = op.Op == token.OR
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	}
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	if wrap {
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		x = &amp;ast.InterfaceType{Methods: &amp;ast.FieldList{List: []*ast.Field{{Type: x}}}}
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		t := check.typ(x)
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>		<span class="comment">// mark t as implicit interface if all went well</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		if t, _ := t.(*Interface); t != nil {
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			t.implicit = true
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		}
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		return t
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	}
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	return check.typ(x)
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>}
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>func (check *Checker) declareTypeParams(tparams []*TypeParam, names []*ast.Ident, scopePos token.Pos) []*TypeParam {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	<span class="comment">// Use Typ[Invalid] for the type constraint to ensure that a type</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	<span class="comment">// is present even if the actual constraint has not been assigned</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	<span class="comment">// yet.</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	<span class="comment">// TODO(gri) Need to systematically review all uses of type parameter</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	<span class="comment">//           constraints to make sure we don&#39;t rely on them if they</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	<span class="comment">//           are not properly set yet.</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	for _, name := range names {
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		tname := NewTypeName(name.Pos(), check.pkg, name.Name, nil)
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		tpar := check.newTypeParam(tname, Typ[Invalid]) <span class="comment">// assigns type to tpar as a side-effect</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		check.declare(check.scope, name, tname, scopePos)
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		tparams = append(tparams, tpar)
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	}
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	if check.conf._Trace &amp;&amp; len(names) &gt; 0 {
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		check.trace(names[0].Pos(), &#34;type params = %v&#34;, tparams[len(tparams)-len(names):])
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	return tparams
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>}
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>func (check *Checker) collectMethods(obj *TypeName) {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	<span class="comment">// get associated methods</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	<span class="comment">// (Checker.collectObjects only collects methods with non-blank names;</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	<span class="comment">// Checker.resolveBaseTypeName ensures that obj is not an alias name</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	<span class="comment">// if it has attached methods.)</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	methods := check.methods[obj]
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	if methods == nil {
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>		return
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	}
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	delete(check.methods, obj)
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	assert(!check.objMap[obj].tdecl.Assign.IsValid()) <span class="comment">// don&#39;t use TypeName.IsAlias (requires fully set up object)</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	<span class="comment">// use an objset to check for name conflicts</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	var mset objset
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;If the base type is a struct type, the non-blank method</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	<span class="comment">// and field names must be distinct.&#34;</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	base := asNamed(obj.typ) <span class="comment">// shouldn&#39;t fail but be conservative</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	if base != nil {
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>		assert(base.TypeArgs().Len() == 0) <span class="comment">// collectMethods should not be called on an instantiated type</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		<span class="comment">// See go.dev/issue/52529: we must delay the expansion of underlying here, as</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>		<span class="comment">// base may not be fully set-up.</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		check.later(func() {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>			check.checkFieldUniqueness(base)
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		}).describef(obj, &#34;verifying field uniqueness for %v&#34;, base)
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		<span class="comment">// Checker.Files may be called multiple times; additional package files</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		<span class="comment">// may add methods to already type-checked types. Add pre-existing methods</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		<span class="comment">// so that we can detect redeclarations.</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		for i := 0; i &lt; base.NumMethods(); i++ {
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			m := base.Method(i)
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>			assert(m.name != &#34;_&#34;)
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>			assert(mset.insert(m) == nil)
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>		}
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	}
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	<span class="comment">// add valid methods</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	for _, m := range methods {
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;For a base type, the non-blank names of methods bound</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		<span class="comment">// to it must be unique.&#34;</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>		assert(m.name != &#34;_&#34;)
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		if alt := mset.insert(m); alt != nil {
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>			if alt.Pos().IsValid() {
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>				check.errorf(m, DuplicateMethod, &#34;method %s.%s already declared at %s&#34;, obj.Name(), m.name, alt.Pos())
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>			} else {
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>				check.errorf(m, DuplicateMethod, &#34;method %s.%s already declared&#34;, obj.Name(), m.name)
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>			}
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>			continue
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>		}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>		if base != nil {
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>			base.AddMethod(m)
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		}
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	}
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>func (check *Checker) checkFieldUniqueness(base *Named) {
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	if t, _ := base.under().(*Struct); t != nil {
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		var mset objset
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>		for i := 0; i &lt; base.NumMethods(); i++ {
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>			m := base.Method(i)
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>			assert(m.name != &#34;_&#34;)
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>			assert(mset.insert(m) == nil)
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>		}
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		<span class="comment">// Check that any non-blank field names of base are distinct from its</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>		<span class="comment">// method names.</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		for _, fld := range t.fields {
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>			if fld.name != &#34;_&#34; {
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>				if alt := mset.insert(fld); alt != nil {
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>					<span class="comment">// Struct fields should already be unique, so we should only</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>					<span class="comment">// encounter an alternate via collision with a method name.</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>					_ = alt.(*Func)
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>					<span class="comment">// For historical consistency, we report the primary error on the</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>					<span class="comment">// method, and the alt decl on the field.</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>					check.errorf(alt, DuplicateFieldAndMethod, &#34;field and method with the same name %s&#34;, fld.name)
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>					check.reportAltDecl(fld)
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>				}
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>			}
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		}
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	}
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>}
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>func (check *Checker) funcDecl(obj *Func, decl *declInfo) {
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	assert(obj.typ == nil)
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	<span class="comment">// func declarations cannot use iota</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	assert(check.iota == nil)
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	sig := new(Signature)
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	obj.typ = sig <span class="comment">// guard against cycles</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	<span class="comment">// Avoid cycle error when referring to method while type-checking the signature.</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>	<span class="comment">// This avoids a nuisance in the best case (non-parameterized receiver type) and</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	<span class="comment">// since the method is not a type, we get an error. If we have a parameterized</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	<span class="comment">// receiver type, instantiating the receiver type leads to the instantiation of</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	<span class="comment">// its methods, and we don&#39;t want a cycle error in that case.</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	<span class="comment">// TODO(gri) review if this is correct and/or whether we still need this?</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	saved := obj.color_
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	obj.color_ = black
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	fdecl := decl.fdecl
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>	check.funcType(sig, fdecl.Recv, fdecl.Type)
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	obj.color_ = saved
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	<span class="comment">// Set the scope&#39;s extent to the complete &#34;func (...) { ... }&#34;</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	<span class="comment">// so that Scope.Innermost works correctly.</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	sig.scope.pos = fdecl.Pos()
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	sig.scope.end = fdecl.End()
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>	if fdecl.Type.TypeParams.NumFields() &gt; 0 &amp;&amp; fdecl.Body == nil {
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>		check.softErrorf(fdecl.Name, BadDecl, &#34;generic function is missing function body&#34;)
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>	}
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	<span class="comment">// function body must be type-checked after global declarations</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	<span class="comment">// (functions implemented elsewhere have no body)</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	if !check.conf.IgnoreFuncBodies &amp;&amp; fdecl.Body != nil {
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>		check.later(func() {
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>			check.funcBody(decl, obj.name, sig, fdecl.Body, nil)
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>		}).describef(obj, &#34;func %s&#34;, obj.name)
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	}
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>}
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>func (check *Checker) declStmt(d ast.Decl) {
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	pkg := check.pkg
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	check.walkDecl(d, func(d decl) {
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		switch d := d.(type) {
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		case constDecl:
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>			top := len(check.delayed)
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>			<span class="comment">// declare all constants</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>			lhs := make([]*Const, len(d.spec.Names))
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>			for i, name := range d.spec.Names {
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>				obj := NewConst(name.Pos(), pkg, name.Name, nil, constant.MakeInt64(int64(d.iota)))
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>				lhs[i] = obj
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>				var init ast.Expr
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>				if i &lt; len(d.init) {
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>					init = d.init[i]
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>				}
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>				check.constDecl(obj, d.typ, init, d.inherited)
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>			}
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>			<span class="comment">// process function literals in init expressions before scope changes</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>			check.processDelayed(top)
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>			<span class="comment">// spec: &#34;The scope of a constant or variable identifier declared</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>			<span class="comment">// inside a function begins at the end of the ConstSpec or VarSpec</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>			<span class="comment">// (ShortVarDecl for short variable declarations) and ends at the</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>			<span class="comment">// end of the innermost containing block.&#34;</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>			scopePos := d.spec.End()
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>			for i, name := range d.spec.Names {
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>				check.declare(check.scope, name, lhs[i], scopePos)
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			}
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>		case varDecl:
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>			top := len(check.delayed)
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>			lhs0 := make([]*Var, len(d.spec.Names))
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>			for i, name := range d.spec.Names {
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>				lhs0[i] = NewVar(name.Pos(), pkg, name.Name, nil)
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>			}
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>			<span class="comment">// initialize all variables</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>			for i, obj := range lhs0 {
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>				var lhs []*Var
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>				var init ast.Expr
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>				switch len(d.spec.Values) {
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>				case len(d.spec.Names):
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>					<span class="comment">// lhs and rhs match</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>					init = d.spec.Values[i]
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>				case 1:
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>					<span class="comment">// rhs is expected to be a multi-valued expression</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>					lhs = lhs0
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>					init = d.spec.Values[0]
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>				default:
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>					if i &lt; len(d.spec.Values) {
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>						init = d.spec.Values[i]
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>					}
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>				}
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>				check.varDecl(obj, lhs, d.spec.Type, init)
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>				if len(d.spec.Values) == 1 {
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>					<span class="comment">// If we have a single lhs variable we are done either way.</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>					<span class="comment">// If we have a single rhs expression, it must be a multi-</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>					<span class="comment">// valued expression, in which case handling the first lhs</span>
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>					<span class="comment">// variable will cause all lhs variables to have a type</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>					<span class="comment">// assigned, and we are done as well.</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>					if debug {
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>						for _, obj := range lhs0 {
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>							assert(obj.typ != nil)
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>						}
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>					}
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>					break
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>				}
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>			}
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>			<span class="comment">// process function literals in init expressions before scope changes</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>			check.processDelayed(top)
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>			<span class="comment">// declare all variables</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>			<span class="comment">// (only at this point are the variable scopes (parents) set)</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>			scopePos := d.spec.End() <span class="comment">// see constant declarations</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>			for i, name := range d.spec.Names {
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>				<span class="comment">// see constant declarations</span>
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>				check.declare(check.scope, name, lhs0[i], scopePos)
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>			}
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>		case typeDecl:
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>			obj := NewTypeName(d.spec.Name.Pos(), pkg, d.spec.Name.Name, nil)
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>			<span class="comment">// spec: &#34;The scope of a type identifier declared inside a function</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>			<span class="comment">// begins at the identifier in the TypeSpec and ends at the end of</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>			<span class="comment">// the innermost containing block.&#34;</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>			scopePos := d.spec.Name.Pos()
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>			check.declare(check.scope, d.spec.Name, obj, scopePos)
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>			<span class="comment">// mark and unmark type before calling typeDecl; its type is still nil (see Checker.objDecl)</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>			obj.setColor(grey + color(check.push(obj)))
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			check.typeDecl(obj, d.spec, nil)
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>			check.pop().setColor(black)
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>		default:
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>			check.errorf(d.node(), InvalidSyntaxTree, &#34;unknown ast.Decl node %T&#34;, d.node())
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>		}
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	})
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>}
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>
</pre><p><a href="decl.go?m=text">View as plain text</a></p>

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
