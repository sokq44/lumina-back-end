<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/debugcall.go - Go Documentation Server</title>

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
<a href="debugcall.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">debugcall.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2018 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Though the debug call function feature is not enabled on</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// ppc64, inserted ppc64 to avoid missing Go declaration error</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// for debugCallPanicked while building runtime.test</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//go:build amd64 || arm64 || ppc64le || ppc64</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>package runtime
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>import (
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>const (
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	debugCallSystemStack = &#34;executing on Go runtime stack&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	debugCallUnknownFunc = &#34;call from unknown function&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	debugCallRuntime     = &#34;call from within the Go runtime&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	debugCallUnsafePoint = &#34;call not at safe point&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>func debugCallV2()
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>func debugCallPanicked(val any)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// debugCallCheck checks whether it is safe to inject a debugger</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// function call with return PC pc. If not, it returns a string</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// explaining why.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>func debugCallCheck(pc uintptr) string {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// No user calls from the system stack.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	if getg() != getg().m.curg {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		return debugCallSystemStack
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	if sp := getcallersp(); !(getg().stack.lo &lt; sp &amp;&amp; sp &lt;= getg().stack.hi) {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		<span class="comment">// Fast syscalls (nanotime) and racecall switch to the</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		<span class="comment">// g0 stack without switching g. We can&#39;t safely make</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		<span class="comment">// a call in this state. (We can&#39;t even safely</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		<span class="comment">// systemstack.)</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		return debugCallSystemStack
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// Switch to the system stack to avoid overflowing the user</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// stack.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	var ret string
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		f := findfunc(pc)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		if !f.valid() {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			ret = debugCallUnknownFunc
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			return
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		name := funcname(f)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		switch name {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		case &#34;debugCall32&#34;,
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			&#34;debugCall64&#34;,
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			&#34;debugCall128&#34;,
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			&#34;debugCall256&#34;,
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			&#34;debugCall512&#34;,
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			&#34;debugCall1024&#34;,
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			&#34;debugCall2048&#34;,
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			&#34;debugCall4096&#34;,
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			&#34;debugCall8192&#34;,
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			&#34;debugCall16384&#34;,
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			&#34;debugCall32768&#34;,
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			&#34;debugCall65536&#34;:
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			<span class="comment">// These functions are allowed so that the debugger can initiate multiple function calls.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>			<span class="comment">// See: https://golang.org/cl/161137/</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			return
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		<span class="comment">// Disallow calls from the runtime. We could</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		<span class="comment">// potentially make this condition tighter (e.g., not</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		<span class="comment">// when locks are held), but there are enough tightly</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		<span class="comment">// coded sequences (e.g., defer handling) that it&#39;s</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		<span class="comment">// better to play it safe.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		if pfx := &#34;runtime.&#34;; len(name) &gt; len(pfx) &amp;&amp; name[:len(pfx)] == pfx {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			ret = debugCallRuntime
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			return
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		<span class="comment">// Check that this isn&#39;t an unsafe-point.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		if pc != f.entry() {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			pc--
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		up := pcdatavalue(f, abi.PCDATA_UnsafePoint, pc)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		if up != abi.UnsafePointSafe {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			<span class="comment">// Not at a safe point.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			ret = debugCallUnsafePoint
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	})
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	return ret
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// debugCallWrap starts a new goroutine to run a debug call and blocks</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// the calling goroutine. On the goroutine, it prepares to recover</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// panics from the debug call, and then calls the call dispatching</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// function at PC dispatch.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// This must be deeply nosplit because there are untyped values on the</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// stack from debugCallV2.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>func debugCallWrap(dispatch uintptr) {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	var lockedExt uint32
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	callerpc := getcallerpc()
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	gp := getg()
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// Lock ourselves to the OS thread.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// Debuggers rely on us running on the same thread until we get to</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// dispatch the function they asked as to.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// We&#39;re going to transfer this to the new G we just created.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	lockOSThread()
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// Create a new goroutine to execute the call on. Run this on</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// the system stack to avoid growing our stack.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		<span class="comment">// TODO(mknyszek): It would be nice to wrap these arguments in an allocated</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		<span class="comment">// closure and start the goroutine with that closure, but the compiler disallows</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		<span class="comment">// implicit closure allocation in the runtime.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		fn := debugCallWrap1
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		newg := newproc1(*(**funcval)(unsafe.Pointer(&amp;fn)), gp, callerpc)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		args := &amp;debugCallWrapArgs{
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			dispatch: dispatch,
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			callingG: gp,
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		newg.param = unsafe.Pointer(args)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// Transfer locked-ness to the new goroutine.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// Save lock state to restore later.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		mp := gp.m
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		if mp != gp.lockedm.ptr() {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			throw(&#34;inconsistent lockedm&#34;)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		<span class="comment">// Save the external lock count and clear it so</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		<span class="comment">// that it can&#39;t be unlocked from the debug call.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		<span class="comment">// Note: we already locked internally to the thread,</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		<span class="comment">// so if we were locked before we&#39;re still locked now.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		lockedExt = mp.lockedExt
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		mp.lockedExt = 0
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		mp.lockedg.set(newg)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		newg.lockedm.set(mp)
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		gp.lockedm = 0
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		<span class="comment">// Mark the calling goroutine as being at an async</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		<span class="comment">// safe-point, since it has a few conservative frames</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		<span class="comment">// at the bottom of the stack. This also prevents</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		<span class="comment">// stack shrinks.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		gp.asyncSafePoint = true
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// Stash newg away so we can execute it below (mcall&#39;s</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// closure can&#39;t capture anything).</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		gp.schedlink.set(newg)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	})
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// Switch to the new goroutine.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	mcall(func(gp *g) {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		<span class="comment">// Get newg.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		newg := gp.schedlink.ptr()
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		gp.schedlink = 0
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// Park the calling goroutine.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		trace := traceAcquire()
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		casGToWaiting(gp, _Grunning, waitReasonDebugCall)
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		if trace.ok() {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			trace.GoPark(traceBlockDebugCall, 1)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			traceRelease(trace)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		dropg()
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		<span class="comment">// Directly execute the new goroutine. The debug</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		<span class="comment">// protocol will continue on the new goroutine, so</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		<span class="comment">// it&#39;s important we not just let the scheduler do</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		<span class="comment">// this or it may resume a different goroutine.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		execute(newg, true)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	})
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// We&#39;ll resume here when the call returns.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// Restore locked state.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	mp := gp.m
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	mp.lockedExt = lockedExt
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	mp.lockedg.set(gp)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	gp.lockedm.set(mp)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// Undo the lockOSThread we did earlier.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	unlockOSThread()
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	gp.asyncSafePoint = false
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>type debugCallWrapArgs struct {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	dispatch uintptr
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	callingG *g
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">// debugCallWrap1 is the continuation of debugCallWrap on the callee</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">// goroutine.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>func debugCallWrap1() {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	gp := getg()
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	args := (*debugCallWrapArgs)(gp.param)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	dispatch, callingG := args.dispatch, args.callingG
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	gp.param = nil
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	<span class="comment">// Dispatch call and trap panics.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	debugCallWrap2(dispatch)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// Resume the caller goroutine.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	getg().schedlink.set(callingG)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	mcall(func(gp *g) {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		callingG := gp.schedlink.ptr()
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		gp.schedlink = 0
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		<span class="comment">// Unlock this goroutine from the M if necessary. The</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		<span class="comment">// calling G will relock.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		if gp.lockedm != 0 {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			gp.lockedm = 0
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			gp.m.lockedg = 0
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		<span class="comment">// Switch back to the calling goroutine. At some point</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		<span class="comment">// the scheduler will schedule us again and we&#39;ll</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		<span class="comment">// finish exiting.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		trace := traceAcquire()
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		casgstatus(gp, _Grunning, _Grunnable)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		if trace.ok() {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			trace.GoSched()
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			traceRelease(trace)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		dropg()
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		lock(&amp;sched.lock)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		globrunqput(gp)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		unlock(&amp;sched.lock)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		trace = traceAcquire()
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		casgstatus(callingG, _Gwaiting, _Grunnable)
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		if trace.ok() {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			trace.GoUnpark(callingG, 0)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			traceRelease(trace)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		execute(callingG, true)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	})
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func debugCallWrap2(dispatch uintptr) {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// Call the dispatch function and trap panics.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	var dispatchF func()
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	dispatchFV := funcval{dispatch}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	*(*unsafe.Pointer)(unsafe.Pointer(&amp;dispatchF)) = noescape(unsafe.Pointer(&amp;dispatchFV))
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	var ok bool
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	defer func() {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		if !ok {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			err := recover()
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			debugCallPanicked(err)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	}()
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	dispatchF()
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	ok = true
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
</pre><p><a href="debugcall.go?m=text">View as plain text</a></p>

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
