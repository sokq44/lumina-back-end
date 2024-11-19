<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/os/exec_posix.go - Go Documentation Server</title>

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
<a href="exec_posix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/os">os</a>/<span class="text-muted">exec_posix.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/os">os</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build unix || (js &amp;&amp; wasm) || wasip1 || windows</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package os
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/itoa&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/syscall/execenv&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;runtime&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;syscall&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// The only signal values guaranteed to be present in the os package on all</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// systems are os.Interrupt (send the process an interrupt) and os.Kill (force</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// the process to exit). On Windows, sending os.Interrupt to a process with</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// os.Process.Signal is not implemented; it will return an error instead of</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// sending a signal.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>var (
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	Interrupt Signal = syscall.SIGINT
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	Kill      Signal = syscall.SIGKILL
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>func startProcess(name string, argv []string, attr *ProcAttr) (p *Process, err error) {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// If there is no SysProcAttr (ie. no Chroot or changed</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// UID/GID), double-check existence of the directory we want</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// to chdir into. We can make the error clearer this way.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	if attr != nil &amp;&amp; attr.Sys == nil &amp;&amp; attr.Dir != &#34;&#34; {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		if _, err := Stat(attr.Dir); err != nil {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>			pe := err.(*PathError)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>			pe.Op = &#34;chdir&#34;
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>			return nil, pe
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	sysattr := &amp;syscall.ProcAttr{
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		Dir: attr.Dir,
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		Env: attr.Env,
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		Sys: attr.Sys,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	}
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	if sysattr.Env == nil {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		sysattr.Env, err = execenv.Default(sysattr.Sys)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		if err != nil {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			return nil, err
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	sysattr.Files = make([]uintptr, 0, len(attr.Files))
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	for _, f := range attr.Files {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		sysattr.Files = append(sysattr.Files, f.Fd())
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	pid, h, e := syscall.StartProcess(name, argv, sysattr)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// Make sure we don&#39;t run the finalizers of attr.Files.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	runtime.KeepAlive(attr)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if e != nil {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		return nil, &amp;PathError{Op: &#34;fork/exec&#34;, Path: name, Err: e}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	return newProcess(pid, h), nil
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>func (p *Process) kill() error {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	return p.Signal(Kill)
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// ProcessState stores information about a process, as reported by Wait.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>type ProcessState struct {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	pid    int                <span class="comment">// The process&#39;s id.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	status syscall.WaitStatus <span class="comment">// System-dependent status info.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	rusage *syscall.Rusage
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// Pid returns the process id of the exited process.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>func (p *ProcessState) Pid() int {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	return p.pid
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>func (p *ProcessState) exited() bool {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	return p.status.Exited()
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>func (p *ProcessState) success() bool {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	return p.status.ExitStatus() == 0
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>func (p *ProcessState) sys() any {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	return p.status
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>func (p *ProcessState) sysUsage() any {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	return p.rusage
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>func (p *ProcessState) String() string {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	if p == nil {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		return &#34;&lt;nil&gt;&#34;
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	status := p.Sys().(syscall.WaitStatus)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	res := &#34;&#34;
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	switch {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	case status.Exited():
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		code := status.ExitStatus()
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		if runtime.GOOS == &#34;windows&#34; &amp;&amp; uint(code) &gt;= 1&lt;&lt;16 { <span class="comment">// windows uses large hex numbers</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			res = &#34;exit status &#34; + itoa.Uitox(uint(code))
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		} else { <span class="comment">// unix systems use small decimal integers</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			res = &#34;exit status &#34; + itoa.Itoa(code) <span class="comment">// unix</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	case status.Signaled():
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		res = &#34;signal: &#34; + status.Signal().String()
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	case status.Stopped():
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		res = &#34;stop signal: &#34; + status.StopSignal().String()
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		if status.StopSignal() == syscall.SIGTRAP &amp;&amp; status.TrapCause() != 0 {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			res += &#34; (trap &#34; + itoa.Itoa(status.TrapCause()) + &#34;)&#34;
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	case status.Continued():
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		res = &#34;continued&#34;
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	if status.CoreDump() {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		res += &#34; (core dumped)&#34;
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	return res
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// ExitCode returns the exit code of the exited process, or -1</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// if the process hasn&#39;t exited or was terminated by a signal.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>func (p *ProcessState) ExitCode() int {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// return -1 if the process hasn&#39;t started.</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	if p == nil {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		return -1
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	return p.status.ExitStatus()
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
</pre><p><a href="exec_posix.go?m=text">View as plain text</a></p>

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
