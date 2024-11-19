<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/stubs.go - Go Documentation Server</title>

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
<a href="stubs.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">stubs.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2014 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// Should be a built-in for unsafe.Pointer?</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	return unsafe.Pointer(uintptr(p) + x)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>}
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// getg returns the pointer to the current g.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// The compiler rewrites calls to this function into instructions</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// that fetch the g directly (from TLS or from the dedicated register).</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>func getg() *g
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// mcall switches from the g to the g0 stack and invokes fn(g),</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// where g is the goroutine that made the call.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// mcall saves g&#39;s current PC/SP in g-&gt;sched so that it can be restored later.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// It is up to fn to arrange for that later execution, typically by recording</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// g in a data structure, causing something to call ready(g) later.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// mcall returns to the original goroutine g later, when g has been rescheduled.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// fn must not return at all; typically it ends by calling schedule, to let the m</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// run other goroutines.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// mcall can only be called from g stacks (not g0, not gsignal).</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// This must NOT be go:noescape: if fn is a stack-allocated closure,</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// fn puts g on a run queue, and g executes before fn returns, the</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// closure will be invalidated while it is still executing.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>func mcall(fn func(*g))
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// systemstack runs fn on a system stack.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// If systemstack is called from the per-OS-thread (g0) stack, or</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// if systemstack is called from the signal handling (gsignal) stack,</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// systemstack calls fn directly and returns.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// Otherwise, systemstack is being called from the limited stack</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// of an ordinary goroutine. In this case, systemstack switches</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// to the per-OS-thread stack, calls fn, and switches back.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// It is common to use a func literal as the argument, in order</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// to share inputs and outputs with the code around the call</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// to system stack:</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//	... set up y ...</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//	systemstack(func() {</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//		x = bigcall(y)</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//	})</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//	... use x ...</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func systemstack(fn func())
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>func badsystemstack() {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	writeErrStr(&#34;fatal: systemstack called from unexpected goroutine&#34;)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// memclrNoHeapPointers clears n bytes starting at ptr.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// Usually you should use typedmemclr. memclrNoHeapPointers should be</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// used only when the caller knows that *ptr contains no heap pointers</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// because either:</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// *ptr is initialized memory and its type is pointer-free, or</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// *ptr is uninitialized memory (e.g., memory that&#39;s being reused</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// for a new allocation) and hence contains only &#34;junk&#34;.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// memclrNoHeapPointers ensures that if ptr is pointer-aligned, and n</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// is a multiple of the pointer size, then any pointer-aligned,</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// pointer-sized portion is cleared atomically. Despite the function</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// name, this is necessary because this function is the underlying</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// implementation of typedmemclr and memclrHasPointers. See the doc of</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// memmove for more details.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// The (CPU-specific) implementations of this function are in memclr_*.s.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_memclrNoHeapPointers reflect.memclrNoHeapPointers</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>func reflect_memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr) {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	memclrNoHeapPointers(ptr, n)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// memmove copies n bytes from &#34;from&#34; to &#34;to&#34;.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// memmove ensures that any pointer in &#34;from&#34; is written to &#34;to&#34; with</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// an indivisible write, so that racy reads cannot observe a</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// half-written pointer. This is necessary to prevent the garbage</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// collector from observing invalid pointers, and differs from memmove</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// in unmanaged languages. However, memmove is only required to do</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// this if &#34;from&#34; and &#34;to&#34; may contain pointers, which can only be the</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// case if &#34;from&#34;, &#34;to&#34;, and &#34;n&#34; are all be word-aligned.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// Implementations are in memmove_*.s.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>func memmove(to, from unsafe.Pointer, n uintptr)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// Outside assembly calls memmove. Make sure it has ABI wrappers.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//go:linkname memmove</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_memmove reflect.memmove</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>func reflect_memmove(to, from unsafe.Pointer, n uintptr) {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	memmove(to, from, n)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// exported value for testing</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>const hashLoad = float32(loadFactorNum) / float32(loadFactorDen)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// in internal/bytealg/equal_*.s</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func memequal(a, b unsafe.Pointer, size uintptr) bool
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// noescape hides a pointer from escape analysis.  noescape is</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// the identity function but escape analysis doesn&#39;t think the</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// output depends on the input.  noescape is inlined and currently</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// compiles down to zero instructions.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// USE CAREFULLY!</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>func noescape(p unsafe.Pointer) unsafe.Pointer {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	x := uintptr(p)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	return unsafe.Pointer(x ^ 0)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">// noEscapePtr hides a pointer from escape analysis. See noescape.</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// USE CAREFULLY!</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>func noEscapePtr[T any](p *T) *T {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	x := uintptr(unsafe.Pointer(p))
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	return (*T)(unsafe.Pointer(x ^ 0))
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// Not all cgocallback frames are actually cgocallback,</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// so not all have these arguments. Mark them uintptr so that the GC</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// does not misinterpret memory when the arguments are not present.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// cgocallback is not called from Go, only from crosscall2.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// This in turn calls cgocallbackg, which is where we&#39;ll find</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">// pointer-declared arguments.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// When fn is nil (frame is saved g), call dropm instead,</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// this is used when the C thread is exiting.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>func cgocallback(fn, frame, ctxt uintptr)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>func gogo(buf *gobuf)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func asminit()
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>func setg(gg *g)
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>func breakpoint()
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// reflectcall calls fn with arguments described by stackArgs, stackArgsSize,</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">// frameSize, and regArgs.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// Arguments passed on the stack and space for return values passed on the stack</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span><span class="comment">// must be laid out at the space pointed to by stackArgs (with total length</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// stackArgsSize) according to the ABI.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// stackRetOffset must be some value &lt;= stackArgsSize that indicates the</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// offset within stackArgs where the return value space begins.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">// frameSize is the total size of the argument frame at stackArgs and must</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// therefore be &gt;= stackArgsSize. It must include additional space for spilling</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// register arguments for stack growth and preemption.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">// TODO(mknyszek): Once we don&#39;t need the additional spill space, remove frameSize,</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">// since frameSize will be redundant with stackArgsSize.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">// Arguments passed in registers must be laid out in regArgs according to the ABI.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">// regArgs will hold any return values passed in registers after the call.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span><span class="comment">// reflectcall copies stack arguments from stackArgs to the goroutine stack, and</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span><span class="comment">// then copies back stackArgsSize-stackRetOffset bytes back to the return space</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// in stackArgs once fn has completed. It also &#34;unspills&#34; argument registers from</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// regArgs before calling fn, and spills them back into regArgs immediately</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// following the call to fn. If there are results being returned on the stack,</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// the caller should pass the argument frame type as stackArgsType so that</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// reflectcall can execute appropriate write barriers during the copy.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// reflectcall expects regArgs.ReturnIsPtr to be populated indicating which</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">// registers on the return path will contain Go pointers. It will then store</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span><span class="comment">// these pointers in regArgs.Ptrs such that they are visible to the GC.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// Package reflect passes a frame type. In package runtime, there is only</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// one call that copies results back, in callbackWrap in syscall_windows.go, and it</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// does NOT pass a frame type, meaning there are no write barriers invoked. See that</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// call site for justification.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// Package reflect accesses this symbol through a linkname.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">// Arguments passed through to reflectcall do not escape. The type is used</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">// only in a very limited callee of reflectcall, the stackArgs are copied, and</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// regArgs is only used in the reflectcall frame.</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>func reflectcall(stackArgsType *_type, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>func procyield(cycles uint32)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>type neverCallThisFunction struct{}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// goexit is the return stub at the top of every goroutine call stack.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// Each goroutine stack is constructed as if goexit called the</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// goroutine&#39;s entry point function, so that when the entry point</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// function returns, it will return to goexit, which will call goexit1</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// to perform the actual exit.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">// This function must never be called directly. Call goexit1 instead.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// gentraceback assumes that goexit terminates the stack. A direct</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// call on the stack will cause gentraceback to stop walking the stack</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// prematurely and if there is leftover state it may panic.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>func goexit(neverCallThisFunction)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// publicationBarrier performs a store/store barrier (a &#34;publication&#34;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// or &#34;export&#34; barrier). Some form of synchronization is required</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// between initializing an object and making that object accessible to</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// another processor. Without synchronization, the initialization</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// writes and the &#34;publication&#34; write may be reordered, allowing the</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// other processor to follow the pointer and observe an uninitialized</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// object. In general, higher-level synchronization should be used,</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">// such as locking or an atomic pointer write. publicationBarrier is</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">// for when those aren&#39;t an option, such as in the implementation of</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">// the memory manager.</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">// There&#39;s no corresponding barrier for the read side because the read</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">// side naturally has a data dependency order. All architectures that</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">// Go supports or seems likely to ever support automatically enforce</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span><span class="comment">// data dependency ordering.</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>func publicationBarrier()
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span><span class="comment">// getcallerpc returns the program counter (PC) of its caller&#39;s caller.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">// getcallersp returns the stack pointer (SP) of its caller&#39;s caller.</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// The implementation may be a compiler intrinsic; there is not</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">// necessarily code implementing this on every platform.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// For example:</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">//	func f(arg1, arg2, arg3 int) {</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span><span class="comment">//		pc := getcallerpc()</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span><span class="comment">//		sp := getcallersp()</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// These two lines find the PC and SP immediately following</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// the call to f (where f will return).</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span><span class="comment">// The call to getcallerpc and getcallersp must be done in the</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">// frame being asked about.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span><span class="comment">// The result of getcallersp is correct at the time of the return,</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span><span class="comment">// but it may be invalidated by any subsequent call to a function</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">// that might relocate the stack in order to grow or shrink it.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span><span class="comment">// A general rule is that the result of getcallersp should be used</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">// immediately and can only be passed to nosplit functions.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>func getcallerpc() uintptr
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>func getcallersp() uintptr <span class="comment">// implemented as an intrinsic on all platforms</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span><span class="comment">// getclosureptr returns the pointer to the current closure.</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span><span class="comment">// getclosureptr can only be used in an assignment statement</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span><span class="comment">// at the entry of a function. Moreover, go:nosplit directive</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span><span class="comment">// must be specified at the declaration of caller function,</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// so that the function prolog does not clobber the closure register.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span><span class="comment">// for example:</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">//	//go:nosplit</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">//	func f(arg1, arg2, arg3 int) {</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">//		dx := getclosureptr()</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">// The compiler rewrites calls to this function into instructions that fetch the</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">// pointer from a well-known register (DX on x86 architecture, etc.) directly.</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">// WARNING: PGO-based devirtualization cannot detect that caller of</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span><span class="comment">// getclosureptr require closure context, and thus must maintain a list of</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">// these functions, which is in</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span><span class="comment">// cmd/compile/internal/devirtualize/pgo.maybeDevirtualizeFunctionCall.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>func getclosureptr() uintptr
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>func asmcgocall(fn, arg unsafe.Pointer) int32
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>func morestack()
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>func morestack_noctxt()
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>func rt0_go()
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">// return0 is a stub used to return 0 from deferproc.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span><span class="comment">// It is called at the very end of deferproc to signal</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span><span class="comment">// the calling Go function that it should not jump</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span><span class="comment">// to deferreturn.</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">// in asm_*.s</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>func return0()
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// in asm_*.s</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// not called directly; definitions here supply type information for traceback.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// These must have the same signature (arg pointer map) as reflectcall.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>func call16(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>func call32(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>func call64(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>func call128(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>func call256(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>func call512(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>func call1024(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>func call2048(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>func call4096(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>func call8192(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>func call16384(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>func call32768(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>func call65536(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>func call131072(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>func call262144(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>func call524288(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>func call1048576(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>func call2097152(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>func call4194304(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>func call8388608(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>func call16777216(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>func call33554432(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>func call67108864(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>func call134217728(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>func call268435456(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func call536870912(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>func call1073741824(typ, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>func systemstack_switch()
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">// alignUp rounds n up to a multiple of a. a must be a power of 2.</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>func alignUp(n, a uintptr) uintptr {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	return (n + a - 1) &amp;^ (a - 1)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// alignDown rounds n down to a multiple of a. a must be a power of 2.</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>func alignDown(n, a uintptr) uintptr {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	return n &amp;^ (a - 1)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// divRoundUp returns ceil(n / a).</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>func divRoundUp(n, a uintptr) uintptr {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	<span class="comment">// a is generally a power of two. This will get inlined and</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	<span class="comment">// the compiler will optimize the division.</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	return (n + a - 1) / a
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// checkASM reports whether assembly runtime checks have passed.</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>func checkASM() bool
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>func memequal_varlen(a, b unsafe.Pointer) bool
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">// bool2int returns 0 if x is false or 1 if x is true.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>func bool2int(x bool) int {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// Avoid branches. In the SSA compiler, this compiles to</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	<span class="comment">// exactly what you would want it to.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	return int(*(*uint8)(unsafe.Pointer(&amp;x)))
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span><span class="comment">// abort crashes the runtime in situations where even throw might not</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span><span class="comment">// work. In general it should do something a debugger will recognize</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span><span class="comment">// (e.g., an INT3 on x86). A crash in abort is recognized by the</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span><span class="comment">// signal handler, which will attempt to tear down the runtime</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span><span class="comment">// immediately.</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>func abort()
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span><span class="comment">// Called from compiled code; declared for vet; do NOT call from Go.</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>func gcWriteBarrier1()
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>func gcWriteBarrier2()
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>func gcWriteBarrier3()
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>func gcWriteBarrier4()
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>func gcWriteBarrier5()
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>func gcWriteBarrier6()
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>func gcWriteBarrier7()
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>func gcWriteBarrier8()
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>func duffzero()
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>func duffcopy()
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">// Called from linker-generated .initarray; declared for go vet; do NOT call from Go.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>func addmoduledata()
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">// Injected by the signal handler for panicking signals.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span><span class="comment">// Initializes any registers that have fixed meaning at calls but</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span><span class="comment">// are scratch in bodies and calls sigpanic.</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span><span class="comment">// On many platforms it just jumps to sigpanic.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>func sigpanic0()
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span><span class="comment">// intArgRegs is used by the various register assignment</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span><span class="comment">// algorithm implementations in the runtime. These include:.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span><span class="comment">// - Finalizers (mfinal.go)</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span><span class="comment">// - Windows callbacks (syscall_windows.go)</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span><span class="comment">// Both are stripped-down versions of the algorithm since they</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span><span class="comment">// only have to deal with a subset of cases (finalizers only</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span><span class="comment">// take a pointer or interface argument, Go Windows callbacks</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span><span class="comment">// don&#39;t support floating point).</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span><span class="comment">// It should be modified with care and are generally only</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span><span class="comment">// modified when testing this package.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span><span class="comment">// It should never be set higher than its internal/abi</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span><span class="comment">// constant counterparts, because the system relies on a</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span><span class="comment">// structure that is at least large enough to hold the</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span><span class="comment">// registers the system supports.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span><span class="comment">// Protected by finlock.</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>var intArgRegs = abi.IntArgRegs
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>
</pre><p><a href="stubs.go?m=text">View as plain text</a></p>

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
