<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mem.go - Go Documentation Server</title>

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
<a href="mem.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mem.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2022 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;unsafe&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// OS memory management abstraction layer</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// Regions of the address space managed by the runtime may be in one of four</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// states at any given time:</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// 1) None - Unreserved and unmapped, the default state of any region.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// 2) Reserved - Owned by the runtime, but accessing it would cause a fault.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//               Does not count against the process&#39; memory footprint.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// 3) Prepared - Reserved, intended not to be backed by physical memory (though</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//               an OS may implement this lazily). Can transition efficiently to</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//               Ready. Accessing memory in such a region is undefined (may</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//               fault, may give back unexpected zeroes, etc.).</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// 4) Ready - may be accessed safely.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// This set of states is more than is strictly necessary to support all the</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// currently supported platforms. One could get by with just None, Reserved, and</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// Ready. However, the Prepared state gives us flexibility for performance</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// purposes. For example, on POSIX-y operating systems, Reserved is usually a</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// private anonymous mmap&#39;d region with PROT_NONE set, and to transition</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// to Ready would require setting PROT_READ|PROT_WRITE. However the</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// underspecification of Prepared lets us use just MADV_FREE to transition from</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// Ready to Prepared. Thus with the Prepared state we can set the permission</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// bits just once early on, we can efficiently tell the OS that it&#39;s free to</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// take pages away from us when we don&#39;t strictly need them.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// This file defines a cross-OS interface for a common set of helpers</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// that transition memory regions between these states. The helpers call into</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// OS-specific implementations that handle errors, while the interface boundary</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// implements cross-OS functionality, like updating runtime accounting.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// sysAlloc transitions an OS-chosen region of memory from None to Ready.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// More specifically, it obtains a large chunk of zeroed memory from the</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// operating system, typically on the order of a hundred kilobytes</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// or a megabyte. This memory is always immediately available for use.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// sysStat must be non-nil.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// Don&#39;t split the stack as this function may be invoked without a valid G,</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// which prevents us from allocating more stack.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>func sysAlloc(n uintptr, sysStat *sysMemStat) unsafe.Pointer {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	sysStat.add(int64(n))
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	gcController.mappedReady.Add(int64(n))
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	return sysAllocOS(n)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// sysUnused transitions a memory region from Ready to Prepared. It notifies the</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// operating system that the physical pages backing this memory region are no</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// longer needed and can be reused for other purposes. The contents of a</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// sysUnused memory region are considered forfeit and the region must not be</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// accessed again until sysUsed is called.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func sysUnused(v unsafe.Pointer, n uintptr) {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	gcController.mappedReady.Add(-int64(n))
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	sysUnusedOS(v, n)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// sysUsed transitions a memory region from Prepared to Ready. It notifies the</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// operating system that the memory region is needed and ensures that the region</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// may be safely accessed. This is typically a no-op on systems that don&#39;t have</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// an explicit commit step and hard over-commit limits, but is critical on</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// Windows, for example.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// This operation is idempotent for memory already in the Prepared state, so</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// it is safe to refer, with v and n, to a range of memory that includes both</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// Prepared and Ready memory. However, the caller must provide the exact amount</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// of Prepared memory for accounting purposes.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>func sysUsed(v unsafe.Pointer, n, prepared uintptr) {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	gcController.mappedReady.Add(int64(prepared))
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	sysUsedOS(v, n)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// sysHugePage does not transition memory regions, but instead provides a</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// hint to the OS that it would be more efficient to back this memory region</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// with pages of a larger size transparently.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func sysHugePage(v unsafe.Pointer, n uintptr) {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	sysHugePageOS(v, n)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// sysNoHugePage does not transition memory regions, but instead provides a</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// hint to the OS that it would be less efficient to back this memory region</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// with pages of a larger size transparently.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>func sysNoHugePage(v unsafe.Pointer, n uintptr) {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	sysNoHugePageOS(v, n)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// sysHugePageCollapse attempts to immediately back the provided memory region</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// with huge pages. It is best-effort and may fail silently.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func sysHugePageCollapse(v unsafe.Pointer, n uintptr) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	sysHugePageCollapseOS(v, n)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// sysFree transitions a memory region from any state to None. Therefore, it</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// returns memory unconditionally. It is used if an out-of-memory error has been</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// detected midway through an allocation or to carve out an aligned section of</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// the address space. It is okay if sysFree is a no-op only if sysReserve always</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// returns a memory region aligned to the heap allocator&#39;s alignment</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// restrictions.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// sysStat must be non-nil.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// Don&#39;t split the stack as this function may be invoked without a valid G,</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// which prevents us from allocating more stack.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func sysFree(v unsafe.Pointer, n uintptr, sysStat *sysMemStat) {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	sysStat.add(-int64(n))
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	gcController.mappedReady.Add(-int64(n))
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	sysFreeOS(v, n)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// sysFault transitions a memory region from Ready to Reserved. It</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// marks a region such that it will always fault if accessed. Used only for</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// debugging the runtime.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// TODO(mknyszek): Currently it&#39;s true that all uses of sysFault transition</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// memory from Ready to Reserved, but this may not be true in the future</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// since on every platform the operation is much more general than that.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// If a transition from Prepared is ever introduced, create a new function</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// that elides the Ready state accounting.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>func sysFault(v unsafe.Pointer, n uintptr) {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	gcController.mappedReady.Add(-int64(n))
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	sysFaultOS(v, n)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">// sysReserve transitions a memory region from None to Reserved. It reserves</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">// address space in such a way that it would cause a fatal fault upon access</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">// (either via permissions or not committing the memory). Such a reservation is</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// thus never backed by physical memory.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">// If the pointer passed to it is non-nil, the caller wants the</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// reservation there, but sysReserve can still choose another</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// location if that one is unavailable.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// NOTE: sysReserve returns OS-aligned memory, but the heap allocator</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// may use larger alignment, so the caller must be careful to realign the</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// memory obtained by sysReserve.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>func sysReserve(v unsafe.Pointer, n uintptr) unsafe.Pointer {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	return sysReserveOS(v, n)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// sysMap transitions a memory region from Reserved to Prepared. It ensures the</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// memory region can be efficiently transitioned to Ready.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">// sysStat must be non-nil.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>func sysMap(v unsafe.Pointer, n uintptr, sysStat *sysMemStat) {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	sysStat.add(int64(n))
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	sysMapOS(v, n)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
</pre><p><a href="mem.go?m=text">View as plain text</a></p>

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
