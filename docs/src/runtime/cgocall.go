<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/cgocall.go - Go Documentation Server</title>

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
<a href="cgocall.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">cgocall.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Cgo call and callback support.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// To call into the C function f from Go, the cgo-generated code calls</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// runtime.cgocall(_cgo_Cfunc_f, frame), where _cgo_Cfunc_f is a</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// gcc-compiled function written by cgo.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// runtime.cgocall (below) calls entersyscall so as not to block</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// other goroutines or the garbage collector, and then calls</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// runtime.asmcgocall(_cgo_Cfunc_f, frame).</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// runtime.asmcgocall (in asm_$GOARCH.s) switches to the m-&gt;g0 stack</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// (assumed to be an operating system-allocated stack, so safe to run</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// gcc-compiled code on) and calls _cgo_Cfunc_f(frame).</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// _cgo_Cfunc_f invokes the actual C function f with arguments</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// taken from the frame structure, records the results in the frame,</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// and returns to runtime.asmcgocall.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// After it regains control, runtime.asmcgocall switches back to the</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// original g (m-&gt;curg)&#39;s stack and returns to runtime.cgocall.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// After it regains control, runtime.cgocall calls exitsyscall, which blocks</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// until this m can run Go code without violating the $GOMAXPROCS limit,</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// and then unlocks g from m.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// The above description skipped over the possibility of the gcc-compiled</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// function f calling back into Go. If that happens, we continue down</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// the rabbit hole during the execution of f.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// To make it possible for gcc-compiled C code to call a Go function p.GoF,</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// cgo writes a gcc-compiled function named GoF (not p.GoF, since gcc doesn&#39;t</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// know about packages).  The gcc-compiled C function f calls GoF.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// GoF initializes &#34;frame&#34;, a structure containing all of its</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// arguments and slots for p.GoF&#39;s results. It calls</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// crosscall2(_cgoexp_GoF, frame, framesize, ctxt) using the gcc ABI.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// crosscall2 (in cgo/asm_$GOARCH.s) is a four-argument adapter from</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// the gcc function call ABI to the gc function call ABI. At this</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// point we&#39;re in the Go runtime, but we&#39;re still running on m.g0&#39;s</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// stack and outside the $GOMAXPROCS limit. crosscall2 calls</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// runtime.cgocallback(_cgoexp_GoF, frame, ctxt) using the gc ABI.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// (crosscall2&#39;s framesize argument is no longer used, but there&#39;s one</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// case where SWIG calls crosscall2 directly and expects to pass this</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// argument. See _cgo_panic.)</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// runtime.cgocallback (in asm_$GOARCH.s) switches from m.g0&#39;s stack</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// to the original g (m.curg)&#39;s stack, on which it calls</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// runtime.cgocallbackg(_cgoexp_GoF, frame, ctxt). As part of the</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// stack switch, runtime.cgocallback saves the current SP as</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// m.g0.sched.sp, so that any use of m.g0&#39;s stack during the execution</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// of the callback will be done below the existing stack frames.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// Before overwriting m.g0.sched.sp, it pushes the old value on the</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// m.g0 stack, so that it can be restored later.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// runtime.cgocallbackg (below) is now running on a real goroutine</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// stack (not an m.g0 stack).  First it calls runtime.exitsyscall, which will</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// block until the $GOMAXPROCS limit allows running this goroutine.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// Once exitsyscall has returned, it is safe to do things like call the memory</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// allocator or invoke the Go callback function.  runtime.cgocallbackg</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// first defers a function to unwind m.g0.sched.sp, so that if p.GoF</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// panics, m.g0.sched.sp will be restored to its old value: the m.g0 stack</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// and the m.curg stack will be unwound in lock step.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// Then it calls _cgoexp_GoF(frame).</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// _cgoexp_GoF, which was generated by cmd/cgo, unpacks the arguments</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// from frame, calls p.GoF, writes the results back to frame, and</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// returns. Now we start unwinding this whole process.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// runtime.cgocallbackg pops but does not execute the deferred</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// function to unwind m.g0.sched.sp, calls runtime.entersyscall, and</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// returns to runtime.cgocallback.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// After it regains control, runtime.cgocallback switches back to</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// m.g0&#39;s stack (the pointer is still in m.g0.sched.sp), restores the old</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// m.g0.sched.sp value from the stack, and returns to crosscall2.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// crosscall2 restores the callee-save registers for gcc and returns</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// to GoF, which unpacks any result values and returns to f.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>package runtime
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>import (
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	&#34;internal/goexperiment&#34;
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// Addresses collected in a cgo backtrace when crashing.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// Length must match arg.Max in x_cgo_callers in runtime/cgo/gcc_traceback.c.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>type cgoCallers [32]uintptr
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// argset matches runtime/cgo/linux_syscall.c:argset_t</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>type argset struct {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	args   unsafe.Pointer
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	retval uintptr
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// wrapper for syscall package to call cgocall for libc (cgo) calls.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//go:linkname syscall_cgocaller syscall.cgocaller</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">//go:uintptrescapes</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>func syscall_cgocaller(fn unsafe.Pointer, args ...uintptr) uintptr {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	as := argset{args: unsafe.Pointer(&amp;args[0])}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	cgocall(fn, unsafe.Pointer(&amp;as))
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	return as.retval
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>var ncgocall uint64 <span class="comment">// number of cgo calls in total for dead m</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// Call from Go to C.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// This must be nosplit because it&#39;s used for syscalls on some</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// platforms. Syscalls may have untyped arguments on the stack, so</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// it&#39;s not safe to grow or scan the stack.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func cgocall(fn, arg unsafe.Pointer) int32 {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	if !iscgo &amp;&amp; GOOS != &#34;solaris&#34; &amp;&amp; GOOS != &#34;illumos&#34; &amp;&amp; GOOS != &#34;windows&#34; {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		throw(&#34;cgocall unavailable&#34;)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if fn == nil {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		throw(&#34;cgocall nil&#34;)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if raceenabled {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		racereleasemerge(unsafe.Pointer(&amp;racecgosync))
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	mp := getg().m
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	mp.ncgocall++
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// Reset traceback.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	mp.cgoCallers[0] = 0
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">// Announce we are entering a system call</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// so that the scheduler knows to create another</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// M to run goroutines while we are in the</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// foreign code.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// The call to asmcgocall is guaranteed not to</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// grow the stack and does not allocate memory,</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// so it is safe to call while &#34;in a system call&#34;, outside</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// the $GOMAXPROCS accounting.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// fn may call back into Go code, in which case we&#39;ll exit the</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// &#34;system call&#34;, run the Go code (which may grow the stack),</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// and then re-enter the &#34;system call&#34; reusing the PC and SP</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// saved by entersyscall here.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	entersyscall()
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// Tell asynchronous preemption that we&#39;re entering external</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// code. We do this after entersyscall because this may block</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// and cause an async preemption to fail, but at this point a</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// sync preemption will succeed (though this is not a matter</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// of correctness).</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	osPreemptExtEnter(mp)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	mp.incgo = true
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// We use ncgo as a check during execution tracing for whether there is</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// any C on the call stack, which there will be after this point. If</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// there isn&#39;t, we can use frame pointer unwinding to collect call</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// stacks efficiently. This will be the case for the first Go-to-C call</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// on a stack, so it&#39;s preferable to update it here, after we emit a</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// trace event in entersyscall above.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	mp.ncgo++
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	errno := asmcgocall(fn, arg)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// Update accounting before exitsyscall because exitsyscall may</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// reschedule us on to a different M.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	mp.incgo = false
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	mp.ncgo--
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	osPreemptExtExit(mp)
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	exitsyscall()
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// Note that raceacquire must be called only after exitsyscall has</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// wired this M to a P.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	if raceenabled {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		raceacquire(unsafe.Pointer(&amp;racecgosync))
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// From the garbage collector&#39;s perspective, time can move</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// backwards in the sequence above. If there&#39;s a callback into</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// Go code, GC will see this function at the call to</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// asmcgocall. When the Go call later returns to C, the</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// syscall PC/SP is rolled back and the GC sees this function</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// back at the call to entersyscall. Normally, fn and arg</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// would be live at entersyscall and dead at asmcgocall, so if</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// time moved backwards, GC would see these arguments as dead</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// and then live. Prevent these undead arguments from crashing</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// GC by forcing them to stay live across this time warp.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	KeepAlive(fn)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	KeepAlive(arg)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	KeepAlive(mp)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	return errno
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// Set or reset the system stack bounds for a callback on sp.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// Must be nosplit because it is called by needm prior to fully initializing</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// the M.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>func callbackUpdateSystemStack(mp *m, sp uintptr, signal bool) {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	g0 := mp.g0
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	if sp &gt; g0.stack.lo &amp;&amp; sp &lt;= g0.stack.hi {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// Stack already in bounds, nothing to do.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		return
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	if mp.ncgo &gt; 0 {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		<span class="comment">// ncgo &gt; 0 indicates that this M was in Go further up the stack</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		<span class="comment">// (it called C and is now receiving a callback). It is not</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		<span class="comment">// safe for the C call to change the stack out from under us.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		<span class="comment">// Note that this case isn&#39;t possible for signal == true, as</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		<span class="comment">// that is always passing a new M from needm.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		<span class="comment">// Stack is bogus, but reset the bounds anyway so we can print.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		hi := g0.stack.hi
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		lo := g0.stack.lo
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		g0.stack.hi = sp + 1024
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		g0.stack.lo = sp - 32*1024
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		g0.stackguard0 = g0.stack.lo + stackGuard
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		g0.stackguard1 = g0.stackguard0
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		print(&#34;M &#34;, mp.id, &#34; procid &#34;, mp.procid, &#34; runtime: cgocallback with sp=&#34;, hex(sp), &#34; out of bounds [&#34;, hex(lo), &#34;, &#34;, hex(hi), &#34;]&#34;)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		exit(2)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// This M does not have Go further up the stack. However, it may have</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// previously called into Go, initializing the stack bounds. Between</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// that call returning and now the stack may have changed (perhaps the</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">// C thread is running a coroutine library). We need to update the</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	<span class="comment">// stack bounds for this case.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	<span class="comment">// Set the stack bounds to match the current stack. If we don&#39;t</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// actually know how big the stack is, like we don&#39;t know how big any</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// scheduling stack is, but we assume there&#39;s at least 32 kB. If we</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// can get a more accurate stack bound from pthread, use that, provided</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// it actually contains SP..</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	g0.stack.hi = sp + 1024
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	g0.stack.lo = sp - 32*1024
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	if !signal &amp;&amp; _cgo_getstackbound != nil {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		<span class="comment">// Don&#39;t adjust if called from the signal handler.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		<span class="comment">// We are on the signal stack, not the pthread stack.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		<span class="comment">// (We could get the stack bounds from sigaltstack, but</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		<span class="comment">// we&#39;re getting out of the signal handler very soon</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		<span class="comment">// anyway. Not worth it.)</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		var bounds [2]uintptr
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		asmcgocall(_cgo_getstackbound, unsafe.Pointer(&amp;bounds))
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		<span class="comment">// getstackbound is an unsupported no-op on Windows.</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		<span class="comment">// Don&#39;t use these bounds if they don&#39;t contain SP. Perhaps we</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		<span class="comment">// were called by something not using the standard thread</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		<span class="comment">// stack.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		if bounds[0] != 0 &amp;&amp; sp &gt; bounds[0] &amp;&amp; sp &lt;= bounds[1] {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			g0.stack.lo = bounds[0]
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			g0.stack.hi = bounds[1]
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	g0.stackguard0 = g0.stack.lo + stackGuard
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	g0.stackguard1 = g0.stackguard0
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span><span class="comment">// Call from C back to Go. fn must point to an ABIInternal Go entry-point.</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>func cgocallbackg(fn, frame unsafe.Pointer, ctxt uintptr) {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	gp := getg()
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	if gp != gp.m.curg {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		println(&#34;runtime: bad g in cgocallback&#34;)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		exit(2)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	sp := gp.m.g0.sched.sp <span class="comment">// system sp saved by cgocallback.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	callbackUpdateSystemStack(gp.m, sp, false)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// The call from C is on gp.m&#39;s g0 stack, so we must ensure</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// that we stay on that M. We have to do this before calling</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// exitsyscall, since it would otherwise be free to move us to</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">// a different M. The call to unlockOSThread is in this function</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// after cgocallbackg1, or in the case of panicking, in unwindm.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	lockOSThread()
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	checkm := gp.m
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">// Save current syscall parameters, so m.syscall can be</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">// used again if callback decide to make syscall.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	syscall := gp.m.syscall
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	<span class="comment">// entersyscall saves the caller&#39;s SP to allow the GC to trace the Go</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">// stack. However, since we&#39;re returning to an earlier stack frame and</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	<span class="comment">// need to pair with the entersyscall() call made by cgocall, we must</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	<span class="comment">// save syscall* and let reentersyscall restore them.</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	savedsp := unsafe.Pointer(gp.syscallsp)
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	savedpc := gp.syscallpc
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	exitsyscall() <span class="comment">// coming out of cgo call</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	gp.m.incgo = false
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	if gp.m.isextra {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		gp.m.isExtraInC = false
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	osPreemptExtExit(gp.m)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	if gp.nocgocallback {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		panic(&#34;runtime: function marked with #cgo nocallback called back into Go&#34;)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	cgocallbackg1(fn, frame, ctxt)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	<span class="comment">// At this point we&#39;re about to call unlockOSThread.</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// The following code must not change to a different m.</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">// This is enforced by checking incgo in the schedule function.</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	gp.m.incgo = true
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	unlockOSThread()
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	if gp.m.isextra {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		gp.m.isExtraInC = true
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	if gp.m != checkm {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		throw(&#34;m changed unexpectedly in cgocallbackg&#34;)
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	osPreemptExtEnter(gp.m)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">// going back to cgo call</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	reentersyscall(savedpc, uintptr(savedsp))
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	gp.m.syscall = syscall
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>func cgocallbackg1(fn, frame unsafe.Pointer, ctxt uintptr) {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	gp := getg()
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	if gp.m.needextram || extraMWaiters.Load() &gt; 0 {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		gp.m.needextram = false
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		systemstack(newextram)
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	if ctxt != 0 {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		s := append(gp.cgoCtxt, ctxt)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		<span class="comment">// Now we need to set gp.cgoCtxt = s, but we could get</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		<span class="comment">// a SIGPROF signal while manipulating the slice, and</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		<span class="comment">// the SIGPROF handler could pick up gp.cgoCtxt while</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		<span class="comment">// tracing up the stack.  We need to ensure that the</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		<span class="comment">// handler always sees a valid slice, so set the</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		<span class="comment">// values in an order such that it always does.</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		p := (*slice)(unsafe.Pointer(&amp;gp.cgoCtxt))
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		atomicstorep(unsafe.Pointer(&amp;p.array), unsafe.Pointer(&amp;s[0]))
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		p.cap = cap(s)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		p.len = len(s)
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		defer func(gp *g) {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			<span class="comment">// Decrease the length of the slice by one, safely.</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			p := (*slice)(unsafe.Pointer(&amp;gp.cgoCtxt))
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>			p.len--
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		}(gp)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	if gp.m.ncgo == 0 {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		<span class="comment">// The C call to Go came from a thread not currently running</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		<span class="comment">// any Go. In the case of -buildmode=c-archive or c-shared,</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		<span class="comment">// this call may be coming in before package initialization</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		<span class="comment">// is complete. Wait until it is.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		&lt;-main_init_done
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	<span class="comment">// Check whether the profiler needs to be turned on or off; this route to</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// run Go code does not use runtime.execute, so bypasses the check there.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	hz := sched.profilehz
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	if gp.m.profilehz != hz {
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		setThreadCPUProfiler(hz)
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">// Add entry to defer stack in case of panic.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	restore := true
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	defer unwindm(&amp;restore)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	if raceenabled {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		raceacquire(unsafe.Pointer(&amp;racecgosync))
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	<span class="comment">// Invoke callback. This function is generated by cmd/cgo and</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	<span class="comment">// will unpack the argument frame and call the Go function.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	var cb func(frame unsafe.Pointer)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	cbFV := funcval{uintptr(fn)}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	*(*unsafe.Pointer)(unsafe.Pointer(&amp;cb)) = noescape(unsafe.Pointer(&amp;cbFV))
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	cb(frame)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	if raceenabled {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		racereleasemerge(unsafe.Pointer(&amp;racecgosync))
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	<span class="comment">// Do not unwind m-&gt;g0-&gt;sched.sp.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	<span class="comment">// Our caller, cgocallback, will do that.</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	restore = false
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>func unwindm(restore *bool) {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	if *restore {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		<span class="comment">// Restore sp saved by cgocallback during</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		<span class="comment">// unwind of g&#39;s stack (see comment at top of file).</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		mp := acquirem()
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		sched := &amp;mp.g0.sched
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		sched.sp = *(*uintptr)(unsafe.Pointer(sched.sp + alignUp(sys.MinFrameSize, sys.StackAlign)))
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		<span class="comment">// Do the accounting that cgocall will not have a chance to do</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		<span class="comment">// during an unwind.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		<span class="comment">// In the case where a Go call originates from C, ncgo is 0</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		<span class="comment">// and there is no matching cgocall to end.</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		if mp.ncgo &gt; 0 {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			mp.incgo = false
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>			mp.ncgo--
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			osPreemptExtExit(mp)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		<span class="comment">// Undo the call to lockOSThread in cgocallbackg, only on the</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		<span class="comment">// panicking path. In normal return case cgocallbackg will call</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		<span class="comment">// unlockOSThread, ensuring no preemption point after the unlock.</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		<span class="comment">// Here we don&#39;t need to worry about preemption, because we&#39;re</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		<span class="comment">// panicking out of the callback and unwinding the g0 stack,</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		<span class="comment">// instead of reentering cgo (which requires the same thread).</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		unlockOSThread()
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		releasem(mp)
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span><span class="comment">// called from assembly.</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>func badcgocallback() {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	throw(&#34;misaligned stack in cgocallback&#34;)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span><span class="comment">// called from (incomplete) assembly.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>func cgounimpl() {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	throw(&#34;cgo not implemented&#34;)
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>var racecgosync uint64 <span class="comment">// represents possible synchronization in C code</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span><span class="comment">// Pointer checking for cgo code.</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span><span class="comment">// We want to detect all cases where a program that does not use</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span><span class="comment">// unsafe makes a cgo call passing a Go pointer to memory that</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span><span class="comment">// contains an unpinned Go pointer. Here a Go pointer is defined as a</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span><span class="comment">// pointer to memory allocated by the Go runtime. Programs that use</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span><span class="comment">// unsafe can evade this restriction easily, so we don&#39;t try to catch</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">// them. The cgo program will rewrite all possibly bad pointer</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span><span class="comment">// arguments to call cgoCheckPointer, where we can catch cases of a Go</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span><span class="comment">// pointer pointing to an unpinned Go pointer.</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span><span class="comment">// Complicating matters, taking the address of a slice or array</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span><span class="comment">// element permits the C program to access all elements of the slice</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span><span class="comment">// or array. In that case we will see a pointer to a single element,</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span><span class="comment">// but we need to check the entire data structure.</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span><span class="comment">// The cgoCheckPointer call takes additional arguments indicating that</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span><span class="comment">// it was called on an address expression. An additional argument of</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span><span class="comment">// true means that it only needs to check a single element. An</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span><span class="comment">// additional argument of a slice or array means that it needs to</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span><span class="comment">// check the entire slice/array, but nothing else. Otherwise, the</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span><span class="comment">// pointer could be anything, and we check the entire heap object,</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span><span class="comment">// which is conservative but safe.</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span><span class="comment">// When and if we implement a moving garbage collector,</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span><span class="comment">// cgoCheckPointer will pin the pointer for the duration of the cgo</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span><span class="comment">// call.  (This is necessary but not sufficient; the cgo program will</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span><span class="comment">// also have to change to pin Go pointers that cannot point to Go</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span><span class="comment">// pointers.)</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span><span class="comment">// cgoCheckPointer checks if the argument contains a Go pointer that</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span><span class="comment">// points to an unpinned Go pointer, and panics if it does.</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>func cgoCheckPointer(ptr any, arg any) {
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	if !goexperiment.CgoCheck2 &amp;&amp; debug.cgocheck == 0 {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		return
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	ep := efaceOf(&amp;ptr)
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	t := ep._type
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	top := true
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	if arg != nil &amp;&amp; (t.Kind_&amp;kindMask == kindPtr || t.Kind_&amp;kindMask == kindUnsafePointer) {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		p := ep.data
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		if t.Kind_&amp;kindDirectIface == 0 {
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>			p = *(*unsafe.Pointer)(p)
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		if p == nil || !cgoIsGoPointer(p) {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>			return
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		aep := efaceOf(&amp;arg)
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		switch aep._type.Kind_ &amp; kindMask {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		case kindBool:
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			if t.Kind_&amp;kindMask == kindUnsafePointer {
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>				<span class="comment">// We don&#39;t know the type of the element.</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>				break
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>			pt := (*ptrtype)(unsafe.Pointer(t))
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>			cgoCheckArg(pt.Elem, p, true, false, cgoCheckPointerFail)
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>			return
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		case kindSlice:
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			<span class="comment">// Check the slice rather than the pointer.</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			ep = aep
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>			t = ep._type
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		case kindArray:
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			<span class="comment">// Check the array rather than the pointer.</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>			<span class="comment">// Pass top as false since we have a pointer</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>			<span class="comment">// to the array.</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>			ep = aep
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>			t = ep._type
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>			top = false
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		default:
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>			throw(&#34;can&#39;t happen&#34;)
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		}
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	}
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	cgoCheckArg(t, ep.data, t.Kind_&amp;kindDirectIface == 0, top, cgoCheckPointerFail)
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>const cgoCheckPointerFail = &#34;cgo argument has Go pointer to unpinned Go pointer&#34;
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>const cgoResultFail = &#34;cgo result is unpinned Go pointer or points to unpinned Go pointer&#34;
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">// cgoCheckArg is the real work of cgoCheckPointer. The argument p</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">// is either a pointer to the value (of type t), or the value itself,</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">// depending on indir. The top parameter is whether we are at the top</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">// level, where Go pointers are allowed. Go pointers to pinned objects are</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">// allowed as long as they don&#39;t reference other unpinned pointers.</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>func cgoCheckArg(t *_type, p unsafe.Pointer, indir, top bool, msg string) {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	if t.PtrBytes == 0 || p == nil {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		<span class="comment">// If the type has no pointers there is nothing to do.</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		return
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	switch t.Kind_ &amp; kindMask {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	default:
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		throw(&#34;can&#39;t happen&#34;)
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	case kindArray:
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		at := (*arraytype)(unsafe.Pointer(t))
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		if !indir {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			if at.Len != 1 {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>				throw(&#34;can&#39;t happen&#34;)
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>			}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			cgoCheckArg(at.Elem, p, at.Elem.Kind_&amp;kindDirectIface == 0, top, msg)
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>			return
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; at.Len; i++ {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>			cgoCheckArg(at.Elem, p, true, top, msg)
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>			p = add(p, at.Elem.Size_)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		}
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	case kindChan, kindMap:
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		<span class="comment">// These types contain internal pointers that will</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		<span class="comment">// always be allocated in the Go heap. It&#39;s never OK</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		<span class="comment">// to pass them to C.</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		panic(errorString(msg))
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	case kindFunc:
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		if indir {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>			p = *(*unsafe.Pointer)(p)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		if !cgoIsGoPointer(p) {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>			return
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		}
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		panic(errorString(msg))
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	case kindInterface:
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		it := *(**_type)(p)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		if it == nil {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>			return
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		}
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		<span class="comment">// A type known at compile time is OK since it&#39;s</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		<span class="comment">// constant. A type not known at compile time will be</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		<span class="comment">// in the heap and will not be OK.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		if inheap(uintptr(unsafe.Pointer(it))) {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>			panic(errorString(msg))
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		p = *(*unsafe.Pointer)(add(p, goarch.PtrSize))
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		if !cgoIsGoPointer(p) {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>			return
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		}
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		if !top &amp;&amp; !isPinned(p) {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			panic(errorString(msg))
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		}
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		cgoCheckArg(it, p, it.Kind_&amp;kindDirectIface == 0, false, msg)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	case kindSlice:
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		st := (*slicetype)(unsafe.Pointer(t))
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		s := (*slice)(p)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		p = s.array
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		if p == nil || !cgoIsGoPointer(p) {
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>			return
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		}
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		if !top &amp;&amp; !isPinned(p) {
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			panic(errorString(msg))
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		if st.Elem.PtrBytes == 0 {
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>			return
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		}
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		for i := 0; i &lt; s.cap; i++ {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>			cgoCheckArg(st.Elem, p, true, false, msg)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>			p = add(p, st.Elem.Size_)
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		}
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	case kindString:
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		ss := (*stringStruct)(p)
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		if !cgoIsGoPointer(ss.str) {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>			return
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		}
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		if !top &amp;&amp; !isPinned(ss.str) {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			panic(errorString(msg))
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		}
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	case kindStruct:
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		st := (*structtype)(unsafe.Pointer(t))
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		if !indir {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>			if len(st.Fields) != 1 {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>				throw(&#34;can&#39;t happen&#34;)
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			cgoCheckArg(st.Fields[0].Typ, p, st.Fields[0].Typ.Kind_&amp;kindDirectIface == 0, top, msg)
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>			return
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		for _, f := range st.Fields {
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			if f.Typ.PtrBytes == 0 {
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>				continue
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>			}
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>			cgoCheckArg(f.Typ, add(p, f.Offset), true, top, msg)
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		}
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	case kindPtr, kindUnsafePointer:
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		if indir {
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>			p = *(*unsafe.Pointer)(p)
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			if p == nil {
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>				return
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>			}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		if !cgoIsGoPointer(p) {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>			return
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		if !top &amp;&amp; !isPinned(p) {
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>			panic(errorString(msg))
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		}
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		cgoCheckUnknownPointer(p, msg)
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	}
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span><span class="comment">// cgoCheckUnknownPointer is called for an arbitrary pointer into Go</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span><span class="comment">// memory. It checks whether that Go memory contains any other</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span><span class="comment">// pointer into unpinned Go memory. If it does, we panic.</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span><span class="comment">// The return values are unused but useful to see in panic tracebacks.</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>func cgoCheckUnknownPointer(p unsafe.Pointer, msg string) (base, i uintptr) {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	if inheap(uintptr(p)) {
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		b, span, _ := findObject(uintptr(p), 0, 0)
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		base = b
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		if base == 0 {
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>			return
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		if goexperiment.AllocHeaders {
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>			tp := span.typePointersOfUnchecked(base)
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>			for {
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>				var addr uintptr
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>				if tp, addr = tp.next(base + span.elemsize); addr == 0 {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>					break
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>				}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>				pp := *(*unsafe.Pointer)(unsafe.Pointer(addr))
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>				if cgoIsGoPointer(pp) &amp;&amp; !isPinned(pp) {
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>					panic(errorString(msg))
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>				}
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>			}
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		} else {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>			n := span.elemsize
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>			hbits := heapBitsForAddr(base, n)
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			for {
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>				var addr uintptr
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>				if hbits, addr = hbits.next(); addr == 0 {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>					break
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>				}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>				pp := *(*unsafe.Pointer)(unsafe.Pointer(addr))
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>				if cgoIsGoPointer(pp) &amp;&amp; !isPinned(pp) {
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>					panic(errorString(msg))
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>				}
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>			}
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		}
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		return
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	}
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	for _, datap := range activeModules() {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		if cgoInRange(p, datap.data, datap.edata) || cgoInRange(p, datap.bss, datap.ebss) {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			<span class="comment">// We have no way to know the size of the object.</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			<span class="comment">// We have to assume that it might contain a pointer.</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>			panic(errorString(msg))
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		}
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		<span class="comment">// In the text or noptr sections, we know that the</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>		<span class="comment">// pointer does not point to a Go pointer.</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	}
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	return
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>}
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span><span class="comment">// cgoIsGoPointer reports whether the pointer is a Go pointer--a</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span><span class="comment">// pointer to Go memory. We only care about Go memory that might</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span><span class="comment">// contain pointers.</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>func cgoIsGoPointer(p unsafe.Pointer) bool {
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	if p == nil {
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		return false
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	}
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	if inHeapOrStack(uintptr(p)) {
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		return true
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	}
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	for _, datap := range activeModules() {
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		if cgoInRange(p, datap.data, datap.edata) || cgoInRange(p, datap.bss, datap.ebss) {
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>			return true
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		}
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	return false
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>}
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span><span class="comment">// cgoInRange reports whether p is between start and end.</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>func cgoInRange(p unsafe.Pointer, start, end uintptr) bool {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	return start &lt;= uintptr(p) &amp;&amp; uintptr(p) &lt; end
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span><span class="comment">// cgoCheckResult is called to check the result parameter of an</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span><span class="comment">// exported Go function. It panics if the result is or contains any</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span><span class="comment">// other pointer into unpinned Go memory.</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>func cgoCheckResult(val any) {
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	if !goexperiment.CgoCheck2 &amp;&amp; debug.cgocheck == 0 {
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>		return
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	}
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	ep := efaceOf(&amp;val)
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	t := ep._type
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	cgoCheckArg(t, ep.data, t.Kind_&amp;kindDirectIface == 0, false, cgoResultFail)
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>}
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>
</pre><p><a href="cgocall.go?m=text">View as plain text</a></p>

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
