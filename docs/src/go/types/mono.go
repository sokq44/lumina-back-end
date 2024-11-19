<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/mono.go - Go Documentation Server</title>

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
<a href="mono.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">mono.go</span>
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
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// This file implements a check to validate that a Go package doesn&#39;t</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// have unbounded recursive instantiation, which is not compatible</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// with compilers using static instantiation (such as</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// monomorphization).</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// It implements a sort of &#34;type flow&#34; analysis by detecting which</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// type parameters are instantiated with other type parameters (or</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// types derived thereof). A package cannot be statically instantiated</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// if the graph has any cycles involving at least one derived type.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// Concretely, we construct a directed, weighted graph. Vertices are</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// used to represent type parameters as well as some defined</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// types. Edges are used to represent how types depend on each other:</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// * Everywhere a type-parameterized function or type is instantiated,</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//   we add edges to each type parameter from the vertices (if any)</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//   representing each type parameter or defined type referenced by</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//   the type argument. If the type argument is just the referenced</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//   type itself, then the edge has weight 0, otherwise 1.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// * For every defined type declared within a type-parameterized</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//   function or method, we add an edge of weight 1 to the defined</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//   type from each ambient type parameter.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// For example, given:</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//	func f[A, B any]() {</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//		type T int</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//		f[T, map[A]B]()</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// we construct vertices representing types A, B, and T. Because of</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// declaration &#34;type T int&#34;, we construct edges T&lt;-A and T&lt;-B with</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// weight 1; and because of instantiation &#34;f[T, map[A]B]&#34; we construct</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// edges A&lt;-T with weight 0, and B&lt;-A and B&lt;-B with weight 1.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// Finally, we look for any positive-weight cycles. Zero-weight cycles</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// are allowed because static instantiation will reach a fixed point.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>type monoGraph struct {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	vertices []monoVertex
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	edges    []monoEdge
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// canon maps method receiver type parameters to their respective</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// receiver type&#39;s type parameters.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	canon map[*TypeParam]*TypeParam
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// nameIdx maps a defined type or (canonical) type parameter to its</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// vertex index.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	nameIdx map[*TypeName]int
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>type monoVertex struct {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	weight int <span class="comment">// weight of heaviest known path to this vertex</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	pre    int <span class="comment">// previous edge (if any) in the above path</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	len    int <span class="comment">// length of the above path</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// obj is the defined type or type parameter represented by this</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// vertex.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	obj *TypeName
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>type monoEdge struct {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	dst, src int
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	weight   int
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	pos token.Pos
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	typ Type
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func (check *Checker) monomorph() {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// We detect unbounded instantiation cycles using a variant of</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// Bellman-Ford&#39;s algorithm. Namely, instead of always running |V|</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// iterations, we run until we either reach a fixed point or we&#39;ve</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// found a path of length |V|. This allows us to terminate earlier</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// when there are no cycles, which should be the common case.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	again := true
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	for again {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		again = false
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		for i, edge := range check.mono.edges {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			src := &amp;check.mono.vertices[edge.src]
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			dst := &amp;check.mono.vertices[edge.dst]
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			<span class="comment">// N.B., we&#39;re looking for the greatest weight paths, unlike</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			<span class="comment">// typical Bellman-Ford.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			w := src.weight + edge.weight
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			if w &lt;= dst.weight {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>				continue
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			dst.pre = i
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			dst.len = src.len + 1
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			if dst.len == len(check.mono.vertices) {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>				check.reportInstanceLoop(edge.dst)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>				return
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			dst.weight = w
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			again = true
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>func (check *Checker) reportInstanceLoop(v int) {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	var stack []int
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	seen := make([]bool, len(check.mono.vertices))
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// We have a path that contains a cycle and ends at v, but v may</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// only be reachable from the cycle, not on the cycle itself. We</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// start by walking backwards along the path until we find a vertex</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// that appears twice.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	for !seen[v] {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		stack = append(stack, v)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		seen[v] = true
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		v = check.mono.edges[check.mono.vertices[v].pre].src
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// Trim any vertices we visited before visiting v the first</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// time. Since v is the first vertex we found within the cycle, any</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// vertices we visited earlier cannot be part of the cycle.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	for stack[0] != v {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		stack = stack[1:]
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// TODO(mdempsky): Pivot stack so we report the cycle from the top?</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	obj0 := check.mono.vertices[v].obj
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	check.error(obj0, InvalidInstanceCycle, &#34;instantiation cycle:&#34;)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	qf := RelativeTo(check.pkg)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	for _, v := range stack {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		edge := check.mono.edges[check.mono.vertices[v].pre]
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		obj := check.mono.vertices[edge.dst].obj
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		switch obj.Type().(type) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		default:
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			panic(&#34;unexpected type&#34;)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		case *Named:
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			check.errorf(atPos(edge.pos), InvalidInstanceCycle, &#34;\t%s implicitly parameterized by %s&#34;, obj.Name(), TypeString(edge.typ, qf)) <span class="comment">// secondary error, \t indented</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		case *TypeParam:
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			check.errorf(atPos(edge.pos), InvalidInstanceCycle, &#34;\t%s instantiated as %s&#34;, obj.Name(), TypeString(edge.typ, qf)) <span class="comment">// secondary error, \t indented</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// recordCanon records that tpar is the canonical type parameter</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">// corresponding to method type parameter mpar.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>func (w *monoGraph) recordCanon(mpar, tpar *TypeParam) {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if w.canon == nil {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		w.canon = make(map[*TypeParam]*TypeParam)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	w.canon[mpar] = tpar
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// recordInstance records that the given type parameters were</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">// instantiated with the corresponding type arguments.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>func (w *monoGraph) recordInstance(pkg *Package, pos token.Pos, tparams []*TypeParam, targs []Type, xlist []ast.Expr) {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	for i, tpar := range tparams {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		pos := pos
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		if i &lt; len(xlist) {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			pos = xlist[i].Pos()
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		w.assign(pkg, pos, tpar, targs[i])
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">// assign records that tpar was instantiated as targ at pos.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>func (w *monoGraph) assign(pkg *Package, pos token.Pos, tpar *TypeParam, targ Type) {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// Go generics do not have an analog to C++`s template-templates,</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// where a template parameter can itself be an instantiable</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// template. So any instantiation cycles must occur within a single</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// package. Accordingly, we can ignore instantiations of imported</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// type parameters.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// TODO(mdempsky): Push this check up into recordInstance? All type</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// parameters in a list will appear in the same package.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	if tpar.Obj().Pkg() != pkg {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		return
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// flow adds an edge from vertex src representing that typ flows to tpar.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	flow := func(src int, typ Type) {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		weight := 1
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		if typ == targ {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			weight = 0
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		w.addEdge(w.typeParamVertex(tpar), src, weight, pos, targ)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// Recursively walk the type argument to find any defined types or</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// type parameters.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	var do func(typ Type)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	do = func(typ Type) {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		switch typ := Unalias(typ).(type) {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		default:
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			panic(&#34;unexpected type&#34;)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		case *TypeParam:
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			assert(typ.Obj().Pkg() == pkg)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			flow(w.typeParamVertex(typ), typ)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		case *Named:
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			if src := w.localNamedVertex(pkg, typ.Origin()); src &gt;= 0 {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>				flow(src, typ)
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			targs := typ.TypeArgs()
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			for i := 0; i &lt; targs.Len(); i++ {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				do(targs.At(i))
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		case *Array:
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			do(typ.Elem())
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		case *Basic:
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			<span class="comment">// ok</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		case *Chan:
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			do(typ.Elem())
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		case *Map:
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			do(typ.Key())
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			do(typ.Elem())
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		case *Pointer:
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			do(typ.Elem())
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		case *Slice:
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			do(typ.Elem())
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		case *Interface:
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			for i := 0; i &lt; typ.NumMethods(); i++ {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>				do(typ.Method(i).Type())
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		case *Signature:
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			tuple := func(tup *Tuple) {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>				for i := 0; i &lt; tup.Len(); i++ {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>					do(tup.At(i).Type())
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>				}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			tuple(typ.Params())
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			tuple(typ.Results())
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		case *Struct:
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			for i := 0; i &lt; typ.NumFields(); i++ {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				do(typ.Field(i).Type())
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	do(targ)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span><span class="comment">// localNamedVertex returns the index of the vertex representing</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">// named, or -1 if named doesn&#39;t need representation.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>func (w *monoGraph) localNamedVertex(pkg *Package, named *Named) int {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	obj := named.Obj()
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	if obj.Pkg() != pkg {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		return -1 <span class="comment">// imported type</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	}
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	root := pkg.Scope()
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	if obj.Parent() == root {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		return -1 <span class="comment">// package scope, no ambient type parameters</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	if idx, ok := w.nameIdx[obj]; ok {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		return idx
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	idx := -1
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	<span class="comment">// Walk the type definition&#39;s scope to find any ambient type</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	<span class="comment">// parameters that it&#39;s implicitly parameterized by.</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	for scope := obj.Parent(); scope != root; scope = scope.Parent() {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		for _, elem := range scope.elems {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			if elem, ok := elem.(*TypeName); ok &amp;&amp; !elem.IsAlias() &amp;&amp; cmpPos(elem.Pos(), obj.Pos()) &lt; 0 {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>				if tpar, ok := elem.Type().(*TypeParam); ok {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>					if idx &lt; 0 {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>						idx = len(w.vertices)
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>						w.vertices = append(w.vertices, monoVertex{obj: obj})
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>					}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>					w.addEdge(idx, w.typeParamVertex(tpar), 1, obj.Pos(), tpar)
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>				}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	if w.nameIdx == nil {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		w.nameIdx = make(map[*TypeName]int)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	w.nameIdx[obj] = idx
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	return idx
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">// typeParamVertex returns the index of the vertex representing tpar.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>func (w *monoGraph) typeParamVertex(tpar *TypeParam) int {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if x, ok := w.canon[tpar]; ok {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		tpar = x
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	obj := tpar.Obj()
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	if idx, ok := w.nameIdx[obj]; ok {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		return idx
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	if w.nameIdx == nil {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		w.nameIdx = make(map[*TypeName]int)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	idx := len(w.vertices)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	w.vertices = append(w.vertices, monoVertex{obj: obj})
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	w.nameIdx[obj] = idx
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	return idx
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>func (w *monoGraph) addEdge(dst, src, weight int, pos token.Pos, typ Type) {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	<span class="comment">// TODO(mdempsky): Deduplicate redundant edges?</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	w.edges = append(w.edges, monoEdge{
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		dst:    dst,
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		src:    src,
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		weight: weight,
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		pos: pos,
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		typ: typ,
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	})
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>
</pre><p><a href="mono.go?m=text">View as plain text</a></p>

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
