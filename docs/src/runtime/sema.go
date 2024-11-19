<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/sema.go - Go Documentation Server</title>

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
<a href="sema.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">sema.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Semaphore implementation exposed to Go.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// Intended use is provide a sleep and wakeup</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// primitive that can be used in the contended case</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// of other synchronization primitives.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// Thus it targets the same goal as Linux&#39;s futex,</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// but it has much simpler semantics.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// That is, don&#39;t think of these as semaphores.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// Think of them as a way to implement sleep and wakeup</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// such that every sleep is paired with a single wakeup,</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// even if, due to races, the wakeup happens before the sleep.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// See Mullender and Cox, ``Semaphores in Plan 9,&#39;&#39;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// https://swtch.com/semaphore.pdf</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>package runtime
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>import (
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	&#34;internal/cpu&#34;
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// Asynchronous semaphore for sync.Mutex.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// A semaRoot holds a balanced tree of sudog with distinct addresses (s.elem).</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// Each of those sudog may in turn point (through s.waitlink) to a list</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// of other sudogs waiting on the same address.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// The operations on the inner lists of sudogs with the same address</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// are all O(1). The scanning of the top-level semaRoot list is O(log n),</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// where n is the number of distinct addresses with goroutines blocked</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// on them that hash to the given semaRoot.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// See golang.org/issue/17953 for a program that worked badly</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// before we introduced the second level of list, and</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// BenchmarkSemTable/OneAddrCollision/* for a benchmark that exercises this.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>type semaRoot struct {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	lock  mutex
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	treap *sudog        <span class="comment">// root of balanced tree of unique waiters.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	nwait atomic.Uint32 <span class="comment">// Number of waiters. Read w/o the lock.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>var semtable semTable
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// Prime to not correlate with any user patterns.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>const semTabSize = 251
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>type semTable [semTabSize]struct {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	root semaRoot
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	pad  [cpu.CacheLinePadSize - unsafe.Sizeof(semaRoot{})]byte
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>func (t *semTable) rootFor(addr *uint32) *semaRoot {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	return &amp;t[(uintptr(unsafe.Pointer(addr))&gt;&gt;3)%semTabSize].root
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//go:linkname sync_runtime_Semacquire sync.runtime_Semacquire</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>func sync_runtime_Semacquire(addr *uint32) {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	semacquire1(addr, false, semaBlockProfile, 0, waitReasonSemacquire)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_Semacquire internal/poll.runtime_Semacquire</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>func poll_runtime_Semacquire(addr *uint32) {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	semacquire1(addr, false, semaBlockProfile, 0, waitReasonSemacquire)
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">//go:linkname sync_runtime_Semrelease sync.runtime_Semrelease</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>func sync_runtime_Semrelease(addr *uint32, handoff bool, skipframes int) {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	semrelease1(addr, handoff, skipframes)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//go:linkname sync_runtime_SemacquireMutex sync.runtime_SemacquireMutex</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>func sync_runtime_SemacquireMutex(addr *uint32, lifo bool, skipframes int) {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	semacquire1(addr, lifo, semaBlockProfile|semaMutexProfile, skipframes, waitReasonSyncMutexLock)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//go:linkname sync_runtime_SemacquireRWMutexR sync.runtime_SemacquireRWMutexR</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>func sync_runtime_SemacquireRWMutexR(addr *uint32, lifo bool, skipframes int) {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	semacquire1(addr, lifo, semaBlockProfile|semaMutexProfile, skipframes, waitReasonSyncRWMutexRLock)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//go:linkname sync_runtime_SemacquireRWMutex sync.runtime_SemacquireRWMutex</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>func sync_runtime_SemacquireRWMutex(addr *uint32, lifo bool, skipframes int) {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	semacquire1(addr, lifo, semaBlockProfile|semaMutexProfile, skipframes, waitReasonSyncRWMutexLock)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_Semrelease internal/poll.runtime_Semrelease</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func poll_runtime_Semrelease(addr *uint32) {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	semrelease(addr)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>func readyWithTime(s *sudog, traceskip int) {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if s.releasetime != 0 {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		s.releasetime = cputicks()
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	goready(s.g, traceskip)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>type semaProfileFlags int
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>const (
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	semaBlockProfile semaProfileFlags = 1 &lt;&lt; iota
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	semaMutexProfile
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// Called from runtime.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>func semacquire(addr *uint32) {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	semacquire1(addr, false, 0, 0, waitReasonSemacquire)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>func semacquire1(addr *uint32, lifo bool, profile semaProfileFlags, skipframes int, reason waitReason) {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	gp := getg()
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	if gp != gp.m.curg {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		throw(&#34;semacquire not on the G stack&#34;)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// Easy case.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	if cansemacquire(addr) {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		return
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// Harder case:</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">//	increment waiter count</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">//	try cansemacquire one more time, return if succeeded</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">//	enqueue itself as a waiter</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">//	sleep</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">//	(waiter descriptor is dequeued by signaler)</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	s := acquireSudog()
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	root := semtable.rootFor(addr)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	t0 := int64(0)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	s.releasetime = 0
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	s.acquiretime = 0
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	s.ticket = 0
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	if profile&amp;semaBlockProfile != 0 &amp;&amp; blockprofilerate &gt; 0 {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		t0 = cputicks()
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		s.releasetime = -1
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	if profile&amp;semaMutexProfile != 0 &amp;&amp; mutexprofilerate &gt; 0 {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		if t0 == 0 {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			t0 = cputicks()
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		s.acquiretime = t0
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	for {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		lockWithRank(&amp;root.lock, lockRankRoot)
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		<span class="comment">// Add ourselves to nwait to disable &#34;easy case&#34; in semrelease.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		root.nwait.Add(1)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		<span class="comment">// Check cansemacquire to avoid missed wakeup.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		if cansemacquire(addr) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			root.nwait.Add(-1)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			unlock(&amp;root.lock)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			break
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// Any semrelease after the cansemacquire knows we&#39;re waiting</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// (we set nwait above), so go to sleep.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		root.queue(addr, s, lifo)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		goparkunlock(&amp;root.lock, reason, traceBlockSync, 4+skipframes)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		if s.ticket != 0 || cansemacquire(addr) {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			break
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	if s.releasetime &gt; 0 {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		blockevent(s.releasetime-t0, 3+skipframes)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	releaseSudog(s)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>func semrelease(addr *uint32) {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	semrelease1(addr, false, 0)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>func semrelease1(addr *uint32, handoff bool, skipframes int) {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	root := semtable.rootFor(addr)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	atomic.Xadd(addr, 1)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// Easy case: no waiters?</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// This check must happen after the xadd, to avoid a missed wakeup</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// (see loop in semacquire).</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	if root.nwait.Load() == 0 {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		return
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// Harder case: search for a waiter and wake it.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	lockWithRank(&amp;root.lock, lockRankRoot)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	if root.nwait.Load() == 0 {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		<span class="comment">// The count is already consumed by another goroutine,</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		<span class="comment">// so no need to wake up another goroutine.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		unlock(&amp;root.lock)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		return
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	s, t0, tailtime := root.dequeue(addr)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	if s != nil {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		root.nwait.Add(-1)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	unlock(&amp;root.lock)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	if s != nil { <span class="comment">// May be slow or even yield, so unlock first</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		acquiretime := s.acquiretime
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		if acquiretime != 0 {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>			<span class="comment">// Charge contention that this (delayed) unlock caused.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			<span class="comment">// If there are N more goroutines waiting beyond the</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			<span class="comment">// one that&#39;s waking up, charge their delay as well, so that</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			<span class="comment">// contention holding up many goroutines shows up as</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			<span class="comment">// more costly than contention holding up a single goroutine.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			<span class="comment">// It would take O(N) time to calculate how long each goroutine</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			<span class="comment">// has been waiting, so instead we charge avg(head-wait, tail-wait)*N.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			<span class="comment">// head-wait is the longest wait and tail-wait is the shortest.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			<span class="comment">// (When we do a lifo insertion, we preserve this property by</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			<span class="comment">// copying the old head&#39;s acquiretime into the inserted new head.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			<span class="comment">// In that case the overall average may be slightly high, but that&#39;s fine:</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			<span class="comment">// the average of the ends is only an approximation to the actual</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			<span class="comment">// average anyway.)</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			<span class="comment">// The root.dequeue above changed the head and tail acquiretime</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			<span class="comment">// to the current time, so the next unlock will not re-count this contention.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			dt0 := t0 - acquiretime
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			dt := dt0
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			if s.waiters != 0 {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>				dtail := t0 - tailtime
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>				dt += (dtail + dt0) / 2 * int64(s.waiters)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			mutexevent(dt, 3+skipframes)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		if s.ticket != 0 {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			throw(&#34;corrupted semaphore ticket&#34;)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		if handoff &amp;&amp; cansemacquire(addr) {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			s.ticket = 1
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		readyWithTime(s, 5+skipframes)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		if s.ticket == 1 &amp;&amp; getg().m.locks == 0 {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			<span class="comment">// Direct G handoff</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			<span class="comment">// readyWithTime has added the waiter G as runnext in the</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			<span class="comment">// current P; we now call the scheduler so that we start running</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			<span class="comment">// the waiter G immediately.</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			<span class="comment">// Note that waiter inherits our time slice: this is desirable</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			<span class="comment">// to avoid having a highly contended semaphore hog the P</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			<span class="comment">// indefinitely. goyield is like Gosched, but it emits a</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			<span class="comment">// &#34;preempted&#34; trace event instead and, more importantly, puts</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			<span class="comment">// the current G on the local runq instead of the global one.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			<span class="comment">// We only do this in the starving regime (handoff=true), as in</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			<span class="comment">// the non-starving case it is possible for a different waiter</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			<span class="comment">// to acquire the semaphore while we are yielding/scheduling,</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			<span class="comment">// and this would be wasteful. We wait instead to enter starving</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			<span class="comment">// regime, and then we start to do direct handoffs of ticket and</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			<span class="comment">// P.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			<span class="comment">// See issue 33747 for discussion.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			goyield()
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>func cansemacquire(addr *uint32) bool {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	for {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		v := atomic.Load(addr)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		if v == 0 {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			return false
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		if atomic.Cas(addr, v, v-1) {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			return true
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">// queue adds s to the blocked goroutines in semaRoot.</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>func (root *semaRoot) queue(addr *uint32, s *sudog, lifo bool) {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	s.g = getg()
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	s.elem = unsafe.Pointer(addr)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	s.next = nil
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	s.prev = nil
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	s.waiters = 0
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	var last *sudog
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	pt := &amp;root.treap
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	for t := *pt; t != nil; t = *pt {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		if t.elem == unsafe.Pointer(addr) {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			<span class="comment">// Already have addr in list.</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			if lifo {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>				<span class="comment">// Substitute s in t&#39;s place in treap.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>				*pt = s
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>				s.ticket = t.ticket
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>				s.acquiretime = t.acquiretime <span class="comment">// preserve head acquiretime as oldest time</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>				s.parent = t.parent
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>				s.prev = t.prev
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>				s.next = t.next
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>				if s.prev != nil {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>					s.prev.parent = s
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>				}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>				if s.next != nil {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>					s.next.parent = s
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>				}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>				<span class="comment">// Add t first in s&#39;s wait list.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>				s.waitlink = t
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>				s.waittail = t.waittail
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>				if s.waittail == nil {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>					s.waittail = t
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>				}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				s.waiters = t.waiters
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>				if s.waiters+1 != 0 {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>					s.waiters++
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>				}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>				t.parent = nil
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>				t.prev = nil
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>				t.next = nil
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>				t.waittail = nil
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			} else {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>				<span class="comment">// Add s to end of t&#39;s wait list.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>				if t.waittail == nil {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>					t.waitlink = s
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>				} else {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>					t.waittail.waitlink = s
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>				}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>				t.waittail = s
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>				s.waitlink = nil
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>				if t.waiters+1 != 0 {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>					t.waiters++
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>				}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			return
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		last = t
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		if uintptr(unsafe.Pointer(addr)) &lt; uintptr(t.elem) {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			pt = &amp;t.prev
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		} else {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			pt = &amp;t.next
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">// Add s as new leaf in tree of unique addrs.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	<span class="comment">// The balanced tree is a treap using ticket as the random heap priority.</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">// That is, it is a binary tree ordered according to the elem addresses,</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// but then among the space of possible binary trees respecting those</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// addresses, it is kept balanced on average by maintaining a heap ordering</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// on the ticket: s.ticket &lt;= both s.prev.ticket and s.next.ticket.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// https://en.wikipedia.org/wiki/Treap</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// https://faculty.washington.edu/aragon/pubs/rst89.pdf</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// s.ticket compared with zero in couple of places, therefore set lowest bit.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">// It will not affect treap&#39;s quality noticeably.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	s.ticket = cheaprand() | 1
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	s.parent = last
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	*pt = s
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	<span class="comment">// Rotate up into tree according to ticket (priority).</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	for s.parent != nil &amp;&amp; s.parent.ticket &gt; s.ticket {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		if s.parent.prev == s {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>			root.rotateRight(s.parent)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		} else {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>			if s.parent.next != s {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>				panic(&#34;semaRoot queue&#34;)
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>			root.rotateLeft(s.parent)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	}
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">// dequeue searches for and finds the first goroutine</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// in semaRoot blocked on addr.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span><span class="comment">// If the sudog was being profiled, dequeue returns the time</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// at which it was woken up as now. Otherwise now is 0.</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// If there are additional entries in the wait list, dequeue</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span><span class="comment">// returns tailtime set to the last entry&#39;s acquiretime.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">// Otherwise tailtime is found.acquiretime.</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>func (root *semaRoot) dequeue(addr *uint32) (found *sudog, now, tailtime int64) {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	ps := &amp;root.treap
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	s := *ps
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	for ; s != nil; s = *ps {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		if s.elem == unsafe.Pointer(addr) {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			goto Found
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		if uintptr(unsafe.Pointer(addr)) &lt; uintptr(s.elem) {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			ps = &amp;s.prev
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		} else {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			ps = &amp;s.next
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	return nil, 0, 0
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>Found:
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	now = int64(0)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	if s.acquiretime != 0 {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		now = cputicks()
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	if t := s.waitlink; t != nil {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		<span class="comment">// Substitute t, also waiting on addr, for s in root tree of unique addrs.</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		*ps = t
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		t.ticket = s.ticket
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		t.parent = s.parent
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		t.prev = s.prev
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		if t.prev != nil {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			t.prev.parent = t
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		t.next = s.next
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		if t.next != nil {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			t.next.parent = t
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		}
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		if t.waitlink != nil {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			t.waittail = s.waittail
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		} else {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			t.waittail = nil
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		t.waiters = s.waiters
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		if t.waiters &gt; 1 {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			t.waiters--
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		<span class="comment">// Set head and tail acquire time to &#39;now&#39;,</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		<span class="comment">// because the caller will take care of charging</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		<span class="comment">// the delays before now for all entries in the list.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		t.acquiretime = now
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		tailtime = s.waittail.acquiretime
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		s.waittail.acquiretime = now
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		s.waitlink = nil
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		s.waittail = nil
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	} else {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		<span class="comment">// Rotate s down to be leaf of tree for removal, respecting priorities.</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		for s.next != nil || s.prev != nil {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			if s.next == nil || s.prev != nil &amp;&amp; s.prev.ticket &lt; s.next.ticket {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>				root.rotateRight(s)
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			} else {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>				root.rotateLeft(s)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		<span class="comment">// Remove s, now a leaf.</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		if s.parent != nil {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			if s.parent.prev == s {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>				s.parent.prev = nil
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			} else {
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>				s.parent.next = nil
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		} else {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			root.treap = nil
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		tailtime = s.acquiretime
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	s.parent = nil
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	s.elem = nil
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	s.next = nil
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	s.prev = nil
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	s.ticket = 0
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	return s, now, tailtime
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span><span class="comment">// rotateLeft rotates the tree rooted at node x.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span><span class="comment">// turning (x a (y b c)) into (y (x a b) c).</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>func (root *semaRoot) rotateLeft(x *sudog) {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	<span class="comment">// p -&gt; (x a (y b c))</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	p := x.parent
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	y := x.next
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	b := y.prev
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	y.prev = x
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	x.parent = y
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	x.next = b
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	if b != nil {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		b.parent = x
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	y.parent = p
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	if p == nil {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		root.treap = y
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	} else if p.prev == x {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		p.prev = y
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	} else {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		if p.next != x {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			throw(&#34;semaRoot rotateLeft&#34;)
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		p.next = y
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span><span class="comment">// rotateRight rotates the tree rooted at node y.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span><span class="comment">// turning (y (x a b) c) into (x a (y b c)).</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>func (root *semaRoot) rotateRight(y *sudog) {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	<span class="comment">// p -&gt; (y (x a b) c)</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	p := y.parent
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	x := y.prev
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	b := x.next
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	x.next = y
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	y.parent = x
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	y.prev = b
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	if b != nil {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		b.parent = y
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	x.parent = p
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	if p == nil {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		root.treap = x
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	} else if p.prev == y {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		p.prev = x
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	} else {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		if p.next != y {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			throw(&#34;semaRoot rotateRight&#34;)
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		p.next = x
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	}
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span><span class="comment">// notifyList is a ticket-based notification list used to implement sync.Cond.</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span><span class="comment">// It must be kept in sync with the sync package.</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>type notifyList struct {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	<span class="comment">// wait is the ticket number of the next waiter. It is atomically</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	<span class="comment">// incremented outside the lock.</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	wait atomic.Uint32
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	<span class="comment">// notify is the ticket number of the next waiter to be notified. It can</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	<span class="comment">// be read outside the lock, but is only written to with lock held.</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	<span class="comment">// Both wait &amp; notify can wrap around, and such cases will be correctly</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	<span class="comment">// handled as long as their &#34;unwrapped&#34; difference is bounded by 2^31.</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	<span class="comment">// For this not to be the case, we&#39;d need to have 2^31+ goroutines</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	<span class="comment">// blocked on the same condvar, which is currently not possible.</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	notify uint32
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	<span class="comment">// List of parked waiters.</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	lock mutex
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	head *sudog
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	tail *sudog
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span><span class="comment">// less checks if a &lt; b, considering a &amp; b running counts that may overflow the</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">// 32-bit range, and that their &#34;unwrapped&#34; difference is always less than 2^31.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>func less(a, b uint32) bool {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	return int32(a-b) &lt; 0
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>}
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span><span class="comment">// notifyListAdd adds the caller to a notify list such that it can receive</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span><span class="comment">// notifications. The caller must eventually call notifyListWait to wait for</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span><span class="comment">// such a notification, passing the returned ticket number.</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span><span class="comment">//go:linkname notifyListAdd sync.runtime_notifyListAdd</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>func notifyListAdd(l *notifyList) uint32 {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	<span class="comment">// This may be called concurrently, for example, when called from</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	<span class="comment">// sync.Cond.Wait while holding a RWMutex in read mode.</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	return l.wait.Add(1) - 1
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>}
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">// notifyListWait waits for a notification. If one has been sent since</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">// notifyListAdd was called, it returns immediately. Otherwise, it blocks.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">//go:linkname notifyListWait sync.runtime_notifyListWait</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>func notifyListWait(l *notifyList, t uint32) {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	lockWithRank(&amp;l.lock, lockRankNotifyList)
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	<span class="comment">// Return right away if this ticket has already been notified.</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	if less(t, l.notify) {
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		unlock(&amp;l.lock)
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		return
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	<span class="comment">// Enqueue itself.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	s := acquireSudog()
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	s.g = getg()
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	s.ticket = t
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	s.releasetime = 0
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	t0 := int64(0)
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	if blockprofilerate &gt; 0 {
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		t0 = cputicks()
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		s.releasetime = -1
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	if l.tail == nil {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		l.head = s
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	} else {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		l.tail.next = s
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	l.tail = s
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	goparkunlock(&amp;l.lock, waitReasonSyncCondWait, traceBlockCondWait, 3)
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	if t0 != 0 {
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		blockevent(s.releasetime-t0, 2)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	}
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	releaseSudog(s)
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span><span class="comment">// notifyListNotifyAll notifies all entries in the list.</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span><span class="comment">//go:linkname notifyListNotifyAll sync.runtime_notifyListNotifyAll</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>func notifyListNotifyAll(l *notifyList) {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	<span class="comment">// Fast-path: if there are no new waiters since the last notification</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	<span class="comment">// we don&#39;t need to acquire the lock.</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	if l.wait.Load() == atomic.Load(&amp;l.notify) {
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		return
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	<span class="comment">// Pull the list out into a local variable, waiters will be readied</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	<span class="comment">// outside the lock.</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	lockWithRank(&amp;l.lock, lockRankNotifyList)
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	s := l.head
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	l.head = nil
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	l.tail = nil
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	<span class="comment">// Update the next ticket to be notified. We can set it to the current</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	<span class="comment">// value of wait because any previous waiters are already in the list</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	<span class="comment">// or will notice that they have already been notified when trying to</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	<span class="comment">// add themselves to the list.</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	atomic.Store(&amp;l.notify, l.wait.Load())
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	unlock(&amp;l.lock)
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	<span class="comment">// Go through the local list and ready all waiters.</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	for s != nil {
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		next := s.next
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		s.next = nil
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		readyWithTime(s, 4)
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		s = next
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>}
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span><span class="comment">// notifyListNotifyOne notifies one entry in the list.</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span><span class="comment">//go:linkname notifyListNotifyOne sync.runtime_notifyListNotifyOne</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>func notifyListNotifyOne(l *notifyList) {
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	<span class="comment">// Fast-path: if there are no new waiters since the last notification</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	<span class="comment">// we don&#39;t need to acquire the lock at all.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	if l.wait.Load() == atomic.Load(&amp;l.notify) {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		return
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	}
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	lockWithRank(&amp;l.lock, lockRankNotifyList)
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	<span class="comment">// Re-check under the lock if we need to do anything.</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	t := l.notify
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	if t == l.wait.Load() {
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		unlock(&amp;l.lock)
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		return
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	<span class="comment">// Update the next notify ticket number.</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	atomic.Store(&amp;l.notify, t+1)
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	<span class="comment">// Try to find the g that needs to be notified.</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	<span class="comment">// If it hasn&#39;t made it to the list yet we won&#39;t find it,</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	<span class="comment">// but it won&#39;t park itself once it sees the new notify number.</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	<span class="comment">// This scan looks linear but essentially always stops quickly.</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	<span class="comment">// Because g&#39;s queue separately from taking numbers,</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	<span class="comment">// there may be minor reorderings in the list, but we</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	<span class="comment">// expect the g we&#39;re looking for to be near the front.</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	<span class="comment">// The g has others in front of it on the list only to the</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	<span class="comment">// extent that it lost the race, so the iteration will not</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	<span class="comment">// be too long. This applies even when the g is missing:</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	<span class="comment">// it hasn&#39;t yet gotten to sleep and has lost the race to</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	<span class="comment">// the (few) other g&#39;s that we find on the list.</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	for p, s := (*sudog)(nil), l.head; s != nil; p, s = s, s.next {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		if s.ticket == t {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>			n := s.next
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>			if p != nil {
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>				p.next = n
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>			} else {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>				l.head = n
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>			}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>			if n == nil {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>				l.tail = p
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>			}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>			unlock(&amp;l.lock)
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>			s.next = nil
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>			readyWithTime(s, 4)
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>			return
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		}
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	}
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	unlock(&amp;l.lock)
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>}
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span><span class="comment">//go:linkname notifyListCheck sync.runtime_notifyListCheck</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>func notifyListCheck(sz uintptr) {
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	if sz != unsafe.Sizeof(notifyList{}) {
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		print(&#34;runtime: bad notifyList size - sync=&#34;, sz, &#34; runtime=&#34;, unsafe.Sizeof(notifyList{}), &#34;\n&#34;)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		throw(&#34;bad notifyList size&#34;)
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	}
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span><span class="comment">//go:linkname sync_nanotime sync.runtime_nanotime</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>func sync_nanotime() int64 {
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	return nanotime()
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>}
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>
</pre><p><a href="sema.go?m=text">View as plain text</a></p>

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
