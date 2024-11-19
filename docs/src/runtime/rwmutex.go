<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/rwmutex.go - Go Documentation Server</title>

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
<a href="rwmutex.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">rwmutex.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2017 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>)
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// This is a copy of sync/rwmutex.go rewritten to work in the runtime.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// A rwmutex is a reader/writer mutual exclusion lock.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// The lock can be held by an arbitrary number of readers or a single writer.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// This is a variant of sync.RWMutex, for the runtime package.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// Like mutex, rwmutex blocks the calling M.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// It does not interact with the goroutine scheduler.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type rwmutex struct {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	rLock      mutex    <span class="comment">// protects readers, readerPass, writer</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	readers    muintptr <span class="comment">// list of pending readers</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	readerPass uint32   <span class="comment">// number of pending readers to skip readers list</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	wLock  mutex    <span class="comment">// serializes writers</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	writer muintptr <span class="comment">// pending writer waiting for completing readers</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	readerCount atomic.Int32 <span class="comment">// number of pending readers</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	readerWait  atomic.Int32 <span class="comment">// number of departing readers</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	readRank  lockRank <span class="comment">// semantic lock rank for read locking</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// Lock ranking an rwmutex has two aspects:</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// Semantic ranking: this rwmutex represents some higher level lock that</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// protects some resource (e.g., allocmLock protects creation of new Ms). The</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// read and write locks of that resource need to be represented in the lock</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// rank.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// Internal ranking: as an implementation detail, rwmutex uses two mutexes:</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// rLock and wLock. These have lock order requirements: wLock must be locked</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// before rLock. This also needs to be represented in the lock rank.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// Semantic ranking is represented by acquiring readRank during read lock and</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// writeRank during write lock.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// wLock is held for the duration of a write lock, so it uses writeRank</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// directly, both for semantic and internal ranking. rLock is only held</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// temporarily inside the rlock/lock methods, so it uses readRankInternal to</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// represent internal ranking. Semantic ranking is represented by a separate</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// acquire of readRank for the duration of a read lock.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// The lock ranking must document this ordering:</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// - readRankInternal is a leaf lock.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// - readRank is taken before readRankInternal.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// - writeRank is taken before readRankInternal.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// - readRank is placed in the lock order wherever a read lock of this rwmutex</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//   belongs.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// - writeRank is placed in the lock order wherever a write lock of this</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//   rwmutex belongs.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func (rw *rwmutex) init(readRank, readRankInternal, writeRank lockRank) {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	rw.readRank = readRank
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	lockInit(&amp;rw.rLock, readRankInternal)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	lockInit(&amp;rw.wLock, writeRank)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>const rwmutexMaxReaders = 1 &lt;&lt; 30
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// rlock locks rw for reading.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>func (rw *rwmutex) rlock() {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// The reader must not be allowed to lose its P or else other</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// things blocking on the lock may consume all of the Ps and</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// deadlock (issue #20903). Alternatively, we could drop the P</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// while sleeping.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	acquirem()
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	acquireLockRank(rw.readRank)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	lockWithRankMayAcquire(&amp;rw.rLock, getLockRank(&amp;rw.rLock))
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if rw.readerCount.Add(1) &lt; 0 {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		<span class="comment">// A writer is pending. Park on the reader queue.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			lock(&amp;rw.rLock)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			if rw.readerPass &gt; 0 {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>				<span class="comment">// Writer finished.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>				rw.readerPass -= 1
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>				unlock(&amp;rw.rLock)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			} else {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>				<span class="comment">// Queue this reader to be woken by</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>				<span class="comment">// the writer.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>				m := getg().m
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>				m.schedlink = rw.readers
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>				rw.readers.set(m)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>				unlock(&amp;rw.rLock)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>				notesleep(&amp;m.park)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>				noteclear(&amp;m.park)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		})
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// runlock undoes a single rlock call on rw.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>func (rw *rwmutex) runlock() {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	if r := rw.readerCount.Add(-1); r &lt; 0 {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		if r+1 == 0 || r+1 == -rwmutexMaxReaders {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			throw(&#34;runlock of unlocked rwmutex&#34;)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		<span class="comment">// A writer is pending.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		if rw.readerWait.Add(-1) == 0 {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			<span class="comment">// The last reader unblocks the writer.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			lock(&amp;rw.rLock)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			w := rw.writer.ptr()
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			if w != nil {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>				notewakeup(&amp;w.park)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			unlock(&amp;rw.rLock)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	releaseLockRank(rw.readRank)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	releasem(getg().m)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// lock locks rw for writing.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func (rw *rwmutex) lock() {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// Resolve competition with other writers and stick to our P.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	lock(&amp;rw.wLock)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	m := getg().m
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// Announce that there is a pending writer.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	r := rw.readerCount.Add(-rwmutexMaxReaders) + rwmutexMaxReaders
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// Wait for any active readers to complete.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	lock(&amp;rw.rLock)
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	if r != 0 &amp;&amp; rw.readerWait.Add(r) != 0 {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		<span class="comment">// Wait for reader to wake us up.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		systemstack(func() {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			rw.writer.set(m)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			unlock(&amp;rw.rLock)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			notesleep(&amp;m.park)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			noteclear(&amp;m.park)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		})
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	} else {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		unlock(&amp;rw.rLock)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// unlock unlocks rw for writing.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>func (rw *rwmutex) unlock() {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">// Announce to readers that there is no active writer.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	r := rw.readerCount.Add(rwmutexMaxReaders)
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	if r &gt;= rwmutexMaxReaders {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		throw(&#34;unlock of unlocked rwmutex&#34;)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// Unblock blocked readers.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	lock(&amp;rw.rLock)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	for rw.readers.ptr() != nil {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		reader := rw.readers.ptr()
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		rw.readers = reader.schedlink
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		reader.schedlink.set(nil)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		notewakeup(&amp;reader.park)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		r -= 1
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// If r &gt; 0, there are pending readers that aren&#39;t on the</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// queue. Tell them to skip waiting.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	rw.readerPass += uint32(r)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	unlock(&amp;rw.rLock)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// Allow other writers to proceed.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	unlock(&amp;rw.wLock)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
</pre><p><a href="rwmutex.go?m=text">View as plain text</a></p>

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
