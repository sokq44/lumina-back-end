<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/sock_posix.go - Go Documentation Server</title>

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
<a href="sock_posix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<span class="text-muted">sock_posix.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build unix || windows</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package net
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;context&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/poll&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;os&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;syscall&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// socket returns a network file descriptor that is ready for</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// asynchronous I/O using the network poller.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>func socket(ctx context.Context, net string, family, sotype, proto int, ipv6only bool, laddr, raddr sockaddr, ctrlCtxFn func(context.Context, string, string, syscall.RawConn) error) (fd *netFD, err error) {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	s, err := sysSocket(family, sotype, proto)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	if err != nil {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>		return nil, err
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	if err = setDefaultSockopts(s, family, sotype, ipv6only); err != nil {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>		poll.CloseFunc(s)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		return nil, err
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	}
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	if fd, err = newFD(s, family, sotype, net); err != nil {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		poll.CloseFunc(s)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		return nil, err
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// This function makes a network file descriptor for the</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// following applications:</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// - An endpoint holder that opens a passive stream</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">//   connection, known as a stream listener</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// - An endpoint holder that opens a destination-unspecific</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">//   datagram connection, known as a datagram listener</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// - An endpoint holder that opens an active stream or a</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">//   destination-specific datagram connection, known as a</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">//   dialer</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// - An endpoint holder that opens the other connection, such</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">//   as talking to the protocol stack inside the kernel</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// For stream and datagram listeners, they will only require</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// named sockets, so we can assume that it&#39;s just a request</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// from stream or datagram listeners when laddr is not nil but</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// raddr is nil. Otherwise we assume it&#39;s just for dialers or</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// the other connection holders.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	if laddr != nil &amp;&amp; raddr == nil {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		switch sotype {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		case syscall.SOCK_STREAM, syscall.SOCK_SEQPACKET:
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			if err := fd.listenStream(ctx, laddr, listenerBacklog(), ctrlCtxFn); err != nil {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>				fd.Close()
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>				return nil, err
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			return fd, nil
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		case syscall.SOCK_DGRAM:
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			if err := fd.listenDatagram(ctx, laddr, ctrlCtxFn); err != nil {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>				fd.Close()
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>				return nil, err
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			return fd, nil
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	if err := fd.dial(ctx, laddr, raddr, ctrlCtxFn); err != nil {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		fd.Close()
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		return nil, err
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return fd, nil
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>func (fd *netFD) ctrlNetwork() string {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	switch fd.net {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	case &#34;unix&#34;, &#34;unixgram&#34;, &#34;unixpacket&#34;:
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		return fd.net
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	switch fd.net[len(fd.net)-1] {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	case &#39;4&#39;, &#39;6&#39;:
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		return fd.net
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	if fd.family == syscall.AF_INET {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		return fd.net + &#34;4&#34;
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	return fd.net + &#34;6&#34;
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>func (fd *netFD) dial(ctx context.Context, laddr, raddr sockaddr, ctrlCtxFn func(context.Context, string, string, syscall.RawConn) error) error {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	var c *rawConn
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	if ctrlCtxFn != nil {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		c = newRawConn(fd)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		var ctrlAddr string
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		if raddr != nil {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			ctrlAddr = raddr.String()
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		} else if laddr != nil {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			ctrlAddr = laddr.String()
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		if err := ctrlCtxFn(ctx, fd.ctrlNetwork(), ctrlAddr, c); err != nil {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			return err
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	var lsa syscall.Sockaddr
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	var err error
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	if laddr != nil {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		if lsa, err = laddr.sockaddr(fd.family); err != nil {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			return err
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		} else if lsa != nil {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			if err = syscall.Bind(fd.pfd.Sysfd, lsa); err != nil {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>				return os.NewSyscallError(&#34;bind&#34;, err)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	var rsa syscall.Sockaddr  <span class="comment">// remote address from the user</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	var crsa syscall.Sockaddr <span class="comment">// remote address we actually connected to</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	if raddr != nil {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		if rsa, err = raddr.sockaddr(fd.family); err != nil {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			return err
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		if crsa, err = fd.connect(ctx, lsa, rsa); err != nil {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			return err
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		fd.isConnected = true
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	} else {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		if err := fd.init(); err != nil {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			return err
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">// Record the local and remote addresses from the actual socket.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// Get the local address by calling Getsockname.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// For the remote address, use</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// 1) the one returned by the connect method, if any; or</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// 2) the one from Getpeername, if it succeeds; or</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// 3) the one passed to us as the raddr parameter.</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	lsa, _ = syscall.Getsockname(fd.pfd.Sysfd)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	if crsa != nil {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		fd.setAddr(fd.addrFunc()(lsa), fd.addrFunc()(crsa))
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	} else if rsa, _ = syscall.Getpeername(fd.pfd.Sysfd); rsa != nil {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		fd.setAddr(fd.addrFunc()(lsa), fd.addrFunc()(rsa))
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	} else {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		fd.setAddr(fd.addrFunc()(lsa), raddr)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	return nil
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>func (fd *netFD) listenStream(ctx context.Context, laddr sockaddr, backlog int, ctrlCtxFn func(context.Context, string, string, syscall.RawConn) error) error {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	var err error
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	if err = setDefaultListenerSockopts(fd.pfd.Sysfd); err != nil {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		return err
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	var lsa syscall.Sockaddr
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if lsa, err = laddr.sockaddr(fd.family); err != nil {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return err
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	if ctrlCtxFn != nil {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		c := newRawConn(fd)
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		if err := ctrlCtxFn(ctx, fd.ctrlNetwork(), laddr.String(), c); err != nil {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			return err
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	if err = syscall.Bind(fd.pfd.Sysfd, lsa); err != nil {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		return os.NewSyscallError(&#34;bind&#34;, err)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if err = listenFunc(fd.pfd.Sysfd, backlog); err != nil {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		return os.NewSyscallError(&#34;listen&#34;, err)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	if err = fd.init(); err != nil {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		return err
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	lsa, _ = syscall.Getsockname(fd.pfd.Sysfd)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	fd.setAddr(fd.addrFunc()(lsa), nil)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	return nil
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>func (fd *netFD) listenDatagram(ctx context.Context, laddr sockaddr, ctrlCtxFn func(context.Context, string, string, syscall.RawConn) error) error {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	switch addr := laddr.(type) {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	case *UDPAddr:
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		<span class="comment">// We provide a socket that listens to a wildcard</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		<span class="comment">// address with reusable UDP port when the given laddr</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		<span class="comment">// is an appropriate UDP multicast address prefix.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		<span class="comment">// This makes it possible for a single UDP listener to</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		<span class="comment">// join multiple different group addresses, for</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		<span class="comment">// multiple UDP listeners that listen on the same UDP</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		<span class="comment">// port to join the same group address.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		if addr.IP != nil &amp;&amp; addr.IP.IsMulticast() {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			if err := setDefaultMulticastSockopts(fd.pfd.Sysfd); err != nil {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>				return err
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			addr := *addr
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			switch fd.family {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			case syscall.AF_INET:
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>				addr.IP = IPv4zero
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			case syscall.AF_INET6:
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>				addr.IP = IPv6unspecified
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>			laddr = &amp;addr
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	var err error
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	var lsa syscall.Sockaddr
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	if lsa, err = laddr.sockaddr(fd.family); err != nil {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		return err
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if ctrlCtxFn != nil {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		c := newRawConn(fd)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		if err := ctrlCtxFn(ctx, fd.ctrlNetwork(), laddr.String(), c); err != nil {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			return err
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	if err = syscall.Bind(fd.pfd.Sysfd, lsa); err != nil {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		return os.NewSyscallError(&#34;bind&#34;, err)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	if err = fd.init(); err != nil {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		return err
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	lsa, _ = syscall.Getsockname(fd.pfd.Sysfd)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	fd.setAddr(fd.addrFunc()(lsa), nil)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	return nil
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
</pre><p><a href="sock_posix.go?m=text">View as plain text</a></p>

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
