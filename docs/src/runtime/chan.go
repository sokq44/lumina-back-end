<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/chan.go - Go Documentation Server</title>

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
<a href="chan.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">chan.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2014 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This file contains the implementation of Go channels.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// Invariants:</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//  At least one of c.sendq and c.recvq is empty,</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//  except for the case of an unbuffered channel with a single goroutine</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//  blocked on it for both sending and receiving using a select statement,</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//  in which case the length of c.sendq and c.recvq is limited only by the</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//  size of the select statement.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// For buffered channels, also:</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//  c.qcount &gt; 0 implies that c.recvq is empty.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//  c.qcount &lt; c.dataqsiz implies that c.sendq is empty.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>import (
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	&#34;runtime/internal/math&#34;
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>const (
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	maxAlign  = 8
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	hchanSize = unsafe.Sizeof(hchan{}) + uintptr(-int(unsafe.Sizeof(hchan{}))&amp;(maxAlign-1))
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	debugChan = false
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>)
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>type hchan struct {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	qcount   uint           <span class="comment">// total data in the queue</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	dataqsiz uint           <span class="comment">// size of the circular queue</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	buf      unsafe.Pointer <span class="comment">// points to an array of dataqsiz elements</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	elemsize uint16
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	closed   uint32
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	elemtype *_type <span class="comment">// element type</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	sendx    uint   <span class="comment">// send index</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	recvx    uint   <span class="comment">// receive index</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	recvq    waitq  <span class="comment">// list of recv waiters</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	sendq    waitq  <span class="comment">// list of send waiters</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// lock protects all fields in hchan, as well as several</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// fields in sudogs blocked on this channel.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// Do not change another G&#39;s status while holding this lock</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// (in particular, do not ready a G), as this can deadlock</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// with stack shrinking.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	lock mutex
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>type waitq struct {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	first *sudog
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	last  *sudog
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_makechan reflect.makechan</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func reflect_makechan(t *chantype, size int) *hchan {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	return makechan(t, size)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>func makechan64(t *chantype, size int64) *hchan {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if int64(int(size)) != size {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		panic(plainError(&#34;makechan: size out of range&#34;))
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	return makechan(t, int(size))
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func makechan(t *chantype, size int) *hchan {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	elem := t.Elem
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// compiler checks this but be safe.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	if elem.Size_ &gt;= 1&lt;&lt;16 {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		throw(&#34;makechan: invalid channel element type&#34;)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	if hchanSize%maxAlign != 0 || elem.Align_ &gt; maxAlign {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		throw(&#34;makechan: bad alignment&#34;)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	mem, overflow := math.MulUintptr(elem.Size_, uintptr(size))
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	if overflow || mem &gt; maxAlloc-hchanSize || size &lt; 0 {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		panic(plainError(&#34;makechan: size out of range&#34;))
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// buf points into the same allocation, elemtype is persistent.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// SudoG&#39;s are referenced from their owning thread so they can&#39;t be collected.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// TODO(dvyukov,rlh): Rethink when collector can move allocated objects.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	var c *hchan
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	switch {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	case mem == 0:
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		<span class="comment">// Queue or element size is zero.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		c = (*hchan)(mallocgc(hchanSize, nil, true))
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		<span class="comment">// Race detector uses this location for synchronization.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		c.buf = c.raceaddr()
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	case elem.PtrBytes == 0:
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		<span class="comment">// Elements do not contain pointers.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		<span class="comment">// Allocate hchan and buf in one call.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		c.buf = add(unsafe.Pointer(c), hchanSize)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	default:
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		<span class="comment">// Elements contain pointers.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		c = new(hchan)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		c.buf = mallocgc(mem, elem, true)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	c.elemsize = uint16(elem.Size_)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	c.elemtype = elem
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	c.dataqsiz = uint(size)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	lockInit(&amp;c.lock, lockRankHchan)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if debugChan {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		print(&#34;makechan: chan=&#34;, c, &#34;; elemsize=&#34;, elem.Size_, &#34;; dataqsiz=&#34;, size, &#34;\n&#34;)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	return c
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// chanbuf(c, i) is pointer to the i&#39;th slot in the buffer.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>func chanbuf(c *hchan, i uint) unsafe.Pointer {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	return add(c.buf, uintptr(i)*uintptr(c.elemsize))
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// full reports whether a send on c would block (that is, the channel is full).</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// It uses a single word-sized read of mutable state, so although</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// the answer is instantaneously true, the correct answer may have changed</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// by the time the calling function receives the return value.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>func full(c *hchan) bool {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// c.dataqsiz is immutable (never written after the channel is created)</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// so it is safe to read at any time during channel operation.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if c.dataqsiz == 0 {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// Assumes that a pointer read is relaxed-atomic.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		return c.recvq.first == nil
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// Assumes that a uint read is relaxed-atomic.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	return c.qcount == c.dataqsiz
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// entry point for c &lt;- x from compiled code.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>func chansend1(c *hchan, elem unsafe.Pointer) {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	chansend(c, elem, true, getcallerpc())
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">/*
<span id="L149" class="ln">   149&nbsp;&nbsp;</span> * generic single channel send/recv
<span id="L150" class="ln">   150&nbsp;&nbsp;</span> * If block is not nil,
<span id="L151" class="ln">   151&nbsp;&nbsp;</span> * then the protocol will not
<span id="L152" class="ln">   152&nbsp;&nbsp;</span> * sleep but return if it could
<span id="L153" class="ln">   153&nbsp;&nbsp;</span> * not complete.
<span id="L154" class="ln">   154&nbsp;&nbsp;</span> *
<span id="L155" class="ln">   155&nbsp;&nbsp;</span> * sleep can wake up with g.param == nil
<span id="L156" class="ln">   156&nbsp;&nbsp;</span> * when a channel involved in the sleep has
<span id="L157" class="ln">   157&nbsp;&nbsp;</span> * been closed.  it is easiest to loop and re-run
<span id="L158" class="ln">   158&nbsp;&nbsp;</span> * the operation; we&#39;ll see that it&#39;s now closed.
<span id="L159" class="ln">   159&nbsp;&nbsp;</span> */</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if c == nil {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		if !block {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			return false
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		gopark(nil, nil, waitReasonChanSendNilChan, traceBlockForever, 2)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		throw(&#34;unreachable&#34;)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	if debugChan {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		print(&#34;chansend: chan=&#34;, c, &#34;\n&#34;)
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	if raceenabled {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		racereadpc(c.raceaddr(), callerpc, abi.FuncPCABIInternal(chansend))
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// Fast path: check for failed non-blocking operation without acquiring the lock.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// After observing that the channel is not closed, we observe that the channel is</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// not ready for sending. Each of these observations is a single word-sized read</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// (first c.closed and second full()).</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// Because a closed channel cannot transition from &#39;ready for sending&#39; to</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// &#39;not ready for sending&#39;, even if the channel is closed between the two observations,</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// they imply a moment between the two when the channel was both not yet closed</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// and not ready for sending. We behave as if we observed the channel at that moment,</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// and report that the send cannot proceed.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">// It is okay if the reads are reordered here: if we observe that the channel is not</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// ready for sending and then observe that it is not closed, that implies that the</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// channel wasn&#39;t closed during the first observation. However, nothing here</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// guarantees forward progress. We rely on the side effects of lock release in</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// chanrecv() and closechan() to update this thread&#39;s view of c.closed and full().</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	if !block &amp;&amp; c.closed == 0 &amp;&amp; full(c) {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		return false
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	var t0 int64
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	if blockprofilerate &gt; 0 {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		t0 = cputicks()
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	lock(&amp;c.lock)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	if c.closed != 0 {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		unlock(&amp;c.lock)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		panic(plainError(&#34;send on closed channel&#34;))
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	if sg := c.recvq.dequeue(); sg != nil {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		<span class="comment">// Found a waiting receiver. We pass the value we want to send</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		<span class="comment">// directly to the receiver, bypassing the channel buffer (if any).</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		send(c, sg, ep, func() { unlock(&amp;c.lock) }, 3)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		return true
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	if c.qcount &lt; c.dataqsiz {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		<span class="comment">// Space is available in the channel buffer. Enqueue the element to send.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		qp := chanbuf(c, c.sendx)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		if raceenabled {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			racenotify(c, c.sendx, nil)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		typedmemmove(c.elemtype, qp, ep)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		c.sendx++
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		if c.sendx == c.dataqsiz {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			c.sendx = 0
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		c.qcount++
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		unlock(&amp;c.lock)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		return true
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	if !block {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		unlock(&amp;c.lock)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		return false
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// Block on the channel. Some receiver will complete our operation for us.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	gp := getg()
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	mysg := acquireSudog()
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	mysg.releasetime = 0
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	if t0 != 0 {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		mysg.releasetime = -1
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// No stack splits between assigning elem and enqueuing mysg</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// on gp.waiting where copystack can find it.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	mysg.elem = ep
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	mysg.waitlink = nil
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	mysg.g = gp
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	mysg.isSelect = false
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	mysg.c = c
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	gp.waiting = mysg
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	gp.param = nil
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	c.sendq.enqueue(mysg)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	<span class="comment">// Signal to anyone trying to shrink our stack that we&#39;re about</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	<span class="comment">// to park on a channel. The window between when this G&#39;s status</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// changes and when we set gp.activeStackChans is not safe for</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">// stack shrinking.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	gp.parkingOnChan.Store(true)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	gopark(chanparkcommit, unsafe.Pointer(&amp;c.lock), waitReasonChanSend, traceBlockChanSend, 2)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// Ensure the value being sent is kept alive until the</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// receiver copies it out. The sudog has a pointer to the</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// stack object, but sudogs aren&#39;t considered as roots of the</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// stack tracer.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	KeepAlive(ep)
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	<span class="comment">// someone woke us up.</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	if mysg != gp.waiting {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		throw(&#34;G waiting list is corrupted&#34;)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	gp.waiting = nil
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	gp.activeStackChans = false
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	closed := !mysg.success
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	gp.param = nil
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	if mysg.releasetime &gt; 0 {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		blockevent(mysg.releasetime-t0, 2)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	mysg.c = nil
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	releaseSudog(mysg)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	if closed {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		if c.closed == 0 {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			throw(&#34;chansend: spurious wakeup&#34;)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		panic(plainError(&#34;send on closed channel&#34;))
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	return true
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">// send processes a send operation on an empty channel c.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span><span class="comment">// The value ep sent by the sender is copied to the receiver sg.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">// The receiver is then woken up to go on its merry way.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span><span class="comment">// Channel c must be empty and locked.  send unlocks c with unlockf.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">// sg must already be dequeued from c.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span><span class="comment">// ep must be non-nil and point to the heap or the caller&#39;s stack.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	if raceenabled {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		if c.dataqsiz == 0 {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			racesync(c, sg)
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		} else {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			<span class="comment">// Pretend we go through the buffer, even though</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			<span class="comment">// we copy directly. Note that we need to increment</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			<span class="comment">// the head/tail locations only when raceenabled.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			racenotify(c, c.recvx, nil)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			racenotify(c, c.recvx, sg)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			c.recvx++
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			if c.recvx == c.dataqsiz {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>				c.recvx = 0
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>			c.sendx = c.recvx <span class="comment">// c.sendx = (c.sendx+1) % c.dataqsiz</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	}
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	if sg.elem != nil {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		sendDirect(c.elemtype, sg, ep)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		sg.elem = nil
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	gp := sg.g
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	unlockf()
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	gp.param = unsafe.Pointer(sg)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	sg.success = true
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	if sg.releasetime != 0 {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		sg.releasetime = cputicks()
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	goready(gp, skip+1)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span><span class="comment">// Sends and receives on unbuffered or empty-buffered channels are the</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">// only operations where one running goroutine writes to the stack of</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span><span class="comment">// another running goroutine. The GC assumes that stack writes only</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// happen when the goroutine is running and are only done by that</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// goroutine. Using a write barrier is sufficient to make up for</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">// violating that assumption, but the write barrier has to work.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span><span class="comment">// typedmemmove will call bulkBarrierPreWrite, but the target bytes</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// are not in the heap, so that will not help. We arrange to call</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">// memmove and typeBitsBulkBarrier instead.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>func sendDirect(t *_type, sg *sudog, src unsafe.Pointer) {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// src is on our stack, dst is a slot on another stack.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// Once we read sg.elem out of sg, it will no longer</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// be updated if the destination&#39;s stack gets copied (shrunk).</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">// So make sure that no preemption points can happen between read &amp; use.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	dst := sg.elem
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	typeBitsBulkBarrier(t, uintptr(dst), uintptr(src), t.Size_)
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	<span class="comment">// No need for cgo write barrier checks because dst is always</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	<span class="comment">// Go memory.</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	memmove(dst, src, t.Size_)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>func recvDirect(t *_type, sg *sudog, dst unsafe.Pointer) {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	<span class="comment">// dst is on our stack or the heap, src is on another stack.</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	<span class="comment">// The channel is locked, so src will not move during this</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	<span class="comment">// operation.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	src := sg.elem
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	typeBitsBulkBarrier(t, uintptr(dst), uintptr(src), t.Size_)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	memmove(dst, src, t.Size_)
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>}
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>func closechan(c *hchan) {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	if c == nil {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		panic(plainError(&#34;close of nil channel&#34;))
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	lock(&amp;c.lock)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	if c.closed != 0 {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		unlock(&amp;c.lock)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		panic(plainError(&#34;close of closed channel&#34;))
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	if raceenabled {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		racewritepc(c.raceaddr(), callerpc, abi.FuncPCABIInternal(closechan))
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		racerelease(c.raceaddr())
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	c.closed = 1
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	var glist gList
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	<span class="comment">// release all readers</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	for {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		sg := c.recvq.dequeue()
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		if sg == nil {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			break
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		}
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		if sg.elem != nil {
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>			typedmemclr(c.elemtype, sg.elem)
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			sg.elem = nil
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		if sg.releasetime != 0 {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			sg.releasetime = cputicks()
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		}
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		gp := sg.g
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		gp.param = unsafe.Pointer(sg)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		sg.success = false
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		if raceenabled {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			raceacquireg(gp, c.raceaddr())
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		glist.push(gp)
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	<span class="comment">// release all writers (they will panic)</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	for {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		sg := c.sendq.dequeue()
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		if sg == nil {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			break
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		sg.elem = nil
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		if sg.releasetime != 0 {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			sg.releasetime = cputicks()
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		gp := sg.g
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		gp.param = unsafe.Pointer(sg)
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		sg.success = false
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		if raceenabled {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			raceacquireg(gp, c.raceaddr())
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		glist.push(gp)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	unlock(&amp;c.lock)
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	<span class="comment">// Ready all Gs now that we&#39;ve dropped the channel lock.</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	for !glist.empty() {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		gp := glist.pop()
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		gp.schedlink = 0
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		goready(gp, 3)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span><span class="comment">// empty reports whether a read from c would block (that is, the channel is</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span><span class="comment">// empty).  It uses a single atomic read of mutable state.</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>func empty(c *hchan) bool {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	<span class="comment">// c.dataqsiz is immutable.</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	if c.dataqsiz == 0 {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		return atomic.Loadp(unsafe.Pointer(&amp;c.sendq.first)) == nil
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	return atomic.Loaduint(&amp;c.qcount) == 0
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>}
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span><span class="comment">// entry points for &lt;- c from compiled code.</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>func chanrecv1(c *hchan, elem unsafe.Pointer) {
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	chanrecv(c, elem, true)
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>func chanrecv2(c *hchan, elem unsafe.Pointer) (received bool) {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	_, received = chanrecv(c, elem, true)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	return
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span><span class="comment">// chanrecv receives on channel c and writes the received data to ep.</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span><span class="comment">// ep may be nil, in which case received data is ignored.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span><span class="comment">// If block == false and no elements are available, returns (false, false).</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span><span class="comment">// Otherwise, if c is closed, zeros *ep and returns (true, false).</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span><span class="comment">// Otherwise, fills in *ep with an element and returns (true, true).</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span><span class="comment">// A non-nil ep must point to the heap or the caller&#39;s stack.</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	<span class="comment">// raceenabled: don&#39;t need to check ep, as it is always on the stack</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	<span class="comment">// or is new memory allocated by reflect.</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	if debugChan {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		print(&#34;chanrecv: chan=&#34;, c, &#34;\n&#34;)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	if c == nil {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		if !block {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>			return
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		gopark(nil, nil, waitReasonChanReceiveNilChan, traceBlockForever, 2)
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		throw(&#34;unreachable&#34;)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	<span class="comment">// Fast path: check for failed non-blocking operation without acquiring the lock.</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	if !block &amp;&amp; empty(c) {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		<span class="comment">// After observing that the channel is not ready for receiving, we observe whether the</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		<span class="comment">// channel is closed.</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		<span class="comment">// Reordering of these checks could lead to incorrect behavior when racing with a close.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		<span class="comment">// For example, if the channel was open and not empty, was closed, and then drained,</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		<span class="comment">// reordered reads could incorrectly indicate &#34;open and empty&#34;. To prevent reordering,</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		<span class="comment">// we use atomic loads for both checks, and rely on emptying and closing to happen in</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		<span class="comment">// separate critical sections under the same lock.  This assumption fails when closing</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		<span class="comment">// an unbuffered channel with a blocked send, but that is an error condition anyway.</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		if atomic.Load(&amp;c.closed) == 0 {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			<span class="comment">// Because a channel cannot be reopened, the later observation of the channel</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			<span class="comment">// being not closed implies that it was also not closed at the moment of the</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>			<span class="comment">// first observation. We behave as if we observed the channel at that moment</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			<span class="comment">// and report that the receive cannot proceed.</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>			return
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		<span class="comment">// The channel is irreversibly closed. Re-check whether the channel has any pending data</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		<span class="comment">// to receive, which could have arrived between the empty and closed checks above.</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		<span class="comment">// Sequential consistency is also required here, when racing with such a send.</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		if empty(c) {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			<span class="comment">// The channel is irreversibly closed and empty.</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>			if raceenabled {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>				raceacquire(c.raceaddr())
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			if ep != nil {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>				typedmemclr(c.elemtype, ep)
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>			}
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			return true, false
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	var t0 int64
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	if blockprofilerate &gt; 0 {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		t0 = cputicks()
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	lock(&amp;c.lock)
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	if c.closed != 0 {
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		if c.qcount == 0 {
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>			if raceenabled {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>				raceacquire(c.raceaddr())
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			unlock(&amp;c.lock)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>			if ep != nil {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>				typedmemclr(c.elemtype, ep)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>			return true, false
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		}
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		<span class="comment">// The channel has been closed, but the channel&#39;s buffer have data.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	} else {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		<span class="comment">// Just found waiting sender with not closed.</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		if sg := c.sendq.dequeue(); sg != nil {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>			<span class="comment">// Found a waiting sender. If buffer is size 0, receive value</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			<span class="comment">// directly from sender. Otherwise, receive from head of queue</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			<span class="comment">// and add sender&#39;s value to the tail of the queue (both map to</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>			<span class="comment">// the same buffer slot because the queue is full).</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>			recv(c, sg, ep, func() { unlock(&amp;c.lock) }, 3)
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>			return true, true
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	if c.qcount &gt; 0 {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		<span class="comment">// Receive directly from queue</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		qp := chanbuf(c, c.recvx)
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		if raceenabled {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>			racenotify(c, c.recvx, nil)
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		}
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		if ep != nil {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>			typedmemmove(c.elemtype, ep, qp)
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		typedmemclr(c.elemtype, qp)
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		c.recvx++
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		if c.recvx == c.dataqsiz {
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>			c.recvx = 0
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		c.qcount--
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		unlock(&amp;c.lock)
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		return true, true
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	}
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	if !block {
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		unlock(&amp;c.lock)
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		return false, false
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	}
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	<span class="comment">// no sender available: block on this channel.</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	gp := getg()
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	mysg := acquireSudog()
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	mysg.releasetime = 0
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	if t0 != 0 {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		mysg.releasetime = -1
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	<span class="comment">// No stack splits between assigning elem and enqueuing mysg</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	<span class="comment">// on gp.waiting where copystack can find it.</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	mysg.elem = ep
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	mysg.waitlink = nil
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	gp.waiting = mysg
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	mysg.g = gp
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	mysg.isSelect = false
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	mysg.c = c
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	gp.param = nil
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	c.recvq.enqueue(mysg)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	<span class="comment">// Signal to anyone trying to shrink our stack that we&#39;re about</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	<span class="comment">// to park on a channel. The window between when this G&#39;s status</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	<span class="comment">// changes and when we set gp.activeStackChans is not safe for</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	<span class="comment">// stack shrinking.</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	gp.parkingOnChan.Store(true)
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	gopark(chanparkcommit, unsafe.Pointer(&amp;c.lock), waitReasonChanReceive, traceBlockChanRecv, 2)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	<span class="comment">// someone woke us up</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	if mysg != gp.waiting {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		throw(&#34;G waiting list is corrupted&#34;)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	gp.waiting = nil
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	gp.activeStackChans = false
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	if mysg.releasetime &gt; 0 {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		blockevent(mysg.releasetime-t0, 2)
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	success := mysg.success
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	gp.param = nil
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	mysg.c = nil
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	releaseSudog(mysg)
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	return true, success
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>}
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span><span class="comment">// recv processes a receive operation on a full channel c.</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span><span class="comment">// There are 2 parts:</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span><span class="comment">//  1. The value sent by the sender sg is put into the channel</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span><span class="comment">//     and the sender is woken up to go on its merry way.</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span><span class="comment">//  2. The value received by the receiver (the current G) is</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span><span class="comment">//     written to ep.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span><span class="comment">// For synchronous channels, both values are the same.</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span><span class="comment">// For asynchronous channels, the receiver gets its data from</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span><span class="comment">// the channel buffer and the sender&#39;s data is put in the</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span><span class="comment">// channel buffer.</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span><span class="comment">// Channel c must be full and locked. recv unlocks c with unlockf.</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span><span class="comment">// sg must already be dequeued from c.</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span><span class="comment">// A non-nil ep must point to the heap or the caller&#39;s stack.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>func recv(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	if c.dataqsiz == 0 {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		if raceenabled {
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>			racesync(c, sg)
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		}
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>		if ep != nil {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			<span class="comment">// copy data from sender</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>			recvDirect(c.elemtype, sg, ep)
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	} else {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		<span class="comment">// Queue is full. Take the item at the</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		<span class="comment">// head of the queue. Make the sender enqueue</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		<span class="comment">// its item at the tail of the queue. Since the</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		<span class="comment">// queue is full, those are both the same slot.</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		qp := chanbuf(c, c.recvx)
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		if raceenabled {
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>			racenotify(c, c.recvx, nil)
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			racenotify(c, c.recvx, sg)
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>		}
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		<span class="comment">// copy data from queue to receiver</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		if ep != nil {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>			typedmemmove(c.elemtype, ep, qp)
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		}
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		<span class="comment">// copy data from sender to queue</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		typedmemmove(c.elemtype, qp, sg.elem)
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		c.recvx++
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		if c.recvx == c.dataqsiz {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>			c.recvx = 0
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		c.sendx = c.recvx <span class="comment">// c.sendx = (c.sendx+1) % c.dataqsiz</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	}
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	sg.elem = nil
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	gp := sg.g
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	unlockf()
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	gp.param = unsafe.Pointer(sg)
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	sg.success = true
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	if sg.releasetime != 0 {
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		sg.releasetime = cputicks()
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	}
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	goready(gp, skip+1)
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>}
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>func chanparkcommit(gp *g, chanLock unsafe.Pointer) bool {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	<span class="comment">// There are unlocked sudogs that point into gp&#39;s stack. Stack</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	<span class="comment">// copying must lock the channels of those sudogs.</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	<span class="comment">// Set activeStackChans here instead of before we try parking</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	<span class="comment">// because we could self-deadlock in stack growth on the</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	<span class="comment">// channel lock.</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	gp.activeStackChans = true
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	<span class="comment">// Mark that it&#39;s safe for stack shrinking to occur now,</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	<span class="comment">// because any thread acquiring this G&#39;s stack for shrinking</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	<span class="comment">// is guaranteed to observe activeStackChans after this store.</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	gp.parkingOnChan.Store(false)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	<span class="comment">// Make sure we unlock after setting activeStackChans and</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	<span class="comment">// unsetting parkingOnChan. The moment we unlock chanLock</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	<span class="comment">// we risk gp getting readied by a channel operation and</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	<span class="comment">// so gp could continue running before everything before</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	<span class="comment">// the unlock is visible (even to gp itself).</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	unlock((*mutex)(chanLock))
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	return true
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>}
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span><span class="comment">// compiler implements</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span><span class="comment">//	select {</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span><span class="comment">//	case c &lt;- v:</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span><span class="comment">//		... foo</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span><span class="comment">//	default:</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span><span class="comment">//		... bar</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span><span class="comment">// as</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span><span class="comment">//	if selectnbsend(c, v) {</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span><span class="comment">//		... foo</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span><span class="comment">//	} else {</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span><span class="comment">//		... bar</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>func selectnbsend(c *hchan, elem unsafe.Pointer) (selected bool) {
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	return chansend(c, elem, false, getcallerpc())
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>}
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span><span class="comment">// compiler implements</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span><span class="comment">//	select {</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span><span class="comment">//	case v, ok = &lt;-c:</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span><span class="comment">//		... foo</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span><span class="comment">//	default:</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span><span class="comment">//		... bar</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span><span class="comment">// as</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span><span class="comment">//	if selected, ok = selectnbrecv(&amp;v, c); selected {</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span><span class="comment">//		... foo</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span><span class="comment">//	} else {</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span><span class="comment">//		... bar</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>func selectnbrecv(elem unsafe.Pointer, c *hchan) (selected, received bool) {
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	return chanrecv(c, elem, false)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_chansend reflect.chansend0</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>func reflect_chansend(c *hchan, elem unsafe.Pointer, nb bool) (selected bool) {
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	return chansend(c, elem, !nb, getcallerpc())
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>}
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_chanrecv reflect.chanrecv</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>func reflect_chanrecv(c *hchan, nb bool, elem unsafe.Pointer) (selected bool, received bool) {
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	return chanrecv(c, elem, !nb)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>}
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_chanlen reflect.chanlen</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>func reflect_chanlen(c *hchan) int {
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	if c == nil {
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		return 0
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	}
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	return int(c.qcount)
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>}
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span><span class="comment">//go:linkname reflectlite_chanlen internal/reflectlite.chanlen</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>func reflectlite_chanlen(c *hchan) int {
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	if c == nil {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		return 0
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	return int(c.qcount)
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>}
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_chancap reflect.chancap</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>func reflect_chancap(c *hchan) int {
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	if c == nil {
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>		return 0
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	}
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	return int(c.dataqsiz)
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>}
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_chanclose reflect.chanclose</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>func reflect_chanclose(c *hchan) {
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	closechan(c)
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>}
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>func (q *waitq) enqueue(sgp *sudog) {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	sgp.next = nil
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	x := q.last
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	if x == nil {
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		sgp.prev = nil
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		q.first = sgp
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		q.last = sgp
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		return
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	sgp.prev = x
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	x.next = sgp
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	q.last = sgp
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>}
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>func (q *waitq) dequeue() *sudog {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	for {
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>		sgp := q.first
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		if sgp == nil {
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>			return nil
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		y := sgp.next
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>		if y == nil {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>			q.first = nil
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>			q.last = nil
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		} else {
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>			y.prev = nil
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>			q.first = y
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>			sgp.next = nil <span class="comment">// mark as removed (see dequeueSudoG)</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>		}
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		<span class="comment">// if a goroutine was put on this queue because of a</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		<span class="comment">// select, there is a small window between the goroutine</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>		<span class="comment">// being woken up by a different case and it grabbing the</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>		<span class="comment">// channel locks. Once it has the lock</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>		<span class="comment">// it removes itself from the queue, so we won&#39;t see it after that.</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>		<span class="comment">// We use a flag in the G struct to tell us when someone</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		<span class="comment">// else has won the race to signal this goroutine but the goroutine</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>		<span class="comment">// hasn&#39;t removed itself from the queue yet.</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>		if sgp.isSelect &amp;&amp; !sgp.g.selectDone.CompareAndSwap(0, 1) {
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>			continue
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		}
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		return sgp
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>	}
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>}
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>func (c *hchan) raceaddr() unsafe.Pointer {
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	<span class="comment">// Treat read-like and write-like operations on the channel to</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	<span class="comment">// happen at this address. Avoid using the address of qcount</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	<span class="comment">// or dataqsiz, because the len() and cap() builtins read</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	<span class="comment">// those addresses, and we don&#39;t want them racing with</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	<span class="comment">// operations like close().</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	return unsafe.Pointer(&amp;c.buf)
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>}
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>func racesync(c *hchan, sg *sudog) {
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	racerelease(chanbuf(c, 0))
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	raceacquireg(sg.g, chanbuf(c, 0))
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	racereleaseg(sg.g, chanbuf(c, 0))
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	raceacquire(chanbuf(c, 0))
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>}
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span><span class="comment">// Notify the race detector of a send or receive involving buffer entry idx</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span><span class="comment">// and a channel c or its communicating partner sg.</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span><span class="comment">// This function handles the special case of c.elemsize==0.</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>func racenotify(c *hchan, idx uint, sg *sudog) {
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	<span class="comment">// We could have passed the unsafe.Pointer corresponding to entry idx</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	<span class="comment">// instead of idx itself.  However, in a future version of this function,</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	<span class="comment">// we can use idx to better handle the case of elemsize==0.</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	<span class="comment">// A future improvement to the detector is to call TSan with c and idx:</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	<span class="comment">// this way, Go will continue to not allocating buffer entries for channels</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	<span class="comment">// of elemsize==0, yet the race detector can be made to handle multiple</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>	<span class="comment">// sync objects underneath the hood (one sync object per idx)</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	qp := chanbuf(c, idx)
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	<span class="comment">// When elemsize==0, we don&#39;t allocate a full buffer for the channel.</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	<span class="comment">// Instead of individual buffer entries, the race detector uses the</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	<span class="comment">// c.buf as the only buffer entry.  This simplification prevents us from</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	<span class="comment">// following the memory model&#39;s happens-before rules (rules that are</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	<span class="comment">// implemented in racereleaseacquire).  Instead, we accumulate happens-before</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	<span class="comment">// information in the synchronization object associated with c.buf.</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>	if c.elemsize == 0 {
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		if sg == nil {
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>			raceacquire(qp)
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>			racerelease(qp)
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		} else {
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>			raceacquireg(sg.g, qp)
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>			racereleaseg(sg.g, qp)
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>		}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>	} else {
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>		if sg == nil {
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>			racereleaseacquire(qp)
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		} else {
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>			racereleaseacquireg(sg.g, qp)
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		}
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>}
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>
</pre><p><a href="chan.go?m=text">View as plain text</a></p>

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
