<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/coro.go - Go Documentation Server</title>

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
<a href="coro.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">coro.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;unsafe&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// A coro represents extra concurrency without extra parallelism,</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// as would be needed for a coroutine implementation.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// The coro does not represent a specific coroutine, only the ability</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// to do coroutine-style control transfers.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// It can be thought of as like a special channel that always has</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// a goroutine blocked on it. If another goroutine calls coroswitch(c),</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// the caller becomes the goroutine blocked in c, and the goroutine</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// formerly blocked in c starts running.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// These switches continue until a call to coroexit(c),</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// which ends the use of the coro by releasing the blocked</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// goroutine in c and exiting the current goroutine.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// Coros are heap allocated and garbage collected, so that user code</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// can hold a pointer to a coro without causing potential dangling</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// pointer errors.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>type coro struct {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	gp guintptr
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	f  func(*coro)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//go:linkname newcoro</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// newcoro creates a new coro containing a</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// goroutine blocked waiting to run f</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// and returns that coro.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>func newcoro(f func(*coro)) *coro {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	c := new(coro)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	c.f = f
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	pc := getcallerpc()
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	gp := getg()
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		start := corostart
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		startfv := *(**funcval)(unsafe.Pointer(&amp;start))
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		gp = newproc1(startfv, gp, pc)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	})
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	gp.coroarg = c
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	gp.waitreason = waitReasonCoroutine
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	casgstatus(gp, _Grunnable, _Gwaiting)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	c.gp.set(gp)
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	return c
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//go:linkname corostart</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// corostart is the entry func for a new coroutine.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// It runs the coroutine user function f passed to corostart</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// and then calls coroexit to remove the extra concurrency.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>func corostart() {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	gp := getg()
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	c := gp.coroarg
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	gp.coroarg = nil
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	c.f(c)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	coroexit(c)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// coroexit is like coroswitch but closes the coro</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// and exits the current goroutine</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>func coroexit(c *coro) {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	gp := getg()
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	gp.coroarg = c
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	gp.coroexit = true
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	mcall(coroswitch_m)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//go:linkname coroswitch</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// coroswitch switches to the goroutine blocked on c</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// and then blocks the current goroutine on c.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>func coroswitch(c *coro) {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	gp := getg()
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	gp.coroarg = c
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	mcall(coroswitch_m)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// coroswitch_m is the implementation of coroswitch</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// that runs on the m stack.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// Note: Coroutine switches are expected to happen at</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// an order of magnitude (or more) higher frequency</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// than regular goroutine switches, so this path is heavily</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// optimized to remove unnecessary work.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// The fast path here is three CAS: the one at the top on gp.atomicstatus,</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// the one in the middle to choose the next g,</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// and the one at the bottom on gnext.atomicstatus.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// It is important not to add more atomic operations or other</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// expensive operations to the fast path.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func coroswitch_m(gp *g) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// TODO(rsc,mknyszek): add tracing support in a lightweight manner.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// Probably the tracer will need a global bool (set and cleared during STW)</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// that this code can check to decide whether to use trace.gen.Load();</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// we do not want to do the atomic load all the time, especially when</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// tracer use is relatively rare.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	c := gp.coroarg
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	gp.coroarg = nil
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	exit := gp.coroexit
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	gp.coroexit = false
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	mp := gp.m
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	if exit {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		gdestroy(gp)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		gp = nil
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	} else {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		<span class="comment">// If we can CAS ourselves directly from running to waiting, so do,</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		<span class="comment">// keeping the control transfer as lightweight as possible.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		gp.waitreason = waitReasonCoroutine
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		if !gp.atomicstatus.CompareAndSwap(_Grunning, _Gwaiting) {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			<span class="comment">// The CAS failed: use casgstatus, which will take care of</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			<span class="comment">// coordinating with the garbage collector about the state change.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			casgstatus(gp, _Grunning, _Gwaiting)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		<span class="comment">// Clear gp.m.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		setMNoWB(&amp;gp.m, nil)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// The goroutine stored in c is the one to run next.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// Swap it with ourselves.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	var gnext *g
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	for {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		<span class="comment">// Note: this is a racy load, but it will eventually</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		<span class="comment">// get the right value, and if it gets the wrong value,</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		<span class="comment">// the c.gp.cas will fail, so no harm done other than</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		<span class="comment">// a wasted loop iteration.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		<span class="comment">// The cas will also sync c.gp&#39;s</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// memory enough that the next iteration of the racy load</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// should see the correct value.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// We are avoiding the atomic load to keep this path</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		<span class="comment">// as lightweight as absolutely possible.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		<span class="comment">// (The atomic load is free on x86 but not free elsewhere.)</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		next := c.gp
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		if next.ptr() == nil {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			throw(&#34;coroswitch on exited coro&#34;)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		var self guintptr
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		self.set(gp)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		if c.gp.cas(next, self) {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			gnext = next.ptr()
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			break
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// Start running next, without heavy scheduling machinery.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// Set mp.curg and gnext.m and then update scheduling state</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// directly if possible.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	setGNoWB(&amp;mp.curg, gnext)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	setMNoWB(&amp;gnext.m, mp)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if !gnext.atomicstatus.CompareAndSwap(_Gwaiting, _Grunning) {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// The CAS failed: use casgstatus, which will take care of</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// coordinating with the garbage collector about the state change.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		casgstatus(gnext, _Gwaiting, _Grunnable)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		casgstatus(gnext, _Grunnable, _Grunning)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// Switch to gnext. Does not return.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	gogo(&amp;gnext.sched)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
</pre><p><a href="coro.go?m=text">View as plain text</a></p>

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
