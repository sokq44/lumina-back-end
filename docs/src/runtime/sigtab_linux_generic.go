<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/sigtab_linux_generic.go - Go Documentation Server</title>

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
<a href="sigtab_linux_generic.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">sigtab_linux_generic.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build !mips &amp;&amp; !mipsle &amp;&amp; !mips64 &amp;&amp; !mips64le &amp;&amp; linux</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>var sigtable = [...]sigTabT{
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	<span class="comment">/* 0 */</span> {0, &#34;SIGNONE: no trap&#34;},
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	<span class="comment">/* 1 */</span> {_SigNotify + _SigKill, &#34;SIGHUP: terminal line hangup&#34;},
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	<span class="comment">/* 2 */</span> {_SigNotify + _SigKill, &#34;SIGINT: interrupt&#34;},
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	<span class="comment">/* 3 */</span> {_SigNotify + _SigThrow, &#34;SIGQUIT: quit&#34;},
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	<span class="comment">/* 4 */</span> {_SigThrow + _SigUnblock, &#34;SIGILL: illegal instruction&#34;},
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	<span class="comment">/* 5 */</span> {_SigThrow + _SigUnblock, &#34;SIGTRAP: trace trap&#34;},
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	<span class="comment">/* 6 */</span> {_SigNotify + _SigThrow, &#34;SIGABRT: abort&#34;},
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">/* 7 */</span> {_SigPanic + _SigUnblock, &#34;SIGBUS: bus error&#34;},
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">/* 8 */</span> {_SigPanic + _SigUnblock, &#34;SIGFPE: floating-point exception&#34;},
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">/* 9 */</span> {0, &#34;SIGKILL: kill&#34;},
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">/* 10 */</span> {_SigNotify, &#34;SIGUSR1: user-defined signal 1&#34;},
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">/* 11 */</span> {_SigPanic + _SigUnblock, &#34;SIGSEGV: segmentation violation&#34;},
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">/* 12 */</span> {_SigNotify, &#34;SIGUSR2: user-defined signal 2&#34;},
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">/* 13 */</span> {_SigNotify, &#34;SIGPIPE: write to broken pipe&#34;},
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">/* 14 */</span> {_SigNotify, &#34;SIGALRM: alarm clock&#34;},
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">/* 15 */</span> {_SigNotify + _SigKill, &#34;SIGTERM: termination&#34;},
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">/* 16 */</span> {_SigThrow + _SigUnblock, &#34;SIGSTKFLT: stack fault&#34;},
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">/* 17 */</span> {_SigNotify + _SigUnblock + _SigIgn, &#34;SIGCHLD: child status has changed&#34;},
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">/* 18 */</span> {_SigNotify + _SigDefault + _SigIgn, &#34;SIGCONT: continue&#34;},
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">/* 19 */</span> {0, &#34;SIGSTOP: stop, unblockable&#34;},
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">/* 20 */</span> {_SigNotify + _SigDefault + _SigIgn, &#34;SIGTSTP: keyboard stop&#34;},
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">/* 21 */</span> {_SigNotify + _SigDefault + _SigIgn, &#34;SIGTTIN: background read from tty&#34;},
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">/* 22 */</span> {_SigNotify + _SigDefault + _SigIgn, &#34;SIGTTOU: background write to tty&#34;},
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">/* 23 */</span> {_SigNotify + _SigIgn, &#34;SIGURG: urgent condition on socket&#34;},
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">/* 24 */</span> {_SigNotify, &#34;SIGXCPU: cpu limit exceeded&#34;},
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">/* 25 */</span> {_SigNotify, &#34;SIGXFSZ: file size limit exceeded&#34;},
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">/* 26 */</span> {_SigNotify, &#34;SIGVTALRM: virtual alarm clock&#34;},
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">/* 27 */</span> {_SigNotify + _SigUnblock, &#34;SIGPROF: profiling alarm clock&#34;},
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">/* 28 */</span> {_SigNotify + _SigIgn, &#34;SIGWINCH: window size change&#34;},
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">/* 29 */</span> {_SigNotify, &#34;SIGIO: i/o now possible&#34;},
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">/* 30 */</span> {_SigNotify, &#34;SIGPWR: power failure restart&#34;},
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">/* 31 */</span> {_SigThrow, &#34;SIGSYS: bad system call&#34;},
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">/* 32 */</span> {_SigSetStack + _SigUnblock, &#34;signal 32&#34;}, <span class="comment">/* SIGCANCEL; see issue 6997 */</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">/* 33 */</span> {_SigSetStack + _SigUnblock, &#34;signal 33&#34;}, <span class="comment">/* SIGSETXID; see issues 3871, 9400, 12498 */</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">/* 34 */</span> {_SigSetStack + _SigUnblock, &#34;signal 34&#34;}, <span class="comment">/* musl SIGSYNCCALL; see issue 39343 */</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">/* 35 */</span> {_SigNotify, &#34;signal 35&#34;},
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">/* 36 */</span> {_SigNotify, &#34;signal 36&#34;},
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">/* 37 */</span> {_SigNotify, &#34;signal 37&#34;},
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">/* 38 */</span> {_SigNotify, &#34;signal 38&#34;},
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">/* 39 */</span> {_SigNotify, &#34;signal 39&#34;},
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">/* 40 */</span> {_SigNotify, &#34;signal 40&#34;},
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">/* 41 */</span> {_SigNotify, &#34;signal 41&#34;},
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">/* 42 */</span> {_SigNotify, &#34;signal 42&#34;},
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">/* 43 */</span> {_SigNotify, &#34;signal 43&#34;},
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">/* 44 */</span> {_SigNotify, &#34;signal 44&#34;},
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">/* 45 */</span> {_SigNotify, &#34;signal 45&#34;},
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">/* 46 */</span> {_SigNotify, &#34;signal 46&#34;},
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">/* 47 */</span> {_SigNotify, &#34;signal 47&#34;},
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">/* 48 */</span> {_SigNotify, &#34;signal 48&#34;},
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">/* 49 */</span> {_SigNotify, &#34;signal 49&#34;},
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">/* 50 */</span> {_SigNotify, &#34;signal 50&#34;},
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">/* 51 */</span> {_SigNotify, &#34;signal 51&#34;},
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">/* 52 */</span> {_SigNotify, &#34;signal 52&#34;},
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">/* 53 */</span> {_SigNotify, &#34;signal 53&#34;},
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">/* 54 */</span> {_SigNotify, &#34;signal 54&#34;},
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">/* 55 */</span> {_SigNotify, &#34;signal 55&#34;},
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">/* 56 */</span> {_SigNotify, &#34;signal 56&#34;},
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">/* 57 */</span> {_SigNotify, &#34;signal 57&#34;},
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">/* 58 */</span> {_SigNotify, &#34;signal 58&#34;},
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">/* 59 */</span> {_SigNotify, &#34;signal 59&#34;},
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">/* 60 */</span> {_SigNotify, &#34;signal 60&#34;},
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">/* 61 */</span> {_SigNotify, &#34;signal 61&#34;},
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">/* 62 */</span> {_SigNotify, &#34;signal 62&#34;},
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">/* 63 */</span> {_SigNotify, &#34;signal 63&#34;},
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">/* 64 */</span> {_SigNotify, &#34;signal 64&#34;},
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
</pre><p><a href="sigtab_linux_generic.go?m=text">View as plain text</a></p>

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
