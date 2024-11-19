<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mgcsweep.go - Go Documentation Server</title>

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
<a href="mgcsweep.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mgcsweep.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Garbage collector: sweeping</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// The sweeper consists of two different algorithms:</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// * The object reclaimer finds and frees unmarked slots in spans. It</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//   can free a whole span if none of the objects are marked, but that</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//   isn&#39;t its goal. This can be driven either synchronously by</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//   mcentral.cacheSpan for mcentral spans, or asynchronously by</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//   sweepone, which looks at all the mcentral lists.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// * The span reclaimer looks for spans that contain no marked objects</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//   and frees whole spans. This is a separate algorithm because</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//   freeing whole spans is the hardest task for the object reclaimer,</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//   but is critical when allocating new spans. The entry point for</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//   this is mheap_.reclaim and it&#39;s driven by a sequential scan of</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//   the page marks bitmap in the heap arenas.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// Both algorithms ultimately call mspan.sweep, which sweeps a single</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// heap span.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>package runtime
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>import (
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	&#34;internal/goexperiment&#34;
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>var sweep sweepdata
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// State of background sweep.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>type sweepdata struct {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	lock   mutex
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	g      *g
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	parked bool
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// active tracks outstanding sweepers and the sweep</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// termination condition.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	active activeSweep
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// centralIndex is the current unswept span class.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// It represents an index into the mcentral span</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// sets. Accessed and updated via its load and</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// update methods. Not protected by a lock.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// Reset at mark termination.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// Used by mheap.nextSpanForSweep.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	centralIndex sweepClass
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// sweepClass is a spanClass and one bit to represent whether we&#39;re currently</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// sweeping partial or full spans.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>type sweepClass uint32
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>const (
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	numSweepClasses            = numSpanClasses * 2
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	sweepClassDone  sweepClass = sweepClass(^uint32(0))
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>func (s *sweepClass) load() sweepClass {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	return sweepClass(atomic.Load((*uint32)(s)))
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>func (s *sweepClass) update(sNew sweepClass) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// Only update *s if its current value is less than sNew,</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// since *s increases monotonically.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	sOld := s.load()
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	for sOld &lt; sNew &amp;&amp; !atomic.Cas((*uint32)(s), uint32(sOld), uint32(sNew)) {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		sOld = s.load()
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): This isn&#39;t the only place we have</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// an atomic monotonically increasing counter. It would</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// be nice to have an &#34;atomic max&#34; which is just implemented</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// as the above on most architectures. Some architectures</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// like RISC-V however have native support for an atomic max.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func (s *sweepClass) clear() {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	atomic.Store((*uint32)(s), 0)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// split returns the underlying span class as well as</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// whether we&#39;re interested in the full or partial</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// unswept lists for that class, indicated as a boolean</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// (true means &#34;full&#34;).</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func (s sweepClass) split() (spc spanClass, full bool) {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	return spanClass(s &gt;&gt; 1), s&amp;1 == 0
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// nextSpanForSweep finds and pops the next span for sweeping from the</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// central sweep buffers. It returns ownership of the span to the caller.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// Returns nil if no such span exists.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>func (h *mheap) nextSpanForSweep() *mspan {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	sg := h.sweepgen
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	for sc := sweep.centralIndex.load(); sc &lt; numSweepClasses; sc++ {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		spc, full := sc.split()
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		c := &amp;h.central[spc].mcentral
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		var s *mspan
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		if full {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			s = c.fullUnswept(sg).pop()
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		} else {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			s = c.partialUnswept(sg).pop()
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		if s != nil {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			<span class="comment">// Write down that we found something so future sweepers</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			<span class="comment">// can start from here.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			sweep.centralIndex.update(sc)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			return s
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// Write down that we found nothing.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	sweep.centralIndex.update(sweepClassDone)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	return nil
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>const sweepDrainedMask = 1 &lt;&lt; 31
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// activeSweep is a type that captures whether sweeping</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// is done, and whether there are any outstanding sweepers.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// Every potential sweeper must call begin() before they look</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// for work, and end() after they&#39;ve finished sweeping.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>type activeSweep struct {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// state is divided into two parts.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// The top bit (masked by sweepDrainedMask) is a boolean</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// value indicating whether all the sweep work has been</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// drained from the queue.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// The rest of the bits are a counter, indicating the</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// number of outstanding concurrent sweepers.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	state atomic.Uint32
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// begin registers a new sweeper. Returns a sweepLocker</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// for acquiring spans for sweeping. Any outstanding sweeper blocks</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// sweep termination.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// If the sweepLocker is invalid, the caller can be sure that all</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// outstanding sweep work has been drained, so there is nothing left</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// to sweep. Note that there may be sweepers currently running, so</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// this does not indicate that all sweeping has completed.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// Even if the sweepLocker is invalid, its sweepGen is always valid.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>func (a *activeSweep) begin() sweepLocker {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	for {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		state := a.state.Load()
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		if state&amp;sweepDrainedMask != 0 {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			return sweepLocker{mheap_.sweepgen, false}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		if a.state.CompareAndSwap(state, state+1) {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			return sweepLocker{mheap_.sweepgen, true}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// end deregisters a sweeper. Must be called once for each time</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// begin is called if the sweepLocker is valid.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>func (a *activeSweep) end(sl sweepLocker) {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	if sl.sweepGen != mheap_.sweepgen {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		throw(&#34;sweeper left outstanding across sweep generations&#34;)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	for {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		state := a.state.Load()
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		if (state&amp;^sweepDrainedMask)-1 &gt;= sweepDrainedMask {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			throw(&#34;mismatched begin/end of activeSweep&#34;)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		if a.state.CompareAndSwap(state, state-1) {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			if state != sweepDrainedMask {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>				return
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			if debug.gcpacertrace &gt; 0 {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>				live := gcController.heapLive.Load()
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>				print(&#34;pacer: sweep done at heap size &#34;, live&gt;&gt;20, &#34;MB; allocated &#34;, (live-mheap_.sweepHeapLiveBasis)&gt;&gt;20, &#34;MB during sweep; swept &#34;, mheap_.pagesSwept.Load(), &#34; pages at &#34;, mheap_.sweepPagesPerByte, &#34; pages/byte\n&#34;)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			return
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// markDrained marks the active sweep cycle as having drained</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// all remaining work. This is safe to be called concurrently</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// with all other methods of activeSweep, though may race.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// Returns true if this call was the one that actually performed</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// the mark.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>func (a *activeSweep) markDrained() bool {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	for {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		state := a.state.Load()
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		if state&amp;sweepDrainedMask != 0 {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			return false
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		if a.state.CompareAndSwap(state, state|sweepDrainedMask) {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			return true
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">// sweepers returns the current number of active sweepers.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>func (a *activeSweep) sweepers() uint32 {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	return a.state.Load() &amp;^ sweepDrainedMask
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// isDone returns true if all sweep work has been drained and no more</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// outstanding sweepers exist. That is, when the sweep phase is</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// completely done.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>func (a *activeSweep) isDone() bool {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	return a.state.Load() == sweepDrainedMask
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// reset sets up the activeSweep for the next sweep cycle.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// The world must be stopped.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>func (a *activeSweep) reset() {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	a.state.Store(0)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// finishsweep_m ensures that all spans are swept.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// The world must be stopped. This ensures there are no sweeps in</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// progress.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>func finishsweep_m() {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">// Sweeping must be complete before marking commences, so</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	<span class="comment">// sweep any unswept spans. If this is a concurrent GC, there</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">// shouldn&#39;t be any spans left to sweep, so this should finish</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// instantly. If GC was forced before the concurrent sweep</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// finished, there may be spans to sweep.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	for sweepone() != ^uintptr(0) {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// Make sure there aren&#39;t any outstanding sweepers left.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	<span class="comment">// At this point, with the world stopped, it means one of two</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// things. Either we were able to preempt a sweeper, or that</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// a sweeper didn&#39;t call sweep.active.end when it should have.</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// Both cases indicate a bug, so throw.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	if sweep.active.sweepers() != 0 {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		throw(&#34;active sweepers found at start of mark phase&#34;)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// Reset all the unswept buffers, which should be empty.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// Do this in sweep termination as opposed to mark termination</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// so that we can catch unswept spans and reclaim blocks as</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// soon as possible.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	sg := mheap_.sweepgen
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	for i := range mheap_.central {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		c := &amp;mheap_.central[i].mcentral
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		c.partialUnswept(sg).reset()
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		c.fullUnswept(sg).reset()
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// Sweeping is done, so there won&#39;t be any new memory to</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// scavenge for a bit.</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// If the scavenger isn&#39;t already awake, wake it up. There&#39;s</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	<span class="comment">// definitely work for it to do at this point.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	scavenger.wake()
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	nextMarkBitArenaEpoch()
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>func bgsweep(c chan int) {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	sweep.g = getg()
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	lockInit(&amp;sweep.lock, lockRankSweep)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	lock(&amp;sweep.lock)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	sweep.parked = true
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	c &lt;- 1
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	goparkunlock(&amp;sweep.lock, waitReasonGCSweepWait, traceBlockGCSweep, 1)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	for {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		<span class="comment">// bgsweep attempts to be a &#34;low priority&#34; goroutine by intentionally</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		<span class="comment">// yielding time. It&#39;s OK if it doesn&#39;t run, because goroutines allocating</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		<span class="comment">// memory will sweep and ensure that all spans are swept before the next</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		<span class="comment">// GC cycle. We really only want to run when we&#39;re idle.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		<span class="comment">// However, calling Gosched after each span swept produces a tremendous</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		<span class="comment">// amount of tracing events, sometimes up to 50% of events in a trace. It&#39;s</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		<span class="comment">// also inefficient to call into the scheduler so much because sweeping a</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		<span class="comment">// single span is in general a very fast operation, taking as little as 30 ns</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		<span class="comment">// on modern hardware. (See #54767.)</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		<span class="comment">// As a result, bgsweep sweeps in batches, and only calls into the scheduler</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		<span class="comment">// at the end of every batch. Furthermore, it only yields its time if there</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		<span class="comment">// isn&#39;t spare idle time available on other cores. If there&#39;s available idle</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		<span class="comment">// time, helping to sweep can reduce allocation latencies by getting ahead of</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		<span class="comment">// the proportional sweeper and having spans ready to go for allocation.</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		const sweepBatchSize = 10
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		nSwept := 0
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		for sweepone() != ^uintptr(0) {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			nSwept++
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			if nSwept%sweepBatchSize == 0 {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>				goschedIfBusy()
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		for freeSomeWbufs(true) {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			<span class="comment">// N.B. freeSomeWbufs is already batched internally.</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			goschedIfBusy()
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		lock(&amp;sweep.lock)
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		if !isSweepDone() {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>			<span class="comment">// This can happen if a GC runs between</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>			<span class="comment">// gosweepone returning ^0 above</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>			<span class="comment">// and the lock being acquired.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>			unlock(&amp;sweep.lock)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>			continue
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		sweep.parked = true
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		goparkunlock(&amp;sweep.lock, waitReasonGCSweepWait, traceBlockGCSweep, 1)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span><span class="comment">// sweepLocker acquires sweep ownership of spans.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>type sweepLocker struct {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	<span class="comment">// sweepGen is the sweep generation of the heap.</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	sweepGen uint32
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	valid    bool
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// sweepLocked represents sweep ownership of a span.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>type sweepLocked struct {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	*mspan
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">// tryAcquire attempts to acquire sweep ownership of span s. If it</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// successfully acquires ownership, it blocks sweep completion.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func (l *sweepLocker) tryAcquire(s *mspan) (sweepLocked, bool) {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	if !l.valid {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		throw(&#34;use of invalid sweepLocker&#34;)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">// Check before attempting to CAS.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	if atomic.Load(&amp;s.sweepgen) != l.sweepGen-2 {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		return sweepLocked{}, false
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	<span class="comment">// Attempt to acquire sweep ownership of s.</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	if !atomic.Cas(&amp;s.sweepgen, l.sweepGen-2, l.sweepGen-1) {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		return sweepLocked{}, false
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	return sweepLocked{s}, true
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// sweepone sweeps some unswept heap span and returns the number of pages returned</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// to the heap, or ^uintptr(0) if there was nothing to sweep.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>func sweepone() uintptr {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	gp := getg()
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	<span class="comment">// Increment locks to ensure that the goroutine is not preempted</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	<span class="comment">// in the middle of sweep thus leaving the span in an inconsistent state for next GC</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	gp.m.locks++
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	<span class="comment">// TODO(austin): sweepone is almost always called in a loop;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	<span class="comment">// lift the sweepLocker into its callers.</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	sl := sweep.active.begin()
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	if !sl.valid {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		gp.m.locks--
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		return ^uintptr(0)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	<span class="comment">// Find a span to sweep.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	npages := ^uintptr(0)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	var noMoreWork bool
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	for {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		s := mheap_.nextSpanForSweep()
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		if s == nil {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			noMoreWork = sweep.active.markDrained()
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			break
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		if state := s.state.get(); state != mSpanInUse {
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			<span class="comment">// This can happen if direct sweeping already</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			<span class="comment">// swept this span, but in that case the sweep</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			<span class="comment">// generation should always be up-to-date.</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			if !(s.sweepgen == sl.sweepGen || s.sweepgen == sl.sweepGen+3) {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>				print(&#34;runtime: bad span s.state=&#34;, state, &#34; s.sweepgen=&#34;, s.sweepgen, &#34; sweepgen=&#34;, sl.sweepGen, &#34;\n&#34;)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>				throw(&#34;non in-use span in unswept list&#34;)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>			continue
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		if s, ok := sl.tryAcquire(s); ok {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			<span class="comment">// Sweep the span we found.</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			npages = s.npages
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			if s.sweep(false) {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>				<span class="comment">// Whole span was freed. Count it toward the</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>				<span class="comment">// page reclaimer credit since these pages can</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>				<span class="comment">// now be used for span allocation.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>				mheap_.reclaimCredit.Add(npages)
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			} else {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>				<span class="comment">// Span is still in-use, so this returned no</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>				<span class="comment">// pages to the heap and the span needs to</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>				<span class="comment">// move to the swept in-use list.</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>				npages = 0
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			break
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	}
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	sweep.active.end(sl)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	if noMoreWork {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		<span class="comment">// The sweep list is empty. There may still be</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		<span class="comment">// concurrent sweeps running, but we&#39;re at least very</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		<span class="comment">// close to done sweeping.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		<span class="comment">// Move the scavenge gen forward (signaling</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		<span class="comment">// that there&#39;s new work to do) and wake the scavenger.</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		<span class="comment">// The scavenger is signaled by the last sweeper because once</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		<span class="comment">// sweeping is done, we will definitely have useful work for</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		<span class="comment">// the scavenger to do, since the scavenger only runs over the</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		<span class="comment">// heap once per GC cycle. This update is not done during sweep</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		<span class="comment">// termination because in some cases there may be a long delay</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		<span class="comment">// between sweep done and sweep termination (e.g. not enough</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		<span class="comment">// allocations to trigger a GC) which would be nice to fill in</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		<span class="comment">// with scavenging work.</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		if debug.scavtrace &gt; 0 {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			systemstack(func() {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>				lock(&amp;mheap_.lock)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>				<span class="comment">// Get released stats.</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>				releasedBg := mheap_.pages.scav.releasedBg.Load()
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>				releasedEager := mheap_.pages.scav.releasedEager.Load()
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>				<span class="comment">// Print the line.</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>				printScavTrace(releasedBg, releasedEager, false)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>				<span class="comment">// Update the stats.</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>				mheap_.pages.scav.releasedBg.Add(-releasedBg)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>				mheap_.pages.scav.releasedEager.Add(-releasedEager)
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>				unlock(&amp;mheap_.lock)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			})
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		scavenger.ready()
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	gp.m.locks--
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	return npages
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span><span class="comment">// isSweepDone reports whether all spans are swept.</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span><span class="comment">// Note that this condition may transition from false to true at any</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span><span class="comment">// time as the sweeper runs. It may transition from true to false if a</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span><span class="comment">// GC runs; to prevent that the caller must be non-preemptible or must</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span><span class="comment">// somehow block GC progress.</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>func isSweepDone() bool {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	return sweep.active.isDone()
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span><span class="comment">// Returns only when span s has been swept.</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>func (s *mspan) ensureSwept() {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	<span class="comment">// Caller must disable preemption.</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	<span class="comment">// Otherwise when this function returns the span can become unswept again</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	<span class="comment">// (if GC is triggered on another goroutine).</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	gp := getg()
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	if gp.m.locks == 0 &amp;&amp; gp.m.mallocing == 0 &amp;&amp; gp != gp.m.g0 {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		throw(&#34;mspan.ensureSwept: m is not locked&#34;)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	<span class="comment">// If this operation fails, then that means that there are</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	<span class="comment">// no more spans to be swept. In this case, either s has already</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	<span class="comment">// been swept, or is about to be acquired for sweeping and swept.</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	sl := sweep.active.begin()
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	if sl.valid {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		<span class="comment">// The caller must be sure that the span is a mSpanInUse span.</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		if s, ok := sl.tryAcquire(s); ok {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			s.sweep(false)
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			sweep.active.end(sl)
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			return
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		sweep.active.end(sl)
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	<span class="comment">// Unfortunately we can&#39;t sweep the span ourselves. Somebody else</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	<span class="comment">// got to it first. We don&#39;t have efficient means to wait, but that&#39;s</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	<span class="comment">// OK, it will be swept fairly soon.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	for {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		spangen := atomic.Load(&amp;s.sweepgen)
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		if spangen == sl.sweepGen || spangen == sl.sweepGen+3 {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			break
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		osyield()
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span><span class="comment">// sweep frees or collects finalizers for blocks not marked in the mark phase.</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span><span class="comment">// It clears the mark bits in preparation for the next GC round.</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span><span class="comment">// Returns true if the span was returned to heap.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span><span class="comment">// If preserve=true, don&#39;t return it to heap nor relink in mcentral lists;</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span><span class="comment">// caller takes care of it.</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>func (sl *sweepLocked) sweep(preserve bool) bool {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s critical that we enter this function with preemption disabled,</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	<span class="comment">// GC must not start while we are in the middle of this function.</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	gp := getg()
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	if gp.m.locks == 0 &amp;&amp; gp.m.mallocing == 0 &amp;&amp; gp != gp.m.g0 {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		throw(&#34;mspan.sweep: m is not locked&#34;)
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	s := sl.mspan
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	if !preserve {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		<span class="comment">// We&#39;ll release ownership of this span. Nil it out to</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		<span class="comment">// prevent the caller from accidentally using it.</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		sl.mspan = nil
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	sweepgen := mheap_.sweepgen
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	if state := s.state.get(); state != mSpanInUse || s.sweepgen != sweepgen-1 {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		print(&#34;mspan.sweep: state=&#34;, state, &#34; sweepgen=&#34;, s.sweepgen, &#34; mheap.sweepgen=&#34;, sweepgen, &#34;\n&#34;)
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		throw(&#34;mspan.sweep: bad span state&#34;)
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	}
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	trace := traceAcquire()
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		trace.GCSweepSpan(s.npages * _PageSize)
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	}
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	mheap_.pagesSwept.Add(int64(s.npages))
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	spc := s.spanclass
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	size := s.elemsize
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	<span class="comment">// The allocBits indicate which unmarked objects don&#39;t need to be</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	<span class="comment">// processed since they were free at the end of the last GC cycle</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	<span class="comment">// and were not allocated since then.</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	<span class="comment">// If the allocBits index is &gt;= s.freeindex and the bit</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	<span class="comment">// is not marked then the object remains unallocated</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	<span class="comment">// since the last GC.</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	<span class="comment">// This situation is analogous to being on a freelist.</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	<span class="comment">// Unlink &amp; free special records for any objects we&#39;re about to free.</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	<span class="comment">// Two complications here:</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	<span class="comment">// 1. An object can have both finalizer and profile special records.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	<span class="comment">//    In such case we need to queue finalizer for execution,</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	<span class="comment">//    mark the object as live and preserve the profile special.</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	<span class="comment">// 2. A tiny object can have several finalizers setup for different offsets.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	<span class="comment">//    If such object is not marked, we need to queue all finalizers at once.</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	<span class="comment">// Both 1 and 2 are possible at the same time.</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	hadSpecials := s.specials != nil
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	siter := newSpecialsIter(s)
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	for siter.valid() {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		<span class="comment">// A finalizer can be set for an inner byte of an object, find object beginning.</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		objIndex := uintptr(siter.s.offset) / size
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		p := s.base() + objIndex*size
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		mbits := s.markBitsForIndex(objIndex)
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		if !mbits.isMarked() {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			<span class="comment">// This object is not marked and has at least one special record.</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>			<span class="comment">// Pass 1: see if it has at least one finalizer.</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>			hasFin := false
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			endOffset := p - s.base() + size
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>			for tmp := siter.s; tmp != nil &amp;&amp; uintptr(tmp.offset) &lt; endOffset; tmp = tmp.next {
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>				if tmp.kind == _KindSpecialFinalizer {
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>					<span class="comment">// Stop freeing of object if it has a finalizer.</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>					mbits.setMarkedNonAtomic()
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>					hasFin = true
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>					break
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>				}
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>			}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>			<span class="comment">// Pass 2: queue all finalizers _or_ handle profile record.</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>			for siter.valid() &amp;&amp; uintptr(siter.s.offset) &lt; endOffset {
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>				<span class="comment">// Find the exact byte for which the special was setup</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>				<span class="comment">// (as opposed to object beginning).</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>				special := siter.s
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>				p := s.base() + uintptr(special.offset)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>				if special.kind == _KindSpecialFinalizer || !hasFin {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>					siter.unlinkAndNext()
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>					freeSpecial(special, unsafe.Pointer(p), size)
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>				} else {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>					<span class="comment">// The object has finalizers, so we&#39;re keeping it alive.</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>					<span class="comment">// All other specials only apply when an object is freed,</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>					<span class="comment">// so just keep the special record.</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>					siter.next()
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>				}
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>			}
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		} else {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>			<span class="comment">// object is still live</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>			if siter.s.kind == _KindSpecialReachable {
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>				special := siter.unlinkAndNext()
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>				(*specialReachable)(unsafe.Pointer(special)).reachable = true
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>				freeSpecial(special, unsafe.Pointer(p), size)
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>			} else {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>				<span class="comment">// keep special record</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>				siter.next()
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>			}
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	if hadSpecials &amp;&amp; s.specials == nil {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		spanHasNoSpecials(s)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	if debug.allocfreetrace != 0 || debug.clobberfree != 0 || raceenabled || msanenabled || asanenabled {
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		<span class="comment">// Find all newly freed objects. This doesn&#39;t have to</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		<span class="comment">// efficient; allocfreetrace has massive overhead.</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		mbits := s.markBitsForBase()
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		abits := s.allocBitsForIndex(0)
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; uintptr(s.nelems); i++ {
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			if !mbits.isMarked() &amp;&amp; (abits.index &lt; uintptr(s.freeindex) || abits.isMarked()) {
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>				x := s.base() + i*s.elemsize
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>				if debug.allocfreetrace != 0 {
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>					tracefree(unsafe.Pointer(x), size)
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>				}
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>				if debug.clobberfree != 0 {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>					clobberfree(unsafe.Pointer(x), size)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>				}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>				<span class="comment">// User arenas are handled on explicit free.</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>				if raceenabled &amp;&amp; !s.isUserArenaChunk {
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>					racefree(unsafe.Pointer(x), size)
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>				}
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>				if msanenabled &amp;&amp; !s.isUserArenaChunk {
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>					msanfree(unsafe.Pointer(x), size)
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>				}
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>				if asanenabled &amp;&amp; !s.isUserArenaChunk {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>					asanpoison(unsafe.Pointer(x), size)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>				}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>			}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>			mbits.advance()
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>			abits.advance()
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	<span class="comment">// Check for zombie objects.</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	if s.freeindex &lt; s.nelems {
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		<span class="comment">// Everything &lt; freeindex is allocated and hence</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>		<span class="comment">// cannot be zombies.</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		<span class="comment">// Check the first bitmap byte, where we have to be</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		<span class="comment">// careful with freeindex.</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		obj := uintptr(s.freeindex)
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		if (*s.gcmarkBits.bytep(obj / 8)&amp;^*s.allocBits.bytep(obj / 8))&gt;&gt;(obj%8) != 0 {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>			s.reportZombies()
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		}
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		<span class="comment">// Check remaining bytes.</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		for i := obj/8 + 1; i &lt; divRoundUp(uintptr(s.nelems), 8); i++ {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>			if *s.gcmarkBits.bytep(i)&amp;^*s.allocBits.bytep(i) != 0 {
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>				s.reportZombies()
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>			}
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		}
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	<span class="comment">// Count the number of free objects in this span.</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	nalloc := uint16(s.countAlloc())
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	nfreed := s.allocCount - nalloc
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	if nalloc &gt; s.allocCount {
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		<span class="comment">// The zombie check above should have caught this in</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>		<span class="comment">// more detail.</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		print(&#34;runtime: nelems=&#34;, s.nelems, &#34; nalloc=&#34;, nalloc, &#34; previous allocCount=&#34;, s.allocCount, &#34; nfreed=&#34;, nfreed, &#34;\n&#34;)
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		throw(&#34;sweep increased allocation count&#34;)
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	}
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	s.allocCount = nalloc
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	s.freeindex = 0 <span class="comment">// reset allocation index to start of span.</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	s.freeIndexForScan = 0
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	if traceEnabled() {
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		getg().m.p.ptr().trace.reclaimed += uintptr(nfreed) * s.elemsize
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	}
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	<span class="comment">// gcmarkBits becomes the allocBits.</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	<span class="comment">// get a fresh cleared gcmarkBits in preparation for next GC</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	s.allocBits = s.gcmarkBits
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	s.gcmarkBits = newMarkBits(uintptr(s.nelems))
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	<span class="comment">// refresh pinnerBits if they exists</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	if s.pinnerBits != nil {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		s.refreshPinnerBits()
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	<span class="comment">// Initialize alloc bits cache.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	s.refillAllocCache(0)
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	<span class="comment">// The span must be in our exclusive ownership until we update sweepgen,</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	<span class="comment">// check for potential races.</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	if state := s.state.get(); state != mSpanInUse || s.sweepgen != sweepgen-1 {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		print(&#34;mspan.sweep: state=&#34;, state, &#34; sweepgen=&#34;, s.sweepgen, &#34; mheap.sweepgen=&#34;, sweepgen, &#34;\n&#34;)
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		throw(&#34;mspan.sweep: bad span state after sweep&#34;)
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	}
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	if s.sweepgen == sweepgen+1 || s.sweepgen == sweepgen+3 {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		throw(&#34;swept cached span&#34;)
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	<span class="comment">// We need to set s.sweepgen = h.sweepgen only when all blocks are swept,</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	<span class="comment">// because of the potential for a concurrent free/SetFinalizer.</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	<span class="comment">// But we need to set it before we make the span available for allocation</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	<span class="comment">// (return it to heap or mcentral), because allocation code assumes that a</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	<span class="comment">// span is already swept if available for allocation.</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	<span class="comment">// Serialization point.</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	<span class="comment">// At this point the mark bits are cleared and allocation ready</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	<span class="comment">// to go so release the span.</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	atomic.Store(&amp;s.sweepgen, sweepgen)
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	if s.isUserArenaChunk {
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		if preserve {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>			<span class="comment">// This is a case that should never be handled by a sweeper that</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>			<span class="comment">// preserves the span for reuse.</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>			throw(&#34;sweep: tried to preserve a user arena span&#34;)
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		}
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		if nalloc &gt; 0 {
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>			<span class="comment">// There still exist pointers into the span or the span hasn&#39;t been</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			<span class="comment">// freed yet. It&#39;s not ready to be reused. Put it back on the</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>			<span class="comment">// full swept list for the next cycle.</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>			mheap_.central[spc].mcentral.fullSwept(sweepgen).push(s)
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			return false
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>		}
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		<span class="comment">// It&#39;s only at this point that the sweeper doesn&#39;t actually need to look</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		<span class="comment">// at this arena anymore, so subtract from pagesInUse now.</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		mheap_.pagesInUse.Add(-s.npages)
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		s.state.set(mSpanDead)
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>		<span class="comment">// The arena is ready to be recycled. Remove it from the quarantine list</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		<span class="comment">// and place it on the ready list. Don&#39;t add it back to any sweep lists.</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>			<span class="comment">// It&#39;s the arena code&#39;s responsibility to get the chunk on the quarantine</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>			<span class="comment">// list by the time all references to the chunk are gone.</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>			if s.list != &amp;mheap_.userArena.quarantineList {
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>				throw(&#34;user arena span is on the wrong list&#34;)
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>			}
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>			lock(&amp;mheap_.lock)
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>			mheap_.userArena.quarantineList.remove(s)
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>			mheap_.userArena.readyList.insert(s)
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>			unlock(&amp;mheap_.lock)
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>		})
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		return false
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	}
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	if spc.sizeclass() != 0 {
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>		<span class="comment">// Handle spans for small objects.</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>		if nfreed &gt; 0 {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>			<span class="comment">// Only mark the span as needing zeroing if we&#39;ve freed any</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>			<span class="comment">// objects, because a fresh span that had been allocated into,</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>			<span class="comment">// wasn&#39;t totally filled, but then swept, still has all of its</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>			<span class="comment">// free slots zeroed.</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			s.needzero = 1
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>			stats := memstats.heapStats.acquire()
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>			atomic.Xadd64(&amp;stats.smallFreeCount[spc.sizeclass()], int64(nfreed))
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>			memstats.heapStats.release()
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>			<span class="comment">// Count the frees in the inconsistent, internal stats.</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>			gcController.totalFree.Add(int64(nfreed) * int64(s.elemsize))
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		}
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		if !preserve {
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>			<span class="comment">// The caller may not have removed this span from whatever</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>			<span class="comment">// unswept set its on but taken ownership of the span for</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>			<span class="comment">// sweeping by updating sweepgen. If this span still is in</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>			<span class="comment">// an unswept set, then the mcentral will pop it off the</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>			<span class="comment">// set, check its sweepgen, and ignore it.</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>			if nalloc == 0 {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>				<span class="comment">// Free totally free span directly back to the heap.</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>				mheap_.freeSpan(s)
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>				return true
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>			}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>			<span class="comment">// Return span back to the right mcentral list.</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>			if nalloc == s.nelems {
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>				mheap_.central[spc].mcentral.fullSwept(sweepgen).push(s)
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			} else {
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>				mheap_.central[spc].mcentral.partialSwept(sweepgen).push(s)
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>			}
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>		}
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	} else if !preserve {
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		<span class="comment">// Handle spans for large objects.</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		if nfreed != 0 {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			<span class="comment">// Free large object span to heap.</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>			<span class="comment">// NOTE(rsc,dvyukov): The original implementation of efence</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>			<span class="comment">// in CL 22060046 used sysFree instead of sysFault, so that</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>			<span class="comment">// the operating system would eventually give the memory</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>			<span class="comment">// back to us again, so that an efence program could run</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>			<span class="comment">// longer without running out of memory. Unfortunately,</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>			<span class="comment">// calling sysFree here without any kind of adjustment of the</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>			<span class="comment">// heap data structures means that when the memory does</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>			<span class="comment">// come back to us, we have the wrong metadata for it, either in</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>			<span class="comment">// the mspan structures or in the garbage collection bitmap.</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>			<span class="comment">// Using sysFault here means that the program will run out of</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>			<span class="comment">// memory fairly quickly in efence mode, but at least it won&#39;t</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>			<span class="comment">// have mysterious crashes due to confused memory reuse.</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>			<span class="comment">// It should be possible to switch back to sysFree if we also</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>			<span class="comment">// implement and then call some kind of mheap.deleteSpan.</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>			if debug.efence &gt; 0 {
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>				s.limit = 0 <span class="comment">// prevent mlookup from finding this span</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>				sysFault(unsafe.Pointer(s.base()), size)
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>			} else {
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>				mheap_.freeSpan(s)
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>			}
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>			if goexperiment.AllocHeaders &amp;&amp; s.largeType != nil &amp;&amp; s.largeType.TFlag&amp;abi.TFlagUnrolledBitmap != 0 {
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>				<span class="comment">// In the allocheaders experiment, the unrolled GCProg bitmap is allocated separately.</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>				<span class="comment">// Free the space for the unrolled bitmap.</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>				systemstack(func() {
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>					s := spanOf(uintptr(unsafe.Pointer(s.largeType)))
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>					mheap_.freeManual(s, spanAllocPtrScalarBits)
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>				})
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>				<span class="comment">// Make sure to zero this pointer without putting the old</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>				<span class="comment">// value in a write buffer, as the old value might be an</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>				<span class="comment">// invalid pointer. See arena.go:(*mheap).allocUserArenaChunk.</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>				*(*uintptr)(unsafe.Pointer(&amp;s.largeType)) = 0
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>			}
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>			<span class="comment">// Count the free in the consistent, external stats.</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>			stats := memstats.heapStats.acquire()
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>			atomic.Xadd64(&amp;stats.largeFreeCount, 1)
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>			atomic.Xadd64(&amp;stats.largeFree, int64(size))
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>			memstats.heapStats.release()
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>			<span class="comment">// Count the free in the inconsistent, internal stats.</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>			gcController.totalFree.Add(int64(size))
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>			return true
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		}
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		<span class="comment">// Add a large span directly onto the full+swept list.</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		mheap_.central[spc].mcentral.fullSwept(sweepgen).push(s)
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	}
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	return false
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>}
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span><span class="comment">// reportZombies reports any marked but free objects in s and throws.</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span><span class="comment">// This generally means one of the following:</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span><span class="comment">// 1. User code converted a pointer to a uintptr and then back</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span><span class="comment">// unsafely, and a GC ran while the uintptr was the only reference to</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span><span class="comment">// an object.</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span><span class="comment">// 2. User code (or a compiler bug) constructed a bad pointer that</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span><span class="comment">// points to a free slot, often a past-the-end pointer.</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span><span class="comment">// 3. The GC two cycles ago missed a pointer and freed a live object,</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span><span class="comment">// but it was still live in the last cycle, so this GC cycle found a</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span><span class="comment">// pointer to that object and marked it.</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>func (s *mspan) reportZombies() {
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	printlock()
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	print(&#34;runtime: marked free object in span &#34;, s, &#34;, elemsize=&#34;, s.elemsize, &#34; freeindex=&#34;, s.freeindex, &#34; (bad use of unsafe.Pointer? try -d=checkptr)\n&#34;)
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	mbits := s.markBitsForBase()
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	abits := s.allocBitsForIndex(0)
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; uintptr(s.nelems); i++ {
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>		addr := s.base() + i*s.elemsize
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>		print(hex(addr))
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		alloc := i &lt; uintptr(s.freeindex) || abits.isMarked()
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		if alloc {
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>			print(&#34; alloc&#34;)
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		} else {
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>			print(&#34; free &#34;)
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>		}
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		if mbits.isMarked() {
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>			print(&#34; marked  &#34;)
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		} else {
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>			print(&#34; unmarked&#34;)
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		}
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>		zombie := mbits.isMarked() &amp;&amp; !alloc
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		if zombie {
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>			print(&#34; zombie&#34;)
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>		}
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		if zombie {
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>			length := s.elemsize
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>			if length &gt; 1024 {
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>				length = 1024
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>			}
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>			hexdumpWords(addr, addr+length, nil)
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>		}
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>		mbits.advance()
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>		abits.advance()
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	}
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	throw(&#34;found pointer to free object&#34;)
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>}
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span><span class="comment">// deductSweepCredit deducts sweep credit for allocating a span of</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span><span class="comment">// size spanBytes. This must be performed *before* the span is</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span><span class="comment">// allocated to ensure the system has enough credit. If necessary, it</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span><span class="comment">// performs sweeping to prevent going in to debt. If the caller will</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span><span class="comment">// also sweep pages (e.g., for a large allocation), it can pass a</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span><span class="comment">// non-zero callerSweepPages to leave that many pages unswept.</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span><span class="comment">// deductSweepCredit makes a worst-case assumption that all spanBytes</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span><span class="comment">// bytes of the ultimately allocated span will be available for object</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span><span class="comment">// allocation.</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span><span class="comment">// deductSweepCredit is the core of the &#34;proportional sweep&#34; system.</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span><span class="comment">// It uses statistics gathered by the garbage collector to perform</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span><span class="comment">// enough sweeping so that all pages are swept during the concurrent</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span><span class="comment">// sweep phase between GC cycles.</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span><span class="comment">// mheap_ must NOT be locked.</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>func deductSweepCredit(spanBytes uintptr, callerSweepPages uintptr) {
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	if mheap_.sweepPagesPerByte == 0 {
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>		<span class="comment">// Proportional sweep is done or disabled.</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		return
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	}
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	trace := traceAcquire()
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		trace.GCSweepStart()
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	}
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	<span class="comment">// Fix debt if necessary.</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>retry:
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	sweptBasis := mheap_.pagesSweptBasis.Load()
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	live := gcController.heapLive.Load()
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	liveBasis := mheap_.sweepHeapLiveBasis
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	newHeapLive := spanBytes
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	if liveBasis &lt; live {
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>		<span class="comment">// Only do this subtraction when we don&#39;t overflow. Otherwise, pagesTarget</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>		<span class="comment">// might be computed as something really huge, causing us to get stuck</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>		<span class="comment">// sweeping here until the next mark phase.</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		<span class="comment">// Overflow can happen here if gcPaceSweeper is called concurrently with</span>
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>		<span class="comment">// sweeping (i.e. not during a STW, like it usually is) because this code</span>
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		<span class="comment">// is intentionally racy. A concurrent call to gcPaceSweeper can happen</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		<span class="comment">// if a GC tuning parameter is modified and we read an older value of</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>		<span class="comment">// heapLive than what was used to set the basis.</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		<span class="comment">// This state should be transient, so it&#39;s fine to just let newHeapLive</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>		<span class="comment">// be a relatively small number. We&#39;ll probably just skip this attempt to</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>		<span class="comment">// sweep.</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>		<span class="comment">// See issue #57523.</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>		newHeapLive += uintptr(live - liveBasis)
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	}
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	pagesTarget := int64(mheap_.sweepPagesPerByte*float64(newHeapLive)) - int64(callerSweepPages)
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>	for pagesTarget &gt; int64(mheap_.pagesSwept.Load()-sweptBasis) {
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>		if sweepone() == ^uintptr(0) {
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>			mheap_.sweepPagesPerByte = 0
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>			break
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>		}
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>		if mheap_.pagesSweptBasis.Load() != sweptBasis {
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>			<span class="comment">// Sweep pacing changed. Recompute debt.</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>			goto retry
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>		}
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	}
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	trace = traceAcquire()
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>	if trace.ok() {
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		trace.GCSweepDone()
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>		traceRelease(trace)
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>	}
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>}
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span><span class="comment">// clobberfree sets the memory content at x to bad content, for debugging</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span><span class="comment">// purposes.</span>
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>func clobberfree(x unsafe.Pointer, size uintptr) {
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>	<span class="comment">// size (span.elemsize) is always a multiple of 4.</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; size; i += 4 {
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>		*(*uint32)(add(x, i)) = 0xdeadbeef
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>	}
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>}
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span><span class="comment">// gcPaceSweeper updates the sweeper&#39;s pacing parameters.</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span><span class="comment">// Must be called whenever the GC&#39;s pacing is updated.</span>
<span id="L959" class="ln">   959&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span><span class="comment">// The world must be stopped, or mheap_.lock must be held.</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>func gcPaceSweeper(trigger uint64) {
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>	assertWorldStoppedOrLockHeld(&amp;mheap_.lock)
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	<span class="comment">// Update sweep pacing.</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	if isSweepDone() {
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		mheap_.sweepPagesPerByte = 0
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>	} else {
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>		<span class="comment">// Concurrent sweep needs to sweep all of the in-use</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>		<span class="comment">// pages by the time the allocated heap reaches the GC</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>		<span class="comment">// trigger. Compute the ratio of in-use pages to sweep</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>		<span class="comment">// per byte allocated, accounting for the fact that</span>
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>		<span class="comment">// some might already be swept.</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>		heapLiveBasis := gcController.heapLive.Load()
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>		heapDistance := int64(trigger) - int64(heapLiveBasis)
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>		<span class="comment">// Add a little margin so rounding errors and</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>		<span class="comment">// concurrent sweep are less likely to leave pages</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>		<span class="comment">// unswept when GC starts.</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>		heapDistance -= 1024 * 1024
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>		if heapDistance &lt; _PageSize {
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>			<span class="comment">// Avoid setting the sweep ratio extremely high</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>			heapDistance = _PageSize
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		}
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		pagesSwept := mheap_.pagesSwept.Load()
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>		pagesInUse := mheap_.pagesInUse.Load()
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>		sweepDistancePages := int64(pagesInUse) - int64(pagesSwept)
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>		if sweepDistancePages &lt;= 0 {
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>			mheap_.sweepPagesPerByte = 0
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>		} else {
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>			mheap_.sweepPagesPerByte = float64(sweepDistancePages) / float64(heapDistance)
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>			mheap_.sweepHeapLiveBasis = heapLiveBasis
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>			<span class="comment">// Write pagesSweptBasis last, since this</span>
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>			<span class="comment">// signals concurrent sweeps to recompute</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>			<span class="comment">// their debt.</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>			mheap_.pagesSweptBasis.Store(pagesSwept)
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		}
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	}
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>}
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>
</pre><p><a href="mgcsweep.go?m=text">View as plain text</a></p>

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
