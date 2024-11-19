<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mgcpacer.go - Go Documentation Server</title>

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
<a href="mgcpacer.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mgcpacer.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2021 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/cpu&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/goexperiment&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	_ &#34;unsafe&#34; <span class="comment">// for go:linkname</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>const (
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	<span class="comment">// gcGoalUtilization is the goal CPU utilization for</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	<span class="comment">// marking as a fraction of GOMAXPROCS.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// Increasing the goal utilization will shorten GC cycles as the GC</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// has more resources behind it, lessening costs from the write barrier,</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// but comes at the cost of increasing mutator latency.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	gcGoalUtilization = gcBackgroundUtilization
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// gcBackgroundUtilization is the fixed CPU utilization for background</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// marking. It must be &lt;= gcGoalUtilization. The difference between</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// gcGoalUtilization and gcBackgroundUtilization will be made up by</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// mark assists. The scheduler will aim to use within 50% of this</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// goal.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// As a general rule, there&#39;s little reason to set gcBackgroundUtilization</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// &lt; gcGoalUtilization. One reason might be in mostly idle applications,</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// where goroutines are unlikely to assist at all, so the actual</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// utilization will be lower than the goal. But this is moot point</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// because the idle mark workers already soak up idle CPU resources.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// These two values are still kept separate however because they are</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// distinct conceptually, and in previous iterations of the pacer the</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// distinction was more important.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	gcBackgroundUtilization = 0.25
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// gcCreditSlack is the amount of scan work credit that can</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// accumulate locally before updating gcController.heapScanWork and,</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// optionally, gcController.bgScanCredit. Lower values give a more</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// accurate assist ratio and make it more likely that assists will</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// successfully steal background credit. Higher values reduce memory</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// contention.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	gcCreditSlack = 2000
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// gcAssistTimeSlack is the nanoseconds of mutator assist time that</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// can accumulate on a P before updating gcController.assistTime.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	gcAssistTimeSlack = 5000
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// gcOverAssistWork determines how many extra units of scan work a GC</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// assist does when an assist happens. This amortizes the cost of an</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// assist by pre-paying for this many bytes of future allocations.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	gcOverAssistWork = 64 &lt;&lt; 10
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// defaultHeapMinimum is the value of heapMinimum for GOGC==100.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	defaultHeapMinimum = (goexperiment.HeapMinimum512KiBInt)*(512&lt;&lt;10) +
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		(1-goexperiment.HeapMinimum512KiBInt)*(4&lt;&lt;20)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// maxStackScanSlack is the bytes of stack space allocated or freed</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// that can accumulate on a P before updating gcController.stackSize.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	maxStackScanSlack = 8 &lt;&lt; 10
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// memoryLimitMinHeapGoalHeadroom is the minimum amount of headroom the</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// pacer gives to the heap goal when operating in the memory-limited regime.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// That is, it&#39;ll reduce the heap goal by this many extra bytes off of the</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// base calculation, at minimum.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	memoryLimitMinHeapGoalHeadroom = 1 &lt;&lt; 20
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// memoryLimitHeapGoalHeadroomPercent is how headroom the memory-limit-based</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// heap goal should have as a percent of the maximum possible heap goal allowed</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// to maintain the memory limit.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	memoryLimitHeapGoalHeadroomPercent = 3
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// gcController implements the GC pacing controller that determines</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// when to trigger concurrent garbage collection and how much marking</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// work to do in mutator assists and background marking.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// It calculates the ratio between the allocation rate (in terms of CPU</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// time) and the GC scan throughput to determine the heap size at which to</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// trigger a GC cycle such that no GC assists are required to finish on time.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// This algorithm thus optimizes GC CPU utilization to the dedicated background</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// mark utilization of 25% of GOMAXPROCS by minimizing GC assists.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// GOMAXPROCS. The high-level design of this algorithm is documented</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// at https://github.com/golang/proposal/blob/master/design/44167-gc-pacer-redesign.md.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// See https://golang.org/s/go15gcpacing for additional historical context.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>var gcController gcControllerState
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>type gcControllerState struct {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// Initialized from GOGC. GOGC=off means no GC.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	gcPercent atomic.Int32
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// memoryLimit is the soft memory limit in bytes.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// Initialized from GOMEMLIMIT. GOMEMLIMIT=off is equivalent to MaxInt64</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// which means no soft memory limit in practice.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// This is an int64 instead of a uint64 to more easily maintain parity with</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// the SetMemoryLimit API, which sets a maximum at MaxInt64. This value</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// should never be negative.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	memoryLimit atomic.Int64
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// heapMinimum is the minimum heap size at which to trigger GC.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// For small heaps, this overrides the usual GOGC*live set rule.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// When there is a very small live set but a lot of allocation, simply</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">// collecting when the heap reaches GOGC*live results in many GC</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// cycles and high total per-GC overhead. This minimum amortizes this</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// per-GC overhead while keeping the heap reasonably small.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// During initialization this is set to 4MB*GOGC/100. In the case of</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// GOGC==0, this will set heapMinimum to 0, resulting in constant</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// collection even when the heap size is small, which is useful for</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// debugging.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	heapMinimum uint64
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// runway is the amount of runway in heap bytes allocated by the</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// application that we want to give the GC once it starts.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// This is computed from consMark during mark termination.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	runway atomic.Uint64
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// consMark is the estimated per-CPU consMark ratio for the application.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// It represents the ratio between the application&#39;s allocation</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// rate, as bytes allocated per CPU-time, and the GC&#39;s scan rate,</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// as bytes scanned per CPU-time.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// The units of this ratio are (B / cpu-ns) / (B / cpu-ns).</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// At a high level, this value is computed as the bytes of memory</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// allocated (cons) per unit of scan work completed (mark) in a GC</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// cycle, divided by the CPU time spent on each activity.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// Updated at the end of each GC cycle, in endCycle.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	consMark float64
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// lastConsMark is the computed cons/mark value for the previous 4 GC</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// cycles. Note that this is *not* the last value of consMark, but the</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// measured cons/mark value in endCycle.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	lastConsMark [4]float64
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">// gcPercentHeapGoal is the goal heapLive for when next GC ends derived</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// from gcPercent.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// Set to ^uint64(0) if gcPercent is disabled.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	gcPercentHeapGoal atomic.Uint64
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// sweepDistMinTrigger is the minimum trigger to ensure a minimum</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// sweep distance.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// This bound is also special because it applies to both the trigger</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// *and* the goal (all other trigger bounds must be based *on* the goal).</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// It is computed ahead of time, at commit time. The theory is that,</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// absent a sudden change to a parameter like gcPercent, the trigger</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// will be chosen to always give the sweeper enough headroom. However,</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// such a change might dramatically and suddenly move up the trigger,</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// in which case we need to ensure the sweeper still has enough headroom.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	sweepDistMinTrigger atomic.Uint64
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// triggered is the point at which the current GC cycle actually triggered.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// Only valid during the mark phase of a GC cycle, otherwise set to ^uint64(0).</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// Updated while the world is stopped.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	triggered uint64
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// lastHeapGoal is the value of heapGoal at the moment the last GC</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// ended. Note that this is distinct from the last value heapGoal had,</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// because it could change if e.g. gcPercent changes.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// Read and written with the world stopped or with mheap_.lock held.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	lastHeapGoal uint64
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// heapLive is the number of bytes considered live by the GC.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// That is: retained by the most recent GC plus allocated</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// since then. heapLive ≤ memstats.totalAlloc-memstats.totalFree, since</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// heapAlloc includes unmarked objects that have not yet been swept (and</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// hence goes up as we allocate and down as we sweep) while heapLive</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// excludes these objects (and hence only goes up between GCs).</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// To reduce contention, this is updated only when obtaining a span</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// from an mcentral and at this point it counts all of the unallocated</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// slots in that span (which will be allocated before that mcache</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// obtains another span from that mcentral). Hence, it slightly</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// overestimates the &#34;true&#34; live heap size. It&#39;s better to overestimate</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// than to underestimate because 1) this triggers the GC earlier than</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">// necessary rather than potentially too late and 2) this leads to a</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// conservative GC rate rather than a GC rate that is potentially too</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// low.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// Whenever this is updated, call traceHeapAlloc() and</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// this gcControllerState&#39;s revise() method.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	heapLive atomic.Uint64
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// heapScan is the number of bytes of &#34;scannable&#34; heap. This is the</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// live heap (as counted by heapLive), but omitting no-scan objects and</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// no-scan tails of objects.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// This value is fixed at the start of a GC cycle. It represents the</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// maximum scannable heap.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	heapScan atomic.Uint64
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// lastHeapScan is the number of bytes of heap that were scanned</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// last GC cycle. It is the same as heapMarked, but only</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// includes the &#34;scannable&#34; parts of objects.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// Updated when the world is stopped.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	lastHeapScan uint64
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	<span class="comment">// lastStackScan is the number of bytes of stack that were scanned</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// last GC cycle.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	lastStackScan atomic.Uint64
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// maxStackScan is the amount of allocated goroutine stack space in</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	<span class="comment">// use by goroutines.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	<span class="comment">// This number tracks allocated goroutine stack space rather than used</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	<span class="comment">// goroutine stack space (i.e. what is actually scanned) because used</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// goroutine stack space is much harder to measure cheaply. By using</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// allocated space, we make an overestimate; this is OK, it&#39;s better</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">// to conservatively overcount than undercount.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	maxStackScan atomic.Uint64
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	<span class="comment">// globalsScan is the total amount of global variable space</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// that is scannable.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	globalsScan atomic.Uint64
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// heapMarked is the number of bytes marked by the previous</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// GC. After mark termination, heapLive == heapMarked, but</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">// unlike heapLive, heapMarked does not change until the</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">// next mark termination.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	heapMarked uint64
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">// heapScanWork is the total heap scan work performed this cycle.</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// stackScanWork is the total stack scan work performed this cycle.</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// globalsScanWork is the total globals scan work performed this cycle.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	<span class="comment">// These are updated atomically during the cycle. Updates occur in</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	<span class="comment">// bounded batches, since they are both written and read</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// throughout the cycle. At the end of the cycle, heapScanWork is how</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	<span class="comment">// much of the retained heap is scannable.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// Currently these are measured in bytes. For most uses, this is an</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// opaque unit of work, but for estimation the definition is important.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	<span class="comment">// Note that stackScanWork includes only stack space scanned, not all</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// of the allocated stack.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	heapScanWork    atomic.Int64
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	stackScanWork   atomic.Int64
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	globalsScanWork atomic.Int64
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// bgScanCredit is the scan work credit accumulated by the concurrent</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	<span class="comment">// background scan. This credit is accumulated by the background scan</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	<span class="comment">// and stolen by mutator assists.  Updates occur in bounded batches,</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// since it is both written and read throughout the cycle.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	bgScanCredit atomic.Int64
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">// assistTime is the nanoseconds spent in mutator assists</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// during this cycle. This is updated atomically, and must also</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// be updated atomically even during a STW, because it is read</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// by sysmon. Updates occur in bounded batches, since it is both</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// written and read throughout the cycle.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	assistTime atomic.Int64
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	<span class="comment">// dedicatedMarkTime is the nanoseconds spent in dedicated mark workers</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	<span class="comment">// during this cycle. This is updated at the end of the concurrent mark</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// phase.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	dedicatedMarkTime atomic.Int64
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	<span class="comment">// fractionalMarkTime is the nanoseconds spent in the fractional mark</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	<span class="comment">// worker during this cycle. This is updated throughout the cycle and</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	<span class="comment">// will be up-to-date if the fractional mark worker is not currently</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	<span class="comment">// running.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	fractionalMarkTime atomic.Int64
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	<span class="comment">// idleMarkTime is the nanoseconds spent in idle marking during this</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	<span class="comment">// cycle. This is updated throughout the cycle.</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	idleMarkTime atomic.Int64
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	<span class="comment">// markStartTime is the absolute start time in nanoseconds</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	<span class="comment">// that assists and background mark workers started.</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	markStartTime int64
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">// dedicatedMarkWorkersNeeded is the number of dedicated mark workers</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	<span class="comment">// that need to be started. This is computed at the beginning of each</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	<span class="comment">// cycle and decremented as dedicated mark workers get started.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	dedicatedMarkWorkersNeeded atomic.Int64
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">// idleMarkWorkers is two packed int32 values in a single uint64.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// These two values are always updated simultaneously.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// The bottom int32 is the current number of idle mark workers executing.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// The top int32 is the maximum number of idle mark workers allowed to</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">// execute concurrently. Normally, this number is just gomaxprocs. However,</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// during periodic GC cycles it is set to 0 because the system is idle</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	<span class="comment">// anyway; there&#39;s no need to go full blast on all of GOMAXPROCS.</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">// The maximum number of idle mark workers is used to prevent new workers</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">// from starting, but it is not a hard maximum. It is possible (but</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	<span class="comment">// exceedingly rare) for the current number of idle mark workers to</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	<span class="comment">// transiently exceed the maximum. This could happen if the maximum changes</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	<span class="comment">// just after a GC ends, and an M with no P.</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	<span class="comment">// Note that if we have no dedicated mark workers, we set this value to</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	<span class="comment">// 1 in this case we only have fractional GC workers which aren&#39;t scheduled</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	<span class="comment">// strictly enough to ensure GC progress. As a result, idle-priority mark</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	<span class="comment">// workers are vital to GC progress in these situations.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	<span class="comment">// For example, consider a situation in which goroutines block on the GC</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	<span class="comment">// (such as via runtime.GOMAXPROCS) and only fractional mark workers are</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	<span class="comment">// scheduled (e.g. GOMAXPROCS=1). Without idle-priority mark workers, the</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	<span class="comment">// last running M might skip scheduling a fractional mark worker if its</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	<span class="comment">// utilization goal is met, such that once it goes to sleep (because there&#39;s</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// nothing to do), there will be nothing else to spin up a new M for the</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	<span class="comment">// fractional worker in the future, stalling GC progress and causing a</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	<span class="comment">// deadlock. However, idle-priority workers will *always* run when there is</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	<span class="comment">// nothing left to do, ensuring the GC makes progress.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	<span class="comment">// See github.com/golang/go/issues/44163 for more details.</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	idleMarkWorkers atomic.Uint64
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	<span class="comment">// assistWorkPerByte is the ratio of scan work to allocated</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// bytes that should be performed by mutator assists. This is</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">// computed at the beginning of each cycle and updated every</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// time heapScan is updated.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	assistWorkPerByte atomic.Float64
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">// assistBytesPerWork is 1/assistWorkPerByte.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">// Note that because this is read and written independently</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// from assistWorkPerByte users may notice a skew between</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// the two values, and such a state should be safe.</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	assistBytesPerWork atomic.Float64
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// fractionalUtilizationGoal is the fraction of wall clock</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// time that should be spent in the fractional mark worker on</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// each P that isn&#39;t running a dedicated worker.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	<span class="comment">// For example, if the utilization goal is 25% and there are</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	<span class="comment">// no dedicated workers, this will be 0.25. If the goal is</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	<span class="comment">// 25%, there is one dedicated worker, and GOMAXPROCS is 5,</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	<span class="comment">// this will be 0.05 to make up the missing 5%.</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	<span class="comment">// If this is zero, no fractional workers are needed.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	fractionalUtilizationGoal float64
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	<span class="comment">// These memory stats are effectively duplicates of fields from</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	<span class="comment">// memstats.heapStats but are updated atomically or with the world</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	<span class="comment">// stopped and don&#39;t provide the same consistency guarantees.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	<span class="comment">// Because the runtime is responsible for managing a memory limit, it&#39;s</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	<span class="comment">// useful to couple these stats more tightly to the gcController, which</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	<span class="comment">// is intimately connected to how that memory limit is maintained.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	heapInUse    sysMemStat    <span class="comment">// bytes in mSpanInUse spans</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	heapReleased sysMemStat    <span class="comment">// bytes released to the OS</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	heapFree     sysMemStat    <span class="comment">// bytes not in any span, but not released to the OS</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	totalAlloc   atomic.Uint64 <span class="comment">// total bytes allocated</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	totalFree    atomic.Uint64 <span class="comment">// total bytes freed</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	mappedReady  atomic.Uint64 <span class="comment">// total virtual memory in the Ready state (see mem.go).</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	<span class="comment">// test indicates that this is a test-only copy of gcControllerState.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	test bool
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	_ cpu.CacheLinePad
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>func (c *gcControllerState) init(gcPercent int32, memoryLimit int64) {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	c.heapMinimum = defaultHeapMinimum
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	c.triggered = ^uint64(0)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	c.setGCPercent(gcPercent)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	c.setMemoryLimit(memoryLimit)
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	c.commit(true) <span class="comment">// No sweep phase in the first GC cycle.</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	<span class="comment">// N.B. Don&#39;t bother calling traceHeapGoal. Tracing is never enabled at</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	<span class="comment">// initialization time.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	<span class="comment">// N.B. No need to call revise; there&#39;s no GC enabled during</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	<span class="comment">// initialization.</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span><span class="comment">// startCycle resets the GC controller&#39;s state and computes estimates</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span><span class="comment">// for a new GC cycle. The caller must hold worldsema and the world</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span><span class="comment">// must be stopped.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>func (c *gcControllerState) startCycle(markStartTime int64, procs int, trigger gcTrigger) {
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	c.heapScanWork.Store(0)
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	c.stackScanWork.Store(0)
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	c.globalsScanWork.Store(0)
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	c.bgScanCredit.Store(0)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	c.assistTime.Store(0)
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	c.dedicatedMarkTime.Store(0)
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	c.fractionalMarkTime.Store(0)
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	c.idleMarkTime.Store(0)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	c.markStartTime = markStartTime
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	c.triggered = c.heapLive.Load()
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	<span class="comment">// Compute the background mark utilization goal. In general,</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	<span class="comment">// this may not come out exactly. We round the number of</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	<span class="comment">// dedicated workers so that the utilization is closest to</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	<span class="comment">// 25%. For small GOMAXPROCS, this would introduce too much</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	<span class="comment">// error, so we add fractional workers in that case.</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	totalUtilizationGoal := float64(procs) * gcBackgroundUtilization
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	dedicatedMarkWorkersNeeded := int64(totalUtilizationGoal + 0.5)
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	utilError := float64(dedicatedMarkWorkersNeeded)/totalUtilizationGoal - 1
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	const maxUtilError = 0.3
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	if utilError &lt; -maxUtilError || utilError &gt; maxUtilError {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		<span class="comment">// Rounding put us more than 30% off our goal. With</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		<span class="comment">// gcBackgroundUtilization of 25%, this happens for</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		<span class="comment">// GOMAXPROCS&lt;=3 or GOMAXPROCS=6. Enable fractional</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		<span class="comment">// workers to compensate.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		if float64(dedicatedMarkWorkersNeeded) &gt; totalUtilizationGoal {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>			<span class="comment">// Too many dedicated workers.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			dedicatedMarkWorkersNeeded--
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		c.fractionalUtilizationGoal = (totalUtilizationGoal - float64(dedicatedMarkWorkersNeeded)) / float64(procs)
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	} else {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		c.fractionalUtilizationGoal = 0
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	<span class="comment">// In STW mode, we just want dedicated workers.</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	if debug.gcstoptheworld &gt; 0 {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		dedicatedMarkWorkersNeeded = int64(procs)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		c.fractionalUtilizationGoal = 0
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	<span class="comment">// Clear per-P state</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	for _, p := range allp {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		p.gcAssistTime = 0
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		p.gcFractionalMarkTime = 0
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	if trigger.kind == gcTriggerTime {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		<span class="comment">// During a periodic GC cycle, reduce the number of idle mark workers</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		<span class="comment">// required. However, we need at least one dedicated mark worker or</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		<span class="comment">// idle GC worker to ensure GC progress in some scenarios (see comment</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		<span class="comment">// on maxIdleMarkWorkers).</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		if dedicatedMarkWorkersNeeded &gt; 0 {
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			c.setMaxIdleMarkWorkers(0)
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		} else {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>			<span class="comment">// TODO(mknyszek): The fundamental reason why we need this is because</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			<span class="comment">// we can&#39;t count on the fractional mark worker to get scheduled.</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>			<span class="comment">// Fix that by ensuring it gets scheduled according to its quota even</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			<span class="comment">// if the rest of the application is idle.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>			c.setMaxIdleMarkWorkers(1)
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	} else {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		<span class="comment">// N.B. gomaxprocs and dedicatedMarkWorkersNeeded are guaranteed not to</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		<span class="comment">// change during a GC cycle.</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		c.setMaxIdleMarkWorkers(int32(procs) - int32(dedicatedMarkWorkersNeeded))
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	<span class="comment">// Compute initial values for controls that are updated</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	<span class="comment">// throughout the cycle.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	c.dedicatedMarkWorkersNeeded.Store(dedicatedMarkWorkersNeeded)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	c.revise()
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	if debug.gcpacertrace &gt; 0 {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		heapGoal := c.heapGoal()
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		assistRatio := c.assistWorkPerByte.Load()
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		print(&#34;pacer: assist ratio=&#34;, assistRatio,
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			&#34; (scan &#34;, gcController.heapScan.Load()&gt;&gt;20, &#34; MB in &#34;,
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			work.initialHeapLive&gt;&gt;20, &#34;-&gt;&#34;,
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>			heapGoal&gt;&gt;20, &#34; MB)&#34;,
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			&#34; workers=&#34;, dedicatedMarkWorkersNeeded,
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			&#34;+&#34;, c.fractionalUtilizationGoal, &#34;\n&#34;)
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span><span class="comment">// revise updates the assist ratio during the GC cycle to account for</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span><span class="comment">// improved estimates. This should be called whenever gcController.heapScan,</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span><span class="comment">// gcController.heapLive, or if any inputs to gcController.heapGoal are</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span><span class="comment">// updated. It is safe to call concurrently, but it may race with other</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span><span class="comment">// calls to revise.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span><span class="comment">// The result of this race is that the two assist ratio values may not line</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span><span class="comment">// up or may be stale. In practice this is OK because the assist ratio</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span><span class="comment">// moves slowly throughout a GC cycle, and the assist ratio is a best-effort</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span><span class="comment">// heuristic anyway. Furthermore, no part of the heuristic depends on</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span><span class="comment">// the two assist ratio values being exact reciprocals of one another, since</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span><span class="comment">// the two values are used to convert values from different sources.</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span><span class="comment">// The worst case result of this raciness is that we may miss a larger shift</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span><span class="comment">// in the ratio (say, if we decide to pace more aggressively against the</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span><span class="comment">// hard heap goal) but even this &#34;hard goal&#34; is best-effort (see #40460).</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span><span class="comment">// The dedicated GC should ensure we don&#39;t exceed the hard goal by too much</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span><span class="comment">// in the rare case we do exceed it.</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span><span class="comment">// It should only be called when gcBlackenEnabled != 0 (because this</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span><span class="comment">// is when assists are enabled and the necessary statistics are</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span><span class="comment">// available).</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>func (c *gcControllerState) revise() {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	gcPercent := c.gcPercent.Load()
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	if gcPercent &lt; 0 {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		<span class="comment">// If GC is disabled but we&#39;re running a forced GC,</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		<span class="comment">// act like GOGC is huge for the below calculations.</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		gcPercent = 100000
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	live := c.heapLive.Load()
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	scan := c.heapScan.Load()
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	work := c.heapScanWork.Load() + c.stackScanWork.Load() + c.globalsScanWork.Load()
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	<span class="comment">// Assume we&#39;re under the soft goal. Pace GC to complete at</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	<span class="comment">// heapGoal assuming the heap is in steady-state.</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	heapGoal := int64(c.heapGoal())
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	<span class="comment">// The expected scan work is computed as the amount of bytes scanned last</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	<span class="comment">// GC cycle (both heap and stack), plus our estimate of globals work for this cycle.</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	scanWorkExpected := int64(c.lastHeapScan + c.lastStackScan.Load() + c.globalsScan.Load())
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	<span class="comment">// maxScanWork is a worst-case estimate of the amount of scan work that</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	<span class="comment">// needs to be performed in this GC cycle. Specifically, it represents</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	<span class="comment">// the case where *all* scannable memory turns out to be live, and</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	<span class="comment">// *all* allocated stack space is scannable.</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	maxStackScan := c.maxStackScan.Load()
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	maxScanWork := int64(scan + maxStackScan + c.globalsScan.Load())
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	if work &gt; scanWorkExpected {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		<span class="comment">// We&#39;ve already done more scan work than expected. Because our expectation</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		<span class="comment">// is based on a steady-state scannable heap size, we assume this means our</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		<span class="comment">// heap is growing. Compute a new heap goal that takes our existing runway</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		<span class="comment">// computed for scanWorkExpected and extrapolates it to maxScanWork, the worst-case</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		<span class="comment">// scan work. This keeps our assist ratio stable if the heap continues to grow.</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		<span class="comment">// The effect of this mechanism is that assists stay flat in the face of heap</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		<span class="comment">// growths. It&#39;s OK to use more memory this cycle to scan all the live heap,</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		<span class="comment">// because the next GC cycle is inevitably going to use *at least* that much</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		<span class="comment">// memory anyway.</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		extHeapGoal := int64(float64(heapGoal-int64(c.triggered))/float64(scanWorkExpected)*float64(maxScanWork)) + int64(c.triggered)
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		scanWorkExpected = maxScanWork
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		<span class="comment">// hardGoal is a hard limit on the amount that we&#39;re willing to push back the</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		<span class="comment">// heap goal, and that&#39;s twice the heap goal (i.e. if GOGC=100 and the heap and/or</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		<span class="comment">// stacks and/or globals grow to twice their size, this limits the current GC cycle&#39;s</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		<span class="comment">// growth to 4x the original live heap&#39;s size).</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		<span class="comment">// This maintains the invariant that we use no more memory than the next GC cycle</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		<span class="comment">// will anyway.</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		hardGoal := int64((1.0 + float64(gcPercent)/100.0) * float64(heapGoal))
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		if extHeapGoal &gt; hardGoal {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>			extHeapGoal = hardGoal
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		heapGoal = extHeapGoal
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	}
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	if int64(live) &gt; heapGoal {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re already past our heap goal, even the extrapolated one.</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		<span class="comment">// Leave ourselves some extra runway, so in the worst case we</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		<span class="comment">// finish by that point.</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		const maxOvershoot = 1.1
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		heapGoal = int64(float64(heapGoal) * maxOvershoot)
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		<span class="comment">// Compute the upper bound on the scan work remaining.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		scanWorkExpected = maxScanWork
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	<span class="comment">// Compute the remaining scan work estimate.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	<span class="comment">// Note that we currently count allocations during GC as both</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	<span class="comment">// scannable heap (heapScan) and scan work completed</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	<span class="comment">// (scanWork), so allocation will change this difference</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	<span class="comment">// slowly in the soft regime and not at all in the hard</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	<span class="comment">// regime.</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	scanWorkRemaining := scanWorkExpected - work
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	if scanWorkRemaining &lt; 1000 {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		<span class="comment">// We set a somewhat arbitrary lower bound on</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		<span class="comment">// remaining scan work since if we aim a little high,</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		<span class="comment">// we can miss by a little.</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		<span class="comment">// We *do* need to enforce that this is at least 1,</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		<span class="comment">// since marking is racy and double-scanning objects</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		<span class="comment">// may legitimately make the remaining scan work</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		<span class="comment">// negative, even in the hard goal regime.</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		scanWorkRemaining = 1000
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	<span class="comment">// Compute the heap distance remaining.</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	heapRemaining := heapGoal - int64(live)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	if heapRemaining &lt;= 0 {
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		<span class="comment">// This shouldn&#39;t happen, but if it does, avoid</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		<span class="comment">// dividing by zero or setting the assist negative.</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		heapRemaining = 1
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	<span class="comment">// Compute the mutator assist ratio so by the time the mutator</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	<span class="comment">// allocates the remaining heap bytes up to heapGoal, it will</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	<span class="comment">// have done (or stolen) the remaining amount of scan work.</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	<span class="comment">// Note that the assist ratio values are updated atomically</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	<span class="comment">// but not together. This means there may be some degree of</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	<span class="comment">// skew between the two values. This is generally OK as the</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	<span class="comment">// values shift relatively slowly over the course of a GC</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	<span class="comment">// cycle.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	assistWorkPerByte := float64(scanWorkRemaining) / float64(heapRemaining)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	assistBytesPerWork := float64(heapRemaining) / float64(scanWorkRemaining)
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	c.assistWorkPerByte.Store(assistWorkPerByte)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	c.assistBytesPerWork.Store(assistBytesPerWork)
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span><span class="comment">// endCycle computes the consMark estimate for the next cycle.</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span><span class="comment">// userForced indicates whether the current GC cycle was forced</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span><span class="comment">// by the application.</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>func (c *gcControllerState) endCycle(now int64, procs int, userForced bool) {
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	<span class="comment">// Record last heap goal for the scavenger.</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	<span class="comment">// We&#39;ll be updating the heap goal soon.</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	gcController.lastHeapGoal = c.heapGoal()
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	<span class="comment">// Compute the duration of time for which assists were turned on.</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	assistDuration := now - c.markStartTime
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	<span class="comment">// Assume background mark hit its utilization goal.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	utilization := gcBackgroundUtilization
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	<span class="comment">// Add assist utilization; avoid divide by zero.</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	if assistDuration &gt; 0 {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		utilization += float64(c.assistTime.Load()) / float64(assistDuration*int64(procs))
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	}
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	if c.heapLive.Load() &lt;= c.triggered {
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		<span class="comment">// Shouldn&#39;t happen, but let&#39;s be very safe about this in case the</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		<span class="comment">// GC is somehow extremely short.</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		<span class="comment">// In this case though, the only reasonable value for c.heapLive-c.triggered</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		<span class="comment">// would be 0, which isn&#39;t really all that useful, i.e. the GC was so short</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		<span class="comment">// that it didn&#39;t matter.</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		<span class="comment">// Ignore this case and don&#39;t update anything.</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		return
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	idleUtilization := 0.0
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	if assistDuration &gt; 0 {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		idleUtilization = float64(c.idleMarkTime.Load()) / float64(assistDuration*int64(procs))
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	<span class="comment">// Determine the cons/mark ratio.</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	<span class="comment">// The units we want for the numerator and denominator are both B / cpu-ns.</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	<span class="comment">// We get this by taking the bytes allocated or scanned, and divide by the amount of</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	<span class="comment">// CPU time it took for those operations. For allocations, that CPU time is</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	<span class="comment">//    assistDuration * procs * (1 - utilization)</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	<span class="comment">// Where utilization includes just background GC workers and assists. It does *not*</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	<span class="comment">// include idle GC work time, because in theory the mutator is free to take that at</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	<span class="comment">// any point.</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	<span class="comment">// For scanning, that CPU time is</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	<span class="comment">//    assistDuration * procs * (utilization + idleUtilization)</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	<span class="comment">// In this case, we *include* idle utilization, because that is additional CPU time that</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	<span class="comment">// the GC had available to it.</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	<span class="comment">// In effect, idle GC time is sort of double-counted here, but it&#39;s very weird compared</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	<span class="comment">// to other kinds of GC work, because of how fluid it is. Namely, because the mutator is</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	<span class="comment">// *always* free to take it.</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	<span class="comment">// So this calculation is really:</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	<span class="comment">//     (heapLive-trigger) / (assistDuration * procs * (1-utilization)) /</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	<span class="comment">//         (scanWork) / (assistDuration * procs * (utilization+idleUtilization))</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	<span class="comment">// Note that because we only care about the ratio, assistDuration and procs cancel out.</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	scanWork := c.heapScanWork.Load() + c.stackScanWork.Load() + c.globalsScanWork.Load()
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	currentConsMark := (float64(c.heapLive.Load()-c.triggered) * (utilization + idleUtilization)) /
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		(float64(scanWork) * (1 - utilization))
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	<span class="comment">// Update our cons/mark estimate. This is the maximum of the value we just computed and the last</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	<span class="comment">// 4 cons/mark values we measured. The reason we take the maximum here is to bias a noisy</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	<span class="comment">// cons/mark measurement toward fewer assists at the expense of additional GC cycles (starting</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	<span class="comment">// earlier).</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	oldConsMark := c.consMark
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	c.consMark = currentConsMark
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	for i := range c.lastConsMark {
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		if c.lastConsMark[i] &gt; c.consMark {
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>			c.consMark = c.lastConsMark[i]
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		}
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	copy(c.lastConsMark[:], c.lastConsMark[1:])
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	c.lastConsMark[len(c.lastConsMark)-1] = currentConsMark
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	if debug.gcpacertrace &gt; 0 {
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>		printlock()
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		goal := gcGoalUtilization * 100
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		print(&#34;pacer: &#34;, int(utilization*100), &#34;% CPU (&#34;, int(goal), &#34; exp.) for &#34;)
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		print(c.heapScanWork.Load(), &#34;+&#34;, c.stackScanWork.Load(), &#34;+&#34;, c.globalsScanWork.Load(), &#34; B work (&#34;, c.lastHeapScan+c.lastStackScan.Load()+c.globalsScan.Load(), &#34; B exp.) &#34;)
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		live := c.heapLive.Load()
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		print(&#34;in &#34;, c.triggered, &#34; B -&gt; &#34;, live, &#34; B (∆goal &#34;, int64(live)-int64(c.lastHeapGoal), &#34;, cons/mark &#34;, oldConsMark, &#34;)&#34;)
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		println()
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		printunlock()
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	}
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span><span class="comment">// enlistWorker encourages another dedicated mark worker to start on</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span><span class="comment">// another P if there are spare worker slots. It is used by putfull</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span><span class="comment">// when more work is made available.</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>func (c *gcControllerState) enlistWorker() {
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	<span class="comment">// If there are idle Ps, wake one so it will run an idle worker.</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	<span class="comment">// NOTE: This is suspected of causing deadlocks. See golang.org/issue/19112.</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	<span class="comment">//	if sched.npidle.Load() != 0 &amp;&amp; sched.nmspinning.Load() == 0 {</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	<span class="comment">//		wakep()</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	<span class="comment">//		return</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	<span class="comment">//	}</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	<span class="comment">// There are no idle Ps. If we need more dedicated workers,</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	<span class="comment">// try to preempt a running P so it will switch to a worker.</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	if c.dedicatedMarkWorkersNeeded.Load() &lt;= 0 {
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>		return
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	}
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	<span class="comment">// Pick a random other P to preempt.</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	if gomaxprocs &lt;= 1 {
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		return
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	}
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	gp := getg()
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	if gp == nil || gp.m == nil || gp.m.p == 0 {
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>		return
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	}
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	myID := gp.m.p.ptr().id
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	for tries := 0; tries &lt; 5; tries++ {
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		id := int32(cheaprandn(uint32(gomaxprocs - 1)))
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		if id &gt;= myID {
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>			id++
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>		}
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>		p := allp[id]
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		if p.status != _Prunning {
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>			continue
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		}
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		if preemptone(p) {
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>			return
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		}
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	}
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>}
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span><span class="comment">// findRunnableGCWorker returns a background mark worker for pp if it</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span><span class="comment">// should be run. This must only be called when gcBlackenEnabled != 0.</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>func (c *gcControllerState) findRunnableGCWorker(pp *p, now int64) (*g, int64) {
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	if gcBlackenEnabled == 0 {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>		throw(&#34;gcControllerState.findRunnable: blackening not enabled&#34;)
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	}
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	<span class="comment">// Since we have the current time, check if the GC CPU limiter</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	<span class="comment">// hasn&#39;t had an update in a while. This check is necessary in</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	<span class="comment">// case the limiter is on but hasn&#39;t been checked in a while and</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	<span class="comment">// so may have left sufficient headroom to turn off again.</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	if now == 0 {
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		now = nanotime()
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	}
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	if gcCPULimiter.needUpdate(now) {
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		gcCPULimiter.update(now)
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	}
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	if !gcMarkWorkAvailable(pp) {
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		<span class="comment">// No work to be done right now. This can happen at</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		<span class="comment">// the end of the mark phase when there are still</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		<span class="comment">// assists tapering off. Don&#39;t bother running a worker</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		<span class="comment">// now because it&#39;ll just return immediately.</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>		return nil, now
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	}
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	<span class="comment">// Grab a worker before we commit to running below.</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	node := (*gcBgMarkWorkerNode)(gcBgMarkWorkerPool.pop())
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	if node == nil {
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		<span class="comment">// There is at least one worker per P, so normally there are</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		<span class="comment">// enough workers to run on all Ps, if necessary. However, once</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		<span class="comment">// a worker enters gcMarkDone it may park without rejoining the</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		<span class="comment">// pool, thus freeing a P with no corresponding worker.</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		<span class="comment">// gcMarkDone never depends on another worker doing work, so it</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		<span class="comment">// is safe to simply do nothing here.</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		<span class="comment">// If gcMarkDone bails out without completing the mark phase,</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		<span class="comment">// it will always do so with queued global work. Thus, that P</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>		<span class="comment">// will be immediately eligible to re-run the worker G it was</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>		<span class="comment">// just using, ensuring work can complete.</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		return nil, now
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	}
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	decIfPositive := func(val *atomic.Int64) bool {
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		for {
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>			v := val.Load()
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>			if v &lt;= 0 {
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>				return false
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>			}
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>			if val.CompareAndSwap(v, v-1) {
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>				return true
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>			}
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>		}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	}
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	if decIfPositive(&amp;c.dedicatedMarkWorkersNeeded) {
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		<span class="comment">// This P is now dedicated to marking until the end of</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		<span class="comment">// the concurrent mark phase.</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>		pp.gcMarkWorkerMode = gcMarkWorkerDedicatedMode
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	} else if c.fractionalUtilizationGoal == 0 {
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>		<span class="comment">// No need for fractional workers.</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>		gcBgMarkWorkerPool.push(&amp;node.node)
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		return nil, now
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	} else {
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>		<span class="comment">// Is this P behind on the fractional utilization</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>		<span class="comment">// goal?</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>		<span class="comment">// This should be kept in sync with pollFractionalWorkerExit.</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		delta := now - c.markStartTime
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		if delta &gt; 0 &amp;&amp; float64(pp.gcFractionalMarkTime)/float64(delta) &gt; c.fractionalUtilizationGoal {
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			<span class="comment">// Nope. No need to run a fractional worker.</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>			gcBgMarkWorkerPool.push(&amp;node.node)
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>			return nil, now
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		}
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>		<span class="comment">// Run a fractional worker.</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>		pp.gcMarkWorkerMode = gcMarkWorkerFractionalMode
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	}
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	<span class="comment">// Run the background mark worker.</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	gp := node.gp.ptr()
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	trace := traceAcquire()
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	casgstatus(gp, _Gwaiting, _Grunnable)
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		trace.GoUnpark(gp, 0)
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	}
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	return gp, now
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>}
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span><span class="comment">// resetLive sets up the controller state for the next mark phase after the end</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span><span class="comment">// of the previous one. Must be called after endCycle and before commit, before</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span><span class="comment">// the world is started.</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span><span class="comment">// The world must be stopped.</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>func (c *gcControllerState) resetLive(bytesMarked uint64) {
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	c.heapMarked = bytesMarked
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	c.heapLive.Store(bytesMarked)
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	c.heapScan.Store(uint64(c.heapScanWork.Load()))
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>	c.lastHeapScan = uint64(c.heapScanWork.Load())
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	c.lastStackScan.Store(uint64(c.stackScanWork.Load()))
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	c.triggered = ^uint64(0) <span class="comment">// Reset triggered.</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	<span class="comment">// heapLive was updated, so emit a trace event.</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	trace := traceAcquire()
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		trace.HeapAlloc(bytesMarked)
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	}
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>}
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span><span class="comment">// markWorkerStop must be called whenever a mark worker stops executing.</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L842" class="ln">   842&nbsp;&nbsp;</span><span class="comment">// It updates mark work accounting in the controller by a duration of</span>
<span id="L843" class="ln">   843&nbsp;&nbsp;</span><span class="comment">// work in nanoseconds and other bookkeeping.</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span><span class="comment">// Safe to execute at any time.</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>func (c *gcControllerState) markWorkerStop(mode gcMarkWorkerMode, duration int64) {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	switch mode {
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	case gcMarkWorkerDedicatedMode:
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		c.dedicatedMarkTime.Add(duration)
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		c.dedicatedMarkWorkersNeeded.Add(1)
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	case gcMarkWorkerFractionalMode:
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		c.fractionalMarkTime.Add(duration)
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	case gcMarkWorkerIdleMode:
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		c.idleMarkTime.Add(duration)
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>		c.removeIdleMarkWorker()
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	default:
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>		throw(&#34;markWorkerStop: unknown mark worker mode&#34;)
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	}
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>}
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>func (c *gcControllerState) update(dHeapLive, dHeapScan int64) {
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	if dHeapLive != 0 {
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		trace := traceAcquire()
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		live := gcController.heapLive.Add(dHeapLive)
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		if trace.ok() {
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>			<span class="comment">// gcController.heapLive changed.</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>			trace.HeapAlloc(live)
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>			traceRelease(trace)
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>		}
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	}
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	if gcBlackenEnabled == 0 {
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		<span class="comment">// Update heapScan when we&#39;re not in a current GC. It is fixed</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		<span class="comment">// at the beginning of a cycle.</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		if dHeapScan != 0 {
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>			gcController.heapScan.Add(dHeapScan)
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>		}
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	} else {
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>		<span class="comment">// gcController.heapLive changed.</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>		c.revise()
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	}
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>}
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>func (c *gcControllerState) addScannableStack(pp *p, amount int64) {
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	if pp == nil {
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>		c.maxStackScan.Add(amount)
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>		return
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	}
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>	pp.maxStackScanDelta += amount
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>	if pp.maxStackScanDelta &gt;= maxStackScanSlack || pp.maxStackScanDelta &lt;= -maxStackScanSlack {
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>		c.maxStackScan.Add(pp.maxStackScanDelta)
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>		pp.maxStackScanDelta = 0
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	}
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>}
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>func (c *gcControllerState) addGlobals(amount int64) {
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	c.globalsScan.Add(amount)
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>}
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span><span class="comment">// heapGoal returns the current heap goal.</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>func (c *gcControllerState) heapGoal() uint64 {
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	goal, _ := c.heapGoalInternal()
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	return goal
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>}
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span><span class="comment">// heapGoalInternal is the implementation of heapGoal which returns additional</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span><span class="comment">// information that is necessary for computing the trigger.</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span><span class="comment">// The returned minTrigger is always &lt;= goal.</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>func (c *gcControllerState) heapGoalInternal() (goal, minTrigger uint64) {
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	<span class="comment">// Start with the goal calculated for gcPercent.</span>
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>	goal = c.gcPercentHeapGoal.Load()
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	<span class="comment">// Check if the memory-limit-based goal is smaller, and if so, pick that.</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	if newGoal := c.memoryLimitHeapGoal(); newGoal &lt; goal {
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		goal = newGoal
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>	} else {
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re not limited by the memory limit goal, so perform a series of</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		<span class="comment">// adjustments that might move the goal forward in a variety of circumstances.</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>		sweepDistTrigger := c.sweepDistMinTrigger.Load()
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		if sweepDistTrigger &gt; goal {
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>			<span class="comment">// Set the goal to maintain a minimum sweep distance since</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>			<span class="comment">// the last call to commit. Note that we never want to do this</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>			<span class="comment">// if we&#39;re in the memory limit regime, because it could push</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>			<span class="comment">// the goal up.</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>			goal = sweepDistTrigger
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>		}
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>		<span class="comment">// Since we ignore the sweep distance trigger in the memory</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		<span class="comment">// limit regime, we need to ensure we don&#39;t propagate it to</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>		<span class="comment">// the trigger, because it could cause a violation of the</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>		<span class="comment">// invariant that the trigger &lt; goal.</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>		minTrigger = sweepDistTrigger
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>		<span class="comment">// Ensure that the heap goal is at least a little larger than</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>		<span class="comment">// the point at which we triggered. This may not be the case if GC</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>		<span class="comment">// start is delayed or if the allocation that pushed gcController.heapLive</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>		<span class="comment">// over trigger is large or if the trigger is really close to</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>		<span class="comment">// GOGC. Assist is proportional to this distance, so enforce a</span>
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>		<span class="comment">// minimum distance, even if it means going over the GOGC goal</span>
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>		<span class="comment">// by a tiny bit.</span>
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		<span class="comment">// Ignore this if we&#39;re in the memory limit regime: we&#39;d prefer to</span>
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>		<span class="comment">// have the GC respond hard about how close we are to the goal than to</span>
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>		<span class="comment">// push the goal back in such a manner that it could cause us to exceed</span>
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>		<span class="comment">// the memory limit.</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>		const minRunway = 64 &lt;&lt; 10
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>		if c.triggered != ^uint64(0) &amp;&amp; goal &lt; c.triggered+minRunway {
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>			goal = c.triggered + minRunway
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		}
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>	}
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>	return
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>}
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>
<span id="L954" class="ln">   954&nbsp;&nbsp;</span><span class="comment">// memoryLimitHeapGoal returns a heap goal derived from memoryLimit.</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>func (c *gcControllerState) memoryLimitHeapGoal() uint64 {
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>	<span class="comment">// Start by pulling out some values we&#39;ll need. Be careful about overflow.</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	var heapFree, heapAlloc, mappedReady uint64
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>	for {
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		heapFree = c.heapFree.load()                         <span class="comment">// Free and unscavenged memory.</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>		heapAlloc = c.totalAlloc.Load() - c.totalFree.Load() <span class="comment">// Heap object bytes in use.</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>		mappedReady = c.mappedReady.Load()                   <span class="comment">// Total unreleased mapped memory.</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>		if heapFree+heapAlloc &lt;= mappedReady {
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>			break
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>		}
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>		<span class="comment">// It is impossible for total unreleased mapped memory to exceed heap memory, but</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		<span class="comment">// because these stats are updated independently, we may observe a partial update</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>		<span class="comment">// including only some values. Thus, we appear to break the invariant. However,</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>		<span class="comment">// this condition is necessarily transient, so just try again. In the case of a</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>		<span class="comment">// persistent accounting error, we&#39;ll deadlock here.</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	}
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>	<span class="comment">// Below we compute a goal from memoryLimit. There are a few things to be aware of.</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	<span class="comment">// Firstly, the memoryLimit does not easily compare to the heap goal: the former</span>
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>	<span class="comment">// is total mapped memory by the runtime that hasn&#39;t been released, while the latter is</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>	<span class="comment">// only heap object memory. Intuitively, the way we convert from one to the other is to</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	<span class="comment">// subtract everything from memoryLimit that both contributes to the memory limit (so,</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>	<span class="comment">// ignore scavenged memory) and doesn&#39;t contain heap objects. This isn&#39;t quite what</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	<span class="comment">// lines up with reality, but it&#39;s a good starting point.</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	<span class="comment">// In practice this computation looks like the following:</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	<span class="comment">//    goal := memoryLimit - ((mappedReady - heapFree - heapAlloc) + max(mappedReady - memoryLimit, 0))</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>	<span class="comment">//                    ^1                                    ^2</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>	<span class="comment">//    goal -= goal / 100 * memoryLimitHeapGoalHeadroomPercent</span>
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	<span class="comment">//    ^3</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	<span class="comment">// Let&#39;s break this down.</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	<span class="comment">// The first term (marker 1) is everything that contributes to the memory limit and isn&#39;t</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	<span class="comment">// or couldn&#39;t become heap objects. It represents, broadly speaking, non-heap overheads.</span>
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>	<span class="comment">// One oddity you may have noticed is that we also subtract out heapFree, i.e. unscavenged</span>
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>	<span class="comment">// memory that may contain heap objects in the future.</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	<span class="comment">// Let&#39;s take a step back. In an ideal world, this term would look something like just</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>	<span class="comment">// the heap goal. That is, we &#34;reserve&#34; enough space for the heap to grow to the heap</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	<span class="comment">// goal, and subtract out everything else. This is of course impossible; the definition</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	<span class="comment">// is circular! However, this impossible definition contains a key insight: the amount</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	<span class="comment">// we&#39;re *going* to use matters just as much as whatever we&#39;re currently using.</span>
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	<span class="comment">// Consider if the heap shrinks to 1/10th its size, leaving behind lots of free and</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>	<span class="comment">// unscavenged memory. mappedReady - heapAlloc will be quite large, because of that free</span>
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>	<span class="comment">// and unscavenged memory, pushing the goal down significantly.</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>	<span class="comment">// heapFree is also safe to exclude from the memory limit because in the steady-state, it&#39;s</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	<span class="comment">// just a pool of memory for future heap allocations, and making new allocations from heapFree</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	<span class="comment">// memory doesn&#39;t increase overall memory use. In transient states, the scavenger and the</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	<span class="comment">// allocator actively manage the pool of heapFree memory to maintain the memory limit.</span>
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>	<span class="comment">// The second term (marker 2) is the amount of memory we&#39;ve exceeded the limit by, and is</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>	<span class="comment">// intended to help recover from such a situation. By pushing the heap goal down, we also</span>
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>	<span class="comment">// push the trigger down, triggering and finishing a GC sooner in order to make room for</span>
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	<span class="comment">// other memory sources. Note that since we&#39;re effectively reducing the heap goal by X bytes,</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	<span class="comment">// we&#39;re actually giving more than X bytes of headroom back, because the heap goal is in</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	<span class="comment">// terms of heap objects, but it takes more than X bytes (e.g. due to fragmentation) to store</span>
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	<span class="comment">// X bytes worth of objects.</span>
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	<span class="comment">// The final adjustment (marker 3) reduces the maximum possible memory limit heap goal by</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	<span class="comment">// memoryLimitHeapGoalPercent. As the name implies, this is to provide additional headroom in</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>	<span class="comment">// the face of pacing inaccuracies, and also to leave a buffer of unscavenged memory so the</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>	<span class="comment">// allocator isn&#39;t constantly scavenging. The reduction amount also has a fixed minimum</span>
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	<span class="comment">// (memoryLimitMinHeapGoalHeadroom, not pictured) because the aforementioned pacing inaccuracies</span>
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	<span class="comment">// disproportionately affect small heaps: as heaps get smaller, the pacer&#39;s inputs get fuzzier.</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	<span class="comment">// Shorter GC cycles and less GC work means noisy external factors like the OS scheduler have a</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>	<span class="comment">// greater impact.</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	memoryLimit := uint64(c.memoryLimit.Load())
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	<span class="comment">// Compute term 1.</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	nonHeapMemory := mappedReady - heapFree - heapAlloc
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	<span class="comment">// Compute term 2.</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	var overage uint64
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>	if mappedReady &gt; memoryLimit {
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>		overage = mappedReady - memoryLimit
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	}
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>	if nonHeapMemory+overage &gt;= memoryLimit {
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re at a point where non-heap memory exceeds the memory limit on its own.</span>
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>		<span class="comment">// There&#39;s honestly not much we can do here but just trigger GCs continuously</span>
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>		<span class="comment">// and let the CPU limiter reign that in. Something has to give at this point.</span>
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>		<span class="comment">// Set it to heapMarked, the lowest possible goal.</span>
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		return c.heapMarked
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>	}
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>	<span class="comment">// Compute the goal.</span>
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>	goal := memoryLimit - (nonHeapMemory + overage)
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>	<span class="comment">// Apply some headroom to the goal to account for pacing inaccuracies and to reduce</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>	<span class="comment">// the impact of scavenging at allocation time in response to a high allocation rate</span>
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>	<span class="comment">// when GOGC=off. See issue #57069. Also, be careful about small limits.</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>	headroom := goal / 100 * memoryLimitHeapGoalHeadroomPercent
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>	if headroom &lt; memoryLimitMinHeapGoalHeadroom {
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>		<span class="comment">// Set a fixed minimum to deal with the particularly large effect pacing inaccuracies</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>		<span class="comment">// have for smaller heaps.</span>
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>		headroom = memoryLimitMinHeapGoalHeadroom
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>	}
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>	if goal &lt; headroom || goal-headroom &lt; headroom {
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>		goal = headroom
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>	} else {
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>		goal = goal - headroom
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>	}
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t let us go below the live heap. A heap goal below the live heap doesn&#39;t make sense.</span>
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>	if goal &lt; c.heapMarked {
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>		goal = c.heapMarked
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>	}
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	return goal
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>}
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>const (
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	<span class="comment">// These constants determine the bounds on the GC trigger as a fraction</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	<span class="comment">// of heap bytes allocated between the start of a GC (heapLive == heapMarked)</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>	<span class="comment">// and the end of a GC (heapLive == heapGoal).</span>
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	<span class="comment">// The constants are obscured in this way for efficiency. The denominator</span>
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	<span class="comment">// of the fraction is always a power-of-two for a quick division, so that</span>
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>	<span class="comment">// the numerator is a single constant integer multiplication.</span>
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	triggerRatioDen = 64
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	<span class="comment">// The minimum trigger constant was chosen empirically: given a sufficiently</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>	<span class="comment">// fast/scalable allocator with 48 Ps that could drive the trigger ratio</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>	<span class="comment">// to &lt;0.05, this constant causes applications to retain the same peak</span>
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>	<span class="comment">// RSS compared to not having this allocator.</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>	minTriggerRatioNum = 45 <span class="comment">// ~0.7</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>	<span class="comment">// The maximum trigger constant is chosen somewhat arbitrarily, but the</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>	<span class="comment">// current constant has served us well over the years.</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	maxTriggerRatioNum = 61 <span class="comment">// ~0.95</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>)
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span><span class="comment">// trigger returns the current point at which a GC should trigger along with</span>
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span><span class="comment">// the heap goal.</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span><span class="comment">// The returned value may be compared against heapLive to determine whether</span>
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span><span class="comment">// the GC should trigger. Thus, the GC trigger condition should be (but may</span>
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span><span class="comment">// not be, in the case of small movements for efficiency) checked whenever</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span><span class="comment">// the heap goal may change.</span>
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>func (c *gcControllerState) trigger() (uint64, uint64) {
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>	goal, minTrigger := c.heapGoalInternal()
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>	<span class="comment">// Invariant: the trigger must always be less than the heap goal.</span>
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>	<span class="comment">// Note that the memory limit sets a hard maximum on our heap goal,</span>
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>	<span class="comment">// but the live heap may grow beyond it.</span>
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>	if c.heapMarked &gt;= goal {
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>		<span class="comment">// The goal should never be smaller than heapMarked, but let&#39;s be</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>		<span class="comment">// defensive about it. The only reasonable trigger here is one that</span>
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>		<span class="comment">// causes a continuous GC cycle at heapMarked, but respect the goal</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>		<span class="comment">// if it came out as smaller than that.</span>
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>		return goal, goal
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>	}
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>	<span class="comment">// Below this point, c.heapMarked &lt; goal.</span>
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>	<span class="comment">// heapMarked is our absolute minimum, and it&#39;s possible the trigger</span>
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>	<span class="comment">// bound we get from heapGoalinternal is less than that.</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>	if minTrigger &lt; c.heapMarked {
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>		minTrigger = c.heapMarked
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>	}
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>	<span class="comment">// If we let the trigger go too low, then if the application</span>
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	<span class="comment">// is allocating very rapidly we might end up in a situation</span>
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	<span class="comment">// where we&#39;re allocating black during a nearly always-on GC.</span>
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>	<span class="comment">// The result of this is a growing heap and ultimately an</span>
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>	<span class="comment">// increase in RSS. By capping us at a point &gt;0, we&#39;re essentially</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>	<span class="comment">// saying that we&#39;re OK using more CPU during the GC to prevent</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>	<span class="comment">// this growth in RSS.</span>
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>	triggerLowerBound := ((goal-c.heapMarked)/triggerRatioDen)*minTriggerRatioNum + c.heapMarked
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>	if minTrigger &lt; triggerLowerBound {
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>		minTrigger = triggerLowerBound
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>	}
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>	<span class="comment">// For small heaps, set the max trigger point at maxTriggerRatio of the way</span>
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>	<span class="comment">// from the live heap to the heap goal. This ensures we always have *some*</span>
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>	<span class="comment">// headroom when the GC actually starts. For larger heaps, set the max trigger</span>
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>	<span class="comment">// point at the goal, minus the minimum heap size.</span>
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>	<span class="comment">// This choice follows from the fact that the minimum heap size is chosen</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>	<span class="comment">// to reflect the costs of a GC with no work to do. With a large heap but</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>	<span class="comment">// very little scan work to perform, this gives us exactly as much runway</span>
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	<span class="comment">// as we would need, in the worst case.</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	maxTrigger := ((goal-c.heapMarked)/triggerRatioDen)*maxTriggerRatioNum + c.heapMarked
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	if goal &gt; defaultHeapMinimum &amp;&amp; goal-defaultHeapMinimum &gt; maxTrigger {
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>		maxTrigger = goal - defaultHeapMinimum
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>	}
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>	maxTrigger = max(maxTrigger, minTrigger)
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>	<span class="comment">// Compute the trigger from our bounds and the runway stored by commit.</span>
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>	var trigger uint64
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>	runway := c.runway.Load()
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>	if runway &gt; goal {
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>		trigger = minTrigger
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>	} else {
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>		trigger = goal - runway
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>	}
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>	trigger = max(trigger, minTrigger)
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>	trigger = min(trigger, maxTrigger)
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>	if trigger &gt; goal {
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>		print(&#34;trigger=&#34;, trigger, &#34; heapGoal=&#34;, goal, &#34;\n&#34;)
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>		print(&#34;minTrigger=&#34;, minTrigger, &#34; maxTrigger=&#34;, maxTrigger, &#34;\n&#34;)
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>		throw(&#34;produced a trigger greater than the heap goal&#34;)
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>	}
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>	return trigger, goal
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>}
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span><span class="comment">// commit recomputes all pacing parameters needed to derive the</span>
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span><span class="comment">// trigger and the heap goal. Namely, the gcPercent-based heap goal,</span>
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span><span class="comment">// and the amount of runway we want to give the GC this cycle.</span>
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span><span class="comment">// This can be called any time. If GC is the in the middle of a</span>
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span><span class="comment">// concurrent phase, it will adjust the pacing of that phase.</span>
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span><span class="comment">// isSweepDone should be the result of calling isSweepDone(),</span>
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span><span class="comment">// unless we&#39;re testing or we know we&#39;re executing during a GC cycle.</span>
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span><span class="comment">// This depends on gcPercent, gcController.heapMarked, and</span>
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span><span class="comment">// gcController.heapLive. These must be up to date.</span>
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span><span class="comment">// Callers must call gcControllerState.revise after calling this</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span><span class="comment">// function if the GC is enabled.</span>
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span><span class="comment">// mheap_.lock must be held or the world must be stopped.</span>
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>func (c *gcControllerState) commit(isSweepDone bool) {
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>	if !c.test {
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>		assertWorldStoppedOrLockHeld(&amp;mheap_.lock)
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>	}
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>	if isSweepDone {
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>		<span class="comment">// The sweep is done, so there aren&#39;t any restrictions on the trigger</span>
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>		<span class="comment">// we need to think about.</span>
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>		c.sweepDistMinTrigger.Store(0)
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>	} else {
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>		<span class="comment">// Concurrent sweep happens in the heap growth</span>
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>		<span class="comment">// from gcController.heapLive to trigger. Make sure we</span>
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>		<span class="comment">// give the sweeper some runway if it doesn&#39;t have enough.</span>
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>		c.sweepDistMinTrigger.Store(c.heapLive.Load() + sweepMinHeapDistance)
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>	}
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>	<span class="comment">// Compute the next GC goal, which is when the allocated heap</span>
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>	<span class="comment">// has grown by GOGC/100 over where it started the last cycle,</span>
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>	<span class="comment">// plus additional runway for non-heap sources of GC work.</span>
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>	gcPercentHeapGoal := ^uint64(0)
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>	if gcPercent := c.gcPercent.Load(); gcPercent &gt;= 0 {
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>		gcPercentHeapGoal = c.heapMarked + (c.heapMarked+c.lastStackScan.Load()+c.globalsScan.Load())*uint64(gcPercent)/100
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>	}
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>	<span class="comment">// Apply the minimum heap size here. It&#39;s defined in terms of gcPercent</span>
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>	<span class="comment">// and is only updated by functions that call commit.</span>
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>	if gcPercentHeapGoal &lt; c.heapMinimum {
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>		gcPercentHeapGoal = c.heapMinimum
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>	}
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>	c.gcPercentHeapGoal.Store(gcPercentHeapGoal)
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>	<span class="comment">// Compute the amount of runway we want the GC to have by using our</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>	<span class="comment">// estimate of the cons/mark ratio.</span>
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>	<span class="comment">// The idea is to take our expected scan work, and multiply it by</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>	<span class="comment">// the cons/mark ratio to determine how long it&#39;ll take to complete</span>
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>	<span class="comment">// that scan work in terms of bytes allocated. This gives us our GC&#39;s</span>
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>	<span class="comment">// runway.</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>	<span class="comment">// However, the cons/mark ratio is a ratio of rates per CPU-second, but</span>
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>	<span class="comment">// here we care about the relative rates for some division of CPU</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>	<span class="comment">// resources among the mutator and the GC.</span>
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>	<span class="comment">// To summarize, we have B / cpu-ns, and we want B / ns. We get that</span>
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>	<span class="comment">// by multiplying by our desired division of CPU resources. We choose</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>	<span class="comment">// to express CPU resources as GOMAPROCS*fraction. Note that because</span>
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>	<span class="comment">// we&#39;re working with a ratio here, we can omit the number of CPU cores,</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>	<span class="comment">// because they&#39;ll appear in the numerator and denominator and cancel out.</span>
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>	<span class="comment">// As a result, this is basically just &#34;weighing&#34; the cons/mark ratio by</span>
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	<span class="comment">// our desired division of resources.</span>
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>	<span class="comment">// Furthermore, by setting the runway so that CPU resources are divided</span>
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>	<span class="comment">// this way, assuming that the cons/mark ratio is correct, we make that</span>
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>	<span class="comment">// division a reality.</span>
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>	c.runway.Store(uint64((c.consMark * (1 - gcGoalUtilization) / (gcGoalUtilization)) * float64(c.lastHeapScan+c.lastStackScan.Load()+c.globalsScan.Load())))
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>}
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span><span class="comment">// setGCPercent updates gcPercent. commit must be called after.</span>
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span><span class="comment">// Returns the old value of gcPercent.</span>
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span><span class="comment">// The world must be stopped, or mheap_.lock must be held.</span>
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>func (c *gcControllerState) setGCPercent(in int32) int32 {
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>	if !c.test {
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>		assertWorldStoppedOrLockHeld(&amp;mheap_.lock)
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>	}
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>	out := c.gcPercent.Load()
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>	if in &lt; 0 {
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>		in = -1
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>	}
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>	c.heapMinimum = defaultHeapMinimum * uint64(in) / 100
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>	c.gcPercent.Store(in)
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>	return out
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>}
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span><span class="comment">//go:linkname setGCPercent runtime/debug.setGCPercent</span>
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>func setGCPercent(in int32) (out int32) {
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>	<span class="comment">// Run on the system stack since we grab the heap lock.</span>
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>		lock(&amp;mheap_.lock)
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>		out = gcController.setGCPercent(in)
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>		gcControllerCommit()
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>		unlock(&amp;mheap_.lock)
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>	})
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>	<span class="comment">// If we just disabled GC, wait for any concurrent GC mark to</span>
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>	<span class="comment">// finish so we always return with no GC running.</span>
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	if in &lt; 0 {
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>		gcWaitOnMark(work.cycles.Load())
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>	}
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>	return out
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>}
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>func readGOGC() int32 {
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>	p := gogetenv(&#34;GOGC&#34;)
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>	if p == &#34;off&#34; {
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>		return -1
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>	}
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>	if n, ok := atoi32(p); ok {
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>		return n
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>	}
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>	return 100
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>}
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span><span class="comment">// setMemoryLimit updates memoryLimit. commit must be called after</span>
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span><span class="comment">// Returns the old value of memoryLimit.</span>
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span><span class="comment">// The world must be stopped, or mheap_.lock must be held.</span>
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>func (c *gcControllerState) setMemoryLimit(in int64) int64 {
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>	if !c.test {
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>		assertWorldStoppedOrLockHeld(&amp;mheap_.lock)
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>	}
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>	out := c.memoryLimit.Load()
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>	if in &gt;= 0 {
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>		c.memoryLimit.Store(in)
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>	}
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>	return out
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>}
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span><span class="comment">//go:linkname setMemoryLimit runtime/debug.setMemoryLimit</span>
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>func setMemoryLimit(in int64) (out int64) {
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>	<span class="comment">// Run on the system stack since we grab the heap lock.</span>
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>		lock(&amp;mheap_.lock)
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>		out = gcController.setMemoryLimit(in)
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>		if in &lt; 0 || out == in {
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>			<span class="comment">// If we&#39;re just checking the value or not changing</span>
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>			<span class="comment">// it, there&#39;s no point in doing the rest.</span>
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>			unlock(&amp;mheap_.lock)
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>			return
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>		}
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>		gcControllerCommit()
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>		unlock(&amp;mheap_.lock)
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>	})
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>	return out
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>}
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>func readGOMEMLIMIT() int64 {
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>	p := gogetenv(&#34;GOMEMLIMIT&#34;)
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>	if p == &#34;&#34; || p == &#34;off&#34; {
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>		return maxInt64
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>	}
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>	n, ok := parseByteCount(p)
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>	if !ok {
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>		print(&#34;GOMEMLIMIT=&#34;, p, &#34;\n&#34;)
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>		throw(&#34;malformed GOMEMLIMIT; see `go doc runtime/debug.SetMemoryLimit`&#34;)
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>	}
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>	return n
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>}
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span><span class="comment">// addIdleMarkWorker attempts to add a new idle mark worker.</span>
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span><span class="comment">// If this returns true, the caller must become an idle mark worker unless</span>
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span><span class="comment">// there&#39;s no background mark worker goroutines in the pool. This case is</span>
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span><span class="comment">// harmless because there are already background mark workers running.</span>
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span><span class="comment">// If this returns false, the caller must NOT become an idle mark worker.</span>
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span><span class="comment">// nosplit because it may be called without a P.</span>
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>func (c *gcControllerState) addIdleMarkWorker() bool {
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>	for {
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>		old := c.idleMarkWorkers.Load()
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>		n, max := int32(old&amp;uint64(^uint32(0))), int32(old&gt;&gt;32)
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>		if n &gt;= max {
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>			<span class="comment">// See the comment on idleMarkWorkers for why</span>
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>			<span class="comment">// n &gt; max is tolerated.</span>
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>			return false
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>		}
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>		if n &lt; 0 {
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>			print(&#34;n=&#34;, n, &#34; max=&#34;, max, &#34;\n&#34;)
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>			throw(&#34;negative idle mark workers&#34;)
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>		}
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>		new := uint64(uint32(n+1)) | (uint64(max) &lt;&lt; 32)
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>		if c.idleMarkWorkers.CompareAndSwap(old, new) {
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>			return true
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>		}
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>	}
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>}
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span><span class="comment">// needIdleMarkWorker is a hint as to whether another idle mark worker is needed.</span>
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span><span class="comment">// The caller must still call addIdleMarkWorker to become one. This is mainly</span>
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span><span class="comment">// useful for a quick check before an expensive operation.</span>
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span><span class="comment">// nosplit because it may be called without a P.</span>
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>func (c *gcControllerState) needIdleMarkWorker() bool {
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>	p := c.idleMarkWorkers.Load()
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>	n, max := int32(p&amp;uint64(^uint32(0))), int32(p&gt;&gt;32)
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>	return n &lt; max
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>}
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span><span class="comment">// removeIdleMarkWorker must be called when a new idle mark worker stops executing.</span>
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>func (c *gcControllerState) removeIdleMarkWorker() {
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>	for {
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>		old := c.idleMarkWorkers.Load()
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>		n, max := int32(old&amp;uint64(^uint32(0))), int32(old&gt;&gt;32)
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>		if n-1 &lt; 0 {
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>			print(&#34;n=&#34;, n, &#34; max=&#34;, max, &#34;\n&#34;)
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>			throw(&#34;negative idle mark workers&#34;)
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>		}
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>		new := uint64(uint32(n-1)) | (uint64(max) &lt;&lt; 32)
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>		if c.idleMarkWorkers.CompareAndSwap(old, new) {
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>			return
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>		}
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>	}
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>}
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span><span class="comment">// setMaxIdleMarkWorkers sets the maximum number of idle mark workers allowed.</span>
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span><span class="comment">// This method is optimistic in that it does not wait for the number of</span>
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span><span class="comment">// idle mark workers to reduce to max before returning; it assumes the workers</span>
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span><span class="comment">// will deschedule themselves.</span>
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>func (c *gcControllerState) setMaxIdleMarkWorkers(max int32) {
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>	for {
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>		old := c.idleMarkWorkers.Load()
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>		n := int32(old &amp; uint64(^uint32(0)))
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>		if n &lt; 0 {
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>			print(&#34;n=&#34;, n, &#34; max=&#34;, max, &#34;\n&#34;)
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>			throw(&#34;negative idle mark workers&#34;)
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>		}
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>		new := uint64(uint32(n)) | (uint64(max) &lt;&lt; 32)
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>		if c.idleMarkWorkers.CompareAndSwap(old, new) {
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>			return
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>		}
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>	}
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>}
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span><span class="comment">// gcControllerCommit is gcController.commit, but passes arguments from live</span>
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span><span class="comment">// (non-test) data. It also updates any consumers of the GC pacing, such as</span>
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span><span class="comment">// sweep pacing and the background scavenger.</span>
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span><span class="comment">// Calls gcController.commit.</span>
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span><span class="comment">// The heap lock must be held, so this must be executed on the system stack.</span>
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>func gcControllerCommit() {
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>	assertWorldStoppedOrLockHeld(&amp;mheap_.lock)
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>	gcController.commit(isSweepDone())
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span>	<span class="comment">// Update mark pacing.</span>
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span>	if gcphase != _GCoff {
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>		gcController.revise()
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>	}
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): This isn&#39;t really accurate any longer because the heap</span>
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span>	<span class="comment">// goal is computed dynamically. Still useful to snapshot, but not as useful.</span>
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span>	trace := traceAcquire()
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span>		trace.HeapGoal()
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>	}
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>	trigger, heapGoal := gcController.trigger()
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>	gcPaceSweeper(trigger)
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>	gcPaceScavenger(gcController.memoryLimit.Load(), heapGoal, gcController.lastHeapGoal)
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>}
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>
</pre><p><a href="mgcpacer.go?m=text">View as plain text</a></p>

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
