<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/metrics/doc.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../index.html">GoDoc</a></div>
<a href="doc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<a href="http://localhost:8080/src/runtime/metrics">metrics</a>/<span class="text-muted">doc.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime/metrics">runtime/metrics</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2020 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Note: run &#39;go generate&#39; (which will run &#39;go test -generate&#39;) to update the &#34;Supported metrics&#34; list.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//go:generate go test -run=Docs -generate</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">/*
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>Package metrics provides a stable interface to access implementation-defined
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>metrics exported by the Go runtime. This package is similar to existing functions
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>like [runtime.ReadMemStats] and [runtime/debug.ReadGCStats], but significantly more general.
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>The set of metrics defined by this package may evolve as the runtime itself
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>evolves, and also enables variation across Go implementations, whose relevant
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>metric sets may not intersect.
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span># Interface
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>Metrics are designated by a string key, rather than, for example, a field name in
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>a struct. The full list of supported metrics is always available in the slice of
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>Descriptions returned by [All]. Each [Description] also includes useful information
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>about the metric.
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>Thus, users of this API are encouraged to sample supported metrics defined by the
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>slice returned by All to remain compatible across Go versions. Of course, situations
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>arise where reading specific metrics is critical. For these cases, users are
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>encouraged to use build tags, and although metrics may be deprecated and removed,
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>users should consider this to be an exceptional and rare event, coinciding with a
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>very large change in a particular Go implementation.
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>Each metric key also has a &#34;kind&#34; (see [ValueKind]) that describes the format of the
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>metric&#39;s value.
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>In the interest of not breaking users of this package, the &#34;kind&#34; for a given metric
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>is guaranteed not to change. If it must change, then a new metric will be introduced
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>with a new key and a new &#34;kind.&#34;
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span># Metric key format
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>As mentioned earlier, metric keys are strings. Their format is simple and well-defined,
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>designed to be both human and machine readable. It is split into two components,
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>separated by a colon: a rooted path and a unit. The choice to include the unit in
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>the key is motivated by compatibility: if a metric&#39;s unit changes, its semantics likely
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>did also, and a new key should be introduced.
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>For more details on the precise definition of the metric key&#39;s path and unit formats, see
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>the documentation of the Name field of the Description struct.
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span># A note about floats
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>This package supports metrics whose values have a floating-point representation. In
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>order to improve ease-of-use, this package promises to never produce the following
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>classes of floating-point values: NaN, infinity.
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span># Supported metrics
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>Below is the full list of supported metrics, ordered lexicographically.
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	/cgo/go-to-c-calls:calls
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		Count of calls made from Go to C by the current process.
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	/cpu/classes/gc/mark/assist:cpu-seconds
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		Estimated total CPU time goroutines spent performing GC
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		tasks to assist the GC and prevent it from falling behind the
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		application. This metric is an overestimate, and not directly
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		comparable to system CPU time measurements. Compare only with
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		other /cpu/classes metrics.
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	/cpu/classes/gc/mark/dedicated:cpu-seconds
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		Estimated total CPU time spent performing GC tasks on processors
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		(as defined by GOMAXPROCS) dedicated to those tasks. This metric
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		is an overestimate, and not directly comparable to system CPU
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		time measurements. Compare only with other /cpu/classes metrics.
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	/cpu/classes/gc/mark/idle:cpu-seconds
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		Estimated total CPU time spent performing GC tasks on spare CPU
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		resources that the Go scheduler could not otherwise find a use
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		for. This should be subtracted from the total GC CPU time to
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		obtain a measure of compulsory GC CPU time. This metric is an
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		overestimate, and not directly comparable to system CPU time
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		measurements. Compare only with other /cpu/classes metrics.
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	/cpu/classes/gc/pause:cpu-seconds
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		Estimated total CPU time spent with the application paused by
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		the GC. Even if only one thread is running during the pause,
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		this is computed as GOMAXPROCS times the pause latency because
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		nothing else can be executing. This is the exact sum of samples
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		in /sched/pauses/total/gc:seconds if each sample is multiplied
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		by GOMAXPROCS at the time it is taken. This metric is an
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		overestimate, and not directly comparable to system CPU time
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		measurements. Compare only with other /cpu/classes metrics.
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	/cpu/classes/gc/total:cpu-seconds
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		Estimated total CPU time spent performing GC tasks. This metric
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		is an overestimate, and not directly comparable to system CPU
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		time measurements. Compare only with other /cpu/classes metrics.
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		Sum of all metrics in /cpu/classes/gc.
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	/cpu/classes/idle:cpu-seconds
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		Estimated total available CPU time not spent executing
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		any Go or Go runtime code. In other words, the part of
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		/cpu/classes/total:cpu-seconds that was unused. This metric is
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		an overestimate, and not directly comparable to system CPU time
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		measurements. Compare only with other /cpu/classes metrics.
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	/cpu/classes/scavenge/assist:cpu-seconds
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		Estimated total CPU time spent returning unused memory to the
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		underlying platform in response eagerly in response to memory
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		pressure. This metric is an overestimate, and not directly
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		comparable to system CPU time measurements. Compare only with
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		other /cpu/classes metrics.
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	/cpu/classes/scavenge/background:cpu-seconds
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		Estimated total CPU time spent performing background tasks to
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		return unused memory to the underlying platform. This metric is
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		an overestimate, and not directly comparable to system CPU time
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		measurements. Compare only with other /cpu/classes metrics.
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	/cpu/classes/scavenge/total:cpu-seconds
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		Estimated total CPU time spent performing tasks that return
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		unused memory to the underlying platform. This metric is an
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		overestimate, and not directly comparable to system CPU time
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		measurements. Compare only with other /cpu/classes metrics.
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		Sum of all metrics in /cpu/classes/scavenge.
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	/cpu/classes/total:cpu-seconds
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		Estimated total available CPU time for user Go code or the Go
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		runtime, as defined by GOMAXPROCS. In other words, GOMAXPROCS
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		integrated over the wall-clock duration this process has been
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		executing for. This metric is an overestimate, and not directly
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		comparable to system CPU time measurements. Compare only with
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		other /cpu/classes metrics. Sum of all metrics in /cpu/classes.
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	/cpu/classes/user:cpu-seconds
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		Estimated total CPU time spent running user Go code. This may
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		also include some small amount of time spent in the Go runtime.
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		This metric is an overestimate, and not directly comparable
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		to system CPU time measurements. Compare only with other
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		/cpu/classes metrics.
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	/gc/cycles/automatic:gc-cycles
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		Count of completed GC cycles generated by the Go runtime.
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	/gc/cycles/forced:gc-cycles
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		Count of completed GC cycles forced by the application.
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	/gc/cycles/total:gc-cycles
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		Count of all completed GC cycles.
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	/gc/gogc:percent
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		Heap size target percentage configured by the user, otherwise
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		100. This value is set by the GOGC environment variable, and the
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		runtime/debug.SetGCPercent function.
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	/gc/gomemlimit:bytes
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		Go runtime memory limit configured by the user, otherwise
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		math.MaxInt64. This value is set by the GOMEMLIMIT environment
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		variable, and the runtime/debug.SetMemoryLimit function.
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	/gc/heap/allocs-by-size:bytes
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		Distribution of heap allocations by approximate size.
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		Bucket counts increase monotonically. Note that this does not
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		include tiny objects as defined by /gc/heap/tiny/allocs:objects,
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		only tiny blocks.
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	/gc/heap/allocs:bytes
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		Cumulative sum of memory allocated to the heap by the
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		application.
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	/gc/heap/allocs:objects
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		Cumulative count of heap allocations triggered by the
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		application. Note that this does not include tiny objects as
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	/gc/heap/frees-by-size:bytes
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		Distribution of freed heap allocations by approximate size.
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		Bucket counts increase monotonically. Note that this does not
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		include tiny objects as defined by /gc/heap/tiny/allocs:objects,
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		only tiny blocks.
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	/gc/heap/frees:bytes
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		Cumulative sum of heap memory freed by the garbage collector.
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	/gc/heap/frees:objects
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		Cumulative count of heap allocations whose storage was freed
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		by the garbage collector. Note that this does not include tiny
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		objects as defined by /gc/heap/tiny/allocs:objects, only tiny
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		blocks.
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	/gc/heap/goal:bytes
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		Heap size target for the end of the GC cycle.
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	/gc/heap/live:bytes
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		Heap memory occupied by live objects that were marked by the
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		previous GC.
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	/gc/heap/objects:objects
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		Number of objects, live or unswept, occupying heap memory.
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	/gc/heap/tiny/allocs:objects
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		Count of small allocations that are packed together into blocks.
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		These allocations are counted separately from other allocations
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		because each individual allocation is not tracked by the
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		runtime, only their block. Each block is already accounted for
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		in allocs-by-size and frees-by-size.
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	/gc/limiter/last-enabled:gc-cycle
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		GC cycle the last time the GC CPU limiter was enabled.
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		This metric is useful for diagnosing the root cause of an
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		out-of-memory error, because the limiter trades memory for CPU
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		time when the GC&#39;s CPU time gets too high. This is most likely
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		to occur with use of SetMemoryLimit. The first GC cycle is cycle
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		1, so a value of 0 indicates that it was never enabled.
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	/gc/pauses:seconds
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		Deprecated. Prefer the identical /sched/pauses/total/gc:seconds.
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	/gc/scan/globals:bytes
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		The total amount of global variable space that is scannable.
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	/gc/scan/heap:bytes
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		The total amount of heap space that is scannable.
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	/gc/scan/stack:bytes
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		The number of bytes of stack that were scanned last GC cycle.
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	/gc/scan/total:bytes
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		The total amount space that is scannable. Sum of all metrics in
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		/gc/scan.
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	/gc/stack/starting-size:bytes
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		The stack size of new goroutines.
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	/godebug/non-default-behavior/execerrdot:events
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the os/exec
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=execerrdot=... setting.
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	/godebug/non-default-behavior/gocachehash:events
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the cmd/go
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=gocachehash=... setting.
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	/godebug/non-default-behavior/gocachetest:events
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the cmd/go
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=gocachetest=... setting.
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	/godebug/non-default-behavior/gocacheverify:events
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the cmd/go
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=gocacheverify=... setting.
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	/godebug/non-default-behavior/gotypesalias:events
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the go/types
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=gotypesalias=... setting.
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	/godebug/non-default-behavior/http2client:events
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the net/http
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=http2client=... setting.
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	/godebug/non-default-behavior/http2server:events
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the net/http
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=http2server=... setting.
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	/godebug/non-default-behavior/httplaxcontentlength:events
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the net/http
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=httplaxcontentlength=...
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		setting.
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	/godebug/non-default-behavior/httpmuxgo121:events
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the net/http
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=httpmuxgo121=... setting.
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	/godebug/non-default-behavior/installgoroot:events
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the go/build
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=installgoroot=... setting.
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	/godebug/non-default-behavior/jstmpllitinterp:events
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		The number of non-default behaviors executed by
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		the html/template package due to a non-default
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		GODEBUG=jstmpllitinterp=... setting.
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	/godebug/non-default-behavior/multipartmaxheaders:events
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		The number of non-default behaviors executed by
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		the mime/multipart package due to a non-default
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		GODEBUG=multipartmaxheaders=... setting.
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	/godebug/non-default-behavior/multipartmaxparts:events
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		The number of non-default behaviors executed by
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		the mime/multipart package due to a non-default
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		GODEBUG=multipartmaxparts=... setting.
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	/godebug/non-default-behavior/multipathtcp:events
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the net package
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		due to a non-default GODEBUG=multipathtcp=... setting.
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	/godebug/non-default-behavior/panicnil:events
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the runtime
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=panicnil=... setting.
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	/godebug/non-default-behavior/randautoseed:events
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the math/rand
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=randautoseed=... setting.
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	/godebug/non-default-behavior/tarinsecurepath:events
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the archive/tar
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=tarinsecurepath=...
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		setting.
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	/godebug/non-default-behavior/tls10server:events
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the crypto/tls
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=tls10server=... setting.
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	/godebug/non-default-behavior/tlsmaxrsasize:events
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the crypto/tls
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=tlsmaxrsasize=... setting.
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	/godebug/non-default-behavior/tlsrsakex:events
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the crypto/tls
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=tlsrsakex=... setting.
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	/godebug/non-default-behavior/tlsunsafeekm:events
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the crypto/tls
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=tlsunsafeekm=... setting.
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	/godebug/non-default-behavior/x509sha1:events
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the crypto/x509
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=x509sha1=... setting.
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	/godebug/non-default-behavior/x509usefallbackroots:events
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the crypto/x509
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=x509usefallbackroots=...
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		setting.
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	/godebug/non-default-behavior/x509usepolicies:events
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the crypto/x509
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=x509usepolicies=...
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		setting.
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	/godebug/non-default-behavior/zipinsecurepath:events
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		The number of non-default behaviors executed by the archive/zip
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		package due to a non-default GODEBUG=zipinsecurepath=...
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		setting.
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	/memory/classes/heap/free:bytes
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		Memory that is completely free and eligible to be returned to
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		the underlying system, but has not been. This metric is the
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		runtime&#39;s estimate of free address space that is backed by
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		physical memory.
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	/memory/classes/heap/objects:bytes
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		Memory occupied by live objects and dead objects that have not
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		yet been marked free by the garbage collector.
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	/memory/classes/heap/released:bytes
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		Memory that is completely free and has been returned to the
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		underlying system. This metric is the runtime&#39;s estimate of free
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		address space that is still mapped into the process, but is not
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		backed by physical memory.
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	/memory/classes/heap/stacks:bytes
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		Memory allocated from the heap that is reserved for stack space,
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		whether or not it is currently in-use. Currently, this
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		represents all stack memory for goroutines. It also includes all
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		OS thread stacks in non-cgo programs. Note that stacks may be
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		allocated differently in the future, and this may change.
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	/memory/classes/heap/unused:bytes
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		Memory that is reserved for heap objects but is not currently
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		used to hold heap objects.
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	/memory/classes/metadata/mcache/free:bytes
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		Memory that is reserved for runtime mcache structures, but not
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		in-use.
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	/memory/classes/metadata/mcache/inuse:bytes
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		Memory that is occupied by runtime mcache structures that are
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		currently being used.
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	/memory/classes/metadata/mspan/free:bytes
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		Memory that is reserved for runtime mspan structures, but not
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		in-use.
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	/memory/classes/metadata/mspan/inuse:bytes
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		Memory that is occupied by runtime mspan structures that are
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		currently being used.
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	/memory/classes/metadata/other:bytes
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		Memory that is reserved for or used to hold runtime metadata.
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	/memory/classes/os-stacks:bytes
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		Stack memory allocated by the underlying operating system.
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		In non-cgo programs this metric is currently zero. This may
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		change in the future.In cgo programs this metric includes
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		OS thread stacks allocated directly from the OS. Currently,
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		this only accounts for one stack in c-shared and c-archive build
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		modes, and other sources of stacks from the OS are not measured.
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		This too may change in the future.
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	/memory/classes/other:bytes
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		Memory used by execution trace buffers, structures for debugging
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		the runtime, finalizer and profiler specials, and more.
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	/memory/classes/profiling/buckets:bytes
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		Memory that is used by the stack trace hash map used for
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		profiling.
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	/memory/classes/total:bytes
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		All memory mapped by the Go runtime into the current process
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		as read-write. Note that this does not include memory mapped
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		by code called via cgo or via the syscall package. Sum of all
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		metrics in /memory/classes.
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	/sched/gomaxprocs:threads
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		The current runtime.GOMAXPROCS setting, or the number of
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		operating system threads that can execute user-level Go code
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		simultaneously.
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	/sched/goroutines:goroutines
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		Count of live goroutines.
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	/sched/latencies:seconds
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		Distribution of the time goroutines have spent in the scheduler
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		in a runnable state before actually running. Bucket counts
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		increase monotonically.
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	/sched/pauses/stopping/gc:seconds
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		Distribution of individual GC-related stop-the-world stopping
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		latencies. This is the time it takes from deciding to stop the
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		world until all Ps are stopped. This is a subset of the total
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		GC-related stop-the-world time (/sched/pauses/total/gc:seconds).
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		During this time, some threads may be executing. Bucket counts
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		increase monotonically.
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	/sched/pauses/stopping/other:seconds
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		Distribution of individual non-GC-related stop-the-world
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		stopping latencies. This is the time it takes from deciding
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		to stop the world until all Ps are stopped. This is a
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		subset of the total non-GC-related stop-the-world time
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		(/sched/pauses/total/other:seconds). During this time, some
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		threads may be executing. Bucket counts increase monotonically.
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	/sched/pauses/total/gc:seconds
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		Distribution of individual GC-related stop-the-world pause
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		latencies. This is the time from deciding to stop the world
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		until the world is started again. Some of this time is spent
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		getting all threads to stop (this is measured directly in
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		/sched/pauses/stopping/gc:seconds), during which some threads
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		may still be running. Bucket counts increase monotonically.
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	/sched/pauses/total/other:seconds
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		Distribution of individual non-GC-related stop-the-world
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		pause latencies. This is the time from deciding to stop the
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		world until the world is started again. Some of this time
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		is spent getting all threads to stop (measured directly in
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		/sched/pauses/stopping/other:seconds). Bucket counts increase
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		monotonically.
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	/sync/mutex/wait/total:seconds
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		Approximate cumulative time goroutines have spent blocked on a
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		sync.Mutex, sync.RWMutex, or runtime-internal lock. This metric
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		is useful for identifying global changes in lock contention.
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		Collect a mutex or block profile using the runtime/pprof package
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		for more detailed contention data.
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>*/</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>package metrics
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>
</pre><p><a href="doc.go?m=text">View as plain text</a></p>

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
