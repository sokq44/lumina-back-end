<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/runtime.go - Go Documentation Server</title>

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
<a href="runtime.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">runtime.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//go:generate go run wincallback.go</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//go:generate go run mkduff.go</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//go:generate go run mkfastlog2table.go</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//go:generate go run mklockrank.go -o lockrank.go</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>var ticks ticksType
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>type ticksType struct {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// lock protects access to start* and val.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	lock       mutex
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	startTicks int64
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	startTime  int64
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	val        atomic.Int64
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>}
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// init initializes ticks to maximize the chance that we have a good ticksPerSecond reference.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// Must not run concurrently with ticksPerSecond.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>func (t *ticksType) init() {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	lock(&amp;ticks.lock)
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	t.startTime = nanotime()
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	t.startTicks = cputicks()
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	unlock(&amp;ticks.lock)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// minTimeForTicksPerSecond is the minimum elapsed time we require to consider our ticksPerSecond</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// measurement to be of decent enough quality for profiling.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// There&#39;s a linear relationship here between minimum time and error from the true value.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// The error from the true ticks-per-second in a linux/amd64 VM seems to be:</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// -   1 ms -&gt; ~0.02% error</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// -   5 ms -&gt; ~0.004% error</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// -  10 ms -&gt; ~0.002% error</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// -  50 ms -&gt; ~0.0003% error</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// - 100 ms -&gt; ~0.0001% error</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// We&#39;re willing to take 0.004% error here, because ticksPerSecond is intended to be used for</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// converting durations, not timestamps. Durations are usually going to be much larger, and so</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// the tiny error doesn&#39;t matter. The error is definitely going to be a problem when trying to</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// use this for timestamps, as it&#39;ll make those timestamps much less likely to line up.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>const minTimeForTicksPerSecond = 5_000_000*(1-osHasLowResClockInt) + 100_000_000*osHasLowResClockInt
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// ticksPerSecond returns a conversion rate between the cputicks clock and the nanotime clock.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// Note: Clocks are hard. Using this as an actual conversion rate for timestamps is ill-advised</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// and should be avoided when possible. Use only for durations, where a tiny error term isn&#39;t going</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// to make a meaningful difference in even a 1ms duration. If an accurate timestamp is needed,</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// use nanotime instead. (The entire Windows platform is a broad exception to this rule, where nanotime</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// produces timestamps on such a coarse granularity that the error from this conversion is actually</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// preferable.)</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// The strategy for computing the conversion rate is to write down nanotime and cputicks as</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// early in process startup as possible. From then, we just need to wait until we get values</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// from nanotime that we can use (some platforms have a really coarse system time granularity).</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// We require some amount of time to pass to ensure that the conversion rate is fairly accurate</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// in aggregate. But because we compute this rate lazily, there&#39;s a pretty good chance a decent</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// amount of time has passed by the time we get here.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// Must be called from a normal goroutine context (running regular goroutine with a P).</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// Called by runtime/pprof in addition to runtime code.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// TODO(mknyszek): This doesn&#39;t account for things like CPU frequency scaling. Consider</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// a more sophisticated and general approach in the future.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>func ticksPerSecond() int64 {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// Get the conversion rate if we&#39;ve already computed it.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	r := ticks.val.Load()
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	if r != 0 {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		return r
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// Compute the conversion rate.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	for {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		lock(&amp;ticks.lock)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		r = ticks.val.Load()
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		if r != 0 {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			unlock(&amp;ticks.lock)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			return r
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		<span class="comment">// Grab the current time in both clocks.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		nowTime := nanotime()
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		nowTicks := cputicks()
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		<span class="comment">// See if we can use these times.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		if nowTicks &gt; ticks.startTicks &amp;&amp; nowTime-ticks.startTime &gt; minTimeForTicksPerSecond {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			<span class="comment">// Perform the calculation with floats. We don&#39;t want to risk overflow.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			r = int64(float64(nowTicks-ticks.startTicks) * 1e9 / float64(nowTime-ticks.startTime))
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			if r == 0 {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>				<span class="comment">// Zero is both a sentinel value and it would be bad if callers used this as</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>				<span class="comment">// a divisor. We tried out best, so just make it 1.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>				r++
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>			}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			ticks.val.Store(r)
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			unlock(&amp;ticks.lock)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			break
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		unlock(&amp;ticks.lock)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		<span class="comment">// Sleep in one millisecond increments until we have a reliable time.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		timeSleep(1_000_000)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	return r
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>var envs []string
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>var argslice []string
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">//go:linkname syscall_runtime_envs syscall.runtime_envs</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>func syscall_runtime_envs() []string { return append([]string{}, envs...) }
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">//go:linkname syscall_Getpagesize syscall.Getpagesize</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func syscall_Getpagesize() int { return int(physPageSize) }
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">//go:linkname os_runtime_args os.runtime_args</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>func os_runtime_args() []string { return append([]string{}, argslice...) }
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">//go:linkname syscall_Exit syscall.Exit</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func syscall_Exit(code int) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	exit(int32(code))
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>var godebugDefault string
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>var godebugUpdate atomic.Pointer[func(string, string)]
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>var godebugEnv atomic.Pointer[string] <span class="comment">// set by parsedebugvars</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>var godebugNewIncNonDefault atomic.Pointer[func(string) func()]
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">//go:linkname godebug_setUpdate internal/godebug.setUpdate</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>func godebug_setUpdate(update func(string, string)) {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	p := new(func(string, string))
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	*p = update
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	godebugUpdate.Store(p)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	godebugNotify(false)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">//go:linkname godebug_setNewIncNonDefault internal/godebug.setNewIncNonDefault</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>func godebug_setNewIncNonDefault(newIncNonDefault func(string) func()) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	p := new(func(string) func())
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	*p = newIncNonDefault
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	godebugNewIncNonDefault.Store(p)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// A godebugInc provides access to internal/godebug&#39;s IncNonDefault function</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// for a given GODEBUG setting.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// Calls before internal/godebug registers itself are dropped on the floor.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>type godebugInc struct {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	name string
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	inc  atomic.Pointer[func()]
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>func (g *godebugInc) IncNonDefault() {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	inc := g.inc.Load()
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	if inc == nil {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		newInc := godebugNewIncNonDefault.Load()
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		if newInc == nil {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			return
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		inc = new(func())
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		*inc = (*newInc)(g.name)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		if raceenabled {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			racereleasemerge(unsafe.Pointer(&amp;g.inc))
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		if !g.inc.CompareAndSwap(nil, inc) {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			inc = g.inc.Load()
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	if raceenabled {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		raceacquire(unsafe.Pointer(&amp;g.inc))
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	(*inc)()
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>func godebugNotify(envChanged bool) {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	update := godebugUpdate.Load()
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	var env string
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	if p := godebugEnv.Load(); p != nil {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		env = *p
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	if envChanged {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		reparsedebugvars(env)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	if update != nil {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		(*update)(godebugDefault, env)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">//go:linkname syscall_runtimeSetenv syscall.runtimeSetenv</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>func syscall_runtimeSetenv(key, value string) {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	setenv_c(key, value)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	if key == &#34;GODEBUG&#34; {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		p := new(string)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		*p = value
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		godebugEnv.Store(p)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		godebugNotify(true)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">//go:linkname syscall_runtimeUnsetenv syscall.runtimeUnsetenv</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>func syscall_runtimeUnsetenv(key string) {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	unsetenv_c(key)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	if key == &#34;GODEBUG&#34; {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		godebugEnv.Store(nil)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		godebugNotify(true)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">// writeErrStr writes a string to descriptor 2.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>func writeErrStr(s string) {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	write(2, unsafe.Pointer(unsafe.StringData(s)), int32(len(s)))
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// auxv is populated on relevant platforms but defined here for all platforms</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// so x/sys/cpu can assume the getAuxv symbol exists without keeping its list</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// of auxv-using GOOS build tags in sync.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// It contains an even number of elements, (tag, value) pairs.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>var auxv []uintptr
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>func getAuxv() []uintptr { return auxv } <span class="comment">// accessed from x/sys/cpu; see issue 57336</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
</pre><p><a href="runtime.go?m=text">View as plain text</a></p>

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
