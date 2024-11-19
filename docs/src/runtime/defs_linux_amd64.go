<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/defs_linux_amd64.go - Go Documentation Server</title>

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
<a href="defs_linux_amd64.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">defs_linux_amd64.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// created by cgo -cdefs and then converted to Go</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// cgo -cdefs defs_linux.go defs1_linux.go</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>package runtime
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>import &#34;unsafe&#34;
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>const (
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	_EINTR  = 0x4
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	_EAGAIN = 0xb
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	_ENOMEM = 0xc
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	_PROT_NONE  = 0x0
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	_PROT_READ  = 0x1
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	_PROT_WRITE = 0x2
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	_PROT_EXEC  = 0x4
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	_MAP_ANON    = 0x20
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	_MAP_PRIVATE = 0x2
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	_MAP_FIXED   = 0x10
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	_MADV_DONTNEED   = 0x4
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	_MADV_FREE       = 0x8
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	_MADV_HUGEPAGE   = 0xe
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	_MADV_NOHUGEPAGE = 0xf
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	_MADV_COLLAPSE   = 0x19
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	_SA_RESTART  = 0x10000000
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	_SA_ONSTACK  = 0x8000000
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	_SA_RESTORER = 0x4000000
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	_SA_SIGINFO  = 0x4
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	_SI_KERNEL = 0x80
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	_SI_TIMER  = -0x2
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	_SIGHUP    = 0x1
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	_SIGINT    = 0x2
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	_SIGQUIT   = 0x3
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	_SIGILL    = 0x4
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	_SIGTRAP   = 0x5
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	_SIGABRT   = 0x6
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	_SIGBUS    = 0x7
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	_SIGFPE    = 0x8
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	_SIGKILL   = 0x9
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	_SIGUSR1   = 0xa
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	_SIGSEGV   = 0xb
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	_SIGUSR2   = 0xc
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	_SIGPIPE   = 0xd
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	_SIGALRM   = 0xe
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	_SIGSTKFLT = 0x10
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	_SIGCHLD   = 0x11
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	_SIGCONT   = 0x12
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	_SIGSTOP   = 0x13
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	_SIGTSTP   = 0x14
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	_SIGTTIN   = 0x15
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	_SIGTTOU   = 0x16
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	_SIGURG    = 0x17
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	_SIGXCPU   = 0x18
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	_SIGXFSZ   = 0x19
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	_SIGVTALRM = 0x1a
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	_SIGPROF   = 0x1b
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	_SIGWINCH  = 0x1c
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	_SIGIO     = 0x1d
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	_SIGPWR    = 0x1e
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	_SIGSYS    = 0x1f
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	_SIGRTMIN = 0x20
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	_FPE_INTDIV = 0x1
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	_FPE_INTOVF = 0x2
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	_FPE_FLTDIV = 0x3
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	_FPE_FLTOVF = 0x4
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	_FPE_FLTUND = 0x5
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	_FPE_FLTRES = 0x6
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	_FPE_FLTINV = 0x7
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	_FPE_FLTSUB = 0x8
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	_BUS_ADRALN = 0x1
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	_BUS_ADRERR = 0x2
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	_BUS_OBJERR = 0x3
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	_SEGV_MAPERR = 0x1
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	_SEGV_ACCERR = 0x2
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	_ITIMER_REAL    = 0x0
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	_ITIMER_VIRTUAL = 0x1
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	_ITIMER_PROF    = 0x2
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	_CLOCK_THREAD_CPUTIME_ID = 0x3
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	_SIGEV_THREAD_ID = 0x4
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	_AF_UNIX    = 0x1
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	_SOCK_DGRAM = 0x2
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>type timespec struct {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	tv_sec  int64
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	tv_nsec int64
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>func (ts *timespec) setNsec(ns int64) {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	ts.tv_sec = ns / 1e9
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	ts.tv_nsec = ns % 1e9
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>type timeval struct {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	tv_sec  int64
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	tv_usec int64
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func (tv *timeval) set_usec(x int32) {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	tv.tv_usec = int64(x)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>type sigactiont struct {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	sa_handler  uintptr
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	sa_flags    uint64
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	sa_restorer uintptr
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	sa_mask     uint64
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>type siginfoFields struct {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	si_signo int32
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	si_errno int32
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	si_code  int32
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// below here is a union; si_addr is the only field we use</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	si_addr uint64
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>type siginfo struct {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	siginfoFields
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// Pad struct to the max size in the kernel.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	_ [_si_max_size - unsafe.Sizeof(siginfoFields{})]byte
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>type itimerspec struct {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	it_interval timespec
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	it_value    timespec
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>type itimerval struct {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	it_interval timeval
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	it_value    timeval
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>type sigeventFields struct {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	value  uintptr
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	signo  int32
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	notify int32
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// below here is a union; sigev_notify_thread_id is the only field we use</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	sigev_notify_thread_id int32
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>type sigevent struct {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	sigeventFields
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// Pad struct to the max size in the kernel.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	_ [_sigev_max_size - unsafe.Sizeof(sigeventFields{})]byte
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// created by cgo -cdefs and then converted to Go</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">// cgo -cdefs defs_linux.go defs1_linux.go</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>const (
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	_O_RDONLY   = 0x0
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	_O_WRONLY   = 0x1
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	_O_CREAT    = 0x40
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	_O_TRUNC    = 0x200
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	_O_NONBLOCK = 0x800
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	_O_CLOEXEC  = 0x80000
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>type usigset struct {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	__val [16]uint64
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>type fpxreg struct {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	significand [4]uint16
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	exponent    uint16
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	padding     [3]uint16
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>type xmmreg struct {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	element [4]uint32
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>type fpstate struct {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	cwd       uint16
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	swd       uint16
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	ftw       uint16
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	fop       uint16
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	rip       uint64
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	rdp       uint64
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	mxcsr     uint32
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	mxcr_mask uint32
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	_st       [8]fpxreg
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	_xmm      [16]xmmreg
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	padding   [24]uint32
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>type fpxreg1 struct {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	significand [4]uint16
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	exponent    uint16
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	padding     [3]uint16
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>type xmmreg1 struct {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	element [4]uint32
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>type fpstate1 struct {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	cwd       uint16
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	swd       uint16
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	ftw       uint16
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	fop       uint16
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	rip       uint64
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	rdp       uint64
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	mxcsr     uint32
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	mxcr_mask uint32
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	_st       [8]fpxreg1
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	_xmm      [16]xmmreg1
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	padding   [24]uint32
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>type fpreg1 struct {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	significand [4]uint16
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	exponent    uint16
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>type stackt struct {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	ss_sp     *byte
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	ss_flags  int32
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	pad_cgo_0 [4]byte
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	ss_size   uintptr
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>type mcontext struct {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	gregs       [23]uint64
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	fpregs      *fpstate
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	__reserved1 [8]uint64
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>type ucontext struct {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	uc_flags     uint64
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	uc_link      *ucontext
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	uc_stack     stackt
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	uc_mcontext  mcontext
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	uc_sigmask   usigset
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	__fpregs_mem fpstate
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>type sigcontext struct {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	r8          uint64
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	r9          uint64
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	r10         uint64
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	r11         uint64
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	r12         uint64
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	r13         uint64
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	r14         uint64
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	r15         uint64
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	rdi         uint64
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	rsi         uint64
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	rbp         uint64
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	rbx         uint64
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	rdx         uint64
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	rax         uint64
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	rcx         uint64
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	rsp         uint64
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	rip         uint64
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	eflags      uint64
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	cs          uint16
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	gs          uint16
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	fs          uint16
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	__pad0      uint16
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	err         uint64
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	trapno      uint64
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	oldmask     uint64
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	cr2         uint64
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	fpstate     *fpstate1
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	__reserved1 [8]uint64
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>type sockaddr_un struct {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	family uint16
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	path   [108]byte
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
</pre><p><a href="defs_linux_amd64.go?m=text">View as plain text</a></p>

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
