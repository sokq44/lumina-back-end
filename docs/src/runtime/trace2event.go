<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/trace2event.go - Go Documentation Server</title>

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
<a href="trace2event.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">trace2event.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Trace event writing API for trace2runtime.go.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package runtime
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// Event types in the trace, args are given in square brackets.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Naming scheme:</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//   - Time range event pairs have suffixes &#34;Begin&#34; and &#34;End&#34;.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//   - &#34;Start&#34;, &#34;Stop&#34;, &#34;Create&#34;, &#34;Destroy&#34;, &#34;Block&#34;, &#34;Unblock&#34;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//     are suffixes reserved for scheduling resources.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// NOTE: If you add an event type, make sure you also update all</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// tables in this file!</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>type traceEv uint8
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>const (
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	traceEvNone traceEv = iota <span class="comment">// unused</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// Structural events.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	traceEvEventBatch <span class="comment">// start of per-M batch of events [generation, M ID, timestamp, batch length]</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	traceEvStacks     <span class="comment">// start of a section of the stack table [...traceEvStack]</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	traceEvStack      <span class="comment">// stack table entry [ID, ...{PC, func string ID, file string ID, line #}]</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	traceEvStrings    <span class="comment">// start of a section of the string dictionary [...traceEvString]</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	traceEvString     <span class="comment">// string dictionary entry [ID, length, string]</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	traceEvCPUSamples <span class="comment">// start of a section of CPU samples [...traceEvCPUSample]</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	traceEvCPUSample  <span class="comment">// CPU profiling sample [timestamp, M ID, P ID, goroutine ID, stack ID]</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	traceEvFrequency  <span class="comment">// timestamp units per sec [freq]</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// Procs.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	traceEvProcsChange <span class="comment">// current value of GOMAXPROCS [timestamp, GOMAXPROCS, stack ID]</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	traceEvProcStart   <span class="comment">// start of P [timestamp, P ID, P seq]</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	traceEvProcStop    <span class="comment">// stop of P [timestamp]</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	traceEvProcSteal   <span class="comment">// P was stolen [timestamp, P ID, P seq, M ID]</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	traceEvProcStatus  <span class="comment">// P status at the start of a generation [timestamp, P ID, status]</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// Goroutines.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	traceEvGoCreate            <span class="comment">// goroutine creation [timestamp, new goroutine ID, new stack ID, stack ID]</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	traceEvGoCreateSyscall     <span class="comment">// goroutine appears in syscall (cgo callback) [timestamp, new goroutine ID]</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	traceEvGoStart             <span class="comment">// goroutine starts running [timestamp, goroutine ID, goroutine seq]</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	traceEvGoDestroy           <span class="comment">// goroutine ends [timestamp]</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	traceEvGoDestroySyscall    <span class="comment">// goroutine ends in syscall (cgo callback) [timestamp]</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	traceEvGoStop              <span class="comment">// goroutine yields its time, but is runnable [timestamp, reason, stack ID]</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	traceEvGoBlock             <span class="comment">// goroutine blocks [timestamp, reason, stack ID]</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	traceEvGoUnblock           <span class="comment">// goroutine is unblocked [timestamp, goroutine ID, goroutine seq, stack ID]</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	traceEvGoSyscallBegin      <span class="comment">// syscall enter [timestamp, P seq, stack ID]</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	traceEvGoSyscallEnd        <span class="comment">// syscall exit [timestamp]</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	traceEvGoSyscallEndBlocked <span class="comment">// syscall exit and it blocked at some point [timestamp]</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	traceEvGoStatus            <span class="comment">// goroutine status at the start of a generation [timestamp, goroutine ID, M ID, status]</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// STW.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	traceEvSTWBegin <span class="comment">// STW start [timestamp, kind]</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	traceEvSTWEnd   <span class="comment">// STW done [timestamp]</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// GC events.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	traceEvGCActive           <span class="comment">// GC active [timestamp, seq]</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	traceEvGCBegin            <span class="comment">// GC start [timestamp, seq, stack ID]</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	traceEvGCEnd              <span class="comment">// GC done [timestamp, seq]</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	traceEvGCSweepActive      <span class="comment">// GC sweep active [timestamp, P ID]</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	traceEvGCSweepBegin       <span class="comment">// GC sweep start [timestamp, stack ID]</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	traceEvGCSweepEnd         <span class="comment">// GC sweep done [timestamp, swept bytes, reclaimed bytes]</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	traceEvGCMarkAssistActive <span class="comment">// GC mark assist active [timestamp, goroutine ID]</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	traceEvGCMarkAssistBegin  <span class="comment">// GC mark assist start [timestamp, stack ID]</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	traceEvGCMarkAssistEnd    <span class="comment">// GC mark assist done [timestamp]</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	traceEvHeapAlloc          <span class="comment">// gcController.heapLive change [timestamp, heap alloc in bytes]</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	traceEvHeapGoal           <span class="comment">// gcController.heapGoal() change [timestamp, heap goal in bytes]</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// Annotations.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	traceEvGoLabel         <span class="comment">// apply string label to current running goroutine [timestamp, label string ID]</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	traceEvUserTaskBegin   <span class="comment">// trace.NewTask [timestamp, internal task ID, internal parent task ID, name string ID, stack ID]</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	traceEvUserTaskEnd     <span class="comment">// end of a task [timestamp, internal task ID, stack ID]</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	traceEvUserRegionBegin <span class="comment">// trace.{Start,With}Region [timestamp, internal task ID, name string ID, stack ID]</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	traceEvUserRegionEnd   <span class="comment">// trace.{End,With}Region [timestamp, internal task ID, name string ID, stack ID]</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	traceEvUserLog         <span class="comment">// trace.Log [timestamp, internal task ID, key string ID, stack, value string ID]</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// traceArg is a simple wrapper type to help ensure that arguments passed</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// to traces are well-formed.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>type traceArg uint64
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// traceEventWriter is the high-level API for writing trace events.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// See the comment on traceWriter about style for more details as to why</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// this type and its methods are structured the way they are.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>type traceEventWriter struct {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	w traceWriter
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// eventWriter creates a new traceEventWriter. It is the main entrypoint for writing trace events.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// Before creating the event writer, this method will emit a status for the current goroutine</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// or proc if it exists, and if it hasn&#39;t had its status emitted yet. goStatus and procStatus indicate</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// what the status of goroutine or P should be immediately *before* the events that are about to</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// be written using the eventWriter (if they exist). No status will be written if there&#39;s no active</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// goroutine or P.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// Callers can elect to pass a constant value here if the status is clear (e.g. a goroutine must have</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// been Runnable before a GoStart). Otherwise, callers can query the status of either the goroutine</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// or P and pass the appropriate status.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// In this case, the default status should be traceGoBad or traceProcBad to help identify bugs sooner.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>func (tl traceLocker) eventWriter(goStatus traceGoStatus, procStatus traceProcStatus) traceEventWriter {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	w := tl.writer()
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	if pp := tl.mp.p.ptr(); pp != nil &amp;&amp; !pp.trace.statusWasTraced(tl.gen) &amp;&amp; pp.trace.acquireStatus(tl.gen) {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		w = w.writeProcStatus(uint64(pp.id), procStatus, pp.trace.inSweep)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	if gp := tl.mp.curg; gp != nil &amp;&amp; !gp.trace.statusWasTraced(tl.gen) &amp;&amp; gp.trace.acquireStatus(tl.gen) {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		w = w.writeGoStatus(uint64(gp.goid), int64(tl.mp.procid), goStatus, gp.inMarkAssist)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	return traceEventWriter{w}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// commit writes out a trace event and calls end. It&#39;s a helper to make the</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// common case of writing out a single event less error-prone.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func (e traceEventWriter) commit(ev traceEv, args ...traceArg) {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	e = e.write(ev, args...)
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	e.end()
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// write writes an event into the trace.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>func (e traceEventWriter) write(ev traceEv, args ...traceArg) traceEventWriter {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	e.w = e.w.event(ev, args...)
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	return e
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">// end finishes writing to the trace. The traceEventWriter must not be used after this call.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>func (e traceEventWriter) end() {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	e.w.end()
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// traceEventWrite is the part of traceEvent that actually writes the event.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>func (w traceWriter) event(ev traceEv, args ...traceArg) traceWriter {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// Make sure we have room.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	w, _ = w.ensure(1 + (len(args)+1)*traceBytesPerNumber)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// Compute the timestamp diff that we&#39;ll put in the trace.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	ts := traceClockNow()
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	if ts &lt;= w.traceBuf.lastTime {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		ts = w.traceBuf.lastTime + 1
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	tsDiff := uint64(ts - w.traceBuf.lastTime)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	w.traceBuf.lastTime = ts
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// Write out event.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	w.byte(byte(ev))
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	w.varint(tsDiff)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	for _, arg := range args {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		w.varint(uint64(arg))
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	return w
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// stack takes a stack trace skipping the provided number of frames.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// It then returns a traceArg representing that stack which may be</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// passed to write.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>func (tl traceLocker) stack(skip int) traceArg {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	return traceArg(traceStack(skip, tl.mp, tl.gen))
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// startPC takes a start PC for a goroutine and produces a unique</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">// stack ID for it.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// It then returns a traceArg representing that stack which may be</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// passed to write.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>func (tl traceLocker) startPC(pc uintptr) traceArg {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// +PCQuantum because makeTraceFrame expects return PCs and subtracts PCQuantum.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	return traceArg(trace.stackTab[tl.gen%2].put([]uintptr{
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		logicalStackSentinel,
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		startPCForTrace(pc) + sys.PCQuantum,
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	}))
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">// string returns a traceArg representing s which may be passed to write.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">// The string is assumed to be relatively short and popular, so it may be</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span><span class="comment">// stored for a while in the string dictionary.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>func (tl traceLocker) string(s string) traceArg {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	return traceArg(trace.stringTab[tl.gen%2].put(tl.gen, s))
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// uniqueString returns a traceArg representing s which may be passed to write.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// The string is assumed to be unique or long, so it will be written out to</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// the trace eagerly.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>func (tl traceLocker) uniqueString(s string) traceArg {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	return traceArg(trace.stringTab[tl.gen%2].emit(tl.gen, s))
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
</pre><p><a href="trace2event.go?m=text">View as plain text</a></p>

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
