<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/iface.go - Go Documentation Server</title>

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
<a href="iface.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">iface.go</span>
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
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>const itabInitSize = 512
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>var (
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	itabLock      mutex                               <span class="comment">// lock for accessing itab table</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	itabTable     = &amp;itabTableInit                    <span class="comment">// pointer to current table</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	itabTableInit = itabTableType{size: itabInitSize} <span class="comment">// starter table</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// Note: change the formula in the mallocgc call in itabAdd if you change these fields.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>type itabTableType struct {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	size    uintptr             <span class="comment">// length of entries array. Always a power of 2.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	count   uintptr             <span class="comment">// current number of filled entries.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	entries [itabInitSize]*itab <span class="comment">// really [size] large</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>func itabHashFunc(inter *interfacetype, typ *_type) uintptr {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// compiler has provided some good hash codes for us.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	return uintptr(inter.Type.Hash ^ typ.Hash)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>func getitab(inter *interfacetype, typ *_type, canfail bool) *itab {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	if len(inter.Methods) == 0 {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		throw(&#34;internal error - misuse of itab&#34;)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// easy case</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	if typ.TFlag&amp;abi.TFlagUncommon == 0 {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		if canfail {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>			return nil
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		name := toRType(&amp;inter.Type).nameOff(inter.Methods[0].Name)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		panic(&amp;TypeAssertionError{nil, typ, &amp;inter.Type, name.Name()})
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	var m *itab
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// First, look in the existing table to see if we can find the itab we need.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// This is by far the most common case, so do it without locks.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// Use atomic to ensure we see any previous writes done by the thread</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// that updates the itabTable field (with atomic.Storep in itabAdd).</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	t := (*itabTableType)(atomic.Loadp(unsafe.Pointer(&amp;itabTable)))
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	if m = t.find(inter, typ); m != nil {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		goto finish
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// Not found.  Grab the lock and try again.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	lock(&amp;itabLock)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	if m = itabTable.find(inter, typ); m != nil {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		unlock(&amp;itabLock)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		goto finish
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// Entry doesn&#39;t exist yet. Make a new entry &amp; add it.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	m = (*itab)(persistentalloc(unsafe.Sizeof(itab{})+uintptr(len(inter.Methods)-1)*goarch.PtrSize, 0, &amp;memstats.other_sys))
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	m.inter = inter
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	m._type = typ
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// The hash is used in type switches. However, compiler statically generates itab&#39;s</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// for all interface/type pairs used in switches (which are added to itabTable</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// in itabsinit). The dynamically-generated itab&#39;s never participate in type switches,</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// and thus the hash is irrelevant.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// Note: m.hash is _not_ the hash used for the runtime itabTable hash table.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	m.hash = 0
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	m.init()
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	itabAdd(m)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	unlock(&amp;itabLock)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>finish:
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	if m.fun[0] != 0 {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		return m
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	if canfail {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		return nil
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// this can only happen if the conversion</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// was already done once using the , ok form</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// and we have a cached negative result.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// The cached result doesn&#39;t record which</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// interface function was missing, so initialize</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// the itab again to get the missing function name.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	panic(&amp;TypeAssertionError{concrete: typ, asserted: &amp;inter.Type, missingMethod: m.init()})
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// find finds the given interface/type pair in t.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// Returns nil if the given interface/type pair isn&#39;t present.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>func (t *itabTableType) find(inter *interfacetype, typ *_type) *itab {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// Implemented using quadratic probing.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// Probe sequence is h(i) = h0 + i*(i+1)/2 mod 2^k.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// We&#39;re guaranteed to hit all table entries using this probe sequence.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	mask := t.size - 1
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	h := itabHashFunc(inter, typ) &amp; mask
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	for i := uintptr(1); ; i++ {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		p := (**itab)(add(unsafe.Pointer(&amp;t.entries), h*goarch.PtrSize))
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">// Use atomic read here so if we see m != nil, we also see</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		<span class="comment">// the initializations of the fields of m.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		<span class="comment">// m := *p</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		m := (*itab)(atomic.Loadp(unsafe.Pointer(p)))
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		if m == nil {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			return nil
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		if m.inter == inter &amp;&amp; m._type == typ {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			return m
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		h += i
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		h &amp;= mask
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// itabAdd adds the given itab to the itab hash table.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// itabLock must be held.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>func itabAdd(m *itab) {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// Bugs can lead to calling this while mallocing is set,</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// typically because this is called while panicking.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// Crash reliably, rather than only when we need to grow</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// the hash table.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	if getg().m.mallocing != 0 {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		throw(&#34;malloc deadlock&#34;)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	t := itabTable
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if t.count &gt;= 3*(t.size/4) { <span class="comment">// 75% load factor</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// Grow hash table.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// t2 = new(itabTableType) + some additional entries</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// We lie and tell malloc we want pointer-free memory because</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		<span class="comment">// all the pointed-to values are not in the heap.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		t2 := (*itabTableType)(mallocgc((2+2*t.size)*goarch.PtrSize, nil, true))
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		t2.size = t.size * 2
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		<span class="comment">// Copy over entries.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		<span class="comment">// Note: while copying, other threads may look for an itab and</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		<span class="comment">// fail to find it. That&#39;s ok, they will then try to get the itab lock</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">// and as a consequence wait until this copying is complete.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		iterate_itabs(t2.add)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		if t2.count != t.count {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			throw(&#34;mismatched count during itab table copy&#34;)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		<span class="comment">// Publish new hash table. Use an atomic write: see comment in getitab.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		atomicstorep(unsafe.Pointer(&amp;itabTable), unsafe.Pointer(t2))
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		<span class="comment">// Adopt the new table as our own.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		t = itabTable
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		<span class="comment">// Note: the old table can be GC&#39;ed here.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	t.add(m)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">// add adds the given itab to itab table t.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// itabLock must be held.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func (t *itabTableType) add(m *itab) {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// See comment in find about the probe sequence.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// Insert new itab in the first empty spot in the probe sequence.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	mask := t.size - 1
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	h := itabHashFunc(m.inter, m._type) &amp; mask
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	for i := uintptr(1); ; i++ {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		p := (**itab)(add(unsafe.Pointer(&amp;t.entries), h*goarch.PtrSize))
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		m2 := *p
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		if m2 == m {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			<span class="comment">// A given itab may be used in more than one module</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			<span class="comment">// and thanks to the way global symbol resolution works, the</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			<span class="comment">// pointed-to itab may already have been inserted into the</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			<span class="comment">// global &#39;hash&#39;.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			return
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		if m2 == nil {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			<span class="comment">// Use atomic write here so if a reader sees m, it also</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			<span class="comment">// sees the correctly initialized fields of m.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			<span class="comment">// NoWB is ok because m is not in heap memory.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			<span class="comment">// *p = m</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			atomic.StorepNoWB(unsafe.Pointer(p), unsafe.Pointer(m))
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			t.count++
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			return
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		h += i
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		h &amp;= mask
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// init fills in the m.fun array with all the code pointers for</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// the m.inter/m._type pair. If the type does not implement the interface,</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// it sets m.fun[0] to 0 and returns the name of an interface function that is missing.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// It is ok to call this multiple times on the same m, even concurrently.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>func (m *itab) init() string {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	inter := m.inter
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	typ := m._type
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	x := typ.Uncommon()
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// both inter and typ have method sorted by name,</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// and interface names are unique,</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// so can iterate over both in lock step;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// the loop is O(ni+nt) not O(ni*nt).</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	ni := len(inter.Methods)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	nt := int(x.Mcount)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	xmhdr := (*[1 &lt;&lt; 16]abi.Method)(add(unsafe.Pointer(x), uintptr(x.Moff)))[:nt:nt]
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	j := 0
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	methods := (*[1 &lt;&lt; 16]unsafe.Pointer)(unsafe.Pointer(&amp;m.fun[0]))[:ni:ni]
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	var fun0 unsafe.Pointer
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>imethods:
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	for k := 0; k &lt; ni; k++ {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		i := &amp;inter.Methods[k]
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		itype := toRType(&amp;inter.Type).typeOff(i.Typ)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		name := toRType(&amp;inter.Type).nameOff(i.Name)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		iname := name.Name()
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		ipkg := pkgPath(name)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		if ipkg == &#34;&#34; {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			ipkg = inter.PkgPath.Name()
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		for ; j &lt; nt; j++ {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			t := &amp;xmhdr[j]
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			rtyp := toRType(typ)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			tname := rtyp.nameOff(t.Name)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			if rtyp.typeOff(t.Mtyp) == itype &amp;&amp; tname.Name() == iname {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>				pkgPath := pkgPath(tname)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				if pkgPath == &#34;&#34; {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>					pkgPath = rtyp.nameOff(x.PkgPath).Name()
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>				}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>				if tname.IsExported() || pkgPath == ipkg {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>					ifn := rtyp.textOff(t.Ifn)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>					if k == 0 {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>						fun0 = ifn <span class="comment">// we&#39;ll set m.fun[0] at the end</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>					} else {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>						methods[k] = ifn
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>					}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>					continue imethods
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>				}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		<span class="comment">// didn&#39;t find method</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		m.fun[0] = 0
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		return iname
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	m.fun[0] = uintptr(fun0)
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	return &#34;&#34;
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>func itabsinit() {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	lockInit(&amp;itabLock, lockRankItab)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	lock(&amp;itabLock)
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	for _, md := range activeModules() {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		for _, i := range md.itablinks {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			itabAdd(i)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	unlock(&amp;itabLock)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span><span class="comment">// panicdottypeE is called when doing an e.(T) conversion and the conversion fails.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span><span class="comment">// have = the dynamic type we have.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">// want = the static type we&#39;re trying to convert to.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// iface = the static type we&#39;re converting from.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>func panicdottypeE(have, want, iface *_type) {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	panic(&amp;TypeAssertionError{iface, have, want, &#34;&#34;})
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">// panicdottypeI is called when doing an i.(T) conversion and the conversion fails.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">// Same args as panicdottypeE, but &#34;have&#34; is the dynamic itab we have.</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>func panicdottypeI(have *itab, want, iface *_type) {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	var t *_type
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	if have != nil {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		t = have._type
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	panicdottypeE(t, want, iface)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span><span class="comment">// panicnildottype is called when doing an i.(T) conversion and the interface i is nil.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span><span class="comment">// want = the static type we&#39;re trying to convert to.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>func panicnildottype(want *_type) {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	panic(&amp;TypeAssertionError{nil, nil, want, &#34;&#34;})
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	<span class="comment">// TODO: Add the static type we&#39;re converting from as well.</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	<span class="comment">// It might generate a better error message.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	<span class="comment">// Just to match other nil conversion errors, we don&#39;t for now.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">// The specialized convTx routines need a type descriptor to use when calling mallocgc.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">// We don&#39;t need the type to be exact, just to have the correct size, alignment, and pointer-ness.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">// However, when debugging, it&#39;d be nice to have some indication in mallocgc where the types came from,</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span><span class="comment">// so we use named types here.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">// We then construct interface values of these types,</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span><span class="comment">// and then extract the type word to use as needed.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>type (
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	uint16InterfacePtr uint16
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	uint32InterfacePtr uint32
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	uint64InterfacePtr uint64
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	stringInterfacePtr string
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	sliceInterfacePtr  []byte
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>var (
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	uint16Eface any = uint16InterfacePtr(0)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	uint32Eface any = uint32InterfacePtr(0)
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	uint64Eface any = uint64InterfacePtr(0)
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	stringEface any = stringInterfacePtr(&#34;&#34;)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	sliceEface  any = sliceInterfacePtr(nil)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	uint16Type *_type = efaceOf(&amp;uint16Eface)._type
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	uint32Type *_type = efaceOf(&amp;uint32Eface)._type
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	uint64Type *_type = efaceOf(&amp;uint64Eface)._type
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	stringType *_type = efaceOf(&amp;stringEface)._type
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	sliceType  *_type = efaceOf(&amp;sliceEface)._type
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">// The conv and assert functions below do very similar things.</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span><span class="comment">// The convXXX functions are guaranteed by the compiler to succeed.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">// The assertXXX functions may fail (either panicking or returning false,</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// depending on whether they are 1-result or 2-result).</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">// The convXXX functions succeed on a nil input, whereas the assertXXX</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">// functions fail on a nil input.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">// convT converts a value of type t, which is pointed to by v, to a pointer that can</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span><span class="comment">// be used as the second word of an interface value.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>func convT(t *_type, v unsafe.Pointer) unsafe.Pointer {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	if raceenabled {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		raceReadObjectPC(t, v, getcallerpc(), abi.FuncPCABIInternal(convT))
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	if msanenabled {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		msanread(v, t.Size_)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if asanenabled {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		asanread(v, t.Size_)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	x := mallocgc(t.Size_, t, true)
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	typedmemmove(t, x, v)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	return x
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>func convTnoptr(t *_type, v unsafe.Pointer) unsafe.Pointer {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// TODO: maybe take size instead of type?</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	if raceenabled {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		raceReadObjectPC(t, v, getcallerpc(), abi.FuncPCABIInternal(convTnoptr))
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	if msanenabled {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		msanread(v, t.Size_)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	if asanenabled {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		asanread(v, t.Size_)
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	x := mallocgc(t.Size_, t, false)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	memmove(x, v, t.Size_)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	return x
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>func convT16(val uint16) (x unsafe.Pointer) {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	if val &lt; uint16(len(staticuint64s)) {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		x = unsafe.Pointer(&amp;staticuint64s[val])
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		if goarch.BigEndian {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			x = add(x, 6)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	} else {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		x = mallocgc(2, uint16Type, false)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		*(*uint16)(x) = val
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	return
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>func convT32(val uint32) (x unsafe.Pointer) {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	if val &lt; uint32(len(staticuint64s)) {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		x = unsafe.Pointer(&amp;staticuint64s[val])
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		if goarch.BigEndian {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			x = add(x, 4)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	} else {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		x = mallocgc(4, uint32Type, false)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		*(*uint32)(x) = val
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	return
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>func convT64(val uint64) (x unsafe.Pointer) {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	if val &lt; uint64(len(staticuint64s)) {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		x = unsafe.Pointer(&amp;staticuint64s[val])
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	} else {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		x = mallocgc(8, uint64Type, false)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		*(*uint64)(x) = val
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	return
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>func convTstring(val string) (x unsafe.Pointer) {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	if val == &#34;&#34; {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		x = unsafe.Pointer(&amp;zeroVal[0])
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	} else {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		x = mallocgc(unsafe.Sizeof(val), stringType, true)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		*(*string)(x) = val
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	return
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>func convTslice(val []byte) (x unsafe.Pointer) {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	<span class="comment">// Note: this must work for any element type, not just byte.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	if (*slice)(unsafe.Pointer(&amp;val)).array == nil {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		x = unsafe.Pointer(&amp;zeroVal[0])
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	} else {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		x = mallocgc(unsafe.Sizeof(val), sliceType, true)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		*(*[]byte)(x) = val
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	return
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>func assertE2I(inter *interfacetype, t *_type) *itab {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	if t == nil {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		<span class="comment">// explicit conversions require non-nil interface value.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		panic(&amp;TypeAssertionError{nil, nil, &amp;inter.Type, &#34;&#34;})
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	return getitab(inter, t, false)
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>func assertE2I2(inter *interfacetype, t *_type) *itab {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	if t == nil {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		return nil
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	return getitab(inter, t, true)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span><span class="comment">// typeAssert builds an itab for the concrete type t and the</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span><span class="comment">// interface type s.Inter. If the conversion is not possible it</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span><span class="comment">// panics if s.CanFail is false and returns nil if s.CanFail is true.</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>func typeAssert(s *abi.TypeAssert, t *_type) *itab {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	var tab *itab
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	if t == nil {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		if !s.CanFail {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>			panic(&amp;TypeAssertionError{nil, nil, &amp;s.Inter.Type, &#34;&#34;})
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	} else {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		tab = getitab(s.Inter, t, s.CanFail)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	if !abi.UseInterfaceSwitchCache(GOARCH) {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		return tab
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	}
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// Maybe update the cache, so the next time the generated code</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	<span class="comment">// doesn&#39;t need to call into the runtime.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	if cheaprand()&amp;1023 != 0 {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		<span class="comment">// Only bother updating the cache ~1 in 1000 times.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		return tab
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	<span class="comment">// Load the current cache.</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	oldC := (*abi.TypeAssertCache)(atomic.Loadp(unsafe.Pointer(&amp;s.Cache)))
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	if cheaprand()&amp;uint32(oldC.Mask) != 0 {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		<span class="comment">// As cache gets larger, choose to update it less often</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		<span class="comment">// so we can amortize the cost of building a new cache.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		return tab
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	<span class="comment">// Make a new cache.</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	newC := buildTypeAssertCache(oldC, t, tab)
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	<span class="comment">// Update cache. Use compare-and-swap so if multiple threads</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	<span class="comment">// are fighting to update the cache, at least one of their</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	<span class="comment">// updates will stick.</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	atomic_casPointer((*unsafe.Pointer)(unsafe.Pointer(&amp;s.Cache)), unsafe.Pointer(oldC), unsafe.Pointer(newC))
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	return tab
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>func buildTypeAssertCache(oldC *abi.TypeAssertCache, typ *_type, tab *itab) *abi.TypeAssertCache {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	oldEntries := unsafe.Slice(&amp;oldC.Entries[0], oldC.Mask+1)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	<span class="comment">// Count the number of entries we need.</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	n := 1
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	for _, e := range oldEntries {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		if e.Typ != 0 {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			n++
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	<span class="comment">// Figure out how big a table we need.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	<span class="comment">// We need at least one more slot than the number of entries</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	<span class="comment">// so that we are guaranteed an empty slot (for termination).</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	newN := n * 2                         <span class="comment">// make it at most 50% full</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	newN = 1 &lt;&lt; sys.Len64(uint64(newN-1)) <span class="comment">// round up to a power of 2</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	<span class="comment">// Allocate the new table.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	newSize := unsafe.Sizeof(abi.TypeAssertCache{}) + uintptr(newN-1)*unsafe.Sizeof(abi.TypeAssertCacheEntry{})
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	newC := (*abi.TypeAssertCache)(mallocgc(newSize, nil, true))
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	newC.Mask = uintptr(newN - 1)
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	newEntries := unsafe.Slice(&amp;newC.Entries[0], newN)
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	<span class="comment">// Fill the new table.</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	addEntry := func(typ *_type, tab *itab) {
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		h := int(typ.Hash) &amp; (newN - 1)
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		for {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			if newEntries[h].Typ == 0 {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>				newEntries[h].Typ = uintptr(unsafe.Pointer(typ))
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>				newEntries[h].Itab = uintptr(unsafe.Pointer(tab))
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>				return
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			h = (h + 1) &amp; (newN - 1)
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		}
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	}
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	for _, e := range oldEntries {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		if e.Typ != 0 {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>			addEntry((*_type)(unsafe.Pointer(e.Typ)), (*itab)(unsafe.Pointer(e.Itab)))
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	addEntry(typ, tab)
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	return newC
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span><span class="comment">// Empty type assert cache. Contains one entry with a nil Typ (which</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span><span class="comment">// causes a cache lookup to fail immediately.)</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>var emptyTypeAssertCache = abi.TypeAssertCache{Mask: 0}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span><span class="comment">// interfaceSwitch compares t against the list of cases in s.</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span><span class="comment">// If t matches case i, interfaceSwitch returns the case index i and</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span><span class="comment">// an itab for the pair &lt;t, s.Cases[i]&gt;.</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span><span class="comment">// If there is no match, return N,nil, where N is the number</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span><span class="comment">// of cases.</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>func interfaceSwitch(s *abi.InterfaceSwitch, t *_type) (int, *itab) {
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	cases := unsafe.Slice(&amp;s.Cases[0], s.NCases)
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	<span class="comment">// Results if we don&#39;t find a match.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	case_ := len(cases)
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	var tab *itab
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	<span class="comment">// Look through each case in order.</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	for i, c := range cases {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		tab = getitab(c, t, true)
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		if tab != nil {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>			case_ = i
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>			break
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	if !abi.UseInterfaceSwitchCache(GOARCH) {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		return case_, tab
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	<span class="comment">// Maybe update the cache, so the next time the generated code</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	<span class="comment">// doesn&#39;t need to call into the runtime.</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	if cheaprand()&amp;1023 != 0 {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		<span class="comment">// Only bother updating the cache ~1 in 1000 times.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		<span class="comment">// This ensures we don&#39;t waste memory on switches, or</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		<span class="comment">// switch arguments, that only happen a few times.</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		return case_, tab
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	<span class="comment">// Load the current cache.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	oldC := (*abi.InterfaceSwitchCache)(atomic.Loadp(unsafe.Pointer(&amp;s.Cache)))
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	if cheaprand()&amp;uint32(oldC.Mask) != 0 {
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		<span class="comment">// As cache gets larger, choose to update it less often</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		<span class="comment">// so we can amortize the cost of building a new cache</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		<span class="comment">// (that cost is linear in oldc.Mask).</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		return case_, tab
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	<span class="comment">// Make a new cache.</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	newC := buildInterfaceSwitchCache(oldC, t, case_, tab)
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	<span class="comment">// Update cache. Use compare-and-swap so if multiple threads</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	<span class="comment">// are fighting to update the cache, at least one of their</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	<span class="comment">// updates will stick.</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	atomic_casPointer((*unsafe.Pointer)(unsafe.Pointer(&amp;s.Cache)), unsafe.Pointer(oldC), unsafe.Pointer(newC))
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	return case_, tab
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span><span class="comment">// buildInterfaceSwitchCache constructs an interface switch cache</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span><span class="comment">// containing all the entries from oldC plus the new entry</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span><span class="comment">// (typ,case_,tab).</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>func buildInterfaceSwitchCache(oldC *abi.InterfaceSwitchCache, typ *_type, case_ int, tab *itab) *abi.InterfaceSwitchCache {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	oldEntries := unsafe.Slice(&amp;oldC.Entries[0], oldC.Mask+1)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	<span class="comment">// Count the number of entries we need.</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	n := 1
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	for _, e := range oldEntries {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		if e.Typ != 0 {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>			n++
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		}
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	}
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	<span class="comment">// Figure out how big a table we need.</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	<span class="comment">// We need at least one more slot than the number of entries</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	<span class="comment">// so that we are guaranteed an empty slot (for termination).</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	newN := n * 2                         <span class="comment">// make it at most 50% full</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	newN = 1 &lt;&lt; sys.Len64(uint64(newN-1)) <span class="comment">// round up to a power of 2</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	<span class="comment">// Allocate the new table.</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	newSize := unsafe.Sizeof(abi.InterfaceSwitchCache{}) + uintptr(newN-1)*unsafe.Sizeof(abi.InterfaceSwitchCacheEntry{})
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	newC := (*abi.InterfaceSwitchCache)(mallocgc(newSize, nil, true))
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	newC.Mask = uintptr(newN - 1)
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	newEntries := unsafe.Slice(&amp;newC.Entries[0], newN)
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	<span class="comment">// Fill the new table.</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	addEntry := func(typ *_type, case_ int, tab *itab) {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		h := int(typ.Hash) &amp; (newN - 1)
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		for {
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>			if newEntries[h].Typ == 0 {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>				newEntries[h].Typ = uintptr(unsafe.Pointer(typ))
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>				newEntries[h].Case = case_
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>				newEntries[h].Itab = uintptr(unsafe.Pointer(tab))
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>				return
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>			h = (h + 1) &amp; (newN - 1)
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		}
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	for _, e := range oldEntries {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		if e.Typ != 0 {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>			addEntry((*_type)(unsafe.Pointer(e.Typ)), e.Case, (*itab)(unsafe.Pointer(e.Itab)))
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	}
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	addEntry(typ, case_, tab)
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	return newC
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>}
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span><span class="comment">// Empty interface switch cache. Contains one entry with a nil Typ (which</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span><span class="comment">// causes a cache lookup to fail immediately.)</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>var emptyInterfaceSwitchCache = abi.InterfaceSwitchCache{Mask: 0}
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_ifaceE2I reflect.ifaceE2I</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>func reflect_ifaceE2I(inter *interfacetype, e eface, dst *iface) {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	*dst = iface{assertE2I(inter, e._type), e.data}
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span><span class="comment">//go:linkname reflectlite_ifaceE2I internal/reflectlite.ifaceE2I</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>func reflectlite_ifaceE2I(inter *interfacetype, e eface, dst *iface) {
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	*dst = iface{assertE2I(inter, e._type), e.data}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>}
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>func iterate_itabs(fn func(*itab)) {
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	<span class="comment">// Note: only runs during stop the world or with itabLock held,</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	<span class="comment">// so no other locks/atomics needed.</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	t := itabTable
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	for i := uintptr(0); i &lt; t.size; i++ {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		m := *(**itab)(add(unsafe.Pointer(&amp;t.entries), i*goarch.PtrSize))
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		if m != nil {
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			fn(m)
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span><span class="comment">// staticuint64s is used to avoid allocating in convTx for small integer values.</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>var staticuint64s = [...]uint64{
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57,
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67,
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77,
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f,
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97,
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7,
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf,
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7,
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf,
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7,
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf,
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	0xd0, 0xd1, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7,
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	0xd8, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde, 0xdf,
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	0xe0, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7,
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	0xe8, 0xe9, 0xea, 0xeb, 0xec, 0xed, 0xee, 0xef,
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7,
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff,
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>}
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span><span class="comment">// The linker redirects a reference of a method that it determined</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span><span class="comment">// unreachable to a reference to this function, so it will throw if</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span><span class="comment">// ever called.</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>func unreachableMethod() {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	throw(&#34;unreachable method called. linker bug?&#34;)
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>
</pre><p><a href="iface.go?m=text">View as plain text</a></p>

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
