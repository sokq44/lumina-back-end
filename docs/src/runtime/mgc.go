<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mgc.go - Go Documentation Server</title>

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
<a href="mgc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mgc.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Garbage collector (GC).</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// The GC runs concurrently with mutator threads, is type accurate (aka precise), allows multiple</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// GC thread to run in parallel. It is a concurrent mark and sweep that uses a write barrier. It is</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// non-generational and non-compacting. Allocation is done using size segregated per P allocation</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// areas to minimize fragmentation while eliminating locks in the common case.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// The algorithm decomposes into several steps.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// This is a high level description of the algorithm being used. For an overview of GC a good</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// place to start is Richard Jones&#39; gchandbook.org.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// The algorithm&#39;s intellectual heritage includes Dijkstra&#39;s on-the-fly algorithm, see</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Edsger W. Dijkstra, Leslie Lamport, A. J. Martin, C. S. Scholten, and E. F. M. Steffens. 1978.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// On-the-fly garbage collection: an exercise in cooperation. Commun. ACM 21, 11 (November 1978),</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// 966-975.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// For journal quality proofs that these steps are complete, correct, and terminate see</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// Hudson, R., and Moss, J.E.B. Copying Garbage Collection without stopping the world.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// Concurrency and Computation: Practice and Experience 15(3-5), 2003.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// 1. GC performs sweep termination.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//    a. Stop the world. This causes all Ps to reach a GC safe-point.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//    b. Sweep any unswept spans. There will only be unswept spans if</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//    this GC cycle was forced before the expected time.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// 2. GC performs the mark phase.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//    a. Prepare for the mark phase by setting gcphase to _GCmark</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//    (from _GCoff), enabling the write barrier, enabling mutator</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//    assists, and enqueueing root mark jobs. No objects may be</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//    scanned until all Ps have enabled the write barrier, which is</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//    accomplished using STW.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//    b. Start the world. From this point, GC work is done by mark</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//    workers started by the scheduler and by assists performed as</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//    part of allocation. The write barrier shades both the</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//    overwritten pointer and the new pointer value for any pointer</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//    writes (see mbarrier.go for details). Newly allocated objects</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//    are immediately marked black.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//    c. GC performs root marking jobs. This includes scanning all</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//    stacks, shading all globals, and shading any heap pointers in</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//    off-heap runtime data structures. Scanning a stack stops a</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//    goroutine, shades any pointers found on its stack, and then</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//    resumes the goroutine.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//    d. GC drains the work queue of grey objects, scanning each grey</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//    object to black and shading all pointers found in the object</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//    (which in turn may add those pointers to the work queue).</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//    e. Because GC work is spread across local caches, GC uses a</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//    distributed termination algorithm to detect when there are no</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//    more root marking jobs or grey objects (see gcMarkDone). At this</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//    point, GC transitions to mark termination.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// 3. GC performs mark termination.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//    a. Stop the world.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//    b. Set gcphase to _GCmarktermination, and disable workers and</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">//    assists.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//    c. Perform housekeeping like flushing mcaches.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// 4. GC performs the sweep phase.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//    a. Prepare for the sweep phase by setting gcphase to _GCoff,</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//    setting up sweep state and disabling the write barrier.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//    b. Start the world. From this point on, newly allocated objects</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//    are white, and allocating sweeps spans before use if necessary.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//    c. GC does concurrent sweeping in the background and in response</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//    to allocation. See description below.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// 5. When sufficient allocation has taken place, replay the sequence</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// starting with 1 above. See discussion of GC rate below.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// Concurrent sweep.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// The sweep phase proceeds concurrently with normal program execution.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// The heap is swept span-by-span both lazily (when a goroutine needs another span)</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// and concurrently in a background goroutine (this helps programs that are not CPU bound).</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// At the end of STW mark termination all spans are marked as &#34;needs sweeping&#34;.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// The background sweeper goroutine simply sweeps spans one-by-one.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// To avoid requesting more OS memory while there are unswept spans, when a</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// goroutine needs another span, it first attempts to reclaim that much memory</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// by sweeping. When a goroutine needs to allocate a new small-object span, it</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// sweeps small-object spans for the same object size until it frees at least</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// one object. When a goroutine needs to allocate large-object span from heap,</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// it sweeps spans until it frees at least that many pages into heap. There is</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// one case where this may not suffice: if a goroutine sweeps and frees two</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// nonadjacent one-page spans to the heap, it will allocate a new two-page</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// span, but there can still be other one-page unswept spans which could be</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// combined into a two-page span.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// It&#39;s critical to ensure that no operations proceed on unswept spans (that would corrupt</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// mark bits in GC bitmap). During GC all mcaches are flushed into the central cache,</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// so they are empty. When a goroutine grabs a new span into mcache, it sweeps it.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// When a goroutine explicitly frees an object or sets a finalizer, it ensures that</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// the span is swept (either by sweeping it, or by waiting for the concurrent sweep to finish).</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// The finalizer goroutine is kicked off only when all spans are swept.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// When the next GC starts, it sweeps all not-yet-swept spans (if any).</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// GC rate.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// Next GC is after we&#39;ve allocated an extra amount of memory proportional to</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">// the amount already in use. The proportion is controlled by GOGC environment variable</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">// (100 by default). If GOGC=100 and we&#39;re using 4M, we&#39;ll GC again when we get to 8M</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// (this mark is computed by the gcController.heapGoal method). This keeps the GC cost in</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// linear proportion to the allocation cost. Adjusting GOGC just changes the linear constant</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// (and also the amount of extra memory used).</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// Oblets</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// In order to prevent long pauses while scanning large objects and to</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// improve parallelism, the garbage collector breaks up scan jobs for</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// objects larger than maxObletBytes into &#34;oblets&#34; of at most</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// maxObletBytes. When scanning encounters the beginning of a large</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// object, it scans only the first oblet and enqueues the remaining</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// oblets as new scan jobs.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>package runtime
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>import (
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	&#34;internal/cpu&#34;
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>const (
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	_DebugGC      = 0
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	_FinBlockSize = 4 * 1024
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// concurrentSweep is a debug flag. Disabling this flag</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// ensures all spans are swept while the world is stopped.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	concurrentSweep = true
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// debugScanConservative enables debug logging for stack</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// frames that are scanned conservatively.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	debugScanConservative = false
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// sweepMinHeapDistance is a lower bound on the heap distance</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// (in bytes) reserved for concurrent sweeping between GC</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// cycles.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	sweepMinHeapDistance = 1024 * 1024
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// heapObjectsCanMove always returns false in the current garbage collector.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// It exists for go4.org/unsafe/assume-no-moving-gc, which is an</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// unfortunate idea that had an even more unfortunate implementation.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">// Every time a new Go release happened, the package stopped building,</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// and the authors had to add a new file with a new //go:build line, and</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// then the entire ecosystem of packages with that as a dependency had to</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">// explicitly update to the new version. Many packages depend on</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// assume-no-moving-gc transitively, through paths like</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// inet.af/netaddr -&gt; go4.org/intern -&gt; assume-no-moving-gc.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// This was causing a significant amount of friction around each new</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">// release, so we added this bool for the package to //go:linkname</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// instead. The bool is still unfortunate, but it&#39;s not as bad as</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// breaking the ecosystem on every new release.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// If the Go garbage collector ever does move heap objects, we can set</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">// this to true to break all the programs using assume-no-moving-gc.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">//go:linkname heapObjectsCanMove</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>func heapObjectsCanMove() bool {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	return false
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>func gcinit() {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	if unsafe.Sizeof(workbuf{}) != _WorkbufSize {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		throw(&#34;size of Workbuf is suboptimal&#34;)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// No sweep on the first cycle.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	sweep.active.state.Store(sweepDrainedMask)
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// Initialize GC pacer state.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// Use the environment variable GOGC for the initial gcPercent value.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// Use the environment variable GOMEMLIMIT for the initial memoryLimit value.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	gcController.init(readGOGC(), readGOMEMLIMIT())
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	work.startSema = 1
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	work.markDoneSema = 1
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	lockInit(&amp;work.sweepWaiters.lock, lockRankSweepWaiters)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	lockInit(&amp;work.assistQueue.lock, lockRankAssistQueue)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	lockInit(&amp;work.wbufSpans.lock, lockRankWbufSpans)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// gcenable is called after the bulk of the runtime initialization,</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// just before we&#39;re about to start letting user code run.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// It kicks off the background sweeper goroutine, the background</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// scavenger goroutine, and enables GC.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>func gcenable() {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// Kick off sweeping and scavenging.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	c := make(chan int, 2)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	go bgsweep(c)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	go bgscavenge(c)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	&lt;-c
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	&lt;-c
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	memstats.enablegc = true <span class="comment">// now that runtime is initialized, GC is okay</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// Garbage collector phase.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// Indicates to write barrier and synchronization task to perform.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>var gcphase uint32
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// The compiler knows about this variable.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// If you change it, you must change builtin/runtime.go, too.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// If you change the first four bytes, you must also change the write</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// barrier insertion code.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>var writeBarrier struct {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	enabled bool    <span class="comment">// compiler emits a check of this before calling write barrier</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	pad     [3]byte <span class="comment">// compiler uses 32-bit load for &#34;enabled&#34; field</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	alignme uint64  <span class="comment">// guarantee alignment so that compiler can use a 32 or 64-bit load</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// gcBlackenEnabled is 1 if mutator assists and background mark</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// workers are allowed to blacken objects. This must only be set when</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// gcphase == _GCmark.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>var gcBlackenEnabled uint32
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>const (
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	_GCoff             = iota <span class="comment">// GC not running; sweeping in background, write barrier disabled</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	_GCmark                   <span class="comment">// GC marking roots and workbufs: allocate black, write barrier ENABLED</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	_GCmarktermination        <span class="comment">// GC mark termination: allocate black, P&#39;s help GC, write barrier ENABLED</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>func setGCPhase(x uint32) {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	atomic.Store(&amp;gcphase, x)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	writeBarrier.enabled = gcphase == _GCmark || gcphase == _GCmarktermination
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span><span class="comment">// gcMarkWorkerMode represents the mode that a concurrent mark worker</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span><span class="comment">// should operate in.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">// Concurrent marking happens through four different mechanisms. One</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// is mutator assists, which happen in response to allocations and are</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">// not scheduled. The other three are variations in the per-P mark</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// workers and are distinguished by gcMarkWorkerMode.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>type gcMarkWorkerMode int
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>const (
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// gcMarkWorkerNotWorker indicates that the next scheduled G is not</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// starting work and the mode should be ignored.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	gcMarkWorkerNotWorker gcMarkWorkerMode = iota
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	<span class="comment">// gcMarkWorkerDedicatedMode indicates that the P of a mark</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// worker is dedicated to running that mark worker. The mark</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">// worker should run without preemption.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	gcMarkWorkerDedicatedMode
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// gcMarkWorkerFractionalMode indicates that a P is currently</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// running the &#34;fractional&#34; mark worker. The fractional worker</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// is necessary when GOMAXPROCS*gcBackgroundUtilization is not</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// an integer and using only dedicated workers would result in</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// utilization too far from the target of gcBackgroundUtilization.</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	<span class="comment">// The fractional worker should run until it is preempted and</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	<span class="comment">// will be scheduled to pick up the fractional part of</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	<span class="comment">// GOMAXPROCS*gcBackgroundUtilization.</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	gcMarkWorkerFractionalMode
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	<span class="comment">// gcMarkWorkerIdleMode indicates that a P is running the mark</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	<span class="comment">// worker because it has nothing else to do. The idle worker</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	<span class="comment">// should run until it is preempted and account its time</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	<span class="comment">// against gcController.idleMarkTime.</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	gcMarkWorkerIdleMode
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// gcMarkWorkerModeStrings are the strings labels of gcMarkWorkerModes</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span><span class="comment">// to use in execution traces.</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>var gcMarkWorkerModeStrings = [...]string{
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	&#34;Not worker&#34;,
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	&#34;GC (dedicated)&#34;,
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	&#34;GC (fractional)&#34;,
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	&#34;GC (idle)&#34;,
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">// pollFractionalWorkerExit reports whether a fractional mark worker</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span><span class="comment">// should self-preempt. It assumes it is called from the fractional</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">// worker.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>func pollFractionalWorkerExit() bool {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">// This should be kept in sync with the fractional worker</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// scheduler logic in findRunnableGCWorker.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	now := nanotime()
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	delta := now - gcController.markStartTime
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	if delta &lt;= 0 {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		return true
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	p := getg().m.p.ptr()
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	selfTime := p.gcFractionalMarkTime + (now - p.gcMarkWorkerStartTime)
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">// Add some slack to the utilization goal so that the</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">// fractional worker isn&#39;t behind again the instant it exits.</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	return float64(selfTime)/float64(delta) &gt; 1.2*gcController.fractionalUtilizationGoal
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>var work workType
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>type workType struct {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	full  lfstack          <span class="comment">// lock-free list of full blocks workbuf</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	_     cpu.CacheLinePad <span class="comment">// prevents false-sharing between full and empty</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	empty lfstack          <span class="comment">// lock-free list of empty blocks workbuf</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	_     cpu.CacheLinePad <span class="comment">// prevents false-sharing between empty and nproc/nwait</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	wbufSpans struct {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		lock mutex
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		<span class="comment">// free is a list of spans dedicated to workbufs, but</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		<span class="comment">// that don&#39;t currently contain any workbufs.</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		free mSpanList
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		<span class="comment">// busy is a list of all spans containing workbufs on</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		<span class="comment">// one of the workbuf lists.</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		busy mSpanList
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	<span class="comment">// Restore 64-bit alignment on 32-bit.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	_ uint32
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// bytesMarked is the number of bytes marked this cycle. This</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">// includes bytes blackened in scanned objects, noscan objects</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// that go straight to black, and permagrey objects scanned by</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	<span class="comment">// markroot during the concurrent scan phase. This is updated</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	<span class="comment">// atomically during the cycle. Updates may be batched</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">// arbitrarily, since the value is only read at the end of the</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	<span class="comment">// cycle.</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// Because of benign races during marking, this number may not</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// be the exact number of marked bytes, but it should be very</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// close.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// Put this field here because it needs 64-bit atomic access</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// (and thus 8-byte alignment even on 32-bit architectures).</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	bytesMarked uint64
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	markrootNext uint32 <span class="comment">// next markroot job</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	markrootJobs uint32 <span class="comment">// number of markroot jobs</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	nproc  uint32
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	tstart int64
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	nwait  uint32
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	<span class="comment">// Number of roots of various root types. Set by gcMarkRootPrepare.</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	<span class="comment">// nStackRoots == len(stackRoots), but we have nStackRoots for</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	<span class="comment">// consistency.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	nDataRoots, nBSSRoots, nSpanRoots, nStackRoots int
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	<span class="comment">// Base indexes of each root type. Set by gcMarkRootPrepare.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	baseData, baseBSS, baseSpans, baseStacks, baseEnd uint32
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	<span class="comment">// stackRoots is a snapshot of all of the Gs that existed</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	<span class="comment">// before the beginning of concurrent marking. The backing</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	<span class="comment">// store of this must not be modified because it might be</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	<span class="comment">// shared with allgs.</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	stackRoots []*g
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	<span class="comment">// Each type of GC state transition is protected by a lock.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	<span class="comment">// Since multiple threads can simultaneously detect the state</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	<span class="comment">// transition condition, any thread that detects a transition</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	<span class="comment">// condition must acquire the appropriate transition lock,</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	<span class="comment">// re-check the transition condition and return if it no</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	<span class="comment">// longer holds or perform the transition if it does.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// Likewise, any transition must invalidate the transition</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	<span class="comment">// condition before releasing the lock. This ensures that each</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	<span class="comment">// transition is performed by exactly one thread and threads</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	<span class="comment">// that need the transition to happen block until it has</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	<span class="comment">// happened.</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	<span class="comment">// startSema protects the transition from &#34;off&#34; to mark or</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	<span class="comment">// mark termination.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	startSema uint32
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	<span class="comment">// markDoneSema protects transitions from mark to mark termination.</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	markDoneSema uint32
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	bgMarkReady note   <span class="comment">// signal background mark worker has started</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	bgMarkDone  uint32 <span class="comment">// cas to 1 when at a background mark completion point</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	<span class="comment">// Background mark completion signaling</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	<span class="comment">// mode is the concurrency mode of the current GC cycle.</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	mode gcMode
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	<span class="comment">// userForced indicates the current GC cycle was forced by an</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	<span class="comment">// explicit user call.</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	userForced bool
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	<span class="comment">// initialHeapLive is the value of gcController.heapLive at the</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	<span class="comment">// beginning of this GC cycle.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	initialHeapLive uint64
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	<span class="comment">// assistQueue is a queue of assists that are blocked because</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	<span class="comment">// there was neither enough credit to steal or enough work to</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	<span class="comment">// do.</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	assistQueue struct {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		lock mutex
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		q    gQueue
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	<span class="comment">// sweepWaiters is a list of blocked goroutines to wake when</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	<span class="comment">// we transition from mark termination to sweep.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	sweepWaiters struct {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		lock mutex
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		list gList
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	<span class="comment">// cycles is the number of completed GC cycles, where a GC</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	<span class="comment">// cycle is sweep termination, mark, mark termination, and</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	<span class="comment">// sweep. This differs from memstats.numgc, which is</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	<span class="comment">// incremented at mark termination.</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	cycles atomic.Uint32
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	<span class="comment">// Timing/utilization stats for this cycle.</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	stwprocs, maxprocs                 int32
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	tSweepTerm, tMark, tMarkTerm, tEnd int64 <span class="comment">// nanotime() of phase start</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	pauseNS int64 <span class="comment">// total STW time this cycle</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	<span class="comment">// debug.gctrace heap sizes for this cycle.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	heap0, heap1, heap2 uint64
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	<span class="comment">// Cumulative estimated CPU usage.</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	cpuStats
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span><span class="comment">// GC runs a garbage collection and blocks the caller until the</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span><span class="comment">// garbage collection is complete. It may also block the entire</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span><span class="comment">// program.</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>func GC() {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	<span class="comment">// We consider a cycle to be: sweep termination, mark, mark</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	<span class="comment">// termination, and sweep. This function shouldn&#39;t return</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	<span class="comment">// until a full cycle has been completed, from beginning to</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	<span class="comment">// end. Hence, we always want to finish up the current cycle</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	<span class="comment">// and start a new one. That means:</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	<span class="comment">// 1. In sweep termination, mark, or mark termination of cycle</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// N, wait until mark termination N completes and transitions</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	<span class="comment">// to sweep N.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	<span class="comment">// 2. In sweep N, help with sweep N.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	<span class="comment">// At this point we can begin a full cycle N+1.</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	<span class="comment">// 3. Trigger cycle N+1 by starting sweep termination N+1.</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	<span class="comment">// 4. Wait for mark termination N+1 to complete.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	<span class="comment">// 5. Help with sweep N+1 until it&#39;s done.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	<span class="comment">// This all has to be written to deal with the fact that the</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	<span class="comment">// GC may move ahead on its own. For example, when we block</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	<span class="comment">// until mark termination N, we may wake up in cycle N+2.</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	<span class="comment">// Wait until the current sweep termination, mark, and mark</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	<span class="comment">// termination complete.</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	n := work.cycles.Load()
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	gcWaitOnMark(n)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	<span class="comment">// We&#39;re now in sweep N or later. Trigger GC cycle N+1, which</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	<span class="comment">// will first finish sweep N if necessary and then enter sweep</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	<span class="comment">// termination N+1.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	gcStart(gcTrigger{kind: gcTriggerCycle, n: n + 1})
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	<span class="comment">// Wait for mark termination N+1 to complete.</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	gcWaitOnMark(n + 1)
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	<span class="comment">// Finish sweep N+1 before returning. We do this both to</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	<span class="comment">// complete the cycle and because runtime.GC() is often used</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	<span class="comment">// as part of tests and benchmarks to get the system into a</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	<span class="comment">// relatively stable and isolated state.</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	for work.cycles.Load() == n+1 &amp;&amp; sweepone() != ^uintptr(0) {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		Gosched()
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	<span class="comment">// Callers may assume that the heap profile reflects the</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	<span class="comment">// just-completed cycle when this returns (historically this</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	<span class="comment">// happened because this was a STW GC), but right now the</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	<span class="comment">// profile still reflects mark termination N, not N+1.</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	<span class="comment">// As soon as all of the sweep frees from cycle N+1 are done,</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	<span class="comment">// we can go ahead and publish the heap profile.</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	<span class="comment">// First, wait for sweeping to finish. (We know there are no</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	<span class="comment">// more spans on the sweep queue, but we may be concurrently</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	<span class="comment">// sweeping spans, so we have to wait.)</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	for work.cycles.Load() == n+1 &amp;&amp; !isSweepDone() {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		Gosched()
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	<span class="comment">// Now we&#39;re really done with sweeping, so we can publish the</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	<span class="comment">// stable heap profile. Only do this if we haven&#39;t already hit</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	<span class="comment">// another mark termination.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	cycle := work.cycles.Load()
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	if cycle == n+1 || (gcphase == _GCmark &amp;&amp; cycle == n+2) {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		mProf_PostSweep()
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	}
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	releasem(mp)
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span><span class="comment">// gcWaitOnMark blocks until GC finishes the Nth mark phase. If GC has</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span><span class="comment">// already completed this mark phase, it returns immediately.</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>func gcWaitOnMark(n uint32) {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	for {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		<span class="comment">// Disable phase transitions.</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		lock(&amp;work.sweepWaiters.lock)
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		nMarks := work.cycles.Load()
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		if gcphase != _GCmark {
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>			<span class="comment">// We&#39;ve already completed this cycle&#39;s mark.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>			nMarks++
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		}
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		if nMarks &gt; n {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			<span class="comment">// We&#39;re done.</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			unlock(&amp;work.sweepWaiters.lock)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>			return
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		<span class="comment">// Wait until sweep termination, mark, and mark</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		<span class="comment">// termination of cycle N complete.</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		work.sweepWaiters.list.push(getg())
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		goparkunlock(&amp;work.sweepWaiters.lock, waitReasonWaitForGCCycle, traceBlockUntilGCEnds, 1)
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>}
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span><span class="comment">// gcMode indicates how concurrent a GC cycle should be.</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>type gcMode int
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>const (
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	gcBackgroundMode gcMode = iota <span class="comment">// concurrent GC and sweep</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	gcForceMode                    <span class="comment">// stop-the-world GC now, concurrent sweep</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	gcForceBlockMode               <span class="comment">// stop-the-world GC now and STW sweep (forced by user)</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">// A gcTrigger is a predicate for starting a GC cycle. Specifically,</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">// it is an exit condition for the _GCoff phase.</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>type gcTrigger struct {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	kind gcTriggerKind
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	now  int64  <span class="comment">// gcTriggerTime: current time</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	n    uint32 <span class="comment">// gcTriggerCycle: cycle number to start</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>}
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>type gcTriggerKind int
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>const (
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	<span class="comment">// gcTriggerHeap indicates that a cycle should be started when</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	<span class="comment">// the heap size reaches the trigger heap size computed by the</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	<span class="comment">// controller.</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	gcTriggerHeap gcTriggerKind = iota
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	<span class="comment">// gcTriggerTime indicates that a cycle should be started when</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	<span class="comment">// it&#39;s been more than forcegcperiod nanoseconds since the</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	<span class="comment">// previous GC cycle.</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	gcTriggerTime
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	<span class="comment">// gcTriggerCycle indicates that a cycle should be started if</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	<span class="comment">// we have not yet started cycle number gcTrigger.n (relative</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	<span class="comment">// to work.cycles).</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	gcTriggerCycle
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span><span class="comment">// test reports whether the trigger condition is satisfied, meaning</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span><span class="comment">// that the exit condition for the _GCoff phase has been met. The exit</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span><span class="comment">// condition should be tested when allocating.</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>func (t gcTrigger) test() bool {
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	if !memstats.enablegc || panicking.Load() != 0 || gcphase != _GCoff {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		return false
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	switch t.kind {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	case gcTriggerHeap:
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		trigger, _ := gcController.trigger()
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		return gcController.heapLive.Load() &gt;= trigger
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	case gcTriggerTime:
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		if gcController.gcPercent.Load() &lt; 0 {
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>			return false
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		lastgc := int64(atomic.Load64(&amp;memstats.last_gc_nanotime))
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		return lastgc != 0 &amp;&amp; t.now-lastgc &gt; forcegcperiod
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	case gcTriggerCycle:
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		<span class="comment">// t.n &gt; work.cycles, but accounting for wraparound.</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		return int32(t.n-work.cycles.Load()) &gt; 0
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	return true
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span><span class="comment">// gcStart starts the GC. It transitions from _GCoff to _GCmark (if</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span><span class="comment">// debug.gcstoptheworld == 0) or performs all of GC (if</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span><span class="comment">// debug.gcstoptheworld != 0).</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span><span class="comment">// This may return without performing this transition in some cases,</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span><span class="comment">// such as when called on a system stack or with locks held.</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>func gcStart(trigger gcTrigger) {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	<span class="comment">// Since this is called from malloc and malloc is called in</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	<span class="comment">// the guts of a number of libraries that might be holding</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	<span class="comment">// locks, don&#39;t attempt to start GC in non-preemptible or</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	<span class="comment">// potentially unstable situations.</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	if gp := getg(); gp == mp.g0 || mp.locks &gt; 1 || mp.preemptoff != &#34;&#34; {
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		releasem(mp)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		return
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	}
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	releasem(mp)
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	mp = nil
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	<span class="comment">// Pick up the remaining unswept/not being swept spans concurrently</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	<span class="comment">// This shouldn&#39;t happen if we&#39;re being invoked in background</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	<span class="comment">// mode since proportional sweep should have just finished</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	<span class="comment">// sweeping everything, but rounding errors, etc, may leave a</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	<span class="comment">// few spans unswept. In forced mode, this is necessary since</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	<span class="comment">// GC can be forced at any point in the sweeping cycle.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	<span class="comment">// We check the transition condition continuously here in case</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	<span class="comment">// this G gets delayed in to the next GC cycle.</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	for trigger.test() &amp;&amp; sweepone() != ^uintptr(0) {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	}
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	<span class="comment">// Perform GC initialization and the sweep termination</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	<span class="comment">// transition.</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	semacquire(&amp;work.startSema)
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	<span class="comment">// Re-check transition condition under transition lock.</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	if !trigger.test() {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		semrelease(&amp;work.startSema)
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		return
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	}
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	<span class="comment">// In gcstoptheworld debug mode, upgrade the mode accordingly.</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	<span class="comment">// We do this after re-checking the transition condition so</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	<span class="comment">// that multiple goroutines that detect the heap trigger don&#39;t</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	<span class="comment">// start multiple STW GCs.</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	mode := gcBackgroundMode
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	if debug.gcstoptheworld == 1 {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		mode = gcForceMode
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	} else if debug.gcstoptheworld == 2 {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		mode = gcForceBlockMode
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	}
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	<span class="comment">// Ok, we&#39;re doing it! Stop everybody else</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	semacquire(&amp;gcsema)
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	semacquire(&amp;worldsema)
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	<span class="comment">// For stats, check if this GC was forced by the user.</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	<span class="comment">// Update it under gcsema to avoid gctrace getting wrong values.</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	work.userForced = trigger.kind == gcTriggerCycle
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	trace := traceAcquire()
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		trace.GCStart()
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	}
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	<span class="comment">// Check that all Ps have finished deferred mcache flushes.</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	for _, p := range allp {
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		if fg := p.mcache.flushGen.Load(); fg != mheap_.sweepgen {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>			println(&#34;runtime: p&#34;, p.id, &#34;flushGen&#34;, fg, &#34;!= sweepgen&#34;, mheap_.sweepgen)
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>			throw(&#34;p mcache not flushed&#34;)
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		}
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	}
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	gcBgMarkStartWorkers()
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	systemstack(gcResetMarkState)
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	work.stwprocs, work.maxprocs = gomaxprocs, gomaxprocs
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	if work.stwprocs &gt; ncpu {
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		<span class="comment">// This is used to compute CPU time of the STW phases,</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		<span class="comment">// so it can&#39;t be more than ncpu, even if GOMAXPROCS is.</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>		work.stwprocs = ncpu
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	}
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	work.heap0 = gcController.heapLive.Load()
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	work.pauseNS = 0
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	work.mode = mode
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	now := nanotime()
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	work.tSweepTerm = now
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	var stw worldStop
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		stw = stopTheWorldWithSema(stwGCSweepTerm)
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	})
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	<span class="comment">// Finish sweep before we start concurrent scan.</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		finishsweep_m()
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	})
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	<span class="comment">// clearpools before we start the GC. If we wait the memory will not be</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	<span class="comment">// reclaimed until the next GC cycle.</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	clearpools()
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	work.cycles.Add(1)
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	<span class="comment">// Assists and workers can start the moment we start</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	<span class="comment">// the world.</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	gcController.startCycle(now, int(gomaxprocs), trigger)
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	<span class="comment">// Notify the CPU limiter that assists may begin.</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	gcCPULimiter.startGCTransition(true, now)
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	<span class="comment">// In STW mode, disable scheduling of user Gs. This may also</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	<span class="comment">// disable scheduling of this goroutine, so it may block as</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	<span class="comment">// soon as we start the world again.</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	if mode != gcBackgroundMode {
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		schedEnableUser(false)
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	}
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	<span class="comment">// Enter concurrent mark phase and enable</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	<span class="comment">// write barriers.</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	<span class="comment">// Because the world is stopped, all Ps will</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	<span class="comment">// observe that write barriers are enabled by</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	<span class="comment">// the time we start the world and begin</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	<span class="comment">// scanning.</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	<span class="comment">// Write barriers must be enabled before assists are</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	<span class="comment">// enabled because they must be enabled before</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	<span class="comment">// any non-leaf heap objects are marked. Since</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	<span class="comment">// allocations are blocked until assists can</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	<span class="comment">// happen, we want to enable assists as early as</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	<span class="comment">// possible.</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	setGCPhase(_GCmark)
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	gcBgMarkPrepare() <span class="comment">// Must happen before assists are enabled.</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	gcMarkRootPrepare()
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>	<span class="comment">// Mark all active tinyalloc blocks. Since we&#39;re</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	<span class="comment">// allocating from these, they need to be black like</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	<span class="comment">// other allocations. The alternative is to blacken</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	<span class="comment">// the tiny block on every allocation from it, which</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	<span class="comment">// would slow down the tiny allocator.</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	gcMarkTinyAllocs()
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	<span class="comment">// At this point all Ps have enabled the write</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	<span class="comment">// barrier, thus maintaining the no white to</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	<span class="comment">// black invariant. Enable mutator assists to</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	<span class="comment">// put back-pressure on fast allocating</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	<span class="comment">// mutators.</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	atomic.Store(&amp;gcBlackenEnabled, 1)
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	<span class="comment">// In STW mode, we could block the instant systemstack</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	<span class="comment">// returns, so make sure we&#39;re not preemptible.</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	mp = acquirem()
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	<span class="comment">// Concurrent mark.</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		now = startTheWorldWithSema(0, stw)
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		work.pauseNS += now - stw.start
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		work.tMark = now
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		sweepTermCpu := int64(work.stwprocs) * (work.tMark - work.tSweepTerm)
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>		work.cpuStats.gcPauseTime += sweepTermCpu
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>		work.cpuStats.gcTotalTime += sweepTermCpu
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>		<span class="comment">// Release the CPU limiter.</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		gcCPULimiter.finishGCTransition(now)
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	})
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	<span class="comment">// Release the world sema before Gosched() in STW mode</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	<span class="comment">// because we will need to reacquire it later but before</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	<span class="comment">// this goroutine becomes runnable again, and we could</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	<span class="comment">// self-deadlock otherwise.</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	semrelease(&amp;worldsema)
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	releasem(mp)
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	<span class="comment">// Make sure we block instead of returning to user code</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	<span class="comment">// in STW mode.</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	if mode != gcBackgroundMode {
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		Gosched()
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	}
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	semrelease(&amp;work.startSema)
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>}
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span><span class="comment">// gcMarkDoneFlushed counts the number of P&#39;s with flushed work.</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span><span class="comment">// Ideally this would be a captured local in gcMarkDone, but forEachP</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span><span class="comment">// escapes its callback closure, so it can&#39;t capture anything.</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span><span class="comment">// This is protected by markDoneSema.</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>var gcMarkDoneFlushed uint32
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span><span class="comment">// gcMarkDone transitions the GC from mark to mark termination if all</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span><span class="comment">// reachable objects have been marked (that is, there are no grey</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span><span class="comment">// objects and can be no more in the future). Otherwise, it flushes</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span><span class="comment">// all local work to the global queues where it can be discovered by</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span><span class="comment">// other workers.</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span><span class="comment">// This should be called when all local mark work has been drained and</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span><span class="comment">// there are no remaining workers. Specifically, when</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span><span class="comment">//	work.nwait == work.nproc &amp;&amp; !gcMarkWorkAvailable(p)</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span><span class="comment">// The calling context must be preemptible.</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span><span class="comment">// Flushing local work is important because idle Ps may have local</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span><span class="comment">// work queued. This is the only way to make that work visible and</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span><span class="comment">// drive GC to completion.</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span><span class="comment">// It is explicitly okay to have write barriers in this function. If</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span><span class="comment">// it does transition to mark termination, then all reachable objects</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span><span class="comment">// have been marked, so the write barrier cannot shade any more</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span><span class="comment">// objects.</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>func gcMarkDone() {
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	<span class="comment">// Ensure only one thread is running the ragged barrier at a</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	<span class="comment">// time.</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	semacquire(&amp;work.markDoneSema)
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>top:
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	<span class="comment">// Re-check transition condition under transition lock.</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s critical that this checks the global work queues are</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	<span class="comment">// empty before performing the ragged barrier. Otherwise,</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	<span class="comment">// there could be global work that a P could take after the P</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	<span class="comment">// has passed the ragged barrier.</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	if !(gcphase == _GCmark &amp;&amp; work.nwait == work.nproc &amp;&amp; !gcMarkWorkAvailable(nil)) {
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>		semrelease(&amp;work.markDoneSema)
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		return
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	}
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	<span class="comment">// forEachP needs worldsema to execute, and we&#39;ll need it to</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	<span class="comment">// stop the world later, so acquire worldsema now.</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	semacquire(&amp;worldsema)
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	<span class="comment">// Flush all local buffers and collect flushedWork flags.</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	gcMarkDoneFlushed = 0
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	forEachP(waitReasonGCMarkTermination, func(pp *p) {
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>		<span class="comment">// Flush the write barrier buffer, since this may add</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		<span class="comment">// work to the gcWork.</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		wbBufFlush1(pp)
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		<span class="comment">// Flush the gcWork, since this may create global work</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>		<span class="comment">// and set the flushedWork flag.</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		<span class="comment">// TODO(austin): Break up these workbufs to</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		<span class="comment">// better distribute work.</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		pp.gcw.dispose()
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		<span class="comment">// Collect the flushedWork flag.</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		if pp.gcw.flushedWork {
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>			atomic.Xadd(&amp;gcMarkDoneFlushed, 1)
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>			pp.gcw.flushedWork = false
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	})
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	if gcMarkDoneFlushed != 0 {
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		<span class="comment">// More grey objects were discovered since the</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		<span class="comment">// previous termination check, so there may be more</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>		<span class="comment">// work to do. Keep going. It&#39;s possible the</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		<span class="comment">// transition condition became true again during the</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		<span class="comment">// ragged barrier, so re-check it.</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>		semrelease(&amp;worldsema)
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		goto top
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	}
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	<span class="comment">// There was no global work, no local work, and no Ps</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	<span class="comment">// communicated work since we took markDoneSema. Therefore</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	<span class="comment">// there are no grey objects and no more objects can be</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	<span class="comment">// shaded. Transition to mark termination.</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	now := nanotime()
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	work.tMarkTerm = now
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	getg().m.preemptoff = &#34;gcing&#34;
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	var stw worldStop
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		stw = stopTheWorldWithSema(stwGCMarkTerm)
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>	})
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	<span class="comment">// The gcphase is _GCmark, it will transition to _GCmarktermination</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	<span class="comment">// below. The important thing is that the wb remains active until</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	<span class="comment">// all marking is complete. This includes writes made by the GC.</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	<span class="comment">// There is sometimes work left over when we enter mark termination due</span>
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	<span class="comment">// to write barriers performed after the completion barrier above.</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	<span class="comment">// Detect this and resume concurrent mark. This is obviously</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	<span class="comment">// unfortunate.</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	<span class="comment">// See issue #27993 for details.</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	<span class="comment">// Switch to the system stack to call wbBufFlush1, though in this case</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	<span class="comment">// it doesn&#39;t matter because we&#39;re non-preemptible anyway.</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	restart := false
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>		for _, p := range allp {
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>			wbBufFlush1(p)
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>			if !p.gcw.empty() {
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>				restart = true
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>				break
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>			}
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>		}
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>	})
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>	if restart {
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>		getg().m.preemptoff = &#34;&#34;
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>			now := startTheWorldWithSema(0, stw)
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>			work.pauseNS += now - stw.start
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>		})
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		semrelease(&amp;worldsema)
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>		goto top
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	}
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	gcComputeStartingStackSize()
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	<span class="comment">// Disable assists and background workers. We must do</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	<span class="comment">// this before waking blocked assists.</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>	atomic.Store(&amp;gcBlackenEnabled, 0)
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	<span class="comment">// Notify the CPU limiter that GC assists will now cease.</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	gcCPULimiter.startGCTransition(false, now)
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	<span class="comment">// Wake all blocked assists. These will run when we</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	<span class="comment">// start the world again.</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	gcWakeAllAssists()
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	<span class="comment">// Likewise, release the transition lock. Blocked</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	<span class="comment">// workers and assists will run when we start the</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	<span class="comment">// world again.</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	semrelease(&amp;work.markDoneSema)
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>	<span class="comment">// In STW mode, re-enable user goroutines. These will be</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>	<span class="comment">// queued to run after we start the world.</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>	schedEnableUser(true)
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>	<span class="comment">// endCycle depends on all gcWork cache stats being flushed.</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	<span class="comment">// The termination algorithm above ensured that up to</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	<span class="comment">// allocations since the ragged barrier.</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	gcController.endCycle(now, int(gomaxprocs), work.userForced)
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	<span class="comment">// Perform mark termination. This will restart the world.</span>
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	gcMarkTermination(stw)
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>}
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span><span class="comment">// World must be stopped and mark assists and background workers must be</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span><span class="comment">// disabled.</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>func gcMarkTermination(stw worldStop) {
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>	<span class="comment">// Start marktermination (write barrier remains enabled for now).</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	setGCPhase(_GCmarktermination)
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>	work.heap1 = gcController.heapLive.Load()
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>	startTime := nanotime()
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	mp.preemptoff = &#34;gcing&#34;
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>	mp.traceback = 2
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>	curgp := mp.curg
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>	<span class="comment">// N.B. The execution tracer is not aware of this status</span>
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>	<span class="comment">// transition and handles it specially based on the</span>
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>	<span class="comment">// wait reason.</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>	casGToWaiting(curgp, _Grunning, waitReasonGarbageCollection)
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>	<span class="comment">// Run gc on the g0 stack. We do this so that the g stack</span>
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>	<span class="comment">// we&#39;re currently running on will no longer change. Cuts</span>
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>	<span class="comment">// the root set down a bit (g0 stacks are not scanned, and</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>	<span class="comment">// we don&#39;t need to scan gc&#39;s internal state).  We also</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>	<span class="comment">// need to switch to g0 so we can shrink the stack.</span>
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>		gcMark(startTime)
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>		<span class="comment">// Must return immediately.</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>		<span class="comment">// The outer function&#39;s stack may have moved</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>		<span class="comment">// during gcMark (it shrinks stacks, including the</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>		<span class="comment">// outer function&#39;s stack), so we must not refer</span>
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		<span class="comment">// to any of its variables. Return back to the</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>		<span class="comment">// non-system stack to pick up the new addresses</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>		<span class="comment">// before continuing.</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>	})
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	var stwSwept bool
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		work.heap2 = work.bytesMarked
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>		if debug.gccheckmark &gt; 0 {
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>			<span class="comment">// Run a full non-parallel, stop-the-world</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>			<span class="comment">// mark using checkmark bits, to check that we</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>			<span class="comment">// didn&#39;t forget to mark anything during the</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>			<span class="comment">// concurrent mark process.</span>
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>			startCheckmarks()
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>			gcResetMarkState()
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>			gcw := &amp;getg().m.p.ptr().gcw
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>			gcDrain(gcw, 0)
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>			wbBufFlush1(getg().m.p.ptr())
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>			gcw.dispose()
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>			endCheckmarks()
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>		}
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>		<span class="comment">// marking is complete so we can turn the write barrier off</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		setGCPhase(_GCoff)
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		stwSwept = gcSweep(work.mode)
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>	})
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>	mp.traceback = 0
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	casgstatus(curgp, _Gwaiting, _Grunning)
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	trace := traceAcquire()
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>		trace.GCDone()
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	}
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>	<span class="comment">// all done</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	mp.preemptoff = &#34;&#34;
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	if gcphase != _GCoff {
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>		throw(&#34;gc done but gcphase != _GCoff&#34;)
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	}
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>	<span class="comment">// Record heapInUse for scavenger.</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>	memstats.lastHeapInUse = gcController.heapInUse.load()
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	<span class="comment">// Update GC trigger and pacing, as well as downstream consumers</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	<span class="comment">// of this pacing information, for the next cycle.</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	systemstack(gcControllerCommit)
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>	<span class="comment">// Update timing memstats</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>	now := nanotime()
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>	sec, nsec, _ := time_now()
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	unixNow := sec*1e9 + int64(nsec)
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	work.pauseNS += now - stw.start
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	work.tEnd = now
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	atomic.Store64(&amp;memstats.last_gc_unix, uint64(unixNow)) <span class="comment">// must be Unix time to make sense to user</span>
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	atomic.Store64(&amp;memstats.last_gc_nanotime, uint64(now)) <span class="comment">// monotonic time for us</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	memstats.pause_ns[memstats.numgc%uint32(len(memstats.pause_ns))] = uint64(work.pauseNS)
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	memstats.pause_end[memstats.numgc%uint32(len(memstats.pause_end))] = uint64(unixNow)
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>	memstats.pause_total_ns += uint64(work.pauseNS)
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	markTermCpu := int64(work.stwprocs) * (work.tEnd - work.tMarkTerm)
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	work.cpuStats.gcPauseTime += markTermCpu
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	work.cpuStats.gcTotalTime += markTermCpu
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>	<span class="comment">// Accumulate CPU stats.</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>	<span class="comment">// Pass gcMarkPhase=true so we can get all the latest GC CPU stats in there too.</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	work.cpuStats.accumulate(now, true)
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>	<span class="comment">// Compute overall GC CPU utilization.</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	<span class="comment">// Omit idle marking time from the overall utilization here since it&#39;s &#34;free&#34;.</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	memstats.gc_cpu_fraction = float64(work.cpuStats.gcTotalTime-work.cpuStats.gcIdleTime) / float64(work.cpuStats.totalTime)
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>	<span class="comment">// Reset assist time and background time stats.</span>
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>	<span class="comment">// Do this now, instead of at the start of the next GC cycle, because</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>	<span class="comment">// these two may keep accumulating even if the GC is not active.</span>
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	scavenge.assistTime.Store(0)
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>	scavenge.backgroundTime.Store(0)
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>	<span class="comment">// Reset idle time stat.</span>
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>	sched.idleTime.Store(0)
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>	if work.userForced {
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>		memstats.numforcedgc++
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>	}
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>	<span class="comment">// Bump GC cycle count and wake goroutines waiting on sweep.</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>	lock(&amp;work.sweepWaiters.lock)
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>	memstats.numgc++
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>	injectglist(&amp;work.sweepWaiters.list)
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>	unlock(&amp;work.sweepWaiters.lock)
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>	<span class="comment">// Increment the scavenge generation now.</span>
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>	<span class="comment">// This moment represents peak heap in use because we&#39;re</span>
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>	<span class="comment">// about to start sweeping.</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>	mheap_.pages.scav.index.nextGen()
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>	<span class="comment">// Release the CPU limiter.</span>
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>	gcCPULimiter.finishGCTransition(now)
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>	<span class="comment">// Finish the current heap profiling cycle and start a new</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>	<span class="comment">// heap profiling cycle. We do this before starting the world</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>	<span class="comment">// so events don&#39;t leak into the wrong cycle.</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	mProf_NextCycle()
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>	<span class="comment">// There may be stale spans in mcaches that need to be swept.</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>	<span class="comment">// Those aren&#39;t tracked in any sweep lists, so we need to</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	<span class="comment">// count them against sweep completion until we ensure all</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	<span class="comment">// those spans have been forced out.</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	<span class="comment">// If gcSweep fully swept the heap (for example if the sweep</span>
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	<span class="comment">// is not concurrent due to a GODEBUG setting), then we expect</span>
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	<span class="comment">// the sweepLocker to be invalid, since sweeping is done.</span>
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	<span class="comment">// N.B. Below we might duplicate some work from gcSweep; this is</span>
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>	<span class="comment">// fine as all that work is idempotent within a GC cycle, and</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	<span class="comment">// we&#39;re still holding worldsema so a new cycle can&#39;t start.</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>	sl := sweep.active.begin()
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>	if !stwSwept &amp;&amp; !sl.valid {
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>		throw(&#34;failed to set sweep barrier&#34;)
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>	} else if stwSwept &amp;&amp; sl.valid {
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>		throw(&#34;non-concurrent sweep failed to drain all sweep queues&#34;)
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>	}
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>		<span class="comment">// The memstats updated above must be updated with the world</span>
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>		<span class="comment">// stopped to ensure consistency of some values, such as</span>
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>		<span class="comment">// sched.idleTime and sched.totaltime. memstats also include</span>
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>		<span class="comment">// the pause time (work,pauseNS), forcing computation of the</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>		<span class="comment">// total pause time before the pause actually ends.</span>
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>		<span class="comment">// Here we reuse the same now for start the world so that the</span>
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>		<span class="comment">// time added to /sched/pauses/total/gc:seconds will be</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>		<span class="comment">// consistent with the value in memstats.</span>
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>		startTheWorldWithSema(now, stw)
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>	})
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>	<span class="comment">// Flush the heap profile so we can start a new cycle next GC.</span>
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	<span class="comment">// This is relatively expensive, so we don&#39;t do it with the</span>
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>	<span class="comment">// world stopped.</span>
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>	mProf_Flush()
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>	<span class="comment">// Prepare workbufs for freeing by the sweeper. We do this</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>	<span class="comment">// asynchronously because it can take non-trivial time.</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>	prepareFreeWorkbufs()
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>	<span class="comment">// Free stack spans. This must be done between GC cycles.</span>
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>	systemstack(freeStackSpans)
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>	<span class="comment">// Ensure all mcaches are flushed. Each P will flush its own</span>
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>	<span class="comment">// mcache before allocating, but idle Ps may not. Since this</span>
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>	<span class="comment">// is necessary to sweep all spans, we need to ensure all</span>
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>	<span class="comment">// mcaches are flushed before we start the next GC cycle.</span>
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>	<span class="comment">// While we&#39;re here, flush the page cache for idle Ps to avoid</span>
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>	<span class="comment">// having pages get stuck on them. These pages are hidden from</span>
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>	<span class="comment">// the scavenger, so in small idle heaps a significant amount</span>
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>	<span class="comment">// of additional memory might be held onto.</span>
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	<span class="comment">// Also, flush the pinner cache, to avoid leaking that memory</span>
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	<span class="comment">// indefinitely.</span>
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>	forEachP(waitReasonFlushProcCaches, func(pp *p) {
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>		pp.mcache.prepareForSweep()
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>		if pp.status == _Pidle {
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>			systemstack(func() {
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>				lock(&amp;mheap_.lock)
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>				pp.pcache.flush(&amp;mheap_.pages)
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>				unlock(&amp;mheap_.lock)
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>			})
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>		}
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>		pp.pinnerCache = nil
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>	})
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>	if sl.valid {
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>		<span class="comment">// Now that we&#39;ve swept stale spans in mcaches, they don&#39;t</span>
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>		<span class="comment">// count against unswept spans.</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>		<span class="comment">// Note: this sweepLocker may not be valid if sweeping had</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>		<span class="comment">// already completed during the STW. See the corresponding</span>
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>		<span class="comment">// begin() call that produced sl.</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>		sweep.active.end(sl)
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	}
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>	<span class="comment">// Print gctrace before dropping worldsema. As soon as we drop</span>
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>	<span class="comment">// worldsema another cycle could start and smash the stats</span>
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>	<span class="comment">// we&#39;re trying to print.</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>	if debug.gctrace &gt; 0 {
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>		util := int(memstats.gc_cpu_fraction * 100)
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>		var sbuf [24]byte
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>		printlock()
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>		print(&#34;gc &#34;, memstats.numgc,
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>			&#34; @&#34;, string(itoaDiv(sbuf[:], uint64(work.tSweepTerm-runtimeInitTime)/1e6, 3)), &#34;s &#34;,
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>			util, &#34;%: &#34;)
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>		prev := work.tSweepTerm
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>		for i, ns := range []int64{work.tMark, work.tMarkTerm, work.tEnd} {
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>			if i != 0 {
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>				print(&#34;+&#34;)
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>			}
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>			print(string(fmtNSAsMS(sbuf[:], uint64(ns-prev))))
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>			prev = ns
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>		}
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>		print(&#34; ms clock, &#34;)
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>		for i, ns := range []int64{
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>			int64(work.stwprocs) * (work.tMark - work.tSweepTerm),
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>			gcController.assistTime.Load(),
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>			gcController.dedicatedMarkTime.Load() + gcController.fractionalMarkTime.Load(),
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>			gcController.idleMarkTime.Load(),
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>			markTermCpu,
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>		} {
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>			if i == 2 || i == 3 {
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>				<span class="comment">// Separate mark time components with /.</span>
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>				print(&#34;/&#34;)
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>			} else if i != 0 {
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>				print(&#34;+&#34;)
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>			}
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>			print(string(fmtNSAsMS(sbuf[:], uint64(ns))))
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>		}
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>		print(&#34; ms cpu, &#34;,
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>			work.heap0&gt;&gt;20, &#34;-&gt;&#34;, work.heap1&gt;&gt;20, &#34;-&gt;&#34;, work.heap2&gt;&gt;20, &#34; MB, &#34;,
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>			gcController.lastHeapGoal&gt;&gt;20, &#34; MB goal, &#34;,
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>			gcController.lastStackScan.Load()&gt;&gt;20, &#34; MB stacks, &#34;,
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>			gcController.globalsScan.Load()&gt;&gt;20, &#34; MB globals, &#34;,
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>			work.maxprocs, &#34; P&#34;)
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>		if work.userForced {
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>			print(&#34; (forced)&#34;)
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>		}
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>		printunlock()
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>	}
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>	<span class="comment">// Set any arena chunks that were deferred to fault.</span>
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>	lock(&amp;userArenaState.lock)
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>	faultList := userArenaState.fault
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>	userArenaState.fault = nil
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>	unlock(&amp;userArenaState.lock)
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>	for _, lc := range faultList {
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>		lc.mspan.setUserArenaChunkToFault()
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>	}
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>	<span class="comment">// Enable huge pages on some metadata if we cross a heap threshold.</span>
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>	if gcController.heapGoal() &gt; minHeapForMetadataHugePages {
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>			mheap_.enableMetadataHugePages()
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>		})
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>	}
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>	semrelease(&amp;worldsema)
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>	semrelease(&amp;gcsema)
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>	<span class="comment">// Careful: another GC cycle may start now.</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>	releasem(mp)
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>	mp = nil
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>	<span class="comment">// now that gc is done, kick off finalizer thread if needed</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>	if !concurrentSweep {
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>		<span class="comment">// give the queued finalizers, if any, a chance to run</span>
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>		Gosched()
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>	}
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>}
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span><span class="comment">// gcBgMarkStartWorkers prepares background mark worker goroutines. These</span>
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span><span class="comment">// goroutines will not run until the mark phase, but they must be started while</span>
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span><span class="comment">// the work is not stopped and from a regular G stack. The caller must hold</span>
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span><span class="comment">// worldsema.</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>func gcBgMarkStartWorkers() {
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>	<span class="comment">// Background marking is performed by per-P G&#39;s. Ensure that each P has</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>	<span class="comment">// a background GC G.</span>
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	<span class="comment">// Worker Gs don&#39;t exit if gomaxprocs is reduced. If it is raised</span>
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>	<span class="comment">// again, we can reuse the old workers; no need to create new workers.</span>
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>	for gcBgMarkWorkerCount &lt; gomaxprocs {
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>		go gcBgMarkWorker()
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>		notetsleepg(&amp;work.bgMarkReady, -1)
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>		noteclear(&amp;work.bgMarkReady)
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>		<span class="comment">// The worker is now guaranteed to be added to the pool before</span>
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>		<span class="comment">// its P&#39;s next findRunnableGCWorker.</span>
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>		gcBgMarkWorkerCount++
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>	}
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>}
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span><span class="comment">// gcBgMarkPrepare sets up state for background marking.</span>
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span><span class="comment">// Mutator assists must not yet be enabled.</span>
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>func gcBgMarkPrepare() {
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>	<span class="comment">// Background marking will stop when the work queues are empty</span>
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>	<span class="comment">// and there are no more workers (note that, since this is</span>
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>	<span class="comment">// concurrent, this may be a transient state, but mark</span>
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>	<span class="comment">// termination will clean it up). Between background workers</span>
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>	<span class="comment">// and assists, we don&#39;t really know how many workers there</span>
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>	<span class="comment">// will be, so we pretend to have an arbitrarily large number</span>
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>	<span class="comment">// of workers, almost all of which are &#34;waiting&#34;. While a</span>
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>	<span class="comment">// worker is working it decrements nwait. If nproc == nwait,</span>
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>	<span class="comment">// there are no workers.</span>
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>	work.nproc = ^uint32(0)
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>	work.nwait = ^uint32(0)
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>}
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span><span class="comment">// gcBgMarkWorkerNode is an entry in the gcBgMarkWorkerPool. It points to a single</span>
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span><span class="comment">// gcBgMarkWorker goroutine.</span>
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>type gcBgMarkWorkerNode struct {
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>	<span class="comment">// Unused workers are managed in a lock-free stack. This field must be first.</span>
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>	node lfnode
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>	<span class="comment">// The g of this worker.</span>
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>	gp guintptr
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	<span class="comment">// Release this m on park. This is used to communicate with the unlock</span>
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	<span class="comment">// function, which cannot access the G&#39;s stack. It is unused outside of</span>
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>	<span class="comment">// gcBgMarkWorker().</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	m muintptr
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>}
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>func gcBgMarkWorker() {
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>	gp := getg()
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>	<span class="comment">// We pass node to a gopark unlock function, so it can&#39;t be on</span>
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>	<span class="comment">// the stack (see gopark). Prevent deadlock from recursively</span>
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>	<span class="comment">// starting GC by disabling preemption.</span>
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>	gp.m.preemptoff = &#34;GC worker init&#34;
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>	node := new(gcBgMarkWorkerNode)
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>	gp.m.preemptoff = &#34;&#34;
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>	node.gp.set(gp)
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>	node.m.set(acquirem())
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>	notewakeup(&amp;work.bgMarkReady)
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>	<span class="comment">// After this point, the background mark worker is generally scheduled</span>
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>	<span class="comment">// cooperatively by gcController.findRunnableGCWorker. While performing</span>
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>	<span class="comment">// work on the P, preemption is disabled because we are working on</span>
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>	<span class="comment">// P-local work buffers. When the preempt flag is set, this puts itself</span>
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>	<span class="comment">// into _Gwaiting to be woken up by gcController.findRunnableGCWorker</span>
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>	<span class="comment">// at the appropriate time.</span>
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>	<span class="comment">// When preemption is enabled (e.g., while in gcMarkDone), this worker</span>
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>	<span class="comment">// may be preempted and schedule as a _Grunnable G from a runq. That is</span>
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>	<span class="comment">// fine; it will eventually gopark again for further scheduling via</span>
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>	<span class="comment">// findRunnableGCWorker.</span>
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>	<span class="comment">// Since we disable preemption before notifying bgMarkReady, we</span>
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>	<span class="comment">// guarantee that this G will be in the worker pool for the next</span>
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>	<span class="comment">// findRunnableGCWorker. This isn&#39;t strictly necessary, but it reduces</span>
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>	<span class="comment">// latency between _GCmark starting and the workers starting.</span>
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>	for {
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>		<span class="comment">// Go to sleep until woken by</span>
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>		<span class="comment">// gcController.findRunnableGCWorker.</span>
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>		gopark(func(g *g, nodep unsafe.Pointer) bool {
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>			node := (*gcBgMarkWorkerNode)(nodep)
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>			if mp := node.m.ptr(); mp != nil {
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>				<span class="comment">// The worker G is no longer running; release</span>
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>				<span class="comment">// the M.</span>
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>				<span class="comment">// N.B. it is _safe_ to release the M as soon</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>				<span class="comment">// as we are no longer performing P-local mark</span>
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>				<span class="comment">// work.</span>
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>				<span class="comment">// However, since we cooperatively stop work</span>
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>				<span class="comment">// when gp.preempt is set, if we releasem in</span>
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>				<span class="comment">// the loop then the following call to gopark</span>
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>				<span class="comment">// would immediately preempt the G. This is</span>
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>				<span class="comment">// also safe, but inefficient: the G must</span>
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>				<span class="comment">// schedule again only to enter gopark and park</span>
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>				<span class="comment">// again. Thus, we defer the release until</span>
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>				<span class="comment">// after parking the G.</span>
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>				releasem(mp)
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>			}
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>			<span class="comment">// Release this G to the pool.</span>
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>			gcBgMarkWorkerPool.push(&amp;node.node)
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>			<span class="comment">// Note that at this point, the G may immediately be</span>
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>			<span class="comment">// rescheduled and may be running.</span>
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>			return true
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>		}, unsafe.Pointer(node), waitReasonGCWorkerIdle, traceBlockSystemGoroutine, 0)
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>		<span class="comment">// Preemption must not occur here, or another G might see</span>
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>		<span class="comment">// p.gcMarkWorkerMode.</span>
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>		<span class="comment">// Disable preemption so we can use the gcw. If the</span>
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>		<span class="comment">// scheduler wants to preempt us, we&#39;ll stop draining,</span>
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>		<span class="comment">// dispose the gcw, and then preempt.</span>
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>		node.m.set(acquirem())
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>		pp := gp.m.p.ptr() <span class="comment">// P can&#39;t change with preemption disabled.</span>
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>		if gcBlackenEnabled == 0 {
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>			println(&#34;worker mode&#34;, pp.gcMarkWorkerMode)
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>			throw(&#34;gcBgMarkWorker: blackening not enabled&#34;)
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>		}
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>		if pp.gcMarkWorkerMode == gcMarkWorkerNotWorker {
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>			throw(&#34;gcBgMarkWorker: mode not set&#34;)
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>		}
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>		startTime := nanotime()
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>		pp.gcMarkWorkerStartTime = startTime
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>		var trackLimiterEvent bool
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>		if pp.gcMarkWorkerMode == gcMarkWorkerIdleMode {
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>			trackLimiterEvent = pp.limiterEvent.start(limiterEventIdleMarkWork, startTime)
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>		}
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>		decnwait := atomic.Xadd(&amp;work.nwait, -1)
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>		if decnwait == work.nproc {
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>			println(&#34;runtime: work.nwait=&#34;, decnwait, &#34;work.nproc=&#34;, work.nproc)
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>			throw(&#34;work.nwait was &gt; work.nproc&#34;)
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>		}
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>			<span class="comment">// Mark our goroutine preemptible so its stack</span>
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>			<span class="comment">// can be scanned. This lets two mark workers</span>
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>			<span class="comment">// scan each other (otherwise, they would</span>
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>			<span class="comment">// deadlock). We must not modify anything on</span>
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>			<span class="comment">// the G stack. However, stack shrinking is</span>
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>			<span class="comment">// disabled for mark workers, so it is safe to</span>
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>			<span class="comment">// read from the G stack.</span>
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>			<span class="comment">// N.B. The execution tracer is not aware of this status</span>
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>			<span class="comment">// transition and handles it specially based on the</span>
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>			<span class="comment">// wait reason.</span>
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>			casGToWaiting(gp, _Grunning, waitReasonGCWorkerActive)
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>			switch pp.gcMarkWorkerMode {
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>			default:
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>				throw(&#34;gcBgMarkWorker: unexpected gcMarkWorkerMode&#34;)
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>			case gcMarkWorkerDedicatedMode:
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>				gcDrainMarkWorkerDedicated(&amp;pp.gcw, true)
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>				if gp.preempt {
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>					<span class="comment">// We were preempted. This is</span>
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>					<span class="comment">// a useful signal to kick</span>
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>					<span class="comment">// everything out of the run</span>
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>					<span class="comment">// queue so it can run</span>
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>					<span class="comment">// somewhere else.</span>
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>					if drainQ, n := runqdrain(pp); n &gt; 0 {
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>						lock(&amp;sched.lock)
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>						globrunqputbatch(&amp;drainQ, int32(n))
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>						unlock(&amp;sched.lock)
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>					}
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>				}
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>				<span class="comment">// Go back to draining, this time</span>
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>				<span class="comment">// without preemption.</span>
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>				gcDrainMarkWorkerDedicated(&amp;pp.gcw, false)
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>			case gcMarkWorkerFractionalMode:
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>				gcDrainMarkWorkerFractional(&amp;pp.gcw)
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>			case gcMarkWorkerIdleMode:
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>				gcDrainMarkWorkerIdle(&amp;pp.gcw)
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>			}
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>			casgstatus(gp, _Gwaiting, _Grunning)
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>		})
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>		<span class="comment">// Account for time and mark us as stopped.</span>
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>		now := nanotime()
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>		duration := now - startTime
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>		gcController.markWorkerStop(pp.gcMarkWorkerMode, duration)
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>		if trackLimiterEvent {
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>			pp.limiterEvent.stop(limiterEventIdleMarkWork, now)
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>		}
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>		if pp.gcMarkWorkerMode == gcMarkWorkerFractionalMode {
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>			atomic.Xaddint64(&amp;pp.gcFractionalMarkTime, duration)
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>		}
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>		<span class="comment">// Was this the last worker and did we run out</span>
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>		<span class="comment">// of work?</span>
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>		incnwait := atomic.Xadd(&amp;work.nwait, +1)
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>		if incnwait &gt; work.nproc {
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>			println(&#34;runtime: p.gcMarkWorkerMode=&#34;, pp.gcMarkWorkerMode,
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>				&#34;work.nwait=&#34;, incnwait, &#34;work.nproc=&#34;, work.nproc)
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>			throw(&#34;work.nwait &gt; work.nproc&#34;)
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>		}
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span>
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span>		<span class="comment">// We&#39;ll releasem after this point and thus this P may run</span>
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>		<span class="comment">// something else. We must clear the worker mode to avoid</span>
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>		<span class="comment">// attributing the mode to a different (non-worker) G in</span>
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>		<span class="comment">// traceGoStart.</span>
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span>		pp.gcMarkWorkerMode = gcMarkWorkerNotWorker
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span>
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span>		<span class="comment">// If this worker reached a background mark completion</span>
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span>		<span class="comment">// point, signal the main GC goroutine.</span>
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span>		if incnwait == work.nproc &amp;&amp; !gcMarkWorkAvailable(nil) {
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span>			<span class="comment">// We don&#39;t need the P-local buffers here, allow</span>
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>			<span class="comment">// preemption because we may schedule like a regular</span>
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>			<span class="comment">// goroutine in gcMarkDone (block on locks, etc).</span>
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>			releasem(node.m.ptr())
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>			node.m.set(nil)
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>			gcMarkDone()
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>		}
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>	}
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>}
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span>
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span><span class="comment">// gcMarkWorkAvailable reports whether executing a mark worker</span>
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span><span class="comment">// on p is potentially useful. p may be nil, in which case it only</span>
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span><span class="comment">// checks the global sources of work.</span>
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span>func gcMarkWorkAvailable(p *p) bool {
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span>	if p != nil &amp;&amp; !p.gcw.empty() {
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span>		return true
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span>	}
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span>	if !work.full.empty() {
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span>		return true <span class="comment">// global work available</span>
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span>	}
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span>	if work.markrootNext &lt; work.markrootJobs {
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span>		return true <span class="comment">// root scan work available</span>
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span>	}
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span>	return false
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span>}
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span>
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span><span class="comment">// gcMark runs the mark (or, for concurrent GC, mark termination)</span>
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span><span class="comment">// All gcWork caches must be empty.</span>
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span><span class="comment">// STW is in effect at this point.</span>
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span>func gcMark(startTime int64) {
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span>	if debug.allocfreetrace &gt; 0 {
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span>		tracegc()
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span>	}
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span>
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span>	if gcphase != _GCmarktermination {
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span>		throw(&#34;in gcMark expecting to see gcphase as _GCmarktermination&#34;)
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span>	}
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span>	work.tstart = startTime
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span>
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span>	<span class="comment">// Check that there&#39;s no marking work remaining.</span>
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span>	if work.full != 0 || work.markrootNext &lt; work.markrootJobs {
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span>		print(&#34;runtime: full=&#34;, hex(work.full), &#34; next=&#34;, work.markrootNext, &#34; jobs=&#34;, work.markrootJobs, &#34; nDataRoots=&#34;, work.nDataRoots, &#34; nBSSRoots=&#34;, work.nBSSRoots, &#34; nSpanRoots=&#34;, work.nSpanRoots, &#34; nStackRoots=&#34;, work.nStackRoots, &#34;\n&#34;)
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span>		panic(&#34;non-empty mark queue after concurrent mark&#34;)
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span>	}
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span>
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span>	if debug.gccheckmark &gt; 0 {
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span>		<span class="comment">// This is expensive when there&#39;s a large number of</span>
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span>		<span class="comment">// Gs, so only do it if checkmark is also enabled.</span>
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span>		gcMarkRootCheck()
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span>	}
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span>
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span>	<span class="comment">// Drop allg snapshot. allgs may have grown, in which case</span>
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span>	<span class="comment">// this is the only reference to the old backing store and</span>
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span>	<span class="comment">// there&#39;s no need to keep it around.</span>
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span>	work.stackRoots = nil
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span>
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span>	<span class="comment">// Clear out buffers and double-check that all gcWork caches</span>
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span>	<span class="comment">// are empty. This should be ensured by gcMarkDone before we</span>
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span>	<span class="comment">// enter mark termination.</span>
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span>	<span class="comment">// TODO: We could clear out buffers just before mark if this</span>
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span>	<span class="comment">// has a non-negligible impact on STW time.</span>
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span>	for _, p := range allp {
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span>		<span class="comment">// The write barrier may have buffered pointers since</span>
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span>		<span class="comment">// the gcMarkDone barrier. However, since the barrier</span>
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span>		<span class="comment">// ensured all reachable objects were marked, all of</span>
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span>		<span class="comment">// these must be pointers to black objects. Hence we</span>
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span>		<span class="comment">// can just discard the write barrier buffer.</span>
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span>		if debug.gccheckmark &gt; 0 {
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span>			<span class="comment">// For debugging, flush the buffer and make</span>
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span>			<span class="comment">// sure it really was all marked.</span>
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span>			wbBufFlush1(p)
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span>		} else {
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span>			p.wbBuf.reset()
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>		}
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span>		gcw := &amp;p.gcw
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span>		if !gcw.empty() {
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span>			printlock()
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span>			print(&#34;runtime: P &#34;, p.id, &#34; flushedWork &#34;, gcw.flushedWork)
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span>			if gcw.wbuf1 == nil {
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span>				print(&#34; wbuf1=&lt;nil&gt;&#34;)
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span>			} else {
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span>				print(&#34; wbuf1.n=&#34;, gcw.wbuf1.nobj)
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span>			}
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span>			if gcw.wbuf2 == nil {
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span>				print(&#34; wbuf2=&lt;nil&gt;&#34;)
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span>			} else {
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span>				print(&#34; wbuf2.n=&#34;, gcw.wbuf2.nobj)
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>			}
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>			print(&#34;\n&#34;)
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>			throw(&#34;P has cached GC work at end of mark termination&#34;)
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>		}
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span>		<span class="comment">// There may still be cached empty buffers, which we</span>
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span>		<span class="comment">// need to flush since we&#39;re going to free them. Also,</span>
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span>		<span class="comment">// there may be non-zero stats because we allocated</span>
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span>		<span class="comment">// black after the gcMarkDone barrier.</span>
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>		gcw.dispose()
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>	}
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span>	<span class="comment">// Flush scanAlloc from each mcache since we&#39;re about to modify</span>
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span>	<span class="comment">// heapScan directly. If we were to flush this later, then scanAlloc</span>
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span>	<span class="comment">// might have incorrect information.</span>
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>	<span class="comment">// Note that it&#39;s not important to retain this information; we know</span>
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>	<span class="comment">// exactly what heapScan is at this point via scanWork.</span>
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>	for _, p := range allp {
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>		c := p.mcache
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span>		if c == nil {
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span>			continue
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span>		}
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span>		c.scanAlloc = 0
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span>	}
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span>
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span>	<span class="comment">// Reset controller state.</span>
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span>	gcController.resetLive(work.bytesMarked)
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span>}
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span>
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span><span class="comment">// gcSweep must be called on the system stack because it acquires the heap</span>
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span><span class="comment">// lock. See mheap for details.</span>
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span><span class="comment">// Returns true if the heap was fully swept by this function.</span>
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span><span class="comment">// The world must be stopped.</span>
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span>func gcSweep(mode gcMode) bool {
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span>
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span>	if gcphase != _GCoff {
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span>		throw(&#34;gcSweep being done but phase is not GCoff&#34;)
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span>	}
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span>
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span>	lock(&amp;mheap_.lock)
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span>	mheap_.sweepgen += 2
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span>	sweep.active.reset()
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span>	mheap_.pagesSwept.Store(0)
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span>	mheap_.sweepArenas = mheap_.allArenas
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span>	mheap_.reclaimIndex.Store(0)
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span>	mheap_.reclaimCredit.Store(0)
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span>	unlock(&amp;mheap_.lock)
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span>
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>	sweep.centralIndex.clear()
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span>
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span>	if !concurrentSweep || mode == gcForceBlockMode {
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span>		<span class="comment">// Special case synchronous sweep.</span>
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span>		<span class="comment">// Record that no proportional sweeping has to happen.</span>
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span>		lock(&amp;mheap_.lock)
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>		mheap_.sweepPagesPerByte = 0
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>		unlock(&amp;mheap_.lock)
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>		<span class="comment">// Flush all mcaches.</span>
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>		for _, pp := range allp {
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>			pp.mcache.prepareForSweep()
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>		}
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>		<span class="comment">// Sweep all spans eagerly.</span>
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span>		for sweepone() != ^uintptr(0) {
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span>		}
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span>		<span class="comment">// Free workbufs eagerly.</span>
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span>		prepareFreeWorkbufs()
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span>		for freeSomeWbufs(false) {
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>		}
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>		<span class="comment">// All &#34;free&#34; events for this mark/sweep cycle have</span>
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span>		<span class="comment">// now happened, so we can make this profile cycle</span>
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span>		<span class="comment">// available immediately.</span>
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span>		mProf_NextCycle()
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span>		mProf_Flush()
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span>		return true
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span>	}
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span>
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span>	<span class="comment">// Background sweep.</span>
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span>	lock(&amp;sweep.lock)
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span>	if sweep.parked {
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span>		sweep.parked = false
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span>		ready(sweep.g, 0, true)
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>	}
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>	unlock(&amp;sweep.lock)
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>	return false
<span id="L1618" class="ln">  1618&nbsp;&nbsp;</span>}
<span id="L1619" class="ln">  1619&nbsp;&nbsp;</span>
<span id="L1620" class="ln">  1620&nbsp;&nbsp;</span><span class="comment">// gcResetMarkState resets global state prior to marking (concurrent</span>
<span id="L1621" class="ln">  1621&nbsp;&nbsp;</span><span class="comment">// or STW) and resets the stack scan state of all Gs.</span>
<span id="L1622" class="ln">  1622&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1623" class="ln">  1623&nbsp;&nbsp;</span><span class="comment">// This is safe to do without the world stopped because any Gs created</span>
<span id="L1624" class="ln">  1624&nbsp;&nbsp;</span><span class="comment">// during or after this will start out in the reset state.</span>
<span id="L1625" class="ln">  1625&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1626" class="ln">  1626&nbsp;&nbsp;</span><span class="comment">// gcResetMarkState must be called on the system stack because it acquires</span>
<span id="L1627" class="ln">  1627&nbsp;&nbsp;</span><span class="comment">// the heap lock. See mheap for details.</span>
<span id="L1628" class="ln">  1628&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1629" class="ln">  1629&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1630" class="ln">  1630&nbsp;&nbsp;</span>func gcResetMarkState() {
<span id="L1631" class="ln">  1631&nbsp;&nbsp;</span>	<span class="comment">// This may be called during a concurrent phase, so lock to make sure</span>
<span id="L1632" class="ln">  1632&nbsp;&nbsp;</span>	<span class="comment">// allgs doesn&#39;t change.</span>
<span id="L1633" class="ln">  1633&nbsp;&nbsp;</span>	forEachG(func(gp *g) {
<span id="L1634" class="ln">  1634&nbsp;&nbsp;</span>		gp.gcscandone = false <span class="comment">// set to true in gcphasework</span>
<span id="L1635" class="ln">  1635&nbsp;&nbsp;</span>		gp.gcAssistBytes = 0
<span id="L1636" class="ln">  1636&nbsp;&nbsp;</span>	})
<span id="L1637" class="ln">  1637&nbsp;&nbsp;</span>
<span id="L1638" class="ln">  1638&nbsp;&nbsp;</span>	<span class="comment">// Clear page marks. This is just 1MB per 64GB of heap, so the</span>
<span id="L1639" class="ln">  1639&nbsp;&nbsp;</span>	<span class="comment">// time here is pretty trivial.</span>
<span id="L1640" class="ln">  1640&nbsp;&nbsp;</span>	lock(&amp;mheap_.lock)
<span id="L1641" class="ln">  1641&nbsp;&nbsp;</span>	arenas := mheap_.allArenas
<span id="L1642" class="ln">  1642&nbsp;&nbsp;</span>	unlock(&amp;mheap_.lock)
<span id="L1643" class="ln">  1643&nbsp;&nbsp;</span>	for _, ai := range arenas {
<span id="L1644" class="ln">  1644&nbsp;&nbsp;</span>		ha := mheap_.arenas[ai.l1()][ai.l2()]
<span id="L1645" class="ln">  1645&nbsp;&nbsp;</span>		for i := range ha.pageMarks {
<span id="L1646" class="ln">  1646&nbsp;&nbsp;</span>			ha.pageMarks[i] = 0
<span id="L1647" class="ln">  1647&nbsp;&nbsp;</span>		}
<span id="L1648" class="ln">  1648&nbsp;&nbsp;</span>	}
<span id="L1649" class="ln">  1649&nbsp;&nbsp;</span>
<span id="L1650" class="ln">  1650&nbsp;&nbsp;</span>	work.bytesMarked = 0
<span id="L1651" class="ln">  1651&nbsp;&nbsp;</span>	work.initialHeapLive = gcController.heapLive.Load()
<span id="L1652" class="ln">  1652&nbsp;&nbsp;</span>}
<span id="L1653" class="ln">  1653&nbsp;&nbsp;</span>
<span id="L1654" class="ln">  1654&nbsp;&nbsp;</span><span class="comment">// Hooks for other packages</span>
<span id="L1655" class="ln">  1655&nbsp;&nbsp;</span>
<span id="L1656" class="ln">  1656&nbsp;&nbsp;</span>var poolcleanup func()
<span id="L1657" class="ln">  1657&nbsp;&nbsp;</span>var boringCaches []unsafe.Pointer <span class="comment">// for crypto/internal/boring</span>
<span id="L1658" class="ln">  1658&nbsp;&nbsp;</span>
<span id="L1659" class="ln">  1659&nbsp;&nbsp;</span><span class="comment">//go:linkname sync_runtime_registerPoolCleanup sync.runtime_registerPoolCleanup</span>
<span id="L1660" class="ln">  1660&nbsp;&nbsp;</span>func sync_runtime_registerPoolCleanup(f func()) {
<span id="L1661" class="ln">  1661&nbsp;&nbsp;</span>	poolcleanup = f
<span id="L1662" class="ln">  1662&nbsp;&nbsp;</span>}
<span id="L1663" class="ln">  1663&nbsp;&nbsp;</span>
<span id="L1664" class="ln">  1664&nbsp;&nbsp;</span><span class="comment">//go:linkname boring_registerCache crypto/internal/boring/bcache.registerCache</span>
<span id="L1665" class="ln">  1665&nbsp;&nbsp;</span>func boring_registerCache(p unsafe.Pointer) {
<span id="L1666" class="ln">  1666&nbsp;&nbsp;</span>	boringCaches = append(boringCaches, p)
<span id="L1667" class="ln">  1667&nbsp;&nbsp;</span>}
<span id="L1668" class="ln">  1668&nbsp;&nbsp;</span>
<span id="L1669" class="ln">  1669&nbsp;&nbsp;</span>func clearpools() {
<span id="L1670" class="ln">  1670&nbsp;&nbsp;</span>	<span class="comment">// clear sync.Pools</span>
<span id="L1671" class="ln">  1671&nbsp;&nbsp;</span>	if poolcleanup != nil {
<span id="L1672" class="ln">  1672&nbsp;&nbsp;</span>		poolcleanup()
<span id="L1673" class="ln">  1673&nbsp;&nbsp;</span>	}
<span id="L1674" class="ln">  1674&nbsp;&nbsp;</span>
<span id="L1675" class="ln">  1675&nbsp;&nbsp;</span>	<span class="comment">// clear boringcrypto caches</span>
<span id="L1676" class="ln">  1676&nbsp;&nbsp;</span>	for _, p := range boringCaches {
<span id="L1677" class="ln">  1677&nbsp;&nbsp;</span>		atomicstorep(p, nil)
<span id="L1678" class="ln">  1678&nbsp;&nbsp;</span>	}
<span id="L1679" class="ln">  1679&nbsp;&nbsp;</span>
<span id="L1680" class="ln">  1680&nbsp;&nbsp;</span>	<span class="comment">// Clear central sudog cache.</span>
<span id="L1681" class="ln">  1681&nbsp;&nbsp;</span>	<span class="comment">// Leave per-P caches alone, they have strictly bounded size.</span>
<span id="L1682" class="ln">  1682&nbsp;&nbsp;</span>	<span class="comment">// Disconnect cached list before dropping it on the floor,</span>
<span id="L1683" class="ln">  1683&nbsp;&nbsp;</span>	<span class="comment">// so that a dangling ref to one entry does not pin all of them.</span>
<span id="L1684" class="ln">  1684&nbsp;&nbsp;</span>	lock(&amp;sched.sudoglock)
<span id="L1685" class="ln">  1685&nbsp;&nbsp;</span>	var sg, sgnext *sudog
<span id="L1686" class="ln">  1686&nbsp;&nbsp;</span>	for sg = sched.sudogcache; sg != nil; sg = sgnext {
<span id="L1687" class="ln">  1687&nbsp;&nbsp;</span>		sgnext = sg.next
<span id="L1688" class="ln">  1688&nbsp;&nbsp;</span>		sg.next = nil
<span id="L1689" class="ln">  1689&nbsp;&nbsp;</span>	}
<span id="L1690" class="ln">  1690&nbsp;&nbsp;</span>	sched.sudogcache = nil
<span id="L1691" class="ln">  1691&nbsp;&nbsp;</span>	unlock(&amp;sched.sudoglock)
<span id="L1692" class="ln">  1692&nbsp;&nbsp;</span>
<span id="L1693" class="ln">  1693&nbsp;&nbsp;</span>	<span class="comment">// Clear central defer pool.</span>
<span id="L1694" class="ln">  1694&nbsp;&nbsp;</span>	<span class="comment">// Leave per-P pools alone, they have strictly bounded size.</span>
<span id="L1695" class="ln">  1695&nbsp;&nbsp;</span>	lock(&amp;sched.deferlock)
<span id="L1696" class="ln">  1696&nbsp;&nbsp;</span>	<span class="comment">// disconnect cached list before dropping it on the floor,</span>
<span id="L1697" class="ln">  1697&nbsp;&nbsp;</span>	<span class="comment">// so that a dangling ref to one entry does not pin all of them.</span>
<span id="L1698" class="ln">  1698&nbsp;&nbsp;</span>	var d, dlink *_defer
<span id="L1699" class="ln">  1699&nbsp;&nbsp;</span>	for d = sched.deferpool; d != nil; d = dlink {
<span id="L1700" class="ln">  1700&nbsp;&nbsp;</span>		dlink = d.link
<span id="L1701" class="ln">  1701&nbsp;&nbsp;</span>		d.link = nil
<span id="L1702" class="ln">  1702&nbsp;&nbsp;</span>	}
<span id="L1703" class="ln">  1703&nbsp;&nbsp;</span>	sched.deferpool = nil
<span id="L1704" class="ln">  1704&nbsp;&nbsp;</span>	unlock(&amp;sched.deferlock)
<span id="L1705" class="ln">  1705&nbsp;&nbsp;</span>}
<span id="L1706" class="ln">  1706&nbsp;&nbsp;</span>
<span id="L1707" class="ln">  1707&nbsp;&nbsp;</span><span class="comment">// Timing</span>
<span id="L1708" class="ln">  1708&nbsp;&nbsp;</span>
<span id="L1709" class="ln">  1709&nbsp;&nbsp;</span><span class="comment">// itoaDiv formats val/(10**dec) into buf.</span>
<span id="L1710" class="ln">  1710&nbsp;&nbsp;</span>func itoaDiv(buf []byte, val uint64, dec int) []byte {
<span id="L1711" class="ln">  1711&nbsp;&nbsp;</span>	i := len(buf) - 1
<span id="L1712" class="ln">  1712&nbsp;&nbsp;</span>	idec := i - dec
<span id="L1713" class="ln">  1713&nbsp;&nbsp;</span>	for val &gt;= 10 || i &gt;= idec {
<span id="L1714" class="ln">  1714&nbsp;&nbsp;</span>		buf[i] = byte(val%10 + &#39;0&#39;)
<span id="L1715" class="ln">  1715&nbsp;&nbsp;</span>		i--
<span id="L1716" class="ln">  1716&nbsp;&nbsp;</span>		if i == idec {
<span id="L1717" class="ln">  1717&nbsp;&nbsp;</span>			buf[i] = &#39;.&#39;
<span id="L1718" class="ln">  1718&nbsp;&nbsp;</span>			i--
<span id="L1719" class="ln">  1719&nbsp;&nbsp;</span>		}
<span id="L1720" class="ln">  1720&nbsp;&nbsp;</span>		val /= 10
<span id="L1721" class="ln">  1721&nbsp;&nbsp;</span>	}
<span id="L1722" class="ln">  1722&nbsp;&nbsp;</span>	buf[i] = byte(val + &#39;0&#39;)
<span id="L1723" class="ln">  1723&nbsp;&nbsp;</span>	return buf[i:]
<span id="L1724" class="ln">  1724&nbsp;&nbsp;</span>}
<span id="L1725" class="ln">  1725&nbsp;&nbsp;</span>
<span id="L1726" class="ln">  1726&nbsp;&nbsp;</span><span class="comment">// fmtNSAsMS nicely formats ns nanoseconds as milliseconds.</span>
<span id="L1727" class="ln">  1727&nbsp;&nbsp;</span>func fmtNSAsMS(buf []byte, ns uint64) []byte {
<span id="L1728" class="ln">  1728&nbsp;&nbsp;</span>	if ns &gt;= 10e6 {
<span id="L1729" class="ln">  1729&nbsp;&nbsp;</span>		<span class="comment">// Format as whole milliseconds.</span>
<span id="L1730" class="ln">  1730&nbsp;&nbsp;</span>		return itoaDiv(buf, ns/1e6, 0)
<span id="L1731" class="ln">  1731&nbsp;&nbsp;</span>	}
<span id="L1732" class="ln">  1732&nbsp;&nbsp;</span>	<span class="comment">// Format two digits of precision, with at most three decimal places.</span>
<span id="L1733" class="ln">  1733&nbsp;&nbsp;</span>	x := ns / 1e3
<span id="L1734" class="ln">  1734&nbsp;&nbsp;</span>	if x == 0 {
<span id="L1735" class="ln">  1735&nbsp;&nbsp;</span>		buf[0] = &#39;0&#39;
<span id="L1736" class="ln">  1736&nbsp;&nbsp;</span>		return buf[:1]
<span id="L1737" class="ln">  1737&nbsp;&nbsp;</span>	}
<span id="L1738" class="ln">  1738&nbsp;&nbsp;</span>	dec := 3
<span id="L1739" class="ln">  1739&nbsp;&nbsp;</span>	for x &gt;= 100 {
<span id="L1740" class="ln">  1740&nbsp;&nbsp;</span>		x /= 10
<span id="L1741" class="ln">  1741&nbsp;&nbsp;</span>		dec--
<span id="L1742" class="ln">  1742&nbsp;&nbsp;</span>	}
<span id="L1743" class="ln">  1743&nbsp;&nbsp;</span>	return itoaDiv(buf, x, dec)
<span id="L1744" class="ln">  1744&nbsp;&nbsp;</span>}
<span id="L1745" class="ln">  1745&nbsp;&nbsp;</span>
<span id="L1746" class="ln">  1746&nbsp;&nbsp;</span><span class="comment">// Helpers for testing GC.</span>
<span id="L1747" class="ln">  1747&nbsp;&nbsp;</span>
<span id="L1748" class="ln">  1748&nbsp;&nbsp;</span><span class="comment">// gcTestMoveStackOnNextCall causes the stack to be moved on a call</span>
<span id="L1749" class="ln">  1749&nbsp;&nbsp;</span><span class="comment">// immediately following the call to this. It may not work correctly</span>
<span id="L1750" class="ln">  1750&nbsp;&nbsp;</span><span class="comment">// if any other work appears after this call (such as returning).</span>
<span id="L1751" class="ln">  1751&nbsp;&nbsp;</span><span class="comment">// Typically the following call should be marked go:noinline so it</span>
<span id="L1752" class="ln">  1752&nbsp;&nbsp;</span><span class="comment">// performs a stack check.</span>
<span id="L1753" class="ln">  1753&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1754" class="ln">  1754&nbsp;&nbsp;</span><span class="comment">// In rare cases this may not cause the stack to move, specifically if</span>
<span id="L1755" class="ln">  1755&nbsp;&nbsp;</span><span class="comment">// there&#39;s a preemption between this call and the next.</span>
<span id="L1756" class="ln">  1756&nbsp;&nbsp;</span>func gcTestMoveStackOnNextCall() {
<span id="L1757" class="ln">  1757&nbsp;&nbsp;</span>	gp := getg()
<span id="L1758" class="ln">  1758&nbsp;&nbsp;</span>	gp.stackguard0 = stackForceMove
<span id="L1759" class="ln">  1759&nbsp;&nbsp;</span>}
<span id="L1760" class="ln">  1760&nbsp;&nbsp;</span>
<span id="L1761" class="ln">  1761&nbsp;&nbsp;</span><span class="comment">// gcTestIsReachable performs a GC and returns a bit set where bit i</span>
<span id="L1762" class="ln">  1762&nbsp;&nbsp;</span><span class="comment">// is set if ptrs[i] is reachable.</span>
<span id="L1763" class="ln">  1763&nbsp;&nbsp;</span>func gcTestIsReachable(ptrs ...unsafe.Pointer) (mask uint64) {
<span id="L1764" class="ln">  1764&nbsp;&nbsp;</span>	<span class="comment">// This takes the pointers as unsafe.Pointers in order to keep</span>
<span id="L1765" class="ln">  1765&nbsp;&nbsp;</span>	<span class="comment">// them live long enough for us to attach specials. After</span>
<span id="L1766" class="ln">  1766&nbsp;&nbsp;</span>	<span class="comment">// that, we drop our references to them.</span>
<span id="L1767" class="ln">  1767&nbsp;&nbsp;</span>
<span id="L1768" class="ln">  1768&nbsp;&nbsp;</span>	if len(ptrs) &gt; 64 {
<span id="L1769" class="ln">  1769&nbsp;&nbsp;</span>		panic(&#34;too many pointers for uint64 mask&#34;)
<span id="L1770" class="ln">  1770&nbsp;&nbsp;</span>	}
<span id="L1771" class="ln">  1771&nbsp;&nbsp;</span>
<span id="L1772" class="ln">  1772&nbsp;&nbsp;</span>	<span class="comment">// Block GC while we attach specials and drop our references</span>
<span id="L1773" class="ln">  1773&nbsp;&nbsp;</span>	<span class="comment">// to ptrs. Otherwise, if a GC is in progress, it could mark</span>
<span id="L1774" class="ln">  1774&nbsp;&nbsp;</span>	<span class="comment">// them reachable via this function before we have a chance to</span>
<span id="L1775" class="ln">  1775&nbsp;&nbsp;</span>	<span class="comment">// drop them.</span>
<span id="L1776" class="ln">  1776&nbsp;&nbsp;</span>	semacquire(&amp;gcsema)
<span id="L1777" class="ln">  1777&nbsp;&nbsp;</span>
<span id="L1778" class="ln">  1778&nbsp;&nbsp;</span>	<span class="comment">// Create reachability specials for ptrs.</span>
<span id="L1779" class="ln">  1779&nbsp;&nbsp;</span>	specials := make([]*specialReachable, len(ptrs))
<span id="L1780" class="ln">  1780&nbsp;&nbsp;</span>	for i, p := range ptrs {
<span id="L1781" class="ln">  1781&nbsp;&nbsp;</span>		lock(&amp;mheap_.speciallock)
<span id="L1782" class="ln">  1782&nbsp;&nbsp;</span>		s := (*specialReachable)(mheap_.specialReachableAlloc.alloc())
<span id="L1783" class="ln">  1783&nbsp;&nbsp;</span>		unlock(&amp;mheap_.speciallock)
<span id="L1784" class="ln">  1784&nbsp;&nbsp;</span>		s.special.kind = _KindSpecialReachable
<span id="L1785" class="ln">  1785&nbsp;&nbsp;</span>		if !addspecial(p, &amp;s.special) {
<span id="L1786" class="ln">  1786&nbsp;&nbsp;</span>			throw(&#34;already have a reachable special (duplicate pointer?)&#34;)
<span id="L1787" class="ln">  1787&nbsp;&nbsp;</span>		}
<span id="L1788" class="ln">  1788&nbsp;&nbsp;</span>		specials[i] = s
<span id="L1789" class="ln">  1789&nbsp;&nbsp;</span>		<span class="comment">// Make sure we don&#39;t retain ptrs.</span>
<span id="L1790" class="ln">  1790&nbsp;&nbsp;</span>		ptrs[i] = nil
<span id="L1791" class="ln">  1791&nbsp;&nbsp;</span>	}
<span id="L1792" class="ln">  1792&nbsp;&nbsp;</span>
<span id="L1793" class="ln">  1793&nbsp;&nbsp;</span>	semrelease(&amp;gcsema)
<span id="L1794" class="ln">  1794&nbsp;&nbsp;</span>
<span id="L1795" class="ln">  1795&nbsp;&nbsp;</span>	<span class="comment">// Force a full GC and sweep.</span>
<span id="L1796" class="ln">  1796&nbsp;&nbsp;</span>	GC()
<span id="L1797" class="ln">  1797&nbsp;&nbsp;</span>
<span id="L1798" class="ln">  1798&nbsp;&nbsp;</span>	<span class="comment">// Process specials.</span>
<span id="L1799" class="ln">  1799&nbsp;&nbsp;</span>	for i, s := range specials {
<span id="L1800" class="ln">  1800&nbsp;&nbsp;</span>		if !s.done {
<span id="L1801" class="ln">  1801&nbsp;&nbsp;</span>			printlock()
<span id="L1802" class="ln">  1802&nbsp;&nbsp;</span>			println(&#34;runtime: object&#34;, i, &#34;was not swept&#34;)
<span id="L1803" class="ln">  1803&nbsp;&nbsp;</span>			throw(&#34;IsReachable failed&#34;)
<span id="L1804" class="ln">  1804&nbsp;&nbsp;</span>		}
<span id="L1805" class="ln">  1805&nbsp;&nbsp;</span>		if s.reachable {
<span id="L1806" class="ln">  1806&nbsp;&nbsp;</span>			mask |= 1 &lt;&lt; i
<span id="L1807" class="ln">  1807&nbsp;&nbsp;</span>		}
<span id="L1808" class="ln">  1808&nbsp;&nbsp;</span>		lock(&amp;mheap_.speciallock)
<span id="L1809" class="ln">  1809&nbsp;&nbsp;</span>		mheap_.specialReachableAlloc.free(unsafe.Pointer(s))
<span id="L1810" class="ln">  1810&nbsp;&nbsp;</span>		unlock(&amp;mheap_.speciallock)
<span id="L1811" class="ln">  1811&nbsp;&nbsp;</span>	}
<span id="L1812" class="ln">  1812&nbsp;&nbsp;</span>
<span id="L1813" class="ln">  1813&nbsp;&nbsp;</span>	return mask
<span id="L1814" class="ln">  1814&nbsp;&nbsp;</span>}
<span id="L1815" class="ln">  1815&nbsp;&nbsp;</span>
<span id="L1816" class="ln">  1816&nbsp;&nbsp;</span><span class="comment">// gcTestPointerClass returns the category of what p points to, one of:</span>
<span id="L1817" class="ln">  1817&nbsp;&nbsp;</span><span class="comment">// &#34;heap&#34;, &#34;stack&#34;, &#34;data&#34;, &#34;bss&#34;, &#34;other&#34;. This is useful for checking</span>
<span id="L1818" class="ln">  1818&nbsp;&nbsp;</span><span class="comment">// that a test is doing what it&#39;s intended to do.</span>
<span id="L1819" class="ln">  1819&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1820" class="ln">  1820&nbsp;&nbsp;</span><span class="comment">// This is nosplit simply to avoid extra pointer shuffling that may</span>
<span id="L1821" class="ln">  1821&nbsp;&nbsp;</span><span class="comment">// complicate a test.</span>
<span id="L1822" class="ln">  1822&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1823" class="ln">  1823&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1824" class="ln">  1824&nbsp;&nbsp;</span>func gcTestPointerClass(p unsafe.Pointer) string {
<span id="L1825" class="ln">  1825&nbsp;&nbsp;</span>	p2 := uintptr(noescape(p))
<span id="L1826" class="ln">  1826&nbsp;&nbsp;</span>	gp := getg()
<span id="L1827" class="ln">  1827&nbsp;&nbsp;</span>	if gp.stack.lo &lt;= p2 &amp;&amp; p2 &lt; gp.stack.hi {
<span id="L1828" class="ln">  1828&nbsp;&nbsp;</span>		return &#34;stack&#34;
<span id="L1829" class="ln">  1829&nbsp;&nbsp;</span>	}
<span id="L1830" class="ln">  1830&nbsp;&nbsp;</span>	if base, _, _ := findObject(p2, 0, 0); base != 0 {
<span id="L1831" class="ln">  1831&nbsp;&nbsp;</span>		return &#34;heap&#34;
<span id="L1832" class="ln">  1832&nbsp;&nbsp;</span>	}
<span id="L1833" class="ln">  1833&nbsp;&nbsp;</span>	for _, datap := range activeModules() {
<span id="L1834" class="ln">  1834&nbsp;&nbsp;</span>		if datap.data &lt;= p2 &amp;&amp; p2 &lt; datap.edata || datap.noptrdata &lt;= p2 &amp;&amp; p2 &lt; datap.enoptrdata {
<span id="L1835" class="ln">  1835&nbsp;&nbsp;</span>			return &#34;data&#34;
<span id="L1836" class="ln">  1836&nbsp;&nbsp;</span>		}
<span id="L1837" class="ln">  1837&nbsp;&nbsp;</span>		if datap.bss &lt;= p2 &amp;&amp; p2 &lt; datap.ebss || datap.noptrbss &lt;= p2 &amp;&amp; p2 &lt;= datap.enoptrbss {
<span id="L1838" class="ln">  1838&nbsp;&nbsp;</span>			return &#34;bss&#34;
<span id="L1839" class="ln">  1839&nbsp;&nbsp;</span>		}
<span id="L1840" class="ln">  1840&nbsp;&nbsp;</span>	}
<span id="L1841" class="ln">  1841&nbsp;&nbsp;</span>	KeepAlive(p)
<span id="L1842" class="ln">  1842&nbsp;&nbsp;</span>	return &#34;other&#34;
<span id="L1843" class="ln">  1843&nbsp;&nbsp;</span>}
<span id="L1844" class="ln">  1844&nbsp;&nbsp;</span>
</pre><p><a href="mgc.go?m=text">View as plain text</a></p>

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
