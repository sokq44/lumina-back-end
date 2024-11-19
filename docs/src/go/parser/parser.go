<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/parser/parser.go - Go Documentation Server</title>

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
<a href="parser.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/parser">parser</a>/<span class="text-muted">parser.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/parser">go/parser</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package parser implements a parser for Go source files. Input may be</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// provided in a variety of forms (see the various Parse* functions); the</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// output is an abstract syntax tree (AST) representing the Go source. The</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// parser is invoked through one of the Parse* functions.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// The parser accepts a larger language than is syntactically permitted by</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// the Go spec, for simplicity, and for improved robustness in the presence</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// of syntax errors. For instance, in method declarations, the receiver is</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// treated like an ordinary parameter list and thus may contain multiple</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// entries where the spec permits exactly one. Consequently, the corresponding</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// field in the AST (ast.FuncDecl.Recv) field is not restricted to one entry.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>package parser
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>import (
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;go/build/constraint&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;go/internal/typeparams&#34;
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	&#34;go/scanner&#34;
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// The parser structure holds the parser&#39;s internal state.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>type parser struct {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	file    *token.File
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	errors  scanner.ErrorList
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	scanner scanner.Scanner
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// Tracing/debugging</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	mode   Mode <span class="comment">// parsing mode</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	trace  bool <span class="comment">// == (mode&amp;Trace != 0)</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	indent int  <span class="comment">// indentation used for tracing output</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// Comments</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	comments    []*ast.CommentGroup
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	leadComment *ast.CommentGroup <span class="comment">// last lead comment</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	lineComment *ast.CommentGroup <span class="comment">// last line comment</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	top         bool              <span class="comment">// in top of file (before package clause)</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	goVersion   string            <span class="comment">// minimum Go version found in //go:build comment</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// Next token</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	pos token.Pos   <span class="comment">// token position</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	tok token.Token <span class="comment">// one token look-ahead</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	lit string      <span class="comment">// token literal</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// Error recovery</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// (used to limit the number of calls to parser.advance</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// w/o making scanning progress - avoids potential endless</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// loops across multiple parser functions during error recovery)</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	syncPos token.Pos <span class="comment">// last synchronization position</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	syncCnt int       <span class="comment">// number of parser.advance calls without progress</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// Non-syntactic parser control</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	exprLev int  <span class="comment">// &lt; 0: in control clause, &gt;= 0: in expression</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	inRhs   bool <span class="comment">// if set, the parser is parsing a rhs expression</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	imports []*ast.ImportSpec <span class="comment">// list of imports</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// nestLev is used to track and limit the recursion depth</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// during parsing.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	nestLev int
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>func (p *parser) init(fset *token.FileSet, filename string, src []byte, mode Mode) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	p.file = fset.AddFile(filename, -1, len(src))
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	eh := func(pos token.Position, msg string) { p.errors.Add(pos, msg) }
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	p.scanner.Init(p.file, src, eh, scanner.ScanComments)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	p.top = true
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	p.mode = mode
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	p.trace = mode&amp;Trace != 0 <span class="comment">// for convenience (p.trace is used frequently)</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	p.next()
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// Parsing support</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func (p *parser) printTrace(a ...any) {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	const dots = &#34;. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . &#34;
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	const n = len(dots)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	pos := p.file.Position(p.pos)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	fmt.Printf(&#34;%5d:%3d: &#34;, pos.Line, pos.Column)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	i := 2 * p.indent
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	for i &gt; n {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		fmt.Print(dots)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		i -= n
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// i &lt;= n</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	fmt.Print(dots[0:i])
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	fmt.Println(a...)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>func trace(p *parser, msg string) *parser {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	p.printTrace(msg, &#34;(&#34;)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	p.indent++
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	return p
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// Usage pattern: defer un(trace(p, &#34;...&#34;))</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func un(p *parser) {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	p.indent--
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	p.printTrace(&#34;)&#34;)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// maxNestLev is the deepest we&#39;re willing to recurse during parsing</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>const maxNestLev int = 1e5
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func incNestLev(p *parser) *parser {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	p.nestLev++
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if p.nestLev &gt; maxNestLev {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		p.error(p.pos, &#34;exceeded max nesting depth&#34;)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		panic(bailout{})
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	return p
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// decNestLev is used to track nesting depth during parsing to prevent stack exhaustion.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// It is used along with incNestLev in a similar fashion to how un and trace are used.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func decNestLev(p *parser) {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	p.nestLev--
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// Advance to the next token.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>func (p *parser) next0() {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// Because of one-token look-ahead, print the previous token</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// when tracing as it provides a more readable output. The</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// very first token (!p.pos.IsValid()) is not initialized</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// (it is token.ILLEGAL), so don&#39;t print it.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	if p.trace &amp;&amp; p.pos.IsValid() {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		s := p.tok.String()
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		switch {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		case p.tok.IsLiteral():
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			p.printTrace(s, p.lit)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		case p.tok.IsOperator(), p.tok.IsKeyword():
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			p.printTrace(&#34;\&#34;&#34; + s + &#34;\&#34;&#34;)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		default:
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			p.printTrace(s)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	for {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		p.pos, p.tok, p.lit = p.scanner.Scan()
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		if p.tok == token.COMMENT {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			if p.top &amp;&amp; strings.HasPrefix(p.lit, &#34;//go:build&#34;) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>				if x, err := constraint.Parse(p.lit); err == nil {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>					p.goVersion = constraint.GoVersion(x)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>				}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			if p.mode&amp;ParseComments == 0 {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>				continue
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		} else {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			<span class="comment">// Found a non-comment; top of file is over.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			p.top = false
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		break
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">// Consume a comment and return it and the line on which it ends.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>func (p *parser) consumeComment() (comment *ast.Comment, endline int) {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// /*-style comments may end on a different line than where they start.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// Scan the comment for &#39;\n&#39; chars and adjust endline accordingly.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	endline = p.file.Line(p.pos)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if p.lit[1] == &#39;*&#39; {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		<span class="comment">// don&#39;t use range here - no need to decode Unicode code points</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		for i := 0; i &lt; len(p.lit); i++ {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			if p.lit[i] == &#39;\n&#39; {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>				endline++
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	comment = &amp;ast.Comment{Slash: p.pos, Text: p.lit}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	p.next0()
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	return
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span><span class="comment">// Consume a group of adjacent comments, add it to the parser&#39;s</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// comments list, and return it together with the line at which</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// the last comment in the group ends. A non-comment token or n</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// empty lines terminate a comment group.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>func (p *parser) consumeCommentGroup(n int) (comments *ast.CommentGroup, endline int) {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	var list []*ast.Comment
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	endline = p.file.Line(p.pos)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	for p.tok == token.COMMENT &amp;&amp; p.file.Line(p.pos) &lt;= endline+n {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		var comment *ast.Comment
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		comment, endline = p.consumeComment()
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		list = append(list, comment)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// add comment group to the comments list</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	comments = &amp;ast.CommentGroup{List: list}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	p.comments = append(p.comments, comments)
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	return
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// Advance to the next non-comment token. In the process, collect</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// any comment groups encountered, and remember the last lead and</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// line comments.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// A lead comment is a comment group that starts and ends in a</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// line without any other tokens and that is followed by a non-comment</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// token on the line immediately after the comment group.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// A line comment is a comment group that follows a non-comment</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// token on the same line, and that has no tokens after it on the line</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// where it ends.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// Lead and line comments may be considered documentation that is</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// stored in the AST.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>func (p *parser) next() {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	p.leadComment = nil
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	p.lineComment = nil
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	prev := p.pos
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	p.next0()
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	if p.tok == token.COMMENT {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		var comment *ast.CommentGroup
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		var endline int
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		if p.file.Line(p.pos) == p.file.Line(prev) {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			<span class="comment">// The comment is on same line as the previous token; it</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			<span class="comment">// cannot be a lead comment but may be a line comment.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			comment, endline = p.consumeCommentGroup(0)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			if p.file.Line(p.pos) != endline || p.tok == token.SEMICOLON || p.tok == token.EOF {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>				<span class="comment">// The next token is on a different line, thus</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>				<span class="comment">// the last comment group is a line comment.</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>				p.lineComment = comment
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		<span class="comment">// consume successor comments, if any</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		endline = -1
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		for p.tok == token.COMMENT {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			comment, endline = p.consumeCommentGroup(1)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		if endline+1 == p.file.Line(p.pos) {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			<span class="comment">// The next token is following on the line immediately after the</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			<span class="comment">// comment group, thus the last comment group is a lead comment.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			p.leadComment = comment
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// A bailout panic is raised to indicate early termination. pos and msg are</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// only populated when bailing out of object resolution.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>type bailout struct {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	pos token.Pos
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	msg string
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>func (p *parser) error(pos token.Pos, msg string) {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	if p.trace {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		defer un(trace(p, &#34;error: &#34;+msg))
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	epos := p.file.Position(pos)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// If AllErrors is not set, discard errors reported on the same line</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	<span class="comment">// as the last recorded error and stop parsing if there are more than</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	<span class="comment">// 10 errors.</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	if p.mode&amp;AllErrors == 0 {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		n := len(p.errors)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		if n &gt; 0 &amp;&amp; p.errors[n-1].Pos.Line == epos.Line {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			return <span class="comment">// discard - likely a spurious error</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		if n &gt; 10 {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			panic(bailout{})
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	p.errors.Add(epos, msg)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>func (p *parser) errorExpected(pos token.Pos, msg string) {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	msg = &#34;expected &#34; + msg
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	if pos == p.pos {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		<span class="comment">// the error happened at the current position;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		<span class="comment">// make the error message more specific</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		switch {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		case p.tok == token.SEMICOLON &amp;&amp; p.lit == &#34;\n&#34;:
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			msg += &#34;, found newline&#34;
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		case p.tok.IsLiteral():
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			<span class="comment">// print 123 rather than &#39;INT&#39;, etc.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			msg += &#34;, found &#34; + p.lit
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		default:
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			msg += &#34;, found &#39;&#34; + p.tok.String() + &#34;&#39;&#34;
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	p.error(pos, msg)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>func (p *parser) expect(tok token.Token) token.Pos {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	pos := p.pos
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	if p.tok != tok {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		p.errorExpected(pos, &#34;&#39;&#34;+tok.String()+&#34;&#39;&#34;)
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	p.next() <span class="comment">// make progress</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	return pos
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// expect2 is like expect, but it returns an invalid position</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">// if the expected token is not found.</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>func (p *parser) expect2(tok token.Token) (pos token.Pos) {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	if p.tok == tok {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		pos = p.pos
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	} else {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		p.errorExpected(p.pos, &#34;&#39;&#34;+tok.String()+&#34;&#39;&#34;)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	p.next() <span class="comment">// make progress</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	return
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span><span class="comment">// expectClosing is like expect but provides a better error message</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span><span class="comment">// for the common case of a missing comma before a newline.</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>func (p *parser) expectClosing(tok token.Token, context string) token.Pos {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	if p.tok != tok &amp;&amp; p.tok == token.SEMICOLON &amp;&amp; p.lit == &#34;\n&#34; {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		p.error(p.pos, &#34;missing &#39;,&#39; before newline in &#34;+context)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		p.next()
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	return p.expect(tok)
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">// expectSemi consumes a semicolon and returns the applicable line comment.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>func (p *parser) expectSemi() (comment *ast.CommentGroup) {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// semicolon is optional before a closing &#39;)&#39; or &#39;}&#39;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	if p.tok != token.RPAREN &amp;&amp; p.tok != token.RBRACE {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		switch p.tok {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		case token.COMMA:
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>			<span class="comment">// permit a &#39;,&#39; instead of a &#39;;&#39; but complain</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			p.errorExpected(p.pos, &#34;&#39;;&#39;&#34;)
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			fallthrough
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		case token.SEMICOLON:
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			if p.lit == &#34;;&#34; {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>				<span class="comment">// explicit semicolon</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>				p.next()
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>				comment = p.lineComment <span class="comment">// use following comments</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			} else {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>				<span class="comment">// artificial semicolon</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>				comment = p.lineComment <span class="comment">// use preceding comments</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>				p.next()
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>			}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			return comment
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		default:
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>			p.errorExpected(p.pos, &#34;&#39;;&#39;&#34;)
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			p.advance(stmtStart)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	return nil
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>func (p *parser) atComma(context string, follow token.Token) bool {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	if p.tok == token.COMMA {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		return true
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	}
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	if p.tok != follow {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		msg := &#34;missing &#39;,&#39;&#34;
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		if p.tok == token.SEMICOLON &amp;&amp; p.lit == &#34;\n&#34; {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			msg += &#34; before newline&#34;
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		}
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		p.error(p.pos, msg+&#34; in &#34;+context)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		return true <span class="comment">// &#34;insert&#34; comma and continue</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	return false
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>func assert(cond bool, msg string) {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	if !cond {
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		panic(&#34;go/parser internal error: &#34; + msg)
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span><span class="comment">// advance consumes tokens until the current token p.tok</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span><span class="comment">// is in the &#39;to&#39; set, or token.EOF. For error recovery.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>func (p *parser) advance(to map[token.Token]bool) {
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	for ; p.tok != token.EOF; p.next() {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		if to[p.tok] {
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>			<span class="comment">// Return only if parser made some progress since last</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			<span class="comment">// sync or if it has not reached 10 advance calls without</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			<span class="comment">// progress. Otherwise consume at least one token to</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			<span class="comment">// avoid an endless parser loop (it is possible that</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			<span class="comment">// both parseOperand and parseStmt call advance and</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			<span class="comment">// correctly do not advance, thus the need for the</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			<span class="comment">// invocation limit p.syncCnt).</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			if p.pos == p.syncPos &amp;&amp; p.syncCnt &lt; 10 {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>				p.syncCnt++
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>				return
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			}
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			if p.pos &gt; p.syncPos {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>				p.syncPos = p.pos
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>				p.syncCnt = 0
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>				return
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			<span class="comment">// Reaching here indicates a parser bug, likely an</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			<span class="comment">// incorrect token list in this function, but it only</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			<span class="comment">// leads to skipping of possibly correct code if a</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			<span class="comment">// previous error is present, and thus is preferred</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			<span class="comment">// over a non-terminating parse.</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>var stmtStart = map[token.Token]bool{
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	token.BREAK:       true,
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	token.CONST:       true,
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	token.CONTINUE:    true,
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	token.DEFER:       true,
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	token.FALLTHROUGH: true,
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	token.FOR:         true,
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	token.GO:          true,
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	token.GOTO:        true,
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	token.IF:          true,
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	token.RETURN:      true,
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	token.SELECT:      true,
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	token.SWITCH:      true,
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	token.TYPE:        true,
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	token.VAR:         true,
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>var declStart = map[token.Token]bool{
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	token.IMPORT: true,
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	token.CONST:  true,
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	token.TYPE:   true,
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	token.VAR:    true,
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>var exprEnd = map[token.Token]bool{
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	token.COMMA:     true,
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	token.COLON:     true,
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	token.SEMICOLON: true,
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	token.RPAREN:    true,
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	token.RBRACK:    true,
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	token.RBRACE:    true,
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span><span class="comment">// safePos returns a valid file position for a given position: If pos</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span><span class="comment">// is valid to begin with, safePos returns pos. If pos is out-of-range,</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span><span class="comment">// safePos returns the EOF position.</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span><span class="comment">// This is hack to work around &#34;artificial&#34; end positions in the AST which</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span><span class="comment">// are computed by adding 1 to (presumably valid) token positions. If the</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span><span class="comment">// token positions are invalid due to parse errors, the resulting end position</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span><span class="comment">// may be past the file&#39;s EOF position, which would lead to panics if used</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span><span class="comment">// later on.</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>func (p *parser) safePos(pos token.Pos) (res token.Pos) {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	defer func() {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		if recover() != nil {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			res = token.Pos(p.file.Base() + p.file.Size()) <span class="comment">// EOF position</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	}()
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	_ = p.file.Offset(pos) <span class="comment">// trigger a panic if position is out-of-range</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	return pos
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span><span class="comment">// Identifiers</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>func (p *parser) parseIdent() *ast.Ident {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	pos := p.pos
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	name := &#34;_&#34;
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	if p.tok == token.IDENT {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		name = p.lit
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		p.next()
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	} else {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		p.expect(token.IDENT) <span class="comment">// use expect() error handling</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	return &amp;ast.Ident{NamePos: pos, Name: name}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>func (p *parser) parseIdentList() (list []*ast.Ident) {
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	if p.trace {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		defer un(trace(p, &#34;IdentList&#34;))
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	list = append(list, p.parseIdent())
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	for p.tok == token.COMMA {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		p.next()
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		list = append(list, p.parseIdent())
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	}
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	return
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span><span class="comment">// Common productions</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span><span class="comment">// If lhs is set, result list elements which are identifiers are not resolved.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>func (p *parser) parseExprList() (list []ast.Expr) {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	if p.trace {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		defer un(trace(p, &#34;ExpressionList&#34;))
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	}
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	list = append(list, p.parseExpr())
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	for p.tok == token.COMMA {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		p.next()
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		list = append(list, p.parseExpr())
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	return
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>func (p *parser) parseList(inRhs bool) []ast.Expr {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	old := p.inRhs
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	p.inRhs = inRhs
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	list := p.parseExprList()
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	p.inRhs = old
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	return list
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span><span class="comment">// Types</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>func (p *parser) parseType() ast.Expr {
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	if p.trace {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Type&#34;))
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	}
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	typ := p.tryIdentOrType()
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	if typ == nil {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		pos := p.pos
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		p.errorExpected(pos, &#34;type&#34;)
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		p.advance(exprEnd)
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		return &amp;ast.BadExpr{From: pos, To: p.pos}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	return typ
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>func (p *parser) parseQualifiedIdent(ident *ast.Ident) ast.Expr {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	if p.trace {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		defer un(trace(p, &#34;QualifiedIdent&#34;))
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	}
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	typ := p.parseTypeName(ident)
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	if p.tok == token.LBRACK {
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		typ = p.parseTypeInstance(typ)
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	return typ
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// If the result is an identifier, it is not resolved.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>func (p *parser) parseTypeName(ident *ast.Ident) ast.Expr {
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	if p.trace {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		defer un(trace(p, &#34;TypeName&#34;))
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	if ident == nil {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		ident = p.parseIdent()
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	if p.tok == token.PERIOD {
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		<span class="comment">// ident is a package name</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		p.next()
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		sel := p.parseIdent()
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		return &amp;ast.SelectorExpr{X: ident, Sel: sel}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	return ident
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span><span class="comment">// &#34;[&#34; has already been consumed, and lbrack is its position.</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span><span class="comment">// If len != nil it is the already consumed array length.</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>func (p *parser) parseArrayType(lbrack token.Pos, len ast.Expr) *ast.ArrayType {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	if p.trace {
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		defer un(trace(p, &#34;ArrayType&#34;))
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	}
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	if len == nil {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		p.exprLev++
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		<span class="comment">// always permit ellipsis for more fault-tolerant parsing</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		if p.tok == token.ELLIPSIS {
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>			len = &amp;ast.Ellipsis{Ellipsis: p.pos}
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>			p.next()
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		} else if p.tok != token.RBRACK {
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>			len = p.parseRhs()
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		p.exprLev--
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	}
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	if p.tok == token.COMMA {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		<span class="comment">// Trailing commas are accepted in type parameter</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		<span class="comment">// lists but not in array type declarations.</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		<span class="comment">// Accept for better error handling but complain.</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		p.error(p.pos, &#34;unexpected comma; expecting ]&#34;)
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		p.next()
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	}
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	p.expect(token.RBRACK)
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	elt := p.parseType()
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	return &amp;ast.ArrayType{Lbrack: lbrack, Len: len, Elt: elt}
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>func (p *parser) parseArrayFieldOrTypeInstance(x *ast.Ident) (*ast.Ident, ast.Expr) {
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	if p.trace {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		defer un(trace(p, &#34;ArrayFieldOrTypeInstance&#34;))
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	lbrack := p.expect(token.LBRACK)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	trailingComma := token.NoPos <span class="comment">// if valid, the position of a trailing comma preceding the &#39;]&#39;</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	var args []ast.Expr
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	if p.tok != token.RBRACK {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		p.exprLev++
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		args = append(args, p.parseRhs())
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		for p.tok == token.COMMA {
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>			comma := p.pos
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>			p.next()
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>			if p.tok == token.RBRACK {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>				trailingComma = comma
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>				break
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>			}
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			args = append(args, p.parseRhs())
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		}
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		p.exprLev--
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	rbrack := p.expect(token.RBRACK)
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	if len(args) == 0 {
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		<span class="comment">// x []E</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		elt := p.parseType()
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		return x, &amp;ast.ArrayType{Lbrack: lbrack, Elt: elt}
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	<span class="comment">// x [P]E or x[P]</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	if len(args) == 1 {
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		elt := p.tryIdentOrType()
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		if elt != nil {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>			<span class="comment">// x [P]E</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>			if trailingComma.IsValid() {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>				<span class="comment">// Trailing commas are invalid in array type fields.</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>				p.error(trailingComma, &#34;unexpected comma; expecting ]&#34;)
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			}
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>			return x, &amp;ast.ArrayType{Lbrack: lbrack, Len: args[0], Elt: elt}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>		}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	<span class="comment">// x[P], x[P1, P2], ...</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	return nil, typeparams.PackIndexExpr(x, lbrack, args, rbrack)
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>func (p *parser) parseFieldDecl() *ast.Field {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	if p.trace {
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		defer un(trace(p, &#34;FieldDecl&#34;))
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	}
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	doc := p.leadComment
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	var names []*ast.Ident
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	var typ ast.Expr
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	switch p.tok {
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	case token.IDENT:
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		name := p.parseIdent()
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		if p.tok == token.PERIOD || p.tok == token.STRING || p.tok == token.SEMICOLON || p.tok == token.RBRACE {
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>			<span class="comment">// embedded type</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>			typ = name
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>			if p.tok == token.PERIOD {
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>				typ = p.parseQualifiedIdent(name)
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>			}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		} else {
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>			<span class="comment">// name1, name2, ... T</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>			names = []*ast.Ident{name}
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>			for p.tok == token.COMMA {
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>				p.next()
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>				names = append(names, p.parseIdent())
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>			}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>			<span class="comment">// Careful dance: We don&#39;t know if we have an embedded instantiated</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>			<span class="comment">// type T[P1, P2, ...] or a field T of array type []E or [P]E.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>			if len(names) == 1 &amp;&amp; p.tok == token.LBRACK {
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>				name, typ = p.parseArrayFieldOrTypeInstance(name)
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>				if name == nil {
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>					names = nil
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>				}
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>			} else {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>				<span class="comment">// T P</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>				typ = p.parseType()
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>			}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		}
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	case token.MUL:
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		star := p.pos
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		p.next()
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		if p.tok == token.LPAREN {
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>			<span class="comment">// *(T)</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>			p.error(p.pos, &#34;cannot parenthesize embedded type&#34;)
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>			p.next()
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>			typ = p.parseQualifiedIdent(nil)
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>			<span class="comment">// expect closing &#39;)&#39; but no need to complain if missing</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>			if p.tok == token.RPAREN {
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>				p.next()
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>			}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		} else {
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			<span class="comment">// *T</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>			typ = p.parseQualifiedIdent(nil)
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		}
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		typ = &amp;ast.StarExpr{Star: star, X: typ}
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	case token.LPAREN:
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		p.error(p.pos, &#34;cannot parenthesize embedded type&#34;)
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		p.next()
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		if p.tok == token.MUL {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			<span class="comment">// (*T)</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>			star := p.pos
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>			p.next()
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			typ = &amp;ast.StarExpr{Star: star, X: p.parseQualifiedIdent(nil)}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>		} else {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>			<span class="comment">// (T)</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>			typ = p.parseQualifiedIdent(nil)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		<span class="comment">// expect closing &#39;)&#39; but no need to complain if missing</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		if p.tok == token.RPAREN {
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>			p.next()
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>		}
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	default:
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		pos := p.pos
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		p.errorExpected(pos, &#34;field name or embedded type&#34;)
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		p.advance(exprEnd)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		typ = &amp;ast.BadExpr{From: pos, To: p.pos}
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	}
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	var tag *ast.BasicLit
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	if p.tok == token.STRING {
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		tag = &amp;ast.BasicLit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>		p.next()
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	}
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	comment := p.expectSemi()
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	field := &amp;ast.Field{Doc: doc, Names: names, Type: typ, Tag: tag, Comment: comment}
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	return field
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>}
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>func (p *parser) parseStructType() *ast.StructType {
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	if p.trace {
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		defer un(trace(p, &#34;StructType&#34;))
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	}
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	pos := p.expect(token.STRUCT)
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	lbrace := p.expect(token.LBRACE)
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	var list []*ast.Field
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	for p.tok == token.IDENT || p.tok == token.MUL || p.tok == token.LPAREN {
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		<span class="comment">// a field declaration cannot start with a &#39;(&#39; but we accept</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		<span class="comment">// it here for more robust parsing and better error messages</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		<span class="comment">// (parseFieldDecl will check and complain if necessary)</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>		list = append(list, p.parseFieldDecl())
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	}
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	rbrace := p.expect(token.RBRACE)
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	return &amp;ast.StructType{
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>		Struct: pos,
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		Fields: &amp;ast.FieldList{
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>			Opening: lbrace,
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>			List:    list,
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>			Closing: rbrace,
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		},
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	}
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>func (p *parser) parsePointerType() *ast.StarExpr {
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	if p.trace {
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>		defer un(trace(p, &#34;PointerType&#34;))
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	}
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	star := p.expect(token.MUL)
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	base := p.parseType()
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	return &amp;ast.StarExpr{Star: star, X: base}
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>func (p *parser) parseDotsType() *ast.Ellipsis {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>	if p.trace {
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		defer un(trace(p, &#34;DotsType&#34;))
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	}
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	pos := p.expect(token.ELLIPSIS)
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	elt := p.parseType()
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	return &amp;ast.Ellipsis{Ellipsis: pos, Elt: elt}
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>}
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>type field struct {
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	name *ast.Ident
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	typ  ast.Expr
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>}
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>func (p *parser) parseParamDecl(name *ast.Ident, typeSetsOK bool) (f field) {
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	<span class="comment">// TODO(rFindley) refactor to be more similar to paramDeclOrNil in the syntax</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	<span class="comment">// package</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	if p.trace {
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>		defer un(trace(p, &#34;ParamDeclOrNil&#34;))
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	}
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	ptok := p.tok
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	if name != nil {
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		p.tok = token.IDENT <span class="comment">// force token.IDENT case in switch below</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	} else if typeSetsOK &amp;&amp; p.tok == token.TILDE {
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>		<span class="comment">// &#34;~&#34; ...</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>		return field{nil, p.embeddedElem(nil)}
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	}
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	switch p.tok {
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	case token.IDENT:
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>		<span class="comment">// name</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		if name != nil {
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>			f.name = name
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>			p.tok = ptok
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		} else {
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>			f.name = p.parseIdent()
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		}
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>		switch p.tok {
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		case token.IDENT, token.MUL, token.ARROW, token.FUNC, token.CHAN, token.MAP, token.STRUCT, token.INTERFACE, token.LPAREN:
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>			<span class="comment">// name type</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>			f.typ = p.parseType()
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>		case token.LBRACK:
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>			<span class="comment">// name &#34;[&#34; type1, ..., typeN &#34;]&#34; or name &#34;[&#34; n &#34;]&#34; type</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>			f.name, f.typ = p.parseArrayFieldOrTypeInstance(f.name)
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		case token.ELLIPSIS:
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>			<span class="comment">// name &#34;...&#34; type</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>			f.typ = p.parseDotsType()
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>			return <span class="comment">// don&#39;t allow ...type &#34;|&#34; ...</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		case token.PERIOD:
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>			<span class="comment">// name &#34;.&#34; ...</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>			f.typ = p.parseQualifiedIdent(f.name)
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>			f.name = nil
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		case token.TILDE:
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>			if typeSetsOK {
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>				f.typ = p.embeddedElem(nil)
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>				return
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>			}
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		case token.OR:
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>			if typeSetsOK {
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>				<span class="comment">// name &#34;|&#34; typeset</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>				f.typ = p.embeddedElem(f.name)
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>				f.name = nil
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>				return
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>			}
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		}
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	case token.MUL, token.ARROW, token.FUNC, token.LBRACK, token.CHAN, token.MAP, token.STRUCT, token.INTERFACE, token.LPAREN:
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		<span class="comment">// type</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>		f.typ = p.parseType()
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	case token.ELLIPSIS:
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		<span class="comment">// &#34;...&#34; type</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>		<span class="comment">// (always accepted)</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		f.typ = p.parseDotsType()
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		return <span class="comment">// don&#39;t allow ...type &#34;|&#34; ...</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	default:
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		<span class="comment">// TODO(rfindley): this is incorrect in the case of type parameter lists</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		<span class="comment">//                 (should be &#34;&#39;]&#39;&#34; in that case)</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		p.errorExpected(p.pos, &#34;&#39;)&#39;&#34;)
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		p.advance(exprEnd)
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	}
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	<span class="comment">// [name] type &#34;|&#34;</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	if typeSetsOK &amp;&amp; p.tok == token.OR &amp;&amp; f.typ != nil {
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>		f.typ = p.embeddedElem(f.typ)
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	}
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	return
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>}
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>func (p *parser) parseParameterList(name0 *ast.Ident, typ0 ast.Expr, closing token.Token) (params []*ast.Field) {
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	if p.trace {
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>		defer un(trace(p, &#34;ParameterList&#34;))
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	}
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	<span class="comment">// Type parameters are the only parameter list closed by &#39;]&#39;.</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	tparams := closing == token.RBRACK
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	pos0 := p.pos
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	if name0 != nil {
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>		pos0 = name0.Pos()
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	} else if typ0 != nil {
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>		pos0 = typ0.Pos()
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>	}
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>	<span class="comment">// Note: The code below matches the corresponding code in the syntax</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	<span class="comment">//       parser closely. Changes must be reflected in either parser.</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	<span class="comment">//       For the code to match, we use the local []field list that</span>
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	<span class="comment">//       corresponds to []syntax.Field. At the end, the list must be</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	<span class="comment">//       converted into an []*ast.Field.</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	var list []field
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	var named int <span class="comment">// number of parameters that have an explicit name and type</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	var typed int <span class="comment">// number of parameters that have an explicit type</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	for name0 != nil || p.tok != closing &amp;&amp; p.tok != token.EOF {
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>		var par field
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>		if typ0 != nil {
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>			if tparams {
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>				typ0 = p.embeddedElem(typ0)
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>			}
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>			par = field{name0, typ0}
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>		} else {
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>			par = p.parseParamDecl(name0, tparams)
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>		}
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>		name0 = nil <span class="comment">// 1st name was consumed if present</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>		typ0 = nil  <span class="comment">// 1st typ was consumed if present</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>		if par.name != nil || par.typ != nil {
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>			list = append(list, par)
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>			if par.name != nil &amp;&amp; par.typ != nil {
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>				named++
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>			}
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>			if par.typ != nil {
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>				typed++
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>			}
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		}
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>		if !p.atComma(&#34;parameter list&#34;, closing) {
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>			break
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>		}
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>		p.next()
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	}
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	if len(list) == 0 {
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		return <span class="comment">// not uncommon</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>	}
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>	<span class="comment">// distribute parameter types (len(list) &gt; 0)</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>	if named == 0 {
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>		<span class="comment">// all unnamed =&gt; found names are type names</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>		for i := 0; i &lt; len(list); i++ {
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>			par := &amp;list[i]
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>			if typ := par.name; typ != nil {
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>				par.typ = typ
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>				par.name = nil
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>			}
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		}
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		if tparams {
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>			<span class="comment">// This is the same error handling as below, adjusted for type parameters only.</span>
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>			<span class="comment">// See comment below for details. (go.dev/issue/64534)</span>
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>			var errPos token.Pos
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>			var msg string
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>			if named == typed <span class="comment">/* same as typed == 0 */</span> {
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>				errPos = p.pos <span class="comment">// position error at closing ]</span>
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>				msg = &#34;missing type constraint&#34;
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>			} else {
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>				errPos = pos0 <span class="comment">// position at opening [ or first name</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>				msg = &#34;missing type parameter name&#34;
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>				if len(list) == 1 {
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>					msg += &#34; or invalid array length&#34;
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>				}
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>			}
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>			p.error(errPos, msg)
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>		}
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>	} else if named != len(list) {
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>		<span class="comment">// some named or we&#39;re in a type parameter list =&gt; all must be named</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>		var errPos token.Pos <span class="comment">// left-most error position (or invalid)</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>		var typ ast.Expr     <span class="comment">// current type (from right to left)</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>		for i := len(list) - 1; i &gt;= 0; i-- {
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>			if par := &amp;list[i]; par.typ != nil {
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>				typ = par.typ
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>				if par.name == nil {
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>					errPos = typ.Pos()
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>					n := ast.NewIdent(&#34;_&#34;)
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>					n.NamePos = errPos <span class="comment">// correct position</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>					par.name = n
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>				}
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>			} else if typ != nil {
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>				par.typ = typ
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>			} else {
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>				<span class="comment">// par.typ == nil &amp;&amp; typ == nil =&gt; we only have a par.name</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>				errPos = par.name.Pos()
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>				par.typ = &amp;ast.BadExpr{From: errPos, To: p.pos}
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>			}
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>		}
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>		if errPos.IsValid() {
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>			var msg string
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>			if tparams {
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>				<span class="comment">// Not all parameters are named because named != len(list).</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>				<span class="comment">// If named == typed we must have parameters that have no types,</span>
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>				<span class="comment">// and they must be at the end of the parameter list, otherwise</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>				<span class="comment">// the types would have been filled in by the right-to-left sweep</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>				<span class="comment">// above and we wouldn&#39;t have an error. Since we are in a type</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>				<span class="comment">// parameter list, the missing types are constraints.</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>				if named == typed {
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>					errPos = p.pos <span class="comment">// position error at closing ]</span>
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>					msg = &#34;missing type constraint&#34;
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>				} else {
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>					msg = &#34;missing type parameter name&#34;
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>					<span class="comment">// go.dev/issue/60812</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>					if len(list) == 1 {
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>						msg += &#34; or invalid array length&#34;
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>					}
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>				}
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>			} else {
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>				msg = &#34;mixed named and unnamed parameters&#34;
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>			}
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>			p.error(errPos, msg)
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>		}
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>	}
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	<span class="comment">// Convert list to []*ast.Field.</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	<span class="comment">// If list contains types only, each type gets its own ast.Field.</span>
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	if named == 0 {
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>		<span class="comment">// parameter list consists of types only</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>		for _, par := range list {
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>			assert(par.typ != nil, &#34;nil type in unnamed parameter list&#34;)
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>			params = append(params, &amp;ast.Field{Type: par.typ})
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>		}
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>		return
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	}
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	<span class="comment">// If the parameter list consists of named parameters with types,</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	<span class="comment">// collect all names with the same types into a single ast.Field.</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>	var names []*ast.Ident
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>	var typ ast.Expr
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	addParams := func() {
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>		assert(typ != nil, &#34;nil type in named parameter list&#34;)
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		field := &amp;ast.Field{Names: names, Type: typ}
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>		params = append(params, field)
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		names = nil
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	}
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>	for _, par := range list {
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>		if par.typ != typ {
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>			if len(names) &gt; 0 {
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>				addParams()
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>			}
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>			typ = par.typ
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		}
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>		names = append(names, par.name)
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	}
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>	if len(names) &gt; 0 {
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>		addParams()
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	}
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>	return
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>}
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>func (p *parser) parseParameters(acceptTParams bool) (tparams, params *ast.FieldList) {
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>	if p.trace {
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Parameters&#34;))
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>	}
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>	if acceptTParams &amp;&amp; p.tok == token.LBRACK {
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>		opening := p.pos
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>		p.next()
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>		<span class="comment">// [T any](params) syntax</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>		list := p.parseParameterList(nil, nil, token.RBRACK)
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>		rbrack := p.expect(token.RBRACK)
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>		tparams = &amp;ast.FieldList{Opening: opening, List: list, Closing: rbrack}
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>		<span class="comment">// Type parameter lists must not be empty.</span>
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>		if tparams.NumFields() == 0 {
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>			p.error(tparams.Closing, &#34;empty type parameter list&#34;)
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>			tparams = nil <span class="comment">// avoid follow-on errors</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>		}
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>	}
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>	opening := p.expect(token.LPAREN)
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>	var fields []*ast.Field
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>	if p.tok != token.RPAREN {
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>		fields = p.parseParameterList(nil, nil, token.RPAREN)
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	}
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>	rparen := p.expect(token.RPAREN)
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>	params = &amp;ast.FieldList{Opening: opening, List: fields, Closing: rparen}
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	return
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>}
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>func (p *parser) parseResult() *ast.FieldList {
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	if p.trace {
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Result&#34;))
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	}
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	if p.tok == token.LPAREN {
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>		_, results := p.parseParameters(false)
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>		return results
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>	}
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>	typ := p.tryIdentOrType()
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>	if typ != nil {
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>		list := make([]*ast.Field, 1)
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>		list[0] = &amp;ast.Field{Type: typ}
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>		return &amp;ast.FieldList{List: list}
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>	}
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>	return nil
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>}
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>func (p *parser) parseFuncType() *ast.FuncType {
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>	if p.trace {
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>		defer un(trace(p, &#34;FuncType&#34;))
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>	}
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	pos := p.expect(token.FUNC)
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>	tparams, params := p.parseParameters(true)
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	if tparams != nil {
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>		p.error(tparams.Pos(), &#34;function type must have no type parameters&#34;)
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>	}
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>	results := p.parseResult()
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>	return &amp;ast.FuncType{Func: pos, Params: params, Results: results}
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>}
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>func (p *parser) parseMethodSpec() *ast.Field {
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>	if p.trace {
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>		defer un(trace(p, &#34;MethodSpec&#34;))
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>	}
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>	doc := p.leadComment
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>	var idents []*ast.Ident
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>	var typ ast.Expr
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>	x := p.parseTypeName(nil)
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>	if ident, _ := x.(*ast.Ident); ident != nil {
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>		switch {
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>		case p.tok == token.LBRACK:
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>			<span class="comment">// generic method or embedded instantiated type</span>
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>			lbrack := p.pos
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>			p.next()
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>			p.exprLev++
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>			x := p.parseExpr()
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>			p.exprLev--
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>			if name0, _ := x.(*ast.Ident); name0 != nil &amp;&amp; p.tok != token.COMMA &amp;&amp; p.tok != token.RBRACK {
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>				<span class="comment">// generic method m[T any]</span>
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>				<span class="comment">// Interface methods do not have type parameters. We parse them for a</span>
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>				<span class="comment">// better error message and improved error recovery.</span>
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>				_ = p.parseParameterList(name0, nil, token.RBRACK)
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>				_ = p.expect(token.RBRACK)
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>				p.error(lbrack, &#34;interface method must have no type parameters&#34;)
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>				<span class="comment">// TODO(rfindley) refactor to share code with parseFuncType.</span>
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>				_, params := p.parseParameters(false)
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>				results := p.parseResult()
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>				idents = []*ast.Ident{ident}
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>				typ = &amp;ast.FuncType{
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>					Func:    token.NoPos,
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>					Params:  params,
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>					Results: results,
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>				}
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>			} else {
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>				<span class="comment">// embedded instantiated type</span>
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>				<span class="comment">// TODO(rfindley) should resolve all identifiers in x.</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>				list := []ast.Expr{x}
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>				if p.atComma(&#34;type argument list&#34;, token.RBRACK) {
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>					p.exprLev++
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>					p.next()
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>					for p.tok != token.RBRACK &amp;&amp; p.tok != token.EOF {
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>						list = append(list, p.parseType())
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>						if !p.atComma(&#34;type argument list&#34;, token.RBRACK) {
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>							break
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>						}
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>						p.next()
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>					}
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>					p.exprLev--
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>				}
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>				rbrack := p.expectClosing(token.RBRACK, &#34;type argument list&#34;)
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>				typ = typeparams.PackIndexExpr(ident, lbrack, list, rbrack)
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>			}
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>		case p.tok == token.LPAREN:
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>			<span class="comment">// ordinary method</span>
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>			<span class="comment">// TODO(rfindley) refactor to share code with parseFuncType.</span>
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>			_, params := p.parseParameters(false)
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>			results := p.parseResult()
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>			idents = []*ast.Ident{ident}
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>			typ = &amp;ast.FuncType{Func: token.NoPos, Params: params, Results: results}
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>		default:
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>			<span class="comment">// embedded type</span>
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>			typ = x
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>		}
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>	} else {
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>		<span class="comment">// embedded, possibly instantiated type</span>
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>		typ = x
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>		if p.tok == token.LBRACK {
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>			<span class="comment">// embedded instantiated interface</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>			typ = p.parseTypeInstance(typ)
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>		}
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>	}
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>	<span class="comment">// Comment is added at the callsite: the field below may joined with</span>
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>	<span class="comment">// additional type specs using &#39;|&#39;.</span>
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>	<span class="comment">// TODO(rfindley) this should be refactored.</span>
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>	<span class="comment">// TODO(rfindley) add more tests for comment handling.</span>
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>	return &amp;ast.Field{Doc: doc, Names: idents, Type: typ}
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>}
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>func (p *parser) embeddedElem(x ast.Expr) ast.Expr {
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>	if p.trace {
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>		defer un(trace(p, &#34;EmbeddedElem&#34;))
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>	}
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>	if x == nil {
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>		x = p.embeddedTerm()
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>	}
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>	for p.tok == token.OR {
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>		t := new(ast.BinaryExpr)
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>		t.OpPos = p.pos
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>		t.Op = token.OR
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>		p.next()
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>		t.X = x
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>		t.Y = p.embeddedTerm()
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>		x = t
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>	}
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>	return x
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>}
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>func (p *parser) embeddedTerm() ast.Expr {
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>	if p.trace {
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>		defer un(trace(p, &#34;EmbeddedTerm&#34;))
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>	}
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>	if p.tok == token.TILDE {
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>		t := new(ast.UnaryExpr)
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>		t.OpPos = p.pos
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>		t.Op = token.TILDE
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>		p.next()
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>		t.X = p.parseType()
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>		return t
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>	}
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>	t := p.tryIdentOrType()
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>	if t == nil {
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>		pos := p.pos
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>		p.errorExpected(pos, &#34;~ term or type&#34;)
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>		p.advance(exprEnd)
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>		return &amp;ast.BadExpr{From: pos, To: p.pos}
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>	}
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	return t
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>}
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>func (p *parser) parseInterfaceType() *ast.InterfaceType {
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>	if p.trace {
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>		defer un(trace(p, &#34;InterfaceType&#34;))
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>	}
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>	pos := p.expect(token.INTERFACE)
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>	lbrace := p.expect(token.LBRACE)
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>	var list []*ast.Field
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>parseElements:
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>	for {
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>		switch {
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>		case p.tok == token.IDENT:
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>			f := p.parseMethodSpec()
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>			if f.Names == nil {
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>				f.Type = p.embeddedElem(f.Type)
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>			}
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>			f.Comment = p.expectSemi()
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>			list = append(list, f)
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>		case p.tok == token.TILDE:
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>			typ := p.embeddedElem(nil)
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>			comment := p.expectSemi()
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>			list = append(list, &amp;ast.Field{Type: typ, Comment: comment})
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>		default:
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>			if t := p.tryIdentOrType(); t != nil {
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>				typ := p.embeddedElem(t)
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>				comment := p.expectSemi()
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>				list = append(list, &amp;ast.Field{Type: typ, Comment: comment})
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>			} else {
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>				break parseElements
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>			}
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>		}
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>	}
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>	<span class="comment">// TODO(rfindley): the error produced here could be improved, since we could</span>
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	<span class="comment">// accept an identifier, &#39;type&#39;, or a &#39;}&#39; at this point.</span>
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	rbrace := p.expect(token.RBRACE)
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	return &amp;ast.InterfaceType{
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>		Interface: pos,
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>		Methods: &amp;ast.FieldList{
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>			Opening: lbrace,
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>			List:    list,
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>			Closing: rbrace,
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>		},
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>	}
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>}
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>func (p *parser) parseMapType() *ast.MapType {
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>	if p.trace {
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>		defer un(trace(p, &#34;MapType&#34;))
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>	}
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>	pos := p.expect(token.MAP)
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>	p.expect(token.LBRACK)
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>	key := p.parseType()
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>	p.expect(token.RBRACK)
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>	value := p.parseType()
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>	return &amp;ast.MapType{Map: pos, Key: key, Value: value}
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>}
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>func (p *parser) parseChanType() *ast.ChanType {
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>	if p.trace {
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>		defer un(trace(p, &#34;ChanType&#34;))
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>	}
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>	pos := p.pos
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>	dir := ast.SEND | ast.RECV
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>	var arrow token.Pos
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>	if p.tok == token.CHAN {
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>		p.next()
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>		if p.tok == token.ARROW {
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>			arrow = p.pos
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>			p.next()
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>			dir = ast.SEND
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>		}
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>	} else {
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>		arrow = p.expect(token.ARROW)
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>		p.expect(token.CHAN)
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>		dir = ast.RECV
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>	}
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>	value := p.parseType()
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>	return &amp;ast.ChanType{Begin: pos, Arrow: arrow, Dir: dir, Value: value}
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>}
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>func (p *parser) parseTypeInstance(typ ast.Expr) ast.Expr {
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>	if p.trace {
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>		defer un(trace(p, &#34;TypeInstance&#34;))
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>	}
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>	opening := p.expect(token.LBRACK)
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>	p.exprLev++
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>	var list []ast.Expr
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>	for p.tok != token.RBRACK &amp;&amp; p.tok != token.EOF {
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>		list = append(list, p.parseType())
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>		if !p.atComma(&#34;type argument list&#34;, token.RBRACK) {
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>			break
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>		}
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>		p.next()
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>	}
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>	p.exprLev--
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>	closing := p.expectClosing(token.RBRACK, &#34;type argument list&#34;)
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>	if len(list) == 0 {
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>		p.errorExpected(closing, &#34;type argument list&#34;)
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>		return &amp;ast.IndexExpr{
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>			X:      typ,
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>			Lbrack: opening,
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>			Index:  &amp;ast.BadExpr{From: opening + 1, To: closing},
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>			Rbrack: closing,
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>		}
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>	}
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>	return typeparams.PackIndexExpr(typ, opening, list, closing)
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>}
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>func (p *parser) tryIdentOrType() ast.Expr {
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>	defer decNestLev(incNestLev(p))
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>	switch p.tok {
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>	case token.IDENT:
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>		typ := p.parseTypeName(nil)
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>		if p.tok == token.LBRACK {
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>			typ = p.parseTypeInstance(typ)
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>		}
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>		return typ
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>	case token.LBRACK:
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>		lbrack := p.expect(token.LBRACK)
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>		return p.parseArrayType(lbrack, nil)
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>	case token.STRUCT:
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>		return p.parseStructType()
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>	case token.MUL:
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>		return p.parsePointerType()
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>	case token.FUNC:
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>		return p.parseFuncType()
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>	case token.INTERFACE:
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>		return p.parseInterfaceType()
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>	case token.MAP:
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>		return p.parseMapType()
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>	case token.CHAN, token.ARROW:
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>		return p.parseChanType()
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>	case token.LPAREN:
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>		lparen := p.pos
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>		p.next()
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>		typ := p.parseType()
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>		rparen := p.expect(token.RPAREN)
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>		return &amp;ast.ParenExpr{Lparen: lparen, X: typ, Rparen: rparen}
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>	}
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>	<span class="comment">// no type found</span>
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>	return nil
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>}
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span><span class="comment">// Blocks</span>
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>func (p *parser) parseStmtList() (list []ast.Stmt) {
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>	if p.trace {
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>		defer un(trace(p, &#34;StatementList&#34;))
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>	}
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>	for p.tok != token.CASE &amp;&amp; p.tok != token.DEFAULT &amp;&amp; p.tok != token.RBRACE &amp;&amp; p.tok != token.EOF {
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>		list = append(list, p.parseStmt())
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>	}
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>	return
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>}
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>func (p *parser) parseBody() *ast.BlockStmt {
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>	if p.trace {
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Body&#34;))
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>	}
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>	lbrace := p.expect(token.LBRACE)
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>	list := p.parseStmtList()
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>	rbrace := p.expect2(token.RBRACE)
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>	return &amp;ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>}
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>func (p *parser) parseBlockStmt() *ast.BlockStmt {
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>	if p.trace {
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>		defer un(trace(p, &#34;BlockStmt&#34;))
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>	}
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>	lbrace := p.expect(token.LBRACE)
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>	list := p.parseStmtList()
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>	rbrace := p.expect2(token.RBRACE)
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>	return &amp;ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>}
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span><span class="comment">// Expressions</span>
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>func (p *parser) parseFuncTypeOrLit() ast.Expr {
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>	if p.trace {
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span>		defer un(trace(p, &#34;FuncTypeOrLit&#34;))
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span>	}
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span>
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span>	typ := p.parseFuncType()
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span>	if p.tok != token.LBRACE {
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span>		<span class="comment">// function type only</span>
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>		return typ
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>	}
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>	p.exprLev++
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>	body := p.parseBody()
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>	p.exprLev--
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>	return &amp;ast.FuncLit{Type: typ, Body: body}
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>}
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span>
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span><span class="comment">// parseOperand may return an expression or a raw type (incl. array</span>
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span><span class="comment">// types of the form [...]T). Callers must verify the result.</span>
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span>func (p *parser) parseOperand() ast.Expr {
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span>	if p.trace {
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Operand&#34;))
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span>	}
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span>
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span>	switch p.tok {
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span>	case token.IDENT:
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span>		x := p.parseIdent()
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span>		return x
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span>
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span>	case token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING:
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span>		x := &amp;ast.BasicLit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span>		p.next()
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span>		return x
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span>
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span>	case token.LPAREN:
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span>		lparen := p.pos
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span>		p.next()
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span>		p.exprLev++
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span>		x := p.parseRhs() <span class="comment">// types may be parenthesized: (some type)</span>
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span>		p.exprLev--
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span>		rparen := p.expect(token.RPAREN)
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span>		return &amp;ast.ParenExpr{Lparen: lparen, X: x, Rparen: rparen}
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span>
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span>	case token.FUNC:
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span>		return p.parseFuncTypeOrLit()
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span>	}
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span>
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span>	if typ := p.tryIdentOrType(); typ != nil { <span class="comment">// do not consume trailing type parameters</span>
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span>		<span class="comment">// could be type for composite literal or conversion</span>
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span>		_, isIdent := typ.(*ast.Ident)
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span>		assert(!isIdent, &#34;type cannot be identifier&#34;)
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span>		return typ
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span>	}
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span>
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span>	<span class="comment">// we have an error</span>
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span>	pos := p.pos
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span>	p.errorExpected(pos, &#34;operand&#34;)
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span>	p.advance(stmtStart)
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span>	return &amp;ast.BadExpr{From: pos, To: p.pos}
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span>}
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span>
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span>func (p *parser) parseSelector(x ast.Expr) ast.Expr {
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span>	if p.trace {
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Selector&#34;))
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span>	}
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span>
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span>	sel := p.parseIdent()
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span>
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span>	return &amp;ast.SelectorExpr{X: x, Sel: sel}
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span>}
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span>
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span>func (p *parser) parseTypeAssertion(x ast.Expr) ast.Expr {
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span>	if p.trace {
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span>		defer un(trace(p, &#34;TypeAssertion&#34;))
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span>	}
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span>
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span>	lparen := p.expect(token.LPAREN)
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span>	var typ ast.Expr
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span>	if p.tok == token.TYPE {
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span>		<span class="comment">// type switch: typ == nil</span>
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span>		p.next()
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>	} else {
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>		typ = p.parseType()
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span>	}
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span>	rparen := p.expect(token.RPAREN)
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span>
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span>	return &amp;ast.TypeAssertExpr{X: x, Type: typ, Lparen: lparen, Rparen: rparen}
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span>}
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span>
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span>func (p *parser) parseIndexOrSliceOrInstance(x ast.Expr) ast.Expr {
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span>	if p.trace {
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span>		defer un(trace(p, &#34;parseIndexOrSliceOrInstance&#34;))
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span>	}
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span>
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span>	lbrack := p.expect(token.LBRACK)
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span>	if p.tok == token.RBRACK {
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>		<span class="comment">// empty index, slice or index expressions are not permitted;</span>
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>		<span class="comment">// accept them for parsing tolerance, but complain</span>
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>		p.errorExpected(p.pos, &#34;operand&#34;)
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>		rbrack := p.pos
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span>		p.next()
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span>		return &amp;ast.IndexExpr{
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span>			X:      x,
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span>			Lbrack: lbrack,
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>			Index:  &amp;ast.BadExpr{From: rbrack, To: rbrack},
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>			Rbrack: rbrack,
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>		}
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span>	}
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span>	p.exprLev++
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span>
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span>	const N = 3 <span class="comment">// change the 3 to 2 to disable 3-index slices</span>
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>	var args []ast.Expr
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>	var index [N]ast.Expr
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>	var colons [N - 1]token.Pos
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>	if p.tok != token.COLON {
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span>		<span class="comment">// We can&#39;t know if we have an index expression or a type instantiation;</span>
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span>		<span class="comment">// so even if we see a (named) type we are not going to be in type context.</span>
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span>		index[0] = p.parseRhs()
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span>	}
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span>	ncolons := 0
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span>	switch p.tok {
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span>	case token.COLON:
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span>		<span class="comment">// slice expression</span>
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span>		for p.tok == token.COLON &amp;&amp; ncolons &lt; len(colons) {
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span>			colons[ncolons] = p.pos
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span>			ncolons++
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span>			p.next()
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span>			if p.tok != token.COLON &amp;&amp; p.tok != token.RBRACK &amp;&amp; p.tok != token.EOF {
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span>				index[ncolons] = p.parseRhs()
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span>			}
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span>		}
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span>	case token.COMMA:
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span>		<span class="comment">// instance expression</span>
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span>		args = append(args, index[0])
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span>		for p.tok == token.COMMA {
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span>			p.next()
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span>			if p.tok != token.RBRACK &amp;&amp; p.tok != token.EOF {
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span>				args = append(args, p.parseType())
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span>			}
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span>		}
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span>	}
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span>
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span>	p.exprLev--
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span>	rbrack := p.expect(token.RBRACK)
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span>
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span>	if ncolons &gt; 0 {
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span>		<span class="comment">// slice expression</span>
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span>		slice3 := false
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span>		if ncolons == 2 {
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>			slice3 = true
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span>			<span class="comment">// Check presence of middle and final index here rather than during type-checking</span>
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span>			<span class="comment">// to prevent erroneous programs from passing through gofmt (was go.dev/issue/7305).</span>
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span>			if index[1] == nil {
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span>				p.error(colons[0], &#34;middle index required in 3-index slice&#34;)
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span>				index[1] = &amp;ast.BadExpr{From: colons[0] + 1, To: colons[1]}
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>			}
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>			if index[2] == nil {
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>				p.error(colons[1], &#34;final index required in 3-index slice&#34;)
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>				index[2] = &amp;ast.BadExpr{From: colons[1] + 1, To: rbrack}
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>			}
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>		}
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>		return &amp;ast.SliceExpr{X: x, Lbrack: lbrack, Low: index[0], High: index[1], Max: index[2], Slice3: slice3, Rbrack: rbrack}
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span>	}
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span>
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span>	if len(args) == 0 {
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span>		<span class="comment">// index expression</span>
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span>		return &amp;ast.IndexExpr{X: x, Lbrack: lbrack, Index: index[0], Rbrack: rbrack}
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>	}
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span>	<span class="comment">// instance expression</span>
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span>	return typeparams.PackIndexExpr(x, lbrack, args, rbrack)
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span>}
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span>
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span>func (p *parser) parseCallOrConversion(fun ast.Expr) *ast.CallExpr {
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span>	if p.trace {
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span>		defer un(trace(p, &#34;CallOrConversion&#34;))
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span>	}
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span>
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span>	lparen := p.expect(token.LPAREN)
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span>	p.exprLev++
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span>	var list []ast.Expr
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>	var ellipsis token.Pos
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>	for p.tok != token.RPAREN &amp;&amp; p.tok != token.EOF &amp;&amp; !ellipsis.IsValid() {
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>		list = append(list, p.parseRhs()) <span class="comment">// builtins may expect a type: make(some type, ...)</span>
<span id="L1618" class="ln">  1618&nbsp;&nbsp;</span>		if p.tok == token.ELLIPSIS {
<span id="L1619" class="ln">  1619&nbsp;&nbsp;</span>			ellipsis = p.pos
<span id="L1620" class="ln">  1620&nbsp;&nbsp;</span>			p.next()
<span id="L1621" class="ln">  1621&nbsp;&nbsp;</span>		}
<span id="L1622" class="ln">  1622&nbsp;&nbsp;</span>		if !p.atComma(&#34;argument list&#34;, token.RPAREN) {
<span id="L1623" class="ln">  1623&nbsp;&nbsp;</span>			break
<span id="L1624" class="ln">  1624&nbsp;&nbsp;</span>		}
<span id="L1625" class="ln">  1625&nbsp;&nbsp;</span>		p.next()
<span id="L1626" class="ln">  1626&nbsp;&nbsp;</span>	}
<span id="L1627" class="ln">  1627&nbsp;&nbsp;</span>	p.exprLev--
<span id="L1628" class="ln">  1628&nbsp;&nbsp;</span>	rparen := p.expectClosing(token.RPAREN, &#34;argument list&#34;)
<span id="L1629" class="ln">  1629&nbsp;&nbsp;</span>
<span id="L1630" class="ln">  1630&nbsp;&nbsp;</span>	return &amp;ast.CallExpr{Fun: fun, Lparen: lparen, Args: list, Ellipsis: ellipsis, Rparen: rparen}
<span id="L1631" class="ln">  1631&nbsp;&nbsp;</span>}
<span id="L1632" class="ln">  1632&nbsp;&nbsp;</span>
<span id="L1633" class="ln">  1633&nbsp;&nbsp;</span>func (p *parser) parseValue() ast.Expr {
<span id="L1634" class="ln">  1634&nbsp;&nbsp;</span>	if p.trace {
<span id="L1635" class="ln">  1635&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Element&#34;))
<span id="L1636" class="ln">  1636&nbsp;&nbsp;</span>	}
<span id="L1637" class="ln">  1637&nbsp;&nbsp;</span>
<span id="L1638" class="ln">  1638&nbsp;&nbsp;</span>	if p.tok == token.LBRACE {
<span id="L1639" class="ln">  1639&nbsp;&nbsp;</span>		return p.parseLiteralValue(nil)
<span id="L1640" class="ln">  1640&nbsp;&nbsp;</span>	}
<span id="L1641" class="ln">  1641&nbsp;&nbsp;</span>
<span id="L1642" class="ln">  1642&nbsp;&nbsp;</span>	x := p.parseExpr()
<span id="L1643" class="ln">  1643&nbsp;&nbsp;</span>
<span id="L1644" class="ln">  1644&nbsp;&nbsp;</span>	return x
<span id="L1645" class="ln">  1645&nbsp;&nbsp;</span>}
<span id="L1646" class="ln">  1646&nbsp;&nbsp;</span>
<span id="L1647" class="ln">  1647&nbsp;&nbsp;</span>func (p *parser) parseElement() ast.Expr {
<span id="L1648" class="ln">  1648&nbsp;&nbsp;</span>	if p.trace {
<span id="L1649" class="ln">  1649&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Element&#34;))
<span id="L1650" class="ln">  1650&nbsp;&nbsp;</span>	}
<span id="L1651" class="ln">  1651&nbsp;&nbsp;</span>
<span id="L1652" class="ln">  1652&nbsp;&nbsp;</span>	x := p.parseValue()
<span id="L1653" class="ln">  1653&nbsp;&nbsp;</span>	if p.tok == token.COLON {
<span id="L1654" class="ln">  1654&nbsp;&nbsp;</span>		colon := p.pos
<span id="L1655" class="ln">  1655&nbsp;&nbsp;</span>		p.next()
<span id="L1656" class="ln">  1656&nbsp;&nbsp;</span>		x = &amp;ast.KeyValueExpr{Key: x, Colon: colon, Value: p.parseValue()}
<span id="L1657" class="ln">  1657&nbsp;&nbsp;</span>	}
<span id="L1658" class="ln">  1658&nbsp;&nbsp;</span>
<span id="L1659" class="ln">  1659&nbsp;&nbsp;</span>	return x
<span id="L1660" class="ln">  1660&nbsp;&nbsp;</span>}
<span id="L1661" class="ln">  1661&nbsp;&nbsp;</span>
<span id="L1662" class="ln">  1662&nbsp;&nbsp;</span>func (p *parser) parseElementList() (list []ast.Expr) {
<span id="L1663" class="ln">  1663&nbsp;&nbsp;</span>	if p.trace {
<span id="L1664" class="ln">  1664&nbsp;&nbsp;</span>		defer un(trace(p, &#34;ElementList&#34;))
<span id="L1665" class="ln">  1665&nbsp;&nbsp;</span>	}
<span id="L1666" class="ln">  1666&nbsp;&nbsp;</span>
<span id="L1667" class="ln">  1667&nbsp;&nbsp;</span>	for p.tok != token.RBRACE &amp;&amp; p.tok != token.EOF {
<span id="L1668" class="ln">  1668&nbsp;&nbsp;</span>		list = append(list, p.parseElement())
<span id="L1669" class="ln">  1669&nbsp;&nbsp;</span>		if !p.atComma(&#34;composite literal&#34;, token.RBRACE) {
<span id="L1670" class="ln">  1670&nbsp;&nbsp;</span>			break
<span id="L1671" class="ln">  1671&nbsp;&nbsp;</span>		}
<span id="L1672" class="ln">  1672&nbsp;&nbsp;</span>		p.next()
<span id="L1673" class="ln">  1673&nbsp;&nbsp;</span>	}
<span id="L1674" class="ln">  1674&nbsp;&nbsp;</span>
<span id="L1675" class="ln">  1675&nbsp;&nbsp;</span>	return
<span id="L1676" class="ln">  1676&nbsp;&nbsp;</span>}
<span id="L1677" class="ln">  1677&nbsp;&nbsp;</span>
<span id="L1678" class="ln">  1678&nbsp;&nbsp;</span>func (p *parser) parseLiteralValue(typ ast.Expr) ast.Expr {
<span id="L1679" class="ln">  1679&nbsp;&nbsp;</span>	defer decNestLev(incNestLev(p))
<span id="L1680" class="ln">  1680&nbsp;&nbsp;</span>
<span id="L1681" class="ln">  1681&nbsp;&nbsp;</span>	if p.trace {
<span id="L1682" class="ln">  1682&nbsp;&nbsp;</span>		defer un(trace(p, &#34;LiteralValue&#34;))
<span id="L1683" class="ln">  1683&nbsp;&nbsp;</span>	}
<span id="L1684" class="ln">  1684&nbsp;&nbsp;</span>
<span id="L1685" class="ln">  1685&nbsp;&nbsp;</span>	lbrace := p.expect(token.LBRACE)
<span id="L1686" class="ln">  1686&nbsp;&nbsp;</span>	var elts []ast.Expr
<span id="L1687" class="ln">  1687&nbsp;&nbsp;</span>	p.exprLev++
<span id="L1688" class="ln">  1688&nbsp;&nbsp;</span>	if p.tok != token.RBRACE {
<span id="L1689" class="ln">  1689&nbsp;&nbsp;</span>		elts = p.parseElementList()
<span id="L1690" class="ln">  1690&nbsp;&nbsp;</span>	}
<span id="L1691" class="ln">  1691&nbsp;&nbsp;</span>	p.exprLev--
<span id="L1692" class="ln">  1692&nbsp;&nbsp;</span>	rbrace := p.expectClosing(token.RBRACE, &#34;composite literal&#34;)
<span id="L1693" class="ln">  1693&nbsp;&nbsp;</span>	return &amp;ast.CompositeLit{Type: typ, Lbrace: lbrace, Elts: elts, Rbrace: rbrace}
<span id="L1694" class="ln">  1694&nbsp;&nbsp;</span>}
<span id="L1695" class="ln">  1695&nbsp;&nbsp;</span>
<span id="L1696" class="ln">  1696&nbsp;&nbsp;</span>func (p *parser) parsePrimaryExpr(x ast.Expr) ast.Expr {
<span id="L1697" class="ln">  1697&nbsp;&nbsp;</span>	if p.trace {
<span id="L1698" class="ln">  1698&nbsp;&nbsp;</span>		defer un(trace(p, &#34;PrimaryExpr&#34;))
<span id="L1699" class="ln">  1699&nbsp;&nbsp;</span>	}
<span id="L1700" class="ln">  1700&nbsp;&nbsp;</span>
<span id="L1701" class="ln">  1701&nbsp;&nbsp;</span>	if x == nil {
<span id="L1702" class="ln">  1702&nbsp;&nbsp;</span>		x = p.parseOperand()
<span id="L1703" class="ln">  1703&nbsp;&nbsp;</span>	}
<span id="L1704" class="ln">  1704&nbsp;&nbsp;</span>	<span class="comment">// We track the nesting here rather than at the entry for the function,</span>
<span id="L1705" class="ln">  1705&nbsp;&nbsp;</span>	<span class="comment">// since it can iteratively produce a nested output, and we want to</span>
<span id="L1706" class="ln">  1706&nbsp;&nbsp;</span>	<span class="comment">// limit how deep a structure we generate.</span>
<span id="L1707" class="ln">  1707&nbsp;&nbsp;</span>	var n int
<span id="L1708" class="ln">  1708&nbsp;&nbsp;</span>	defer func() { p.nestLev -= n }()
<span id="L1709" class="ln">  1709&nbsp;&nbsp;</span>	for n = 1; ; n++ {
<span id="L1710" class="ln">  1710&nbsp;&nbsp;</span>		incNestLev(p)
<span id="L1711" class="ln">  1711&nbsp;&nbsp;</span>		switch p.tok {
<span id="L1712" class="ln">  1712&nbsp;&nbsp;</span>		case token.PERIOD:
<span id="L1713" class="ln">  1713&nbsp;&nbsp;</span>			p.next()
<span id="L1714" class="ln">  1714&nbsp;&nbsp;</span>			switch p.tok {
<span id="L1715" class="ln">  1715&nbsp;&nbsp;</span>			case token.IDENT:
<span id="L1716" class="ln">  1716&nbsp;&nbsp;</span>				x = p.parseSelector(x)
<span id="L1717" class="ln">  1717&nbsp;&nbsp;</span>			case token.LPAREN:
<span id="L1718" class="ln">  1718&nbsp;&nbsp;</span>				x = p.parseTypeAssertion(x)
<span id="L1719" class="ln">  1719&nbsp;&nbsp;</span>			default:
<span id="L1720" class="ln">  1720&nbsp;&nbsp;</span>				pos := p.pos
<span id="L1721" class="ln">  1721&nbsp;&nbsp;</span>				p.errorExpected(pos, &#34;selector or type assertion&#34;)
<span id="L1722" class="ln">  1722&nbsp;&nbsp;</span>				<span class="comment">// TODO(rFindley) The check for token.RBRACE below is a targeted fix</span>
<span id="L1723" class="ln">  1723&nbsp;&nbsp;</span>				<span class="comment">//                to error recovery sufficient to make the x/tools tests to</span>
<span id="L1724" class="ln">  1724&nbsp;&nbsp;</span>				<span class="comment">//                pass with the new parsing logic introduced for type</span>
<span id="L1725" class="ln">  1725&nbsp;&nbsp;</span>				<span class="comment">//                parameters. Remove this once error recovery has been</span>
<span id="L1726" class="ln">  1726&nbsp;&nbsp;</span>				<span class="comment">//                more generally reconsidered.</span>
<span id="L1727" class="ln">  1727&nbsp;&nbsp;</span>				if p.tok != token.RBRACE {
<span id="L1728" class="ln">  1728&nbsp;&nbsp;</span>					p.next() <span class="comment">// make progress</span>
<span id="L1729" class="ln">  1729&nbsp;&nbsp;</span>				}
<span id="L1730" class="ln">  1730&nbsp;&nbsp;</span>				sel := &amp;ast.Ident{NamePos: pos, Name: &#34;_&#34;}
<span id="L1731" class="ln">  1731&nbsp;&nbsp;</span>				x = &amp;ast.SelectorExpr{X: x, Sel: sel}
<span id="L1732" class="ln">  1732&nbsp;&nbsp;</span>			}
<span id="L1733" class="ln">  1733&nbsp;&nbsp;</span>		case token.LBRACK:
<span id="L1734" class="ln">  1734&nbsp;&nbsp;</span>			x = p.parseIndexOrSliceOrInstance(x)
<span id="L1735" class="ln">  1735&nbsp;&nbsp;</span>		case token.LPAREN:
<span id="L1736" class="ln">  1736&nbsp;&nbsp;</span>			x = p.parseCallOrConversion(x)
<span id="L1737" class="ln">  1737&nbsp;&nbsp;</span>		case token.LBRACE:
<span id="L1738" class="ln">  1738&nbsp;&nbsp;</span>			<span class="comment">// operand may have returned a parenthesized complit</span>
<span id="L1739" class="ln">  1739&nbsp;&nbsp;</span>			<span class="comment">// type; accept it but complain if we have a complit</span>
<span id="L1740" class="ln">  1740&nbsp;&nbsp;</span>			t := ast.Unparen(x)
<span id="L1741" class="ln">  1741&nbsp;&nbsp;</span>			<span class="comment">// determine if &#39;{&#39; belongs to a composite literal or a block statement</span>
<span id="L1742" class="ln">  1742&nbsp;&nbsp;</span>			switch t.(type) {
<span id="L1743" class="ln">  1743&nbsp;&nbsp;</span>			case *ast.BadExpr, *ast.Ident, *ast.SelectorExpr:
<span id="L1744" class="ln">  1744&nbsp;&nbsp;</span>				if p.exprLev &lt; 0 {
<span id="L1745" class="ln">  1745&nbsp;&nbsp;</span>					return x
<span id="L1746" class="ln">  1746&nbsp;&nbsp;</span>				}
<span id="L1747" class="ln">  1747&nbsp;&nbsp;</span>				<span class="comment">// x is possibly a composite literal type</span>
<span id="L1748" class="ln">  1748&nbsp;&nbsp;</span>			case *ast.IndexExpr, *ast.IndexListExpr:
<span id="L1749" class="ln">  1749&nbsp;&nbsp;</span>				if p.exprLev &lt; 0 {
<span id="L1750" class="ln">  1750&nbsp;&nbsp;</span>					return x
<span id="L1751" class="ln">  1751&nbsp;&nbsp;</span>				}
<span id="L1752" class="ln">  1752&nbsp;&nbsp;</span>				<span class="comment">// x is possibly a composite literal type</span>
<span id="L1753" class="ln">  1753&nbsp;&nbsp;</span>			case *ast.ArrayType, *ast.StructType, *ast.MapType:
<span id="L1754" class="ln">  1754&nbsp;&nbsp;</span>				<span class="comment">// x is a composite literal type</span>
<span id="L1755" class="ln">  1755&nbsp;&nbsp;</span>			default:
<span id="L1756" class="ln">  1756&nbsp;&nbsp;</span>				return x
<span id="L1757" class="ln">  1757&nbsp;&nbsp;</span>			}
<span id="L1758" class="ln">  1758&nbsp;&nbsp;</span>			if t != x {
<span id="L1759" class="ln">  1759&nbsp;&nbsp;</span>				p.error(t.Pos(), &#34;cannot parenthesize type in composite literal&#34;)
<span id="L1760" class="ln">  1760&nbsp;&nbsp;</span>				<span class="comment">// already progressed, no need to advance</span>
<span id="L1761" class="ln">  1761&nbsp;&nbsp;</span>			}
<span id="L1762" class="ln">  1762&nbsp;&nbsp;</span>			x = p.parseLiteralValue(x)
<span id="L1763" class="ln">  1763&nbsp;&nbsp;</span>		default:
<span id="L1764" class="ln">  1764&nbsp;&nbsp;</span>			return x
<span id="L1765" class="ln">  1765&nbsp;&nbsp;</span>		}
<span id="L1766" class="ln">  1766&nbsp;&nbsp;</span>	}
<span id="L1767" class="ln">  1767&nbsp;&nbsp;</span>}
<span id="L1768" class="ln">  1768&nbsp;&nbsp;</span>
<span id="L1769" class="ln">  1769&nbsp;&nbsp;</span>func (p *parser) parseUnaryExpr() ast.Expr {
<span id="L1770" class="ln">  1770&nbsp;&nbsp;</span>	defer decNestLev(incNestLev(p))
<span id="L1771" class="ln">  1771&nbsp;&nbsp;</span>
<span id="L1772" class="ln">  1772&nbsp;&nbsp;</span>	if p.trace {
<span id="L1773" class="ln">  1773&nbsp;&nbsp;</span>		defer un(trace(p, &#34;UnaryExpr&#34;))
<span id="L1774" class="ln">  1774&nbsp;&nbsp;</span>	}
<span id="L1775" class="ln">  1775&nbsp;&nbsp;</span>
<span id="L1776" class="ln">  1776&nbsp;&nbsp;</span>	switch p.tok {
<span id="L1777" class="ln">  1777&nbsp;&nbsp;</span>	case token.ADD, token.SUB, token.NOT, token.XOR, token.AND, token.TILDE:
<span id="L1778" class="ln">  1778&nbsp;&nbsp;</span>		pos, op := p.pos, p.tok
<span id="L1779" class="ln">  1779&nbsp;&nbsp;</span>		p.next()
<span id="L1780" class="ln">  1780&nbsp;&nbsp;</span>		x := p.parseUnaryExpr()
<span id="L1781" class="ln">  1781&nbsp;&nbsp;</span>		return &amp;ast.UnaryExpr{OpPos: pos, Op: op, X: x}
<span id="L1782" class="ln">  1782&nbsp;&nbsp;</span>
<span id="L1783" class="ln">  1783&nbsp;&nbsp;</span>	case token.ARROW:
<span id="L1784" class="ln">  1784&nbsp;&nbsp;</span>		<span class="comment">// channel type or receive expression</span>
<span id="L1785" class="ln">  1785&nbsp;&nbsp;</span>		arrow := p.pos
<span id="L1786" class="ln">  1786&nbsp;&nbsp;</span>		p.next()
<span id="L1787" class="ln">  1787&nbsp;&nbsp;</span>
<span id="L1788" class="ln">  1788&nbsp;&nbsp;</span>		<span class="comment">// If the next token is token.CHAN we still don&#39;t know if it</span>
<span id="L1789" class="ln">  1789&nbsp;&nbsp;</span>		<span class="comment">// is a channel type or a receive operation - we only know</span>
<span id="L1790" class="ln">  1790&nbsp;&nbsp;</span>		<span class="comment">// once we have found the end of the unary expression. There</span>
<span id="L1791" class="ln">  1791&nbsp;&nbsp;</span>		<span class="comment">// are two cases:</span>
<span id="L1792" class="ln">  1792&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1793" class="ln">  1793&nbsp;&nbsp;</span>		<span class="comment">//   &lt;- type  =&gt; (&lt;-type) must be channel type</span>
<span id="L1794" class="ln">  1794&nbsp;&nbsp;</span>		<span class="comment">//   &lt;- expr  =&gt; &lt;-(expr) is a receive from an expression</span>
<span id="L1795" class="ln">  1795&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1796" class="ln">  1796&nbsp;&nbsp;</span>		<span class="comment">// In the first case, the arrow must be re-associated with</span>
<span id="L1797" class="ln">  1797&nbsp;&nbsp;</span>		<span class="comment">// the channel type parsed already:</span>
<span id="L1798" class="ln">  1798&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1799" class="ln">  1799&nbsp;&nbsp;</span>		<span class="comment">//   &lt;- (chan type)    =&gt;  (&lt;-chan type)</span>
<span id="L1800" class="ln">  1800&nbsp;&nbsp;</span>		<span class="comment">//   &lt;- (chan&lt;- type)  =&gt;  (&lt;-chan (&lt;-type))</span>
<span id="L1801" class="ln">  1801&nbsp;&nbsp;</span>
<span id="L1802" class="ln">  1802&nbsp;&nbsp;</span>		x := p.parseUnaryExpr()
<span id="L1803" class="ln">  1803&nbsp;&nbsp;</span>
<span id="L1804" class="ln">  1804&nbsp;&nbsp;</span>		<span class="comment">// determine which case we have</span>
<span id="L1805" class="ln">  1805&nbsp;&nbsp;</span>		if typ, ok := x.(*ast.ChanType); ok {
<span id="L1806" class="ln">  1806&nbsp;&nbsp;</span>			<span class="comment">// (&lt;-type)</span>
<span id="L1807" class="ln">  1807&nbsp;&nbsp;</span>
<span id="L1808" class="ln">  1808&nbsp;&nbsp;</span>			<span class="comment">// re-associate position info and &lt;-</span>
<span id="L1809" class="ln">  1809&nbsp;&nbsp;</span>			dir := ast.SEND
<span id="L1810" class="ln">  1810&nbsp;&nbsp;</span>			for ok &amp;&amp; dir == ast.SEND {
<span id="L1811" class="ln">  1811&nbsp;&nbsp;</span>				if typ.Dir == ast.RECV {
<span id="L1812" class="ln">  1812&nbsp;&nbsp;</span>					<span class="comment">// error: (&lt;-type) is (&lt;-(&lt;-chan T))</span>
<span id="L1813" class="ln">  1813&nbsp;&nbsp;</span>					p.errorExpected(typ.Arrow, &#34;&#39;chan&#39;&#34;)
<span id="L1814" class="ln">  1814&nbsp;&nbsp;</span>				}
<span id="L1815" class="ln">  1815&nbsp;&nbsp;</span>				arrow, typ.Begin, typ.Arrow = typ.Arrow, arrow, arrow
<span id="L1816" class="ln">  1816&nbsp;&nbsp;</span>				dir, typ.Dir = typ.Dir, ast.RECV
<span id="L1817" class="ln">  1817&nbsp;&nbsp;</span>				typ, ok = typ.Value.(*ast.ChanType)
<span id="L1818" class="ln">  1818&nbsp;&nbsp;</span>			}
<span id="L1819" class="ln">  1819&nbsp;&nbsp;</span>			if dir == ast.SEND {
<span id="L1820" class="ln">  1820&nbsp;&nbsp;</span>				p.errorExpected(arrow, &#34;channel type&#34;)
<span id="L1821" class="ln">  1821&nbsp;&nbsp;</span>			}
<span id="L1822" class="ln">  1822&nbsp;&nbsp;</span>
<span id="L1823" class="ln">  1823&nbsp;&nbsp;</span>			return x
<span id="L1824" class="ln">  1824&nbsp;&nbsp;</span>		}
<span id="L1825" class="ln">  1825&nbsp;&nbsp;</span>
<span id="L1826" class="ln">  1826&nbsp;&nbsp;</span>		<span class="comment">// &lt;-(expr)</span>
<span id="L1827" class="ln">  1827&nbsp;&nbsp;</span>		return &amp;ast.UnaryExpr{OpPos: arrow, Op: token.ARROW, X: x}
<span id="L1828" class="ln">  1828&nbsp;&nbsp;</span>
<span id="L1829" class="ln">  1829&nbsp;&nbsp;</span>	case token.MUL:
<span id="L1830" class="ln">  1830&nbsp;&nbsp;</span>		<span class="comment">// pointer type or unary &#34;*&#34; expression</span>
<span id="L1831" class="ln">  1831&nbsp;&nbsp;</span>		pos := p.pos
<span id="L1832" class="ln">  1832&nbsp;&nbsp;</span>		p.next()
<span id="L1833" class="ln">  1833&nbsp;&nbsp;</span>		x := p.parseUnaryExpr()
<span id="L1834" class="ln">  1834&nbsp;&nbsp;</span>		return &amp;ast.StarExpr{Star: pos, X: x}
<span id="L1835" class="ln">  1835&nbsp;&nbsp;</span>	}
<span id="L1836" class="ln">  1836&nbsp;&nbsp;</span>
<span id="L1837" class="ln">  1837&nbsp;&nbsp;</span>	return p.parsePrimaryExpr(nil)
<span id="L1838" class="ln">  1838&nbsp;&nbsp;</span>}
<span id="L1839" class="ln">  1839&nbsp;&nbsp;</span>
<span id="L1840" class="ln">  1840&nbsp;&nbsp;</span>func (p *parser) tokPrec() (token.Token, int) {
<span id="L1841" class="ln">  1841&nbsp;&nbsp;</span>	tok := p.tok
<span id="L1842" class="ln">  1842&nbsp;&nbsp;</span>	if p.inRhs &amp;&amp; tok == token.ASSIGN {
<span id="L1843" class="ln">  1843&nbsp;&nbsp;</span>		tok = token.EQL
<span id="L1844" class="ln">  1844&nbsp;&nbsp;</span>	}
<span id="L1845" class="ln">  1845&nbsp;&nbsp;</span>	return tok, tok.Precedence()
<span id="L1846" class="ln">  1846&nbsp;&nbsp;</span>}
<span id="L1847" class="ln">  1847&nbsp;&nbsp;</span>
<span id="L1848" class="ln">  1848&nbsp;&nbsp;</span><span class="comment">// parseBinaryExpr parses a (possibly) binary expression.</span>
<span id="L1849" class="ln">  1849&nbsp;&nbsp;</span><span class="comment">// If x is non-nil, it is used as the left operand.</span>
<span id="L1850" class="ln">  1850&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1851" class="ln">  1851&nbsp;&nbsp;</span><span class="comment">// TODO(rfindley): parseBinaryExpr has become overloaded. Consider refactoring.</span>
<span id="L1852" class="ln">  1852&nbsp;&nbsp;</span>func (p *parser) parseBinaryExpr(x ast.Expr, prec1 int) ast.Expr {
<span id="L1853" class="ln">  1853&nbsp;&nbsp;</span>	if p.trace {
<span id="L1854" class="ln">  1854&nbsp;&nbsp;</span>		defer un(trace(p, &#34;BinaryExpr&#34;))
<span id="L1855" class="ln">  1855&nbsp;&nbsp;</span>	}
<span id="L1856" class="ln">  1856&nbsp;&nbsp;</span>
<span id="L1857" class="ln">  1857&nbsp;&nbsp;</span>	if x == nil {
<span id="L1858" class="ln">  1858&nbsp;&nbsp;</span>		x = p.parseUnaryExpr()
<span id="L1859" class="ln">  1859&nbsp;&nbsp;</span>	}
<span id="L1860" class="ln">  1860&nbsp;&nbsp;</span>	<span class="comment">// We track the nesting here rather than at the entry for the function,</span>
<span id="L1861" class="ln">  1861&nbsp;&nbsp;</span>	<span class="comment">// since it can iteratively produce a nested output, and we want to</span>
<span id="L1862" class="ln">  1862&nbsp;&nbsp;</span>	<span class="comment">// limit how deep a structure we generate.</span>
<span id="L1863" class="ln">  1863&nbsp;&nbsp;</span>	var n int
<span id="L1864" class="ln">  1864&nbsp;&nbsp;</span>	defer func() { p.nestLev -= n }()
<span id="L1865" class="ln">  1865&nbsp;&nbsp;</span>	for n = 1; ; n++ {
<span id="L1866" class="ln">  1866&nbsp;&nbsp;</span>		incNestLev(p)
<span id="L1867" class="ln">  1867&nbsp;&nbsp;</span>		op, oprec := p.tokPrec()
<span id="L1868" class="ln">  1868&nbsp;&nbsp;</span>		if oprec &lt; prec1 {
<span id="L1869" class="ln">  1869&nbsp;&nbsp;</span>			return x
<span id="L1870" class="ln">  1870&nbsp;&nbsp;</span>		}
<span id="L1871" class="ln">  1871&nbsp;&nbsp;</span>		pos := p.expect(op)
<span id="L1872" class="ln">  1872&nbsp;&nbsp;</span>		y := p.parseBinaryExpr(nil, oprec+1)
<span id="L1873" class="ln">  1873&nbsp;&nbsp;</span>		x = &amp;ast.BinaryExpr{X: x, OpPos: pos, Op: op, Y: y}
<span id="L1874" class="ln">  1874&nbsp;&nbsp;</span>	}
<span id="L1875" class="ln">  1875&nbsp;&nbsp;</span>}
<span id="L1876" class="ln">  1876&nbsp;&nbsp;</span>
<span id="L1877" class="ln">  1877&nbsp;&nbsp;</span><span class="comment">// The result may be a type or even a raw type ([...]int).</span>
<span id="L1878" class="ln">  1878&nbsp;&nbsp;</span>func (p *parser) parseExpr() ast.Expr {
<span id="L1879" class="ln">  1879&nbsp;&nbsp;</span>	if p.trace {
<span id="L1880" class="ln">  1880&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Expression&#34;))
<span id="L1881" class="ln">  1881&nbsp;&nbsp;</span>	}
<span id="L1882" class="ln">  1882&nbsp;&nbsp;</span>
<span id="L1883" class="ln">  1883&nbsp;&nbsp;</span>	return p.parseBinaryExpr(nil, token.LowestPrec+1)
<span id="L1884" class="ln">  1884&nbsp;&nbsp;</span>}
<span id="L1885" class="ln">  1885&nbsp;&nbsp;</span>
<span id="L1886" class="ln">  1886&nbsp;&nbsp;</span>func (p *parser) parseRhs() ast.Expr {
<span id="L1887" class="ln">  1887&nbsp;&nbsp;</span>	old := p.inRhs
<span id="L1888" class="ln">  1888&nbsp;&nbsp;</span>	p.inRhs = true
<span id="L1889" class="ln">  1889&nbsp;&nbsp;</span>	x := p.parseExpr()
<span id="L1890" class="ln">  1890&nbsp;&nbsp;</span>	p.inRhs = old
<span id="L1891" class="ln">  1891&nbsp;&nbsp;</span>	return x
<span id="L1892" class="ln">  1892&nbsp;&nbsp;</span>}
<span id="L1893" class="ln">  1893&nbsp;&nbsp;</span>
<span id="L1894" class="ln">  1894&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L1895" class="ln">  1895&nbsp;&nbsp;</span><span class="comment">// Statements</span>
<span id="L1896" class="ln">  1896&nbsp;&nbsp;</span>
<span id="L1897" class="ln">  1897&nbsp;&nbsp;</span><span class="comment">// Parsing modes for parseSimpleStmt.</span>
<span id="L1898" class="ln">  1898&nbsp;&nbsp;</span>const (
<span id="L1899" class="ln">  1899&nbsp;&nbsp;</span>	basic = iota
<span id="L1900" class="ln">  1900&nbsp;&nbsp;</span>	labelOk
<span id="L1901" class="ln">  1901&nbsp;&nbsp;</span>	rangeOk
<span id="L1902" class="ln">  1902&nbsp;&nbsp;</span>)
<span id="L1903" class="ln">  1903&nbsp;&nbsp;</span>
<span id="L1904" class="ln">  1904&nbsp;&nbsp;</span><span class="comment">// parseSimpleStmt returns true as 2nd result if it parsed the assignment</span>
<span id="L1905" class="ln">  1905&nbsp;&nbsp;</span><span class="comment">// of a range clause (with mode == rangeOk). The returned statement is an</span>
<span id="L1906" class="ln">  1906&nbsp;&nbsp;</span><span class="comment">// assignment with a right-hand side that is a single unary expression of</span>
<span id="L1907" class="ln">  1907&nbsp;&nbsp;</span><span class="comment">// the form &#34;range x&#34;. No guarantees are given for the left-hand side.</span>
<span id="L1908" class="ln">  1908&nbsp;&nbsp;</span>func (p *parser) parseSimpleStmt(mode int) (ast.Stmt, bool) {
<span id="L1909" class="ln">  1909&nbsp;&nbsp;</span>	if p.trace {
<span id="L1910" class="ln">  1910&nbsp;&nbsp;</span>		defer un(trace(p, &#34;SimpleStmt&#34;))
<span id="L1911" class="ln">  1911&nbsp;&nbsp;</span>	}
<span id="L1912" class="ln">  1912&nbsp;&nbsp;</span>
<span id="L1913" class="ln">  1913&nbsp;&nbsp;</span>	x := p.parseList(false)
<span id="L1914" class="ln">  1914&nbsp;&nbsp;</span>
<span id="L1915" class="ln">  1915&nbsp;&nbsp;</span>	switch p.tok {
<span id="L1916" class="ln">  1916&nbsp;&nbsp;</span>	case
<span id="L1917" class="ln">  1917&nbsp;&nbsp;</span>		token.DEFINE, token.ASSIGN, token.ADD_ASSIGN,
<span id="L1918" class="ln">  1918&nbsp;&nbsp;</span>		token.SUB_ASSIGN, token.MUL_ASSIGN, token.QUO_ASSIGN,
<span id="L1919" class="ln">  1919&nbsp;&nbsp;</span>		token.REM_ASSIGN, token.AND_ASSIGN, token.OR_ASSIGN,
<span id="L1920" class="ln">  1920&nbsp;&nbsp;</span>		token.XOR_ASSIGN, token.SHL_ASSIGN, token.SHR_ASSIGN, token.AND_NOT_ASSIGN:
<span id="L1921" class="ln">  1921&nbsp;&nbsp;</span>		<span class="comment">// assignment statement, possibly part of a range clause</span>
<span id="L1922" class="ln">  1922&nbsp;&nbsp;</span>		pos, tok := p.pos, p.tok
<span id="L1923" class="ln">  1923&nbsp;&nbsp;</span>		p.next()
<span id="L1924" class="ln">  1924&nbsp;&nbsp;</span>		var y []ast.Expr
<span id="L1925" class="ln">  1925&nbsp;&nbsp;</span>		isRange := false
<span id="L1926" class="ln">  1926&nbsp;&nbsp;</span>		if mode == rangeOk &amp;&amp; p.tok == token.RANGE &amp;&amp; (tok == token.DEFINE || tok == token.ASSIGN) {
<span id="L1927" class="ln">  1927&nbsp;&nbsp;</span>			pos := p.pos
<span id="L1928" class="ln">  1928&nbsp;&nbsp;</span>			p.next()
<span id="L1929" class="ln">  1929&nbsp;&nbsp;</span>			y = []ast.Expr{&amp;ast.UnaryExpr{OpPos: pos, Op: token.RANGE, X: p.parseRhs()}}
<span id="L1930" class="ln">  1930&nbsp;&nbsp;</span>			isRange = true
<span id="L1931" class="ln">  1931&nbsp;&nbsp;</span>		} else {
<span id="L1932" class="ln">  1932&nbsp;&nbsp;</span>			y = p.parseList(true)
<span id="L1933" class="ln">  1933&nbsp;&nbsp;</span>		}
<span id="L1934" class="ln">  1934&nbsp;&nbsp;</span>		return &amp;ast.AssignStmt{Lhs: x, TokPos: pos, Tok: tok, Rhs: y}, isRange
<span id="L1935" class="ln">  1935&nbsp;&nbsp;</span>	}
<span id="L1936" class="ln">  1936&nbsp;&nbsp;</span>
<span id="L1937" class="ln">  1937&nbsp;&nbsp;</span>	if len(x) &gt; 1 {
<span id="L1938" class="ln">  1938&nbsp;&nbsp;</span>		p.errorExpected(x[0].Pos(), &#34;1 expression&#34;)
<span id="L1939" class="ln">  1939&nbsp;&nbsp;</span>		<span class="comment">// continue with first expression</span>
<span id="L1940" class="ln">  1940&nbsp;&nbsp;</span>	}
<span id="L1941" class="ln">  1941&nbsp;&nbsp;</span>
<span id="L1942" class="ln">  1942&nbsp;&nbsp;</span>	switch p.tok {
<span id="L1943" class="ln">  1943&nbsp;&nbsp;</span>	case token.COLON:
<span id="L1944" class="ln">  1944&nbsp;&nbsp;</span>		<span class="comment">// labeled statement</span>
<span id="L1945" class="ln">  1945&nbsp;&nbsp;</span>		colon := p.pos
<span id="L1946" class="ln">  1946&nbsp;&nbsp;</span>		p.next()
<span id="L1947" class="ln">  1947&nbsp;&nbsp;</span>		if label, isIdent := x[0].(*ast.Ident); mode == labelOk &amp;&amp; isIdent {
<span id="L1948" class="ln">  1948&nbsp;&nbsp;</span>			<span class="comment">// Go spec: The scope of a label is the body of the function</span>
<span id="L1949" class="ln">  1949&nbsp;&nbsp;</span>			<span class="comment">// in which it is declared and excludes the body of any nested</span>
<span id="L1950" class="ln">  1950&nbsp;&nbsp;</span>			<span class="comment">// function.</span>
<span id="L1951" class="ln">  1951&nbsp;&nbsp;</span>			stmt := &amp;ast.LabeledStmt{Label: label, Colon: colon, Stmt: p.parseStmt()}
<span id="L1952" class="ln">  1952&nbsp;&nbsp;</span>			return stmt, false
<span id="L1953" class="ln">  1953&nbsp;&nbsp;</span>		}
<span id="L1954" class="ln">  1954&nbsp;&nbsp;</span>		<span class="comment">// The label declaration typically starts at x[0].Pos(), but the label</span>
<span id="L1955" class="ln">  1955&nbsp;&nbsp;</span>		<span class="comment">// declaration may be erroneous due to a token after that position (and</span>
<span id="L1956" class="ln">  1956&nbsp;&nbsp;</span>		<span class="comment">// before the &#39;:&#39;). If SpuriousErrors is not set, the (only) error</span>
<span id="L1957" class="ln">  1957&nbsp;&nbsp;</span>		<span class="comment">// reported for the line is the illegal label error instead of the token</span>
<span id="L1958" class="ln">  1958&nbsp;&nbsp;</span>		<span class="comment">// before the &#39;:&#39; that caused the problem. Thus, use the (latest) colon</span>
<span id="L1959" class="ln">  1959&nbsp;&nbsp;</span>		<span class="comment">// position for error reporting.</span>
<span id="L1960" class="ln">  1960&nbsp;&nbsp;</span>		p.error(colon, &#34;illegal label declaration&#34;)
<span id="L1961" class="ln">  1961&nbsp;&nbsp;</span>		return &amp;ast.BadStmt{From: x[0].Pos(), To: colon + 1}, false
<span id="L1962" class="ln">  1962&nbsp;&nbsp;</span>
<span id="L1963" class="ln">  1963&nbsp;&nbsp;</span>	case token.ARROW:
<span id="L1964" class="ln">  1964&nbsp;&nbsp;</span>		<span class="comment">// send statement</span>
<span id="L1965" class="ln">  1965&nbsp;&nbsp;</span>		arrow := p.pos
<span id="L1966" class="ln">  1966&nbsp;&nbsp;</span>		p.next()
<span id="L1967" class="ln">  1967&nbsp;&nbsp;</span>		y := p.parseRhs()
<span id="L1968" class="ln">  1968&nbsp;&nbsp;</span>		return &amp;ast.SendStmt{Chan: x[0], Arrow: arrow, Value: y}, false
<span id="L1969" class="ln">  1969&nbsp;&nbsp;</span>
<span id="L1970" class="ln">  1970&nbsp;&nbsp;</span>	case token.INC, token.DEC:
<span id="L1971" class="ln">  1971&nbsp;&nbsp;</span>		<span class="comment">// increment or decrement</span>
<span id="L1972" class="ln">  1972&nbsp;&nbsp;</span>		s := &amp;ast.IncDecStmt{X: x[0], TokPos: p.pos, Tok: p.tok}
<span id="L1973" class="ln">  1973&nbsp;&nbsp;</span>		p.next()
<span id="L1974" class="ln">  1974&nbsp;&nbsp;</span>		return s, false
<span id="L1975" class="ln">  1975&nbsp;&nbsp;</span>	}
<span id="L1976" class="ln">  1976&nbsp;&nbsp;</span>
<span id="L1977" class="ln">  1977&nbsp;&nbsp;</span>	<span class="comment">// expression</span>
<span id="L1978" class="ln">  1978&nbsp;&nbsp;</span>	return &amp;ast.ExprStmt{X: x[0]}, false
<span id="L1979" class="ln">  1979&nbsp;&nbsp;</span>}
<span id="L1980" class="ln">  1980&nbsp;&nbsp;</span>
<span id="L1981" class="ln">  1981&nbsp;&nbsp;</span>func (p *parser) parseCallExpr(callType string) *ast.CallExpr {
<span id="L1982" class="ln">  1982&nbsp;&nbsp;</span>	x := p.parseRhs() <span class="comment">// could be a conversion: (some type)(x)</span>
<span id="L1983" class="ln">  1983&nbsp;&nbsp;</span>	if t := ast.Unparen(x); t != x {
<span id="L1984" class="ln">  1984&nbsp;&nbsp;</span>		p.error(x.Pos(), fmt.Sprintf(&#34;expression in %s must not be parenthesized&#34;, callType))
<span id="L1985" class="ln">  1985&nbsp;&nbsp;</span>		x = t
<span id="L1986" class="ln">  1986&nbsp;&nbsp;</span>	}
<span id="L1987" class="ln">  1987&nbsp;&nbsp;</span>	if call, isCall := x.(*ast.CallExpr); isCall {
<span id="L1988" class="ln">  1988&nbsp;&nbsp;</span>		return call
<span id="L1989" class="ln">  1989&nbsp;&nbsp;</span>	}
<span id="L1990" class="ln">  1990&nbsp;&nbsp;</span>	if _, isBad := x.(*ast.BadExpr); !isBad {
<span id="L1991" class="ln">  1991&nbsp;&nbsp;</span>		<span class="comment">// only report error if it&#39;s a new one</span>
<span id="L1992" class="ln">  1992&nbsp;&nbsp;</span>		p.error(p.safePos(x.End()), fmt.Sprintf(&#34;expression in %s must be function call&#34;, callType))
<span id="L1993" class="ln">  1993&nbsp;&nbsp;</span>	}
<span id="L1994" class="ln">  1994&nbsp;&nbsp;</span>	return nil
<span id="L1995" class="ln">  1995&nbsp;&nbsp;</span>}
<span id="L1996" class="ln">  1996&nbsp;&nbsp;</span>
<span id="L1997" class="ln">  1997&nbsp;&nbsp;</span>func (p *parser) parseGoStmt() ast.Stmt {
<span id="L1998" class="ln">  1998&nbsp;&nbsp;</span>	if p.trace {
<span id="L1999" class="ln">  1999&nbsp;&nbsp;</span>		defer un(trace(p, &#34;GoStmt&#34;))
<span id="L2000" class="ln">  2000&nbsp;&nbsp;</span>	}
<span id="L2001" class="ln">  2001&nbsp;&nbsp;</span>
<span id="L2002" class="ln">  2002&nbsp;&nbsp;</span>	pos := p.expect(token.GO)
<span id="L2003" class="ln">  2003&nbsp;&nbsp;</span>	call := p.parseCallExpr(&#34;go&#34;)
<span id="L2004" class="ln">  2004&nbsp;&nbsp;</span>	p.expectSemi()
<span id="L2005" class="ln">  2005&nbsp;&nbsp;</span>	if call == nil {
<span id="L2006" class="ln">  2006&nbsp;&nbsp;</span>		return &amp;ast.BadStmt{From: pos, To: pos + 2} <span class="comment">// len(&#34;go&#34;)</span>
<span id="L2007" class="ln">  2007&nbsp;&nbsp;</span>	}
<span id="L2008" class="ln">  2008&nbsp;&nbsp;</span>
<span id="L2009" class="ln">  2009&nbsp;&nbsp;</span>	return &amp;ast.GoStmt{Go: pos, Call: call}
<span id="L2010" class="ln">  2010&nbsp;&nbsp;</span>}
<span id="L2011" class="ln">  2011&nbsp;&nbsp;</span>
<span id="L2012" class="ln">  2012&nbsp;&nbsp;</span>func (p *parser) parseDeferStmt() ast.Stmt {
<span id="L2013" class="ln">  2013&nbsp;&nbsp;</span>	if p.trace {
<span id="L2014" class="ln">  2014&nbsp;&nbsp;</span>		defer un(trace(p, &#34;DeferStmt&#34;))
<span id="L2015" class="ln">  2015&nbsp;&nbsp;</span>	}
<span id="L2016" class="ln">  2016&nbsp;&nbsp;</span>
<span id="L2017" class="ln">  2017&nbsp;&nbsp;</span>	pos := p.expect(token.DEFER)
<span id="L2018" class="ln">  2018&nbsp;&nbsp;</span>	call := p.parseCallExpr(&#34;defer&#34;)
<span id="L2019" class="ln">  2019&nbsp;&nbsp;</span>	p.expectSemi()
<span id="L2020" class="ln">  2020&nbsp;&nbsp;</span>	if call == nil {
<span id="L2021" class="ln">  2021&nbsp;&nbsp;</span>		return &amp;ast.BadStmt{From: pos, To: pos + 5} <span class="comment">// len(&#34;defer&#34;)</span>
<span id="L2022" class="ln">  2022&nbsp;&nbsp;</span>	}
<span id="L2023" class="ln">  2023&nbsp;&nbsp;</span>
<span id="L2024" class="ln">  2024&nbsp;&nbsp;</span>	return &amp;ast.DeferStmt{Defer: pos, Call: call}
<span id="L2025" class="ln">  2025&nbsp;&nbsp;</span>}
<span id="L2026" class="ln">  2026&nbsp;&nbsp;</span>
<span id="L2027" class="ln">  2027&nbsp;&nbsp;</span>func (p *parser) parseReturnStmt() *ast.ReturnStmt {
<span id="L2028" class="ln">  2028&nbsp;&nbsp;</span>	if p.trace {
<span id="L2029" class="ln">  2029&nbsp;&nbsp;</span>		defer un(trace(p, &#34;ReturnStmt&#34;))
<span id="L2030" class="ln">  2030&nbsp;&nbsp;</span>	}
<span id="L2031" class="ln">  2031&nbsp;&nbsp;</span>
<span id="L2032" class="ln">  2032&nbsp;&nbsp;</span>	pos := p.pos
<span id="L2033" class="ln">  2033&nbsp;&nbsp;</span>	p.expect(token.RETURN)
<span id="L2034" class="ln">  2034&nbsp;&nbsp;</span>	var x []ast.Expr
<span id="L2035" class="ln">  2035&nbsp;&nbsp;</span>	if p.tok != token.SEMICOLON &amp;&amp; p.tok != token.RBRACE {
<span id="L2036" class="ln">  2036&nbsp;&nbsp;</span>		x = p.parseList(true)
<span id="L2037" class="ln">  2037&nbsp;&nbsp;</span>	}
<span id="L2038" class="ln">  2038&nbsp;&nbsp;</span>	p.expectSemi()
<span id="L2039" class="ln">  2039&nbsp;&nbsp;</span>
<span id="L2040" class="ln">  2040&nbsp;&nbsp;</span>	return &amp;ast.ReturnStmt{Return: pos, Results: x}
<span id="L2041" class="ln">  2041&nbsp;&nbsp;</span>}
<span id="L2042" class="ln">  2042&nbsp;&nbsp;</span>
<span id="L2043" class="ln">  2043&nbsp;&nbsp;</span>func (p *parser) parseBranchStmt(tok token.Token) *ast.BranchStmt {
<span id="L2044" class="ln">  2044&nbsp;&nbsp;</span>	if p.trace {
<span id="L2045" class="ln">  2045&nbsp;&nbsp;</span>		defer un(trace(p, &#34;BranchStmt&#34;))
<span id="L2046" class="ln">  2046&nbsp;&nbsp;</span>	}
<span id="L2047" class="ln">  2047&nbsp;&nbsp;</span>
<span id="L2048" class="ln">  2048&nbsp;&nbsp;</span>	pos := p.expect(tok)
<span id="L2049" class="ln">  2049&nbsp;&nbsp;</span>	var label *ast.Ident
<span id="L2050" class="ln">  2050&nbsp;&nbsp;</span>	if tok != token.FALLTHROUGH &amp;&amp; p.tok == token.IDENT {
<span id="L2051" class="ln">  2051&nbsp;&nbsp;</span>		label = p.parseIdent()
<span id="L2052" class="ln">  2052&nbsp;&nbsp;</span>	}
<span id="L2053" class="ln">  2053&nbsp;&nbsp;</span>	p.expectSemi()
<span id="L2054" class="ln">  2054&nbsp;&nbsp;</span>
<span id="L2055" class="ln">  2055&nbsp;&nbsp;</span>	return &amp;ast.BranchStmt{TokPos: pos, Tok: tok, Label: label}
<span id="L2056" class="ln">  2056&nbsp;&nbsp;</span>}
<span id="L2057" class="ln">  2057&nbsp;&nbsp;</span>
<span id="L2058" class="ln">  2058&nbsp;&nbsp;</span>func (p *parser) makeExpr(s ast.Stmt, want string) ast.Expr {
<span id="L2059" class="ln">  2059&nbsp;&nbsp;</span>	if s == nil {
<span id="L2060" class="ln">  2060&nbsp;&nbsp;</span>		return nil
<span id="L2061" class="ln">  2061&nbsp;&nbsp;</span>	}
<span id="L2062" class="ln">  2062&nbsp;&nbsp;</span>	if es, isExpr := s.(*ast.ExprStmt); isExpr {
<span id="L2063" class="ln">  2063&nbsp;&nbsp;</span>		return es.X
<span id="L2064" class="ln">  2064&nbsp;&nbsp;</span>	}
<span id="L2065" class="ln">  2065&nbsp;&nbsp;</span>	found := &#34;simple statement&#34;
<span id="L2066" class="ln">  2066&nbsp;&nbsp;</span>	if _, isAss := s.(*ast.AssignStmt); isAss {
<span id="L2067" class="ln">  2067&nbsp;&nbsp;</span>		found = &#34;assignment&#34;
<span id="L2068" class="ln">  2068&nbsp;&nbsp;</span>	}
<span id="L2069" class="ln">  2069&nbsp;&nbsp;</span>	p.error(s.Pos(), fmt.Sprintf(&#34;expected %s, found %s (missing parentheses around composite literal?)&#34;, want, found))
<span id="L2070" class="ln">  2070&nbsp;&nbsp;</span>	return &amp;ast.BadExpr{From: s.Pos(), To: p.safePos(s.End())}
<span id="L2071" class="ln">  2071&nbsp;&nbsp;</span>}
<span id="L2072" class="ln">  2072&nbsp;&nbsp;</span>
<span id="L2073" class="ln">  2073&nbsp;&nbsp;</span><span class="comment">// parseIfHeader is an adjusted version of parser.header</span>
<span id="L2074" class="ln">  2074&nbsp;&nbsp;</span><span class="comment">// in cmd/compile/internal/syntax/parser.go, which has</span>
<span id="L2075" class="ln">  2075&nbsp;&nbsp;</span><span class="comment">// been tuned for better error handling.</span>
<span id="L2076" class="ln">  2076&nbsp;&nbsp;</span>func (p *parser) parseIfHeader() (init ast.Stmt, cond ast.Expr) {
<span id="L2077" class="ln">  2077&nbsp;&nbsp;</span>	if p.tok == token.LBRACE {
<span id="L2078" class="ln">  2078&nbsp;&nbsp;</span>		p.error(p.pos, &#34;missing condition in if statement&#34;)
<span id="L2079" class="ln">  2079&nbsp;&nbsp;</span>		cond = &amp;ast.BadExpr{From: p.pos, To: p.pos}
<span id="L2080" class="ln">  2080&nbsp;&nbsp;</span>		return
<span id="L2081" class="ln">  2081&nbsp;&nbsp;</span>	}
<span id="L2082" class="ln">  2082&nbsp;&nbsp;</span>	<span class="comment">// p.tok != token.LBRACE</span>
<span id="L2083" class="ln">  2083&nbsp;&nbsp;</span>
<span id="L2084" class="ln">  2084&nbsp;&nbsp;</span>	prevLev := p.exprLev
<span id="L2085" class="ln">  2085&nbsp;&nbsp;</span>	p.exprLev = -1
<span id="L2086" class="ln">  2086&nbsp;&nbsp;</span>
<span id="L2087" class="ln">  2087&nbsp;&nbsp;</span>	if p.tok != token.SEMICOLON {
<span id="L2088" class="ln">  2088&nbsp;&nbsp;</span>		<span class="comment">// accept potential variable declaration but complain</span>
<span id="L2089" class="ln">  2089&nbsp;&nbsp;</span>		if p.tok == token.VAR {
<span id="L2090" class="ln">  2090&nbsp;&nbsp;</span>			p.next()
<span id="L2091" class="ln">  2091&nbsp;&nbsp;</span>			p.error(p.pos, &#34;var declaration not allowed in if initializer&#34;)
<span id="L2092" class="ln">  2092&nbsp;&nbsp;</span>		}
<span id="L2093" class="ln">  2093&nbsp;&nbsp;</span>		init, _ = p.parseSimpleStmt(basic)
<span id="L2094" class="ln">  2094&nbsp;&nbsp;</span>	}
<span id="L2095" class="ln">  2095&nbsp;&nbsp;</span>
<span id="L2096" class="ln">  2096&nbsp;&nbsp;</span>	var condStmt ast.Stmt
<span id="L2097" class="ln">  2097&nbsp;&nbsp;</span>	var semi struct {
<span id="L2098" class="ln">  2098&nbsp;&nbsp;</span>		pos token.Pos
<span id="L2099" class="ln">  2099&nbsp;&nbsp;</span>		lit string <span class="comment">// &#34;;&#34; or &#34;\n&#34;; valid if pos.IsValid()</span>
<span id="L2100" class="ln">  2100&nbsp;&nbsp;</span>	}
<span id="L2101" class="ln">  2101&nbsp;&nbsp;</span>	if p.tok != token.LBRACE {
<span id="L2102" class="ln">  2102&nbsp;&nbsp;</span>		if p.tok == token.SEMICOLON {
<span id="L2103" class="ln">  2103&nbsp;&nbsp;</span>			semi.pos = p.pos
<span id="L2104" class="ln">  2104&nbsp;&nbsp;</span>			semi.lit = p.lit
<span id="L2105" class="ln">  2105&nbsp;&nbsp;</span>			p.next()
<span id="L2106" class="ln">  2106&nbsp;&nbsp;</span>		} else {
<span id="L2107" class="ln">  2107&nbsp;&nbsp;</span>			p.expect(token.SEMICOLON)
<span id="L2108" class="ln">  2108&nbsp;&nbsp;</span>		}
<span id="L2109" class="ln">  2109&nbsp;&nbsp;</span>		if p.tok != token.LBRACE {
<span id="L2110" class="ln">  2110&nbsp;&nbsp;</span>			condStmt, _ = p.parseSimpleStmt(basic)
<span id="L2111" class="ln">  2111&nbsp;&nbsp;</span>		}
<span id="L2112" class="ln">  2112&nbsp;&nbsp;</span>	} else {
<span id="L2113" class="ln">  2113&nbsp;&nbsp;</span>		condStmt = init
<span id="L2114" class="ln">  2114&nbsp;&nbsp;</span>		init = nil
<span id="L2115" class="ln">  2115&nbsp;&nbsp;</span>	}
<span id="L2116" class="ln">  2116&nbsp;&nbsp;</span>
<span id="L2117" class="ln">  2117&nbsp;&nbsp;</span>	if condStmt != nil {
<span id="L2118" class="ln">  2118&nbsp;&nbsp;</span>		cond = p.makeExpr(condStmt, &#34;boolean expression&#34;)
<span id="L2119" class="ln">  2119&nbsp;&nbsp;</span>	} else if semi.pos.IsValid() {
<span id="L2120" class="ln">  2120&nbsp;&nbsp;</span>		if semi.lit == &#34;\n&#34; {
<span id="L2121" class="ln">  2121&nbsp;&nbsp;</span>			p.error(semi.pos, &#34;unexpected newline, expecting { after if clause&#34;)
<span id="L2122" class="ln">  2122&nbsp;&nbsp;</span>		} else {
<span id="L2123" class="ln">  2123&nbsp;&nbsp;</span>			p.error(semi.pos, &#34;missing condition in if statement&#34;)
<span id="L2124" class="ln">  2124&nbsp;&nbsp;</span>		}
<span id="L2125" class="ln">  2125&nbsp;&nbsp;</span>	}
<span id="L2126" class="ln">  2126&nbsp;&nbsp;</span>
<span id="L2127" class="ln">  2127&nbsp;&nbsp;</span>	<span class="comment">// make sure we have a valid AST</span>
<span id="L2128" class="ln">  2128&nbsp;&nbsp;</span>	if cond == nil {
<span id="L2129" class="ln">  2129&nbsp;&nbsp;</span>		cond = &amp;ast.BadExpr{From: p.pos, To: p.pos}
<span id="L2130" class="ln">  2130&nbsp;&nbsp;</span>	}
<span id="L2131" class="ln">  2131&nbsp;&nbsp;</span>
<span id="L2132" class="ln">  2132&nbsp;&nbsp;</span>	p.exprLev = prevLev
<span id="L2133" class="ln">  2133&nbsp;&nbsp;</span>	return
<span id="L2134" class="ln">  2134&nbsp;&nbsp;</span>}
<span id="L2135" class="ln">  2135&nbsp;&nbsp;</span>
<span id="L2136" class="ln">  2136&nbsp;&nbsp;</span>func (p *parser) parseIfStmt() *ast.IfStmt {
<span id="L2137" class="ln">  2137&nbsp;&nbsp;</span>	defer decNestLev(incNestLev(p))
<span id="L2138" class="ln">  2138&nbsp;&nbsp;</span>
<span id="L2139" class="ln">  2139&nbsp;&nbsp;</span>	if p.trace {
<span id="L2140" class="ln">  2140&nbsp;&nbsp;</span>		defer un(trace(p, &#34;IfStmt&#34;))
<span id="L2141" class="ln">  2141&nbsp;&nbsp;</span>	}
<span id="L2142" class="ln">  2142&nbsp;&nbsp;</span>
<span id="L2143" class="ln">  2143&nbsp;&nbsp;</span>	pos := p.expect(token.IF)
<span id="L2144" class="ln">  2144&nbsp;&nbsp;</span>
<span id="L2145" class="ln">  2145&nbsp;&nbsp;</span>	init, cond := p.parseIfHeader()
<span id="L2146" class="ln">  2146&nbsp;&nbsp;</span>	body := p.parseBlockStmt()
<span id="L2147" class="ln">  2147&nbsp;&nbsp;</span>
<span id="L2148" class="ln">  2148&nbsp;&nbsp;</span>	var else_ ast.Stmt
<span id="L2149" class="ln">  2149&nbsp;&nbsp;</span>	if p.tok == token.ELSE {
<span id="L2150" class="ln">  2150&nbsp;&nbsp;</span>		p.next()
<span id="L2151" class="ln">  2151&nbsp;&nbsp;</span>		switch p.tok {
<span id="L2152" class="ln">  2152&nbsp;&nbsp;</span>		case token.IF:
<span id="L2153" class="ln">  2153&nbsp;&nbsp;</span>			else_ = p.parseIfStmt()
<span id="L2154" class="ln">  2154&nbsp;&nbsp;</span>		case token.LBRACE:
<span id="L2155" class="ln">  2155&nbsp;&nbsp;</span>			else_ = p.parseBlockStmt()
<span id="L2156" class="ln">  2156&nbsp;&nbsp;</span>			p.expectSemi()
<span id="L2157" class="ln">  2157&nbsp;&nbsp;</span>		default:
<span id="L2158" class="ln">  2158&nbsp;&nbsp;</span>			p.errorExpected(p.pos, &#34;if statement or block&#34;)
<span id="L2159" class="ln">  2159&nbsp;&nbsp;</span>			else_ = &amp;ast.BadStmt{From: p.pos, To: p.pos}
<span id="L2160" class="ln">  2160&nbsp;&nbsp;</span>		}
<span id="L2161" class="ln">  2161&nbsp;&nbsp;</span>	} else {
<span id="L2162" class="ln">  2162&nbsp;&nbsp;</span>		p.expectSemi()
<span id="L2163" class="ln">  2163&nbsp;&nbsp;</span>	}
<span id="L2164" class="ln">  2164&nbsp;&nbsp;</span>
<span id="L2165" class="ln">  2165&nbsp;&nbsp;</span>	return &amp;ast.IfStmt{If: pos, Init: init, Cond: cond, Body: body, Else: else_}
<span id="L2166" class="ln">  2166&nbsp;&nbsp;</span>}
<span id="L2167" class="ln">  2167&nbsp;&nbsp;</span>
<span id="L2168" class="ln">  2168&nbsp;&nbsp;</span>func (p *parser) parseCaseClause() *ast.CaseClause {
<span id="L2169" class="ln">  2169&nbsp;&nbsp;</span>	if p.trace {
<span id="L2170" class="ln">  2170&nbsp;&nbsp;</span>		defer un(trace(p, &#34;CaseClause&#34;))
<span id="L2171" class="ln">  2171&nbsp;&nbsp;</span>	}
<span id="L2172" class="ln">  2172&nbsp;&nbsp;</span>
<span id="L2173" class="ln">  2173&nbsp;&nbsp;</span>	pos := p.pos
<span id="L2174" class="ln">  2174&nbsp;&nbsp;</span>	var list []ast.Expr
<span id="L2175" class="ln">  2175&nbsp;&nbsp;</span>	if p.tok == token.CASE {
<span id="L2176" class="ln">  2176&nbsp;&nbsp;</span>		p.next()
<span id="L2177" class="ln">  2177&nbsp;&nbsp;</span>		list = p.parseList(true)
<span id="L2178" class="ln">  2178&nbsp;&nbsp;</span>	} else {
<span id="L2179" class="ln">  2179&nbsp;&nbsp;</span>		p.expect(token.DEFAULT)
<span id="L2180" class="ln">  2180&nbsp;&nbsp;</span>	}
<span id="L2181" class="ln">  2181&nbsp;&nbsp;</span>
<span id="L2182" class="ln">  2182&nbsp;&nbsp;</span>	colon := p.expect(token.COLON)
<span id="L2183" class="ln">  2183&nbsp;&nbsp;</span>	body := p.parseStmtList()
<span id="L2184" class="ln">  2184&nbsp;&nbsp;</span>
<span id="L2185" class="ln">  2185&nbsp;&nbsp;</span>	return &amp;ast.CaseClause{Case: pos, List: list, Colon: colon, Body: body}
<span id="L2186" class="ln">  2186&nbsp;&nbsp;</span>}
<span id="L2187" class="ln">  2187&nbsp;&nbsp;</span>
<span id="L2188" class="ln">  2188&nbsp;&nbsp;</span>func isTypeSwitchAssert(x ast.Expr) bool {
<span id="L2189" class="ln">  2189&nbsp;&nbsp;</span>	a, ok := x.(*ast.TypeAssertExpr)
<span id="L2190" class="ln">  2190&nbsp;&nbsp;</span>	return ok &amp;&amp; a.Type == nil
<span id="L2191" class="ln">  2191&nbsp;&nbsp;</span>}
<span id="L2192" class="ln">  2192&nbsp;&nbsp;</span>
<span id="L2193" class="ln">  2193&nbsp;&nbsp;</span>func (p *parser) isTypeSwitchGuard(s ast.Stmt) bool {
<span id="L2194" class="ln">  2194&nbsp;&nbsp;</span>	switch t := s.(type) {
<span id="L2195" class="ln">  2195&nbsp;&nbsp;</span>	case *ast.ExprStmt:
<span id="L2196" class="ln">  2196&nbsp;&nbsp;</span>		<span class="comment">// x.(type)</span>
<span id="L2197" class="ln">  2197&nbsp;&nbsp;</span>		return isTypeSwitchAssert(t.X)
<span id="L2198" class="ln">  2198&nbsp;&nbsp;</span>	case *ast.AssignStmt:
<span id="L2199" class="ln">  2199&nbsp;&nbsp;</span>		<span class="comment">// v := x.(type)</span>
<span id="L2200" class="ln">  2200&nbsp;&nbsp;</span>		if len(t.Lhs) == 1 &amp;&amp; len(t.Rhs) == 1 &amp;&amp; isTypeSwitchAssert(t.Rhs[0]) {
<span id="L2201" class="ln">  2201&nbsp;&nbsp;</span>			switch t.Tok {
<span id="L2202" class="ln">  2202&nbsp;&nbsp;</span>			case token.ASSIGN:
<span id="L2203" class="ln">  2203&nbsp;&nbsp;</span>				<span class="comment">// permit v = x.(type) but complain</span>
<span id="L2204" class="ln">  2204&nbsp;&nbsp;</span>				p.error(t.TokPos, &#34;expected &#39;:=&#39;, found &#39;=&#39;&#34;)
<span id="L2205" class="ln">  2205&nbsp;&nbsp;</span>				fallthrough
<span id="L2206" class="ln">  2206&nbsp;&nbsp;</span>			case token.DEFINE:
<span id="L2207" class="ln">  2207&nbsp;&nbsp;</span>				return true
<span id="L2208" class="ln">  2208&nbsp;&nbsp;</span>			}
<span id="L2209" class="ln">  2209&nbsp;&nbsp;</span>		}
<span id="L2210" class="ln">  2210&nbsp;&nbsp;</span>	}
<span id="L2211" class="ln">  2211&nbsp;&nbsp;</span>	return false
<span id="L2212" class="ln">  2212&nbsp;&nbsp;</span>}
<span id="L2213" class="ln">  2213&nbsp;&nbsp;</span>
<span id="L2214" class="ln">  2214&nbsp;&nbsp;</span>func (p *parser) parseSwitchStmt() ast.Stmt {
<span id="L2215" class="ln">  2215&nbsp;&nbsp;</span>	if p.trace {
<span id="L2216" class="ln">  2216&nbsp;&nbsp;</span>		defer un(trace(p, &#34;SwitchStmt&#34;))
<span id="L2217" class="ln">  2217&nbsp;&nbsp;</span>	}
<span id="L2218" class="ln">  2218&nbsp;&nbsp;</span>
<span id="L2219" class="ln">  2219&nbsp;&nbsp;</span>	pos := p.expect(token.SWITCH)
<span id="L2220" class="ln">  2220&nbsp;&nbsp;</span>
<span id="L2221" class="ln">  2221&nbsp;&nbsp;</span>	var s1, s2 ast.Stmt
<span id="L2222" class="ln">  2222&nbsp;&nbsp;</span>	if p.tok != token.LBRACE {
<span id="L2223" class="ln">  2223&nbsp;&nbsp;</span>		prevLev := p.exprLev
<span id="L2224" class="ln">  2224&nbsp;&nbsp;</span>		p.exprLev = -1
<span id="L2225" class="ln">  2225&nbsp;&nbsp;</span>		if p.tok != token.SEMICOLON {
<span id="L2226" class="ln">  2226&nbsp;&nbsp;</span>			s2, _ = p.parseSimpleStmt(basic)
<span id="L2227" class="ln">  2227&nbsp;&nbsp;</span>		}
<span id="L2228" class="ln">  2228&nbsp;&nbsp;</span>		if p.tok == token.SEMICOLON {
<span id="L2229" class="ln">  2229&nbsp;&nbsp;</span>			p.next()
<span id="L2230" class="ln">  2230&nbsp;&nbsp;</span>			s1 = s2
<span id="L2231" class="ln">  2231&nbsp;&nbsp;</span>			s2 = nil
<span id="L2232" class="ln">  2232&nbsp;&nbsp;</span>			if p.tok != token.LBRACE {
<span id="L2233" class="ln">  2233&nbsp;&nbsp;</span>				<span class="comment">// A TypeSwitchGuard may declare a variable in addition</span>
<span id="L2234" class="ln">  2234&nbsp;&nbsp;</span>				<span class="comment">// to the variable declared in the initial SimpleStmt.</span>
<span id="L2235" class="ln">  2235&nbsp;&nbsp;</span>				<span class="comment">// Introduce extra scope to avoid redeclaration errors:</span>
<span id="L2236" class="ln">  2236&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L2237" class="ln">  2237&nbsp;&nbsp;</span>				<span class="comment">//	switch t := 0; t := x.(T) { ... }</span>
<span id="L2238" class="ln">  2238&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L2239" class="ln">  2239&nbsp;&nbsp;</span>				<span class="comment">// (this code is not valid Go because the first t</span>
<span id="L2240" class="ln">  2240&nbsp;&nbsp;</span>				<span class="comment">// cannot be accessed and thus is never used, the extra</span>
<span id="L2241" class="ln">  2241&nbsp;&nbsp;</span>				<span class="comment">// scope is needed for the correct error message).</span>
<span id="L2242" class="ln">  2242&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L2243" class="ln">  2243&nbsp;&nbsp;</span>				<span class="comment">// If we don&#39;t have a type switch, s2 must be an expression.</span>
<span id="L2244" class="ln">  2244&nbsp;&nbsp;</span>				<span class="comment">// Having the extra nested but empty scope won&#39;t affect it.</span>
<span id="L2245" class="ln">  2245&nbsp;&nbsp;</span>				s2, _ = p.parseSimpleStmt(basic)
<span id="L2246" class="ln">  2246&nbsp;&nbsp;</span>			}
<span id="L2247" class="ln">  2247&nbsp;&nbsp;</span>		}
<span id="L2248" class="ln">  2248&nbsp;&nbsp;</span>		p.exprLev = prevLev
<span id="L2249" class="ln">  2249&nbsp;&nbsp;</span>	}
<span id="L2250" class="ln">  2250&nbsp;&nbsp;</span>
<span id="L2251" class="ln">  2251&nbsp;&nbsp;</span>	typeSwitch := p.isTypeSwitchGuard(s2)
<span id="L2252" class="ln">  2252&nbsp;&nbsp;</span>	lbrace := p.expect(token.LBRACE)
<span id="L2253" class="ln">  2253&nbsp;&nbsp;</span>	var list []ast.Stmt
<span id="L2254" class="ln">  2254&nbsp;&nbsp;</span>	for p.tok == token.CASE || p.tok == token.DEFAULT {
<span id="L2255" class="ln">  2255&nbsp;&nbsp;</span>		list = append(list, p.parseCaseClause())
<span id="L2256" class="ln">  2256&nbsp;&nbsp;</span>	}
<span id="L2257" class="ln">  2257&nbsp;&nbsp;</span>	rbrace := p.expect(token.RBRACE)
<span id="L2258" class="ln">  2258&nbsp;&nbsp;</span>	p.expectSemi()
<span id="L2259" class="ln">  2259&nbsp;&nbsp;</span>	body := &amp;ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
<span id="L2260" class="ln">  2260&nbsp;&nbsp;</span>
<span id="L2261" class="ln">  2261&nbsp;&nbsp;</span>	if typeSwitch {
<span id="L2262" class="ln">  2262&nbsp;&nbsp;</span>		return &amp;ast.TypeSwitchStmt{Switch: pos, Init: s1, Assign: s2, Body: body}
<span id="L2263" class="ln">  2263&nbsp;&nbsp;</span>	}
<span id="L2264" class="ln">  2264&nbsp;&nbsp;</span>
<span id="L2265" class="ln">  2265&nbsp;&nbsp;</span>	return &amp;ast.SwitchStmt{Switch: pos, Init: s1, Tag: p.makeExpr(s2, &#34;switch expression&#34;), Body: body}
<span id="L2266" class="ln">  2266&nbsp;&nbsp;</span>}
<span id="L2267" class="ln">  2267&nbsp;&nbsp;</span>
<span id="L2268" class="ln">  2268&nbsp;&nbsp;</span>func (p *parser) parseCommClause() *ast.CommClause {
<span id="L2269" class="ln">  2269&nbsp;&nbsp;</span>	if p.trace {
<span id="L2270" class="ln">  2270&nbsp;&nbsp;</span>		defer un(trace(p, &#34;CommClause&#34;))
<span id="L2271" class="ln">  2271&nbsp;&nbsp;</span>	}
<span id="L2272" class="ln">  2272&nbsp;&nbsp;</span>
<span id="L2273" class="ln">  2273&nbsp;&nbsp;</span>	pos := p.pos
<span id="L2274" class="ln">  2274&nbsp;&nbsp;</span>	var comm ast.Stmt
<span id="L2275" class="ln">  2275&nbsp;&nbsp;</span>	if p.tok == token.CASE {
<span id="L2276" class="ln">  2276&nbsp;&nbsp;</span>		p.next()
<span id="L2277" class="ln">  2277&nbsp;&nbsp;</span>		lhs := p.parseList(false)
<span id="L2278" class="ln">  2278&nbsp;&nbsp;</span>		if p.tok == token.ARROW {
<span id="L2279" class="ln">  2279&nbsp;&nbsp;</span>			<span class="comment">// SendStmt</span>
<span id="L2280" class="ln">  2280&nbsp;&nbsp;</span>			if len(lhs) &gt; 1 {
<span id="L2281" class="ln">  2281&nbsp;&nbsp;</span>				p.errorExpected(lhs[0].Pos(), &#34;1 expression&#34;)
<span id="L2282" class="ln">  2282&nbsp;&nbsp;</span>				<span class="comment">// continue with first expression</span>
<span id="L2283" class="ln">  2283&nbsp;&nbsp;</span>			}
<span id="L2284" class="ln">  2284&nbsp;&nbsp;</span>			arrow := p.pos
<span id="L2285" class="ln">  2285&nbsp;&nbsp;</span>			p.next()
<span id="L2286" class="ln">  2286&nbsp;&nbsp;</span>			rhs := p.parseRhs()
<span id="L2287" class="ln">  2287&nbsp;&nbsp;</span>			comm = &amp;ast.SendStmt{Chan: lhs[0], Arrow: arrow, Value: rhs}
<span id="L2288" class="ln">  2288&nbsp;&nbsp;</span>		} else {
<span id="L2289" class="ln">  2289&nbsp;&nbsp;</span>			<span class="comment">// RecvStmt</span>
<span id="L2290" class="ln">  2290&nbsp;&nbsp;</span>			if tok := p.tok; tok == token.ASSIGN || tok == token.DEFINE {
<span id="L2291" class="ln">  2291&nbsp;&nbsp;</span>				<span class="comment">// RecvStmt with assignment</span>
<span id="L2292" class="ln">  2292&nbsp;&nbsp;</span>				if len(lhs) &gt; 2 {
<span id="L2293" class="ln">  2293&nbsp;&nbsp;</span>					p.errorExpected(lhs[0].Pos(), &#34;1 or 2 expressions&#34;)
<span id="L2294" class="ln">  2294&nbsp;&nbsp;</span>					<span class="comment">// continue with first two expressions</span>
<span id="L2295" class="ln">  2295&nbsp;&nbsp;</span>					lhs = lhs[0:2]
<span id="L2296" class="ln">  2296&nbsp;&nbsp;</span>				}
<span id="L2297" class="ln">  2297&nbsp;&nbsp;</span>				pos := p.pos
<span id="L2298" class="ln">  2298&nbsp;&nbsp;</span>				p.next()
<span id="L2299" class="ln">  2299&nbsp;&nbsp;</span>				rhs := p.parseRhs()
<span id="L2300" class="ln">  2300&nbsp;&nbsp;</span>				comm = &amp;ast.AssignStmt{Lhs: lhs, TokPos: pos, Tok: tok, Rhs: []ast.Expr{rhs}}
<span id="L2301" class="ln">  2301&nbsp;&nbsp;</span>			} else {
<span id="L2302" class="ln">  2302&nbsp;&nbsp;</span>				<span class="comment">// lhs must be single receive operation</span>
<span id="L2303" class="ln">  2303&nbsp;&nbsp;</span>				if len(lhs) &gt; 1 {
<span id="L2304" class="ln">  2304&nbsp;&nbsp;</span>					p.errorExpected(lhs[0].Pos(), &#34;1 expression&#34;)
<span id="L2305" class="ln">  2305&nbsp;&nbsp;</span>					<span class="comment">// continue with first expression</span>
<span id="L2306" class="ln">  2306&nbsp;&nbsp;</span>				}
<span id="L2307" class="ln">  2307&nbsp;&nbsp;</span>				comm = &amp;ast.ExprStmt{X: lhs[0]}
<span id="L2308" class="ln">  2308&nbsp;&nbsp;</span>			}
<span id="L2309" class="ln">  2309&nbsp;&nbsp;</span>		}
<span id="L2310" class="ln">  2310&nbsp;&nbsp;</span>	} else {
<span id="L2311" class="ln">  2311&nbsp;&nbsp;</span>		p.expect(token.DEFAULT)
<span id="L2312" class="ln">  2312&nbsp;&nbsp;</span>	}
<span id="L2313" class="ln">  2313&nbsp;&nbsp;</span>
<span id="L2314" class="ln">  2314&nbsp;&nbsp;</span>	colon := p.expect(token.COLON)
<span id="L2315" class="ln">  2315&nbsp;&nbsp;</span>	body := p.parseStmtList()
<span id="L2316" class="ln">  2316&nbsp;&nbsp;</span>
<span id="L2317" class="ln">  2317&nbsp;&nbsp;</span>	return &amp;ast.CommClause{Case: pos, Comm: comm, Colon: colon, Body: body}
<span id="L2318" class="ln">  2318&nbsp;&nbsp;</span>}
<span id="L2319" class="ln">  2319&nbsp;&nbsp;</span>
<span id="L2320" class="ln">  2320&nbsp;&nbsp;</span>func (p *parser) parseSelectStmt() *ast.SelectStmt {
<span id="L2321" class="ln">  2321&nbsp;&nbsp;</span>	if p.trace {
<span id="L2322" class="ln">  2322&nbsp;&nbsp;</span>		defer un(trace(p, &#34;SelectStmt&#34;))
<span id="L2323" class="ln">  2323&nbsp;&nbsp;</span>	}
<span id="L2324" class="ln">  2324&nbsp;&nbsp;</span>
<span id="L2325" class="ln">  2325&nbsp;&nbsp;</span>	pos := p.expect(token.SELECT)
<span id="L2326" class="ln">  2326&nbsp;&nbsp;</span>	lbrace := p.expect(token.LBRACE)
<span id="L2327" class="ln">  2327&nbsp;&nbsp;</span>	var list []ast.Stmt
<span id="L2328" class="ln">  2328&nbsp;&nbsp;</span>	for p.tok == token.CASE || p.tok == token.DEFAULT {
<span id="L2329" class="ln">  2329&nbsp;&nbsp;</span>		list = append(list, p.parseCommClause())
<span id="L2330" class="ln">  2330&nbsp;&nbsp;</span>	}
<span id="L2331" class="ln">  2331&nbsp;&nbsp;</span>	rbrace := p.expect(token.RBRACE)
<span id="L2332" class="ln">  2332&nbsp;&nbsp;</span>	p.expectSemi()
<span id="L2333" class="ln">  2333&nbsp;&nbsp;</span>	body := &amp;ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
<span id="L2334" class="ln">  2334&nbsp;&nbsp;</span>
<span id="L2335" class="ln">  2335&nbsp;&nbsp;</span>	return &amp;ast.SelectStmt{Select: pos, Body: body}
<span id="L2336" class="ln">  2336&nbsp;&nbsp;</span>}
<span id="L2337" class="ln">  2337&nbsp;&nbsp;</span>
<span id="L2338" class="ln">  2338&nbsp;&nbsp;</span>func (p *parser) parseForStmt() ast.Stmt {
<span id="L2339" class="ln">  2339&nbsp;&nbsp;</span>	if p.trace {
<span id="L2340" class="ln">  2340&nbsp;&nbsp;</span>		defer un(trace(p, &#34;ForStmt&#34;))
<span id="L2341" class="ln">  2341&nbsp;&nbsp;</span>	}
<span id="L2342" class="ln">  2342&nbsp;&nbsp;</span>
<span id="L2343" class="ln">  2343&nbsp;&nbsp;</span>	pos := p.expect(token.FOR)
<span id="L2344" class="ln">  2344&nbsp;&nbsp;</span>
<span id="L2345" class="ln">  2345&nbsp;&nbsp;</span>	var s1, s2, s3 ast.Stmt
<span id="L2346" class="ln">  2346&nbsp;&nbsp;</span>	var isRange bool
<span id="L2347" class="ln">  2347&nbsp;&nbsp;</span>	if p.tok != token.LBRACE {
<span id="L2348" class="ln">  2348&nbsp;&nbsp;</span>		prevLev := p.exprLev
<span id="L2349" class="ln">  2349&nbsp;&nbsp;</span>		p.exprLev = -1
<span id="L2350" class="ln">  2350&nbsp;&nbsp;</span>		if p.tok != token.SEMICOLON {
<span id="L2351" class="ln">  2351&nbsp;&nbsp;</span>			if p.tok == token.RANGE {
<span id="L2352" class="ln">  2352&nbsp;&nbsp;</span>				<span class="comment">// &#34;for range x&#34; (nil lhs in assignment)</span>
<span id="L2353" class="ln">  2353&nbsp;&nbsp;</span>				pos := p.pos
<span id="L2354" class="ln">  2354&nbsp;&nbsp;</span>				p.next()
<span id="L2355" class="ln">  2355&nbsp;&nbsp;</span>				y := []ast.Expr{&amp;ast.UnaryExpr{OpPos: pos, Op: token.RANGE, X: p.parseRhs()}}
<span id="L2356" class="ln">  2356&nbsp;&nbsp;</span>				s2 = &amp;ast.AssignStmt{Rhs: y}
<span id="L2357" class="ln">  2357&nbsp;&nbsp;</span>				isRange = true
<span id="L2358" class="ln">  2358&nbsp;&nbsp;</span>			} else {
<span id="L2359" class="ln">  2359&nbsp;&nbsp;</span>				s2, isRange = p.parseSimpleStmt(rangeOk)
<span id="L2360" class="ln">  2360&nbsp;&nbsp;</span>			}
<span id="L2361" class="ln">  2361&nbsp;&nbsp;</span>		}
<span id="L2362" class="ln">  2362&nbsp;&nbsp;</span>		if !isRange &amp;&amp; p.tok == token.SEMICOLON {
<span id="L2363" class="ln">  2363&nbsp;&nbsp;</span>			p.next()
<span id="L2364" class="ln">  2364&nbsp;&nbsp;</span>			s1 = s2
<span id="L2365" class="ln">  2365&nbsp;&nbsp;</span>			s2 = nil
<span id="L2366" class="ln">  2366&nbsp;&nbsp;</span>			if p.tok != token.SEMICOLON {
<span id="L2367" class="ln">  2367&nbsp;&nbsp;</span>				s2, _ = p.parseSimpleStmt(basic)
<span id="L2368" class="ln">  2368&nbsp;&nbsp;</span>			}
<span id="L2369" class="ln">  2369&nbsp;&nbsp;</span>			p.expectSemi()
<span id="L2370" class="ln">  2370&nbsp;&nbsp;</span>			if p.tok != token.LBRACE {
<span id="L2371" class="ln">  2371&nbsp;&nbsp;</span>				s3, _ = p.parseSimpleStmt(basic)
<span id="L2372" class="ln">  2372&nbsp;&nbsp;</span>			}
<span id="L2373" class="ln">  2373&nbsp;&nbsp;</span>		}
<span id="L2374" class="ln">  2374&nbsp;&nbsp;</span>		p.exprLev = prevLev
<span id="L2375" class="ln">  2375&nbsp;&nbsp;</span>	}
<span id="L2376" class="ln">  2376&nbsp;&nbsp;</span>
<span id="L2377" class="ln">  2377&nbsp;&nbsp;</span>	body := p.parseBlockStmt()
<span id="L2378" class="ln">  2378&nbsp;&nbsp;</span>	p.expectSemi()
<span id="L2379" class="ln">  2379&nbsp;&nbsp;</span>
<span id="L2380" class="ln">  2380&nbsp;&nbsp;</span>	if isRange {
<span id="L2381" class="ln">  2381&nbsp;&nbsp;</span>		as := s2.(*ast.AssignStmt)
<span id="L2382" class="ln">  2382&nbsp;&nbsp;</span>		<span class="comment">// check lhs</span>
<span id="L2383" class="ln">  2383&nbsp;&nbsp;</span>		var key, value ast.Expr
<span id="L2384" class="ln">  2384&nbsp;&nbsp;</span>		switch len(as.Lhs) {
<span id="L2385" class="ln">  2385&nbsp;&nbsp;</span>		case 0:
<span id="L2386" class="ln">  2386&nbsp;&nbsp;</span>			<span class="comment">// nothing to do</span>
<span id="L2387" class="ln">  2387&nbsp;&nbsp;</span>		case 1:
<span id="L2388" class="ln">  2388&nbsp;&nbsp;</span>			key = as.Lhs[0]
<span id="L2389" class="ln">  2389&nbsp;&nbsp;</span>		case 2:
<span id="L2390" class="ln">  2390&nbsp;&nbsp;</span>			key, value = as.Lhs[0], as.Lhs[1]
<span id="L2391" class="ln">  2391&nbsp;&nbsp;</span>		default:
<span id="L2392" class="ln">  2392&nbsp;&nbsp;</span>			p.errorExpected(as.Lhs[len(as.Lhs)-1].Pos(), &#34;at most 2 expressions&#34;)
<span id="L2393" class="ln">  2393&nbsp;&nbsp;</span>			return &amp;ast.BadStmt{From: pos, To: p.safePos(body.End())}
<span id="L2394" class="ln">  2394&nbsp;&nbsp;</span>		}
<span id="L2395" class="ln">  2395&nbsp;&nbsp;</span>		<span class="comment">// parseSimpleStmt returned a right-hand side that</span>
<span id="L2396" class="ln">  2396&nbsp;&nbsp;</span>		<span class="comment">// is a single unary expression of the form &#34;range x&#34;</span>
<span id="L2397" class="ln">  2397&nbsp;&nbsp;</span>		x := as.Rhs[0].(*ast.UnaryExpr).X
<span id="L2398" class="ln">  2398&nbsp;&nbsp;</span>		return &amp;ast.RangeStmt{
<span id="L2399" class="ln">  2399&nbsp;&nbsp;</span>			For:    pos,
<span id="L2400" class="ln">  2400&nbsp;&nbsp;</span>			Key:    key,
<span id="L2401" class="ln">  2401&nbsp;&nbsp;</span>			Value:  value,
<span id="L2402" class="ln">  2402&nbsp;&nbsp;</span>			TokPos: as.TokPos,
<span id="L2403" class="ln">  2403&nbsp;&nbsp;</span>			Tok:    as.Tok,
<span id="L2404" class="ln">  2404&nbsp;&nbsp;</span>			Range:  as.Rhs[0].Pos(),
<span id="L2405" class="ln">  2405&nbsp;&nbsp;</span>			X:      x,
<span id="L2406" class="ln">  2406&nbsp;&nbsp;</span>			Body:   body,
<span id="L2407" class="ln">  2407&nbsp;&nbsp;</span>		}
<span id="L2408" class="ln">  2408&nbsp;&nbsp;</span>	}
<span id="L2409" class="ln">  2409&nbsp;&nbsp;</span>
<span id="L2410" class="ln">  2410&nbsp;&nbsp;</span>	<span class="comment">// regular for statement</span>
<span id="L2411" class="ln">  2411&nbsp;&nbsp;</span>	return &amp;ast.ForStmt{
<span id="L2412" class="ln">  2412&nbsp;&nbsp;</span>		For:  pos,
<span id="L2413" class="ln">  2413&nbsp;&nbsp;</span>		Init: s1,
<span id="L2414" class="ln">  2414&nbsp;&nbsp;</span>		Cond: p.makeExpr(s2, &#34;boolean or range expression&#34;),
<span id="L2415" class="ln">  2415&nbsp;&nbsp;</span>		Post: s3,
<span id="L2416" class="ln">  2416&nbsp;&nbsp;</span>		Body: body,
<span id="L2417" class="ln">  2417&nbsp;&nbsp;</span>	}
<span id="L2418" class="ln">  2418&nbsp;&nbsp;</span>}
<span id="L2419" class="ln">  2419&nbsp;&nbsp;</span>
<span id="L2420" class="ln">  2420&nbsp;&nbsp;</span>func (p *parser) parseStmt() (s ast.Stmt) {
<span id="L2421" class="ln">  2421&nbsp;&nbsp;</span>	defer decNestLev(incNestLev(p))
<span id="L2422" class="ln">  2422&nbsp;&nbsp;</span>
<span id="L2423" class="ln">  2423&nbsp;&nbsp;</span>	if p.trace {
<span id="L2424" class="ln">  2424&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Statement&#34;))
<span id="L2425" class="ln">  2425&nbsp;&nbsp;</span>	}
<span id="L2426" class="ln">  2426&nbsp;&nbsp;</span>
<span id="L2427" class="ln">  2427&nbsp;&nbsp;</span>	switch p.tok {
<span id="L2428" class="ln">  2428&nbsp;&nbsp;</span>	case token.CONST, token.TYPE, token.VAR:
<span id="L2429" class="ln">  2429&nbsp;&nbsp;</span>		s = &amp;ast.DeclStmt{Decl: p.parseDecl(stmtStart)}
<span id="L2430" class="ln">  2430&nbsp;&nbsp;</span>	case
<span id="L2431" class="ln">  2431&nbsp;&nbsp;</span>		<span class="comment">// tokens that may start an expression</span>
<span id="L2432" class="ln">  2432&nbsp;&nbsp;</span>		token.IDENT, token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING, token.FUNC, token.LPAREN, <span class="comment">// operands</span>
<span id="L2433" class="ln">  2433&nbsp;&nbsp;</span>		token.LBRACK, token.STRUCT, token.MAP, token.CHAN, token.INTERFACE, <span class="comment">// composite types</span>
<span id="L2434" class="ln">  2434&nbsp;&nbsp;</span>		token.ADD, token.SUB, token.MUL, token.AND, token.XOR, token.ARROW, token.NOT: <span class="comment">// unary operators</span>
<span id="L2435" class="ln">  2435&nbsp;&nbsp;</span>		s, _ = p.parseSimpleStmt(labelOk)
<span id="L2436" class="ln">  2436&nbsp;&nbsp;</span>		<span class="comment">// because of the required look-ahead, labeled statements are</span>
<span id="L2437" class="ln">  2437&nbsp;&nbsp;</span>		<span class="comment">// parsed by parseSimpleStmt - don&#39;t expect a semicolon after</span>
<span id="L2438" class="ln">  2438&nbsp;&nbsp;</span>		<span class="comment">// them</span>
<span id="L2439" class="ln">  2439&nbsp;&nbsp;</span>		if _, isLabeledStmt := s.(*ast.LabeledStmt); !isLabeledStmt {
<span id="L2440" class="ln">  2440&nbsp;&nbsp;</span>			p.expectSemi()
<span id="L2441" class="ln">  2441&nbsp;&nbsp;</span>		}
<span id="L2442" class="ln">  2442&nbsp;&nbsp;</span>	case token.GO:
<span id="L2443" class="ln">  2443&nbsp;&nbsp;</span>		s = p.parseGoStmt()
<span id="L2444" class="ln">  2444&nbsp;&nbsp;</span>	case token.DEFER:
<span id="L2445" class="ln">  2445&nbsp;&nbsp;</span>		s = p.parseDeferStmt()
<span id="L2446" class="ln">  2446&nbsp;&nbsp;</span>	case token.RETURN:
<span id="L2447" class="ln">  2447&nbsp;&nbsp;</span>		s = p.parseReturnStmt()
<span id="L2448" class="ln">  2448&nbsp;&nbsp;</span>	case token.BREAK, token.CONTINUE, token.GOTO, token.FALLTHROUGH:
<span id="L2449" class="ln">  2449&nbsp;&nbsp;</span>		s = p.parseBranchStmt(p.tok)
<span id="L2450" class="ln">  2450&nbsp;&nbsp;</span>	case token.LBRACE:
<span id="L2451" class="ln">  2451&nbsp;&nbsp;</span>		s = p.parseBlockStmt()
<span id="L2452" class="ln">  2452&nbsp;&nbsp;</span>		p.expectSemi()
<span id="L2453" class="ln">  2453&nbsp;&nbsp;</span>	case token.IF:
<span id="L2454" class="ln">  2454&nbsp;&nbsp;</span>		s = p.parseIfStmt()
<span id="L2455" class="ln">  2455&nbsp;&nbsp;</span>	case token.SWITCH:
<span id="L2456" class="ln">  2456&nbsp;&nbsp;</span>		s = p.parseSwitchStmt()
<span id="L2457" class="ln">  2457&nbsp;&nbsp;</span>	case token.SELECT:
<span id="L2458" class="ln">  2458&nbsp;&nbsp;</span>		s = p.parseSelectStmt()
<span id="L2459" class="ln">  2459&nbsp;&nbsp;</span>	case token.FOR:
<span id="L2460" class="ln">  2460&nbsp;&nbsp;</span>		s = p.parseForStmt()
<span id="L2461" class="ln">  2461&nbsp;&nbsp;</span>	case token.SEMICOLON:
<span id="L2462" class="ln">  2462&nbsp;&nbsp;</span>		<span class="comment">// Is it ever possible to have an implicit semicolon</span>
<span id="L2463" class="ln">  2463&nbsp;&nbsp;</span>		<span class="comment">// producing an empty statement in a valid program?</span>
<span id="L2464" class="ln">  2464&nbsp;&nbsp;</span>		<span class="comment">// (handle correctly anyway)</span>
<span id="L2465" class="ln">  2465&nbsp;&nbsp;</span>		s = &amp;ast.EmptyStmt{Semicolon: p.pos, Implicit: p.lit == &#34;\n&#34;}
<span id="L2466" class="ln">  2466&nbsp;&nbsp;</span>		p.next()
<span id="L2467" class="ln">  2467&nbsp;&nbsp;</span>	case token.RBRACE:
<span id="L2468" class="ln">  2468&nbsp;&nbsp;</span>		<span class="comment">// a semicolon may be omitted before a closing &#34;}&#34;</span>
<span id="L2469" class="ln">  2469&nbsp;&nbsp;</span>		s = &amp;ast.EmptyStmt{Semicolon: p.pos, Implicit: true}
<span id="L2470" class="ln">  2470&nbsp;&nbsp;</span>	default:
<span id="L2471" class="ln">  2471&nbsp;&nbsp;</span>		<span class="comment">// no statement found</span>
<span id="L2472" class="ln">  2472&nbsp;&nbsp;</span>		pos := p.pos
<span id="L2473" class="ln">  2473&nbsp;&nbsp;</span>		p.errorExpected(pos, &#34;statement&#34;)
<span id="L2474" class="ln">  2474&nbsp;&nbsp;</span>		p.advance(stmtStart)
<span id="L2475" class="ln">  2475&nbsp;&nbsp;</span>		s = &amp;ast.BadStmt{From: pos, To: p.pos}
<span id="L2476" class="ln">  2476&nbsp;&nbsp;</span>	}
<span id="L2477" class="ln">  2477&nbsp;&nbsp;</span>
<span id="L2478" class="ln">  2478&nbsp;&nbsp;</span>	return
<span id="L2479" class="ln">  2479&nbsp;&nbsp;</span>}
<span id="L2480" class="ln">  2480&nbsp;&nbsp;</span>
<span id="L2481" class="ln">  2481&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L2482" class="ln">  2482&nbsp;&nbsp;</span><span class="comment">// Declarations</span>
<span id="L2483" class="ln">  2483&nbsp;&nbsp;</span>
<span id="L2484" class="ln">  2484&nbsp;&nbsp;</span>type parseSpecFunction func(doc *ast.CommentGroup, keyword token.Token, iota int) ast.Spec
<span id="L2485" class="ln">  2485&nbsp;&nbsp;</span>
<span id="L2486" class="ln">  2486&nbsp;&nbsp;</span>func (p *parser) parseImportSpec(doc *ast.CommentGroup, _ token.Token, _ int) ast.Spec {
<span id="L2487" class="ln">  2487&nbsp;&nbsp;</span>	if p.trace {
<span id="L2488" class="ln">  2488&nbsp;&nbsp;</span>		defer un(trace(p, &#34;ImportSpec&#34;))
<span id="L2489" class="ln">  2489&nbsp;&nbsp;</span>	}
<span id="L2490" class="ln">  2490&nbsp;&nbsp;</span>
<span id="L2491" class="ln">  2491&nbsp;&nbsp;</span>	var ident *ast.Ident
<span id="L2492" class="ln">  2492&nbsp;&nbsp;</span>	switch p.tok {
<span id="L2493" class="ln">  2493&nbsp;&nbsp;</span>	case token.IDENT:
<span id="L2494" class="ln">  2494&nbsp;&nbsp;</span>		ident = p.parseIdent()
<span id="L2495" class="ln">  2495&nbsp;&nbsp;</span>	case token.PERIOD:
<span id="L2496" class="ln">  2496&nbsp;&nbsp;</span>		ident = &amp;ast.Ident{NamePos: p.pos, Name: &#34;.&#34;}
<span id="L2497" class="ln">  2497&nbsp;&nbsp;</span>		p.next()
<span id="L2498" class="ln">  2498&nbsp;&nbsp;</span>	}
<span id="L2499" class="ln">  2499&nbsp;&nbsp;</span>
<span id="L2500" class="ln">  2500&nbsp;&nbsp;</span>	pos := p.pos
<span id="L2501" class="ln">  2501&nbsp;&nbsp;</span>	var path string
<span id="L2502" class="ln">  2502&nbsp;&nbsp;</span>	if p.tok == token.STRING {
<span id="L2503" class="ln">  2503&nbsp;&nbsp;</span>		path = p.lit
<span id="L2504" class="ln">  2504&nbsp;&nbsp;</span>		p.next()
<span id="L2505" class="ln">  2505&nbsp;&nbsp;</span>	} else if p.tok.IsLiteral() {
<span id="L2506" class="ln">  2506&nbsp;&nbsp;</span>		p.error(pos, &#34;import path must be a string&#34;)
<span id="L2507" class="ln">  2507&nbsp;&nbsp;</span>		p.next()
<span id="L2508" class="ln">  2508&nbsp;&nbsp;</span>	} else {
<span id="L2509" class="ln">  2509&nbsp;&nbsp;</span>		p.error(pos, &#34;missing import path&#34;)
<span id="L2510" class="ln">  2510&nbsp;&nbsp;</span>		p.advance(exprEnd)
<span id="L2511" class="ln">  2511&nbsp;&nbsp;</span>	}
<span id="L2512" class="ln">  2512&nbsp;&nbsp;</span>	comment := p.expectSemi()
<span id="L2513" class="ln">  2513&nbsp;&nbsp;</span>
<span id="L2514" class="ln">  2514&nbsp;&nbsp;</span>	<span class="comment">// collect imports</span>
<span id="L2515" class="ln">  2515&nbsp;&nbsp;</span>	spec := &amp;ast.ImportSpec{
<span id="L2516" class="ln">  2516&nbsp;&nbsp;</span>		Doc:     doc,
<span id="L2517" class="ln">  2517&nbsp;&nbsp;</span>		Name:    ident,
<span id="L2518" class="ln">  2518&nbsp;&nbsp;</span>		Path:    &amp;ast.BasicLit{ValuePos: pos, Kind: token.STRING, Value: path},
<span id="L2519" class="ln">  2519&nbsp;&nbsp;</span>		Comment: comment,
<span id="L2520" class="ln">  2520&nbsp;&nbsp;</span>	}
<span id="L2521" class="ln">  2521&nbsp;&nbsp;</span>	p.imports = append(p.imports, spec)
<span id="L2522" class="ln">  2522&nbsp;&nbsp;</span>
<span id="L2523" class="ln">  2523&nbsp;&nbsp;</span>	return spec
<span id="L2524" class="ln">  2524&nbsp;&nbsp;</span>}
<span id="L2525" class="ln">  2525&nbsp;&nbsp;</span>
<span id="L2526" class="ln">  2526&nbsp;&nbsp;</span>func (p *parser) parseValueSpec(doc *ast.CommentGroup, keyword token.Token, iota int) ast.Spec {
<span id="L2527" class="ln">  2527&nbsp;&nbsp;</span>	if p.trace {
<span id="L2528" class="ln">  2528&nbsp;&nbsp;</span>		defer un(trace(p, keyword.String()+&#34;Spec&#34;))
<span id="L2529" class="ln">  2529&nbsp;&nbsp;</span>	}
<span id="L2530" class="ln">  2530&nbsp;&nbsp;</span>
<span id="L2531" class="ln">  2531&nbsp;&nbsp;</span>	idents := p.parseIdentList()
<span id="L2532" class="ln">  2532&nbsp;&nbsp;</span>	var typ ast.Expr
<span id="L2533" class="ln">  2533&nbsp;&nbsp;</span>	var values []ast.Expr
<span id="L2534" class="ln">  2534&nbsp;&nbsp;</span>	switch keyword {
<span id="L2535" class="ln">  2535&nbsp;&nbsp;</span>	case token.CONST:
<span id="L2536" class="ln">  2536&nbsp;&nbsp;</span>		<span class="comment">// always permit optional type and initialization for more tolerant parsing</span>
<span id="L2537" class="ln">  2537&nbsp;&nbsp;</span>		if p.tok != token.EOF &amp;&amp; p.tok != token.SEMICOLON &amp;&amp; p.tok != token.RPAREN {
<span id="L2538" class="ln">  2538&nbsp;&nbsp;</span>			typ = p.tryIdentOrType()
<span id="L2539" class="ln">  2539&nbsp;&nbsp;</span>			if p.tok == token.ASSIGN {
<span id="L2540" class="ln">  2540&nbsp;&nbsp;</span>				p.next()
<span id="L2541" class="ln">  2541&nbsp;&nbsp;</span>				values = p.parseList(true)
<span id="L2542" class="ln">  2542&nbsp;&nbsp;</span>			}
<span id="L2543" class="ln">  2543&nbsp;&nbsp;</span>		}
<span id="L2544" class="ln">  2544&nbsp;&nbsp;</span>	case token.VAR:
<span id="L2545" class="ln">  2545&nbsp;&nbsp;</span>		if p.tok != token.ASSIGN {
<span id="L2546" class="ln">  2546&nbsp;&nbsp;</span>			typ = p.parseType()
<span id="L2547" class="ln">  2547&nbsp;&nbsp;</span>		}
<span id="L2548" class="ln">  2548&nbsp;&nbsp;</span>		if p.tok == token.ASSIGN {
<span id="L2549" class="ln">  2549&nbsp;&nbsp;</span>			p.next()
<span id="L2550" class="ln">  2550&nbsp;&nbsp;</span>			values = p.parseList(true)
<span id="L2551" class="ln">  2551&nbsp;&nbsp;</span>		}
<span id="L2552" class="ln">  2552&nbsp;&nbsp;</span>	default:
<span id="L2553" class="ln">  2553&nbsp;&nbsp;</span>		panic(&#34;unreachable&#34;)
<span id="L2554" class="ln">  2554&nbsp;&nbsp;</span>	}
<span id="L2555" class="ln">  2555&nbsp;&nbsp;</span>	comment := p.expectSemi()
<span id="L2556" class="ln">  2556&nbsp;&nbsp;</span>
<span id="L2557" class="ln">  2557&nbsp;&nbsp;</span>	spec := &amp;ast.ValueSpec{
<span id="L2558" class="ln">  2558&nbsp;&nbsp;</span>		Doc:     doc,
<span id="L2559" class="ln">  2559&nbsp;&nbsp;</span>		Names:   idents,
<span id="L2560" class="ln">  2560&nbsp;&nbsp;</span>		Type:    typ,
<span id="L2561" class="ln">  2561&nbsp;&nbsp;</span>		Values:  values,
<span id="L2562" class="ln">  2562&nbsp;&nbsp;</span>		Comment: comment,
<span id="L2563" class="ln">  2563&nbsp;&nbsp;</span>	}
<span id="L2564" class="ln">  2564&nbsp;&nbsp;</span>	return spec
<span id="L2565" class="ln">  2565&nbsp;&nbsp;</span>}
<span id="L2566" class="ln">  2566&nbsp;&nbsp;</span>
<span id="L2567" class="ln">  2567&nbsp;&nbsp;</span>func (p *parser) parseGenericType(spec *ast.TypeSpec, openPos token.Pos, name0 *ast.Ident, typ0 ast.Expr) {
<span id="L2568" class="ln">  2568&nbsp;&nbsp;</span>	if p.trace {
<span id="L2569" class="ln">  2569&nbsp;&nbsp;</span>		defer un(trace(p, &#34;parseGenericType&#34;))
<span id="L2570" class="ln">  2570&nbsp;&nbsp;</span>	}
<span id="L2571" class="ln">  2571&nbsp;&nbsp;</span>
<span id="L2572" class="ln">  2572&nbsp;&nbsp;</span>	list := p.parseParameterList(name0, typ0, token.RBRACK)
<span id="L2573" class="ln">  2573&nbsp;&nbsp;</span>	closePos := p.expect(token.RBRACK)
<span id="L2574" class="ln">  2574&nbsp;&nbsp;</span>	spec.TypeParams = &amp;ast.FieldList{Opening: openPos, List: list, Closing: closePos}
<span id="L2575" class="ln">  2575&nbsp;&nbsp;</span>	<span class="comment">// Let the type checker decide whether to accept type parameters on aliases:</span>
<span id="L2576" class="ln">  2576&nbsp;&nbsp;</span>	<span class="comment">// see go.dev/issue/46477.</span>
<span id="L2577" class="ln">  2577&nbsp;&nbsp;</span>	if p.tok == token.ASSIGN {
<span id="L2578" class="ln">  2578&nbsp;&nbsp;</span>		<span class="comment">// type alias</span>
<span id="L2579" class="ln">  2579&nbsp;&nbsp;</span>		spec.Assign = p.pos
<span id="L2580" class="ln">  2580&nbsp;&nbsp;</span>		p.next()
<span id="L2581" class="ln">  2581&nbsp;&nbsp;</span>	}
<span id="L2582" class="ln">  2582&nbsp;&nbsp;</span>	spec.Type = p.parseType()
<span id="L2583" class="ln">  2583&nbsp;&nbsp;</span>}
<span id="L2584" class="ln">  2584&nbsp;&nbsp;</span>
<span id="L2585" class="ln">  2585&nbsp;&nbsp;</span>func (p *parser) parseTypeSpec(doc *ast.CommentGroup, _ token.Token, _ int) ast.Spec {
<span id="L2586" class="ln">  2586&nbsp;&nbsp;</span>	if p.trace {
<span id="L2587" class="ln">  2587&nbsp;&nbsp;</span>		defer un(trace(p, &#34;TypeSpec&#34;))
<span id="L2588" class="ln">  2588&nbsp;&nbsp;</span>	}
<span id="L2589" class="ln">  2589&nbsp;&nbsp;</span>
<span id="L2590" class="ln">  2590&nbsp;&nbsp;</span>	name := p.parseIdent()
<span id="L2591" class="ln">  2591&nbsp;&nbsp;</span>	spec := &amp;ast.TypeSpec{Doc: doc, Name: name}
<span id="L2592" class="ln">  2592&nbsp;&nbsp;</span>
<span id="L2593" class="ln">  2593&nbsp;&nbsp;</span>	if p.tok == token.LBRACK {
<span id="L2594" class="ln">  2594&nbsp;&nbsp;</span>		<span class="comment">// spec.Name &#34;[&#34; ...</span>
<span id="L2595" class="ln">  2595&nbsp;&nbsp;</span>		<span class="comment">// array/slice type or type parameter list</span>
<span id="L2596" class="ln">  2596&nbsp;&nbsp;</span>		lbrack := p.pos
<span id="L2597" class="ln">  2597&nbsp;&nbsp;</span>		p.next()
<span id="L2598" class="ln">  2598&nbsp;&nbsp;</span>		if p.tok == token.IDENT {
<span id="L2599" class="ln">  2599&nbsp;&nbsp;</span>			<span class="comment">// We may have an array type or a type parameter list.</span>
<span id="L2600" class="ln">  2600&nbsp;&nbsp;</span>			<span class="comment">// In either case we expect an expression x (which may</span>
<span id="L2601" class="ln">  2601&nbsp;&nbsp;</span>			<span class="comment">// just be a name, or a more complex expression) which</span>
<span id="L2602" class="ln">  2602&nbsp;&nbsp;</span>			<span class="comment">// we can analyze further.</span>
<span id="L2603" class="ln">  2603&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L2604" class="ln">  2604&nbsp;&nbsp;</span>			<span class="comment">// A type parameter list may have a type bound starting</span>
<span id="L2605" class="ln">  2605&nbsp;&nbsp;</span>			<span class="comment">// with a &#34;[&#34; as in: P []E. In that case, simply parsing</span>
<span id="L2606" class="ln">  2606&nbsp;&nbsp;</span>			<span class="comment">// an expression would lead to an error: P[] is invalid.</span>
<span id="L2607" class="ln">  2607&nbsp;&nbsp;</span>			<span class="comment">// But since index or slice expressions are never constant</span>
<span id="L2608" class="ln">  2608&nbsp;&nbsp;</span>			<span class="comment">// and thus invalid array length expressions, if the name</span>
<span id="L2609" class="ln">  2609&nbsp;&nbsp;</span>			<span class="comment">// is followed by &#34;[&#34; it must be the start of an array or</span>
<span id="L2610" class="ln">  2610&nbsp;&nbsp;</span>			<span class="comment">// slice constraint. Only if we don&#39;t see a &#34;[&#34; do we</span>
<span id="L2611" class="ln">  2611&nbsp;&nbsp;</span>			<span class="comment">// need to parse a full expression. Notably, name &lt;- x</span>
<span id="L2612" class="ln">  2612&nbsp;&nbsp;</span>			<span class="comment">// is not a concern because name &lt;- x is a statement and</span>
<span id="L2613" class="ln">  2613&nbsp;&nbsp;</span>			<span class="comment">// not an expression.</span>
<span id="L2614" class="ln">  2614&nbsp;&nbsp;</span>			var x ast.Expr = p.parseIdent()
<span id="L2615" class="ln">  2615&nbsp;&nbsp;</span>			if p.tok != token.LBRACK {
<span id="L2616" class="ln">  2616&nbsp;&nbsp;</span>				<span class="comment">// To parse the expression starting with name, expand</span>
<span id="L2617" class="ln">  2617&nbsp;&nbsp;</span>				<span class="comment">// the call sequence we would get by passing in name</span>
<span id="L2618" class="ln">  2618&nbsp;&nbsp;</span>				<span class="comment">// to parser.expr, and pass in name to parsePrimaryExpr.</span>
<span id="L2619" class="ln">  2619&nbsp;&nbsp;</span>				p.exprLev++
<span id="L2620" class="ln">  2620&nbsp;&nbsp;</span>				lhs := p.parsePrimaryExpr(x)
<span id="L2621" class="ln">  2621&nbsp;&nbsp;</span>				x = p.parseBinaryExpr(lhs, token.LowestPrec+1)
<span id="L2622" class="ln">  2622&nbsp;&nbsp;</span>				p.exprLev--
<span id="L2623" class="ln">  2623&nbsp;&nbsp;</span>			}
<span id="L2624" class="ln">  2624&nbsp;&nbsp;</span>			<span class="comment">// Analyze expression x. If we can split x into a type parameter</span>
<span id="L2625" class="ln">  2625&nbsp;&nbsp;</span>			<span class="comment">// name, possibly followed by a type parameter type, we consider</span>
<span id="L2626" class="ln">  2626&nbsp;&nbsp;</span>			<span class="comment">// this the start of a type parameter list, with some caveats:</span>
<span id="L2627" class="ln">  2627&nbsp;&nbsp;</span>			<span class="comment">// a single name followed by &#34;]&#34; tilts the decision towards an</span>
<span id="L2628" class="ln">  2628&nbsp;&nbsp;</span>			<span class="comment">// array declaration; a type parameter type that could also be</span>
<span id="L2629" class="ln">  2629&nbsp;&nbsp;</span>			<span class="comment">// an ordinary expression but which is followed by a comma tilts</span>
<span id="L2630" class="ln">  2630&nbsp;&nbsp;</span>			<span class="comment">// the decision towards a type parameter list.</span>
<span id="L2631" class="ln">  2631&nbsp;&nbsp;</span>			if pname, ptype := extractName(x, p.tok == token.COMMA); pname != nil &amp;&amp; (ptype != nil || p.tok != token.RBRACK) {
<span id="L2632" class="ln">  2632&nbsp;&nbsp;</span>				<span class="comment">// spec.Name &#34;[&#34; pname ...</span>
<span id="L2633" class="ln">  2633&nbsp;&nbsp;</span>				<span class="comment">// spec.Name &#34;[&#34; pname ptype ...</span>
<span id="L2634" class="ln">  2634&nbsp;&nbsp;</span>				<span class="comment">// spec.Name &#34;[&#34; pname ptype &#34;,&#34; ...</span>
<span id="L2635" class="ln">  2635&nbsp;&nbsp;</span>				p.parseGenericType(spec, lbrack, pname, ptype) <span class="comment">// ptype may be nil</span>
<span id="L2636" class="ln">  2636&nbsp;&nbsp;</span>			} else {
<span id="L2637" class="ln">  2637&nbsp;&nbsp;</span>				<span class="comment">// spec.Name &#34;[&#34; pname &#34;]&#34; ...</span>
<span id="L2638" class="ln">  2638&nbsp;&nbsp;</span>				<span class="comment">// spec.Name &#34;[&#34; x ...</span>
<span id="L2639" class="ln">  2639&nbsp;&nbsp;</span>				spec.Type = p.parseArrayType(lbrack, x)
<span id="L2640" class="ln">  2640&nbsp;&nbsp;</span>			}
<span id="L2641" class="ln">  2641&nbsp;&nbsp;</span>		} else {
<span id="L2642" class="ln">  2642&nbsp;&nbsp;</span>			<span class="comment">// array type</span>
<span id="L2643" class="ln">  2643&nbsp;&nbsp;</span>			spec.Type = p.parseArrayType(lbrack, nil)
<span id="L2644" class="ln">  2644&nbsp;&nbsp;</span>		}
<span id="L2645" class="ln">  2645&nbsp;&nbsp;</span>	} else {
<span id="L2646" class="ln">  2646&nbsp;&nbsp;</span>		<span class="comment">// no type parameters</span>
<span id="L2647" class="ln">  2647&nbsp;&nbsp;</span>		if p.tok == token.ASSIGN {
<span id="L2648" class="ln">  2648&nbsp;&nbsp;</span>			<span class="comment">// type alias</span>
<span id="L2649" class="ln">  2649&nbsp;&nbsp;</span>			spec.Assign = p.pos
<span id="L2650" class="ln">  2650&nbsp;&nbsp;</span>			p.next()
<span id="L2651" class="ln">  2651&nbsp;&nbsp;</span>		}
<span id="L2652" class="ln">  2652&nbsp;&nbsp;</span>		spec.Type = p.parseType()
<span id="L2653" class="ln">  2653&nbsp;&nbsp;</span>	}
<span id="L2654" class="ln">  2654&nbsp;&nbsp;</span>
<span id="L2655" class="ln">  2655&nbsp;&nbsp;</span>	spec.Comment = p.expectSemi()
<span id="L2656" class="ln">  2656&nbsp;&nbsp;</span>
<span id="L2657" class="ln">  2657&nbsp;&nbsp;</span>	return spec
<span id="L2658" class="ln">  2658&nbsp;&nbsp;</span>}
<span id="L2659" class="ln">  2659&nbsp;&nbsp;</span>
<span id="L2660" class="ln">  2660&nbsp;&nbsp;</span><span class="comment">// extractName splits the expression x into (name, expr) if syntactically</span>
<span id="L2661" class="ln">  2661&nbsp;&nbsp;</span><span class="comment">// x can be written as name expr. The split only happens if expr is a type</span>
<span id="L2662" class="ln">  2662&nbsp;&nbsp;</span><span class="comment">// element (per the isTypeElem predicate) or if force is set.</span>
<span id="L2663" class="ln">  2663&nbsp;&nbsp;</span><span class="comment">// If x is just a name, the result is (name, nil). If the split succeeds,</span>
<span id="L2664" class="ln">  2664&nbsp;&nbsp;</span><span class="comment">// the result is (name, expr). Otherwise the result is (nil, x).</span>
<span id="L2665" class="ln">  2665&nbsp;&nbsp;</span><span class="comment">// Examples:</span>
<span id="L2666" class="ln">  2666&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L2667" class="ln">  2667&nbsp;&nbsp;</span><span class="comment">//	x           force    name    expr</span>
<span id="L2668" class="ln">  2668&nbsp;&nbsp;</span><span class="comment">//	------------------------------------</span>
<span id="L2669" class="ln">  2669&nbsp;&nbsp;</span><span class="comment">//	P*[]int     T/F      P       *[]int</span>
<span id="L2670" class="ln">  2670&nbsp;&nbsp;</span><span class="comment">//	P*E         T        P       *E</span>
<span id="L2671" class="ln">  2671&nbsp;&nbsp;</span><span class="comment">//	P*E         F        nil     P*E</span>
<span id="L2672" class="ln">  2672&nbsp;&nbsp;</span><span class="comment">//	P([]int)    T/F      P       []int</span>
<span id="L2673" class="ln">  2673&nbsp;&nbsp;</span><span class="comment">//	P(E)        T        P       E</span>
<span id="L2674" class="ln">  2674&nbsp;&nbsp;</span><span class="comment">//	P(E)        F        nil     P(E)</span>
<span id="L2675" class="ln">  2675&nbsp;&nbsp;</span><span class="comment">//	P*E|F|~G    T/F      P       *E|F|~G</span>
<span id="L2676" class="ln">  2676&nbsp;&nbsp;</span><span class="comment">//	P*E|F|G     T        P       *E|F|G</span>
<span id="L2677" class="ln">  2677&nbsp;&nbsp;</span><span class="comment">//	P*E|F|G     F        nil     P*E|F|G</span>
<span id="L2678" class="ln">  2678&nbsp;&nbsp;</span>func extractName(x ast.Expr, force bool) (*ast.Ident, ast.Expr) {
<span id="L2679" class="ln">  2679&nbsp;&nbsp;</span>	switch x := x.(type) {
<span id="L2680" class="ln">  2680&nbsp;&nbsp;</span>	case *ast.Ident:
<span id="L2681" class="ln">  2681&nbsp;&nbsp;</span>		return x, nil
<span id="L2682" class="ln">  2682&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L2683" class="ln">  2683&nbsp;&nbsp;</span>		switch x.Op {
<span id="L2684" class="ln">  2684&nbsp;&nbsp;</span>		case token.MUL:
<span id="L2685" class="ln">  2685&nbsp;&nbsp;</span>			if name, _ := x.X.(*ast.Ident); name != nil &amp;&amp; (force || isTypeElem(x.Y)) {
<span id="L2686" class="ln">  2686&nbsp;&nbsp;</span>				<span class="comment">// x = name *x.Y</span>
<span id="L2687" class="ln">  2687&nbsp;&nbsp;</span>				return name, &amp;ast.StarExpr{Star: x.OpPos, X: x.Y}
<span id="L2688" class="ln">  2688&nbsp;&nbsp;</span>			}
<span id="L2689" class="ln">  2689&nbsp;&nbsp;</span>		case token.OR:
<span id="L2690" class="ln">  2690&nbsp;&nbsp;</span>			if name, lhs := extractName(x.X, force || isTypeElem(x.Y)); name != nil &amp;&amp; lhs != nil {
<span id="L2691" class="ln">  2691&nbsp;&nbsp;</span>				<span class="comment">// x = name lhs|x.Y</span>
<span id="L2692" class="ln">  2692&nbsp;&nbsp;</span>				op := *x
<span id="L2693" class="ln">  2693&nbsp;&nbsp;</span>				op.X = lhs
<span id="L2694" class="ln">  2694&nbsp;&nbsp;</span>				return name, &amp;op
<span id="L2695" class="ln">  2695&nbsp;&nbsp;</span>			}
<span id="L2696" class="ln">  2696&nbsp;&nbsp;</span>		}
<span id="L2697" class="ln">  2697&nbsp;&nbsp;</span>	case *ast.CallExpr:
<span id="L2698" class="ln">  2698&nbsp;&nbsp;</span>		if name, _ := x.Fun.(*ast.Ident); name != nil {
<span id="L2699" class="ln">  2699&nbsp;&nbsp;</span>			if len(x.Args) == 1 &amp;&amp; x.Ellipsis == token.NoPos &amp;&amp; (force || isTypeElem(x.Args[0])) {
<span id="L2700" class="ln">  2700&nbsp;&nbsp;</span>				<span class="comment">// x = name &#34;(&#34; x.ArgList[0] &#34;)&#34;</span>
<span id="L2701" class="ln">  2701&nbsp;&nbsp;</span>				return name, x.Args[0]
<span id="L2702" class="ln">  2702&nbsp;&nbsp;</span>			}
<span id="L2703" class="ln">  2703&nbsp;&nbsp;</span>		}
<span id="L2704" class="ln">  2704&nbsp;&nbsp;</span>	}
<span id="L2705" class="ln">  2705&nbsp;&nbsp;</span>	return nil, x
<span id="L2706" class="ln">  2706&nbsp;&nbsp;</span>}
<span id="L2707" class="ln">  2707&nbsp;&nbsp;</span>
<span id="L2708" class="ln">  2708&nbsp;&nbsp;</span><span class="comment">// isTypeElem reports whether x is a (possibly parenthesized) type element expression.</span>
<span id="L2709" class="ln">  2709&nbsp;&nbsp;</span><span class="comment">// The result is false if x could be a type element OR an ordinary (value) expression.</span>
<span id="L2710" class="ln">  2710&nbsp;&nbsp;</span>func isTypeElem(x ast.Expr) bool {
<span id="L2711" class="ln">  2711&nbsp;&nbsp;</span>	switch x := x.(type) {
<span id="L2712" class="ln">  2712&nbsp;&nbsp;</span>	case *ast.ArrayType, *ast.StructType, *ast.FuncType, *ast.InterfaceType, *ast.MapType, *ast.ChanType:
<span id="L2713" class="ln">  2713&nbsp;&nbsp;</span>		return true
<span id="L2714" class="ln">  2714&nbsp;&nbsp;</span>	case *ast.BinaryExpr:
<span id="L2715" class="ln">  2715&nbsp;&nbsp;</span>		return isTypeElem(x.X) || isTypeElem(x.Y)
<span id="L2716" class="ln">  2716&nbsp;&nbsp;</span>	case *ast.UnaryExpr:
<span id="L2717" class="ln">  2717&nbsp;&nbsp;</span>		return x.Op == token.TILDE
<span id="L2718" class="ln">  2718&nbsp;&nbsp;</span>	case *ast.ParenExpr:
<span id="L2719" class="ln">  2719&nbsp;&nbsp;</span>		return isTypeElem(x.X)
<span id="L2720" class="ln">  2720&nbsp;&nbsp;</span>	}
<span id="L2721" class="ln">  2721&nbsp;&nbsp;</span>	return false
<span id="L2722" class="ln">  2722&nbsp;&nbsp;</span>}
<span id="L2723" class="ln">  2723&nbsp;&nbsp;</span>
<span id="L2724" class="ln">  2724&nbsp;&nbsp;</span>func (p *parser) parseGenDecl(keyword token.Token, f parseSpecFunction) *ast.GenDecl {
<span id="L2725" class="ln">  2725&nbsp;&nbsp;</span>	if p.trace {
<span id="L2726" class="ln">  2726&nbsp;&nbsp;</span>		defer un(trace(p, &#34;GenDecl(&#34;+keyword.String()+&#34;)&#34;))
<span id="L2727" class="ln">  2727&nbsp;&nbsp;</span>	}
<span id="L2728" class="ln">  2728&nbsp;&nbsp;</span>
<span id="L2729" class="ln">  2729&nbsp;&nbsp;</span>	doc := p.leadComment
<span id="L2730" class="ln">  2730&nbsp;&nbsp;</span>	pos := p.expect(keyword)
<span id="L2731" class="ln">  2731&nbsp;&nbsp;</span>	var lparen, rparen token.Pos
<span id="L2732" class="ln">  2732&nbsp;&nbsp;</span>	var list []ast.Spec
<span id="L2733" class="ln">  2733&nbsp;&nbsp;</span>	if p.tok == token.LPAREN {
<span id="L2734" class="ln">  2734&nbsp;&nbsp;</span>		lparen = p.pos
<span id="L2735" class="ln">  2735&nbsp;&nbsp;</span>		p.next()
<span id="L2736" class="ln">  2736&nbsp;&nbsp;</span>		for iota := 0; p.tok != token.RPAREN &amp;&amp; p.tok != token.EOF; iota++ {
<span id="L2737" class="ln">  2737&nbsp;&nbsp;</span>			list = append(list, f(p.leadComment, keyword, iota))
<span id="L2738" class="ln">  2738&nbsp;&nbsp;</span>		}
<span id="L2739" class="ln">  2739&nbsp;&nbsp;</span>		rparen = p.expect(token.RPAREN)
<span id="L2740" class="ln">  2740&nbsp;&nbsp;</span>		p.expectSemi()
<span id="L2741" class="ln">  2741&nbsp;&nbsp;</span>	} else {
<span id="L2742" class="ln">  2742&nbsp;&nbsp;</span>		list = append(list, f(nil, keyword, 0))
<span id="L2743" class="ln">  2743&nbsp;&nbsp;</span>	}
<span id="L2744" class="ln">  2744&nbsp;&nbsp;</span>
<span id="L2745" class="ln">  2745&nbsp;&nbsp;</span>	return &amp;ast.GenDecl{
<span id="L2746" class="ln">  2746&nbsp;&nbsp;</span>		Doc:    doc,
<span id="L2747" class="ln">  2747&nbsp;&nbsp;</span>		TokPos: pos,
<span id="L2748" class="ln">  2748&nbsp;&nbsp;</span>		Tok:    keyword,
<span id="L2749" class="ln">  2749&nbsp;&nbsp;</span>		Lparen: lparen,
<span id="L2750" class="ln">  2750&nbsp;&nbsp;</span>		Specs:  list,
<span id="L2751" class="ln">  2751&nbsp;&nbsp;</span>		Rparen: rparen,
<span id="L2752" class="ln">  2752&nbsp;&nbsp;</span>	}
<span id="L2753" class="ln">  2753&nbsp;&nbsp;</span>}
<span id="L2754" class="ln">  2754&nbsp;&nbsp;</span>
<span id="L2755" class="ln">  2755&nbsp;&nbsp;</span>func (p *parser) parseFuncDecl() *ast.FuncDecl {
<span id="L2756" class="ln">  2756&nbsp;&nbsp;</span>	if p.trace {
<span id="L2757" class="ln">  2757&nbsp;&nbsp;</span>		defer un(trace(p, &#34;FunctionDecl&#34;))
<span id="L2758" class="ln">  2758&nbsp;&nbsp;</span>	}
<span id="L2759" class="ln">  2759&nbsp;&nbsp;</span>
<span id="L2760" class="ln">  2760&nbsp;&nbsp;</span>	doc := p.leadComment
<span id="L2761" class="ln">  2761&nbsp;&nbsp;</span>	pos := p.expect(token.FUNC)
<span id="L2762" class="ln">  2762&nbsp;&nbsp;</span>
<span id="L2763" class="ln">  2763&nbsp;&nbsp;</span>	var recv *ast.FieldList
<span id="L2764" class="ln">  2764&nbsp;&nbsp;</span>	if p.tok == token.LPAREN {
<span id="L2765" class="ln">  2765&nbsp;&nbsp;</span>		_, recv = p.parseParameters(false)
<span id="L2766" class="ln">  2766&nbsp;&nbsp;</span>	}
<span id="L2767" class="ln">  2767&nbsp;&nbsp;</span>
<span id="L2768" class="ln">  2768&nbsp;&nbsp;</span>	ident := p.parseIdent()
<span id="L2769" class="ln">  2769&nbsp;&nbsp;</span>
<span id="L2770" class="ln">  2770&nbsp;&nbsp;</span>	tparams, params := p.parseParameters(true)
<span id="L2771" class="ln">  2771&nbsp;&nbsp;</span>	if recv != nil &amp;&amp; tparams != nil {
<span id="L2772" class="ln">  2772&nbsp;&nbsp;</span>		<span class="comment">// Method declarations do not have type parameters. We parse them for a</span>
<span id="L2773" class="ln">  2773&nbsp;&nbsp;</span>		<span class="comment">// better error message and improved error recovery.</span>
<span id="L2774" class="ln">  2774&nbsp;&nbsp;</span>		p.error(tparams.Opening, &#34;method must have no type parameters&#34;)
<span id="L2775" class="ln">  2775&nbsp;&nbsp;</span>		tparams = nil
<span id="L2776" class="ln">  2776&nbsp;&nbsp;</span>	}
<span id="L2777" class="ln">  2777&nbsp;&nbsp;</span>	results := p.parseResult()
<span id="L2778" class="ln">  2778&nbsp;&nbsp;</span>
<span id="L2779" class="ln">  2779&nbsp;&nbsp;</span>	var body *ast.BlockStmt
<span id="L2780" class="ln">  2780&nbsp;&nbsp;</span>	switch p.tok {
<span id="L2781" class="ln">  2781&nbsp;&nbsp;</span>	case token.LBRACE:
<span id="L2782" class="ln">  2782&nbsp;&nbsp;</span>		body = p.parseBody()
<span id="L2783" class="ln">  2783&nbsp;&nbsp;</span>		p.expectSemi()
<span id="L2784" class="ln">  2784&nbsp;&nbsp;</span>	case token.SEMICOLON:
<span id="L2785" class="ln">  2785&nbsp;&nbsp;</span>		p.next()
<span id="L2786" class="ln">  2786&nbsp;&nbsp;</span>		if p.tok == token.LBRACE {
<span id="L2787" class="ln">  2787&nbsp;&nbsp;</span>			<span class="comment">// opening { of function declaration on next line</span>
<span id="L2788" class="ln">  2788&nbsp;&nbsp;</span>			p.error(p.pos, &#34;unexpected semicolon or newline before {&#34;)
<span id="L2789" class="ln">  2789&nbsp;&nbsp;</span>			body = p.parseBody()
<span id="L2790" class="ln">  2790&nbsp;&nbsp;</span>			p.expectSemi()
<span id="L2791" class="ln">  2791&nbsp;&nbsp;</span>		}
<span id="L2792" class="ln">  2792&nbsp;&nbsp;</span>	default:
<span id="L2793" class="ln">  2793&nbsp;&nbsp;</span>		p.expectSemi()
<span id="L2794" class="ln">  2794&nbsp;&nbsp;</span>	}
<span id="L2795" class="ln">  2795&nbsp;&nbsp;</span>
<span id="L2796" class="ln">  2796&nbsp;&nbsp;</span>	decl := &amp;ast.FuncDecl{
<span id="L2797" class="ln">  2797&nbsp;&nbsp;</span>		Doc:  doc,
<span id="L2798" class="ln">  2798&nbsp;&nbsp;</span>		Recv: recv,
<span id="L2799" class="ln">  2799&nbsp;&nbsp;</span>		Name: ident,
<span id="L2800" class="ln">  2800&nbsp;&nbsp;</span>		Type: &amp;ast.FuncType{
<span id="L2801" class="ln">  2801&nbsp;&nbsp;</span>			Func:       pos,
<span id="L2802" class="ln">  2802&nbsp;&nbsp;</span>			TypeParams: tparams,
<span id="L2803" class="ln">  2803&nbsp;&nbsp;</span>			Params:     params,
<span id="L2804" class="ln">  2804&nbsp;&nbsp;</span>			Results:    results,
<span id="L2805" class="ln">  2805&nbsp;&nbsp;</span>		},
<span id="L2806" class="ln">  2806&nbsp;&nbsp;</span>		Body: body,
<span id="L2807" class="ln">  2807&nbsp;&nbsp;</span>	}
<span id="L2808" class="ln">  2808&nbsp;&nbsp;</span>	return decl
<span id="L2809" class="ln">  2809&nbsp;&nbsp;</span>}
<span id="L2810" class="ln">  2810&nbsp;&nbsp;</span>
<span id="L2811" class="ln">  2811&nbsp;&nbsp;</span>func (p *parser) parseDecl(sync map[token.Token]bool) ast.Decl {
<span id="L2812" class="ln">  2812&nbsp;&nbsp;</span>	if p.trace {
<span id="L2813" class="ln">  2813&nbsp;&nbsp;</span>		defer un(trace(p, &#34;Declaration&#34;))
<span id="L2814" class="ln">  2814&nbsp;&nbsp;</span>	}
<span id="L2815" class="ln">  2815&nbsp;&nbsp;</span>
<span id="L2816" class="ln">  2816&nbsp;&nbsp;</span>	var f parseSpecFunction
<span id="L2817" class="ln">  2817&nbsp;&nbsp;</span>	switch p.tok {
<span id="L2818" class="ln">  2818&nbsp;&nbsp;</span>	case token.IMPORT:
<span id="L2819" class="ln">  2819&nbsp;&nbsp;</span>		f = p.parseImportSpec
<span id="L2820" class="ln">  2820&nbsp;&nbsp;</span>
<span id="L2821" class="ln">  2821&nbsp;&nbsp;</span>	case token.CONST, token.VAR:
<span id="L2822" class="ln">  2822&nbsp;&nbsp;</span>		f = p.parseValueSpec
<span id="L2823" class="ln">  2823&nbsp;&nbsp;</span>
<span id="L2824" class="ln">  2824&nbsp;&nbsp;</span>	case token.TYPE:
<span id="L2825" class="ln">  2825&nbsp;&nbsp;</span>		f = p.parseTypeSpec
<span id="L2826" class="ln">  2826&nbsp;&nbsp;</span>
<span id="L2827" class="ln">  2827&nbsp;&nbsp;</span>	case token.FUNC:
<span id="L2828" class="ln">  2828&nbsp;&nbsp;</span>		return p.parseFuncDecl()
<span id="L2829" class="ln">  2829&nbsp;&nbsp;</span>
<span id="L2830" class="ln">  2830&nbsp;&nbsp;</span>	default:
<span id="L2831" class="ln">  2831&nbsp;&nbsp;</span>		pos := p.pos
<span id="L2832" class="ln">  2832&nbsp;&nbsp;</span>		p.errorExpected(pos, &#34;declaration&#34;)
<span id="L2833" class="ln">  2833&nbsp;&nbsp;</span>		p.advance(sync)
<span id="L2834" class="ln">  2834&nbsp;&nbsp;</span>		return &amp;ast.BadDecl{From: pos, To: p.pos}
<span id="L2835" class="ln">  2835&nbsp;&nbsp;</span>	}
<span id="L2836" class="ln">  2836&nbsp;&nbsp;</span>
<span id="L2837" class="ln">  2837&nbsp;&nbsp;</span>	return p.parseGenDecl(p.tok, f)
<span id="L2838" class="ln">  2838&nbsp;&nbsp;</span>}
<span id="L2839" class="ln">  2839&nbsp;&nbsp;</span>
<span id="L2840" class="ln">  2840&nbsp;&nbsp;</span><span class="comment">// ----------------------------------------------------------------------------</span>
<span id="L2841" class="ln">  2841&nbsp;&nbsp;</span><span class="comment">// Source files</span>
<span id="L2842" class="ln">  2842&nbsp;&nbsp;</span>
<span id="L2843" class="ln">  2843&nbsp;&nbsp;</span>func (p *parser) parseFile() *ast.File {
<span id="L2844" class="ln">  2844&nbsp;&nbsp;</span>	if p.trace {
<span id="L2845" class="ln">  2845&nbsp;&nbsp;</span>		defer un(trace(p, &#34;File&#34;))
<span id="L2846" class="ln">  2846&nbsp;&nbsp;</span>	}
<span id="L2847" class="ln">  2847&nbsp;&nbsp;</span>
<span id="L2848" class="ln">  2848&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t bother parsing the rest if we had errors scanning the first token.</span>
<span id="L2849" class="ln">  2849&nbsp;&nbsp;</span>	<span class="comment">// Likely not a Go source file at all.</span>
<span id="L2850" class="ln">  2850&nbsp;&nbsp;</span>	if p.errors.Len() != 0 {
<span id="L2851" class="ln">  2851&nbsp;&nbsp;</span>		return nil
<span id="L2852" class="ln">  2852&nbsp;&nbsp;</span>	}
<span id="L2853" class="ln">  2853&nbsp;&nbsp;</span>
<span id="L2854" class="ln">  2854&nbsp;&nbsp;</span>	<span class="comment">// package clause</span>
<span id="L2855" class="ln">  2855&nbsp;&nbsp;</span>	doc := p.leadComment
<span id="L2856" class="ln">  2856&nbsp;&nbsp;</span>	pos := p.expect(token.PACKAGE)
<span id="L2857" class="ln">  2857&nbsp;&nbsp;</span>	<span class="comment">// Go spec: The package clause is not a declaration;</span>
<span id="L2858" class="ln">  2858&nbsp;&nbsp;</span>	<span class="comment">// the package name does not appear in any scope.</span>
<span id="L2859" class="ln">  2859&nbsp;&nbsp;</span>	ident := p.parseIdent()
<span id="L2860" class="ln">  2860&nbsp;&nbsp;</span>	if ident.Name == &#34;_&#34; &amp;&amp; p.mode&amp;DeclarationErrors != 0 {
<span id="L2861" class="ln">  2861&nbsp;&nbsp;</span>		p.error(p.pos, &#34;invalid package name _&#34;)
<span id="L2862" class="ln">  2862&nbsp;&nbsp;</span>	}
<span id="L2863" class="ln">  2863&nbsp;&nbsp;</span>	p.expectSemi()
<span id="L2864" class="ln">  2864&nbsp;&nbsp;</span>
<span id="L2865" class="ln">  2865&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t bother parsing the rest if we had errors parsing the package clause.</span>
<span id="L2866" class="ln">  2866&nbsp;&nbsp;</span>	<span class="comment">// Likely not a Go source file at all.</span>
<span id="L2867" class="ln">  2867&nbsp;&nbsp;</span>	if p.errors.Len() != 0 {
<span id="L2868" class="ln">  2868&nbsp;&nbsp;</span>		return nil
<span id="L2869" class="ln">  2869&nbsp;&nbsp;</span>	}
<span id="L2870" class="ln">  2870&nbsp;&nbsp;</span>
<span id="L2871" class="ln">  2871&nbsp;&nbsp;</span>	var decls []ast.Decl
<span id="L2872" class="ln">  2872&nbsp;&nbsp;</span>	if p.mode&amp;PackageClauseOnly == 0 {
<span id="L2873" class="ln">  2873&nbsp;&nbsp;</span>		<span class="comment">// import decls</span>
<span id="L2874" class="ln">  2874&nbsp;&nbsp;</span>		for p.tok == token.IMPORT {
<span id="L2875" class="ln">  2875&nbsp;&nbsp;</span>			decls = append(decls, p.parseGenDecl(token.IMPORT, p.parseImportSpec))
<span id="L2876" class="ln">  2876&nbsp;&nbsp;</span>		}
<span id="L2877" class="ln">  2877&nbsp;&nbsp;</span>
<span id="L2878" class="ln">  2878&nbsp;&nbsp;</span>		if p.mode&amp;ImportsOnly == 0 {
<span id="L2879" class="ln">  2879&nbsp;&nbsp;</span>			<span class="comment">// rest of package body</span>
<span id="L2880" class="ln">  2880&nbsp;&nbsp;</span>			prev := token.IMPORT
<span id="L2881" class="ln">  2881&nbsp;&nbsp;</span>			for p.tok != token.EOF {
<span id="L2882" class="ln">  2882&nbsp;&nbsp;</span>				<span class="comment">// Continue to accept import declarations for error tolerance, but complain.</span>
<span id="L2883" class="ln">  2883&nbsp;&nbsp;</span>				if p.tok == token.IMPORT &amp;&amp; prev != token.IMPORT {
<span id="L2884" class="ln">  2884&nbsp;&nbsp;</span>					p.error(p.pos, &#34;imports must appear before other declarations&#34;)
<span id="L2885" class="ln">  2885&nbsp;&nbsp;</span>				}
<span id="L2886" class="ln">  2886&nbsp;&nbsp;</span>				prev = p.tok
<span id="L2887" class="ln">  2887&nbsp;&nbsp;</span>
<span id="L2888" class="ln">  2888&nbsp;&nbsp;</span>				decls = append(decls, p.parseDecl(declStart))
<span id="L2889" class="ln">  2889&nbsp;&nbsp;</span>			}
<span id="L2890" class="ln">  2890&nbsp;&nbsp;</span>		}
<span id="L2891" class="ln">  2891&nbsp;&nbsp;</span>	}
<span id="L2892" class="ln">  2892&nbsp;&nbsp;</span>
<span id="L2893" class="ln">  2893&nbsp;&nbsp;</span>	f := &amp;ast.File{
<span id="L2894" class="ln">  2894&nbsp;&nbsp;</span>		Doc:       doc,
<span id="L2895" class="ln">  2895&nbsp;&nbsp;</span>		Package:   pos,
<span id="L2896" class="ln">  2896&nbsp;&nbsp;</span>		Name:      ident,
<span id="L2897" class="ln">  2897&nbsp;&nbsp;</span>		Decls:     decls,
<span id="L2898" class="ln">  2898&nbsp;&nbsp;</span>		FileStart: token.Pos(p.file.Base()),
<span id="L2899" class="ln">  2899&nbsp;&nbsp;</span>		FileEnd:   token.Pos(p.file.Base() + p.file.Size()),
<span id="L2900" class="ln">  2900&nbsp;&nbsp;</span>		Imports:   p.imports,
<span id="L2901" class="ln">  2901&nbsp;&nbsp;</span>		Comments:  p.comments,
<span id="L2902" class="ln">  2902&nbsp;&nbsp;</span>		GoVersion: p.goVersion,
<span id="L2903" class="ln">  2903&nbsp;&nbsp;</span>	}
<span id="L2904" class="ln">  2904&nbsp;&nbsp;</span>	var declErr func(token.Pos, string)
<span id="L2905" class="ln">  2905&nbsp;&nbsp;</span>	if p.mode&amp;DeclarationErrors != 0 {
<span id="L2906" class="ln">  2906&nbsp;&nbsp;</span>		declErr = p.error
<span id="L2907" class="ln">  2907&nbsp;&nbsp;</span>	}
<span id="L2908" class="ln">  2908&nbsp;&nbsp;</span>	if p.mode&amp;SkipObjectResolution == 0 {
<span id="L2909" class="ln">  2909&nbsp;&nbsp;</span>		resolveFile(f, p.file, declErr)
<span id="L2910" class="ln">  2910&nbsp;&nbsp;</span>	}
<span id="L2911" class="ln">  2911&nbsp;&nbsp;</span>
<span id="L2912" class="ln">  2912&nbsp;&nbsp;</span>	return f
<span id="L2913" class="ln">  2913&nbsp;&nbsp;</span>}
<span id="L2914" class="ln">  2914&nbsp;&nbsp;</span>
</pre><p><a href="parser.go?m=text">View as plain text</a></p>

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
