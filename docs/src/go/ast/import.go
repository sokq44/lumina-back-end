<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/ast/import.go - Go Documentation Server</title>

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
<a href="import.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/ast">ast</a>/<span class="text-muted">import.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/ast">go/ast</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package ast
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// SortImports sorts runs of consecutive import lines in import blocks in f.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// It also removes duplicate imports when it is possible to do so without data loss.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>func SortImports(fset *token.FileSet, f *File) {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	for _, d := range f.Decls {
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>		d, ok := d.(*GenDecl)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>		if !ok || d.Tok != token.IMPORT {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>			<span class="comment">// Not an import declaration, so we&#39;re done.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>			<span class="comment">// Imports are always first.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>			break
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>		}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>		if !d.Lparen.IsValid() {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>			<span class="comment">// Not a block: sorted by default.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>			continue
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		<span class="comment">// Identify and sort runs of specs on successive lines.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		i := 0
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		specs := d.Specs[:0]
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		for j, s := range d.Specs {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>			if j &gt; i &amp;&amp; lineAt(fset, s.Pos()) &gt; 1+lineAt(fset, d.Specs[j-1].End()) {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>				<span class="comment">// j begins a new run. End this one.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>				specs = append(specs, sortSpecs(fset, f, d.Specs[i:j])...)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>				i = j
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>			}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		specs = append(specs, sortSpecs(fset, f, d.Specs[i:])...)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		d.Specs = specs
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		<span class="comment">// Deduping can leave a blank line before the rparen; clean that up.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		if len(d.Specs) &gt; 0 {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>			lastSpec := d.Specs[len(d.Specs)-1]
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>			lastLine := lineAt(fset, lastSpec.Pos())
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			rParenLine := lineAt(fset, d.Rparen)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			for rParenLine &gt; lastLine+1 {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>				rParenLine--
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>				fset.File(d.Rparen).MergeLine(rParenLine)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>func lineAt(fset *token.FileSet, pos token.Pos) int {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	return fset.PositionFor(pos, false).Line
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>func importPath(s Spec) string {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	t, err := strconv.Unquote(s.(*ImportSpec).Path.Value)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	if err == nil {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		return t
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	return &#34;&#34;
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>func importName(s Spec) string {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	n := s.(*ImportSpec).Name
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	if n == nil {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	return n.Name
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>func importComment(s Spec) string {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	c := s.(*ImportSpec).Comment
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	if c == nil {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	return c.Text()
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// collapse indicates whether prev may be removed, leaving only next.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>func collapse(prev, next Spec) bool {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	if importPath(next) != importPath(prev) || importName(next) != importName(prev) {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		return false
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	return prev.(*ImportSpec).Comment == nil
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>type posSpan struct {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	Start token.Pos
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	End   token.Pos
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>type cgPos struct {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	left bool <span class="comment">// true if comment is to the left of the spec, false otherwise.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	cg   *CommentGroup
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>func sortSpecs(fset *token.FileSet, f *File, specs []Spec) []Spec {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// Can&#39;t short-circuit here even if specs are already sorted,</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// since they might yet need deduplication.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// A lone import, however, may be safely ignored.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	if len(specs) &lt;= 1 {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		return specs
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// Record positions for specs.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	pos := make([]posSpan, len(specs))
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	for i, s := range specs {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		pos[i] = posSpan{s.Pos(), s.End()}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// Identify comments in this range.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	begSpecs := pos[0].Start
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	endSpecs := pos[len(pos)-1].End
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	beg := fset.File(begSpecs).LineStart(lineAt(fset, begSpecs))
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	endLine := lineAt(fset, endSpecs)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	endFile := fset.File(endSpecs)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	var end token.Pos
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	if endLine == endFile.LineCount() {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		end = endSpecs
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	} else {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		end = endFile.LineStart(endLine + 1) <span class="comment">// beginning of next line</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	first := len(f.Comments)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	last := -1
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	for i, g := range f.Comments {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		if g.End() &gt;= end {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			break
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		<span class="comment">// g.End() &lt; end</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		if beg &lt;= g.Pos() {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			<span class="comment">// comment is within the range [beg, end[ of import declarations</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			if i &lt; first {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>				first = i
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			if i &gt; last {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>				last = i
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	var comments []*CommentGroup
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	if last &gt;= 0 {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		comments = f.Comments[first : last+1]
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// Assign each comment to the import spec on the same line.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	importComments := map[*ImportSpec][]cgPos{}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	specIndex := 0
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	for _, g := range comments {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		for specIndex+1 &lt; len(specs) &amp;&amp; pos[specIndex+1].Start &lt;= g.Pos() {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			specIndex++
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		var left bool
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// A block comment can appear before the first import spec.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		if specIndex == 0 &amp;&amp; pos[specIndex].Start &gt; g.Pos() {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			left = true
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		} else if specIndex+1 &lt; len(specs) &amp;&amp; <span class="comment">// Or it can appear on the left of an import spec.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			lineAt(fset, pos[specIndex].Start)+1 == lineAt(fset, g.Pos()) {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			specIndex++
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			left = true
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		s := specs[specIndex].(*ImportSpec)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		importComments[s] = append(importComments[s], cgPos{left: left, cg: g})
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// Sort the import specs by import path.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// Remove duplicates, when possible without data loss.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// Reassign the import paths to have the same position sequence.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// Reassign each comment to the spec on the same line.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// Sort the comments by new position.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	sort.Slice(specs, func(i, j int) bool {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		ipath := importPath(specs[i])
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		jpath := importPath(specs[j])
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		if ipath != jpath {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			return ipath &lt; jpath
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		iname := importName(specs[i])
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		jname := importName(specs[j])
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		if iname != jname {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			return iname &lt; jname
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		return importComment(specs[i]) &lt; importComment(specs[j])
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	})
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// Dedup. Thanks to our sorting, we can just consider</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// adjacent pairs of imports.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	deduped := specs[:0]
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	for i, s := range specs {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		if i == len(specs)-1 || !collapse(s, specs[i+1]) {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			deduped = append(deduped, s)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		} else {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			p := s.Pos()
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			fset.File(p).MergeLine(lineAt(fset, p))
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	specs = deduped
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// Fix up comment positions</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	for i, s := range specs {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		s := s.(*ImportSpec)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		if s.Name != nil {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			s.Name.NamePos = pos[i].Start
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		s.Path.ValuePos = pos[i].Start
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		s.EndPos = pos[i].End
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		for _, g := range importComments[s] {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			for _, c := range g.cg.List {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>				if g.left {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>					c.Slash = pos[i].Start - 1
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>				} else {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>					<span class="comment">// An import spec can have both block comment and a line comment</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>					<span class="comment">// to its right. In that case, both of them will have the same pos.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>					<span class="comment">// But while formatting the AST, the line comment gets moved to</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>					<span class="comment">// after the block comment.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>					c.Slash = pos[i].End
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>				}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	sort.Slice(comments, func(i, j int) bool {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		return comments[i].Pos() &lt; comments[j].Pos()
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	})
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	return specs
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
</pre><p><a href="import.go?m=text">View as plain text</a></p>

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
