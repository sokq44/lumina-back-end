<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/sync/poolqueue.go - Go Documentation Server</title>

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
<a href="poolqueue.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/sync">sync</a>/<span class="text-muted">poolqueue.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/sync">sync</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2019 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package sync
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;sync/atomic&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// poolDequeue is a lock-free fixed-size single-producer,</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// multi-consumer queue. The single producer can both push and pop</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// from the head, and consumers can pop from the tail.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// It has the added feature that it nils out unused slots to avoid</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// unnecessary retention of objects. This is important for sync.Pool,</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// but not typically a property considered in the literature.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>type poolDequeue struct {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// headTail packs together a 32-bit head index and a 32-bit</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// tail index. Both are indexes into vals modulo len(vals)-1.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// tail = index of oldest data in queue</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// head = index of next slot to fill</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// Slots in the range [tail, head) are owned by consumers.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// A consumer continues to own a slot outside this range until</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// it nils the slot, at which point ownership passes to the</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// producer.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// The head index is stored in the most-significant bits so</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// that we can atomically add to it and the overflow is</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// harmless.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	headTail atomic.Uint64
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// vals is a ring buffer of interface{} values stored in this</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// dequeue. The size of this must be a power of 2.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// vals[i].typ is nil if the slot is empty and non-nil</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// otherwise. A slot is still in use until *both* the tail</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// index has moved beyond it and typ has been set to nil. This</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// is set to nil atomically by the consumer and read</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// atomically by the producer.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	vals []eface
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>type eface struct {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	typ, val unsafe.Pointer
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>const dequeueBits = 32
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// dequeueLimit is the maximum size of a poolDequeue.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// This must be at most (1&lt;&lt;dequeueBits)/2 because detecting fullness</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// depends on wrapping around the ring buffer without wrapping around</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// the index. We divide by 4 so this fits in an int on 32-bit.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>const dequeueLimit = (1 &lt;&lt; dequeueBits) / 4
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// dequeueNil is used in poolDequeue to represent interface{}(nil).</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// Since we use nil to represent empty slots, we need a sentinel value</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// to represent nil.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>type dequeueNil *struct{}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>func (d *poolDequeue) unpack(ptrs uint64) (head, tail uint32) {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	const mask = 1&lt;&lt;dequeueBits - 1
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	head = uint32((ptrs &gt;&gt; dequeueBits) &amp; mask)
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	tail = uint32(ptrs &amp; mask)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	return
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func (d *poolDequeue) pack(head, tail uint32) uint64 {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	const mask = 1&lt;&lt;dequeueBits - 1
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return (uint64(head) &lt;&lt; dequeueBits) |
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		uint64(tail&amp;mask)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// pushHead adds val at the head of the queue. It returns false if the</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// queue is full. It must only be called by a single producer.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>func (d *poolDequeue) pushHead(val any) bool {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	ptrs := d.headTail.Load()
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	head, tail := d.unpack(ptrs)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	if (tail+uint32(len(d.vals)))&amp;(1&lt;&lt;dequeueBits-1) == head {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		<span class="comment">// Queue is full.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		return false
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	slot := &amp;d.vals[head&amp;uint32(len(d.vals)-1)]
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// Check if the head slot has been released by popTail.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	typ := atomic.LoadPointer(&amp;slot.typ)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	if typ != nil {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		<span class="comment">// Another goroutine is still cleaning up the tail, so</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		<span class="comment">// the queue is actually still full.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		return false
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// The head slot is free, so we own it.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	if val == nil {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		val = dequeueNil(nil)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	*(*any)(unsafe.Pointer(slot)) = val
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// Increment head. This passes ownership of slot to popTail</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// and acts as a store barrier for writing the slot.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	d.headTail.Add(1 &lt;&lt; dequeueBits)
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	return true
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// popHead removes and returns the element at the head of the queue.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// It returns false if the queue is empty. It must only be called by a</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// single producer.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>func (d *poolDequeue) popHead() (any, bool) {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	var slot *eface
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	for {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		ptrs := d.headTail.Load()
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		head, tail := d.unpack(ptrs)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		if tail == head {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			<span class="comment">// Queue is empty.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			return nil, false
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		<span class="comment">// Confirm tail and decrement head. We do this before</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		<span class="comment">// reading the value to take back ownership of this</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		<span class="comment">// slot.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		head--
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		ptrs2 := d.pack(head, tail)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			<span class="comment">// We successfully took back slot.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			slot = &amp;d.vals[head&amp;uint32(len(d.vals)-1)]
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			break
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	val := *(*any)(unsafe.Pointer(slot))
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	if val == dequeueNil(nil) {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		val = nil
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// Zero the slot. Unlike popTail, this isn&#39;t racing with</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// pushHead, so we don&#39;t need to be careful here.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	*slot = eface{}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	return val, true
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// popTail removes and returns the element at the tail of the queue.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// It returns false if the queue is empty. It may be called by any</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// number of consumers.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>func (d *poolDequeue) popTail() (any, bool) {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	var slot *eface
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	for {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		ptrs := d.headTail.Load()
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		head, tail := d.unpack(ptrs)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		if tail == head {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			<span class="comment">// Queue is empty.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			return nil, false
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// Confirm head and tail (for our speculative check</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// above) and increment tail. If this succeeds, then</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		<span class="comment">// we own the slot at tail.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		ptrs2 := d.pack(head, tail+1)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			<span class="comment">// Success.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			slot = &amp;d.vals[tail&amp;uint32(len(d.vals)-1)]
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			break
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// We now own slot.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	val := *(*any)(unsafe.Pointer(slot))
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if val == dequeueNil(nil) {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		val = nil
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// Tell pushHead that we&#39;re done with this slot. Zeroing the</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// slot is also important so we don&#39;t leave behind references</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// that could keep this object live longer than necessary.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// We write to val first and then publish that we&#39;re done with</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// this slot by atomically writing to typ.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	slot.val = nil
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	atomic.StorePointer(&amp;slot.typ, nil)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// At this point pushHead owns the slot.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	return val, true
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// poolChain is a dynamically-sized version of poolDequeue.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// This is implemented as a doubly-linked list queue of poolDequeues</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// where each dequeue is double the size of the previous one. Once a</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// dequeue fills up, this allocates a new one and only ever pushes to</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// the latest dequeue. Pops happen from the other end of the list and</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">// once a dequeue is exhausted, it gets removed from the list.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>type poolChain struct {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// head is the poolDequeue to push to. This is only accessed</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// by the producer, so doesn&#39;t need to be synchronized.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	head *poolChainElt
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// tail is the poolDequeue to popTail from. This is accessed</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// by consumers, so reads and writes must be atomic.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	tail *poolChainElt
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>type poolChainElt struct {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	poolDequeue
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// next and prev link to the adjacent poolChainElts in this</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// poolChain.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">// next is written atomically by the producer and read</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	<span class="comment">// atomically by the consumer. It only transitions from nil to</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// non-nil.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// prev is written atomically by the consumer and read</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// atomically by the producer. It only transitions from</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	<span class="comment">// non-nil to nil.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	next, prev *poolChainElt
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>func storePoolChainElt(pp **poolChainElt, v *poolChainElt) {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(pp)), unsafe.Pointer(v))
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>func loadPoolChainElt(pp **poolChainElt) *poolChainElt {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	return (*poolChainElt)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(pp))))
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>func (c *poolChain) pushHead(val any) {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	d := c.head
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	if d == nil {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		<span class="comment">// Initialize the chain.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		const initSize = 8 <span class="comment">// Must be a power of 2</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		d = new(poolChainElt)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		d.vals = make([]eface, initSize)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		c.head = d
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		storePoolChainElt(&amp;c.tail, d)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	if d.pushHead(val) {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		return
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// The current dequeue is full. Allocate a new one of twice</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// the size.</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	newSize := len(d.vals) * 2
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	if newSize &gt;= dequeueLimit {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		<span class="comment">// Can&#39;t make it any bigger.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		newSize = dequeueLimit
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	d2 := &amp;poolChainElt{prev: d}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	d2.vals = make([]eface, newSize)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	c.head = d2
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	storePoolChainElt(&amp;d.next, d2)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	d2.pushHead(val)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>func (c *poolChain) popHead() (any, bool) {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	d := c.head
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	for d != nil {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		if val, ok := d.popHead(); ok {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			return val, ok
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		<span class="comment">// There may still be unconsumed elements in the</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		<span class="comment">// previous dequeue, so try backing up.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		d = loadPoolChainElt(&amp;d.prev)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	return nil, false
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>func (c *poolChain) popTail() (any, bool) {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	d := loadPoolChainElt(&amp;c.tail)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	if d == nil {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		return nil, false
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	for {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		<span class="comment">// It&#39;s important that we load the next pointer</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		<span class="comment">// *before* popping the tail. In general, d may be</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		<span class="comment">// transiently empty, but if next is non-nil before</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		<span class="comment">// the pop and the pop fails, then d is permanently</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		<span class="comment">// empty, which is the only condition under which it&#39;s</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		<span class="comment">// safe to drop d from the chain.</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		d2 := loadPoolChainElt(&amp;d.next)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		if val, ok := d.popTail(); ok {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			return val, ok
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		if d2 == nil {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			<span class="comment">// This is the only dequeue. It&#39;s empty right</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			<span class="comment">// now, but could be pushed to in the future.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			return nil, false
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		<span class="comment">// The tail of the chain has been drained, so move on</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		<span class="comment">// to the next dequeue. Try to drop it from the chain</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		<span class="comment">// so the next pop doesn&#39;t have to look at the empty</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		<span class="comment">// dequeue again.</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&amp;c.tail)), unsafe.Pointer(d), unsafe.Pointer(d2)) {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			<span class="comment">// We won the race. Clear the prev pointer so</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			<span class="comment">// the garbage collector can collect the empty</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			<span class="comment">// dequeue and so popHead doesn&#39;t back up</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			<span class="comment">// further than necessary.</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			storePoolChainElt(&amp;d2.prev, nil)
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		d = d2
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
</pre><p><a href="poolqueue.go?m=text">View as plain text</a></p>

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
