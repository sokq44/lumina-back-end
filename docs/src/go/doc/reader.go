<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/doc/reader.go - Go Documentation Server</title>

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
<a href="reader.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/doc">doc</a>/<span class="text-muted">reader.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/doc">go/doc</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package doc
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/lazyregexp&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;path&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// function/method sets</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// Internally, we treat functions like methods and collect them in method sets.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// A methodSet describes a set of methods. Entries where Decl == nil are conflict</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// entries (more than one method with the same name at the same embedding level).</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>type methodSet map[string]*Func
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// recvString returns a string representation of recv of the form &#34;T&#34;, &#34;*T&#34;,</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// &#34;T[A, ...]&#34;, &#34;*T[A, ...]&#34; or &#34;BADRECV&#34; (if not a proper receiver type).</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>func recvString(recv ast.Expr) string {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	switch t := recv.(type) {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	case *ast.Ident:
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		return t.Name
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	case *ast.StarExpr:
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		return &#34;*&#34; + recvString(t.X)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	case *ast.IndexExpr:
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		<span class="comment">// Generic type with one parameter.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		return fmt.Sprintf(&#34;%s[%s]&#34;, recvString(t.X), recvParam(t.Index))
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	case *ast.IndexListExpr:
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		<span class="comment">// Generic type with multiple parameters.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		if len(t.Indices) &gt; 0 {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>			var b strings.Builder
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>			b.WriteString(recvString(t.X))
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>			b.WriteByte(&#39;[&#39;)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			b.WriteString(recvParam(t.Indices[0]))
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			for _, e := range t.Indices[1:] {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>				b.WriteString(&#34;, &#34;)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>				b.WriteString(recvParam(e))
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			b.WriteByte(&#39;]&#39;)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			return b.String()
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	return &#34;BADRECV&#34;
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func recvParam(p ast.Expr) string {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if id, ok := p.(*ast.Ident); ok {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		return id.Name
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	return &#34;BADPARAM&#34;
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// set creates the corresponding Func for f and adds it to mset.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// If there are multiple f&#39;s with the same name, set keeps the first</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// one with documentation; conflicts are ignored. The boolean</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// specifies whether to leave the AST untouched.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>func (mset methodSet) set(f *ast.FuncDecl, preserveAST bool) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	name := f.Name.Name
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if g := mset[name]; g != nil &amp;&amp; g.Doc != &#34;&#34; {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		<span class="comment">// A function with the same name has already been registered;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		<span class="comment">// since it has documentation, assume f is simply another</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		<span class="comment">// implementation and ignore it. This does not happen if the</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		<span class="comment">// caller is using go/build.ScanDir to determine the list of</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		<span class="comment">// files implementing a package.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		return
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// function doesn&#39;t exist or has no documentation; use f</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	recv := &#34;&#34;
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	if f.Recv != nil {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		var typ ast.Expr
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		<span class="comment">// be careful in case of incorrect ASTs</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		if list := f.Recv.List; len(list) == 1 {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			typ = list[0].Type
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		recv = recvString(typ)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	mset[name] = &amp;Func{
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		Doc:  f.Doc.Text(),
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		Name: name,
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		Decl: f,
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		Recv: recv,
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		Orig: recv,
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if !preserveAST {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		f.Doc = nil <span class="comment">// doc consumed - remove from AST</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// add adds method m to the method set; m is ignored if the method set</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// already contains a method with the same name at the same or a higher</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// level than m.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>func (mset methodSet) add(m *Func) {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	old := mset[m.Name]
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	if old == nil || m.Level &lt; old.Level {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		mset[m.Name] = m
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		return
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	if m.Level == old.Level {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		<span class="comment">// conflict - mark it using a method with nil Decl</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		mset[m.Name] = &amp;Func{
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			Name:  m.Name,
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			Level: m.Level,
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// Named types</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// baseTypeName returns the name of the base type of x (or &#34;&#34;)</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// and whether the type is imported or not.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func baseTypeName(x ast.Expr) (name string, imported bool) {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	switch t := x.(type) {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	case *ast.Ident:
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		return t.Name, false
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	case *ast.IndexExpr:
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		return baseTypeName(t.X)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	case *ast.IndexListExpr:
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		return baseTypeName(t.X)
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	case *ast.SelectorExpr:
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		if _, ok := t.X.(*ast.Ident); ok {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			<span class="comment">// only possible for qualified type names;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			<span class="comment">// assume type is imported</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			return t.Sel.Name, true
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	case *ast.ParenExpr:
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		return baseTypeName(t.X)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	case *ast.StarExpr:
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		return baseTypeName(t.X)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	return &#34;&#34;, false
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// An embeddedSet describes a set of embedded types.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>type embeddedSet map[*namedType]bool
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// A namedType represents a named unqualified (package local, or possibly</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// predeclared) type. The namedType for a type name is always found via</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// reader.lookupType.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>type namedType struct {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	doc  string       <span class="comment">// doc comment for type</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	name string       <span class="comment">// type name</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	decl *ast.GenDecl <span class="comment">// nil if declaration hasn&#39;t been seen yet</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	isEmbedded bool        <span class="comment">// true if this type is embedded</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	isStruct   bool        <span class="comment">// true if this type is a struct</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	embedded   embeddedSet <span class="comment">// true if the embedded type is a pointer</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// associated declarations</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	values  []*Value <span class="comment">// consts and vars</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	funcs   methodSet
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	methods methodSet
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span><span class="comment">// AST reader</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">// reader accumulates documentation for a single package.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// It modifies the AST: Comments (declaration documentation)</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// that have been collected by the reader are set to nil</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// in the respective AST nodes so that they are not printed</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">// twice (once when printing the documentation and once when</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// printing the corresponding AST node).</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>type reader struct {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	mode Mode
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// package properties</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	doc       string <span class="comment">// package documentation, if any</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	filenames []string
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	notes     map[string][]*Note
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// imports</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	imports      map[string]int
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	hasDotImp    bool <span class="comment">// if set, package contains a dot import</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	importByName map[string]string
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// declarations</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	values []*Value <span class="comment">// consts and vars</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	order  int      <span class="comment">// sort order of const and var declarations (when we can&#39;t use a name)</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	types  map[string]*namedType
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	funcs  methodSet
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// support for package-local shadowing of predeclared types</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	shadowedPredecl map[string]bool
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	fixmap          map[string][]*ast.InterfaceType
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>func (r *reader) isVisible(name string) bool {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	return r.mode&amp;AllDecls != 0 || token.IsExported(name)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">// lookupType returns the base type with the given name.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// If the base type has not been encountered yet, a new</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// type with the given name but no associated declaration</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// is added to the type map.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>func (r *reader) lookupType(name string) *namedType {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	if name == &#34;&#34; || name == &#34;_&#34; {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		return nil <span class="comment">// no type docs for anonymous types</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	if typ, found := r.types[name]; found {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		return typ
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// type not found - add one without declaration</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	typ := &amp;namedType{
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		name:     name,
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		embedded: make(embeddedSet),
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		funcs:    make(methodSet),
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		methods:  make(methodSet),
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	r.types[name] = typ
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	return typ
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// recordAnonymousField registers fieldType as the type of an</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// anonymous field in the parent type. If the field is imported</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// (qualified name) or the parent is nil, the field is ignored.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// The function returns the field name.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>func (r *reader) recordAnonymousField(parent *namedType, fieldType ast.Expr) (fname string) {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	fname, imp := baseTypeName(fieldType)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	if parent == nil || imp {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		return
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	if ftype := r.lookupType(fname); ftype != nil {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		ftype.isEmbedded = true
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		_, ptr := fieldType.(*ast.StarExpr)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		parent.embedded[ftype] = ptr
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	return
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>func (r *reader) readDoc(comment *ast.CommentGroup) {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// By convention there should be only one package comment</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// but collect all of them if there are more than one.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	text := comment.Text()
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	if r.doc == &#34;&#34; {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		r.doc = text
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		return
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	r.doc += &#34;\n&#34; + text
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>func (r *reader) remember(predecl string, typ *ast.InterfaceType) {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	if r.fixmap == nil {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		r.fixmap = make(map[string][]*ast.InterfaceType)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	r.fixmap[predecl] = append(r.fixmap[predecl], typ)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>func specNames(specs []ast.Spec) []string {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	names := make([]string, 0, len(specs)) <span class="comment">// reasonable estimate</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	for _, s := range specs {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		<span class="comment">// s guaranteed to be an *ast.ValueSpec by readValue</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		for _, ident := range s.(*ast.ValueSpec).Names {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			names = append(names, ident.Name)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	}
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	return names
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span><span class="comment">// readValue processes a const or var declaration.</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>func (r *reader) readValue(decl *ast.GenDecl) {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	<span class="comment">// determine if decl should be associated with a type</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	<span class="comment">// Heuristic: For each typed entry, determine the type name, if any.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">//            If there is exactly one type name that is sufficiently</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	<span class="comment">//            frequent, associate the decl with the respective type.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	domName := &#34;&#34;
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	domFreq := 0
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	prev := &#34;&#34;
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	n := 0
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	for _, spec := range decl.Specs {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		s, ok := spec.(*ast.ValueSpec)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		if !ok {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			continue <span class="comment">// should not happen, but be conservative</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		name := &#34;&#34;
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		switch {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		case s.Type != nil:
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			<span class="comment">// a type is present; determine its name</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			if n, imp := baseTypeName(s.Type); !imp {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>				name = n
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		case decl.Tok == token.CONST &amp;&amp; len(s.Values) == 0:
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			<span class="comment">// no type or value is present but we have a constant declaration;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			<span class="comment">// use the previous type name (possibly the empty string)</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			name = prev
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		if name != &#34;&#34; {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			<span class="comment">// entry has a named type</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			if domName != &#34;&#34; &amp;&amp; domName != name {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>				<span class="comment">// more than one type name - do not associate</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>				<span class="comment">// with any type</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>				domName = &#34;&#34;
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>				break
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			domName = name
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>			domFreq++
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		prev = name
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		n++
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	<span class="comment">// nothing to do w/o a legal declaration</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	if n == 0 {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		return
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	<span class="comment">// determine values list with which to associate the Value for this decl</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	values := &amp;r.values
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	const threshold = 0.75
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	if domName != &#34;&#34; &amp;&amp; r.isVisible(domName) &amp;&amp; domFreq &gt;= int(float64(len(decl.Specs))*threshold) {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		<span class="comment">// typed entries are sufficiently frequent</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		if typ := r.lookupType(domName); typ != nil {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			values = &amp;typ.values <span class="comment">// associate with that type</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	*values = append(*values, &amp;Value{
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		Doc:   decl.Doc.Text(),
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		Names: specNames(decl.Specs),
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		Decl:  decl,
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		order: r.order,
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	})
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	if r.mode&amp;PreserveAST == 0 {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		decl.Doc = nil <span class="comment">// doc consumed - remove from AST</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// Note: It&#39;s important that the order used here is global because the cleanupTypes</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// methods may move values associated with types back into the global list. If the</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">// order is list-specific, sorting is not deterministic because the same order value</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	<span class="comment">// may appear multiple times (was bug, found when fixing #16153).</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	r.order++
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// fields returns a struct&#39;s fields or an interface&#39;s methods.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>func fields(typ ast.Expr) (list []*ast.Field, isStruct bool) {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	var fields *ast.FieldList
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	switch t := typ.(type) {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	case *ast.StructType:
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		fields = t.Fields
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		isStruct = true
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	case *ast.InterfaceType:
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		fields = t.Methods
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	if fields != nil {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		list = fields.List
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	return
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// readType processes a type declaration.</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>func (r *reader) readType(decl *ast.GenDecl, spec *ast.TypeSpec) {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	typ := r.lookupType(spec.Name.Name)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	if typ == nil {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		return <span class="comment">// no name or blank name - ignore the type</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	<span class="comment">// A type should be added at most once, so typ.decl</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// should be nil - if it is not, simply overwrite it.</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	typ.decl = decl
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	<span class="comment">// compute documentation</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	doc := spec.Doc
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	if doc == nil {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		<span class="comment">// no doc associated with the spec, use the declaration doc, if any</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		doc = decl.Doc
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	if r.mode&amp;PreserveAST == 0 {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		spec.Doc = nil <span class="comment">// doc consumed - remove from AST</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		decl.Doc = nil <span class="comment">// doc consumed - remove from AST</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	typ.doc = doc.Text()
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// record anonymous fields (they may contribute methods)</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	<span class="comment">// (some fields may have been recorded already when filtering</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	<span class="comment">// exports, but that&#39;s ok)</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	var list []*ast.Field
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	list, typ.isStruct = fields(spec.Type)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	for _, field := range list {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		if len(field.Names) == 0 {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			r.recordAnonymousField(typ, field.Type)
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">// isPredeclared reports whether n denotes a predeclared type.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>func (r *reader) isPredeclared(n string) bool {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	return predeclaredTypes[n] &amp;&amp; r.types[n] == nil
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span><span class="comment">// readFunc processes a func or method declaration.</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>func (r *reader) readFunc(fun *ast.FuncDecl) {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	<span class="comment">// strip function body if requested.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	if r.mode&amp;PreserveAST == 0 {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		fun.Body = nil
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	<span class="comment">// associate methods with the receiver type, if any</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	if fun.Recv != nil {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		<span class="comment">// method</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		if len(fun.Recv.List) == 0 {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			<span class="comment">// should not happen (incorrect AST); (See issue 17788)</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			<span class="comment">// don&#39;t show this method</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			return
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		recvTypeName, imp := baseTypeName(fun.Recv.List[0].Type)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		if imp {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			<span class="comment">// should not happen (incorrect AST);</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			<span class="comment">// don&#39;t show this method</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			return
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		if typ := r.lookupType(recvTypeName); typ != nil {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			typ.methods.set(fun, r.mode&amp;PreserveAST != 0)
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		<span class="comment">// otherwise ignore the method</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri): There may be exported methods of non-exported types</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		<span class="comment">// that can be called because of exported values (consts, vars, or</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		<span class="comment">// function results) of that type. Could determine if that is the</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		<span class="comment">// case and then show those methods in an appropriate section.</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		return
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	<span class="comment">// Associate factory functions with the first visible result type, as long as</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	<span class="comment">// others are predeclared types.</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	if fun.Type.Results.NumFields() &gt;= 1 {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		var typ *namedType <span class="comment">// type to associate the function with</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		numResultTypes := 0
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		for _, res := range fun.Type.Results.List {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>			factoryType := res.Type
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			if t, ok := factoryType.(*ast.ArrayType); ok {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>				<span class="comment">// We consider functions that return slices or arrays of type</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>				<span class="comment">// T (or pointers to T) as factory functions of T.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>				factoryType = t.Elt
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			if n, imp := baseTypeName(factoryType); !imp &amp;&amp; r.isVisible(n) &amp;&amp; !r.isPredeclared(n) {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>				if lookupTypeParam(n, fun.Type.TypeParams) != nil {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>					<span class="comment">// Issue #49477: don&#39;t associate fun with its type parameter result.</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>					<span class="comment">// A type parameter is not a defined type.</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>					continue
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>				}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>				if t := r.lookupType(n); t != nil {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>					typ = t
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>					numResultTypes++
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>					if numResultTypes &gt; 1 {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>						break
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>					}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>				}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		<span class="comment">// If there is exactly one result type,</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		<span class="comment">// associate the function with that type.</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		if numResultTypes == 1 {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			typ.funcs.set(fun, r.mode&amp;PreserveAST != 0)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			return
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	<span class="comment">// just an ordinary function</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	r.funcs.set(fun, r.mode&amp;PreserveAST != 0)
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span><span class="comment">// lookupTypeParam searches for type parameters named name within the tparams</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span><span class="comment">// field list, returning the relevant identifier if found, or nil if not.</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>func lookupTypeParam(name string, tparams *ast.FieldList) *ast.Ident {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	if tparams == nil {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		return nil
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	for _, field := range tparams.List {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		for _, id := range field.Names {
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			if id.Name == name {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>				return id
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		}
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	return nil
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>}
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>var (
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	noteMarker    = `([A-Z][A-Z]+)\(([^)]+)\):?`                <span class="comment">// MARKER(uid), MARKER at least 2 chars, uid at least 1 char</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	noteMarkerRx  = lazyregexp.New(`^[ \t]*` + noteMarker)      <span class="comment">// MARKER(uid) at text start</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	noteCommentRx = lazyregexp.New(`^/[/*][ \t]*` + noteMarker) <span class="comment">// MARKER(uid) at comment start</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>)
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span><span class="comment">// clean replaces each sequence of space, \r, or \t characters</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span><span class="comment">// with a single space and removes any trailing and leading spaces.</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>func clean(s string) string {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	var b []byte
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	p := byte(&#39; &#39;)
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i++ {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		q := s[i]
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		if q == &#39;\r&#39; || q == &#39;\t&#39; {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			q = &#39; &#39;
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		if q != &#39; &#39; || p != &#39; &#39; {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>			b = append(b, q)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			p = q
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	<span class="comment">// remove trailing blank, if any</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	if n := len(b); n &gt; 0 &amp;&amp; p == &#39; &#39; {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		b = b[0 : n-1]
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	return string(b)
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span><span class="comment">// readNote collects a single note from a sequence of comments.</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>func (r *reader) readNote(list []*ast.Comment) {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	text := (&amp;ast.CommentGroup{List: list}).Text()
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	if m := noteMarkerRx.FindStringSubmatchIndex(text); m != nil {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		<span class="comment">// The note body starts after the marker.</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		<span class="comment">// We remove any formatting so that we don&#39;t</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		<span class="comment">// get spurious line breaks/indentation when</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		<span class="comment">// showing the TODO body.</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		body := clean(text[m[1]:])
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		if body != &#34;&#34; {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>			marker := text[m[2]:m[3]]
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			r.notes[marker] = append(r.notes[marker], &amp;Note{
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>				Pos:  list[0].Pos(),
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>				End:  list[len(list)-1].End(),
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>				UID:  text[m[4]:m[5]],
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>				Body: body,
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>			})
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span><span class="comment">// readNotes extracts notes from comments.</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">// A note must start at the beginning of a comment with &#34;MARKER(uid):&#34;</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">// and is followed by the note body (e.g., &#34;// BUG(gri): fix this&#34;).</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">// The note ends at the end of the comment group or at the start of</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">// another note in the same comment group, whichever comes first.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>func (r *reader) readNotes(comments []*ast.CommentGroup) {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	for _, group := range comments {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		i := -1 <span class="comment">// comment index of most recent note start, valid if &gt;= 0</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		list := group.List
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		for j, c := range list {
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>			if noteCommentRx.MatchString(c.Text) {
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>				if i &gt;= 0 {
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>					r.readNote(list[i:j])
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>				}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>				i = j
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>			}
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		}
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		if i &gt;= 0 {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			r.readNote(list[i:])
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span><span class="comment">// readFile adds the AST for a source file to the reader.</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>func (r *reader) readFile(src *ast.File) {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	<span class="comment">// add package documentation</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	if src.Doc != nil {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		r.readDoc(src.Doc)
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		if r.mode&amp;PreserveAST == 0 {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>			src.Doc = nil <span class="comment">// doc consumed - remove from AST</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	<span class="comment">// add all declarations but for functions which are processed in a separate pass</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	for _, decl := range src.Decls {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		switch d := decl.(type) {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		case *ast.GenDecl:
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>			switch d.Tok {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>			case token.IMPORT:
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>				<span class="comment">// imports are handled individually</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>				for _, spec := range d.Specs {
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>					if s, ok := spec.(*ast.ImportSpec); ok {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>						if import_, err := strconv.Unquote(s.Path.Value); err == nil {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>							r.imports[import_] = 1
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>							var name string
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>							if s.Name != nil {
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>								name = s.Name.Name
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>								if name == &#34;.&#34; {
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>									r.hasDotImp = true
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>								}
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>							}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>							if name != &#34;.&#34; {
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>								if name == &#34;&#34; {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>									name = assumedPackageName(import_)
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>								}
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>								old, ok := r.importByName[name]
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>								if !ok {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>									r.importByName[name] = import_
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>								} else if old != import_ &amp;&amp; old != &#34;&#34; {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>									r.importByName[name] = &#34;&#34; <span class="comment">// ambiguous</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>								}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>							}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>						}
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>					}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>				}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>			case token.CONST, token.VAR:
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>				<span class="comment">// constants and variables are always handled as a group</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>				r.readValue(d)
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			case token.TYPE:
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>				<span class="comment">// types are handled individually</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>				if len(d.Specs) == 1 &amp;&amp; !d.Lparen.IsValid() {
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>					<span class="comment">// common case: single declaration w/o parentheses</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>					<span class="comment">// (if a single declaration is parenthesized,</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>					<span class="comment">// create a new fake declaration below, so that</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>					<span class="comment">// go/doc type declarations always appear w/o</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>					<span class="comment">// parentheses)</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>					if s, ok := d.Specs[0].(*ast.TypeSpec); ok {
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>						r.readType(d, s)
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>					}
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>					break
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>				}
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>				for _, spec := range d.Specs {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>					if s, ok := spec.(*ast.TypeSpec); ok {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>						<span class="comment">// use an individual (possibly fake) declaration</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>						<span class="comment">// for each type; this also ensures that each type</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>						<span class="comment">// gets to (re-)use the declaration documentation</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>						<span class="comment">// if there&#39;s none associated with the spec itself</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>						fake := &amp;ast.GenDecl{
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>							Doc: d.Doc,
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>							<span class="comment">// don&#39;t use the existing TokPos because it</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>							<span class="comment">// will lead to the wrong selection range for</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>							<span class="comment">// the fake declaration if there are more</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>							<span class="comment">// than one type in the group (this affects</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>							<span class="comment">// src/cmd/godoc/godoc.go&#39;s posLink_urlFunc)</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>							TokPos: s.Pos(),
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>							Tok:    token.TYPE,
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>							Specs:  []ast.Spec{s},
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>						}
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>						r.readType(fake, s)
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>					}
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>				}
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>			}
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		}
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	}
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	<span class="comment">// collect MARKER(...): annotations</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	r.readNotes(src.Comments)
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	if r.mode&amp;PreserveAST == 0 {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		src.Comments = nil <span class="comment">// consumed unassociated comments - remove from AST</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>func (r *reader) readPackage(pkg *ast.Package, mode Mode) {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	<span class="comment">// initialize reader</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	r.filenames = make([]string, len(pkg.Files))
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	r.imports = make(map[string]int)
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	r.mode = mode
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	r.types = make(map[string]*namedType)
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	r.funcs = make(methodSet)
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	r.notes = make(map[string][]*Note)
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	r.importByName = make(map[string]string)
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	<span class="comment">// sort package files before reading them so that the</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	<span class="comment">// result does not depend on map iteration order</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	i := 0
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	for filename := range pkg.Files {
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		r.filenames[i] = filename
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		i++
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	}
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	sort.Strings(r.filenames)
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	<span class="comment">// process files in sorted order</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	for _, filename := range r.filenames {
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		f := pkg.Files[filename]
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>		if mode&amp;AllDecls == 0 {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>			r.fileExports(f)
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		r.readFile(f)
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	}
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	for name, path := range r.importByName {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		if path == &#34;&#34; {
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>			delete(r.importByName, name)
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		}
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	}
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	<span class="comment">// process functions now that we have better type information</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	for _, f := range pkg.Files {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		for _, decl := range f.Decls {
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>			if d, ok := decl.(*ast.FuncDecl); ok {
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>				r.readFunc(d)
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>			}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		}
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	}
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>}
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span><span class="comment">// Types</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>func customizeRecv(f *Func, recvTypeName string, embeddedIsPtr bool, level int) *Func {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	if f == nil || f.Decl == nil || f.Decl.Recv == nil || len(f.Decl.Recv.List) != 1 {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		return f <span class="comment">// shouldn&#39;t happen, but be safe</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	}
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	<span class="comment">// copy existing receiver field and set new type</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	newField := *f.Decl.Recv.List[0]
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	origPos := newField.Type.Pos()
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	_, origRecvIsPtr := newField.Type.(*ast.StarExpr)
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	newIdent := &amp;ast.Ident{NamePos: origPos, Name: recvTypeName}
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	var typ ast.Expr = newIdent
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	if !embeddedIsPtr &amp;&amp; origRecvIsPtr {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>		newIdent.NamePos++ <span class="comment">// &#39;*&#39; is one character</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		typ = &amp;ast.StarExpr{Star: origPos, X: newIdent}
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	}
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	newField.Type = typ
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	<span class="comment">// copy existing receiver field list and set new receiver field</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	newFieldList := *f.Decl.Recv
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	newFieldList.List = []*ast.Field{&amp;newField}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	<span class="comment">// copy existing function declaration and set new receiver field list</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	newFuncDecl := *f.Decl
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	newFuncDecl.Recv = &amp;newFieldList
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	<span class="comment">// copy existing function documentation and set new declaration</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	newF := *f
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	newF.Decl = &amp;newFuncDecl
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	newF.Recv = recvString(typ)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	<span class="comment">// the Orig field never changes</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	newF.Level = level
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	return &amp;newF
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>}
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span><span class="comment">// collectEmbeddedMethods collects the embedded methods of typ in mset.</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>func (r *reader) collectEmbeddedMethods(mset methodSet, typ *namedType, recvTypeName string, embeddedIsPtr bool, level int, visited embeddedSet) {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	visited[typ] = true
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	for embedded, isPtr := range typ.embedded {
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>		<span class="comment">// Once an embedded type is embedded as a pointer type</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>		<span class="comment">// all embedded types in those types are treated like</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>		<span class="comment">// pointer types for the purpose of the receiver type</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		<span class="comment">// computation; i.e., embeddedIsPtr is sticky for this</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		<span class="comment">// embedding hierarchy.</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>		thisEmbeddedIsPtr := embeddedIsPtr || isPtr
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		for _, m := range embedded.methods {
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			<span class="comment">// only top-level methods are embedded</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>			if m.Level == 0 {
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>				mset.add(customizeRecv(m, recvTypeName, thisEmbeddedIsPtr, level))
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>			}
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>		}
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		if !visited[embedded] {
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>			r.collectEmbeddedMethods(mset, embedded, recvTypeName, thisEmbeddedIsPtr, level+1, visited)
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		}
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	}
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	delete(visited, typ)
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>}
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span><span class="comment">// computeMethodSets determines the actual method sets for each type encountered.</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>func (r *reader) computeMethodSets() {
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	for _, t := range r.types {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>		<span class="comment">// collect embedded methods for t</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		if t.isStruct {
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>			<span class="comment">// struct</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>			r.collectEmbeddedMethods(t.methods, t, t.name, false, 1, make(embeddedSet))
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		} else {
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>			<span class="comment">// interface</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) fix this</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>		}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	}
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	<span class="comment">// For any predeclared names that are declared locally, don&#39;t treat them as</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	<span class="comment">// exported fields anymore.</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	for predecl := range r.shadowedPredecl {
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		for _, ityp := range r.fixmap[predecl] {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			removeAnonymousField(predecl, ityp)
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>		}
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	}
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>}
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span><span class="comment">// cleanupTypes removes the association of functions and methods with</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span><span class="comment">// types that have no declaration. Instead, these functions and methods</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span><span class="comment">// are shown at the package level. It also removes types with missing</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span><span class="comment">// declarations or which are not visible.</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>func (r *reader) cleanupTypes() {
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	for _, t := range r.types {
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>		visible := r.isVisible(t.name)
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>		predeclared := predeclaredTypes[t.name]
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>		if t.decl == nil &amp;&amp; (predeclared || visible &amp;&amp; (t.isEmbedded || r.hasDotImp)) {
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>			<span class="comment">// t.name is a predeclared type (and was not redeclared in this package),</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>			<span class="comment">// or it was embedded somewhere but its declaration is missing (because</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>			<span class="comment">// the AST is incomplete), or we have a dot-import (and all bets are off):</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>			<span class="comment">// move any associated values, funcs, and methods back to the top-level so</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>			<span class="comment">// that they are not lost.</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>			<span class="comment">// 1) move values</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>			r.values = append(r.values, t.values...)
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>			<span class="comment">// 2) move factory functions</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>			for name, f := range t.funcs {
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>				<span class="comment">// in a correct AST, package-level function names</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>				<span class="comment">// are all different - no need to check for conflicts</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>				r.funcs[name] = f
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>			}
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>			<span class="comment">// 3) move methods</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			if !predeclared {
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>				for name, m := range t.methods {
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>					<span class="comment">// don&#39;t overwrite functions with the same name - drop them</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>					if _, found := r.funcs[name]; !found {
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>						r.funcs[name] = m
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>					}
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>				}
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>			}
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		}
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>		<span class="comment">// remove types w/o declaration or which are not visible</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>		if t.decl == nil || !visible {
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>			delete(r.types, t.name)
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>		}
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	}
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>}
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span><span class="comment">// Sorting</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>type data struct {
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	n    int
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	swap func(i, j int)
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	less func(i, j int) bool
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>}
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>func (d *data) Len() int           { return d.n }
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>func (d *data) Swap(i, j int)      { d.swap(i, j) }
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>func (d *data) Less(i, j int) bool { return d.less(i, j) }
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span><span class="comment">// sortBy is a helper function for sorting.</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>func sortBy(less func(i, j int) bool, swap func(i, j int), n int) {
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	sort.Sort(&amp;data{n, swap, less})
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>}
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>func sortedKeys(m map[string]int) []string {
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	list := make([]string, len(m))
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>	i := 0
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	for key := range m {
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		list[i] = key
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		i++
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	}
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	sort.Strings(list)
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	return list
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span><span class="comment">// sortingName returns the name to use when sorting d into place.</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>func sortingName(d *ast.GenDecl) string {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	if len(d.Specs) == 1 {
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>		if s, ok := d.Specs[0].(*ast.ValueSpec); ok {
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>			return s.Names[0].Name
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	}
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	return &#34;&#34;
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>}
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>func sortedValues(m []*Value, tok token.Token) []*Value {
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	list := make([]*Value, len(m)) <span class="comment">// big enough in any case</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	i := 0
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	for _, val := range m {
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		if val.Decl.Tok == tok {
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>			list[i] = val
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>			i++
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		}
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	}
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	list = list[0:i]
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	sortBy(
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		func(i, j int) bool {
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>			if ni, nj := sortingName(list[i].Decl), sortingName(list[j].Decl); ni != nj {
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>				return ni &lt; nj
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>			}
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>			return list[i].order &lt; list[j].order
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		},
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		func(i, j int) { list[i], list[j] = list[j], list[i] },
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		len(list),
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	)
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	return list
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>}
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>func sortedTypes(m map[string]*namedType, allMethods bool) []*Type {
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	list := make([]*Type, len(m))
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	i := 0
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	for _, t := range m {
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>		list[i] = &amp;Type{
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>			Doc:     t.doc,
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>			Name:    t.name,
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>			Decl:    t.decl,
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>			Consts:  sortedValues(t.values, token.CONST),
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			Vars:    sortedValues(t.values, token.VAR),
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>			Funcs:   sortedFuncs(t.funcs, true),
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			Methods: sortedFuncs(t.methods, allMethods),
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>		}
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>		i++
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	}
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	sortBy(
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>		func(i, j int) bool { return list[i].Name &lt; list[j].Name },
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		func(i, j int) { list[i], list[j] = list[j], list[i] },
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>		len(list),
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>	)
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	return list
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>}
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>func removeStar(s string) string {
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	if len(s) &gt; 0 &amp;&amp; s[0] == &#39;*&#39; {
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>		return s[1:]
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	}
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	return s
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>}
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>func sortedFuncs(m methodSet, allMethods bool) []*Func {
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	list := make([]*Func, len(m))
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	i := 0
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	for _, m := range m {
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>		<span class="comment">// determine which methods to include</span>
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		switch {
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		case m.Decl == nil:
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>			<span class="comment">// exclude conflict entry</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>		case allMethods, m.Level == 0, !token.IsExported(removeStar(m.Orig)):
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>			<span class="comment">// forced inclusion, method not embedded, or method</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>			<span class="comment">// embedded but original receiver type not exported</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>			list[i] = m
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>			i++
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>		}
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	}
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	list = list[0:i]
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	sortBy(
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		func(i, j int) bool { return list[i].Name &lt; list[j].Name },
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>		func(i, j int) { list[i], list[j] = list[j], list[i] },
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>		len(list),
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>	)
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>	return list
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>}
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span><span class="comment">// noteBodies returns a list of note body strings given a list of notes.</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span><span class="comment">// This is only used to populate the deprecated Package.Bugs field.</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>func noteBodies(notes []*Note) []string {
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	var list []string
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	for _, n := range notes {
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		list = append(list, n.Body)
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>	}
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>	return list
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>}
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span><span class="comment">// Predeclared identifiers</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>
<span id="L949" class="ln">   949&nbsp;&nbsp;</span><span class="comment">// IsPredeclared reports whether s is a predeclared identifier.</span>
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>func IsPredeclared(s string) bool {
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>	return predeclaredTypes[s] || predeclaredFuncs[s] || predeclaredConstants[s]
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>}
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>var predeclaredTypes = map[string]bool{
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>	&#34;any&#34;:        true,
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>	&#34;bool&#34;:       true,
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	&#34;byte&#34;:       true,
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>	&#34;comparable&#34;: true,
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>	&#34;complex64&#34;:  true,
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>	&#34;complex128&#34;: true,
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	&#34;error&#34;:      true,
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>	&#34;float32&#34;:    true,
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>	&#34;float64&#34;:    true,
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	&#34;int&#34;:        true,
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	&#34;int8&#34;:       true,
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>	&#34;int16&#34;:      true,
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>	&#34;int32&#34;:      true,
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>	&#34;int64&#34;:      true,
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>	&#34;rune&#34;:       true,
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	&#34;string&#34;:     true,
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>	&#34;uint&#34;:       true,
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>	&#34;uint8&#34;:      true,
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	&#34;uint16&#34;:     true,
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>	&#34;uint32&#34;:     true,
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>	&#34;uint64&#34;:     true,
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	&#34;uintptr&#34;:    true,
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>}
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>var predeclaredFuncs = map[string]bool{
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	&#34;append&#34;:  true,
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	&#34;cap&#34;:     true,
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	&#34;close&#34;:   true,
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>	&#34;complex&#34;: true,
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>	&#34;copy&#34;:    true,
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	&#34;delete&#34;:  true,
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>	&#34;imag&#34;:    true,
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	&#34;len&#34;:     true,
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>	&#34;make&#34;:    true,
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	&#34;new&#34;:     true,
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	&#34;panic&#34;:   true,
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>	&#34;print&#34;:   true,
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>	&#34;println&#34;: true,
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	&#34;real&#34;:    true,
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	&#34;recover&#34;: true,
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>}
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>var predeclaredConstants = map[string]bool{
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	&#34;false&#34;: true,
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	&#34;iota&#34;:  true,
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	&#34;nil&#34;:   true,
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>	&#34;true&#34;:  true,
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>}
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span><span class="comment">// assumedPackageName returns the assumed package name</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span><span class="comment">// for a given import path. This is a copy of</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span><span class="comment">// golang.org/x/tools/internal/imports.ImportPathToAssumedName.</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>func assumedPackageName(importPath string) string {
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	notIdentifier := func(ch rune) bool {
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>		return !(&#39;a&#39; &lt;= ch &amp;&amp; ch &lt;= &#39;z&#39; || &#39;A&#39; &lt;= ch &amp;&amp; ch &lt;= &#39;Z&#39; ||
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>			&#39;0&#39; &lt;= ch &amp;&amp; ch &lt;= &#39;9&#39; ||
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>			ch == &#39;_&#39; ||
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>			ch &gt;= utf8.RuneSelf &amp;&amp; (unicode.IsLetter(ch) || unicode.IsDigit(ch)))
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	}
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	base := path.Base(importPath)
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	if strings.HasPrefix(base, &#34;v&#34;) {
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>		if _, err := strconv.Atoi(base[1:]); err == nil {
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>			dir := path.Dir(importPath)
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>			if dir != &#34;.&#34; {
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>				base = path.Base(dir)
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>			}
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>		}
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	}
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>	base = strings.TrimPrefix(base, &#34;go-&#34;)
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>	if i := strings.IndexFunc(base, notIdentifier); i &gt;= 0 {
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>		base = base[:i]
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>	}
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	return base
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>}
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>
</pre><p><a href="reader.go?m=text">View as plain text</a></p>

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
