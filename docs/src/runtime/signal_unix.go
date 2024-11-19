<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/signal_unix.go - Go Documentation Server</title>

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
<a href="signal_unix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">signal_unix.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2012 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build unix</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;runtime/internal/sys&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// sigTabT is the type of an entry in the global sigtable array.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// sigtable is inherently system dependent, and appears in OS-specific files,</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// but sigTabT is the same for all Unixy systems.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// The sigtable array is indexed by a system signal number to get the flags</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// and printable name of each signal.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>type sigTabT struct {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	flags int32
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	name  string
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//go:linkname os_sigpipe os.sigpipe</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>func os_sigpipe() {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	systemstack(sigpipe)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>func signame(sig uint32) string {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	if sig &gt;= uint32(len(sigtable)) {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	return sigtable[sig].name
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>const (
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	_SIG_DFL uintptr = 0
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	_SIG_IGN uintptr = 1
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// sigPreempt is the signal used for non-cooperative preemption.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// There&#39;s no good way to choose this signal, but there are some</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// heuristics:</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// 1. It should be a signal that&#39;s passed-through by debuggers by</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// default. On Linux, this is SIGALRM, SIGURG, SIGCHLD, SIGIO,</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// SIGVTALRM, SIGPROF, and SIGWINCH, plus some glibc-internal signals.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// 2. It shouldn&#39;t be used internally by libc in mixed Go/C binaries</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// because libc may assume it&#39;s the only thing that can handle these</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// signals. For example SIGCANCEL or SIGSETXID.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// 3. It should be a signal that can happen spuriously without</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// consequences. For example, SIGALRM is a bad choice because the</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// signal handler can&#39;t tell if it was caused by the real process</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// alarm or not (arguably this means the signal is broken, but I</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// digress). SIGUSR1 and SIGUSR2 are also bad because those are often</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// used in meaningful ways by applications.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// 4. We need to deal with platforms without real-time signals (like</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// macOS), so those are out.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// We use SIGURG because it meets all of these criteria, is extremely</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// unlikely to be used by an application for its &#34;real&#34; meaning (both</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// because out-of-band data is basically unused and because SIGURG</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// doesn&#39;t report which socket has the condition, making it pretty</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// useless), and even if it is, the application has to be ready for</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// spurious SIGURG. SIGIO wouldn&#39;t be a bad choice either, but is more</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// likely to be used for real.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>const sigPreempt = _SIGURG
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// Stores the signal handlers registered before Go installed its own.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// These signal handlers will be invoked in cases where Go doesn&#39;t want to</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// handle a particular signal (e.g., signal occurred on a non-Go thread).</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// See sigfwdgo for more information on when the signals are forwarded.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// This is read by the signal handler; accesses should use</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// atomic.Loaduintptr and atomic.Storeuintptr.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>var fwdSig [_NSIG]uintptr
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// handlingSig is indexed by signal number and is non-zero if we are</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// currently handling the signal. Or, to put it another way, whether</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// the signal handler is currently set to the Go signal handler or not.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// This is uint32 rather than bool so that we can use atomic instructions.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>var handlingSig [_NSIG]uint32
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// channels for synchronizing signal mask updates with the signal mask</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// thread</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>var (
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	disableSigChan  chan uint32
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	enableSigChan   chan uint32
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	maskUpdatedChan chan struct{}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>func init() {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// _NSIG is the number of signals on this operating system.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// sigtable should describe what to do for all the possible signals.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	if len(sigtable) != _NSIG {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		print(&#34;runtime: len(sigtable)=&#34;, len(sigtable), &#34; _NSIG=&#34;, _NSIG, &#34;\n&#34;)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		throw(&#34;bad sigtable len&#34;)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>var signalsOK bool
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// Initialize signals.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// Called by libpreinit so runtime may not be initialized.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>func initsig(preinit bool) {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if !preinit {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		<span class="comment">// It&#39;s now OK for signal handlers to run.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		signalsOK = true
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// For c-archive/c-shared this is called by libpreinit with</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// preinit == true.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	if (isarchive || islibrary) &amp;&amp; !preinit {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		return
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	for i := uint32(0); i &lt; _NSIG; i++ {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		t := &amp;sigtable[i]
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		if t.flags == 0 || t.flags&amp;_SigDefault != 0 {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			continue
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		<span class="comment">// We don&#39;t need to use atomic operations here because</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		<span class="comment">// there shouldn&#39;t be any other goroutines running yet.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		fwdSig[i] = getsig(i)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		if !sigInstallGoHandler(i) {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			<span class="comment">// Even if we are not installing a signal handler,</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			<span class="comment">// set SA_ONSTACK if necessary.</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			if fwdSig[i] != _SIG_DFL &amp;&amp; fwdSig[i] != _SIG_IGN {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>				setsigstack(i)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			} else if fwdSig[i] == _SIG_IGN {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>				sigInitIgnored(i)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			continue
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		handlingSig[i] = 1
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		setsig(i, abi.FuncPCABIInternal(sighandler))
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>func sigInstallGoHandler(sig uint32) bool {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// For some signals, we respect an inherited SIG_IGN handler</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// rather than insist on installing our own default handler.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// Even these signals can be fetched using the os/signal package.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	switch sig {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	case _SIGHUP, _SIGINT:
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		if atomic.Loaduintptr(&amp;fwdSig[sig]) == _SIG_IGN {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			return false
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	if (GOOS == &#34;linux&#34; || GOOS == &#34;android&#34;) &amp;&amp; !iscgo &amp;&amp; sig == sigPerThreadSyscall {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		<span class="comment">// sigPerThreadSyscall is the same signal used by glibc for</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		<span class="comment">// per-thread syscalls on Linux. We use it for the same purpose</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// in non-cgo binaries.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		return true
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	t := &amp;sigtable[sig]
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	if t.flags&amp;_SigSetStack != 0 {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		return false
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// When built using c-archive or c-shared, only install signal</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// handlers for synchronous signals and SIGPIPE and sigPreempt.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	if (isarchive || islibrary) &amp;&amp; t.flags&amp;_SigPanic == 0 &amp;&amp; sig != _SIGPIPE &amp;&amp; sig != sigPreempt {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		return false
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	return true
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// sigenable enables the Go signal handler to catch the signal sig.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// It is only called while holding the os/signal.handlers lock,</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// via os/signal.enableSignal and signal_enable.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>func sigenable(sig uint32) {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	if sig &gt;= uint32(len(sigtable)) {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		return
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// SIGPROF is handled specially for profiling.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	if sig == _SIGPROF {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		return
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	t := &amp;sigtable[sig]
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	if t.flags&amp;_SigNotify != 0 {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		ensureSigM()
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		enableSigChan &lt;- sig
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		&lt;-maskUpdatedChan
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		if atomic.Cas(&amp;handlingSig[sig], 0, 1) {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			atomic.Storeuintptr(&amp;fwdSig[sig], getsig(sig))
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			setsig(sig, abi.FuncPCABIInternal(sighandler))
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// sigdisable disables the Go signal handler for the signal sig.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// It is only called while holding the os/signal.handlers lock,</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// via os/signal.disableSignal and signal_disable.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>func sigdisable(sig uint32) {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	if sig &gt;= uint32(len(sigtable)) {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		return
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	<span class="comment">// SIGPROF is handled specially for profiling.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	if sig == _SIGPROF {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		return
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	t := &amp;sigtable[sig]
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	if t.flags&amp;_SigNotify != 0 {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		ensureSigM()
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		disableSigChan &lt;- sig
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		&lt;-maskUpdatedChan
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		<span class="comment">// If initsig does not install a signal handler for a</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		<span class="comment">// signal, then to go back to the state before Notify</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		<span class="comment">// we should remove the one we installed.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		if !sigInstallGoHandler(sig) {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			atomic.Store(&amp;handlingSig[sig], 0)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			setsig(sig, atomic.Loaduintptr(&amp;fwdSig[sig]))
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span><span class="comment">// sigignore ignores the signal sig.</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span><span class="comment">// It is only called while holding the os/signal.handlers lock,</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span><span class="comment">// via os/signal.ignoreSignal and signal_ignore.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>func sigignore(sig uint32) {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	if sig &gt;= uint32(len(sigtable)) {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		return
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// SIGPROF is handled specially for profiling.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	if sig == _SIGPROF {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		return
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	t := &amp;sigtable[sig]
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	if t.flags&amp;_SigNotify != 0 {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		atomic.Store(&amp;handlingSig[sig], 0)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		setsig(sig, _SIG_IGN)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// clearSignalHandlers clears all signal handlers that are not ignored</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span><span class="comment">// back to the default. This is called by the child after a fork, so that</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span><span class="comment">// we can enable the signal mask for the exec without worrying about</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">// running a signal handler in the child.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>func clearSignalHandlers() {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	for i := uint32(0); i &lt; _NSIG; i++ {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		if atomic.Load(&amp;handlingSig[i]) != 0 {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			setsig(i, _SIG_DFL)
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span><span class="comment">// setProcessCPUProfilerTimer is called when the profiling timer changes.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span><span class="comment">// It is called with prof.signalLock held. hz is the new timer, and is 0 if</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// profiling is being disabled. Enable or disable the signal as</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span><span class="comment">// required for -buildmode=c-archive.</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>func setProcessCPUProfilerTimer(hz int32) {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	if hz != 0 {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		<span class="comment">// Enable the Go signal handler if not enabled.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		if atomic.Cas(&amp;handlingSig[_SIGPROF], 0, 1) {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			h := getsig(_SIGPROF)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			<span class="comment">// If no signal handler was installed before, then we record</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			<span class="comment">// _SIG_IGN here. When we turn off profiling (below) we&#39;ll start</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			<span class="comment">// ignoring SIGPROF signals. We do this, rather than change</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			<span class="comment">// to SIG_DFL, because there may be a pending SIGPROF</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			<span class="comment">// signal that has not yet been delivered to some other thread.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			<span class="comment">// If we change to SIG_DFL when turning off profiling, the</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			<span class="comment">// program will crash when that SIGPROF is delivered. We assume</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			<span class="comment">// that programs that use profiling don&#39;t want to crash on a</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			<span class="comment">// stray SIGPROF. See issue 19320.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			<span class="comment">// We do the change here instead of when turning off profiling,</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			<span class="comment">// because there we may race with a signal handler running</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			<span class="comment">// concurrently, in particular, sigfwdgo may observe _SIG_DFL and</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			<span class="comment">// die. See issue 43828.</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			if h == _SIG_DFL {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>				h = _SIG_IGN
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			atomic.Storeuintptr(&amp;fwdSig[_SIGPROF], h)
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			setsig(_SIGPROF, abi.FuncPCABIInternal(sighandler))
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		var it itimerval
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		it.it_interval.tv_sec = 0
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		it.it_interval.set_usec(1000000 / hz)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		it.it_value = it.it_interval
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		setitimer(_ITIMER_PROF, &amp;it, nil)
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	} else {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		setitimer(_ITIMER_PROF, &amp;itimerval{}, nil)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		<span class="comment">// If the Go signal handler should be disabled by default,</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		<span class="comment">// switch back to the signal handler that was installed</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		<span class="comment">// when we enabled profiling. We don&#39;t try to handle the case</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		<span class="comment">// of a program that changes the SIGPROF handler while Go</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		<span class="comment">// profiling is enabled.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		if !sigInstallGoHandler(_SIGPROF) {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>			if atomic.Cas(&amp;handlingSig[_SIGPROF], 1, 0) {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>				h := atomic.Loaduintptr(&amp;fwdSig[_SIGPROF])
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>				setsig(_SIGPROF, h)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">// setThreadCPUProfilerHz makes any thread-specific changes required to</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span><span class="comment">// implement profiling at a rate of hz.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// No changes required on Unix systems when using setitimer.</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>func setThreadCPUProfilerHz(hz int32) {
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	getg().m.profilehz = hz
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>func sigpipe() {
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	if signal_ignored(_SIGPIPE) || sigsend(_SIGPIPE) {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		return
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	dieFromSignal(_SIGPIPE)
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// doSigPreempt handles a preemption signal on gp.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>func doSigPreempt(gp *g, ctxt *sigctxt) {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	<span class="comment">// Check if this G wants to be preempted and is safe to</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	<span class="comment">// preempt.</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	if wantAsyncPreempt(gp) {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		if ok, newpc := isAsyncSafePoint(gp, ctxt.sigpc(), ctxt.sigsp(), ctxt.siglr()); ok {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>			<span class="comment">// Adjust the PC and inject a call to asyncPreempt.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			ctxt.pushCall(abi.FuncPCABI0(asyncPreempt), newpc)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	<span class="comment">// Acknowledge the preemption.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	gp.m.preemptGen.Add(1)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	gp.m.signalPending.Store(0)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	if GOOS == &#34;darwin&#34; || GOOS == &#34;ios&#34; {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		pendingPreemptSignals.Add(-1)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>const preemptMSupported = true
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// preemptM sends a preemption request to mp. This request may be</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span><span class="comment">// handled asynchronously and may be coalesced with other requests to</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">// the M. When the request is received, if the running G or P are</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// marked for preemption and the goroutine is at an asynchronous</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// safe-point, it will preempt the goroutine. It always atomically</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">// increments mp.preemptGen after handling a preemption request.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>func preemptM(mp *m) {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// On Darwin, don&#39;t try to preempt threads during exec.</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	<span class="comment">// Issue #41702.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	if GOOS == &#34;darwin&#34; || GOOS == &#34;ios&#34; {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		execLock.rlock()
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	if mp.signalPending.CompareAndSwap(0, 1) {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		if GOOS == &#34;darwin&#34; || GOOS == &#34;ios&#34; {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			pendingPreemptSignals.Add(1)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		<span class="comment">// If multiple threads are preempting the same M, it may send many</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		<span class="comment">// signals to the same M such that it hardly make progress, causing</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		<span class="comment">// live-lock problem. Apparently this could happen on darwin. See</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		<span class="comment">// issue #37741.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		<span class="comment">// Only send a signal if there isn&#39;t already one pending.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		signalM(mp, sigPreempt)
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	if GOOS == &#34;darwin&#34; || GOOS == &#34;ios&#34; {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		execLock.runlock()
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	}
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>}
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">// sigFetchG fetches the value of G safely when running in a signal handler.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">// On some architectures, the g value may be clobbered when running in a VDSO.</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// See issue #32912.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>func sigFetchG(c *sigctxt) *g {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	switch GOARCH {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	case &#34;arm&#34;, &#34;arm64&#34;, &#34;loong64&#34;, &#34;ppc64&#34;, &#34;ppc64le&#34;, &#34;riscv64&#34;, &#34;s390x&#34;:
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		if !iscgo &amp;&amp; inVDSOPage(c.sigpc()) {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			<span class="comment">// When using cgo, we save the g on TLS and load it from there</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			<span class="comment">// in sigtramp. Just use that.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			<span class="comment">// Otherwise, before making a VDSO call we save the g to the</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			<span class="comment">// bottom of the signal stack. Fetch from there.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			<span class="comment">// TODO: in efence mode, stack is sysAlloc&#39;d, so this wouldn&#39;t</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			<span class="comment">// work.</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			sp := getcallersp()
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			s := spanOf(sp)
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			if s != nil &amp;&amp; s.state.get() == mSpanManual &amp;&amp; s.base() &lt; sp &amp;&amp; sp &lt; s.limit {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>				gp := *(**g)(unsafe.Pointer(s.base()))
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>				return gp
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			return nil
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	return getg()
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span><span class="comment">// sigtrampgo is called from the signal handler function, sigtramp,</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span><span class="comment">// written in assembly code.</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span><span class="comment">// This is called by the signal handler, and the world may be stopped.</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span><span class="comment">// It must be nosplit because getg() is still the G that was running</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span><span class="comment">// (if any) when the signal was delivered, but it&#39;s (usually) called</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span><span class="comment">// on the gsignal stack. Until this switches the G to gsignal, the</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span><span class="comment">// stack bounds check won&#39;t work.</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>func sigtrampgo(sig uint32, info *siginfo, ctx unsafe.Pointer) {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	if sigfwdgo(sig, info, ctx) {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		return
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	c := &amp;sigctxt{info, ctx}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	gp := sigFetchG(c)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	setg(gp)
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	if gp == nil || (gp.m != nil &amp;&amp; gp.m.isExtraInC) {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		if sig == _SIGPROF {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			<span class="comment">// Some platforms (Linux) have per-thread timers, which we use in</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>			<span class="comment">// combination with the process-wide timer. Avoid double-counting.</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			if validSIGPROF(nil, c) {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>				sigprofNonGoPC(c.sigpc())
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			return
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		if sig == sigPreempt &amp;&amp; preemptMSupported &amp;&amp; debug.asyncpreemptoff == 0 {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>			<span class="comment">// This is probably a signal from preemptM sent</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>			<span class="comment">// while executing Go code but received while</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>			<span class="comment">// executing non-Go code.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			<span class="comment">// We got past sigfwdgo, so we know that there is</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			<span class="comment">// no non-Go signal handler for sigPreempt.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			<span class="comment">// The default behavior for sigPreempt is to ignore</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>			<span class="comment">// the signal, so badsignal will be a no-op anyway.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			if GOOS == &#34;darwin&#34; || GOOS == &#34;ios&#34; {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>				pendingPreemptSignals.Add(-1)
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			return
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		c.fixsigcode(sig)
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		<span class="comment">// Set g to nil here and badsignal will use g0 by needm.</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		<span class="comment">// TODO: reuse the current m here by using the gsignal and adjustSignalStack,</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		<span class="comment">// since the current g maybe a normal goroutine and actually running on the signal stack,</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		<span class="comment">// it may hit stack split that is not expected here.</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		if gp != nil {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			setg(nil)
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		badsignal(uintptr(sig), c)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		<span class="comment">// Restore g</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		if gp != nil {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			setg(gp)
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		return
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	setg(gp.m.gsignal)
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	<span class="comment">// If some non-Go code called sigaltstack, adjust.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	var gsignalStack gsignalStack
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	setStack := adjustSignalStack(sig, gp.m, &amp;gsignalStack)
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	if setStack {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		gp.m.gsignal.stktopsp = getcallersp()
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	}
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	if gp.stackguard0 == stackFork {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		signalDuringFork(sig)
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	c.fixsigcode(sig)
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	sighandler(sig, info, ctx, gp)
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	setg(gp)
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	if setStack {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		restoreGsignalStack(&amp;gsignalStack)
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	}
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span><span class="comment">// If the signal handler receives a SIGPROF signal on a non-Go thread,</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span><span class="comment">// it tries to collect a traceback into sigprofCallers.</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span><span class="comment">// sigprofCallersUse is set to non-zero while sigprofCallers holds a traceback.</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>var sigprofCallers cgoCallers
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>var sigprofCallersUse uint32
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span><span class="comment">// sigprofNonGo is called if we receive a SIGPROF signal on a non-Go thread,</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span><span class="comment">// and the signal handler collected a stack trace in sigprofCallers.</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span><span class="comment">// When this is called, sigprofCallersUse will be non-zero.</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span><span class="comment">// g is nil, and what we can do is very limited.</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span><span class="comment">// It is called from the signal handling functions written in assembly code that</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span><span class="comment">// are active for cgo programs, cgoSigtramp and sigprofNonGoWrapper, which have</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span><span class="comment">// not verified that the SIGPROF delivery corresponds to the best available</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span><span class="comment">// profiling source for this thread.</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>func sigprofNonGo(sig uint32, info *siginfo, ctx unsafe.Pointer) {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	if prof.hz.Load() != 0 {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		c := &amp;sigctxt{info, ctx}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		<span class="comment">// Some platforms (Linux) have per-thread timers, which we use in</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		<span class="comment">// combination with the process-wide timer. Avoid double-counting.</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		if validSIGPROF(nil, c) {
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			n := 0
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>			for n &lt; len(sigprofCallers) &amp;&amp; sigprofCallers[n] != 0 {
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>				n++
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>			}
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>			cpuprof.addNonGo(sigprofCallers[:n])
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	}
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	atomic.Store(&amp;sigprofCallersUse, 0)
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>}
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span><span class="comment">// sigprofNonGoPC is called when a profiling signal arrived on a</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span><span class="comment">// non-Go thread and we have a single PC value, not a stack trace.</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span><span class="comment">// g is nil, and what we can do is very limited.</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>func sigprofNonGoPC(pc uintptr) {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	if prof.hz.Load() != 0 {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		stk := []uintptr{
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>			pc,
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			abi.FuncPCABIInternal(_ExternalCode) + sys.PCQuantum,
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		cpuprof.addNonGo(stk)
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>}
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span><span class="comment">// adjustSignalStack adjusts the current stack guard based on the</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span><span class="comment">// stack pointer that is actually in use while handling a signal.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span><span class="comment">// We do this in case some non-Go code called sigaltstack.</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span><span class="comment">// This reports whether the stack was adjusted, and if so stores the old</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// signal stack in *gsigstack.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>func adjustSignalStack(sig uint32, mp *m, gsigStack *gsignalStack) bool {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	sp := uintptr(unsafe.Pointer(&amp;sig))
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	if sp &gt;= mp.gsignal.stack.lo &amp;&amp; sp &lt; mp.gsignal.stack.hi {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		return false
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	}
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	var st stackt
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	sigaltstack(nil, &amp;st)
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	stsp := uintptr(unsafe.Pointer(st.ss_sp))
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	if st.ss_flags&amp;_SS_DISABLE == 0 &amp;&amp; sp &gt;= stsp &amp;&amp; sp &lt; stsp+st.ss_size {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		setGsignalStack(&amp;st, gsigStack)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		return true
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	if sp &gt;= mp.g0.stack.lo &amp;&amp; sp &lt; mp.g0.stack.hi {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		<span class="comment">// The signal was delivered on the g0 stack.</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		<span class="comment">// This can happen when linked with C code</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		<span class="comment">// using the thread sanitizer, which collects</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		<span class="comment">// signals then delivers them itself by calling</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		<span class="comment">// the signal handler directly when C code,</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		<span class="comment">// including C code called via cgo, calls a</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		<span class="comment">// TSAN-intercepted function such as malloc.</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		<span class="comment">// We check this condition last as g0.stack.lo</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		<span class="comment">// may be not very accurate (see mstart).</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		st := stackt{ss_size: mp.g0.stack.hi - mp.g0.stack.lo}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		setSignalstackSP(&amp;st, mp.g0.stack.lo)
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		setGsignalStack(&amp;st, gsigStack)
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		return true
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	<span class="comment">// sp is not within gsignal stack, g0 stack, or sigaltstack. Bad.</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	setg(nil)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	needm(true)
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	if st.ss_flags&amp;_SS_DISABLE != 0 {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		noSignalStack(sig)
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	} else {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		sigNotOnStack(sig, sp, mp)
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	dropm()
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	return false
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>}
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span><span class="comment">// crashing is the number of m&#39;s we have waited for when implementing</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span><span class="comment">// GOTRACEBACK=crash when a signal is received.</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>var crashing atomic.Int32
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span><span class="comment">// testSigtrap and testSigusr1 are used by the runtime tests. If</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span><span class="comment">// non-nil, it is called on SIGTRAP/SIGUSR1. If it returns true, the</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span><span class="comment">// normal behavior on this signal is suppressed.</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>var testSigtrap func(info *siginfo, ctxt *sigctxt, gp *g) bool
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>var testSigusr1 func(gp *g) bool
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span><span class="comment">// sighandler is invoked when a signal occurs. The global g will be</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span><span class="comment">// set to a gsignal goroutine and we will be running on the alternate</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span><span class="comment">// signal stack. The parameter gp will be the value of the global g</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span><span class="comment">// when the signal occurred. The sig, info, and ctxt parameters are</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span><span class="comment">// from the system signal handler: they are the parameters passed when</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span><span class="comment">// the SA is passed to the sigaction system call.</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span><span class="comment">// The garbage collector may have stopped the world, so write barriers</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span><span class="comment">// are not allowed.</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>func sighandler(sig uint32, info *siginfo, ctxt unsafe.Pointer, gp *g) {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	<span class="comment">// The g executing the signal handler. This is almost always</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	<span class="comment">// mp.gsignal. See delayedSignal for an exception.</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	gsignal := getg()
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	mp := gsignal.m
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	c := &amp;sigctxt{info, ctxt}
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	<span class="comment">// Cgo TSAN (not the Go race detector) intercepts signals and calls the</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	<span class="comment">// signal handler at a later time. When the signal handler is called, the</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	<span class="comment">// memory may have changed, but the signal context remains old. The</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	<span class="comment">// unmatched signal context and memory makes it unsafe to unwind or inspect</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	<span class="comment">// the stack. So we ignore delayed non-fatal signals that will cause a stack</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	<span class="comment">// inspection (profiling signal and preemption signal).</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	<span class="comment">// cgo_yield is only non-nil for TSAN, and is specifically used to trigger</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	<span class="comment">// signal delivery. We use that as an indicator of delayed signals.</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	<span class="comment">// For delayed signals, the handler is called on the g0 stack (see</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	<span class="comment">// adjustSignalStack).</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	delayedSignal := *cgo_yield != nil &amp;&amp; mp != nil &amp;&amp; gsignal.stack == mp.g0.stack
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	if sig == _SIGPROF {
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		<span class="comment">// Some platforms (Linux) have per-thread timers, which we use in</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		<span class="comment">// combination with the process-wide timer. Avoid double-counting.</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		if !delayedSignal &amp;&amp; validSIGPROF(mp, c) {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>			sigprof(c.sigpc(), c.sigsp(), c.siglr(), gp, mp)
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		return
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	}
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	if sig == _SIGTRAP &amp;&amp; testSigtrap != nil &amp;&amp; testSigtrap(info, (*sigctxt)(noescape(unsafe.Pointer(c))), gp) {
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		return
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	}
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	if sig == _SIGUSR1 &amp;&amp; testSigusr1 != nil &amp;&amp; testSigusr1(gp) {
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		return
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	}
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	if (GOOS == &#34;linux&#34; || GOOS == &#34;android&#34;) &amp;&amp; sig == sigPerThreadSyscall {
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		<span class="comment">// sigPerThreadSyscall is the same signal used by glibc for</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		<span class="comment">// per-thread syscalls on Linux. We use it for the same purpose</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		<span class="comment">// in non-cgo binaries. Since this signal is not _SigNotify,</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		<span class="comment">// there is nothing more to do once we run the syscall.</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		runPerThreadSyscall()
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		return
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	}
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	if sig == sigPreempt &amp;&amp; debug.asyncpreemptoff == 0 &amp;&amp; !delayedSignal {
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		<span class="comment">// Might be a preemption signal.</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		doSigPreempt(gp, c)
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		<span class="comment">// Even if this was definitely a preemption signal, it</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		<span class="comment">// may have been coalesced with another signal, so we</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		<span class="comment">// still let it through to the application.</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	flags := int32(_SigThrow)
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	if sig &lt; uint32(len(sigtable)) {
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		flags = sigtable[sig].flags
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	}
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	if !c.sigFromUser() &amp;&amp; flags&amp;_SigPanic != 0 &amp;&amp; (gp.throwsplit || gp != mp.curg) {
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		<span class="comment">// We can&#39;t safely sigpanic because it may grow the</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		<span class="comment">// stack. Abort in the signal handler instead.</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		<span class="comment">// Also don&#39;t inject a sigpanic if we are not on a</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		<span class="comment">// user G stack. Either we&#39;re in the runtime, or we&#39;re</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		<span class="comment">// running C code. Either way we cannot recover.</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		flags = _SigThrow
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	if isAbortPC(c.sigpc()) {
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		<span class="comment">// On many architectures, the abort function just</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		<span class="comment">// causes a memory fault. Don&#39;t turn that into a panic.</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		flags = _SigThrow
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	}
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	if !c.sigFromUser() &amp;&amp; flags&amp;_SigPanic != 0 {
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		<span class="comment">// The signal is going to cause a panic.</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		<span class="comment">// Arrange the stack so that it looks like the point</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		<span class="comment">// where the signal occurred made a call to the</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		<span class="comment">// function sigpanic. Then set the PC to sigpanic.</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		<span class="comment">// Have to pass arguments out of band since</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		<span class="comment">// augmenting the stack frame would break</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		<span class="comment">// the unwinding code.</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>		gp.sig = sig
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>		gp.sigcode0 = uintptr(c.sigcode())
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		gp.sigcode1 = c.fault()
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		gp.sigpc = c.sigpc()
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		c.preparePanic(sig, gp)
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		return
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	}
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	if c.sigFromUser() || flags&amp;_SigNotify != 0 {
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		if sigsend(sig) {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>			return
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>		}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	}
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	if c.sigFromUser() &amp;&amp; signal_ignored(sig) {
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		return
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	}
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	if flags&amp;_SigKill != 0 {
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>		dieFromSignal(sig)
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	}
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	<span class="comment">// _SigThrow means that we should exit now.</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	<span class="comment">// If we get here with _SigPanic, it means that the signal</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	<span class="comment">// was sent to us by a program (c.sigFromUser() is true);</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	<span class="comment">// in that case, if we didn&#39;t handle it in sigsend, we exit now.</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	if flags&amp;(_SigThrow|_SigPanic) == 0 {
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		return
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	mp.throwing = throwTypeRuntime
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	mp.caughtsig.set(gp)
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	if crashing.Load() == 0 {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		startpanic_m()
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	}
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	gp = fatalsignal(sig, c, gp, mp)
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	level, _, docrash := gotraceback()
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	if level &gt; 0 {
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		goroutineheader(gp)
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		tracebacktrap(c.sigpc(), c.sigsp(), c.siglr(), gp)
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>		if crashing.Load() &gt; 0 &amp;&amp; gp != mp.curg &amp;&amp; mp.curg != nil &amp;&amp; readgstatus(mp.curg)&amp;^_Gscan == _Grunning {
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>			<span class="comment">// tracebackothers on original m skipped this one; trace it now.</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>			goroutineheader(mp.curg)
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>			traceback(^uintptr(0), ^uintptr(0), 0, mp.curg)
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		} else if crashing.Load() == 0 {
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>			tracebackothers(gp)
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>			print(&#34;\n&#34;)
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		}
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		dumpregs(c)
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	}
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	if docrash {
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>		isCrashThread := false
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		if crashing.CompareAndSwap(0, 1) {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>			isCrashThread = true
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		} else {
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>			crashing.Add(1)
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		if crashing.Load() &lt; mcount()-int32(extraMLength.Load()) {
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>			<span class="comment">// There are other m&#39;s that need to dump their stacks.</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>			<span class="comment">// Relay SIGQUIT to the next m by sending it to the current process.</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			<span class="comment">// All m&#39;s that have already received SIGQUIT have signal masks blocking</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>			<span class="comment">// receipt of any signals, so the SIGQUIT will go to an m that hasn&#39;t seen it yet.</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>			<span class="comment">// The first m will wait until all ms received the SIGQUIT, then crash/exit.</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>			<span class="comment">// Just in case the relaying gets botched, each m involved in</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>			<span class="comment">// the relay sleeps for 5 seconds and then does the crash/exit itself.</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>			<span class="comment">// The faulting m is crashing first so it is the faulting thread in the core dump (see issue #63277):</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>			<span class="comment">// in expected operation, the first m will wait until the last m has received the SIGQUIT,</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			<span class="comment">// and then run crash/exit and the process is gone.</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>			<span class="comment">// However, if it spends more than 5 seconds to send SIGQUIT to all ms,</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>			<span class="comment">// any of ms may crash/exit the process after waiting for 5 seconds.</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>			print(&#34;\n-----\n\n&#34;)
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>			raiseproc(_SIGQUIT)
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		}
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>		if isCrashThread {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>			i := 0
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>			for (crashing.Load() &lt; mcount()-int32(extraMLength.Load())) &amp;&amp; i &lt; 10 {
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>				i++
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>				usleep(500 * 1000)
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>			}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>		} else {
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>			usleep(5 * 1000 * 1000)
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>		}
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		printDebugLog()
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		crash()
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>	}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	printDebugLog()
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	exit(2)
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>}
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>func fatalsignal(sig uint32, c *sigctxt, gp *g, mp *m) *g {
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	if sig &lt; uint32(len(sigtable)) {
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>		print(sigtable[sig].name, &#34;\n&#34;)
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	} else {
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		print(&#34;Signal &#34;, sig, &#34;\n&#34;)
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	}
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	if isSecureMode() {
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		exit(2)
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	}
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	print(&#34;PC=&#34;, hex(c.sigpc()), &#34; m=&#34;, mp.id, &#34; sigcode=&#34;, c.sigcode())
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	if sig == _SIGSEGV || sig == _SIGBUS {
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		print(&#34; addr=&#34;, hex(c.fault()))
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	}
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	print(&#34;\n&#34;)
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	if mp.incgo &amp;&amp; gp == mp.g0 &amp;&amp; mp.curg != nil {
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>		print(&#34;signal arrived during cgo execution\n&#34;)
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		<span class="comment">// Switch to curg so that we get a traceback of the Go code</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		<span class="comment">// leading up to the cgocall, which switched from curg to g0.</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>		gp = mp.curg
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	}
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	if sig == _SIGILL || sig == _SIGFPE {
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		<span class="comment">// It would be nice to know how long the instruction is.</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		<span class="comment">// Unfortunately, that&#39;s complicated to do in general (mostly for x86</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>		<span class="comment">// and s930x, but other archs have non-standard instruction lengths also).</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>		<span class="comment">// Opt to print 16 bytes, which covers most instructions.</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>		const maxN = 16
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		n := uintptr(maxN)
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		<span class="comment">// We have to be careful, though. If we&#39;re near the end of</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>		<span class="comment">// a page and the following page isn&#39;t mapped, we could</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		<span class="comment">// segfault. So make sure we don&#39;t straddle a page (even though</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		<span class="comment">// that could lead to printing an incomplete instruction).</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>		<span class="comment">// We&#39;re assuming here we can read at least the page containing the PC.</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		<span class="comment">// I suppose it is possible that the page is mapped executable but not readable?</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		pc := c.sigpc()
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		if n &gt; physPageSize-pc%physPageSize {
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>			n = physPageSize - pc%physPageSize
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>		}
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		print(&#34;instruction bytes:&#34;)
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		b := (*[maxN]byte)(unsafe.Pointer(pc))
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; n; i++ {
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>			print(&#34; &#34;, hex(b[i]))
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		}
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		println()
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	}
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	print(&#34;\n&#34;)
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	return gp
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span><span class="comment">// sigpanic turns a synchronous signal into a run-time panic.</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span><span class="comment">// If the signal handler sees a synchronous panic, it arranges the</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span><span class="comment">// stack to look like the function where the signal occurred called</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span><span class="comment">// sigpanic, sets the signal&#39;s PC value to sigpanic, and returns from</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span><span class="comment">// the signal handler. The effect is that the program will act as</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span><span class="comment">// though the function that got the signal simply called sigpanic</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span><span class="comment">// instead.</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span><span class="comment">// This must NOT be nosplit because the linker doesn&#39;t know where</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span><span class="comment">// sigpanic calls can be injected.</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span><span class="comment">// The signal handler must not inject a call to sigpanic if</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span><span class="comment">// getg().throwsplit, since sigpanic may need to grow the stack.</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span><span class="comment">// This is exported via linkname to assembly in runtime/cgo.</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span><span class="comment">//go:linkname sigpanic</span>
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>func sigpanic() {
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	gp := getg()
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	if !canpanic() {
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		throw(&#34;unexpected signal during runtime execution&#34;)
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	}
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	switch gp.sig {
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	case _SIGBUS:
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>		if gp.sigcode0 == _BUS_ADRERR &amp;&amp; gp.sigcode1 &lt; 0x1000 {
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>			panicmem()
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		}
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		<span class="comment">// Support runtime/debug.SetPanicOnFault.</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		if gp.paniconfault {
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>			panicmemAddr(gp.sigcode1)
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>		}
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>		print(&#34;unexpected fault address &#34;, hex(gp.sigcode1), &#34;\n&#34;)
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>		throw(&#34;fault&#34;)
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	case _SIGSEGV:
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>		if (gp.sigcode0 == 0 || gp.sigcode0 == _SEGV_MAPERR || gp.sigcode0 == _SEGV_ACCERR) &amp;&amp; gp.sigcode1 &lt; 0x1000 {
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>			panicmem()
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>		}
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>		<span class="comment">// Support runtime/debug.SetPanicOnFault.</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>		if gp.paniconfault {
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>			panicmemAddr(gp.sigcode1)
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>		}
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>		if inUserArenaChunk(gp.sigcode1) {
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>			<span class="comment">// We could check that the arena chunk is explicitly set to fault,</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			<span class="comment">// but the fact that we faulted on accessing it is enough to prove</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>			<span class="comment">// that it is.</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			print(&#34;accessed data from freed user arena &#34;, hex(gp.sigcode1), &#34;\n&#34;)
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>		} else {
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>			print(&#34;unexpected fault address &#34;, hex(gp.sigcode1), &#34;\n&#34;)
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>		}
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		throw(&#34;fault&#34;)
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	case _SIGFPE:
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>		switch gp.sigcode0 {
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		case _FPE_INTDIV:
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>			panicdivide()
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		case _FPE_INTOVF:
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>			panicoverflow()
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>		}
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>		panicfloat()
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	}
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	if gp.sig &gt;= uint32(len(sigtable)) {
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>		<span class="comment">// can&#39;t happen: we looked up gp.sig in sigtable to decide to call sigpanic</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>		throw(&#34;unexpected signal value&#34;)
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	}
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	panic(errorString(sigtable[gp.sig].name))
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>}
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span><span class="comment">// dieFromSignal kills the program with a signal.</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span><span class="comment">// This provides the expected exit status for the shell.</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span><span class="comment">// This is only called with fatal signals expected to kill the process.</span>
<span id="L916" class="ln">   916&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L917" class="ln">   917&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>func dieFromSignal(sig uint32) {
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>	unblocksig(sig)
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>	<span class="comment">// Mark the signal as unhandled to ensure it is forwarded.</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	atomic.Store(&amp;handlingSig[sig], 0)
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	raise(sig)
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	<span class="comment">// That should have killed us. On some systems, though, raise</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	<span class="comment">// sends the signal to the whole process rather than to just</span>
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	<span class="comment">// the current thread, which means that the signal may not yet</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	<span class="comment">// have been delivered. Give other threads a chance to run and</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>	<span class="comment">// pick up the signal.</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>	osyield()
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>	osyield()
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>	osyield()
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	<span class="comment">// If that didn&#39;t work, try _SIG_DFL.</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>	setsig(sig, _SIG_DFL)
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>	raise(sig)
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	osyield()
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	osyield()
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	osyield()
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>	<span class="comment">// If we are still somehow running, just exit with the wrong status.</span>
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>	exit(2)
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>}
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span><span class="comment">// raisebadsignal is called when a signal is received on a non-Go</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span><span class="comment">// thread, and the Go program does not want to handle it (that is, the</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span><span class="comment">// program has not called os/signal.Notify for the signal).</span>
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>func raisebadsignal(sig uint32, c *sigctxt) {
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>	if sig == _SIGPROF {
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>		<span class="comment">// Ignore profiling signals that arrive on non-Go threads.</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>		return
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>	}
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>	var handler uintptr
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>	if sig &gt;= _NSIG {
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>		handler = _SIG_DFL
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>	} else {
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		handler = atomic.Loaduintptr(&amp;fwdSig[sig])
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>	}
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>	<span class="comment">// Reset the signal handler and raise the signal.</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>	<span class="comment">// We are currently running inside a signal handler, so the</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	<span class="comment">// signal is blocked. We need to unblock it before raising the</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	<span class="comment">// signal, or the signal we raise will be ignored until we return</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>	<span class="comment">// from the signal handler. We know that the signal was unblocked</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>	<span class="comment">// before entering the handler, or else we would not have received</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>	<span class="comment">// it. That means that we don&#39;t have to worry about blocking it</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>	<span class="comment">// again.</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	unblocksig(sig)
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>	setsig(sig, handler)
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;re linked into a non-Go program we want to try to</span>
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>	<span class="comment">// avoid modifying the original context in which the signal</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>	<span class="comment">// was raised. If the handler is the default, we know it</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	<span class="comment">// is non-recoverable, so we don&#39;t have to worry about</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>	<span class="comment">// re-installing sighandler. At this point we can just</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	<span class="comment">// return and the signal will be re-raised and caught by</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>	<span class="comment">// the default handler with the correct context.</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	<span class="comment">// On FreeBSD, the libthr sigaction code prevents</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	<span class="comment">// this from working so we fall through to raise.</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>	if GOOS != &#34;freebsd&#34; &amp;&amp; (isarchive || islibrary) &amp;&amp; handler == _SIG_DFL &amp;&amp; !c.sigFromUser() {
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>		return
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	}
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>	raise(sig)
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>	<span class="comment">// Give the signal a chance to be delivered.</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	<span class="comment">// In almost all real cases the program is about to crash,</span>
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>	<span class="comment">// so sleeping here is not a waste of time.</span>
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>	usleep(1000)
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	<span class="comment">// If the signal didn&#39;t cause the program to exit, restore the</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>	<span class="comment">// Go signal handler and carry on.</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	<span class="comment">// We may receive another instance of the signal before we</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	<span class="comment">// restore the Go handler, but that is not so bad: we know</span>
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	<span class="comment">// that the Go program has been ignoring the signal.</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	setsig(sig, abi.FuncPCABIInternal(sighandler))
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>}
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>func crash() {
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	dieFromSignal(_SIGABRT)
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>}
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span><span class="comment">// ensureSigM starts one global, sleeping thread to make sure at least one thread</span>
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span><span class="comment">// is available to catch signals enabled for os/signal.</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>func ensureSigM() {
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>	if maskUpdatedChan != nil {
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>		return
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	}
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	maskUpdatedChan = make(chan struct{})
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	disableSigChan = make(chan uint32)
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	enableSigChan = make(chan uint32)
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	go func() {
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>		<span class="comment">// Signal masks are per-thread, so make sure this goroutine stays on one</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>		<span class="comment">// thread.</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>		LockOSThread()
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		defer UnlockOSThread()
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>		<span class="comment">// The sigBlocked mask contains the signals not active for os/signal,</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>		<span class="comment">// initially all signals except the essential. When signal.Notify()/Stop is called,</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>		<span class="comment">// sigenable/sigdisable in turn notify this thread to update its signal</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		<span class="comment">// mask accordingly.</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>		sigBlocked := sigset_all
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>		for i := range sigtable {
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>			if !blockableSig(uint32(i)) {
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>				sigdelset(&amp;sigBlocked, i)
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>			}
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>		}
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>		sigprocmask(_SIG_SETMASK, &amp;sigBlocked, nil)
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		for {
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>			select {
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>			case sig := &lt;-enableSigChan:
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>				if sig &gt; 0 {
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>					sigdelset(&amp;sigBlocked, int(sig))
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>				}
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>			case sig := &lt;-disableSigChan:
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>				if sig &gt; 0 &amp;&amp; blockableSig(sig) {
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>					sigaddset(&amp;sigBlocked, int(sig))
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>				}
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>			}
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>			sigprocmask(_SIG_SETMASK, &amp;sigBlocked, nil)
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>			maskUpdatedChan &lt;- struct{}{}
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>		}
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>	}()
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>}
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span><span class="comment">// This is called when we receive a signal when there is no signal stack.</span>
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span><span class="comment">// This can only happen if non-Go code calls sigaltstack to disable the</span>
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span><span class="comment">// signal stack.</span>
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>func noSignalStack(sig uint32) {
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>	println(&#34;signal&#34;, sig, &#34;received on thread with no signal stack&#34;)
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>	throw(&#34;non-Go code disabled sigaltstack&#34;)
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>}
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span><span class="comment">// This is called if we receive a signal when there is a signal stack</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span><span class="comment">// but we are not on it. This can only happen if non-Go code called</span>
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span><span class="comment">// sigaction without setting the SS_ONSTACK flag.</span>
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>func sigNotOnStack(sig uint32, sp uintptr, mp *m) {
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>	println(&#34;signal&#34;, sig, &#34;received but handler not on signal stack&#34;)
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>	print(&#34;mp.gsignal stack [&#34;, hex(mp.gsignal.stack.lo), &#34; &#34;, hex(mp.gsignal.stack.hi), &#34;], &#34;)
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>	print(&#34;mp.g0 stack [&#34;, hex(mp.g0.stack.lo), &#34; &#34;, hex(mp.g0.stack.hi), &#34;], sp=&#34;, hex(sp), &#34;\n&#34;)
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>	throw(&#34;non-Go code set up signal handler without SA_ONSTACK flag&#34;)
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>}
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span><span class="comment">// signalDuringFork is called if we receive a signal while doing a fork.</span>
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span><span class="comment">// We do not want signals at that time, as a signal sent to the process</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span><span class="comment">// group may be delivered to the child process, causing confusion.</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span><span class="comment">// This should never be called, because we block signals across the fork;</span>
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span><span class="comment">// this function is just a safety check. See issue 18600 for background.</span>
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>func signalDuringFork(sig uint32) {
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	println(&#34;signal&#34;, sig, &#34;received during fork&#34;)
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	throw(&#34;signal received during fork&#34;)
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>}
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span><span class="comment">// This runs on a foreign stack, without an m or a g. No stack split.</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span><span class="comment">//go:norace</span>
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>func badsignal(sig uintptr, c *sigctxt) {
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>	if !iscgo &amp;&amp; !cgoHasExtraM {
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>		<span class="comment">// There is no extra M. needm will not be able to grab</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>		<span class="comment">// an M. Instead of hanging, just crash.</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>		<span class="comment">// Cannot call split-stack function as there is no G.</span>
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>		writeErrStr(&#34;fatal: bad g in signal handler\n&#34;)
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>		exit(2)
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>		*(*uintptr)(unsafe.Pointer(uintptr(123))) = 2
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>	}
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>	needm(true)
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>	if !sigsend(uint32(sig)) {
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>		<span class="comment">// A foreign thread received the signal sig, and the</span>
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>		<span class="comment">// Go code does not want to handle it.</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>		raisebadsignal(uint32(sig), c)
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>	}
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>	dropm()
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>}
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>func sigfwd(fn uintptr, sig uint32, info *siginfo, ctx unsafe.Pointer)
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span><span class="comment">// Determines if the signal should be handled by Go and if not, forwards the</span>
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span><span class="comment">// signal to the handler that was installed before Go&#39;s. Returns whether the</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span><span class="comment">// signal was forwarded.</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span><span class="comment">// This is called by the signal handler, and the world may be stopped.</span>
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>func sigfwdgo(sig uint32, info *siginfo, ctx unsafe.Pointer) bool {
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>	if sig &gt;= uint32(len(sigtable)) {
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>		return false
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>	}
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>	fwdFn := atomic.Loaduintptr(&amp;fwdSig[sig])
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>	flags := sigtable[sig].flags
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>	<span class="comment">// If we aren&#39;t handling the signal, forward it.</span>
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>	if atomic.Load(&amp;handlingSig[sig]) == 0 || !signalsOK {
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>		<span class="comment">// If the signal is ignored, doing nothing is the same as forwarding.</span>
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>		if fwdFn == _SIG_IGN || (fwdFn == _SIG_DFL &amp;&amp; flags&amp;_SigIgn != 0) {
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>			return true
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>		}
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>		<span class="comment">// We are not handling the signal and there is no other handler to forward to.</span>
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>		<span class="comment">// Crash with the default behavior.</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>		if fwdFn == _SIG_DFL {
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>			setsig(sig, _SIG_DFL)
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>			dieFromSignal(sig)
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>			return false
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>		}
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>		sigfwd(fwdFn, sig, info, ctx)
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>		return true
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>	}
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>	<span class="comment">// This function and its caller sigtrampgo assumes SIGPIPE is delivered on the</span>
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>	<span class="comment">// originating thread. This property does not hold on macOS (golang.org/issue/33384),</span>
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>	<span class="comment">// so we have no choice but to ignore SIGPIPE.</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>	if (GOOS == &#34;darwin&#34; || GOOS == &#34;ios&#34;) &amp;&amp; sig == _SIGPIPE {
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>		return true
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	}
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	<span class="comment">// If there is no handler to forward to, no need to forward.</span>
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>	if fwdFn == _SIG_DFL {
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>		return false
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>	}
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>	c := &amp;sigctxt{info, ctx}
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>	<span class="comment">// Only forward synchronous signals and SIGPIPE.</span>
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>	<span class="comment">// Unfortunately, user generated SIGPIPEs will also be forwarded, because si_code</span>
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>	<span class="comment">// is set to _SI_USER even for a SIGPIPE raised from a write to a closed socket</span>
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>	<span class="comment">// or pipe.</span>
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>	if (c.sigFromUser() || flags&amp;_SigPanic == 0) &amp;&amp; sig != _SIGPIPE {
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>		return false
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>	}
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>	<span class="comment">// Determine if the signal occurred inside Go code. We test that:</span>
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>	<span class="comment">//   (1) we weren&#39;t in VDSO page,</span>
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>	<span class="comment">//   (2) we were in a goroutine (i.e., m.curg != nil), and</span>
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>	<span class="comment">//   (3) we weren&#39;t in CGO.</span>
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>	<span class="comment">//   (4) we weren&#39;t in dropped extra m.</span>
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>	gp := sigFetchG(c)
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>	if gp != nil &amp;&amp; gp.m != nil &amp;&amp; gp.m.curg != nil &amp;&amp; !gp.m.isExtraInC &amp;&amp; !gp.m.incgo {
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>		return false
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>	}
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>	<span class="comment">// Signal not handled by Go, forward it.</span>
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>	if fwdFn != _SIG_IGN {
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>		sigfwd(fwdFn, sig, info, ctx)
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>	}
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>	return true
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>}
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span><span class="comment">// sigsave saves the current thread&#39;s signal mask into *p.</span>
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span><span class="comment">// This is used to preserve the non-Go signal mask when a non-Go</span>
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span><span class="comment">// thread calls a Go function.</span>
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span><span class="comment">// This is nosplit and nowritebarrierrec because it is called by needm</span>
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span><span class="comment">// which may be called on a non-Go thread with no g available.</span>
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>func sigsave(p *sigset) {
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>	sigprocmask(_SIG_SETMASK, nil, p)
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>}
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span><span class="comment">// msigrestore sets the current thread&#39;s signal mask to sigmask.</span>
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span><span class="comment">// This is used to restore the non-Go signal mask when a non-Go thread</span>
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span><span class="comment">// calls a Go function.</span>
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span><span class="comment">// This is nosplit and nowritebarrierrec because it is called by dropm</span>
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span><span class="comment">// after g has been cleared.</span>
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>func msigrestore(sigmask sigset) {
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>	sigprocmask(_SIG_SETMASK, &amp;sigmask, nil)
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>}
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span><span class="comment">// sigsetAllExiting is used by sigblock(true) when a thread is</span>
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span><span class="comment">// exiting.</span>
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>var sigsetAllExiting = func() sigset {
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>	res := sigset_all
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>	<span class="comment">// Apply GOOS-specific overrides here, rather than in osinit,</span>
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>	<span class="comment">// because osinit may be called before sigsetAllExiting is</span>
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>	<span class="comment">// initialized (#51913).</span>
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>	if GOOS == &#34;linux&#34; &amp;&amp; iscgo {
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>		<span class="comment">// #42494 glibc and musl reserve some signals for</span>
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>		<span class="comment">// internal use and require they not be blocked by</span>
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>		<span class="comment">// the rest of a normal C runtime. When the go runtime</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>		<span class="comment">// blocks...unblocks signals, temporarily, the blocked</span>
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>		<span class="comment">// interval of time is generally very short. As such,</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>		<span class="comment">// these expectations of *libc code are mostly met by</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>		<span class="comment">// the combined go+cgo system of threads. However,</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>		<span class="comment">// when go causes a thread to exit, via a return from</span>
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>		<span class="comment">// mstart(), the combined runtime can deadlock if</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>		<span class="comment">// these signals are blocked. Thus, don&#39;t block these</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>		<span class="comment">// signals when exiting threads.</span>
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>		<span class="comment">// - glibc: SIGCANCEL (32), SIGSETXID (33)</span>
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>		<span class="comment">// - musl: SIGTIMER (32), SIGCANCEL (33), SIGSYNCCALL (34)</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>		sigdelset(&amp;res, 32)
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>		sigdelset(&amp;res, 33)
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>		sigdelset(&amp;res, 34)
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>	}
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>	return res
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>}()
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span><span class="comment">// sigblock blocks signals in the current thread&#39;s signal mask.</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span><span class="comment">// This is used to block signals while setting up and tearing down g</span>
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span><span class="comment">// when a non-Go thread calls a Go function. When a thread is exiting</span>
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span><span class="comment">// we use the sigsetAllExiting value, otherwise the OS specific</span>
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span><span class="comment">// definition of sigset_all is used.</span>
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span><span class="comment">// This is nosplit and nowritebarrierrec because it is called by needm</span>
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span><span class="comment">// which may be called on a non-Go thread with no g available.</span>
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>func sigblock(exiting bool) {
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>	if exiting {
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>		sigprocmask(_SIG_SETMASK, &amp;sigsetAllExiting, nil)
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>		return
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>	}
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>	sigprocmask(_SIG_SETMASK, &amp;sigset_all, nil)
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>}
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span><span class="comment">// unblocksig removes sig from the current thread&#39;s signal mask.</span>
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span><span class="comment">// This is nosplit and nowritebarrierrec because it is called from</span>
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span><span class="comment">// dieFromSignal, which can be called by sigfwdgo while running in the</span>
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span><span class="comment">// signal handler, on the signal stack, with no g available.</span>
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>func unblocksig(sig uint32) {
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>	var set sigset
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>	sigaddset(&amp;set, int(sig))
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>	sigprocmask(_SIG_UNBLOCK, &amp;set, nil)
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>}
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span><span class="comment">// minitSignals is called when initializing a new m to set the</span>
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span><span class="comment">// thread&#39;s alternate signal stack and signal mask.</span>
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>func minitSignals() {
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>	minitSignalStack()
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>	minitSignalMask()
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>}
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span><span class="comment">// minitSignalStack is called when initializing a new m to set the</span>
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span><span class="comment">// alternate signal stack. If the alternate signal stack is not set</span>
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span><span class="comment">// for the thread (the normal case) then set the alternate signal</span>
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span><span class="comment">// stack to the gsignal stack. If the alternate signal stack is set</span>
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span><span class="comment">// for the thread (the case when a non-Go thread sets the alternate</span>
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span><span class="comment">// signal stack and then calls a Go function) then set the gsignal</span>
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span><span class="comment">// stack to the alternate signal stack. We also set the alternate</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span><span class="comment">// signal stack to the gsignal stack if cgo is not used (regardless</span>
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span><span class="comment">// of whether it is already set). Record which choice was made in</span>
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span><span class="comment">// newSigstack, so that it can be undone in unminit.</span>
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>func minitSignalStack() {
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>	mp := getg().m
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>	var st stackt
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>	sigaltstack(nil, &amp;st)
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>	if st.ss_flags&amp;_SS_DISABLE != 0 || !iscgo {
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>		signalstack(&amp;mp.gsignal.stack)
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>		mp.newSigstack = true
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>	} else {
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>		setGsignalStack(&amp;st, &amp;mp.goSigStack)
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>		mp.newSigstack = false
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>	}
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>}
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span><span class="comment">// minitSignalMask is called when initializing a new m to set the</span>
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span><span class="comment">// thread&#39;s signal mask. When this is called all signals have been</span>
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span><span class="comment">// blocked for the thread.  This starts with m.sigmask, which was set</span>
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span><span class="comment">// either from initSigmask for a newly created thread or by calling</span>
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span><span class="comment">// sigsave if this is a non-Go thread calling a Go function. It</span>
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span><span class="comment">// removes all essential signals from the mask, thus causing those</span>
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span><span class="comment">// signals to not be blocked. Then it sets the thread&#39;s signal mask.</span>
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span><span class="comment">// After this is called the thread can receive signals.</span>
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>func minitSignalMask() {
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>	nmask := getg().m.sigmask
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>	for i := range sigtable {
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>		if !blockableSig(uint32(i)) {
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>			sigdelset(&amp;nmask, i)
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>		}
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>	}
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>	sigprocmask(_SIG_SETMASK, &amp;nmask, nil)
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>}
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span><span class="comment">// unminitSignals is called from dropm, via unminit, to undo the</span>
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span><span class="comment">// effect of calling minit on a non-Go thread.</span>
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>func unminitSignals() {
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>	if getg().m.newSigstack {
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>		st := stackt{ss_flags: _SS_DISABLE}
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>		sigaltstack(&amp;st, nil)
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>	} else {
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>		<span class="comment">// We got the signal stack from someone else. Restore</span>
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>		<span class="comment">// the Go-allocated stack in case this M gets reused</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>		<span class="comment">// for another thread (e.g., it&#39;s an extram). Also, on</span>
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>		<span class="comment">// Android, libc allocates a signal stack for all</span>
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>		<span class="comment">// threads, so it&#39;s important to restore the Go stack</span>
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>		<span class="comment">// even on Go-created threads so we can free it.</span>
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>		restoreGsignalStack(&amp;getg().m.goSigStack)
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>	}
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>}
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span><span class="comment">// blockableSig reports whether sig may be blocked by the signal mask.</span>
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span><span class="comment">// We never want to block the signals marked _SigUnblock;</span>
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span><span class="comment">// these are the synchronous signals that turn into a Go panic.</span>
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span><span class="comment">// We never want to block the preemption signal if it is being used.</span>
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span><span class="comment">// In a Go program--not a c-archive/c-shared--we never want to block</span>
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span><span class="comment">// the signals marked _SigKill or _SigThrow, as otherwise it&#39;s possible</span>
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span><span class="comment">// for all running threads to block them and delay their delivery until</span>
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span><span class="comment">// we start a new thread. When linked into a C program we let the C code</span>
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span><span class="comment">// decide on the disposition of those signals.</span>
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>func blockableSig(sig uint32) bool {
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>	flags := sigtable[sig].flags
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>	if flags&amp;_SigUnblock != 0 {
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>		return false
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>	}
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>	if sig == sigPreempt &amp;&amp; preemptMSupported &amp;&amp; debug.asyncpreemptoff == 0 {
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>		return false
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>	}
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>	if isarchive || islibrary {
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>		return true
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>	}
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>	return flags&amp;(_SigKill|_SigThrow) == 0
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>}
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span><span class="comment">// gsignalStack saves the fields of the gsignal stack changed by</span>
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span><span class="comment">// setGsignalStack.</span>
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>type gsignalStack struct {
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>	stack       stack
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>	stackguard0 uintptr
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>	stackguard1 uintptr
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>	stktopsp    uintptr
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>}
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span><span class="comment">// setGsignalStack sets the gsignal stack of the current m to an</span>
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span><span class="comment">// alternate signal stack returned from the sigaltstack system call.</span>
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span><span class="comment">// It saves the old values in *old for use by restoreGsignalStack.</span>
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span><span class="comment">// This is used when handling a signal if non-Go code has set the</span>
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span><span class="comment">// alternate signal stack.</span>
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>func setGsignalStack(st *stackt, old *gsignalStack) {
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>	gp := getg()
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>	if old != nil {
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>		old.stack = gp.m.gsignal.stack
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>		old.stackguard0 = gp.m.gsignal.stackguard0
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>		old.stackguard1 = gp.m.gsignal.stackguard1
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>		old.stktopsp = gp.m.gsignal.stktopsp
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>	}
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>	stsp := uintptr(unsafe.Pointer(st.ss_sp))
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>	gp.m.gsignal.stack.lo = stsp
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>	gp.m.gsignal.stack.hi = stsp + st.ss_size
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>	gp.m.gsignal.stackguard0 = stsp + stackGuard
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>	gp.m.gsignal.stackguard1 = stsp + stackGuard
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>}
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span><span class="comment">// restoreGsignalStack restores the gsignal stack to the value it had</span>
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span><span class="comment">// before entering the signal handler.</span>
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>func restoreGsignalStack(st *gsignalStack) {
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>	gp := getg().m.gsignal
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>	gp.stack = st.stack
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>	gp.stackguard0 = st.stackguard0
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>	gp.stackguard1 = st.stackguard1
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>	gp.stktopsp = st.stktopsp
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>}
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span><span class="comment">// signalstack sets the current thread&#39;s alternate signal stack to s.</span>
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>func signalstack(s *stack) {
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>	st := stackt{ss_size: s.hi - s.lo}
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>	setSignalstackSP(&amp;st, s.lo)
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>	sigaltstack(&amp;st, nil)
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>}
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span><span class="comment">// setsigsegv is used on darwin/arm64 to fake a segmentation fault.</span>
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span><span class="comment">// This is exported via linkname to assembly in runtime/cgo.</span>
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span><span class="comment">//go:linkname setsigsegv</span>
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>func setsigsegv(pc uintptr) {
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>	gp := getg()
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>	gp.sig = _SIGSEGV
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>	gp.sigpc = pc
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>	gp.sigcode0 = _SEGV_MAPERR
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>	gp.sigcode1 = 0 <span class="comment">// TODO: emulate si_addr</span>
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>}
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>
</pre><p><a href="signal_unix.go?m=text">View as plain text</a></p>

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
