<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mwbbuf.go - Go Documentation Server</title>

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
<a href="mwbbuf.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mwbbuf.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2017 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This implements the write barrier buffer. The write barrier itself</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// is gcWriteBarrier and is implemented in assembly.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// See mbarrier.go for algorithmic details on the write barrier. This</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// file deals only with the buffer.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// The write barrier has a fast path and a slow path. The fast path</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// simply enqueues to a per-P write barrier buffer. It&#39;s written in</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// assembly and doesn&#39;t clobber any general purpose registers, so it</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// doesn&#39;t have the usual overheads of a Go call.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// When the buffer fills up, the write barrier invokes the slow path</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// (wbBufFlush) to flush the buffer to the GC work queues. In this</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// path, since the compiler didn&#39;t spill registers, we spill *all*</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// registers and disallow any GC safe points that could observe the</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// stack frame (since we don&#39;t know the types of the spilled</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// registers).</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>package runtime
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>import (
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>)
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// testSmallBuf forces a small write barrier buffer to stress write</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// barrier flushing.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>const testSmallBuf = false
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// wbBuf is a per-P buffer of pointers queued by the write barrier.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// This buffer is flushed to the GC workbufs when it fills up and on</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// various GC transitions.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// This is closely related to a &#34;sequential store buffer&#34; (SSB),</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// except that SSBs are usually used for maintaining remembered sets,</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// while this is used for marking.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>type wbBuf struct {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// next points to the next slot in buf. It must not be a</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// pointer type because it can point past the end of buf and</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// must be updated without write barriers.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// This is a pointer rather than an index to optimize the</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// write barrier assembly.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	next uintptr
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// end points to just past the end of buf. It must not be a</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// pointer type because it points past the end of buf and must</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// be updated without write barriers.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	end uintptr
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// buf stores a series of pointers to execute write barriers on.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	buf [wbBufEntries]uintptr
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>const (
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// wbBufEntries is the maximum number of pointers that can be</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// stored in the write barrier buffer.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// This trades latency for throughput amortization. Higher</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// values amortize flushing overhead more, but increase the</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// latency of flushing. Higher values also increase the cache</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// footprint of the buffer.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// TODO: What is the latency cost of this? Tune this value.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	wbBufEntries = 512
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// Maximum number of entries that we need to ask from the</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// buffer in a single call.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	wbMaxEntriesPerCall = 8
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// reset empties b by resetting its next and end pointers.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>func (b *wbBuf) reset() {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	start := uintptr(unsafe.Pointer(&amp;b.buf[0]))
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	b.next = start
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	if testSmallBuf {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		<span class="comment">// For testing, make the buffer smaller but more than</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		<span class="comment">// 1 write barrier&#39;s worth, so it tests both the</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		<span class="comment">// immediate flush and delayed flush cases.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		b.end = uintptr(unsafe.Pointer(&amp;b.buf[wbMaxEntriesPerCall+1]))
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	} else {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		b.end = start + uintptr(len(b.buf))*unsafe.Sizeof(b.buf[0])
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	if (b.end-b.next)%unsafe.Sizeof(b.buf[0]) != 0 {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		throw(&#34;bad write barrier buffer bounds&#34;)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// discard resets b&#39;s next pointer, but not its end pointer.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// This must be nosplit because it&#39;s called by wbBufFlush.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>func (b *wbBuf) discard() {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	b.next = uintptr(unsafe.Pointer(&amp;b.buf[0]))
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// empty reports whether b contains no pointers.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func (b *wbBuf) empty() bool {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	return b.next == uintptr(unsafe.Pointer(&amp;b.buf[0]))
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// getX returns space in the write barrier buffer to store X pointers.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// getX will flush the buffer if necessary. Callers should use this as:</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//	buf := &amp;getg().m.p.ptr().wbBuf</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">//	p := buf.get2()</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">//	p[0], p[1] = old, new</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">//	... actual memory write ...</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// The caller must ensure there are no preemption points during the</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// above sequence. There must be no preemption points while buf is in</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// use because it is a per-P resource. There must be no preemption</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// points between the buffer put and the write to memory because this</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// could allow a GC phase change, which could result in missed write</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// barriers.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// getX must be nowritebarrierrec to because write barriers here would</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// corrupt the write barrier buffer. It (and everything it calls, if</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// it called anything) has to be nosplit to avoid scheduling on to a</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// different P and a different buffer.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func (b *wbBuf) get1() *[1]uintptr {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	if b.next+goarch.PtrSize &gt; b.end {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		wbBufFlush()
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	p := (*[1]uintptr)(unsafe.Pointer(b.next))
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	b.next += goarch.PtrSize
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	return p
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>func (b *wbBuf) get2() *[2]uintptr {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	if b.next+2*goarch.PtrSize &gt; b.end {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		wbBufFlush()
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	p := (*[2]uintptr)(unsafe.Pointer(b.next))
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	b.next += 2 * goarch.PtrSize
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	return p
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// wbBufFlush flushes the current P&#39;s write barrier buffer to the GC</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">// workbufs.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// This must not have write barriers because it is part of the write</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// barrier implementation.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// This and everything it calls must be nosplit because 1) the stack</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">// contains untyped slots from gcWriteBarrier and 2) there must not be</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// a GC safe point between the write barrier test in the caller and</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// flushing the buffer.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// TODO: A &#34;go:nosplitrec&#34; annotation would be perfect for this.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>func wbBufFlush() {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// Note: Every possible return from this function must reset</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// the buffer&#39;s next pointer to prevent buffer overflow.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if getg().m.dying &gt; 0 {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re going down. Not much point in write barriers</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		<span class="comment">// and this way we can allow write barriers in the</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		<span class="comment">// panic path.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		getg().m.p.ptr().wbBuf.discard()
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		return
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// Switch to the system stack so we don&#39;t have to worry about</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// safe points.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		wbBufFlush1(getg().m.p.ptr())
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	})
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span><span class="comment">// wbBufFlush1 flushes p&#39;s write barrier buffer to the GC work queue.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// This must not have write barriers because it is part of the write</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// barrier implementation, so this may lead to infinite loops or</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// buffer corruption.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// This must be non-preemptible because it uses the P&#39;s workbuf.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>func wbBufFlush1(pp *p) {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// Get the buffered pointers.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	start := uintptr(unsafe.Pointer(&amp;pp.wbBuf.buf[0]))
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	n := (pp.wbBuf.next - start) / unsafe.Sizeof(pp.wbBuf.buf[0])
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	ptrs := pp.wbBuf.buf[:n]
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// Poison the buffer to make extra sure nothing is enqueued</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// while we&#39;re processing the buffer.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	pp.wbBuf.next = 0
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	if useCheckmark {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		<span class="comment">// Slow path for checkmark mode.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		for _, ptr := range ptrs {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			shade(ptr)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		pp.wbBuf.reset()
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		return
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// Mark all of the pointers in the buffer and record only the</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// pointers we greyed. We use the buffer itself to temporarily</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	<span class="comment">// record greyed pointers.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	<span class="comment">// TODO: Should scanobject/scanblock just stuff pointers into</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	<span class="comment">// the wbBuf? Then this would become the sole greying path.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// TODO: We could avoid shading any of the &#34;new&#34; pointers in</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">// the buffer if the stack has been shaded, or even avoid</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	<span class="comment">// putting them in the buffer at all (which would double its</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">// capacity). This is slightly complicated with the buffer; we</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	<span class="comment">// could track whether any un-shaded goroutine has used the</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// buffer, or just track globally whether there are any</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	<span class="comment">// un-shaded stacks and flush after each stack scan.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	gcw := &amp;pp.gcw
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	pos := 0
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	for _, ptr := range ptrs {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		if ptr &lt; minLegalPointer {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			<span class="comment">// nil pointers are very common, especially</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			<span class="comment">// for the &#34;old&#34; values. Filter out these and</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			<span class="comment">// other &#34;obvious&#34; non-heap pointers ASAP.</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			<span class="comment">// TODO: Should we filter out nils in the fast</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			<span class="comment">// path to reduce the rate of flushes?</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			continue
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		obj, span, objIndex := findObject(ptr, 0, 0)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		if obj == 0 {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			continue
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		<span class="comment">// TODO: Consider making two passes where the first</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		<span class="comment">// just prefetches the mark bits.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		mbits := span.markBitsForIndex(objIndex)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		if mbits.isMarked() {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			continue
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		mbits.setMarked()
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		<span class="comment">// Mark span.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		arena, pageIdx, pageMask := pageIndexOf(span.base())
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		if arena.pageMarks[pageIdx]&amp;pageMask == 0 {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			atomic.Or8(&amp;arena.pageMarks[pageIdx], pageMask)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		if span.spanclass.noscan() {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			gcw.bytesMarked += uint64(span.elemsize)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			continue
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		ptrs[pos] = obj
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		pos++
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	<span class="comment">// Enqueue the greyed objects.</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	gcw.putBatch(ptrs[:pos])
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	pp.wbBuf.reset()
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
</pre><p><a href="mwbbuf.go?m=text">View as plain text</a></p>

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
