<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/resolver.go - Go Documentation Server</title>

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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">resolver.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package types
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;go/constant&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/internal/typeparams&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// A declInfo describes a package-level const, type, var, or func declaration.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>type declInfo struct {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	file      *Scope        <span class="comment">// scope of file containing this declaration</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	lhs       []*Var        <span class="comment">// lhs of n:1 variable declarations, or nil</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	vtyp      ast.Expr      <span class="comment">// type, or nil (for const and var declarations only)</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	init      ast.Expr      <span class="comment">// init/orig expression, or nil (for const and var declarations only)</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	inherited bool          <span class="comment">// if set, the init expression is inherited from a previous constant declaration</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	tdecl     *ast.TypeSpec <span class="comment">// type declaration, or nil</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	fdecl     *ast.FuncDecl <span class="comment">// func declaration, or nil</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// The deps field tracks initialization expression dependencies.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	deps map[Object]bool <span class="comment">// lazily initialized</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// hasInitializer reports whether the declared object has an initialization</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// expression or function body.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>func (d *declInfo) hasInitializer() bool {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	return d.init != nil || d.fdecl != nil &amp;&amp; d.fdecl.Body != nil
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// addDep adds obj to the set of objects d&#39;s init expression depends on.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>func (d *declInfo) addDep(obj Object) {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	m := d.deps
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	if m == nil {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		m = make(map[Object]bool)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		d.deps = m
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	m[obj] = true
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// arityMatch checks that the lhs and rhs of a const or var decl</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// have the appropriate number of names and init exprs. For const</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// decls, init is the value spec providing the init exprs; for</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// var decls, init is nil (the init exprs are in s in this case).</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func (check *Checker) arityMatch(s, init *ast.ValueSpec) {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	l := len(s.Names)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	r := len(s.Values)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	if init != nil {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		r = len(init.Values)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	const code = WrongAssignCount
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	switch {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	case init == nil &amp;&amp; r == 0:
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		<span class="comment">// var decl w/o init expr</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		if s.Type == nil {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			check.error(s, code, &#34;missing type or init expr&#34;)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	case l &lt; r:
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		if l &lt; len(s.Values) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			<span class="comment">// init exprs from s</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>			n := s.Values[l]
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			check.errorf(n, code, &#34;extra init expr %s&#34;, n)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) avoid declared and not used error here</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		} else {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			<span class="comment">// init exprs &#34;inherited&#34;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			check.errorf(s, code, &#34;extra init expr at %s&#34;, check.fset.Position(init.Pos()))
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) avoid declared and not used error here</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	case l &gt; r &amp;&amp; (init != nil || r != 1):
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		n := s.Names[r]
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		check.errorf(n, code, &#34;missing init expr for %s&#34;, n)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>func validatedImportPath(path string) (string, error) {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	s, err := strconv.Unquote(path)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	if err != nil {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		return &#34;&#34;, err
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	if s == &#34;&#34; {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		return &#34;&#34;, fmt.Errorf(&#34;empty string&#34;)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	const illegalChars = `!&#34;#$%&amp;&#39;()*,:;&lt;=&gt;?[\]^{|}` + &#34;`\uFFFD&#34;
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	for _, r := range s {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		if !unicode.IsGraphic(r) || unicode.IsSpace(r) || strings.ContainsRune(illegalChars, r) {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			return s, fmt.Errorf(&#34;invalid character %#U&#34;, r)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	return s, nil
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// declarePkgObj declares obj in the package scope, records its ident -&gt; obj mapping,</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// and updates check.objMap. The object must not be a function or method.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>func (check *Checker) declarePkgObj(ident *ast.Ident, obj Object, d *declInfo) {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	assert(ident.Name == obj.Name())
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;A package-scope or file-scope identifier with name init</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">// may only be declared to be a function with this (func()) signature.&#34;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	if ident.Name == &#34;init&#34; {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		check.error(ident, InvalidInitDecl, &#34;cannot declare init - must be func&#34;)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		return
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;The main package must have package name main and declare</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// a function main that takes no arguments and returns no value.&#34;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	if ident.Name == &#34;main&#34; &amp;&amp; check.pkg.name == &#34;main&#34; {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		check.error(ident, InvalidMainDecl, &#34;cannot declare main - must be func&#34;)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		return
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	check.declare(check.pkg.scope, ident, obj, nopos)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	check.objMap[obj] = d
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	obj.setOrder(uint32(len(check.objMap)))
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// filename returns a filename suitable for debugging output.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>func (check *Checker) filename(fileNo int) string {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	file := check.files[fileNo]
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if pos := file.Pos(); pos.IsValid() {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		return check.fset.File(pos).Name()
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;file[%d]&#34;, fileNo)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>func (check *Checker) importPackage(at positioner, path, dir string) *Package {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// If we already have a package for the given (path, dir)</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// pair, use it instead of doing a full import.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// Checker.impMap only caches packages that are marked Complete</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// or fake (dummy packages for failed imports). Incomplete but</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// non-fake packages do require an import to complete them.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	key := importKey{path, dir}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	imp := check.impMap[key]
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	if imp != nil {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		return imp
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">// no package yet =&gt; import it</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	if path == &#34;C&#34; &amp;&amp; (check.conf.FakeImportC || check.conf.go115UsesCgo) {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		imp = NewPackage(&#34;C&#34;, &#34;C&#34;)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		imp.fake = true <span class="comment">// package scope is not populated</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		imp.cgo = check.conf.go115UsesCgo
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	} else {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		<span class="comment">// ordinary import</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		var err error
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		if importer := check.conf.Importer; importer == nil {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			err = fmt.Errorf(&#34;Config.Importer not installed&#34;)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		} else if importerFrom, ok := importer.(ImporterFrom); ok {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			imp, err = importerFrom.ImportFrom(path, dir, 0)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			if imp == nil &amp;&amp; err == nil {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>				err = fmt.Errorf(&#34;Config.Importer.ImportFrom(%s, %s, 0) returned nil but no error&#34;, path, dir)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		} else {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			imp, err = importer.Import(path)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			if imp == nil &amp;&amp; err == nil {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>				err = fmt.Errorf(&#34;Config.Importer.Import(%s) returned nil but no error&#34;, path)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// make sure we have a valid package name</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		<span class="comment">// (errors here can only happen through manipulation of packages after creation)</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		if err == nil &amp;&amp; imp != nil &amp;&amp; (imp.name == &#34;_&#34; || imp.name == &#34;&#34;) {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			err = fmt.Errorf(&#34;invalid package name: %q&#34;, imp.name)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			imp = nil <span class="comment">// create fake package below</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		if err != nil {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			check.errorf(at, BrokenImport, &#34;could not import %s (%s)&#34;, path, err)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			if imp == nil {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>				<span class="comment">// create a new fake package</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>				<span class="comment">// come up with a sensible package name (heuristic)</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>				name := path
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>				if i := len(name); i &gt; 0 &amp;&amp; name[i-1] == &#39;/&#39; {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>					name = name[:i-1]
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>				}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>				if i := strings.LastIndex(name, &#34;/&#34;); i &gt;= 0 {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>					name = name[i+1:]
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>				}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>				imp = NewPackage(path, name)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			<span class="comment">// continue to use the package as best as we can</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			imp.fake = true <span class="comment">// avoid follow-up lookup failures</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// package should be complete or marked fake, but be cautious</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	if imp.complete || imp.fake {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		check.impMap[key] = imp
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		<span class="comment">// Once we&#39;ve formatted an error message, keep the pkgPathMap</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		<span class="comment">// up-to-date on subsequent imports. It is used for package</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		<span class="comment">// qualification in error messages.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		if check.pkgPathMap != nil {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			check.markImports(imp)
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		return imp
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// something went wrong (importer may have returned incomplete package without error)</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	return nil
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// collectObjects collects all file and package objects and inserts them</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// into their respective scopes. It also performs imports and associates</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// methods with receiver base type names.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>func (check *Checker) collectObjects() {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	pkg := check.pkg
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// pkgImports is the set of packages already imported by any package file seen</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	<span class="comment">// so far. Used to avoid duplicate entries in pkg.imports. Allocate and populate</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	<span class="comment">// it (pkg.imports may not be empty if we are checking test files incrementally).</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	<span class="comment">// Note that pkgImports is keyed by package (and thus package path), not by an</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	<span class="comment">// importKey value. Two different importKey values may map to the same package</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// which is why we cannot use the check.impMap here.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	var pkgImports = make(map[*Package]bool)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	for _, imp := range pkg.imports {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		pkgImports[imp] = true
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	type methodInfo struct {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		obj  *Func      <span class="comment">// method</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		ptr  bool       <span class="comment">// true if pointer receiver</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		recv *ast.Ident <span class="comment">// receiver type name</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	var methods []methodInfo <span class="comment">// collected methods with valid receivers and non-blank _ names</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	var fileScopes []*Scope
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	for fileNo, file := range check.files {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		<span class="comment">// The package identifier denotes the current package,</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		<span class="comment">// but there is no corresponding package object.</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		check.recordDef(file.Name, nil)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		<span class="comment">// Use the actual source file extent rather than *ast.File extent since the</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		<span class="comment">// latter doesn&#39;t include comments which appear at the start or end of the file.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		<span class="comment">// Be conservative and use the *ast.File extent if we don&#39;t have a *token.File.</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		pos, end := file.Pos(), file.End()
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		if f := check.fset.File(file.Pos()); f != nil {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			pos, end = token.Pos(f.Base()), token.Pos(f.Base()+f.Size())
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		fileScope := NewScope(pkg.scope, pos, end, check.filename(fileNo))
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		fileScopes = append(fileScopes, fileScope)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		check.recordScope(file, fileScope)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		<span class="comment">// determine file directory, necessary to resolve imports</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		<span class="comment">// FileName may be &#34;&#34; (typically for tests) in which case</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		<span class="comment">// we get &#34;.&#34; as the directory which is what we would want.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		fileDir := dir(check.fset.Position(file.Name.Pos()).Filename)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		check.walkDecls(file.Decls, func(d decl) {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			switch d := d.(type) {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			case importDecl:
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>				<span class="comment">// import package</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>				if d.spec.Path.Value == &#34;&#34; {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>					return <span class="comment">// error reported by parser</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>				}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>				path, err := validatedImportPath(d.spec.Path.Value)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				if err != nil {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>					check.errorf(d.spec.Path, BadImportPath, &#34;invalid import path (%s)&#34;, err)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>					return
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>				}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>				imp := check.importPackage(d.spec.Path, path, fileDir)
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				if imp == nil {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>					return
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>				}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>				<span class="comment">// local name overrides imported package name</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>				name := imp.name
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>				if d.spec.Name != nil {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>					name = d.spec.Name.Name
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>					if path == &#34;C&#34; {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>						<span class="comment">// match 1.17 cmd/compile (not prescribed by spec)</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>						check.error(d.spec.Name, ImportCRenamed, `cannot rename import &#34;C&#34;`)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>						return
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>					}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>				}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>				if name == &#34;init&#34; {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>					check.error(d.spec, InvalidInitDecl, &#34;cannot import package as init - init must be a func&#34;)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>					return
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>				}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>				<span class="comment">// add package to list of explicit imports</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>				<span class="comment">// (this functionality is provided as a convenience</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>				<span class="comment">// for clients; it is not needed for type-checking)</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>				if !pkgImports[imp] {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>					pkgImports[imp] = true
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>					pkg.imports = append(pkg.imports, imp)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>				}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>				pkgName := NewPkgName(d.spec.Pos(), pkg, name, imp)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>				if d.spec.Name != nil {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>					<span class="comment">// in a dot-import, the dot represents the package</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>					check.recordDef(d.spec.Name, pkgName)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>				} else {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>					check.recordImplicit(d.spec, pkgName)
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>				}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>				if imp.fake {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>					<span class="comment">// match 1.17 cmd/compile (not prescribed by spec)</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>					pkgName.used = true
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>				}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>				<span class="comment">// add import to file scope</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>				check.imports = append(check.imports, pkgName)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>				if name == &#34;.&#34; {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>					<span class="comment">// dot-import</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>					if check.dotImportMap == nil {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>						check.dotImportMap = make(map[dotImportKey]*PkgName)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>					}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>					<span class="comment">// merge imported scope with file scope</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>					for name, obj := range imp.scope.elems {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>						<span class="comment">// Note: Avoid eager resolve(name, obj) here, so we only</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>						<span class="comment">// resolve dot-imported objects as needed.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>						<span class="comment">// A package scope may contain non-exported objects,</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>						<span class="comment">// do not import them!</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>						if token.IsExported(name) {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>							<span class="comment">// declare dot-imported object</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>							<span class="comment">// (Do not use check.declare because it modifies the object</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>							<span class="comment">// via Object.setScopePos, which leads to a race condition;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>							<span class="comment">// the object may be imported into more than one file scope</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>							<span class="comment">// concurrently. See go.dev/issue/32154.)</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>							if alt := fileScope.Lookup(name); alt != nil {
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>								check.errorf(d.spec.Name, DuplicateDecl, &#34;%s redeclared in this block&#34;, alt.Name())
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>								check.reportAltDecl(alt)
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>							} else {
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>								fileScope.insert(name, obj)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>								check.dotImportMap[dotImportKey{fileScope, name}] = pkgName
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>							}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>						}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>					}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>				} else {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>					<span class="comment">// declare imported package object in file scope</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>					<span class="comment">// (no need to provide s.Name since we called check.recordDef earlier)</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>					check.declare(fileScope, nil, pkgName, nopos)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>				}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			case constDecl:
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>				<span class="comment">// declare all constants</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>				for i, name := range d.spec.Names {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>					obj := NewConst(name.Pos(), pkg, name.Name, nil, constant.MakeInt64(int64(d.iota)))
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>					var init ast.Expr
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>					if i &lt; len(d.init) {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>						init = d.init[i]
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>					}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>					d := &amp;declInfo{file: fileScope, vtyp: d.typ, init: init, inherited: d.inherited}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>					check.declarePkgObj(name, obj, d)
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>				}
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			case varDecl:
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>				lhs := make([]*Var, len(d.spec.Names))
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>				<span class="comment">// If there&#39;s exactly one rhs initializer, use</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>				<span class="comment">// the same declInfo d1 for all lhs variables</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>				<span class="comment">// so that each lhs variable depends on the same</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>				<span class="comment">// rhs initializer (n:1 var declaration).</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>				var d1 *declInfo
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>				if len(d.spec.Values) == 1 {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>					<span class="comment">// The lhs elements are only set up after the for loop below,</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>					<span class="comment">// but that&#39;s ok because declareVar only collects the declInfo</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>					<span class="comment">// for a later phase.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>					d1 = &amp;declInfo{file: fileScope, lhs: lhs, vtyp: d.spec.Type, init: d.spec.Values[0]}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>				}
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>				<span class="comment">// declare all variables</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>				for i, name := range d.spec.Names {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>					obj := NewVar(name.Pos(), pkg, name.Name, nil)
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>					lhs[i] = obj
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>					di := d1
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>					if di == nil {
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>						<span class="comment">// individual assignments</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>						var init ast.Expr
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>						if i &lt; len(d.spec.Values) {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>							init = d.spec.Values[i]
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>						}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>						di = &amp;declInfo{file: fileScope, vtyp: d.spec.Type, init: init}
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>					}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>					check.declarePkgObj(name, obj, di)
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>				}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			case typeDecl:
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>				_ = d.spec.TypeParams.NumFields() != 0 &amp;&amp; check.verifyVersionf(d.spec.TypeParams.List[0], go1_18, &#34;type parameter&#34;)
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>				obj := NewTypeName(d.spec.Name.Pos(), pkg, d.spec.Name.Name, nil)
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>				check.declarePkgObj(d.spec.Name, obj, &amp;declInfo{file: fileScope, tdecl: d.spec})
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			case funcDecl:
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>				name := d.decl.Name.Name
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>				obj := NewFunc(d.decl.Name.Pos(), pkg, name, nil)
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>				hasTParamError := false <span class="comment">// avoid duplicate type parameter errors</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>				if d.decl.Recv.NumFields() == 0 {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>					<span class="comment">// regular function</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>					if d.decl.Recv != nil {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>						check.error(d.decl.Recv, BadRecv, &#34;method has no receiver&#34;)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>						<span class="comment">// treat as function</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>					}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>					if name == &#34;init&#34; || (name == &#34;main&#34; &amp;&amp; check.pkg.name == &#34;main&#34;) {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>						code := InvalidInitDecl
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>						if name == &#34;main&#34; {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>							code = InvalidMainDecl
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>						}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>						if d.decl.Type.TypeParams.NumFields() != 0 {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>							check.softErrorf(d.decl.Type.TypeParams.List[0], code, &#34;func %s must have no type parameters&#34;, name)
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>							hasTParamError = true
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>						}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>						if t := d.decl.Type; t.Params.NumFields() != 0 || t.Results != nil {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>							<span class="comment">// TODO(rFindley) Should this be a hard error?</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>							check.softErrorf(d.decl.Name, code, &#34;func %s must have no arguments and no return values&#34;, name)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>						}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>					}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>					if name == &#34;init&#34; {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>						<span class="comment">// don&#39;t declare init functions in the package scope - they are invisible</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>						obj.parent = pkg.scope
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>						check.recordDef(d.decl.Name, obj)
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>						<span class="comment">// init functions must have a body</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>						if d.decl.Body == nil {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>							<span class="comment">// TODO(gri) make this error message consistent with the others above</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>							check.softErrorf(obj, MissingInitBody, &#34;missing function body&#34;)
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>						}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>					} else {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>						check.declare(pkg.scope, d.decl.Name, obj, nopos)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>					}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>				} else {
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>					<span class="comment">// method</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>					<span class="comment">// TODO(rFindley) earlier versions of this code checked that methods</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>					<span class="comment">//                have no type parameters, but this is checked later</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>					<span class="comment">//                when type checking the function type. Confirm that</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>					<span class="comment">//                we don&#39;t need to check tparams here.</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>					ptr, recv, _ := check.unpackRecv(d.decl.Recv.List[0].Type, false)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>					<span class="comment">// (Methods with invalid receiver cannot be associated to a type, and</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>					<span class="comment">// methods with blank _ names are never found; no need to collect any</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>					<span class="comment">// of them. They will still be type-checked with all the other functions.)</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>					if recv != nil &amp;&amp; name != &#34;_&#34; {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>						methods = append(methods, methodInfo{obj, ptr, recv})
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>					}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>					check.recordDef(d.decl.Name, obj)
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>				}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>				_ = d.decl.Type.TypeParams.NumFields() != 0 &amp;&amp; !hasTParamError &amp;&amp; check.verifyVersionf(d.decl.Type.TypeParams.List[0], go1_18, &#34;type parameter&#34;)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>				info := &amp;declInfo{file: fileScope, fdecl: d.decl}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>				<span class="comment">// Methods are not package-level objects but we still track them in the</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>				<span class="comment">// object map so that we can handle them like regular functions (if the</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>				<span class="comment">// receiver is invalid); also we need their fdecl info when associating</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>				<span class="comment">// them with their receiver base type, below.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>				check.objMap[obj] = info
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>				obj.setOrder(uint32(len(check.objMap)))
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		})
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	<span class="comment">// verify that objects in package and file scopes have different names</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	for _, scope := range fileScopes {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		for name, obj := range scope.elems {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			if alt := pkg.scope.Lookup(name); alt != nil {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>				obj = resolve(name, obj)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>				if pkg, ok := obj.(*PkgName); ok {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>					check.errorf(alt, DuplicateDecl, &#34;%s already declared through import of %s&#34;, alt.Name(), pkg.Imported())
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>					check.reportAltDecl(pkg)
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>				} else {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>					check.errorf(alt, DuplicateDecl, &#34;%s already declared through dot-import of %s&#34;, alt.Name(), obj.Pkg())
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>					<span class="comment">// TODO(gri) dot-imported objects don&#39;t have a position; reportAltDecl won&#39;t print anything</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>					check.reportAltDecl(obj)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>				}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>			}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	<span class="comment">// Now that we have all package scope objects and all methods,</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	<span class="comment">// associate methods with receiver base type name where possible.</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	<span class="comment">// Ignore methods that have an invalid receiver. They will be</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	<span class="comment">// type-checked later, with regular functions.</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	if methods == nil {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		return <span class="comment">// nothing to do</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	check.methods = make(map[*TypeName][]*Func)
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	for i := range methods {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		m := &amp;methods[i]
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		<span class="comment">// Determine the receiver base type and associate m with it.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		ptr, base := check.resolveBaseTypeName(m.ptr, m.recv, fileScopes)
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		if base != nil {
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>			m.obj.hasPtrRecv_ = ptr
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			check.methods[base] = append(check.methods[base], m.obj)
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span><span class="comment">// unpackRecv unpacks a receiver type and returns its components: ptr indicates whether</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span><span class="comment">// rtyp is a pointer receiver, rname is the receiver type name, and tparams are its</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span><span class="comment">// type parameters, if any. The type parameters are only unpacked if unpackParams is</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span><span class="comment">// set. If rname is nil, the receiver is unusable (i.e., the source has a bug which we</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span><span class="comment">// cannot easily work around).</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>func (check *Checker) unpackRecv(rtyp ast.Expr, unpackParams bool) (ptr bool, rname *ast.Ident, tparams []*ast.Ident) {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>L: <span class="comment">// unpack receiver type</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	<span class="comment">// This accepts invalid receivers such as ***T and does not</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	<span class="comment">// work for other invalid receivers, but we don&#39;t care. The</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	<span class="comment">// validity of receiver expressions is checked elsewhere.</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	for {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		switch t := rtyp.(type) {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		case *ast.ParenExpr:
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			rtyp = t.X
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		case *ast.StarExpr:
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			ptr = true
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			rtyp = t.X
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		default:
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			break L
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	<span class="comment">// unpack type parameters, if any</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	switch rtyp.(type) {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	case *ast.IndexExpr, *ast.IndexListExpr:
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		ix := typeparams.UnpackIndexExpr(rtyp)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		rtyp = ix.X
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		if unpackParams {
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			for _, arg := range ix.Indices {
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>				var par *ast.Ident
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>				switch arg := arg.(type) {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>				case *ast.Ident:
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>					par = arg
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>				case *ast.BadExpr:
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>					<span class="comment">// ignore - error already reported by parser</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>				case nil:
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>					check.error(ix.Orig, InvalidSyntaxTree, &#34;parameterized receiver contains nil parameters&#34;)
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>				default:
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>					check.errorf(arg, BadDecl, &#34;receiver type parameter %s must be an identifier&#34;, arg)
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>				}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>				if par == nil {
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>					par = &amp;ast.Ident{NamePos: arg.Pos(), Name: &#34;_&#34;}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>				}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>				tparams = append(tparams, par)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>			}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		}
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	<span class="comment">// unpack receiver name</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	if name, _ := rtyp.(*ast.Ident); name != nil {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		rname = name
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	}
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	return
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span><span class="comment">// resolveBaseTypeName returns the non-alias base type name for typ, and whether</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span><span class="comment">// there was a pointer indirection to get to it. The base type name must be declared</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span><span class="comment">// in package scope, and there can be at most one pointer indirection. If no such type</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// name exists, the returned base is nil.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>func (check *Checker) resolveBaseTypeName(seenPtr bool, typ ast.Expr, fileScopes []*Scope) (ptr bool, base *TypeName) {
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	<span class="comment">// Algorithm: Starting from a type expression, which may be a name,</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	<span class="comment">// we follow that type through alias declarations until we reach a</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	<span class="comment">// non-alias type name. If we encounter anything but pointer types or</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	<span class="comment">// parentheses we&#39;re done. If we encounter more than one pointer type</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	<span class="comment">// we&#39;re done.</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	ptr = seenPtr
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	var seen map[*TypeName]bool
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	for {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		<span class="comment">// Note: this differs from types2, but is necessary. The syntax parser</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		<span class="comment">// strips unnecessary parens.</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		typ = unparen(typ)
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		<span class="comment">// check if we have a pointer type</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		if pexpr, _ := typ.(*ast.StarExpr); pexpr != nil {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>			<span class="comment">// if we&#39;ve already seen a pointer, we&#39;re done</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>			if ptr {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>				return false, nil
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>			}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>			ptr = true
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>			typ = unparen(pexpr.X) <span class="comment">// continue with pointer base type</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		<span class="comment">// typ must be a name, or a C.name cgo selector.</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		var name string
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		switch typ := typ.(type) {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		case *ast.Ident:
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>			name = typ.Name
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		case *ast.SelectorExpr:
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>			<span class="comment">// C.struct_foo is a valid type name for packages using cgo.</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>			<span class="comment">// Detect this case, and adjust name so that the correct TypeName is</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>			<span class="comment">// resolved below.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>			if ident, _ := typ.X.(*ast.Ident); ident != nil &amp;&amp; ident.Name == &#34;C&#34; {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>				<span class="comment">// Check whether &#34;C&#34; actually resolves to an import of &#34;C&#34;, by looking</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>				<span class="comment">// in the appropriate file scope.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>				var obj Object
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>				for _, scope := range fileScopes {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>					if scope.Contains(ident.Pos()) {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>						obj = scope.Lookup(ident.Name)
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>					}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>				}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>				<span class="comment">// If Config.go115UsesCgo is set, the typechecker will resolve Cgo</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>				<span class="comment">// selectors to their cgo name. We must do the same here.</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>				if pname, _ := obj.(*PkgName); pname != nil {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>					if pname.imported.cgo { <span class="comment">// only set if Config.go115UsesCgo is set</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>						name = &#34;_Ctype_&#34; + typ.Sel.Name
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>					}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>				}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>			}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>			if name == &#34;&#34; {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>				return false, nil
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		default:
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>			return false, nil
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		<span class="comment">// name must denote an object found in the current package scope</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		<span class="comment">// (note that dot-imported objects are not in the package scope!)</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		obj := check.pkg.scope.Lookup(name)
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		if obj == nil {
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>			return false, nil
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		}
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		<span class="comment">// the object must be a type name...</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		tname, _ := obj.(*TypeName)
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		if tname == nil {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			return false, nil
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		}
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		<span class="comment">// ... which we have not seen before</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		if seen[tname] {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>			return false, nil
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		<span class="comment">// we&#39;re done if tdecl defined tname as a new type</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		<span class="comment">// (rather than an alias)</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		tdecl := check.objMap[tname].tdecl <span class="comment">// must exist for objects in package scope</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		if !tdecl.Assign.IsValid() {
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			return ptr, tname
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>		}
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		<span class="comment">// otherwise, continue resolving</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		typ = tdecl.Type
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		if seen == nil {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>			seen = make(map[*TypeName]bool)
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		}
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		seen[tname] = true
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span><span class="comment">// packageObjects typechecks all package objects, but not function bodies.</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>func (check *Checker) packageObjects() {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	<span class="comment">// process package objects in source order for reproducible results</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	objList := make([]Object, len(check.objMap))
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	i := 0
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	for obj := range check.objMap {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		objList[i] = obj
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		i++
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	}
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	sort.Sort(inSourceOrder(objList))
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	<span class="comment">// add new methods to already type-checked types (from a prior Checker.Files call)</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	for _, obj := range objList {
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		if obj, _ := obj.(*TypeName); obj != nil &amp;&amp; obj.typ != nil {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>			check.collectMethods(obj)
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		}
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	}
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	if check.enableAlias {
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		<span class="comment">// With Alias nodes we can process declarations in any order.</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		for _, obj := range objList {
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>			check.objDecl(obj, nil)
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	} else {
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		<span class="comment">// Without Alias nodes, we process non-alias type declarations first, followed by</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		<span class="comment">// alias declarations, and then everything else. This appears to avoid most situations</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		<span class="comment">// where the type of an alias is needed before it is available.</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>		<span class="comment">// There may still be cases where this is not good enough (see also go.dev/issue/25838).</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		<span class="comment">// In those cases Checker.ident will report an error (&#34;invalid use of type alias&#34;).</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		var aliasList []*TypeName
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		var othersList []Object <span class="comment">// everything that&#39;s not a type</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>		<span class="comment">// phase 1: non-alias type declarations</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		for _, obj := range objList {
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>			if tname, _ := obj.(*TypeName); tname != nil {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>				if check.objMap[tname].tdecl.Assign.IsValid() {
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>					aliasList = append(aliasList, tname)
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>				} else {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>					check.objDecl(obj, nil)
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>				}
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>			} else {
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>				othersList = append(othersList, obj)
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>			}
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		<span class="comment">// phase 2: alias type declarations</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		for _, obj := range aliasList {
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>			check.objDecl(obj, nil)
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		}
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		<span class="comment">// phase 3: all other declarations</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		for _, obj := range othersList {
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>			check.objDecl(obj, nil)
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		}
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	}
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	<span class="comment">// At this point we may have a non-empty check.methods map; this means that not all</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	<span class="comment">// entries were deleted at the end of typeDecl because the respective receiver base</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	<span class="comment">// types were not found. In that case, an error was reported when declaring those</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	<span class="comment">// methods. We can now safely discard this map.</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	check.methods = nil
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>}
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span><span class="comment">// inSourceOrder implements the sort.Sort interface.</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>type inSourceOrder []Object
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>func (a inSourceOrder) Len() int           { return len(a) }
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>func (a inSourceOrder) Less(i, j int) bool { return a[i].order() &lt; a[j].order() }
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>func (a inSourceOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span><span class="comment">// unusedImports checks for unused imports.</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>func (check *Checker) unusedImports() {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	<span class="comment">// If function bodies are not checked, packages&#39; uses are likely missing - don&#39;t check.</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	if check.conf.IgnoreFuncBodies {
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		return
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	}
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;It is illegal (...) to directly import a package without referring to</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	<span class="comment">// any of its exported identifiers. To import a package solely for its side-effects</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	<span class="comment">// (initialization), use the blank identifier as explicit package name.&#34;</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	for _, obj := range check.imports {
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		if !obj.used &amp;&amp; obj.name != &#34;_&#34; {
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>			check.errorUnusedPkg(obj)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		}
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	}
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>}
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>func (check *Checker) errorUnusedPkg(obj *PkgName) {
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	<span class="comment">// If the package was imported with a name other than the final</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	<span class="comment">// import path element, show it explicitly in the error message.</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	<span class="comment">// Note that this handles both renamed imports and imports of</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	<span class="comment">// packages containing unconventional package declarations.</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	<span class="comment">// Note that this uses / always, even on Windows, because Go import</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	<span class="comment">// paths always use forward slashes.</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	path := obj.imported.path
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	elem := path
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	if i := strings.LastIndex(elem, &#34;/&#34;); i &gt;= 0 {
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		elem = elem[i+1:]
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	}
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	if obj.name == &#34;&#34; || obj.name == &#34;.&#34; || obj.name == elem {
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		check.softErrorf(obj, UnusedImport, &#34;%q imported and not used&#34;, path)
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	} else {
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		check.softErrorf(obj, UnusedImport, &#34;%q imported as %s and not used&#34;, path, obj.name)
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	}
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>}
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span><span class="comment">// dir makes a good-faith attempt to return the directory</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span><span class="comment">// portion of path. If path is empty, the result is &#34;.&#34;.</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span><span class="comment">// (Per the go/build package dependency tests, we cannot import</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span><span class="comment">// path/filepath and simply use filepath.Dir.)</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>func dir(path string) string {
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	if i := strings.LastIndexAny(path, `/\`); i &gt; 0 {
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		return path[:i]
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	}
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	<span class="comment">// i &lt;= 0</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	return &#34;.&#34;
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>}
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>
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
