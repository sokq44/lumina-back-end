<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/cmd/trace/v2/pprof.go - Go Documentation Server</title>

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
<a href="pprof.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/cmd">cmd</a>/<a href="http://localhost:8080/src/cmd/trace">trace</a>/<a href="http://localhost:8080/src/cmd/trace/v2">v2</a>/<span class="text-muted">pprof.go</span>
  </h1>





  <h2>
    Documentation: <a href="../../../../pkg/cmd/trace/v2">cmd/trace/v2</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2014 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Serving of pprof-like profiles.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package trace
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;cmp&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/trace&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;internal/trace/traceviewer&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	tracev2 &#34;internal/trace/v2&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;net/http&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;slices&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>func pprofByGoroutine(compute computePprofFunc, t *parsedTrace) traceviewer.ProfileFunc {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	return func(r *http.Request) ([]traceviewer.ProfileRecord, error) {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>		name := r.FormValue(&#34;name&#34;)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>		gToIntervals, err := pprofMatchingGoroutines(name, t)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		if err != nil {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>			return nil, err
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		return compute(gToIntervals, t.events)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>func pprofByRegion(compute computePprofFunc, t *parsedTrace) traceviewer.ProfileFunc {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	return func(r *http.Request) ([]traceviewer.ProfileRecord, error) {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		filter, err := newRegionFilter(r)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		if err != nil {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>			return nil, err
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		gToIntervals, err := pprofMatchingRegions(filter, t)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		if err != nil {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>			return nil, err
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		return compute(gToIntervals, t.events)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// pprofMatchingGoroutines returns the ids of goroutines of the matching name and its interval.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// If the id string is empty, returns nil without an error.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>func pprofMatchingGoroutines(name string, t *parsedTrace) (map[tracev2.GoID][]interval, error) {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	res := make(map[tracev2.GoID][]interval)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	for _, g := range t.summary.Goroutines {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		if g.Name != name {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			continue
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		endTime := g.EndTime
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		if g.EndTime == 0 {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			endTime = t.endTime() <span class="comment">// Use the trace end time, since the goroutine is still live then.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		res[g.ID] = []interval{{start: g.StartTime, end: endTime}}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if len(res) == 0 {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return nil, fmt.Errorf(&#34;failed to find matching goroutines for name: %s&#34;, name)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	return res, nil
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// pprofMatchingRegions returns the time intervals of matching regions</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// grouped by the goroutine id. If the filter is nil, returns nil without an error.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>func pprofMatchingRegions(filter *regionFilter, t *parsedTrace) (map[tracev2.GoID][]interval, error) {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	if filter == nil {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		return nil, nil
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	gToIntervals := make(map[tracev2.GoID][]interval)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	for _, g := range t.summary.Goroutines {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		for _, r := range g.Regions {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			if !filter.match(t, r) {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>				continue
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			gToIntervals[g.ID] = append(gToIntervals[g.ID], regionInterval(t, r))
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	for g, intervals := range gToIntervals {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		<span class="comment">// In order to remove nested regions and</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		<span class="comment">// consider only the outermost regions,</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		<span class="comment">// first, we sort based on the start time</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		<span class="comment">// and then scan through to select only the outermost regions.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		slices.SortFunc(intervals, func(a, b interval) int {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			if c := cmp.Compare(a.start, b.start); c != 0 {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>				return c
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			return cmp.Compare(a.end, b.end)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		})
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		var lastTimestamp tracev2.Time
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		var n int
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		<span class="comment">// Select only the outermost regions.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		for _, i := range intervals {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			if lastTimestamp &lt;= i.start {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>				intervals[n] = i <span class="comment">// new non-overlapping region starts.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>				lastTimestamp = i.end
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>				n++
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			<span class="comment">// Otherwise, skip because this region overlaps with a previous region.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		gToIntervals[g] = intervals[:n]
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	return gToIntervals, nil
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>type computePprofFunc func(gToIntervals map[tracev2.GoID][]interval, events []tracev2.Event) ([]traceviewer.ProfileRecord, error)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// computePprofIO returns a computePprofFunc that generates IO pprof-like profile (time spent in</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// IO wait, currently only network blocking event).</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>func computePprofIO() computePprofFunc {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	return makeComputePprofFunc(tracev2.GoWaiting, func(reason string) bool {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		return reason == &#34;network&#34;
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	})
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// computePprofBlock returns a computePprofFunc that generates blocking pprof-like profile</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// (time spent blocked on synchronization primitives).</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>func computePprofBlock() computePprofFunc {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	return makeComputePprofFunc(tracev2.GoWaiting, func(reason string) bool {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		return strings.Contains(reason, &#34;chan&#34;) || strings.Contains(reason, &#34;sync&#34;) || strings.Contains(reason, &#34;select&#34;)
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	})
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// computePprofSyscall returns a computePprofFunc that generates a syscall pprof-like</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// profile (time spent in syscalls).</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>func computePprofSyscall() computePprofFunc {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	return makeComputePprofFunc(tracev2.GoSyscall, func(_ string) bool {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		return true
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	})
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// computePprofSched returns a computePprofFunc that generates a scheduler latency pprof-like profile</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// (time between a goroutine become runnable and actually scheduled for execution).</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>func computePprofSched() computePprofFunc {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	return makeComputePprofFunc(tracev2.GoRunnable, func(_ string) bool {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		return true
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	})
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// makeComputePprofFunc returns a computePprofFunc that generates a profile of time goroutines spend</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// in a particular state for the specified reasons.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>func makeComputePprofFunc(state tracev2.GoState, trackReason func(string) bool) computePprofFunc {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	return func(gToIntervals map[tracev2.GoID][]interval, events []tracev2.Event) ([]traceviewer.ProfileRecord, error) {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		stacks := newStackMap()
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		tracking := make(map[tracev2.GoID]*tracev2.Event)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		for i := range events {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			ev := &amp;events[i]
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			<span class="comment">// Filter out any non-state-transitions and events without stacks.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			if ev.Kind() != tracev2.EventStateTransition {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>				continue
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			stack := ev.Stack()
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			if stack == tracev2.NoStack {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>				continue
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			<span class="comment">// The state transition has to apply to a goroutine.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			st := ev.StateTransition()
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			if st.Resource.Kind != tracev2.ResourceGoroutine {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>				continue
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>			id := st.Resource.Goroutine()
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			_, new := st.Goroutine()
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			<span class="comment">// Check if we&#39;re tracking this goroutine.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			startEv := tracking[id]
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			if startEv == nil {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>				<span class="comment">// We&#39;re not. Start tracking if the new state</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>				<span class="comment">// matches what we want and the transition is</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>				<span class="comment">// for one of the reasons we care about.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>				if new == state &amp;&amp; trackReason(st.Reason) {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>					tracking[id] = ev
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>				}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>				continue
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			<span class="comment">// We&#39;re tracking this goroutine.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			if new == state {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>				<span class="comment">// We&#39;re tracking this goroutine, but it&#39;s just transitioning</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>				<span class="comment">// to the same state (this is a no-ip</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>				continue
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			<span class="comment">// The goroutine has transitioned out of the state we care about,</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			<span class="comment">// so remove it from tracking and record the stack.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			delete(tracking, id)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			overlapping := pprofOverlappingDuration(gToIntervals, id, interval{startEv.Time(), ev.Time()})
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			if overlapping &gt; 0 {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>				rec := stacks.getOrAdd(startEv.Stack())
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>				rec.Count++
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>				rec.Time += overlapping
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		return stacks.profile(), nil
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">// pprofOverlappingDuration returns the overlapping duration between</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">// the time intervals in gToIntervals and the specified event.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">// If gToIntervals is nil, this simply returns the event&#39;s duration.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>func pprofOverlappingDuration(gToIntervals map[tracev2.GoID][]interval, id tracev2.GoID, sample interval) time.Duration {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	if gToIntervals == nil { <span class="comment">// No filtering.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		return sample.duration()
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	intervals := gToIntervals[id]
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	if len(intervals) == 0 {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		return 0
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	var overlapping time.Duration
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	for _, i := range intervals {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		if o := i.overlap(sample); o &gt; 0 {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			overlapping += o
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	return overlapping
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// interval represents a time interval in the trace.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>type interval struct {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	start, end tracev2.Time
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>func (i interval) duration() time.Duration {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	return i.end.Sub(i.start)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>func (i1 interval) overlap(i2 interval) time.Duration {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">// Assume start1 &lt;= end1 and start2 &lt;= end2</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	if i1.end &lt; i2.start || i2.end &lt; i1.start {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		return 0
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	if i1.start &lt; i2.start { <span class="comment">// choose the later one</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		i1.start = i2.start
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	if i1.end &gt; i2.end { <span class="comment">// choose the earlier one</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		i1.end = i2.end
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	return i1.duration()
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">// pprofMaxStack is the extent of the deduplication we&#39;re willing to do.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// Because slices aren&#39;t comparable and we want to leverage maps for deduplication,</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">// we have to choose a fixed constant upper bound on the amount of frames we want</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">// to support. In practice this is fine because there&#39;s a maximum depth to these</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span><span class="comment">// stacks anyway.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>const pprofMaxStack = 128
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// stackMap is a map of tracev2.Stack to some value V.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>type stackMap struct {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// stacks contains the full list of stacks in the set, however</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">// it is insufficient for deduplication because tracev2.Stack</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	<span class="comment">// equality is only optimistic. If two tracev2.Stacks are equal,</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">// then they are guaranteed to be equal in content. If they are</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// not equal, then they might still be equal in content.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	stacks map[tracev2.Stack]*traceviewer.ProfileRecord
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// pcs is the source-of-truth for deduplication. It is a map of</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// the actual PCs in the stack to a tracev2.Stack.</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	pcs map[[pprofMaxStack]uint64]tracev2.Stack
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>func newStackMap() *stackMap {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	return &amp;stackMap{
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		stacks: make(map[tracev2.Stack]*traceviewer.ProfileRecord),
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		pcs:    make(map[[pprofMaxStack]uint64]tracev2.Stack),
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>func (m *stackMap) getOrAdd(stack tracev2.Stack) *traceviewer.ProfileRecord {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">// Fast path: check to see if this exact stack is already in the map.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	if rec, ok := m.stacks[stack]; ok {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		return rec
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	<span class="comment">// Slow path: the stack may still be in the map.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	<span class="comment">// Grab the stack&#39;s PCs as the source-of-truth.</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	var pcs [pprofMaxStack]uint64
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	pcsForStack(stack, &amp;pcs)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	<span class="comment">// Check the source-of-truth.</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	var rec *traceviewer.ProfileRecord
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	if existing, ok := m.pcs[pcs]; ok {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		<span class="comment">// In the map.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		rec = m.stacks[existing]
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		delete(m.stacks, existing)
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	} else {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		<span class="comment">// Not in the map.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		rec = new(traceviewer.ProfileRecord)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">// Insert regardless of whether we have a match in m.pcs.</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// Even if we have a match, we want to keep the newest version</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	<span class="comment">// of that stack, since we&#39;re much more likely tos see it again</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">// as we iterate through the trace linearly. Simultaneously, we</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">// are likely to never see the old stack again.</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	m.pcs[pcs] = stack
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	m.stacks[stack] = rec
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	return rec
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>func (m *stackMap) profile() []traceviewer.ProfileRecord {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	prof := make([]traceviewer.ProfileRecord, 0, len(m.stacks))
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	for stack, record := range m.stacks {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		rec := *record
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		i := 0
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		stack.Frames(func(frame tracev2.StackFrame) bool {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>			rec.Stack = append(rec.Stack, &amp;trace.Frame{
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>				PC:   frame.PC,
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>				Fn:   frame.Func,
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>				File: frame.File,
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>				Line: int(frame.Line),
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			})
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>			i++
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			<span class="comment">// Cut this off at pprofMaxStack because that&#39;s as far</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			<span class="comment">// as our deduplication goes.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			return i &lt; pprofMaxStack
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		})
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		prof = append(prof, rec)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	return prof
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// pcsForStack extracts the first pprofMaxStack PCs from stack into pcs.</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>func pcsForStack(stack tracev2.Stack, pcs *[pprofMaxStack]uint64) {
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	i := 0
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	stack.Frames(func(frame tracev2.StackFrame) bool {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		pcs[i] = frame.PC
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		i++
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		return i &lt; len(pcs)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	})
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
</pre><p><a href="pprof.go?m=text">View as plain text</a></p>

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
