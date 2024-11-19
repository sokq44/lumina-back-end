<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mfixalloc.go - Go Documentation Server</title>

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
<a href="mfixalloc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mfixalloc.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Fixed-size object allocator. Returned memory is not zeroed.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// See malloc.go for overview.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package runtime
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// fixalloc is a simple free-list allocator for fixed size objects.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Malloc uses a FixAlloc wrapped around sysAlloc to manage its</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// mcache and mspan objects.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// Memory returned by fixalloc.alloc is zeroed by default, but the</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// caller may take responsibility for zeroing allocations by setting</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// the zero flag to false. This is only safe if the memory never</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// contains heap pointers.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// The caller is responsible for locking around FixAlloc calls.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// Callers can keep state in the object but the first word is</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// smashed by freeing and reallocating.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// Consider marking fixalloc&#39;d types not in heap by embedding</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// runtime/internal/sys.NotInHeap.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>type fixalloc struct {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	size   uintptr
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	first  func(arg, p unsafe.Pointer) <span class="comment">// called first time p is returned</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	arg    unsafe.Pointer
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	list   *mlink
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	chunk  uintptr <span class="comment">// use uintptr instead of unsafe.Pointer to avoid write barriers</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	nchunk uint32  <span class="comment">// bytes remaining in current chunk</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	nalloc uint32  <span class="comment">// size of new chunks in bytes</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	inuse  uintptr <span class="comment">// in-use bytes now</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	stat   *sysMemStat
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	zero   bool <span class="comment">// zero allocations</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>}
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// A generic linked list of blocks.  (Typically the block is bigger than sizeof(MLink).)</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// Since assignments to mlink.next will result in a write barrier being performed</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// this cannot be used by some of the internal GC structures. For example when</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// the sweeper is placing an unmarked object on the free list it does not want the</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// write barrier to be called since that could result in the object being reachable.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>type mlink struct {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	_    sys.NotInHeap
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	next *mlink
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// Initialize f to allocate objects of the given size,</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// using the allocator to obtain chunks of memory.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>func (f *fixalloc) init(size uintptr, first func(arg, p unsafe.Pointer), arg unsafe.Pointer, stat *sysMemStat) {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	if size &gt; _FixAllocChunk {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		throw(&#34;runtime: fixalloc size too large&#34;)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	size = max(size, unsafe.Sizeof(mlink{}))
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	f.size = size
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	f.first = first
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	f.arg = arg
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	f.list = nil
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	f.chunk = 0
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	f.nchunk = 0
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	f.nalloc = uint32(_FixAllocChunk / size * size) <span class="comment">// Round _FixAllocChunk down to an exact multiple of size to eliminate tail waste</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	f.inuse = 0
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	f.stat = stat
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	f.zero = true
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>func (f *fixalloc) alloc() unsafe.Pointer {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	if f.size == 0 {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		print(&#34;runtime: use of FixAlloc_Alloc before FixAlloc_Init\n&#34;)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		throw(&#34;runtime: internal error&#34;)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if f.list != nil {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		v := unsafe.Pointer(f.list)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		f.list = f.list.next
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		f.inuse += f.size
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		if f.zero {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			memclrNoHeapPointers(v, f.size)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		return v
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	if uintptr(f.nchunk) &lt; f.size {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		f.chunk = uintptr(persistentalloc(uintptr(f.nalloc), 0, f.stat))
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		f.nchunk = f.nalloc
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	v := unsafe.Pointer(f.chunk)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	if f.first != nil {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		f.first(f.arg, v)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	f.chunk = f.chunk + f.size
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	f.nchunk -= uint32(f.size)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	f.inuse += f.size
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	return v
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>func (f *fixalloc) free(p unsafe.Pointer) {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	f.inuse -= f.size
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	v := (*mlink)(p)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	v.next = f.list
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	f.list = v
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
</pre><p><a href="mfixalloc.go?m=text">View as plain text</a></p>

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
