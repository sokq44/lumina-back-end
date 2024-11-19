<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/cgo_unix.go - Go Documentation Server</title>

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
<a href="cgo_unix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<span class="text-muted">cgo_unix.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net">net</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file is called cgo_unix.go, but to allow syscalls-to-libc-based</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// implementations to share the code, it does not use cgo directly.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Instead of C.foo it uses _C_foo, which is defined in either</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// cgo_unix_cgo.go or cgo_unix_syscall.go</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//go:build !netgo &amp;&amp; ((cgo &amp;&amp; unix) || darwin)</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>package net
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>import (
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;context&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;net/netip&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;syscall&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;golang.org/x/net/dns/dnsmessage&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// cgoAvailable set to true to indicate that the cgo resolver</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// is available on this system.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>const cgoAvailable = true
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// An addrinfoErrno represents a getaddrinfo, getnameinfo-specific</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// error number. It&#39;s a signed number and a zero value is a non-error</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// by convention.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>type addrinfoErrno int
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>func (eai addrinfoErrno) Error() string   { return _C_gai_strerror(_C_int(eai)) }
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>func (eai addrinfoErrno) Temporary() bool { return eai == _C_EAI_AGAIN }
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>func (eai addrinfoErrno) Timeout() bool   { return false }
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// isAddrinfoErrno is just for testing purposes.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>func (eai addrinfoErrno) isAddrinfoErrno() {}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// doBlockingWithCtx executes a blocking function in a separate goroutine when the provided</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// context is cancellable. It is intended for use with calls that don&#39;t support context</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// cancellation (cgo, syscalls). blocking func may still be running after this function finishes.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func doBlockingWithCtx[T any](ctx context.Context, blocking func() (T, error)) (T, error) {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	if ctx.Done() == nil {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		return blocking()
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	type result struct {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		res T
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		err error
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	res := make(chan result, 1)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	go func() {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		var r result
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		r.res, r.err = blocking()
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		res &lt;- r
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	}()
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	select {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	case r := &lt;-res:
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		return r.res, r.err
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	case &lt;-ctx.Done():
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		var zero T
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		return zero, mapErr(ctx.Err())
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>func cgoLookupHost(ctx context.Context, name string) (hosts []string, err error) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	addrs, err := cgoLookupIP(ctx, &#34;ip&#34;, name)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if err != nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		return nil, err
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	for _, addr := range addrs {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		hosts = append(hosts, addr.String())
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	return hosts, nil
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>func cgoLookupPort(ctx context.Context, network, service string) (port int, err error) {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	var hints _C_struct_addrinfo
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	switch network {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	case &#34;ip&#34;: <span class="comment">// no hints</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	case &#34;tcp&#34;, &#34;tcp4&#34;, &#34;tcp6&#34;:
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		*_C_ai_socktype(&amp;hints) = _C_SOCK_STREAM
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		*_C_ai_protocol(&amp;hints) = _C_IPPROTO_TCP
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	case &#34;udp&#34;, &#34;udp4&#34;, &#34;udp6&#34;:
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		*_C_ai_socktype(&amp;hints) = _C_SOCK_DGRAM
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		*_C_ai_protocol(&amp;hints) = _C_IPPROTO_UDP
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	default:
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		return 0, &amp;DNSError{Err: &#34;unknown network&#34;, Name: network + &#34;/&#34; + service}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	switch ipVersion(network) {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	case &#39;4&#39;:
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		*_C_ai_family(&amp;hints) = _C_AF_INET
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	case &#39;6&#39;:
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		*_C_ai_family(&amp;hints) = _C_AF_INET6
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	return doBlockingWithCtx(ctx, func() (int, error) {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		return cgoLookupServicePort(&amp;hints, network, service)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	})
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func cgoLookupServicePort(hints *_C_struct_addrinfo, network, service string) (port int, err error) {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	cservice, err := syscall.ByteSliceFromString(service)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	if err != nil {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		return 0, &amp;DNSError{Err: err.Error(), Name: network + &#34;/&#34; + service}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// Lowercase the C service name.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	for i, b := range cservice[:len(service)] {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		cservice[i] = lowerASCII(b)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	var res *_C_struct_addrinfo
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	gerrno, err := _C_getaddrinfo(nil, (*_C_char)(unsafe.Pointer(&amp;cservice[0])), hints, &amp;res)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	if gerrno != 0 {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		isTemporary := false
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		switch gerrno {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		case _C_EAI_SYSTEM:
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			if err == nil { <span class="comment">// see golang.org/issue/6232</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>				err = syscall.EMFILE
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		case _C_EAI_SERVICE, _C_EAI_NONAME: <span class="comment">// Darwin returns EAI_NONAME.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			return 0, &amp;DNSError{Err: &#34;unknown port&#34;, Name: network + &#34;/&#34; + service, IsNotFound: true}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		default:
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			err = addrinfoErrno(gerrno)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			isTemporary = addrinfoErrno(gerrno).Temporary()
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		return 0, &amp;DNSError{Err: err.Error(), Name: network + &#34;/&#34; + service, IsTemporary: isTemporary}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	defer _C_freeaddrinfo(res)
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	for r := res; r != nil; r = *_C_ai_next(r) {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		switch *_C_ai_family(r) {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		case _C_AF_INET:
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			sa := (*syscall.RawSockaddrInet4)(unsafe.Pointer(*_C_ai_addr(r)))
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			p := (*[2]byte)(unsafe.Pointer(&amp;sa.Port))
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			return int(p[0])&lt;&lt;8 | int(p[1]), nil
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		case _C_AF_INET6:
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			sa := (*syscall.RawSockaddrInet6)(unsafe.Pointer(*_C_ai_addr(r)))
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			p := (*[2]byte)(unsafe.Pointer(&amp;sa.Port))
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			return int(p[0])&lt;&lt;8 | int(p[1]), nil
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	return 0, &amp;DNSError{Err: &#34;unknown port&#34;, Name: network + &#34;/&#34; + service, IsNotFound: true}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>func cgoLookupHostIP(network, name string) (addrs []IPAddr, err error) {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	acquireThread()
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	defer releaseThread()
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	var hints _C_struct_addrinfo
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	*_C_ai_flags(&amp;hints) = cgoAddrInfoFlags
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	*_C_ai_socktype(&amp;hints) = _C_SOCK_STREAM
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	*_C_ai_family(&amp;hints) = _C_AF_UNSPEC
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	switch ipVersion(network) {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	case &#39;4&#39;:
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		*_C_ai_family(&amp;hints) = _C_AF_INET
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	case &#39;6&#39;:
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		*_C_ai_family(&amp;hints) = _C_AF_INET6
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	h, err := syscall.BytePtrFromString(name)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	if err != nil {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		return nil, &amp;DNSError{Err: err.Error(), Name: name}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	var res *_C_struct_addrinfo
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	gerrno, err := _C_getaddrinfo((*_C_char)(unsafe.Pointer(h)), nil, &amp;hints, &amp;res)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	if gerrno != 0 {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		isErrorNoSuchHost := false
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		isTemporary := false
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		switch gerrno {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		case _C_EAI_SYSTEM:
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			if err == nil {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>				<span class="comment">// err should not be nil, but sometimes getaddrinfo returns</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>				<span class="comment">// gerrno == _C_EAI_SYSTEM with err == nil on Linux.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>				<span class="comment">// The report claims that it happens when we have too many</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>				<span class="comment">// open files, so use syscall.EMFILE (too many open files in system).</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>				<span class="comment">// Most system calls would return ENFILE (too many open files),</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>				<span class="comment">// so at the least EMFILE should be easy to recognize if this</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>				<span class="comment">// comes up again. golang.org/issue/6232.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>				err = syscall.EMFILE
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		case _C_EAI_NONAME, _C_EAI_NODATA:
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			err = errNoSuchHost
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			isErrorNoSuchHost = true
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		default:
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			err = addrinfoErrno(gerrno)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			isTemporary = addrinfoErrno(gerrno).Temporary()
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		return nil, &amp;DNSError{Err: err.Error(), Name: name, IsNotFound: isErrorNoSuchHost, IsTemporary: isTemporary}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	defer _C_freeaddrinfo(res)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	for r := res; r != nil; r = *_C_ai_next(r) {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		<span class="comment">// We only asked for SOCK_STREAM, but check anyhow.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		if *_C_ai_socktype(r) != _C_SOCK_STREAM {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			continue
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		switch *_C_ai_family(r) {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		case _C_AF_INET:
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			sa := (*syscall.RawSockaddrInet4)(unsafe.Pointer(*_C_ai_addr(r)))
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			addr := IPAddr{IP: copyIP(sa.Addr[:])}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			addrs = append(addrs, addr)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		case _C_AF_INET6:
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			sa := (*syscall.RawSockaddrInet6)(unsafe.Pointer(*_C_ai_addr(r)))
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			addr := IPAddr{IP: copyIP(sa.Addr[:]), Zone: zoneCache.name(int(sa.Scope_id))}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			addrs = append(addrs, addr)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	return addrs, nil
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>func cgoLookupIP(ctx context.Context, network, name string) (addrs []IPAddr, err error) {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	return doBlockingWithCtx(ctx, func() ([]IPAddr, error) {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		return cgoLookupHostIP(network, name)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	})
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// These are roughly enough for the following:</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">//	 Source		Encoding			Maximum length of single name entry</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">//	 Unicast DNS		ASCII or			&lt;=253 + a NUL terminator</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">//				Unicode in RFC 5892		252 * total number of labels + delimiters + a NUL terminator</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">//	 Multicast DNS	UTF-8 in RFC 5198 or		&lt;=253 + a NUL terminator</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">//				the same as unicast DNS ASCII	&lt;=253 + a NUL terminator</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">//	 Local database	various				depends on implementation</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>const (
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	nameinfoLen    = 64
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	maxNameinfoLen = 4096
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>func cgoLookupPTR(ctx context.Context, addr string) (names []string, err error) {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	ip, err := netip.ParseAddr(addr)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	if err != nil {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		return nil, &amp;DNSError{Err: &#34;invalid address&#34;, Name: addr}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	sa, salen := cgoSockaddr(IP(ip.AsSlice()), ip.Zone())
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	if sa == nil {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		return nil, &amp;DNSError{Err: &#34;invalid address &#34; + ip.String(), Name: addr}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	return doBlockingWithCtx(ctx, func() ([]string, error) {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		return cgoLookupAddrPTR(addr, sa, salen)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	})
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>func cgoLookupAddrPTR(addr string, sa *_C_struct_sockaddr, salen _C_socklen_t) (names []string, err error) {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	acquireThread()
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	defer releaseThread()
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	var gerrno int
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	var b []byte
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	for l := nameinfoLen; l &lt;= maxNameinfoLen; l *= 2 {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		b = make([]byte, l)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		gerrno, err = cgoNameinfoPTR(b, sa, salen)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		if gerrno == 0 || gerrno != _C_EAI_OVERFLOW {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			break
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	if gerrno != 0 {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		isErrorNoSuchHost := false
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		isTemporary := false
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		switch gerrno {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		case _C_EAI_SYSTEM:
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			if err == nil { <span class="comment">// see golang.org/issue/6232</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				err = syscall.EMFILE
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		case _C_EAI_NONAME:
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			err = errNoSuchHost
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			isErrorNoSuchHost = true
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		default:
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			err = addrinfoErrno(gerrno)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			isTemporary = addrinfoErrno(gerrno).Temporary()
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		return nil, &amp;DNSError{Err: err.Error(), Name: addr, IsTemporary: isTemporary, IsNotFound: isErrorNoSuchHost}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	for i := 0; i &lt; len(b); i++ {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		if b[i] == 0 {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			b = b[:i]
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			break
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	return []string{absDomainName(string(b))}, nil
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>func cgoSockaddr(ip IP, zone string) (*_C_struct_sockaddr, _C_socklen_t) {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	if ip4 := ip.To4(); ip4 != nil {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		return cgoSockaddrInet4(ip4), _C_socklen_t(syscall.SizeofSockaddrInet4)
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	if ip6 := ip.To16(); ip6 != nil {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		return cgoSockaddrInet6(ip6, zoneCache.index(zone)), _C_socklen_t(syscall.SizeofSockaddrInet6)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	return nil, 0
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>func cgoLookupCNAME(ctx context.Context, name string) (cname string, err error, completed bool) {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	resources, err := resSearch(ctx, name, int(dnsmessage.TypeCNAME), int(dnsmessage.ClassINET))
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	if err != nil {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		return
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	cname, err = parseCNAMEFromResources(resources)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	if err != nil {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		return &#34;&#34;, err, false
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	return cname, nil, true
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// resSearch will make a call to the &#39;res_nsearch&#39; routine in the C library</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// and parse the output as a slice of DNS resources.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>func resSearch(ctx context.Context, hostname string, rtype, class int) ([]dnsmessage.Resource, error) {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	return doBlockingWithCtx(ctx, func() ([]dnsmessage.Resource, error) {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		return cgoResSearch(hostname, rtype, class)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	})
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>func cgoResSearch(hostname string, rtype, class int) ([]dnsmessage.Resource, error) {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	acquireThread()
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	defer releaseThread()
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	resStateSize := unsafe.Sizeof(_C_struct___res_state{})
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	var state *_C_struct___res_state
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	if resStateSize &gt; 0 {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		mem := _C_malloc(resStateSize)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		defer _C_free(mem)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		memSlice := unsafe.Slice((*byte)(mem), resStateSize)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		clear(memSlice)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		state = (*_C_struct___res_state)(unsafe.Pointer(&amp;memSlice[0]))
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	if err := _C_res_ninit(state); err != nil {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		return nil, errors.New(&#34;res_ninit failure: &#34; + err.Error())
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	defer _C_res_nclose(state)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// Some res_nsearch implementations (like macOS) do not set errno.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// They set h_errno, which is not per-thread and useless to us.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// res_nsearch returns the size of the DNS response packet.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// But if the DNS response packet contains failure-like response codes,</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">// res_search returns -1 even though it has copied the packet into buf,</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	<span class="comment">// giving us no way to find out how big the packet is.</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	<span class="comment">// For now, we are willing to take res_search&#39;s word that there&#39;s nothing</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	<span class="comment">// useful in the response, even though there *is* a response.</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	bufSize := maxDNSPacketSize
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	buf := (*_C_uchar)(_C_malloc(uintptr(bufSize)))
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	defer _C_free(unsafe.Pointer(buf))
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	s, err := syscall.BytePtrFromString(hostname)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	if err != nil {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		return nil, err
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	var size int
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	for {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		size, _ = _C_res_nsearch(state, (*_C_char)(unsafe.Pointer(s)), class, rtype, buf, bufSize)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		if size &lt;= 0 || size &gt; 0xffff {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			return nil, errors.New(&#34;res_nsearch failure&#34;)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		if size &lt;= bufSize {
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			break
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		<span class="comment">// Allocate a bigger buffer to fit the entire msg.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		_C_free(unsafe.Pointer(buf))
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		bufSize = size
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		buf = (*_C_uchar)(_C_malloc(uintptr(bufSize)))
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	var p dnsmessage.Parser
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	if _, err := p.Start(unsafe.Slice((*byte)(unsafe.Pointer(buf)), size)); err != nil {
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		return nil, err
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	p.SkipAllQuestions()
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	resources, err := p.AllAnswers()
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	if err != nil {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		return nil, err
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	return resources, nil
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
</pre><p><a href="cgo_unix.go?m=text">View as plain text</a></p>

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
