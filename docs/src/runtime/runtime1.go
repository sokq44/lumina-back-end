<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/runtime1.go - Go Documentation Server</title>

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
<a href="runtime1.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">runtime1.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/bytealg&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// Keep a cached value to make gotraceback fast,</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// since we call it on every call to gentraceback.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// The cached value is a uint32 in which the low bits</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// are the &#34;crash&#34; and &#34;all&#34; settings and the remaining</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// bits are the traceback value (0 off, 1 on, 2 include system).</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>const (
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	tracebackCrash = 1 &lt;&lt; iota
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	tracebackAll
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	tracebackShift = iota
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>var traceback_cache uint32 = 2 &lt;&lt; tracebackShift
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>var traceback_env uint32
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// gotraceback returns the current traceback settings.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// If level is 0, suppress all tracebacks.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// If level is 1, show tracebacks, but exclude runtime frames.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// If level is 2, show tracebacks including runtime frames.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// If all is set, print all goroutine stacks. Otherwise, print just the current goroutine.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// If crash is set, crash (core dump, etc) after tracebacking.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>func gotraceback() (level int32, all, crash bool) {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	gp := getg()
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	t := atomic.Load(&amp;traceback_cache)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	crash = t&amp;tracebackCrash != 0
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	all = gp.m.throwing &gt;= throwTypeUser || t&amp;tracebackAll != 0
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	if gp.m.traceback != 0 {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		level = int32(gp.m.traceback)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	} else if gp.m.throwing &gt;= throwTypeRuntime {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		<span class="comment">// Always include runtime frames in runtime throws unless</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		<span class="comment">// otherwise overridden by m.traceback.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		level = 2
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	} else {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		level = int32(t &gt;&gt; tracebackShift)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	return
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>var (
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	argc int32
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	argv **byte
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// nosplit for use in linux startup sysargs.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>func argv_index(argv **byte, i int32) *byte {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	return *(**byte)(add(unsafe.Pointer(argv), uintptr(i)*goarch.PtrSize))
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>func args(c int32, v **byte) {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	argc = c
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	argv = v
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	sysargs(c, v)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func goargs() {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	if GOOS == &#34;windows&#34; {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		return
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	argslice = make([]string, argc)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	for i := int32(0); i &lt; argc; i++ {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		argslice[i] = gostringnocopy(argv_index(argv, i))
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>func goenvs_unix() {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// TODO(austin): ppc64 in dynamic linking mode doesn&#39;t</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// guarantee env[] will immediately follow argv. Might cause</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// problems.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	n := int32(0)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	for argv_index(argv, argc+1+n) != nil {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		n++
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	envs = make([]string, n)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	for i := int32(0); i &lt; n; i++ {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		envs[i] = gostring(argv_index(argv, argc+1+i))
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func environ() []string {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	return envs
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// TODO: These should be locals in testAtomic64, but we don&#39;t 8-byte</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// align stack variables on 386.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>var test_z64, test_x64 uint64
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func testAtomic64() {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	test_z64 = 42
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	test_x64 = 0
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	if atomic.Cas64(&amp;test_z64, test_x64, 1) {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		throw(&#34;cas64 failed&#34;)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	if test_x64 != 0 {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		throw(&#34;cas64 failed&#34;)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	test_x64 = 42
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if !atomic.Cas64(&amp;test_z64, test_x64, 1) {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		throw(&#34;cas64 failed&#34;)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	if test_x64 != 42 || test_z64 != 1 {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		throw(&#34;cas64 failed&#34;)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	if atomic.Load64(&amp;test_z64) != 1 {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		throw(&#34;load64 failed&#34;)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	atomic.Store64(&amp;test_z64, (1&lt;&lt;40)+1)
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	if atomic.Load64(&amp;test_z64) != (1&lt;&lt;40)+1 {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		throw(&#34;store64 failed&#34;)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	if atomic.Xadd64(&amp;test_z64, (1&lt;&lt;40)+1) != (2&lt;&lt;40)+2 {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		throw(&#34;xadd64 failed&#34;)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	if atomic.Load64(&amp;test_z64) != (2&lt;&lt;40)+2 {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		throw(&#34;xadd64 failed&#34;)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	if atomic.Xchg64(&amp;test_z64, (3&lt;&lt;40)+3) != (2&lt;&lt;40)+2 {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		throw(&#34;xchg64 failed&#34;)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	if atomic.Load64(&amp;test_z64) != (3&lt;&lt;40)+3 {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		throw(&#34;xchg64 failed&#34;)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>func check() {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	var (
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		a     int8
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		b     uint8
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		c     int16
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		d     uint16
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		e     int32
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		f     uint32
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		g     int64
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		h     uint64
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		i, i1 float32
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		j, j1 float64
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		k     unsafe.Pointer
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		l     *uint16
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		m     [4]byte
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	type x1t struct {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		x uint8
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	type y1t struct {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		x1 x1t
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		y  uint8
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	var x1 x1t
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	var y1 y1t
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	if unsafe.Sizeof(a) != 1 {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		throw(&#34;bad a&#34;)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	if unsafe.Sizeof(b) != 1 {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		throw(&#34;bad b&#34;)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	if unsafe.Sizeof(c) != 2 {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		throw(&#34;bad c&#34;)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	if unsafe.Sizeof(d) != 2 {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		throw(&#34;bad d&#34;)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	if unsafe.Sizeof(e) != 4 {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		throw(&#34;bad e&#34;)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	if unsafe.Sizeof(f) != 4 {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		throw(&#34;bad f&#34;)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	if unsafe.Sizeof(g) != 8 {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		throw(&#34;bad g&#34;)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	if unsafe.Sizeof(h) != 8 {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		throw(&#34;bad h&#34;)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	if unsafe.Sizeof(i) != 4 {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		throw(&#34;bad i&#34;)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	if unsafe.Sizeof(j) != 8 {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		throw(&#34;bad j&#34;)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	if unsafe.Sizeof(k) != goarch.PtrSize {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		throw(&#34;bad k&#34;)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	if unsafe.Sizeof(l) != goarch.PtrSize {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		throw(&#34;bad l&#34;)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	if unsafe.Sizeof(x1) != 1 {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		throw(&#34;bad unsafe.Sizeof x1&#34;)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	if unsafe.Offsetof(y1.y) != 1 {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		throw(&#34;bad offsetof y1.y&#34;)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	if unsafe.Sizeof(y1) != 2 {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		throw(&#34;bad unsafe.Sizeof y1&#34;)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	if timediv(12345*1000000000+54321, 1000000000, &amp;e) != 12345 || e != 54321 {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		throw(&#34;bad timediv&#34;)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	var z uint32
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	z = 1
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	if !atomic.Cas(&amp;z, 1, 2) {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		throw(&#34;cas1&#34;)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	if z != 2 {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		throw(&#34;cas2&#34;)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	z = 4
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	if atomic.Cas(&amp;z, 5, 6) {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		throw(&#34;cas3&#34;)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	if z != 4 {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		throw(&#34;cas4&#34;)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	z = 0xffffffff
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	if !atomic.Cas(&amp;z, 0xffffffff, 0xfffffffe) {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		throw(&#34;cas5&#34;)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	if z != 0xfffffffe {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		throw(&#34;cas6&#34;)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	m = [4]byte{1, 1, 1, 1}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	atomic.Or8(&amp;m[1], 0xf0)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	if m[0] != 1 || m[1] != 0xf1 || m[2] != 1 || m[3] != 1 {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		throw(&#34;atomicor8&#34;)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	m = [4]byte{0xff, 0xff, 0xff, 0xff}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	atomic.And8(&amp;m[1], 0x1)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	if m[0] != 0xff || m[1] != 0x1 || m[2] != 0xff || m[3] != 0xff {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		throw(&#34;atomicand8&#34;)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	*(*uint64)(unsafe.Pointer(&amp;j)) = ^uint64(0)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	if j == j {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		throw(&#34;float64nan&#34;)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	if !(j != j) {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		throw(&#34;float64nan1&#34;)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	*(*uint64)(unsafe.Pointer(&amp;j1)) = ^uint64(1)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	if j == j1 {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		throw(&#34;float64nan2&#34;)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	if !(j != j1) {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		throw(&#34;float64nan3&#34;)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	*(*uint32)(unsafe.Pointer(&amp;i)) = ^uint32(0)
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	if i == i {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		throw(&#34;float32nan&#34;)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	if i == i {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		throw(&#34;float32nan1&#34;)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	*(*uint32)(unsafe.Pointer(&amp;i1)) = ^uint32(1)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	if i == i1 {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		throw(&#34;float32nan2&#34;)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	if i == i1 {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		throw(&#34;float32nan3&#34;)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	testAtomic64()
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	if fixedStack != round2(fixedStack) {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		throw(&#34;FixedStack is not power-of-2&#34;)
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	if !checkASM() {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		throw(&#34;assembly checks failed&#34;)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>type dbgVar struct {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	name   string
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	value  *int32        <span class="comment">// for variables that can only be set at startup</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	atomic *atomic.Int32 <span class="comment">// for variables that can be changed during execution</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	def    int32         <span class="comment">// default value (ideally zero)</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">// Holds variables parsed from GODEBUG env var,</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">// except for &#34;memprofilerate&#34; since there is an</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">// existing int var for that value, which may</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// already have an initial value.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>var debug struct {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	cgocheck                int32
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	clobberfree             int32
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	disablethp              int32
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	dontfreezetheworld      int32
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	efence                  int32
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	gccheckmark             int32
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	gcpacertrace            int32
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	gcshrinkstackoff        int32
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	gcstoptheworld          int32
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	gctrace                 int32
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	invalidptr              int32
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	madvdontneed            int32 <span class="comment">// for Linux; issue 28466</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	runtimeContentionStacks atomic.Int32
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	scavtrace               int32
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	scheddetail             int32
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	schedtrace              int32
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	tracebackancestors      int32
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	asyncpreemptoff         int32
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	harddecommit            int32
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	adaptivestackstart      int32
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	tracefpunwindoff        int32
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	traceadvanceperiod      int32
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// debug.malloc is used as a combined debug check</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// in the malloc function and should be set</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// if any of the below debug options is != 0.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	malloc         bool
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	allocfreetrace int32
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	inittrace      int32
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	sbrk           int32
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	panicnil atomic.Int32
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>var dbgvars = []*dbgVar{
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	{name: &#34;allocfreetrace&#34;, value: &amp;debug.allocfreetrace},
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	{name: &#34;clobberfree&#34;, value: &amp;debug.clobberfree},
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	{name: &#34;cgocheck&#34;, value: &amp;debug.cgocheck},
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	{name: &#34;disablethp&#34;, value: &amp;debug.disablethp},
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	{name: &#34;dontfreezetheworld&#34;, value: &amp;debug.dontfreezetheworld},
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	{name: &#34;efence&#34;, value: &amp;debug.efence},
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	{name: &#34;gccheckmark&#34;, value: &amp;debug.gccheckmark},
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	{name: &#34;gcpacertrace&#34;, value: &amp;debug.gcpacertrace},
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	{name: &#34;gcshrinkstackoff&#34;, value: &amp;debug.gcshrinkstackoff},
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	{name: &#34;gcstoptheworld&#34;, value: &amp;debug.gcstoptheworld},
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	{name: &#34;gctrace&#34;, value: &amp;debug.gctrace},
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	{name: &#34;invalidptr&#34;, value: &amp;debug.invalidptr},
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	{name: &#34;madvdontneed&#34;, value: &amp;debug.madvdontneed},
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	{name: &#34;runtimecontentionstacks&#34;, atomic: &amp;debug.runtimeContentionStacks},
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	{name: &#34;sbrk&#34;, value: &amp;debug.sbrk},
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	{name: &#34;scavtrace&#34;, value: &amp;debug.scavtrace},
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	{name: &#34;scheddetail&#34;, value: &amp;debug.scheddetail},
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	{name: &#34;schedtrace&#34;, value: &amp;debug.schedtrace},
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	{name: &#34;tracebackancestors&#34;, value: &amp;debug.tracebackancestors},
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	{name: &#34;asyncpreemptoff&#34;, value: &amp;debug.asyncpreemptoff},
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	{name: &#34;inittrace&#34;, value: &amp;debug.inittrace},
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	{name: &#34;harddecommit&#34;, value: &amp;debug.harddecommit},
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	{name: &#34;adaptivestackstart&#34;, value: &amp;debug.adaptivestackstart},
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	{name: &#34;tracefpunwindoff&#34;, value: &amp;debug.tracefpunwindoff},
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	{name: &#34;panicnil&#34;, atomic: &amp;debug.panicnil},
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	{name: &#34;traceadvanceperiod&#34;, value: &amp;debug.traceadvanceperiod},
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>func parsedebugvars() {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	<span class="comment">// defaults</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	debug.cgocheck = 1
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	debug.invalidptr = 1
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	debug.adaptivestackstart = 1 <span class="comment">// set this to 0 to turn larger initial goroutine stacks off</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	if GOOS == &#34;linux&#34; {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		<span class="comment">// On Linux, MADV_FREE is faster than MADV_DONTNEED,</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		<span class="comment">// but doesn&#39;t affect many of the statistics that</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		<span class="comment">// MADV_DONTNEED does until the memory is actually</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		<span class="comment">// reclaimed. This generally leads to poor user</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		<span class="comment">// experience, like confusing stats in top and other</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		<span class="comment">// monitoring tools; and bad integration with</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		<span class="comment">// management systems that respond to memory usage.</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		<span class="comment">// Hence, default to MADV_DONTNEED.</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		debug.madvdontneed = 1
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	debug.traceadvanceperiod = defaultTraceAdvancePeriod
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	godebug := gogetenv(&#34;GODEBUG&#34;)
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	p := new(string)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	*p = godebug
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	godebugEnv.Store(p)
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	<span class="comment">// apply runtime defaults, if any</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	for _, v := range dbgvars {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		if v.def != 0 {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			<span class="comment">// Every var should have either v.value or v.atomic set.</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			if v.value != nil {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>				*v.value = v.def
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			} else if v.atomic != nil {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>				v.atomic.Store(v.def)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	<span class="comment">// apply compile-time GODEBUG settings</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	parsegodebug(godebugDefault, nil)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	<span class="comment">// apply environment settings</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	parsegodebug(godebug, nil)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	debug.malloc = (debug.allocfreetrace | debug.inittrace | debug.sbrk) != 0
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	setTraceback(gogetenv(&#34;GOTRACEBACK&#34;))
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	traceback_env = traceback_cache
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>}
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span><span class="comment">// reparsedebugvars reparses the runtime&#39;s debug variables</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span><span class="comment">// because the environment variable has been changed to env.</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>func reparsedebugvars(env string) {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	seen := make(map[string]bool)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	<span class="comment">// apply environment settings</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	parsegodebug(env, seen)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	<span class="comment">// apply compile-time GODEBUG settings for as-yet-unseen variables</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	parsegodebug(godebugDefault, seen)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	<span class="comment">// apply defaults for as-yet-unseen variables</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	for _, v := range dbgvars {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		if v.atomic != nil &amp;&amp; !seen[v.name] {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			v.atomic.Store(0)
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span><span class="comment">// parsegodebug parses the godebug string, updating variables listed in dbgvars.</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span><span class="comment">// If seen == nil, this is startup time and we process the string left to right</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span><span class="comment">// overwriting older settings with newer ones.</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span><span class="comment">// If seen != nil, $GODEBUG has changed and we are doing an</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span><span class="comment">// incremental update. To avoid flapping in the case where a value is</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span><span class="comment">// set multiple times (perhaps in the default and the environment,</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span><span class="comment">// or perhaps twice in the environment), we process the string right-to-left</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span><span class="comment">// and only change values not already seen. After doing this for both</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span><span class="comment">// the environment and the default settings, the caller must also call</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span><span class="comment">// cleargodebug(seen) to reset any now-unset values back to their defaults.</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>func parsegodebug(godebug string, seen map[string]bool) {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	for p := godebug; p != &#34;&#34;; {
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		var field string
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		if seen == nil {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			<span class="comment">// startup: process left to right, overwriting older settings with newer</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			i := bytealg.IndexByteString(p, &#39;,&#39;)
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			if i &lt; 0 {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>				field, p = p, &#34;&#34;
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			} else {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>				field, p = p[:i], p[i+1:]
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		} else {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>			<span class="comment">// incremental update: process right to left, updating and skipping seen</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			i := len(p) - 1
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			for i &gt;= 0 &amp;&amp; p[i] != &#39;,&#39; {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>				i--
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			if i &lt; 0 {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>				p, field = &#34;&#34;, p
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			} else {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>				p, field = p[:i], p[i+1:]
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>			}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		i := bytealg.IndexByteString(field, &#39;=&#39;)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		if i &lt; 0 {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			continue
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		key, value := field[:i], field[i+1:]
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		if seen[key] {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			continue
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		if seen != nil {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			seen[key] = true
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		<span class="comment">// Update MemProfileRate directly here since it</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		<span class="comment">// is int, not int32, and should only be updated</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		<span class="comment">// if specified in GODEBUG.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		if seen == nil &amp;&amp; key == &#34;memprofilerate&#34; {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			if n, ok := atoi(value); ok {
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>				MemProfileRate = n
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			}
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		} else {
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>			for _, v := range dbgvars {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>				if v.name == key {
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>					if n, ok := atoi32(value); ok {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>						if seen == nil &amp;&amp; v.value != nil {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>							*v.value = n
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>						} else if v.atomic != nil {
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>							v.atomic.Store(n)
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>						}
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>					}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>				}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>			}
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		}
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	if debug.cgocheck &gt; 1 {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		throw(&#34;cgocheck &gt; 1 mode is no longer supported at runtime. Use GOEXPERIMENT=cgocheck2 at build time instead.&#34;)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span><span class="comment">//go:linkname setTraceback runtime/debug.SetTraceback</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>func setTraceback(level string) {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	var t uint32
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	switch level {
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	case &#34;none&#34;:
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		t = 0
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	case &#34;single&#34;, &#34;&#34;:
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		t = 1 &lt;&lt; tracebackShift
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	case &#34;all&#34;:
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		t = 1&lt;&lt;tracebackShift | tracebackAll
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	case &#34;system&#34;:
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		t = 2&lt;&lt;tracebackShift | tracebackAll
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	case &#34;crash&#34;:
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		t = 2&lt;&lt;tracebackShift | tracebackAll | tracebackCrash
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	case &#34;wer&#34;:
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		if GOOS == &#34;windows&#34; {
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>			t = 2&lt;&lt;tracebackShift | tracebackAll | tracebackCrash
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>			enableWER()
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			break
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		}
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		fallthrough
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	default:
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		t = tracebackAll
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		if n, ok := atoi(level); ok &amp;&amp; n == int(uint32(n)) {
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>			t |= uint32(n) &lt;&lt; tracebackShift
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	<span class="comment">// when C owns the process, simply exit&#39;ing the process on fatal errors</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	<span class="comment">// and panics is surprising. Be louder and abort instead.</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	if islibrary || isarchive {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		t |= tracebackCrash
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	}
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	t |= traceback_env
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	atomic.Store(&amp;traceback_cache, t)
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span><span class="comment">// Poor mans 64-bit division.</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span><span class="comment">// This is a very special function, do not use it if you are not sure what you are doing.</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span><span class="comment">// int64 division is lowered into _divv() call on 386, which does not fit into nosplit functions.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span><span class="comment">// Handles overflow in a time-specific manner.</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span><span class="comment">// This keeps us within no-split stack limits on 32-bit processors.</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>func timediv(v int64, div int32, rem *int32) int32 {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	res := int32(0)
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	for bit := 30; bit &gt;= 0; bit-- {
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		if v &gt;= int64(div)&lt;&lt;uint(bit) {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			v = v - (int64(div) &lt;&lt; uint(bit))
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>			<span class="comment">// Before this for loop, res was 0, thus all these</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			<span class="comment">// power of 2 increments are now just bitsets.</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>			res |= 1 &lt;&lt; uint(bit)
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	}
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	if v &gt;= int64(div) {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		if rem != nil {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>			*rem = 0
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		return 0x7fffffff
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	if rem != nil {
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		*rem = int32(v)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	}
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	return res
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span><span class="comment">// Helpers for Go. Must be NOSPLIT, must only call NOSPLIT functions, and must not block.</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>func acquirem() *m {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	gp := getg()
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	gp.m.locks++
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	return gp.m
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>}
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>func releasem(mp *m) {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	gp := getg()
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	mp.locks--
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	if mp.locks == 0 &amp;&amp; gp.preempt {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		<span class="comment">// restore the preemption request in case we&#39;ve cleared it in newstack</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		gp.stackguard0 = stackPreempt
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	}
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_typelinks reflect.typelinks</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>func reflect_typelinks() ([]unsafe.Pointer, [][]int32) {
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	modules := activeModules()
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	sections := []unsafe.Pointer{unsafe.Pointer(modules[0].types)}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	ret := [][]int32{modules[0].typelinks}
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	for _, md := range modules[1:] {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		sections = append(sections, unsafe.Pointer(md.types))
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		ret = append(ret, md.typelinks)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	}
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	return sections, ret
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span><span class="comment">// reflect_resolveNameOff resolves a name offset from a base pointer.</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_resolveNameOff reflect.resolveNameOff</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>func reflect_resolveNameOff(ptrInModule unsafe.Pointer, off int32) unsafe.Pointer {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	return unsafe.Pointer(resolveNameOff(ptrInModule, nameOff(off)).Bytes)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span><span class="comment">// reflect_resolveTypeOff resolves an *rtype offset from a base type.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_resolveTypeOff reflect.resolveTypeOff</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>func reflect_resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer {
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	return unsafe.Pointer(toRType((*_type)(rtype)).typeOff(typeOff(off)))
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>}
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span><span class="comment">// reflect_resolveTextOff resolves a function pointer offset from a base type.</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_resolveTextOff reflect.resolveTextOff</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>func reflect_resolveTextOff(rtype unsafe.Pointer, off int32) unsafe.Pointer {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	return toRType((*_type)(rtype)).textOff(textOff(off))
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span><span class="comment">// reflectlite_resolveNameOff resolves a name offset from a base pointer.</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span><span class="comment">//go:linkname reflectlite_resolveNameOff internal/reflectlite.resolveNameOff</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>func reflectlite_resolveNameOff(ptrInModule unsafe.Pointer, off int32) unsafe.Pointer {
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	return unsafe.Pointer(resolveNameOff(ptrInModule, nameOff(off)).Bytes)
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>}
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span><span class="comment">// reflectlite_resolveTypeOff resolves an *rtype offset from a base type.</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span><span class="comment">//go:linkname reflectlite_resolveTypeOff internal/reflectlite.resolveTypeOff</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>func reflectlite_resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer {
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	return unsafe.Pointer(toRType((*_type)(rtype)).typeOff(typeOff(off)))
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>}
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span><span class="comment">// reflect_addReflectOff adds a pointer to the reflection offset lookup map.</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_addReflectOff reflect.addReflectOff</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>func reflect_addReflectOff(ptr unsafe.Pointer) int32 {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	reflectOffsLock()
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	if reflectOffs.m == nil {
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		reflectOffs.m = make(map[int32]unsafe.Pointer)
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		reflectOffs.minv = make(map[unsafe.Pointer]int32)
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		reflectOffs.next = -1
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	id, found := reflectOffs.minv[ptr]
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	if !found {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		id = reflectOffs.next
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		reflectOffs.next-- <span class="comment">// use negative offsets as IDs to aid debugging</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		reflectOffs.m[id] = ptr
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		reflectOffs.minv[ptr] = id
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	reflectOffsUnlock()
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	return id
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>}
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>
</pre><p><a href="runtime1.go?m=text">View as plain text</a></p>

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
