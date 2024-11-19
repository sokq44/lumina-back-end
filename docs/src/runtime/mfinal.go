<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mfinal.go - Go Documentation Server</title>

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
<a href="mfinal.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mfinal.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Garbage collector: finalizers and block profiling.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/goexperiment&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// finblock is an array of finalizers to be executed. finblocks are</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// arranged in a linked list for the finalizer queue.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// finblock is allocated from non-GC&#39;d memory, so any heap pointers</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// must be specially handled. GC currently assumes that the finalizer</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// queue does not grow during marking (but it can shrink).</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>type finblock struct {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	_       sys.NotInHeap
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	alllink *finblock
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	next    *finblock
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	cnt     uint32
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	_       int32
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	fin     [(_FinBlockSize - 2*goarch.PtrSize - 2*4) / unsafe.Sizeof(finalizer{})]finalizer
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>var fingStatus atomic.Uint32
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// finalizer goroutine status.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>const (
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	fingUninitialized uint32 = iota
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	fingCreated       uint32 = 1 &lt;&lt; (iota - 1)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	fingRunningFinalizer
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	fingWait
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	fingWake
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>var finlock mutex  <span class="comment">// protects the following variables</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>var fing *g        <span class="comment">// goroutine that runs finalizers</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>var finq *finblock <span class="comment">// list of finalizers that are to be executed</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>var finc *finblock <span class="comment">// cache of free blocks</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>var finptrmask [_FinBlockSize / goarch.PtrSize / 8]byte
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>var allfin *finblock <span class="comment">// list of all blocks</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// NOTE: Layout known to queuefinalizer.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>type finalizer struct {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	fn   *funcval       <span class="comment">// function to call (may be a heap pointer)</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	arg  unsafe.Pointer <span class="comment">// ptr to object (may be a heap pointer)</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	nret uintptr        <span class="comment">// bytes of return values from fn</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	fint *_type         <span class="comment">// type of first argument of fn</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	ot   *ptrtype       <span class="comment">// type of ptr to object (may be a heap pointer)</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>var finalizer1 = [...]byte{
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// Each Finalizer is 5 words, ptr ptr INT ptr ptr (INT = uintptr here)</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// Each byte describes 8 words.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// Need 8 Finalizers described by 5 bytes before pattern repeats:</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr INT ptr ptr</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr INT ptr ptr</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr INT ptr ptr</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr INT ptr ptr</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr INT ptr ptr</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr INT ptr ptr</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr INT ptr ptr</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr INT ptr ptr</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// aka</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr INT ptr ptr ptr ptr INT</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr ptr ptr INT ptr ptr ptr</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">//	ptr INT ptr ptr ptr ptr INT ptr</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">//	ptr ptr ptr INT ptr ptr ptr ptr</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">//	INT ptr ptr ptr ptr INT ptr ptr</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// Assumptions about Finalizer layout checked below.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	1&lt;&lt;0 | 1&lt;&lt;1 | 0&lt;&lt;2 | 1&lt;&lt;3 | 1&lt;&lt;4 | 1&lt;&lt;5 | 1&lt;&lt;6 | 0&lt;&lt;7,
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	1&lt;&lt;0 | 1&lt;&lt;1 | 1&lt;&lt;2 | 1&lt;&lt;3 | 0&lt;&lt;4 | 1&lt;&lt;5 | 1&lt;&lt;6 | 1&lt;&lt;7,
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	1&lt;&lt;0 | 0&lt;&lt;1 | 1&lt;&lt;2 | 1&lt;&lt;3 | 1&lt;&lt;4 | 1&lt;&lt;5 | 0&lt;&lt;6 | 1&lt;&lt;7,
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	1&lt;&lt;0 | 1&lt;&lt;1 | 1&lt;&lt;2 | 0&lt;&lt;3 | 1&lt;&lt;4 | 1&lt;&lt;5 | 1&lt;&lt;6 | 1&lt;&lt;7,
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	0&lt;&lt;0 | 1&lt;&lt;1 | 1&lt;&lt;2 | 1&lt;&lt;3 | 1&lt;&lt;4 | 0&lt;&lt;5 | 1&lt;&lt;6 | 1&lt;&lt;7,
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// lockRankMayQueueFinalizer records the lock ranking effects of a</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// function that may call queuefinalizer.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func lockRankMayQueueFinalizer() {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	lockWithRankMayAcquire(&amp;finlock, getLockRank(&amp;finlock))
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>func queuefinalizer(p unsafe.Pointer, fn *funcval, nret uintptr, fint *_type, ot *ptrtype) {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if gcphase != _GCoff {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		<span class="comment">// Currently we assume that the finalizer queue won&#39;t</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		<span class="comment">// grow during marking so we don&#39;t have to rescan it</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		<span class="comment">// during mark termination. If we ever need to lift</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		<span class="comment">// this assumption, we can do it by adding the</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		<span class="comment">// necessary barriers to queuefinalizer (which it may</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		<span class="comment">// have automatically).</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		throw(&#34;queuefinalizer during GC&#34;)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	lock(&amp;finlock)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	if finq == nil || finq.cnt == uint32(len(finq.fin)) {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		if finc == nil {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			finc = (*finblock)(persistentalloc(_FinBlockSize, 0, &amp;memstats.gcMiscSys))
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			finc.alllink = allfin
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			allfin = finc
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			if finptrmask[0] == 0 {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>				<span class="comment">// Build pointer mask for Finalizer array in block.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>				<span class="comment">// Check assumptions made in finalizer1 array above.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>				if (unsafe.Sizeof(finalizer{}) != 5*goarch.PtrSize ||
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>					unsafe.Offsetof(finalizer{}.fn) != 0 ||
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>					unsafe.Offsetof(finalizer{}.arg) != goarch.PtrSize ||
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>					unsafe.Offsetof(finalizer{}.nret) != 2*goarch.PtrSize ||
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>					unsafe.Offsetof(finalizer{}.fint) != 3*goarch.PtrSize ||
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>					unsafe.Offsetof(finalizer{}.ot) != 4*goarch.PtrSize) {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>					throw(&#34;finalizer out of sync&#34;)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>				}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>				for i := range finptrmask {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>					finptrmask[i] = finalizer1[i%len(finalizer1)]
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>				}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		block := finc
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		finc = block.next
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		block.next = finq
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		finq = block
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	f := &amp;finq.fin[finq.cnt]
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	atomic.Xadd(&amp;finq.cnt, +1) <span class="comment">// Sync with markroots</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	f.fn = fn
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	f.nret = nret
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	f.fint = fint
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	f.ot = ot
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	f.arg = p
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	unlock(&amp;finlock)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	fingStatus.Or(fingWake)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>func iterate_finq(callback func(*funcval, unsafe.Pointer, uintptr, *_type, *ptrtype)) {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	for fb := allfin; fb != nil; fb = fb.alllink {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		for i := uint32(0); i &lt; fb.cnt; i++ {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			f := &amp;fb.fin[i]
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			callback(f.fn, f.arg, f.nret, f.fint, f.ot)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>func wakefing() *g {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	if ok := fingStatus.CompareAndSwap(fingCreated|fingWait|fingWake, fingCreated); ok {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		return fing
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	return nil
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>func createfing() {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// start the finalizer goroutine exactly once</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if fingStatus.Load() == fingUninitialized &amp;&amp; fingStatus.CompareAndSwap(fingUninitialized, fingCreated) {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		go runfinq()
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>func finalizercommit(gp *g, lock unsafe.Pointer) bool {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	unlock((*mutex)(lock))
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// fingStatus should be modified after fing is put into a waiting state</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// to avoid waking fing in running state, even if it is about to be parked.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	fingStatus.Or(fingWait)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	return true
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// This is the goroutine that runs all of the finalizers.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>func runfinq() {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	var (
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		frame    unsafe.Pointer
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		framecap uintptr
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		argRegs  int
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	)
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	gp := getg()
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	lock(&amp;finlock)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	fing = gp
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	unlock(&amp;finlock)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	for {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		lock(&amp;finlock)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		fb := finq
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		finq = nil
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		if fb == nil {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			gopark(finalizercommit, unsafe.Pointer(&amp;finlock), waitReasonFinalizerWait, traceBlockSystemGoroutine, 1)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			continue
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		argRegs = intArgRegs
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		unlock(&amp;finlock)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		if raceenabled {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			racefingo()
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		for fb != nil {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			for i := fb.cnt; i &gt; 0; i-- {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>				f := &amp;fb.fin[i-1]
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>				var regs abi.RegArgs
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>				<span class="comment">// The args may be passed in registers or on stack. Even for</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>				<span class="comment">// the register case, we still need the spill slots.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>				<span class="comment">// TODO: revisit if we remove spill slots.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>				<span class="comment">// Unfortunately because we can have an arbitrary</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>				<span class="comment">// amount of returns and it would be complex to try and</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>				<span class="comment">// figure out how many of those can get passed in registers,</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>				<span class="comment">// just conservatively assume none of them do.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>				framesz := unsafe.Sizeof((any)(nil)) + f.nret
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>				if framecap &lt; framesz {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>					<span class="comment">// The frame does not contain pointers interesting for GC,</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>					<span class="comment">// all not yet finalized objects are stored in finq.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>					<span class="comment">// If we do not mark it as FlagNoScan,</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>					<span class="comment">// the last finalized object is not collected.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>					frame = mallocgc(framesz, nil, true)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>					framecap = framesz
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>				}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>				if f.fint == nil {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>					throw(&#34;missing type in runfinq&#34;)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>				}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>				r := frame
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>				if argRegs &gt; 0 {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>					r = unsafe.Pointer(&amp;regs.Ints)
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>				} else {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>					<span class="comment">// frame is effectively uninitialized</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>					<span class="comment">// memory. That means we have to clear</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>					<span class="comment">// it before writing to it to avoid</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>					<span class="comment">// confusing the write barrier.</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>					*(*[2]uintptr)(frame) = [2]uintptr{}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>				}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>				switch f.fint.Kind_ &amp; kindMask {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>				case kindPtr:
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>					<span class="comment">// direct use of pointer</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>					*(*unsafe.Pointer)(r) = f.arg
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>				case kindInterface:
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>					ityp := (*interfacetype)(unsafe.Pointer(f.fint))
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>					<span class="comment">// set up with empty interface</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>					(*eface)(r)._type = &amp;f.ot.Type
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>					(*eface)(r).data = f.arg
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>					if len(ityp.Methods) != 0 {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>						<span class="comment">// convert to interface with methods</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>						<span class="comment">// this conversion is guaranteed to succeed - we checked in SetFinalizer</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>						(*iface)(r).tab = assertE2I(ityp, (*eface)(r)._type)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>					}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>				default:
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>					throw(&#34;bad kind in runfinq&#34;)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>				}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				fingStatus.Or(fingRunningFinalizer)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>				reflectcall(nil, unsafe.Pointer(f.fn), frame, uint32(framesz), uint32(framesz), uint32(framesz), &amp;regs)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>				fingStatus.And(^fingRunningFinalizer)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>				<span class="comment">// Drop finalizer queue heap references</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>				<span class="comment">// before hiding them from markroot.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>				<span class="comment">// This also ensures these will be</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				<span class="comment">// clear if we reuse the finalizer.</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>				f.fn = nil
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>				f.arg = nil
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>				f.ot = nil
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>				atomic.Store(&amp;fb.cnt, i-1)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			next := fb.next
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			lock(&amp;finlock)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			fb.next = finc
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			finc = fb
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			unlock(&amp;finlock)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			fb = next
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>func isGoPointerWithoutSpan(p unsafe.Pointer) bool {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	<span class="comment">// 0-length objects are okay.</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	if p == unsafe.Pointer(&amp;zerobase) {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		return true
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	<span class="comment">// Global initializers might be linker-allocated.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">//	var Foo = &amp;Object{}</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	<span class="comment">//	func main() {</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	<span class="comment">//		runtime.SetFinalizer(Foo, nil)</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	<span class="comment">//	}</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	<span class="comment">// The relevant segments are: noptrdata, data, bss, noptrbss.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">// We cannot assume they are in any order or even contiguous,</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// due to external linking.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	for datap := &amp;firstmoduledata; datap != nil; datap = datap.next {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		if datap.noptrdata &lt;= uintptr(p) &amp;&amp; uintptr(p) &lt; datap.enoptrdata ||
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			datap.data &lt;= uintptr(p) &amp;&amp; uintptr(p) &lt; datap.edata ||
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			datap.bss &lt;= uintptr(p) &amp;&amp; uintptr(p) &lt; datap.ebss ||
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			datap.noptrbss &lt;= uintptr(p) &amp;&amp; uintptr(p) &lt; datap.enoptrbss {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			return true
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	return false
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span><span class="comment">// blockUntilEmptyFinalizerQueue blocks until either the finalizer</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span><span class="comment">// queue is emptied (and the finalizers have executed) or the timeout</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">// is reached. Returns true if the finalizer queue was emptied.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">// This is used by the runtime and sync tests.</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>func blockUntilEmptyFinalizerQueue(timeout int64) bool {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	start := nanotime()
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	for nanotime()-start &lt; timeout {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		lock(&amp;finlock)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		<span class="comment">// We know the queue has been drained when both finq is nil</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		<span class="comment">// and the finalizer g has stopped executing.</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		empty := finq == nil
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		empty = empty &amp;&amp; readgstatus(fing) == _Gwaiting &amp;&amp; fing.waitreason == waitReasonFinalizerWait
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		unlock(&amp;finlock)
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		if empty {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			return true
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		Gosched()
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	return false
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span><span class="comment">// SetFinalizer sets the finalizer associated with obj to the provided</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span><span class="comment">// finalizer function. When the garbage collector finds an unreachable block</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">// with an associated finalizer, it clears the association and runs</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span><span class="comment">// finalizer(obj) in a separate goroutine. This makes obj reachable again,</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// but now without an associated finalizer. Assuming that SetFinalizer</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// is not called again, the next time the garbage collector sees</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">// that obj is unreachable, it will free obj.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// SetFinalizer(obj, nil) clears any finalizer associated with obj.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">// The argument obj must be a pointer to an object allocated by calling</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// new, by taking the address of a composite literal, or by taking the</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">// address of a local variable.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// The argument finalizer must be a function that takes a single argument</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">// to which obj&#39;s type can be assigned, and can have arbitrary ignored return</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// values. If either of these is not true, SetFinalizer may abort the</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// program.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span><span class="comment">// Finalizers are run in dependency order: if A points at B, both have</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">// finalizers, and they are otherwise unreachable, only the finalizer</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">// for A runs; once A is freed, the finalizer for B can run.</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// If a cyclic structure includes a block with a finalizer, that</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// cycle is not guaranteed to be garbage collected and the finalizer</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// is not guaranteed to run, because there is no ordering that</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// respects the dependencies.</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// The finalizer is scheduled to run at some arbitrary time after the</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// program can no longer reach the object to which obj points.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// There is no guarantee that finalizers will run before a program exits,</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// so typically they are useful only for releasing non-memory resources</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// associated with an object during a long-running program.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// For example, an [os.File] object could use a finalizer to close the</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// associated operating system file descriptor when a program discards</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// an os.File without calling Close, but it would be a mistake</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">// to depend on a finalizer to flush an in-memory I/O buffer such as a</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// [bufio.Writer], because the buffer would not be flushed at program exit.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// It is not guaranteed that a finalizer will run if the size of *obj is</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// zero bytes, because it may share same address with other zero-size</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span><span class="comment">// objects in memory. See https://go.dev/ref/spec#Size_and_alignment_guarantees.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// It is not guaranteed that a finalizer will run for objects allocated</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// in initializers for package-level variables. Such objects may be</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">// linker-allocated, not heap-allocated.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span><span class="comment">// Note that because finalizers may execute arbitrarily far into the future</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">// after an object is no longer referenced, the runtime is allowed to perform</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span><span class="comment">// a space-saving optimization that batches objects together in a single</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span><span class="comment">// allocation slot. The finalizer for an unreferenced object in such an</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span><span class="comment">// allocation may never run if it always exists in the same batch as a</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span><span class="comment">// referenced object. Typically, this batching only happens for tiny</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span><span class="comment">// (on the order of 16 bytes or less) and pointer-free objects.</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span><span class="comment">// A finalizer may run as soon as an object becomes unreachable.</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span><span class="comment">// In order to use finalizers correctly, the program must ensure that</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span><span class="comment">// the object is reachable until it is no longer required.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span><span class="comment">// Objects stored in global variables, or that can be found by tracing</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span><span class="comment">// pointers from a global variable, are reachable. For other objects,</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span><span class="comment">// pass the object to a call of the [KeepAlive] function to mark the</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span><span class="comment">// last point in the function where the object must be reachable.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span><span class="comment">// For example, if p points to a struct, such as os.File, that contains</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span><span class="comment">// a file descriptor d, and p has a finalizer that closes that file</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span><span class="comment">// descriptor, and if the last use of p in a function is a call to</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span><span class="comment">// syscall.Write(p.d, buf, size), then p may be unreachable as soon as</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span><span class="comment">// the program enters [syscall.Write]. The finalizer may run at that moment,</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span><span class="comment">// closing p.d, causing syscall.Write to fail because it is writing to</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">// a closed file descriptor (or, worse, to an entirely different</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span><span class="comment">// file descriptor opened by a different goroutine). To avoid this problem,</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">// call KeepAlive(p) after the call to syscall.Write.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// A single goroutine runs all finalizers for a program, sequentially.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">// If a finalizer must run for a long time, it should do so by starting</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span><span class="comment">// a new goroutine.</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span><span class="comment">// In the terminology of the Go memory model, a call</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span><span class="comment">// SetFinalizer(x, f) “synchronizes before” the finalization call f(x).</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span><span class="comment">// However, there is no guarantee that KeepAlive(x) or any other use of x</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span><span class="comment">// “synchronizes before” f(x), so in general a finalizer should use a mutex</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span><span class="comment">// or other synchronization mechanism if it needs to access mutable state in x.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span><span class="comment">// For example, consider a finalizer that inspects a mutable field in x</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span><span class="comment">// that is modified from time to time in the main program before x</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span><span class="comment">// becomes unreachable and the finalizer is invoked.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span><span class="comment">// The modifications in the main program and the inspection in the finalizer</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span><span class="comment">// need to use appropriate synchronization, such as mutexes or atomic updates,</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span><span class="comment">// to avoid read-write races.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>func SetFinalizer(obj any, finalizer any) {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	if debug.sbrk != 0 {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		<span class="comment">// debug.sbrk never frees memory, so no finalizers run</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		<span class="comment">// (and we don&#39;t have the data structures to record them).</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		return
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	e := efaceOf(&amp;obj)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	etyp := e._type
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	if etyp == nil {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		throw(&#34;runtime.SetFinalizer: first argument is nil&#34;)
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	if etyp.Kind_&amp;kindMask != kindPtr {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		throw(&#34;runtime.SetFinalizer: first argument is &#34; + toRType(etyp).string() + &#34;, not pointer&#34;)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	ot := (*ptrtype)(unsafe.Pointer(etyp))
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	if ot.Elem == nil {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		throw(&#34;nil elem type!&#34;)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	if inUserArenaChunk(uintptr(e.data)) {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		<span class="comment">// Arena-allocated objects are not eligible for finalizers.</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		throw(&#34;runtime.SetFinalizer: first argument was allocated into an arena&#34;)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	<span class="comment">// find the containing object</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	base, span, _ := findObject(uintptr(e.data), 0, 0)
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	if base == 0 {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		if isGoPointerWithoutSpan(e.data) {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>			return
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		throw(&#34;runtime.SetFinalizer: pointer not in allocated block&#34;)
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	<span class="comment">// Move base forward if we&#39;ve got an allocation header.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	if goexperiment.AllocHeaders &amp;&amp; !span.spanclass.noscan() &amp;&amp; !heapBitsInSpan(span.elemsize) &amp;&amp; span.spanclass.sizeclass() != 0 {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		base += mallocHeaderSize
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	if uintptr(e.data) != base {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		<span class="comment">// As an implementation detail we allow to set finalizers for an inner byte</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		<span class="comment">// of an object if it could come from tiny alloc (see mallocgc for details).</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		if ot.Elem == nil || ot.Elem.PtrBytes != 0 || ot.Elem.Size_ &gt;= maxTinySize {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			throw(&#34;runtime.SetFinalizer: pointer not at beginning of allocated block&#34;)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	f := efaceOf(&amp;finalizer)
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	ftyp := f._type
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	if ftyp == nil {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		<span class="comment">// switch to system stack and remove finalizer</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>			removefinalizer(e.data)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		})
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		return
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	if ftyp.Kind_&amp;kindMask != kindFunc {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		throw(&#34;runtime.SetFinalizer: second argument is &#34; + toRType(ftyp).string() + &#34;, not a function&#34;)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	ft := (*functype)(unsafe.Pointer(ftyp))
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	if ft.IsVariadic() {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		throw(&#34;runtime.SetFinalizer: cannot pass &#34; + toRType(etyp).string() + &#34; to finalizer &#34; + toRType(ftyp).string() + &#34; because dotdotdot&#34;)
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	if ft.InCount != 1 {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		throw(&#34;runtime.SetFinalizer: cannot pass &#34; + toRType(etyp).string() + &#34; to finalizer &#34; + toRType(ftyp).string())
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	fint := ft.InSlice()[0]
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	switch {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	case fint == etyp:
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		<span class="comment">// ok - same type</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		goto okarg
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	case fint.Kind_&amp;kindMask == kindPtr:
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		if (fint.Uncommon() == nil || etyp.Uncommon() == nil) &amp;&amp; (*ptrtype)(unsafe.Pointer(fint)).Elem == ot.Elem {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			<span class="comment">// ok - not same type, but both pointers,</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			<span class="comment">// one or the other is unnamed, and same element type, so assignable.</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			goto okarg
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	case fint.Kind_&amp;kindMask == kindInterface:
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		ityp := (*interfacetype)(unsafe.Pointer(fint))
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		if len(ityp.Methods) == 0 {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>			<span class="comment">// ok - satisfies empty interface</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>			goto okarg
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		if itab := assertE2I2(ityp, efaceOf(&amp;obj)._type); itab != nil {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			goto okarg
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	}
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	throw(&#34;runtime.SetFinalizer: cannot pass &#34; + toRType(etyp).string() + &#34; to finalizer &#34; + toRType(ftyp).string())
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>okarg:
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	<span class="comment">// compute size needed for return parameters</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	nret := uintptr(0)
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	for _, t := range ft.OutSlice() {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		nret = alignUp(nret, uintptr(t.Align_)) + t.Size_
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	nret = alignUp(nret, goarch.PtrSize)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	<span class="comment">// make sure we have a finalizer goroutine</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	createfing()
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		if !addfinalizer(e.data, (*funcval)(f.data), nret, fint, ot) {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			throw(&#34;runtime.SetFinalizer: finalizer already set&#34;)
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	})
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>}
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span><span class="comment">// Mark KeepAlive as noinline so that it is easily detectable as an intrinsic.</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span><span class="comment">//go:noinline</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span><span class="comment">// KeepAlive marks its argument as currently reachable.</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span><span class="comment">// This ensures that the object is not freed, and its finalizer is not run,</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span><span class="comment">// before the point in the program where KeepAlive is called.</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span><span class="comment">// A very simplified example showing where KeepAlive is required:</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span><span class="comment">//	type File struct { d int }</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span><span class="comment">//	d, err := syscall.Open(&#34;/file/path&#34;, syscall.O_RDONLY, 0)</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span><span class="comment">//	// ... do something if err != nil ...</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span><span class="comment">//	p := &amp;File{d}</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span><span class="comment">//	runtime.SetFinalizer(p, func(p *File) { syscall.Close(p.d) })</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span><span class="comment">//	var buf [10]byte</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span><span class="comment">//	n, err := syscall.Read(p.d, buf[:])</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span><span class="comment">//	// Ensure p is not finalized until Read returns.</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span><span class="comment">//	runtime.KeepAlive(p)</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span><span class="comment">//	// No more uses of p after this point.</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">// Without the KeepAlive call, the finalizer could run at the start of</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">// [syscall.Read], closing the file descriptor before syscall.Read makes</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">// the actual system call.</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">// Note: KeepAlive should only be used to prevent finalizers from</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">// running prematurely. In particular, when used with [unsafe.Pointer],</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span><span class="comment">// the rules for valid uses of unsafe.Pointer still apply.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>func KeepAlive(x any) {
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	<span class="comment">// Introduce a use of x that the compiler can&#39;t eliminate.</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	<span class="comment">// This makes sure x is alive on entry. We need x to be alive</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	<span class="comment">// on entry for &#34;defer runtime.KeepAlive(x)&#34;; see issue 21402.</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	if cgoAlwaysFalse {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		println(x)
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>}
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>
</pre><p><a href="mfinal.go?m=text">View as plain text</a></p>

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
