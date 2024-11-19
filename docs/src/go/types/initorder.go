<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/initorder.go - Go Documentation Server</title>

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
<a href="initorder.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">initorder.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;container/heap&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// initOrder computes the Info.InitOrder for package variables.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>func (check *Checker) initOrder() {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	<span class="comment">// An InitOrder may already have been computed if a package is</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">// built from several calls to (*Checker).Files. Clear it.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	check.Info.InitOrder = check.Info.InitOrder[:0]
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// Compute the object dependency graph and initialize</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// a priority queue with the list of graph nodes.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	pq := nodeQueue(dependencyGraph(check.objMap))
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	heap.Init(&amp;pq)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	const debug = false
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	if debug {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		fmt.Printf(&#34;Computing initialization order for %s\n\n&#34;, check.pkg)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		fmt.Println(&#34;Object dependency graph:&#34;)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		for obj, d := range check.objMap {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>			<span class="comment">// only print objects that may appear in the dependency graph</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>			if obj, _ := obj.(dependency); obj != nil {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>				if len(d.deps) &gt; 0 {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>					fmt.Printf(&#34;\t%s depends on\n&#34;, obj.Name())
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>					for dep := range d.deps {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>						fmt.Printf(&#34;\t\t%s\n&#34;, dep.Name())
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>					}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>				} else {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>					fmt.Printf(&#34;\t%s has no dependencies\n&#34;, obj.Name())
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>				}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>			}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		fmt.Println()
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		fmt.Println(&#34;Transposed object dependency graph (functions eliminated):&#34;)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		for _, n := range pq {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			fmt.Printf(&#34;\t%s depends on %d nodes\n&#34;, n.obj.Name(), n.ndeps)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			for p := range n.pred {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>				fmt.Printf(&#34;\t\t%s is dependent\n&#34;, p.obj.Name())
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		fmt.Println()
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		fmt.Println(&#34;Processing nodes:&#34;)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// Determine initialization order by removing the highest priority node</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// (the one with the fewest dependencies) and its edges from the graph,</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// repeatedly, until there are no nodes left.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// In a valid Go program, those nodes always have zero dependencies (after</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// removing all incoming dependencies), otherwise there are initialization</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// cycles.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	emitted := make(map[*declInfo]bool)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	for len(pq) &gt; 0 {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		<span class="comment">// get the next node</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		n := heap.Pop(&amp;pq).(*graphNode)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		if debug {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			fmt.Printf(&#34;\t%s (src pos %d) depends on %d nodes now\n&#34;,
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>				n.obj.Name(), n.obj.order(), n.ndeps)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		<span class="comment">// if n still depends on other nodes, we have a cycle</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		if n.ndeps &gt; 0 {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			cycle := findPath(check.objMap, n.obj, n.obj, make(map[Object]bool))
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			<span class="comment">// If n.obj is not part of the cycle (e.g., n.obj-&gt;b-&gt;c-&gt;d-&gt;c),</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			<span class="comment">// cycle will be nil. Don&#39;t report anything in that case since</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			<span class="comment">// the cycle is reported when the algorithm gets to an object</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			<span class="comment">// in the cycle.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			<span class="comment">// Furthermore, once an object in the cycle is encountered,</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			<span class="comment">// the cycle will be broken (dependency count will be reduced</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			<span class="comment">// below), and so the remaining nodes in the cycle don&#39;t trigger</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			<span class="comment">// another error (unless they are part of multiple cycles).</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			if cycle != nil {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>				check.reportCycle(cycle)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			<span class="comment">// Ok to continue, but the variable initialization order</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			<span class="comment">// will be incorrect at this point since it assumes no</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			<span class="comment">// cycle errors.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		<span class="comment">// reduce dependency count of all dependent nodes</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		<span class="comment">// and update priority queue</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		for p := range n.pred {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			p.ndeps--
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			heap.Fix(&amp;pq, p.index)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		<span class="comment">// record the init order for variables with initializers only</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		v, _ := n.obj.(*Var)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		info := check.objMap[v]
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		if v == nil || !info.hasInitializer() {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			continue
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		<span class="comment">// n:1 variable declarations such as: a, b = f()</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">// introduce a node for each lhs variable (here: a, b);</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		<span class="comment">// but they all have the same initializer - emit only</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		<span class="comment">// one, for the first variable seen</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		if emitted[info] {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			continue <span class="comment">// initializer already emitted, if any</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		emitted[info] = true
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		infoLhs := info.lhs <span class="comment">// possibly nil (see declInfo.lhs field comment)</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		if infoLhs == nil {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			infoLhs = []*Var{v}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		init := &amp;Initializer{infoLhs, info.init}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		check.Info.InitOrder = append(check.Info.InitOrder, init)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	if debug {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		fmt.Println()
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		fmt.Println(&#34;Initialization order:&#34;)
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		for _, init := range check.Info.InitOrder {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			fmt.Printf(&#34;\t%s\n&#34;, init)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		fmt.Println()
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// findPath returns the (reversed) list of objects []Object{to, ... from}</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">// such that there is a path of object dependencies from &#39;from&#39; to &#39;to&#39;.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">// If there is no such path, the result is nil.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>func findPath(objMap map[Object]*declInfo, from, to Object, seen map[Object]bool) []Object {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	if seen[from] {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		return nil
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	seen[from] = true
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	for d := range objMap[from].deps {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		if d == to {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			return []Object{d}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		if P := findPath(objMap, d, to, seen); P != nil {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			return append(P, d)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	return nil
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">// reportCycle reports an error for the given cycle.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>func (check *Checker) reportCycle(cycle []Object) {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	obj := cycle[0]
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// report a more concise error for self references</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	if len(cycle) == 1 {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		check.errorf(obj, InvalidInitCycle, &#34;initialization cycle: %s refers to itself&#34;, obj.Name())
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		return
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	check.errorf(obj, InvalidInitCycle, &#34;initialization cycle for %s&#34;, obj.Name())
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// subtle loop: print cycle[i] for i = 0, n-1, n-2, ... 1 for len(cycle) = n</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	for i := len(cycle) - 1; i &gt;= 0; i-- {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		check.errorf(obj, InvalidInitCycle, &#34;\t%s refers to&#34;, obj.Name()) <span class="comment">// secondary error, \t indented</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		obj = cycle[i]
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// print cycle[0] again to close the cycle</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	check.errorf(obj, InvalidInitCycle, &#34;\t%s&#34;, obj.Name())
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">// Object dependency graph</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// A dependency is an object that may be a dependency in an initialization</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">// expression. Only constants, variables, and functions can be dependencies.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">// Constants are here because constant expression cycles are reported during</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">// initialization order computation.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>type dependency interface {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	Object
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	isDependency()
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span><span class="comment">// A graphNode represents a node in the object dependency graph.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// Each node p in n.pred represents an edge p-&gt;n, and each node</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// s in n.succ represents an edge n-&gt;s; with a-&gt;b indicating that</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// a depends on b.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>type graphNode struct {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	obj        dependency <span class="comment">// object represented by this node</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	pred, succ nodeSet    <span class="comment">// consumers and dependencies of this node (lazily initialized)</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	index      int        <span class="comment">// node index in graph slice/priority queue</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	ndeps      int        <span class="comment">// number of outstanding dependencies before this object can be initialized</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// cost returns the cost of removing this node, which involves copying each</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// predecessor to each successor (and vice-versa).</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>func (n *graphNode) cost() int {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	return len(n.pred) * len(n.succ)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>type nodeSet map[*graphNode]bool
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>func (s *nodeSet) add(p *graphNode) {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	if *s == nil {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		*s = make(nodeSet)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	(*s)[p] = true
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// dependencyGraph computes the object dependency graph from the given objMap,</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// with any function nodes removed. The resulting graph contains only constants</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// and variables.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>func dependencyGraph(objMap map[Object]*declInfo) []*graphNode {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// M is the dependency (Object) -&gt; graphNode mapping</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	M := make(map[dependency]*graphNode)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	for obj := range objMap {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// only consider nodes that may be an initialization dependency</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		if obj, _ := obj.(dependency); obj != nil {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			M[obj] = &amp;graphNode{obj: obj}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">// compute edges for graph M</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	<span class="comment">// (We need to include all nodes, even isolated ones, because they still need</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// to be scheduled for initialization in correct order relative to other nodes.)</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	for obj, n := range M {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		<span class="comment">// for each dependency obj -&gt; d (= deps[i]), create graph edges n-&gt;s and s-&gt;n</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		for d := range objMap[obj].deps {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			<span class="comment">// only consider nodes that may be an initialization dependency</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			if d, _ := d.(dependency); d != nil {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>				d := M[d]
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>				n.succ.add(d)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>				d.pred.add(n)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	var G, funcG []*graphNode <span class="comment">// separate non-functions and functions</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	for _, n := range M {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		if _, ok := n.obj.(*Func); ok {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			funcG = append(funcG, n)
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		} else {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			G = append(G, n)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// remove function nodes and collect remaining graph nodes in G</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	<span class="comment">// (Mutually recursive functions may introduce cycles among themselves</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// which are permitted. Yet such cycles may incorrectly inflate the dependency</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// count for variables which in turn may not get scheduled for initialization</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// in correct order.)</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	<span class="comment">// Note that because we recursively copy predecessors and successors</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	<span class="comment">// throughout the function graph, the cost of removing a function at</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// position X is proportional to cost * (len(funcG)-X). Therefore, we should</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">// remove high-cost functions last.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	sort.Slice(funcG, func(i, j int) bool {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		return funcG[i].cost() &lt; funcG[j].cost()
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	})
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	for _, n := range funcG {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		<span class="comment">// connect each predecessor p of n with each successor s</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		<span class="comment">// and drop the function node (don&#39;t collect it in G)</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		for p := range n.pred {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			<span class="comment">// ignore self-cycles</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			if p != n {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>				<span class="comment">// Each successor s of n becomes a successor of p, and</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				<span class="comment">// each predecessor p of n becomes a predecessor of s.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>				for s := range n.succ {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>					<span class="comment">// ignore self-cycles</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>					if s != n {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>						p.succ.add(s)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>						s.pred.add(p)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>					}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>				}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>				delete(p.succ, n) <span class="comment">// remove edge to n</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		for s := range n.succ {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			delete(s.pred, n) <span class="comment">// remove edge to n</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	<span class="comment">// fill in index and ndeps fields</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	for i, n := range G {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		n.index = i
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		n.ndeps = len(n.succ)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	return G
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span><span class="comment">// Priority queue</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span><span class="comment">// nodeQueue implements the container/heap interface;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span><span class="comment">// a nodeQueue may be used as a priority queue.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>type nodeQueue []*graphNode
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>func (a nodeQueue) Len() int { return len(a) }
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>func (a nodeQueue) Swap(i, j int) {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	x, y := a[i], a[j]
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	a[i], a[j] = y, x
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	x.index, y.index = j, i
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>func (a nodeQueue) Less(i, j int) bool {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	x, y := a[i], a[j]
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	<span class="comment">// nodes are prioritized by number of incoming dependencies (1st key)</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	<span class="comment">// and source order (2nd key)</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	return x.ndeps &lt; y.ndeps || x.ndeps == y.ndeps &amp;&amp; x.obj.order() &lt; y.obj.order()
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>func (a *nodeQueue) Push(x any) {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	panic(&#34;unreachable&#34;)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>func (a *nodeQueue) Pop() any {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	n := len(*a)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	x := (*a)[n-1]
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	x.index = -1 <span class="comment">// for safety</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	*a = (*a)[:n-1]
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	return x
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
</pre><p><a href="initorder.go?m=text">View as plain text</a></p>

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
