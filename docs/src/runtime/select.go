<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/select.go - Go Documentation Server</title>

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
<a href="select.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">select.go</span>
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
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This file contains the implementation of Go select statements.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>const debugSelect = false
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// Select case descriptor.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Known to compiler.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// Changes here must also be made in src/cmd/compile/internal/walk/select.go&#39;s scasetype.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>type scase struct {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	c    *hchan         <span class="comment">// chan</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	elem unsafe.Pointer <span class="comment">// data element</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>var (
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	chansendpc = abi.FuncPCABIInternal(chansend)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	chanrecvpc = abi.FuncPCABIInternal(chanrecv)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>func selectsetpc(pc *uintptr) {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	*pc = getcallerpc()
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>func sellock(scases []scase, lockorder []uint16) {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	var c *hchan
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	for _, o := range lockorder {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		c0 := scases[o].c
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		if c0 != c {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>			c = c0
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>			lock(&amp;c.lock)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>}
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>func selunlock(scases []scase, lockorder []uint16) {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// We must be very careful here to not touch sel after we have unlocked</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// the last lock, because sel can be freed right after the last unlock.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// Consider the following situation.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// First M calls runtime·park() in runtime·selectgo() passing the sel.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// Once runtime·park() has unlocked the last lock, another M makes</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// the G that calls select runnable again and schedules it for execution.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// When the G runs on another M, it locks all the locks and frees sel.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// Now if the first M touches sel, it will access freed memory.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	for i := len(lockorder) - 1; i &gt;= 0; i-- {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		c := scases[lockorder[i]].c
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		if i &gt; 0 &amp;&amp; c == scases[lockorder[i-1]].c {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			continue <span class="comment">// will unlock it on the next iteration</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		unlock(&amp;c.lock)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>func selparkcommit(gp *g, _ unsafe.Pointer) bool {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// There are unlocked sudogs that point into gp&#39;s stack. Stack</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// copying must lock the channels of those sudogs.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// Set activeStackChans here instead of before we try parking</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// because we could self-deadlock in stack growth on a</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// channel lock.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	gp.activeStackChans = true
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// Mark that it&#39;s safe for stack shrinking to occur now,</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// because any thread acquiring this G&#39;s stack for shrinking</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// is guaranteed to observe activeStackChans after this store.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	gp.parkingOnChan.Store(false)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// Make sure we unlock after setting activeStackChans and</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// unsetting parkingOnChan. The moment we unlock any of the</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// channel locks we risk gp getting readied by a channel operation</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// and so gp could continue running before everything before the</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// unlock is visible (even to gp itself).</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// This must not access gp&#39;s stack (see gopark). In</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// particular, it must not access the *hselect. That&#39;s okay,</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// because by the time this is called, gp.waiting has all</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// channels in lock order.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	var lastc *hchan
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	for sg := gp.waiting; sg != nil; sg = sg.waitlink {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		if sg.c != lastc &amp;&amp; lastc != nil {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			<span class="comment">// As soon as we unlock the channel, fields in</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			<span class="comment">// any sudog with that channel may change,</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			<span class="comment">// including c and waitlink. Since multiple</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			<span class="comment">// sudogs may have the same channel, we unlock</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			<span class="comment">// only after we&#39;ve passed the last instance</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			<span class="comment">// of a channel.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			unlock(&amp;lastc.lock)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		lastc = sg.c
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if lastc != nil {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		unlock(&amp;lastc.lock)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	return true
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>func block() {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	gopark(nil, nil, waitReasonSelectNoCases, traceBlockForever, 1) <span class="comment">// forever</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// selectgo implements the select statement.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// cas0 points to an array of type [ncases]scase, and order0 points to</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// an array of type [2*ncases]uint16 where ncases must be &lt;= 65536.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// Both reside on the goroutine&#39;s stack (regardless of any escaping in</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// selectgo).</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// For race detector builds, pc0 points to an array of type</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">// [ncases]uintptr (also on the stack); for other builds, it&#39;s set to</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">// nil.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// selectgo returns the index of the chosen scase, which matches the</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// ordinal position of its respective select{recv,send,default} call.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// Also, if the chosen scase was a receive operation, it reports whether</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// a value was received.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>func selectgo(cas0 *scase, order0 *uint16, pc0 *uintptr, nsends, nrecvs int, block bool) (int, bool) {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	if debugSelect {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		print(&#34;select: cas0=&#34;, cas0, &#34;\n&#34;)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// NOTE: In order to maintain a lean stack size, the number of scases</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// is capped at 65536.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	cas1 := (*[1 &lt;&lt; 16]scase)(unsafe.Pointer(cas0))
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	order1 := (*[1 &lt;&lt; 17]uint16)(unsafe.Pointer(order0))
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	ncases := nsends + nrecvs
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	scases := cas1[:ncases:ncases]
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	pollorder := order1[:ncases:ncases]
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	lockorder := order1[ncases:][:ncases:ncases]
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// NOTE: pollorder/lockorder&#39;s underlying array was not zero-initialized by compiler.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// Even when raceenabled is true, there might be select</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// statements in packages compiled without -race (e.g.,</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// ensureSigM in runtime/signal_unix.go).</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	var pcs []uintptr
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	if raceenabled &amp;&amp; pc0 != nil {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		pc1 := (*[1 &lt;&lt; 16]uintptr)(unsafe.Pointer(pc0))
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		pcs = pc1[:ncases:ncases]
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	casePC := func(casi int) uintptr {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		if pcs == nil {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			return 0
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		return pcs[casi]
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	var t0 int64
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if blockprofilerate &gt; 0 {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		t0 = cputicks()
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// The compiler rewrites selects that statically have</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// only 0 or 1 cases plus default into simpler constructs.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// The only way we can end up with such small sel.ncase</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// values here is for a larger select in which most channels</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// have been nilled out. The general code handles those</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// cases correctly, and they are rare enough not to bother</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// optimizing (and needing to test).</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// generate permuted order</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	norder := 0
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	for i := range scases {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		cas := &amp;scases[i]
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		<span class="comment">// Omit cases without channels from the poll and lock orders.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		if cas.c == nil {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			cas.elem = nil <span class="comment">// allow GC</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			continue
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		j := cheaprandn(uint32(norder + 1))
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		pollorder[norder] = pollorder[j]
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		pollorder[j] = uint16(i)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		norder++
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	pollorder = pollorder[:norder]
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	lockorder = lockorder[:norder]
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// sort the cases by Hchan address to get the locking order.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// simple heap sort, to guarantee n log n time and constant stack footprint.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	for i := range lockorder {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		j := i
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		<span class="comment">// Start with the pollorder to permute cases on the same channel.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		c := scases[pollorder[i]].c
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		for j &gt; 0 &amp;&amp; scases[lockorder[(j-1)/2]].c.sortkey() &lt; c.sortkey() {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			k := (j - 1) / 2
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			lockorder[j] = lockorder[k]
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			j = k
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		lockorder[j] = pollorder[i]
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	for i := len(lockorder) - 1; i &gt;= 0; i-- {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		o := lockorder[i]
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		c := scases[o].c
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		lockorder[i] = lockorder[0]
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		j := 0
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		for {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			k := j*2 + 1
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			if k &gt;= i {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>				break
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			if k+1 &lt; i &amp;&amp; scases[lockorder[k]].c.sortkey() &lt; scases[lockorder[k+1]].c.sortkey() {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>				k++
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			if c.sortkey() &lt; scases[lockorder[k]].c.sortkey() {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>				lockorder[j] = lockorder[k]
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>				j = k
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>				continue
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			break
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		lockorder[j] = o
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	if debugSelect {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		for i := 0; i+1 &lt; len(lockorder); i++ {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			if scases[lockorder[i]].c.sortkey() &gt; scases[lockorder[i+1]].c.sortkey() {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>				print(&#34;i=&#34;, i, &#34; x=&#34;, lockorder[i], &#34; y=&#34;, lockorder[i+1], &#34;\n&#34;)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				throw(&#34;select: broken sort&#34;)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// lock all the channels involved in the select</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	sellock(scases, lockorder)
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	var (
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		gp     *g
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		sg     *sudog
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		c      *hchan
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		k      *scase
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		sglist *sudog
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		sgnext *sudog
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		qp     unsafe.Pointer
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		nextp  **sudog
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// pass 1 - look for something already waiting</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	var casi int
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	var cas *scase
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	var caseSuccess bool
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	var caseReleaseTime int64 = -1
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	var recvOK bool
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	for _, casei := range pollorder {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		casi = int(casei)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		cas = &amp;scases[casi]
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		c = cas.c
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		if casi &gt;= nsends {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			sg = c.sendq.dequeue()
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			if sg != nil {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>				goto recv
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			if c.qcount &gt; 0 {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>				goto bufrecv
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			if c.closed != 0 {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>				goto rclose
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		} else {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			if raceenabled {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>				racereadpc(c.raceaddr(), casePC(casi), chansendpc)
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			}
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			if c.closed != 0 {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>				goto sclose
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			sg = c.recvq.dequeue()
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			if sg != nil {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>				goto send
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			if c.qcount &lt; c.dataqsiz {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>				goto bufsend
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	if !block {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		selunlock(scases, lockorder)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		casi = -1
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		goto retc
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	<span class="comment">// pass 2 - enqueue on all chans</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	gp = getg()
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	if gp.waiting != nil {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		throw(&#34;gp.waiting != nil&#34;)
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	nextp = &amp;gp.waiting
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	for _, casei := range lockorder {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		casi = int(casei)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		cas = &amp;scases[casi]
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		c = cas.c
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		sg := acquireSudog()
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		sg.g = gp
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		sg.isSelect = true
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		<span class="comment">// No stack splits between assigning elem and enqueuing</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		<span class="comment">// sg on gp.waiting where copystack can find it.</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		sg.elem = cas.elem
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		sg.releasetime = 0
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		if t0 != 0 {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			sg.releasetime = -1
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		sg.c = c
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		<span class="comment">// Construct waiting list in lock order.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		*nextp = sg
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		nextp = &amp;sg.waitlink
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		if casi &lt; nsends {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>			c.sendq.enqueue(sg)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		} else {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>			c.recvq.enqueue(sg)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	<span class="comment">// wait for someone to wake us up</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	gp.param = nil
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	<span class="comment">// Signal to anyone trying to shrink our stack that we&#39;re about</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	<span class="comment">// to park on a channel. The window between when this G&#39;s status</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	<span class="comment">// changes and when we set gp.activeStackChans is not safe for</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// stack shrinking.</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	gp.parkingOnChan.Store(true)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	gopark(selparkcommit, nil, waitReasonSelect, traceBlockSelect, 1)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	gp.activeStackChans = false
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	sellock(scases, lockorder)
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	gp.selectDone.Store(0)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	sg = (*sudog)(gp.param)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	gp.param = nil
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// pass 3 - dequeue from unsuccessful chans</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// otherwise they stack up on quiet channels</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// record the successful case, if any.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// We singly-linked up the SudoGs in lock order.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	casi = -1
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	cas = nil
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	caseSuccess = false
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	sglist = gp.waiting
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	<span class="comment">// Clear all elem before unlinking from gp.waiting.</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	for sg1 := gp.waiting; sg1 != nil; sg1 = sg1.waitlink {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		sg1.isSelect = false
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		sg1.elem = nil
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		sg1.c = nil
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	gp.waiting = nil
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	for _, casei := range lockorder {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		k = &amp;scases[casei]
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		if sg == sglist {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			<span class="comment">// sg has already been dequeued by the G that woke us up.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			casi = int(casei)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			cas = k
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			caseSuccess = sglist.success
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			if sglist.releasetime &gt; 0 {
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>				caseReleaseTime = sglist.releasetime
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>			}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		} else {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			c = k.c
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			if int(casei) &lt; nsends {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>				c.sendq.dequeueSudoG(sglist)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>			} else {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>				c.recvq.dequeueSudoG(sglist)
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		}
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		sgnext = sglist.waitlink
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		sglist.waitlink = nil
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		releaseSudog(sglist)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		sglist = sgnext
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	if cas == nil {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		throw(&#34;selectgo: bad wakeup&#34;)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	c = cas.c
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	if debugSelect {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		print(&#34;wait-return: cas0=&#34;, cas0, &#34; c=&#34;, c, &#34; cas=&#34;, cas, &#34; send=&#34;, casi &lt; nsends, &#34;\n&#34;)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	if casi &lt; nsends {
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		if !caseSuccess {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			goto sclose
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	} else {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		recvOK = caseSuccess
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	if raceenabled {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		if casi &lt; nsends {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			raceReadObjectPC(c.elemtype, cas.elem, casePC(casi), chansendpc)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		} else if cas.elem != nil {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			raceWriteObjectPC(c.elemtype, cas.elem, casePC(casi), chanrecvpc)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	if msanenabled {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		if casi &lt; nsends {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			msanread(cas.elem, c.elemtype.Size_)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		} else if cas.elem != nil {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			msanwrite(cas.elem, c.elemtype.Size_)
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	if asanenabled {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		if casi &lt; nsends {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			asanread(cas.elem, c.elemtype.Size_)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		} else if cas.elem != nil {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			asanwrite(cas.elem, c.elemtype.Size_)
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	selunlock(scases, lockorder)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	goto retc
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>bufrecv:
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	<span class="comment">// can receive from buffer</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	if raceenabled {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		if cas.elem != nil {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			raceWriteObjectPC(c.elemtype, cas.elem, casePC(casi), chanrecvpc)
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		racenotify(c, c.recvx, nil)
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	if msanenabled &amp;&amp; cas.elem != nil {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		msanwrite(cas.elem, c.elemtype.Size_)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	if asanenabled &amp;&amp; cas.elem != nil {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		asanwrite(cas.elem, c.elemtype.Size_)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	recvOK = true
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	qp = chanbuf(c, c.recvx)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	if cas.elem != nil {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		typedmemmove(c.elemtype, cas.elem, qp)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	typedmemclr(c.elemtype, qp)
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	c.recvx++
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	if c.recvx == c.dataqsiz {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		c.recvx = 0
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	c.qcount--
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	selunlock(scases, lockorder)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	goto retc
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>bufsend:
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	<span class="comment">// can send to buffer</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	if raceenabled {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		racenotify(c, c.sendx, nil)
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		raceReadObjectPC(c.elemtype, cas.elem, casePC(casi), chansendpc)
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	if msanenabled {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		msanread(cas.elem, c.elemtype.Size_)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	if asanenabled {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		asanread(cas.elem, c.elemtype.Size_)
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	typedmemmove(c.elemtype, chanbuf(c, c.sendx), cas.elem)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	c.sendx++
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	if c.sendx == c.dataqsiz {
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		c.sendx = 0
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	c.qcount++
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	selunlock(scases, lockorder)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	goto retc
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>recv:
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	<span class="comment">// can receive from sleeping sender (sg)</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	recv(c, sg, cas.elem, func() { selunlock(scases, lockorder) }, 2)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	if debugSelect {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		print(&#34;syncrecv: cas0=&#34;, cas0, &#34; c=&#34;, c, &#34;\n&#34;)
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	recvOK = true
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	goto retc
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>rclose:
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	<span class="comment">// read at end of closed channel</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	selunlock(scases, lockorder)
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	recvOK = false
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	if cas.elem != nil {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		typedmemclr(c.elemtype, cas.elem)
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	}
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	if raceenabled {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		raceacquire(c.raceaddr())
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	}
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	goto retc
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>send:
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	<span class="comment">// can send to a sleeping receiver (sg)</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	if raceenabled {
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		raceReadObjectPC(c.elemtype, cas.elem, casePC(casi), chansendpc)
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	if msanenabled {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		msanread(cas.elem, c.elemtype.Size_)
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	if asanenabled {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		asanread(cas.elem, c.elemtype.Size_)
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	send(c, sg, cas.elem, func() { selunlock(scases, lockorder) }, 2)
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	if debugSelect {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		print(&#34;syncsend: cas0=&#34;, cas0, &#34; c=&#34;, c, &#34;\n&#34;)
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	goto retc
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>retc:
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	if caseReleaseTime &gt; 0 {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		blockevent(caseReleaseTime-t0, 1)
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	return casi, recvOK
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>sclose:
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	<span class="comment">// send on closed channel</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	selunlock(scases, lockorder)
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	panic(plainError(&#34;send on closed channel&#34;))
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>func (c *hchan) sortkey() uintptr {
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	return uintptr(unsafe.Pointer(c))
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span><span class="comment">// A runtimeSelect is a single case passed to rselect.</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span><span class="comment">// This must match ../reflect/value.go:/runtimeSelect</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>type runtimeSelect struct {
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	dir selectDir
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	typ unsafe.Pointer <span class="comment">// channel type (not used here)</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	ch  *hchan         <span class="comment">// channel</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	val unsafe.Pointer <span class="comment">// ptr to data (SendDir) or ptr to receive buffer (RecvDir)</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>}
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span><span class="comment">// These values must match ../reflect/value.go:/SelectDir.</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>type selectDir int
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>const (
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	_             selectDir = iota
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	selectSend              <span class="comment">// case Chan &lt;- Send</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	selectRecv              <span class="comment">// case &lt;-Chan:</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	selectDefault           <span class="comment">// default</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>)
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_rselect reflect.rselect</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>func reflect_rselect(cases []runtimeSelect) (int, bool) {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	if len(cases) == 0 {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		block()
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	sel := make([]scase, len(cases))
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	orig := make([]int, len(cases))
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	nsends, nrecvs := 0, 0
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	dflt := -1
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	for i, rc := range cases {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		var j int
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		switch rc.dir {
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		case selectDefault:
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>			dflt = i
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			continue
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		case selectSend:
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>			j = nsends
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			nsends++
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		case selectRecv:
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			nrecvs++
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>			j = len(cases) - nrecvs
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		sel[j] = scase{c: rc.ch, elem: rc.val}
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		orig[j] = i
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	<span class="comment">// Only a default case.</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	if nsends+nrecvs == 0 {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		return dflt, false
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	<span class="comment">// Compact sel and orig if necessary.</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	if nsends+nrecvs &lt; len(cases) {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		copy(sel[nsends:], sel[len(cases)-nrecvs:])
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		copy(orig[nsends:], orig[len(cases)-nrecvs:])
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	}
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	order := make([]uint16, 2*(nsends+nrecvs))
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	var pc0 *uintptr
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	if raceenabled {
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		pcs := make([]uintptr, nsends+nrecvs)
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		for i := range pcs {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>			selectsetpc(&amp;pcs[i])
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		pc0 = &amp;pcs[0]
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	chosen, recvOK := selectgo(&amp;sel[0], &amp;order[0], pc0, nsends, nrecvs, dflt == -1)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	<span class="comment">// Translate chosen back to caller&#39;s ordering.</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	if chosen &lt; 0 {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		chosen = dflt
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	} else {
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		chosen = orig[chosen]
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	}
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	return chosen, recvOK
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>func (q *waitq) dequeueSudoG(sgp *sudog) {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	x := sgp.prev
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	y := sgp.next
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	if x != nil {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		if y != nil {
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			<span class="comment">// middle of queue</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>			x.next = y
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>			y.prev = x
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>			sgp.next = nil
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>			sgp.prev = nil
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>			return
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		}
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		<span class="comment">// end of queue</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		x.next = nil
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		q.last = x
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		sgp.prev = nil
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		return
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	}
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	if y != nil {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		<span class="comment">// start of queue</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>		y.prev = nil
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		q.first = y
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		sgp.next = nil
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		return
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	}
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	<span class="comment">// x==y==nil. Either sgp is the only element in the queue,</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	<span class="comment">// or it has already been removed. Use q.first to disambiguate.</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	if q.first == sgp {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		q.first = nil
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		q.last = nil
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	}
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>
</pre><p><a href="select.go?m=text">View as plain text</a></p>

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
