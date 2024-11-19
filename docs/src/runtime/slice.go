<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/slice.go - Go Documentation Server</title>

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
<a href="slice.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">slice.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;runtime/internal/math&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>type slice struct {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	array unsafe.Pointer
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	len   int
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	cap   int
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>}
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// A notInHeapSlice is a slice backed by runtime/internal/sys.NotInHeap memory.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>type notInHeapSlice struct {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	array *notInHeap
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	len   int
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	cap   int
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>}
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>func panicmakeslicelen() {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	panic(errorString(&#34;makeslice: len out of range&#34;))
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>func panicmakeslicecap() {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	panic(errorString(&#34;makeslice: cap out of range&#34;))
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// makeslicecopy allocates a slice of &#34;tolen&#34; elements of type &#34;et&#34;,</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// then copies &#34;fromlen&#34; elements of type &#34;et&#34; into that new allocation from &#34;from&#34;.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>func makeslicecopy(et *_type, tolen int, fromlen int, from unsafe.Pointer) unsafe.Pointer {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	var tomem, copymem uintptr
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	if uintptr(tolen) &gt; uintptr(fromlen) {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		var overflow bool
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		tomem, overflow = math.MulUintptr(et.Size_, uintptr(tolen))
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		if overflow || tomem &gt; maxAlloc || tolen &lt; 0 {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>			panicmakeslicelen()
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		copymem = et.Size_ * uintptr(fromlen)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	} else {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		<span class="comment">// fromlen is a known good length providing and equal or greater than tolen,</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		<span class="comment">// thereby making tolen a good slice length too as from and to slices have the</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		<span class="comment">// same element width.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		tomem = et.Size_ * uintptr(tolen)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		copymem = tomem
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	var to unsafe.Pointer
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	if et.PtrBytes == 0 {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		to = mallocgc(tomem, nil, false)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		if copymem &lt; tomem {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			memclrNoHeapPointers(add(to, copymem), tomem-copymem)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	} else {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		<span class="comment">// Note: can&#39;t use rawmem (which avoids zeroing of memory), because then GC can scan uninitialized memory.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		to = mallocgc(tomem, et, true)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		if copymem &gt; 0 &amp;&amp; writeBarrier.enabled {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			<span class="comment">// Only shade the pointers in old.array since we know the destination slice to</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			<span class="comment">// only contains nil pointers because it has been cleared during alloc.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			<span class="comment">// It&#39;s safe to pass a type to this function as an optimization because</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			<span class="comment">// from and to only ever refer to memory representing whole values of</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			<span class="comment">// type et. See the comment on bulkBarrierPreWrite.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>			bulkBarrierPreWriteSrcOnly(uintptr(to), uintptr(from), copymem, et)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	if raceenabled {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		pc := abi.FuncPCABIInternal(makeslicecopy)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		racereadrangepc(from, copymem, callerpc, pc)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if msanenabled {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		msanread(from, copymem)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	if asanenabled {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		asanread(from, copymem)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	memmove(to, from, copymem)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	return to
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>func makeslice(et *_type, len, cap int) unsafe.Pointer {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	mem, overflow := math.MulUintptr(et.Size_, uintptr(cap))
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	if overflow || mem &gt; maxAlloc || len &lt; 0 || len &gt; cap {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		<span class="comment">// NOTE: Produce a &#39;len out of range&#39; error instead of a</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		<span class="comment">// &#39;cap out of range&#39; error when someone does make([]T, bignumber).</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		<span class="comment">// &#39;cap out of range&#39; is true too, but since the cap is only being</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		<span class="comment">// supplied implicitly, saying len is clearer.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		<span class="comment">// See golang.org/issue/4085.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		mem, overflow := math.MulUintptr(et.Size_, uintptr(len))
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		if overflow || mem &gt; maxAlloc || len &lt; 0 {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			panicmakeslicelen()
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		panicmakeslicecap()
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	return mallocgc(mem, et, true)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>func makeslice64(et *_type, len64, cap64 int64) unsafe.Pointer {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	len := int(len64)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	if int64(len) != len64 {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		panicmakeslicelen()
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	cap := int(cap64)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	if int64(cap) != cap64 {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		panicmakeslicecap()
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	return makeslice(et, len, cap)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// growslice allocates new backing store for a slice.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// arguments:</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">//	oldPtr = pointer to the slice&#39;s backing array</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">//	newLen = new length (= oldLen + num)</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">//	oldCap = original slice&#39;s capacity.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">//	   num = number of elements being added</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">//	    et = element type</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">// return values:</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">//	newPtr = pointer to the new backing store</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">//	newLen = same value as the argument</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">//	newCap = capacity of the new backing store</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// Requires that uint(newLen) &gt; uint(oldCap).</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// Assumes the original slice length is newLen - num</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// A new backing store is allocated with space for at least newLen elements.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// Existing entries [0, oldLen) are copied over to the new backing store.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// Added entries [oldLen, newLen) are not initialized by growslice</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// (although for pointer-containing element types, they are zeroed). They</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// must be initialized by the caller.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// Trailing entries [newLen, newCap) are zeroed.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// growslice&#39;s odd calling convention makes the generated code that calls</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// this function simpler. In particular, it accepts and returns the</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">// new length so that the old length is not live (does not need to be</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">// spilled/restored) and the new length is returned (also does not need</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// to be spilled/restored).</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	oldLen := newLen - num
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	if raceenabled {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		racereadrangepc(oldPtr, uintptr(oldLen*int(et.Size_)), callerpc, abi.FuncPCABIInternal(growslice))
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if msanenabled {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		msanread(oldPtr, uintptr(oldLen*int(et.Size_)))
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	if asanenabled {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		asanread(oldPtr, uintptr(oldLen*int(et.Size_)))
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	if newLen &lt; 0 {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		panic(errorString(&#34;growslice: len out of range&#34;))
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	if et.Size_ == 0 {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		<span class="comment">// append should not create a slice with nil pointer but non-zero len.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		<span class="comment">// We assume that append doesn&#39;t need to preserve oldPtr in this case.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		return slice{unsafe.Pointer(&amp;zerobase), newLen, newLen}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	newcap := nextslicecap(newLen, oldCap)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	var overflow bool
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	var lenmem, newlenmem, capmem uintptr
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// Specialize for common values of et.Size.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// For 1 we don&#39;t need any division/multiplication.</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// For goarch.PtrSize, compiler will optimize division/multiplication into a shift by a constant.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// For powers of 2, use a variable shift.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	noscan := et.PtrBytes == 0
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	switch {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	case et.Size_ == 1:
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		lenmem = uintptr(oldLen)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		newlenmem = uintptr(newLen)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		capmem = roundupsize(uintptr(newcap), noscan)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		overflow = uintptr(newcap) &gt; maxAlloc
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		newcap = int(capmem)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	case et.Size_ == goarch.PtrSize:
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		lenmem = uintptr(oldLen) * goarch.PtrSize
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		newlenmem = uintptr(newLen) * goarch.PtrSize
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		capmem = roundupsize(uintptr(newcap)*goarch.PtrSize, noscan)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		overflow = uintptr(newcap) &gt; maxAlloc/goarch.PtrSize
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		newcap = int(capmem / goarch.PtrSize)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	case isPowerOfTwo(et.Size_):
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		var shift uintptr
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		if goarch.PtrSize == 8 {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			<span class="comment">// Mask shift for better code generation.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			shift = uintptr(sys.TrailingZeros64(uint64(et.Size_))) &amp; 63
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		} else {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			shift = uintptr(sys.TrailingZeros32(uint32(et.Size_))) &amp; 31
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		lenmem = uintptr(oldLen) &lt;&lt; shift
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		newlenmem = uintptr(newLen) &lt;&lt; shift
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		capmem = roundupsize(uintptr(newcap)&lt;&lt;shift, noscan)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		overflow = uintptr(newcap) &gt; (maxAlloc &gt;&gt; shift)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		newcap = int(capmem &gt;&gt; shift)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		capmem = uintptr(newcap) &lt;&lt; shift
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	default:
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		lenmem = uintptr(oldLen) * et.Size_
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		newlenmem = uintptr(newLen) * et.Size_
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		capmem, overflow = math.MulUintptr(et.Size_, uintptr(newcap))
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		capmem = roundupsize(capmem, noscan)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		newcap = int(capmem / et.Size_)
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		capmem = uintptr(newcap) * et.Size_
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	<span class="comment">// The check of overflow in addition to capmem &gt; maxAlloc is needed</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">// to prevent an overflow which can be used to trigger a segfault</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	<span class="comment">// on 32bit architectures with this example program:</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	<span class="comment">// type T [1&lt;&lt;27 + 1]int64</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// var d T</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// var s []T</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">// func main() {</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">//   s = append(s, d, d, d, d)</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	<span class="comment">//   print(len(s), &#34;\n&#34;)</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">// }</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	if overflow || capmem &gt; maxAlloc {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		panic(errorString(&#34;growslice: len out of range&#34;))
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	var p unsafe.Pointer
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	if et.PtrBytes == 0 {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		p = mallocgc(capmem, nil, false)
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		<span class="comment">// The append() that calls growslice is going to overwrite from oldLen to newLen.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		<span class="comment">// Only clear the part that will not be overwritten.</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		<span class="comment">// The reflect_growslice() that calls growslice will manually clear</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		<span class="comment">// the region not cleared here.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		memclrNoHeapPointers(add(p, newlenmem), capmem-newlenmem)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	} else {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		<span class="comment">// Note: can&#39;t use rawmem (which avoids zeroing of memory), because then GC can scan uninitialized memory.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		p = mallocgc(capmem, et, true)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		if lenmem &gt; 0 &amp;&amp; writeBarrier.enabled {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			<span class="comment">// Only shade the pointers in oldPtr since we know the destination slice p</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			<span class="comment">// only contains nil pointers because it has been cleared during alloc.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			<span class="comment">// It&#39;s safe to pass a type to this function as an optimization because</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			<span class="comment">// from and to only ever refer to memory representing whole values of</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			<span class="comment">// type et. See the comment on bulkBarrierPreWrite.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			bulkBarrierPreWriteSrcOnly(uintptr(p), uintptr(oldPtr), lenmem-et.Size_+et.PtrBytes, et)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	memmove(p, oldPtr, lenmem)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	return slice{p, newLen, newcap}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">// nextslicecap computes the next appropriate slice length.</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>func nextslicecap(newLen, oldCap int) int {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	newcap := oldCap
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	doublecap := newcap + newcap
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	if newLen &gt; doublecap {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		return newLen
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	const threshold = 256
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	if oldCap &lt; threshold {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		return doublecap
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	for {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		<span class="comment">// Transition from growing 2x for small slices</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		<span class="comment">// to growing 1.25x for large slices. This formula</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		<span class="comment">// gives a smooth-ish transition between the two.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		newcap += (newcap + 3*threshold) &gt;&gt; 2
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		<span class="comment">// We need to check `newcap &gt;= newLen` and whether `newcap` overflowed.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		<span class="comment">// newLen is guaranteed to be larger than zero, hence</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		<span class="comment">// when newcap overflows then `uint(newcap) &gt; uint(newLen)`.</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		<span class="comment">// This allows to check for both with the same comparison.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		if uint(newcap) &gt;= uint(newLen) {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			break
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// Set newcap to the requested cap when</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">// the newcap calculation overflowed.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	if newcap &lt;= 0 {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		return newLen
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	return newcap
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_growslice reflect.growslice</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>func reflect_growslice(et *_type, old slice, num int) slice {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	<span class="comment">// Semantically equivalent to slices.Grow, except that the caller</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	<span class="comment">// is responsible for ensuring that old.len+num &gt; old.cap.</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	num -= old.cap - old.len <span class="comment">// preserve memory of old[old.len:old.cap]</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	new := growslice(old.array, old.cap+num, old.cap, num, et)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	<span class="comment">// growslice does not zero out new[old.cap:new.len] since it assumes that</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	<span class="comment">// the memory will be overwritten by an append() that called growslice.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	<span class="comment">// Since the caller of reflect_growslice is not append(),</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	<span class="comment">// zero out this region before returning the slice to the reflect package.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	if et.PtrBytes == 0 {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		oldcapmem := uintptr(old.cap) * et.Size_
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		newlenmem := uintptr(new.len) * et.Size_
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		memclrNoHeapPointers(add(new.array, oldcapmem), newlenmem-oldcapmem)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	new.len = old.len <span class="comment">// preserve the old length</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	return new
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>func isPowerOfTwo(x uintptr) bool {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	return x&amp;(x-1) == 0
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span><span class="comment">// slicecopy is used to copy from a string or slice of pointerless elements into a slice.</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>func slicecopy(toPtr unsafe.Pointer, toLen int, fromPtr unsafe.Pointer, fromLen int, width uintptr) int {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	if fromLen == 0 || toLen == 0 {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		return 0
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	n := fromLen
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	if toLen &lt; n {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		n = toLen
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	if width == 0 {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		return n
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	size := uintptr(n) * width
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	if raceenabled {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		pc := abi.FuncPCABIInternal(slicecopy)
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		racereadrangepc(fromPtr, size, callerpc, pc)
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		racewriterangepc(toPtr, size, callerpc, pc)
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	if msanenabled {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		msanread(fromPtr, size)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		msanwrite(toPtr, size)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	if asanenabled {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		asanread(fromPtr, size)
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		asanwrite(toPtr, size)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	if size == 1 { <span class="comment">// common case worth about 2x to do here</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		<span class="comment">// TODO: is this still worth it with new memmove impl?</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		*(*byte)(toPtr) = *(*byte)(fromPtr) <span class="comment">// known to be a byte pointer</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	} else {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		memmove(toPtr, fromPtr, size)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	return n
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">//go:linkname bytealg_MakeNoZero internal/bytealg.MakeNoZero</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>func bytealg_MakeNoZero(len int) []byte {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	if uintptr(len) &gt; maxAlloc {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		panicmakeslicelen()
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	return unsafe.Slice((*byte)(mallocgc(uintptr(len), nil, false)), len)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
</pre><p><a href="slice.go?m=text">View as plain text</a></p>

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
