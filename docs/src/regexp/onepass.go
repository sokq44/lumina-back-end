<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/regexp/onepass.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../index.html">GoDoc</a></div>
<a href="onepass.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/regexp">regexp</a>/<span class="text-muted">onepass.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/regexp">regexp</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2014 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package regexp
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;regexp/syntax&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// &#34;One-pass&#34; regexp execution.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// Some regexps can be analyzed to determine that they never need</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// backtracking: they are guaranteed to run in one pass over the string</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// without bothering to save all the usual NFA state.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// Detect those and execute them more quickly.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// A onePassProg is a compiled one-pass regular expression program.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// It is the same as syntax.Prog except for the use of onePassInst.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>type onePassProg struct {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	Inst   []onePassInst
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	Start  int <span class="comment">// index of start instruction</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	NumCap int <span class="comment">// number of InstCapture insts in re</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// A onePassInst is a single instruction in a one-pass regular expression program.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// It is the same as syntax.Inst except for the new &#39;Next&#39; field.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>type onePassInst struct {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	syntax.Inst
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	Next []uint32
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// onePassPrefix returns a literal string that all matches for the</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// regexp must start with. Complete is true if the prefix</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// is the entire match. Pc is the index of the last rune instruction</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// in the string. The onePassPrefix skips over the mandatory</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// EmptyBeginText.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>func onePassPrefix(p *syntax.Prog) (prefix string, complete bool, pc uint32) {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	i := &amp;p.Inst[p.Start]
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	if i.Op != syntax.InstEmptyWidth || (syntax.EmptyOp(i.Arg))&amp;syntax.EmptyBeginText == 0 {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		return &#34;&#34;, i.Op == syntax.InstMatch, uint32(p.Start)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	pc = i.Out
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	i = &amp;p.Inst[pc]
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	for i.Op == syntax.InstNop {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		pc = i.Out
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		i = &amp;p.Inst[pc]
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// Avoid allocation of buffer if prefix is empty.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	if iop(i) != syntax.InstRune || len(i.Rune) != 1 {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		return &#34;&#34;, i.Op == syntax.InstMatch, uint32(p.Start)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// Have prefix; gather characters.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	var buf strings.Builder
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	for iop(i) == syntax.InstRune &amp;&amp; len(i.Rune) == 1 &amp;&amp; syntax.Flags(i.Arg)&amp;syntax.FoldCase == 0 &amp;&amp; i.Rune[0] != utf8.RuneError {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		buf.WriteRune(i.Rune[0])
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		pc, i = i.Out, &amp;p.Inst[i.Out]
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if i.Op == syntax.InstEmptyWidth &amp;&amp;
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		syntax.EmptyOp(i.Arg)&amp;syntax.EmptyEndText != 0 &amp;&amp;
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		p.Inst[i.Out].Op == syntax.InstMatch {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		complete = true
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	return buf.String(), complete, pc
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// onePassNext selects the next actionable state of the prog, based on the input character.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// It should only be called when i.Op == InstAlt or InstAltMatch, and from the one-pass machine.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// One of the alternates may ultimately lead without input to end of line. If the instruction</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// is InstAltMatch the path to the InstMatch is in i.Out, the normal node in i.Next.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>func onePassNext(i *onePassInst, r rune) uint32 {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	next := i.MatchRunePos(r)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	if next &gt;= 0 {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		return i.Next[next]
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if i.Op == syntax.InstAltMatch {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		return i.Out
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	return 0
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>func iop(i *syntax.Inst) syntax.InstOp {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	op := i.Op
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	switch op {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	case syntax.InstRune1, syntax.InstRuneAny, syntax.InstRuneAnyNotNL:
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		op = syntax.InstRune
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	return op
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// Sparse Array implementation is used as a queueOnePass.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>type queueOnePass struct {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	sparse          []uint32
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	dense           []uint32
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	size, nextIndex uint32
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>func (q *queueOnePass) empty() bool {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	return q.nextIndex &gt;= q.size
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>func (q *queueOnePass) next() (n uint32) {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	n = q.dense[q.nextIndex]
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	q.nextIndex++
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	return
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>func (q *queueOnePass) clear() {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	q.size = 0
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	q.nextIndex = 0
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>func (q *queueOnePass) contains(u uint32) bool {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	if u &gt;= uint32(len(q.sparse)) {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		return false
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	return q.sparse[u] &lt; q.size &amp;&amp; q.dense[q.sparse[u]] == u
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func (q *queueOnePass) insert(u uint32) {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	if !q.contains(u) {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		q.insertNew(u)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>func (q *queueOnePass) insertNew(u uint32) {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	if u &gt;= uint32(len(q.sparse)) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		return
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	q.sparse[u] = q.size
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	q.dense[q.size] = u
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	q.size++
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>func newQueue(size int) (q *queueOnePass) {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	return &amp;queueOnePass{
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		sparse: make([]uint32, size),
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		dense:  make([]uint32, size),
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// mergeRuneSets merges two non-intersecting runesets, and returns the merged result,</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// and a NextIp array. The idea is that if a rune matches the OnePassRunes at index</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// i, NextIp[i/2] is the target. If the input sets intersect, an empty runeset and a</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// NextIp array with the single element mergeFailed is returned.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// The code assumes that both inputs contain ordered and non-intersecting rune pairs.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>const mergeFailed = uint32(0xffffffff)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>var (
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	noRune = []rune{}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	noNext = []uint32{mergeFailed}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>func mergeRuneSets(leftRunes, rightRunes *[]rune, leftPC, rightPC uint32) ([]rune, []uint32) {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	leftLen := len(*leftRunes)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	rightLen := len(*rightRunes)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if leftLen&amp;0x1 != 0 || rightLen&amp;0x1 != 0 {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		panic(&#34;mergeRuneSets odd length []rune&#34;)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	var (
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		lx, rx int
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	merged := make([]rune, 0)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	next := make([]uint32, 0)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	ok := true
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	defer func() {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		if !ok {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			merged = nil
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			next = nil
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}()
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	ix := -1
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	extend := func(newLow *int, newArray *[]rune, pc uint32) bool {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		if ix &gt; 0 &amp;&amp; (*newArray)[*newLow] &lt;= merged[ix] {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			return false
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		merged = append(merged, (*newArray)[*newLow], (*newArray)[*newLow+1])
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		*newLow += 2
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		ix += 2
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		next = append(next, pc)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		return true
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	for lx &lt; leftLen || rx &lt; rightLen {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		switch {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		case rx &gt;= rightLen:
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			ok = extend(&amp;lx, leftRunes, leftPC)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		case lx &gt;= leftLen:
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			ok = extend(&amp;rx, rightRunes, rightPC)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		case (*rightRunes)[rx] &lt; (*leftRunes)[lx]:
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			ok = extend(&amp;rx, rightRunes, rightPC)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		default:
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			ok = extend(&amp;lx, leftRunes, leftPC)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		if !ok {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			return noRune, noNext
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	return merged, next
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// cleanupOnePass drops working memory, and restores certain shortcut instructions.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>func cleanupOnePass(prog *onePassProg, original *syntax.Prog) {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	for ix, instOriginal := range original.Inst {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		switch instOriginal.Op {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		case syntax.InstAlt, syntax.InstAltMatch, syntax.InstRune:
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		case syntax.InstCapture, syntax.InstEmptyWidth, syntax.InstNop, syntax.InstMatch, syntax.InstFail:
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			prog.Inst[ix].Next = nil
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		case syntax.InstRune1, syntax.InstRuneAny, syntax.InstRuneAnyNotNL:
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			prog.Inst[ix].Next = nil
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			prog.Inst[ix] = onePassInst{Inst: instOriginal}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// onePassCopy creates a copy of the original Prog, as we&#39;ll be modifying it.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>func onePassCopy(prog *syntax.Prog) *onePassProg {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	p := &amp;onePassProg{
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		Start:  prog.Start,
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		NumCap: prog.NumCap,
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		Inst:   make([]onePassInst, len(prog.Inst)),
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	for i, inst := range prog.Inst {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		p.Inst[i] = onePassInst{Inst: inst}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">// rewrites one or more common Prog constructs that enable some otherwise</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">// non-onepass Progs to be onepass. A:BD (for example) means an InstAlt at</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	<span class="comment">// ip A, that points to ips B &amp; C.</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">// A:BC + B:DA =&gt; A:BC + B:CD</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// A:BC + B:DC =&gt; A:DC + B:DC</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	for pc := range p.Inst {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		switch p.Inst[pc].Op {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		default:
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			continue
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		case syntax.InstAlt, syntax.InstAltMatch:
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			<span class="comment">// A:Bx + B:Ay</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			p_A_Other := &amp;p.Inst[pc].Out
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			p_A_Alt := &amp;p.Inst[pc].Arg
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			<span class="comment">// make sure a target is another Alt</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			instAlt := p.Inst[*p_A_Alt]
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			if !(instAlt.Op == syntax.InstAlt || instAlt.Op == syntax.InstAltMatch) {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>				p_A_Alt, p_A_Other = p_A_Other, p_A_Alt
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>				instAlt = p.Inst[*p_A_Alt]
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>				if !(instAlt.Op == syntax.InstAlt || instAlt.Op == syntax.InstAltMatch) {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>					continue
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>				}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			instOther := p.Inst[*p_A_Other]
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			<span class="comment">// Analyzing both legs pointing to Alts is for another day</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			if instOther.Op == syntax.InstAlt || instOther.Op == syntax.InstAltMatch {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>				<span class="comment">// too complicated</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>				continue
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			<span class="comment">// simple empty transition loop</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			<span class="comment">// A:BC + B:DA =&gt; A:BC + B:DC</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			p_B_Alt := &amp;p.Inst[*p_A_Alt].Out
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			p_B_Other := &amp;p.Inst[*p_A_Alt].Arg
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			patch := false
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			if instAlt.Out == uint32(pc) {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>				patch = true
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			} else if instAlt.Arg == uint32(pc) {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				patch = true
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>				p_B_Alt, p_B_Other = p_B_Other, p_B_Alt
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			if patch {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>				*p_B_Alt = *p_A_Other
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			<span class="comment">// empty transition to common target</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			<span class="comment">// A:BC + B:DC =&gt; A:DC + B:DC</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			if *p_A_Other == *p_B_Alt {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>				*p_A_Alt = *p_B_Other
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	return p
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">// runeSlice exists to permit sorting the case-folded rune sets.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>type runeSlice []rune
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>func (p runeSlice) Len() int           { return len(p) }
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>func (p runeSlice) Less(i, j int) bool { return p[i] &lt; p[j] }
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>func (p runeSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>var anyRuneNotNL = []rune{0, &#39;\n&#39; - 1, &#39;\n&#39; + 1, unicode.MaxRune}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>var anyRune = []rune{0, unicode.MaxRune}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span><span class="comment">// makeOnePass creates a onepass Prog, if possible. It is possible if at any alt,</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span><span class="comment">// the match engine can always tell which branch to take. The routine may modify</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span><span class="comment">// p if it is turned into a onepass Prog. If it isn&#39;t possible for this to be a</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span><span class="comment">// onepass Prog, the Prog nil is returned. makeOnePass is recursive</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span><span class="comment">// to the size of the Prog.</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>func makeOnePass(p *onePassProg) *onePassProg {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">// If the machine is very long, it&#39;s not worth the time to check if we can use one pass.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	if len(p.Inst) &gt;= 1000 {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		return nil
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	var (
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		instQueue    = newQueue(len(p.Inst))
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		visitQueue   = newQueue(len(p.Inst))
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		check        func(uint32, []bool) bool
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		onePassRunes = make([][]rune, len(p.Inst))
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	<span class="comment">// check that paths from Alt instructions are unambiguous, and rebuild the new</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	<span class="comment">// program as a onepass program</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	check = func(pc uint32, m []bool) (ok bool) {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		ok = true
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		inst := &amp;p.Inst[pc]
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		if visitQueue.contains(pc) {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			return
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		visitQueue.insert(pc)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		switch inst.Op {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		case syntax.InstAlt, syntax.InstAltMatch:
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			ok = check(inst.Out, m) &amp;&amp; check(inst.Arg, m)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			<span class="comment">// check no-input paths to InstMatch</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			matchOut := m[inst.Out]
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			matchArg := m[inst.Arg]
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			if matchOut &amp;&amp; matchArg {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>				ok = false
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>				break
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			<span class="comment">// Match on empty goes in inst.Out</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>			if matchArg {
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>				inst.Out, inst.Arg = inst.Arg, inst.Out
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>				matchOut, matchArg = matchArg, matchOut
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			if matchOut {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>				m[pc] = true
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>				inst.Op = syntax.InstAltMatch
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			<span class="comment">// build a dispatch operator from the two legs of the alt.</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			onePassRunes[pc], inst.Next = mergeRuneSets(
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>				&amp;onePassRunes[inst.Out], &amp;onePassRunes[inst.Arg], inst.Out, inst.Arg)
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			if len(inst.Next) &gt; 0 &amp;&amp; inst.Next[0] == mergeFailed {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>				ok = false
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>				break
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>			}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		case syntax.InstCapture, syntax.InstNop:
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>			ok = check(inst.Out, m)
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>			m[pc] = m[inst.Out]
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			<span class="comment">// pass matching runes back through these no-ops.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>			onePassRunes[pc] = append([]rune{}, onePassRunes[inst.Out]...)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>			inst.Next = make([]uint32, len(onePassRunes[pc])/2+1)
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			for i := range inst.Next {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>				inst.Next[i] = inst.Out
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		case syntax.InstEmptyWidth:
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			ok = check(inst.Out, m)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			m[pc] = m[inst.Out]
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>			onePassRunes[pc] = append([]rune{}, onePassRunes[inst.Out]...)
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>			inst.Next = make([]uint32, len(onePassRunes[pc])/2+1)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			for i := range inst.Next {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>				inst.Next[i] = inst.Out
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		case syntax.InstMatch, syntax.InstFail:
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			m[pc] = inst.Op == syntax.InstMatch
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		case syntax.InstRune:
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			m[pc] = false
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			if len(inst.Next) &gt; 0 {
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>				break
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>			}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			instQueue.insert(inst.Out)
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			if len(inst.Rune) == 0 {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>				onePassRunes[pc] = []rune{}
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>				inst.Next = []uint32{inst.Out}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>				break
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			runes := make([]rune, 0)
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			if len(inst.Rune) == 1 &amp;&amp; syntax.Flags(inst.Arg)&amp;syntax.FoldCase != 0 {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>				r0 := inst.Rune[0]
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>				runes = append(runes, r0, r0)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>				for r1 := unicode.SimpleFold(r0); r1 != r0; r1 = unicode.SimpleFold(r1) {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>					runes = append(runes, r1, r1)
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>				}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>				sort.Sort(runeSlice(runes))
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>			} else {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>				runes = append(runes, inst.Rune...)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			onePassRunes[pc] = runes
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			inst.Next = make([]uint32, len(onePassRunes[pc])/2+1)
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			for i := range inst.Next {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>				inst.Next[i] = inst.Out
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			inst.Op = syntax.InstRune
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		case syntax.InstRune1:
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			m[pc] = false
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			if len(inst.Next) &gt; 0 {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>				break
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			instQueue.insert(inst.Out)
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			runes := []rune{}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			<span class="comment">// expand case-folded runes</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			if syntax.Flags(inst.Arg)&amp;syntax.FoldCase != 0 {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>				r0 := inst.Rune[0]
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>				runes = append(runes, r0, r0)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>				for r1 := unicode.SimpleFold(r0); r1 != r0; r1 = unicode.SimpleFold(r1) {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>					runes = append(runes, r1, r1)
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>				}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>				sort.Sort(runeSlice(runes))
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>			} else {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>				runes = append(runes, inst.Rune[0], inst.Rune[0])
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			onePassRunes[pc] = runes
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			inst.Next = make([]uint32, len(onePassRunes[pc])/2+1)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			for i := range inst.Next {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>				inst.Next[i] = inst.Out
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			inst.Op = syntax.InstRune
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		case syntax.InstRuneAny:
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>			m[pc] = false
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			if len(inst.Next) &gt; 0 {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>				break
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>			instQueue.insert(inst.Out)
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			onePassRunes[pc] = append([]rune{}, anyRune...)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>			inst.Next = []uint32{inst.Out}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		case syntax.InstRuneAnyNotNL:
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>			m[pc] = false
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			if len(inst.Next) &gt; 0 {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>				break
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>			instQueue.insert(inst.Out)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			onePassRunes[pc] = append([]rune{}, anyRuneNotNL...)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>			inst.Next = make([]uint32, len(onePassRunes[pc])/2+1)
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>			for i := range inst.Next {
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>				inst.Next[i] = inst.Out
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>			}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		}
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		return
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	instQueue.clear()
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	instQueue.insert(uint32(p.Start))
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	m := make([]bool, len(p.Inst))
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	for !instQueue.empty() {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		visitQueue.clear()
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		pc := instQueue.next()
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		if !check(pc, m) {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>			p = nil
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			break
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	if p != nil {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		for i := range p.Inst {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>			p.Inst[i].Rune = onePassRunes[i]
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	return p
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>}
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span><span class="comment">// compileOnePass returns a new *syntax.Prog suitable for onePass execution if the original Prog</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span><span class="comment">// can be recharacterized as a one-pass regexp program, or syntax.nil if the</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">// Prog cannot be converted. For a one pass prog, the fundamental condition that must</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span><span class="comment">// be true is: at any InstAlt, there must be no ambiguity about what branch to  take.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>func compileOnePass(prog *syntax.Prog) (p *onePassProg) {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	if prog.Start == 0 {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		return nil
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	<span class="comment">// onepass regexp is anchored</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	if prog.Inst[prog.Start].Op != syntax.InstEmptyWidth ||
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		syntax.EmptyOp(prog.Inst[prog.Start].Arg)&amp;syntax.EmptyBeginText != syntax.EmptyBeginText {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		return nil
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	<span class="comment">// every instruction leading to InstMatch must be EmptyEndText</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	for _, inst := range prog.Inst {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		opOut := prog.Inst[inst.Out].Op
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		switch inst.Op {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		default:
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			if opOut == syntax.InstMatch {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>				return nil
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		case syntax.InstAlt, syntax.InstAltMatch:
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			if opOut == syntax.InstMatch || prog.Inst[inst.Arg].Op == syntax.InstMatch {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>				return nil
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			}
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		case syntax.InstEmptyWidth:
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			if opOut == syntax.InstMatch {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>				if syntax.EmptyOp(inst.Arg)&amp;syntax.EmptyEndText == syntax.EmptyEndText {
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>					continue
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>				}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>				return nil
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>			}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		}
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	<span class="comment">// Creates a slightly optimized copy of the original Prog</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	<span class="comment">// that cleans up some Prog idioms that block valid onepass programs</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	p = onePassCopy(prog)
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	<span class="comment">// checkAmbiguity on InstAlts, build onepass Prog if possible</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	p = makeOnePass(p)
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	if p != nil {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		cleanupOnePass(p, prog)
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	return p
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
</pre><p><a href="onepass.go?m=text">View as plain text</a></p>

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
