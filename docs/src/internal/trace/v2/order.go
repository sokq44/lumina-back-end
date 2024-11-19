<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/internal/trace/v2/order.go - Go Documentation Server</title>

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
<a href="order.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/internal">internal</a>/<a href="http://localhost:8080/src/internal/trace">trace</a>/<a href="http://localhost:8080/src/internal/trace/v2">v2</a>/<span class="text-muted">order.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/internal/trace/v2">internal/trace/v2</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2023 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package trace
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/trace/v2/event&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/trace/v2/event/go122&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;internal/trace/v2/version&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// ordering emulates Go scheduler state for both validation and</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// for putting events in the right order.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type ordering struct {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	gStates     map[GoID]*gState
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	pStates     map[ProcID]*pState <span class="comment">// TODO: The keys are dense, so this can be a slice.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	mStates     map[ThreadID]*mState
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	activeTasks map[TaskID]taskState
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	gcSeq       uint64
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	gcState     gcState
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	initialGen  uint64
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// Some events like GoDestroySyscall produce two events instead of one.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// extraEvent is this extra space. advance must not be called unless</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// the extraEvent has been consumed with consumeExtraEvent.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Replace this with a more formal queue.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	extraEvent Event
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// consumeExtraEvent consumes the extra event.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>func (o *ordering) consumeExtraEvent() Event {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	if o.extraEvent.Kind() == EventBad {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		return Event{}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	r := o.extraEvent
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	o.extraEvent = Event{}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	return r
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// advance checks if it&#39;s valid to proceed with ev which came from thread m.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// Returns the schedCtx at the point of the event, whether it&#39;s OK to advance</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// with this event, and any error encountered in validation.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// It assumes the gen value passed to it is monotonically increasing across calls.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// If any error is returned, then the trace is broken and trace parsing must cease.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// If it&#39;s not valid to advance with ev, but no error was encountered, the caller</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// should attempt to advance with other candidate events from other threads. If the</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// caller runs out of candidates, the trace is invalid.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>func (o *ordering) advance(ev *baseEvent, evt *evTable, m ThreadID, gen uint64) (schedCtx, bool, error) {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	if o.initialGen == 0 {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		<span class="comment">// Set the initial gen if necessary.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		o.initialGen = gen
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	var curCtx, newCtx schedCtx
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	curCtx.M = m
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	newCtx.M = m
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	if m == NoThread {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		curCtx.P = NoProc
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		curCtx.G = NoGoroutine
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		newCtx = curCtx
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	} else {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		<span class="comment">// Pull out or create the mState for this event.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		ms, ok := o.mStates[m]
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		if !ok {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			ms = &amp;mState{
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>				g: NoGoroutine,
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>				p: NoProc,
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			o.mStates[m] = ms
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		curCtx.P = ms.p
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		curCtx.G = ms.g
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		newCtx = curCtx
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		defer func() {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			<span class="comment">// Update the mState for this event.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			ms.p = newCtx.P
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			ms.g = newCtx.G
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		}()
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	switch typ := ev.typ; typ {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// Handle procs.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	case go122.EvProcStatus:
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		pid := ProcID(ev.args[0])
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		status := go122.ProcStatus(ev.args[1])
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		if int(status) &gt;= len(go122ProcStatus2ProcState) {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;invalid status for proc %d: %d&#34;, pid, status)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		oldState := go122ProcStatus2ProcState[status]
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		if s, ok := o.pStates[pid]; ok {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			if status == go122.ProcSyscallAbandoned &amp;&amp; s.status == go122.ProcSyscall {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>				<span class="comment">// ProcSyscallAbandoned is a special case of ProcSyscall. It indicates a</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>				<span class="comment">// potential loss of information, but if we&#39;re already in ProcSyscall,</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>				<span class="comment">// we haven&#39;t lost the relevant information. Promote the status and advance.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>				oldState = ProcRunning
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>				ev.args[1] = uint64(go122.ProcSyscall)
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			} else if status == go122.ProcSyscallAbandoned &amp;&amp; s.status == go122.ProcSyscallAbandoned {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>				<span class="comment">// If we&#39;re passing through ProcSyscallAbandoned, then there&#39;s no promotion</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>				<span class="comment">// to do. We&#39;ve lost the M that this P is associated with. However it got there,</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>				<span class="comment">// it&#39;s going to appear as idle in the API, so pass through as idle.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>				oldState = ProcIdle
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>				ev.args[1] = uint64(go122.ProcSyscallAbandoned)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			} else if s.status != status {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>				return curCtx, false, fmt.Errorf(&#34;inconsistent status for proc %d: old %v vs. new %v&#34;, pid, s.status, status)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			s.seq = makeSeq(gen, 0) <span class="comment">// Reset seq.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		} else {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			o.pStates[pid] = &amp;pState{id: pid, status: status, seq: makeSeq(gen, 0)}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			if gen == o.initialGen {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>				oldState = ProcUndetermined
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			} else {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>				oldState = ProcNotExist
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		ev.extra(version.Go122)[0] = uint64(oldState) <span class="comment">// Smuggle in the old state for StateTransition.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		<span class="comment">// Bind the proc to the new context, if it&#39;s running.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		if status == go122.ProcRunning || status == go122.ProcSyscall {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			newCtx.P = pid
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		<span class="comment">// If we&#39;re advancing through ProcSyscallAbandoned *but* oldState is running then we&#39;ve</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		<span class="comment">// promoted it to ProcSyscall. However, because it&#39;s ProcSyscallAbandoned, we know this</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		<span class="comment">// P is about to get stolen and its status very likely isn&#39;t being emitted by the same</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		<span class="comment">// thread it was bound to. Since this status is Running -&gt; Running and Running is binding,</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// we need to make sure we emit it in the right context: the context to which it is bound.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// Find it, and set our current context to it.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		if status == go122.ProcSyscallAbandoned &amp;&amp; oldState == ProcRunning {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			<span class="comment">// N.B. This is slow but it should be fairly rare.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			found := false
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			for mid, ms := range o.mStates {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>				if ms.p == pid {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>					curCtx.M = mid
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>					curCtx.P = pid
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>					curCtx.G = ms.g
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>					found = true
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>				}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			if !found {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				return curCtx, false, fmt.Errorf(&#34;failed to find sched context for proc %d that&#39;s about to be stolen&#34;, pid)
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	case go122.EvProcStart:
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		pid := ProcID(ev.args[0])
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		seq := makeSeq(gen, ev.args[1])
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		<span class="comment">// Try to advance. We might fail here due to sequencing, because the P hasn&#39;t</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// had a status emitted, or because we already have a P and we&#39;re in a syscall,</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// and we haven&#39;t observed that it was stolen from us yet.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		state, ok := o.pStates[pid]
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		if !ok || state.status != go122.ProcIdle || !seq.succeeds(state.seq) || curCtx.P != NoProc {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			<span class="comment">// We can&#39;t make an inference as to whether this is bad. We could just be seeing</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			<span class="comment">// a ProcStart on a different M before the proc&#39;s state was emitted, or before we</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			<span class="comment">// got to the right point in the trace.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			<span class="comment">// Note that we also don&#39;t advance here if we have a P and we&#39;re in a syscall.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			return curCtx, false, nil
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// We can advance this P. Check some invariants.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		<span class="comment">// We might have a goroutine if a goroutine is exiting a syscall.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		reqs := event.SchedReqs{Thread: event.MustHave, Proc: event.MustNotHave, Goroutine: event.MayHave}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, reqs); err != nil {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		state.status = go122.ProcRunning
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		state.seq = seq
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		newCtx.P = pid
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	case go122.EvProcStop:
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		<span class="comment">// We must be able to advance this P.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		<span class="comment">// There are 2 ways a P can stop: ProcStop and ProcSteal. ProcStop is used when the P</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		<span class="comment">// is stopped by the same M that started it, while ProcSteal is used when another M</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		<span class="comment">// steals the P by stopping it from a distance.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		<span class="comment">// Since a P is bound to an M, and we&#39;re stopping on the same M we started, it must</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		<span class="comment">// always be possible to advance the current M&#39;s P from a ProcStop. This is also why</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		<span class="comment">// ProcStop doesn&#39;t need a sequence number.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		state, ok := o.pStates[curCtx.P]
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		if !ok {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;event %s for proc (%v) that doesn&#39;t exist&#34;, go122.EventString(typ), curCtx.P)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		if state.status != go122.ProcRunning &amp;&amp; state.status != go122.ProcSyscall {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;%s event for proc that&#39;s not %s or %s&#34;, go122.EventString(typ), go122.ProcRunning, go122.ProcSyscall)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		reqs := event.SchedReqs{Thread: event.MustHave, Proc: event.MustHave, Goroutine: event.MayHave}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, reqs); err != nil {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		state.status = go122.ProcIdle
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		newCtx.P = NoProc
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	case go122.EvProcSteal:
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		pid := ProcID(ev.args[0])
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		seq := makeSeq(gen, ev.args[1])
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		state, ok := o.pStates[pid]
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		if !ok || (state.status != go122.ProcSyscall &amp;&amp; state.status != go122.ProcSyscallAbandoned) || !seq.succeeds(state.seq) {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			<span class="comment">// We can&#39;t make an inference as to whether this is bad. We could just be seeing</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			<span class="comment">// a ProcStart on a different M before the proc&#39;s state was emitted, or before we</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			<span class="comment">// got to the right point in the trace.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			return curCtx, false, nil
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		<span class="comment">// We can advance this P. Check some invariants.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		reqs := event.SchedReqs{Thread: event.MustHave, Proc: event.MayHave, Goroutine: event.MayHave}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, reqs); err != nil {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// Smuggle in the P state that let us advance so we can surface information to the event.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		<span class="comment">// Specifically, we need to make sure that the event is interpreted not as a transition of</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		<span class="comment">// ProcRunning -&gt; ProcIdle but ProcIdle -&gt; ProcIdle instead.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		<span class="comment">// ProcRunning is binding, but we may be running with a P on the current M and we can&#39;t</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		<span class="comment">// bind another P. This P is about to go ProcIdle anyway.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		oldStatus := state.status
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		ev.extra(version.Go122)[0] = uint64(oldStatus)
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		<span class="comment">// Update the P&#39;s status and sequence number.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		state.status = go122.ProcIdle
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		state.seq = seq
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		<span class="comment">// If we&#39;ve lost information then don&#39;t try to do anything with the M.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		<span class="comment">// It may have moved on and we can&#39;t be sure.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		if oldStatus == go122.ProcSyscallAbandoned {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			return curCtx, true, nil
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		<span class="comment">// Validate that the M we&#39;re stealing from is what we expect.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		mid := ThreadID(ev.args[2]) <span class="comment">// The M we&#39;re stealing from.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		if mid == curCtx.M {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			<span class="comment">// We&#39;re stealing from ourselves. This behaves like a ProcStop.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			if curCtx.P != pid {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>				return curCtx, false, fmt.Errorf(&#34;tried to self-steal proc %d (thread %d), but got proc %d instead&#34;, pid, mid, curCtx.P)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			newCtx.P = NoProc
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			return curCtx, true, nil
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re stealing from some other M.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		mState, ok := o.mStates[mid]
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		if !ok {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;stole proc from non-existent thread %d&#34;, mid)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		<span class="comment">// Make sure we&#39;re actually stealing the right P.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		if mState.p != pid {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;tried to steal proc %d from thread %d, but got proc %d instead&#34;, pid, mid, mState.p)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		<span class="comment">// Tell the M it has no P so it can proceed.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		<span class="comment">// This is safe because we know the P was in a syscall and</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		<span class="comment">// the other M must be trying to get out of the syscall.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		<span class="comment">// GoSyscallEndBlocked cannot advance until the corresponding</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		<span class="comment">// M loses its P.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		mState.p = NoProc
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	<span class="comment">// Handle goroutines.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	case go122.EvGoStatus:
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		gid := GoID(ev.args[0])
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		mid := ThreadID(ev.args[1])
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		status := go122.GoStatus(ev.args[2])
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		if int(status) &gt;= len(go122GoStatus2GoState) {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;invalid status for goroutine %d: %d&#34;, gid, status)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		oldState := go122GoStatus2GoState[status]
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		if s, ok := o.gStates[gid]; ok {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			if s.status != status {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>				return curCtx, false, fmt.Errorf(&#34;inconsistent status for goroutine %d: old %v vs. new %v&#34;, gid, s.status, status)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			s.seq = makeSeq(gen, 0) <span class="comment">// Reset seq.</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		} else if gen == o.initialGen {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			<span class="comment">// Set the state.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			o.gStates[gid] = &amp;gState{id: gid, status: status, seq: makeSeq(gen, 0)}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			oldState = GoUndetermined
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		} else {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;found goroutine status for new goroutine after the first generation: id=%v status=%v&#34;, gid, status)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		ev.extra(version.Go122)[0] = uint64(oldState) <span class="comment">// Smuggle in the old state for StateTransition.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		switch status {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		case go122.GoRunning:
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			<span class="comment">// Bind the goroutine to the new context, since it&#39;s running.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			newCtx.G = gid
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		case go122.GoSyscall:
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			if mid == NoThread {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				return curCtx, false, fmt.Errorf(&#34;found goroutine %d in syscall without a thread&#34;, gid)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			<span class="comment">// Is the syscall on this thread? If so, bind it to the context.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			<span class="comment">// Otherwise, we&#39;re talking about a G sitting in a syscall on an M.</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			<span class="comment">// Validate the named M.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			if mid == curCtx.M {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>				if gen != o.initialGen &amp;&amp; curCtx.G != gid {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>					<span class="comment">// If this isn&#39;t the first generation, we *must* have seen this</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>					<span class="comment">// binding occur already. Even if the G was blocked in a syscall</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>					<span class="comment">// for multiple generations since trace start, we would have seen</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>					<span class="comment">// a previous GoStatus event that bound the goroutine to an M.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>					return curCtx, false, fmt.Errorf(&#34;inconsistent thread for syscalling goroutine %d: thread has goroutine %d&#34;, gid, curCtx.G)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>				}
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>				newCtx.G = gid
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>				break
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>			}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>			<span class="comment">// Now we&#39;re talking about a thread and goroutine that have been</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>			<span class="comment">// blocked on a syscall for the entire generation. This case must</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			<span class="comment">// not have a P; the runtime makes sure that all Ps are traced at</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>			<span class="comment">// the beginning of a generation, which involves taking a P back</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			<span class="comment">// from every thread.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			ms, ok := o.mStates[mid]
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			if ok {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>				<span class="comment">// This M has been seen. That means we must have seen this</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>				<span class="comment">// goroutine go into a syscall on this thread at some point.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>				if ms.g != gid {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>					<span class="comment">// But the G on the M doesn&#39;t match. Something&#39;s wrong.</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>					return curCtx, false, fmt.Errorf(&#34;inconsistent thread for syscalling goroutine %d: thread has goroutine %d&#34;, gid, ms.g)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>				}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>				<span class="comment">// This case is just a Syscall-&gt;Syscall event, which needs to</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>				<span class="comment">// appear as having the G currently bound to this M.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>				curCtx.G = ms.g
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			} else if !ok {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>				<span class="comment">// The M hasn&#39;t been seen yet. That means this goroutine</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>				<span class="comment">// has just been sitting in a syscall on this M. Create</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>				<span class="comment">// a state for it.</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>				o.mStates[mid] = &amp;mState{g: gid, p: NoProc}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>				<span class="comment">// Don&#39;t set curCtx.G in this case because this event is the</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>				<span class="comment">// binding event (and curCtx represents the &#34;before&#34; state).</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>			}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>			<span class="comment">// Update the current context to the M we&#39;re talking about.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			curCtx.M = mid
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	case go122.EvGoCreate:
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		<span class="comment">// Goroutines must be created on a running P, but may or may not be created</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		<span class="comment">// by a running goroutine.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		reqs := event.SchedReqs{Thread: event.MustHave, Proc: event.MustHave, Goroutine: event.MayHave}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, reqs); err != nil {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		<span class="comment">// If we have a goroutine, it must be running.</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		if state, ok := o.gStates[curCtx.G]; ok &amp;&amp; state.status != go122.GoRunning {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;%s event for goroutine that&#39;s not %s&#34;, go122.EventString(typ), GoRunning)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		<span class="comment">// This goroutine created another. Add a state for it.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		newgid := GoID(ev.args[0])
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		if _, ok := o.gStates[newgid]; ok {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;tried to create goroutine (%v) that already exists&#34;, newgid)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		o.gStates[newgid] = &amp;gState{id: newgid, status: go122.GoRunnable, seq: makeSeq(gen, 0)}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	case go122.EvGoDestroy, go122.EvGoStop, go122.EvGoBlock:
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		<span class="comment">// These are goroutine events that all require an active running</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		<span class="comment">// goroutine on some thread. They must *always* be advance-able,</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		<span class="comment">// since running goroutines are bound to their M.</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		state, ok := o.gStates[curCtx.G]
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		if !ok {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;event %s for goroutine (%v) that doesn&#39;t exist&#34;, go122.EventString(typ), curCtx.G)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		if state.status != go122.GoRunning {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;%s event for goroutine that&#39;s not %s&#34;, go122.EventString(typ), GoRunning)
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		<span class="comment">// Handle each case slightly differently; we just group them together</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		<span class="comment">// because they have shared preconditions.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		switch typ {
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		case go122.EvGoDestroy:
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			<span class="comment">// This goroutine is exiting itself.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			delete(o.gStates, curCtx.G)
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			newCtx.G = NoGoroutine
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		case go122.EvGoStop:
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			<span class="comment">// Goroutine stopped (yielded). It&#39;s runnable but not running on this M.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			state.status = go122.GoRunnable
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>			newCtx.G = NoGoroutine
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		case go122.EvGoBlock:
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>			<span class="comment">// Goroutine blocked. It&#39;s waiting now and not running on this M.</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			state.status = go122.GoWaiting
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			newCtx.G = NoGoroutine
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		}
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	case go122.EvGoStart:
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		gid := GoID(ev.args[0])
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		seq := makeSeq(gen, ev.args[1])
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		state, ok := o.gStates[gid]
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		if !ok || state.status != go122.GoRunnable || !seq.succeeds(state.seq) {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			<span class="comment">// We can&#39;t make an inference as to whether this is bad. We could just be seeing</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			<span class="comment">// a GoStart on a different M before the goroutine was created, before it had its</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			<span class="comment">// state emitted, or before we got to the right point in the trace yet.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			return curCtx, false, nil
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		<span class="comment">// We can advance this goroutine. Check some invariants.</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		reqs := event.SchedReqs{Thread: event.MustHave, Proc: event.MustHave, Goroutine: event.MustNotHave}
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, reqs); err != nil {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		state.status = go122.GoRunning
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		state.seq = seq
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		newCtx.G = gid
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	case go122.EvGoUnblock:
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		<span class="comment">// N.B. These both reference the goroutine to unblock, not the current goroutine.</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		gid := GoID(ev.args[0])
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		seq := makeSeq(gen, ev.args[1])
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		state, ok := o.gStates[gid]
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		if !ok || state.status != go122.GoWaiting || !seq.succeeds(state.seq) {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			<span class="comment">// We can&#39;t make an inference as to whether this is bad. We could just be seeing</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			<span class="comment">// a GoUnblock on a different M before the goroutine was created and blocked itself,</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			<span class="comment">// before it had its state emitted, or before we got to the right point in the trace yet.</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			return curCtx, false, nil
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		state.status = go122.GoRunnable
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		state.seq = seq
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		<span class="comment">// N.B. No context to validate. Basically anything can unblock</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		<span class="comment">// a goroutine (e.g. sysmon).</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	case go122.EvGoSyscallBegin:
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		<span class="comment">// Entering a syscall requires an active running goroutine with a</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		<span class="comment">// proc on some thread. It is always advancable.</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		state, ok := o.gStates[curCtx.G]
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		if !ok {
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;event %s for goroutine (%v) that doesn&#39;t exist&#34;, go122.EventString(typ), curCtx.G)
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		}
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		if state.status != go122.GoRunning {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;%s event for goroutine that&#39;s not %s&#34;, go122.EventString(typ), GoRunning)
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		}
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		<span class="comment">// Goroutine entered a syscall. It&#39;s still running on this P and M.</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		state.status = go122.GoSyscall
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		pState, ok := o.pStates[curCtx.P]
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		if !ok {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;uninitialized proc %d found during %s&#34;, curCtx.P, go122.EventString(typ))
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		pState.status = go122.ProcSyscall
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		<span class="comment">// Validate the P sequence number on the event and advance it.</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		<span class="comment">// We have a P sequence number for what is supposed to be a goroutine event</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		<span class="comment">// so that we can correctly model P stealing. Without this sequence number here,</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		<span class="comment">// the syscall from which a ProcSteal event is stealing can be ambiguous in the</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		<span class="comment">// face of broken timestamps. See the go122-syscall-steal-proc-ambiguous test for</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		<span class="comment">// more details.</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		<span class="comment">// Note that because this sequence number only exists as a tool for disambiguation,</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		<span class="comment">// we can enforce that we have the right sequence number at this point; we don&#39;t need</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		<span class="comment">// to back off and see if any other events will advance. This is a running P.</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		pSeq := makeSeq(gen, ev.args[0])
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		if !pSeq.succeeds(pState.seq) {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;failed to advance %s: can&#39;t make sequence: %s -&gt; %s&#34;, go122.EventString(typ), pState.seq, pSeq)
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		pState.seq = pSeq
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	case go122.EvGoSyscallEnd:
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		<span class="comment">// This event is always advance-able because it happens on the same</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		<span class="comment">// thread that EvGoSyscallStart happened, and the goroutine can&#39;t leave</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		<span class="comment">// that thread until its done.</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		state, ok := o.gStates[curCtx.G]
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		if !ok {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;event %s for goroutine (%v) that doesn&#39;t exist&#34;, go122.EventString(typ), curCtx.G)
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		if state.status != go122.GoSyscall {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;%s event for goroutine that&#39;s not %s&#34;, go122.EventString(typ), GoRunning)
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		state.status = go122.GoRunning
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		<span class="comment">// Transfer the P back to running from syscall.</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		pState, ok := o.pStates[curCtx.P]
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		if !ok {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;uninitialized proc %d found during %s&#34;, curCtx.P, go122.EventString(typ))
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		if pState.status != go122.ProcSyscall {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;expected proc %d in state %v, but got %v instead&#34;, curCtx.P, go122.ProcSyscall, pState.status)
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		pState.status = go122.ProcRunning
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	case go122.EvGoSyscallEndBlocked:
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		<span class="comment">// This event becomes advanceable when its P is not in a syscall state</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		<span class="comment">// (lack of a P altogether is also acceptable for advancing).</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		<span class="comment">// The transfer out of ProcSyscall can happen either voluntarily via</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		<span class="comment">// ProcStop or involuntarily via ProcSteal. We may also acquire a new P</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		<span class="comment">// before we get here (after the transfer out) but that&#39;s OK: that new</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		<span class="comment">// P won&#39;t be in the ProcSyscall state anymore.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		<span class="comment">// Basically: while we have a preemptible P, don&#39;t advance, because we</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		<span class="comment">// *know* from the event that we&#39;re going to lose it at some point during</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		<span class="comment">// the syscall. We shouldn&#39;t advance until that happens.</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		if curCtx.P != NoProc {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			pState, ok := o.pStates[curCtx.P]
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			if !ok {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>				return curCtx, false, fmt.Errorf(&#34;uninitialized proc %d found during %s&#34;, curCtx.P, go122.EventString(typ))
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>			}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			if pState.status == go122.ProcSyscall {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>				return curCtx, false, nil
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		<span class="comment">// As mentioned above, we may have a P here if we ProcStart</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		<span class="comment">// before this event.</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.SchedReqs{Thread: event.MustHave, Proc: event.MayHave, Goroutine: event.MustHave}); err != nil {
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		state, ok := o.gStates[curCtx.G]
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		if !ok {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;event %s for goroutine (%v) that doesn&#39;t exist&#34;, go122.EventString(typ), curCtx.G)
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		}
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		if state.status != go122.GoSyscall {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;%s event for goroutine that&#39;s not %s&#34;, go122.EventString(typ), GoRunning)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		newCtx.G = NoGoroutine
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		state.status = go122.GoRunnable
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	case go122.EvGoCreateSyscall:
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		<span class="comment">// This event indicates that a goroutine is effectively</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		<span class="comment">// being created out of a cgo callback. Such a goroutine</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		<span class="comment">// is &#39;created&#39; in the syscall state.</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.SchedReqs{Thread: event.MustHave, Proc: event.MayHave, Goroutine: event.MustNotHave}); err != nil {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		}
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		<span class="comment">// This goroutine is effectively being created. Add a state for it.</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		newgid := GoID(ev.args[0])
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		if _, ok := o.gStates[newgid]; ok {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;tried to create goroutine (%v) in syscall that already exists&#34;, newgid)
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		o.gStates[newgid] = &amp;gState{id: newgid, status: go122.GoSyscall, seq: makeSeq(gen, 0)}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		<span class="comment">// Goroutine is executing. Bind it to the context.</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		newCtx.G = newgid
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	case go122.EvGoDestroySyscall:
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		<span class="comment">// This event indicates that a goroutine created for a</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		<span class="comment">// cgo callback is disappearing, either because the callback</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		<span class="comment">// ending or the C thread that called it is being destroyed.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		<span class="comment">// Also, treat this as if we lost our P too.</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		<span class="comment">// The thread ID may be reused by the platform and we&#39;ll get</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		<span class="comment">// really confused if we try to steal the P is this is running</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		<span class="comment">// with later. The new M with the same ID could even try to</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		<span class="comment">// steal back this P from itself!</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		<span class="comment">// The runtime is careful to make sure that any GoCreateSyscall</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		<span class="comment">// event will enter the runtime emitting events for reacquiring a P.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		<span class="comment">// Note: we might have a P here. The P might not be released</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		<span class="comment">// eagerly by the runtime, and it might get stolen back later</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		<span class="comment">// (or never again, if the program is going to exit).</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.SchedReqs{Thread: event.MustHave, Proc: event.MayHave, Goroutine: event.MustHave}); err != nil {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		<span class="comment">// Check to make sure the goroutine exists in the right state.</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		state, ok := o.gStates[curCtx.G]
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		if !ok {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;event %s for goroutine (%v) that doesn&#39;t exist&#34;, go122.EventString(typ), curCtx.G)
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		}
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		if state.status != go122.GoSyscall {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;%s event for goroutine that&#39;s not %v&#34;, go122.EventString(typ), GoSyscall)
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		<span class="comment">// This goroutine is exiting itself.</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		delete(o.gStates, curCtx.G)
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		newCtx.G = NoGoroutine
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		<span class="comment">// If we have a proc, then we&#39;re dissociating from it now. See the comment at the top of the case.</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		if curCtx.P != NoProc {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>			pState, ok := o.pStates[curCtx.P]
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>			if !ok {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>				return curCtx, false, fmt.Errorf(&#34;found invalid proc %d during %s&#34;, curCtx.P, go122.EventString(typ))
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>			}
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>			if pState.status != go122.ProcSyscall {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>				return curCtx, false, fmt.Errorf(&#34;proc %d in unexpected state %s during %s&#34;, curCtx.P, pState.status, go122.EventString(typ))
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>			}
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>			<span class="comment">// See the go122-create-syscall-reuse-thread-id test case for more details.</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>			pState.status = go122.ProcSyscallAbandoned
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>			newCtx.P = NoProc
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>			<span class="comment">// Queue an extra self-ProcSteal event.</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>			o.extraEvent = Event{
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>				table: evt,
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>				ctx:   curCtx,
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>				base: baseEvent{
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>					typ:  go122.EvProcSteal,
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>					time: ev.time,
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>				},
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>			o.extraEvent.base.args[0] = uint64(curCtx.P)
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>			o.extraEvent.base.extra(version.Go122)[0] = uint64(go122.ProcSyscall)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	<span class="comment">// Handle tasks. Tasks are interesting because:</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	<span class="comment">// - There&#39;s no Begin event required to reference a task.</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	<span class="comment">// - End for a particular task ID can appear multiple times.</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	<span class="comment">// As a result, there&#39;s very little to validate. The only</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	<span class="comment">// thing we have to be sure of is that a task didn&#39;t begin</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	<span class="comment">// after it had already begun. Task IDs are allowed to be</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	<span class="comment">// reused, so we don&#39;t care about a Begin after an End.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	case go122.EvUserTaskBegin:
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		id := TaskID(ev.args[0])
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		if _, ok := o.activeTasks[id]; ok {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;task ID conflict: %d&#34;, id)
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		}
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		<span class="comment">// Get the parent ID, but don&#39;t validate it. There&#39;s no guarantee</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		<span class="comment">// we actually have information on whether it&#39;s active.</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		parentID := TaskID(ev.args[1])
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		if parentID == BackgroundTask {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>			<span class="comment">// Note: a value of 0 here actually means no parent, *not* the</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>			<span class="comment">// background task. Automatic background task attachment only</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>			<span class="comment">// applies to regions.</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>			parentID = NoTask
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			ev.args[1] = uint64(NoTask)
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		}
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		<span class="comment">// Validate the name and record it. We&#39;ll need to pass it through to</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		<span class="comment">// EvUserTaskEnd.</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		nameID := stringID(ev.args[2])
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		name, ok := evt.strings.get(nameID)
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		if !ok {
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;invalid string ID %v for %v event&#34;, nameID, typ)
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		}
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		o.activeTasks[id] = taskState{name: name, parentID: parentID}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		return curCtx, true, validateCtx(curCtx, event.UserGoReqs)
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	case go122.EvUserTaskEnd:
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>		id := TaskID(ev.args[0])
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		if ts, ok := o.activeTasks[id]; ok {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>			<span class="comment">// Smuggle the task info. This may happen in a different generation,</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>			<span class="comment">// which may not have the name in its string table. Add it to the extra</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>			<span class="comment">// strings table so we can look it up later.</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>			ev.extra(version.Go122)[0] = uint64(ts.parentID)
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>			ev.extra(version.Go122)[1] = uint64(evt.addExtraString(ts.name))
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			delete(o.activeTasks, id)
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		} else {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>			<span class="comment">// Explicitly clear the task info.</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>			ev.extra(version.Go122)[0] = uint64(NoTask)
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>			ev.extra(version.Go122)[1] = uint64(evt.addExtraString(&#34;&#34;))
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		}
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		return curCtx, true, validateCtx(curCtx, event.UserGoReqs)
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	<span class="comment">// Handle user regions.</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	case go122.EvUserRegionBegin:
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		}
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>		tid := TaskID(ev.args[0])
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		nameID := stringID(ev.args[1])
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		name, ok := evt.strings.get(nameID)
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		if !ok {
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;invalid string ID %v for %v event&#34;, nameID, typ)
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		gState, ok := o.gStates[curCtx.G]
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		if !ok {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;encountered EvUserRegionBegin without known state for current goroutine %d&#34;, curCtx.G)
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		}
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		if err := gState.beginRegion(userRegion{tid, name}); err != nil {
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		}
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	case go122.EvUserRegionEnd:
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>		tid := TaskID(ev.args[0])
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		nameID := stringID(ev.args[1])
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		name, ok := evt.strings.get(nameID)
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		if !ok {
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;invalid string ID %v for %v event&#34;, nameID, typ)
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		}
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		gState, ok := o.gStates[curCtx.G]
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		if !ok {
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;encountered EvUserRegionEnd without known state for current goroutine %d&#34;, curCtx.G)
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		}
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		if err := gState.endRegion(userRegion{tid, name}); err != nil {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		}
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	<span class="comment">// Handle the GC mark phase.</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	<span class="comment">// We have sequence numbers for both start and end because they</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	<span class="comment">// can happen on completely different threads. We want an explicit</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	<span class="comment">// partial order edge between start and end here, otherwise we&#39;re</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	<span class="comment">// relying entirely on timestamps to make sure we don&#39;t advance a</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	<span class="comment">// GCEnd for a _different_ GC cycle if timestamps are wildly broken.</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	case go122.EvGCActive:
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		seq := ev.args[0]
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>		if gen == o.initialGen {
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>			if o.gcState != gcUndetermined {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>				return curCtx, false, fmt.Errorf(&#34;GCActive in the first generation isn&#39;t first GC event&#34;)
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			}
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			o.gcSeq = seq
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>			o.gcState = gcRunning
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>			return curCtx, true, nil
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		}
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>		if seq != o.gcSeq+1 {
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>			<span class="comment">// This is not the right GC cycle.</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			return curCtx, false, nil
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		}
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		if o.gcState != gcRunning {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;encountered GCActive while GC was not in progress&#34;)
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		}
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>		o.gcSeq = seq
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		}
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	case go122.EvGCBegin:
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		seq := ev.args[0]
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		if o.gcState == gcUndetermined {
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>			o.gcSeq = seq
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>			o.gcState = gcRunning
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>			return curCtx, true, nil
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		}
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		if seq != o.gcSeq+1 {
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>			<span class="comment">// This is not the right GC cycle.</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>			return curCtx, false, nil
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		}
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		if o.gcState == gcRunning {
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;encountered GCBegin while GC was already in progress&#34;)
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>		}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		o.gcSeq = seq
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		o.gcState = gcRunning
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>		}
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	case go122.EvGCEnd:
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>		seq := ev.args[0]
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>		if seq != o.gcSeq+1 {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>			<span class="comment">// This is not the right GC cycle.</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>			return curCtx, false, nil
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>		}
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		if o.gcState == gcNotRunning {
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;encountered GCEnd when GC was not in progress&#34;)
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>		}
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		if o.gcState == gcUndetermined {
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;encountered GCEnd when GC was in an undetermined state&#34;)
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>		}
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		o.gcSeq = seq
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		o.gcState = gcNotRunning
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		}
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	<span class="comment">// Handle simple instantaneous events that require a G.</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	case go122.EvGoLabel, go122.EvProcsChange, go122.EvUserLog:
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		}
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	<span class="comment">// Handle allocation states, which don&#39;t require a G.</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	case go122.EvHeapAlloc, go122.EvHeapGoal:
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.SchedReqs{Thread: event.MustHave, Proc: event.MustHave, Goroutine: event.MayHave}); err != nil {
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		}
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	<span class="comment">// Handle sweep, which is bound to a P and doesn&#39;t require a G.</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	case go122.EvGCSweepBegin:
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.SchedReqs{Thread: event.MustHave, Proc: event.MustHave, Goroutine: event.MayHave}); err != nil {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>		}
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		if err := o.pStates[curCtx.P].beginRange(makeRangeType(typ, 0)); err != nil {
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>	case go122.EvGCSweepActive:
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		pid := ProcID(ev.args[0])
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		<span class="comment">// N.B. In practice Ps can&#39;t block while they&#39;re sweeping, so this can only</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		<span class="comment">// ever reference curCtx.P. However, be lenient about this like we are with</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>		<span class="comment">// GCMarkAssistActive; there&#39;s no reason the runtime couldn&#39;t change to block</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>		<span class="comment">// in the middle of a sweep.</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>		pState, ok := o.pStates[pid]
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>		if !ok {
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;encountered GCSweepActive for unknown proc %d&#34;, pid)
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		}
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		if err := pState.activeRange(makeRangeType(typ, 0), gen == o.initialGen); err != nil {
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>		}
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	case go122.EvGCSweepEnd:
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.SchedReqs{Thread: event.MustHave, Proc: event.MustHave, Goroutine: event.MayHave}); err != nil {
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>		}
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>		_, err := o.pStates[curCtx.P].endRange(typ)
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		if err != nil {
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		}
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	<span class="comment">// Handle special goroutine-bound event ranges.</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	case go122.EvSTWBegin, go122.EvGCMarkAssistBegin:
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>		}
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		desc := stringID(0)
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>		if typ == go122.EvSTWBegin {
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>			desc = stringID(ev.args[0])
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>		}
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>		gState, ok := o.gStates[curCtx.G]
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		if !ok {
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;encountered event of type %d without known state for current goroutine %d&#34;, typ, curCtx.G)
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		}
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		if err := gState.beginRange(makeRangeType(typ, desc)); err != nil {
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		}
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>	case go122.EvGCMarkAssistActive:
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		gid := GoID(ev.args[0])
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>		<span class="comment">// N.B. Like GoStatus, this can happen at any time, because it can</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>		<span class="comment">// reference a non-running goroutine. Don&#39;t check anything about the</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>		<span class="comment">// current scheduler context.</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		gState, ok := o.gStates[gid]
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		if !ok {
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;uninitialized goroutine %d found during %s&#34;, gid, go122.EventString(typ))
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		}
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		if err := gState.activeRange(makeRangeType(typ, 0), gen == o.initialGen); err != nil {
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		}
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	case go122.EvSTWEnd, go122.EvGCMarkAssistEnd:
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		if err := validateCtx(curCtx, event.UserGoReqs); err != nil {
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		}
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		gState, ok := o.gStates[curCtx.G]
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		if !ok {
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>			return curCtx, false, fmt.Errorf(&#34;encountered event of type %d without known state for current goroutine %d&#34;, typ, curCtx.G)
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		}
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		desc, err := gState.endRange(typ)
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		if err != nil {
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>			return curCtx, false, err
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>		if typ == go122.EvSTWEnd {
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>			<span class="comment">// Smuggle the kind into the event.</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			<span class="comment">// Don&#39;t use ev.extra here so we have symmetry with STWBegin.</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>			ev.args[0] = uint64(desc)
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		}
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>		return curCtx, true, nil
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	}
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	return curCtx, false, fmt.Errorf(&#34;bad event type found while ordering: %v&#34;, ev.typ)
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>}
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span><span class="comment">// schedCtx represents the scheduling resources associated with an event.</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>type schedCtx struct {
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	G GoID
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	P ProcID
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	M ThreadID
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>}
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span><span class="comment">// validateCtx ensures that ctx conforms to some reqs, returning an error if</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span><span class="comment">// it doesn&#39;t.</span>
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>func validateCtx(ctx schedCtx, reqs event.SchedReqs) error {
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	<span class="comment">// Check thread requirements.</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	if reqs.Thread == event.MustHave &amp;&amp; ctx.M == NoThread {
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;expected a thread but didn&#39;t have one&#34;)
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	} else if reqs.Thread == event.MustNotHave &amp;&amp; ctx.M != NoThread {
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;expected no thread but had one&#34;)
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	}
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	<span class="comment">// Check proc requirements.</span>
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	if reqs.Proc == event.MustHave &amp;&amp; ctx.P == NoProc {
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;expected a proc but didn&#39;t have one&#34;)
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	} else if reqs.Proc == event.MustNotHave &amp;&amp; ctx.P != NoProc {
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;expected no proc but had one&#34;)
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	}
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	<span class="comment">// Check goroutine requirements.</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	if reqs.Goroutine == event.MustHave &amp;&amp; ctx.G == NoGoroutine {
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;expected a goroutine but didn&#39;t have one&#34;)
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	} else if reqs.Goroutine == event.MustNotHave &amp;&amp; ctx.G != NoGoroutine {
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;expected no goroutine but had one&#34;)
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	}
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	return nil
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>}
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span><span class="comment">// gcState is a trinary variable for the current state of the GC.</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span><span class="comment">// The third state besides &#34;enabled&#34; and &#34;disabled&#34; is &#34;undetermined.&#34;</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>type gcState uint8
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>const (
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	gcUndetermined gcState = iota
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	gcNotRunning
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	gcRunning
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>)
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>
<span id="L897" class="ln">   897&nbsp;&nbsp;</span><span class="comment">// String returns a human-readable string for the GC state.</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>func (s gcState) String() string {
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	switch s {
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>	case gcUndetermined:
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		return &#34;Undetermined&#34;
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	case gcNotRunning:
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>		return &#34;NotRunning&#34;
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	case gcRunning:
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>		return &#34;Running&#34;
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	}
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	return &#34;Bad&#34;
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>}
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span><span class="comment">// userRegion represents a unique user region when attached to some gState.</span>
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>type userRegion struct {
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	<span class="comment">// name must be a resolved string because the string ID for the same</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	<span class="comment">// string may change across generations, but we care about checking</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	<span class="comment">// the value itself.</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	taskID TaskID
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>	name   string
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>}
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span><span class="comment">// rangeType is a way to classify special ranges of time.</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L921" class="ln">   921&nbsp;&nbsp;</span><span class="comment">// These typically correspond 1:1 with &#34;Begin&#34; events, but</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span><span class="comment">// they may have an optional subtype that describes the range</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span><span class="comment">// in more detail.</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>type rangeType struct {
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	typ  event.Type <span class="comment">// &#34;Begin&#34; event.</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	desc stringID   <span class="comment">// Optional subtype.</span>
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>}
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span><span class="comment">// makeRangeType constructs a new rangeType.</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>func makeRangeType(typ event.Type, desc stringID) rangeType {
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>	if styp := go122.Specs()[typ].StartEv; styp != go122.EvNone {
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>		typ = styp
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>	}
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	return rangeType{typ, desc}
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>}
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span><span class="comment">// gState is the state of a goroutine at a point in the trace.</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>type gState struct {
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	id     GoID
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	status go122.GoStatus
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>	seq    seqCounter
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>	<span class="comment">// regions are the active user regions for this goroutine.</span>
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>	regions []userRegion
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>	<span class="comment">// rangeState is the state of special time ranges bound to this goroutine.</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>	rangeState
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>}
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>
<span id="L950" class="ln">   950&nbsp;&nbsp;</span><span class="comment">// beginRegion starts a user region on the goroutine.</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>func (s *gState) beginRegion(r userRegion) error {
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>	s.regions = append(s.regions, r)
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>	return nil
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>}
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span><span class="comment">// endRegion ends a user region on the goroutine.</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>func (s *gState) endRegion(r userRegion) error {
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>	if len(s.regions) == 0 {
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		<span class="comment">// We do not know about regions that began before tracing started.</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>		return nil
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	}
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>	if next := s.regions[len(s.regions)-1]; next != r {
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;misuse of region in goroutine %v: region end %v when the inner-most active region start event is %v&#34;, s.id, r, next)
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	}
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	s.regions = s.regions[:len(s.regions)-1]
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>	return nil
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>}
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span><span class="comment">// pState is the state of a proc at a point in the trace.</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>type pState struct {
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>	id     ProcID
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>	status go122.ProcStatus
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	seq    seqCounter
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>	<span class="comment">// rangeState is the state of special time ranges bound to this proc.</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	rangeState
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>}
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span><span class="comment">// mState is the state of a thread at a point in the trace.</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>type mState struct {
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	g GoID   <span class="comment">// Goroutine bound to this M. (The goroutine&#39;s state is Executing.)</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	p ProcID <span class="comment">// Proc bound to this M. (The proc&#39;s state is Executing.)</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>}
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>
<span id="L985" class="ln">   985&nbsp;&nbsp;</span><span class="comment">// rangeState represents the state of special time ranges.</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>type rangeState struct {
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	<span class="comment">// inFlight contains the rangeTypes of any ranges bound to a resource.</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>	inFlight []rangeType
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>}
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>
<span id="L991" class="ln">   991&nbsp;&nbsp;</span><span class="comment">// beginRange begins a special range in time on the goroutine.</span>
<span id="L992" class="ln">   992&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span><span class="comment">// Returns an error if the range is already in progress.</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>func (s *rangeState) beginRange(typ rangeType) error {
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>	if s.hasRange(typ) {
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;discovered event already in-flight for when starting event %v&#34;, go122.Specs()[typ.typ].Name)
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	}
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	s.inFlight = append(s.inFlight, typ)
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	return nil
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>}
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span><span class="comment">// activeRange marks special range in time on the goroutine as active in the</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span><span class="comment">// initial generation, or confirms that it is indeed active in later generations.</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>func (s *rangeState) activeRange(typ rangeType, isInitialGen bool) error {
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	if isInitialGen {
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>		if s.hasRange(typ) {
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;found named active range already in first gen: %v&#34;, typ)
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>		}
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>		s.inFlight = append(s.inFlight, typ)
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>	} else if !s.hasRange(typ) {
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;resource is missing active range: %v %v&#34;, go122.Specs()[typ.typ].Name, s.inFlight)
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	}
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	return nil
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>}
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span><span class="comment">// hasRange returns true if a special time range on the goroutine as in progress.</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>func (s *rangeState) hasRange(typ rangeType) bool {
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	for _, ftyp := range s.inFlight {
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>		if ftyp == typ {
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>			return true
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		}
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	}
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	return false
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>}
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span><span class="comment">// endsRange ends a special range in time on the goroutine.</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span><span class="comment">// This must line up with the start event type  of the range the goroutine is currently in.</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>func (s *rangeState) endRange(typ event.Type) (stringID, error) {
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>	st := go122.Specs()[typ].StartEv
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	idx := -1
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	for i, r := range s.inFlight {
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		if r.typ == st {
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>			idx = i
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>			break
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>		}
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>	}
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	if idx &lt; 0 {
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>		return 0, fmt.Errorf(&#34;tried to end event %v, but not in-flight&#34;, go122.Specs()[st].Name)
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>	}
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>	<span class="comment">// Swap remove.</span>
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>	desc := s.inFlight[idx].desc
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>	s.inFlight[idx], s.inFlight[len(s.inFlight)-1] = s.inFlight[len(s.inFlight)-1], s.inFlight[idx]
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>	s.inFlight = s.inFlight[:len(s.inFlight)-1]
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>	return desc, nil
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>}
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span><span class="comment">// seqCounter represents a global sequence counter for a resource.</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>type seqCounter struct {
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>	gen uint64 <span class="comment">// The generation for the local sequence counter seq.</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>	seq uint64 <span class="comment">// The sequence number local to the generation.</span>
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>}
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span><span class="comment">// makeSeq creates a new seqCounter.</span>
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>func makeSeq(gen, seq uint64) seqCounter {
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>	return seqCounter{gen: gen, seq: seq}
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>}
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span><span class="comment">// succeeds returns true if a is the immediate successor of b.</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>func (a seqCounter) succeeds(b seqCounter) bool {
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>	return a.gen == b.gen &amp;&amp; a.seq == b.seq+1
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>}
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span><span class="comment">// String returns a debug string representation of the seqCounter.</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>func (c seqCounter) String() string {
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;%d (gen=%d)&#34;, c.seq, c.gen)
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>}
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>func dumpOrdering(order *ordering) string {
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	var sb strings.Builder
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	for id, state := range order.gStates {
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>		fmt.Fprintf(&amp;sb, &#34;G %d [status=%s seq=%s]\n&#34;, id, state.status, state.seq)
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	}
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	fmt.Fprintln(&amp;sb)
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	for id, state := range order.pStates {
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>		fmt.Fprintf(&amp;sb, &#34;P %d [status=%s seq=%s]\n&#34;, id, state.status, state.seq)
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	}
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>	fmt.Fprintln(&amp;sb)
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	for id, state := range order.mStates {
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>		fmt.Fprintf(&amp;sb, &#34;M %d [g=%d p=%d]\n&#34;, id, state.g, state.p)
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>	}
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>	fmt.Fprintln(&amp;sb)
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>	fmt.Fprintf(&amp;sb, &#34;GC %d %s\n&#34;, order.gcSeq, order.gcState)
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>	return sb.String()
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>}
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span><span class="comment">// taskState represents an active task.</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>type taskState struct {
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>	<span class="comment">// name is the type of the active task.</span>
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>	name string
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>	<span class="comment">// parentID is the parent ID of the active task.</span>
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>	parentID TaskID
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>}
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>
</pre><p><a href="order.go?m=text">View as plain text</a></p>

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
