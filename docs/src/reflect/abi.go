<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/reflect/abi.go - Go Documentation Server</title>

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
<a href="abi.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/reflect">reflect</a>/<span class="text-muted">abi.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/reflect">reflect</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2021 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package reflect
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// These variables are used by the register assignment</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// algorithm in this file.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// They should be modified with care (no other reflect code</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// may be executing) and are generally only modified</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// when testing this package.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// They should never be set higher than their internal/abi</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// constant counterparts, because the system relies on a</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// structure that is at least large enough to hold the</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// registers the system supports.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// Currently they&#39;re set to zero because using the actual</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// constants will break every part of the toolchain that</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// uses reflect to call functions (e.g. go test, or anything</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// that uses text/template). The values that are currently</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// commented out there should be the actual values once</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// we&#39;re ready to use the register ABI everywhere.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>var (
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	intArgRegs   = abi.IntArgRegs
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	floatArgRegs = abi.FloatArgRegs
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	floatRegSize = uintptr(abi.EffectiveFloatRegSize)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// abiStep represents an ABI &#34;instruction.&#34; Each instruction</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// describes one part of how to translate between a Go value</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// in memory and a call frame.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>type abiStep struct {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	kind abiStepKind
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// offset and size together describe a part of a Go value</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// in memory.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	offset uintptr
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	size   uintptr <span class="comment">// size in bytes of the part</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// These fields describe the ABI side of the translation.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	stkOff uintptr <span class="comment">// stack offset, used if kind == abiStepStack</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	ireg   int     <span class="comment">// integer register index, used if kind == abiStepIntReg or kind == abiStepPointer</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	freg   int     <span class="comment">// FP register index, used if kind == abiStepFloatReg</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// abiStepKind is the &#34;op-code&#34; for an abiStep instruction.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>type abiStepKind int
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>const (
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	abiStepBad      abiStepKind = iota
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	abiStepStack                <span class="comment">// copy to/from stack</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	abiStepIntReg               <span class="comment">// copy to/from integer register</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	abiStepPointer              <span class="comment">// copy pointer to/from integer register</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	abiStepFloatReg             <span class="comment">// copy to/from FP register</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// abiSeq represents a sequence of ABI instructions for copying</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// from a series of reflect.Values to a call frame (for call arguments)</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// or vice-versa (for call results).</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// An abiSeq should be populated by calling its addArg method.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>type abiSeq struct {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// steps is the set of instructions.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// The instructions are grouped together by whole arguments,</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// with the starting index for the instructions</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// of the i&#39;th Go value available in valueStart.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// For instance, if this abiSeq represents 3 arguments</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// passed to a function, then the 2nd argument&#39;s steps</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// begin at steps[valueStart[1]].</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// Because reflect accepts Go arguments in distinct</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// Values and each Value is stored separately, each abiStep</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// that begins a new argument will have its offset</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// field == 0.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	steps      []abiStep
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	valueStart []int
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	stackBytes   uintptr <span class="comment">// stack space used</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	iregs, fregs int     <span class="comment">// registers used</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>func (a *abiSeq) dump() {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	for i, p := range a.steps {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		println(&#34;part&#34;, i, p.kind, p.offset, p.size, p.stkOff, p.ireg, p.freg)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	print(&#34;values &#34;)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	for _, i := range a.valueStart {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		print(i, &#34; &#34;)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	println()
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	println(&#34;stack&#34;, a.stackBytes)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	println(&#34;iregs&#34;, a.iregs)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	println(&#34;fregs&#34;, a.fregs)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// stepsForValue returns the ABI instructions for translating</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// the i&#39;th Go argument or return value represented by this</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// abiSeq to the Go ABI.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>func (a *abiSeq) stepsForValue(i int) []abiStep {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	s := a.valueStart[i]
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	var e int
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	if i == len(a.valueStart)-1 {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		e = len(a.steps)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	} else {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		e = a.valueStart[i+1]
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	return a.steps[s:e]
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// addArg extends the abiSeq with a new Go value of type t.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// If the value was stack-assigned, returns the single</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// abiStep describing that translation, and nil otherwise.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func (a *abiSeq) addArg(t *abi.Type) *abiStep {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// We&#39;ll always be adding a new value, so do that first.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	pStart := len(a.steps)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	a.valueStart = append(a.valueStart, pStart)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	if t.Size() == 0 {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		<span class="comment">// If the size of the argument type is zero, then</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		<span class="comment">// in order to degrade gracefully into ABI0, we need</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		<span class="comment">// to stack-assign this type. The reason is that</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		<span class="comment">// although zero-sized types take up no space on the</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		<span class="comment">// stack, they do cause the next argument to be aligned.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// So just do that here, but don&#39;t bother actually</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// generating a new ABI step for it (there&#39;s nothing to</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// actually copy).</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		<span class="comment">// We cannot handle this in the recursive case of</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		<span class="comment">// regAssign because zero-sized *fields* of a</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		<span class="comment">// non-zero-sized struct do not cause it to be</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		<span class="comment">// stack-assigned. So we need a special case here</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		<span class="comment">// at the top.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		a.stackBytes = align(a.stackBytes, uintptr(t.Align()))
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		return nil
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// Hold a copy of &#34;a&#34; so that we can roll back if</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">// register assignment fails.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	aOld := *a
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	if !a.regAssign(t, 0) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		<span class="comment">// Register assignment failed. Roll back any changes</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		<span class="comment">// and stack-assign.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		*a = aOld
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		a.stackAssign(t.Size(), uintptr(t.Align()))
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return &amp;a.steps[len(a.steps)-1]
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	return nil
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// addRcvr extends the abiSeq with a new method call</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// receiver according to the interface calling convention.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// If the receiver was stack-assigned, returns the single</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// abiStep describing that translation, and nil otherwise.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// Returns true if the receiver is a pointer.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>func (a *abiSeq) addRcvr(rcvr *abi.Type) (*abiStep, bool) {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// The receiver is always one word.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	a.valueStart = append(a.valueStart, len(a.steps))
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	var ok, ptr bool
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	if ifaceIndir(rcvr) || rcvr.Pointers() {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		ok = a.assignIntN(0, goarch.PtrSize, 1, 0b1)
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		ptr = true
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	} else {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		<span class="comment">// TODO(mknyszek): Is this case even possible?</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		<span class="comment">// The interface data work never contains a non-pointer</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		<span class="comment">// value. This case was copied over from older code</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		<span class="comment">// in the reflect package which only conditionally added</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		<span class="comment">// a pointer bit to the reflect.(Value).Call stack frame&#39;s</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		<span class="comment">// GC bitmap.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		ok = a.assignIntN(0, goarch.PtrSize, 1, 0b0)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		ptr = false
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	if !ok {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		a.stackAssign(goarch.PtrSize, goarch.PtrSize)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		return &amp;a.steps[len(a.steps)-1], ptr
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	return nil, ptr
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// regAssign attempts to reserve argument registers for a value of</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// type t, stored at some offset.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// It returns whether or not the assignment succeeded, but</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">// leaves any changes it made to a.steps behind, so the caller</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span><span class="comment">// must undo that work by adjusting a.steps if it fails.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// This method along with the assign* methods represent the</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// complete register-assignment algorithm for the Go ABI.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>func (a *abiSeq) regAssign(t *abi.Type, offset uintptr) bool {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	switch Kind(t.Kind()) {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	case UnsafePointer, Pointer, Chan, Map, Func:
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		return a.assignIntN(offset, t.Size(), 1, 0b1)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	case Bool, Int, Uint, Int8, Uint8, Int16, Uint16, Int32, Uint32, Uintptr:
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		return a.assignIntN(offset, t.Size(), 1, 0b0)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	case Int64, Uint64:
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		switch goarch.PtrSize {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		case 4:
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			return a.assignIntN(offset, 4, 2, 0b0)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		case 8:
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			return a.assignIntN(offset, 8, 1, 0b0)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	case Float32, Float64:
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		return a.assignFloatN(offset, t.Size(), 1)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	case Complex64:
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		return a.assignFloatN(offset, 4, 2)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	case Complex128:
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		return a.assignFloatN(offset, 8, 2)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	case String:
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		return a.assignIntN(offset, goarch.PtrSize, 2, 0b01)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	case Interface:
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		return a.assignIntN(offset, goarch.PtrSize, 2, 0b10)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	case Slice:
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		return a.assignIntN(offset, goarch.PtrSize, 3, 0b001)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	case Array:
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		tt := (*arrayType)(unsafe.Pointer(t))
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		switch tt.Len {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		case 0:
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			<span class="comment">// There&#39;s nothing to assign, so don&#39;t modify</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			<span class="comment">// a.steps but succeed so the caller doesn&#39;t</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			<span class="comment">// try to stack-assign this value.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			return true
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		case 1:
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			return a.regAssign(tt.Elem, offset)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		default:
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			return false
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	case Struct:
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		st := (*structType)(unsafe.Pointer(t))
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		for i := range st.Fields {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			f := &amp;st.Fields[i]
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			if !a.regAssign(f.Typ, offset+f.Offset) {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>				return false
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		return true
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	default:
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		print(&#34;t.Kind == &#34;, t.Kind(), &#34;\n&#34;)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		panic(&#34;unknown type kind&#34;)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	panic(&#34;unhandled register assignment path&#34;)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span><span class="comment">// assignIntN assigns n values to registers, each &#34;size&#34; bytes large,</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span><span class="comment">// from the data at [offset, offset+n*size) in memory. Each value at</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// [offset+i*size, offset+(i+1)*size) for i &lt; n is assigned to the</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// next n integer registers.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span><span class="comment">// Bit i in ptrMap indicates whether the i&#39;th value is a pointer.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span><span class="comment">// n must be &lt;= 8.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// Returns whether assignment succeeded.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>func (a *abiSeq) assignIntN(offset, size uintptr, n int, ptrMap uint8) bool {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	if n &gt; 8 || n &lt; 0 {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		panic(&#34;invalid n&#34;)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	if ptrMap != 0 &amp;&amp; size != goarch.PtrSize {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		panic(&#34;non-empty pointer map passed for non-pointer-size values&#34;)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	if a.iregs+n &gt; intArgRegs {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		return false
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		kind := abiStepIntReg
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		if ptrMap&amp;(uint8(1)&lt;&lt;i) != 0 {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			kind = abiStepPointer
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		a.steps = append(a.steps, abiStep{
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			kind:   kind,
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			offset: offset + uintptr(i)*size,
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			size:   size,
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			ireg:   a.iregs,
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		})
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		a.iregs++
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	return true
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span><span class="comment">// assignFloatN assigns n values to registers, each &#34;size&#34; bytes large,</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">// from the data at [offset, offset+n*size) in memory. Each value at</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span><span class="comment">// [offset+i*size, offset+(i+1)*size) for i &lt; n is assigned to the</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">// next n floating-point registers.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">// Returns whether assignment succeeded.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>func (a *abiSeq) assignFloatN(offset, size uintptr, n int) bool {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	if n &lt; 0 {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		panic(&#34;invalid n&#34;)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	if a.fregs+n &gt; floatArgRegs || floatRegSize &lt; size {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		return false
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		a.steps = append(a.steps, abiStep{
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			kind:   abiStepFloatReg,
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			offset: offset + uintptr(i)*size,
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			size:   size,
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			freg:   a.fregs,
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		})
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		a.fregs++
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	return true
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>}
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">// stackAssign reserves space for one value that is &#34;size&#34; bytes</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span><span class="comment">// large with alignment &#34;alignment&#34; to the stack.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// Should not be called directly; use addArg instead.</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>func (a *abiSeq) stackAssign(size, alignment uintptr) {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	a.stackBytes = align(a.stackBytes, alignment)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	a.steps = append(a.steps, abiStep{
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		kind:   abiStepStack,
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		offset: 0, <span class="comment">// Only used for whole arguments, so the memory offset is 0.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		size:   size,
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		stkOff: a.stackBytes,
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	})
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	a.stackBytes += size
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span><span class="comment">// abiDesc describes the ABI for a function or method.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>type abiDesc struct {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	<span class="comment">// call and ret represent the translation steps for</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">// the call and return paths of a Go function.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	call, ret abiSeq
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// These fields describe the stack space allocated</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// for the call. stackCallArgsSize is the amount of space</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// reserved for arguments but not return values. retOffset</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// is the offset at which return values begin, and</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// spill is the size in bytes of additional space reserved</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// to spill argument registers into in case of preemption in</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// reflectcall&#39;s stack frame.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	stackCallArgsSize, retOffset, spill uintptr
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	<span class="comment">// stackPtrs is a bitmap that indicates whether</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	<span class="comment">// each word in the ABI stack space (stack-assigned</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	<span class="comment">// args + return values) is a pointer. Used</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	<span class="comment">// as the heap pointer bitmap for stack space</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	<span class="comment">// passed to reflectcall.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	stackPtrs *bitVector
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	<span class="comment">// inRegPtrs is a bitmap whose i&#39;th bit indicates</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	<span class="comment">// whether the i&#39;th integer argument register contains</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	<span class="comment">// a pointer. Used by makeFuncStub and methodValueCall</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	<span class="comment">// to make result pointers visible to the GC.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	<span class="comment">// outRegPtrs is the same, but for result values.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	<span class="comment">// Used by reflectcall to make result pointers visible</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	<span class="comment">// to the GC.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	inRegPtrs, outRegPtrs abi.IntArgRegBitmap
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>func (a *abiDesc) dump() {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	println(&#34;ABI&#34;)
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	println(&#34;call&#34;)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	a.call.dump()
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	println(&#34;ret&#34;)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	a.ret.dump()
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	println(&#34;stackCallArgsSize&#34;, a.stackCallArgsSize)
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	println(&#34;retOffset&#34;, a.retOffset)
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	println(&#34;spill&#34;, a.spill)
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	print(&#34;inRegPtrs:&#34;)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	dumpPtrBitMap(a.inRegPtrs)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	println()
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	print(&#34;outRegPtrs:&#34;)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	dumpPtrBitMap(a.outRegPtrs)
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	println()
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>}
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>func dumpPtrBitMap(b abi.IntArgRegBitmap) {
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	for i := 0; i &lt; intArgRegs; i++ {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		x := 0
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		if b.Get(i) {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			x = 1
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		print(&#34; &#34;, x)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>func newAbiDesc(t *funcType, rcvr *abi.Type) abiDesc {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	<span class="comment">// We need to add space for this argument to</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	<span class="comment">// the frame so that it can spill args into it.</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	<span class="comment">// The size of this space is just the sum of the sizes</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	<span class="comment">// of each register-allocated type.</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): Remove this when we no longer have</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	<span class="comment">// caller reserved spill space.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	spill := uintptr(0)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	<span class="comment">// Compute gc program &amp; stack bitmap for stack arguments</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	stackPtrs := new(bitVector)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	<span class="comment">// Compute the stack frame pointer bitmap and register</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	<span class="comment">// pointer bitmap for arguments.</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	inRegPtrs := abi.IntArgRegBitmap{}
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	<span class="comment">// Compute abiSeq for input parameters.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	var in abiSeq
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	if rcvr != nil {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		stkStep, isPtr := in.addRcvr(rcvr)
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		if stkStep != nil {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			if isPtr {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>				stackPtrs.append(1)
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			} else {
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>				stackPtrs.append(0)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		} else {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			spill += goarch.PtrSize
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	for i, arg := range t.InSlice() {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		stkStep := in.addArg(arg)
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		if stkStep != nil {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			addTypeBits(stackPtrs, stkStep.stkOff, arg)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		} else {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			spill = align(spill, uintptr(arg.Align()))
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>			spill += arg.Size()
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			for _, st := range in.stepsForValue(i) {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>				if st.kind == abiStepPointer {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>					inRegPtrs.Set(st.ireg)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>				}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	spill = align(spill, goarch.PtrSize)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	<span class="comment">// From the input parameters alone, we now know</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	<span class="comment">// the stackCallArgsSize and retOffset.</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	stackCallArgsSize := in.stackBytes
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	retOffset := align(in.stackBytes, goarch.PtrSize)
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	<span class="comment">// Compute the stack frame pointer bitmap and register</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// pointer bitmap for return values.</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	outRegPtrs := abi.IntArgRegBitmap{}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	<span class="comment">// Compute abiSeq for output parameters.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	var out abiSeq
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	<span class="comment">// Stack-assigned return values do not share</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	<span class="comment">// space with arguments like they do with registers,</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	<span class="comment">// so we need to inject a stack offset here.</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	<span class="comment">// Fake it by artificially extending stackBytes by</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	<span class="comment">// the return offset.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	out.stackBytes = retOffset
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	for i, res := range t.OutSlice() {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		stkStep := out.addArg(res)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		if stkStep != nil {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			addTypeBits(stackPtrs, stkStep.stkOff, res)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		} else {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			for _, st := range out.stepsForValue(i) {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>				if st.kind == abiStepPointer {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>					outRegPtrs.Set(st.ireg)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>				}
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	<span class="comment">// Undo the faking from earlier so that stackBytes</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	<span class="comment">// is accurate.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	out.stackBytes -= retOffset
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	return abiDesc{in, out, stackCallArgsSize, retOffset, spill, stackPtrs, inRegPtrs, outRegPtrs}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span><span class="comment">// intFromReg loads an argSize sized integer from reg and places it at to.</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span><span class="comment">// argSize must be non-zero, fit in a register, and a power-of-two.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>func intFromReg(r *abi.RegArgs, reg int, argSize uintptr, to unsafe.Pointer) {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	memmove(to, r.IntRegArgAddr(reg, argSize), argSize)
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span><span class="comment">// intToReg loads an argSize sized integer and stores it into reg.</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span><span class="comment">// argSize must be non-zero, fit in a register, and a power-of-two.</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>func intToReg(r *abi.RegArgs, reg int, argSize uintptr, from unsafe.Pointer) {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	memmove(r.IntRegArgAddr(reg, argSize), from, argSize)
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span><span class="comment">// floatFromReg loads a float value from its register representation in r.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span><span class="comment">// argSize must be 4 or 8.</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>func floatFromReg(r *abi.RegArgs, reg int, argSize uintptr, to unsafe.Pointer) {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	switch argSize {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	case 4:
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		*(*float32)(to) = archFloat32FromReg(r.Floats[reg])
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	case 8:
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		*(*float64)(to) = *(*float64)(unsafe.Pointer(&amp;r.Floats[reg]))
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	default:
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		panic(&#34;bad argSize&#34;)
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span><span class="comment">// floatToReg stores a float value in its register representation in r.</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span><span class="comment">// argSize must be either 4 or 8.</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>func floatToReg(r *abi.RegArgs, reg int, argSize uintptr, from unsafe.Pointer) {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	switch argSize {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	case 4:
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		r.Floats[reg] = archFloat32ToReg(*(*float32)(from))
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	case 8:
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		r.Floats[reg] = *(*uint64)(from)
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	default:
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		panic(&#34;bad argSize&#34;)
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>
</pre><p><a href="abi.go?m=text">View as plain text</a></p>

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
