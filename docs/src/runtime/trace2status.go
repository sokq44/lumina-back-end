<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/trace2status.go - Go Documentation Server</title>

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
<a href="trace2status.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">trace2status.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Trace goroutine and P status management.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package runtime
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import &#34;runtime/internal/atomic&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// traceGoStatus is the status of a goroutine.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// They correspond directly to the various goroutine</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// statuses.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>type traceGoStatus uint8
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>const (
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	traceGoBad traceGoStatus = iota
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	traceGoRunnable
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	traceGoRunning
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	traceGoSyscall
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	traceGoWaiting
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// traceProcStatus is the status of a P.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// They mostly correspond to the various P statuses.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>type traceProcStatus uint8
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>const (
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	traceProcBad traceProcStatus = iota
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	traceProcRunning
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	traceProcIdle
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	traceProcSyscall
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// traceProcSyscallAbandoned is a special case of</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// traceProcSyscall. It&#39;s used in the very specific case</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// where the first a P is mentioned in a generation is</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// part of a ProcSteal event. If that&#39;s the first time</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// it&#39;s mentioned, then there&#39;s no GoSyscallBegin to</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// connect the P stealing back to at that point. This</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// special state indicates this to the parser, so it</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// doesn&#39;t try to find a GoSyscallEndBlocked that</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// corresponds with the ProcSteal.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	traceProcSyscallAbandoned
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// writeGoStatus emits a GoStatus event as well as any active ranges on the goroutine.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>func (w traceWriter) writeGoStatus(goid uint64, mid int64, status traceGoStatus, markAssist bool) traceWriter {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// The status should never be bad. Some invariant must have been violated.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	if status == traceGoBad {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		print(&#34;runtime: goid=&#34;, goid, &#34;\n&#34;)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		throw(&#34;attempted to trace a bad status for a goroutine&#34;)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// Trace the status.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	w = w.event(traceEvGoStatus, traceArg(goid), traceArg(uint64(mid)), traceArg(status))
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// Trace any special ranges that are in-progress.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	if markAssist {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		w = w.event(traceEvGCMarkAssistActive, traceArg(goid))
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	return w
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// writeProcStatusForP emits a ProcStatus event for the provided p based on its status.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// The caller must fully own pp and it must be prevented from transitioning (e.g. this can be</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// called by a forEachP callback or from a STW).</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func (w traceWriter) writeProcStatusForP(pp *p, inSTW bool) traceWriter {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	if !pp.trace.acquireStatus(w.gen) {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		return w
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	var status traceProcStatus
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	switch pp.status {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	case _Pidle, _Pgcstop:
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		status = traceProcIdle
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		if pp.status == _Pgcstop &amp;&amp; inSTW {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			<span class="comment">// N.B. a P that is running and currently has the world stopped will be</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			<span class="comment">// in _Pgcstop, but we model it as running in the tracer.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			status = traceProcRunning
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	case _Prunning:
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		status = traceProcRunning
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		<span class="comment">// There&#39;s a short window wherein the goroutine may have entered _Gsyscall</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		<span class="comment">// but it still owns the P (it&#39;s not in _Psyscall yet). The goroutine entering</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		<span class="comment">// _Gsyscall is the tracer&#39;s signal that the P its bound to is also in a syscall,</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		<span class="comment">// so we need to emit a status that matches. See #64318.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		if w.mp.p.ptr() == pp &amp;&amp; w.mp.curg != nil &amp;&amp; readgstatus(w.mp.curg)&amp;^_Gscan == _Gsyscall {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			status = traceProcSyscall
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	case _Psyscall:
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		status = traceProcSyscall
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	default:
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		throw(&#34;attempt to trace invalid or unsupported P status&#34;)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	w = w.writeProcStatus(uint64(pp.id), status, pp.trace.inSweep)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	return w
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// writeProcStatus emits a ProcStatus event with all the provided information.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// The caller must have taken ownership of a P&#39;s status writing, and the P must be</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// prevented from transitioning.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>func (w traceWriter) writeProcStatus(pid uint64, status traceProcStatus, inSweep bool) traceWriter {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">// The status should never be bad. Some invariant must have been violated.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	if status == traceProcBad {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		print(&#34;runtime: pid=&#34;, pid, &#34;\n&#34;)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		throw(&#34;attempted to trace a bad status for a proc&#34;)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// Trace the status.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	w = w.event(traceEvProcStatus, traceArg(pid), traceArg(status))
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// Trace any special ranges that are in-progress.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	if inSweep {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		w = w.event(traceEvGCSweepActive, traceArg(pid))
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	return w
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// goStatusToTraceGoStatus translates the internal status to tracGoStatus.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// status must not be _Gdead or any status whose name has the suffix &#34;_unused.&#34;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>func goStatusToTraceGoStatus(status uint32, wr waitReason) traceGoStatus {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// N.B. Ignore the _Gscan bit. We don&#39;t model it in the tracer.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	var tgs traceGoStatus
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	switch status &amp;^ _Gscan {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	case _Grunnable:
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		tgs = traceGoRunnable
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	case _Grunning, _Gcopystack:
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		tgs = traceGoRunning
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	case _Gsyscall:
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		tgs = traceGoSyscall
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	case _Gwaiting, _Gpreempted:
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		<span class="comment">// There are a number of cases where a G might end up in</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		<span class="comment">// _Gwaiting but it&#39;s actually running in a non-preemptive</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		<span class="comment">// state but needs to present itself as preempted to the</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		<span class="comment">// garbage collector. In these cases, we&#39;re not going to</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		<span class="comment">// emit an event, and we want these goroutines to appear in</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		<span class="comment">// the final trace as if they&#39;re running, not blocked.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		tgs = traceGoWaiting
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		if status == _Gwaiting &amp;&amp;
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			wr == waitReasonStoppingTheWorld ||
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			wr == waitReasonGCMarkTermination ||
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			wr == waitReasonGarbageCollection ||
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			wr == waitReasonTraceProcStatus ||
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			wr == waitReasonPageTraceFlush ||
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			wr == waitReasonGCWorkerActive {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			tgs = traceGoRunning
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	case _Gdead:
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		throw(&#34;tried to trace dead goroutine&#34;)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	default:
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		throw(&#34;tried to trace goroutine with invalid or unsupported status&#34;)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	return tgs
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// traceSchedResourceState is shared state for scheduling resources (i.e. fields common to</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// both Gs and Ps).</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>type traceSchedResourceState struct {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// statusTraced indicates whether a status event was traced for this resource</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// a particular generation.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// There are 3 of these because when transitioning across generations, traceAdvance</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// needs to be able to reliably observe whether a status was traced for the previous</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// generation, while we need to clear the value for the next generation.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	statusTraced [3]atomic.Uint32
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// seq is the sequence counter for this scheduling resource&#39;s events.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// The purpose of the sequence counter is to establish a partial order between</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// events that don&#39;t obviously happen serially (same M) in the stream ofevents.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// There are two of these so that we can reset the counter on each generation.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// This saves space in the resulting trace by keeping the counter small and allows</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// GoStatus and GoCreate events to omit a sequence number (implicitly 0).</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	seq [2]uint64
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">// acquireStatus acquires the right to emit a Status event for the scheduling resource.</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>func (r *traceSchedResourceState) acquireStatus(gen uintptr) bool {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if !r.statusTraced[gen%3].CompareAndSwap(0, 1) {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		return false
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	r.readyNextGen(gen)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	return true
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// readyNextGen readies r for the generation following gen.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>func (r *traceSchedResourceState) readyNextGen(gen uintptr) {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	nextGen := traceNextGen(gen)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	r.seq[nextGen%2] = 0
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	r.statusTraced[nextGen%3].Store(0)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// statusWasTraced returns true if the sched resource&#39;s status was already acquired for tracing.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>func (r *traceSchedResourceState) statusWasTraced(gen uintptr) bool {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	return r.statusTraced[gen%3].Load() != 0
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">// setStatusTraced indicates that the resource&#39;s status was already traced, for example</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// when a goroutine is created.</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>func (r *traceSchedResourceState) setStatusTraced(gen uintptr) {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	r.statusTraced[gen%3].Store(1)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// nextSeq returns the next sequence number for the resource.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>func (r *traceSchedResourceState) nextSeq(gen uintptr) traceArg {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	r.seq[gen%2]++
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	return traceArg(r.seq[gen%2])
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
</pre><p><a href="trace2status.go?m=text">View as plain text</a></p>

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
