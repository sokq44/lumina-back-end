<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/symtabinl.go - Go Documentation Server</title>

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
<a href="symtabinl.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">symtabinl.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2023 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;internal/abi&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// inlinedCall is the encoding of entries in the FUNCDATA_InlTree table.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>type inlinedCall struct {
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	funcID    abi.FuncID <span class="comment">// type of the called function</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	_         [3]byte
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	nameOff   int32 <span class="comment">// offset into pclntab for name of called function</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	parentPc  int32 <span class="comment">// position of an instruction whose source position is the call site (offset from entry)</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	startLine int32 <span class="comment">// line number of start of function (func keyword/TEXT directive)</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>}
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// An inlineUnwinder iterates over the stack of inlined calls at a PC by</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// decoding the inline table. The last step of iteration is always the frame of</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// the physical function, so there&#39;s always at least one frame.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// This is typically used as:</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//	for u, uf := newInlineUnwinder(...); uf.valid(); uf = u.next(uf) { ... }</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// Implementation note: This is used in contexts that disallow write barriers.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// Hence, the constructor returns this by value and pointer receiver methods</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// must not mutate pointer fields. Also, we keep the mutable state in a separate</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// struct mostly to keep both structs SSA-able, which generates much better</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// code.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>type inlineUnwinder struct {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	f       funcInfo
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	inlTree *[1 &lt;&lt; 20]inlinedCall
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// An inlineFrame is a position in an inlineUnwinder.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>type inlineFrame struct {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// pc is the PC giving the file/line metadata of the current frame. This is</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// always a &#34;call PC&#34; (not a &#34;return PC&#34;). This is 0 when the iterator is</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// exhausted.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	pc uintptr
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// index is the index of the current record in inlTree, or -1 if we are in</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// the outermost function.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	index int32
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// newInlineUnwinder creates an inlineUnwinder initially set to the inner-most</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// inlined frame at PC. PC should be a &#34;call PC&#34; (not a &#34;return PC&#34;).</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// This unwinder uses non-strict handling of PC because it&#39;s assumed this is</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// only ever used for symbolic debugging. If things go really wrong, it&#39;ll just</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// fall back to the outermost frame.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func newInlineUnwinder(f funcInfo, pc uintptr) (inlineUnwinder, inlineFrame) {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	inldata := funcdata(f, abi.FUNCDATA_InlTree)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	if inldata == nil {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		return inlineUnwinder{f: f}, inlineFrame{pc: pc, index: -1}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	inlTree := (*[1 &lt;&lt; 20]inlinedCall)(inldata)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	u := inlineUnwinder{f: f, inlTree: inlTree}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	return u, u.resolveInternal(pc)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>func (u *inlineUnwinder) resolveInternal(pc uintptr) inlineFrame {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	return inlineFrame{
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		pc: pc,
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		<span class="comment">// Conveniently, this returns -1 if there&#39;s an error, which is the same</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		<span class="comment">// value we use for the outermost frame.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		index: pcdatavalue1(u.f, abi.PCDATA_InlTreeIndex, pc, false),
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>func (uf inlineFrame) valid() bool {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return uf.pc != 0
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// next returns the frame representing uf&#39;s logical caller.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>func (u *inlineUnwinder) next(uf inlineFrame) inlineFrame {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	if uf.index &lt; 0 {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		uf.pc = 0
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		return uf
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	parentPc := u.inlTree[uf.index].parentPc
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	return u.resolveInternal(u.f.entry() + uintptr(parentPc))
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// isInlined returns whether uf is an inlined frame.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>func (u *inlineUnwinder) isInlined(uf inlineFrame) bool {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	return uf.index &gt;= 0
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// srcFunc returns the srcFunc representing the given frame.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>func (u *inlineUnwinder) srcFunc(uf inlineFrame) srcFunc {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	if uf.index &lt; 0 {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		return u.f.srcFunc()
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	t := &amp;u.inlTree[uf.index]
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	return srcFunc{
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		u.f.datap,
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		t.nameOff,
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		t.startLine,
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		t.funcID,
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// fileLine returns the file name and line number of the call within the given</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// frame. As a convenience, for the innermost frame, it returns the file and</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// line of the PC this unwinder was started at (often this is a call to another</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// physical function).</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// It returns &#34;?&#34;, 0 if something goes wrong.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>func (u *inlineUnwinder) fileLine(uf inlineFrame) (file string, line int) {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	file, line32 := funcline1(u.f, uf.pc, false)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	return file, int(line32)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
</pre><p><a href="symtabinl.go?m=text">View as plain text</a></p>

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
