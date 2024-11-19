<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mbarrier.go - Go Documentation Server</title>

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
<a href="mbarrier.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mbarrier.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2015 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Garbage collector: write barriers.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// For the concurrent garbage collector, the Go compiler implements</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// updates to pointer-valued fields that may be in heap objects by</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// emitting calls to write barriers. The main write barrier for</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// individual pointer writes is gcWriteBarrier and is implemented in</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// assembly. This file contains write barrier entry points for bulk</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// operations. See also mwbbuf.go.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>package runtime
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>import (
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;internal/goexperiment&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// Go uses a hybrid barrier that combines a Yuasa-style deletion</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// barrier—which shades the object whose reference is being</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// overwritten—with Dijkstra insertion barrier—which shades the object</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// whose reference is being written. The insertion part of the barrier</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// is necessary while the calling goroutine&#39;s stack is grey. In</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// pseudocode, the barrier is:</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//     writePointer(slot, ptr):</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//         shade(*slot)</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//         if current stack is grey:</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//             shade(ptr)</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//         *slot = ptr</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// slot is the destination in Go code.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// ptr is the value that goes into the slot in Go code.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// Shade indicates that it has seen a white pointer by adding the referent</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// to wbuf as well as marking it.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// The two shades and the condition work together to prevent a mutator</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// from hiding an object from the garbage collector:</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// 1. shade(*slot) prevents a mutator from hiding an object by moving</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// the sole pointer to it from the heap to its stack. If it attempts</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// to unlink an object from the heap, this will shade it.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// 2. shade(ptr) prevents a mutator from hiding an object by moving</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// the sole pointer to it from its stack into a black object in the</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// heap. If it attempts to install the pointer into a black object,</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// this will shade it.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// 3. Once a goroutine&#39;s stack is black, the shade(ptr) becomes</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// unnecessary. shade(ptr) prevents hiding an object by moving it from</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// the stack to the heap, but this requires first having a pointer</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// hidden on the stack. Immediately after a stack is scanned, it only</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// points to shaded objects, so it&#39;s not hiding anything, and the</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// shade(*slot) prevents it from hiding any other pointers on its</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// stack.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// For a detailed description of this barrier and proof of</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// correctness, see https://github.com/golang/proposal/blob/master/design/17503-eliminate-rescan.md</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// Dealing with memory ordering:</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// Both the Yuasa and Dijkstra barriers can be made conditional on the</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// color of the object containing the slot. We chose not to make these</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// conditional because the cost of ensuring that the object holding</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// the slot doesn&#39;t concurrently change color without the mutator</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// noticing seems prohibitive.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// Consider the following example where the mutator writes into</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// a slot and then loads the slot&#39;s mark bit while the GC thread</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// writes to the slot&#39;s mark bit and then as part of scanning reads</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// the slot.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// Initially both [slot] and [slotmark] are 0 (nil)</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// Mutator thread          GC thread</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// st [slot], ptr          st [slotmark], 1</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// ld r1, [slotmark]       ld r2, [slot]</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// Without an expensive memory barrier between the st and the ld, the final</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// result on most HW (including 386/amd64) can be r1==r2==0. This is a classic</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// example of what can happen when loads are allowed to be reordered with older</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// stores (avoiding such reorderings lies at the heart of the classic</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// Peterson/Dekker algorithms for mutual exclusion). Rather than require memory</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// barriers, which will slow down both the mutator and the GC, we always grey</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// the ptr object regardless of the slot&#39;s color.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// Another place where we intentionally omit memory barriers is when</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// accessing mheap_.arena_used to check if a pointer points into the</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// heap. On relaxed memory machines, it&#39;s possible for a mutator to</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// extend the size of the heap by updating arena_used, allocate an</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// object from this new region, and publish a pointer to that object,</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// but for tracing running on another processor to observe the pointer</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// but use the old value of arena_used. In this case, tracing will not</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// mark the object, even though it&#39;s reachable. However, the mutator</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// is guaranteed to execute a write barrier when it publishes the</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// pointer, so it will take care of marking the object. A general</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// consequence of this is that the garbage collector may cache the</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// value of mheap_.arena_used. (See issue #9984.)</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// Stack writes:</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// The compiler omits write barriers for writes to the current frame,</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// but if a stack pointer has been passed down the call stack, the</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// compiler will generate a write barrier for writes through that</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// pointer (because it doesn&#39;t know it&#39;s not a heap pointer).</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// Global writes:</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// The Go garbage collector requires write barriers when heap pointers</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// are stored in globals. Many garbage collectors ignore writes to</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// globals and instead pick up global -&gt; heap pointers during</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// termination. This increases pause time, so we instead rely on write</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// barriers for writes to globals so that we don&#39;t have to rescan</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// global during mark termination.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// Publication ordering:</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// The write barrier is *pre-publication*, meaning that the write</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// barrier happens prior to the *slot = ptr write that may make ptr</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// reachable by some goroutine that currently cannot reach it.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">// Signal handler pointer writes:</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">// In general, the signal handler cannot safely invoke the write</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// barrier because it may run without a P or even during the write</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// barrier.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// There is exactly one exception: profbuf.go omits a barrier during</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// signal handler profile logging. That&#39;s safe only because of the</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// deletion barrier. See profbuf.go for a detailed argument. If we</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// remove the deletion barrier, we&#39;ll have to work out a new way to</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// handle the profile logging.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// typedmemmove copies a value of type typ to dst from src.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// Must be nosplit, see #16026.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// TODO: Perfect for go:nosplitrec since we can&#39;t have a safe point</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// anywhere in the bulk barrier or memmove.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func typedmemmove(typ *abi.Type, dst, src unsafe.Pointer) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if dst == src {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if writeBarrier.enabled &amp;&amp; typ.PtrBytes != 0 {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// This always copies a full value of type typ so it&#39;s safe</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// to pass typ along as an optimization. See the comment on</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		<span class="comment">// bulkBarrierPreWrite.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		bulkBarrierPreWrite(uintptr(dst), uintptr(src), typ.PtrBytes, typ)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// There&#39;s a race here: if some other goroutine can write to</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// src, it may change some pointer in src after we&#39;ve</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// performed the write barrier but before we perform the</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// memory copy. This safe because the write performed by that</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// other goroutine must also be accompanied by a write</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// barrier, so at worst we&#39;ve unnecessarily greyed the old</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// pointer that was in src.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	memmove(dst, src, typ.Size_)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if goexperiment.CgoCheck2 {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		cgoCheckMemmove2(typ, dst, src, 0, typ.Size_)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// wbZero performs the write barrier operations necessary before</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// zeroing a region of memory at address dst of type typ.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">// Does not actually do the zeroing.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>func wbZero(typ *_type, dst unsafe.Pointer) {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// This always copies a full value of type typ so it&#39;s safe</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// to pass typ along as an optimization. See the comment on</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// bulkBarrierPreWrite.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	bulkBarrierPreWrite(uintptr(dst), 0, typ.PtrBytes, typ)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// wbMove performs the write barrier operations necessary before</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// copying a region of memory from src to dst of type typ.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// Does not actually do the copying.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>func wbMove(typ *_type, dst, src unsafe.Pointer) {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// This always copies a full value of type typ so it&#39;s safe to</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// pass a type here.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// See the comment on bulkBarrierPreWrite.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	bulkBarrierPreWrite(uintptr(dst), uintptr(src), typ.PtrBytes, typ)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_typedmemmove reflect.typedmemmove</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>func reflect_typedmemmove(typ *_type, dst, src unsafe.Pointer) {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	if raceenabled {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		raceWriteObjectPC(typ, dst, getcallerpc(), abi.FuncPCABIInternal(reflect_typedmemmove))
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		raceReadObjectPC(typ, src, getcallerpc(), abi.FuncPCABIInternal(reflect_typedmemmove))
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	if msanenabled {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		msanwrite(dst, typ.Size_)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		msanread(src, typ.Size_)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	if asanenabled {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		asanwrite(dst, typ.Size_)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		asanread(src, typ.Size_)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	typedmemmove(typ, dst, src)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">//go:linkname reflectlite_typedmemmove internal/reflectlite.typedmemmove</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>func reflectlite_typedmemmove(typ *_type, dst, src unsafe.Pointer) {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	reflect_typedmemmove(typ, dst, src)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// reflectcallmove is invoked by reflectcall to copy the return values</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// out of the stack and into the heap, invoking the necessary write</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// barriers. dst, src, and size describe the return value area to</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// copy. typ describes the entire frame (not just the return values).</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// typ may be nil, which indicates write barriers are not needed.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// It must be nosplit and must only call nosplit functions because the</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// stack map of reflectcall is wrong.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>func reflectcallmove(typ *_type, dst, src unsafe.Pointer, size uintptr, regs *abi.RegArgs) {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	if writeBarrier.enabled &amp;&amp; typ != nil &amp;&amp; typ.PtrBytes != 0 &amp;&amp; size &gt;= goarch.PtrSize {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		<span class="comment">// Pass nil for the type. dst does not point to value of type typ,</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		<span class="comment">// but rather points into one, so applying the optimization is not</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		<span class="comment">// safe. See the comment on this function.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		bulkBarrierPreWrite(uintptr(dst), uintptr(src), size, nil)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	memmove(dst, src, size)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// Move pointers returned in registers to a place where the GC can see them.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	for i := range regs.Ints {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		if regs.ReturnIsPtr.Get(i) {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			regs.Ptrs[i] = unsafe.Pointer(regs.Ints[i])
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>func typedslicecopy(typ *_type, dstPtr unsafe.Pointer, dstLen int, srcPtr unsafe.Pointer, srcLen int) int {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	n := dstLen
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	if n &gt; srcLen {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		n = srcLen
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	if n == 0 {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		return 0
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// The compiler emits calls to typedslicecopy before</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// instrumentation runs, so unlike the other copying and</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// assignment operations, it&#39;s not instrumented in the calling</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// code and needs its own instrumentation.</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	if raceenabled {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		pc := abi.FuncPCABIInternal(slicecopy)
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		racewriterangepc(dstPtr, uintptr(n)*typ.Size_, callerpc, pc)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		racereadrangepc(srcPtr, uintptr(n)*typ.Size_, callerpc, pc)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	if msanenabled {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		msanwrite(dstPtr, uintptr(n)*typ.Size_)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		msanread(srcPtr, uintptr(n)*typ.Size_)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	if asanenabled {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		asanwrite(dstPtr, uintptr(n)*typ.Size_)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		asanread(srcPtr, uintptr(n)*typ.Size_)
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	if goexperiment.CgoCheck2 {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		cgoCheckSliceCopy(typ, dstPtr, srcPtr, n)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	if dstPtr == srcPtr {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		return n
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	<span class="comment">// Note: No point in checking typ.PtrBytes here:</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	<span class="comment">// compiler only emits calls to typedslicecopy for types with pointers,</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">// and growslice and reflect_typedslicecopy check for pointers</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// before calling typedslicecopy.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	size := uintptr(n) * typ.Size_
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	if writeBarrier.enabled {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		<span class="comment">// This always copies one or more full values of type typ so</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		<span class="comment">// it&#39;s safe to pass typ along as an optimization. See the comment on</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		<span class="comment">// bulkBarrierPreWrite.</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		pwsize := size - typ.Size_ + typ.PtrBytes
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		bulkBarrierPreWrite(uintptr(dstPtr), uintptr(srcPtr), pwsize, typ)
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">// See typedmemmove for a discussion of the race between the</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">// barrier and memmove.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	memmove(dstPtr, srcPtr, size)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	return n
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_typedslicecopy reflect.typedslicecopy</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>func reflect_typedslicecopy(elemType *_type, dst, src slice) int {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	if elemType.PtrBytes == 0 {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		return slicecopy(dst.array, dst.len, src.array, src.len, elemType.Size_)
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	}
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	return typedslicecopy(elemType, dst.array, dst.len, src.array, src.len)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">// typedmemclr clears the typed memory at ptr with type typ. The</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// memory at ptr must already be initialized (and hence in type-safe</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">// state). If the memory is being initialized for the first time, see</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">// memclrNoHeapPointers.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">// If the caller knows that typ has pointers, it can alternatively</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span><span class="comment">// call memclrHasPointers.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span><span class="comment">// TODO: A &#34;go:nosplitrec&#34; annotation would be perfect for this.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>func typedmemclr(typ *_type, ptr unsafe.Pointer) {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	if writeBarrier.enabled &amp;&amp; typ.PtrBytes != 0 {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		<span class="comment">// This always clears a whole value of type typ, so it&#39;s</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		<span class="comment">// safe to pass a type here and apply the optimization.</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		<span class="comment">// See the comment on bulkBarrierPreWrite.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		bulkBarrierPreWrite(uintptr(ptr), 0, typ.PtrBytes, typ)
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	memclrNoHeapPointers(ptr, typ.Size_)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_typedmemclr reflect.typedmemclr</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func reflect_typedmemclr(typ *_type, ptr unsafe.Pointer) {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	typedmemclr(typ, ptr)
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_typedmemclrpartial reflect.typedmemclrpartial</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>func reflect_typedmemclrpartial(typ *_type, ptr unsafe.Pointer, off, size uintptr) {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	if writeBarrier.enabled &amp;&amp; typ.PtrBytes != 0 {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		<span class="comment">// Pass nil for the type. ptr does not point to value of type typ,</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		<span class="comment">// but rather points into one so it&#39;s not safe to apply the optimization.</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		<span class="comment">// See the comment on this function in the reflect package and the</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		<span class="comment">// comment on bulkBarrierPreWrite.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		bulkBarrierPreWrite(uintptr(ptr), 0, size, nil)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	memclrNoHeapPointers(ptr, size)
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_typedarrayclear reflect.typedarrayclear</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>func reflect_typedarrayclear(typ *_type, ptr unsafe.Pointer, len int) {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	size := typ.Size_ * uintptr(len)
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	if writeBarrier.enabled &amp;&amp; typ.PtrBytes != 0 {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		<span class="comment">// This always clears whole elements of an array, so it&#39;s</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		<span class="comment">// safe to pass a type here. See the comment on bulkBarrierPreWrite.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		bulkBarrierPreWrite(uintptr(ptr), 0, size, typ)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	memclrNoHeapPointers(ptr, size)
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span><span class="comment">// memclrHasPointers clears n bytes of typed memory starting at ptr.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">// The caller must ensure that the type of the object at ptr has</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// pointers, usually by checking typ.PtrBytes. However, ptr</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// does not have to point to the start of the allocation.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>func memclrHasPointers(ptr unsafe.Pointer, n uintptr) {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	<span class="comment">// Pass nil for the type since we don&#39;t have one here anyway.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	bulkBarrierPreWrite(uintptr(ptr), 0, n, nil)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	memclrNoHeapPointers(ptr, n)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
</pre><p><a href="mbarrier.go?m=text">View as plain text</a></p>

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
