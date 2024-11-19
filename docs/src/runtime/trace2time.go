<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/trace2time.go - Go Documentation Server</title>

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
<a href="trace2time.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">trace2time.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Trace time and clock.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package runtime
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import &#34;internal/goarch&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// Timestamps in trace are produced through either nanotime or cputicks</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// and divided by traceTimeDiv. nanotime is used everywhere except on</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// platforms where osHasLowResClock is true, because the system clock</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// isn&#39;t granular enough to get useful information out of a trace in</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// many cases.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// This makes absolute values of timestamp diffs smaller, and so they are</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// encoded in fewer bytes.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// The target resolution in all cases is 64 nanoseconds.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// This is based on the fact that fundamentally the execution tracer won&#39;t emit</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// events more frequently than roughly every 200 ns or so, because that&#39;s roughly</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// how long it takes to call through the scheduler.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// We could be more aggressive and bump this up to 128 ns while still getting</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// useful data, but the extra bit doesn&#39;t save us that much and the headroom is</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// nice to have.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// Hitting this target resolution is easy in the nanotime case: just pick a</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// division of 64. In the cputicks case it&#39;s a bit more complex.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// For x86, on a 3 GHz machine, we&#39;d want to divide by 3*64 to hit our target.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// To keep the division operation efficient, we round that up to 4*64, or 256.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// Given what cputicks represents, we use this on all other platforms except</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// for PowerPC.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// The suggested increment frequency for PowerPC&#39;s time base register is</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// 512 MHz according to Power ISA v2.07 section 6.2, so we use 32 on ppc64</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// and ppc64le.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>const traceTimeDiv = (1-osHasLowResClockInt)*64 + osHasLowResClockInt*(256-224*(goarch.IsPpc64|goarch.IsPpc64le))
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// traceTime represents a timestamp for the trace.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>type traceTime uint64
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// traceClockNow returns a monotonic timestamp. The clock this function gets</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// the timestamp from is specific to tracing, and shouldn&#39;t be mixed with other</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// clock sources.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// nosplit because it&#39;s called from exitsyscall, which is nosplit.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>func traceClockNow() traceTime {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	if osHasLowResClock {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		return traceTime(cputicks() / traceTimeDiv)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	return traceTime(nanotime() / traceTimeDiv)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// traceClockUnitsPerSecond estimates the number of trace clock units per</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// second that elapse.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>func traceClockUnitsPerSecond() uint64 {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	if osHasLowResClock {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re using cputicks as our clock, so we need a real estimate.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		return uint64(ticksPerSecond() / traceTimeDiv)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// Our clock is nanotime, so it&#39;s just the constant time division.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// (trace clock units / nanoseconds) * (1e9 nanoseconds / 1 second)</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	return uint64(1.0 / float64(traceTimeDiv) * 1e9)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// traceFrequency writes a batch with a single EvFrequency event.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// freq is the number of trace clock units per second.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>func traceFrequency(gen uintptr) {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	w := unsafeTraceWriter(gen, nil)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// Ensure we have a place to write to.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	w, _ = w.ensure(1 + traceBytesPerNumber <span class="comment">/* traceEvFrequency + frequency */</span>)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// Write out the string.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	w.byte(byte(traceEvFrequency))
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	w.varint(traceClockUnitsPerSecond())
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// Immediately flush the buffer.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		lock(&amp;trace.lock)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		traceBufFlush(w.traceBuf, gen)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		unlock(&amp;trace.lock)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	})
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
</pre><p><a href="trace2time.go?m=text">View as plain text</a></p>

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
