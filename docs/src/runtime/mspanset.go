<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mspanset.go - Go Documentation Server</title>

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
<a href="mspanset.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mspanset.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/cpu&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// A spanSet is a set of *mspans.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// spanSet is safe for concurrent push and pop operations.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>type spanSet struct {
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// A spanSet is a two-level data structure consisting of a</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// growable spine that points to fixed-sized blocks. The spine</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// can be accessed without locks, but adding a block or</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// growing it requires taking the spine lock.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// Because each mspan covers at least 8K of heap and takes at</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// most 8 bytes in the spanSet, the growth of the spine is</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// quite limited.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// The spine and all blocks are allocated off-heap, which</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// allows this to be used in the memory manager and avoids the</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// need for write barriers on all of these. spanSetBlocks are</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// managed in a pool, though never freed back to the operating</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// system. We never release spine memory because there could be</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// concurrent lock-free access and we&#39;re likely to reuse it</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// anyway. (In principle, we could do this during STW.)</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	spineLock mutex
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	spine     atomicSpanSetSpinePointer <span class="comment">// *[N]atomic.Pointer[spanSetBlock]</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	spineLen  atomic.Uintptr            <span class="comment">// Spine array length</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	spineCap  uintptr                   <span class="comment">// Spine array cap, accessed under spineLock</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// index is the head and tail of the spanSet in a single field.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// The head and the tail both represent an index into the logical</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// concatenation of all blocks, with the head always behind or</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// equal to the tail (indicating an empty set). This field is</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// always accessed atomically.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// The head and the tail are only 32 bits wide, which means we</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// can only support up to 2^32 pushes before a reset. If every</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// span in the heap were stored in this set, and each span were</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// the minimum size (1 runtime page, 8 KiB), then roughly the</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// smallest heap which would be unrepresentable is 32 TiB in size.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	index atomicHeadTailIndex
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>const (
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	spanSetBlockEntries = 512 <span class="comment">// 4KB on 64-bit</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	spanSetInitSpineCap = 256 <span class="comment">// Enough for 1GB heap on 64-bit</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>type spanSetBlock struct {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// Free spanSetBlocks are managed via a lock-free stack.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	lfnode
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// popped is the number of pop operations that have occurred on</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// this block. This number is used to help determine when a block</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// may be safely recycled.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	popped atomic.Uint32
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// spans is the set of spans in this block.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	spans [spanSetBlockEntries]atomicMSpanPointer
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// push adds span s to buffer b. push is safe to call concurrently</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// with other push and pop operations.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>func (b *spanSet) push(s *mspan) {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// Obtain our slot.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	cursor := uintptr(b.index.incTail().tail() - 1)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	top, bottom := cursor/spanSetBlockEntries, cursor%spanSetBlockEntries
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// Do we need to add a block?</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	spineLen := b.spineLen.Load()
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	var block *spanSetBlock
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>retry:
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	if top &lt; spineLen {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		block = b.spine.Load().lookup(top).Load()
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	} else {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		<span class="comment">// Add a new block to the spine, potentially growing</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		<span class="comment">// the spine.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		lock(&amp;b.spineLock)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		<span class="comment">// spineLen cannot change until we release the lock,</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		<span class="comment">// but may have changed while we were waiting.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		spineLen = b.spineLen.Load()
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		if top &lt; spineLen {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			unlock(&amp;b.spineLock)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			goto retry
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		spine := b.spine.Load()
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		if spineLen == b.spineCap {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			<span class="comment">// Grow the spine.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			newCap := b.spineCap * 2
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			if newCap == 0 {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>				newCap = spanSetInitSpineCap
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>			newSpine := persistentalloc(newCap*goarch.PtrSize, cpu.CacheLineSize, &amp;memstats.gcMiscSys)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			if b.spineCap != 0 {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>				<span class="comment">// Blocks are allocated off-heap, so</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>				<span class="comment">// no write barriers.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>				memmove(newSpine, spine.p, b.spineCap*goarch.PtrSize)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			spine = spanSetSpinePointer{newSpine}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			<span class="comment">// Spine is allocated off-heap, so no write barrier.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			b.spine.StoreNoWB(spine)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			b.spineCap = newCap
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			<span class="comment">// We can&#39;t immediately free the old spine</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			<span class="comment">// since a concurrent push with a lower index</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			<span class="comment">// could still be reading from it. We let it</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			<span class="comment">// leak because even a 1TB heap would waste</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			<span class="comment">// less than 2MB of memory on old spines. If</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			<span class="comment">// this is a problem, we could free old spines</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			<span class="comment">// during STW.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		<span class="comment">// Allocate a new block from the pool.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		block = spanSetBlockPool.alloc()
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		<span class="comment">// Add it to the spine.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		<span class="comment">// Blocks are allocated off-heap, so no write barrier.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		spine.lookup(top).StoreNoWB(block)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		b.spineLen.Store(spineLen + 1)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		unlock(&amp;b.spineLock)
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// We have a block. Insert the span atomically, since there may be</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// concurrent readers via the block API.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	block.spans[bottom].StoreNoWB(s)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// pop removes and returns a span from buffer b, or nil if b is empty.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// pop is safe to call concurrently with other pop and push operations.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>func (b *spanSet) pop() *mspan {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	var head, tail uint32
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>claimLoop:
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	for {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		headtail := b.index.load()
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		head, tail = headtail.split()
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		if head &gt;= tail {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			<span class="comment">// The buf is empty, as far as we can tell.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			return nil
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		<span class="comment">// Check if the head position we want to claim is actually</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		<span class="comment">// backed by a block.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		spineLen := b.spineLen.Load()
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		if spineLen &lt;= uintptr(head)/spanSetBlockEntries {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			<span class="comment">// We&#39;re racing with a spine growth and the allocation of</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			<span class="comment">// a new block (and maybe a new spine!), and trying to grab</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			<span class="comment">// the span at the index which is currently being pushed.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			<span class="comment">// Instead of spinning, let&#39;s just notify the caller that</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			<span class="comment">// there&#39;s nothing currently here. Spinning on this is</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			<span class="comment">// almost definitely not worth it.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			return nil
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		<span class="comment">// Try to claim the current head by CASing in an updated head.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		<span class="comment">// This may fail transiently due to a push which modifies the</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		<span class="comment">// tail, so keep trying while the head isn&#39;t changing.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		want := head
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		for want == head {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			if b.index.cas(headtail, makeHeadTailIndex(want+1, tail)) {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>				break claimLoop
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			headtail = b.index.load()
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			head, tail = headtail.split()
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		<span class="comment">// We failed to claim the spot we were after and the head changed,</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		<span class="comment">// meaning a popper got ahead of us. Try again from the top because</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		<span class="comment">// the buf may not be empty.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	top, bottom := head/spanSetBlockEntries, head%spanSetBlockEntries
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// We may be reading a stale spine pointer, but because the length</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// grows monotonically and we&#39;ve already verified it, we&#39;ll definitely</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// be reading from a valid block.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	blockp := b.spine.Load().lookup(uintptr(top))
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// Given that the spine length is correct, we know we will never</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// see a nil block here, since the length is always updated after</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// the block is set.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	block := blockp.Load()
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	s := block.spans[bottom].Load()
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	for s == nil {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		<span class="comment">// We raced with the span actually being set, but given that we</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		<span class="comment">// know a block for this span exists, the race window here is</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		<span class="comment">// extremely small. Try again.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		s = block.spans[bottom].Load()
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// Clear the pointer. This isn&#39;t strictly necessary, but defensively</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// avoids accidentally re-using blocks which could lead to memory</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// corruption. This way, we&#39;ll get a nil pointer access instead.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	block.spans[bottom].StoreNoWB(nil)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// Increase the popped count. If we are the last possible popper</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// in the block (note that bottom need not equal spanSetBlockEntries-1</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	<span class="comment">// due to races) then it&#39;s our responsibility to free the block.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// If we increment popped to spanSetBlockEntries, we can be sure that</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// we&#39;re the last popper for this block, and it&#39;s thus safe to free it.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// Every other popper must have crossed this barrier (and thus finished</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// popping its corresponding mspan) by the time we get here. Because</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">// we&#39;re the last popper, we also don&#39;t have to worry about concurrent</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">// pushers (there can&#39;t be any). Note that we may not be the popper</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	<span class="comment">// which claimed the last slot in the block, we&#39;re just the last one</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// to finish popping.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	if block.popped.Add(1) == spanSetBlockEntries {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		<span class="comment">// Clear the block&#39;s pointer.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		blockp.StoreNoWB(nil)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		<span class="comment">// Return the block to the block pool.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		spanSetBlockPool.free(block)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	return s
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// reset resets a spanSet which is empty. It will also clean up</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// any left over blocks.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// Throws if the buf is not empty.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// reset may not be called concurrently with any other operations</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// on the span set.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>func (b *spanSet) reset() {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	head, tail := b.index.load().split()
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	if head &lt; tail {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		print(&#34;head = &#34;, head, &#34;, tail = &#34;, tail, &#34;\n&#34;)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		throw(&#34;attempt to clear non-empty span set&#34;)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	top := head / spanSetBlockEntries
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	if uintptr(top) &lt; b.spineLen.Load() {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		<span class="comment">// If the head catches up to the tail and the set is empty,</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		<span class="comment">// we may not clean up the block containing the head and tail</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		<span class="comment">// since it may be pushed into again. In order to avoid leaking</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		<span class="comment">// memory since we&#39;re going to reset the head and tail, clean</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		<span class="comment">// up such a block now, if it exists.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		blockp := b.spine.Load().lookup(uintptr(top))
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		block := blockp.Load()
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		if block != nil {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			<span class="comment">// Check the popped value.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			if block.popped.Load() == 0 {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>				<span class="comment">// popped should never be zero because that means we have</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>				<span class="comment">// pushed at least one value but not yet popped if this</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>				<span class="comment">// block pointer is not nil.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>				throw(&#34;span set block with unpopped elements found in reset&#34;)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			if block.popped.Load() == spanSetBlockEntries {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>				<span class="comment">// popped should also never be equal to spanSetBlockEntries</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				<span class="comment">// because the last popper should have made the block pointer</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>				<span class="comment">// in this slot nil.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>				throw(&#34;fully empty unfreed span set block found in reset&#34;)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			<span class="comment">// Clear the pointer to the block.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			blockp.StoreNoWB(nil)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			<span class="comment">// Return the block to the block pool.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			spanSetBlockPool.free(block)
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	b.index.reset()
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	b.spineLen.Store(0)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span><span class="comment">// atomicSpanSetSpinePointer is an atomically-accessed spanSetSpinePointer.</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span><span class="comment">// It has the same semantics as atomic.UnsafePointer.</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>type atomicSpanSetSpinePointer struct {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	a atomic.UnsafePointer
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span><span class="comment">// Loads the spanSetSpinePointer and returns it.</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">// It has the same semantics as atomic.UnsafePointer.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>func (s *atomicSpanSetSpinePointer) Load() spanSetSpinePointer {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	return spanSetSpinePointer{s.a.Load()}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">// Stores the spanSetSpinePointer.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span><span class="comment">// It has the same semantics as [atomic.UnsafePointer].</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>func (s *atomicSpanSetSpinePointer) StoreNoWB(p spanSetSpinePointer) {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	s.a.StoreNoWB(p.p)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">// spanSetSpinePointer represents a pointer to a contiguous block of atomic.Pointer[spanSetBlock].</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>type spanSetSpinePointer struct {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	p unsafe.Pointer
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span><span class="comment">// lookup returns &amp;s[idx].</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>func (s spanSetSpinePointer) lookup(idx uintptr) *atomic.Pointer[spanSetBlock] {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	return (*atomic.Pointer[spanSetBlock])(add(s.p, goarch.PtrSize*idx))
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span><span class="comment">// spanSetBlockPool is a global pool of spanSetBlocks.</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>var spanSetBlockPool spanSetBlockAlloc
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">// spanSetBlockAlloc represents a concurrent pool of spanSetBlocks.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>type spanSetBlockAlloc struct {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	stack lfstack
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// alloc tries to grab a spanSetBlock out of the pool, and if it fails</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// persistentallocs a new one and returns it.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>func (p *spanSetBlockAlloc) alloc() *spanSetBlock {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	if s := (*spanSetBlock)(p.stack.pop()); s != nil {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		return s
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	return (*spanSetBlock)(persistentalloc(unsafe.Sizeof(spanSetBlock{}), cpu.CacheLineSize, &amp;memstats.gcMiscSys))
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">// free returns a spanSetBlock back to the pool.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>func (p *spanSetBlockAlloc) free(block *spanSetBlock) {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	block.popped.Store(0)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	p.stack.push(&amp;block.lfnode)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span><span class="comment">// headTailIndex represents a combined 32-bit head and 32-bit tail</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">// of a queue into a single 64-bit value.</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>type headTailIndex uint64
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// makeHeadTailIndex creates a headTailIndex value from a separate</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">// head and tail.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>func makeHeadTailIndex(head, tail uint32) headTailIndex {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	return headTailIndex(uint64(head)&lt;&lt;32 | uint64(tail))
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// head returns the head of a headTailIndex value.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func (h headTailIndex) head() uint32 {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	return uint32(h &gt;&gt; 32)
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// tail returns the tail of a headTailIndex value.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>func (h headTailIndex) tail() uint32 {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	return uint32(h)
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// split splits the headTailIndex value into its parts.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>func (h headTailIndex) split() (head uint32, tail uint32) {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	return h.head(), h.tail()
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// atomicHeadTailIndex is an atomically-accessed headTailIndex.</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>type atomicHeadTailIndex struct {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	u atomic.Uint64
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// load atomically reads a headTailIndex value.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>func (h *atomicHeadTailIndex) load() headTailIndex {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	return headTailIndex(h.u.Load())
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span><span class="comment">// cas atomically compares-and-swaps a headTailIndex value.</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>func (h *atomicHeadTailIndex) cas(old, new headTailIndex) bool {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	return h.u.CompareAndSwap(uint64(old), uint64(new))
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// incHead atomically increments the head of a headTailIndex.</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>func (h *atomicHeadTailIndex) incHead() headTailIndex {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	return headTailIndex(h.u.Add(1 &lt;&lt; 32))
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">// decHead atomically decrements the head of a headTailIndex.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>func (h *atomicHeadTailIndex) decHead() headTailIndex {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	return headTailIndex(h.u.Add(-(1 &lt;&lt; 32)))
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span><span class="comment">// incTail atomically increments the tail of a headTailIndex.</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>func (h *atomicHeadTailIndex) incTail() headTailIndex {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	ht := headTailIndex(h.u.Add(1))
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	<span class="comment">// Check for overflow.</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	if ht.tail() == 0 {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		print(&#34;runtime: head = &#34;, ht.head(), &#34;, tail = &#34;, ht.tail(), &#34;\n&#34;)
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		throw(&#34;headTailIndex overflow&#34;)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	return ht
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span><span class="comment">// reset clears the headTailIndex to (0, 0).</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>func (h *atomicHeadTailIndex) reset() {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	h.u.Store(0)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">// atomicMSpanPointer is an atomic.Pointer[mspan]. Can&#39;t use generics because it&#39;s NotInHeap.</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>type atomicMSpanPointer struct {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	p atomic.UnsafePointer
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">// Load returns the *mspan.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>func (p *atomicMSpanPointer) Load() *mspan {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	return (*mspan)(p.p.Load())
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span><span class="comment">// Store stores an *mspan.</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>func (p *atomicMSpanPointer) StoreNoWB(s *mspan) {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	p.p.StoreNoWB(unsafe.Pointer(s))
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>
</pre><p><a href="mspanset.go?m=text">View as plain text</a></p>

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
