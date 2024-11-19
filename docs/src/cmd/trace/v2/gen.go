<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/cmd/trace/v2/gen.go - Go Documentation Server</title>

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
<a href="gen.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/cmd">cmd</a>/<a href="http://localhost:8080/src/cmd/trace">trace</a>/<a href="http://localhost:8080/src/cmd/trace/v2">v2</a>/<span class="text-muted">gen.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/trace&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/trace/traceviewer&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	tracev2 &#34;internal/trace/v2&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// generator is an interface for generating a JSON trace for the trace viewer</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// from a trace. Each method in this interface is a handler for a kind of event</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// that is interesting to render in the UI via the JSON trace.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type generator interface {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// Global parts.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	Sync() <span class="comment">// Notifies the generator of an EventSync event.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	StackSample(ctx *traceContext, ev *tracev2.Event)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	GlobalRange(ctx *traceContext, ev *tracev2.Event)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	GlobalMetric(ctx *traceContext, ev *tracev2.Event)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// Goroutine parts.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	GoroutineLabel(ctx *traceContext, ev *tracev2.Event)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	GoroutineRange(ctx *traceContext, ev *tracev2.Event)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	GoroutineTransition(ctx *traceContext, ev *tracev2.Event)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// Proc parts.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	ProcRange(ctx *traceContext, ev *tracev2.Event)
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	ProcTransition(ctx *traceContext, ev *tracev2.Event)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// User annotations.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	Log(ctx *traceContext, ev *tracev2.Event)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// Finish indicates the end of the trace and finalizes generation.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	Finish(ctx *traceContext)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// runGenerator produces a trace into ctx by running the generator over the parsed trace.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func runGenerator(ctx *traceContext, g generator, parsed *parsedTrace, opts *genOpts) {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	for i := range parsed.events {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		ev := &amp;parsed.events[i]
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		switch ev.Kind() {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		case tracev2.EventSync:
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			g.Sync()
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		case tracev2.EventStackSample:
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			g.StackSample(ctx, ev)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		case tracev2.EventRangeBegin, tracev2.EventRangeActive, tracev2.EventRangeEnd:
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			r := ev.Range()
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			switch r.Scope.Kind {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			case tracev2.ResourceGoroutine:
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>				g.GoroutineRange(ctx, ev)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			case tracev2.ResourceProc:
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>				g.ProcRange(ctx, ev)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			case tracev2.ResourceNone:
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>				g.GlobalRange(ctx, ev)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		case tracev2.EventMetric:
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			g.GlobalMetric(ctx, ev)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		case tracev2.EventLabel:
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			l := ev.Label()
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			if l.Resource.Kind == tracev2.ResourceGoroutine {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>				g.GoroutineLabel(ctx, ev)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		case tracev2.EventStateTransition:
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			switch ev.StateTransition().Resource.Kind {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			case tracev2.ResourceProc:
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>				g.ProcTransition(ctx, ev)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			case tracev2.ResourceGoroutine:
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>				g.GoroutineTransition(ctx, ev)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		case tracev2.EventLog:
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			g.Log(ctx, ev)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	for i, task := range opts.tasks {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		emitTask(ctx, task, i)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		if opts.mode&amp;traceviewer.ModeGoroutineOriented != 0 {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			for _, region := range task.Regions {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>				emitRegion(ctx, region)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	g.Finish(ctx)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// emitTask emits information about a task into the trace viewer&#39;s event stream.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// sortIndex sets the order in which this task will appear related to other tasks,</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// lowest first.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>func emitTask(ctx *traceContext, task *trace.UserTaskSummary, sortIndex int) {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// Collect information about the task.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	var startStack, endStack tracev2.Stack
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	var startG, endG tracev2.GoID
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	startTime, endTime := ctx.startTime, ctx.endTime
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	if task.Start != nil {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		startStack = task.Start.Stack()
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		startG = task.Start.Goroutine()
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		startTime = task.Start.Time()
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	if task.End != nil {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		endStack = task.End.Stack()
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		endG = task.End.Goroutine()
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		endTime = task.End.Time()
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	arg := struct {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		ID     uint64 `json:&#34;id&#34;`
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		StartG uint64 `json:&#34;start_g,omitempty&#34;`
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		EndG   uint64 `json:&#34;end_g,omitempty&#34;`
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	}{
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		ID:     uint64(task.ID),
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		StartG: uint64(startG),
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		EndG:   uint64(endG),
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// Emit the task slice and notify the emitter of the task.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	ctx.Task(uint64(task.ID), fmt.Sprintf(&#34;T%d %s&#34;, task.ID, task.Name), sortIndex)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	ctx.TaskSlice(traceviewer.SliceEvent{
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		Name:     task.Name,
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		Ts:       ctx.elapsed(startTime),
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		Dur:      endTime.Sub(startTime),
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		Resource: uint64(task.ID),
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		Stack:    ctx.Stack(viewerFrames(startStack)),
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		EndStack: ctx.Stack(viewerFrames(endStack)),
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		Arg:      arg,
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	})
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// Emit an arrow from the parent to the child.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	if task.Parent != nil &amp;&amp; task.Start != nil &amp;&amp; task.Start.Kind() == tracev2.EventTaskBegin {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		ctx.TaskArrow(traceviewer.ArrowEvent{
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			Name:         &#34;newTask&#34;,
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			Start:        ctx.elapsed(task.Start.Time()),
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			End:          ctx.elapsed(task.Start.Time()),
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			FromResource: uint64(task.Parent.ID),
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			ToResource:   uint64(task.ID),
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			FromStack:    ctx.Stack(viewerFrames(task.Start.Stack())),
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		})
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// emitRegion emits goroutine-based slice events to the UI. The caller</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// must be emitting for a goroutine-oriented trace.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// TODO(mknyszek): Make regions part of the regular generator loop and</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// treat them like ranges so that we can emit regions in traces oriented</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// by proc or thread.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>func emitRegion(ctx *traceContext, region *trace.UserRegionSummary) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	if region.Name == &#34;&#34; {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		return
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// Collect information about the region.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	var startStack, endStack tracev2.Stack
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	goroutine := tracev2.NoGoroutine
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	startTime, endTime := ctx.startTime, ctx.endTime
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	if region.Start != nil {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		startStack = region.Start.Stack()
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		startTime = region.Start.Time()
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		goroutine = region.Start.Goroutine()
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	if region.End != nil {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		endStack = region.End.Stack()
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		endTime = region.End.Time()
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		goroutine = region.End.Goroutine()
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	if goroutine == tracev2.NoGoroutine {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		return
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	arg := struct {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		TaskID uint64 `json:&#34;taskid&#34;`
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}{
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		TaskID: uint64(region.TaskID),
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	ctx.AsyncSlice(traceviewer.AsyncSliceEvent{
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		SliceEvent: traceviewer.SliceEvent{
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			Name:     region.Name,
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			Ts:       ctx.elapsed(startTime),
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			Dur:      endTime.Sub(startTime),
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			Resource: uint64(goroutine),
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			Stack:    ctx.Stack(viewerFrames(startStack)),
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			EndStack: ctx.Stack(viewerFrames(endStack)),
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			Arg:      arg,
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		},
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		Category:       &#34;Region&#34;,
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		Scope:          fmt.Sprintf(&#34;%x&#34;, region.TaskID),
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		TaskColorIndex: uint64(region.TaskID),
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	})
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// Building blocks for generators.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">// stackSampleGenerator implements a generic handler for stack sample events.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span><span class="comment">// The provided resource is the resource the stack sample should count against.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>type stackSampleGenerator[R resource] struct {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// getResource is a function to extract a resource ID from a stack sample event.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	getResource func(*tracev2.Event) R
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// StackSample implements a stack sample event handler. It expects ev to be one such event.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>func (g *stackSampleGenerator[R]) StackSample(ctx *traceContext, ev *tracev2.Event) {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	id := g.getResource(ev)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	if id == R(noResource) {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		<span class="comment">// We have nowhere to put this in the UI.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		return
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	ctx.Instant(traceviewer.InstantEvent{
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		Name:     &#34;CPU profile sample&#34;,
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		Ts:       ctx.elapsed(ev.Time()),
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		Resource: uint64(id),
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		Stack:    ctx.Stack(viewerFrames(ev.Stack())),
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	})
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// globalRangeGenerator implements a generic handler for EventRange* events that pertain</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// to tracev2.ResourceNone (the global scope).</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>type globalRangeGenerator struct {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	ranges   map[string]activeRange
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	seenSync bool
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// Sync notifies the generator of an EventSync event.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>func (g *globalRangeGenerator) Sync() {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	g.seenSync = true
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// GlobalRange implements a handler for EventRange* events whose Scope.Kind is ResourceNone.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// It expects ev to be one such event.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>func (g *globalRangeGenerator) GlobalRange(ctx *traceContext, ev *tracev2.Event) {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	if g.ranges == nil {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		g.ranges = make(map[string]activeRange)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	r := ev.Range()
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	switch ev.Kind() {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	case tracev2.EventRangeBegin:
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		g.ranges[r.Name] = activeRange{ev.Time(), ev.Stack()}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	case tracev2.EventRangeActive:
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		<span class="comment">// If we&#39;ve seen a Sync event, then Active events are always redundant.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		if !g.seenSync {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			<span class="comment">// Otherwise, they extend back to the start of the trace.</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			g.ranges[r.Name] = activeRange{ctx.startTime, ev.Stack()}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	case tracev2.EventRangeEnd:
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		<span class="comment">// Only emit GC events, because we have nowhere to</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		<span class="comment">// put other events.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		ar := g.ranges[r.Name]
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		if strings.Contains(r.Name, &#34;GC&#34;) {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			ctx.Slice(traceviewer.SliceEvent{
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>				Name:     r.Name,
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>				Ts:       ctx.elapsed(ar.time),
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>				Dur:      ev.Time().Sub(ar.time),
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>				Resource: trace.GCP,
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>				Stack:    ctx.Stack(viewerFrames(ar.stack)),
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>				EndStack: ctx.Stack(viewerFrames(ev.Stack())),
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			})
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		delete(g.ranges, r.Name)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span><span class="comment">// Finish flushes any outstanding ranges at the end of the trace.</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>func (g *globalRangeGenerator) Finish(ctx *traceContext) {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	for name, ar := range g.ranges {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		if !strings.Contains(name, &#34;GC&#34;) {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			continue
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		ctx.Slice(traceviewer.SliceEvent{
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			Name:     name,
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			Ts:       ctx.elapsed(ar.time),
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			Dur:      ctx.endTime.Sub(ar.time),
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			Resource: trace.GCP,
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			Stack:    ctx.Stack(viewerFrames(ar.stack)),
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		})
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// globalMetricGenerator implements a generic handler for Metric events.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>type globalMetricGenerator struct {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// GlobalMetric implements an event handler for EventMetric events. ev must be one such event.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>func (g *globalMetricGenerator) GlobalMetric(ctx *traceContext, ev *tracev2.Event) {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	m := ev.Metric()
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	switch m.Name {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	case &#34;/memory/classes/heap/objects:bytes&#34;:
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		ctx.HeapAlloc(ctx.elapsed(ev.Time()), m.Value.Uint64())
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	case &#34;/gc/heap/goal:bytes&#34;:
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		ctx.HeapGoal(ctx.elapsed(ev.Time()), m.Value.Uint64())
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	case &#34;/sched/gomaxprocs:threads&#34;:
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		ctx.Gomaxprocs(m.Value.Uint64())
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span><span class="comment">// procRangeGenerator implements a generic handler for EventRange* events whose Scope.Kind is</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span><span class="comment">// ResourceProc.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>type procRangeGenerator struct {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	ranges   map[tracev2.Range]activeRange
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	seenSync bool
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">// Sync notifies the generator of an EventSync event.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>func (g *procRangeGenerator) Sync() {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	g.seenSync = true
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">// ProcRange implements a handler for EventRange* events whose Scope.Kind is ResourceProc.</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">// It expects ev to be one such event.</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>func (g *procRangeGenerator) ProcRange(ctx *traceContext, ev *tracev2.Event) {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	if g.ranges == nil {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		g.ranges = make(map[tracev2.Range]activeRange)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	}
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	r := ev.Range()
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	switch ev.Kind() {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	case tracev2.EventRangeBegin:
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		g.ranges[r] = activeRange{ev.Time(), ev.Stack()}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	case tracev2.EventRangeActive:
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		<span class="comment">// If we&#39;ve seen a Sync event, then Active events are always redundant.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		if !g.seenSync {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			<span class="comment">// Otherwise, they extend back to the start of the trace.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			g.ranges[r] = activeRange{ctx.startTime, ev.Stack()}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	case tracev2.EventRangeEnd:
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		<span class="comment">// Emit proc-based ranges.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		ar := g.ranges[r]
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		ctx.Slice(traceviewer.SliceEvent{
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			Name:     r.Name,
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			Ts:       ctx.elapsed(ar.time),
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			Dur:      ev.Time().Sub(ar.time),
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>			Resource: uint64(r.Scope.Proc()),
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>			Stack:    ctx.Stack(viewerFrames(ar.stack)),
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			EndStack: ctx.Stack(viewerFrames(ev.Stack())),
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		})
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		delete(g.ranges, r)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// Finish flushes any outstanding ranges at the end of the trace.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>func (g *procRangeGenerator) Finish(ctx *traceContext) {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	for r, ar := range g.ranges {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		ctx.Slice(traceviewer.SliceEvent{
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			Name:     r.Name,
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			Ts:       ctx.elapsed(ar.time),
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			Dur:      ctx.endTime.Sub(ar.time),
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			Resource: uint64(r.Scope.Proc()),
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			Stack:    ctx.Stack(viewerFrames(ar.stack)),
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		})
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// activeRange represents an active EventRange* range.</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>type activeRange struct {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	time  tracev2.Time
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	stack tracev2.Stack
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// completedRange represents a completed EventRange* range.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>type completedRange struct {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	name       string
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	startTime  tracev2.Time
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	endTime    tracev2.Time
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	startStack tracev2.Stack
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	endStack   tracev2.Stack
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	arg        any
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>}
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>type logEventGenerator[R resource] struct {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	<span class="comment">// getResource is a function to extract a resource ID from a Log event.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	getResource func(*tracev2.Event) R
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>}
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span><span class="comment">// Log implements a log event handler. It expects ev to be one such event.</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>func (g *logEventGenerator[R]) Log(ctx *traceContext, ev *tracev2.Event) {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	id := g.getResource(ev)
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	if id == R(noResource) {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		<span class="comment">// We have nowhere to put this in the UI.</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		return
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	<span class="comment">// Construct the name to present.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	log := ev.Log()
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	name := log.Message
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	if log.Category != &#34;&#34; {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		name = &#34;[&#34; + log.Category + &#34;] &#34; + name
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	<span class="comment">// Emit an instant event.</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	ctx.Instant(traceviewer.InstantEvent{
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		Name:     name,
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		Ts:       ctx.elapsed(ev.Time()),
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		Category: &#34;user event&#34;,
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		Resource: uint64(id),
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		Stack:    ctx.Stack(viewerFrames(ev.Stack())),
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	})
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
</pre><p><a href="gen.go?m=text">View as plain text</a></p>

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
