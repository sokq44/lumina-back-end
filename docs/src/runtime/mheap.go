<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mheap.go - Go Documentation Server</title>

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
<a href="mheap.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mheap.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Page heap.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// See malloc.go for overview.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package runtime
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/cpu&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;internal/goexperiment&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>const (
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// minPhysPageSize is a lower-bound on the physical page size. The</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// true physical page size may be larger than this. In contrast,</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// sys.PhysPageSize is an upper-bound on the physical page size.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	minPhysPageSize = 4096
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// maxPhysPageSize is the maximum page size the runtime supports.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	maxPhysPageSize = 512 &lt;&lt; 10
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// maxPhysHugePageSize sets an upper-bound on the maximum huge page size</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// that the runtime supports.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	maxPhysHugePageSize = pallocChunkBytes
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// pagesPerReclaimerChunk indicates how many pages to scan from the</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// pageInUse bitmap at a time. Used by the page reclaimer.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// Higher values reduce contention on scanning indexes (such as</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// h.reclaimIndex), but increase the minimum latency of the</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// operation.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// The time required to scan this many pages can vary a lot depending</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// on how many spans are actually freed. Experimentally, it can</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// scan for pages at ~300 GB/ms on a 2.6GHz Core i7, but can only</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// free spans at ~32 MB/ms. Using 512 pages bounds this at</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// roughly 100µs.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// Must be a multiple of the pageInUse bitmap element size and</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// must also evenly divide pagesPerArena.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	pagesPerReclaimerChunk = 512
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// physPageAlignedStacks indicates whether stack allocations must be</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// physical page aligned. This is a requirement for MAP_STACK on</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// OpenBSD.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	physPageAlignedStacks = GOOS == &#34;openbsd&#34;
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// Main malloc heap.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// The heap itself is the &#34;free&#34; and &#34;scav&#34; treaps,</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// but all the other global data is here too.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// mheap must not be heap-allocated because it contains mSpanLists,</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// which must not be heap-allocated.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>type mheap struct {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// lock must only be acquired on the system stack, otherwise a g</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// could self-deadlock if its stack grows with the lock held.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	lock mutex
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	pages pageAlloc <span class="comment">// page allocation data structure</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	sweepgen uint32 <span class="comment">// sweep generation, see comment in mspan; written during STW</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// allspans is a slice of all mspans ever created. Each mspan</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// appears exactly once.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// The memory for allspans is manually managed and can be</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// reallocated and move as the heap grows.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// In general, allspans is protected by mheap_.lock, which</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// prevents concurrent access as well as freeing the backing</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// store. Accesses during STW might not hold the lock, but</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// must ensure that allocation cannot happen around the</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// access (since that may free the backing store).</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	allspans []*mspan <span class="comment">// all spans out there</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// Proportional sweep</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// These parameters represent a linear function from gcController.heapLive</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// to page sweep count. The proportional sweep system works to</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// stay in the black by keeping the current page sweep count</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// above this line at the current gcController.heapLive.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// The line has slope sweepPagesPerByte and passes through a</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// basis point at (sweepHeapLiveBasis, pagesSweptBasis). At</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// any given time, the system is at (gcController.heapLive,</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// pagesSwept) in this space.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// It is important that the line pass through a point we</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// control rather than simply starting at a 0,0 origin</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// because that lets us adjust sweep pacing at any time while</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// accounting for current progress. If we could only adjust</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// the slope, it would create a discontinuity in debt if any</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// progress has already been made.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	pagesInUse         atomic.Uintptr <span class="comment">// pages of spans in stats mSpanInUse</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	pagesSwept         atomic.Uint64  <span class="comment">// pages swept this cycle</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	pagesSweptBasis    atomic.Uint64  <span class="comment">// pagesSwept to use as the origin of the sweep ratio</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	sweepHeapLiveBasis uint64         <span class="comment">// value of gcController.heapLive to use as the origin of sweep ratio; written with lock, read without</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	sweepPagesPerByte  float64        <span class="comment">// proportional sweep ratio; written with lock, read without</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// Page reclaimer state</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// reclaimIndex is the page index in allArenas of next page to</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// reclaim. Specifically, it refers to page (i %</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// pagesPerArena) of arena allArenas[i / pagesPerArena].</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// If this is &gt;= 1&lt;&lt;63, the page reclaimer is done scanning</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// the page marks.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	reclaimIndex atomic.Uint64
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// reclaimCredit is spare credit for extra pages swept. Since</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// the page reclaimer works in large chunks, it may reclaim</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// more than requested. Any spare pages released go to this</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// credit pool.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	reclaimCredit atomic.Uintptr
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	_ cpu.CacheLinePad <span class="comment">// prevents false-sharing between arenas and preceding variables</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// arenas is the heap arena map. It points to the metadata for</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// the heap for every arena frame of the entire usable virtual</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// address space.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// Use arenaIndex to compute indexes into this array.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// For regions of the address space that are not backed by the</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// Go heap, the arena map contains nil.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// Modifications are protected by mheap_.lock. Reads can be</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// performed without locking; however, a given entry can</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// transition from nil to non-nil at any time when the lock</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// isn&#39;t held. (Entries never transitions back to nil.)</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// In general, this is a two-level mapping consisting of an L1</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">// map and possibly many L2 maps. This saves space when there</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// are a huge number of arena frames. However, on many</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// platforms (even 64-bit), arenaL1Bits is 0, making this</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// effectively a single-level map. In this case, arenas[0]</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">// will never be nil.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	arenas [1 &lt;&lt; arenaL1Bits]*[1 &lt;&lt; arenaL2Bits]*heapArena
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// arenasHugePages indicates whether arenas&#39; L2 entries are eligible</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// to be backed by huge pages.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	arenasHugePages bool
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// heapArenaAlloc is pre-reserved space for allocating heapArena</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// objects. This is only used on 32-bit, where we pre-reserve</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// this space to avoid interleaving it with the heap itself.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	heapArenaAlloc linearAlloc
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// arenaHints is a list of addresses at which to attempt to</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// add more heap arenas. This is initially populated with a</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// set of general hint addresses, and grown with the bounds of</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// actual heap arena ranges.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	arenaHints *arenaHint
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// arena is a pre-reserved space for allocating heap arenas</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// (the actual arenas). This is only used on 32-bit.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	arena linearAlloc
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// allArenas is the arenaIndex of every mapped arena. This can</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// be used to iterate through the address space.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// Access is protected by mheap_.lock. However, since this is</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// append-only and old backing arrays are never freed, it is</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// safe to acquire mheap_.lock, copy the slice header, and</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// then release mheap_.lock.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	allArenas []arenaIdx
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// sweepArenas is a snapshot of allArenas taken at the</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// beginning of the sweep cycle. This can be read safely by</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// simply blocking GC (by disabling preemption).</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	sweepArenas []arenaIdx
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// markArenas is a snapshot of allArenas taken at the beginning</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// of the mark cycle. Because allArenas is append-only, neither</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// this slice nor its contents will change during the mark, so</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// it can be read safely.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	markArenas []arenaIdx
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// curArena is the arena that the heap is currently growing</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// into. This should always be physPageSize-aligned.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	curArena struct {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		base, end uintptr
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// central free lists for small size classes.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// the padding makes sure that the mcentrals are</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// spaced CacheLinePadSize bytes apart, so that each mcentral.lock</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// gets its own cache line.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// central is indexed by spanClass.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	central [numSpanClasses]struct {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		mcentral mcentral
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		pad      [(cpu.CacheLinePadSize - unsafe.Sizeof(mcentral{})%cpu.CacheLinePadSize) % cpu.CacheLinePadSize]byte
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	spanalloc              fixalloc <span class="comment">// allocator for span*</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	cachealloc             fixalloc <span class="comment">// allocator for mcache*</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	specialfinalizeralloc  fixalloc <span class="comment">// allocator for specialfinalizer*</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	specialprofilealloc    fixalloc <span class="comment">// allocator for specialprofile*</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	specialReachableAlloc  fixalloc <span class="comment">// allocator for specialReachable</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	specialPinCounterAlloc fixalloc <span class="comment">// allocator for specialPinCounter</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	speciallock            mutex    <span class="comment">// lock for special record allocators.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	arenaHintAlloc         fixalloc <span class="comment">// allocator for arenaHints</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// User arena state.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	<span class="comment">// Protected by mheap_.lock.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	userArena struct {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// arenaHints is a list of addresses at which to attempt to</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		<span class="comment">// add more heap arenas for user arena chunks. This is initially</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		<span class="comment">// populated with a set of general hint addresses, and grown with</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		<span class="comment">// the bounds of actual heap arena ranges.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		arenaHints *arenaHint
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		<span class="comment">// quarantineList is a list of user arena spans that have been set to fault, but</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		<span class="comment">// are waiting for all pointers into them to go away. Sweeping handles</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		<span class="comment">// identifying when this is true, and moves the span to the ready list.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		quarantineList mSpanList
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		<span class="comment">// readyList is a list of empty user arena spans that are ready for reuse.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		readyList mSpanList
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	unused *specialfinalizer <span class="comment">// never set, just here to force the specialfinalizer type into DWARF</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>var mheap_ mheap
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">// A heapArena stores metadata for a heap arena. heapArenas are stored</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">// outside of the Go heap and accessed via the mheap_.arenas index.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>type heapArena struct {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// heapArenaPtrScalar contains pointer/scalar data about the heap for this heap arena.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	heapArenaPtrScalar
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">// spans maps from virtual address page ID within this arena to *mspan.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	<span class="comment">// For allocated spans, their pages map to the span itself.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// For free spans, only the lowest and highest pages map to the span itself.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	<span class="comment">// Internal pages map to an arbitrary span.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// For pages that have never been allocated, spans entries are nil.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// Modifications are protected by mheap.lock. Reads can be</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// performed without locking, but ONLY from indexes that are</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	<span class="comment">// known to contain in-use or stack spans. This means there</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	<span class="comment">// must not be a safe-point between establishing that an</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// address is live and looking it up in the spans array.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	spans [pagesPerArena]*mspan
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">// pageInUse is a bitmap that indicates which spans are in</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// state mSpanInUse. This bitmap is indexed by page number,</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// but only the bit corresponding to the first page in each</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// span is used.</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// Reads and writes are atomic.</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	pageInUse [pagesPerArena / 8]uint8
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	<span class="comment">// pageMarks is a bitmap that indicates which spans have any</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// marked objects on them. Like pageInUse, only the bit</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	<span class="comment">// corresponding to the first page in each span is used.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	<span class="comment">// Writes are done atomically during marking. Reads are</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	<span class="comment">// non-atomic and lock-free since they only occur during</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	<span class="comment">// sweeping (and hence never race with writes).</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	<span class="comment">// This is used to quickly find whole spans that can be freed.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	<span class="comment">// TODO(austin): It would be nice if this was uint64 for</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	<span class="comment">// faster scanning, but we don&#39;t have 64-bit atomic bit</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	<span class="comment">// operations.</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	pageMarks [pagesPerArena / 8]uint8
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	<span class="comment">// pageSpecials is a bitmap that indicates which spans have</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	<span class="comment">// specials (finalizers or other). Like pageInUse, only the bit</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	<span class="comment">// corresponding to the first page in each span is used.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	<span class="comment">// Writes are done atomically whenever a special is added to</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	<span class="comment">// a span and whenever the last special is removed from a span.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	<span class="comment">// Reads are done atomically to find spans containing specials</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	<span class="comment">// during marking.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	pageSpecials [pagesPerArena / 8]uint8
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// checkmarks stores the debug.gccheckmark state. It is only</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// used if debug.gccheckmark &gt; 0.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	checkmarks *checkmarksMap
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">// zeroedBase marks the first byte of the first page in this</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// arena which hasn&#39;t been used yet and is therefore already</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	<span class="comment">// zero. zeroedBase is relative to the arena base.</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">// Increases monotonically until it hits heapArenaBytes.</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">// This field is sufficient to determine if an allocation</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	<span class="comment">// needs to be zeroed because the page allocator follows an</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	<span class="comment">// address-ordered first-fit policy.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">// Read atomically and written with an atomic CAS.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	zeroedBase uintptr
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// arenaHint is a hint for where to grow the heap arenas. See</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// mheap_.arenaHints.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>type arenaHint struct {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	_    sys.NotInHeap
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	addr uintptr
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	down bool
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	next *arenaHint
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">// An mspan is a run of pages.</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span><span class="comment">// When a mspan is in the heap free treap, state == mSpanFree</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span><span class="comment">// and heapmap(s-&gt;start) == span, heapmap(s-&gt;start+s-&gt;npages-1) == span.</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span><span class="comment">// If the mspan is in the heap scav treap, then in addition to the</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span><span class="comment">// above scavenged == true. scavenged == false in all other cases.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span><span class="comment">// When a mspan is allocated, state == mSpanInUse or mSpanManual</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">// and heapmap(i) == span for all s-&gt;start &lt;= i &lt; s-&gt;start+s-&gt;npages.</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// Every mspan is in one doubly-linked list, either in the mheap&#39;s</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// busy list or one of the mcentral&#39;s span lists.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span><span class="comment">// An mspan representing actual memory has state mSpanInUse,</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// mSpanManual, or mSpanFree. Transitions between these states are</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">// constrained as follows:</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">//   - A span may transition from free to in-use or manual during any GC</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">//     phase.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">//   - During sweeping (gcphase == _GCoff), a span may transition from</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">//     in-use to free (as a result of sweeping) or manual to free (as a</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">//     result of stacks being freed).</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span><span class="comment">//   - During GC (gcphase != _GCoff), a span *must not* transition from</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">//     manual or in-use to free. Because concurrent GC may read a pointer</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">//     and then look up its span, the span state must be monotonic.</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// Setting mspan.state to mSpanInUse or mSpanManual must be done</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// atomically and only after all other span fields are valid.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// Likewise, if inspecting a span is contingent on it being</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// mSpanInUse, the state should be loaded atomically and checked</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// before depending on other fields. This allows the garbage collector</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// to safely deal with potentially invalid pointers, since resolving</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// such pointers may race with a span being allocated.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>type mSpanState uint8
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>const (
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	mSpanDead   mSpanState = iota
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	mSpanInUse             <span class="comment">// allocated for garbage collected heap</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	mSpanManual            <span class="comment">// allocated for manual management (e.g., stack allocator)</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// mSpanStateNames are the names of the span states, indexed by</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// mSpanState.</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>var mSpanStateNames = []string{
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	&#34;mSpanDead&#34;,
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	&#34;mSpanInUse&#34;,
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	&#34;mSpanManual&#34;,
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span><span class="comment">// mSpanStateBox holds an atomic.Uint8 to provide atomic operations on</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">// an mSpanState. This is a separate type to disallow accidental comparison</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span><span class="comment">// or assignment with mSpanState.</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>type mSpanStateBox struct {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	s atomic.Uint8
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span><span class="comment">// It is nosplit to match get, below.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>func (b *mSpanStateBox) set(s mSpanState) {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	b.s.Store(uint8(s))
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span><span class="comment">// It is nosplit because it&#39;s called indirectly by typedmemclr,</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span><span class="comment">// which must not be preempted.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>func (b *mSpanStateBox) get() mSpanState {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	return mSpanState(b.s.Load())
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">// mSpanList heads a linked list of spans.</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>type mSpanList struct {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	_     sys.NotInHeap
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	first *mspan <span class="comment">// first span in list, or nil if none</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	last  *mspan <span class="comment">// last span in list, or nil if none</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>type mspan struct {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	_    sys.NotInHeap
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	next *mspan     <span class="comment">// next span in list, or nil if none</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	prev *mspan     <span class="comment">// previous span in list, or nil if none</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	list *mSpanList <span class="comment">// For debugging.</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	startAddr uintptr <span class="comment">// address of first byte of span aka s.base()</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	npages    uintptr <span class="comment">// number of pages in span</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	manualFreeList gclinkptr <span class="comment">// list of free objects in mSpanManual spans</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	<span class="comment">// freeindex is the slot index between 0 and nelems at which to begin scanning</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	<span class="comment">// for the next free object in this span.</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	<span class="comment">// Each allocation scans allocBits starting at freeindex until it encounters a 0</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	<span class="comment">// indicating a free object. freeindex is then adjusted so that subsequent scans begin</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	<span class="comment">// just past the newly discovered free object.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	<span class="comment">// If freeindex == nelem, this span has no free objects.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	<span class="comment">// allocBits is a bitmap of objects in this span.</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	<span class="comment">// If n &gt;= freeindex and allocBits[n/8] &amp; (1&lt;&lt;(n%8)) is 0</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	<span class="comment">// then object n is free;</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	<span class="comment">// otherwise, object n is allocated. Bits starting at nelem are</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	<span class="comment">// undefined and should never be referenced.</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	<span class="comment">// Object n starts at address n*elemsize + (start &lt;&lt; pageShift).</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	freeindex uint16
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	<span class="comment">// TODO: Look up nelems from sizeclass and remove this field if it</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	<span class="comment">// helps performance.</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	nelems uint16 <span class="comment">// number of object in the span.</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	<span class="comment">// freeIndexForScan is like freeindex, except that freeindex is</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	<span class="comment">// used by the allocator whereas freeIndexForScan is used by the</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	<span class="comment">// GC scanner. They are two fields so that the GC sees the object</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	<span class="comment">// is allocated only when the object and the heap bits are</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	<span class="comment">// initialized (see also the assignment of freeIndexForScan in</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	<span class="comment">// mallocgc, and issue 54596).</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	freeIndexForScan uint16
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	<span class="comment">// Cache of the allocBits at freeindex. allocCache is shifted</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	<span class="comment">// such that the lowest bit corresponds to the bit freeindex.</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	<span class="comment">// allocCache holds the complement of allocBits, thus allowing</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	<span class="comment">// ctz (count trailing zero) to use it directly.</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	<span class="comment">// allocCache may contain bits beyond s.nelems; the caller must ignore</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// these.</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	allocCache uint64
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	<span class="comment">// allocBits and gcmarkBits hold pointers to a span&#39;s mark and</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	<span class="comment">// allocation bits. The pointers are 8 byte aligned.</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	<span class="comment">// There are three arenas where this data is held.</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	<span class="comment">// free: Dirty arenas that are no longer accessed</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	<span class="comment">//       and can be reused.</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	<span class="comment">// next: Holds information to be used in the next GC cycle.</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	<span class="comment">// current: Information being used during this GC cycle.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	<span class="comment">// previous: Information being used during the last GC cycle.</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	<span class="comment">// A new GC cycle starts with the call to finishsweep_m.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	<span class="comment">// finishsweep_m moves the previous arena to the free arena,</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	<span class="comment">// the current arena to the previous arena, and</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	<span class="comment">// the next arena to the current arena.</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	<span class="comment">// The next arena is populated as the spans request</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	<span class="comment">// memory to hold gcmarkBits for the next GC cycle as well</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	<span class="comment">// as allocBits for newly allocated spans.</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	<span class="comment">// The pointer arithmetic is done &#34;by hand&#34; instead of using</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	<span class="comment">// arrays to avoid bounds checks along critical performance</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	<span class="comment">// paths.</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	<span class="comment">// The sweep will free the old allocBits and set allocBits to the</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	<span class="comment">// gcmarkBits. The gcmarkBits are replaced with a fresh zeroed</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	<span class="comment">// out memory.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	allocBits  *gcBits
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	gcmarkBits *gcBits
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	pinnerBits *gcBits <span class="comment">// bitmap for pinned objects; accessed atomically</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	<span class="comment">// sweep generation:</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	<span class="comment">// if sweepgen == h-&gt;sweepgen - 2, the span needs sweeping</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	<span class="comment">// if sweepgen == h-&gt;sweepgen - 1, the span is currently being swept</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	<span class="comment">// if sweepgen == h-&gt;sweepgen, the span is swept and ready to use</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	<span class="comment">// if sweepgen == h-&gt;sweepgen + 1, the span was cached before sweep began and is still cached, and needs sweeping</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	<span class="comment">// if sweepgen == h-&gt;sweepgen + 3, the span was swept and then cached and is still cached</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	<span class="comment">// h-&gt;sweepgen is incremented by 2 after every GC</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	sweepgen              uint32
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	divMul                uint32        <span class="comment">// for divide by elemsize</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	allocCount            uint16        <span class="comment">// number of allocated objects</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	spanclass             spanClass     <span class="comment">// size class and noscan (uint8)</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	state                 mSpanStateBox <span class="comment">// mSpanInUse etc; accessed atomically (get/set methods)</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	needzero              uint8         <span class="comment">// needs to be zeroed before allocation</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	isUserArenaChunk      bool          <span class="comment">// whether or not this span represents a user arena</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	allocCountBeforeCache uint16        <span class="comment">// a copy of allocCount that is stored just before this span is cached</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	elemsize              uintptr       <span class="comment">// computed from sizeclass or from npages</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	limit                 uintptr       <span class="comment">// end of data in span</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	speciallock           mutex         <span class="comment">// guards specials list and changes to pinnerBits</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	specials              *special      <span class="comment">// linked list of special records sorted by offset.</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	userArenaChunkFree    addrRange     <span class="comment">// interval for managing chunk allocation</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	largeType             *_type        <span class="comment">// malloc header for large objects.</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>func (s *mspan) base() uintptr {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	return s.startAddr
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>func (s *mspan) layout() (size, n, total uintptr) {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	total = s.npages &lt;&lt; _PageShift
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	size = s.elemsize
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	if size &gt; 0 {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		n = total / size
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	return
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span><span class="comment">// recordspan adds a newly allocated span to h.allspans.</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span><span class="comment">// This only happens the first time a span is allocated from</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span><span class="comment">// mheap.spanalloc (it is not called when a span is reused).</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span><span class="comment">// Write barriers are disallowed here because it can be called from</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span><span class="comment">// gcWork when allocating new workbufs. However, because it&#39;s an</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span><span class="comment">// indirect call from the fixalloc initializer, the compiler can&#39;t see</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span><span class="comment">// this.</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span><span class="comment">// The heap lock must be held.</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>func recordspan(vh unsafe.Pointer, p unsafe.Pointer) {
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	h := (*mheap)(vh)
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	s := (*mspan)(p)
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	assertLockHeld(&amp;h.lock)
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	if len(h.allspans) &gt;= cap(h.allspans) {
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		n := 64 * 1024 / goarch.PtrSize
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		if n &lt; cap(h.allspans)*3/2 {
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			n = cap(h.allspans) * 3 / 2
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		}
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		var new []*mspan
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		sp := (*slice)(unsafe.Pointer(&amp;new))
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		sp.array = sysAlloc(uintptr(n)*goarch.PtrSize, &amp;memstats.other_sys)
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		if sp.array == nil {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			throw(&#34;runtime: cannot allocate memory&#34;)
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		sp.len = len(h.allspans)
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		sp.cap = n
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		if len(h.allspans) &gt; 0 {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			copy(new, h.allspans)
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		}
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		oldAllspans := h.allspans
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		*(*notInHeapSlice)(unsafe.Pointer(&amp;h.allspans)) = *(*notInHeapSlice)(unsafe.Pointer(&amp;new))
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		if len(oldAllspans) != 0 {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>			sysFree(unsafe.Pointer(&amp;oldAllspans[0]), uintptr(cap(oldAllspans))*unsafe.Sizeof(oldAllspans[0]), &amp;memstats.other_sys)
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		}
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	h.allspans = h.allspans[:len(h.allspans)+1]
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	h.allspans[len(h.allspans)-1] = s
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// A spanClass represents the size class and noscan-ness of a span.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span><span class="comment">// Each size class has a noscan spanClass and a scan spanClass. The</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span><span class="comment">// noscan spanClass contains only noscan objects, which do not contain</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span><span class="comment">// pointers and thus do not need to be scanned by the garbage</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span><span class="comment">// collector.</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>type spanClass uint8
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>const (
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	numSpanClasses = _NumSizeClasses &lt;&lt; 1
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	tinySpanClass  = spanClass(tinySizeClass&lt;&lt;1 | 1)
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>func makeSpanClass(sizeclass uint8, noscan bool) spanClass {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	return spanClass(sizeclass&lt;&lt;1) | spanClass(bool2int(noscan))
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>func (sc spanClass) sizeclass() int8 {
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	return int8(sc &gt;&gt; 1)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>}
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>func (sc spanClass) noscan() bool {
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	return sc&amp;1 != 0
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>}
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span><span class="comment">// arenaIndex returns the index into mheap_.arenas of the arena</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span><span class="comment">// containing metadata for p. This index combines of an index into the</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span><span class="comment">// L1 map and an index into the L2 map and should be used as</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span><span class="comment">// mheap_.arenas[ai.l1()][ai.l2()].</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span><span class="comment">// If p is outside the range of valid heap addresses, either l1() or</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span><span class="comment">// l2() will be out of bounds.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span><span class="comment">// It is nosplit because it&#39;s called by spanOf and several other</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span><span class="comment">// nosplit functions.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>func arenaIndex(p uintptr) arenaIdx {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	return arenaIdx((p - arenaBaseOffset) / heapArenaBytes)
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span><span class="comment">// arenaBase returns the low address of the region covered by heap</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span><span class="comment">// arena i.</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>func arenaBase(i arenaIdx) uintptr {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	return uintptr(i)*heapArenaBytes + arenaBaseOffset
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>}
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>type arenaIdx uint
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span><span class="comment">// l1 returns the &#34;l1&#34; portion of an arenaIdx.</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span><span class="comment">// Marked nosplit because it&#39;s called by spanOf and other nosplit</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span><span class="comment">// functions.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>func (i arenaIdx) l1() uint {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	if arenaL1Bits == 0 {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		<span class="comment">// Let the compiler optimize this away if there&#39;s no</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		<span class="comment">// L1 map.</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		return 0
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	} else {
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		return uint(i) &gt;&gt; arenaL1Shift
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	}
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>}
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span><span class="comment">// l2 returns the &#34;l2&#34; portion of an arenaIdx.</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span><span class="comment">// Marked nosplit because it&#39;s called by spanOf and other nosplit funcs.</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span><span class="comment">// functions.</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>func (i arenaIdx) l2() uint {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	if arenaL1Bits == 0 {
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		return uint(i)
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	} else {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		return uint(i) &amp; (1&lt;&lt;arenaL2Bits - 1)
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>}
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span><span class="comment">// inheap reports whether b is a pointer into a (potentially dead) heap object.</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span><span class="comment">// It returns false for pointers into mSpanManual spans.</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span><span class="comment">// Non-preemptible because it is used by write barriers.</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>func inheap(b uintptr) bool {
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	return spanOfHeap(b) != nil
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span><span class="comment">// inHeapOrStack is a variant of inheap that returns true for pointers</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span><span class="comment">// into any allocated heap span.</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>func inHeapOrStack(b uintptr) bool {
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	s := spanOf(b)
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	if s == nil || b &lt; s.base() {
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		return false
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	}
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	switch s.state.get() {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	case mSpanInUse, mSpanManual:
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		return b &lt; s.limit
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	default:
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		return false
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>}
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span><span class="comment">// spanOf returns the span of p. If p does not point into the heap</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span><span class="comment">// arena or no span has ever contained p, spanOf returns nil.</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span><span class="comment">// If p does not point to allocated memory, this may return a non-nil</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span><span class="comment">// span that does *not* contain p. If this is a possibility, the</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span><span class="comment">// caller should either call spanOfHeap or check the span bounds</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span><span class="comment">// explicitly.</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span><span class="comment">// Must be nosplit because it has callers that are nosplit.</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>func spanOf(p uintptr) *mspan {
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	<span class="comment">// This function looks big, but we use a lot of constant</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	<span class="comment">// folding around arenaL1Bits to get it under the inlining</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	<span class="comment">// budget. Also, many of the checks here are safety checks</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	<span class="comment">// that Go needs to do anyway, so the generated code is quite</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	<span class="comment">// short.</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	ri := arenaIndex(p)
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	if arenaL1Bits == 0 {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		<span class="comment">// If there&#39;s no L1, then ri.l1() can&#39;t be out of bounds but ri.l2() can.</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		if ri.l2() &gt;= uint(len(mheap_.arenas[0])) {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			return nil
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		}
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	} else {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		<span class="comment">// If there&#39;s an L1, then ri.l1() can be out of bounds but ri.l2() can&#39;t.</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		if ri.l1() &gt;= uint(len(mheap_.arenas)) {
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>			return nil
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	}
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	l2 := mheap_.arenas[ri.l1()]
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	if arenaL1Bits != 0 &amp;&amp; l2 == nil { <span class="comment">// Should never happen if there&#39;s no L1.</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		return nil
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	}
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	ha := l2[ri.l2()]
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	if ha == nil {
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		return nil
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	return ha.spans[(p/pageSize)%pagesPerArena]
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>}
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span><span class="comment">// spanOfUnchecked is equivalent to spanOf, but the caller must ensure</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span><span class="comment">// that p points into an allocated heap arena.</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span><span class="comment">// Must be nosplit because it has callers that are nosplit.</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>func spanOfUnchecked(p uintptr) *mspan {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	ai := arenaIndex(p)
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	return mheap_.arenas[ai.l1()][ai.l2()].spans[(p/pageSize)%pagesPerArena]
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>}
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span><span class="comment">// spanOfHeap is like spanOf, but returns nil if p does not point to a</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span><span class="comment">// heap object.</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span><span class="comment">// Must be nosplit because it has callers that are nosplit.</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>func spanOfHeap(p uintptr) *mspan {
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	s := spanOf(p)
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	<span class="comment">// s is nil if it&#39;s never been allocated. Otherwise, we check</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	<span class="comment">// its state first because we don&#39;t trust this pointer, so we</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	<span class="comment">// have to synchronize with span initialization. Then, it&#39;s</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	<span class="comment">// still possible we picked up a stale span pointer, so we</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	<span class="comment">// have to check the span&#39;s bounds.</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	if s == nil || s.state.get() != mSpanInUse || p &lt; s.base() || p &gt;= s.limit {
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		return nil
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>	}
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	return s
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>}
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span><span class="comment">// pageIndexOf returns the arena, page index, and page mask for pointer p.</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span><span class="comment">// The caller must ensure p is in the heap.</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>func pageIndexOf(p uintptr) (arena *heapArena, pageIdx uintptr, pageMask uint8) {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	ai := arenaIndex(p)
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	arena = mheap_.arenas[ai.l1()][ai.l2()]
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	pageIdx = ((p / pageSize) / 8) % uintptr(len(arena.pageInUse))
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	pageMask = byte(1 &lt;&lt; ((p / pageSize) % 8))
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	return
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span><span class="comment">// Initialize the heap.</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>func (h *mheap) init() {
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	lockInit(&amp;h.lock, lockRankMheap)
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	lockInit(&amp;h.speciallock, lockRankMheapSpecial)
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	h.spanalloc.init(unsafe.Sizeof(mspan{}), recordspan, unsafe.Pointer(h), &amp;memstats.mspan_sys)
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	h.cachealloc.init(unsafe.Sizeof(mcache{}), nil, nil, &amp;memstats.mcache_sys)
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	h.specialfinalizeralloc.init(unsafe.Sizeof(specialfinalizer{}), nil, nil, &amp;memstats.other_sys)
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	h.specialprofilealloc.init(unsafe.Sizeof(specialprofile{}), nil, nil, &amp;memstats.other_sys)
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	h.specialReachableAlloc.init(unsafe.Sizeof(specialReachable{}), nil, nil, &amp;memstats.other_sys)
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	h.specialPinCounterAlloc.init(unsafe.Sizeof(specialPinCounter{}), nil, nil, &amp;memstats.other_sys)
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	h.arenaHintAlloc.init(unsafe.Sizeof(arenaHint{}), nil, nil, &amp;memstats.other_sys)
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t zero mspan allocations. Background sweeping can</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	<span class="comment">// inspect a span concurrently with allocating it, so it&#39;s</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	<span class="comment">// important that the span&#39;s sweepgen survive across freeing</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	<span class="comment">// and re-allocating a span to prevent background sweeping</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	<span class="comment">// from improperly cas&#39;ing it from 0.</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	<span class="comment">// This is safe because mspan contains no heap pointers.</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	h.spanalloc.zero = false
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	<span class="comment">// h-&gt;mapcache needs no init</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	for i := range h.central {
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		h.central[i].mcentral.init(spanClass(i))
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	}
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	h.pages.init(&amp;h.lock, &amp;memstats.gcMiscSys, false)
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>}
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span><span class="comment">// reclaim sweeps and reclaims at least npage pages into the heap.</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span><span class="comment">// It is called before allocating npage pages to keep growth in check.</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span><span class="comment">// reclaim implements the page-reclaimer half of the sweeper.</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span><span class="comment">// h.lock must NOT be held.</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>func (h *mheap) reclaim(npage uintptr) {
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>	<span class="comment">// TODO(austin): Half of the time spent freeing spans is in</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	<span class="comment">// locking/unlocking the heap (even with low contention). We</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	<span class="comment">// could make the slow path here several times faster by</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	<span class="comment">// batching heap frees.</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	<span class="comment">// Bail early if there&#39;s no more reclaim work.</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	if h.reclaimIndex.Load() &gt;= 1&lt;&lt;63 {
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		return
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	}
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	<span class="comment">// Disable preemption so the GC can&#39;t start while we&#39;re</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	<span class="comment">// sweeping, so we can read h.sweepArenas, and so</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	<span class="comment">// traceGCSweepStart/Done pair on the P.</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	trace := traceAcquire()
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		trace.GCSweepStart()
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	}
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	arenas := h.sweepArenas
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	locked := false
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	for npage &gt; 0 {
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		<span class="comment">// Pull from accumulated credit first.</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>		if credit := h.reclaimCredit.Load(); credit &gt; 0 {
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>			take := credit
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>			if take &gt; npage {
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>				<span class="comment">// Take only what we need.</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>				take = npage
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>			}
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>			if h.reclaimCredit.CompareAndSwap(credit, credit-take) {
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>				npage -= take
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>			}
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>			continue
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		}
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		<span class="comment">// Claim a chunk of work.</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>		idx := uintptr(h.reclaimIndex.Add(pagesPerReclaimerChunk) - pagesPerReclaimerChunk)
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		if idx/pagesPerArena &gt;= uintptr(len(arenas)) {
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>			<span class="comment">// Page reclaiming is done.</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>			h.reclaimIndex.Store(1 &lt;&lt; 63)
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>			break
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>		}
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		if !locked {
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>			<span class="comment">// Lock the heap for reclaimChunk.</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>			lock(&amp;h.lock)
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>			locked = true
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>		}
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		<span class="comment">// Scan this chunk.</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		nfound := h.reclaimChunk(arenas, idx, pagesPerReclaimerChunk)
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		if nfound &lt;= npage {
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>			npage -= nfound
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		} else {
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>			<span class="comment">// Put spare pages toward global credit.</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>			h.reclaimCredit.Add(nfound - npage)
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>			npage = 0
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		}
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	}
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	if locked {
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		unlock(&amp;h.lock)
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>	trace = traceAcquire()
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		trace.GCSweepDone()
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	}
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	releasem(mp)
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span><span class="comment">// reclaimChunk sweeps unmarked spans that start at page indexes [pageIdx, pageIdx+n).</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span><span class="comment">// It returns the number of pages returned to the heap.</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span><span class="comment">// h.lock must be held and the caller must be non-preemptible. Note: h.lock may be</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span><span class="comment">// temporarily unlocked and re-locked in order to do sweeping or if tracing is</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span><span class="comment">// enabled.</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>func (h *mheap) reclaimChunk(arenas []arenaIdx, pageIdx, n uintptr) uintptr {
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	<span class="comment">// The heap lock must be held because this accesses the</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	<span class="comment">// heapArena.spans arrays using potentially non-live pointers.</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	<span class="comment">// In particular, if a span were freed and merged concurrently</span>
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	<span class="comment">// with this probing heapArena.spans, it would be possible to</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	<span class="comment">// observe arbitrary, stale span pointers.</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	assertLockHeld(&amp;h.lock)
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	n0 := n
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	var nFreed uintptr
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	sl := sweep.active.begin()
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	if !sl.valid {
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>		return 0
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	}
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	for n &gt; 0 {
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		ai := arenas[pageIdx/pagesPerArena]
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		ha := h.arenas[ai.l1()][ai.l2()]
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>		<span class="comment">// Get a chunk of the bitmap to work on.</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>		arenaPage := uint(pageIdx % pagesPerArena)
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>		inUse := ha.pageInUse[arenaPage/8:]
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>		marked := ha.pageMarks[arenaPage/8:]
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>		if uintptr(len(inUse)) &gt; n/8 {
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>			inUse = inUse[:n/8]
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>			marked = marked[:n/8]
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>		}
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>		<span class="comment">// Scan this bitmap chunk for spans that are in-use</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>		<span class="comment">// but have no marked objects on them.</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>		for i := range inUse {
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>			inUseUnmarked := atomic.Load8(&amp;inUse[i]) &amp;^ marked[i]
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			if inUseUnmarked == 0 {
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>				continue
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			}
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>			for j := uint(0); j &lt; 8; j++ {
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>				if inUseUnmarked&amp;(1&lt;&lt;j) != 0 {
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>					s := ha.spans[arenaPage+uint(i)*8+j]
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>					if s, ok := sl.tryAcquire(s); ok {
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>						npages := s.npages
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>						unlock(&amp;h.lock)
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>						if s.sweep(false) {
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>							nFreed += npages
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>						}
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>						lock(&amp;h.lock)
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>						<span class="comment">// Reload inUse. It&#39;s possible nearby</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>						<span class="comment">// spans were freed when we dropped the</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>						<span class="comment">// lock and we don&#39;t want to get stale</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>						<span class="comment">// pointers from the spans array.</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>						inUseUnmarked = atomic.Load8(&amp;inUse[i]) &amp;^ marked[i]
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>					}
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>				}
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>			}
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>		}
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>		<span class="comment">// Advance.</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>		pageIdx += uintptr(len(inUse) * 8)
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		n -= uintptr(len(inUse) * 8)
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>	}
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>	sweep.active.end(sl)
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>	trace := traceAcquire()
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>		unlock(&amp;h.lock)
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		<span class="comment">// Account for pages scanned but not reclaimed.</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>		trace.GCSweepSpan((n0 - nFreed) * pageSize)
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>		lock(&amp;h.lock)
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	}
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	assertLockHeld(&amp;h.lock) <span class="comment">// Must be locked on return.</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	return nFreed
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>}
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span><span class="comment">// spanAllocType represents the type of allocation to make, or</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span><span class="comment">// the type of allocation to be freed.</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>type spanAllocType uint8
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>const (
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>	spanAllocHeap          spanAllocType = iota <span class="comment">// heap span</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>	spanAllocStack                              <span class="comment">// stack span</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	spanAllocPtrScalarBits                      <span class="comment">// unrolled GC prog bitmap span</span>
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	spanAllocWorkBuf                            <span class="comment">// work buf span</span>
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>)
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>
<span id="L942" class="ln">   942&nbsp;&nbsp;</span><span class="comment">// manual returns true if the span allocation is manually managed.</span>
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>func (s spanAllocType) manual() bool {
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>	return s != spanAllocHeap
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>}
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span><span class="comment">// alloc allocates a new span of npage pages from the GC&#39;d heap.</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L949" class="ln">   949&nbsp;&nbsp;</span><span class="comment">// spanclass indicates the span&#39;s size class and scannability.</span>
<span id="L950" class="ln">   950&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span><span class="comment">// Returns a span that has been fully initialized. span.needzero indicates</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span><span class="comment">// whether the span has been zeroed. Note that it may not be.</span>
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>func (h *mheap) alloc(npages uintptr, spanclass spanClass) *mspan {
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t do any operations that lock the heap on the G stack.</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>	<span class="comment">// It might trigger stack growth, and the stack growth code needs</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>	<span class="comment">// to be able to allocate heap.</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	var s *mspan
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		<span class="comment">// To prevent excessive heap growth, before allocating n pages</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>		<span class="comment">// we need to sweep and reclaim at least n pages.</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>		if !isSweepDone() {
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>			h.reclaim(npages)
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>		}
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>		s = h.allocSpan(npages, spanAllocHeap, spanclass)
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	})
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>	return s
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>}
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span><span class="comment">// allocManual allocates a manually-managed span of npage pages.</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span><span class="comment">// allocManual returns nil if allocation fails.</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L972" class="ln">   972&nbsp;&nbsp;</span><span class="comment">// allocManual adds the bytes used to *stat, which should be a</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span><span class="comment">// memstats in-use field. Unlike allocations in the GC&#39;d heap, the</span>
<span id="L974" class="ln">   974&nbsp;&nbsp;</span><span class="comment">// allocation does *not* count toward heapInUse.</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span><span class="comment">// The memory backing the returned span may not be zeroed if</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span><span class="comment">// span.needzero is set.</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span><span class="comment">// allocManual must be called on the system stack because it may</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span><span class="comment">// acquire the heap lock via allocSpan. See mheap for details.</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span><span class="comment">// If new code is written to call allocManual, do NOT use an</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span><span class="comment">// existing spanAllocType value and instead declare a new one.</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L985" class="ln">   985&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>func (h *mheap) allocManual(npages uintptr, typ spanAllocType) *mspan {
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	if !typ.manual() {
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>		throw(&#34;manual span allocation called with non-manually-managed type&#34;)
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	}
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	return h.allocSpan(npages, typ, 0)
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>}
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span><span class="comment">// setSpans modifies the span map so [spanOf(base), spanOf(base+npage*pageSize))</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span><span class="comment">// is s.</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>func (h *mheap) setSpans(base, npage uintptr, s *mspan) {
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	p := base / pageSize
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	ai := arenaIndex(base)
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	ha := h.arenas[ai.l1()][ai.l2()]
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	for n := uintptr(0); n &lt; npage; n++ {
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>		i := (p + n) % pagesPerArena
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>		if i == 0 {
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>			ai = arenaIndex(base + n*pageSize)
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>			ha = h.arenas[ai.l1()][ai.l2()]
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>		}
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>		ha.spans[i] = s
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	}
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>}
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span><span class="comment">// allocNeedsZero checks if the region of address space [base, base+npage*pageSize),</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span><span class="comment">// assumed to be allocated, needs to be zeroed, updating heap arena metadata for</span>
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span><span class="comment">// future allocations.</span>
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span><span class="comment">// This must be called each time pages are allocated from the heap, even if the page</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span><span class="comment">// allocator can otherwise prove the memory it&#39;s allocating is already zero because</span>
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span><span class="comment">// they&#39;re fresh from the operating system. It updates heapArena metadata that is</span>
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span><span class="comment">// critical for future page allocations.</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span><span class="comment">// There are no locking constraints on this method.</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>func (h *mheap) allocNeedsZero(base, npage uintptr) (needZero bool) {
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>	for npage &gt; 0 {
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		ai := arenaIndex(base)
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>		ha := h.arenas[ai.l1()][ai.l2()]
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>		zeroedBase := atomic.Loaduintptr(&amp;ha.zeroedBase)
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		arenaBase := base % heapArenaBytes
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>		if arenaBase &lt; zeroedBase {
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>			<span class="comment">// We extended into the non-zeroed part of the</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>			<span class="comment">// arena, so this region needs to be zeroed before use.</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>			<span class="comment">// zeroedBase is monotonically increasing, so if we see this now then</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>			<span class="comment">// we can be sure we need to zero this memory region.</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>			<span class="comment">// We still need to update zeroedBase for this arena, and</span>
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>			<span class="comment">// potentially more arenas.</span>
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>			needZero = true
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>		}
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>		<span class="comment">// We may observe arenaBase &gt; zeroedBase if we&#39;re racing with one or more</span>
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>		<span class="comment">// allocations which are acquiring memory directly before us in the address</span>
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>		<span class="comment">// space. But, because we know no one else is acquiring *this* memory, it&#39;s</span>
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>		<span class="comment">// still safe to not zero.</span>
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		<span class="comment">// Compute how far into the arena we extend into, capped</span>
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>		<span class="comment">// at heapArenaBytes.</span>
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>		arenaLimit := arenaBase + npage*pageSize
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>		if arenaLimit &gt; heapArenaBytes {
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>			arenaLimit = heapArenaBytes
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>		}
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>		<span class="comment">// Increase ha.zeroedBase so it&#39;s &gt;= arenaLimit.</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>		<span class="comment">// We may be racing with other updates.</span>
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>		for arenaLimit &gt; zeroedBase {
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>			if atomic.Casuintptr(&amp;ha.zeroedBase, zeroedBase, arenaLimit) {
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>				break
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>			}
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>			zeroedBase = atomic.Loaduintptr(&amp;ha.zeroedBase)
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>			<span class="comment">// Double check basic conditions of zeroedBase.</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>			if zeroedBase &lt;= arenaLimit &amp;&amp; zeroedBase &gt; arenaBase {
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>				<span class="comment">// The zeroedBase moved into the space we were trying to</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>				<span class="comment">// claim. That&#39;s very bad, and indicates someone allocated</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>				<span class="comment">// the same region we did.</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>				throw(&#34;potentially overlapping in-use allocations detected&#34;)
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>			}
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>		}
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>		<span class="comment">// Move base forward and subtract from npage to move into</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>		<span class="comment">// the next arena, or finish.</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>		base += arenaLimit - arenaBase
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>		npage -= (arenaLimit - arenaBase) / pageSize
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>	}
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>	return
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>}
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span><span class="comment">// tryAllocMSpan attempts to allocate an mspan object from</span>
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span><span class="comment">// the P-local cache, but may fail.</span>
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span><span class="comment">// h.lock need not be held.</span>
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span><span class="comment">// This caller must ensure that its P won&#39;t change underneath</span>
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span><span class="comment">// it during this function. Currently to ensure that we enforce</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span><span class="comment">// that the function is run on the system stack, because that&#39;s</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span><span class="comment">// the only place it is used now. In the future, this requirement</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span><span class="comment">// may be relaxed if its use is necessary elsewhere.</span>
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>func (h *mheap) tryAllocMSpan() *mspan {
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>	pp := getg().m.p.ptr()
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>	<span class="comment">// If we don&#39;t have a p or the cache is empty, we can&#39;t do</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	<span class="comment">// anything here.</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>	if pp == nil || pp.mspancache.len == 0 {
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>		return nil
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>	}
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>	<span class="comment">// Pull off the last entry in the cache.</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>	s := pp.mspancache.buf[pp.mspancache.len-1]
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>	pp.mspancache.len--
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>	return s
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>}
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span><span class="comment">// allocMSpanLocked allocates an mspan object.</span>
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span><span class="comment">// h.lock must be held.</span>
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span><span class="comment">// allocMSpanLocked must be called on the system stack because</span>
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span><span class="comment">// its caller holds the heap lock. See mheap for details.</span>
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span><span class="comment">// Running on the system stack also ensures that we won&#39;t</span>
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span><span class="comment">// switch Ps during this function. See tryAllocMSpan for details.</span>
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>func (h *mheap) allocMSpanLocked() *mspan {
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>	assertLockHeld(&amp;h.lock)
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>	pp := getg().m.p.ptr()
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>	if pp == nil {
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>		<span class="comment">// We don&#39;t have a p so just do the normal thing.</span>
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>		return (*mspan)(h.spanalloc.alloc())
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>	}
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>	<span class="comment">// Refill the cache if necessary.</span>
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>	if pp.mspancache.len == 0 {
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>		const refillCount = len(pp.mspancache.buf) / 2
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>		for i := 0; i &lt; refillCount; i++ {
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>			pp.mspancache.buf[i] = (*mspan)(h.spanalloc.alloc())
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>		}
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>		pp.mspancache.len = refillCount
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	}
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	<span class="comment">// Pull off the last entry in the cache.</span>
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>	s := pp.mspancache.buf[pp.mspancache.len-1]
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>	pp.mspancache.len--
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>	return s
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>}
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span><span class="comment">// freeMSpanLocked free an mspan object.</span>
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span><span class="comment">// h.lock must be held.</span>
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span><span class="comment">// freeMSpanLocked must be called on the system stack because</span>
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span><span class="comment">// its caller holds the heap lock. See mheap for details.</span>
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span><span class="comment">// Running on the system stack also ensures that we won&#39;t</span>
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span><span class="comment">// switch Ps during this function. See tryAllocMSpan for details.</span>
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>func (h *mheap) freeMSpanLocked(s *mspan) {
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>	assertLockHeld(&amp;h.lock)
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	pp := getg().m.p.ptr()
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	<span class="comment">// First try to free the mspan directly to the cache.</span>
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>	if pp != nil &amp;&amp; pp.mspancache.len &lt; len(pp.mspancache.buf) {
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>		pp.mspancache.buf[pp.mspancache.len] = s
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>		pp.mspancache.len++
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>		return
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>	}
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>	<span class="comment">// Failing that (or if we don&#39;t have a p), just free it to</span>
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>	<span class="comment">// the heap.</span>
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>	h.spanalloc.free(unsafe.Pointer(s))
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>}
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span><span class="comment">// allocSpan allocates an mspan which owns npages worth of memory.</span>
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span><span class="comment">// If typ.manual() == false, allocSpan allocates a heap span of class spanclass</span>
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span><span class="comment">// and updates heap accounting. If manual == true, allocSpan allocates a</span>
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span><span class="comment">// manually-managed span (spanclass is ignored), and the caller is</span>
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span><span class="comment">// responsible for any accounting related to its use of the span. Either</span>
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span><span class="comment">// way, allocSpan will atomically add the bytes in the newly allocated</span>
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span><span class="comment">// span to *sysStat.</span>
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span><span class="comment">// The returned span is fully initialized.</span>
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span><span class="comment">// h.lock must not be held.</span>
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span><span class="comment">// allocSpan must be called on the system stack both because it acquires</span>
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span><span class="comment">// the heap lock and because it must block GC transitions.</span>
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>func (h *mheap) allocSpan(npages uintptr, typ spanAllocType, spanclass spanClass) (s *mspan) {
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>	<span class="comment">// Function-global state.</span>
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>	gp := getg()
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>	base, scav := uintptr(0), uintptr(0)
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>	growth := uintptr(0)
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>	<span class="comment">// On some platforms we need to provide physical page aligned stack</span>
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>	<span class="comment">// allocations. Where the page size is less than the physical page</span>
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>	<span class="comment">// size, we already manage to do this by default.</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>	needPhysPageAlign := physPageAlignedStacks &amp;&amp; typ == spanAllocStack &amp;&amp; pageSize &lt; physPageSize
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>	<span class="comment">// If the allocation is small enough, try the page cache!</span>
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>	<span class="comment">// The page cache does not support aligned allocations, so we cannot use</span>
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>	<span class="comment">// it if we need to provide a physical page aligned stack allocation.</span>
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>	pp := gp.m.p.ptr()
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>	if !needPhysPageAlign &amp;&amp; pp != nil &amp;&amp; npages &lt; pageCachePages/4 {
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>		c := &amp;pp.pcache
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>		<span class="comment">// If the cache is empty, refill it.</span>
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>		if c.empty() {
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>			lock(&amp;h.lock)
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>			*c = h.pages.allocToCache()
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>			unlock(&amp;h.lock)
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>		}
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>		<span class="comment">// Try to allocate from the cache.</span>
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>		base, scav = c.alloc(npages)
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>		if base != 0 {
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>			s = h.tryAllocMSpan()
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>			if s != nil {
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>				goto HaveSpan
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>			}
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>			<span class="comment">// We have a base but no mspan, so we need</span>
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>			<span class="comment">// to lock the heap.</span>
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>		}
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>	}
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>	<span class="comment">// For one reason or another, we couldn&#39;t get the</span>
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>	<span class="comment">// whole job done without the heap lock.</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>	lock(&amp;h.lock)
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>	if needPhysPageAlign {
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>		<span class="comment">// Overallocate by a physical page to allow for later alignment.</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>		extraPages := physPageSize / pageSize
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>		<span class="comment">// Find a big enough region first, but then only allocate the</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>		<span class="comment">// aligned portion. We can&#39;t just allocate and then free the</span>
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>		<span class="comment">// edges because we need to account for scavenged memory, and</span>
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>		<span class="comment">// that&#39;s difficult with alloc.</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>		<span class="comment">// Note that we skip updates to searchAddr here. It&#39;s OK if</span>
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>		<span class="comment">// it&#39;s stale and higher than normal; it&#39;ll operate correctly,</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>		<span class="comment">// just come with a performance cost.</span>
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>		base, _ = h.pages.find(npages + extraPages)
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>		if base == 0 {
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>			var ok bool
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>			growth, ok = h.grow(npages + extraPages)
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>			if !ok {
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>				unlock(&amp;h.lock)
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>				return nil
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>			}
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>			base, _ = h.pages.find(npages + extraPages)
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>			if base == 0 {
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>				throw(&#34;grew heap, but no adequate free space found&#34;)
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>			}
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>		}
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>		base = alignUp(base, physPageSize)
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>		scav = h.pages.allocRange(base, npages)
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>	}
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>	if base == 0 {
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>		<span class="comment">// Try to acquire a base address.</span>
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>		base, scav = h.pages.alloc(npages)
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>		if base == 0 {
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>			var ok bool
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>			growth, ok = h.grow(npages)
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>			if !ok {
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>				unlock(&amp;h.lock)
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>				return nil
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>			}
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>			base, scav = h.pages.alloc(npages)
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>			if base == 0 {
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>				throw(&#34;grew heap, but no adequate free space found&#34;)
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>			}
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>		}
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>	}
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>	if s == nil {
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>		<span class="comment">// We failed to get an mspan earlier, so grab</span>
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>		<span class="comment">// one now that we have the heap lock.</span>
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>		s = h.allocMSpanLocked()
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>	}
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>	unlock(&amp;h.lock)
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>HaveSpan:
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>	<span class="comment">// Decide if we need to scavenge in response to what we just allocated.</span>
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>	<span class="comment">// Specifically, we track the maximum amount of memory to scavenge of all</span>
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>	<span class="comment">// the alternatives below, assuming that the maximum satisfies *all*</span>
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>	<span class="comment">// conditions we check (e.g. if we need to scavenge X to satisfy the</span>
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>	<span class="comment">// memory limit and Y to satisfy heap-growth scavenging, and Y &gt; X, then</span>
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	<span class="comment">// it&#39;s fine to pick Y, because the memory limit is still satisfied).</span>
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s fine to do this after allocating because we expect any scavenged</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	<span class="comment">// pages not to get touched until we return. Simultaneously, it&#39;s important</span>
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>	<span class="comment">// to do this before calling sysUsed because that may commit address space.</span>
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>	bytesToScavenge := uintptr(0)
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>	forceScavenge := false
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>	if limit := gcController.memoryLimit.Load(); !gcCPULimiter.limiting() {
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>		<span class="comment">// Assist with scavenging to maintain the memory limit by the amount</span>
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>		<span class="comment">// that we expect to page in.</span>
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>		inuse := gcController.mappedReady.Load()
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>		<span class="comment">// Be careful about overflow, especially with uintptrs. Even on 32-bit platforms</span>
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>		<span class="comment">// someone can set a really big memory limit that isn&#39;t maxInt64.</span>
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>		if uint64(scav)+inuse &gt; uint64(limit) {
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>			bytesToScavenge = uintptr(uint64(scav) + inuse - uint64(limit))
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>			forceScavenge = true
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>		}
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>	}
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>	if goal := scavenge.gcPercentGoal.Load(); goal != ^uint64(0) &amp;&amp; growth &gt; 0 {
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>		<span class="comment">// We just caused a heap growth, so scavenge down what will soon be used.</span>
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>		<span class="comment">// By scavenging inline we deal with the failure to allocate out of</span>
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>		<span class="comment">// memory fragments by scavenging the memory fragments that are least</span>
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>		<span class="comment">// likely to be re-used.</span>
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>		<span class="comment">// Only bother with this because we&#39;re not using a memory limit. We don&#39;t</span>
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>		<span class="comment">// care about heap growths as long as we&#39;re under the memory limit, and the</span>
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>		<span class="comment">// previous check for scaving already handles that.</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>		if retained := heapRetained(); retained+uint64(growth) &gt; goal {
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>			<span class="comment">// The scavenging algorithm requires the heap lock to be dropped so it</span>
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>			<span class="comment">// can acquire it only sparingly. This is a potentially expensive operation</span>
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>			<span class="comment">// so it frees up other goroutines to allocate in the meanwhile. In fact,</span>
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>			<span class="comment">// they can make use of the growth we just created.</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>			todo := growth
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>			if overage := uintptr(retained + uint64(growth) - goal); todo &gt; overage {
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>				todo = overage
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>			}
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>			if todo &gt; bytesToScavenge {
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>				bytesToScavenge = todo
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>			}
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>		}
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>	}
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>	<span class="comment">// There are a few very limited circumstances where we won&#39;t have a P here.</span>
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s OK to simply skip scavenging in these cases. Something else will notice</span>
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>	<span class="comment">// and pick up the tab.</span>
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>	var now int64
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>	if pp != nil &amp;&amp; bytesToScavenge &gt; 0 {
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>		<span class="comment">// Measure how long we spent scavenging and add that measurement to the assist</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>		<span class="comment">// time so we can track it for the GC CPU limiter.</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>		<span class="comment">// Limiter event tracking might be disabled if we end up here</span>
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>		<span class="comment">// while on a mark worker.</span>
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>		start := nanotime()
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>		track := pp.limiterEvent.start(limiterEventScavengeAssist, start)
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>		<span class="comment">// Scavenge, but back out if the limiter turns on.</span>
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>		released := h.pages.scavenge(bytesToScavenge, func() bool {
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>			return gcCPULimiter.limiting()
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>		}, forceScavenge)
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>		mheap_.pages.scav.releasedEager.Add(released)
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>		<span class="comment">// Finish up accounting.</span>
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>		now = nanotime()
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>		if track {
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>			pp.limiterEvent.stop(limiterEventScavengeAssist, now)
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>		}
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>		scavenge.assistTime.Add(now - start)
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>	}
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>	<span class="comment">// Initialize the span.</span>
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>	h.initSpan(s, typ, spanclass, base, npages)
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>	<span class="comment">// Commit and account for any scavenged memory that the span now owns.</span>
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>	nbytes := npages * pageSize
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>	if scav != 0 {
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>		<span class="comment">// sysUsed all the pages that are actually available</span>
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>		<span class="comment">// in the span since some of them might be scavenged.</span>
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>		sysUsed(unsafe.Pointer(base), nbytes, scav)
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>		gcController.heapReleased.add(-int64(scav))
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>	}
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>	<span class="comment">// Update stats.</span>
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>	gcController.heapFree.add(-int64(nbytes - scav))
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>	if typ == spanAllocHeap {
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>		gcController.heapInUse.add(int64(nbytes))
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>	}
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>	<span class="comment">// Update consistent stats.</span>
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>	stats := memstats.heapStats.acquire()
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>	atomic.Xaddint64(&amp;stats.committed, int64(scav))
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>	atomic.Xaddint64(&amp;stats.released, -int64(scav))
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>	switch typ {
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>	case spanAllocHeap:
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>		atomic.Xaddint64(&amp;stats.inHeap, int64(nbytes))
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>	case spanAllocStack:
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>		atomic.Xaddint64(&amp;stats.inStacks, int64(nbytes))
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>	case spanAllocPtrScalarBits:
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>		atomic.Xaddint64(&amp;stats.inPtrScalarBits, int64(nbytes))
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>	case spanAllocWorkBuf:
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>		atomic.Xaddint64(&amp;stats.inWorkBufs, int64(nbytes))
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>	}
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>	memstats.heapStats.release()
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>	pageTraceAlloc(pp, now, base, npages)
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>	return s
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>}
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span><span class="comment">// initSpan initializes a blank span s which will represent the range</span>
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span><span class="comment">// [base, base+npages*pageSize). typ is the type of span being allocated.</span>
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>func (h *mheap) initSpan(s *mspan, typ spanAllocType, spanclass spanClass, base, npages uintptr) {
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>	<span class="comment">// At this point, both s != nil and base != 0, and the heap</span>
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>	<span class="comment">// lock is no longer held. Initialize the span.</span>
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>	s.init(base, npages)
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>	if h.allocNeedsZero(base, npages) {
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>		s.needzero = 1
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>	}
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>	nbytes := npages * pageSize
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>	if typ.manual() {
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>		s.manualFreeList = 0
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>		s.nelems = 0
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>		s.limit = s.base() + s.npages*pageSize
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>		s.state.set(mSpanManual)
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>	} else {
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>		<span class="comment">// We must set span properties before the span is published anywhere</span>
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>		<span class="comment">// since we&#39;re not holding the heap lock.</span>
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>		s.spanclass = spanclass
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>		if sizeclass := spanclass.sizeclass(); sizeclass == 0 {
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>			s.elemsize = nbytes
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>			s.nelems = 1
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>			s.divMul = 0
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>		} else {
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>			s.elemsize = uintptr(class_to_size[sizeclass])
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>			if goexperiment.AllocHeaders &amp;&amp; !s.spanclass.noscan() &amp;&amp; heapBitsInSpan(s.elemsize) {
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>				<span class="comment">// In the allocheaders experiment, reserve space for the pointer/scan bitmap at the end.</span>
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>				s.nelems = uint16((nbytes - (nbytes / goarch.PtrSize / 8)) / s.elemsize)
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>			} else {
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>				s.nelems = uint16(nbytes / s.elemsize)
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>			}
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>			s.divMul = class_to_divmagic[sizeclass]
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>		}
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>		<span class="comment">// Initialize mark and allocation structures.</span>
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>		s.freeindex = 0
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>		s.freeIndexForScan = 0
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>		s.allocCache = ^uint64(0) <span class="comment">// all 1s indicating all free.</span>
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>		s.gcmarkBits = newMarkBits(uintptr(s.nelems))
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>		s.allocBits = newAllocBits(uintptr(s.nelems))
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>		<span class="comment">// It&#39;s safe to access h.sweepgen without the heap lock because it&#39;s</span>
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>		<span class="comment">// only ever updated with the world stopped and we run on the</span>
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>		<span class="comment">// systemstack which blocks a STW transition.</span>
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>		atomic.Store(&amp;s.sweepgen, h.sweepgen)
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>		<span class="comment">// Now that the span is filled in, set its state. This</span>
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>		<span class="comment">// is a publication barrier for the other fields in</span>
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>		<span class="comment">// the span. While valid pointers into this span</span>
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>		<span class="comment">// should never be visible until the span is returned,</span>
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>		<span class="comment">// if the garbage collector finds an invalid pointer,</span>
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>		<span class="comment">// access to the span may race with initialization of</span>
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>		<span class="comment">// the span. We resolve this race by atomically</span>
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>		<span class="comment">// setting the state after the span is fully</span>
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>		<span class="comment">// initialized, and atomically checking the state in</span>
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span>		<span class="comment">// any situation where a pointer is suspect.</span>
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span>		s.state.set(mSpanInUse)
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>	}
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>	<span class="comment">// Publish the span in various locations.</span>
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span>
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span>	<span class="comment">// This is safe to call without the lock held because the slots</span>
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span>	<span class="comment">// related to this span will only ever be read or modified by</span>
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span>	<span class="comment">// this thread until pointers into the span are published (and</span>
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span>	<span class="comment">// we execute a publication barrier at the end of this function</span>
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span>	<span class="comment">// before that happens) or pageInUse is updated.</span>
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>	h.setSpans(s.base(), npages, s)
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>	if !typ.manual() {
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>		<span class="comment">// Mark in-use span in arena page bitmap.</span>
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>		<span class="comment">// This publishes the span to the page sweeper, so</span>
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>		<span class="comment">// it&#39;s imperative that the span be completely initialized</span>
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>		<span class="comment">// prior to this line.</span>
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>		arena, pageIdx, pageMask := pageIndexOf(s.base())
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span>		atomic.Or8(&amp;arena.pageInUse[pageIdx], pageMask)
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span>
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span>		<span class="comment">// Update related page sweeper stats.</span>
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span>		h.pagesInUse.Add(npages)
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span>	}
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span>
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span>	<span class="comment">// Make sure the newly allocated span will be observed</span>
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span>	<span class="comment">// by the GC before pointers into the span are published.</span>
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span>	publicationBarrier()
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span>}
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span>
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span><span class="comment">// Try to add at least npage pages of memory to the heap,</span>
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span><span class="comment">// returning how much the heap grew by and whether it worked.</span>
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span><span class="comment">// h.lock must be held.</span>
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span>func (h *mheap) grow(npage uintptr) (uintptr, bool) {
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span>	assertLockHeld(&amp;h.lock)
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span>
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span>	<span class="comment">// We must grow the heap in whole palloc chunks.</span>
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span>	<span class="comment">// We call sysMap below but note that because we</span>
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span>	<span class="comment">// round up to pallocChunkPages which is on the order</span>
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span>	<span class="comment">// of MiB (generally &gt;= to the huge page size) we</span>
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span>	<span class="comment">// won&#39;t be calling it too much.</span>
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span>	ask := alignUp(npage, pallocChunkPages) * pageSize
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span>
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span>	totalGrowth := uintptr(0)
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span>	<span class="comment">// This may overflow because ask could be very large</span>
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span>	<span class="comment">// and is otherwise unrelated to h.curArena.base.</span>
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span>	end := h.curArena.base + ask
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span>	nBase := alignUp(end, physPageSize)
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span>	if nBase &gt; h.curArena.end || <span class="comment">/* overflow */</span> end &lt; h.curArena.base {
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span>		<span class="comment">// Not enough room in the current arena. Allocate more</span>
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span>		<span class="comment">// arena space. This may not be contiguous with the</span>
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span>		<span class="comment">// current arena, so we have to request the full ask.</span>
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span>		av, asize := h.sysAlloc(ask, &amp;h.arenaHints, true)
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span>		if av == nil {
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span>			inUse := gcController.heapFree.load() + gcController.heapReleased.load() + gcController.heapInUse.load()
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span>			print(&#34;runtime: out of memory: cannot allocate &#34;, ask, &#34;-byte block (&#34;, inUse, &#34; in use)\n&#34;)
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span>			return 0, false
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span>		}
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span>
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span>		if uintptr(av) == h.curArena.end {
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span>			<span class="comment">// The new space is contiguous with the old</span>
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span>			<span class="comment">// space, so just extend the current space.</span>
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span>			h.curArena.end = uintptr(av) + asize
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span>		} else {
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span>			<span class="comment">// The new space is discontiguous. Track what</span>
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span>			<span class="comment">// remains of the current space and switch to</span>
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span>			<span class="comment">// the new space. This should be rare.</span>
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span>			if size := h.curArena.end - h.curArena.base; size != 0 {
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span>				<span class="comment">// Transition this space from Reserved to Prepared and mark it</span>
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span>				<span class="comment">// as released since we&#39;ll be able to start using it after updating</span>
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span>				<span class="comment">// the page allocator and releasing the lock at any time.</span>
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span>				sysMap(unsafe.Pointer(h.curArena.base), size, &amp;gcController.heapReleased)
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span>				<span class="comment">// Update stats.</span>
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span>				stats := memstats.heapStats.acquire()
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span>				atomic.Xaddint64(&amp;stats.released, int64(size))
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span>				memstats.heapStats.release()
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span>				<span class="comment">// Update the page allocator&#39;s structures to make this</span>
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span>				<span class="comment">// space ready for allocation.</span>
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span>				h.pages.grow(h.curArena.base, size)
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span>				totalGrowth += size
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span>			}
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span>			<span class="comment">// Switch to the new space.</span>
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span>			h.curArena.base = uintptr(av)
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>			h.curArena.end = uintptr(av) + asize
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>		}
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span>
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span>		<span class="comment">// Recalculate nBase.</span>
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span>		<span class="comment">// We know this won&#39;t overflow, because sysAlloc returned</span>
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span>		<span class="comment">// a valid region starting at h.curArena.base which is at</span>
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span>		<span class="comment">// least ask bytes in size.</span>
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span>		nBase = alignUp(h.curArena.base+ask, physPageSize)
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span>	}
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span>
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span>	<span class="comment">// Grow into the current arena.</span>
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span>	v := h.curArena.base
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span>	h.curArena.base = nBase
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span>
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span>	<span class="comment">// Transition the space we&#39;re going to use from Reserved to Prepared.</span>
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>	<span class="comment">// The allocation is always aligned to the heap arena</span>
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>	<span class="comment">// size which is always &gt; physPageSize, so its safe to</span>
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>	<span class="comment">// just add directly to heapReleased.</span>
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span>	sysMap(unsafe.Pointer(v), nBase-v, &amp;gcController.heapReleased)
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span>
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span>	<span class="comment">// The memory just allocated counts as both released</span>
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span>	<span class="comment">// and idle, even though it&#39;s not yet backed by spans.</span>
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>	stats := memstats.heapStats.acquire()
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>	atomic.Xaddint64(&amp;stats.released, int64(nBase-v))
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>	memstats.heapStats.release()
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span>
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span>	<span class="comment">// Update the page allocator&#39;s structures to make this</span>
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span>	<span class="comment">// space ready for allocation.</span>
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span>	h.pages.grow(v, nBase-v)
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>	totalGrowth += nBase - v
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>	return totalGrowth, true
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>}
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span><span class="comment">// Free the span back into the heap.</span>
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span>func (h *mheap) freeSpan(s *mspan) {
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span>		pageTraceFree(getg().m.p.ptr(), 0, s.base(), s.npages)
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span>
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span>		lock(&amp;h.lock)
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span>		if msanenabled {
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span>			<span class="comment">// Tell msan that this entire span is no longer in use.</span>
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span>			base := unsafe.Pointer(s.base())
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span>			bytes := s.npages &lt;&lt; _PageShift
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span>			msanfree(base, bytes)
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span>		}
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span>		if asanenabled {
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span>			<span class="comment">// Tell asan that this entire span is no longer in use.</span>
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span>			base := unsafe.Pointer(s.base())
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span>			bytes := s.npages &lt;&lt; _PageShift
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span>			asanpoison(base, bytes)
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span>		}
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span>		h.freeSpanLocked(s, spanAllocHeap)
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span>		unlock(&amp;h.lock)
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span>	})
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span>}
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span>
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span><span class="comment">// freeManual frees a manually-managed span returned by allocManual.</span>
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span><span class="comment">// typ must be the same as the spanAllocType passed to the allocManual that</span>
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span><span class="comment">// allocated s.</span>
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span><span class="comment">// This must only be called when gcphase == _GCoff. See mSpanState for</span>
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span><span class="comment">// an explanation.</span>
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span><span class="comment">// freeManual must be called on the system stack because it acquires</span>
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span><span class="comment">// the heap lock. See mheap for details.</span>
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>func (h *mheap) freeManual(s *mspan, typ spanAllocType) {
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span>	pageTraceFree(getg().m.p.ptr(), 0, s.base(), s.npages)
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span>
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span>	s.needzero = 1
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span>	lock(&amp;h.lock)
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span>	h.freeSpanLocked(s, typ)
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>	unlock(&amp;h.lock)
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>}
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>func (h *mheap) freeSpanLocked(s *mspan, typ spanAllocType) {
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>	assertLockHeld(&amp;h.lock)
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>	switch s.state.get() {
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span>	case mSpanManual:
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span>		if s.allocCount != 0 {
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span>			throw(&#34;mheap.freeSpanLocked - invalid stack free&#34;)
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span>		}
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span>	case mSpanInUse:
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>		if s.isUserArenaChunk {
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>			throw(&#34;mheap.freeSpanLocked - invalid free of user arena chunk&#34;)
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span>		}
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span>		if s.allocCount != 0 || s.sweepgen != h.sweepgen {
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span>			print(&#34;mheap.freeSpanLocked - span &#34;, s, &#34; ptr &#34;, hex(s.base()), &#34; allocCount &#34;, s.allocCount, &#34; sweepgen &#34;, s.sweepgen, &#34;/&#34;, h.sweepgen, &#34;\n&#34;)
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span>			throw(&#34;mheap.freeSpanLocked - invalid free&#34;)
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span>		}
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span>		h.pagesInUse.Add(-s.npages)
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span>
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span>		<span class="comment">// Clear in-use bit in arena page bitmap.</span>
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span>		arena, pageIdx, pageMask := pageIndexOf(s.base())
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span>		atomic.And8(&amp;arena.pageInUse[pageIdx], ^pageMask)
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span>	default:
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span>		throw(&#34;mheap.freeSpanLocked - invalid span state&#34;)
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>	}
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>	<span class="comment">// Update stats.</span>
<span id="L1618" class="ln">  1618&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1619" class="ln">  1619&nbsp;&nbsp;</span>	<span class="comment">// Mirrors the code in allocSpan.</span>
<span id="L1620" class="ln">  1620&nbsp;&nbsp;</span>	nbytes := s.npages * pageSize
<span id="L1621" class="ln">  1621&nbsp;&nbsp;</span>	gcController.heapFree.add(int64(nbytes))
<span id="L1622" class="ln">  1622&nbsp;&nbsp;</span>	if typ == spanAllocHeap {
<span id="L1623" class="ln">  1623&nbsp;&nbsp;</span>		gcController.heapInUse.add(-int64(nbytes))
<span id="L1624" class="ln">  1624&nbsp;&nbsp;</span>	}
<span id="L1625" class="ln">  1625&nbsp;&nbsp;</span>	<span class="comment">// Update consistent stats.</span>
<span id="L1626" class="ln">  1626&nbsp;&nbsp;</span>	stats := memstats.heapStats.acquire()
<span id="L1627" class="ln">  1627&nbsp;&nbsp;</span>	switch typ {
<span id="L1628" class="ln">  1628&nbsp;&nbsp;</span>	case spanAllocHeap:
<span id="L1629" class="ln">  1629&nbsp;&nbsp;</span>		atomic.Xaddint64(&amp;stats.inHeap, -int64(nbytes))
<span id="L1630" class="ln">  1630&nbsp;&nbsp;</span>	case spanAllocStack:
<span id="L1631" class="ln">  1631&nbsp;&nbsp;</span>		atomic.Xaddint64(&amp;stats.inStacks, -int64(nbytes))
<span id="L1632" class="ln">  1632&nbsp;&nbsp;</span>	case spanAllocPtrScalarBits:
<span id="L1633" class="ln">  1633&nbsp;&nbsp;</span>		atomic.Xaddint64(&amp;stats.inPtrScalarBits, -int64(nbytes))
<span id="L1634" class="ln">  1634&nbsp;&nbsp;</span>	case spanAllocWorkBuf:
<span id="L1635" class="ln">  1635&nbsp;&nbsp;</span>		atomic.Xaddint64(&amp;stats.inWorkBufs, -int64(nbytes))
<span id="L1636" class="ln">  1636&nbsp;&nbsp;</span>	}
<span id="L1637" class="ln">  1637&nbsp;&nbsp;</span>	memstats.heapStats.release()
<span id="L1638" class="ln">  1638&nbsp;&nbsp;</span>
<span id="L1639" class="ln">  1639&nbsp;&nbsp;</span>	<span class="comment">// Mark the space as free.</span>
<span id="L1640" class="ln">  1640&nbsp;&nbsp;</span>	h.pages.free(s.base(), s.npages)
<span id="L1641" class="ln">  1641&nbsp;&nbsp;</span>
<span id="L1642" class="ln">  1642&nbsp;&nbsp;</span>	<span class="comment">// Free the span structure. We no longer have a use for it.</span>
<span id="L1643" class="ln">  1643&nbsp;&nbsp;</span>	s.state.set(mSpanDead)
<span id="L1644" class="ln">  1644&nbsp;&nbsp;</span>	h.freeMSpanLocked(s)
<span id="L1645" class="ln">  1645&nbsp;&nbsp;</span>}
<span id="L1646" class="ln">  1646&nbsp;&nbsp;</span>
<span id="L1647" class="ln">  1647&nbsp;&nbsp;</span><span class="comment">// scavengeAll acquires the heap lock (blocking any additional</span>
<span id="L1648" class="ln">  1648&nbsp;&nbsp;</span><span class="comment">// manipulation of the page allocator) and iterates over the whole</span>
<span id="L1649" class="ln">  1649&nbsp;&nbsp;</span><span class="comment">// heap, scavenging every free page available.</span>
<span id="L1650" class="ln">  1650&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1651" class="ln">  1651&nbsp;&nbsp;</span><span class="comment">// Must run on the system stack because it acquires the heap lock.</span>
<span id="L1652" class="ln">  1652&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1653" class="ln">  1653&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1654" class="ln">  1654&nbsp;&nbsp;</span>func (h *mheap) scavengeAll() {
<span id="L1655" class="ln">  1655&nbsp;&nbsp;</span>	<span class="comment">// Disallow malloc or panic while holding the heap lock. We do</span>
<span id="L1656" class="ln">  1656&nbsp;&nbsp;</span>	<span class="comment">// this here because this is a non-mallocgc entry-point to</span>
<span id="L1657" class="ln">  1657&nbsp;&nbsp;</span>	<span class="comment">// the mheap API.</span>
<span id="L1658" class="ln">  1658&nbsp;&nbsp;</span>	gp := getg()
<span id="L1659" class="ln">  1659&nbsp;&nbsp;</span>	gp.m.mallocing++
<span id="L1660" class="ln">  1660&nbsp;&nbsp;</span>
<span id="L1661" class="ln">  1661&nbsp;&nbsp;</span>	<span class="comment">// Force scavenge everything.</span>
<span id="L1662" class="ln">  1662&nbsp;&nbsp;</span>	released := h.pages.scavenge(^uintptr(0), nil, true)
<span id="L1663" class="ln">  1663&nbsp;&nbsp;</span>
<span id="L1664" class="ln">  1664&nbsp;&nbsp;</span>	gp.m.mallocing--
<span id="L1665" class="ln">  1665&nbsp;&nbsp;</span>
<span id="L1666" class="ln">  1666&nbsp;&nbsp;</span>	if debug.scavtrace &gt; 0 {
<span id="L1667" class="ln">  1667&nbsp;&nbsp;</span>		printScavTrace(0, released, true)
<span id="L1668" class="ln">  1668&nbsp;&nbsp;</span>	}
<span id="L1669" class="ln">  1669&nbsp;&nbsp;</span>}
<span id="L1670" class="ln">  1670&nbsp;&nbsp;</span>
<span id="L1671" class="ln">  1671&nbsp;&nbsp;</span><span class="comment">//go:linkname runtime_debug_freeOSMemory runtime/debug.freeOSMemory</span>
<span id="L1672" class="ln">  1672&nbsp;&nbsp;</span>func runtime_debug_freeOSMemory() {
<span id="L1673" class="ln">  1673&nbsp;&nbsp;</span>	GC()
<span id="L1674" class="ln">  1674&nbsp;&nbsp;</span>	systemstack(func() { mheap_.scavengeAll() })
<span id="L1675" class="ln">  1675&nbsp;&nbsp;</span>}
<span id="L1676" class="ln">  1676&nbsp;&nbsp;</span>
<span id="L1677" class="ln">  1677&nbsp;&nbsp;</span><span class="comment">// Initialize a new span with the given start and npages.</span>
<span id="L1678" class="ln">  1678&nbsp;&nbsp;</span>func (span *mspan) init(base uintptr, npages uintptr) {
<span id="L1679" class="ln">  1679&nbsp;&nbsp;</span>	<span class="comment">// span is *not* zeroed.</span>
<span id="L1680" class="ln">  1680&nbsp;&nbsp;</span>	span.next = nil
<span id="L1681" class="ln">  1681&nbsp;&nbsp;</span>	span.prev = nil
<span id="L1682" class="ln">  1682&nbsp;&nbsp;</span>	span.list = nil
<span id="L1683" class="ln">  1683&nbsp;&nbsp;</span>	span.startAddr = base
<span id="L1684" class="ln">  1684&nbsp;&nbsp;</span>	span.npages = npages
<span id="L1685" class="ln">  1685&nbsp;&nbsp;</span>	span.allocCount = 0
<span id="L1686" class="ln">  1686&nbsp;&nbsp;</span>	span.spanclass = 0
<span id="L1687" class="ln">  1687&nbsp;&nbsp;</span>	span.elemsize = 0
<span id="L1688" class="ln">  1688&nbsp;&nbsp;</span>	span.speciallock.key = 0
<span id="L1689" class="ln">  1689&nbsp;&nbsp;</span>	span.specials = nil
<span id="L1690" class="ln">  1690&nbsp;&nbsp;</span>	span.needzero = 0
<span id="L1691" class="ln">  1691&nbsp;&nbsp;</span>	span.freeindex = 0
<span id="L1692" class="ln">  1692&nbsp;&nbsp;</span>	span.freeIndexForScan = 0
<span id="L1693" class="ln">  1693&nbsp;&nbsp;</span>	span.allocBits = nil
<span id="L1694" class="ln">  1694&nbsp;&nbsp;</span>	span.gcmarkBits = nil
<span id="L1695" class="ln">  1695&nbsp;&nbsp;</span>	span.pinnerBits = nil
<span id="L1696" class="ln">  1696&nbsp;&nbsp;</span>	span.state.set(mSpanDead)
<span id="L1697" class="ln">  1697&nbsp;&nbsp;</span>	lockInit(&amp;span.speciallock, lockRankMspanSpecial)
<span id="L1698" class="ln">  1698&nbsp;&nbsp;</span>}
<span id="L1699" class="ln">  1699&nbsp;&nbsp;</span>
<span id="L1700" class="ln">  1700&nbsp;&nbsp;</span>func (span *mspan) inList() bool {
<span id="L1701" class="ln">  1701&nbsp;&nbsp;</span>	return span.list != nil
<span id="L1702" class="ln">  1702&nbsp;&nbsp;</span>}
<span id="L1703" class="ln">  1703&nbsp;&nbsp;</span>
<span id="L1704" class="ln">  1704&nbsp;&nbsp;</span><span class="comment">// Initialize an empty doubly-linked list.</span>
<span id="L1705" class="ln">  1705&nbsp;&nbsp;</span>func (list *mSpanList) init() {
<span id="L1706" class="ln">  1706&nbsp;&nbsp;</span>	list.first = nil
<span id="L1707" class="ln">  1707&nbsp;&nbsp;</span>	list.last = nil
<span id="L1708" class="ln">  1708&nbsp;&nbsp;</span>}
<span id="L1709" class="ln">  1709&nbsp;&nbsp;</span>
<span id="L1710" class="ln">  1710&nbsp;&nbsp;</span>func (list *mSpanList) remove(span *mspan) {
<span id="L1711" class="ln">  1711&nbsp;&nbsp;</span>	if span.list != list {
<span id="L1712" class="ln">  1712&nbsp;&nbsp;</span>		print(&#34;runtime: failed mSpanList.remove span.npages=&#34;, span.npages,
<span id="L1713" class="ln">  1713&nbsp;&nbsp;</span>			&#34; span=&#34;, span, &#34; prev=&#34;, span.prev, &#34; span.list=&#34;, span.list, &#34; list=&#34;, list, &#34;\n&#34;)
<span id="L1714" class="ln">  1714&nbsp;&nbsp;</span>		throw(&#34;mSpanList.remove&#34;)
<span id="L1715" class="ln">  1715&nbsp;&nbsp;</span>	}
<span id="L1716" class="ln">  1716&nbsp;&nbsp;</span>	if list.first == span {
<span id="L1717" class="ln">  1717&nbsp;&nbsp;</span>		list.first = span.next
<span id="L1718" class="ln">  1718&nbsp;&nbsp;</span>	} else {
<span id="L1719" class="ln">  1719&nbsp;&nbsp;</span>		span.prev.next = span.next
<span id="L1720" class="ln">  1720&nbsp;&nbsp;</span>	}
<span id="L1721" class="ln">  1721&nbsp;&nbsp;</span>	if list.last == span {
<span id="L1722" class="ln">  1722&nbsp;&nbsp;</span>		list.last = span.prev
<span id="L1723" class="ln">  1723&nbsp;&nbsp;</span>	} else {
<span id="L1724" class="ln">  1724&nbsp;&nbsp;</span>		span.next.prev = span.prev
<span id="L1725" class="ln">  1725&nbsp;&nbsp;</span>	}
<span id="L1726" class="ln">  1726&nbsp;&nbsp;</span>	span.next = nil
<span id="L1727" class="ln">  1727&nbsp;&nbsp;</span>	span.prev = nil
<span id="L1728" class="ln">  1728&nbsp;&nbsp;</span>	span.list = nil
<span id="L1729" class="ln">  1729&nbsp;&nbsp;</span>}
<span id="L1730" class="ln">  1730&nbsp;&nbsp;</span>
<span id="L1731" class="ln">  1731&nbsp;&nbsp;</span>func (list *mSpanList) isEmpty() bool {
<span id="L1732" class="ln">  1732&nbsp;&nbsp;</span>	return list.first == nil
<span id="L1733" class="ln">  1733&nbsp;&nbsp;</span>}
<span id="L1734" class="ln">  1734&nbsp;&nbsp;</span>
<span id="L1735" class="ln">  1735&nbsp;&nbsp;</span>func (list *mSpanList) insert(span *mspan) {
<span id="L1736" class="ln">  1736&nbsp;&nbsp;</span>	if span.next != nil || span.prev != nil || span.list != nil {
<span id="L1737" class="ln">  1737&nbsp;&nbsp;</span>		println(&#34;runtime: failed mSpanList.insert&#34;, span, span.next, span.prev, span.list)
<span id="L1738" class="ln">  1738&nbsp;&nbsp;</span>		throw(&#34;mSpanList.insert&#34;)
<span id="L1739" class="ln">  1739&nbsp;&nbsp;</span>	}
<span id="L1740" class="ln">  1740&nbsp;&nbsp;</span>	span.next = list.first
<span id="L1741" class="ln">  1741&nbsp;&nbsp;</span>	if list.first != nil {
<span id="L1742" class="ln">  1742&nbsp;&nbsp;</span>		<span class="comment">// The list contains at least one span; link it in.</span>
<span id="L1743" class="ln">  1743&nbsp;&nbsp;</span>		<span class="comment">// The last span in the list doesn&#39;t change.</span>
<span id="L1744" class="ln">  1744&nbsp;&nbsp;</span>		list.first.prev = span
<span id="L1745" class="ln">  1745&nbsp;&nbsp;</span>	} else {
<span id="L1746" class="ln">  1746&nbsp;&nbsp;</span>		<span class="comment">// The list contains no spans, so this is also the last span.</span>
<span id="L1747" class="ln">  1747&nbsp;&nbsp;</span>		list.last = span
<span id="L1748" class="ln">  1748&nbsp;&nbsp;</span>	}
<span id="L1749" class="ln">  1749&nbsp;&nbsp;</span>	list.first = span
<span id="L1750" class="ln">  1750&nbsp;&nbsp;</span>	span.list = list
<span id="L1751" class="ln">  1751&nbsp;&nbsp;</span>}
<span id="L1752" class="ln">  1752&nbsp;&nbsp;</span>
<span id="L1753" class="ln">  1753&nbsp;&nbsp;</span>func (list *mSpanList) insertBack(span *mspan) {
<span id="L1754" class="ln">  1754&nbsp;&nbsp;</span>	if span.next != nil || span.prev != nil || span.list != nil {
<span id="L1755" class="ln">  1755&nbsp;&nbsp;</span>		println(&#34;runtime: failed mSpanList.insertBack&#34;, span, span.next, span.prev, span.list)
<span id="L1756" class="ln">  1756&nbsp;&nbsp;</span>		throw(&#34;mSpanList.insertBack&#34;)
<span id="L1757" class="ln">  1757&nbsp;&nbsp;</span>	}
<span id="L1758" class="ln">  1758&nbsp;&nbsp;</span>	span.prev = list.last
<span id="L1759" class="ln">  1759&nbsp;&nbsp;</span>	if list.last != nil {
<span id="L1760" class="ln">  1760&nbsp;&nbsp;</span>		<span class="comment">// The list contains at least one span.</span>
<span id="L1761" class="ln">  1761&nbsp;&nbsp;</span>		list.last.next = span
<span id="L1762" class="ln">  1762&nbsp;&nbsp;</span>	} else {
<span id="L1763" class="ln">  1763&nbsp;&nbsp;</span>		<span class="comment">// The list contains no spans, so this is also the first span.</span>
<span id="L1764" class="ln">  1764&nbsp;&nbsp;</span>		list.first = span
<span id="L1765" class="ln">  1765&nbsp;&nbsp;</span>	}
<span id="L1766" class="ln">  1766&nbsp;&nbsp;</span>	list.last = span
<span id="L1767" class="ln">  1767&nbsp;&nbsp;</span>	span.list = list
<span id="L1768" class="ln">  1768&nbsp;&nbsp;</span>}
<span id="L1769" class="ln">  1769&nbsp;&nbsp;</span>
<span id="L1770" class="ln">  1770&nbsp;&nbsp;</span><span class="comment">// takeAll removes all spans from other and inserts them at the front</span>
<span id="L1771" class="ln">  1771&nbsp;&nbsp;</span><span class="comment">// of list.</span>
<span id="L1772" class="ln">  1772&nbsp;&nbsp;</span>func (list *mSpanList) takeAll(other *mSpanList) {
<span id="L1773" class="ln">  1773&nbsp;&nbsp;</span>	if other.isEmpty() {
<span id="L1774" class="ln">  1774&nbsp;&nbsp;</span>		return
<span id="L1775" class="ln">  1775&nbsp;&nbsp;</span>	}
<span id="L1776" class="ln">  1776&nbsp;&nbsp;</span>
<span id="L1777" class="ln">  1777&nbsp;&nbsp;</span>	<span class="comment">// Reparent everything in other to list.</span>
<span id="L1778" class="ln">  1778&nbsp;&nbsp;</span>	for s := other.first; s != nil; s = s.next {
<span id="L1779" class="ln">  1779&nbsp;&nbsp;</span>		s.list = list
<span id="L1780" class="ln">  1780&nbsp;&nbsp;</span>	}
<span id="L1781" class="ln">  1781&nbsp;&nbsp;</span>
<span id="L1782" class="ln">  1782&nbsp;&nbsp;</span>	<span class="comment">// Concatenate the lists.</span>
<span id="L1783" class="ln">  1783&nbsp;&nbsp;</span>	if list.isEmpty() {
<span id="L1784" class="ln">  1784&nbsp;&nbsp;</span>		*list = *other
<span id="L1785" class="ln">  1785&nbsp;&nbsp;</span>	} else {
<span id="L1786" class="ln">  1786&nbsp;&nbsp;</span>		<span class="comment">// Neither list is empty. Put other before list.</span>
<span id="L1787" class="ln">  1787&nbsp;&nbsp;</span>		other.last.next = list.first
<span id="L1788" class="ln">  1788&nbsp;&nbsp;</span>		list.first.prev = other.last
<span id="L1789" class="ln">  1789&nbsp;&nbsp;</span>		list.first = other.first
<span id="L1790" class="ln">  1790&nbsp;&nbsp;</span>	}
<span id="L1791" class="ln">  1791&nbsp;&nbsp;</span>
<span id="L1792" class="ln">  1792&nbsp;&nbsp;</span>	other.first, other.last = nil, nil
<span id="L1793" class="ln">  1793&nbsp;&nbsp;</span>}
<span id="L1794" class="ln">  1794&nbsp;&nbsp;</span>
<span id="L1795" class="ln">  1795&nbsp;&nbsp;</span>const (
<span id="L1796" class="ln">  1796&nbsp;&nbsp;</span>	_KindSpecialFinalizer = 1
<span id="L1797" class="ln">  1797&nbsp;&nbsp;</span>	_KindSpecialProfile   = 2
<span id="L1798" class="ln">  1798&nbsp;&nbsp;</span>	<span class="comment">// _KindSpecialReachable is a special used for tracking</span>
<span id="L1799" class="ln">  1799&nbsp;&nbsp;</span>	<span class="comment">// reachability during testing.</span>
<span id="L1800" class="ln">  1800&nbsp;&nbsp;</span>	_KindSpecialReachable = 3
<span id="L1801" class="ln">  1801&nbsp;&nbsp;</span>	<span class="comment">// _KindSpecialPinCounter is a special used for objects that are pinned</span>
<span id="L1802" class="ln">  1802&nbsp;&nbsp;</span>	<span class="comment">// multiple times</span>
<span id="L1803" class="ln">  1803&nbsp;&nbsp;</span>	_KindSpecialPinCounter = 4
<span id="L1804" class="ln">  1804&nbsp;&nbsp;</span>	<span class="comment">// Note: The finalizer special must be first because if we&#39;re freeing</span>
<span id="L1805" class="ln">  1805&nbsp;&nbsp;</span>	<span class="comment">// an object, a finalizer special will cause the freeing operation</span>
<span id="L1806" class="ln">  1806&nbsp;&nbsp;</span>	<span class="comment">// to abort, and we want to keep the other special records around</span>
<span id="L1807" class="ln">  1807&nbsp;&nbsp;</span>	<span class="comment">// if that happens.</span>
<span id="L1808" class="ln">  1808&nbsp;&nbsp;</span>)
<span id="L1809" class="ln">  1809&nbsp;&nbsp;</span>
<span id="L1810" class="ln">  1810&nbsp;&nbsp;</span>type special struct {
<span id="L1811" class="ln">  1811&nbsp;&nbsp;</span>	_      sys.NotInHeap
<span id="L1812" class="ln">  1812&nbsp;&nbsp;</span>	next   *special <span class="comment">// linked list in span</span>
<span id="L1813" class="ln">  1813&nbsp;&nbsp;</span>	offset uint16   <span class="comment">// span offset of object</span>
<span id="L1814" class="ln">  1814&nbsp;&nbsp;</span>	kind   byte     <span class="comment">// kind of special</span>
<span id="L1815" class="ln">  1815&nbsp;&nbsp;</span>}
<span id="L1816" class="ln">  1816&nbsp;&nbsp;</span>
<span id="L1817" class="ln">  1817&nbsp;&nbsp;</span><span class="comment">// spanHasSpecials marks a span as having specials in the arena bitmap.</span>
<span id="L1818" class="ln">  1818&nbsp;&nbsp;</span>func spanHasSpecials(s *mspan) {
<span id="L1819" class="ln">  1819&nbsp;&nbsp;</span>	arenaPage := (s.base() / pageSize) % pagesPerArena
<span id="L1820" class="ln">  1820&nbsp;&nbsp;</span>	ai := arenaIndex(s.base())
<span id="L1821" class="ln">  1821&nbsp;&nbsp;</span>	ha := mheap_.arenas[ai.l1()][ai.l2()]
<span id="L1822" class="ln">  1822&nbsp;&nbsp;</span>	atomic.Or8(&amp;ha.pageSpecials[arenaPage/8], uint8(1)&lt;&lt;(arenaPage%8))
<span id="L1823" class="ln">  1823&nbsp;&nbsp;</span>}
<span id="L1824" class="ln">  1824&nbsp;&nbsp;</span>
<span id="L1825" class="ln">  1825&nbsp;&nbsp;</span><span class="comment">// spanHasNoSpecials marks a span as having no specials in the arena bitmap.</span>
<span id="L1826" class="ln">  1826&nbsp;&nbsp;</span>func spanHasNoSpecials(s *mspan) {
<span id="L1827" class="ln">  1827&nbsp;&nbsp;</span>	arenaPage := (s.base() / pageSize) % pagesPerArena
<span id="L1828" class="ln">  1828&nbsp;&nbsp;</span>	ai := arenaIndex(s.base())
<span id="L1829" class="ln">  1829&nbsp;&nbsp;</span>	ha := mheap_.arenas[ai.l1()][ai.l2()]
<span id="L1830" class="ln">  1830&nbsp;&nbsp;</span>	atomic.And8(&amp;ha.pageSpecials[arenaPage/8], ^(uint8(1) &lt;&lt; (arenaPage % 8)))
<span id="L1831" class="ln">  1831&nbsp;&nbsp;</span>}
<span id="L1832" class="ln">  1832&nbsp;&nbsp;</span>
<span id="L1833" class="ln">  1833&nbsp;&nbsp;</span><span class="comment">// Adds the special record s to the list of special records for</span>
<span id="L1834" class="ln">  1834&nbsp;&nbsp;</span><span class="comment">// the object p. All fields of s should be filled in except for</span>
<span id="L1835" class="ln">  1835&nbsp;&nbsp;</span><span class="comment">// offset &amp; next, which this routine will fill in.</span>
<span id="L1836" class="ln">  1836&nbsp;&nbsp;</span><span class="comment">// Returns true if the special was successfully added, false otherwise.</span>
<span id="L1837" class="ln">  1837&nbsp;&nbsp;</span><span class="comment">// (The add will fail only if a record with the same p and s-&gt;kind</span>
<span id="L1838" class="ln">  1838&nbsp;&nbsp;</span><span class="comment">// already exists.)</span>
<span id="L1839" class="ln">  1839&nbsp;&nbsp;</span>func addspecial(p unsafe.Pointer, s *special) bool {
<span id="L1840" class="ln">  1840&nbsp;&nbsp;</span>	span := spanOfHeap(uintptr(p))
<span id="L1841" class="ln">  1841&nbsp;&nbsp;</span>	if span == nil {
<span id="L1842" class="ln">  1842&nbsp;&nbsp;</span>		throw(&#34;addspecial on invalid pointer&#34;)
<span id="L1843" class="ln">  1843&nbsp;&nbsp;</span>	}
<span id="L1844" class="ln">  1844&nbsp;&nbsp;</span>
<span id="L1845" class="ln">  1845&nbsp;&nbsp;</span>	<span class="comment">// Ensure that the span is swept.</span>
<span id="L1846" class="ln">  1846&nbsp;&nbsp;</span>	<span class="comment">// Sweeping accesses the specials list w/o locks, so we have</span>
<span id="L1847" class="ln">  1847&nbsp;&nbsp;</span>	<span class="comment">// to synchronize with it. And it&#39;s just much safer.</span>
<span id="L1848" class="ln">  1848&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L1849" class="ln">  1849&nbsp;&nbsp;</span>	span.ensureSwept()
<span id="L1850" class="ln">  1850&nbsp;&nbsp;</span>
<span id="L1851" class="ln">  1851&nbsp;&nbsp;</span>	offset := uintptr(p) - span.base()
<span id="L1852" class="ln">  1852&nbsp;&nbsp;</span>	kind := s.kind
<span id="L1853" class="ln">  1853&nbsp;&nbsp;</span>
<span id="L1854" class="ln">  1854&nbsp;&nbsp;</span>	lock(&amp;span.speciallock)
<span id="L1855" class="ln">  1855&nbsp;&nbsp;</span>
<span id="L1856" class="ln">  1856&nbsp;&nbsp;</span>	<span class="comment">// Find splice point, check for existing record.</span>
<span id="L1857" class="ln">  1857&nbsp;&nbsp;</span>	iter, exists := span.specialFindSplicePoint(offset, kind)
<span id="L1858" class="ln">  1858&nbsp;&nbsp;</span>	if !exists {
<span id="L1859" class="ln">  1859&nbsp;&nbsp;</span>		<span class="comment">// Splice in record, fill in offset.</span>
<span id="L1860" class="ln">  1860&nbsp;&nbsp;</span>		s.offset = uint16(offset)
<span id="L1861" class="ln">  1861&nbsp;&nbsp;</span>		s.next = *iter
<span id="L1862" class="ln">  1862&nbsp;&nbsp;</span>		*iter = s
<span id="L1863" class="ln">  1863&nbsp;&nbsp;</span>		spanHasSpecials(span)
<span id="L1864" class="ln">  1864&nbsp;&nbsp;</span>	}
<span id="L1865" class="ln">  1865&nbsp;&nbsp;</span>
<span id="L1866" class="ln">  1866&nbsp;&nbsp;</span>	unlock(&amp;span.speciallock)
<span id="L1867" class="ln">  1867&nbsp;&nbsp;</span>	releasem(mp)
<span id="L1868" class="ln">  1868&nbsp;&nbsp;</span>	return !exists <span class="comment">// already exists</span>
<span id="L1869" class="ln">  1869&nbsp;&nbsp;</span>}
<span id="L1870" class="ln">  1870&nbsp;&nbsp;</span>
<span id="L1871" class="ln">  1871&nbsp;&nbsp;</span><span class="comment">// Removes the Special record of the given kind for the object p.</span>
<span id="L1872" class="ln">  1872&nbsp;&nbsp;</span><span class="comment">// Returns the record if the record existed, nil otherwise.</span>
<span id="L1873" class="ln">  1873&nbsp;&nbsp;</span><span class="comment">// The caller must FixAlloc_Free the result.</span>
<span id="L1874" class="ln">  1874&nbsp;&nbsp;</span>func removespecial(p unsafe.Pointer, kind uint8) *special {
<span id="L1875" class="ln">  1875&nbsp;&nbsp;</span>	span := spanOfHeap(uintptr(p))
<span id="L1876" class="ln">  1876&nbsp;&nbsp;</span>	if span == nil {
<span id="L1877" class="ln">  1877&nbsp;&nbsp;</span>		throw(&#34;removespecial on invalid pointer&#34;)
<span id="L1878" class="ln">  1878&nbsp;&nbsp;</span>	}
<span id="L1879" class="ln">  1879&nbsp;&nbsp;</span>
<span id="L1880" class="ln">  1880&nbsp;&nbsp;</span>	<span class="comment">// Ensure that the span is swept.</span>
<span id="L1881" class="ln">  1881&nbsp;&nbsp;</span>	<span class="comment">// Sweeping accesses the specials list w/o locks, so we have</span>
<span id="L1882" class="ln">  1882&nbsp;&nbsp;</span>	<span class="comment">// to synchronize with it. And it&#39;s just much safer.</span>
<span id="L1883" class="ln">  1883&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L1884" class="ln">  1884&nbsp;&nbsp;</span>	span.ensureSwept()
<span id="L1885" class="ln">  1885&nbsp;&nbsp;</span>
<span id="L1886" class="ln">  1886&nbsp;&nbsp;</span>	offset := uintptr(p) - span.base()
<span id="L1887" class="ln">  1887&nbsp;&nbsp;</span>
<span id="L1888" class="ln">  1888&nbsp;&nbsp;</span>	var result *special
<span id="L1889" class="ln">  1889&nbsp;&nbsp;</span>	lock(&amp;span.speciallock)
<span id="L1890" class="ln">  1890&nbsp;&nbsp;</span>
<span id="L1891" class="ln">  1891&nbsp;&nbsp;</span>	iter, exists := span.specialFindSplicePoint(offset, kind)
<span id="L1892" class="ln">  1892&nbsp;&nbsp;</span>	if exists {
<span id="L1893" class="ln">  1893&nbsp;&nbsp;</span>		s := *iter
<span id="L1894" class="ln">  1894&nbsp;&nbsp;</span>		*iter = s.next
<span id="L1895" class="ln">  1895&nbsp;&nbsp;</span>		result = s
<span id="L1896" class="ln">  1896&nbsp;&nbsp;</span>	}
<span id="L1897" class="ln">  1897&nbsp;&nbsp;</span>	if span.specials == nil {
<span id="L1898" class="ln">  1898&nbsp;&nbsp;</span>		spanHasNoSpecials(span)
<span id="L1899" class="ln">  1899&nbsp;&nbsp;</span>	}
<span id="L1900" class="ln">  1900&nbsp;&nbsp;</span>	unlock(&amp;span.speciallock)
<span id="L1901" class="ln">  1901&nbsp;&nbsp;</span>	releasem(mp)
<span id="L1902" class="ln">  1902&nbsp;&nbsp;</span>	return result
<span id="L1903" class="ln">  1903&nbsp;&nbsp;</span>}
<span id="L1904" class="ln">  1904&nbsp;&nbsp;</span>
<span id="L1905" class="ln">  1905&nbsp;&nbsp;</span><span class="comment">// Find a splice point in the sorted list and check for an already existing</span>
<span id="L1906" class="ln">  1906&nbsp;&nbsp;</span><span class="comment">// record. Returns a pointer to the next-reference in the list predecessor.</span>
<span id="L1907" class="ln">  1907&nbsp;&nbsp;</span><span class="comment">// Returns true, if the referenced item is an exact match.</span>
<span id="L1908" class="ln">  1908&nbsp;&nbsp;</span>func (span *mspan) specialFindSplicePoint(offset uintptr, kind byte) (**special, bool) {
<span id="L1909" class="ln">  1909&nbsp;&nbsp;</span>	<span class="comment">// Find splice point, check for existing record.</span>
<span id="L1910" class="ln">  1910&nbsp;&nbsp;</span>	iter := &amp;span.specials
<span id="L1911" class="ln">  1911&nbsp;&nbsp;</span>	found := false
<span id="L1912" class="ln">  1912&nbsp;&nbsp;</span>	for {
<span id="L1913" class="ln">  1913&nbsp;&nbsp;</span>		s := *iter
<span id="L1914" class="ln">  1914&nbsp;&nbsp;</span>		if s == nil {
<span id="L1915" class="ln">  1915&nbsp;&nbsp;</span>			break
<span id="L1916" class="ln">  1916&nbsp;&nbsp;</span>		}
<span id="L1917" class="ln">  1917&nbsp;&nbsp;</span>		if offset == uintptr(s.offset) &amp;&amp; kind == s.kind {
<span id="L1918" class="ln">  1918&nbsp;&nbsp;</span>			found = true
<span id="L1919" class="ln">  1919&nbsp;&nbsp;</span>			break
<span id="L1920" class="ln">  1920&nbsp;&nbsp;</span>		}
<span id="L1921" class="ln">  1921&nbsp;&nbsp;</span>		if offset &lt; uintptr(s.offset) || (offset == uintptr(s.offset) &amp;&amp; kind &lt; s.kind) {
<span id="L1922" class="ln">  1922&nbsp;&nbsp;</span>			break
<span id="L1923" class="ln">  1923&nbsp;&nbsp;</span>		}
<span id="L1924" class="ln">  1924&nbsp;&nbsp;</span>		iter = &amp;s.next
<span id="L1925" class="ln">  1925&nbsp;&nbsp;</span>	}
<span id="L1926" class="ln">  1926&nbsp;&nbsp;</span>	return iter, found
<span id="L1927" class="ln">  1927&nbsp;&nbsp;</span>}
<span id="L1928" class="ln">  1928&nbsp;&nbsp;</span>
<span id="L1929" class="ln">  1929&nbsp;&nbsp;</span><span class="comment">// The described object has a finalizer set for it.</span>
<span id="L1930" class="ln">  1930&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1931" class="ln">  1931&nbsp;&nbsp;</span><span class="comment">// specialfinalizer is allocated from non-GC&#39;d memory, so any heap</span>
<span id="L1932" class="ln">  1932&nbsp;&nbsp;</span><span class="comment">// pointers must be specially handled.</span>
<span id="L1933" class="ln">  1933&nbsp;&nbsp;</span>type specialfinalizer struct {
<span id="L1934" class="ln">  1934&nbsp;&nbsp;</span>	_       sys.NotInHeap
<span id="L1935" class="ln">  1935&nbsp;&nbsp;</span>	special special
<span id="L1936" class="ln">  1936&nbsp;&nbsp;</span>	fn      *funcval <span class="comment">// May be a heap pointer.</span>
<span id="L1937" class="ln">  1937&nbsp;&nbsp;</span>	nret    uintptr
<span id="L1938" class="ln">  1938&nbsp;&nbsp;</span>	fint    *_type   <span class="comment">// May be a heap pointer, but always live.</span>
<span id="L1939" class="ln">  1939&nbsp;&nbsp;</span>	ot      *ptrtype <span class="comment">// May be a heap pointer, but always live.</span>
<span id="L1940" class="ln">  1940&nbsp;&nbsp;</span>}
<span id="L1941" class="ln">  1941&nbsp;&nbsp;</span>
<span id="L1942" class="ln">  1942&nbsp;&nbsp;</span><span class="comment">// Adds a finalizer to the object p. Returns true if it succeeded.</span>
<span id="L1943" class="ln">  1943&nbsp;&nbsp;</span>func addfinalizer(p unsafe.Pointer, f *funcval, nret uintptr, fint *_type, ot *ptrtype) bool {
<span id="L1944" class="ln">  1944&nbsp;&nbsp;</span>	lock(&amp;mheap_.speciallock)
<span id="L1945" class="ln">  1945&nbsp;&nbsp;</span>	s := (*specialfinalizer)(mheap_.specialfinalizeralloc.alloc())
<span id="L1946" class="ln">  1946&nbsp;&nbsp;</span>	unlock(&amp;mheap_.speciallock)
<span id="L1947" class="ln">  1947&nbsp;&nbsp;</span>	s.special.kind = _KindSpecialFinalizer
<span id="L1948" class="ln">  1948&nbsp;&nbsp;</span>	s.fn = f
<span id="L1949" class="ln">  1949&nbsp;&nbsp;</span>	s.nret = nret
<span id="L1950" class="ln">  1950&nbsp;&nbsp;</span>	s.fint = fint
<span id="L1951" class="ln">  1951&nbsp;&nbsp;</span>	s.ot = ot
<span id="L1952" class="ln">  1952&nbsp;&nbsp;</span>	if addspecial(p, &amp;s.special) {
<span id="L1953" class="ln">  1953&nbsp;&nbsp;</span>		<span class="comment">// This is responsible for maintaining the same</span>
<span id="L1954" class="ln">  1954&nbsp;&nbsp;</span>		<span class="comment">// GC-related invariants as markrootSpans in any</span>
<span id="L1955" class="ln">  1955&nbsp;&nbsp;</span>		<span class="comment">// situation where it&#39;s possible that markrootSpans</span>
<span id="L1956" class="ln">  1956&nbsp;&nbsp;</span>		<span class="comment">// has already run but mark termination hasn&#39;t yet.</span>
<span id="L1957" class="ln">  1957&nbsp;&nbsp;</span>		if gcphase != _GCoff {
<span id="L1958" class="ln">  1958&nbsp;&nbsp;</span>			base, span, _ := findObject(uintptr(p), 0, 0)
<span id="L1959" class="ln">  1959&nbsp;&nbsp;</span>			mp := acquirem()
<span id="L1960" class="ln">  1960&nbsp;&nbsp;</span>			gcw := &amp;mp.p.ptr().gcw
<span id="L1961" class="ln">  1961&nbsp;&nbsp;</span>			<span class="comment">// Mark everything reachable from the object</span>
<span id="L1962" class="ln">  1962&nbsp;&nbsp;</span>			<span class="comment">// so it&#39;s retained for the finalizer.</span>
<span id="L1963" class="ln">  1963&nbsp;&nbsp;</span>			if !span.spanclass.noscan() {
<span id="L1964" class="ln">  1964&nbsp;&nbsp;</span>				scanobject(base, gcw)
<span id="L1965" class="ln">  1965&nbsp;&nbsp;</span>			}
<span id="L1966" class="ln">  1966&nbsp;&nbsp;</span>			<span class="comment">// Mark the finalizer itself, since the</span>
<span id="L1967" class="ln">  1967&nbsp;&nbsp;</span>			<span class="comment">// special isn&#39;t part of the GC&#39;d heap.</span>
<span id="L1968" class="ln">  1968&nbsp;&nbsp;</span>			scanblock(uintptr(unsafe.Pointer(&amp;s.fn)), goarch.PtrSize, &amp;oneptrmask[0], gcw, nil)
<span id="L1969" class="ln">  1969&nbsp;&nbsp;</span>			releasem(mp)
<span id="L1970" class="ln">  1970&nbsp;&nbsp;</span>		}
<span id="L1971" class="ln">  1971&nbsp;&nbsp;</span>		return true
<span id="L1972" class="ln">  1972&nbsp;&nbsp;</span>	}
<span id="L1973" class="ln">  1973&nbsp;&nbsp;</span>
<span id="L1974" class="ln">  1974&nbsp;&nbsp;</span>	<span class="comment">// There was an old finalizer</span>
<span id="L1975" class="ln">  1975&nbsp;&nbsp;</span>	lock(&amp;mheap_.speciallock)
<span id="L1976" class="ln">  1976&nbsp;&nbsp;</span>	mheap_.specialfinalizeralloc.free(unsafe.Pointer(s))
<span id="L1977" class="ln">  1977&nbsp;&nbsp;</span>	unlock(&amp;mheap_.speciallock)
<span id="L1978" class="ln">  1978&nbsp;&nbsp;</span>	return false
<span id="L1979" class="ln">  1979&nbsp;&nbsp;</span>}
<span id="L1980" class="ln">  1980&nbsp;&nbsp;</span>
<span id="L1981" class="ln">  1981&nbsp;&nbsp;</span><span class="comment">// Removes the finalizer (if any) from the object p.</span>
<span id="L1982" class="ln">  1982&nbsp;&nbsp;</span>func removefinalizer(p unsafe.Pointer) {
<span id="L1983" class="ln">  1983&nbsp;&nbsp;</span>	s := (*specialfinalizer)(unsafe.Pointer(removespecial(p, _KindSpecialFinalizer)))
<span id="L1984" class="ln">  1984&nbsp;&nbsp;</span>	if s == nil {
<span id="L1985" class="ln">  1985&nbsp;&nbsp;</span>		return <span class="comment">// there wasn&#39;t a finalizer to remove</span>
<span id="L1986" class="ln">  1986&nbsp;&nbsp;</span>	}
<span id="L1987" class="ln">  1987&nbsp;&nbsp;</span>	lock(&amp;mheap_.speciallock)
<span id="L1988" class="ln">  1988&nbsp;&nbsp;</span>	mheap_.specialfinalizeralloc.free(unsafe.Pointer(s))
<span id="L1989" class="ln">  1989&nbsp;&nbsp;</span>	unlock(&amp;mheap_.speciallock)
<span id="L1990" class="ln">  1990&nbsp;&nbsp;</span>}
<span id="L1991" class="ln">  1991&nbsp;&nbsp;</span>
<span id="L1992" class="ln">  1992&nbsp;&nbsp;</span><span class="comment">// The described object is being heap profiled.</span>
<span id="L1993" class="ln">  1993&nbsp;&nbsp;</span>type specialprofile struct {
<span id="L1994" class="ln">  1994&nbsp;&nbsp;</span>	_       sys.NotInHeap
<span id="L1995" class="ln">  1995&nbsp;&nbsp;</span>	special special
<span id="L1996" class="ln">  1996&nbsp;&nbsp;</span>	b       *bucket
<span id="L1997" class="ln">  1997&nbsp;&nbsp;</span>}
<span id="L1998" class="ln">  1998&nbsp;&nbsp;</span>
<span id="L1999" class="ln">  1999&nbsp;&nbsp;</span><span class="comment">// Set the heap profile bucket associated with addr to b.</span>
<span id="L2000" class="ln">  2000&nbsp;&nbsp;</span>func setprofilebucket(p unsafe.Pointer, b *bucket) {
<span id="L2001" class="ln">  2001&nbsp;&nbsp;</span>	lock(&amp;mheap_.speciallock)
<span id="L2002" class="ln">  2002&nbsp;&nbsp;</span>	s := (*specialprofile)(mheap_.specialprofilealloc.alloc())
<span id="L2003" class="ln">  2003&nbsp;&nbsp;</span>	unlock(&amp;mheap_.speciallock)
<span id="L2004" class="ln">  2004&nbsp;&nbsp;</span>	s.special.kind = _KindSpecialProfile
<span id="L2005" class="ln">  2005&nbsp;&nbsp;</span>	s.b = b
<span id="L2006" class="ln">  2006&nbsp;&nbsp;</span>	if !addspecial(p, &amp;s.special) {
<span id="L2007" class="ln">  2007&nbsp;&nbsp;</span>		throw(&#34;setprofilebucket: profile already set&#34;)
<span id="L2008" class="ln">  2008&nbsp;&nbsp;</span>	}
<span id="L2009" class="ln">  2009&nbsp;&nbsp;</span>}
<span id="L2010" class="ln">  2010&nbsp;&nbsp;</span>
<span id="L2011" class="ln">  2011&nbsp;&nbsp;</span><span class="comment">// specialReachable tracks whether an object is reachable on the next</span>
<span id="L2012" class="ln">  2012&nbsp;&nbsp;</span><span class="comment">// GC cycle. This is used by testing.</span>
<span id="L2013" class="ln">  2013&nbsp;&nbsp;</span>type specialReachable struct {
<span id="L2014" class="ln">  2014&nbsp;&nbsp;</span>	special   special
<span id="L2015" class="ln">  2015&nbsp;&nbsp;</span>	done      bool
<span id="L2016" class="ln">  2016&nbsp;&nbsp;</span>	reachable bool
<span id="L2017" class="ln">  2017&nbsp;&nbsp;</span>}
<span id="L2018" class="ln">  2018&nbsp;&nbsp;</span>
<span id="L2019" class="ln">  2019&nbsp;&nbsp;</span><span class="comment">// specialPinCounter tracks whether an object is pinned multiple times.</span>
<span id="L2020" class="ln">  2020&nbsp;&nbsp;</span>type specialPinCounter struct {
<span id="L2021" class="ln">  2021&nbsp;&nbsp;</span>	special special
<span id="L2022" class="ln">  2022&nbsp;&nbsp;</span>	counter uintptr
<span id="L2023" class="ln">  2023&nbsp;&nbsp;</span>}
<span id="L2024" class="ln">  2024&nbsp;&nbsp;</span>
<span id="L2025" class="ln">  2025&nbsp;&nbsp;</span><span class="comment">// specialsIter helps iterate over specials lists.</span>
<span id="L2026" class="ln">  2026&nbsp;&nbsp;</span>type specialsIter struct {
<span id="L2027" class="ln">  2027&nbsp;&nbsp;</span>	pprev **special
<span id="L2028" class="ln">  2028&nbsp;&nbsp;</span>	s     *special
<span id="L2029" class="ln">  2029&nbsp;&nbsp;</span>}
<span id="L2030" class="ln">  2030&nbsp;&nbsp;</span>
<span id="L2031" class="ln">  2031&nbsp;&nbsp;</span>func newSpecialsIter(span *mspan) specialsIter {
<span id="L2032" class="ln">  2032&nbsp;&nbsp;</span>	return specialsIter{&amp;span.specials, span.specials}
<span id="L2033" class="ln">  2033&nbsp;&nbsp;</span>}
<span id="L2034" class="ln">  2034&nbsp;&nbsp;</span>
<span id="L2035" class="ln">  2035&nbsp;&nbsp;</span>func (i *specialsIter) valid() bool {
<span id="L2036" class="ln">  2036&nbsp;&nbsp;</span>	return i.s != nil
<span id="L2037" class="ln">  2037&nbsp;&nbsp;</span>}
<span id="L2038" class="ln">  2038&nbsp;&nbsp;</span>
<span id="L2039" class="ln">  2039&nbsp;&nbsp;</span>func (i *specialsIter) next() {
<span id="L2040" class="ln">  2040&nbsp;&nbsp;</span>	i.pprev = &amp;i.s.next
<span id="L2041" class="ln">  2041&nbsp;&nbsp;</span>	i.s = *i.pprev
<span id="L2042" class="ln">  2042&nbsp;&nbsp;</span>}
<span id="L2043" class="ln">  2043&nbsp;&nbsp;</span>
<span id="L2044" class="ln">  2044&nbsp;&nbsp;</span><span class="comment">// unlinkAndNext removes the current special from the list and moves</span>
<span id="L2045" class="ln">  2045&nbsp;&nbsp;</span><span class="comment">// the iterator to the next special. It returns the unlinked special.</span>
<span id="L2046" class="ln">  2046&nbsp;&nbsp;</span>func (i *specialsIter) unlinkAndNext() *special {
<span id="L2047" class="ln">  2047&nbsp;&nbsp;</span>	cur := i.s
<span id="L2048" class="ln">  2048&nbsp;&nbsp;</span>	i.s = cur.next
<span id="L2049" class="ln">  2049&nbsp;&nbsp;</span>	*i.pprev = i.s
<span id="L2050" class="ln">  2050&nbsp;&nbsp;</span>	return cur
<span id="L2051" class="ln">  2051&nbsp;&nbsp;</span>}
<span id="L2052" class="ln">  2052&nbsp;&nbsp;</span>
<span id="L2053" class="ln">  2053&nbsp;&nbsp;</span><span class="comment">// freeSpecial performs any cleanup on special s and deallocates it.</span>
<span id="L2054" class="ln">  2054&nbsp;&nbsp;</span><span class="comment">// s must already be unlinked from the specials list.</span>
<span id="L2055" class="ln">  2055&nbsp;&nbsp;</span>func freeSpecial(s *special, p unsafe.Pointer, size uintptr) {
<span id="L2056" class="ln">  2056&nbsp;&nbsp;</span>	switch s.kind {
<span id="L2057" class="ln">  2057&nbsp;&nbsp;</span>	case _KindSpecialFinalizer:
<span id="L2058" class="ln">  2058&nbsp;&nbsp;</span>		sf := (*specialfinalizer)(unsafe.Pointer(s))
<span id="L2059" class="ln">  2059&nbsp;&nbsp;</span>		queuefinalizer(p, sf.fn, sf.nret, sf.fint, sf.ot)
<span id="L2060" class="ln">  2060&nbsp;&nbsp;</span>		lock(&amp;mheap_.speciallock)
<span id="L2061" class="ln">  2061&nbsp;&nbsp;</span>		mheap_.specialfinalizeralloc.free(unsafe.Pointer(sf))
<span id="L2062" class="ln">  2062&nbsp;&nbsp;</span>		unlock(&amp;mheap_.speciallock)
<span id="L2063" class="ln">  2063&nbsp;&nbsp;</span>	case _KindSpecialProfile:
<span id="L2064" class="ln">  2064&nbsp;&nbsp;</span>		sp := (*specialprofile)(unsafe.Pointer(s))
<span id="L2065" class="ln">  2065&nbsp;&nbsp;</span>		mProf_Free(sp.b, size)
<span id="L2066" class="ln">  2066&nbsp;&nbsp;</span>		lock(&amp;mheap_.speciallock)
<span id="L2067" class="ln">  2067&nbsp;&nbsp;</span>		mheap_.specialprofilealloc.free(unsafe.Pointer(sp))
<span id="L2068" class="ln">  2068&nbsp;&nbsp;</span>		unlock(&amp;mheap_.speciallock)
<span id="L2069" class="ln">  2069&nbsp;&nbsp;</span>	case _KindSpecialReachable:
<span id="L2070" class="ln">  2070&nbsp;&nbsp;</span>		sp := (*specialReachable)(unsafe.Pointer(s))
<span id="L2071" class="ln">  2071&nbsp;&nbsp;</span>		sp.done = true
<span id="L2072" class="ln">  2072&nbsp;&nbsp;</span>		<span class="comment">// The creator frees these.</span>
<span id="L2073" class="ln">  2073&nbsp;&nbsp;</span>	case _KindSpecialPinCounter:
<span id="L2074" class="ln">  2074&nbsp;&nbsp;</span>		lock(&amp;mheap_.speciallock)
<span id="L2075" class="ln">  2075&nbsp;&nbsp;</span>		mheap_.specialPinCounterAlloc.free(unsafe.Pointer(s))
<span id="L2076" class="ln">  2076&nbsp;&nbsp;</span>		unlock(&amp;mheap_.speciallock)
<span id="L2077" class="ln">  2077&nbsp;&nbsp;</span>	default:
<span id="L2078" class="ln">  2078&nbsp;&nbsp;</span>		throw(&#34;bad special kind&#34;)
<span id="L2079" class="ln">  2079&nbsp;&nbsp;</span>		panic(&#34;not reached&#34;)
<span id="L2080" class="ln">  2080&nbsp;&nbsp;</span>	}
<span id="L2081" class="ln">  2081&nbsp;&nbsp;</span>}
<span id="L2082" class="ln">  2082&nbsp;&nbsp;</span>
<span id="L2083" class="ln">  2083&nbsp;&nbsp;</span><span class="comment">// gcBits is an alloc/mark bitmap. This is always used as gcBits.x.</span>
<span id="L2084" class="ln">  2084&nbsp;&nbsp;</span>type gcBits struct {
<span id="L2085" class="ln">  2085&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L2086" class="ln">  2086&nbsp;&nbsp;</span>	x uint8
<span id="L2087" class="ln">  2087&nbsp;&nbsp;</span>}
<span id="L2088" class="ln">  2088&nbsp;&nbsp;</span>
<span id="L2089" class="ln">  2089&nbsp;&nbsp;</span><span class="comment">// bytep returns a pointer to the n&#39;th byte of b.</span>
<span id="L2090" class="ln">  2090&nbsp;&nbsp;</span>func (b *gcBits) bytep(n uintptr) *uint8 {
<span id="L2091" class="ln">  2091&nbsp;&nbsp;</span>	return addb(&amp;b.x, n)
<span id="L2092" class="ln">  2092&nbsp;&nbsp;</span>}
<span id="L2093" class="ln">  2093&nbsp;&nbsp;</span>
<span id="L2094" class="ln">  2094&nbsp;&nbsp;</span><span class="comment">// bitp returns a pointer to the byte containing bit n and a mask for</span>
<span id="L2095" class="ln">  2095&nbsp;&nbsp;</span><span class="comment">// selecting that bit from *bytep.</span>
<span id="L2096" class="ln">  2096&nbsp;&nbsp;</span>func (b *gcBits) bitp(n uintptr) (bytep *uint8, mask uint8) {
<span id="L2097" class="ln">  2097&nbsp;&nbsp;</span>	return b.bytep(n / 8), 1 &lt;&lt; (n % 8)
<span id="L2098" class="ln">  2098&nbsp;&nbsp;</span>}
<span id="L2099" class="ln">  2099&nbsp;&nbsp;</span>
<span id="L2100" class="ln">  2100&nbsp;&nbsp;</span>const gcBitsChunkBytes = uintptr(64 &lt;&lt; 10)
<span id="L2101" class="ln">  2101&nbsp;&nbsp;</span>const gcBitsHeaderBytes = unsafe.Sizeof(gcBitsHeader{})
<span id="L2102" class="ln">  2102&nbsp;&nbsp;</span>
<span id="L2103" class="ln">  2103&nbsp;&nbsp;</span>type gcBitsHeader struct {
<span id="L2104" class="ln">  2104&nbsp;&nbsp;</span>	free uintptr <span class="comment">// free is the index into bits of the next free byte.</span>
<span id="L2105" class="ln">  2105&nbsp;&nbsp;</span>	next uintptr <span class="comment">// *gcBits triggers recursive type bug. (issue 14620)</span>
<span id="L2106" class="ln">  2106&nbsp;&nbsp;</span>}
<span id="L2107" class="ln">  2107&nbsp;&nbsp;</span>
<span id="L2108" class="ln">  2108&nbsp;&nbsp;</span>type gcBitsArena struct {
<span id="L2109" class="ln">  2109&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L2110" class="ln">  2110&nbsp;&nbsp;</span>	<span class="comment">// gcBitsHeader // side step recursive type bug (issue 14620) by including fields by hand.</span>
<span id="L2111" class="ln">  2111&nbsp;&nbsp;</span>	free uintptr <span class="comment">// free is the index into bits of the next free byte; read/write atomically</span>
<span id="L2112" class="ln">  2112&nbsp;&nbsp;</span>	next *gcBitsArena
<span id="L2113" class="ln">  2113&nbsp;&nbsp;</span>	bits [gcBitsChunkBytes - gcBitsHeaderBytes]gcBits
<span id="L2114" class="ln">  2114&nbsp;&nbsp;</span>}
<span id="L2115" class="ln">  2115&nbsp;&nbsp;</span>
<span id="L2116" class="ln">  2116&nbsp;&nbsp;</span>var gcBitsArenas struct {
<span id="L2117" class="ln">  2117&nbsp;&nbsp;</span>	lock     mutex
<span id="L2118" class="ln">  2118&nbsp;&nbsp;</span>	free     *gcBitsArena
<span id="L2119" class="ln">  2119&nbsp;&nbsp;</span>	next     *gcBitsArena <span class="comment">// Read atomically. Write atomically under lock.</span>
<span id="L2120" class="ln">  2120&nbsp;&nbsp;</span>	current  *gcBitsArena
<span id="L2121" class="ln">  2121&nbsp;&nbsp;</span>	previous *gcBitsArena
<span id="L2122" class="ln">  2122&nbsp;&nbsp;</span>}
<span id="L2123" class="ln">  2123&nbsp;&nbsp;</span>
<span id="L2124" class="ln">  2124&nbsp;&nbsp;</span><span class="comment">// tryAlloc allocates from b or returns nil if b does not have enough room.</span>
<span id="L2125" class="ln">  2125&nbsp;&nbsp;</span><span class="comment">// This is safe to call concurrently.</span>
<span id="L2126" class="ln">  2126&nbsp;&nbsp;</span>func (b *gcBitsArena) tryAlloc(bytes uintptr) *gcBits {
<span id="L2127" class="ln">  2127&nbsp;&nbsp;</span>	if b == nil || atomic.Loaduintptr(&amp;b.free)+bytes &gt; uintptr(len(b.bits)) {
<span id="L2128" class="ln">  2128&nbsp;&nbsp;</span>		return nil
<span id="L2129" class="ln">  2129&nbsp;&nbsp;</span>	}
<span id="L2130" class="ln">  2130&nbsp;&nbsp;</span>	<span class="comment">// Try to allocate from this block.</span>
<span id="L2131" class="ln">  2131&nbsp;&nbsp;</span>	end := atomic.Xadduintptr(&amp;b.free, bytes)
<span id="L2132" class="ln">  2132&nbsp;&nbsp;</span>	if end &gt; uintptr(len(b.bits)) {
<span id="L2133" class="ln">  2133&nbsp;&nbsp;</span>		return nil
<span id="L2134" class="ln">  2134&nbsp;&nbsp;</span>	}
<span id="L2135" class="ln">  2135&nbsp;&nbsp;</span>	<span class="comment">// There was enough room.</span>
<span id="L2136" class="ln">  2136&nbsp;&nbsp;</span>	start := end - bytes
<span id="L2137" class="ln">  2137&nbsp;&nbsp;</span>	return &amp;b.bits[start]
<span id="L2138" class="ln">  2138&nbsp;&nbsp;</span>}
<span id="L2139" class="ln">  2139&nbsp;&nbsp;</span>
<span id="L2140" class="ln">  2140&nbsp;&nbsp;</span><span class="comment">// newMarkBits returns a pointer to 8 byte aligned bytes</span>
<span id="L2141" class="ln">  2141&nbsp;&nbsp;</span><span class="comment">// to be used for a span&#39;s mark bits.</span>
<span id="L2142" class="ln">  2142&nbsp;&nbsp;</span>func newMarkBits(nelems uintptr) *gcBits {
<span id="L2143" class="ln">  2143&nbsp;&nbsp;</span>	blocksNeeded := (nelems + 63) / 64
<span id="L2144" class="ln">  2144&nbsp;&nbsp;</span>	bytesNeeded := blocksNeeded * 8
<span id="L2145" class="ln">  2145&nbsp;&nbsp;</span>
<span id="L2146" class="ln">  2146&nbsp;&nbsp;</span>	<span class="comment">// Try directly allocating from the current head arena.</span>
<span id="L2147" class="ln">  2147&nbsp;&nbsp;</span>	head := (*gcBitsArena)(atomic.Loadp(unsafe.Pointer(&amp;gcBitsArenas.next)))
<span id="L2148" class="ln">  2148&nbsp;&nbsp;</span>	if p := head.tryAlloc(bytesNeeded); p != nil {
<span id="L2149" class="ln">  2149&nbsp;&nbsp;</span>		return p
<span id="L2150" class="ln">  2150&nbsp;&nbsp;</span>	}
<span id="L2151" class="ln">  2151&nbsp;&nbsp;</span>
<span id="L2152" class="ln">  2152&nbsp;&nbsp;</span>	<span class="comment">// There&#39;s not enough room in the head arena. We may need to</span>
<span id="L2153" class="ln">  2153&nbsp;&nbsp;</span>	<span class="comment">// allocate a new arena.</span>
<span id="L2154" class="ln">  2154&nbsp;&nbsp;</span>	lock(&amp;gcBitsArenas.lock)
<span id="L2155" class="ln">  2155&nbsp;&nbsp;</span>	<span class="comment">// Try the head arena again, since it may have changed. Now</span>
<span id="L2156" class="ln">  2156&nbsp;&nbsp;</span>	<span class="comment">// that we hold the lock, the list head can&#39;t change, but its</span>
<span id="L2157" class="ln">  2157&nbsp;&nbsp;</span>	<span class="comment">// free position still can.</span>
<span id="L2158" class="ln">  2158&nbsp;&nbsp;</span>	if p := gcBitsArenas.next.tryAlloc(bytesNeeded); p != nil {
<span id="L2159" class="ln">  2159&nbsp;&nbsp;</span>		unlock(&amp;gcBitsArenas.lock)
<span id="L2160" class="ln">  2160&nbsp;&nbsp;</span>		return p
<span id="L2161" class="ln">  2161&nbsp;&nbsp;</span>	}
<span id="L2162" class="ln">  2162&nbsp;&nbsp;</span>
<span id="L2163" class="ln">  2163&nbsp;&nbsp;</span>	<span class="comment">// Allocate a new arena. This may temporarily drop the lock.</span>
<span id="L2164" class="ln">  2164&nbsp;&nbsp;</span>	fresh := newArenaMayUnlock()
<span id="L2165" class="ln">  2165&nbsp;&nbsp;</span>	<span class="comment">// If newArenaMayUnlock dropped the lock, another thread may</span>
<span id="L2166" class="ln">  2166&nbsp;&nbsp;</span>	<span class="comment">// have put a fresh arena on the &#34;next&#34; list. Try allocating</span>
<span id="L2167" class="ln">  2167&nbsp;&nbsp;</span>	<span class="comment">// from next again.</span>
<span id="L2168" class="ln">  2168&nbsp;&nbsp;</span>	if p := gcBitsArenas.next.tryAlloc(bytesNeeded); p != nil {
<span id="L2169" class="ln">  2169&nbsp;&nbsp;</span>		<span class="comment">// Put fresh back on the free list.</span>
<span id="L2170" class="ln">  2170&nbsp;&nbsp;</span>		<span class="comment">// TODO: Mark it &#34;already zeroed&#34;</span>
<span id="L2171" class="ln">  2171&nbsp;&nbsp;</span>		fresh.next = gcBitsArenas.free
<span id="L2172" class="ln">  2172&nbsp;&nbsp;</span>		gcBitsArenas.free = fresh
<span id="L2173" class="ln">  2173&nbsp;&nbsp;</span>		unlock(&amp;gcBitsArenas.lock)
<span id="L2174" class="ln">  2174&nbsp;&nbsp;</span>		return p
<span id="L2175" class="ln">  2175&nbsp;&nbsp;</span>	}
<span id="L2176" class="ln">  2176&nbsp;&nbsp;</span>
<span id="L2177" class="ln">  2177&nbsp;&nbsp;</span>	<span class="comment">// Allocate from the fresh arena. We haven&#39;t linked it in yet, so</span>
<span id="L2178" class="ln">  2178&nbsp;&nbsp;</span>	<span class="comment">// this cannot race and is guaranteed to succeed.</span>
<span id="L2179" class="ln">  2179&nbsp;&nbsp;</span>	p := fresh.tryAlloc(bytesNeeded)
<span id="L2180" class="ln">  2180&nbsp;&nbsp;</span>	if p == nil {
<span id="L2181" class="ln">  2181&nbsp;&nbsp;</span>		throw(&#34;markBits overflow&#34;)
<span id="L2182" class="ln">  2182&nbsp;&nbsp;</span>	}
<span id="L2183" class="ln">  2183&nbsp;&nbsp;</span>
<span id="L2184" class="ln">  2184&nbsp;&nbsp;</span>	<span class="comment">// Add the fresh arena to the &#34;next&#34; list.</span>
<span id="L2185" class="ln">  2185&nbsp;&nbsp;</span>	fresh.next = gcBitsArenas.next
<span id="L2186" class="ln">  2186&nbsp;&nbsp;</span>	atomic.StorepNoWB(unsafe.Pointer(&amp;gcBitsArenas.next), unsafe.Pointer(fresh))
<span id="L2187" class="ln">  2187&nbsp;&nbsp;</span>
<span id="L2188" class="ln">  2188&nbsp;&nbsp;</span>	unlock(&amp;gcBitsArenas.lock)
<span id="L2189" class="ln">  2189&nbsp;&nbsp;</span>	return p
<span id="L2190" class="ln">  2190&nbsp;&nbsp;</span>}
<span id="L2191" class="ln">  2191&nbsp;&nbsp;</span>
<span id="L2192" class="ln">  2192&nbsp;&nbsp;</span><span class="comment">// newAllocBits returns a pointer to 8 byte aligned bytes</span>
<span id="L2193" class="ln">  2193&nbsp;&nbsp;</span><span class="comment">// to be used for this span&#39;s alloc bits.</span>
<span id="L2194" class="ln">  2194&nbsp;&nbsp;</span><span class="comment">// newAllocBits is used to provide newly initialized spans</span>
<span id="L2195" class="ln">  2195&nbsp;&nbsp;</span><span class="comment">// allocation bits. For spans not being initialized the</span>
<span id="L2196" class="ln">  2196&nbsp;&nbsp;</span><span class="comment">// mark bits are repurposed as allocation bits when</span>
<span id="L2197" class="ln">  2197&nbsp;&nbsp;</span><span class="comment">// the span is swept.</span>
<span id="L2198" class="ln">  2198&nbsp;&nbsp;</span>func newAllocBits(nelems uintptr) *gcBits {
<span id="L2199" class="ln">  2199&nbsp;&nbsp;</span>	return newMarkBits(nelems)
<span id="L2200" class="ln">  2200&nbsp;&nbsp;</span>}
<span id="L2201" class="ln">  2201&nbsp;&nbsp;</span>
<span id="L2202" class="ln">  2202&nbsp;&nbsp;</span><span class="comment">// nextMarkBitArenaEpoch establishes a new epoch for the arenas</span>
<span id="L2203" class="ln">  2203&nbsp;&nbsp;</span><span class="comment">// holding the mark bits. The arenas are named relative to the</span>
<span id="L2204" class="ln">  2204&nbsp;&nbsp;</span><span class="comment">// current GC cycle which is demarcated by the call to finishweep_m.</span>
<span id="L2205" class="ln">  2205&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L2206" class="ln">  2206&nbsp;&nbsp;</span><span class="comment">// All current spans have been swept.</span>
<span id="L2207" class="ln">  2207&nbsp;&nbsp;</span><span class="comment">// During that sweep each span allocated room for its gcmarkBits in</span>
<span id="L2208" class="ln">  2208&nbsp;&nbsp;</span><span class="comment">// gcBitsArenas.next block. gcBitsArenas.next becomes the gcBitsArenas.current</span>
<span id="L2209" class="ln">  2209&nbsp;&nbsp;</span><span class="comment">// where the GC will mark objects and after each span is swept these bits</span>
<span id="L2210" class="ln">  2210&nbsp;&nbsp;</span><span class="comment">// will be used to allocate objects.</span>
<span id="L2211" class="ln">  2211&nbsp;&nbsp;</span><span class="comment">// gcBitsArenas.current becomes gcBitsArenas.previous where the span&#39;s</span>
<span id="L2212" class="ln">  2212&nbsp;&nbsp;</span><span class="comment">// gcAllocBits live until all the spans have been swept during this GC cycle.</span>
<span id="L2213" class="ln">  2213&nbsp;&nbsp;</span><span class="comment">// The span&#39;s sweep extinguishes all the references to gcBitsArenas.previous</span>
<span id="L2214" class="ln">  2214&nbsp;&nbsp;</span><span class="comment">// by pointing gcAllocBits into the gcBitsArenas.current.</span>
<span id="L2215" class="ln">  2215&nbsp;&nbsp;</span><span class="comment">// The gcBitsArenas.previous is released to the gcBitsArenas.free list.</span>
<span id="L2216" class="ln">  2216&nbsp;&nbsp;</span>func nextMarkBitArenaEpoch() {
<span id="L2217" class="ln">  2217&nbsp;&nbsp;</span>	lock(&amp;gcBitsArenas.lock)
<span id="L2218" class="ln">  2218&nbsp;&nbsp;</span>	if gcBitsArenas.previous != nil {
<span id="L2219" class="ln">  2219&nbsp;&nbsp;</span>		if gcBitsArenas.free == nil {
<span id="L2220" class="ln">  2220&nbsp;&nbsp;</span>			gcBitsArenas.free = gcBitsArenas.previous
<span id="L2221" class="ln">  2221&nbsp;&nbsp;</span>		} else {
<span id="L2222" class="ln">  2222&nbsp;&nbsp;</span>			<span class="comment">// Find end of previous arenas.</span>
<span id="L2223" class="ln">  2223&nbsp;&nbsp;</span>			last := gcBitsArenas.previous
<span id="L2224" class="ln">  2224&nbsp;&nbsp;</span>			for last = gcBitsArenas.previous; last.next != nil; last = last.next {
<span id="L2225" class="ln">  2225&nbsp;&nbsp;</span>			}
<span id="L2226" class="ln">  2226&nbsp;&nbsp;</span>			last.next = gcBitsArenas.free
<span id="L2227" class="ln">  2227&nbsp;&nbsp;</span>			gcBitsArenas.free = gcBitsArenas.previous
<span id="L2228" class="ln">  2228&nbsp;&nbsp;</span>		}
<span id="L2229" class="ln">  2229&nbsp;&nbsp;</span>	}
<span id="L2230" class="ln">  2230&nbsp;&nbsp;</span>	gcBitsArenas.previous = gcBitsArenas.current
<span id="L2231" class="ln">  2231&nbsp;&nbsp;</span>	gcBitsArenas.current = gcBitsArenas.next
<span id="L2232" class="ln">  2232&nbsp;&nbsp;</span>	atomic.StorepNoWB(unsafe.Pointer(&amp;gcBitsArenas.next), nil) <span class="comment">// newMarkBits calls newArena when needed</span>
<span id="L2233" class="ln">  2233&nbsp;&nbsp;</span>	unlock(&amp;gcBitsArenas.lock)
<span id="L2234" class="ln">  2234&nbsp;&nbsp;</span>}
<span id="L2235" class="ln">  2235&nbsp;&nbsp;</span>
<span id="L2236" class="ln">  2236&nbsp;&nbsp;</span><span class="comment">// newArenaMayUnlock allocates and zeroes a gcBits arena.</span>
<span id="L2237" class="ln">  2237&nbsp;&nbsp;</span><span class="comment">// The caller must hold gcBitsArena.lock. This may temporarily release it.</span>
<span id="L2238" class="ln">  2238&nbsp;&nbsp;</span>func newArenaMayUnlock() *gcBitsArena {
<span id="L2239" class="ln">  2239&nbsp;&nbsp;</span>	var result *gcBitsArena
<span id="L2240" class="ln">  2240&nbsp;&nbsp;</span>	if gcBitsArenas.free == nil {
<span id="L2241" class="ln">  2241&nbsp;&nbsp;</span>		unlock(&amp;gcBitsArenas.lock)
<span id="L2242" class="ln">  2242&nbsp;&nbsp;</span>		result = (*gcBitsArena)(sysAlloc(gcBitsChunkBytes, &amp;memstats.gcMiscSys))
<span id="L2243" class="ln">  2243&nbsp;&nbsp;</span>		if result == nil {
<span id="L2244" class="ln">  2244&nbsp;&nbsp;</span>			throw(&#34;runtime: cannot allocate memory&#34;)
<span id="L2245" class="ln">  2245&nbsp;&nbsp;</span>		}
<span id="L2246" class="ln">  2246&nbsp;&nbsp;</span>		lock(&amp;gcBitsArenas.lock)
<span id="L2247" class="ln">  2247&nbsp;&nbsp;</span>	} else {
<span id="L2248" class="ln">  2248&nbsp;&nbsp;</span>		result = gcBitsArenas.free
<span id="L2249" class="ln">  2249&nbsp;&nbsp;</span>		gcBitsArenas.free = gcBitsArenas.free.next
<span id="L2250" class="ln">  2250&nbsp;&nbsp;</span>		memclrNoHeapPointers(unsafe.Pointer(result), gcBitsChunkBytes)
<span id="L2251" class="ln">  2251&nbsp;&nbsp;</span>	}
<span id="L2252" class="ln">  2252&nbsp;&nbsp;</span>	result.next = nil
<span id="L2253" class="ln">  2253&nbsp;&nbsp;</span>	<span class="comment">// If result.bits is not 8 byte aligned adjust index so</span>
<span id="L2254" class="ln">  2254&nbsp;&nbsp;</span>	<span class="comment">// that &amp;result.bits[result.free] is 8 byte aligned.</span>
<span id="L2255" class="ln">  2255&nbsp;&nbsp;</span>	if unsafe.Offsetof(gcBitsArena{}.bits)&amp;7 == 0 {
<span id="L2256" class="ln">  2256&nbsp;&nbsp;</span>		result.free = 0
<span id="L2257" class="ln">  2257&nbsp;&nbsp;</span>	} else {
<span id="L2258" class="ln">  2258&nbsp;&nbsp;</span>		result.free = 8 - (uintptr(unsafe.Pointer(&amp;result.bits[0])) &amp; 7)
<span id="L2259" class="ln">  2259&nbsp;&nbsp;</span>	}
<span id="L2260" class="ln">  2260&nbsp;&nbsp;</span>	return result
<span id="L2261" class="ln">  2261&nbsp;&nbsp;</span>}
<span id="L2262" class="ln">  2262&nbsp;&nbsp;</span>
</pre><p><a href="mheap.go?m=text">View as plain text</a></p>

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
