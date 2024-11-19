<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/reflect/deepequal.go - Go Documentation Server</title>

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
<a href="deepequal.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/reflect">reflect</a>/<span class="text-muted">deepequal.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/reflect">reflect</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Deep equality test via reflection</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package reflect
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/bytealg&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// During deepValueEqual, must keep track of checks that are</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// in progress. The comparison algorithm assumes that all</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// checks in progress are true when it reencounters them.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Visited comparisons are stored in a map indexed by visit.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type visit struct {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	a1  unsafe.Pointer
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	a2  unsafe.Pointer
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	typ Type
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// Tests for deep equality using reflected types. The map argument tracks</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// comparisons that have already been seen, which allows short circuiting on</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// recursive types.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>func deepValueEqual(v1, v2 Value, visited map[visit]bool) bool {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	if !v1.IsValid() || !v2.IsValid() {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		return v1.IsValid() == v2.IsValid()
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	if v1.Type() != v2.Type() {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		return false
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// We want to avoid putting more in the visited map than we need to.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// For any possible reference cycle that might be encountered,</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// hard(v1, v2) needs to return true for at least one of the types in the cycle,</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// and it&#39;s safe and valid to get Value&#39;s internal pointer.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	hard := func(v1, v2 Value) bool {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		switch v1.Kind() {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		case Pointer:
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>			if v1.typ().PtrBytes == 0 {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>				<span class="comment">// not-in-heap pointers can&#39;t be cyclic.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>				<span class="comment">// At least, all of our current uses of runtime/internal/sys.NotInHeap</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>				<span class="comment">// have that property. The runtime ones aren&#39;t cyclic (and we don&#39;t use</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>				<span class="comment">// DeepEqual on them anyway), and the cgo-generated ones are</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>				<span class="comment">// all empty structs.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>				return false
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			fallthrough
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		case Map, Slice, Interface:
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			<span class="comment">// Nil pointers cannot be cyclic. Avoid putting them in the visited map.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			return !v1.IsNil() &amp;&amp; !v2.IsNil()
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		return false
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	if hard(v1, v2) {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		<span class="comment">// For a Pointer or Map value, we need to check flagIndir,</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		<span class="comment">// which we do by calling the pointer method.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		<span class="comment">// For Slice or Interface, flagIndir is always set,</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		<span class="comment">// and using v.ptr suffices.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		ptrval := func(v Value) unsafe.Pointer {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			switch v.Kind() {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			case Pointer, Map:
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>				return v.pointer()
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			default:
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>				return v.ptr
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		addr1 := ptrval(v1)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		addr2 := ptrval(v2)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		if uintptr(addr1) &gt; uintptr(addr2) {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			<span class="comment">// Canonicalize order to reduce number of entries in visited.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			<span class="comment">// Assumes non-moving garbage collector.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			addr1, addr2 = addr2, addr1
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		<span class="comment">// Short circuit if references are already seen.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		typ := v1.Type()
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		v := visit{addr1, addr2, typ}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		if visited[v] {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			return true
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		<span class="comment">// Remember for later.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		visited[v] = true
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	switch v1.Kind() {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	case Array:
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		for i := 0; i &lt; v1.Len(); i++ {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			if !deepValueEqual(v1.Index(i), v2.Index(i), visited) {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>				return false
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		return true
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	case Slice:
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		if v1.IsNil() != v2.IsNil() {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			return false
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		if v1.Len() != v2.Len() {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			return false
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		if v1.UnsafePointer() == v2.UnsafePointer() {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			return true
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		<span class="comment">// Special case for []byte, which is common.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		if v1.Type().Elem().Kind() == Uint8 {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			return bytealg.Equal(v1.Bytes(), v2.Bytes())
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		for i := 0; i &lt; v1.Len(); i++ {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			if !deepValueEqual(v1.Index(i), v2.Index(i), visited) {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>				return false
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		return true
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	case Interface:
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		if v1.IsNil() || v2.IsNil() {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			return v1.IsNil() == v2.IsNil()
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		return deepValueEqual(v1.Elem(), v2.Elem(), visited)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	case Pointer:
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		if v1.UnsafePointer() == v2.UnsafePointer() {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			return true
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		return deepValueEqual(v1.Elem(), v2.Elem(), visited)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	case Struct:
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		for i, n := 0, v1.NumField(); i &lt; n; i++ {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			if !deepValueEqual(v1.Field(i), v2.Field(i), visited) {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>				return false
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		return true
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	case Map:
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		if v1.IsNil() != v2.IsNil() {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			return false
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		if v1.Len() != v2.Len() {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			return false
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		if v1.UnsafePointer() == v2.UnsafePointer() {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			return true
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		iter := v1.MapRange()
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		for iter.Next() {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			val1 := iter.Value()
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			val2 := v2.MapIndex(iter.Key())
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			if !val1.IsValid() || !val2.IsValid() || !deepValueEqual(val1, val2, visited) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>				return false
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		return true
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	case Func:
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		if v1.IsNil() &amp;&amp; v2.IsNil() {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			return true
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// Can&#39;t do better than this:</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		return false
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	case Int, Int8, Int16, Int32, Int64:
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		return v1.Int() == v2.Int()
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	case Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		return v1.Uint() == v2.Uint()
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	case String:
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		return v1.String() == v2.String()
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	case Bool:
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		return v1.Bool() == v2.Bool()
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	case Float32, Float64:
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		return v1.Float() == v2.Float()
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	case Complex64, Complex128:
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		return v1.Complex() == v2.Complex()
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	default:
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		<span class="comment">// Normal equality suffices</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		return valueInterface(v1, false) == valueInterface(v2, false)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">// DeepEqual reports whether x and y are “deeply equal,” defined as follows.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">// Two values of identical type are deeply equal if one of the following cases applies.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span><span class="comment">// Values of distinct types are never deeply equal.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">// Array values are deeply equal when their corresponding elements are deeply equal.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span><span class="comment">// Struct values are deeply equal if their corresponding fields,</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span><span class="comment">// both exported and unexported, are deeply equal.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// Func values are deeply equal if both are nil; otherwise they are not deeply equal.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// Interface values are deeply equal if they hold deeply equal concrete values.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// Map values are deeply equal when all of the following are true:</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// they are both nil or both non-nil, they have the same length,</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">// and either they are the same map object or their corresponding keys</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span><span class="comment">// (matched using Go equality) map to deeply equal values.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// Pointer values are deeply equal if they are equal using Go&#39;s == operator</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// or if they point to deeply equal values.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// Slice values are deeply equal when all of the following are true:</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// they are both nil or both non-nil, they have the same length,</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// and either they point to the same initial entry of the same underlying array</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">// (that is, &amp;x[0] == &amp;y[0]) or their corresponding elements (up to length) are deeply equal.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">// Note that a non-nil empty slice and a nil slice (for example, []byte{} and []byte(nil))</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">// are not deeply equal.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// Other values - numbers, bools, strings, and channels - are deeply equal</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// if they are equal using Go&#39;s == operator.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// In general DeepEqual is a recursive relaxation of Go&#39;s == operator.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// However, this idea is impossible to implement without some inconsistency.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// Specifically, it is possible for a value to be unequal to itself,</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// either because it is of func type (uncomparable in general)</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// or because it is a floating-point NaN value (not equal to itself in floating-point comparison),</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// or because it is an array, struct, or interface containing</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// such a value.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// On the other hand, pointer values are always equal to themselves,</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// even if they point at or contain such problematic values,</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// because they compare equal using Go&#39;s == operator, and that</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">// is a sufficient condition to be deeply equal, regardless of content.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">// DeepEqual has been defined so that the same short-cut applies</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// to slices and maps: if x and y are the same slice or the same map,</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// they are deeply equal regardless of content.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// As DeepEqual traverses the data values it may find a cycle. The</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// second and subsequent times that DeepEqual compares two pointer</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// values that have been compared before, it treats the values as</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// equal rather than examining the values to which they point.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// This ensures that DeepEqual terminates.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>func DeepEqual(x, y any) bool {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	if x == nil || y == nil {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		return x == y
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	v1 := ValueOf(x)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	v2 := ValueOf(y)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	if v1.Type() != v2.Type() {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		return false
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	return deepValueEqual(v1, v2, make(map[visit]bool))
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
</pre><p><a href="deepequal.go?m=text">View as plain text</a></p>

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
