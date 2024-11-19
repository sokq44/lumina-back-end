<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/heapdump.go - Go Documentation Server</title>

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
<a href="heapdump.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">heapdump.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Implementation of runtime/debug.WriteHeapDump. Writes all</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// objects in the heap plus additional info (roots, threads,</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// finalizers, etc.) to a file.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// The format of the dumped file is described at</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// https://golang.org/s/go15heapdump.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>package runtime
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>import (
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;internal/goexperiment&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//go:linkname runtime_debug_WriteHeapDump runtime/debug.WriteHeapDump</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>func runtime_debug_WriteHeapDump(fd uintptr) {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	stw := stopTheWorld(stwWriteHeapDump)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// Keep m on this G&#39;s stack instead of the system stack.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// Both readmemstats_m and writeheapdump_m have pretty large</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// peak stack depths and we risk blowing the system stack.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// This is safe because the world is stopped, so we don&#39;t</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// need to worry about anyone shrinking and therefore moving</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// our stack.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	var m MemStats
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		<span class="comment">// Call readmemstats_m here instead of deeper in</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		<span class="comment">// writeheapdump_m because we might blow the system stack</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		<span class="comment">// otherwise.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		readmemstats_m(&amp;m)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		writeheapdump_m(fd, &amp;m)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	})
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	startTheWorld(stw)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>const (
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	fieldKindEol       = 0
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	fieldKindPtr       = 1
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	fieldKindIface     = 2
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	fieldKindEface     = 3
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	tagEOF             = 0
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	tagObject          = 1
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	tagOtherRoot       = 2
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	tagType            = 3
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	tagGoroutine       = 4
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	tagStackFrame      = 5
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	tagParams          = 6
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	tagFinalizer       = 7
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	tagItab            = 8
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	tagOSThread        = 9
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	tagMemStats        = 10
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	tagQueuedFinalizer = 11
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	tagData            = 12
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	tagBSS             = 13
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	tagDefer           = 14
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	tagPanic           = 15
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	tagMemProf         = 16
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	tagAllocSample     = 17
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>var dumpfd uintptr <span class="comment">// fd to write the dump to.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>var tmpbuf []byte
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// buffer of pending write data</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>const (
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	bufSize = 4096
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>var buf [bufSize]byte
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>var nbuf uintptr
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func dwrite(data unsafe.Pointer, len uintptr) {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if len == 0 {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		return
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	if nbuf+len &lt;= bufSize {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		copy(buf[nbuf:], (*[bufSize]byte)(data)[:len])
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		nbuf += len
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		return
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	write(dumpfd, unsafe.Pointer(&amp;buf), int32(nbuf))
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	if len &gt;= bufSize {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		write(dumpfd, data, int32(len))
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		nbuf = 0
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	} else {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		copy(buf[:], (*[bufSize]byte)(data)[:len])
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		nbuf = len
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>func dwritebyte(b byte) {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	dwrite(unsafe.Pointer(&amp;b), 1)
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>func flush() {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	write(dumpfd, unsafe.Pointer(&amp;buf), int32(nbuf))
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	nbuf = 0
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// Cache of types that have been serialized already.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// We use a type&#39;s hash field to pick a bucket.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// Inside a bucket, we keep a list of types that</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// have been serialized so far, most recently used first.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// Note: when a bucket overflows we may end up</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// serializing a type more than once. That&#39;s ok.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>const (
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	typeCacheBuckets = 256
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	typeCacheAssoc   = 4
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>type typeCacheBucket struct {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	t [typeCacheAssoc]*_type
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>var typecache [typeCacheBuckets]typeCacheBucket
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// dump a uint64 in a varint format parseable by encoding/binary.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>func dumpint(v uint64) {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	var buf [10]byte
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	var n int
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	for v &gt;= 0x80 {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		buf[n] = byte(v | 0x80)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		n++
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		v &gt;&gt;= 7
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	buf[n] = byte(v)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	n++
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	dwrite(unsafe.Pointer(&amp;buf), uintptr(n))
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>func dumpbool(b bool) {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	if b {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		dumpint(1)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	} else {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		dumpint(0)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// dump varint uint64 length followed by memory contents.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>func dumpmemrange(data unsafe.Pointer, len uintptr) {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	dumpint(uint64(len))
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	dwrite(data, len)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>func dumpslice(b []byte) {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	dumpint(uint64(len(b)))
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	if len(b) &gt; 0 {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		dwrite(unsafe.Pointer(&amp;b[0]), uintptr(len(b)))
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func dumpstr(s string) {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	dumpmemrange(unsafe.Pointer(unsafe.StringData(s)), uintptr(len(s)))
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// dump information for a type.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>func dumptype(t *_type) {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	if t == nil {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		return
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;ve definitely serialized the type before,</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// no need to do it again.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	b := &amp;typecache[t.Hash&amp;(typeCacheBuckets-1)]
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	if t == b.t[0] {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		return
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	for i := 1; i &lt; typeCacheAssoc; i++ {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		if t == b.t[i] {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			<span class="comment">// Move-to-front</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			for j := i; j &gt; 0; j-- {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>				b.t[j] = b.t[j-1]
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			b.t[0] = t
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			return
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// Might not have been dumped yet. Dump it and</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">// remember we did so.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	for j := typeCacheAssoc - 1; j &gt; 0; j-- {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		b.t[j] = b.t[j-1]
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	b.t[0] = t
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// dump the type</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	dumpint(tagType)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(t))))
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	dumpint(uint64(t.Size_))
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	rt := toRType(t)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	if x := t.Uncommon(); x == nil || rt.nameOff(x.PkgPath).Name() == &#34;&#34; {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		dumpstr(rt.string())
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	} else {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		pkgpath := rt.nameOff(x.PkgPath).Name()
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		name := rt.name()
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(len(pkgpath)) + 1 + uintptr(len(name))))
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		dwrite(unsafe.Pointer(unsafe.StringData(pkgpath)), uintptr(len(pkgpath)))
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		dwritebyte(&#39;.&#39;)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		dwrite(unsafe.Pointer(unsafe.StringData(name)), uintptr(len(name)))
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	dumpbool(t.Kind_&amp;kindDirectIface == 0 || t.PtrBytes != 0)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// dump an object.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>func dumpobj(obj unsafe.Pointer, size uintptr, bv bitvector) {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	dumpint(tagObject)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(obj)))
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	dumpmemrange(obj, size)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	dumpfields(bv)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>func dumpotherroot(description string, to unsafe.Pointer) {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	dumpint(tagOtherRoot)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	dumpstr(description)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(to)))
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>func dumpfinalizer(obj unsafe.Pointer, fn *funcval, fint *_type, ot *ptrtype) {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	dumpint(tagFinalizer)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(obj)))
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(fn))))
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(fn.fn))))
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(fint))))
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(ot))))
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>type childInfo struct {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// Information passed up from the callee frame about</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// the layout of the outargs region.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	argoff uintptr   <span class="comment">// where the arguments start in the frame</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	arglen uintptr   <span class="comment">// size of args region</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	args   bitvector <span class="comment">// if args.n &gt;= 0, pointer map of args region</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	sp     *uint8    <span class="comment">// callee sp</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	depth  uintptr   <span class="comment">// depth in call stack (0 == most recent)</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// dump kinds &amp; offsets of interesting fields in bv.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>func dumpbv(cbv *bitvector, offset uintptr) {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; uintptr(cbv.n); i++ {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		if cbv.ptrbit(i) == 1 {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			dumpint(fieldKindPtr)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			dumpint(uint64(offset + i*goarch.PtrSize))
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>func dumpframe(s *stkframe, child *childInfo) {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	f := s.fn
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	<span class="comment">// Figure out what we can about our stack map</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	pc := s.pc
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	pcdata := int32(-1) <span class="comment">// Use the entry map at function entry</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	if pc != f.entry() {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		pc--
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		pcdata = pcdatavalue(f, abi.PCDATA_StackMapIndex, pc)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	if pcdata == -1 {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		<span class="comment">// We do not have a valid pcdata value but there might be a</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		<span class="comment">// stackmap for this function. It is likely that we are looking</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		<span class="comment">// at the function prologue, assume so and hope for the best.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		pcdata = 0
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	stkmap := (*stackmap)(funcdata(f, abi.FUNCDATA_LocalsPointerMaps))
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	var bv bitvector
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	if stkmap != nil &amp;&amp; stkmap.n &gt; 0 {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		bv = stackmapdata(stkmap, pcdata)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	} else {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		bv.n = -1
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	<span class="comment">// Dump main body of stack frame.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	dumpint(tagStackFrame)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	dumpint(uint64(s.sp))                              <span class="comment">// lowest address in frame</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	dumpint(uint64(child.depth))                       <span class="comment">// # of frames deep on the stack</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(child.sp)))) <span class="comment">// sp of child, or 0 if bottom of stack</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	dumpmemrange(unsafe.Pointer(s.sp), s.fp-s.sp)      <span class="comment">// frame contents</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	dumpint(uint64(f.entry()))
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	dumpint(uint64(s.pc))
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	dumpint(uint64(s.continpc))
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	name := funcname(f)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	if name == &#34;&#34; {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		name = &#34;unknown function&#34;
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	dumpstr(name)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// Dump fields in the outargs section</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	if child.args.n &gt;= 0 {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		dumpbv(&amp;child.args, child.argoff)
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	} else {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		<span class="comment">// conservative - everything might be a pointer</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		for off := child.argoff; off &lt; child.argoff+child.arglen; off += goarch.PtrSize {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			dumpint(fieldKindPtr)
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			dumpint(uint64(off))
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	<span class="comment">// Dump fields in the local vars section</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if stkmap == nil {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		<span class="comment">// No locals information, dump everything.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		for off := child.arglen; off &lt; s.varp-s.sp; off += goarch.PtrSize {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>			dumpint(fieldKindPtr)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>			dumpint(uint64(off))
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	} else if stkmap.n &lt; 0 {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		<span class="comment">// Locals size information, dump just the locals.</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		size := uintptr(-stkmap.n)
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		for off := s.varp - size - s.sp; off &lt; s.varp-s.sp; off += goarch.PtrSize {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			dumpint(fieldKindPtr)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>			dumpint(uint64(off))
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	} else if stkmap.n &gt; 0 {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		<span class="comment">// Locals bitmap information, scan just the pointers in</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		<span class="comment">// locals.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		dumpbv(&amp;bv, s.varp-uintptr(bv.n)*goarch.PtrSize-s.sp)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	dumpint(fieldKindEol)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// Record arg info for parent.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	child.argoff = s.argp - s.fp
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	child.arglen = s.argBytes()
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	child.sp = (*uint8)(unsafe.Pointer(s.sp))
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	child.depth++
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	stkmap = (*stackmap)(funcdata(f, abi.FUNCDATA_ArgsPointerMaps))
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	if stkmap != nil {
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		child.args = stackmapdata(stkmap, pcdata)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	} else {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		child.args.n = -1
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	return
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>func dumpgoroutine(gp *g) {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	var sp, pc, lr uintptr
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	if gp.syscallsp != 0 {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		sp = gp.syscallsp
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		pc = gp.syscallpc
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		lr = 0
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	} else {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		sp = gp.sched.sp
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		pc = gp.sched.pc
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		lr = gp.sched.lr
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	dumpint(tagGoroutine)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(gp))))
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	dumpint(uint64(sp))
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	dumpint(gp.goid)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	dumpint(uint64(gp.gopc))
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	dumpint(uint64(readgstatus(gp)))
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	dumpbool(isSystemGoroutine(gp, false))
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	dumpbool(false) <span class="comment">// isbackground</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	dumpint(uint64(gp.waitsince))
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	dumpstr(gp.waitreason.String())
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(gp.sched.ctxt)))
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(gp.m))))
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(gp._defer))))
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(gp._panic))))
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	<span class="comment">// dump stack</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	var child childInfo
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	child.args.n = -1
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	child.arglen = 0
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	child.sp = nil
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	child.depth = 0
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	var u unwinder
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	for u.initAt(pc, sp, lr, gp, 0); u.valid(); u.next() {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		dumpframe(&amp;u.frame, &amp;child)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	<span class="comment">// dump defer &amp; panic records</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	for d := gp._defer; d != nil; d = d.link {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		dumpint(tagDefer)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(unsafe.Pointer(d))))
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(unsafe.Pointer(gp))))
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		dumpint(uint64(d.sp))
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		dumpint(uint64(d.pc))
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		fn := *(**funcval)(unsafe.Pointer(&amp;d.fn))
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(unsafe.Pointer(fn))))
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		if d.fn == nil {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			<span class="comment">// d.fn can be nil for open-coded defers</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			dumpint(uint64(0))
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		} else {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			dumpint(uint64(uintptr(unsafe.Pointer(fn.fn))))
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(unsafe.Pointer(d.link))))
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	}
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	for p := gp._panic; p != nil; p = p.link {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		dumpint(tagPanic)
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(unsafe.Pointer(p))))
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(unsafe.Pointer(gp))))
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		eface := efaceOf(&amp;p.arg)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(unsafe.Pointer(eface._type))))
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(eface.data)))
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		dumpint(0) <span class="comment">// was p-&gt;defer, no longer recorded</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(unsafe.Pointer(p.link))))
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>func dumpgs() {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	<span class="comment">// goroutines &amp; stacks</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	forEachG(func(gp *g) {
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		status := readgstatus(gp) <span class="comment">// The world is stopped so gp will not be in a scan state.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		switch status {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		default:
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			print(&#34;runtime: unexpected G.status &#34;, hex(status), &#34;\n&#34;)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			throw(&#34;dumpgs in STW - bad status&#34;)
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		case _Gdead:
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			<span class="comment">// ok</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		case _Grunnable,
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>			_Gsyscall,
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			_Gwaiting:
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			dumpgoroutine(gp)
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	})
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>func finq_callback(fn *funcval, obj unsafe.Pointer, nret uintptr, fint *_type, ot *ptrtype) {
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	dumpint(tagQueuedFinalizer)
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(obj)))
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(fn))))
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(fn.fn))))
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(fint))))
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(ot))))
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>func dumproots() {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	<span class="comment">// To protect mheap_.allspans.</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// TODO(mwhudson): dump datamask etc from all objects</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	<span class="comment">// data segment</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	dumpint(tagData)
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	dumpint(uint64(firstmoduledata.data))
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	dumpmemrange(unsafe.Pointer(firstmoduledata.data), firstmoduledata.edata-firstmoduledata.data)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	dumpfields(firstmoduledata.gcdatamask)
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	<span class="comment">// bss segment</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	dumpint(tagBSS)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	dumpint(uint64(firstmoduledata.bss))
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	dumpmemrange(unsafe.Pointer(firstmoduledata.bss), firstmoduledata.ebss-firstmoduledata.bss)
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	dumpfields(firstmoduledata.gcbssmask)
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	<span class="comment">// mspan.types</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	for _, s := range mheap_.allspans {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		if s.state.get() == mSpanInUse {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			<span class="comment">// Finalizers</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			for sp := s.specials; sp != nil; sp = sp.next {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>				if sp.kind != _KindSpecialFinalizer {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>					continue
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>				}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>				spf := (*specialfinalizer)(unsafe.Pointer(sp))
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>				p := unsafe.Pointer(s.base() + uintptr(spf.special.offset))
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>				dumpfinalizer(p, spf.fn, spf.fint, spf.ot)
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>			}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	<span class="comment">// Finalizer queue</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	iterate_finq(finq_callback)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span><span class="comment">// Bit vector of free marks.</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span><span class="comment">// Needs to be as big as the largest number of objects per span.</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>var freemark [_PageSize / 8]bool
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>func dumpobjs() {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	<span class="comment">// To protect mheap_.allspans.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	for _, s := range mheap_.allspans {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		if s.state.get() != mSpanInUse {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			continue
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		p := s.base()
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		size := s.elemsize
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		n := (s.npages &lt;&lt; _PageShift) / size
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		if n &gt; uintptr(len(freemark)) {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>			throw(&#34;freemark array doesn&#39;t have enough entries&#34;)
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		for freeIndex := uint16(0); freeIndex &lt; s.nelems; freeIndex++ {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>			if s.isFree(uintptr(freeIndex)) {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>				freemark[freeIndex] = true
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		for j := uintptr(0); j &lt; n; j, p = j+1, p+size {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			if freemark[j] {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>				freemark[j] = false
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>				continue
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			dumpobj(unsafe.Pointer(p), size, makeheapobjbv(p, size))
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>func dumpparams() {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	dumpint(tagParams)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	x := uintptr(1)
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	if *(*byte)(unsafe.Pointer(&amp;x)) == 1 {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		dumpbool(false) <span class="comment">// little-endian ptrs</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	} else {
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		dumpbool(true) <span class="comment">// big-endian ptrs</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	}
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	dumpint(goarch.PtrSize)
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	var arenaStart, arenaEnd uintptr
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	for i1 := range mheap_.arenas {
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		if mheap_.arenas[i1] == nil {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			continue
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		for i, ha := range mheap_.arenas[i1] {
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>			if ha == nil {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>				continue
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>			}
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>			base := arenaBase(arenaIdx(i1)&lt;&lt;arenaL1Shift | arenaIdx(i))
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			if arenaStart == 0 || base &lt; arenaStart {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>				arenaStart = base
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			}
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			if base+heapArenaBytes &gt; arenaEnd {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>				arenaEnd = base + heapArenaBytes
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>			}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	dumpint(uint64(arenaStart))
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	dumpint(uint64(arenaEnd))
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	dumpstr(goarch.GOARCH)
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	dumpstr(buildVersion)
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	dumpint(uint64(ncpu))
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>}
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>func itab_callback(tab *itab) {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	t := tab._type
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	dumptype(t)
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	dumpint(tagItab)
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(tab))))
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(t))))
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>func dumpitabs() {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	iterate_itabs(itab_callback)
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>}
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>func dumpms() {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	for mp := allm; mp != nil; mp = mp.alllink {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		dumpint(tagOSThread)
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		dumpint(uint64(uintptr(unsafe.Pointer(mp))))
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		dumpint(uint64(mp.id))
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		dumpint(mp.procid)
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>func dumpmemstats(m *MemStats) {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	<span class="comment">// These ints should be identical to the exported</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	<span class="comment">// MemStats structure and should be ordered the same</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	<span class="comment">// way too.</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	dumpint(tagMemStats)
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	dumpint(m.Alloc)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	dumpint(m.TotalAlloc)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	dumpint(m.Sys)
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	dumpint(m.Lookups)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	dumpint(m.Mallocs)
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	dumpint(m.Frees)
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	dumpint(m.HeapAlloc)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	dumpint(m.HeapSys)
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	dumpint(m.HeapIdle)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	dumpint(m.HeapInuse)
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	dumpint(m.HeapReleased)
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	dumpint(m.HeapObjects)
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	dumpint(m.StackInuse)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	dumpint(m.StackSys)
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	dumpint(m.MSpanInuse)
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	dumpint(m.MSpanSys)
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	dumpint(m.MCacheInuse)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	dumpint(m.MCacheSys)
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	dumpint(m.BuckHashSys)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	dumpint(m.GCSys)
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	dumpint(m.OtherSys)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	dumpint(m.NextGC)
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	dumpint(m.LastGC)
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	dumpint(m.PauseTotalNs)
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	for i := 0; i &lt; 256; i++ {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		dumpint(m.PauseNs[i])
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	dumpint(uint64(m.NumGC))
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>}
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>func dumpmemprof_callback(b *bucket, nstk uintptr, pstk *uintptr, size, allocs, frees uintptr) {
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	stk := (*[100000]uintptr)(unsafe.Pointer(pstk))
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	dumpint(tagMemProf)
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	dumpint(uint64(uintptr(unsafe.Pointer(b))))
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	dumpint(uint64(size))
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	dumpint(uint64(nstk))
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; nstk; i++ {
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		pc := stk[i]
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		f := findfunc(pc)
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		if !f.valid() {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>			var buf [64]byte
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>			n := len(buf)
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>			n--
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>			buf[n] = &#39;)&#39;
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>			if pc == 0 {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>				n--
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>				buf[n] = &#39;0&#39;
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>			} else {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>				for pc &gt; 0 {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>					n--
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>					buf[n] = &#34;0123456789abcdef&#34;[pc&amp;15]
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>					pc &gt;&gt;= 4
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>				}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>			}
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>			n--
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>			buf[n] = &#39;x&#39;
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			n--
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			buf[n] = &#39;0&#39;
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>			n--
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>			buf[n] = &#39;(&#39;
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>			dumpslice(buf[n:])
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			dumpstr(&#34;?&#34;)
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>			dumpint(0)
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		} else {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>			dumpstr(funcname(f))
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>			if i &gt; 0 &amp;&amp; pc &gt; f.entry() {
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>				pc--
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>			}
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>			file, line := funcline(f, pc)
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			dumpstr(file)
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>			dumpint(uint64(line))
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>		}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	dumpint(uint64(allocs))
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	dumpint(uint64(frees))
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>func dumpmemprof() {
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	<span class="comment">// To protect mheap_.allspans.</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	iterate_memprof(dumpmemprof_callback)
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	for _, s := range mheap_.allspans {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		if s.state.get() != mSpanInUse {
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>			continue
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		}
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		for sp := s.specials; sp != nil; sp = sp.next {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>			if sp.kind != _KindSpecialProfile {
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>				continue
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>			}
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>			spp := (*specialprofile)(unsafe.Pointer(sp))
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>			p := s.base() + uintptr(spp.special.offset)
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>			dumpint(tagAllocSample)
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>			dumpint(uint64(p))
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>			dumpint(uint64(uintptr(unsafe.Pointer(spp.b))))
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	}
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>}
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>var dumphdr = []byte(&#34;go1.7 heap dump\n&#34;)
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>func mdump(m *MemStats) {
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	<span class="comment">// make sure we&#39;re done sweeping</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	for _, s := range mheap_.allspans {
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		if s.state.get() == mSpanInUse {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>			s.ensureSwept()
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		}
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	}
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	memclrNoHeapPointers(unsafe.Pointer(&amp;typecache), unsafe.Sizeof(typecache))
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	dwrite(unsafe.Pointer(&amp;dumphdr[0]), uintptr(len(dumphdr)))
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	dumpparams()
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	dumpitabs()
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	dumpobjs()
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	dumpgs()
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	dumpms()
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	dumproots()
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	dumpmemstats(m)
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	dumpmemprof()
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	dumpint(tagEOF)
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	flush()
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>}
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>func writeheapdump_m(fd uintptr, m *MemStats) {
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	assertWorldStopped()
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	gp := getg()
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	casGToWaiting(gp.m.curg, _Grunning, waitReasonDumpingHeap)
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	<span class="comment">// Set dump file.</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	dumpfd = fd
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	<span class="comment">// Call dump routine.</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	mdump(m)
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	<span class="comment">// Reset dump file.</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	dumpfd = 0
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	if tmpbuf != nil {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>		sysFree(unsafe.Pointer(&amp;tmpbuf[0]), uintptr(len(tmpbuf)), &amp;memstats.other_sys)
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>		tmpbuf = nil
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	}
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	casgstatus(gp.m.curg, _Gwaiting, _Grunning)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span><span class="comment">// dumpint() the kind &amp; offset of each field in an object.</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>func dumpfields(bv bitvector) {
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	dumpbv(&amp;bv, 0)
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	dumpint(fieldKindEol)
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>}
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>func makeheapobjbv(p uintptr, size uintptr) bitvector {
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	<span class="comment">// Extend the temp buffer if necessary.</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	nptr := size / goarch.PtrSize
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	if uintptr(len(tmpbuf)) &lt; nptr/8+1 {
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		if tmpbuf != nil {
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>			sysFree(unsafe.Pointer(&amp;tmpbuf[0]), uintptr(len(tmpbuf)), &amp;memstats.other_sys)
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		}
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		n := nptr/8 + 1
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>		p := sysAlloc(n, &amp;memstats.other_sys)
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		if p == nil {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>			throw(&#34;heapdump: out of memory&#34;)
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		}
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>		tmpbuf = (*[1 &lt;&lt; 30]byte)(p)[:n]
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	}
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	<span class="comment">// Convert heap bitmap to pointer bitmap.</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; nptr/8+1; i++ {
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		tmpbuf[i] = 0
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	}
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	if goexperiment.AllocHeaders {
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		s := spanOf(p)
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>		tp := s.typePointersOf(p, size)
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		for {
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>			var addr uintptr
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>			if tp, addr = tp.next(p + size); addr == 0 {
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>				break
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>			}
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>			i := (addr - p) / goarch.PtrSize
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>			tmpbuf[i/8] |= 1 &lt;&lt; (i % 8)
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		}
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	} else {
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>		hbits := heapBitsForAddr(p, size)
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		for {
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>			var addr uintptr
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>			hbits, addr = hbits.next()
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>			if addr == 0 {
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>				break
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>			}
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>			i := (addr - p) / goarch.PtrSize
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>			tmpbuf[i/8] |= 1 &lt;&lt; (i % 8)
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		}
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	}
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	return bitvector{int32(nptr), &amp;tmpbuf[0]}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>}
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>
</pre><p><a href="heapdump.go?m=text">View as plain text</a></p>

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
