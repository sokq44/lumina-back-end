<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/lock_futex.go - Go Documentation Server</title>

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
<a href="lock_futex.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">lock_futex.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build dragonfly || freebsd || linux</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// This implementation depends on OS-specific implementations of</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//	futexsleep(addr *uint32, val uint32, ns int64)</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//		Atomically,</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//			if *addr == val { sleep }</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//		Might be woken up spuriously; that&#39;s allowed.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//		Don&#39;t sleep longer than ns; ns &lt; 0 means forever.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//	futexwakeup(addr *uint32, cnt uint32)</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//		If any procs are sleeping on addr, wake up at most cnt.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>const (
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	mutex_unlocked = 0
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	mutex_locked   = 1
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	mutex_sleeping = 2
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	active_spin     = 4
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	active_spin_cnt = 30
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	passive_spin    = 1
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>)
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// Possible lock states are mutex_unlocked, mutex_locked and mutex_sleeping.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// mutex_sleeping means that there is presumably at least one sleeping thread.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// Note that there can be spinning threads during all states - they do not</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// affect mutex&#39;s state.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// We use the uintptr mutex.key and note.key as a uint32.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func key32(p *uintptr) *uint32 {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	return (*uint32)(unsafe.Pointer(p))
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>func mutexContended(l *mutex) bool {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	return atomic.Load(key32(&amp;l.key)) &gt; mutex_locked
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>func lock(l *mutex) {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	lockWithRank(l, getLockRank(l))
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>func lock2(l *mutex) {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	gp := getg()
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	if gp.m.locks &lt; 0 {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		throw(&#34;runtime·lock: lock count&#34;)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	gp.m.locks++
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// Speculative grab for lock.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	v := atomic.Xchg(key32(&amp;l.key), mutex_locked)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if v == mutex_unlocked {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		return
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// wait is either MUTEX_LOCKED or MUTEX_SLEEPING</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// depending on whether there is a thread sleeping</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// on this mutex. If we ever change l-&gt;key from</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// MUTEX_SLEEPING to some other value, we must be</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// careful to change it back to MUTEX_SLEEPING before</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// returning, to ensure that the sleeping thread gets</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// its wakeup call.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	wait := v
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	timer := &amp;lockTimer{lock: l}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	timer.begin()
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// On uniprocessors, no point spinning.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// On multiprocessors, spin for ACTIVE_SPIN attempts.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	spin := 0
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	if ncpu &gt; 1 {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		spin = active_spin
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	for {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		<span class="comment">// Try for lock, spinning.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		for i := 0; i &lt; spin; i++ {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			for l.key == mutex_unlocked {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>				if atomic.Cas(key32(&amp;l.key), mutex_unlocked, wait) {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>					timer.end()
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>					return
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>				}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			procyield(active_spin_cnt)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		<span class="comment">// Try for lock, rescheduling.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		for i := 0; i &lt; passive_spin; i++ {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			for l.key == mutex_unlocked {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>				if atomic.Cas(key32(&amp;l.key), mutex_unlocked, wait) {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>					timer.end()
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>					return
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>				}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			osyield()
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		<span class="comment">// Sleep.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		v = atomic.Xchg(key32(&amp;l.key), mutex_sleeping)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		if v == mutex_unlocked {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			timer.end()
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			return
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		wait = mutex_sleeping
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		futexsleep(key32(&amp;l.key), mutex_sleeping, -1)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>func unlock(l *mutex) {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	unlockWithRank(l)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func unlock2(l *mutex) {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	v := atomic.Xchg(key32(&amp;l.key), mutex_unlocked)
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	if v == mutex_unlocked {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		throw(&#34;unlock of unlocked lock&#34;)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if v == mutex_sleeping {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		futexwakeup(key32(&amp;l.key), 1)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	gp := getg()
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	gp.m.mLockProfile.recordUnlock(l)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	gp.m.locks--
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	if gp.m.locks &lt; 0 {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		throw(&#34;runtime·unlock: lock count&#34;)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	if gp.m.locks == 0 &amp;&amp; gp.preempt { <span class="comment">// restore the preemption request in case we&#39;ve cleared it in newstack</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		gp.stackguard0 = stackPreempt
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// One-time notifications.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>func noteclear(n *note) {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	n.key = 0
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>func notewakeup(n *note) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	old := atomic.Xchg(key32(&amp;n.key), 1)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	if old != 0 {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		print(&#34;notewakeup - double wakeup (&#34;, old, &#34;)\n&#34;)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		throw(&#34;notewakeup - double wakeup&#34;)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	futexwakeup(key32(&amp;n.key), 1)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>func notesleep(n *note) {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	gp := getg()
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	if gp != gp.m.g0 {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		throw(&#34;notesleep not on g0&#34;)
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	ns := int64(-1)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	if *cgo_yield != nil {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		<span class="comment">// Sleep for an arbitrary-but-moderate interval to poll libc interceptors.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		ns = 10e6
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	for atomic.Load(key32(&amp;n.key)) == 0 {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		gp.m.blocked = true
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		futexsleep(key32(&amp;n.key), 0, ns)
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		if *cgo_yield != nil {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			asmcgocall(*cgo_yield, nil)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		gp.m.blocked = false
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">// May run with m.p==nil if called from notetsleep, so write barriers</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">// are not allowed.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>func notetsleep_internal(n *note, ns int64) bool {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	gp := getg()
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	if ns &lt; 0 {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		if *cgo_yield != nil {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			<span class="comment">// Sleep for an arbitrary-but-moderate interval to poll libc interceptors.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			ns = 10e6
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		for atomic.Load(key32(&amp;n.key)) == 0 {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			gp.m.blocked = true
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			futexsleep(key32(&amp;n.key), 0, ns)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			if *cgo_yield != nil {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>				asmcgocall(*cgo_yield, nil)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			gp.m.blocked = false
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		return true
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	if atomic.Load(key32(&amp;n.key)) != 0 {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		return true
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	deadline := nanotime() + ns
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	for {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		if *cgo_yield != nil &amp;&amp; ns &gt; 10e6 {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			ns = 10e6
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		gp.m.blocked = true
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		futexsleep(key32(&amp;n.key), 0, ns)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		if *cgo_yield != nil {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			asmcgocall(*cgo_yield, nil)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		gp.m.blocked = false
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		if atomic.Load(key32(&amp;n.key)) != 0 {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			break
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		now := nanotime()
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		if now &gt;= deadline {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			break
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		ns = deadline - now
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	return atomic.Load(key32(&amp;n.key)) != 0
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>func notetsleep(n *note, ns int64) bool {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	gp := getg()
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	if gp != gp.m.g0 &amp;&amp; gp.m.preemptoff != &#34;&#34; {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		throw(&#34;notetsleep not on g0&#34;)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	return notetsleep_internal(n, ns)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">// same as runtime·notetsleep, but called on user g (not g0)</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">// calls only nosplit functions between entersyscallblock/exitsyscall.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>func notetsleepg(n *note, ns int64) bool {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	gp := getg()
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if gp == gp.m.g0 {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		throw(&#34;notetsleepg on g0&#34;)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	entersyscallblock()
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	ok := notetsleep_internal(n, ns)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	exitsyscall()
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	return ok
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>func beforeIdle(int64, int64) (*g, bool) {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	return nil, false
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>func checkTimeouts() {}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
</pre><p><a href="lock_futex.go?m=text">View as plain text</a></p>

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
