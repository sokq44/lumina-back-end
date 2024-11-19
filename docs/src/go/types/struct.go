<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/struct.go - Go Documentation Server</title>

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
<a href="struct.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">struct.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package types
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// API</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// A Struct represents a struct type.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type Struct struct {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	fields []*Var   <span class="comment">// fields != nil indicates the struct is set up (possibly with len(fields) == 0)</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	tags   []string <span class="comment">// field tags; nil if there are no tags</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>}
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// NewStruct returns a new struct with the given fields and corresponding field tags.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// If a field with index i has a tag, tags[i] must be that tag, but len(tags) may be</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// only as long as required to hold the tag with the largest index i. Consequently,</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// if no field has a tag, tags may be nil.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>func NewStruct(fields []*Var, tags []string) *Struct {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	var fset objset
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	for _, f := range fields {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		if f.name != &#34;_&#34; &amp;&amp; fset.insert(f) != nil {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>			panic(&#34;multiple fields with the same name&#34;)
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	if len(tags) &gt; len(fields) {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		panic(&#34;more tags than fields&#34;)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	s := &amp;Struct{fields: fields, tags: tags}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	s.markComplete()
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	return s
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// NumFields returns the number of fields in the struct (including blank and embedded fields).</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func (s *Struct) NumFields() int { return len(s.fields) }
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// Field returns the i&#39;th field for 0 &lt;= i &lt; NumFields().</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>func (s *Struct) Field(i int) *Var { return s.fields[i] }
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// Tag returns the i&#39;th field tag for 0 &lt;= i &lt; NumFields().</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>func (s *Struct) Tag(i int) string {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	if i &lt; len(s.tags) {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		return s.tags[i]
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	return &#34;&#34;
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>func (t *Struct) Underlying() Type { return t }
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>func (t *Struct) String() string   { return TypeString(t, nil) }
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// Implementation</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>func (s *Struct) markComplete() {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if s.fields == nil {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		s.fields = make([]*Var, 0)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>func (check *Checker) structType(styp *Struct, e *ast.StructType) {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	list := e.Fields
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	if list == nil {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		styp.markComplete()
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		return
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// struct fields and tags</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	var fields []*Var
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	var tags []string
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// for double-declaration checks</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	var fset objset
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// current field typ and tag</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	var typ Type
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	var tag string
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	add := func(ident *ast.Ident, embedded bool, pos token.Pos) {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		if tag != &#34;&#34; &amp;&amp; tags == nil {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			tags = make([]string, len(fields))
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		if tags != nil {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			tags = append(tags, tag)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		name := ident.Name
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		fld := NewField(pos, check.pkg, name, typ, embedded)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		<span class="comment">// spec: &#34;Within a struct, non-blank field names must be unique.&#34;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		if name == &#34;_&#34; || check.declareInSet(&amp;fset, pos, fld) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			fields = append(fields, fld)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			check.recordDef(ident, fld)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// addInvalid adds an embedded field of invalid type to the struct for</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// fields with errors; this keeps the number of struct fields in sync</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// with the source as long as the fields are _ or have different names</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// (go.dev/issue/25627).</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	addInvalid := func(ident *ast.Ident, pos token.Pos) {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		typ = Typ[Invalid]
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		tag = &#34;&#34;
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		add(ident, true, pos)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	for _, f := range list.List {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		typ = check.varType(f.Type)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		tag = check.tag(f.Tag)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		if len(f.Names) &gt; 0 {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			<span class="comment">// named fields</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			for _, name := range f.Names {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>				add(name, false, name.Pos())
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		} else {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			<span class="comment">// embedded field</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			<span class="comment">// spec: &#34;An embedded type must be specified as a type name T or as a</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			<span class="comment">// pointer to a non-interface type name *T, and T itself may not be a</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			<span class="comment">// pointer type.&#34;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			pos := f.Type.Pos() <span class="comment">// position of type, for errors</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			name := embeddedFieldIdent(f.Type)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			if name == nil {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>				check.errorf(f.Type, InvalidSyntaxTree, &#34;embedded field type %s has no name&#34;, f.Type)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>				name = ast.NewIdent(&#34;_&#34;)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>				name.NamePos = pos
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>				addInvalid(name, pos)
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>				continue
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			add(name, true, name.Pos()) <span class="comment">// struct{p.T} field has position of T</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			<span class="comment">// Because we have a name, typ must be of the form T or *T, where T is the name</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			<span class="comment">// of a (named or alias) type, and t (= deref(typ)) must be the type of T.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			<span class="comment">// We must delay this check to the end because we don&#39;t want to instantiate</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			<span class="comment">// (via under(t)) a possibly incomplete type.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			<span class="comment">// for use in the closure below</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			embeddedTyp := typ
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			embeddedPos := f.Type
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			check.later(func() {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>				t, isPtr := deref(embeddedTyp)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>				switch u := under(t).(type) {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				case *Basic:
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>					if !isValid(t) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>						<span class="comment">// error was reported before</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>						return
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>					}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>					<span class="comment">// unsafe.Pointer is treated like a regular pointer</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>					if u.kind == UnsafePointer {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>						check.error(embeddedPos, InvalidPtrEmbed, &#34;embedded field type cannot be unsafe.Pointer&#34;)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>					}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>				case *Pointer:
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>					check.error(embeddedPos, InvalidPtrEmbed, &#34;embedded field type cannot be a pointer&#34;)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>				case *Interface:
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>					if isTypeParam(t) {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>						<span class="comment">// The error code here is inconsistent with other error codes for</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>						<span class="comment">// invalid embedding, because this restriction may be relaxed in the</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>						<span class="comment">// future, and so it did not warrant a new error code.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>						check.error(embeddedPos, MisplacedTypeParam, &#34;embedded field type cannot be a (pointer to a) type parameter&#34;)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>						break
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>					}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>					if isPtr {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>						check.error(embeddedPos, InvalidPtrEmbed, &#34;embedded field type cannot be a pointer to an interface&#34;)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>					}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>				}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			}).describef(embeddedPos, &#34;check embedded type %s&#34;, embeddedTyp)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	styp.fields = fields
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	styp.tags = tags
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	styp.markComplete()
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>func embeddedFieldIdent(e ast.Expr) *ast.Ident {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	switch e := e.(type) {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	case *ast.Ident:
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		return e
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	case *ast.StarExpr:
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		<span class="comment">// *T is valid, but **T is not</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		if _, ok := e.X.(*ast.StarExpr); !ok {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			return embeddedFieldIdent(e.X)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	case *ast.SelectorExpr:
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		return e.Sel
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	case *ast.IndexExpr:
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		return embeddedFieldIdent(e.X)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	case *ast.IndexListExpr:
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		return embeddedFieldIdent(e.X)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	return nil <span class="comment">// invalid embedded field</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>func (check *Checker) declareInSet(oset *objset, pos token.Pos, obj Object) bool {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	if alt := oset.insert(obj); alt != nil {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		check.errorf(atPos(pos), DuplicateDecl, &#34;%s redeclared&#34;, obj.Name())
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		check.reportAltDecl(alt)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		return false
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	return true
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>func (check *Checker) tag(t *ast.BasicLit) string {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	if t != nil {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		if t.Kind == token.STRING {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			if val, err := strconv.Unquote(t.Value); err == nil {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>				return val
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		check.errorf(t, InvalidSyntaxTree, &#34;incorrect tag syntax: %q&#34;, t.Value)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	return &#34;&#34;
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
</pre><p><a href="struct.go?m=text">View as plain text</a></p>

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
