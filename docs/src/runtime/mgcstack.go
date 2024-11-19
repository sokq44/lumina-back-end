<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/mgcstack.go - Go Documentation Server</title>

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
<a href="mgcstack.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">mgcstack.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Garbage collector: stack objects and stack tracing</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// See the design doc at https://docs.google.com/document/d/1un-Jn47yByHL7I0aVIP_uVCMxjdM5mpelJhiKlIqxkE/edit?usp=sharing</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Also see issue 22350.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// Stack tracing solves the problem of determining which parts of the</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// stack are live and should be scanned. It runs as part of scanning</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// a single goroutine stack.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// Normally determining which parts of the stack are live is easy to</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// do statically, as user code has explicit references (reads and</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// writes) to stack variables. The compiler can do a simple dataflow</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// analysis to determine liveness of stack variables at every point in</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// the code. See cmd/compile/internal/gc/plive.go for that analysis.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// However, when we take the address of a stack variable, determining</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// whether that variable is still live is less clear. We can still</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// look for static accesses, but accesses through a pointer to the</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// variable are difficult in general to track statically. That pointer</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// can be passed among functions on the stack, conditionally retained,</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// etc.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// Instead, we will track pointers to stack variables dynamically.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// All pointers to stack-allocated variables will themselves be on the</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// stack somewhere (or in associated locations, like defer records), so</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// we can find them all efficiently.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// Stack tracing is organized as a mini garbage collection tracing</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// pass. The objects in this garbage collection are all the variables</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// on the stack whose address is taken, and which themselves contain a</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// pointer. We call these variables &#34;stack objects&#34;.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// We begin by determining all the stack objects on the stack and all</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// the statically live pointers that may point into the stack. We then</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// process each pointer to see if it points to a stack object. If it</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// does, we scan that stack object. It may contain pointers into the</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// heap, in which case those pointers are passed to the main garbage</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// collection. It may also contain pointers into the stack, in which</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// case we add them to our set of stack pointers.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// Once we&#39;re done processing all the pointers (including the ones we</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// added during processing), we&#39;ve found all the stack objects that</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// are live. Any dead stack objects are not scanned and their contents</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// will not keep heap objects live. Unlike the main garbage</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// collection, we can&#39;t sweep the dead stack objects; they live on in</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// a moribund state until the stack frame that contains them is</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// popped.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// A stack can look like this:</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// +----------+</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// | foo()    |</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// | +------+ |</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// | |  A   | | &lt;---\</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// | +------+ |     |</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// |          |     |</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// | +------+ |     |</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// | |  B   | |     |</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// | +------+ |     |</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// |          |     |</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// +----------+     |</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// | bar()    |     |</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// | +------+ |     |</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// | |  C   | | &lt;-\ |</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// | +----|-+ |   | |</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// |      |   |   | |</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// | +----v-+ |   | |</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// | |  D  ---------/</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// | +------+ |   |</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// |          |   |</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// +----------+   |</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// | baz()    |   |</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// | +------+ |   |</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// | |  E  -------/</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// | +------+ |</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// |      ^   |</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// | F: --/   |</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// |          |</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// +----------+</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// foo() calls bar() calls baz(). Each has a frame on the stack.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// foo() has stack objects A and B.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// bar() has stack objects C and D, with C pointing to D and D pointing to A.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// baz() has a stack object E pointing to C, and a local variable F pointing to E.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// Starting from the pointer in local variable F, we will eventually</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// scan all of E, C, D, and A (in that order). B is never scanned</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// because there is no live pointer to it. If B is also statically</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// dead (meaning that foo() never accesses B again after it calls</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// bar()), then B&#39;s pointers into the heap are not considered live.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>package runtime
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>import (
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>const stackTraceDebug = false
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// Buffer for pointers found during stack tracing.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// Must be smaller than or equal to workbuf.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>type stackWorkBuf struct {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	stackWorkBufHdr
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	obj [(_WorkbufSize - unsafe.Sizeof(stackWorkBufHdr{})) / goarch.PtrSize]uintptr
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// Header declaration must come after the buf declaration above, because of issue #14620.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>type stackWorkBufHdr struct {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	workbufhdr
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	next *stackWorkBuf <span class="comment">// linked list of workbufs</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// Note: we could theoretically repurpose lfnode.next as this next pointer.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// It would save 1 word, but that probably isn&#39;t worth busting open</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// the lfnode API.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// Buffer for stack objects found on a goroutine stack.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// Must be smaller than or equal to workbuf.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>type stackObjectBuf struct {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	stackObjectBufHdr
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	obj [(_WorkbufSize - unsafe.Sizeof(stackObjectBufHdr{})) / unsafe.Sizeof(stackObject{})]stackObject
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>type stackObjectBufHdr struct {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	_ sys.NotInHeap
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	workbufhdr
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	next *stackObjectBuf
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>func init() {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	if unsafe.Sizeof(stackWorkBuf{}) &gt; unsafe.Sizeof(workbuf{}) {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		panic(&#34;stackWorkBuf too big&#34;)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	if unsafe.Sizeof(stackObjectBuf{}) &gt; unsafe.Sizeof(workbuf{}) {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		panic(&#34;stackObjectBuf too big&#34;)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// A stackObject represents a variable on the stack that has had</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// its address taken.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>type stackObject struct {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	_     sys.NotInHeap
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	off   uint32             <span class="comment">// offset above stack.lo</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	size  uint32             <span class="comment">// size of object</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	r     *stackObjectRecord <span class="comment">// info of the object (for ptr/nonptr bits). nil if object has been scanned.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	left  *stackObject       <span class="comment">// objects with lower addresses</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	right *stackObject       <span class="comment">// objects with higher addresses</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// obj.r = r, but with no write barrier.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func (obj *stackObject) setRecord(r *stackObjectRecord) {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// Types of stack objects are always in read-only memory, not the heap.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// So not using a write barrier is ok.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	*(*uintptr)(unsafe.Pointer(&amp;obj.r)) = uintptr(unsafe.Pointer(r))
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// A stackScanState keeps track of the state used during the GC walk</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// of a goroutine.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>type stackScanState struct {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// stack limits</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	stack stack
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// conservative indicates that the next frame must be scanned conservatively.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// This applies only to the innermost frame at an async safe-point.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	conservative bool
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// buf contains the set of possible pointers to stack objects.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// Organized as a LIFO linked list of buffers.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// All buffers except possibly the head buffer are full.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	buf     *stackWorkBuf
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	freeBuf *stackWorkBuf <span class="comment">// keep around one free buffer for allocation hysteresis</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// cbuf contains conservative pointers to stack objects. If</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// all pointers to a stack object are obtained via</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// conservative scanning, then the stack object may be dead</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// and may contain dead pointers, so it must be scanned</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// defensively.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	cbuf *stackWorkBuf
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// list of stack objects</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// Objects are in increasing address order.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	head  *stackObjectBuf
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	tail  *stackObjectBuf
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	nobjs int
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// root of binary tree for fast object lookup by address</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// Initialized by buildIndex.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	root *stackObject
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// Add p as a potential pointer to a stack object.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// p must be a stack address.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>func (s *stackScanState) putPtr(p uintptr, conservative bool) {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	if p &lt; s.stack.lo || p &gt;= s.stack.hi {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		throw(&#34;address not a stack address&#34;)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	head := &amp;s.buf
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	if conservative {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		head = &amp;s.cbuf
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	buf := *head
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if buf == nil {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		<span class="comment">// Initial setup.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		buf = (*stackWorkBuf)(unsafe.Pointer(getempty()))
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		buf.nobj = 0
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		buf.next = nil
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		*head = buf
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	} else if buf.nobj == len(buf.obj) {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		if s.freeBuf != nil {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			buf = s.freeBuf
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			s.freeBuf = nil
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		} else {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			buf = (*stackWorkBuf)(unsafe.Pointer(getempty()))
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		buf.nobj = 0
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		buf.next = *head
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		*head = buf
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	buf.obj[buf.nobj] = p
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	buf.nobj++
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// Remove and return a potential pointer to a stack object.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">// Returns 0 if there are no more pointers available.</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">// This prefers non-conservative pointers so we scan stack objects</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">// precisely if there are any non-conservative pointers to them.</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>func (s *stackScanState) getPtr() (p uintptr, conservative bool) {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	for _, head := range []**stackWorkBuf{&amp;s.buf, &amp;s.cbuf} {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		buf := *head
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		if buf == nil {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			<span class="comment">// Never had any data.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			continue
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		if buf.nobj == 0 {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			if s.freeBuf != nil {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>				<span class="comment">// Free old freeBuf.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>				putempty((*workbuf)(unsafe.Pointer(s.freeBuf)))
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			<span class="comment">// Move buf to the freeBuf.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			s.freeBuf = buf
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			buf = buf.next
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			*head = buf
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			if buf == nil {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>				<span class="comment">// No more data in this list.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				continue
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		buf.nobj--
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		return buf.obj[buf.nobj], head == &amp;s.cbuf
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// No more data in either list.</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	if s.freeBuf != nil {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		putempty((*workbuf)(unsafe.Pointer(s.freeBuf)))
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		s.freeBuf = nil
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	return 0, false
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">// addObject adds a stack object at addr of type typ to the set of stack objects.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>func (s *stackScanState) addObject(addr uintptr, r *stackObjectRecord) {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	x := s.tail
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	if x == nil {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		<span class="comment">// initial setup</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		x = (*stackObjectBuf)(unsafe.Pointer(getempty()))
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		x.next = nil
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		s.head = x
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		s.tail = x
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	if x.nobj &gt; 0 &amp;&amp; uint32(addr-s.stack.lo) &lt; x.obj[x.nobj-1].off+x.obj[x.nobj-1].size {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		throw(&#34;objects added out of order or overlapping&#34;)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	if x.nobj == len(x.obj) {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		<span class="comment">// full buffer - allocate a new buffer, add to end of linked list</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		y := (*stackObjectBuf)(unsafe.Pointer(getempty()))
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		y.next = nil
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		x.next = y
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		s.tail = y
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		x = y
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	obj := &amp;x.obj[x.nobj]
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	x.nobj++
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	obj.off = uint32(addr - s.stack.lo)
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	obj.size = uint32(r.size)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	obj.setRecord(r)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// obj.left and obj.right will be initialized by buildIndex before use.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	s.nobjs++
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span><span class="comment">// buildIndex initializes s.root to a binary search tree.</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">// It should be called after all addObject calls but before</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">// any call of findObject.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>func (s *stackScanState) buildIndex() {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	s.root, _, _ = binarySearchTree(s.head, 0, s.nobjs)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">// Build a binary search tree with the n objects in the list</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">// x.obj[idx], x.obj[idx+1], ..., x.next.obj[0], ...</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// Returns the root of that tree, and the buf+idx of the nth object after x.obj[idx].</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// (The first object that was not included in the binary search tree.)</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// If n == 0, returns nil, x.</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>func binarySearchTree(x *stackObjectBuf, idx int, n int) (root *stackObject, restBuf *stackObjectBuf, restIdx int) {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	if n == 0 {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		return nil, x, idx
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	var left, right *stackObject
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	left, x, idx = binarySearchTree(x, idx, n/2)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	root = &amp;x.obj[idx]
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	idx++
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	if idx == len(x.obj) {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		x = x.next
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		idx = 0
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	right, x, idx = binarySearchTree(x, idx, n-n/2-1)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	root.left = left
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	root.right = right
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	return root, x, idx
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// findObject returns the stack object containing address a, if any.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">// Must have called buildIndex previously.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>func (s *stackScanState) findObject(a uintptr) *stackObject {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	off := uint32(a - s.stack.lo)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	obj := s.root
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	for {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		if obj == nil {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			return nil
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		if off &lt; obj.off {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>			obj = obj.left
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			continue
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		if off &gt;= obj.off+obj.size {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			obj = obj.right
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			continue
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		return obj
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
</pre><p><a href="mgcstack.go?m=text">View as plain text</a></p>

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
