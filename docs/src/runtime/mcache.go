<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mcache.go - Go Documentation Server</title>

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
<a href="mcache.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mcache.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// Per-thread (in Go, per-P) cache for small objects.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// This includes a small object cache and local allocation stats.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// No locking needed because it is per-thread (per-P).</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// mcaches are allocated from non-GC&#39;d memory, so any heap pointers</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// must be specially handled.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>type mcache struct {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// The following members are accessed on every malloc,</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// so they are grouped here for better caching.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	nextSample uintptr <span class="comment">// trigger heap sample after allocating this many bytes</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	scanAlloc  uintptr <span class="comment">// bytes of scannable heap allocated</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// Allocator cache for tiny objects w/o pointers.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// See &#34;Tiny allocator&#34; comment in malloc.go.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// tiny points to the beginning of the current tiny block, or</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// nil if there is no current tiny block.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// tiny is a heap pointer. Since mcache is in non-GC&#39;d memory,</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// we handle it by clearing it in releaseAll during mark</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// termination.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// tinyAllocs is the number of tiny allocations performed</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// by the P that owns this mcache.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	tiny       uintptr
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	tinyoffset uintptr
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	tinyAllocs uintptr
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// The rest is not accessed on every malloc.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	alloc [numSpanClasses]*mspan <span class="comment">// spans to allocate from, indexed by spanClass</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	stackcache [_NumStackOrders]stackfreelist
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// flushGen indicates the sweepgen during which this mcache</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// was last flushed. If flushGen != mheap_.sweepgen, the spans</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// in this mcache are stale and need to the flushed so they</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// can be swept. This is done in acquirep.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	flushGen atomic.Uint32
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// A gclink is a node in a linked list of blocks, like mlink,</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// but it is opaque to the garbage collector.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// The GC does not trace the pointers during collection,</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// and the compiler does not emit write barriers for assignments</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// of gclinkptr values. Code should store references to gclinks</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// as gclinkptr, not as *gclink.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>type gclink struct {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	next gclinkptr
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// A gclinkptr is a pointer to a gclink, but it is opaque</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// to the garbage collector.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>type gclinkptr uintptr
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// ptr returns the *gclink form of p.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// The result should be used for accessing fields, not stored</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// in other data structures.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>func (p gclinkptr) ptr() *gclink {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return (*gclink)(unsafe.Pointer(p))
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>type stackfreelist struct {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	list gclinkptr <span class="comment">// linked list of free stacks</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	size uintptr   <span class="comment">// total size of stacks in list</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// dummy mspan that contains no free objects.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>var emptymspan mspan
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>func allocmcache() *mcache {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	var c *mcache
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		lock(&amp;mheap_.lock)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		c = (*mcache)(mheap_.cachealloc.alloc())
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		c.flushGen.Store(mheap_.sweepgen)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		unlock(&amp;mheap_.lock)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	})
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	for i := range c.alloc {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		c.alloc[i] = &amp;emptymspan
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	c.nextSample = nextSample()
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	return c
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// freemcache releases resources associated with this</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// mcache and puts the object onto a free list.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// In some cases there is no way to simply release</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// resources, such as statistics, so donate them to</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// a different mcache (the recipient).</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>func freemcache(c *mcache) {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		c.releaseAll()
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		stackcache_clear(c)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		<span class="comment">// NOTE(rsc,rlh): If gcworkbuffree comes back, we need to coordinate</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		<span class="comment">// with the stealing of gcworkbufs during garbage collection to avoid</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		<span class="comment">// a race where the workbuf is double-freed.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		<span class="comment">// gcworkbuffree(c.gcworkbuf)</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		lock(&amp;mheap_.lock)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		mheap_.cachealloc.free(unsafe.Pointer(c))
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		unlock(&amp;mheap_.lock)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	})
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// getMCache is a convenience function which tries to obtain an mcache.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// Returns nil if we&#39;re not bootstrapping or we don&#39;t have a P. The caller&#39;s</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// P must not change, so we must be in a non-preemptible state.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>func getMCache(mp *m) *mcache {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// Grab the mcache, since that&#39;s where stats live.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	pp := mp.p.ptr()
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	var c *mcache
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	if pp == nil {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		<span class="comment">// We will be called without a P while bootstrapping,</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		<span class="comment">// in which case we use mcache0, which is set in mallocinit.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		<span class="comment">// mcache0 is cleared when bootstrapping is complete,</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// by procresize.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		c = mcache0
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	} else {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		c = pp.mcache
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	return c
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// refill acquires a new span of span class spc for c. This span will</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// have at least one free object. The current span in c must be full.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// Must run in a non-preemptible context since otherwise the owner of</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// c could change.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>func (c *mcache) refill(spc spanClass) {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// Return the current cached span to the central lists.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	s := c.alloc[spc]
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	if s.allocCount != s.nelems {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		throw(&#34;refill of span with free space remaining&#34;)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	if s != &amp;emptymspan {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		<span class="comment">// Mark this span as no longer cached.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		if s.sweepgen != mheap_.sweepgen+3 {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			throw(&#34;bad sweepgen in refill&#34;)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		mheap_.central[spc].mcentral.uncacheSpan(s)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		<span class="comment">// Count up how many slots were used and record it.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		stats := memstats.heapStats.acquire()
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		slotsUsed := int64(s.allocCount) - int64(s.allocCountBeforeCache)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		atomic.Xadd64(&amp;stats.smallAllocCount[spc.sizeclass()], slotsUsed)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		<span class="comment">// Flush tinyAllocs.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		if spc == tinySpanClass {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			atomic.Xadd64(&amp;stats.tinyAllocCount, int64(c.tinyAllocs))
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			c.tinyAllocs = 0
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		memstats.heapStats.release()
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		<span class="comment">// Count the allocs in inconsistent, internal stats.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		bytesAllocated := slotsUsed * int64(s.elemsize)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		gcController.totalAlloc.Add(bytesAllocated)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		<span class="comment">// Clear the second allocCount just to be safe.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		s.allocCountBeforeCache = 0
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// Get a new cached span from the central lists.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	s = mheap_.central[spc].mcentral.cacheSpan()
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	if s == nil {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		throw(&#34;out of memory&#34;)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	if s.allocCount == s.nelems {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		throw(&#34;span has no free space&#34;)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// Indicate that this span is cached and prevent asynchronous</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// sweeping in the next sweep phase.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	s.sweepgen = mheap_.sweepgen + 3
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// Store the current alloc count for accounting later.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	s.allocCountBeforeCache = s.allocCount
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// Update heapLive and flush scanAlloc.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// We have not yet allocated anything new into the span, but we</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// assume that all of its slots will get used, so this makes</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// heapLive an overestimate.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// When the span gets uncached, we&#39;ll fix up this overestimate</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// if necessary (see releaseAll).</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// We pick an overestimate here because an underestimate leads</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// the pacer to believe that it&#39;s in better shape than it is,</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">// which appears to lead to more memory used. See #53738 for</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">// more details.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	usedBytes := uintptr(s.allocCount) * s.elemsize
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	gcController.update(int64(s.npages*pageSize)-int64(usedBytes), int64(c.scanAlloc))
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	c.scanAlloc = 0
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	c.alloc[spc] = s
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// allocLarge allocates a span for a large object.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>func (c *mcache) allocLarge(size uintptr, noscan bool) *mspan {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	if size+_PageSize &lt; size {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		throw(&#34;out of memory&#34;)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	npages := size &gt;&gt; _PageShift
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	if size&amp;_PageMask != 0 {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		npages++
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// Deduct credit for this span allocation and sweep if</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// necessary. mHeap_Alloc will also sweep npages, so this only</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// pays the debt down to npage pages.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	deductSweepCredit(npages*_PageSize, npages)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	spc := makeSpanClass(0, noscan)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	s := mheap_.alloc(npages, spc)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	if s == nil {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		throw(&#34;out of memory&#34;)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	<span class="comment">// Count the alloc in consistent, external stats.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	stats := memstats.heapStats.acquire()
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	atomic.Xadd64(&amp;stats.largeAlloc, int64(npages*pageSize))
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	atomic.Xadd64(&amp;stats.largeAllocCount, 1)
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	memstats.heapStats.release()
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// Count the alloc in inconsistent, internal stats.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	gcController.totalAlloc.Add(int64(npages * pageSize))
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// Update heapLive.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	gcController.update(int64(s.npages*pageSize), 0)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// Put the large span in the mcentral swept list so that it&#39;s</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// visible to the background sweeper.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	mheap_.central[spc].mcentral.fullSwept(mheap_.sweepgen).push(s)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	s.limit = s.base() + size
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	s.initHeapBits(false)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	return s
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>func (c *mcache) releaseAll() {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// Take this opportunity to flush scanAlloc.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	scanAlloc := int64(c.scanAlloc)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	c.scanAlloc = 0
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	sg := mheap_.sweepgen
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	dHeapLive := int64(0)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	for i := range c.alloc {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		s := c.alloc[i]
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		if s != &amp;emptymspan {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			slotsUsed := int64(s.allocCount) - int64(s.allocCountBeforeCache)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			s.allocCountBeforeCache = 0
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			<span class="comment">// Adjust smallAllocCount for whatever was allocated.</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			stats := memstats.heapStats.acquire()
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			atomic.Xadd64(&amp;stats.smallAllocCount[spanClass(i).sizeclass()], slotsUsed)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			memstats.heapStats.release()
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			<span class="comment">// Adjust the actual allocs in inconsistent, internal stats.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			<span class="comment">// We assumed earlier that the full span gets allocated.</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			gcController.totalAlloc.Add(slotsUsed * int64(s.elemsize))
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			if s.sweepgen != sg+1 {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>				<span class="comment">// refill conservatively counted unallocated slots in gcController.heapLive.</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>				<span class="comment">// Undo this.</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>				<span class="comment">// If this span was cached before sweep, then gcController.heapLive was totally</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>				<span class="comment">// recomputed since caching this span, so we don&#39;t do this for stale spans.</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>				dHeapLive -= int64(s.nelems-s.allocCount) * int64(s.elemsize)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			<span class="comment">// Release the span to the mcentral.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			mheap_.central[i].mcentral.uncacheSpan(s)
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			c.alloc[i] = &amp;emptymspan
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// Clear tinyalloc pool.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	c.tiny = 0
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	c.tinyoffset = 0
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">// Flush tinyAllocs.</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	stats := memstats.heapStats.acquire()
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	atomic.Xadd64(&amp;stats.tinyAllocCount, int64(c.tinyAllocs))
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	c.tinyAllocs = 0
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	memstats.heapStats.release()
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">// Update heapLive and heapScan.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	gcController.update(dHeapLive, scanAlloc)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// prepareForSweep flushes c if the system has entered a new sweep phase</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// since c was populated. This must happen between the sweep phase</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// starting and the first allocation from c.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>func (c *mcache) prepareForSweep() {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	<span class="comment">// Alternatively, instead of making sure we do this on every P</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	<span class="comment">// between starting the world and allocating on that P, we</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	<span class="comment">// could leave allocate-black on, allow allocation to continue</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// as usual, use a ragged barrier at the beginning of sweep to</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	<span class="comment">// ensure all cached spans are swept, and then disable</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	<span class="comment">// allocate-black. However, with this approach it&#39;s difficult</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	<span class="comment">// to avoid spilling mark bits into the *next* GC cycle.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	sg := mheap_.sweepgen
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	flushGen := c.flushGen.Load()
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	if flushGen == sg {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		return
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	} else if flushGen != sg-2 {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		println(&#34;bad flushGen&#34;, flushGen, &#34;in prepareForSweep; sweepgen&#34;, sg)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		throw(&#34;bad flushGen&#34;)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	c.releaseAll()
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	stackcache_clear(c)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	c.flushGen.Store(mheap_.sweepgen) <span class="comment">// Synchronizes with gcStart</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
</pre><p><a href="mcache.go?m=text">View as plain text</a></p>

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
