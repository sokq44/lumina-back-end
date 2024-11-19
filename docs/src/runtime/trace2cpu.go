<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/trace2cpu.go - Go Documentation Server</title>

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
<a href="trace2cpu.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">trace2cpu.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// CPU profile -&gt; trace</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package runtime
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// traceInitReadCPU initializes CPU profile -&gt; tracer state for tracing.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// Returns a profBuf for reading from.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>func traceInitReadCPU() {
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	if traceEnabled() {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>		throw(&#34;traceInitReadCPU called with trace enabled&#34;)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	}
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// Create new profBuf for CPU samples that will be emitted as events.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// Format: after the timestamp, header is [pp.id, gp.goid, mp.procid].</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	trace.cpuLogRead[0] = newProfBuf(3, profBufWordCount, profBufTagCount)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	trace.cpuLogRead[1] = newProfBuf(3, profBufWordCount, profBufTagCount)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// We must not acquire trace.signalLock outside of a signal handler: a</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// profiling signal may arrive at any time and try to acquire it, leading to</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// deadlock. Because we can&#39;t use that lock to protect updates to</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// trace.cpuLogWrite (only use of the structure it references), reads and</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// writes of the pointer must be atomic. (And although this field is never</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// the sole pointer to the profBuf value, it&#39;s best to allow a write barrier</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// here.)</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	trace.cpuLogWrite[0].Store(trace.cpuLogRead[0])
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	trace.cpuLogWrite[1].Store(trace.cpuLogRead[1])
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// traceStartReadCPU creates a goroutine to start reading CPU profile</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// data into an active trace.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// traceAdvanceSema must be held.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>func traceStartReadCPU() {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	if !traceEnabled() {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		throw(&#34;traceStartReadCPU called with trace disabled&#34;)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// Spin up the logger goroutine.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	trace.cpuSleep = newWakeableSleep()
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	done := make(chan struct{}, 1)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	go func() {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		for traceEnabled() {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			<span class="comment">// Sleep here because traceReadCPU is non-blocking. This mirrors</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			<span class="comment">// how the runtime/pprof package obtains CPU profile data.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			<span class="comment">// We can&#39;t do a blocking read here because Darwin can&#39;t do a</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			<span class="comment">// wakeup from a signal handler, so all CPU profiling is just</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			<span class="comment">// non-blocking. See #61768 for more details.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			<span class="comment">// Like the runtime/pprof package, even if that bug didn&#39;t exist</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			<span class="comment">// we would still want to do a goroutine-level sleep in between</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>			<span class="comment">// reads to avoid frequent wakeups.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			trace.cpuSleep.sleep(100_000_000)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			tl := traceAcquire()
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			if !tl.ok() {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>				<span class="comment">// Tracing disabled.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>				break
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			keepGoing := traceReadCPU(tl.gen)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			traceRelease(tl)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			if !keepGoing {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>				break
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		done &lt;- struct{}{}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}()
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	trace.cpuLogDone = done
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// traceStopReadCPU blocks until the trace CPU reading goroutine exits.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// traceAdvanceSema must be held, and tracing must be disabled.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>func traceStopReadCPU() {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	if traceEnabled() {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		throw(&#34;traceStopReadCPU called with trace enabled&#34;)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// Once we close the profbuf, we&#39;ll be in one of two situations:</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// - The logger goroutine has already exited because it observed</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">//   that the trace is disabled.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// - The logger goroutine is asleep.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// Wake the goroutine so it can observe that their the buffer is</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// closed an exit.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	trace.cpuLogWrite[0].Store(nil)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	trace.cpuLogWrite[1].Store(nil)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	trace.cpuLogRead[0].close()
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	trace.cpuLogRead[1].close()
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	trace.cpuSleep.wake()
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// Wait until the logger goroutine exits.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	&lt;-trace.cpuLogDone
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// Clear state for the next trace.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	trace.cpuLogDone = nil
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	trace.cpuLogRead[0] = nil
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	trace.cpuLogRead[1] = nil
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	trace.cpuSleep.close()
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// traceReadCPU attempts to read from the provided profBuf[gen%2] and write</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// into the trace. Returns true if there might be more to read or false</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// if the profBuf is closed or the caller should otherwise stop reading.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// The caller is responsible for ensuring that gen does not change. Either</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// the caller must be in a traceAcquire/traceRelease block, or must be calling</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// with traceAdvanceSema held.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// No more than one goroutine may be in traceReadCPU for the same</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">// profBuf at a time.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// Must not run on the system stack because profBuf.read performs race</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// operations.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>func traceReadCPU(gen uintptr) bool {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	var pcBuf [traceStackSize]uintptr
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	data, tags, eof := trace.cpuLogRead[gen%2].read(profBufNonBlocking)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	for len(data) &gt; 0 {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		if len(data) &lt; 4 || data[0] &gt; uint64(len(data)) {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			break <span class="comment">// truncated profile</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		if data[0] &lt; 4 || tags != nil &amp;&amp; len(tags) &lt; 1 {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			break <span class="comment">// malformed profile</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		if len(tags) &lt; 1 {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			break <span class="comment">// mismatched profile records and tags</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		<span class="comment">// Deserialize the data in the profile buffer.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		recordLen := data[0]
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		timestamp := data[1]
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		ppid := data[2] &gt;&gt; 1
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		if hasP := (data[2] &amp; 0b1) != 0; !hasP {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			ppid = ^uint64(0)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		goid := data[3]
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		mpid := data[4]
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		stk := data[5:recordLen]
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">// Overflow records always have their headers contain</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		<span class="comment">// all zeroes.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		isOverflowRecord := len(stk) == 1 &amp;&amp; data[2] == 0 &amp;&amp; data[3] == 0 &amp;&amp; data[4] == 0
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		<span class="comment">// Move the data iterator forward.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		data = data[recordLen:]
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		<span class="comment">// No support here for reporting goroutine tags at the moment; if</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		<span class="comment">// that information is to be part of the execution trace, we&#39;d</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		<span class="comment">// probably want to see when the tags are applied and when they</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		<span class="comment">// change, instead of only seeing them when we get a CPU sample.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		tags = tags[1:]
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		if isOverflowRecord {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			<span class="comment">// Looks like an overflow record from the profBuf. Not much to</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			<span class="comment">// do here, we only want to report full records.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			continue
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		<span class="comment">// Construct the stack for insertion to the stack table.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		nstk := 1
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		pcBuf[0] = logicalStackSentinel
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		for ; nstk &lt; len(pcBuf) &amp;&amp; nstk-1 &lt; len(stk); nstk++ {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			pcBuf[nstk] = uintptr(stk[nstk-1])
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		<span class="comment">// Write out a trace event.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		w := unsafeTraceWriter(gen, trace.cpuBuf[gen%2])
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		<span class="comment">// Ensure we have a place to write to.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		var flushed bool
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		w, flushed = w.ensure(2 + 5*traceBytesPerNumber <span class="comment">/* traceEvCPUSamples + traceEvCPUSample + timestamp + g + m + p + stack ID */</span>)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		if flushed {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			<span class="comment">// Annotate the batch as containing strings.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			w.byte(byte(traceEvCPUSamples))
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		<span class="comment">// Add the stack to the table.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		stackID := trace.stackTab[gen%2].put(pcBuf[:nstk])
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		<span class="comment">// Write out the CPU sample.</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		w.byte(byte(traceEvCPUSample))
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		w.varint(timestamp)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		w.varint(mpid)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		w.varint(ppid)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		w.varint(goid)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		w.varint(stackID)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		trace.cpuBuf[gen%2] = w.traceBuf
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	return !eof
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// traceCPUFlush flushes trace.cpuBuf[gen%2]. The caller must be certain that gen</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// has completed and that there are no more writers to it.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// Must run on the systemstack because it flushes buffers and acquires trace.lock</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// to do so.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>func traceCPUFlush(gen uintptr) {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// Flush any remaining trace buffers containing CPU samples.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	if buf := trace.cpuBuf[gen%2]; buf != nil {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		lock(&amp;trace.lock)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		traceBufFlush(buf, gen)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		unlock(&amp;trace.lock)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		trace.cpuBuf[gen%2] = nil
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// traceCPUSample writes a CPU profile sample stack to the execution tracer&#39;s</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// profiling buffer. It is called from a signal handler, so is limited in what</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// it can do. mp must be the thread that is currently stopped in a signal.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>func traceCPUSample(gp *g, mp *m, pp *p, stk []uintptr) {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	if !traceEnabled() {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// Tracing is usually turned off; don&#39;t spend time acquiring the signal</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		<span class="comment">// lock unless it&#39;s active.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		return
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	if mp == nil {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		<span class="comment">// Drop samples that don&#39;t have an identifiable thread. We can&#39;t render</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		<span class="comment">// this in any useful way anyway.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		return
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// We&#39;re going to conditionally write to one of two buffers based on the</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// generation. To make sure we write to the correct one, we need to make</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// sure this thread&#39;s trace seqlock is held. If it already is, then we&#39;re</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">// in the tracer and we can just take advantage of that. If it isn&#39;t, then</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">// we need to acquire it and read the generation.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	locked := false
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	if mp.trace.seqlock.Load()%2 == 0 {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		mp.trace.seqlock.Add(1)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		locked = true
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	gen := trace.gen.Load()
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	if gen == 0 {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		<span class="comment">// Tracing is disabled, as it turns out. Release the seqlock if necessary</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		<span class="comment">// and exit.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		if locked {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			mp.trace.seqlock.Add(1)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		return
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	now := traceClockNow()
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	<span class="comment">// The &#34;header&#34; here is the ID of the M that was running the profiled code,</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// followed by the IDs of the P and goroutine. (For normal CPU profiling, it&#39;s</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// usually the number of samples with the given stack.) Near syscalls, pp</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// may be nil. Reporting goid of 0 is fine for either g0 or a nil gp.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	var hdr [3]uint64
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	if pp != nil {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		<span class="comment">// Overflow records in profBuf have all header values set to zero. Make</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		<span class="comment">// sure that real headers have at least one bit set.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		hdr[0] = uint64(pp.id)&lt;&lt;1 | 0b1
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	} else {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		hdr[0] = 0b10
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	if gp != nil {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		hdr[1] = gp.goid
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	if mp != nil {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		hdr[2] = uint64(mp.procid)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// Allow only one writer at a time</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	for !trace.signalLock.CompareAndSwap(0, 1) {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		<span class="comment">// TODO: Is it safe to osyield here? https://go.dev/issue/52672</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		osyield()
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	if log := trace.cpuLogWrite[gen%2].Load(); log != nil {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		<span class="comment">// Note: we don&#39;t pass a tag pointer here (how should profiling tags</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		<span class="comment">// interact with the execution tracer?), but if we did we&#39;d need to be</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		<span class="comment">// careful about write barriers. See the long comment in profBuf.write.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		log.write(nil, int64(now), hdr[:], stk)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	trace.signalLock.Store(0)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	<span class="comment">// Release the seqlock if we acquired it earlier.</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	if locked {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		mp.trace.seqlock.Add(1)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
</pre><p><a href="trace2cpu.go?m=text">View as plain text</a></p>

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
