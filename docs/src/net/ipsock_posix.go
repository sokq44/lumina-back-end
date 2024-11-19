<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/ipsock_posix.go - Go Documentation Server</title>

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
<a href="ipsock_posix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<span class="text-muted">ipsock_posix.go</span>
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
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/poll&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;net/netip&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;runtime&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;syscall&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// probe probes IPv4, IPv6 and IPv4-mapped IPv6 communication</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// capabilities which are controlled by the IPV6_V6ONLY socket option</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// and kernel configuration.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// Should we try to use the IPv4 socket interface if we&#39;re only</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// dealing with IPv4 sockets? As long as the host system understands</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// IPv4-mapped IPv6, it&#39;s okay to pass IPv4-mapped IPv6 addresses to</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// the IPv6 interface. That simplifies our code and is most</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// general. Unfortunately, we need to run on kernels built without</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// IPv6 support too. So probe the kernel to figure it out.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>func (p *ipStackCapabilities) probe() {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	switch runtime.GOOS {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	case &#34;js&#34;, &#34;wasip1&#34;:
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		<span class="comment">// Both ipv4 and ipv6 are faked; see net_fake.go.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		p.ipv4Enabled = true
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		p.ipv6Enabled = true
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		p.ipv4MappedIPv6Enabled = true
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		return
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	s, err := sysSocket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	switch err {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	case syscall.EAFNOSUPPORT, syscall.EPROTONOSUPPORT:
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	case nil:
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		poll.CloseFunc(s)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		p.ipv4Enabled = true
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	var probes = []struct {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		laddr TCPAddr
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		value int
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	}{
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		<span class="comment">// IPv6 communication capability</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		{laddr: TCPAddr{IP: ParseIP(&#34;::1&#34;)}, value: 1},
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		<span class="comment">// IPv4-mapped IPv6 address communication capability</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		{laddr: TCPAddr{IP: IPv4(127, 0, 0, 1)}, value: 0},
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	switch runtime.GOOS {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	case &#34;dragonfly&#34;, &#34;openbsd&#34;:
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		<span class="comment">// The latest DragonFly BSD and OpenBSD kernels don&#39;t</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		<span class="comment">// support IPV6_V6ONLY=0. They always return an error</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		<span class="comment">// and we don&#39;t need to probe the capability.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		probes = probes[:1]
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	for i := range probes {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		s, err := sysSocket(syscall.AF_INET6, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		if err != nil {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			continue
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		defer poll.CloseFunc(s)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, probes[i].value)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		sa, err := probes[i].laddr.sockaddr(syscall.AF_INET6)
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		if err != nil {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			continue
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		if err := syscall.Bind(s, sa); err != nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			continue
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		if i == 0 {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			p.ipv6Enabled = true
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		} else {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			p.ipv4MappedIPv6Enabled = true
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// favoriteAddrFamily returns the appropriate address family for the</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// given network, laddr, raddr and mode.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// If mode indicates &#34;listen&#34; and laddr is a wildcard, we assume that</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// the user wants to make a passive-open connection with a wildcard</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// address family, both AF_INET and AF_INET6, and a wildcard address</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// like the following:</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">//   - A listen for a wildcard communication domain, &#34;tcp&#34; or</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//     &#34;udp&#34;, with a wildcard address: If the platform supports</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">//     both IPv6 and IPv4-mapped IPv6 communication capabilities,</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">//     or does not support IPv4, we use a dual stack, AF_INET6 and</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">//     IPV6_V6ONLY=0, wildcard address listen. The dual stack</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">//     wildcard address listen may fall back to an IPv6-only,</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">//     AF_INET6 and IPV6_V6ONLY=1, wildcard address listen.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">//     Otherwise we prefer an IPv4-only, AF_INET, wildcard address</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//     listen.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">//   - A listen for a wildcard communication domain, &#34;tcp&#34; or</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">//     &#34;udp&#34;, with an IPv4 wildcard address: same as above.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">//   - A listen for a wildcard communication domain, &#34;tcp&#34; or</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">//     &#34;udp&#34;, with an IPv6 wildcard address: same as above.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//   - A listen for an IPv4 communication domain, &#34;tcp4&#34; or &#34;udp4&#34;,</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">//     with an IPv4 wildcard address: We use an IPv4-only, AF_INET,</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">//     wildcard address listen.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">//   - A listen for an IPv6 communication domain, &#34;tcp6&#34; or &#34;udp6&#34;,</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//     with an IPv6 wildcard address: We use an IPv6-only, AF_INET6</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//     and IPV6_V6ONLY=1, wildcard address listen.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">// Otherwise guess: If the addresses are IPv4 then returns AF_INET,</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">// or else returns AF_INET6. It also returns a boolean value what</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// designates IPV6_V6ONLY option.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// Note that the latest DragonFly BSD and OpenBSD kernels allow</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// neither &#34;net.inet6.ip6.v6only=1&#34; change nor IPPROTO_IPV6 level</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// IPV6_V6ONLY socket option setting.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>func favoriteAddrFamily(network string, laddr, raddr sockaddr, mode string) (family int, ipv6only bool) {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	switch network[len(network)-1] {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	case &#39;4&#39;:
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		return syscall.AF_INET, false
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	case &#39;6&#39;:
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		return syscall.AF_INET6, true
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if mode == &#34;listen&#34; &amp;&amp; (laddr == nil || laddr.isWildcard()) {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		if supportsIPv4map() || !supportsIPv4() {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			return syscall.AF_INET6, false
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		if laddr == nil {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			return syscall.AF_INET, false
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		return laddr.family(), false
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	if (laddr == nil || laddr.family() == syscall.AF_INET) &amp;&amp;
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		(raddr == nil || raddr.family() == syscall.AF_INET) {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		return syscall.AF_INET, false
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	return syscall.AF_INET6, false
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>func internetSocket(ctx context.Context, net string, laddr, raddr sockaddr, sotype, proto int, mode string, ctrlCtxFn func(context.Context, string, string, syscall.RawConn) error) (fd *netFD, err error) {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	switch runtime.GOOS {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	case &#34;aix&#34;, &#34;windows&#34;, &#34;openbsd&#34;, &#34;js&#34;, &#34;wasip1&#34;:
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		if mode == &#34;dial&#34; &amp;&amp; raddr.isWildcard() {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			raddr = raddr.toLocal(net)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	family, ipv6only := favoriteAddrFamily(net, laddr, raddr, mode)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	return socket(ctx, net, family, sotype, proto, ipv6only, laddr, raddr, ctrlCtxFn)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>func ipToSockaddrInet4(ip IP, port int) (syscall.SockaddrInet4, error) {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	if len(ip) == 0 {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		ip = IPv4zero
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	ip4 := ip.To4()
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	if ip4 == nil {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		return syscall.SockaddrInet4{}, &amp;AddrError{Err: &#34;non-IPv4 address&#34;, Addr: ip.String()}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	sa := syscall.SockaddrInet4{Port: port}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	copy(sa.Addr[:], ip4)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	return sa, nil
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>func ipToSockaddrInet6(ip IP, port int, zone string) (syscall.SockaddrInet6, error) {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// In general, an IP wildcard address, which is either</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// &#34;0.0.0.0&#34; or &#34;::&#34;, means the entire IP addressing</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// space. For some historical reason, it is used to</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// specify &#34;any available address&#34; on some operations</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// of IP node.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// When the IP node supports IPv4-mapped IPv6 address,</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// we allow a listener to listen to the wildcard</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// address of both IP addressing spaces by specifying</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// IPv6 wildcard address.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	if len(ip) == 0 || ip.Equal(IPv4zero) {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		ip = IPv6zero
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// We accept any IPv6 address including IPv4-mapped</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">// IPv6 address.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	ip6 := ip.To16()
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	if ip6 == nil {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		return syscall.SockaddrInet6{}, &amp;AddrError{Err: &#34;non-IPv6 address&#34;, Addr: ip.String()}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	sa := syscall.SockaddrInet6{Port: port, ZoneId: uint32(zoneCache.index(zone))}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	copy(sa.Addr[:], ip6)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	return sa, nil
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>func ipToSockaddr(family int, ip IP, port int, zone string) (syscall.Sockaddr, error) {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	switch family {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	case syscall.AF_INET:
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		sa, err := ipToSockaddrInet4(ip, port)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		if err != nil {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			return nil, err
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		return &amp;sa, nil
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	case syscall.AF_INET6:
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		sa, err := ipToSockaddrInet6(ip, port, zone)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		if err != nil {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			return nil, err
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		return &amp;sa, nil
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	return nil, &amp;AddrError{Err: &#34;invalid address family&#34;, Addr: ip.String()}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>func addrPortToSockaddrInet4(ap netip.AddrPort) (syscall.SockaddrInet4, error) {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// ipToSockaddrInet4 has special handling here for zero length slices.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// We do not, because netip has no concept of a generic zero IP address.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	addr := ap.Addr()
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	if !addr.Is4() {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		return syscall.SockaddrInet4{}, &amp;AddrError{Err: &#34;non-IPv4 address&#34;, Addr: addr.String()}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	sa := syscall.SockaddrInet4{
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		Addr: addr.As4(),
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		Port: int(ap.Port()),
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	return sa, nil
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>func addrPortToSockaddrInet6(ap netip.AddrPort) (syscall.SockaddrInet6, error) {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// ipToSockaddrInet6 has special handling here for zero length slices.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// We do not, because netip has no concept of a generic zero IP address.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">// addr is allowed to be an IPv4 address, because As16 will convert it</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">// to an IPv4-mapped IPv6 address.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">// The error message is kept consistent with ipToSockaddrInet6.</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	addr := ap.Addr()
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	if !addr.IsValid() {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		return syscall.SockaddrInet6{}, &amp;AddrError{Err: &#34;non-IPv6 address&#34;, Addr: addr.String()}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	sa := syscall.SockaddrInet6{
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		Addr:   addr.As16(),
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		Port:   int(ap.Port()),
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		ZoneId: uint32(zoneCache.index(addr.Zone())),
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	return sa, nil
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
</pre><p><a href="ipsock_posix.go?m=text">View as plain text</a></p>

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
