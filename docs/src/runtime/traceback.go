<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/traceback.go - Go Documentation Server</title>

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
<a href="traceback.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">traceback.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/bytealg&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// The code in this file implements stack trace walking for all architectures.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// The most important fact about a given architecture is whether it uses a link register.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// On systems with link registers, the prologue for a non-leaf function stores the</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// incoming value of LR at the bottom of the newly allocated stack frame.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// On systems without link registers (x86), the architecture pushes a return PC during</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// the call instruction, so the return PC ends up above the stack frame.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// In this file, the return PC is always called LR, no matter how it was found.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>const usesLR = sys.MinFrameSize &gt; 0
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>const (
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// tracebackInnerFrames is the number of innermost frames to print in a</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// stack trace. The total maximum frames is tracebackInnerFrames +</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// tracebackOuterFrames.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	tracebackInnerFrames = 50
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// tracebackOuterFrames is the number of outermost frames to print in a</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// stack trace.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	tracebackOuterFrames = 50
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// unwindFlags control the behavior of various unwinders.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>type unwindFlags uint8
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>const (
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// unwindPrintErrors indicates that if unwinding encounters an error, it</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// should print a message and stop without throwing. This is used for things</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// like stack printing, where it&#39;s better to get incomplete information than</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// to crash. This is also used in situations where everything may not be</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// stopped nicely and the stack walk may not be able to complete, such as</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// during profiling signals or during a crash.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// If neither unwindPrintErrors or unwindSilentErrors are set, unwinding</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// performs extra consistency checks and throws on any error.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// Note that there are a small number of fatal situations that will throw</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// regardless of unwindPrintErrors or unwindSilentErrors.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	unwindPrintErrors unwindFlags = 1 &lt;&lt; iota
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// unwindSilentErrors silently ignores errors during unwinding.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	unwindSilentErrors
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// unwindTrap indicates that the initial PC and SP are from a trap, not a</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// return PC from a call.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// The unwindTrap flag is updated during unwinding. If set, frame.pc is the</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// address of a faulting instruction instead of the return address of a</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// call. It also means the liveness at pc may not be known.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// TODO: Distinguish frame.continpc, which is really the stack map PC, from</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// the actual continuation PC, which is computed differently depending on</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// this flag and a few other things.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	unwindTrap
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// unwindJumpStack indicates that, if the traceback is on a system stack, it</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// should resume tracing at the user stack when the system stack is</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// exhausted.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	unwindJumpStack
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// An unwinder iterates the physical stack frames of a Go sack.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// Typical use of an unwinder looks like:</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//	var u unwinder</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//	for u.init(gp, 0); u.valid(); u.next() {</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//		// ... use frame info in u ...</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// Implementation note: This is carefully structured to be pointer-free because</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// tracebacks happen in places that disallow write barriers (e.g., signals).</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// Even if this is stack-allocated, its pointer-receiver methods don&#39;t know that</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// their receiver is on the stack, so they still emit write barriers. Here we</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// address that by carefully avoiding any pointers in this type. Another</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// approach would be to split this into a mutable part that&#39;s passed by pointer</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// but contains no pointers itself and an immutable part that&#39;s passed and</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// returned by value and can contain pointers. We could potentially hide that</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// we&#39;re doing that in trivial methods that are inlined into the caller that has</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// the stack allocation, but that&#39;s fragile.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>type unwinder struct {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// frame is the current physical stack frame, or all 0s if</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// there is no frame.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	frame stkframe
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// g is the G who&#39;s stack is being unwound. If the</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// unwindJumpStack flag is set and the unwinder jumps stacks,</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// this will be different from the initial G.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	g guintptr
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// cgoCtxt is the index into g.cgoCtxt of the next frame on the cgo stack.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// The cgo stack is unwound in tandem with the Go stack as we find marker frames.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	cgoCtxt int
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">// calleeFuncID is the function ID of the caller of the current</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// frame.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	calleeFuncID abi.FuncID
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// flags are the flags to this unwind. Some of these are updated as we</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// unwind (see the flags documentation).</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	flags unwindFlags
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// init initializes u to start unwinding gp&#39;s stack and positions the</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// iterator on gp&#39;s innermost frame. gp must not be the current G.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// A single unwinder can be reused for multiple unwinds.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>func (u *unwinder) init(gp *g, flags unwindFlags) {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// Implementation note: This starts the iterator on the first frame and we</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// provide a &#34;valid&#34; method. Alternatively, this could start in a &#34;before</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// the first frame&#34; state and &#34;next&#34; could return whether it was able to</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// move to the next frame, but that&#39;s both more awkward to use in a &#34;for&#34;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// loop and is harder to implement because we have to do things differently</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// for the first frame.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	u.initAt(^uintptr(0), ^uintptr(0), ^uintptr(0), gp, flags)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func (u *unwinder) initAt(pc0, sp0, lr0 uintptr, gp *g, flags unwindFlags) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t call this &#34;g&#34;; it&#39;s too easy get &#34;g&#34; and &#34;gp&#34; confused.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if ourg := getg(); ourg == gp &amp;&amp; ourg == ourg.m.curg {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// The starting sp has been passed in as a uintptr, and the caller may</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// have other uintptr-typed stack references as well.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// If during one of the calls that got us here or during one of the</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		<span class="comment">// callbacks below the stack must be grown, all these uintptr references</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		<span class="comment">// to the stack will not be updated, and traceback will continue</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		<span class="comment">// to inspect the old stack memory, which may no longer be valid.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		<span class="comment">// Even if all the variables were updated correctly, it is not clear that</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		<span class="comment">// we want to expose a traceback that begins on one stack and ends</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		<span class="comment">// on another stack. That could confuse callers quite a bit.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		<span class="comment">// Instead, we require that initAt and any other function that</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">// accepts an sp for the current goroutine (typically obtained by</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		<span class="comment">// calling getcallersp) must not run on that goroutine&#39;s stack but</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		<span class="comment">// instead on the g0 stack.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		throw(&#34;cannot trace user goroutine on its own stack&#34;)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	if pc0 == ^uintptr(0) &amp;&amp; sp0 == ^uintptr(0) { <span class="comment">// Signal to fetch saved values from gp.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		if gp.syscallsp != 0 {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			pc0 = gp.syscallpc
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			sp0 = gp.syscallsp
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			if usesLR {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>				lr0 = 0
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		} else {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			pc0 = gp.sched.pc
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			sp0 = gp.sched.sp
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			if usesLR {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>				lr0 = gp.sched.lr
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	var frame stkframe
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	frame.pc = pc0
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	frame.sp = sp0
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	if usesLR {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		frame.lr = lr0
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// If the PC is zero, it&#39;s likely a nil function call.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// Start in the caller&#39;s frame.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	if frame.pc == 0 {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		if usesLR {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			frame.pc = *(*uintptr)(unsafe.Pointer(frame.sp))
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			frame.lr = 0
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		} else {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			frame.pc = *(*uintptr)(unsafe.Pointer(frame.sp))
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			frame.sp += goarch.PtrSize
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// runtime/internal/atomic functions call into kernel helpers on</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// arm &lt; 7. See runtime/internal/atomic/sys_linux_arm.s.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">// Start in the caller&#39;s frame.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	if GOARCH == &#34;arm&#34; &amp;&amp; goarm &lt; 7 &amp;&amp; GOOS == &#34;linux&#34; &amp;&amp; frame.pc&amp;0xffff0000 == 0xffff0000 {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		<span class="comment">// Note that the calls are simple BL without pushing the return</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		<span class="comment">// address, so we use LR directly.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		<span class="comment">// The kernel helpers are frameless leaf functions, so SP and</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		<span class="comment">// LR are not touched.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		frame.pc = frame.lr
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		frame.lr = 0
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	f := findfunc(frame.pc)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	if !f.valid() {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		if flags&amp;unwindSilentErrors == 0 {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>			print(&#34;runtime: g &#34;, gp.goid, &#34; gp=&#34;, gp, &#34;: unknown pc &#34;, hex(frame.pc), &#34;\n&#34;)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			tracebackHexdump(gp.stack, &amp;frame, 0)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		if flags&amp;(unwindPrintErrors|unwindSilentErrors) == 0 {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			throw(&#34;unknown pc&#34;)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		*u = unwinder{}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		return
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	frame.fn = f
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// Populate the unwinder.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	*u = unwinder{
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		frame:        frame,
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		g:            gp.guintptr(),
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		cgoCtxt:      len(gp.cgoCtxt) - 1,
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		calleeFuncID: abi.FuncIDNormal,
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		flags:        flags,
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	isSyscall := frame.pc == pc0 &amp;&amp; frame.sp == sp0 &amp;&amp; pc0 == gp.syscallpc &amp;&amp; sp0 == gp.syscallsp
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	u.resolveInternal(true, isSyscall)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>func (u *unwinder) valid() bool {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	return u.frame.pc != 0
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// resolveInternal fills in u.frame based on u.frame.fn, pc, and sp.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// innermost indicates that this is the first resolve on this stack. If</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">// innermost is set, isSyscall indicates that the PC/SP was retrieved from</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">// gp.syscall*; this is otherwise ignored.</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">// On entry, u.frame contains:</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">//   - fn is the running function.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">//   - pc is the PC in the running function.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">//   - sp is the stack pointer at that program counter.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span><span class="comment">//   - For the innermost frame on LR machines, lr is the program counter that called fn.</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span><span class="comment">// On return, u.frame contains:</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span><span class="comment">//   - fp is the stack pointer of the caller.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">//   - lr is the program counter that called fn.</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">//   - varp, argp, and continpc are populated for the current frame.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// If fn is a stack-jumping function, resolveInternal can change the entire</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// frame state to follow that stack jump.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">// This is internal to unwinder.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func (u *unwinder) resolveInternal(innermost, isSyscall bool) {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	frame := &amp;u.frame
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	gp := u.g.ptr()
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	f := frame.fn
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	if f.pcsp == 0 {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		<span class="comment">// No frame information, must be external function, like race support.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		<span class="comment">// See golang.org/issue/13568.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		u.finishInternal()
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		return
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// Compute function info flags.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	flag := f.flag
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	if f.funcID == abi.FuncID_cgocallback {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		<span class="comment">// cgocallback does write SP to switch from the g0 to the curg stack,</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		<span class="comment">// but it carefully arranges that during the transition BOTH stacks</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		<span class="comment">// have cgocallback frame valid for unwinding through.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		<span class="comment">// So we don&#39;t need to exclude it with the other SP-writing functions.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		flag &amp;^= abi.FuncFlagSPWrite
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	if isSyscall {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		<span class="comment">// Some Syscall functions write to SP, but they do so only after</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		<span class="comment">// saving the entry PC/SP using entersyscall.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		<span class="comment">// Since we are using the entry PC/SP, the later SP write doesn&#39;t matter.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		flag &amp;^= abi.FuncFlagSPWrite
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	<span class="comment">// Found an actual function.</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	<span class="comment">// Derive frame pointer.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	if frame.fp == 0 {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		<span class="comment">// Jump over system stack transitions. If we&#39;re on g0 and there&#39;s a user</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		<span class="comment">// goroutine, try to jump. Otherwise this is a regular call.</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		<span class="comment">// We also defensively check that this won&#39;t switch M&#39;s on us,</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		<span class="comment">// which could happen at critical points in the scheduler.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		<span class="comment">// This ensures gp.m doesn&#39;t change from a stack jump.</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		if u.flags&amp;unwindJumpStack != 0 &amp;&amp; gp == gp.m.g0 &amp;&amp; gp.m.curg != nil &amp;&amp; gp.m.curg.m == gp.m {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			switch f.funcID {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			case abi.FuncID_morestack:
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>				<span class="comment">// morestack does not return normally -- newstack()</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>				<span class="comment">// gogo&#39;s to curg.sched. Match that.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>				<span class="comment">// This keeps morestack() from showing up in the backtrace,</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>				<span class="comment">// but that makes some sense since it&#39;ll never be returned</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>				<span class="comment">// to.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>				gp = gp.m.curg
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>				u.g.set(gp)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>				frame.pc = gp.sched.pc
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>				frame.fn = findfunc(frame.pc)
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				f = frame.fn
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>				flag = f.flag
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>				frame.lr = gp.sched.lr
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>				frame.sp = gp.sched.sp
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>				u.cgoCtxt = len(gp.cgoCtxt) - 1
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			case abi.FuncID_systemstack:
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>				<span class="comment">// systemstack returns normally, so just follow the</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>				<span class="comment">// stack transition.</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>				if usesLR &amp;&amp; funcspdelta(f, frame.pc) == 0 {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>					<span class="comment">// We&#39;re at the function prologue and the stack</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>					<span class="comment">// switch hasn&#39;t happened, or epilogue where we&#39;re</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>					<span class="comment">// about to return. Just unwind normally.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>					<span class="comment">// Do this only on LR machines because on x86</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>					<span class="comment">// systemstack doesn&#39;t have an SP delta (the CALL</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>					<span class="comment">// instruction opens the frame), therefore no way</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>					<span class="comment">// to check.</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>					flag &amp;^= abi.FuncFlagSPWrite
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>					break
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>				}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>				gp = gp.m.curg
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>				u.g.set(gp)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>				frame.sp = gp.sched.sp
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>				u.cgoCtxt = len(gp.cgoCtxt) - 1
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>				flag &amp;^= abi.FuncFlagSPWrite
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>			}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		frame.fp = frame.sp + uintptr(funcspdelta(f, frame.pc))
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		if !usesLR {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			<span class="comment">// On x86, call instruction pushes return PC before entering new function.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			frame.fp += goarch.PtrSize
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">// Derive link register.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	if flag&amp;abi.FuncFlagTopFrame != 0 {
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		<span class="comment">// This function marks the top of the stack. Stop the traceback.</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		frame.lr = 0
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	} else if flag&amp;abi.FuncFlagSPWrite != 0 &amp;&amp; (!innermost || u.flags&amp;(unwindPrintErrors|unwindSilentErrors) != 0) {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		<span class="comment">// The function we are in does a write to SP that we don&#39;t know</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		<span class="comment">// how to encode in the spdelta table. Examples include context</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		<span class="comment">// switch routines like runtime.gogo but also any code that switches</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		<span class="comment">// to the g0 stack to run host C code.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		<span class="comment">// We can&#39;t reliably unwind the SP (we might not even be on</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		<span class="comment">// the stack we think we are), so stop the traceback here.</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		<span class="comment">// The one exception (encoded in the complex condition above) is that</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		<span class="comment">// we assume if we&#39;re doing a precise traceback, and this is the</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		<span class="comment">// innermost frame, that the SPWRITE function voluntarily preempted itself on entry</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		<span class="comment">// during the stack growth check. In that case, the function has</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		<span class="comment">// not yet had a chance to do any writes to SP and is safe to unwind.</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		<span class="comment">// isAsyncSafePoint does not allow assembly functions to be async preempted,</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		<span class="comment">// and preemptPark double-checks that SPWRITE functions are not async preempted.</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		<span class="comment">// So for GC stack traversal, we can safely ignore SPWRITE for the innermost frame,</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		<span class="comment">// but farther up the stack we&#39;d better not find any.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		<span class="comment">// This is somewhat imprecise because we&#39;re just guessing that we&#39;re in the stack</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		<span class="comment">// growth check. It would be better if SPWRITE were encoded in the spdelta</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		<span class="comment">// table so we would know for sure that we were still in safe code.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		<span class="comment">// uSE uPE inn | action</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		<span class="comment">//  T   _   _  | frame.lr = 0</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		<span class="comment">//  F   T   _  | frame.lr = 0</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		<span class="comment">//  F   F   F  | print; panic</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		<span class="comment">//  F   F   T  | ignore SPWrite</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		if u.flags&amp;(unwindPrintErrors|unwindSilentErrors) == 0 &amp;&amp; !innermost {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			println(&#34;traceback: unexpected SPWRITE function&#34;, funcname(f))
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			throw(&#34;traceback&#34;)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		frame.lr = 0
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	} else {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		var lrPtr uintptr
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		if usesLR {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			if innermost &amp;&amp; frame.sp &lt; frame.fp || frame.lr == 0 {
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>				lrPtr = frame.sp
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>				frame.lr = *(*uintptr)(unsafe.Pointer(lrPtr))
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		} else {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			if frame.lr == 0 {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>				lrPtr = frame.fp - goarch.PtrSize
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>				frame.lr = *(*uintptr)(unsafe.Pointer(lrPtr))
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	frame.varp = frame.fp
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	if !usesLR {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		<span class="comment">// On x86, call instruction pushes return PC before entering new function.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		frame.varp -= goarch.PtrSize
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	<span class="comment">// For architectures with frame pointers, if there&#39;s</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	<span class="comment">// a frame, then there&#39;s a saved frame pointer here.</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	<span class="comment">// NOTE: This code is not as general as it looks.</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	<span class="comment">// On x86, the ABI is to save the frame pointer word at the</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	<span class="comment">// top of the stack frame, so we have to back down over it.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	<span class="comment">// On arm64, the frame pointer should be at the bottom of</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	<span class="comment">// the stack (with R29 (aka FP) = RSP), in which case we would</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	<span class="comment">// not want to do the subtraction here. But we started out without</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	<span class="comment">// any frame pointer, and when we wanted to add it, we didn&#39;t</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	<span class="comment">// want to break all the assembly doing direct writes to 8(RSP)</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	<span class="comment">// to set the first parameter to a called function.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	<span class="comment">// So we decided to write the FP link *below* the stack pointer</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	<span class="comment">// (with R29 = RSP - 8 in Go functions).</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	<span class="comment">// This is technically ABI-compatible but not standard.</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	<span class="comment">// And it happens to end up mimicking the x86 layout.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	<span class="comment">// Other architectures may make different decisions.</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	if frame.varp &gt; frame.sp &amp;&amp; framepointer_enabled {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		frame.varp -= goarch.PtrSize
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	frame.argp = frame.fp + sys.MinFrameSize
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	<span class="comment">// Determine frame&#39;s &#39;continuation PC&#39;, where it can continue.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	<span class="comment">// Normally this is the return address on the stack, but if sigpanic</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	<span class="comment">// is immediately below this function on the stack, then the frame</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	<span class="comment">// stopped executing due to a trap, and frame.pc is probably not</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	<span class="comment">// a safe point for looking up liveness information. In this panicking case,</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	<span class="comment">// the function either doesn&#39;t return at all (if it has no defers or if the</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	<span class="comment">// defers do not recover) or it returns from one of the calls to</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	<span class="comment">// deferproc a second time (if the corresponding deferred func recovers).</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	<span class="comment">// In the latter case, use a deferreturn call site as the continuation pc.</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	frame.continpc = frame.pc
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	if u.calleeFuncID == abi.FuncID_sigpanic {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		if frame.fn.deferreturn != 0 {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			frame.continpc = frame.fn.entry() + uintptr(frame.fn.deferreturn) + 1
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			<span class="comment">// Note: this may perhaps keep return variables alive longer than</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>			<span class="comment">// strictly necessary, as we are using &#34;function has a defer statement&#34;</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			<span class="comment">// as a proxy for &#34;function actually deferred something&#34;. It seems</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>			<span class="comment">// to be a minor drawback. (We used to actually look through the</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			<span class="comment">// gp._defer for a defer corresponding to this function, but that</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>			<span class="comment">// is hard to do with defer records on the stack during a stack copy.)</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			<span class="comment">// Note: the +1 is to offset the -1 that</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>			<span class="comment">// stack.go:getStackMap does to back up a return</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			<span class="comment">// address make sure the pc is in the CALL instruction.</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		} else {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			frame.continpc = 0
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	}
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>func (u *unwinder) next() {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	frame := &amp;u.frame
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	f := frame.fn
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	gp := u.g.ptr()
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	<span class="comment">// Do not unwind past the bottom of the stack.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	if frame.lr == 0 {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		u.finishInternal()
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		return
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	flr := findfunc(frame.lr)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	if !flr.valid() {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		<span class="comment">// This happens if you get a profiling interrupt at just the wrong time.</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		<span class="comment">// In that context it is okay to stop early.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		<span class="comment">// But if no error flags are set, we&#39;re doing a garbage collection and must</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		<span class="comment">// get everything, so crash loudly.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		fail := u.flags&amp;(unwindPrintErrors|unwindSilentErrors) == 0
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		doPrint := u.flags&amp;unwindSilentErrors == 0
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		if doPrint &amp;&amp; gp.m.incgo &amp;&amp; f.funcID == abi.FuncID_sigpanic {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			<span class="comment">// We can inject sigpanic</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>			<span class="comment">// calls directly into C code,</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			<span class="comment">// in which case we&#39;ll see a C</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			<span class="comment">// return PC. Don&#39;t complain.</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>			doPrint = false
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		if fail || doPrint {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>			print(&#34;runtime: g &#34;, gp.goid, &#34;: unexpected return pc for &#34;, funcname(f), &#34; called from &#34;, hex(frame.lr), &#34;\n&#34;)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			tracebackHexdump(gp.stack, frame, 0)
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		if fail {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>			throw(&#34;unknown caller pc&#34;)
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		frame.lr = 0
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		u.finishInternal()
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		return
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	if frame.pc == frame.lr &amp;&amp; frame.sp == frame.fp {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		<span class="comment">// If the next frame is identical to the current frame, we cannot make progress.</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		print(&#34;runtime: traceback stuck. pc=&#34;, hex(frame.pc), &#34; sp=&#34;, hex(frame.sp), &#34;\n&#34;)
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		tracebackHexdump(gp.stack, frame, frame.sp)
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		throw(&#34;traceback stuck&#34;)
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	injectedCall := f.funcID == abi.FuncID_sigpanic || f.funcID == abi.FuncID_asyncPreempt || f.funcID == abi.FuncID_debugCallV2
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	if injectedCall {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		u.flags |= unwindTrap
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	} else {
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		u.flags &amp;^= unwindTrap
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	}
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	<span class="comment">// Unwind to next frame.</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	u.calleeFuncID = f.funcID
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	frame.fn = flr
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	frame.pc = frame.lr
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	frame.lr = 0
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	frame.sp = frame.fp
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	frame.fp = 0
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	<span class="comment">// On link register architectures, sighandler saves the LR on stack</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	<span class="comment">// before faking a call.</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	if usesLR &amp;&amp; injectedCall {
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		x := *(*uintptr)(unsafe.Pointer(frame.sp))
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		frame.sp += alignUp(sys.MinFrameSize, sys.StackAlign)
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		f = findfunc(frame.pc)
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		frame.fn = f
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		if !f.valid() {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			frame.pc = x
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		} else if funcspdelta(f, frame.pc) == 0 {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			frame.lr = x
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	u.resolveInternal(false, false)
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span><span class="comment">// finishInternal is an unwinder-internal helper called after the stack has been</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span><span class="comment">// exhausted. It sets the unwinder to an invalid state and checks that it</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span><span class="comment">// successfully unwound the entire stack.</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>func (u *unwinder) finishInternal() {
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	u.frame.pc = 0
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	<span class="comment">// Note that panic != nil is okay here: there can be leftover panics,</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	<span class="comment">// because the defers on the panic stack do not nest in frame order as</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	<span class="comment">// they do on the defer stack. If you have:</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	<span class="comment">//	frame 1 defers d1</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	<span class="comment">//	frame 2 defers d2</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	<span class="comment">//	frame 3 defers d3</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	<span class="comment">//	frame 4 panics</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	<span class="comment">//	frame 4&#39;s panic starts running defers</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	<span class="comment">//	frame 5, running d3, defers d4</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	<span class="comment">//	frame 5 panics</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	<span class="comment">//	frame 5&#39;s panic starts running defers</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	<span class="comment">//	frame 6, running d4, garbage collects</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	<span class="comment">//	frame 6, running d2, garbage collects</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	<span class="comment">// During the execution of d4, the panic stack is d4 -&gt; d3, which</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	<span class="comment">// is nested properly, and we&#39;ll treat frame 3 as resumable, because we</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	<span class="comment">// can find d3. (And in fact frame 3 is resumable. If d4 recovers</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	<span class="comment">// and frame 5 continues running, d3, d3 can recover and we&#39;ll</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	<span class="comment">// resume execution in (returning from) frame 3.)</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	<span class="comment">// During the execution of d2, however, the panic stack is d2 -&gt; d3,</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	<span class="comment">// which is inverted. The scan will match d2 to frame 2 but having</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	<span class="comment">// d2 on the stack until then means it will not match d3 to frame 3.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	<span class="comment">// This is okay: if we&#39;re running d2, then all the defers after d2 have</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	<span class="comment">// completed and their corresponding frames are dead. Not finding d3</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	<span class="comment">// for frame 3 means we&#39;ll set frame 3&#39;s continpc == 0, which is correct</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	<span class="comment">// (frame 3 is dead). At the end of the walk the panic stack can thus</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	<span class="comment">// contain defers (d3 in this case) for dead frames. The inversion here</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	<span class="comment">// always indicates a dead frame, and the effect of the inversion on the</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	<span class="comment">// scan is to hide those dead frames, so the scan is still okay:</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	<span class="comment">// what&#39;s left on the panic stack are exactly (and only) the dead frames.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	<span class="comment">// We require callback != nil here because only when callback != nil</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	<span class="comment">// do we know that gentraceback is being called in a &#34;must be correct&#34;</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	<span class="comment">// context as opposed to a &#34;best effort&#34; context. The tracebacks with</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	<span class="comment">// callbacks only happen when everything is stopped nicely.</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	<span class="comment">// At other times, such as when gathering a stack for a profiling signal</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	<span class="comment">// or when printing a traceback during a crash, everything may not be</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	<span class="comment">// stopped nicely, and the stack walk may not be able to complete.</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	gp := u.g.ptr()
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	if u.flags&amp;(unwindPrintErrors|unwindSilentErrors) == 0 &amp;&amp; u.frame.sp != gp.stktopsp {
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		print(&#34;runtime: g&#34;, gp.goid, &#34;: frame.sp=&#34;, hex(u.frame.sp), &#34; top=&#34;, hex(gp.stktopsp), &#34;\n&#34;)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		print(&#34;\tstack=[&#34;, hex(gp.stack.lo), &#34;-&#34;, hex(gp.stack.hi), &#34;\n&#34;)
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		throw(&#34;traceback did not unwind completely&#34;)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span><span class="comment">// symPC returns the PC that should be used for symbolizing the current frame.</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span><span class="comment">// Specifically, this is the PC of the last instruction executed in this frame.</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span><span class="comment">// If this frame did a normal call, then frame.pc is a return PC, so this will</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span><span class="comment">// return frame.pc-1, which points into the CALL instruction. If the frame was</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span><span class="comment">// interrupted by a signal (e.g., profiler, segv, etc) then frame.pc is for the</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span><span class="comment">// trapped instruction, so this returns frame.pc. See issue #34123. Finally,</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span><span class="comment">// frame.pc can be at function entry when the frame is initialized without</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span><span class="comment">// actually running code, like in runtime.mstart, in which case this returns</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span><span class="comment">// frame.pc because that&#39;s the best we can do.</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>func (u *unwinder) symPC() uintptr {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	if u.flags&amp;unwindTrap == 0 &amp;&amp; u.frame.pc &gt; u.frame.fn.entry() {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		<span class="comment">// Regular call.</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		return u.frame.pc - 1
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	}
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	<span class="comment">// Trapping instruction or we&#39;re at the function entry point.</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	return u.frame.pc
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>}
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span><span class="comment">// cgoCallers populates pcBuf with the cgo callers of the current frame using</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span><span class="comment">// the registered cgo unwinder. It returns the number of PCs written to pcBuf.</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span><span class="comment">// If the current frame is not a cgo frame or if there&#39;s no registered cgo</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span><span class="comment">// unwinder, it returns 0.</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>func (u *unwinder) cgoCallers(pcBuf []uintptr) int {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	if cgoTraceback == nil || u.frame.fn.funcID != abi.FuncID_cgocallback || u.cgoCtxt &lt; 0 {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		<span class="comment">// We don&#39;t have a cgo unwinder (typical case), or we do but we&#39;re not</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		<span class="comment">// in a cgo frame or we&#39;re out of cgo context.</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		return 0
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	ctxt := u.g.ptr().cgoCtxt[u.cgoCtxt]
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	u.cgoCtxt--
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	cgoContextPCs(ctxt, pcBuf)
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	for i, pc := range pcBuf {
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		if pc == 0 {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>			return i
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	return len(pcBuf)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span><span class="comment">// tracebackPCs populates pcBuf with the return addresses for each frame from u</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span><span class="comment">// and returns the number of PCs written to pcBuf. The returned PCs correspond</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span><span class="comment">// to &#34;logical frames&#34; rather than &#34;physical frames&#34;; that is if A is inlined</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span><span class="comment">// into B, this will still return a PCs for both A and B. This also includes PCs</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span><span class="comment">// generated by the cgo unwinder, if one is registered.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span><span class="comment">// If skip != 0, this skips this many logical frames.</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span><span class="comment">// Callers should set the unwindSilentErrors flag on u.</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>func tracebackPCs(u *unwinder, skip int, pcBuf []uintptr) int {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	var cgoBuf [32]uintptr
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	n := 0
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	for ; n &lt; len(pcBuf) &amp;&amp; u.valid(); u.next() {
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		f := u.frame.fn
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		cgoN := u.cgoCallers(cgoBuf[:])
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		<span class="comment">// TODO: Why does &amp;u.cache cause u to escape? (Same in traceback2)</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		for iu, uf := newInlineUnwinder(f, u.symPC()); n &lt; len(pcBuf) &amp;&amp; uf.valid(); uf = iu.next(uf) {
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			sf := iu.srcFunc(uf)
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>			if sf.funcID == abi.FuncIDWrapper &amp;&amp; elideWrapperCalling(u.calleeFuncID) {
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>				<span class="comment">// ignore wrappers</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>			} else if skip &gt; 0 {
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>				skip--
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>			} else {
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>				<span class="comment">// Callers expect the pc buffer to contain return addresses</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>				<span class="comment">// and do the -1 themselves, so we add 1 to the call PC to</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>				<span class="comment">// create a return PC.</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>				pcBuf[n] = uf.pc + 1
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>				n++
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>			}
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			u.calleeFuncID = sf.funcID
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>		<span class="comment">// Add cgo frames (if we&#39;re done skipping over the requested number of</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		<span class="comment">// Go frames).</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		if skip == 0 {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>			n += copy(pcBuf[n:], cgoBuf[:cgoN])
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	return n
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>}
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span><span class="comment">// printArgs prints function arguments in traceback.</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>func printArgs(f funcInfo, argp unsafe.Pointer, pc uintptr) {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	<span class="comment">// The &#34;instruction&#34; of argument printing is encoded in _FUNCDATA_ArgInfo.</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	<span class="comment">// See cmd/compile/internal/ssagen.emitArgInfo for the description of the</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	<span class="comment">// encoding.</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	<span class="comment">// These constants need to be in sync with the compiler.</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	const (
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		_endSeq         = 0xff
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		_startAgg       = 0xfe
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		_endAgg         = 0xfd
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		_dotdotdot      = 0xfc
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		_offsetTooLarge = 0xfb
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	)
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	const (
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		limit    = 10                       <span class="comment">// print no more than 10 args/components</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		maxDepth = 5                        <span class="comment">// no more than 5 layers of nesting</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		maxLen   = (maxDepth*3+2)*limit + 1 <span class="comment">// max length of _FUNCDATA_ArgInfo (see the compiler side for reasoning)</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	)
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	p := (*[maxLen]uint8)(funcdata(f, abi.FUNCDATA_ArgInfo))
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	if p == nil {
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		return
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	}
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	liveInfo := funcdata(f, abi.FUNCDATA_ArgLiveInfo)
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	liveIdx := pcdatavalue(f, abi.PCDATA_ArgLiveIndex, pc)
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	startOffset := uint8(0xff) <span class="comment">// smallest offset that needs liveness info (slots with a lower offset is always live)</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	if liveInfo != nil {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		startOffset = *(*uint8)(liveInfo)
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	}
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	isLive := func(off, slotIdx uint8) bool {
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		if liveInfo == nil || liveIdx &lt;= 0 {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>			return true <span class="comment">// no liveness info, always live</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		if off &lt; startOffset {
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>			return true
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		}
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		bits := *(*uint8)(add(liveInfo, uintptr(liveIdx)+uintptr(slotIdx/8)))
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		return bits&amp;(1&lt;&lt;(slotIdx%8)) != 0
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	}
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	print1 := func(off, sz, slotIdx uint8) {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>		x := readUnaligned64(add(argp, uintptr(off)))
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		<span class="comment">// mask out irrelevant bits</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		if sz &lt; 8 {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			shift := 64 - sz*8
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			if goarch.BigEndian {
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>				x = x &gt;&gt; shift
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>			} else {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>				x = x &lt;&lt; shift &gt;&gt; shift
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>			}
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		}
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		print(hex(x))
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		if !isLive(off, slotIdx) {
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>			print(&#34;?&#34;)
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>		}
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	}
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	start := true
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	printcomma := func() {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		if !start {
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>			print(&#34;, &#34;)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	}
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	pi := 0
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	slotIdx := uint8(0) <span class="comment">// register arg spill slot index</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>printloop:
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	for {
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		o := p[pi]
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		pi++
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		switch o {
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		case _endSeq:
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>			break printloop
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		case _startAgg:
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>			printcomma()
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>			print(&#34;{&#34;)
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>			start = true
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>			continue
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>		case _endAgg:
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>			print(&#34;}&#34;)
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>		case _dotdotdot:
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>			printcomma()
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>			print(&#34;...&#34;)
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>		case _offsetTooLarge:
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			printcomma()
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>			print(&#34;_&#34;)
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		default:
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>			printcomma()
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>			sz := p[pi]
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			pi++
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>			print1(o, sz, slotIdx)
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>			if o &gt;= startOffset {
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>				slotIdx++
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>			}
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		}
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		start = false
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	}
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>}
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span><span class="comment">// funcNamePiecesForPrint returns the function name for printing to the user.</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span><span class="comment">// It returns three pieces so it doesn&#39;t need an allocation for string</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span><span class="comment">// concatenation.</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>func funcNamePiecesForPrint(name string) (string, string, string) {
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	<span class="comment">// Replace the shape name in generic function with &#34;...&#34;.</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	i := bytealg.IndexByteString(name, &#39;[&#39;)
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		return name, &#34;&#34;, &#34;&#34;
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	j := len(name) - 1
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	for name[j] != &#39;]&#39; {
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		j--
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	if j &lt;= i {
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		return name, &#34;&#34;, &#34;&#34;
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	}
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	return name[:i], &#34;[...]&#34;, name[j+1:]
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>}
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span><span class="comment">// funcNameForPrint returns the function name for printing to the user.</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>func funcNameForPrint(name string) string {
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	a, b, c := funcNamePiecesForPrint(name)
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	return a + b + c
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span><span class="comment">// printFuncName prints a function name. name is the function name in</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span><span class="comment">// the binary&#39;s func data table.</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>func printFuncName(name string) {
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	if name == &#34;runtime.gopanic&#34; {
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>		print(&#34;panic&#34;)
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>		return
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	}
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	a, b, c := funcNamePiecesForPrint(name)
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	print(a, b, c)
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>}
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>func printcreatedby(gp *g) {
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	<span class="comment">// Show what created goroutine, except main goroutine (goid 1).</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	pc := gp.gopc
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	f := findfunc(pc)
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	if f.valid() &amp;&amp; showframe(f.srcFunc(), gp, false, abi.FuncIDNormal) &amp;&amp; gp.goid != 1 {
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>		printcreatedby1(f, pc, gp.parentGoid)
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	}
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>}
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>func printcreatedby1(f funcInfo, pc uintptr, goid uint64) {
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	print(&#34;created by &#34;)
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>	printFuncName(funcname(f))
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	if goid != 0 {
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		print(&#34; in goroutine &#34;, goid)
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	}
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	print(&#34;\n&#34;)
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	tracepc := pc <span class="comment">// back up to CALL instruction for funcline.</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	if pc &gt; f.entry() {
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		tracepc -= sys.PCQuantum
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	}
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	file, line := funcline(f, tracepc)
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	print(&#34;\t&#34;, file, &#34;:&#34;, line)
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	if pc &gt; f.entry() {
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		print(&#34; +&#34;, hex(pc-f.entry()))
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	}
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	print(&#34;\n&#34;)
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>}
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>func traceback(pc, sp, lr uintptr, gp *g) {
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	traceback1(pc, sp, lr, gp, 0)
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>}
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span><span class="comment">// tracebacktrap is like traceback but expects that the PC and SP were obtained</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span><span class="comment">// from a trap, not from gp-&gt;sched or gp-&gt;syscallpc/gp-&gt;syscallsp or getcallerpc/getcallersp.</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span><span class="comment">// Because they are from a trap instead of from a saved pair,</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span><span class="comment">// the initial PC must not be rewound to the previous instruction.</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span><span class="comment">// (All the saved pairs record a PC that is a return address, so we</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span><span class="comment">// rewind it into the CALL instruction.)</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span><span class="comment">// If gp.m.libcall{g,pc,sp} information is available, it uses that information in preference to</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span><span class="comment">// the pc/sp/lr passed in.</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>func tracebacktrap(pc, sp, lr uintptr, gp *g) {
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	if gp.m.libcallsp != 0 {
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re in C code somewhere, traceback from the saved position.</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		traceback1(gp.m.libcallpc, gp.m.libcallsp, 0, gp.m.libcallg.ptr(), 0)
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		return
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	}
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	traceback1(pc, sp, lr, gp, unwindTrap)
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>}
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>func traceback1(pc, sp, lr uintptr, gp *g, flags unwindFlags) {
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	<span class="comment">// If the goroutine is in cgo, and we have a cgo traceback, print that.</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	if iscgo &amp;&amp; gp.m != nil &amp;&amp; gp.m.ncgo &gt; 0 &amp;&amp; gp.syscallsp != 0 &amp;&amp; gp.m.cgoCallers != nil &amp;&amp; gp.m.cgoCallers[0] != 0 {
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		<span class="comment">// Lock cgoCallers so that a signal handler won&#39;t</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		<span class="comment">// change it, copy the array, reset it, unlock it.</span>
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		<span class="comment">// We are locked to the thread and are not running</span>
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>		<span class="comment">// concurrently with a signal handler.</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>		<span class="comment">// We just have to stop a signal handler from interrupting</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>		<span class="comment">// in the middle of our copy.</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		gp.m.cgoCallersUse.Store(1)
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		cgoCallers := *gp.m.cgoCallers
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>		gp.m.cgoCallers[0] = 0
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		gp.m.cgoCallersUse.Store(0)
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>		printCgoTraceback(&amp;cgoCallers)
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	}
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	if readgstatus(gp)&amp;^_Gscan == _Gsyscall {
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>		<span class="comment">// Override registers if blocked in system call.</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		pc = gp.syscallpc
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>		sp = gp.syscallsp
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		flags &amp;^= unwindTrap
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	}
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	if gp.m != nil &amp;&amp; gp.m.vdsoSP != 0 {
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		<span class="comment">// Override registers if running in VDSO. This comes after the</span>
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		<span class="comment">// _Gsyscall check to cover VDSO calls after entersyscall.</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		pc = gp.m.vdsoPC
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		sp = gp.m.vdsoSP
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		flags &amp;^= unwindTrap
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	}
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	<span class="comment">// Print traceback.</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	<span class="comment">// We print the first tracebackInnerFrames frames, and the last</span>
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	<span class="comment">// tracebackOuterFrames frames. There are many possible approaches to this.</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	<span class="comment">// There are various complications to this:</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	<span class="comment">// - We&#39;d prefer to walk the stack once because in really bad situations</span>
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	<span class="comment">//   traceback may crash (and we want as much output as possible) or the stack</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	<span class="comment">//   may be changing.</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	<span class="comment">// - Each physical frame can represent several logical frames, so we might</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	<span class="comment">//   have to pause in the middle of a physical frame and pick up in the middle</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	<span class="comment">//   of a physical frame.</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	<span class="comment">// - The cgo symbolizer can expand a cgo PC to more than one logical frame,</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	<span class="comment">//   and involves juggling state on the C side that we don&#39;t manage. Since its</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	<span class="comment">//   expansion state is managed on the C side, we can&#39;t capture the expansion</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	<span class="comment">//   state part way through, and because the output strings are managed on the</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	<span class="comment">//   C side, we can&#39;t capture the output. Thus, our only choice is to replay a</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	<span class="comment">//   whole expansion, potentially discarding some of it.</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>	<span class="comment">// Rejected approaches:</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>	<span class="comment">// - Do two passes where the first pass just counts and the second pass does</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	<span class="comment">//   all the printing. This is undesirable if the stack is corrupted or changing</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	<span class="comment">//   because we won&#39;t see a partial stack if we panic.</span>
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	<span class="comment">// - Keep a ring buffer of the last N logical frames and use this to print</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	<span class="comment">//   the bottom frames once we reach the end of the stack. This works, but</span>
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	<span class="comment">//   requires keeping a surprising amount of state on the stack, and we have</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	<span class="comment">//   to run the cgo symbolizer twice—once to count frames, and a second to</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	<span class="comment">//   print them—since we can&#39;t retain the strings it returns.</span>
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	<span class="comment">// Instead, we print the outer frames, and if we reach that limit, we clone</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	<span class="comment">// the unwinder, count the remaining frames, and then skip forward and</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>	<span class="comment">// finish printing from the clone. This makes two passes over the outer part</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	<span class="comment">// of the stack, but the single pass over the inner part ensures that&#39;s</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	<span class="comment">// printed immediately and not revisited. It keeps minimal state on the</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	<span class="comment">// stack. And through a combination of skip counts and limits, we can do all</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	<span class="comment">// of the steps we need with a single traceback printer implementation.</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	<span class="comment">// We could be more lax about exactly how many frames we print, for example</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	<span class="comment">// always stopping and resuming on physical frame boundaries, or at least</span>
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>	<span class="comment">// cgo expansion boundaries. It&#39;s not clear that&#39;s much simpler.</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	flags |= unwindPrintErrors
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	var u unwinder
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	tracebackWithRuntime := func(showRuntime bool) int {
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		const maxInt int = 0x7fffffff
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>		u.initAt(pc, sp, lr, gp, flags)
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		n, lastN := traceback2(&amp;u, showRuntime, 0, tracebackInnerFrames)
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		if n &lt; tracebackInnerFrames {
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>			<span class="comment">// We printed the whole stack.</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>			return n
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		}
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>		<span class="comment">// Clone the unwinder and figure out how many frames are left. This</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>		<span class="comment">// count will include any logical frames already printed for u&#39;s current</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>		<span class="comment">// physical frame.</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>		u2 := u
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>		remaining, _ := traceback2(&amp;u, showRuntime, maxInt, 0)
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>		elide := remaining - lastN - tracebackOuterFrames
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>		if elide &gt; 0 {
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>			print(&#34;...&#34;, elide, &#34; frames elided...\n&#34;)
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>			traceback2(&amp;u2, showRuntime, lastN+elide, tracebackOuterFrames)
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>		} else if elide &lt;= 0 {
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>			<span class="comment">// There are tracebackOuterFrames or fewer frames left to print.</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>			<span class="comment">// Just print the rest of the stack.</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>			traceback2(&amp;u2, showRuntime, lastN, tracebackOuterFrames)
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>		}
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>		return n
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>	}
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	<span class="comment">// By default, omits runtime frames. If that means we print nothing at all,</span>
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	<span class="comment">// repeat forcing all frames printed.</span>
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	if tracebackWithRuntime(false) == 0 {
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		tracebackWithRuntime(true)
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>	}
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>	printcreatedby(gp)
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>	if gp.ancestors == nil {
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>		return
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>	}
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>	for _, ancestor := range *gp.ancestors {
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		printAncestorTraceback(ancestor)
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>	}
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>}
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>
<span id="L953" class="ln">   953&nbsp;&nbsp;</span><span class="comment">// traceback2 prints a stack trace starting at u. It skips the first &#34;skip&#34;</span>
<span id="L954" class="ln">   954&nbsp;&nbsp;</span><span class="comment">// logical frames, after which it prints at most &#34;max&#34; logical frames. It</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span><span class="comment">// returns n, which is the number of logical frames skipped and printed, and</span>
<span id="L956" class="ln">   956&nbsp;&nbsp;</span><span class="comment">// lastN, which is the number of logical frames skipped or printed just in the</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span><span class="comment">// physical frame that u references.</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>func traceback2(u *unwinder, showRuntime bool, skip, max int) (n, lastN int) {
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>	<span class="comment">// commitFrame commits to a logical frame and returns whether this frame</span>
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>	<span class="comment">// should be printed and whether iteration should stop.</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	commitFrame := func() (pr, stop bool) {
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>		if skip == 0 &amp;&amp; max == 0 {
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>			<span class="comment">// Stop</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>			return false, true
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>		}
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		n++
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>		lastN++
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>		if skip &gt; 0 {
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>			<span class="comment">// Skip</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>			skip--
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>			return false, false
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>		}
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>		<span class="comment">// Print</span>
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>		max--
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>		return true, false
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	}
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	gp := u.g.ptr()
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>	level, _, _ := gotraceback()
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	var cgoBuf [32]uintptr
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	for ; u.valid(); u.next() {
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		lastN = 0
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		f := u.frame.fn
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>		for iu, uf := newInlineUnwinder(f, u.symPC()); uf.valid(); uf = iu.next(uf) {
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>			sf := iu.srcFunc(uf)
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>			callee := u.calleeFuncID
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>			u.calleeFuncID = sf.funcID
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>			if !(showRuntime || showframe(sf, gp, n == 0, callee)) {
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>				continue
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>			}
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>			if pr, stop := commitFrame(); stop {
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>				return
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>			} else if !pr {
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>				continue
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>			}
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>			name := sf.name()
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>			file, line := iu.fileLine(uf)
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>			<span class="comment">// Print during crash.</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>			<span class="comment">//	main(0x1, 0x2, 0x3)</span>
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>			<span class="comment">//		/home/rsc/go/src/runtime/x.go:23 +0xf</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>			printFuncName(name)
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>			print(&#34;(&#34;)
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>			if iu.isInlined(uf) {
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>				print(&#34;...&#34;)
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>			} else {
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>				argp := unsafe.Pointer(u.frame.argp)
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>				printArgs(f, argp, u.symPC())
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>			}
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>			print(&#34;)\n&#34;)
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>			print(&#34;\t&#34;, file, &#34;:&#34;, line)
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>			if !iu.isInlined(uf) {
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>				if u.frame.pc &gt; f.entry() {
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>					print(&#34; +&#34;, hex(u.frame.pc-f.entry()))
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>				}
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>				if gp.m != nil &amp;&amp; gp.m.throwing &gt;= throwTypeRuntime &amp;&amp; gp == gp.m.curg || level &gt;= 2 {
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>					print(&#34; fp=&#34;, hex(u.frame.fp), &#34; sp=&#34;, hex(u.frame.sp), &#34; pc=&#34;, hex(u.frame.pc))
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>				}
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>			}
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>			print(&#34;\n&#34;)
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		}
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		<span class="comment">// Print cgo frames.</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>		if cgoN := u.cgoCallers(cgoBuf[:]); cgoN &gt; 0 {
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>			var arg cgoSymbolizerArg
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>			anySymbolized := false
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>			stop := false
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>			for _, pc := range cgoBuf[:cgoN] {
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>				if cgoSymbolizer == nil {
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>					if pr, stop := commitFrame(); stop {
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>						break
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>					} else if pr {
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>						print(&#34;non-Go function at pc=&#34;, hex(pc), &#34;\n&#34;)
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>					}
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>				} else {
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>					stop = printOneCgoTraceback(pc, commitFrame, &amp;arg)
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>					anySymbolized = true
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>					if stop {
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>						break
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>					}
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>				}
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>			}
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>			if anySymbolized {
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>				<span class="comment">// Free symbolization state.</span>
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>				arg.pc = 0
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>				callCgoSymbolizer(&amp;arg)
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>			}
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>			if stop {
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>				return
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>			}
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>		}
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>	}
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>	return n, 0
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>}
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span><span class="comment">// printAncestorTraceback prints the traceback of the given ancestor.</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span><span class="comment">// TODO: Unify this with gentraceback and CallersFrames.</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>func printAncestorTraceback(ancestor ancestorInfo) {
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>	print(&#34;[originating from goroutine &#34;, ancestor.goid, &#34;]:\n&#34;)
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>	for fidx, pc := range ancestor.pcs {
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>		f := findfunc(pc) <span class="comment">// f previously validated</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>		if showfuncinfo(f.srcFunc(), fidx == 0, abi.FuncIDNormal) {
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>			printAncestorTracebackFuncInfo(f, pc)
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>		}
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>	}
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>	if len(ancestor.pcs) == tracebackInnerFrames {
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>		print(&#34;...additional frames elided...\n&#34;)
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	}
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	<span class="comment">// Show what created goroutine, except main goroutine (goid 1).</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>	f := findfunc(ancestor.gopc)
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	if f.valid() &amp;&amp; showfuncinfo(f.srcFunc(), false, abi.FuncIDNormal) &amp;&amp; ancestor.goid != 1 {
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>		<span class="comment">// In ancestor mode, we&#39;ll already print the goroutine ancestor.</span>
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>		<span class="comment">// Pass 0 for the goid parameter so we don&#39;t print it again.</span>
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>		printcreatedby1(f, ancestor.gopc, 0)
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	}
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>}
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span><span class="comment">// printAncestorTracebackFuncInfo prints the given function info at a given pc</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span><span class="comment">// within an ancestor traceback. The precision of this info is reduced</span>
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span><span class="comment">// due to only have access to the pcs at the time of the caller</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span><span class="comment">// goroutine being created.</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>func printAncestorTracebackFuncInfo(f funcInfo, pc uintptr) {
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>	u, uf := newInlineUnwinder(f, pc)
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>	file, line := u.fileLine(uf)
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	printFuncName(u.srcFunc(uf).name())
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>	print(&#34;(...)\n&#34;)
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>	print(&#34;\t&#34;, file, &#34;:&#34;, line)
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>	if pc &gt; f.entry() {
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>		print(&#34; +&#34;, hex(pc-f.entry()))
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>	}
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>	print(&#34;\n&#34;)
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>}
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>func callers(skip int, pcbuf []uintptr) int {
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>	sp := getcallersp()
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>	pc := getcallerpc()
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	gp := getg()
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>	var n int
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>		var u unwinder
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>		u.initAt(pc, sp, 0, gp, unwindSilentErrors)
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>		n = tracebackPCs(&amp;u, skip, pcbuf)
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>	})
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>	return n
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>}
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>func gcallers(gp *g, skip int, pcbuf []uintptr) int {
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>	var u unwinder
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>	u.init(gp, unwindSilentErrors)
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>	return tracebackPCs(&amp;u, skip, pcbuf)
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>}
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span><span class="comment">// showframe reports whether the frame with the given characteristics should</span>
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span><span class="comment">// be printed during a traceback.</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>func showframe(sf srcFunc, gp *g, firstFrame bool, calleeID abi.FuncID) bool {
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>	mp := getg().m
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>	if mp.throwing &gt;= throwTypeRuntime &amp;&amp; gp != nil &amp;&amp; (gp == mp.curg || gp == mp.caughtsig.ptr()) {
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>		return true
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>	}
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	return showfuncinfo(sf, firstFrame, calleeID)
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>}
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span><span class="comment">// showfuncinfo reports whether a function with the given characteristics should</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span><span class="comment">// be printed during a traceback.</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>func showfuncinfo(sf srcFunc, firstFrame bool, calleeID abi.FuncID) bool {
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>	level, _, _ := gotraceback()
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>	if level &gt; 1 {
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>		<span class="comment">// Show all frames.</span>
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>		return true
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>	}
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>	if sf.funcID == abi.FuncIDWrapper &amp;&amp; elideWrapperCalling(calleeID) {
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>		return false
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>	}
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>	name := sf.name()
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>	<span class="comment">// Special case: always show runtime.gopanic frame</span>
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	<span class="comment">// in the middle of a stack trace, so that we can</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	<span class="comment">// see the boundary between ordinary code and</span>
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	<span class="comment">// panic-induced deferred code.</span>
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>	<span class="comment">// See golang.org/issue/5832.</span>
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>	if name == &#34;runtime.gopanic&#34; &amp;&amp; !firstFrame {
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>		return true
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>	}
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>	return bytealg.IndexByteString(name, &#39;.&#39;) &gt;= 0 &amp;&amp; (!hasPrefix(name, &#34;runtime.&#34;) || isExportedRuntime(name))
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>}
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span><span class="comment">// isExportedRuntime reports whether name is an exported runtime function.</span>
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span><span class="comment">// It is only for runtime functions, so ASCII A-Z is fine.</span>
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span><span class="comment">// TODO: this handles exported functions but not exported methods.</span>
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>func isExportedRuntime(name string) bool {
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>	const n = len(&#34;runtime.&#34;)
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>	return len(name) &gt; n &amp;&amp; name[:n] == &#34;runtime.&#34; &amp;&amp; &#39;A&#39; &lt;= name[n] &amp;&amp; name[n] &lt;= &#39;Z&#39;
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>}
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span><span class="comment">// elideWrapperCalling reports whether a wrapper function that called</span>
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span><span class="comment">// function id should be elided from stack traces.</span>
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>func elideWrapperCalling(id abi.FuncID) bool {
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>	<span class="comment">// If the wrapper called a panic function instead of the</span>
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>	<span class="comment">// wrapped function, we want to include it in stacks.</span>
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>	return !(id == abi.FuncID_gopanic || id == abi.FuncID_sigpanic || id == abi.FuncID_panicwrap)
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>}
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>var gStatusStrings = [...]string{
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>	_Gidle:      &#34;idle&#34;,
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>	_Grunnable:  &#34;runnable&#34;,
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>	_Grunning:   &#34;running&#34;,
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>	_Gsyscall:   &#34;syscall&#34;,
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>	_Gwaiting:   &#34;waiting&#34;,
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>	_Gdead:      &#34;dead&#34;,
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>	_Gcopystack: &#34;copystack&#34;,
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>	_Gpreempted: &#34;preempted&#34;,
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>}
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>func goroutineheader(gp *g) {
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>	level, _, _ := gotraceback()
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>	gpstatus := readgstatus(gp)
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>	isScan := gpstatus&amp;_Gscan != 0
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>	gpstatus &amp;^= _Gscan <span class="comment">// drop the scan bit</span>
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>	<span class="comment">// Basic string status</span>
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>	var status string
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>	if 0 &lt;= gpstatus &amp;&amp; gpstatus &lt; uint32(len(gStatusStrings)) {
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>		status = gStatusStrings[gpstatus]
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>	} else {
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>		status = &#34;???&#34;
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>	}
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>	<span class="comment">// Override.</span>
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>	if gpstatus == _Gwaiting &amp;&amp; gp.waitreason != waitReasonZero {
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>		status = gp.waitreason.String()
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>	}
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>	<span class="comment">// approx time the G is blocked, in minutes</span>
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>	var waitfor int64
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>	if (gpstatus == _Gwaiting || gpstatus == _Gsyscall) &amp;&amp; gp.waitsince != 0 {
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>		waitfor = (nanotime() - gp.waitsince) / 60e9
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>	}
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>	print(&#34;goroutine &#34;, gp.goid)
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>	if gp.m != nil &amp;&amp; gp.m.throwing &gt;= throwTypeRuntime &amp;&amp; gp == gp.m.curg || level &gt;= 2 {
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>		print(&#34; gp=&#34;, gp)
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>		if gp.m != nil {
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>			print(&#34; m=&#34;, gp.m.id, &#34; mp=&#34;, gp.m)
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>		} else {
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>			print(&#34; m=nil&#34;)
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>		}
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>	}
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>	print(&#34; [&#34;, status)
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>	if isScan {
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>		print(&#34; (scan)&#34;)
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>	}
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>	if waitfor &gt;= 1 {
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>		print(&#34;, &#34;, waitfor, &#34; minutes&#34;)
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>	}
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>	if gp.lockedm != 0 {
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>		print(&#34;, locked to thread&#34;)
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>	}
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>	print(&#34;]:\n&#34;)
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>}
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>func tracebackothers(me *g) {
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>	level, _, _ := gotraceback()
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>	<span class="comment">// Show the current goroutine first, if we haven&#39;t already.</span>
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	curgp := getg().m.curg
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>	if curgp != nil &amp;&amp; curgp != me {
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>		goroutineheader(curgp)
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>		traceback(^uintptr(0), ^uintptr(0), 0, curgp)
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>	}
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>	<span class="comment">// We can&#39;t call locking forEachG here because this may be during fatal</span>
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>	<span class="comment">// throw/panic, where locking could be out-of-order or a direct</span>
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>	<span class="comment">// deadlock.</span>
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>	<span class="comment">// Instead, use forEachGRace, which requires no locking. We don&#39;t lock</span>
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>	<span class="comment">// against concurrent creation of new Gs, but even with allglock we may</span>
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>	<span class="comment">// miss Gs created after this loop.</span>
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>	forEachGRace(func(gp *g) {
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>		if gp == me || gp == curgp || readgstatus(gp) == _Gdead || isSystemGoroutine(gp, false) &amp;&amp; level &lt; 2 {
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>			return
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>		}
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>		goroutineheader(gp)
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>		<span class="comment">// Note: gp.m == getg().m occurs when tracebackothers is called</span>
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>		<span class="comment">// from a signal handler initiated during a systemstack call.</span>
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>		<span class="comment">// The original G is still in the running state, and we want to</span>
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>		<span class="comment">// print its stack.</span>
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>		if gp.m != getg().m &amp;&amp; readgstatus(gp)&amp;^_Gscan == _Grunning {
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>			print(&#34;\tgoroutine running on other thread; stack unavailable\n&#34;)
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>			printcreatedby(gp)
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>		} else {
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>			traceback(^uintptr(0), ^uintptr(0), 0, gp)
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>		}
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>	})
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>}
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span><span class="comment">// tracebackHexdump hexdumps part of stk around frame.sp and frame.fp</span>
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span><span class="comment">// for debugging purposes. If the address bad is included in the</span>
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span><span class="comment">// hexdumped range, it will mark it as well.</span>
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>func tracebackHexdump(stk stack, frame *stkframe, bad uintptr) {
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>	const expand = 32 * goarch.PtrSize
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>	const maxExpand = 256 * goarch.PtrSize
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	<span class="comment">// Start around frame.sp.</span>
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	lo, hi := frame.sp, frame.sp
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>	<span class="comment">// Expand to include frame.fp.</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	if frame.fp != 0 &amp;&amp; frame.fp &lt; lo {
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>		lo = frame.fp
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>	}
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>	if frame.fp != 0 &amp;&amp; frame.fp &gt; hi {
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>		hi = frame.fp
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>	}
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>	<span class="comment">// Expand a bit more.</span>
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>	lo, hi = lo-expand, hi+expand
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>	<span class="comment">// But don&#39;t go too far from frame.sp.</span>
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>	if lo &lt; frame.sp-maxExpand {
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>		lo = frame.sp - maxExpand
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>	}
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>	if hi &gt; frame.sp+maxExpand {
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>		hi = frame.sp + maxExpand
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>	}
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>	<span class="comment">// And don&#39;t go outside the stack bounds.</span>
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>	if lo &lt; stk.lo {
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>		lo = stk.lo
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>	}
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>	if hi &gt; stk.hi {
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>		hi = stk.hi
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>	}
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>	<span class="comment">// Print the hex dump.</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>	print(&#34;stack: frame={sp:&#34;, hex(frame.sp), &#34;, fp:&#34;, hex(frame.fp), &#34;} stack=[&#34;, hex(stk.lo), &#34;,&#34;, hex(stk.hi), &#34;)\n&#34;)
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>	hexdumpWords(lo, hi, func(p uintptr) byte {
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>		switch p {
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>		case frame.fp:
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>			return &#39;&gt;&#39;
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>		case frame.sp:
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>			return &#39;&lt;&#39;
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>		case bad:
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>			return &#39;!&#39;
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>		}
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>		return 0
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>	})
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>}
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span><span class="comment">// isSystemGoroutine reports whether the goroutine g must be omitted</span>
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span><span class="comment">// in stack dumps and deadlock detector. This is any goroutine that</span>
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span><span class="comment">// starts at a runtime.* entry point, except for runtime.main,</span>
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span><span class="comment">// runtime.handleAsyncEvent (wasm only) and sometimes runtime.runfinq.</span>
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span><span class="comment">// If fixed is true, any goroutine that can vary between user and</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span><span class="comment">// system (that is, the finalizer goroutine) is considered a user</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span><span class="comment">// goroutine.</span>
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>func isSystemGoroutine(gp *g, fixed bool) bool {
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>	<span class="comment">// Keep this in sync with internal/trace.IsSystemGoroutine.</span>
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>	f := findfunc(gp.startpc)
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>	if !f.valid() {
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>		return false
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>	}
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>	if f.funcID == abi.FuncID_runtime_main || f.funcID == abi.FuncID_corostart || f.funcID == abi.FuncID_handleAsyncEvent {
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>		return false
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>	}
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>	if f.funcID == abi.FuncID_runfinq {
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>		<span class="comment">// We include the finalizer goroutine if it&#39;s calling</span>
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>		<span class="comment">// back into user code.</span>
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>		if fixed {
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>			<span class="comment">// This goroutine can vary. In fixed mode,</span>
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>			<span class="comment">// always consider it a user goroutine.</span>
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>			return false
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>		}
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>		return fingStatus.Load()&amp;fingRunningFinalizer == 0
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>	}
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>	return hasPrefix(funcname(f), &#34;runtime.&#34;)
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>}
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span><span class="comment">// SetCgoTraceback records three C functions to use to gather</span>
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span><span class="comment">// traceback information from C code and to convert that traceback</span>
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span><span class="comment">// information into symbolic information. These are used when printing</span>
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span><span class="comment">// stack traces for a program that uses cgo.</span>
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span><span class="comment">// The traceback and context functions may be called from a signal</span>
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span><span class="comment">// handler, and must therefore use only async-signal safe functions.</span>
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span><span class="comment">// The symbolizer function may be called while the program is</span>
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span><span class="comment">// crashing, and so must be cautious about using memory.  None of the</span>
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span><span class="comment">// functions may call back into Go.</span>
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span><span class="comment">// The context function will be called with a single argument, a</span>
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span><span class="comment">// pointer to a struct:</span>
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span><span class="comment">//	struct {</span>
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span><span class="comment">//		Context uintptr</span>
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span><span class="comment">// In C syntax, this struct will be</span>
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span><span class="comment">//	struct {</span>
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span><span class="comment">//		uintptr_t Context;</span>
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span><span class="comment">//	};</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span><span class="comment">// If the Context field is 0, the context function is being called to</span>
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span><span class="comment">// record the current traceback context. It should record in the</span>
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span><span class="comment">// Context field whatever information is needed about the current</span>
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span><span class="comment">// point of execution to later produce a stack trace, probably the</span>
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span><span class="comment">// stack pointer and PC. In this case the context function will be</span>
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span><span class="comment">// called from C code.</span>
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span><span class="comment">// If the Context field is not 0, then it is a value returned by a</span>
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span><span class="comment">// previous call to the context function. This case is called when the</span>
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span><span class="comment">// context is no longer needed; that is, when the Go code is returning</span>
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span><span class="comment">// to its C code caller. This permits the context function to release</span>
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span><span class="comment">// any associated resources.</span>
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span><span class="comment">// While it would be correct for the context function to record a</span>
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span><span class="comment">// complete a stack trace whenever it is called, and simply copy that</span>
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span><span class="comment">// out in the traceback function, in a typical program the context</span>
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span><span class="comment">// function will be called many times without ever recording a</span>
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span><span class="comment">// traceback for that context. Recording a complete stack trace in a</span>
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span><span class="comment">// call to the context function is likely to be inefficient.</span>
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span><span class="comment">// The traceback function will be called with a single argument, a</span>
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span><span class="comment">// pointer to a struct:</span>
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span><span class="comment">//	struct {</span>
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span><span class="comment">//		Context    uintptr</span>
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span><span class="comment">//		SigContext uintptr</span>
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span><span class="comment">//		Buf        *uintptr</span>
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span><span class="comment">//		Max        uintptr</span>
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span><span class="comment">// In C syntax, this struct will be</span>
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span><span class="comment">//	struct {</span>
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span><span class="comment">//		uintptr_t  Context;</span>
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span><span class="comment">//		uintptr_t  SigContext;</span>
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span><span class="comment">//		uintptr_t* Buf;</span>
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span><span class="comment">//		uintptr_t  Max;</span>
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span><span class="comment">//	};</span>
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span><span class="comment">// The Context field will be zero to gather a traceback from the</span>
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span><span class="comment">// current program execution point. In this case, the traceback</span>
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span><span class="comment">// function will be called from C code.</span>
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span><span class="comment">// Otherwise Context will be a value previously returned by a call to</span>
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span><span class="comment">// the context function. The traceback function should gather a stack</span>
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span><span class="comment">// trace from that saved point in the program execution. The traceback</span>
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span><span class="comment">// function may be called from an execution thread other than the one</span>
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span><span class="comment">// that recorded the context, but only when the context is known to be</span>
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span><span class="comment">// valid and unchanging. The traceback function may also be called</span>
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span><span class="comment">// deeper in the call stack on the same thread that recorded the</span>
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span><span class="comment">// context. The traceback function may be called multiple times with</span>
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span><span class="comment">// the same Context value; it will usually be appropriate to cache the</span>
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span><span class="comment">// result, if possible, the first time this is called for a specific</span>
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span><span class="comment">// context value.</span>
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span><span class="comment">// If the traceback function is called from a signal handler on a Unix</span>
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span><span class="comment">// system, SigContext will be the signal context argument passed to</span>
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span><span class="comment">// the signal handler (a C ucontext_t* cast to uintptr_t). This may be</span>
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span><span class="comment">// used to start tracing at the point where the signal occurred. If</span>
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span><span class="comment">// the traceback function is not called from a signal handler,</span>
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span><span class="comment">// SigContext will be zero.</span>
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span><span class="comment">// Buf is where the traceback information should be stored. It should</span>
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span><span class="comment">// be PC values, such that Buf[0] is the PC of the caller, Buf[1] is</span>
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span><span class="comment">// the PC of that function&#39;s caller, and so on.  Max is the maximum</span>
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span><span class="comment">// number of entries to store.  The function should store a zero to</span>
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span><span class="comment">// indicate the top of the stack, or that the caller is on a different</span>
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span><span class="comment">// stack, presumably a Go stack.</span>
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span><span class="comment">// Unlike runtime.Callers, the PC values returned should, when passed</span>
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span><span class="comment">// to the symbolizer function, return the file/line of the call</span>
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span><span class="comment">// instruction.  No additional subtraction is required or appropriate.</span>
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span><span class="comment">// On all platforms, the traceback function is invoked when a call from</span>
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span><span class="comment">// Go to C to Go requests a stack trace. On linux/amd64, linux/ppc64le,</span>
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span><span class="comment">// linux/arm64, and freebsd/amd64, the traceback function is also invoked</span>
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span><span class="comment">// when a signal is received by a thread that is executing a cgo call.</span>
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span><span class="comment">// The traceback function should not make assumptions about when it is</span>
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span><span class="comment">// called, as future versions of Go may make additional calls.</span>
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span><span class="comment">// The symbolizer function will be called with a single argument, a</span>
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span><span class="comment">// pointer to a struct:</span>
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span><span class="comment">//	struct {</span>
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span><span class="comment">//		PC      uintptr // program counter to fetch information for</span>
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span><span class="comment">//		File    *byte   // file name (NUL terminated)</span>
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span><span class="comment">//		Lineno  uintptr // line number</span>
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span><span class="comment">//		Func    *byte   // function name (NUL terminated)</span>
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span><span class="comment">//		Entry   uintptr // function entry point</span>
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span><span class="comment">//		More    uintptr // set non-zero if more info for this PC</span>
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span><span class="comment">//		Data    uintptr // unused by runtime, available for function</span>
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span><span class="comment">// In C syntax, this struct will be</span>
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span><span class="comment">//	struct {</span>
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span><span class="comment">//		uintptr_t PC;</span>
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span><span class="comment">//		char*     File;</span>
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span><span class="comment">//		uintptr_t Lineno;</span>
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span><span class="comment">//		char*     Func;</span>
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span><span class="comment">//		uintptr_t Entry;</span>
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span><span class="comment">//		uintptr_t More;</span>
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span><span class="comment">//		uintptr_t Data;</span>
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span><span class="comment">//	};</span>
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span><span class="comment">// The PC field will be a value returned by a call to the traceback</span>
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span><span class="comment">// function.</span>
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span><span class="comment">// The first time the function is called for a particular traceback,</span>
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span><span class="comment">// all the fields except PC will be 0. The function should fill in the</span>
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span><span class="comment">// other fields if possible, setting them to 0/nil if the information</span>
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span><span class="comment">// is not available. The Data field may be used to store any useful</span>
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span><span class="comment">// information across calls. The More field should be set to non-zero</span>
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span><span class="comment">// if there is more information for this PC, zero otherwise. If More</span>
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span><span class="comment">// is set non-zero, the function will be called again with the same</span>
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span><span class="comment">// PC, and may return different information (this is intended for use</span>
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span><span class="comment">// with inlined functions). If More is zero, the function will be</span>
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span><span class="comment">// called with the next PC value in the traceback. When the traceback</span>
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span><span class="comment">// is complete, the function will be called once more with PC set to</span>
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span><span class="comment">// zero; this may be used to free any information. Each call will</span>
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span><span class="comment">// leave the fields of the struct set to the same values they had upon</span>
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span><span class="comment">// return, except for the PC field when the More field is zero. The</span>
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span><span class="comment">// function must not keep a copy of the struct pointer between calls.</span>
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span><span class="comment">// When calling SetCgoTraceback, the version argument is the version</span>
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span><span class="comment">// number of the structs that the functions expect to receive.</span>
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span><span class="comment">// Currently this must be zero.</span>
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span><span class="comment">// The symbolizer function may be nil, in which case the results of</span>
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span><span class="comment">// the traceback function will be displayed as numbers. If the</span>
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span><span class="comment">// traceback function is nil, the symbolizer function will never be</span>
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span><span class="comment">// called. The context function may be nil, in which case the</span>
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span><span class="comment">// traceback function will only be called with the context field set</span>
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span><span class="comment">// to zero.  If the context function is nil, then calls from Go to C</span>
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span><span class="comment">// to Go will not show a traceback for the C portion of the call stack.</span>
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span><span class="comment">// SetCgoTraceback should be called only once, ideally from an init function.</span>
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span>func SetCgoTraceback(version int, traceback, context, symbolizer unsafe.Pointer) {
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span>	if version != 0 {
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span>		panic(&#34;unsupported version&#34;)
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span>	}
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span>
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span>	if cgoTraceback != nil &amp;&amp; cgoTraceback != traceback ||
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span>		cgoContext != nil &amp;&amp; cgoContext != context ||
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span>		cgoSymbolizer != nil &amp;&amp; cgoSymbolizer != symbolizer {
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span>		panic(&#34;call SetCgoTraceback only once&#34;)
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span>	}
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span>
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span>	cgoTraceback = traceback
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span>	cgoContext = context
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>	cgoSymbolizer = symbolizer
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span>	<span class="comment">// The context function is called when a C function calls a Go</span>
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span>	<span class="comment">// function. As such it is only called by C code in runtime/cgo.</span>
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span>	if _cgo_set_context_function != nil {
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span>		cgocall(_cgo_set_context_function, context)
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span>	}
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span>}
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span>
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span>var cgoTraceback unsafe.Pointer
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span>var cgoContext unsafe.Pointer
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span>var cgoSymbolizer unsafe.Pointer
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span>
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span><span class="comment">// cgoTracebackArg is the type passed to cgoTraceback.</span>
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span>type cgoTracebackArg struct {
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>	context    uintptr
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>	sigContext uintptr
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>	buf        *uintptr
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>	max        uintptr
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span>}
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span>
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span><span class="comment">// cgoContextArg is the type passed to the context function.</span>
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span>type cgoContextArg struct {
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>	context uintptr
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>}
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span><span class="comment">// cgoSymbolizerArg is the type passed to cgoSymbolizer.</span>
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span>type cgoSymbolizerArg struct {
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span>	pc       uintptr
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span>	file     *byte
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>	lineno   uintptr
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>	funcName *byte
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>	entry    uintptr
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>	more     uintptr
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span>	data     uintptr
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span>}
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span>
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span><span class="comment">// printCgoTraceback prints a traceback of callers.</span>
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span>func printCgoTraceback(callers *cgoCallers) {
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span>	if cgoSymbolizer == nil {
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span>		for _, c := range callers {
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span>			if c == 0 {
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span>				break
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span>			}
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span>			print(&#34;non-Go function at pc=&#34;, hex(c), &#34;\n&#34;)
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span>		}
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span>		return
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span>	}
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span>
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span>	commitFrame := func() (pr, stop bool) { return true, false }
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span>	var arg cgoSymbolizerArg
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span>	for _, c := range callers {
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span>		if c == 0 {
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span>			break
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span>		}
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span>		printOneCgoTraceback(c, commitFrame, &amp;arg)
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span>	}
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span>	arg.pc = 0
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span>	callCgoSymbolizer(&amp;arg)
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span>}
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span>
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span><span class="comment">// printOneCgoTraceback prints the traceback of a single cgo caller.</span>
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span><span class="comment">// This can print more than one line because of inlining.</span>
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span><span class="comment">// It returns the &#34;stop&#34; result of commitFrame.</span>
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span>func printOneCgoTraceback(pc uintptr, commitFrame func() (pr, stop bool), arg *cgoSymbolizerArg) bool {
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span>	arg.pc = pc
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span>	for {
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span>		if pr, stop := commitFrame(); stop {
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>			return true
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span>		} else if !pr {
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span>			continue
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span>		}
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span>
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span>		callCgoSymbolizer(arg)
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>		if arg.funcName != nil {
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>			<span class="comment">// Note that we don&#39;t print any argument</span>
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>			<span class="comment">// information here, not even parentheses.</span>
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>			<span class="comment">// The symbolizer must add that if appropriate.</span>
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>			println(gostringnocopy(arg.funcName))
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>		} else {
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>			println(&#34;non-Go function&#34;)
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span>		}
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span>		print(&#34;\t&#34;)
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span>		if arg.file != nil {
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span>			print(gostringnocopy(arg.file), &#34;:&#34;, arg.lineno, &#34; &#34;)
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span>		}
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>		print(&#34;pc=&#34;, hex(pc), &#34;\n&#34;)
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>		if arg.more == 0 {
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span>			return false
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span>		}
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span>	}
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span>}
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span>
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span><span class="comment">// callCgoSymbolizer calls the cgoSymbolizer function.</span>
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span>func callCgoSymbolizer(arg *cgoSymbolizerArg) {
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span>	call := cgocall
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span>	if panicking.Load() &gt; 0 || getg().m.curg != getg() {
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span>		<span class="comment">// We do not want to call into the scheduler when panicking</span>
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span>		<span class="comment">// or when on the system stack.</span>
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span>		call = asmcgocall
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>	}
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>	if msanenabled {
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>		msanwrite(unsafe.Pointer(arg), unsafe.Sizeof(cgoSymbolizerArg{}))
<span id="L1618" class="ln">  1618&nbsp;&nbsp;</span>	}
<span id="L1619" class="ln">  1619&nbsp;&nbsp;</span>	if asanenabled {
<span id="L1620" class="ln">  1620&nbsp;&nbsp;</span>		asanwrite(unsafe.Pointer(arg), unsafe.Sizeof(cgoSymbolizerArg{}))
<span id="L1621" class="ln">  1621&nbsp;&nbsp;</span>	}
<span id="L1622" class="ln">  1622&nbsp;&nbsp;</span>	call(cgoSymbolizer, noescape(unsafe.Pointer(arg)))
<span id="L1623" class="ln">  1623&nbsp;&nbsp;</span>}
<span id="L1624" class="ln">  1624&nbsp;&nbsp;</span>
<span id="L1625" class="ln">  1625&nbsp;&nbsp;</span><span class="comment">// cgoContextPCs gets the PC values from a cgo traceback.</span>
<span id="L1626" class="ln">  1626&nbsp;&nbsp;</span>func cgoContextPCs(ctxt uintptr, buf []uintptr) {
<span id="L1627" class="ln">  1627&nbsp;&nbsp;</span>	if cgoTraceback == nil {
<span id="L1628" class="ln">  1628&nbsp;&nbsp;</span>		return
<span id="L1629" class="ln">  1629&nbsp;&nbsp;</span>	}
<span id="L1630" class="ln">  1630&nbsp;&nbsp;</span>	call := cgocall
<span id="L1631" class="ln">  1631&nbsp;&nbsp;</span>	if panicking.Load() &gt; 0 || getg().m.curg != getg() {
<span id="L1632" class="ln">  1632&nbsp;&nbsp;</span>		<span class="comment">// We do not want to call into the scheduler when panicking</span>
<span id="L1633" class="ln">  1633&nbsp;&nbsp;</span>		<span class="comment">// or when on the system stack.</span>
<span id="L1634" class="ln">  1634&nbsp;&nbsp;</span>		call = asmcgocall
<span id="L1635" class="ln">  1635&nbsp;&nbsp;</span>	}
<span id="L1636" class="ln">  1636&nbsp;&nbsp;</span>	arg := cgoTracebackArg{
<span id="L1637" class="ln">  1637&nbsp;&nbsp;</span>		context: ctxt,
<span id="L1638" class="ln">  1638&nbsp;&nbsp;</span>		buf:     (*uintptr)(noescape(unsafe.Pointer(&amp;buf[0]))),
<span id="L1639" class="ln">  1639&nbsp;&nbsp;</span>		max:     uintptr(len(buf)),
<span id="L1640" class="ln">  1640&nbsp;&nbsp;</span>	}
<span id="L1641" class="ln">  1641&nbsp;&nbsp;</span>	if msanenabled {
<span id="L1642" class="ln">  1642&nbsp;&nbsp;</span>		msanwrite(unsafe.Pointer(&amp;arg), unsafe.Sizeof(arg))
<span id="L1643" class="ln">  1643&nbsp;&nbsp;</span>	}
<span id="L1644" class="ln">  1644&nbsp;&nbsp;</span>	if asanenabled {
<span id="L1645" class="ln">  1645&nbsp;&nbsp;</span>		asanwrite(unsafe.Pointer(&amp;arg), unsafe.Sizeof(arg))
<span id="L1646" class="ln">  1646&nbsp;&nbsp;</span>	}
<span id="L1647" class="ln">  1647&nbsp;&nbsp;</span>	call(cgoTraceback, noescape(unsafe.Pointer(&amp;arg)))
<span id="L1648" class="ln">  1648&nbsp;&nbsp;</span>}
<span id="L1649" class="ln">  1649&nbsp;&nbsp;</span>
</pre><p><a href="traceback.go?m=text">View as plain text</a></p>

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
