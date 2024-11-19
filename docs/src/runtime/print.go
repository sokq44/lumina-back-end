<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/print.go - Go Documentation Server</title>

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
<a href="print.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">print.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// The compiler knows that a print of a value of this type</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// should use printhex instead of printuint (decimal).</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>type hex uint64
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>func bytes(s string) (ret []byte) {
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	rp := (*slice)(unsafe.Pointer(&amp;ret))
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	sp := stringStructOf(&amp;s)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	rp.array = sp.str
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	rp.len = sp.len
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	rp.cap = sp.len
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	return
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>}
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>var (
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// printBacklog is a circular buffer of messages written with the builtin</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// print* functions, for use in postmortem analysis of core dumps.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	printBacklog      [512]byte
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	printBacklogIndex int
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// recordForPanic maintains a circular buffer of messages written by the</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// runtime leading up to a process crash, allowing the messages to be</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// extracted from a core dump.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// The text written during a process crash (following &#34;panic&#34; or &#34;fatal</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// error&#34;) is not saved, since the goroutine stacks will generally be readable</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// from the runtime data structures in the core file.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>func recordForPanic(b []byte) {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	printlock()
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	if panicking.Load() == 0 {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		<span class="comment">// Not actively crashing: maintain circular buffer of print output.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		for i := 0; i &lt; len(b); {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>			n := copy(printBacklog[printBacklogIndex:], b[i:])
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			i += n
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			printBacklogIndex += n
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			printBacklogIndex %= len(printBacklog)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	printunlock()
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>var debuglock mutex
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// The compiler emits calls to printlock and printunlock around</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// the multiple calls that implement a single Go print or println</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// statement. Some of the print helpers (printslice, for example)</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// call print recursively. There is also the problem of a crash</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// happening during the print routines and needing to acquire</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// the print lock to print information about the crash.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// For both these reasons, let a thread acquire the printlock &#39;recursively&#39;.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>func printlock() {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	mp := getg().m
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	mp.locks++ <span class="comment">// do not reschedule between printlock++ and lock(&amp;debuglock).</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	mp.printlock++
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	if mp.printlock == 1 {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		lock(&amp;debuglock)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	mp.locks-- <span class="comment">// now we know debuglock is held and holding up mp.locks for us.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>func printunlock() {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	mp := getg().m
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	mp.printlock--
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	if mp.printlock == 0 {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		unlock(&amp;debuglock)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// write to goroutine-local buffer if diverting output,</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// or else standard error.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>func gwrite(b []byte) {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	if len(b) == 0 {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		return
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	recordForPanic(b)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	gp := getg()
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t use the writebuf if gp.m is dying. We want anything</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// written through gwrite to appear in the terminal rather</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// than be written to in some buffer, if we&#39;re in a panicking state.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// Note that we can&#39;t just clear writebuf in the gp.m.dying case</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// because a panic isn&#39;t allowed to have any write barriers.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if gp == nil || gp.writebuf == nil || gp.m.dying &gt; 0 {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		writeErr(b)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		return
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	n := copy(gp.writebuf[len(gp.writebuf):cap(gp.writebuf)], b)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	gp.writebuf = gp.writebuf[:len(gp.writebuf)+n]
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func printsp() {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	printstring(&#34; &#34;)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>func printnl() {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	printstring(&#34;\n&#34;)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func printbool(v bool) {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	if v {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		printstring(&#34;true&#34;)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	} else {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		printstring(&#34;false&#34;)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>func printfloat(v float64) {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	switch {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	case v != v:
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		printstring(&#34;NaN&#34;)
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		return
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	case v+v == v &amp;&amp; v &gt; 0:
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		printstring(&#34;+Inf&#34;)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		return
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	case v+v == v &amp;&amp; v &lt; 0:
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		printstring(&#34;-Inf&#34;)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		return
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	const n = 7 <span class="comment">// digits printed</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	var buf [n + 7]byte
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	buf[0] = &#39;+&#39;
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	e := 0 <span class="comment">// exp</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	if v == 0 {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		if 1/v &lt; 0 {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			buf[0] = &#39;-&#39;
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	} else {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		if v &lt; 0 {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			v = -v
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			buf[0] = &#39;-&#39;
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		<span class="comment">// normalize</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		for v &gt;= 10 {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			e++
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			v /= 10
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		for v &lt; 1 {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			e--
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			v *= 10
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// round</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		h := 5.0
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		for i := 0; i &lt; n; i++ {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			h /= 10
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		v += h
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		if v &gt;= 10 {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			e++
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			v /= 10
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// format +d.dddd+edd</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		s := int(v)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		buf[i+2] = byte(s + &#39;0&#39;)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		v -= float64(s)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		v *= 10
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	buf[1] = buf[2]
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	buf[2] = &#39;.&#39;
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	buf[n+2] = &#39;e&#39;
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	buf[n+3] = &#39;+&#39;
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	if e &lt; 0 {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		e = -e
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		buf[n+3] = &#39;-&#39;
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	buf[n+4] = byte(e/100) + &#39;0&#39;
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	buf[n+5] = byte(e/10)%10 + &#39;0&#39;
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	buf[n+6] = byte(e%10) + &#39;0&#39;
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	gwrite(buf[:])
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>func printcomplex(c complex128) {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	print(&#34;(&#34;, real(c), imag(c), &#34;i)&#34;)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>func printuint(v uint64) {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	var buf [100]byte
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	i := len(buf)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	for i--; i &gt; 0; i-- {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		buf[i] = byte(v%10 + &#39;0&#39;)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		if v &lt; 10 {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			break
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		v /= 10
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	gwrite(buf[i:])
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>func printint(v int64) {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if v &lt; 0 {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		printstring(&#34;-&#34;)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		v = -v
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	printuint(uint64(v))
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>var minhexdigits = 0 <span class="comment">// protected by printlock</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>func printhex(v uint64) {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	const dig = &#34;0123456789abcdef&#34;
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	var buf [100]byte
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	i := len(buf)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	for i--; i &gt; 0; i-- {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		buf[i] = dig[v%16]
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		if v &lt; 16 &amp;&amp; len(buf)-i &gt;= minhexdigits {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			break
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		v /= 16
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	i--
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	buf[i] = &#39;x&#39;
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	i--
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	buf[i] = &#39;0&#39;
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	gwrite(buf[i:])
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>func printpointer(p unsafe.Pointer) {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	printhex(uint64(uintptr(p)))
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>func printuintptr(p uintptr) {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	printhex(uint64(p))
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>func printstring(s string) {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	gwrite(bytes(s))
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>func printslice(s []byte) {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	sp := (*slice)(unsafe.Pointer(&amp;s))
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	print(&#34;[&#34;, len(s), &#34;/&#34;, cap(s), &#34;]&#34;)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	printpointer(sp.array)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>func printeface(e eface) {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	print(&#34;(&#34;, e._type, &#34;,&#34;, e.data, &#34;)&#34;)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>func printiface(i iface) {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	print(&#34;(&#34;, i.tab, &#34;,&#34;, i.data, &#34;)&#34;)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">// hexdumpWords prints a word-oriented hex dump of [p, end).</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">// If mark != nil, it will be called with each printed word&#39;s address</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">// and should return a character mark to appear just before that</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span><span class="comment">// word&#39;s value. It can return 0 to indicate no mark.</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>func hexdumpWords(p, end uintptr, mark func(uintptr) byte) {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	printlock()
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	var markbuf [1]byte
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	markbuf[0] = &#39; &#39;
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	minhexdigits = int(unsafe.Sizeof(uintptr(0)) * 2)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	for i := uintptr(0); p+i &lt; end; i += goarch.PtrSize {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		if i%16 == 0 {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			if i != 0 {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>				println()
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			print(hex(p+i), &#34;: &#34;)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		if mark != nil {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			markbuf[0] = mark(p + i)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			if markbuf[0] == 0 {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>				markbuf[0] = &#39; &#39;
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		gwrite(markbuf[:])
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		val := *(*uintptr)(unsafe.Pointer(p + i))
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		print(hex(val))
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		print(&#34; &#34;)
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		<span class="comment">// Can we symbolize val?</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		fn := findfunc(val)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		if fn.valid() {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			print(&#34;&lt;&#34;, funcname(fn), &#34;+&#34;, hex(val-fn.entry()), &#34;&gt; &#34;)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	minhexdigits = 0
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	println()
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	printunlock()
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
</pre><p><a href="print.go?m=text">View as plain text</a></p>

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
