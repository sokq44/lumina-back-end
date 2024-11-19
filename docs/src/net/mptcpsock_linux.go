<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/mptcpsock_linux.go - Go Documentation Server</title>

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
<a href="mptcpsock_linux.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<span class="text-muted">mptcpsock_linux.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net">net</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2023 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package net
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;context&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/poll&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/syscall/unix&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;syscall&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>var (
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	mptcpOnce      sync.Once
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	mptcpAvailable bool
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	hasSOLMPTCP    bool
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// These constants aren&#39;t in the syscall package, which is frozen</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>const (
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	_IPPROTO_MPTCP = 0x106
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	_SOL_MPTCP     = 0x11c
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	_MPTCP_INFO    = 0x1
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>func supportsMultipathTCP() bool {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	mptcpOnce.Do(initMPTCPavailable)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	return mptcpAvailable
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// Check that MPTCP is supported by attempting to create an MPTCP socket and by</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// looking at the returned error if any.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>func initMPTCPavailable() {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	s, err := sysSocket(syscall.AF_INET, syscall.SOCK_STREAM, _IPPROTO_MPTCP)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	switch {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	case errors.Is(err, syscall.EPROTONOSUPPORT): <span class="comment">// Not supported: &gt;= v5.6</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	case errors.Is(err, syscall.EINVAL): <span class="comment">// Not supported: &lt; v5.6</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	case err == nil: <span class="comment">// Supported and no error</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		poll.CloseFunc(s)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		fallthrough
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	default:
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		<span class="comment">// another error: MPTCP was not available but it might be later</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		mptcpAvailable = true
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	major, minor := unix.KernelVersion()
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// SOL_MPTCP only supported from kernel 5.16</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	hasSOLMPTCP = major &gt; 5 || (major == 5 &amp;&amp; minor &gt;= 16)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func (sd *sysDialer) dialMPTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error) {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	if supportsMultipathTCP() {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		if conn, err := sd.doDialTCPProto(ctx, laddr, raddr, _IPPROTO_MPTCP); err == nil {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			return conn, nil
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// Fallback to dialTCP if Multipath TCP isn&#39;t supported on this operating</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// system. But also fallback in case of any error with MPTCP.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// Possible MPTCP specific error: ENOPROTOOPT (sysctl net.mptcp.enabled=0)</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// But just in case MPTCP is blocked differently (SELinux, etc.), just</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// retry with &#34;plain&#34; TCP.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	return sd.dialTCP(ctx, laddr, raddr)
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>func (sl *sysListener) listenMPTCP(ctx context.Context, laddr *TCPAddr) (*TCPListener, error) {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if supportsMultipathTCP() {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		if dial, err := sl.listenTCPProto(ctx, laddr, _IPPROTO_MPTCP); err == nil {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			return dial, nil
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// Fallback to listenTCP if Multipath TCP isn&#39;t supported on this operating</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// system. But also fallback in case of any error with MPTCP.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// Possible MPTCP specific error: ENOPROTOOPT (sysctl net.mptcp.enabled=0)</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// But just in case MPTCP is blocked differently (SELinux, etc.), just</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// retry with &#34;plain&#34; TCP.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	return sl.listenTCP(ctx, laddr)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// hasFallenBack reports whether the MPTCP connection has fallen back to &#34;plain&#34;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// TCP.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// A connection can fallback to TCP for different reasons, e.g. the other peer</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// doesn&#39;t support it, a middle box &#34;accidentally&#34; drops the option, etc.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// If the MPTCP protocol has not been requested when creating the socket, this</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// method will return true: MPTCP is not being used.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// Kernel &gt;= 5.16 returns EOPNOTSUPP/ENOPROTOOPT in case of fallback.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// Older kernels will always return them even if MPTCP is used: not usable.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func hasFallenBack(fd *netFD) bool {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	_, err := fd.pfd.GetsockoptInt(_SOL_MPTCP, _MPTCP_INFO)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// 2 expected errors in case of fallback depending on the address family</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">//   - AF_INET:  EOPNOTSUPP</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">//   - AF_INET6: ENOPROTOOPT</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	return err == syscall.EOPNOTSUPP || err == syscall.ENOPROTOOPT
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// isUsingMPTCPProto reports whether the socket protocol is MPTCP.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// Compared to hasFallenBack method, here only the socket protocol being used is</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// checked: it can be MPTCP but it doesn&#39;t mean MPTCP is used on the wire, maybe</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// a fallback to TCP has been done.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>func isUsingMPTCPProto(fd *netFD) bool {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	proto, _ := fd.pfd.GetsockoptInt(syscall.SOL_SOCKET, syscall.SO_PROTOCOL)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	return proto == _IPPROTO_MPTCP
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// isUsingMultipathTCP reports whether MPTCP is still being used.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// Please look at the description of hasFallenBack (kernel &gt;=5.16) and</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// isUsingMPTCPProto methods for more details about what is being checked here.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>func isUsingMultipathTCP(fd *netFD) bool {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	if hasSOLMPTCP {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		return !hasFallenBack(fd)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	return isUsingMPTCPProto(fd)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
</pre><p><a href="mptcpsock_linux.go?m=text">View as plain text</a></p>

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
