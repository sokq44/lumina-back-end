<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/signal_amd64.go - Go Documentation Server</title>

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
<a href="signal_amd64.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">signal_amd64.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2013 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build amd64 &amp;&amp; (darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris)</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>func dumpregs(c *sigctxt) {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	print(&#34;rax    &#34;, hex(c.rax()), &#34;\n&#34;)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	print(&#34;rbx    &#34;, hex(c.rbx()), &#34;\n&#34;)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	print(&#34;rcx    &#34;, hex(c.rcx()), &#34;\n&#34;)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	print(&#34;rdx    &#34;, hex(c.rdx()), &#34;\n&#34;)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	print(&#34;rdi    &#34;, hex(c.rdi()), &#34;\n&#34;)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	print(&#34;rsi    &#34;, hex(c.rsi()), &#34;\n&#34;)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	print(&#34;rbp    &#34;, hex(c.rbp()), &#34;\n&#34;)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	print(&#34;rsp    &#34;, hex(c.rsp()), &#34;\n&#34;)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	print(&#34;r8     &#34;, hex(c.r8()), &#34;\n&#34;)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	print(&#34;r9     &#34;, hex(c.r9()), &#34;\n&#34;)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	print(&#34;r10    &#34;, hex(c.r10()), &#34;\n&#34;)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	print(&#34;r11    &#34;, hex(c.r11()), &#34;\n&#34;)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	print(&#34;r12    &#34;, hex(c.r12()), &#34;\n&#34;)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	print(&#34;r13    &#34;, hex(c.r13()), &#34;\n&#34;)
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	print(&#34;r14    &#34;, hex(c.r14()), &#34;\n&#34;)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	print(&#34;r15    &#34;, hex(c.r15()), &#34;\n&#34;)
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	print(&#34;rip    &#34;, hex(c.rip()), &#34;\n&#34;)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	print(&#34;rflags &#34;, hex(c.rflags()), &#34;\n&#34;)
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	print(&#34;cs     &#34;, hex(c.cs()), &#34;\n&#34;)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	print(&#34;fs     &#34;, hex(c.fs()), &#34;\n&#34;)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	print(&#34;gs     &#34;, hex(c.gs()), &#34;\n&#34;)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>func (c *sigctxt) sigpc() uintptr { return uintptr(c.rip()) }
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func (c *sigctxt) setsigpc(x uint64) { c.set_rip(x) }
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>func (c *sigctxt) sigsp() uintptr    { return uintptr(c.rsp()) }
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>func (c *sigctxt) siglr() uintptr    { return 0 }
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>func (c *sigctxt) fault() uintptr    { return uintptr(c.sigaddr()) }
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// preparePanic sets up the stack to look like a call to sigpanic.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>func (c *sigctxt) preparePanic(sig uint32, gp *g) {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// Work around Leopard bug that doesn&#39;t set FPE_INTDIV.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// Look at instruction to see if it is a divide.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// Not necessary in Snow Leopard (si_code will be != 0).</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	if GOOS == &#34;darwin&#34; &amp;&amp; sig == _SIGFPE &amp;&amp; gp.sigcode0 == 0 {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		pc := (*[4]byte)(unsafe.Pointer(gp.sigpc))
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		i := 0
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		if pc[i]&amp;0xF0 == 0x40 { <span class="comment">// 64-bit REX prefix</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			i++
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		} else if pc[i] == 0x66 { <span class="comment">// 16-bit instruction prefix</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			i++
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		if pc[i] == 0xF6 || pc[i] == 0xF7 {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			gp.sigcode0 = _FPE_INTDIV
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	pc := uintptr(c.rip())
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	sp := uintptr(c.rsp())
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// In case we are panicking from external code, we need to initialize</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// Go special registers. We inject sigpanic0 (instead of sigpanic),</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// which takes care of that.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	if shouldPushSigpanic(gp, pc, *(*uintptr)(unsafe.Pointer(sp))) {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		c.pushCall(abi.FuncPCABI0(sigpanic0), pc)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	} else {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		<span class="comment">// Not safe to push the call. Just clobber the frame.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		c.set_rip(uint64(abi.FuncPCABI0(sigpanic0)))
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>func (c *sigctxt) pushCall(targetPC, resumePC uintptr) {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// Make it look like we called target at resumePC.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	sp := uintptr(c.rsp())
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	sp -= goarch.PtrSize
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	*(*uintptr)(unsafe.Pointer(sp)) = resumePC
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	c.set_rsp(uint64(sp))
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	c.set_rip(uint64(targetPC))
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
</pre><p><a href="signal_amd64.go?m=text">View as plain text</a></p>

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
