<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/os_linux.go - Go Documentation Server</title>

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
<a href="os_linux.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">os_linux.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;runtime/internal/syscall&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// sigPerThreadSyscall is the same signal (SIGSETXID) used by glibc for</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// per-thread syscalls on Linux. We use it for the same purpose in non-cgo</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// binaries.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>const sigPerThreadSyscall = _SIGRTMIN + 1
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>type mOS struct {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// profileTimer holds the ID of the POSIX interval timer for profiling CPU</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// usage on this thread.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// It is valid when the profileTimerValid field is true. A thread</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// creates and manages its own timer, and these fields are read and written</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// only by this thread. But because some of the reads on profileTimerValid</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// are in signal handling code, this field should be atomic type.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	profileTimer      int32
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	profileTimerValid atomic.Bool
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// needPerThreadSyscall indicates that a per-thread syscall is required</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// for doAllThreadsSyscall.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	needPerThreadSyscall atomic.Uint8
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>func futex(addr unsafe.Pointer, op int32, val uint32, ts, addr2 unsafe.Pointer, val3 uint32) int32
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// Linux futex.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//	futexsleep(uint32 *addr, uint32 val)</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//	futexwakeup(uint32 *addr)</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// Futexsleep atomically checks if *addr == val and if so, sleeps on addr.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// Futexwakeup wakes up threads sleeping on addr.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// Futexsleep is allowed to wake up spuriously.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>const (
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	_FUTEX_PRIVATE_FLAG = 128
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	_FUTEX_WAIT_PRIVATE = 0 | _FUTEX_PRIVATE_FLAG
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	_FUTEX_WAKE_PRIVATE = 1 | _FUTEX_PRIVATE_FLAG
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// Atomically,</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//	if(*addr == val) sleep</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// Might be woken up spuriously; that&#39;s allowed.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// Don&#39;t sleep longer than ns; ns &lt; 0 means forever.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>func futexsleep(addr *uint32, val uint32, ns int64) {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// Some Linux kernels have a bug where futex of</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// FUTEX_WAIT returns an internal error code</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// as an errno. Libpthread ignores the return value</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// here, and so can we: as it says a few lines up,</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// spurious wakeups are allowed.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	if ns &lt; 0 {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		futex(unsafe.Pointer(addr), _FUTEX_WAIT_PRIVATE, val, nil, nil, 0)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		return
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	var ts timespec
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	ts.setNsec(ns)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	futex(unsafe.Pointer(addr), _FUTEX_WAIT_PRIVATE, val, unsafe.Pointer(&amp;ts), nil, 0)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// If any procs are sleeping on addr, wake up at most cnt.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>func futexwakeup(addr *uint32, cnt uint32) {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	ret := futex(unsafe.Pointer(addr), _FUTEX_WAKE_PRIVATE, cnt, nil, nil, 0)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	if ret &gt;= 0 {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		return
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// I don&#39;t know that futex wakeup can return</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// EAGAIN or EINTR, but if it does, it would be</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// safe to loop and call futex again.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		print(&#34;futexwakeup addr=&#34;, addr, &#34; returned &#34;, ret, &#34;\n&#34;)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	})
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	*(*int32)(unsafe.Pointer(uintptr(0x1006))) = 0x1006
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func getproccount() int32 {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// This buffer is huge (8 kB) but we are on the system stack</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// and there should be plenty of space (64 kB).</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// Also this is a leaf, so we&#39;re not holding up the memory for long.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// See golang.org/issue/11823.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// The suggested behavior here is to keep trying with ever-larger</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// buffers, but we don&#39;t have a dynamic memory allocator at the</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// moment, so that&#39;s a bit tricky and seems like overkill.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	const maxCPUs = 64 * 1024
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	var buf [maxCPUs / 8]byte
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	r := sched_getaffinity(0, unsafe.Sizeof(buf), &amp;buf[0])
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	if r &lt; 0 {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		return 1
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	n := int32(0)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	for _, v := range buf[:r] {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		for v != 0 {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			n += int32(v &amp; 1)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			v &gt;&gt;= 1
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	if n == 0 {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		n = 1
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	return n
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// Clone, the Linux rfork.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>const (
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	_CLONE_VM             = 0x100
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	_CLONE_FS             = 0x200
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	_CLONE_FILES          = 0x400
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	_CLONE_SIGHAND        = 0x800
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	_CLONE_PTRACE         = 0x2000
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	_CLONE_VFORK          = 0x4000
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	_CLONE_PARENT         = 0x8000
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	_CLONE_THREAD         = 0x10000
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	_CLONE_NEWNS          = 0x20000
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	_CLONE_SYSVSEM        = 0x40000
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	_CLONE_SETTLS         = 0x80000
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	_CLONE_PARENT_SETTID  = 0x100000
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	_CLONE_CHILD_CLEARTID = 0x200000
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	_CLONE_UNTRACED       = 0x800000
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	_CLONE_CHILD_SETTID   = 0x1000000
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	_CLONE_STOPPED        = 0x2000000
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	_CLONE_NEWUTS         = 0x4000000
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	_CLONE_NEWIPC         = 0x8000000
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// As of QEMU 2.8.0 (5ea2fc84d), user emulation requires all six of these</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// flags to be set when creating a thread; attempts to share the other</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">// five but leave SYSVSEM unshared will fail with -EINVAL.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// In non-QEMU environments CLONE_SYSVSEM is inconsequential as we do not</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// use System V semaphores.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	cloneFlags = _CLONE_VM | <span class="comment">/* share memory */</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		_CLONE_FS | <span class="comment">/* share cwd, etc */</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		_CLONE_FILES | <span class="comment">/* share fd table */</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		_CLONE_SIGHAND | <span class="comment">/* share sig handler table */</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		_CLONE_SYSVSEM | <span class="comment">/* share SysV semaphore undo lists (see issue #20763) */</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		_CLONE_THREAD <span class="comment">/* revisit - okay for now */</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>func clone(flags int32, stk, mp, gp, fn unsafe.Pointer) int32
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// May run with m.p==nil, so write barriers are not allowed.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrier</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>func newosproc(mp *m) {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	stk := unsafe.Pointer(mp.g0.stack.hi)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">/*
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	 * note: strace gets confused if we use CLONE_PTRACE here.
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	 */</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	if false {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		print(&#34;newosproc stk=&#34;, stk, &#34; m=&#34;, mp, &#34; g=&#34;, mp.g0, &#34; clone=&#34;, abi.FuncPCABI0(clone), &#34; id=&#34;, mp.id, &#34; ostk=&#34;, &amp;mp, &#34;\n&#34;)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// Disable signals during clone, so that the new thread starts</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// with signals disabled. It will enable them in minit.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	var oset sigset
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	sigprocmask(_SIG_SETMASK, &amp;sigset_all, &amp;oset)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	ret := retryOnEAGAIN(func() int32 {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		r := clone(cloneFlags, stk, unsafe.Pointer(mp), unsafe.Pointer(mp.g0), unsafe.Pointer(abi.FuncPCABI0(mstart)))
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		<span class="comment">// clone returns positive TID, negative errno.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		<span class="comment">// We don&#39;t care about the TID.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		if r &gt;= 0 {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			return 0
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		return -r
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	})
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	sigprocmask(_SIG_SETMASK, &amp;oset, nil)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	if ret != 0 {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		print(&#34;runtime: failed to create new OS thread (have &#34;, mcount(), &#34; already; errno=&#34;, ret, &#34;)\n&#34;)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		if ret == _EAGAIN {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			println(&#34;runtime: may need to increase max user processes (ulimit -u)&#34;)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		throw(&#34;newosproc&#34;)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// Version of newosproc that doesn&#39;t require a valid G.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>func newosproc0(stacksize uintptr, fn unsafe.Pointer) {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	stack := sysAlloc(stacksize, &amp;memstats.stacks_sys)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	if stack == nil {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		writeErrStr(failallocatestack)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		exit(1)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	ret := clone(cloneFlags, unsafe.Pointer(uintptr(stack)+stacksize), nil, nil, fn)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	if ret &lt; 0 {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		writeErrStr(failthreadcreate)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		exit(1)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>const (
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	_AT_NULL     = 0  <span class="comment">// End of vector</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	_AT_PAGESZ   = 6  <span class="comment">// System physical page size</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	_AT_PLATFORM = 15 <span class="comment">// string identifying platform</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	_AT_HWCAP    = 16 <span class="comment">// hardware capability bit vector</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	_AT_SECURE   = 23 <span class="comment">// secure mode boolean</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	_AT_RANDOM   = 25 <span class="comment">// introduced in 2.6.29</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	_AT_HWCAP2   = 26 <span class="comment">// hardware capability bit vector 2</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>var procAuxv = []byte(&#34;/proc/self/auxv\x00&#34;)
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>var addrspace_vec [1]byte
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>func mincore(addr unsafe.Pointer, n uintptr, dst *byte) int32
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>var auxvreadbuf [128]uintptr
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>func sysargs(argc int32, argv **byte) {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	n := argc + 1
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// skip over argv, envp to get to auxv</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	for argv_index(argv, n) != nil {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		n++
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// skip NULL separator</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	n++
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// now argv+n is auxv</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	auxvp := (*[1 &lt;&lt; 28]uintptr)(add(unsafe.Pointer(argv), uintptr(n)*goarch.PtrSize))
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	if pairs := sysauxv(auxvp[:]); pairs != 0 {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		auxv = auxvp[: pairs*2 : pairs*2]
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		return
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// In some situations we don&#39;t get a loader-provided</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// auxv, such as when loaded as a library on Android.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// Fall back to /proc/self/auxv.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	fd := open(&amp;procAuxv[0], 0 <span class="comment">/* O_RDONLY */</span>, 0)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	if fd &lt; 0 {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		<span class="comment">// On Android, /proc/self/auxv might be unreadable (issue 9229), so we fallback to</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		<span class="comment">// try using mincore to detect the physical page size.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		<span class="comment">// mincore should return EINVAL when address is not a multiple of system page size.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		const size = 256 &lt;&lt; 10 <span class="comment">// size of memory region to allocate</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		p, err := mmap(nil, size, _PROT_READ|_PROT_WRITE, _MAP_ANON|_MAP_PRIVATE, -1, 0)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		if err != 0 {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			return
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		var n uintptr
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		for n = 4 &lt;&lt; 10; n &lt; size; n &lt;&lt;= 1 {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			err := mincore(unsafe.Pointer(uintptr(p)+n), 1, &amp;addrspace_vec[0])
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			if err == 0 {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				physPageSize = n
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>				break
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		if physPageSize == 0 {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			physPageSize = size
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		munmap(p, size)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		return
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	n = read(fd, noescape(unsafe.Pointer(&amp;auxvreadbuf[0])), int32(unsafe.Sizeof(auxvreadbuf)))
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	closefd(fd)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	if n &lt; 0 {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		return
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	<span class="comment">// Make sure buf is terminated, even if we didn&#39;t read</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">// the whole file.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	auxvreadbuf[len(auxvreadbuf)-2] = _AT_NULL
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	pairs := sysauxv(auxvreadbuf[:])
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	auxv = auxvreadbuf[: pairs*2 : pairs*2]
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span><span class="comment">// secureMode holds the value of AT_SECURE passed in the auxiliary vector.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>var secureMode bool
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>func sysauxv(auxv []uintptr) (pairs int) {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	var i int
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	for ; auxv[i] != _AT_NULL; i += 2 {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		tag, val := auxv[i], auxv[i+1]
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		switch tag {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		case _AT_RANDOM:
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			<span class="comment">// The kernel provides a pointer to 16-bytes</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			<span class="comment">// worth of random data.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			startupRand = (*[16]byte)(unsafe.Pointer(val))[:]
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		case _AT_PAGESZ:
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			physPageSize = val
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		case _AT_SECURE:
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>			secureMode = val == 1
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		archauxv(tag, val)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		vdsoauxv(tag, val)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	return i / 2
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>var sysTHPSizePath = []byte(&#34;/sys/kernel/mm/transparent_hugepage/hpage_pmd_size\x00&#34;)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>func getHugePageSize() uintptr {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	var numbuf [20]byte
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	fd := open(&amp;sysTHPSizePath[0], 0 <span class="comment">/* O_RDONLY */</span>, 0)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	if fd &lt; 0 {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		return 0
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	ptr := noescape(unsafe.Pointer(&amp;numbuf[0]))
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	n := read(fd, ptr, int32(len(numbuf)))
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	closefd(fd)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if n &lt;= 0 {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		return 0
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	n-- <span class="comment">// remove trailing newline</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	v, ok := atoi(slicebytetostringtmp((*byte)(ptr), int(n)))
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	if !ok || v &lt; 0 {
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		v = 0
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	if v&amp;(v-1) != 0 {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		<span class="comment">// v is not a power of 2</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		return 0
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	return uintptr(v)
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>func osinit() {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	ncpu = getproccount()
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	physHugePageSize = getHugePageSize()
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	osArchInit()
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>var urandom_dev = []byte(&#34;/dev/urandom\x00&#34;)
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>func readRandom(r []byte) int {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	fd := open(&amp;urandom_dev[0], 0 <span class="comment">/* O_RDONLY */</span>, 0)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	n := read(fd, unsafe.Pointer(&amp;r[0]), int32(len(r)))
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	closefd(fd)
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	return int(n)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>func goenvs() {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	goenvs_unix()
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// Called to do synchronous initialization of Go code built with</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span><span class="comment">// -buildmode=c-archive or -buildmode=c-shared.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">// None of the Go runtime is initialized.</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>func libpreinit() {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	initsig(true)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span><span class="comment">// Called to initialize a new m (including the bootstrap m).</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span><span class="comment">// Called on the parent thread (main thread in case of bootstrap), can allocate memory.</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>func mpreinit(mp *m) {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	mp.gsignal = malg(32 * 1024) <span class="comment">// Linux wants &gt;= 2K</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	mp.gsignal.m = mp
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>func gettid() uint32
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span><span class="comment">// Called to initialize a new m (including the bootstrap m).</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span><span class="comment">// Called on the new thread, cannot allocate memory.</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>func minit() {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	minitSignals()
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	<span class="comment">// Cgo-created threads and the bootstrap m are missing a</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	<span class="comment">// procid. We need this for asynchronous preemption and it&#39;s</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	<span class="comment">// useful in debuggers.</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	getg().m.procid = uint64(gettid())
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>}
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span><span class="comment">// Called from dropm to undo the effect of an minit.</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>func unminit() {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	unminitSignals()
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	getg().m.procid = 0
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span><span class="comment">// Called from exitm, but not from drop, to undo the effect of thread-owned</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span><span class="comment">// resources in minit, semacreate, or elsewhere. Do not take locks after calling this.</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>func mdestroy(mp *m) {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>}
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span><span class="comment">//#ifdef GOARCH_386</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span><span class="comment">//#define sa_handler k_sa_handler</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span><span class="comment">//#endif</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>func sigreturn__sigaction()
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>func sigtramp() <span class="comment">// Called via C ABI</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>func cgoSigtramp()
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>func sigaltstack(new, old *stackt)
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>func setitimer(mode int32, new, old *itimerval)
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>func timer_create(clockid int32, sevp *sigevent, timerid *int32) int32
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>func timer_settime(timerid int32, flags int32, new, old *itimerspec) int32
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>func timer_delete(timerid int32) int32
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>func rtsigprocmask(how int32, new, old *sigset, size int32)
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>func sigprocmask(how int32, new, old *sigset) {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	rtsigprocmask(how, new, old, int32(unsafe.Sizeof(*new)))
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>func raise(sig uint32)
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>func raiseproc(sig uint32)
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>func sched_getaffinity(pid, len uintptr, buf *byte) int32
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>func osyield()
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>func osyield_no_g() {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	osyield()
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>func pipe2(flags int32) (r, w int32, errno int32)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>func fcntl(fd, cmd, arg int32) (ret int32, errno int32) {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	r, _, err := syscall.Syscall6(syscall.SYS_FCNTL, uintptr(fd), uintptr(cmd), uintptr(arg), 0, 0, 0)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	return int32(r), int32(err)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>const (
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	_si_max_size    = 128
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	_sigev_max_size = 64
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>)
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>func setsig(i uint32, fn uintptr) {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	var sa sigactiont
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	sa.sa_flags = _SA_SIGINFO | _SA_ONSTACK | _SA_RESTORER | _SA_RESTART
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	sigfillset(&amp;sa.sa_mask)
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	<span class="comment">// Although Linux manpage says &#34;sa_restorer element is obsolete and</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	<span class="comment">// should not be used&#34;. x86_64 kernel requires it. Only use it on</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	<span class="comment">// x86.</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	if GOARCH == &#34;386&#34; || GOARCH == &#34;amd64&#34; {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		sa.sa_restorer = abi.FuncPCABI0(sigreturn__sigaction)
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	if fn == abi.FuncPCABIInternal(sighandler) { <span class="comment">// abi.FuncPCABIInternal(sighandler) matches the callers in signal_unix.go</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		if iscgo {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			fn = abi.FuncPCABI0(cgoSigtramp)
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		} else {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			fn = abi.FuncPCABI0(sigtramp)
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		}
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	sa.sa_handler = fn
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	sigaction(i, &amp;sa, nil)
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>}
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>func setsigstack(i uint32) {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	var sa sigactiont
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	sigaction(i, nil, &amp;sa)
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	if sa.sa_flags&amp;_SA_ONSTACK != 0 {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		return
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	sa.sa_flags |= _SA_ONSTACK
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	sigaction(i, &amp;sa, nil)
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span><span class="comment">//go:nowritebarrierrec</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>func getsig(i uint32) uintptr {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	var sa sigactiont
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	sigaction(i, nil, &amp;sa)
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	return sa.sa_handler
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span><span class="comment">// setSignalstackSP sets the ss_sp field of a stackt.</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>func setSignalstackSP(s *stackt, sp uintptr) {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	*(*uintptr)(unsafe.Pointer(&amp;s.ss_sp)) = sp
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>func (c *sigctxt) fixsigcode(sig uint32) {
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span><span class="comment">// sysSigaction calls the rt_sigaction system call.</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>func sysSigaction(sig uint32, new, old *sigactiont) {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	if rt_sigaction(uintptr(sig), new, old, unsafe.Sizeof(sigactiont{}.sa_mask)) != 0 {
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		<span class="comment">// Workaround for bugs in QEMU user mode emulation.</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		<span class="comment">// QEMU turns calls to the sigaction system call into</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		<span class="comment">// calls to the C library sigaction call; the C</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		<span class="comment">// library call rejects attempts to call sigaction for</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		<span class="comment">// SIGCANCEL (32) or SIGSETXID (33).</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		<span class="comment">// QEMU rejects calling sigaction on SIGRTMAX (64).</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		<span class="comment">// Just ignore the error in these case. There isn&#39;t</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		<span class="comment">// anything we can do about it anyhow.</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		if sig != 32 &amp;&amp; sig != 33 &amp;&amp; sig != 64 {
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>			<span class="comment">// Use system stack to avoid split stack overflow on ppc64/ppc64le.</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>			systemstack(func() {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>				throw(&#34;sigaction failed&#34;)
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			})
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	}
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">// rt_sigaction is implemented in assembly.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">//go:noescape</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>func rt_sigaction(sig uintptr, new, old *sigactiont, size uintptr) int32
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>func getpid() int
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>func tgkill(tgid, tid, sig int)
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span><span class="comment">// signalM sends a signal to mp.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>func signalM(mp *m, sig int) {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	tgkill(getpid(), int(mp.procid), sig)
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>}
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span><span class="comment">// validSIGPROF compares this signal delivery&#39;s code against the signal sources</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span><span class="comment">// that the profiler uses, returning whether the delivery should be processed.</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span><span class="comment">// To be processed, a signal delivery from a known profiling mechanism should</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span><span class="comment">// correspond to the best profiling mechanism available to this thread. Signals</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span><span class="comment">// from other sources are always considered valid.</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>func validSIGPROF(mp *m, c *sigctxt) bool {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	code := int32(c.sigcode())
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	setitimer := code == _SI_KERNEL
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	timer_create := code == _SI_TIMER
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	if !(setitimer || timer_create) {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		<span class="comment">// The signal doesn&#39;t correspond to a profiling mechanism that the</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		<span class="comment">// runtime enables itself. There&#39;s no reason to process it, but there&#39;s</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		<span class="comment">// no reason to ignore it either.</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		return true
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	if mp == nil {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		<span class="comment">// Since we don&#39;t have an M, we can&#39;t check if there&#39;s an active</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		<span class="comment">// per-thread timer for this thread. We don&#39;t know how long this thread</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		<span class="comment">// has been around, and if it happened to interact with the Go scheduler</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		<span class="comment">// at a time when profiling was active (causing it to have a per-thread</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		<span class="comment">// timer). But it may have never interacted with the Go scheduler, or</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		<span class="comment">// never while profiling was active. To avoid double-counting, process</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		<span class="comment">// only signals from setitimer.</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		<span class="comment">// When a custom cgo traceback function has been registered (on</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		<span class="comment">// platforms that support runtime.SetCgoTraceback), SIGPROF signals</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		<span class="comment">// delivered to a thread that cannot find a matching M do this check in</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		<span class="comment">// the assembly implementations of runtime.cgoSigtramp.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		return setitimer
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	<span class="comment">// Having an M means the thread interacts with the Go scheduler, and we can</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	<span class="comment">// check whether there&#39;s an active per-thread timer for this thread.</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	if mp.profileTimerValid.Load() {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		<span class="comment">// If this M has its own per-thread CPU profiling interval timer, we</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		<span class="comment">// should track the SIGPROF signals that come from that timer (for</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		<span class="comment">// accurate reporting of its CPU usage; see issue 35057) and ignore any</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		<span class="comment">// that it gets from the process-wide setitimer (to not over-count its</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		<span class="comment">// CPU consumption).</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		return timer_create
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	<span class="comment">// No active per-thread timer means the only valid profiler is setitimer.</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	return setitimer
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>func setProcessCPUProfiler(hz int32) {
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	setProcessCPUProfilerTimer(hz)
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>func setThreadCPUProfiler(hz int32) {
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	mp := getg().m
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	mp.profilehz = hz
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	<span class="comment">// destroy any active timer</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	if mp.profileTimerValid.Load() {
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		timerid := mp.profileTimer
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		mp.profileTimerValid.Store(false)
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		mp.profileTimer = 0
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		ret := timer_delete(timerid)
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		if ret != 0 {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			print(&#34;runtime: failed to disable profiling timer; timer_delete(&#34;, timerid, &#34;) errno=&#34;, -ret, &#34;\n&#34;)
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			throw(&#34;timer_delete&#34;)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	if hz == 0 {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		<span class="comment">// If the goal was to disable profiling for this thread, then the job&#39;s done.</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		return
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	}
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	<span class="comment">// The period of the timer should be 1/Hz. For every &#34;1/Hz&#34; of additional</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	<span class="comment">// work, the user should expect one additional sample in the profile.</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	<span class="comment">// But to scale down to very small amounts of application work, to observe</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	<span class="comment">// even CPU usage of &#34;one tenth&#34; of the requested period, set the initial</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	<span class="comment">// timing delay in a different way: So that &#34;one tenth&#34; of a period of CPU</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	<span class="comment">// spend shows up as a 10% chance of one sample (for an expected value of</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	<span class="comment">// 0.1 samples), and so that &#34;two and six tenths&#34; periods of CPU spend show</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	<span class="comment">// up as a 60% chance of 3 samples and a 40% chance of 2 samples (for an</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	<span class="comment">// expected value of 2.6). Set the initial delay to a value in the unifom</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	<span class="comment">// random distribution between 0 and the desired period. And because &#34;0&#34;</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	<span class="comment">// means &#34;disable timer&#34;, add 1 so the half-open interval [0,period) turns</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	<span class="comment">// into (0,period].</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	<span class="comment">// Otherwise, this would show up as a bias away from short-lived threads and</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	<span class="comment">// from threads that are only occasionally active: for example, when the</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	<span class="comment">// garbage collector runs on a mostly-idle system, the additional threads it</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	<span class="comment">// activates may do a couple milliseconds of GC-related work and nothing</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	<span class="comment">// else in the few seconds that the profiler observes.</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	spec := new(itimerspec)
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	spec.it_value.setNsec(1 + int64(cheaprandn(uint32(1e9/hz))))
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	spec.it_interval.setNsec(1e9 / int64(hz))
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	var timerid int32
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	var sevp sigevent
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	sevp.notify = _SIGEV_THREAD_ID
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	sevp.signo = _SIGPROF
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	sevp.sigev_notify_thread_id = int32(mp.procid)
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	ret := timer_create(_CLOCK_THREAD_CPUTIME_ID, &amp;sevp, &amp;timerid)
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	if ret != 0 {
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		<span class="comment">// If we cannot create a timer for this M, leave profileTimerValid false</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		<span class="comment">// to fall back to the process-wide setitimer profiler.</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		return
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	}
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	ret = timer_settime(timerid, 0, spec, nil)
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	if ret != 0 {
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		print(&#34;runtime: failed to configure profiling timer; timer_settime(&#34;, timerid,
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>			&#34;, 0, {interval: {&#34;,
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>			spec.it_interval.tv_sec, &#34;s + &#34;, spec.it_interval.tv_nsec, &#34;ns} value: {&#34;,
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>			spec.it_value.tv_sec, &#34;s + &#34;, spec.it_value.tv_nsec, &#34;ns}}, nil) errno=&#34;, -ret, &#34;\n&#34;)
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>		throw(&#34;timer_settime&#34;)
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	}
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	mp.profileTimer = timerid
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	mp.profileTimerValid.Store(true)
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>}
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span><span class="comment">// perThreadSyscallArgs contains the system call number, arguments, and</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span><span class="comment">// expected return values for a system call to be executed on all threads.</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>type perThreadSyscallArgs struct {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	trap uintptr
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	a1   uintptr
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	a2   uintptr
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	a3   uintptr
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	a4   uintptr
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	a5   uintptr
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	a6   uintptr
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	r1   uintptr
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	r2   uintptr
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>}
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span><span class="comment">// perThreadSyscall is the system call to execute for the ongoing</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span><span class="comment">// doAllThreadsSyscall.</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span><span class="comment">// perThreadSyscall may only be written while mp.needPerThreadSyscall == 0 on</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span><span class="comment">// all Ms.</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>var perThreadSyscall perThreadSyscallArgs
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span><span class="comment">// syscall_runtime_doAllThreadsSyscall and executes a specified system call on</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span><span class="comment">// all Ms.</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span><span class="comment">// The system call is expected to succeed and return the same value on every</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span><span class="comment">// thread. If any threads do not match, the runtime throws.</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span><span class="comment">//go:linkname syscall_runtime_doAllThreadsSyscall syscall.runtime_doAllThreadsSyscall</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span><span class="comment">//go:uintptrescapes</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>func syscall_runtime_doAllThreadsSyscall(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	if iscgo {
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		<span class="comment">// In cgo, we are not aware of threads created in C, so this approach will not work.</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>		panic(&#34;doAllThreadsSyscall not supported with cgo enabled&#34;)
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	<span class="comment">// STW to guarantee that user goroutines see an atomic change to thread</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	<span class="comment">// state. Without STW, goroutines could migrate Ms while change is in</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	<span class="comment">// progress and e.g., see state old -&gt; new -&gt; old -&gt; new.</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	<span class="comment">// N.B. Internally, this function does not depend on STW to</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	<span class="comment">// successfully change every thread. It is only needed for user</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	<span class="comment">// expectations, per above.</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	stw := stopTheWorld(stwAllThreadsSyscall)
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	<span class="comment">// This function depends on several properties:</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	<span class="comment">// 1. All OS threads that already exist are associated with an M in</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	<span class="comment">//    allm. i.e., we won&#39;t miss any pre-existing threads.</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	<span class="comment">// 2. All Ms listed in allm will eventually have an OS thread exist.</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>	<span class="comment">//    i.e., they will set procid and be able to receive signals.</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	<span class="comment">// 3. OS threads created after we read allm will clone from a thread</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	<span class="comment">//    that has executed the system call. i.e., they inherit the</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	<span class="comment">//    modified state.</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	<span class="comment">// We achieve these through different mechanisms:</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	<span class="comment">// 1. Addition of new Ms to allm in allocm happens before clone of its</span>
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	<span class="comment">//    OS thread later in newm.</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	<span class="comment">// 2. newm does acquirem to avoid being preempted, ensuring that new Ms</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	<span class="comment">//    created in allocm will eventually reach OS thread clone later in</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	<span class="comment">//    newm.</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	<span class="comment">// 3. We take allocmLock for write here to prevent allocation of new Ms</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	<span class="comment">//    while this function runs. Per (1), this prevents clone of OS</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	<span class="comment">//    threads that are not yet in allm.</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	allocmLock.lock()
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	<span class="comment">// Disable preemption, preventing us from changing Ms, as we handle</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	<span class="comment">// this M specially.</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	<span class="comment">// N.B. STW and lock() above do this as well, this is added for extra</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	<span class="comment">// clarity.</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	acquirem()
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	<span class="comment">// N.B. allocmLock also prevents concurrent execution of this function,</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	<span class="comment">// serializing use of perThreadSyscall, mp.needPerThreadSyscall, and</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	<span class="comment">// ensuring all threads execute system calls from multiple calls in the</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	<span class="comment">// same order.</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	r1, r2, errno := syscall.Syscall6(trap, a1, a2, a3, a4, a5, a6)
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	if GOARCH == &#34;ppc64&#34; || GOARCH == &#34;ppc64le&#34; {
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		<span class="comment">// TODO(https://go.dev/issue/51192 ): ppc64 doesn&#39;t use r2.</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		r2 = 0
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	if errno != 0 {
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		releasem(getg().m)
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		allocmLock.unlock()
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>		startTheWorld(stw)
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		return r1, r2, errno
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	}
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	perThreadSyscall = perThreadSyscallArgs{
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		trap: trap,
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		a1:   a1,
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>		a2:   a2,
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>		a3:   a3,
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		a4:   a4,
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>		a5:   a5,
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		a6:   a6,
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		r1:   r1,
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>		r2:   r2,
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>	}
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	<span class="comment">// Wait for all threads to start.</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	<span class="comment">// As described above, some Ms have been added to allm prior to</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	<span class="comment">// allocmLock, but not yet completed OS clone and set procid.</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	<span class="comment">// At minimum we must wait for a thread to set procid before we can</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	<span class="comment">// send it a signal.</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>	<span class="comment">// We take this one step further and wait for all threads to start</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	<span class="comment">// before sending any signals. This prevents system calls from getting</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	<span class="comment">// applied twice: once in the parent and once in the child, like so:</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	<span class="comment">//          A                     B                  C</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	<span class="comment">//                         add C to allm</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	<span class="comment">// doAllThreadsSyscall</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	<span class="comment">//   allocmLock.lock()</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	<span class="comment">//   signal B</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	<span class="comment">//                         &lt;receive signal&gt;</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	<span class="comment">//                         execute syscall</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>	<span class="comment">//                         &lt;signal return&gt;</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	<span class="comment">//                         clone C</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	<span class="comment">//                                             &lt;thread start&gt;</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	<span class="comment">//                                             set procid</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	<span class="comment">//   signal C</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	<span class="comment">//                                             &lt;receive signal&gt;</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	<span class="comment">//                                             execute syscall</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	<span class="comment">//                                             &lt;signal return&gt;</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	<span class="comment">// In this case, thread C inherited the syscall-modified state from</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	<span class="comment">// thread B and did not need to execute the syscall, but did anyway</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	<span class="comment">// because doAllThreadsSyscall could not be sure whether it was</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	<span class="comment">// required.</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	<span class="comment">// Some system calls may not be idempotent, so we ensure each thread</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	<span class="comment">// executes the system call exactly once.</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	for mp := allm; mp != nil; mp = mp.alllink {
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		for atomic.Load64(&amp;mp.procid) == 0 {
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>			<span class="comment">// Thread is starting.</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>			osyield()
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		}
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	}
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	<span class="comment">// Signal every other thread, where they will execute perThreadSyscall</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	<span class="comment">// from the signal handler.</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	gp := getg()
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	tid := gp.m.procid
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	for mp := allm; mp != nil; mp = mp.alllink {
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		if atomic.Load64(&amp;mp.procid) == tid {
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>			<span class="comment">// Our thread already performed the syscall.</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>			continue
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		}
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		mp.needPerThreadSyscall.Store(1)
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		signalM(mp, sigPerThreadSyscall)
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	}
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	<span class="comment">// Wait for all threads to complete.</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>	for mp := allm; mp != nil; mp = mp.alllink {
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		if mp.procid == tid {
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>			continue
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		}
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		for mp.needPerThreadSyscall.Load() != 0 {
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>			osyield()
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	perThreadSyscall = perThreadSyscallArgs{}
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	releasem(getg().m)
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	allocmLock.unlock()
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	startTheWorld(stw)
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	return r1, r2, errno
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>}
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span><span class="comment">// runPerThreadSyscall runs perThreadSyscall for this M if required.</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span><span class="comment">// This function throws if the system call returns with anything other than the</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span><span class="comment">// expected values.</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>func runPerThreadSyscall() {
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	gp := getg()
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	if gp.m.needPerThreadSyscall.Load() == 0 {
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		return
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	}
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	args := perThreadSyscall
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	r1, r2, errno := syscall.Syscall6(args.trap, args.a1, args.a2, args.a3, args.a4, args.a5, args.a6)
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	if GOARCH == &#34;ppc64&#34; || GOARCH == &#34;ppc64le&#34; {
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>		<span class="comment">// TODO(https://go.dev/issue/51192 ): ppc64 doesn&#39;t use r2.</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>		r2 = 0
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	}
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	if errno != 0 || r1 != args.r1 || r2 != args.r2 {
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		print(&#34;trap:&#34;, args.trap, &#34;, a123456=[&#34;, args.a1, &#34;,&#34;, args.a2, &#34;,&#34;, args.a3, &#34;,&#34;, args.a4, &#34;,&#34;, args.a5, &#34;,&#34;, args.a6, &#34;]\n&#34;)
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		print(&#34;results: got {r1=&#34;, r1, &#34;,r2=&#34;, r2, &#34;,errno=&#34;, errno, &#34;}, want {r1=&#34;, args.r1, &#34;,r2=&#34;, args.r2, &#34;,errno=0}\n&#34;)
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>		fatal(&#34;AllThreadsSyscall6 results differ between threads; runtime corrupted&#34;)
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	}
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	gp.m.needPerThreadSyscall.Store(0)
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>}
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>const (
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	_SI_USER  = 0
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	_SI_TKILL = -6
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>)
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span><span class="comment">// sigFromUser reports whether the signal was sent because of a call</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span><span class="comment">// to kill or tgkill.</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>func (c *sigctxt) sigFromUser() bool {
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>	code := int32(c.sigcode())
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	return code == _SI_USER || code == _SI_TKILL
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>}
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
</pre><p><a href="os_linux.go?m=text">View as plain text</a></p>

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
