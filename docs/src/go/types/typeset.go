<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/typeset.go - Go Documentation Server</title>

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
<a href="typeset.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">typeset.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// API</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// A _TypeSet represents the type set of an interface.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// Because of existing language restrictions, methods can be &#34;factored out&#34;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// from the terms. The actual type set is the intersection of the type set</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// implied by the methods and the type set described by the terms and the</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// comparable bit. To test whether a type is included in a type set</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// (&#34;implements&#34; relation), the type must implement all methods _and_ be</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// an element of the type set described by the terms and the comparable bit.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// If the term list describes the set of all types and comparable is true,</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// only comparable types are meant; in all other cases comparable is false.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>type _TypeSet struct {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	methods    []*Func  <span class="comment">// all methods of the interface; sorted by unique ID</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	terms      termlist <span class="comment">// type terms of the type set</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	comparable bool     <span class="comment">// invariant: !comparable || terms.isAll()</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// IsEmpty reports whether type set s is the empty set.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>func (s *_TypeSet) IsEmpty() bool { return s.terms.isEmpty() }
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// IsAll reports whether type set s is the set of all types (corresponding to the empty interface).</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>func (s *_TypeSet) IsAll() bool { return s.IsMethodSet() &amp;&amp; len(s.methods) == 0 }
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// IsMethodSet reports whether the interface t is fully described by its method set.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>func (s *_TypeSet) IsMethodSet() bool { return !s.comparable &amp;&amp; s.terms.isAll() }
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// IsComparable reports whether each type in the set is comparable.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func (s *_TypeSet) IsComparable(seen map[Type]bool) bool {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	if s.terms.isAll() {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		return s.comparable
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	return s.is(func(t *term) bool {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		return t != nil &amp;&amp; comparable(t.typ, false, seen, nil)
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	})
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// NumMethods returns the number of methods available.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>func (s *_TypeSet) NumMethods() int { return len(s.methods) }
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// Method returns the i&#39;th method of type set s for 0 &lt;= i &lt; s.NumMethods().</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// The methods are ordered by their unique ID.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>func (s *_TypeSet) Method(i int) *Func { return s.methods[i] }
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// LookupMethod returns the index of and method with matching package and name, or (-1, nil).</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>func (s *_TypeSet) LookupMethod(pkg *Package, name string, foldCase bool) (int, *Func) {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	return lookupMethod(s.methods, pkg, name, foldCase)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>func (s *_TypeSet) String() string {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	switch {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	case s.IsEmpty():
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		return &#34;∅&#34;
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	case s.IsAll():
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		return &#34;𝓤&#34;
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	hasMethods := len(s.methods) &gt; 0
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	hasTerms := s.hasTerms()
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	var buf strings.Builder
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	buf.WriteByte(&#39;{&#39;)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	if s.comparable {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		buf.WriteString(&#34;comparable&#34;)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		if hasMethods || hasTerms {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			buf.WriteString(&#34;; &#34;)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	for i, m := range s.methods {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		if i &gt; 0 {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			buf.WriteString(&#34;; &#34;)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		buf.WriteString(m.String())
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	if hasMethods &amp;&amp; hasTerms {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		buf.WriteString(&#34;; &#34;)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	if hasTerms {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		buf.WriteString(s.terms.String())
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	buf.WriteString(&#34;}&#34;)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	return buf.String()
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// Implementation</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// hasTerms reports whether the type set has specific type terms.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>func (s *_TypeSet) hasTerms() bool { return !s.terms.isEmpty() &amp;&amp; !s.terms.isAll() }
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// subsetOf reports whether s1 ⊆ s2.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func (s1 *_TypeSet) subsetOf(s2 *_TypeSet) bool { return s1.terms.subsetOf(s2.terms) }
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// TODO(gri) TypeSet.is and TypeSet.underIs should probably also go into termlist.go</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// is calls f with the specific type terms of s and reports whether</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// all calls to f returned true. If there are no specific terms, is</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// returns the result of f(nil).</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>func (s *_TypeSet) is(f func(*term) bool) bool {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	if !s.hasTerms() {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		return f(nil)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	for _, t := range s.terms {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		assert(t.typ != nil)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		if !f(t) {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			return false
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	return true
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// underIs calls f with the underlying types of the specific type terms</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// of s and reports whether all calls to f returned true. If there are</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// no specific terms, underIs returns the result of f(nil).</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>func (s *_TypeSet) underIs(f func(Type) bool) bool {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if !s.hasTerms() {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		return f(nil)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	for _, t := range s.terms {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		assert(t.typ != nil)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// x == under(x) for ~x terms</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		u := t.typ
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		if !t.tilde {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			u = under(u)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		if debug {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			assert(Identical(u, under(u)))
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		if !f(u) {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			return false
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	return true
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// topTypeSet may be used as type set for the empty interface.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>var topTypeSet = _TypeSet{terms: allTermlist}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">// computeInterfaceTypeSet may be called with check == nil.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>func computeInterfaceTypeSet(check *Checker, pos token.Pos, ityp *Interface) *_TypeSet {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	if ityp.tset != nil {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		return ityp.tset
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// If the interface is not fully set up yet, the type set will</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// not be complete, which may lead to errors when using the</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// type set (e.g. missing method). Don&#39;t compute a partial type</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// set (and don&#39;t store it!), so that we still compute the full</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// type set eventually. Instead, return the top type set and</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// let any follow-on errors play out.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// TODO(gri) Consider recording when this happens and reporting</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// it as an error (but only if there were no other errors so to</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// to not have unnecessary follow-on errors).</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	if !ityp.complete {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		return &amp;topTypeSet
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	if check != nil &amp;&amp; check.conf._Trace {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		<span class="comment">// Types don&#39;t generally have position information.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		<span class="comment">// If we don&#39;t have a valid pos provided, try to use</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		<span class="comment">// one close enough.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		if !pos.IsValid() &amp;&amp; len(ityp.methods) &gt; 0 {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			pos = ityp.methods[0].pos
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		check.trace(pos, &#34;-- type set for %s&#34;, ityp)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		check.indent++
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		defer func() {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			check.indent--
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			check.trace(pos, &#34;=&gt; %s &#34;, ityp.typeSet())
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		}()
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">// An infinitely expanding interface (due to a cycle) is detected</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// elsewhere (Checker.validType), so here we simply assume we only</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// have valid interfaces. Mark the interface as complete to avoid</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// infinite recursion if the validType check occurs later for some</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// reason.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	ityp.tset = &amp;_TypeSet{terms: allTermlist} <span class="comment">// TODO(gri) is this sufficient?</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	var unionSets map[*Union]*_TypeSet
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	if check != nil {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		if check.unionTypeSets == nil {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			check.unionTypeSets = make(map[*Union]*_TypeSet)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		unionSets = check.unionTypeSets
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	} else {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		unionSets = make(map[*Union]*_TypeSet)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// Methods of embedded interfaces are collected unchanged; i.e., the identity</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// of a method I.m&#39;s Func Object of an interface I is the same as that of</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// the method m in an interface that embeds interface I. On the other hand,</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// if a method is embedded via multiple overlapping embedded interfaces, we</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">// don&#39;t provide a guarantee which &#34;original m&#34; got chosen for the embedding</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">// interface. See also go.dev/issue/34421.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// If we don&#39;t care to provide this identity guarantee anymore, instead of</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// reusing the original method in embeddings, we can clone the method&#39;s Func</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// Object and give it the position of a corresponding embedded interface. Then</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// we can get rid of the mpos map below and simply use the cloned method&#39;s</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	<span class="comment">// position.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	var seen objset
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	var allMethods []*Func
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	mpos := make(map[*Func]token.Pos) <span class="comment">// method specification or method embedding position, for good error messages</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	addMethod := func(pos token.Pos, m *Func, explicit bool) {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		switch other := seen.insert(m); {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		case other == nil:
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			allMethods = append(allMethods, m)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			mpos[m] = pos
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		case explicit:
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			if check != nil {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>				check.errorf(atPos(pos), DuplicateDecl, &#34;duplicate method %s&#34;, m.name)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>				check.errorf(atPos(mpos[other.(*Func)]), DuplicateDecl, &#34;\tother declaration of %s&#34;, m.name) <span class="comment">// secondary error, \t indented</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		default:
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			<span class="comment">// We have a duplicate method name in an embedded (not explicitly declared) method.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			<span class="comment">// Check method signatures after all types are computed (go.dev/issue/33656).</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			<span class="comment">// If we&#39;re pre-go1.14 (overlapping embeddings are not permitted), report that</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			<span class="comment">// error here as well (even though we could do it eagerly) because it&#39;s the same</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			<span class="comment">// error message.</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			if check != nil {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>				check.later(func() {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>					if !check.allowVersion(m.pkg, atPos(pos), go1_14) || !Identical(m.typ, other.Type()) {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>						check.errorf(atPos(pos), DuplicateDecl, &#34;duplicate method %s&#34;, m.name)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>						check.errorf(atPos(mpos[other.(*Func)]), DuplicateDecl, &#34;\tother declaration of %s&#34;, m.name) <span class="comment">// secondary error, \t indented</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>					}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>				}).describef(atPos(pos), &#34;duplicate method check for %s&#34;, m.name)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	for _, m := range ityp.methods {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		addMethod(m.pos, m, true)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// collect embedded elements</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	allTerms := allTermlist
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	allComparable := false
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	for i, typ := range ityp.embeddeds {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		<span class="comment">// The embedding position is nil for imported interfaces</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		<span class="comment">// and also for interface copies after substitution (but</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		<span class="comment">// in that case we don&#39;t need to report errors again).</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		var pos token.Pos <span class="comment">// embedding position</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		if ityp.embedPos != nil {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			pos = (*ityp.embedPos)[i]
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		var comparable bool
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		var terms termlist
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		switch u := under(typ).(type) {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		case *Interface:
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			<span class="comment">// For now we don&#39;t permit type parameters as constraints.</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			assert(!isTypeParam(typ))
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			tset := computeInterfaceTypeSet(check, pos, u)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			<span class="comment">// If typ is local, an error was already reported where typ is specified/defined.</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			if check != nil &amp;&amp; check.isImportedConstraint(typ) &amp;&amp; !check.verifyVersionf(atPos(pos), go1_18, &#34;embedding constraint interface %s&#34;, typ) {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>				continue
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			comparable = tset.comparable
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			for _, m := range tset.methods {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>				addMethod(pos, m, false) <span class="comment">// use embedding position pos rather than m.pos</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			terms = tset.terms
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		case *Union:
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			if check != nil &amp;&amp; !check.verifyVersionf(atPos(pos), go1_18, &#34;embedding interface element %s&#34;, u) {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>				continue
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			tset := computeUnionTypeSet(check, unionSets, pos, u)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			if tset == &amp;invalidTypeSet {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>				continue <span class="comment">// ignore invalid unions</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			assert(!tset.comparable)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			assert(len(tset.methods) == 0)
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			terms = tset.terms
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		default:
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			if !isValid(u) {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>				continue
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			if check != nil &amp;&amp; !check.verifyVersionf(atPos(pos), go1_18, &#34;embedding non-interface type %s&#34;, typ) {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>				continue
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			terms = termlist{{false, typ}}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		<span class="comment">// The type set of an interface is the intersection of the type sets of all its elements.</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		<span class="comment">// Due to language restrictions, only embedded interfaces can add methods, they are handled</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		<span class="comment">// separately. Here we only need to intersect the term lists and comparable bits.</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		allTerms, allComparable = intersectTermLists(allTerms, allComparable, terms, comparable)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	ityp.tset.comparable = allComparable
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if len(allMethods) != 0 {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		sortMethods(allMethods)
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		ityp.tset.methods = allMethods
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	}
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	ityp.tset.terms = allTerms
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	return ityp.tset
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">// TODO(gri) The intersectTermLists function belongs to the termlist implementation.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">//           The comparable type set may also be best represented as a term (using</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">//           a special type).</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span><span class="comment">// intersectTermLists computes the intersection of two term lists and respective comparable bits.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span><span class="comment">// xcomp, ycomp are valid only if xterms.isAll() and yterms.isAll() respectively.</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>func intersectTermLists(xterms termlist, xcomp bool, yterms termlist, ycomp bool) (termlist, bool) {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	terms := xterms.intersect(yterms)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	<span class="comment">// If one of xterms or yterms is marked as comparable,</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// the result must only include comparable types.</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	comp := xcomp || ycomp
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	if comp &amp;&amp; !terms.isAll() {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		<span class="comment">// only keep comparable terms</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		i := 0
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		for _, t := range terms {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			assert(t.typ != nil)
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			if comparable(t.typ, false <span class="comment">/* strictly comparable */</span>, nil, nil) {
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>				terms[i] = t
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>				i++
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>			}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		terms = terms[:i]
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		if !terms.isAll() {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>			comp = false
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	assert(!comp || terms.isAll()) <span class="comment">// comparable invariant</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	return terms, comp
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>func sortMethods(list []*Func) {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	sort.Sort(byUniqueMethodName(list))
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>func assertSortedMethods(list []*Func) {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	if !debug {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		panic(&#34;assertSortedMethods called outside debug mode&#34;)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	if !sort.IsSorted(byUniqueMethodName(list)) {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		panic(&#34;methods not sorted&#34;)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// byUniqueMethodName method lists can be sorted by their unique method names.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>type byUniqueMethodName []*Func
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>func (a byUniqueMethodName) Len() int           { return len(a) }
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>func (a byUniqueMethodName) Less(i, j int) bool { return a[i].less(&amp;a[j].object) }
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>func (a byUniqueMethodName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// invalidTypeSet is a singleton type set to signal an invalid type set</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">// due to an error. It&#39;s also a valid empty type set, so consumers of</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">// type sets may choose to ignore it.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>var invalidTypeSet _TypeSet
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span><span class="comment">// computeUnionTypeSet may be called with check == nil.</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span><span class="comment">// The result is &amp;invalidTypeSet if the union overflows.</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>func computeUnionTypeSet(check *Checker, unionSets map[*Union]*_TypeSet, pos token.Pos, utyp *Union) *_TypeSet {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	if tset, _ := unionSets[utyp]; tset != nil {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		return tset
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	<span class="comment">// avoid infinite recursion (see also computeInterfaceTypeSet)</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	unionSets[utyp] = new(_TypeSet)
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	var allTerms termlist
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	for _, t := range utyp.terms {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		var terms termlist
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		u := under(t.typ)
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		if ui, _ := u.(*Interface); ui != nil {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			<span class="comment">// For now we don&#39;t permit type parameters as constraints.</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>			assert(!isTypeParam(t.typ))
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			terms = computeInterfaceTypeSet(check, pos, ui).terms
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		} else if !isValid(u) {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			continue
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		} else {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			if t.tilde &amp;&amp; !Identical(t.typ, u) {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>				<span class="comment">// There is no underlying type which is t.typ.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>				<span class="comment">// The corresponding type set is empty.</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>				t = nil <span class="comment">// ∅ term</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			terms = termlist{(*term)(t)}
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		<span class="comment">// The type set of a union expression is the union</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		<span class="comment">// of the type sets of each term.</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		allTerms = allTerms.union(terms)
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		if len(allTerms) &gt; maxTermCount {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			if check != nil {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>				check.errorf(atPos(pos), InvalidUnion, &#34;cannot handle more than %d union terms (implementation limitation)&#34;, maxTermCount)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			unionSets[utyp] = &amp;invalidTypeSet
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			return unionSets[utyp]
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	unionSets[utyp].terms = allTerms
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	return unionSets[utyp]
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>
</pre><p><a href="typeset.go?m=text">View as plain text</a></p>

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
