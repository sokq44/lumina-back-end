<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/exithook.go - Go Documentation Server</title>

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
<a href="exithook.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">exithook.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2022 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// addExitHook registers the specified function &#39;f&#39; to be run at</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// program termination (e.g. when someone invokes os.Exit(), or when</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// main.main returns). Hooks are run in reverse order of registration:</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// first hook added is the last one run.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// CAREFUL: the expectation is that addExitHook should only be called</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// from a safe context (e.g. not an error/panic path or signal</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// handler, preemption enabled, allocation allowed, write barriers</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// allowed, etc), and that the exit function &#39;f&#39; will be invoked under</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// similar circumstances. That is the say, we are expecting that &#39;f&#39;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// uses normal / high-level Go code as opposed to one of the more</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// restricted dialects used for the trickier parts of the runtime.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>func addExitHook(f func(), runOnNonZeroExit bool) {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	exitHooks.hooks = append(exitHooks.hooks, exitHook{f: f, runOnNonZeroExit: runOnNonZeroExit})
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>}
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// exitHook stores a function to be run on program exit, registered</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// by the utility runtime.addExitHook.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>type exitHook struct {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	f                func() <span class="comment">// func to run</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	runOnNonZeroExit bool   <span class="comment">// whether to run on non-zero exit code</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// exitHooks stores state related to hook functions registered to</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// run when program execution terminates.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>var exitHooks struct {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	hooks            []exitHook
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	runningExitHooks bool
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// runExitHooks runs any registered exit hook functions (funcs</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// previously registered using runtime.addExitHook). Here &#39;exitCode&#39;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// is the status code being passed to os.Exit, or zero if the program</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// is terminating normally without calling os.Exit.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>func runExitHooks(exitCode int) {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	if exitHooks.runningExitHooks {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		throw(&#34;internal error: exit hook invoked exit&#34;)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	exitHooks.runningExitHooks = true
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	runExitHook := func(f func()) (caughtPanic bool) {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		defer func() {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			if x := recover(); x != nil {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>				caughtPanic = true
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		}()
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		f()
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		return
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	finishPageTrace()
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	for i := range exitHooks.hooks {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		h := exitHooks.hooks[len(exitHooks.hooks)-i-1]
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		if exitCode != 0 &amp;&amp; !h.runOnNonZeroExit {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			continue
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		if caughtPanic := runExitHook(h.f); caughtPanic {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			throw(&#34;internal error: exit hook invoked panic&#34;)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	exitHooks.hooks = nil
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	exitHooks.runningExitHooks = false
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
</pre><p><a href="exithook.go?m=text">View as plain text</a></p>

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
