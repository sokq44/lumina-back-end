<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/trace2buf.go - Go Documentation Server</title>

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
<a href="trace2buf.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">trace2buf.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2023 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build goexperiment.exectracer2</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Trace buffer management.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package runtime
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// Maximum number of bytes required to encode uint64 in base-128.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>const traceBytesPerNumber = 10
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// traceWriter is the interface for writing all trace data.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// This type is passed around as a value, and all of its methods return</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// a new traceWriter. This allows for chaining together calls in a fluent-style</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// API. This is partly stylistic, and very slightly for performance, since</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// the compiler can destructure this value and pass it between calls as</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// just regular arguments. However, this style is not load-bearing, and</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// we can change it if it&#39;s deemed too error-prone.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>type traceWriter struct {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	traceLocker
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	*traceBuf
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// write returns an a traceWriter that writes into the current M&#39;s stream.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>func (tl traceLocker) writer() traceWriter {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	return traceWriter{traceLocker: tl, traceBuf: tl.mp.trace.buf[tl.gen%2]}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// unsafeTraceWriter produces a traceWriter that doesn&#39;t lock the trace.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// It should only be used in contexts where either:</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// - Another traceLocker is held.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// - trace.gen is prevented from advancing.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// buf may be nil.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>func unsafeTraceWriter(gen uintptr, buf *traceBuf) traceWriter {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	return traceWriter{traceLocker: traceLocker{gen: gen}, traceBuf: buf}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// end writes the buffer back into the m.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>func (w traceWriter) end() {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	if w.mp == nil {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		<span class="comment">// Tolerate a nil mp. It makes code that creates traceWriters directly</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		<span class="comment">// less error-prone.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		return
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	w.mp.trace.buf[w.gen%2] = w.traceBuf
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// ensure makes sure that at least maxSize bytes are available to write.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// Returns whether the buffer was flushed.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>func (w traceWriter) ensure(maxSize int) (traceWriter, bool) {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	refill := w.traceBuf == nil || !w.available(maxSize)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if refill {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		w = w.refill()
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	return w, refill
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// flush puts w.traceBuf on the queue of full buffers.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>func (w traceWriter) flush() traceWriter {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		lock(&amp;trace.lock)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		if w.traceBuf != nil {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			traceBufFlush(w.traceBuf, w.gen)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		unlock(&amp;trace.lock)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	})
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	w.traceBuf = nil
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	return w
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// refill puts w.traceBuf on the queue of full buffers and refresh&#39;s w&#39;s buffer.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func (w traceWriter) refill() traceWriter {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		lock(&amp;trace.lock)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		if w.traceBuf != nil {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			traceBufFlush(w.traceBuf, w.gen)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		if trace.empty != nil {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			w.traceBuf = trace.empty
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			trace.empty = w.traceBuf.link
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			unlock(&amp;trace.lock)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		} else {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			unlock(&amp;trace.lock)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			w.traceBuf = (*traceBuf)(sysAlloc(unsafe.Sizeof(traceBuf{}), &amp;memstats.other_sys))
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			if w.traceBuf == nil {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>				throw(&#34;trace: out of memory&#34;)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	})
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// Initialize the buffer.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	ts := traceClockNow()
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	if ts &lt;= w.traceBuf.lastTime {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		ts = w.traceBuf.lastTime + 1
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	w.traceBuf.lastTime = ts
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	w.traceBuf.link = nil
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	w.traceBuf.pos = 0
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// Tolerate a nil mp.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	mID := ^uint64(0)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	if w.mp != nil {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		mID = uint64(w.mp.procid)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// Write the buffer&#39;s header.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	w.byte(byte(traceEvEventBatch))
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	w.varint(uint64(w.gen))
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	w.varint(uint64(mID))
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	w.varint(uint64(ts))
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	w.traceBuf.lenPos = w.varintReserve()
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	return w
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// traceBufQueue is a FIFO of traceBufs.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>type traceBufQueue struct {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	head, tail *traceBuf
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// push queues buf into queue of buffers.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func (q *traceBufQueue) push(buf *traceBuf) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	buf.link = nil
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if q.head == nil {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		q.head = buf
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	} else {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		q.tail.link = buf
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	q.tail = buf
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// pop dequeues from the queue of buffers.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>func (q *traceBufQueue) pop() *traceBuf {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	buf := q.head
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	if buf == nil {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		return nil
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	q.head = buf.link
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	if q.head == nil {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		q.tail = nil
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	buf.link = nil
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	return buf
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>func (q *traceBufQueue) empty() bool {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	return q.head == nil
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// traceBufHeader is per-P tracing buffer.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>type traceBufHeader struct {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	link     *traceBuf <span class="comment">// in trace.empty/full</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	lastTime traceTime <span class="comment">// when we wrote the last event</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	pos      int       <span class="comment">// next write offset in arr</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	lenPos   int       <span class="comment">// position of batch length value</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// traceBuf is per-M tracing buffer.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// TODO(mknyszek): Rename traceBuf to traceBatch, since they map 1:1 with event batches.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>type traceBuf struct {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	traceBufHeader
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	arr [64&lt;&lt;10 - unsafe.Sizeof(traceBufHeader{})]byte <span class="comment">// underlying buffer for traceBufHeader.buf</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// byte appends v to buf.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>func (buf *traceBuf) byte(v byte) {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	buf.arr[buf.pos] = v
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	buf.pos++
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">// varint appends v to buf in little-endian-base-128 encoding.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>func (buf *traceBuf) varint(v uint64) {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	pos := buf.pos
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	arr := buf.arr[pos : pos+traceBytesPerNumber]
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	for i := range arr {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		if v &lt; 0x80 {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			pos += i + 1
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			arr[i] = byte(v)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			break
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		arr[i] = 0x80 | byte(v)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		v &gt;&gt;= 7
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	buf.pos = pos
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// varintReserve reserves enough space in buf to hold any varint.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// Space reserved this way can be filled in with the varintAt method.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>func (buf *traceBuf) varintReserve() int {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	p := buf.pos
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	buf.pos += traceBytesPerNumber
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	return p
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// stringData appends s&#39;s data directly to buf.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>func (buf *traceBuf) stringData(s string) {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	buf.pos += copy(buf.arr[buf.pos:], s)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>func (buf *traceBuf) available(size int) bool {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	return len(buf.arr)-buf.pos &gt;= size
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// varintAt writes varint v at byte position pos in buf. This always</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// consumes traceBytesPerNumber bytes. This is intended for when the caller</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// needs to reserve space for a varint but can&#39;t populate it until later.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">// Use varintReserve to reserve this space.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>func (buf *traceBuf) varintAt(pos int, v uint64) {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	for i := 0; i &lt; traceBytesPerNumber; i++ {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		if i &lt; traceBytesPerNumber-1 {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			buf.arr[pos] = 0x80 | byte(v)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		} else {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			buf.arr[pos] = byte(v)
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		v &gt;&gt;= 7
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		pos++
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	if v != 0 {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		throw(&#34;v could not fit in traceBytesPerNumber&#34;)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">// traceBufFlush flushes a trace buffer.</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">// Must run on the system stack because trace.lock must be held.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>func traceBufFlush(buf *traceBuf, gen uintptr) {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	assertLockHeld(&amp;trace.lock)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// Write out the non-header length of the batch in the header.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// Note: the length of the header is not included to make it easier</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">// to calculate this value when deserializing and reserializing the</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	<span class="comment">// trace. Varints can have additional padding of zero bits that is</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// quite difficult to preserve, and if we include the header we</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	<span class="comment">// force serializers to do more work. Nothing else actually needs</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// padding.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	buf.varintAt(buf.lenPos, uint64(buf.pos-(buf.lenPos+traceBytesPerNumber)))
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	trace.full[gen%2].push(buf)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	<span class="comment">// Notify the scheduler that there&#39;s work available and that the trace</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	<span class="comment">// reader should be scheduled.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	if !trace.workAvailable.Load() {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		trace.workAvailable.Store(true)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
</pre><p><a href="trace2buf.go?m=text">View as plain text</a></p>

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
