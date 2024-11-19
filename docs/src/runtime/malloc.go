<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/malloc.go - Go Documentation Server</title>

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
<a href="malloc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">malloc.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Memory allocator.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This was originally based on tcmalloc, but has diverged quite a bit.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// http://goog-perftools.sourceforge.net/doc/tcmalloc.html</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// The main allocator works in runs of pages.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// Small allocation sizes (up to and including 32 kB) are</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// rounded to one of about 70 size classes, each of which</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// has its own free set of objects of exactly that size.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// Any free page of memory can be split into a set of objects</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// of one size class, which are then managed using a free bitmap.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// The allocator&#39;s data structures are:</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//	fixalloc: a free-list allocator for fixed-size off-heap objects,</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//		used to manage storage used by the allocator.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//	mheap: the malloc heap, managed at page (8192-byte) granularity.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//	mspan: a run of in-use pages managed by the mheap.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//	mcentral: collects all spans of a given size class.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//	mcache: a per-P cache of mspans with free space.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//	mstats: allocation statistics.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// Allocating a small object proceeds up a hierarchy of caches:</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//	1. Round the size up to one of the small size classes</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//	   and look in the corresponding mspan in this P&#39;s mcache.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//	   Scan the mspan&#39;s free bitmap to find a free slot.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//	   If there is a free slot, allocate it.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//	   This can all be done without acquiring a lock.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//	2. If the mspan has no free slots, obtain a new mspan</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//	   from the mcentral&#39;s list of mspans of the required size</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//	   class that have free space.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//	   Obtaining a whole span amortizes the cost of locking</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//	   the mcentral.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//	3. If the mcentral&#39;s mspan list is empty, obtain a run</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//	   of pages from the mheap to use for the mspan.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//	4. If the mheap is empty or has no page runs large enough,</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//	   allocate a new group of pages (at least 1MB) from the</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//	   operating system. Allocating a large run of pages</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//	   amortizes the cost of talking to the operating system.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// Sweeping an mspan and freeing objects on it proceeds up a similar</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// hierarchy:</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//	1. If the mspan is being swept in response to allocation, it</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//	   is returned to the mcache to satisfy the allocation.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//	2. Otherwise, if the mspan still has allocated objects in it,</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//	   it is placed on the mcentral free list for the mspan&#39;s size</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//	   class.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//	3. Otherwise, if all objects in the mspan are free, the mspan&#39;s</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//	   pages are returned to the mheap and the mspan is now dead.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// Allocating and freeing a large object uses the mheap</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// directly, bypassing the mcache and mcentral.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// If mspan.needzero is false, then free object slots in the mspan are</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// already zeroed. Otherwise if needzero is true, objects are zeroed as</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// they are allocated. There are various benefits to delaying zeroing</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// this way:</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">//	1. Stack frame allocation can avoid zeroing altogether.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//	2. It exhibits better temporal locality, since the program is</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//	   probably about to write to the memory.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//	3. We don&#39;t zero pages that never get reused.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// Virtual memory layout</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// The heap consists of a set of arenas, which are 64MB on 64-bit and</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// 4MB on 32-bit (heapArenaBytes). Each arena&#39;s start address is also</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// aligned to the arena size.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// Each arena has an associated heapArena object that stores the</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// metadata for that arena: the heap bitmap for all words in the arena</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// and the span map for all pages in the arena. heapArena objects are</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// themselves allocated off-heap.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// Since arenas are aligned, the address space can be viewed as a</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// series of arena frames. The arena map (mheap_.arenas) maps from</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// arena frame number to *heapArena, or nil for parts of the address</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// space not backed by the Go heap. The arena map is structured as a</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// two-level array consisting of a &#34;L1&#34; arena map and many &#34;L2&#34; arena</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// maps; however, since arenas are large, on many architectures, the</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// arena map consists of a single, large L2 map.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// The arena map covers the entire possible address space, allowing</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// the Go heap to use any part of the address space. The allocator</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// attempts to keep arenas contiguous so that large spans (and hence</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// large objects) can cross arenas.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>package runtime
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>import (
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	&#34;internal/goexperiment&#34;
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	&#34;internal/goos&#34;
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	&#34;runtime/internal/math&#34;
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>const (
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	maxTinySize   = _TinySize
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	tinySizeClass = _TinySizeClass
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	maxSmallSize  = _MaxSmallSize
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	pageShift = _PageShift
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	pageSize  = _PageSize
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	_PageSize = 1 &lt;&lt; _PageShift
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	_PageMask = _PageSize - 1
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// _64bit = 1 on 64-bit systems, 0 on 32-bit systems</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	_64bit = 1 &lt;&lt; (^uintptr(0) &gt;&gt; 63) / 2
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// Tiny allocator parameters, see &#34;Tiny allocator&#34; comment in malloc.go.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	_TinySize      = 16
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	_TinySizeClass = int8(2)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	_FixAllocChunk = 16 &lt;&lt; 10 <span class="comment">// Chunk size for FixAlloc</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// Per-P, per order stack segment cache size.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	_StackCacheSize = 32 * 1024
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// Number of orders that get caching. Order 0 is FixedStack</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// and each successive order is twice as large.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// We want to cache 2KB, 4KB, 8KB, and 16KB stacks. Larger stacks</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// will be allocated directly.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// Since FixedStack is different on different systems, we</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// must vary NumStackOrders to keep the same maximum cached size.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">//   OS               | FixedStack | NumStackOrders</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">//   -----------------+------------+---------------</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">//   linux/darwin/bsd | 2KB        | 4</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">//   windows/32       | 4KB        | 3</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">//   windows/64       | 8KB        | 2</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">//   plan9            | 4KB        | 3</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	_NumStackOrders = 4 - goarch.PtrSize/4*goos.IsWindows - 1*goos.IsPlan9
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// heapAddrBits is the number of bits in a heap address. On</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// amd64, addresses are sign-extended beyond heapAddrBits. On</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// other arches, they are zero-extended.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// On most 64-bit platforms, we limit this to 48 bits based on a</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// combination of hardware and OS limitations.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// amd64 hardware limits addresses to 48 bits, sign-extended</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// to 64 bits. Addresses where the top 16 bits are not either</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// all 0 or all 1 are &#34;non-canonical&#34; and invalid. Because of</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// these &#34;negative&#34; addresses, we offset addresses by 1&lt;&lt;47</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// (arenaBaseOffset) on amd64 before computing indexes into</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// the heap arenas index. In 2017, amd64 hardware added</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// support for 57 bit addresses; however, currently only Linux</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// supports this extension and the kernel will never choose an</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// address above 1&lt;&lt;47 unless mmap is called with a hint</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// address above 1&lt;&lt;47 (which we never do).</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// arm64 hardware (as of ARMv8) limits user addresses to 48</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// bits, in the range [0, 1&lt;&lt;48).</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// ppc64, mips64, and s390x support arbitrary 64 bit addresses</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// in hardware. On Linux, Go leans on stricter OS limits. Based</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// on Linux&#39;s processor.h, the user address space is limited as</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// follows on 64-bit architectures:</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// Architecture  Name              Maximum Value (exclusive)</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// ---------------------------------------------------------------------</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// amd64         TASK_SIZE_MAX     0x007ffffffff000 (47 bit addresses)</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// arm64         TASK_SIZE_64      0x01000000000000 (48 bit addresses)</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// ppc64{,le}    TASK_SIZE_USER64  0x00400000000000 (46 bit addresses)</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// mips64{,le}   TASK_SIZE64       0x00010000000000 (40 bit addresses)</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// s390x         TASK_SIZE         1&lt;&lt;64 (64 bit addresses)</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// These limits may increase over time, but are currently at</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// most 48 bits except on s390x. On all architectures, Linux</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// starts placing mmap&#39;d regions at addresses that are</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// significantly below 48 bits, so even if it&#39;s possible to</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">// exceed Go&#39;s 48 bit limit, it&#39;s extremely unlikely in</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// practice.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// On 32-bit platforms, we accept the full 32-bit address</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// space because doing so is cheap.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// mips32 only has access to the low 2GB of virtual memory, so</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// we further limit it to 31 bits.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// On ios/arm64, although 64-bit pointers are presumably</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// available, pointers are truncated to 33 bits in iOS &lt;14.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// Furthermore, only the top 4 GiB of the address space are</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// actually available to the application. In iOS &gt;=14, more</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// of the address space is available, and the OS can now</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// provide addresses outside of those 33 bits. Pick 40 bits</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// as a reasonable balance between address space usage by the</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	<span class="comment">// page allocator, and flexibility for what mmap&#39;d regions</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// we&#39;ll accept for the heap. We can&#39;t just move to the full</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// 48 bits because this uses too much address space for older</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// iOS versions.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Once iOS &lt;14 is deprecated, promote ios/arm64</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// to a 48-bit address space like every other arm64 platform.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">// WebAssembly currently has a limit of 4GB linear memory.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	heapAddrBits = (_64bit*(1-goarch.IsWasm)*(1-goos.IsIos*goarch.IsArm64))*48 + (1-_64bit+goarch.IsWasm)*(32-(goarch.IsMips+goarch.IsMipsle)) + 40*goos.IsIos*goarch.IsArm64
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// maxAlloc is the maximum size of an allocation. On 64-bit,</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// it&#39;s theoretically possible to allocate 1&lt;&lt;heapAddrBits bytes. On</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// 32-bit, however, this is one less than 1&lt;&lt;32 because the</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	<span class="comment">// number of bytes in the address space doesn&#39;t actually fit</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	<span class="comment">// in a uintptr.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	maxAlloc = (1 &lt;&lt; heapAddrBits) - (1-_64bit)*1
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// The number of bits in a heap address, the size of heap</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// arenas, and the L1 and L2 arena map sizes are related by</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	<span class="comment">//   (1 &lt;&lt; addr bits) = arena size * L1 entries * L2 entries</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	<span class="comment">// Currently, we balance these as follows:</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	<span class="comment">//       Platform  Addr bits  Arena size  L1 entries   L2 entries</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// --------------  ---------  ----------  ----------  -----------</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">//       */64-bit         48        64MB           1    4M (32MB)</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// windows/64-bit         48         4MB          64    1M  (8MB)</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">//      ios/arm64         33         4MB           1  2048  (8KB)</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">//       */32-bit         32         4MB           1  1024  (4KB)</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">//     */mips(le)         31         4MB           1   512  (2KB)</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">// heapArenaBytes is the size of a heap arena. The heap</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// consists of mappings of size heapArenaBytes, aligned to</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// heapArenaBytes. The initial heap mapping is one arena.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	<span class="comment">// This is currently 64MB on 64-bit non-Windows and 4MB on</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	<span class="comment">// 32-bit and on Windows. We use smaller arenas on Windows</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// because all committed memory is charged to the process,</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	<span class="comment">// even if it&#39;s not touched. Hence, for processes with small</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// heaps, the mapped arena space needs to be commensurate.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// This is particularly important with the race detector,</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// since it significantly amplifies the cost of committed</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">// memory.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	heapArenaBytes = 1 &lt;&lt; logHeapArenaBytes
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	heapArenaWords = heapArenaBytes / goarch.PtrSize
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// logHeapArenaBytes is log_2 of heapArenaBytes. For clarity,</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// prefer using heapArenaBytes where possible (we need the</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// constant to compute some other constants).</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	logHeapArenaBytes = (6+20)*(_64bit*(1-goos.IsWindows)*(1-goarch.IsWasm)*(1-goos.IsIos*goarch.IsArm64)) + (2+20)*(_64bit*goos.IsWindows) + (2+20)*(1-_64bit) + (2+20)*goarch.IsWasm + (2+20)*goos.IsIos*goarch.IsArm64
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// heapArenaBitmapWords is the size of each heap arena&#39;s bitmap in uintptrs.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	heapArenaBitmapWords = heapArenaWords / (8 * goarch.PtrSize)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	pagesPerArena = heapArenaBytes / pageSize
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// arenaL1Bits is the number of bits of the arena number</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// covered by the first level arena map.</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// This number should be small, since the first level arena</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	<span class="comment">// map requires PtrSize*(1&lt;&lt;arenaL1Bits) of space in the</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	<span class="comment">// binary&#39;s BSS. It can be zero, in which case the first level</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	<span class="comment">// index is effectively unused. There is a performance benefit</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// to this, since the generated code can be more efficient,</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	<span class="comment">// but comes at the cost of having a large L2 mapping.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	<span class="comment">// We use the L1 map on 64-bit Windows because the arena size</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	<span class="comment">// is small, but the address space is still 48 bits, and</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	<span class="comment">// there&#39;s a high cost to having a large L2.</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	arenaL1Bits = 6 * (_64bit * goos.IsWindows)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">// arenaL2Bits is the number of bits of the arena number</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	<span class="comment">// covered by the second level arena index.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	<span class="comment">// The size of each arena map allocation is proportional to</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	<span class="comment">// 1&lt;&lt;arenaL2Bits, so it&#39;s important that this not be too</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	<span class="comment">// large. 48 bits leads to 32MB arena index allocations, which</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	<span class="comment">// is about the practical threshold.</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	arenaL2Bits = heapAddrBits - logHeapArenaBytes - arenaL1Bits
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">// arenaL1Shift is the number of bits to shift an arena frame</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	<span class="comment">// number by to compute an index into the first level arena map.</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	arenaL1Shift = arenaL2Bits
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	<span class="comment">// arenaBits is the total bits in a combined arena map index.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">// This is split between the index into the L1 arena map and</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// the L2 arena map.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	arenaBits = arenaL1Bits + arenaL2Bits
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">// arenaBaseOffset is the pointer value that corresponds to</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// index 0 in the heap arena map.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// On amd64, the address space is 48 bits, sign extended to 64</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	<span class="comment">// bits. This offset lets us handle &#34;negative&#34; addresses (or</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">// high addresses if viewed as unsigned).</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">// On aix/ppc64, this offset allows to keep the heapAddrBits to</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	<span class="comment">// 48. Otherwise, it would be 60 in order to handle mmap addresses</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	<span class="comment">// (in range 0x0a00000000000000 - 0x0afffffffffffff). But in this</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	<span class="comment">// case, the memory reserved in (s *pageAlloc).init for chunks</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">// is causing important slowdowns.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	<span class="comment">// On other platforms, the user address space is contiguous</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	<span class="comment">// and starts at 0, so no offset is necessary.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	arenaBaseOffset = 0xffff800000000000*goarch.IsAmd64 + 0x0a00000000000000*goos.IsAix
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	<span class="comment">// A typed version of this constant that will make it into DWARF (for viewcore).</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	arenaBaseOffsetUintptr = uintptr(arenaBaseOffset)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	<span class="comment">// Max number of threads to run garbage collection.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	<span class="comment">// 2, 3, and 4 are all plausible maximums depending</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	<span class="comment">// on the hardware details of the machine. The garbage</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// collector scales well to 32 cpus.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	_MaxGcproc = 32
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	<span class="comment">// minLegalPointer is the smallest possible legal pointer.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	<span class="comment">// This is the smallest possible architectural page size,</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	<span class="comment">// since we assume that the first page is never mapped.</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	<span class="comment">// This should agree with minZeroPage in the compiler.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	minLegalPointer uintptr = 4096
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">// minHeapForMetadataHugePages sets a threshold on when certain kinds of</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// heap metadata, currently the arenas map L2 entries and page alloc bitmap</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	<span class="comment">// mappings, are allowed to be backed by huge pages. If the heap goal ever</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	<span class="comment">// exceeds this threshold, then huge pages are enabled.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	<span class="comment">// These numbers are chosen with the assumption that huge pages are on the</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">// order of a few MiB in size.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// The kind of metadata this applies to has a very low overhead when compared</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// to address space used, but their constant overheads for small heaps would</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// be very high if they were to be backed by huge pages (e.g. a few MiB makes</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// a huge difference for an 8 MiB heap, but barely any difference for a 1 GiB</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// heap). The benefit of huge pages is also not worth it for small heaps,</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// because only a very, very small part of the metadata is used for small heaps.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	<span class="comment">// N.B. If the heap goal exceeds the threshold then shrinks to a very small size</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	<span class="comment">// again, then huge pages will still be enabled for this mapping. The reason is that</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	<span class="comment">// there&#39;s no point unless we&#39;re also returning the physical memory for these</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	<span class="comment">// metadata mappings back to the OS. That would be quite complex to do in general</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	<span class="comment">// as the heap is likely fragmented after a reduction in heap size.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	minHeapForMetadataHugePages = 1 &lt;&lt; 30
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// physPageSize is the size in bytes of the OS&#39;s physical pages.</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// Mapping and unmapping operations must be done at multiples of</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// physPageSize.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// This must be set by the OS init code (typically in osinit) before</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// mallocinit.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>var physPageSize uintptr
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// physHugePageSize is the size in bytes of the OS&#39;s default physical huge</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">// page size whose allocation is opaque to the application. It is assumed</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// and verified to be a power of two.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// If set, this must be set by the OS init code (typically in osinit) before</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// mallocinit. However, setting it at all is optional, and leaving the default</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span><span class="comment">// value is always safe (though potentially less efficient).</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// Since physHugePageSize is always assumed to be a power of two,</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// physHugePageShift is defined as physHugePageSize == 1 &lt;&lt; physHugePageShift.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">// The purpose of physHugePageShift is to avoid doing divisions in</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">// performance critical functions.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>var (
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	physHugePageSize  uintptr
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	physHugePageShift uint
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>func mallocinit() {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	if class_to_size[_TinySizeClass] != _TinySize {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		throw(&#34;bad TinySizeClass&#34;)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	if heapArenaBitmapWords&amp;(heapArenaBitmapWords-1) != 0 {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		<span class="comment">// heapBits expects modular arithmetic on bitmap</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		<span class="comment">// addresses to work.</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		throw(&#34;heapArenaBitmapWords not a power of 2&#34;)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	}
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	<span class="comment">// Check physPageSize.</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	if physPageSize == 0 {
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		<span class="comment">// The OS init code failed to fetch the physical page size.</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		throw(&#34;failed to get system page size&#34;)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	if physPageSize &gt; maxPhysPageSize {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		print(&#34;system page size (&#34;, physPageSize, &#34;) is larger than maximum page size (&#34;, maxPhysPageSize, &#34;)\n&#34;)
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		throw(&#34;bad system page size&#34;)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	if physPageSize &lt; minPhysPageSize {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		print(&#34;system page size (&#34;, physPageSize, &#34;) is smaller than minimum page size (&#34;, minPhysPageSize, &#34;)\n&#34;)
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		throw(&#34;bad system page size&#34;)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	}
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	if physPageSize&amp;(physPageSize-1) != 0 {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		print(&#34;system page size (&#34;, physPageSize, &#34;) must be a power of 2\n&#34;)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		throw(&#34;bad system page size&#34;)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	if physHugePageSize&amp;(physHugePageSize-1) != 0 {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		print(&#34;system huge page size (&#34;, physHugePageSize, &#34;) must be a power of 2\n&#34;)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		throw(&#34;bad system huge page size&#34;)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	if physHugePageSize &gt; maxPhysHugePageSize {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		<span class="comment">// physHugePageSize is greater than the maximum supported huge page size.</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		<span class="comment">// Don&#39;t throw here, like in the other cases, since a system configured</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		<span class="comment">// in this way isn&#39;t wrong, we just don&#39;t have the code to support them.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		<span class="comment">// Instead, silently set the huge page size to zero.</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		physHugePageSize = 0
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	if physHugePageSize != 0 {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		<span class="comment">// Since physHugePageSize is a power of 2, it suffices to increase</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		<span class="comment">// physHugePageShift until 1&lt;&lt;physHugePageShift == physHugePageSize.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		for 1&lt;&lt;physHugePageShift != physHugePageSize {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			physHugePageShift++
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	}
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	if pagesPerArena%pagesPerSpanRoot != 0 {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		print(&#34;pagesPerArena (&#34;, pagesPerArena, &#34;) is not divisible by pagesPerSpanRoot (&#34;, pagesPerSpanRoot, &#34;)\n&#34;)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		throw(&#34;bad pagesPerSpanRoot&#34;)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	if pagesPerArena%pagesPerReclaimerChunk != 0 {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		print(&#34;pagesPerArena (&#34;, pagesPerArena, &#34;) is not divisible by pagesPerReclaimerChunk (&#34;, pagesPerReclaimerChunk, &#34;)\n&#34;)
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		throw(&#34;bad pagesPerReclaimerChunk&#34;)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	if goexperiment.AllocHeaders {
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		<span class="comment">// Check that the minimum size (exclusive) for a malloc header is also</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		<span class="comment">// a size class boundary. This is important to making sure checks align</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		<span class="comment">// across different parts of the runtime.</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		minSizeForMallocHeaderIsSizeClass := false
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		for i := 0; i &lt; len(class_to_size); i++ {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			if minSizeForMallocHeader == uintptr(class_to_size[i]) {
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>				minSizeForMallocHeaderIsSizeClass = true
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>				break
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		if !minSizeForMallocHeaderIsSizeClass {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			throw(&#34;min size of malloc header is not a size class boundary&#34;)
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		<span class="comment">// Check that the pointer bitmap for all small sizes without a malloc header</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		<span class="comment">// fits in a word.</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		if minSizeForMallocHeader/goarch.PtrSize &gt; 8*goarch.PtrSize {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			throw(&#34;max pointer/scan bitmap size for headerless objects is too large&#34;)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	if minTagBits &gt; taggedPointerBits {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		throw(&#34;taggedPointerbits too small&#34;)
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	<span class="comment">// Initialize the heap.</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	mheap_.init()
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	mcache0 = allocmcache()
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	lockInit(&amp;gcBitsArenas.lock, lockRankGcBitsArenas)
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	lockInit(&amp;profInsertLock, lockRankProfInsert)
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	lockInit(&amp;profBlockLock, lockRankProfBlock)
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	lockInit(&amp;profMemActiveLock, lockRankProfMemActive)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	for i := range profMemFutureLock {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		lockInit(&amp;profMemFutureLock[i], lockRankProfMemFuture)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	lockInit(&amp;globalAlloc.mutex, lockRankGlobalAlloc)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	<span class="comment">// Create initial arena growth hints.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	if goarch.PtrSize == 8 {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		<span class="comment">// On a 64-bit machine, we pick the following hints</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		<span class="comment">// because:</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		<span class="comment">// 1. Starting from the middle of the address space</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		<span class="comment">// makes it easier to grow out a contiguous range</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		<span class="comment">// without running in to some other mapping.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		<span class="comment">// 2. This makes Go heap addresses more easily</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		<span class="comment">// recognizable when debugging.</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		<span class="comment">// 3. Stack scanning in gccgo is still conservative,</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		<span class="comment">// so it&#39;s important that addresses be distinguishable</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		<span class="comment">// from other data.</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		<span class="comment">// Starting at 0x00c0 means that the valid memory addresses</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		<span class="comment">// will begin 0x00c0, 0x00c1, ...</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		<span class="comment">// In little-endian, that&#39;s c0 00, c1 00, ... None of those are valid</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		<span class="comment">// UTF-8 sequences, and they are otherwise as far away from</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		<span class="comment">// ff (likely a common byte) as possible. If that fails, we try other 0xXXc0</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		<span class="comment">// addresses. An earlier attempt to use 0x11f8 caused out of memory errors</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		<span class="comment">// on OS X during thread allocations.  0x00c0 causes conflicts with</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		<span class="comment">// AddressSanitizer which reserves all memory up to 0x0100.</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		<span class="comment">// These choices reduce the odds of a conservative garbage collector</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		<span class="comment">// not collecting memory because some non-pointer block of memory</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		<span class="comment">// had a bit pattern that matched a memory address.</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		<span class="comment">// However, on arm64, we ignore all this advice above and slam the</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		<span class="comment">// allocation at 0x40 &lt;&lt; 32 because when using 4k pages with 3-level</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		<span class="comment">// translation buffers, the user address space is limited to 39 bits</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		<span class="comment">// On ios/arm64, the address space is even smaller.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		<span class="comment">// On AIX, mmaps starts at 0x0A00000000000000 for 64-bit.</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		<span class="comment">// processes.</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		<span class="comment">// Space mapped for user arenas comes immediately after the range</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		<span class="comment">// originally reserved for the regular heap when race mode is not</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		<span class="comment">// enabled because user arena chunks can never be used for regular heap</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		<span class="comment">// allocations and we want to avoid fragmenting the address space.</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		<span class="comment">// In race mode we have no choice but to just use the same hints because</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		<span class="comment">// the race detector requires that the heap be mapped contiguously.</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		for i := 0x7f; i &gt;= 0; i-- {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			var p uintptr
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			switch {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			case raceenabled:
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>				<span class="comment">// The TSAN runtime requires the heap</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>				<span class="comment">// to be in the range [0x00c000000000,</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>				<span class="comment">// 0x00e000000000).</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>				p = uintptr(i)&lt;&lt;32 | uintptrMask&amp;(0x00c0&lt;&lt;32)
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>				if p &gt;= uintptrMask&amp;0x00e000000000 {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>					continue
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>				}
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>			case GOARCH == &#34;arm64&#34; &amp;&amp; GOOS == &#34;ios&#34;:
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>				p = uintptr(i)&lt;&lt;40 | uintptrMask&amp;(0x0013&lt;&lt;28)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			case GOARCH == &#34;arm64&#34;:
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>				p = uintptr(i)&lt;&lt;40 | uintptrMask&amp;(0x0040&lt;&lt;32)
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>			case GOOS == &#34;aix&#34;:
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>				if i == 0 {
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>					<span class="comment">// We don&#39;t use addresses directly after 0x0A00000000000000</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>					<span class="comment">// to avoid collisions with others mmaps done by non-go programs.</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>					continue
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>				}
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>				p = uintptr(i)&lt;&lt;40 | uintptrMask&amp;(0xa0&lt;&lt;52)
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			default:
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>				p = uintptr(i)&lt;&lt;40 | uintptrMask&amp;(0x00c0&lt;&lt;32)
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>			}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>			<span class="comment">// Switch to generating hints for user arenas if we&#39;ve gone</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>			<span class="comment">// through about half the hints. In race mode, take only about</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			<span class="comment">// a quarter; we don&#39;t have very much space to work with.</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			hintList := &amp;mheap_.arenaHints
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>			if (!raceenabled &amp;&amp; i &gt; 0x3f) || (raceenabled &amp;&amp; i &gt; 0x5f) {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>				hintList = &amp;mheap_.userArena.arenaHints
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>			}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			hint := (*arenaHint)(mheap_.arenaHintAlloc.alloc())
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>			hint.addr = p
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			hint.next, *hintList = *hintList, hint
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	} else {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		<span class="comment">// On a 32-bit machine, we&#39;re much more concerned</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		<span class="comment">// about keeping the usable heap contiguous.</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		<span class="comment">// Hence:</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		<span class="comment">// 1. We reserve space for all heapArenas up front so</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		<span class="comment">// they don&#39;t get interleaved with the heap. They&#39;re</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		<span class="comment">// ~258MB, so this isn&#39;t too bad. (We could reserve a</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		<span class="comment">// smaller amount of space up front if this is a</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		<span class="comment">// problem.)</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		<span class="comment">// 2. We hint the heap to start right above the end of</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		<span class="comment">// the binary so we have the best chance of keeping it</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		<span class="comment">// contiguous.</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		<span class="comment">// 3. We try to stake out a reasonably large initial</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		<span class="comment">// heap reservation.</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		const arenaMetaSize = (1 &lt;&lt; arenaBits) * unsafe.Sizeof(heapArena{})
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		meta := uintptr(sysReserve(nil, arenaMetaSize))
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		if meta != 0 {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>			mheap_.heapArenaAlloc.init(meta, arenaMetaSize, true)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		<span class="comment">// We want to start the arena low, but if we&#39;re linked</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		<span class="comment">// against C code, it&#39;s possible global constructors</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		<span class="comment">// have called malloc and adjusted the process&#39; brk.</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		<span class="comment">// Query the brk so we can avoid trying to map the</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		<span class="comment">// region over it (which will cause the kernel to put</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		<span class="comment">// the region somewhere else, likely at a high</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		<span class="comment">// address).</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		procBrk := sbrk0()
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		<span class="comment">// If we ask for the end of the data segment but the</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		<span class="comment">// operating system requires a little more space</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		<span class="comment">// before we can start allocating, it will give out a</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		<span class="comment">// slightly higher pointer. Except QEMU, which is</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		<span class="comment">// buggy, as usual: it won&#39;t adjust the pointer</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		<span class="comment">// upward. So adjust it upward a little bit ourselves:</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		<span class="comment">// 1/4 MB to get away from the running binary image.</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		p := firstmoduledata.end
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		if p &lt; procBrk {
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>			p = procBrk
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		if mheap_.heapArenaAlloc.next &lt;= p &amp;&amp; p &lt; mheap_.heapArenaAlloc.end {
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>			p = mheap_.heapArenaAlloc.end
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		}
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		p = alignUp(p+(256&lt;&lt;10), heapArenaBytes)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		<span class="comment">// Because we&#39;re worried about fragmentation on</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		<span class="comment">// 32-bit, we try to make a large initial reservation.</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		arenaSizes := []uintptr{
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>			512 &lt;&lt; 20,
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>			256 &lt;&lt; 20,
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>			128 &lt;&lt; 20,
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		for _, arenaSize := range arenaSizes {
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>			a, size := sysReserveAligned(unsafe.Pointer(p), arenaSize, heapArenaBytes)
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>			if a != nil {
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>				mheap_.arena.init(uintptr(a), size, false)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>				p = mheap_.arena.end <span class="comment">// For hint below</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>				break
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		hint := (*arenaHint)(mheap_.arenaHintAlloc.alloc())
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		hint.addr = p
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		hint.next, mheap_.arenaHints = mheap_.arenaHints, hint
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		<span class="comment">// Place the hint for user arenas just after the large reservation.</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		<span class="comment">// While this potentially competes with the hint above, in practice we probably</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		<span class="comment">// aren&#39;t going to be getting this far anyway on 32-bit platforms.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		userArenaHint := (*arenaHint)(mheap_.arenaHintAlloc.alloc())
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		userArenaHint.addr = p
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		userArenaHint.next, mheap_.userArena.arenaHints = mheap_.userArena.arenaHints, userArenaHint
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	}
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	<span class="comment">// Initialize the memory limit here because the allocator is going to look at it</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	<span class="comment">// but we haven&#39;t called gcinit yet and we&#39;re definitely going to allocate memory before then.</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	gcController.memoryLimit.Store(maxInt64)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span><span class="comment">// sysAlloc allocates heap arena space for at least n bytes. The</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span><span class="comment">// returned pointer is always heapArenaBytes-aligned and backed by</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span><span class="comment">// h.arenas metadata. The returned size is always a multiple of</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span><span class="comment">// heapArenaBytes. sysAlloc returns nil on failure.</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span><span class="comment">// There is no corresponding free function.</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span><span class="comment">// hintList is a list of hint addresses for where to allocate new</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span><span class="comment">// heap arenas. It must be non-nil.</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span><span class="comment">// register indicates whether the heap arena should be registered</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span><span class="comment">// in allArenas.</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span><span class="comment">// sysAlloc returns a memory region in the Reserved state. This region must</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span><span class="comment">// be transitioned to Prepared and then Ready before use.</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span><span class="comment">// h must be locked.</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>func (h *mheap) sysAlloc(n uintptr, hintList **arenaHint, register bool) (v unsafe.Pointer, size uintptr) {
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	assertLockHeld(&amp;h.lock)
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	n = alignUp(n, heapArenaBytes)
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	if hintList == &amp;h.arenaHints {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		<span class="comment">// First, try the arena pre-reservation.</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		<span class="comment">// Newly-used mappings are considered released.</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		<span class="comment">// Only do this if we&#39;re using the regular heap arena hints.</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		<span class="comment">// This behavior is only for the heap.</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		v = h.arena.alloc(n, heapArenaBytes, &amp;gcController.heapReleased)
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		if v != nil {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>			size = n
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>			goto mapped
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		}
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	}
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	<span class="comment">// Try to grow the heap at a hint address.</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	for *hintList != nil {
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		hint := *hintList
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		p := hint.addr
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		if hint.down {
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>			p -= n
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		}
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		if p+n &lt; p {
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>			<span class="comment">// We can&#39;t use this, so don&#39;t ask.</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>			v = nil
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		} else if arenaIndex(p+n-1) &gt;= 1&lt;&lt;arenaBits {
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>			<span class="comment">// Outside addressable heap. Can&#39;t use.</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>			v = nil
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>		} else {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>			v = sysReserve(unsafe.Pointer(p), n)
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		if p == uintptr(v) {
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>			<span class="comment">// Success. Update the hint.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>			if !hint.down {
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>				p += n
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>			}
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>			hint.addr = p
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>			size = n
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>			break
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		}
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		<span class="comment">// Failed. Discard this hint and try the next.</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		<span class="comment">// TODO: This would be cleaner if sysReserve could be</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		<span class="comment">// told to only return the requested address. In</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		<span class="comment">// particular, this is already how Windows behaves, so</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		<span class="comment">// it would simplify things there.</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		if v != nil {
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>			sysFreeOS(v, n)
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		}
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		*hintList = hint.next
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		h.arenaHintAlloc.free(unsafe.Pointer(hint))
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	}
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	if size == 0 {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		if raceenabled {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			<span class="comment">// The race detector assumes the heap lives in</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			<span class="comment">// [0x00c000000000, 0x00e000000000), but we</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>			<span class="comment">// just ran out of hints in this region. Give</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>			<span class="comment">// a nice failure.</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>			throw(&#34;too many address space collisions for -race mode&#34;)
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>		}
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		<span class="comment">// All of the hints failed, so we&#39;ll take any</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		<span class="comment">// (sufficiently aligned) address the kernel will give</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		<span class="comment">// us.</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>		v, size = sysReserveAligned(nil, n, heapArenaBytes)
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		if v == nil {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>			return nil, 0
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>		}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		<span class="comment">// Create new hints for extending this region.</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		hint := (*arenaHint)(h.arenaHintAlloc.alloc())
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		hint.addr, hint.down = uintptr(v), true
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		hint.next, mheap_.arenaHints = mheap_.arenaHints, hint
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		hint = (*arenaHint)(h.arenaHintAlloc.alloc())
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>		hint.addr = uintptr(v) + size
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>		hint.next, mheap_.arenaHints = mheap_.arenaHints, hint
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	}
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	<span class="comment">// Check for bad pointers or pointers we can&#39;t use.</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	{
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		var bad string
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		p := uintptr(v)
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		if p+size &lt; p {
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>			bad = &#34;region exceeds uintptr range&#34;
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>		} else if arenaIndex(p) &gt;= 1&lt;&lt;arenaBits {
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>			bad = &#34;base outside usable address space&#34;
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		} else if arenaIndex(p+size-1) &gt;= 1&lt;&lt;arenaBits {
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>			bad = &#34;end outside usable address space&#34;
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		}
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>		if bad != &#34;&#34; {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>			<span class="comment">// This should be impossible on most architectures,</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>			<span class="comment">// but it would be really confusing to debug.</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>			print(&#34;runtime: memory allocated by OS [&#34;, hex(p), &#34;, &#34;, hex(p+size), &#34;) not in usable address space: &#34;, bad, &#34;\n&#34;)
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			throw(&#34;memory reservation exceeds address space limit&#34;)
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		}
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	if uintptr(v)&amp;(heapArenaBytes-1) != 0 {
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		throw(&#34;misrounded allocation in sysAlloc&#34;)
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	}
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>mapped:
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	<span class="comment">// Create arena metadata.</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	for ri := arenaIndex(uintptr(v)); ri &lt;= arenaIndex(uintptr(v)+size-1); ri++ {
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		l2 := h.arenas[ri.l1()]
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		if l2 == nil {
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>			<span class="comment">// Allocate an L2 arena map.</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>			<span class="comment">// Use sysAllocOS instead of sysAlloc or persistentalloc because there&#39;s no</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>			<span class="comment">// statistic we can comfortably account for this space in. With this structure,</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>			<span class="comment">// we rely on demand paging to avoid large overheads, but tracking which memory</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>			<span class="comment">// is paged in is too expensive. Trying to account for the whole region means</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>			<span class="comment">// that it will appear like an enormous memory overhead in statistics, even though</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>			<span class="comment">// it is not.</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>			l2 = (*[1 &lt;&lt; arenaL2Bits]*heapArena)(sysAllocOS(unsafe.Sizeof(*l2)))
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>			if l2 == nil {
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>				throw(&#34;out of memory allocating heap arena map&#34;)
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>			}
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>			if h.arenasHugePages {
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>				sysHugePage(unsafe.Pointer(l2), unsafe.Sizeof(*l2))
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			} else {
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>				sysNoHugePage(unsafe.Pointer(l2), unsafe.Sizeof(*l2))
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>			}
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>			atomic.StorepNoWB(unsafe.Pointer(&amp;h.arenas[ri.l1()]), unsafe.Pointer(l2))
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>		}
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		if l2[ri.l2()] != nil {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			throw(&#34;arena already initialized&#34;)
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>		}
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		var r *heapArena
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>		r = (*heapArena)(h.heapArenaAlloc.alloc(unsafe.Sizeof(*r), goarch.PtrSize, &amp;memstats.gcMiscSys))
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		if r == nil {
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>			r = (*heapArena)(persistentalloc(unsafe.Sizeof(*r), goarch.PtrSize, &amp;memstats.gcMiscSys))
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>			if r == nil {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>				throw(&#34;out of memory allocating heap arena metadata&#34;)
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>			}
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		}
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>		<span class="comment">// Register the arena in allArenas if requested.</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>		if register {
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>			if len(h.allArenas) == cap(h.allArenas) {
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>				size := 2 * uintptr(cap(h.allArenas)) * goarch.PtrSize
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>				if size == 0 {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>					size = physPageSize
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>				}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>				newArray := (*notInHeap)(persistentalloc(size, goarch.PtrSize, &amp;memstats.gcMiscSys))
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>				if newArray == nil {
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>					throw(&#34;out of memory allocating allArenas&#34;)
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>				}
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>				oldSlice := h.allArenas
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>				*(*notInHeapSlice)(unsafe.Pointer(&amp;h.allArenas)) = notInHeapSlice{newArray, len(h.allArenas), int(size / goarch.PtrSize)}
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>				copy(h.allArenas, oldSlice)
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>				<span class="comment">// Do not free the old backing array because</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>				<span class="comment">// there may be concurrent readers. Since we</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>				<span class="comment">// double the array each time, this can lead</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>				<span class="comment">// to at most 2x waste.</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			}
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>			h.allArenas = h.allArenas[:len(h.allArenas)+1]
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>			h.allArenas[len(h.allArenas)-1] = ri
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		}
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>		<span class="comment">// Store atomically just in case an object from the</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		<span class="comment">// new heap arena becomes visible before the heap lock</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>		<span class="comment">// is released (which shouldn&#39;t happen, but there&#39;s</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		<span class="comment">// little downside to this).</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>		atomic.StorepNoWB(unsafe.Pointer(&amp;l2[ri.l2()]), unsafe.Pointer(r))
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	}
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	<span class="comment">// Tell the race detector about the new heap memory.</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	if raceenabled {
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		racemapshadow(v, size)
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	}
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	return
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>}
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span><span class="comment">// sysReserveAligned is like sysReserve, but the returned pointer is</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span><span class="comment">// aligned to align bytes. It may reserve either n or n+align bytes,</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span><span class="comment">// so it returns the size that was reserved.</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>func sysReserveAligned(v unsafe.Pointer, size, align uintptr) (unsafe.Pointer, uintptr) {
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	<span class="comment">// Since the alignment is rather large in uses of this</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	<span class="comment">// function, we&#39;re not likely to get it by chance, so we ask</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	<span class="comment">// for a larger region and remove the parts we don&#39;t need.</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	retries := 0
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>retry:
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	p := uintptr(sysReserve(v, size+align))
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	switch {
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	case p == 0:
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		return nil, 0
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	case p&amp;(align-1) == 0:
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		return unsafe.Pointer(p), size + align
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	case GOOS == &#34;windows&#34;:
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		<span class="comment">// On Windows we can&#39;t release pieces of a</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		<span class="comment">// reservation, so we release the whole thing and</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		<span class="comment">// re-reserve the aligned sub-region. This may race,</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		<span class="comment">// so we may have to try again.</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		sysFreeOS(unsafe.Pointer(p), size+align)
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		p = alignUp(p, align)
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		p2 := sysReserve(unsafe.Pointer(p), size)
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>		if p != uintptr(p2) {
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>			<span class="comment">// Must have raced. Try again.</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			sysFreeOS(p2, size)
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>			if retries++; retries == 100 {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>				throw(&#34;failed to allocate aligned heap memory; too many retries&#34;)
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>			}
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>			goto retry
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>		<span class="comment">// Success.</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		return p2, size
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	default:
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		<span class="comment">// Trim off the unaligned parts.</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>		pAligned := alignUp(p, align)
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		sysFreeOS(unsafe.Pointer(p), pAligned-p)
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>		end := pAligned + size
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		endLen := (p + size + align) - end
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		if endLen &gt; 0 {
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>			sysFreeOS(unsafe.Pointer(end), endLen)
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		}
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		return unsafe.Pointer(pAligned), size
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	}
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>}
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span><span class="comment">// enableMetadataHugePages enables huge pages for various sources of heap metadata.</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span><span class="comment">// A note on latency: for sufficiently small heaps (&lt;10s of GiB) this function will take constant</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span><span class="comment">// time, but may take time proportional to the size of the mapped heap beyond that.</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L871" class="ln">   871&nbsp;&nbsp;</span><span class="comment">// This function is idempotent.</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span><span class="comment">// The heap lock must not be held over this operation, since it will briefly acquire</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span><span class="comment">// the heap lock.</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span><span class="comment">// Must be called on the system stack because it acquires the heap lock.</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>func (h *mheap) enableMetadataHugePages() {
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	<span class="comment">// Enable huge pages for page structure.</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	h.pages.enableChunkHugePages()
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	<span class="comment">// Grab the lock and set arenasHugePages if it&#39;s not.</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	<span class="comment">// Once arenasHugePages is set, all new L2 entries will be eligible for</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	<span class="comment">// huge pages. We&#39;ll set all the old entries after we release the lock.</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	lock(&amp;h.lock)
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>	if h.arenasHugePages {
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>		unlock(&amp;h.lock)
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>		return
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>	}
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	h.arenasHugePages = true
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	unlock(&amp;h.lock)
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	<span class="comment">// N.B. The arenas L1 map is quite small on all platforms, so it&#39;s fine to</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	<span class="comment">// just iterate over the whole thing.</span>
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	for i := range h.arenas {
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		l2 := (*[1 &lt;&lt; arenaL2Bits]*heapArena)(atomic.Loadp(unsafe.Pointer(&amp;h.arenas[i])))
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>		if l2 == nil {
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>			continue
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		}
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>		sysHugePage(unsafe.Pointer(l2), unsafe.Sizeof(*l2))
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>	}
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>}
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span><span class="comment">// base address for all 0-byte allocations</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>var zerobase uintptr
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span><span class="comment">// nextFreeFast returns the next free object if one is quickly available.</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span><span class="comment">// Otherwise it returns 0.</span>
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>func nextFreeFast(s *mspan) gclinkptr {
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	theBit := sys.TrailingZeros64(s.allocCache) <span class="comment">// Is there a free object in the allocCache?</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	if theBit &lt; 64 {
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>		result := s.freeindex + uint16(theBit)
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		if result &lt; s.nelems {
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>			freeidx := result + 1
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>			if freeidx%64 == 0 &amp;&amp; freeidx != s.nelems {
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>				return 0
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>			}
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>			s.allocCache &gt;&gt;= uint(theBit + 1)
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>			s.freeindex = freeidx
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>			s.allocCount++
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>			return gclinkptr(uintptr(result)*s.elemsize + s.base())
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>		}
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	}
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	return 0
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>}
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span><span class="comment">// nextFree returns the next free object from the cached span if one is available.</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span><span class="comment">// Otherwise it refills the cache with a span with an available object and</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span><span class="comment">// returns that object along with a flag indicating that this was a heavy</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span><span class="comment">// weight allocation. If it is a heavy weight allocation the caller must</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span><span class="comment">// determine whether a new GC cycle needs to be started or if the GC is active</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span><span class="comment">// whether this goroutine needs to assist the GC.</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span><span class="comment">// Must run in a non-preemptible context since otherwise the owner of</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span><span class="comment">// c could change.</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>func (c *mcache) nextFree(spc spanClass) (v gclinkptr, s *mspan, shouldhelpgc bool) {
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	s = c.alloc[spc]
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	shouldhelpgc = false
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>	freeIndex := s.nextFreeIndex()
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>	if freeIndex == s.nelems {
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>		<span class="comment">// The span is full.</span>
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>		if s.allocCount != s.nelems {
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>			println(&#34;runtime: s.allocCount=&#34;, s.allocCount, &#34;s.nelems=&#34;, s.nelems)
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>			throw(&#34;s.allocCount != s.nelems &amp;&amp; freeIndex == s.nelems&#34;)
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>		}
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>		c.refill(spc)
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		shouldhelpgc = true
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>		s = c.alloc[spc]
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>		freeIndex = s.nextFreeIndex()
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>	}
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>	if freeIndex &gt;= s.nelems {
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>		throw(&#34;freeIndex is not valid&#34;)
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	}
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>	v = gclinkptr(uintptr(freeIndex)*s.elemsize + s.base())
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>	s.allocCount++
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	if s.allocCount &gt; s.nelems {
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>		println(&#34;s.allocCount=&#34;, s.allocCount, &#34;s.nelems=&#34;, s.nelems)
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>		throw(&#34;s.allocCount &gt; s.nelems&#34;)
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	}
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	return
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>}
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span><span class="comment">// Allocate an object of size bytes.</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span><span class="comment">// Small objects are allocated from the per-P cache&#39;s free lists.</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span><span class="comment">// Large objects (&gt; 32 kB) are allocated straight from the heap.</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>	if gcphase == _GCmarktermination {
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>		throw(&#34;mallocgc called with gcphase == _GCmarktermination&#34;)
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>	}
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	if size == 0 {
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>		return unsafe.Pointer(&amp;zerobase)
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	}
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s possible for any malloc to trigger sweeping, which may in</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	<span class="comment">// turn queue finalizers. Record this dynamic lock edge.</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	lockRankMayQueueFinalizer()
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>	userSize := size
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	if asanenabled {
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>		<span class="comment">// Refer to ASAN runtime library, the malloc() function allocates extra memory,</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>		<span class="comment">// the redzone, around the user requested memory region. And the redzones are marked</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>		<span class="comment">// as unaddressable. We perform the same operations in Go to detect the overflows or</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>		<span class="comment">// underflows.</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>		size += computeRZlog(size)
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>	}
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	if debug.malloc {
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>		if debug.sbrk != 0 {
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>			align := uintptr(16)
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>			if typ != nil {
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>				<span class="comment">// TODO(austin): This should be just</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>				<span class="comment">//   align = uintptr(typ.align)</span>
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>				<span class="comment">// but that&#39;s only 4 on 32-bit platforms,</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>				<span class="comment">// even if there&#39;s a uint64 field in typ (see #599).</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>				<span class="comment">// This causes 64-bit atomic accesses to panic.</span>
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>				<span class="comment">// Hence, we use stricter alignment that matches</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>				<span class="comment">// the normal allocator better.</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>				if size&amp;7 == 0 {
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>					align = 8
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>				} else if size&amp;3 == 0 {
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>					align = 4
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>				} else if size&amp;1 == 0 {
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>					align = 2
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>				} else {
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>					align = 1
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>				}
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>			}
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>			return persistentalloc(size, align, &amp;memstats.other_sys)
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>		}
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>		if inittrace.active &amp;&amp; inittrace.id == getg().goid {
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>			<span class="comment">// Init functions are executed sequentially in a single goroutine.</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>			inittrace.allocs += 1
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>		}
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	}
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	<span class="comment">// assistG is the G to charge for this allocation, or nil if</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>	<span class="comment">// GC is not currently active.</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>	assistG := deductAssistCredit(size)
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>	<span class="comment">// Set mp.mallocing to keep from being preempted by GC.</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	if mp.mallocing != 0 {
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>		throw(&#34;malloc deadlock&#34;)
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	}
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	if mp.gsignal == getg() {
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		throw(&#34;malloc during signal&#34;)
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>	}
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	mp.mallocing = 1
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>	shouldhelpgc := false
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	dataSize := userSize
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>	c := getMCache(mp)
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>	if c == nil {
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>		throw(&#34;mallocgc called without a P or outside bootstrapping&#34;)
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>	}
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>	var span *mspan
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>	var header **_type
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>	var x unsafe.Pointer
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>	noscan := typ == nil || typ.PtrBytes == 0
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>	<span class="comment">// In some cases block zeroing can profitably (for latency reduction purposes)</span>
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>	<span class="comment">// be delayed till preemption is possible; delayedZeroing tracks that state.</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>	delayedZeroing := false
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>	<span class="comment">// Determine if it&#39;s a &#39;small&#39; object that goes into a size-classed span.</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>	<span class="comment">// Note: This comparison looks a little strange, but it exists to smooth out</span>
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>	<span class="comment">// the crossover between the largest size class and large objects that have</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>	<span class="comment">// their own spans. The small window of object sizes between maxSmallSize-mallocHeaderSize</span>
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>	<span class="comment">// and maxSmallSize will be considered large, even though they might fit in</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>	<span class="comment">// a size class. In practice this is completely fine, since the largest small</span>
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>	<span class="comment">// size class has a single object in it already, precisely to make the transition</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>	<span class="comment">// to large objects smooth.</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>	if size &lt;= maxSmallSize-mallocHeaderSize {
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>		if noscan &amp;&amp; size &lt; maxTinySize {
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>			<span class="comment">// Tiny allocator.</span>
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>			<span class="comment">// Tiny allocator combines several tiny allocation requests</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>			<span class="comment">// into a single memory block. The resulting memory block</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>			<span class="comment">// is freed when all subobjects are unreachable. The subobjects</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>			<span class="comment">// must be noscan (don&#39;t have pointers), this ensures that</span>
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>			<span class="comment">// the amount of potentially wasted memory is bounded.</span>
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>			<span class="comment">// Size of the memory block used for combining (maxTinySize) is tunable.</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>			<span class="comment">// Current setting is 16 bytes, which relates to 2x worst case memory</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>			<span class="comment">// wastage (when all but one subobjects are unreachable).</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>			<span class="comment">// 8 bytes would result in no wastage at all, but provides less</span>
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>			<span class="comment">// opportunities for combining.</span>
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>			<span class="comment">// 32 bytes provides more opportunities for combining,</span>
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>			<span class="comment">// but can lead to 4x worst case wastage.</span>
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>			<span class="comment">// The best case winning is 8x regardless of block size.</span>
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>			<span class="comment">// Objects obtained from tiny allocator must not be freed explicitly.</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>			<span class="comment">// So when an object will be freed explicitly, we ensure that</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>			<span class="comment">// its size &gt;= maxTinySize.</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>			<span class="comment">// SetFinalizer has a special case for objects potentially coming</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>			<span class="comment">// from tiny allocator, it such case it allows to set finalizers</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>			<span class="comment">// for an inner byte of a memory block.</span>
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>			<span class="comment">// The main targets of tiny allocator are small strings and</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>			<span class="comment">// standalone escaping variables. On a json benchmark</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>			<span class="comment">// the allocator reduces number of allocations by ~12% and</span>
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>			<span class="comment">// reduces heap size by ~20%.</span>
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>			off := c.tinyoffset
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>			<span class="comment">// Align tiny pointer for required (conservative) alignment.</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>			if size&amp;7 == 0 {
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>				off = alignUp(off, 8)
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>			} else if goarch.PtrSize == 4 &amp;&amp; size == 12 {
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>				<span class="comment">// Conservatively align 12-byte objects to 8 bytes on 32-bit</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>				<span class="comment">// systems so that objects whose first field is a 64-bit</span>
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>				<span class="comment">// value is aligned to 8 bytes and does not cause a fault on</span>
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>				<span class="comment">// atomic access. See issue 37262.</span>
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>				<span class="comment">// TODO(mknyszek): Remove this workaround if/when issue 36606</span>
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>				<span class="comment">// is resolved.</span>
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>				off = alignUp(off, 8)
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>			} else if size&amp;3 == 0 {
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>				off = alignUp(off, 4)
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>			} else if size&amp;1 == 0 {
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>				off = alignUp(off, 2)
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>			}
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>			if off+size &lt;= maxTinySize &amp;&amp; c.tiny != 0 {
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>				<span class="comment">// The object fits into existing tiny block.</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>				x = unsafe.Pointer(c.tiny + off)
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>				c.tinyoffset = off + size
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>				c.tinyAllocs++
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>				mp.mallocing = 0
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>				releasem(mp)
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>				return x
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>			}
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>			<span class="comment">// Allocate a new maxTinySize block.</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>			span = c.alloc[tinySpanClass]
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>			v := nextFreeFast(span)
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>			if v == 0 {
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>				v, span, shouldhelpgc = c.nextFree(tinySpanClass)
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>			}
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>			x = unsafe.Pointer(v)
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>			(*[2]uint64)(x)[0] = 0
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>			(*[2]uint64)(x)[1] = 0
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>			<span class="comment">// See if we need to replace the existing tiny block with the new one</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>			<span class="comment">// based on amount of remaining free space.</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>			if !raceenabled &amp;&amp; (size &lt; c.tinyoffset || c.tiny == 0) {
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>				<span class="comment">// Note: disabled when race detector is on, see comment near end of this function.</span>
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>				c.tiny = uintptr(x)
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>				c.tinyoffset = size
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>			}
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>			size = maxTinySize
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>		} else {
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>			hasHeader := !noscan &amp;&amp; !heapBitsInSpan(size)
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>			if goexperiment.AllocHeaders &amp;&amp; hasHeader {
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>				size += mallocHeaderSize
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>			}
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>			var sizeclass uint8
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>			if size &lt;= smallSizeMax-8 {
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>				sizeclass = size_to_class8[divRoundUp(size, smallSizeDiv)]
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>			} else {
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>				sizeclass = size_to_class128[divRoundUp(size-smallSizeMax, largeSizeDiv)]
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>			}
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>			size = uintptr(class_to_size[sizeclass])
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>			spc := makeSpanClass(sizeclass, noscan)
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>			span = c.alloc[spc]
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>			v := nextFreeFast(span)
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>			if v == 0 {
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>				v, span, shouldhelpgc = c.nextFree(spc)
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>			}
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>			x = unsafe.Pointer(v)
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>			if needzero &amp;&amp; span.needzero != 0 {
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>				memclrNoHeapPointers(x, size)
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>			}
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>			if goexperiment.AllocHeaders &amp;&amp; hasHeader {
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>				header = (**_type)(x)
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>				x = add(x, mallocHeaderSize)
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>				size -= mallocHeaderSize
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>			}
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>		}
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>	} else {
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>		shouldhelpgc = true
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>		<span class="comment">// For large allocations, keep track of zeroed state so that</span>
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>		<span class="comment">// bulk zeroing can be happen later in a preemptible context.</span>
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>		span = c.allocLarge(size, noscan)
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>		span.freeindex = 1
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>		span.allocCount = 1
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>		size = span.elemsize
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>		x = unsafe.Pointer(span.base())
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>		if needzero &amp;&amp; span.needzero != 0 {
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>			if noscan {
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>				delayedZeroing = true
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>			} else {
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>				memclrNoHeapPointers(x, size)
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>			}
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>		}
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>		if goexperiment.AllocHeaders &amp;&amp; !noscan {
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>			header = &amp;span.largeType
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>		}
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>	}
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>	if !noscan {
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>		if goexperiment.AllocHeaders {
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>			c.scanAlloc += heapSetType(uintptr(x), dataSize, typ, header, span)
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>		} else {
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>			var scanSize uintptr
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>			heapBitsSetType(uintptr(x), size, dataSize, typ)
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>			if dataSize &gt; typ.Size_ {
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>				<span class="comment">// Array allocation. If there are any</span>
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>				<span class="comment">// pointers, GC has to scan to the last</span>
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>				<span class="comment">// element.</span>
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>				if typ.PtrBytes != 0 {
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>					scanSize = dataSize - typ.Size_ + typ.PtrBytes
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>				}
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>			} else {
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>				scanSize = typ.PtrBytes
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>			}
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>			c.scanAlloc += scanSize
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>		}
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>	}
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>	<span class="comment">// Ensure that the stores above that initialize x to</span>
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>	<span class="comment">// type-safe memory and set the heap bits occur before</span>
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>	<span class="comment">// the caller can make x observable to the garbage</span>
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>	<span class="comment">// collector. Otherwise, on weakly ordered machines,</span>
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>	<span class="comment">// the garbage collector could follow a pointer to x,</span>
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>	<span class="comment">// but see uninitialized memory or stale heap bits.</span>
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>	publicationBarrier()
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>	<span class="comment">// As x and the heap bits are initialized, update</span>
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>	<span class="comment">// freeIndexForScan now so x is seen by the GC</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>	<span class="comment">// (including conservative scan) as an allocated object.</span>
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>	<span class="comment">// While this pointer can&#39;t escape into user code as a</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>	<span class="comment">// _live_ pointer until we return, conservative scanning</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>	<span class="comment">// may find a dead pointer that happens to point into this</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>	<span class="comment">// object. Delaying this update until now ensures that</span>
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>	<span class="comment">// conservative scanning considers this pointer dead until</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>	<span class="comment">// this point.</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>	span.freeIndexForScan = span.freeindex
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>	<span class="comment">// Allocate black during GC.</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>	<span class="comment">// All slots hold nil so no scanning is needed.</span>
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>	<span class="comment">// This may be racing with GC so do it atomically if there can be</span>
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>	<span class="comment">// a race marking the bit.</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>	if gcphase != _GCoff {
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>		gcmarknewobject(span, uintptr(x))
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>	}
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>	if raceenabled {
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>		racemalloc(x, size)
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>	}
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	if msanenabled {
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>		msanmalloc(x, size)
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>	}
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>	if asanenabled {
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>		<span class="comment">// We should only read/write the memory with the size asked by the user.</span>
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>		<span class="comment">// The rest of the allocated memory should be poisoned, so that we can report</span>
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>		<span class="comment">// errors when accessing poisoned memory.</span>
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>		<span class="comment">// The allocated memory is larger than required userSize, it will also include</span>
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>		<span class="comment">// redzone and some other padding bytes.</span>
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>		rzBeg := unsafe.Add(x, userSize)
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>		asanpoison(rzBeg, size-userSize)
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>		asanunpoison(x, userSize)
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>	}
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>	<span class="comment">// If !goexperiment.AllocHeaders, &#34;size&#34; doesn&#39;t include the</span>
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>	<span class="comment">// allocation header, so use span.elemsize as the &#34;full&#34; size</span>
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>	<span class="comment">// for various computations below.</span>
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): We should really count the header as part</span>
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>	<span class="comment">// of gc_sys or something, but it&#39;s risky to change the</span>
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>	<span class="comment">// accounting so much right now. Just pretend its internal</span>
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>	<span class="comment">// fragmentation and match the GC&#39;s accounting by using the</span>
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>	<span class="comment">// whole allocation slot.</span>
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>	fullSize := size
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>	if goexperiment.AllocHeaders {
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>		fullSize = span.elemsize
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>	}
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>	if rate := MemProfileRate; rate &gt; 0 {
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>		<span class="comment">// Note cache c only valid while m acquired; see #47302</span>
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>		<span class="comment">// N.B. Use the full size because that matches how the GC</span>
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>		<span class="comment">// will update the mem profile on the &#34;free&#34; side.</span>
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>		if rate != 1 &amp;&amp; fullSize &lt; c.nextSample {
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>			c.nextSample -= fullSize
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>		} else {
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>			profilealloc(mp, x, fullSize)
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>		}
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>	}
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	mp.mallocing = 0
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	releasem(mp)
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	<span class="comment">// Pointerfree data can be zeroed late in a context where preemption can occur.</span>
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>	<span class="comment">// x will keep the memory alive.</span>
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>	if delayedZeroing {
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>		if !noscan {
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>			throw(&#34;delayed zeroing on data that may contain pointers&#34;)
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>		}
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>		if goexperiment.AllocHeaders &amp;&amp; header != nil {
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>			throw(&#34;unexpected malloc header in delayed zeroing of large object&#34;)
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>		}
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>		<span class="comment">// N.B. size == fullSize always in this case.</span>
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>		memclrNoHeapPointersChunked(size, x) <span class="comment">// This is a possible preemption point: see #47302</span>
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>	}
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>	if debug.malloc {
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>		if debug.allocfreetrace != 0 {
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>			tracealloc(x, size, typ)
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>		}
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>		if inittrace.active &amp;&amp; inittrace.id == getg().goid {
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>			<span class="comment">// Init functions are executed sequentially in a single goroutine.</span>
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>			inittrace.bytes += uint64(fullSize)
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>		}
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>	}
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>	if assistG != nil {
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>		<span class="comment">// Account for internal fragmentation in the assist</span>
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>		<span class="comment">// debt now that we know it.</span>
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>		<span class="comment">// N.B. Use the full size because that&#39;s how the rest</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>		<span class="comment">// of the GC accounts for bytes marked.</span>
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>		assistG.gcAssistBytes -= int64(fullSize - dataSize)
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>	}
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>	if shouldhelpgc {
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>		if t := (gcTrigger{kind: gcTriggerHeap}); t.test() {
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>			gcStart(t)
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>		}
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>	}
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>	if raceenabled &amp;&amp; noscan &amp;&amp; dataSize &lt; maxTinySize {
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>		<span class="comment">// Pad tinysize allocations so they are aligned with the end</span>
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>		<span class="comment">// of the tinyalloc region. This ensures that any arithmetic</span>
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>		<span class="comment">// that goes off the top end of the object will be detectable</span>
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>		<span class="comment">// by checkptr (issue 38872).</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>		<span class="comment">// Note that we disable tinyalloc when raceenabled for this to work.</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>		<span class="comment">// TODO: This padding is only performed when the race detector</span>
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>		<span class="comment">// is enabled. It would be nice to enable it if any package</span>
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>		<span class="comment">// was compiled with checkptr, but there&#39;s no easy way to</span>
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>		<span class="comment">// detect that (especially at compile time).</span>
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>		<span class="comment">// TODO: enable this padding for all allocations, not just</span>
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>		<span class="comment">// tinyalloc ones. It&#39;s tricky because of pointer maps.</span>
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>		<span class="comment">// Maybe just all noscan objects?</span>
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>		x = add(x, size-dataSize)
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>	}
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>	return x
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>}
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span><span class="comment">// deductAssistCredit reduces the current G&#39;s assist credit</span>
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span><span class="comment">// by size bytes, and assists the GC if necessary.</span>
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span><span class="comment">// Caller must be preemptible.</span>
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span><span class="comment">// Returns the G for which the assist credit was accounted.</span>
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>func deductAssistCredit(size uintptr) *g {
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>	var assistG *g
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>	if gcBlackenEnabled != 0 {
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>		<span class="comment">// Charge the current user G for this allocation.</span>
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>		assistG = getg()
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>		if assistG.m.curg != nil {
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>			assistG = assistG.m.curg
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>		}
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>		<span class="comment">// Charge the allocation against the G. We&#39;ll account</span>
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>		<span class="comment">// for internal fragmentation at the end of mallocgc.</span>
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>		assistG.gcAssistBytes -= int64(size)
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>		if assistG.gcAssistBytes &lt; 0 {
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>			<span class="comment">// This G is in debt. Assist the GC to correct</span>
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>			<span class="comment">// this before allocating. This must happen</span>
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>			<span class="comment">// before disabling preemption.</span>
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>			gcAssistAlloc(assistG)
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>		}
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>	}
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>	return assistG
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>}
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span><span class="comment">// memclrNoHeapPointersChunked repeatedly calls memclrNoHeapPointers</span>
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span><span class="comment">// on chunks of the buffer to be zeroed, with opportunities for preemption</span>
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span><span class="comment">// along the way.  memclrNoHeapPointers contains no safepoints and also</span>
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span><span class="comment">// cannot be preemptively scheduled, so this provides a still-efficient</span>
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span><span class="comment">// block copy that can also be preempted on a reasonable granularity.</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span><span class="comment">// Use this with care; if the data being cleared is tagged to contain</span>
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span><span class="comment">// pointers, this allows the GC to run before it is all cleared.</span>
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>func memclrNoHeapPointersChunked(size uintptr, x unsafe.Pointer) {
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>	v := uintptr(x)
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>	<span class="comment">// got this from benchmarking. 128k is too small, 512k is too large.</span>
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>	const chunkBytes = 256 * 1024
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>	vsize := v + size
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>	for voff := v; voff &lt; vsize; voff = voff + chunkBytes {
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>		if getg().preempt {
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>			<span class="comment">// may hold locks, e.g., profiling</span>
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>			goschedguarded()
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>		}
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>		<span class="comment">// clear min(avail, lump) bytes</span>
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>		n := vsize - voff
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>		if n &gt; chunkBytes {
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>			n = chunkBytes
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>		}
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>		memclrNoHeapPointers(unsafe.Pointer(voff), n)
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>	}
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>}
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span><span class="comment">// implementation of new builtin</span>
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span><span class="comment">// compiler (both frontend and SSA backend) knows the signature</span>
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span><span class="comment">// of this function.</span>
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>func newobject(typ *_type) unsafe.Pointer {
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>	return mallocgc(typ.Size_, typ, true)
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>}
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_unsafe_New reflect.unsafe_New</span>
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>func reflect_unsafe_New(typ *_type) unsafe.Pointer {
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>	return mallocgc(typ.Size_, typ, true)
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>}
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span><span class="comment">//go:linkname reflectlite_unsafe_New internal/reflectlite.unsafe_New</span>
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>func reflectlite_unsafe_New(typ *_type) unsafe.Pointer {
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>	return mallocgc(typ.Size_, typ, true)
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>}
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span><span class="comment">// newarray allocates an array of n elements of type typ.</span>
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>func newarray(typ *_type, n int) unsafe.Pointer {
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>	if n == 1 {
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>		return mallocgc(typ.Size_, typ, true)
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>	}
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>	mem, overflow := math.MulUintptr(typ.Size_, uintptr(n))
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>	if overflow || mem &gt; maxAlloc || n &lt; 0 {
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>		panic(plainError(&#34;runtime: allocation size out of range&#34;))
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>	}
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>	return mallocgc(mem, typ, true)
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>}
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_unsafe_NewArray reflect.unsafe_NewArray</span>
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>func reflect_unsafe_NewArray(typ *_type, n int) unsafe.Pointer {
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>	return newarray(typ, n)
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>}
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>func profilealloc(mp *m, x unsafe.Pointer, size uintptr) {
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>	c := getMCache(mp)
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>	if c == nil {
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>		throw(&#34;profilealloc called without a P or outside bootstrapping&#34;)
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>	}
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>	c.nextSample = nextSample()
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>	mProf_Malloc(x, size)
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>}
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span><span class="comment">// nextSample returns the next sampling point for heap profiling. The goal is</span>
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span><span class="comment">// to sample allocations on average every MemProfileRate bytes, but with a</span>
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span><span class="comment">// completely random distribution over the allocation timeline; this</span>
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span><span class="comment">// corresponds to a Poisson process with parameter MemProfileRate. In Poisson</span>
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span><span class="comment">// processes, the distance between two samples follows the exponential</span>
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span><span class="comment">// distribution (exp(MemProfileRate)), so the best return value is a random</span>
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span><span class="comment">// number taken from an exponential distribution whose mean is MemProfileRate.</span>
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span>func nextSample() uintptr {
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span>	if MemProfileRate == 1 {
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span>		<span class="comment">// Callers assign our return value to</span>
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span>		<span class="comment">// mcache.next_sample, but next_sample is not used</span>
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span>		<span class="comment">// when the rate is 1. So avoid the math below and</span>
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>		<span class="comment">// just return something.</span>
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>		return 0
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>	}
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>	if GOOS == &#34;plan9&#34; {
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>		<span class="comment">// Plan 9 doesn&#39;t support floating point in note handler.</span>
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>		if gp := getg(); gp == gp.m.gsignal {
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>			return nextSampleNoFP()
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>		}
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>	}
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span>
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span>	return uintptr(fastexprand(MemProfileRate))
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span>}
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span>
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span><span class="comment">// fastexprand returns a random number from an exponential distribution with</span>
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span><span class="comment">// the specified mean.</span>
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span>func fastexprand(mean int) int32 {
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span>	<span class="comment">// Avoid overflow. Maximum possible step is</span>
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span>	<span class="comment">// -ln(1/(1&lt;&lt;randomBitCount)) * mean, approximately 20 * mean.</span>
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span>	switch {
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span>	case mean &gt; 0x7000000:
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span>		mean = 0x7000000
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span>	case mean == 0:
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span>		return 0
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span>	}
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span>
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span>	<span class="comment">// Take a random sample of the exponential distribution exp(-mean*x).</span>
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span>	<span class="comment">// The probability distribution function is mean*exp(-mean*x), so the CDF is</span>
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span>	<span class="comment">// p = 1 - exp(-mean*x), so</span>
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span>	<span class="comment">// q = 1 - p == exp(-mean*x)</span>
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span>	<span class="comment">// log_e(q) = -mean*x</span>
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span>	<span class="comment">// -log_e(q)/mean = x</span>
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span>	<span class="comment">// x = -log_e(q) * mean</span>
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span>	<span class="comment">// x = log_2(q) * (-log_e(2)) * mean    ; Using log_2 for efficiency</span>
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span>	const randomBitCount = 26
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span>	q := cheaprandn(1&lt;&lt;randomBitCount) + 1
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span>	qlog := fastlog2(float64(q)) - randomBitCount
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span>	if qlog &gt; 0 {
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span>		qlog = 0
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span>	}
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span>	const minusLog2 = -0.6931471805599453 <span class="comment">// -ln(2)</span>
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span>	return int32(qlog*(minusLog2*float64(mean))) + 1
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span>}
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span>
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span><span class="comment">// nextSampleNoFP is similar to nextSample, but uses older,</span>
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span><span class="comment">// simpler code to avoid floating point.</span>
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span>func nextSampleNoFP() uintptr {
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span>	<span class="comment">// Set first allocation sample size.</span>
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span>	rate := MemProfileRate
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span>	if rate &gt; 0x3fffffff { <span class="comment">// make 2*rate not overflow</span>
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span>		rate = 0x3fffffff
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span>	}
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span>	if rate != 0 {
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span>		return uintptr(cheaprandn(uint32(2 * rate)))
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span>	}
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span>	return 0
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span>}
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span>
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span>type persistentAlloc struct {
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span>	base *notInHeap
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span>	off  uintptr
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span>}
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span>
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span>var globalAlloc struct {
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span>	mutex
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span>	persistentAlloc
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span>}
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span>
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span><span class="comment">// persistentChunkSize is the number of bytes we allocate when we grow</span>
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span><span class="comment">// a persistentAlloc.</span>
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span>const persistentChunkSize = 256 &lt;&lt; 10
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span>
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span><span class="comment">// persistentChunks is a list of all the persistent chunks we have</span>
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span><span class="comment">// allocated. The list is maintained through the first word in the</span>
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span><span class="comment">// persistent chunk. This is updated atomically.</span>
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>var persistentChunks *notInHeap
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span><span class="comment">// Wrapper around sysAlloc that can allocate small chunks.</span>
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span><span class="comment">// There is no associated free operation.</span>
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span><span class="comment">// Intended for things like function/type/debug-related persistent data.</span>
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span><span class="comment">// If align is 0, uses default align (currently 8).</span>
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span><span class="comment">// The returned memory will be zeroed.</span>
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span><span class="comment">// sysStat must be non-nil.</span>
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span><span class="comment">// Consider marking persistentalloc&#39;d types not in heap by embedding</span>
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span><span class="comment">// runtime/internal/sys.NotInHeap.</span>
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span>func persistentalloc(size, align uintptr, sysStat *sysMemStat) unsafe.Pointer {
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span>	var p *notInHeap
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span>		p = persistentalloc1(size, align, sysStat)
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>	})
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>	return unsafe.Pointer(p)
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>}
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span><span class="comment">// Must run on system stack because stack growth can (re)invoke it.</span>
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span><span class="comment">// See issue 9174.</span>
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>func persistentalloc1(size, align uintptr, sysStat *sysMemStat) *notInHeap {
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>	const (
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>		maxBlock = 64 &lt;&lt; 10 <span class="comment">// VM reservation granularity is 64K on windows</span>
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span>	)
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span>
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span>	if size == 0 {
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span>		throw(&#34;persistentalloc: size == 0&#34;)
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>	}
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>	if align != 0 {
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>		if align&amp;(align-1) != 0 {
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>			throw(&#34;persistentalloc: align is not a power of 2&#34;)
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span>		}
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span>		if align &gt; _PageSize {
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span>			throw(&#34;persistentalloc: align is too large&#34;)
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span>		}
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span>	} else {
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span>		align = 8
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span>	}
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span>
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span>	if size &gt;= maxBlock {
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span>		return (*notInHeap)(sysAlloc(size, sysStat))
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span>	}
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span>
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span>	var persistent *persistentAlloc
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span>	if mp != nil &amp;&amp; mp.p != 0 {
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span>		persistent = &amp;mp.p.ptr().palloc
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span>	} else {
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span>		lock(&amp;globalAlloc.mutex)
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span>		persistent = &amp;globalAlloc.persistentAlloc
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span>	}
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span>	persistent.off = alignUp(persistent.off, align)
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span>	if persistent.off+size &gt; persistentChunkSize || persistent.base == nil {
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span>		persistent.base = (*notInHeap)(sysAlloc(persistentChunkSize, &amp;memstats.other_sys))
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span>		if persistent.base == nil {
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span>			if persistent == &amp;globalAlloc.persistentAlloc {
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span>				unlock(&amp;globalAlloc.mutex)
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span>			}
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span>			throw(&#34;runtime: cannot allocate memory&#34;)
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span>		}
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span>
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span>		<span class="comment">// Add the new chunk to the persistentChunks list.</span>
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span>		for {
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span>			chunks := uintptr(unsafe.Pointer(persistentChunks))
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span>			*(*uintptr)(unsafe.Pointer(persistent.base)) = chunks
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>			if atomic.Casuintptr((*uintptr)(unsafe.Pointer(&amp;persistentChunks)), chunks, uintptr(unsafe.Pointer(persistent.base))) {
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span>				break
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span>			}
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span>		}
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span>		persistent.off = alignUp(goarch.PtrSize, align)
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span>	}
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>	p := persistent.base.add(persistent.off)
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>	persistent.off += size
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>	releasem(mp)
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>	if persistent == &amp;globalAlloc.persistentAlloc {
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>		unlock(&amp;globalAlloc.mutex)
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>	}
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span>	if sysStat != &amp;memstats.other_sys {
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span>		sysStat.add(int64(size))
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span>		memstats.other_sys.add(-int64(size))
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span>	}
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span>	return p
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>}
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span><span class="comment">// inPersistentAlloc reports whether p points to memory allocated by</span>
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span><span class="comment">// persistentalloc. This must be nosplit because it is called by the</span>
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span><span class="comment">// cgo checker code, which is called by the write barrier code.</span>
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span>func inPersistentAlloc(p uintptr) bool {
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span>	chunk := atomic.Loaduintptr((*uintptr)(unsafe.Pointer(&amp;persistentChunks)))
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span>	for chunk != 0 {
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span>		if p &gt;= chunk &amp;&amp; p &lt; chunk+persistentChunkSize {
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span>			return true
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span>		}
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span>		chunk = *(*uintptr)(unsafe.Pointer(chunk))
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>	}
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>	return false
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>}
<span id="L1618" class="ln">  1618&nbsp;&nbsp;</span>
<span id="L1619" class="ln">  1619&nbsp;&nbsp;</span><span class="comment">// linearAlloc is a simple linear allocator that pre-reserves a region</span>
<span id="L1620" class="ln">  1620&nbsp;&nbsp;</span><span class="comment">// of memory and then optionally maps that region into the Ready state</span>
<span id="L1621" class="ln">  1621&nbsp;&nbsp;</span><span class="comment">// as needed.</span>
<span id="L1622" class="ln">  1622&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1623" class="ln">  1623&nbsp;&nbsp;</span><span class="comment">// The caller is responsible for locking.</span>
<span id="L1624" class="ln">  1624&nbsp;&nbsp;</span>type linearAlloc struct {
<span id="L1625" class="ln">  1625&nbsp;&nbsp;</span>	next   uintptr <span class="comment">// next free byte</span>
<span id="L1626" class="ln">  1626&nbsp;&nbsp;</span>	mapped uintptr <span class="comment">// one byte past end of mapped space</span>
<span id="L1627" class="ln">  1627&nbsp;&nbsp;</span>	end    uintptr <span class="comment">// end of reserved space</span>
<span id="L1628" class="ln">  1628&nbsp;&nbsp;</span>
<span id="L1629" class="ln">  1629&nbsp;&nbsp;</span>	mapMemory bool <span class="comment">// transition memory from Reserved to Ready if true</span>
<span id="L1630" class="ln">  1630&nbsp;&nbsp;</span>}
<span id="L1631" class="ln">  1631&nbsp;&nbsp;</span>
<span id="L1632" class="ln">  1632&nbsp;&nbsp;</span>func (l *linearAlloc) init(base, size uintptr, mapMemory bool) {
<span id="L1633" class="ln">  1633&nbsp;&nbsp;</span>	if base+size &lt; base {
<span id="L1634" class="ln">  1634&nbsp;&nbsp;</span>		<span class="comment">// Chop off the last byte. The runtime isn&#39;t prepared</span>
<span id="L1635" class="ln">  1635&nbsp;&nbsp;</span>		<span class="comment">// to deal with situations where the bounds could overflow.</span>
<span id="L1636" class="ln">  1636&nbsp;&nbsp;</span>		<span class="comment">// Leave that memory reserved, though, so we don&#39;t map it</span>
<span id="L1637" class="ln">  1637&nbsp;&nbsp;</span>		<span class="comment">// later.</span>
<span id="L1638" class="ln">  1638&nbsp;&nbsp;</span>		size -= 1
<span id="L1639" class="ln">  1639&nbsp;&nbsp;</span>	}
<span id="L1640" class="ln">  1640&nbsp;&nbsp;</span>	l.next, l.mapped = base, base
<span id="L1641" class="ln">  1641&nbsp;&nbsp;</span>	l.end = base + size
<span id="L1642" class="ln">  1642&nbsp;&nbsp;</span>	l.mapMemory = mapMemory
<span id="L1643" class="ln">  1643&nbsp;&nbsp;</span>}
<span id="L1644" class="ln">  1644&nbsp;&nbsp;</span>
<span id="L1645" class="ln">  1645&nbsp;&nbsp;</span>func (l *linearAlloc) alloc(size, align uintptr, sysStat *sysMemStat) unsafe.Pointer {
<span id="L1646" class="ln">  1646&nbsp;&nbsp;</span>	p := alignUp(l.next, align)
<span id="L1647" class="ln">  1647&nbsp;&nbsp;</span>	if p+size &gt; l.end {
<span id="L1648" class="ln">  1648&nbsp;&nbsp;</span>		return nil
<span id="L1649" class="ln">  1649&nbsp;&nbsp;</span>	}
<span id="L1650" class="ln">  1650&nbsp;&nbsp;</span>	l.next = p + size
<span id="L1651" class="ln">  1651&nbsp;&nbsp;</span>	if pEnd := alignUp(l.next-1, physPageSize); pEnd &gt; l.mapped {
<span id="L1652" class="ln">  1652&nbsp;&nbsp;</span>		if l.mapMemory {
<span id="L1653" class="ln">  1653&nbsp;&nbsp;</span>			<span class="comment">// Transition from Reserved to Prepared to Ready.</span>
<span id="L1654" class="ln">  1654&nbsp;&nbsp;</span>			n := pEnd - l.mapped
<span id="L1655" class="ln">  1655&nbsp;&nbsp;</span>			sysMap(unsafe.Pointer(l.mapped), n, sysStat)
<span id="L1656" class="ln">  1656&nbsp;&nbsp;</span>			sysUsed(unsafe.Pointer(l.mapped), n, n)
<span id="L1657" class="ln">  1657&nbsp;&nbsp;</span>		}
<span id="L1658" class="ln">  1658&nbsp;&nbsp;</span>		l.mapped = pEnd
<span id="L1659" class="ln">  1659&nbsp;&nbsp;</span>	}
<span id="L1660" class="ln">  1660&nbsp;&nbsp;</span>	return unsafe.Pointer(p)
<span id="L1661" class="ln">  1661&nbsp;&nbsp;</span>}
<span id="L1662" class="ln">  1662&nbsp;&nbsp;</span>
<span id="L1663" class="ln">  1663&nbsp;&nbsp;</span><span class="comment">// notInHeap is off-heap memory allocated by a lower-level allocator</span>
<span id="L1664" class="ln">  1664&nbsp;&nbsp;</span><span class="comment">// like sysAlloc or persistentAlloc.</span>
<span id="L1665" class="ln">  1665&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1666" class="ln">  1666&nbsp;&nbsp;</span><span class="comment">// In general, it&#39;s better to use real types which embed</span>
<span id="L1667" class="ln">  1667&nbsp;&nbsp;</span><span class="comment">// runtime/internal/sys.NotInHeap, but this serves as a generic type</span>
<span id="L1668" class="ln">  1668&nbsp;&nbsp;</span><span class="comment">// for situations where that isn&#39;t possible (like in the allocators).</span>
<span id="L1669" class="ln">  1669&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1670" class="ln">  1670&nbsp;&nbsp;</span><span class="comment">// TODO: Use this as the return type of sysAlloc, persistentAlloc, etc?</span>
<span id="L1671" class="ln">  1671&nbsp;&nbsp;</span>type notInHeap struct{ _ sys.NotInHeap }
<span id="L1672" class="ln">  1672&nbsp;&nbsp;</span>
<span id="L1673" class="ln">  1673&nbsp;&nbsp;</span>func (p *notInHeap) add(bytes uintptr) *notInHeap {
<span id="L1674" class="ln">  1674&nbsp;&nbsp;</span>	return (*notInHeap)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + bytes))
<span id="L1675" class="ln">  1675&nbsp;&nbsp;</span>}
<span id="L1676" class="ln">  1676&nbsp;&nbsp;</span>
<span id="L1677" class="ln">  1677&nbsp;&nbsp;</span><span class="comment">// computeRZlog computes the size of the redzone.</span>
<span id="L1678" class="ln">  1678&nbsp;&nbsp;</span><span class="comment">// Refer to the implementation of the compiler-rt.</span>
<span id="L1679" class="ln">  1679&nbsp;&nbsp;</span>func computeRZlog(userSize uintptr) uintptr {
<span id="L1680" class="ln">  1680&nbsp;&nbsp;</span>	switch {
<span id="L1681" class="ln">  1681&nbsp;&nbsp;</span>	case userSize &lt;= (64 - 16):
<span id="L1682" class="ln">  1682&nbsp;&nbsp;</span>		return 16 &lt;&lt; 0
<span id="L1683" class="ln">  1683&nbsp;&nbsp;</span>	case userSize &lt;= (128 - 32):
<span id="L1684" class="ln">  1684&nbsp;&nbsp;</span>		return 16 &lt;&lt; 1
<span id="L1685" class="ln">  1685&nbsp;&nbsp;</span>	case userSize &lt;= (512 - 64):
<span id="L1686" class="ln">  1686&nbsp;&nbsp;</span>		return 16 &lt;&lt; 2
<span id="L1687" class="ln">  1687&nbsp;&nbsp;</span>	case userSize &lt;= (4096 - 128):
<span id="L1688" class="ln">  1688&nbsp;&nbsp;</span>		return 16 &lt;&lt; 3
<span id="L1689" class="ln">  1689&nbsp;&nbsp;</span>	case userSize &lt;= (1&lt;&lt;14)-256:
<span id="L1690" class="ln">  1690&nbsp;&nbsp;</span>		return 16 &lt;&lt; 4
<span id="L1691" class="ln">  1691&nbsp;&nbsp;</span>	case userSize &lt;= (1&lt;&lt;15)-512:
<span id="L1692" class="ln">  1692&nbsp;&nbsp;</span>		return 16 &lt;&lt; 5
<span id="L1693" class="ln">  1693&nbsp;&nbsp;</span>	case userSize &lt;= (1&lt;&lt;16)-1024:
<span id="L1694" class="ln">  1694&nbsp;&nbsp;</span>		return 16 &lt;&lt; 6
<span id="L1695" class="ln">  1695&nbsp;&nbsp;</span>	default:
<span id="L1696" class="ln">  1696&nbsp;&nbsp;</span>		return 16 &lt;&lt; 7
<span id="L1697" class="ln">  1697&nbsp;&nbsp;</span>	}
<span id="L1698" class="ln">  1698&nbsp;&nbsp;</span>}
<span id="L1699" class="ln">  1699&nbsp;&nbsp;</span>
</pre><p><a href="malloc.go?m=text">View as plain text</a></p>

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
