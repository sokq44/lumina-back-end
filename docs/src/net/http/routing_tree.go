<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/http/routing_tree.go - Go Documentation Server</title>

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
<a href="routing_tree.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<a href="http://localhost:8080/src/net/http">http</a>/<span class="text-muted">routing_tree.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net/http">net/http</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2023 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements a decision tree for fast matching of requests to</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// patterns.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// The root of the tree branches on the host of the request.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// The next level branches on the method.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// The remaining levels branch on consecutive segments of the path.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// The &#34;more specific wins&#34; precedence rule can result in backtracking.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// For example, given the patterns</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//     /a/b/z</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//     /a/{x}/c</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// we will first try to match the path &#34;/a/b/c&#34; with /a/b/z, and</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// when that fails we will try against /a/{x}/c.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>package http
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>import (
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// A routingNode is a node in the decision tree.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// The same struct is used for leaf and interior nodes.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>type routingNode struct {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// A leaf node holds a single pattern and the Handler it was registered</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// with.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	pattern *pattern
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	handler Handler
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// An interior node maps parts of the incoming request to child nodes.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// special children keys:</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">//     &#34;/&#34;	trailing slash (resulting from {$})</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">//	   &#34;&#34;   single wildcard</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">//	   &#34;*&#34;  multi wildcard</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	children   mapping[string, *routingNode]
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	emptyChild *routingNode <span class="comment">// optimization: child with key &#34;&#34;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// addPattern adds a pattern and its associated Handler to the tree</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// at root.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>func (root *routingNode) addPattern(p *pattern, h Handler) {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// First level of tree is host.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	n := root.addChild(p.host)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// Second level of tree is method.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	n = n.addChild(p.method)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// Remaining levels are path.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	n.addSegments(p.segments, p, h)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// addSegments adds the given segments to the tree rooted at n.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// If there are no segments, then n is a leaf node that holds</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// the given pattern and handler.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>func (n *routingNode) addSegments(segs []segment, p *pattern, h Handler) {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	if len(segs) == 0 {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		n.set(p, h)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		return
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	seg := segs[0]
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	if seg.multi {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		if len(segs) != 1 {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			panic(&#34;multi wildcard not last&#34;)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		n.addChild(&#34;*&#34;).set(p, h)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	} else if seg.wild {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		n.addChild(&#34;&#34;).addSegments(segs[1:], p, h)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	} else {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		n.addChild(seg.s).addSegments(segs[1:], p, h)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// set sets the pattern and handler for n, which</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// must be a leaf node.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>func (n *routingNode) set(p *pattern, h Handler) {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	if n.pattern != nil || n.handler != nil {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		panic(&#34;non-nil leaf fields&#34;)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	n.pattern = p
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	n.handler = h
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// addChild adds a child node with the given key to n</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// if one does not exist, and returns the child.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>func (n *routingNode) addChild(key string) *routingNode {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	if key == &#34;&#34; {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		if n.emptyChild == nil {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			n.emptyChild = &amp;routingNode{}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		return n.emptyChild
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	if c := n.findChild(key); c != nil {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		return c
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	c := &amp;routingNode{}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	n.children.add(key, c)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	return c
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// findChild returns the child of n with the given key, or nil</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// if there is no child with that key.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>func (n *routingNode) findChild(key string) *routingNode {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	if key == &#34;&#34; {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		return n.emptyChild
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	r, _ := n.children.find(key)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	return r
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// match returns the leaf node under root that matches the arguments, and a list</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// of values for pattern wildcards in the order that the wildcards appear.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// For example, if the request path is &#34;/a/b/c&#34; and the pattern is &#34;/{x}/b/{y}&#34;,</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">// then the second return value will be []string{&#34;a&#34;, &#34;c&#34;}.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>func (root *routingNode) match(host, method, path string) (*routingNode, []string) {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	if host != &#34;&#34; {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		<span class="comment">// There is a host. If there is a pattern that specifies that host and it</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		<span class="comment">// matches, we are done. If the pattern doesn&#39;t match, fall through to</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		<span class="comment">// try patterns with no host.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		if l, m := root.findChild(host).matchMethodAndPath(method, path); l != nil {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			return l, m
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	return root.emptyChild.matchMethodAndPath(method, path)
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// matchMethodAndPath matches the method and path.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// Its return values are the same as [routingNode.match].</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// The receiver should be a child of the root.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>func (n *routingNode) matchMethodAndPath(method, path string) (*routingNode, []string) {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	if n == nil {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		return nil, nil
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	if l, m := n.findChild(method).matchPath(path, nil); l != nil {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// Exact match of method name.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		return l, m
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	if method == &#34;HEAD&#34; {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		<span class="comment">// GET matches HEAD too.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		if l, m := n.findChild(&#34;GET&#34;).matchPath(path, nil); l != nil {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			return l, m
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// No exact match; try patterns with no method.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	return n.emptyChild.matchPath(path, nil)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// matchPath matches a path.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// Its return values are the same as [routingNode.match].</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// matchPath calls itself recursively. The matches argument holds the wildcard matches</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// found so far.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func (n *routingNode) matchPath(path string, matches []string) (*routingNode, []string) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if n == nil {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return nil, nil
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// If path is empty, then we are done.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// If n is a leaf node, we found a match; return it.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// If n is an interior node (which means it has a nil pattern),</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// then we failed to match.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	if path == &#34;&#34; {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		if n.pattern == nil {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			return nil, nil
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		return n, matches
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// Get the first segment of path.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	seg, rest := firstSegment(path)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// First try matching against patterns that have a literal for this position.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// We know by construction that such patterns are more specific than those</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// with a wildcard at this position (they are either more specific, equivalent,</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// or overlap, and we ruled out the first two when the patterns were registered).</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	if n, m := n.findChild(seg).matchPath(rest, matches); n != nil {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		return n, m
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// If matching a literal fails, try again with patterns that have a single</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// wildcard (represented by an empty string in the child mapping).</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// Again, by construction, patterns with a single wildcard must be more specific than</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// those with a multi wildcard.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// We skip this step if the segment is a trailing slash, because single wildcards</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// don&#39;t match trailing slashes.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	if seg != &#34;/&#34; {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		if n, m := n.emptyChild.matchPath(rest, append(matches, seg)); n != nil {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			return n, m
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// Lastly, match the pattern (there can be at most one) that has a multi</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// wildcard in this position to the rest of the path.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	if c := n.findChild(&#34;*&#34;); c != nil {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		<span class="comment">// Don&#39;t record a match for a nameless wildcard (which arises from a</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		<span class="comment">// trailing slash in the pattern).</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		if c.pattern.lastSegment().s != &#34;&#34; {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			matches = append(matches, pathUnescape(path[1:])) <span class="comment">// remove initial slash</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		return c, matches
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	return nil, nil
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// firstSegment splits path into its first segment, and the rest.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// The path must begin with &#34;/&#34;.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// If path consists of only a slash, firstSegment returns (&#34;/&#34;, &#34;&#34;).</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">// The segment is returned unescaped, if possible.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>func firstSegment(path string) (seg, rest string) {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	if path == &#34;/&#34; {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		return &#34;/&#34;, &#34;&#34;
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	path = path[1:] <span class="comment">// drop initial slash</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	i := strings.IndexByte(path, &#39;/&#39;)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		i = len(path)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	return pathUnescape(path[:i]), path[i:]
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// matchingMethods adds to methodSet all the methods that would result in a</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// match if passed to routingNode.match with the given host and path.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>func (root *routingNode) matchingMethods(host, path string, methodSet map[string]bool) {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	if host != &#34;&#34; {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		root.findChild(host).matchingMethodsPath(path, methodSet)
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	root.emptyChild.matchingMethodsPath(path, methodSet)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	if methodSet[&#34;GET&#34;] {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		methodSet[&#34;HEAD&#34;] = true
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>func (n *routingNode) matchingMethodsPath(path string, set map[string]bool) {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	if n == nil {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		return
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	n.children.eachPair(func(method string, c *routingNode) bool {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		if p, _ := c.matchPath(path, nil); p != nil {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			set[method] = true
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		return true
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	})
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t look at the empty child. If there were an empty</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// child, it would match on any method, but we only</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	<span class="comment">// call this when we fail to match on a method.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>
</pre><p><a href="routing_tree.go?m=text">View as plain text</a></p>

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
