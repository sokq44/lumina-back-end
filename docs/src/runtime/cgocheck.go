<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/cgocheck.go - Go Documentation Server</title>

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
<a href="cgocheck.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">cgocheck.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2015 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Code to check that pointer writes follow the cgo rules.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// These functions are invoked when GOEXPERIMENT=cgocheck2 is enabled.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>package runtime
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>import (
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/goexperiment&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>const cgoWriteBarrierFail = &#34;unpinned Go pointer stored into non-Go memory&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// cgoCheckPtrWrite is called whenever a pointer is stored into memory.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// It throws if the program is storing an unpinned Go pointer into non-Go</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// memory.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// This is called from generated code when GOEXPERIMENT=cgocheck2 is enabled.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>func cgoCheckPtrWrite(dst *unsafe.Pointer, src unsafe.Pointer) {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	if !mainStarted {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		<span class="comment">// Something early in startup hates this function.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		<span class="comment">// Don&#39;t start doing any actual checking until the</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		<span class="comment">// runtime has set itself up.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		return
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	if !cgoIsGoPointer(src) {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		return
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	if cgoIsGoPointer(unsafe.Pointer(dst)) {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		return
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// If we are running on the system stack then dst might be an</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// address on the stack, which is OK.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	gp := getg()
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	if gp == gp.m.g0 || gp == gp.m.gsignal {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		return
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// Allocating memory can write to various mfixalloc structs</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// that look like they are non-Go memory.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	if gp.m.mallocing != 0 {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		return
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// If the object is pinned, it&#39;s safe to store it in C memory. The GC</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// ensures it will not be moved or freed.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	if isPinned(src) {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		return
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s OK if writing to memory allocated by persistentalloc.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// Do this check last because it is more expensive and rarely true.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// If it is false the expense doesn&#39;t matter since we are crashing.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	if inPersistentAlloc(uintptr(unsafe.Pointer(dst))) {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		return
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		println(&#34;write of unpinned Go pointer&#34;, hex(uintptr(src)), &#34;to non-Go memory&#34;, hex(uintptr(unsafe.Pointer(dst))))
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		throw(cgoWriteBarrierFail)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	})
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// cgoCheckMemmove is called when moving a block of memory.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// It throws if the program is copying a block that contains an unpinned Go</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// pointer into non-Go memory.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// This is called from generated code when GOEXPERIMENT=cgocheck2 is enabled.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>func cgoCheckMemmove(typ *_type, dst, src unsafe.Pointer) {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	cgoCheckMemmove2(typ, dst, src, 0, typ.Size_)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// cgoCheckMemmove2 is called when moving a block of memory.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// dst and src point off bytes into the value to copy.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// size is the number of bytes to copy.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// It throws if the program is copying a block that contains an unpinned Go</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// pointer into non-Go memory.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>func cgoCheckMemmove2(typ *_type, dst, src unsafe.Pointer, off, size uintptr) {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	if typ.PtrBytes == 0 {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		return
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if !cgoIsGoPointer(src) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		return
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	if cgoIsGoPointer(dst) {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		return
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	cgoCheckTypedBlock(typ, src, off, size)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// cgoCheckSliceCopy is called when copying n elements of a slice.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// src and dst are pointers to the first element of the slice.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// typ is the element type of the slice.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// It throws if the program is copying slice elements that contain unpinned Go</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// pointers into non-Go memory.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func cgoCheckSliceCopy(typ *_type, dst, src unsafe.Pointer, n int) {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	if typ.PtrBytes == 0 {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		return
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	if !cgoIsGoPointer(src) {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		return
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	if cgoIsGoPointer(dst) {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		return
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	p := src
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		cgoCheckTypedBlock(typ, p, 0, typ.Size_)
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		p = add(p, typ.Size_)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// cgoCheckTypedBlock checks the block of memory at src, for up to size bytes,</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">// and throws if it finds an unpinned Go pointer. The type of the memory is typ,</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// and src is off bytes into that type.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>func cgoCheckTypedBlock(typ *_type, src unsafe.Pointer, off, size uintptr) {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// Anything past typ.PtrBytes is not a pointer.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	if typ.PtrBytes &lt;= off {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		return
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	if ptrdataSize := typ.PtrBytes - off; size &gt; ptrdataSize {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		size = ptrdataSize
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	if typ.Kind_&amp;kindGCProg == 0 {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		cgoCheckBits(src, typ.GCData, off, size)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		return
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// The type has a GC program. Try to find GC bits somewhere else.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	for _, datap := range activeModules() {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		if cgoInRange(src, datap.data, datap.edata) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			doff := uintptr(src) - datap.data
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			cgoCheckBits(add(src, -doff), datap.gcdatamask.bytedata, off+doff, size)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			return
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		if cgoInRange(src, datap.bss, datap.ebss) {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			boff := uintptr(src) - datap.bss
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			cgoCheckBits(add(src, -boff), datap.gcbssmask.bytedata, off+boff, size)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			return
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	s := spanOfUnchecked(uintptr(src))
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	if s.state.get() == mSpanManual {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		<span class="comment">// There are no heap bits for value stored on the stack.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		<span class="comment">// For a channel receive src might be on the stack of some</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// other goroutine, so we can&#39;t unwind the stack even if</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		<span class="comment">// we wanted to.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		<span class="comment">// We can&#39;t expand the GC program without extra storage</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		<span class="comment">// space we can&#39;t easily get.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		<span class="comment">// Fortunately we have the type information.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			cgoCheckUsingType(typ, src, off, size)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		})
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		return
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// src must be in the regular heap.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	if goexperiment.AllocHeaders {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		tp := s.typePointersOf(uintptr(src), size)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		for {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			var addr uintptr
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			if tp, addr = tp.next(uintptr(src) + size); addr == 0 {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>				break
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			v := *(*unsafe.Pointer)(unsafe.Pointer(addr))
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			if cgoIsGoPointer(v) &amp;&amp; !isPinned(v) {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				throw(cgoWriteBarrierFail)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	} else {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		hbits := heapBitsForAddr(uintptr(src), size)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		for {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			var addr uintptr
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			if hbits, addr = hbits.next(); addr == 0 {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>				break
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			v := *(*unsafe.Pointer)(unsafe.Pointer(addr))
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			if cgoIsGoPointer(v) &amp;&amp; !isPinned(v) {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>				throw(cgoWriteBarrierFail)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>			}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// cgoCheckBits checks the block of memory at src, for up to size</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// bytes, and throws if it finds an unpinned Go pointer. The gcbits mark each</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// pointer value. The src pointer is off bytes into the gcbits.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>func cgoCheckBits(src unsafe.Pointer, gcbits *byte, off, size uintptr) {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	skipMask := off / goarch.PtrSize / 8
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	skipBytes := skipMask * goarch.PtrSize * 8
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	ptrmask := addb(gcbits, skipMask)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	src = add(src, skipBytes)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	off -= skipBytes
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	size += off
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	var bits uint32
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; size; i += goarch.PtrSize {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		if i&amp;(goarch.PtrSize*8-1) == 0 {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			bits = uint32(*ptrmask)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			ptrmask = addb(ptrmask, 1)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		} else {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			bits &gt;&gt;= 1
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		if off &gt; 0 {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			off -= goarch.PtrSize
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		} else {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			if bits&amp;1 != 0 {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>				v := *(*unsafe.Pointer)(add(src, i))
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>				if cgoIsGoPointer(v) &amp;&amp; !isPinned(v) {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>					throw(cgoWriteBarrierFail)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>				}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span><span class="comment">// cgoCheckUsingType is like cgoCheckTypedBlock, but is a last ditch</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span><span class="comment">// fall back to look for pointers in src using the type information.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span><span class="comment">// We only use this when looking at a value on the stack when the type</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">// uses a GC program, because otherwise it&#39;s more efficient to use the</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// GC bits. This is called on the system stack.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>func cgoCheckUsingType(typ *_type, src unsafe.Pointer, off, size uintptr) {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	if typ.PtrBytes == 0 {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		return
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	<span class="comment">// Anything past typ.PtrBytes is not a pointer.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	if typ.PtrBytes &lt;= off {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		return
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	if ptrdataSize := typ.PtrBytes - off; size &gt; ptrdataSize {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		size = ptrdataSize
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	if typ.Kind_&amp;kindGCProg == 0 {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		cgoCheckBits(src, typ.GCData, off, size)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		return
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	switch typ.Kind_ &amp; kindMask {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	default:
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		throw(&#34;can&#39;t happen&#34;)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	case kindArray:
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		at := (*arraytype)(unsafe.Pointer(typ))
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; at.Len; i++ {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			if off &lt; at.Elem.Size_ {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>				cgoCheckUsingType(at.Elem, src, off, size)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			src = add(src, at.Elem.Size_)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			skipped := off
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			if skipped &gt; at.Elem.Size_ {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>				skipped = at.Elem.Size_
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			checked := at.Elem.Size_ - skipped
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			off -= skipped
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			if size &lt;= checked {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>				return
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			size -= checked
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	case kindStruct:
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		st := (*structtype)(unsafe.Pointer(typ))
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		for _, f := range st.Fields {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			if off &lt; f.Typ.Size_ {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>				cgoCheckUsingType(f.Typ, src, off, size)
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			src = add(src, f.Typ.Size_)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			skipped := off
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			if skipped &gt; f.Typ.Size_ {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>				skipped = f.Typ.Size_
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			checked := f.Typ.Size_ - skipped
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			off -= skipped
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			if size &lt;= checked {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>				return
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			size -= checked
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
</pre><p><a href="cgocheck.go?m=text">View as plain text</a></p>

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
