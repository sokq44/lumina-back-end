<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/coverage/hooks.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../index.html">GoDoc</a></div>
<a href="hooks.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<a href="http://localhost:8080/src/runtime/coverage">coverage</a>/<span class="text-muted">hooks.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime/coverage">runtime/coverage</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2022 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package coverage
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import _ &#34;unsafe&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// initHook is invoked from the main package &#34;init&#34; routine in</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// programs built with &#34;-cover&#34;. This function is intended to be</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// called only by the compiler.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// If &#39;istest&#39; is false, it indicates we&#39;re building a regular program</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// (&#34;go build -cover ...&#34;), in which case we immediately try to write</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// out the meta-data file, and register emitCounterData as an exit</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// hook.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// If &#39;istest&#39; is true (indicating that the program in question is a</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// Go test binary), then we tentatively queue up both emitMetaData and</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// emitCounterData as exit hooks. In the normal case (e.g. regular &#34;go</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// test -cover&#34; run) the testmain.go boilerplate will run at the end</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// of the test, write out the coverage percentage, and then invoke</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// markProfileEmitted() to indicate that no more work needs to be</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// done. If however that call is never made, this is a sign that the</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// test binary is being used as a replacement binary for the tool</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// being tested, hence we do want to run exit hooks when the program</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// terminates.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>func initHook(istest bool) {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// Note: hooks are run in reverse registration order, so</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// register the counter data hook before the meta-data hook</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// (in the case where two hooks are needed).</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	runOnNonZeroExit := true
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	runtime_addExitHook(emitCounterData, runOnNonZeroExit)
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	if istest {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		runtime_addExitHook(emitMetaData, runOnNonZeroExit)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	} else {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		emitMetaData()
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//go:linkname runtime_addExitHook runtime.addExitHook</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func runtime_addExitHook(f func(), runOnNonZeroExit bool)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
</pre><p><a href="hooks.go?m=text">View as plain text</a></p>

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
