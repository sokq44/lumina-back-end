<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/debuglog.go - Go Documentation Server</title>

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
<a href="debuglog.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">debuglog.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file provides an internal debug logging facility. The debug</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// log is a lightweight, in-memory, per-M ring buffer. By default, the</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// runtime prints the debug log on panic.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// To print something to the debug log, call dlog to obtain a dlogger</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// and use the methods on that to add values. The values will be</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// space-separated in the output (much like println).</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// This facility can be enabled by passing -tags debuglog when</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// building. Without this tag, dlog calls compile to nothing.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>package runtime
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>import (
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// debugLogBytes is the size of each per-M ring buffer. This is</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// allocated off-heap to avoid blowing up the M and hence the GC&#39;d</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// heap size.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>const debugLogBytes = 16 &lt;&lt; 10
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// debugLogStringLimit is the maximum number of bytes in a string.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// Above this, the string will be truncated with &#34;..(n more bytes)..&#34;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>const debugLogStringLimit = debugLogBytes / 8
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// dlog returns a debug logger. The caller can use methods on the</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// returned logger to add values, which will be space-separated in the</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// final output, much like println. The caller must call end() to</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// finish the message.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// dlog can be used from highly-constrained corners of the runtime: it</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// is safe to use in the signal handler, from within the write</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// barrier, from within the stack implementation, and in places that</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// must be recursively nosplit.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// This will be compiled away if built without the debuglog build tag.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// However, argument construction may not be. If any of the arguments</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// are not literals or trivial expressions, consider protecting the</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// call with &#34;if dlogEnabled&#34;.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>func dlog() *dlogger {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		return nil
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// Get the time.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	tick, nano := uint64(cputicks()), uint64(nanotime())
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// Try to get a cached logger.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	l := getCachedDlogger()
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// If we couldn&#39;t get a cached logger, try to get one from the</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// global pool.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if l == nil {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		allp := (*uintptr)(unsafe.Pointer(&amp;allDloggers))
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		all := (*dlogger)(unsafe.Pointer(atomic.Loaduintptr(allp)))
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		for l1 := all; l1 != nil; l1 = l1.allLink {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			if l1.owned.Load() == 0 &amp;&amp; l1.owned.CompareAndSwap(0, 1) {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>				l = l1
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>				break
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// If that failed, allocate a new logger.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	if l == nil {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		<span class="comment">// Use sysAllocOS instead of sysAlloc because we want to interfere</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		<span class="comment">// with the runtime as little as possible, and sysAlloc updates accounting.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		l = (*dlogger)(sysAllocOS(unsafe.Sizeof(dlogger{})))
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		if l == nil {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			throw(&#34;failed to allocate debug log&#34;)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		l.w.r.data = &amp;l.w.data
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		l.owned.Store(1)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		<span class="comment">// Prepend to allDloggers list.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		headp := (*uintptr)(unsafe.Pointer(&amp;allDloggers))
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		for {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			head := atomic.Loaduintptr(headp)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			l.allLink = (*dlogger)(unsafe.Pointer(head))
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			if atomic.Casuintptr(headp, head, uintptr(unsafe.Pointer(l))) {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>				break
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// If the time delta is getting too high, write a new sync</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// packet. We set the limit so we don&#39;t write more than 6</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// bytes of delta in the record header.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	const deltaLimit = 1&lt;&lt;(3*7) - 1 <span class="comment">// ~2ms between sync packets</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	if tick-l.w.tick &gt; deltaLimit || nano-l.w.nano &gt; deltaLimit {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		l.w.writeSync(tick, nano)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// Reserve space for framing header.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	l.w.ensure(debugLogHeaderSize)
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	l.w.write += debugLogHeaderSize
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">// Write record header.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	l.w.uvarint(tick - l.w.tick)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	l.w.uvarint(nano - l.w.nano)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	gp := getg()
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	if gp != nil &amp;&amp; gp.m != nil &amp;&amp; gp.m.p != 0 {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		l.w.varint(int64(gp.m.p.ptr().id))
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	} else {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		l.w.varint(-1)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	return l
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// A dlogger writes to the debug log.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// To obtain a dlogger, call dlog(). When done with the dlogger, call</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// end().</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>type dlogger struct {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	w debugLogWriter
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// allLink is the next dlogger in the allDloggers list.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	allLink *dlogger
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// owned indicates that this dlogger is owned by an M. This is</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// accessed atomically.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	owned atomic.Uint32
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// allDloggers is a list of all dloggers, linked through</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">// dlogger.allLink. This is accessed atomically. This is prepend only,</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// so it doesn&#39;t need to protect against ABA races.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>var allDloggers *dlogger
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>func (l *dlogger) end() {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		return
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// Fill in framing header.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	size := l.w.write - l.w.r.end
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	if !l.w.writeFrameAt(l.w.r.end, size) {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		throw(&#34;record too large&#34;)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// Commit the record.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	l.w.r.end = l.w.write
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// Attempt to return this logger to the cache.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	if putCachedDlogger(l) {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		return
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// Return the logger to the global pool.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	l.owned.Store(0)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>const (
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	debugLogUnknown = 1 + iota
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	debugLogBoolTrue
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	debugLogBoolFalse
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	debugLogInt
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	debugLogUint
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	debugLogHex
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	debugLogPtr
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	debugLogString
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	debugLogConstString
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	debugLogStringOverflow
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	debugLogPC
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	debugLogTraceback
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>func (l *dlogger) b(x bool) *dlogger {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		return l
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	if x {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		l.w.byte(debugLogBoolTrue)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	} else {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		l.w.byte(debugLogBoolFalse)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	return l
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>func (l *dlogger) i(x int) *dlogger {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	return l.i64(int64(x))
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>func (l *dlogger) i8(x int8) *dlogger {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	return l.i64(int64(x))
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>func (l *dlogger) i16(x int16) *dlogger {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	return l.i64(int64(x))
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>func (l *dlogger) i32(x int32) *dlogger {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	return l.i64(int64(x))
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>func (l *dlogger) i64(x int64) *dlogger {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		return l
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	l.w.byte(debugLogInt)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	l.w.varint(x)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	return l
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>func (l *dlogger) u(x uint) *dlogger {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	return l.u64(uint64(x))
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>func (l *dlogger) uptr(x uintptr) *dlogger {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	return l.u64(uint64(x))
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>func (l *dlogger) u8(x uint8) *dlogger {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	return l.u64(uint64(x))
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>func (l *dlogger) u16(x uint16) *dlogger {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	return l.u64(uint64(x))
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>func (l *dlogger) u32(x uint32) *dlogger {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	return l.u64(uint64(x))
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func (l *dlogger) u64(x uint64) *dlogger {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		return l
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	l.w.byte(debugLogUint)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	l.w.uvarint(x)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	return l
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>func (l *dlogger) hex(x uint64) *dlogger {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		return l
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	l.w.byte(debugLogHex)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	l.w.uvarint(x)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	return l
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>}
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>func (l *dlogger) p(x any) *dlogger {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		return l
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	l.w.byte(debugLogPtr)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	if x == nil {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		l.w.uvarint(0)
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	} else {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		v := efaceOf(&amp;x)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		switch v._type.Kind_ &amp; kindMask {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		case kindChan, kindFunc, kindMap, kindPtr, kindUnsafePointer:
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			l.w.uvarint(uint64(uintptr(v.data)))
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		default:
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			throw(&#34;not a pointer type&#34;)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	return l
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>func (l *dlogger) s(x string) *dlogger {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		return l
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	strData := unsafe.StringData(x)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	datap := &amp;firstmoduledata
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	if len(x) &gt; 4 &amp;&amp; datap.etext &lt;= uintptr(unsafe.Pointer(strData)) &amp;&amp; uintptr(unsafe.Pointer(strData)) &lt; datap.end {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		<span class="comment">// String constants are in the rodata section, which</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		<span class="comment">// isn&#39;t recorded in moduledata. But it has to be</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		<span class="comment">// somewhere between etext and end.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		l.w.byte(debugLogConstString)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		l.w.uvarint(uint64(len(x)))
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		l.w.uvarint(uint64(uintptr(unsafe.Pointer(strData)) - datap.etext))
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	} else {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		l.w.byte(debugLogString)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		<span class="comment">// We can&#39;t use unsafe.Slice as it may panic, which isn&#39;t safe</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		<span class="comment">// in this (potentially) nowritebarrier context.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		var b []byte
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		bb := (*slice)(unsafe.Pointer(&amp;b))
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		bb.array = unsafe.Pointer(strData)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		bb.len, bb.cap = len(x), len(x)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		if len(b) &gt; debugLogStringLimit {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>			b = b[:debugLogStringLimit]
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		l.w.uvarint(uint64(len(b)))
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		l.w.bytes(b)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		if len(b) != len(x) {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			l.w.byte(debugLogStringOverflow)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			l.w.uvarint(uint64(len(x) - len(b)))
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	return l
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>func (l *dlogger) pc(x uintptr) *dlogger {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		return l
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	l.w.byte(debugLogPC)
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	l.w.uvarint(uint64(x))
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	return l
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>func (l *dlogger) traceback(x []uintptr) *dlogger {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		return l
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	l.w.byte(debugLogTraceback)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	l.w.uvarint(uint64(len(x)))
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	for _, pc := range x {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		l.w.uvarint(uint64(pc))
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	return l
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// A debugLogWriter is a ring buffer of binary debug log records.</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// A log record consists of a 2-byte framing header and a sequence of</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// fields. The framing header gives the size of the record as a little</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// endian 16-bit value. Each field starts with a byte indicating its</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// type, followed by type-specific data. If the size in the framing</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// header is 0, it&#39;s a sync record consisting of two little endian</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// 64-bit values giving a new time base.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">// Because this is a ring buffer, new records will eventually</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// overwrite old records. Hence, it maintains a reader that consumes</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span><span class="comment">// the log as it gets overwritten. That reader state is where an</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// actual log reader would start.</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>type debugLogWriter struct {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	_     sys.NotInHeap
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	write uint64
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	data  debugLogBuf
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	<span class="comment">// tick and nano are the time bases from the most recently</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	<span class="comment">// written sync record.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	tick, nano uint64
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	<span class="comment">// r is a reader that consumes records as they get overwritten</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	<span class="comment">// by the writer. It also acts as the initial reader state</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	<span class="comment">// when printing the log.</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	r debugLogReader
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	<span class="comment">// buf is a scratch buffer for encoding. This is here to</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	<span class="comment">// reduce stack usage.</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	buf [10]byte
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>type debugLogBuf struct {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	b [debugLogBytes]byte
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>const (
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	<span class="comment">// debugLogHeaderSize is the number of bytes in the framing</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	<span class="comment">// header of every dlog record.</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	debugLogHeaderSize = 2
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	<span class="comment">// debugLogSyncSize is the number of bytes in a sync record.</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	debugLogSyncSize = debugLogHeaderSize + 2*8
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>func (l *debugLogWriter) ensure(n uint64) {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	for l.write+n &gt;= l.r.begin+uint64(len(l.data.b)) {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		<span class="comment">// Consume record at begin.</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		if l.r.skip() == ^uint64(0) {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			<span class="comment">// Wrapped around within a record.</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			<span class="comment">// TODO(austin): It would be better to just</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			<span class="comment">// eat the whole buffer at this point, but we</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			<span class="comment">// have to communicate that to the reader</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			<span class="comment">// somehow.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			throw(&#34;record wrapped around&#34;)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>func (l *debugLogWriter) writeFrameAt(pos, size uint64) bool {
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	l.data.b[pos%uint64(len(l.data.b))] = uint8(size)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	l.data.b[(pos+1)%uint64(len(l.data.b))] = uint8(size &gt;&gt; 8)
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	return size &lt;= 0xFFFF
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>func (l *debugLogWriter) writeSync(tick, nano uint64) {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	l.tick, l.nano = tick, nano
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	l.ensure(debugLogHeaderSize)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	l.writeFrameAt(l.write, 0)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	l.write += debugLogHeaderSize
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	l.writeUint64LE(tick)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	l.writeUint64LE(nano)
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	l.r.end = l.write
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>func (l *debugLogWriter) writeUint64LE(x uint64) {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	var b [8]byte
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	b[0] = byte(x)
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	b[1] = byte(x &gt;&gt; 8)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	b[2] = byte(x &gt;&gt; 16)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	b[3] = byte(x &gt;&gt; 24)
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	b[4] = byte(x &gt;&gt; 32)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	b[5] = byte(x &gt;&gt; 40)
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	b[6] = byte(x &gt;&gt; 48)
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	b[7] = byte(x &gt;&gt; 56)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	l.bytes(b[:])
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>func (l *debugLogWriter) byte(x byte) {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	l.ensure(1)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	pos := l.write
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	l.write++
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	l.data.b[pos%uint64(len(l.data.b))] = x
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>func (l *debugLogWriter) bytes(x []byte) {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	l.ensure(uint64(len(x)))
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	pos := l.write
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	l.write += uint64(len(x))
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	for len(x) &gt; 0 {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		n := copy(l.data.b[pos%uint64(len(l.data.b)):], x)
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		pos += uint64(n)
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		x = x[n:]
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	}
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>func (l *debugLogWriter) varint(x int64) {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	var u uint64
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		u = (^uint64(x) &lt;&lt; 1) | 1 <span class="comment">// complement i, bit 0 is 1</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	} else {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		u = (uint64(x) &lt;&lt; 1) <span class="comment">// do not complement i, bit 0 is 0</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	l.uvarint(u)
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>func (l *debugLogWriter) uvarint(u uint64) {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	i := 0
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	for u &gt;= 0x80 {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		l.buf[i] = byte(u) | 0x80
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		u &gt;&gt;= 7
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		i++
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	l.buf[i] = byte(u)
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	i++
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	l.bytes(l.buf[:i])
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>type debugLogReader struct {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	data *debugLogBuf
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	<span class="comment">// begin and end are the positions in the log of the beginning</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	<span class="comment">// and end of the log data, modulo len(data).</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	begin, end uint64
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	<span class="comment">// tick and nano are the current time base at begin.</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	tick, nano uint64
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>func (r *debugLogReader) skip() uint64 {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	<span class="comment">// Read size at pos.</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	if r.begin+debugLogHeaderSize &gt; r.end {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		return ^uint64(0)
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	size := uint64(r.readUint16LEAt(r.begin))
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	if size == 0 {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		<span class="comment">// Sync packet.</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		r.tick = r.readUint64LEAt(r.begin + debugLogHeaderSize)
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		r.nano = r.readUint64LEAt(r.begin + debugLogHeaderSize + 8)
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		size = debugLogSyncSize
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	if r.begin+size &gt; r.end {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		return ^uint64(0)
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	r.begin += size
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	return size
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>func (r *debugLogReader) readUint16LEAt(pos uint64) uint16 {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	return uint16(r.data.b[pos%uint64(len(r.data.b))]) |
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		uint16(r.data.b[(pos+1)%uint64(len(r.data.b))])&lt;&lt;8
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>func (r *debugLogReader) readUint64LEAt(pos uint64) uint64 {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	var b [8]byte
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	for i := range b {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		b[i] = r.data.b[pos%uint64(len(r.data.b))]
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		pos++
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	}
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	return uint64(b[0]) | uint64(b[1])&lt;&lt;8 |
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		uint64(b[2])&lt;&lt;16 | uint64(b[3])&lt;&lt;24 |
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		uint64(b[4])&lt;&lt;32 | uint64(b[5])&lt;&lt;40 |
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		uint64(b[6])&lt;&lt;48 | uint64(b[7])&lt;&lt;56
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>func (r *debugLogReader) peek() (tick uint64) {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	<span class="comment">// Consume any sync records.</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	size := uint64(0)
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	for size == 0 {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		if r.begin+debugLogHeaderSize &gt; r.end {
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			return ^uint64(0)
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		size = uint64(r.readUint16LEAt(r.begin))
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		if size != 0 {
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>			break
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		if r.begin+debugLogSyncSize &gt; r.end {
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>			return ^uint64(0)
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		<span class="comment">// Sync packet.</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		r.tick = r.readUint64LEAt(r.begin + debugLogHeaderSize)
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		r.nano = r.readUint64LEAt(r.begin + debugLogHeaderSize + 8)
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		r.begin += debugLogSyncSize
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	}
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	<span class="comment">// Peek tick delta.</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	if r.begin+size &gt; r.end {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		return ^uint64(0)
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	pos := r.begin + debugLogHeaderSize
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	var u uint64
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	for i := uint(0); ; i += 7 {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		b := r.data.b[pos%uint64(len(r.data.b))]
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		pos++
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		u |= uint64(b&amp;^0x80) &lt;&lt; i
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		if b&amp;0x80 == 0 {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>			break
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	if pos &gt; r.begin+size {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		return ^uint64(0)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	return r.tick + u
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>}
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>func (r *debugLogReader) header() (end, tick, nano uint64, p int) {
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	<span class="comment">// Read size. We&#39;ve already skipped sync packets and checked</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	<span class="comment">// bounds in peek.</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	size := uint64(r.readUint16LEAt(r.begin))
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	end = r.begin + size
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	r.begin += debugLogHeaderSize
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	<span class="comment">// Read tick, nano, and p.</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	tick = r.uvarint() + r.tick
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	nano = r.uvarint() + r.nano
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	p = int(r.varint())
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	return
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>}
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>func (r *debugLogReader) uvarint() uint64 {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	var u uint64
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	for i := uint(0); ; i += 7 {
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		b := r.data.b[r.begin%uint64(len(r.data.b))]
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		r.begin++
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		u |= uint64(b&amp;^0x80) &lt;&lt; i
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		if b&amp;0x80 == 0 {
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>			break
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	return u
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>}
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>func (r *debugLogReader) varint() int64 {
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	u := r.uvarint()
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	var v int64
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	if u&amp;1 == 0 {
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		v = int64(u &gt;&gt; 1)
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	} else {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		v = ^int64(u &gt;&gt; 1)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	return v
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>}
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>func (r *debugLogReader) printVal() bool {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	typ := r.data.b[r.begin%uint64(len(r.data.b))]
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	r.begin++
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	switch typ {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	default:
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		print(&#34;&lt;unknown field type &#34;, hex(typ), &#34; pos &#34;, r.begin-1, &#34; end &#34;, r.end, &#34;&gt;\n&#34;)
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		return false
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	case debugLogUnknown:
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		print(&#34;&lt;unknown kind&gt;&#34;)
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	case debugLogBoolTrue:
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		print(true)
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	case debugLogBoolFalse:
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>		print(false)
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	case debugLogInt:
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		print(r.varint())
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	case debugLogUint:
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		print(r.uvarint())
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	case debugLogHex, debugLogPtr:
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		print(hex(r.uvarint()))
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	case debugLogString:
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		sl := r.uvarint()
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		if r.begin+sl &gt; r.end {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>			r.begin = r.end
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>			print(&#34;&lt;string length corrupted&gt;&#34;)
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>			break
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		}
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		for sl &gt; 0 {
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>			b := r.data.b[r.begin%uint64(len(r.data.b)):]
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>			if uint64(len(b)) &gt; sl {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>				b = b[:sl]
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>			}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>			r.begin += uint64(len(b))
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>			sl -= uint64(len(b))
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>			gwrite(b)
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	case debugLogConstString:
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		len, ptr := int(r.uvarint()), uintptr(r.uvarint())
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		ptr += firstmoduledata.etext
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		<span class="comment">// We can&#39;t use unsafe.String as it may panic, which isn&#39;t safe</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		<span class="comment">// in this (potentially) nowritebarrier context.</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		str := stringStruct{
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>			str: unsafe.Pointer(ptr),
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>			len: len,
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		}
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		s := *(*string)(unsafe.Pointer(&amp;str))
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		print(s)
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	case debugLogStringOverflow:
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		print(&#34;..(&#34;, r.uvarint(), &#34; more bytes)..&#34;)
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	case debugLogPC:
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		printDebugLogPC(uintptr(r.uvarint()), false)
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	case debugLogTraceback:
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		n := int(r.uvarint())
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		for i := 0; i &lt; n; i++ {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>			print(&#34;\n\t&#34;)
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			<span class="comment">// gentraceback PCs are always return PCs.</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>			<span class="comment">// Convert them to call PCs.</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>			<span class="comment">// TODO(austin): Expand inlined frames.</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>			printDebugLogPC(uintptr(r.uvarint()), true)
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		}
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	return true
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>}
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span><span class="comment">// printDebugLog prints the debug log.</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>func printDebugLog() {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	if !dlogEnabled {
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		return
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	<span class="comment">// This function should not panic or throw since it is used in</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	<span class="comment">// the fatal panic path and this may deadlock.</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	printlock()
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	<span class="comment">// Get the list of all debug logs.</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	allp := (*uintptr)(unsafe.Pointer(&amp;allDloggers))
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	all := (*dlogger)(unsafe.Pointer(atomic.Loaduintptr(allp)))
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	<span class="comment">// Count the logs.</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	n := 0
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	for l := all; l != nil; l = l.allLink {
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>		n++
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	}
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	if n == 0 {
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		printunlock()
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		return
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	}
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	<span class="comment">// Prepare read state for all logs.</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	type readState struct {
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		debugLogReader
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		first    bool
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		lost     uint64
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		nextTick uint64
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	}
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	<span class="comment">// Use sysAllocOS instead of sysAlloc because we want to interfere</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	<span class="comment">// with the runtime as little as possible, and sysAlloc updates accounting.</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>	state1 := sysAllocOS(unsafe.Sizeof(readState{}) * uintptr(n))
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	if state1 == nil {
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		println(&#34;failed to allocate read state for&#34;, n, &#34;logs&#34;)
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		printunlock()
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>		return
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	}
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	state := (*[1 &lt;&lt; 20]readState)(state1)[:n]
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	{
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>		l := all
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>		for i := range state {
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			s := &amp;state[i]
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>			s.debugLogReader = l.w.r
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>			s.first = true
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>			s.lost = l.w.r.begin
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>			s.nextTick = s.peek()
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			l = l.allLink
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>		}
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	}
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	<span class="comment">// Print records.</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	for {
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		<span class="comment">// Find the next record.</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		var best struct {
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>			tick uint64
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>			i    int
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>		}
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>		best.tick = ^uint64(0)
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		for i := range state {
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>			if state[i].nextTick &lt; best.tick {
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>				best.tick = state[i].nextTick
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>				best.i = i
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>			}
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		}
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		if best.tick == ^uint64(0) {
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>			break
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		}
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>		<span class="comment">// Print record.</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		s := &amp;state[best.i]
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		if s.first {
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>			print(&#34;&gt;&gt; begin log &#34;, best.i)
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>			if s.lost != 0 {
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>				print(&#34;; lost first &#34;, s.lost&gt;&gt;10, &#34;KB&#34;)
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>			}
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			print(&#34; &lt;&lt;\n&#34;)
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>			s.first = false
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		}
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		end, _, nano, p := s.header()
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		oldEnd := s.end
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>		s.end = end
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		print(&#34;[&#34;)
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		var tmpbuf [21]byte
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>		pnano := int64(nano) - runtimeInitTime
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>		if pnano &lt; 0 {
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>			<span class="comment">// Logged before runtimeInitTime was set.</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>			pnano = 0
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>		}
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		pnanoBytes := itoaDiv(tmpbuf[:], uint64(pnano), 9)
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		print(slicebytetostringtmp((*byte)(noescape(unsafe.Pointer(&amp;pnanoBytes[0]))), len(pnanoBytes)))
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>		print(&#34; P &#34;, p, &#34;] &#34;)
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>		for i := 0; s.begin &lt; s.end; i++ {
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>			if i &gt; 0 {
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>				print(&#34; &#34;)
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>			}
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>			if !s.printVal() {
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>				<span class="comment">// Abort this P log.</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>				print(&#34;&lt;aborting P log&gt;&#34;)
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>				end = oldEnd
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>				break
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>			}
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>		}
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		println()
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		<span class="comment">// Move on to the next record.</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>		s.begin = end
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>		s.end = oldEnd
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		s.nextTick = s.peek()
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	}
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	printunlock()
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>}
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span><span class="comment">// printDebugLogPC prints a single symbolized PC. If returnPC is true,</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span><span class="comment">// pc is a return PC that must first be converted to a call PC.</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>func printDebugLogPC(pc uintptr, returnPC bool) {
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	fn := findfunc(pc)
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	if returnPC &amp;&amp; (!fn.valid() || pc &gt; fn.entry()) {
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>		<span class="comment">// TODO(austin): Don&#39;t back up if the previous frame</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		<span class="comment">// was a sigpanic.</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		pc--
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	}
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	print(hex(pc))
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	if !fn.valid() {
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		print(&#34; [unknown PC]&#34;)
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	} else {
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		name := funcname(fn)
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		file, line := funcline(fn, pc)
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>		print(&#34; [&#34;, name, &#34;+&#34;, hex(pc-fn.entry()),
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>			&#34; &#34;, file, &#34;:&#34;, line, &#34;]&#34;)
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	}
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>}
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>
</pre><p><a href="debuglog.go?m=text">View as plain text</a></p>

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
