<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/panic.go - Go Documentation Server</title>

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
<a href="panic.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">panic.go</span>
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
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// throwType indicates the current type of ongoing throw, which affects the</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// amount of detail printed to stderr. Higher values include more detail.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>type throwType uint32
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>const (
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// throwTypeNone means that we are not throwing.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	throwTypeNone throwType = iota
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// throwTypeUser is a throw due to a problem with the application.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// These throws do not include runtime frames, system goroutines, or</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// frame metadata.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	throwTypeUser
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// throwTypeRuntime is a throw due to a problem with Go itself.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// These throws include as much information as possible to aid in</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// debugging the runtime, including runtime frames, system goroutines,</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// and frame metadata.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	throwTypeRuntime
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// We have two different ways of doing defers. The older way involves creating a</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// defer record at the time that a defer statement is executing and adding it to a</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// defer chain. This chain is inspected by the deferreturn call at all function</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// exits in order to run the appropriate defer calls. A cheaper way (which we call</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// open-coded defers) is used for functions in which no defer statements occur in</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// loops. In that case, we simply store the defer function/arg information into</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// specific stack slots at the point of each defer statement, as well as setting a</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// bit in a bitmask. At each function exit, we add inline code to directly make</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// the appropriate defer calls based on the bitmask and fn/arg information stored</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// on the stack. During panic/Goexit processing, the appropriate defer calls are</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// made using extra funcdata info that indicates the exact stack slots that</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// contain the bitmask and defer fn/args.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// Check to make sure we can really generate a panic. If the panic</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// was generated from the runtime, or from inside malloc, then convert</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// to a throw of msg.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// pc should be the program counter of the compiler-generated code that</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// triggered this panic.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>func panicCheck1(pc uintptr, msg string) {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	if goarch.IsWasm == 0 &amp;&amp; hasPrefix(funcname(findfunc(pc)), &#34;runtime.&#34;) {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		<span class="comment">// Note: wasm can&#39;t tail call, so we can&#39;t get the original caller&#39;s pc.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		throw(msg)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// TODO: is this redundant? How could we be in malloc</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// but not in the runtime? runtime/internal/*, maybe?</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	gp := getg()
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if gp != nil &amp;&amp; gp.m != nil &amp;&amp; gp.m.mallocing != 0 {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		throw(msg)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// Same as above, but calling from the runtime is allowed.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// Using this function is necessary for any panic that may be</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// generated by runtime.sigpanic, since those are always called by the</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// runtime.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>func panicCheck2(err string) {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// panic allocates, so to avoid recursive malloc, turn panics</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// during malloc into throws.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	gp := getg()
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	if gp != nil &amp;&amp; gp.m != nil &amp;&amp; gp.m.mallocing != 0 {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		throw(err)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// Many of the following panic entry-points turn into throws when they</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// happen in various runtime contexts. These should never happen in</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// the runtime, and if they do, they indicate a serious issue and</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// should not be caught by user code.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// The panic{Index,Slice,divide,shift} functions are called by</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// code generated by the compiler for out of bounds index expressions,</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// out of bounds slice expressions, division by zero, and shift by negative.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// The panicdivide (again), panicoverflow, panicfloat, and panicmem</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// functions are called by the signal handler when a signal occurs</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// indicating the respective problem.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// Since panic{Index,Slice,shift} are never called directly, and</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// since the runtime package should never have an out of bounds slice</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// or array reference or negative shift, if we see those functions called from the</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// runtime package we turn the panic into a throw. That will dump the</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// entire runtime stack for easier debugging.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// The entry points called by the signal handler will be called from</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// runtime.sigpanic, so we can&#39;t disallow calls from the runtime to</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// these (they always look like they&#39;re called from the runtime).</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// Hence, for these, we just check for clearly bad runtime conditions.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// The panic{Index,Slice} functions are implemented in assembly and tail call</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// to the goPanic{Index,Slice} functions below. This is done so we can use</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// a space-minimal register calling convention.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// failures in the comparisons for s[x], 0 &lt;= x &lt; y (y == len(s))</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//go:yeswritebarrierrec</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>func goPanicIndex(x int, y int) {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;index out of range&#34;)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsIndex})
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">//go:yeswritebarrierrec</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>func goPanicIndexU(x uint, y int) {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;index out of range&#34;)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsIndex})
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// failures in the comparisons for s[:x], 0 &lt;= x &lt;= y (y == len(s) or cap(s))</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">//go:yeswritebarrierrec</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>func goPanicSliceAlen(x int, y int) {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSliceAlen})
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">//go:yeswritebarrierrec</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>func goPanicSliceAlenU(x uint, y int) {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSliceAlen})
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">//go:yeswritebarrierrec</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>func goPanicSliceAcap(x int, y int) {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSliceAcap})
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">//go:yeswritebarrierrec</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>func goPanicSliceAcapU(x uint, y int) {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSliceAcap})
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// failures in the comparisons for s[x:y], 0 &lt;= x &lt;= y</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">//go:yeswritebarrierrec</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func goPanicSliceB(x int, y int) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSliceB})
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">//go:yeswritebarrierrec</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>func goPanicSliceBU(x uint, y int) {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSliceB})
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// failures in the comparisons for s[::x], 0 &lt;= x &lt;= y (y == len(s) or cap(s))</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>func goPanicSlice3Alen(x int, y int) {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSlice3Alen})
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>func goPanicSlice3AlenU(x uint, y int) {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSlice3Alen})
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>func goPanicSlice3Acap(x int, y int) {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSlice3Acap})
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>func goPanicSlice3AcapU(x uint, y int) {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSlice3Acap})
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">// failures in the comparisons for s[:x:y], 0 &lt;= x &lt;= y</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>func goPanicSlice3B(x int, y int) {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSlice3B})
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>func goPanicSlice3BU(x uint, y int) {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSlice3B})
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// failures in the comparisons for s[x:y:], 0 &lt;= x &lt;= y</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>func goPanicSlice3C(x int, y int) {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSlice3C})
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>func goPanicSlice3CU(x uint, y int) {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice bounds out of range&#34;)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSlice3C})
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// failures in the conversion ([x]T)(s) or (*[x]T)(s), 0 &lt;= x &lt;= y, y == len(s)</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>func goPanicSliceConvert(x int, y int) {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;slice length too short to convert to array or pointer to array&#34;)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsConvert})
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// Implemented in assembly, as they take arguments in registers.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// Declared here to mark them as ABIInternal.</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>func panicIndex(x int, y int)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>func panicIndexU(x uint, y int)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>func panicSliceAlen(x int, y int)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>func panicSliceAlenU(x uint, y int)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>func panicSliceAcap(x int, y int)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>func panicSliceAcapU(x uint, y int)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>func panicSliceB(x int, y int)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>func panicSliceBU(x uint, y int)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>func panicSlice3Alen(x int, y int)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>func panicSlice3AlenU(x uint, y int)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>func panicSlice3Acap(x int, y int)
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>func panicSlice3AcapU(x uint, y int)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>func panicSlice3B(x int, y int)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>func panicSlice3BU(x uint, y int)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>func panicSlice3C(x int, y int)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>func panicSlice3CU(x uint, y int)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>func panicSliceConvert(x int, y int)
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>var shiftError = error(errorString(&#34;negative shift amount&#34;))
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">//go:yeswritebarrierrec</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>func panicshift() {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	panicCheck1(getcallerpc(), &#34;negative shift amount&#34;)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	panic(shiftError)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>var divideError = error(errorString(&#34;integer divide by zero&#34;))
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">//go:yeswritebarrierrec</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>func panicdivide() {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	panicCheck2(&#34;integer divide by zero&#34;)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	panic(divideError)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>var overflowError = error(errorString(&#34;integer overflow&#34;))
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>func panicoverflow() {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	panicCheck2(&#34;integer overflow&#34;)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	panic(overflowError)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>var floatError = error(errorString(&#34;floating point error&#34;))
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>func panicfloat() {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	panicCheck2(&#34;floating point error&#34;)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	panic(floatError)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>var memoryError = error(errorString(&#34;invalid memory address or nil pointer dereference&#34;))
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>func panicmem() {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	panicCheck2(&#34;invalid memory address or nil pointer dereference&#34;)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	panic(memoryError)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>func panicmemAddr(addr uintptr) {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	panicCheck2(&#34;invalid memory address or nil pointer dereference&#34;)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	panic(errorAddressString{msg: &#34;invalid memory address or nil pointer dereference&#34;, addr: addr})
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">// Create a new deferred function fn, which has no arguments and results.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span><span class="comment">// The compiler turns a defer statement into a call to this.</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>func deferproc(fn func()) {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	gp := getg()
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	if gp.m.curg != gp {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		<span class="comment">// go code on the system stack can&#39;t defer</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		throw(&#34;defer on system stack&#34;)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	d := newdefer()
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	d.link = gp._defer
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	gp._defer = d
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	d.fn = fn
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	d.pc = getcallerpc()
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	<span class="comment">// We must not be preempted between calling getcallersp and</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	<span class="comment">// storing it to d.sp because getcallersp&#39;s result is a</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">// uintptr stack pointer.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	d.sp = getcallersp()
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	<span class="comment">// deferproc returns 0 normally.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	<span class="comment">// a deferred func that stops a panic</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">// makes the deferproc return 1.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// the code the compiler generates always</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// checks the return value and jumps to the</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// end of the function if deferproc returns != 0.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	return0()
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// No code can go here - the C return register has</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">// been set and must not be clobbered.</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>var rangeExitError = error(errorString(&#34;range function continued iteration after exit&#34;))
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">//go:noinline</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>func panicrangeexit() {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	panic(rangeExitError)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">// deferrangefunc is called by functions that are about to</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">// execute a range-over-function loop in which the loop body</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// may execute a defer statement. That defer needs to add to</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// the chain for the current function, not the func literal synthesized</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// to represent the loop body. To do that, the original function</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// calls deferrangefunc to obtain an opaque token representing</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">// the current frame, and then the loop body uses deferprocat</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span><span class="comment">// instead of deferproc to add to that frame&#39;s defer lists.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// The token is an &#39;any&#39; with underlying type *atomic.Pointer[_defer].</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">// It is the atomically-updated head of a linked list of _defer structs</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">// representing deferred calls. At the same time, we create a _defer</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">// struct on the main g._defer list with d.head set to this head pointer.</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span><span class="comment">// The g._defer list is now a linked list of deferred calls,</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span><span class="comment">// but an atomic list hanging off:</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span><span class="comment">//		g._defer =&gt; d4 -&gt; d3 -&gt; drangefunc -&gt; d2 -&gt; d1 -&gt; nil</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span><span class="comment">//	                             | .head</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span><span class="comment">//	                             |</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">//	                             +--&gt; dY -&gt; dX -&gt; nil</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// with each -&gt; indicating a d.link pointer, and where drangefunc</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// has the d.rangefunc = true bit set.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">// Note that the function being ranged over may have added</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span><span class="comment">// its own defers (d4 and d3), so drangefunc need not be at the</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// top of the list when deferprocat is used. This is why we pass</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">// the atomic head explicitly.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// To keep misbehaving programs from crashing the runtime,</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">// deferprocat pushes new defers onto the .head list atomically.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// The fact that it is a separate list from the main goroutine</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">// defer list means that the main goroutine&#39;s defers can still</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// be handled non-atomically.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">// In the diagram, dY and dX are meant to be processed when</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span><span class="comment">// drangefunc would be processed, which is to say the defer order</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">// should be d4, d3, dY, dX, d2, d1. To make that happen,</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">// when defer processing reaches a d with rangefunc=true,</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// it calls deferconvert to atomically take the extras</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// away from d.head and then adds them to the main list.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// That is, deferconvert changes this list:</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">//		g._defer =&gt; drangefunc -&gt; d2 -&gt; d1 -&gt; nil</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">//	                 | .head</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">//	                 |</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">//	                 +--&gt; dY -&gt; dX -&gt; nil</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// into this list:</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">//	g._defer =&gt; dY -&gt; dX -&gt; d2 -&gt; d1 -&gt; nil</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">// It also poisons *drangefunc.head so that any future</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span><span class="comment">// deferprocat using that head will throw.</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span><span class="comment">// (The atomic head is ordinary garbage collected memory so that</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// it&#39;s not a problem if user code holds onto it beyond</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span><span class="comment">// the lifetime of drangefunc.)</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// TODO: We could arrange for the compiler to call into the</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// runtime after the loop finishes normally, to do an eager</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">// deferconvert, which would catch calling the loop body</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">// and having it defer after the loop is done. If we have a</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span><span class="comment">// more general catch of loop body misuse, though, this</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">// might not be worth worrying about in addition.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span><span class="comment">// See also ../cmd/compile/internal/rangefunc/rewrite.go.</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>func deferrangefunc() any {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	gp := getg()
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	if gp.m.curg != gp {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		<span class="comment">// go code on the system stack can&#39;t defer</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		throw(&#34;defer on system stack&#34;)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	d := newdefer()
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	d.link = gp._defer
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	gp._defer = d
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	d.pc = getcallerpc()
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// We must not be preempted between calling getcallersp and</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	<span class="comment">// storing it to d.sp because getcallersp&#39;s result is a</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	<span class="comment">// uintptr stack pointer.</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	d.sp = getcallersp()
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	d.rangefunc = true
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	d.head = new(atomic.Pointer[_defer])
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	return d.head
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// badDefer returns a fixed bad defer pointer for poisoning an atomic defer list head.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>func badDefer() *_defer {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	return (*_defer)(unsafe.Pointer(uintptr(1)))
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span><span class="comment">// deferprocat is like deferproc but adds to the atomic list represented by frame.</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span><span class="comment">// See the doc comment for deferrangefunc for details.</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>func deferprocat(fn func(), frame any) {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	head := frame.(*atomic.Pointer[_defer])
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	if raceenabled {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		racewritepc(unsafe.Pointer(head), getcallerpc(), abi.FuncPCABIInternal(deferprocat))
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	d1 := newdefer()
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	d1.fn = fn
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	for {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		d1.link = head.Load()
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		if d1.link == badDefer() {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			throw(&#34;defer after range func returned&#34;)
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		if head.CompareAndSwap(d1.link, d1) {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			break
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	<span class="comment">// Must be last - see deferproc above.</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	return0()
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span><span class="comment">// deferconvert converts a rangefunc defer list into an ordinary list.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span><span class="comment">// See the doc comment for deferrangefunc for details.</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>func deferconvert(d *_defer) *_defer {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	head := d.head
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	if raceenabled {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		racereadpc(unsafe.Pointer(head), getcallerpc(), abi.FuncPCABIInternal(deferconvert))
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	tail := d.link
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	d.rangefunc = false
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	d0 := d
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	for {
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		d = head.Load()
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		if head.CompareAndSwap(d, badDefer()) {
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			break
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	}
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	if d == nil {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		freedefer(d0)
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		return tail
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	for d1 := d; ; d1 = d1.link {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		d1.sp = d0.sp
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		d1.pc = d0.pc
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		if d1.link == nil {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>			d1.link = tail
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>			break
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	freedefer(d0)
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	return d
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span><span class="comment">// deferprocStack queues a new deferred function with a defer record on the stack.</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span><span class="comment">// The defer record must have its fn field initialized.</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span><span class="comment">// All other fields can contain junk.</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span><span class="comment">// Nosplit because of the uninitialized pointer fields on the stack.</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>func deferprocStack(d *_defer) {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	gp := getg()
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	if gp.m.curg != gp {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		<span class="comment">// go code on the system stack can&#39;t defer</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		throw(&#34;defer on system stack&#34;)
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	<span class="comment">// fn is already set.</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	<span class="comment">// The other fields are junk on entry to deferprocStack and</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	<span class="comment">// are initialized here.</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	d.heap = false
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	d.rangefunc = false
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	d.sp = getcallersp()
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	d.pc = getcallerpc()
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	<span class="comment">// The lines below implement:</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	<span class="comment">//   d.panic = nil</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	<span class="comment">//   d.fd = nil</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	<span class="comment">//   d.link = gp._defer</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	<span class="comment">//   d.head = nil</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	<span class="comment">//   gp._defer = d</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	<span class="comment">// But without write barriers. The first three are writes to</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	<span class="comment">// the stack so they don&#39;t need a write barrier, and furthermore</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	<span class="comment">// are to uninitialized memory, so they must not use a write barrier.</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	<span class="comment">// The fourth write does not require a write barrier because we</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	<span class="comment">// explicitly mark all the defer structures, so we don&#39;t need to</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	<span class="comment">// keep track of pointers to them with a write barrier.</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	*(*uintptr)(unsafe.Pointer(&amp;d.link)) = uintptr(unsafe.Pointer(gp._defer))
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	*(*uintptr)(unsafe.Pointer(&amp;d.head)) = 0
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	*(*uintptr)(unsafe.Pointer(&amp;gp._defer)) = uintptr(unsafe.Pointer(d))
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	return0()
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	<span class="comment">// No code can go here - the C return register has</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	<span class="comment">// been set and must not be clobbered.</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>}
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span><span class="comment">// Each P holds a pool for defers.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span><span class="comment">// Allocate a Defer, usually using per-P pool.</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span><span class="comment">// Each defer must be released with freedefer.  The defer is not</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span><span class="comment">// added to any defer chain yet.</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>func newdefer() *_defer {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	var d *_defer
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	pp := mp.p.ptr()
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	if len(pp.deferpool) == 0 &amp;&amp; sched.deferpool != nil {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		lock(&amp;sched.deferlock)
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		for len(pp.deferpool) &lt; cap(pp.deferpool)/2 &amp;&amp; sched.deferpool != nil {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			d := sched.deferpool
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			sched.deferpool = d.link
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			d.link = nil
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			pp.deferpool = append(pp.deferpool, d)
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		unlock(&amp;sched.deferlock)
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	if n := len(pp.deferpool); n &gt; 0 {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		d = pp.deferpool[n-1]
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		pp.deferpool[n-1] = nil
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		pp.deferpool = pp.deferpool[:n-1]
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	}
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	releasem(mp)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	mp, pp = nil, nil
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	if d == nil {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		<span class="comment">// Allocate new defer.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		d = new(_defer)
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	d.heap = true
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	return d
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>}
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span><span class="comment">// Free the given defer.</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span><span class="comment">// The defer cannot be used after this call.</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span><span class="comment">// This is nosplit because the incoming defer is in a perilous state.</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span><span class="comment">// It&#39;s not on any defer list, so stack copying won&#39;t adjust stack</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span><span class="comment">// pointers in it (namely, d.link). Hence, if we were to copy the</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span><span class="comment">// stack, d could then contain a stale pointer.</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>func freedefer(d *_defer) {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	d.link = nil
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	<span class="comment">// After this point we can copy the stack.</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	if d.fn != nil {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		freedeferfn()
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	}
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	if !d.heap {
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		return
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	}
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	pp := mp.p.ptr()
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	if len(pp.deferpool) == cap(pp.deferpool) {
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		<span class="comment">// Transfer half of local cache to the central cache.</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		var first, last *_defer
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		for len(pp.deferpool) &gt; cap(pp.deferpool)/2 {
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>			n := len(pp.deferpool)
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			d := pp.deferpool[n-1]
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>			pp.deferpool[n-1] = nil
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			pp.deferpool = pp.deferpool[:n-1]
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>			if first == nil {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>				first = d
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>			} else {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>				last.link = d
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>			}
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>			last = d
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		lock(&amp;sched.deferlock)
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		last.link = sched.deferpool
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		sched.deferpool = first
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		unlock(&amp;sched.deferlock)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	}
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	*d = _defer{}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	pp.deferpool = append(pp.deferpool, d)
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	releasem(mp)
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	mp, pp = nil, nil
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span><span class="comment">// Separate function so that it can split stack.</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span><span class="comment">// Windows otherwise runs out of stack space.</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>func freedeferfn() {
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	<span class="comment">// fn must be cleared before d is unlinked from gp.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	throw(&#34;freedefer with d.fn != nil&#34;)
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span><span class="comment">// deferreturn runs deferred functions for the caller&#39;s frame.</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span><span class="comment">// The compiler inserts a call to this at the end of any</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span><span class="comment">// function which calls defer.</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>func deferreturn() {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	var p _panic
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	p.deferreturn = true
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	p.start(getcallerpc(), unsafe.Pointer(getcallersp()))
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	for {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		fn, ok := p.nextDefer()
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		if !ok {
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>			break
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		fn()
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	}
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>}
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span><span class="comment">// Goexit terminates the goroutine that calls it. No other goroutine is affected.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span><span class="comment">// Goexit runs all deferred calls before terminating the goroutine. Because Goexit</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span><span class="comment">// is not a panic, any recover calls in those deferred functions will return nil.</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span><span class="comment">// Calling Goexit from the main goroutine terminates that goroutine</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span><span class="comment">// without func main returning. Since func main has not returned,</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span><span class="comment">// the program continues execution of other goroutines.</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span><span class="comment">// If all other goroutines exit, the program crashes.</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>func Goexit() {
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	<span class="comment">// Create a panic object for Goexit, so we can recognize when it might be</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	<span class="comment">// bypassed by a recover().</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	var p _panic
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	p.goexit = true
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	p.start(getcallerpc(), unsafe.Pointer(getcallersp()))
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	for {
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		fn, ok := p.nextDefer()
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		if !ok {
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>			break
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		}
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		fn()
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	goexit1()
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span><span class="comment">// Call all Error and String methods before freezing the world.</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span><span class="comment">// Used when crashing with panicking.</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>func preprintpanics(p *_panic) {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	defer func() {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		text := &#34;panic while printing panic value&#34;
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		switch r := recover().(type) {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		case nil:
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>			<span class="comment">// nothing to do</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		case string:
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>			throw(text + &#34;: &#34; + r)
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>		default:
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>			throw(text + &#34;: type &#34; + toRType(efaceOf(&amp;r)._type).string())
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		}
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	}()
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	for p != nil {
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		switch v := p.arg.(type) {
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		case error:
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>			p.arg = v.Error()
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		case stringer:
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>			p.arg = v.String()
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		}
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>		p = p.link
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>}
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span><span class="comment">// Print all currently active panics. Used when crashing.</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span><span class="comment">// Should only be called after preprintpanics.</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>func printpanics(p *_panic) {
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	if p.link != nil {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		printpanics(p.link)
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		if !p.link.goexit {
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>			print(&#34;\t&#34;)
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		}
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	}
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	if p.goexit {
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		return
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	}
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	print(&#34;panic: &#34;)
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	printany(p.arg)
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	if p.recovered {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		print(&#34; [recovered]&#34;)
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	print(&#34;\n&#34;)
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>}
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span><span class="comment">// readvarintUnsafe reads the uint32 in varint format starting at fd, and returns the</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span><span class="comment">// uint32 and a pointer to the byte following the varint.</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span><span class="comment">// The implementation is the same with runtime.readvarint, except that this function</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span><span class="comment">// uses unsafe.Pointer for speed.</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>func readvarintUnsafe(fd unsafe.Pointer) (uint32, unsafe.Pointer) {
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	var r uint32
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	var shift int
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	for {
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		b := *(*uint8)(fd)
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		fd = add(fd, unsafe.Sizeof(b))
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		if b &lt; 128 {
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>			return r + uint32(b)&lt;&lt;shift, fd
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		}
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		r += uint32(b&amp;0x7F) &lt;&lt; (shift &amp; 31)
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		shift += 7
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		if shift &gt; 28 {
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>			panic(&#34;Bad varint&#34;)
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>		}
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	}
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span><span class="comment">// A PanicNilError happens when code calls panic(nil).</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span><span class="comment">// Before Go 1.21, programs that called panic(nil) observed recover returning nil.</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span><span class="comment">// Starting in Go 1.21, programs that call panic(nil) observe recover returning a *PanicNilError.</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span><span class="comment">// Programs can change back to the old behavior by setting GODEBUG=panicnil=1.</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>type PanicNilError struct {
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	<span class="comment">// This field makes PanicNilError structurally different from</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	<span class="comment">// any other struct in this package, and the _ makes it different</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	<span class="comment">// from any struct in other packages too.</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	<span class="comment">// This avoids any accidental conversions being possible</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	<span class="comment">// between this struct and some other struct sharing the same fields,</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	<span class="comment">// like happened in go.dev/issue/56603.</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	_ [0]*PanicNilError
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>}
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>func (*PanicNilError) Error() string { return &#34;panic called with nil argument&#34; }
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>func (*PanicNilError) RuntimeError() {}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>var panicnil = &amp;godebugInc{name: &#34;panicnil&#34;}
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span><span class="comment">// The implementation of the predeclared function panic.</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>func gopanic(e any) {
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	if e == nil {
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		if debug.panicnil.Load() != 1 {
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>			e = new(PanicNilError)
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		} else {
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>			panicnil.IncNonDefault()
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		}
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>	}
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	gp := getg()
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	if gp.m.curg != gp {
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>		print(&#34;panic: &#34;)
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		printany(e)
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		throw(&#34;panic on system stack&#34;)
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	}
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	if gp.m.mallocing != 0 {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		print(&#34;panic: &#34;)
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		printany(e)
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		throw(&#34;panic during malloc&#34;)
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	}
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	if gp.m.preemptoff != &#34;&#34; {
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		print(&#34;panic: &#34;)
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		printany(e)
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		print(&#34;preempt off reason: &#34;)
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		print(gp.m.preemptoff)
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		throw(&#34;panic during preemptoff&#34;)
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	}
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	if gp.m.locks != 0 {
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>		print(&#34;panic: &#34;)
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		printany(e)
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>		print(&#34;\n&#34;)
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		throw(&#34;panic holding locks&#34;)
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	}
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	var p _panic
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	p.arg = e
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	runningPanicDefers.Add(1)
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	p.start(getcallerpc(), unsafe.Pointer(getcallersp()))
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	for {
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		fn, ok := p.nextDefer()
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>		if !ok {
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>			break
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		}
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		fn()
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	}
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	<span class="comment">// ran out of deferred calls - old-school panic now</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	<span class="comment">// Because it is unsafe to call arbitrary user code after freezing</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>	<span class="comment">// the world, we call preprintpanics to invoke all necessary Error</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	<span class="comment">// and String methods to prepare the panic strings before startpanic.</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>	preprintpanics(&amp;p)
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>	fatalpanic(&amp;p)   <span class="comment">// should not return</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	*(*int)(nil) = 0 <span class="comment">// not reached</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>}
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span><span class="comment">// start initializes a panic to start unwinding the stack.</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span><span class="comment">// If p.goexit is true, then start may return multiple times.</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>func (p *_panic) start(pc uintptr, sp unsafe.Pointer) {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	gp := getg()
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	<span class="comment">// Record the caller&#39;s PC and SP, so recovery can identify panics</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	<span class="comment">// that have been recovered. Also, so that if p is from Goexit, we</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	<span class="comment">// can restart its defer processing loop if a recovered panic tries</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	<span class="comment">// to jump past it.</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	p.startPC = getcallerpc()
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	p.startSP = unsafe.Pointer(getcallersp())
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	if p.deferreturn {
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>		p.sp = sp
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		if s := (*savedOpenDeferState)(gp.param); s != nil {
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			<span class="comment">// recovery saved some state for us, so that we can resume</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>			<span class="comment">// calling open-coded defers without unwinding the stack.</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>			gp.param = nil
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>			p.retpc = s.retpc
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>			p.deferBitsPtr = (*byte)(add(sp, s.deferBitsOffset))
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>			p.slotsPtr = add(sp, s.slotsOffset)
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		}
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>		return
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	}
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	p.link = gp._panic
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	gp._panic = (*_panic)(noescape(unsafe.Pointer(p)))
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	<span class="comment">// Initialize state machine, and find the first frame with a defer.</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>	<span class="comment">// Note: We could use startPC and startSP here, but callers will</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	<span class="comment">// never have defer statements themselves. By starting at their</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	<span class="comment">// caller instead, we avoid needing to unwind through an extra</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	<span class="comment">// frame. It also somewhat simplifies the terminating condition for</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	<span class="comment">// deferreturn.</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	p.lr, p.fp = pc, sp
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	p.nextFrame()
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>}
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span><span class="comment">// nextDefer returns the next deferred function to invoke, if any.</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span><span class="comment">// Note: The &#34;ok bool&#34; result is necessary to correctly handle when</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span><span class="comment">// the deferred function itself was nil (e.g., &#34;defer (func())(nil)&#34;).</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>func (p *_panic) nextDefer() (func(), bool) {
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	gp := getg()
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	if !p.deferreturn {
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		if gp._panic != p {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>			throw(&#34;bad panic stack&#34;)
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		}
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		if p.recovered {
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>			mcall(recovery) <span class="comment">// does not return</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>			throw(&#34;recovery failed&#34;)
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	<span class="comment">// The assembler adjusts p.argp in wrapper functions that shouldn&#39;t</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>	<span class="comment">// be visible to recover(), so we need to restore it each iteration.</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	p.argp = add(p.startSP, sys.MinFrameSize)
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	for {
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		for p.deferBitsPtr != nil {
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>			bits := *p.deferBitsPtr
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>			<span class="comment">// Check whether any open-coded defers are still pending.</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>			<span class="comment">// Note: We need to check this upfront (rather than after</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>			<span class="comment">// clearing the top bit) because it&#39;s possible that Goexit</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>			<span class="comment">// invokes a deferred call, and there were still more pending</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>			<span class="comment">// open-coded defers in the frame; but then the deferred call</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>			<span class="comment">// panic and invoked the remaining defers in the frame, before</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>			<span class="comment">// recovering and restarting the Goexit loop.</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>			if bits == 0 {
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>				p.deferBitsPtr = nil
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>				break
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>			}
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>			<span class="comment">// Find index of top bit set.</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>			i := 7 - uintptr(sys.LeadingZeros8(bits))
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>			<span class="comment">// Clear bit and store it back.</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>			bits &amp;^= 1 &lt;&lt; i
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>			*p.deferBitsPtr = bits
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>			return *(*func())(add(p.slotsPtr, i*goarch.PtrSize)), true
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		}
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	Recheck:
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>		if d := gp._defer; d != nil &amp;&amp; d.sp == uintptr(p.sp) {
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>			if d.rangefunc {
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>				gp._defer = deferconvert(d)
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>				goto Recheck
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>			}
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>			fn := d.fn
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>			d.fn = nil
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>			<span class="comment">// TODO(mdempsky): Instead of having each deferproc call have</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>			<span class="comment">// its own &#34;deferreturn(); return&#34; sequence, we should just make</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>			<span class="comment">// them reuse the one we emit for open-coded defers.</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			p.retpc = d.pc
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			<span class="comment">// Unlink and free.</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>			gp._defer = d.link
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>			freedefer(d)
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>			return fn, true
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>		}
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		if !p.nextFrame() {
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>			return nil, false
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		}
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	}
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>}
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span><span class="comment">// nextFrame finds the next frame that contains deferred calls, if any.</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>func (p *_panic) nextFrame() (ok bool) {
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	if p.lr == 0 {
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>		return false
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	}
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	gp := getg()
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>		var limit uintptr
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>		if d := gp._defer; d != nil {
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>			limit = d.sp
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		}
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		var u unwinder
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		u.initAt(p.lr, uintptr(p.fp), 0, gp, 0)
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>		for {
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>			if !u.valid() {
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>				p.lr = 0
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>				return <span class="comment">// ok == false</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>			}
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>			<span class="comment">// TODO(mdempsky): If we populate u.frame.fn.deferreturn for</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>			<span class="comment">// every frame containing a defer (not just open-coded defers),</span>
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>			<span class="comment">// then we can simply loop until we find the next frame where</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>			<span class="comment">// it&#39;s non-zero.</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>			if u.frame.sp == limit {
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>				break <span class="comment">// found a frame with linked defers</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>			}
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>			if p.initOpenCodedDefers(u.frame.fn, unsafe.Pointer(u.frame.varp)) {
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>				break <span class="comment">// found a frame with open-coded defers</span>
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>			}
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>			u.next()
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>		}
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		p.lr = u.frame.lr
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		p.sp = unsafe.Pointer(u.frame.sp)
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>		p.fp = unsafe.Pointer(u.frame.fp)
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>		ok = true
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>	})
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>	return
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>}
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>func (p *_panic) initOpenCodedDefers(fn funcInfo, varp unsafe.Pointer) bool {
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>	fd := funcdata(fn, abi.FUNCDATA_OpenCodedDeferInfo)
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>	if fd == nil {
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>		return false
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>	}
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	if fn.deferreturn == 0 {
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>		throw(&#34;missing deferreturn&#34;)
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>	}
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	deferBitsOffset, fd := readvarintUnsafe(fd)
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>	deferBitsPtr := (*uint8)(add(varp, -uintptr(deferBitsOffset)))
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>	if *deferBitsPtr == 0 {
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>		return false <span class="comment">// has open-coded defers, but none pending</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	}
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>	slotsOffset, fd := readvarintUnsafe(fd)
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>	p.retpc = fn.entry() + uintptr(fn.deferreturn)
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	p.deferBitsPtr = deferBitsPtr
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>	p.slotsPtr = add(varp, -uintptr(slotsOffset))
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	return true
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>}
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span><span class="comment">// The implementation of the predeclared function recover.</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span><span class="comment">// Cannot split the stack because it needs to reliably</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span><span class="comment">// find the stack segment of its caller.</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span><span class="comment">// TODO(rsc): Once we commit to CopyStackAlways,</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span><span class="comment">// this doesn&#39;t need to be nosplit.</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>func gorecover(argp uintptr) any {
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	<span class="comment">// Must be in a function running as part of a deferred call during the panic.</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>	<span class="comment">// Must be called from the topmost function of the call</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	<span class="comment">// (the function used in the defer statement).</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>	<span class="comment">// p.argp is the argument pointer of that topmost deferred function call.</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	<span class="comment">// Compare against argp reported by caller.</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	<span class="comment">// If they match, the caller is the one who can recover.</span>
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>	gp := getg()
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>	p := gp._panic
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	if p != nil &amp;&amp; !p.goexit &amp;&amp; !p.recovered &amp;&amp; argp == uintptr(p.argp) {
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>		p.recovered = true
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		return p.arg
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	}
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	return nil
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>}
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span><span class="comment">//go:linkname sync_throw sync.throw</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>func sync_throw(s string) {
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>	throw(s)
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>}
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span><span class="comment">//go:linkname sync_fatal sync.fatal</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>func sync_fatal(s string) {
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	fatal(s)
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>}
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span><span class="comment">// throw triggers a fatal error that dumps a stack trace and exits.</span>
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span><span class="comment">// throw should be used for runtime-internal fatal errors where Go itself,</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span><span class="comment">// rather than user code, may be at fault for the failure.</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>func throw(s string) {
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	<span class="comment">// Everything throw does should be recursively nosplit so it</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	<span class="comment">// can be called even when it&#39;s unsafe to grow the stack.</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>		print(&#34;fatal error: &#34;, s, &#34;\n&#34;)
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	})
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	fatalthrow(throwTypeRuntime)
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>}
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span><span class="comment">// fatal triggers a fatal error that dumps a stack trace and exits.</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span><span class="comment">// fatal is equivalent to throw, but is used when user code is expected to be</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span><span class="comment">// at fault for the failure, such as racing map writes.</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span><span class="comment">// fatal does not include runtime frames, system goroutines, or frame metadata</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span><span class="comment">// (fp, sp, pc) in the stack trace unless GOTRACEBACK=system or higher.</span>
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>func fatal(s string) {
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>	<span class="comment">// Everything fatal does should be recursively nosplit so it</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>	<span class="comment">// can be called even when it&#39;s unsafe to grow the stack.</span>
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>		print(&#34;fatal error: &#34;, s, &#34;\n&#34;)
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>	})
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>	fatalthrow(throwTypeUser)
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>}
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span><span class="comment">// runningPanicDefers is non-zero while running deferred functions for panic.</span>
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span><span class="comment">// This is used to try hard to get a panic stack trace out when exiting.</span>
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>var runningPanicDefers atomic.Uint32
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span><span class="comment">// panicking is non-zero when crashing the program for an unrecovered panic.</span>
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>var panicking atomic.Uint32
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span><span class="comment">// paniclk is held while printing the panic information and stack trace,</span>
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span><span class="comment">// so that two concurrent panics don&#39;t overlap their output.</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>var paniclk mutex
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span><span class="comment">// Unwind the stack after a deferred function calls recover</span>
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span><span class="comment">// after a panic. Then arrange to continue running as though</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span><span class="comment">// the caller of the deferred function returned normally.</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span><span class="comment">// However, if unwinding the stack would skip over a Goexit call, we</span>
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span><span class="comment">// return into the Goexit loop instead, so it can continue processing</span>
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span><span class="comment">// defers instead.</span>
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>func recovery(gp *g) {
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>	p := gp._panic
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>	pc, sp, fp := p.retpc, uintptr(p.sp), uintptr(p.fp)
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	p0, saveOpenDeferState := p, p.deferBitsPtr != nil &amp;&amp; *p.deferBitsPtr != 0
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>	<span class="comment">// Unwind the panic stack.</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>	for ; p != nil &amp;&amp; uintptr(p.startSP) &lt; sp; p = p.link {
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>		<span class="comment">// Don&#39;t allow jumping past a pending Goexit.</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>		<span class="comment">// Instead, have its _panic.start() call return again.</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>		<span class="comment">// TODO(mdempsky): In this case, Goexit will resume walking the</span>
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>		<span class="comment">// stack where it left off, which means it will need to rewalk</span>
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>		<span class="comment">// frames that we&#39;ve already processed.</span>
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>		<span class="comment">// There&#39;s a similar issue with nested panics, when the inner</span>
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>		<span class="comment">// panic supercedes the outer panic. Again, we end up needing to</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>		<span class="comment">// walk the same stack frames.</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>		<span class="comment">// These are probably pretty rare occurrences in practice, and</span>
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>		<span class="comment">// they don&#39;t seem any worse than the existing logic. But if we</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>		<span class="comment">// move the unwinding state into _panic, we could detect when we</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>		<span class="comment">// run into where the last panic started, and then just pick up</span>
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>		<span class="comment">// where it left off instead.</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>		<span class="comment">// With how subtle defer handling is, this might not actually be</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>		<span class="comment">// worthwhile though.</span>
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>		if p.goexit {
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>			pc, sp = p.startPC, uintptr(p.startSP)
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>			saveOpenDeferState = false <span class="comment">// goexit is unwinding the stack anyway</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>			break
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>		}
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>		runningPanicDefers.Add(-1)
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>	}
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>	gp._panic = p
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	if p == nil { <span class="comment">// must be done with signal</span>
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>		gp.sig = 0
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	}
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>	if gp.param != nil {
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>		throw(&#34;unexpected gp.param&#34;)
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>	}
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>	if saveOpenDeferState {
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>		<span class="comment">// If we&#39;re returning to deferreturn and there are more open-coded</span>
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>		<span class="comment">// defers for it to call, save enough state for it to be able to</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>		<span class="comment">// pick up where p0 left off.</span>
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>		gp.param = unsafe.Pointer(&amp;savedOpenDeferState{
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>			retpc: p0.retpc,
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>			<span class="comment">// We need to save deferBitsPtr and slotsPtr too, but those are</span>
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>			<span class="comment">// stack pointers. To avoid issues around heap objects pointing</span>
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>			<span class="comment">// to the stack, save them as offsets from SP.</span>
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>			deferBitsOffset: uintptr(unsafe.Pointer(p0.deferBitsPtr)) - uintptr(p0.sp),
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>			slotsOffset:     uintptr(p0.slotsPtr) - uintptr(p0.sp),
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>		})
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>	}
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>	<span class="comment">// TODO(mdempsky): Currently, we rely on frames containing &#34;defer&#34;</span>
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	<span class="comment">// to end with &#34;CALL deferreturn; RET&#34;. This allows deferreturn to</span>
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	<span class="comment">// finish running any pending defers in the frame.</span>
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>	<span class="comment">// But we should be able to tell whether there are still pending</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>	<span class="comment">// defers here. If there aren&#39;t, we can just jump directly to the</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>	<span class="comment">// &#34;RET&#34; instruction. And if there are, we don&#39;t need an actual</span>
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>	<span class="comment">// &#34;CALL deferreturn&#34; instruction; we can simulate it with something</span>
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>	<span class="comment">// like:</span>
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>	<span class="comment">//	if usesLR {</span>
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>	<span class="comment">//		lr = pc</span>
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>	<span class="comment">//	} else {</span>
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>	<span class="comment">//		sp -= sizeof(pc)</span>
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>	<span class="comment">//		*(*uintptr)(sp) = pc</span>
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>	<span class="comment">//	}</span>
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>	<span class="comment">//	pc = funcPC(deferreturn)</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>	<span class="comment">// So that we effectively tail call into deferreturn, such that it</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>	<span class="comment">// then returns to the simple &#34;RET&#34; epilogue. That would save the</span>
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	<span class="comment">// overhead of the &#34;deferreturn&#34; call when there aren&#39;t actually any</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	<span class="comment">// pending defers left, and shrink the TEXT size of compiled</span>
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	<span class="comment">// binaries. (Admittedly, both of these are modest savings.)</span>
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>	<span class="comment">// Ensure we&#39;re recovering within the appropriate stack.</span>
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>	if sp != 0 &amp;&amp; (sp &lt; gp.stack.lo || gp.stack.hi &lt; sp) {
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>		print(&#34;recover: &#34;, hex(sp), &#34; not in [&#34;, hex(gp.stack.lo), &#34;, &#34;, hex(gp.stack.hi), &#34;]\n&#34;)
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>		throw(&#34;bad recovery&#34;)
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>	}
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>	<span class="comment">// Make the deferproc for this d return again,</span>
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>	<span class="comment">// this time returning 1. The calling function will</span>
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>	<span class="comment">// jump to the standard return epilogue.</span>
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>	gp.sched.sp = sp
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>	gp.sched.pc = pc
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>	gp.sched.lr = 0
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>	<span class="comment">// Restore the bp on platforms that support frame pointers.</span>
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>	<span class="comment">// N.B. It&#39;s fine to not set anything for platforms that don&#39;t</span>
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>	<span class="comment">// support frame pointers, since nothing consumes them.</span>
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>	switch {
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>	case goarch.IsAmd64 != 0:
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>		<span class="comment">// on x86, fp actually points one word higher than the top of</span>
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>		<span class="comment">// the frame since the return address is saved on the stack by</span>
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>		<span class="comment">// the caller</span>
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>		gp.sched.bp = fp - 2*goarch.PtrSize
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>	case goarch.IsArm64 != 0:
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>		<span class="comment">// on arm64, the architectural bp points one word higher</span>
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>		<span class="comment">// than the sp. fp is totally useless to us here, because it</span>
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>		<span class="comment">// only gets us to the caller&#39;s fp.</span>
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>		gp.sched.bp = sp - goarch.PtrSize
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>	}
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>	gp.sched.ret = 1
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>	gogo(&amp;gp.sched)
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>}
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span><span class="comment">// fatalthrow implements an unrecoverable runtime throw. It freezes the</span>
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span><span class="comment">// system, prints stack traces starting from its caller, and terminates the</span>
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span><span class="comment">// process.</span>
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>func fatalthrow(t throwType) {
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>	pc := getcallerpc()
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>	sp := getcallersp()
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>	gp := getg()
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>	if gp.m.throwing == throwTypeNone {
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>		gp.m.throwing = t
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>	}
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>	<span class="comment">// Switch to the system stack to avoid any stack growth, which may make</span>
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>	<span class="comment">// things worse if the runtime is in a bad state.</span>
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>		if isSecureMode() {
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>			exit(2)
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>		}
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>		startpanic_m()
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>		if dopanic_m(gp, pc, sp) {
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>			<span class="comment">// crash uses a decent amount of nosplit stack and we&#39;re already</span>
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>			<span class="comment">// low on stack in throw, so crash on the system stack (unlike</span>
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>			<span class="comment">// fatalpanic).</span>
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>			crash()
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>		}
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>		exit(2)
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>	})
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>	*(*int)(nil) = 0 <span class="comment">// not reached</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>}
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span><span class="comment">// fatalpanic implements an unrecoverable panic. It is like fatalthrow, except</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span><span class="comment">// that if msgs != nil, fatalpanic also prints panic messages and decrements</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span><span class="comment">// runningPanicDefers once main is blocked from exiting.</span>
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>func fatalpanic(msgs *_panic) {
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>	pc := getcallerpc()
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>	sp := getcallersp()
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>	gp := getg()
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>	var docrash bool
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>	<span class="comment">// Switch to the system stack to avoid any stack growth, which</span>
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>	<span class="comment">// may make things worse if the runtime is in a bad state.</span>
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>		if startpanic_m() &amp;&amp; msgs != nil {
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>			<span class="comment">// There were panic messages and startpanic_m</span>
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>			<span class="comment">// says it&#39;s okay to try to print them.</span>
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>			<span class="comment">// startpanic_m set panicking, which will</span>
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>			<span class="comment">// block main from exiting, so now OK to</span>
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>			<span class="comment">// decrement runningPanicDefers.</span>
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>			runningPanicDefers.Add(-1)
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>			printpanics(msgs)
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>		}
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>		docrash = dopanic_m(gp, pc, sp)
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>	})
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>	if docrash {
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>		<span class="comment">// By crashing outside the above systemstack call, debuggers</span>
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>		<span class="comment">// will not be confused when generating a backtrace.</span>
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>		<span class="comment">// Function crash is marked nosplit to avoid stack growth.</span>
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>		crash()
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>	}
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>		exit(2)
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>	})
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>	*(*int)(nil) = 0 <span class="comment">// not reached</span>
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>}
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span><span class="comment">// startpanic_m prepares for an unrecoverable panic.</span>
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span><span class="comment">// It returns true if panic messages should be printed, or false if</span>
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span><span class="comment">// the runtime is in bad shape and should just print stacks.</span>
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span><span class="comment">// It must not have write barriers even though the write barrier</span>
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span><span class="comment">// explicitly ignores writes once dying &gt; 0. Write barriers still</span>
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span><span class="comment">// assume that g.m.p != nil, and this function may not have P</span>
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span><span class="comment">// in some contexts (e.g. a panic in a signal handler for a signal</span>
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span><span class="comment">// sent to an M with no P).</span>
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>func startpanic_m() bool {
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>	gp := getg()
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>	if mheap_.cachealloc.size == 0 { <span class="comment">// very early</span>
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>		print(&#34;runtime: panic before malloc heap initialized\n&#34;)
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>	}
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	<span class="comment">// Disallow malloc during an unrecoverable panic. A panic</span>
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>	<span class="comment">// could happen in a signal handler, or in a throw, or inside</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	<span class="comment">// malloc itself. We want to catch if an allocation ever does</span>
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>	<span class="comment">// happen (even if we&#39;re not in one of these situations).</span>
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>	gp.m.mallocing++
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;re dying because of a bad lock count, set it to a</span>
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>	<span class="comment">// good lock count so we don&#39;t recursively panic below.</span>
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>	if gp.m.locks &lt; 0 {
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>		gp.m.locks = 1
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>	}
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>	switch gp.m.dying {
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>	case 0:
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>		<span class="comment">// Setting dying &gt;0 has the side-effect of disabling this G&#39;s writebuf.</span>
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>		gp.m.dying = 1
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>		panicking.Add(1)
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>		lock(&amp;paniclk)
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>		if debug.schedtrace &gt; 0 || debug.scheddetail &gt; 0 {
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>			schedtrace(true)
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>		}
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>		freezetheworld()
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>		return true
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>	case 1:
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>		<span class="comment">// Something failed while panicking.</span>
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>		<span class="comment">// Just print a stack trace and exit.</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>		gp.m.dying = 2
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>		print(&#34;panic during panic\n&#34;)
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>		return false
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>	case 2:
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>		<span class="comment">// This is a genuine bug in the runtime, we couldn&#39;t even</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>		<span class="comment">// print the stack trace successfully.</span>
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>		gp.m.dying = 3
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>		print(&#34;stack trace unavailable\n&#34;)
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>		exit(4)
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>		fallthrough
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>	default:
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>		<span class="comment">// Can&#39;t even print! Just exit.</span>
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>		exit(5)
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>		return false <span class="comment">// Need to return something.</span>
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>	}
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>}
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>var didothers bool
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>var deadlock mutex
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span><span class="comment">// gp is the crashing g running on this M, but may be a user G, while getg() is</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span><span class="comment">// always g0.</span>
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>func dopanic_m(gp *g, pc, sp uintptr) bool {
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>	if gp.sig != 0 {
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>		signame := signame(gp.sig)
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>		if signame != &#34;&#34; {
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>			print(&#34;[signal &#34;, signame)
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>		} else {
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>			print(&#34;[signal &#34;, hex(gp.sig))
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>		}
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>		print(&#34; code=&#34;, hex(gp.sigcode0), &#34; addr=&#34;, hex(gp.sigcode1), &#34; pc=&#34;, hex(gp.sigpc), &#34;]\n&#34;)
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>	}
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>	level, all, docrash := gotraceback()
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>	if level &gt; 0 {
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>		if gp != gp.m.curg {
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>			all = true
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>		}
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>		if gp != gp.m.g0 {
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>			print(&#34;\n&#34;)
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>			goroutineheader(gp)
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>			traceback(pc, sp, 0, gp)
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>		} else if level &gt;= 2 || gp.m.throwing &gt;= throwTypeRuntime {
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>			print(&#34;\nruntime stack:\n&#34;)
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>			traceback(pc, sp, 0, gp)
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>		}
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>		if !didothers &amp;&amp; all {
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>			didothers = true
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>			tracebackothers(gp)
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>		}
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>	}
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>	unlock(&amp;paniclk)
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>	if panicking.Add(-1) != 0 {
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>		<span class="comment">// Some other m is panicking too.</span>
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>		<span class="comment">// Let it print what it needs to print.</span>
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>		<span class="comment">// Wait forever without chewing up cpu.</span>
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>		<span class="comment">// It will exit when it&#39;s done.</span>
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>		lock(&amp;deadlock)
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>		lock(&amp;deadlock)
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>	}
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>	printDebugLog()
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>	return docrash
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>}
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span><span class="comment">// canpanic returns false if a signal should throw instead of</span>
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span><span class="comment">// panicking.</span>
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>func canpanic() bool {
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>	gp := getg()
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>	<span class="comment">// Is it okay for gp to panic instead of crashing the program?</span>
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>	<span class="comment">// Yes, as long as it is running Go code, not runtime code,</span>
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>	<span class="comment">// and not stuck in a system call.</span>
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>	if gp != mp.curg {
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>		releasem(mp)
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>		return false
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>	}
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>	<span class="comment">// N.B. mp.locks != 1 instead of 0 to account for acquirem.</span>
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>	if mp.locks != 1 || mp.mallocing != 0 || mp.throwing != throwTypeNone || mp.preemptoff != &#34;&#34; || mp.dying != 0 {
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>		releasem(mp)
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>		return false
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>	}
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>	status := readgstatus(gp)
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>	if status&amp;^_Gscan != _Grunning || gp.syscallsp != 0 {
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>		releasem(mp)
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>		return false
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>	}
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>	if GOOS == &#34;windows&#34; &amp;&amp; mp.libcallsp != 0 {
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>		releasem(mp)
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>		return false
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>	}
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>	releasem(mp)
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>	return true
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>}
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span><span class="comment">// shouldPushSigpanic reports whether pc should be used as sigpanic&#39;s</span>
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span><span class="comment">// return PC (pushing a frame for the call). Otherwise, it should be</span>
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span><span class="comment">// left alone so that LR is used as sigpanic&#39;s return PC, effectively</span>
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span><span class="comment">// replacing the top-most frame with sigpanic. This is used by</span>
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span><span class="comment">// preparePanic.</span>
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>func shouldPushSigpanic(gp *g, pc, lr uintptr) bool {
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>	if pc == 0 {
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>		<span class="comment">// Probably a call to a nil func. The old LR is more</span>
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>		<span class="comment">// useful in the stack trace. Not pushing the frame</span>
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>		<span class="comment">// will make the trace look like a call to sigpanic</span>
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>		<span class="comment">// instead. (Otherwise the trace will end at sigpanic</span>
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>		<span class="comment">// and we won&#39;t get to see who faulted.)</span>
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>		return false
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>	}
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>	<span class="comment">// If we don&#39;t recognize the PC as code, but we do recognize</span>
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>	<span class="comment">// the link register as code, then this assumes the panic was</span>
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>	<span class="comment">// caused by a call to non-code. In this case, we want to</span>
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>	<span class="comment">// ignore this call to make unwinding show the context.</span>
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>	<span class="comment">// If we running C code, we&#39;re not going to recognize pc as a</span>
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>	<span class="comment">// Go function, so just assume it&#39;s good. Otherwise, traceback</span>
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>	<span class="comment">// may try to read a stale LR that looks like a Go code</span>
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>	<span class="comment">// pointer and wander into the woods.</span>
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>	if gp.m.incgo || findfunc(pc).valid() {
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>		<span class="comment">// This wasn&#39;t a bad call, so use PC as sigpanic&#39;s</span>
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>		<span class="comment">// return PC.</span>
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>		return true
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>	}
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>	if findfunc(lr).valid() {
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>		<span class="comment">// This was a bad call, but the LR is good, so use the</span>
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>		<span class="comment">// LR as sigpanic&#39;s return PC.</span>
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>		return false
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>	}
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span>	<span class="comment">// Neither the PC or LR is good. Hopefully pushing a frame</span>
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span>	<span class="comment">// will work.</span>
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>	return true
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>}
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span><span class="comment">// isAbortPC reports whether pc is the program counter at which</span>
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span><span class="comment">// runtime.abort raises a signal.</span>
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span><span class="comment">// It is nosplit because it&#39;s part of the isgoexception</span>
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span><span class="comment">// implementation.</span>
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>func isAbortPC(pc uintptr) bool {
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>	f := findfunc(pc)
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>	if !f.valid() {
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>		return false
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>	}
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>	return f.funcID == abi.FuncID_abort
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>}
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>
</pre><p><a href="panic.go?m=text">View as plain text</a></p>

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
