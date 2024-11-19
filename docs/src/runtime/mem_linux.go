<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mem_linux.go - Go Documentation Server</title>

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
<a href="mem_linux.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mem_linux.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2010 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>const (
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	_EACCES = 13
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	_EINVAL = 22
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Don&#39;t split the stack as this method may be invoked without a valid G, which</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// prevents us from allocating more stack.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>func sysAllocOS(n uintptr) unsafe.Pointer {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	p, err := mmap(nil, n, _PROT_READ|_PROT_WRITE, _MAP_ANON|_MAP_PRIVATE, -1, 0)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	if err != 0 {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>		if err == _EACCES {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>			print(&#34;runtime: mmap: access denied\n&#34;)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>			exit(2)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		if err == _EAGAIN {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>			print(&#34;runtime: mmap: too much locked memory (check &#39;ulimit -l&#39;).\n&#34;)
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>			exit(2)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		return nil
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	return p
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>var adviseUnused = uint32(_MADV_FREE)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>const madviseUnsupported = 0
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>func sysUnusedOS(v unsafe.Pointer, n uintptr) {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	if uintptr(v)&amp;(physPageSize-1) != 0 || n&amp;(physPageSize-1) != 0 {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		<span class="comment">// madvise will round this to any physical page</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		<span class="comment">// *covered* by this range, so an unaligned madvise</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		<span class="comment">// will release more memory than intended.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		throw(&#34;unaligned sysUnused&#34;)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	advise := atomic.Load(&amp;adviseUnused)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	if debug.madvdontneed != 0 &amp;&amp; advise != madviseUnsupported {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		advise = _MADV_DONTNEED
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	switch advise {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	case _MADV_FREE:
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		if madvise(v, n, _MADV_FREE) == 0 {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			break
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		atomic.Store(&amp;adviseUnused, _MADV_DONTNEED)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		fallthrough
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	case _MADV_DONTNEED:
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		<span class="comment">// MADV_FREE was added in Linux 4.5. Fall back on MADV_DONTNEED if it&#39;s</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		<span class="comment">// not supported.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		if madvise(v, n, _MADV_DONTNEED) == 0 {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			break
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		atomic.Store(&amp;adviseUnused, madviseUnsupported)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		fallthrough
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	case madviseUnsupported:
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		<span class="comment">// Since Linux 3.18, support for madvise is optional.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		<span class="comment">// Fall back on mmap if it&#39;s not supported.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		<span class="comment">// _MAP_ANON|_MAP_FIXED|_MAP_PRIVATE will unmap all the</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		<span class="comment">// pages in the old mapping, and remap the memory region.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		mmap(v, n, _PROT_READ|_PROT_WRITE, _MAP_ANON|_MAP_FIXED|_MAP_PRIVATE, -1, 0)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	if debug.harddecommit &gt; 0 {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		p, err := mmap(v, n, _PROT_NONE, _MAP_ANON|_MAP_FIXED|_MAP_PRIVATE, -1, 0)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		if p != v || err != 0 {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			throw(&#34;runtime: cannot disable permissions in address space&#34;)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>func sysUsedOS(v unsafe.Pointer, n uintptr) {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	if debug.harddecommit &gt; 0 {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		p, err := mmap(v, n, _PROT_READ|_PROT_WRITE, _MAP_ANON|_MAP_FIXED|_MAP_PRIVATE, -1, 0)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		if err == _ENOMEM {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			throw(&#34;runtime: out of memory&#34;)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		if p != v || err != 0 {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			throw(&#34;runtime: cannot remap pages in address space&#34;)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		return
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func sysHugePageOS(v unsafe.Pointer, n uintptr) {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	if physHugePageSize != 0 {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		<span class="comment">// Round v up to a huge page boundary.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		beg := alignUp(uintptr(v), physHugePageSize)
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		<span class="comment">// Round v+n down to a huge page boundary.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		end := alignDown(uintptr(v)+n, physHugePageSize)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		if beg &lt; end {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			madvise(unsafe.Pointer(beg), end-beg, _MADV_HUGEPAGE)
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>func sysNoHugePageOS(v unsafe.Pointer, n uintptr) {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	if uintptr(v)&amp;(physPageSize-1) != 0 {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		<span class="comment">// The Linux implementation requires that the address</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		<span class="comment">// addr be page-aligned, and allows length to be zero.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		throw(&#34;unaligned sysNoHugePageOS&#34;)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	madvise(v, n, _MADV_NOHUGEPAGE)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>func sysHugePageCollapseOS(v unsafe.Pointer, n uintptr) {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	if uintptr(v)&amp;(physPageSize-1) != 0 {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		<span class="comment">// The Linux implementation requires that the address</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		<span class="comment">// addr be page-aligned, and allows length to be zero.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		throw(&#34;unaligned sysHugePageCollapseOS&#34;)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	if physHugePageSize == 0 {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		return
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// N.B. If you find yourself debugging this code, note that</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// this call can fail with EAGAIN because it&#39;s best-effort.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// Also, when it returns an error, it&#39;s only for the last</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// huge page in the region requested.</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// It can also sometimes return EINVAL if the corresponding</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// region hasn&#39;t been backed by physical memory. This is</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// difficult to guarantee in general, and it also means</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// there&#39;s no way to distinguish whether this syscall is</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// actually available. Oops.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// Anyway, that&#39;s why this call just doesn&#39;t bother checking</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// any errors.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	madvise(v, n, _MADV_COLLAPSE)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// Don&#39;t split the stack as this function may be invoked without a valid G,</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// which prevents us from allocating more stack.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>func sysFreeOS(v unsafe.Pointer, n uintptr) {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	munmap(v, n)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func sysFaultOS(v unsafe.Pointer, n uintptr) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	mmap(v, n, _PROT_NONE, _MAP_ANON|_MAP_PRIVATE|_MAP_FIXED, -1, 0)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>func sysReserveOS(v unsafe.Pointer, n uintptr) unsafe.Pointer {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	p, err := mmap(v, n, _PROT_NONE, _MAP_ANON|_MAP_PRIVATE, -1, 0)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	if err != 0 {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		return nil
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	return p
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>func sysMapOS(v unsafe.Pointer, n uintptr) {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	p, err := mmap(v, n, _PROT_READ|_PROT_WRITE, _MAP_ANON|_MAP_FIXED|_MAP_PRIVATE, -1, 0)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	if err == _ENOMEM {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		throw(&#34;runtime: out of memory&#34;)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	if p != v || err != 0 {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		print(&#34;runtime: mmap(&#34;, v, &#34;, &#34;, n, &#34;) returned &#34;, p, &#34;, &#34;, err, &#34;\n&#34;)
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		throw(&#34;runtime: cannot map pages in arena address space&#34;)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// Disable huge pages if the GODEBUG for it is set.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// Note that there are a few sysHugePage calls that can override this, but</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// they&#39;re all for GC metadata.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	if debug.disablethp != 0 {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		sysNoHugePageOS(v, n)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
</pre><p><a href="mem_linux.go?m=text">View as plain text</a></p>

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
