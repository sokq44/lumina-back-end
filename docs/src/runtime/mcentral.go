<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mcentral.go - Go Documentation Server</title>

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
<a href="mcentral.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mcentral.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Central free lists.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// See malloc.go for an overview.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// The mcentral doesn&#39;t actually contain the list of free objects; the mspan does.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// Each mcentral is two lists of mspans: those with free objects (c-&gt;nonempty)</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// and those that are completely allocated (c-&gt;empty).</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>package runtime
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>import (
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// Central list of free objects of a given size.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>type mcentral struct {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	_         sys.NotInHeap
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	spanclass spanClass
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// partial and full contain two mspan sets: one of swept in-use</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// spans, and one of unswept in-use spans. These two trade</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// roles on each GC cycle. The unswept set is drained either by</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// allocation or by the background sweeper in every GC cycle,</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// so only two roles are necessary.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// sweepgen is increased by 2 on each GC cycle, so the swept</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// spans are in partial[sweepgen/2%2] and the unswept spans are in</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// partial[1-sweepgen/2%2]. Sweeping pops spans from the</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// unswept set and pushes spans that are still in-use on the</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// swept set. Likewise, allocating an in-use span pushes it</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// on the swept set.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// Some parts of the sweeper can sweep arbitrary spans, and hence</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// can&#39;t remove them from the unswept set, but will add the span</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// to the appropriate swept list. As a result, the parts of the</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// sweeper and mcentral that do consume from the unswept list may</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// encounter swept spans, and these should be ignored.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	partial [2]spanSet <span class="comment">// list of spans with a free object</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	full    [2]spanSet <span class="comment">// list of spans with no free objects</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// Initialize a single central free list.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>func (c *mcentral) init(spc spanClass) {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	c.spanclass = spc
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	lockInit(&amp;c.partial[0].spineLock, lockRankSpanSetSpine)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	lockInit(&amp;c.partial[1].spineLock, lockRankSpanSetSpine)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	lockInit(&amp;c.full[0].spineLock, lockRankSpanSetSpine)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	lockInit(&amp;c.full[1].spineLock, lockRankSpanSetSpine)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// partialUnswept returns the spanSet which holds partially-filled</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// unswept spans for this sweepgen.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func (c *mcentral) partialUnswept(sweepgen uint32) *spanSet {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	return &amp;c.partial[1-sweepgen/2%2]
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// partialSwept returns the spanSet which holds partially-filled</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// swept spans for this sweepgen.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>func (c *mcentral) partialSwept(sweepgen uint32) *spanSet {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	return &amp;c.partial[sweepgen/2%2]
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// fullUnswept returns the spanSet which holds unswept spans without any</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// free slots for this sweepgen.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>func (c *mcentral) fullUnswept(sweepgen uint32) *spanSet {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	return &amp;c.full[1-sweepgen/2%2]
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// fullSwept returns the spanSet which holds swept spans without any</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// free slots for this sweepgen.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>func (c *mcentral) fullSwept(sweepgen uint32) *spanSet {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	return &amp;c.full[sweepgen/2%2]
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// Allocate a span to use in an mcache.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>func (c *mcentral) cacheSpan() *mspan {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// Deduct credit for this span allocation and sweep if necessary.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	spanBytes := uintptr(class_to_allocnpages[c.spanclass.sizeclass()]) * _PageSize
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	deductSweepCredit(spanBytes, 0)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	traceDone := false
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	trace := traceAcquire()
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		trace.GCSweepStart()
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// If we sweep spanBudget spans without finding any free</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// space, just allocate a fresh span. This limits the amount</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// of time we can spend trying to find free space and</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// amortizes the cost of small object sweeping over the</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// benefit of having a full free span to allocate from. By</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// setting this to 100, we limit the space overhead to 1%.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// TODO(austin,mknyszek): This still has bad worst-case</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// throughput. For example, this could find just one free slot</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// on the 100th swept span. That limits allocation latency, but</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// still has very poor throughput. We could instead keep a</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// running free-to-used budget and switch to fresh span</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// allocation if the budget runs low.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	spanBudget := 100
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	var s *mspan
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	var sl sweepLocker
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// Try partial swept spans first.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	sg := mheap_.sweepgen
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	if s = c.partialSwept(sg).pop(); s != nil {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		goto havespan
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	sl = sweep.active.begin()
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	if sl.valid {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		<span class="comment">// Now try partial unswept spans.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		for ; spanBudget &gt;= 0; spanBudget-- {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			s = c.partialUnswept(sg).pop()
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			if s == nil {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>				break
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			if s, ok := sl.tryAcquire(s); ok {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>				<span class="comment">// We got ownership of the span, so let&#39;s sweep it and use it.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>				s.sweep(true)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>				sweep.active.end(sl)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>				goto havespan
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			<span class="comment">// We failed to get ownership of the span, which means it&#39;s being or</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			<span class="comment">// has been swept by an asynchronous sweeper that just couldn&#39;t remove it</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			<span class="comment">// from the unswept list. That sweeper took ownership of the span and</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			<span class="comment">// responsibility for either freeing it to the heap or putting it on the</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			<span class="comment">// right swept list. Either way, we should just ignore it (and it&#39;s unsafe</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			<span class="comment">// for us to do anything else).</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		<span class="comment">// Now try full unswept spans, sweeping them and putting them into the</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		<span class="comment">// right list if we fail to get a span.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		for ; spanBudget &gt;= 0; spanBudget-- {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			s = c.fullUnswept(sg).pop()
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			if s == nil {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>				break
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			if s, ok := sl.tryAcquire(s); ok {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>				<span class="comment">// We got ownership of the span, so let&#39;s sweep it.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>				s.sweep(true)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				<span class="comment">// Check if there&#39;s any free space.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>				freeIndex := s.nextFreeIndex()
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>				if freeIndex != s.nelems {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>					s.freeindex = freeIndex
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>					sweep.active.end(sl)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>					goto havespan
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>				}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>				<span class="comment">// Add it to the swept list, because sweeping didn&#39;t give us any free space.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>				c.fullSwept(sg).push(s.mspan)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			<span class="comment">// See comment for partial unswept spans.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		sweep.active.end(sl)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	trace = traceAcquire()
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		trace.GCSweepDone()
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		traceDone = true
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// We failed to get a span from the mcentral so get one from mheap.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	s = c.grow()
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	if s == nil {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		return nil
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// At this point s is a span that should have free slots.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>havespan:
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	if !traceDone {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		trace := traceAcquire()
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		if trace.ok() {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			trace.GCSweepDone()
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			traceRelease(trace)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	n := int(s.nelems) - int(s.allocCount)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if n == 0 || s.freeindex == s.nelems || s.allocCount == s.nelems {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		throw(&#34;span has no free objects&#34;)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	freeByteBase := s.freeindex &amp;^ (64 - 1)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	whichByte := freeByteBase / 8
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// Init alloc bits cache.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	s.refillAllocCache(whichByte)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// Adjust the allocCache so that s.freeindex corresponds to the low bit in</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// s.allocCache.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	s.allocCache &gt;&gt;= s.freeindex % 64
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	return s
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// Return span from an mcache.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">// s must have a span class corresponding to this</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">// mcentral and it must not be empty.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>func (c *mcentral) uncacheSpan(s *mspan) {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	if s.allocCount == 0 {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		throw(&#34;uncaching span but s.allocCount == 0&#34;)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	sg := mheap_.sweepgen
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	stale := s.sweepgen == sg+1
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// Fix up sweepgen.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	if stale {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		<span class="comment">// Span was cached before sweep began. It&#39;s our</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		<span class="comment">// responsibility to sweep it.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		<span class="comment">// Set sweepgen to indicate it&#39;s not cached but needs</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// sweeping and can&#39;t be allocated from. sweep will</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		<span class="comment">// set s.sweepgen to indicate s is swept.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		atomic.Store(&amp;s.sweepgen, sg-1)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	} else {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		<span class="comment">// Indicate that s is no longer cached.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		atomic.Store(&amp;s.sweepgen, sg)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// Put the span in the appropriate place.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	if stale {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		<span class="comment">// It&#39;s stale, so just sweep it. Sweeping will put it on</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		<span class="comment">// the right list.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		<span class="comment">// We don&#39;t use a sweepLocker here. Stale cached spans</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		<span class="comment">// aren&#39;t in the global sweep lists, so mark termination</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		<span class="comment">// itself holds up sweep completion until all mcaches</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		<span class="comment">// have been swept.</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		ss := sweepLocked{s}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		ss.sweep(false)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	} else {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		if int(s.nelems)-int(s.allocCount) &gt; 0 {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			<span class="comment">// Put it back on the partial swept list.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			c.partialSwept(sg).push(s)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		} else {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			<span class="comment">// There&#39;s no free space and it&#39;s not stale, so put it on the</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			<span class="comment">// full swept list.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			c.fullSwept(sg).push(s)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">// grow allocates a new empty span from the heap and initializes it for c&#39;s size class.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>func (c *mcentral) grow() *mspan {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	npages := uintptr(class_to_allocnpages[c.spanclass.sizeclass()])
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	size := uintptr(class_to_size[c.spanclass.sizeclass()])
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	s := mheap_.alloc(npages, c.spanclass)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	if s == nil {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		return nil
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">// Use division by multiplication and shifts to quickly compute:</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// n := (npages &lt;&lt; _PageShift) / size</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	n := s.divideByElemSize(npages &lt;&lt; _PageShift)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	s.limit = s.base() + size*n
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	s.initHeapBits(false)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	return s
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
</pre><p><a href="mcentral.go?m=text">View as plain text</a></p>

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
