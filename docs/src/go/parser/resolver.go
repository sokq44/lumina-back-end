<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/parser/resolver.go - Go Documentation Server</title>

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
<a href="resolver.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/parser">parser</a>/<span class="text-muted">resolver.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/parser">go/parser</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2021 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package parser
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>const debugResolve = false
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// resolveFile walks the given file to resolve identifiers within the file</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// scope, updating ast.Ident.Obj fields with declaration information.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// If declErr is non-nil, it is used to report declaration errors during</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// resolution. tok is used to format position in error messages.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>func resolveFile(file *ast.File, handle *token.File, declErr func(token.Pos, string)) {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	pkgScope := ast.NewScope(nil)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	r := &amp;resolver{
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>		handle:   handle,
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		declErr:  declErr,
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		topScope: pkgScope,
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		pkgScope: pkgScope,
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		depth:    1,
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	for _, decl := range file.Decls {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		ast.Walk(r, decl)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	r.closeScope()
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	assert(r.topScope == nil, &#34;unbalanced scopes&#34;)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	assert(r.labelScope == nil, &#34;unbalanced label scopes&#34;)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// resolve global identifiers within the same file</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	i := 0
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	for _, ident := range r.unresolved {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		<span class="comment">// i &lt;= index for current ident</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		assert(ident.Obj == unresolved, &#34;object already resolved&#34;)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		ident.Obj = r.pkgScope.Lookup(ident.Name) <span class="comment">// also removes unresolved sentinel</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		if ident.Obj == nil {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			r.unresolved[i] = ident
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			i++
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		} else if debugResolve {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			pos := ident.Obj.Decl.(interface{ Pos() token.Pos }).Pos()
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			r.trace(&#34;resolved %s@%v to package object %v&#34;, ident.Name, ident.Pos(), pos)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	file.Scope = r.pkgScope
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	file.Unresolved = r.unresolved[0:i]
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>const maxScopeDepth int = 1e3
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>type resolver struct {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	handle  *token.File
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	declErr func(token.Pos, string)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// Ordinary identifier scopes</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	pkgScope   *ast.Scope   <span class="comment">// pkgScope.Outer == nil</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	topScope   *ast.Scope   <span class="comment">// top-most scope; may be pkgScope</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	unresolved []*ast.Ident <span class="comment">// unresolved identifiers</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	depth      int          <span class="comment">// scope depth</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// Label scopes</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// (maintained by open/close LabelScope)</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	labelScope  *ast.Scope     <span class="comment">// label scope for current function</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	targetStack [][]*ast.Ident <span class="comment">// stack of unresolved labels</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>func (r *resolver) trace(format string, args ...any) {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	fmt.Println(strings.Repeat(&#34;. &#34;, r.depth) + r.sprintf(format, args...))
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func (r *resolver) sprintf(format string, args ...any) string {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	for i, arg := range args {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		switch arg := arg.(type) {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		case token.Pos:
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			args[i] = r.handle.Position(arg)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	return fmt.Sprintf(format, args...)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>func (r *resolver) openScope(pos token.Pos) {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	r.depth++
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	if r.depth &gt; maxScopeDepth {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		panic(bailout{pos: pos, msg: &#34;exceeded max scope depth during object resolution&#34;})
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	if debugResolve {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		r.trace(&#34;opening scope @%v&#34;, pos)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	r.topScope = ast.NewScope(r.topScope)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>func (r *resolver) closeScope() {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	r.depth--
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	if debugResolve {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		r.trace(&#34;closing scope&#34;)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	r.topScope = r.topScope.Outer
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>func (r *resolver) openLabelScope() {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	r.labelScope = ast.NewScope(r.labelScope)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	r.targetStack = append(r.targetStack, nil)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func (r *resolver) closeLabelScope() {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// resolve labels</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	n := len(r.targetStack) - 1
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	scope := r.labelScope
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	for _, ident := range r.targetStack[n] {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		ident.Obj = scope.Lookup(ident.Name)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		if ident.Obj == nil &amp;&amp; r.declErr != nil {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			r.declErr(ident.Pos(), fmt.Sprintf(&#34;label %s undefined&#34;, ident.Name))
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// pop label scope</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	r.targetStack = r.targetStack[0:n]
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	r.labelScope = r.labelScope.Outer
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>func (r *resolver) declare(decl, data any, scope *ast.Scope, kind ast.ObjKind, idents ...*ast.Ident) {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	for _, ident := range idents {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		if ident.Obj != nil {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			panic(fmt.Sprintf(&#34;%v: identifier %s already declared or resolved&#34;, ident.Pos(), ident.Name))
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		obj := ast.NewObj(kind, ident.Name)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// remember the corresponding declaration for redeclaration</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// errors and global variable resolution/typechecking phase</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		obj.Decl = decl
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		obj.Data = data
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		<span class="comment">// Identifiers (for receiver type parameters) are written to the scope, but</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		<span class="comment">// never set as the resolved object. See go.dev/issue/50956.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		if _, ok := decl.(*ast.Ident); !ok {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			ident.Obj = obj
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		if ident.Name != &#34;_&#34; {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			if debugResolve {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>				r.trace(&#34;declaring %s@%v&#34;, ident.Name, ident.Pos())
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			if alt := scope.Insert(obj); alt != nil &amp;&amp; r.declErr != nil {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				prevDecl := &#34;&#34;
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>				if pos := alt.Pos(); pos.IsValid() {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>					prevDecl = r.sprintf(&#34;\n\tprevious declaration at %v&#34;, pos)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>				}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>				r.declErr(ident.Pos(), fmt.Sprintf(&#34;%s redeclared in this block%s&#34;, ident.Name, prevDecl))
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>func (r *resolver) shortVarDecl(decl *ast.AssignStmt) {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// Go spec: A short variable declaration may redeclare variables</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// provided they were originally declared in the same block with</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// the same type, and at least one of the non-blank variables is new.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	n := 0 <span class="comment">// number of new variables</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	for _, x := range decl.Lhs {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		if ident, isIdent := x.(*ast.Ident); isIdent {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			assert(ident.Obj == nil, &#34;identifier already declared or resolved&#34;)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			obj := ast.NewObj(ast.Var, ident.Name)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>			<span class="comment">// remember corresponding assignment for other tools</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			obj.Decl = decl
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			ident.Obj = obj
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			if ident.Name != &#34;_&#34; {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>				if debugResolve {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>					r.trace(&#34;declaring %s@%v&#34;, ident.Name, ident.Pos())
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>				}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>				if alt := r.topScope.Insert(obj); alt != nil {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>					ident.Obj = alt <span class="comment">// redeclaration</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>				} else {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>					n++ <span class="comment">// new declaration</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>				}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	if n == 0 &amp;&amp; r.declErr != nil {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		r.declErr(decl.Lhs[0].Pos(), &#34;no new variables on left side of :=&#34;)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// The unresolved object is a sentinel to mark identifiers that have been added</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// to the list of unresolved identifiers. The sentinel is only used for verifying</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// internal consistency.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>var unresolved = new(ast.Object)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// If x is an identifier, resolve attempts to resolve x by looking up</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">// the object it denotes. If no object is found and collectUnresolved is</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span><span class="comment">// set, x is marked as unresolved and collected in the list of unresolved</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span><span class="comment">// identifiers.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>func (r *resolver) resolve(ident *ast.Ident, collectUnresolved bool) {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	if ident.Obj != nil {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		panic(r.sprintf(&#34;%v: identifier %s already declared or resolved&#34;, ident.Pos(), ident.Name))
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// &#39;_&#39; should never refer to existing declarations, because it has special</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// handling in the spec.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	if ident.Name == &#34;_&#34; {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		return
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	for s := r.topScope; s != nil; s = s.Outer {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		if obj := s.Lookup(ident.Name); obj != nil {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			if debugResolve {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>				r.trace(&#34;resolved %v:%s to %v&#34;, ident.Pos(), ident.Name, obj)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			assert(obj.Name != &#34;&#34;, &#34;obj with no name&#34;)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			<span class="comment">// Identifiers (for receiver type parameters) are written to the scope,</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			<span class="comment">// but never set as the resolved object. See go.dev/issue/50956.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			if _, ok := obj.Decl.(*ast.Ident); !ok {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>				ident.Obj = obj
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			return
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	<span class="comment">// all local scopes are known, so any unresolved identifier</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// must be found either in the file scope, package scope</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// (perhaps in another file), or universe scope --- collect</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">// them so that they can be resolved later</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	if collectUnresolved {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		ident.Obj = unresolved
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		r.unresolved = append(r.unresolved, ident)
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>func (r *resolver) walkExprs(list []ast.Expr) {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	for _, node := range list {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		ast.Walk(r, node)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>func (r *resolver) walkLHS(list []ast.Expr) {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	for _, expr := range list {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		expr := ast.Unparen(expr)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		if _, ok := expr.(*ast.Ident); !ok &amp;&amp; expr != nil {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			ast.Walk(r, expr)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>func (r *resolver) walkStmts(list []ast.Stmt) {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	for _, stmt := range list {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		ast.Walk(r, stmt)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>func (r *resolver) Visit(node ast.Node) ast.Visitor {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	if debugResolve &amp;&amp; node != nil {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		r.trace(&#34;node %T@%v&#34;, node, node.Pos())
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	switch n := node.(type) {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">// Expressions.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	case *ast.Ident:
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		r.resolve(n, true)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	case *ast.FuncLit:
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		r.walkFuncType(n.Type)
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		r.walkBody(n.Body)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	case *ast.SelectorExpr:
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		ast.Walk(r, n.X)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		<span class="comment">// Note: don&#39;t try to resolve n.Sel, as we don&#39;t support qualified</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		<span class="comment">// resolution.</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	case *ast.StructType:
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		r.walkFieldList(n.Fields, ast.Var)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	case *ast.FuncType:
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		r.walkFuncType(n)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	case *ast.CompositeLit:
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		if n.Type != nil {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			ast.Walk(r, n.Type)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		for _, e := range n.Elts {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			if kv, _ := e.(*ast.KeyValueExpr); kv != nil {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>				<span class="comment">// See go.dev/issue/45160: try to resolve composite lit keys, but don&#39;t</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>				<span class="comment">// collect them as unresolved if resolution failed. This replicates</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>				<span class="comment">// existing behavior when resolving during parsing.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>				if ident, _ := kv.Key.(*ast.Ident); ident != nil {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>					r.resolve(ident, false)
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>				} else {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>					ast.Walk(r, kv.Key)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>				}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>				ast.Walk(r, kv.Value)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			} else {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>				ast.Walk(r, e)
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	case *ast.InterfaceType:
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		r.walkFieldList(n.Methods, ast.Fun)
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	<span class="comment">// Statements</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	case *ast.LabeledStmt:
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		r.declare(n, nil, r.labelScope, ast.Lbl, n.Label)
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		ast.Walk(r, n.Stmt)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	case *ast.AssignStmt:
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		r.walkExprs(n.Rhs)
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		if n.Tok == token.DEFINE {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>			r.shortVarDecl(n)
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		} else {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			r.walkExprs(n.Lhs)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	case *ast.BranchStmt:
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		<span class="comment">// add to list of unresolved targets</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		if n.Tok != token.FALLTHROUGH &amp;&amp; n.Label != nil {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>			depth := len(r.targetStack) - 1
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			r.targetStack[depth] = append(r.targetStack[depth], n.Label)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	case *ast.BlockStmt:
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		r.walkStmts(n.List)
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	case *ast.IfStmt:
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		if n.Init != nil {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			ast.Walk(r, n.Init)
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		ast.Walk(r, n.Cond)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		ast.Walk(r, n.Body)
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		if n.Else != nil {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			ast.Walk(r, n.Else)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	case *ast.CaseClause:
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		r.walkExprs(n.List)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		r.walkStmts(n.Body)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	case *ast.SwitchStmt:
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		if n.Init != nil {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>			ast.Walk(r, n.Init)
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		}
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		if n.Tag != nil {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			<span class="comment">// The scope below reproduces some unnecessary behavior of the parser,</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			<span class="comment">// opening an extra scope in case this is a type switch. It&#39;s not needed</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			<span class="comment">// for expression switches.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			<span class="comment">// TODO: remove this once we&#39;ve matched the parser resolution exactly.</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>			if n.Init != nil {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>				r.openScope(n.Tag.Pos())
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>				defer r.closeScope()
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			}
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			ast.Walk(r, n.Tag)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		if n.Body != nil {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			r.walkStmts(n.Body.List)
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		}
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	case *ast.TypeSwitchStmt:
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		if n.Init != nil {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			r.openScope(n.Pos())
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			defer r.closeScope()
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			ast.Walk(r, n.Init)
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		r.openScope(n.Assign.Pos())
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		ast.Walk(r, n.Assign)
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		<span class="comment">// s.Body consists only of case clauses, so does not get its own</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		<span class="comment">// scope.</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		if n.Body != nil {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			r.walkStmts(n.Body.List)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	case *ast.CommClause:
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		if n.Comm != nil {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			ast.Walk(r, n.Comm)
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		}
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		r.walkStmts(n.Body)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	case *ast.SelectStmt:
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		<span class="comment">// as for switch statements, select statement bodies don&#39;t get their own</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		<span class="comment">// scope.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		if n.Body != nil {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			r.walkStmts(n.Body.List)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	case *ast.ForStmt:
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		if n.Init != nil {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			ast.Walk(r, n.Init)
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		if n.Cond != nil {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			ast.Walk(r, n.Cond)
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		if n.Post != nil {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>			ast.Walk(r, n.Post)
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		ast.Walk(r, n.Body)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	case *ast.RangeStmt:
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		ast.Walk(r, n.X)
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		var lhs []ast.Expr
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		if n.Key != nil {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>			lhs = append(lhs, n.Key)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		if n.Value != nil {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			lhs = append(lhs, n.Value)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		if len(lhs) &gt; 0 {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>			if n.Tok == token.DEFINE {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>				<span class="comment">// Note: we can&#39;t exactly match the behavior of object resolution</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>				<span class="comment">// during the parsing pass here, as it uses the position of the RANGE</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>				<span class="comment">// token for the RHS OpPos. That information is not contained within</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>				<span class="comment">// the AST.</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>				as := &amp;ast.AssignStmt{
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>					Lhs:    lhs,
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>					Tok:    token.DEFINE,
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>					TokPos: n.TokPos,
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>					Rhs:    []ast.Expr{&amp;ast.UnaryExpr{Op: token.RANGE, X: n.X}},
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>				}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>				<span class="comment">// TODO(rFindley): this walkLHS reproduced the parser resolution, but</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>				<span class="comment">// is it necessary? By comparison, for a normal AssignStmt we don&#39;t</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>				<span class="comment">// walk the LHS in case there is an invalid identifier list.</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>				r.walkLHS(lhs)
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>				r.shortVarDecl(as)
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>			} else {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>				r.walkExprs(lhs)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		ast.Walk(r, n.Body)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	<span class="comment">// Declarations</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	case *ast.GenDecl:
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		switch n.Tok {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		case token.CONST, token.VAR:
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			for i, spec := range n.Specs {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>				spec := spec.(*ast.ValueSpec)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>				kind := ast.Con
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>				if n.Tok == token.VAR {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>					kind = ast.Var
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>				}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>				r.walkExprs(spec.Values)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>				if spec.Type != nil {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>					ast.Walk(r, spec.Type)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>				}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>				r.declare(spec, i, r.topScope, kind, spec.Names...)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			}
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		case token.TYPE:
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			for _, spec := range n.Specs {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>				spec := spec.(*ast.TypeSpec)
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>				<span class="comment">// Go spec: The scope of a type identifier declared inside a function begins</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>				<span class="comment">// at the identifier in the TypeSpec and ends at the end of the innermost</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>				<span class="comment">// containing block.</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>				r.declare(spec, nil, r.topScope, ast.Typ, spec.Name)
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>				if spec.TypeParams != nil {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>					r.openScope(spec.Pos())
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>					defer r.closeScope()
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>					r.walkTParams(spec.TypeParams)
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>				}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>				ast.Walk(r, spec.Type)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		}
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	case *ast.FuncDecl:
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		<span class="comment">// Open the function scope.</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		r.openScope(n.Pos())
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		defer r.closeScope()
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		r.walkRecv(n.Recv)
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		<span class="comment">// Type parameters are walked normally: they can reference each other, and</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		<span class="comment">// can be referenced by normal parameters.</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		if n.Type.TypeParams != nil {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>			r.walkTParams(n.Type.TypeParams)
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>			<span class="comment">// TODO(rFindley): need to address receiver type parameters.</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		<span class="comment">// Resolve and declare parameters in a specific order to get duplicate</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		<span class="comment">// declaration errors in the correct location.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		r.resolveList(n.Type.Params)
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		r.resolveList(n.Type.Results)
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		r.declareList(n.Recv, ast.Var)
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		r.declareList(n.Type.Params, ast.Var)
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		r.declareList(n.Type.Results, ast.Var)
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		r.walkBody(n.Body)
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		if n.Recv == nil &amp;&amp; n.Name.Name != &#34;init&#34; {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>			r.declare(n, nil, r.pkgScope, ast.Fun, n.Name)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	default:
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		return r
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	return nil
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>func (r *resolver) walkFuncType(typ *ast.FuncType) {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	<span class="comment">// typ.TypeParams must be walked separately for FuncDecls.</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	r.resolveList(typ.Params)
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	r.resolveList(typ.Results)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	r.declareList(typ.Params, ast.Var)
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	r.declareList(typ.Results, ast.Var)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>func (r *resolver) resolveList(list *ast.FieldList) {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	if list == nil {
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		return
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	for _, f := range list.List {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		if f.Type != nil {
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			ast.Walk(r, f.Type)
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		}
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	}
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>func (r *resolver) declareList(list *ast.FieldList, kind ast.ObjKind) {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	if list == nil {
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		return
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	for _, f := range list.List {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		r.declare(f, nil, r.topScope, kind, f.Names...)
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	}
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>}
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>func (r *resolver) walkRecv(recv *ast.FieldList) {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	<span class="comment">// If our receiver has receiver type parameters, we must declare them before</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	<span class="comment">// trying to resolve the rest of the receiver, and avoid re-resolving the</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	<span class="comment">// type parameter identifiers.</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	if recv == nil || len(recv.List) == 0 {
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		return <span class="comment">// nothing to do</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	}
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	typ := recv.List[0].Type
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	if ptr, ok := typ.(*ast.StarExpr); ok {
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		typ = ptr.X
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	}
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	var declareExprs []ast.Expr <span class="comment">// exprs to declare</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	var resolveExprs []ast.Expr <span class="comment">// exprs to resolve</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	switch typ := typ.(type) {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	case *ast.IndexExpr:
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		declareExprs = []ast.Expr{typ.Index}
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		resolveExprs = append(resolveExprs, typ.X)
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	case *ast.IndexListExpr:
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		declareExprs = typ.Indices
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		resolveExprs = append(resolveExprs, typ.X)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	default:
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		resolveExprs = append(resolveExprs, typ)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	for _, expr := range declareExprs {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		if id, _ := expr.(*ast.Ident); id != nil {
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>			r.declare(expr, nil, r.topScope, ast.Typ, id)
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		} else {
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>			<span class="comment">// The receiver type parameter expression is invalid, but try to resolve</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>			<span class="comment">// it anyway for consistency.</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>			resolveExprs = append(resolveExprs, expr)
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	}
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	for _, expr := range resolveExprs {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		if expr != nil {
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>			ast.Walk(r, expr)
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	<span class="comment">// The receiver is invalid, but try to resolve it anyway for consistency.</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	for _, f := range recv.List[1:] {
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		if f.Type != nil {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>			ast.Walk(r, f.Type)
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	}
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>func (r *resolver) walkFieldList(list *ast.FieldList, kind ast.ObjKind) {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	if list == nil {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		return
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	}
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	r.resolveList(list)
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	r.declareList(list, kind)
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>}
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span><span class="comment">// walkTParams is like walkFieldList, but declares type parameters eagerly so</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span><span class="comment">// that they may be resolved in the constraint expressions held in the field</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span><span class="comment">// Type.</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>func (r *resolver) walkTParams(list *ast.FieldList) {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	r.declareList(list, ast.Typ)
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	r.resolveList(list)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>}
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>func (r *resolver) walkBody(body *ast.BlockStmt) {
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	if body == nil {
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		return
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	r.openLabelScope()
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	defer r.closeLabelScope()
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	r.walkStmts(body.List)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
</pre><p><a href="resolver.go?m=text">View as plain text</a></p>

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
