<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/metrics/description.go - Go Documentation Server</title>

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
<a href="description.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<a href="http://localhost:8080/src/runtime/metrics">metrics</a>/<span class="text-muted">description.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package metrics
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;internal/godebugs&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// Description describes a runtime metric.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>type Description struct {
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	<span class="comment">// Name is the full name of the metric which includes the unit.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	<span class="comment">// The format of the metric may be described by the following regular expression.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	<span class="comment">// 	^(?P&lt;name&gt;/[^:]+):(?P&lt;unit&gt;[^:*/]+(?:[*/][^:*/]+)*)$</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">// The format splits the name into two components, separated by a colon: a path which always</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// starts with a /, and a machine-parseable unit. The name may contain any valid Unicode</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// codepoint in between / characters, but by convention will try to stick to lowercase</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// characters and hyphens. An example of such a path might be &#34;/memory/heap/free&#34;.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// The unit is by convention a series of lowercase English unit names (singular or plural)</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// without prefixes delimited by &#39;*&#39; or &#39;/&#39;. The unit names may contain any valid Unicode</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// codepoint that is not a delimiter.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// Examples of units might be &#34;seconds&#34;, &#34;bytes&#34;, &#34;bytes/second&#34;, &#34;cpu-seconds&#34;,</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// &#34;byte*cpu-seconds&#34;, and &#34;bytes/second/second&#34;.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// For histograms, multiple units may apply. For instance, the units of the buckets and</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// the count. By convention, for histograms, the units of the count are always &#34;samples&#34;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// with the type of sample evident by the metric&#39;s name, while the unit in the name</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// specifies the buckets&#39; unit.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// A complete name might look like &#34;/memory/heap/free:bytes&#34;.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	Name string
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// Description is an English language sentence describing the metric.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	Description string
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// Kind is the kind of value for this metric.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// The purpose of this field is to allow users to filter out metrics whose values are</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// types which their application may not understand.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	Kind ValueKind
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// Cumulative is whether or not the metric is cumulative. If a cumulative metric is just</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// a single number, then it increases monotonically. If the metric is a distribution,</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// then each bucket count increases monotonically.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// This flag thus indicates whether or not it&#39;s useful to compute a rate from this value.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	Cumulative bool
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// The English language descriptions below must be kept in sync with the</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// descriptions of each metric in doc.go by running &#39;go generate&#39;.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>var allDesc = []Description{
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	{
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		Name:        &#34;/cgo/go-to-c-calls:calls&#34;,
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		Description: &#34;Count of calls made from Go to C by the current process.&#34;,
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	},
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	{
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/gc/mark/assist:cpu-seconds&#34;,
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		Description: &#34;Estimated total CPU time goroutines spent performing GC tasks &#34; +
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			&#34;to assist the GC and prevent it from falling behind the application. &#34; +
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			&#34;This metric is an overestimate, and not directly comparable to &#34; +
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			&#34;system CPU time measurements. Compare only with other /cpu/classes &#34; +
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			&#34;metrics.&#34;,
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	},
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	{
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/gc/mark/dedicated:cpu-seconds&#34;,
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		Description: &#34;Estimated total CPU time spent performing GC tasks on &#34; +
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			&#34;processors (as defined by GOMAXPROCS) dedicated to those tasks. &#34; +
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			&#34;This metric is an overestimate, and not directly comparable to &#34; +
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			&#34;system CPU time measurements. Compare only with other /cpu/classes &#34; +
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			&#34;metrics.&#34;,
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	},
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	{
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/gc/mark/idle:cpu-seconds&#34;,
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		Description: &#34;Estimated total CPU time spent performing GC tasks on &#34; +
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			&#34;spare CPU resources that the Go scheduler could not otherwise find &#34; +
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			&#34;a use for. This should be subtracted from the total GC CPU time to &#34; +
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			&#34;obtain a measure of compulsory GC CPU time. &#34; +
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			&#34;This metric is an overestimate, and not directly comparable to &#34; +
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			&#34;system CPU time measurements. Compare only with other /cpu/classes &#34; +
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			&#34;metrics.&#34;,
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	},
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	{
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/gc/pause:cpu-seconds&#34;,
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		Description: &#34;Estimated total CPU time spent with the application paused by &#34; +
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			&#34;the GC. Even if only one thread is running during the pause, this is &#34; +
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			&#34;computed as GOMAXPROCS times the pause latency because nothing else &#34; +
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			&#34;can be executing. This is the exact sum of samples in &#34; +
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			&#34;/sched/pauses/total/gc:seconds if each sample is multiplied by &#34; +
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			&#34;GOMAXPROCS at the time it is taken. This metric is an overestimate, &#34; +
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			&#34;and not directly comparable to system CPU time measurements. Compare &#34; +
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			&#34;only with other /cpu/classes metrics.&#34;,
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	},
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	{
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/gc/total:cpu-seconds&#34;,
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		Description: &#34;Estimated total CPU time spent performing GC tasks. &#34; +
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			&#34;This metric is an overestimate, and not directly comparable to &#34; +
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			&#34;system CPU time measurements. Compare only with other /cpu/classes &#34; +
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			&#34;metrics. Sum of all metrics in /cpu/classes/gc.&#34;,
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	},
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	{
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/idle:cpu-seconds&#34;,
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		Description: &#34;Estimated total available CPU time not spent executing any Go or Go runtime code. &#34; +
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			&#34;In other words, the part of /cpu/classes/total:cpu-seconds that was unused. &#34; +
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			&#34;This metric is an overestimate, and not directly comparable to &#34; +
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			&#34;system CPU time measurements. Compare only with other /cpu/classes &#34; +
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			&#34;metrics.&#34;,
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	},
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	{
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/scavenge/assist:cpu-seconds&#34;,
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		Description: &#34;Estimated total CPU time spent returning unused memory to the &#34; +
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			&#34;underlying platform in response eagerly in response to memory pressure. &#34; +
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			&#34;This metric is an overestimate, and not directly comparable to &#34; +
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			&#34;system CPU time measurements. Compare only with other /cpu/classes &#34; +
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			&#34;metrics.&#34;,
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	},
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	{
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/scavenge/background:cpu-seconds&#34;,
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		Description: &#34;Estimated total CPU time spent performing background tasks &#34; +
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			&#34;to return unused memory to the underlying platform. &#34; +
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			&#34;This metric is an overestimate, and not directly comparable to &#34; +
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			&#34;system CPU time measurements. Compare only with other /cpu/classes &#34; +
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			&#34;metrics.&#34;,
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	},
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	{
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/scavenge/total:cpu-seconds&#34;,
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		Description: &#34;Estimated total CPU time spent performing tasks that return &#34; +
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			&#34;unused memory to the underlying platform. &#34; +
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			&#34;This metric is an overestimate, and not directly comparable to &#34; +
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			&#34;system CPU time measurements. Compare only with other /cpu/classes &#34; +
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			&#34;metrics. Sum of all metrics in /cpu/classes/scavenge.&#34;,
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	},
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	{
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/total:cpu-seconds&#34;,
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		Description: &#34;Estimated total available CPU time for user Go code &#34; +
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			&#34;or the Go runtime, as defined by GOMAXPROCS. In other words, GOMAXPROCS &#34; +
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			&#34;integrated over the wall-clock duration this process has been executing for. &#34; +
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			&#34;This metric is an overestimate, and not directly comparable to &#34; +
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			&#34;system CPU time measurements. Compare only with other /cpu/classes &#34; +
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			&#34;metrics. Sum of all metrics in /cpu/classes.&#34;,
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	},
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	{
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		Name: &#34;/cpu/classes/user:cpu-seconds&#34;,
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		Description: &#34;Estimated total CPU time spent running user Go code. This may &#34; +
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			&#34;also include some small amount of time spent in the Go runtime. &#34; +
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			&#34;This metric is an overestimate, and not directly comparable to &#34; +
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			&#34;system CPU time measurements. Compare only with other /cpu/classes &#34; +
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			&#34;metrics.&#34;,
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		Kind:       KindFloat64,
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	},
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	{
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		Name:        &#34;/gc/cycles/automatic:gc-cycles&#34;,
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		Description: &#34;Count of completed GC cycles generated by the Go runtime.&#34;,
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	},
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	{
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		Name:        &#34;/gc/cycles/forced:gc-cycles&#34;,
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		Description: &#34;Count of completed GC cycles forced by the application.&#34;,
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	},
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	{
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		Name:        &#34;/gc/cycles/total:gc-cycles&#34;,
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		Description: &#34;Count of all completed GC cycles.&#34;,
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	},
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	{
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		Name: &#34;/gc/gogc:percent&#34;,
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		Description: &#34;Heap size target percentage configured by the user, otherwise 100. This &#34; +
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			&#34;value is set by the GOGC environment variable, and the runtime/debug.SetGCPercent &#34; +
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			&#34;function.&#34;,
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		Kind: KindUint64,
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	},
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	{
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		Name: &#34;/gc/gomemlimit:bytes&#34;,
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		Description: &#34;Go runtime memory limit configured by the user, otherwise &#34; +
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			&#34;math.MaxInt64. This value is set by the GOMEMLIMIT environment variable, and &#34; +
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			&#34;the runtime/debug.SetMemoryLimit function.&#34;,
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		Kind: KindUint64,
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	},
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	{
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		Name: &#34;/gc/heap/allocs-by-size:bytes&#34;,
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		Description: &#34;Distribution of heap allocations by approximate size. &#34; +
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			&#34;Bucket counts increase monotonically. &#34; +
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			&#34;Note that this does not include tiny objects as defined by &#34; +
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			&#34;/gc/heap/tiny/allocs:objects, only tiny blocks.&#34;,
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		Kind:       KindFloat64Histogram,
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	},
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	{
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		Name:        &#34;/gc/heap/allocs:bytes&#34;,
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		Description: &#34;Cumulative sum of memory allocated to the heap by the application.&#34;,
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	},
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	{
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		Name: &#34;/gc/heap/allocs:objects&#34;,
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		Description: &#34;Cumulative count of heap allocations triggered by the application. &#34; +
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			&#34;Note that this does not include tiny objects as defined by &#34; +
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			&#34;/gc/heap/tiny/allocs:objects, only tiny blocks.&#34;,
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		Kind:       KindUint64,
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	},
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	{
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		Name: &#34;/gc/heap/frees-by-size:bytes&#34;,
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		Description: &#34;Distribution of freed heap allocations by approximate size. &#34; +
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			&#34;Bucket counts increase monotonically. &#34; +
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			&#34;Note that this does not include tiny objects as defined by &#34; +
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			&#34;/gc/heap/tiny/allocs:objects, only tiny blocks.&#34;,
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		Kind:       KindFloat64Histogram,
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	},
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	{
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		Name:        &#34;/gc/heap/frees:bytes&#34;,
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		Description: &#34;Cumulative sum of heap memory freed by the garbage collector.&#34;,
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	},
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	{
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		Name: &#34;/gc/heap/frees:objects&#34;,
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		Description: &#34;Cumulative count of heap allocations whose storage was freed &#34; +
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			&#34;by the garbage collector. &#34; +
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			&#34;Note that this does not include tiny objects as defined by &#34; +
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			&#34;/gc/heap/tiny/allocs:objects, only tiny blocks.&#34;,
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		Kind:       KindUint64,
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	},
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	{
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		Name:        &#34;/gc/heap/goal:bytes&#34;,
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		Description: &#34;Heap size target for the end of the GC cycle.&#34;,
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	},
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	{
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		Name:        &#34;/gc/heap/live:bytes&#34;,
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		Description: &#34;Heap memory occupied by live objects that were marked by the previous GC.&#34;,
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	},
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	{
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		Name:        &#34;/gc/heap/objects:objects&#34;,
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		Description: &#34;Number of objects, live or unswept, occupying heap memory.&#34;,
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	},
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	{
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		Name: &#34;/gc/heap/tiny/allocs:objects&#34;,
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		Description: &#34;Count of small allocations that are packed together into blocks. &#34; +
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			&#34;These allocations are counted separately from other allocations &#34; +
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			&#34;because each individual allocation is not tracked by the runtime, &#34; +
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			&#34;only their block. Each block is already accounted for in &#34; +
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			&#34;allocs-by-size and frees-by-size.&#34;,
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		Kind:       KindUint64,
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		Cumulative: true,
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	},
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	{
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		Name: &#34;/gc/limiter/last-enabled:gc-cycle&#34;,
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		Description: &#34;GC cycle the last time the GC CPU limiter was enabled. &#34; +
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			&#34;This metric is useful for diagnosing the root cause of an out-of-memory &#34; +
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			&#34;error, because the limiter trades memory for CPU time when the GC&#39;s CPU &#34; +
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			&#34;time gets too high. This is most likely to occur with use of SetMemoryLimit. &#34; +
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			&#34;The first GC cycle is cycle 1, so a value of 0 indicates that it was never enabled.&#34;,
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		Kind: KindUint64,
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	},
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	{
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		Name:        &#34;/gc/pauses:seconds&#34;,
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		Description: &#34;Deprecated. Prefer the identical /sched/pauses/total/gc:seconds.&#34;,
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		Kind:        KindFloat64Histogram,
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	},
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	{
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		Name:        &#34;/gc/scan/globals:bytes&#34;,
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		Description: &#34;The total amount of global variable space that is scannable.&#34;,
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	},
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	{
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		Name:        &#34;/gc/scan/heap:bytes&#34;,
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		Description: &#34;The total amount of heap space that is scannable.&#34;,
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	},
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	{
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		Name:        &#34;/gc/scan/stack:bytes&#34;,
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		Description: &#34;The number of bytes of stack that were scanned last GC cycle.&#34;,
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	},
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	{
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		Name:        &#34;/gc/scan/total:bytes&#34;,
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		Description: &#34;The total amount space that is scannable. Sum of all metrics in /gc/scan.&#34;,
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	},
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	{
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		Name:        &#34;/gc/stack/starting-size:bytes&#34;,
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		Description: &#34;The stack size of new goroutines.&#34;,
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		Cumulative:  false,
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	},
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	{
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		Name: &#34;/memory/classes/heap/free:bytes&#34;,
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		Description: &#34;Memory that is completely free and eligible to be returned to the underlying system, &#34; +
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			&#34;but has not been. This metric is the runtime&#39;s estimate of free address space that is backed by &#34; +
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			&#34;physical memory.&#34;,
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		Kind: KindUint64,
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	},
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	{
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		Name:        &#34;/memory/classes/heap/objects:bytes&#34;,
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		Description: &#34;Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector.&#34;,
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	},
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	{
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		Name: &#34;/memory/classes/heap/released:bytes&#34;,
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		Description: &#34;Memory that is completely free and has been returned to the underlying system. This &#34; +
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			&#34;metric is the runtime&#39;s estimate of free address space that is still mapped into the process, &#34; +
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>			&#34;but is not backed by physical memory.&#34;,
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		Kind: KindUint64,
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	},
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	{
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		Name: &#34;/memory/classes/heap/stacks:bytes&#34;,
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		Description: &#34;Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. &#34; +
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			&#34;Currently, this represents all stack memory for goroutines. It also includes all OS thread stacks in non-cgo programs. &#34; +
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			&#34;Note that stacks may be allocated differently in the future, and this may change.&#34;,
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		Kind: KindUint64,
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	},
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	{
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		Name:        &#34;/memory/classes/heap/unused:bytes&#34;,
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		Description: &#34;Memory that is reserved for heap objects but is not currently used to hold heap objects.&#34;,
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	},
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	{
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		Name:        &#34;/memory/classes/metadata/mcache/free:bytes&#34;,
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		Description: &#34;Memory that is reserved for runtime mcache structures, but not in-use.&#34;,
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	},
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	{
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		Name:        &#34;/memory/classes/metadata/mcache/inuse:bytes&#34;,
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		Description: &#34;Memory that is occupied by runtime mcache structures that are currently being used.&#34;,
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	},
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	{
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		Name:        &#34;/memory/classes/metadata/mspan/free:bytes&#34;,
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		Description: &#34;Memory that is reserved for runtime mspan structures, but not in-use.&#34;,
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	},
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	{
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		Name:        &#34;/memory/classes/metadata/mspan/inuse:bytes&#34;,
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		Description: &#34;Memory that is occupied by runtime mspan structures that are currently being used.&#34;,
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	},
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	{
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		Name:        &#34;/memory/classes/metadata/other:bytes&#34;,
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		Description: &#34;Memory that is reserved for or used to hold runtime metadata.&#34;,
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	},
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	{
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		Name: &#34;/memory/classes/os-stacks:bytes&#34;,
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		Description: &#34;Stack memory allocated by the underlying operating system. &#34; +
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			&#34;In non-cgo programs this metric is currently zero. This may change in the future.&#34; +
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			&#34;In cgo programs this metric includes OS thread stacks allocated directly from the OS. &#34; +
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			&#34;Currently, this only accounts for one stack in c-shared and c-archive build modes, &#34; +
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			&#34;and other sources of stacks from the OS are not measured. This too may change in the future.&#34;,
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		Kind: KindUint64,
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	},
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	{
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		Name:        &#34;/memory/classes/other:bytes&#34;,
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		Description: &#34;Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more.&#34;,
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	},
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	{
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		Name:        &#34;/memory/classes/profiling/buckets:bytes&#34;,
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		Description: &#34;Memory that is used by the stack trace hash map used for profiling.&#34;,
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	},
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	{
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		Name:        &#34;/memory/classes/total:bytes&#34;,
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		Description: &#34;All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes.&#34;,
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	},
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	{
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		Name:        &#34;/sched/gomaxprocs:threads&#34;,
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		Description: &#34;The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously.&#34;,
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	},
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	{
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		Name:        &#34;/sched/goroutines:goroutines&#34;,
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		Description: &#34;Count of live goroutines.&#34;,
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		Kind:        KindUint64,
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	},
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	{
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		Name:        &#34;/sched/latencies:seconds&#34;,
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		Description: &#34;Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. Bucket counts increase monotonically.&#34;,
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		Kind:        KindFloat64Histogram,
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	},
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	{
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		Name:        &#34;/sched/pauses/stopping/gc:seconds&#34;,
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		Description: &#34;Distribution of individual GC-related stop-the-world stopping latencies. This is the time it takes from deciding to stop the world until all Ps are stopped. This is a subset of the total GC-related stop-the-world time (/sched/pauses/total/gc:seconds). During this time, some threads may be executing. Bucket counts increase monotonically.&#34;,
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		Kind:        KindFloat64Histogram,
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	},
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	{
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		Name:        &#34;/sched/pauses/stopping/other:seconds&#34;,
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		Description: &#34;Distribution of individual non-GC-related stop-the-world stopping latencies. This is the time it takes from deciding to stop the world until all Ps are stopped. This is a subset of the total non-GC-related stop-the-world time (/sched/pauses/total/other:seconds). During this time, some threads may be executing. Bucket counts increase monotonically.&#34;,
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		Kind:        KindFloat64Histogram,
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	},
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	{
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		Name:        &#34;/sched/pauses/total/gc:seconds&#34;,
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		Description: &#34;Distribution of individual GC-related stop-the-world pause latencies. This is the time from deciding to stop the world until the world is started again. Some of this time is spent getting all threads to stop (this is measured directly in /sched/pauses/stopping/gc:seconds), during which some threads may still be running. Bucket counts increase monotonically.&#34;,
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		Kind:        KindFloat64Histogram,
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	},
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	{
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		Name:        &#34;/sched/pauses/total/other:seconds&#34;,
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		Description: &#34;Distribution of individual non-GC-related stop-the-world pause latencies. This is the time from deciding to stop the world until the world is started again. Some of this time is spent getting all threads to stop (measured directly in /sched/pauses/stopping/other:seconds). Bucket counts increase monotonically.&#34;,
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		Kind:        KindFloat64Histogram,
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	},
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	{
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		Name:        &#34;/sync/mutex/wait/total:seconds&#34;,
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		Description: &#34;Approximate cumulative time goroutines have spent blocked on a sync.Mutex, sync.RWMutex, or runtime-internal lock. This metric is useful for identifying global changes in lock contention. Collect a mutex or block profile using the runtime/pprof package for more detailed contention data.&#34;,
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		Kind:        KindFloat64,
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		Cumulative:  true,
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	},
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>func init() {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	<span class="comment">// Insert all the non-default-reporting GODEBUGs into the table,</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	<span class="comment">// preserving the overall sort order.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	i := 0
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	for i &lt; len(allDesc) &amp;&amp; allDesc[i].Name &lt; &#34;/godebug/&#34; {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		i++
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	more := make([]Description, i, len(allDesc)+len(godebugs.All))
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	copy(more, allDesc)
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	for _, info := range godebugs.All {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		if !info.Opaque {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			more = append(more, Description{
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>				Name: &#34;/godebug/non-default-behavior/&#34; + info.Name + &#34;:events&#34;,
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>				Description: &#34;The number of non-default behaviors executed by the &#34; +
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>					info.Package + &#34; package &#34; + &#34;due to a non-default &#34; +
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>					&#34;GODEBUG=&#34; + info.Name + &#34;=... setting.&#34;,
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>				Kind:       KindUint64,
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>				Cumulative: true,
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>			})
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	allDesc = append(more, allDesc[i:]...)
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span><span class="comment">// All returns a slice of containing metric descriptions for all supported metrics.</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>func All() []Description {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	return allDesc
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>
</pre><p><a href="description.go?m=text">View as plain text</a></p>

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
