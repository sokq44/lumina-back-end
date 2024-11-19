<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/netpoll.go - Go Documentation Server</title>

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
<a href="netpoll.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">netpoll.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2013 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build unix || (js &amp;&amp; wasm) || wasip1 || windows</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// Integrated network poller (platform-independent part).</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// A particular implementation (epoll/kqueue/port/AIX/Windows)</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// must define the following functions:</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// func netpollinit()</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//     Initialize the poller. Only called once.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// func netpollopen(fd uintptr, pd *pollDesc) int32</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//     Arm edge-triggered notifications for fd. The pd argument is to pass</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//     back to netpollready when fd is ready. Return an errno value.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// func netpollclose(fd uintptr) int32</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//     Disable notifications for fd. Return an errno value.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// func netpoll(delta int64) (gList, int32)</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//     Poll the network. If delta &lt; 0, block indefinitely. If delta == 0,</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//     poll without blocking. If delta &gt; 0, block for up to delta nanoseconds.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//     Return a list of goroutines built by calling netpollready,</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//     and a delta to add to netpollWaiters when all goroutines are ready.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//     This will never return an empty list with a non-zero delta.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// func netpollBreak()</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//     Wake up the network poller, assumed to be blocked in netpoll.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// func netpollIsPollDescriptor(fd uintptr) bool</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//     Reports whether fd is a file descriptor used by the poller.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// Error codes returned by runtime_pollReset and runtime_pollWait.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// These must match the values in internal/poll/fd_poll_runtime.go.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>const (
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	pollNoError        = 0 <span class="comment">// no error</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	pollErrClosing     = 1 <span class="comment">// descriptor is closed</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	pollErrTimeout     = 2 <span class="comment">// I/O timeout</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	pollErrNotPollable = 3 <span class="comment">// general error polling descriptor</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// pollDesc contains 2 binary semaphores, rg and wg, to park reader and writer</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// goroutines respectively. The semaphore can be in the following states:</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//	pdReady - io readiness notification is pending;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//	          a goroutine consumes the notification by changing the state to pdNil.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//	pdWait - a goroutine prepares to park on the semaphore, but not yet parked;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//	         the goroutine commits to park by changing the state to G pointer,</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//	         or, alternatively, concurrent io notification changes the state to pdReady,</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//	         or, alternatively, concurrent timeout/close changes the state to pdNil.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//	G pointer - the goroutine is blocked on the semaphore;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//	            io notification or timeout/close changes the state to pdReady or pdNil respectively</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//	            and unparks the goroutine.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//	pdNil - none of the above.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>const (
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	pdNil   uintptr = 0
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	pdReady uintptr = 1
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	pdWait  uintptr = 2
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>const pollBlockSize = 4 * 1024
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// Network poller descriptor.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// No heap pointers.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>type pollDesc struct {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	_     sys.NotInHeap
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	link  *pollDesc      <span class="comment">// in pollcache, protected by pollcache.lock</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	fd    uintptr        <span class="comment">// constant for pollDesc usage lifetime</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	fdseq atomic.Uintptr <span class="comment">// protects against stale pollDesc</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// atomicInfo holds bits from closing, rd, and wd,</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// which are only ever written while holding the lock,</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// summarized for use by netpollcheckerr,</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// which cannot acquire the lock.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// After writing these fields under lock in a way that</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// might change the summary, code must call publishInfo</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// before releasing the lock.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// Code that changes fields and then calls netpollunblock</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// (while still holding the lock) must call publishInfo</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// before calling netpollunblock, because publishInfo is what</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// stops netpollblock from blocking anew</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// (by changing the result of netpollcheckerr).</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// atomicInfo also holds the eventErr bit,</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// recording whether a poll event on the fd got an error;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// atomicInfo is the only source of truth for that bit.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	atomicInfo atomic.Uint32 <span class="comment">// atomic pollInfo</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// rg, wg are accessed atomically and hold g pointers.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// (Using atomic.Uintptr here is similar to using guintptr elsewhere.)</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	rg atomic.Uintptr <span class="comment">// pdReady, pdWait, G waiting for read or pdNil</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	wg atomic.Uintptr <span class="comment">// pdReady, pdWait, G waiting for write or pdNil</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	lock    mutex <span class="comment">// protects the following fields</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	closing bool
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	user    uint32    <span class="comment">// user settable cookie</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	rseq    uintptr   <span class="comment">// protects from stale read timers</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	rt      timer     <span class="comment">// read deadline timer (set if rt.f != nil)</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	rd      int64     <span class="comment">// read deadline (a nanotime in the future, -1 when expired)</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	wseq    uintptr   <span class="comment">// protects from stale write timers</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	wt      timer     <span class="comment">// write deadline timer</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	wd      int64     <span class="comment">// write deadline (a nanotime in the future, -1 when expired)</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	self    *pollDesc <span class="comment">// storage for indirect interface. See (*pollDesc).makeArg.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">// pollInfo is the bits needed by netpollcheckerr, stored atomically,</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// mostly duplicating state that is manipulated under lock in pollDesc.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// The one exception is the pollEventErr bit, which is maintained only</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// in the pollInfo.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>type pollInfo uint32
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>const (
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	pollClosing = 1 &lt;&lt; iota
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	pollEventErr
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	pollExpiredReadDeadline
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	pollExpiredWriteDeadline
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	pollFDSeq <span class="comment">// 20 bit field, low 20 bits of fdseq field</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>const (
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	pollFDSeqBits = 20                   <span class="comment">// number of bits in pollFDSeq</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	pollFDSeqMask = 1&lt;&lt;pollFDSeqBits - 1 <span class="comment">// mask for pollFDSeq</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>func (i pollInfo) closing() bool              { return i&amp;pollClosing != 0 }
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>func (i pollInfo) eventErr() bool             { return i&amp;pollEventErr != 0 }
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>func (i pollInfo) expiredReadDeadline() bool  { return i&amp;pollExpiredReadDeadline != 0 }
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>func (i pollInfo) expiredWriteDeadline() bool { return i&amp;pollExpiredWriteDeadline != 0 }
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// info returns the pollInfo corresponding to pd.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>func (pd *pollDesc) info() pollInfo {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	return pollInfo(pd.atomicInfo.Load())
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// publishInfo updates pd.atomicInfo (returned by pd.info)</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// using the other values in pd.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// It must be called while holding pd.lock,</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// and it must be called after changing anything</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// that might affect the info bits.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// In practice this means after changing closing</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// or changing rd or wd from &lt; 0 to &gt;= 0.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>func (pd *pollDesc) publishInfo() {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	var info uint32
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if pd.closing {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		info |= pollClosing
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if pd.rd &lt; 0 {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		info |= pollExpiredReadDeadline
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if pd.wd &lt; 0 {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		info |= pollExpiredWriteDeadline
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	info |= uint32(pd.fdseq.Load()&amp;pollFDSeqMask) &lt;&lt; pollFDSeq
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// Set all of x except the pollEventErr bit.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	x := pd.atomicInfo.Load()
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	for !pd.atomicInfo.CompareAndSwap(x, (x&amp;pollEventErr)|info) {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		x = pd.atomicInfo.Load()
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// setEventErr sets the result of pd.info().eventErr() to b.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// We only change the error bit if seq == 0 or if seq matches pollFDSeq</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// (issue #59545).</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>func (pd *pollDesc) setEventErr(b bool, seq uintptr) {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	mSeq := uint32(seq &amp; pollFDSeqMask)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	x := pd.atomicInfo.Load()
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	xSeq := (x &gt;&gt; pollFDSeq) &amp; pollFDSeqMask
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	if seq != 0 &amp;&amp; xSeq != mSeq {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		return
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	for (x&amp;pollEventErr != 0) != b &amp;&amp; !pd.atomicInfo.CompareAndSwap(x, x^pollEventErr) {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		x = pd.atomicInfo.Load()
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		xSeq := (x &gt;&gt; pollFDSeq) &amp; pollFDSeqMask
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		if seq != 0 &amp;&amp; xSeq != mSeq {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			return
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>type pollCache struct {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	lock  mutex
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	first *pollDesc
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// PollDesc objects must be type-stable,</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// because we can get ready notification from epoll/kqueue</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// after the descriptor is closed/reused.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// Stale notifications are detected using seq variable,</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// seq is incremented when deadlines are changed or descriptor is reused.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>var (
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	netpollInitLock mutex
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	netpollInited   atomic.Uint32
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	pollcache      pollCache
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	netpollWaiters atomic.Uint32
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_pollServerInit internal/poll.runtime_pollServerInit</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>func poll_runtime_pollServerInit() {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	netpollGenericInit()
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>func netpollGenericInit() {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	if netpollInited.Load() == 0 {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		lockInit(&amp;netpollInitLock, lockRankNetpollInit)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		lock(&amp;netpollInitLock)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		if netpollInited.Load() == 0 {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			netpollinit()
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			netpollInited.Store(1)
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		unlock(&amp;netpollInitLock)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>func netpollinited() bool {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	return netpollInited.Load() != 0
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_isPollServerDescriptor internal/poll.runtime_isPollServerDescriptor</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// poll_runtime_isPollServerDescriptor reports whether fd is a</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// descriptor being used by netpoll.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>func poll_runtime_isPollServerDescriptor(fd uintptr) bool {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	return netpollIsPollDescriptor(fd)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_pollOpen internal/poll.runtime_pollOpen</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>func poll_runtime_pollOpen(fd uintptr) (*pollDesc, int) {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	pd := pollcache.alloc()
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	lock(&amp;pd.lock)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	wg := pd.wg.Load()
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if wg != pdNil &amp;&amp; wg != pdReady {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		throw(&#34;runtime: blocked write on free polldesc&#34;)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	rg := pd.rg.Load()
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	if rg != pdNil &amp;&amp; rg != pdReady {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		throw(&#34;runtime: blocked read on free polldesc&#34;)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	pd.fd = fd
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	if pd.fdseq.Load() == 0 {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		<span class="comment">// The value 0 is special in setEventErr, so don&#39;t use it.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		pd.fdseq.Store(1)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	pd.closing = false
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	pd.setEventErr(false, 0)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	pd.rseq++
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	pd.rg.Store(pdNil)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	pd.rd = 0
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	pd.wseq++
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	pd.wg.Store(pdNil)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	pd.wd = 0
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	pd.self = pd
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	pd.publishInfo()
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	unlock(&amp;pd.lock)
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	errno := netpollopen(fd, pd)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	if errno != 0 {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		pollcache.free(pd)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		return nil, int(errno)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	return pd, 0
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_pollClose internal/poll.runtime_pollClose</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>func poll_runtime_pollClose(pd *pollDesc) {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	if !pd.closing {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		throw(&#34;runtime: close polldesc w/o unblock&#34;)
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	wg := pd.wg.Load()
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	if wg != pdNil &amp;&amp; wg != pdReady {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		throw(&#34;runtime: blocked write on closing polldesc&#34;)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	rg := pd.rg.Load()
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	if rg != pdNil &amp;&amp; rg != pdReady {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		throw(&#34;runtime: blocked read on closing polldesc&#34;)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	netpollclose(pd.fd)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	pollcache.free(pd)
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>func (c *pollCache) free(pd *pollDesc) {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// pd can&#39;t be shared here, but lock anyhow because</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// that&#39;s what publishInfo documents.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	lock(&amp;pd.lock)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">// Increment the fdseq field, so that any currently</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// running netpoll calls will not mark pd as ready.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	fdseq := pd.fdseq.Load()
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	fdseq = (fdseq + 1) &amp; (1&lt;&lt;taggedPointerBits - 1)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	pd.fdseq.Store(fdseq)
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	pd.publishInfo()
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	unlock(&amp;pd.lock)
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	lock(&amp;c.lock)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	pd.link = c.first
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	c.first = pd
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	unlock(&amp;c.lock)
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>}
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">// poll_runtime_pollReset, which is internal/poll.runtime_pollReset,</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span><span class="comment">// prepares a descriptor for polling in mode, which is &#39;r&#39; or &#39;w&#39;.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">// This returns an error code; the codes are defined above.</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_pollReset internal/poll.runtime_pollReset</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>func poll_runtime_pollReset(pd *pollDesc, mode int) int {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	errcode := netpollcheckerr(pd, int32(mode))
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	if errcode != pollNoError {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		return errcode
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	if mode == &#39;r&#39; {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		pd.rg.Store(pdNil)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	} else if mode == &#39;w&#39; {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		pd.wg.Store(pdNil)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	return pollNoError
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">// poll_runtime_pollWait, which is internal/poll.runtime_pollWait,</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span><span class="comment">// waits for a descriptor to be ready for reading or writing,</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// according to mode, which is &#39;r&#39; or &#39;w&#39;.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">// This returns an error code; the codes are defined above.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_pollWait internal/poll.runtime_pollWait</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func poll_runtime_pollWait(pd *pollDesc, mode int) int {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	errcode := netpollcheckerr(pd, int32(mode))
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	if errcode != pollNoError {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		return errcode
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	<span class="comment">// As for now only Solaris, illumos, AIX and wasip1 use level-triggered IO.</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	if GOOS == &#34;solaris&#34; || GOOS == &#34;illumos&#34; || GOOS == &#34;aix&#34; || GOOS == &#34;wasip1&#34; {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		netpollarm(pd, mode)
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	for !netpollblock(pd, int32(mode), false) {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		errcode = netpollcheckerr(pd, int32(mode))
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		if errcode != pollNoError {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>			return errcode
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		<span class="comment">// Can happen if timeout has fired and unblocked us,</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		<span class="comment">// but before we had a chance to run, timeout has been reset.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		<span class="comment">// Pretend it has not happened and retry.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	return pollNoError
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>}
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_pollWaitCanceled internal/poll.runtime_pollWaitCanceled</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>func poll_runtime_pollWaitCanceled(pd *pollDesc, mode int) {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	<span class="comment">// This function is used only on windows after a failed attempt to cancel</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	<span class="comment">// a pending async IO operation. Wait for ioready, ignore closing or timeouts.</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	for !netpollblock(pd, int32(mode), true) {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_pollSetDeadline internal/poll.runtime_pollSetDeadline</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>func poll_runtime_pollSetDeadline(pd *pollDesc, d int64, mode int) {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	lock(&amp;pd.lock)
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	if pd.closing {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		unlock(&amp;pd.lock)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		return
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	rd0, wd0 := pd.rd, pd.wd
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	combo0 := rd0 &gt; 0 &amp;&amp; rd0 == wd0
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	if d &gt; 0 {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		d += nanotime()
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		if d &lt;= 0 {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			<span class="comment">// If the user has a deadline in the future, but the delay calculation</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			<span class="comment">// overflows, then set the deadline to the maximum possible value.</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			d = 1&lt;&lt;63 - 1
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	if mode == &#39;r&#39; || mode == &#39;r&#39;+&#39;w&#39; {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		pd.rd = d
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	if mode == &#39;w&#39; || mode == &#39;r&#39;+&#39;w&#39; {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		pd.wd = d
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	pd.publishInfo()
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	combo := pd.rd &gt; 0 &amp;&amp; pd.rd == pd.wd
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	rtf := netpollReadDeadline
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	if combo {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		rtf = netpollDeadline
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	if pd.rt.f == nil {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		if pd.rd &gt; 0 {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			pd.rt.f = rtf
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			<span class="comment">// Copy current seq into the timer arg.</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			<span class="comment">// Timer func will check the seq against current descriptor seq,</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			<span class="comment">// if they differ the descriptor was reused or timers were reset.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			pd.rt.arg = pd.makeArg()
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			pd.rt.seq = pd.rseq
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			resettimer(&amp;pd.rt, pd.rd)
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		}
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	} else if pd.rd != rd0 || combo != combo0 {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		pd.rseq++ <span class="comment">// invalidate current timers</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		if pd.rd &gt; 0 {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			modtimer(&amp;pd.rt, pd.rd, 0, rtf, pd.makeArg(), pd.rseq)
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		} else {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			deltimer(&amp;pd.rt)
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			pd.rt.f = nil
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		}
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	if pd.wt.f == nil {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		if pd.wd &gt; 0 &amp;&amp; !combo {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			pd.wt.f = netpollWriteDeadline
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			pd.wt.arg = pd.makeArg()
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			pd.wt.seq = pd.wseq
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			resettimer(&amp;pd.wt, pd.wd)
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		}
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	} else if pd.wd != wd0 || combo != combo0 {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		pd.wseq++ <span class="comment">// invalidate current timers</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		if pd.wd &gt; 0 &amp;&amp; !combo {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			modtimer(&amp;pd.wt, pd.wd, 0, netpollWriteDeadline, pd.makeArg(), pd.wseq)
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		} else {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>			deltimer(&amp;pd.wt)
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			pd.wt.f = nil
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	<span class="comment">// If we set the new deadline in the past, unblock currently pending IO if any.</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	<span class="comment">// Note that pd.publishInfo has already been called, above, immediately after modifying rd and wd.</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	delta := int32(0)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	var rg, wg *g
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	if pd.rd &lt; 0 {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		rg = netpollunblock(pd, &#39;r&#39;, false, &amp;delta)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	if pd.wd &lt; 0 {
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		wg = netpollunblock(pd, &#39;w&#39;, false, &amp;delta)
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	unlock(&amp;pd.lock)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	if rg != nil {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		netpollgoready(rg, 3)
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	if wg != nil {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		netpollgoready(wg, 3)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	netpollAdjustWaiters(delta)
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span><span class="comment">//go:linkname poll_runtime_pollUnblock internal/poll.runtime_pollUnblock</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>func poll_runtime_pollUnblock(pd *pollDesc) {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	lock(&amp;pd.lock)
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	if pd.closing {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		throw(&#34;runtime: unblock on closing polldesc&#34;)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	pd.closing = true
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	pd.rseq++
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	pd.wseq++
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	var rg, wg *g
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	pd.publishInfo()
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	delta := int32(0)
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	rg = netpollunblock(pd, &#39;r&#39;, false, &amp;delta)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	wg = netpollunblock(pd, &#39;w&#39;, false, &amp;delta)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	if pd.rt.f != nil {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		deltimer(&amp;pd.rt)
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		pd.rt.f = nil
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	if pd.wt.f != nil {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		deltimer(&amp;pd.wt)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		pd.wt.f = nil
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	unlock(&amp;pd.lock)
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	if rg != nil {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		netpollgoready(rg, 3)
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	if wg != nil {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		netpollgoready(wg, 3)
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	netpollAdjustWaiters(delta)
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>}
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span><span class="comment">// netpollready is called by the platform-specific netpoll function.</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span><span class="comment">// It declares that the fd associated with pd is ready for I/O.</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span><span class="comment">// The toRun argument is used to build a list of goroutines to return</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span><span class="comment">// from netpoll. The mode argument is &#39;r&#39;, &#39;w&#39;, or &#39;r&#39;+&#39;w&#39; to indicate</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span><span class="comment">// whether the fd is ready for reading or writing or both.</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span><span class="comment">// This returns a delta to apply to netpollWaiters.</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span><span class="comment">// This may run while the world is stopped, so write barriers are not allowed.</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>func netpollready(toRun *gList, pd *pollDesc, mode int32) int32 {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	delta := int32(0)
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	var rg, wg *g
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	if mode == &#39;r&#39; || mode == &#39;r&#39;+&#39;w&#39; {
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		rg = netpollunblock(pd, &#39;r&#39;, true, &amp;delta)
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	}
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	if mode == &#39;w&#39; || mode == &#39;r&#39;+&#39;w&#39; {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		wg = netpollunblock(pd, &#39;w&#39;, true, &amp;delta)
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	}
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	if rg != nil {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		toRun.push(rg)
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	if wg != nil {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		toRun.push(wg)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	return delta
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>func netpollcheckerr(pd *pollDesc, mode int32) int {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	info := pd.info()
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	if info.closing() {
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		return pollErrClosing
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	if (mode == &#39;r&#39; &amp;&amp; info.expiredReadDeadline()) || (mode == &#39;w&#39; &amp;&amp; info.expiredWriteDeadline()) {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		return pollErrTimeout
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	<span class="comment">// Report an event scanning error only on a read event.</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	<span class="comment">// An error on a write event will be captured in a subsequent</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	<span class="comment">// write call that is able to report a more specific error.</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	if mode == &#39;r&#39; &amp;&amp; info.eventErr() {
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		return pollErrNotPollable
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	}
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	return pollNoError
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>}
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>func netpollblockcommit(gp *g, gpp unsafe.Pointer) bool {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	r := atomic.Casuintptr((*uintptr)(gpp), pdWait, uintptr(unsafe.Pointer(gp)))
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	if r {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		<span class="comment">// Bump the count of goroutines waiting for the poller.</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		<span class="comment">// The scheduler uses this to decide whether to block</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		<span class="comment">// waiting for the poller if there is nothing else to do.</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		netpollAdjustWaiters(1)
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	return r
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>func netpollgoready(gp *g, traceskip int) {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	goready(gp, traceskip+1)
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>}
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">// returns true if IO is ready, or false if timed out or closed</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">// waitio - wait only for completed IO, ignore errors</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span><span class="comment">// Concurrent calls to netpollblock in the same mode are forbidden, as pollDesc</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span><span class="comment">// can hold only a single waiting goroutine for each mode.</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>func netpollblock(pd *pollDesc, mode int32, waitio bool) bool {
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	gpp := &amp;pd.rg
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	if mode == &#39;w&#39; {
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		gpp = &amp;pd.wg
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	<span class="comment">// set the gpp semaphore to pdWait</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	for {
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		<span class="comment">// Consume notification if already ready.</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		if gpp.CompareAndSwap(pdReady, pdNil) {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>			return true
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		if gpp.CompareAndSwap(pdNil, pdWait) {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>			break
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		<span class="comment">// Double check that this isn&#39;t corrupt; otherwise we&#39;d loop</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		<span class="comment">// forever.</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		if v := gpp.Load(); v != pdReady &amp;&amp; v != pdNil {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>			throw(&#34;runtime: double wait&#34;)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	<span class="comment">// need to recheck error states after setting gpp to pdWait</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	<span class="comment">// this is necessary because runtime_pollUnblock/runtime_pollSetDeadline/deadlineimpl</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	<span class="comment">// do the opposite: store to closing/rd/wd, publishInfo, load of rg/wg</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	if waitio || netpollcheckerr(pd, mode) == pollNoError {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		gopark(netpollblockcommit, unsafe.Pointer(gpp), waitReasonIOWait, traceBlockNet, 5)
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	<span class="comment">// be careful to not lose concurrent pdReady notification</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	old := gpp.Swap(pdNil)
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	if old &gt; pdWait {
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		throw(&#34;runtime: corrupted polldesc&#34;)
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	return old == pdReady
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>}
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span><span class="comment">// netpollunblock moves either pd.rg (if mode == &#39;r&#39;) or</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span><span class="comment">// pd.wg (if mode == &#39;w&#39;) into the pdReady state.</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span><span class="comment">// This returns any goroutine blocked on pd.{rg,wg}.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span><span class="comment">// It adds any adjustment to netpollWaiters to *delta;</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span><span class="comment">// this adjustment should be applied after the goroutine has</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span><span class="comment">// been marked ready.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>func netpollunblock(pd *pollDesc, mode int32, ioready bool, delta *int32) *g {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	gpp := &amp;pd.rg
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	if mode == &#39;w&#39; {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		gpp = &amp;pd.wg
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	for {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		old := gpp.Load()
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		if old == pdReady {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>			return nil
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		}
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		if old == pdNil &amp;&amp; !ioready {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>			<span class="comment">// Only set pdReady for ioready. runtime_pollWait</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>			<span class="comment">// will check for timeout/cancel before waiting.</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>			return nil
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		}
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		new := pdNil
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		if ioready {
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>			new = pdReady
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		if gpp.CompareAndSwap(old, new) {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>			if old == pdWait {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>				old = pdNil
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>			} else if old != pdNil {
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>				*delta -= 1
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>			}
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>			return (*g)(unsafe.Pointer(old))
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		}
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	}
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>}
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>func netpolldeadlineimpl(pd *pollDesc, seq uintptr, read, write bool) {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	lock(&amp;pd.lock)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	<span class="comment">// Seq arg is seq when the timer was set.</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	<span class="comment">// If it&#39;s stale, ignore the timer event.</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	currentSeq := pd.rseq
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	if !read {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		currentSeq = pd.wseq
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	if seq != currentSeq {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		<span class="comment">// The descriptor was reused or timers were reset.</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		unlock(&amp;pd.lock)
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		return
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	delta := int32(0)
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	var rg *g
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	if read {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		if pd.rd &lt;= 0 || pd.rt.f == nil {
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>			throw(&#34;runtime: inconsistent read deadline&#34;)
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		}
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		pd.rd = -1
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		pd.publishInfo()
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		rg = netpollunblock(pd, &#39;r&#39;, false, &amp;delta)
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	var wg *g
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	if write {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		if pd.wd &lt;= 0 || pd.wt.f == nil &amp;&amp; !read {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>			throw(&#34;runtime: inconsistent write deadline&#34;)
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		pd.wd = -1
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		pd.publishInfo()
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		wg = netpollunblock(pd, &#39;w&#39;, false, &amp;delta)
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	unlock(&amp;pd.lock)
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	if rg != nil {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		netpollgoready(rg, 0)
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	}
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	if wg != nil {
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		netpollgoready(wg, 0)
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	netpollAdjustWaiters(delta)
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>}
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>func netpollDeadline(arg any, seq uintptr) {
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	netpolldeadlineimpl(arg.(*pollDesc), seq, true, true)
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>}
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>func netpollReadDeadline(arg any, seq uintptr) {
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	netpolldeadlineimpl(arg.(*pollDesc), seq, true, false)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>}
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>func netpollWriteDeadline(arg any, seq uintptr) {
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	netpolldeadlineimpl(arg.(*pollDesc), seq, false, true)
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>}
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span><span class="comment">// netpollAnyWaiters reports whether any goroutines are waiting for I/O.</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>func netpollAnyWaiters() bool {
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	return netpollWaiters.Load() &gt; 0
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>}
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span><span class="comment">// netpollAdjustWaiters adds delta to netpollWaiters.</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>func netpollAdjustWaiters(delta int32) {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	if delta != 0 {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		netpollWaiters.Add(delta)
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	}
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>func (c *pollCache) alloc() *pollDesc {
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	lock(&amp;c.lock)
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	if c.first == nil {
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		const pdSize = unsafe.Sizeof(pollDesc{})
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		n := pollBlockSize / pdSize
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		if n == 0 {
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>			n = 1
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		}
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		<span class="comment">// Must be in non-GC memory because can be referenced</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>		<span class="comment">// only from epoll/kqueue internals.</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		mem := persistentalloc(n*pdSize, 0, &amp;memstats.other_sys)
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; n; i++ {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			pd := (*pollDesc)(add(mem, i*pdSize))
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			pd.link = c.first
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>			c.first = pd
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		}
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	}
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	pd := c.first
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	c.first = pd.link
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	lockInit(&amp;pd.lock, lockRankPollDesc)
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	unlock(&amp;c.lock)
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	return pd
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>}
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span><span class="comment">// makeArg converts pd to an interface{}.</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span><span class="comment">// makeArg does not do any allocation. Normally, such</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span><span class="comment">// a conversion requires an allocation because pointers to</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span><span class="comment">// types which embed runtime/internal/sys.NotInHeap (which pollDesc is)</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span><span class="comment">// must be stored in interfaces indirectly. See issue 42076.</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>func (pd *pollDesc) makeArg() (i any) {
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	x := (*eface)(unsafe.Pointer(&amp;i))
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	x._type = pdType
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	x.data = unsafe.Pointer(&amp;pd.self)
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	return
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>}
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>var (
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	pdEface any    = (*pollDesc)(nil)
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	pdType  *_type = efaceOf(&amp;pdEface)._type
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>)
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>
</pre><p><a href="netpoll.go?m=text">View as plain text</a></p>

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
