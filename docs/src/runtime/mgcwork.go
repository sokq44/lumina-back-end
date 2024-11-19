<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mgcwork.go - Go Documentation Server</title>

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
<a href="mgcwork.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mgcwork.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>const (
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	_WorkbufSize = 2048 <span class="comment">// in bytes; larger values result in less contention</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">// workbufAlloc is the number of bytes to allocate at a time</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// for new workbufs. This must be a multiple of pageSize and</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// should be a multiple of _WorkbufSize.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// Larger values reduce workbuf allocation overhead. Smaller</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// values reduce heap fragmentation.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	workbufAlloc = 32 &lt;&lt; 10
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>func init() {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	if workbufAlloc%pageSize != 0 || workbufAlloc%_WorkbufSize != 0 {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		throw(&#34;bad workbufAlloc&#34;)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// Garbage collector work pool abstraction.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// This implements a producer/consumer model for pointers to grey</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// objects. A grey object is one that is marked and on a work</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// queue. A black object is marked and not on a work queue.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// Write barriers, root discovery, stack scanning, and object scanning</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// produce pointers to grey objects. Scanning consumes pointers to</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// grey objects, thus blackening them, and then scans them,</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// potentially producing new pointers to grey objects.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// A gcWork provides the interface to produce and consume work for the</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// garbage collector.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// A gcWork can be used on the stack as follows:</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//	(preemption must be disabled)</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//	gcw := &amp;getg().m.p.ptr().gcw</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//	.. call gcw.put() to produce and gcw.tryGet() to consume ..</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// It&#39;s important that any use of gcWork during the mark phase prevent</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// the garbage collector from transitioning to mark termination since</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// gcWork may locally hold GC work buffers. This can be done by</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// disabling preemption (systemstack or acquirem).</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>type gcWork struct {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// wbuf1 and wbuf2 are the primary and secondary work buffers.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// This can be thought of as a stack of both work buffers&#39;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// pointers concatenated. When we pop the last pointer, we</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// shift the stack up by one work buffer by bringing in a new</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// full buffer and discarding an empty one. When we fill both</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// buffers, we shift the stack down by one work buffer by</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// bringing in a new empty buffer and discarding a full one.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// This way we have one buffer&#39;s worth of hysteresis, which</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// amortizes the cost of getting or putting a work buffer over</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// at least one buffer of work and reduces contention on the</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// global work lists.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// wbuf1 is always the buffer we&#39;re currently pushing to and</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// popping from and wbuf2 is the buffer that will be discarded</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// next.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// Invariant: Both wbuf1 and wbuf2 are nil or neither are.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	wbuf1, wbuf2 *workbuf
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// Bytes marked (blackened) on this gcWork. This is aggregated</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// into work.bytesMarked by dispose.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	bytesMarked uint64
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// Heap scan work performed on this gcWork. This is aggregated into</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// gcController by dispose and may also be flushed by callers.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// Other types of scan work are flushed immediately.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	heapScanWork int64
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// flushedWork indicates that a non-empty work buffer was</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// flushed to the global work list since the last gcMarkDone</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// termination check. Specifically, this indicates that this</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// gcWork may have communicated work to another gcWork.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	flushedWork bool
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// Most of the methods of gcWork are go:nowritebarrierrec because the</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// write barrier itself can invoke gcWork methods but the methods are</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// not generally re-entrant. Hence, if a gcWork method invoked the</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// write barrier while the gcWork was in an inconsistent state, and</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// the write barrier in turn invoked a gcWork method, it could</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// permanently corrupt the gcWork.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>func (w *gcWork) init() {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	w.wbuf1 = getempty()
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	wbuf2 := trygetfull()
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	if wbuf2 == nil {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		wbuf2 = getempty()
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	w.wbuf2 = wbuf2
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// put enqueues a pointer for the garbage collector to trace.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// obj must point to the beginning of a heap object or an oblet.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func (w *gcWork) put(obj uintptr) {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	flushed := false
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	wbuf := w.wbuf1
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// Record that this may acquire the wbufSpans or heap lock to</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// allocate a workbuf.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	lockWithRankMayAcquire(&amp;work.wbufSpans.lock, lockRankWbufSpans)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	lockWithRankMayAcquire(&amp;mheap_.lock, lockRankMheap)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	if wbuf == nil {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		w.init()
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		wbuf = w.wbuf1
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		<span class="comment">// wbuf is empty at this point.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	} else if wbuf.nobj == len(wbuf.obj) {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		w.wbuf1, w.wbuf2 = w.wbuf2, w.wbuf1
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		wbuf = w.wbuf1
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		if wbuf.nobj == len(wbuf.obj) {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			putfull(wbuf)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			w.flushedWork = true
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			wbuf = getempty()
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			w.wbuf1 = wbuf
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			flushed = true
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	wbuf.obj[wbuf.nobj] = obj
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	wbuf.nobj++
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// If we put a buffer on full, let the GC controller know so</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// it can encourage more workers to run. We delay this until</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// the end of put so that w is in a consistent state, since</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// enlistWorker may itself manipulate w.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	if flushed &amp;&amp; gcphase == _GCmark {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		gcController.enlistWorker()
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// putFast does a put and reports whether it can be done quickly</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// otherwise it returns false and the caller needs to call put.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func (w *gcWork) putFast(obj uintptr) bool {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	wbuf := w.wbuf1
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	if wbuf == nil || wbuf.nobj == len(wbuf.obj) {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		return false
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	wbuf.obj[wbuf.nobj] = obj
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	wbuf.nobj++
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	return true
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// putBatch performs a put on every pointer in obj. See put for</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// constraints on these pointers.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>func (w *gcWork) putBatch(obj []uintptr) {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	if len(obj) == 0 {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		return
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	flushed := false
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	wbuf := w.wbuf1
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	if wbuf == nil {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		w.init()
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		wbuf = w.wbuf1
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	for len(obj) &gt; 0 {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		for wbuf.nobj == len(wbuf.obj) {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			putfull(wbuf)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			w.flushedWork = true
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			w.wbuf1, w.wbuf2 = w.wbuf2, getempty()
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			wbuf = w.wbuf1
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			flushed = true
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		n := copy(wbuf.obj[wbuf.nobj:], obj)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		wbuf.nobj += n
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		obj = obj[n:]
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	if flushed &amp;&amp; gcphase == _GCmark {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		gcController.enlistWorker()
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// tryGet dequeues a pointer for the garbage collector to trace.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// If there are no pointers remaining in this gcWork or in the global</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// queue, tryGet returns 0.  Note that there may still be pointers in</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// other gcWork instances or other caches.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>func (w *gcWork) tryGet() uintptr {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	wbuf := w.wbuf1
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	if wbuf == nil {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		w.init()
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		wbuf = w.wbuf1
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		<span class="comment">// wbuf is empty at this point.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if wbuf.nobj == 0 {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		w.wbuf1, w.wbuf2 = w.wbuf2, w.wbuf1
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		wbuf = w.wbuf1
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		if wbuf.nobj == 0 {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			owbuf := wbuf
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			wbuf = trygetfull()
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			if wbuf == nil {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>				return 0
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			putempty(owbuf)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			w.wbuf1 = wbuf
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	wbuf.nobj--
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	return wbuf.obj[wbuf.nobj]
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// tryGetFast dequeues a pointer for the garbage collector to trace</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// if one is readily available. Otherwise it returns 0 and</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// the caller is expected to call tryGet().</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>func (w *gcWork) tryGetFast() uintptr {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	wbuf := w.wbuf1
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	if wbuf == nil || wbuf.nobj == 0 {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		return 0
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	wbuf.nobj--
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	return wbuf.obj[wbuf.nobj]
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">// dispose returns any cached pointers to the global queue.</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// The buffers are being put on the full queue so that the</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">// write barriers will not simply reacquire them before the</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// GC can inspect them. This helps reduce the mutator&#39;s</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// ability to hide pointers during the concurrent mark phase.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func (w *gcWork) dispose() {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	if wbuf := w.wbuf1; wbuf != nil {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		if wbuf.nobj == 0 {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			putempty(wbuf)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		} else {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			putfull(wbuf)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			w.flushedWork = true
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		w.wbuf1 = nil
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		wbuf = w.wbuf2
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		if wbuf.nobj == 0 {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			putempty(wbuf)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		} else {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			putfull(wbuf)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			w.flushedWork = true
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		w.wbuf2 = nil
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	if w.bytesMarked != 0 {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		<span class="comment">// dispose happens relatively infrequently. If this</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		<span class="comment">// atomic becomes a problem, we should first try to</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		<span class="comment">// dispose less and if necessary aggregate in a per-P</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		<span class="comment">// counter.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		atomic.Xadd64(&amp;work.bytesMarked, int64(w.bytesMarked))
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		w.bytesMarked = 0
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	if w.heapScanWork != 0 {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		gcController.heapScanWork.Add(w.heapScanWork)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		w.heapScanWork = 0
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">// balance moves some work that&#39;s cached in this gcWork back on the</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">// global queue.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>func (w *gcWork) balance() {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	if w.wbuf1 == nil {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		return
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	if wbuf := w.wbuf2; wbuf.nobj != 0 {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		putfull(wbuf)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		w.flushedWork = true
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		w.wbuf2 = getempty()
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	} else if wbuf := w.wbuf1; wbuf.nobj &gt; 4 {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		w.wbuf1 = handoff(wbuf)
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		w.flushedWork = true <span class="comment">// handoff did putfull</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	} else {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		return
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	<span class="comment">// We flushed a buffer to the full list, so wake a worker.</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	if gcphase == _GCmark {
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		gcController.enlistWorker()
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// empty reports whether w has no mark work available.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>func (w *gcWork) empty() bool {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	return w.wbuf1 == nil || (w.wbuf1.nobj == 0 &amp;&amp; w.wbuf2.nobj == 0)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// Internally, the GC work pool is kept in arrays in work buffers.</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">// The gcWork interface caches a work buffer until full (or empty) to</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">// avoid contending on the global work buffer lists.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>type workbufhdr struct {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	node lfnode <span class="comment">// must be first</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	nobj int
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>type workbuf struct {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	workbufhdr
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// account for the above fields</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	obj [(_WorkbufSize - unsafe.Sizeof(workbufhdr{})) / goarch.PtrSize]uintptr
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span><span class="comment">// workbuf factory routines. These funcs are used to manage the</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// workbufs.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">// If the GC asks for some work these are the only routines that</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">// make wbufs available to the GC.</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func (b *workbuf) checknonempty() {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	if b.nobj == 0 {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		throw(&#34;workbuf is empty&#34;)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>func (b *workbuf) checkempty() {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	if b.nobj != 0 {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		throw(&#34;workbuf is not empty&#34;)
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// getempty pops an empty work buffer off the work.empty list,</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// allocating new buffers if none are available.</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>func getempty() *workbuf {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	var b *workbuf
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	if work.empty != 0 {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		b = (*workbuf)(work.empty.pop())
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		if b != nil {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			b.checkempty()
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	<span class="comment">// Record that this may acquire the wbufSpans or heap lock to</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	<span class="comment">// allocate a workbuf.</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	lockWithRankMayAcquire(&amp;work.wbufSpans.lock, lockRankWbufSpans)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	lockWithRankMayAcquire(&amp;mheap_.lock, lockRankMheap)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	if b == nil {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		<span class="comment">// Allocate more workbufs.</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		var s *mspan
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		if work.wbufSpans.free.first != nil {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			lock(&amp;work.wbufSpans.lock)
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			s = work.wbufSpans.free.first
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			if s != nil {
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>				work.wbufSpans.free.remove(s)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>				work.wbufSpans.busy.insert(s)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			unlock(&amp;work.wbufSpans.lock)
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		}
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		if s == nil {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			systemstack(func() {
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>				s = mheap_.allocManual(workbufAlloc/pageSize, spanAllocWorkBuf)
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			})
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			if s == nil {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>				throw(&#34;out of memory&#34;)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			<span class="comment">// Record the new span in the busy list.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			lock(&amp;work.wbufSpans.lock)
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>			work.wbufSpans.busy.insert(s)
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			unlock(&amp;work.wbufSpans.lock)
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		<span class="comment">// Slice up the span into new workbufs. Return one and</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		<span class="comment">// put the rest on the empty list.</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		for i := uintptr(0); i+_WorkbufSize &lt;= workbufAlloc; i += _WorkbufSize {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			newb := (*workbuf)(unsafe.Pointer(s.base() + i))
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			newb.nobj = 0
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			lfnodeValidate(&amp;newb.node)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			if i == 0 {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>				b = newb
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			} else {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>				putempty(newb)
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	return b
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span><span class="comment">// putempty puts a workbuf onto the work.empty list.</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span><span class="comment">// Upon entry this goroutine owns b. The lfstack.push relinquishes ownership.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>func putempty(b *workbuf) {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	b.checkempty()
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	work.empty.push(&amp;b.node)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>}
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span><span class="comment">// putfull puts the workbuf on the work.full list for the GC.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span><span class="comment">// putfull accepts partially full buffers so the GC can avoid competing</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span><span class="comment">// with the mutators for ownership of partially full buffers.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>func putfull(b *workbuf) {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	b.checknonempty()
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	work.full.push(&amp;b.node)
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span><span class="comment">// trygetfull tries to get a full or partially empty workbuffer.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span><span class="comment">// If one is not immediately available return nil.</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>func trygetfull() *workbuf {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	b := (*workbuf)(work.full.pop())
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	if b != nil {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		b.checknonempty()
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		return b
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	return b
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>func handoff(b *workbuf) *workbuf {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	<span class="comment">// Make new buffer with half of b&#39;s pointers.</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	b1 := getempty()
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	n := b.nobj / 2
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	b.nobj -= n
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	b1.nobj = n
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	memmove(unsafe.Pointer(&amp;b1.obj[0]), unsafe.Pointer(&amp;b.obj[b.nobj]), uintptr(n)*unsafe.Sizeof(b1.obj[0]))
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	<span class="comment">// Put b on full list - let first half of b get stolen.</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	putfull(b)
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	return b1
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span><span class="comment">// prepareFreeWorkbufs moves busy workbuf spans to free list so they</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span><span class="comment">// can be freed to the heap. This must only be called when all</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span><span class="comment">// workbufs are on the empty list.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>func prepareFreeWorkbufs() {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	lock(&amp;work.wbufSpans.lock)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	if work.full != 0 {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		throw(&#34;cannot free workbufs when work.full != 0&#34;)
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	<span class="comment">// Since all workbufs are on the empty list, we don&#39;t care</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	<span class="comment">// which ones are in which spans. We can wipe the entire empty</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	<span class="comment">// list and move all workbuf spans to the free list.</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	work.empty = 0
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	work.wbufSpans.free.takeAll(&amp;work.wbufSpans.busy)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	unlock(&amp;work.wbufSpans.lock)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>}
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span><span class="comment">// freeSomeWbufs frees some workbufs back to the heap and returns</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span><span class="comment">// true if it should be called again to free more.</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>func freeSomeWbufs(preemptible bool) bool {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	const batchSize = 64 <span class="comment">// ~1–2 µs per span.</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	lock(&amp;work.wbufSpans.lock)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	if gcphase != _GCoff || work.wbufSpans.free.isEmpty() {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		unlock(&amp;work.wbufSpans.lock)
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		return false
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		gp := getg().m.curg
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		for i := 0; i &lt; batchSize &amp;&amp; !(preemptible &amp;&amp; gp.preempt); i++ {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			span := work.wbufSpans.free.first
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			if span == nil {
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>				break
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			work.wbufSpans.free.remove(span)
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			mheap_.freeManual(span, spanAllocWorkBuf)
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	})
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	more := !work.wbufSpans.free.isEmpty()
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	unlock(&amp;work.wbufSpans.lock)
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	return more
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>
</pre><p><a href="mgcwork.go?m=text">View as plain text</a></p>

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
