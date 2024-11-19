<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/printer/nodes.go - Go Documentation Server</title>

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
<a href="nodes.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/printer">printer</a>/<span class="text-muted">nodes.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/printer">go/printer</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements printing of AST nodes; specifically</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// expressions, statements, declarations, and files. It uses</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// the print functionality implemented in printer.go.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package printer
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// Formatting issues:</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// - better comment formatting for /*-style comments at the end of a line (e.g. a declaration)</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//   when the comment spans multiple lines; if such a comment is just two lines, formatting is</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//   not idempotent</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// - formatting of expression lists</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// - should use blank instead of tab to separate one-line function bodies from</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//   the function header unless there is a group of consecutive one-liners</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// Common AST nodes.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// Print as many newlines as necessary (but at least min newlines) to get to</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// the current line. ws is printed before the first line break. If newSection</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// is set, the first line break is printed as formfeed. Returns 0 if no line</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// breaks were printed, returns 1 if there was exactly one newline printed,</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// and returns a value &gt; 1 if there was a formfeed or more than one newline</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// printed.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// TODO(gri): linebreak may add too many lines if the next statement at &#34;line&#34;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// is preceded by comments because the computation of n assumes</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// the current position before the comment and the target position</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// after the comment. Thus, after interspersing such comments, the</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// space taken up by them is not considered to reduce the number of</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// linebreaks. At the moment there is no easy way to know about</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// future (not yet interspersed) comments in this function.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>func (p *printer) linebreak(line, min int, ws whiteSpace, newSection bool) (nbreaks int) {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	n := max(nlimit(line-p.pos.Line), min)
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	if n &gt; 0 {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		p.print(ws)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		if newSection {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			p.print(formfeed)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			n--
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			nbreaks = 2
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		nbreaks += n
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		for ; n &gt; 0; n-- {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			p.print(newline)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	return
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// setComment sets g as the next comment if g != nil and if node comments</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// are enabled - this mode is used when printing source code fragments such</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// as exports only. It assumes that there is no pending comment in p.comments</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// and at most one pending comment in the p.comment cache.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>func (p *printer) setComment(g *ast.CommentGroup) {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	if g == nil || !p.useNodeComments {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		return
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if p.comments == nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		<span class="comment">// initialize p.comments lazily</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		p.comments = make([]*ast.CommentGroup, 1)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	} else if p.cindex &lt; len(p.comments) {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		<span class="comment">// for some reason there are pending comments; this</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		<span class="comment">// should never happen - handle gracefully and flush</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		<span class="comment">// all comments up to g, ignore anything after that</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		p.flush(p.posFor(g.List[0].Pos()), token.ILLEGAL)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		p.comments = p.comments[0:1]
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		<span class="comment">// in debug mode, report error</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		p.internalError(&#34;setComment found pending comments&#34;)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	p.comments[0] = g
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	p.cindex = 0
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// don&#39;t overwrite any pending comment in the p.comment cache</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// (there may be a pending comment when a line comment is</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// immediately followed by a lead comment with no other</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// tokens between)</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	if p.commentOffset == infinity {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		p.nextComment() <span class="comment">// get comment ready for use</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>type exprListMode uint
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>const (
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	commaTerm exprListMode = 1 &lt;&lt; iota <span class="comment">// list is optionally terminated by a comma</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	noIndent                           <span class="comment">// no extra indentation in multi-line lists</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// If indent is set, a multi-line identifier list is indented after the</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// first linebreak encountered.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>func (p *printer) identList(list []*ast.Ident, indent bool) {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// convert into an expression list so we can re-use exprList formatting</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	xlist := make([]ast.Expr, len(list))
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	for i, x := range list {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		xlist[i] = x
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	var mode exprListMode
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	if !indent {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		mode = noIndent
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	p.exprList(token.NoPos, xlist, 1, mode, token.NoPos, false)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>const filteredMsg = &#34;contains filtered or unexported fields&#34;
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// Print a list of expressions. If the list spans multiple</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// source lines, the original line breaks are respected between</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// expressions.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// TODO(gri) Consider rewriting this to be independent of []ast.Expr</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// so that we can use the algorithm for any kind of list</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">//	(e.g., pass list via a channel over which to range).</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>func (p *printer) exprList(prev0 token.Pos, list []ast.Expr, depth int, mode exprListMode, next0 token.Pos, isIncomplete bool) {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if len(list) == 0 {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		if isIncomplete {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			prev := p.posFor(prev0)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			next := p.posFor(next0)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			if prev.IsValid() &amp;&amp; prev.Line == next.Line {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>				p.print(&#34;/* &#34; + filteredMsg + &#34; */&#34;)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			} else {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>				p.print(newline)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>				p.print(indent, &#34;// &#34;+filteredMsg, unindent, newline)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		return
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	prev := p.posFor(prev0)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	next := p.posFor(next0)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	line := p.lineFor(list[0].Pos())
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	endLine := p.lineFor(list[len(list)-1].End())
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	if prev.IsValid() &amp;&amp; prev.Line == line &amp;&amp; line == endLine {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		<span class="comment">// all list entries on a single line</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		for i, x := range list {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			if i &gt; 0 {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>				<span class="comment">// use position of expression following the comma as</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>				<span class="comment">// comma position for correct comment placement</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>				p.setPos(x.Pos())
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>				p.print(token.COMMA, blank)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			p.expr0(x, depth)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		if isIncomplete {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			p.print(token.COMMA, blank, &#34;/* &#34;+filteredMsg+&#34; */&#34;)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		return
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// list entries span multiple lines;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// use source code positions to guide line breaks</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t add extra indentation if noIndent is set;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// i.e., pretend that the first line is already indented.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	ws := ignore
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	if mode&amp;noIndent == 0 {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		ws = indent
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// The first linebreak is always a formfeed since this section must not</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// depend on any previous formatting.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	prevBreak := -1 <span class="comment">// index of last expression that was followed by a linebreak</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	if prev.IsValid() &amp;&amp; prev.Line &lt; line &amp;&amp; p.linebreak(line, 0, ws, true) &gt; 0 {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		ws = ignore
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		prevBreak = 0
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// initialize expression/key size: a zero value indicates expr/key doesn&#39;t fit on a single line</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	size := 0
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// We use the ratio between the geometric mean of the previous key sizes and</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// the current size to determine if there should be a break in the alignment.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// To compute the geometric mean we accumulate the ln(size) values (lnsum)</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// and the number of sizes included (count).</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	lnsum := 0.0
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	count := 0
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// print all list elements</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	prevLine := prev.Line
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	for i, x := range list {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		line = p.lineFor(x.Pos())
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		<span class="comment">// Determine if the next linebreak, if any, needs to use formfeed:</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		<span class="comment">// in general, use the entire node size to make the decision; for</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		<span class="comment">// key:value expressions, use the key size.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri) for a better result, should probably incorporate both</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		<span class="comment">//           the key and the node size into the decision process</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		useFF := true
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		<span class="comment">// Determine element size: All bets are off if we don&#39;t have</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		<span class="comment">// position information for the previous and next token (likely</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		<span class="comment">// generated code - simply ignore the size in this case by setting</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		<span class="comment">// it to 0).</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		prevSize := size
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		const infinity = 1e6 <span class="comment">// larger than any source line</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		size = p.nodeSize(x, infinity)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		pair, isPair := x.(*ast.KeyValueExpr)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		if size &lt;= infinity &amp;&amp; prev.IsValid() &amp;&amp; next.IsValid() {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			<span class="comment">// x fits on a single line</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			if isPair {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>				size = p.nodeSize(pair.Key, infinity) <span class="comment">// size &lt;= infinity</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		} else {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			<span class="comment">// size too large or we don&#39;t have good layout information</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			size = 0
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		<span class="comment">// If the previous line and the current line had single-</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		<span class="comment">// line-expressions and the key sizes are small or the</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		<span class="comment">// ratio between the current key and the geometric mean</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		<span class="comment">// if the previous key sizes does not exceed a threshold,</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		<span class="comment">// align columns and do not use formfeed.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		if prevSize &gt; 0 &amp;&amp; size &gt; 0 {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			const smallSize = 40
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			if count == 0 || prevSize &lt;= smallSize &amp;&amp; size &lt;= smallSize {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>				useFF = false
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			} else {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>				const r = 2.5                               <span class="comment">// threshold</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>				geomean := math.Exp(lnsum / float64(count)) <span class="comment">// count &gt; 0</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>				ratio := float64(size) / geomean
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>				useFF = r*ratio &lt;= 1 || r &lt;= ratio
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		needsLinebreak := 0 &lt; prevLine &amp;&amp; prevLine &lt; line
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		if i &gt; 0 {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			<span class="comment">// Use position of expression following the comma as</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			<span class="comment">// comma position for correct comment placement, but</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			<span class="comment">// only if the expression is on the same line.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			if !needsLinebreak {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>				p.setPos(x.Pos())
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			p.print(token.COMMA)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			needsBlank := true
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			if needsLinebreak {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>				<span class="comment">// Lines are broken using newlines so comments remain aligned</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>				<span class="comment">// unless useFF is set or there are multiple expressions on</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>				<span class="comment">// the same line in which case formfeed is used.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>				nbreaks := p.linebreak(line, 0, ws, useFF || prevBreak+1 &lt; i)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>				if nbreaks &gt; 0 {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>					ws = ignore
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>					prevBreak = i
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>					needsBlank = false <span class="comment">// we got a line break instead</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>				}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>				<span class="comment">// If there was a new section or more than one new line</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>				<span class="comment">// (which means that the tabwriter will implicitly break</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>				<span class="comment">// the section), reset the geomean variables since we are</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>				<span class="comment">// starting a new group of elements with the next element.</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				if nbreaks &gt; 1 {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>					lnsum = 0
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>					count = 0
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>				}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			if needsBlank {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				p.print(blank)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		if len(list) &gt; 1 &amp;&amp; isPair &amp;&amp; size &gt; 0 &amp;&amp; needsLinebreak {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			<span class="comment">// We have a key:value expression that fits onto one line</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			<span class="comment">// and it&#39;s not on the same line as the prior expression:</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			<span class="comment">// Use a column for the key such that consecutive entries</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			<span class="comment">// can align if possible.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			<span class="comment">// (needsLinebreak is set if we started a new line before)</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			p.expr(pair.Key)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			p.setPos(pair.Colon)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			p.print(token.COLON, vtab)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			p.expr(pair.Value)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		} else {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			p.expr0(x, depth)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		if size &gt; 0 {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			lnsum += math.Log(float64(size))
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			count++
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		prevLine = line
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	if mode&amp;commaTerm != 0 &amp;&amp; next.IsValid() &amp;&amp; p.pos.Line &lt; next.Line {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		<span class="comment">// Print a terminating comma if the next token is on a new line.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		p.print(token.COMMA)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		if isIncomplete {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			p.print(newline)
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			p.print(&#34;// &#34; + filteredMsg)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		if ws == ignore &amp;&amp; mode&amp;noIndent == 0 {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			<span class="comment">// unindent if we indented</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			p.print(unindent)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		p.print(formfeed) <span class="comment">// terminating comma needs a line break to look good</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		return
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	if isIncomplete {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		p.print(token.COMMA, newline)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		p.print(&#34;// &#34;+filteredMsg, newline)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	if ws == ignore &amp;&amp; mode&amp;noIndent == 0 {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		<span class="comment">// unindent if we indented</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		p.print(unindent)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>type paramMode int
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>const (
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	funcParam paramMode = iota
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	funcTParam
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	typeTParam
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>func (p *printer) parameters(fields *ast.FieldList, mode paramMode) {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	openTok, closeTok := token.LPAREN, token.RPAREN
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	if mode != funcParam {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		openTok, closeTok = token.LBRACK, token.RBRACK
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	p.setPos(fields.Opening)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	p.print(openTok)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	if len(fields.List) &gt; 0 {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		prevLine := p.lineFor(fields.Opening)
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		ws := indent
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		for i, par := range fields.List {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>			<span class="comment">// determine par begin and end line (may be different</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			<span class="comment">// if there are multiple parameter names for this par</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			<span class="comment">// or the type is on a separate line)</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			parLineBeg := p.lineFor(par.Pos())
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			parLineEnd := p.lineFor(par.End())
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			<span class="comment">// separating &#34;,&#34; if needed</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			needsLinebreak := 0 &lt; prevLine &amp;&amp; prevLine &lt; parLineBeg
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>			if i &gt; 0 {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>				<span class="comment">// use position of parameter following the comma as</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>				<span class="comment">// comma position for correct comma placement, but</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>				<span class="comment">// only if the next parameter is on the same line</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>				if !needsLinebreak {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>					p.setPos(par.Pos())
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>				}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>				p.print(token.COMMA)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>			}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			<span class="comment">// separator if needed (linebreak or blank)</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			if needsLinebreak &amp;&amp; p.linebreak(parLineBeg, 0, ws, true) &gt; 0 {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>				<span class="comment">// break line if the opening &#34;(&#34; or previous parameter ended on a different line</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>				ws = ignore
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			} else if i &gt; 0 {
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>				p.print(blank)
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>			}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>			<span class="comment">// parameter names</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			if len(par.Names) &gt; 0 {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>				<span class="comment">// Very subtle: If we indented before (ws == ignore), identList</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>				<span class="comment">// won&#39;t indent again. If we didn&#39;t (ws == indent), identList will</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>				<span class="comment">// indent if the identList spans multiple lines, and it will outdent</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>				<span class="comment">// again at the end (and still ws == indent). Thus, a subsequent indent</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>				<span class="comment">// by a linebreak call after a type, or in the next multi-line identList</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>				<span class="comment">// will do the right thing.</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>				p.identList(par.Names, ws == indent)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>				p.print(blank)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>			}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			<span class="comment">// parameter type</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			p.expr(stripParensAlways(par.Type))
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			prevLine = parLineEnd
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		<span class="comment">// if the closing &#34;)&#34; is on a separate line from the last parameter,</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		<span class="comment">// print an additional &#34;,&#34; and line break</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		if closing := p.lineFor(fields.Closing); 0 &lt; prevLine &amp;&amp; prevLine &lt; closing {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			p.print(token.COMMA)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			p.linebreak(closing, 0, ignore, true)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		} else if mode == typeTParam &amp;&amp; fields.NumFields() == 1 &amp;&amp; combinesWithName(fields.List[0].Type) {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			<span class="comment">// A type parameter list [P T] where the name P and the type expression T syntactically</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>			<span class="comment">// combine to another valid (value) expression requires a trailing comma, as in [P *T,]</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			<span class="comment">// (or an enclosing interface as in [P interface(*T)]), so that the type parameter list</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>			<span class="comment">// is not parsed as an array length [P*T].</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			p.print(token.COMMA)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		<span class="comment">// unindent if we indented</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		if ws == ignore {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			p.print(unindent)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	}
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	p.setPos(fields.Closing)
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	p.print(closeTok)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span><span class="comment">// combinesWithName reports whether a name followed by the expression x</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span><span class="comment">// syntactically combines to another valid (value) expression. For instance</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span><span class="comment">// using *T for x, &#34;name *T&#34; syntactically appears as the expression x*T.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span><span class="comment">// On the other hand, using  P|Q or *P|~Q for x, &#34;name P|Q&#34; or name *P|~Q&#34;</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span><span class="comment">// cannot be combined into a valid (value) expression.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>func combinesWithName(x ast.Expr) bool {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	switch x := x.(type) {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	case *ast.StarExpr:
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		<span class="comment">// name *x.X combines to name*x.X if x.X is not a type element</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		return !isTypeElem(x.X)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		return combinesWithName(x.X) &amp;&amp; !isTypeElem(x.Y)
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	case *ast.ParenExpr:
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		<span class="comment">// name(x) combines but we are making sure at</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		<span class="comment">// the call site that x is never parenthesized.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		panic(&#34;unexpected parenthesized expression&#34;)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	return false
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>}
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span><span class="comment">// isTypeElem reports whether x is a (possibly parenthesized) type element expression.</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span><span class="comment">// The result is false if x could be a type element OR an ordinary (value) expression.</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>func isTypeElem(x ast.Expr) bool {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	switch x := x.(type) {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	case *ast.ArrayType, *ast.StructType, *ast.FuncType, *ast.InterfaceType, *ast.MapType, *ast.ChanType:
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		return true
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	case *ast.UnaryExpr:
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		return x.Op == token.TILDE
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		return isTypeElem(x.X) || isTypeElem(x.Y)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	case *ast.ParenExpr:
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		return isTypeElem(x.X)
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	return false
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>func (p *printer) signature(sig *ast.FuncType) {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	if sig.TypeParams != nil {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		p.parameters(sig.TypeParams, funcTParam)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	if sig.Params != nil {
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		p.parameters(sig.Params, funcParam)
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	} else {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		p.print(token.LPAREN, token.RPAREN)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	res := sig.Results
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	n := res.NumFields()
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	if n &gt; 0 {
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		<span class="comment">// res != nil</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		p.print(blank)
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		if n == 1 &amp;&amp; res.List[0].Names == nil {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			<span class="comment">// single anonymous res; no ()&#39;s</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			p.expr(stripParensAlways(res.List[0].Type))
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>			return
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		p.parameters(res, funcParam)
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>func identListSize(list []*ast.Ident, maxSize int) (size int) {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	for i, x := range list {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		if i &gt; 0 {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			size += len(&#34;, &#34;)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		}
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		size += utf8.RuneCountInString(x.Name)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		if size &gt;= maxSize {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>			break
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	return
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>func (p *printer) isOneLineFieldList(list []*ast.Field) bool {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	if len(list) != 1 {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		return false <span class="comment">// allow only one field</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	f := list[0]
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	if f.Tag != nil || f.Comment != nil {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		return false <span class="comment">// don&#39;t allow tags or comments</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	<span class="comment">// only name(s) and type</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	const maxSize = 30 <span class="comment">// adjust as appropriate, this is an approximate value</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	namesSize := identListSize(f.Names, maxSize)
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	if namesSize &gt; 0 {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		namesSize = 1 <span class="comment">// blank between names and types</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	}
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	typeSize := p.nodeSize(f.Type, maxSize)
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	return namesSize+typeSize &lt;= maxSize
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>func (p *printer) setLineComment(text string) {
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	p.setComment(&amp;ast.CommentGroup{List: []*ast.Comment{{Slash: token.NoPos, Text: text}}})
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>func (p *printer) fieldList(fields *ast.FieldList, isStruct, isIncomplete bool) {
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	lbrace := fields.Opening
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	list := fields.List
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	rbrace := fields.Closing
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	hasComments := isIncomplete || p.commentBefore(p.posFor(rbrace))
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	srcIsOneLine := lbrace.IsValid() &amp;&amp; rbrace.IsValid() &amp;&amp; p.lineFor(lbrace) == p.lineFor(rbrace)
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	if !hasComments &amp;&amp; srcIsOneLine {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		<span class="comment">// possibly a one-line struct/interface</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		if len(list) == 0 {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>			<span class="comment">// no blank between keyword and {} in this case</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			p.setPos(lbrace)
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>			p.print(token.LBRACE)
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			p.setPos(rbrace)
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			p.print(token.RBRACE)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			return
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		} else if p.isOneLineFieldList(list) {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			<span class="comment">// small enough - print on one line</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>			<span class="comment">// (don&#39;t use identList and ignore source line breaks)</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>			p.setPos(lbrace)
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>			p.print(token.LBRACE, blank)
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>			f := list[0]
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			if isStruct {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>				for i, x := range f.Names {
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>					if i &gt; 0 {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>						<span class="comment">// no comments so no need for comma position</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>						p.print(token.COMMA, blank)
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>					}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>					p.expr(x)
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>				}
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>				if len(f.Names) &gt; 0 {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>					p.print(blank)
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>				}
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>				p.expr(f.Type)
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			} else { <span class="comment">// interface</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>				if len(f.Names) &gt; 0 {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>					name := f.Names[0] <span class="comment">// method name</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>					p.expr(name)
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>					p.signature(f.Type.(*ast.FuncType)) <span class="comment">// don&#39;t print &#34;func&#34;</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>				} else {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>					<span class="comment">// embedded interface</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>					p.expr(f.Type)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>				}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>			}
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>			p.print(blank)
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			p.setPos(rbrace)
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>			p.print(token.RBRACE)
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			return
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	}
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	<span class="comment">// hasComments || !srcIsOneLine</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	p.print(blank)
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	p.setPos(lbrace)
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	p.print(token.LBRACE, indent)
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	if hasComments || len(list) &gt; 0 {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		p.print(formfeed)
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	}
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	if isStruct {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		sep := vtab
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		if len(list) == 1 {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			sep = blank
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		}
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		var line int
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		for i, f := range list {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>			if i &gt; 0 {
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>				p.linebreak(p.lineFor(f.Pos()), 1, ignore, p.linesFrom(line) &gt; 0)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>			}
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>			extraTabs := 0
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>			p.setComment(f.Doc)
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>			p.recordLine(&amp;line)
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>			if len(f.Names) &gt; 0 {
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>				<span class="comment">// named fields</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>				p.identList(f.Names, false)
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>				p.print(sep)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>				p.expr(f.Type)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>				extraTabs = 1
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>			} else {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>				<span class="comment">// anonymous field</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>				p.expr(f.Type)
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>				extraTabs = 2
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>			}
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>			if f.Tag != nil {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>				if len(f.Names) &gt; 0 &amp;&amp; sep == vtab {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>					p.print(sep)
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>				}
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>				p.print(sep)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>				p.expr(f.Tag)
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>				extraTabs = 0
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>			}
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>			if f.Comment != nil {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>				for ; extraTabs &gt; 0; extraTabs-- {
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>					p.print(sep)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>				}
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>				p.setComment(f.Comment)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>			}
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		if isIncomplete {
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>			if len(list) &gt; 0 {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>				p.print(formfeed)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>			}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>			p.flush(p.posFor(rbrace), token.RBRACE) <span class="comment">// make sure we don&#39;t lose the last line comment</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>			p.setLineComment(&#34;// &#34; + filteredMsg)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	} else { <span class="comment">// interface</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		var line int
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		var prev *ast.Ident <span class="comment">// previous &#34;type&#34; identifier</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		for i, f := range list {
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>			var name *ast.Ident <span class="comment">// first name, or nil</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>			if len(f.Names) &gt; 0 {
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>				name = f.Names[0]
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>			}
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>			if i &gt; 0 {
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>				<span class="comment">// don&#39;t do a line break (min == 0) if we are printing a list of types</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>				<span class="comment">// TODO(gri) this doesn&#39;t work quite right if the list of types is</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>				<span class="comment">//           spread across multiple lines</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>				min := 1
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>				if prev != nil &amp;&amp; name == prev {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>					min = 0
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>				}
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>				p.linebreak(p.lineFor(f.Pos()), min, ignore, p.linesFrom(line) &gt; 0)
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			}
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			p.setComment(f.Doc)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>			p.recordLine(&amp;line)
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>			if name != nil {
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>				<span class="comment">// method</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>				p.expr(name)
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>				p.signature(f.Type.(*ast.FuncType)) <span class="comment">// don&#39;t print &#34;func&#34;</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>				prev = nil
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			} else {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>				<span class="comment">// embedded interface</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>				p.expr(f.Type)
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>				prev = nil
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>			p.setComment(f.Comment)
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		}
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		if isIncomplete {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>			if len(list) &gt; 0 {
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>				p.print(formfeed)
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>			}
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>			p.flush(p.posFor(rbrace), token.RBRACE) <span class="comment">// make sure we don&#39;t lose the last line comment</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			p.setLineComment(&#34;// contains filtered or unexported methods&#34;)
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	p.print(unindent, formfeed)
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	p.setPos(rbrace)
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	p.print(token.RBRACE)
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span><span class="comment">// Expressions</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>func walkBinary(e *ast.BinaryExpr) (has4, has5 bool, maxProblem int) {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	switch e.Op.Precedence() {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	case 4:
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		has4 = true
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	case 5:
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		has5 = true
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	switch l := e.X.(type) {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		if l.Op.Precedence() &lt; e.Op.Precedence() {
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>			<span class="comment">// parens will be inserted.</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>			<span class="comment">// pretend this is an *ast.ParenExpr and do nothing.</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>			break
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		h4, h5, mp := walkBinary(l)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		has4 = has4 || h4
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		has5 = has5 || h5
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		maxProblem = max(maxProblem, mp)
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	}
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	switch r := e.Y.(type) {
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>		if r.Op.Precedence() &lt;= e.Op.Precedence() {
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>			<span class="comment">// parens will be inserted.</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>			<span class="comment">// pretend this is an *ast.ParenExpr and do nothing.</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>			break
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		}
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		h4, h5, mp := walkBinary(r)
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		has4 = has4 || h4
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		has5 = has5 || h5
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		maxProblem = max(maxProblem, mp)
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	case *ast.StarExpr:
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		if e.Op == token.QUO { <span class="comment">// `*/`</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>			maxProblem = 5
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	case *ast.UnaryExpr:
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		switch e.Op.String() + r.Op.String() {
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		case &#34;/*&#34;, &#34;&amp;&amp;&#34;, &#34;&amp;^&#34;:
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>			maxProblem = 5
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		case &#34;++&#34;, &#34;--&#34;:
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>			maxProblem = max(maxProblem, 4)
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		}
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	return
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>}
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>func cutoff(e *ast.BinaryExpr, depth int) int {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	has4, has5, maxProblem := walkBinary(e)
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	if maxProblem &gt; 0 {
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		return maxProblem + 1
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	}
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	if has4 &amp;&amp; has5 {
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		if depth == 1 {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			return 5
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		}
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>		return 4
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	if depth == 1 {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		return 6
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	}
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	return 4
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>}
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>func diffPrec(expr ast.Expr, prec int) int {
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	x, ok := expr.(*ast.BinaryExpr)
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	if !ok || prec != x.Op.Precedence() {
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		return 1
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	}
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	return 0
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>}
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>func reduceDepth(depth int) int {
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>	depth--
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	if depth &lt; 1 {
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		depth = 1
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	}
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	return depth
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>}
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span><span class="comment">// Format the binary expression: decide the cutoff and then format.</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span><span class="comment">// Let&#39;s call depth == 1 Normal mode, and depth &gt; 1 Compact mode.</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span><span class="comment">// (Algorithm suggestion by Russ Cox.)</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span><span class="comment">// The precedences are:</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span><span class="comment">//	5             *  /  %  &lt;&lt;  &gt;&gt;  &amp;  &amp;^</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span><span class="comment">//	4             +  -  |  ^</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span><span class="comment">//	3             ==  !=  &lt;  &lt;=  &gt;  &gt;=</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span><span class="comment">//	2             &amp;&amp;</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span><span class="comment">//	1             ||</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span><span class="comment">// The only decision is whether there will be spaces around levels 4 and 5.</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span><span class="comment">// There are never spaces at level 6 (unary), and always spaces at levels 3 and below.</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span><span class="comment">// To choose the cutoff, look at the whole expression but excluding primary</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span><span class="comment">// expressions (function calls, parenthesized exprs), and apply these rules:</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span><span class="comment">//  1. If there is a binary operator with a right side unary operand</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span><span class="comment">//     that would clash without a space, the cutoff must be (in order):</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span><span class="comment">//     /*	6</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span><span class="comment">//     &amp;&amp;	6</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span><span class="comment">//     &amp;^	6</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span><span class="comment">//     ++	5</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span><span class="comment">//     --	5</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span><span class="comment">//     (Comparison operators always have spaces around them.)</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span><span class="comment">//  2. If there is a mix of level 5 and level 4 operators, then the cutoff</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span><span class="comment">//     is 5 (use spaces to distinguish precedence) in Normal mode</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span><span class="comment">//     and 4 (never use spaces) in Compact mode.</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span><span class="comment">//  3. If there are no level 4 operators or no level 5 operators, then the</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span><span class="comment">//     cutoff is 6 (always use spaces) in Normal mode</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span><span class="comment">//     and 4 (never use spaces) in Compact mode.</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>func (p *printer) binaryExpr(x *ast.BinaryExpr, prec1, cutoff, depth int) {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	prec := x.Op.Precedence()
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	if prec &lt; prec1 {
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		<span class="comment">// parenthesis needed</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>		<span class="comment">// Note: The parser inserts an ast.ParenExpr node; thus this case</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		<span class="comment">//       can only occur if the AST is created in a different way.</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		p.print(token.LPAREN)
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>		p.expr0(x, reduceDepth(depth)) <span class="comment">// parentheses undo one level of depth</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		p.print(token.RPAREN)
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		return
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	}
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	printBlank := prec &lt; cutoff
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	ws := indent
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	p.expr1(x.X, prec, depth+diffPrec(x.X, prec))
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	if printBlank {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		p.print(blank)
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>	}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	xline := p.pos.Line <span class="comment">// before the operator (it may be on the next line!)</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	yline := p.lineFor(x.Y.Pos())
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	p.setPos(x.OpPos)
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	p.print(x.Op)
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	if xline != yline &amp;&amp; xline &gt; 0 &amp;&amp; yline &gt; 0 {
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>		<span class="comment">// at least one line break, but respect an extra empty line</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>		<span class="comment">// in the source</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		if p.linebreak(yline, 1, ws, true) &gt; 0 {
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>			ws = ignore
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>			printBlank = false <span class="comment">// no blank after line break</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		}
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	}
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	if printBlank {
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		p.print(blank)
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	}
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	p.expr1(x.Y, prec+1, depth+1)
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	if ws == ignore {
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		p.print(unindent)
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	}
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>}
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>func isBinary(expr ast.Expr) bool {
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	_, ok := expr.(*ast.BinaryExpr)
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	return ok
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>}
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>func (p *printer) expr1(expr ast.Expr, prec1, depth int) {
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	p.setPos(expr.Pos())
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>	switch x := expr.(type) {
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	case *ast.BadExpr:
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>		p.print(&#34;BadExpr&#34;)
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	case *ast.Ident:
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		p.print(x)
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		if depth &lt; 1 {
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>			p.internalError(&#34;depth &lt; 1:&#34;, depth)
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>			depth = 1
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		}
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		p.binaryExpr(x, prec1, cutoff(x, depth), depth)
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	case *ast.KeyValueExpr:
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>		p.expr(x.Key)
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		p.setPos(x.Colon)
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		p.print(token.COLON, blank)
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		p.expr(x.Value)
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	case *ast.StarExpr:
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		const prec = token.UnaryPrec
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		if prec &lt; prec1 {
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>			<span class="comment">// parenthesis needed</span>
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>			p.print(token.LPAREN)
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>			p.print(token.MUL)
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>			p.expr(x.X)
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			p.print(token.RPAREN)
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		} else {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>			<span class="comment">// no parenthesis needed</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>			p.print(token.MUL)
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>			p.expr(x.X)
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	case *ast.UnaryExpr:
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>		const prec = token.UnaryPrec
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		if prec &lt; prec1 {
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>			<span class="comment">// parenthesis needed</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>			p.print(token.LPAREN)
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>			p.expr(x)
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>			p.print(token.RPAREN)
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		} else {
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>			<span class="comment">// no parenthesis needed</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>			p.print(x.Op)
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>			if x.Op == token.RANGE {
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>				<span class="comment">// TODO(gri) Remove this code if it cannot be reached.</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>				p.print(blank)
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>			}
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>			p.expr1(x.X, prec, depth)
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		}
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	case *ast.BasicLit:
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>		if p.Config.Mode&amp;normalizeNumbers != 0 {
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>			x = normalizedNumber(x)
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		}
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		p.print(x)
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	case *ast.FuncLit:
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>		p.setPos(x.Type.Pos())
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>		p.print(token.FUNC)
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>		<span class="comment">// See the comment in funcDecl about how the header size is computed.</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>		startCol := p.out.Column - len(&#34;func&#34;)
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>		p.signature(x.Type)
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>		p.funcBody(p.distanceFrom(x.Type.Pos(), startCol), blank, x.Body)
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	case *ast.ParenExpr:
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>		if _, hasParens := x.X.(*ast.ParenExpr); hasParens {
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>			<span class="comment">// don&#39;t print parentheses around an already parenthesized expression</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>			<span class="comment">// TODO(gri) consider making this more general and incorporate precedence levels</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>			p.expr0(x.X, depth)
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>		} else {
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			p.print(token.LPAREN)
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>			p.expr0(x.X, reduceDepth(depth)) <span class="comment">// parentheses undo one level of depth</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			p.setPos(x.Rparen)
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>			p.print(token.RPAREN)
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>		}
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	case *ast.SelectorExpr:
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>		p.selectorExpr(x, depth, false)
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	case *ast.TypeAssertExpr:
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>		p.expr1(x.X, token.HighestPrec, depth)
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		p.print(token.PERIOD)
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		p.setPos(x.Lparen)
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>		p.print(token.LPAREN)
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>		if x.Type != nil {
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>			p.expr(x.Type)
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>		} else {
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>			p.print(token.TYPE)
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>		}
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>		p.setPos(x.Rparen)
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>		p.print(token.RPAREN)
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>	case *ast.IndexExpr:
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri): should treat[] like parentheses and undo one level of depth</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>		p.expr1(x.X, token.HighestPrec, 1)
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>		p.setPos(x.Lbrack)
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		p.print(token.LBRACK)
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>		p.expr0(x.Index, depth+1)
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		p.setPos(x.Rbrack)
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		p.print(token.RBRACK)
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>	case *ast.IndexListExpr:
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri): as for IndexExpr, should treat [] like parentheses and undo</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>		<span class="comment">// one level of depth</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>		p.expr1(x.X, token.HighestPrec, 1)
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>		p.setPos(x.Lbrack)
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>		p.print(token.LBRACK)
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>		p.exprList(x.Lbrack, x.Indices, depth+1, commaTerm, x.Rbrack, false)
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>		p.setPos(x.Rbrack)
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>		p.print(token.RBRACK)
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>	case *ast.SliceExpr:
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>		<span class="comment">// TODO(gri): should treat[] like parentheses and undo one level of depth</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>		p.expr1(x.X, token.HighestPrec, 1)
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>		p.setPos(x.Lbrack)
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>		p.print(token.LBRACK)
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>		indices := []ast.Expr{x.Low, x.High}
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>		if x.Max != nil {
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>			indices = append(indices, x.Max)
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>		}
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>		<span class="comment">// determine if we need extra blanks around &#39;:&#39;</span>
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>		var needsBlanks bool
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		if depth &lt;= 1 {
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>			var indexCount int
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>			var hasBinaries bool
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>			for _, x := range indices {
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>				if x != nil {
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>					indexCount++
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>					if isBinary(x) {
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>						hasBinaries = true
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>					}
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>				}
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>			}
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			if indexCount &gt; 1 &amp;&amp; hasBinaries {
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>				needsBlanks = true
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>			}
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>		}
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>		for i, x := range indices {
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>			if i &gt; 0 {
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>				if indices[i-1] != nil &amp;&amp; needsBlanks {
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>					p.print(blank)
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>				}
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>				p.print(token.COLON)
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>				if x != nil &amp;&amp; needsBlanks {
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>					p.print(blank)
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>				}
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>			}
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>			if x != nil {
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>				p.expr0(x, depth+1)
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>			}
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>		}
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>		p.setPos(x.Rbrack)
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>		p.print(token.RBRACK)
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	case *ast.CallExpr:
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>		if len(x.Args) &gt; 1 {
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>			depth++
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>		}
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>		<span class="comment">// Conversions to literal function types or &lt;-chan</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>		<span class="comment">// types require parentheses around the type.</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>		paren := false
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>		switch t := x.Fun.(type) {
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		case *ast.FuncType:
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>			paren = true
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>		case *ast.ChanType:
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>			paren = t.Dir == ast.RECV
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>		}
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>		if paren {
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>			p.print(token.LPAREN)
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>		}
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>		wasIndented := p.possibleSelectorExpr(x.Fun, token.HighestPrec, depth)
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>		if paren {
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>			p.print(token.RPAREN)
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>		}
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		p.setPos(x.Lparen)
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>		p.print(token.LPAREN)
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>		if x.Ellipsis.IsValid() {
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>			p.exprList(x.Lparen, x.Args, depth, 0, x.Ellipsis, false)
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>			p.setPos(x.Ellipsis)
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>			p.print(token.ELLIPSIS)
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>			if x.Rparen.IsValid() &amp;&amp; p.lineFor(x.Ellipsis) &lt; p.lineFor(x.Rparen) {
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>				p.print(token.COMMA, formfeed)
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>			}
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>		} else {
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>			p.exprList(x.Lparen, x.Args, depth, commaTerm, x.Rparen, false)
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>		}
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>		p.setPos(x.Rparen)
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>		p.print(token.RPAREN)
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>		if wasIndented {
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>			p.print(unindent)
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>		}
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	case *ast.CompositeLit:
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>		<span class="comment">// composite literal elements that are composite literals themselves may have the type omitted</span>
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>		if x.Type != nil {
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>			p.expr1(x.Type, token.HighestPrec, depth)
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>		}
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>		p.level++
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>		p.setPos(x.Lbrace)
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>		p.print(token.LBRACE)
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		p.exprList(x.Lbrace, x.Elts, 1, commaTerm, x.Rbrace, x.Incomplete)
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>		<span class="comment">// do not insert extra line break following a /*-style comment</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		<span class="comment">// before the closing &#39;}&#39; as it might break the code if there</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>		<span class="comment">// is no trailing &#39;,&#39;</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		mode := noExtraLinebreak
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>		<span class="comment">// do not insert extra blank following a /*-style comment</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>		<span class="comment">// before the closing &#39;}&#39; unless the literal is empty</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>		if len(x.Elts) &gt; 0 {
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>			mode |= noExtraBlank
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>		}
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>		<span class="comment">// need the initial indent to print lone comments with</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>		<span class="comment">// the proper level of indentation</span>
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		p.print(indent, unindent, mode)
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>		p.setPos(x.Rbrace)
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>		p.print(token.RBRACE, mode)
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>		p.level--
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	case *ast.Ellipsis:
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>		p.print(token.ELLIPSIS)
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>		if x.Elt != nil {
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>			p.expr(x.Elt)
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		}
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>	case *ast.ArrayType:
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>		p.print(token.LBRACK)
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>		if x.Len != nil {
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>			p.expr(x.Len)
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>		}
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>		p.print(token.RBRACK)
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>		p.expr(x.Elt)
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>	case *ast.StructType:
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>		p.print(token.STRUCT)
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>		p.fieldList(x.Fields, true, x.Incomplete)
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>	case *ast.FuncType:
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>		p.print(token.FUNC)
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>		p.signature(x)
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>	case *ast.InterfaceType:
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>		p.print(token.INTERFACE)
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>		p.fieldList(x.Methods, false, x.Incomplete)
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>	case *ast.MapType:
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>		p.print(token.MAP, token.LBRACK)
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>		p.expr(x.Key)
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>		p.print(token.RBRACK)
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>		p.expr(x.Value)
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	case *ast.ChanType:
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>		switch x.Dir {
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>		case ast.SEND | ast.RECV:
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>			p.print(token.CHAN)
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>		case ast.RECV:
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>			p.print(token.ARROW, token.CHAN) <span class="comment">// x.Arrow and x.Pos() are the same</span>
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>		case ast.SEND:
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>			p.print(token.CHAN)
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>			p.setPos(x.Arrow)
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>			p.print(token.ARROW)
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>		}
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>		p.print(blank)
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>		p.expr(x.Value)
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>	default:
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>		panic(&#34;unreachable&#34;)
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>	}
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>}
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span><span class="comment">// normalizedNumber rewrites base prefixes and exponents</span>
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span><span class="comment">// of numbers to use lower-case letters (0X123 to 0x123 and 1.2E3 to 1.2e3),</span>
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span><span class="comment">// and removes leading 0&#39;s from integer imaginary literals (0765i to 765i).</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span><span class="comment">// It leaves hexadecimal digits alone.</span>
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span><span class="comment">// normalizedNumber doesn&#39;t modify the ast.BasicLit value lit points to.</span>
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span><span class="comment">// If lit is not a number or a number in canonical format already,</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span><span class="comment">// lit is returned as is. Otherwise a new ast.BasicLit is created.</span>
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>func normalizedNumber(lit *ast.BasicLit) *ast.BasicLit {
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>	if lit.Kind != token.INT &amp;&amp; lit.Kind != token.FLOAT &amp;&amp; lit.Kind != token.IMAG {
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>		return lit <span class="comment">// not a number - nothing to do</span>
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>	}
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	if len(lit.Value) &lt; 2 {
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>		return lit <span class="comment">// only one digit (common case) - nothing to do</span>
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>	}
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>	<span class="comment">// len(lit.Value) &gt;= 2</span>
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>	<span class="comment">// We ignore lit.Kind because for lit.Kind == token.IMAG the literal may be an integer</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>	<span class="comment">// or floating-point value, decimal or not. Instead, just consider the literal pattern.</span>
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>	x := lit.Value
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>	switch x[:2] {
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>	default:
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>		<span class="comment">// 0-prefix octal, decimal int, or float (possibly with &#39;i&#39; suffix)</span>
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>		if i := strings.LastIndexByte(x, &#39;E&#39;); i &gt;= 0 {
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>			x = x[:i] + &#34;e&#34; + x[i+1:]
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>			break
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>		}
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>		<span class="comment">// remove leading 0&#39;s from integer (but not floating-point) imaginary literals</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>		if x[len(x)-1] == &#39;i&#39; &amp;&amp; !strings.ContainsAny(x, &#34;.e&#34;) {
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>			x = strings.TrimLeft(x, &#34;0_&#34;)
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>			if x == &#34;i&#34; {
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>				x = &#34;0i&#34;
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>			}
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>		}
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	case &#34;0X&#34;:
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>		x = &#34;0x&#34; + x[2:]
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>		<span class="comment">// possibly a hexadecimal float</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>		if i := strings.LastIndexByte(x, &#39;P&#39;); i &gt;= 0 {
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>			x = x[:i] + &#34;p&#34; + x[i+1:]
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>		}
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>	case &#34;0x&#34;:
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>		<span class="comment">// possibly a hexadecimal float</span>
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>		i := strings.LastIndexByte(x, &#39;P&#39;)
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>		if i == -1 {
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>			return lit <span class="comment">// nothing to do</span>
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>		}
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>		x = x[:i] + &#34;p&#34; + x[i+1:]
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>	case &#34;0O&#34;:
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>		x = &#34;0o&#34; + x[2:]
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>	case &#34;0o&#34;:
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>		return lit <span class="comment">// nothing to do</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>	case &#34;0B&#34;:
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>		x = &#34;0b&#34; + x[2:]
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	case &#34;0b&#34;:
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>		return lit <span class="comment">// nothing to do</span>
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>	}
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>	return &amp;ast.BasicLit{ValuePos: lit.ValuePos, Kind: lit.Kind, Value: x}
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>}
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>func (p *printer) possibleSelectorExpr(expr ast.Expr, prec1, depth int) bool {
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>	if x, ok := expr.(*ast.SelectorExpr); ok {
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>		return p.selectorExpr(x, depth, true)
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>	}
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>	p.expr1(expr, prec1, depth)
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>	return false
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>}
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span><span class="comment">// selectorExpr handles an *ast.SelectorExpr node and reports whether x spans</span>
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span><span class="comment">// multiple lines.</span>
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>func (p *printer) selectorExpr(x *ast.SelectorExpr, depth int, isMethod bool) bool {
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>	p.expr1(x.X, token.HighestPrec, depth)
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>	p.print(token.PERIOD)
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>	if line := p.lineFor(x.Sel.Pos()); p.pos.IsValid() &amp;&amp; p.pos.Line &lt; line {
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>		p.print(indent, newline)
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>		p.setPos(x.Sel.Pos())
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>		p.print(x.Sel)
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>		if !isMethod {
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>			p.print(unindent)
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>		}
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>		return true
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>	}
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>	p.setPos(x.Sel.Pos())
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>	p.print(x.Sel)
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>	return false
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>}
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>func (p *printer) expr0(x ast.Expr, depth int) {
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>	p.expr1(x, token.LowestPrec, depth)
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>}
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>func (p *printer) expr(x ast.Expr) {
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>	const depth = 1
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>	p.expr1(x, token.LowestPrec, depth)
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>}
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span><span class="comment">// Statements</span>
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span><span class="comment">// Print the statement list indented, but without a newline after the last statement.</span>
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span><span class="comment">// Extra line breaks between statements in the source are respected but at most one</span>
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span><span class="comment">// empty line is printed between statements.</span>
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>func (p *printer) stmtList(list []ast.Stmt, nindent int, nextIsRBrace bool) {
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>	if nindent &gt; 0 {
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>		p.print(indent)
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>	}
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>	var line int
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>	i := 0
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>	for _, s := range list {
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>		<span class="comment">// ignore empty statements (was issue 3466)</span>
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>		if _, isEmpty := s.(*ast.EmptyStmt); !isEmpty {
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>			<span class="comment">// nindent == 0 only for lists of switch/select case clauses;</span>
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>			<span class="comment">// in those cases each clause is a new section</span>
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>			if len(p.output) &gt; 0 {
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>				<span class="comment">// only print line break if we are not at the beginning of the output</span>
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>				<span class="comment">// (i.e., we are not printing only a partial program)</span>
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>				p.linebreak(p.lineFor(s.Pos()), 1, ignore, i == 0 || nindent == 0 || p.linesFrom(line) &gt; 0)
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>			}
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>			p.recordLine(&amp;line)
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>			p.stmt(s, nextIsRBrace &amp;&amp; i == len(list)-1)
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>			<span class="comment">// labeled statements put labels on a separate line, but here</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>			<span class="comment">// we only care about the start line of the actual statement</span>
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>			<span class="comment">// without label - correct line for each label</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>			for t := s; ; {
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>				lt, _ := t.(*ast.LabeledStmt)
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>				if lt == nil {
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>					break
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>				}
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>				line++
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>				t = lt.Stmt
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>			}
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>			i++
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>		}
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>	}
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>	if nindent &gt; 0 {
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>		p.print(unindent)
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>	}
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>}
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span><span class="comment">// block prints an *ast.BlockStmt; it always spans at least two lines.</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>func (p *printer) block(b *ast.BlockStmt, nindent int) {
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>	p.setPos(b.Lbrace)
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	p.print(token.LBRACE)
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>	p.stmtList(b.List, nindent, true)
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>	p.linebreak(p.lineFor(b.Rbrace), 1, ignore, true)
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>	p.setPos(b.Rbrace)
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>	p.print(token.RBRACE)
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>}
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>func isTypeName(x ast.Expr) bool {
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>	switch t := x.(type) {
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>	case *ast.Ident:
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>		return true
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>	case *ast.SelectorExpr:
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>		return isTypeName(t.X)
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>	}
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>	return false
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>}
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>func stripParens(x ast.Expr) ast.Expr {
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>	if px, strip := x.(*ast.ParenExpr); strip {
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>		<span class="comment">// parentheses must not be stripped if there are any</span>
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>		<span class="comment">// unparenthesized composite literals starting with</span>
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>		<span class="comment">// a type name</span>
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>		ast.Inspect(px.X, func(node ast.Node) bool {
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>			switch x := node.(type) {
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>			case *ast.ParenExpr:
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>				<span class="comment">// parentheses protect enclosed composite literals</span>
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>				return false
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>			case *ast.CompositeLit:
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>				if isTypeName(x.Type) {
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>					strip = false <span class="comment">// do not strip parentheses</span>
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>				}
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>				return false
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>			}
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>			<span class="comment">// in all other cases, keep inspecting</span>
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>			return true
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>		})
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>		if strip {
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>			return stripParens(px.X)
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>		}
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	}
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	return x
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>}
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>func stripParensAlways(x ast.Expr) ast.Expr {
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>	if x, ok := x.(*ast.ParenExpr); ok {
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>		return stripParensAlways(x.X)
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>	}
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>	return x
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>}
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>func (p *printer) controlClause(isForStmt bool, init ast.Stmt, expr ast.Expr, post ast.Stmt) {
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>	p.print(blank)
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>	needsBlank := false
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>	if init == nil &amp;&amp; post == nil {
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>		<span class="comment">// no semicolons required</span>
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>		if expr != nil {
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>			p.expr(stripParens(expr))
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>			needsBlank = true
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>		}
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>	} else {
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>		<span class="comment">// all semicolons required</span>
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>		<span class="comment">// (they are not separators, print them explicitly)</span>
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>		if init != nil {
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>			p.stmt(init, false)
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>		}
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>		p.print(token.SEMICOLON, blank)
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>		if expr != nil {
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>			p.expr(stripParens(expr))
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>			needsBlank = true
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>		}
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>		if isForStmt {
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>			p.print(token.SEMICOLON, blank)
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>			needsBlank = false
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>			if post != nil {
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>				p.stmt(post, false)
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>				needsBlank = true
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>			}
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>		}
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>	}
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>	if needsBlank {
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>		p.print(blank)
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>	}
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>}
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span><span class="comment">// indentList reports whether an expression list would look better if it</span>
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span><span class="comment">// were indented wholesale (starting with the very first element, rather</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span><span class="comment">// than starting at the first line break).</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>func (p *printer) indentList(list []ast.Expr) bool {
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>	<span class="comment">// Heuristic: indentList reports whether there are more than one multi-</span>
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>	<span class="comment">// line element in the list, or if there is any element that is not</span>
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>	<span class="comment">// starting on the same line as the previous one ends.</span>
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>	if len(list) &gt;= 2 {
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>		var b = p.lineFor(list[0].Pos())
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>		var e = p.lineFor(list[len(list)-1].End())
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>		if 0 &lt; b &amp;&amp; b &lt; e {
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>			<span class="comment">// list spans multiple lines</span>
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>			n := 0 <span class="comment">// multi-line element count</span>
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>			line := b
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>			for _, x := range list {
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>				xb := p.lineFor(x.Pos())
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>				xe := p.lineFor(x.End())
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>				if line &lt; xb {
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>					<span class="comment">// x is not starting on the same</span>
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>					<span class="comment">// line as the previous one ended</span>
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>					return true
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>				}
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>				if xb &lt; xe {
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>					<span class="comment">// x is a multi-line element</span>
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>					n++
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>				}
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>				line = xe
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>			}
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>			return n &gt; 1
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>		}
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>	}
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>	return false
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>}
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>func (p *printer) stmt(stmt ast.Stmt, nextIsRBrace bool) {
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>	p.setPos(stmt.Pos())
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>	switch s := stmt.(type) {
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>	case *ast.BadStmt:
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>		p.print(&#34;BadStmt&#34;)
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>	case *ast.DeclStmt:
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>		p.decl(s.Decl)
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>	case *ast.EmptyStmt:
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>		<span class="comment">// nothing to do</span>
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>	case *ast.LabeledStmt:
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>		<span class="comment">// a &#34;correcting&#34; unindent immediately following a line break</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>		<span class="comment">// is applied before the line break if there is no comment</span>
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>		<span class="comment">// between (see writeWhitespace)</span>
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>		p.print(unindent)
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>		p.expr(s.Label)
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>		p.setPos(s.Colon)
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>		p.print(token.COLON, indent)
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>		if e, isEmpty := s.Stmt.(*ast.EmptyStmt); isEmpty {
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>			if !nextIsRBrace {
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>				p.print(newline)
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>				p.setPos(e.Pos())
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>				p.print(token.SEMICOLON)
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>				break
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>			}
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>		} else {
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>			p.linebreak(p.lineFor(s.Stmt.Pos()), 1, ignore, true)
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>		}
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>		p.stmt(s.Stmt, nextIsRBrace)
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>	case *ast.ExprStmt:
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>		const depth = 1
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>		p.expr0(s.X, depth)
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>	case *ast.SendStmt:
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>		const depth = 1
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>		p.expr0(s.Chan, depth)
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>		p.print(blank)
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>		p.setPos(s.Arrow)
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>		p.print(token.ARROW, blank)
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>		p.expr0(s.Value, depth)
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>	case *ast.IncDecStmt:
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>		const depth = 1
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>		p.expr0(s.X, depth+1)
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>		p.setPos(s.TokPos)
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>		p.print(s.Tok)
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>	case *ast.AssignStmt:
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>		var depth = 1
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>		if len(s.Lhs) &gt; 1 &amp;&amp; len(s.Rhs) &gt; 1 {
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>			depth++
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>		}
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>		p.exprList(s.Pos(), s.Lhs, depth, 0, s.TokPos, false)
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>		p.print(blank)
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>		p.setPos(s.TokPos)
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>		p.print(s.Tok, blank)
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>		p.exprList(s.TokPos, s.Rhs, depth, 0, token.NoPos, false)
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>	case *ast.GoStmt:
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>		p.print(token.GO, blank)
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>		p.expr(s.Call)
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>	case *ast.DeferStmt:
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>		p.print(token.DEFER, blank)
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>		p.expr(s.Call)
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>	case *ast.ReturnStmt:
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>		p.print(token.RETURN)
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>		if s.Results != nil {
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>			p.print(blank)
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>			<span class="comment">// Use indentList heuristic to make corner cases look</span>
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>			<span class="comment">// better (issue 1207). A more systematic approach would</span>
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>			<span class="comment">// always indent, but this would cause significant</span>
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>			<span class="comment">// reformatting of the code base and not necessarily</span>
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>			<span class="comment">// lead to more nicely formatted code in general.</span>
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>			if p.indentList(s.Results) {
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>				p.print(indent)
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span>				<span class="comment">// Use NoPos so that a newline never goes before</span>
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span>				<span class="comment">// the results (see issue #32854).</span>
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>				p.exprList(token.NoPos, s.Results, 1, noIndent, token.NoPos, false)
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>				p.print(unindent)
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>			} else {
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span>				p.exprList(token.NoPos, s.Results, 1, 0, token.NoPos, false)
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span>			}
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span>		}
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span>
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span>	case *ast.BranchStmt:
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span>		p.print(s.Tok)
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>		if s.Label != nil {
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>			p.print(blank)
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>			p.expr(s.Label)
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>		}
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>	case *ast.BlockStmt:
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>		p.block(s, 1)
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>	case *ast.IfStmt:
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span>		p.print(token.IF)
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span>		p.controlClause(false, s.Init, s.Cond, nil)
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span>		p.block(s.Body, 1)
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span>		if s.Else != nil {
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span>			p.print(blank, token.ELSE, blank)
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span>			switch s.Else.(type) {
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span>			case *ast.BlockStmt, *ast.IfStmt:
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span>				p.stmt(s.Else, nextIsRBrace)
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span>			default:
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span>				<span class="comment">// This can only happen with an incorrectly</span>
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span>				<span class="comment">// constructed AST. Permit it but print so</span>
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span>				<span class="comment">// that it can be parsed without errors.</span>
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span>				p.print(token.LBRACE, indent, formfeed)
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span>				p.stmt(s.Else, true)
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span>				p.print(unindent, formfeed, token.RBRACE)
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span>			}
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span>		}
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span>
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span>	case *ast.CaseClause:
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span>		if s.List != nil {
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span>			p.print(token.CASE, blank)
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span>			p.exprList(s.Pos(), s.List, 1, 0, s.Colon, false)
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span>		} else {
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span>			p.print(token.DEFAULT)
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span>		}
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span>		p.setPos(s.Colon)
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span>		p.print(token.COLON)
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span>		p.stmtList(s.Body, 1, nextIsRBrace)
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span>
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span>	case *ast.SwitchStmt:
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span>		p.print(token.SWITCH)
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span>		p.controlClause(false, s.Init, s.Tag, nil)
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span>		p.block(s.Body, 0)
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span>
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span>	case *ast.TypeSwitchStmt:
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span>		p.print(token.SWITCH)
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span>		if s.Init != nil {
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span>			p.print(blank)
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span>			p.stmt(s.Init, false)
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span>			p.print(token.SEMICOLON)
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span>		}
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span>		p.print(blank)
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span>		p.stmt(s.Assign, false)
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span>		p.print(blank)
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span>		p.block(s.Body, 0)
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span>
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span>	case *ast.CommClause:
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span>		if s.Comm != nil {
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span>			p.print(token.CASE, blank)
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span>			p.stmt(s.Comm, false)
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span>		} else {
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span>			p.print(token.DEFAULT)
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span>		}
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span>		p.setPos(s.Colon)
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span>		p.print(token.COLON)
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span>		p.stmtList(s.Body, 1, nextIsRBrace)
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span>
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span>	case *ast.SelectStmt:
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span>		p.print(token.SELECT, blank)
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span>		body := s.Body
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span>		if len(body.List) == 0 &amp;&amp; !p.commentBefore(p.posFor(body.Rbrace)) {
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span>			<span class="comment">// print empty select statement w/o comments on one line</span>
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span>			p.setPos(body.Lbrace)
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span>			p.print(token.LBRACE)
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span>			p.setPos(body.Rbrace)
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>			p.print(token.RBRACE)
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>		} else {
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span>			p.block(body, 0)
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span>		}
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span>
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span>	case *ast.ForStmt:
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span>		p.print(token.FOR)
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span>		p.controlClause(true, s.Init, s.Cond, s.Post)
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span>		p.block(s.Body, 1)
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span>
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span>	case *ast.RangeStmt:
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span>		p.print(token.FOR, blank)
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span>		if s.Key != nil {
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span>			p.expr(s.Key)
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span>			if s.Value != nil {
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>				<span class="comment">// use position of value following the comma as</span>
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>				<span class="comment">// comma position for correct comment placement</span>
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>				p.setPos(s.Value.Pos())
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>				p.print(token.COMMA, blank)
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span>				p.expr(s.Value)
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span>			}
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span>			p.print(blank)
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span>			p.setPos(s.TokPos)
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>			p.print(s.Tok, blank)
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>		}
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>		p.print(token.RANGE, blank)
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span>		p.expr(stripParens(s.X))
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span>		p.print(blank)
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span>		p.block(s.Body, 1)
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span>
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>	default:
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>		panic(&#34;unreachable&#34;)
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>	}
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>}
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span>
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span><span class="comment">// Declarations</span>
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span>
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span><span class="comment">// The keepTypeColumn function determines if the type column of a series of</span>
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span><span class="comment">// consecutive const or var declarations must be kept, or if initialization</span>
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span><span class="comment">// values (V) can be placed in the type column (T) instead. The i&#39;th entry</span>
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span><span class="comment">// in the result slice is true if the type column in spec[i] must be kept.</span>
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span><span class="comment">// For example, the declaration:</span>
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span><span class="comment">//		const (</span>
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span><span class="comment">//			foobar int = 42 // comment</span>
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span><span class="comment">//			x          = 7  // comment</span>
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span><span class="comment">//			foo</span>
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span><span class="comment">//	             bar = 991</span>
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span><span class="comment">//		)</span>
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span><span class="comment">// leads to the type/values matrix below. A run of value columns (V) can</span>
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span><span class="comment">// be moved into the type column if there is no type for any of the values</span>
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span><span class="comment">// in that column (we only move entire columns so that they align properly).</span>
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span><span class="comment">//		matrix        formatted     result</span>
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span><span class="comment">//	                   matrix</span>
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span><span class="comment">//		T  V    -&gt;    T  V     -&gt;   true      there is a T and so the type</span>
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span><span class="comment">//		-  V          -  V          true      column must be kept</span>
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span><span class="comment">//		-  -          -  -          false</span>
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span><span class="comment">//		-  V          V  -          false     V is moved into T column</span>
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span>func keepTypeColumn(specs []ast.Spec) []bool {
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span>	m := make([]bool, len(specs))
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span>
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span>	populate := func(i, j int, keepType bool) {
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span>		if keepType {
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span>			for ; i &lt; j; i++ {
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>				m[i] = true
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span>			}
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span>		}
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span>	}
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span>
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span>	i0 := -1 <span class="comment">// if i0 &gt;= 0 we are in a run and i0 is the start of the run</span>
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>	var keepType bool
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>	for i, s := range specs {
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>		t := s.(*ast.ValueSpec)
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>		if t.Values != nil {
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>			if i0 &lt; 0 {
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>				<span class="comment">// start of a run of ValueSpecs with non-nil Values</span>
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>				i0 = i
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span>				keepType = false
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span>			}
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span>		} else {
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span>			if i0 &gt;= 0 {
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span>				<span class="comment">// end of a run</span>
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>				populate(i0, i, keepType)
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>				i0 = -1
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span>			}
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span>		}
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span>		if t.Type != nil {
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span>			keepType = true
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span>		}
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span>	}
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span>	if i0 &gt;= 0 {
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span>		<span class="comment">// end of a run</span>
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span>		populate(i0, len(specs), keepType)
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span>	}
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span>
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span>	return m
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>}
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>func (p *printer) valueSpec(s *ast.ValueSpec, keepType bool) {
<span id="L1618" class="ln">  1618&nbsp;&nbsp;</span>	p.setComment(s.Doc)
<span id="L1619" class="ln">  1619&nbsp;&nbsp;</span>	p.identList(s.Names, false) <span class="comment">// always present</span>
<span id="L1620" class="ln">  1620&nbsp;&nbsp;</span>	extraTabs := 3
<span id="L1621" class="ln">  1621&nbsp;&nbsp;</span>	if s.Type != nil || keepType {
<span id="L1622" class="ln">  1622&nbsp;&nbsp;</span>		p.print(vtab)
<span id="L1623" class="ln">  1623&nbsp;&nbsp;</span>		extraTabs--
<span id="L1624" class="ln">  1624&nbsp;&nbsp;</span>	}
<span id="L1625" class="ln">  1625&nbsp;&nbsp;</span>	if s.Type != nil {
<span id="L1626" class="ln">  1626&nbsp;&nbsp;</span>		p.expr(s.Type)
<span id="L1627" class="ln">  1627&nbsp;&nbsp;</span>	}
<span id="L1628" class="ln">  1628&nbsp;&nbsp;</span>	if s.Values != nil {
<span id="L1629" class="ln">  1629&nbsp;&nbsp;</span>		p.print(vtab, token.ASSIGN, blank)
<span id="L1630" class="ln">  1630&nbsp;&nbsp;</span>		p.exprList(token.NoPos, s.Values, 1, 0, token.NoPos, false)
<span id="L1631" class="ln">  1631&nbsp;&nbsp;</span>		extraTabs--
<span id="L1632" class="ln">  1632&nbsp;&nbsp;</span>	}
<span id="L1633" class="ln">  1633&nbsp;&nbsp;</span>	if s.Comment != nil {
<span id="L1634" class="ln">  1634&nbsp;&nbsp;</span>		for ; extraTabs &gt; 0; extraTabs-- {
<span id="L1635" class="ln">  1635&nbsp;&nbsp;</span>			p.print(vtab)
<span id="L1636" class="ln">  1636&nbsp;&nbsp;</span>		}
<span id="L1637" class="ln">  1637&nbsp;&nbsp;</span>		p.setComment(s.Comment)
<span id="L1638" class="ln">  1638&nbsp;&nbsp;</span>	}
<span id="L1639" class="ln">  1639&nbsp;&nbsp;</span>}
<span id="L1640" class="ln">  1640&nbsp;&nbsp;</span>
<span id="L1641" class="ln">  1641&nbsp;&nbsp;</span>func sanitizeImportPath(lit *ast.BasicLit) *ast.BasicLit {
<span id="L1642" class="ln">  1642&nbsp;&nbsp;</span>	<span class="comment">// Note: An unmodified AST generated by go/parser will already</span>
<span id="L1643" class="ln">  1643&nbsp;&nbsp;</span>	<span class="comment">// contain a backward- or double-quoted path string that does</span>
<span id="L1644" class="ln">  1644&nbsp;&nbsp;</span>	<span class="comment">// not contain any invalid characters, and most of the work</span>
<span id="L1645" class="ln">  1645&nbsp;&nbsp;</span>	<span class="comment">// here is not needed. However, a modified or generated AST</span>
<span id="L1646" class="ln">  1646&nbsp;&nbsp;</span>	<span class="comment">// may possibly contain non-canonical paths. Do the work in</span>
<span id="L1647" class="ln">  1647&nbsp;&nbsp;</span>	<span class="comment">// all cases since it&#39;s not too hard and not speed-critical.</span>
<span id="L1648" class="ln">  1648&nbsp;&nbsp;</span>
<span id="L1649" class="ln">  1649&nbsp;&nbsp;</span>	<span class="comment">// if we don&#39;t have a proper string, be conservative and return whatever we have</span>
<span id="L1650" class="ln">  1650&nbsp;&nbsp;</span>	if lit.Kind != token.STRING {
<span id="L1651" class="ln">  1651&nbsp;&nbsp;</span>		return lit
<span id="L1652" class="ln">  1652&nbsp;&nbsp;</span>	}
<span id="L1653" class="ln">  1653&nbsp;&nbsp;</span>	s, err := strconv.Unquote(lit.Value)
<span id="L1654" class="ln">  1654&nbsp;&nbsp;</span>	if err != nil {
<span id="L1655" class="ln">  1655&nbsp;&nbsp;</span>		return lit
<span id="L1656" class="ln">  1656&nbsp;&nbsp;</span>	}
<span id="L1657" class="ln">  1657&nbsp;&nbsp;</span>
<span id="L1658" class="ln">  1658&nbsp;&nbsp;</span>	<span class="comment">// if the string is an invalid path, return whatever we have</span>
<span id="L1659" class="ln">  1659&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1660" class="ln">  1660&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;Implementation restriction: A compiler may restrict</span>
<span id="L1661" class="ln">  1661&nbsp;&nbsp;</span>	<span class="comment">// ImportPaths to non-empty strings using only characters belonging</span>
<span id="L1662" class="ln">  1662&nbsp;&nbsp;</span>	<span class="comment">// to Unicode&#39;s L, M, N, P, and S general categories (the Graphic</span>
<span id="L1663" class="ln">  1663&nbsp;&nbsp;</span>	<span class="comment">// characters without spaces) and may also exclude the characters</span>
<span id="L1664" class="ln">  1664&nbsp;&nbsp;</span>	<span class="comment">// !&#34;#$%&amp;&#39;()*,:;&lt;=&gt;?[\]^`{|} and the Unicode replacement character</span>
<span id="L1665" class="ln">  1665&nbsp;&nbsp;</span>	<span class="comment">// U+FFFD.&#34;</span>
<span id="L1666" class="ln">  1666&nbsp;&nbsp;</span>	if s == &#34;&#34; {
<span id="L1667" class="ln">  1667&nbsp;&nbsp;</span>		return lit
<span id="L1668" class="ln">  1668&nbsp;&nbsp;</span>	}
<span id="L1669" class="ln">  1669&nbsp;&nbsp;</span>	const illegalChars = `!&#34;#$%&amp;&#39;()*,:;&lt;=&gt;?[\]^{|}` + &#34;`\uFFFD&#34;
<span id="L1670" class="ln">  1670&nbsp;&nbsp;</span>	for _, r := range s {
<span id="L1671" class="ln">  1671&nbsp;&nbsp;</span>		if !unicode.IsGraphic(r) || unicode.IsSpace(r) || strings.ContainsRune(illegalChars, r) {
<span id="L1672" class="ln">  1672&nbsp;&nbsp;</span>			return lit
<span id="L1673" class="ln">  1673&nbsp;&nbsp;</span>		}
<span id="L1674" class="ln">  1674&nbsp;&nbsp;</span>	}
<span id="L1675" class="ln">  1675&nbsp;&nbsp;</span>
<span id="L1676" class="ln">  1676&nbsp;&nbsp;</span>	<span class="comment">// otherwise, return the double-quoted path</span>
<span id="L1677" class="ln">  1677&nbsp;&nbsp;</span>	s = strconv.Quote(s)
<span id="L1678" class="ln">  1678&nbsp;&nbsp;</span>	if s == lit.Value {
<span id="L1679" class="ln">  1679&nbsp;&nbsp;</span>		return lit <span class="comment">// nothing wrong with lit</span>
<span id="L1680" class="ln">  1680&nbsp;&nbsp;</span>	}
<span id="L1681" class="ln">  1681&nbsp;&nbsp;</span>	return &amp;ast.BasicLit{ValuePos: lit.ValuePos, Kind: token.STRING, Value: s}
<span id="L1682" class="ln">  1682&nbsp;&nbsp;</span>}
<span id="L1683" class="ln">  1683&nbsp;&nbsp;</span>
<span id="L1684" class="ln">  1684&nbsp;&nbsp;</span><span class="comment">// The parameter n is the number of specs in the group. If doIndent is set,</span>
<span id="L1685" class="ln">  1685&nbsp;&nbsp;</span><span class="comment">// multi-line identifier lists in the spec are indented when the first</span>
<span id="L1686" class="ln">  1686&nbsp;&nbsp;</span><span class="comment">// linebreak is encountered.</span>
<span id="L1687" class="ln">  1687&nbsp;&nbsp;</span>func (p *printer) spec(spec ast.Spec, n int, doIndent bool) {
<span id="L1688" class="ln">  1688&nbsp;&nbsp;</span>	switch s := spec.(type) {
<span id="L1689" class="ln">  1689&nbsp;&nbsp;</span>	case *ast.ImportSpec:
<span id="L1690" class="ln">  1690&nbsp;&nbsp;</span>		p.setComment(s.Doc)
<span id="L1691" class="ln">  1691&nbsp;&nbsp;</span>		if s.Name != nil {
<span id="L1692" class="ln">  1692&nbsp;&nbsp;</span>			p.expr(s.Name)
<span id="L1693" class="ln">  1693&nbsp;&nbsp;</span>			p.print(blank)
<span id="L1694" class="ln">  1694&nbsp;&nbsp;</span>		}
<span id="L1695" class="ln">  1695&nbsp;&nbsp;</span>		p.expr(sanitizeImportPath(s.Path))
<span id="L1696" class="ln">  1696&nbsp;&nbsp;</span>		p.setComment(s.Comment)
<span id="L1697" class="ln">  1697&nbsp;&nbsp;</span>		p.setPos(s.EndPos)
<span id="L1698" class="ln">  1698&nbsp;&nbsp;</span>
<span id="L1699" class="ln">  1699&nbsp;&nbsp;</span>	case *ast.ValueSpec:
<span id="L1700" class="ln">  1700&nbsp;&nbsp;</span>		if n != 1 {
<span id="L1701" class="ln">  1701&nbsp;&nbsp;</span>			p.internalError(&#34;expected n = 1; got&#34;, n)
<span id="L1702" class="ln">  1702&nbsp;&nbsp;</span>		}
<span id="L1703" class="ln">  1703&nbsp;&nbsp;</span>		p.setComment(s.Doc)
<span id="L1704" class="ln">  1704&nbsp;&nbsp;</span>		p.identList(s.Names, doIndent) <span class="comment">// always present</span>
<span id="L1705" class="ln">  1705&nbsp;&nbsp;</span>		if s.Type != nil {
<span id="L1706" class="ln">  1706&nbsp;&nbsp;</span>			p.print(blank)
<span id="L1707" class="ln">  1707&nbsp;&nbsp;</span>			p.expr(s.Type)
<span id="L1708" class="ln">  1708&nbsp;&nbsp;</span>		}
<span id="L1709" class="ln">  1709&nbsp;&nbsp;</span>		if s.Values != nil {
<span id="L1710" class="ln">  1710&nbsp;&nbsp;</span>			p.print(blank, token.ASSIGN, blank)
<span id="L1711" class="ln">  1711&nbsp;&nbsp;</span>			p.exprList(token.NoPos, s.Values, 1, 0, token.NoPos, false)
<span id="L1712" class="ln">  1712&nbsp;&nbsp;</span>		}
<span id="L1713" class="ln">  1713&nbsp;&nbsp;</span>		p.setComment(s.Comment)
<span id="L1714" class="ln">  1714&nbsp;&nbsp;</span>
<span id="L1715" class="ln">  1715&nbsp;&nbsp;</span>	case *ast.TypeSpec:
<span id="L1716" class="ln">  1716&nbsp;&nbsp;</span>		p.setComment(s.Doc)
<span id="L1717" class="ln">  1717&nbsp;&nbsp;</span>		p.expr(s.Name)
<span id="L1718" class="ln">  1718&nbsp;&nbsp;</span>		if s.TypeParams != nil {
<span id="L1719" class="ln">  1719&nbsp;&nbsp;</span>			p.parameters(s.TypeParams, typeTParam)
<span id="L1720" class="ln">  1720&nbsp;&nbsp;</span>		}
<span id="L1721" class="ln">  1721&nbsp;&nbsp;</span>		if n == 1 {
<span id="L1722" class="ln">  1722&nbsp;&nbsp;</span>			p.print(blank)
<span id="L1723" class="ln">  1723&nbsp;&nbsp;</span>		} else {
<span id="L1724" class="ln">  1724&nbsp;&nbsp;</span>			p.print(vtab)
<span id="L1725" class="ln">  1725&nbsp;&nbsp;</span>		}
<span id="L1726" class="ln">  1726&nbsp;&nbsp;</span>		if s.Assign.IsValid() {
<span id="L1727" class="ln">  1727&nbsp;&nbsp;</span>			p.print(token.ASSIGN, blank)
<span id="L1728" class="ln">  1728&nbsp;&nbsp;</span>		}
<span id="L1729" class="ln">  1729&nbsp;&nbsp;</span>		p.expr(s.Type)
<span id="L1730" class="ln">  1730&nbsp;&nbsp;</span>		p.setComment(s.Comment)
<span id="L1731" class="ln">  1731&nbsp;&nbsp;</span>
<span id="L1732" class="ln">  1732&nbsp;&nbsp;</span>	default:
<span id="L1733" class="ln">  1733&nbsp;&nbsp;</span>		panic(&#34;unreachable&#34;)
<span id="L1734" class="ln">  1734&nbsp;&nbsp;</span>	}
<span id="L1735" class="ln">  1735&nbsp;&nbsp;</span>}
<span id="L1736" class="ln">  1736&nbsp;&nbsp;</span>
<span id="L1737" class="ln">  1737&nbsp;&nbsp;</span>func (p *printer) genDecl(d *ast.GenDecl) {
<span id="L1738" class="ln">  1738&nbsp;&nbsp;</span>	p.setComment(d.Doc)
<span id="L1739" class="ln">  1739&nbsp;&nbsp;</span>	p.setPos(d.Pos())
<span id="L1740" class="ln">  1740&nbsp;&nbsp;</span>	p.print(d.Tok, blank)
<span id="L1741" class="ln">  1741&nbsp;&nbsp;</span>
<span id="L1742" class="ln">  1742&nbsp;&nbsp;</span>	if d.Lparen.IsValid() || len(d.Specs) != 1 {
<span id="L1743" class="ln">  1743&nbsp;&nbsp;</span>		<span class="comment">// group of parenthesized declarations</span>
<span id="L1744" class="ln">  1744&nbsp;&nbsp;</span>		p.setPos(d.Lparen)
<span id="L1745" class="ln">  1745&nbsp;&nbsp;</span>		p.print(token.LPAREN)
<span id="L1746" class="ln">  1746&nbsp;&nbsp;</span>		if n := len(d.Specs); n &gt; 0 {
<span id="L1747" class="ln">  1747&nbsp;&nbsp;</span>			p.print(indent, formfeed)
<span id="L1748" class="ln">  1748&nbsp;&nbsp;</span>			if n &gt; 1 &amp;&amp; (d.Tok == token.CONST || d.Tok == token.VAR) {
<span id="L1749" class="ln">  1749&nbsp;&nbsp;</span>				<span class="comment">// two or more grouped const/var declarations:</span>
<span id="L1750" class="ln">  1750&nbsp;&nbsp;</span>				<span class="comment">// determine if the type column must be kept</span>
<span id="L1751" class="ln">  1751&nbsp;&nbsp;</span>				keepType := keepTypeColumn(d.Specs)
<span id="L1752" class="ln">  1752&nbsp;&nbsp;</span>				var line int
<span id="L1753" class="ln">  1753&nbsp;&nbsp;</span>				for i, s := range d.Specs {
<span id="L1754" class="ln">  1754&nbsp;&nbsp;</span>					if i &gt; 0 {
<span id="L1755" class="ln">  1755&nbsp;&nbsp;</span>						p.linebreak(p.lineFor(s.Pos()), 1, ignore, p.linesFrom(line) &gt; 0)
<span id="L1756" class="ln">  1756&nbsp;&nbsp;</span>					}
<span id="L1757" class="ln">  1757&nbsp;&nbsp;</span>					p.recordLine(&amp;line)
<span id="L1758" class="ln">  1758&nbsp;&nbsp;</span>					p.valueSpec(s.(*ast.ValueSpec), keepType[i])
<span id="L1759" class="ln">  1759&nbsp;&nbsp;</span>				}
<span id="L1760" class="ln">  1760&nbsp;&nbsp;</span>			} else {
<span id="L1761" class="ln">  1761&nbsp;&nbsp;</span>				var line int
<span id="L1762" class="ln">  1762&nbsp;&nbsp;</span>				for i, s := range d.Specs {
<span id="L1763" class="ln">  1763&nbsp;&nbsp;</span>					if i &gt; 0 {
<span id="L1764" class="ln">  1764&nbsp;&nbsp;</span>						p.linebreak(p.lineFor(s.Pos()), 1, ignore, p.linesFrom(line) &gt; 0)
<span id="L1765" class="ln">  1765&nbsp;&nbsp;</span>					}
<span id="L1766" class="ln">  1766&nbsp;&nbsp;</span>					p.recordLine(&amp;line)
<span id="L1767" class="ln">  1767&nbsp;&nbsp;</span>					p.spec(s, n, false)
<span id="L1768" class="ln">  1768&nbsp;&nbsp;</span>				}
<span id="L1769" class="ln">  1769&nbsp;&nbsp;</span>			}
<span id="L1770" class="ln">  1770&nbsp;&nbsp;</span>			p.print(unindent, formfeed)
<span id="L1771" class="ln">  1771&nbsp;&nbsp;</span>		}
<span id="L1772" class="ln">  1772&nbsp;&nbsp;</span>		p.setPos(d.Rparen)
<span id="L1773" class="ln">  1773&nbsp;&nbsp;</span>		p.print(token.RPAREN)
<span id="L1774" class="ln">  1774&nbsp;&nbsp;</span>
<span id="L1775" class="ln">  1775&nbsp;&nbsp;</span>	} else if len(d.Specs) &gt; 0 {
<span id="L1776" class="ln">  1776&nbsp;&nbsp;</span>		<span class="comment">// single declaration</span>
<span id="L1777" class="ln">  1777&nbsp;&nbsp;</span>		p.spec(d.Specs[0], 1, true)
<span id="L1778" class="ln">  1778&nbsp;&nbsp;</span>	}
<span id="L1779" class="ln">  1779&nbsp;&nbsp;</span>}
<span id="L1780" class="ln">  1780&nbsp;&nbsp;</span>
<span id="L1781" class="ln">  1781&nbsp;&nbsp;</span><span class="comment">// sizeCounter is an io.Writer which counts the number of bytes written,</span>
<span id="L1782" class="ln">  1782&nbsp;&nbsp;</span><span class="comment">// as well as whether a newline character was seen.</span>
<span id="L1783" class="ln">  1783&nbsp;&nbsp;</span>type sizeCounter struct {
<span id="L1784" class="ln">  1784&nbsp;&nbsp;</span>	hasNewline bool
<span id="L1785" class="ln">  1785&nbsp;&nbsp;</span>	size       int
<span id="L1786" class="ln">  1786&nbsp;&nbsp;</span>}
<span id="L1787" class="ln">  1787&nbsp;&nbsp;</span>
<span id="L1788" class="ln">  1788&nbsp;&nbsp;</span>func (c *sizeCounter) Write(p []byte) (int, error) {
<span id="L1789" class="ln">  1789&nbsp;&nbsp;</span>	if !c.hasNewline {
<span id="L1790" class="ln">  1790&nbsp;&nbsp;</span>		for _, b := range p {
<span id="L1791" class="ln">  1791&nbsp;&nbsp;</span>			if b == &#39;\n&#39; || b == &#39;\f&#39; {
<span id="L1792" class="ln">  1792&nbsp;&nbsp;</span>				c.hasNewline = true
<span id="L1793" class="ln">  1793&nbsp;&nbsp;</span>				break
<span id="L1794" class="ln">  1794&nbsp;&nbsp;</span>			}
<span id="L1795" class="ln">  1795&nbsp;&nbsp;</span>		}
<span id="L1796" class="ln">  1796&nbsp;&nbsp;</span>	}
<span id="L1797" class="ln">  1797&nbsp;&nbsp;</span>	c.size += len(p)
<span id="L1798" class="ln">  1798&nbsp;&nbsp;</span>	return len(p), nil
<span id="L1799" class="ln">  1799&nbsp;&nbsp;</span>}
<span id="L1800" class="ln">  1800&nbsp;&nbsp;</span>
<span id="L1801" class="ln">  1801&nbsp;&nbsp;</span><span class="comment">// nodeSize determines the size of n in chars after formatting.</span>
<span id="L1802" class="ln">  1802&nbsp;&nbsp;</span><span class="comment">// The result is &lt;= maxSize if the node fits on one line with at</span>
<span id="L1803" class="ln">  1803&nbsp;&nbsp;</span><span class="comment">// most maxSize chars and the formatted output doesn&#39;t contain</span>
<span id="L1804" class="ln">  1804&nbsp;&nbsp;</span><span class="comment">// any control chars. Otherwise, the result is &gt; maxSize.</span>
<span id="L1805" class="ln">  1805&nbsp;&nbsp;</span>func (p *printer) nodeSize(n ast.Node, maxSize int) (size int) {
<span id="L1806" class="ln">  1806&nbsp;&nbsp;</span>	<span class="comment">// nodeSize invokes the printer, which may invoke nodeSize</span>
<span id="L1807" class="ln">  1807&nbsp;&nbsp;</span>	<span class="comment">// recursively. For deep composite literal nests, this can</span>
<span id="L1808" class="ln">  1808&nbsp;&nbsp;</span>	<span class="comment">// lead to an exponential algorithm. Remember previous</span>
<span id="L1809" class="ln">  1809&nbsp;&nbsp;</span>	<span class="comment">// results to prune the recursion (was issue 1628).</span>
<span id="L1810" class="ln">  1810&nbsp;&nbsp;</span>	if size, found := p.nodeSizes[n]; found {
<span id="L1811" class="ln">  1811&nbsp;&nbsp;</span>		return size
<span id="L1812" class="ln">  1812&nbsp;&nbsp;</span>	}
<span id="L1813" class="ln">  1813&nbsp;&nbsp;</span>
<span id="L1814" class="ln">  1814&nbsp;&nbsp;</span>	size = maxSize + 1 <span class="comment">// assume n doesn&#39;t fit</span>
<span id="L1815" class="ln">  1815&nbsp;&nbsp;</span>	p.nodeSizes[n] = size
<span id="L1816" class="ln">  1816&nbsp;&nbsp;</span>
<span id="L1817" class="ln">  1817&nbsp;&nbsp;</span>	<span class="comment">// nodeSize computation must be independent of particular</span>
<span id="L1818" class="ln">  1818&nbsp;&nbsp;</span>	<span class="comment">// style so that we always get the same decision; print</span>
<span id="L1819" class="ln">  1819&nbsp;&nbsp;</span>	<span class="comment">// in RawFormat</span>
<span id="L1820" class="ln">  1820&nbsp;&nbsp;</span>	cfg := Config{Mode: RawFormat}
<span id="L1821" class="ln">  1821&nbsp;&nbsp;</span>	var counter sizeCounter
<span id="L1822" class="ln">  1822&nbsp;&nbsp;</span>	if err := cfg.fprint(&amp;counter, p.fset, n, p.nodeSizes); err != nil {
<span id="L1823" class="ln">  1823&nbsp;&nbsp;</span>		return
<span id="L1824" class="ln">  1824&nbsp;&nbsp;</span>	}
<span id="L1825" class="ln">  1825&nbsp;&nbsp;</span>	if counter.size &lt;= maxSize &amp;&amp; !counter.hasNewline {
<span id="L1826" class="ln">  1826&nbsp;&nbsp;</span>		<span class="comment">// n fits in a single line</span>
<span id="L1827" class="ln">  1827&nbsp;&nbsp;</span>		size = counter.size
<span id="L1828" class="ln">  1828&nbsp;&nbsp;</span>		p.nodeSizes[n] = size
<span id="L1829" class="ln">  1829&nbsp;&nbsp;</span>	}
<span id="L1830" class="ln">  1830&nbsp;&nbsp;</span>	return
<span id="L1831" class="ln">  1831&nbsp;&nbsp;</span>}
<span id="L1832" class="ln">  1832&nbsp;&nbsp;</span>
<span id="L1833" class="ln">  1833&nbsp;&nbsp;</span><span class="comment">// numLines returns the number of lines spanned by node n in the original source.</span>
<span id="L1834" class="ln">  1834&nbsp;&nbsp;</span>func (p *printer) numLines(n ast.Node) int {
<span id="L1835" class="ln">  1835&nbsp;&nbsp;</span>	if from := n.Pos(); from.IsValid() {
<span id="L1836" class="ln">  1836&nbsp;&nbsp;</span>		if to := n.End(); to.IsValid() {
<span id="L1837" class="ln">  1837&nbsp;&nbsp;</span>			return p.lineFor(to) - p.lineFor(from) + 1
<span id="L1838" class="ln">  1838&nbsp;&nbsp;</span>		}
<span id="L1839" class="ln">  1839&nbsp;&nbsp;</span>	}
<span id="L1840" class="ln">  1840&nbsp;&nbsp;</span>	return infinity
<span id="L1841" class="ln">  1841&nbsp;&nbsp;</span>}
<span id="L1842" class="ln">  1842&nbsp;&nbsp;</span>
<span id="L1843" class="ln">  1843&nbsp;&nbsp;</span><span class="comment">// bodySize is like nodeSize but it is specialized for *ast.BlockStmt&#39;s.</span>
<span id="L1844" class="ln">  1844&nbsp;&nbsp;</span>func (p *printer) bodySize(b *ast.BlockStmt, maxSize int) int {
<span id="L1845" class="ln">  1845&nbsp;&nbsp;</span>	pos1 := b.Pos()
<span id="L1846" class="ln">  1846&nbsp;&nbsp;</span>	pos2 := b.Rbrace
<span id="L1847" class="ln">  1847&nbsp;&nbsp;</span>	if pos1.IsValid() &amp;&amp; pos2.IsValid() &amp;&amp; p.lineFor(pos1) != p.lineFor(pos2) {
<span id="L1848" class="ln">  1848&nbsp;&nbsp;</span>		<span class="comment">// opening and closing brace are on different lines - don&#39;t make it a one-liner</span>
<span id="L1849" class="ln">  1849&nbsp;&nbsp;</span>		return maxSize + 1
<span id="L1850" class="ln">  1850&nbsp;&nbsp;</span>	}
<span id="L1851" class="ln">  1851&nbsp;&nbsp;</span>	if len(b.List) &gt; 5 {
<span id="L1852" class="ln">  1852&nbsp;&nbsp;</span>		<span class="comment">// too many statements - don&#39;t make it a one-liner</span>
<span id="L1853" class="ln">  1853&nbsp;&nbsp;</span>		return maxSize + 1
<span id="L1854" class="ln">  1854&nbsp;&nbsp;</span>	}
<span id="L1855" class="ln">  1855&nbsp;&nbsp;</span>	<span class="comment">// otherwise, estimate body size</span>
<span id="L1856" class="ln">  1856&nbsp;&nbsp;</span>	bodySize := p.commentSizeBefore(p.posFor(pos2))
<span id="L1857" class="ln">  1857&nbsp;&nbsp;</span>	for i, s := range b.List {
<span id="L1858" class="ln">  1858&nbsp;&nbsp;</span>		if bodySize &gt; maxSize {
<span id="L1859" class="ln">  1859&nbsp;&nbsp;</span>			break <span class="comment">// no need to continue</span>
<span id="L1860" class="ln">  1860&nbsp;&nbsp;</span>		}
<span id="L1861" class="ln">  1861&nbsp;&nbsp;</span>		if i &gt; 0 {
<span id="L1862" class="ln">  1862&nbsp;&nbsp;</span>			bodySize += 2 <span class="comment">// space for a semicolon and blank</span>
<span id="L1863" class="ln">  1863&nbsp;&nbsp;</span>		}
<span id="L1864" class="ln">  1864&nbsp;&nbsp;</span>		bodySize += p.nodeSize(s, maxSize)
<span id="L1865" class="ln">  1865&nbsp;&nbsp;</span>	}
<span id="L1866" class="ln">  1866&nbsp;&nbsp;</span>	return bodySize
<span id="L1867" class="ln">  1867&nbsp;&nbsp;</span>}
<span id="L1868" class="ln">  1868&nbsp;&nbsp;</span>
<span id="L1869" class="ln">  1869&nbsp;&nbsp;</span><span class="comment">// funcBody prints a function body following a function header of given headerSize.</span>
<span id="L1870" class="ln">  1870&nbsp;&nbsp;</span><span class="comment">// If the header&#39;s and block&#39;s size are &#34;small enough&#34; and the block is &#34;simple enough&#34;,</span>
<span id="L1871" class="ln">  1871&nbsp;&nbsp;</span><span class="comment">// the block is printed on the current line, without line breaks, spaced from the header</span>
<span id="L1872" class="ln">  1872&nbsp;&nbsp;</span><span class="comment">// by sep. Otherwise the block&#39;s opening &#34;{&#34; is printed on the current line, followed by</span>
<span id="L1873" class="ln">  1873&nbsp;&nbsp;</span><span class="comment">// lines for the block&#39;s statements and its closing &#34;}&#34;.</span>
<span id="L1874" class="ln">  1874&nbsp;&nbsp;</span>func (p *printer) funcBody(headerSize int, sep whiteSpace, b *ast.BlockStmt) {
<span id="L1875" class="ln">  1875&nbsp;&nbsp;</span>	if b == nil {
<span id="L1876" class="ln">  1876&nbsp;&nbsp;</span>		return
<span id="L1877" class="ln">  1877&nbsp;&nbsp;</span>	}
<span id="L1878" class="ln">  1878&nbsp;&nbsp;</span>
<span id="L1879" class="ln">  1879&nbsp;&nbsp;</span>	<span class="comment">// save/restore composite literal nesting level</span>
<span id="L1880" class="ln">  1880&nbsp;&nbsp;</span>	defer func(level int) {
<span id="L1881" class="ln">  1881&nbsp;&nbsp;</span>		p.level = level
<span id="L1882" class="ln">  1882&nbsp;&nbsp;</span>	}(p.level)
<span id="L1883" class="ln">  1883&nbsp;&nbsp;</span>	p.level = 0
<span id="L1884" class="ln">  1884&nbsp;&nbsp;</span>
<span id="L1885" class="ln">  1885&nbsp;&nbsp;</span>	const maxSize = 100
<span id="L1886" class="ln">  1886&nbsp;&nbsp;</span>	if headerSize+p.bodySize(b, maxSize) &lt;= maxSize {
<span id="L1887" class="ln">  1887&nbsp;&nbsp;</span>		p.print(sep)
<span id="L1888" class="ln">  1888&nbsp;&nbsp;</span>		p.setPos(b.Lbrace)
<span id="L1889" class="ln">  1889&nbsp;&nbsp;</span>		p.print(token.LBRACE)
<span id="L1890" class="ln">  1890&nbsp;&nbsp;</span>		if len(b.List) &gt; 0 {
<span id="L1891" class="ln">  1891&nbsp;&nbsp;</span>			p.print(blank)
<span id="L1892" class="ln">  1892&nbsp;&nbsp;</span>			for i, s := range b.List {
<span id="L1893" class="ln">  1893&nbsp;&nbsp;</span>				if i &gt; 0 {
<span id="L1894" class="ln">  1894&nbsp;&nbsp;</span>					p.print(token.SEMICOLON, blank)
<span id="L1895" class="ln">  1895&nbsp;&nbsp;</span>				}
<span id="L1896" class="ln">  1896&nbsp;&nbsp;</span>				p.stmt(s, i == len(b.List)-1)
<span id="L1897" class="ln">  1897&nbsp;&nbsp;</span>			}
<span id="L1898" class="ln">  1898&nbsp;&nbsp;</span>			p.print(blank)
<span id="L1899" class="ln">  1899&nbsp;&nbsp;</span>		}
<span id="L1900" class="ln">  1900&nbsp;&nbsp;</span>		p.print(noExtraLinebreak)
<span id="L1901" class="ln">  1901&nbsp;&nbsp;</span>		p.setPos(b.Rbrace)
<span id="L1902" class="ln">  1902&nbsp;&nbsp;</span>		p.print(token.RBRACE, noExtraLinebreak)
<span id="L1903" class="ln">  1903&nbsp;&nbsp;</span>		return
<span id="L1904" class="ln">  1904&nbsp;&nbsp;</span>	}
<span id="L1905" class="ln">  1905&nbsp;&nbsp;</span>
<span id="L1906" class="ln">  1906&nbsp;&nbsp;</span>	if sep != ignore {
<span id="L1907" class="ln">  1907&nbsp;&nbsp;</span>		p.print(blank) <span class="comment">// always use blank</span>
<span id="L1908" class="ln">  1908&nbsp;&nbsp;</span>	}
<span id="L1909" class="ln">  1909&nbsp;&nbsp;</span>	p.block(b, 1)
<span id="L1910" class="ln">  1910&nbsp;&nbsp;</span>}
<span id="L1911" class="ln">  1911&nbsp;&nbsp;</span>
<span id="L1912" class="ln">  1912&nbsp;&nbsp;</span><span class="comment">// distanceFrom returns the column difference between p.out (the current output</span>
<span id="L1913" class="ln">  1913&nbsp;&nbsp;</span><span class="comment">// position) and startOutCol. If the start position is on a different line from</span>
<span id="L1914" class="ln">  1914&nbsp;&nbsp;</span><span class="comment">// the current position (or either is unknown), the result is infinity.</span>
<span id="L1915" class="ln">  1915&nbsp;&nbsp;</span>func (p *printer) distanceFrom(startPos token.Pos, startOutCol int) int {
<span id="L1916" class="ln">  1916&nbsp;&nbsp;</span>	if startPos.IsValid() &amp;&amp; p.pos.IsValid() &amp;&amp; p.posFor(startPos).Line == p.pos.Line {
<span id="L1917" class="ln">  1917&nbsp;&nbsp;</span>		return p.out.Column - startOutCol
<span id="L1918" class="ln">  1918&nbsp;&nbsp;</span>	}
<span id="L1919" class="ln">  1919&nbsp;&nbsp;</span>	return infinity
<span id="L1920" class="ln">  1920&nbsp;&nbsp;</span>}
<span id="L1921" class="ln">  1921&nbsp;&nbsp;</span>
<span id="L1922" class="ln">  1922&nbsp;&nbsp;</span>func (p *printer) funcDecl(d *ast.FuncDecl) {
<span id="L1923" class="ln">  1923&nbsp;&nbsp;</span>	p.setComment(d.Doc)
<span id="L1924" class="ln">  1924&nbsp;&nbsp;</span>	p.setPos(d.Pos())
<span id="L1925" class="ln">  1925&nbsp;&nbsp;</span>	p.print(token.FUNC, blank)
<span id="L1926" class="ln">  1926&nbsp;&nbsp;</span>	<span class="comment">// We have to save startCol only after emitting FUNC; otherwise it can be on a</span>
<span id="L1927" class="ln">  1927&nbsp;&nbsp;</span>	<span class="comment">// different line (all whitespace preceding the FUNC is emitted only when the</span>
<span id="L1928" class="ln">  1928&nbsp;&nbsp;</span>	<span class="comment">// FUNC is emitted).</span>
<span id="L1929" class="ln">  1929&nbsp;&nbsp;</span>	startCol := p.out.Column - len(&#34;func &#34;)
<span id="L1930" class="ln">  1930&nbsp;&nbsp;</span>	if d.Recv != nil {
<span id="L1931" class="ln">  1931&nbsp;&nbsp;</span>		p.parameters(d.Recv, funcParam) <span class="comment">// method: print receiver</span>
<span id="L1932" class="ln">  1932&nbsp;&nbsp;</span>		p.print(blank)
<span id="L1933" class="ln">  1933&nbsp;&nbsp;</span>	}
<span id="L1934" class="ln">  1934&nbsp;&nbsp;</span>	p.expr(d.Name)
<span id="L1935" class="ln">  1935&nbsp;&nbsp;</span>	p.signature(d.Type)
<span id="L1936" class="ln">  1936&nbsp;&nbsp;</span>	p.funcBody(p.distanceFrom(d.Pos(), startCol), vtab, d.Body)
<span id="L1937" class="ln">  1937&nbsp;&nbsp;</span>}
<span id="L1938" class="ln">  1938&nbsp;&nbsp;</span>
<span id="L1939" class="ln">  1939&nbsp;&nbsp;</span>func (p *printer) decl(decl ast.Decl) {
<span id="L1940" class="ln">  1940&nbsp;&nbsp;</span>	switch d := decl.(type) {
<span id="L1941" class="ln">  1941&nbsp;&nbsp;</span>	case *ast.BadDecl:
<span id="L1942" class="ln">  1942&nbsp;&nbsp;</span>		p.setPos(d.Pos())
<span id="L1943" class="ln">  1943&nbsp;&nbsp;</span>		p.print(&#34;BadDecl&#34;)
<span id="L1944" class="ln">  1944&nbsp;&nbsp;</span>	case *ast.GenDecl:
<span id="L1945" class="ln">  1945&nbsp;&nbsp;</span>		p.genDecl(d)
<span id="L1946" class="ln">  1946&nbsp;&nbsp;</span>	case *ast.FuncDecl:
<span id="L1947" class="ln">  1947&nbsp;&nbsp;</span>		p.funcDecl(d)
<span id="L1948" class="ln">  1948&nbsp;&nbsp;</span>	default:
<span id="L1949" class="ln">  1949&nbsp;&nbsp;</span>		panic(&#34;unreachable&#34;)
<span id="L1950" class="ln">  1950&nbsp;&nbsp;</span>	}
<span id="L1951" class="ln">  1951&nbsp;&nbsp;</span>}
<span id="L1952" class="ln">  1952&nbsp;&nbsp;</span>
<span id="L1953" class="ln">  1953&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L1954" class="ln">  1954&nbsp;&nbsp;</span><span class="comment">// Files</span>
<span id="L1955" class="ln">  1955&nbsp;&nbsp;</span>
<span id="L1956" class="ln">  1956&nbsp;&nbsp;</span>func declToken(decl ast.Decl) (tok token.Token) {
<span id="L1957" class="ln">  1957&nbsp;&nbsp;</span>	tok = token.ILLEGAL
<span id="L1958" class="ln">  1958&nbsp;&nbsp;</span>	switch d := decl.(type) {
<span id="L1959" class="ln">  1959&nbsp;&nbsp;</span>	case *ast.GenDecl:
<span id="L1960" class="ln">  1960&nbsp;&nbsp;</span>		tok = d.Tok
<span id="L1961" class="ln">  1961&nbsp;&nbsp;</span>	case *ast.FuncDecl:
<span id="L1962" class="ln">  1962&nbsp;&nbsp;</span>		tok = token.FUNC
<span id="L1963" class="ln">  1963&nbsp;&nbsp;</span>	}
<span id="L1964" class="ln">  1964&nbsp;&nbsp;</span>	return
<span id="L1965" class="ln">  1965&nbsp;&nbsp;</span>}
<span id="L1966" class="ln">  1966&nbsp;&nbsp;</span>
<span id="L1967" class="ln">  1967&nbsp;&nbsp;</span>func (p *printer) declList(list []ast.Decl) {
<span id="L1968" class="ln">  1968&nbsp;&nbsp;</span>	tok := token.ILLEGAL
<span id="L1969" class="ln">  1969&nbsp;&nbsp;</span>	for _, d := range list {
<span id="L1970" class="ln">  1970&nbsp;&nbsp;</span>		prev := tok
<span id="L1971" class="ln">  1971&nbsp;&nbsp;</span>		tok = declToken(d)
<span id="L1972" class="ln">  1972&nbsp;&nbsp;</span>		<span class="comment">// If the declaration token changed (e.g., from CONST to TYPE)</span>
<span id="L1973" class="ln">  1973&nbsp;&nbsp;</span>		<span class="comment">// or the next declaration has documentation associated with it,</span>
<span id="L1974" class="ln">  1974&nbsp;&nbsp;</span>		<span class="comment">// print an empty line between top-level declarations.</span>
<span id="L1975" class="ln">  1975&nbsp;&nbsp;</span>		<span class="comment">// (because p.linebreak is called with the position of d, which</span>
<span id="L1976" class="ln">  1976&nbsp;&nbsp;</span>		<span class="comment">// is past any documentation, the minimum requirement is satisfied</span>
<span id="L1977" class="ln">  1977&nbsp;&nbsp;</span>		<span class="comment">// even w/o the extra getDoc(d) nil-check - leave it in case the</span>
<span id="L1978" class="ln">  1978&nbsp;&nbsp;</span>		<span class="comment">// linebreak logic improves - there&#39;s already a TODO).</span>
<span id="L1979" class="ln">  1979&nbsp;&nbsp;</span>		if len(p.output) &gt; 0 {
<span id="L1980" class="ln">  1980&nbsp;&nbsp;</span>			<span class="comment">// only print line break if we are not at the beginning of the output</span>
<span id="L1981" class="ln">  1981&nbsp;&nbsp;</span>			<span class="comment">// (i.e., we are not printing only a partial program)</span>
<span id="L1982" class="ln">  1982&nbsp;&nbsp;</span>			min := 1
<span id="L1983" class="ln">  1983&nbsp;&nbsp;</span>			if prev != tok || getDoc(d) != nil {
<span id="L1984" class="ln">  1984&nbsp;&nbsp;</span>				min = 2
<span id="L1985" class="ln">  1985&nbsp;&nbsp;</span>			}
<span id="L1986" class="ln">  1986&nbsp;&nbsp;</span>			<span class="comment">// start a new section if the next declaration is a function</span>
<span id="L1987" class="ln">  1987&nbsp;&nbsp;</span>			<span class="comment">// that spans multiple lines (see also issue #19544)</span>
<span id="L1988" class="ln">  1988&nbsp;&nbsp;</span>			p.linebreak(p.lineFor(d.Pos()), min, ignore, tok == token.FUNC &amp;&amp; p.numLines(d) &gt; 1)
<span id="L1989" class="ln">  1989&nbsp;&nbsp;</span>		}
<span id="L1990" class="ln">  1990&nbsp;&nbsp;</span>		p.decl(d)
<span id="L1991" class="ln">  1991&nbsp;&nbsp;</span>	}
<span id="L1992" class="ln">  1992&nbsp;&nbsp;</span>}
<span id="L1993" class="ln">  1993&nbsp;&nbsp;</span>
<span id="L1994" class="ln">  1994&nbsp;&nbsp;</span>func (p *printer) file(src *ast.File) {
<span id="L1995" class="ln">  1995&nbsp;&nbsp;</span>	p.setComment(src.Doc)
<span id="L1996" class="ln">  1996&nbsp;&nbsp;</span>	p.setPos(src.Pos())
<span id="L1997" class="ln">  1997&nbsp;&nbsp;</span>	p.print(token.PACKAGE, blank)
<span id="L1998" class="ln">  1998&nbsp;&nbsp;</span>	p.expr(src.Name)
<span id="L1999" class="ln">  1999&nbsp;&nbsp;</span>	p.declList(src.Decls)
<span id="L2000" class="ln">  2000&nbsp;&nbsp;</span>	p.print(newline)
<span id="L2001" class="ln">  2001&nbsp;&nbsp;</span>}
<span id="L2002" class="ln">  2002&nbsp;&nbsp;</span>
</pre><p><a href="nodes.go?m=text">View as plain text</a></p>

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
