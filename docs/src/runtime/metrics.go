<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/metrics.go - Go Documentation Server</title>

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
<a href="metrics.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">metrics.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2020 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Metrics implementation exported to runtime/metrics.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/godebugs&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>var (
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	<span class="comment">// metrics is a map of runtime/metrics keys to data used by the runtime</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	<span class="comment">// to sample each metric&#39;s value. metricsInit indicates it has been</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">// initialized.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// These fields are protected by metricsSema which should be</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// locked/unlocked with metricsLock() / metricsUnlock().</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	metricsSema uint32 = 1
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	metricsInit bool
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	metrics     map[string]metricData
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	sizeClassBuckets []float64
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	timeHistBuckets  []float64
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>type metricData struct {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// deps is the set of runtime statistics that this metric</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// depends on. Before compute is called, the statAggregate</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// which will be passed must ensure() these dependencies.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	deps statDepSet
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// compute is a function that populates a metricValue</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// given a populated statAggregate structure.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	compute func(in *statAggregate, out *metricValue)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>func metricsLock() {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// Acquire the metricsSema but with handoff. Operations are typically</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// expensive enough that queueing up goroutines and handing off between</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// them will be noticeably better-behaved.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	semacquire1(&amp;metricsSema, true, 0, 0, waitReasonSemacquire)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	if raceenabled {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		raceacquire(unsafe.Pointer(&amp;metricsSema))
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>func metricsUnlock() {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	if raceenabled {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		racerelease(unsafe.Pointer(&amp;metricsSema))
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	semrelease(&amp;metricsSema)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// initMetrics initializes the metrics map if it hasn&#39;t been yet.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// metricsSema must be held.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func initMetrics() {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	if metricsInit {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		return
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	sizeClassBuckets = make([]float64, _NumSizeClasses, _NumSizeClasses+1)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// Skip size class 0 which is a stand-in for large objects, but large</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// objects are tracked separately (and they actually get placed in</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// the last bucket, not the first).</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	sizeClassBuckets[0] = 1 <span class="comment">// The smallest allocation is 1 byte in size.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	for i := 1; i &lt; _NumSizeClasses; i++ {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		<span class="comment">// Size classes have an inclusive upper-bound</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		<span class="comment">// and exclusive lower bound (e.g. 48-byte size class is</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		<span class="comment">// (32, 48]) whereas we want and inclusive lower-bound</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		<span class="comment">// and exclusive upper-bound (e.g. 48-byte size class is</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		<span class="comment">// [33, 49)). We can achieve this by shifting all bucket</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		<span class="comment">// boundaries up by 1.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		<span class="comment">// Also, a float64 can precisely represent integers with</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		<span class="comment">// value up to 2^53 and size classes are relatively small</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		<span class="comment">// (nowhere near 2^48 even) so this will give us exact</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		<span class="comment">// boundaries.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		sizeClassBuckets[i] = float64(class_to_size[i] + 1)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	sizeClassBuckets = append(sizeClassBuckets, float64Inf())
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	timeHistBuckets = timeHistogramMetricsBuckets()
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	metrics = map[string]metricData{
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		&#34;/cgo/go-to-c-calls:calls&#34;: {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>				out.scalar = uint64(NumCgoCall())
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			},
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		},
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		&#34;/cpu/classes/gc/mark/assist:cpu-seconds&#34;: {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.gcAssistTime))
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			},
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		},
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		&#34;/cpu/classes/gc/mark/dedicated:cpu-seconds&#34;: {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.gcDedicatedTime))
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			},
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		},
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		&#34;/cpu/classes/gc/mark/idle:cpu-seconds&#34;: {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.gcIdleTime))
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			},
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		},
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		&#34;/cpu/classes/gc/pause:cpu-seconds&#34;: {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.gcPauseTime))
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			},
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		},
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		&#34;/cpu/classes/gc/total:cpu-seconds&#34;: {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.gcTotalTime))
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			},
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		},
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		&#34;/cpu/classes/idle:cpu-seconds&#34;: {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.idleTime))
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			},
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		},
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		&#34;/cpu/classes/scavenge/assist:cpu-seconds&#34;: {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.scavengeAssistTime))
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			},
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		},
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		&#34;/cpu/classes/scavenge/background:cpu-seconds&#34;: {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.scavengeBgTime))
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			},
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		},
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		&#34;/cpu/classes/scavenge/total:cpu-seconds&#34;: {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.scavengeTotalTime))
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			},
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		},
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		&#34;/cpu/classes/total:cpu-seconds&#34;: {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.totalTime))
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			},
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		},
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		&#34;/cpu/classes/user:cpu-seconds&#34;: {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			deps: makeStatDepSet(cpuStatsDep),
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(in.cpuStats.userTime))
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			},
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		},
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		&#34;/gc/cycles/automatic:gc-cycles&#34;: {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>				out.scalar = in.sysStats.gcCyclesDone - in.sysStats.gcCyclesForced
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			},
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		},
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		&#34;/gc/cycles/forced:gc-cycles&#34;: {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>				out.scalar = in.sysStats.gcCyclesForced
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			},
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		},
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		&#34;/gc/cycles/total:gc-cycles&#34;: {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				out.scalar = in.sysStats.gcCyclesDone
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			},
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		},
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		&#34;/gc/scan/globals:bytes&#34;: {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			deps: makeStatDepSet(gcStatsDep),
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>				out.scalar = in.gcStats.globalsScan
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			},
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		},
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		&#34;/gc/scan/heap:bytes&#34;: {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			deps: makeStatDepSet(gcStatsDep),
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>				out.scalar = in.gcStats.heapScan
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			},
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		},
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		&#34;/gc/scan/stack:bytes&#34;: {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			deps: makeStatDepSet(gcStatsDep),
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>				out.scalar = in.gcStats.stackScan
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			},
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		},
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		&#34;/gc/scan/total:bytes&#34;: {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			deps: makeStatDepSet(gcStatsDep),
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>				out.scalar = in.gcStats.totalScan
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			},
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		},
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		&#34;/gc/heap/allocs-by-size:bytes&#34;: {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>				hist := out.float64HistOrInit(sizeClassBuckets)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				hist.counts[len(hist.counts)-1] = in.heapStats.largeAllocCount
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>				<span class="comment">// Cut off the first index which is ostensibly for size class 0,</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>				<span class="comment">// but large objects are tracked separately so it&#39;s actually unused.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>				for i, count := range in.heapStats.smallAllocCount[1:] {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>					hist.counts[i] = count
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>				}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			},
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		},
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		&#34;/gc/heap/allocs:bytes&#34;: {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>				out.scalar = in.heapStats.totalAllocated
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			},
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		},
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		&#34;/gc/heap/allocs:objects&#34;: {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>				out.scalar = in.heapStats.totalAllocs
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			},
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		},
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		&#34;/gc/heap/frees-by-size:bytes&#34;: {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>				hist := out.float64HistOrInit(sizeClassBuckets)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>				hist.counts[len(hist.counts)-1] = in.heapStats.largeFreeCount
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>				<span class="comment">// Cut off the first index which is ostensibly for size class 0,</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>				<span class="comment">// but large objects are tracked separately so it&#39;s actually unused.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>				for i, count := range in.heapStats.smallFreeCount[1:] {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>					hist.counts[i] = count
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			},
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		},
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		&#34;/gc/heap/frees:bytes&#34;: {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				out.scalar = in.heapStats.totalFreed
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			},
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		},
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		&#34;/gc/heap/frees:objects&#34;: {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>				out.scalar = in.heapStats.totalFrees
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			},
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		},
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		&#34;/gc/heap/goal:bytes&#34;: {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>				out.scalar = in.sysStats.heapGoal
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			},
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		},
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		&#34;/gc/gomemlimit:bytes&#34;: {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>				out.scalar = uint64(gcController.memoryLimit.Load())
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			},
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		},
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		&#34;/gc/gogc:percent&#34;: {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>				out.scalar = uint64(gcController.gcPercent.Load())
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			},
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		},
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		&#34;/gc/heap/live:bytes&#34;: {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>				out.scalar = gcController.heapMarked
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			},
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		},
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		&#34;/gc/heap/objects:objects&#34;: {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>				out.scalar = in.heapStats.numObjects
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			},
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		},
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		&#34;/gc/heap/tiny/allocs:objects&#34;: {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>				out.scalar = in.heapStats.tinyAllocCount
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>			},
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		},
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		&#34;/gc/limiter/last-enabled:gc-cycle&#34;: {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>				out.scalar = uint64(gcCPULimiter.lastEnabledCycle.Load())
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>			},
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		},
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		&#34;/gc/pauses:seconds&#34;: {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>				<span class="comment">// N.B. this is identical to /sched/pauses/total/gc:seconds.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>				sched.stwTotalTimeGC.write(out)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>			},
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		},
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		&#34;/gc/stack/starting-size:bytes&#34;: {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>				out.scalar = uint64(startingStackSize)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			},
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		},
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		&#34;/memory/classes/heap/free:bytes&#34;: {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>				out.scalar = uint64(in.heapStats.committed - in.heapStats.inHeap -
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>					in.heapStats.inStacks - in.heapStats.inWorkBufs -
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>					in.heapStats.inPtrScalarBits)
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			},
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		},
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		&#34;/memory/classes/heap/objects:bytes&#34;: {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>				out.scalar = in.heapStats.inObjects
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			},
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		},
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		&#34;/memory/classes/heap/released:bytes&#34;: {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>				out.scalar = uint64(in.heapStats.released)
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>			},
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		},
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		&#34;/memory/classes/heap/stacks:bytes&#34;: {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>				out.scalar = uint64(in.heapStats.inStacks)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			},
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		},
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		&#34;/memory/classes/heap/unused:bytes&#34;: {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep),
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>				out.scalar = uint64(in.heapStats.inHeap) - in.heapStats.inObjects
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			},
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		},
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		&#34;/memory/classes/metadata/mcache/free:bytes&#34;: {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>				out.scalar = in.sysStats.mCacheSys - in.sysStats.mCacheInUse
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>			},
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		},
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		&#34;/memory/classes/metadata/mcache/inuse:bytes&#34;: {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>				out.scalar = in.sysStats.mCacheInUse
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			},
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		},
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		&#34;/memory/classes/metadata/mspan/free:bytes&#34;: {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>				out.scalar = in.sysStats.mSpanSys - in.sysStats.mSpanInUse
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			},
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		},
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		&#34;/memory/classes/metadata/mspan/inuse:bytes&#34;: {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>				out.scalar = in.sysStats.mSpanInUse
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			},
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		},
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		&#34;/memory/classes/metadata/other:bytes&#34;: {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep, sysStatsDep),
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>				out.scalar = uint64(in.heapStats.inWorkBufs+in.heapStats.inPtrScalarBits) + in.sysStats.gcMiscSys
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			},
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		},
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		&#34;/memory/classes/os-stacks:bytes&#34;: {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>				out.scalar = in.sysStats.stacksSys
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			},
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		},
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		&#34;/memory/classes/other:bytes&#34;: {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>				out.scalar = in.sysStats.otherSys
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			},
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		},
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		&#34;/memory/classes/profiling/buckets:bytes&#34;: {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			deps: makeStatDepSet(sysStatsDep),
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>				out.scalar = in.sysStats.buckHashSys
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>			},
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		},
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		&#34;/memory/classes/total:bytes&#34;: {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			deps: makeStatDepSet(heapStatsDep, sysStatsDep),
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>			compute: func(in *statAggregate, out *metricValue) {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>				out.scalar = uint64(in.heapStats.committed+in.heapStats.released) +
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>					in.sysStats.stacksSys + in.sysStats.mSpanSys +
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>					in.sysStats.mCacheSys + in.sysStats.buckHashSys +
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>					in.sysStats.gcMiscSys + in.sysStats.otherSys
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>			},
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		},
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		&#34;/sched/gomaxprocs:threads&#34;: {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>				out.scalar = uint64(gomaxprocs)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			},
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		},
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		&#34;/sched/goroutines:goroutines&#34;: {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>				out.kind = metricKindUint64
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>				out.scalar = uint64(gcount())
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>			},
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		},
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		&#34;/sched/latencies:seconds&#34;: {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>				sched.timeToRun.write(out)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>			},
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		},
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		&#34;/sched/pauses/stopping/gc:seconds&#34;: {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>				sched.stwStoppingTimeGC.write(out)
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			},
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		},
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		&#34;/sched/pauses/stopping/other:seconds&#34;: {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>				sched.stwStoppingTimeOther.write(out)
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			},
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		},
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		&#34;/sched/pauses/total/gc:seconds&#34;: {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>				sched.stwTotalTimeGC.write(out)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			},
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		},
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		&#34;/sched/pauses/total/other:seconds&#34;: {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>				sched.stwTotalTimeOther.write(out)
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>			},
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		},
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		&#34;/sync/mutex/wait/total:seconds&#34;: {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			compute: func(_ *statAggregate, out *metricValue) {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>				out.kind = metricKindFloat64
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>				out.scalar = float64bits(nsToSec(totalMutexWaitTimeNanos()))
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			},
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		},
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	for _, info := range godebugs.All {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		if !info.Opaque {
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			metrics[&#34;/godebug/non-default-behavior/&#34;+info.Name+&#34;:events&#34;] = metricData{compute: compute0}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	metricsInit = true
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>func compute0(_ *statAggregate, out *metricValue) {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	out.kind = metricKindUint64
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	out.scalar = 0
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>type metricReader func() uint64
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>func (f metricReader) compute(_ *statAggregate, out *metricValue) {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	out.kind = metricKindUint64
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	out.scalar = f()
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>}
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span><span class="comment">//go:linkname godebug_registerMetric internal/godebug.registerMetric</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>func godebug_registerMetric(name string, read func() uint64) {
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	metricsLock()
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	initMetrics()
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	d, ok := metrics[name]
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	if !ok {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		throw(&#34;runtime: unexpected metric registration for &#34; + name)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	d.compute = metricReader(read).compute
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	metrics[name] = d
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	metricsUnlock()
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span><span class="comment">// statDep is a dependency on a group of statistics</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span><span class="comment">// that a metric might have.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>type statDep uint
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>const (
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	heapStatsDep statDep = iota <span class="comment">// corresponds to heapStatsAggregate</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	sysStatsDep                 <span class="comment">// corresponds to sysStatsAggregate</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	cpuStatsDep                 <span class="comment">// corresponds to cpuStatsAggregate</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	gcStatsDep                  <span class="comment">// corresponds to gcStatsAggregate</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	numStatsDeps
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>)
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">// statDepSet represents a set of statDeps.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span><span class="comment">// Under the hood, it&#39;s a bitmap.</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>type statDepSet [1]uint64
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span><span class="comment">// makeStatDepSet creates a new statDepSet from a list of statDeps.</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>func makeStatDepSet(deps ...statDep) statDepSet {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	var s statDepSet
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	for _, d := range deps {
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		s[d/64] |= 1 &lt;&lt; (d % 64)
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	return s
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">// difference returns set difference of s from b as a new set.</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>func (s statDepSet) difference(b statDepSet) statDepSet {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	var c statDepSet
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	for i := range s {
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		c[i] = s[i] &amp;^ b[i]
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	return c
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span><span class="comment">// union returns the union of the two sets as a new set.</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>func (s statDepSet) union(b statDepSet) statDepSet {
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	var c statDepSet
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	for i := range s {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		c[i] = s[i] | b[i]
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	}
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	return c
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>}
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span><span class="comment">// empty returns true if there are no dependencies in the set.</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>func (s *statDepSet) empty() bool {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	for _, c := range s {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		if c != 0 {
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			return false
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	return true
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>}
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span><span class="comment">// has returns true if the set contains a given statDep.</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>func (s *statDepSet) has(d statDep) bool {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	return s[d/64]&amp;(1&lt;&lt;(d%64)) != 0
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span><span class="comment">// heapStatsAggregate represents memory stats obtained from the</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span><span class="comment">// runtime. This set of stats is grouped together because they</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span><span class="comment">// depend on each other in some way to make sense of the runtime&#39;s</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span><span class="comment">// current heap memory use. They&#39;re also sharded across Ps, so it</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span><span class="comment">// makes sense to grab them all at once.</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>type heapStatsAggregate struct {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	heapStatsDelta
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	<span class="comment">// Derived from values in heapStatsDelta.</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	<span class="comment">// inObjects is the bytes of memory occupied by objects,</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	inObjects uint64
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	<span class="comment">// numObjects is the number of live objects in the heap.</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	numObjects uint64
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	<span class="comment">// totalAllocated is the total bytes of heap objects allocated</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	<span class="comment">// over the lifetime of the program.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	totalAllocated uint64
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	<span class="comment">// totalFreed is the total bytes of heap objects freed</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	<span class="comment">// over the lifetime of the program.</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	totalFreed uint64
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	<span class="comment">// totalAllocs is the number of heap objects allocated over</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	<span class="comment">// the lifetime of the program.</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	totalAllocs uint64
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	<span class="comment">// totalFrees is the number of heap objects freed over</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	<span class="comment">// the lifetime of the program.</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	totalFrees uint64
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span><span class="comment">// compute populates the heapStatsAggregate with values from the runtime.</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>func (a *heapStatsAggregate) compute() {
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	memstats.heapStats.read(&amp;a.heapStatsDelta)
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	<span class="comment">// Calculate derived stats.</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	a.totalAllocs = a.largeAllocCount
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	a.totalFrees = a.largeFreeCount
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	a.totalAllocated = a.largeAlloc
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	a.totalFreed = a.largeFree
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	for i := range a.smallAllocCount {
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		na := a.smallAllocCount[i]
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		nf := a.smallFreeCount[i]
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		a.totalAllocs += na
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		a.totalFrees += nf
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		a.totalAllocated += na * uint64(class_to_size[i])
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		a.totalFreed += nf * uint64(class_to_size[i])
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	}
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	a.inObjects = a.totalAllocated - a.totalFreed
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	a.numObjects = a.totalAllocs - a.totalFrees
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span><span class="comment">// sysStatsAggregate represents system memory stats obtained</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span><span class="comment">// from the runtime. This set of stats is grouped together because</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span><span class="comment">// they&#39;re all relatively cheap to acquire and generally independent</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span><span class="comment">// of one another and other runtime memory stats. The fact that they</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span><span class="comment">// may be acquired at different times, especially with respect to</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span><span class="comment">// heapStatsAggregate, means there could be some skew, but because of</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span><span class="comment">// these stats are independent, there&#39;s no real consistency issue here.</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>type sysStatsAggregate struct {
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	stacksSys      uint64
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	mSpanSys       uint64
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	mSpanInUse     uint64
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	mCacheSys      uint64
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	mCacheInUse    uint64
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	buckHashSys    uint64
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	gcMiscSys      uint64
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	otherSys       uint64
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	heapGoal       uint64
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	gcCyclesDone   uint64
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	gcCyclesForced uint64
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>}
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span><span class="comment">// compute populates the sysStatsAggregate with values from the runtime.</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>func (a *sysStatsAggregate) compute() {
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	a.stacksSys = memstats.stacks_sys.load()
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	a.buckHashSys = memstats.buckhash_sys.load()
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	a.gcMiscSys = memstats.gcMiscSys.load()
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	a.otherSys = memstats.other_sys.load()
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	a.heapGoal = gcController.heapGoal()
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	a.gcCyclesDone = uint64(memstats.numgc)
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	a.gcCyclesForced = uint64(memstats.numforcedgc)
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		lock(&amp;mheap_.lock)
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		a.mSpanSys = memstats.mspan_sys.load()
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		a.mSpanInUse = uint64(mheap_.spanalloc.inuse)
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		a.mCacheSys = memstats.mcache_sys.load()
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		a.mCacheInUse = uint64(mheap_.cachealloc.inuse)
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		unlock(&amp;mheap_.lock)
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	})
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>}
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span><span class="comment">// cpuStatsAggregate represents CPU stats obtained from the runtime</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span><span class="comment">// acquired together to avoid skew and inconsistencies.</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>type cpuStatsAggregate struct {
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	cpuStats
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span><span class="comment">// compute populates the cpuStatsAggregate with values from the runtime.</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>func (a *cpuStatsAggregate) compute() {
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	a.cpuStats = work.cpuStats
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Update the CPU stats again so that we&#39;re not</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	<span class="comment">// just relying on the STW snapshot. The issue here is that currently</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	<span class="comment">// this will cause non-monotonicity in the &#34;user&#34; CPU time metric.</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	<span class="comment">// a.cpuStats.accumulate(nanotime(), gcphase == _GCmark)</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>}
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span><span class="comment">// gcStatsAggregate represents various GC stats obtained from the runtime</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span><span class="comment">// acquired together to avoid skew and inconsistencies.</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>type gcStatsAggregate struct {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	heapScan    uint64
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	stackScan   uint64
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	globalsScan uint64
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	totalScan   uint64
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>}
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span><span class="comment">// compute populates the gcStatsAggregate with values from the runtime.</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>func (a *gcStatsAggregate) compute() {
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	a.heapScan = gcController.heapScan.Load()
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	a.stackScan = gcController.lastStackScan.Load()
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	a.globalsScan = gcController.globalsScan.Load()
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	a.totalScan = a.heapScan + a.stackScan + a.globalsScan
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span><span class="comment">// nsToSec takes a duration in nanoseconds and converts it to seconds as</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span><span class="comment">// a float64.</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>func nsToSec(ns int64) float64 {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	return float64(ns) / 1e9
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>}
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span><span class="comment">// statAggregate is the main driver of the metrics implementation.</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span><span class="comment">// It contains multiple aggregates of runtime statistics, as well</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span><span class="comment">// as a set of these aggregates that it has populated. The aggregates</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span><span class="comment">// are populated lazily by its ensure method.</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>type statAggregate struct {
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	ensured   statDepSet
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	heapStats heapStatsAggregate
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	sysStats  sysStatsAggregate
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	cpuStats  cpuStatsAggregate
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	gcStats   gcStatsAggregate
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>}
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span><span class="comment">// ensure populates statistics aggregates determined by deps if they</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span><span class="comment">// haven&#39;t yet been populated.</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>func (a *statAggregate) ensure(deps *statDepSet) {
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	missing := deps.difference(a.ensured)
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	if missing.empty() {
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		return
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	}
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	for i := statDep(0); i &lt; numStatsDeps; i++ {
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		if !missing.has(i) {
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>			continue
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>		}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		switch i {
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		case heapStatsDep:
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>			a.heapStats.compute()
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		case sysStatsDep:
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>			a.sysStats.compute()
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		case cpuStatsDep:
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>			a.cpuStats.compute()
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>		case gcStatsDep:
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			a.gcStats.compute()
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		}
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	a.ensured = a.ensured.union(missing)
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>}
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span><span class="comment">// metricKind is a runtime copy of runtime/metrics.ValueKind and</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span><span class="comment">// must be kept structurally identical to that type.</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>type metricKind int
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>const (
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	<span class="comment">// These values must be kept identical to their corresponding Kind* values</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	<span class="comment">// in the runtime/metrics package.</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	metricKindBad metricKind = iota
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	metricKindUint64
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	metricKindFloat64
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	metricKindFloat64Histogram
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>)
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span><span class="comment">// metricSample is a runtime copy of runtime/metrics.Sample and</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span><span class="comment">// must be kept structurally identical to that type.</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>type metricSample struct {
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	name  string
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	value metricValue
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>}
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span><span class="comment">// metricValue is a runtime copy of runtime/metrics.Sample and</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span><span class="comment">// must be kept structurally identical to that type.</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>type metricValue struct {
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	kind    metricKind
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	scalar  uint64         <span class="comment">// contains scalar values for scalar Kinds.</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	pointer unsafe.Pointer <span class="comment">// contains non-scalar values.</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>}
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span><span class="comment">// float64HistOrInit tries to pull out an existing float64Histogram</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span><span class="comment">// from the value, but if none exists, then it allocates one with</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span><span class="comment">// the given buckets.</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>func (v *metricValue) float64HistOrInit(buckets []float64) *metricFloat64Histogram {
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>	var hist *metricFloat64Histogram
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	if v.kind == metricKindFloat64Histogram &amp;&amp; v.pointer != nil {
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>		hist = (*metricFloat64Histogram)(v.pointer)
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>	} else {
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		v.kind = metricKindFloat64Histogram
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		hist = new(metricFloat64Histogram)
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>		v.pointer = unsafe.Pointer(hist)
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	hist.buckets = buckets
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	if len(hist.counts) != len(hist.buckets)-1 {
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>		hist.counts = make([]uint64, len(buckets)-1)
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	}
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	return hist
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span><span class="comment">// metricFloat64Histogram is a runtime copy of runtime/metrics.Float64Histogram</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span><span class="comment">// and must be kept structurally identical to that type.</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>type metricFloat64Histogram struct {
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	counts  []uint64
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	buckets []float64
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>}
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span><span class="comment">// agg is used by readMetrics, and is protected by metricsSema.</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span><span class="comment">// Managed as a global variable because its pointer will be</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span><span class="comment">// an argument to a dynamically-defined function, and we&#39;d</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span><span class="comment">// like to avoid it escaping to the heap.</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>var agg statAggregate
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>type metricName struct {
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	name string
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	kind metricKind
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>}
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span><span class="comment">// readMetricNames is the implementation of runtime/metrics.readMetricNames,</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span><span class="comment">// used by the runtime/metrics test and otherwise unreferenced.</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span><span class="comment">//go:linkname readMetricNames runtime/metrics_test.runtime_readMetricNames</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>func readMetricNames() []string {
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	metricsLock()
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	initMetrics()
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	n := len(metrics)
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	metricsUnlock()
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	list := make([]string, 0, n)
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	metricsLock()
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	for name := range metrics {
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		list = append(list, name)
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	}
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	metricsUnlock()
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	return list
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>}
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span><span class="comment">// readMetrics is the implementation of runtime/metrics.Read.</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span><span class="comment">//go:linkname readMetrics runtime/metrics.runtime_readMetrics</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>func readMetrics(samplesp unsafe.Pointer, len int, cap int) {
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	metricsLock()
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>	<span class="comment">// Ensure the map is initialized.</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	initMetrics()
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	<span class="comment">// Read the metrics.</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	readMetricsLocked(samplesp, len, cap)
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	metricsUnlock()
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span><span class="comment">// readMetricsLocked is the internal, locked portion of readMetrics.</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span><span class="comment">// Broken out for more robust testing. metricsLock must be held and</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span><span class="comment">// initMetrics must have been called already.</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>func readMetricsLocked(samplesp unsafe.Pointer, len int, cap int) {
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	<span class="comment">// Construct a slice from the args.</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	sl := slice{samplesp, len, cap}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	samples := *(*[]metricSample)(unsafe.Pointer(&amp;sl))
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	<span class="comment">// Clear agg defensively.</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	agg = statAggregate{}
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	<span class="comment">// Sample.</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	for i := range samples {
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		sample := &amp;samples[i]
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		data, ok := metrics[sample.name]
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>		if !ok {
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>			sample.value.kind = metricKindBad
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>			continue
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		}
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		<span class="comment">// Ensure we have all the stats we need.</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		<span class="comment">// agg is populated lazily.</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>		agg.ensure(&amp;data.deps)
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>		<span class="comment">// Compute the value based on the stats we have.</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>		data.compute(&amp;agg, &amp;sample.value)
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	}
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>}
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>
</pre><p><a href="metrics.go?m=text">View as plain text</a></p>

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
