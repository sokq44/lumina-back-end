<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mgcscavenge.go - Go Documentation Server</title>

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
<a href="mgcscavenge.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mgcscavenge.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Scavenging free pages.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This file implements scavenging (the release of physical pages backing mapped</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// memory) of free and unused pages in the heap as a way to deal with page-level</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// fragmentation and reduce the RSS of Go applications.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// Scavenging in Go happens on two fronts: there&#39;s the background</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// (asynchronous) scavenger and the allocation-time (synchronous) scavenger.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// The former happens on a goroutine much like the background sweeper which is</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// soft-capped at using scavengePercent of the mutator&#39;s time, based on</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// order-of-magnitude estimates of the costs of scavenging. The latter happens</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// when allocating pages from the heap.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// The scavenger&#39;s primary goal is to bring the estimated heap RSS of the</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// application down to a goal.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// Before we consider what this looks like, we need to split the world into two</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// halves. One in which a memory limit is not set, and one in which it is.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// For the former, the goal is defined as:</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//   (retainExtraPercent+100) / 100 * (heapGoal / lastHeapGoal) * lastHeapInUse</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// Essentially, we wish to have the application&#39;s RSS track the heap goal, but</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// the heap goal is defined in terms of bytes of objects, rather than pages like</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// RSS. As a result, we need to take into account for fragmentation internal to</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// spans. heapGoal / lastHeapGoal defines the ratio between the current heap goal</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// and the last heap goal, which tells us by how much the heap is growing and</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// shrinking. We estimate what the heap will grow to in terms of pages by taking</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// this ratio and multiplying it by heapInUse at the end of the last GC, which</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// allows us to account for this additional fragmentation. Note that this</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// procedure makes the assumption that the degree of fragmentation won&#39;t change</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// dramatically over the next GC cycle. Overestimating the amount of</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// fragmentation simply results in higher memory use, which will be accounted</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// for by the next pacing up date. Underestimating the fragmentation however</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// could lead to performance degradation. Handling this case is not within the</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// scope of the scavenger. Situations where the amount of fragmentation balloons</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// over the course of a single GC cycle should be considered pathologies,</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// flagged as bugs, and fixed appropriately.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// An additional factor of retainExtraPercent is added as a buffer to help ensure</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// that there&#39;s more unscavenged memory to allocate out of, since each allocation</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// out of scavenged memory incurs a potentially expensive page fault.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// If a memory limit is set, then we wish to pick a scavenge goal that maintains</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// that memory limit. For that, we look at total memory that has been committed</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// (memstats.mappedReady) and try to bring that down below the limit. In this case,</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// we want to give buffer space in the *opposite* direction. When the application</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// is close to the limit, we want to make sure we push harder to keep it under, so</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// if we target below the memory limit, we ensure that the background scavenger is</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// giving the situation the urgency it deserves.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// In this case, the goal is defined as:</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//    (100-reduceExtraPercent) / 100 * memoryLimit</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// We compute both of these goals, and check whether either of them have been met.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// The background scavenger continues operating as long as either one of the goals</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// has not been met.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// The goals are updated after each GC.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// Synchronous scavenging happens for one of two reasons: if an allocation would</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// exceed the memory limit or whenever the heap grows in size, for some</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// definition of heap-growth. The intuition behind this second reason is that the</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// application had to grow the heap because existing fragments were not sufficiently</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// large to satisfy a page-level memory allocation, so we scavenge those fragments</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// eagerly to offset the growth in RSS that results.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// Lastly, not all pages are available for scavenging at all times and in all cases.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// The background scavenger and heap-growth scavenger only release memory in chunks</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// that have not been densely-allocated for at least 1 full GC cycle. The reason</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// behind this is likelihood of reuse: the Go heap is allocated in a first-fit order</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// and by the end of the GC mark phase, the heap tends to be densely packed. Releasing</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// memory in these densely packed chunks while they&#39;re being packed is counter-productive,</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// and worse, it breaks up huge pages on systems that support them. The scavenger (invoked</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// during memory allocation) further ensures that chunks it identifies as &#34;dense&#34; are</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// immediately eligible for being backed by huge pages. Note that for the most part these</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// density heuristics are best-effort heuristics. It&#39;s totally possible (but unlikely)</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// that a chunk that just became dense is scavenged in the case of a race between memory</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// allocation and scavenging.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// When synchronously scavenging for the memory limit or for debug.FreeOSMemory, these</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// &#34;dense&#34; packing heuristics are ignored (in other words, scavenging is &#34;forced&#34;) because</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// in these scenarios returning memory to the OS is more important than keeping CPU</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// overheads low.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>package runtime
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>import (
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	&#34;internal/goos&#34;
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>const (
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// The background scavenger is paced according to these parameters.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// scavengePercent represents the portion of mutator time we&#39;re willing</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// to spend on scavenging in percent.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	scavengePercent = 1 <span class="comment">// 1%</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// retainExtraPercent represents the amount of memory over the heap goal</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">// that the scavenger should keep as a buffer space for the allocator.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// This constant is used when we do not have a memory limit set.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// The purpose of maintaining this overhead is to have a greater pool of</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// unscavenged memory available for allocation (since using scavenged memory</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// incurs an additional cost), to account for heap fragmentation and</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// the ever-changing layout of the heap.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	retainExtraPercent = 10
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// reduceExtraPercent represents the amount of memory under the limit</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// that the scavenger should target. For example, 5 means we target 95%</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// of the limit.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// The purpose of shooting lower than the limit is to ensure that, once</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// close to the limit, the scavenger is working hard to maintain it. If</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// we have a memory limit set but are far away from it, there&#39;s no harm</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// in leaving up to 100-retainExtraPercent live, and it&#39;s more efficient</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// anyway, for the same reasons that retainExtraPercent exists.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	reduceExtraPercent = 5
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// maxPagesPerPhysPage is the maximum number of supported runtime pages per</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// physical page, based on maxPhysPageSize.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	maxPagesPerPhysPage = maxPhysPageSize / pageSize
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// scavengeCostRatio is the approximate ratio between the costs of using previously</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// scavenged memory and scavenging memory.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// For most systems the cost of scavenging greatly outweighs the costs</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// associated with using scavenged memory, making this constant 0. On other systems</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// (especially ones where &#34;sysUsed&#34; is not just a no-op) this cost is non-trivial.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// This ratio is used as part of multiplicative factor to help the scavenger account</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// for the additional costs of using scavenged memory in its pacing.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	scavengeCostRatio = 0.7 * (goos.IsDarwin + goos.IsIos)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">// scavChunkHiOcFrac indicates the fraction of pages that need to be allocated</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// in the chunk in a single GC cycle for it to be considered high density.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	scavChunkHiOccFrac  = 0.96875
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	scavChunkHiOccPages = uint16(scavChunkHiOccFrac * pallocChunkPages)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// heapRetained returns an estimate of the current heap RSS.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>func heapRetained() uint64 {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	return gcController.heapInUse.load() + gcController.heapFree.load()
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// gcPaceScavenger updates the scavenger&#39;s pacing, particularly</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// its rate and RSS goal. For this, it requires the current heapGoal,</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// and the heapGoal for the previous GC cycle.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">// The RSS goal is based on the current heap goal with a small overhead</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// to accommodate non-determinism in the allocator.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">// The pacing is based on scavengePageRate, which applies to both regular and</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// huge pages. See that constant for more information.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// Must be called whenever GC pacing is updated.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// mheap_.lock must be held or the world must be stopped.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>func gcPaceScavenger(memoryLimit int64, heapGoal, lastHeapGoal uint64) {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	assertWorldStoppedOrLockHeld(&amp;mheap_.lock)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// As described at the top of this file, there are two scavenge goals here: one</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// for gcPercent and one for memoryLimit. Let&#39;s handle the latter first because</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// it&#39;s simpler.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// We want to target retaining (100-reduceExtraPercent)% of the heap.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	memoryLimitGoal := uint64(float64(memoryLimit) * (1 - reduceExtraPercent/100.0))
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// mappedReady is comparable to memoryLimit, and represents how much total memory</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// the Go runtime has committed now (estimated).</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	mappedReady := gcController.mappedReady.Load()
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;re below the goal already indicate that we don&#39;t need the background</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// scavenger for the memory limit. This may seems worrisome at first, but note</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// that the allocator will assist the background scavenger in the face of a memory</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// limit, so we&#39;ll be safe even if we stop the scavenger when we shouldn&#39;t have.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if mappedReady &lt;= memoryLimitGoal {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		scavenge.memoryLimitGoal.Store(^uint64(0))
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	} else {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		scavenge.memoryLimitGoal.Store(memoryLimitGoal)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// Now handle the gcPercent goal.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;re called before the first GC completed, disable scavenging.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// We never scavenge before the 2nd GC cycle anyway (we don&#39;t have enough</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// information about the heap yet) so this is fine, and avoids a fault</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// or garbage data later.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	if lastHeapGoal == 0 {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		scavenge.gcPercentGoal.Store(^uint64(0))
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		return
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// Compute our scavenging goal.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	goalRatio := float64(heapGoal) / float64(lastHeapGoal)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	gcPercentGoal := uint64(float64(memstats.lastHeapInUse) * goalRatio)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// Add retainExtraPercent overhead to retainedGoal. This calculation</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// looks strange but the purpose is to arrive at an integer division</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// (e.g. if retainExtraPercent = 12.5, then we get a divisor of 8)</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// that also avoids the overflow from a multiplication.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	gcPercentGoal += gcPercentGoal / (1.0 / (retainExtraPercent / 100.0))
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">// Align it to a physical page boundary to make the following calculations</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">// a bit more exact.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	gcPercentGoal = (gcPercentGoal + uint64(physPageSize) - 1) &amp;^ (uint64(physPageSize) - 1)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// Represents where we are now in the heap&#39;s contribution to RSS in bytes.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// Guaranteed to always be a multiple of physPageSize on systems where</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	<span class="comment">// physPageSize &lt;= pageSize since we map new heap memory at a size larger than</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	<span class="comment">// any physPageSize and released memory in multiples of the physPageSize.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	<span class="comment">// However, certain functions recategorize heap memory as other stats (e.g.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// stacks) and this happens in multiples of pageSize, so on systems</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// where physPageSize &gt; pageSize the calculations below will not be exact.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">// Generally this is OK since we&#39;ll be off by at most one regular</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	<span class="comment">// physical page.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	heapRetainedNow := heapRetained()
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;re already below our goal, or within one page of our goal, then indicate</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	<span class="comment">// that we don&#39;t need the background scavenger for maintaining a memory overhead</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// proportional to the heap goal.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	if heapRetainedNow &lt;= gcPercentGoal || heapRetainedNow-gcPercentGoal &lt; uint64(physPageSize) {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		scavenge.gcPercentGoal.Store(^uint64(0))
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	} else {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		scavenge.gcPercentGoal.Store(gcPercentGoal)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>var scavenge struct {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// gcPercentGoal is the amount of retained heap memory (measured by</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// heapRetained) that the runtime will try to maintain by returning</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	<span class="comment">// memory to the OS. This goal is derived from gcController.gcPercent</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	<span class="comment">// by choosing to retain enough memory to allocate heap memory up to</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// the heap goal.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	gcPercentGoal atomic.Uint64
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// memoryLimitGoal is the amount of memory retained by the runtime (</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// measured by gcController.mappedReady) that the runtime will try to</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">// maintain by returning memory to the OS. This goal is derived from</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	<span class="comment">// gcController.memoryLimit by choosing to target the memory limit or</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// some lower target to keep the scavenger working.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	memoryLimitGoal atomic.Uint64
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// assistTime is the time spent by the allocator scavenging in the last GC cycle.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// This is reset once a GC cycle ends.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	assistTime atomic.Int64
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// backgroundTime is the time spent by the background scavenger in the last GC cycle.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	<span class="comment">// This is reset once a GC cycle ends.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	backgroundTime atomic.Int64
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>const (
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// It doesn&#39;t really matter what value we start at, but we can&#39;t be zero, because</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// that&#39;ll cause divide-by-zero issues. Pick something conservative which we&#39;ll</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	<span class="comment">// also use as a fallback.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	startingScavSleepRatio = 0.001
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// Spend at least 1 ms scavenging, otherwise the corresponding</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	<span class="comment">// sleep time to maintain our desired utilization is too low to</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	<span class="comment">// be reliable.</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	minScavWorkTime = 1e6
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span><span class="comment">// Sleep/wait state of the background scavenger.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>var scavenger scavengerState
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>type scavengerState struct {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	<span class="comment">// lock protects all fields below.</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	lock mutex
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	<span class="comment">// g is the goroutine the scavenger is bound to.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	g *g
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	<span class="comment">// parked is whether or not the scavenger is parked.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	parked bool
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	<span class="comment">// timer is the timer used for the scavenger to sleep.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	timer *timer
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">// sysmonWake signals to sysmon that it should wake the scavenger.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	sysmonWake atomic.Uint32
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// targetCPUFraction is the target CPU overhead for the scavenger.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	targetCPUFraction float64
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">// sleepRatio is the ratio of time spent doing scavenging work to</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// time spent sleeping. This is used to decide how long the scavenger</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	<span class="comment">// should sleep for in between batches of work. It is set by</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">// critSleepController in order to maintain a CPU overhead of</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">// targetCPUFraction.</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	<span class="comment">// Lower means more sleep, higher means more aggressive scavenging.</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	sleepRatio float64
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">// sleepController controls sleepRatio.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	<span class="comment">// See sleepRatio for more details.</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	sleepController piController
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	<span class="comment">// controllerCooldown is the time left in nanoseconds during which we avoid</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	<span class="comment">// using the controller and we hold sleepRatio at a conservative</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	<span class="comment">// value. Used if the controller&#39;s assumptions fail to hold.</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	controllerCooldown int64
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	<span class="comment">// printControllerReset instructs printScavTrace to signal that</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// the controller was reset.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	printControllerReset bool
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	<span class="comment">// sleepStub is a stub used for testing to avoid actually having</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	<span class="comment">// the scavenger sleep.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	<span class="comment">// Unlike the other stubs, this is not populated if left nil</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	<span class="comment">// Instead, it is called when non-nil because any valid implementation</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	<span class="comment">// of this function basically requires closing over this scavenger</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// state, and allocating a closure is not allowed in the runtime as</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">// a matter of policy.</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	sleepStub func(n int64) int64
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	<span class="comment">// scavenge is a function that scavenges n bytes of memory.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">// Returns how many bytes of memory it actually scavenged, as</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	<span class="comment">// well as the time it took in nanoseconds. Usually mheap.pages.scavenge</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">// with nanotime called around it, but stubbed out for testing.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// Like mheap.pages.scavenge, if it scavenges less than n bytes of</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// memory, the caller may assume the heap is exhausted of scavengable</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// memory for now.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// If this is nil, it is populated with the real thing in init.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	scavenge func(n uintptr) (uintptr, int64)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">// shouldStop is a callback called in the work loop and provides a</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	<span class="comment">// point that can force the scavenger to stop early, for example because</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	<span class="comment">// the scavenge policy dictates too much has been scavenged already.</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	<span class="comment">// If this is nil, it is populated with the real thing in init.</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	shouldStop func() bool
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	<span class="comment">// gomaxprocs returns the current value of gomaxprocs. Stub for testing.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	<span class="comment">// If this is nil, it is populated with the real thing in init.</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	gomaxprocs func() int32
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// init initializes a scavenger state and wires to the current G.</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// Must be called from a regular goroutine that can allocate.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>func (s *scavengerState) init() {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	if s.g != nil {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		throw(&#34;scavenger state is already wired&#34;)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	lockInit(&amp;s.lock, lockRankScavenge)
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	s.g = getg()
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	s.timer = new(timer)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	s.timer.arg = s
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	s.timer.f = func(s any, _ uintptr) {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		s.(*scavengerState).wake()
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// input: fraction of CPU time actually used.</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	<span class="comment">// setpoint: ideal CPU fraction.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	<span class="comment">// output: ratio of time worked to time slept (determines sleep time).</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	<span class="comment">// The output of this controller is somewhat indirect to what we actually</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	<span class="comment">// want to achieve: how much time to sleep for. The reason for this definition</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	<span class="comment">// is to ensure that the controller&#39;s outputs have a direct relationship with</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	<span class="comment">// its inputs (as opposed to an inverse relationship), making it somewhat</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	<span class="comment">// easier to reason about for tuning purposes.</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	s.sleepController = piController{
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		<span class="comment">// Tuned loosely via Ziegler-Nichols process.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		kp: 0.3375,
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		ti: 3.2e6,
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		tt: 1e9, <span class="comment">// 1 second reset time.</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		<span class="comment">// These ranges seem wide, but we want to give the controller plenty of</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		<span class="comment">// room to hunt for the optimal value.</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		min: 0.001,  <span class="comment">// 1:1000</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		max: 1000.0, <span class="comment">// 1000:1</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	s.sleepRatio = startingScavSleepRatio
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	<span class="comment">// Install real functions if stubs aren&#39;t present.</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	if s.scavenge == nil {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		s.scavenge = func(n uintptr) (uintptr, int64) {
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			start := nanotime()
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			r := mheap_.pages.scavenge(n, nil, false)
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			end := nanotime()
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			if start &gt;= end {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>				return r, 0
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			scavenge.backgroundTime.Add(end - start)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			return r, end - start
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	}
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	if s.shouldStop == nil {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		s.shouldStop = func() bool {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			<span class="comment">// If background scavenging is disabled or if there&#39;s no work to do just stop.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			return heapRetained() &lt;= scavenge.gcPercentGoal.Load() &amp;&amp;
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>				gcController.mappedReady.Load() &lt;= scavenge.memoryLimitGoal.Load()
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	if s.gomaxprocs == nil {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		s.gomaxprocs = func() int32 {
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			return gomaxprocs
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span><span class="comment">// park parks the scavenger goroutine.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>func (s *scavengerState) park() {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	lock(&amp;s.lock)
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	if getg() != s.g {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		throw(&#34;tried to park scavenger from another goroutine&#34;)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	s.parked = true
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	goparkunlock(&amp;s.lock, waitReasonGCScavengeWait, traceBlockSystemGoroutine, 2)
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span><span class="comment">// ready signals to sysmon that the scavenger should be awoken.</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>func (s *scavengerState) ready() {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	s.sysmonWake.Store(1)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span><span class="comment">// wake immediately unparks the scavenger if necessary.</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span><span class="comment">// Safe to run without a P.</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>func (s *scavengerState) wake() {
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	lock(&amp;s.lock)
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	if s.parked {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		<span class="comment">// Unset sysmonWake, since the scavenger is now being awoken.</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		s.sysmonWake.Store(0)
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		<span class="comment">// s.parked is unset to prevent a double wake-up.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		s.parked = false
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		<span class="comment">// Ready the goroutine by injecting it. We use injectglist instead</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		<span class="comment">// of ready or goready in order to allow us to run this function</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		<span class="comment">// without a P. injectglist also avoids placing the goroutine in</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		<span class="comment">// the current P&#39;s runnext slot, which is desirable to prevent</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		<span class="comment">// the scavenger from interfering with user goroutine scheduling</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		<span class="comment">// too much.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		var list gList
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		list.push(s.g)
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		injectglist(&amp;list)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	unlock(&amp;s.lock)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span><span class="comment">// sleep puts the scavenger to sleep based on the amount of time that it worked</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span><span class="comment">// in nanoseconds.</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span><span class="comment">// Note that this function should only be called by the scavenger.</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span><span class="comment">// The scavenger may be woken up earlier by a pacing change, and it may not go</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">// to sleep at all if there&#39;s a pending pacing change.</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>func (s *scavengerState) sleep(worked float64) {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	lock(&amp;s.lock)
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	if getg() != s.g {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		throw(&#34;tried to sleep scavenger from another goroutine&#34;)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	if worked &lt; minScavWorkTime {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		<span class="comment">// This means there wasn&#39;t enough work to actually fill up minScavWorkTime.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		<span class="comment">// That&#39;s fine; we shouldn&#39;t try to do anything with this information</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		<span class="comment">// because it&#39;s going result in a short enough sleep request that things</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		<span class="comment">// will get messy. Just assume we did at least this much work.</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		<span class="comment">// All this means is that we&#39;ll sleep longer than we otherwise would have.</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		worked = minScavWorkTime
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	<span class="comment">// Multiply the critical time by 1 + the ratio of the costs of using</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	<span class="comment">// scavenged memory vs. scavenging memory. This forces us to pay down</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	<span class="comment">// the cost of reusing this memory eagerly by sleeping for a longer period</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	<span class="comment">// of time and scavenging less frequently. More concretely, we avoid situations</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	<span class="comment">// where we end up scavenging so often that we hurt allocation performance</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	<span class="comment">// because of the additional overheads of using scavenged memory.</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	worked *= 1 + scavengeCostRatio
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	<span class="comment">// sleepTime is the amount of time we&#39;re going to sleep, based on the amount</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	<span class="comment">// of time we worked, and the sleepRatio.</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	sleepTime := int64(worked / s.sleepRatio)
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	var slept int64
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	if s.sleepStub == nil {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		<span class="comment">// Set the timer.</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		<span class="comment">// This must happen here instead of inside gopark</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		<span class="comment">// because we can&#39;t close over any variables without</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		<span class="comment">// failing escape analysis.</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		start := nanotime()
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		resetTimer(s.timer, start+sleepTime)
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		<span class="comment">// Mark ourselves as asleep and go to sleep.</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		s.parked = true
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		goparkunlock(&amp;s.lock, waitReasonSleep, traceBlockSleep, 2)
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		<span class="comment">// How long we actually slept for.</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		slept = nanotime() - start
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		lock(&amp;s.lock)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		<span class="comment">// Stop the timer here because s.wake is unable to do it for us.</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		<span class="comment">// We don&#39;t really care if we succeed in stopping the timer. One</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		<span class="comment">// reason we might fail is that we&#39;ve already woken up, but the timer</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		<span class="comment">// might be in the process of firing on some other P; essentially we&#39;re</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		<span class="comment">// racing with it. That&#39;s totally OK. Double wake-ups are perfectly safe.</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		stopTimer(s.timer)
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		unlock(&amp;s.lock)
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	} else {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		unlock(&amp;s.lock)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		slept = s.sleepStub(sleepTime)
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	<span class="comment">// Stop here if we&#39;re cooling down from the controller.</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	if s.controllerCooldown &gt; 0 {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		<span class="comment">// worked and slept aren&#39;t exact measures of time, but it&#39;s OK to be a bit</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		<span class="comment">// sloppy here. We&#39;re just hoping we&#39;re avoiding some transient bad behavior.</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		t := slept + int64(worked)
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		if t &gt; s.controllerCooldown {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>			s.controllerCooldown = 0
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		} else {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			s.controllerCooldown -= t
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		}
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		return
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	<span class="comment">// idealFraction is the ideal % of overall application CPU time that we</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	<span class="comment">// spend scavenging.</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	idealFraction := float64(scavengePercent) / 100.0
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	<span class="comment">// Calculate the CPU time spent.</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	<span class="comment">// This may be slightly inaccurate with respect to GOMAXPROCS, but we&#39;re</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	<span class="comment">// recomputing this often enough relative to GOMAXPROCS changes in general</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	<span class="comment">// (it only changes when the world is stopped, and not during a GC) that</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	<span class="comment">// that small inaccuracy is in the noise.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	cpuFraction := worked / ((float64(slept) + worked) * float64(s.gomaxprocs()))
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	<span class="comment">// Update the critSleepRatio, adjusting until we reach our ideal fraction.</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	var ok bool
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	s.sleepRatio, ok = s.sleepController.next(cpuFraction, idealFraction, float64(slept)+worked)
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	if !ok {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		<span class="comment">// The core assumption of the controller, that we can get a proportional</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		<span class="comment">// response, broke down. This may be transient, so temporarily switch to</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		<span class="comment">// sleeping a fixed, conservative amount.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		s.sleepRatio = startingScavSleepRatio
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		s.controllerCooldown = 5e9 <span class="comment">// 5 seconds.</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		<span class="comment">// Signal the scav trace printer to output this.</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		s.controllerFailed()
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	}
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span><span class="comment">// controllerFailed indicates that the scavenger&#39;s scheduling</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span><span class="comment">// controller failed.</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>func (s *scavengerState) controllerFailed() {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	lock(&amp;s.lock)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	s.printControllerReset = true
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	unlock(&amp;s.lock)
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span><span class="comment">// run is the body of the main scavenging loop.</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span><span class="comment">// Returns the number of bytes released and the estimated time spent</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span><span class="comment">// releasing those bytes.</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span><span class="comment">// Must be run on the scavenger goroutine.</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>func (s *scavengerState) run() (released uintptr, worked float64) {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	lock(&amp;s.lock)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	if getg() != s.g {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		throw(&#34;tried to run scavenger from another goroutine&#34;)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	unlock(&amp;s.lock)
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	for worked &lt; minScavWorkTime {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		<span class="comment">// If something from outside tells us to stop early, stop.</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		if s.shouldStop() {
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>			break
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		<span class="comment">// scavengeQuantum is the amount of memory we try to scavenge</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		<span class="comment">// in one go. A smaller value means the scavenger is more responsive</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		<span class="comment">// to the scheduler in case of e.g. preemption. A larger value means</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		<span class="comment">// that the overheads of scavenging are better amortized, so better</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		<span class="comment">// scavenging throughput.</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		<span class="comment">// The current value is chosen assuming a cost of ~10µs/physical page</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		<span class="comment">// (this is somewhat pessimistic), which implies a worst-case latency of</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		<span class="comment">// about 160µs for 4 KiB physical pages. The current value is biased</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		<span class="comment">// toward latency over throughput.</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		const scavengeQuantum = 64 &lt;&lt; 10
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		<span class="comment">// Accumulate the amount of time spent scavenging.</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		r, duration := s.scavenge(scavengeQuantum)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		<span class="comment">// On some platforms we may see end &gt;= start if the time it takes to scavenge</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		<span class="comment">// memory is less than the minimum granularity of its clock (e.g. Windows) or</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		<span class="comment">// due to clock bugs.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		<span class="comment">// In this case, just assume scavenging takes 10 µs per regular physical page</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		<span class="comment">// (determined empirically), and conservatively ignore the impact of huge pages</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		<span class="comment">// on timing.</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		const approxWorkedNSPerPhysicalPage = 10e3
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		if duration == 0 {
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>			worked += approxWorkedNSPerPhysicalPage * float64(r/physPageSize)
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		} else {
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>			<span class="comment">// TODO(mknyszek): If duration is small compared to worked, it could be</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>			<span class="comment">// rounded down to zero. Probably not a problem in practice because the</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>			<span class="comment">// values are all within a few orders of magnitude of each other but maybe</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>			<span class="comment">// worth worrying about.</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>			worked += float64(duration)
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>		}
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		released += r
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		<span class="comment">// scavenge does not return until it either finds the requisite amount of</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		<span class="comment">// memory to scavenge, or exhausts the heap. If we haven&#39;t found enough</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		<span class="comment">// to scavenge, then the heap must be exhausted.</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		if r &lt; scavengeQuantum {
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			break
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		}
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		<span class="comment">// When using fake time just do one loop.</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		if faketime != 0 {
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>			break
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>		}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	}
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	if released &gt; 0 &amp;&amp; released &lt; physPageSize {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		<span class="comment">// If this happens, it means that we may have attempted to release part</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		<span class="comment">// of a physical page, but the likely effect of that is that it released</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		<span class="comment">// the whole physical page, some of which may have still been in-use.</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		<span class="comment">// This could lead to memory corruption. Throw.</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		throw(&#34;released less than one physical page of memory&#34;)
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	}
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	return
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span><span class="comment">// Background scavenger.</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span><span class="comment">// The background scavenger maintains the RSS of the application below</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span><span class="comment">// the line described by the proportional scavenging statistics in</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span><span class="comment">// the mheap struct.</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>func bgscavenge(c chan int) {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	scavenger.init()
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	c &lt;- 1
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	scavenger.park()
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	for {
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		released, workTime := scavenger.run()
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		if released == 0 {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>			scavenger.park()
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>			continue
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		}
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		mheap_.pages.scav.releasedBg.Add(released)
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		scavenger.sleep(workTime)
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	}
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>}
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span><span class="comment">// scavenge scavenges nbytes worth of free pages, starting with the</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span><span class="comment">// highest address first. Successive calls continue from where it left</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span><span class="comment">// off until the heap is exhausted. force makes all memory available to</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span><span class="comment">// scavenge, ignoring huge page heuristics.</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span><span class="comment">// Returns the amount of memory scavenged in bytes.</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span><span class="comment">// scavenge always tries to scavenge nbytes worth of memory, and will</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span><span class="comment">// only fail to do so if the heap is exhausted for now.</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>func (p *pageAlloc) scavenge(nbytes uintptr, shouldStop func() bool, force bool) uintptr {
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	released := uintptr(0)
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	for released &lt; nbytes {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		ci, pageIdx := p.scav.index.find(force)
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		if ci == 0 {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>			break
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		}
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>			released += p.scavengeOne(ci, pageIdx, nbytes-released)
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		})
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		if shouldStop != nil &amp;&amp; shouldStop() {
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>			break
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		}
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	return released
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>}
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span><span class="comment">// printScavTrace prints a scavenge trace line to standard error.</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span><span class="comment">// released should be the amount of memory released since the last time this</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span><span class="comment">// was called, and forced indicates whether the scavenge was forced by the</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span><span class="comment">// application.</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span><span class="comment">// scavenger.lock must be held.</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>func printScavTrace(releasedBg, releasedEager uintptr, forced bool) {
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	assertLockHeld(&amp;scavenger.lock)
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	printlock()
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	print(&#34;scav &#34;,
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		releasedBg&gt;&gt;10, &#34; KiB work (bg), &#34;,
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		releasedEager&gt;&gt;10, &#34; KiB work (eager), &#34;,
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		gcController.heapReleased.load()&gt;&gt;10, &#34; KiB now, &#34;,
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		(gcController.heapInUse.load()*100)/heapRetained(), &#34;% util&#34;,
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	)
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	if forced {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>		print(&#34; (forced)&#34;)
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	} else if scavenger.printControllerReset {
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>		print(&#34; [controller reset]&#34;)
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		scavenger.printControllerReset = false
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	}
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	println()
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	printunlock()
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>}
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span><span class="comment">// scavengeOne walks over the chunk at chunk index ci and searches for</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span><span class="comment">// a contiguous run of pages to scavenge. It will try to scavenge</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span><span class="comment">// at most max bytes at once, but may scavenge more to avoid</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span><span class="comment">// breaking huge pages. Once it scavenges some memory it returns</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span><span class="comment">// how much it scavenged in bytes.</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span><span class="comment">// searchIdx is the page index to start searching from in ci.</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span><span class="comment">// Returns the number of bytes scavenged.</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span><span class="comment">// Must run on the systemstack because it acquires p.mheapLock.</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>func (p *pageAlloc) scavengeOne(ci chunkIdx, searchIdx uint, max uintptr) uintptr {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	<span class="comment">// Calculate the maximum number of pages to scavenge.</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	<span class="comment">// This should be alignUp(max, pageSize) / pageSize but max can and will</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	<span class="comment">// be ^uintptr(0), so we need to be very careful not to overflow here.</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	<span class="comment">// Rather than use alignUp, calculate the number of pages rounded down</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	<span class="comment">// first, then add back one if necessary.</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	maxPages := max / pageSize
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	if max%pageSize != 0 {
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		maxPages++
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	}
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	<span class="comment">// Calculate the minimum number of pages we can scavenge.</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	<span class="comment">// Because we can only scavenge whole physical pages, we must</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	<span class="comment">// ensure that we scavenge at least minPages each time, aligned</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	<span class="comment">// to minPages*pageSize.</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	minPages := physPageSize / pageSize
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	if minPages &lt; 1 {
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		minPages = 1
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	}
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	lock(p.mheapLock)
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	if p.summary[len(p.summary)-1][ci].max() &gt;= uint(minPages) {
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		<span class="comment">// We only bother looking for a candidate if there at least</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>		<span class="comment">// minPages free pages at all.</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		base, npages := p.chunkOf(ci).findScavengeCandidate(searchIdx, minPages, maxPages)
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		<span class="comment">// If we found something, scavenge it and return!</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		if npages != 0 {
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>			<span class="comment">// Compute the full address for the start of the range.</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>			addr := chunkBase(ci) + uintptr(base)*pageSize
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>			<span class="comment">// Mark the range we&#39;re about to scavenge as allocated, because</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>			<span class="comment">// we don&#39;t want any allocating goroutines to grab it while</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>			<span class="comment">// the scavenging is in progress. Be careful here -- just do the</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>			<span class="comment">// bare minimum to avoid stepping on our own scavenging stats.</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>			p.chunkOf(ci).allocRange(base, npages)
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>			p.update(addr, uintptr(npages), true, true)
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>			<span class="comment">// With that done, it&#39;s safe to unlock.</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>			unlock(p.mheapLock)
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>			if !p.test {
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>				pageTraceScav(getg().m.p.ptr(), 0, addr, uintptr(npages))
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>				<span class="comment">// Only perform sys* operations if we&#39;re not in a test.</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>				<span class="comment">// It&#39;s dangerous to do so otherwise.</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>				sysUnused(unsafe.Pointer(addr), uintptr(npages)*pageSize)
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>				<span class="comment">// Update global accounting only when not in test, otherwise</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>				<span class="comment">// the runtime&#39;s accounting will be wrong.</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>				nbytes := int64(npages * pageSize)
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>				gcController.heapReleased.add(nbytes)
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>				gcController.heapFree.add(-nbytes)
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>				stats := memstats.heapStats.acquire()
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>				atomic.Xaddint64(&amp;stats.committed, -nbytes)
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>				atomic.Xaddint64(&amp;stats.released, nbytes)
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>				memstats.heapStats.release()
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>			}
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>			<span class="comment">// Relock the heap, because now we need to make these pages</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>			<span class="comment">// available allocation. Free them back to the page allocator.</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>			lock(p.mheapLock)
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>			if b := (offAddr{addr}); b.lessThan(p.searchAddr) {
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>				p.searchAddr = b
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>			}
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			p.chunkOf(ci).free(base, npages)
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>			p.update(addr, uintptr(npages), true, false)
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>			<span class="comment">// Mark the range as scavenged.</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>			p.chunkOf(ci).scavenged.setRange(base, npages)
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>			unlock(p.mheapLock)
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>			return uintptr(npages) * pageSize
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		}
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	}
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	<span class="comment">// Mark this chunk as having no free pages.</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	p.scav.index.setEmpty(ci)
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	unlock(p.mheapLock)
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	return 0
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>}
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span><span class="comment">// fillAligned returns x but with all zeroes in m-aligned</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span><span class="comment">// groups of m bits set to 1 if any bit in the group is non-zero.</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span><span class="comment">// For example, fillAligned(0x0100a3, 8) == 0xff00ff.</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span><span class="comment">// Note that if m == 1, this is a no-op.</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span><span class="comment">// m must be a power of 2 &lt;= maxPagesPerPhysPage.</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>func fillAligned(x uint64, m uint) uint64 {
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	apply := func(x uint64, c uint64) uint64 {
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		<span class="comment">// The technique used it here is derived from</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>		<span class="comment">// https://graphics.stanford.edu/~seander/bithacks.html#ZeroInWord</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		<span class="comment">// and extended for more than just bytes (like nibbles</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		<span class="comment">// and uint16s) by using an appropriate constant.</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		<span class="comment">// To summarize the technique, quoting from that page:</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>		<span class="comment">// &#34;[It] works by first zeroing the high bits of the [8]</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		<span class="comment">// bytes in the word. Subsequently, it adds a number that</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		<span class="comment">// will result in an overflow to the high bit of a byte if</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		<span class="comment">// any of the low bits were initially set. Next the high</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		<span class="comment">// bits of the original word are ORed with these values;</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		<span class="comment">// thus, the high bit of a byte is set iff any bit in the</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		<span class="comment">// byte was set. Finally, we determine if any of these high</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		<span class="comment">// bits are zero by ORing with ones everywhere except the</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		<span class="comment">// high bits and inverting the result.&#34;</span>
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		return ^((((x &amp; c) + c) | x) | c)
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>	<span class="comment">// Transform x to contain a 1 bit at the top of each m-aligned</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	<span class="comment">// group of m zero bits.</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>	switch m {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	case 1:
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>		return x
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	case 2:
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		x = apply(x, 0x5555555555555555)
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	case 4:
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		x = apply(x, 0x7777777777777777)
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	case 8:
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		x = apply(x, 0x7f7f7f7f7f7f7f7f)
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	case 16:
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		x = apply(x, 0x7fff7fff7fff7fff)
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	case 32:
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		x = apply(x, 0x7fffffff7fffffff)
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	case 64: <span class="comment">// == maxPagesPerPhysPage</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>		x = apply(x, 0x7fffffffffffffff)
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	default:
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		throw(&#34;bad m value&#34;)
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	}
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	<span class="comment">// Now, the top bit of each m-aligned group in x is set</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>	<span class="comment">// that group was all zero in the original x.</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	<span class="comment">// From each group of m bits subtract 1.</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	<span class="comment">// Because we know only the top bits of each</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	<span class="comment">// m-aligned group are set, we know this will</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	<span class="comment">// set each group to have all the bits set except</span>
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	<span class="comment">// the top bit, so just OR with the original</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	<span class="comment">// result to set all the bits.</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	return ^((x - (x &gt;&gt; (m - 1))) | x)
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>}
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span><span class="comment">// findScavengeCandidate returns a start index and a size for this pallocData</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span><span class="comment">// segment which represents a contiguous region of free and unscavenged memory.</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span><span class="comment">// searchIdx indicates the page index within this chunk to start the search, but</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span><span class="comment">// note that findScavengeCandidate searches backwards through the pallocData. As</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span><span class="comment">// a result, it will return the highest scavenge candidate in address order.</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span><span class="comment">// min indicates a hard minimum size and alignment for runs of pages. That is,</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span><span class="comment">// findScavengeCandidate will not return a region smaller than min pages in size,</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span><span class="comment">// or that is min pages or greater in size but not aligned to min. min must be</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span><span class="comment">// a non-zero power of 2 &lt;= maxPagesPerPhysPage.</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span><span class="comment">// max is a hint for how big of a region is desired. If max &gt;= pallocChunkPages, then</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span><span class="comment">// findScavengeCandidate effectively returns entire free and unscavenged regions.</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span><span class="comment">// If max &lt; pallocChunkPages, it may truncate the returned region such that size is</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span><span class="comment">// max. However, findScavengeCandidate may still return a larger region if, for</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span><span class="comment">// example, it chooses to preserve huge pages, or if max is not aligned to min (it</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span><span class="comment">// will round up). That is, even if max is small, the returned size is not guaranteed</span>
<span id="L894" class="ln">   894&nbsp;&nbsp;</span><span class="comment">// to be equal to max. max is allowed to be less than min, in which case it is as if</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span><span class="comment">// max == min.</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>func (m *pallocData) findScavengeCandidate(searchIdx uint, minimum, max uintptr) (uint, uint) {
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	if minimum&amp;(minimum-1) != 0 || minimum == 0 {
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		print(&#34;runtime: min = &#34;, minimum, &#34;\n&#34;)
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>		throw(&#34;min must be a non-zero power of 2&#34;)
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>	} else if minimum &gt; maxPagesPerPhysPage {
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		print(&#34;runtime: min = &#34;, minimum, &#34;\n&#34;)
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>		throw(&#34;min too large&#34;)
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>	}
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	<span class="comment">// max may not be min-aligned, so we might accidentally truncate to</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	<span class="comment">// a max value which causes us to return a non-min-aligned value.</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	<span class="comment">// To prevent this, align max up to a multiple of min (which is always</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	<span class="comment">// a power of 2). This also prevents max from ever being less than</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	<span class="comment">// min, unless it&#39;s zero, so handle that explicitly.</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	if max == 0 {
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>		max = minimum
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>	} else {
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>		max = alignUp(max, minimum)
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	}
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	i := int(searchIdx / 64)
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>	<span class="comment">// Start by quickly skipping over blocks of non-free or scavenged pages.</span>
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>	for ; i &gt;= 0; i-- {
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		<span class="comment">// 1s are scavenged OR non-free =&gt; 0s are unscavenged AND free</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>		x := fillAligned(m.scavenged[i]|m.pallocBits[i], uint(minimum))
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>		if x != ^uint64(0) {
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>			break
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>		}
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	}
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>		<span class="comment">// Failed to find any free/unscavenged pages.</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>		return 0, 0
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	}
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	<span class="comment">// We have something in the 64-bit chunk at i, but it could</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>	<span class="comment">// extend further. Loop until we find the extent of it.</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>	<span class="comment">// 1s are scavenged OR non-free =&gt; 0s are unscavenged AND free</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>	x := fillAligned(m.scavenged[i]|m.pallocBits[i], uint(minimum))
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>	z1 := uint(sys.LeadingZeros64(^x))
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	run, end := uint(0), uint(i)*64+(64-z1)
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>	if x&lt;&lt;z1 != 0 {
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>		<span class="comment">// After shifting out z1 bits, we still have 1s,</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>		<span class="comment">// so the run ends inside this word.</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>		run = uint(sys.LeadingZeros64(x &lt;&lt; z1))
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	} else {
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>		<span class="comment">// After shifting out z1 bits, we have no more 1s.</span>
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		<span class="comment">// This means the run extends to the bottom of the</span>
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		<span class="comment">// word so it may extend into further words.</span>
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>		run = 64 - z1
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>		for j := i - 1; j &gt;= 0; j-- {
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>			x := fillAligned(m.scavenged[j]|m.pallocBits[j], uint(minimum))
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>			run += uint(sys.LeadingZeros64(x))
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>			if x != 0 {
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>				<span class="comment">// The run stopped in this word.</span>
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>				break
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>			}
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>		}
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>	}
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>	<span class="comment">// Split the run we found if it&#39;s larger than max but hold on to</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>	<span class="comment">// our original length, since we may need it later.</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>	size := min(run, uint(max))
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	start := end - size
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>	<span class="comment">// Each huge page is guaranteed to fit in a single palloc chunk.</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Support larger huge page sizes.</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Consider taking pages-per-huge-page as a parameter</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>	<span class="comment">// so we can write tests for this.</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	if physHugePageSize &gt; pageSize &amp;&amp; physHugePageSize &gt; physPageSize {
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>		<span class="comment">// We have huge pages, so let&#39;s ensure we don&#39;t break one by scavenging</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		<span class="comment">// over a huge page boundary. If the range [start, start+size) overlaps with</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>		<span class="comment">// a free-and-unscavenged huge page, we want to grow the region we scavenge</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>		<span class="comment">// to include that huge page.</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>		<span class="comment">// Compute the huge page boundary above our candidate.</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>		pagesPerHugePage := physHugePageSize / pageSize
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>		hugePageAbove := uint(alignUp(uintptr(start), pagesPerHugePage))
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>		<span class="comment">// If that boundary is within our current candidate, then we may be breaking</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>		<span class="comment">// a huge page.</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>		if hugePageAbove &lt;= end {
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>			<span class="comment">// Compute the huge page boundary below our candidate.</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>			hugePageBelow := uint(alignDown(uintptr(start), pagesPerHugePage))
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>			if hugePageBelow &gt;= end-run {
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>				<span class="comment">// We&#39;re in danger of breaking apart a huge page since start+size crosses</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>				<span class="comment">// a huge page boundary and rounding down start to the nearest huge</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>				<span class="comment">// page boundary is included in the full run we found. Include the entire</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>				<span class="comment">// huge page in the bound by rounding down to the huge page size.</span>
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>				size = size + (start - hugePageBelow)
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>				start = hugePageBelow
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>			}
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>		}
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	}
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	return start, size
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>}
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span><span class="comment">// scavengeIndex is a structure for efficiently managing which pageAlloc chunks have</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span><span class="comment">// memory available to scavenge.</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>type scavengeIndex struct {
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	<span class="comment">// chunks is a scavChunkData-per-chunk structure that indicates the presence of pages</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	<span class="comment">// available for scavenging. Updates to the index are serialized by the pageAlloc lock.</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	<span class="comment">// It tracks chunk occupancy and a generation counter per chunk. If a chunk&#39;s occupancy</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	<span class="comment">// never exceeds pallocChunkDensePages over the course of a single GC cycle, the chunk</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>	<span class="comment">// becomes eligible for scavenging on the next cycle. If a chunk ever hits this density</span>
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>	<span class="comment">// threshold it immediately becomes unavailable for scavenging in the current cycle as</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>	<span class="comment">// well as the next.</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	<span class="comment">// [min, max) represents the range of chunks that is safe to access (i.e. will not cause</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	<span class="comment">// a fault). As an optimization minHeapIdx represents the true minimum chunk that has been</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	<span class="comment">// mapped, since min is likely rounded down to include the system page containing minHeapIdx.</span>
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>	<span class="comment">// For a chunk size of 4 MiB this structure will only use 2 MiB for a 1 TiB contiguous heap.</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>	chunks     []atomicScavChunkData
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>	min, max   atomic.Uintptr
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	minHeapIdx atomic.Uintptr
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	<span class="comment">// searchAddr* is the maximum address (in the offset address space, so we have a linear</span>
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	<span class="comment">// view of the address space; see mranges.go:offAddr) containing memory available to</span>
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	<span class="comment">// scavenge. It is a hint to the find operation to avoid O(n^2) behavior in repeated lookups.</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	<span class="comment">// searchAddr* is always inclusive and should be the base address of the highest runtime</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>	<span class="comment">// page available for scavenging.</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	<span class="comment">// searchAddrForce is managed by find and free.</span>
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	<span class="comment">// searchAddrBg is managed by find and nextGen.</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>	<span class="comment">// Normally, find monotonically decreases searchAddr* as it finds no more free pages to</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>	<span class="comment">// scavenge. However, mark, when marking a new chunk at an index greater than the current</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	<span class="comment">// searchAddr, sets searchAddr to the *negative* index into chunks of that page. The trick here</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>	<span class="comment">// is that concurrent calls to find will fail to monotonically decrease searchAddr*, and so they</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	<span class="comment">// won&#39;t barge over new memory becoming available to scavenge. Furthermore, this ensures</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	<span class="comment">// that some future caller of find *must* observe the new high index. That caller</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>	<span class="comment">// (or any other racing with it), then makes searchAddr positive before continuing, bringing</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	<span class="comment">// us back to our monotonically decreasing steady-state.</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>	<span class="comment">// A pageAlloc lock serializes updates between min, max, and searchAddr, so abs(searchAddr)</span>
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>	<span class="comment">// is always guaranteed to be &gt;= min and &lt; max (converted to heap addresses).</span>
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>	<span class="comment">// searchAddrBg is increased only on each new generation and is mainly used by the</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>	<span class="comment">// background scavenger and heap-growth scavenging. searchAddrForce is increased continuously</span>
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	<span class="comment">// as memory gets freed and is mainly used by eager memory reclaim such as debug.FreeOSMemory</span>
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>	<span class="comment">// and scavenging to maintain the memory limit.</span>
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>	searchAddrBg    atomicOffAddr
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>	searchAddrForce atomicOffAddr
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>	<span class="comment">// freeHWM is the highest address (in offset address space) that was freed</span>
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>	<span class="comment">// this generation.</span>
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>	freeHWM offAddr
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>	<span class="comment">// Generation counter. Updated by nextGen at the end of each mark phase.</span>
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>	gen uint32
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>	<span class="comment">// test indicates whether or not we&#39;re in a test.</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>	test bool
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>}
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span><span class="comment">// init initializes the scavengeIndex.</span>
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span><span class="comment">// Returns the amount added to sysStat.</span>
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>func (s *scavengeIndex) init(test bool, sysStat *sysMemStat) uintptr {
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>	s.searchAddrBg.Clear()
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>	s.searchAddrForce.Clear()
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>	s.freeHWM = minOffAddr
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>	s.test = test
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>	return s.sysInit(test, sysStat)
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>}
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span><span class="comment">// sysGrow updates the index&#39;s backing store in response to a heap growth.</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span><span class="comment">// Returns the amount of memory added to sysStat.</span>
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>func (s *scavengeIndex) grow(base, limit uintptr, sysStat *sysMemStat) uintptr {
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>	<span class="comment">// Update minHeapIdx. Note that even if there&#39;s no mapping work to do,</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	<span class="comment">// we may still have a new, lower minimum heap address.</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	minHeapIdx := s.minHeapIdx.Load()
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>	if baseIdx := uintptr(chunkIndex(base)); minHeapIdx == 0 || baseIdx &lt; minHeapIdx {
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>		s.minHeapIdx.Store(baseIdx)
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	}
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	return s.sysGrow(base, limit, sysStat)
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>}
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span><span class="comment">// find returns the highest chunk index that may contain pages available to scavenge.</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span><span class="comment">// It also returns an offset to start searching in the highest chunk.</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>func (s *scavengeIndex) find(force bool) (chunkIdx, uint) {
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>	cursor := &amp;s.searchAddrBg
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>	if force {
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>		cursor = &amp;s.searchAddrForce
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>	}
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>	searchAddr, marked := cursor.Load()
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>	if searchAddr == minOffAddr.addr() {
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>		<span class="comment">// We got a cleared search addr.</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>		return 0, 0
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>	}
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>	<span class="comment">// Starting from searchAddr&#39;s chunk, iterate until we find a chunk with pages to scavenge.</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>	gen := s.gen
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>	min := chunkIdx(s.minHeapIdx.Load())
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>	start := chunkIndex(searchAddr)
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>	<span class="comment">// N.B. We&#39;ll never map the 0&#39;th chunk, so minHeapIdx ensures this loop overflow.</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>	for i := start; i &gt;= min; i-- {
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>		<span class="comment">// Skip over chunks.</span>
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>		if !s.chunks[i].load().shouldScavenge(gen, force) {
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>			continue
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>		}
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re still scavenging this chunk.</span>
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>		if i == start {
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>			return i, chunkPageIndex(searchAddr)
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>		}
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>		<span class="comment">// Try to reduce searchAddr to newSearchAddr.</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>		newSearchAddr := chunkBase(i) + pallocChunkBytes - pageSize
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>		if marked {
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>			<span class="comment">// Attempt to be the first one to decrease the searchAddr</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>			<span class="comment">// after an increase. If we fail, that means there was another</span>
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>			<span class="comment">// increase, or somebody else got to it before us. Either way,</span>
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>			<span class="comment">// it doesn&#39;t matter. We may lose some performance having an</span>
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>			<span class="comment">// incorrect search address, but it&#39;s far more important that</span>
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>			<span class="comment">// we don&#39;t miss updates.</span>
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>			cursor.StoreUnmark(searchAddr, newSearchAddr)
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>		} else {
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>			<span class="comment">// Decrease searchAddr.</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>			cursor.StoreMin(newSearchAddr)
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>		}
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>		return i, pallocChunkPages - 1
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>	}
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>	<span class="comment">// Clear searchAddr, because we&#39;ve exhausted the heap.</span>
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	cursor.Clear()
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	return 0, 0
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>}
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span><span class="comment">// alloc updates metadata for chunk at index ci with the fact that</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span><span class="comment">// an allocation of npages occurred. It also eagerly attempts to collapse</span>
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span><span class="comment">// the chunk&#39;s memory into hugepage if the chunk has become sufficiently</span>
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span><span class="comment">// dense and we&#39;re not allocating the whole chunk at once (which suggests</span>
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span><span class="comment">// the allocation is part of a bigger one and it&#39;s probably not worth</span>
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span><span class="comment">// eagerly collapsing).</span>
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span><span class="comment">// alloc may only run concurrently with find.</span>
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>func (s *scavengeIndex) alloc(ci chunkIdx, npages uint) {
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>	sc := s.chunks[ci].load()
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>	sc.alloc(npages, s.gen)
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Consider eagerly backing memory with huge pages</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>	<span class="comment">// here and track whether we believe this chunk is backed by huge pages.</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>	<span class="comment">// In the past we&#39;ve attempted to use sysHugePageCollapse (which uses</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>	<span class="comment">// MADV_COLLAPSE on Linux, and is unsupported elswhere) for this purpose,</span>
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	<span class="comment">// but that caused performance issues in production environments.</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	s.chunks[ci].store(sc)
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>}
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span><span class="comment">// free updates metadata for chunk at index ci with the fact that</span>
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span><span class="comment">// a free of npages occurred.</span>
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span><span class="comment">// free may only run concurrently with find.</span>
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>func (s *scavengeIndex) free(ci chunkIdx, page, npages uint) {
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>	sc := s.chunks[ci].load()
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>	sc.free(npages, s.gen)
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>	s.chunks[ci].store(sc)
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>	<span class="comment">// Update scavenge search addresses.</span>
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>	addr := chunkBase(ci) + uintptr(page+npages-1)*pageSize
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>	if s.freeHWM.lessThan(offAddr{addr}) {
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>		s.freeHWM = offAddr{addr}
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>	}
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>	<span class="comment">// N.B. Because free is serialized, it&#39;s not necessary to do a</span>
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>	<span class="comment">// full CAS here. free only ever increases searchAddr, while</span>
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>	<span class="comment">// find only ever decreases it. Since we only ever race with</span>
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>	<span class="comment">// decreases, even if the value we loaded is stale, the actual</span>
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>	<span class="comment">// value will never be larger.</span>
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>	searchAddr, _ := s.searchAddrForce.Load()
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>	if (offAddr{searchAddr}).lessThan(offAddr{addr}) {
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>		s.searchAddrForce.StoreMarked(addr)
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>	}
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>}
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span><span class="comment">// nextGen moves the scavenger forward one generation. Must be called</span>
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span><span class="comment">// once per GC cycle, but may be called more often to force more memory</span>
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span><span class="comment">// to be released.</span>
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span><span class="comment">// nextGen may only run concurrently with find.</span>
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>func (s *scavengeIndex) nextGen() {
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>	s.gen++
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>	searchAddr, _ := s.searchAddrBg.Load()
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>	if (offAddr{searchAddr}).lessThan(s.freeHWM) {
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>		s.searchAddrBg.StoreMarked(s.freeHWM.addr())
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>	}
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>	s.freeHWM = minOffAddr
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>}
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span><span class="comment">// setEmpty marks that the scavenger has finished looking at ci</span>
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span><span class="comment">// for now to prevent the scavenger from getting stuck looking</span>
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span><span class="comment">// at the same chunk.</span>
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span><span class="comment">// setEmpty may only run concurrently with find.</span>
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>func (s *scavengeIndex) setEmpty(ci chunkIdx) {
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>	val := s.chunks[ci].load()
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>	val.setEmpty()
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>	s.chunks[ci].store(val)
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>}
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span><span class="comment">// atomicScavChunkData is an atomic wrapper around a scavChunkData</span>
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span><span class="comment">// that stores it in its packed form.</span>
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>type atomicScavChunkData struct {
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>	value atomic.Uint64
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>}
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span><span class="comment">// load loads and unpacks a scavChunkData.</span>
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>func (sc *atomicScavChunkData) load() scavChunkData {
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>	return unpackScavChunkData(sc.value.Load())
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>}
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span><span class="comment">// store packs and writes a new scavChunkData. store must be serialized</span>
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span><span class="comment">// with other calls to store.</span>
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>func (sc *atomicScavChunkData) store(ssc scavChunkData) {
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>	sc.value.Store(ssc.pack())
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>}
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span><span class="comment">// scavChunkData tracks information about a palloc chunk for</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span><span class="comment">// scavenging. It packs well into 64 bits.</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span><span class="comment">// The zero value always represents a valid newly-grown chunk.</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>type scavChunkData struct {
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>	<span class="comment">// inUse indicates how many pages in this chunk are currently</span>
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>	<span class="comment">// allocated.</span>
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>	<span class="comment">// Only the first 10 bits are used.</span>
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>	inUse uint16
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>	<span class="comment">// lastInUse indicates how many pages in this chunk were allocated</span>
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>	<span class="comment">// when we transitioned from gen-1 to gen.</span>
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>	<span class="comment">// Only the first 10 bits are used.</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>	lastInUse uint16
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>	<span class="comment">// gen is the generation counter from a scavengeIndex from the</span>
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>	<span class="comment">// last time this scavChunkData was updated.</span>
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	gen uint32
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>	<span class="comment">// scavChunkFlags represents additional flags</span>
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>	<span class="comment">// Note: only 6 bits are available.</span>
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>	scavChunkFlags
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>}
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span><span class="comment">// unpackScavChunkData unpacks a scavChunkData from a uint64.</span>
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>func unpackScavChunkData(sc uint64) scavChunkData {
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>	return scavChunkData{
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>		inUse:          uint16(sc),
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>		lastInUse:      uint16(sc&gt;&gt;16) &amp; scavChunkInUseMask,
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>		gen:            uint32(sc &gt;&gt; 32),
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>		scavChunkFlags: scavChunkFlags(uint8(sc&gt;&gt;(16+logScavChunkInUseMax)) &amp; scavChunkFlagsMask),
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>	}
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>}
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span><span class="comment">// pack returns sc packed into a uint64.</span>
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>func (sc scavChunkData) pack() uint64 {
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>	return uint64(sc.inUse) |
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>		(uint64(sc.lastInUse) &lt;&lt; 16) |
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>		(uint64(sc.scavChunkFlags) &lt;&lt; (16 + logScavChunkInUseMax)) |
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>		(uint64(sc.gen) &lt;&lt; 32)
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>}
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>const (
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>	<span class="comment">// scavChunkHasFree indicates whether the chunk has anything left to</span>
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>	<span class="comment">// scavenge. This is the opposite of &#34;empty,&#34; used elsewhere in this</span>
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>	<span class="comment">// file. The reason we say &#34;HasFree&#34; here is so the zero value is</span>
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>	<span class="comment">// correct for a newly-grown chunk. (New memory is scavenged.)</span>
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>	scavChunkHasFree scavChunkFlags = 1 &lt;&lt; iota
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>	<span class="comment">// scavChunkMaxFlags is the maximum number of flags we can have, given how</span>
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>	<span class="comment">// a scavChunkData is packed into 8 bytes.</span>
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>	scavChunkMaxFlags  = 6
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>	scavChunkFlagsMask = (1 &lt;&lt; scavChunkMaxFlags) - 1
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>	<span class="comment">// logScavChunkInUseMax is the number of bits needed to represent the number</span>
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	<span class="comment">// of pages allocated in a single chunk. This is 1 more than log2 of the</span>
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	<span class="comment">// number of pages in the chunk because we need to represent a fully-allocated</span>
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>	<span class="comment">// chunk.</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	logScavChunkInUseMax = logPallocChunkPages + 1
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>	scavChunkInUseMask   = (1 &lt;&lt; logScavChunkInUseMax) - 1
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>)
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span><span class="comment">// scavChunkFlags is a set of bit-flags for the scavenger for each palloc chunk.</span>
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>type scavChunkFlags uint8
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span><span class="comment">// isEmpty returns true if the hasFree flag is unset.</span>
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>func (sc *scavChunkFlags) isEmpty() bool {
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>	return (*sc)&amp;scavChunkHasFree == 0
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>}
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span><span class="comment">// setEmpty clears the hasFree flag.</span>
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>func (sc *scavChunkFlags) setEmpty() {
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>	*sc &amp;^= scavChunkHasFree
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>}
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span><span class="comment">// setNonEmpty sets the hasFree flag.</span>
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>func (sc *scavChunkFlags) setNonEmpty() {
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>	*sc |= scavChunkHasFree
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>}
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span><span class="comment">// shouldScavenge returns true if the corresponding chunk should be interrogated</span>
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span><span class="comment">// by the scavenger.</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>func (sc scavChunkData) shouldScavenge(currGen uint32, force bool) bool {
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>	if sc.isEmpty() {
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>		<span class="comment">// Nothing to scavenge.</span>
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>		return false
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>	}
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>	if force {
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re forcing the memory to be scavenged.</span>
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>		return true
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>	}
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>	if sc.gen == currGen {
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>		<span class="comment">// In the current generation, if either the current or last generation</span>
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>		<span class="comment">// is dense, then skip scavenging. Inverting that, we should scavenge</span>
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>		<span class="comment">// if both the current and last generation were not dense.</span>
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>		return sc.inUse &lt; scavChunkHiOccPages &amp;&amp; sc.lastInUse &lt; scavChunkHiOccPages
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>	}
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;re one or more generations ahead, we know inUse represents the current</span>
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>	<span class="comment">// state of the chunk, since otherwise it would&#39;ve been updated already.</span>
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>	return sc.inUse &lt; scavChunkHiOccPages
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>}
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span><span class="comment">// alloc updates sc given that npages were allocated in the corresponding chunk.</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>func (sc *scavChunkData) alloc(npages uint, newGen uint32) {
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>	if uint(sc.inUse)+npages &gt; pallocChunkPages {
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>		print(&#34;runtime: inUse=&#34;, sc.inUse, &#34; npages=&#34;, npages, &#34;\n&#34;)
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>		throw(&#34;too many pages allocated in chunk?&#34;)
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>	}
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>	if sc.gen != newGen {
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>		sc.lastInUse = sc.inUse
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>		sc.gen = newGen
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>	}
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>	sc.inUse += uint16(npages)
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>	if sc.inUse == pallocChunkPages {
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>		<span class="comment">// There&#39;s nothing for the scavenger to take from here.</span>
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>		sc.setEmpty()
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>	}
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>}
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span><span class="comment">// free updates sc given that npages was freed in the corresponding chunk.</span>
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>func (sc *scavChunkData) free(npages uint, newGen uint32) {
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>	if uint(sc.inUse) &lt; npages {
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>		print(&#34;runtime: inUse=&#34;, sc.inUse, &#34; npages=&#34;, npages, &#34;\n&#34;)
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>		throw(&#34;allocated pages below zero?&#34;)
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>	}
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>	if sc.gen != newGen {
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>		sc.lastInUse = sc.inUse
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>		sc.gen = newGen
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>	}
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>	sc.inUse -= uint16(npages)
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>	<span class="comment">// The scavenger can no longer be done with this chunk now that</span>
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>	<span class="comment">// new memory has been freed into it.</span>
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>	sc.setNonEmpty()
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>}
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>type piController struct {
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>	kp float64 <span class="comment">// Proportional constant.</span>
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>	ti float64 <span class="comment">// Integral time constant.</span>
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>	tt float64 <span class="comment">// Reset time.</span>
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>	min, max float64 <span class="comment">// Output boundaries.</span>
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>	<span class="comment">// PI controller state.</span>
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>	errIntegral float64 <span class="comment">// Integral of the error from t=0 to now.</span>
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>	<span class="comment">// Error flags.</span>
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>	errOverflow   bool <span class="comment">// Set if errIntegral ever overflowed.</span>
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>	inputOverflow bool <span class="comment">// Set if an operation with the input overflowed.</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>}
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span><span class="comment">// next provides a new sample to the controller.</span>
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span><span class="comment">// input is the sample, setpoint is the desired point, and period is how much</span>
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span><span class="comment">// time (in whatever unit makes the most sense) has passed since the last sample.</span>
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span><span class="comment">// Returns a new value for the variable it&#39;s controlling, and whether the operation</span>
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span><span class="comment">// completed successfully. One reason this might fail is if error has been growing</span>
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span><span class="comment">// in an unbounded manner, to the point of overflow.</span>
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span><span class="comment">// In the specific case of an error overflow occurs, the errOverflow field will be</span>
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span><span class="comment">// set and the rest of the controller&#39;s internal state will be fully reset.</span>
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>func (c *piController) next(input, setpoint, period float64) (float64, bool) {
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>	<span class="comment">// Compute the raw output value.</span>
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>	prop := c.kp * (setpoint - input)
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>	rawOutput := prop + c.errIntegral
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>	<span class="comment">// Clamp rawOutput into output.</span>
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>	output := rawOutput
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>	if isInf(output) || isNaN(output) {
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>		<span class="comment">// The input had a large enough magnitude that either it was already</span>
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>		<span class="comment">// overflowed, or some operation with it overflowed.</span>
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>		<span class="comment">// Set a flag and reset. That&#39;s the safest thing to do.</span>
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>		c.reset()
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>		c.inputOverflow = true
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>		return c.min, false
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>	}
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>	if output &lt; c.min {
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>		output = c.min
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>	} else if output &gt; c.max {
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>		output = c.max
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>	}
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>	<span class="comment">// Update the controller&#39;s state.</span>
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>	if c.ti != 0 &amp;&amp; c.tt != 0 {
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>		c.errIntegral += (c.kp*period/c.ti)*(setpoint-input) + (period/c.tt)*(output-rawOutput)
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>		if isInf(c.errIntegral) || isNaN(c.errIntegral) {
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>			<span class="comment">// So much error has accumulated that we managed to overflow.</span>
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>			<span class="comment">// The assumptions around the controller have likely broken down.</span>
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>			<span class="comment">// Set a flag and reset. That&#39;s the safest thing to do.</span>
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>			c.reset()
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>			c.errOverflow = true
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>			return c.min, false
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>		}
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>	}
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>	return output, true
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>}
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span><span class="comment">// reset resets the controller state, except for controller error flags.</span>
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>func (c *piController) reset() {
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>	c.errIntegral = 0
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>}
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>
</pre><p><a href="mgcscavenge.go?m=text">View as plain text</a></p>

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
