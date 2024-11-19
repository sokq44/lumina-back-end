<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mbitmap_allocheaders.go - Go Documentation Server</title>

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
<a href="mbitmap_allocheaders.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mbitmap_allocheaders.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build goexperiment.allocheaders</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Garbage collector: type and heap bitmaps.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// Stack, data, and bss bitmaps</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// Stack frames and global variables in the data and bss sections are</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// described by bitmaps with 1 bit per pointer-sized word. A &#34;1&#34; bit</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// means the word is a live pointer to be visited by the GC (referred to</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// as &#34;pointer&#34;). A &#34;0&#34; bit means the word should be ignored by GC</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// (referred to as &#34;scalar&#34;, though it could be a dead pointer value).</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Heap bitmaps</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// The heap bitmap comprises 1 bit for each pointer-sized word in the heap,</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// recording whether a pointer is stored in that word or not. This bitmap</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// is stored at the end of a span for small objects and is unrolled at</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// runtime from type metadata for all larger objects. Objects without</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// pointers have neither a bitmap nor associated type metadata.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// Bits in all cases correspond to words in little-endian order.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// For small objects, if s is the mspan for the span starting at &#34;start&#34;,</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// then s.heapBits() returns a slice containing the bitmap for the whole span.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// That is, s.heapBits()[0] holds the goarch.PtrSize*8 bits for the first</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// goarch.PtrSize*8 words from &#34;start&#34; through &#34;start+63*ptrSize&#34; in the span.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// On a related note, small objects are always small enough that their bitmap</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// fits in goarch.PtrSize*8 bits, so writing out bitmap data takes two bitmap</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// writes at most (because object boundaries don&#39;t generally lie on</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// s.heapBits()[i] boundaries).</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// For larger objects, if t is the type for the object starting at &#34;start&#34;,</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// within some span whose mspan is s, then the bitmap at t.GCData is &#34;tiled&#34;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// from &#34;start&#34; through &#34;start+s.elemsize&#34;.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// Specifically, the first bit of t.GCData corresponds to the word at &#34;start&#34;,</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// the second to the word after &#34;start&#34;, and so on up to t.PtrBytes. At t.PtrBytes,</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// we skip to &#34;start+t.Size_&#34; and begin again from there. This process is</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// repeated until we hit &#34;start+s.elemsize&#34;.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// This tiling algorithm supports array data, since the type always refers to</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// the element type of the array. Single objects are considered the same as</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// single-element arrays.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// The tiling algorithm may scan data past the end of the compiler-recognized</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// object, but any unused data within the allocation slot (i.e. within s.elemsize)</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// is zeroed, so the GC just observes nil pointers.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// Note that this &#34;tiled&#34; bitmap isn&#39;t stored anywhere; it is generated on-the-fly.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// For objects without their own span, the type metadata is stored in the first</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// word before the object at the beginning of the allocation slot. For objects</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// with their own span, the type metadata is stored in the mspan.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// The bitmap for small unallocated objects in scannable spans is not maintained</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// (can be junk).</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>package runtime
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>import (
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>const (
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// A malloc header is functionally a single type pointer, but</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// we need to use 8 here to ensure 8-byte alignment of allocations</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// on 32-bit platforms. It&#39;s wasteful, but a lot of code relies on</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// 8-byte alignment for 8-byte atomics.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	mallocHeaderSize = 8
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// The minimum object size that has a malloc header, exclusive.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// The size of this value controls overheads from the malloc header.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// The minimum size is bound by writeHeapBitsSmall, which assumes that the</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// pointer bitmap for objects of a size smaller than this doesn&#39;t cross</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// more than one pointer-word boundary. This sets an upper-bound on this</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// value at the number of bits in a uintptr, multiplied by the pointer</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// size in bytes.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// We choose a value here that has a natural cutover point in terms of memory</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// overheads. This value just happens to be the maximum possible value this</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// can be.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// A span with heap bits in it will have 128 bytes of heap bits on 64-bit</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// platforms, and 256 bytes of heap bits on 32-bit platforms. The first size</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// class where malloc headers match this overhead for 64-bit platforms is</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// 512 bytes (8 KiB / 512 bytes * 8 bytes-per-header = 128 bytes of overhead).</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// On 32-bit platforms, this same point is the 256 byte size class</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// (8 KiB / 256 bytes * 8 bytes-per-header = 256 bytes of overhead).</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// Guaranteed to be exactly at a size class boundary. The reason this value is</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// an exclusive minimum is subtle. Suppose we&#39;re allocating a 504-byte object</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// and its rounded up to 512 bytes for the size class. If minSizeForMallocHeader</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// is 512 and an inclusive minimum, then a comparison against minSizeForMallocHeader</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// by the two values would produce different results. In other words, the comparison</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// would not be invariant to size-class rounding. Eschewing this property means a</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// more complex check or possibly storing additional state to determine whether a</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// span has malloc headers.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	minSizeForMallocHeader = goarch.PtrSize * ptrBits
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// heapBitsInSpan returns true if the size of an object implies its ptr/scalar</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// data is stored at the end of the span, and is accessible via span.heapBits.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// Note: this works for both rounded-up sizes (span.elemsize) and unrounded</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// type sizes because minSizeForMallocHeader is guaranteed to be at a size</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// class boundary.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func heapBitsInSpan(userSize uintptr) bool {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// N.B. minSizeForMallocHeader is an exclusive minimum so that this function is</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// invariant under size-class rounding on its input.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	return userSize &lt;= minSizeForMallocHeader
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// heapArenaPtrScalar contains the per-heapArena pointer/scalar metadata for the GC.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>type heapArenaPtrScalar struct {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// N.B. This is no longer necessary with allocation headers.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// typePointers is an iterator over the pointers in a heap object.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// Iteration through this type implements the tiling algorithm described at the</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// top of this file.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>type typePointers struct {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// elem is the address of the current array element of type typ being iterated over.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// Objects that are not arrays are treated as single-element arrays, in which case</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// this value does not change.</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	elem uintptr
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// addr is the address the iterator is currently working from and describes</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// the address of the first word referenced by mask.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	addr uintptr
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// mask is a bitmask where each bit corresponds to pointer-words after addr.</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// Bit 0 is the pointer-word at addr, Bit 1 is the next word, and so on.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// If a bit is 1, then there is a pointer at that word.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// nextFast and next mask out bits in this mask as their pointers are processed.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	mask uintptr
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// typ is a pointer to the type information for the heap object&#39;s type.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// This may be nil if the object is in a span where heapBitsInSpan(span.elemsize) is true.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	typ *_type
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// typePointersOf returns an iterator over all heap pointers in the range [addr, addr+size).</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// addr and addr+size must be in the range [span.base(), span.limit).</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">// Note: addr+size must be passed as the limit argument to the iterator&#39;s next method on</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// each iteration. This slightly awkward API is to allow typePointers to be destructured</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// by the compiler.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// nosplit because it is used during write barriers and must not be preempted.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func (span *mspan) typePointersOf(addr, size uintptr) typePointers {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	base := span.objBase(addr)
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	tp := span.typePointersOfUnchecked(base)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if base == addr &amp;&amp; size == span.elemsize {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		return tp
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	return tp.fastForward(addr-tp.addr, addr+size)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// typePointersOfUnchecked is like typePointersOf, but assumes addr is the base</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">// of an allocation slot in a span (the start of the object if no header, the</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// header otherwise). It returns an iterator that generates all pointers</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// in the range [addr, addr+span.elemsize).</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">// nosplit because it is used during write barriers and must not be preempted.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>func (span *mspan) typePointersOfUnchecked(addr uintptr) typePointers {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	const doubleCheck = false
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	if doubleCheck &amp;&amp; span.objBase(addr) != addr {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		print(&#34;runtime: addr=&#34;, addr, &#34; base=&#34;, span.objBase(addr), &#34;\n&#34;)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		throw(&#34;typePointersOfUnchecked consisting of non-base-address for object&#34;)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	spc := span.spanclass
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if spc.noscan() {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		return typePointers{}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	if heapBitsInSpan(span.elemsize) {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		<span class="comment">// Handle header-less objects.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		return typePointers{elem: addr, addr: addr, mask: span.heapBitsSmallForAddr(addr)}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// All of these objects have a header.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	var typ *_type
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	if spc.sizeclass() != 0 {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		<span class="comment">// Pull the allocation header from the first word of the object.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		typ = *(**_type)(unsafe.Pointer(addr))
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		addr += mallocHeaderSize
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	} else {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		typ = span.largeType
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	gcdata := typ.GCData
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	return typePointers{elem: addr, addr: addr, mask: readUintptr(gcdata), typ: typ}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// typePointersOfType is like typePointersOf, but assumes addr points to one or more</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// contiguous instances of the provided type. The provided type must not be nil and</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// it must not have its type metadata encoded as a gcprog.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// It returns an iterator that tiles typ.GCData starting from addr. It&#39;s the caller&#39;s</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// responsibility to limit iteration.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// nosplit because its callers are nosplit and require all their callees to be nosplit.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>func (span *mspan) typePointersOfType(typ *abi.Type, addr uintptr) typePointers {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	const doubleCheck = false
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	if doubleCheck &amp;&amp; (typ == nil || typ.Kind_&amp;kindGCProg != 0) {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		throw(&#34;bad type passed to typePointersOfType&#34;)
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	if span.spanclass.noscan() {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		return typePointers{}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">// Since we have the type, pretend we have a header.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	gcdata := typ.GCData
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	return typePointers{elem: addr, addr: addr, mask: readUintptr(gcdata), typ: typ}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// nextFast is the fast path of next. nextFast is written to be inlineable and,</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// as the name implies, fast.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// Callers that are performance-critical should iterate using the following</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">// pattern:</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">//	for {</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">//		var addr uintptr</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">//		if tp, addr = tp.nextFast(); addr == 0 {</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">//			if tp, addr = tp.next(limit); addr == 0 {</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">//				break</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span><span class="comment">//			}</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span><span class="comment">//		}</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span><span class="comment">//		// Use addr.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span><span class="comment">//		...</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">// nosplit because it is used during write barriers and must not be preempted.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>func (tp typePointers) nextFast() (typePointers, uintptr) {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// TESTQ/JEQ</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	if tp.mask == 0 {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		return tp, 0
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	<span class="comment">// BSFQ</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	var i int
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	if goarch.PtrSize == 8 {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		i = sys.TrailingZeros64(uint64(tp.mask))
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	} else {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		i = sys.TrailingZeros32(uint32(tp.mask))
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// BTCQ</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	tp.mask ^= uintptr(1) &lt;&lt; (i &amp; (ptrBits - 1))
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// LEAQ (XX)(XX*8)</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	return tp, tp.addr + uintptr(i)*goarch.PtrSize
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span><span class="comment">// next advances the pointers iterator, returning the updated iterator and</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">// the address of the next pointer.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span><span class="comment">// limit must be the same each time it is passed to next.</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span><span class="comment">// nosplit because it is used during write barriers and must not be preempted.</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>func (tp typePointers) next(limit uintptr) (typePointers, uintptr) {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	for {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		if tp.mask != 0 {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			return tp.nextFast()
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		<span class="comment">// Stop if we don&#39;t actually have type information.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		if tp.typ == nil {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			return typePointers{}, 0
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		<span class="comment">// Advance to the next element if necessary.</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		if tp.addr+goarch.PtrSize*ptrBits &gt;= tp.elem+tp.typ.PtrBytes {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			tp.elem += tp.typ.Size_
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			tp.addr = tp.elem
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		} else {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			tp.addr += ptrBits * goarch.PtrSize
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		<span class="comment">// Check if we&#39;ve exceeded the limit with the last update.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		if tp.addr &gt;= limit {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			return typePointers{}, 0
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		<span class="comment">// Grab more bits and try again.</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		tp.mask = readUintptr(addb(tp.typ.GCData, (tp.addr-tp.elem)/goarch.PtrSize/8))
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		if tp.addr+goarch.PtrSize*ptrBits &gt; limit {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			bits := (tp.addr + goarch.PtrSize*ptrBits - limit) / goarch.PtrSize
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			tp.mask &amp;^= ((1 &lt;&lt; (bits)) - 1) &lt;&lt; (ptrBits - bits)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// fastForward moves the iterator forward by n bytes. n must be a multiple</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// of goarch.PtrSize. limit must be the same limit passed to next for this</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// iterator.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">// nosplit because it is used during write barriers and must not be preempted.</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>func (tp typePointers) fastForward(n, limit uintptr) typePointers {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// Basic bounds check.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	target := tp.addr + n
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	if target &gt;= limit {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		return typePointers{}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	if tp.typ == nil {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		<span class="comment">// Handle small objects.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		<span class="comment">// Clear any bits before the target address.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		tp.mask &amp;^= (1 &lt;&lt; ((target - tp.addr) / goarch.PtrSize)) - 1
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		<span class="comment">// Clear any bits past the limit.</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		if tp.addr+goarch.PtrSize*ptrBits &gt; limit {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			bits := (tp.addr + goarch.PtrSize*ptrBits - limit) / goarch.PtrSize
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			tp.mask &amp;^= ((1 &lt;&lt; (bits)) - 1) &lt;&lt; (ptrBits - bits)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		return tp
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// Move up elem and addr.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// Offsets within an element are always at a ptrBits*goarch.PtrSize boundary.</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	if n &gt;= tp.typ.Size_ {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		<span class="comment">// elem needs to be moved to the element containing</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		<span class="comment">// tp.addr + n.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		oldelem := tp.elem
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		tp.elem += (tp.addr - tp.elem + n) / tp.typ.Size_ * tp.typ.Size_
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		tp.addr = tp.elem + alignDown(n-(tp.elem-oldelem), ptrBits*goarch.PtrSize)
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	} else {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		tp.addr += alignDown(n, ptrBits*goarch.PtrSize)
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	if tp.addr-tp.elem &gt;= tp.typ.PtrBytes {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re starting in the non-pointer area of an array.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		<span class="comment">// Move up to the next element.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		tp.elem += tp.typ.Size_
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		tp.addr = tp.elem
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		tp.mask = readUintptr(tp.typ.GCData)
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		<span class="comment">// We may have exceeded the limit after this. Bail just like next does.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		if tp.addr &gt;= limit {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>			return typePointers{}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		}
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	} else {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		<span class="comment">// Grab the mask, but then clear any bits before the target address and any</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		<span class="comment">// bits over the limit.</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		tp.mask = readUintptr(addb(tp.typ.GCData, (tp.addr-tp.elem)/goarch.PtrSize/8))
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		tp.mask &amp;^= (1 &lt;&lt; ((target - tp.addr) / goarch.PtrSize)) - 1
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	if tp.addr+goarch.PtrSize*ptrBits &gt; limit {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		bits := (tp.addr + goarch.PtrSize*ptrBits - limit) / goarch.PtrSize
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		tp.mask &amp;^= ((1 &lt;&lt; (bits)) - 1) &lt;&lt; (ptrBits - bits)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	return tp
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span><span class="comment">// objBase returns the base pointer for the object containing addr in span.</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span><span class="comment">// Assumes that addr points into a valid part of span (span.base() &lt;= addr &lt; span.limit).</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>func (span *mspan) objBase(addr uintptr) uintptr {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	return span.base() + span.objIndex(addr)*span.elemsize
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span><span class="comment">// bulkBarrierPreWrite executes a write barrier</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span><span class="comment">// for every pointer slot in the memory range [src, src+size),</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span><span class="comment">// using pointer/scalar information from [dst, dst+size).</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span><span class="comment">// This executes the write barriers necessary before a memmove.</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span><span class="comment">// src, dst, and size must be pointer-aligned.</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span><span class="comment">// The range [dst, dst+size) must lie within a single object.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span><span class="comment">// It does not perform the actual writes.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span><span class="comment">// As a special case, src == 0 indicates that this is being used for a</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span><span class="comment">// memclr. bulkBarrierPreWrite will pass 0 for the src of each write</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span><span class="comment">// barrier.</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span><span class="comment">// Callers should call bulkBarrierPreWrite immediately before</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">// calling memmove(dst, src, size). This function is marked nosplit</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span><span class="comment">// to avoid being preempted; the GC must not stop the goroutine</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">// between the memmove and the execution of the barriers.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">// The caller is also responsible for cgo pointer checks if this</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// may be writing Go pointers into non-Go memory.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span><span class="comment">// Pointer data is not maintained for allocations containing</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span><span class="comment">// no pointers at all; any caller of bulkBarrierPreWrite must first</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span><span class="comment">// make sure the underlying allocation contains pointers, usually</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span><span class="comment">// by checking typ.PtrBytes.</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span><span class="comment">// The typ argument is the type of the space at src and dst (and the</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span><span class="comment">// element type if src and dst refer to arrays) and it is optional.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span><span class="comment">// If typ is nil, the barrier will still behave as expected and typ</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span><span class="comment">// is used purely as an optimization. However, it must be used with</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span><span class="comment">// care.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span><span class="comment">// If typ is not nil, then src and dst must point to one or more values</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span><span class="comment">// of type typ. The caller must ensure that the ranges [src, src+size)</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span><span class="comment">// and [dst, dst+size) refer to one or more whole values of type src and</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span><span class="comment">// dst (leaving off the pointerless tail of the space is OK). If this</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span><span class="comment">// precondition is not followed, this function will fail to scan the</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span><span class="comment">// right pointers.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span><span class="comment">// When in doubt, pass nil for typ. That is safe and will always work.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span><span class="comment">// Callers must perform cgo checks if goexperiment.CgoCheck2.</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>func bulkBarrierPreWrite(dst, src, size uintptr, typ *abi.Type) {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	if (dst|src|size)&amp;(goarch.PtrSize-1) != 0 {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		throw(&#34;bulkBarrierPreWrite: unaligned arguments&#34;)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	if !writeBarrier.enabled {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		return
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	s := spanOf(dst)
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	if s == nil {
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		<span class="comment">// If dst is a global, use the data or BSS bitmaps to</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		<span class="comment">// execute write barriers.</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		for _, datap := range activeModules() {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			if datap.data &lt;= dst &amp;&amp; dst &lt; datap.edata {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>				bulkBarrierBitmap(dst, src, size, dst-datap.data, datap.gcdatamask.bytedata)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>				return
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>			}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		}
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		for _, datap := range activeModules() {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>			if datap.bss &lt;= dst &amp;&amp; dst &lt; datap.ebss {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>				bulkBarrierBitmap(dst, src, size, dst-datap.bss, datap.gcbssmask.bytedata)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>				return
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>			}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		return
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	} else if s.state.get() != mSpanInUse || dst &lt; s.base() || s.limit &lt;= dst {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		<span class="comment">// dst was heap memory at some point, but isn&#39;t now.</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		<span class="comment">// It can&#39;t be a global. It must be either our stack,</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		<span class="comment">// or in the case of direct channel sends, it could be</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		<span class="comment">// another stack. Either way, no need for barriers.</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		<span class="comment">// This will also catch if dst is in a freed span,</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		<span class="comment">// though that should never have.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		return
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	buf := &amp;getg().m.p.ptr().wbBuf
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	<span class="comment">// Double-check that the bitmaps generated in the two possible paths match.</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	const doubleCheck = false
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	if doubleCheck {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		doubleCheckTypePointersOfType(s, typ, dst, size)
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	var tp typePointers
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	if typ != nil &amp;&amp; typ.Kind_&amp;kindGCProg == 0 {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		tp = s.typePointersOfType(typ, dst)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	} else {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		tp = s.typePointersOf(dst, size)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	if src == 0 {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		for {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>			var addr uintptr
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>			if tp, addr = tp.next(dst + size); addr == 0 {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>				break
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>			dstx := (*uintptr)(unsafe.Pointer(addr))
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			p := buf.get1()
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			p[0] = *dstx
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	} else {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		for {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			var addr uintptr
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			if tp, addr = tp.next(dst + size); addr == 0 {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>				break
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			dstx := (*uintptr)(unsafe.Pointer(addr))
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			srcx := (*uintptr)(unsafe.Pointer(src + (addr - dst)))
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			p := buf.get2()
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			p[0] = *dstx
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>			p[1] = *srcx
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		}
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span><span class="comment">// bulkBarrierPreWriteSrcOnly is like bulkBarrierPreWrite but</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span><span class="comment">// does not execute write barriers for [dst, dst+size).</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span><span class="comment">// In addition to the requirements of bulkBarrierPreWrite</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span><span class="comment">// callers need to ensure [dst, dst+size) is zeroed.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span><span class="comment">// This is used for special cases where e.g. dst was just</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span><span class="comment">// created and zeroed with malloc.</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span><span class="comment">// The type of the space can be provided purely as an optimization.</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span><span class="comment">// See bulkBarrierPreWrite&#39;s comment for more details -- use this</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span><span class="comment">// optimization with great care.</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>func bulkBarrierPreWriteSrcOnly(dst, src, size uintptr, typ *abi.Type) {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	if (dst|src|size)&amp;(goarch.PtrSize-1) != 0 {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		throw(&#34;bulkBarrierPreWrite: unaligned arguments&#34;)
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	if !writeBarrier.enabled {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		return
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	buf := &amp;getg().m.p.ptr().wbBuf
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	s := spanOf(dst)
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	<span class="comment">// Double-check that the bitmaps generated in the two possible paths match.</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	const doubleCheck = false
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	if doubleCheck {
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		doubleCheckTypePointersOfType(s, typ, dst, size)
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	var tp typePointers
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	if typ != nil &amp;&amp; typ.Kind_&amp;kindGCProg == 0 {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		tp = s.typePointersOfType(typ, dst)
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	} else {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		tp = s.typePointersOf(dst, size)
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	}
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	for {
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		var addr uintptr
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		if tp, addr = tp.next(dst + size); addr == 0 {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>			break
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		srcx := (*uintptr)(unsafe.Pointer(addr - dst + src))
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		p := buf.get1()
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		p[0] = *srcx
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">// initHeapBits initializes the heap bitmap for a span.</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">// TODO(mknyszek): This should set the heap bits for single pointer</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">// allocations eagerly to avoid calling heapSetType at allocation time,</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">// just to write one bit.</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>func (s *mspan) initHeapBits(forceClear bool) {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	if (!s.spanclass.noscan() &amp;&amp; heapBitsInSpan(s.elemsize)) || s.isUserArenaChunk {
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		b := s.heapBits()
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		for i := range b {
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>			b[i] = 0
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		}
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">// bswapIfBigEndian swaps the byte order of the uintptr on goarch.BigEndian platforms,</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span><span class="comment">// and leaves it alone elsewhere.</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>func bswapIfBigEndian(x uintptr) uintptr {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	if goarch.BigEndian {
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		if goarch.PtrSize == 8 {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			return uintptr(sys.Bswap64(uint64(x)))
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		}
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		return uintptr(sys.Bswap32(uint32(x)))
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	return x
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>}
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>type writeUserArenaHeapBits struct {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	offset uintptr <span class="comment">// offset in span that the low bit of mask represents the pointer state of.</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	mask   uintptr <span class="comment">// some pointer bits starting at the address addr.</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	valid  uintptr <span class="comment">// number of bits in buf that are valid (including low)</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	low    uintptr <span class="comment">// number of low-order bits to not overwrite</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>func (s *mspan) writeUserArenaHeapBits(addr uintptr) (h writeUserArenaHeapBits) {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	offset := addr - s.base()
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	<span class="comment">// We start writing bits maybe in the middle of a heap bitmap word.</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	<span class="comment">// Remember how many bits into the word we started, so we can be sure</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	<span class="comment">// not to overwrite the previous bits.</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	h.low = offset / goarch.PtrSize % ptrBits
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	<span class="comment">// round down to heap word that starts the bitmap word.</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	h.offset = offset - h.low*goarch.PtrSize
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	<span class="comment">// We don&#39;t have any bits yet.</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	h.mask = 0
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	h.valid = h.low
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	return
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span><span class="comment">// write appends the pointerness of the next valid pointer slots</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span><span class="comment">// using the low valid bits of bits. 1=pointer, 0=scalar.</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>func (h writeUserArenaHeapBits) write(s *mspan, bits, valid uintptr) writeUserArenaHeapBits {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	if h.valid+valid &lt;= ptrBits {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		<span class="comment">// Fast path - just accumulate the bits.</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		h.mask |= bits &lt;&lt; h.valid
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		h.valid += valid
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		return h
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	<span class="comment">// Too many bits to fit in this word. Write the current word</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	<span class="comment">// out and move on to the next word.</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	data := h.mask | bits&lt;&lt;h.valid       <span class="comment">// mask for this word</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	h.mask = bits &gt;&gt; (ptrBits - h.valid) <span class="comment">// leftover for next word</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	h.valid += valid - ptrBits           <span class="comment">// have h.valid+valid bits, writing ptrBits of them</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	<span class="comment">// Flush mask to the memory bitmap.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	idx := h.offset / (ptrBits * goarch.PtrSize)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	m := uintptr(1)&lt;&lt;h.low - 1
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	bitmap := s.heapBits()
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	bitmap[idx] = bswapIfBigEndian(bswapIfBigEndian(bitmap[idx])&amp;m | data)
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	<span class="comment">// Note: no synchronization required for this write because</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	<span class="comment">// the allocator has exclusive access to the page, and the bitmap</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	<span class="comment">// entries are all for a single page. Also, visibility of these</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	<span class="comment">// writes is guaranteed by the publication barrier in mallocgc.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	<span class="comment">// Move to next word of bitmap.</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	h.offset += ptrBits * goarch.PtrSize
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	h.low = 0
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	return h
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>}
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span><span class="comment">// Add padding of size bytes.</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>func (h writeUserArenaHeapBits) pad(s *mspan, size uintptr) writeUserArenaHeapBits {
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	if size == 0 {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		return h
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	words := size / goarch.PtrSize
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	for words &gt; ptrBits {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		h = h.write(s, 0, ptrBits)
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		words -= ptrBits
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	}
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	return h.write(s, 0, words)
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>}
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span><span class="comment">// Flush the bits that have been written, and add zeros as needed</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span><span class="comment">// to cover the full object [addr, addr+size).</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>func (h writeUserArenaHeapBits) flush(s *mspan, addr, size uintptr) {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	offset := addr - s.base()
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	<span class="comment">// zeros counts the number of bits needed to represent the object minus the</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	<span class="comment">// number of bits we&#39;ve already written. This is the number of 0 bits</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	<span class="comment">// that need to be added.</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	zeros := (offset+size-h.offset)/goarch.PtrSize - h.valid
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	<span class="comment">// Add zero bits up to the bitmap word boundary</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	if zeros &gt; 0 {
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		z := ptrBits - h.valid
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		if z &gt; zeros {
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>			z = zeros
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		}
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		h.valid += z
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		zeros -= z
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	}
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	<span class="comment">// Find word in bitmap that we&#39;re going to write.</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	bitmap := s.heapBits()
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	idx := h.offset / (ptrBits * goarch.PtrSize)
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	<span class="comment">// Write remaining bits.</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	if h.valid != h.low {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		m := uintptr(1)&lt;&lt;h.low - 1      <span class="comment">// don&#39;t clear existing bits below &#34;low&#34;</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		m |= ^(uintptr(1)&lt;&lt;h.valid - 1) <span class="comment">// don&#39;t clear existing bits above &#34;valid&#34;</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		bitmap[idx] = bswapIfBigEndian(bswapIfBigEndian(bitmap[idx])&amp;m | h.mask)
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	}
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	if zeros == 0 {
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		return
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	}
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	<span class="comment">// Advance to next bitmap word.</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	h.offset += ptrBits * goarch.PtrSize
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	<span class="comment">// Continue on writing zeros for the rest of the object.</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	<span class="comment">// For standard use of the ptr bits this is not required, as</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	<span class="comment">// the bits are read from the beginning of the object. Some uses,</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	<span class="comment">// like noscan spans, oblets, bulk write barriers, and cgocheck, might</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	<span class="comment">// start mid-object, so these writes are still required.</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	for {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		<span class="comment">// Write zero bits.</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		idx := h.offset / (ptrBits * goarch.PtrSize)
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		if zeros &lt; ptrBits {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>			bitmap[idx] = bswapIfBigEndian(bswapIfBigEndian(bitmap[idx]) &amp;^ (uintptr(1)&lt;&lt;zeros - 1))
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			break
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		} else if zeros == ptrBits {
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>			bitmap[idx] = 0
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>			break
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		} else {
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>			bitmap[idx] = 0
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>			zeros -= ptrBits
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		}
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		h.offset += ptrBits * goarch.PtrSize
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	}
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>}
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span><span class="comment">// heapBits returns the heap ptr/scalar bits stored at the end of the span for</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span><span class="comment">// small object spans and heap arena spans.</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span><span class="comment">// Note that the uintptr of each element means something different for small object</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span><span class="comment">// spans and for heap arena spans. Small object spans are easy: they&#39;re never interpreted</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span><span class="comment">// as anything but uintptr, so they&#39;re immune to differences in endianness. However, the</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span><span class="comment">// heapBits for user arena spans is exposed through a dummy type descriptor, so the byte</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span><span class="comment">// ordering needs to match the same byte ordering the compiler would emit. The compiler always</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span><span class="comment">// emits the bitmap data in little endian byte ordering, so on big endian platforms these</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span><span class="comment">// uintptrs will have their byte orders swapped from what they normally would be.</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span><span class="comment">// heapBitsInSpan(span.elemsize) or span.isUserArenaChunk must be true.</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>func (span *mspan) heapBits() []uintptr {
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	const doubleCheck = false
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	if doubleCheck &amp;&amp; !span.isUserArenaChunk {
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>		if span.spanclass.noscan() {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>			throw(&#34;heapBits called for noscan&#34;)
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		}
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		if span.elemsize &gt; minSizeForMallocHeader {
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>			throw(&#34;heapBits called for span class that should have a malloc header&#34;)
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		}
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	}
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	<span class="comment">// Find the bitmap at the end of the span.</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	<span class="comment">// Nearly every span with heap bits is exactly one page in size. Arenas are the only exception.</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	if span.npages == 1 {
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		<span class="comment">// This will be inlined and constant-folded down.</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		return heapBitsSlice(span.base(), pageSize)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	}
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	return heapBitsSlice(span.base(), span.npages*pageSize)
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>}
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span><span class="comment">// Helper for constructing a slice for the span&#39;s heap bits.</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>func heapBitsSlice(spanBase, spanSize uintptr) []uintptr {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	bitmapSize := spanSize / goarch.PtrSize / 8
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	elems := int(bitmapSize / goarch.PtrSize)
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	var sl notInHeapSlice
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	sl = notInHeapSlice{(*notInHeap)(unsafe.Pointer(spanBase + spanSize - bitmapSize)), elems, elems}
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	return *(*[]uintptr)(unsafe.Pointer(&amp;sl))
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>}
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span><span class="comment">// heapBitsSmallForAddr loads the heap bits for the object stored at addr from span.heapBits.</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span><span class="comment">// addr must be the base pointer of an object in the span. heapBitsInSpan(span.elemsize)</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span><span class="comment">// must be true.</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>func (span *mspan) heapBitsSmallForAddr(addr uintptr) uintptr {
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	spanSize := span.npages * pageSize
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	bitmapSize := spanSize / goarch.PtrSize / 8
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	hbits := (*byte)(unsafe.Pointer(span.base() + spanSize - bitmapSize))
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	<span class="comment">// These objects are always small enough that their bitmaps</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	<span class="comment">// fit in a single word, so just load the word or two we need.</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	<span class="comment">// Mirrors mspan.writeHeapBitsSmall.</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	<span class="comment">// We should be using heapBits(), but unfortunately it introduces</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	<span class="comment">// both bounds checks panics and throw which causes us to exceed</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	<span class="comment">// the nosplit limit in quite a few cases.</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	i := (addr - span.base()) / goarch.PtrSize / ptrBits
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	j := (addr - span.base()) / goarch.PtrSize % ptrBits
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	bits := span.elemsize / goarch.PtrSize
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	word0 := (*uintptr)(unsafe.Pointer(addb(hbits, goarch.PtrSize*(i+0))))
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	word1 := (*uintptr)(unsafe.Pointer(addb(hbits, goarch.PtrSize*(i+1))))
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	var read uintptr
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	if j+bits &gt; ptrBits {
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>		<span class="comment">// Two reads.</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>		bits0 := ptrBits - j
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		bits1 := bits - bits0
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		read = *word0 &gt;&gt; j
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>		read |= (*word1 &amp; ((1 &lt;&lt; bits1) - 1)) &lt;&lt; bits0
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	} else {
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		<span class="comment">// One read.</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>		read = (*word0 &gt;&gt; j) &amp; ((1 &lt;&lt; bits) - 1)
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>	}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	return read
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>}
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span><span class="comment">// writeHeapBitsSmall writes the heap bits for small objects whose ptr/scalar data is</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span><span class="comment">// stored as a bitmap at the end of the span.</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span><span class="comment">// Assumes dataSize is &lt;= ptrBits*goarch.PtrSize. x must be a pointer into the span.</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span><span class="comment">// heapBitsInSpan(dataSize) must be true. dataSize must be &gt;= typ.Size_.</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>func (span *mspan) writeHeapBitsSmall(x, dataSize uintptr, typ *_type) (scanSize uintptr) {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	<span class="comment">// The objects here are always really small, so a single load is sufficient.</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>	src0 := readUintptr(typ.GCData)
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	<span class="comment">// Create repetitions of the bitmap if we have a small array.</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	bits := span.elemsize / goarch.PtrSize
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	scanSize = typ.PtrBytes
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	src := src0
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	switch typ.Size_ {
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	case goarch.PtrSize:
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		src = (1 &lt;&lt; (dataSize / goarch.PtrSize)) - 1
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	default:
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		for i := typ.Size_; i &lt; dataSize; i += typ.Size_ {
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>			src |= src0 &lt;&lt; (i / goarch.PtrSize)
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			scanSize += typ.Size_
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		}
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	}
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	<span class="comment">// Since we&#39;re never writing more than one uintptr&#39;s worth of bits, we&#39;re either going</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	<span class="comment">// to do one or two writes.</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	dst := span.heapBits()
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	o := (x - span.base()) / goarch.PtrSize
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	i := o / ptrBits
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	j := o % ptrBits
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	if j+bits &gt; ptrBits {
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		<span class="comment">// Two writes.</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>		bits0 := ptrBits - j
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		bits1 := bits - bits0
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		dst[i+0] = dst[i+0]&amp;(^uintptr(0)&gt;&gt;bits0) | (src &lt;&lt; j)
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>		dst[i+1] = dst[i+1]&amp;^((1&lt;&lt;bits1)-1) | (src &gt;&gt; bits0)
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	} else {
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>		<span class="comment">// One write.</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		dst[i] = (dst[i] &amp;^ (((1 &lt;&lt; bits) - 1) &lt;&lt; j)) | (src &lt;&lt; j)
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	}
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	const doubleCheck = false
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	if doubleCheck {
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		srcRead := span.heapBitsSmallForAddr(x)
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		if srcRead != src {
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>			print(&#34;runtime: x=&#34;, hex(x), &#34; i=&#34;, i, &#34; j=&#34;, j, &#34; bits=&#34;, bits, &#34;\n&#34;)
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>			print(&#34;runtime: dataSize=&#34;, dataSize, &#34; typ.Size_=&#34;, typ.Size_, &#34; typ.PtrBytes=&#34;, typ.PtrBytes, &#34;\n&#34;)
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>			print(&#34;runtime: src0=&#34;, hex(src0), &#34; src=&#34;, hex(src), &#34; srcRead=&#34;, hex(srcRead), &#34;\n&#34;)
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>			throw(&#34;bad pointer bits written for small object&#34;)
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		}
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	}
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	return
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>}
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span><span class="comment">// For !goexperiment.AllocHeaders.</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>func heapBitsSetType(x, size, dataSize uintptr, typ *_type) {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>}
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span><span class="comment">// heapSetType records that the new allocation [x, x+size)</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span><span class="comment">// holds in [x, x+dataSize) one or more values of type typ.</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span><span class="comment">// (The number of values is given by dataSize / typ.Size.)</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span><span class="comment">// If dataSize &lt; size, the fragment [x+dataSize, x+size) is</span>
<span id="L842" class="ln">   842&nbsp;&nbsp;</span><span class="comment">// recorded as non-pointer data.</span>
<span id="L843" class="ln">   843&nbsp;&nbsp;</span><span class="comment">// It is known that the type has pointers somewhere;</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span><span class="comment">// malloc does not call heapSetType when there are no pointers.</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span><span class="comment">// There can be read-write races between heapSetType and things</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span><span class="comment">// that read the heap metadata like scanobject. However, since</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span><span class="comment">// heapSetType is only used for objects that have not yet been</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span><span class="comment">// made reachable, readers will ignore bits being modified by this</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span><span class="comment">// function. This does mean this function cannot transiently modify</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span><span class="comment">// shared memory that belongs to neighboring objects. Also, on weakly-ordered</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span><span class="comment">// machines, callers must execute a store/store (publication) barrier</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span><span class="comment">// between calling this function and making the object reachable.</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>func heapSetType(x, dataSize uintptr, typ *_type, header **_type, span *mspan) (scanSize uintptr) {
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	const doubleCheck = false
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	gctyp := typ
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	if header == nil {
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		if doubleCheck &amp;&amp; (!heapBitsInSpan(dataSize) || !heapBitsInSpan(span.elemsize)) {
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>			throw(&#34;tried to write heap bits, but no heap bits in span&#34;)
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		}
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		<span class="comment">// Handle the case where we have no malloc header.</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		scanSize = span.writeHeapBitsSmall(x, dataSize, typ)
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	} else {
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		if typ.Kind_&amp;kindGCProg != 0 {
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>			<span class="comment">// Allocate space to unroll the gcprog. This space will consist of</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>			<span class="comment">// a dummy _type value and the unrolled gcprog. The dummy _type will</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>			<span class="comment">// refer to the bitmap, and the mspan will refer to the dummy _type.</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>			if span.spanclass.sizeclass() != 0 {
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>				throw(&#34;GCProg for type that isn&#39;t large&#34;)
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>			}
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>			spaceNeeded := alignUp(unsafe.Sizeof(_type{}), goarch.PtrSize)
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>			heapBitsOff := spaceNeeded
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>			spaceNeeded += alignUp(typ.PtrBytes/goarch.PtrSize/8, goarch.PtrSize)
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>			npages := alignUp(spaceNeeded, pageSize) / pageSize
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>			var progSpan *mspan
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>			systemstack(func() {
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>				progSpan = mheap_.allocManual(npages, spanAllocPtrScalarBits)
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>				memclrNoHeapPointers(unsafe.Pointer(progSpan.base()), progSpan.npages*pageSize)
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>			})
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>			<span class="comment">// Write a dummy _type in the new space.</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>			<span class="comment">// We only need to write size, PtrBytes, and GCData, since that&#39;s all</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>			<span class="comment">// the GC cares about.</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>			gctyp = (*_type)(unsafe.Pointer(progSpan.base()))
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>			gctyp.Size_ = typ.Size_
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>			gctyp.PtrBytes = typ.PtrBytes
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>			gctyp.GCData = (*byte)(add(unsafe.Pointer(progSpan.base()), heapBitsOff))
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			gctyp.TFlag = abi.TFlagUnrolledBitmap
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			<span class="comment">// Expand the GC program into space reserved at the end of the new span.</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>			runGCProg(addb(typ.GCData, 4), gctyp.GCData)
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>		}
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		<span class="comment">// Write out the header.</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>		*header = gctyp
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>		scanSize = span.elemsize
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	}
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>	if doubleCheck {
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		doubleCheckHeapPointers(x, dataSize, gctyp, header, span)
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>		<span class="comment">// To exercise the less common path more often, generate</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>		<span class="comment">// a random interior pointer and make sure iterating from</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>		<span class="comment">// that point works correctly too.</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>		maxIterBytes := span.elemsize
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>		if header == nil {
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>			maxIterBytes = dataSize
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>		}
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>		off := alignUp(uintptr(cheaprand())%dataSize, goarch.PtrSize)
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>		size := dataSize - off
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>		if size == 0 {
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>			off -= goarch.PtrSize
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>			size += goarch.PtrSize
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		}
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>		interior := x + off
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		size -= alignDown(uintptr(cheaprand())%size, goarch.PtrSize)
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		if size == 0 {
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>			size = goarch.PtrSize
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>		}
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		<span class="comment">// Round up the type to the size of the type.</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>		size = (size + gctyp.Size_ - 1) / gctyp.Size_ * gctyp.Size_
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>		if interior+size &gt; x+maxIterBytes {
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>			size = x + maxIterBytes - interior
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>		}
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>		doubleCheckHeapPointersInterior(x, interior, size, dataSize, gctyp, header, span)
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	}
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	return
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>}
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>func doubleCheckHeapPointers(x, dataSize uintptr, typ *_type, header **_type, span *mspan) {
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>	<span class="comment">// Check that scanning the full object works.</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>	tp := span.typePointersOfUnchecked(span.objBase(x))
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	maxIterBytes := span.elemsize
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>	if header == nil {
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>		maxIterBytes = dataSize
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>	}
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	bad := false
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; maxIterBytes; i += goarch.PtrSize {
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>		<span class="comment">// Compute the pointer bit we want at offset i.</span>
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		want := false
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		if i &lt; span.elemsize {
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>			off := i % typ.Size_
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>			if off &lt; typ.PtrBytes {
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>				j := off / goarch.PtrSize
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>				want = *addb(typ.GCData, j/8)&gt;&gt;(j%8)&amp;1 != 0
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>			}
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>		}
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		if want {
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>			var addr uintptr
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>			tp, addr = tp.next(x + span.elemsize)
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			if addr == 0 {
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>				println(&#34;runtime: found bad iterator&#34;)
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>			}
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>			if addr != x+i {
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>				print(&#34;runtime: addr=&#34;, hex(addr), &#34; x+i=&#34;, hex(x+i), &#34;\n&#34;)
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>				bad = true
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>			}
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		}
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>	}
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	if !bad {
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>		var addr uintptr
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>		tp, addr = tp.next(x + span.elemsize)
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>		if addr == 0 {
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>			return
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		}
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>		println(&#34;runtime: extra pointer:&#34;, hex(addr))
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>	}
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>	print(&#34;runtime: hasHeader=&#34;, header != nil, &#34; typ.Size_=&#34;, typ.Size_, &#34; hasGCProg=&#34;, typ.Kind_&amp;kindGCProg != 0, &#34;\n&#34;)
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	print(&#34;runtime: x=&#34;, hex(x), &#34; dataSize=&#34;, dataSize, &#34; elemsize=&#34;, span.elemsize, &#34;\n&#34;)
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>	print(&#34;runtime: typ=&#34;, unsafe.Pointer(typ), &#34; typ.PtrBytes=&#34;, typ.PtrBytes, &#34;\n&#34;)
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>	print(&#34;runtime: limit=&#34;, hex(x+span.elemsize), &#34;\n&#34;)
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	tp = span.typePointersOfUnchecked(x)
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>	dumpTypePointers(tp)
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>	for {
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>		var addr uintptr
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>		if tp, addr = tp.next(x + span.elemsize); addr == 0 {
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>			println(&#34;runtime: would&#39;ve stopped here&#34;)
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>			dumpTypePointers(tp)
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>			break
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>		}
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		print(&#34;runtime: addr=&#34;, hex(addr), &#34;\n&#34;)
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		dumpTypePointers(tp)
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>	}
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	throw(&#34;heapSetType: pointer entry not correct&#34;)
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>}
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>func doubleCheckHeapPointersInterior(x, interior, size, dataSize uintptr, typ *_type, header **_type, span *mspan) {
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	bad := false
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	if interior &lt; x {
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>		print(&#34;runtime: interior=&#34;, hex(interior), &#34; x=&#34;, hex(x), &#34;\n&#34;)
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>		throw(&#34;found bad interior pointer&#34;)
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	}
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	off := interior - x
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>	tp := span.typePointersOf(interior, size)
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	for i := off; i &lt; off+size; i += goarch.PtrSize {
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>		<span class="comment">// Compute the pointer bit we want at offset i.</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>		want := false
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>		if i &lt; span.elemsize {
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>			off := i % typ.Size_
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>			if off &lt; typ.PtrBytes {
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>				j := off / goarch.PtrSize
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>				want = *addb(typ.GCData, j/8)&gt;&gt;(j%8)&amp;1 != 0
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>			}
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>		}
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>		if want {
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>			var addr uintptr
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>			tp, addr = tp.next(interior + size)
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>			if addr == 0 {
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>				println(&#34;runtime: found bad iterator&#34;)
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>				bad = true
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>			}
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>			if addr != x+i {
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>				print(&#34;runtime: addr=&#34;, hex(addr), &#34; x+i=&#34;, hex(x+i), &#34;\n&#34;)
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>				bad = true
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>			}
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>		}
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	}
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>	if !bad {
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>		var addr uintptr
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		tp, addr = tp.next(interior + size)
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>		if addr == 0 {
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>			return
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>		}
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		println(&#34;runtime: extra pointer:&#34;, hex(addr))
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	}
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>	print(&#34;runtime: hasHeader=&#34;, header != nil, &#34; typ.Size_=&#34;, typ.Size_, &#34;\n&#34;)
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	print(&#34;runtime: x=&#34;, hex(x), &#34; dataSize=&#34;, dataSize, &#34; elemsize=&#34;, span.elemsize, &#34; interior=&#34;, hex(interior), &#34; size=&#34;, size, &#34;\n&#34;)
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	print(&#34;runtime: limit=&#34;, hex(interior+size), &#34;\n&#34;)
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>	tp = span.typePointersOf(interior, size)
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	dumpTypePointers(tp)
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	for {
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		var addr uintptr
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>		if tp, addr = tp.next(interior + size); addr == 0 {
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>			println(&#34;runtime: would&#39;ve stopped here&#34;)
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>			dumpTypePointers(tp)
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>			break
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>		}
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>		print(&#34;runtime: addr=&#34;, hex(addr), &#34;\n&#34;)
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>		dumpTypePointers(tp)
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>	}
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>	print(&#34;runtime: want: &#34;)
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>	for i := off; i &lt; off+size; i += goarch.PtrSize {
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>		<span class="comment">// Compute the pointer bit we want at offset i.</span>
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>		want := false
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>		if i &lt; dataSize {
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>			off := i % typ.Size_
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>			if off &lt; typ.PtrBytes {
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>				j := off / goarch.PtrSize
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>				want = *addb(typ.GCData, j/8)&gt;&gt;(j%8)&amp;1 != 0
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>			}
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>		}
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>		if want {
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>			print(&#34;1&#34;)
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>		} else {
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>			print(&#34;0&#34;)
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>		}
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>	}
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>	println()
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>	throw(&#34;heapSetType: pointer entry not correct&#34;)
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>}
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>func doubleCheckTypePointersOfType(s *mspan, typ *_type, addr, size uintptr) {
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>	if typ == nil || typ.Kind_&amp;kindGCProg != 0 {
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>		return
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>	}
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	if typ.Kind_&amp;kindMask == kindInterface {
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>		<span class="comment">// Interfaces are unfortunately inconsistently handled</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>		<span class="comment">// when it comes to the type pointer, so it&#39;s easy to</span>
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>		<span class="comment">// produce a lot of false positives here.</span>
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>		return
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	}
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>	tp0 := s.typePointersOfType(typ, addr)
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	tp1 := s.typePointersOf(addr, size)
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>	failed := false
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	for {
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>		var addr0, addr1 uintptr
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>		tp0, addr0 = tp0.next(addr + size)
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>		tp1, addr1 = tp1.next(addr + size)
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>		if addr0 != addr1 {
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>			failed = true
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>			break
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>		}
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>		if addr0 == 0 {
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>			break
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>		}
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>	}
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>	if failed {
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>		tp0 := s.typePointersOfType(typ, addr)
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>		tp1 := s.typePointersOf(addr, size)
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>		print(&#34;runtime: addr=&#34;, hex(addr), &#34; size=&#34;, size, &#34;\n&#34;)
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>		print(&#34;runtime: type=&#34;, toRType(typ).string(), &#34;\n&#34;)
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>		dumpTypePointers(tp0)
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>		dumpTypePointers(tp1)
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>		for {
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>			var addr0, addr1 uintptr
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>			tp0, addr0 = tp0.next(addr + size)
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>			tp1, addr1 = tp1.next(addr + size)
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>			print(&#34;runtime: &#34;, hex(addr0), &#34; &#34;, hex(addr1), &#34;\n&#34;)
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>			if addr0 == 0 &amp;&amp; addr1 == 0 {
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>				break
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>			}
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>		}
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>		throw(&#34;mismatch between typePointersOfType and typePointersOf&#34;)
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>	}
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>}
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>func dumpTypePointers(tp typePointers) {
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>	print(&#34;runtime: tp.elem=&#34;, hex(tp.elem), &#34; tp.typ=&#34;, unsafe.Pointer(tp.typ), &#34;\n&#34;)
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>	print(&#34;runtime: tp.addr=&#34;, hex(tp.addr), &#34; tp.mask=&#34;)
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; ptrBits; i++ {
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>		if tp.mask&amp;(uintptr(1)&lt;&lt;i) != 0 {
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>			print(&#34;1&#34;)
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>		} else {
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>			print(&#34;0&#34;)
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>		}
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>	}
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>	println()
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>}
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span><span class="comment">// Testing.</span>
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span><span class="comment">// Returns GC type info for the pointer stored in ep for testing.</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span><span class="comment">// If ep points to the stack, only static live information will be returned</span>
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span><span class="comment">// (i.e. not for objects which are only dynamically live stack objects).</span>
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>func getgcmask(ep any) (mask []byte) {
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>	e := *efaceOf(&amp;ep)
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>	p := e.data
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>	t := e._type
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>	var et *_type
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>	if t.Kind_&amp;kindMask != kindPtr {
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>		throw(&#34;bad argument to getgcmask: expected type to be a pointer to the value type whose mask is being queried&#34;)
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>	}
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>	et = (*ptrtype)(unsafe.Pointer(t)).Elem
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>	<span class="comment">// data or bss</span>
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	for _, datap := range activeModules() {
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>		<span class="comment">// data</span>
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>		if datap.data &lt;= uintptr(p) &amp;&amp; uintptr(p) &lt; datap.edata {
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>			bitmap := datap.gcdatamask.bytedata
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>			n := et.Size_
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>			mask = make([]byte, n/goarch.PtrSize)
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>			for i := uintptr(0); i &lt; n; i += goarch.PtrSize {
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>				off := (uintptr(p) + i - datap.data) / goarch.PtrSize
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>				mask[i/goarch.PtrSize] = (*addb(bitmap, off/8) &gt;&gt; (off % 8)) &amp; 1
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>			}
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>			return
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>		}
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>		<span class="comment">// bss</span>
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>		if datap.bss &lt;= uintptr(p) &amp;&amp; uintptr(p) &lt; datap.ebss {
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>			bitmap := datap.gcbssmask.bytedata
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>			n := et.Size_
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>			mask = make([]byte, n/goarch.PtrSize)
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>			for i := uintptr(0); i &lt; n; i += goarch.PtrSize {
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>				off := (uintptr(p) + i - datap.bss) / goarch.PtrSize
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>				mask[i/goarch.PtrSize] = (*addb(bitmap, off/8) &gt;&gt; (off % 8)) &amp; 1
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>			}
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>			return
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>		}
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>	}
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>	<span class="comment">// heap</span>
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>	if base, s, _ := findObject(uintptr(p), 0, 0); base != 0 {
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>		if s.spanclass.noscan() {
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>			return nil
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>		}
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>		limit := base + s.elemsize
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>		<span class="comment">// Move the base up to the iterator&#39;s start, because</span>
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>		<span class="comment">// we want to hide evidence of a malloc header from the</span>
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>		<span class="comment">// caller.</span>
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>		tp := s.typePointersOfUnchecked(base)
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>		base = tp.addr
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>		<span class="comment">// Unroll the full bitmap the GC would actually observe.</span>
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>		maskFromHeap := make([]byte, (limit-base)/goarch.PtrSize)
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>		for {
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>			var addr uintptr
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>			if tp, addr = tp.next(limit); addr == 0 {
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>				break
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>			}
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>			maskFromHeap[(addr-base)/goarch.PtrSize] = 1
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>		}
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>		<span class="comment">// Double-check that every part of the ptr/scalar we&#39;re not</span>
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>		<span class="comment">// showing the caller is zeroed. This keeps us honest that</span>
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>		<span class="comment">// that information is actually irrelevant.</span>
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>		for i := limit; i &lt; s.elemsize; i++ {
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>			if *(*byte)(unsafe.Pointer(i)) != 0 {
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>				throw(&#34;found non-zeroed tail of allocation&#34;)
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>			}
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>		}
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>		<span class="comment">// Callers (and a check we&#39;re about to run) expects this mask</span>
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>		<span class="comment">// to end at the last pointer.</span>
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>		for len(maskFromHeap) &gt; 0 &amp;&amp; maskFromHeap[len(maskFromHeap)-1] == 0 {
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>			maskFromHeap = maskFromHeap[:len(maskFromHeap)-1]
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>		}
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>		if et.Kind_&amp;kindGCProg == 0 {
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>			<span class="comment">// Unroll again, but this time from the type information.</span>
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>			maskFromType := make([]byte, (limit-base)/goarch.PtrSize)
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>			tp = s.typePointersOfType(et, base)
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>			for {
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>				var addr uintptr
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>				if tp, addr = tp.next(limit); addr == 0 {
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>					break
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>				}
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>				maskFromType[(addr-base)/goarch.PtrSize] = 1
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>			}
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>			<span class="comment">// Validate that the prefix of maskFromType is equal to</span>
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>			<span class="comment">// maskFromHeap. maskFromType may contain more pointers than</span>
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>			<span class="comment">// maskFromHeap produces because maskFromHeap may be able to</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>			<span class="comment">// get exact type information for certain classes of objects.</span>
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>			<span class="comment">// With maskFromType, we&#39;re always just tiling the type bitmap</span>
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>			<span class="comment">// through to the elemsize.</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>			<span class="comment">// It&#39;s OK if maskFromType has pointers in elemsize that extend</span>
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>			<span class="comment">// past the actual populated space; we checked above that all</span>
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>			<span class="comment">// that space is zeroed, so just the GC will just see nil pointers.</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>			differs := false
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>			for i := range maskFromHeap {
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>				if maskFromHeap[i] != maskFromType[i] {
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>					differs = true
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>					break
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>				}
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>			}
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>			if differs {
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>				print(&#34;runtime: heap mask=&#34;)
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>				for _, b := range maskFromHeap {
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>					print(b)
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>				}
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>				println()
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>				print(&#34;runtime: type mask=&#34;)
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>				for _, b := range maskFromType {
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>					print(b)
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>				}
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>				println()
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>				print(&#34;runtime: type=&#34;, toRType(et).string(), &#34;\n&#34;)
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>				throw(&#34;found two different masks from two different methods&#34;)
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>			}
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>		}
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>		<span class="comment">// Select the heap mask to return. We may not have a type mask.</span>
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>		mask = maskFromHeap
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>		<span class="comment">// Make sure we keep ep alive. We may have stopped referencing</span>
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>		<span class="comment">// ep&#39;s data pointer sometime before this point and it&#39;s possible</span>
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>		<span class="comment">// for that memory to get freed.</span>
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>		KeepAlive(ep)
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>		return
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>	}
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>	<span class="comment">// stack</span>
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>	if gp := getg(); gp.m.curg.stack.lo &lt;= uintptr(p) &amp;&amp; uintptr(p) &lt; gp.m.curg.stack.hi {
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>		found := false
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>		var u unwinder
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>		for u.initAt(gp.m.curg.sched.pc, gp.m.curg.sched.sp, 0, gp.m.curg, 0); u.valid(); u.next() {
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>			if u.frame.sp &lt;= uintptr(p) &amp;&amp; uintptr(p) &lt; u.frame.varp {
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>				found = true
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>				break
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>			}
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>		}
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>		if found {
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>			locals, _, _ := u.frame.getStackMap(false)
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>			if locals.n == 0 {
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>				return
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>			}
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>			size := uintptr(locals.n) * goarch.PtrSize
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>			n := (*ptrtype)(unsafe.Pointer(t)).Elem.Size_
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>			mask = make([]byte, n/goarch.PtrSize)
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>			for i := uintptr(0); i &lt; n; i += goarch.PtrSize {
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>				off := (uintptr(p) + i - u.frame.varp + size) / goarch.PtrSize
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>				mask[i/goarch.PtrSize] = locals.ptrbit(off)
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>			}
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>		}
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>		return
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>	}
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>	<span class="comment">// otherwise, not something the GC knows about.</span>
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>	<span class="comment">// possibly read-only data, like malloc(0).</span>
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>	<span class="comment">// must not have pointers</span>
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>	return
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>}
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span><span class="comment">// userArenaHeapBitsSetType is the equivalent of heapSetType but for</span>
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span><span class="comment">// non-slice-backing-store Go values allocated in a user arena chunk. It</span>
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span><span class="comment">// sets up the type metadata for the value with type typ allocated at address ptr.</span>
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span><span class="comment">// base is the base address of the arena chunk.</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>func userArenaHeapBitsSetType(typ *_type, ptr unsafe.Pointer, s *mspan) {
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>	base := s.base()
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>	h := s.writeUserArenaHeapBits(uintptr(ptr))
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>	p := typ.GCData <span class="comment">// start of 1-bit pointer mask (or GC program)</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>	var gcProgBits uintptr
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>	if typ.Kind_&amp;kindGCProg != 0 {
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>		<span class="comment">// Expand gc program, using the object itself for storage.</span>
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>		gcProgBits = runGCProg(addb(p, 4), (*byte)(ptr))
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>		p = (*byte)(ptr)
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>	}
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>	nb := typ.PtrBytes / goarch.PtrSize
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; nb; i += ptrBits {
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>		k := nb - i
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>		if k &gt; ptrBits {
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>			k = ptrBits
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>		}
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>		<span class="comment">// N.B. On big endian platforms we byte swap the data that we</span>
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>		<span class="comment">// read from GCData, which is always stored in little-endian order</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>		<span class="comment">// by the compiler. writeUserArenaHeapBits handles data in</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>		<span class="comment">// a platform-ordered way for efficiency, but stores back the</span>
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>		<span class="comment">// data in little endian order, since we expose the bitmap through</span>
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>		<span class="comment">// a dummy type.</span>
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>		h = h.write(s, readUintptr(addb(p, i/8)), k)
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>	}
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>	<span class="comment">// Note: we call pad here to ensure we emit explicit 0 bits</span>
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>	<span class="comment">// for the pointerless tail of the object. This ensures that</span>
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>	<span class="comment">// there&#39;s only a single noMorePtrs mark for the next object</span>
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>	<span class="comment">// to clear. We don&#39;t need to do this to clear stale noMorePtrs</span>
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>	<span class="comment">// markers from previous uses because arena chunk pointer bitmaps</span>
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>	<span class="comment">// are always fully cleared when reused.</span>
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>	h = h.pad(s, typ.Size_-typ.PtrBytes)
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>	h.flush(s, uintptr(ptr), typ.Size_)
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>	if typ.Kind_&amp;kindGCProg != 0 {
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>		<span class="comment">// Zero out temporary ptrmask buffer inside object.</span>
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>		memclrNoHeapPointers(ptr, (gcProgBits+7)/8)
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>	}
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>	<span class="comment">// Update the PtrBytes value in the type information. After this</span>
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>	<span class="comment">// point, the GC will observe the new bitmap.</span>
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>	s.largeType.PtrBytes = uintptr(ptr) - base + typ.PtrBytes
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>	<span class="comment">// Double-check that the bitmap was written out correctly.</span>
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>	const doubleCheck = false
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>	if doubleCheck {
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>		doubleCheckHeapPointersInterior(uintptr(ptr), uintptr(ptr), typ.Size_, typ.Size_, typ, &amp;s.largeType, s)
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>	}
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>}
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span><span class="comment">// For !goexperiment.AllocHeaders, to pass TestIntendedInlining.</span>
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>func writeHeapBitsForAddr() {
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>	panic(&#34;not implemented&#34;)
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>}
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span><span class="comment">// For !goexperiment.AllocHeaders.</span>
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>type heapBits struct {
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>}
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span><span class="comment">// For !goexperiment.AllocHeaders.</span>
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>func heapBitsForAddr(addr, size uintptr) heapBits {
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>	panic(&#34;not implemented&#34;)
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>}
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span><span class="comment">// For !goexperiment.AllocHeaders.</span>
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>func (h heapBits) next() (heapBits, uintptr) {
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>	panic(&#34;not implemented&#34;)
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>}
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span><span class="comment">// For !goexperiment.AllocHeaders.</span>
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>func (h heapBits) nextFast() (heapBits, uintptr) {
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>	panic(&#34;not implemented&#34;)
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>}
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>
</pre><p><a href="mbitmap_allocheaders.go?m=text">View as plain text</a></p>

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
