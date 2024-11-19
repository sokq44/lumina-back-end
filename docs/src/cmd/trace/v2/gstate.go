<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/cmd/trace/v2/gstate.go - Go Documentation Server</title>

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
<a href="gstate.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/cmd">cmd</a>/<a href="http://localhost:8080/src/cmd/trace">trace</a>/<a href="http://localhost:8080/src/cmd/trace/v2">v2</a>/<span class="text-muted">gstate.go</span>
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
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/trace/traceviewer/format&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	tracev2 &#34;internal/trace/v2&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// resource is a generic constraint interface for resource IDs.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>type resource interface {
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	tracev2.GoID | tracev2.ProcID | tracev2.ThreadID
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>}
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// noResource indicates the lack of a resource.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>const noResource = -1
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// gState represents the trace viewer state of a goroutine in a trace.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// The type parameter on this type is the resource which is used to construct</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// a timeline of events. e.g. R=ProcID for a proc-oriented view, R=GoID for</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// a goroutine-oriented view, etc.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>type gState[R resource] struct {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	baseName  string
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	named     bool   <span class="comment">// Whether baseName has been set.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	label     string <span class="comment">// EventLabel extension.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	isSystemG bool
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	executing R <span class="comment">// The resource this goroutine is executing on. (Could be itself.)</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// lastStopStack is the stack trace at the point of the last</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// call to the stop method. This tends to be a more reliable way</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// of picking up stack traces, since the parser doesn&#39;t provide</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// a stack for every state transition event.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	lastStopStack tracev2.Stack
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// activeRanges is the set of all active ranges on the goroutine.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	activeRanges map[string]activeRange
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// completedRanges is a list of ranges that completed since before the</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// goroutine stopped executing. These are flushed on every stop or block.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	completedRanges []completedRange
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// startRunning is the most recent event that caused a goroutine to</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// transition to GoRunning.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	startRunningTime tracev2.Time
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// startSyscall is the most recent event that caused a goroutine to</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// transition to GoSyscall.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	syscall struct {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		time   tracev2.Time
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		stack  tracev2.Stack
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		active bool
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// startBlockReason is the StateTransition.Reason of the most recent</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// event that caused a gorotuine to transition to GoWaiting.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	startBlockReason string
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// startCause is the event that allowed this goroutine to start running.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s used to generate flow events. This is typically something like</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// an unblock event or a goroutine creation event.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// startCause.resource is the resource on which startCause happened, but is</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// listed separately because the cause may have happened on a resource that</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// isn&#39;t R (or perhaps on some abstract nebulous resource, like trace.NetpollP).</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	startCause struct {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		time     tracev2.Time
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		name     string
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		resource uint64
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		stack    tracev2.Stack
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// newGState constructs a new goroutine state for the goroutine</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// identified by the provided ID.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func newGState[R resource](goID tracev2.GoID) *gState[R] {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	return &amp;gState[R]{
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		baseName:     fmt.Sprintf(&#34;G%d&#34;, goID),
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		executing:    R(noResource),
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		activeRanges: make(map[string]activeRange),
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// augmentName attempts to use stk to augment the name of the goroutine</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// with stack information. This stack must be related to the goroutine</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// in some way, but it doesn&#39;t really matter which stack.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>func (gs *gState[R]) augmentName(stk tracev2.Stack) {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	if gs.named {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		return
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	if stk == tracev2.NoStack {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		return
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	name := lastFunc(stk)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	gs.baseName += fmt.Sprintf(&#34; %s&#34;, name)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	gs.named = true
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	gs.isSystemG = trace.IsSystemGoroutine(name)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// setLabel adds an additional label to the goroutine&#39;s name.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>func (gs *gState[R]) setLabel(label string) {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	gs.label = label
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// name returns a name for the goroutine.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func (gs *gState[R]) name() string {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	name := gs.baseName
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if gs.label != &#34;&#34; {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		name += &#34; (&#34; + gs.label + &#34;)&#34;
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	return name
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// setStartCause sets the reason a goroutine will be allowed to start soon.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// For example, via unblocking or exiting a blocked syscall.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>func (gs *gState[R]) setStartCause(ts tracev2.Time, name string, resource uint64, stack tracev2.Stack) {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	gs.startCause.time = ts
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	gs.startCause.name = name
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	gs.startCause.resource = resource
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	gs.startCause.stack = stack
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// created indicates that this goroutine was just created by the provided creator.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func (gs *gState[R]) created(ts tracev2.Time, creator R, stack tracev2.Stack) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	if creator == R(noResource) {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		return
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	gs.setStartCause(ts, &#34;go&#34;, uint64(creator), stack)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">// start indicates that a goroutine has started running on a proc.</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>func (gs *gState[R]) start(ts tracev2.Time, resource R, ctx *traceContext) {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// Set the time for all the active ranges.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	for name := range gs.activeRanges {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		gs.activeRanges[name] = activeRange{ts, tracev2.NoStack}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	if gs.startCause.name != &#34;&#34; {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		<span class="comment">// It has a start cause. Emit a flow event.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		ctx.Arrow(traceviewer.ArrowEvent{
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			Name:         gs.startCause.name,
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			Start:        ctx.elapsed(gs.startCause.time),
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			End:          ctx.elapsed(ts),
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			FromResource: uint64(gs.startCause.resource),
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			ToResource:   uint64(resource),
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			FromStack:    ctx.Stack(viewerFrames(gs.startCause.stack)),
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		})
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		gs.startCause.time = 0
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		gs.startCause.name = &#34;&#34;
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		gs.startCause.resource = 0
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		gs.startCause.stack = tracev2.NoStack
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	gs.executing = resource
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	gs.startRunningTime = ts
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// syscallBegin indicates that the goroutine entered a syscall on a proc.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>func (gs *gState[R]) syscallBegin(ts tracev2.Time, resource R, stack tracev2.Stack) {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	gs.syscall.time = ts
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	gs.syscall.stack = stack
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	gs.syscall.active = true
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	if gs.executing == R(noResource) {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		gs.executing = resource
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		gs.startRunningTime = ts
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// syscallEnd ends the syscall slice, wherever the syscall is at. This is orthogonal</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// to blockedSyscallEnd -- both must be called when a syscall ends and that syscall</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">// blocked. They&#39;re kept separate because syscallEnd indicates the point at which the</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">// goroutine is no longer executing on the resource (e.g. a proc) whereas blockedSyscallEnd</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">// is the point at which the goroutine actually exited the syscall regardless of which</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span><span class="comment">// resource that happened on.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>func (gs *gState[R]) syscallEnd(ts tracev2.Time, blocked bool, ctx *traceContext) {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	if !gs.syscall.active {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		return
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	blockString := &#34;no&#34;
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	if blocked {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		blockString = &#34;yes&#34;
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	gs.completedRanges = append(gs.completedRanges, completedRange{
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		name:       &#34;syscall&#34;,
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		startTime:  gs.syscall.time,
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		endTime:    ts,
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		startStack: gs.syscall.stack,
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		arg:        format.BlockedArg{Blocked: blockString},
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	})
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	gs.syscall.active = false
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	gs.syscall.time = 0
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	gs.syscall.stack = tracev2.NoStack
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// blockedSyscallEnd indicates the point at which the blocked syscall ended. This is distinct</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">// and orthogonal to syscallEnd; both must be called if the syscall blocked. This sets up an instant</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">// to emit a flow event from, indicating explicitly that this goroutine was unblocked by the system.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>func (gs *gState[R]) blockedSyscallEnd(ts tracev2.Time, stack tracev2.Stack, ctx *traceContext) {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	name := &#34;exit blocked syscall&#34;
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	gs.setStartCause(ts, name, trace.SyscallP, stack)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// Emit an syscall exit instant event for the &#34;Syscall&#34; lane.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	ctx.Instant(traceviewer.InstantEvent{
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		Name:     name,
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		Ts:       ctx.elapsed(ts),
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		Resource: trace.SyscallP,
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		Stack:    ctx.Stack(viewerFrames(stack)),
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	})
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// unblock indicates that the goroutine gs represents has been unblocked.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>func (gs *gState[R]) unblock(ts tracev2.Time, stack tracev2.Stack, resource R, ctx *traceContext) {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	name := &#34;unblock&#34;
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	viewerResource := uint64(resource)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	if gs.startBlockReason != &#34;&#34; {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		name = fmt.Sprintf(&#34;%s (%s)&#34;, name, gs.startBlockReason)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	if strings.Contains(gs.startBlockReason, &#34;network&#34;) {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		<span class="comment">// Attribute the network instant to the nebulous &#34;NetpollP&#34; if</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		<span class="comment">// resource isn&#39;t a thread, because there&#39;s a good chance that</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		<span class="comment">// resource isn&#39;t going to be valid in this case.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		<span class="comment">// TODO(mknyszek): Handle this invalidness in a more general way.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		if _, ok := any(resource).(tracev2.ThreadID); !ok {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			<span class="comment">// Emit an unblock instant event for the &#34;Network&#34; lane.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			viewerResource = trace.NetpollP
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		ctx.Instant(traceviewer.InstantEvent{
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			Name:     name,
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			Ts:       ctx.elapsed(ts),
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			Resource: viewerResource,
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			Stack:    ctx.Stack(viewerFrames(stack)),
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		})
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	gs.startBlockReason = &#34;&#34;
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if viewerResource != 0 {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		gs.setStartCause(ts, name, viewerResource, stack)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// block indicates that the goroutine has stopped executing on a proc -- specifically,</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// it blocked for some reason.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>func (gs *gState[R]) block(ts tracev2.Time, stack tracev2.Stack, reason string, ctx *traceContext) {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	gs.startBlockReason = reason
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	gs.stop(ts, stack, ctx)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// stop indicates that the goroutine has stopped executing on a proc.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>func (gs *gState[R]) stop(ts tracev2.Time, stack tracev2.Stack, ctx *traceContext) {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// Emit the execution time slice.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	var stk int
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	if gs.lastStopStack != tracev2.NoStack {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		stk = ctx.Stack(viewerFrames(gs.lastStopStack))
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// Check invariants.</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	if gs.startRunningTime == 0 {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		panic(&#34;silently broken trace or generator invariant (startRunningTime != 0) not held&#34;)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	if gs.executing == R(noResource) {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		panic(&#34;non-executing goroutine stopped&#34;)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	ctx.Slice(traceviewer.SliceEvent{
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		Name:     gs.name(),
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		Ts:       ctx.elapsed(gs.startRunningTime),
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		Dur:      ts.Sub(gs.startRunningTime),
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		Resource: uint64(gs.executing),
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		Stack:    stk,
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	})
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">// Flush completed ranges.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	for _, cr := range gs.completedRanges {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		ctx.Slice(traceviewer.SliceEvent{
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			Name:     cr.name,
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			Ts:       ctx.elapsed(cr.startTime),
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			Dur:      cr.endTime.Sub(cr.startTime),
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			Resource: uint64(gs.executing),
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			Stack:    ctx.Stack(viewerFrames(cr.startStack)),
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			EndStack: ctx.Stack(viewerFrames(cr.endStack)),
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			Arg:      cr.arg,
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		})
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	gs.completedRanges = gs.completedRanges[:0]
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">// Continue in-progress ranges.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	for name, r := range gs.activeRanges {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		<span class="comment">// Check invariant.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		if r.time == 0 {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			panic(&#34;silently broken trace or generator invariant (activeRanges time != 0) not held&#34;)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		ctx.Slice(traceviewer.SliceEvent{
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			Name:     name,
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			Ts:       ctx.elapsed(r.time),
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			Dur:      ts.Sub(r.time),
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			Resource: uint64(gs.executing),
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			Stack:    ctx.Stack(viewerFrames(r.stack)),
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		})
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">// Clear the range info.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	for name := range gs.activeRanges {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		gs.activeRanges[name] = activeRange{0, tracev2.NoStack}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	gs.startRunningTime = 0
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	gs.lastStopStack = stack
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	gs.executing = R(noResource)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// finalize writes out any in-progress slices as if the goroutine stopped.</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">// This must only be used once the trace has been fully processed and no</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">// further events will be processed. This method may leave the gState in</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">// an inconsistent state.</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>func (gs *gState[R]) finish(ctx *traceContext) {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	if gs.executing != R(noResource) {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		gs.syscallEnd(ctx.endTime, false, ctx)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		gs.stop(ctx.endTime, tracev2.NoStack, ctx)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">// rangeBegin indicates the start of a special range of time.</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>func (gs *gState[R]) rangeBegin(ts tracev2.Time, name string, stack tracev2.Stack) {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if gs.executing != R(noResource) {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		<span class="comment">// If we&#39;re executing, start the slice from here.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		gs.activeRanges[name] = activeRange{ts, stack}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	} else {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		<span class="comment">// If the goroutine isn&#39;t executing, there&#39;s no place for</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		<span class="comment">// us to create a slice from. Wait until it starts executing.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		gs.activeRanges[name] = activeRange{0, stack}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">// rangeActive indicates that a special range of time has been in progress.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>func (gs *gState[R]) rangeActive(name string) {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	if gs.executing != R(noResource) {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		<span class="comment">// If we&#39;re executing, and the range is active, then start</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		<span class="comment">// from wherever the goroutine started running from.</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		gs.activeRanges[name] = activeRange{gs.startRunningTime, tracev2.NoStack}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	} else {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		<span class="comment">// If the goroutine isn&#39;t executing, there&#39;s no place for</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		<span class="comment">// us to create a slice from. Wait until it starts executing.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		gs.activeRanges[name] = activeRange{0, tracev2.NoStack}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// rangeEnd indicates the end of a special range of time.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>func (gs *gState[R]) rangeEnd(ts tracev2.Time, name string, stack tracev2.Stack, ctx *traceContext) {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	if gs.executing != R(noResource) {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		r := gs.activeRanges[name]
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		gs.completedRanges = append(gs.completedRanges, completedRange{
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			name:       name,
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			startTime:  r.time,
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			endTime:    ts,
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			startStack: r.stack,
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			endStack:   stack,
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		})
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	delete(gs.activeRanges, name)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>}
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>func lastFunc(s tracev2.Stack) string {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	var last tracev2.StackFrame
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	s.Frames(func(f tracev2.StackFrame) bool {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		last = f
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		return true
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	})
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	return last.Func
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
</pre><p><a href="gstate.go?m=text">View as plain text</a></p>

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
