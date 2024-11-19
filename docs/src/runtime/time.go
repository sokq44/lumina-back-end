<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/time.go - Go Documentation Server</title>

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
<a href="time.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">time.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Time-related runtime and pieces of package time.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// Package time knows the layout of this structure.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// If this struct changes, adjust ../time/sleep.go:/runtimeTimer.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type timer struct {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// If this timer is on a heap, which P&#39;s heap it is on.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// puintptr rather than *p to match uintptr in the versions</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// of this struct defined in other packages.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	pp puintptr
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// Timer wakes up at when, and then at when+period, ... (period &gt; 0 only)</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// each time calling f(arg, now) in the timer goroutine, so f must be</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// a well-behaved function and not block.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// when must be positive on an active timer.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	when   int64
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	period int64
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	f      func(any, uintptr)
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	arg    any
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	seq    uintptr
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// What to set the when field to in timerModifiedXX status.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	nextwhen int64
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// The status field holds one of the values below.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	status atomic.Uint32
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// Code outside this file has to be careful in using a timer value.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// The pp, status, and nextwhen fields may only be used by code in this file.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// Code that creates a new timer value can set the when, period, f,</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// arg, and seq fields.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// A new timer value may be passed to addtimer (called by time.startTimer).</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// After doing that no fields may be touched.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// An active timer (one that has been passed to addtimer) may be</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// passed to deltimer (time.stopTimer), after which it is no longer an</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// active timer. It is an inactive timer.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// In an inactive timer the period, f, arg, and seq fields may be modified,</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// but not the when field.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// It&#39;s OK to just drop an inactive timer and let the GC collect it.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// It&#39;s not OK to pass an inactive timer to addtimer.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// Only newly allocated timer values may be passed to addtimer.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// An active timer may be passed to modtimer. No fields may be touched.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// It remains an active timer.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// An inactive timer may be passed to resettimer to turn into an</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// active timer with an updated when field.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// It&#39;s OK to pass a newly allocated timer value to resettimer.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// Timer operations are addtimer, deltimer, modtimer, resettimer,</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// cleantimers, adjusttimers, and runtimer.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// We don&#39;t permit calling addtimer/deltimer/modtimer/resettimer simultaneously,</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// but adjusttimers and runtimer can be called at the same time as any of those.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// Active timers live in heaps attached to P, in the timers field.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// Inactive timers live there too temporarily, until they are removed.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// addtimer:</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//   timerNoStatus   -&gt; timerWaiting</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//   anything else   -&gt; panic: invalid value</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// deltimer:</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//   timerWaiting         -&gt; timerModifying -&gt; timerDeleted</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//   timerModifiedEarlier -&gt; timerModifying -&gt; timerDeleted</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//   timerModifiedLater   -&gt; timerModifying -&gt; timerDeleted</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//   timerNoStatus        -&gt; do nothing</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">//   timerDeleted         -&gt; do nothing</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//   timerRemoving        -&gt; do nothing</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//   timerRemoved         -&gt; do nothing</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">//   timerRunning         -&gt; wait until status changes</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//   timerMoving          -&gt; wait until status changes</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//   timerModifying       -&gt; wait until status changes</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// modtimer:</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//   timerWaiting    -&gt; timerModifying -&gt; timerModifiedXX</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">//   timerModifiedXX -&gt; timerModifying -&gt; timerModifiedYY</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">//   timerNoStatus   -&gt; timerModifying -&gt; timerWaiting</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">//   timerRemoved    -&gt; timerModifying -&gt; timerWaiting</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">//   timerDeleted    -&gt; timerModifying -&gt; timerModifiedXX</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">//   timerRunning    -&gt; wait until status changes</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">//   timerMoving     -&gt; wait until status changes</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//   timerRemoving   -&gt; wait until status changes</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">//   timerModifying  -&gt; wait until status changes</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// cleantimers (looks in P&#39;s timer heap):</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">//   timerDeleted    -&gt; timerRemoving -&gt; timerRemoved</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">//   timerModifiedXX -&gt; timerMoving -&gt; timerWaiting</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// adjusttimers (looks in P&#39;s timer heap):</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">//   timerDeleted    -&gt; timerRemoving -&gt; timerRemoved</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">//   timerModifiedXX -&gt; timerMoving -&gt; timerWaiting</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// runtimer (looks in P&#39;s timer heap):</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">//   timerNoStatus   -&gt; panic: uninitialized timer</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">//   timerWaiting    -&gt; timerWaiting or</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">//   timerWaiting    -&gt; timerRunning -&gt; timerNoStatus or</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">//   timerWaiting    -&gt; timerRunning -&gt; timerWaiting</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//   timerModifying  -&gt; wait until status changes</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//   timerModifiedXX -&gt; timerMoving -&gt; timerWaiting</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">//   timerDeleted    -&gt; timerRemoving -&gt; timerRemoved</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">//   timerRunning    -&gt; panic: concurrent runtimer calls</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">//   timerRemoved    -&gt; panic: inconsistent timer heap</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">//   timerRemoving   -&gt; panic: inconsistent timer heap</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">//   timerMoving     -&gt; panic: inconsistent timer heap</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// Values for the timer status field.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>const (
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// Timer has no status set yet.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	timerNoStatus = iota
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// Waiting for timer to fire.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// The timer is in some P&#39;s heap.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	timerWaiting
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// Running the timer function.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// A timer will only have this status briefly.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	timerRunning
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// The timer is deleted and should be removed.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// It should not be run, but it is still in some P&#39;s heap.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	timerDeleted
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// The timer is being removed.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// The timer will only have this status briefly.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	timerRemoving
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// The timer has been stopped.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// It is not in any P&#39;s heap.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	timerRemoved
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// The timer is being modified.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// The timer will only have this status briefly.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	timerModifying
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// The timer has been modified to an earlier time.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// The new when value is in the nextwhen field.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// The timer is in some P&#39;s heap, possibly in the wrong place.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	timerModifiedEarlier
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// The timer has been modified to the same or a later time.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// The new when value is in the nextwhen field.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// The timer is in some P&#39;s heap, possibly in the wrong place.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	timerModifiedLater
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// The timer has been modified and is being moved.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// The timer will only have this status briefly.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	timerMoving
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>)
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// maxWhen is the maximum value for timer&#39;s when field.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>const maxWhen = 1&lt;&lt;63 - 1
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// verifyTimers can be set to true to add debugging checks that the</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// timer heaps are valid.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>const verifyTimers = false
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">// Package time APIs.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// Godoc uses the comments in package time, not these.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// time.now is implemented in assembly.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// timeSleep puts the current goroutine to sleep for at least ns nanoseconds.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">//go:linkname timeSleep time.Sleep</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>func timeSleep(ns int64) {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	if ns &lt;= 0 {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		return
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	gp := getg()
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	t := gp.timer
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if t == nil {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		t = new(timer)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		gp.timer = t
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	t.f = goroutineReady
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	t.arg = gp
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	t.nextwhen = nanotime() + ns
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	if t.nextwhen &lt; 0 { <span class="comment">// check for overflow.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		t.nextwhen = maxWhen
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	gopark(resetForSleep, unsafe.Pointer(t), waitReasonSleep, traceBlockSleep, 1)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// resetForSleep is called after the goroutine is parked for timeSleep.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// We can&#39;t call resettimer in timeSleep itself because if this is a short</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// sleep and there are many goroutines then the P can wind up running the</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// timer function, goroutineReady, before the goroutine has been parked.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>func resetForSleep(gp *g, ut unsafe.Pointer) bool {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	t := (*timer)(ut)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	resettimer(t, t.nextwhen)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	return true
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// startTimer adds t to the timer heap.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">//go:linkname startTimer time.startTimer</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>func startTimer(t *timer) {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	if raceenabled {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		racerelease(unsafe.Pointer(t))
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	addtimer(t)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// stopTimer stops a timer.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">// It reports whether t was stopped before being run.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">//go:linkname stopTimer time.stopTimer</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>func stopTimer(t *timer) bool {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	return deltimer(t)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// resetTimer resets an inactive timer, adding it to the heap.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// Reports whether the timer was modified before it was run.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">//go:linkname resetTimer time.resetTimer</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>func resetTimer(t *timer, when int64) bool {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	if raceenabled {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		racerelease(unsafe.Pointer(t))
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	return resettimer(t, when)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">// modTimer modifies an existing timer.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span><span class="comment">//go:linkname modTimer time.modTimer</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>func modTimer(t *timer, when, period int64, f func(any, uintptr), arg any, seq uintptr) {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	modtimer(t, when, period, f, arg, seq)
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// Go runtime.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// Ready the goroutine arg.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>func goroutineReady(arg any, seq uintptr) {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	goready(arg.(*g), 0)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span><span class="comment">// Note: this changes some unsynchronized operations to synchronized operations</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span><span class="comment">// addtimer adds a timer to the current P.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// This should only be called with a newly created timer.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// That avoids the risk of changing the when field of a timer in some P&#39;s heap,</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// which could cause the heap to become unsorted.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>func addtimer(t *timer) {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	<span class="comment">// when must be positive. A negative value will cause runtimer to</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">// overflow during its delta calculation and never expire other runtime</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// timers. Zero will cause checkTimers to fail to notice the timer.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	if t.when &lt;= 0 {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		throw(&#34;timer when must be positive&#34;)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	if t.period &lt; 0 {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		throw(&#34;timer period must be non-negative&#34;)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	if t.status.Load() != timerNoStatus {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		throw(&#34;addtimer called with initialized timer&#34;)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	t.status.Store(timerWaiting)
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	when := t.when
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	<span class="comment">// Disable preemption while using pp to avoid changing another P&#39;s heap.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	pp := getg().m.p.ptr()
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	lock(&amp;pp.timersLock)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	cleantimers(pp)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	doaddtimer(pp, t)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	unlock(&amp;pp.timersLock)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	wakeNetPoller(when)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	releasem(mp)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">// doaddtimer adds t to the current P&#39;s heap.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span><span class="comment">// The caller must have locked the timers for pp.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>func doaddtimer(pp *p, t *timer) {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// Timers rely on the network poller, so make sure the poller</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// has started.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	if netpollInited.Load() == 0 {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		netpollGenericInit()
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	if t.pp != 0 {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		throw(&#34;doaddtimer: P already set in timer&#34;)
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	t.pp.set(pp)
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	i := len(pp.timers)
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	pp.timers = append(pp.timers, t)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	siftupTimer(pp.timers, i)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	if t == pp.timers[0] {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		pp.timer0When.Store(t.when)
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	pp.numTimers.Add(1)
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// deltimer deletes the timer t. It may be on some other P, so we can&#39;t</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// actually remove it from the timers heap. We can only mark it as deleted.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">// It will be removed in due course by the P whose heap it is on.</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span><span class="comment">// Reports whether the timer was removed before it was run.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>func deltimer(t *timer) bool {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	for {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		switch s := t.status.Load(); s {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		case timerWaiting, timerModifiedLater:
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>			<span class="comment">// Prevent preemption while the timer is in timerModifying.</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			<span class="comment">// This could lead to a self-deadlock. See #38070.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			mp := acquirem()
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			if t.status.CompareAndSwap(s, timerModifying) {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>				<span class="comment">// Must fetch t.pp before changing status,</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>				<span class="comment">// as cleantimers in another goroutine</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>				<span class="comment">// can clear t.pp of a timerDeleted timer.</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>				tpp := t.pp.ptr()
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>				if !t.status.CompareAndSwap(timerModifying, timerDeleted) {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>					badTimer()
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>				}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>				releasem(mp)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>				tpp.deletedTimers.Add(1)
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>				<span class="comment">// Timer was not yet run.</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>				return true
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>			} else {
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>				releasem(mp)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>			}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		case timerModifiedEarlier:
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			<span class="comment">// Prevent preemption while the timer is in timerModifying.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>			<span class="comment">// This could lead to a self-deadlock. See #38070.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>			mp := acquirem()
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			if t.status.CompareAndSwap(s, timerModifying) {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>				<span class="comment">// Must fetch t.pp before setting status</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>				<span class="comment">// to timerDeleted.</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>				tpp := t.pp.ptr()
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>				if !t.status.CompareAndSwap(timerModifying, timerDeleted) {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>					badTimer()
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>				}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>				releasem(mp)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>				tpp.deletedTimers.Add(1)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>				<span class="comment">// Timer was not yet run.</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>				return true
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>			} else {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>				releasem(mp)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>			}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		case timerDeleted, timerRemoving, timerRemoved:
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			<span class="comment">// Timer was already run.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			return false
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		case timerRunning, timerMoving:
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			<span class="comment">// The timer is being run or moved, by a different P.</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			<span class="comment">// Wait for it to complete.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			osyield()
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		case timerNoStatus:
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>			<span class="comment">// Removing timer that was never added or</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			<span class="comment">// has already been run. Also see issue 21874.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			return false
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		case timerModifying:
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>			<span class="comment">// Simultaneous calls to deltimer and modtimer.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			<span class="comment">// Wait for the other call to complete.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			osyield()
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		default:
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			badTimer()
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span><span class="comment">// dodeltimer removes timer i from the current P&#39;s heap.</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span><span class="comment">// We are locked on the P when this is called.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span><span class="comment">// It returns the smallest changed index in pp.timers.</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span><span class="comment">// The caller must have locked the timers for pp.</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>func dodeltimer(pp *p, i int) int {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	if t := pp.timers[i]; t.pp.ptr() != pp {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		throw(&#34;dodeltimer: wrong P&#34;)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	} else {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		t.pp = 0
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	last := len(pp.timers) - 1
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	if i != last {
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		pp.timers[i] = pp.timers[last]
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	pp.timers[last] = nil
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	pp.timers = pp.timers[:last]
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	smallestChanged := i
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	if i != last {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		<span class="comment">// Moving to i may have moved the last timer to a new parent,</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		<span class="comment">// so sift up to preserve the heap guarantee.</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		smallestChanged = siftupTimer(pp.timers, i)
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		siftdownTimer(pp.timers, i)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	}
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	if i == 0 {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		updateTimer0When(pp)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	n := pp.numTimers.Add(-1)
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	if n == 0 {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		<span class="comment">// If there are no timers, then clearly none are modified.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		pp.timerModifiedEarliest.Store(0)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	return smallestChanged
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span><span class="comment">// dodeltimer0 removes timer 0 from the current P&#39;s heap.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span><span class="comment">// We are locked on the P when this is called.</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span><span class="comment">// It reports whether it saw no problems due to races.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span><span class="comment">// The caller must have locked the timers for pp.</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>func dodeltimer0(pp *p) {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	if t := pp.timers[0]; t.pp.ptr() != pp {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		throw(&#34;dodeltimer0: wrong P&#34;)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	} else {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		t.pp = 0
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	last := len(pp.timers) - 1
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	if last &gt; 0 {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		pp.timers[0] = pp.timers[last]
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	pp.timers[last] = nil
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	pp.timers = pp.timers[:last]
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	if last &gt; 0 {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		siftdownTimer(pp.timers, 0)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	updateTimer0When(pp)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	n := pp.numTimers.Add(-1)
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	if n == 0 {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		<span class="comment">// If there are no timers, then clearly none are modified.</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		pp.timerModifiedEarliest.Store(0)
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span><span class="comment">// modtimer modifies an existing timer.</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span><span class="comment">// This is called by the netpoll code or time.Ticker.Reset or time.Timer.Reset.</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span><span class="comment">// Reports whether the timer was modified before it was run.</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>func modtimer(t *timer, when, period int64, f func(any, uintptr), arg any, seq uintptr) bool {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	if when &lt;= 0 {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		throw(&#34;timer when must be positive&#34;)
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	if period &lt; 0 {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		throw(&#34;timer period must be non-negative&#34;)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	status := uint32(timerNoStatus)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	wasRemoved := false
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	var pending bool
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	var mp *m
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>loop:
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	for {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		switch status = t.status.Load(); status {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		case timerWaiting, timerModifiedEarlier, timerModifiedLater:
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			<span class="comment">// Prevent preemption while the timer is in timerModifying.</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>			<span class="comment">// This could lead to a self-deadlock. See #38070.</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			mp = acquirem()
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			if t.status.CompareAndSwap(status, timerModifying) {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>				pending = true <span class="comment">// timer not yet run</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>				break loop
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>			releasem(mp)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		case timerNoStatus, timerRemoved:
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			<span class="comment">// Prevent preemption while the timer is in timerModifying.</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>			<span class="comment">// This could lead to a self-deadlock. See #38070.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			mp = acquirem()
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>			<span class="comment">// Timer was already run and t is no longer in a heap.</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>			<span class="comment">// Act like addtimer.</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>			if t.status.CompareAndSwap(status, timerModifying) {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>				wasRemoved = true
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>				pending = false <span class="comment">// timer already run or stopped</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>				break loop
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			releasem(mp)
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		case timerDeleted:
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			<span class="comment">// Prevent preemption while the timer is in timerModifying.</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			<span class="comment">// This could lead to a self-deadlock. See #38070.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			mp = acquirem()
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			if t.status.CompareAndSwap(status, timerModifying) {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>				t.pp.ptr().deletedTimers.Add(-1)
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>				pending = false <span class="comment">// timer already stopped</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>				break loop
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			releasem(mp)
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		case timerRunning, timerRemoving, timerMoving:
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>			<span class="comment">// The timer is being run or moved, by a different P.</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			<span class="comment">// Wait for it to complete.</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>			osyield()
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		case timerModifying:
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>			<span class="comment">// Multiple simultaneous calls to modtimer.</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>			<span class="comment">// Wait for the other call to complete.</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>			osyield()
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		default:
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			badTimer()
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	}
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	t.period = period
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	t.f = f
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	t.arg = arg
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	t.seq = seq
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	if wasRemoved {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		t.when = when
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		pp := getg().m.p.ptr()
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		lock(&amp;pp.timersLock)
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		doaddtimer(pp, t)
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		unlock(&amp;pp.timersLock)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		if !t.status.CompareAndSwap(timerModifying, timerWaiting) {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			badTimer()
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		releasem(mp)
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		wakeNetPoller(when)
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	} else {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		<span class="comment">// The timer is in some other P&#39;s heap, so we can&#39;t change</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		<span class="comment">// the when field. If we did, the other P&#39;s heap would</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		<span class="comment">// be out of order. So we put the new when value in the</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		<span class="comment">// nextwhen field, and let the other P set the when field</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		<span class="comment">// when it is prepared to resort the heap.</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		t.nextwhen = when
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		newStatus := uint32(timerModifiedLater)
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		if when &lt; t.when {
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>			newStatus = timerModifiedEarlier
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		tpp := t.pp.ptr()
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		if newStatus == timerModifiedEarlier {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>			updateTimerModifiedEarliest(tpp, when)
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		<span class="comment">// Set the new status of the timer.</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		if !t.status.CompareAndSwap(timerModifying, newStatus) {
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			badTimer()
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		releasem(mp)
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		<span class="comment">// If the new status is earlier, wake up the poller.</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		if newStatus == timerModifiedEarlier {
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			wakeNetPoller(when)
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	}
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	return pending
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span><span class="comment">// resettimer resets the time when a timer should fire.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span><span class="comment">// If used for an inactive timer, the timer will become active.</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span><span class="comment">// This should be called instead of addtimer if the timer value has been,</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// or may have been, used previously.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">// Reports whether the timer was modified before it was run.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>func resettimer(t *timer, when int64) bool {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	return modtimer(t, when, t.period, t.f, t.arg, t.seq)
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span><span class="comment">// cleantimers cleans up the head of the timer queue. This speeds up</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span><span class="comment">// programs that create and delete timers; leaving them in the heap</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span><span class="comment">// slows down addtimer. Reports whether no timer problems were found.</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span><span class="comment">// The caller must have locked the timers for pp.</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>func cleantimers(pp *p) {
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	gp := getg()
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	for {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		if len(pp.timers) == 0 {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>			return
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		<span class="comment">// This loop can theoretically run for a while, and because</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		<span class="comment">// it is holding timersLock it cannot be preempted.</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		<span class="comment">// If someone is trying to preempt us, just return.</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		<span class="comment">// We can clean the timers later.</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		if gp.preemptStop {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>			return
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		}
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		t := pp.timers[0]
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		if t.pp.ptr() != pp {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>			throw(&#34;cleantimers: bad p&#34;)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		switch s := t.status.Load(); s {
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		case timerDeleted:
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>			if !t.status.CompareAndSwap(s, timerRemoving) {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>				continue
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>			}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>			dodeltimer0(pp)
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>			if !t.status.CompareAndSwap(timerRemoving, timerRemoved) {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>				badTimer()
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>			}
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>			pp.deletedTimers.Add(-1)
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		case timerModifiedEarlier, timerModifiedLater:
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>			if !t.status.CompareAndSwap(s, timerMoving) {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>				continue
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>			<span class="comment">// Now we can change the when field.</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>			t.when = t.nextwhen
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>			<span class="comment">// Move t to the right position.</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>			dodeltimer0(pp)
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>			doaddtimer(pp, t)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>			if !t.status.CompareAndSwap(timerMoving, timerWaiting) {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>				badTimer()
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>			}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		default:
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>			<span class="comment">// Head of timers does not need adjustment.</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			return
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	}
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span><span class="comment">// moveTimers moves a slice of timers to pp. The slice has been taken</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span><span class="comment">// from a different P.</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span><span class="comment">// This is currently called when the world is stopped, but the caller</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span><span class="comment">// is expected to have locked the timers for pp.</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>func moveTimers(pp *p, timers []*timer) {
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	for _, t := range timers {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	loop:
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		for {
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>			switch s := t.status.Load(); s {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>			case timerWaiting:
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>				if !t.status.CompareAndSwap(s, timerMoving) {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>					continue
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>				}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>				t.pp = 0
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>				doaddtimer(pp, t)
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>				if !t.status.CompareAndSwap(timerMoving, timerWaiting) {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>					badTimer()
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>				}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>				break loop
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>			case timerModifiedEarlier, timerModifiedLater:
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>				if !t.status.CompareAndSwap(s, timerMoving) {
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>					continue
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>				}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>				t.when = t.nextwhen
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>				t.pp = 0
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>				doaddtimer(pp, t)
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>				if !t.status.CompareAndSwap(timerMoving, timerWaiting) {
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>					badTimer()
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>				}
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>				break loop
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			case timerDeleted:
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>				if !t.status.CompareAndSwap(s, timerRemoved) {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>					continue
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>				}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>				t.pp = 0
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>				<span class="comment">// We no longer need this timer in the heap.</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>				break loop
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>			case timerModifying:
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>				<span class="comment">// Loop until the modification is complete.</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>				osyield()
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>			case timerNoStatus, timerRemoved:
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>				<span class="comment">// We should not see these status values in a timers heap.</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>				badTimer()
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>			case timerRunning, timerRemoving, timerMoving:
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>				<span class="comment">// Some other P thinks it owns this timer,</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>				<span class="comment">// which should not happen.</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>				badTimer()
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>			default:
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>				badTimer()
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>			}
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		}
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	}
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>}
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span><span class="comment">// adjusttimers looks through the timers in the current P&#39;s heap for</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span><span class="comment">// any timers that have been modified to run earlier, and puts them in</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span><span class="comment">// the correct place in the heap. While looking for those timers,</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span><span class="comment">// it also moves timers that have been modified to run later,</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span><span class="comment">// and removes deleted timers. The caller must have locked the timers for pp.</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>func adjusttimers(pp *p, now int64) {
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	<span class="comment">// If we haven&#39;t yet reached the time of the first timerModifiedEarlier</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	<span class="comment">// timer, don&#39;t do anything. This speeds up programs that adjust</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	<span class="comment">// a lot of timers back and forth if the timers rarely expire.</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	<span class="comment">// We&#39;ll postpone looking through all the adjusted timers until</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	<span class="comment">// one would actually expire.</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	first := pp.timerModifiedEarliest.Load()
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	if first == 0 || first &gt; now {
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		if verifyTimers {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>			verifyTimerHeap(pp)
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		}
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		return
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	}
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	<span class="comment">// We are going to clear all timerModifiedEarlier timers.</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	pp.timerModifiedEarliest.Store(0)
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	var moved []*timer
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	for i := 0; i &lt; len(pp.timers); i++ {
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		t := pp.timers[i]
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		if t.pp.ptr() != pp {
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>			throw(&#34;adjusttimers: bad p&#34;)
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		}
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		switch s := t.status.Load(); s {
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		case timerDeleted:
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>			if t.status.CompareAndSwap(s, timerRemoving) {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>				changed := dodeltimer(pp, i)
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>				if !t.status.CompareAndSwap(timerRemoving, timerRemoved) {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>					badTimer()
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>				}
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>				pp.deletedTimers.Add(-1)
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>				<span class="comment">// Go back to the earliest changed heap entry.</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>				<span class="comment">// &#34;- 1&#34; because the loop will add 1.</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>				i = changed - 1
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>			}
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		case timerModifiedEarlier, timerModifiedLater:
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			if t.status.CompareAndSwap(s, timerMoving) {
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>				<span class="comment">// Now we can change the when field.</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>				t.when = t.nextwhen
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>				<span class="comment">// Take t off the heap, and hold onto it.</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>				<span class="comment">// We don&#39;t add it back yet because the</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>				<span class="comment">// heap manipulation could cause our</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>				<span class="comment">// loop to skip some other timer.</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>				changed := dodeltimer(pp, i)
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>				moved = append(moved, t)
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>				<span class="comment">// Go back to the earliest changed heap entry.</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>				<span class="comment">// &#34;- 1&#34; because the loop will add 1.</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>				i = changed - 1
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>			}
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>		case timerNoStatus, timerRunning, timerRemoving, timerRemoved, timerMoving:
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>			badTimer()
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		case timerWaiting:
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>			<span class="comment">// OK, nothing to do.</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		case timerModifying:
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>			<span class="comment">// Check again after modification is complete.</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>			osyield()
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>			i--
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		default:
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>			badTimer()
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>		}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	}
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	if len(moved) &gt; 0 {
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		addAdjustedTimers(pp, moved)
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	}
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	if verifyTimers {
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>		verifyTimerHeap(pp)
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	}
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>}
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span><span class="comment">// addAdjustedTimers adds any timers we adjusted in adjusttimers</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span><span class="comment">// back to the timer heap.</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>func addAdjustedTimers(pp *p, moved []*timer) {
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	for _, t := range moved {
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		doaddtimer(pp, t)
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		if !t.status.CompareAndSwap(timerMoving, timerWaiting) {
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>			badTimer()
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		}
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	}
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>}
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span><span class="comment">// nobarrierWakeTime looks at P&#39;s timers and returns the time when we</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span><span class="comment">// should wake up the netpoller. It returns 0 if there are no timers.</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span><span class="comment">// This function is invoked when dropping a P, and must run without</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span><span class="comment">// any write barriers.</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>func nobarrierWakeTime(pp *p) int64 {
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	next := pp.timer0When.Load()
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	nextAdj := pp.timerModifiedEarliest.Load()
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	if next == 0 || (nextAdj != 0 &amp;&amp; nextAdj &lt; next) {
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		next = nextAdj
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	}
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	return next
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span><span class="comment">// runtimer examines the first timer in timers. If it is ready based on now,</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span><span class="comment">// it runs the timer and removes or updates it.</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span><span class="comment">// Returns 0 if it ran a timer, -1 if there are no more timers, or the time</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span><span class="comment">// when the first timer should run.</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span><span class="comment">// The caller must have locked the timers for pp.</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span><span class="comment">// If a timer is run, this will temporarily unlock the timers.</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>func runtimer(pp *p, now int64) int64 {
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>	for {
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		t := pp.timers[0]
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>		if t.pp.ptr() != pp {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>			throw(&#34;runtimer: bad p&#34;)
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		}
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		switch s := t.status.Load(); s {
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>		case timerWaiting:
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>			if t.when &gt; now {
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>				<span class="comment">// Not ready to run.</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>				return t.when
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>			}
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>			if !t.status.CompareAndSwap(s, timerRunning) {
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>				continue
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>			}
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>			<span class="comment">// Note that runOneTimer may temporarily unlock</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>			<span class="comment">// pp.timersLock.</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>			runOneTimer(pp, t, now)
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>			return 0
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>		case timerDeleted:
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>			if !t.status.CompareAndSwap(s, timerRemoving) {
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>				continue
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>			}
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>			dodeltimer0(pp)
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			if !t.status.CompareAndSwap(timerRemoving, timerRemoved) {
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>				badTimer()
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>			}
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>			pp.deletedTimers.Add(-1)
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>			if len(pp.timers) == 0 {
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>				return -1
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>			}
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		case timerModifiedEarlier, timerModifiedLater:
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>			if !t.status.CompareAndSwap(s, timerMoving) {
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>				continue
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>			}
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>			t.when = t.nextwhen
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>			dodeltimer0(pp)
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>			doaddtimer(pp, t)
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>			if !t.status.CompareAndSwap(timerMoving, timerWaiting) {
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>				badTimer()
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>			}
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		case timerModifying:
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>			<span class="comment">// Wait for modification to complete.</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>			osyield()
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		case timerNoStatus, timerRemoved:
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>			<span class="comment">// Should not see a new or inactive timer on the heap.</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>			badTimer()
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		case timerRunning, timerRemoving, timerMoving:
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>			<span class="comment">// These should only be set when timers are locked,</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>			<span class="comment">// and we didn&#39;t do it.</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>			badTimer()
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		default:
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>			badTimer()
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		}
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	}
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>}
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span><span class="comment">// runOneTimer runs a single timer.</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span><span class="comment">// The caller must have locked the timers for pp.</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span><span class="comment">// This will temporarily unlock the timers while running the timer function.</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>func runOneTimer(pp *p, t *timer, now int64) {
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	if raceenabled {
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>		ppcur := getg().m.p.ptr()
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>		if ppcur.timerRaceCtx == 0 {
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			ppcur.timerRaceCtx = racegostart(abi.FuncPCABIInternal(runtimer) + sys.PCQuantum)
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		}
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		raceacquirectx(ppcur.timerRaceCtx, unsafe.Pointer(t))
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	}
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	f := t.f
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	arg := t.arg
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	seq := t.seq
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	if t.period &gt; 0 {
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>		<span class="comment">// Leave in heap but adjust next time to fire.</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		delta := t.when - now
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>		t.when += t.period * (1 + -delta/t.period)
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		if t.when &lt; 0 { <span class="comment">// check for overflow.</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>			t.when = maxWhen
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>		}
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		siftdownTimer(pp.timers, 0)
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		if !t.status.CompareAndSwap(timerRunning, timerWaiting) {
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>			badTimer()
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		}
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		updateTimer0When(pp)
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	} else {
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		<span class="comment">// Remove from heap.</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>		dodeltimer0(pp)
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>		if !t.status.CompareAndSwap(timerRunning, timerNoStatus) {
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>			badTimer()
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>		}
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	}
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	if raceenabled {
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>		<span class="comment">// Temporarily use the current P&#39;s racectx for g0.</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>		gp := getg()
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>		if gp.racectx != 0 {
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>			throw(&#34;runOneTimer: unexpected racectx&#34;)
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>		}
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>		gp.racectx = gp.m.p.ptr().timerRaceCtx
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	}
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	unlock(&amp;pp.timersLock)
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	f(arg, seq)
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	lock(&amp;pp.timersLock)
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>	if raceenabled {
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>		gp := getg()
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>		gp.racectx = 0
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	}
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>}
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span><span class="comment">// clearDeletedTimers removes all deleted timers from the P&#39;s timer heap.</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span><span class="comment">// This is used to avoid clogging up the heap if the program</span>
<span id="L897" class="ln">   897&nbsp;&nbsp;</span><span class="comment">// starts a lot of long-running timers and then stops them.</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span><span class="comment">// For example, this can happen via context.WithTimeout.</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span><span class="comment">// This is the only function that walks through the entire timer heap,</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span><span class="comment">// other than moveTimers which only runs when the world is stopped.</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span><span class="comment">// The caller must have locked the timers for pp.</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>func clearDeletedTimers(pp *p) {
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	<span class="comment">// We are going to clear all timerModifiedEarlier timers.</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	<span class="comment">// Do this now in case new ones show up while we are looping.</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	pp.timerModifiedEarliest.Store(0)
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	cdel := int32(0)
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	to := 0
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>	changedHeap := false
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	timers := pp.timers
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>nextTimer:
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	for _, t := range timers {
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		for {
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>			switch s := t.status.Load(); s {
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>			case timerWaiting:
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>				if changedHeap {
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>					timers[to] = t
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>					siftupTimer(timers, to)
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>				}
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>				to++
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>				continue nextTimer
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>			case timerModifiedEarlier, timerModifiedLater:
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>				if t.status.CompareAndSwap(s, timerMoving) {
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>					t.when = t.nextwhen
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>					timers[to] = t
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>					siftupTimer(timers, to)
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>					to++
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>					changedHeap = true
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>					if !t.status.CompareAndSwap(timerMoving, timerWaiting) {
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>						badTimer()
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>					}
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>					continue nextTimer
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>				}
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>			case timerDeleted:
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>				if t.status.CompareAndSwap(s, timerRemoving) {
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>					t.pp = 0
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>					cdel++
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>					if !t.status.CompareAndSwap(timerRemoving, timerRemoved) {
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>						badTimer()
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>					}
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>					changedHeap = true
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>					continue nextTimer
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>				}
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>			case timerModifying:
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>				<span class="comment">// Loop until modification complete.</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>				osyield()
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>			case timerNoStatus, timerRemoved:
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>				<span class="comment">// We should not see these status values in a timer heap.</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>				badTimer()
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			case timerRunning, timerRemoving, timerMoving:
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>				<span class="comment">// Some other P thinks it owns this timer,</span>
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>				<span class="comment">// which should not happen.</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>				badTimer()
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>			default:
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>				badTimer()
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>			}
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		}
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>	}
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>	<span class="comment">// Set remaining slots in timers slice to nil,</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>	<span class="comment">// so that the timer values can be garbage collected.</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	for i := to; i &lt; len(timers); i++ {
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>		timers[i] = nil
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>	}
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>	pp.deletedTimers.Add(-cdel)
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>	pp.numTimers.Add(-cdel)
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>	timers = timers[:to]
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>	pp.timers = timers
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	updateTimer0When(pp)
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>	if verifyTimers {
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>		verifyTimerHeap(pp)
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>	}
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>}
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span><span class="comment">// verifyTimerHeap verifies that the timer heap is in a valid state.</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span><span class="comment">// This is only for debugging, and is only called if verifyTimers is true.</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span><span class="comment">// The caller must have locked the timers.</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>func verifyTimerHeap(pp *p) {
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>	for i, t := range pp.timers {
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>		if i == 0 {
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>			<span class="comment">// First timer has no parent.</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>			continue
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>		}
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>		<span class="comment">// The heap is 4-ary. See siftupTimer and siftdownTimer.</span>
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>		p := (i - 1) / 4
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>		if t.when &lt; pp.timers[p].when {
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>			print(&#34;bad timer heap at &#34;, i, &#34;: &#34;, p, &#34;: &#34;, pp.timers[p].when, &#34;, &#34;, i, &#34;: &#34;, t.when, &#34;\n&#34;)
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>			throw(&#34;bad timer heap&#34;)
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		}
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	}
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	if numTimers := int(pp.numTimers.Load()); len(pp.timers) != numTimers {
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>		println(&#34;timer heap len&#34;, len(pp.timers), &#34;!= numTimers&#34;, numTimers)
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>		throw(&#34;bad timer heap len&#34;)
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	}
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>}
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span><span class="comment">// updateTimer0When sets the P&#39;s timer0When field.</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span><span class="comment">// The caller must have locked the timers for pp.</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>func updateTimer0When(pp *p) {
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	if len(pp.timers) == 0 {
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>		pp.timer0When.Store(0)
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	} else {
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>		pp.timer0When.Store(pp.timers[0].when)
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>	}
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>}
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span><span class="comment">// updateTimerModifiedEarliest updates the recorded nextwhen field of the</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span><span class="comment">// earlier timerModifiedEarier value.</span>
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span><span class="comment">// The timers for pp will not be locked.</span>
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>func updateTimerModifiedEarliest(pp *p, nextwhen int64) {
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	for {
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>		old := pp.timerModifiedEarliest.Load()
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>		if old != 0 &amp;&amp; old &lt; nextwhen {
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>			return
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		}
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		if pp.timerModifiedEarliest.CompareAndSwap(old, nextwhen) {
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>			return
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		}
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	}
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>}
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span><span class="comment">// timeSleepUntil returns the time when the next timer should fire. Returns</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span><span class="comment">// maxWhen if there are no timers.</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span><span class="comment">// This is only called by sysmon and checkdead.</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>func timeSleepUntil() int64 {
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>	next := int64(maxWhen)
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	<span class="comment">// Prevent allp slice changes. This is like retake.</span>
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>	lock(&amp;allpLock)
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>	for _, pp := range allp {
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>		if pp == nil {
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>			<span class="comment">// This can happen if procresize has grown</span>
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>			<span class="comment">// allp but not yet created new Ps.</span>
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>			continue
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		}
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>		w := pp.timer0When.Load()
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>		if w != 0 &amp;&amp; w &lt; next {
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>			next = w
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>		}
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>		w = pp.timerModifiedEarliest.Load()
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>		if w != 0 &amp;&amp; w &lt; next {
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>			next = w
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>		}
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>	}
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>	unlock(&amp;allpLock)
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>	return next
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>}
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span><span class="comment">// Heap maintenance algorithms.</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span><span class="comment">// These algorithms check for slice index errors manually.</span>
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span><span class="comment">// Slice index error can happen if the program is using racy</span>
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span><span class="comment">// access to timers. We don&#39;t want to panic here, because</span>
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span><span class="comment">// it will cause the program to crash with a mysterious</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span><span class="comment">// &#34;panic holding locks&#34; message. Instead, we panic while not</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span><span class="comment">// holding a lock.</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span><span class="comment">// siftupTimer puts the timer at position i in the right place</span>
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span><span class="comment">// in the heap by moving it up toward the top of the heap.</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span><span class="comment">// It returns the smallest changed index.</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>func siftupTimer(t []*timer, i int) int {
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	if i &gt;= len(t) {
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>		badTimer()
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	}
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	when := t[i].when
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	if when &lt;= 0 {
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>		badTimer()
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	}
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>	tmp := t[i]
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	for i &gt; 0 {
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>		p := (i - 1) / 4 <span class="comment">// parent</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>		if when &gt;= t[p].when {
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>			break
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>		}
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>		t[i] = t[p]
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>		i = p
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>	}
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	if tmp != t[i] {
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>		t[i] = tmp
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>	}
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>	return i
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>}
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span><span class="comment">// siftdownTimer puts the timer at position i in the right place</span>
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span><span class="comment">// in the heap by moving it down toward the bottom of the heap.</span>
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>func siftdownTimer(t []*timer, i int) {
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>	n := len(t)
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>	if i &gt;= n {
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>		badTimer()
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	}
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>	when := t[i].when
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	if when &lt;= 0 {
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>		badTimer()
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>	}
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>	tmp := t[i]
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>	for {
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>		c := i*4 + 1 <span class="comment">// left child</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>		c3 := c + 2  <span class="comment">// mid child</span>
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>		if c &gt;= n {
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>			break
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>		}
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>		w := t[c].when
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>		if c+1 &lt; n &amp;&amp; t[c+1].when &lt; w {
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>			w = t[c+1].when
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>			c++
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>		}
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>		if c3 &lt; n {
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>			w3 := t[c3].when
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>			if c3+1 &lt; n &amp;&amp; t[c3+1].when &lt; w3 {
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>				w3 = t[c3+1].when
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>				c3++
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>			}
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>			if w3 &lt; w {
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>				w = w3
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>				c = c3
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>			}
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>		}
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>		if w &gt;= when {
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>			break
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>		}
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>		t[i] = t[c]
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>		i = c
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>	}
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>	if tmp != t[i] {
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>		t[i] = tmp
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>	}
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>}
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span><span class="comment">// badTimer is called if the timer data structures have been corrupted,</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span><span class="comment">// presumably due to racy use by the program. We panic here rather than</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span><span class="comment">// panicking due to invalid slice access while holding locks.</span>
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span><span class="comment">// See issue #25686.</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>func badTimer() {
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	throw(&#34;timer data corruption&#34;)
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>}
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>
</pre><p><a href="time.go?m=text">View as plain text</a></p>

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
