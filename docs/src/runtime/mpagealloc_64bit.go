<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mpagealloc_64bit.go - Go Documentation Server</title>

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
<a href="mpagealloc_64bit.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mpagealloc_64bit.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build amd64 || arm64 || loong64 || mips64 || mips64le || ppc64 || ppc64le || riscv64 || s390x</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>const (
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	<span class="comment">// The number of levels in the radix tree.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	summaryLevels = 5
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">// Constants for testing.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	pageAlloc32Bit = 0
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	pageAlloc64Bit = 1
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// Number of bits needed to represent all indices into the L1 of the</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// chunks map.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// See (*pageAlloc).chunks for more details. Update the documentation</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// there should this number change.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	pallocChunksL1Bits = 13
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// levelBits is the number of bits in the radix for a given level in the super summary</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// structure.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// The sum of all the entries of levelBits should equal heapAddrBits.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>var levelBits = [summaryLevels]uint{
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	summaryL0Bits,
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	summaryLevelBits,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	summaryLevelBits,
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	summaryLevelBits,
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	summaryLevelBits,
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// levelShift is the number of bits to shift to acquire the radix for a given level</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// in the super summary structure.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// With levelShift, one can compute the index of the summary at level l related to a</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// pointer p by doing:</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//	p &gt;&gt; levelShift[l]</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>var levelShift = [summaryLevels]uint{
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	heapAddrBits - summaryL0Bits,
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	heapAddrBits - summaryL0Bits - 1*summaryLevelBits,
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	heapAddrBits - summaryL0Bits - 2*summaryLevelBits,
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	heapAddrBits - summaryL0Bits - 3*summaryLevelBits,
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	heapAddrBits - summaryL0Bits - 4*summaryLevelBits,
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// levelLogPages is log2 the maximum number of runtime pages in the address space</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// a summary in the given level represents.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// The leaf level always represents exactly log2 of 1 chunk&#39;s worth of pages.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>var levelLogPages = [summaryLevels]uint{
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	logPallocChunkPages + 4*summaryLevelBits,
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	logPallocChunkPages + 3*summaryLevelBits,
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	logPallocChunkPages + 2*summaryLevelBits,
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	logPallocChunkPages + 1*summaryLevelBits,
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	logPallocChunkPages,
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// sysInit performs architecture-dependent initialization of fields</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// in pageAlloc. pageAlloc should be uninitialized except for sysStat</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// if any runtime statistic should be updated.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>func (p *pageAlloc) sysInit(test bool) {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// Reserve memory for each level. This will get mapped in</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// as R/W by setArenas.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	for l, shift := range levelShift {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		entries := 1 &lt;&lt; (heapAddrBits - shift)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		<span class="comment">// Reserve b bytes of memory anywhere in the address space.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		b := alignUp(uintptr(entries)*pallocSumBytes, physPageSize)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		r := sysReserve(nil, b)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		if r == nil {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			throw(&#34;failed to reserve page summary memory&#34;)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		<span class="comment">// Put this reservation into a slice.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		sl := notInHeapSlice{(*notInHeap)(r), 0, entries}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		p.summary[l] = *(*[]pallocSum)(unsafe.Pointer(&amp;sl))
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// sysGrow performs architecture-dependent operations on heap</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// growth for the page allocator, such as mapping in new memory</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// for summaries. It also updates the length of the slices in</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// p.summary.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// base is the base of the newly-added heap memory and limit is</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// the first address past the end of the newly-added heap memory.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// Both must be aligned to pallocChunkBytes.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// The caller must update p.start and p.end after calling sysGrow.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>func (p *pageAlloc) sysGrow(base, limit uintptr) {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	if base%pallocChunkBytes != 0 || limit%pallocChunkBytes != 0 {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		print(&#34;runtime: base = &#34;, hex(base), &#34;, limit = &#34;, hex(limit), &#34;\n&#34;)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		throw(&#34;sysGrow bounds not aligned to pallocChunkBytes&#34;)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	<span class="comment">// addrRangeToSummaryRange converts a range of addresses into a range</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// of summary indices which must be mapped to support those addresses</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">// in the summary range.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	addrRangeToSummaryRange := func(level int, r addrRange) (int, int) {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		sumIdxBase, sumIdxLimit := addrsToSummaryRange(level, r.base.addr(), r.limit.addr())
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		return blockAlignSummaryRange(level, sumIdxBase, sumIdxLimit)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// summaryRangeToSumAddrRange converts a range of indices in any</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// level of p.summary into page-aligned addresses which cover that</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// range of indices.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	summaryRangeToSumAddrRange := func(level, sumIdxBase, sumIdxLimit int) addrRange {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		baseOffset := alignDown(uintptr(sumIdxBase)*pallocSumBytes, physPageSize)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		limitOffset := alignUp(uintptr(sumIdxLimit)*pallocSumBytes, physPageSize)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		base := unsafe.Pointer(&amp;p.summary[level][0])
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		return addrRange{
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			offAddr{uintptr(add(base, baseOffset))},
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			offAddr{uintptr(add(base, limitOffset))},
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// addrRangeToSumAddrRange is a convenience function that converts</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// an address range r to the address range of the given summary level</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// that stores the summaries for r.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	addrRangeToSumAddrRange := func(level int, r addrRange) addrRange {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		sumIdxBase, sumIdxLimit := addrRangeToSummaryRange(level, r)
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		return summaryRangeToSumAddrRange(level, sumIdxBase, sumIdxLimit)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// Find the first inUse index which is strictly greater than base.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// Because this function will never be asked remap the same memory</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// twice, this index is effectively the index at which we would insert</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// this new growth, and base will never overlap/be contained within</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// any existing range.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// This will be used to look at what memory in the summary array is already</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">// mapped before and after this new range.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	inUseIndex := p.inUse.findSucc(base)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// Walk up the radix tree and map summaries in as needed.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	for l := range p.summary {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		<span class="comment">// Figure out what part of the summary array this new address space needs.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		needIdxBase, needIdxLimit := addrRangeToSummaryRange(l, makeAddrRange(base, limit))
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		<span class="comment">// Update the summary slices with a new upper-bound. This ensures</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		<span class="comment">// we get tight bounds checks on at least the top bound.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		<span class="comment">// We must do this regardless of whether we map new memory.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		if needIdxLimit &gt; len(p.summary[l]) {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			p.summary[l] = p.summary[l][:needIdxLimit]
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		<span class="comment">// Compute the needed address range in the summary array for level l.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		need := summaryRangeToSumAddrRange(l, needIdxBase, needIdxLimit)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		<span class="comment">// Prune need down to what needs to be newly mapped. Some parts of it may</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		<span class="comment">// already be mapped by what inUse describes due to page alignment requirements</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		<span class="comment">// for mapping. Because this function will never be asked to remap the same</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		<span class="comment">// memory twice, it should never be possible to prune in such a way that causes</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		<span class="comment">// need to be split.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		if inUseIndex &gt; 0 {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			need = need.subtract(addrRangeToSumAddrRange(l, p.inUse.ranges[inUseIndex-1]))
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		if inUseIndex &lt; len(p.inUse.ranges) {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			need = need.subtract(addrRangeToSumAddrRange(l, p.inUse.ranges[inUseIndex]))
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		<span class="comment">// It&#39;s possible that after our pruning above, there&#39;s nothing new to map.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		if need.size() == 0 {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			continue
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		<span class="comment">// Map and commit need.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		sysMap(unsafe.Pointer(need.base.addr()), need.size(), p.sysStat)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		sysUsed(unsafe.Pointer(need.base.addr()), need.size(), need.size())
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		p.summaryMappedReady += need.size()
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// Update the scavenge index.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	p.summaryMappedReady += p.scav.index.sysGrow(base, limit, p.sysStat)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// sysGrow increases the index&#39;s backing store in response to a heap growth.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// Returns the amount of memory added to sysStat.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>func (s *scavengeIndex) sysGrow(base, limit uintptr, sysStat *sysMemStat) uintptr {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	if base%pallocChunkBytes != 0 || limit%pallocChunkBytes != 0 {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		print(&#34;runtime: base = &#34;, hex(base), &#34;, limit = &#34;, hex(limit), &#34;\n&#34;)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		throw(&#34;sysGrow bounds not aligned to pallocChunkBytes&#34;)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	scSize := unsafe.Sizeof(atomicScavChunkData{})
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// Map and commit the pieces of chunks that we need.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// We always map the full range of the minimum heap address to the</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// maximum heap address. We don&#39;t do this for the summary structure</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// because it&#39;s quite large and a discontiguous heap could cause a</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// lot of memory to be used. In this situation, the worst case overhead</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	<span class="comment">// is in the single-digit MiB if we map the whole thing.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// The base address of the backing store is always page-aligned,</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// because it comes from the OS, so it&#39;s sufficient to align the</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// index.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	haveMin := s.min.Load()
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	haveMax := s.max.Load()
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	needMin := alignDown(uintptr(chunkIndex(base)), physPageSize/scSize)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	needMax := alignUp(uintptr(chunkIndex(limit)), physPageSize/scSize)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// We need a contiguous range, so extend the range if there&#39;s no overlap.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	if needMax &lt; haveMin {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		needMax = haveMin
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	if haveMax != 0 &amp;&amp; needMin &gt; haveMax {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		needMin = haveMax
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// Avoid a panic from indexing one past the last element.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	chunksBase := uintptr(unsafe.Pointer(&amp;s.chunks[0]))
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	have := makeAddrRange(chunksBase+haveMin*scSize, chunksBase+haveMax*scSize)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	need := makeAddrRange(chunksBase+needMin*scSize, chunksBase+needMax*scSize)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// Subtract any overlap from rounding. We can&#39;t re-map memory because</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	<span class="comment">// it&#39;ll be zeroed.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	need = need.subtract(have)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;ve got something to map, map it, and update the slice bounds.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	if need.size() != 0 {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		sysMap(unsafe.Pointer(need.base.addr()), need.size(), sysStat)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		sysUsed(unsafe.Pointer(need.base.addr()), need.size(), need.size())
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		<span class="comment">// Update the indices only after the new memory is valid.</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		if haveMax == 0 || needMin &lt; haveMin {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			s.min.Store(needMin)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		if needMax &gt; haveMax {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			s.max.Store(needMax)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	return need.size()
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// sysInit initializes the scavengeIndex&#39; chunks array.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// Returns the amount of memory added to sysStat.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>func (s *scavengeIndex) sysInit(test bool, sysStat *sysMemStat) uintptr {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	n := uintptr(1&lt;&lt;heapAddrBits) / pallocChunkBytes
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	nbytes := n * unsafe.Sizeof(atomicScavChunkData{})
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	r := sysReserve(nil, nbytes)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	sl := notInHeapSlice{(*notInHeap)(r), int(n), int(n)}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	s.chunks = *(*[]atomicScavChunkData)(unsafe.Pointer(&amp;sl))
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	return 0 <span class="comment">// All memory above is mapped Reserved.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
</pre><p><a href="mpagealloc_64bit.go?m=text">View as plain text</a></p>

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
