<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mpagealloc.go - Go Documentation Server</title>

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
<a href="mpagealloc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mpagealloc.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2019 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Page allocator.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// The page allocator manages mapped pages (defined by pageSize, NOT</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// physPageSize) for allocation and re-use. It is embedded into mheap.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// Pages are managed using a bitmap that is sharded into chunks.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// In the bitmap, 1 means in-use, and 0 means free. The bitmap spans the</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// process&#39;s address space. Chunks are managed in a sparse-array-style structure</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// similar to mheap.arenas, since the bitmap may be large on some systems.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// The bitmap is efficiently searched by using a radix tree in combination</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// with fast bit-wise intrinsics. Allocation is performed using an address-ordered</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// first-fit approach.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// Each entry in the radix tree is a summary that describes three properties of</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// a particular region of the address space: the number of contiguous free pages</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// at the start and end of the region it represents, and the maximum number of</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// contiguous free pages found anywhere in that region.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// Each level of the radix tree is stored as one contiguous array, which represents</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// a different granularity of subdivision of the processes&#39; address space. Thus, this</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// radix tree is actually implicit in these large arrays, as opposed to having explicit</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// dynamically-allocated pointer-based node structures. Naturally, these arrays may be</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// quite large for system with large address spaces, so in these cases they are mapped</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// into memory as needed. The leaf summaries of the tree correspond to a bitmap chunk.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// The root level (referred to as L0 and index 0 in pageAlloc.summary) has each</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// summary represent the largest section of address space (16 GiB on 64-bit systems),</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// with each subsequent level representing successively smaller subsections until we</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// reach the finest granularity at the leaves, a chunk.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// More specifically, each summary in each level (except for leaf summaries)</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// represents some number of entries in the following level. For example, each</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// summary in the root level may represent a 16 GiB region of address space,</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// and in the next level there could be 8 corresponding entries which represent 2</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// GiB subsections of that 16 GiB region, each of which could correspond to 8</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// entries in the next level which each represent 256 MiB regions, and so on.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// Thus, this design only scales to heaps so large, but can always be extended to</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// larger heaps by simply adding levels to the radix tree, which mostly costs</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// additional virtual address space. The choice of managing large arrays also means</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// that a large amount of virtual address space may be reserved by the runtime.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>package runtime
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>import (
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>const (
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// The size of a bitmap chunk, i.e. the amount of bits (that is, pages) to consider</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// in the bitmap at once.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	pallocChunkPages    = 1 &lt;&lt; logPallocChunkPages
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	pallocChunkBytes    = pallocChunkPages * pageSize
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	logPallocChunkPages = 9
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	logPallocChunkBytes = logPallocChunkPages + pageShift
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// The number of radix bits for each level.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// The value of 3 is chosen such that the block of summaries we need to scan at</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// each level fits in 64 bytes (2^3 summaries * 8 bytes per summary), which is</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// close to the L1 cache line width on many systems. Also, a value of 3 fits 4 tree</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// levels perfectly into the 21-bit pallocBits summary field at the root level.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// The following equation explains how each of the constants relate:</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// summaryL0Bits + (summaryLevels-1)*summaryLevelBits + logPallocChunkBytes = heapAddrBits</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// summaryLevels is an architecture-dependent value defined in mpagealloc_*.go.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	summaryLevelBits = 3
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	summaryL0Bits    = heapAddrBits - logPallocChunkBytes - (summaryLevels-1)*summaryLevelBits
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// pallocChunksL2Bits is the number of bits of the chunk index number</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// covered by the second level of the chunks map.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// See (*pageAlloc).chunks for more details. Update the documentation</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// there should this change.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	pallocChunksL2Bits  = heapAddrBits - logPallocChunkBytes - pallocChunksL1Bits
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	pallocChunksL1Shift = pallocChunksL2Bits
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// maxSearchAddr returns the maximum searchAddr value, which indicates</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// that the heap has no free space.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// This function exists just to make it clear that this is the maximum address</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// for the page allocator&#39;s search space. See maxOffAddr for details.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// It&#39;s a function (rather than a variable) because it needs to be</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// usable before package runtime&#39;s dynamic initialization is complete.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// See #51913 for details.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>func maxSearchAddr() offAddr { return maxOffAddr }
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// Global chunk index.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// Represents an index into the leaf level of the radix tree.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// Similar to arenaIndex, except instead of arenas, it divides the address</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// space into chunks.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>type chunkIdx uint
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// chunkIndex returns the global index of the palloc chunk containing the</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// pointer p.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>func chunkIndex(p uintptr) chunkIdx {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	return chunkIdx((p - arenaBaseOffset) / pallocChunkBytes)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// chunkBase returns the base address of the palloc chunk at index ci.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>func chunkBase(ci chunkIdx) uintptr {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	return uintptr(ci)*pallocChunkBytes + arenaBaseOffset
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">// chunkPageIndex computes the index of the page that contains p,</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// relative to the chunk which contains p.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>func chunkPageIndex(p uintptr) uint {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	return uint(p % pallocChunkBytes / pageSize)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// l1 returns the index into the first level of (*pageAlloc).chunks.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>func (i chunkIdx) l1() uint {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if pallocChunksL1Bits == 0 {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		<span class="comment">// Let the compiler optimize this away if there&#39;s no</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		<span class="comment">// L1 map.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		return 0
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	} else {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		return uint(i) &gt;&gt; pallocChunksL1Shift
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// l2 returns the index into the second level of (*pageAlloc).chunks.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>func (i chunkIdx) l2() uint {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	if pallocChunksL1Bits == 0 {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		return uint(i)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	} else {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		return uint(i) &amp; (1&lt;&lt;pallocChunksL2Bits - 1)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// offAddrToLevelIndex converts an address in the offset address space</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// to the index into summary[level] containing addr.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>func offAddrToLevelIndex(level int, addr offAddr) int {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	return int((addr.a - arenaBaseOffset) &gt;&gt; levelShift[level])
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// levelIndexToOffAddr converts an index into summary[level] into</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// the corresponding address in the offset address space.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>func levelIndexToOffAddr(level, idx int) offAddr {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	return offAddr{(uintptr(idx) &lt;&lt; levelShift[level]) + arenaBaseOffset}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">// addrsToSummaryRange converts base and limit pointers into a range</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// of entries for the given summary level.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// The returned range is inclusive on the lower bound and exclusive on</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// the upper bound.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>func addrsToSummaryRange(level int, base, limit uintptr) (lo int, hi int) {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// This is slightly more nuanced than just a shift for the exclusive</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// upper-bound. Note that the exclusive upper bound may be within a</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// summary at this level, meaning if we just do the obvious computation</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// hi will end up being an inclusive upper bound. Unfortunately, just</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// adding 1 to that is too broad since we might be on the very edge</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// of a summary&#39;s max page count boundary for this level</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// (1 &lt;&lt; levelLogPages[level]). So, make limit an inclusive upper bound</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// then shift, then add 1, so we get an exclusive upper bound at the end.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	lo = int((base - arenaBaseOffset) &gt;&gt; levelShift[level])
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	hi = int(((limit-1)-arenaBaseOffset)&gt;&gt;levelShift[level]) + 1
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	return
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// blockAlignSummaryRange aligns indices into the given level to that</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// level&#39;s block width (1 &lt;&lt; levelBits[level]). It assumes lo is inclusive</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">// and hi is exclusive, and so aligns them down and up respectively.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>func blockAlignSummaryRange(level int, lo, hi int) (int, int) {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	e := uintptr(1) &lt;&lt; levelBits[level]
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	return int(alignDown(uintptr(lo), e)), int(alignUp(uintptr(hi), e))
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>type pageAlloc struct {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// Radix tree of summaries.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// Each slice&#39;s cap represents the whole memory reservation.</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// Each slice&#39;s len reflects the allocator&#39;s maximum known</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// mapped heap address for that level.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// The backing store of each summary level is reserved in init</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">// and may or may not be committed in grow (small address spaces</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// may commit all the memory in init).</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// The purpose of keeping len &lt;= cap is to enforce bounds checks</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// on the top end of the slice so that instead of an unknown</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// runtime segmentation fault, we get a much friendlier out-of-bounds</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// error.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// To iterate over a summary level, use inUse to determine which ranges</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// are currently available. Otherwise one might try to access</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// memory which is only Reserved which may result in a hard fault.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// We may still get segmentation faults &lt; len since some of that</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// memory may not be committed yet.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	summary [summaryLevels][]pallocSum
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// chunks is a slice of bitmap chunks.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// The total size of chunks is quite large on most 64-bit platforms</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// (O(GiB) or more) if flattened, so rather than making one large mapping</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// (which has problems on some platforms, even when PROT_NONE) we use a</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">// two-level sparse array approach similar to the arena index in mheap.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	<span class="comment">// To find the chunk containing a memory address `a`, do:</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">//   chunkOf(chunkIndex(a))</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// Below is a table describing the configuration for chunks for various</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// heapAddrBits supported by the runtime.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	<span class="comment">// heapAddrBits | L1 Bits | L2 Bits | L2 Entry Size</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	<span class="comment">// ------------------------------------------------</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	<span class="comment">// 32           | 0       | 10      | 128 KiB</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// 33 (iOS)     | 0       | 11      | 256 KiB</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// 48           | 13      | 13      | 1 MiB</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	<span class="comment">// There&#39;s no reason to use the L1 part of chunks on 32-bit, the</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">// address space is small so the L2 is small. For platforms with a</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	<span class="comment">// 48-bit address space, we pick the L1 such that the L2 is 1 MiB</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// in size, which is a good balance between low granularity without</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	<span class="comment">// making the impact on BSS too high (note the L1 is stored directly</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// in pageAlloc).</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// To iterate over the bitmap, use inUse to determine which ranges</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">// are currently available. Otherwise one might iterate over unused</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">// ranges.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	<span class="comment">// Protected by mheapLock.</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Consider changing the definition of the bitmap</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// such that 1 means free and 0 means in-use so that summaries and</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// the bitmaps align better on zero-values.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	chunks [1 &lt;&lt; pallocChunksL1Bits]*[1 &lt;&lt; pallocChunksL2Bits]pallocData
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// The address to start an allocation search with. It must never</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	<span class="comment">// point to any memory that is not contained in inUse, i.e.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// inUse.contains(searchAddr.addr()) must always be true. The one</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// exception to this rule is that it may take on the value of</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// maxOffAddr to indicate that the heap is exhausted.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	<span class="comment">// We guarantee that all valid heap addresses below this value</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// are allocated and not worth searching.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	searchAddr offAddr
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// start and end represent the chunk indices</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// which pageAlloc knows about. It assumes</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// chunks in the range [start, end) are</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	<span class="comment">// currently ready to use.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	start, end chunkIdx
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">// inUse is a slice of ranges of address space which are</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	<span class="comment">// known by the page allocator to be currently in-use (passed</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">// to grow).</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// We care much more about having a contiguous heap in these cases</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// and take additional measures to ensure that, so in nearly all</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// cases this should have just 1 element.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	<span class="comment">// All access is protected by the mheapLock.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	inUse addrRanges
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// scav stores the scavenger state.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	scav struct {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		<span class="comment">// index is an efficient index of chunks that have pages available to</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		<span class="comment">// scavenge.</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		index scavengeIndex
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		<span class="comment">// releasedBg is the amount of memory released in the background this</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		<span class="comment">// scavenge cycle.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		releasedBg atomic.Uintptr
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		<span class="comment">// releasedEager is the amount of memory released eagerly this scavenge</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		<span class="comment">// cycle.</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		releasedEager atomic.Uintptr
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	<span class="comment">// mheap_.lock. This level of indirection makes it possible</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	<span class="comment">// to test pageAlloc independently of the runtime allocator.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	mheapLock *mutex
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	<span class="comment">// sysStat is the runtime memstat to update when new system</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	<span class="comment">// memory is committed by the pageAlloc for allocation metadata.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	sysStat *sysMemStat
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// summaryMappedReady is the number of bytes mapped in the Ready state</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// in the summary structure. Used only for testing currently.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">// Protected by mheapLock.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	summaryMappedReady uintptr
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// chunkHugePages indicates whether page bitmap chunks should be backed</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	<span class="comment">// by huge pages.</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	chunkHugePages bool
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">// Whether or not this struct is being used in tests.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	test bool
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>func (p *pageAlloc) init(mheapLock *mutex, sysStat *sysMemStat, test bool) {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	if levelLogPages[0] &gt; logMaxPackedValue {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		<span class="comment">// We can&#39;t represent 1&lt;&lt;levelLogPages[0] pages, the maximum number</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		<span class="comment">// of pages we need to represent at the root level, in a summary, which</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		<span class="comment">// is a big problem. Throw.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		print(&#34;runtime: root level max pages = &#34;, 1&lt;&lt;levelLogPages[0], &#34;\n&#34;)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		print(&#34;runtime: summary max pages = &#34;, maxPackedValue, &#34;\n&#34;)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		throw(&#34;root level max pages doesn&#39;t fit in summary&#34;)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	p.sysStat = sysStat
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// Initialize p.inUse.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	p.inUse.init(sysStat)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	<span class="comment">// System-dependent initialization.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	p.sysInit(test)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	<span class="comment">// Start with the searchAddr in a state indicating there&#39;s no free memory.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	p.searchAddr = maxSearchAddr()
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// Set the mheapLock.</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	p.mheapLock = mheapLock
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	<span class="comment">// Initialize the scavenge index.</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	p.summaryMappedReady += p.scav.index.init(test, sysStat)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	<span class="comment">// Set if we&#39;re in a test.</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	p.test = test
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// tryChunkOf returns the bitmap data for the given chunk.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// Returns nil if the chunk data has not been mapped.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>func (p *pageAlloc) tryChunkOf(ci chunkIdx) *pallocData {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	l2 := p.chunks[ci.l1()]
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	if l2 == nil {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		return nil
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	return &amp;l2[ci.l2()]
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// chunkOf returns the chunk at the given chunk index.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// The chunk index must be valid or this method may throw.</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>func (p *pageAlloc) chunkOf(ci chunkIdx) *pallocData {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	return &amp;p.chunks[ci.l1()][ci.l2()]
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// grow sets up the metadata for the address range [base, base+size).</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// It may allocate metadata, in which case *p.sysStat will be updated.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// p.mheapLock must be held.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>func (p *pageAlloc) grow(base, size uintptr) {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	assertLockHeld(p.mheapLock)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	<span class="comment">// Round up to chunks, since we can&#39;t deal with increments smaller</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	<span class="comment">// than chunks. Also, sysGrow expects aligned values.</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	limit := alignUp(base+size, pallocChunkBytes)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	base = alignDown(base, pallocChunkBytes)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	<span class="comment">// Grow the summary levels in a system-dependent manner.</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	<span class="comment">// We just update a bunch of additional metadata here.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	p.sysGrow(base, limit)
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// Grow the scavenge index.</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	p.summaryMappedReady += p.scav.index.grow(base, limit, p.sysStat)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	<span class="comment">// Update p.start and p.end.</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	<span class="comment">// If no growth happened yet, start == 0. This is generally</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	<span class="comment">// safe since the zero page is unmapped.</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	firstGrowth := p.start == 0
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	start, end := chunkIndex(base), chunkIndex(limit)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	if firstGrowth || start &lt; p.start {
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		p.start = start
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	if end &gt; p.end {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		p.end = end
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	<span class="comment">// Note that [base, limit) will never overlap with any existing</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// range inUse because grow only ever adds never-used memory</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	<span class="comment">// regions to the page allocator.</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	p.inUse.add(makeAddrRange(base, limit))
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	<span class="comment">// A grow operation is a lot like a free operation, so if our</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	<span class="comment">// chunk ends up below p.searchAddr, update p.searchAddr to the</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">// new address, just like in free.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	if b := (offAddr{base}); b.lessThan(p.searchAddr) {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		p.searchAddr = b
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	<span class="comment">// Add entries into chunks, which is sparse, if needed. Then,</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	<span class="comment">// initialize the bitmap.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	<span class="comment">// Newly-grown memory is always considered scavenged.</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	<span class="comment">// Set all the bits in the scavenged bitmaps high.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	for c := chunkIndex(base); c &lt; chunkIndex(limit); c++ {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		if p.chunks[c.l1()] == nil {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			<span class="comment">// Create the necessary l2 entry.</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			const l2Size = unsafe.Sizeof(*p.chunks[0])
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			r := sysAlloc(l2Size, p.sysStat)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			if r == nil {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>				throw(&#34;pageAlloc: out of memory&#34;)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			if !p.test {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>				<span class="comment">// Make the chunk mapping eligible or ineligible</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>				<span class="comment">// for huge pages, depending on what our current</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>				<span class="comment">// state is.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>				if p.chunkHugePages {
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>					sysHugePage(r, l2Size)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>				} else {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>					sysNoHugePage(r, l2Size)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>				}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			<span class="comment">// Store the new chunk block but avoid a write barrier.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			<span class="comment">// grow is used in call chains that disallow write barriers.</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			*(*uintptr)(unsafe.Pointer(&amp;p.chunks[c.l1()])) = uintptr(r)
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		p.chunkOf(c).scavenged.setRange(0, pallocChunkPages)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	<span class="comment">// Update summaries accordingly. The grow acts like a free, so</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	<span class="comment">// we need to ensure this newly-free memory is visible in the</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	<span class="comment">// summaries.</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	p.update(base, size/pageSize, true, false)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span><span class="comment">// enableChunkHugePages enables huge pages for the chunk bitmap mappings (disabled by default).</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span><span class="comment">// This function is idempotent.</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span><span class="comment">// A note on latency: for sufficiently small heaps (&lt;10s of GiB) this function will take constant</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span><span class="comment">// time, but may take time proportional to the size of the mapped heap beyond that.</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span><span class="comment">// The heap lock must not be held over this operation, since it will briefly acquire</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span><span class="comment">// the heap lock.</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span><span class="comment">// Must be called on the system stack because it acquires the heap lock.</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>func (p *pageAlloc) enableChunkHugePages() {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	<span class="comment">// Grab the heap lock to turn on huge pages for new chunks and clone the current</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	<span class="comment">// heap address space ranges.</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	<span class="comment">// After the lock is released, we can be sure that bitmaps for any new chunks may</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	<span class="comment">// be backed with huge pages, and we have the address space for the rest of the</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	<span class="comment">// chunks. At the end of this function, all chunk metadata should be backed by huge</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	<span class="comment">// pages.</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	lock(&amp;mheap_.lock)
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	if p.chunkHugePages {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		unlock(&amp;mheap_.lock)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		return
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	p.chunkHugePages = true
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	var inUse addrRanges
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	inUse.sysStat = p.sysStat
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	p.inUse.cloneInto(&amp;inUse)
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	unlock(&amp;mheap_.lock)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	<span class="comment">// This might seem like a lot of work, but all these loops are for generality.</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	<span class="comment">// For a 1 GiB contiguous heap, a 48-bit address space, 13 L1 bits, a palloc chunk size</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	<span class="comment">// of 4 MiB, and adherence to the default set of heap address hints, this will result in</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	<span class="comment">// exactly 1 call to sysHugePage.</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	for _, r := range p.inUse.ranges {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		for i := chunkIndex(r.base.addr()).l1(); i &lt; chunkIndex(r.limit.addr()-1).l1(); i++ {
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>			<span class="comment">// N.B. We can assume that p.chunks[i] is non-nil and in a mapped part of p.chunks</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			<span class="comment">// because it&#39;s derived from inUse, which never shrinks.</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			sysHugePage(unsafe.Pointer(p.chunks[i]), unsafe.Sizeof(*p.chunks[0]))
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span><span class="comment">// update updates heap metadata. It must be called each time the bitmap</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span><span class="comment">// is updated.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span><span class="comment">// If contig is true, update does some optimizations assuming that there was</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span><span class="comment">// a contiguous allocation or free between addr and addr+npages. alloc indicates</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span><span class="comment">// whether the operation performed was an allocation or a free.</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span><span class="comment">// p.mheapLock must be held.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>func (p *pageAlloc) update(base, npages uintptr, contig, alloc bool) {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	assertLockHeld(p.mheapLock)
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	<span class="comment">// base, limit, start, and end are inclusive.</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	limit := base + npages*pageSize - 1
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	sc, ec := chunkIndex(base), chunkIndex(limit)
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	<span class="comment">// Handle updating the lowest level first.</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	if sc == ec {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		<span class="comment">// Fast path: the allocation doesn&#39;t span more than one chunk,</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		<span class="comment">// so update this one and if the summary didn&#39;t change, return.</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		x := p.summary[len(p.summary)-1][sc]
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		y := p.chunkOf(sc).summarize()
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		if x == y {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			return
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		}
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		p.summary[len(p.summary)-1][sc] = y
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	} else if contig {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		<span class="comment">// Slow contiguous path: the allocation spans more than one chunk</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		<span class="comment">// and at least one summary is guaranteed to change.</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		summary := p.summary[len(p.summary)-1]
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		<span class="comment">// Update the summary for chunk sc.</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		summary[sc] = p.chunkOf(sc).summarize()
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		<span class="comment">// Update the summaries for chunks in between, which are</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		<span class="comment">// either totally allocated or freed.</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		whole := p.summary[len(p.summary)-1][sc+1 : ec]
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		if alloc {
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>			<span class="comment">// Should optimize into a memclr.</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>			for i := range whole {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>				whole[i] = 0
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		} else {
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>			for i := range whole {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>				whole[i] = freeChunkSum
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		<span class="comment">// Update the summary for chunk ec.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		summary[ec] = p.chunkOf(ec).summarize()
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	} else {
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		<span class="comment">// Slow general path: the allocation spans more than one chunk</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		<span class="comment">// and at least one summary is guaranteed to change.</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		<span class="comment">// We can&#39;t assume a contiguous allocation happened, so walk over</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		<span class="comment">// every chunk in the range and manually recompute the summary.</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		summary := p.summary[len(p.summary)-1]
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		for c := sc; c &lt;= ec; c++ {
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>			summary[c] = p.chunkOf(c).summarize()
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	<span class="comment">// Walk up the radix tree and update the summaries appropriately.</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	changed := true
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	for l := len(p.summary) - 2; l &gt;= 0 &amp;&amp; changed; l-- {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		<span class="comment">// Update summaries at level l from summaries at level l+1.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		changed = false
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		<span class="comment">// &#34;Constants&#34; for the previous level which we</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		<span class="comment">// need to compute the summary from that level.</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		logEntriesPerBlock := levelBits[l+1]
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		logMaxPages := levelLogPages[l+1]
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		<span class="comment">// lo and hi describe all the parts of the level we need to look at.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		lo, hi := addrsToSummaryRange(l, base, limit+1)
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		<span class="comment">// Iterate over each block, updating the corresponding summary in the less-granular level.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		for i := lo; i &lt; hi; i++ {
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>			children := p.summary[l+1][i&lt;&lt;logEntriesPerBlock : (i+1)&lt;&lt;logEntriesPerBlock]
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			sum := mergeSummaries(children, logMaxPages)
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>			old := p.summary[l][i]
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>			if old != sum {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>				changed = true
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>				p.summary[l][i] = sum
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>}
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span><span class="comment">// allocRange marks the range of memory [base, base+npages*pageSize) as</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span><span class="comment">// allocated. It also updates the summaries to reflect the newly-updated</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span><span class="comment">// bitmap.</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span><span class="comment">// Returns the amount of scavenged memory in bytes present in the</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span><span class="comment">// allocated range.</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span><span class="comment">// p.mheapLock must be held.</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>func (p *pageAlloc) allocRange(base, npages uintptr) uintptr {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	assertLockHeld(p.mheapLock)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	limit := base + npages*pageSize - 1
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	sc, ec := chunkIndex(base), chunkIndex(limit)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	si, ei := chunkPageIndex(base), chunkPageIndex(limit)
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	scav := uint(0)
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	if sc == ec {
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		<span class="comment">// The range doesn&#39;t cross any chunk boundaries.</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		chunk := p.chunkOf(sc)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		scav += chunk.scavenged.popcntRange(si, ei+1-si)
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		chunk.allocRange(si, ei+1-si)
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		p.scav.index.alloc(sc, ei+1-si)
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	} else {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		<span class="comment">// The range crosses at least one chunk boundary.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		chunk := p.chunkOf(sc)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		scav += chunk.scavenged.popcntRange(si, pallocChunkPages-si)
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		chunk.allocRange(si, pallocChunkPages-si)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		p.scav.index.alloc(sc, pallocChunkPages-si)
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		for c := sc + 1; c &lt; ec; c++ {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			chunk := p.chunkOf(c)
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>			scav += chunk.scavenged.popcntRange(0, pallocChunkPages)
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>			chunk.allocAll()
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>			p.scav.index.alloc(c, pallocChunkPages)
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		chunk = p.chunkOf(ec)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		scav += chunk.scavenged.popcntRange(0, ei+1)
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		chunk.allocRange(0, ei+1)
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		p.scav.index.alloc(ec, ei+1)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	}
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	p.update(base, npages, true, true)
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	return uintptr(scav) * pageSize
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span><span class="comment">// findMappedAddr returns the smallest mapped offAddr that is</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span><span class="comment">// &gt;= addr. That is, if addr refers to mapped memory, then it is</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span><span class="comment">// returned. If addr is higher than any mapped region, then</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span><span class="comment">// it returns maxOffAddr.</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span><span class="comment">// p.mheapLock must be held.</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>func (p *pageAlloc) findMappedAddr(addr offAddr) offAddr {
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	assertLockHeld(p.mheapLock)
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;re not in a test, validate first by checking mheap_.arenas.</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	<span class="comment">// This is a fast path which is only safe to use outside of testing.</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	ai := arenaIndex(addr.addr())
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	if p.test || mheap_.arenas[ai.l1()] == nil || mheap_.arenas[ai.l1()][ai.l2()] == nil {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		vAddr, ok := p.inUse.findAddrGreaterEqual(addr.addr())
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		if ok {
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>			return offAddr{vAddr}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		} else {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>			<span class="comment">// The candidate search address is greater than any</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>			<span class="comment">// known address, which means we definitely have no</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			<span class="comment">// free memory left.</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			return maxOffAddr
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		}
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	return addr
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span><span class="comment">// find searches for the first (address-ordered) contiguous free region of</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span><span class="comment">// npages in size and returns a base address for that region.</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span><span class="comment">// It uses p.searchAddr to prune its search and assumes that no palloc chunks</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span><span class="comment">// below chunkIndex(p.searchAddr) contain any free memory at all.</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span><span class="comment">// find also computes and returns a candidate p.searchAddr, which may or</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span><span class="comment">// may not prune more of the address space than p.searchAddr already does.</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span><span class="comment">// This candidate is always a valid p.searchAddr.</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span><span class="comment">// find represents the slow path and the full radix tree search.</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span><span class="comment">// Returns a base address of 0 on failure, in which case the candidate</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span><span class="comment">// searchAddr returned is invalid and must be ignored.</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span><span class="comment">// p.mheapLock must be held.</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>func (p *pageAlloc) find(npages uintptr) (uintptr, offAddr) {
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	assertLockHeld(p.mheapLock)
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	<span class="comment">// Search algorithm.</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	<span class="comment">// This algorithm walks each level l of the radix tree from the root level</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	<span class="comment">// to the leaf level. It iterates over at most 1 &lt;&lt; levelBits[l] of entries</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	<span class="comment">// in a given level in the radix tree, and uses the summary information to</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	<span class="comment">// find either:</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	<span class="comment">//  1) That a given subtree contains a large enough contiguous region, at</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	<span class="comment">//     which point it continues iterating on the next level, or</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	<span class="comment">//  2) That there are enough contiguous boundary-crossing bits to satisfy</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	<span class="comment">//     the allocation, at which point it knows exactly where to start</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	<span class="comment">//     allocating from.</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	<span class="comment">// i tracks the index into the current level l&#39;s structure for the</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	<span class="comment">// contiguous 1 &lt;&lt; levelBits[l] entries we&#39;re actually interested in.</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	<span class="comment">// NOTE: Technically this search could allocate a region which crosses</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	<span class="comment">// the arenaBaseOffset boundary, which when arenaBaseOffset != 0, is</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	<span class="comment">// a discontinuity. However, the only way this could happen is if the</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	<span class="comment">// page at the zero address is mapped, and this is impossible on</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	<span class="comment">// every system we support where arenaBaseOffset != 0. So, the</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	<span class="comment">// discontinuity is already encoded in the fact that the OS will never</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	<span class="comment">// map the zero page for us, and this function doesn&#39;t try to handle</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	<span class="comment">// this case in any way.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	<span class="comment">// i is the beginning of the block of entries we&#39;re searching at the</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	<span class="comment">// current level.</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	i := 0
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	<span class="comment">// firstFree is the region of address space that we are certain to</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	<span class="comment">// find the first free page in the heap. base and bound are the inclusive</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	<span class="comment">// bounds of this window, and both are addresses in the linearized, contiguous</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	<span class="comment">// view of the address space (with arenaBaseOffset pre-added). At each level,</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	<span class="comment">// this window is narrowed as we find the memory region containing the</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	<span class="comment">// first free page of memory. To begin with, the range reflects the</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	<span class="comment">// full process address space.</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	<span class="comment">// firstFree is updated by calling foundFree each time free space in the</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	<span class="comment">// heap is discovered.</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	<span class="comment">// At the end of the search, base.addr() is the best new</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	<span class="comment">// searchAddr we could deduce in this search.</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	firstFree := struct {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>		base, bound offAddr
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	}{
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		base:  minOffAddr,
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		bound: maxOffAddr,
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	}
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	<span class="comment">// foundFree takes the given address range [addr, addr+size) and</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	<span class="comment">// updates firstFree if it is a narrower range. The input range must</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	<span class="comment">// either be fully contained within firstFree or not overlap with it</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	<span class="comment">// at all.</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	<span class="comment">// This way, we&#39;ll record the first summary we find with any free</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	<span class="comment">// pages on the root level and narrow that down if we descend into</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	<span class="comment">// that summary. But as soon as we need to iterate beyond that summary</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	<span class="comment">// in a level to find a large enough range, we&#39;ll stop narrowing.</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	foundFree := func(addr offAddr, size uintptr) {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>		if firstFree.base.lessEqual(addr) &amp;&amp; addr.add(size-1).lessEqual(firstFree.bound) {
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			<span class="comment">// This range fits within the current firstFree window, so narrow</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>			<span class="comment">// down the firstFree window to the base and bound of this range.</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>			firstFree.base = addr
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>			firstFree.bound = addr.add(size - 1)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		} else if !(addr.add(size-1).lessThan(firstFree.base) || firstFree.bound.lessThan(addr)) {
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>			<span class="comment">// This range only partially overlaps with the firstFree range,</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>			<span class="comment">// so throw.</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>			print(&#34;runtime: addr = &#34;, hex(addr.addr()), &#34;, size = &#34;, size, &#34;\n&#34;)
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>			print(&#34;runtime: base = &#34;, hex(firstFree.base.addr()), &#34;, bound = &#34;, hex(firstFree.bound.addr()), &#34;\n&#34;)
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>			throw(&#34;range partially overlaps&#34;)
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		}
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	}
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	<span class="comment">// lastSum is the summary which we saw on the previous level that made us</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	<span class="comment">// move on to the next level. Used to print additional information in the</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	<span class="comment">// case of a catastrophic failure.</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>	<span class="comment">// lastSumIdx is that summary&#39;s index in the previous level.</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	lastSum := packPallocSum(0, 0, 0)
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	lastSumIdx := -1
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>nextLevel:
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	for l := 0; l &lt; len(p.summary); l++ {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>		<span class="comment">// For the root level, entriesPerBlock is the whole level.</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		entriesPerBlock := 1 &lt;&lt; levelBits[l]
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>		logMaxPages := levelLogPages[l]
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>		<span class="comment">// We&#39;ve moved into a new level, so let&#39;s update i to our new</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		<span class="comment">// starting index. This is a no-op for level 0.</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		i &lt;&lt;= levelBits[l]
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		<span class="comment">// Slice out the block of entries we care about.</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		entries := p.summary[l][i : i+entriesPerBlock]
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		<span class="comment">// Determine j0, the first index we should start iterating from.</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		<span class="comment">// The searchAddr may help us eliminate iterations if we followed the</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>		<span class="comment">// searchAddr on the previous level or we&#39;re on the root level, in which</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		<span class="comment">// case the searchAddr should be the same as i after levelShift.</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		j0 := 0
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		if searchIdx := offAddrToLevelIndex(l, p.searchAddr); searchIdx&amp;^(entriesPerBlock-1) == i {
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>			j0 = searchIdx &amp; (entriesPerBlock - 1)
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		}
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>		<span class="comment">// Run over the level entries looking for</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		<span class="comment">// a contiguous run of at least npages either</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>		<span class="comment">// within an entry or across entries.</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>		<span class="comment">// base contains the page index (relative to</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		<span class="comment">// the first entry&#39;s first page) of the currently</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		<span class="comment">// considered run of consecutive pages.</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		<span class="comment">// size contains the size of the currently considered</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		<span class="comment">// run of consecutive pages.</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		var base, size uint
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>		for j := j0; j &lt; len(entries); j++ {
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>			sum := entries[j]
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>			if sum == 0 {
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>				<span class="comment">// A full entry means we broke any streak and</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>				<span class="comment">// that we should skip it altogether.</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>				size = 0
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>				continue
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			}
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>			<span class="comment">// We&#39;ve encountered a non-zero summary which means</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>			<span class="comment">// free memory, so update firstFree.</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>			foundFree(levelIndexToOffAddr(l, i+j), (uintptr(1)&lt;&lt;logMaxPages)*pageSize)
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>			s := sum.start()
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>			if size+s &gt;= uint(npages) {
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>				<span class="comment">// If size == 0 we don&#39;t have a run yet,</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>				<span class="comment">// which means base isn&#39;t valid. So, set</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>				<span class="comment">// base to the first page in this block.</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>				if size == 0 {
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>					base = uint(j) &lt;&lt; logMaxPages
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>				}
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>				<span class="comment">// We hit npages; we&#39;re done!</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>				size += s
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>				break
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>			}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>			if sum.max() &gt;= uint(npages) {
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>				<span class="comment">// The entry itself contains npages contiguous</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>				<span class="comment">// free pages, so continue on the next level</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>				<span class="comment">// to find that run.</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>				i += j
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>				lastSumIdx = i
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>				lastSum = sum
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>				continue nextLevel
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>			}
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>			if size == 0 || s &lt; 1&lt;&lt;logMaxPages {
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>				<span class="comment">// We either don&#39;t have a current run started, or this entry</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>				<span class="comment">// isn&#39;t totally free (meaning we can&#39;t continue the current</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>				<span class="comment">// one), so try to begin a new run by setting size and base</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>				<span class="comment">// based on sum.end.</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>				size = sum.end()
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>				base = uint(j+1)&lt;&lt;logMaxPages - size
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>				continue
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>			}
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>			<span class="comment">// The entry is completely free, so continue the run.</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>			size += 1 &lt;&lt; logMaxPages
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>		}
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>		if size &gt;= uint(npages) {
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>			<span class="comment">// We found a sufficiently large run of free pages straddling</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>			<span class="comment">// some boundary, so compute the address and return it.</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>			addr := levelIndexToOffAddr(l, i).add(uintptr(base) * pageSize).addr()
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>			return addr, p.findMappedAddr(firstFree.base)
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>		}
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		if l == 0 {
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>			<span class="comment">// We&#39;re at level zero, so that means we&#39;ve exhausted our search.</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>			return 0, maxSearchAddr()
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		}
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re not at level zero, and we exhausted the level we were looking in.</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>		<span class="comment">// This means that either our calculations were wrong or the level above</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		<span class="comment">// lied to us. In either case, dump some useful state and throw.</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		print(&#34;runtime: summary[&#34;, l-1, &#34;][&#34;, lastSumIdx, &#34;] = &#34;, lastSum.start(), &#34;, &#34;, lastSum.max(), &#34;, &#34;, lastSum.end(), &#34;\n&#34;)
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>		print(&#34;runtime: level = &#34;, l, &#34;, npages = &#34;, npages, &#34;, j0 = &#34;, j0, &#34;\n&#34;)
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		print(&#34;runtime: p.searchAddr = &#34;, hex(p.searchAddr.addr()), &#34;, i = &#34;, i, &#34;\n&#34;)
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		print(&#34;runtime: levelShift[level] = &#34;, levelShift[l], &#34;, levelBits[level] = &#34;, levelBits[l], &#34;\n&#34;)
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>		for j := 0; j &lt; len(entries); j++ {
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>			sum := entries[j]
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>			print(&#34;runtime: summary[&#34;, l, &#34;][&#34;, i+j, &#34;] = (&#34;, sum.start(), &#34;, &#34;, sum.max(), &#34;, &#34;, sum.end(), &#34;)\n&#34;)
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		}
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		throw(&#34;bad summary data&#34;)
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	}
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	<span class="comment">// Since we&#39;ve gotten to this point, that means we haven&#39;t found a</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>	<span class="comment">// sufficiently-sized free region straddling some boundary (chunk or larger).</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	<span class="comment">// This means the last summary we inspected must have had a large enough &#34;max&#34;</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	<span class="comment">// value, so look inside the chunk to find a suitable run.</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	<span class="comment">// After iterating over all levels, i must contain a chunk index which</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	<span class="comment">// is what the final level represents.</span>
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	ci := chunkIdx(i)
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	j, searchIdx := p.chunkOf(ci).find(npages, 0)
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>	if j == ^uint(0) {
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>		<span class="comment">// We couldn&#39;t find any space in this chunk despite the summaries telling</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		<span class="comment">// us it should be there. There&#39;s likely a bug, so dump some state and throw.</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		sum := p.summary[len(p.summary)-1][i]
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>		print(&#34;runtime: summary[&#34;, len(p.summary)-1, &#34;][&#34;, i, &#34;] = (&#34;, sum.start(), &#34;, &#34;, sum.max(), &#34;, &#34;, sum.end(), &#34;)\n&#34;)
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		print(&#34;runtime: npages = &#34;, npages, &#34;\n&#34;)
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		throw(&#34;bad summary data&#34;)
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	}
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	<span class="comment">// Compute the address at which the free space starts.</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	addr := chunkBase(ci) + uintptr(j)*pageSize
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	<span class="comment">// Since we actually searched the chunk, we may have</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	<span class="comment">// found an even narrower free window.</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	searchAddr := chunkBase(ci) + uintptr(searchIdx)*pageSize
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	foundFree(offAddr{searchAddr}, chunkBase(ci+1)-searchAddr)
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	return addr, p.findMappedAddr(firstFree.base)
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>}
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span><span class="comment">// alloc allocates npages worth of memory from the page heap, returning the base</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span><span class="comment">// address for the allocation and the amount of scavenged memory in bytes</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span><span class="comment">// contained in the region [base address, base address + npages*pageSize).</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span><span class="comment">// Returns a 0 base address on failure, in which case other returned values</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span><span class="comment">// should be ignored.</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span><span class="comment">// p.mheapLock must be held.</span>
<span id="L871" class="ln">   871&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span><span class="comment">// Must run on the system stack because p.mheapLock must be held.</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>func (p *pageAlloc) alloc(npages uintptr) (addr uintptr, scav uintptr) {
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	assertLockHeld(p.mheapLock)
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	<span class="comment">// If the searchAddr refers to a region which has a higher address than</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	<span class="comment">// any known chunk, then we know we&#39;re out of memory.</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	if chunkIndex(p.searchAddr.addr()) &gt;= p.end {
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>		return 0, 0
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	}
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	<span class="comment">// If npages has a chance of fitting in the chunk where the searchAddr is,</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	<span class="comment">// search it directly.</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	searchAddr := minOffAddr
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	if pallocChunkPages-chunkPageIndex(p.searchAddr.addr()) &gt;= uint(npages) {
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>		<span class="comment">// npages is guaranteed to be no greater than pallocChunkPages here.</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>		i := chunkIndex(p.searchAddr.addr())
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>		if max := p.summary[len(p.summary)-1][i].max(); max &gt;= uint(npages) {
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			j, searchIdx := p.chunkOf(i).find(npages, chunkPageIndex(p.searchAddr.addr()))
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>			if j == ^uint(0) {
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>				print(&#34;runtime: max = &#34;, max, &#34;, npages = &#34;, npages, &#34;\n&#34;)
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>				print(&#34;runtime: searchIdx = &#34;, chunkPageIndex(p.searchAddr.addr()), &#34;, p.searchAddr = &#34;, hex(p.searchAddr.addr()), &#34;\n&#34;)
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>				throw(&#34;bad summary data&#34;)
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>			}
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>			addr = chunkBase(i) + uintptr(j)*pageSize
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>			searchAddr = offAddr{chunkBase(i) + uintptr(searchIdx)*pageSize}
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>			goto Found
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		}
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	}
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	<span class="comment">// We failed to use a searchAddr for one reason or another, so try</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>	<span class="comment">// the slow path.</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	addr, searchAddr = p.find(npages)
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	if addr == 0 {
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>		if npages == 1 {
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>			<span class="comment">// We failed to find a single free page, the smallest unit</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>			<span class="comment">// of allocation. This means we know the heap is completely</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>			<span class="comment">// exhausted. Otherwise, the heap still might have free</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>			<span class="comment">// space in it, just not enough contiguous space to</span>
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>			<span class="comment">// accommodate npages.</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>			p.searchAddr = maxSearchAddr()
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>		}
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>		return 0, 0
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	}
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>Found:
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>	<span class="comment">// Go ahead and actually mark the bits now that we have an address.</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>	scav = p.allocRange(addr, npages)
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>	<span class="comment">// If we found a higher searchAddr, we know that all the</span>
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>	<span class="comment">// heap memory before that searchAddr in an offset address space is</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	<span class="comment">// allocated, so bump p.searchAddr up to the new one.</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	if p.searchAddr.lessThan(searchAddr) {
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>		p.searchAddr = searchAddr
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	}
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	return addr, scav
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>}
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span><span class="comment">// free returns npages worth of memory starting at base back to the page heap.</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span><span class="comment">// p.mheapLock must be held.</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span><span class="comment">// Must run on the system stack because p.mheapLock must be held.</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>func (p *pageAlloc) free(base, npages uintptr) {
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>	assertLockHeld(p.mheapLock)
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;re freeing pages below the p.searchAddr, update searchAddr.</span>
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	if b := (offAddr{base}); b.lessThan(p.searchAddr) {
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		p.searchAddr = b
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>	}
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>	limit := base + npages*pageSize - 1
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>	if npages == 1 {
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>		<span class="comment">// Fast path: we&#39;re clearing a single bit, and we know exactly</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>		<span class="comment">// where it is, so mark it directly.</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>		i := chunkIndex(base)
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>		pi := chunkPageIndex(base)
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		p.chunkOf(i).free1(pi)
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>		p.scav.index.free(i, pi, 1)
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>	} else {
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>		<span class="comment">// Slow path: we&#39;re clearing more bits so we may need to iterate.</span>
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>		sc, ec := chunkIndex(base), chunkIndex(limit)
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>		si, ei := chunkPageIndex(base), chunkPageIndex(limit)
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>		if sc == ec {
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>			<span class="comment">// The range doesn&#39;t cross any chunk boundaries.</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>			p.chunkOf(sc).free(si, ei+1-si)
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>			p.scav.index.free(sc, si, ei+1-si)
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>		} else {
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>			<span class="comment">// The range crosses at least one chunk boundary.</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>			p.chunkOf(sc).free(si, pallocChunkPages-si)
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>			p.scav.index.free(sc, si, pallocChunkPages-si)
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>			for c := sc + 1; c &lt; ec; c++ {
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>				p.chunkOf(c).freeAll()
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>				p.scav.index.free(c, 0, pallocChunkPages)
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>			}
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>			p.chunkOf(ec).free(0, ei+1)
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>			p.scav.index.free(ec, 0, ei+1)
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>		}
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>	}
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>	p.update(base, npages, true, false)
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>}
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>const (
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	pallocSumBytes = unsafe.Sizeof(pallocSum(0))
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	<span class="comment">// maxPackedValue is the maximum value that any of the three fields in</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>	<span class="comment">// the pallocSum may take on.</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	maxPackedValue    = 1 &lt;&lt; logMaxPackedValue
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	logMaxPackedValue = logPallocChunkPages + (summaryLevels-1)*summaryLevelBits
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>	freeChunkSum = pallocSum(uint64(pallocChunkPages) |
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>		uint64(pallocChunkPages&lt;&lt;logMaxPackedValue) |
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>		uint64(pallocChunkPages&lt;&lt;(2*logMaxPackedValue)))
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>)
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span><span class="comment">// pallocSum is a packed summary type which packs three numbers: start, max,</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span><span class="comment">// and end into a single 8-byte value. Each of these values are a summary of</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span><span class="comment">// a bitmap and are thus counts, each of which may have a maximum value of</span>
<span id="L991" class="ln">   991&nbsp;&nbsp;</span><span class="comment">// 2^21 - 1, or all three may be equal to 2^21. The latter case is represented</span>
<span id="L992" class="ln">   992&nbsp;&nbsp;</span><span class="comment">// by just setting the 64th bit.</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>type pallocSum uint64
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span><span class="comment">// packPallocSum takes a start, max, and end value and produces a pallocSum.</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>func packPallocSum(start, max, end uint) pallocSum {
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	if max == maxPackedValue {
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>		return pallocSum(uint64(1 &lt;&lt; 63))
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	}
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	return pallocSum((uint64(start) &amp; (maxPackedValue - 1)) |
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>		((uint64(max) &amp; (maxPackedValue - 1)) &lt;&lt; logMaxPackedValue) |
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>		((uint64(end) &amp; (maxPackedValue - 1)) &lt;&lt; (2 * logMaxPackedValue)))
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>}
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span><span class="comment">// start extracts the start value from a packed sum.</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>func (p pallocSum) start() uint {
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	if uint64(p)&amp;uint64(1&lt;&lt;63) != 0 {
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>		return maxPackedValue
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>	}
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>	return uint(uint64(p) &amp; (maxPackedValue - 1))
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>}
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span><span class="comment">// max extracts the max value from a packed sum.</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>func (p pallocSum) max() uint {
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	if uint64(p)&amp;uint64(1&lt;&lt;63) != 0 {
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>		return maxPackedValue
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	}
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	return uint((uint64(p) &gt;&gt; logMaxPackedValue) &amp; (maxPackedValue - 1))
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>}
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span><span class="comment">// end extracts the end value from a packed sum.</span>
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>func (p pallocSum) end() uint {
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	if uint64(p)&amp;uint64(1&lt;&lt;63) != 0 {
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>		return maxPackedValue
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>	}
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	return uint((uint64(p) &gt;&gt; (2 * logMaxPackedValue)) &amp; (maxPackedValue - 1))
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>}
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span><span class="comment">// unpack unpacks all three values from the summary.</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>func (p pallocSum) unpack() (uint, uint, uint) {
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	if uint64(p)&amp;uint64(1&lt;&lt;63) != 0 {
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>		return maxPackedValue, maxPackedValue, maxPackedValue
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>	}
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>	return uint(uint64(p) &amp; (maxPackedValue - 1)),
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>		uint((uint64(p) &gt;&gt; logMaxPackedValue) &amp; (maxPackedValue - 1)),
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>		uint((uint64(p) &gt;&gt; (2 * logMaxPackedValue)) &amp; (maxPackedValue - 1))
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>}
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span><span class="comment">// mergeSummaries merges consecutive summaries which may each represent at</span>
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span><span class="comment">// most 1 &lt;&lt; logMaxPagesPerSum pages each together into one.</span>
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>func mergeSummaries(sums []pallocSum, logMaxPagesPerSum uint) pallocSum {
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>	<span class="comment">// Merge the summaries in sums into one.</span>
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>	<span class="comment">// We do this by keeping a running summary representing the merged</span>
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>	<span class="comment">// summaries of sums[:i] in start, most, and end.</span>
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>	start, most, end := sums[0].unpack()
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>	for i := 1; i &lt; len(sums); i++ {
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>		<span class="comment">// Merge in sums[i].</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>		si, mi, ei := sums[i].unpack()
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>		<span class="comment">// Merge in sums[i].start only if the running summary is</span>
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>		<span class="comment">// completely free, otherwise this summary&#39;s start</span>
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>		<span class="comment">// plays no role in the combined sum.</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>		if start == uint(i)&lt;&lt;logMaxPagesPerSum {
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>			start += si
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>		}
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>		<span class="comment">// Recompute the max value of the running sum by looking</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>		<span class="comment">// across the boundary between the running sum and sums[i]</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>		<span class="comment">// and at the max sums[i], taking the greatest of those two</span>
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>		<span class="comment">// and the max of the running sum.</span>
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>		most = max(most, end+si, mi)
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>		<span class="comment">// Merge in end by checking if this new summary is totally</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>		<span class="comment">// free. If it is, then we want to extend the running sum&#39;s</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>		<span class="comment">// end by the new summary. If not, then we have some alloc&#39;d</span>
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>		<span class="comment">// pages in there and we just want to take the end value in</span>
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>		<span class="comment">// sums[i].</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>		if ei == 1&lt;&lt;logMaxPagesPerSum {
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>			end += 1 &lt;&lt; logMaxPagesPerSum
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>		} else {
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>			end = ei
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>		}
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	}
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	return packPallocSum(start, most, end)
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>}
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>
</pre><p><a href="mpagealloc.go?m=text">View as plain text</a></p>

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
