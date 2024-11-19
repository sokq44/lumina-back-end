<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/trace2map.go - Go Documentation Server</title>

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
<a href="trace2map.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">trace2map.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build goexperiment.exectracer2</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Simple hash table for tracing. Provides a mapping</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// between variable-length data and a unique ID. Subsequent</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// puts of the same data will return the same ID.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// Uses a region-based allocation scheme and assumes that the</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// table doesn&#39;t ever grow very big.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// This is definitely not a general-purpose hash table! It avoids</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// doing any high-level Go operations so it&#39;s safe to use even in</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// sensitive contexts.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>package runtime
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>import (
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>type traceMap struct {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	lock mutex <span class="comment">// Must be acquired on the system stack</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	seq  atomic.Uint64
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	mem  traceRegionAlloc
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	tab  [1 &lt;&lt; 13]atomic.UnsafePointer <span class="comment">// *traceMapNode (can&#39;t use generics because it&#39;s notinheap)</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>type traceMapNode struct {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	_    sys.NotInHeap
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	link atomic.UnsafePointer <span class="comment">// *traceMapNode (can&#39;t use generics because it&#39;s notinheap)</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	hash uintptr
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	id   uint64
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	data []byte
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// next is a type-safe wrapper around link.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func (n *traceMapNode) next() *traceMapNode {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	return (*traceMapNode)(n.link.Load())
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// stealID steals an ID from the table, ensuring that it will not</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// appear in the table anymore.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>func (tab *traceMap) stealID() uint64 {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	return tab.seq.Add(1)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// put inserts the data into the table.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// It&#39;s always safe to noescape data because its bytes are always copied.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// Returns a unique ID for the data and whether this is the first time</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// the data has been added to the map.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func (tab *traceMap) put(data unsafe.Pointer, size uintptr) (uint64, bool) {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if size == 0 {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		return 0, false
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	hash := memhash(data, 0, size)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// First, search the hashtable w/o the mutex.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	if id := tab.find(data, size, hash); id != 0 {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		return id, false
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// Now, double check under the mutex.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// Switch to the system stack so we can acquire tab.lock</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	var id uint64
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	var added bool
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		lock(&amp;tab.lock)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		if id = tab.find(data, size, hash); id != 0 {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			unlock(&amp;tab.lock)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			return
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		<span class="comment">// Create new record.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		id = tab.seq.Add(1)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		vd := tab.newTraceMapNode(data, size, hash, id)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		<span class="comment">// Insert it into the table.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		<span class="comment">// Update the link first, since the node isn&#39;t published yet.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		<span class="comment">// Then, store the node in the table as the new first node</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		<span class="comment">// for the bucket.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		part := int(hash % uintptr(len(tab.tab)))
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		vd.link.StoreNoWB(tab.tab[part].Load())
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		tab.tab[part].StoreNoWB(unsafe.Pointer(vd))
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		unlock(&amp;tab.lock)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		added = true
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	})
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	return id, added
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// find looks up data in the table, assuming hash is a hash of data.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// Returns 0 if the data is not found, and the unique ID for it if it is.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>func (tab *traceMap) find(data unsafe.Pointer, size, hash uintptr) uint64 {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	part := int(hash % uintptr(len(tab.tab)))
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	for vd := tab.bucket(part); vd != nil; vd = vd.next() {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		<span class="comment">// Synchronization not necessary. Once published to the table, these</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		<span class="comment">// values are immutable.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		if vd.hash == hash &amp;&amp; uintptr(len(vd.data)) == size {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			if memequal(unsafe.Pointer(&amp;vd.data[0]), data, size) {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>				return vd.id
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	return 0
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// bucket is a type-safe wrapper for looking up a value in tab.tab.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>func (tab *traceMap) bucket(part int) *traceMapNode {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	return (*traceMapNode)(tab.tab[part].Load())
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>func (tab *traceMap) newTraceMapNode(data unsafe.Pointer, size, hash uintptr, id uint64) *traceMapNode {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// Create data array.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	sl := notInHeapSlice{
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		array: tab.mem.alloc(size),
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		len:   int(size),
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		cap:   int(size),
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	memmove(unsafe.Pointer(sl.array), data, size)
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// Create metadata structure.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	meta := (*traceMapNode)(unsafe.Pointer(tab.mem.alloc(unsafe.Sizeof(traceMapNode{}))))
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	*(*notInHeapSlice)(unsafe.Pointer(&amp;meta.data)) = sl
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	meta.id = id
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	meta.hash = hash
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	return meta
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">// reset drops all allocated memory from the table and resets it.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// tab.lock must be held. Must run on the system stack because of this.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>func (tab *traceMap) reset() {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	assertLockHeld(&amp;tab.lock)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	tab.mem.drop()
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	tab.seq.Store(0)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// Clear table without write barriers. The table consists entirely</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// of notinheap pointers, so this is fine.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">// Write barriers may theoretically call into the tracer and acquire</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// the lock again, and this lock ordering is expressed in the static</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// lock ranking checker.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	memclrNoHeapPointers(unsafe.Pointer(&amp;tab.tab), unsafe.Sizeof(tab.tab))
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
</pre><p><a href="trace2map.go?m=text">View as plain text</a></p>

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
