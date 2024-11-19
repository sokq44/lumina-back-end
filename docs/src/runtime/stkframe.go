<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/stkframe.go - Go Documentation Server</title>

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
<a href="stkframe.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">stkframe.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2022 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// A stkframe holds information about a single physical stack frame.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>type stkframe struct {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	<span class="comment">// fn is the function being run in this frame. If there is</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">// inlining, this is the outermost function.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	fn funcInfo
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// pc is the program counter within fn.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// The meaning of this is subtle:</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// - Typically, this frame performed a regular function call</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">//   and this is the return PC (just after the CALL</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">//   instruction). In this case, pc-1 reflects the CALL</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">//   instruction itself and is the correct source of symbolic</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">//   information.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// - If this frame &#34;called&#34; sigpanic, then pc is the</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">//   instruction that panicked, and pc is the correct address</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">//   to use for symbolic information.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// - If this is the innermost frame, then PC is where</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">//   execution will continue, but it may not be the</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">//   instruction following a CALL. This may be from</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">//   cooperative preemption, in which case this is the</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">//   instruction after the call to morestack. Or this may be</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">//   from a signal or an un-started goroutine, in which case</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">//   PC could be any instruction, including the first</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">//   instruction in a function. Conventionally, we use pc-1</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">//   for symbolic information, unless pc == fn.entry(), in</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">//   which case we use pc.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	pc uintptr
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// continpc is the PC where execution will continue in fn, or</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// 0 if execution will not continue in this frame.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// This is usually the same as pc, unless this frame &#34;called&#34;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// sigpanic, in which case it&#39;s either the address of</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// deferreturn or 0 if this frame will never execute again.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// This is the PC to use to look up GC liveness for this frame.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	continpc uintptr
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	lr   uintptr <span class="comment">// program counter at caller aka link register</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	sp   uintptr <span class="comment">// stack pointer at pc</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	fp   uintptr <span class="comment">// stack pointer at caller aka frame pointer</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	varp uintptr <span class="comment">// top of local variables</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	argp uintptr <span class="comment">// pointer to function arguments</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// reflectMethodValue is a partial duplicate of reflect.makeFuncImpl</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// and reflect.methodValue.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>type reflectMethodValue struct {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	fn     uintptr
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	stack  *bitvector <span class="comment">// ptrmap for both args and results</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	argLen uintptr    <span class="comment">// just args</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// argBytes returns the argument frame size for a call to frame.fn.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func (frame *stkframe) argBytes() uintptr {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	if frame.fn.args != abi.ArgsSizeUnknown {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		return uintptr(frame.fn.args)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// This is an uncommon and complicated case. Fall back to fully</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// fetching the argument map to compute its size.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	argMap, _ := frame.argMapInternal()
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	return uintptr(argMap.n) * goarch.PtrSize
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// argMapInternal is used internally by stkframe to fetch special</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// argument maps.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// argMap.n is always populated with the size of the argument map.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// argMap.bytedata is only populated for dynamic argument maps (used</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// by reflect). If the caller requires the argument map, it should use</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// this if non-nil, and otherwise fetch the argument map using the</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// current PC.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// hasReflectStackObj indicates that this frame also has a reflect</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// function stack object, which the caller must synthesize.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>func (frame *stkframe) argMapInternal() (argMap bitvector, hasReflectStackObj bool) {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	f := frame.fn
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if f.args != abi.ArgsSizeUnknown {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		argMap.n = f.args / goarch.PtrSize
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		return
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// Extract argument bitmaps for reflect stubs from the calls they made to reflect.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	switch funcname(f) {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	case &#34;reflect.makeFuncStub&#34;, &#34;reflect.methodValueCall&#34;:
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		<span class="comment">// These take a *reflect.methodValue as their</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		<span class="comment">// context register and immediately save it to 0(SP).</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		<span class="comment">// Get the methodValue from 0(SP).</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		arg0 := frame.sp + sys.MinFrameSize
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		minSP := frame.fp
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		if !usesLR {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			<span class="comment">// The CALL itself pushes a word.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			<span class="comment">// Undo that adjustment.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			minSP -= goarch.PtrSize
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		if arg0 &gt;= minSP {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			<span class="comment">// The function hasn&#39;t started yet.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			<span class="comment">// This only happens if f was the</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			<span class="comment">// start function of a new goroutine</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			<span class="comment">// that hasn&#39;t run yet *and* f takes</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			<span class="comment">// no arguments and has no results</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			<span class="comment">// (otherwise it will get wrapped in a</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			<span class="comment">// closure). In this case, we can&#39;t</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			<span class="comment">// reach into its locals because it</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			<span class="comment">// doesn&#39;t have locals yet, but we</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			<span class="comment">// also know its argument map is</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			<span class="comment">// empty.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			if frame.pc != f.entry() {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>				print(&#34;runtime: confused by &#34;, funcname(f), &#34;: no frame (sp=&#34;, hex(frame.sp), &#34; fp=&#34;, hex(frame.fp), &#34;) at entry+&#34;, hex(frame.pc-f.entry()), &#34;\n&#34;)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>				throw(&#34;reflect mismatch&#34;)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			return bitvector{}, false <span class="comment">// No locals, so also no stack objects</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		hasReflectStackObj = true
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		mv := *(**reflectMethodValue)(unsafe.Pointer(arg0))
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// Figure out whether the return values are valid.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// Reflect will update this value after it copies</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// in the return values.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		retValid := *(*bool)(unsafe.Pointer(arg0 + 4*goarch.PtrSize))
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		if mv.fn != f.entry() {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			print(&#34;runtime: confused by &#34;, funcname(f), &#34;\n&#34;)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			throw(&#34;reflect mismatch&#34;)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		argMap = *mv.stack
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		if !retValid {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			<span class="comment">// argMap.n includes the results, but</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			<span class="comment">// those aren&#39;t valid, so drop them.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			n := int32((mv.argLen &amp;^ (goarch.PtrSize - 1)) / goarch.PtrSize)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			if n &lt; argMap.n {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				argMap.n = n
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	return
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// getStackMap returns the locals and arguments live pointer maps, and</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// stack object list for frame.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>func (frame *stkframe) getStackMap(debug bool) (locals, args bitvector, objs []stackObjectRecord) {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	targetpc := frame.continpc
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if targetpc == 0 {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		<span class="comment">// Frame is dead. Return empty bitvectors.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		return
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	f := frame.fn
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	pcdata := int32(-1)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	if targetpc != f.entry() {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		<span class="comment">// Back up to the CALL. If we&#39;re at the function entry</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// point, we want to use the entry map (-1), even if</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		<span class="comment">// the first instruction of the function changes the</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		<span class="comment">// stack map.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		targetpc--
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		pcdata = pcdatavalue(f, abi.PCDATA_StackMapIndex, targetpc)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	if pcdata == -1 {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		<span class="comment">// We do not have a valid pcdata value but there might be a</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		<span class="comment">// stackmap for this function. It is likely that we are looking</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		<span class="comment">// at the function prologue, assume so and hope for the best.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		pcdata = 0
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// Local variables.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	size := frame.varp - frame.sp
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	var minsize uintptr
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	switch goarch.ArchFamily {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	case goarch.ARM64:
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		minsize = sys.StackAlign
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	default:
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		minsize = sys.MinFrameSize
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	if size &gt; minsize {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		stackid := pcdata
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		stkmap := (*stackmap)(funcdata(f, abi.FUNCDATA_LocalsPointerMaps))
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		if stkmap == nil || stkmap.n &lt;= 0 {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			print(&#34;runtime: frame &#34;, funcname(f), &#34; untyped locals &#34;, hex(frame.varp-size), &#34;+&#34;, hex(size), &#34;\n&#34;)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			throw(&#34;missing stackmap&#34;)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		<span class="comment">// If nbit == 0, there&#39;s no work to do.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		if stkmap.nbit &gt; 0 {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			if stackid &lt; 0 || stackid &gt;= stkmap.n {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>				<span class="comment">// don&#39;t know where we are</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>				print(&#34;runtime: pcdata is &#34;, stackid, &#34; and &#34;, stkmap.n, &#34; locals stack map entries for &#34;, funcname(f), &#34; (targetpc=&#34;, hex(targetpc), &#34;)\n&#34;)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>				throw(&#34;bad symbol table&#34;)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			locals = stackmapdata(stkmap, stackid)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			if stackDebug &gt;= 3 &amp;&amp; debug {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>				print(&#34;      locals &#34;, stackid, &#34;/&#34;, stkmap.n, &#34; &#34;, locals.n, &#34; words &#34;, locals.bytedata, &#34;\n&#34;)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		} else if stackDebug &gt;= 3 &amp;&amp; debug {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			print(&#34;      no locals to adjust\n&#34;)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// Arguments. First fetch frame size and special-case argument maps.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	var isReflect bool
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	args, isReflect = frame.argMapInternal()
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	if args.n &gt; 0 &amp;&amp; args.bytedata == nil {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		<span class="comment">// Non-empty argument frame, but not a special map.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// Fetch the argument map at pcdata.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		stackmap := (*stackmap)(funcdata(f, abi.FUNCDATA_ArgsPointerMaps))
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		if stackmap == nil || stackmap.n &lt;= 0 {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			print(&#34;runtime: frame &#34;, funcname(f), &#34; untyped args &#34;, hex(frame.argp), &#34;+&#34;, hex(args.n*goarch.PtrSize), &#34;\n&#34;)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			throw(&#34;missing stackmap&#34;)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		if pcdata &lt; 0 || pcdata &gt;= stackmap.n {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			<span class="comment">// don&#39;t know where we are</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			print(&#34;runtime: pcdata is &#34;, pcdata, &#34; and &#34;, stackmap.n, &#34; args stack map entries for &#34;, funcname(f), &#34; (targetpc=&#34;, hex(targetpc), &#34;)\n&#34;)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			throw(&#34;bad symbol table&#34;)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		if stackmap.nbit == 0 {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			args.n = 0
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		} else {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			args = stackmapdata(stackmap, pcdata)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// stack objects.</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	if (GOARCH == &#34;amd64&#34; || GOARCH == &#34;arm64&#34; || GOARCH == &#34;loong64&#34; || GOARCH == &#34;ppc64&#34; || GOARCH == &#34;ppc64le&#34; || GOARCH == &#34;riscv64&#34;) &amp;&amp;
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		unsafe.Sizeof(abi.RegArgs{}) &gt; 0 &amp;&amp; isReflect {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		<span class="comment">// For reflect.makeFuncStub and reflect.methodValueCall,</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		<span class="comment">// we need to fake the stack object record.</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		<span class="comment">// These frames contain an internal/abi.RegArgs at a hard-coded offset.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		<span class="comment">// This offset matches the assembly code on amd64 and arm64.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		objs = methodValueCallFrameObjs[:]
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	} else {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		p := funcdata(f, abi.FUNCDATA_StackObjects)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		if p != nil {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			n := *(*uintptr)(p)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			p = add(p, goarch.PtrSize)
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			r0 := (*stackObjectRecord)(noescape(p))
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			objs = unsafe.Slice(r0, int(n))
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			<span class="comment">// Note: the noescape above is needed to keep</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			<span class="comment">// getStackMap from &#34;leaking param content:</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			<span class="comment">// frame&#34;.  That leak propagates up to getgcmask, then</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			<span class="comment">// GCMask, then verifyGCInfo, which converts the stack</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			<span class="comment">// gcinfo tests into heap gcinfo tests :(</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	return
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>var methodValueCallFrameObjs [1]stackObjectRecord <span class="comment">// initialized in stackobjectinit</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>func stkobjinit() {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	var abiRegArgsEface any = abi.RegArgs{}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	abiRegArgsType := efaceOf(&amp;abiRegArgsEface)._type
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	if abiRegArgsType.Kind_&amp;kindGCProg != 0 {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		throw(&#34;abiRegArgsType needs GC Prog, update methodValueCallFrameObjs&#34;)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	<span class="comment">// Set methodValueCallFrameObjs[0].gcdataoff so that</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	<span class="comment">// stackObjectRecord.gcdata() will work correctly with it.</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	ptr := uintptr(unsafe.Pointer(&amp;methodValueCallFrameObjs[0]))
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	var mod *moduledata
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	for datap := &amp;firstmoduledata; datap != nil; datap = datap.next {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		if datap.gofunc &lt;= ptr &amp;&amp; ptr &lt; datap.end {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			mod = datap
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			break
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	if mod == nil {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		throw(&#34;methodValueCallFrameObjs is not in a module&#34;)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	methodValueCallFrameObjs[0] = stackObjectRecord{
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		off:       -int32(alignUp(abiRegArgsType.Size_, 8)), <span class="comment">// It&#39;s always the highest address local.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		size:      int32(abiRegArgsType.Size_),
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		_ptrdata:  int32(abiRegArgsType.PtrBytes),
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		gcdataoff: uint32(uintptr(unsafe.Pointer(abiRegArgsType.GCData)) - mod.rodata),
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
</pre><p><a href="stkframe.go?m=text">View as plain text</a></p>

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
