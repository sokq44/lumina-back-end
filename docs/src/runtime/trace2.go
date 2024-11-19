<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/trace2.go - Go Documentation Server</title>

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
<a href="trace2.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">trace2.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Go execution tracer.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// The tracer captures a wide range of execution events like goroutine</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// creation/blocking/unblocking, syscall enter/exit/block, GC-related events,</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// changes of heap size, processor start/stop, etc and writes them to a buffer</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// in a compact form. A precise nanosecond-precision timestamp and a stack</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// trace is captured for most events.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// Tracer invariants (to keep the synchronization making sense):</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// - An m that has a trace buffer must be on either the allm or sched.freem lists.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// - Any trace buffer mutation must either be happening in traceAdvance or between</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//   a traceAcquire and a subsequent traceRelease.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// - traceAdvance cannot return until the previous generation&#39;s buffers are all flushed.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// See https://go.dev/issue/60773 for a link to the full design.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>package runtime
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>import (
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// Trace state.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// trace is global tracing context.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>var trace struct {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// trace.lock must only be acquired on the system stack where</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// stack splits cannot happen while it is held.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	lock mutex
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// Trace buffer management.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// First we check the empty list for any free buffers. If not, buffers</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// are allocated directly from the OS. Once they&#39;re filled up and/or</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// flushed, they end up on the full queue for trace.gen%2.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// The trace reader takes buffers off the full list one-by-one and</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// places them into reading until they&#39;re finished being read from.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// Then they&#39;re placed onto the empty list.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// Protected by trace.lock.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	reading       *traceBuf <span class="comment">// buffer currently handed off to user</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	empty         *traceBuf <span class="comment">// stack of empty buffers</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	full          [2]traceBufQueue
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	workAvailable atomic.Bool
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// State for the trace reader goroutine.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// Protected by trace.lock.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	readerGen     atomic.Uintptr <span class="comment">// the generation the reader is currently reading for</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	flushedGen    atomic.Uintptr <span class="comment">// the last completed generation</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	headerWritten bool           <span class="comment">// whether ReadTrace has emitted trace header</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// doneSema is used to synchronize the reader and traceAdvance. Specifically,</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// it notifies traceAdvance that the reader is done with a generation.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// Both semaphores are 0 by default (so, acquires block). traceAdvance</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// attempts to acquire for gen%2 after flushing the last buffers for gen.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// Meanwhile the reader releases the sema for gen%2 when it has finished</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// processing gen.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	doneSema [2]uint32
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// Trace data tables for deduplicating data going into the trace.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// There are 2 of each: one for gen%2, one for 1-gen%2.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	stackTab  [2]traceStackTable  <span class="comment">// maps stack traces to unique ids</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	stringTab [2]traceStringTable <span class="comment">// maps strings to unique ids</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// cpuLogRead accepts CPU profile samples from the signal handler where</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// they&#39;re generated. There are two profBufs here: one for gen%2, one for</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// 1-gen%2. These profBufs use a three-word header to hold the IDs of the P, G,</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// and M (respectively) that were active at the time of the sample. Because</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// profBuf uses a record with all zeros in its header to indicate overflow,</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// we make sure to make the P field always non-zero: The ID of a real P will</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// start at bit 1, and bit 0 will be set. Samples that arrive while no P is</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// running (such as near syscalls) will set the first header field to 0b10.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// This careful handling of the first header field allows us to store ID of</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// the active G directly in the second field, even though that will be 0</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// when sampling g0.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// Initialization and teardown of these fields is protected by traceAdvanceSema.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	cpuLogRead  [2]*profBuf
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	signalLock  atomic.Uint32              <span class="comment">// protects use of the following member, only usable in signal handlers</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	cpuLogWrite [2]atomic.Pointer[profBuf] <span class="comment">// copy of cpuLogRead for use in signal handlers, set without signalLock</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	cpuSleep    *wakeableSleep
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	cpuLogDone  &lt;-chan struct{}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	cpuBuf      [2]*traceBuf
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	reader atomic.Pointer[g] <span class="comment">// goroutine that called ReadTrace, or nil</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// Fast mappings from enumerations to string IDs that are prepopulated</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// in the trace.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	markWorkerLabels [2][len(gcMarkWorkerModeStrings)]traceArg
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	goStopReasons    [2][len(traceGoStopReasonStrings)]traceArg
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	goBlockReasons   [2][len(traceBlockReasonStrings)]traceArg
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// Trace generation counter.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	gen            atomic.Uintptr
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	lastNonZeroGen uintptr <span class="comment">// last non-zero value of gen</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// shutdown is set when we are waiting for trace reader to finish after setting gen to 0</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// Writes protected by trace.lock.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	shutdown atomic.Bool
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// Number of goroutines in syscall exiting slow path.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	exitingSyscall atomic.Int32
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// seqGC is the sequence counter for GC begin/end.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// Mutated only during stop-the-world.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	seqGC uint64
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// Trace public API.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>var (
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	traceAdvanceSema  uint32 = 1
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	traceShutdownSema uint32 = 1
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>)
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// StartTrace enables tracing for the current process.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// While tracing, the data will be buffered and available via [ReadTrace].</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// StartTrace returns an error if tracing is already enabled.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// Most clients should use the [runtime/trace] package or the [testing] package&#39;s</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// -test.trace flag instead of calling StartTrace directly.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func StartTrace() error {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	if traceEnabled() || traceShuttingDown() {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		return errorString(&#34;tracing is already enabled&#34;)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// Block until cleanup of the last trace is done.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	semacquire(&amp;traceShutdownSema)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	semrelease(&amp;traceShutdownSema)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// Hold traceAdvanceSema across trace start, since we&#39;ll want it on</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// the other side of tracing being enabled globally.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	semacquire(&amp;traceAdvanceSema)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">// Initialize CPU profile -&gt; trace ingestion.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	traceInitReadCPU()
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// Compute the first generation for this StartTrace.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// Note: we start from the last non-zero generation rather than 1 so we</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// can avoid resetting all the arrays indexed by gen%2 or gen%3. There&#39;s</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// more than one of each per m, p, and goroutine.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	firstGen := traceNextGen(trace.lastNonZeroGen)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// Reset GC sequencer.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	trace.seqGC = 1
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// Reset trace reader state.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	trace.headerWritten = false
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	trace.readerGen.Store(firstGen)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	trace.flushedGen.Store(0)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// Register some basic strings in the string tables.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	traceRegisterLabelsAndReasons(firstGen)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// Stop the world.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// The purpose of stopping the world is to make sure that no goroutine is in a</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// context where it could emit an event by bringing all goroutines to a safe point</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// with no opportunity to transition.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// The exception to this rule are goroutines that are concurrently exiting a syscall.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// Those will all be forced into the syscalling slow path, and we&#39;ll just make sure</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// that we don&#39;t observe any goroutines in that critical section before starting</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// the world again.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// A good follow-up question to this is why stopping the world is necessary at all</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// given that we have traceAcquire and traceRelease. Unfortunately, those only help</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// us when tracing is already active (for performance, so when tracing is off the</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// tracing seqlock is left untouched). The main issue here is subtle: we&#39;re going to</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// want to obtain a correct starting status for each goroutine, but there are windows</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// of time in which we could read and emit an incorrect status. Specifically:</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">//	trace := traceAcquire()</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">//  // &lt;----&gt; problem window</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">//	casgstatus(gp, _Gwaiting, _Grunnable)</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">//	if trace.ok() {</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">//		trace.GoUnpark(gp, 2)</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">//		traceRelease(trace)</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">//	}</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// More precisely, if we readgstatus for a gp while another goroutine is in the problem</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// window and that goroutine didn&#39;t observe that tracing had begun, then we might write</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// a GoStatus(GoWaiting) event for that goroutine, but it won&#39;t trace an event marking</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// the transition from GoWaiting to GoRunnable. The trace will then be broken, because</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// future events will be emitted assuming the tracer sees GoRunnable.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// In short, what we really need here is to make sure that the next time *any goroutine*</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// hits a traceAcquire, it sees that the trace is enabled.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// Note also that stopping the world is necessary to make sure sweep-related events are</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// coherent. Since the world is stopped and sweeps are non-preemptible, we can never start</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// the world and see an unpaired sweep &#39;end&#39; event. Other parts of the tracer rely on this.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	stw := stopTheWorld(stwStartTrace)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// Prevent sysmon from running any code that could generate events.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	lock(&amp;sched.sysmonlock)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// Reset mSyscallID on all Ps while we have them stationary and the trace is disabled.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	for _, pp := range allp {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		pp.trace.mSyscallID = -1
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// Start tracing.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// After this executes, other Ms may start creating trace buffers and emitting</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// data into them.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	trace.gen.Store(firstGen)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	<span class="comment">// Wait for exitingSyscall to drain.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// It may not monotonically decrease to zero, but in the limit it will always become</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// zero because the world is stopped and there are no available Ps for syscall-exited</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">// goroutines to run on.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">// Because we set gen before checking this, and because exitingSyscall is always incremented</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	<span class="comment">// *after* traceAcquire (which checks gen), we can be certain that when exitingSyscall is zero</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// that any goroutine that goes to exit a syscall from then on *must* observe the new gen.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// The critical section on each goroutine here is going to be quite short, so the likelihood</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// that we observe a zero value is high.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	for trace.exitingSyscall.Load() != 0 {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		osyield()
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	<span class="comment">// Record some initial pieces of information.</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// N.B. This will also emit a status event for this goroutine.</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	tl := traceAcquire()
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	tl.Gomaxprocs(gomaxprocs)  <span class="comment">// Get this as early in the trace as possible. See comment in traceAdvance.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	tl.STWStart(stwStartTrace) <span class="comment">// We didn&#39;t trace this above, so trace it now.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// Record the fact that a GC is active, if applicable.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if gcphase == _GCmark || gcphase == _GCmarktermination {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		tl.GCActive()
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">// Record the heap goal so we have it at the very beginning of the trace.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	tl.HeapGoal()
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	<span class="comment">// Make sure a ProcStatus is emitted for every P, while we&#39;re here.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	for _, pp := range allp {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		tl.writer().writeProcStatusForP(pp, pp == tl.mp.p.ptr()).end()
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	traceRelease(tl)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	unlock(&amp;sched.sysmonlock)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	startTheWorld(stw)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	traceStartReadCPU()
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	traceAdvancer.start()
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	semrelease(&amp;traceAdvanceSema)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	return nil
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">// StopTrace stops tracing, if it was previously enabled.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">// StopTrace only returns after all the reads for the trace have completed.</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>func StopTrace() {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	traceAdvance(true)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span><span class="comment">// traceAdvance moves tracing to the next generation, and cleans up the current generation,</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span><span class="comment">// ensuring that it&#39;s flushed out before returning. If stopTrace is true, it disables tracing</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span><span class="comment">// altogether instead of advancing to the next generation.</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span><span class="comment">// traceAdvanceSema must not be held.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>func traceAdvance(stopTrace bool) {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	semacquire(&amp;traceAdvanceSema)
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	<span class="comment">// Get the gen that we&#39;re advancing from. In this function we don&#39;t really care much</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	<span class="comment">// about the generation we&#39;re advancing _into_ since we&#39;ll do all the cleanup in this</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	<span class="comment">// generation for the next advancement.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	gen := trace.gen.Load()
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	if gen == 0 {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		<span class="comment">// We may end up here traceAdvance is called concurrently with StopTrace.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		semrelease(&amp;traceAdvanceSema)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		return
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	<span class="comment">// Write an EvFrequency event for this generation.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// N.B. This may block for quite a while to get a good frequency estimate, so make sure we do</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// this here and not e.g. on the trace reader.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	traceFrequency(gen)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// Collect all the untraced Gs.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	type untracedG struct {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		gp           *g
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		goid         uint64
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		mid          int64
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		status       uint32
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		waitreason   waitReason
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		inMarkAssist bool
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	var untracedGs []untracedG
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	forEachGRace(func(gp *g) {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		<span class="comment">// Make absolutely sure all Gs are ready for the next</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		<span class="comment">// generation. We need to do this even for dead Gs because</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		<span class="comment">// they may come alive with a new identity, and its status</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		<span class="comment">// traced bookkeeping might end up being stale.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		<span class="comment">// We may miss totally new goroutines, but they&#39;ll always</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		<span class="comment">// have clean bookkeeping.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		gp.trace.readyNextGen(gen)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		<span class="comment">// If the status was traced, nothing else to do.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		if gp.trace.statusWasTraced(gen) {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>			return
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		<span class="comment">// Scribble down information about this goroutine.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		ug := untracedG{gp: gp, mid: -1}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			me := getg().m.curg
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			<span class="comment">// We don&#39;t have to handle this G status transition because we</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>			<span class="comment">// already eliminated ourselves from consideration above.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>			casGToWaiting(me, _Grunning, waitReasonTraceGoroutineStatus)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			<span class="comment">// We need to suspend and take ownership of the G to safely read its</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			<span class="comment">// goid. Note that we can&#39;t actually emit the event at this point</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			<span class="comment">// because we might stop the G in a window where it&#39;s unsafe to write</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			<span class="comment">// events based on the G&#39;s status. We need the global trace buffer flush</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			<span class="comment">// coming up to make sure we&#39;re not racing with the G.</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>			<span class="comment">// It should be very unlikely that we try to preempt a running G here.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			<span class="comment">// The only situation that we might is that we&#39;re racing with a G</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			<span class="comment">// that&#39;s running for the first time in this generation. Therefore,</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>			<span class="comment">// this should be relatively fast.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>			s := suspendG(gp)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>			if !s.dead {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>				ug.goid = s.g.goid
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>				if s.g.m != nil {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>					ug.mid = int64(s.g.m.procid)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>				}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>				ug.status = readgstatus(s.g) &amp;^ _Gscan
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>				ug.waitreason = s.g.waitreason
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>				ug.inMarkAssist = s.g.inMarkAssist
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			resumeG(s)
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			casgstatus(me, _Gwaiting, _Grunning)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		})
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		if ug.goid != 0 {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>			untracedGs = append(untracedGs, ug)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	})
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	if !stopTrace {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		<span class="comment">// Re-register runtime goroutine labels and stop/block reasons.</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		traceRegisterLabelsAndReasons(traceNextGen(gen))
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	}
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	<span class="comment">// Now that we&#39;ve done some of the heavy stuff, prevent the world from stopping.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	<span class="comment">// This is necessary to ensure the consistency of the STW events. If we&#39;re feeling</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	<span class="comment">// adventurous we could lift this restriction and add a STWActive event, but the</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	<span class="comment">// cost of maintaining this consistency is low. We&#39;re not going to hold this semaphore</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	<span class="comment">// for very long and most STW periods are very short.</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	<span class="comment">// Once we hold worldsema, prevent preemption as well so we&#39;re not interrupted partway</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	<span class="comment">// through this. We want to get this done as soon as possible.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	semacquire(&amp;worldsema)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	<span class="comment">// Advance the generation or stop the trace.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	trace.lastNonZeroGen = gen
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	if stopTrace {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>			<span class="comment">// Ordering is important here. Set shutdown first, then disable tracing,</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>			<span class="comment">// so that conditions like (traceEnabled() || traceShuttingDown()) have</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			<span class="comment">// no opportunity to be false. Hold the trace lock so this update appears</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			<span class="comment">// atomic to the trace reader.</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			lock(&amp;trace.lock)
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			trace.shutdown.Store(true)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			trace.gen.Store(0)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			unlock(&amp;trace.lock)
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		})
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	} else {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		trace.gen.Store(traceNextGen(gen))
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// Emit a ProcsChange event so we have one on record for each generation.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	<span class="comment">// Let&#39;s emit it as soon as possible so that downstream tools can rely on the value</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	<span class="comment">// being there fairly soon in a generation.</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s important that we do this before allowing stop-the-worlds again,</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	<span class="comment">// because the procs count could change.</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	if !stopTrace {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		tl := traceAcquire()
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		tl.Gomaxprocs(gomaxprocs)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		traceRelease(tl)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	<span class="comment">// Emit a GCActive event in the new generation if necessary.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s important that we do this before allowing stop-the-worlds again,</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	<span class="comment">// because that could emit global GC-related events.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	if !stopTrace &amp;&amp; (gcphase == _GCmark || gcphase == _GCmarktermination) {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		tl := traceAcquire()
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		tl.GCActive()
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		traceRelease(tl)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	<span class="comment">// Preemption is OK again after this. If the world stops or whatever it&#39;s fine.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	<span class="comment">// We&#39;re just cleaning up the last generation after this point.</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	<span class="comment">// We also don&#39;t care if the GC starts again after this for the same reasons.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	releasem(mp)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	semrelease(&amp;worldsema)
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	<span class="comment">// Snapshot allm and freem.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	<span class="comment">// Snapshotting after the generation counter update is sufficient.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	<span class="comment">// Because an m must be on either allm or sched.freem if it has an active trace</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	<span class="comment">// buffer, new threads added to allm after this point must necessarily observe</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	<span class="comment">// the new generation number (sched.lock acts as a barrier).</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	<span class="comment">// Threads that exit before this point and are on neither list explicitly</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	<span class="comment">// flush their own buffers in traceThreadDestroy.</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	<span class="comment">// Snapshotting freem is necessary because Ms can continue to emit events</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	<span class="comment">// while they&#39;re still on that list. Removal from sched.freem is serialized with</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	<span class="comment">// this snapshot, so either we&#39;ll capture an m on sched.freem and race with</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	<span class="comment">// the removal to flush its buffers (resolved by traceThreadDestroy acquiring</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	<span class="comment">// the thread&#39;s seqlock, which one of us must win, so at least its old gen buffer</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	<span class="comment">// will be flushed in time for the new generation) or it will have flushed its</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	<span class="comment">// buffers before we snapshotted it to begin with.</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	lock(&amp;sched.lock)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	mToFlush := allm
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	for mp := mToFlush; mp != nil; mp = mp.alllink {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		mp.trace.link = mp.alllink
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	for mp := sched.freem; mp != nil; mp = mp.freelink {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		mp.trace.link = mToFlush
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		mToFlush = mp
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	unlock(&amp;sched.lock)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// Iterate over our snapshot, flushing every buffer until we&#39;re done.</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	<span class="comment">// Because trace writers read the generation while the seqlock is</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	<span class="comment">// held, we can be certain that when there are no writers there are</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	<span class="comment">// also no stale generation values left. Therefore, it&#39;s safe to flush</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	<span class="comment">// any buffers that remain in that generation&#39;s slot.</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	const debugDeadlock = false
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		<span class="comment">// Track iterations for some rudimentary deadlock detection.</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		i := 0
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		detectedDeadlock := false
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		for mToFlush != nil {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>			prev := &amp;mToFlush
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			for mp := *prev; mp != nil; {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>				if mp.trace.seqlock.Load()%2 != 0 {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>					<span class="comment">// The M is writing. Come back to it later.</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>					prev = &amp;mp.trace.link
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>					mp = mp.trace.link
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>					continue
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>				}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>				<span class="comment">// Flush the trace buffer.</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>				<span class="comment">// trace.lock needed for traceBufFlush, but also to synchronize</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>				<span class="comment">// with traceThreadDestroy, which flushes both buffers unconditionally.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>				lock(&amp;trace.lock)
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>				bufp := &amp;mp.trace.buf[gen%2]
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>				if *bufp != nil {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>					traceBufFlush(*bufp, gen)
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>					*bufp = nil
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>				}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>				unlock(&amp;trace.lock)
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>				<span class="comment">// Remove the m from the flush list.</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>				*prev = mp.trace.link
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>				mp.trace.link = nil
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>				mp = *prev
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			<span class="comment">// Yield only if we&#39;re going to be going around the loop again.</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			if mToFlush != nil {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>				osyield()
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			if debugDeadlock {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>				<span class="comment">// Try to detect a deadlock. We probably shouldn&#39;t loop here</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>				<span class="comment">// this many times.</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>				if i &gt; 100000 &amp;&amp; !detectedDeadlock {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>					detectedDeadlock = true
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>					println(&#34;runtime: failing to flush&#34;)
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>					for mp := mToFlush; mp != nil; mp = mp.trace.link {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>						print(&#34;runtime: m=&#34;, mp.id, &#34;\n&#34;)
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>					}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>				}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>				i++
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	})
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	<span class="comment">// At this point, the old generation is fully flushed minus stack and string</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	<span class="comment">// tables, CPU samples, and goroutines that haven&#39;t run at all during the last</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	<span class="comment">// generation.</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	<span class="comment">// Check to see if any Gs still haven&#39;t had events written out for them.</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	statusWriter := unsafeTraceWriter(gen, nil)
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	for _, ug := range untracedGs {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		if ug.gp.trace.statusWasTraced(gen) {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>			<span class="comment">// It was traced, we don&#39;t need to do anything.</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			continue
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		<span class="comment">// It still wasn&#39;t traced. Because we ensured all Ms stopped writing trace</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		<span class="comment">// events to the last generation, that must mean the G never had its status</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		<span class="comment">// traced in gen between when we recorded it and now. If that&#39;s true, the goid</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		<span class="comment">// and status we recorded then is exactly what we want right now.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		status := goStatusToTraceGoStatus(ug.status, ug.waitreason)
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		statusWriter = statusWriter.writeGoStatus(ug.goid, ug.mid, status, ug.inMarkAssist)
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	statusWriter.flush().end()
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	<span class="comment">// Read everything out of the last gen&#39;s CPU profile buffer.</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	traceReadCPU(gen)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		<span class="comment">// Flush CPU samples, stacks, and strings for the last generation. This is safe,</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		<span class="comment">// because we&#39;re now certain no M is writing to the last generation.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		<span class="comment">// Ordering is important here. traceCPUFlush may generate new stacks and dumping</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		<span class="comment">// stacks may generate new strings.</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		traceCPUFlush(gen)
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		trace.stackTab[gen%2].dump(gen)
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		trace.stringTab[gen%2].reset(gen)
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		<span class="comment">// That&#39;s it. This generation is done producing buffers.</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		lock(&amp;trace.lock)
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		trace.flushedGen.Store(gen)
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		unlock(&amp;trace.lock)
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	})
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	if stopTrace {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		semacquire(&amp;traceShutdownSema)
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		<span class="comment">// Finish off CPU profile reading.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		traceStopReadCPU()
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	} else {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		<span class="comment">// Go over each P and emit a status event for it if necessary.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		<span class="comment">// We do this at the beginning of the new generation instead of the</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		<span class="comment">// end like we do for goroutines because forEachP doesn&#39;t give us a</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		<span class="comment">// hook to skip Ps that have already been traced. Since we have to</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		<span class="comment">// preempt all Ps anyway, might as well stay consistent with StartTrace</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		<span class="comment">// which does this during the STW.</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		semacquire(&amp;worldsema)
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		forEachP(waitReasonTraceProcStatus, func(pp *p) {
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>			tl := traceAcquire()
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>			if !pp.trace.statusWasTraced(tl.gen) {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>				tl.writer().writeProcStatusForP(pp, false).end()
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>			}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>			traceRelease(tl)
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		})
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		<span class="comment">// Perform status reset on dead Ps because they just appear as idle.</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		<span class="comment">// Holding worldsema prevents allp from changing.</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		<span class="comment">// TODO(mknyszek): Consider explicitly emitting ProcCreate and ProcDestroy</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		<span class="comment">// events to indicate whether a P exists, rather than just making its</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		<span class="comment">// existence implicit.</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		for _, pp := range allp[len(allp):cap(allp)] {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>			pp.trace.readyNextGen(traceNextGen(gen))
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		semrelease(&amp;worldsema)
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	<span class="comment">// Block until the trace reader has finished processing the last generation.</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	semacquire(&amp;trace.doneSema[gen%2])
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	if raceenabled {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		raceacquire(unsafe.Pointer(&amp;trace.doneSema[gen%2]))
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	}
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	<span class="comment">// Double-check that things look as we expect after advancing and perform some</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	<span class="comment">// final cleanup if the trace has fully stopped.</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		lock(&amp;trace.lock)
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		if !trace.full[gen%2].empty() {
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>			throw(&#34;trace: non-empty full trace buffer for done generation&#34;)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		if stopTrace {
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>			if !trace.full[1-(gen%2)].empty() {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>				throw(&#34;trace: non-empty full trace buffer for next generation&#34;)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>			}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>			if trace.reading != nil || trace.reader.Load() != nil {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>				throw(&#34;trace: reading after shutdown&#34;)
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>			}
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>			<span class="comment">// Free all the empty buffers.</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>			for trace.empty != nil {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>				buf := trace.empty
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>				trace.empty = buf.link
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>				sysFree(unsafe.Pointer(buf), unsafe.Sizeof(*buf), &amp;memstats.other_sys)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>			}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>			<span class="comment">// Clear trace.shutdown and other flags.</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>			trace.headerWritten = false
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>			trace.shutdown.Store(false)
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		unlock(&amp;trace.lock)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	})
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	if stopTrace {
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		<span class="comment">// Clear the sweep state on every P for the next time tracing is enabled.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		<span class="comment">// It may be stale in the next trace because we may have ended tracing in</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		<span class="comment">// the middle of a sweep on a P.</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		<span class="comment">// It&#39;s fine not to call forEachP here because tracing is disabled and we</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		<span class="comment">// know at this point that nothing is calling into the tracer, but we do</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		<span class="comment">// need to look at dead Ps too just because GOMAXPROCS could have been called</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		<span class="comment">// at any point since we stopped tracing, and we have to ensure there&#39;s no</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		<span class="comment">// bad state on dead Ps too. Prevent a STW and a concurrent GOMAXPROCS that</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		<span class="comment">// might mutate allp by making ourselves briefly non-preemptible.</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		mp := acquirem()
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		for _, pp := range allp[:cap(allp)] {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>			pp.trace.inSweep = false
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			pp.trace.maySweep = false
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			pp.trace.swept = 0
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>			pp.trace.reclaimed = 0
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		releasem(mp)
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	}
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	<span class="comment">// Release the advance semaphore. If stopTrace is true we&#39;re still holding onto</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	<span class="comment">// traceShutdownSema.</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	<span class="comment">// Do a direct handoff. Don&#39;t let one caller of traceAdvance starve</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	<span class="comment">// other calls to traceAdvance.</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	semrelease1(&amp;traceAdvanceSema, true, 0)
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	if stopTrace {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		<span class="comment">// Stop the traceAdvancer. We can&#39;t be holding traceAdvanceSema here because</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		<span class="comment">// we&#39;ll deadlock (we&#39;re blocked on the advancer goroutine exiting, but it</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		<span class="comment">// may be currently trying to acquire traceAdvanceSema).</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		traceAdvancer.stop()
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		semrelease(&amp;traceShutdownSema)
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	}
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>func traceNextGen(gen uintptr) uintptr {
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	if gen == ^uintptr(0) {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		<span class="comment">// gen is used both %2 and %3 and we want both patterns to continue when we loop around.</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		<span class="comment">// ^uint32(0) and ^uint64(0) are both odd and multiples of 3. Therefore the next generation</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		<span class="comment">// we want is even and one more than a multiple of 3. The smallest such number is 4.</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		return 4
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	}
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	return gen + 1
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span><span class="comment">// traceRegisterLabelsAndReasons re-registers mark worker labels and</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span><span class="comment">// goroutine stop/block reasons in the string table for the provided</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span><span class="comment">// generation. Note: the provided generation must not have started yet.</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>func traceRegisterLabelsAndReasons(gen uintptr) {
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	for i, label := range gcMarkWorkerModeStrings[:] {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		trace.markWorkerLabels[gen%2][i] = traceArg(trace.stringTab[gen%2].put(gen, label))
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	}
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	for i, str := range traceBlockReasonStrings[:] {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		trace.goBlockReasons[gen%2][i] = traceArg(trace.stringTab[gen%2].put(gen, str))
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	}
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	for i, str := range traceGoStopReasonStrings[:] {
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		trace.goStopReasons[gen%2][i] = traceArg(trace.stringTab[gen%2].put(gen, str))
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	}
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span><span class="comment">// ReadTrace returns the next chunk of binary tracing data, blocking until data</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span><span class="comment">// is available. If tracing is turned off and all the data accumulated while it</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span><span class="comment">// was on has been returned, ReadTrace returns nil. The caller must copy the</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span><span class="comment">// returned data before calling ReadTrace again.</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span><span class="comment">// ReadTrace must be called from one goroutine at a time.</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>func ReadTrace() []byte {
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>top:
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	var buf []byte
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	var park bool
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		buf, park = readTrace0()
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	})
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	if park {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		gopark(func(gp *g, _ unsafe.Pointer) bool {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			if !trace.reader.CompareAndSwapNoWB(nil, gp) {
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>				<span class="comment">// We&#39;re racing with another reader.</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>				<span class="comment">// Wake up and handle this case.</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>				return false
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>			}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>			if g2 := traceReader(); gp == g2 {
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>				<span class="comment">// New data arrived between unlocking</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>				<span class="comment">// and the CAS and we won the wake-up</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>				<span class="comment">// race, so wake up directly.</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>				return false
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>			} else if g2 != nil {
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>				printlock()
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>				println(&#34;runtime: got trace reader&#34;, g2, g2.goid)
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>				throw(&#34;unexpected trace reader&#34;)
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>			}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			return true
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>		}, nil, waitReasonTraceReaderBlocked, traceBlockSystemGoroutine, 2)
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		goto top
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	}
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	return buf
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>}
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span><span class="comment">// readTrace0 is ReadTrace&#39;s continuation on g0. This must run on the</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span><span class="comment">// system stack because it acquires trace.lock.</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>func readTrace0() (buf []byte, park bool) {
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	if raceenabled {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		<span class="comment">// g0 doesn&#39;t have a race context. Borrow the user G&#39;s.</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		if getg().racectx != 0 {
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>			throw(&#34;expected racectx == 0&#34;)
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		}
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		getg().racectx = getg().m.curg.racectx
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>		<span class="comment">// (This defer should get open-coded, which is safe on</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>		<span class="comment">// the system stack.)</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		defer func() { getg().racectx = 0 }()
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	}
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	<span class="comment">// This function must not allocate while holding trace.lock:</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	<span class="comment">// allocation can call heap allocate, which will try to emit a trace</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	<span class="comment">// event while holding heap lock.</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	lock(&amp;trace.lock)
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	if trace.reader.Load() != nil {
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		<span class="comment">// More than one goroutine reads trace. This is bad.</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		<span class="comment">// But we rather do not crash the program because of tracing,</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>		<span class="comment">// because tracing can be enabled at runtime on prod servers.</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		unlock(&amp;trace.lock)
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>		println(&#34;runtime: ReadTrace called from multiple goroutines simultaneously&#34;)
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		return nil, false
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	}
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	<span class="comment">// Recycle the old buffer.</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	if buf := trace.reading; buf != nil {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		buf.link = trace.empty
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		trace.empty = buf
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>		trace.reading = nil
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	}
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	<span class="comment">// Write trace header.</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	if !trace.headerWritten {
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		trace.headerWritten = true
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		unlock(&amp;trace.lock)
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>		return []byte(&#34;go 1.22 trace\x00\x00\x00&#34;), false
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	}
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	<span class="comment">// Read the next buffer.</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	if trace.readerGen.Load() == 0 {
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>		trace.readerGen.Store(1)
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	}
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	var gen uintptr
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	for {
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		assertLockHeld(&amp;trace.lock)
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>		gen = trace.readerGen.Load()
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		<span class="comment">// Check to see if we need to block for more data in this generation</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		<span class="comment">// or if we need to move our generation forward.</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		if !trace.full[gen%2].empty() {
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>			break
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		}
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>		<span class="comment">// Most of the time readerGen is one generation ahead of flushedGen, as the</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		<span class="comment">// current generation is being read from. Then, once the last buffer is flushed</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		<span class="comment">// into readerGen, flushedGen will rise to meet it. At this point, the tracer</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>		<span class="comment">// is waiting on the reader to finish flushing the last generation so that it</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>		<span class="comment">// can continue to advance.</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		if trace.flushedGen.Load() == gen {
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>			if trace.shutdown.Load() {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>				unlock(&amp;trace.lock)
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>				<span class="comment">// Wake up anyone waiting for us to be done with this generation.</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>				<span class="comment">// Do this after reading trace.shutdown, because the thread we&#39;re</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>				<span class="comment">// waking up is going to clear trace.shutdown.</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>				if raceenabled {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>					<span class="comment">// Model synchronization on trace.doneSema, which te race</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>					<span class="comment">// detector does not see. This is required to avoid false</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>					<span class="comment">// race reports on writer passed to trace.Start.</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>					racerelease(unsafe.Pointer(&amp;trace.doneSema[gen%2]))
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>				}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>				semrelease(&amp;trace.doneSema[gen%2])
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>				<span class="comment">// We&#39;re shutting down, and the last generation is fully</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>				<span class="comment">// read. We&#39;re done.</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>				return nil, false
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>			}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>			<span class="comment">// The previous gen has had all of its buffers flushed, and</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>			<span class="comment">// there&#39;s nothing else for us to read. Advance the generation</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>			<span class="comment">// we&#39;re reading from and try again.</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>			trace.readerGen.Store(trace.gen.Load())
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>			unlock(&amp;trace.lock)
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>			<span class="comment">// Wake up anyone waiting for us to be done with this generation.</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>			<span class="comment">// Do this after reading gen to make sure we can&#39;t have the trace</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>			<span class="comment">// advance until we&#39;ve read it.</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>			if raceenabled {
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>				<span class="comment">// See comment above in the shutdown case.</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>				racerelease(unsafe.Pointer(&amp;trace.doneSema[gen%2]))
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>			}
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>			semrelease(&amp;trace.doneSema[gen%2])
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>			<span class="comment">// Reacquire the lock and go back to the top of the loop.</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>			lock(&amp;trace.lock)
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>			continue
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		}
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>		<span class="comment">// Wait for new data.</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		<span class="comment">// We don&#39;t simply use a note because the scheduler</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>		<span class="comment">// executes this goroutine directly when it wakes up</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		<span class="comment">// (also a note would consume an M).</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>		<span class="comment">// Before we drop the lock, clear the workAvailable flag. Work can</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		<span class="comment">// only be queued with trace.lock held, so this is at least true until</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>		<span class="comment">// we drop the lock.</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		trace.workAvailable.Store(false)
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		unlock(&amp;trace.lock)
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>		return nil, true
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	}
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	<span class="comment">// Pull a buffer.</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	tbuf := trace.full[gen%2].pop()
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	trace.reading = tbuf
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	unlock(&amp;trace.lock)
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	return tbuf.arr[:tbuf.pos], false
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>}
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span><span class="comment">// traceReader returns the trace reader that should be woken up, if any.</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span><span class="comment">// Callers should first check (traceEnabled() || traceShuttingDown()).</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span><span class="comment">// This must run on the system stack because it acquires trace.lock.</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>func traceReader() *g {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>	gp := traceReaderAvailable()
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	if gp == nil || !trace.reader.CompareAndSwapNoWB(gp, nil) {
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		return nil
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	}
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	return gp
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>}
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>
<span id="L843" class="ln">   843&nbsp;&nbsp;</span><span class="comment">// traceReaderAvailable returns the trace reader if it is not currently</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span><span class="comment">// scheduled and should be. Callers should first check that</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span><span class="comment">// (traceEnabled() || traceShuttingDown()) is true.</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>func traceReaderAvailable() *g {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	<span class="comment">// There are three conditions under which we definitely want to schedule</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	<span class="comment">// the reader:</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	<span class="comment">// - The reader is lagging behind in finishing off the last generation.</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	<span class="comment">//   In this case, trace buffers could even be empty, but the trace</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	<span class="comment">//   advancer will be waiting on the reader, so we have to make sure</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	<span class="comment">//   to schedule the reader ASAP.</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	<span class="comment">// - The reader has pending work to process for it&#39;s reader generation</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	<span class="comment">//   (assuming readerGen is not lagging behind). Note that we also want</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	<span class="comment">//   to be careful *not* to schedule the reader if there&#39;s no work to do.</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	<span class="comment">// - The trace is shutting down. The trace stopper blocks on the reader</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	<span class="comment">//   to finish, much like trace advancement.</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	<span class="comment">// We also want to be careful not to schedule the reader if there&#39;s no</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	<span class="comment">// reason to.</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	if trace.flushedGen.Load() == trace.readerGen.Load() || trace.workAvailable.Load() || trace.shutdown.Load() {
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		return trace.reader.Load()
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	}
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	return nil
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>}
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span><span class="comment">// Trace advancer goroutine.</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>var traceAdvancer traceAdvancerState
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>type traceAdvancerState struct {
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	timer *wakeableSleep
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	done  chan struct{}
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>}
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span><span class="comment">// start starts a new traceAdvancer.</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>func (s *traceAdvancerState) start() {
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	<span class="comment">// Start a goroutine to periodically advance the trace generation.</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	s.done = make(chan struct{})
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	s.timer = newWakeableSleep()
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	go func() {
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>		for traceEnabled() {
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>			<span class="comment">// Set a timer to wake us up</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>			s.timer.sleep(int64(debug.traceadvanceperiod))
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>			<span class="comment">// Try to advance the trace.</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>			traceAdvance(false)
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>		}
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>		s.done &lt;- struct{}{}
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>	}()
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>}
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span><span class="comment">// stop stops a traceAdvancer and blocks until it exits.</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>func (s *traceAdvancerState) stop() {
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	s.timer.wake()
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	&lt;-s.done
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	close(s.done)
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	s.timer.close()
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>}
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span><span class="comment">// traceAdvancePeriod is the approximate period between</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span><span class="comment">// new generations.</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>const defaultTraceAdvancePeriod = 1e9 <span class="comment">// 1 second.</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span><span class="comment">// wakeableSleep manages a wakeable goroutine sleep.</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span><span class="comment">// Users of this type must call init before first use and</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span><span class="comment">// close to free up resources. Once close is called, init</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span><span class="comment">// must be called before another use.</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>type wakeableSleep struct {
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	timer *timer
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	<span class="comment">// lock protects access to wakeup, but not send/recv on it.</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	lock   mutex
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	wakeup chan struct{}
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>}
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>
<span id="L917" class="ln">   917&nbsp;&nbsp;</span><span class="comment">// newWakeableSleep initializes a new wakeableSleep and returns it.</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>func newWakeableSleep() *wakeableSleep {
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>	s := new(wakeableSleep)
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>	lockInit(&amp;s.lock, lockRankWakeableSleep)
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>	s.wakeup = make(chan struct{}, 1)
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	s.timer = new(timer)
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	s.timer.arg = s
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	s.timer.f = func(s any, _ uintptr) {
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>		s.(*wakeableSleep).wake()
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	}
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	return s
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>}
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span><span class="comment">// sleep sleeps for the provided duration in nanoseconds or until</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span><span class="comment">// another goroutine calls wake.</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span><span class="comment">// Must not be called by more than one goroutine at a time and</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span><span class="comment">// must not be called concurrently with close.</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>func (s *wakeableSleep) sleep(ns int64) {
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>	resetTimer(s.timer, nanotime()+ns)
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>	lock(&amp;s.lock)
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	if raceenabled {
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>		raceacquire(unsafe.Pointer(&amp;s.lock))
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	}
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>	wakeup := s.wakeup
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>	if raceenabled {
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>		racerelease(unsafe.Pointer(&amp;s.lock))
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>	}
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>	unlock(&amp;s.lock)
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>	&lt;-wakeup
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>	stopTimer(s.timer)
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>}
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>
<span id="L950" class="ln">   950&nbsp;&nbsp;</span><span class="comment">// wake awakens any goroutine sleeping on the timer.</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span><span class="comment">// Safe for concurrent use with all other methods.</span>
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>func (s *wakeableSleep) wake() {
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>	<span class="comment">// Grab the wakeup channel, which may be nil if we&#39;re</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>	<span class="comment">// racing with close.</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>	lock(&amp;s.lock)
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	if raceenabled {
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>		raceacquire(unsafe.Pointer(&amp;s.lock))
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>	}
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>	if s.wakeup != nil {
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>		<span class="comment">// Non-blocking send.</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>		<span class="comment">// Others may also write to this channel and we don&#39;t</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>		<span class="comment">// want to block on the receiver waking up. This also</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>		<span class="comment">// effectively batches together wakeup notifications.</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		select {
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>		case s.wakeup &lt;- struct{}{}:
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>		default:
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>		}
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	}
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>	if raceenabled {
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>		racerelease(unsafe.Pointer(&amp;s.lock))
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	}
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>	unlock(&amp;s.lock)
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>}
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span><span class="comment">// close wakes any goroutine sleeping on the timer and prevents</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span><span class="comment">// further sleeping on it.</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span><span class="comment">// Once close is called, the wakeableSleep must no longer be used.</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span><span class="comment">// It must only be called once no goroutine is sleeping on the</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span><span class="comment">// timer *and* nothing else will call wake concurrently.</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>func (s *wakeableSleep) close() {
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	<span class="comment">// Set wakeup to nil so that a late timer ends up being a no-op.</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>	lock(&amp;s.lock)
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	if raceenabled {
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>		raceacquire(unsafe.Pointer(&amp;s.lock))
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	}
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	wakeup := s.wakeup
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>	s.wakeup = nil
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	<span class="comment">// Close the channel.</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	close(wakeup)
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	if raceenabled {
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>		racerelease(unsafe.Pointer(&amp;s.lock))
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	}
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	unlock(&amp;s.lock)
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	return
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>}
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>
</pre><p><a href="trace2.go?m=text">View as plain text</a></p>

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
