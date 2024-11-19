<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/udpsock_posix.go - Go Documentation Server</title>

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
<a href="udpsock_posix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<span class="text-muted">udpsock_posix.go</span>
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
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;net/netip&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;syscall&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>func sockaddrToUDP(sa syscall.Sockaddr) Addr {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	switch sa := sa.(type) {
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	case *syscall.SockaddrInet4:
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>		return &amp;UDPAddr{IP: sa.Addr[0:], Port: sa.Port}
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	case *syscall.SockaddrInet6:
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>		return &amp;UDPAddr{IP: sa.Addr[0:], Port: sa.Port, Zone: zoneCache.name(int(sa.ZoneId))}
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	}
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	return nil
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>}
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>func (a *UDPAddr) family() int {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	if a == nil || len(a.IP) &lt;= IPv4len {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		return syscall.AF_INET
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	if a.IP.To4() != nil {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		return syscall.AF_INET
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	return syscall.AF_INET6
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>func (a *UDPAddr) sockaddr(family int) (syscall.Sockaddr, error) {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	if a == nil {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		return nil, nil
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	return ipToSockaddr(family, a.IP, a.Port, a.Zone)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func (a *UDPAddr) toLocal(net string) sockaddr {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	return &amp;UDPAddr{loopbackIP(net), a.Port, a.Zone}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>func (c *UDPConn) readFrom(b []byte, addr *UDPAddr) (int, *UDPAddr, error) {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	var n int
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	var err error
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	switch c.fd.family {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	case syscall.AF_INET:
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		var from syscall.SockaddrInet4
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		n, err = c.fd.readFromInet4(b, &amp;from)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		if err == nil {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			ip := from.Addr <span class="comment">// copy from.Addr; ip escapes, so this line allocates 4 bytes</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>			*addr = UDPAddr{IP: ip[:], Port: from.Port}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	case syscall.AF_INET6:
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		var from syscall.SockaddrInet6
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		n, err = c.fd.readFromInet6(b, &amp;from)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		if err == nil {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			ip := from.Addr <span class="comment">// copy from.Addr; ip escapes, so this line allocates 16 bytes</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			*addr = UDPAddr{IP: ip[:], Port: from.Port, Zone: zoneCache.name(int(from.ZoneId))}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if err != nil {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		<span class="comment">// No sockaddr, so don&#39;t return UDPAddr.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		addr = nil
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	return n, addr, err
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func (c *UDPConn) readFromAddrPort(b []byte) (n int, addr netip.AddrPort, err error) {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	var ip netip.Addr
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	var port int
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	switch c.fd.family {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	case syscall.AF_INET:
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		var from syscall.SockaddrInet4
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		n, err = c.fd.readFromInet4(b, &amp;from)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		if err == nil {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			ip = netip.AddrFrom4(from.Addr)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			port = from.Port
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	case syscall.AF_INET6:
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		var from syscall.SockaddrInet6
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		n, err = c.fd.readFromInet6(b, &amp;from)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		if err == nil {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			ip = netip.AddrFrom16(from.Addr).WithZone(zoneCache.name(int(from.ZoneId)))
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			port = from.Port
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	if err == nil {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		addr = netip.AddrPortFrom(ip, uint16(port))
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	return n, addr, err
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func (c *UDPConn) readMsg(b, oob []byte) (n, oobn, flags int, addr netip.AddrPort, err error) {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	switch c.fd.family {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	case syscall.AF_INET:
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		var sa syscall.SockaddrInet4
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		n, oobn, flags, err = c.fd.readMsgInet4(b, oob, 0, &amp;sa)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		ip := netip.AddrFrom4(sa.Addr)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		addr = netip.AddrPortFrom(ip, uint16(sa.Port))
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	case syscall.AF_INET6:
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		var sa syscall.SockaddrInet6
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		n, oobn, flags, err = c.fd.readMsgInet6(b, oob, 0, &amp;sa)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		ip := netip.AddrFrom16(sa.Addr).WithZone(zoneCache.name(int(sa.ZoneId)))
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		addr = netip.AddrPortFrom(ip, uint16(sa.Port))
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	return
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func (c *UDPConn) writeTo(b []byte, addr *UDPAddr) (int, error) {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	if c.fd.isConnected {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		return 0, ErrWriteToConnected
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	if addr == nil {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		return 0, errMissingAddress
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	switch c.fd.family {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	case syscall.AF_INET:
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		sa, err := ipToSockaddrInet4(addr.IP, addr.Port)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		if err != nil {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			return 0, err
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		return c.fd.writeToInet4(b, &amp;sa)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	case syscall.AF_INET6:
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		sa, err := ipToSockaddrInet6(addr.IP, addr.Port, addr.Zone)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		if err != nil {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			return 0, err
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		return c.fd.writeToInet6(b, &amp;sa)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	default:
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		return 0, &amp;AddrError{Err: &#34;invalid address family&#34;, Addr: addr.IP.String()}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>func (c *UDPConn) writeToAddrPort(b []byte, addr netip.AddrPort) (int, error) {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	if c.fd.isConnected {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		return 0, ErrWriteToConnected
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	if !addr.IsValid() {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		return 0, errMissingAddress
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	switch c.fd.family {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	case syscall.AF_INET:
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		sa, err := addrPortToSockaddrInet4(addr)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		if err != nil {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			return 0, err
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		return c.fd.writeToInet4(b, &amp;sa)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	case syscall.AF_INET6:
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		sa, err := addrPortToSockaddrInet6(addr)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		if err != nil {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			return 0, err
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		return c.fd.writeToInet6(b, &amp;sa)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	default:
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		return 0, &amp;AddrError{Err: &#34;invalid address family&#34;, Addr: addr.Addr().String()}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>func (c *UDPConn) writeMsg(b, oob []byte, addr *UDPAddr) (n, oobn int, err error) {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	if c.fd.isConnected &amp;&amp; addr != nil {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		return 0, 0, ErrWriteToConnected
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	if !c.fd.isConnected &amp;&amp; addr == nil {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		return 0, 0, errMissingAddress
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	sa, err := addr.sockaddr(c.fd.family)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	if err != nil {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		return 0, 0, err
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	return c.fd.writeMsg(b, oob, sa)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>func (c *UDPConn) writeMsgAddrPort(b, oob []byte, addr netip.AddrPort) (n, oobn int, err error) {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	if c.fd.isConnected &amp;&amp; addr.IsValid() {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		return 0, 0, ErrWriteToConnected
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	if !c.fd.isConnected &amp;&amp; !addr.IsValid() {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		return 0, 0, errMissingAddress
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	switch c.fd.family {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	case syscall.AF_INET:
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		sa, err := addrPortToSockaddrInet4(addr)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		if err != nil {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			return 0, 0, err
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		return c.fd.writeMsgInet4(b, oob, &amp;sa)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	case syscall.AF_INET6:
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		sa, err := addrPortToSockaddrInet6(addr)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		if err != nil {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			return 0, 0, err
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		return c.fd.writeMsgInet6(b, oob, &amp;sa)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	default:
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		return 0, 0, &amp;AddrError{Err: &#34;invalid address family&#34;, Addr: addr.Addr().String()}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>func (sd *sysDialer) dialUDP(ctx context.Context, laddr, raddr *UDPAddr) (*UDPConn, error) {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	ctrlCtxFn := sd.Dialer.ControlContext
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	if ctrlCtxFn == nil &amp;&amp; sd.Dialer.Control != nil {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		ctrlCtxFn = func(cxt context.Context, network, address string, c syscall.RawConn) error {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			return sd.Dialer.Control(network, address, c)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	fd, err := internetSocket(ctx, sd.network, laddr, raddr, syscall.SOCK_DGRAM, 0, &#34;dial&#34;, ctrlCtxFn)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	if err != nil {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		return nil, err
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	return newUDPConn(fd), nil
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>func (sl *sysListener) listenUDP(ctx context.Context, laddr *UDPAddr) (*UDPConn, error) {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	var ctrlCtxFn func(cxt context.Context, network, address string, c syscall.RawConn) error
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	if sl.ListenConfig.Control != nil {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		ctrlCtxFn = func(cxt context.Context, network, address string, c syscall.RawConn) error {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			return sl.ListenConfig.Control(network, address, c)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	fd, err := internetSocket(ctx, sl.network, laddr, nil, syscall.SOCK_DGRAM, 0, &#34;listen&#34;, ctrlCtxFn)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	if err != nil {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		return nil, err
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	return newUDPConn(fd), nil
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>func (sl *sysListener) listenMulticastUDP(ctx context.Context, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error) {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	var ctrlCtxFn func(cxt context.Context, network, address string, c syscall.RawConn) error
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	if sl.ListenConfig.Control != nil {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		ctrlCtxFn = func(cxt context.Context, network, address string, c syscall.RawConn) error {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			return sl.ListenConfig.Control(network, address, c)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	fd, err := internetSocket(ctx, sl.network, gaddr, nil, syscall.SOCK_DGRAM, 0, &#34;listen&#34;, ctrlCtxFn)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	if err != nil {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		return nil, err
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	c := newUDPConn(fd)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	if ip4 := gaddr.IP.To4(); ip4 != nil {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		if err := listenIPv4MulticastUDP(c, ifi, ip4); err != nil {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			c.Close()
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			return nil, err
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	} else {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		if err := listenIPv6MulticastUDP(c, ifi, gaddr.IP); err != nil {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			c.Close()
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			return nil, err
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	return c, nil
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>func listenIPv4MulticastUDP(c *UDPConn, ifi *Interface, ip IP) error {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	if ifi != nil {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		if err := setIPv4MulticastInterface(c.fd, ifi); err != nil {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			return err
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	if err := setIPv4MulticastLoopback(c.fd, false); err != nil {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		return err
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	if err := joinIPv4Group(c.fd, ifi, ip); err != nil {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		return err
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	return nil
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>func listenIPv6MulticastUDP(c *UDPConn, ifi *Interface, ip IP) error {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	if ifi != nil {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		if err := setIPv6MulticastInterface(c.fd, ifi); err != nil {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			return err
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	if err := setIPv6MulticastLoopback(c.fd, false); err != nil {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		return err
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	if err := joinIPv6Group(c.fd, ifi, ip); err != nil {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		return err
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	return nil
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
</pre><p><a href="udpsock_posix.go?m=text">View as plain text</a></p>

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
