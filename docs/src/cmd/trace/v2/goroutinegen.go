<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/cmd/trace/v2/goroutinegen.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../../index.html">GoDoc</a></div>
<a href="goroutinegen.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/cmd">cmd</a>/<a href="http://localhost:8080/src/cmd/trace">trace</a>/<a href="http://localhost:8080/src/cmd/trace/v2">v2</a>/<span class="text-muted">goroutinegen.go</span>
  </h1>





  <h2>
    Documentation: <a href="../../../../pkg/cmd/trace/v2">cmd/trace/v2</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2023 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package trace
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	tracev2 &#34;internal/trace/v2&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>)
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>var _ generator = &amp;goroutineGenerator{}
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>type goroutineGenerator struct {
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	globalRangeGenerator
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	globalMetricGenerator
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	stackSampleGenerator[tracev2.GoID]
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	logEventGenerator[tracev2.GoID]
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	gStates map[tracev2.GoID]*gState[tracev2.GoID]
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	focus   tracev2.GoID
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	filter  map[tracev2.GoID]struct{}
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>func newGoroutineGenerator(ctx *traceContext, focus tracev2.GoID, filter map[tracev2.GoID]struct{}) *goroutineGenerator {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	gg := new(goroutineGenerator)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	rg := func(ev *tracev2.Event) tracev2.GoID {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		return ev.Goroutine()
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	gg.stackSampleGenerator.getResource = rg
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	gg.logEventGenerator.getResource = rg
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	gg.gStates = make(map[tracev2.GoID]*gState[tracev2.GoID])
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	gg.focus = focus
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	gg.filter = filter
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// Enable a filter on the emitter.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	if filter != nil {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		ctx.SetResourceFilter(func(resource uint64) bool {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>			_, ok := filter[tracev2.GoID(resource)]
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>			return ok
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		})
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	return gg
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>func (g *goroutineGenerator) Sync() {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	g.globalRangeGenerator.Sync()
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>func (g *goroutineGenerator) GoroutineLabel(ctx *traceContext, ev *tracev2.Event) {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	l := ev.Label()
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	g.gStates[l.Resource.Goroutine()].setLabel(l.Label)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func (g *goroutineGenerator) GoroutineRange(ctx *traceContext, ev *tracev2.Event) {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	r := ev.Range()
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	switch ev.Kind() {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	case tracev2.EventRangeBegin:
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		g.gStates[r.Scope.Goroutine()].rangeBegin(ev.Time(), r.Name, ev.Stack())
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	case tracev2.EventRangeActive:
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		g.gStates[r.Scope.Goroutine()].rangeActive(r.Name)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	case tracev2.EventRangeEnd:
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		gs := g.gStates[r.Scope.Goroutine()]
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		gs.rangeEnd(ev.Time(), r.Name, ev.Stack(), ctx)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>func (g *goroutineGenerator) GoroutineTransition(ctx *traceContext, ev *tracev2.Event) {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	st := ev.StateTransition()
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	goID := st.Resource.Goroutine()
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// If we haven&#39;t seen this goroutine before, create a new</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// gState for it.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	gs, ok := g.gStates[goID]
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	if !ok {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		gs = newGState[tracev2.GoID](goID)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		g.gStates[goID] = gs
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// Try to augment the name of the goroutine.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	gs.augmentName(st.Stack)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// Handle the goroutine state transition.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	from, to := st.Goroutine()
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	if from == to {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		<span class="comment">// Filter out no-op events.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		return
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	if from.Executing() &amp;&amp; !to.Executing() {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		if to == tracev2.GoWaiting {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			<span class="comment">// Goroutine started blocking.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			gs.block(ev.Time(), ev.Stack(), st.Reason, ctx)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		} else {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			gs.stop(ev.Time(), ev.Stack(), ctx)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if !from.Executing() &amp;&amp; to.Executing() {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		start := ev.Time()
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		if from == tracev2.GoUndetermined {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			<span class="comment">// Back-date the event to the start of the trace.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			start = ctx.startTime
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		gs.start(start, goID, ctx)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	if from == tracev2.GoWaiting {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">// Goroutine unblocked.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		gs.unblock(ev.Time(), ev.Stack(), ev.Goroutine(), ctx)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	if from == tracev2.GoNotExist &amp;&amp; to == tracev2.GoRunnable {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		<span class="comment">// Goroutine was created.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		gs.created(ev.Time(), ev.Goroutine(), ev.Stack())
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	if from == tracev2.GoSyscall &amp;&amp; to != tracev2.GoRunning {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		<span class="comment">// Exiting blocked syscall.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		gs.syscallEnd(ev.Time(), true, ctx)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		gs.blockedSyscallEnd(ev.Time(), ev.Stack(), ctx)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	} else if from == tracev2.GoSyscall {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		<span class="comment">// Check if we&#39;re exiting a syscall in a non-blocking way.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		gs.syscallEnd(ev.Time(), false, ctx)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// Handle syscalls.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if to == tracev2.GoSyscall {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		start := ev.Time()
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		if from == tracev2.GoUndetermined {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			<span class="comment">// Back-date the event to the start of the trace.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			start = ctx.startTime
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		<span class="comment">// Write down that we&#39;ve entered a syscall. Note: we might have no G or P here</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		<span class="comment">// if we&#39;re in a cgo callback or this is a transition from GoUndetermined</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		<span class="comment">// (i.e. the G has been blocked in a syscall).</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		gs.syscallBegin(start, goID, ev.Stack())
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// Note down the goroutine transition.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	_, inMarkAssist := gs.activeRanges[&#34;GC mark assist&#34;]
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	ctx.GoroutineTransition(ctx.elapsed(ev.Time()), viewerGState(from, inMarkAssist), viewerGState(to, inMarkAssist))
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>func (g *goroutineGenerator) ProcRange(ctx *traceContext, ev *tracev2.Event) {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Extend procRangeGenerator to support rendering proc ranges</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// that overlap with a goroutine&#39;s execution.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>func (g *goroutineGenerator) ProcTransition(ctx *traceContext, ev *tracev2.Event) {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// Not needed. All relevant information for goroutines can be derived from goroutine transitions.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>func (g *goroutineGenerator) Finish(ctx *traceContext) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	ctx.SetResourceType(&#34;G&#34;)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// Finish off global ranges.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	g.globalRangeGenerator.Finish(ctx)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// Finish off all the goroutine slices.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	for id, gs := range g.gStates {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		gs.finish(ctx)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		<span class="comment">// Tell the emitter about the goroutines we want to render.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		ctx.Resource(uint64(id), gs.name())
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// Set the goroutine to focus on.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	if g.focus != tracev2.NoGoroutine {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		ctx.Focus(uint64(g.focus))
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
</pre><p><a href="goroutinegen.go?m=text">View as plain text</a></p>

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
