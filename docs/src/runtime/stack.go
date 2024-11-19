<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/stack.go - Go Documentation Server</title>

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
<a href="stack.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">stack.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2013 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/cpu&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/goos&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">/*
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>Stack layout parameters.
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>Included both by runtime (compiled via 6c) and linkers (compiled via gcc).
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>The per-goroutine g-&gt;stackguard is set to point StackGuard bytes
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>above the bottom of the stack.  Each function compares its stack
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>pointer against g-&gt;stackguard to check for overflow.  To cut one
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>instruction from the check sequence for functions with tiny frames,
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>the stack is allowed to protrude StackSmall bytes below the stack
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>guard.  Functions with large frames don&#39;t bother with the check and
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>always call morestack.  The sequences are (for amd64, others are
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>similar):
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	guard = g-&gt;stackguard
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	frame = function&#39;s stack frame size
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	argsize = size of function arguments (call + return)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	stack frame size &lt;= StackSmall:
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		CMPQ guard, SP
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		JHI 3(PC)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		MOVQ m-&gt;morearg, $(argsize &lt;&lt; 32)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		CALL morestack(SB)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	stack frame size &gt; StackSmall but &lt; StackBig
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		LEAQ (frame-StackSmall)(SP), R0
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		CMPQ guard, R0
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		JHI 3(PC)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		MOVQ m-&gt;morearg, $(argsize &lt;&lt; 32)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		CALL morestack(SB)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	stack frame size &gt;= StackBig:
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		MOVQ m-&gt;morearg, $((argsize &lt;&lt; 32) | frame)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		CALL morestack(SB)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>The bottom StackGuard - StackSmall bytes are important: there has
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>to be enough room to execute functions that refuse to check for
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>stack overflow, either because they need to be adjacent to the
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>actual caller&#39;s frame (deferproc) or because they handle the imminent
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>stack overflow (morestack).
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>For example, deferproc might call malloc, which does one of the
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>above checks (without allocating a full frame), which might trigger
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>a call to morestack.  This sequence needs to fit in the bottom
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>section of the stack.  On amd64, morestack&#39;s frame is 40 bytes, and
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>deferproc&#39;s frame is 56 bytes.  That fits well within the
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>StackGuard - StackSmall bytes at the bottom.
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>The linkers explore all possible call traces involving non-splitting
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>functions to make sure that this limit cannot be violated.
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>*/</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>const (
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// stackSystem is a number of additional bytes to add</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// to each stack below the usual guard area for OS-specific</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// purposes like signal handling. Used on Windows, Plan 9,</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// and iOS because they do not use a separate stack.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	stackSystem = goos.IsWindows*512*goarch.PtrSize + goos.IsPlan9*512 + goos.IsIos*goarch.IsArm64*1024
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// The minimum size of stack used by Go code</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	stackMin = 2048
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// The minimum stack size to allocate.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// The hackery here rounds fixedStack0 up to a power of 2.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	fixedStack0 = stackMin + stackSystem
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	fixedStack1 = fixedStack0 - 1
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	fixedStack2 = fixedStack1 | (fixedStack1 &gt;&gt; 1)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	fixedStack3 = fixedStack2 | (fixedStack2 &gt;&gt; 2)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	fixedStack4 = fixedStack3 | (fixedStack3 &gt;&gt; 4)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	fixedStack5 = fixedStack4 | (fixedStack4 &gt;&gt; 8)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	fixedStack6 = fixedStack5 | (fixedStack5 &gt;&gt; 16)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	fixedStack  = fixedStack6 + 1
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// stackNosplit is the maximum number of bytes that a chain of NOSPLIT</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// functions can use.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// This arithmetic must match that in cmd/internal/objabi/stack.go:StackNosplit.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	stackNosplit = abi.StackNosplitBase * sys.StackGuardMultiplier
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// The stack guard is a pointer this many bytes above the</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// bottom of the stack.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// The guard leaves enough room for a stackNosplit chain of NOSPLIT calls</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// plus one stackSmall frame plus stackSystem bytes for the OS.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// This arithmetic must match that in cmd/internal/objabi/stack.go:StackLimit.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	stackGuard = stackNosplit + stackSystem + abi.StackSmall
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>)
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>const (
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// stackDebug == 0: no logging</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">//            == 1: logging of per-stack operations</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">//            == 2: logging of per-frame operations</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	<span class="comment">//            == 3: logging of per-word updates</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">//            == 4: logging of per-word reads</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	stackDebug       = 0
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	stackFromSystem  = 0 <span class="comment">// allocate stacks from system memory instead of the heap</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	stackFaultOnFree = 0 <span class="comment">// old stacks are mapped noaccess to detect use after free</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	stackNoCache     = 0 <span class="comment">// disable per-P small stack caches</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// check the BP links during traceback.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	debugCheckBP = false
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>var (
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	stackPoisonCopy = 0 <span class="comment">// fill stack that should not be accessed with garbage, to detect bad dereferences during copy</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>const (
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	uintptrMask = 1&lt;&lt;(8*goarch.PtrSize) - 1
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// The values below can be stored to g.stackguard0 to force</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// the next stack check to fail.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// These are all larger than any real SP.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// Goroutine preemption request.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// 0xfffffade in hex.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	stackPreempt = uintptrMask &amp; -1314
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// Thread is forking. Causes a split stack check failure.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// 0xfffffb2e in hex.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	stackFork = uintptrMask &amp; -1234
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// Force a stack movement. Used for debugging.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// 0xfffffeed in hex.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	stackForceMove = uintptrMask &amp; -275
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// stackPoisonMin is the lowest allowed stack poison value.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	stackPoisonMin = uintptrMask &amp; -4096
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// Global pool of spans that have free stacks.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// Stacks are assigned an order according to size.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">//	order = log_2(size/FixedStack)</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// There is a free list for each order.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>var stackpool [_NumStackOrders]struct {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	item stackpoolItem
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	_    [(cpu.CacheLinePadSize - unsafe.Sizeof(stackpoolItem{})%cpu.CacheLinePadSize) % cpu.CacheLinePadSize]byte
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>type stackpoolItem struct {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	_    sys.NotInHeap
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	mu   mutex
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	span mSpanList
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">// Global pool of large stack spans.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>var stackLarge struct {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	lock mutex
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	free [heapAddrBits - pageShift]mSpanList <span class="comment">// free lists by log_2(s.npages)</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>func stackinit() {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	if _StackCacheSize&amp;_PageMask != 0 {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		throw(&#34;cache size must be a multiple of page size&#34;)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	for i := range stackpool {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		stackpool[i].item.span.init()
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		lockInit(&amp;stackpool[i].item.mu, lockRankStackpool)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	for i := range stackLarge.free {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		stackLarge.free[i].init()
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		lockInit(&amp;stackLarge.lock, lockRankStackLarge)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">// stacklog2 returns ⌊log_2(n)⌋.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>func stacklog2(n uintptr) int {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	log2 := 0
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	for n &gt; 1 {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		n &gt;&gt;= 1
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		log2++
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	return log2
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// Allocates a stack from the free pool. Must be called with</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// stackpool[order].item.mu held.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>func stackpoolalloc(order uint8) gclinkptr {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	list := &amp;stackpool[order].item.span
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	s := list.first
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	lockWithRankMayAcquire(&amp;mheap_.lock, lockRankMheap)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	if s == nil {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		<span class="comment">// no free stacks. Allocate another span worth.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		s = mheap_.allocManual(_StackCacheSize&gt;&gt;_PageShift, spanAllocStack)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		if s == nil {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			throw(&#34;out of memory&#34;)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		if s.allocCount != 0 {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			throw(&#34;bad allocCount&#34;)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		if s.manualFreeList.ptr() != nil {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			throw(&#34;bad manualFreeList&#34;)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		osStackAlloc(s)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		s.elemsize = fixedStack &lt;&lt; order
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; _StackCacheSize; i += s.elemsize {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			x := gclinkptr(s.base() + i)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			x.ptr().next = s.manualFreeList
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			s.manualFreeList = x
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		list.insert(s)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	x := s.manualFreeList
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	if x.ptr() == nil {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		throw(&#34;span has no free stacks&#34;)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	s.manualFreeList = x.ptr().next
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	s.allocCount++
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	if s.manualFreeList.ptr() == nil {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		<span class="comment">// all stacks in s are allocated.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		list.remove(s)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	return x
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// Adds stack x to the free pool. Must be called with stackpool[order].item.mu held.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>func stackpoolfree(x gclinkptr, order uint8) {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	s := spanOfUnchecked(uintptr(x))
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	if s.state.get() != mSpanManual {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		throw(&#34;freeing stack not in a stack span&#34;)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	if s.manualFreeList.ptr() == nil {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		<span class="comment">// s will now have a free stack</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		stackpool[order].item.span.insert(s)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	x.ptr().next = s.manualFreeList
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	s.manualFreeList = x
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	s.allocCount--
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	if gcphase == _GCoff &amp;&amp; s.allocCount == 0 {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		<span class="comment">// Span is completely free. Return it to the heap</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		<span class="comment">// immediately if we&#39;re sweeping.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		<span class="comment">// If GC is active, we delay the free until the end of</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		<span class="comment">// GC to avoid the following type of situation:</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		<span class="comment">// 1) GC starts, scans a SudoG but does not yet mark the SudoG.elem pointer</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		<span class="comment">// 2) The stack that pointer points to is copied</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		<span class="comment">// 3) The old stack is freed</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		<span class="comment">// 4) The containing span is marked free</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		<span class="comment">// 5) GC attempts to mark the SudoG.elem pointer. The</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		<span class="comment">//    marking fails because the pointer looks like a</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		<span class="comment">//    pointer into a free span.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		<span class="comment">// By not freeing, we prevent step #4 until GC is done.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		stackpool[order].item.span.remove(s)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		s.manualFreeList = 0
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		osStackFree(s)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		mheap_.freeManual(s, spanAllocStack)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span><span class="comment">// stackcacherefill/stackcacherelease implement a global pool of stack segments.</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">// The pool is required to prevent unlimited growth of per-thread caches.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>func stackcacherefill(c *mcache, order uint8) {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	if stackDebug &gt;= 1 {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		print(&#34;stackcacherefill order=&#34;, order, &#34;\n&#34;)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">// Grab some stacks from the global cache.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	<span class="comment">// Grab half of the allowed capacity (to prevent thrashing).</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	var list gclinkptr
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	var size uintptr
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	lock(&amp;stackpool[order].item.mu)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	for size &lt; _StackCacheSize/2 {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		x := stackpoolalloc(order)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		x.ptr().next = list
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		list = x
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		size += fixedStack &lt;&lt; order
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	unlock(&amp;stackpool[order].item.mu)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	c.stackcache[order].list = list
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	c.stackcache[order].size = size
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>func stackcacherelease(c *mcache, order uint8) {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	if stackDebug &gt;= 1 {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		print(&#34;stackcacherelease order=&#34;, order, &#34;\n&#34;)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	x := c.stackcache[order].list
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	size := c.stackcache[order].size
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	lock(&amp;stackpool[order].item.mu)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	for size &gt; _StackCacheSize/2 {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		y := x.ptr().next
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		stackpoolfree(x, order)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		x = y
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		size -= fixedStack &lt;&lt; order
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	unlock(&amp;stackpool[order].item.mu)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	c.stackcache[order].list = x
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	c.stackcache[order].size = size
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>func stackcache_clear(c *mcache) {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	if stackDebug &gt;= 1 {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		print(&#34;stackcache clear\n&#34;)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	for order := uint8(0); order &lt; _NumStackOrders; order++ {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		lock(&amp;stackpool[order].item.mu)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		x := c.stackcache[order].list
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		for x.ptr() != nil {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			y := x.ptr().next
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			stackpoolfree(x, order)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>			x = y
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		c.stackcache[order].list = 0
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		c.stackcache[order].size = 0
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		unlock(&amp;stackpool[order].item.mu)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">// stackalloc allocates an n byte stack.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// stackalloc must run on the system stack because it uses per-P</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">// resources and must not split the stack.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func stackalloc(n uint32) stack {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// Stackalloc must be called on scheduler stack, so that we</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// never try to grow the stack during the code that stackalloc runs.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// Doing so would cause a deadlock (issue 1547).</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	thisg := getg()
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	if thisg != thisg.m.g0 {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		throw(&#34;stackalloc not on scheduler stack&#34;)
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	if n&amp;(n-1) != 0 {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		throw(&#34;stack size not a power of 2&#34;)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	if stackDebug &gt;= 1 {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		print(&#34;stackalloc &#34;, n, &#34;\n&#34;)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	if debug.efence != 0 || stackFromSystem != 0 {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		n = uint32(alignUp(uintptr(n), physPageSize))
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		v := sysAlloc(uintptr(n), &amp;memstats.stacks_sys)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		if v == nil {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			throw(&#34;out of memory (stackalloc)&#34;)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		return stack{uintptr(v), uintptr(v) + uintptr(n)}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	<span class="comment">// Small stacks are allocated with a fixed-size free-list allocator.</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	<span class="comment">// If we need a stack of a bigger size, we fall back on allocating</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	<span class="comment">// a dedicated span.</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	var v unsafe.Pointer
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	if n &lt; fixedStack&lt;&lt;_NumStackOrders &amp;&amp; n &lt; _StackCacheSize {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		order := uint8(0)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		n2 := n
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		for n2 &gt; fixedStack {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			order++
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			n2 &gt;&gt;= 1
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		var x gclinkptr
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		if stackNoCache != 0 || thisg.m.p == 0 || thisg.m.preemptoff != &#34;&#34; {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			<span class="comment">// thisg.m.p == 0 can happen in the guts of exitsyscall</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			<span class="comment">// or procresize. Just get a stack from the global pool.</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			<span class="comment">// Also don&#39;t touch stackcache during gc</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			<span class="comment">// as it&#39;s flushed concurrently.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			lock(&amp;stackpool[order].item.mu)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			x = stackpoolalloc(order)
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			unlock(&amp;stackpool[order].item.mu)
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		} else {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			c := thisg.m.p.ptr().mcache
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			x = c.stackcache[order].list
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			if x.ptr() == nil {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>				stackcacherefill(c, order)
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>				x = c.stackcache[order].list
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>			c.stackcache[order].list = x.ptr().next
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			c.stackcache[order].size -= uintptr(n)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		v = unsafe.Pointer(x)
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	} else {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		var s *mspan
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		npage := uintptr(n) &gt;&gt; _PageShift
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		log2npage := stacklog2(npage)
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		<span class="comment">// Try to get a stack from the large stack cache.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		lock(&amp;stackLarge.lock)
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		if !stackLarge.free[log2npage].isEmpty() {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			s = stackLarge.free[log2npage].first
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			stackLarge.free[log2npage].remove(s)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		unlock(&amp;stackLarge.lock)
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		lockWithRankMayAcquire(&amp;mheap_.lock, lockRankMheap)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		if s == nil {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			<span class="comment">// Allocate a new stack from the heap.</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			s = mheap_.allocManual(npage, spanAllocStack)
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			if s == nil {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>				throw(&#34;out of memory&#34;)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>			}
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			osStackAlloc(s)
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			s.elemsize = uintptr(n)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		v = unsafe.Pointer(s.base())
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	if raceenabled {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		racemalloc(v, uintptr(n))
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	if msanenabled {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		msanmalloc(v, uintptr(n))
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	if asanenabled {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		asanunpoison(v, uintptr(n))
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	if stackDebug &gt;= 1 {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		print(&#34;  allocated &#34;, v, &#34;\n&#34;)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	return stack{uintptr(v), uintptr(v) + uintptr(n)}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span><span class="comment">// stackfree frees an n byte stack allocation at stk.</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span><span class="comment">// stackfree must run on the system stack because it uses per-P</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span><span class="comment">// resources and must not split the stack.</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>func stackfree(stk stack) {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	gp := getg()
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	v := unsafe.Pointer(stk.lo)
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	n := stk.hi - stk.lo
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	if n&amp;(n-1) != 0 {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		throw(&#34;stack not a power of 2&#34;)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	if stk.lo+n &lt; stk.hi {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		throw(&#34;bad stack size&#34;)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	if stackDebug &gt;= 1 {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		println(&#34;stackfree&#34;, v, n)
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		memclrNoHeapPointers(v, n) <span class="comment">// for testing, clobber stack data</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	if debug.efence != 0 || stackFromSystem != 0 {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		if debug.efence != 0 || stackFaultOnFree != 0 {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			sysFault(v, n)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		} else {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			sysFree(v, n, &amp;memstats.stacks_sys)
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		return
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	}
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	if msanenabled {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		msanfree(v, n)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	if asanenabled {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		asanpoison(v, n)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	if n &lt; fixedStack&lt;&lt;_NumStackOrders &amp;&amp; n &lt; _StackCacheSize {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		order := uint8(0)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		n2 := n
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		for n2 &gt; fixedStack {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			order++
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			n2 &gt;&gt;= 1
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		x := gclinkptr(v)
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		if stackNoCache != 0 || gp.m.p == 0 || gp.m.preemptoff != &#34;&#34; {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			lock(&amp;stackpool[order].item.mu)
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			stackpoolfree(x, order)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			unlock(&amp;stackpool[order].item.mu)
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		} else {
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			c := gp.m.p.ptr().mcache
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			if c.stackcache[order].size &gt;= _StackCacheSize {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>				stackcacherelease(c, order)
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			}
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			x.ptr().next = c.stackcache[order].list
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			c.stackcache[order].list = x
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			c.stackcache[order].size += n
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	} else {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		s := spanOfUnchecked(uintptr(v))
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		if s.state.get() != mSpanManual {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>			println(hex(s.base()), v)
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>			throw(&#34;bad span state&#34;)
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		if gcphase == _GCoff {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			<span class="comment">// Free the stack immediately if we&#39;re</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>			<span class="comment">// sweeping.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>			osStackFree(s)
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			mheap_.freeManual(s, spanAllocStack)
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		} else {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>			<span class="comment">// If the GC is running, we can&#39;t return a</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>			<span class="comment">// stack span to the heap because it could be</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			<span class="comment">// reused as a heap span, and this state</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			<span class="comment">// change would race with GC. Add it to the</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>			<span class="comment">// large stack cache instead.</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>			log2npage := stacklog2(s.npages)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			lock(&amp;stackLarge.lock)
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>			stackLarge.free[log2npage].insert(s)
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			unlock(&amp;stackLarge.lock)
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>}
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>var maxstacksize uintptr = 1 &lt;&lt; 20 <span class="comment">// enough until runtime.main sets it for real</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>var maxstackceiling = maxstacksize
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>var ptrnames = []string{
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	0: &#34;scalar&#34;,
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	1: &#34;ptr&#34;,
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span><span class="comment">// Stack frame layout</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">// (x86)</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span><span class="comment">// +------------------+</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span><span class="comment">// | args from caller |</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span><span class="comment">// +------------------+ &lt;- frame-&gt;argp</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span><span class="comment">// |  return address  |</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span><span class="comment">// +------------------+</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span><span class="comment">// |  caller&#39;s BP (*) | (*) if framepointer_enabled &amp;&amp; varp &gt; sp</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span><span class="comment">// +------------------+ &lt;- frame-&gt;varp</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span><span class="comment">// |     locals       |</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span><span class="comment">// +------------------+</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span><span class="comment">// |  args to callee  |</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span><span class="comment">// +------------------+ &lt;- frame-&gt;sp</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span><span class="comment">// (arm)</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">// +------------------+</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">// | args from caller |</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">// +------------------+ &lt;- frame-&gt;argp</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">// | caller&#39;s retaddr |</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">// +------------------+</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">// |  caller&#39;s FP (*) | (*) on ARM64, if framepointer_enabled &amp;&amp; varp &gt; sp</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span><span class="comment">// +------------------+ &lt;- frame-&gt;varp</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span><span class="comment">// |     locals       |</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span><span class="comment">// +------------------+</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span><span class="comment">// |  args to callee  |</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span><span class="comment">// +------------------+</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span><span class="comment">// |  return address  |</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span><span class="comment">// +------------------+ &lt;- frame-&gt;sp</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// varp &gt; sp means that the function has a frame;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">// varp == sp means frameless function.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>type adjustinfo struct {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	old   stack
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	delta uintptr <span class="comment">// ptr distance from old to new stack (newbase - oldbase)</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	<span class="comment">// sghi is the highest sudog.elem on the stack.</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	sghi uintptr
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span><span class="comment">// adjustpointer checks whether *vpp is in the old stack described by adjinfo.</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span><span class="comment">// If so, it rewrites *vpp to point into the new stack.</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>func adjustpointer(adjinfo *adjustinfo, vpp unsafe.Pointer) {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	pp := (*uintptr)(vpp)
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	p := *pp
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	if stackDebug &gt;= 4 {
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		print(&#34;        &#34;, pp, &#34;:&#34;, hex(p), &#34;\n&#34;)
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	if adjinfo.old.lo &lt;= p &amp;&amp; p &lt; adjinfo.old.hi {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		*pp = p + adjinfo.delta
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		if stackDebug &gt;= 3 {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>			print(&#34;        adjust ptr &#34;, pp, &#34;:&#34;, hex(p), &#34; -&gt; &#34;, hex(*pp), &#34;\n&#34;)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		}
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	}
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>}
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span><span class="comment">// Information from the compiler about the layout of stack frames.</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span><span class="comment">// Note: this type must agree with reflect.bitVector.</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>type bitvector struct {
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	n        int32 <span class="comment">// # of bits</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	bytedata *uint8
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span><span class="comment">// ptrbit returns the i&#39;th bit in bv.</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span><span class="comment">// ptrbit is less efficient than iterating directly over bitvector bits,</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span><span class="comment">// and should only be used in non-performance-critical code.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span><span class="comment">// See adjustpointers for an example of a high-efficiency walk of a bitvector.</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>func (bv *bitvector) ptrbit(i uintptr) uint8 {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	b := *(addb(bv.bytedata, i/8))
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	return (b &gt;&gt; (i % 8)) &amp; 1
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span><span class="comment">// bv describes the memory starting at address scanp.</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span><span class="comment">// Adjust any pointers contained therein.</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>func adjustpointers(scanp unsafe.Pointer, bv *bitvector, adjinfo *adjustinfo, f funcInfo) {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	minp := adjinfo.old.lo
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	maxp := adjinfo.old.hi
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	delta := adjinfo.delta
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	num := uintptr(bv.n)
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	<span class="comment">// If this frame might contain channel receive slots, use CAS</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	<span class="comment">// to adjust pointers. If the slot hasn&#39;t been received into</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	<span class="comment">// yet, it may contain stack pointers and a concurrent send</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	<span class="comment">// could race with adjusting those pointers. (The sent value</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	<span class="comment">// itself can never contain stack pointers.)</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	useCAS := uintptr(scanp) &lt; adjinfo.sghi
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; num; i += 8 {
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		if stackDebug &gt;= 4 {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>			for j := uintptr(0); j &lt; 8; j++ {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>				print(&#34;        &#34;, add(scanp, (i+j)*goarch.PtrSize), &#34;:&#34;, ptrnames[bv.ptrbit(i+j)], &#34;:&#34;, hex(*(*uintptr)(add(scanp, (i+j)*goarch.PtrSize))), &#34; # &#34;, i, &#34; &#34;, *addb(bv.bytedata, i/8), &#34;\n&#34;)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>			}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		}
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		b := *(addb(bv.bytedata, i/8))
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		for b != 0 {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>			j := uintptr(sys.TrailingZeros8(b))
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>			b &amp;= b - 1
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>			pp := (*uintptr)(add(scanp, (i+j)*goarch.PtrSize))
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		retry:
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			p := *pp
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			if f.valid() &amp;&amp; 0 &lt; p &amp;&amp; p &lt; minLegalPointer &amp;&amp; debug.invalidptr != 0 {
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>				<span class="comment">// Looks like a junk value in a pointer slot.</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>				<span class="comment">// Live analysis wrong?</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>				getg().m.traceback = 2
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>				print(&#34;runtime: bad pointer in frame &#34;, funcname(f), &#34; at &#34;, pp, &#34;: &#34;, hex(p), &#34;\n&#34;)
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>				throw(&#34;invalid pointer found on stack&#34;)
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			if minp &lt;= p &amp;&amp; p &lt; maxp {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>				if stackDebug &gt;= 3 {
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>					print(&#34;adjust ptr &#34;, hex(p), &#34; &#34;, funcname(f), &#34;\n&#34;)
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>				}
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>				if useCAS {
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>					ppu := (*unsafe.Pointer)(unsafe.Pointer(pp))
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>					if !atomic.Casp1(ppu, unsafe.Pointer(p), unsafe.Pointer(p+delta)) {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>						goto retry
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>					}
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>				} else {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>					*pp = p + delta
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>				}
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			}
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span><span class="comment">// Note: the argument/return area is adjusted by the callee.</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>func adjustframe(frame *stkframe, adjinfo *adjustinfo) {
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	if frame.continpc == 0 {
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		<span class="comment">// Frame is dead.</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		return
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	}
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	f := frame.fn
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	if stackDebug &gt;= 2 {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>		print(&#34;    adjusting &#34;, funcname(f), &#34; frame=[&#34;, hex(frame.sp), &#34;,&#34;, hex(frame.fp), &#34;] pc=&#34;, hex(frame.pc), &#34; continpc=&#34;, hex(frame.continpc), &#34;\n&#34;)
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	<span class="comment">// Adjust saved frame pointer if there is one.</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	if (goarch.ArchFamily == goarch.AMD64 || goarch.ArchFamily == goarch.ARM64) &amp;&amp; frame.argp-frame.varp == 2*goarch.PtrSize {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		if stackDebug &gt;= 3 {
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>			print(&#34;      saved bp\n&#34;)
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		}
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		if debugCheckBP {
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>			<span class="comment">// Frame pointers should always point to the next higher frame on</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>			<span class="comment">// the Go stack (or be nil, for the top frame on the stack).</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>			bp := *(*uintptr)(unsafe.Pointer(frame.varp))
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>			if bp != 0 &amp;&amp; (bp &lt; adjinfo.old.lo || bp &gt;= adjinfo.old.hi) {
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>				println(&#34;runtime: found invalid frame pointer&#34;)
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>				print(&#34;bp=&#34;, hex(bp), &#34; min=&#34;, hex(adjinfo.old.lo), &#34; max=&#34;, hex(adjinfo.old.hi), &#34;\n&#34;)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>				throw(&#34;bad frame pointer&#34;)
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>			}
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>		<span class="comment">// On AMD64, this is the caller&#39;s frame pointer saved in the current</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		<span class="comment">// frame.</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		<span class="comment">// On ARM64, this is the frame pointer of the caller&#39;s caller saved</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		<span class="comment">// by the caller in its frame (one word below its SP).</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>		adjustpointer(adjinfo, unsafe.Pointer(frame.varp))
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	}
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	locals, args, objs := frame.getStackMap(true)
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	<span class="comment">// Adjust local variables if stack frame has been allocated.</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	if locals.n &gt; 0 {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		size := uintptr(locals.n) * goarch.PtrSize
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		adjustpointers(unsafe.Pointer(frame.varp-size), &amp;locals, adjinfo, f)
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	<span class="comment">// Adjust arguments.</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	if args.n &gt; 0 {
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		if stackDebug &gt;= 3 {
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>			print(&#34;      args\n&#34;)
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		}
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		adjustpointers(unsafe.Pointer(frame.argp), &amp;args, adjinfo, funcInfo{})
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	}
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	<span class="comment">// Adjust pointers in all stack objects (whether they are live or not).</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	<span class="comment">// See comments in mgcmark.go:scanframeworker.</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	if frame.varp != 0 {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		for i := range objs {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			obj := &amp;objs[i]
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			off := obj.off
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>			base := frame.varp <span class="comment">// locals base pointer</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>			if off &gt;= 0 {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>				base = frame.argp <span class="comment">// arguments and return values base pointer</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>			}
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>			p := base + uintptr(off)
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			if p &lt; frame.sp {
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>				<span class="comment">// Object hasn&#39;t been allocated in the frame yet.</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>				<span class="comment">// (Happens when the stack bounds check fails and</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>				<span class="comment">// we call into morestack.)</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>				continue
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>			}
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			ptrdata := obj.ptrdata()
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>			gcdata := obj.gcdata()
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>			var s *mspan
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>			if obj.useGCProg() {
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>				<span class="comment">// See comments in mgcmark.go:scanstack</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>				s = materializeGCProg(ptrdata, gcdata)
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>				gcdata = (*byte)(unsafe.Pointer(s.startAddr))
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>			}
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>			for i := uintptr(0); i &lt; ptrdata; i += goarch.PtrSize {
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>				if *addb(gcdata, i/(8*goarch.PtrSize))&gt;&gt;(i/goarch.PtrSize&amp;7)&amp;1 != 0 {
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>					adjustpointer(adjinfo, unsafe.Pointer(p+i))
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>				}
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>			}
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>			if s != nil {
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>				dematerializeGCProg(s)
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>			}
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		}
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>}
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>func adjustctxt(gp *g, adjinfo *adjustinfo) {
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	adjustpointer(adjinfo, unsafe.Pointer(&amp;gp.sched.ctxt))
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	if !framepointer_enabled {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		return
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	}
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	if debugCheckBP {
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>		bp := gp.sched.bp
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		if bp != 0 &amp;&amp; (bp &lt; adjinfo.old.lo || bp &gt;= adjinfo.old.hi) {
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>			println(&#34;runtime: found invalid top frame pointer&#34;)
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>			print(&#34;bp=&#34;, hex(bp), &#34; min=&#34;, hex(adjinfo.old.lo), &#34; max=&#34;, hex(adjinfo.old.hi), &#34;\n&#34;)
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>			throw(&#34;bad top frame pointer&#34;)
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		}
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	}
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	oldfp := gp.sched.bp
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	adjustpointer(adjinfo, unsafe.Pointer(&amp;gp.sched.bp))
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	if GOARCH == &#34;arm64&#34; {
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		<span class="comment">// On ARM64, the frame pointer is saved one word *below* the SP,</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		<span class="comment">// which is not copied or adjusted in any frame. Do it explicitly</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		<span class="comment">// here.</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		if oldfp == gp.sched.sp-goarch.PtrSize {
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>			memmove(unsafe.Pointer(gp.sched.bp), unsafe.Pointer(oldfp), goarch.PtrSize)
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>			adjustpointer(adjinfo, unsafe.Pointer(gp.sched.bp))
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>		}
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	}
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>}
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>func adjustdefers(gp *g, adjinfo *adjustinfo) {
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	<span class="comment">// Adjust pointers in the Defer structs.</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	<span class="comment">// We need to do this first because we need to adjust the</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	<span class="comment">// defer.link fields so we always work on the new stack.</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	adjustpointer(adjinfo, unsafe.Pointer(&amp;gp._defer))
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	for d := gp._defer; d != nil; d = d.link {
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		adjustpointer(adjinfo, unsafe.Pointer(&amp;d.fn))
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>		adjustpointer(adjinfo, unsafe.Pointer(&amp;d.sp))
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		adjustpointer(adjinfo, unsafe.Pointer(&amp;d.link))
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	}
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>}
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>func adjustpanics(gp *g, adjinfo *adjustinfo) {
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	<span class="comment">// Panics are on stack and already adjusted.</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	<span class="comment">// Update pointer to head of list in G.</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	adjustpointer(adjinfo, unsafe.Pointer(&amp;gp._panic))
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>}
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>func adjustsudogs(gp *g, adjinfo *adjustinfo) {
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	<span class="comment">// the data elements pointed to by a SudoG structure</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>	<span class="comment">// might be in the stack.</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>	for s := gp.waiting; s != nil; s = s.waitlink {
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		adjustpointer(adjinfo, unsafe.Pointer(&amp;s.elem))
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	}
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>}
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>func fillstack(stk stack, b byte) {
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	for p := stk.lo; p &lt; stk.hi; p++ {
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>		*(*byte)(unsafe.Pointer(p)) = b
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	}
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>}
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>func findsghi(gp *g, stk stack) uintptr {
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	var sghi uintptr
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	for sg := gp.waiting; sg != nil; sg = sg.waitlink {
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		p := uintptr(sg.elem) + uintptr(sg.c.elemsize)
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>		if stk.lo &lt;= p &amp;&amp; p &lt; stk.hi &amp;&amp; p &gt; sghi {
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>			sghi = p
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>		}
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	}
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	return sghi
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>}
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span><span class="comment">// syncadjustsudogs adjusts gp&#39;s sudogs and copies the part of gp&#39;s</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span><span class="comment">// stack they refer to while synchronizing with concurrent channel</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span><span class="comment">// operations. It returns the number of bytes of stack copied.</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>func syncadjustsudogs(gp *g, used uintptr, adjinfo *adjustinfo) uintptr {
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	if gp.waiting == nil {
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>		return 0
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	}
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	<span class="comment">// Lock channels to prevent concurrent send/receive.</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	var lastc *hchan
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	for sg := gp.waiting; sg != nil; sg = sg.waitlink {
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		if sg.c != lastc {
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>			<span class="comment">// There is a ranking cycle here between gscan bit and</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>			<span class="comment">// hchan locks. Normally, we only allow acquiring hchan</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>			<span class="comment">// locks and then getting a gscan bit. In this case, we</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>			<span class="comment">// already have the gscan bit. We allow acquiring hchan</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>			<span class="comment">// locks here as a special case, since a deadlock can&#39;t</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>			<span class="comment">// happen because the G involved must already be</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>			<span class="comment">// suspended. So, we get a special hchan lock rank here</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>			<span class="comment">// that is lower than gscan, but doesn&#39;t allow acquiring</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>			<span class="comment">// any other locks other than hchan.</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>			lockWithRank(&amp;sg.c.lock, lockRankHchanLeaf)
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>		}
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		lastc = sg.c
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	}
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	<span class="comment">// Adjust sudogs.</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	adjustsudogs(gp, adjinfo)
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	<span class="comment">// Copy the part of the stack the sudogs point in to</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	<span class="comment">// while holding the lock to prevent races on</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	<span class="comment">// send/receive slots.</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	var sgsize uintptr
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	if adjinfo.sghi != 0 {
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		oldBot := adjinfo.old.hi - used
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		newBot := oldBot + adjinfo.delta
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		sgsize = adjinfo.sghi - oldBot
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		memmove(unsafe.Pointer(newBot), unsafe.Pointer(oldBot), sgsize)
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	}
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	<span class="comment">// Unlock channels.</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	lastc = nil
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	for sg := gp.waiting; sg != nil; sg = sg.waitlink {
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>		if sg.c != lastc {
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>			unlock(&amp;sg.c.lock)
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>		}
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		lastc = sg.c
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	}
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	return sgsize
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span><span class="comment">// Copies gp&#39;s stack to a new stack of a different size.</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span><span class="comment">// Caller must have changed gp status to Gcopystack.</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>func copystack(gp *g, newsize uintptr) {
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	if gp.syscallsp != 0 {
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		throw(&#34;stack growth not allowed in system call&#34;)
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	}
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	old := gp.stack
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	if old.lo == 0 {
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>		throw(&#34;nil stackbase&#34;)
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	}
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	used := old.hi - gp.sched.sp
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	<span class="comment">// Add just the difference to gcController.addScannableStack.</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	<span class="comment">// g0 stacks never move, so this will never account for them.</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s also fine if we have no P, addScannableStack can deal with</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	<span class="comment">// that case.</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	gcController.addScannableStack(getg().m.p.ptr(), int64(newsize)-int64(old.hi-old.lo))
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	<span class="comment">// allocate new stack</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	new := stackalloc(uint32(newsize))
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	if stackPoisonCopy != 0 {
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		fillstack(new, 0xfd)
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	}
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	if stackDebug &gt;= 1 {
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>		print(&#34;copystack gp=&#34;, gp, &#34; [&#34;, hex(old.lo), &#34; &#34;, hex(old.hi-used), &#34; &#34;, hex(old.hi), &#34;]&#34;, &#34; -&gt; [&#34;, hex(new.lo), &#34; &#34;, hex(new.hi-used), &#34; &#34;, hex(new.hi), &#34;]/&#34;, newsize, &#34;\n&#34;)
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	}
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	<span class="comment">// Compute adjustment.</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	var adjinfo adjustinfo
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	adjinfo.old = old
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	adjinfo.delta = new.hi - old.hi
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	<span class="comment">// Adjust sudogs, synchronizing with channel ops if necessary.</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	ncopy := used
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	if !gp.activeStackChans {
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>		if newsize &lt; old.hi-old.lo &amp;&amp; gp.parkingOnChan.Load() {
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>			<span class="comment">// It&#39;s not safe for someone to shrink this stack while we&#39;re actively</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>			<span class="comment">// parking on a channel, but it is safe to grow since we do that</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			<span class="comment">// ourselves and explicitly don&#39;t want to synchronize with channels</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>			<span class="comment">// since we could self-deadlock.</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			throw(&#34;racy sudog adjustment due to parking on channel&#34;)
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>		}
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>		adjustsudogs(gp, &amp;adjinfo)
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	} else {
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		<span class="comment">// sudogs may be pointing in to the stack and gp has</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>		<span class="comment">// released channel locks, so other goroutines could</span>
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>		<span class="comment">// be writing to gp&#39;s stack. Find the highest such</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		<span class="comment">// pointer so we can handle everything there and below</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>		<span class="comment">// carefully. (This shouldn&#39;t be far from the bottom</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		<span class="comment">// of the stack, so there&#39;s little cost in handling</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		<span class="comment">// everything below it carefully.)</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>		adjinfo.sghi = findsghi(gp, old)
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>		<span class="comment">// Synchronize with channel ops and copy the part of</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>		<span class="comment">// the stack they may interact with.</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>		ncopy -= syncadjustsudogs(gp, used, &amp;adjinfo)
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	}
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	<span class="comment">// Copy the stack (or the rest of it) to the new location</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	memmove(unsafe.Pointer(new.hi-ncopy), unsafe.Pointer(old.hi-ncopy), ncopy)
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	<span class="comment">// Adjust remaining structures that have pointers into stacks.</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	<span class="comment">// We have to do most of these before we traceback the new</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	<span class="comment">// stack because gentraceback uses them.</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	adjustctxt(gp, &amp;adjinfo)
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>	adjustdefers(gp, &amp;adjinfo)
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>	adjustpanics(gp, &amp;adjinfo)
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>	if adjinfo.sghi != 0 {
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>		adjinfo.sghi += adjinfo.delta
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>	}
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	<span class="comment">// Swap out old stack for new one</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	gp.stack = new
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	gp.stackguard0 = new.lo + stackGuard <span class="comment">// NOTE: might clobber a preempt request</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	gp.sched.sp = new.hi - used
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	gp.stktopsp += adjinfo.delta
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	<span class="comment">// Adjust pointers in the new stack.</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>	var u unwinder
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>	for u.init(gp, 0); u.valid(); u.next() {
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>		adjustframe(&amp;u.frame, &amp;adjinfo)
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>	}
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	<span class="comment">// free old stack</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>	if stackPoisonCopy != 0 {
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>		fillstack(old, 0xfc)
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>	}
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	stackfree(old)
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>}
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>
<span id="L941" class="ln">   941&nbsp;&nbsp;</span><span class="comment">// round x up to a power of 2.</span>
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>func round2(x int32) int32 {
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>	s := uint(0)
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>	for 1&lt;&lt;s &lt; x {
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>		s++
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>	}
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>	return 1 &lt;&lt; s
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>}
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>
<span id="L950" class="ln">   950&nbsp;&nbsp;</span><span class="comment">// Called from runtime·morestack when more stack is needed.</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span><span class="comment">// Allocate larger stack and relocate to new stack.</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span><span class="comment">// Stack growth is multiplicative, for constant amortized cost.</span>
<span id="L953" class="ln">   953&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L954" class="ln">   954&nbsp;&nbsp;</span><span class="comment">// g-&gt;atomicstatus will be Grunning or Gscanrunning upon entry.</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span><span class="comment">// If the scheduler is trying to stop this g, then it will set preemptStop.</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span><span class="comment">// This must be nowritebarrierrec because it can be called as part of</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span><span class="comment">// stack growth from other nowritebarrierrec functions, but the</span>
<span id="L959" class="ln">   959&nbsp;&nbsp;</span><span class="comment">// compiler doesn&#39;t check this.</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>func newstack() {
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>	thisg := getg()
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	<span class="comment">// TODO: double check all gp. shouldn&#39;t be getg().</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	if thisg.m.morebuf.g.ptr().stackguard0 == stackFork {
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		throw(&#34;stack growth after fork&#34;)
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>	}
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>	if thisg.m.morebuf.g.ptr() != thisg.m.curg {
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>		print(&#34;runtime: newstack called from g=&#34;, hex(thisg.m.morebuf.g), &#34;\n&#34;+&#34;\tm=&#34;, thisg.m, &#34; m-&gt;curg=&#34;, thisg.m.curg, &#34; m-&gt;g0=&#34;, thisg.m.g0, &#34; m-&gt;gsignal=&#34;, thisg.m.gsignal, &#34;\n&#34;)
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>		morebuf := thisg.m.morebuf
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>		traceback(morebuf.pc, morebuf.sp, morebuf.lr, morebuf.g.ptr())
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>		throw(&#34;runtime: wrong goroutine in newstack&#34;)
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	}
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>	gp := thisg.m.curg
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>	if thisg.m.curg.throwsplit {
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>		<span class="comment">// Update syscallsp, syscallpc in case traceback uses them.</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>		morebuf := thisg.m.morebuf
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>		gp.syscallsp = morebuf.sp
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>		gp.syscallpc = morebuf.pc
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		pcname, pcoff := &#34;(unknown)&#34;, uintptr(0)
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		f := findfunc(gp.sched.pc)
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>		if f.valid() {
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>			pcname = funcname(f)
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>			pcoff = gp.sched.pc - f.entry()
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>		}
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>		print(&#34;runtime: newstack at &#34;, pcname, &#34;+&#34;, hex(pcoff),
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>			&#34; sp=&#34;, hex(gp.sched.sp), &#34; stack=[&#34;, hex(gp.stack.lo), &#34;, &#34;, hex(gp.stack.hi), &#34;]\n&#34;,
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>			&#34;\tmorebuf={pc:&#34;, hex(morebuf.pc), &#34; sp:&#34;, hex(morebuf.sp), &#34; lr:&#34;, hex(morebuf.lr), &#34;}\n&#34;,
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>			&#34;\tsched={pc:&#34;, hex(gp.sched.pc), &#34; sp:&#34;, hex(gp.sched.sp), &#34; lr:&#34;, hex(gp.sched.lr), &#34; ctxt:&#34;, gp.sched.ctxt, &#34;}\n&#34;)
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>		thisg.m.traceback = 2 <span class="comment">// Include runtime frames</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>		traceback(morebuf.pc, morebuf.sp, morebuf.lr, gp)
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		throw(&#34;runtime: stack split at bad time&#34;)
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	}
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	morebuf := thisg.m.morebuf
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	thisg.m.morebuf.pc = 0
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	thisg.m.morebuf.lr = 0
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>	thisg.m.morebuf.sp = 0
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>	thisg.m.morebuf.g = 0
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>	<span class="comment">// NOTE: stackguard0 may change underfoot, if another thread</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	<span class="comment">// is about to try to preempt gp. Read it just once and use that same</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	<span class="comment">// value now and below.</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	stackguard0 := atomic.Loaduintptr(&amp;gp.stackguard0)
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>	<span class="comment">// Be conservative about where we preempt.</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>	<span class="comment">// We are interested in preempting user Go code, not runtime code.</span>
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;re holding locks, mallocing, or preemption is disabled, don&#39;t</span>
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	<span class="comment">// preempt.</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	<span class="comment">// This check is very early in newstack so that even the status change</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	<span class="comment">// from Grunning to Gwaiting and back doesn&#39;t happen in this case.</span>
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	<span class="comment">// That status change by itself can be viewed as a small preemption,</span>
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	<span class="comment">// because the GC might change Gwaiting to Gscanwaiting, and then</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	<span class="comment">// this goroutine has to wait for the GC to finish before continuing.</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	<span class="comment">// If the GC is in some way dependent on this goroutine (for example,</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>	<span class="comment">// it needs a lock held by the goroutine), that small preemption turns</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>	<span class="comment">// into a real deadlock.</span>
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	preempt := stackguard0 == stackPreempt
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	if preempt {
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		if !canPreemptM(thisg.m) {
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>			<span class="comment">// Let the goroutine keep running for now.</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>			<span class="comment">// gp-&gt;preempt is set, so it will be preempted next time.</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>			gp.stackguard0 = gp.stack.lo + stackGuard
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>			gogo(&amp;gp.sched) <span class="comment">// never return</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>		}
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	}
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	if gp.stack.lo == 0 {
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>		throw(&#34;missing stack in newstack&#34;)
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>	}
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>	sp := gp.sched.sp
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	if goarch.ArchFamily == goarch.AMD64 || goarch.ArchFamily == goarch.I386 || goarch.ArchFamily == goarch.WASM {
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>		<span class="comment">// The call to morestack cost a word.</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>		sp -= goarch.PtrSize
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	}
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>	if stackDebug &gt;= 1 || sp &lt; gp.stack.lo {
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>		print(&#34;runtime: newstack sp=&#34;, hex(sp), &#34; stack=[&#34;, hex(gp.stack.lo), &#34;, &#34;, hex(gp.stack.hi), &#34;]\n&#34;,
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>			&#34;\tmorebuf={pc:&#34;, hex(morebuf.pc), &#34; sp:&#34;, hex(morebuf.sp), &#34; lr:&#34;, hex(morebuf.lr), &#34;}\n&#34;,
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>			&#34;\tsched={pc:&#34;, hex(gp.sched.pc), &#34; sp:&#34;, hex(gp.sched.sp), &#34; lr:&#34;, hex(gp.sched.lr), &#34; ctxt:&#34;, gp.sched.ctxt, &#34;}\n&#34;)
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>	}
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>	if sp &lt; gp.stack.lo {
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>		print(&#34;runtime: gp=&#34;, gp, &#34;, goid=&#34;, gp.goid, &#34;, gp-&gt;status=&#34;, hex(readgstatus(gp)), &#34;\n &#34;)
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>		print(&#34;runtime: split stack overflow: &#34;, hex(sp), &#34; &lt; &#34;, hex(gp.stack.lo), &#34;\n&#34;)
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>		throw(&#34;runtime: split stack overflow&#34;)
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>	}
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>	if preempt {
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>		if gp == thisg.m.g0 {
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>			throw(&#34;runtime: preempt g0&#34;)
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>		}
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>		if thisg.m.p == 0 &amp;&amp; thisg.m.locks == 0 {
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>			throw(&#34;runtime: g is running but p is not&#34;)
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>		}
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>		if gp.preemptShrink {
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>			<span class="comment">// We&#39;re at a synchronous safe point now, so</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>			<span class="comment">// do the pending stack shrink.</span>
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>			gp.preemptShrink = false
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>			shrinkstack(gp)
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>		}
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>		if gp.preemptStop {
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>			preemptPark(gp) <span class="comment">// never returns</span>
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>		}
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>		<span class="comment">// Act like goroutine called runtime.Gosched.</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>		gopreempt_m(gp) <span class="comment">// never return</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	}
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	<span class="comment">// Allocate a bigger segment and move the stack.</span>
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	oldsize := gp.stack.hi - gp.stack.lo
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	newsize := oldsize * 2
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	<span class="comment">// Make sure we grow at least as much as needed to fit the new frame.</span>
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>	<span class="comment">// (This is just an optimization - the caller of morestack will</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	<span class="comment">// recheck the bounds on return.)</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>	if f := findfunc(gp.sched.pc); f.valid() {
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>		max := uintptr(funcMaxSPDelta(f))
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>		needed := max + stackGuard
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>		used := gp.stack.hi - gp.sched.sp
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>		for newsize-used &lt; needed {
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>			newsize *= 2
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>		}
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	}
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>	if stackguard0 == stackForceMove {
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>		<span class="comment">// Forced stack movement used for debugging.</span>
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>		<span class="comment">// Don&#39;t double the stack (or we may quickly run out</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>		<span class="comment">// if this is done repeatedly).</span>
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>		newsize = oldsize
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>	}
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>	if newsize &gt; maxstacksize || newsize &gt; maxstackceiling {
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>		if maxstacksize &lt; maxstackceiling {
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>			print(&#34;runtime: goroutine stack exceeds &#34;, maxstacksize, &#34;-byte limit\n&#34;)
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>		} else {
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>			print(&#34;runtime: goroutine stack exceeds &#34;, maxstackceiling, &#34;-byte limit\n&#34;)
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>		}
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>		print(&#34;runtime: sp=&#34;, hex(sp), &#34; stack=[&#34;, hex(gp.stack.lo), &#34;, &#34;, hex(gp.stack.hi), &#34;]\n&#34;)
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>		throw(&#34;stack overflow&#34;)
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>	}
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>	<span class="comment">// The goroutine must be executing in order to call newstack,</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>	<span class="comment">// so it must be Grunning (or Gscanrunning).</span>
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>	casgstatus(gp, _Grunning, _Gcopystack)
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>	<span class="comment">// The concurrent GC will not scan the stack while we are doing the copy since</span>
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>	<span class="comment">// the gp is in a Gcopystack status.</span>
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>	copystack(gp, newsize)
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>	if stackDebug &gt;= 1 {
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>		print(&#34;stack grow done\n&#34;)
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>	}
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>	casgstatus(gp, _Gcopystack, _Grunning)
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>	gogo(&amp;gp.sched)
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>}
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>func nilfunc() {
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	*(*uint8)(nil) = 0
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>}
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span><span class="comment">// adjust Gobuf as if it executed a call to fn</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span><span class="comment">// and then stopped before the first instruction in fn.</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>func gostartcallfn(gobuf *gobuf, fv *funcval) {
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>	var fn unsafe.Pointer
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>	if fv != nil {
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>		fn = unsafe.Pointer(fv.fn)
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>	} else {
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>		fn = unsafe.Pointer(abi.FuncPCABIInternal(nilfunc))
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>	}
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>	gostartcall(gobuf, fn, unsafe.Pointer(fv))
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>}
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span><span class="comment">// isShrinkStackSafe returns whether it&#39;s safe to attempt to shrink</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span><span class="comment">// gp&#39;s stack. Shrinking the stack is only safe when we have precise</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span><span class="comment">// pointer maps for all frames on the stack.</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>func isShrinkStackSafe(gp *g) bool {
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	<span class="comment">// We can&#39;t copy the stack if we&#39;re in a syscall.</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	<span class="comment">// The syscall might have pointers into the stack and</span>
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	<span class="comment">// often we don&#39;t have precise pointer maps for the innermost</span>
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>	<span class="comment">// frames.</span>
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>	<span class="comment">// We also can&#39;t copy the stack if we&#39;re at an asynchronous</span>
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>	<span class="comment">// safe-point because we don&#39;t have precise pointer maps for</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>	<span class="comment">// all frames.</span>
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>	<span class="comment">// We also can&#39;t *shrink* the stack in the window between the</span>
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>	<span class="comment">// goroutine calling gopark to park on a channel and</span>
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>	<span class="comment">// gp.activeStackChans being set.</span>
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>	return gp.syscallsp == 0 &amp;&amp; !gp.asyncSafePoint &amp;&amp; !gp.parkingOnChan.Load()
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>}
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span><span class="comment">// Maybe shrink the stack being used by gp.</span>
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span><span class="comment">// gp must be stopped and we must own its stack. It may be in</span>
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span><span class="comment">// _Grunning, but only if this is our own user G.</span>
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>func shrinkstack(gp *g) {
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>	if gp.stack.lo == 0 {
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>		throw(&#34;missing stack in shrinkstack&#34;)
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>	}
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>	if s := readgstatus(gp); s&amp;_Gscan == 0 {
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>		<span class="comment">// We don&#39;t own the stack via _Gscan. We could still</span>
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>		<span class="comment">// own it if this is our own user G and we&#39;re on the</span>
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>		<span class="comment">// system stack.</span>
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>		if !(gp == getg().m.curg &amp;&amp; getg() != getg().m.curg &amp;&amp; s == _Grunning) {
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>			<span class="comment">// We don&#39;t own the stack.</span>
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>			throw(&#34;bad status in shrinkstack&#34;)
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>		}
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>	}
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>	if !isShrinkStackSafe(gp) {
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>		throw(&#34;shrinkstack at bad time&#34;)
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>	}
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>	<span class="comment">// Check for self-shrinks while in a libcall. These may have</span>
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>	<span class="comment">// pointers into the stack disguised as uintptrs, but these</span>
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>	<span class="comment">// code paths should all be nosplit.</span>
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>	if gp == getg().m.curg &amp;&amp; gp.m.libcallsp != 0 {
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>		throw(&#34;shrinking stack in libcall&#34;)
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>	}
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>	if debug.gcshrinkstackoff &gt; 0 {
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>		return
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>	}
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>	f := findfunc(gp.startpc)
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>	if f.valid() &amp;&amp; f.funcID == abi.FuncID_gcBgMarkWorker {
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re not allowed to shrink the gcBgMarkWorker</span>
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>		<span class="comment">// stack (see gcBgMarkWorker for explanation).</span>
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>		return
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>	}
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>	oldsize := gp.stack.hi - gp.stack.lo
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>	newsize := oldsize / 2
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t shrink the allocation below the minimum-sized stack</span>
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>	<span class="comment">// allocation.</span>
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>	if newsize &lt; fixedStack {
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>		return
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>	}
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>	<span class="comment">// Compute how much of the stack is currently in use and only</span>
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>	<span class="comment">// shrink the stack if gp is using less than a quarter of its</span>
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>	<span class="comment">// current stack. The currently used stack includes everything</span>
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>	<span class="comment">// down to the SP plus the stack guard space that ensures</span>
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>	<span class="comment">// there&#39;s room for nosplit functions.</span>
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>	avail := gp.stack.hi - gp.stack.lo
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>	if used := gp.stack.hi - gp.sched.sp + stackNosplit; used &gt;= avail/4 {
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>		return
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>	}
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>	if stackDebug &gt; 0 {
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>		print(&#34;shrinking stack &#34;, oldsize, &#34;-&gt;&#34;, newsize, &#34;\n&#34;)
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>	}
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>	copystack(gp, newsize)
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>}
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span><span class="comment">// freeStackSpans frees unused stack spans at the end of GC.</span>
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>func freeStackSpans() {
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>	<span class="comment">// Scan stack pools for empty stack spans.</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>	for order := range stackpool {
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>		lock(&amp;stackpool[order].item.mu)
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>		list := &amp;stackpool[order].item.span
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>		for s := list.first; s != nil; {
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>			next := s.next
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>			if s.allocCount == 0 {
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>				list.remove(s)
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>				s.manualFreeList = 0
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>				osStackFree(s)
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>				mheap_.freeManual(s, spanAllocStack)
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>			}
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>			s = next
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>		}
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>		unlock(&amp;stackpool[order].item.mu)
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>	}
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>	<span class="comment">// Free large stack spans.</span>
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>	lock(&amp;stackLarge.lock)
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>	for i := range stackLarge.free {
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>		for s := stackLarge.free[i].first; s != nil; {
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>			next := s.next
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>			stackLarge.free[i].remove(s)
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>			osStackFree(s)
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>			mheap_.freeManual(s, spanAllocStack)
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>			s = next
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>		}
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>	}
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>	unlock(&amp;stackLarge.lock)
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>}
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span><span class="comment">// A stackObjectRecord is generated by the compiler for each stack object in a stack frame.</span>
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span><span class="comment">// This record must match the generator code in cmd/compile/internal/liveness/plive.go:emitStackObjects.</span>
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>type stackObjectRecord struct {
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>	<span class="comment">// offset in frame</span>
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>	<span class="comment">// if negative, offset from varp</span>
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>	<span class="comment">// if non-negative, offset from argp</span>
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>	off       int32
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>	size      int32
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>	_ptrdata  int32  <span class="comment">// ptrdata, or -ptrdata is GC prog is used</span>
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>	gcdataoff uint32 <span class="comment">// offset to gcdata from moduledata.rodata</span>
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>}
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>func (r *stackObjectRecord) useGCProg() bool {
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>	return r._ptrdata &lt; 0
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>}
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>func (r *stackObjectRecord) ptrdata() uintptr {
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>	x := r._ptrdata
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>		return uintptr(-x)
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	}
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	return uintptr(x)
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>}
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span><span class="comment">// gcdata returns pointer map or GC prog of the type.</span>
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>func (r *stackObjectRecord) gcdata() *byte {
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>	ptr := uintptr(unsafe.Pointer(r))
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>	var mod *moduledata
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>	for datap := &amp;firstmoduledata; datap != nil; datap = datap.next {
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>		if datap.gofunc &lt;= ptr &amp;&amp; ptr &lt; datap.end {
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>			mod = datap
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>			break
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>		}
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>	}
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>	<span class="comment">// If you get a panic here due to a nil mod,</span>
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>	<span class="comment">// you may have made a copy of a stackObjectRecord.</span>
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>	<span class="comment">// You must use the original pointer.</span>
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>	res := mod.rodata + uintptr(r.gcdataoff)
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>	return (*byte)(unsafe.Pointer(res))
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>}
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span><span class="comment">// This is exported as ABI0 via linkname so obj can call it.</span>
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span><span class="comment">//go:linkname morestackc</span>
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>func morestackc() {
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>	throw(&#34;attempt to execute system stack code on user stack&#34;)
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>}
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span><span class="comment">// startingStackSize is the amount of stack that new goroutines start with.</span>
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span><span class="comment">// It is a power of 2, and between _FixedStack and maxstacksize, inclusive.</span>
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span><span class="comment">// startingStackSize is updated every GC by tracking the average size of</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span><span class="comment">// stacks scanned during the GC.</span>
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>var startingStackSize uint32 = fixedStack
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>func gcComputeStartingStackSize() {
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>	if debug.adaptivestackstart == 0 {
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>		return
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>	}
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>	<span class="comment">// For details, see the design doc at</span>
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>	<span class="comment">// https://docs.google.com/document/d/1YDlGIdVTPnmUiTAavlZxBI1d9pwGQgZT7IKFKlIXohQ/edit?usp=sharing</span>
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>	<span class="comment">// The basic algorithm is to track the average size of stacks</span>
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>	<span class="comment">// and start goroutines with stack equal to that average size.</span>
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>	<span class="comment">// Starting at the average size uses at most 2x the space that</span>
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>	<span class="comment">// an ideal algorithm would have used.</span>
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>	<span class="comment">// This is just a heuristic to avoid excessive stack growth work</span>
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>	<span class="comment">// early in a goroutine&#39;s lifetime. See issue 18138. Stacks that</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>	<span class="comment">// are allocated too small can still grow, and stacks allocated</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>	<span class="comment">// too large can still shrink.</span>
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>	var scannedStackSize uint64
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>	var scannedStacks uint64
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>	for _, p := range allp {
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>		scannedStackSize += p.scannedStackSize
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>		scannedStacks += p.scannedStacks
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>		<span class="comment">// Reset for next time</span>
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>		p.scannedStackSize = 0
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>		p.scannedStacks = 0
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>	}
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>	if scannedStacks == 0 {
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>		startingStackSize = fixedStack
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>		return
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>	}
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>	avg := scannedStackSize/scannedStacks + stackGuard
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>	<span class="comment">// Note: we add stackGuard to ensure that a goroutine that</span>
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>	<span class="comment">// uses the average space will not trigger a growth.</span>
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>	if avg &gt; uint64(maxstacksize) {
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>		avg = uint64(maxstacksize)
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>	}
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>	if avg &lt; fixedStack {
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>		avg = fixedStack
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>	}
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>	<span class="comment">// Note: maxstacksize fits in 30 bits, so avg also does.</span>
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>	startingStackSize = uint32(round2(int32(avg)))
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>}
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>
</pre><p><a href="stack.go?m=text">View as plain text</a></p>

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
