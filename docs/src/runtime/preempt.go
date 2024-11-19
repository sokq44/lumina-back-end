<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/preempt.go - Go Documentation Server</title>

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
<a href="preempt.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">preempt.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2019 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Goroutine preemption</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// A goroutine can be preempted at any safe-point. Currently, there</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// are a few categories of safe-points:</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// 1. A blocked safe-point occurs for the duration that a goroutine is</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//    descheduled, blocked on synchronization, or in a system call.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// 2. Synchronous safe-points occur when a running goroutine checks</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//    for a preemption request.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// 3. Asynchronous safe-points occur at any instruction in user code</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//    where the goroutine can be safely paused and a conservative</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//    stack and register scan can find stack roots. The runtime can</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//    stop a goroutine at an async safe-point using a signal.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// At both blocked and synchronous safe-points, a goroutine&#39;s CPU</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// state is minimal and the garbage collector has complete information</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// about its entire stack. This makes it possible to deschedule a</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// goroutine with minimal space, and to precisely scan a goroutine&#39;s</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// stack.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// Synchronous safe-points are implemented by overloading the stack</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// bound check in function prologues. To preempt a goroutine at the</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// next synchronous safe-point, the runtime poisons the goroutine&#39;s</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// stack bound to a value that will cause the next stack bound check</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// to fail and enter the stack growth implementation, which will</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// detect that it was actually a preemption and redirect to preemption</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// handling.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// Preemption at asynchronous safe-points is implemented by suspending</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// the thread using an OS mechanism (e.g., signals) and inspecting its</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// state to determine if the goroutine was at an asynchronous</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// safe-point. Since the thread suspension itself is generally</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// asynchronous, it also checks if the running goroutine wants to be</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// preempted, since this could have changed. If all conditions are</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// satisfied, it adjusts the signal context to make it look like the</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// signaled thread just called asyncPreempt and resumes the thread.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// asyncPreempt spills all registers and enters the scheduler.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// (An alternative would be to preempt in the signal handler itself.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// This would let the OS save and restore the register state and the</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// runtime would only need to know how to extract potentially</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// pointer-containing registers from the signal context. However, this</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// would consume an M for every preempted G, and the scheduler itself</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// is not designed to run from a signal handler, as it tends to</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// allocate memory and start threads in the preemption path.)</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>package runtime
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>import (
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>type suspendGState struct {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	g *g
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// dead indicates the goroutine was not suspended because it</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// is dead. This goroutine could be reused after the dead</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// state was observed, so the caller must not assume that it</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// remains dead.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	dead bool
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// stopped indicates that this suspendG transitioned the G to</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// _Gwaiting via g.preemptStop and thus is responsible for</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// readying it when done.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	stopped bool
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// suspendG suspends goroutine gp at a safe-point and returns the</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// state of the suspended goroutine. The caller gets read access to</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// the goroutine until it calls resumeG.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// It is safe for multiple callers to attempt to suspend the same</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// goroutine at the same time. The goroutine may execute between</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// subsequent successful suspend operations. The current</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// implementation grants exclusive access to the goroutine, and hence</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// multiple callers will serialize. However, the intent is to grant</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// shared read access, so please don&#39;t depend on exclusive access.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// This must be called from the system stack and the user goroutine on</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// the current M (if any) must be in a preemptible state. This</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// prevents deadlocks where two goroutines attempt to suspend each</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// other and both are in non-preemptible states. There are other ways</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// to resolve this deadlock, but this seems simplest.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// TODO(austin): What if we instead required this to be called from a</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// user goroutine? Then we could deschedule the goroutine while</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// waiting instead of blocking the thread. If two goroutines tried to</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// suspend each other, one of them would win and the other wouldn&#39;t</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// complete the suspend until it was resumed. We would have to be</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// careful that they couldn&#39;t actually queue up suspend for each other</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// and then both be suspended. This would also avoid the need for a</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// kernel context switch in the synchronous case because we could just</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// directly schedule the waiter. The context switch is unavoidable in</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// the signal case.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>func suspendG(gp *g) suspendGState {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	if mp := getg().m; mp.curg != nil &amp;&amp; readgstatus(mp.curg) == _Grunning {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">// Since we&#39;re on the system stack of this M, the user</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		<span class="comment">// G is stuck at an unsafe point. If another goroutine</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		<span class="comment">// were to try to preempt m.curg, it could deadlock.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		throw(&#34;suspendG from non-preemptible goroutine&#34;)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// See https://golang.org/cl/21503 for justification of the yield delay.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	const yieldDelay = 10 * 1000
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	var nextYield int64
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// Drive the goroutine to a preemption point.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	stopped := false
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	var asyncM *m
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	var asyncGen uint32
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	var nextPreemptM int64
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	for i := 0; ; i++ {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		switch s := readgstatus(gp); s {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		default:
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			if s&amp;_Gscan != 0 {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>				<span class="comment">// Someone else is suspending it. Wait</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>				<span class="comment">// for them to finish.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>				<span class="comment">// TODO: It would be nicer if we could</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>				<span class="comment">// coalesce suspends.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>				break
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			dumpgstatus(gp)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			throw(&#34;invalid g status&#34;)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		case _Gdead:
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			<span class="comment">// Nothing to suspend.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			<span class="comment">// preemptStop may need to be cleared, but</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			<span class="comment">// doing that here could race with goroutine</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			<span class="comment">// reuse. Instead, goexit0 clears it.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			return suspendGState{dead: true}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		case _Gcopystack:
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			<span class="comment">// The stack is being copied. We need to wait</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			<span class="comment">// until this is done.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		case _Gpreempted:
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			<span class="comment">// We (or someone else) suspended the G. Claim</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			<span class="comment">// ownership of it by transitioning it to</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			<span class="comment">// _Gwaiting.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			if !casGFromPreempted(gp, _Gpreempted, _Gwaiting) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>				break
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			<span class="comment">// We stopped the G, so we have to ready it later.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			stopped = true
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			s = _Gwaiting
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			fallthrough
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		case _Grunnable, _Gsyscall, _Gwaiting:
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			<span class="comment">// Claim goroutine by setting scan bit.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			<span class="comment">// This may race with execution or readying of gp.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			<span class="comment">// The scan bit keeps it from transition state.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			if !castogscanstatus(gp, s, s|_Gscan) {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>				break
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			<span class="comment">// Clear the preemption request. It&#39;s safe to</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			<span class="comment">// reset the stack guard because we hold the</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			<span class="comment">// _Gscan bit and thus own the stack.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			gp.preemptStop = false
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			gp.preempt = false
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			gp.stackguard0 = gp.stack.lo + stackGuard
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			<span class="comment">// The goroutine was already at a safe-point</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			<span class="comment">// and we&#39;ve now locked that in.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			<span class="comment">// TODO: It would be much better if we didn&#39;t</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			<span class="comment">// leave it in _Gscan, but instead gently</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			<span class="comment">// prevented its scheduling until resumption.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			<span class="comment">// Maybe we only use this to bump a suspended</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			<span class="comment">// count and the scheduler skips suspended</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			<span class="comment">// goroutines? That wouldn&#39;t be enough for</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			<span class="comment">// {_Gsyscall,_Gwaiting} -&gt; _Grunning. Maybe</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			<span class="comment">// for all those transitions we need to check</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			<span class="comment">// suspended and deschedule?</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			return suspendGState{g: gp, stopped: stopped}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		case _Grunning:
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			<span class="comment">// Optimization: if there is already a pending preemption request</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			<span class="comment">// (from the previous loop iteration), don&#39;t bother with the atomics.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			if gp.preemptStop &amp;&amp; gp.preempt &amp;&amp; gp.stackguard0 == stackPreempt &amp;&amp; asyncM == gp.m &amp;&amp; asyncM.preemptGen.Load() == asyncGen {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>				break
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			<span class="comment">// Temporarily block state transitions.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			if !castogscanstatus(gp, _Grunning, _Gscanrunning) {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>				break
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			<span class="comment">// Request synchronous preemption.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			gp.preemptStop = true
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			gp.preempt = true
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			gp.stackguard0 = stackPreempt
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			<span class="comment">// Prepare for asynchronous preemption.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			asyncM2 := gp.m
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			asyncGen2 := asyncM2.preemptGen.Load()
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			needAsync := asyncM != asyncM2 || asyncGen != asyncGen2
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			asyncM = asyncM2
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			asyncGen = asyncGen2
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			casfrom_Gscanstatus(gp, _Gscanrunning, _Grunning)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			<span class="comment">// Send asynchronous preemption. We do this</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			<span class="comment">// after CASing the G back to _Grunning</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			<span class="comment">// because preemptM may be synchronous and we</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			<span class="comment">// don&#39;t want to catch the G just spinning on</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			<span class="comment">// its status.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			if preemptMSupported &amp;&amp; debug.asyncpreemptoff == 0 &amp;&amp; needAsync {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>				<span class="comment">// Rate limit preemptM calls. This is</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				<span class="comment">// particularly important on Windows</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>				<span class="comment">// where preemptM is actually</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>				<span class="comment">// synchronous and the spin loop here</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>				<span class="comment">// can lead to live-lock.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>				now := nanotime()
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>				if now &gt;= nextPreemptM {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>					nextPreemptM = now + yieldDelay/2
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>					preemptM(asyncM)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>				}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		<span class="comment">// TODO: Don&#39;t busy wait. This loop should really only</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		<span class="comment">// be a simple read/decide/CAS loop that only fails if</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		<span class="comment">// there&#39;s an active race. Once the CAS succeeds, we</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		<span class="comment">// should queue up the preemption (which will require</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		<span class="comment">// it to be reliable in the _Grunning case, not</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		<span class="comment">// best-effort) and then sleep until we&#39;re notified</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		<span class="comment">// that the goroutine is suspended.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		if i == 0 {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			nextYield = nanotime() + yieldDelay
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		if nanotime() &lt; nextYield {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			procyield(10)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		} else {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			osyield()
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			nextYield = nanotime() + yieldDelay/2
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// resumeG undoes the effects of suspendG, allowing the suspended</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// goroutine to continue from its current safe-point.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>func resumeG(state suspendGState) {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	if state.dead {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		<span class="comment">// We didn&#39;t actually stop anything.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		return
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	gp := state.g
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	switch s := readgstatus(gp); s {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	default:
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		dumpgstatus(gp)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		throw(&#34;unexpected g status&#34;)
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	case _Grunnable | _Gscan,
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		_Gwaiting | _Gscan,
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		_Gsyscall | _Gscan:
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		casfrom_Gscanstatus(gp, s, s&amp;^_Gscan)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	if state.stopped {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		<span class="comment">// We stopped it, so we need to re-schedule it.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		ready(gp, 0, true)
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// canPreemptM reports whether mp is in a state that is safe to preempt.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">// It is nosplit because it has nosplit callers.</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>func canPreemptM(mp *m) bool {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	return mp.locks == 0 &amp;&amp; mp.mallocing == 0 &amp;&amp; mp.preemptoff == &#34;&#34; &amp;&amp; mp.p.ptr().status == _Prunning
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">//go:generate go run mkpreempt.go</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">// asyncPreempt saves all user registers and calls asyncPreempt2.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span><span class="comment">// When stack scanning encounters an asyncPreempt frame, it scans that</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span><span class="comment">// frame and its parent frame conservatively.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span><span class="comment">// asyncPreempt is implemented in assembly.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>func asyncPreempt()
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>func asyncPreempt2() {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	gp := getg()
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	gp.asyncSafePoint = true
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	if gp.preemptStop {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		mcall(preemptPark)
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	} else {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		mcall(gopreempt_m)
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	gp.asyncSafePoint = false
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>}
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">// asyncPreemptStack is the bytes of stack space required to inject an</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span><span class="comment">// asyncPreempt call.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>var asyncPreemptStack = ^uintptr(0)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>func init() {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	f := findfunc(abi.FuncPCABI0(asyncPreempt))
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	total := funcMaxSPDelta(f)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	f = findfunc(abi.FuncPCABIInternal(asyncPreempt2))
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	total += funcMaxSPDelta(f)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	<span class="comment">// Add some overhead for return PCs, etc.</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	asyncPreemptStack = uintptr(total) + 8*goarch.PtrSize
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	if asyncPreemptStack &gt; stackNosplit {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		<span class="comment">// We need more than the nosplit limit. This isn&#39;t</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		<span class="comment">// unsafe, but it may limit asynchronous preemption.</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		<span class="comment">// This may be a problem if we start using more</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		<span class="comment">// registers. In that case, we should store registers</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		<span class="comment">// in a context object. If we pre-allocate one per P,</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		<span class="comment">// asyncPreempt can spill just a few registers to the</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		<span class="comment">// stack, then grab its context object and spill into</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		<span class="comment">// it. When it enters the runtime, it would allocate a</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		<span class="comment">// new context for the P.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		print(&#34;runtime: asyncPreemptStack=&#34;, asyncPreemptStack, &#34;\n&#34;)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		throw(&#34;async stack too large&#34;)
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// wantAsyncPreempt returns whether an asynchronous preemption is</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// queued for gp.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>func wantAsyncPreempt(gp *g) bool {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	<span class="comment">// Check both the G and the P.</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	return (gp.preempt || gp.m.p != 0 &amp;&amp; gp.m.p.ptr().preempt) &amp;&amp; readgstatus(gp)&amp;^_Gscan == _Grunning
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// isAsyncSafePoint reports whether gp at instruction PC is an</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// asynchronous safe point. This indicates that:</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// 1. It&#39;s safe to suspend gp and conservatively scan its stack and</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// registers. There are no potentially hidden pointer values and it&#39;s</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// not in the middle of an atomic sequence like a write barrier.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// 2. gp has enough stack space to inject the asyncPreempt call.</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// 3. It&#39;s generally safe to interact with the runtime, even if we&#39;re</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// in a signal handler stopped here. For example, there are no runtime</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// locks held, so acquiring a runtime lock won&#39;t self-deadlock.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// In some cases the PC is safe for asynchronous preemption but it</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span><span class="comment">// also needs to adjust the resumption PC. The new PC is returned in</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// the second result.</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>func isAsyncSafePoint(gp *g, pc, sp, lr uintptr) (bool, uintptr) {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	mp := gp.m
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	<span class="comment">// Only user Gs can have safe-points. We check this first</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	<span class="comment">// because it&#39;s extremely common that we&#39;ll catch mp in the</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	<span class="comment">// scheduler processing this G preemption.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	if mp.curg != gp {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		return false, 0
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	<span class="comment">// Check M state.</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	if mp.p == 0 || !canPreemptM(mp) {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		return false, 0
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	}
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	<span class="comment">// Check stack space.</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	if sp &lt; gp.stack.lo || sp-gp.stack.lo &lt; asyncPreemptStack {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		return false, 0
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	<span class="comment">// Check if PC is an unsafe-point.</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	f := findfunc(pc)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	if !f.valid() {
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		<span class="comment">// Not Go code.</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		return false, 0
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	if (GOARCH == &#34;mips&#34; || GOARCH == &#34;mipsle&#34; || GOARCH == &#34;mips64&#34; || GOARCH == &#34;mips64le&#34;) &amp;&amp; lr == pc+8 &amp;&amp; funcspdelta(f, pc) == 0 {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		<span class="comment">// We probably stopped at a half-executed CALL instruction,</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		<span class="comment">// where the LR is updated but the PC has not. If we preempt</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		<span class="comment">// here we&#39;ll see a seemingly self-recursive call, which is in</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		<span class="comment">// fact not.</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		<span class="comment">// This is normally ok, as we use the return address saved on</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		<span class="comment">// stack for unwinding, not the LR value. But if this is a</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		<span class="comment">// call to morestack, we haven&#39;t created the frame, and we&#39;ll</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		<span class="comment">// use the LR for unwinding, which will be bad.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		return false, 0
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	up, startpc := pcdatavalue2(f, abi.PCDATA_UnsafePoint, pc)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	if up == abi.UnsafePointUnsafe {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		<span class="comment">// Unsafe-point marked by compiler. This includes</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		<span class="comment">// atomic sequences (e.g., write barrier) and nosplit</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		<span class="comment">// functions (except at calls).</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		return false, 0
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	if fd := funcdata(f, abi.FUNCDATA_LocalsPointerMaps); fd == nil || f.flag&amp;abi.FuncFlagAsm != 0 {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		<span class="comment">// This is assembly code. Don&#39;t assume it&#39;s well-formed.</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		<span class="comment">// TODO: Empirically we still need the fd == nil check. Why?</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		<span class="comment">// TODO: Are there cases that are safe but don&#39;t have a</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		<span class="comment">// locals pointer map, like empty frame functions?</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		<span class="comment">// It might be possible to preempt any assembly functions</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		<span class="comment">// except the ones that have funcFlag_SPWRITE set in f.flag.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		return false, 0
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	<span class="comment">// Check the inner-most name</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	u, uf := newInlineUnwinder(f, pc)
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	name := u.srcFunc(uf).name()
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	if hasPrefix(name, &#34;runtime.&#34;) ||
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		hasPrefix(name, &#34;runtime/internal/&#34;) ||
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		hasPrefix(name, &#34;reflect.&#34;) {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		<span class="comment">// For now we never async preempt the runtime or</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		<span class="comment">// anything closely tied to the runtime. Known issues</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		<span class="comment">// include: various points in the scheduler (&#34;don&#39;t</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		<span class="comment">// preempt between here and here&#34;), much of the defer</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		<span class="comment">// implementation (untyped info on stack), bulk write</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		<span class="comment">// barriers (write barrier check),</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		<span class="comment">// reflect.{makeFuncStub,methodValueCall}.</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		<span class="comment">// TODO(austin): We should improve this, or opt things</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		<span class="comment">// in incrementally.</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		return false, 0
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	switch up {
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	case abi.UnsafePointRestart1, abi.UnsafePointRestart2:
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		<span class="comment">// Restartable instruction sequence. Back off PC to</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		<span class="comment">// the start PC.</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		if startpc == 0 || startpc &gt; pc || pc-startpc &gt; 20 {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>			throw(&#34;bad restart PC&#34;)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		return true, startpc
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	case abi.UnsafePointRestartAtEntry:
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		<span class="comment">// Restart from the function entry at resumption.</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		return true, f.entry()
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	return true, pc
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
</pre><p><a href="preempt.go?m=text">View as plain text</a></p>

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
