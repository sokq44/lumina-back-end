<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/trace2stack.go - Go Documentation Server</title>

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
<a href="trace2stack.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">trace2stack.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build goexperiment.exectracer2</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Trace stack table and acquisition.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package runtime
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>const (
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// Maximum number of PCs in a single stack trace.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// Since events contain only stack id rather than whole stack trace,</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// we can allow quite large values here.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	traceStackSize = 128
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// logicalStackSentinel is a sentinel value at pcBuf[0] signifying that</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// pcBuf[1:] holds a logical stack requiring no further processing. Any other</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// value at pcBuf[0] represents a skip value to apply to the physical stack in</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// pcBuf[1:] after inline expansion.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	logicalStackSentinel = ^uintptr(0)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// traceStack captures a stack trace and registers it in the trace stack table.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// It then returns its unique ID.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// skip controls the number of leaf frames to omit in order to hide tracer internals</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// from stack traces, see CL 5523.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// Avoid calling this function directly. gen needs to be the current generation</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// that this stack trace is being written out for, which needs to be synchronized with</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// generations moving forward. Prefer traceEventWriter.stack.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>func traceStack(skip int, mp *m, gen uintptr) uint64 {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	var pcBuf [traceStackSize]uintptr
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	gp := getg()
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	curgp := gp.m.curg
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	nstk := 1
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	if tracefpunwindoff() || mp.hasCgoOnStack() {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		<span class="comment">// Slow path: Unwind using default unwinder. Used when frame pointer</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		<span class="comment">// unwinding is unavailable or disabled (tracefpunwindoff), or might</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		<span class="comment">// produce incomplete results or crashes (hasCgoOnStack). Note that no</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		<span class="comment">// cgo callback related crashes have been observed yet. The main</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		<span class="comment">// motivation is to take advantage of a potentially registered cgo</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		<span class="comment">// symbolizer.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		pcBuf[0] = logicalStackSentinel
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		if curgp == gp {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			nstk += callers(skip+1, pcBuf[1:])
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		} else if curgp != nil {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			nstk += gcallers(curgp, skip, pcBuf[1:])
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	} else {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		<span class="comment">// Fast path: Unwind using frame pointers.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		pcBuf[0] = uintptr(skip)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		if curgp == gp {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			nstk += fpTracebackPCs(unsafe.Pointer(getfp()), pcBuf[1:])
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		} else if curgp != nil {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			<span class="comment">// We&#39;re called on the g0 stack through mcall(fn) or systemstack(fn). To</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			<span class="comment">// behave like gcallers above, we start unwinding from sched.bp, which</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			<span class="comment">// points to the caller frame of the leaf frame on g&#39;s stack. The return</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			<span class="comment">// address of the leaf frame is stored in sched.pc, which we manually</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			<span class="comment">// capture here.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			pcBuf[1] = curgp.sched.pc
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			nstk += 1 + fpTracebackPCs(unsafe.Pointer(curgp.sched.bp), pcBuf[2:])
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	if nstk &gt; 0 {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		nstk-- <span class="comment">// skip runtime.goexit</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	if nstk &gt; 0 &amp;&amp; curgp.goid == 1 {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		nstk-- <span class="comment">// skip runtime.main</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	id := trace.stackTab[gen%2].put(pcBuf[:nstk])
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	return id
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// traceStackTable maps stack traces (arrays of PC&#39;s) to unique uint32 ids.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// It is lock-free for reading.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>type traceStackTable struct {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	tab traceMap
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// put returns a unique id for the stack trace pcs and caches it in the table,</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// if it sees the trace for the first time.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func (t *traceStackTable) put(pcs []uintptr) uint64 {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	if len(pcs) == 0 {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		return 0
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	id, _ := t.tab.put(noescape(unsafe.Pointer(&amp;pcs[0])), uintptr(len(pcs))*unsafe.Sizeof(uintptr(0)))
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	return id
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// dump writes all previously cached stacks to trace buffers,</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// releases all memory and resets state. It must only be called once the caller</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// can guarantee that there are no more writers to the table.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// This must run on the system stack because it flushes buffers and thus</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// may acquire trace.lock.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>func (t *traceStackTable) dump(gen uintptr) {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	w := unsafeTraceWriter(gen, nil)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// Iterate over the table.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// Do not acquire t.tab.lock. There&#39;s a conceptual lock cycle between acquiring this lock</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// here and allocation-related locks. Specifically, this lock may be acquired when an event</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// is emitted in allocation paths. Simultaneously, we might allocate here with the lock held,</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// creating a cycle. In practice, this cycle is never exercised. Because the table is only</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// dumped once there are no more writers, it&#39;s not possible for the cycle to occur. However</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// the lockrank mode is not sophisticated enough to identify this, and if it&#39;s not possible</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// for that cycle to happen, then it&#39;s also not possible for this to race with writers to</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// the table.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	for i := range t.tab.tab {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		stk := t.tab.bucket(i)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		for ; stk != nil; stk = stk.next() {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			stack := unsafe.Slice((*uintptr)(unsafe.Pointer(&amp;stk.data[0])), uintptr(len(stk.data))/unsafe.Sizeof(uintptr(0)))
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			<span class="comment">// N.B. This might allocate, but that&#39;s OK because we&#39;re not writing to the M&#39;s buffer,</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			<span class="comment">// but one we&#39;re about to create (with ensure).</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			frames := makeTraceFrames(gen, fpunwindExpand(stack))
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			<span class="comment">// Returns the maximum number of bytes required to hold the encoded stack, given that</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			<span class="comment">// it contains N frames.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			maxBytes := 1 + (2+4*len(frames))*traceBytesPerNumber
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			<span class="comment">// Estimate the size of this record. This</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			<span class="comment">// bound is pretty loose, but avoids counting</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			<span class="comment">// lots of varint sizes.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			<span class="comment">// Add 1 because we might also write traceEvStacks.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			var flushed bool
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			w, flushed = w.ensure(1 + maxBytes)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			if flushed {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>				w.byte(byte(traceEvStacks))
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			<span class="comment">// Emit stack event.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			w.byte(byte(traceEvStack))
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			w.varint(uint64(stk.id))
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			w.varint(uint64(len(frames)))
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			for _, frame := range frames {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>				w.varint(uint64(frame.PC))
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>				w.varint(frame.funcID)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>				w.varint(frame.fileID)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>				w.varint(frame.line)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// Still, hold the lock over reset. The callee expects it, even though it&#39;s</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// not strictly necessary.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	lock(&amp;t.tab.lock)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	t.tab.reset()
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	unlock(&amp;t.tab.lock)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	w.flush().end()
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">// makeTraceFrames returns the frames corresponding to pcs. It may</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// allocate and may emit trace events.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>func makeTraceFrames(gen uintptr, pcs []uintptr) []traceFrame {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	frames := make([]traceFrame, 0, len(pcs))
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	ci := CallersFrames(pcs)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	for {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		f, more := ci.Next()
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		frames = append(frames, makeTraceFrame(gen, f))
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		if !more {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			return frames
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>type traceFrame struct {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	PC     uintptr
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	funcID uint64
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	fileID uint64
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	line   uint64
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// makeTraceFrame sets up a traceFrame for a frame.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>func makeTraceFrame(gen uintptr, f Frame) traceFrame {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	var frame traceFrame
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	frame.PC = f.PC
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	fn := f.Function
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	const maxLen = 1 &lt;&lt; 10
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	if len(fn) &gt; maxLen {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		fn = fn[len(fn)-maxLen:]
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	frame.funcID = trace.stringTab[gen%2].put(gen, fn)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	frame.line = uint64(f.Line)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	file := f.File
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	if len(file) &gt; maxLen {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		file = file[len(file)-maxLen:]
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	frame.fileID = trace.stringTab[gen%2].put(gen, file)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	return frame
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// tracefpunwindoff returns true if frame pointer unwinding for the tracer is</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// disabled via GODEBUG or not supported by the architecture.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>func tracefpunwindoff() bool {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	return debug.tracefpunwindoff != 0 || (goarch.ArchFamily != goarch.AMD64 &amp;&amp; goarch.ArchFamily != goarch.ARM64)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// fpTracebackPCs populates pcBuf with the return addresses for each frame and</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// returns the number of PCs written to pcBuf. The returned PCs correspond to</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// &#34;physical frames&#34; rather than &#34;logical frames&#34;; that is if A is inlined into</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// B, this will return a PC for only B.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>func fpTracebackPCs(fp unsafe.Pointer, pcBuf []uintptr) (i int) {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	for i = 0; i &lt; len(pcBuf) &amp;&amp; fp != nil; i++ {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// return addr sits one word above the frame pointer</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		pcBuf[i] = *(*uintptr)(unsafe.Pointer(uintptr(fp) + goarch.PtrSize))
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		<span class="comment">// follow the frame pointer to the next one</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		fp = unsafe.Pointer(*(*uintptr)(fp))
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	return i
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// fpunwindExpand checks if pcBuf contains logical frames (which include inlined</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// frames) or physical frames (produced by frame pointer unwinding) using a</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// sentinel value in pcBuf[0]. Logical frames are simply returned without the</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// sentinel. Physical frames are turned into logical frames via inline unwinding</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// and by applying the skip value that&#39;s stored in pcBuf[0].</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>func fpunwindExpand(pcBuf []uintptr) []uintptr {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	if len(pcBuf) &gt; 0 &amp;&amp; pcBuf[0] == logicalStackSentinel {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		<span class="comment">// pcBuf contains logical rather than inlined frames, skip has already been</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		<span class="comment">// applied, just return it without the sentinel value in pcBuf[0].</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		return pcBuf[1:]
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	var (
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		lastFuncID = abi.FuncIDNormal
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		newPCBuf   = make([]uintptr, 0, traceStackSize)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		skip       = pcBuf[0]
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		<span class="comment">// skipOrAdd skips or appends retPC to newPCBuf and returns true if more</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		<span class="comment">// pcs can be added.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		skipOrAdd = func(retPC uintptr) bool {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			if skip &gt; 0 {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>				skip--
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			} else {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>				newPCBuf = append(newPCBuf, retPC)
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			return len(newPCBuf) &lt; cap(newPCBuf)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>outer:
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	for _, retPC := range pcBuf[1:] {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		callPC := retPC - 1
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		fi := findfunc(callPC)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		if !fi.valid() {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			<span class="comment">// There is no funcInfo if callPC belongs to a C function. In this case</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			<span class="comment">// we still keep the pc, but don&#39;t attempt to expand inlined frames.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			if more := skipOrAdd(retPC); !more {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				break outer
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			continue
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		u, uf := newInlineUnwinder(fi, callPC)
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		for ; uf.valid(); uf = u.next(uf) {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			sf := u.srcFunc(uf)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			if sf.funcID == abi.FuncIDWrapper &amp;&amp; elideWrapperCalling(lastFuncID) {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>				<span class="comment">// ignore wrappers</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			} else if more := skipOrAdd(uf.pc + 1); !more {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>				break outer
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			lastFuncID = sf.funcID
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	return newPCBuf
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// startPCForTrace returns the start PC of a goroutine for tracing purposes.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">// If pc is a wrapper, it returns the PC of the wrapped function. Otherwise it</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">// returns pc.</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>func startPCForTrace(pc uintptr) uintptr {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	f := findfunc(pc)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	if !f.valid() {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		return pc <span class="comment">// may happen for locked g in extra M since its pc is 0.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	w := funcdata(f, abi.FUNCDATA_WrapInfo)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	if w == nil {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		return pc <span class="comment">// not a wrapper</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	return f.datap.textAddr(*(*uint32)(w))
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
</pre><p><a href="trace2stack.go?m=text">View as plain text</a></p>

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
