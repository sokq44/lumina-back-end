<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/text/template/parse/lex.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../../index.html">GoDoc</a></div>
<a href="lex.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/text">text</a>/<a href="http://localhost:8080/src/text/template">template</a>/<a href="http://localhost:8080/src/text/template/parse">parse</a>/<span class="text-muted">lex.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/text/template/parse">text/template/parse</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package parse
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// item represents a token or text string returned from the scanner.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>type item struct {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	typ  itemType <span class="comment">// The type of this item.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	pos  Pos      <span class="comment">// The starting position, in bytes, of this item in the input string.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	val  string   <span class="comment">// The value of this item.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	line int      <span class="comment">// The line number at the start of this item.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>}
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>func (i item) String() string {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	switch {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	case i.typ == itemEOF:
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		return &#34;EOF&#34;
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	case i.typ == itemError:
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		return i.val
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	case i.typ &gt; itemKeyword:
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		return fmt.Sprintf(&#34;&lt;%s&gt;&#34;, i.val)
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	case len(i.val) &gt; 10:
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		return fmt.Sprintf(&#34;%.10q...&#34;, i.val)
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;%q&#34;, i.val)
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// itemType identifies the type of lex items.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>type itemType int
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>const (
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	itemError        itemType = iota <span class="comment">// error occurred; value is text of error</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	itemBool                         <span class="comment">// boolean constant</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	itemChar                         <span class="comment">// printable ASCII character; grab bag for comma etc.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	itemCharConstant                 <span class="comment">// character constant</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	itemComment                      <span class="comment">// comment text</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	itemComplex                      <span class="comment">// complex constant (1+2i); imaginary is just a number</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	itemAssign                       <span class="comment">// equals (&#39;=&#39;) introducing an assignment</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	itemDeclare                      <span class="comment">// colon-equals (&#39;:=&#39;) introducing a declaration</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	itemEOF
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	itemField      <span class="comment">// alphanumeric identifier starting with &#39;.&#39;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	itemIdentifier <span class="comment">// alphanumeric identifier not starting with &#39;.&#39;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	itemLeftDelim  <span class="comment">// left action delimiter</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	itemLeftParen  <span class="comment">// &#39;(&#39; inside action</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	itemNumber     <span class="comment">// simple number, including imaginary</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	itemPipe       <span class="comment">// pipe symbol</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	itemRawString  <span class="comment">// raw quoted string (includes quotes)</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	itemRightDelim <span class="comment">// right action delimiter</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	itemRightParen <span class="comment">// &#39;)&#39; inside action</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	itemSpace      <span class="comment">// run of spaces separating arguments</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	itemString     <span class="comment">// quoted string (includes quotes)</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	itemText       <span class="comment">// plain text</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	itemVariable   <span class="comment">// variable starting with &#39;$&#39;, such as &#39;$&#39; or  &#39;$1&#39; or &#39;$hello&#39;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// Keywords appear after all the rest.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	itemKeyword  <span class="comment">// used only to delimit the keywords</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	itemBlock    <span class="comment">// block keyword</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	itemBreak    <span class="comment">// break keyword</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	itemContinue <span class="comment">// continue keyword</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	itemDot      <span class="comment">// the cursor, spelled &#39;.&#39;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	itemDefine   <span class="comment">// define keyword</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	itemElse     <span class="comment">// else keyword</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	itemEnd      <span class="comment">// end keyword</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	itemIf       <span class="comment">// if keyword</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	itemNil      <span class="comment">// the untyped nil constant, easiest to treat as a keyword</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	itemRange    <span class="comment">// range keyword</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	itemTemplate <span class="comment">// template keyword</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	itemWith     <span class="comment">// with keyword</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>var key = map[string]itemType{
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	&#34;.&#34;:        itemDot,
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	&#34;block&#34;:    itemBlock,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	&#34;break&#34;:    itemBreak,
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	&#34;continue&#34;: itemContinue,
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	&#34;define&#34;:   itemDefine,
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	&#34;else&#34;:     itemElse,
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	&#34;end&#34;:      itemEnd,
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	&#34;if&#34;:       itemIf,
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	&#34;range&#34;:    itemRange,
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	&#34;nil&#34;:      itemNil,
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	&#34;template&#34;: itemTemplate,
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	&#34;with&#34;:     itemWith,
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>const eof = -1
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// Trimming spaces.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// If the action begins &#34;{{- &#34; rather than &#34;{{&#34;, then all space/tab/newlines</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// preceding the action are trimmed; conversely if it ends &#34; -}}&#34; the</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// leading spaces are trimmed. This is done entirely in the lexer; the</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// parser never sees it happen. We require an ASCII space (&#39; &#39;, \t, \r, \n)</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// to be present to avoid ambiguity with things like &#34;{{-3}}&#34;. It reads</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// better with the space present anyway. For simplicity, only ASCII</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// does the job.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>const (
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	spaceChars    = &#34; \t\r\n&#34;  <span class="comment">// These are the space characters defined by Go itself.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	trimMarker    = &#39;-&#39;        <span class="comment">// Attached to left/right delimiter, trims trailing spaces from preceding/following text.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	trimMarkerLen = Pos(1 + 1) <span class="comment">// marker plus space before or after</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// stateFn represents the state of the scanner as a function that returns the next state.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>type stateFn func(*lexer) stateFn
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// lexer holds the state of the scanner.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>type lexer struct {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	name         string <span class="comment">// the name of the input; used only for error reports</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	input        string <span class="comment">// the string being scanned</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	leftDelim    string <span class="comment">// start of action marker</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	rightDelim   string <span class="comment">// end of action marker</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	pos          Pos    <span class="comment">// current position in the input</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	start        Pos    <span class="comment">// start position of this item</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	atEOF        bool   <span class="comment">// we have hit the end of input and returned eof</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	parenDepth   int    <span class="comment">// nesting depth of ( ) exprs</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	line         int    <span class="comment">// 1+number of newlines seen</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	startLine    int    <span class="comment">// start line of this item</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	item         item   <span class="comment">// item to return to parser</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	insideAction bool   <span class="comment">// are we inside an action?</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	options      lexOptions
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// lexOptions control behavior of the lexer. All default to false.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>type lexOptions struct {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	emitComment bool <span class="comment">// emit itemComment tokens.</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	breakOK     bool <span class="comment">// break keyword allowed</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	continueOK  bool <span class="comment">// continue keyword allowed</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// next returns the next rune in the input.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>func (l *lexer) next() rune {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	if int(l.pos) &gt;= len(l.input) {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		l.atEOF = true
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		return eof
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	l.pos += Pos(w)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	if r == &#39;\n&#39; {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		l.line++
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	return r
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// peek returns but does not consume the next rune in the input.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>func (l *lexer) peek() rune {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	r := l.next()
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	l.backup()
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	return r
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// backup steps back one rune.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>func (l *lexer) backup() {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if !l.atEOF &amp;&amp; l.pos &gt; 0 {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		r, w := utf8.DecodeLastRuneInString(l.input[:l.pos])
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		l.pos -= Pos(w)
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		<span class="comment">// Correct newline count.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		if r == &#39;\n&#39; {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			l.line--
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// thisItem returns the item at the current input point with the specified type</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">// and advances the input.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>func (l *lexer) thisItem(t itemType) item {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	i := item{t, l.start, l.input[l.start:l.pos], l.startLine}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	l.start = l.pos
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	l.startLine = l.line
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	return i
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">// emit passes the trailing text as an item back to the parser.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>func (l *lexer) emit(t itemType) stateFn {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	return l.emitItem(l.thisItem(t))
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">// emitItem passes the specified item to the parser.</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>func (l *lexer) emitItem(i item) stateFn {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	l.item = i
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	return nil
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// ignore skips over the pending input before this point.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// It tracks newlines in the ignored text, so use it only</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// for text that is skipped without calling l.next.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>func (l *lexer) ignore() {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	l.line += strings.Count(l.input[l.start:l.pos], &#34;\n&#34;)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	l.start = l.pos
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	l.startLine = l.line
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// accept consumes the next rune if it&#39;s from the valid set.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>func (l *lexer) accept(valid string) bool {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	if strings.ContainsRune(valid, l.next()) {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		return true
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	l.backup()
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	return false
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// acceptRun consumes a run of runes from the valid set.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>func (l *lexer) acceptRun(valid string) {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	for strings.ContainsRune(valid, l.next()) {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	l.backup()
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// errorf returns an error token and terminates the scan by passing</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// back a nil pointer that will be the next state, terminating l.nextItem.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>func (l *lexer) errorf(format string, args ...any) stateFn {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	l.item = item{itemError, l.start, fmt.Sprintf(format, args...), l.startLine}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	l.start = 0
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	l.pos = 0
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	l.input = l.input[:0]
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	return nil
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// nextItem returns the next item from the input.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// Called by the parser, not in the lexing goroutine.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>func (l *lexer) nextItem() item {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	l.item = item{itemEOF, l.pos, &#34;EOF&#34;, l.startLine}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	state := lexText
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	if l.insideAction {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		state = lexInsideAction
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	for {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		state = state(l)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		if state == nil {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			return l.item
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span><span class="comment">// lex creates a new scanner for the input string.</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>func lex(name, input, left, right string) *lexer {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if left == &#34;&#34; {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		left = leftDelim
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	if right == &#34;&#34; {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		right = rightDelim
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	l := &amp;lexer{
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		name:         name,
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		input:        input,
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		leftDelim:    left,
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		rightDelim:   right,
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		line:         1,
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		startLine:    1,
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		insideAction: false,
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	return l
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// state functions</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>const (
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	leftDelim    = &#34;{{&#34;
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	rightDelim   = &#34;}}&#34;
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	leftComment  = &#34;/*&#34;
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	rightComment = &#34;*/&#34;
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>)
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">// lexText scans until an opening action delimiter, &#34;{{&#34;.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>func lexText(l *lexer) stateFn {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	if x := strings.Index(l.input[l.pos:], l.leftDelim); x &gt;= 0 {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		if x &gt; 0 {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			l.pos += Pos(x)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			<span class="comment">// Do we trim any trailing space?</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			trimLength := Pos(0)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			delimEnd := l.pos + Pos(len(l.leftDelim))
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			if hasLeftTrimMarker(l.input[delimEnd:]) {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>				trimLength = rightTrimLength(l.input[l.start:l.pos])
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			l.pos -= trimLength
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			l.line += strings.Count(l.input[l.start:l.pos], &#34;\n&#34;)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			i := l.thisItem(itemText)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			l.pos += trimLength
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			l.ignore()
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			if len(i.val) &gt; 0 {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>				return l.emitItem(i)
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		return lexLeftDelim
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	l.pos = Pos(len(l.input))
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// Correctly reached EOF.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	if l.pos &gt; l.start {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		l.line += strings.Count(l.input[l.start:l.pos], &#34;\n&#34;)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		return l.emit(itemText)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	return l.emit(itemEOF)
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">// rightTrimLength returns the length of the spaces at the end of the string.</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>func rightTrimLength(s string) Pos {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	return Pos(len(s) - len(strings.TrimRight(s, spaceChars)))
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">// atRightDelim reports whether the lexer is at a right delimiter, possibly preceded by a trim marker.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>func (l *lexer) atRightDelim() (delim, trimSpaces bool) {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if hasRightTrimMarker(l.input[l.pos:]) &amp;&amp; strings.HasPrefix(l.input[l.pos+trimMarkerLen:], l.rightDelim) { <span class="comment">// With trim marker.</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		return true, true
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	if strings.HasPrefix(l.input[l.pos:], l.rightDelim) { <span class="comment">// Without trim marker.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		return true, false
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	return false, false
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">// leftTrimLength returns the length of the spaces at the beginning of the string.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>func leftTrimLength(s string) Pos {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	return Pos(len(s) - len(strings.TrimLeft(s, spaceChars)))
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span><span class="comment">// lexLeftDelim scans the left delimiter, which is known to be present, possibly with a trim marker.</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span><span class="comment">// (The text to be trimmed has already been emitted.)</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>func lexLeftDelim(l *lexer) stateFn {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	l.pos += Pos(len(l.leftDelim))
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	trimSpace := hasLeftTrimMarker(l.input[l.pos:])
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	afterMarker := Pos(0)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	if trimSpace {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		afterMarker = trimMarkerLen
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	if strings.HasPrefix(l.input[l.pos+afterMarker:], leftComment) {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		l.pos += afterMarker
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		l.ignore()
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		return lexComment
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	i := l.thisItem(itemLeftDelim)
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	l.insideAction = true
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	l.pos += afterMarker
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	l.ignore()
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	l.parenDepth = 0
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	return l.emitItem(i)
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">// lexComment scans a comment. The left comment marker is known to be present.</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>func lexComment(l *lexer) stateFn {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	l.pos += Pos(len(leftComment))
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	x := strings.Index(l.input[l.pos:], rightComment)
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		return l.errorf(&#34;unclosed comment&#34;)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	l.pos += Pos(x + len(rightComment))
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	delim, trimSpace := l.atRightDelim()
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	if !delim {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		return l.errorf(&#34;comment ends before closing delimiter&#34;)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	i := l.thisItem(itemComment)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	if trimSpace {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		l.pos += trimMarkerLen
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	l.pos += Pos(len(l.rightDelim))
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	if trimSpace {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		l.pos += leftTrimLength(l.input[l.pos:])
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	l.ignore()
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	if l.options.emitComment {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		return l.emitItem(i)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	return lexText
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">// lexRightDelim scans the right delimiter, which is known to be present, possibly with a trim marker.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>func lexRightDelim(l *lexer) stateFn {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	_, trimSpace := l.atRightDelim()
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	if trimSpace {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		l.pos += trimMarkerLen
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		l.ignore()
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	l.pos += Pos(len(l.rightDelim))
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	i := l.thisItem(itemRightDelim)
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	if trimSpace {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		l.pos += leftTrimLength(l.input[l.pos:])
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		l.ignore()
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	l.insideAction = false
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	return l.emitItem(i)
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span><span class="comment">// lexInsideAction scans the elements inside action delimiters.</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>func lexInsideAction(l *lexer) stateFn {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	<span class="comment">// Either number, quoted string, or identifier.</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">// Spaces separate arguments; runs of spaces turn into itemSpace.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	<span class="comment">// Pipe symbols separate and are emitted.</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	delim, _ := l.atRightDelim()
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	if delim {
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		if l.parenDepth == 0 {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			return lexRightDelim
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		return l.errorf(&#34;unclosed left paren&#34;)
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	switch r := l.next(); {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	case r == eof:
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		return l.errorf(&#34;unclosed action&#34;)
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	case isSpace(r):
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		l.backup() <span class="comment">// Put space back in case we have &#34; -}}&#34;.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		return lexSpace
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	case r == &#39;=&#39;:
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		return l.emit(itemAssign)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	case r == &#39;:&#39;:
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		if l.next() != &#39;=&#39; {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			return l.errorf(&#34;expected :=&#34;)
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		return l.emit(itemDeclare)
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	case r == &#39;|&#39;:
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		return l.emit(itemPipe)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	case r == &#39;&#34;&#39;:
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		return lexQuote
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	case r == &#39;`&#39;:
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		return lexRawQuote
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	case r == &#39;$&#39;:
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		return lexVariable
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	case r == &#39;\&#39;&#39;:
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		return lexChar
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	case r == &#39;.&#39;:
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		<span class="comment">// special look-ahead for &#34;.field&#34; so we don&#39;t break l.backup().</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		if l.pos &lt; Pos(len(l.input)) {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>			r := l.input[l.pos]
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			if r &lt; &#39;0&#39; || &#39;9&#39; &lt; r {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>				return lexField
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		fallthrough <span class="comment">// &#39;.&#39; can start a number.</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	case r == &#39;+&#39; || r == &#39;-&#39; || (&#39;0&#39; &lt;= r &amp;&amp; r &lt;= &#39;9&#39;):
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		l.backup()
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		return lexNumber
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	case isAlphaNumeric(r):
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		l.backup()
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		return lexIdentifier
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	case r == &#39;(&#39;:
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		l.parenDepth++
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		return l.emit(itemLeftParen)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	case r == &#39;)&#39;:
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		l.parenDepth--
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		if l.parenDepth &lt; 0 {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>			return l.errorf(&#34;unexpected right paren&#34;)
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		return l.emit(itemRightParen)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	case r &lt;= unicode.MaxASCII &amp;&amp; unicode.IsPrint(r):
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		return l.emit(itemChar)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	default:
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		return l.errorf(&#34;unrecognized character in action: %#U&#34;, r)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span><span class="comment">// lexSpace scans a run of space characters.</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span><span class="comment">// We have not consumed the first space, which is known to be present.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span><span class="comment">// Take care if there is a trim-marked right delimiter, which starts with a space.</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>func lexSpace(l *lexer) stateFn {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	var r rune
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	var numSpaces int
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	for {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		r = l.peek()
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		if !isSpace(r) {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>			break
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		l.next()
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		numSpaces++
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	<span class="comment">// Be careful about a trim-marked closing delimiter, which has a minus</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	<span class="comment">// after a space. We know there is a space, so check for the &#39;-&#39; that might follow.</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	if hasRightTrimMarker(l.input[l.pos-1:]) &amp;&amp; strings.HasPrefix(l.input[l.pos-1+trimMarkerLen:], l.rightDelim) {
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		l.backup() <span class="comment">// Before the space.</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		if numSpaces == 1 {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			return lexRightDelim <span class="comment">// On the delim, so go right to that.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	return l.emit(itemSpace)
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span><span class="comment">// lexIdentifier scans an alphanumeric.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>func lexIdentifier(l *lexer) stateFn {
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	for {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		switch r := l.next(); {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		case isAlphaNumeric(r):
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			<span class="comment">// absorb.</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		default:
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			l.backup()
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			word := l.input[l.start:l.pos]
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>			if !l.atTerminator() {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>				return l.errorf(&#34;bad character %#U&#34;, r)
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>			}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>			switch {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>			case key[word] &gt; itemKeyword:
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>				item := key[word]
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>				if item == itemBreak &amp;&amp; !l.options.breakOK || item == itemContinue &amp;&amp; !l.options.continueOK {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>					return l.emit(itemIdentifier)
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>				}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>				return l.emit(item)
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>			case word[0] == &#39;.&#39;:
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>				return l.emit(itemField)
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			case word == &#34;true&#34;, word == &#34;false&#34;:
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>				return l.emit(itemBool)
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>			default:
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>				return l.emit(itemIdentifier)
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span><span class="comment">// lexField scans a field: .Alphanumeric.</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span><span class="comment">// The . has been scanned.</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>func lexField(l *lexer) stateFn {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	return lexFieldOrVariable(l, itemField)
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span><span class="comment">// lexVariable scans a Variable: $Alphanumeric.</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span><span class="comment">// The $ has been scanned.</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>func lexVariable(l *lexer) stateFn {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	if l.atTerminator() { <span class="comment">// Nothing interesting follows -&gt; &#34;$&#34;.</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		return l.emit(itemVariable)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	}
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	return lexFieldOrVariable(l, itemVariable)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span><span class="comment">// lexFieldOrVariable scans a field or variable: [.$]Alphanumeric.</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">// The . or $ has been scanned.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>func lexFieldOrVariable(l *lexer, typ itemType) stateFn {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	if l.atTerminator() { <span class="comment">// Nothing interesting follows -&gt; &#34;.&#34; or &#34;$&#34;.</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		if typ == itemVariable {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>			return l.emit(itemVariable)
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		}
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		return l.emit(itemDot)
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	}
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	var r rune
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	for {
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		r = l.next()
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		if !isAlphaNumeric(r) {
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			l.backup()
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>			break
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		}
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	if !l.atTerminator() {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		return l.errorf(&#34;bad character %#U&#34;, r)
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	}
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	return l.emit(typ)
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>}
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span><span class="comment">// atTerminator reports whether the input is at valid termination character to</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span><span class="comment">// appear after an identifier. Breaks .X.Y into two pieces. Also catches cases</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span><span class="comment">// like &#34;$x+2&#34; not being acceptable without a space, in case we decide one</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span><span class="comment">// day to implement arithmetic.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>func (l *lexer) atTerminator() bool {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	r := l.peek()
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	if isSpace(r) {
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		return true
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	}
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	switch r {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	case eof, &#39;.&#39;, &#39;,&#39;, &#39;|&#39;, &#39;:&#39;, &#39;)&#39;, &#39;(&#39;:
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		return true
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	return strings.HasPrefix(l.input[l.pos:], l.rightDelim)
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span><span class="comment">// lexChar scans a character constant. The initial quote is already</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span><span class="comment">// scanned. Syntax checking is done by the parser.</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>func lexChar(l *lexer) stateFn {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>Loop:
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	for {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		switch l.next() {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		case &#39;\\&#39;:
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>			if r := l.next(); r != eof &amp;&amp; r != &#39;\n&#39; {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>				break
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>			}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>			fallthrough
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		case eof, &#39;\n&#39;:
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>			return l.errorf(&#34;unterminated character constant&#34;)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		case &#39;\&#39;&#39;:
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>			break Loop
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		}
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	}
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	return l.emit(itemCharConstant)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span><span class="comment">// lexNumber scans a number: decimal, octal, hex, float, or imaginary. This</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span><span class="comment">// isn&#39;t a perfect number scanner - for instance it accepts &#34;.&#34; and &#34;0x0.2&#34;</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span><span class="comment">// and &#34;089&#34; - but when it&#39;s wrong the input is invalid and the parser (via</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span><span class="comment">// strconv) will notice.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>func lexNumber(l *lexer) stateFn {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	if !l.scanNumber() {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		return l.errorf(&#34;bad number syntax: %q&#34;, l.input[l.start:l.pos])
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	}
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	if sign := l.peek(); sign == &#39;+&#39; || sign == &#39;-&#39; {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		<span class="comment">// Complex: 1+2i. No spaces, must end in &#39;i&#39;.</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		if !l.scanNumber() || l.input[l.pos-1] != &#39;i&#39; {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>			return l.errorf(&#34;bad number syntax: %q&#34;, l.input[l.start:l.pos])
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		return l.emit(itemComplex)
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	}
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	return l.emit(itemNumber)
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>func (l *lexer) scanNumber() bool {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	<span class="comment">// Optional leading sign.</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	l.accept(&#34;+-&#34;)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	<span class="comment">// Is it hex?</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	digits := &#34;0123456789_&#34;
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	if l.accept(&#34;0&#34;) {
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		<span class="comment">// Note: Leading 0 does not mean octal in floats.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		if l.accept(&#34;xX&#34;) {
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>			digits = &#34;0123456789abcdefABCDEF_&#34;
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		} else if l.accept(&#34;oO&#34;) {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>			digits = &#34;01234567_&#34;
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		} else if l.accept(&#34;bB&#34;) {
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>			digits = &#34;01_&#34;
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		}
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	}
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	l.acceptRun(digits)
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	if l.accept(&#34;.&#34;) {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		l.acceptRun(digits)
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	}
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	if len(digits) == 10+1 &amp;&amp; l.accept(&#34;eE&#34;) {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>		l.accept(&#34;+-&#34;)
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		l.acceptRun(&#34;0123456789_&#34;)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	if len(digits) == 16+6+1 &amp;&amp; l.accept(&#34;pP&#34;) {
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		l.accept(&#34;+-&#34;)
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		l.acceptRun(&#34;0123456789_&#34;)
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	<span class="comment">// Is it imaginary?</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	l.accept(&#34;i&#34;)
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	<span class="comment">// Next thing mustn&#39;t be alphanumeric.</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	if isAlphaNumeric(l.peek()) {
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		l.next()
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>		return false
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	}
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	return true
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>}
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span><span class="comment">// lexQuote scans a quoted string.</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>func lexQuote(l *lexer) stateFn {
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>Loop:
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	for {
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		switch l.next() {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>		case &#39;\\&#39;:
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>			if r := l.next(); r != eof &amp;&amp; r != &#39;\n&#39; {
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>				break
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>			}
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>			fallthrough
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		case eof, &#39;\n&#39;:
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>			return l.errorf(&#34;unterminated quoted string&#34;)
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		case &#39;&#34;&#39;:
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>			break Loop
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	}
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	return l.emit(itemString)
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span><span class="comment">// lexRawQuote scans a raw quoted string.</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>func lexRawQuote(l *lexer) stateFn {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>Loop:
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	for {
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		switch l.next() {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		case eof:
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>			return l.errorf(&#34;unterminated raw quoted string&#34;)
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		case &#39;`&#39;:
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>			break Loop
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		}
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	return l.emit(itemRawString)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>}
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span><span class="comment">// isSpace reports whether r is a space character.</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>func isSpace(r rune) bool {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	return r == &#39; &#39; || r == &#39;\t&#39; || r == &#39;\r&#39; || r == &#39;\n&#39;
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span><span class="comment">// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>func isAlphaNumeric(r rune) bool {
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	return r == &#39;_&#39; || unicode.IsLetter(r) || unicode.IsDigit(r)
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>}
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>func hasLeftTrimMarker(s string) bool {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	return len(s) &gt;= 2 &amp;&amp; s[0] == trimMarker &amp;&amp; isSpace(rune(s[1]))
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>}
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>func hasRightTrimMarker(s string) bool {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	return len(s) &gt;= 2 &amp;&amp; isSpace(rune(s[0])) &amp;&amp; s[1] == trimMarker
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>
</pre><p><a href="lex.go?m=text">View as plain text</a></p>

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
