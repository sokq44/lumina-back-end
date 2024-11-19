<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/os/signal/doc.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../index.html">GoDoc</a></div>
<a href="doc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/os">os</a>/<a href="http://localhost:8080/src/os/signal">signal</a>/<span class="text-muted">doc.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/os/signal">os/signal</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2015 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">/*
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>Package signal implements access to incoming signals.
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>Signals are primarily used on Unix-like systems. For the use of this
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package on Windows and Plan 9, see below.
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span># Types of signals
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>The signals SIGKILL and SIGSTOP may not be caught by a program, and
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>therefore cannot be affected by this package.
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>Synchronous signals are signals triggered by errors in program
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>execution: SIGBUS, SIGFPE, and SIGSEGV. These are only considered
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>synchronous when caused by program execution, not when sent using
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>[os.Process.Kill] or the kill program or some similar mechanism. In
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>general, except as discussed below, Go programs will convert a
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>synchronous signal into a run-time panic.
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>The remaining signals are asynchronous signals. They are not
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>triggered by program errors, but are instead sent from the kernel or
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>from some other program.
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>Of the asynchronous signals, the SIGHUP signal is sent when a program
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>loses its controlling terminal. The SIGINT signal is sent when the
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>user at the controlling terminal presses the interrupt character,
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>which by default is ^C (Control-C). The SIGQUIT signal is sent when
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>the user at the controlling terminal presses the quit character, which
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>by default is ^\ (Control-Backslash). In general you can cause a
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>program to simply exit by pressing ^C, and you can cause it to exit
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>with a stack dump by pressing ^\.
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span># Default behavior of signals in Go programs
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>By default, a synchronous signal is converted into a run-time panic. A
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>SIGHUP, SIGINT, or SIGTERM signal causes the program to exit. A
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>SIGQUIT, SIGILL, SIGTRAP, SIGABRT, SIGSTKFLT, SIGEMT, or SIGSYS signal
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>causes the program to exit with a stack dump. A SIGTSTP, SIGTTIN, or
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>SIGTTOU signal gets the system default behavior (these signals are
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>used by the shell for job control). The SIGPROF signal is handled
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>directly by the Go runtime to implement runtime.CPUProfile. Other
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>signals will be caught but no action will be taken.
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>If the Go program is started with either SIGHUP or SIGINT ignored
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>(signal handler set to SIG_IGN), they will remain ignored.
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>If the Go program is started with a non-empty signal mask, that will
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>generally be honored. However, some signals are explicitly unblocked:
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>the synchronous signals, SIGILL, SIGTRAP, SIGSTKFLT, SIGCHLD, SIGPROF,
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>and, on Linux, signals 32 (SIGCANCEL) and 33 (SIGSETXID)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>(SIGCANCEL and SIGSETXID are used internally by glibc). Subprocesses
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>started by [os.Exec], or by [os/exec], will inherit the
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>modified signal mask.
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span># Changing the behavior of signals in Go programs
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>The functions in this package allow a program to change the way Go
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>programs handle signals.
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>Notify disables the default behavior for a given set of asynchronous
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>signals and instead delivers them over one or more registered
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>channels. Specifically, it applies to the signals SIGHUP, SIGINT,
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>SIGQUIT, SIGABRT, and SIGTERM. It also applies to the job control
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>signals SIGTSTP, SIGTTIN, and SIGTTOU, in which case the system
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>default behavior does not occur. It also applies to some signals that
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>otherwise cause no action: SIGUSR1, SIGUSR2, SIGPIPE, SIGALRM,
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>SIGCHLD, SIGCONT, SIGURG, SIGXCPU, SIGXFSZ, SIGVTALRM, SIGWINCH,
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>SIGIO, SIGPWR, SIGSYS, SIGINFO, SIGTHR, SIGWAITING, SIGLWP, SIGFREEZE,
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>SIGTHAW, SIGLOST, SIGXRES, SIGJVM1, SIGJVM2, and any real time signals
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>used on the system. Note that not all of these signals are available
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>on all systems.
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>If the program was started with SIGHUP or SIGINT ignored, and Notify
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>is called for either signal, a signal handler will be installed for
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>that signal and it will no longer be ignored. If, later, Reset or
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>Ignore is called for that signal, or Stop is called on all channels
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>passed to Notify for that signal, the signal will once again be
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>ignored. Reset will restore the system default behavior for the
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>signal, while Ignore will cause the system to ignore the signal
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>entirely.
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>If the program is started with a non-empty signal mask, some signals
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>will be explicitly unblocked as described above. If Notify is called
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>for a blocked signal, it will be unblocked. If, later, Reset is
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>called for that signal, or Stop is called on all channels passed to
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>Notify for that signal, the signal will once again be blocked.
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span># SIGPIPE
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>When a Go program writes to a broken pipe, the kernel will raise a
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>SIGPIPE signal.
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>If the program has not called Notify to receive SIGPIPE signals, then
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>the behavior depends on the file descriptor number. A write to a
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>broken pipe on file descriptors 1 or 2 (standard output or standard
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>error) will cause the program to exit with a SIGPIPE signal. A write
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>to a broken pipe on some other file descriptor will take no action on
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>the SIGPIPE signal, and the write will fail with an EPIPE error.
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>If the program has called Notify to receive SIGPIPE signals, the file
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>descriptor number does not matter. The SIGPIPE signal will be
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>delivered to the Notify channel, and the write will fail with an EPIPE
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>error.
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>This means that, by default, command line programs will behave like
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>typical Unix command line programs, while other programs will not
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>crash with SIGPIPE when writing to a closed network connection.
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span># Go programs that use cgo or SWIG
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>In a Go program that includes non-Go code, typically C/C++ code
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>accessed using cgo or SWIG, Go&#39;s startup code normally runs first. It
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>configures the signal handlers as expected by the Go runtime, before
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>the non-Go startup code runs. If the non-Go startup code wishes to
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>install its own signal handlers, it must take certain steps to keep Go
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>working well. This section documents those steps and the overall
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>effect changes to signal handler settings by the non-Go code can have
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>on Go programs. In rare cases, the non-Go code may run before the Go
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>code, in which case the next section also applies.
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>If the non-Go code called by the Go program does not change any signal
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>handlers or masks, then the behavior is the same as for a pure Go
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>program.
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>If the non-Go code installs any signal handlers, it must use the
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>SA_ONSTACK flag with sigaction. Failing to do so is likely to cause
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>the program to crash if the signal is received. Go programs routinely
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>run with a limited stack, and therefore set up an alternate signal
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>stack.
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>If the non-Go code installs a signal handler for any of the
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>synchronous signals (SIGBUS, SIGFPE, SIGSEGV), then it should record
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>the existing Go signal handler. If those signals occur while
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>executing Go code, it should invoke the Go signal handler (whether the
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>signal occurs while executing Go code can be determined by looking at
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>the PC passed to the signal handler). Otherwise some Go run-time
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>panics will not occur as expected.
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>If the non-Go code installs a signal handler for any of the
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>asynchronous signals, it may invoke the Go signal handler or not as it
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>chooses. Naturally, if it does not invoke the Go signal handler, the
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>Go behavior described above will not occur. This can be an issue with
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>the SIGPROF signal in particular.
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>The non-Go code should not change the signal mask on any threads
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>created by the Go runtime. If the non-Go code starts new threads of
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>its own, it may set the signal mask as it pleases.
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>If the non-Go code starts a new thread, changes the signal mask, and
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>then invokes a Go function in that thread, the Go runtime will
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>automatically unblock certain signals: the synchronous signals,
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>SIGILL, SIGTRAP, SIGSTKFLT, SIGCHLD, SIGPROF, SIGCANCEL, and
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>SIGSETXID. When the Go function returns, the non-Go signal mask will
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>be restored.
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>If the Go signal handler is invoked on a non-Go thread not running Go
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>code, the handler generally forwards the signal to the non-Go code, as
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>follows. If the signal is SIGPROF, the Go handler does
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>nothing. Otherwise, the Go handler removes itself, unblocks the
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>signal, and raises it again, to invoke any non-Go handler or default
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>system handler. If the program does not exit, the Go handler then
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>reinstalls itself and continues execution of the program.
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>If a SIGPIPE signal is received, the Go program will invoke the
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>special handling described above if the SIGPIPE is received on a Go
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>thread.  If the SIGPIPE is received on a non-Go thread the signal will
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>be forwarded to the non-Go handler, if any; if there is none the
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>default system handler will cause the program to terminate.
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span># Non-Go programs that call Go code
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>When Go code is built with options like -buildmode=c-shared, it will
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>be run as part of an existing non-Go program. The non-Go code may
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>have already installed signal handlers when the Go code starts (that
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>may also happen in unusual cases when using cgo or SWIG; in that case,
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>the discussion here applies).  For -buildmode=c-archive the Go runtime
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>will initialize signals at global constructor time.  For
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>-buildmode=c-shared the Go runtime will initialize signals when the
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>shared library is loaded.
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>If the Go runtime sees an existing signal handler for the SIGCANCEL or
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>SIGSETXID signals (which are used only on Linux), it will turn on
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>the SA_ONSTACK flag and otherwise keep the signal handler.
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>For the synchronous signals and SIGPIPE, the Go runtime will install a
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>signal handler. It will save any existing signal handler. If a
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>synchronous signal arrives while executing non-Go code, the Go runtime
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>will invoke the existing signal handler instead of the Go signal
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>handler.
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>Go code built with -buildmode=c-archive or -buildmode=c-shared will
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>not install any other signal handlers by default. If there is an
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>existing signal handler, the Go runtime will turn on the SA_ONSTACK
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>flag and otherwise keep the signal handler. If Notify is called for an
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>asynchronous signal, a Go signal handler will be installed for that
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>signal. If, later, Reset is called for that signal, the original
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>handling for that signal will be reinstalled, restoring the non-Go
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>signal handler if any.
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>Go code built without -buildmode=c-archive or -buildmode=c-shared will
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>install a signal handler for the asynchronous signals listed above,
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>and save any existing signal handler. If a signal is delivered to a
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>non-Go thread, it will act as described above, except that if there is
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>an existing non-Go signal handler, that handler will be installed
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>before raising the signal.
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span># Windows
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>On Windows a ^C (Control-C) or ^BREAK (Control-Break) normally cause
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>the program to exit. If Notify is called for [os.Interrupt], ^C or ^BREAK
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>will cause [os.Interrupt] to be sent on the channel, and the program will
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>not exit. If Reset is called, or Stop is called on all channels passed
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>to Notify, then the default behavior will be restored.
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>Additionally, if Notify is called, and Windows sends CTRL_CLOSE_EVENT,
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>CTRL_LOGOFF_EVENT or CTRL_SHUTDOWN_EVENT to the process, Notify will
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>return syscall.SIGTERM. Unlike Control-C and Control-Break, Notify does
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>not change process behavior when either CTRL_CLOSE_EVENT,
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>CTRL_LOGOFF_EVENT or CTRL_SHUTDOWN_EVENT is received - the process will
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>still get terminated unless it exits. But receiving syscall.SIGTERM will
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>give the process an opportunity to clean up before termination.
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span># Plan 9
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>On Plan 9, signals have type syscall.Note, which is a string. Calling
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>Notify with a syscall.Note will cause that value to be sent on the
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>channel when that string is posted as a note.
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>*/</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>package signal
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>
</pre><p><a href="doc.go?m=text">View as plain text</a></p>

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
