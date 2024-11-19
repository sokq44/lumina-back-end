<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/trace2runtime.go - Go Documentation Server</title>

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
<a href="trace2runtime.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">trace2runtime.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Runtime -&gt; tracer API.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package runtime
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	_ &#34;unsafe&#34; <span class="comment">// for go:linkname</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// gTraceState is per-G state for the tracer.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>type gTraceState struct {
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	traceSchedResourceState
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>}
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// reset resets the gTraceState for a new goroutine.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>func (s *gTraceState) reset() {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	s.seq = [2]uint64{}
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// N.B. s.statusTraced is managed and cleared separately.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>}
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// mTraceState is per-M state for the tracer.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>type mTraceState struct {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	seqlock atomic.Uintptr <span class="comment">// seqlock indicating that this M is writing to a trace buffer.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	buf     [2]*traceBuf   <span class="comment">// Per-M traceBuf for writing. Indexed by trace.gen%2.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	link    *m             <span class="comment">// Snapshot of alllink or freelink.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// pTraceState is per-P state for the tracer.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>type pTraceState struct {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	traceSchedResourceState
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// mSyscallID is the ID of the M this was bound to before entering a syscall.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	mSyscallID int64
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// maySweep indicates the sweep events should be traced.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// This is used to defer the sweep start event until a span</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// has actually been swept.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	maySweep bool
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// inSweep indicates that at least one sweep event has been traced.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	inSweep bool
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// swept and reclaimed track the number of bytes swept and reclaimed</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// by sweeping in the current sweep loop (while maySweep was true).</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	swept, reclaimed uintptr
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// traceLockInit initializes global trace locks.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>func traceLockInit() {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// Sharing a lock rank here is fine because they should never be accessed</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// together. If they are, we want to find out immediately.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	lockInit(&amp;trace.stringTab[0].lock, lockRankTraceStrings)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	lockInit(&amp;trace.stringTab[0].tab.lock, lockRankTraceStrings)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	lockInit(&amp;trace.stringTab[1].lock, lockRankTraceStrings)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	lockInit(&amp;trace.stringTab[1].tab.lock, lockRankTraceStrings)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	lockInit(&amp;trace.stackTab[0].tab.lock, lockRankTraceStackTab)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	lockInit(&amp;trace.stackTab[1].tab.lock, lockRankTraceStackTab)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	lockInit(&amp;trace.lock, lockRankTrace)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// lockRankMayTraceFlush records the lock ranking effects of a</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// potential call to traceFlush.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// nosplit because traceAcquire is nosplit.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>func lockRankMayTraceFlush() {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	lockWithRankMayAcquire(&amp;trace.lock, getLockRank(&amp;trace.lock))
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// traceBlockReason is an enumeration of reasons a goroutine might block.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// This is the interface the rest of the runtime uses to tell the</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// tracer why a goroutine blocked. The tracer then propagates this information</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// into the trace however it sees fit.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// Note that traceBlockReasons should not be compared, since reasons that are</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// distinct by name may *not* be distinct by value.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>type traceBlockReason uint8
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>const (
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	traceBlockGeneric traceBlockReason = iota
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	traceBlockForever
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	traceBlockNet
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	traceBlockSelect
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	traceBlockCondWait
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	traceBlockSync
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	traceBlockChanSend
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	traceBlockChanRecv
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	traceBlockGCMarkAssist
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	traceBlockGCSweep
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	traceBlockSystemGoroutine
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	traceBlockPreempted
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	traceBlockDebugCall
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	traceBlockUntilGCEnds
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	traceBlockSleep
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>var traceBlockReasonStrings = [...]string{
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	traceBlockGeneric:         &#34;unspecified&#34;,
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	traceBlockForever:         &#34;forever&#34;,
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	traceBlockNet:             &#34;network&#34;,
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	traceBlockSelect:          &#34;select&#34;,
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	traceBlockCondWait:        &#34;sync.(*Cond).Wait&#34;,
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	traceBlockSync:            &#34;sync&#34;,
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	traceBlockChanSend:        &#34;chan send&#34;,
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	traceBlockChanRecv:        &#34;chan receive&#34;,
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	traceBlockGCMarkAssist:    &#34;GC mark assist wait for work&#34;,
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	traceBlockGCSweep:         &#34;GC background sweeper wait&#34;,
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	traceBlockSystemGoroutine: &#34;system goroutine wait&#34;,
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	traceBlockPreempted:       &#34;preempted&#34;,
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	traceBlockDebugCall:       &#34;wait for debug call&#34;,
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	traceBlockUntilGCEnds:     &#34;wait until GC ends&#34;,
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	traceBlockSleep:           &#34;sleep&#34;,
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// traceGoStopReason is an enumeration of reasons a goroutine might yield.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// Note that traceGoStopReasons should not be compared, since reasons that are</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// distinct by name may *not* be distinct by value.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>type traceGoStopReason uint8
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>const (
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	traceGoStopGeneric traceGoStopReason = iota
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	traceGoStopGoSched
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	traceGoStopPreempted
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>var traceGoStopReasonStrings = [...]string{
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	traceGoStopGeneric:   &#34;unspecified&#34;,
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	traceGoStopGoSched:   &#34;runtime.Gosched&#34;,
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	traceGoStopPreempted: &#34;preempted&#34;,
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// traceEnabled returns true if the trace is currently enabled.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>func traceEnabled() bool {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	return trace.gen.Load() != 0
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// traceShuttingDown returns true if the trace is currently shutting down.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>func traceShuttingDown() bool {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	return trace.shutdown.Load()
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">// traceLocker represents an M writing trace events. While a traceLocker value</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">// is valid, the tracer observes all operations on the G/M/P or trace events being</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// written as happening atomically.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>type traceLocker struct {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	mp  *m
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	gen uintptr
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// debugTraceReentrancy checks if the trace is reentrant.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// This is optional because throwing in a function makes it instantly</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// not inlineable, and we want traceAcquire to be inlineable for</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// low overhead when the trace is disabled.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>const debugTraceReentrancy = false
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// traceAcquire prepares this M for writing one or more trace events.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// nosplit because it&#39;s called on the syscall path when stack movement is forbidden.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>func traceAcquire() traceLocker {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	if !traceEnabled() {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		return traceLocker{}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	return traceAcquireEnabled()
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">// traceAcquireEnabled is the traceEnabled path for traceAcquire. It&#39;s explicitly</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span><span class="comment">// broken out to make traceAcquire inlineable to keep the overhead of the tracer</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">// when it&#39;s disabled low.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">// nosplit because it&#39;s called by traceAcquire, which is nosplit.</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>func traceAcquireEnabled() traceLocker {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// Any time we acquire a traceLocker, we may flush a trace buffer. But</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">// buffer flushes are rare. Record the lock edge even if it doesn&#39;t happen</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// this time.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	lockRankMayTraceFlush()
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// Prevent preemption.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// Acquire the trace seqlock. This prevents traceAdvance from moving forward</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// until all Ms are observed to be outside of their seqlock critical section.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// Note: The seqlock is mutated here and also in traceCPUSample. If you update</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// usage of the seqlock here, make sure to also look at what traceCPUSample is</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// doing.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	seq := mp.trace.seqlock.Add(1)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	if debugTraceReentrancy &amp;&amp; seq%2 != 1 {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		throw(&#34;bad use of trace.seqlock or tracer is reentrant&#34;)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// N.B. This load of gen appears redundant with the one in traceEnabled.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// However, it&#39;s very important that the gen we use for writing to the trace</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// is acquired under a traceLocker so traceAdvance can make sure no stale</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">// gen values are being used.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	<span class="comment">// Because we&#39;re doing this load again, it also means that the trace</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// might end up being disabled when we load it. In that case we need to undo</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// what we did and bail.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	gen := trace.gen.Load()
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	if gen == 0 {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		mp.trace.seqlock.Add(1)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		releasem(mp)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		return traceLocker{}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	return traceLocker{mp, gen}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// ok returns true if the traceLocker is valid (i.e. tracing is enabled).</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// nosplit because it&#39;s called on the syscall path when stack movement is forbidden.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>func (tl traceLocker) ok() bool {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	return tl.gen != 0
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// traceRelease indicates that this M is done writing trace events.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">// nosplit because it&#39;s called on the syscall path when stack movement is forbidden.</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>func traceRelease(tl traceLocker) {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	seq := tl.mp.trace.seqlock.Add(1)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	if debugTraceReentrancy &amp;&amp; seq%2 != 0 {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		print(&#34;runtime: seq=&#34;, seq, &#34;\n&#34;)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		throw(&#34;bad use of trace.seqlock&#34;)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	releasem(tl.mp)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">// traceExitingSyscall marks a goroutine as exiting the syscall slow path.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// Must be paired with a traceExitedSyscall call.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>func traceExitingSyscall() {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	trace.exitingSyscall.Add(1)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span><span class="comment">// traceExitedSyscall marks a goroutine as having exited the syscall slow path.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>func traceExitedSyscall() {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	trace.exitingSyscall.Add(-1)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span><span class="comment">// Gomaxprocs emits a ProcsChange event.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>func (tl traceLocker) Gomaxprocs(procs int32) {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvProcsChange, traceArg(procs), tl.stack(1))
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">// ProcStart traces a ProcStart event.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">// Must be called with a valid P.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>func (tl traceLocker) ProcStart() {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	pp := tl.mp.p.ptr()
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// Procs are typically started within the scheduler when there is no user goroutine. If there is a user goroutine,</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	<span class="comment">// it must be in _Gsyscall because the only time a goroutine is allowed to have its Proc moved around from under it</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	<span class="comment">// is during a syscall.</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	tl.eventWriter(traceGoSyscall, traceProcIdle).commit(traceEvProcStart, traceArg(pp.id), pp.trace.nextSeq(tl.gen))
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span><span class="comment">// ProcStop traces a ProcStop event.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>func (tl traceLocker) ProcStop(pp *p) {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">// The only time a goroutine is allowed to have its Proc moved around</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	<span class="comment">// from under it is during a syscall.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	tl.eventWriter(traceGoSyscall, traceProcRunning).commit(traceEvProcStop)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// GCActive traces a GCActive event.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">// Must be emitted by an actively running goroutine on an active P. This restriction can be changed</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">// easily and only depends on where it&#39;s currently called.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>func (tl traceLocker) GCActive() {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGCActive, traceArg(trace.seqGC))
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	<span class="comment">// N.B. Only one GC can be running at a time, so this is naturally</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	<span class="comment">// serialized by the caller.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	trace.seqGC++
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">// GCStart traces a GCBegin event.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span><span class="comment">// Must be emitted by an actively running goroutine on an active P. This restriction can be changed</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span><span class="comment">// easily and only depends on where it&#39;s currently called.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>func (tl traceLocker) GCStart() {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGCBegin, traceArg(trace.seqGC), tl.stack(3))
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	<span class="comment">// N.B. Only one GC can be running at a time, so this is naturally</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">// serialized by the caller.</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	trace.seqGC++
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span><span class="comment">// GCDone traces a GCEnd event.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">// Must be emitted by an actively running goroutine on an active P. This restriction can be changed</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">// easily and only depends on where it&#39;s currently called.</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>func (tl traceLocker) GCDone() {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGCEnd, traceArg(trace.seqGC))
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	<span class="comment">// N.B. Only one GC can be running at a time, so this is naturally</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	<span class="comment">// serialized by the caller.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	trace.seqGC++
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">// STWStart traces a STWBegin event.</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>func (tl traceLocker) STWStart(reason stwReason) {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// Although the current P may be in _Pgcstop here, we model the P as running during the STW. This deviates from the</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	<span class="comment">// runtime&#39;s state tracking, but it&#39;s more accurate and doesn&#39;t result in any loss of information.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvSTWBegin, tl.string(reason.String()), tl.stack(2))
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span><span class="comment">// STWDone traces a STWEnd event.</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>func (tl traceLocker) STWDone() {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	<span class="comment">// Although the current P may be in _Pgcstop here, we model the P as running during the STW. This deviates from the</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	<span class="comment">// runtime&#39;s state tracking, but it&#39;s more accurate and doesn&#39;t result in any loss of information.</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvSTWEnd)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// GCSweepStart prepares to trace a sweep loop. This does not</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// emit any events until traceGCSweepSpan is called.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span><span class="comment">// GCSweepStart must be paired with traceGCSweepDone and there</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// must be no preemption points between these two calls.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">// Must be called with a valid P.</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>func (tl traceLocker) GCSweepStart() {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// Delay the actual GCSweepBegin event until the first span</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// sweep. If we don&#39;t sweep anything, don&#39;t emit any events.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	pp := tl.mp.p.ptr()
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	if pp.trace.maySweep {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		throw(&#34;double traceGCSweepStart&#34;)
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	pp.trace.maySweep, pp.trace.swept, pp.trace.reclaimed = true, 0, 0
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// GCSweepSpan traces the sweep of a single span. If this is</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// the first span swept since traceGCSweepStart was called, this</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// will emit a GCSweepBegin event.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// This may be called outside a traceGCSweepStart/traceGCSweepDone</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// pair; however, it will not emit any trace events in this case.</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// Must be called with a valid P.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>func (tl traceLocker) GCSweepSpan(bytesSwept uintptr) {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	pp := tl.mp.p.ptr()
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	if pp.trace.maySweep {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		if pp.trace.swept == 0 {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGCSweepBegin, tl.stack(1))
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			pp.trace.inSweep = true
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		pp.trace.swept += bytesSwept
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">// GCSweepDone finishes tracing a sweep loop. If any memory was</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// swept (i.e. traceGCSweepSpan emitted an event) then this will emit</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// a GCSweepEnd event.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">// Must be called with a valid P.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>func (tl traceLocker) GCSweepDone() {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	pp := tl.mp.p.ptr()
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	if !pp.trace.maySweep {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		throw(&#34;missing traceGCSweepStart&#34;)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	if pp.trace.inSweep {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGCSweepEnd, traceArg(pp.trace.swept), traceArg(pp.trace.reclaimed))
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		pp.trace.inSweep = false
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	pp.trace.maySweep = false
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span><span class="comment">// GCMarkAssistStart emits a MarkAssistBegin event.</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>func (tl traceLocker) GCMarkAssistStart() {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGCMarkAssistBegin, tl.stack(1))
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span><span class="comment">// GCMarkAssistDone emits a MarkAssistEnd event.</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>func (tl traceLocker) GCMarkAssistDone() {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGCMarkAssistEnd)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">// GoCreate emits a GoCreate event.</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>func (tl traceLocker) GoCreate(newg *g, pc uintptr) {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	newg.trace.setStatusTraced(tl.gen)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGoCreate, traceArg(newg.goid), tl.startPC(pc), tl.stack(2))
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>}
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span><span class="comment">// GoStart emits a GoStart event.</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span><span class="comment">// Must be called with a valid P.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>func (tl traceLocker) GoStart() {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	gp := getg().m.curg
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	pp := gp.m.p
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	w := tl.eventWriter(traceGoRunnable, traceProcRunning)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	w = w.write(traceEvGoStart, traceArg(gp.goid), gp.trace.nextSeq(tl.gen))
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	if pp.ptr().gcMarkWorkerMode != gcMarkWorkerNotWorker {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		w = w.write(traceEvGoLabel, trace.markWorkerLabels[tl.gen%2][pp.ptr().gcMarkWorkerMode])
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	w.end()
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span><span class="comment">// GoEnd emits a GoDestroy event.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span><span class="comment">// TODO(mknyszek): Rename this to GoDestroy.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>func (tl traceLocker) GoEnd() {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGoDestroy)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span><span class="comment">// GoSched emits a GoStop event with a GoSched reason.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>func (tl traceLocker) GoSched() {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	tl.GoStop(traceGoStopGoSched)
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span><span class="comment">// GoPreempt emits a GoStop event with a GoPreempted reason.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>func (tl traceLocker) GoPreempt() {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	tl.GoStop(traceGoStopPreempted)
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span><span class="comment">// GoStop emits a GoStop event with the provided reason.</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>func (tl traceLocker) GoStop(reason traceGoStopReason) {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGoStop, traceArg(trace.goStopReasons[tl.gen%2][reason]), tl.stack(1))
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span><span class="comment">// GoPark emits a GoBlock event with the provided reason.</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span><span class="comment">// TODO(mknyszek): Replace traceBlockReason with waitReason. It&#39;s silly</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span><span class="comment">// that we have both, and waitReason is way more descriptive.</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>func (tl traceLocker) GoPark(reason traceBlockReason, skip int) {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGoBlock, traceArg(trace.goBlockReasons[tl.gen%2][reason]), tl.stack(skip))
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>}
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span><span class="comment">// GoUnpark emits a GoUnblock event.</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>func (tl traceLocker) GoUnpark(gp *g, skip int) {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	<span class="comment">// Emit a GoWaiting status if necessary for the unblocked goroutine.</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	w := tl.eventWriter(traceGoRunning, traceProcRunning)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	if !gp.trace.statusWasTraced(tl.gen) &amp;&amp; gp.trace.acquireStatus(tl.gen) {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		<span class="comment">// Careful: don&#39;t use the event writer. We never want status or in-progress events</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		<span class="comment">// to trigger more in-progress events.</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		w.w = w.w.writeGoStatus(gp.goid, -1, traceGoWaiting, gp.inMarkAssist)
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	w.commit(traceEvGoUnblock, traceArg(gp.goid), gp.trace.nextSeq(tl.gen), tl.stack(skip))
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span><span class="comment">// GoSysCall emits a GoSyscallBegin event.</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span><span class="comment">// Must be called with a valid P.</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>func (tl traceLocker) GoSysCall() {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	var skip int
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	switch {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	case tracefpunwindoff():
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		<span class="comment">// Unwind by skipping 1 frame relative to gp.syscallsp which is captured 3</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		<span class="comment">// results by hard coding the number of frames in between our caller and the</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		<span class="comment">// actual syscall, see cases below.</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		<span class="comment">// TODO(felixge): Implement gp.syscallbp to avoid this workaround?</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		skip = 1
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	case GOOS == &#34;solaris&#34; || GOOS == &#34;illumos&#34;:
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		<span class="comment">// These platforms don&#39;t use a libc_read_trampoline.</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		skip = 3
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	default:
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		<span class="comment">// Skip the extra trampoline frame used on most systems.</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		skip = 4
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	<span class="comment">// Scribble down the M that the P is currently attached to.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	pp := tl.mp.p.ptr()
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	pp.trace.mSyscallID = int64(tl.mp.procid)
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvGoSyscallBegin, pp.trace.nextSeq(tl.gen), tl.stack(skip))
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span><span class="comment">// GoSysExit emits a GoSyscallEnd event, possibly along with a GoSyscallBlocked event</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span><span class="comment">// if lostP is true.</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span><span class="comment">// lostP must be true in all cases that a goroutine loses its P during a syscall.</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span><span class="comment">// This means it&#39;s not sufficient to check if it has no P. In particular, it needs to be</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span><span class="comment">// true in the following cases:</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span><span class="comment">// - The goroutine lost its P, it ran some other code, and then got it back. It&#39;s now running with that P.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span><span class="comment">// - The goroutine lost its P and was unable to reacquire it, and is now running without a P.</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span><span class="comment">// - The goroutine lost its P and acquired a different one, and is now running with that P.</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>func (tl traceLocker) GoSysExit(lostP bool) {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	ev := traceEvGoSyscallEnd
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	procStatus := traceProcSyscall <span class="comment">// Procs implicitly enter traceProcSyscall on GoSyscallBegin.</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	if lostP {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		ev = traceEvGoSyscallEndBlocked
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		procStatus = traceProcRunning <span class="comment">// If a G has a P when emitting this event, it reacquired a P and is indeed running.</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	} else {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		tl.mp.p.ptr().trace.mSyscallID = -1
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	tl.eventWriter(traceGoSyscall, procStatus).commit(ev)
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>}
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span><span class="comment">// ProcSteal indicates that our current M stole a P from another M.</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span><span class="comment">// inSyscall indicates that we&#39;re stealing the P from a syscall context.</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span><span class="comment">// The caller must have ownership of pp.</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>func (tl traceLocker) ProcSteal(pp *p, inSyscall bool) {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	<span class="comment">// Grab the M ID we stole from.</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	mStolenFrom := pp.trace.mSyscallID
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	pp.trace.mSyscallID = -1
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	<span class="comment">// The status of the proc and goroutine, if we need to emit one here, is not evident from the</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	<span class="comment">// context of just emitting this event alone. There are two cases. Either we&#39;re trying to steal</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	<span class="comment">// the P just to get its attention (e.g. STW or sysmon retake) or we&#39;re trying to steal a P for</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	<span class="comment">// ourselves specifically to keep running. The two contexts look different, but can be summarized</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	<span class="comment">// fairly succinctly. In the former, we&#39;re a regular running goroutine and proc, if we have either.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	<span class="comment">// In the latter, we&#39;re a goroutine in a syscall.</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	goStatus := traceGoRunning
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	procStatus := traceProcRunning
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	if inSyscall {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		goStatus = traceGoSyscall
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		procStatus = traceProcSyscallAbandoned
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	w := tl.eventWriter(goStatus, procStatus)
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	<span class="comment">// Emit the status of the P we&#39;re stealing. We may have *just* done this when creating the event</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	<span class="comment">// writer but it&#39;s not guaranteed, even if inSyscall is true. Although it might seem like from a</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	<span class="comment">// syscall context we&#39;re always stealing a P for ourselves, we may have not wired it up yet (so</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	<span class="comment">// it wouldn&#39;t be visible to eventWriter) or we may not even intend to wire it up to ourselves</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	<span class="comment">// at all (e.g. entersyscall_gcwait).</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	if !pp.trace.statusWasTraced(tl.gen) &amp;&amp; pp.trace.acquireStatus(tl.gen) {
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		<span class="comment">// Careful: don&#39;t use the event writer. We never want status or in-progress events</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		<span class="comment">// to trigger more in-progress events.</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		w.w = w.w.writeProcStatus(uint64(pp.id), traceProcSyscallAbandoned, pp.trace.inSweep)
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	w.commit(traceEvProcSteal, traceArg(pp.id), pp.trace.nextSeq(tl.gen), traceArg(mStolenFrom))
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span><span class="comment">// GoSysBlock is a no-op in the new tracer.</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>func (tl traceLocker) GoSysBlock(pp *p) {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>}
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">// HeapAlloc emits a HeapAlloc event.</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>func (tl traceLocker) HeapAlloc(live uint64) {
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvHeapAlloc, traceArg(live))
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span><span class="comment">// HeapGoal reads the current heap goal and emits a HeapGoal event.</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>func (tl traceLocker) HeapGoal() {
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	heapGoal := gcController.heapGoal()
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	if heapGoal == ^uint64(0) {
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		<span class="comment">// Heap-based triggering is disabled.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		heapGoal = 0
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvHeapGoal, traceArg(heapGoal))
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>}
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span><span class="comment">// OneNewExtraM is a no-op in the new tracer. This is worth keeping around though because</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span><span class="comment">// it&#39;s a good place to insert a thread-level event about the new extra M.</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>func (tl traceLocker) OneNewExtraM(_ *g) {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span><span class="comment">// GoCreateSyscall indicates that a goroutine has transitioned from dead to GoSyscall.</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span><span class="comment">// Unlike GoCreate, the caller must be running on gp.</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span><span class="comment">// This occurs when C code calls into Go. On pthread platforms it occurs only when</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span><span class="comment">// a C thread calls into Go code for the first time.</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>func (tl traceLocker) GoCreateSyscall(gp *g) {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	<span class="comment">// N.B. We should never trace a status for this goroutine (which we&#39;re currently running on),</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	<span class="comment">// since we want this to appear like goroutine creation.</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	gp.trace.setStatusTraced(tl.gen)
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	tl.eventWriter(traceGoBad, traceProcBad).commit(traceEvGoCreateSyscall, traceArg(gp.goid))
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span><span class="comment">// GoDestroySyscall indicates that a goroutine has transitioned from GoSyscall to dead.</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span><span class="comment">// Must not have a P.</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span><span class="comment">// This occurs when Go code returns back to C. On pthread platforms it occurs only when</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span><span class="comment">// the C thread is destroyed.</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>func (tl traceLocker) GoDestroySyscall() {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	<span class="comment">// N.B. If we trace a status here, we must never have a P, and we must be on a goroutine</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	<span class="comment">// that is in the syscall state.</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	tl.eventWriter(traceGoSyscall, traceProcBad).commit(traceEvGoDestroySyscall)
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>}
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span><span class="comment">// To access runtime functions from runtime/trace.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span><span class="comment">// See runtime/trace/annotation.go</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span><span class="comment">// trace_userTaskCreate emits a UserTaskCreate event.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span><span class="comment">//go:linkname trace_userTaskCreate runtime/trace.userTaskCreate</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>func trace_userTaskCreate(id, parentID uint64, taskType string) {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	tl := traceAcquire()
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	if !tl.ok() {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		<span class="comment">// Need to do this check because the caller won&#39;t have it.</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		return
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	}
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvUserTaskBegin, traceArg(id), traceArg(parentID), tl.string(taskType), tl.stack(3))
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	traceRelease(tl)
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>}
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span><span class="comment">// trace_userTaskEnd emits a UserTaskEnd event.</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span><span class="comment">//go:linkname trace_userTaskEnd runtime/trace.userTaskEnd</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>func trace_userTaskEnd(id uint64) {
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	tl := traceAcquire()
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	if !tl.ok() {
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		<span class="comment">// Need to do this check because the caller won&#39;t have it.</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		return
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	}
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvUserTaskEnd, traceArg(id), tl.stack(2))
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	traceRelease(tl)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span><span class="comment">// trace_userTaskEnd emits a UserRegionBegin or UserRegionEnd event,</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span><span class="comment">// depending on mode (0 == Begin, 1 == End).</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span><span class="comment">// TODO(mknyszek): Just make this two functions.</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span><span class="comment">//go:linkname trace_userRegion runtime/trace.userRegion</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>func trace_userRegion(id, mode uint64, name string) {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	tl := traceAcquire()
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	if !tl.ok() {
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		<span class="comment">// Need to do this check because the caller won&#39;t have it.</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		return
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	}
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	var ev traceEv
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	switch mode {
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	case 0:
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		ev = traceEvUserRegionBegin
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	case 1:
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		ev = traceEvUserRegionEnd
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	default:
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>		return
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	}
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(ev, traceArg(id), tl.string(name), tl.stack(3))
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	traceRelease(tl)
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>}
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span><span class="comment">// trace_userTaskEnd emits a UserRegionBegin or UserRegionEnd event.</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span><span class="comment">//go:linkname trace_userLog runtime/trace.userLog</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>func trace_userLog(id uint64, category, message string) {
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	tl := traceAcquire()
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	if !tl.ok() {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		<span class="comment">// Need to do this check because the caller won&#39;t have it.</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		return
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	tl.eventWriter(traceGoRunning, traceProcRunning).commit(traceEvUserLog, traceArg(id), tl.string(category), tl.uniqueString(message), tl.stack(3))
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	traceRelease(tl)
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>}
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span><span class="comment">// traceProcFree is called when a P is destroyed.</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span><span class="comment">// This must run on the system stack to match the old tracer.</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>func traceProcFree(_ *p) {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span><span class="comment">// traceThreadDestroy is called when a thread is removed from</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span><span class="comment">// sched.freem.</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span><span class="comment">// mp must not be able to emit trace events anymore.</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span><span class="comment">// sched.lock must be held to synchronize with traceAdvance.</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>func traceThreadDestroy(mp *m) {
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	assertLockHeld(&amp;sched.lock)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	<span class="comment">// Flush all outstanding buffers to maintain the invariant</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	<span class="comment">// that an M only has active buffers while on sched.freem</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	<span class="comment">// or allm.</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	<span class="comment">// Perform a traceAcquire/traceRelease on behalf of mp to</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	<span class="comment">// synchronize with the tracer trying to flush our buffer</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	<span class="comment">// as well.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	seq := mp.trace.seqlock.Add(1)
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	if debugTraceReentrancy &amp;&amp; seq%2 != 1 {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		throw(&#34;bad use of trace.seqlock or tracer is reentrant&#34;)
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	}
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		lock(&amp;trace.lock)
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		for i := range mp.trace.buf {
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>			if mp.trace.buf[i] != nil {
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>				<span class="comment">// N.B. traceBufFlush accepts a generation, but it</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>				<span class="comment">// really just cares about gen%2.</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>				traceBufFlush(mp.trace.buf[i], uintptr(i))
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>				mp.trace.buf[i] = nil
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>			}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		}
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		unlock(&amp;trace.lock)
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	})
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	seq1 := mp.trace.seqlock.Add(1)
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	if seq1 != seq+1 {
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		print(&#34;runtime: seq1=&#34;, seq1, &#34;\n&#34;)
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>		throw(&#34;bad use of trace.seqlock&#34;)
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	}
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span><span class="comment">// Not used in the new tracer; solely for compatibility with the old tracer.</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span><span class="comment">// nosplit because it&#39;s called from exitsyscall without a P.</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>func (_ traceLocker) RecordSyscallExitedTime(_ *g, _ *p) {
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>}
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>
</pre><p><a href="trace2runtime.go?m=text">View as plain text</a></p>

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
