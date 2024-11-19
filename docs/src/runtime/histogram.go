<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/histogram.go - Go Documentation Server</title>

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
<a href="histogram.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">histogram.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>const (
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	<span class="comment">// For the time histogram type, we use an HDR histogram.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	<span class="comment">// Values are placed in buckets based solely on the most</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	<span class="comment">// significant set bit. Thus, buckets are power-of-2 sized.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">// Values are then placed into sub-buckets based on the value of</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// the next timeHistSubBucketBits most significant bits. Thus,</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// sub-buckets are linear within a bucket.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// Therefore, the number of sub-buckets (timeHistNumSubBuckets)</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// defines the error. This error may be computed as</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// 1/timeHistNumSubBuckets*100%. For example, for 16 sub-buckets</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// per bucket the error is approximately 6%.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// The number of buckets (timeHistNumBuckets), on the</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// other hand, defines the range. To avoid producing a large number</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// of buckets that are close together, especially for small numbers</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// (e.g. 1, 2, 3, 4, 5 ns) that aren&#39;t very useful, timeHistNumBuckets</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// is defined in terms of the least significant bit (timeHistMinBucketBits)</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// that needs to be set before we start bucketing and the most</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// significant bit (timeHistMaxBucketBits) that we bucket before we just</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// dump it into a catch-all bucket.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// As an example, consider the configuration:</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">//    timeHistMinBucketBits = 9</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">//    timeHistMaxBucketBits = 48</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">//    timeHistSubBucketBits = 2</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// Then:</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">//    011000001</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">//    ^--</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">//    │ ^</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">//    │ └---- Next 2 bits -&gt; sub-bucket 3</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">//    └------- Bit 9 unset -&gt; bucket 0</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">//    110000001</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">//    ^--</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">//    │ ^</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">//    │ └---- Next 2 bits -&gt; sub-bucket 2</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">//    └------- Bit 9 set -&gt; bucket 1</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">//    1000000010</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">//    ^-- ^</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">//    │ ^ └-- Lower bits ignored</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">//    │ └---- Next 2 bits -&gt; sub-bucket 0</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">//    └------- Bit 10 set -&gt; bucket 2</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// Following this pattern, bucket 38 will have the bit 46 set. We don&#39;t</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// have any buckets for higher values, so we spill the rest into an overflow</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// bucket containing values of 2^47-1 nanoseconds or approx. 1 day or more.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// This range is more than enough to handle durations produced by the runtime.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	timeHistMinBucketBits = 9
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	timeHistMaxBucketBits = 48 <span class="comment">// Note that this is exclusive; 1 higher than the actual range.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	timeHistSubBucketBits = 2
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	timeHistNumSubBuckets = 1 &lt;&lt; timeHistSubBucketBits
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	timeHistNumBuckets    = timeHistMaxBucketBits - timeHistMinBucketBits + 1
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// Two extra buckets, one for underflow, one for overflow.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	timeHistTotalBuckets = timeHistNumBuckets*timeHistNumSubBuckets + 2
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// timeHistogram represents a distribution of durations in</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// nanoseconds.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// The accuracy and range of the histogram is defined by the</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// timeHistSubBucketBits and timeHistNumBuckets constants.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// It is an HDR histogram with exponentially-distributed</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// buckets and linearly distributed sub-buckets.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// The histogram is safe for concurrent reads and writes.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>type timeHistogram struct {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	counts [timeHistNumBuckets * timeHistNumSubBuckets]atomic.Uint64
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// underflow counts all the times we got a negative duration</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// sample. Because of how time works on some platforms, it&#39;s</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// possible to measure negative durations. We could ignore them,</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// but we record them anyway because it&#39;s better to have some</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// signal that it&#39;s happening than just missing samples.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	underflow atomic.Uint64
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// overflow counts all the times we got a duration that exceeded</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// the range counts represents.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	overflow atomic.Uint64
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// record adds the given duration to the distribution.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// Disallow preemptions and stack growths because this function</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// may run in sensitive locations.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func (h *timeHistogram) record(duration int64) {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	<span class="comment">// If the duration is negative, capture that in underflow.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	if duration &lt; 0 {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		h.underflow.Add(1)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		return
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// bucketBit is the target bit for the bucket which is usually the</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// highest 1 bit, but if we&#39;re less than the minimum, is the highest</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// 1 bit of the minimum (which will be zero in the duration).</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// bucket is the bucket index, which is the bucketBit minus the</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// highest bit of the minimum, plus one to leave room for the catch-all</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// bucket for samples lower than the minimum.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	var bucketBit, bucket uint
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	if l := sys.Len64(uint64(duration)); l &lt; timeHistMinBucketBits {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		bucketBit = timeHistMinBucketBits
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		bucket = 0 <span class="comment">// bucketBit - timeHistMinBucketBits</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	} else {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		bucketBit = uint(l)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		bucket = bucketBit - timeHistMinBucketBits + 1
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// If the bucket we computed is greater than the number of buckets,</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// count that in overflow.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	if bucket &gt;= timeHistNumBuckets {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		h.overflow.Add(1)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		return
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// The sub-bucket index is just next timeHistSubBucketBits after the bucketBit.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	subBucket := uint(duration&gt;&gt;(bucketBit-1-timeHistSubBucketBits)) % timeHistNumSubBuckets
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	h.counts[bucket*timeHistNumSubBuckets+subBucket].Add(1)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// write dumps the histogram to the passed metricValue as a float64 histogram.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>func (h *timeHistogram) write(out *metricValue) {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	hist := out.float64HistOrInit(timeHistBuckets)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// The bottom-most bucket, containing negative values, is tracked</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// separately as underflow, so fill that in manually and then iterate</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// over the rest.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	hist.counts[0] = h.underflow.Load()
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	for i := range h.counts {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		hist.counts[i+1] = h.counts[i].Load()
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	hist.counts[len(hist.counts)-1] = h.overflow.Load()
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>const (
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	fInf    = 0x7FF0000000000000
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	fNegInf = 0xFFF0000000000000
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>func float64Inf() float64 {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	inf := uint64(fInf)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	return *(*float64)(unsafe.Pointer(&amp;inf))
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func float64NegInf() float64 {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	inf := uint64(fNegInf)
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	return *(*float64)(unsafe.Pointer(&amp;inf))
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">// timeHistogramMetricsBuckets generates a slice of boundaries for</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// the timeHistogram. These boundaries are represented in seconds,</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// not nanoseconds like the timeHistogram represents durations.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>func timeHistogramMetricsBuckets() []float64 {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	b := make([]float64, timeHistTotalBuckets+1)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// Underflow bucket.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	b[0] = float64NegInf()
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	for j := 0; j &lt; timeHistNumSubBuckets; j++ {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		<span class="comment">// No bucket bit for the first few buckets. Just sub-bucket bits after the</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		<span class="comment">// min bucket bit.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		bucketNanos := uint64(j) &lt;&lt; (timeHistMinBucketBits - 1 - timeHistSubBucketBits)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		<span class="comment">// Convert nanoseconds to seconds via a division.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		<span class="comment">// These values will all be exactly representable by a float64.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		b[j+1] = float64(bucketNanos) / 1e9
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// Generate the rest of the buckets. It&#39;s easier to reason</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// about if we cut out the 0&#39;th bucket.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	for i := timeHistMinBucketBits; i &lt; timeHistMaxBucketBits; i++ {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		for j := 0; j &lt; timeHistNumSubBuckets; j++ {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			<span class="comment">// Set the bucket bit.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			bucketNanos := uint64(1) &lt;&lt; (i - 1)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			<span class="comment">// Set the sub-bucket bits.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			bucketNanos |= uint64(j) &lt;&lt; (i - 1 - timeHistSubBucketBits)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			<span class="comment">// The index for this bucket is going to be the (i+1)&#39;th bucket</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			<span class="comment">// (note that we&#39;re starting from zero, but handled the first bucket</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			<span class="comment">// earlier, so we need to compensate), and the j&#39;th sub bucket.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			<span class="comment">// Add 1 because we left space for -Inf.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			bucketIndex := (i-timeHistMinBucketBits+1)*timeHistNumSubBuckets + j + 1
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			<span class="comment">// Convert nanoseconds to seconds via a division.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			<span class="comment">// These values will all be exactly representable by a float64.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			b[bucketIndex] = float64(bucketNanos) / 1e9
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// Overflow bucket.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	b[len(b)-2] = float64(uint64(1)&lt;&lt;(timeHistMaxBucketBits-1)) / 1e9
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	b[len(b)-1] = float64Inf()
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	return b
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
</pre><p><a href="histogram.go?m=text">View as plain text</a></p>

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
