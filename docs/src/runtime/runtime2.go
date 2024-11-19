<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/runtime2.go - Go Documentation Server</title>

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
<a href="runtime2.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">runtime2.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/chacha8rand&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// defined constants</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>const (
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// G status</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// Beyond indicating the general state of a G, the G status</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// acts like a lock on the goroutine&#39;s stack (and hence its</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// ability to execute user code).</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// If you add to this list, add to the list</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// of &#34;okay during garbage collection&#34; status</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// in mgcmark.go too.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// TODO(austin): The _Gscan bit could be much lighter-weight.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// For example, we could choose not to run _Gscanrunnable</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// goroutines found in the run queue, rather than CAS-looping</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// until they become _Grunnable. And transitions like</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// _Gscanwaiting -&gt; _Gscanrunnable are actually okay because</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// they don&#39;t affect stack ownership.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// _Gidle means this goroutine was just allocated and has not</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// yet been initialized.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	_Gidle = iota <span class="comment">// 0</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// _Grunnable means this goroutine is on a run queue. It is</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// not currently executing user code. The stack is not owned.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	_Grunnable <span class="comment">// 1</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// _Grunning means this goroutine may execute user code. The</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// stack is owned by this goroutine. It is not on a run queue.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// It is assigned an M and a P (g.m and g.m.p are valid).</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	_Grunning <span class="comment">// 2</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// _Gsyscall means this goroutine is executing a system call.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// It is not executing user code. The stack is owned by this</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// goroutine. It is not on a run queue. It is assigned an M.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	_Gsyscall <span class="comment">// 3</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// _Gwaiting means this goroutine is blocked in the runtime.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// It is not executing user code. It is not on a run queue,</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// but should be recorded somewhere (e.g., a channel wait</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// queue) so it can be ready()d when necessary. The stack is</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// not owned *except* that a channel operation may read or</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// write parts of the stack under the appropriate channel</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// lock. Otherwise, it is not safe to access the stack after a</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// goroutine enters _Gwaiting (e.g., it may get moved).</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	_Gwaiting <span class="comment">// 4</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// _Gmoribund_unused is currently unused, but hardcoded in gdb</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// scripts.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	_Gmoribund_unused <span class="comment">// 5</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// _Gdead means this goroutine is currently unused. It may be</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// just exited, on a free list, or just being initialized. It</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// is not executing user code. It may or may not have a stack</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// allocated. The G and its stack (if any) are owned by the M</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// that is exiting the G or that obtained the G from the free</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// list.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	_Gdead <span class="comment">// 6</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// _Genqueue_unused is currently unused.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	_Genqueue_unused <span class="comment">// 7</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// _Gcopystack means this goroutine&#39;s stack is being moved. It</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// is not executing user code and is not on a run queue. The</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// stack is owned by the goroutine that put it in _Gcopystack.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	_Gcopystack <span class="comment">// 8</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// _Gpreempted means this goroutine stopped itself for a</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// suspendG preemption. It is like _Gwaiting, but nothing is</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// yet responsible for ready()ing it. Some suspendG must CAS</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// the status to _Gwaiting to take responsibility for</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// ready()ing this G.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	_Gpreempted <span class="comment">// 9</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// _Gscan combined with one of the above states other than</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// _Grunning indicates that GC is scanning the stack. The</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// goroutine is not executing user code and the stack is owned</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// by the goroutine that set the _Gscan bit.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// _Gscanrunning is different: it is used to briefly block</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// state transitions while GC signals the G to scan its own</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// stack. This is otherwise like _Grunning.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// atomicstatus&amp;~Gscan gives the state the goroutine will</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// return to when the scan completes.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	_Gscan          = 0x1000
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	_Gscanrunnable  = _Gscan + _Grunnable  <span class="comment">// 0x1001</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	_Gscanrunning   = _Gscan + _Grunning   <span class="comment">// 0x1002</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	_Gscansyscall   = _Gscan + _Gsyscall   <span class="comment">// 0x1003</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	_Gscanwaiting   = _Gscan + _Gwaiting   <span class="comment">// 0x1004</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	_Gscanpreempted = _Gscan + _Gpreempted <span class="comment">// 0x1009</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>const (
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// P status</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// _Pidle means a P is not being used to run user code or the</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// scheduler. Typically, it&#39;s on the idle P list and available</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// to the scheduler, but it may just be transitioning between</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// other states.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// The P is owned by the idle list or by whatever is</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// transitioning its state. Its run queue is empty.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	_Pidle = iota
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// _Prunning means a P is owned by an M and is being used to</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// run user code or the scheduler. Only the M that owns this P</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// is allowed to change the P&#39;s status from _Prunning. The M</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// may transition the P to _Pidle (if it has no more work to</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// do), _Psyscall (when entering a syscall), or _Pgcstop (to</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// halt for the GC). The M may also hand ownership of the P</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// off directly to another M (e.g., to schedule a locked G).</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	_Prunning
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// _Psyscall means a P is not running user code. It has</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// affinity to an M in a syscall but is not owned by it and</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// may be stolen by another M. This is similar to _Pidle but</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// uses lightweight transitions and maintains M affinity.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// Leaving _Psyscall must be done with a CAS, either to steal</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// or retake the P. Note that there&#39;s an ABA hazard: even if</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// an M successfully CASes its original P back to _Prunning</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// after a syscall, it must understand the P may have been</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// used by another M in the interim.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	_Psyscall
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// _Pgcstop means a P is halted for STW and owned by the M</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">// that stopped the world. The M that stopped the world</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// continues to use its P, even in _Pgcstop. Transitioning</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// from _Prunning to _Pgcstop causes an M to release its P and</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// park.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// The P retains its run queue and startTheWorld will restart</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// the scheduler on Ps with non-empty run queues.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	_Pgcstop
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// _Pdead means a P is no longer used (GOMAXPROCS shrank). We</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// reuse Ps if GOMAXPROCS increases. A dead P is mostly</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// stripped of its resources, though a few things remain</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// (e.g., trace buffers).</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	_Pdead
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// Mutual exclusion locks.  In the uncontended case,</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// as fast as spin locks (just a few user-level instructions),</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">// but on the contention path they sleep in the kernel.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// A zeroed Mutex is unlocked (no need to initialize each lock).</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// Initialization is helpful for static lock ranking, but not required.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>type mutex struct {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// Empty struct if lock ranking is disabled, otherwise includes the lock rank</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	lockRankStruct
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// Futex-based impl treats it as uint32 key,</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// while sema-based impl as M* waitm.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// Used to be a union, but unions break precise GC.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	key uintptr
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// sleep and wakeup on one-time events.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">// before any calls to notesleep or notewakeup,</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// must call noteclear to initialize the Note.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// then, exactly one thread can call notesleep</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">// and exactly one thread can call notewakeup (once).</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">// once notewakeup has been called, the notesleep</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">// will return.  future notesleep will return immediately.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span><span class="comment">// subsequent noteclear must be called only after</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">// previous notesleep has returned, e.g. it&#39;s disallowed</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">// to call noteclear straight after notewakeup.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span><span class="comment">// notetsleep is like notesleep but wakes up after</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span><span class="comment">// a given number of nanoseconds even if the event</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// has not yet happened.  if a goroutine uses notetsleep to</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// wake up early, it must wait to call noteclear until it</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// can be sure that no other goroutine is calling</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// notewakeup.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// notesleep/notetsleep are generally called on g0,</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// notetsleepg is similar to notetsleep but is called on user g.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>type note struct {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// Futex-based impl treats it as uint32 key,</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// while sema-based impl as M* waitm.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// Used to be a union, but unions break precise GC.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	key uintptr
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>type funcval struct {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	fn uintptr
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// variable-size, fn-specific data here</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>type iface struct {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	tab  *itab
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	data unsafe.Pointer
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>type eface struct {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	_type *_type
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	data  unsafe.Pointer
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>func efaceOf(ep *any) *eface {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	return (*eface)(unsafe.Pointer(ep))
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">// The guintptr, muintptr, and puintptr are all used to bypass write barriers.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">// It is particularly important to avoid write barriers when the current P has</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// been released, because the GC thinks the world is stopped, and an</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// unexpected write barrier would not be synchronized with the GC,</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// which can lead to a half-executed write barrier that has marked the object</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// but not queued it. If the GC skips the object and completes before the</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// queuing can occur, it will incorrectly free the object.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// We tried using special assignment functions invoked only when not</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// holding a running P, but then some updates to a particular memory</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// word went through write barriers and some did not. This breaks the</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// write barrier shadow checking mode, and it is also scary: better to have</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// a word that is completely ignored by the GC than to have one for which</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// only a few updates are ignored.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">// Gs and Ps are always reachable via true pointers in the</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">// allgs and allp lists or (during allocation before they reach those lists)</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">// from stack variables.</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">// Ms are always reachable via true pointers either from allm or</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">// freem. Unlike Gs and Ps we do free Ms, so it&#39;s important that</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span><span class="comment">// nothing ever hold an muintptr across a safe point.</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span><span class="comment">// A guintptr holds a goroutine pointer, but typed as a uintptr</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span><span class="comment">// to bypass write barriers. It is used in the Gobuf goroutine state</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">// and in scheduling lists that are manipulated without a P.</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">// The Gobuf.g goroutine pointer is almost always updated by assembly code.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// In one of the few places it is updated by Go code - func save - it must be</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// treated as a uintptr to avoid a write barrier being emitted at a bad time.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">// Instead of figuring out how to emit the write barriers missing in the</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">// assembly manipulation, we change the type of the field to uintptr,</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span><span class="comment">// so that it does not require write barriers at all.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span><span class="comment">// Goroutine structs are published in the allg list and never freed.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// That will keep the goroutine structs from being collected.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// There is never a time that Gobuf.g&#39;s contain the only references</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// to a goroutine: the publishing of the goroutine in allg comes first.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span><span class="comment">// Goroutine pointers are also kept in non-GC-visible places like TLS,</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span><span class="comment">// so I can&#39;t see them ever moving. If we did want to start moving data</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">// in the GC, we&#39;d need to allocate the goroutine structs from an</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// alternate arena. Using guintptr doesn&#39;t make that problem any worse.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span><span class="comment">// Note that pollDesc.rg, pollDesc.wg also store g in uintptr form,</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span><span class="comment">// so they would need to be updated too if g&#39;s start moving.</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>type guintptr uintptr
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>func (gp guintptr) ptr() *g { return (*g)(unsafe.Pointer(gp)) }
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>func (gp *guintptr) set(g *g) { *gp = guintptr(unsafe.Pointer(g)) }
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>func (gp *guintptr) cas(old, new guintptr) bool {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	return atomic.Casuintptr((*uintptr)(unsafe.Pointer(gp)), uintptr(old), uintptr(new))
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>func (gp *g) guintptr() guintptr {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	return guintptr(unsafe.Pointer(gp))
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// setGNoWB performs *gp = new without a write barrier.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">// For times when it&#39;s impractical to use a guintptr.</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>func setGNoWB(gp **g, new *g) {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	(*guintptr)(unsafe.Pointer(gp)).set(new)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>type puintptr uintptr
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>func (pp puintptr) ptr() *p { return (*p)(unsafe.Pointer(pp)) }
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>func (pp *puintptr) set(p *p) { *pp = puintptr(unsafe.Pointer(p)) }
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span><span class="comment">// muintptr is a *m that is not tracked by the garbage collector.</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">// Because we do free Ms, there are some additional constrains on</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">// muintptrs:</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span><span class="comment">//  1. Never hold an muintptr locally across a safe point.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">//  2. Any muintptr in the heap must be owned by the M itself so it can</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">//     ensure it is not in use when the last true *m is released.</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>type muintptr uintptr
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>func (mp muintptr) ptr() *m { return (*m)(unsafe.Pointer(mp)) }
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>func (mp *muintptr) set(m *m) { *mp = muintptr(unsafe.Pointer(m)) }
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// setMNoWB performs *mp = new without a write barrier.</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">// For times when it&#39;s impractical to use an muintptr.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>func setMNoWB(mp **m, new *m) {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	(*muintptr)(unsafe.Pointer(mp)).set(new)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>type gobuf struct {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// The offsets of sp, pc, and g are known to (hard-coded in) libmach.</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// ctxt is unusual with respect to GC: it may be a</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	<span class="comment">// heap-allocated funcval, so GC needs to track it, but it</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	<span class="comment">// needs to be set and cleared from assembly, where it&#39;s</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">// difficult to have write barriers. However, ctxt is really a</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	<span class="comment">// saved, live register, and we only ever exchange it between</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">// the real register and the gobuf. Hence, we treat it as a</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// root during stack scanning, which means assembly that saves</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// and restores it doesn&#39;t need write barriers. It&#39;s still</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// typed as a pointer so that any other writes from Go get</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// write barriers.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	sp   uintptr
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	pc   uintptr
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	g    guintptr
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	ctxt unsafe.Pointer
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	ret  uintptr
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	lr   uintptr
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	bp   uintptr <span class="comment">// for framepointer-enabled architectures</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// sudog (pseudo-g) represents a g in a wait list, such as for sending/receiving</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// on a channel.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// sudog is necessary because the g ↔ synchronization object relation</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// is many-to-many. A g can be on many wait lists, so there may be</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// many sudogs for one g; and many gs may be waiting on the same</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// synchronization object, so there may be many sudogs for one object.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// sudogs are allocated from a special pool. Use acquireSudog and</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// releaseSudog to allocate and free them.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>type sudog struct {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	<span class="comment">// The following fields are protected by the hchan.lock of the</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	<span class="comment">// channel this sudog is blocking on. shrinkstack depends on</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	<span class="comment">// this for sudogs involved in channel ops.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	g *g
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	next *sudog
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	prev *sudog
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	elem unsafe.Pointer <span class="comment">// data element (may point to stack)</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	<span class="comment">// The following fields are never accessed concurrently.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	<span class="comment">// For channels, waitlink is only accessed by g.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// For semaphores, all fields (including the ones above)</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	<span class="comment">// are only accessed when holding a semaRoot lock.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	acquiretime int64
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	releasetime int64
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	ticket      uint32
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	<span class="comment">// isSelect indicates g is participating in a select, so</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	<span class="comment">// g.selectDone must be CAS&#39;d to win the wake-up race.</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	isSelect bool
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	<span class="comment">// success indicates whether communication over channel c</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	<span class="comment">// succeeded. It is true if the goroutine was awoken because a</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	<span class="comment">// value was delivered over channel c, and false if awoken</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	<span class="comment">// because c was closed.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	success bool
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	<span class="comment">// waiters is a count of semaRoot waiting list other than head of list,</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	<span class="comment">// clamped to a uint16 to fit in unused space.</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	<span class="comment">// Only meaningful at the head of the list.</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	<span class="comment">// (If we wanted to be overly clever, we could store a high 16 bits</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">// in the second entry in the list.)</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	waiters uint16
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	parent   *sudog <span class="comment">// semaRoot binary tree</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	waitlink *sudog <span class="comment">// g.waiting list or semaRoot</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	waittail *sudog <span class="comment">// semaRoot</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	c        *hchan <span class="comment">// channel</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>}
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>type libcall struct {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	fn   uintptr
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	n    uintptr <span class="comment">// number of parameters</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	args uintptr <span class="comment">// parameters</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	r1   uintptr <span class="comment">// return values</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	r2   uintptr
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	err  uintptr <span class="comment">// error number</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span><span class="comment">// Stack describes a Go execution stack.</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span><span class="comment">// The bounds of the stack are exactly [lo, hi),</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span><span class="comment">// with no implicit data structures on either side.</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>type stack struct {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	lo uintptr
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	hi uintptr
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span><span class="comment">// heldLockInfo gives info on a held lock and the rank of that lock</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>type heldLockInfo struct {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	lockAddr uintptr
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	rank     lockRank
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>type g struct {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	<span class="comment">// Stack parameters.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	<span class="comment">// stack describes the actual stack memory: [stack.lo, stack.hi).</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	<span class="comment">// stackguard0 is the stack pointer compared in the Go stack growth prologue.</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	<span class="comment">// It is stack.lo+StackGuard normally, but can be StackPreempt to trigger a preemption.</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	<span class="comment">// stackguard1 is the stack pointer compared in the //go:systemstack stack growth prologue.</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	<span class="comment">// It is stack.lo+StackGuard on g0 and gsignal stacks.</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	<span class="comment">// It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	stack       stack   <span class="comment">// offset known to runtime/cgo</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	stackguard0 uintptr <span class="comment">// offset known to liblink</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	stackguard1 uintptr <span class="comment">// offset known to liblink</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	_panic    *_panic <span class="comment">// innermost panic - offset known to liblink</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	_defer    *_defer <span class="comment">// innermost defer</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	m         *m      <span class="comment">// current m; offset known to arm liblink</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	sched     gobuf
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	syscallsp uintptr <span class="comment">// if status==Gsyscall, syscallsp = sched.sp to use during gc</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	syscallpc uintptr <span class="comment">// if status==Gsyscall, syscallpc = sched.pc to use during gc</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	stktopsp  uintptr <span class="comment">// expected sp at top of stack, to check in traceback</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// param is a generic pointer parameter field used to pass</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	<span class="comment">// values in particular contexts where other storage for the</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	<span class="comment">// parameter would be difficult to find. It is currently used</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	<span class="comment">// in four ways:</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	<span class="comment">// 1. When a channel operation wakes up a blocked goroutine, it sets param to</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	<span class="comment">//    point to the sudog of the completed blocking operation.</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	<span class="comment">// 2. By gcAssistAlloc1 to signal back to its caller that the goroutine completed</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	<span class="comment">//    the GC cycle. It is unsafe to do so in any other way, because the goroutine&#39;s</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	<span class="comment">//    stack may have moved in the meantime.</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	<span class="comment">// 3. By debugCallWrap to pass parameters to a new goroutine because allocating a</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	<span class="comment">//    closure in the runtime is forbidden.</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	<span class="comment">// 4. When a panic is recovered and control returns to the respective frame,</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	<span class="comment">//    param may point to a savedOpenDeferState.</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	param        unsafe.Pointer
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	atomicstatus atomic.Uint32
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	stackLock    uint32 <span class="comment">// sigprof/scang lock; TODO: fold in to atomicstatus</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	goid         uint64
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	schedlink    guintptr
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	waitsince    int64      <span class="comment">// approx time when the g become blocked</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	waitreason   waitReason <span class="comment">// if status==Gwaiting</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	preempt       bool <span class="comment">// preemption signal, duplicates stackguard0 = stackpreempt</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	preemptStop   bool <span class="comment">// transition to _Gpreempted on preemption; otherwise, just deschedule</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	preemptShrink bool <span class="comment">// shrink stack at synchronous safe point</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	<span class="comment">// asyncSafePoint is set if g is stopped at an asynchronous</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	<span class="comment">// safe point. This means there are frames on the stack</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	<span class="comment">// without precise pointer information.</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	asyncSafePoint bool
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	paniconfault bool <span class="comment">// panic (instead of crash) on unexpected fault address</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	gcscandone   bool <span class="comment">// g has scanned stack; protected by _Gscan bit in status</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	throwsplit   bool <span class="comment">// must not split stack</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	<span class="comment">// activeStackChans indicates that there are unlocked channels</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	<span class="comment">// pointing into this goroutine&#39;s stack. If true, stack</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	<span class="comment">// copying needs to acquire channel locks to protect these</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	<span class="comment">// areas of the stack.</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	activeStackChans bool
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	<span class="comment">// parkingOnChan indicates that the goroutine is about to</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	<span class="comment">// park on a chansend or chanrecv. Used to signal an unsafe point</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	<span class="comment">// for stack shrinking.</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	parkingOnChan atomic.Bool
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	<span class="comment">// inMarkAssist indicates whether the goroutine is in mark assist.</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	<span class="comment">// Used by the execution tracer.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	inMarkAssist bool
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	coroexit     bool <span class="comment">// argument to coroswitch_m</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	raceignore    int8  <span class="comment">// ignore race detection events</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	nocgocallback bool  <span class="comment">// whether disable callback from C</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	tracking      bool  <span class="comment">// whether we&#39;re tracking this G for sched latency statistics</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	trackingSeq   uint8 <span class="comment">// used to decide whether to track this G</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	trackingStamp int64 <span class="comment">// timestamp of when the G last started being tracked</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	runnableTime  int64 <span class="comment">// the amount of time spent runnable, cleared when running, only used when tracking</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	lockedm       muintptr
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	sig           uint32
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	writebuf      []byte
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	sigcode0      uintptr
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	sigcode1      uintptr
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	sigpc         uintptr
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	parentGoid    uint64          <span class="comment">// goid of goroutine that created this goroutine</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	gopc          uintptr         <span class="comment">// pc of go statement that created this goroutine</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	ancestors     *[]ancestorInfo <span class="comment">// ancestor information goroutine(s) that created this goroutine (only used if debug.tracebackancestors)</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	startpc       uintptr         <span class="comment">// pc of goroutine function</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	racectx       uintptr
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	waiting       *sudog         <span class="comment">// sudog structures this g is waiting on (that have a valid elem ptr); in lock order</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	cgoCtxt       []uintptr      <span class="comment">// cgo traceback context</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	labels        unsafe.Pointer <span class="comment">// profiler labels</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	timer         *timer         <span class="comment">// cached timer for time.Sleep</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	selectDone    atomic.Uint32  <span class="comment">// are we participating in a select and did someone win the race?</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	coroarg *coro <span class="comment">// argument during coroutine transfers</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	<span class="comment">// goroutineProfiled indicates the status of this goroutine&#39;s stack for the</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	<span class="comment">// current in-progress goroutine profile</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	goroutineProfiled goroutineProfileStateHolder
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	<span class="comment">// Per-G tracer state.</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	trace gTraceState
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	<span class="comment">// Per-G GC state</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	<span class="comment">// gcAssistBytes is this G&#39;s GC assist credit in terms of</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	<span class="comment">// bytes allocated. If this is positive, then the G has credit</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	<span class="comment">// to allocate gcAssistBytes bytes without assisting. If this</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	<span class="comment">// is negative, then the G must correct this by performing</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	<span class="comment">// scan work. We track this in bytes to make it fast to update</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	<span class="comment">// and check for debt in the malloc hot path. The assist ratio</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	<span class="comment">// determines how this corresponds to scan work debt.</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	gcAssistBytes int64
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>}
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span><span class="comment">// gTrackingPeriod is the number of transitions out of _Grunning between</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span><span class="comment">// latency tracking runs.</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>const gTrackingPeriod = 8
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>const (
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	<span class="comment">// tlsSlots is the number of pointer-sized slots reserved for TLS on some platforms,</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	<span class="comment">// like Windows.</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	tlsSlots = 6
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	tlsSize  = tlsSlots * goarch.PtrSize
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>)
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">// Values for m.freeWait.</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>const (
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	freeMStack = 0 <span class="comment">// M done, free stack and reference.</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	freeMRef   = 1 <span class="comment">// M done, free reference.</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	freeMWait  = 2 <span class="comment">// M still in use.</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>)
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>type m struct {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	g0      *g     <span class="comment">// goroutine with scheduling stack</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	morebuf gobuf  <span class="comment">// gobuf arg to morestack</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	divmod  uint32 <span class="comment">// div/mod denominator for arm - known to liblink</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	_       uint32 <span class="comment">// align next field to 8 bytes</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	<span class="comment">// Fields not known to debuggers.</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	procid        uint64            <span class="comment">// for debuggers, but offset not hard-coded</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	gsignal       *g                <span class="comment">// signal-handling g</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	goSigStack    gsignalStack      <span class="comment">// Go-allocated signal handling stack</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	sigmask       sigset            <span class="comment">// storage for saved signal mask</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	tls           [tlsSlots]uintptr <span class="comment">// thread-local storage (for x86 extern register)</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	mstartfn      func()
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	curg          *g       <span class="comment">// current running goroutine</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	caughtsig     guintptr <span class="comment">// goroutine running during fatal signal</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	p             puintptr <span class="comment">// attached p for executing go code (nil if not executing go code)</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	nextp         puintptr
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	oldp          puintptr <span class="comment">// the p that was attached before executing a syscall</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	id            int64
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	mallocing     int32
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	throwing      throwType
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	preemptoff    string <span class="comment">// if != &#34;&#34;, keep curg running on this m</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	locks         int32
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	dying         int32
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	profilehz     int32
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	spinning      bool <span class="comment">// m is out of work and is actively looking for work</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	blocked       bool <span class="comment">// m is blocked on a note</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	newSigstack   bool <span class="comment">// minit on C thread called sigaltstack</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	printlock     int8
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	incgo         bool          <span class="comment">// m is executing a cgo call</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	isextra       bool          <span class="comment">// m is an extra m</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	isExtraInC    bool          <span class="comment">// m is an extra m that is not executing Go code</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	isExtraInSig  bool          <span class="comment">// m is an extra m in a signal handler</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	freeWait      atomic.Uint32 <span class="comment">// Whether it is safe to free g0 and delete m (one of freeMRef, freeMStack, freeMWait)</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	needextram    bool
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	traceback     uint8
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	ncgocall      uint64        <span class="comment">// number of cgo calls in total</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	ncgo          int32         <span class="comment">// number of cgo calls currently in progress</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	cgoCallersUse atomic.Uint32 <span class="comment">// if non-zero, cgoCallers in use temporarily</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	cgoCallers    *cgoCallers   <span class="comment">// cgo traceback if crashing in cgo call</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	park          note
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	alllink       *m <span class="comment">// on allm</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	schedlink     muintptr
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	lockedg       guintptr
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	createstack   [32]uintptr <span class="comment">// stack that created this thread, it&#39;s used for StackRecord.Stack0, so it must align with it.</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	lockedExt     uint32      <span class="comment">// tracking for external LockOSThread</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	lockedInt     uint32      <span class="comment">// tracking for internal lockOSThread</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	nextwaitm     muintptr    <span class="comment">// next m waiting for lock</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	mLockProfile mLockProfile <span class="comment">// fields relating to runtime.lock contention</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	<span class="comment">// wait* are used to carry arguments from gopark into park_m, because</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	<span class="comment">// there&#39;s no stack to put them on. That is their sole purpose.</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	waitunlockf          func(*g, unsafe.Pointer) bool
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	waitlock             unsafe.Pointer
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	waitTraceBlockReason traceBlockReason
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	waitTraceSkip        int
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	syscalltick uint32
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	freelink    *m <span class="comment">// on sched.freem</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	trace       mTraceState
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	<span class="comment">// these are here because they are too large to be on the stack</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	<span class="comment">// of low-level NOSPLIT functions.</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	libcall   libcall
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	libcallpc uintptr <span class="comment">// for cpu profiler</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	libcallsp uintptr
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	libcallg  guintptr
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	syscall   libcall <span class="comment">// stores syscall parameters on windows</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	vdsoSP uintptr <span class="comment">// SP for traceback while in VDSO call (0 if not in call)</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	vdsoPC uintptr <span class="comment">// PC for traceback while in VDSO call</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	<span class="comment">// preemptGen counts the number of completed preemption</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	<span class="comment">// signals. This is used to detect when a preemption is</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	<span class="comment">// requested, but fails.</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	preemptGen atomic.Uint32
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	<span class="comment">// Whether this is a pending preemption signal on this M.</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	signalPending atomic.Uint32
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	<span class="comment">// pcvalue lookup cache</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	pcvalueCache pcvalueCache
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	dlogPerM
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	mOS
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	chacha8   chacha8rand.State
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	cheaprand uint64
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	<span class="comment">// Up to 10 locks held by this m, maintained by the lock ranking code.</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	locksHeldLen int
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	locksHeld    [10]heldLockInfo
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>}
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>type p struct {
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	id          int32
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	status      uint32 <span class="comment">// one of pidle/prunning/...</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	link        puintptr
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	schedtick   uint32     <span class="comment">// incremented on every scheduler call</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	syscalltick uint32     <span class="comment">// incremented on every system call</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	sysmontick  sysmontick <span class="comment">// last tick observed by sysmon</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	m           muintptr   <span class="comment">// back-link to associated m (nil if idle)</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	mcache      *mcache
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	pcache      pageCache
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	raceprocctx uintptr
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	deferpool    []*_defer <span class="comment">// pool of available defer structs (see panic.go)</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	deferpoolbuf [32]*_defer
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	<span class="comment">// Cache of goroutine ids, amortizes accesses to runtime·sched.goidgen.</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	goidcache    uint64
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	goidcacheend uint64
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	<span class="comment">// Queue of runnable goroutines. Accessed without lock.</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	runqhead uint32
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	runqtail uint32
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	runq     [256]guintptr
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	<span class="comment">// runnext, if non-nil, is a runnable G that was ready&#39;d by</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	<span class="comment">// the current G and should be run next instead of what&#39;s in</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	<span class="comment">// runq if there&#39;s time remaining in the running G&#39;s time</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	<span class="comment">// slice. It will inherit the time left in the current time</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	<span class="comment">// slice. If a set of goroutines is locked in a</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	<span class="comment">// communicate-and-wait pattern, this schedules that set as a</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	<span class="comment">// unit and eliminates the (potentially large) scheduling</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	<span class="comment">// latency that otherwise arises from adding the ready&#39;d</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	<span class="comment">// goroutines to the end of the run queue.</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	<span class="comment">// Note that while other P&#39;s may atomically CAS this to zero,</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	<span class="comment">// only the owner P can CAS it to a valid G.</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	runnext guintptr
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	<span class="comment">// Available G&#39;s (status == Gdead)</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	gFree struct {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		gList
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		n int32
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	}
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	sudogcache []*sudog
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	sudogbuf   [128]*sudog
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	<span class="comment">// Cache of mspan objects from the heap.</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	mspancache struct {
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		<span class="comment">// We need an explicit length here because this field is used</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>		<span class="comment">// in allocation codepaths where write barriers are not allowed,</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		<span class="comment">// and eliminating the write barrier/keeping it eliminated from</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		<span class="comment">// slice updates is tricky, more so than just managing the length</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		<span class="comment">// ourselves.</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>		len int
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>		buf [128]*mspan
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	}
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	<span class="comment">// Cache of a single pinner object to reduce allocations from repeated</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	<span class="comment">// pinner creation.</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	pinnerCache *pinner
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	trace pTraceState
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	palloc persistentAlloc <span class="comment">// per-P to avoid mutex</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	<span class="comment">// The when field of the first entry on the timer heap.</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	<span class="comment">// This is 0 if the timer heap is empty.</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	timer0When atomic.Int64
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	<span class="comment">// The earliest known nextwhen field of a timer with</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	<span class="comment">// timerModifiedEarlier status. Because the timer may have been</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	<span class="comment">// modified again, there need not be any timer with this value.</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	<span class="comment">// This is 0 if there are no timerModifiedEarlier timers.</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	timerModifiedEarliest atomic.Int64
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	<span class="comment">// Per-P GC state</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	gcAssistTime         int64 <span class="comment">// Nanoseconds in assistAlloc</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	gcFractionalMarkTime int64 <span class="comment">// Nanoseconds in fractional mark worker (atomic)</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	<span class="comment">// limiterEvent tracks events for the GC CPU limiter.</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	limiterEvent limiterEvent
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	<span class="comment">// gcMarkWorkerMode is the mode for the next mark worker to run in.</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	<span class="comment">// That is, this is used to communicate with the worker goroutine</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	<span class="comment">// selected for immediate execution by</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	<span class="comment">// gcController.findRunnableGCWorker. When scheduling other goroutines,</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	<span class="comment">// this field must be set to gcMarkWorkerNotWorker.</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	gcMarkWorkerMode gcMarkWorkerMode
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	<span class="comment">// gcMarkWorkerStartTime is the nanotime() at which the most recent</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	<span class="comment">// mark worker started.</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	gcMarkWorkerStartTime int64
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	<span class="comment">// gcw is this P&#39;s GC work buffer cache. The work buffer is</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	<span class="comment">// filled by write barriers, drained by mutator assists, and</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	<span class="comment">// disposed on certain GC state transitions.</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	gcw gcWork
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	<span class="comment">// wbBuf is this P&#39;s GC write barrier buffer.</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	<span class="comment">// TODO: Consider caching this in the running G.</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	wbBuf wbBuf
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	runSafePointFn uint32 <span class="comment">// if 1, run sched.safePointFn at next safe point</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	<span class="comment">// statsSeq is a counter indicating whether this P is currently</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	<span class="comment">// writing any stats. Its value is even when not, odd when it is.</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	statsSeq atomic.Uint32
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	<span class="comment">// Lock for timers. We normally access the timers while running</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	<span class="comment">// on this P, but the scheduler can also do it from a different P.</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	timersLock mutex
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	<span class="comment">// Actions to take at some time. This is used to implement the</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	<span class="comment">// standard library&#39;s time package.</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	<span class="comment">// Must hold timersLock to access.</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	timers []*timer
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	<span class="comment">// Number of timers in P&#39;s heap.</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	numTimers atomic.Uint32
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	<span class="comment">// Number of timerDeleted timers in P&#39;s heap.</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	deletedTimers atomic.Uint32
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	<span class="comment">// Race context used while executing timer functions.</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	timerRaceCtx uintptr
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	<span class="comment">// maxStackScanDelta accumulates the amount of stack space held by</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	<span class="comment">// live goroutines (i.e. those eligible for stack scanning).</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	<span class="comment">// Flushed to gcController.maxStackScan once maxStackScanSlack</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>	<span class="comment">// or -maxStackScanSlack is reached.</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	maxStackScanDelta int64
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>	<span class="comment">// gc-time statistics about current goroutines</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>	<span class="comment">// Note that this differs from maxStackScan in that this</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	<span class="comment">// accumulates the actual stack observed to be used at GC time (hi - sp),</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	<span class="comment">// not an instantaneous measure of the total stack size that might need</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	<span class="comment">// to be scanned (hi - lo).</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	scannedStackSize uint64 <span class="comment">// stack size of goroutines scanned by this P</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	scannedStacks    uint64 <span class="comment">// number of goroutines scanned by this P</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	<span class="comment">// preempt is set to indicate that this P should be enter the</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	<span class="comment">// scheduler ASAP (regardless of what G is running on it).</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>	preempt bool
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	<span class="comment">// pageTraceBuf is a buffer for writing out page allocation/free/scavenge traces.</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	<span class="comment">// Used only if GOEXPERIMENT=pagetrace.</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	pageTraceBuf pageTraceBuf
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	<span class="comment">// Padding is no longer needed. False sharing is now not a worry because p is large enough</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	<span class="comment">// that its size class is an integer multiple of the cache line size (for any of our architectures).</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>}
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>type schedt struct {
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	goidgen   atomic.Uint64
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	lastpoll  atomic.Int64 <span class="comment">// time of last network poll, 0 if currently polling</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	pollUntil atomic.Int64 <span class="comment">// time to which current poll is sleeping</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	lock mutex
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	<span class="comment">// When increasing nmidle, nmidlelocked, nmsys, or nmfreed, be</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	<span class="comment">// sure to call checkdead().</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	midle        muintptr <span class="comment">// idle m&#39;s waiting for work</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	nmidle       int32    <span class="comment">// number of idle m&#39;s waiting for work</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	nmidlelocked int32    <span class="comment">// number of locked m&#39;s waiting for work</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	mnext        int64    <span class="comment">// number of m&#39;s that have been created and next M ID</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	maxmcount    int32    <span class="comment">// maximum number of m&#39;s allowed (or die)</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	nmsys        int32    <span class="comment">// number of system m&#39;s not counted for deadlock</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	nmfreed      int64    <span class="comment">// cumulative number of freed m&#39;s</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	ngsys atomic.Int32 <span class="comment">// number of system goroutines</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	pidle        puintptr <span class="comment">// idle p&#39;s</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	npidle       atomic.Int32
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	nmspinning   atomic.Int32  <span class="comment">// See &#34;Worker thread parking/unparking&#34; comment in proc.go.</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	needspinning atomic.Uint32 <span class="comment">// See &#34;Delicate dance&#34; comment in proc.go. Boolean. Must hold sched.lock to set to 1.</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	<span class="comment">// Global runnable queue.</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	runq     gQueue
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	runqsize int32
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>	<span class="comment">// disable controls selective disabling of the scheduler.</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	<span class="comment">// Use schedEnableUser to control this.</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	<span class="comment">// disable is protected by sched.lock.</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	disable struct {
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		<span class="comment">// user disables scheduling of user goroutines.</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		user     bool
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		runnable gQueue <span class="comment">// pending runnable Gs</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		n        int32  <span class="comment">// length of runnable</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	}
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	<span class="comment">// Global cache of dead G&#39;s.</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	gFree struct {
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		lock    mutex
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>		stack   gList <span class="comment">// Gs with stacks</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>		noStack gList <span class="comment">// Gs without stacks</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>		n       int32
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>	}
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	<span class="comment">// Central cache of sudog structs.</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	sudoglock  mutex
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	sudogcache *sudog
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	<span class="comment">// Central pool of available defer structs.</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	deferlock mutex
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	deferpool *_defer
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	<span class="comment">// freem is the list of m&#39;s waiting to be freed when their</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	<span class="comment">// m.exited is set. Linked through m.freelink.</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	freem *m
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	gcwaiting  atomic.Bool <span class="comment">// gc is waiting to run</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	stopwait   int32
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	stopnote   note
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	sysmonwait atomic.Bool
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	sysmonnote note
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	<span class="comment">// safePointFn should be called on each P at the next GC</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	<span class="comment">// safepoint if p.runSafePointFn is set.</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	safePointFn   func(*p)
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	safePointWait int32
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	safePointNote note
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	profilehz int32 <span class="comment">// cpu profiling rate</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	procresizetime int64 <span class="comment">// nanotime() of last change to gomaxprocs</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	totaltime      int64 <span class="comment">// ∫gomaxprocs dt up to procresizetime</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	<span class="comment">// sysmonlock protects sysmon&#39;s actions on the runtime.</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	<span class="comment">// Acquire and hold this mutex to block sysmon from interacting</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	<span class="comment">// with the rest of the runtime.</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	sysmonlock mutex
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	<span class="comment">// timeToRun is a distribution of scheduling latencies, defined</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	<span class="comment">// as the sum of time a G spends in the _Grunnable state before</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	<span class="comment">// it transitions to _Grunning.</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	timeToRun timeHistogram
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>	<span class="comment">// idleTime is the total CPU time Ps have &#34;spent&#34; idle.</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>	<span class="comment">// Reset on each GC cycle.</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>	idleTime atomic.Int64
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	<span class="comment">// totalMutexWaitTime is the sum of time goroutines have spent in _Gwaiting</span>
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	<span class="comment">// with a waitreason of the form waitReasonSync{RW,}Mutex{R,}Lock.</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	totalMutexWaitTime atomic.Int64
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	<span class="comment">// stwStoppingTimeGC/Other are distributions of stop-the-world stopping</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	<span class="comment">// latencies, defined as the time taken by stopTheWorldWithSema to get</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	<span class="comment">// all Ps to stop. stwStoppingTimeGC covers all GC-related STWs,</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>	<span class="comment">// stwStoppingTimeOther covers the others.</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	stwStoppingTimeGC    timeHistogram
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	stwStoppingTimeOther timeHistogram
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	<span class="comment">// stwTotalTimeGC/Other are distributions of stop-the-world total</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	<span class="comment">// latencies, defined as the total time from stopTheWorldWithSema to</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	<span class="comment">// startTheWorldWithSema. This is a superset of</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	<span class="comment">// stwStoppingTimeGC/Other. stwTotalTimeGC covers all GC-related STWs,</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	<span class="comment">// stwTotalTimeOther covers the others.</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	stwTotalTimeGC    timeHistogram
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	stwTotalTimeOther timeHistogram
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	<span class="comment">// totalRuntimeLockWaitTime (plus the value of lockWaitTime on each M in</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	<span class="comment">// allm) is the sum of time goroutines have spent in _Grunnable and with an</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	<span class="comment">// M, but waiting for locks within the runtime. This field stores the value</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	<span class="comment">// for Ms that have exited.</span>
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>	totalRuntimeLockWaitTime atomic.Int64
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>}
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span><span class="comment">// Values for the flags field of a sigTabT.</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>const (
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>	_SigNotify   = 1 &lt;&lt; iota <span class="comment">// let signal.Notify have signal, even if from kernel</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	_SigKill                 <span class="comment">// if signal.Notify doesn&#39;t take it, exit quietly</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	_SigThrow                <span class="comment">// if signal.Notify doesn&#39;t take it, exit loudly</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	_SigPanic                <span class="comment">// if the signal is from the kernel, panic</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	_SigDefault              <span class="comment">// if the signal isn&#39;t explicitly requested, don&#39;t monitor it</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	_SigGoExit               <span class="comment">// cause all runtime procs to exit (only used on Plan 9).</span>
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	_SigSetStack             <span class="comment">// Don&#39;t explicitly install handler, but add SA_ONSTACK to existing libc handler</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	_SigUnblock              <span class="comment">// always unblock; see blockableSig</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>	_SigIgn                  <span class="comment">// _SIG_DFL action is to ignore the signal</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>)
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span><span class="comment">// Layout of in-memory per-function information prepared by linker</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span><span class="comment">// See https://golang.org/s/go12symtab.</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span><span class="comment">// Keep in sync with linker (../cmd/link/internal/ld/pcln.go:/pclntab)</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span><span class="comment">// and with package debug/gosym and with symtab.go in package runtime.</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>type _func struct {
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>	sys.NotInHeap <span class="comment">// Only in static data</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	entryOff uint32 <span class="comment">// start pc, as offset from moduledata.text/pcHeader.textStart</span>
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	nameOff  int32  <span class="comment">// function name, as index into moduledata.funcnametab.</span>
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>	args        int32  <span class="comment">// in/out args size</span>
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>	deferreturn uint32 <span class="comment">// offset of start of a deferreturn call instruction from entry, if any.</span>
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>	pcsp      uint32
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>	pcfile    uint32
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>	pcln      uint32
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>	npcdata   uint32
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>	cuOffset  uint32     <span class="comment">// runtime.cutab offset of this function&#39;s CU</span>
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>	startLine int32      <span class="comment">// line number of start of function (func keyword/TEXT directive)</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>	funcID    abi.FuncID <span class="comment">// set for certain special runtime functions</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>	flag      abi.FuncFlag
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>	_         [1]byte <span class="comment">// pad</span>
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>	nfuncdata uint8   <span class="comment">// must be last, must end on a uint32-aligned boundary</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>	<span class="comment">// The end of the struct is followed immediately by two variable-length</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	<span class="comment">// arrays that reference the pcdata and funcdata locations for this</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>	<span class="comment">// function.</span>
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>	<span class="comment">// pcdata contains the offset into moduledata.pctab for the start of</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	<span class="comment">// that index&#39;s table. e.g.,</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>	<span class="comment">// &amp;moduledata.pctab[_func.pcdata[_PCDATA_UnsafePoint]] is the start of</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>	<span class="comment">// the unsafe point table.</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	<span class="comment">// An offset of 0 indicates that there is no table.</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>	<span class="comment">// pcdata [npcdata]uint32</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>	<span class="comment">// funcdata contains the offset past moduledata.gofunc which contains a</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	<span class="comment">// pointer to that index&#39;s funcdata. e.g.,</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>	<span class="comment">// *(moduledata.gofunc +  _func.funcdata[_FUNCDATA_ArgsPointerMaps]) is</span>
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>	<span class="comment">// the argument pointer map.</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>	<span class="comment">// An offset of ^uint32(0) indicates that there is no entry.</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	<span class="comment">// funcdata [nfuncdata]uint32</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>}
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span><span class="comment">// Pseudo-Func that is returned for PCs that occur in inlined code.</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span><span class="comment">// A *Func can be either a *_func or a *funcinl, and they are distinguished</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span><span class="comment">// by the first uintptr.</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span><span class="comment">// TODO(austin): Can we merge this with inlinedCall?</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>type funcinl struct {
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	ones      uint32  <span class="comment">// set to ^0 to distinguish from _func</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>	entry     uintptr <span class="comment">// entry of the real (the &#34;outermost&#34;) frame</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	name      string
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>	file      string
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	line      int32
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	startLine int32
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>}
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span><span class="comment">// layout of Itab known to compilers</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span><span class="comment">// allocated in non-garbage-collected memory</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span><span class="comment">// Needs to be in sync with</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span><span class="comment">// ../cmd/compile/internal/reflectdata/reflect.go:/^func.WritePluginTable.</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>type itab struct {
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	inter *interfacetype
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	_type *_type
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	hash  uint32 <span class="comment">// copy of _type.hash. Used for type switches.</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>	_     [4]byte
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>	fun   [1]uintptr <span class="comment">// variable sized. fun[0]==0 means _type does not implement inter.</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>}
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span><span class="comment">// Lock-free stack node.</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span><span class="comment">// Also known to export_test.go.</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>type lfnode struct {
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	next    uint64
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>	pushcnt uintptr
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>}
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>type forcegcstate struct {
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	lock mutex
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	g    *g
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	idle atomic.Bool
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>}
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span><span class="comment">// A _defer holds an entry on the list of deferred calls.</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span><span class="comment">// If you add a field here, add code to clear it in deferProcStack.</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span><span class="comment">// This struct must match the code in cmd/compile/internal/ssagen/ssa.go:deferstruct</span>
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span><span class="comment">// and cmd/compile/internal/ssagen/ssa.go:(*state).call.</span>
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span><span class="comment">// Some defers will be allocated on the stack and some on the heap.</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span><span class="comment">// All defers are logically part of the stack, so write barriers to</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span><span class="comment">// initialize them are not required. All defers must be manually scanned,</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span><span class="comment">// and for heap defers, marked.</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>type _defer struct {
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>	heap      bool
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	rangefunc bool    <span class="comment">// true for rangefunc list</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	sp        uintptr <span class="comment">// sp at time of defer</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>	pc        uintptr <span class="comment">// pc at time of defer</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	fn        func()  <span class="comment">// can be nil for open-coded defers</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	link      *_defer <span class="comment">// next defer on G; can point to either heap or stack!</span>
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>	<span class="comment">// If rangefunc is true, *head is the head of the atomic linked list</span>
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	<span class="comment">// during a range-over-func execution.</span>
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>	head *atomic.Pointer[_defer]
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>}
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span><span class="comment">// A _panic holds information about an active panic.</span>
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span><span class="comment">// A _panic value must only ever live on the stack.</span>
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span><span class="comment">// The argp and link fields are stack pointers, but don&#39;t need special</span>
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span><span class="comment">// handling during stack growth: because they are pointer-typed and</span>
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span><span class="comment">// _panic values only live on the stack, regular stack pointer</span>
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span><span class="comment">// adjustment takes care of them.</span>
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>type _panic struct {
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>	argp unsafe.Pointer <span class="comment">// pointer to arguments of deferred call run during panic; cannot move - known to liblink</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>	arg  any            <span class="comment">// argument to panic</span>
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>	link *_panic        <span class="comment">// link to earlier panic</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>	<span class="comment">// startPC and startSP track where _panic.start was called.</span>
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>	startPC uintptr
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>	startSP unsafe.Pointer
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>	<span class="comment">// The current stack frame that we&#39;re running deferred calls for.</span>
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>	sp unsafe.Pointer
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>	lr uintptr
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>	fp unsafe.Pointer
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>	<span class="comment">// retpc stores the PC where the panic should jump back to, if the</span>
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>	<span class="comment">// function last returned by _panic.next() recovers the panic.</span>
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>	retpc uintptr
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>	<span class="comment">// Extra state for handling open-coded defers.</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	deferBitsPtr *uint8
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>	slotsPtr     unsafe.Pointer
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>	recovered   bool <span class="comment">// whether this panic has been recovered</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	goexit      bool
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	deferreturn bool
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>}
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span><span class="comment">// savedOpenDeferState tracks the extra state from _panic that&#39;s</span>
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span><span class="comment">// necessary for deferreturn to pick up where gopanic left off,</span>
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span><span class="comment">// without needing to unwind the stack.</span>
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>type savedOpenDeferState struct {
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>	retpc           uintptr
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	deferBitsOffset uintptr
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>	slotsOffset     uintptr
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>}
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span><span class="comment">// ancestorInfo records details of where a goroutine was started.</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>type ancestorInfo struct {
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>	pcs  []uintptr <span class="comment">// pcs from the stack of this goroutine</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>	goid uint64    <span class="comment">// goroutine id of this goroutine; original goroutine possibly dead</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	gopc uintptr   <span class="comment">// pc of go statement that created this goroutine</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>}
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span><span class="comment">// A waitReason explains why a goroutine has been stopped.</span>
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span><span class="comment">// See gopark. Do not re-use waitReasons, add new ones.</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>type waitReason uint8
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>const (
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>	waitReasonZero                  waitReason = iota <span class="comment">// &#34;&#34;</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>	waitReasonGCAssistMarking                         <span class="comment">// &#34;GC assist marking&#34;</span>
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>	waitReasonIOWait                                  <span class="comment">// &#34;IO wait&#34;</span>
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>	waitReasonChanReceiveNilChan                      <span class="comment">// &#34;chan receive (nil chan)&#34;</span>
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	waitReasonChanSendNilChan                         <span class="comment">// &#34;chan send (nil chan)&#34;</span>
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>	waitReasonDumpingHeap                             <span class="comment">// &#34;dumping heap&#34;</span>
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	waitReasonGarbageCollection                       <span class="comment">// &#34;garbage collection&#34;</span>
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>	waitReasonGarbageCollectionScan                   <span class="comment">// &#34;garbage collection scan&#34;</span>
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>	waitReasonPanicWait                               <span class="comment">// &#34;panicwait&#34;</span>
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>	waitReasonSelect                                  <span class="comment">// &#34;select&#34;</span>
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>	waitReasonSelectNoCases                           <span class="comment">// &#34;select (no cases)&#34;</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>	waitReasonGCAssistWait                            <span class="comment">// &#34;GC assist wait&#34;</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>	waitReasonGCSweepWait                             <span class="comment">// &#34;GC sweep wait&#34;</span>
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>	waitReasonGCScavengeWait                          <span class="comment">// &#34;GC scavenge wait&#34;</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>	waitReasonChanReceive                             <span class="comment">// &#34;chan receive&#34;</span>
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>	waitReasonChanSend                                <span class="comment">// &#34;chan send&#34;</span>
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>	waitReasonFinalizerWait                           <span class="comment">// &#34;finalizer wait&#34;</span>
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>	waitReasonForceGCIdle                             <span class="comment">// &#34;force gc (idle)&#34;</span>
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>	waitReasonSemacquire                              <span class="comment">// &#34;semacquire&#34;</span>
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>	waitReasonSleep                                   <span class="comment">// &#34;sleep&#34;</span>
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>	waitReasonSyncCondWait                            <span class="comment">// &#34;sync.Cond.Wait&#34;</span>
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>	waitReasonSyncMutexLock                           <span class="comment">// &#34;sync.Mutex.Lock&#34;</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>	waitReasonSyncRWMutexRLock                        <span class="comment">// &#34;sync.RWMutex.RLock&#34;</span>
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>	waitReasonSyncRWMutexLock                         <span class="comment">// &#34;sync.RWMutex.Lock&#34;</span>
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>	waitReasonTraceReaderBlocked                      <span class="comment">// &#34;trace reader (blocked)&#34;</span>
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>	waitReasonWaitForGCCycle                          <span class="comment">// &#34;wait for GC cycle&#34;</span>
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>	waitReasonGCWorkerIdle                            <span class="comment">// &#34;GC worker (idle)&#34;</span>
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	waitReasonGCWorkerActive                          <span class="comment">// &#34;GC worker (active)&#34;</span>
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	waitReasonPreempted                               <span class="comment">// &#34;preempted&#34;</span>
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>	waitReasonDebugCall                               <span class="comment">// &#34;debug call&#34;</span>
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>	waitReasonGCMarkTermination                       <span class="comment">// &#34;GC mark termination&#34;</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>	waitReasonStoppingTheWorld                        <span class="comment">// &#34;stopping the world&#34;</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>	waitReasonFlushProcCaches                         <span class="comment">// &#34;flushing proc caches&#34;</span>
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>	waitReasonTraceGoroutineStatus                    <span class="comment">// &#34;trace goroutine status&#34;</span>
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>	waitReasonTraceProcStatus                         <span class="comment">// &#34;trace proc status&#34;</span>
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>	waitReasonPageTraceFlush                          <span class="comment">// &#34;page trace flush&#34;</span>
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>	waitReasonCoroutine                               <span class="comment">// &#34;coroutine&#34;</span>
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>)
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>var waitReasonStrings = [...]string{
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>	waitReasonZero:                  &#34;&#34;,
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>	waitReasonGCAssistMarking:       &#34;GC assist marking&#34;,
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>	waitReasonIOWait:                &#34;IO wait&#34;,
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>	waitReasonChanReceiveNilChan:    &#34;chan receive (nil chan)&#34;,
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>	waitReasonChanSendNilChan:       &#34;chan send (nil chan)&#34;,
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>	waitReasonDumpingHeap:           &#34;dumping heap&#34;,
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	waitReasonGarbageCollection:     &#34;garbage collection&#34;,
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	waitReasonGarbageCollectionScan: &#34;garbage collection scan&#34;,
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	waitReasonPanicWait:             &#34;panicwait&#34;,
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>	waitReasonSelect:                &#34;select&#34;,
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>	waitReasonSelectNoCases:         &#34;select (no cases)&#34;,
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>	waitReasonGCAssistWait:          &#34;GC assist wait&#34;,
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>	waitReasonGCSweepWait:           &#34;GC sweep wait&#34;,
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>	waitReasonGCScavengeWait:        &#34;GC scavenge wait&#34;,
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>	waitReasonChanReceive:           &#34;chan receive&#34;,
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>	waitReasonChanSend:              &#34;chan send&#34;,
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>	waitReasonFinalizerWait:         &#34;finalizer wait&#34;,
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>	waitReasonForceGCIdle:           &#34;force gc (idle)&#34;,
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>	waitReasonSemacquire:            &#34;semacquire&#34;,
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>	waitReasonSleep:                 &#34;sleep&#34;,
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>	waitReasonSyncCondWait:          &#34;sync.Cond.Wait&#34;,
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>	waitReasonSyncMutexLock:         &#34;sync.Mutex.Lock&#34;,
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>	waitReasonSyncRWMutexRLock:      &#34;sync.RWMutex.RLock&#34;,
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>	waitReasonSyncRWMutexLock:       &#34;sync.RWMutex.Lock&#34;,
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>	waitReasonTraceReaderBlocked:    &#34;trace reader (blocked)&#34;,
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>	waitReasonWaitForGCCycle:        &#34;wait for GC cycle&#34;,
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>	waitReasonGCWorkerIdle:          &#34;GC worker (idle)&#34;,
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>	waitReasonGCWorkerActive:        &#34;GC worker (active)&#34;,
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>	waitReasonPreempted:             &#34;preempted&#34;,
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>	waitReasonDebugCall:             &#34;debug call&#34;,
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>	waitReasonGCMarkTermination:     &#34;GC mark termination&#34;,
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>	waitReasonStoppingTheWorld:      &#34;stopping the world&#34;,
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>	waitReasonFlushProcCaches:       &#34;flushing proc caches&#34;,
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>	waitReasonTraceGoroutineStatus:  &#34;trace goroutine status&#34;,
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>	waitReasonTraceProcStatus:       &#34;trace proc status&#34;,
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>	waitReasonPageTraceFlush:        &#34;page trace flush&#34;,
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>	waitReasonCoroutine:             &#34;coroutine&#34;,
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>}
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>func (w waitReason) String() string {
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>	if w &lt; 0 || w &gt;= waitReason(len(waitReasonStrings)) {
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>		return &#34;unknown wait reason&#34;
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>	}
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>	return waitReasonStrings[w]
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>}
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>func (w waitReason) isMutexWait() bool {
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>	return w == waitReasonSyncMutexLock ||
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>		w == waitReasonSyncRWMutexRLock ||
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>		w == waitReasonSyncRWMutexLock
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>}
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>var (
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>	allm       *m
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>	gomaxprocs int32
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>	ncpu       int32
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>	forcegc    forcegcstate
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>	sched      schedt
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>	newprocs   int32
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>	<span class="comment">// allpLock protects P-less reads and size changes of allp, idlepMask,</span>
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>	<span class="comment">// and timerpMask, and all writes to allp.</span>
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>	allpLock mutex
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>	<span class="comment">// len(allp) == gomaxprocs; may change at safe points, otherwise</span>
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>	<span class="comment">// immutable.</span>
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>	allp []*p
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>	<span class="comment">// Bitmask of Ps in _Pidle list, one bit per P. Reads and writes must</span>
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>	<span class="comment">// be atomic. Length may change at safe points.</span>
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>	<span class="comment">// Each P must update only its own bit. In order to maintain</span>
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>	<span class="comment">// consistency, a P going idle must the idle mask simultaneously with</span>
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>	<span class="comment">// updates to the idle P list under the sched.lock, otherwise a racing</span>
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>	<span class="comment">// pidleget may clear the mask before pidleput sets the mask,</span>
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>	<span class="comment">// corrupting the bitmap.</span>
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>	<span class="comment">// N.B., procresize takes ownership of all Ps in stopTheWorldWithSema.</span>
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>	idlepMask pMask
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>	<span class="comment">// Bitmask of Ps that may have a timer, one bit per P. Reads and writes</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>	<span class="comment">// must be atomic. Length may change at safe points.</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>	timerpMask pMask
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>	<span class="comment">// Pool of GC parked background workers. Entries are type</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>	<span class="comment">// *gcBgMarkWorkerNode.</span>
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>	gcBgMarkWorkerPool lfstack
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>	<span class="comment">// Total number of gcBgMarkWorker goroutines. Protected by worldsema.</span>
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>	gcBgMarkWorkerCount int32
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>	<span class="comment">// Information about what cpu features are available.</span>
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>	<span class="comment">// Packages outside the runtime should not use these</span>
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>	<span class="comment">// as they are not an external api.</span>
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>	<span class="comment">// Set on startup in asm_{386,amd64}.s</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>	processorVersionInfo uint32
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>	isIntel              bool
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>	<span class="comment">// set by cmd/link on arm systems</span>
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	goarm       uint8
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>	goarmsoftfp uint8
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>)
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span><span class="comment">// Set by the linker so the runtime can determine the buildmode.</span>
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>var (
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>	islibrary bool <span class="comment">// -buildmode=c-shared</span>
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>	isarchive bool <span class="comment">// -buildmode=c-archive</span>
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>)
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span><span class="comment">// Must agree with internal/buildcfg.FramePointerEnabled.</span>
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>const framepointer_enabled = GOARCH == &#34;amd64&#34; || GOARCH == &#34;arm64&#34;
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>
</pre><p><a href="runtime2.go?m=text">View as plain text</a></p>

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
