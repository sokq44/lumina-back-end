<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/tcpsock_posix.go - Go Documentation Server</title>

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
<a href="tcpsock_posix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<span class="text-muted">tcpsock_posix.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net">net</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build unix || js || wasip1 || windows</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package net
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;context&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;os&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;syscall&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>func sockaddrToTCP(sa syscall.Sockaddr) Addr {
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	switch sa := sa.(type) {
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	case *syscall.SockaddrInet4:
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>		return &amp;TCPAddr{IP: sa.Addr[0:], Port: sa.Port}
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	case *syscall.SockaddrInet6:
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>		return &amp;TCPAddr{IP: sa.Addr[0:], Port: sa.Port, Zone: zoneCache.name(int(sa.ZoneId))}
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	return nil
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>func (a *TCPAddr) family() int {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	if a == nil || len(a.IP) &lt;= IPv4len {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		return syscall.AF_INET
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	if a.IP.To4() != nil {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		return syscall.AF_INET
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	return syscall.AF_INET6
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>func (a *TCPAddr) sockaddr(family int) (syscall.Sockaddr, error) {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	if a == nil {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		return nil, nil
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	return ipToSockaddr(family, a.IP, a.Port, a.Zone)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func (a *TCPAddr) toLocal(net string) sockaddr {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	return &amp;TCPAddr{loopbackIP(net), a.Port, a.Zone}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>func (c *TCPConn) readFrom(r io.Reader) (int64, error) {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	if n, err, handled := spliceFrom(c.fd, r); handled {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		return n, err
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	if n, err, handled := sendFile(c.fd, r); handled {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		return n, err
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	return genericReadFrom(c, r)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>func (c *TCPConn) writeTo(w io.Writer) (int64, error) {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	if n, err, handled := spliceTo(w, c.fd); handled {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		return n, err
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	return genericWriteTo(c, w)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>func (sd *sysDialer) dialTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error) {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if h := sd.testHookDialTCP; h != nil {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		return h(ctx, sd.network, laddr, raddr)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	if h := testHookDialTCP; h != nil {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		return h(ctx, sd.network, laddr, raddr)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	return sd.doDialTCP(ctx, laddr, raddr)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>func (sd *sysDialer) doDialTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error) {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	return sd.doDialTCPProto(ctx, laddr, raddr, 0)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>func (sd *sysDialer) doDialTCPProto(ctx context.Context, laddr, raddr *TCPAddr, proto int) (*TCPConn, error) {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	ctrlCtxFn := sd.Dialer.ControlContext
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if ctrlCtxFn == nil &amp;&amp; sd.Dialer.Control != nil {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		ctrlCtxFn = func(cxt context.Context, network, address string, c syscall.RawConn) error {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			return sd.Dialer.Control(network, address, c)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	fd, err := internetSocket(ctx, sd.network, laddr, raddr, syscall.SOCK_STREAM, proto, &#34;dial&#34;, ctrlCtxFn)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// TCP has a rarely used mechanism called a &#39;simultaneous connection&#39; in</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// which Dial(&#34;tcp&#34;, addr1, addr2) run on the machine at addr1 can</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// connect to a simultaneous Dial(&#34;tcp&#34;, addr2, addr1) run on the machine</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// at addr2, without either machine executing Listen. If laddr == nil,</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// it means we want the kernel to pick an appropriate originating local</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// address. Some Linux kernels cycle blindly through a fixed range of</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// local ports, regardless of destination port. If a kernel happens to</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// pick local port 50001 as the source for a Dial(&#34;tcp&#34;, &#34;&#34;, &#34;localhost:50001&#34;),</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// then the Dial will succeed, having simultaneously connected to itself.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// This can only happen when we are letting the kernel pick a port (laddr == nil)</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// and when there is no listener for the destination address.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s hard to argue this is anything other than a kernel bug. If we</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// see this happen, rather than expose the buggy effect to users, we</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	<span class="comment">// close the fd and try again. If it happens twice more, we relent and</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// use the result. See also:</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">//	https://golang.org/issue/2690</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">//	https://stackoverflow.com/questions/4949858/</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// The opposite can also happen: if we ask the kernel to pick an appropriate</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	<span class="comment">// originating local address, sometimes it picks one that is already in use.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// So if the error is EADDRNOTAVAIL, we have to try again too, just for</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">// a different reason.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// The kernel socket code is no doubt enjoying watching us squirm.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	for i := 0; i &lt; 2 &amp;&amp; (laddr == nil || laddr.Port == 0) &amp;&amp; (selfConnect(fd, err) || spuriousENOTAVAIL(err)); i++ {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		if err == nil {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			fd.Close()
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		fd, err = internetSocket(ctx, sd.network, laddr, raddr, syscall.SOCK_STREAM, proto, &#34;dial&#34;, ctrlCtxFn)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	if err != nil {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		return nil, err
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	return newTCPConn(fd, sd.Dialer.KeepAlive, testHookSetKeepAlive), nil
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func selfConnect(fd *netFD, err error) bool {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// If the connect failed, we clearly didn&#39;t connect to ourselves.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	if err != nil {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		return false
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// The socket constructor can return an fd with raddr nil under certain</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// unknown conditions. The errors in the calls there to Getpeername</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// are discarded, but we can&#39;t catch the problem there because those</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// calls are sometimes legally erroneous with a &#34;socket not connected&#34;.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// Since this code (selfConnect) is already trying to work around</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// a problem, we make sure if this happens we recognize trouble and</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// ask the DialTCP routine to try again.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// TODO: try to understand what&#39;s really going on.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	if fd.laddr == nil || fd.raddr == nil {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		return true
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	l := fd.laddr.(*TCPAddr)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	r := fd.raddr.(*TCPAddr)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	return l.Port == r.Port &amp;&amp; l.IP.Equal(r.IP)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>func spuriousENOTAVAIL(err error) bool {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	if op, ok := err.(*OpError); ok {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		err = op.Err
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	if sys, ok := err.(*os.SyscallError); ok {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		err = sys.Err
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	return err == syscall.EADDRNOTAVAIL
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>func (ln *TCPListener) ok() bool { return ln != nil &amp;&amp; ln.fd != nil }
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>func (ln *TCPListener) accept() (*TCPConn, error) {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	fd, err := ln.fd.accept()
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	if err != nil {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		return nil, err
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	return newTCPConn(fd, ln.lc.KeepAlive, nil), nil
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>func (ln *TCPListener) close() error {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	return ln.fd.Close()
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>func (ln *TCPListener) file() (*os.File, error) {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	f, err := ln.fd.dup()
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	if err != nil {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		return nil, err
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	return f, nil
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>func (sl *sysListener) listenTCP(ctx context.Context, laddr *TCPAddr) (*TCPListener, error) {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	return sl.listenTCPProto(ctx, laddr, 0)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>func (sl *sysListener) listenTCPProto(ctx context.Context, laddr *TCPAddr, proto int) (*TCPListener, error) {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	var ctrlCtxFn func(cxt context.Context, network, address string, c syscall.RawConn) error
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	if sl.ListenConfig.Control != nil {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		ctrlCtxFn = func(cxt context.Context, network, address string, c syscall.RawConn) error {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			return sl.ListenConfig.Control(network, address, c)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	fd, err := internetSocket(ctx, sl.network, laddr, nil, syscall.SOCK_STREAM, proto, &#34;listen&#34;, ctrlCtxFn)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	if err != nil {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		return nil, err
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	return &amp;TCPListener{fd: fd, lc: sl.ListenConfig}, nil
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
</pre><p><a href="tcpsock_posix.go?m=text">View as plain text</a></p>

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
