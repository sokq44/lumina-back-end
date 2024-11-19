<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mgcmark.go - Go Documentation Server</title>

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
<a href="mgcmark.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mgcmark.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Garbage collector: marking and scanning</span>
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
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>const (
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	fixedRootFinalizers = iota
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	fixedRootFreeGStacks
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	fixedRootCount
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// rootBlockBytes is the number of bytes to scan per data or</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// BSS root.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	rootBlockBytes = 256 &lt;&lt; 10
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// maxObletBytes is the maximum bytes of an object to scan at</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// once. Larger objects will be split up into &#34;oblets&#34; of at</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// most this size. Since we can scan 1–2 MB/ms, 128 KB bounds</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// scan preemption at ~100 µs.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// This must be &gt; _MaxSmallSize so that the object base is the</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// span base.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	maxObletBytes = 128 &lt;&lt; 10
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// drainCheckThreshold specifies how many units of work to do</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// between self-preemption checks in gcDrain. Assuming a scan</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// rate of 1 MB/ms, this is ~100 µs. Lower values have higher</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// overhead in the scan loop (the scheduler check may perform</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// a syscall, so its overhead is nontrivial). Higher values</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// make the system less responsive to incoming work.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	drainCheckThreshold = 100000
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// pagesPerSpanRoot indicates how many pages to scan from a span root</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// at a time. Used by special root marking.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// Higher values improve throughput by increasing locality, but</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// increase the minimum latency of a marking operation.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// Must be a multiple of the pageInUse bitmap element size and</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// must also evenly divide pagesPerArena.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	pagesPerSpanRoot = 512
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// gcMarkRootPrepare queues root scanning jobs (stacks, globals, and</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// some miscellany) and initializes scanning-related state.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// The world must be stopped.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>func gcMarkRootPrepare() {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// Compute how many data and BSS root blocks there are.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	nBlocks := func(bytes uintptr) int {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		return int(divRoundUp(bytes, rootBlockBytes))
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	work.nDataRoots = 0
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	work.nBSSRoots = 0
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// Scan globals.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	for _, datap := range activeModules() {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		nDataRoots := nBlocks(datap.edata - datap.data)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		if nDataRoots &gt; work.nDataRoots {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			work.nDataRoots = nDataRoots
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	for _, datap := range activeModules() {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		nBSSRoots := nBlocks(datap.ebss - datap.bss)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		if nBSSRoots &gt; work.nBSSRoots {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			work.nBSSRoots = nBSSRoots
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// Scan span roots for finalizer specials.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// We depend on addfinalizer to mark objects that get</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// finalizers after root marking.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// We&#39;re going to scan the whole heap (that was available at the time the</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// mark phase started, i.e. markArenas) for in-use spans which have specials.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// Break up the work into arenas, and further into chunks.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// Snapshot allArenas as markArenas. This snapshot is safe because allArenas</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// is append-only.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	mheap_.markArenas = mheap_.allArenas[:len(mheap_.allArenas):len(mheap_.allArenas)]
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	work.nSpanRoots = len(mheap_.markArenas) * (pagesPerArena / pagesPerSpanRoot)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// Scan stacks.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// Gs may be created after this point, but it&#39;s okay that we</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// ignore them because they begin life without any roots, so</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// there&#39;s nothing to scan, and any roots they create during</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// the concurrent phase will be caught by the write barrier.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	work.stackRoots = allGsSnapshot()
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	work.nStackRoots = len(work.stackRoots)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	work.markrootNext = 0
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	work.markrootJobs = uint32(fixedRootCount + work.nDataRoots + work.nBSSRoots + work.nSpanRoots + work.nStackRoots)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// Calculate base indexes of each root type</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	work.baseData = uint32(fixedRootCount)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	work.baseBSS = work.baseData + uint32(work.nDataRoots)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	work.baseSpans = work.baseBSS + uint32(work.nBSSRoots)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	work.baseStacks = work.baseSpans + uint32(work.nSpanRoots)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	work.baseEnd = work.baseStacks + uint32(work.nStackRoots)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// gcMarkRootCheck checks that all roots have been scanned. It is</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// purely for debugging.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>func gcMarkRootCheck() {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if work.markrootNext &lt; work.markrootJobs {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		print(work.markrootNext, &#34; of &#34;, work.markrootJobs, &#34; markroot jobs done\n&#34;)
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		throw(&#34;left over markroot jobs&#34;)
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// Check that stacks have been scanned.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// We only check the first nStackRoots Gs that we should have scanned.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// Since we don&#39;t care about newer Gs (see comment in</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// gcMarkRootPrepare), no locking is required.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	i := 0
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	forEachGRace(func(gp *g) {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		if i &gt;= work.nStackRoots {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			return
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		if !gp.gcscandone {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			println(&#34;gp&#34;, gp, &#34;goid&#34;, gp.goid,
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>				&#34;status&#34;, readgstatus(gp),
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>				&#34;gcscandone&#34;, gp.gcscandone)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			throw(&#34;scan missed a g&#34;)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		i++
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	})
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// ptrmask for an allocation containing a single pointer.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>var oneptrmask = [...]uint8{1}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">// markroot scans the i&#39;th root.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// Preemption must be disabled (because this uses a gcWork).</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// Returns the amount of GC work credit produced by the operation.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">// If flushBgCredit is true, then that credit is also flushed</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// to the background credit pool.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">// nowritebarrier is only advisory here.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>func markroot(gcw *gcWork, i uint32, flushBgCredit bool) int64 {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// Note: if you add a case here, please also update heapdump.go:dumproots.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	var workDone int64
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	var workCounter *atomic.Int64
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	switch {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	case work.baseData &lt;= i &amp;&amp; i &lt; work.baseBSS:
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		workCounter = &amp;gcController.globalsScanWork
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		for _, datap := range activeModules() {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			workDone += markrootBlock(datap.data, datap.edata-datap.data, datap.gcdatamask.bytedata, gcw, int(i-work.baseData))
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	case work.baseBSS &lt;= i &amp;&amp; i &lt; work.baseSpans:
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		workCounter = &amp;gcController.globalsScanWork
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		for _, datap := range activeModules() {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			workDone += markrootBlock(datap.bss, datap.ebss-datap.bss, datap.gcbssmask.bytedata, gcw, int(i-work.baseBSS))
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	case i == fixedRootFinalizers:
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		for fb := allfin; fb != nil; fb = fb.alllink {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			cnt := uintptr(atomic.Load(&amp;fb.cnt))
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			scanblock(uintptr(unsafe.Pointer(&amp;fb.fin[0])), cnt*unsafe.Sizeof(fb.fin[0]), &amp;finptrmask[0], gcw, nil)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	case i == fixedRootFreeGStacks:
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		<span class="comment">// Switch to the system stack so we can call</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		<span class="comment">// stackfree.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		systemstack(markrootFreeGStacks)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	case work.baseSpans &lt;= i &amp;&amp; i &lt; work.baseStacks:
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		<span class="comment">// mark mspan.specials</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		markrootSpans(gcw, int(i-work.baseSpans))
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	default:
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		<span class="comment">// the rest is scanning goroutine stacks</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		workCounter = &amp;gcController.stackScanWork
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		if i &lt; work.baseStacks || work.baseEnd &lt;= i {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			printlock()
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			print(&#34;runtime: markroot index &#34;, i, &#34; not in stack roots range [&#34;, work.baseStacks, &#34;, &#34;, work.baseEnd, &#34;)\n&#34;)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>			throw(&#34;markroot: bad index&#34;)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		gp := work.stackRoots[i-work.baseStacks]
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		<span class="comment">// remember when we&#39;ve first observed the G blocked</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		<span class="comment">// needed only to output in traceback</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		status := readgstatus(gp) <span class="comment">// We are not in a scan state</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		if (status == _Gwaiting || status == _Gsyscall) &amp;&amp; gp.waitsince == 0 {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			gp.waitsince = work.tstart
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		<span class="comment">// scanstack must be done on the system stack in case</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		<span class="comment">// we&#39;re trying to scan our own stack.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			<span class="comment">// If this is a self-scan, put the user G in</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			<span class="comment">// _Gwaiting to prevent self-deadlock. It may</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			<span class="comment">// already be in _Gwaiting if this is a mark</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			<span class="comment">// worker or we&#39;re in mark termination.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			userG := getg().m.curg
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			selfScan := gp == userG &amp;&amp; readgstatus(userG) == _Grunning
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			if selfScan {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>				casGToWaiting(userG, _Grunning, waitReasonGarbageCollectionScan)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			<span class="comment">// TODO: suspendG blocks (and spins) until gp</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			<span class="comment">// stops, which may take a while for</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			<span class="comment">// running goroutines. Consider doing this in</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			<span class="comment">// two phases where the first is non-blocking:</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			<span class="comment">// we scan the stacks we can and ask running</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			<span class="comment">// goroutines to scan themselves; and the</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			<span class="comment">// second blocks.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			stopped := suspendG(gp)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			if stopped.dead {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>				gp.gcscandone = true
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>				return
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			if gp.gcscandone {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>				throw(&#34;g already scanned&#34;)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			workDone += scanstack(gp, gcw)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			gp.gcscandone = true
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			resumeG(stopped)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			if selfScan {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>				casgstatus(userG, _Gwaiting, _Grunning)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		})
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	if workCounter != nil &amp;&amp; workDone != 0 {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		workCounter.Add(workDone)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		if flushBgCredit {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			gcFlushBgCredit(workDone)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	return workDone
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">// markrootBlock scans the shard&#39;th shard of the block of memory [b0,</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// b0+n0), with the given pointer mask.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span><span class="comment">// Returns the amount of work done.</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>func markrootBlock(b0, n0 uintptr, ptrmask0 *uint8, gcw *gcWork, shard int) int64 {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	if rootBlockBytes%(8*goarch.PtrSize) != 0 {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		<span class="comment">// This is necessary to pick byte offsets in ptrmask0.</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		throw(&#34;rootBlockBytes must be a multiple of 8*ptrSize&#34;)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	<span class="comment">// Note that if b0 is toward the end of the address space,</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	<span class="comment">// then b0 + rootBlockBytes might wrap around.</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	<span class="comment">// These tests are written to avoid any possible overflow.</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	off := uintptr(shard) * rootBlockBytes
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	if off &gt;= n0 {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		return 0
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	b := b0 + off
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	ptrmask := (*uint8)(add(unsafe.Pointer(ptrmask0), uintptr(shard)*(rootBlockBytes/(8*goarch.PtrSize))))
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	n := uintptr(rootBlockBytes)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	if off+n &gt; n0 {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		n = n0 - off
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">// Scan this shard.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	scanblock(b, n, ptrmask, gcw, nil)
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	return int64(n)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">// markrootFreeGStacks frees stacks of dead Gs.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">// This does not free stacks of dead Gs cached on Ps, but having a few</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span><span class="comment">// cached stacks around isn&#39;t a problem.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>func markrootFreeGStacks() {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// Take list of dead Gs with stacks.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	lock(&amp;sched.gFree.lock)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	list := sched.gFree.stack
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	sched.gFree.stack = gList{}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	unlock(&amp;sched.gFree.lock)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	if list.empty() {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		return
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	<span class="comment">// Free stacks.</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	q := gQueue{list.head, list.head}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	for gp := list.head.ptr(); gp != nil; gp = gp.schedlink.ptr() {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		stackfree(gp.stack)
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		gp.stack.lo = 0
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		gp.stack.hi = 0
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		<span class="comment">// Manipulate the queue directly since the Gs are</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		<span class="comment">// already all linked the right way.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		q.tail.set(gp)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	<span class="comment">// Put Gs back on the free list.</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	lock(&amp;sched.gFree.lock)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	sched.gFree.noStack.pushAll(q)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	unlock(&amp;sched.gFree.lock)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span><span class="comment">// markrootSpans marks roots for one shard of markArenas.</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>func markrootSpans(gcw *gcWork, shard int) {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// Objects with finalizers have two GC-related invariants:</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// 1) Everything reachable from the object must be marked.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	<span class="comment">// This ensures that when we pass the object to its finalizer,</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	<span class="comment">// everything the finalizer can reach will be retained.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	<span class="comment">// 2) Finalizer specials (which are not in the garbage</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">// collected heap) are roots. In practice, this means the fn</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// field must be scanned.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	sg := mheap_.sweepgen
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// Find the arena and page index into that arena for this shard.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	ai := mheap_.markArenas[shard/(pagesPerArena/pagesPerSpanRoot)]
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	ha := mheap_.arenas[ai.l1()][ai.l2()]
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	arenaPage := uint(uintptr(shard) * pagesPerSpanRoot % pagesPerArena)
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	<span class="comment">// Construct slice of bitmap which we&#39;ll iterate over.</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	specialsbits := ha.pageSpecials[arenaPage/8:]
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	specialsbits = specialsbits[:pagesPerSpanRoot/8]
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	for i := range specialsbits {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		<span class="comment">// Find set bits, which correspond to spans with specials.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		specials := atomic.Load8(&amp;specialsbits[i])
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		if specials == 0 {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>			continue
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		for j := uint(0); j &lt; 8; j++ {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>			if specials&amp;(1&lt;&lt;j) == 0 {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>				continue
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>			}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>			<span class="comment">// Find the span for this bit.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			<span class="comment">// This value is guaranteed to be non-nil because having</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			<span class="comment">// specials implies that the span is in-use, and since we&#39;re</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			<span class="comment">// currently marking we can be sure that we don&#39;t have to worry</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			<span class="comment">// about the span being freed and re-used.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			s := ha.spans[arenaPage+uint(i)*8+j]
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>			<span class="comment">// The state must be mSpanInUse if the specials bit is set, so</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			<span class="comment">// sanity check that.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			if state := s.state.get(); state != mSpanInUse {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>				print(&#34;s.state = &#34;, state, &#34;\n&#34;)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>				throw(&#34;non in-use span found with specials bit set&#34;)
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			<span class="comment">// Check that this span was swept (it may be cached or uncached).</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			if !useCheckmark &amp;&amp; !(s.sweepgen == sg || s.sweepgen == sg+3) {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>				<span class="comment">// sweepgen was updated (+2) during non-checkmark GC pass</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>				print(&#34;sweep &#34;, s.sweepgen, &#34; &#34;, sg, &#34;\n&#34;)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>				throw(&#34;gc: unswept span&#34;)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			<span class="comment">// Lock the specials to prevent a special from being</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			<span class="comment">// removed from the list while we&#39;re traversing it.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			lock(&amp;s.speciallock)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			for sp := s.specials; sp != nil; sp = sp.next {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>				if sp.kind != _KindSpecialFinalizer {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>					continue
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>				}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>				<span class="comment">// don&#39;t mark finalized object, but scan it so we</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>				<span class="comment">// retain everything it points to.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>				spf := (*specialfinalizer)(unsafe.Pointer(sp))
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>				<span class="comment">// A finalizer can be set for an inner byte of an object, find object beginning.</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>				p := s.base() + uintptr(spf.special.offset)/s.elemsize*s.elemsize
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>				<span class="comment">// Mark everything that can be reached from</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>				<span class="comment">// the object (but *not* the object itself or</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>				<span class="comment">// we&#39;ll never collect it).</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>				if !s.spanclass.noscan() {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>					scanobject(p, gcw)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>				}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>				<span class="comment">// The special itself is a root.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>				scanblock(uintptr(unsafe.Pointer(&amp;spf.fn)), goarch.PtrSize, &amp;oneptrmask[0], gcw, nil)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			}
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			unlock(&amp;s.speciallock)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span><span class="comment">// gcAssistAlloc performs GC work to make gp&#39;s assist debt positive.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span><span class="comment">// gp must be the calling user goroutine.</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span><span class="comment">// This must be called with preemption enabled.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>func gcAssistAlloc(gp *g) {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t assist in non-preemptible contexts. These are</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	<span class="comment">// generally fragile and won&#39;t allow the assist to block.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	if getg() == gp.m.g0 {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		return
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	if mp := getg().m; mp.locks &gt; 0 || mp.preemptoff != &#34;&#34; {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		return
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	<span class="comment">// This extremely verbose boolean indicates whether we&#39;ve</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	<span class="comment">// entered mark assist from the perspective of the tracer.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	<span class="comment">// In the old tracer, this is just before we call gcAssistAlloc1</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	<span class="comment">// *and* tracing is enabled. Because the old tracer doesn&#39;t</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	<span class="comment">// do any extra tracking, we need to be careful to not emit an</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	<span class="comment">// &#34;end&#34; event if there was no corresponding &#34;begin&#34; for the</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	<span class="comment">// mark assist.</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	<span class="comment">// In the new tracer, this is just before we call gcAssistAlloc1</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	<span class="comment">// *regardless* of whether tracing is enabled. This is because</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	<span class="comment">// the new tracer allows for tracing to begin (and advance</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	<span class="comment">// generations) in the middle of a GC mark phase, so we need to</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	<span class="comment">// record some state so that the tracer can pick it up to ensure</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	<span class="comment">// a consistent trace result.</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Hide the details of inMarkAssist in tracer</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	<span class="comment">// functions and simplify all the state tracking. This is a lot.</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	enteredMarkAssistForTracing := false
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>retry:
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	if gcCPULimiter.limiting() {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		<span class="comment">// If the CPU limiter is enabled, intentionally don&#39;t</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		<span class="comment">// assist to reduce the amount of CPU time spent in the GC.</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		if enteredMarkAssistForTracing {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>			trace := traceAcquire()
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			if trace.ok() {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>				trace.GCMarkAssistDone()
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>				<span class="comment">// Set this *after* we trace the end to make sure</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>				<span class="comment">// that we emit an in-progress event if this is</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>				<span class="comment">// the first event for the goroutine in the trace</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>				<span class="comment">// or trace generation. Also, do this between</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>				<span class="comment">// acquire/release because this is part of the</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>				<span class="comment">// goroutine&#39;s trace state, and it must be atomic</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>				<span class="comment">// with respect to the tracer.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>				gp.inMarkAssist = false
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>				traceRelease(trace)
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			} else {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>				<span class="comment">// This state is tracked even if tracing isn&#39;t enabled.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>				<span class="comment">// It&#39;s only used by the new tracer.</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>				<span class="comment">// See the comment on enteredMarkAssistForTracing.</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>				gp.inMarkAssist = false
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		return
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	<span class="comment">// Compute the amount of scan work we need to do to make the</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	<span class="comment">// balance positive. When the required amount of work is low,</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	<span class="comment">// we over-assist to build up credit for future allocations</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	<span class="comment">// and amortize the cost of assisting.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	assistWorkPerByte := gcController.assistWorkPerByte.Load()
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	assistBytesPerWork := gcController.assistBytesPerWork.Load()
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	debtBytes := -gp.gcAssistBytes
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	scanWork := int64(assistWorkPerByte * float64(debtBytes))
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	if scanWork &lt; gcOverAssistWork {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		scanWork = gcOverAssistWork
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		debtBytes = int64(assistBytesPerWork * float64(scanWork))
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	<span class="comment">// Steal as much credit as we can from the background GC&#39;s</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	<span class="comment">// scan credit. This is racy and may drop the background</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	<span class="comment">// credit below 0 if two mutators steal at the same time. This</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	<span class="comment">// will just cause steals to fail until credit is accumulated</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	<span class="comment">// again, so in the long run it doesn&#39;t really matter, but we</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	<span class="comment">// do have to handle the negative credit case.</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	bgScanCredit := gcController.bgScanCredit.Load()
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	stolen := int64(0)
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	if bgScanCredit &gt; 0 {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		if bgScanCredit &lt; scanWork {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			stolen = bgScanCredit
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			gp.gcAssistBytes += 1 + int64(assistBytesPerWork*float64(stolen))
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		} else {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			stolen = scanWork
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>			gp.gcAssistBytes += debtBytes
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		gcController.bgScanCredit.Add(-stolen)
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		scanWork -= stolen
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		if scanWork == 0 {
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>			<span class="comment">// We were able to steal all of the credit we</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>			<span class="comment">// needed.</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			if enteredMarkAssistForTracing {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>				trace := traceAcquire()
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>				if trace.ok() {
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>					trace.GCMarkAssistDone()
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>					<span class="comment">// Set this *after* we trace the end to make sure</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>					<span class="comment">// that we emit an in-progress event if this is</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>					<span class="comment">// the first event for the goroutine in the trace</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>					<span class="comment">// or trace generation. Also, do this between</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>					<span class="comment">// acquire/release because this is part of the</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>					<span class="comment">// goroutine&#39;s trace state, and it must be atomic</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>					<span class="comment">// with respect to the tracer.</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>					gp.inMarkAssist = false
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>					traceRelease(trace)
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>				} else {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>					<span class="comment">// This state is tracked even if tracing isn&#39;t enabled.</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>					<span class="comment">// It&#39;s only used by the new tracer.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>					<span class="comment">// See the comment on enteredMarkAssistForTracing.</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>					gp.inMarkAssist = false
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>				}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			return
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		}
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	if !enteredMarkAssistForTracing {
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		trace := traceAcquire()
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		if trace.ok() {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>			if !goexperiment.ExecTracer2 {
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>				<span class="comment">// In the old tracer, enter mark assist tracing only</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>				<span class="comment">// if we actually traced an event. Otherwise a goroutine</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>				<span class="comment">// waking up from mark assist post-GC might end up</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>				<span class="comment">// writing a stray &#34;end&#34; event.</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>				<span class="comment">// This means inMarkAssist will not be meaningful</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>				<span class="comment">// in the old tracer; that&#39;s OK, it&#39;s unused.</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>				<span class="comment">// See the comment on enteredMarkAssistForTracing.</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>				enteredMarkAssistForTracing = true
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			trace.GCMarkAssistStart()
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>			<span class="comment">// Set this *after* we trace the start, otherwise we may</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>			<span class="comment">// emit an in-progress event for an assist we&#39;re about to start.</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>			gp.inMarkAssist = true
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			traceRelease(trace)
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		} else {
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			gp.inMarkAssist = true
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		if goexperiment.ExecTracer2 {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>			<span class="comment">// In the new tracer, set enter mark assist tracing if we</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>			<span class="comment">// ever pass this point, because we must manage inMarkAssist</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>			<span class="comment">// correctly.</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>			<span class="comment">// See the comment on enteredMarkAssistForTracing.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>			enteredMarkAssistForTracing = true
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	}
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	<span class="comment">// Perform assist work</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		gcAssistAlloc1(gp, scanWork)
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		<span class="comment">// The user stack may have moved, so this can&#39;t touch</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		<span class="comment">// anything on it until it returns from systemstack.</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	})
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	completed := gp.param != nil
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	gp.param = nil
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	if completed {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		gcMarkDone()
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	}
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	if gp.gcAssistBytes &lt; 0 {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		<span class="comment">// We were unable steal enough credit or perform</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		<span class="comment">// enough work to pay off the assist debt. We need to</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		<span class="comment">// do one of these before letting the mutator allocate</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		<span class="comment">// more to prevent over-allocation.</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		<span class="comment">// If this is because we were preempted, reschedule</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		<span class="comment">// and try some more.</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		if gp.preempt {
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>			Gosched()
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>			goto retry
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		}
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		<span class="comment">// Add this G to an assist queue and park. When the GC</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		<span class="comment">// has more background credit, it will satisfy queued</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		<span class="comment">// assists before flushing to the global credit pool.</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		<span class="comment">// Note that this does *not* get woken up when more</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		<span class="comment">// work is added to the work list. The theory is that</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		<span class="comment">// there wasn&#39;t enough work to do anyway, so we might</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		<span class="comment">// as well let background marking take care of the</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		<span class="comment">// work that is available.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		if !gcParkAssist() {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>			goto retry
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		}
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		<span class="comment">// At this point either background GC has satisfied</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		<span class="comment">// this G&#39;s assist debt, or the GC cycle is over.</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	}
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	if enteredMarkAssistForTracing {
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		trace := traceAcquire()
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		if trace.ok() {
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>			trace.GCMarkAssistDone()
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>			<span class="comment">// Set this *after* we trace the end to make sure</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>			<span class="comment">// that we emit an in-progress event if this is</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>			<span class="comment">// the first event for the goroutine in the trace</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>			<span class="comment">// or trace generation. Also, do this between</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>			<span class="comment">// acquire/release because this is part of the</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			<span class="comment">// goroutine&#39;s trace state, and it must be atomic</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>			<span class="comment">// with respect to the tracer.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>			gp.inMarkAssist = false
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>			traceRelease(trace)
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		} else {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>			<span class="comment">// This state is tracked even if tracing isn&#39;t enabled.</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>			<span class="comment">// It&#39;s only used by the new tracer.</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>			<span class="comment">// See the comment on enteredMarkAssistForTracing.</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>			gp.inMarkAssist = false
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		}
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	}
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>}
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span><span class="comment">// gcAssistAlloc1 is the part of gcAssistAlloc that runs on the system</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span><span class="comment">// stack. This is a separate function to make it easier to see that</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span><span class="comment">// we&#39;re not capturing anything from the user stack, since the user</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span><span class="comment">// stack may move while we&#39;re in this function.</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span><span class="comment">// gcAssistAlloc1 indicates whether this assist completed the mark</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span><span class="comment">// phase by setting gp.param to non-nil. This can&#39;t be communicated on</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span><span class="comment">// the stack since it may move.</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>func gcAssistAlloc1(gp *g, scanWork int64) {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	<span class="comment">// Clear the flag indicating that this assist completed the</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	<span class="comment">// mark phase.</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	gp.param = nil
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	if atomic.Load(&amp;gcBlackenEnabled) == 0 {
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		<span class="comment">// The gcBlackenEnabled check in malloc races with the</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		<span class="comment">// store that clears it but an atomic check in every malloc</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		<span class="comment">// would be a performance hit.</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		<span class="comment">// Instead we recheck it here on the non-preemptible system</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		<span class="comment">// stack to determine if we should perform an assist.</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		<span class="comment">// GC is done, so ignore any remaining debt.</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		gp.gcAssistBytes = 0
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>		return
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	<span class="comment">// Track time spent in this assist. Since we&#39;re on the</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	<span class="comment">// system stack, this is non-preemptible, so we can</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	<span class="comment">// just measure start and end time.</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	<span class="comment">// Limiter event tracking might be disabled if we end up here</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	<span class="comment">// while on a mark worker.</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	startTime := nanotime()
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	trackLimiterEvent := gp.m.p.ptr().limiterEvent.start(limiterEventMarkAssist, startTime)
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	decnwait := atomic.Xadd(&amp;work.nwait, -1)
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	if decnwait == work.nproc {
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		println(&#34;runtime: work.nwait =&#34;, decnwait, &#34;work.nproc=&#34;, work.nproc)
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		throw(&#34;nwait &gt; work.nprocs&#34;)
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	}
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	<span class="comment">// gcDrainN requires the caller to be preemptible.</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	casGToWaiting(gp, _Grunning, waitReasonGCAssistMarking)
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	<span class="comment">// drain own cached work first in the hopes that it</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	<span class="comment">// will be more cache friendly.</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	gcw := &amp;getg().m.p.ptr().gcw
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	workDone := gcDrainN(gcw, scanWork)
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	casgstatus(gp, _Gwaiting, _Grunning)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	<span class="comment">// Record that we did this much scan work.</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	<span class="comment">// Back out the number of bytes of assist credit that</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	<span class="comment">// this scan work counts for. The &#34;1+&#34; is a poor man&#39;s</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	<span class="comment">// round-up, to ensure this adds credit even if</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	<span class="comment">// assistBytesPerWork is very low.</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	assistBytesPerWork := gcController.assistBytesPerWork.Load()
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	gp.gcAssistBytes += 1 + int64(assistBytesPerWork*float64(workDone))
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	<span class="comment">// If this is the last worker and we ran out of work,</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	<span class="comment">// signal a completion point.</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	incnwait := atomic.Xadd(&amp;work.nwait, +1)
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	if incnwait &gt; work.nproc {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		println(&#34;runtime: work.nwait=&#34;, incnwait,
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>			&#34;work.nproc=&#34;, work.nproc)
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		throw(&#34;work.nwait &gt; work.nproc&#34;)
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	}
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	if incnwait == work.nproc &amp;&amp; !gcMarkWorkAvailable(nil) {
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		<span class="comment">// This has reached a background completion point. Set</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		<span class="comment">// gp.param to a non-nil value to indicate this. It</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		<span class="comment">// doesn&#39;t matter what we set it to (it just has to be</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		<span class="comment">// a valid pointer).</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		gp.param = unsafe.Pointer(gp)
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	}
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	now := nanotime()
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	duration := now - startTime
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	pp := gp.m.p.ptr()
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	pp.gcAssistTime += duration
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	if trackLimiterEvent {
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>		pp.limiterEvent.stop(limiterEventMarkAssist, now)
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	}
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	if pp.gcAssistTime &gt; gcAssistTimeSlack {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		gcController.assistTime.Add(pp.gcAssistTime)
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>		gcCPULimiter.update(now)
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		pp.gcAssistTime = 0
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	}
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>}
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span><span class="comment">// gcWakeAllAssists wakes all currently blocked assists. This is used</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span><span class="comment">// at the end of a GC cycle. gcBlackenEnabled must be false to prevent</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span><span class="comment">// new assists from going to sleep after this point.</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>func gcWakeAllAssists() {
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	lock(&amp;work.assistQueue.lock)
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	list := work.assistQueue.q.popList()
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	injectglist(&amp;list)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	unlock(&amp;work.assistQueue.lock)
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>}
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span><span class="comment">// gcParkAssist puts the current goroutine on the assist queue and parks.</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span><span class="comment">// gcParkAssist reports whether the assist is now satisfied. If it</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span><span class="comment">// returns false, the caller must retry the assist.</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>func gcParkAssist() bool {
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	lock(&amp;work.assistQueue.lock)
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	<span class="comment">// If the GC cycle finished while we were getting the lock,</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	<span class="comment">// exit the assist. The cycle can&#39;t finish while we hold the</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	<span class="comment">// lock.</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>	if atomic.Load(&amp;gcBlackenEnabled) == 0 {
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>		unlock(&amp;work.assistQueue.lock)
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		return true
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	}
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	gp := getg()
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	oldList := work.assistQueue.q
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	work.assistQueue.q.pushBack(gp)
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	<span class="comment">// Recheck for background credit now that this G is in</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	<span class="comment">// the queue, but can still back out. This avoids a</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	<span class="comment">// race in case background marking has flushed more</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	<span class="comment">// credit since we checked above.</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	if gcController.bgScanCredit.Load() &gt; 0 {
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		work.assistQueue.q = oldList
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		if oldList.tail != 0 {
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>			oldList.tail.ptr().schedlink.set(nil)
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		}
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		unlock(&amp;work.assistQueue.lock)
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>		return false
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	}
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	<span class="comment">// Park.</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	goparkunlock(&amp;work.assistQueue.lock, waitReasonGCAssistWait, traceBlockGCMarkAssist, 2)
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	return true
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>}
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span><span class="comment">// gcFlushBgCredit flushes scanWork units of background scan work</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span><span class="comment">// credit. This first satisfies blocked assists on the</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span><span class="comment">// work.assistQueue and then flushes any remaining credit to</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span><span class="comment">// gcController.bgScanCredit.</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span><span class="comment">// Write barriers are disallowed because this is used by gcDrain after</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span><span class="comment">// it has ensured that all work is drained and this must preserve that</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span><span class="comment">// condition.</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>func gcFlushBgCredit(scanWork int64) {
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	if work.assistQueue.q.empty() {
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		<span class="comment">// Fast path; there are no blocked assists. There&#39;s a</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		<span class="comment">// small window here where an assist may add itself to</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>		<span class="comment">// the blocked queue and park. If that happens, we&#39;ll</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>		<span class="comment">// just get it on the next flush.</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		gcController.bgScanCredit.Add(scanWork)
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		return
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	}
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	assistBytesPerWork := gcController.assistBytesPerWork.Load()
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	scanBytes := int64(float64(scanWork) * assistBytesPerWork)
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	lock(&amp;work.assistQueue.lock)
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>	for !work.assistQueue.q.empty() &amp;&amp; scanBytes &gt; 0 {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		gp := work.assistQueue.q.pop()
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		<span class="comment">// Note that gp.gcAssistBytes is negative because gp</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		<span class="comment">// is in debt. Think carefully about the signs below.</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>		if scanBytes+gp.gcAssistBytes &gt;= 0 {
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>			<span class="comment">// Satisfy this entire assist debt.</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>			scanBytes += gp.gcAssistBytes
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>			gp.gcAssistBytes = 0
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>			<span class="comment">// It&#39;s important that we *not* put gp in</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>			<span class="comment">// runnext. Otherwise, it&#39;s possible for user</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>			<span class="comment">// code to exploit the GC worker&#39;s high</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>			<span class="comment">// scheduler priority to get itself always run</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>			<span class="comment">// before other goroutines and always in the</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>			<span class="comment">// fresh quantum started by GC.</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>			ready(gp, 0, false)
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		} else {
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>			<span class="comment">// Partially satisfy this assist.</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>			gp.gcAssistBytes += scanBytes
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>			scanBytes = 0
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>			<span class="comment">// As a heuristic, we move this assist to the</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>			<span class="comment">// back of the queue so that large assists</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>			<span class="comment">// can&#39;t clog up the assist queue and</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>			<span class="comment">// substantially delay small assists.</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			work.assistQueue.q.pushBack(gp)
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>			break
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		}
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	}
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	if scanBytes &gt; 0 {
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		<span class="comment">// Convert from scan bytes back to work.</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>		assistWorkPerByte := gcController.assistWorkPerByte.Load()
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		scanWork = int64(float64(scanBytes) * assistWorkPerByte)
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>		gcController.bgScanCredit.Add(scanWork)
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	}
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	unlock(&amp;work.assistQueue.lock)
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>}
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span><span class="comment">// scanstack scans gp&#39;s stack, greying all pointers found on the stack.</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span><span class="comment">// Returns the amount of scan work performed, but doesn&#39;t update</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span><span class="comment">// gcController.stackScanWork or flush any credit. Any background credit produced</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span><span class="comment">// by this function should be flushed by its caller. scanstack itself can&#39;t</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span><span class="comment">// safely flush because it may result in trying to wake up a goroutine that</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span><span class="comment">// was just scanned, resulting in a self-deadlock.</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span><span class="comment">// scanstack will also shrink the stack if it is safe to do so. If it</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span><span class="comment">// is not, it schedules a stack shrink for the next synchronous safe</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span><span class="comment">// point.</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span><span class="comment">// scanstack is marked go:systemstack because it must not be preempted</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span><span class="comment">// while using a workbuf.</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>func scanstack(gp *g, gcw *gcWork) int64 {
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	if readgstatus(gp)&amp;_Gscan == 0 {
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>		print(&#34;runtime:scanstack: gp=&#34;, gp, &#34;, goid=&#34;, gp.goid, &#34;, gp-&gt;atomicstatus=&#34;, hex(readgstatus(gp)), &#34;\n&#34;)
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		throw(&#34;scanstack - bad status&#34;)
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	}
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	switch readgstatus(gp) &amp;^ _Gscan {
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	default:
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		print(&#34;runtime: gp=&#34;, gp, &#34;, goid=&#34;, gp.goid, &#34;, gp-&gt;atomicstatus=&#34;, readgstatus(gp), &#34;\n&#34;)
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		throw(&#34;mark - bad status&#34;)
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	case _Gdead:
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		return 0
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	case _Grunning:
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>		print(&#34;runtime: gp=&#34;, gp, &#34;, goid=&#34;, gp.goid, &#34;, gp-&gt;atomicstatus=&#34;, readgstatus(gp), &#34;\n&#34;)
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>		throw(&#34;scanstack: goroutine not stopped&#34;)
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>	case _Grunnable, _Gsyscall, _Gwaiting:
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		<span class="comment">// ok</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	}
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	if gp == getg() {
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>		throw(&#34;can&#39;t scan our own stack&#34;)
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	}
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	<span class="comment">// scannedSize is the amount of work we&#39;ll be reporting.</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	<span class="comment">// It is less than the allocated size (which is hi-lo).</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	var sp uintptr
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	if gp.syscallsp != 0 {
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		sp = gp.syscallsp <span class="comment">// If in a system call this is the stack pointer (gp.sched.sp can be 0 in this case on Windows).</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	} else {
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		sp = gp.sched.sp
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	}
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	scannedSize := gp.stack.hi - sp
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>	<span class="comment">// Keep statistics for initial stack size calculation.</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	<span class="comment">// Note that this accumulates the scanned size, not the allocated size.</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	p := getg().m.p.ptr()
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	p.scannedStackSize += uint64(scannedSize)
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	p.scannedStacks++
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	if isShrinkStackSafe(gp) {
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		<span class="comment">// Shrink the stack if not much of it is being used.</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		shrinkstack(gp)
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	} else {
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>		<span class="comment">// Otherwise, shrink the stack at the next sync safe point.</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>		gp.preemptShrink = true
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	}
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	var state stackScanState
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	state.stack = gp.stack
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	if stackTraceDebug {
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>		println(&#34;stack trace goroutine&#34;, gp.goid)
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	}
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	if debugScanConservative &amp;&amp; gp.asyncSafePoint {
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>		print(&#34;scanning async preempted goroutine &#34;, gp.goid, &#34; stack [&#34;, hex(gp.stack.lo), &#34;,&#34;, hex(gp.stack.hi), &#34;)\n&#34;)
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>	}
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>	<span class="comment">// Scan the saved context register. This is effectively a live</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>	<span class="comment">// register that gets moved back and forth between the</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	<span class="comment">// register and sched.ctxt without a write barrier.</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	if gp.sched.ctxt != nil {
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>		scanblock(uintptr(unsafe.Pointer(&amp;gp.sched.ctxt)), goarch.PtrSize, &amp;oneptrmask[0], gcw, &amp;state)
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	}
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	<span class="comment">// Scan the stack. Accumulate a list of stack objects.</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	var u unwinder
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	for u.init(gp, 0); u.valid(); u.next() {
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		scanframeworker(&amp;u.frame, &amp;state, gcw)
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	}
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>	<span class="comment">// Find additional pointers that point into the stack from the heap.</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	<span class="comment">// Currently this includes defers and panics. See also function copystack.</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	<span class="comment">// Find and trace other pointers in defer records.</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	for d := gp._defer; d != nil; d = d.link {
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>		if d.fn != nil {
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>			<span class="comment">// Scan the func value, which could be a stack allocated closure.</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>			<span class="comment">// See issue 30453.</span>
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>			scanblock(uintptr(unsafe.Pointer(&amp;d.fn)), goarch.PtrSize, &amp;oneptrmask[0], gcw, &amp;state)
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>		}
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>		if d.link != nil {
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>			<span class="comment">// The link field of a stack-allocated defer record might point</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>			<span class="comment">// to a heap-allocated defer record. Keep that heap record live.</span>
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>			scanblock(uintptr(unsafe.Pointer(&amp;d.link)), goarch.PtrSize, &amp;oneptrmask[0], gcw, &amp;state)
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		}
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		<span class="comment">// Retain defers records themselves.</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>		<span class="comment">// Defer records might not be reachable from the G through regular heap</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>		<span class="comment">// tracing because the defer linked list might weave between the stack and the heap.</span>
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		if d.heap {
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>			scanblock(uintptr(unsafe.Pointer(&amp;d)), goarch.PtrSize, &amp;oneptrmask[0], gcw, &amp;state)
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>		}
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	}
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	if gp._panic != nil {
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>		<span class="comment">// Panics are always stack allocated.</span>
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>		state.putPtr(uintptr(unsafe.Pointer(gp._panic)), false)
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	}
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>	<span class="comment">// Find and scan all reachable stack objects.</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>	<span class="comment">// The state&#39;s pointer queue prioritizes precise pointers over</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>	<span class="comment">// conservative pointers so that we&#39;ll prefer scanning stack</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	<span class="comment">// objects precisely.</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>	state.buildIndex()
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>	for {
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>		p, conservative := state.getPtr()
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>		if p == 0 {
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>			break
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>		}
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		obj := state.findObject(p)
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		if obj == nil {
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>			continue
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>		}
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>		r := obj.r
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>		if r == nil {
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>			<span class="comment">// We&#39;ve already scanned this object.</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>			continue
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		}
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>		obj.setRecord(nil) <span class="comment">// Don&#39;t scan it again.</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>		if stackTraceDebug {
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			printlock()
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>			print(&#34;  live stkobj at&#34;, hex(state.stack.lo+uintptr(obj.off)), &#34;of size&#34;, obj.size)
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>			if conservative {
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>				print(&#34; (conservative)&#34;)
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>			}
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>			println()
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>			printunlock()
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		}
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>		gcdata := r.gcdata()
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>		var s *mspan
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>		if r.useGCProg() {
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>			<span class="comment">// This path is pretty unlikely, an object large enough</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>			<span class="comment">// to have a GC program allocated on the stack.</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>			<span class="comment">// We need some space to unpack the program into a straight</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>			<span class="comment">// bitmask, which we allocate/free here.</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>			<span class="comment">// TODO: it would be nice if there were a way to run a GC</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>			<span class="comment">// program without having to store all its bits. We&#39;d have</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>			<span class="comment">// to change from a Lempel-Ziv style program to something else.</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>			<span class="comment">// Or we can forbid putting objects on stacks if they require</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>			<span class="comment">// a gc program (see issue 27447).</span>
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>			s = materializeGCProg(r.ptrdata(), gcdata)
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>			gcdata = (*byte)(unsafe.Pointer(s.startAddr))
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>		}
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>		b := state.stack.lo + uintptr(obj.off)
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>		if conservative {
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>			scanConservative(b, r.ptrdata(), gcdata, gcw, &amp;state)
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>		} else {
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>			scanblock(b, r.ptrdata(), gcdata, gcw, &amp;state)
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>		}
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		if s != nil {
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>			dematerializeGCProg(s)
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>		}
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>	}
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>	<span class="comment">// Deallocate object buffers.</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	<span class="comment">// (Pointer buffers were all deallocated in the loop above.)</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	for state.head != nil {
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>		x := state.head
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>		state.head = x.next
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>		if stackTraceDebug {
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>			for i := 0; i &lt; x.nobj; i++ {
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>				obj := &amp;x.obj[i]
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>				if obj.r == nil { <span class="comment">// reachable</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>					continue
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>				}
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>				println(&#34;  dead stkobj at&#34;, hex(gp.stack.lo+uintptr(obj.off)), &#34;of size&#34;, obj.r.size)
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>				<span class="comment">// Note: not necessarily really dead - only reachable-from-ptr dead.</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>			}
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>		}
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>		x.nobj = 0
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>		putempty((*workbuf)(unsafe.Pointer(x)))
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	}
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	if state.buf != nil || state.cbuf != nil || state.freeBuf != nil {
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>		throw(&#34;remaining pointer buffers&#34;)
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	}
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>	return int64(scannedSize)
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>}
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span><span class="comment">// Scan a stack frame: local variables and function arguments/results.</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>func scanframeworker(frame *stkframe, state *stackScanState, gcw *gcWork) {
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	if _DebugGC &gt; 1 &amp;&amp; frame.continpc != 0 {
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>		print(&#34;scanframe &#34;, funcname(frame.fn), &#34;\n&#34;)
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	}
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>	isAsyncPreempt := frame.fn.valid() &amp;&amp; frame.fn.funcID == abi.FuncID_asyncPreempt
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	isDebugCall := frame.fn.valid() &amp;&amp; frame.fn.funcID == abi.FuncID_debugCallV2
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	if state.conservative || isAsyncPreempt || isDebugCall {
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		if debugScanConservative {
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>			println(&#34;conservatively scanning function&#34;, funcname(frame.fn), &#34;at PC&#34;, hex(frame.continpc))
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		}
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>		<span class="comment">// Conservatively scan the frame. Unlike the precise</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>		<span class="comment">// case, this includes the outgoing argument space</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>		<span class="comment">// since we may have stopped while this function was</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>		<span class="comment">// setting up a call.</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>		<span class="comment">// TODO: We could narrow this down if the compiler</span>
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		<span class="comment">// produced a single map per function of stack slots</span>
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>		<span class="comment">// and registers that ever contain a pointer.</span>
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>		if frame.varp != 0 {
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>			size := frame.varp - frame.sp
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>			if size &gt; 0 {
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>				scanConservative(frame.sp, size, nil, gcw, state)
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>			}
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>		}
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		<span class="comment">// Scan arguments to this frame.</span>
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>		if n := frame.argBytes(); n != 0 {
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>			<span class="comment">// TODO: We could pass the entry argument map</span>
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>			<span class="comment">// to narrow this down further.</span>
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>			scanConservative(frame.argp, n, nil, gcw, state)
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>		}
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>		if isAsyncPreempt || isDebugCall {
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>			<span class="comment">// This function&#39;s frame contained the</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>			<span class="comment">// registers for the asynchronously stopped</span>
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>			<span class="comment">// parent frame. Scan the parent</span>
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>			<span class="comment">// conservatively.</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>			state.conservative = true
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>		} else {
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>			<span class="comment">// We only wanted to scan those two frames</span>
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>			<span class="comment">// conservatively. Clear the flag for future</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>			<span class="comment">// frames.</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>			state.conservative = false
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>		}
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>		return
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>	}
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>	locals, args, objs := frame.getStackMap(false)
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	<span class="comment">// Scan local variables if stack frame has been allocated.</span>
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>	if locals.n &gt; 0 {
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>		size := uintptr(locals.n) * goarch.PtrSize
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>		scanblock(frame.varp-size, size, locals.bytedata, gcw, state)
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	}
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>	<span class="comment">// Scan arguments.</span>
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	if args.n &gt; 0 {
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>		scanblock(frame.argp, uintptr(args.n)*goarch.PtrSize, args.bytedata, gcw, state)
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	}
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	<span class="comment">// Add all stack objects to the stack object list.</span>
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>	if frame.varp != 0 {
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>		<span class="comment">// varp is 0 for defers, where there are no locals.</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>		<span class="comment">// In that case, there can&#39;t be a pointer to its args, either.</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>		<span class="comment">// (And all args would be scanned above anyway.)</span>
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>		for i := range objs {
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>			obj := &amp;objs[i]
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>			off := obj.off
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>			base := frame.varp <span class="comment">// locals base pointer</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>			if off &gt;= 0 {
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>				base = frame.argp <span class="comment">// arguments and return values base pointer</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>			}
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>			ptr := base + uintptr(off)
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>			if ptr &lt; frame.sp {
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>				<span class="comment">// object hasn&#39;t been allocated in the frame yet.</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>				continue
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>			}
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>			if stackTraceDebug {
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>				println(&#34;stkobj at&#34;, hex(ptr), &#34;of size&#34;, obj.size)
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>			}
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>			state.addObject(ptr, obj)
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>		}
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	}
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>}
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>type gcDrainFlags int
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>const (
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>	gcDrainUntilPreempt gcDrainFlags = 1 &lt;&lt; iota
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>	gcDrainFlushBgCredit
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>	gcDrainIdle
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>	gcDrainFractional
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>)
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span><span class="comment">// gcDrainMarkWorkerIdle is a wrapper for gcDrain that exists to better account</span>
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span><span class="comment">// mark time in profiles.</span>
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>func gcDrainMarkWorkerIdle(gcw *gcWork) {
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>	gcDrain(gcw, gcDrainIdle|gcDrainUntilPreempt|gcDrainFlushBgCredit)
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>}
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span><span class="comment">// gcDrainMarkWorkerDedicated is a wrapper for gcDrain that exists to better account</span>
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span><span class="comment">// mark time in profiles.</span>
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>func gcDrainMarkWorkerDedicated(gcw *gcWork, untilPreempt bool) {
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>	flags := gcDrainFlushBgCredit
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>	if untilPreempt {
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>		flags |= gcDrainUntilPreempt
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	}
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>	gcDrain(gcw, flags)
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>}
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span><span class="comment">// gcDrainMarkWorkerFractional is a wrapper for gcDrain that exists to better account</span>
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span><span class="comment">// mark time in profiles.</span>
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>func gcDrainMarkWorkerFractional(gcw *gcWork) {
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>	gcDrain(gcw, gcDrainFractional|gcDrainUntilPreempt|gcDrainFlushBgCredit)
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>}
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span><span class="comment">// gcDrain scans roots and objects in work buffers, blackening grey</span>
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span><span class="comment">// objects until it is unable to get more work. It may return before</span>
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span><span class="comment">// GC is done; it&#39;s the caller&#39;s responsibility to balance work from</span>
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span><span class="comment">// other Ps.</span>
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span><span class="comment">// If flags&amp;gcDrainUntilPreempt != 0, gcDrain returns when g.preempt</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span><span class="comment">// is set.</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span><span class="comment">// If flags&amp;gcDrainIdle != 0, gcDrain returns when there is other work</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span><span class="comment">// to do.</span>
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span><span class="comment">// If flags&amp;gcDrainFractional != 0, gcDrain self-preempts when</span>
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span><span class="comment">// pollFractionalWorkerExit() returns true. This implies</span>
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span><span class="comment">// gcDrainNoBlock.</span>
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span><span class="comment">// If flags&amp;gcDrainFlushBgCredit != 0, gcDrain flushes scan work</span>
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span><span class="comment">// credit to gcController.bgScanCredit every gcCreditSlack units of</span>
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span><span class="comment">// scan work.</span>
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span><span class="comment">// gcDrain will always return if there is a pending STW or forEachP.</span>
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span><span class="comment">// Disabling write barriers is necessary to ensure that after we&#39;ve</span>
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span><span class="comment">// confirmed that we&#39;ve drained gcw, that we don&#39;t accidentally end</span>
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span><span class="comment">// up flipping that condition by immediately adding work in the form</span>
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span><span class="comment">// of a write barrier buffer flush.</span>
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span><span class="comment">// Don&#39;t set nowritebarrierrec because it&#39;s safe for some callees to</span>
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span><span class="comment">// have write barriers enabled.</span>
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>func gcDrain(gcw *gcWork, flags gcDrainFlags) {
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>	if !writeBarrier.enabled {
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>		throw(&#34;gcDrain phase incorrect&#34;)
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>	}
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>	<span class="comment">// N.B. We must be running in a non-preemptible context, so it&#39;s</span>
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>	<span class="comment">// safe to hold a reference to our P here.</span>
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>	gp := getg().m.curg
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>	pp := gp.m.p.ptr()
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>	preemptible := flags&amp;gcDrainUntilPreempt != 0
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>	flushBgCredit := flags&amp;gcDrainFlushBgCredit != 0
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>	idle := flags&amp;gcDrainIdle != 0
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>	initScanWork := gcw.heapScanWork
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>	<span class="comment">// checkWork is the scan work before performing the next</span>
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>	<span class="comment">// self-preempt check.</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>	checkWork := int64(1&lt;&lt;63 - 1)
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>	var check func() bool
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>	if flags&amp;(gcDrainIdle|gcDrainFractional) != 0 {
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>		checkWork = initScanWork + drainCheckThreshold
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>		if idle {
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>			check = pollWork
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>		} else if flags&amp;gcDrainFractional != 0 {
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>			check = pollFractionalWorkerExit
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>		}
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>	}
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>	<span class="comment">// Drain root marking jobs.</span>
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>	if work.markrootNext &lt; work.markrootJobs {
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>		<span class="comment">// Stop if we&#39;re preemptible, if someone wants to STW, or if</span>
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>		<span class="comment">// someone is calling forEachP.</span>
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>		for !(gp.preempt &amp;&amp; (preemptible || sched.gcwaiting.Load() || pp.runSafePointFn != 0)) {
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>			job := atomic.Xadd(&amp;work.markrootNext, +1) - 1
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>			if job &gt;= work.markrootJobs {
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>				break
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>			}
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>			markroot(gcw, job, flushBgCredit)
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>			if check != nil &amp;&amp; check() {
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>				goto done
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>			}
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>		}
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>	}
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>	<span class="comment">// Drain heap marking jobs.</span>
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>	<span class="comment">// Stop if we&#39;re preemptible, if someone wants to STW, or if</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>	<span class="comment">// someone is calling forEachP.</span>
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Consider always checking gp.preempt instead</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>	<span class="comment">// of having the preempt flag, and making an exception for certain</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>	<span class="comment">// mark workers in retake. That might be simpler than trying to</span>
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>	<span class="comment">// enumerate all the reasons why we might want to preempt, even</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>	<span class="comment">// if we&#39;re supposed to be mostly non-preemptible.</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>	for !(gp.preempt &amp;&amp; (preemptible || sched.gcwaiting.Load() || pp.runSafePointFn != 0)) {
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>		<span class="comment">// Try to keep work available on the global queue. We used to</span>
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>		<span class="comment">// check if there were waiting workers, but it&#39;s better to</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>		<span class="comment">// just keep work available than to make workers wait. In the</span>
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>		<span class="comment">// worst case, we&#39;ll do O(log(_WorkbufSize)) unnecessary</span>
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>		<span class="comment">// balances.</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>		if work.full == 0 {
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>			gcw.balance()
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>		}
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>		b := gcw.tryGetFast()
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>		if b == 0 {
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>			b = gcw.tryGet()
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>			if b == 0 {
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>				<span class="comment">// Flush the write barrier</span>
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>				<span class="comment">// buffer; this may create</span>
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>				<span class="comment">// more work.</span>
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>				wbBufFlush()
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>				b = gcw.tryGet()
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>			}
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>		}
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>		if b == 0 {
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>			<span class="comment">// Unable to get work.</span>
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>			break
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>		}
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>		scanobject(b, gcw)
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>		<span class="comment">// Flush background scan work credit to the global</span>
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>		<span class="comment">// account if we&#39;ve accumulated enough locally so</span>
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>		<span class="comment">// mutator assists can draw on it.</span>
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>		if gcw.heapScanWork &gt;= gcCreditSlack {
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>			gcController.heapScanWork.Add(gcw.heapScanWork)
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>			if flushBgCredit {
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>				gcFlushBgCredit(gcw.heapScanWork - initScanWork)
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>				initScanWork = 0
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>			}
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>			checkWork -= gcw.heapScanWork
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>			gcw.heapScanWork = 0
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>			if checkWork &lt;= 0 {
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>				checkWork += drainCheckThreshold
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>				if check != nil &amp;&amp; check() {
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>					break
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>				}
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>			}
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>		}
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>	}
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>done:
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>	<span class="comment">// Flush remaining scan work credit.</span>
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>	if gcw.heapScanWork &gt; 0 {
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>		gcController.heapScanWork.Add(gcw.heapScanWork)
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>		if flushBgCredit {
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>			gcFlushBgCredit(gcw.heapScanWork - initScanWork)
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>		}
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>		gcw.heapScanWork = 0
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	}
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>}
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span><span class="comment">// gcDrainN blackens grey objects until it has performed roughly</span>
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span><span class="comment">// scanWork units of scan work or the G is preempted. This is</span>
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span><span class="comment">// best-effort, so it may perform less work if it fails to get a work</span>
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span><span class="comment">// buffer. Otherwise, it will perform at least n units of work, but</span>
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span><span class="comment">// may perform more because scanning is always done in whole object</span>
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span><span class="comment">// increments. It returns the amount of scan work performed.</span>
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span><span class="comment">// The caller goroutine must be in a preemptible state (e.g.,</span>
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span><span class="comment">// _Gwaiting) to prevent deadlocks during stack scanning. As a</span>
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span><span class="comment">// consequence, this must be called on the system stack.</span>
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>func gcDrainN(gcw *gcWork, scanWork int64) int64 {
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>	if !writeBarrier.enabled {
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>		throw(&#34;gcDrainN phase incorrect&#34;)
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>	}
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>	<span class="comment">// There may already be scan work on the gcw, which we don&#39;t</span>
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>	<span class="comment">// want to claim was done by this call.</span>
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>	workFlushed := -gcw.heapScanWork
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>	<span class="comment">// In addition to backing out because of a preemption, back out</span>
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>	<span class="comment">// if the GC CPU limiter is enabled.</span>
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>	gp := getg().m.curg
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>	for !gp.preempt &amp;&amp; !gcCPULimiter.limiting() &amp;&amp; workFlushed+gcw.heapScanWork &lt; scanWork {
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>		<span class="comment">// See gcDrain comment.</span>
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>		if work.full == 0 {
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>			gcw.balance()
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>		}
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>		b := gcw.tryGetFast()
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>		if b == 0 {
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>			b = gcw.tryGet()
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>			if b == 0 {
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>				<span class="comment">// Flush the write barrier buffer;</span>
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>				<span class="comment">// this may create more work.</span>
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>				wbBufFlush()
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>				b = gcw.tryGet()
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>			}
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>		}
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>		if b == 0 {
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>			<span class="comment">// Try to do a root job.</span>
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>			if work.markrootNext &lt; work.markrootJobs {
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>				job := atomic.Xadd(&amp;work.markrootNext, +1) - 1
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>				if job &lt; work.markrootJobs {
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>					workFlushed += markroot(gcw, job, false)
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>					continue
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>				}
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>			}
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>			<span class="comment">// No heap or root jobs.</span>
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>			break
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>		}
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>		scanobject(b, gcw)
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>		<span class="comment">// Flush background scan work credit.</span>
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>		if gcw.heapScanWork &gt;= gcCreditSlack {
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>			gcController.heapScanWork.Add(gcw.heapScanWork)
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>			workFlushed += gcw.heapScanWork
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>			gcw.heapScanWork = 0
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>		}
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>	}
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>	<span class="comment">// Unlike gcDrain, there&#39;s no need to flush remaining work</span>
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>	<span class="comment">// here because this never flushes to bgScanCredit and</span>
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>	<span class="comment">// gcw.dispose will flush any remaining work to scanWork.</span>
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>	return workFlushed + gcw.heapScanWork
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>}
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span><span class="comment">// scanblock scans b as scanobject would, but using an explicit</span>
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span><span class="comment">// pointer bitmap instead of the heap bitmap.</span>
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span><span class="comment">// This is used to scan non-heap roots, so it does not update</span>
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span><span class="comment">// gcw.bytesMarked or gcw.heapScanWork.</span>
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span><span class="comment">// If stk != nil, possible stack pointers are also reported to stk.putPtr.</span>
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>func scanblock(b0, n0 uintptr, ptrmask *uint8, gcw *gcWork, stk *stackScanState) {
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>	<span class="comment">// Use local copies of original parameters, so that a stack trace</span>
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>	<span class="comment">// due to one of the throws below shows the original block</span>
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>	<span class="comment">// base and extent.</span>
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>	b := b0
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>	n := n0
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; n; {
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>		<span class="comment">// Find bits for the next word.</span>
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>		bits := uint32(*addb(ptrmask, i/(goarch.PtrSize*8)))
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>		if bits == 0 {
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>			i += goarch.PtrSize * 8
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>			continue
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>		}
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>		for j := 0; j &lt; 8 &amp;&amp; i &lt; n; j++ {
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>			if bits&amp;1 != 0 {
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>				<span class="comment">// Same work as in scanobject; see comments there.</span>
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>				p := *(*uintptr)(unsafe.Pointer(b + i))
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>				if p != 0 {
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>					if obj, span, objIndex := findObject(p, b, i); obj != 0 {
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>						greyobject(obj, b, i, span, gcw, objIndex)
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>					} else if stk != nil &amp;&amp; p &gt;= stk.stack.lo &amp;&amp; p &lt; stk.stack.hi {
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>						stk.putPtr(p, false)
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>					}
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>				}
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>			}
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>			bits &gt;&gt;= 1
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>			i += goarch.PtrSize
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>		}
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>	}
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>}
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span><span class="comment">// scanobject scans the object starting at b, adding pointers to gcw.</span>
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span><span class="comment">// b must point to the beginning of a heap object or an oblet.</span>
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span><span class="comment">// scanobject consults the GC bitmap for the pointer mask and the</span>
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span><span class="comment">// spans for the size of the object.</span>
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>func scanobject(b uintptr, gcw *gcWork) {
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>	<span class="comment">// Prefetch object before we scan it.</span>
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>	<span class="comment">// This will overlap fetching the beginning of the object with initial</span>
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>	<span class="comment">// setup before we start scanning the object.</span>
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>	sys.Prefetch(b)
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>	<span class="comment">// Find the bits for b and the size of the object at b.</span>
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>	<span class="comment">// b is either the beginning of an object, in which case this</span>
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>	<span class="comment">// is the size of the object to scan, or it points to an</span>
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>	<span class="comment">// oblet, in which case we compute the size to scan below.</span>
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>	s := spanOfUnchecked(b)
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>	n := s.elemsize
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>	if n == 0 {
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>		throw(&#34;scanobject n == 0&#34;)
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>	}
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>	if s.spanclass.noscan() {
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>		<span class="comment">// Correctness-wise this is ok, but it&#39;s inefficient</span>
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>		<span class="comment">// if noscan objects reach here.</span>
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>		throw(&#34;scanobject of a noscan object&#34;)
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>	}
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>	var tp typePointers
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>	if n &gt; maxObletBytes {
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>		<span class="comment">// Large object. Break into oblets for better</span>
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>		<span class="comment">// parallelism and lower latency.</span>
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>		if b == s.base() {
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>			<span class="comment">// Enqueue the other oblets to scan later.</span>
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>			<span class="comment">// Some oblets may be in b&#39;s scalar tail, but</span>
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>			<span class="comment">// these will be marked as &#34;no more pointers&#34;,</span>
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>			<span class="comment">// so we&#39;ll drop out immediately when we go to</span>
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>			<span class="comment">// scan those.</span>
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>			for oblet := b + maxObletBytes; oblet &lt; s.base()+s.elemsize; oblet += maxObletBytes {
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>				if !gcw.putFast(oblet) {
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span>					gcw.put(oblet)
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span>				}
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>			}
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>		}
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span>		<span class="comment">// Compute the size of the oblet. Since this object</span>
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span>		<span class="comment">// must be a large object, s.base() is the beginning</span>
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span>		<span class="comment">// of the object.</span>
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span>		n = s.base() + s.elemsize - b
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span>		n = min(n, maxObletBytes)
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span>		if goexperiment.AllocHeaders {
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>			tp = s.typePointersOfUnchecked(s.base())
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>			tp = tp.fastForward(b-tp.addr, b+n)
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>		}
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>	} else {
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>		if goexperiment.AllocHeaders {
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>			tp = s.typePointersOfUnchecked(b)
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>		}
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>	}
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span>	var hbits heapBits
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span>	if !goexperiment.AllocHeaders {
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span>		hbits = heapBitsForAddr(b, n)
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span>	}
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span>	var scanSize uintptr
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span>	for {
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span>		var addr uintptr
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span>		if goexperiment.AllocHeaders {
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span>			if tp, addr = tp.nextFast(); addr == 0 {
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span>				if tp, addr = tp.next(b + n); addr == 0 {
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span>					break
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span>				}
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span>			}
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span>		} else {
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span>			if hbits, addr = hbits.nextFast(); addr == 0 {
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span>				if hbits, addr = hbits.next(); addr == 0 {
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span>					break
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span>				}
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span>			}
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span>		}
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span>
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span>		<span class="comment">// Keep track of farthest pointer we found, so we can</span>
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span>		<span class="comment">// update heapScanWork. TODO: is there a better metric,</span>
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span>		<span class="comment">// now that we can skip scalar portions pretty efficiently?</span>
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span>		scanSize = addr - b + goarch.PtrSize
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span>
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span>		<span class="comment">// Work here is duplicated in scanblock and above.</span>
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span>		<span class="comment">// If you make changes here, make changes there too.</span>
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span>		obj := *(*uintptr)(unsafe.Pointer(addr))
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span>
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span>		<span class="comment">// At this point we have extracted the next potential pointer.</span>
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span>		<span class="comment">// Quickly filter out nil and pointers back to the current object.</span>
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span>		if obj != 0 &amp;&amp; obj-b &gt;= n {
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span>			<span class="comment">// Test if obj points into the Go heap and, if so,</span>
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span>			<span class="comment">// mark the object.</span>
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span>			<span class="comment">// Note that it&#39;s possible for findObject to</span>
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span>			<span class="comment">// fail if obj points to a just-allocated heap</span>
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span>			<span class="comment">// object because of a race with growing the</span>
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span>			<span class="comment">// heap. In this case, we know the object was</span>
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span>			<span class="comment">// just allocated and hence will be marked by</span>
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span>			<span class="comment">// allocation itself.</span>
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span>			if obj, span, objIndex := findObject(obj, b, addr-b); obj != 0 {
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span>				greyobject(obj, b, addr-b, span, gcw, objIndex)
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span>			}
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span>		}
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span>	}
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span>	gcw.bytesMarked += uint64(n)
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span>	gcw.heapScanWork += int64(scanSize)
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span>}
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span>
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span><span class="comment">// scanConservative scans block [b, b+n) conservatively, treating any</span>
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span><span class="comment">// pointer-like value in the block as a pointer.</span>
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span><span class="comment">// If ptrmask != nil, only words that are marked in ptrmask are</span>
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span><span class="comment">// considered as potential pointers.</span>
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span><span class="comment">// If state != nil, it&#39;s assumed that [b, b+n) is a block in the stack</span>
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span><span class="comment">// and may contain pointers to stack objects.</span>
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span>func scanConservative(b, n uintptr, ptrmask *uint8, gcw *gcWork, state *stackScanState) {
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span>	if debugScanConservative {
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span>		printlock()
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span>		print(&#34;conservatively scanning [&#34;, hex(b), &#34;,&#34;, hex(b+n), &#34;)\n&#34;)
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span>		hexdumpWords(b, b+n, func(p uintptr) byte {
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span>			if ptrmask != nil {
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>				word := (p - b) / goarch.PtrSize
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>				bits := *addb(ptrmask, word/8)
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span>				if (bits&gt;&gt;(word%8))&amp;1 == 0 {
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span>					return &#39;$&#39;
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span>				}
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span>			}
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span>
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span>			val := *(*uintptr)(unsafe.Pointer(p))
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span>			if state != nil &amp;&amp; state.stack.lo &lt;= val &amp;&amp; val &lt; state.stack.hi {
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span>				return &#39;@&#39;
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span>			}
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span>
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span>			span := spanOfHeap(val)
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span>			if span == nil {
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span>				return &#39; &#39;
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>			}
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>			idx := span.objIndex(val)
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>			if span.isFree(idx) {
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>				return &#39; &#39;
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span>			}
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span>			return &#39;*&#39;
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span>		})
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span>		printunlock()
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>	}
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; n; i += goarch.PtrSize {
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span>		if ptrmask != nil {
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span>			word := i / goarch.PtrSize
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span>			bits := *addb(ptrmask, word/8)
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span>			if bits == 0 {
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>				<span class="comment">// Skip 8 words (the loop increment will do the 8th)</span>
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>				<span class="comment">// This must be the first time we&#39;ve</span>
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>				<span class="comment">// seen this word of ptrmask, so i</span>
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span>				<span class="comment">// must be 8-word-aligned, but check</span>
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span>				<span class="comment">// our reasoning just in case.</span>
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span>				if i%(goarch.PtrSize*8) != 0 {
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span>					throw(&#34;misaligned mask&#34;)
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span>				}
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span>				i += goarch.PtrSize*8 - goarch.PtrSize
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span>				continue
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span>			}
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span>			if (bits&gt;&gt;(word%8))&amp;1 == 0 {
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span>				continue
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span>			}
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span>		}
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span>
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span>		val := *(*uintptr)(unsafe.Pointer(b + i))
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span>
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span>		<span class="comment">// Check if val points into the stack.</span>
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span>		if state != nil &amp;&amp; state.stack.lo &lt;= val &amp;&amp; val &lt; state.stack.hi {
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span>			<span class="comment">// val may point to a stack object. This</span>
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span>			<span class="comment">// object may be dead from last cycle and</span>
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span>			<span class="comment">// hence may contain pointers to unallocated</span>
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span>			<span class="comment">// objects, but unlike heap objects we can&#39;t</span>
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span>			<span class="comment">// tell if it&#39;s already dead. Hence, if all</span>
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span>			<span class="comment">// pointers to this object are from</span>
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span>			<span class="comment">// conservative scanning, we have to scan it</span>
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span>			<span class="comment">// defensively, too.</span>
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span>			state.putPtr(val, true)
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span>			continue
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span>		}
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span>
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span>		<span class="comment">// Check if val points to a heap span.</span>
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span>		span := spanOfHeap(val)
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span>		if span == nil {
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span>			continue
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span>		}
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span>		<span class="comment">// Check if val points to an allocated object.</span>
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span>		idx := span.objIndex(val)
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span>		if span.isFree(idx) {
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span>			continue
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span>		}
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>		<span class="comment">// val points to an allocated object. Mark it.</span>
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>		obj := span.base() + idx*span.elemsize
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>		greyobject(obj, b, i, span, gcw, idx)
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>	}
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>}
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span><span class="comment">// Shade the object if it isn&#39;t already.</span>
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span><span class="comment">// The object is not nil and known to be in the heap.</span>
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span><span class="comment">// Preemption must be disabled.</span>
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>func shade(b uintptr) {
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>	if obj, span, objIndex := findObject(b, 0, 0); obj != 0 {
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span>		gcw := &amp;getg().m.p.ptr().gcw
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span>		greyobject(obj, 0, 0, span, gcw, objIndex)
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span>	}
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span>}
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span>
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span><span class="comment">// obj is the start of an object with mark mbits.</span>
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span><span class="comment">// If it isn&#39;t already marked, mark it and enqueue into gcw.</span>
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span><span class="comment">// base and off are for debugging only and could be removed.</span>
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span><span class="comment">// See also wbBufFlush1, which partially duplicates this logic.</span>
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>func greyobject(obj, base, off uintptr, span *mspan, gcw *gcWork, objIndex uintptr) {
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>	<span class="comment">// obj should be start of allocation, and so must be at least pointer-aligned.</span>
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>	if obj&amp;(goarch.PtrSize-1) != 0 {
<span id="L1618" class="ln">  1618&nbsp;&nbsp;</span>		throw(&#34;greyobject: obj not pointer-aligned&#34;)
<span id="L1619" class="ln">  1619&nbsp;&nbsp;</span>	}
<span id="L1620" class="ln">  1620&nbsp;&nbsp;</span>	mbits := span.markBitsForIndex(objIndex)
<span id="L1621" class="ln">  1621&nbsp;&nbsp;</span>
<span id="L1622" class="ln">  1622&nbsp;&nbsp;</span>	if useCheckmark {
<span id="L1623" class="ln">  1623&nbsp;&nbsp;</span>		if setCheckmark(obj, base, off, mbits) {
<span id="L1624" class="ln">  1624&nbsp;&nbsp;</span>			<span class="comment">// Already marked.</span>
<span id="L1625" class="ln">  1625&nbsp;&nbsp;</span>			return
<span id="L1626" class="ln">  1626&nbsp;&nbsp;</span>		}
<span id="L1627" class="ln">  1627&nbsp;&nbsp;</span>	} else {
<span id="L1628" class="ln">  1628&nbsp;&nbsp;</span>		if debug.gccheckmark &gt; 0 &amp;&amp; span.isFree(objIndex) {
<span id="L1629" class="ln">  1629&nbsp;&nbsp;</span>			print(&#34;runtime: marking free object &#34;, hex(obj), &#34; found at *(&#34;, hex(base), &#34;+&#34;, hex(off), &#34;)\n&#34;)
<span id="L1630" class="ln">  1630&nbsp;&nbsp;</span>			gcDumpObject(&#34;base&#34;, base, off)
<span id="L1631" class="ln">  1631&nbsp;&nbsp;</span>			gcDumpObject(&#34;obj&#34;, obj, ^uintptr(0))
<span id="L1632" class="ln">  1632&nbsp;&nbsp;</span>			getg().m.traceback = 2
<span id="L1633" class="ln">  1633&nbsp;&nbsp;</span>			throw(&#34;marking free object&#34;)
<span id="L1634" class="ln">  1634&nbsp;&nbsp;</span>		}
<span id="L1635" class="ln">  1635&nbsp;&nbsp;</span>
<span id="L1636" class="ln">  1636&nbsp;&nbsp;</span>		<span class="comment">// If marked we have nothing to do.</span>
<span id="L1637" class="ln">  1637&nbsp;&nbsp;</span>		if mbits.isMarked() {
<span id="L1638" class="ln">  1638&nbsp;&nbsp;</span>			return
<span id="L1639" class="ln">  1639&nbsp;&nbsp;</span>		}
<span id="L1640" class="ln">  1640&nbsp;&nbsp;</span>		mbits.setMarked()
<span id="L1641" class="ln">  1641&nbsp;&nbsp;</span>
<span id="L1642" class="ln">  1642&nbsp;&nbsp;</span>		<span class="comment">// Mark span.</span>
<span id="L1643" class="ln">  1643&nbsp;&nbsp;</span>		arena, pageIdx, pageMask := pageIndexOf(span.base())
<span id="L1644" class="ln">  1644&nbsp;&nbsp;</span>		if arena.pageMarks[pageIdx]&amp;pageMask == 0 {
<span id="L1645" class="ln">  1645&nbsp;&nbsp;</span>			atomic.Or8(&amp;arena.pageMarks[pageIdx], pageMask)
<span id="L1646" class="ln">  1646&nbsp;&nbsp;</span>		}
<span id="L1647" class="ln">  1647&nbsp;&nbsp;</span>
<span id="L1648" class="ln">  1648&nbsp;&nbsp;</span>		<span class="comment">// If this is a noscan object, fast-track it to black</span>
<span id="L1649" class="ln">  1649&nbsp;&nbsp;</span>		<span class="comment">// instead of greying it.</span>
<span id="L1650" class="ln">  1650&nbsp;&nbsp;</span>		if span.spanclass.noscan() {
<span id="L1651" class="ln">  1651&nbsp;&nbsp;</span>			gcw.bytesMarked += uint64(span.elemsize)
<span id="L1652" class="ln">  1652&nbsp;&nbsp;</span>			return
<span id="L1653" class="ln">  1653&nbsp;&nbsp;</span>		}
<span id="L1654" class="ln">  1654&nbsp;&nbsp;</span>	}
<span id="L1655" class="ln">  1655&nbsp;&nbsp;</span>
<span id="L1656" class="ln">  1656&nbsp;&nbsp;</span>	<span class="comment">// We&#39;re adding obj to P&#39;s local workbuf, so it&#39;s likely</span>
<span id="L1657" class="ln">  1657&nbsp;&nbsp;</span>	<span class="comment">// this object will be processed soon by the same P.</span>
<span id="L1658" class="ln">  1658&nbsp;&nbsp;</span>	<span class="comment">// Even if the workbuf gets flushed, there will likely still be</span>
<span id="L1659" class="ln">  1659&nbsp;&nbsp;</span>	<span class="comment">// some benefit on platforms with inclusive shared caches.</span>
<span id="L1660" class="ln">  1660&nbsp;&nbsp;</span>	sys.Prefetch(obj)
<span id="L1661" class="ln">  1661&nbsp;&nbsp;</span>	<span class="comment">// Queue the obj for scanning.</span>
<span id="L1662" class="ln">  1662&nbsp;&nbsp;</span>	if !gcw.putFast(obj) {
<span id="L1663" class="ln">  1663&nbsp;&nbsp;</span>		gcw.put(obj)
<span id="L1664" class="ln">  1664&nbsp;&nbsp;</span>	}
<span id="L1665" class="ln">  1665&nbsp;&nbsp;</span>}
<span id="L1666" class="ln">  1666&nbsp;&nbsp;</span>
<span id="L1667" class="ln">  1667&nbsp;&nbsp;</span><span class="comment">// gcDumpObject dumps the contents of obj for debugging and marks the</span>
<span id="L1668" class="ln">  1668&nbsp;&nbsp;</span><span class="comment">// field at byte offset off in obj.</span>
<span id="L1669" class="ln">  1669&nbsp;&nbsp;</span>func gcDumpObject(label string, obj, off uintptr) {
<span id="L1670" class="ln">  1670&nbsp;&nbsp;</span>	s := spanOf(obj)
<span id="L1671" class="ln">  1671&nbsp;&nbsp;</span>	print(label, &#34;=&#34;, hex(obj))
<span id="L1672" class="ln">  1672&nbsp;&nbsp;</span>	if s == nil {
<span id="L1673" class="ln">  1673&nbsp;&nbsp;</span>		print(&#34; s=nil\n&#34;)
<span id="L1674" class="ln">  1674&nbsp;&nbsp;</span>		return
<span id="L1675" class="ln">  1675&nbsp;&nbsp;</span>	}
<span id="L1676" class="ln">  1676&nbsp;&nbsp;</span>	print(&#34; s.base()=&#34;, hex(s.base()), &#34; s.limit=&#34;, hex(s.limit), &#34; s.spanclass=&#34;, s.spanclass, &#34; s.elemsize=&#34;, s.elemsize, &#34; s.state=&#34;)
<span id="L1677" class="ln">  1677&nbsp;&nbsp;</span>	if state := s.state.get(); 0 &lt;= state &amp;&amp; int(state) &lt; len(mSpanStateNames) {
<span id="L1678" class="ln">  1678&nbsp;&nbsp;</span>		print(mSpanStateNames[state], &#34;\n&#34;)
<span id="L1679" class="ln">  1679&nbsp;&nbsp;</span>	} else {
<span id="L1680" class="ln">  1680&nbsp;&nbsp;</span>		print(&#34;unknown(&#34;, state, &#34;)\n&#34;)
<span id="L1681" class="ln">  1681&nbsp;&nbsp;</span>	}
<span id="L1682" class="ln">  1682&nbsp;&nbsp;</span>
<span id="L1683" class="ln">  1683&nbsp;&nbsp;</span>	skipped := false
<span id="L1684" class="ln">  1684&nbsp;&nbsp;</span>	size := s.elemsize
<span id="L1685" class="ln">  1685&nbsp;&nbsp;</span>	if s.state.get() == mSpanManual &amp;&amp; size == 0 {
<span id="L1686" class="ln">  1686&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re printing something from a stack frame. We</span>
<span id="L1687" class="ln">  1687&nbsp;&nbsp;</span>		<span class="comment">// don&#39;t know how big it is, so just show up to an</span>
<span id="L1688" class="ln">  1688&nbsp;&nbsp;</span>		<span class="comment">// including off.</span>
<span id="L1689" class="ln">  1689&nbsp;&nbsp;</span>		size = off + goarch.PtrSize
<span id="L1690" class="ln">  1690&nbsp;&nbsp;</span>	}
<span id="L1691" class="ln">  1691&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; size; i += goarch.PtrSize {
<span id="L1692" class="ln">  1692&nbsp;&nbsp;</span>		<span class="comment">// For big objects, just print the beginning (because</span>
<span id="L1693" class="ln">  1693&nbsp;&nbsp;</span>		<span class="comment">// that usually hints at the object&#39;s type) and the</span>
<span id="L1694" class="ln">  1694&nbsp;&nbsp;</span>		<span class="comment">// fields around off.</span>
<span id="L1695" class="ln">  1695&nbsp;&nbsp;</span>		if !(i &lt; 128*goarch.PtrSize || off-16*goarch.PtrSize &lt; i &amp;&amp; i &lt; off+16*goarch.PtrSize) {
<span id="L1696" class="ln">  1696&nbsp;&nbsp;</span>			skipped = true
<span id="L1697" class="ln">  1697&nbsp;&nbsp;</span>			continue
<span id="L1698" class="ln">  1698&nbsp;&nbsp;</span>		}
<span id="L1699" class="ln">  1699&nbsp;&nbsp;</span>		if skipped {
<span id="L1700" class="ln">  1700&nbsp;&nbsp;</span>			print(&#34; ...\n&#34;)
<span id="L1701" class="ln">  1701&nbsp;&nbsp;</span>			skipped = false
<span id="L1702" class="ln">  1702&nbsp;&nbsp;</span>		}
<span id="L1703" class="ln">  1703&nbsp;&nbsp;</span>		print(&#34; *(&#34;, label, &#34;+&#34;, i, &#34;) = &#34;, hex(*(*uintptr)(unsafe.Pointer(obj + i))))
<span id="L1704" class="ln">  1704&nbsp;&nbsp;</span>		if i == off {
<span id="L1705" class="ln">  1705&nbsp;&nbsp;</span>			print(&#34; &lt;==&#34;)
<span id="L1706" class="ln">  1706&nbsp;&nbsp;</span>		}
<span id="L1707" class="ln">  1707&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L1708" class="ln">  1708&nbsp;&nbsp;</span>	}
<span id="L1709" class="ln">  1709&nbsp;&nbsp;</span>	if skipped {
<span id="L1710" class="ln">  1710&nbsp;&nbsp;</span>		print(&#34; ...\n&#34;)
<span id="L1711" class="ln">  1711&nbsp;&nbsp;</span>	}
<span id="L1712" class="ln">  1712&nbsp;&nbsp;</span>}
<span id="L1713" class="ln">  1713&nbsp;&nbsp;</span>
<span id="L1714" class="ln">  1714&nbsp;&nbsp;</span><span class="comment">// gcmarknewobject marks a newly allocated object black. obj must</span>
<span id="L1715" class="ln">  1715&nbsp;&nbsp;</span><span class="comment">// not contain any non-nil pointers.</span>
<span id="L1716" class="ln">  1716&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1717" class="ln">  1717&nbsp;&nbsp;</span><span class="comment">// This is nosplit so it can manipulate a gcWork without preemption.</span>
<span id="L1718" class="ln">  1718&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1719" class="ln">  1719&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L1720" class="ln">  1720&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1721" class="ln">  1721&nbsp;&nbsp;</span>func gcmarknewobject(span *mspan, obj uintptr) {
<span id="L1722" class="ln">  1722&nbsp;&nbsp;</span>	if useCheckmark { <span class="comment">// The world should be stopped so this should not happen.</span>
<span id="L1723" class="ln">  1723&nbsp;&nbsp;</span>		throw(&#34;gcmarknewobject called while doing checkmark&#34;)
<span id="L1724" class="ln">  1724&nbsp;&nbsp;</span>	}
<span id="L1725" class="ln">  1725&nbsp;&nbsp;</span>
<span id="L1726" class="ln">  1726&nbsp;&nbsp;</span>	<span class="comment">// Mark object.</span>
<span id="L1727" class="ln">  1727&nbsp;&nbsp;</span>	objIndex := span.objIndex(obj)
<span id="L1728" class="ln">  1728&nbsp;&nbsp;</span>	span.markBitsForIndex(objIndex).setMarked()
<span id="L1729" class="ln">  1729&nbsp;&nbsp;</span>
<span id="L1730" class="ln">  1730&nbsp;&nbsp;</span>	<span class="comment">// Mark span.</span>
<span id="L1731" class="ln">  1731&nbsp;&nbsp;</span>	arena, pageIdx, pageMask := pageIndexOf(span.base())
<span id="L1732" class="ln">  1732&nbsp;&nbsp;</span>	if arena.pageMarks[pageIdx]&amp;pageMask == 0 {
<span id="L1733" class="ln">  1733&nbsp;&nbsp;</span>		atomic.Or8(&amp;arena.pageMarks[pageIdx], pageMask)
<span id="L1734" class="ln">  1734&nbsp;&nbsp;</span>	}
<span id="L1735" class="ln">  1735&nbsp;&nbsp;</span>
<span id="L1736" class="ln">  1736&nbsp;&nbsp;</span>	gcw := &amp;getg().m.p.ptr().gcw
<span id="L1737" class="ln">  1737&nbsp;&nbsp;</span>	gcw.bytesMarked += uint64(span.elemsize)
<span id="L1738" class="ln">  1738&nbsp;&nbsp;</span>}
<span id="L1739" class="ln">  1739&nbsp;&nbsp;</span>
<span id="L1740" class="ln">  1740&nbsp;&nbsp;</span><span class="comment">// gcMarkTinyAllocs greys all active tiny alloc blocks.</span>
<span id="L1741" class="ln">  1741&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1742" class="ln">  1742&nbsp;&nbsp;</span><span class="comment">// The world must be stopped.</span>
<span id="L1743" class="ln">  1743&nbsp;&nbsp;</span>func gcMarkTinyAllocs() {
<span id="L1744" class="ln">  1744&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L1745" class="ln">  1745&nbsp;&nbsp;</span>
<span id="L1746" class="ln">  1746&nbsp;&nbsp;</span>	for _, p := range allp {
<span id="L1747" class="ln">  1747&nbsp;&nbsp;</span>		c := p.mcache
<span id="L1748" class="ln">  1748&nbsp;&nbsp;</span>		if c == nil || c.tiny == 0 {
<span id="L1749" class="ln">  1749&nbsp;&nbsp;</span>			continue
<span id="L1750" class="ln">  1750&nbsp;&nbsp;</span>		}
<span id="L1751" class="ln">  1751&nbsp;&nbsp;</span>		_, span, objIndex := findObject(c.tiny, 0, 0)
<span id="L1752" class="ln">  1752&nbsp;&nbsp;</span>		gcw := &amp;p.gcw
<span id="L1753" class="ln">  1753&nbsp;&nbsp;</span>		greyobject(c.tiny, 0, 0, span, gcw, objIndex)
<span id="L1754" class="ln">  1754&nbsp;&nbsp;</span>	}
<span id="L1755" class="ln">  1755&nbsp;&nbsp;</span>}
<span id="L1756" class="ln">  1756&nbsp;&nbsp;</span>
</pre><p><a href="mgcmark.go?m=text">View as plain text</a></p>

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
